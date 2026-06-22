package release

import (
	"fmt"
	"strings"
	"testing"
)

// TestReleaseWorkflowGatesGoreleaserOnManifestTag locks the wiring half of AC-2:
// the real release.yml must run `spacedock-release manifest-tag-gate` over the
// tagged manifests in a job the goreleaser carrier `needs:`, so a tag whose semver
// disagrees with the tagged commit's plugin.json is BLOCKED before goreleaser
// fires — not left to an advisory manual step. Both halves are parsed from the
// workflow YAML (not from instruction prose):
//   - some job runs the SHA-free manifest-tag-gate step over both plugin manifests
//     with the tag ($GITHUB_REF_NAME), AND
//   - the goreleaser-carrying job `needs:` that gate job.
func TestReleaseWorkflowGatesGoreleaserOnManifestTag(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	if err := assertReleaseWorkflowGatesGoreleaserOnManifestTag(release); err != nil {
		t.Fatal(err)
	}
}

// TestReleaseWorkflowManifestTagGuardRejectsDroppedStep is the adversarial twin
// for the gate-step half: remove the manifest-tag-gate invocation and the guard
// must RED, because goreleaser would then cut without the manifest/tag check ever
// running (the cycle-1 advisory-only hole). A guard that stays green here is a hole.
func TestReleaseWorkflowManifestTagGuardRejectsDroppedStep(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	if err := assertReleaseWorkflowGatesGoreleaserOnManifestTag(release); err != nil {
		t.Fatalf("real release.yml unexpectedly fails the manifest-tag-gate guard before mutation: %v", err)
	}

	adversarial := strings.Replace(release,
		`go run ./cmd/spacedock-release manifest-tag-gate "$GITHUB_REF_NAME" .claude-plugin/plugin.json .codex-plugin/plugin.json`,
		`echo "skipping manifest-tag gate"`,
		1)
	if adversarial == release {
		t.Fatal("fixture workflow missing the manifest-tag-gate invocation to drop")
	}
	if err := assertReleaseWorkflowGatesGoreleaserOnManifestTag(adversarial); err == nil {
		t.Fatal("manifest-tag-gate guard accepted a workflow with the gate step dropped")
	}
}

// TestReleaseWorkflowManifestTagGuardRejectsDroppedNeedsEdge is the adversarial
// twin for the needs-edge half: drop the goreleaser job's `needs:` on the gate job
// and the guard must RED, because goreleaser would run in parallel with (or before)
// the gate rather than waiting on it.
func TestReleaseWorkflowManifestTagGuardRejectsDroppedNeedsEdge(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	adversarial := strings.Replace(release,
		"  goreleaser:\n    needs: e2e-gate\n",
		"  goreleaser:\n",
		1)
	if adversarial == release {
		t.Fatal("fixture workflow missing the goreleaser `needs: e2e-gate` edge to drop")
	}
	if err := assertReleaseWorkflowGatesGoreleaserOnManifestTag(adversarial); err == nil {
		t.Fatal("manifest-tag-gate guard accepted a goreleaser job with the needs edge to the gate job dropped")
	}
}

// TestReleaseWorkflowManifestTagGateSkipsPreReleases locks the pre-release carve-
// out: the gate step must carry `if: !contains(github.ref, '-')`, matching the
// stamp step, so a `vX.Y.Z-pre.N` tag (whose manifest legitimately differs from
// the hyphenated tag semver) does NOT self-block. Without the skip every pre-release
// reds the gate (the audit confirmed `0.23.0-pre.1` vs manifest `0.23.0` exits 1).
func TestReleaseWorkflowManifestTagGateSkipsPreReleases(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	step := manifestTagGateStep(release)
	if step == nil {
		t.Fatal("release.yml has no manifest-tag-gate step")
	}
	if !ifSkipsPreRelease(step.ifCond) {
		t.Fatalf("manifest-tag-gate step does not skip pre-releases; if = %q, want a `!contains(github.ref, '-')` guard", step.ifCond)
	}
}

// ifSkipsPreRelease reports whether a step `if:` expression carries the
// `!contains(github.ref, '-')` pre-release skip (whitespace-insensitive).
func ifSkipsPreRelease(ifCond string) bool {
	normalized := strings.ReplaceAll(ifCond, " ", "")
	return strings.Contains(normalized, "!contains(github.ref,'-')")
}

// manifestTagGateStep returns the workflow step that invokes the manifest-tag-gate
// subcommand, or nil when none does.
func manifestTagGateStep(workflow string) *workflowStep {
	for _, job := range parseWorkflowJobs(workflow) {
		for i := range job.steps {
			for _, command := range executableShellCommands(job.steps[i].run) {
				if isManifestTagGateInvocation(command) {
					return &job.steps[i]
				}
			}
		}
	}
	return nil
}

// isManifestTagGateInvocation reports whether command runs the manifest-tag-gate
// subcommand over the tag ($GITHUB_REF_NAME) and both plugin manifests, so a step
// that drops the tag binding or a manifest is not recognized.
func isManifestTagGateInvocation(command string) bool {
	return strings.Contains(command, `go run ./cmd/spacedock-release manifest-tag-gate `) &&
		strings.Contains(command, `"$GITHUB_REF_NAME"`) &&
		strings.Contains(command, `.claude-plugin/plugin.json`) &&
		strings.Contains(command, `.codex-plugin/plugin.json`)
}

// assertReleaseWorkflowGatesGoreleaserOnManifestTag binds the AC-2 wiring to the
// parsed job graph: some job runs the manifest-tag-gate step over the tagged
// manifests, and the goreleaser carrier `needs:` that job. The predicate's own
// tag-vs-manifest comparison is unit-tested separately; this guard proves the
// WIRING gates the cut on it.
func assertReleaseWorkflowGatesGoreleaserOnManifestTag(workflow string) error {
	jobs := parseWorkflowJobs(workflow)

	var goreleaserCarriers []workflowJob
	gateJobNames := map[string]bool{}
	for _, job := range jobs {
		carriesGoreleaser := false
		invokesGate := false
		for _, step := range job.steps {
			if strings.HasPrefix(step.uses, "goreleaser/goreleaser-action@") {
				carriesGoreleaser = true
			}
			for _, command := range executableShellCommands(step.run) {
				if isManifestTagGateInvocation(command) {
					invokesGate = true
				}
			}
		}
		if carriesGoreleaser {
			goreleaserCarriers = append(goreleaserCarriers, job)
		}
		if invokesGate {
			gateJobNames[job.name] = true
		}
	}

	if len(goreleaserCarriers) == 0 {
		return fmt.Errorf("release.yml has no job carrying the goreleaser action")
	}
	if len(gateJobNames) == 0 {
		return fmt.Errorf("release.yml has no job that runs `spacedock-release manifest-tag-gate \"$GITHUB_REF_NAME\"` over both plugin manifests")
	}

	for _, carrier := range goreleaserCarriers {
		needsAGate := false
		for _, need := range carrier.needs {
			if gateJobNames[need] {
				needsAGate = true
			}
		}
		if !needsAGate {
			return fmt.Errorf("release.yml goreleaser job %q does not declare needs: on the manifest-tag-gate job — a tag/manifest mismatch could reach goreleaser unchecked", carrier.name)
		}
	}
	return nil
}
