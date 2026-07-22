// ABOUTME: Recorder operations, pointer-CAS checks, digest binding, and result adoption.
// ABOUTME: These operations never model application state or invoke workflow effects.
package gates

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	jsoncanonicalizer "github.com/cyberphone/json-canonicalization/go/src/webpki.org/jsoncanonicalizer"
	"gopkg.in/yaml.v3"
)

type RecordInput struct {
	BriefingPath    string
	ResultPath      string
	AssociationPath string
	Actor           string
	AdoptionNote    string
	Decision        string
	Reason          string
	Directive       string
}

type artifactRef struct {
	ID  string `json:"id"`
	URI string `json:"uri,omitempty"`
	Rev string `json:"rev"`
}

type briefingManifest struct {
	Type      string        `json:"type"`
	Version   string        `json:"version"`
	ID        string        `json:"id"`
	Question  string        `json:"question"`
	Artifacts []artifactRef `json:"artifacts"`
}

type providerResult struct {
	Type        string      `json:"type"`
	Status      string      `json:"status"`
	Briefing    string      `json:"briefing"`
	Artifact    artifactRef `json:"artifact"`
	Resolution  Resolution  `json:"resolution"`
	Annotations []struct {
		Type     string `json:"type"`
		ID       string `json:"id"`
		Briefing string `json:"briefing"`
	} `json:"annotations"`
	Binding bool   `json:"binding"`
	Actor   string `json:"actor"`
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
	sources := 0
	for _, source := range []string{input.BriefingPath, input.ResultPath, input.Decision} {
		if source != "" {
			sources++
		}
	}
	if sources != 1 {
		return fmt.Errorf("gate record requires exactly one of --briefing, --result, or --decision")
	}
	if input.BriefingPath != "" && (input.ResultPath != "" || input.AssociationPath != "" || input.Actor != "" || input.AdoptionNote != "" || input.Decision != "" || input.Reason != "" || input.Directive != "") || input.ResultPath != "" && (input.Decision != "" || input.Reason != "" || input.Directive != "") || input.Decision != "" && (input.AssociationPath != "" || input.AdoptionNote != "") {
		return fmt.Errorf("gate record flags do not match the selected semantic source")
	}
	if input.BriefingPath != "" && filepath.Base(input.BriefingPath) != "briefing.json" {
		return fmt.Errorf("--briefing must name a canonical package manifest named briefing.json")
	}
	unlock, err := lockEntity(entityPath)
	if err != nil {
		return err
	}
	defer unlock()
	if input.BriefingPath != "" {
		return recordBriefingLocked(entityPath, input.BriefingPath)
	}
	if input.ResultPath != "" {
		return recordResultLocked(entityPath, input)
	}
	return recordChatLocked(entityPath, input)
}

func RecordBriefing(entityPath, briefingPath string) error {
	return RecordSemantic(entityPath, RecordInput{BriefingPath: briefingPath})
}

func recordBriefingLocked(entityPath, briefingPath string) error {
	binding, err := bindingFromManifest(entityPath, briefingPath)
	if err != nil {
		return err
	}
	stage, err := entityStatus(entityPath)
	if err != nil {
		return err
	}
	doc, oldNode, readErr := Read(entityPath)
	if readErr != nil && !strings.Contains(readErr.Error(), "no gates record") {
		return readErr
	}
	if doc == nil {
		doc = &Document{Version: 1}
	}
	record, lookupErr := recordForStage(doc, stage)
	if lookupErr != nil && !strings.Contains(lookupErr.Error(), "no logical gate") {
		return lookupErr
	}
	if record == nil {
		gateID, attemptID, err := initialIDs(binding.ID, entityPath, stage)
		if err != nil {
			return err
		}
		record = &GateRecord{ID: gateID, Stage: stage, Attempts: []Attempt{{ID: attemptID, Briefing: binding}}}
		doc.Records = append(doc.Records, *record)
		doc.Current = Selection{Gate: gateID}
		return writeDocument(entityPath, oldNode, doc)
	}
	previous := &record.Attempts[len(record.Attempts)-1]
	if previous.Resolution == nil {
		if sameBinding(previous.Briefing, binding) {
			return nil
		}
		previous.Briefing = binding
		if err := Validate(doc); err != nil {
			return err
		}
		if err := ValidateTransition(oldNode, doc); err != nil {
			return err
		}
		return writeDocument(entityPath, oldNode, doc)
	}
	nextID, err := successorAttemptID(previous.ID)
	if err != nil {
		return err
	}
	next := Attempt{ID: nextID, Briefing: binding}
	record.Attempts = append(record.Attempts, next)
	doc.Current.Gate = record.ID
	if err := Validate(doc); err != nil {
		return err
	}
	if err := ValidateTransition(oldNode, doc); err != nil {
		return err
	}
	return writeDocument(entityPath, oldNode, doc)
}

