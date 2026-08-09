// ABOUTME: Structural guards binding runtime-live-e2e.yml's live test steps to the
// ABOUTME: pinned-gotestsum clean-log + jsonl-archive shape, with no firehose regression.
package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// transformedLiveSteps are every live test step that must run through gotestsum:
// a single `gotestsum --jsonfile <name>-detail.jsonl --format <clean> -- <args>`
// invocation — one run producing the clean step log AND the archived -json detail,
// with the exit code preserved.
var transformedLiveSteps = []string{
	"Run live Claude E2E",
	"Run live Codex shared scenarios",
}

// TestLiveWorkflowStepsUseGotestsumOneRunShape pins each transformed live step to
// the source-side discipline: gotestsum runs the suite once, archives the -json
// detail to a `--jsonfile`, and preserves the exit. This is the workflow binding
// behind AC-1..AC-4 — the behavioral test proves the SHAPE works; this proves
// every live step IS that shape.
func TestLiveWorkflowStepsUseGotestsumOneRunShape(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	steps := parseWorkflowSteps(live)

	byName := map[string]string{}
	for _, s := range steps {
		if s.run != "" {
			byName[s.name] = s.run
		}
	}

	for _, name := range transformedLiveSteps {
		run, ok := byName[name]
		if !ok {
			t.Errorf("workflow missing transformed live step %q", name)
			continue
		}
		for _, want := range []string{
			"gotestsum",     // the one-run clean+archive renderer
			"--jsonfile",    // the archived -json event stream
			"-detail.jsonl", // the archive filename
			" -- ",          // gotestsum passes the go test args after --
		} {
			if !strings.Contains(run, want) {
				t.Errorf("live step %q lost the gotestsum one-run-archive element %q:\n%s", name, want, run)
			}
		}
	}
}

// TestLiveWorkflowInstallsPinnedGotestsum is AC-5: every live job installs
// gotestsum via the pinned, sha256-verifying install script — never a floating
// `@latest` or an unverified download. The install script itself must pin a fixed
// version and verify a checksum before trusting the binary.
func TestLiveWorkflowInstallsPinnedGotestsum(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	steps := parseWorkflowSteps(live)

	installs := 0
	for _, s := range steps {
		if strings.Contains(s.run, "install-gotestsum.sh") {
			installs++
		}
	}
	// The offline gate and the two live jobs each install gotestsum.
	if installs != 3 {
		t.Errorf("expected gotestsum in the offline gate and both live jobs, found %d installs", installs)
	}

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
