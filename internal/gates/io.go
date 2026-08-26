// ABOUTME: Fence-safe canonical v1 gate reads and atomic gates-subtree writes.
// ABOUTME: Writes preserve every byte outside the binary-owned gates mapping.
package gates

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/gitsource"
	"gopkg.in/yaml.v3"
)

// ErrNoGateRecord marks an entity with no gates record at all. Callers may
// take a gate-free (legacy) path only on this exact condition; any other read
// failure is hostile authority and must be refused byte-clean.
var ErrNoGateRecord = errors.New("entity has no gates record")

// Warning is a non-authoritative compatibility finding from a gate read.
// Path is the stable gates-node path of the application mapping and Field is
// the unknown key observed there. Warnings are sorted and de-duplicated by
// ReadDiagnostics; they never participate in eligibility or writes.
type Warning struct {
	Path  string `json:"path"`
	Field string `json:"field"`
}

func Read(path string) (*Document, *yaml.Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	doc, node, _, err := readDataDiagnostics(data)
	return doc, node, err
}

// ReadDiagnostics reads a gate document and returns compatibility warnings in
// addition to the original gates node. The node is never filtered or replaced;
// callers may safely use it as the compare-and-swap/write expectation. Unknown
// keys are tolerated only in the exact application mappings described by the
// v1 gate contract.
func ReadDiagnostics(path string) (*Document, *yaml.Node, []Warning, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, nil, err
	}
	return readDataDiagnostics(data)
}

func readData(data []byte) (*Document, *yaml.Node, error) {
	doc, node, _, err := readDataDiagnostics(data)
	return doc, node, err
}

func readDataDiagnostics(data []byte) (*Document, *yaml.Node, []Warning, error) {
	root, _, _, err := frontmatterNode(data)
	if err != nil {
		return nil, nil, nil, err
	}
	gatesNode := mappingValue(root, "gates")
	if gatesNode == nil {
		return nil, nil, nil, ErrNoGateRecord
	}
	filtered := cloneYAMLNode(gatesNode)
	warnings, err := filterApplicationMappings(filtered)
	if err != nil {
		return nil, nil, nil, err
	}
	encoded, err := yaml.Marshal(filtered)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("encode gates for validation: %w", err)
	}
	var doc Document
	decoder := yaml.NewDecoder(bytes.NewReader(encoded))
	decoder.KnownFields(true)
	if err := decoder.Decode(&doc); err != nil {
		return nil, nil, nil, fmt.Errorf("decode canonical gates v1: %w", err)
	}
	if err := Validate(&doc); err != nil {
		return nil, nil, nil, err
	}
	return &doc, gatesNode, warnings, nil
}

// cloneYAMLNode clones a node through YAML encoding, which retains scalar tags
// and aliases without ever sharing Content slices with the source node. The
// clone is deliberately used only for strict validation; the source node stays
// untouched for CAS and byte-preserving writes.
func cloneYAMLNode(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}
	encoded, err := yaml.Marshal(node)
	if err != nil {
		// A node obtained from yaml.Unmarshal is marshalable. Keep this helper
		// total for package-internal callers; the caller's subsequent strict
		// decode will report the meaningful failure if a future node shape is not.
		return &yaml.Node{}
	}
	var out yaml.Node
	if err := yaml.Unmarshal(encoded, &out); err != nil || len(out.Content) == 0 {
		return &yaml.Node{}
	}
	if out.Kind == yaml.DocumentNode {
		return out.Content[0]
	}
	return &out
}

// retiredAttemptKey is the one attempt-level key a read drops instead of
// refusing. Its writer was cut with provider-backed closure, so no live command
// produces it, but frozen archived attempts still carry the bytes. It is
// retired rather than unknown: dropping it silently keeps those records
// readable without widening attempt-level tolerance to any other key.
const retiredAttemptKey = "provider-evidence"

