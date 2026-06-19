// ABOUTME: SPACEDOCK_FO_MODEL launcher-export tests — the var reaches a launched FO
// ABOUTME: (normalized from --model), is omitted when unresolvable, and rides --env-pass.
package cli

import (
	"bytes"
	"context"
	"testing"
)

// TestClaudeFrontDoorExportsFOModel: the launcher parses --model off the host
// passthrough, normalizes it, and exports SPACEDOCK_FO_MODEL into the launch env
// so «fo.tier»() can read the FO's own model at boot. The normalized value (not
// the raw flag) is what reaches the FO — a `--model sonnet[1m]` lands as `sonnet`.
// An unresolvable or absent --model leaves the var unset, which the FO resolves
// fail-safe to level-2-only.
func TestClaudeFrontDoorExportsFOModel(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		wantValue string
		wantSet   bool
	}{
		{"space form haiku", []string{"--", "--model", "haiku"}, "haiku", true},
		{"equals form sonnet", []string{"--", "--model=sonnet"}, "sonnet", true},
		{"extended suffix normalizes", []string{"--", "--model", "sonnet[1m]"}, "sonnet", true},
		{"full id normalizes to family", []string{"--", "--model", "claude-opus-4-8[1m]"}, "opus", true},
		{"default alias resolves to opus", []string{"--", "--model", "default"}, "opus", true},
		{"absent --model → var unset (fail-safe)", []string{"--", "-p", "x"}, "", false},
		{"unresolvable --model → var unset (fail-safe)", []string{"--", "--model", "gpt-4-turbo"}, "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHost{manifest: compatibleManifest(t)}
			var stdout, stderr bytes.Buffer

			code := runClaude(context.Background(), tc.args, t.TempDir(), fake, lookFound, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
			}
			got, ok := envValue(fake.launchedEnv, spacedockFOModelEnv)
			if ok != tc.wantSet || got != tc.wantValue {
				t.Fatalf("%s in launch env = %q, %v; want %q, %v (env=%v)",
					spacedockFOModelEnv, got, ok, tc.wantValue, tc.wantSet, fake.launchedEnv)
			}
		})
	}
}

// TestClaudeSafehouseForwardsFOModel: under a .safehouse wrap a resolved
// SPACEDOCK_FO_MODEL rides `--env-pass SPACEDOCK_FO_MODEL` among the safehouse
// flags (before `--`), so the var survives safehouse's env sanitization and
// reaches the sandboxed FO. The launcher binary env-pass is already there; this is
// the second forwarded var. The inner program stays `claude`.
func TestClaudeSafehouseForwardsFOModel(t *testing.T) {
	dir := safehouseFixtureDir(t)
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--", "--model", "haiku"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	// The FO_MODEL env-pass must appear before the `--` separator (a safehouse flag).
	if !envPassBeforeDash(fake.launchedArg, spacedockFOModelEnv) {
		t.Fatalf("wrapped argv missing `--env-pass %s` before `--`: %v", spacedockFOModelEnv, fake.launchedArg)
	}
	// And the value must be in the forwarded env.
	if got, ok := envValue(fake.launchedEnv, spacedockFOModelEnv); !ok || got != "haiku" {
		t.Fatalf("%s in launch env = %q, %v; want haiku, true", spacedockFOModelEnv, got, ok)
	}
}

// TestClaudeSafehouseOmitsFOModelEnvPassWhenUnresolved: with no resolvable
// --model, the wrap carries NO `--env-pass SPACEDOCK_FO_MODEL` — never a stale
// pass-through of an unset var — mirroring the SPACEDOCK_BIN omit-on-failure.
func TestClaudeSafehouseOmitsFOModelEnvPassWhenUnresolved(t *testing.T) {
	dir := safehouseFixtureDir(t)
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--", "-p", "x"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if envPassBeforeDash(fake.launchedArg, spacedockFOModelEnv) {
		t.Fatalf("wrapped argv carried `--env-pass %s` despite no resolvable model: %v", spacedockFOModelEnv, fake.launchedArg)
	}
}

// TestCodexFrontDoorExportsFOModel: the codex front door also parses --model and
// exports SPACEDOCK_FO_MODEL — the tier mechanism is host-neutral on the launcher
// side even though the live routing slice is Claude-only.
func TestCodexFrontDoorExportsFOModel(t *testing.T) {
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--", "--model", "haiku"}, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if got, ok := envValue(fake.launchedEnv, spacedockFOModelEnv); !ok || got != "haiku" {
		t.Fatalf("%s in codex launch env = %q, %v; want haiku, true", spacedockFOModelEnv, got, ok)
	}
}

// envPassBeforeDash reports whether the argv carries a `--env-pass {key}` pair
// before the safehouse `--` separator (so it is a safehouse flag, not an inner
// host arg).
func envPassBeforeDash(argv []string, key string) bool {
	for i := 0; i+1 < len(argv); i++ {
		if argv[i] == "--" {
			return false
		}
		if argv[i] == "--env-pass" && argv[i+1] == key {
			return true
		}
	}
	return false
}
