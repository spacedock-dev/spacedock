// ABOUTME: Mechanical provider-neutral gate-room preparation and atomic binding.
// ABOUTME: Successful preparation publishes only request plus canonical Briefing.
package gates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/spacedock-dev/spacedock/internal/gitsource"
)

const preparedBriefingLocator = "gate-briefing.json"

var prepareWriteBinding = writeDocument

type PrepareInput struct {
	WorkflowDir string
	Question    string
	Artifact    string
	Summary     string
	References  []string
}

type PrepareResult struct {
	Room     string
	Briefing string
	Digest   string
	State    string
}

type preparedReference struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	URI       string `json:"uri"`
	MediaType string `json:"mediaType"`
	Rev       string `json:"rev"`
}

type preparedBriefing struct {
	Type      string              `json:"type"`
	Version   string              `json:"version"`
	ID        string              `json:"id"`
	Question  string              `json:"question"`
	Artifacts []briefingArtifact  `json:"artifacts"`
	Context   []preparedReference `json:"context,omitempty"`
}

// Prepare validates selected committed sources, constructs one canonical open
// room under the entity lock, publishes it, and binds it to the current stage.
func Prepare(entityPath string, input PrepareInput) (PrepareResult, error) {
	if strings.TrimSpace(input.Question) == "" {
		return PrepareResult{}, fmt.Errorf("--question must be nonblank")
	}
	if !utf8.ValidString(input.Question) {
		return PrepareResult{}, fmt.Errorf("--question must be valid UTF-8")
	}
	if !utf8.ValidString(input.Summary) {
		return PrepareResult{}, fmt.Errorf("--summary must be valid UTF-8")
	}
	if strings.TrimSpace(input.Summary) == "" {
		return PrepareResult{}, fmt.Errorf("--summary must be nonblank")
	}
	if !isMarkdownPath(input.Artifact) {
		return PrepareResult{}, fmt.Errorf("--artifact must name a .md or .markdown file")
	}
	if strings.TrimSpace(input.WorkflowDir) == "" {
		return PrepareResult{}, fmt.Errorf("gate prepare requires a workflow directory")
	}
	entityRoot, err := entityResolveRoot(input.WorkflowDir)
	if err != nil {
		return PrepareResult{}, fmt.Errorf("resolve entity root: %w", err)
	}
	paths := append([]string{input.Artifact}, input.References...)
	normalized := make([]string, 0, len(paths))
	seen := map[string]bool{}
	for i, selected := range paths {
		path, err := resolveSelectedSource(selected, entityRoot, i == 0)
		if err != nil {
			return PrepareResult{}, fmt.Errorf("resolve selected source: %w", err)
		}
		if seen[path] {
			// The artifact is always the primary presentation item (sources[0]);
			// a --reference that resolves to the same path as --artifact (or a
			// repeat reference) is redundant, not an error. Drop the duplicate
			// rather than rejecting the prepare — the FO legitimately references the
			// entity under review alongside a gate-review artifact that is the same
			// file, and a hard reject there blocks an otherwise-conforming gate.
			continue
		}
		seen[path] = true
		normalized = append(normalized, path)
	}

	entityPath, err = filepath.Abs(entityPath)
	if err != nil {
		return PrepareResult{}, fmt.Errorf("resolve entity: %w", err)
	}
	if err := refuseNewFlatCompanion(entityPath, input.WorkflowDir); err != nil {
		return PrepareResult{}, err
	}
	unlock, err := lockEntity(entityPath)
	if err != nil {
		return PrepareResult{}, err
	}
	defer unlock()

	stage, err := entityStatus(entityPath)
	if err != nil {
		return PrepareResult{}, err
	}
	if err := validatePreparedStage(input.WorkflowDir, stage); err != nil {
		return PrepareResult{}, err
	}
	doc, oldNode, readErr := Read(entityPath)
	if readErr != nil && !strings.Contains(readErr.Error(), "no gates record") {
		return PrepareResult{}, readErr
	}
	if doc == nil {
		doc = &Document{Version: 1}
	}
	if err := validateRetainedAuthority(entityPath, input.WorkflowDir, doc); err != nil {
		return PrepareResult{}, err
	}
	entityID, err := entityIdentity(entityPath, input.WorkflowDir)
	if err != nil {
		return PrepareResult{}, err
	}
	gateID, attemptID, attemptNumber, record, previous, err := prepareTarget(doc, entityID, stage)
	if err != nil {
		return PrepareResult{}, err
	}
	briefingID := "briefing:" + strings.TrimPrefix(gateID, "gate:") +
		":attempt-" + strconv.Itoa(attemptNumber) + ":revision-1"
	room, err := preparedRoomPath(entityPath, stage, attemptNumber)
	if err != nil {
		return PrepareResult{}, err
	}
	room, err = filepath.Abs(room)
	if err != nil {
		return PrepareResult{}, fmt.Errorf("resolve prepared room: %w", err)
	}

	roots := gitsource.Roots{Main: input.WorkflowDir, State: filepath.Dir(entityPath)}
	sources := make([]gitsource.Source, 0, len(normalized))
	for i, selected := range normalized {
		var source gitsource.Source
		usedRecovered := false
		if selected == entityPath {
			prior, ok, err := preparedEntityReplaySource(entityPath, roots, previous, i)
			if err != nil {
				return PrepareResult{}, err
			}
			if ok {
				source = prior
				usedRecovered = true
			}
		}
		if !usedRecovered {
			var err error
			source, err = gitsource.Inspect(roots, selected)
			if err != nil {
				flag := "--artifact"
				if i > 0 {
					flag = "--reference"
				}
				return PrepareResult{}, fmt.Errorf("%s %s: %w", flag, selected, err)
			}
		}
		sources = append(sources, source)
	}
	if replay, matched, err := preparedReplay(entityPath, previous, briefingID, input.Question, input.Summary, sources); err != nil {
		return PrepareResult{}, err
	} else if matched {
		return replay, nil
	}
	primarySummary := input.Summary
	manifest := preparedBriefing{
		Type:     "Briefing",
		Version:  "1",
		ID:       briefingID,
		Question: input.Question,
		Artifacts: []briefingArtifact{{
			artifactRef: artifactRef{
				ID:  preparedItemID("artifact", briefingID, 1),
				URI: sources[0].URI,
				Rev: sources[0].Rev,
			},
			MediaType: "text/markdown",
			Summary:   &primarySummary,
		}},
	}
	for i, source := range sources[1:] {
		manifest.Context = append(manifest.Context, preparedReference{
			Type:      "Reference",
			ID:        preparedItemID("reference", briefingID, i+2),
			URI:       source.URI,
			MediaType: mediaType(normalized[i+1]),
			Rev:       source.Rev,
		})
	}
	briefingBytes, err := indentedJSON(manifest)
	if err != nil {
		return PrepareResult{}, err
	}
	briefingDigest, err := CanonicalDigest(briefingBytes)
	if err != nil {
		return PrepareResult{}, fmt.Errorf("canonicalize prepared Briefing: %w", err)
	}
	request := gateRoomRequest{
		Type:     "spacedock-gate-presentation-request",
		Version:  "1",
		Gate:     gateID,
		Attempt:  attemptID,
		Actor:    "person:captain",
		Approver: "person:captain",
	}
	request.Briefing.Locator = preparedBriefingLocator
	request.Briefing.ID = briefingID
	request.Briefing.Digest = briefingDigest
	requestBytes, err := indentedJSON(request)
	if err != nil {
		return PrepareResult{}, err
	}
	requestDigest, err := CanonicalDigest(requestBytes)
	if err != nil {
		return PrepareResult{}, fmt.Errorf("canonicalize prepared request: %w", err)
	}
	roomRef, err := relativeRoomRef(entityPath, room)
	if err != nil {
		return PrepareResult{}, err
	}
	binding := Briefing{
		ID:            briefingID,
		Digest:        briefingDigest,
		RequestDigest: requestDigest,
		RoomRef:       roomRef,
	}

	if previous != nil && attemptState(previous) == "open" && previous.Briefing.RequestDigest != "" &&
		!sameBinding(previous.Briefing, binding) {
		return PrepareResult{}, fmt.Errorf("open gate room binding is frozen and cannot be rebound")
	}
	if err := validatePreparedCandidate(entityPath, roots, room, binding, gateID, attemptID, briefingBytes, requestBytes); err != nil {
		return PrepareResult{}, err
	}
	if err := validatePreparedRoomAncestry(entityPath, room); err != nil {
		return PrepareResult{}, err
	}
	created, createdParents, err := publishPreparedRoom(room, briefingBytes, requestBytes)
	if err != nil {
		return PrepareResult{}, err
	}
	if err := validatePreparedRoomEntries(room); err != nil {
		if created {
			rollbackPreparedRoom(room, createdParents)
		}
		return PrepareResult{}, err
	}

	if previous == nil {
		record.Attempts = append(record.Attempts, Attempt{ID: attemptID, Briefing: binding})
	} else if attemptState(previous) == "open" {
		previous.Briefing = binding
	} else {
		if attemptState(previous) == "closed" && previous.Application != nil && previous.Application.State == "pending" {
			previous.Application.State = "superseded"
		}
		record.Attempts = append(record.Attempts, Attempt{ID: attemptID, Briefing: binding})
	}
	if err := Validate(doc); err != nil {
		if created {
			rollbackPreparedRoom(room, createdParents)
		}
		return PrepareResult{}, err
	}
	if oldNode != nil {
		if err := ValidateTransition(oldNode, doc); err != nil {
			if created {
				rollbackPreparedRoom(room, createdParents)
			}
			return PrepareResult{}, err
		}
	}
	if err := prepareWriteBinding(entityPath, oldNode, doc); err != nil {
		if created {
			rollbackPreparedRoom(room, createdParents)
		}
		return PrepareResult{}, err
	}
	return PrepareResult{Room: room, Briefing: briefingID, Digest: briefingDigest, State: "open"}, nil
}

