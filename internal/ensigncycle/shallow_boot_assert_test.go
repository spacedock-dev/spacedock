package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// shallowBootSnapshot is one durable observation of the two-phase shallow-boot
// fixture. Transcript prose never stands in for a state transition: the oracle
// reads Git, both entity locations, the recording gh shim, and dispatch traces.
type shallowBootSnapshot struct {
	gitHead             string
	gitPorcelain        string
	gateEntity          string
	mergedActive        string
	mergedArchive       string
	ghCalls             string
	gateArchived        bool
	gateWorktreeCreated bool
	teamWorkerOnDisk    bool
}

type shallowBootObservation struct {
	initial         shallowBootSnapshot
	greeting        shallowBootSnapshot
	greetingMessage string
	engage          shallowBootSnapshot
	engageMessage   string
}

var (
	mergedTerminalStatus = regexp.MustCompile(`(?im)^status:\s*done\s*$`)
	mergedVerdictPassed  = regexp.MustCompile(`(?im)^verdict:\s*PASSED\s*$`)
	mergedModBlockClear  = regexp.MustCompile(`(?im)^mod-block:\s*$`)
)

// assertShallowBoot grades the #480 phase boundary: greeting is a byte-for-byte,
// no-gh observation; first engage performs the PR sweep, archives a complete
// terminal entity, and renders the already-ready gate without mutating it.
func assertShallowBoot(o shallowBootObservation) error {
	if err := assertShallowBootGreeting(o.initial, o.greeting, o.greetingMessage); err != nil {
		return err
	}
	return assertShallowBootEngage(o.initial, o.greeting, o.engage, o.engageMessage)
}

func assertShallowBootGreeting(initial, greeting shallowBootSnapshot, message string) error {
	if greeting.gitHead != initial.gitHead {
		return fmt.Errorf("Git HEAD changed during the greeting")
	}
	if greeting.gitPorcelain != initial.gitPorcelain {
		return fmt.Errorf("Git porcelain changed during the greeting")
	}
	if greeting.gateEntity != initial.gateEntity {
		return fmt.Errorf("the gated entity changed during the greeting")
	}
	if greeting.mergedActive != initial.mergedActive {
		return fmt.Errorf("the active merged-PR entity changed during the greeting")
	}
	if greeting.mergedArchive != initial.mergedArchive {
		return fmt.Errorf("the merged-PR archive changed during the greeting")
	}
	if greeting.ghCalls != initial.ghCalls {
		return fmt.Errorf("the greeting called gh; PR probing belongs to first engage")
	}
	if greeting.gateArchived != initial.gateArchived || greeting.gateWorktreeCreated != initial.gateWorktreeCreated {
		return fmt.Errorf("the greeting archived or dispatched the gated entity")
	}
	if greeting.teamWorkerOnDisk != initial.teamWorkerOnDisk {
		return fmt.Errorf("the greeting created a worker team")
	}

	lower := strings.ToLower(message)
	if !strings.Contains(lower, "gate check") || !strings.Contains(lower, "review") || !strings.Contains(lower, "engage") {
		return fmt.Errorf("the greeting did not name the ready Gate Check review and offer engage")
	}
	if strings.Contains(lower, "gate review:") || strings.Contains(lower, "decision:") {
		return fmt.Errorf("the greeting rendered the gate review before engage")
	}
	return nil
}

