// ABOUTME: Adjacent Spacedock plugin discovery and front-door selection tests.
// ABOUTME: Covers host-manifest qualification, precedence, gating, and Codex preflight safety.
package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type resolvingHost struct {
	fakeHost
	resolveCalls     int
	resolvedManifest string
}

func (h *resolvingHost) ResolveManifest(host string) (string, error) {
	h.resolveCalls++
	manifest, err := h.fakeHost.ResolveManifest(host)
	h.resolvedManifest = manifest
	return manifest, err
}

func localPluginCheckout(t *testing.T, host string) (root, executable string) {
	t.Helper()
	root = t.TempDir()
	executable = filepath.Join(root, "spacedock")
	if err := os.WriteFile(executable, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeLocalPluginManifest(t, root, host, "spacedock", "./skills/")
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	return resolvedRoot, executable
}

func writeLocalPluginManifest(t *testing.T, root, host, name, skills string) string {
	t.Helper()
	manifestDir := "." + host + "-plugin"
	manifestPath := filepath.Join(root, manifestDir, "plugin.json")
	if err := os.MkdirAll(filepath.Dir(manifestPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if skills != "" {
		skillsPath := skills
		if !filepath.IsAbs(skillsPath) {
			skillsPath = filepath.Join(root, filepath.FromSlash(skillsPath))
		}
		if err := os.MkdirAll(skillsPath, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	body := fmt.Sprintf(`{"name":%q,"version":%q,"skills":%q}`+"\n", name, displayVersion(), skills)
	if err := os.WriteFile(manifestPath, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return manifestPath
}

func TestAdjacentPluginAutomaticallySelectedForClaude(t *testing.T) {
	root, executable := localPluginCheckout(t, "claude")
	withExecutablePath(t, executable, nil)
	host := &resolvingHost{fakeHost: fakeHost{resolveErr: fmt.Errorf("installed resolver must not run")}}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), host, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if host.resolveCalls != 0 {
		t.Fatalf("installed resolver called %d times, want 0 for adjacent checkout", host.resolveCalls)
	}
	want := []string{"claude", "--agent", "spacedock:first-officer", "--permission-mode", "auto", "--plugin-dir", root, wantBootstrapPrompt}
	if !equalArgv(host.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", host.launchedArg, want)
	}
}

func TestAdjacentPluginAutomaticallySelectedForCodex(t *testing.T) {
	t.Setenv("CODEX_HOME", t.TempDir())
	root, executable := localPluginCheckout(t, "codex")
	withExecutablePath(t, executable, nil)
	host := &codexPluginDirHost{fakeHost: fakeHost{manifest: filepath.Join(root, ".codex-plugin", "plugin.json")}}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), nil, t.TempDir(), host, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if host.inspectErr != nil {
		t.Fatalf("marketplace inspection failed: %v", host.inspectErr)
	}
	wantManifest, err := os.ReadFile(filepath.Join(root, ".codex-plugin", "plugin.json"))
	if err != nil {
		t.Fatalf("read adjacent root manifest: %v", err)
	}
	if !bytes.Equal(host.installedManifest, wantManifest) {
		t.Fatalf("staged plugin manifest = %q, want the adjacent root's manifest %q", host.installedManifest, wantManifest)
	}
	if host.launchedArg == nil {
		t.Fatal("Codex launch seam was not reached")
	}
}

func TestInvalidAdjacentPluginFallsBackToInstalledProvider(t *testing.T) {
	for _, hostName := range []string{"claude", "codex"} {
		for _, tc := range []struct {
			name  string
			setup func(t *testing.T, root string)
		}{
			{name: "missing manifest"},
			{name: "wrong manifest name", setup: func(t *testing.T, root string) {
				writeLocalPluginManifest(t, root, hostName, "not-spacedock", "./skills/")
			}},
			{name: "missing skills directory", setup: func(t *testing.T, root string) {
				writeLocalPluginManifest(t, root, hostName, "spacedock", "./skills/")
				if err := os.RemoveAll(filepath.Join(root, "skills")); err != nil {
					t.Fatal(err)
				}
			}},
			{name: "release-style bin directory", setup: func(t *testing.T, root string) {
				writeLocalPluginManifest(t, filepath.Dir(root), hostName, "spacedock", "./skills/")
			}},
		} {
			t.Run(hostName+"/"+tc.name, func(t *testing.T) {
				t.Setenv("CODEX_HOME", t.TempDir())
				root := t.TempDir()
				if tc.name == "release-style bin directory" {
					root = filepath.Join(root, "bin")
					if err := os.MkdirAll(root, 0o755); err != nil {
						t.Fatal(err)
					}
				}
				if tc.setup != nil {
					tc.setup(t, root)
				}
				executable := filepath.Join(root, "spacedock")
				if err := os.WriteFile(executable, []byte("#!/bin/sh\n"), 0o755); err != nil {
					t.Fatal(err)
				}
				withExecutablePath(t, executable, nil)
				h := &resolvingHost{fakeHost: fakeHost{manifest: compatibleManifest(t)}}
				var stdout, stderr bytes.Buffer
				var code int
				if hostName == "claude" {
					code = runClaude(context.Background(), nil, t.TempDir(), h, lookFound, &stdout, &stderr)
					if strings.Contains(strings.Join(h.launchedArg, " "), "--plugin-dir") {
						t.Fatalf("invalid adjacent checkout reached Claude argv: %v", h.launchedArg)
					}
				} else {
					code = runCodex(context.Background(), nil, t.TempDir(), h, lookFound, &stdout, &stderr)
					if len(h.installCmds) != 0 {
						t.Fatalf("invalid adjacent checkout invoked Codex install: %v", h.installCmds)
					}
				}
				if code != 0 || h.resolveCalls == 0 {
					t.Fatalf("installed fallback: exit=%d resolveCalls=%d stderr=%q", code, h.resolveCalls, stderr.String())
				}
			})
		}
	}
}

func TestMissingAdjacentPluginPreservesAutoInstall(t *testing.T) {
	for _, hostName := range []string{"claude", "codex"} {
		t.Run(hostName, func(t *testing.T) {
			t.Setenv("CODEX_HOME", t.TempDir())
			executable := filepath.Join(t.TempDir(), "spacedock")
			if err := os.WriteFile(executable, []byte("#!/bin/sh\n"), 0o755); err != nil {
				t.Fatal(err)
			}
			withExecutablePath(t, executable, nil)
			h := &resolvingHost{fakeHost: fakeHost{manifestAfterInstall: compatibleManifest(t)}}
			var stdout, stderr bytes.Buffer
			var code int
			if hostName == "claude" {
				code = runClaude(context.Background(), nil, t.TempDir(), h, lookFound, &stdout, &stderr)
			} else {
				code = runCodex(context.Background(), nil, t.TempDir(), h, lookFound, &stdout, &stderr)
			}
			if code != 0 {
				t.Fatalf("exit = %d, want installed fallback to heal (stderr=%q)", code, stderr.String())
			}
			if len(h.installCmds) != 3 || h.installCmds[0] != hostName {
				t.Fatalf("install seam = %v, want one %s auto-install", h.installCmds, hostName)
			}
			if h.resolveCalls != 2 || h.launchedArg == nil {
				t.Fatalf("fallback did not re-gate then launch: resolveCalls=%d argv=%v", h.resolveCalls, h.launchedArg)
			}
		})
	}
}

func TestExplicitPluginDirPrecedesAdjacentPlugin(t *testing.T) {
	adjacent, executable := localPluginCheckout(t, "claude")
	explicit, _ := localPluginCheckout(t, "claude")
	withExecutablePath(t, executable, nil)
	host := &resolveErrHost{}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--plugin-dir", explicit}, t.TempDir(), host, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	joined := strings.Join(host.launchedArg, " ")
	if !strings.Contains(joined, "--plugin-dir "+explicit) || strings.Contains(joined, "--plugin-dir "+adjacent) {
		t.Fatalf("explicit override did not win: %v", host.launchedArg)
	}
}

func TestExplicitCodexPluginDirPrecedesAdjacentPlugin(t *testing.T) {
	t.Setenv("CODEX_HOME", t.TempDir())
	adjacent, executable := localPluginCheckout(t, "codex")
	explicit, _ := localPluginCheckout(t, "codex")
	withExecutablePath(t, executable, nil)
	// localPluginCheckout produces byte-identical manifests for both roots, so a
	// staged-manifest comparison can't tell which checkout won without a marker.
	// The marker is an extra JSON field alongside the real (gate-compatible)
	// version — a .json write/read stays outside the instruction-file boundary
	// guard's tracked skill/agent surface.
	mustWrite(t, filepath.Join(adjacent, ".codex-plugin", "plugin.json"), `{"name":"spacedock","version":"`+displayVersion()+`","skills":"./skills/","marker":"adjacent"}`)
	mustWrite(t, filepath.Join(explicit, ".codex-plugin", "plugin.json"), `{"name":"spacedock","version":"`+displayVersion()+`","skills":"./skills/","marker":"explicit"}`)
	host := &codexPluginDirHost{fakeHost: fakeHost{manifest: filepath.Join(explicit, ".codex-plugin", "plugin.json")}}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--plugin-dir", explicit}, t.TempDir(), host, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if !strings.Contains(string(host.installedManifest), `"marker":"explicit"`) {
		t.Fatalf("staged manifest = %q, want the explicit checkout's marker (explicit must win over adjacent %q)", host.installedManifest, adjacent)
	}
}

func TestClaudePostFencePluginDirDoesNotBypassGate(t *testing.T) {
	additional := t.TempDir()
	writeLocalPluginManifest(t, additional, "claude", "additional-tools", "./skills/")
	installedManifest := compatibleManifest(t)
	host := &resolvingHost{fakeHost: fakeHost{manifest: installedManifest}}
	withExecutablePath(t, filepath.Join(t.TempDir(), "missing-spacedock"), os.ErrNotExist)
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--", "--plugin-dir", additional}, t.TempDir(), host, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if host.resolveCalls == 0 {
		t.Fatal("post-fence additional plugin bypassed the installed Spacedock gate")
	}
	if host.resolvedManifest != installedManifest {
		t.Fatalf("selected Spacedock manifest = %q, want installed provider %q", host.resolvedManifest, installedManifest)
	}
	want := []string{"claude", "--agent", "spacedock:first-officer", "--permission-mode", "auto", "--plugin-dir", additional, wantBootstrapPrompt}
	if !equalArgv(host.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", host.launchedArg, want)
	}
}

func TestCodexInvalidExplicitPluginFailsBeforeMutation(t *testing.T) {
	codexHome := t.TempDir()
	t.Setenv("CODEX_HOME", codexHome)
	checkout := t.TempDir()
	writeLocalPluginManifest(t, checkout, "codex", "not-spacedock", "./skills/")
	host := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--plugin-dir", checkout}, t.TempDir(), host, lookFound, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("exit = 0, want invalid-checkout failure (stderr=%q)", stderr.String())
	}
	if !strings.Contains(stderr.String(), ".codex-plugin/plugin.json") || !strings.Contains(stderr.String(), "name") {
		t.Fatalf("diagnostic is not actionable: %q", stderr.String())
	}
	if len(host.installCmds) != 0 {
		t.Fatalf("install seam called before validation: %v", host.installCmds)
	}
	entries, err := os.ReadDir(codexHome)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("Codex config mutated before validation: %v", entries)
	}
}
