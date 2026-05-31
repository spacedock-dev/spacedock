// ABOUTME: Release-notes core — filter the commit log of workflow-state noise,
// ABOUTME: build the LLM changelog prompt, and gate tag-cutting on confirmation.
package release

import (
	"fmt"
	"strings"
)

// workflowNoisePrefixes are the commit-subject prefixes the release changelog
// must drop: the FO/ensign workflow-state commits (`dispatch:`/`advance:`/
// `merge:`/`archive:`) and the two CI version commits this repo's pipeline
// writes (`release: stamp …`, `next: bump …`). They are matched against the
// commit SUBJECT — the text after the `<sha> ` prefix of a `git log --oneline`
// line — so a real commit whose conventional scope merely mentions a noise word
// (e.g. `fix(dispatch): …`) is kept.
var workflowNoisePrefixes = []string{
	"dispatch:",
	"advance:",
	"merge:",
	"archive:",
	"release: stamp",
	"next: bump",
}

// FilterCommitLog drops the workflow-state commits the changelog prompt is told
// to ignore, so the `claude` input (and the no-claude fallback output) is
// already clean. The input is a `git log --oneline` block (`<sha> <subject>`
// per line); lines whose subject starts with a workflow-noise prefix are
// removed and the remaining lines are returned in order.
func FilterCommitLog(rawLog string) string {
	var kept []string
	for _, line := range strings.Split(rawLog, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if isWorkflowNoise(commitSubject(line)) {
			continue
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n")
}

// commitSubject returns the subject of a `git log --oneline` line — the text
// after the leading short-sha and its space. A line with no space is returned
// unchanged.
func commitSubject(oneline string) string {
	if i := strings.IndexByte(oneline, ' '); i >= 0 {
		return oneline[i+1:]
	}
	return oneline
}

func isWorkflowNoise(subject string) bool {
	for _, p := range workflowNoisePrefixes {
		if strings.HasPrefix(subject, p) {
			return true
		}
	}
	return false
}

// BuildChangelogPrompt returns the captain-confirmed changelog prompt for the
// given release version, with the ignore-list adapted to this repo's
// workflow-state noise. The prompt is reused verbatim from the proven
// scripts/release.sh flow; the parenthetical ignore-list names this repo's
// noise prefixes so the LLM is told to drop them even if FilterCommitLog misses
// a class.
func BuildChangelogPrompt(version string) string {
	return fmt.Sprintf("Summarize these git commits into a release changelog for spacedock v%s. "+
		"Plain text only — no markdown headers, no bold/italic. "+
		"Start with one sentence describing the major theme of this release. "+
		"Then list individual changes as '- ' bullet lines. "+
		"For each bullet, lead with the user value (what upgrading gives you), then briefly describe what changed at a high level. "+
		"Ignore workflow-state commits (dispatch/advance/merge/archive entity commits, the "+
		"`release: stamp` and `next: bump` CI commits, and entity-file changes under docs/dev/.spacedock-state/). "+
		"Group related commits into single entries.", version)
}

// NotesIO injects the side-effecting dependencies of GenerateNotes so the
// generation core stays unit-testable: RawLog yields the `git log` range, and
// Claude runs the summarizer over (prompt, filtered-log) input — returning an
// error when `claude` is unavailable, which triggers the filtered-raw-log
// fallback.
type NotesIO struct {
	RawLog func() (string, error)
	Claude func(prompt, input string) (string, error)
}

// GenerateNotes produces the release notes for version: it resolves the commit
// log, filters the workflow-state noise, and feeds the filtered log to `claude`
// with the captain-confirmed prompt. When `claude` is unavailable (the Claude
// hook returns an error), it falls back to the filtered raw log so the flow
// still emits usable notes rather than failing.
func GenerateNotes(version string, io NotesIO) (string, error) {
	raw, err := io.RawLog()
	if err != nil {
		return "", fmt.Errorf("resolve commit log: %w", err)
	}
	filtered := FilterCommitLog(raw)
	if io.Claude == nil {
		return filtered, nil
	}
	notes, err := io.Claude(BuildChangelogPrompt(version), filtered)
	if err != nil {
		return filtered, nil
	}
	return notes, nil
}

// TagIO injects the confirm-then-tag boundary so the gate stays unit-testable:
// Confirm presents the proposed notes and returns the (possibly captain-edited)
// final body plus whether to proceed; CutTag creates the annotated tag with
// that body.
type TagIO struct {
	Confirm func(proposed string) (body string, ok bool)
	CutTag  func(body string) error
}

// ConfirmAndTag presents the proposed notes for review and cuts the annotated
// tag only on explicit confirmation, using the captain's edited body. A decline
// leaves no tag.
func ConfirmAndTag(proposed string, io TagIO) error {
	body, ok := io.Confirm(proposed)
	if !ok {
		return nil
	}
	return io.CutTag(body)
}
