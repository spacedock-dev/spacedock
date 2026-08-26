// ABOUTME: Recorder operations, pointer-CAS checks, and digest binding.
// ABOUTME: These operations never model application state or invoke workflow effects.
package gates

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	jsoncanonicalizer "github.com/cyberphone/json-canonicalization/go/src/webpki.org/jsoncanonicalizer"
	"gopkg.in/yaml.v3"
)

type RecordInput struct {
	BriefingPath          string
	LogPath, Round        string
	Actor, Decision       string
	Reason, WorkflowDir   string
	ConnQuote, ConnSource string
}

type WithdrawInput struct {
	Reason      string
	WorkflowDir string
}

type artifactRef struct {
	ID  string `json:"id"`
	URI string `json:"uri,omitempty"`
	Rev string `json:"rev"`
}

type briefingArtifact struct {
	artifactRef
	MediaType string  `json:"mediaType,omitempty"`
	Summary   *string `json:"summary,omitempty"`
}

type briefingManifest struct {
	Type      string             `json:"type"`
	Version   string             `json:"version"`
	ID        string             `json:"id"`
	Question  string             `json:"question"`
	Artifacts []briefingArtifact `json:"artifacts"`
	Context   []json.RawMessage  `json:"context"`
}

type gateRoomRequest struct {
	Type     string `json:"type"`
	Version  string `json:"version"`
	Gate     string `json:"gate"`
	Attempt  string `json:"attempt"`
	Briefing struct {
		Locator string `json:"locator"`
		ID      string `json:"id"`
		Digest  string `json:"digest"`
	} `json:"briefing"`
	Actor    string `json:"actor"`
	Approver string `json:"approver"`
}

type presentedItem struct {
	Type string `json:"type"`
	artifactRef
}

func decodeGateRoomRequest(data []byte) (*gateRoomRequest, error) {
	var request gateRoomRequest
	if err := decodeAuthorityJSON(data, "parse gate room request", &request); err != nil {
		return nil, err
	}
	return &request, nil
}

// RecordSemantic accepts exactly one decision source while the entity lock is
// held. Lifecycle operation, current-stage target, CAS, and ids are derived by
// the recorder rather than supplied in a transaction envelope.
func RecordSemantic(entityPath string, input RecordInput) error {
	_, err := RecordSemanticSummary(entityPath, input)
	return err
}

// RecordSemanticSummary performs one semantic record and returns the exact
// post-write summary while the entity lock is still held.
func RecordSemanticSummary(entityPath string, input RecordInput) (Summary, error) {
	if input.Round != "" {
		if input.BriefingPath == "" || input.LogPath == "" {
			return Summary{}, fmt.Errorf("gate record --round requires --briefing and --log")
		}
		if input.Actor != "" || input.Decision != "" || input.Reason != "" || input.ConnQuote != "" || input.ConnSource != "" {
			return Summary{}, fmt.Errorf("gate record --round is incompatible with gate-closing flags")
		}
		if filepath.Base(input.BriefingPath) != "briefing.json" || filepath.Base(input.LogPath) != "briefing.review.jsonl" {
			return Summary{}, fmt.Errorf("--round inputs must name briefing.json and briefing.review.jsonl")
		}
		if filepath.Base(entityPath) != "index.md" {
			return Summary{}, fmt.Errorf("gate record --round requires folder-form entity <slug>/index.md because review artifacts accumulate beside the entity; convert with `git mv <slug>.md <slug>/index.md` AND rewrite every `room-ref: ./<slug>/` to `room-ref: ./` in the same commit")
		}
		unlock, err := lockEntity(entityPath)
		if err != nil {
			return Summary{}, err
		}
		defer unlock()
		return Summary{}, recordRoundLockedWith(entityPath, input, nil, atomicWrite)
	}
	if input.LogPath != "" {
		return Summary{}, fmt.Errorf("--log requires --round")
	}
	if input.BriefingPath != "" {
		return Summary{}, fmt.Errorf("gate record --briefing requires --round")
	}
	if input.Decision == "" {
		return Summary{}, fmt.Errorf("gate record requires --decision")
	}
	unlock, err := lockEntity(entityPath)
	if err != nil {
		return Summary{}, err
	}
	defer unlock()
	if err := recordChatLocked(entityPath, input); err != nil {
		return Summary{}, err
	}
	doc, _, err := Read(entityPath)
	if err != nil {
		return Summary{}, err
	}
	stage, err := entityStatus(entityPath)
	if err != nil {
		return Summary{}, err
	}
	return CurrentSummary(doc, stage), nil
}

