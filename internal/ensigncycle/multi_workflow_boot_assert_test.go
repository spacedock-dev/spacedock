package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

type multiWorkflowBootObservation struct {
	commands            []string
	finalMessage        string
	workflowDirs        []string
	entityBefore        []string
	entityAfter         []string
	gitHeadBefore       string
	gitHeadAfter        string
	gitStatusBefore     string
	gitStatusAfter      string
	convergenceArtifact bool
}

var (
	shellCommandSegments = regexp.MustCompile(`\r?\n|;|&&|\|\||\|`)
	spacedockLauncher    = `(?:spacedock(?:_launcher)?|"[^"\r\n]*/spacedock"|[^ \t'"\r\n]*/spacedock|\$\{SPACEDOCK_BIN[^}]*\}|"?\$[A-Za-z_][A-Za-z0-9_]*"?)`
	spacedockStatusCall  = regexp.MustCompile(spacedockLauncher + `[ \t]+status\b`)
	spacedockStateCall   = regexp.MustCompile(spacedockLauncher + `[ \t]+state\b`)
	retryHelperCall      = regexp.MustCompile(`(?:^|[ \t'"\\])(?:jq|python3)(?:[ \t'"\\]|$)|(?:^|[ \t'"\\])go[ \t]+run(?:[ \t'"\\]|$)`)
)

func assertMultiWorkflowBoot(o multiWorkflowBootObservation) error {
	bootIdentifyCalls := 0
	statusRetries := 0
	helperRetries := 0
	convergenceCalls := 0
	workflowSpecificCalls := 0
	for _, command := range o.commands {
		for _, segment := range shellCommandSegments.Split(command, -1) {
			segment = strings.TrimSpace(segment)
			if segment == "" {
				continue
			}
			if spacedockStatusCall.MatchString(segment) {
				if strings.Contains(segment, "--workflow-dir") {
					workflowSpecificCalls++
				}
				if strings.Contains(segment, "--boot") && strings.Contains(segment, "--identify") && strings.Contains(segment, "--json") {
					bootIdentifyCalls++
				} else {
					statusRetries++
				}
			}
			if retryHelperCall.MatchString(segment) {
				helperRetries++
			}
			if spacedockStateCall.MatchString(segment) {
				convergenceCalls++
			}
		}
	}
	if bootIdentifyCalls != 1 {
		return fmt.Errorf("boot identify calls = %d, want exactly 1", bootIdentifyCalls)
	}
	if statusRetries != 0 || helperRetries != 0 {
		return fmt.Errorf("startup retried before workflow selection: status=%d helper=%d, want 0/0", statusRetries, helperRetries)
	}
	if convergenceCalls != 0 {
		return fmt.Errorf("startup converged workflow state %d time(s), want 0 before selection", convergenceCalls)
	}
	if workflowSpecificCalls != 0 {
		return fmt.Errorf("startup selected/deep-booted a workflow %d time(s), want 0 before operator selection", workflowSpecificCalls)
	}
	if !containsExactLine(o.finalMessage, multiWorkflowSelectionGreeting) {
		return fmt.Errorf("greeting missing exact workflow-selection sentence %q", multiWorkflowSelectionGreeting)
	}
	for _, workflowDir := range o.workflowDirs {
		if !strings.Contains(o.finalMessage, workflowDir) {
			return fmt.Errorf("greeting did not name discovered workflow %q", workflowDir)
		}
	}
	if len(o.entityBefore) != len(o.entityAfter) {
		return fmt.Errorf("entity snapshot count changed: before=%d after=%d", len(o.entityBefore), len(o.entityAfter))
	}
	for i := range o.entityBefore {
		if o.entityBefore[i] != o.entityAfter[i] {
			return fmt.Errorf("workflow entity %d changed before selection", i)
		}
	}
	if o.gitHeadBefore != o.gitHeadAfter {
		return fmt.Errorf("project git HEAD changed before selection: %s -> %s", o.gitHeadBefore, o.gitHeadAfter)
	}
	if o.gitStatusBefore != "" || o.gitStatusAfter != o.gitStatusBefore {
		return fmt.Errorf("project git status changed before selection: before=%q after=%q", o.gitStatusBefore, o.gitStatusAfter)
	}
	if o.convergenceArtifact {
		return fmt.Errorf("startup left a worktree/archive convergence artifact before selection")
	}
	return nil
}

func containsExactLine(text, want string) bool {
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSuffix(line, "\r") == want {
			return true
		}
	}
	return false
}

func gatherMultiWorkflowBootObservation(t *testing.T, fx multiWorkflowBootFixture, commands []string, finalMessage string) multiWorkflowBootObservation {
	t.Helper()
	obs := multiWorkflowBootObservation{
		commands:        commands,
		finalMessage:    finalMessage,
		workflowDirs:    append([]string(nil), fx.workflowDirs...),
		entityBefore:    append([]string(nil), fx.entityBefore...),
		gitHeadBefore:   fx.gitHeadBefore,
		gitStatusBefore: fx.gitStatusBefore,
		gitHeadAfter:    strings.TrimSpace(git(t, fx.root, "rev-parse", "HEAD")),
		gitStatusAfter:  strings.TrimSpace(git(t, fx.root, "status", "--short")),
	}
	for _, path := range fx.entityPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			obs.entityAfter = append(obs.entityAfter, "")
			continue
		}
		obs.entityAfter = append(obs.entityAfter, string(data))
	}
	for _, pattern := range []string{
		filepath.Join(fx.root, ".worktrees", "*"),
		filepath.Join(fx.root, "*", ".worktrees", "*"),
		filepath.Join(fx.root, "*", "_archive", "*"),
	} {
		if matches, _ := filepath.Glob(pattern); len(matches) > 0 {
			obs.convergenceArtifact = true
		}
	}
	return obs
}

