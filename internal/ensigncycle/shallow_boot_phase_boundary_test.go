// ABOUTME: Pins neutral shallow-boot phase tasks and the greet-read-only/engage-mutates oracle.
// ABOUTME: Synthetic trajectories isolate every forbidden early mutation and incomplete engage result.
package ensigncycle

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestShallowBootTeamRegistryDistinguishesLeaderFromWorker(t *testing.T) {
	root := t.TempDir()
	leaderOnly := filepath.Join(root, "leader-only.json")
	writeFile(t, leaderOnly, `{"members":[{"name":"team-lead"}]}`)
	if shallowBootTeamHasWorker([]string{leaderOnly}) {
		t.Fatal("Claude's transport-created leader-only registry is not a worker dispatch")
	}

	withWorker := filepath.Join(root, "with-worker.json")
	writeFile(t, withWorker, `{"members":[{"name":"team-lead"},{"name":"ensign-gate-check"}]}`)
	if !shallowBootTeamHasWorker([]string{withWorker}) {
		t.Fatal("a non-lead team member must count as a worker dispatch")
	}
}

func TestShallowBootPhaseTasksStayMinimal(t *testing.T) {
	if claudeShallowBootGreetingTask != "" {
		t.Fatalf("Claude greeting task = %q, want no scenario input before the default greeting", claudeShallowBootGreetingTask)
	}
	if codexShallowBootGreetingTask != "Stop after the greeting." {
		t.Fatalf("Codex greeting task = %q, want exact headless stop clause", codexShallowBootGreetingTask)
	}
	if shallowBootEngageTask != "engage ." {
		t.Fatalf("engage task = %q, want exact operator input", shallowBootEngageTask)
	}
	wantArgv := []string{
		"codex", "--skip-compat-check", codexShallowBootGreetingTask, "--",
		"exec", "--json", "--enable", "multi_agent_v2",
		"--dangerously-bypass-approvals-and-sandbox",
		"--cd", "/tmp/workflow",
		"--output-last-message", "/tmp/final-message.txt",
	}
	if got := codexShallowBootFrontDoorArgv("/tmp/workflow", "/tmp/final-message.txt", codexShallowBootGreetingTask); !reflect.DeepEqual(got, wantArgv) {
		t.Fatalf("Codex shallow-boot front-door argv = %#v, want %#v", got, wantArgv)
	}

	for name, task := range map[string]string{
		"codex greeting": codexShallowBootGreetingTask,
		"engage":         shallowBootEngageTask,
	} {
		lower := strings.ToLower(task)
		for _, forbidden := range []string{"workflow", "pr", "merge", "sweep", "advance", "terminal", "archive", "gate review", "expected", "test"} {
			if strings.Contains(lower, forbidden) {
				t.Errorf("%s task %q contains behavior-directing word %q", name, task, forbidden)
			}
		}
	}
}