func preparedEntityReplaySource(entityPath string, roots gitsource.Roots, previous *Attempt, ordinal int) (gitsource.Source, bool, error) {
	if previous == nil || attemptState(previous) != "open" || previous.Briefing.RequestDigest == "" {
		return gitsource.Source{}, false, nil
	}
	manifest, err := boundBriefingManifest(entityPath, previous.Briefing)
	if err != nil {
		return gitsource.Source{}, false, err
	}
	items, err := canonicalPresentationItems(manifest)
	if err != nil {
		return gitsource.Source{}, false, err
	}
	if ordinal >= len(items) || !strings.HasPrefix(items[ordinal].URI, "git-root://") {
		return gitsource.Source{}, false, nil
	}
	source := gitsource.Source{URI: items[ordinal].URI, Rev: items[ordinal].Rev}
	frozen, err := gitsource.Resolve(roots, source.URI, source.Rev)
	if err != nil {
		return gitsource.Source{}, false, err
	}
	current, err := os.ReadFile(entityPath)
	if err != nil {
		return gitsource.Source{}, false, err
	}
	frozenOutside, err := entityWithoutGates(frozen)
	if err != nil {
		return gitsource.Source{}, false, err
	}
	currentOutside, err := entityWithoutGates(current)
	if err != nil {
		return gitsource.Source{}, false, err
	}
	if !bytes.Equal(frozenOutside, currentOutside) {
		return gitsource.Source{}, false, nil
	}
	return source, true, nil
}

