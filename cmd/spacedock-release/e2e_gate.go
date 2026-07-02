// ABOUTME: `spacedock-release e2e-gate <commit>` — release-time precondition that
// ABOUTME: blocks the cut unless a green Runtime Live E2E run exists for the commit.
package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spacedock-dev/spacedock/internal/release"
)

// runListFunc fetches the `gh run list` JSON for the release commit. It is a seam
// so the gate's exit-code + step-summary behavior is testable against fixtures
// without a live gh.
type runListFunc func(commit string) ([]byte, error)

// ghRunListForCommit queries the Runtime Live E2E runs whose head commit is the
// release commit; `-c <commit>` binds the query to the exact tagged commit so a
// green run on some other line cannot satisfy the gate. The query deliberately
// omits `--status success`: gh's `--status` filter can lag a run's `conclusion`
// right after a re-run flips it to green (observed at the v0.23.0 cut — the
// filtered query returned `[]` while the same query without `--status` already
// saw the success), so the pre-filter costs genuine re-run-to-green cuts and
// buys nothing. Excluding parked/waiting runs is EvaluateE2EGate's job: its
// predicate requires conclusion == "success", so an empty-conclusion parked run
// never qualifies regardless of what this query returns.
func ghRunListForCommit(commit string) ([]byte, error) {
	cmd := exec.Command("gh", "run", "list",
		"--workflow", "Runtime Live E2E",
		"-c", commit,
		"--json", "databaseId,headSha,conclusion,status",
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gh run list: %w", err)
	}
	return out, nil
}

// runE2EGate evaluates the release-time e2e gate for the commit given as the sole
// argument. It reads the captain-waiver reason from SPACEDOCK_E2E_GATE_WAIVER,
// fetches the run list via list, runs the pure decision predicate, records the
// outcome to $GITHUB_STEP_SUMMARY (the audit trail), and returns the process exit
// code: 0 when the gate passes (green run matched, or explicitly waived), 1 when
// it blocks. A query or parse failure with no waiver blocks the cut.
func runE2EGate(args []string, list runListFunc) int {
	if len(args) != 1 || args[0] == "" {
		fmt.Fprintln(os.Stderr, "spacedock-release e2e-gate: need exactly one <release-commit-sha>")
		return 2
	}
	commit := args[0]
	waiver := os.Getenv("SPACEDOCK_E2E_GATE_WAIVER")

	var runListJSON []byte
	if waiver == "" {
		// Only consult gh when there is no waiver; a waiver short-circuits the
		// predicate before the run list, so an emergency cut survives a gh failure.
		out, err := list(commit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "spacedock-release e2e-gate: %v\n", err)
			recordGateSummary(fmt.Sprintf("e2e gate BLOCKED for %s: could not query Runtime Live E2E runs (%v)", commit, err))
			return 1
		}
		runListJSON = out
	}

	dec, err := release.EvaluateE2EGate(runListJSON, commit, waiver)
	if err != nil {
		fmt.Fprintf(os.Stderr, "spacedock-release e2e-gate: %v\n", err)
		recordGateSummary(fmt.Sprintf("e2e gate BLOCKED for %s: %v", commit, err))
		return 1
	}

	recordGateSummary(dec.Reason)
	if !dec.Pass {
		fmt.Fprintf(os.Stderr, "spacedock-release e2e-gate: %s\n", dec.Reason)
		return 1
	}
	fmt.Println(dec.Reason)
	return 0
}

// recordGateSummary appends the gate decision to $GITHUB_STEP_SUMMARY when set
// (the GitHub Actions step-summary file) so every cut — passed, blocked, or
// waived — leaves an auditable record. Outside CI the env var is unset and this
// is a no-op.
func recordGateSummary(reason string) {
	path := os.Getenv("GITHUB_STEP_SUMMARY")
	if path == "" {
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "### Release e2e gate\n\n%s\n", reason)
}
