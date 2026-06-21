// ABOUTME: Pure detector that fails LOUD and EARLY when the live FO booted the
// ABOUTME: wrong root (cd off the fixture / a workflow-dir outside it) — PR #365.
package ensigncycle

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

// detectWrongRootBoot scans a captured FO stream for the wrong-root wander that
// PR #365's opus run hit: a CI env leak (GITHUB_WORKSPACE naming the real repo)
// lured the FO to `cd` into the real checkout and boot its docs/dev workflow
// instead of the test's tmpdir fixture, after which it greeted-and-stopped with
// dispatchable:[] — surfacing only as a confusing pre-TeamCreate timeout. This
// detector turns that silent wander into a legible "FO booted the wrong root"
// failure naming the expected fixture root vs the actual wandered-to path.
//
// It is model-agnostic (it reads the tool-call stream, not any model-specific
// phrasing) and pure (stream + fixtureRoot in, error out), with its own offline
// test. The wander signatures it keys on, all observable in the boot stream:
//
//   - a `cd <abspath>` whose target escapes the fixture root,
//   - a `spacedock status --boot --workflow-dir <PATH>` whose PATH escapes it, and
//   - a `Read <dir>/README.md` (the boot's workflow-README read) outside it.
//
// It deliberately does NOT flag the legitimate real-repo paths a correct boot
// touches: the FO Reads its contract skills from the --plugin-dir checkout (the
// real repo) by design, so a contract Read outside the fixture is NOT a wander —
// only the WORKFLOW root (where it boots / cd's / reads the workflow README) must
// stay under the fixture.
func detectWrongRootBoot(stream, fixtureRoot string) error {
	clean := filepath.Clean(fixtureRoot)
	// On macOS the fixture root is a `/var/folders/...` symlink while the FO, having
	// cd'd in, reports the EvalSymlinks-resolved `/private/var/...` form — the same
	// directory. Resolve the fixture root so an under-fixture path in either form is
	// recognized, matching the EvalSymlinks guard the sibling live runners use.
	if resolved, err := filepath.EvalSymlinks(clean); err == nil {
		clean = resolved
	}
	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}
		var e streamEntry
		if json.Unmarshal([]byte(line), &e) != nil {
			continue
		}
		b := e.toolUseBlock()
		if b == nil {
			continue
		}
		switch b.Name {
		case "Bash":
			if target, ok := wanderTarget(b.Input.Command, clean); ok {
				return fmt.Errorf("FO booted the wrong root: expected the fixture root %q, but the boot command %q targets %q (outside the fixture) — a CI env leak likely lured the FO off its launch cwd",
					clean, strings.TrimSpace(b.Input.Command), target)
			}
		case "Read":
			// The FO reads {workflow_dir}/README.md at boot (Startup step 4). A
			// workflow README read OUTSIDE the fixture means it booted the wrong
			// workflow. Contract skills live under {plugin_dir}/skills/...references/,
			// never a bare <root>/README.md, so this does not flag a contract read.
			if target, ok := wanderWorkflowReadme(b.Input.FilePath, clean); ok {
				return fmt.Errorf("FO booted the wrong root: expected the fixture root %q, but it read the workflow README at %q (outside the fixture) — a CI env leak likely lured the FO off its launch cwd",
					clean, target)
			}
		}
	}
	return nil
}

// wanderWorkflowReadme returns the off-fixture absolute path of a workflow README
// read, when filePath is an absolute `<dir>/README.md` outside fixtureRoot. ok is
// false for a relative path, a README under the fixture, or any non-README read
// (a contract-skill Read under {plugin_dir}/skills is not a workflow README).
func wanderWorkflowReadme(filePath, fixtureRoot string) (string, bool) {
	if filePath == "" || filepath.Base(filePath) != "README.md" || !filepath.IsAbs(filePath) {
		return "", false
	}
	p := filepath.Clean(filePath)
	if isUnder(p, fixtureRoot) {
		return "", false
	}
	return p, true
}

// wanderTarget returns the off-fixture absolute path a boot command targets, when
// the command is a `cd <abspath>` or a `--workflow-dir <abspath>` resolving outside
// fixtureRoot. ok is false when the command names no such escaping path (it stays
// under the fixture, uses a relative path, or is an ordinary command).
func wanderTarget(command, fixtureRoot string) (string, bool) {
	for _, tok := range bootPathArgs(command) {
		if !filepath.IsAbs(tok) {
			continue
		}
		p := filepath.Clean(tok)
		if p == fixtureRoot || isUnder(p, fixtureRoot) {
			continue
		}
		return p, true
	}
	return "", false
}

// bootPathArgs pulls the path arguments a boot command supplies: the target of a
// leading `cd`, and the value after `--workflow-dir`. It splits on whitespace (the
// boot commands are simple `cd …`, `spacedock status --boot --workflow-dir …`
// forms; quoting is not exercised by the real boot stream). A compound boot chains
// the path token straight into a shell separator with no preceding space —
// `cd <root>; ls`, `cd <root>&& ls` — so each extracted token has any trailing
// `;`/`&&`/`|` trimmed before it reaches the isUnder check.
func bootPathArgs(command string) []string {
	fields := strings.Fields(command)
	var paths []string
	for i, f := range fields {
		switch {
		case f == "cd" && i+1 < len(fields):
			paths = append(paths, trimBootSeparator(fields[i+1]))
		case f == "--workflow-dir" && i+1 < len(fields):
			paths = append(paths, trimBootSeparator(fields[i+1]))
		case strings.HasPrefix(f, "--workflow-dir="):
			paths = append(paths, trimBootSeparator(strings.TrimPrefix(f, "--workflow-dir=")))
		}
	}
	return paths
}

// trimBootSeparator strips a single trailing shell separator (`;`, `&&`, `|`) that a
// compound boot glues onto a path token with no preceding space. It leaves an
// already-clean path untouched.
func trimBootSeparator(tok string) string {
	for _, sep := range []string{"&&", ";", "|"} {
		if strings.HasSuffix(tok, sep) {
			return strings.TrimSuffix(tok, sep)
		}
	}
	return tok
}

// isUnder reports whether path p is nested under dir (both already cleaned).
func isUnder(p, dir string) bool {
	rel, err := filepath.Rel(dir, p)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
