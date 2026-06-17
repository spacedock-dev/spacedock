//go:build live

// ABOUTME: N>=3 wrapper of the proven driveHaikuLoopOnce hand-loop drive; runs the loop
// ABOUTME: count times into per-run durable artifact dirs so every stream survives for AC-4/AC-5 grading.
package ensigncycle

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

// TestLiveHaikuLoopSpikeN wraps the proven single-drive driveHaikuLoopOnce in a
// count loop for the AC-5 must-build decision rule (a loop step that breaks on
// ANY of the N>=3 drives is must-build; "holds reliably" requires clean across
// ALL N). It is purely additive over the implementation-proven single drive: it
// does NOT touch driveHaikuLoopOnce, only calls it N times, each into its OWN
// run-{i} artifact subdir so all N tool-call streams are DURABLE on disk for the
// FO's classification (the single-drive path reused a fixed stream filename, which
// N runs would overwrite — the per-run subdir is the only difference).
//
// N defaults to 3 and is overridable via HAIKU_LOOP_N for "more if results vary".
// Each run's grade is logged so the FO reads the per-run breaks/holds matrix off
// the test log AND can re-grade from the durable stream files independently.
func TestLiveHaikuLoopSpikeN(t *testing.T) {
	binary := spacedockBinary(t)
	env := isolatedClaudeEnv(t, os.Getenv("HOME")) // t.Skip when no credential
	env = withBinaryOnPath(env, binary)

	n, err := strconv.Atoi(envOr("HAIKU_LOOP_N", "3"))
	if err != nil || n < 3 {
		t.Fatalf("HAIKU_LOOP_N=%q must be an integer >= 3 (the AC-5 decision rule needs N>=3)", envOr("HAIKU_LOOP_N", "3"))
	}

	base := claudeLiveArtifactDir(t, "haiku-loop-spike")
	grades := make([]haikuLoopGrade, 0, n)
	for i := 1; i <= n; i++ {
		runDir := filepath.Join(base, fmt.Sprintf("run-%d", i))
		if mkErr := os.MkdirAll(runDir, 0o755); mkErr != nil {
			t.Fatal(mkErr)
		}
		started := time.Now()
		grade, artifacts := driveHaikuLoopOnce(t, binary, env, runDir)
		t.Logf("haiku loop run %d/%d finished in %s; artifacts: %s\n  grade: %+v",
			i, n, time.Since(started).Round(time.Second), artifacts, grade)
		grades = append(grades, grade)
	}

	// AC-5 any-break classification, computed from the durable grades. A field that
	// is false on ANY run marks that loop step as breaks-on-any-N (must-build); a
	// field true across ALL N is holds-across-all-N. This is the drive's OWN
	// measurement; the external prior (trigger-carrying steps) is applied by the FO
	// on top of it, never folded into these counts.
	t.Logf("=== AC-5 any-break matrix over N=%d drives ===", n)
	report := func(label string, get func(haikuLoopGrade) bool) {
		breaks := 0
		for _, g := range grades {
			if !get(g) {
				breaks++
			}
		}
		verdict := "HOLDS-across-all-N"
		if breaks > 0 {
			verdict = "BREAKS-on-" + strconv.Itoa(breaks) + "-of-" + strconv.Itoa(n)
		}
		t.Logf("  %-22s %s", label, verdict)
	}
	report("entityLocated", func(g haikuLoopGrade) bool { return g.entityLocated })
	report("builtMarker", func(g haikuLoopGrade) bool { return g.builtMarker })
	report("statusDone", func(g haikuLoopGrade) bool { return g.statusDone })
	report("entityArchived", func(g haikuLoopGrade) bool { return g.entityArchived })
	report("verdictSet", func(g haikuLoopGrade) bool { return g.verdictSet })
	report("pathScopedCommit", func(g haikuLoopGrade) bool { return g.pathScopedCommit })
	report("dispatchBuildCalled", func(g haikuLoopGrade) bool { return g.dispatchBuildCalled })
	report("opusAgentSpawned", func(g haikuLoopGrade) bool { return g.opusAgentSpawned })
	report("gateVerdictFromL3", func(g haikuLoopGrade) bool { return g.gateVerdictFromL3 })

	// The wrapper does NOT fail on a break: a broken step is the SIGNAL the spike
	// measures (it becomes a must-build verb), not a harness failure. The FO reads
	// the matrix and the per-run streams to classify. The only hard failure is a
	// stalled/un-launchable drive, which driveHaikuLoopOnce already t.Fatals.
}