func entityWithoutGates(data []byte) ([]byte, error) {
	root, fmStart, fmEnd, err := frontmatterNode(data)
	if err != nil {
		return nil, err
	}
	start, end, ok := topLevelRange(root, fmStart, fmEnd, "gates")
	if !ok {
		return append([]byte(nil), data...), nil
	}
	startByte, endByte := lineOffset(data, start), lineOffset(data, end)
	return append(append([]byte(nil), data[:startByte]...), data[endByte:]...), nil
}

func preparedReplay(entityPath string, previous *Attempt, briefingID, question, summary string, sources []gitsource.Source) (PrepareResult, bool, error) {
	if previous == nil || attemptState(previous) != "open" || previous.Briefing.RequestDigest == "" {
		return PrepareResult{}, false, nil
	}
	manifest, err := boundBriefingManifest(entityPath, previous.Briefing)
	if err != nil {
		return PrepareResult{}, false, err
	}
	if manifest.ID != briefingID || manifest.Question != question || len(manifest.Artifacts) != 1 ||
		manifest.Artifacts[0].Summary == nil || *manifest.Artifacts[0].Summary != summary {
		return PrepareResult{}, false, nil
	}
	items, err := canonicalPresentationItems(manifest)
	if err != nil {
		return PrepareResult{}, false, err
	}
	if len(items) != len(sources) {
		return PrepareResult{}, false, nil
	}
	for i, item := range items {
		wantType := "Reference"
		if i == 0 {
			wantType = "Artifact"
		}
		if item.Type != wantType {
			return PrepareResult{}, false, nil
		}
		same, err := gitsource.SameLogicalRevision(
			gitsource.Source{URI: item.URI, Rev: item.Rev},
			sources[i],
		)
		if err != nil {
			return PrepareResult{}, false, err
		}
		if !same {
			return PrepareResult{}, false, nil
		}
	}
	room, err := filepath.Abs(filepath.Join(filepath.Dir(entityPath), filepath.FromSlash(previous.Briefing.RoomRef)))
	if err != nil {
		return PrepareResult{}, false, fmt.Errorf("resolve prepared room: %w", err)
	}
	if err := validatePreparedRoomEntries(room); err != nil {
		return PrepareResult{}, false, err
	}
	return PrepareResult{
		Room:     room,
		Briefing: previous.Briefing.ID,
		Digest:   previous.Briefing.Digest,
		State:    "open",
	}, true, nil
}

