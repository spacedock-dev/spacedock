// ABOUTME: Runs the capability-check lines shipped in skills/review/SKILL.md against stub
// ABOUTME: programs on a controlled PATH: refusal, support, and a missing program.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var supportsLineRe = regexp.MustCompile(`^(spacedock-review|subspace-relay-client) --supports \S+ >/dev/null 2>&1 \|\| echo "UNSUPPORTED:`)

func reviewSupportLines(t *testing.T) map[string]string {
	t.Helper()
	path := filepath.Join(repoRoot(t), "skills", "review", "SKILL.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := map[string]string{}
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if m := supportsLineRe.FindStringSubmatch(line); m != nil {
			if _, dup := lines[m[1]]; dup {
				t.Fatalf("two runnable check lines for %s", m[1])
			}
			lines[m[1]] = line
		}
	}
	if len(lines) != 2 {
		t.Fatalf("want one runnable check line for each of two programs, found %v", lines)
	}
	return lines
}

func stubProgram(t *testing.T, dir, name, version, supported, log string) {
	t.Helper()
	body := "#!/bin/sh\necho \"$*\" >> '" + log + "'\n" +
		"[ \"$1\" = --version ] && { echo '" + version + "'; exit 0; }\n" +
		"[ \"$1\" = --supports ] && [ \"$2\" = '" + supported + "' ] && exit 0\n" +
		"exit 1\n"
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
}

func runCheck(t *testing.T, line, pathDir string) string {
	t.Helper()
	cmd := exec.Command("sh", "-c", line)
	cmd.Env = []string{"PATH=" + pathDir + ":/usr/bin:/bin"}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("check line failed to run: %v\n%s", err, out)
	}
	return strings.TrimSpace(string(out))
}

func TestReviewSupportCheckReportsUnsupportedWithoutCallingPublish(t *testing.T) {
	tokens := map[string]string{"spacedock-review": "review-invocation-prepare-v1", "subspace-relay-client": "relay-publish-v1"}
	for program, line := range reviewSupportLines(t) {
		t.Run(program, func(t *testing.T) {
			token := tokens[program]

			dir, log := t.TempDir(), filepath.Join(t.TempDir(), "calls")
			stubProgram(t, dir, program, "0.16.0-beta.1", "an-old-capability", log)
			out := runCheck(t, line, dir)
			if want := "UNSUPPORTED: " + program + " 0.16.0-beta.1"; out != want {
				t.Errorf("old program: output %q, want %q", out, want)
			}
			calls, _ := os.ReadFile(log)
			if want := "--supports " + token + "\n--version\n"; string(calls) != want {
				t.Errorf("old program was called with %q, want exactly %q", calls, want)
			}

			dir, log = t.TempDir(), filepath.Join(t.TempDir(), "calls")
			stubProgram(t, dir, program, "0.16.0-beta.2", token, log)
			if out := runCheck(t, line, dir); out != "" {
				t.Errorf("supporting program: output %q, want none", out)
			}

			if out := runCheck(t, line, t.TempDir()); out != "UNSUPPORTED: "+program+" not installed" {
				t.Errorf("missing program: output %q", out)
			}
		})
	}
}