// Withdraw retires the selected prepared attempt without inventing a Resolution
// or application. The prepared-room requirement lives here, not in Validate.
// Validate reads frontmatter only, so it cannot tell a prepared room from an
// archived opaque ref.
func Withdraw(entityPath string, input WithdrawInput) (Summary, error) {
	if strings.TrimSpace(input.Reason) == "" {
		return Summary{}, fmt.Errorf("gate withdraw requires a nonblank --reason")
	}
	if !utf8.ValidString(input.Reason) {
		return Summary{}, fmt.Errorf("--reason must be valid UTF-8")
	}
	unlock, err := lockEntity(entityPath)
	if err != nil {
		return Summary{}, err
	}
	defer unlock()

	doc, _, err := Read(entityPath)
	if err != nil {
		return Summary{}, err
	}
	workflowDir := input.WorkflowDir
	if workflowDir == "" {
		workflowDir = nearestWorkflowDir(filepath.Dir(entityPath))
	}
	if err := validateRetainedAuthority(entityPath, workflowDir, doc); err != nil {
		return Summary{}, err
	}
	doc, oldNode, _, attempt, err := currentStageAttempt(entityPath, workflowDir)
	if err != nil {
		return Summary{}, err
	}
	if state := attemptState(attempt); state != "open" {
		return Summary{}, fmt.Errorf("attempt %s is frozen %s", attempt.ID, state)
	}
	if !preparedRoomBinding(entityPath, attempt.Briefing) {
		return Summary{}, fmt.Errorf("current attempt has no prepared gate room")
	}
	room, err := filepath.Abs(filepath.Join(filepath.Dir(entityPath), filepath.FromSlash(attempt.Briefing.RoomRef)))
	if err != nil {
		return Summary{}, fmt.Errorf("resolve bound gate room: %w", err)
	}
	if err := validatePreparedRoomEntries(room, attempt.Briefing); err != nil {
		return Summary{}, err
	}
	attempt.Withdrawal = &Withdrawal{
		By:     "agent:first-officer",
		At:     time.Now().UTC().Format(time.RFC3339Nano),
		Reason: input.Reason,
	}
	if err := Validate(doc); err != nil {
		return Summary{}, err
	}
	if err := ValidateTransition(oldNode, doc); err != nil {
		return Summary{}, err
	}
	if err := writeDocument(entityPath, oldNode, doc); err != nil {
		return Summary{}, err
	}
	stage, err := entityStatus(entityPath)
	if err != nil {
		return Summary{}, err
	}
	return CurrentSummary(doc, stage), nil
}