// filterApplicationMappings drops the retired attempt-level key and removes
// non-canonical keys only at gates.records[*].attempts[*].application. It
// validates the application value shape before filtering so null, sequence, and
// scalar legacy values remain strict errors. Unknown keys are reported once per
// path/field and sorted for deterministic operator output; the retired key is
// silent and raises no warning.
func filterApplicationMappings(gatesNode *yaml.Node) ([]Warning, error) {
	var warnings []Warning
	if gatesNode == nil || gatesNode.Kind != yaml.MappingNode {
		return warnings, nil
	}
	records := mappingValue(gatesNode, "records")
	if records == nil || records.Kind != yaml.SequenceNode {
		return warnings, nil
	}
	for ri, record := range records.Content {
		if record == nil || record.Kind != yaml.MappingNode {
			continue
		}
		attempts := mappingValue(record, "attempts")
		if attempts == nil || attempts.Kind != yaml.SequenceNode {
			continue
		}
		for ai, attempt := range attempts.Content {
			if attempt == nil || attempt.Kind != yaml.MappingNode {
				continue
			}
			keptAttempt := make([]*yaml.Node, 0, len(attempt.Content))
			for i := 0; i+1 < len(attempt.Content); i += 2 {
				if attempt.Content[i].Value == retiredAttemptKey {
					continue
				}
				keptAttempt = append(keptAttempt, attempt.Content[i], attempt.Content[i+1])
			}
			attempt.Content = keptAttempt
			for i := 0; i+1 < len(attempt.Content); i += 2 {
				if attempt.Content[i].Value != "application" {
					continue
				}
				application := attempt.Content[i+1]
				if application == nil || application.Kind != yaml.MappingNode {
					return nil, fmt.Errorf("gates.records[%d].attempts[%d].application must be a mapping", ri, ai)
				}
				path := fmt.Sprintf("gates.records[%d].attempts[%d].application", ri, ai)
				kept := make([]*yaml.Node, 0, len(application.Content))
				seen := make(map[string]bool)
				for j := 0; j+1 < len(application.Content); j += 2 {
					key, value := application.Content[j], application.Content[j+1]
					if key.Value == "target-stage" || key.Value == "state" {
						kept = append(kept, key, value)
						continue
					}
					if !seen[key.Value] {
						warnings = append(warnings, Warning{Path: path, Field: key.Value})
						seen[key.Value] = true
					}
				}
				application.Content = kept
			}
		}
	}
	sort.Slice(warnings, func(i, j int) bool {
		if warnings[i].Path != warnings[j].Path {
			return warnings[i].Path < warnings[j].Path
		}
		return warnings[i].Field < warnings[j].Field
	})
	return warnings, nil
}

func validateRetainedAuthority(entityPath, workflowDir string, doc *Document) error {
	return validateRetainedAuthorityExcept(entityPath, workflowDir, doc, "", "")
}

func validateRetainedAuthorityExcept(entityPath, workflowDir string, doc *Document, skipGate, skipAttempt string) error {
	roots := gitsource.Roots{Main: workflowDir, State: filepath.Dir(entityPath)}
	for ri := range doc.Records {
		record := &doc.Records[ri]
		for ai := range record.Attempts {
			attempt := &record.Attempts[ai]
			if record.ID == skipGate && attempt.ID == skipAttempt {
				continue
			}
			// A prepared room retains authority whether or not it has a
			// request. The request checks below belong to a retained
			// request-backed binding alone. The Briefing digest and the
			// selected Git objects are what every gate presents, so every
			// prepared room gets those two checks. A skip here gives the
			// one-file room less validation than the two-file room had, and
			// that inverts the point of the change.
			if !preparedRoomBinding(entityPath, attempt.Briefing) {
				continue
			}
			if attempt.Briefing.RequestDigest != "" {
				room := filepath.Join(filepath.Dir(entityPath), filepath.FromSlash(attempt.Briefing.RoomRef))
				requestBytes, err := os.ReadFile(filepath.Join(room, "request.json"))
				if err != nil {
					return fmt.Errorf("attempt %s retained request.json: %w", attempt.ID, err)
				}
				requestDigest, err := CanonicalDigest(requestBytes)
				if err != nil || requestDigest != attempt.Briefing.RequestDigest {
					return fmt.Errorf("attempt %s retained request.json does not match its frozen digest", attempt.ID)
				}
				request, err := decodeGateRoomRequest(requestBytes)
				if err != nil {
					return err
				}
				if request.Type != "spacedock-gate-presentation-request" || request.Version != "1" ||
					request.Gate != record.ID || request.Attempt != attempt.ID ||
					request.Briefing.ID != attempt.Briefing.ID || request.Briefing.Digest != attempt.Briefing.Digest ||
					request.Actor != "person:captain" || request.Approver != "person:captain" {
					return fmt.Errorf("attempt %s request does not bind its gate, attempt, Briefing, and captain authority", attempt.ID)
				}
			}
			manifest, err := boundBriefingManifest(entityPath, attempt.Briefing)
			if err != nil {
				return err
			}
			items, err := canonicalPresentationItems(manifest)
			if err != nil {
				return err
			}
			if err := validatePresentationGitSources(roots, items); err != nil {
				return fmt.Errorf("attempt %s %w", attempt.ID, err)
			}
		}
	}
	return nil
}

