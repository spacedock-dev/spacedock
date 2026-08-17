// ABOUTME: Entered working stages remain their own dispatch target until their
// ABOUTME: latest stage report is structurally complete and committed cleanly.
package status

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// enteredStageAwaitingCompletion reports whether a current working stage must
// be dispatched before ordinary successor projection resumes. Initial stages
// and gate, terminal, and worktree suppression are handled by the caller.
func enteredStageAwaitingCompletion(e *entity, stage Stage) bool {
	if stage.initial || stage.gate || stage.terminal {
		return false
	}
	return !hasCompleteCommittedStageReport(e.path, stage.Name)
}

// gatePreparable reports whether a gated stage's promotion proof is satisfied.
// A non-initial gated stage owes a complete committed stage report. An INITIAL
// gated stage had no prior stage to write one: the committed clean seed IS the
// artifact the captain reviews, so durability alone is the proof. The exception
// keys on stage.initial ONLY — never on "no report exists" — so no later stage
// gains a path around its completion proof.
func gatePreparable(path string, stage Stage) bool {
	if stage.initial {
		return entityPathCleanInHEAD(path)
	}
	return hasCompleteCommittedStageReport(path, stage.Name)
}

// hasCompleteCommittedStageReport applies the status-owned, mechanical half of
// completion proof. The existing report selector chooses only the latest exact
// stage token. Its checklist must be non-empty, contain no FAILED or blank
// items, give every item a non-empty evidence/rationale line, and end with a
// non-empty Summary. The current entity bytes must also be tracked and clean
// against the local HEAD; sibling dirt is deliberately outside that pathspec.
func hasCompleteCommittedStageReport(path, stage string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return hasCompleteStageReport(data, stage) && entityPathCleanInHEAD(path)
}

func hasCompleteStageReport(data []byte, stage string) bool {
	lines := splitLines(string(data))
	start, end, ok := selectStageReport(lines, stage)
	if !ok {
		return false
	}
	items := extractChecklist(lines, start, end)
	if len(items) == 0 {
		return false
	}
	for _, item := range items {
		if (item.status != "DONE" && item.status != "SKIPPED") || strings.TrimSpace(item.text) == "" || !checklistItemHasEvidence(lines, item) {
			return false
		}
	}
	return stageReportHasSummary(lines, start, end)
}

func checklistItemHasEvidence(lines []string, item checklistItem) bool {
	for line := item.start + 1; line <= item.end; line++ {
		if strings.TrimSpace(lines[line-1]) != "" {
			return true
		}
	}
	return false
}

func stageReportHasSummary(lines []string, start, end int) bool {
	for line := start; line <= end; line++ {
		if strings.TrimSpace(lines[line-1]) != "### Summary" {
			continue
		}
		for bodyLine := line + 1; bodyLine <= end; bodyLine++ {
			level, _, heading := parseATXHeading(lines[bodyLine-1])
			if heading && level <= 3 {
				break
			}
			if strings.TrimSpace(lines[bodyLine-1]) != "" {
				return true
			}
		}
		return false
	}
	return false
}

// entityPathCleanInHEAD is deliberately literal and path-scoped. A tracked
// entity whose working-tree/index bytes differ from HEAD is not durable proof;
// an unrelated dirty sibling does not affect the answer.
func entityPathCleanInHEAD(path string) bool {
	gitRoot, rel, ok := entityGitPath(path)
	if !ok {
		return false
	}
	rel = filepath.ToSlash(rel)
	tracked := exec.Command("git", "--literal-pathspecs", "ls-files", "--error-unmatch", "--", rel)
	tracked.Dir = gitRoot
	if err := tracked.Run(); err != nil {
		return false
	}
	clean := exec.Command("git", "--literal-pathspecs", "diff", "--quiet", "HEAD", "--", rel)
	clean.Dir = gitRoot
	return clean.Run() == nil
}

func entityGitPath(path string) (gitRoot, rel string, ok bool) {
	gitRoot = FindGitRoot(filepath.Dir(path))
	if !hasGitEntry(gitRoot) {
		return "", "", false
	}
	rel, err := filepath.Rel(gitRoot, path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", "", false
	}
	return gitRoot, rel, true
}
