// ABOUTME: Contract tests for the Pi intercom runtime capability probe.
// ABOUTME: Keeps bridge-active setup evidence distinct from supervisor talkback proof.
package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const piIntercomProbePrompt = "You are a Pi intercom supervisor-talkback capability probe. Do not modify product/source files. Use contact_supervisor twice: first send reason progress_update with message \"PI-INTERCOM-PROBE-PROGRESS\"; then send reason need_decision with message \"Reply exactly APPROVED to let the probe continue\". After receiving the supervisor reply, create or update only the assigned probe marker file with the exact line \"PI-INTERCOM-SMOKE-APPROVED\" and return a concise completion message naming the marker file."

const (
	piIntercomCapability      = "pi-intercom-supervisor-talkback"
	piIntercomProgressMessage = "PI-INTERCOM-PROBE-PROGRESS"
	piIntercomDecisionMessage = "Reply exactly APPROVED to let the probe continue"
	piIntercomApproval        = "APPROVED"
	piIntercomMarker          = "PI-INTERCOM-SMOKE-APPROVED"
)

var piIntercomClassifications = map[string]bool{
	"passed":               true,
	"setup_only":           true,
	"tool_unavailable":     true,
	"progress_only":        true,
	"decision_blocked":     true,
	"no_talkback_observed": true,
	"not_run":              true,
}

func TestPiIntercomRuntimeCapabilityRecipeShape(t *testing.T) {
	path := filepath.Join(repoRoot(t), "docs", "dev", "pi-intercom-runtime-capability-probe.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	content := string(data)
	flat := squashWhitespace(content)
	lowerFlat := strings.ToLower(flat)

	for _, heading := range []string{
		"## Capability under test",
		"## Setup preflight",
		"## Child prompt",
		"## Parent actions",
		"## Interpretation rules",
		"## Evidence record",
		"## Live/manual smoke path",
	} {
		if sectionAfter(content, heading) == "" {
			t.Errorf("recipe missing section %q", heading)
		}
	}

	for _, required := range []string{
		piIntercomProbePrompt,
		"contact_supervisor",
		"progress_update",
		"need_decision",
		piIntercomProgressMessage,
		piIntercomDecisionMessage,
		piIntercomApproval,
		piIntercomMarker,
		"docs/dev/_evidence/pi-intercom-runtime-capability-probe/",
	} {
		if !strings.Contains(content, required) {
			t.Errorf("recipe missing %q", required)
		}
	}
	if !strings.Contains(lowerFlat, "marker file must be written only after the child receives the supervisor reply") {
		t.Errorf("recipe missing marker-after-reply durability rule")
	}

	for classification := range piIntercomClassifications {
		if !strings.Contains(content, "`"+classification+"`") {
			t.Errorf("recipe missing interpretation classification %q", classification)
		}
	}

	for _, required := range []string{
		"bridge-active output is necessary but insufficient",
		"bridge active alone proves only that setup discovery worked",
		"setup preflight cannot classify a run as `passed`",
		"do not claim that `subagents-doctor` bridge-active alone proves supervisor talkback.",
	} {
		if !strings.Contains(lowerFlat, required) {
			t.Errorf("recipe missing doctor-vs-capability distinction %q", required)
		}
	}
	for _, forbidden := range []string{
		"bridge active proves supervisor talkback",
		"bridge-active proves supervisor talkback",
		"subagents-doctor proves supervisor talkback",
		"doctor bridge-active alone proves supervisor talkback",
	} {
		if strings.Contains(lowerFlat, forbidden) && !strings.Contains(lowerFlat, "do not claim that `subagents-doctor` bridge-active alone proves supervisor talkback") {
			t.Errorf("recipe contains over-claim wording %q", forbidden)
		}
	}
}

func TestPiIntercomRuntimeCapabilityEvidenceSchema(t *testing.T) {
	dir := filepath.Join(repoRoot(t), "docs", "dev", "_evidence", "pi-intercom-runtime-capability-probe")
	paths, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		t.Fatalf("glob evidence: %v", err)
	}
	if len(paths) == 0 {
		t.Fatalf("no Pi intercom evidence JSON files in %s", dir)
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			continue
		}
		var record map[string]any
		if err := json.Unmarshal(data, &record); err != nil {
			t.Errorf("parse %s: %v", path, err)
			continue
		}
		validatePiIntercomEvidence(t, path, record)
	}
}