func recordChatLocked(entityPath string, input RecordInput) error {
	if input.Actor == "" {
		return fmt.Errorf("--decision requires --actor ID")
	}
	if input.Decision != "approve" && input.Decision != "revise" && input.Decision != "hold" {
		return fmt.Errorf("--decision must be approve, revise, or hold")
	}
	if (input.Decision == "revise" || input.Decision == "hold") && strings.TrimSpace(input.Reason) == "" {
		return fmt.Errorf("%s decision requires --reason", input.Decision)
	}
	if input.Actor != "person:captain" && input.Actor != "agent:first-officer" {
		return fmt.Errorf("unsupported chat decision actor %q", input.Actor)
	}
	if input.Actor == "agent:first-officer" && strings.TrimSpace(input.Reason) == "" {
		return fmt.Errorf("delegated First Officer decision requires --reason")
	}
	// The two record shapes are disjoint by grammar: an FO-rendered chat
	// decision (any of approve/revise/hold) always cites the grant it acted
	// under; a captain decision cites no grant at all, since the captain needs
	// no delegated authority to answer its own gate. The citation attributes;
	// it never authorizes — auto-continue's boundary negative pins that a
	// citation on a no-conn journey still reds.
	hasConnQuote, hasConnSource := strings.TrimSpace(input.ConnQuote) != "", strings.TrimSpace(input.ConnSource) != ""
	if input.Actor == "agent:first-officer" && (!hasConnQuote || !hasConnSource) {
		return fmt.Errorf("delegated First Officer decision requires --conn-quote and --conn-source")
	}
	if input.Actor == "person:captain" && (hasConnQuote || hasConnSource) {
		return fmt.Errorf("--conn-quote and --conn-source are refused on a person:captain decision")
	}
	doc, oldNode, record, attempt, err := currentStageAttempt(entityPath, input.WorkflowDir)
	if err != nil {
		return err
	}
	workflowDir := input.WorkflowDir
	if workflowDir == "" {
		workflowDir = nearestWorkflowDir(filepath.Dir(entityPath))
	}
	if err := validateRetainedAuthority(entityPath, workflowDir, doc); err != nil {
		return err
	}
	if state := attemptState(attempt); state != "open" {
		return fmt.Errorf("attempt %s is frozen %s", attempt.ID, state)
	}
	if !digestRE.MatchString(attempt.Briefing.Digest) {
		return fmt.Errorf("open attempt %s has no verifiable digest", attempt.ID)
	}
	resolution := &Resolution{
		Type:     "Resolution",
		ID:       chatResolutionID(record.ID, attempt.ID),
		Briefing: attempt.Briefing.ID,
		By:       input.Actor,
		At:       time.Now().UTC().Format(time.RFC3339Nano),
		Decision: input.Decision,
		Reason:   input.Reason,
	}
	if input.Actor == "agent:first-officer" {
		resolution.Conn = &Conn{Quote: input.ConnQuote, Source: input.ConnSource}
	}
	if err := closeAttempt(entityPath, input.WorkflowDir, doc, oldNode, record, attempt, resolution); err != nil {
		return err
	}
	return writeDocument(entityPath, oldNode, doc)
}

// lockEntity makes the pointer comparison and gates-only rename one
// process-atomic critical section. Contention fails closed; there is no retry
// loop or lease lifecycle for callers to coordinate.
func lockEntity(path string) (func(), error) {
	lockPath := path + ".gates.lock"
	f, err := os.OpenFile(lockPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		if os.IsExist(err) {
			return nil, fmt.Errorf("concurrent gate writer holds %s", lockPath)
		}
		return nil, err
	}
	return func() {
		_ = f.Close()
		_ = os.Remove(lockPath)
	}, nil
}

func findRecord(doc *Document, id string) *GateRecord {
	for i := range doc.Records {
		if doc.Records[i].ID == id {
			return &doc.Records[i]
		}
	}
	return nil
}

func sameBinding(left, right Briefing) bool {
	return left.ID == right.ID && left.Digest == right.Digest &&
		left.RequestDigest == right.RequestDigest && left.RoomRef == right.RoomRef
}

func currentStageAttempt(entityPath, workflowDir string) (*Document, *yaml.Node, *GateRecord, *Attempt, error) {
	doc, oldNode, err := Read(entityPath)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	stage, err := entityStatus(entityPath)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if workflowDir == "" {
		workflowDir = nearestWorkflowDir(filepath.Dir(entityPath))
	}
	stages, err := applicationStages(filepath.Join(workflowDir, "README.md"))
	if err != nil {
		return nil, nil, nil, nil, err
	}
	stageIndex := applicationStageIndex(stages, stage)
	if stageIndex < 0 {
		return nil, nil, nil, nil, fmt.Errorf("workflow stage %s is not defined in %s", stage, workflowDir)
	}
	if !stages[stageIndex].Gate || stages[stageIndex].Terminal {
		return nil, nil, nil, nil, fmt.Errorf("current workflow stage %s is not an actionable gate:true stage", stage)
	}
	record, err := recordForStage(doc, stage)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	attempt := &record.Attempts[len(record.Attempts)-1]
	if err := validateRecordStage(attempt.Briefing.ID, stage); err != nil {
		return nil, nil, nil, nil, err
	}
	return doc, oldNode, record, attempt, nil
}

