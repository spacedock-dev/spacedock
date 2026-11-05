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

func TestReviewModeRefusesAnUnreadableChoiceAndClearRepairsIt(t *testing.T) {
	for name, content := range map[string]string{
		"corrupt":      "{not json",
		"unknown mode": `{"mode":"banana"}`,
	} {
		t.Run(name, func(t *testing.T) {
			env := reviewModeEnv(t)
			project := t.TempDir()
			_, path := reviewSettingPath(env, project)
			if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
				t.Fatal(err)
			}

			if code, out, _ := runIn(t, project, env, "review-mode", "get", "--json"); code != 1 {
				t.Fatalf("get on %s file: code=%d out=%q, want exit 1", name, code, out)
			}
			if code, _, errOut := runIn(t, project, env, "review-mode", "set", "remote"); code != 0 {
				t.Fatalf("set over a %s file: code=%d stderr=%q, want it to replace the choice", name, code, errOut)
			}
			if _, out, _ := runIn(t, project, env, "review-mode", "get", "--json"); !strings.Contains(out, `"mode":"remote"`) {
				t.Fatalf("get after set over a %s file: %q, want remote", name, out)
			}
			if code, _, errOut := runIn(t, project, env, "review-mode", "clear"); code != 0 {
				t.Fatalf("clear on %s file: code=%d stderr=%q, want it to repair the choice", name, code, errOut)
			}
			if code, out, _ := runIn(t, project, env, "review-mode", "get", "--json"); code != 0 || !strings.Contains(out, `"mode":""`) {
				t.Fatalf("get after clear: code=%d out=%q, want an unset mode", code, out)
			}
		})
	}
}

func TestReviewModeSetReportsAWriteFailure(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("a read-only directory does not stop root")
	}
	root := t.TempDir()
	dir := filepath.Join(root, "subspace", "review-mode")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o700) })
	env := append(os.Environ(), "XDG_STATE_HOME="+root)

	if code, out, _ := runIn(t, t.TempDir(), env, "review-mode", "set", "remote"); code != 1 {
		t.Fatalf("set into a read-only state directory returned %d (%q), want 1", code, out)
	}
}
