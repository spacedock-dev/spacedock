// ABOUTME: Pins the initial First Officer worker-spawn and completion boundary.
// ABOUTME: Rejects ceremonial substitutes that previously let validation start.
package contractlint

import (
	"strings"
	"testing"
)

const initialWorkerSpawnGuard = "On exit 0, the next host action MUST be `«worker.spawn»` with every helper-emitted field unchanged."

const initialWorkerHandleGuard = "Record the returned handle before narration, a file edit, status change, report read, gate action, or wait."

const initialWorkerCompletionGuard = "Do not advance to validation until `«completion-signal»` arrives and the entity-file stage report passes the completion gate."

const initialWorkerFalseEvidenceGuard = "A successful dispatch build, narration, direct status change, or self-authored report is not worker evidence."

const initialWorkerEmptyWaitGuard = "An empty wait without a completion signal is not worker evidence."

const initialWorkerChecklistGuard = "Build a numbered checklist of one to three dispatch-specific linchpin signals"

const initialWorkerChecklistFallback = "When neither source supplies an item, use the target stage's declared requirement as one linchpin; do not pad."

const initialWorkerStartedTransition = "status={next_stage} started"

func TestInitialWorkerSpawnGuardPrecedesCompletionAndValidation(t *testing.T) {
	body := readRepoFile(t, "skills/first-officer/references/fo-dispatch-core.md")

	spawn := strings.Index(body, initialWorkerSpawnGuard)
	handle := strings.Index(body, initialWorkerHandleGuard)
	completion := strings.Index(body, initialWorkerCompletionGuard)
	falseEvidence := strings.Index(body, initialWorkerFalseEvidenceGuard)
	emptyWait := strings.Index(body, initialWorkerEmptyWaitGuard)
	if spawn < 0 {
		t.Fatalf("initial dispatch contract missing real spawn guard %q", initialWorkerSpawnGuard)
	}
	if handle < 0 {
		t.Fatalf("initial dispatch contract missing returned-handle guard %q", initialWorkerHandleGuard)
	}
	if completion < 0 {
		t.Fatalf("initial dispatch contract missing completion-before-validation guard %q", initialWorkerCompletionGuard)
	}
	if falseEvidence < 0 {
		t.Fatalf("initial dispatch contract missing false-evidence guard %q", initialWorkerFalseEvidenceGuard)
	}
	if emptyWait < 0 {
		t.Fatalf("initial dispatch contract missing empty-wait guard %q", initialWorkerEmptyWaitGuard)
	}
	if !(spawn < handle && handle < completion && completion < falseEvidence && falseEvidence < emptyWait) {
		t.Fatalf("initial dispatch guard order = spawn:%d handle:%d completion:%d false-evidence:%d empty-wait:%d", spawn, handle, completion, falseEvidence, emptyWait)
	}
}

func TestInitialWorkerDispatchAlwaysBuildsANonemptyChecklist(t *testing.T) {
	body := readRepoFile(t, "skills/first-officer/references/fo-dispatch-core.md")
	for _, want := range []string{initialWorkerChecklistGuard, initialWorkerChecklistFallback} {
		if !strings.Contains(body, want) {
			t.Errorf("initial dispatch contract missing nonempty-checklist rule %q", want)
		}
	}
}

func TestInitialWorkerDispatchStampsBothStageTransitionsStarted(t *testing.T) {
	body := readRepoFile(t, "skills/first-officer/references/fo-dispatch-core.md")
	if got := strings.Count(body, initialWorkerStartedTransition); got != 2 {
		t.Fatalf("initial dispatch contract has %d started stage transitions, want 2", got)
	}
}
