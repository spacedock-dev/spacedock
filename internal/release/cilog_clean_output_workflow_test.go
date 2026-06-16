// ABOUTME: Structural guards binding runtime-live-e2e.yml's live test steps to the
// ABOUTME: clean-stdout + one-run jsonl-archive shape, with no firehose regression.
package release

import (
	"strings"
	"testing"
)

// transformedLiveSteps are every live test step that must carry the clean-output +
// one-run archive shape: a single `go test -c` compile, the test binary tee'd to
// `go tool test2json` for the .jsonl archive, and a `${PIPESTATUS[0]}` exit
// capture so the trailing clean-view grep cannot mask a test failure.
var transformedLiveSteps = []string{
	"Run live ensign cycle",
	"Run live Claude shared scenarios",
	"Run live Codex shared scenarios",
	"Run Pi shared scenario coverage guard",
	"Run live Pi front-door smoke",
}

// TestLiveWorkflowStepsUseOneRunCleanArchiveShape pins each transformed live step
// to the source-side discipline: compile once, run once, archive the -json detail
// via the toolchain's test2json, and preserve the exit code. This is the workflow
// binding behind AC-1..AC-4 — the behavioral test proves the SHAPE works; this
// proves every live step IS that shape.
func TestLiveWorkflowStepsUseOneRunCleanArchiveShape(t *testing.T) {
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
			"go test -c",        // compile once, no run
			"go tool test2json", // toolchain archive renderer (no third-party dep)
			"-detail.jsonl",     // the archived -json event stream
			"${PIPESTATUS[0]}",  // capture the test binary's exit past the grep
			"grep -vE",          // the clean-view filter to the step log
		} {
			if !strings.Contains(run, want) {
				t.Errorf("live step %q lost the one-run clean-archive element %q:\n%s", name, want, run)
			}
		}
	}
}

// TestLiveWorkflowHasNoFirehoseRegression is the adversarial twin: the old
// `-v | tee *-transcript.txt` firehose (the verbose surface that flooded
// FO/ensign reads) must not reappear on any executable line, and no `go test -v`
// may pipe to a `tee`d transcript again.
func TestLiveWorkflowHasNoFirehoseRegression(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	for _, step := range parseWorkflowSteps(live) {
		for _, command := range executableShellCommands(step.run) {
			if strings.Contains(command, "tee") && strings.Contains(command, "transcript.txt") {
				t.Errorf("live step %q re-introduces a -v transcript firehose: %q", step.name, command)
			}
			// A `go test ... -v` (the driver firehose) piped to tee is the exact
			// shape replaced; the binary-side `-test.v` feeding test2json is fine.
			if strings.Contains(command, "go test ") && strings.Contains(command, " -v") && strings.Contains(command, "tee") {
				t.Errorf("live step %q runs `go test -v | tee` again (firehose to stdout): %q", step.name, command)
			}
		}
	}
}

// TestLiveWorkflowPinsNoFloatingTool is AC-5: the chosen mechanism adds no
// third-party build/test dependency, so there is no `go run <pkg>@latest` (or any
// floating @latest) anywhere in the live test steps. If a future change adopts a
// tool, this guard forces it to be pinned, not floating.
func TestLiveWorkflowPinsNoFloatingTool(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	for _, step := range parseWorkflowSteps(live) {
		for _, command := range executableShellCommands(step.run) {
			if strings.Contains(command, "@latest") {
				t.Errorf("live step %q uses a floating @latest dependency (must be pinned): %q", step.name, command)
			}
			if strings.Contains(command, "gotestsum") {
				t.Errorf("live step %q references gotestsum; the chosen mechanism is stdlib-only (go tool test2json) — if a tool is adopted it must be pinned and this guard updated: %q", step.name, command)
			}
		}
	}
}
