package ensigncycle

import (
	"fmt"
	"testing"
)

type mergeModRecoveryObservation struct {
	activeAfter   string
	archivedAfter string
	gitClean      bool
}

func assertMergeModRecovery(o mergeModRecoveryObservation) error {
	if o.activeAfter != "" {
		return fmt.Errorf("merge-recovery remains active instead of archived")
	}
	if o.archivedAfter == "" {
		return fmt.Errorf("merge-recovery has no archived durable record")
	}
	if !mergedTerminalStatus.MatchString(o.archivedAfter) {
		return fmt.Errorf("archived merge-recovery is not status: done")
	}
	if !mergedVerdictPassed.MatchString(o.archivedAfter) {
		return fmt.Errorf("archived merge-recovery is not verdict: PASSED")
	}
	if !mergedModBlockClear.MatchString(o.archivedAfter) {
		return fmt.Errorf("archived merge-recovery retains its mod-block")
	}
	if !o.gitClean {
		return fmt.Errorf("merge-recovery left the workflow git worktree dirty")
	}
	return nil
}

func TestMergeModRecoveryDurableOracle(t *testing.T) {
	good := mergeModRecoveryObservation{
		archivedAfter: "---\nstatus: done\nverdict: PASSED\nmod-block:\n---\n",
		gitClean:      true,
	}
	if err := assertMergeModRecovery(good); err != nil {
		t.Fatalf("positive recovery: %v", err)
	}
	controls := []mergeModRecoveryObservation{
		{activeAfter: "status: implementation", archivedAfter: good.archivedAfter, gitClean: true},
		{gitClean: true},
		{archivedAfter: "status: implementation\nverdict: PASSED\nmod-block:\n", gitClean: true},
		{archivedAfter: "status: done\nverdict: PASSED\nmod-block: merge:pr-merge\n", gitClean: true},
		{archivedAfter: good.archivedAfter, gitClean: false},
	}
	for i, control := range controls {
		if err := assertMergeModRecovery(control); err == nil {
			t.Errorf("broken recovery control %d passed", i)
		}
	}
}
