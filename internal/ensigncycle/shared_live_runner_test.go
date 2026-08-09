//go:build live

package ensigncycle

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

type liveJourneyTODO struct{ target, owner string }
type sharedRuntimeScenario struct{ name string }

func liveTODO(target, owner string) liveJourneyTODO {
	return liveJourneyTODO{target: target, owner: owner}
}

func liveJourney[Builder, Assertion any](t *testing.T, id, fixtureID string, builder Builder, todos []liveJourneyTODO, exercise func(*testing.T, liveDriver, sharedRuntimeScenario, Builder, Assertion), assertion Assertion) {
	t.Helper()
	if id == "" || fixtureID == "" || exercise == nil {
		t.Fatal("live journey metadata is incomplete")
	}
	build, buildCalls := countedLiveFunction(t, builder)
	assert, assertionCalls := countedLiveFunction(t, assertion)
	newDriver, target := liveDriverForRuntime(t)
	for _, todo := range todos {
		if todo.target == target {
			t.Skipf("TODO(%s): %s/%s lacks passing live evidence", todo.owner, target, id)
		}
	}
	if os.Getenv("SPACEDOCK_LIVE_RUNTIME") == "claude" {
		t.Parallel()
	}
	driver := newDriver()
	exercise(t, driver, sharedRuntimeScenario{name: id}, build, assert)
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

func liveDriverForRuntime(t *testing.T) (func() liveDriver, string) {
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
		return func() liveDriver { return codexAsLiveDriver{t: t, runner: newCodexLiveRunner(t)} }, "codex"
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

func TestClaudeLiveModelMapsToStableTODORole(t *testing.T) {
	tests := []struct {
		model, want string
		wantErr     bool
	}{
		{model: "sonnet", want: "claude-sonnet"},
		{model: "claude-sonnet-5", want: "claude-sonnet"},
		{model: "claude-opus-4-8", want: "claude-opus"},
		{model: "claude-future-unknown", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.model, func(t *testing.T) {
			got, err := claudeLiveRole(test.model)
			if (err != nil) != test.wantErr || got != test.want {
				t.Fatalf("claudeLiveRole(%q) = %q, %v; want %q, error=%t", test.model, got, err, test.want, test.wantErr)
			}
		})
	}
}

//spacedock:live-journey id=full-ensign-cycle fixture=realistic-lifecycle
func TestLiveCommonFullEnsignCycle(t *testing.T) {
	liveJourney(t, "full-ensign-cycle", "realistic-lifecycle", writeRealisticLifecycleFixture, nil, runFullEnsignCycleJourney, someCommitNamesOnly)
}

//spacedock:live-journey id=gate-guardrail fixture=recorded-gate/held
func TestLiveCommonGateGuardrail(t *testing.T) {
	liveJourney(t, "gate-guardrail", "recorded-gate/held", writeGateWorkflow, []liveJourneyTODO{liveTODO("codex", "3zzpdw704df1g8pg1x9thzmw"), liveTODO("pi", "3zzpdw704df1g8pg1x9thzmw")}, runGateStopScenario, assertGateHeld)
}

//spacedock:live-journey id=default-headless-gate-stop fixture=recorded-gate/pre-gate
func TestLiveCommonDefaultHeadlessGateStop(t *testing.T) {
	liveJourney(t, "default-headless-gate-stop", "recorded-gate/pre-gate", writePreGateWorkflow, []liveJourneyTODO{liveTODO("claude-sonnet", "26nk8qd48zknqnn4kc123sez"), liveTODO("codex", "26nk8qd48zknqnn4kc123sez"), liveTODO("pi", "26nk8qd48zknqnn4kc123sez")}, runGateStopScenario, assertGateHeld)
}

//spacedock:live-journey id=withdrawn-gate-recovery fixture=recorded-gate/withdrawn
func TestLiveCommonWithdrawnGateRecovery(t *testing.T) {
	liveJourney(t, "withdrawn-gate-recovery", "recorded-gate/withdrawn", writeWithdrawnGateFixture, []liveJourneyTODO{liveTODO("codex", "47gnqfm1ft6f2hcahz98m2jv")}, runClaudeWithdrawnGateRecoveryScenario, assertWithdrawnGateRecovery)
}

//spacedock:live-journey id=recorded-gate-lifecycle fixture=recorded-gate/prepared
func TestLiveCommonRecordedGateLifecycle(t *testing.T) {
	liveJourney(t, "recorded-gate-lifecycle", "recorded-gate/prepared", writeCommonPreparedRecordedGateFixture, []liveJourneyTODO{liveTODO("claude-opus", "a732sahay8wzgqrd2yr0xxr7")}, runClaudeRecordedGateLifecycleScenario, assertRecordedGateLifecycle)
}

//spacedock:live-journey id=rejection-flow fixture=rejection/before-validation-1
func TestLiveCommonRejectionFlow(t *testing.T) {
	liveJourney(t, "rejection-flow", "rejection/before-validation-1", writeRejectionWorkflow, []liveJourneyTODO{liveTODO("claude-sonnet", "zbcj98qfwtax61vxdzrf615e"), liveTODO("codex", "zbcj98qfwtax61vxdzrf615e"), liveTODO("pi", "zbcj98qfwtax61vxdzrf615e")}, runClaudeRejectionFlowScenario, assertRejectionFlow)
}

//spacedock:live-journey id=feedback-3-cycle-escalation fixture=rejection/before-validation-3
func TestLiveCommonFeedbackThreeCycleEscalation(t *testing.T) {
	liveJourney(t, "feedback-3-cycle-escalation", "rejection/before-validation-3", writeEscalationWorkflow, nil, runClaudeFeedback3CycleEscalationScenario, assertThirdCycleEscalation)
}

//spacedock:live-journey id=merge-hook-guardrail fixture=merge-hook/blocked
func TestLiveCommonMergeHookGuardrail(t *testing.T) {
	liveJourney(t, "merge-hook-guardrail", "merge-hook/blocked", writeMergeHookGuardWorkflow, nil, runClaudeMergeHookGuardrailScenario, assertMergeHookGuardHeld)
}

//spacedock:live-journey id=filing fixture=filing/empty-workflow
func TestLiveCommonFiling(t *testing.T) {
	liveJourney(t, "filing", "filing/empty-workflow", writeFilingWorkflow, nil, runClaudeFilingScenario, assertFilingCommands)
}

//spacedock:live-journey id=shallow-boot fixture=boot/held-gate
func TestLiveCommonShallowBoot(t *testing.T) {
	liveJourney(t, "shallow-boot", "boot/held-gate", writeShallowBootWorkflow, nil, runClaudeShallowBootScenario, assertShallowBoot)
}

//spacedock:live-journey id=zero-discovery fixture=boot/no-workflow
func TestLiveCommonZeroDiscovery(t *testing.T) {
	liveJourney(t, "zero-discovery", "boot/no-workflow", writeZeroDiscoveryFixture, nil, runZeroDiscoveryJourney, detectBroadSearchCommands)
}

//spacedock:live-journey id=auto-continue-after-implementation fixture=auto-continue/single-root,auto-continue/split-root
func TestLiveCommonAutoContinueAfterImplementation(t *testing.T) {
	liveJourney(t, "auto-continue-after-implementation", "auto-continue/single-root,auto-continue/split-root", autoContinueFixtureVariants, nil, runAutoContinueJourney, assertAutoContinue)
}

//spacedock:live-journey id=self-evidence-merge-triage fixture=merge-triage/unapproved-live-evidence
func TestLiveCommonSelfEvidenceMergeTriage(t *testing.T) {
	liveJourney(t, "self-evidence-merge-triage", "merge-triage/unapproved-live-evidence", writeMergeTriageWorkflow, nil, runClaudeSelfEvidenceMergeTriageScenario, assertSelfEvidenceMergeTriage)
}

//spacedock:live-journey id=smallest-sufficient-mechanism fixture=mechanism-choice/mixed-authority
func TestLiveCommonSmallestSufficientMechanism(t *testing.T) {
	liveJourney(t, "smallest-sufficient-mechanism", "mechanism-choice/mixed-authority", writeSmallestMechanismWorkflow, []liveJourneyTODO{liveTODO("claude-sonnet", "9adv48yhye5s2vkhwd7ge52d"), liveTODO("codex", "9adv48yhye5s2vkhwd7ge52d"), liveTODO("pi", "9adv48yhye5s2vkhwd7ge52d")}, runClaudeSmallestSufficientMechanismScenario, assertDurableSmallestMechanism)
}

//spacedock:live-journey id=keep-moving-posture fixture=keep-moving/mixed-events
func TestLiveCommonKeepMovingPosture(t *testing.T) {
	liveJourney(t, "keep-moving-posture", "keep-moving/mixed-events", writeKeepMovingWorkflow, []liveJourneyTODO{liveTODO("claude-sonnet", "9adv48yhye5s2vkhwd7ge52d"), liveTODO("codex", "9adv48yhye5s2vkhwd7ge52d"), liveTODO("pi", "9adv48yhye5s2vkhwd7ge52d")}, runClaudeKeepMovingScenario, assertDurableKeepMoving)
}

//spacedock:live-journey id=ac-value-reanchor fixture=ac-reanchor/means-pass-value-regressed
func TestLiveCommonACValueReanchor(t *testing.T) {
	liveJourney(t, "ac-value-reanchor", "ac-reanchor/means-pass-value-regressed", authorACReanchorScenario, nil, runACValueReanchorJourney, assertACReanchorScenario)
}

//spacedock:live-journey id=owned-conflict-owner-handoff fixture=conflict-owner/stamped-checkout
func TestLiveCommonOwnedConflictOwnerHandoff(t *testing.T) {
	liveJourney(t, "owned-conflict-owner-handoff", "conflict-owner/stamped-checkout", writeConflictOwnerFixture, []liveJourneyTODO{liveTODO("claude-sonnet", "d8qmey415fsb5q9h6q639ngf"), liveTODO("claude-opus", "d8qmey415fsb5q9h6q639ngf"), liveTODO("pi", "d8qmey415fsb5q9h6q639ngf")}, runConflictOwnerHandoffJourney, assertConflictOwnerHandoff)
}
