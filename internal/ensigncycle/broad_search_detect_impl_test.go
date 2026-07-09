// ABOUTME: Pure detector that reds when a zero-discover FO boot broad-searches the
// ABOUTME: filesystem to hunt a workflow instead of report-and-stop — lean-boot guard.
package ensigncycle

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

// detectBroadSearchAtBoot scans a captured FO boot stream for a broad filesystem
// sweep at boot: a `find`/`grep -r`/`ls -R` (or a recursive Glob/repo-wide Grep)
// hunting a workflow or a contract reference file across the project root instead
// of proceeding on the fixture/plugin content already in hand. Two boot shapes
// trigger it: (1) after `spacedock status --discover` returns zero workflows, the
// contract's terminal zero branch is report-and-stop, not a sweep to hunt one down
// (the captain observed this 2026-06-14); (2) after a normal boot resolves its
// workflow, the FO goes looking for a contract reference file (e.g. a skill's
// cross-referenced doc) that lives only in the --plugin-dir checkout, not the
// fixture — also a sweep, not a scoped read. A broad sweep at boot is both a
// discipline violation and a cost/latency regression — the opposite of lean boot.
//
// It is model-agnostic (it reads the tool-call stream, not any model phrasing) and
// pure (stream + fixtureRoot in, error out), with its own offline test. It iterates
// ALL tool_use blocks of each entry (not just the first), so a sweep cannot evade
// it by riding as a second block of a multi-tool turn. The reddable sweep
// signatures, all observable in the boot stream:
//
//   - a `Bash` command invoking find / grep -r / rg / fd / ls -R whose target is the
//     project root or a broad ancestor (not a scoped path under a resolved workflow),
//   - a `Glob` tool_use with a recursive hunting pattern (e.g. **/README.md),
//   - a `Grep` tool_use whose path is the project root / unset (a repo-wide search).
//
// A correct boot touches none of these. A boot MAY also run a plain, non-recursive
// `ls` (e.g. `ls -la` of its own cwd, or `ls {root}`) while orienting — the
// captain's decision (2026-07-02) is that flat `ls` is not a hunt, so the detector
// never reds it. The detector passes all of that and reds only the genuine sweep.
// It guards STRICTLY a substituted sweep — a scoped search under an
// already-resolved workflow dir is legitimate and must pass. The banned axis is
// recursion/hunting, not the `ls` binary or the root path (keyed on recursion for
// `ls`; on the target being the root / an unscoped path for
// find/rg/fd/grep -r/Glob/Grep).
func detectBroadSearchAtBoot(stream, fixtureRoot string) error {
	clean := filepath.Clean(fixtureRoot)
	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}
		var e streamEntry
		if json.Unmarshal([]byte(line), &e) != nil {
			continue
		}
		for _, b := range e.toolUseBlocks() {
			switch b.Name {
			case "Bash":
				if sig, ok := broadSweepCommand(b.Input.Command, clean); ok {
					return fmt.Errorf("FO broad-searched the filesystem at boot: command %q runs %s over the project root %q — a boot-preamble filesystem sweep, not the scenario's own assertion",
						strings.TrimSpace(b.Input.Command), sig, clean)
				}
			case "Glob":
				if recursiveHuntPattern(b.Input.Pattern) {
					return fmt.Errorf("FO broad-searched the filesystem at boot: a recursive Glob %q hunts the project root — a boot-preamble filesystem sweep, not the scenario's own assertion",
						b.Input.Pattern)
				}
			case "Grep":
				if repoWideGrep(b.Input.Path, clean) {
					return fmt.Errorf("FO broad-searched the filesystem at boot: a repo-wide Grep for %q (path %q) hunts the project root — a boot-preamble filesystem sweep, not the scenario's own assertion",
						b.Input.Pattern, b.Input.Path)
				}
			}
		}
	}
	return nil
}

