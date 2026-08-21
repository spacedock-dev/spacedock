//go:build live

package ensigncycle

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

type liveJourneyGap struct{ kind, target, owner string }
type sharedRuntimeScenario struct {
	name  string
	gap   liveJourneyGap
	grade liveGrade
}

func liveTODO(target, owner string) liveJourneyGap {
	return liveJourneyGap{kind: "todo", target: target, owner: owner}
}

func liveXFail(target, owner string) liveJourneyGap {
	return liveJourneyGap{kind: "xfail", target: target, owner: owner}
}

func liveJourney[Builder, Assertion any](t *testing.T, id, fixtureID string, builder Builder, gaps []liveJourneyGap, exercise func(*testing.T, liveDriver, sharedRuntimeScenario, Builder, Assertion), assertion Assertion) {
	t.Helper()
	if id == "" || fixtureID == "" || exercise == nil {
		t.Fatal("live journey metadata is incomplete")
	}
	build, buildCalls := countedLiveFunction(t, builder)
	assert, assertionCalls := countedLiveFunction(t, assertion)
	newDriver, target := liveDriverForRuntime(t, id)
	var selected liveJourneyGap
	for _, gap := range gaps {
		if gap.target == target {
			selected = gap
			if gap.kind == "todo" {
				t.Skipf("TODO(%s): %s/%s lacks passing live evidence", gap.owner, target, id)
			}
		}
	}
	if liveRuntimeRunsParallel(os.Getenv("SPACEDOCK_LIVE_RUNTIME")) {
		t.Parallel()
	}
	driver := newDriver()
	exercise(t, driver, sharedRuntimeScenario{name: id, gap: selected}, build, assert)
	if *buildCalls == 0 || *assertionCalls == 0 {
		t.Fatalf("live journey %s executed builder/assertion %d/%d times", id, *buildCalls, *assertionCalls)
	}
}

func countedLiveFunction[Function any](t *testing.T, function Function) (Function, *int) {
	calls, rv := 0, reflect.ValueOf(function)
	if !rv.IsValid() || rv.Kind() != reflect.Func || rv.IsNil() {
		t.Fatal("live journey builder/assertion is not a function")
	}
	return reflect.MakeFunc(rv.Type(), func(args []reflect.Value) []reflect.Value {
		calls++
		return rv.Call(args)
	}).Interface().(Function), &calls
}

func noLiveGrade(liveResult) {}

func liveRuntimeRunsParallel(runtime string) bool {
	return runtime == "claude" || runtime == "codex" || runtime == "pi"
}

func liveDriverForRuntime(t *testing.T, id string) (func() liveDriver, string) {
	t.Helper()
	switch runtime := os.Getenv("SPACEDOCK_LIVE_RUNTIME"); runtime {
	case "claude":
		role, err := claudeLiveRole(envOr("SPACEDOCK_LIVE_MODEL", "sonnet"))
		if err != nil {
			t.Fatal(err)
			return nil, ""
		}
		return func() liveDriver { return newClaudeLiveRunner(t) }, role
	case "codex":
		return func() liveDriver { return codexAsLiveDriver{t: t, runner: newCodexLiveRunner(t, id)} }, "codex"
	case "pi":
		return func() liveDriver { return newPiSharedLiveDriver(t) }, "pi"
	default:
		t.Fatalf("SPACEDOCK_LIVE_RUNTIME=%q, want claude, codex, or pi", runtime)
		return nil, ""
	}
}

func claudeLiveRole(model string) (string, error) {
	switch model {
	case "sonnet", "claude-sonnet-5":
		return "claude-sonnet", nil
	case "claude-opus-4-8":
		return "claude-opus", nil
	default:
		return "", fmt.Errorf("SPACEDOCK_LIVE_MODEL=%q, want sonnet, claude-sonnet-5, or claude-opus-4-8", model)
	}
}