func validateRecordStage(briefingID, currentStage string) error {
	briefingStage, ok := canonicalBriefingStage(briefingID)
	if !ok {
		return fmt.Errorf("Briefing id %s is not a canonical stage-qualified v1 identity", briefingID)
	}
	if briefingStage != currentStage {
		return fmt.Errorf("Briefing stage %s does not match current workflow stage %s", briefingStage, currentStage)
	}
	return nil
}

var canonicalBriefingID = regexp.MustCompile(`^briefing:(.+):([^:]+):attempt-[1-9][0-9]*:revision-[1-9][0-9]*$`)

func canonicalBriefingStage(id string) (string, bool) {
	matches := canonicalBriefingID.FindStringSubmatch(id)
	if matches == nil {
		return "", false
	}
	return matches[2], true
}

func closeAttempt(entityPath, workflowDir string, doc *Document, oldNode *yaml.Node, record *GateRecord, attempt *Attempt, resolution *Resolution) error {
	if err := validateResolution(resolution, attempt.Briefing.ID); err != nil {
		return err
	}
	application, err := applicationForDecision(entityPath, workflowDir, record.Stage, resolution.Decision)
	if err != nil {
		return err
	}
	attempt.Resolution = resolution
	attempt.Application = application
	if err := Validate(doc); err != nil {
		return err
	}
	return ValidateTransition(oldNode, doc)
}

func applicationForDecision(entityPath, workflowDir, stage, decision string) (*Application, error) {
	switch decision {
	case "hold":
		// Holds are complete Resolutions and carry no application.
		return nil, nil
	case "revise":
		// Validate the stage taxonomy for advisory-round callers, while keeping
		// the feedback-to route outside the durable application object.
		if workflowDir == "" {
			workflowDir = filepath.Dir(entityPath)
		}
		stages, err := applicationStages(filepath.Join(workflowDir, "README.md"))
		if err != nil {
			return nil, err
		}
		if applicationStageIndex(stages, stage) < 0 {
			return nil, fmt.Errorf("workflow stage %s is not defined in %s", stage, workflowDir)
		}
		return nil, nil
	case "approve":
	default:
		return nil, fmt.Errorf("unsupported application decision %q", decision)
	}
	if workflowDir == "" {
		workflowDir = filepath.Dir(entityPath)
	}
	stages, err := applicationStages(filepath.Join(workflowDir, "README.md"))
	if err != nil {
		return nil, err
	}
	i := applicationStageIndex(stages, stage)
	if i < 0 {
		return nil, fmt.Errorf("workflow stage %s is not defined in %s", stage, workflowDir)
	}
	if i+1 >= len(stages) || strings.TrimSpace(stages[i+1].Name) == "" {
		return nil, fmt.Errorf("workflow stage %s has no advance target", stage)
	}
	return &Application{TargetStage: stages[i+1].Name, State: "pending"}, nil
}

type applicationStage struct {
	Name       string
	FeedbackTo string
	Gate       bool
	Terminal   bool
}

func applicationStageIndex(stages []applicationStage, name string) int {
	return slices.IndexFunc(stages, func(stage applicationStage) bool { return stage.Name == name })
}

