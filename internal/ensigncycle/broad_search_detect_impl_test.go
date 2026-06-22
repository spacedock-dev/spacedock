// ABOUTME: Pure detector that reds when a zero-discover FO boot broad-searches the
// ABOUTME: filesystem to hunt a workflow instead of report-and-stop — lean-boot guard.
package ensigncycle

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

// detectBroadSearchAtBoot scans a captured FO boot stream for the lean-boot
// violation the captain observed (2026-06-14): after `spacedock status --discover`
// returned zero workflows, an FO ran a broad `find`/`grep -r`/`ls -R` filesystem
// sweep to hunt a workflow down instead of obeying the contract's terminal zero
// branch (Startup step 3: zero → report no workflow found and STOP). A broad sweep
// at boot is both a discipline violation (the zero branch is report-and-stop) and a
// cost/latency regression — the opposite of lean boot.
//
// Like detectWrongRootBoot it is model-agnostic (it reads the tool-call stream, not
// any model phrasing) and pure (stream + fixtureRoot in, error out), with its own
// offline test. It iterates ALL tool_use blocks of each entry (not just the first),
// so a sweep cannot evade it by riding as a second block of a multi-tool turn. The
// reddable sweep signatures, all observable in the boot stream:
//
//   - a `Bash` command invoking find / grep -r / rg / fd / ls -R whose target is the
//     project root or a broad ancestor (not a scoped path under a resolved workflow),
//   - a `Glob` tool_use with a recursive workflow-hunting pattern (e.g. **/README.md),
//   - a `Grep` tool_use whose path is the project root / unset (a repo-wide search).
//
// A correct zero-discover boot touches none of these: it runs --version,
// git rev-parse, status --discover (zero), then reports no-workflow-found and stops.
// The detector passes that and reds the sweep. It guards STRICTLY the zero-discover
// branch's substituted sweep — a scoped search under an already-resolved workflow
// dir is legitimate and must pass (keyed on the target being the root / a recursive
// pattern, not on the tool name alone).
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
					return fmt.Errorf("FO broad-searched the filesystem at boot: command %q runs %s over the project root %q — after a zero `status --discover` the Startup zero branch is report-and-stop, not a filesystem sweep to hunt a workflow",
						strings.TrimSpace(b.Input.Command), sig, clean)
				}
			case "Glob":
				if recursiveHuntPattern(b.Input.Pattern) {
					return fmt.Errorf("FO broad-searched the filesystem at boot: a recursive Glob %q hunts a workflow project-wide — after a zero `status --discover` the Startup zero branch is report-and-stop",
						b.Input.Pattern)
				}
			case "Grep":
				if repoWideGrep(b.Input.Path, clean) {
					return fmt.Errorf("FO broad-searched the filesystem at boot: a repo-wide Grep for %q (path %q) hunts a workflow project-wide — after a zero `status --discover` the Startup zero branch is report-and-stop",
						b.Input.Pattern, b.Input.Path)
				}
			}
		}
	}
	return nil
}

// broadSweepTools are the directory-listing/search shell tools a workflow-hunting
// sweep uses, broad whenever they target the project root (the zero branch forbids
// ANY root-scoped find/ls, not just a recursive one). `grep` still needs `-r` to be
// a sweep — a non-recursive `grep pattern file` is a single-file read, not a hunt.
var broadSweepTools = []string{"find", "ls", "rg", "fd"}

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
	case tool == "ls" && hasRecursiveFlag(fields):
		sig = "ls -R" // name the recursive form distinctly; a plain root `ls` reds as "ls"
	case contains(broadSweepTools, tool):
		sig = tool
	case tool == "grep" && hasRecursiveFlag(fields):
		sig = "grep -r"
	default:
		return "", false
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

// hasRecursiveFlag reports whether a grep/ls command carries its recursive flag
// (-r/-R, possibly bundled like -rn or -Rl).
func hasRecursiveFlag(fields []string) bool {
	for _, f := range fields[1:] {
		if !strings.HasPrefix(f, "-") || strings.HasPrefix(f, "--") {
			continue
		}
		if strings.ContainsAny(f, "rR") {
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