func codexExecutedCommands(jsonl string) []string {
	var commands []string
	for _, line := range strings.Split(jsonl, "\n") {
		var event codexCommandItem
		if json.Unmarshal([]byte(line), &event) == nil && event.Type == "item.completed" && event.Item.Type == "command_execution" {
			commands = append(commands, event.Item.Command)
		}
	}
	return commands
}

func claudeExecutedCommands(stream string) []string {
	var commands []string
	for _, line := range strings.Split(stream, "\n") {
		var event struct {
			Message *struct {
				Content []struct {
					Type  string `json:"type"`
					Name  string `json:"name"`
					Input struct {
						Command string `json:"command"`
					} `json:"input"`
				} `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal([]byte(line), &event) != nil || event.Message == nil {
			continue
		}
		for _, block := range event.Message.Content {
			if block.Type == "tool_use" && block.Name == "Bash" {
				commands = append(commands, block.Input.Command)
			}
		}
	}
	return commands
}

func TestMultiWorkflowBootAfterOnlyInvariant(t *testing.T) {
	paths := []string{"/project/alpha", "/project/beta"}
	entity := multiWorkflowBootEntity("alpha")
	good := multiWorkflowBootObservation{
		commands:        []string{`${SPACEDOCK_BIN:-spacedock} status --boot --identify --json`},
		finalMessage:    multiWorkflowSelectionGreeting + "\n" + strings.Join(paths, "\n"),
		workflowDirs:    paths,
		entityBefore:    []string{entity, entity},
		entityAfter:     []string{entity, entity},
		gitHeadBefore:   "abc",
		gitHeadAfter:    "abc",
		gitStatusBefore: "",
		gitStatusAfter:  "",
	}
	if err := assertMultiWorkflowBoot(good); err != nil {
		t.Fatalf("after-only multi-workflow baseline must pass: %v", err)
	}
	absoluteLauncher := good
	absoluteLauncher.commands = []string{`"/tmp/bin/spacedock" status --boot --identify --json`}
	if err := assertMultiWorkflowBoot(absoluteLauncher); err != nil {
		t.Fatalf("quoted absolute launcher baseline must pass: %v", err)
	}

	broken := map[string]func(*multiWorkflowBootObservation){
		"duplicate identify": func(o *multiWorkflowBootObservation) { o.commands = append(o.commands, o.commands[0]) },
		"status retry": func(o *multiWorkflowBootObservation) {
			o.commands = append(o.commands, `spacedock status --boot --json`)
		},
		"jq retry": func(o *multiWorkflowBootObservation) {
			o.commands = append(o.commands, `/bin/bash -lc 'jq . boot.json'`)
		},
		"python retry": func(o *multiWorkflowBootObservation) {
			o.commands = append(o.commands, `/bin/bash -lc 'python3 check.py'`)
		},
		"go-run retry": func(o *multiWorkflowBootObservation) {
			o.commands = append(o.commands, `/bin/bash -lc 'go run ./cmd/spacedock status --boot --identify --json'`)
		},
		"convergence": func(o *multiWorkflowBootObservation) { o.commands = append(o.commands, `spacedock state ready`) },
		"workflow boot": func(o *multiWorkflowBootObservation) {
			o.commands = []string{`spacedock status --boot --identify --json --workflow-dir /project/alpha`}
		},
		"wrong greeting": func(o *multiWorkflowBootObservation) { o.finalMessage = strings.Join(paths, "\n") },
		"embedded greeting": func(o *multiWorkflowBootObservation) {
			o.finalMessage = "Selection: " + multiWorkflowSelectionGreeting + "\n" + strings.Join(paths, "\n")
		},
		"unnamed workflow": func(o *multiWorkflowBootObservation) {
			o.finalMessage = multiWorkflowSelectionGreeting + "\n" + paths[0]
		},
		"entity mutation": func(o *multiWorkflowBootObservation) { o.entityAfter[0] += "mutated\n" },
		"git mutation":    func(o *multiWorkflowBootObservation) { o.gitStatusAfter = " M alpha/held-task.md" },
		"artifact":        func(o *multiWorkflowBootObservation) { o.convergenceArtifact = true },
	}
	for name, mutate := range broken {
		t.Run(name, func(t *testing.T) {
			o := good
			o.commands = append([]string(nil), good.commands...)
			o.entityAfter = append([]string(nil), good.entityAfter...)
			mutate(&o)
			if err := assertMultiWorkflowBoot(o); err == nil {
				t.Fatalf("broken %s observation passed", name)
			}
		})
	}
}
