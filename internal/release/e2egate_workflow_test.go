package release

import (
	"fmt"
	"strings"
	"testing"
)

// TestReleaseWorkflowGatesGoreleaserOnE2E locks AC-1: the real release.yml wires
// the e2e gate so a v* tag cannot reach goreleaser without a green live e2e for
// the release commit. Two halves, both parsed from the workflow YAML (not from
// any instruction-file prose):
//   - the goreleaser-carrying job declares `needs:` including `e2e-gate`, and
//   - the e2e-gate job resolves the tagged commit SHA and runs the
//     `spacedock-release e2e-gate` step over that SHA (which itself queries
//     `gh run list --workflow "Runtime Live E2E" --status success -c <sha>`).
func TestReleaseWorkflowGatesGoreleaserOnE2E(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	if err := assertReleaseWorkflowGatesGoreleaserOnE2E(release); err != nil {
		t.Fatal(err)
	}
}

// TestReleaseWorkflowE2EGateGuardRejectsDroppedNeedsEdge is the adversarial twin
// for the needs-edge half: string-substitute the `needs: e2e-gate` edge off the
// goreleaser job and the guard must RED, because goreleaser would then cut
// without the gate ever running. A guard that stays green here is a hole.
func TestReleaseWorkflowE2EGateGuardRejectsDroppedNeedsEdge(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	if err := assertReleaseWorkflowGatesGoreleaserOnE2E(release); err != nil {
		t.Fatalf("real release.yml unexpectedly fails the e2e-gate guard before mutation: %v", err)
	}

	// The goreleaser job header in the real file carries the needs edge as its
	// first line under the job key. Drop exactly that line.
	adversarial := strings.Replace(release,
		"  goreleaser:\n    needs: e2e-gate\n",
		"  goreleaser:\n",
		1)
	if adversarial == release {
		t.Fatal("fixture workflow missing the goreleaser `needs: e2e-gate` edge to drop")
	}

	if err := assertReleaseWorkflowGatesGoreleaserOnE2E(adversarial); err == nil {
		t.Fatal("e2e-gate guard accepted a goreleaser job with the needs: e2e-gate edge dropped")
	}
}

// TestReleaseWorkflowE2EGateGuardRejectsWeakenedShaMatch is the adversarial twin
// for the SHA-binding half: weaken the e2e-gate step so it no longer binds the
// query to the tagged commit SHA (drop the commit-resolving `git rev-list` and
// the `-c "$RELEASE_COMMIT"` binding, accepting "some green run somewhere"). The
// guard must RED — a gate that accepts a green run on any commit does not prove
// the live matrix passed for the commit being released.
func TestReleaseWorkflowE2EGateGuardRejectsWeakenedShaMatch(t *testing.T) {
	release := readWorkflow(t, "release.yml")

	for _, tc := range []struct {
		name string
		from string
		to   string
	}{
		{
			"drop the commit-resolving rev-list",
			`RELEASE_COMMIT="$(git rev-list -1 "$GITHUB_REF_NAME")"`,
			`RELEASE_COMMIT=""`,
		},
		{
			"drop the SHA-bound gate invocation",
			`go run ./cmd/spacedock-release e2e-gate "$RELEASE_COMMIT"`,
			`echo "skipping e2e gate"`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			adversarial := strings.Replace(release, tc.from, tc.to, 1)
			if adversarial == release {
				t.Fatalf("fixture workflow missing the SHA-match anchor %q to weaken", tc.from)
			}
			if err := assertReleaseWorkflowGatesGoreleaserOnE2E(adversarial); err == nil {
				t.Fatalf("e2e-gate guard accepted a workflow with %s", tc.name)
			}
		})
	}
}

// assertReleaseWorkflowGatesGoreleaserOnE2E binds the AC-1 gate to the parsed
// job graph + the e2e-gate step's resolved commands. It requires:
//   - the goreleaser-action carrier needs the e2e-gate job, and
//   - some job resolves the tagged commit SHA (git rev-list -1 $GITHUB_REF_NAME)
//     and runs `spacedock-release e2e-gate "$RELEASE_COMMIT"` over it.
//
// The predicate's own SHA/conclusion matching is unit-tested separately; this
// guard proves the WIRING gates the cut on that predicate, bound to the tag.
func assertReleaseWorkflowGatesGoreleaserOnE2E(workflow string) error {
	jobs := parseWorkflowJobs(workflow)

	var goreleaserCarriers []workflowJob
	gateJobNames := map[string]bool{}
	for _, job := range jobs {
		carriesGoreleaser := false
		resolvesCommit, invokesGate := false, false
		for _, step := range job.steps {
			if strings.HasPrefix(step.uses, "goreleaser/goreleaser-action@") {
				carriesGoreleaser = true
			}
			for _, command := range executableShellCommands(step.run) {
				if strings.Contains(command, `git rev-list -1 "$GITHUB_REF_NAME"`) {
					resolvesCommit = true
				}
				if isE2EGateInvocation(command) {
					invokesGate = true
				}
			}
		}
		if carriesGoreleaser {
			goreleaserCarriers = append(goreleaserCarriers, job)
		}
		if resolvesCommit && invokesGate {
			gateJobNames[job.name] = true
		}
	}

	if len(goreleaserCarriers) == 0 {
		return fmt.Errorf("release.yml has no job carrying the goreleaser action")
	}
	if len(gateJobNames) == 0 {
		return fmt.Errorf("release.yml has no job that resolves the tagged commit SHA and runs `spacedock-release e2e-gate \"$RELEASE_COMMIT\"` over it")
	}

	for _, carrier := range goreleaserCarriers {
		needsAGate := false
		for _, need := range carrier.needs {
			if gateJobNames[need] {
				needsAGate = true
			}
		}
		if !needsAGate {
			return fmt.Errorf("release.yml goreleaser job %q does not declare needs: on the e2e-gate job — a v* tag could reach goreleaser without a green live e2e for the release commit", carrier.name)
		}
	}
	return nil
}

// isE2EGateInvocation reports whether command runs the SHA-bound e2e-gate
// subcommand: it must invoke `spacedock-release e2e-gate` AND pass the resolved
// release commit ($RELEASE_COMMIT), so a gate that drops the SHA binding (and
// would accept a green run on any commit) is not recognized.
func isE2EGateInvocation(command string) bool {
	return strings.Contains(command, `go run ./cmd/spacedock-release e2e-gate `) &&
		strings.Contains(command, `"$RELEASE_COMMIT"`)
}
