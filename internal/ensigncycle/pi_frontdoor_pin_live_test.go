//go:build live

package ensigncycle

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// piStampedSpacedockBinary builds the checkout binary as a release-shaped
// artifact: the goreleaser stamp form of the checkout manifest's own version
// (0.28.0-pre0, an existing tag), so the pinned install source resolves to a
// real remote tag in the live run.
func piStampedSpacedockBinary(t *testing.T, repo string) string {
	t.Helper()
	out := filepath.Join(t.TempDir(), "spacedock")
	cmd := exec.Command("go", "build", "-o", out,
		"-ldflags", "-X github.com/spacedock-dev/spacedock/internal/cli.Version=0.28.0-pre0",
		"./cmd/spacedock")
	cmd.Dir = repo
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build stamped spacedock for the pinned-package live journey: %v\n%s", err, b)
	}
	return out
}

//spacedock:live-proof id=pi-front-door-pinned-package lane=pi-live
func TestLivePiFrontDoorInstallsPinnedPackage(t *testing.T) {
	repo := repoRoot(t)
	piSubagentsRoot := piSubagentsPackageRoot(t)
	binary := piStampedSpacedockBinary(t, repo)

	// The incident fixture: a fresh pi-home whose settings.json registers NO
	// Spacedock package (the post-skew state the v0.27.2 abort left behind).
	piHome := t.TempDir()
	sessionDir := t.TempDir()
	cleanHome := t.TempDir()
	decision := seedPiLiveAuth(t, piHome, os.Getenv("HOME"), os.Getenv("CODEX_AUTH_JSON"), os.Getenv("OPENAI_API_KEY"), os.Getenv("SPACEDOCK_PI_LIVE_REQUIRED"))
	artifactDir := filepath.Join(piLiveArtifactDir(t, "pi-frontdoor-pin"), "run")
	if err := os.MkdirAll(filepath.Join(artifactDir, "sessions"), 0o755); err != nil {
		t.Fatal(err)
	}
	env := piLiveEnvForAuth(piHome, sessionDir, cleanHome, filepath.Dir(binary), piSubagentsRoot, os.Getenv("OPENAI_API_KEY"), decision.mode)
	model := piLiveChildModel(decision)

	// The ordinary installed front door: no --plugin-dir, no SPACEDOCK_REPO_ROOT.
	runPiLiveCommand(t, artifactDir, repo, env, binary,
		"pi",
		"Reply with the single word READY.",
		"--",
		"--print",
		"--model", model,
		"--session-dir", filepath.Join(artifactDir, "sessions"),
	)

	// AC-1 durable proof: the front door's one repair registered the binary's
	// own pinned ref.
	entry := piSpacedockSettingsEntry(t, piHome)
	want := "git:github.com/spacedock-dev/spacedock@v0.28.0-pre0"
	if !strings.Contains(entry, want) {
		t.Fatalf("pinned package entry %q not registered after the repair; got entry %q", want, entry)
	}
}

// piSpacedockSettingsEntry returns the spacedock package entry (string or
// {"source": ...} form) from the pi agent dir's settings.json.
func piSpacedockSettingsEntry(t *testing.T, piHome string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(piHome, "settings.json"))
	if err != nil {
		t.Fatalf("read pi settings: %v", err)
	}
	var settings struct {
		Packages []json.RawMessage `json:"packages"`
	}
	if err := json.Unmarshal(data, &settings); err != nil {
		t.Fatalf("parse pi settings: %v", err)
	}
	for _, raw := range settings.Packages {
		var s string
		if json.Unmarshal(raw, &s) == nil && strings.Contains(s, "spacedock") {
			return s
		}
		var obj struct {
			Source string `json:"source"`
		}
		if json.Unmarshal(raw, &obj) == nil && strings.Contains(obj.Source, "spacedock") {
			return obj.Source
		}
	}
	return ""
}
