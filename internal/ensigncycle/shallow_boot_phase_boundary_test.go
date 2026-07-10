// ABOUTME: Pins neutral shallow-boot phase tasks and the greet-read-only/engage-mutates oracle.
// ABOUTME: Synthetic trajectories isolate every forbidden early mutation and incomplete engage result.
package ensigncycle

import (
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestGitInitPersistsIdentityForIsolatedLiveStateCommits(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "fixture.txt"), "fixture\n")
	gitInit(t, root)
	for key, want := range map[string]string{
		"user.email": "t@t",
		"user.name":  "t",
	} {
		out, err := exec.Command("git", "-C", root, "config", "--local", "--get", key).CombinedOutput()
		if err != nil {
			t.Fatalf("read persisted %s: %v\n%s", key, err, out)
		}
		if got := strings.TrimSpace(string(out)); got != want {
			t.Fatalf("persisted %s = %q, want %q", key, got, want)
		}
	}

	writeFile(t, filepath.Join(root, "fixture.txt"), "changed\n")
	for _, args := range [][]string{
		{"-C", root, "add", "fixture.txt"},
		{"-C", root, "commit", "-m", "raw isolated commit"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Env = append(cleanEnviron("HOME"), "HOME="+t.TempDir())
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("raw git %v with isolated HOME: %v\n%s", args, err, out)
		}
	}
}

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
	if codexShallowBootGreetingTask != "" {
		t.Fatalf("Codex greeting task = %q, want no operator input before the default greeting", codexShallowBootGreetingTask)
	}
	if shallowBootEngageTask != "engage ." {
		t.Fatalf("engage task = %q, want exact operator input", shallowBootEngageTask)
	}
	wantArgv := []string{
		"codex", "--skip-compat-check", "--",
		"exec", "-c", `model_reasoning_effort="low"`, "--json", "--enable", "multi_agent_v2",
		"--dangerously-bypass-approvals-and-sandbox",
		"--cd", "/tmp/workflow",
		"--output-last-message", "/tmp/final-message.txt",
	}
	if got := codexShallowBootFrontDoorArgv("/tmp/workflow", "/tmp/final-message.txt", codexShallowBootGreetingTask); !reflect.DeepEqual(got, wantArgv) {
		t.Fatalf("Codex shallow-boot front-door argv = %#v, want %#v", got, wantArgv)
	}
	wantEngageArgv := []string{
		"codex", "--skip-compat-check", "engage .", "--",
		"exec", "-c", `model_reasoning_effort="low"`, "--json", "--enable", "multi_agent_v2",
		"--dangerously-bypass-approvals-and-sandbox",
		"--cd", "/tmp/workflow",
		"--output-last-message", "/tmp/final-message.txt",
	}
	if got := codexShallowBootFrontDoorArgv("/tmp/workflow", "/tmp/final-message.txt", shallowBootEngageTask); !reflect.DeepEqual(got, wantEngageArgv) {
		t.Fatalf("Codex engage front-door argv = %#v, want %#v", got, wantEngageArgv)
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

func TestShallowBootGreetingNamesReadyGateAndEngageWithoutStageToken(t *testing.T) {
	good := goodShallowBootObservation()
	message := "1 ready gate: gate-check. Use engage to continue."
	if err := assertShallowBootGreeting(good.initial, good.greeting, message); err != nil {
		t.Fatalf("ready gate plus engage must satisfy the greeting contract without a prose-only stage token: %v", err)
	}

	for name, incomplete := range map[string]string{
		"missing ready gate": "A workflow is ready. Use engage to continue.",
		"missing engage":     "1 ready gate: gate-check.",
	} {
		t.Run(name, func(t *testing.T) {
			if err := assertShallowBootGreeting(good.initial, good.greeting, incomplete); err == nil {
				t.Fatalf("incomplete greeting %q unexpectedly passed", incomplete)
			}
		})
	}
}

func exerciseShallowBootTwoPhaseTrajectory(t *testing.T) {
	good := goodShallowBootObservation()
	if err := assertShallowBoot(good); err != nil {
		t.Fatalf("correct greet-then-engage trajectory must pass: %v", err)
	}
	canonicalSlug := goodShallowBootObservation()
	canonicalSlug.greetingMessage = "gate-check is ready at review. Say engage to continue."
	if err := assertShallowBoot(canonicalSlug); err != nil {
		t.Fatalf("canonical ready-gate slug must satisfy greeting naming: %v", err)
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
			name: "greeting names only a future gate",
			mutate: func(o *shallowBootObservation) {
				o.greetingMessage = "merged-pr will reach the review gate. Say engage to continue."
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
			gitHead:       "initial-head",
			gitPorcelain:  " D merged-pr.md\n?? _archive/\n",
			gateEntity:    gate,
			mergedArchive: "---\nid: merged-pr\nstatus: done\ncompleted: 2026-06-13T00:00:00Z\nverdict: PASSED\npr: \"#42\"\nmod-block:\nworktree:\n---\n",
			ghCalls:       "pr view 42 --json state --jq .state\n",
		},
		engageMessage: "Gate review: Gate Check — review\nDecision: approve or reject?",
	}
}
