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
//   - a `spacedock status --boot --workflow-dir <PATH>` whose PATH escapes the
//     fixture root, standalone and fatal,
//   - a `Read <dir>/README.md` (the boot's workflow-README read) outside it,
//     standalone and fatal, and
//   - a `cd <abspath>` whose target escapes the fixture root, but ONLY when the
//     SAME command also carries a workflow-operative token (see
//     hasWorkflowOperativeSignature) — a bare `cd <outside>` alone is sonnet's
//     speculative repo-root sniff, a harmless version-probe that neither
//     persists past the command nor drives any workflow operation.
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
				return fmt.Errorf("FO booted the wrong root: expected the fixture root %q, but the boot command %q targets %q (outside the fixture) — the FO's boot operated outside its launch cwd",
					clean, strings.TrimSpace(b.Input.Command), target)
			}
		case "Read":
			// The FO reads {workflow_dir}/README.md at boot (Startup step 4). A
			// workflow README read OUTSIDE the fixture means it booted the wrong
			// workflow. Contract skills live under {plugin_dir}/skills/...references/,
			// never a bare <root>/README.md, so this does not flag a contract read.
			if target, ok := wanderWorkflowReadme(b.Input.FilePath, clean); ok {
				return fmt.Errorf("FO booted the wrong root: expected the fixture root %q, but it read the workflow README at %q (outside the fixture) — the FO's boot operated outside its launch cwd",
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
// under the fixture, uses a relative path, or is an ordinary command) — or when
// the only escaping path is an uncorroborated bare `cd` (see
// hasWorkflowOperativeSignature).
func wanderTarget(command, fixtureRoot string) (string, bool) {
	corroborated := hasWorkflowOperativeSignature(command)
	for _, arg := range bootPathArgs(command) {
		if !filepath.IsAbs(arg.path) {
			continue
		}
		p := filepath.Clean(arg.path)
		if p == fixtureRoot || isUnder(p, fixtureRoot) {
			continue
		}
		if arg.cd && !corroborated {
			// A bare cd escaping the fixture, alone, is sonnet's speculative
			// repo-root sniff — a harmless version-probe that does not persist
			// (the FO's very next command runs from the fixture root again) and
			// drives no workflow operation. Standalone --workflow-dir escapes
			// (arg.cd == false) are unaffected: they stay first-and-fatal.
			continue
		}
		return p, true
	}
	return "", false
}

// workflowOperativeSubstrings are same-command tokens that turn a bare `cd
// <outside-fixture>` into a genuine wander: the command goes on to operate the
// workflow, rather than merely probing a version or toplevel from wherever the
// cd landed.
var workflowOperativeSubstrings = []string{
	"--workflow-dir",
	"--boot",
	"--discover",
	"status --read",
	"state commit",
}

// hasWorkflowOperativeSignature reports whether command carries a
// workflow-operative token: a spacedock workflow flag, a state-checkout commit,
// `spacedock new`, or a README/entity-path reference. Paired with a same-command
// `cd <outside-fixture>`, this distinguishes an actual attempt to operate the
// workflow from outside the fixture from sonnet's harmless speculative
// repo-root sniff.
func hasWorkflowOperativeSignature(command string) bool {
	for _, s := range workflowOperativeSubstrings {
		if strings.Contains(command, s) {
			return true
		}
	}
	for _, f := range strings.Fields(command) {
		f = trimBootSeparator(f)
		if f == "new" {
			return true
		}
		if strings.Contains(f, "README") || strings.HasSuffix(strings.ToLower(f), ".md") {
			return true
		}
	}
	return false
}

// bootPathArg is one path argument a boot command supplies, tagged with
// whether it came from a `cd` (subject to the corroboration gate) or a
// `--workflow-dir` (standalone-fatal, ungated).
type bootPathArg struct {
	path string
	cd   bool
}

// bootPathArgs pulls the path arguments a boot command supplies: the target of a
// leading `cd`, and the value after `--workflow-dir`. It splits on whitespace (the
// boot commands are simple `cd …`, `spacedock status --boot --workflow-dir …`
// forms; quoting is not exercised by the real boot stream). A compound boot chains
// the path token straight into a shell separator with no preceding space —
// `cd <root>; ls`, `cd <root>&& ls` — so each extracted token has any trailing
// `;`/`&&`/`|` trimmed before it reaches the isUnder check.
func bootPathArgs(command string) []bootPathArg {
	fields := strings.Fields(command)
	var paths []bootPathArg
	for i, f := range fields {
		switch {
		case f == "cd" && i+1 < len(fields):
			paths = append(paths, bootPathArg{path: trimBootSeparator(fields[i+1]), cd: true})
		case f == "--workflow-dir" && i+1 < len(fields):
			paths = append(paths, bootPathArg{path: trimBootSeparator(fields[i+1])})
		case strings.HasPrefix(f, "--workflow-dir="):
			paths = append(paths, bootPathArg{path: trimBootSeparator(strings.TrimPrefix(f, "--workflow-dir="))})
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
