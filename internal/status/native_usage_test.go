// ABOUTME: Usage/parse-error parity — bad flags exit 1 (never 2) with the same
// ABOUTME: Error: ... stderr the oracle emits, locking the {0,1} exit domain.
package status

import (
	"path/filepath"
	"testing"
)

// usageCases are argv shapes that must exit 1 with a stderr Error: message in
// both native and oracle — never exit 2.
var usageCases = []struct {
	name string
	args func(root string) []string
}{
	{"bad-where-no-operator", func(r string) []string { return []string{"--workflow-dir", r, "--where", "status"} }},
	{"where-missing-arg", func(r string) []string { return []string{"--workflow-dir", r, "--where"} }},
	{"fields-and-all-fields", func(r string) []string {
		return []string{"--workflow-dir", r, "--fields", "a", "--all-fields"}
	}},
	{"boot-with-next", func(r string) []string { return []string{"--workflow-dir", r, "--boot", "--next"} }},
	{"next-id-with-set", func(r string) []string {
		return []string{"--workflow-dir", r, "--next-id", "--set", "x", "y=z"}
	}},
	{"resolve-missing-arg", func(r string) []string { return []string{"--workflow-dir", r, "--resolve"} }},
	{"workflow-dir-missing-arg", func(r string) []string { return []string{"--workflow-dir"} }},
	{"id-material-without-next-id", func(r string) []string {
		return []string{"--workflow-dir", r, "--id-seed", "x"}
	}},
	{"root-without-discover-or-resolve", func(r string) []string {
		return []string{"--workflow-dir", r, "--root", r}
	}},
	{"page-zero", func(r string) []string { return []string{"--workflow-dir", r, "--page", "0"} }},
	{"page-negative", func(r string) []string { return []string{"--workflow-dir", r, "--page", "-1"} }},
	{"page-non-integer", func(r string) []string { return []string{"--workflow-dir", r, "--page", "abc"} }},
	{"limit-negative", func(r string) []string { return []string{"--workflow-dir", r, "--limit", "-1"} }},
	{"limit-non-integer", func(r string) []string { return []string{"--workflow-dir", r, "--limit", "abc"} }},
	{"page-with-limit-zero", func(r string) []string {
		return []string{"--workflow-dir", r, "--page", "2", "--limit", "0"}
	}},
	{"page-with-next", func(r string) []string { return []string{"--workflow-dir", r, "--next", "--page", "2"} }},
	{"page-with-boot", func(r string) []string { return []string{"--workflow-dir", r, "--boot", "--page", "2"} }},
	{"limit-with-validate", func(r string) []string {
		return []string{"--workflow-dir", r, "--validate", "--limit", "5"}
	}},
	{"limit-with-read", func(r string) []string {
		return []string{"--workflow-dir", r, "--read", "003-wire-cli", "--limit", "5"}
	}},
	{"page-with-resolve", func(r string) []string {
		return []string{"--workflow-dir", r, "--resolve", "003-wire-cli", "--page", "2"}
	}},
	{"page-with-short-id", func(r string) []string {
		return []string{"--workflow-dir", r, "--short-id", "003-wire-cli", "--page", "2"}
	}},
	{"page-with-next-id", func(r string) []string {
		return []string{"--workflow-dir", r, "--next-id", "--page", "2"}
	}},
	{"limit-with-set", func(r string) []string {
		return []string{"--workflow-dir", r, "--set", "003-wire-cli", "score=0.5", "--limit", "5"}
	}},
	{"page-with-archive", func(r string) []string {
		return []string{"--workflow-dir", r, "--archive", "003-wire-cli", "--page", "2"}
	}},
}

func TestNativeUsageErrorsExitOneNotTwo(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)

	for _, tc := range usageCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			args := tc.args(root)
			nOut, nErr, nCode := runNative(t, root, env, args...)

			if nCode == 2 {
				t.Fatalf("native exit=2 for usage error %q; must be 1 (never 2)", tc.name)
			}
			if nCode != 1 {
				t.Fatalf("usage error %q exit=%d, want 1", tc.name, nCode)
			}
			assertEnvelopeGolden(t, "native-usage-"+tc.name, goldenEnvelope{
				stdout: normalize(nOut, root), stderr: normalize(nErr, root), exit: nCode,
			})
			if nOut != "" {
				t.Fatalf("stdout must be empty on usage error: native=%q", nOut)
			}
		})
	}
}
