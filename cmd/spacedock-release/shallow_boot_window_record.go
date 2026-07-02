// ABOUTME: `spacedock-release shallow-boot-window-record` extracts the AC-1
// ABOUTME: boot-window observation from an archived stream, for the AC-4 backfill.
package main

import (
	"fmt"
	"os"

	"github.com/spacedock-dev/spacedock/internal/ensigncycle"
	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

// shallowBootWindowRecord applies AC-1's extraction logic
// (ensigncycle.BuildShallowBootWindowRecord) to an ARCHIVED claude-stream.jsonl,
// standalone rather than at live-test time, and emits the resulting
// shallow-boot-window record into --out — the same mechanism the AC-4 backfill
// uses to add the observation to a historical run's already-recovered metrics
// dir before that run's ledger is rebuilt.
func shallowBootWindowRecord(args []string) int {
	if len(args) < 5 {
		fmt.Fprintln(os.Stderr, "spacedock-release shallow-boot-window-record: need <claude-stream.jsonl> --model <model> --out <metrics-dir>")
		return 2
	}
	streamPath := args[0]
	var model, outDir string
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--model":
			i++
			if i >= len(args) {
				fmt.Fprintln(os.Stderr, "spacedock-release shallow-boot-window-record: --model needs a value")
				return 2
			}
			model = args[i]
		case "--out":
			i++
			if i >= len(args) {
				fmt.Fprintln(os.Stderr, "spacedock-release shallow-boot-window-record: --out needs a value")
				return 2
			}
			outDir = args[i]
		default:
			fmt.Fprintf(os.Stderr, "spacedock-release shallow-boot-window-record: unknown argument %q\n", args[i])
			return 2
		}
	}
	if streamPath == "" || model == "" || outDir == "" {
		fmt.Fprintln(os.Stderr, "spacedock-release shallow-boot-window-record: need <claude-stream.jsonl> --model <model> --out <metrics-dir>")
		return 2
	}

	data, err := os.ReadFile(streamPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", streamPath, err)
		return 1
	}
	turns, err := journeymetrics.ParseClaudeTurns(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse claude turns: %v\n", err)
		return 1
	}
	record, err := ensigncycle.BuildShallowBootWindowRecord(turns, model)
	if err != nil {
		fmt.Fprintf(os.Stderr, "build shallow-boot-window record: %v\n", err)
		return 1
	}
	if err := journeymetrics.EmitRecord(outDir, record); err != nil {
		fmt.Fprintf(os.Stderr, "emit record: %v\n", err)
		return 1
	}
	fmt.Printf("wrote shallow-boot-window record to %s (turns=%d, tokens=%+v)\n", outDir, record.Turns, record.Tokens)
	return 0
}