func validatePresentationGitSources(roots gitsource.Roots, items []presentedItem) error {
	gitItems := 0
	for _, item := range items {
		if strings.HasPrefix(item.URI, "git-root://") {
			gitItems++
			if _, err := gitsource.Resolve(roots, item.URI, item.Rev); err != nil {
				return fmt.Errorf("selected source: %w", err)
			}
		}
	}
	if gitItems != 0 && gitItems != len(items) {
		return fmt.Errorf("mixes Git-root and non-Git selected source identities")
	}
	return nil
}

func entityStatus(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	root, _, _, err := frontmatterNode(data)
	if err != nil {
		return "", err
	}
	status := mappingValue(root, "status")
	if status == nil || strings.TrimSpace(status.Value) == "" {
		return "", fmt.Errorf("entity has no current workflow status")
	}
	return status.Value, nil
}

func writeDocument(path string, expected *yaml.Node, doc *Document) error {
	return writeEntityDocument(path, expected, nil, doc, nil)
}

// entityField is one extra top-level frontmatter scalar written in the SAME
// locked candidate replacement as the gates/status swap. Its use is confined to
// the terminal delivery write: verdict and completed land with the authority
// spend so done-but-undelivered is unrepresentable. fields are inserted after
// the last frontmatter line when absent (insert: true), never byte-shifting
// existing fields.
type entityField struct {
	key   string
	value string
}

// entityDocumentWriteFn is the ONE atomic candidate replacement at the gates
// writer's write site (atomicWrite in production). Declared as a var so a
// package-level assertion can observe that the terminal delivery write's
// candidate carries all four field changes (application state, status,
// verdict, completed) at once, and can fail the replacement mid-write to
// prove the original bytes stay intact.
var entityDocumentWriteFn = atomicWrite

func writeEntityDocument(path string, expected *yaml.Node, expectedStatus *string, doc *Document, status *string, fields ...entityField) error {
	if err := Validate(doc); err != nil {
		return err
	}
	if expected == nil {
		expected = &yaml.Node{}
	}
	gatesBlock, err := yaml.Marshal(struct {
		Gates *Document `yaml:"gates"`
	}{Gates: doc})
	if err != nil {
		return err
	}
	replacements := []topLevelReplacement{{key: "gates", data: gatesBlock, insert: status == nil}}
	if status != nil {
		statusBlock, err := yaml.Marshal(struct {
			Status string `yaml:"status"`
		}{*status})
		if err != nil {
			return err
		}
		replacements = append(replacements, topLevelReplacement{key: "status", data: statusBlock})
	}
	for _, f := range fields {
		if strings.ContainsAny(f.key, ":\r\n") || strings.ContainsAny(f.value, "\r\n") {
			return fmt.Errorf("entity field %q is not a plain one-line scalar", f.key)
		}
		// Emit the same plain-scalar line runSet writes (e.g. completed:
		// 2026-08-01T00:00:00Z, not the quoted form yaml.Marshal picks for
		// timestamp-shaped strings), so the locked write is byte-identical in
		// shape to the --set it replaces.
		replacements = append(replacements, topLevelReplacement{key: f.key, data: []byte(f.key + ": " + f.value + "\n"), insert: true})
	}
	return mutateEntity(path, entityExpectation{Gates: expected, Status: expectedStatus}, func(original []byte) ([]byte, error) {
		out, err := replaceTopLevels(original, replacements...)
		if err != nil {
			return nil, err
		}
		if _, _, err := readData(out); err != nil {
			return nil, fmt.Errorf("validate rebuilt gates: %w", err)
		}
		if status != nil || len(fields) > 0 {
			parsed, _, _, err := frontmatterNode(out)
			if err != nil {
				return nil, fmt.Errorf("validate rebuilt workflow status")
			}
			if status != nil && (mappingValue(parsed, "status") == nil || mappingValue(parsed, "status").Value != *status) {
				return nil, fmt.Errorf("validate rebuilt workflow status")
			}
			for _, f := range fields {
				if node := mappingValue(parsed, f.key); node == nil || node.Value != f.value {
					return nil, fmt.Errorf("validate rebuilt %s field", f.key)
				}
			}
		}
		return out, nil
	}, entityDocumentWriteFn)
}