//spacedock:live-journey id=owned-conflict-owner-handoff fixture=conflict-owner/stamped-checkout
func TestLiveCommonOwnedConflictOwnerHandoff(t *testing.T) {
	liveJourney(t, "owned-conflict-owner-handoff", "conflict-owner/stamped-checkout", writeConflictOwnerFixture, []liveJourneyGap{liveXFail("pi", "fe7bfjz9sb8wyckmnnm3ncjx")}, runConflictOwnerHandoffJourney, assertConflictOwnerHandoff)
}

// The pi XFAIL below is registered because the Pi rejection-flow target times out
// before completing the expected stop (registered owner p17swb3375rt525fn7f8xt7e).
// Pi-dialect extractors now exist (piRecordedRejectionRound, etc.) so the XFAIL
// can XPASS once the timeout repair lets the FO complete the two-validation cycle.
//
//spacedock:live-journey id=rejection-flow fixture=rejection/before-validation-1
func TestLiveCommonRejectionFlow(t *testing.T) {
	liveJourney(t, "rejection-flow", "rejection/before-validation-1", writeRejectionWorkflow, []liveJourneyGap{liveXFail("pi", "p17swb3375rt525fn7f8xt7e")}, runClaudeRejectionFlowScenario, assertRejectionFlow)
}

//spacedock:live-journey id=auto-continue-after-implementation fixture=auto-continue/single-root,auto-continue/split-root
func TestLiveCommonAutoContinueAfterImplementation(t *testing.T) {
	liveJourney(t, "auto-continue-after-implementation", "auto-continue/single-root,auto-continue/split-root", autoContinueFixtureVariants, nil, runAutoContinueJourney, assertAutoContinue)
}

//spacedock:live-journey id=keep-moving-posture fixture=keep-moving/mixed-events
func TestLiveCommonKeepMovingPosture(t *testing.T) {
	liveJourney(t, "keep-moving-posture", "keep-moving/mixed-events", writeKeepMovingWorkflow, []liveJourneyGap{liveXFail("claude-sonnet", "060xp69y61yhrww23g3wvwqy"), liveXFail("pi", "x02375wsg6q61xek7p0t36j2")}, runClaudeKeepMovingScenario, assertDurableKeepMoving)
}

//spacedock:live-journey id=default-headless-gate-stop fixture=recorded-gate/pre-gate
func TestLiveCommonDefaultHeadlessGateStop(t *testing.T) {
	liveJourney(t, "default-headless-gate-stop", "recorded-gate/pre-gate", writePreGateWorkflow, nil, runGateStopScenario, assertGateHeld)
}

// AUDIT(2026-08-16): finding 5 — the direct-commit marker matches ANY git commit and
// no assertion checks the strategy doc landed. The pi XFAIL below is structural
// (finding 11): a claude-dialect trace reads the pi stream, so it can never XPASS.
//
//spacedock:live-journey id=smallest-sufficient-mechanism fixture=mechanism-choice/mixed-authority
func TestLiveCommonSmallestSufficientMechanism(t *testing.T) {
	liveJourney(t, "smallest-sufficient-mechanism", "mechanism-choice/mixed-authority", writeSmallestMechanismWorkflow, []liveJourneyGap{liveXFail("pi", "h30c9jrfcf21fdh2qs5z58sd")}, runClaudeSmallestSufficientMechanismScenario, assertDurableSmallestMechanism)
}

//spacedock:live-journey id=recorded-gate-lifecycle fixture=recorded-gate/prepared
func TestLiveCommonRecordedGateLifecycle(t *testing.T) {
	liveJourney(t, "recorded-gate-lifecycle", "recorded-gate/prepared", writeCommonPreparedRecordedGateFixture, nil, runClaudeRecordedGateLifecycleScenario, assertRecordedGateLifecycle)
}

// AUDIT(2026-08-16): finding 4 — token-and-one-commit grade: a stage-skipping run
// with a hand-written report grades green; integrationTransitionCommitted, the one
// field that distinguishes the full loop from skip-to-end, is never used here.
//
//spacedock:live-journey id=full-ensign-cycle fixture=realistic-lifecycle
func TestLiveCommonFullEnsignCycle(t *testing.T) {
	liveJourney(t, "full-ensign-cycle", "realistic-lifecycle", writeRealisticLifecycleFixture, nil, runFullEnsignCycleJourney, someCommitNamesOnly)
}