func prepareTarget(doc *Document, entityID, stage string) (gateID, attemptID string, attemptNumber int, record *GateRecord, previous *Attempt, err error) {
	record, lookupErr := recordForStage(doc, stage)
	if lookupErr != nil && !strings.Contains(lookupErr.Error(), "no logical gate") {
		return "", "", 0, nil, nil, lookupErr
	}
	if record == nil {
		entity := entityID
		if cut := strings.LastIndex(entity, ":"); cut >= 0 {
			entity = entity[cut+1:]
		}
		if entity == "" {
			return "", "", 0, nil, nil, fmt.Errorf("cannot derive gate identity from entity")
		}
		gateID = "gate:" + entityID + ":" + stage
		attemptNumber = 1
		attemptID = "gate-attempt:" + entity + "-" + stage + "-1"
		doc.Records = append(doc.Records, GateRecord{ID: gateID, Stage: stage})
		record = &doc.Records[len(doc.Records)-1]
		return gateID, attemptID, attemptNumber, record, nil, nil
	}
	gateID = record.ID
	previous = &record.Attempts[len(record.Attempts)-1]
	attemptNumber, err = attemptSequence(previous.ID)
	if err != nil {
		return "", "", 0, nil, nil, err
	}
	if attemptState(previous) == "open" {
		return gateID, previous.ID, attemptNumber, record, previous, nil
	}
	attemptNumber++
	attemptID, err = successorAttemptID(previous.ID)
	return gateID, attemptID, attemptNumber, record, previous, err
}

func entityIdentity(path, workflowDir string) (string, error) {
	readme, err := os.ReadFile(filepath.Join(workflowDir, "README.md"))
	if err != nil {
		return "", err
	}
	workflow, _, _, err := frontmatterNode(readme)
	if err != nil {
		return "", err
	}
	if style := mappingValue(workflow, "id-style"); style != nil && strings.TrimSpace(style.Value) == "slug" {
		return entitySlug(path), nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	root, _, _, err := frontmatterNode(data)
	if err != nil {
		return "", err
	}
	id := mappingValue(root, "id")
	if id == nil || strings.TrimSpace(id.Value) == "" {
		return "", fmt.Errorf("entity has no identity")
	}
	return id.Value, nil
}

func validatePreparedCandidate(entityPath string, roots gitsource.Roots, room string, binding Briefing, gateID, attemptID string, briefingBytes, requestBytes []byte) error {
	tmp, err := os.MkdirTemp("", "spacedock-gate-prepare-validate-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	briefingPath := filepath.Join(tmp, preparedBriefingLocator)
	if err := os.WriteFile(briefingPath, briefingBytes, 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(tmp, "request.json"), requestBytes, 0o644); err != nil {
		return err
	}
	manifest, err := parseBriefingManifest(briefingBytes)
	if err != nil {
		return err
	}
	if err := validatePreparedSummary(manifest); err != nil {
		return err
	}
	items, err := canonicalPresentationItems(manifest)
	if err != nil {
		return err
	}
	for _, item := range items {
		if _, err := gitsource.Resolve(roots, item.URI, item.Rev); err != nil {
			return err
		}
	}
	candidateBinding := binding
	candidateBinding.RoomRef = "./"
	return validateGateRoomRequest(briefingPath, candidateBinding, gateID, attemptID)
}

func validatePreparedRoomAncestry(entityPath, room string) error {
	trustedHome := filepath.Dir(entityPath)
	parent := filepath.Dir(room)
	rel, err := filepath.Rel(trustedHome, parent)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("prepared room parent escapes the trusted entity home")
	}
	current := trustedHome
	for _, component := range strings.Split(rel, string(filepath.Separator)) {
		if component == "." || component == "" {
			continue
		}
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("prepared room parent %s is a symlink", current)
		}
		if !info.IsDir() {
			return fmt.Errorf("prepared room parent %s is not a directory", current)
		}
	}
	return nil
}

