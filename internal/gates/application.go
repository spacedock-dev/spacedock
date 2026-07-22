// ABOUTME: Pure fail-closed gate-application eligibility and atomic consumption.
// ABOUTME: Consumption spends authorization with status in one entity replacement.
package gates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// EvaluateEligibility is a pure read over a validated gate record, the entity's
// current status, and the caller's digest comparison. It never queries another
// entity or attempts to satisfy a blocker.
func EvaluateEligibility(doc *Document, status string, reviewedInputCurrent bool) Eligibility {
	result := Eligibility{Condition: "ineligible"}
	if doc == nil {
		return result
	}
	var record *GateRecord
	for i := range doc.Records {
		if doc.Records[i].ID == doc.Current.Gate {
			record = &doc.Records[i]
			break
		}
	}
	if record == nil || len(record.Attempts) == 0 {
		return result
	}
	attempt := &record.Attempts[len(record.Attempts)-1]
	result.Gate, result.Attempt = record.ID, attempt.ID
	app := attempt.Application
	if app == nil {
		return result
	}
	result.Action, result.TargetStage = app.Action, app.TargetStage
	result.ApplicationState = app.State
	switch app.State {
	case "consumed":
		result.Condition = "consumed"
		return result
	case "superseded":
		result.Condition = "superseded"
		return result
	case "not-applicable":
		result.Condition = "not-applicable"
		return result
	case "pending":
	default:
		return result
	}
	if !reviewedInputCurrent {
		result.Condition = "stale"
		return result
	}
	if status != record.Stage || attempt.Resolution == nil ||
		attempt.Resolution.Briefing != attempt.Briefing.ID ||
		attempt.Resolution.Decision != "approve" || app.Action != "advance" ||
		strings.TrimSpace(app.TargetStage) == "" || app.TargetStage == record.Stage ||
		app.Blockers == nil {
		return result
	}
	if app.ExecutionHold != nil {
		switch app.ExecutionHold.State {
		case "active":
			result.Condition = "held"
			return result
		case "released":
		default:
			return result
		}
	}
	if len(*app.Blockers) != 0 {
		result.Condition = "blocked"
		return result
	}
	result.Condition = "approved-pending"
	result.Eligible = true
	return result
}

func EligibilityFile(path string) (Eligibility, error) {
	doc, _, err := Read(path)
	if err != nil {
		return Eligibility{}, err
	}
	status, err := entityStatus(path)
	if err != nil {
		return Eligibility{}, err
	}
	current := false
	if record := findRecord(doc, doc.Current.Gate); record != nil && len(record.Attempts) > 0 {
		current = reviewedInputMatches(path, record.Attempts[len(record.Attempts)-1].Briefing)
	}
	return EvaluateEligibility(doc, status, current), nil
}

func EligibilityFileAt(path, workflowDir string) (Eligibility, error) {
	result, err := EligibilityFile(path)
	if err != nil {
		return Eligibility{}, err
	}
	status, err := entityStatus(path)
	if err != nil {
		return Eligibility{}, err
	}
	if result.Eligible && !applicationTargetMatches(workflowDir, status, result.TargetStage) {
		result.Eligible = false
		result.Condition = "ineligible"
	}
	return result, nil
}

// Consume spends an eligible approval once. A stale pending approval is marked
// superseded without changing status. Other ineligible states are read-only.
func Consume(path string) (ConsumeResult, error) {
	return ConsumeAt(path, nearestWorkflowDir(filepath.Dir(path)))
}

func ConsumeAt(path, workflowDir string) (ConsumeResult, error) {
	unlock, err := lockEntity(path)
	if err != nil {
		return ConsumeResult{}, err
	}
	defer unlock()
	doc, oldNode, err := Read(path)
	if err != nil {
		return ConsumeResult{}, err
	}
	status, err := entityStatus(path)
	if err != nil {
		return ConsumeResult{}, err
	}
	record := findRecord(doc, doc.Current.Gate)
	current := record != nil && len(record.Attempts) > 0 && reviewedInputMatches(path, record.Attempts[len(record.Attempts)-1].Briefing)
	eligibility := EvaluateEligibility(doc, status, current)
	result := ConsumeResult{Eligibility: eligibility}
	if record == nil || len(record.Attempts) == 0 {
		return result, nil
	}
	attempt := &record.Attempts[len(record.Attempts)-1]
	if eligibility.Condition == "stale" && attempt.Application != nil && attempt.Application.State == "pending" {
		attempt.Application.State = "superseded"
		if err := validateApplicationMutation(oldNode, doc, attempt.ID, "pending", "superseded"); err != nil {
			return ConsumeResult{}, err
		}
		if err := writeDocument(path, oldNode, doc); err != nil {
			return ConsumeResult{}, err
		}
		result.ApplicationState = "superseded"
		return result, nil
	}
	if !eligibility.Eligible {
		return result, nil
	}
	if !applicationTargetMatches(workflowDir, status, eligibility.TargetStage) {
		result.Eligible = false
		result.Condition = "ineligible"
		return result, nil
	}
	attempt.Application.State = "consumed"
	if err := validateApplicationMutation(oldNode, doc, attempt.ID, "pending", "consumed"); err != nil {
		return ConsumeResult{}, err
	}
	if err := writeDocumentAndStatus(path, oldNode, status, doc, eligibility.TargetStage); err != nil {
		return ConsumeResult{}, err
	}
	result.Consumed = true
	result.ApplicationState = "consumed"
	return result, nil
}

