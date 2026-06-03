// ABOUTME: Contract tests for Codex no-wait completion evidence and runtime wording.
// ABOUTME: Keeps queued notifications distinct from autonomous idle wake-up.
package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

const codexIdleNotificationProbePrompt = "You are a no-write Codex idle-wake probe. Do not read or write files, do not run tools, and do not modify state. Sleep for 30 seconds, then reply exactly: Done: idle-wake Codex notification probe completed."

var codexIdleNotificationClassifications = map[string]bool{
	"foreground_wait":          true,
	"queued_flush":             true,
	"autonomous_idle_wake":     true,
	"no_notification_observed": true,
}

func TestCodexIdleNotificationRuntimeContract(t *testing.T) {
	root := skillsRoot(t)
	path := filepath.Join(root, "first-officer", "references", "codex-first-officer-runtime.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	region := sectionAfter(string(data), "## Awaiting Completion")
	if region == "" {
		t.Fatal("Codex runtime missing `## Awaiting Completion` section")
	}

	for _, required := range []string{
		"### Foreground wait",
		"### Queued notification flushed by later activity",
		"### Autonomous idle FO wake-up",
	} {
		if !strings.Contains(region, required) {
			t.Errorf("Awaiting Completion missing outcome heading %q", required)
		}
	}

	scheduling := paragraphContaining(region, "Before calling `wait_agent`")
	if scheduling == "" {
		t.Fatal("Awaiting Completion missing scheduling priority paragraph before `wait_agent`")
	}
	scheduling = squashWhitespace(scheduling)
	for _, required := range []string{
		"ready final-status notifications",
		"gate decisions",
		"state transitions",
		"newly dispatchable work",
		"no other dispatchable or gate-processing work is available",
		"unresolved worker completion is the next useful idle action",
	} {
		if !strings.Contains(scheduling, required) {
			t.Errorf("wait_agent scheduling paragraph missing %q", required)
		}
	}

	lower := strings.ToLower(squashWhitespace(region))
	for _, forbidden := range []string{
		"must call `wait_agent` after every dispatch",
		"must call wait_agent after every dispatch",
		"always call `wait_agent`",
		"always call wait_agent",
		"foreground-wait after every dispatch",
		"wait after every dispatch",
	} {
		if strings.Contains(lower, forbidden) {
			t.Errorf("Awaiting Completion contains blanket foreground-wait wording %q", forbidden)
		}
	}
}

func TestCodexIdleNotificationRecipeShape(t *testing.T) {
	path := filepath.Join(repoRoot(t), "docs", "dev", "codex-idle-notification-probe.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	content := string(data)
	contentFlat := squashWhitespace(content)

	for _, heading := range []string{
		"## Foreground wait comparison",
		"## No-wait idle probe",
		"## Queued-notification flush check",
	} {
		if sectionAfter(content, heading) == "" {
			t.Errorf("recipe missing section %q", heading)
		}
	}
	if !strings.Contains(content, codexIdleNotificationProbePrompt) {
		t.Errorf("recipe missing exact no-write worker prompt %q", codexIdleNotificationProbePrompt)
	}
	if idleWindowSeconds(content) < 90 {
		t.Errorf("recipe idle window is less than 90 seconds")
	}
	for _, required := range []string{
		"avoid captain messages",
		"shell-outs",
		"terminal jobs",
		"tool calls",
		"retry the same handle",
		"A queued notification flushed by later activity is `queued_flush`, not `autonomous_idle_wake`.",
	} {
		if !strings.Contains(contentFlat, required) {
			t.Errorf("recipe missing %q", required)
		}
	}
	for classification := range codexIdleNotificationClassifications {
		if !strings.Contains(content, "`"+classification+"`") {
			t.Errorf("recipe missing interpretation classification %q", classification)
		}
	}
}

func TestCodexIdleNotificationEvidenceSchema(t *testing.T) {
	dir := filepath.Join(repoRoot(t), "docs", "dev", "_evidence", "codex-idle-notification-probe")
	paths, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		t.Fatalf("glob evidence: %v", err)
	}
	if len(paths) == 0 {
		t.Fatalf("no Codex idle-notification evidence JSON files in %s", dir)
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

		if got := requireInt(t, path, record, "schema_version"); got != 1 {
			t.Errorf("%s: schema_version = %d, want 1", path, got)
		}
		runID := requireString(t, path, record, "run_id")
		requireStringEqual(t, path, record, "host", "codex")
		requireString(t, path, record, "codex_cli_version")
		requireString(t, path, record, "worker_handle")
		requireString(t, path, record, "worker_prompt")
		requireInt(t, path, record, "worker_delay_seconds")
		requireInt(t, path, record, "idle_window_seconds")
		requireOptionalBool(t, path, record, "notification_delivered_before_user_message")
		autonomousWake := requireBool(t, path, record, "autonomous_wake_observed")
		classification := requireString(t, path, record, "classification")
		if !codexIdleNotificationClassifications[classification] {
			t.Errorf("%s: classification %q is outside the allowed enum", path, classification)
		}
		requireString(t, path, record, "interpretation")

		for _, field := range []string{
			"session_started_at_utc",
			"spawned_at_utc",
			"fo_turn_ended_at_utc",
			"first_user_activity_at_utc",
			"final_status_notification_at_utc",
		} {
			requireOptionalRFC3339(t, path, record, field)
		}

		if runID == "2026-06-03-dogfood" && (!autonomousWake || record["first_user_activity_at_utc"] != nil) && classification == "autonomous_idle_wake" {
			t.Errorf("%s: dogfood queued-delivery evidence must not be classified as autonomous idle wake-up", path)
		}
	}
}

func paragraphContaining(text, needle string) string {
	for _, paragraph := range strings.Split(text, "\n\n") {
		if strings.Contains(paragraph, needle) {
			return paragraph
		}
	}
	return ""
}

func idleWindowSeconds(text string) int {
	re := regexp.MustCompile(`(?i)(?:minimum idle window|idle window)[^0-9\n]*(\d+)\s+seconds`)
	best := 0
	for _, match := range re.FindAllStringSubmatch(text, -1) {
		if len(match) != 2 {
			continue
		}
		var n int
		for _, r := range match[1] {
			n = n*10 + int(r-'0')
		}
		if n > best {
			best = n
		}
	}
	return best
}

func squashWhitespace(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

func requireString(t *testing.T, path string, record map[string]any, field string) string {
	t.Helper()
	value, ok := record[field]
	if !ok {
		t.Errorf("%s: missing required field %q", path, field)
		return ""
	}
	s, ok := value.(string)
	if !ok || s == "" {
		t.Errorf("%s: field %q must be a non-empty string", path, field)
		return ""
	}
	return s
}

func requireStringEqual(t *testing.T, path string, record map[string]any, field, want string) {
	t.Helper()
	if got := requireString(t, path, record, field); got != "" && got != want {
		t.Errorf("%s: field %q = %q, want %q", path, field, got, want)
	}
}

func requireInt(t *testing.T, path string, record map[string]any, field string) int {
	t.Helper()
	value, ok := record[field]
	if !ok {
		t.Errorf("%s: missing required field %q", path, field)
		return 0
	}
	f, ok := value.(float64)
	if !ok || f != float64(int(f)) {
		t.Errorf("%s: field %q must be an integer", path, field)
		return 0
	}
	return int(f)
}

func requireBool(t *testing.T, path string, record map[string]any, field string) bool {
	t.Helper()
	value, ok := record[field]
	if !ok {
		t.Errorf("%s: missing required field %q", path, field)
		return false
	}
	b, ok := value.(bool)
	if !ok {
		t.Errorf("%s: field %q must be a boolean", path, field)
		return false
	}
	return b
}

func requireOptionalBool(t *testing.T, path string, record map[string]any, field string) {
	t.Helper()
	value, ok := record[field]
	if !ok {
		t.Errorf("%s: missing required field %q", path, field)
		return
	}
	if value == nil {
		return
	}
	if _, ok := value.(bool); !ok {
		t.Errorf("%s: field %q must be a boolean or null", path, field)
	}
}

func requireOptionalRFC3339(t *testing.T, path string, record map[string]any, field string) {
	t.Helper()
	value, ok := record[field]
	if !ok {
		t.Errorf("%s: missing required field %q", path, field)
		return
	}
	if value == nil {
		return
	}
	s, ok := value.(string)
	if !ok || s == "" {
		t.Errorf("%s: field %q must be an RFC3339 string or null", path, field)
		return
	}
	if _, err := time.Parse(time.RFC3339, s); err != nil {
		t.Errorf("%s: field %q is not RFC3339: %v", path, field, err)
	}
}
