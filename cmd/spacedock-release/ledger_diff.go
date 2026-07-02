// ABOUTME: `spacedock-release ledger-diff` is the AC-4 pre-upload safeguard CLI —
// ABOUTME: it exits non-zero if a rebuilt ledger changes anything beyond the named additions.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
	"github.com/spacedock-dev/spacedock/internal/release"
)

// ledgerDiff is the AC-4 safeguard gate: it must run and PASS before any `gh
// release upload --clobber` of a rebuilt ledger. Exit 0 means the rebuilt
// ledger equals the original plus exactly the --added scenario ids, with every
// pre-existing scenario byte-for-byte unchanged; any other outcome exits 1 and
// the caller must NOT upload.
func ledgerDiff(args []string) int {
	if len(args) < 4 {
		fmt.Fprintln(os.Stderr, "spacedock-release ledger-diff: need <original.json> <rebuilt.json> --added <scenario-id>[,<scenario-id>...]")
		return 2
	}
	originalPath, rebuiltPath := args[0], args[1]
	var added string
	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "--added":
			i++
			if i >= len(args) {
				fmt.Fprintln(os.Stderr, "spacedock-release ledger-diff: --added needs a value")
				return 2
			}
			added = args[i]
		default:
			fmt.Fprintf(os.Stderr, "spacedock-release ledger-diff: unknown argument %q\n", args[i])
			return 2
		}
	}
	if originalPath == "" || rebuiltPath == "" || added == "" {
		fmt.Fprintln(os.Stderr, "spacedock-release ledger-diff: need <original.json> <rebuilt.json> --added <scenario-id>[,<scenario-id>...]")
		return 2
	}

	original, err := readLedgerFile(originalPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read original ledger %s: %v\n", originalPath, err)
		return 1
	}
	rebuilt, err := readLedgerFile(rebuiltPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read rebuilt ledger %s: %v\n", rebuiltPath, err)
		return 1
	}

	if err := release.DiffAddedScenarios(original, rebuilt, strings.Split(added, ",")...); err != nil {
		fmt.Fprintf(os.Stderr, "ledger-diff: UNSAFE TO UPLOAD — %v\n", err)
		return 1
	}
	fmt.Printf("ledger-diff: safe to upload — %s equals %s plus exactly %s\n", rebuiltPath, originalPath, added)
	return 0
}

func readLedgerFile(path string) (journeymetrics.Ledger, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return journeymetrics.Ledger{}, err
	}
	var ledger journeymetrics.Ledger
	if err := json.Unmarshal(data, &ledger); err != nil {
		return journeymetrics.Ledger{}, err
	}
	return ledger, nil
}