func applicationStages(readme string) ([]applicationStage, error) {
	data, err := os.ReadFile(readme)
	if err != nil {
		return nil, fmt.Errorf("read workflow stages: %w", err)
	}
	root, _, _, err := frontmatterNode(data)
	if err != nil {
		return nil, err
	}
	stages := mappingValue(root, "stages")
	if stages == nil {
		return nil, fmt.Errorf("workflow has no stages mapping")
	}
	states := mappingValue(stages, "states")
	if states == nil || states.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("workflow has no stages.states list")
	}
	result := make([]applicationStage, 0, len(states.Content))
	for _, item := range states.Content {
		if item.Kind != yaml.MappingNode {
			continue
		}
		name, feedback := mappingValue(item, "name"), mappingValue(item, "feedback-to")
		if name == nil || strings.TrimSpace(name.Value) == "" {
			continue
		}
		stage := applicationStage{Name: name.Value}
		if feedback != nil {
			stage.FeedbackTo = feedback.Value
		}
		if gate := mappingValue(item, "gate"); gate != nil {
			if err := gate.Decode(&stage.Gate); err != nil {
				return nil, fmt.Errorf("workflow stage %s has invalid gate flag: %w", stage.Name, err)
			}
		}
		if terminal := mappingValue(item, "terminal"); terminal != nil {
			if err := terminal.Decode(&stage.Terminal); err != nil {
				return nil, fmt.Errorf("workflow stage %s has invalid terminal flag: %w", stage.Name, err)
			}
		}
		result = append(result, stage)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("workflow has no application stages")
	}
	return result, nil
}

func chatResolutionID(gateID, attemptID string) string {
	sequence := attemptID
	if cut := strings.LastIndex(attemptID, "-"); cut >= 0 && cut+1 < len(attemptID) {
		sequence = attemptID[cut+1:]
	}
	return "resolution:spacedock:" + strings.TrimPrefix(gateID, "gate:") + ":" + sequence
}

func recordForStage(doc *Document, stage string) (*GateRecord, error) {
	var found *GateRecord
	for i := range doc.Records {
		if doc.Records[i].Stage != stage {
			continue
		}
		if found != nil {
			return nil, fmt.Errorf("multiple logical gates claim workflow stage %s", stage)
		}
		found = &doc.Records[i]
	}
	if found == nil {
		return nil, fmt.Errorf("no logical gate exists for workflow stage %s", stage)
	}
	return found, nil
}

func successorAttemptID(previous string) (string, error) {
	cut := strings.LastIndex(previous, "-")
	if cut < 0 || cut == len(previous)-1 {
		return "", fmt.Errorf("cannot derive successor id from %s", previous)
	}
	sequence, err := strconv.Atoi(previous[cut+1:])
	if err != nil || sequence < 1 {
		return "", fmt.Errorf("cannot derive successor id from %s", previous)
	}
	return previous[:cut+1] + strconv.Itoa(sequence+1), nil
}

func boundBriefingManifest(entityPath string, binding Briefing) (*briefingManifest, error) {
	data, _, err := boundBriefingBytes(entityPath, binding)
	if err != nil {
		return nil, err
	}
	manifest, err := parseBriefingManifest(data)
	if err != nil {
		return nil, err
	}
	if manifest.ID != binding.ID {
		return nil, fmt.Errorf("bound canonical Briefing identity does not match the current binding")
	}
	if preparedRoomBinding(entityPath, binding) {
		if err := validatePreparedSummary(manifest); err != nil {
			return nil, err
		}
	}
	return manifest, nil
}