func assertShallowBootEngage(initial, greeting, engage shallowBootSnapshot, message string) error {
	if engage.gitHead == greeting.gitHead {
		return fmt.Errorf("first engage did not commit the merged-PR state transition")
	}
	if engage.gitPorcelain != initial.gitPorcelain {
		return fmt.Errorf("Git porcelain is not clean after first engage")
	}
	if engage.gateEntity != initial.gateEntity {
		return fmt.Errorf("the gated entity changed during engage — it was dispatched or self-resolved")
	}
	if !reviewStatus.MatchString(engage.gateEntity) {
		return fmt.Errorf("the gated entity is no longer at status: review")
	}
	if completedSet.MatchString(engage.gateEntity) || verdictSetFM.MatchString(engage.gateEntity) {
		return fmt.Errorf("the gated entity has completed/verdict set — engage self-approved it")
	}
	if engage.gateArchived {
		return fmt.Errorf("the gated entity was archived during engage")
	}
	if engage.gateWorktreeCreated {
		return fmt.Errorf("a worktree was created for the gated entity — engage dispatched it")
	}
	if engage.teamWorkerOnDisk {
		return fmt.Errorf("a worker member exists in the team registry — engage dispatched during a gate-only trajectory")
	}

	if engage.mergedActive != "" {
		return fmt.Errorf("the merged-PR entity remains active — first engage did not archive it")
	}
	if engage.mergedArchive == "" {
		return fmt.Errorf("the merged-PR entity was not archived on first engage")
	}
	if !mergedTerminalStatus.MatchString(engage.mergedArchive) {
		return fmt.Errorf("the archived merged-PR entity is not at status: done")
	}
	if !mergedVerdictPassed.MatchString(engage.mergedArchive) {
		return fmt.Errorf("the archived merged-PR entity has no verdict: PASSED")
	}
	if !mergedModBlockClear.MatchString(engage.mergedArchive) {
		return fmt.Errorf("the archived merged-PR entity still carries a mod-block")
	}
	if engage.ghCalls == greeting.ghCalls {
		return fmt.Errorf("first engage made no gh call")
	}

	lower := strings.ToLower(message)
	if !strings.Contains(lower, "gate review:") || !strings.Contains(lower, "decision:") {
		return fmt.Errorf("the engage response did not render a gate review and decision prompt")
	}
	return nil
}

func gatherShallowBootSnapshot(t *testing.T, workflowRoot, teamRoot string, fx shallowBootFixture) shallowBootSnapshot {
	t.Helper()
	snapshot := shallowBootSnapshot{
		gateEntity:    readFileAllowMissing(fx.gateEntityPath),
		mergedActive:  readFileAllowMissing(fx.mergedEntityPath),
		mergedArchive: readFileAllowMissing(fx.mergedArchive),
		ghCalls:       readFileAllowMissing(fx.stubGhLog),
	}

	head, err := exec.Command("git", "-C", workflowRoot, "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		t.Fatalf("snapshot shallow-boot Git HEAD: %v\n%s", err, head)
	}
	snapshot.gitHead = strings.TrimSpace(string(head))
	porcelain, err := exec.Command("git", "-C", workflowRoot, "status", "--porcelain=v1").CombinedOutput()
	if err != nil {
		t.Fatalf("snapshot shallow-boot Git porcelain: %v\n%s", err, porcelain)
	}
	snapshot.gitPorcelain = string(porcelain)

	if _, err := os.Stat(fx.gateEntityArchivePath(workflowRoot)); err == nil {
		snapshot.gateArchived = true
	}
	if matches, _ := filepath.Glob(filepath.Join(workflowRoot, ".worktrees", "*gate-check*")); len(matches) > 0 {
		snapshot.gateWorktreeCreated = true
	}
	if teamRoot != "" {
		if matches, _ := filepath.Glob(filepath.Join(teamRoot, "*", "config.json")); shallowBootTeamHasWorker(matches) {
			snapshot.teamWorkerOnDisk = true
		}
	}
	return snapshot
}

// shallowBootTeamHasWorker ignores Claude's transport-created, leader-only
// session registry and detects an actual dispatched teammate. Recent interactive
// Claude versions create the lead record before the model's greeting even though
// no TeamCreate/Agent action occurred.
func shallowBootTeamHasWorker(configPaths []string) bool {
	for _, path := range configPaths {
		body, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var config struct {
			Members []struct {
				Name string `json:"name"`
			} `json:"members"`
		}
		if json.Unmarshal(body, &config) != nil {
			continue
		}
		for _, member := range config.Members {
			if member.Name != "" && member.Name != "team-lead" {
				return true
			}
		}
	}
	return false
}

func (fx shallowBootFixture) gateEntityArchivePath(workflowRoot string) string {
	return filepath.Join(workflowRoot, "_archive", "gate-check.md")
}

func readFileAllowMissing(path string) string {
	if path == "" {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}