func validatePreparedRoomEntries(room string) error {
	entries, err := os.ReadDir(room)
	if err != nil {
		return err
	}
	expected := map[string]bool{preparedBriefingLocator: true, "request.json": true}
	if len(entries) != len(expected) {
		return fmt.Errorf("prepared room must contain exactly two regular files")
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !expected[entry.Name()] || !info.Mode().IsRegular() {
			return fmt.Errorf("prepared room entry %s is not an expected regular file", entry.Name())
		}
	}
	return nil
}

func publishPreparedRoom(room string, briefing, request []byte) (bool, []string, error) {
	info, err := os.Lstat(room)
	if err == nil {
		if !info.IsDir() {
			return false, nil, fmt.Errorf("prepared room target is occupied")
		}
		entries, readErr := os.ReadDir(room)
		if readErr != nil || len(entries) != 2 {
			return false, nil, fmt.Errorf("prepared room target is occupied by divergent content")
		}
		currentBriefing, briefingErr := os.ReadFile(filepath.Join(room, preparedBriefingLocator))
		currentRequest, requestErr := os.ReadFile(filepath.Join(room, "request.json"))
		if briefingErr != nil || requestErr != nil || !bytes.Equal(currentBriefing, briefing) || !bytes.Equal(currentRequest, request) {
			return false, nil, fmt.Errorf("prepared room target is occupied by divergent content")
		}
		return false, nil, nil
	}
	if !os.IsNotExist(err) {
		return false, nil, err
	}
	parent := filepath.Dir(room)
	createdParents, err := createPreparedParents(parent)
	if err != nil {
		return false, nil, err
	}
	tmp, err := os.MkdirTemp(parent, ".prepare-*")
	if err != nil {
		removePreparedParents(createdParents)
		return false, nil, err
	}
	defer os.RemoveAll(tmp)
	if err := writeSyncedFile(filepath.Join(tmp, preparedBriefingLocator), briefing); err != nil {
		removePreparedParents(createdParents)
		return false, nil, err
	}
	if err := writeSyncedFile(filepath.Join(tmp, "request.json"), request); err != nil {
		removePreparedParents(createdParents)
		return false, nil, err
	}
	if err := os.Rename(tmp, room); err != nil {
		removePreparedParents(createdParents)
		return false, nil, err
	}
	return true, createdParents, nil
}

func createPreparedParents(path string) ([]string, error) {
	var missing []string
	for current := path; ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err == nil {
			if !info.IsDir() {
				return nil, fmt.Errorf("prepared room parent is not a directory")
			}
			break
		}
		if !os.IsNotExist(err) {
			return nil, err
		}
		missing = append(missing, current)
		parent := filepath.Dir(current)
		if parent == current {
			return nil, fmt.Errorf("prepared room has no existing parent")
		}
	}
	created := make([]string, 0, len(missing))
	for i := len(missing) - 1; i >= 0; i-- {
		if err := os.Mkdir(missing[i], 0o755); err != nil {
			removePreparedParents(created)
			return nil, err
		}
		created = append(created, missing[i])
	}
	return created, nil
}

func rollbackPreparedRoom(room string, createdParents []string) {
	_ = os.RemoveAll(room)
	removePreparedParents(createdParents)
}

func removePreparedParents(created []string) {
	for i := len(created) - 1; i >= 0; i-- {
		_ = os.Remove(created[i])
	}
}

func validatePreparedSummary(manifest *briefingManifest) error {
	if len(manifest.Artifacts) == 0 || manifest.Artifacts[0].Summary == nil || strings.TrimSpace(*manifest.Artifacts[0].Summary) == "" {
		return fmt.Errorf("request-backed Briefing requires a nonblank primary Artifact summary")
	}
	for _, artifact := range manifest.Artifacts[1:] {
		if artifact.Summary != nil {
			return fmt.Errorf("request-backed Briefing References must not carry summaries")
		}
	}
	return nil
}

