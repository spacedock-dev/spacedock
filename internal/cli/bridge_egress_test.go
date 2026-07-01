package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBridgeEgressEmitHiddenCLISilentAndWrites(t *testing.T) {
	root := t.TempDir()
	payload := `{"event":"SessionStart","session_id":"ses-1","source":"startup"}`
	var stdout, stderr bytes.Buffer

	code := run(context.Background(),
		[]string{"bridge", "egress", "emit", "--host", "claude"},
		nil, root, strings.NewReader(payload), &stdout, &stderr, &fakeRunner{}, nil)

	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if stdout.Len() != 0 || stderr.Len() != 0 {
		t.Fatalf("bridge egress should be silent, stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	data, err := os.ReadFile(filepath.Join(root, "_bridge", "events.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	var event struct {
		Host      string `json:"host"`
		Event     string `json:"event"`
		SessionID string `json:"session_id"`
		ActorID   string `json:"actor_id"`
		Detail    struct {
			Source string `json:"source"`
		} `json:"detail"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(data), &event); err != nil {
		t.Fatalf("event JSON: %v\n%s", err, data)
	}
	if event.Host != "claude" || event.Event != "SessionStart" || event.SessionID != "ses-1" || event.ActorID != "ses-1" || event.Detail.Source != "startup" {
		t.Fatalf("event mismatch: %+v", event)
	}
}

func TestBridgeEgressEmitMalformedPayloadSilentNoop(t *testing.T) {
	root := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := run(context.Background(),
		[]string{"bridge", "egress", "emit", "--host", "claude"},
		nil, root, strings.NewReader(`{`), &stdout, &stderr, &fakeRunner{}, nil)

	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if stdout.Len() != 0 || stderr.Len() != 0 {
		t.Fatalf("bridge egress should be silent, stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	if _, err := os.Stat(filepath.Join(root, "_bridge", "events.jsonl")); !os.IsNotExist(err) {
		t.Fatalf("events.jsonl exists after malformed payload: %v", err)
	}
}

func TestBridgeCommandStaysOutOfTopLevelHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if strings.Contains(stdout.String(), "bridge") {
		t.Fatalf("top-level help exposes hidden bridge command:\n%s", stdout.String())
	}
}
