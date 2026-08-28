package ensigncycle

import (
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
)

const mergedFloorPatch = 178

var claudeVersionPattern = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)`)

// mergedClaudeHost reports whether a parsed Claude Code version supports the
// current in-process named-background Agent substrate. parsed is false when the
// version line has no leading semantic version.
func mergedClaudeHost(versionOutput string) (merged, parsed bool) {
	match := claudeVersionPattern.FindStringSubmatch(strings.TrimSpace(versionOutput))
	if match == nil {
		return false, false
	}
	version := [3]int{}
	for i := range version {
		value, err := strconv.Atoi(match[i+1])
		if err != nil {
			return false, false
		}
		version[i] = value
	}
	floor := [3]int{2, 1, mergedFloorPatch}
	for i := range version {
		if version[i] != floor[i] {
			return version[i] > floor[i], true
		}
	}
	return true, true
}

func cleanupKeepMovingRoot(t *testing.T, root string, failed bool) {
	t.Helper()
	if failed {
		t.Logf("retained failing keep-moving Git root: %s", root)
		return
	}
	if err := os.RemoveAll(root); err != nil {
		t.Errorf("remove successful keep-moving Git root: %v", err)
	}
}

func TestCleanupKeepMovingRootRetainsOnlyFailures(t *testing.T) {
	root := t.TempDir()
	cleanupKeepMovingRoot(t, root, true)
	if _, err := os.Stat(root); err != nil {
		t.Fatalf("failed workflow was not retained: %v", err)
	}
	cleanupKeepMovingRoot(t, root, false)
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatalf("successful workflow still exists: %v", err)
	}
}

func codexLiveFrontDoorArgv(pluginDir, workflowRoot, finalPath, prompt string) []string {
	return []string{
		"codex",
		prompt,
		"--",
		"exec",
		"--json",
		"--dangerously-bypass-approvals-and-sandbox",
		"--cd", workflowRoot,
		"--output-last-message", finalPath,
	}
}

func codexLiveFrontDoorArgvForScenario(pluginDir, workflowRoot, finalPath, prompt, _ string) []string {
	return codexLiveFrontDoorArgv(pluginDir, workflowRoot, finalPath, prompt)
}

func argvHasAdjacent(args []string, left, right string) bool {
	for i := 0; i+1 < len(args); i++ {
		if args[i] == left && args[i+1] == right {
			return true
		}
	}
	return false
}

func TestCodexLiveRunnerLeavesCollaborationConfigToSpacedockLauncher(t *testing.T) {
	args := codexLiveFrontDoorArgv("/tmp/plugin", "/tmp/workflow", "/tmp/final-message.txt", "run the scenario")
	if argvHasAdjacent(args, "--enable", "multi_agent_v2") || argvHasAdjacent(args, "--disable", "multi_agent_v2") {
		t.Fatalf("Codex live argv must leave collaboration configuration to the Spacedock launcher; args=%v", args)
	}
}

func TestCodexLiveRunnerUsesSpacedockFrontDoorBeforeHostArgs(t *testing.T) {
	args := codexLiveFrontDoorArgv("/tmp/plugin", "/tmp/workflow", "/tmp/final-message.txt", "run the scenario")
	fence := -1
	for i, arg := range args {
		if arg == "--" {
			fence = i
			break
		}
	}
	if fence < 0 {
		t.Fatalf("Codex live argv has no host-argument fence: %v", args)
	}
	if args[0] != "codex" {
		t.Fatalf("Codex front door is not first: %v", args)
	}
	for _, arg := range args {
		if arg == "--plugin-dir" || arg == "/tmp/plugin" || arg == "--skip-compat-check" {
			t.Fatalf("common live runner bypassed the installed stable package: %v", args)
		}
	}
	if args[fence+1] != "exec" {
		t.Fatalf("Codex host argv does not start with exec: %v", args)
	}
	if !argvHasAdjacent(args, "--dangerously-bypass-approvals-and-sandbox", "--cd") {
		t.Fatalf("Codex live argv does not preserve bypass-permission posture before workflow root: %v", args)
	}
}

func TestCodexLiveRunnerUsesCommonPostureForFiling(t *testing.T) {
	base := codexLiveFrontDoorArgv("/tmp/plugin", "/tmp/workflow", "/tmp/final-message.txt", "run the scenario")
	for _, scenario := range []string{"filing", "shallow-boot"} {
		if got := codexLiveFrontDoorArgvForScenario("/tmp/plugin", "/tmp/workflow", "/tmp/final-message.txt", "run the scenario", scenario); !slices.Equal(got, base) {
			t.Fatalf("%s argv changed: got=%v want=%v", scenario, got, base)
		}
	}
}

func TestMergedClaudeHost(t *testing.T) {
	for _, tc := range []struct {
		name           string
		version        string
		merged, parsed bool
	}{
		{"last native-team version", "2.1.177 (Claude Code)", false, true},
		{"merged floor", "2.1.178 (Claude Code)", true, true},
		{"current merged version", "2.1.181 (Claude Code)", true, true},
		{"future minor", "2.2.0 (Claude Code)", true, true},
		{"older minor", "2.0.250 (Claude Code)", false, true},
		{"unparseable", "garbage", false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			merged, parsed := mergedClaudeHost(tc.version)
			if merged != tc.merged || parsed != tc.parsed {
				t.Fatalf("mergedClaudeHost(%q) = (%t, %t), want (%t, %t)", tc.version, merged, parsed, tc.merged, tc.parsed)
			}
		})
	}
}
