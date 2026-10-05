package ensigncycle

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRejectionRoundPublicationExecution(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	const plain = `${SPACEDOCK_BIN:-spacedock} gate record rejection-task --round validation/1 --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl`
	// Minimal retained pre3 and 0.27.3 shapes; the latter quotes the launcher
	// across outer-shell quoting boundaries after a heredoc.
	nested := `/bin/bash -lc "cat <<'PY'
retained heredoc
PY
\""'${SPACEDOCK_BIN:-spacedock}" gate record rejection-task --round validation/1 --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl'`
	for name, command := range map[string]string{
		"pre3":                          "/bin/bash -lc '" + plain + "'",
		"0273":                          nested,
		"echo only":                     "echo '" + plain + "'",
		"failed recorder outer success": strings.Replace(plain, "briefing.json", "absent.json", 1) + "; true",
		"duplicate":                     plain + "; " + plain,
		"second round":                  plain + "; " + strings.Replace(plain, "validation/1", "validation/2", 1),
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			entity := writeRejectionWorkflow(t, root)
			writeFile(t, filepath.Join(root, "rejection-task", "inputs", "briefing.review.jsonl"), rejectionCompleteLog())
			log := filepath.Join(t.TempDir(), "command.log")
			writeFile(t, log, "")
			shim := writeRecordedGateLoggingShim(t, binary, log)
			cmd := exec.Command("/bin/bash", "-c", command)
			cmd.Dir, cmd.Env = root, withRecordedGateEnv(os.Environ(), "SPACEDOCK_BIN", filepath.Join(shim, "spacedock"))
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("shell: %v\n%s", err, output)
			}
			rounds := recorderRejectionRoundPublications(readFile(t, log))
			want := name == "pre3" || name == "0273"
			if err := assertSingleRejectionRoundPublication(rounds); (err == nil) != want {
				t.Fatalf("publications=%v: %v; want accepted=%v", rounds, err, want)
			}
			if want {
				if err := assertRejectionRecordedRound(root, entity, "backlog", true); err != nil {
					t.Fatal(err)
				}
				if err := os.RemoveAll(filepath.Join(root, "rejection-task", "review")); err != nil {
					t.Fatal(err)
				}
				if err := assertRejectionRecordedRound(root, entity, "backlog", true); err == nil {
					t.Fatal("accepted missing round room")
				}
			}
		})
	}
}

func TestRejectionRecorderThroughClaudeLauncher(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	root := t.TempDir()
	writeRejectionWorkflow(t, root)
	log := filepath.Join(t.TempDir(), "command.log")
	shim := writeRecordedGateLoggingShim(t, binary, log)
	command := `"$SPACEDOCK_BIN" gate record rejection-task --round validation/1 --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl`
	// The actual launcher pins SPACEDOCK_BIN; the existing shell startup shim
	// must restore recorder interception inside this host stand-in.
	host := filepath.Join(shim, "claude")
	writeFile(t, host, "#!/bin/sh\nexec /bin/bash -c \"$RECORDER_COMMAND\"\n")
	if err := os.Chmod(host, 0755); err != nil {
		t.Fatal(err)
	}
	env := withSpacedockShimShellEnv(t, withRecordedGateEnv(os.Environ(), "PATH", shim+string(os.PathListSeparator)+os.Getenv("PATH")), shim)
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	run := func(command string) string {
		cmd := exec.Command(binary, "claude", "--plugin-dir", repo, "--skip-compat-check", "--")
		cmd.Dir, cmd.Env = root, withRecordedGateEnv(env, "RECORDER_COMMAND", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("launcher/host shell: %v\n%s", err, output)
		}
		return string(output)
	}
	run(command + "; echo exit=$?")
	writeFile(t, filepath.Join(root, "rejection-task/inputs/briefing.review.jsonl"), rejectionCompleteLog())
	run(command + "; echo exit=$?")
	if !strings.Contains(readFile(t, log), "exit=1\tgate record") || !strings.Contains(readFile(t, log), "exit=0\tgate record") {
		t.Fatal("actual recorder calls were not intercepted")
	}
	if err := assertSingleRejectionRoundPublication(recorderRejectionRoundPublications(readFile(t, log))); err != nil {
		t.Fatal(err)
	}
}
