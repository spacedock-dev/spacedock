package contractlint

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var runtimeBindingCapabilities = []string{
	"worker.spawn",
	"addressable-worker",
	"async-dispatch",
	"worker-identity",
	"completion-signal",
	"worker.shutdown",
	"context-budget",
	"roster-reconcile",
}

func TestCodexAndPiFirstOfficerRuntimeBindingBlocks(t *testing.T) {
	for _, rel := range []string{
		filepath.Join("skills", "first-officer", "references", "codex-first-officer-runtime.md"),
		filepath.Join("skills", "first-officer", "references", "pi-first-officer-runtime.md"),
	} {
		t.Run(rel, func(t *testing.T) {
			text := readRepoFile(t, rel)
			if count := strings.Count(text, "## Runtime implementation\n"); count != 1 {
				t.Fatalf("%s has %d `## Runtime implementation` sections, want exactly 1", rel, count)
			}
			section := markdownSectionFromText(t, text, "## Runtime implementation")
			got := bindingBulletCapabilities(section)
			if !reflect.DeepEqual(got, runtimeBindingCapabilities) {
				t.Fatalf("%s runtime binding capabilities mismatch:\n  got:  %v\n  want: %v", rel, got, runtimeBindingCapabilities)
			}
		})
	}
}

func TestCodexAndPiFirstOfficerRuntimeLifecycleHeadingsRemoved(t *testing.T) {
	cases := map[string][]string{
		filepath.Join("skills", "first-officer", "references", "codex-first-officer-runtime.md"): {
			"## Dispatch",
			"## Awaiting Completion",
			"## Reuse And Feedback Routing",
		},
		filepath.Join("skills", "first-officer", "references", "pi-first-officer-runtime.md"): {
			"## Runtime Shape",
			"## Dispatch",
			"## Awaiting Completion",
			"## Follow-up and Reuse",
			"## Shutdown",
			"### Model Resolution",
			"### Canonical Model Space",
		},
	}
	for rel, banned := range cases {
		text := readRepoFile(t, rel)
		for _, heading := range banned {
			if strings.Contains(text, heading+"\n") {
				t.Errorf("%s still contains lifecycle heading %q", rel, heading)
			}
		}
	}
}

func TestCodexAndPiFirstOfficerRuntimeRejectMutableStepAndNegativeContrast(t *testing.T) {
	for _, rel := range []string{
		filepath.Join("skills", "first-officer", "references", "codex-first-officer-runtime.md"),
		filepath.Join("skills", "first-officer", "references", "pi-first-officer-runtime.md"),
	} {
		text := readRepoFile(t, rel)
		for _, banned := range []string{
			"Dispatch step",
			"Event Loop step",
			"Merge-and-Cleanup step",
			"step 10",
			"does not expose Claude Code team-tool signatures",
			"Do not call or ask workers to call Claude team tools",
			"Pi has no such enum",
			"Claude-centric enum",
			"Do not create teams",
			"Codex declares none",
		} {
			if strings.Contains(text, banned) {
				t.Errorf("%s contains mutable step coupling or negative host contrast %q", rel, banned)
			}
		}
	}
}

func TestPiToolNamesStayInRuntimeBindingOrHarnessSections(t *testing.T) {
	rel := filepath.Join("skills", "first-officer", "references", "pi-first-officer-runtime.md")
	text := readRepoFile(t, rel)
	allowed := markdownSectionFromText(t, text, "## Runtime implementation") + "\n" +
		markdownSectionFromText(t, text, "## Live Harness Isolation")
	outside := text
	for _, heading := range []string{"## Runtime implementation", "## Live Harness Isolation"} {
		_, remainder := extractMarkdownSection(t, outside, heading)
		outside = remainder
	}

	for _, want := range []string{
		"subagent(",
		"intercom(",
		"member_spawn",
		"delegate",
		"message_dm",
		"member_shutdown",
		"team_done",
		`context: "fresh"`,
		"cwd: <resolved repo root>",
		"acceptance",
	} {
		if !strings.Contains(allowed, want) {
			t.Errorf("%s allowed runtime/harness sections missing %q", rel, want)
		}
	}

	for _, banned := range []string{
		"subagent(",
		"intercom(",
		"member_spawn",
		"delegate",
		"message_dm",
		"member_shutdown",
		"team_done",
		`context: "fresh"`,
		"cwd: <resolved repo root>",
		"acceptance",
	} {
		if strings.Contains(outside, banned) {
			t.Errorf("%s contains Pi tool token outside runtime implementation or harness section: %q", rel, banned)
		}
	}
}

func TestCodexAndPiFirstOfficerRuntimeSemanticsPreserved(t *testing.T) {
	cases := map[string][]string{
		filepath.Join("skills", "first-officer", "references", "codex-first-officer-runtime.md"): {
			"live tool surface",
			"`spawn_agent(task_name,message,fork_turns)`",
			"sanitize the helper-emitted `name`",
			"`wait_agent(timeout_ms)`",
			"wait timeout return is normal and retryable",
			"worker is not failed, closed, or redispatched",
			"next foreground wait is reinstalled",
			"queued/activity-driven delivery",
			"autonomous FO wake-up",
			"re-run the kept-alive validation reviewer through the same turn-starting `«addressable-worker»` binding",
		},
		filepath.Join("skills", "first-officer", "references", "pi-first-officer-runtime.md"): {
			"Pi-native substrate selected by the launch/test harness",
			"`subagent(...)` with explicit `context: \"fresh\"` and `cwd: <resolved repo root>`",
			"do not use the `subagent(... acceptance: ...)` contract",
			"`intercom({action:\"list\"})`",
			"Pi's model-space binding is provider/model strings",
			"file verification remains the completion gate",
			"Fresh redispatch remains the default",
			"non-fresh resume is only an explicit manual/debug exception",
			"completed child invocation needs no mailbox shutdown",
			"`member_shutdown`",
		},
	}
	for rel, wants := range cases {
		text := readRepoFile(t, rel)
		for _, want := range wants {
			if !strings.Contains(text, want) {
				t.Errorf("%s missing preserved runtime semantic %q", rel, want)
			}
		}
	}
}

func readRepoFile(t *testing.T, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot(t), rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(data)
}

func markdownSectionFromText(t *testing.T, text, heading string) string {
	t.Helper()
	section, _ := extractMarkdownSection(t, text, heading)
	return section
}

func bindingBulletCapabilities(section string) []string {
	re := regexp.MustCompile(`(?m)^- ` + "`" + `«([^»]+)»` + "`" + `:`)
	out := []string{}
	for _, match := range re.FindAllStringSubmatch(section, -1) {
		out = append(out, match[1])
	}
	return out
}
