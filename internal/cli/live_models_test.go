// ABOUTME: `spacedock live-models` fixes the exact live lane ids it prints.
package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// TestLiveModelsCommandPrintsPinnedLaneModels fixes the exact ids `spacedock
// live-models` prints. The expected ids are authored here — not read from
// internal/release — so this fails if a lane constant changes without the print
// moving in lockstep (e.g. reverting the Pi lane to the old gpt-5.6-luna), or if
// a lane/auth key is added, dropped, or reordered.
func TestLiveModelsCommandPrintsPinnedLaneModels(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"live-models"}, nil, "", nil, &stdout, &stderr, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("live-models exit=%d stderr=%q", code, stderr.String())
	}
	const want = "claude.sonnet=claude-sonnet-5\n" +
		"claude.opus=claude-opus-4-8\n" +
		"codex.exec=gpt-6-luna\n" +
		"pi.oauth=openai-codex/gpt-6-luna:max\n" +
		"pi.api-key=openai/gpt-6-luna:max\n"
	if got := stdout.String(); got != want {
		t.Fatalf("live-models stdout = %q, want %q", got, want)
	}
}

// TestLiveModelsCommandGetPrintsOneID proves `--get` resolves one lane key to its
// id and that an unknown key exits 2 with a diagnostic rather than printing an
// empty model. It fails if a named key stops resolving or if the unknown-key
// path silently succeeds.
func TestLiveModelsCommandGetPrintsOneID(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run(context.Background(), []string{"live-models", "--get", "pi.oauth"}, nil, "", nil, &stdout, &stderr, &status.NativeRunner{}, nil); code != 0 {
		t.Fatalf("--get exit=%d stderr=%q", code, stderr.String())
	}
	if got, want := stdout.String(), "openai-codex/gpt-6-luna:max\n"; got != want {
		t.Fatalf("--get stdout = %q, want %q", got, want)
	}

	stdout.Reset()
	stderr.Reset()
	if code := run(context.Background(), []string{"live-models", "--get", "pi.nope"}, nil, "", nil, &stdout, &stderr, &status.NativeRunner{}, nil); code != 2 {
		t.Fatalf("unknown-key exit=%d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "unknown key") {
		t.Fatalf("unknown-key stderr = %q, want an unknown-key diagnostic", stderr.String())
	}
}