type entityExpectation struct {
	Bytes  []byte
	Gates  *yaml.Node
	Status *string
}

func mutateEntity(path string, expected entityExpectation, build func([]byte) ([]byte, error), replace func(string, []byte) error) error {
	current, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if (expected.Bytes == nil) == (expected.Gates == nil) || expected.Status != nil && expected.Gates == nil {
		return fmt.Errorf("entity mutation requires exactly one expectation mode")
	}
	if expected.Bytes != nil {
		if !bytes.Equal(expected.Bytes, current) {
			return fmt.Errorf("entity changed during locked update")
		}
	} else {
		root, _, _, err := frontmatterNode(current)
		if err != nil {
			return err
		}
		wantGates := expected.Gates
		if wantGates.Kind == 0 {
			wantGates = nil
		}
		if !sameYAMLNode(wantGates, mappingValue(root, "gates")) {
			return fmt.Errorf("gates record changed during locked update")
		}
		if expected.Status != nil {
			status := mappingValue(root, "status")
			if status == nil || status.Value != *expected.Status {
				return fmt.Errorf("workflow status changed during locked update")
			}
		}
	}
	next, err := build(current)
	if err != nil || bytes.Equal(next, current) {
		return err
	}
	return replace(path, next)
}

type topLevelReplacement struct {
	key        string
	start, end int
	data       []byte
	insert     bool
}

func replaceTopLevels(original []byte, replacements ...topLevelReplacement) ([]byte, error) {
	root, fmStart, fmEnd, err := frontmatterNode(original)
	if err != nil {
		return nil, err
	}
	for i := range replacements {
		start, end, ok := topLevelRange(root, fmStart, fmEnd, replacements[i].key)
		if !ok {
			if !replacements[i].insert {
				return nil, fmt.Errorf("entity has no %s field", replacements[i].key)
			}
			start, end = fmEnd, fmEnd
		}
		replacements[i].start, replacements[i].end = lineOffset(original, start), lineOffset(original, end)
		if bytes.Contains(original, []byte("\r\n")) {
			replacements[i].data = []byte(strings.ReplaceAll(string(replacements[i].data), "\n", "\r\n"))
		}
	}
	// Later positions replace first so earlier offsets stay valid; equal
	// offsets (insertions at the frontmatter end) apply last-listed first so
	// the final buffer lists them in the order the caller gave.
	sort.SliceStable(replacements, func(i, j int) bool {
		if replacements[i].start == replacements[j].start {
			return i > j
		}
		return replacements[i].start > replacements[j].start
	})
	out := append([]byte(nil), original...)
	for _, replacement := range replacements {
		out = append(append(append([]byte{}, out[:replacement.start]...), replacement.data...), out[replacement.end:]...)
	}
	return out, nil
}

func topLevelRange(root *yaml.Node, fmStart, fmEnd int, key string) (int, int, bool) {
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value == key {
			end := fmEnd
			if i+2 < len(root.Content) {
				end = fmStart + root.Content[i+2].Line
			}
			return fmStart + root.Content[i].Line, end, true
		}
	}
	return 0, 0, false
}

