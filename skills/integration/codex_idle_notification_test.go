// ABOUTME: Contract tests for Codex no-wait completion evidence.
// ABOUTME: Keeps queued notifications distinct from autonomous idle wake-up.
package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var codexIdleNotificationClassifications = map[string]bool{
	"async_idle_monitoring":    true,
	"queued_flush":             true,
	"autonomous_idle_wake":     true,
	"no_notification_observed": true,
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