// boundBriefingPath resolves the binding shapes to the canonical Briefing file,
// in one place, so every reader agrees. A retained request-backed binding
// resolves through the frozen request locator, after the request proves its own
// digest and Briefing identity. A prepared room reads its own locator, the
// reserved name first and the earlier name second. A legacy binding names the
// Briefing file itself.
func boundBriefingPath(entityPath string, binding Briefing) (string, error) {
	retained := filepath.Join(filepath.Dir(entityPath), filepath.FromSlash(binding.RoomRef))
	if binding.RequestDigest != "" {
		requestBytes, err := os.ReadFile(filepath.Join(retained, "request.json"))
		if err != nil {
			return "", fmt.Errorf("read bound gate room request: %w", err)
		}
		requestDigest, err := CanonicalDigest(requestBytes)
		if err != nil {
			return "", fmt.Errorf("canonicalize bound gate room request: %w", err)
		}
		if requestDigest != binding.RequestDigest {
			return "", fmt.Errorf("bound gate room request does not match the frozen request digest")
		}
		request, err := decodeGateRoomRequest(requestBytes)
		if err != nil {
			return "", err
		}
		briefingPath, err := resolveBriefingLocator(retained, request.Briefing.Locator)
		if err != nil {
			return "", err
		}
		if request.Briefing.ID != binding.ID || request.Briefing.Digest != binding.Digest {
			return "", fmt.Errorf("bound gate room request does not match the frozen Briefing identity and digest")
		}
		return briefingPath, nil
	}
	if preparedRoomBinding(entityPath, binding) {
		for _, locator := range preparedLocators {
			if path, err := resolveBriefingLocator(retained, locator); err == nil {
				return path, nil
			}
		}
		// Report against the reserved name, so the error names what is absent.
		return resolveBriefingLocator(retained, preparedBriefingLocator)
	}
	return retained, nil
}

func boundBriefingBytes(entityPath string, binding Briefing) ([]byte, string, error) {
	briefingPath, err := boundBriefingPath(entityPath, binding)
	if err != nil {
		return nil, "", err
	}
	data, err := os.ReadFile(briefingPath)
	if err != nil {
		return nil, "", fmt.Errorf("read bound canonical Briefing: %w", err)
	}
	digest, err := CanonicalDigest(data)
	if err != nil {
		return nil, "", fmt.Errorf("canonicalize bound Briefing: %w", err)
	}
	if digest != binding.Digest {
		return nil, "", fmt.Errorf("bound canonical Briefing bytes do not match the frozen digest")
	}
	return data, briefingPath, nil
}

func resolveBriefingLocator(room, locator string) (string, error) {
	if locator == "" || strings.Contains(locator, "\\") || strings.HasPrefix(locator, "/") ||
		path.Clean(locator) != locator || locator == "." || locator == ".." || strings.HasPrefix(locator, "../") {
		return "", fmt.Errorf("gate room request has an invalid canonical Briefing locator")
	}
	roomAbs, err := filepath.Abs(room)
	if err != nil {
		return "", err
	}
	candidate := filepath.Join(roomAbs, filepath.FromSlash(locator))
	realRoom, err := filepath.EvalSymlinks(roomAbs)
	if err != nil {
		return "", fmt.Errorf("resolve gate room: %w", err)
	}
	realCandidate, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		return "", fmt.Errorf("resolve canonical Briefing locator: %w", err)
	}
	rel, err := filepath.Rel(realRoom, realCandidate)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("canonical Briefing locator escapes the gate room")
	}
	info, err := os.Stat(realCandidate)
	if err != nil || !info.Mode().IsRegular() {
		return "", fmt.Errorf("canonical Briefing locator does not name a regular file")
	}
	return realCandidate, nil
}

func parseBriefingManifest(data []byte) (*briefingManifest, error) {
	var manifest briefingManifest
	if err := decodeAuthorityJSON(data, "parse Briefing", &manifest); err != nil {
		return nil, err
	}
	if manifest.Type != "Briefing" || manifest.Version != "1" || manifest.ID == "" || strings.TrimSpace(manifest.Question) == "" || len(manifest.Artifacts) == 0 {
		return nil, fmt.Errorf("--briefing must name a complete Briefing v1")
	}
	seen := map[string]bool{}
	for _, artifact := range manifest.Artifacts {
		if artifact.ID == "" || artifact.URI == "" || !digestRE.MatchString(artifact.Rev) || seen[artifact.ID] {
			return nil, fmt.Errorf("--briefing has an incomplete or duplicate artifact binding")
		}
		seen[artifact.ID] = true
	}
	return &manifest, nil
}

