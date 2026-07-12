// ABOUTME: Conditional Zellij metadata forwarding through Safehouse's named env allowlist.
package cli

import (
	"bytes"
	"context"
	"testing"
)

func TestSafehouseEnvPassFlagsConditionallyAddsZellijTargetingNames(t *testing.T) {
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)

	t.Run("ZELLIJ=0 forwards the complete targeting trio", func(t *testing.T) {
		parent := []string{
			"ZELLIJ=0",
			"ZELLIJ_PANE_ID=51",
			"ZELLIJ_SESSION_NAME=excellent-pheasant",
		}
		want := []string{"--env-pass", "SPACEDOCK_BIN,ZELLIJ,ZELLIJ_PANE_ID,ZELLIJ_SESSION_NAME"}
		if got := safehouseEnvPassFlags(parent); !equalArgv(got, want) {
			t.Fatalf("safehouse env-pass flags = %v, want %v", got, want)
		}
	})

	t.Run("non-Zellij parent leaves argv at the established launcher-bin allowlist", func(t *testing.T) {
		parent := []string{
			"ZELLIJ_PANE_ID=51",
			"ZELLIJ_SESSION_NAME=excellent-pheasant",
		}
		want := []string{"--env-pass", spacedockBinEnv}
		if got := safehouseEnvPassFlags(parent); !equalArgv(got, want) {
			t.Fatalf("safehouse env-pass flags = %v, want %v", got, want)
		}
	})
}

func TestWrappedHostsUseConditionalZellijEnvPassFlags(t *testing.T) {
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)

	cases := []struct {
		name   string
		parent []string
		want   []string
	}{
		{
			name: "Zellij parent",
			parent: []string{
				"ZELLIJ=0",
				"ZELLIJ_PANE_ID=51",
				"ZELLIJ_SESSION_NAME=excellent-pheasant",
				"SAFEHOUSE_ENV_PASS=EXTRA_TARGET",
			},
			want: []string{"--env-pass", "SPACEDOCK_BIN,ZELLIJ,ZELLIJ_PANE_ID,ZELLIJ_SESSION_NAME"},
		},
		{
			name: "non-Zellij parent",
			parent: []string{
				"ZELLIJ_PANE_ID=51",
				"ZELLIJ_SESSION_NAME=excellent-pheasant",
				"SAFEHOUSE_ENV_PASS=EXTRA_TARGET",
			},
			want: []string{"--env-pass", spacedockBinEnv},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name+"/claude", func(t *testing.T) {
			dir := safehouseFixtureDir(t)
			fake := &fakeHost{manifest: compatibleManifest(t)}
			var stdout, stderr bytes.Buffer

			code := runClaudeWithEnv(context.Background(), nil, dir, tc.parent, fake, lookFound, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit=%d stderr=%q", code, stderr.String())
			}
			if got := wrappedSafehouseExtra(t, fake.launchedArg); !equalArgv(got, tc.want) {
				t.Fatalf("safehouse extra = %v, want %v", got, tc.want)
			}
			if got, ok := envValue(fake.launchedEnv, "SAFEHOUSE_ENV_PASS"); !ok || got != "EXTRA_TARGET" {
				t.Fatalf("SAFEHOUSE_ENV_PASS in launch env = %q, %v; want %q, true", got, ok, "EXTRA_TARGET")
			}
		})

		t.Run(tc.name+"/codex", func(t *testing.T) {
			dir := safehouseFixtureDir(t)
			fake := &fakeHost{manifest: compatibleManifest(t)}
			var stdout, stderr bytes.Buffer

			code := runCodexWithEnv(context.Background(), nil, dir, tc.parent, fake, lookFound, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit=%d stderr=%q", code, stderr.String())
			}
			if got := wrappedSafehouseExtra(t, fake.launchedArg); !equalArgv(got, tc.want) {
				t.Fatalf("safehouse extra = %v, want %v", got, tc.want)
			}
			if got, ok := envValue(fake.launchedEnv, "SAFEHOUSE_ENV_PASS"); !ok || got != "EXTRA_TARGET" {
				t.Fatalf("SAFEHOUSE_ENV_PASS in launch env = %q, %v; want %q, true", got, ok, "EXTRA_TARGET")
			}
		})

		t.Run(tc.name+"/pi", func(t *testing.T) {
			repo := t.TempDir()
			writePiSkillFixtures(t, repo)
			pkg := t.TempDir()
			writePiSubagentsFixtures(t, pkg)
			ops := piSafehouseReadyOps(repo, pkg)
			var stdout, stderr bytes.Buffer
			parent := append(piTestEnv(pkg, t.TempDir()), tc.parent...)

			code := runPi(context.Background(), []string{"--plugin-dir", repo}, safehouseFixtureDir(t), parent, ops, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit=%d stderr=%q", code, stderr.String())
			}
			if got := wrappedSafehouseExtra(t, ops.launched); !equalArgv(got, tc.want) {
				t.Fatalf("safehouse extra = %v, want %v", got, tc.want)
			}
			if got, ok := envValue(ops.launchedEnv, "SAFEHOUSE_ENV_PASS"); !ok || got != "EXTRA_TARGET" {
				t.Fatalf("SAFEHOUSE_ENV_PASS in launch env = %q, %v; want %q, true", got, ok, "EXTRA_TARGET")
			}
		})
	}
}

func wrappedSafehouseExtra(t *testing.T, argv []string) []string {
	t.Helper()
	if len(argv) < 2 || argv[0] != "safehouse" || argv[1] != "--trust-workdir-config" {
		t.Fatalf("wrapped argv = %v, want safehouse prefix", argv)
	}
	for i := 2; i < len(argv); i++ {
		if argv[i] == "--" {
			return argv[2:i]
		}
	}
	t.Fatalf("wrapped argv lacks -- separator: %v", argv)
	return nil
}