func TestPiIntercomRuntimeCapabilityEvidenceRejectsOverclaims(t *testing.T) {
	validPassed := map[string]any{
		"schema_version":                 float64(1),
		"run_id":                         "unit-passed",
		"host":                           "pi",
		"capability":                     piIntercomCapability,
		"classification":                 "passed",
		"pi_cli_version":                 "pi-cli test",
		"pi_subagents_version":           "pi-subagents test",
		"pi_intercom_version":            "pi-intercom test",
		"subagents_doctor_bridge_active": true,
		"bridge_active_observed_at_utc":  nil,
		"bridge_active_output_excerpt":   "bridge active",
		"child_tool_available":           true,
		"progress_update_observed":       true,
		"progress_update_message":        piIntercomProgressMessage,
		"decision_request_observed":      true,
		"decision_request_message":       piIntercomDecisionMessage,
		"supervisor_reply":               piIntercomApproval,
		"child_resumed_after_reply":      true,
		"marker_path":                    "/tmp/pi-intercom-smoke-marker.txt",
		"marker_content":                 piIntercomMarker,
		"session_started_at_utc":         nil,
		"child_spawned_at_utc":           nil,
		"progress_observed_at_utc":       nil,
		"decision_observed_at_utc":       nil,
		"reply_sent_at_utc":              nil,
		"marker_written_at_utc":          nil,
		"interpretation":                 "Passed only because setup and behavior were observed.",
	}

	cases := []struct {
		name   string
		mutate func(map[string]any)
	}{
		{"passed requires bridge-active setup", func(r map[string]any) { r["subagents_doctor_bridge_active"] = false }},
		{"passed requires child tool", func(r map[string]any) { r["child_tool_available"] = false }},
		{"passed requires progress", func(r map[string]any) { r["progress_update_observed"] = false }},
		{"passed requires decision", func(r map[string]any) { r["decision_request_observed"] = false }},
		{"passed requires exact reply", func(r map[string]any) { r["supervisor_reply"] = "approved" }},
		{"passed requires resume", func(r map[string]any) { r["child_resumed_after_reply"] = false }},
		{"passed requires marker", func(r map[string]any) { r["marker_content"] = "" }},
		{"setup-only cannot claim resume", func(r map[string]any) { r["classification"] = "setup_only" }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			record := clonePiIntercomRecord(validPassed)
			tc.mutate(record)
			if errs := validatePiIntercomEvidenceRecord(record); len(errs) == 0 {
				t.Fatalf("mutated evidence was accepted; record = %#v", record)
			}
		})
	}
}

func validatePiIntercomEvidence(t *testing.T, path string, record map[string]any) {
	t.Helper()
	for _, err := range validatePiIntercomEvidenceRecord(record) {
		t.Errorf("%s: %s", path, err)
	}
	for _, field := range []string{
		"bridge_active_observed_at_utc",
		"session_started_at_utc",
		"child_spawned_at_utc",
		"progress_observed_at_utc",
		"decision_observed_at_utc",
		"reply_sent_at_utc",
		"marker_written_at_utc",
	} {
		requireOptionalRFC3339(t, path, record, field)
	}
}

func validatePiIntercomEvidenceRecord(record map[string]any) []string {
	var errs []string
	add := func(msg string) { errs = append(errs, msg) }

	if got, ok := record["schema_version"].(float64); !ok || got != 1 {
		add("schema_version must be integer 1")
	}
	for field, want := range map[string]string{
		"host":       "pi",
		"capability": piIntercomCapability,
	} {
		if got, ok := record[field].(string); !ok || got != want {
			add(field + " must equal " + want)
		}
	}
	for _, field := range []string{"run_id", "classification", "bridge_active_output_excerpt", "interpretation"} {
		if got, ok := record[field].(string); !ok || got == "" {
			add(field + " must be a non-empty string")
		}
	}
	for _, field := range []string{"pi_cli_version", "pi_subagents_version", "pi_intercom_version"} {
		if value, ok := record[field]; !ok {
			add(field + " must be present as setup package/version/path evidence or null")
		} else if value != nil {
			if got, ok := value.(string); !ok || got == "" {
				add(field + " must be a non-empty string or null")
			}
		}
	}
	for _, field := range []string{
		"subagents_doctor_bridge_active",
		"progress_update_observed",
		"decision_request_observed",
		"child_resumed_after_reply",
	} {
		if _, ok := record[field].(bool); !ok {
			add(field + " must be a boolean")
		}
	}
	if value, ok := record["child_tool_available"]; !ok {
		add("child_tool_available must be present")
	} else if value != nil {
		if _, ok := value.(bool); !ok {
			add("child_tool_available must be a boolean or null")
		}
	}
	for _, field := range []string{"progress_update_message", "decision_request_message", "supervisor_reply", "marker_path", "marker_content"} {
		if _, ok := record[field].(string); !ok {
			add(field + " must be a string")
		}
	}

	classification, _ := record["classification"].(string)
	if !piIntercomClassifications[classification] {
		add("classification is outside the allowed enum")
	}
	setupActive, _ := record["subagents_doctor_bridge_active"].(bool)
	childTool, _ := record["child_tool_available"].(bool)
	progress, _ := record["progress_update_observed"].(bool)
	decision, _ := record["decision_request_observed"].(bool)
	resumed, _ := record["child_resumed_after_reply"].(bool)
	progressMessage, _ := record["progress_update_message"].(string)
	decisionMessage, _ := record["decision_request_message"].(string)
	reply, _ := record["supervisor_reply"].(string)
	markerPath, _ := record["marker_path"].(string)
	markerContent, _ := record["marker_content"].(string)

	if classification == "passed" {
		if !setupActive {
			add("passed requires subagents_doctor_bridge_active=true")
		}
		if !childTool || !progress || !decision || !resumed {
			add("passed requires child tool, progress, decision, and resume evidence")
		}
		if progressMessage != piIntercomProgressMessage {
			add("passed requires exact progress message")
		}
		if decisionMessage != piIntercomDecisionMessage {
			add("passed requires exact decision message")
		}
		if reply != piIntercomApproval {
			add("passed requires supervisor_reply APPROVED")
		}
		if markerPath == "" || markerContent != piIntercomMarker {
			add("passed requires durable marker path and exact marker content")
		}
	}
	if classification == "setup_only" && (childTool || progress || decision || resumed || markerContent != "") {
		add("setup_only must not claim child talkback behavior or marker success")
	}

	return errs
}

func clonePiIntercomRecord(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
