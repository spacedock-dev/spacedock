// ABOUTME: An initial-stage gate briefing pins the exact committed seed bytes,
// ABOUTME: so an unchanged seed replays one digest and an edited seed moves it.
package gates

import (
	"os"
	"strings"
	"testing"
)

// TestPrepareFromInitialSeedPinsCommittedBytes is AC-4. prepareFixture's stage
// is initial+gate over an entity with no stage report, so the seed itself is the
// artifact. Preparing an unchanged seed twice replays one identical digest under
// one identical briefing ID; dropping the artifact pin lets the replay match a
// seed it no longer describes. Once the seed moves, the pin refuses to rebind
// the open attempt at all, and a fresh attempt over the edited seed digests
// differently.
func TestPrepareFromInitialSeedPinsCommittedBytes(t *testing.T) {
	workflow, state, entity, _, _ := prepareFixture(t, "flat")
	input := PrepareInput{WorkflowDir: workflow, Question: "Approve this seed?", Artifact: entity, Summary: "Seed review"}

	first, err := Prepare(entity, input)
	if err != nil {
		t.Fatalf("prepare from the seed: %v", err)
	}
	replay, err := Prepare(entity, input)
	if err != nil {
		t.Fatalf("re-prepare the unchanged seed: %v", err)
	}
	if replay.Briefing != first.Briefing || replay.Digest != first.Digest {
		t.Fatalf("unchanged seed gave %s/%s, want a replay of %s/%s", replay.Briefing, replay.Digest, first.Briefing, first.Digest)
	}

	body, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(entity, append(body, []byte("\nA later edit to the seed.\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	prepareGitRun(t, state, "add", ".")
	prepareGitRun(t, state, "commit", "-q", "-m", "edit the seed")
	// Same attempt, same briefing ID, different seed bytes: the binding is pinned
	// to the bytes it was cut over, so it refuses rather than silently redigesting.
	if _, err := Prepare(entity, input); err == nil || !strings.Contains(err.Error(), "frozen and cannot be rebound") {
		t.Fatalf("re-prepare over an edited seed error = %v, want the frozen-binding refusal", err)
	}
	if _, err := Withdraw(entity, WithdrawInput{Reason: "seed edited", WorkflowDir: workflow}); err != nil {
		t.Fatalf("withdraw the open attempt: %v", err)
	}
	prepareGitRun(t, state, "add", ".")
	prepareGitRun(t, state, "commit", "-q", "-m", "record the withdrawal")
	moved, err := Prepare(entity, input)
	if err != nil {
		t.Fatalf("prepare the edited seed on a fresh attempt: %v", err)
	}
	if moved.Digest == first.Digest {
		t.Fatalf("digest %s survived a modified seed: the artifact is not pinned", moved.Digest)
	}
}
