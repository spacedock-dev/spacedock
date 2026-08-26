// ABOUTME: Entered working stages remain their own dispatch target until their
// ABOUTME: latest stage report is structurally complete and committed cleanly.
package status

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// enteredStageAwaitingCompletion reports whether a current working stage must
// be dispatched before ordinary successor projection resumes. Initial stages
// and gate, terminal, and worktree suppression are handled by the caller.
func enteredStageAwaitingCompletion(e *entity, stage Stage) bool {
	return enteredStageCompletionFailure(e, stage) != nil
}

func enteredStageCompletionFailure(e *entity, stage Stage) *completionFailure {
	if stage.initial || stage.gate || stage.terminal {
		return nil
	}
	return completionFailureForPath(e.path, stage.Name)
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
	return completionFailureForPath(path, stage) == nil
}

type completionFailureKind uint8

const (
	completionMissingReport completionFailureKind = iota
	completionIncompleteReport
	completionUntracked
	completionDirty
	completionInspectionFailed
)

type completionFailure struct {
	kind    completionFailureKind
	line    int
	item    string
	path    string
	gitRoot string
	detail  string
}

func (f *completionFailure) diagnostic(stage string) string {
	switch f.kind {
	case completionMissingReport:
		return fmt.Sprintf("missing current-stage report: no heading whose first stage token is %q; add %q", stage, "## Stage Report: "+stage)
	case completionIncompleteReport:
		return "incomplete current-stage report: " + f.detail
	case completionUntracked:
		return fmt.Sprintf("untracked completion artifact: %s is not tracked in local Git root %s; add and commit that path", f.path, f.gitRoot)
	case completionDirty:
		return fmt.Sprintf("dirty completion artifact: %s differs from local HEAD in Git root %s; commit that path (this guard does not require a remote push)", f.path, f.gitRoot)
	default:
		return "unable to inspect completion artifact: " + f.detail
	}
}

func completionFailureForPath(path, stage string) *completionFailure {
	data, err := os.ReadFile(path)
	if err != nil {
		return &completionFailure{kind: completionInspectionFailed, detail: fmt.Sprintf("cannot read %s: %v", path, err)}
	}
	if failure := stageReportFailure(data, stage); failure != nil {
		return failure
	}
	return entityGitFailure(path)
}

func hasCompleteStageReport(data []byte, stage string) bool {
	return stageReportFailure(data, stage) == nil
}

func stageReportFailure(data []byte, stage string) *completionFailure {
	lines := splitLines(string(data))
	start, end, ok := selectStageReport(lines, stage)
	if !ok {
		return &completionFailure{kind: completionMissingReport}
	}
	items := extractChecklist(lines, start, end)
	var first *completionFailure
	for _, item := range items {
		detail := ""
		switch {
		case item.status == "FAILED":
			detail = fmt.Sprintf("line %d FAILED item %q is unresolved", item.start, item.text)
		case strings.TrimSpace(item.text) == "":
			detail = fmt.Sprintf("line %d %s item has blank text", item.start, item.status)
		case !checklistItemHasEvidence(lines, item):
			detail = fmt.Sprintf("line %d %s item %q has no evidence or rationale line", item.start, item.status, item.text)
		}
		if detail != "" {
			first = &completionFailure{kind: completionIncompleteReport, line: item.start, item: item.text, detail: detail}
			break
		}
	}
	if nearMiss := firstChecklistNearMiss(lines, start, end); nearMiss != nil && (first == nil || nearMiss.line < first.line) {
		first = &completionFailure{
			kind: completionIncompleteReport, line: nearMiss.line, item: nearMiss.text,
			detail: fmt.Sprintf("line %d %q is not canonical; use %q", nearMiss.line, nearMiss.text, "- "+nearMiss.status+": <item text>"),
		}
	}
	if first != nil {
		return first
	}
	if len(items) == 0 {
		return &completionFailure{kind: completionIncompleteReport, detail: "no recognized checklist items; add a canonical - DONE: or - SKIPPED: item"}
	}
	summaryLine, found, complete := stageReportSummary(lines, start, end)
	if !found {
		return &completionFailure{kind: completionIncompleteReport, detail: `missing non-empty "### Summary"`}
	}
	if !complete {
		return &completionFailure{kind: completionIncompleteReport, line: summaryLine, detail: fmt.Sprintf("line %d %q has no content", summaryLine, "### Summary")}
	}
	return nil
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
	_, _, complete := stageReportSummary(lines, start, end)
	return complete
}

func stageReportSummary(lines []string, start, end int) (summaryLine int, found, complete bool) {
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
				return line, true, true
			}
		}
		return line, true, false
	}
	return 0, false, false
}

// entityPathCleanInHEAD is deliberately literal and path-scoped. A tracked
// entity whose working-tree/index bytes differ from HEAD is not durable proof;
// an unrelated dirty sibling does not affect the answer.
func entityPathCleanInHEAD(path string) bool {
	return entityGitFailure(path) == nil
}

func entityGitFailure(path string) *completionFailure {
	gitRoot, rel, ok := entityGitPath(path)
	if !ok {
		return &completionFailure{kind: completionInspectionFailed, path: path, detail: fmt.Sprintf("%s is not inside a local Git worktree", path)}
	}
	rel = filepath.ToSlash(rel)
	tracked := exec.Command("git", "--literal-pathspecs", "ls-files", "--error-unmatch", "--", rel)
	tracked.Dir = gitRoot
	if err := tracked.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 1 {
			return &completionFailure{kind: completionUntracked, path: path, gitRoot: gitRoot}
		}
		return &completionFailure{kind: completionInspectionFailed, path: path, gitRoot: gitRoot, detail: fmt.Sprintf("git ls-files failed in %s: %v", gitRoot, err)}
	}
	clean := exec.Command("git", "--literal-pathspecs", "diff", "--quiet", "HEAD", "--", rel)
	clean.Dir = gitRoot
	err := clean.Run()
	if err == nil {
		return nil
	}
	if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 1 {
		return &completionFailure{kind: completionDirty, path: path, gitRoot: gitRoot}
	}
	return &completionFailure{kind: completionInspectionFailed, path: path, gitRoot: gitRoot, detail: fmt.Sprintf("git diff failed in %s: %v", gitRoot, err)}
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
