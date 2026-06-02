// ABOUTME: Golden harness — freezes the certified-parity native output (3 channels
// ABOUTME: + optional dispatch body) to testdata/golden, normalized for root/home.
package dispatch

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// update regenerates the checked-in dispatch goldens from the current native
// output instead of comparing against them. Run:
// go test ./internal/dispatch -run TestBuild... -update
var update = flag.Bool("update", false, "regenerate golden files from native output")

// goldenEnvelope is the certified native output frozen per case: the three
// channels plus the optional dispatch body, all root/home-normalized.
type goldenEnvelope struct {
	res  runResult
	body string // "" when the case asserts no dispatch body
}

// normPaths replaces the per-run absolute fixture roots with stable placeholders
// so a golden captured under one t.TempDir compares against a run under another.
// Both the as-spelled and the realpath'd spelling (macOS /var->/private/var) map
// to the same placeholder, mirroring the status harness's normalize().
func normPaths(s, root, home string) string {
	for _, sub := range []struct {
		dir         string
		placeholder string
	}{
		{home, "<HOME>"},
		{root, "<ROOT>"},
	} {
		if sub.dir == "" {
			continue
		}
		if real := realpath(sub.dir); real != sub.dir {
			s = strings.ReplaceAll(s, real, sub.placeholder)
		}
		s = strings.ReplaceAll(s, sub.dir, sub.placeholder)
	}
	return s
}

// realpath resolves symlinks the way the oracle's os.path.realpath did, so the
// placeholder substitution accounts for the macOS /var->/private/var rewrite.
func realpath(p string) string {
	resolved, err := filepath.EvalSymlinks(p)
	if err != nil {
		return p
	}
	return resolved
}

// normRun root/home-normalizes all three channels of a run result.
func normRun(res runResult, root, home string) runResult {
	return runResult{
		stdout: normPaths(res.stdout, root, home),
		stderr: normPaths(res.stderr, root, home),
		exit:   res.exit,
	}
}

// goldenSep delimits the named sections of a golden envelope file.
const goldenSep = "\n===== %s =====\n"

var goldenSectionRe = regexp.MustCompile(`(?m)^===== (\w+) =====$`)

// marshalGolden serializes a normalized envelope to the on-disk golden text.
func marshalGolden(env goldenEnvelope) string {
	var b strings.Builder
	fmt.Fprintf(&b, "===== exit =====\n%d\n", env.res.exit)
	fmt.Fprintf(&b, "===== stdout =====\n%s", env.res.stdout)
	fmt.Fprintf(&b, "\n===== stderr =====\n%s", env.res.stderr)
	if env.body != "" {
		fmt.Fprintf(&b, "\n===== body =====\n%s", env.body)
	}
	return b.String()
}

// goldenDispatchPath returns the testdata/golden path for a named case.
func goldenDispatchPath(name string) string {
	return filepath.Join("testdata", "golden", name+".txt")
}

// captureGolden writes the normalized envelope to its golden file (-update mode).
func captureGolden(t *testing.T, name string, env goldenEnvelope) {
	t.Helper()
	path := goldenDispatchPath(name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir golden dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(marshalGolden(env)), 0o644); err != nil {
		t.Fatalf("write golden %s: %v", name, err)
	}
}

// readGolden parses a golden envelope file back into its sections.
func readGolden(t *testing.T, name string) goldenEnvelope {
	t.Helper()
	raw, err := os.ReadFile(goldenDispatchPath(name))
	if err != nil {
		t.Fatalf("read golden %s (regenerate with -update): %v", name, err)
	}
	return parseGolden(t, name, string(raw))
}

// parseGolden splits a marshalled envelope back into its channels + body.
func parseGolden(t *testing.T, name, raw string) goldenEnvelope {
	t.Helper()
	idx := goldenSectionRe.FindAllStringSubmatchIndex(raw, -1)
	if len(idx) == 0 {
		t.Fatalf("golden %s has no sections:\n%s", name, raw)
	}
	sections := map[string]string{}
	for i, m := range idx {
		key := raw[m[2]:m[3]]
		// Body of a section runs from just after its header line to the start of
		// the next header (or EOF). The header line is followed by a newline.
		start := m[1] + 1
		end := len(raw)
		if i+1 < len(idx) {
			// The next header is preceded by a "\n" the marshaller inserted; the
			// regex matched at the start of the header line, so trim the one
			// leading newline that separates this section's content from it.
			end = idx[i+1][0] - 1
		}
		if start > end {
			start = end
		}
		sections[key] = raw[start:end]
	}
	var env goldenEnvelope
	env.res.stdout = sections["stdout"]
	env.res.stderr = sections["stderr"]
	if body, ok := sections["body"]; ok {
		env.body = body
	}
	if _, err := fmt.Sscanf(strings.TrimSpace(sections["exit"]), "%d", &env.res.exit); err != nil {
		t.Fatalf("golden %s exit not an int: %q", name, sections["exit"])
	}
	return env
}

// assertGolden compares the normalized native run (and optional body) against the
// frozen golden, or captures the golden under -update. native and nativeBody must
// already be root/home-normalized by the caller via normRun/normPaths.
func assertGolden(t *testing.T, name string, native goldenEnvelope) {
	t.Helper()
	if *update {
		captureGolden(t, name, native)
		return
	}
	want := readGolden(t, name)
	if native.res.stdout != want.res.stdout {
		t.Errorf("%s: stdout mismatch\n--- native ---\n%q\n--- golden ---\n%q", name, native.res.stdout, want.res.stdout)
	}
	if native.res.stderr != want.res.stderr {
		t.Errorf("%s: stderr mismatch\n--- native ---\n%q\n--- golden ---\n%q", name, native.res.stderr, want.res.stderr)
	}
	if native.res.exit != want.res.exit {
		t.Errorf("%s: exit mismatch native=%d golden=%d", name, native.res.exit, want.res.exit)
	}
	if native.body != want.body {
		t.Errorf("%s: dispatch body mismatch\n--- native ---\n%s\n--- golden ---\n%s", name, native.body, want.body)
	}
}
