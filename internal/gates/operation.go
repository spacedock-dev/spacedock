// ABOUTME: Recorder operations, pointer-CAS checks, digest binding, and room-backed results.
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
	"github.com/spacedock-dev/spacedock/internal/gitsource"
	"gopkg.in/yaml.v3"
)

type RecordInput struct {
	BriefingPath, RoomPath            string
	LogPath, FeedbackCyclePath, Round string
	Actor, Decision                   string
	Reason, WorkflowDir               string
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

type providerResult struct {
	Type        string       `json:"type"`
	Briefing    string       `json:"briefing"`
	Artifact    artifactRef  `json:"artifact"`
	Resolution  Resolution   `json:"resolution"`
	Annotations []Annotation `json:"annotations"`
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

type presentedInventory struct {
	Items []presentedItem `json:"items"`
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

func decodeProviderResult(data []byte) (*providerResult, error) {
	var result providerResult
	if err := decodeAuthorityJSON(data, "parse Result", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func decodePresentedInventory(data []byte) (*presentedInventory, error) {
	var inventory presentedInventory
	if err := decodeAuthorityJSON(data, "parse presented inventory", &inventory); err != nil {
		return nil, err
	}
	return &inventory, nil
}

type resultAssociation struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Result  struct {
		Digest   string `json:"digest"`
		Briefing string `json:"briefing"`
	} `json:"result"`
	Actor     string `json:"actor"`
	Canonical struct {
		Briefing  string        `json:"briefing"`
		Revision  string        `json:"revision"`
		Artifacts []artifactRef `json:"artifacts"`
	} `json:"canonical"`
	Presentation []struct {
		Provider  artifactRef `json:"provider"`
		Canonical artifactRef `json:"canonical"`
	} `json:"presentation"`
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
		if input.RoomPath != "" || input.Actor != "" || input.Decision != "" || input.Reason != "" {
			return Summary{}, fmt.Errorf("gate record --round is incompatible with gate-closing flags")
		}
		if filepath.Base(input.BriefingPath) != "briefing.json" || filepath.Base(input.LogPath) != "briefing.review.jsonl" {
			return Summary{}, fmt.Errorf("--round inputs must name briefing.json and briefing.review.jsonl")
		}
		if filepath.Base(entityPath) != "index.md" {
			return Summary{}, fmt.Errorf("gate record --round requires folder-form entity <slug>/index.md because review artifacts accumulate beside the entity")
		}
		unlock, err := lockEntity(entityPath)
		if err != nil {
			return Summary{}, err
		}
		defer unlock()
		return Summary{}, recordRoundLockedWith(entityPath, input, nil, atomicWrite)
	}
	if input.LogPath != "" || input.FeedbackCyclePath != "" {
		return Summary{}, fmt.Errorf("--log and --feedback-cycle require --round")
	}
	if input.BriefingPath != "" {
		return Summary{}, fmt.Errorf("gate record --briefing requires --round")
	}
	sources := 0
	for _, source := range []string{input.RoomPath, input.Decision} {
		if source != "" {
			sources++
		}
	}
	if sources != 1 {
		return Summary{}, fmt.Errorf("gate record requires exactly one of --room or --decision")
	}
	if input.RoomPath != "" && (input.Actor != "" || input.Decision != "" || input.Reason != "") {
		return Summary{}, fmt.Errorf("gate record flags do not match the selected semantic source")
	}
	unlock, err := lockEntity(entityPath)
	if err != nil {
		return Summary{}, err
	}
	defer unlock()
	var recordErr error
	if input.RoomPath != "" {
		recordErr = recordRoomLocked(entityPath, input)
	} else {
		recordErr = recordChatLocked(entityPath, input)
	}
	if recordErr != nil {
		return Summary{}, recordErr
	}
	doc, _, err := Read(entityPath)
	if err != nil {
		return Summary{}, err
	}
	return CurrentSummary(doc), nil
}

// Withdraw retires the selected request-backed prepared attempt without
// inventing a Resolution or application.
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
	doc, oldNode, record, attempt, err := currentStageAttempt(entityPath, workflowDir)
	if err != nil {
		return Summary{}, err
	}
	if doc.Current.Gate != record.ID {
		return Summary{}, fmt.Errorf("current workflow stage does not match gates.current selection")
	}
	if state := attemptState(attempt); state != "open" {
		return Summary{}, fmt.Errorf("attempt %s is frozen %s", attempt.ID, state)
	}
	if attempt.Briefing.RequestDigest == "" {
		return Summary{}, fmt.Errorf("current attempt is not request-backed")
	}
	room, err := filepath.Abs(filepath.Join(filepath.Dir(entityPath), filepath.FromSlash(attempt.Briefing.RoomRef)))
	if err != nil {
		return Summary{}, fmt.Errorf("resolve bound gate room: %w", err)
	}
	if err := validatePreparedRoomEntries(room); err != nil {
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
	return CurrentSummary(doc), nil
}

func recordRoomLocked(entityPath string, input RecordInput) error {
	doc, oldNode, record, attempt, err := currentStageAttempt(entityPath, input.WorkflowDir)
	if err != nil {
		return err
	}
	workflowDir := input.WorkflowDir
	if workflowDir == "" {
		workflowDir = nearestWorkflowDir(filepath.Dir(entityPath))
	}
	if err := validateRetainedAuthorityExcept(entityPath, workflowDir, doc, record.ID, attempt.ID); err != nil {
		return err
	}
	if state := attemptState(attempt); state != "open" {
		return fmt.Errorf("attempt %s is frozen %s", attempt.ID, state)
	}
	if !digestRE.MatchString(attempt.Briefing.Digest) {
		return fmt.Errorf("open attempt %s has no verifiable digest", attempt.ID)
	}
	if attempt.Briefing.RequestDigest == "" {
		return fmt.Errorf("current attempt has no frozen request digest for room-backed recording")
	}
	roomPath, err := filepath.Abs(input.RoomPath)
	if err != nil {
		return fmt.Errorf("resolve gate room: %w", err)
	}
	boundRoomPath, err := filepath.Abs(filepath.Join(filepath.Dir(entityPath), filepath.FromSlash(attempt.Briefing.RoomRef)))
	if err != nil {
		return fmt.Errorf("resolve bound gate room: %w", err)
	}
	if filepath.Clean(roomPath) != filepath.Clean(boundRoomPath) {
		return fmt.Errorf("--room is not the current attempt's bound gate room")
	}
	requestBytes, err := os.ReadFile(filepath.Join(roomPath, "request.json"))
	if err != nil {
		return fmt.Errorf("read gate room request: %w", err)
	}
	requestDigest, err := CanonicalDigest(requestBytes)
	if err != nil {
		return fmt.Errorf("canonicalize gate room request: %w", err)
	}
	if requestDigest != attempt.Briefing.RequestDigest {
		return fmt.Errorf("gate room request does not match the frozen request digest")
	}
	request, err := decodeGateRoomRequest(requestBytes)
	if err != nil {
		return err
	}
	if request.Type != "spacedock-gate-presentation-request" || request.Version != "1" ||
		request.Gate != record.ID || request.Attempt != attempt.ID ||
		request.Briefing.ID != attempt.Briefing.ID || request.Briefing.Digest != attempt.Briefing.Digest ||
		request.Actor != "person:captain" || request.Approver != "person:captain" {
		return fmt.Errorf("gate room request does not bind the current gate, attempt, Briefing, and captain authority")
	}
	resultPath := filepath.Join(roomPath, "provider", "result.json")
	resultBytes, err := os.ReadFile(resultPath)
	if err != nil {
		return fmt.Errorf("read provider Result: %w", err)
	}
	result, err := decodeProviderResult(resultBytes)
	if err != nil {
		return err
	}
	var envelope map[string]json.RawMessage
	if err := decodeAuthorityJSON(resultBytes, "parse Result envelope", &envelope); err != nil {
		return err
	}
	for _, field := range []string{"status", "binding", "actor", "approver", "resolutionId"} {
		if _, present := envelope[field]; present {
			return fmt.Errorf("advisory Result remains evidence; closing a gate requires a separate minimal binding Result")
		}
	}
	for field := range envelope {
		switch field {
		case "type", "briefing", "artifact", "resolution", "annotations":
		default:
			return fmt.Errorf("binding Result has unknown top-level field %q", field)
		}
	}
	resolutionDecoder := json.NewDecoder(bytes.NewReader(envelope["resolution"]))
	resolutionDecoder.DisallowUnknownFields()
	if err := resolutionDecoder.Decode(&result.Resolution); err != nil {
		return fmt.Errorf("parse Result Resolution: %w", err)
	}
	if result.Briefing != request.Briefing.ID {
		return fmt.Errorf("binding Result does not bind the gate room's canonical Briefing")
	}
	if result.Resolution.By != request.Approver {
		return fmt.Errorf("binding Result Resolution.by does not match the gate room request authority")
	}
	manifest, err := boundBriefingManifest(entityPath, attempt.Briefing)
	if err != nil {
		return err
	}
	canonicalItems, err := canonicalPresentationItems(manifest)
	if err != nil {
		return err
	}
	roots := gitsource.Roots{Main: workflowDir, State: filepath.Dir(entityPath)}
	gitItems := 0
	for _, item := range canonicalItems {
		if strings.HasPrefix(item.URI, "git-root://") {
			gitItems++
			if _, err := gitsource.Resolve(roots, item.URI, item.Rev); err != nil {
				return fmt.Errorf("resolve selected source: %w", err)
			}
		}
	}
	if gitItems != 0 && gitItems != len(canonicalItems) {
		return fmt.Errorf("canonical Briefing mixes Git-root and non-Git selected source identities")
	}
	presentedBytes, err := os.ReadFile(filepath.Join(roomPath, "provider", "presented-inventory.json"))
	if err != nil {
		return fmt.Errorf("read presented inventory: %w", err)
	}
	presented, err := decodePresentedInventory(presentedBytes)
	if err != nil {
		return err
	}
	association, err := deriveAssociation(resultBytes, result, presented, request.Approver, attempt.Briefing, canonicalItems)
	if err != nil {
		return err
	}
	if err := verifyAssociation(resultBytes, result, association, request.Approver, attempt.Briefing, canonicalItems); err != nil {
		return err
	}
	resolution := result.Resolution
	resolution.Briefing = attempt.Briefing.ID
	attempt.ProviderEvidence = &ProviderEvidence{
		ResultDigest:             RawDigest(resultBytes),
		PresentedInventoryDigest: RawDigest(presentedBytes),
	}
	if err := closeAttempt(entityPath, input.WorkflowDir, doc, oldNode, record, attempt, &resolution); err != nil {
		return err
	}
	return writeDocument(entityPath, oldNode, doc)
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
	return left.ID == right.ID && left.Digest == right.Digest && left.DigestDomain == right.DigestDomain &&
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
	doc.Current.Gate = record.ID
	if err := Validate(doc); err != nil {
		return err
	}
	return ValidateTransition(oldNode, doc)
}

func applicationForDecision(entityPath, workflowDir, stage, decision string) (*Application, error) {
	switch decision {
	case "hold":
		return &Application{Action: "none", State: "not-applicable"}, nil
	case "approve", "revise":
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
	if decision == "revise" {
		target := stages[i].FeedbackTo
		if target == "" {
			target = stage
		}
		return &Application{Action: "feedback", TargetStage: target, State: "pending"}, nil
	}
	if i+1 >= len(stages) || strings.TrimSpace(stages[i+1].Name) == "" {
		return nil, fmt.Errorf("workflow stage %s has no advance target", stage)
	}
	blockers := []Blocker{}
	return &Application{Action: "advance", TargetStage: stages[i+1].Name, State: "pending", Blockers: &blockers}, nil
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

func verifyAssociation(resultBytes []byte, result *providerResult, association *resultAssociation, actor string, binding Briefing, inventory []presentedItem) error {
	if result.Type != "review-v1-result" || result.Briefing == "" || result.Artifact.ID == "" || !digestRE.MatchString(result.Artifact.Rev) {
		return fmt.Errorf("gate room Result is not a complete review-v1-result")
	}
	if association.Type != "spacedock-result-association" || association.Version != "1" {
		return fmt.Errorf("derived presentation association is not spacedock-result-association v1")
	}
	if association.Result.Digest != RawDigest(resultBytes) || association.Result.Briefing != result.Briefing {
		return fmt.Errorf("association does not bind the exact Result bytes and provider Briefing")
	}
	if association.Actor != actor || result.Resolution.By != actor {
		return fmt.Errorf("Result actor is not authorized by the retained association")
	}
	if err := verifyProviderResolution(result); err != nil {
		return err
	}
	if association.Canonical.Briefing != binding.ID || association.Canonical.Revision != binding.Digest {
		return fmt.Errorf("association does not bind the current canonical Briefing revision")
	}
	if len(inventory) == 0 || len(association.Canonical.Artifacts) != len(inventory) || len(association.Presentation) != len(inventory) {
		return fmt.Errorf("association does not cover the complete presentation mapping")
	}
	canonical := make(map[string]presentedItem, len(inventory))
	for i, artifact := range inventory {
		declared := association.Canonical.Artifacts[i]
		if artifact.ID == "" || !digestRE.MatchString(artifact.Rev) || canonical[artifact.ID].ID != "" || declared.ID != artifact.ID || declared.Rev != artifact.Rev {
			return fmt.Errorf("association canonical artifacts do not match the bound Briefing inventory")
		}
		canonical[artifact.ID] = artifact
	}
	seenCanonical := map[string]bool{}
	seenProvider := map[string]bool{}
	resultArtifactPresent := false
	for _, mapping := range association.Presentation {
		canonicalItem := canonical[mapping.Canonical.ID]
		if mapping.Provider.ID == "" || !digestRE.MatchString(mapping.Provider.Rev) || canonicalItem.Rev != mapping.Canonical.Rev || seenCanonical[mapping.Canonical.ID] || seenProvider[mapping.Provider.ID] {
			return fmt.Errorf("association does not cover the complete presentation mapping")
		}
		seenCanonical[mapping.Canonical.ID] = true
		seenProvider[mapping.Provider.ID] = true
		if canonicalItem.Type == "Artifact" && mapping.Provider.ID == result.Artifact.ID && mapping.Provider.Rev == result.Artifact.Rev {
			resultArtifactPresent = true
		}
	}
	if len(seenCanonical) != len(canonical) || !resultArtifactPresent {
		return fmt.Errorf("association does not cover the exact Result artifact and complete presentation mapping")
	}
	return nil
}

func deriveAssociation(resultBytes []byte, result *providerResult, presented *presentedInventory, actor string, binding Briefing, inventory []presentedItem) (*resultAssociation, error) {
	if len(inventory) == 0 || len(presented.Items) != len(inventory) {
		return nil, fmt.Errorf("presented inventory does not cover the complete presentation mapping")
	}
	canonical := make(map[string]presentedItem, len(inventory))
	for _, item := range inventory {
		if item.ID == "" || !digestRE.MatchString(item.Rev) || canonical[item.ID].ID != "" {
			return nil, fmt.Errorf("canonical Briefing has incomplete or duplicate presentation identity")
		}
		canonical[item.ID] = item
	}
	association := &resultAssociation{Type: "spacedock-result-association", Version: "1", Actor: actor}
	association.Result.Digest = RawDigest(resultBytes)
	association.Result.Briefing = result.Briefing
	association.Canonical.Briefing = binding.ID
	association.Canonical.Revision = binding.Digest
	for _, item := range inventory {
		association.Canonical.Artifacts = append(association.Canonical.Artifacts, item.artifactRef)
	}
	seen := make(map[string]bool, len(presented.Items))
	for _, item := range presented.Items {
		expected := canonical[item.ID]
		if (item.Type != "Artifact" && item.Type != "Reference") || expected.ID == "" ||
			item.Type != expected.Type || item.Rev != expected.Rev || seen[item.ID] {
			return nil, fmt.Errorf("presented inventory does not cover the complete presentation mapping")
		}
		seen[item.ID] = true
		association.Presentation = append(association.Presentation, struct {
			Provider  artifactRef `json:"provider"`
			Canonical artifactRef `json:"canonical"`
		}{
			Provider:  item.artifactRef,
			Canonical: expected.artifactRef,
		})
	}
	return association, nil
}

func verifyProviderResolution(result *providerResult) error {
	if result.Resolution.Briefing != result.Briefing {
		return fmt.Errorf("provider Resolution does not bind its provider Briefing")
	}
	annotations := map[string]string{}
	for _, annotation := range result.Annotations {
		if annotation.Type != "Annotation" || annotation.ID == "" || annotations[annotation.ID] != "" {
			return fmt.Errorf("provider annotations have missing or duplicate identity")
		}
		annotations[annotation.ID] = annotation.Briefing
	}
	for _, included := range result.Resolution.Includes {
		if annotations[included] != result.Briefing {
			return fmt.Errorf("Resolution includes must name a provider Annotation from the same Briefing")
		}
	}
	return validateResolution(&result.Resolution, result.Briefing)
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

func validateGateRoomRequest(briefingPath string, binding Briefing, gateID, attemptID string) error {
	if binding.RequestDigest == "" {
		return nil
	}
	_, requestBytes, request, err := requestForBriefing(briefingPath)
	if err != nil {
		return err
	}
	if request == nil {
		return fmt.Errorf("request-backed Briefing has no request.json")
	}
	requestDigest, err := CanonicalDigest(requestBytes)
	if err != nil {
		return fmt.Errorf("canonicalize gate room request: %w", err)
	}
	if requestDigest != binding.RequestDigest {
		return fmt.Errorf("gate room request does not match the bound request digest")
	}
	if request.Type != "spacedock-gate-presentation-request" || request.Version != "1" ||
		request.Gate != gateID || request.Attempt != attemptID ||
		request.Briefing.ID != binding.ID || request.Briefing.Digest != binding.Digest ||
		request.Actor != "person:captain" || request.Approver != "person:captain" {
		return fmt.Errorf("gate room request does not bind the derived gate, attempt, Briefing, and captain authority")
	}
	return nil
}

func boundBriefingManifest(entityPath string, binding Briefing) (*briefingManifest, error) {
	if binding.DigestDomain != "canonical-bytes" {
		return nil, fmt.Errorf("Result association requires a canonical-bytes Briefing binding")
	}
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
	if binding.RequestDigest != "" {
		if err := validatePreparedSummary(manifest); err != nil {
			return nil, err
		}
	}
	return manifest, nil
}

func boundBriefingBytes(entityPath string, binding Briefing) ([]byte, string, error) {
	retained := filepath.Join(filepath.Dir(entityPath), filepath.FromSlash(binding.RoomRef))
	briefingPath := retained
	if binding.RequestDigest != "" {
		requestBytes, err := os.ReadFile(filepath.Join(retained, "request.json"))
		if err != nil {
			return nil, "", fmt.Errorf("read bound gate room request: %w", err)
		}
		requestDigest, err := CanonicalDigest(requestBytes)
		if err != nil {
			return nil, "", fmt.Errorf("canonicalize bound gate room request: %w", err)
		}
		if requestDigest != binding.RequestDigest {
			return nil, "", fmt.Errorf("bound gate room request does not match the frozen request digest")
		}
		request, err := decodeGateRoomRequest(requestBytes)
		if err != nil {
			return nil, "", err
		}
		briefingPath, err = resolveBriefingLocator(retained, request.Briefing.Locator)
		if err != nil {
			return nil, "", err
		}
		if request.Briefing.ID != binding.ID || request.Briefing.Digest != binding.Digest {
			return nil, "", fmt.Errorf("bound gate room request does not match the frozen Briefing identity and digest")
		}
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

func requestForBriefing(briefingPath string) (string, []byte, *gateRoomRequest, error) {
	briefingPath, err := filepath.Abs(briefingPath)
	if err != nil {
		return "", nil, nil, err
	}
	if resolved, resolveErr := filepath.EvalSymlinks(briefingPath); resolveErr == nil {
		briefingPath = resolved
	}
	for room := filepath.Dir(briefingPath); ; room = filepath.Dir(room) {
		requestPath := filepath.Join(room, "request.json")
		requestBytes, readErr := os.ReadFile(requestPath)
		if readErr == nil {
			locators := requestLocatorCandidates(requestBytes)
			for _, locator := range locators {
				located, locateErr := resolveBriefingLocator(room, locator)
				if locateErr == nil && filepath.Clean(located) == filepath.Clean(briefingPath) {
					request, err := decodeGateRoomRequest(requestBytes)
					if err != nil {
						return "", nil, nil, err
					}
					return room, requestBytes, request, nil
				}
			}
		}
		if readErr != nil && !os.IsNotExist(readErr) {
			return "", nil, nil, fmt.Errorf("read gate room request: %w", readErr)
		}
		parent := filepath.Dir(room)
		if parent == room {
			return "", nil, nil, nil
		}
	}
}

func requestLocatorCandidates(data []byte) []string {
	var locators []string
	for _, briefing := range duplicatePreservingObjectValues(data, "briefing") {
		for _, raw := range duplicatePreservingObjectValues(briefing, "locator") {
			var locator string
			if json.Unmarshal(raw, &locator) == nil {
				locators = append(locators, locator)
			}
		}
	}
	return locators
}

// duplicatePreservingObjectValues is only a membership probe. Matching request
// bytes still pass through decodeGateRoomRequest's strict duplicate rejection.
func duplicatePreservingObjectValues(data []byte, want string) []json.RawMessage {
	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		return nil
	}
	var values []json.RawMessage
	for decoder.More() {
		keyToken, err := decoder.Token()
		if err != nil {
			return values
		}
		key, ok := keyToken.(string)
		if !ok {
			return values
		}
		var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			return values
		}
		if key == want {
			values = append(values, raw)
		}
	}
	return values
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