// broadSweepTools are the recursive-by-default search shell tools a workflow-hunting
// sweep uses, broad whenever they target the project root (not scoped under an
// already-resolved workflow dir). `ls` is NOT here — a plain `ls` is a single-level
// listing, not a hunt; it reds separately, only when its own recursive flag (`-R`)
// or a globstar path argument is present. `grep` still needs `-r` to be a sweep — a
// non-recursive `grep pattern file` is a single-file read, not a hunt.
var broadSweepTools = []string{"find", "rg", "fd"}

// broadSweepCommand reports whether a Bash command is a broad filesystem sweep
// aimed at the project root or a broad ancestor (not a scoped path under a resolved
// workflow dir). It returns the matched signature for the error message. A search
// whose only path arguments stay UNDER the fixture root's subtree (e.g. a scoped
// grep inside an already-resolved workflow dir) is legitimate and not flagged.
func broadSweepCommand(command, fixtureRoot string) (string, bool) {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return "", false
	}
	tool := fields[0]

	var sig string
	switch {
	case tool == "ls" && (hasRecursiveFlag(fields, "R") || hasGlobstarPathArg(fields)):
		sig = "ls -R" // -R, or a globstar path arg — ls's recursive-descent equivalents
	case contains(broadSweepTools, tool):
		sig = tool
	case tool == "grep" && hasRecursiveFlag(fields, "rR"):
		sig = "grep -r"
	default:
		return "", false // includes a plain, non-recursive `ls` — not a hunt
	}

	// A sweep is broad when it targets the project root / a broad ancestor — i.e. it
	// names no path scoped under the fixture's subtree. A path equal to the fixture
	// root itself is the root sweep (red); a path strictly under it (a resolved
	// workflow dir) is scoped (pass). A sweep with NO path arg defaults to cwd (the
	// boot root) — also broad.
	for _, tok := range fields[1:] {
		if strings.HasPrefix(tok, "-") {
			continue
		}
		if !filepath.IsAbs(tok) {
			continue
		}
		p := filepath.Clean(tok)
		if p != fixtureRoot && isUnder(p, fixtureRoot) {
			return "", false // scoped under a resolved workflow dir — legitimate
		}
	}
	return sig, true
}

// hasRecursiveFlag reports whether a command carries a short flag drawn from
// flagChars (possibly bundled, e.g. -rn or -ltr), the command's own recursive-flag
// alphabet: grep accepts both -r and -R, but ls accepts only -R — for ls, -r is
// reverse-sort, not recursion, so a caller must not share one alphabet across tools.
func hasRecursiveFlag(fields []string, flagChars string) bool {
	for _, f := range fields[1:] {
		if !strings.HasPrefix(f, "-") || strings.HasPrefix(f, "--") {
			continue
		}
		if strings.ContainsAny(f, flagChars) {
			return true
		}
	}
	return false
}

// hasGlobstarPathArg reports whether any non-flag argument contains a `**`
// globstar — ls's recursive-descent equivalent of the banned recursive Glob
// pattern (e.g. `ls **/README.md`).
func hasGlobstarPathArg(fields []string) bool {
	for _, f := range fields[1:] {
		if !strings.HasPrefix(f, "-") && strings.Contains(f, "**") {
			return true
		}
	}
	return false
}

// recursiveHuntPattern reports whether a Glob pattern recursively hunts a workflow
// project-wide — a `**/` recursive descent (e.g. **/README.md, **/*.md). A scoped
// non-recursive pattern is not a project-wide hunt.
func recursiveHuntPattern(pattern string) bool {
	return strings.Contains(pattern, "**")
}

// repoWideGrep reports whether a Grep tool_use searches the project root / a broad
// ancestor (an unset path is repo-wide by default; a path equal to the fixture root
// is the root sweep). A path scoped under the fixture root is legitimate.
func repoWideGrep(path, fixtureRoot string) bool {
	if path == "" {
		return true
	}
	if !filepath.IsAbs(path) {
		return false
	}
	p := filepath.Clean(path)
	return p == fixtureRoot || !isUnder(p, fixtureRoot)
}

// contains reports whether s is in list.
func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