func canonicalPresentationItems(manifest *briefingManifest) ([]presentedItem, error) {
	references, err := referenceInventory(manifest.Context)
	if err != nil {
		return nil, err
	}
	inventory := make([]presentedItem, 0, len(manifest.Artifacts)+len(references))
	seen := make(map[string]bool, len(inventory)+len(references))
	for _, artifact := range manifest.Artifacts {
		inventory = append(inventory, presentedItem{Type: "Artifact", artifactRef: artifact.artifactRef})
		seen[artifact.ID] = true
	}
	for _, reference := range references {
		if seen[reference.ID] {
			return nil, fmt.Errorf("canonical Briefing has an incomplete or duplicate Reference binding")
		}
		seen[reference.ID] = true
		inventory = append(inventory, presentedItem{Type: "Reference", artifactRef: reference})
	}
	return inventory, nil
}

func referenceInventory(nodes []json.RawMessage) ([]artifactRef, error) {
	var references []artifactRef
	for _, raw := range nodes {
		var node struct {
			Type     string            `json:"type"`
			ID       string            `json:"id"`
			URI      string            `json:"uri"`
			Rev      string            `json:"rev"`
			Summary  json.RawMessage   `json:"summary"`
			Children []json.RawMessage `json:"children"`
		}
		if err := json.Unmarshal(raw, &node); err != nil {
			return nil, fmt.Errorf("parse Briefing context: %w", err)
		}
		if node.Type == "Reference" {
			if node.ID == "" || node.URI == "" || !digestRE.MatchString(node.Rev) {
				return nil, fmt.Errorf("canonical Briefing has an incomplete or duplicate Reference binding")
			}
			if node.Summary != nil {
				return nil, fmt.Errorf("canonical Briefing References must not carry summaries")
			}
			references = append(references, artifactRef{ID: node.ID, URI: node.URI, Rev: node.Rev})
		}
		nested, err := referenceInventory(node.Children)
		if err != nil {
			return nil, err
		}
		references = append(references, nested...)
	}
	return references, nil
}

func ValidateTransition(oldNode *yaml.Node, next *Document) error {
	var old Document
	if err := oldNode.Decode(&old); err != nil {
		return err
	}
	for _, oldRecord := range old.Records {
		nr := findRecord(next, oldRecord.ID)
		if nr == nil {
			return fmt.Errorf("gate %s cannot be deleted", oldRecord.ID)
		}
		for _, oldAttempt := range oldRecord.Attempts {
			state := attemptState(&oldAttempt)
			if state != "closed" && state != "withdrawn" {
				continue
			}
			var found *Attempt
			for i := range nr.Attempts {
				if nr.Attempts[i].ID == oldAttempt.ID {
					found = &nr.Attempts[i]
				}
			}
			unchanged := found != nil && nodesEqual(oldAttempt, *found)
			closedSuperseded := state == "closed" && found != nil && pendingApplicationSuperseded(oldAttempt, *found)
			if !unchanged && !closedSuperseded {
				return fmt.Errorf("frozen %s attempt %s cannot be deleted or mutated", state, oldAttempt.ID)
			}
		}
	}
	return nil
}

func pendingApplicationSuperseded(oldAttempt, nextAttempt Attempt) bool {
	if oldAttempt.Application == nil || oldAttempt.Application.State != "pending" || nextAttempt.Application == nil || nextAttempt.Application.State != "superseded" {
		return false
	}
	oldAttempt.Application.State = "superseded"
	return nodesEqual(oldAttempt, nextAttempt)
}

func nodesEqual(a, b Attempt) bool {
	ab, _ := yaml.Marshal(a)
	bb, _ := yaml.Marshal(b)
	return bytes.Equal(ab, bb)
}

func RawDigest(data []byte) string {
	s := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(s[:])
}

func CanonicalDigest(data []byte) (string, error) {
	if err := rejectDuplicateMembers(data); err != nil {
		return "", err
	}
	canonical, err := jsoncanonicalizer.Transform(data)
	if err != nil {
		return "", err
	}
	return RawDigest(canonical), nil
}
