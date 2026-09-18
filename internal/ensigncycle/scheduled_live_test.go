//go:build live

package ensigncycle

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func TestLiveScheduled(t *testing.T) {
	runtime := os.Getenv("SPACEDOCK_LIVE_RUNTIME")
	if runtime != "claude" && runtime != "codex" {
		t.Skip("scheduled suite requires claude or codex")
	}
	// Rounded whole-test seconds from run 34996910090; priorities, not budgets.
	rows := []struct {
		name          string
		codex, claude int
		run           func(*testing.T)
	}{
		{"TestLiveCommonACValueReanchor", 100, 360, TestLiveCommonACValueReanchor},
		{"TestLiveCommonAutoContinueAfterImplementation", 340, 440, TestLiveCommonAutoContinueAfterImplementation},
		{"TestLiveCommonDefaultHeadlessGateStop", 280, 240, TestLiveCommonDefaultHeadlessGateStop},
		{"TestLiveCommonFeedbackThreeCycleEscalation", 70, 370, TestLiveCommonFeedbackThreeCycleEscalation},
		{"TestLiveCommonFiling", 40, 70, TestLiveCommonFiling},
		{"TestLiveCommonFullEnsignCycle", 110, 170, TestLiveCommonFullEnsignCycle},
		{"TestLiveCommonGateGuardrail", 90, 90, TestLiveCommonGateGuardrail},
		{"TestLiveCommonKeepMovingPosture", 270, 340, TestLiveCommonKeepMovingPosture},
		{"TestLiveCommonMergeHookGuardrail", 50, 30, TestLiveCommonMergeHookGuardrail},
		{"TestLiveCommonOwnedConflictOwnerHandoff", 130, 300, TestLiveCommonOwnedConflictOwnerHandoff},
		{"TestLiveCommonRecordedGateLifecycle", 160, 190, TestLiveCommonRecordedGateLifecycle},
		{"TestLiveCommonRejectionFlow", 340, 430, TestLiveCommonRejectionFlow},
		{"TestLiveCommonSelfEvidenceMergeTriage", 70, 120, TestLiveCommonSelfEvidenceMergeTriage},
		{"TestLiveCommonShallowBoot", 20, 30, TestLiveCommonShallowBoot},
		{"TestLiveCommonSmallestSufficientMechanism", 240, 250, TestLiveCommonSmallestSufficientMechanism},
		{"TestLiveCommonWithdrawnGateRecovery", 90, 100, TestLiveCommonWithdrawnGateRecovery},
		{"TestLiveCommonZeroDiscovery", 30, 20, TestLiveCommonZeroDiscovery},
	}
	jobs := make([]liveScheduledTest, 0, 20)
	for _, row := range rows {
		hint := row.codex
		if runtime == "claude" {
			hint = row.claude
		}
		jobs = append(jobs, liveScheduledTest{row.name, hint, row.run})
	}
	if runtime == "claude" {
		jobs = append(jobs,
			liveScheduledTest{"TestLiveBareReachable", 110, TestLiveBareReachable},
			liveScheduledTest{"TestLiveBreakGlassShimRecovery", 270, TestLiveBreakGlassShimRecovery},
			liveScheduledTest{"TestLiveMergedTeamModeDispatch", 140, TestLiveMergedTeamModeDispatch},
		)
	}
	if runtime == "codex" {
		jobs = append(jobs, liveScheduledTest{"TestLiveSemanticNamesCodex", 110, TestLiveSemanticNamesCodex})
	}
	orderLiveTests(jobs)
	var mu sync.Mutex
	next := 0
	for slot := 0; slot < 3; slot++ {
		t.Run(fmt.Sprintf("slot-%d", slot), func(t *testing.T) {
			t.Parallel()
			for {
				mu.Lock()
				if next == len(jobs) {
					mu.Unlock()
					return
				}
				job := jobs[next]
				next++
				mu.Unlock()
				t.Run(job.name, job.run)
			}
		})
	}
}