//spacedock:live-journey id=ac-value-reanchor fixture=ac-reanchor/means-pass-value-regressed
func TestLiveCommonACValueReanchor(t *testing.T) {
	liveJourney(t, "ac-value-reanchor", "ac-reanchor/means-pass-value-regressed", authorACReanchorScenario, nil, runACValueReanchorJourney, assertACReanchorScenario)
}

//spacedock:live-journey id=withdrawn-gate-recovery fixture=recorded-gate/withdrawn
func TestLiveCommonWithdrawnGateRecovery(t *testing.T) {
	liveJourney(t, "withdrawn-gate-recovery", "recorded-gate/withdrawn", writeWithdrawnGateFixture, nil, runClaudeWithdrawnGateRecoveryScenario, assertWithdrawnGateRecovery)
}

//spacedock:live-journey id=gate-guardrail fixture=recorded-gate/held
func TestLiveCommonGateGuardrail(t *testing.T) {
	liveJourney(t, "gate-guardrail", "recorded-gate/held", writeGateWorkflow, nil, runGateStopScenario, assertGateHeld)
}

//spacedock:live-journey id=feedback-3-cycle-escalation fixture=rejection/before-validation-3
func TestLiveCommonFeedbackThreeCycleEscalation(t *testing.T) {
	liveJourney(t, "feedback-3-cycle-escalation", "rejection/before-validation-3", writeEscalationWorkflow, nil, runClaudeFeedback3CycleEscalationScenario, assertThirdCycleEscalation)
}

// AUDIT(2026-08-16): finding 13 — the diagnosis check is attribution-blind: a final
// message naming TestZeroDiscover while concluding the inherited label grades green.
//
//spacedock:live-journey id=self-evidence-merge-triage fixture=merge-triage/unapproved-live-evidence
func TestLiveCommonSelfEvidenceMergeTriage(t *testing.T) {
	liveJourney(t, "self-evidence-merge-triage", "merge-triage/unapproved-live-evidence", writeMergeTriageWorkflow, nil, runClaudeSelfEvidenceMergeTriageScenario, assertSelfEvidenceMergeTriage)
}

// AUDIT(2026-08-16): finding 7 — the final-response requirement is graded against
// finalMessage+stream, where the refused command's tool-result echo satisfies it.
//
//spacedock:live-journey id=merge-hook-guardrail fixture=merge-hook/blocked
func TestLiveCommonMergeHookGuardrail(t *testing.T) {
	liveJourney(t, "merge-hook-guardrail", "merge-hook/blocked", writeMergeHookGuardWorkflow, nil, runClaudeMergeHookGuardrailScenario, assertMergeHookGuardHeld)
}

//spacedock:live-journey id=filing fixture=filing/empty-workflow
func TestLiveCommonFiling(t *testing.T) {
	liveJourney(t, "filing", "filing/empty-workflow", writeFilingWorkflow, nil, runClaudeFilingScenario, assertFilingCommands)
}

// AUDIT(2026-08-16): finding 6 — grades Bash commands only; Glob/Grep broad search
// at boot is detected only by the non-grading failure diagnostic.
//
//spacedock:live-journey id=zero-discovery fixture=boot/no-workflow
func TestLiveCommonZeroDiscovery(t *testing.T) {
	liveJourney(t, "zero-discovery", "boot/no-workflow", writeZeroDiscoveryFixture, nil, runZeroDiscoveryJourney, detectBroadSearchCommands)
}

//spacedock:live-journey id=shallow-boot fixture=boot/held-gate
func TestLiveCommonShallowBoot(t *testing.T) {
	liveJourney(t, "shallow-boot", "boot/held-gate", writeShallowBootWorkflow, nil, runClaudeShallowBootScenario, assertShallowBoot)
}
