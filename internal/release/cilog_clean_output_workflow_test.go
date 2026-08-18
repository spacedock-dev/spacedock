// ABOUTME: Structural guards binding runtime-live-e2e.yml's live test steps to the
// ABOUTME: pinned-gotestsum clean-log + jsonl-archive shape, with no firehose regression.
package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLiveWorkflowInstallsPinnedGotestsum is AC-5: every live job installs
// gotestsum via the pinned, sha256-verifying install script — never a floating
// `@latest` or an unverified download. The install script itself must pin a fixed
// version and verify a checksum before trusting the binary.
func TestLiveWorkflowInstallsPinnedGotestsum(t *testing.T) {
	// Read the EXECUTABLE commands, not the raw text: commenting out every
	// verify/pin line leaves the phrases in the file but renders the script inert,
	// so a raw strings.Contains would stay green on a script that no longer
	// verifies anything. executableShellCommands strips comment/blank lines and
	// joins continuations, so a commented-out verification disappears here.
	active := strings.Join(executableShellCommands(readGotestsumInstallScript(t)), "\n")
	if !strings.Contains(active, `VERSION="`) {
		t.Error("install-gotestsum.sh does not pin a fixed VERSION on an executable line")
	}
	if strings.Contains(active, "@latest") || strings.Contains(active, "/latest/") {
		t.Error("install-gotestsum.sh uses a floating @latest reference (must be a pinned release)")
	}
	if !strings.Contains(active, "shasum -a 256 -c") {
		t.Error("install-gotestsum.sh does not sha256-verify the downloaded tarball before use")
	}
	// At least one concrete pinned checksum must be present (the linux_amd64 one
	// the CI runners use); a bare `sha=""` would defeat the verification.
	if !strings.Contains(active, "linux_amd64)") || !strings.Contains(active, "sha=\"") {
		t.Error("install-gotestsum.sh is missing a concrete pinned checksum for the CI runner platform")
	}
}

// TestLiveWorkflowHasNoFirehoseRegression is the adversarial twin: the old
// `-v | tee *-transcript.txt` firehose (the verbose surface that flooded
// FO/ensign reads) must not reappear, and no `go test ... -v` may pipe to a tee'd
// transcript again on any executable line.
func TestLiveWorkflowHasNoFirehoseRegression(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	for _, step := range parseWorkflowSteps(live) {
		for _, command := range executableShellCommands(step.run) {
			if strings.Contains(command, "tee") && strings.Contains(command, "transcript.txt") {
				t.Errorf("live step %q re-introduces a -v transcript firehose: %q", step.name, command)
			}
			if strings.Contains(command, "go test ") && strings.Contains(command, " -v") && strings.Contains(command, "tee") {
				t.Errorf("live step %q runs `go test -v | tee` again (firehose to stdout): %q", step.name, command)
			}
		}
	}
}

// TestLiveWorkflowPinsNoFloatingTool is AC-5: no `go run <pkg>@latest` (or any
// floating @latest) anywhere in the live test steps — the gotestsum dependency is
// pinned and installed as a verified prebuilt, not floated.
func TestLiveWorkflowPinsNoFloatingTool(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	for _, step := range parseWorkflowSteps(live) {
		for _, command := range executableShellCommands(step.run) {
			if strings.Contains(command, "@latest") {
				t.Errorf("live step %q uses a floating @latest dependency (must be pinned): %q", step.name, command)
			}
		}
	}
}

func readGotestsumInstallScript(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "scripts", "install-gotestsum.sh"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