func exerciseShallowBootTwoPhaseTrajectory(t *testing.T) {
	good := goodShallowBootObservation()
	if err := assertShallowBoot(good); err != nil {
		t.Fatalf("correct greet-then-engage trajectory must pass: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*shallowBootObservation)
	}{
		{
			name: "greeting archives merged entity",
			mutate: func(o *shallowBootObservation) {
				o.greeting.mergedActive = ""
				o.greeting.mergedArchive = o.engage.mergedArchive
			},
		},
		{
			name: "greeting changes merged entity bytes",
			mutate: func(o *shallowBootObservation) {
				o.greeting.mergedActive += "mutation\n"
			},
		},
		{
			name: "greeting changes gate entity bytes",
			mutate: func(o *shallowBootObservation) {
				o.greeting.gateEntity += "mutation\n"
			},
		},
		{
			name: "greeting changes git head",
			mutate: func(o *shallowBootObservation) {
				o.greeting.gitHead = "different-head"
			},
		},
		{
			name: "greeting changes git porcelain",
			mutate: func(o *shallowBootObservation) {
				o.greeting.gitPorcelain = " M merged-pr.md\n"
			},
		},
		{
			name: "greeting calls gh",
			mutate: func(o *shallowBootObservation) {
				o.greeting.ghCalls = "pr view 42 --json state --jq .state\n"
			},
		},
		{
			name: "engage omits state commit",
			mutate: func(o *shallowBootObservation) {
				o.engage.gitHead = o.greeting.gitHead
			},
		},
		{
			name: "engage leaves dirty worktree",
			mutate: func(o *shallowBootObservation) {
				o.engage.gitPorcelain = " D merged-pr.md\n?? _archive/\n"
			},
		},
		{
			name: "engage misses advancement",
			mutate: func(o *shallowBootObservation) {
				o.engage.mergedActive = o.initial.mergedActive
				o.engage.mergedArchive = ""
			},
		},
		{
			name: "engage leaves nonterminal status",
			mutate: func(o *shallowBootObservation) {
				o.engage.mergedArchive = strings.Replace(o.engage.mergedArchive, "status: done", "status: implementation", 1)
			},
		},
		{
			name: "engage omits passed verdict",
			mutate: func(o *shallowBootObservation) {
				o.engage.mergedArchive = strings.Replace(o.engage.mergedArchive, "verdict: PASSED", "verdict:", 1)
			},
		},
		{
			name: "engage leaves mod block",
			mutate: func(o *shallowBootObservation) {
				o.engage.mergedArchive = strings.Replace(o.engage.mergedArchive, "mod-block:\n", "mod-block: merge:pr-merge\n", 1)
			},
		},
		{
			name: "engage changes gate entity",
			mutate: func(o *shallowBootObservation) {
				o.engage.gateEntity += "mutation\n"
			},
		},
		{
			name: "engage archives gate",
			mutate: func(o *shallowBootObservation) {
				o.engage.gateArchived = true
			},
		},
		{
			name: "engage dispatches gate worktree",
			mutate: func(o *shallowBootObservation) {
				o.engage.gateWorktreeCreated = true
			},
		},
		{
			name: "engage creates team",
			mutate: func(o *shallowBootObservation) {
				o.engage.teamWorkerOnDisk = true
			},
		},
		{
			name: "engage never calls gh",
			mutate: func(o *shallowBootObservation) {
				o.engage.ghCalls = ""
			},
		},
		{
			name: "greeting renders gate review",
			mutate: func(o *shallowBootObservation) {
				o.greetingMessage += "\nGate review: Gate Check\nDecision: approve?"
			},
		},
		{
			name: "greeting omits ready gate",
			mutate: func(o *shallowBootObservation) {
				o.greetingMessage = "A workflow is ready. Say engage to continue."
			},
		},
		{
			name: "greeting omits engage offer",
			mutate: func(o *shallowBootObservation) {
				o.greetingMessage = "Gate Check is ready."
			},
		},
		{
			name: "engage omits decision prompt",
			mutate: func(o *shallowBootObservation) {
				o.engageMessage = "Gate review: Gate Check"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			broken := goodShallowBootObservation()
			tt.mutate(&broken)
			if err := assertShallowBoot(broken); err == nil {
				t.Fatal("broken phase trajectory passed")
			}
		})
	}
}

func goodShallowBootObservation() shallowBootObservation {
	gate := shallowBootGateEntity()
	merged := shallowBootMergedEntity()
	initial := shallowBootSnapshot{
		gitHead:      "initial-head",
		gateEntity:   gate,
		mergedActive: merged,
	}
	return shallowBootObservation{
		initial:         initial,
		greeting:        initial,
		greetingMessage: "Gate Check is ready for review. Say engage to continue.",
		engage: shallowBootSnapshot{
			gitHead:       "engage-head",
			gateEntity:    gate,
			mergedArchive: "---\nid: merged-pr\nstatus: done\ncompleted: 2026-06-13T00:00:00Z\nverdict: PASSED\npr: \"#42\"\nmod-block:\nworktree:\n---\n",
			ghCalls:       "pr view 42 --json state --jq .state\n",
		},
		engageMessage: "Gate review: Gate Check — review\nDecision: approve or reject?",
	}
}
