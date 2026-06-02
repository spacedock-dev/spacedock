// ABOUTME: AC-1/AC-2 golden read parity — native stdout for the five read
// ABOUTME: subcommands matches goldens captured from the oracle, after normalization.
package status

import (
	"path/filepath"
	"testing"
)

// readCases are the five FO-load-bearing read subcommands compared byte-for-byte
// (post-normalization) against goldens captured from the oracle at retirement.
var readCases = []struct {
	name   string
	golden string
	extra  []string // args after --workflow-dir <root>
}{
	{name: "default", golden: "seq-default.txt", extra: nil},
	{name: "next", golden: "seq-next.txt", extra: []string{"--next"}},
	{name: "validate", golden: "seq-validate.txt", extra: []string{"--validate"}},
	{name: "resolve", golden: "seq-resolve.txt", extra: []string{"--resolve", "003-wire-cli"}},
	{name: "short-id", golden: "seq-short-id.txt", extra: []string{"--short-id", "003-wire-cli"}},
}

func TestGoldenRead(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)

	for _, tc := range readCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"--workflow-dir", root}, tc.extra...)

			if *update {
				out, stderr, code := runNative(t, root, env, args...)
				if code != 0 {
					t.Fatalf("native exit=%d stderr=%q while capturing golden", code, stderr)
				}
				writeGolden(t, tc.golden, normalize(out, root))
				return
			}

			out, stderr, code := runNative(t, root, env, args...)
			if code != 0 {
				t.Fatalf("native exit=%d stderr=%q", code, stderr)
			}
			got := normalize(out, root)
			want := readGolden(t, tc.golden)
			if got != want {
				t.Fatalf("read parity mismatch for %s\n--- got ---\n%s\n--- want ---\n%s", tc.name, got, want)
			}
		})
	}
}