func sameYAMLNode(left, right *yaml.Node) bool {
	if left == nil || right == nil {
		return left == right
	}
	lb, lerr := yaml.Marshal(left)
	rb, rerr := yaml.Marshal(right)
	return lerr == nil && rerr == nil && bytes.Equal(lb, rb)
}

func normalizedLines(data []byte) []string {
	return strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
}

func lineOffset(data []byte, line int) int {
	if line <= 0 {
		return 0
	}
	seen := 0
	for i, b := range data {
		if b != '\n' {
			continue
		}
		seen++
		if seen == line {
			return i + 1
		}
	}
	return len(data)
}

func frontmatterNode(data []byte) (*yaml.Node, int, int, error) {
	lines := normalizedLines(data)
	fmStart, fmEnd := -1, -1
	for i, line := range lines {
		if strings.TrimSpace(line) != "---" {
			continue
		}
		if fmStart < 0 {
			fmStart = i
		} else {
			fmEnd = i
			break
		}
	}
	if fmStart < 0 || fmEnd < 0 {
		return nil, 0, 0, fmt.Errorf("entity has no complete frontmatter")
	}
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(strings.Join(lines[fmStart+1:fmEnd], "\n")+"\n"), &doc); err != nil {
		return nil, 0, 0, fmt.Errorf("parse frontmatter: %w", err)
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, 0, 0, fmt.Errorf("frontmatter must be a mapping")
	}
	return doc.Content[0], fmStart, fmEnd, nil
}

func mappingValue(m *yaml.Node, key string) *yaml.Node {
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			return m.Content[i+1]
		}
	}
	return nil
}

func atomicWrite(path string, data []byte) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".gates-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if err := tmp.Chmod(info.Mode()); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(name, path)
}

type roundRoomBytes struct {
	Exists        bool
	Briefing, Log []byte
}

func readRoundRoom(room string) (roundRoomBytes, error) {
	info, err := os.Lstat(room)
	if os.IsNotExist(err) {
		return roundRoomBytes{}, nil
	}
	if err != nil || !info.IsDir() {
		return roundRoomBytes{}, fmt.Errorf("round target is occupied")
	}
	entries, err := os.ReadDir(room)
	if err != nil || len(entries) != 2 {
		return roundRoomBytes{}, fmt.Errorf("round target is not a canonical room")
	}
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			return roundRoomBytes{}, fmt.Errorf("round target is not a canonical room")
		}
	}
	result := roundRoomBytes{Exists: true}
	result.Briefing, err = os.ReadFile(filepath.Join(room, "briefing.json"))
	if err == nil {
		result.Log, err = os.ReadFile(filepath.Join(room, "briefing.review.jsonl"))
	}
	return result, err
}

func publishRound(room string, next roundRoomBytes, commitEntity func(bool) error) error {
	current, err := readRoundRoom(room)
	if err != nil {
		return err
	}
	if !next.Exists {
		return fmt.Errorf("round publication requires retained bytes")
	}
	if current.Exists {
		if !bytes.Equal(current.Briefing, next.Briefing) || !bytes.Equal(current.Log, next.Log) {
			return fmt.Errorf("round room is immutable and differs from supplied bytes")
		}
		return commitEntity(true)
	}
	parent := filepath.Dir(room)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	cleanup := func() {
		_ = os.Remove(parent)
		_ = os.Remove(filepath.Dir(parent))
	}
	tmp, err := os.MkdirTemp(parent, ".round-*")
	if err != nil {
		cleanup()
		return err
	}
	defer os.RemoveAll(tmp)
	defer cleanup()
	for name, data := range map[string][]byte{"briefing.json": next.Briefing, "briefing.review.jsonl": next.Log} {
		if err = writeSyncedFile(filepath.Join(tmp, name), data); err != nil {
			return err
		}
	}
	if err = os.Rename(tmp, room); err != nil {
		return err
	}
	if err := commitEntity(false); err != nil {
		_ = os.RemoveAll(room)
		return err
	}
	return nil
}

func writeSyncedFile(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err = file.Write(data); err == nil {
		err = file.Sync()
	}
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	return err
}
