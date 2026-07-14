package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestBridgeIngressWakeHiddenCLIPrintsJSONNoop(t *testing.T) {
	root := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := run(context.Background(),
		[]string{"bridge", "ingress", "wake", "--host", "codex", "--repo-root", root, "--members", "a,b"},
		nil, filepath.Join(root, "elsewhere"), strings.NewReader(""), &stdout, &stderr, &fakeRunner{}, nil)

	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	var got struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &got); err != nil {
		t.Fatalf("stdout JSON: %v\n%s", err, stdout.String())
	}
	if got.Status != "noop" {
		t.Fatalf("status = %q, want noop", got.Status)
	}
}