func validatePreparedStage(workflowDir, stage string) error {
	if stage == "" || strings.TrimSpace(stage) != stage ||
		stage == "." || stage == ".." || strings.ContainsAny(stage, `/\`) ||
		filepath.Base(stage) != stage {
		return fmt.Errorf("workflow stage %q is not a safe room path segment", stage)
	}
	stages, err := applicationStages(filepath.Join(workflowDir, "README.md"))
	if err != nil {
		return err
	}
	stageIndex := applicationStageIndex(stages, stage)
	if stageIndex < 0 {
		return fmt.Errorf("workflow stage %s is not defined in %s", stage, workflowDir)
	}
	if !stages[stageIndex].Gate || stages[stageIndex].Terminal {
		return fmt.Errorf("workflow stage %s is not an actionable gate", stage)
	}
	return nil
}

// refuseNewFlatCompanion refuses to create the first prepared room beside a flat
// entity, but only where the workflow declares folder form. A room lands at
// <state>/<slug>/review/..., so preparing beside <slug>.md puts a <slug>/
// companion next to it and writes refs relative to the state root as
// ./<slug>/review/... — refs that break the moment the entity becomes
// <slug>/index.md. Only a workflow that declares `entity-form: folder` has
// promised that its entities are folder form, so only there is a flat entity a
// mistake worth blocking a gate over. A workflow that declares nothing accepts
// either shape: that is every workflow today, and flat is still what filing
// mints by default.
//
// An entity that already holds rooms is grandfathered even under the
// declaration. Its refs are correct for the form it is in, so refusing would
// block the next gate on an entity that works — the green-then-blocked sequence
// this guard exists to prevent — rather than protect anything. A declaration
// added to a workflow that already holds hybrids therefore stops new ones
// without stranding the old; `status --validate` reports those with the
// conversion remedy.
func refuseNewFlatCompanion(entityPath, workflowDir string) error {
	if filepath.Base(entityPath) == "index.md" {
		return nil
	}
	if !workflowDeclaresFolderForm(workflowDir) {
		return nil
	}
	slug := entitySlug(entityPath)
	if _, err := os.Stat(filepath.Join(filepath.Dir(entityPath), slug, "review")); err == nil {
		return nil
	}
	return fmt.Errorf("gate prepare requires folder-form entity %s/index.md because this workflow declares `entity-form: folder` and review artifacts accumulate beside the entity; file it as %s/index.md, or move it with `git mv %s.md %s/index.md`", slug, slug, slug, slug)
}

// workflowDeclaresFolderForm reports whether the workflow README declares
// `entity-form: folder`. Preparation only reads the key; commissioning and the
// filing-time default own writing it. An unreadable or unparseable README is
// not a declaration — validatePreparedStage reads the same file below and owns
// the error message for a workflow that cannot be read at all.
func workflowDeclaresFolderForm(workflowDir string) bool {
	readme, err := os.ReadFile(filepath.Join(workflowDir, "README.md"))
	if err != nil {
		return false
	}
	workflow, _, _, err := frontmatterNode(readme)
	if err != nil {
		return false
	}
	form := mappingValue(workflow, "entity-form")
	return form != nil && strings.TrimSpace(form.Value) == "folder"
}

func preparedRoomPath(entityPath, stage string, attempt int) (string, error) {
	slug := entitySlug(entityPath)
	home := filepath.Dir(entityPath)
	if filepath.Base(entityPath) != "index.md" {
		home = filepath.Join(home, slug)
	}
	reviewRoot := filepath.Join(home, "review")
	room := filepath.Join(reviewRoot, stage, "briefing-"+strconv.Itoa(attempt))
	rel, err := filepath.Rel(reviewRoot, room)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("prepared room escapes the entity review directory")
	}
	return room, nil
}

func relativeRoomRef(entityPath, room string) (string, error) {
	ref, err := filepath.Rel(filepath.Dir(entityPath), room)
	if err != nil {
		return "", fmt.Errorf("resolve prepared room reference: %w", err)
	}
	ref = filepath.ToSlash(ref)
	if !strings.HasPrefix(ref, ".") {
		ref = "./" + ref
	}
	return ref, nil
}

func entitySlug(entityPath string) string {
	if filepath.Base(entityPath) == "index.md" {
		return filepath.Base(filepath.Dir(entityPath))
	}
	return strings.TrimSuffix(filepath.Base(entityPath), filepath.Ext(entityPath))
}

func attemptSequence(id string) (int, error) {
	cut := strings.LastIndex(id, "-")
	if cut < 0 || cut == len(id)-1 {
		return 0, fmt.Errorf("cannot derive attempt sequence from %s", id)
	}
	n, err := strconv.Atoi(id[cut+1:])
	if err != nil || n < 1 {
		return 0, fmt.Errorf("cannot derive attempt sequence from %s", id)
	}
	return n, nil
}

func preparedItemID(kind, briefingID string, ordinal int) string {
	return kind + ":" + strings.TrimPrefix(briefingID, "briefing:") + ":item-" + strconv.Itoa(ordinal)
}

func mediaType(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".markdown":
		return "text/markdown"
	case ".json":
		return "application/json"
	case ".yaml", ".yml":
		return "application/yaml"
	case ".txt", ".log":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}

// entityResolveRoot returns the directory relative --artifact and --reference
// paths resolve against. In a split-root workflow it is the state-checkout
// entity root, computed once from the README `state:` field of workflowDir —
// the same root the entity path is derived from. In a single-root workflow
// (state: absent, empty, or $inline) it is workflowDir itself, matching the
// prior cwd-based behavior when cwd and the workflow root coincide. The
// absolute/.. rejection mirrors status.ClassifyState but stays local to gates
// to avoid a status→gates import cycle.
func entityResolveRoot(workflowDir string) (string, error) {
	readme, err := os.ReadFile(filepath.Join(workflowDir, "README.md"))
	if err != nil {
		return "", err
	}
	root, _, _, err := frontmatterNode(readme)
	if err != nil {
		return "", err
	}
	state := mappingValue(root, "state")
	if state == nil {
		return workflowDir, nil
	}
	value := strings.TrimSpace(state.Value)
	if value == "" || value == "$inline" {
		return workflowDir, nil
	}
	if filepath.IsAbs(value) {
		return "", fmt.Errorf("state: must be a path relative to the workflow README directory, not absolute: %s", value)
	}
	cleaned := filepath.Clean(value)
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("state: must not escape the workflow README directory: %s", value)
	}
	return filepath.Join(workflowDir, cleaned), nil
}

// resolveSelectedSource makes a selected source path absolute. Relative paths
// resolve against the entity root (the state-checkout root in split-root, the
// workflow dir in single-root); absolute paths pass through cleaned. This is
// the single resolution site — the CLI passes relative paths through unchanged.
//
// A relative path that already carries the state-checkout basename (e.g.
// ".spacedock-state/auto-continue-task/index.md" in a split-root workflow where
// entityRoot is "<workflow>/.spacedock-state") is workflow-rooted, not
// entity-rooted: joining it under entityRoot would double the basename. Resolve
// it against the workflow directory (entityRoot's parent) instead.
//
// The artifact resolves strictly against the entity root — a wrong-root
// artifact is rejected (TestPrepareWrongRootRelativeArtifactFails). A
// reference may legitimately live at the workflow root (e.g.
// recorder-contract.md alongside a split-root state checkout), so references
// fall back to the workflow directory (entityRoot's parent) when the entity-root
// join does not exist. A genuinely missing file reports the entity-root seek
// path so the error shape is preserved.
func resolveSelectedSource(selected, entityRoot string, isArtifact bool) (string, error) {
	if filepath.IsAbs(selected) {
		return filepath.Clean(selected), nil
	}
	selected = filepath.Clean(selected)
	base := filepath.Base(entityRoot)
	// A path carrying the state-checkout basename is workflow-rooted; resolve
	// against the workflow directory so the basename is not doubled. This
	// applies to both artifact and reference.
	if base != "." && strings.HasPrefix(selected, base+string(filepath.Separator)) {
		return filepath.Clean(filepath.Join(filepath.Dir(entityRoot), selected)), nil
	}
	entityJoin := filepath.Clean(filepath.Join(entityRoot, selected))
	if isArtifact {
		return entityJoin, nil
	}
	// Reference: fall back to the workflow root if the entity-root join is absent.
	if info, err := os.Lstat(entityJoin); err == nil && info.Mode().IsRegular() {
		return entityJoin, nil
	}
	workflowJoin := filepath.Clean(filepath.Join(filepath.Dir(entityRoot), selected))
	if info, err := os.Lstat(workflowJoin); err == nil && info.Mode().IsRegular() {
		return workflowJoin, nil
	}
	return entityJoin, nil // not found; report the entity-root seek path
}

func isMarkdownPath(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".markdown"
}

func indentedJSON(value any) ([]byte, error) {
	body, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(body, '\n'), nil
}
