// ABOUTME: Offline-testable end-state helpers for the live cycle test: locate the
// ABOUTME: entity (in place OR archived) and scan the git log for a path-scoped commit.
package ensigncycle

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	// liveStageReportHeading anchors the appended stage-report heading for ANY
	// stage. The deterministic skeleton drives the fixture ONE stage in place so it
	// can pin the exact `## Stage Report: backlog` heading, but a real full FO cycle
	// runs the entity all the way to the TERMINAL stage, so the ensign that finishes
	// the cycle writes `## Stage Report: done`. The live test gates on the protocol
	// SHAPE (a stage-report section exists), not the stage word — `\S` after the
	// colon requires a named stage while rejecting an entity with no stage report at
	// all (the incomplete-cycle shape, which goes red).
	liveStageReportHeading = regexp.MustCompile(`(?m)^## Stage Report: \S`)
	// frontmatterField anchors the terminal `status: done` frontmatter line the
	// FO writes when the entity reaches the terminal stage. Anchored at the line
	// start so a `status:` mention in prose cannot satisfy it.
	frontmatterField = regexp.MustCompile(`(?im)^status:\s*done\s*$`)
	// verdictSet anchors a finalized `verdict:` line carrying a non-empty value.
	// The exact verdict WORD is FO judgment that varies by model (sonnet wrote
	// `verdict: done`, opus wrote `verdict: passed`); both completed the full
	// cycle. The live test gates on MECHANICAL completion, so it only requires the
	// verdict be SET — `\S` after the colon rejects an empty/whitespace-only
	// `verdict:` (the incomplete-cycle shape) while accepting any decided value.
	// `[^\S\n]*` is horizontal whitespace only so it cannot consume the line break
	// and let `\S` reach into the next frontmatter line.
	verdictSet = regexp.MustCompile(`(?im)^verdict:[^\S\n]*\S.*$`)
)

// A real FO driving the fixture entity to the TERMINAL `done` stage ARCHIVES it:
// the flat `make-it-work.md` moves to `_archive/make-it-work.md`. These helpers
// locate the entity wherever the completed cycle left it and inspect the git log
// for the path-scoped state commit, so the live assertions match the REAL
// completed-and-archived end-state rather than the scripted skeleton's single
// in-place append. They live under the DEFAULT build tags (no //go:build live) so
// the live test reuses them AND an offline unit test exercises the locate/scan
// logic against a staged fixture without spending a model.

// locateEntity returns the contents of the fixture entity after the cycle,
// searching the three end-state locations in order: the original flat path, the
// flat archive `_archive/<slug>.md`, and the folder archive `_archive/<slug>/index.md`.
// found reports whether any of them existed. An INCOMPLETE cycle that never
// archived still resolves the original path; a completed cycle resolves the
// archive. When none exist (the entity vanished entirely), found is false and the
// caller fails loudly rather than asserting against empty content.
func locateEntity(root, slug string) (content string, where string, found bool) {
	candidates := []string{
		filepath.Join(root, slug+".md"),
		filepath.Join(root, "_archive", slug+".md"),
		filepath.Join(root, "_archive", slug, "index.md"),
	}
	for _, p := range candidates {
		if b, err := os.ReadFile(p); err == nil {
			return string(b), p, true
		}
	}
	return "", "", false
}

// someCommitNamesOnly reports whether ANY commit in the entity's history named
// ONLY the entity slug — the path-scoped state-commit invariant at the cycle
// level. A full FO-to-done cycle makes MULTIPLE commits (the ensign's path-scoped
// state commit, then the FO's archive/finalize commits), so HEAD is no longer the
// path-scoped one; this scans the whole log for at least one path-scoped commit
// instead of pinning HEAD. The haiku INCOMPLETE cycle committed `[README.md
// make-it-work.md]` (a sibling sweep, not path-scoped) and produced no other
// commit, so this returns false for it — keeping the live test red on an
// incomplete cycle.
func someCommitNamesOnly(t *testing.T, root, slug string) bool {
	t.Helper()
	return pathScopedCommitCount(t, root, slug) >= 1
}

// pathScopedCommitCount counts the commits in the entity's history that named ONLY
// the entity slug (the path-scoped state-commit invariant). It is the shared walk
// behind someCommitNamesOnly (>=1) and the integration-transition grade (>=2: the
// implementation->integration commit AND the integration->done+archive commit must
// both be path-scoped). A `git add -A` / sibling sweep is NOT path-scoped and is not
// counted.
func pathScopedCommitCount(t *testing.T, root, slug string) int {
	t.Helper()
	// One commit per line, name-only files separated by tabs after a leading
	// marker so per-commit boundaries are unambiguous.
	out := git(t, root, "log", "--pretty=format:@@COMMIT@@", "--name-only")
	target := slug + ".md"
	count := 0
	var files []string
	flush := func() {
		if len(files) == 1 && filepath.Base(files[0]) == target {
			count++
		}
		files = files[:0]
	}
	for _, ln := range strings.Split(out, "\n") {
		ln = strings.TrimSpace(ln)
		if ln == "@@COMMIT@@" {
			flush()
			continue
		}
		if ln != "" {
			files = append(files, ln)
		}
	}
	flush()
	return count
}

// integrationTransitionCommitted is the BLOCKER-fix grade field: it proves the Haiku
// FO held the FULL implementation->integration->done loop rather than shortcutting
// implementation->done. It returns true only when BOTH hold:
//
//  1. the state-checkout git log contains at least one commit whose committed entity
//     blob carried `status: integration` (`git log -G'^status: integration$' -- {slug}.md`
//     returns >=1 commit), AND
//  2. there were at least TWO path-scoped commits (the implementation->integration set
//     and the integration->done+archive set).
//
// A run that jumped implementation->done produces exactly ONE path-scoped commit (the
// terminalize) and NO committed blob ever carried `status: integration`, so BOTH
// conditions fail — even though statusDone/verdictSet/builtMarker/pathScopedCommit
// (which only read the TERMINAL state or the FIRST path-scoped commit) all pass. This
// is the ONLY field that distinguishes "held the full loop" from "skipped to the end
// and still produced a clean terminal".
func integrationTransitionCommitted(t *testing.T, root, slug string) bool {
	t.Helper()
	target := slug + ".md"
	out := git(t, root, "log", "-G", `^status: integration$`, "--pretty=format:%H", "--", target)
	committedIntegrationBlob := strings.TrimSpace(out) != ""
	return committedIntegrationBlob && pathScopedCommitCount(t, root, slug) >= 2
}
