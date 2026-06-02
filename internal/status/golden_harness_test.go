// ABOUTME: Golden harness — freezes the certified-parity native output (3 channels)
// ABOUTME: to testdata/golden, normalized for the per-run root, for the graduated suite.
package status

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// goldenEnvelope is the certified native output frozen per case: the three
// channels, root-normalized via the shared normalize().
type goldenEnvelope struct {
	stdout string
	stderr string
	exit   int
}

// goldenSectionRe splits a marshalled envelope back into its named sections.
var goldenSectionRe = regexp.MustCompile(`(?m)^===== (\w+) =====$`)

// marshalEnvelope serializes a normalized envelope to the on-disk golden text.
func marshalEnvelope(env goldenEnvelope) string {
	var b strings.Builder
	fmt.Fprintf(&b, "===== exit =====\n%d\n", env.exit)
	fmt.Fprintf(&b, "===== stdout =====\n%s", env.stdout)
	fmt.Fprintf(&b, "\n===== stderr =====\n%s", env.stderr)
	return b.String()
}

// parseEnvelope splits a marshalled envelope back into its channels.
func parseEnvelope(t *testing.T, name, raw string) goldenEnvelope {
	t.Helper()
	idx := goldenSectionRe.FindAllStringSubmatchIndex(raw, -1)
	if len(idx) == 0 {
		t.Fatalf("golden %s has no sections:\n%s", name, raw)
	}
	sections := map[string]string{}
	for i, m := range idx {
		key := raw[m[2]:m[3]]
		start := m[1] + 1
		end := len(raw)
		if i+1 < len(idx) {
			end = idx[i+1][0] - 1
		}
		if start > end {
			start = end
		}
		sections[key] = raw[start:end]
	}
	env := goldenEnvelope{stdout: sections["stdout"], stderr: sections["stderr"]}
	if _, err := fmt.Sscanf(strings.TrimSpace(sections["exit"]), "%d", &env.exit); err != nil {
		t.Fatalf("golden %s exit not an int: %q", name, sections["exit"])
	}
	return env
}

// envelopePath returns the testdata/golden path for a named envelope golden.
func envelopePath(name string) string {
	return filepath.Join("testdata", "golden", name+".txt")
}

// assertEnvelopeGolden compares a normalized 3-channel native run against the
// frozen golden, or captures it under -update. The caller must root-normalize the
// channels via normalize() before passing them.
func assertEnvelopeGolden(t *testing.T, name string, native goldenEnvelope) {
	t.Helper()
	if *update {
		path := envelopePath(name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir golden dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(marshalEnvelope(native)), 0o644); err != nil {
			t.Fatalf("write golden %s: %v", name, err)
		}
		return
	}
	raw, err := os.ReadFile(envelopePath(name))
	if err != nil {
		t.Fatalf("read golden %s (regenerate with -update): %v", name, err)
	}
	want := parseEnvelope(t, name, string(raw))
	if native.stdout != want.stdout {
		t.Errorf("%s: stdout mismatch\n--- native ---\n%q\n--- golden ---\n%q", name, native.stdout, want.stdout)
	}
	if native.stderr != want.stderr {
		t.Errorf("%s: stderr mismatch\n--- native ---\n%q\n--- golden ---\n%q", name, native.stderr, want.stderr)
	}
	if native.exit != want.exit {
		t.Errorf("%s: exit mismatch native=%d golden=%d", name, native.exit, want.exit)
	}
}

// assertTextGolden captures/compares a single normalized text blob (e.g. a
// mutated entity's frontmatter) against testdata/golden/<name>.txt. The caller
// must normalize the text before passing it.
func assertTextGolden(t *testing.T, name, got string) {
	t.Helper()
	path := envelopePath(name)
	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir golden dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("write golden %s: %v", name, err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s (regenerate with -update): %v", name, err)
	}
	if got != string(want) {
		t.Errorf("%s: text mismatch\n--- native ---\n%s\n--- golden ---\n%s", name, got, string(want))
	}
}