func recordResultLocked(entityPath string, input RecordInput) error {
	if input.AssociationPath == "" || input.Actor == "" {
		return fmt.Errorf("--result requires --association FILE and --actor ID")
	}
	doc, oldNode, record, attempt, err := currentStageAttempt(entityPath)
	if err != nil {
		return err
	}
	if attempt.Resolution != nil {
		return fmt.Errorf("attempt %s is frozen closed", attempt.ID)
	}
	if !digestRE.MatchString(attempt.Briefing.Digest) {
		return fmt.Errorf("open attempt %s has no verifiable digest", attempt.ID)
	}
	resultBytes, err := os.ReadFile(input.ResultPath)
	if err != nil {
		return err
	}
	associationBytes, err := os.ReadFile(input.AssociationPath)
	if err != nil {
		return err
	}
	var result providerResult
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		return fmt.Errorf("parse Result: %w", err)
	}
	var association resultAssociation
	if err := json.Unmarshal(associationBytes, &association); err != nil {
		return fmt.Errorf("parse result association: %w", err)
	}
	manifest, err := boundBriefingManifest(entityPath, attempt.Briefing)
	if err != nil {
		return err
	}
	if err := verifyAssociation(resultBytes, &result, &association, input.Actor, attempt.Briefing, manifest.Artifacts); err != nil {
		return err
	}
	if !result.Binding && strings.TrimSpace(input.AdoptionNote) == "" {
		return fmt.Errorf("advisory Result requires --adoption-note naming its authorizer")
	}
	resolution := result.Resolution
	resolution.Briefing = attempt.Briefing.ID
	resolution.Adoption = input.AdoptionNote
	if err := closeAttempt(doc, oldNode, record, attempt, &resolution); err != nil {
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
	delegated := strings.HasPrefix(input.Actor, "agent:")
	if delegated && strings.TrimSpace(input.Directive) == "" {
		return fmt.Errorf("delegated chat decision requires --directive")
	}
	if input.Actor == "agent:first-officer" && input.Decision == "approve" && strings.TrimSpace(input.Reason) == "" {
		return fmt.Errorf("delegated First Officer approval requires --reason")
	}
	doc, oldNode, record, attempt, err := currentStageAttempt(entityPath)
	if err != nil {
		return err
	}
	if attempt.Resolution != nil {
		return fmt.Errorf("attempt %s is frozen closed", attempt.ID)
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
		Adoption: input.Directive,
	}
	if err := closeAttempt(doc, oldNode, record, attempt, resolution); err != nil {
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
	return left.ID == right.ID && left.Digest == right.Digest && left.DigestDomain == right.DigestDomain && left.RoomRef == right.RoomRef
}

func currentStageAttempt(entityPath string) (*Document, *yaml.Node, *GateRecord, *Attempt, error) {
	doc, oldNode, err := Read(entityPath)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	stage, err := entityStatus(entityPath)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	record, err := recordForStage(doc, stage)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	attempt := &record.Attempts[len(record.Attempts)-1]
	return doc, oldNode, record, attempt, nil
}

func closeAttempt(doc *Document, oldNode *yaml.Node, record *GateRecord, attempt *Attempt, resolution *Resolution) error {
	if err := validateResolution(resolution, attempt.Briefing.ID); err != nil {
		return err
	}
	attempt.Resolution = resolution
	doc.Current.Gate = record.ID
	if err := Validate(doc); err != nil {
		return err
	}
	return ValidateTransition(oldNode, doc)
}

func verifyAssociation(resultBytes []byte, result *providerResult, association *resultAssociation, actor string, binding Briefing, inventory []artifactRef) error {
	if result.Type != "review-v1-result" || result.Briefing == "" || result.Artifact.ID == "" || !digestRE.MatchString(result.Artifact.Rev) {
		return fmt.Errorf("--result is not a complete review-v1-result")
	}
	if association.Type != "spacedock-result-association" || association.Version != "1" {
		return fmt.Errorf("--association is not a spacedock-result-association v1")
	}
	if association.Result.Digest != RawDigest(resultBytes) || association.Result.Briefing != result.Briefing {
		return fmt.Errorf("association does not bind the exact Result bytes and provider Briefing")
	}
	if association.Actor != actor || result.Actor != actor || result.Resolution.By != actor {
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
	canonical := make(map[string]string, len(inventory))
	for i, artifact := range inventory {
		declared := association.Canonical.Artifacts[i]
		if artifact.ID == "" || !digestRE.MatchString(artifact.Rev) || canonical[artifact.ID] != "" || declared.ID != artifact.ID || declared.Rev != artifact.Rev {
			return fmt.Errorf("association canonical artifacts do not match the bound Briefing inventory")
		}
		canonical[artifact.ID] = artifact.Rev
	}
	seenCanonical := map[string]bool{}
	seenProvider := map[string]bool{}
	resultArtifactPresent := false
	for _, mapping := range association.Presentation {
		if mapping.Provider.ID == "" || !digestRE.MatchString(mapping.Provider.Rev) || canonical[mapping.Canonical.ID] != mapping.Canonical.Rev || seenCanonical[mapping.Canonical.ID] || seenProvider[mapping.Provider.ID] {
			return fmt.Errorf("association does not cover the complete presentation mapping")
		}
		seenCanonical[mapping.Canonical.ID] = true
		seenProvider[mapping.Provider.ID] = true
		if mapping.Provider.ID == result.Artifact.ID && mapping.Provider.Rev == result.Artifact.Rev {
			resultArtifactPresent = true
		}
	}
	if len(seenCanonical) != len(canonical) || !resultArtifactPresent {
		return fmt.Errorf("association does not cover the exact Result artifact and complete presentation mapping")
	}
	return nil
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

func initialIDs(briefingID, entityPath, stage string) (string, string, error) {
	if strings.HasPrefix(briefingID, "briefing:") {
		body := strings.TrimPrefix(briefingID, "briefing:")
		if marker := strings.Index(body, ":attempt-"); marker > 0 {
			base := body[:marker]
			parts := strings.Split(base, ":")
			attemptNumber := strings.SplitN(body[marker+len(":attempt-"):], ":", 2)[0]
			if len(parts) >= 2 && parts[len(parts)-1] == stage {
				if _, err := strconv.Atoi(attemptNumber); err == nil {
					entity := parts[len(parts)-2]
					return "gate:" + base, "gate-attempt:" + entity + "-" + stage + "-1", nil
				}
			}
		}
	}
	entity := strings.TrimSuffix(filepath.Base(entityPath), filepath.Ext(entityPath))
	if entity == "index" {
		entity = filepath.Base(filepath.Dir(entityPath))
	}
	entity = strings.Trim(strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' {
			return r
		}
		return '-'
	}, entity), "-")
	if entity == "" || stage == "" {
		return "", "", fmt.Errorf("cannot derive gate identity from entity and Briefing")
	}
	return "gate:" + entity + ":" + stage, "gate-attempt:" + entity + "-" + stage + "-1", nil
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

func bindingFromManifest(entityPath, briefingPath string) (Briefing, error) {
	data, err := os.ReadFile(briefingPath)
	if err != nil {
		return Briefing{}, err
	}
	manifest, err := parseBriefingManifest(data)
	if err != nil {
		return Briefing{}, err
	}
	digest, err := CanonicalDigest(data)
	if err != nil {
		return Briefing{}, fmt.Errorf("canonicalize briefing: %w", err)
	}
	roomRef, err := filepath.Rel(filepath.Dir(entityPath), filepath.Dir(briefingPath))
	if err != nil {
		return Briefing{}, fmt.Errorf("resolve briefing room: %w", err)
	}
	roomRef = filepath.ToSlash(roomRef)
	if roomRef == "." {
		roomRef = "./"
	} else if !strings.HasPrefix(roomRef, ".") {
		roomRef = "./" + roomRef
	}
	return Briefing{ID: manifest.ID, Digest: digest, DigestDomain: "canonical-bytes", RoomRef: roomRef}, nil
}

func boundBriefingManifest(entityPath string, binding Briefing) (*briefingManifest, error) {
	if binding.DigestDomain != "canonical-bytes" {
		return nil, fmt.Errorf("Result association requires a canonical-bytes Briefing binding")
	}
	path := filepath.Join(filepath.Dir(entityPath), filepath.FromSlash(binding.RoomRef), "briefing.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read bound canonical Briefing: %w", err)
	}
	digest, err := CanonicalDigest(data)
	if err != nil {
		return nil, fmt.Errorf("canonicalize bound Briefing: %w", err)
	}
	if digest != binding.Digest {
		return nil, fmt.Errorf("bound canonical Briefing bytes do not match the frozen digest")
	}
	manifest, err := parseBriefingManifest(data)
	if err != nil {
		return nil, err
	}
	if manifest.ID != binding.ID {
		return nil, fmt.Errorf("bound canonical Briefing identity does not match the current binding")
	}
	return manifest, nil
}

func parseBriefingManifest(data []byte) (*briefingManifest, error) {
	var manifest briefingManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse Briefing: %w", err)
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
			if attemptState(&oldAttempt) != "closed" {
				continue
			}
			var found *Attempt
			for i := range nr.Attempts {
				if nr.Attempts[i].ID == oldAttempt.ID {
					found = &nr.Attempts[i]
				}
			}
			if found == nil || !nodesEqual(oldAttempt, *found) {
				return fmt.Errorf("frozen closed attempt %s cannot be deleted or mutated", oldAttempt.ID)
			}
		}
	}
	return nil
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
	canonical, err := jsoncanonicalizer.Transform(data)
	if err != nil {
		return "", err
	}
	return RawDigest(canonical), nil
}