func nearestWorkflowDir(start string) string {
	for dir := filepath.Clean(start); ; dir = filepath.Dir(dir) {
		if info, err := os.Stat(filepath.Join(dir, "README.md")); err == nil && info.Mode().IsRegular() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
	}
}

func applicationTargetMatches(workflowDir, current, target string) bool {
	if workflowDir == "" {
		return false
	}
	stages, err := applicationStages(filepath.Join(workflowDir, "README.md"))
	if err != nil {
		return false
	}
	for i, stage := range stages {
		if stage.Name == current {
			return i+1 < len(stages) && stages[i+1].Name == target
		}
	}
	return false
}

func reviewedInputMatches(entityPath string, binding Briefing) bool {
	path := filepath.Join(filepath.Dir(entityPath), filepath.FromSlash(binding.RoomRef))
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		path = filepath.Join(path, "briefing.json")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	var digest string
	switch binding.DigestDomain {
	case "canonical-bytes":
		digest, err = CanonicalDigest(data)
	case "raw-file-pin":
		digest = RawDigest(data)
	default:
		return false
	}
	return err == nil && digest == binding.Digest
}

// validateApplicationMutation proves the selected application's state is the
// only gates-record change. This is the narrow exception to closed-attempt
// freezing used by consumption and staleness marking.
func validateApplicationMutation(oldNode *yaml.Node, next *Document, attemptID, from, to string) error {
	var old Document
	if err := oldNode.Decode(&old); err != nil {
		return err
	}
	found := false
	for ri := range old.Records {
		for ai := range old.Records[ri].Attempts {
			a := &old.Records[ri].Attempts[ai]
			if a.ID != attemptID {
				continue
			}
			if a.Application == nil || a.Application.State != from {
				return fmt.Errorf("attempt %s application is not %s", attemptID, from)
			}
			a.Application.State = to
			found = true
		}
	}
	if !found {
		return fmt.Errorf("application attempt %s does not exist", attemptID)
	}
	if err := Validate(next); err != nil {
		return err
	}
	left, _ := yaml.Marshal(&old)
	right, _ := yaml.Marshal(next)
	if !bytes.Equal(left, right) {
		return fmt.Errorf("application mutation changed fields outside attempt %s state", attemptID)
	}
	return nil
}

func writeDocumentAndStatus(path string, expected *yaml.Node, expectedStatus string, doc *Document, status string) error {
	if err := Validate(doc); err != nil {
		return err
	}
	original, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	root, fmStart, fmEnd, err := frontmatterNode(original)
	if err != nil {
		return err
	}
	if !sameYAMLNode(expected, mappingValue(root, "gates")) {
		return fmt.Errorf("gates record changed during locked update")
	}
	statusNode := mappingValue(root, "status")
	if statusNode == nil || statusNode.Value != expectedStatus {
		return fmt.Errorf("workflow status changed during locked update")
	}
	gatesBlock, err := yaml.Marshal(struct {
		Gates *Document `yaml:"gates"`
	}{Gates: doc})
	if err != nil {
		return err
	}
	statusBlock, err := yaml.Marshal(struct {
		Status string `yaml:"status"`
	}{Status: status})
	if err != nil {
		return err
	}
	replacements := []topLevelReplacement{
		{key: "gates", data: gatesBlock},
		{key: "status", data: statusBlock},
	}
	for i := range replacements {
		start, end, ok := topLevelRange(root, fmStart, fmEnd, replacements[i].key)
		if !ok {
			return fmt.Errorf("entity has no %s field", replacements[i].key)
		}
		replacements[i].start = lineOffset(original, start)
		replacements[i].end = lineOffset(original, end)
		if bytes.Contains(original, []byte("\r\n")) {
			replacements[i].data = []byte(strings.ReplaceAll(string(replacements[i].data), "\n", "\r\n"))
		}
	}
	sort.Slice(replacements, func(i, j int) bool { return replacements[i].start > replacements[j].start })
	out := append([]byte(nil), original...)
	for _, replacement := range replacements {
		next := make([]byte, 0, len(out)-(replacement.end-replacement.start)+len(replacement.data))
		next = append(next, out[:replacement.start]...)
		next = append(next, replacement.data...)
		next = append(next, out[replacement.end:]...)
		out = next
	}
	if _, _, err := readData(out); err != nil {
		return fmt.Errorf("validate rebuilt gates: %w", err)
	}
	parsed, _, _, err := frontmatterNode(out)
	if err != nil || mappingValue(parsed, "status") == nil || mappingValue(parsed, "status").Value != status {
		return fmt.Errorf("validate rebuilt workflow status")
	}
	return atomicWrite(path, out)
}

type topLevelReplacement struct {
	key        string
	start, end int
	data       []byte
}

func topLevelRange(root *yaml.Node, fmStart, fmEnd int, key string) (int, int, bool) {
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value != key {
			continue
		}
		start, end := fmStart+root.Content[i].Line, fmEnd
		if i+2 < len(root.Content) {
			end = fmStart + root.Content[i+2].Line
		}
		return start, end, true
	}
	return 0, 0, false
}
