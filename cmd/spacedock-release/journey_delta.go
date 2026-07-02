// ABOUTME: `spacedock-release journey-delta` renders and posts the AC-3 per-PR
// ABOUTME: journey-cost delta comment against the previously published release ledger.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
	"github.com/spacedock-dev/spacedock/internal/release"
)

// postCommentFunc posts (or updates) the rendered PR comment. It is a seam so
// the sticky-update argv is testable without a live `gh` call.
type postCommentFunc func(args []string) error

// ghPostComment invokes the real `gh` CLI. It is the production postCommentFunc;
// tests substitute a stub that records args instead.
func ghPostComment(args []string) error {
	cmd := exec.Command("gh", args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}

// journeyDelta computes and posts the AC-3 per-PR journey-cost delta comment:
// it diffs the current PR run's freshly emitted records (--metrics-dir) against
// the latest-by-captured_at observation per scenario/model in the previously
// published release ledger (the positional <previous-ledger.json> argument),
// renders the result as a single Markdown comment, and posts it via the sticky
// --edit-last/--create-if-none gh invocation so repeated pushes update ONE
// comment instead of appending a new one each time.
func journeyDelta(args []string, post postCommentFunc) int {
	if len(args) < 5 {
		fmt.Fprintln(os.Stderr, "spacedock-release journey-delta: need <previous-ledger.json> --metrics-dir <dir> --pr <number>")
		return 2
	}
	ledgerPath := args[0]
	var metricsDir, prNumber string
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--metrics-dir":
			i++
			if i >= len(args) {
				fmt.Fprintln(os.Stderr, "spacedock-release journey-delta: --metrics-dir needs a value")
				return 2
			}
			metricsDir = args[i]
		case "--pr":
			i++
			if i >= len(args) {
				fmt.Fprintln(os.Stderr, "spacedock-release journey-delta: --pr needs a value")
				return 2
			}
			prNumber = args[i]
		default:
			fmt.Fprintf(os.Stderr, "spacedock-release journey-delta: unknown argument %q\n", args[i])
			return 2
		}
	}
	if ledgerPath == "" || metricsDir == "" || prNumber == "" {
		fmt.Fprintln(os.Stderr, "spacedock-release journey-delta: need <previous-ledger.json> --metrics-dir <dir> --pr <number>")
		return 2
	}

	ledgerData, err := os.ReadFile(ledgerPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read previous ledger: %v\n", err)
		return 1
	}
	var ledger journeymetrics.Ledger
	if err := json.Unmarshal(ledgerData, &ledger); err != nil {
		fmt.Fprintf(os.Stderr, "parse previous ledger: %v\n", err)
		return 1
	}

	current, err := journeymetrics.ReadRecordsDir(metricsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read current run metrics: %v\n", err)
		return 1
	}

	deltas := release.ComputeJourneyDeltas(ledger, current)
	body := release.RenderJourneyDeltaComment(deltas)

	bodyFile, err := os.CreateTemp("", "journey-delta-comment-*.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create comment body file: %v\n", err)
		return 1
	}
	defer os.Remove(bodyFile.Name())
	if _, err := bodyFile.WriteString(body); err != nil {
		bodyFile.Close()
		fmt.Fprintf(os.Stderr, "write comment body file: %v\n", err)
		return 1
	}
	if err := bodyFile.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "close comment body file: %v\n", err)
		return 1
	}

	if err := post(release.JourneyDeltaCommentArgs(prNumber, bodyFile.Name())); err != nil {
		fmt.Fprintf(os.Stderr, "post PR comment: %v\n", err)
		return 1
	}
	fmt.Printf("posted journey-cost delta comment on PR %s (%d scenario/model rows)\n", prNumber, len(deltas))
	return 0
}
