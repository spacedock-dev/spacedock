// ABOUTME: Fence-safe canonical v1 gate reads and atomic gates-subtree writes.
// ABOUTME: Writes preserve every byte outside the binary-owned gates mapping.
package gates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/gitsource"
	"gopkg.in/yaml.v3"
)

func Read(path string) (*Document, *yaml.Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	return readData(data)
}

func readData(data []byte) (*Document, *yaml.Node, error) {
	root, _, _, err := frontmatterNode(data)
	if err != nil {
		return nil, nil, err
	}
	gatesNode := mappingValue(root, "gates")
	if gatesNode == nil {
		return nil, nil, fmt.Errorf("entity has no gates record")
	}
	encoded, err := yaml.Marshal(gatesNode)
	if err != nil {
		return nil, nil, fmt.Errorf("encode gates for validation: %w", err)
	}
	var doc Document
	decoder := yaml.NewDecoder(bytes.NewReader(encoded))
	decoder.KnownFields(true)
	if err := decoder.Decode(&doc); err != nil {
		return nil, nil, fmt.Errorf("decode canonical gates v1: %w", err)
	}
	if err := Validate(&doc); err != nil {
		return nil, nil, err
	}
	return &doc, gatesNode, nil
}

func SummaryFile(path string) (Summary, error) {
	return SummaryFileAt(path, nearestWorkflowDir(filepath.Dir(path)))
}

func SummaryFileAt(path, workflowDir string) (Summary, error) {
	doc, _, err := Read(path)
	if err != nil {
		return Summary{}, err
	}
	if err := validateRetainedAuthority(path, workflowDir, doc); err != nil {
		return Summary{}, err
	}
	return CurrentSummary(doc), nil
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
			if attempt.Briefing.RequestDigest == "" {
				continue
			}
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
			manifest, err := boundBriefingManifest(entityPath, attempt.Briefing)
			if err != nil {
				return err
			}
			items, err := canonicalPresentationItems(manifest)
			if err != nil {
				return err
			}
			gitItems := 0
			for _, item := range items {
				if strings.HasPrefix(item.URI, "git-root://") {
					gitItems++
					if _, err := gitsource.Resolve(roots, item.URI, item.Rev); err != nil {
						return fmt.Errorf("attempt %s selected source: %w", attempt.ID, err)
					}
				}
			}
			if gitItems != 0 && gitItems != len(items) {
				return fmt.Errorf("attempt %s mixes Git-root and non-Git selected source identities", attempt.ID)
			}
			if attempt.ProviderEvidence == nil {
				continue
			}
			resultBytes, err := os.ReadFile(filepath.Join(room, "provider", "result.json"))
			if err != nil {
				return fmt.Errorf("attempt %s retained provider/result.json: %w", attempt.ID, err)
			}
			inventoryBytes, err := os.ReadFile(filepath.Join(room, "provider", "presented-inventory.json"))
			if err != nil {
				return fmt.Errorf("attempt %s retained provider/presented-inventory.json: %w", attempt.ID, err)
			}
			if RawDigest(resultBytes) != attempt.ProviderEvidence.ResultDigest {
				return fmt.Errorf("attempt %s retained provider/result.json does not match its frozen digest", attempt.ID)
			}
			if RawDigest(inventoryBytes) != attempt.ProviderEvidence.PresentedInventoryDigest {
				return fmt.Errorf("attempt %s retained provider/presented-inventory.json does not match its frozen digest", attempt.ID)
			}
			result, err := decodeProviderResult(resultBytes)
			if err != nil {
				return err
			}
			inventory, err := decodePresentedInventory(inventoryBytes)
			if err != nil {
				return err
			}
			association, err := deriveAssociation(resultBytes, result, inventory, request.Approver, attempt.Briefing, items)
			if err != nil {
				return err
			}
			if err := verifyAssociation(resultBytes, result, association, request.Approver, attempt.Briefing, items); err != nil {
				return err
			}
		}
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

