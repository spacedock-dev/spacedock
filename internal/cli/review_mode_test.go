package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runIn(t *testing.T, dir string, env []string, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := run(context.Background(), args, env, dir, strings.NewReader(""), &stdout, &stderr, nil, nil)
	return code, stdout.String(), stderr.String()
}

func reviewModeEnv(t *testing.T) []string {
	t.Helper()
	return append(os.Environ(), "XDG_STATE_HOME="+t.TempDir())
}

func TestReviewModeIsUnsetUntilChosenAndScopedToTheProject(t *testing.T) {
	env := reviewModeEnv(t)
	project, other := t.TempDir(), t.TempDir()

	code, out, _ := runIn(t, project, env, "review-mode", "get", "--json")
	if code != 0 || !strings.Contains(out, `"mode":""`) {
		t.Fatalf("fresh project: code=%d out=%q, want an unset mode", code, out)
	}
	if code, _, errOut := runIn(t, project, env, "review-mode", "set", "remote"); code != 0 {
		t.Fatalf("set remote: code=%d stderr=%q", code, errOut)
	}
	_, out, _ = runIn(t, project, env, "review-mode", "get", "--json")
	var got reviewSetting
	if err := json.Unmarshal([]byte(out), &got); err != nil || got.Mode != "remote" || got.DecidedAt == "" {
		t.Fatalf("after set: %q (err %v), want mode remote with a decision time", out, err)
	}
	if _, out, _ = runIn(t, other, env, "review-mode", "get", "--json"); !strings.Contains(out, `"mode":""`) {
		t.Fatalf("another project inherited the choice: %q", out)
	}
	runIn(t, project, env, "review-mode", "clear")
	if _, out, _ = runIn(t, project, env, "review-mode", "get", "--json"); !strings.Contains(out, `"mode":""`) {
		t.Fatalf("after clear: %q, want unset", out)
	}
	if code, _, _ := runIn(t, project, env, "review-mode", "set", "everywhere"); code != 2 {
		t.Fatalf("set with an unknown mode returned %d, want usage error 2", code)
	}
}

func TestRemoteReviewRefusesWithoutRemoteConsent(t *testing.T) {
	env := reviewModeEnv(t)
	project := t.TempDir()
	bin := t.TempDir()
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$0.args\"\n"
	if err := os.WriteFile(filepath.Join(bin, "spacedock-remote-review"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	env = append(env, "PATH="+os.Getenv("PATH"))
	argsFile := filepath.Join(bin, "spacedock-remote-review.args")

	for _, mode := range []string{"", "local"} {
		if mode != "" {
			runIn(t, project, env, "review-mode", "set", mode)
		}
		code, _, errOut := runIn(t, project, env, "remote-review", "plan.md")
		if code != 3 || !strings.Contains(errOut, "review-mode set remote") {
			t.Fatalf("mode %q: code=%d stderr=%q, want refusal 3 naming the consent command", mode, code, errOut)
		}
		if _, err := os.Stat(argsFile); err == nil {
			t.Fatalf("mode %q: the uploader ran without remote consent", mode)
		}
	}

	runIn(t, project, env, "review-mode", "set", "remote")
	if code, _, errOut := runIn(t, project, env, "remote-review", "plan.md", "--dry-run"); code != 0 {
		t.Fatalf("with consent: code=%d stderr=%q", code, errOut)
	}
	got, err := os.ReadFile(argsFile)
	if err != nil || string(got) != "publish\nplan.md\n--dry-run\n" {
		t.Fatalf("uploader argv = %q (err %v), want publish plan.md --dry-run", got, err)
	}
}