func writeEntityDocument(path string, expected *yaml.Node, expectedStatus *string, doc *Document, status *string) error {
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
	replacements := []topLevelReplacement{{key: "gates", data: gatesBlock}}
	if status != nil {
		statusBlock, err := yaml.Marshal(struct {
			Status string `yaml:"status"`
		}{*status})
		if err != nil {
			return err
		}
		replacements = append(replacements, topLevelReplacement{key: "status", data: statusBlock})
	}
	return mutateEntity(path, entityExpectation{Gates: expected, Status: expectedStatus}, func(original []byte) ([]byte, error) {
		out, err := replaceTopLevels(original, status == nil, replacements...)
		if err != nil {
			return nil, err
		}
		if _, _, err := readData(out); err != nil {
			return nil, fmt.Errorf("validate rebuilt gates: %w", err)
		}
		if status != nil {
			parsed, _, _, err := frontmatterNode(out)
			if err != nil || mappingValue(parsed, "status") == nil || mappingValue(parsed, "status").Value != *status {
				return nil, fmt.Errorf("validate rebuilt workflow status")
			}
		}
		return out, nil
	}, atomicWrite)
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
}

func replaceTopLevels(original []byte, insertMissing bool, replacements ...topLevelReplacement) ([]byte, error) {
	root, fmStart, fmEnd, err := frontmatterNode(original)
	if err != nil {
		return nil, err
	}
	for i := range replacements {
		start, end, ok := topLevelRange(root, fmStart, fmEnd, replacements[i].key)
		if !ok {
			if !insertMissing {
				return nil, fmt.Errorf("entity has no %s field", replacements[i].key)
			}
			start, end = fmEnd, fmEnd
		}
		replacements[i].start, replacements[i].end = lineOffset(original, start), lineOffset(original, end)
		if bytes.Contains(original, []byte("\r\n")) {
			replacements[i].data = []byte(strings.ReplaceAll(string(replacements[i].data), "\n", "\r\n"))
		}
	}
	sort.Slice(replacements, func(i, j int) bool { return replacements[i].start > replacements[j].start })
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

func spliceFeedbackCycle(data []byte, line string, cycle int, project bool) ([]byte, error) {
	_, _, fmEnd, err := frontmatterNode(data)
	if err != nil {
		return nil, err
	}
	prefix := fmt.Sprintf("- Cycle %d:", cycle)
	insert, headings, exact, cycles := len(data), 0, 0, 0
	inFence, inSection := false, false
	for offset := lineOffset(data, fmEnd+1); offset < len(data); {
		end := bytes.IndexByte(data[offset:], '\n')
		if end < 0 {
			end = len(data) - offset
		}
		text := strings.TrimSuffix(string(data[offset:offset+end]), "\r")
		trim := strings.TrimSpace(text)
		if strings.HasPrefix(trim, "```") || strings.HasPrefix(trim, "~~~") {
			inFence = !inFence
		} else if !inFence {
			level := strings.IndexFunc(text, func(r rune) bool { return r != '#' })
			if level == 3 && strings.TrimSpace(text[level:]) == "Feedback Cycles" {
				headings++
				inSection = true
			} else if inSection && level > 0 && level <= 3 {
				insert = offset
				inSection = false
			}
			if inSection {
				if text == line {
					exact++
				}
				if strings.HasPrefix(text, prefix) {
					cycles++
				}
			}
		}
		offset += end + 1
	}
	if headings > 1 || !project && cycles != 0 || project && (cycles != exact || cycles > 1) {
		return nil, fmt.Errorf("Feedback Cycles projection conflicts with %s", prefix)
	}
	if !project || exact == 1 {
		return data, nil
	}
	if headings == 0 {
		sep := "\n\n"
		if bytes.HasSuffix(data, []byte("\n\n")) {
			sep = ""
		} else if bytes.HasSuffix(data, []byte("\n")) {
			sep = "\n"
		}
		return append(data, []byte(sep+"### Feedback Cycles\n\n"+line+"\n")...), nil
	}
	sep := ""
	if insert > 0 && data[insert-1] != '\n' {
		sep = "\n"
	}
	add := []byte(sep + line + "\n\n")
	return append(append(append([]byte{}, data[:insert]...), add...), data[insert:]...), nil
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
