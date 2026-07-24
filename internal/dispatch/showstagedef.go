// ABOUTME: show-stage-def extracts a workflow README's ### {stage} subsection,
// ABOUTME: matching the oracle's extract_stage_subsection + cmd_show_stage_def.
package dispatch

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// headingDecorationChars are the inline-markdown decoration characters stripped
// from a heading before tokenizing, matching _HEADING_DECORATION_CHARS.
const headingDecorationChars = "`*_~"

// headingTokens returns the content tokens of a `### `-prefixed heading, or nil
// when the line is not such a heading. Strips inline decoration (“ ` “, `*`,
// `_`, `~`) and treats `(` and `[` as token terminators so trailing annotations
// like `*(captain-interactive)*` or `[terminal]` do not merge with the name.
// Matches _heading_tokens.
func headingTokens(line string) []string {
	stripped := strings.TrimSpace(line)
	if !strings.HasPrefix(stripped, "### ") {
		return nil
	}
	rest := strings.TrimSpace(stripped[4:])
	rest = stripDecoration(rest)
	rest = strings.ReplaceAll(rest, "(", " ")
	rest = strings.ReplaceAll(rest, "[", " ")
	return strings.Fields(rest)
}

// stripDecoration removes every decoration character from s (str.translate with
// a delete table).
func stripDecoration(s string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(headingDecorationChars, r) {
			return -1
		}
		return r
	}, s)
}

// headingFirstToken returns the first content token of a `### ` heading, or ""
// when the line is not a heading or carries no tokens. Matches
// _heading_first_token (which returns None; "" is the no-match sentinel here).
func headingFirstToken(line string) string {
	tokens := headingTokens(line)
	if len(tokens) == 0 {
		return ""
	}
	return tokens[0]
}

// stageHeadingError is the parser diagnostic raised when a `###` heading
// mentions the stage name as a token but the stage name is not the first token.
// It carries the oracle's ValueError text so cmd-level wrappers can prefix it.
type stageHeadingError struct{ msg string }

func (e *stageHeadingError) Error() string { return e.msg }

type sourceSpan struct{ start, end int }

// extractStageSubsection returns the full ### {stage} subsection from a workflow
// README. Heading match is permissive: any `###` line whose first content token
// (after stripping “ ` “, `*`, `_`, `~` and treating `(` / `[` as token
// terminators) equals stage is a match. When no permissive match is found but
// some `###` line mentions the stage name as a token, a *stageHeadingError is
// returned surfacing the malformed heading. Returns ("", nil) only when no
// `###` line mentions the stage name at all (genuinely missing stage). Matches
// extract_stage_subsection.
func extractStageSubsection(readmePath, stage string) (string, error) {
	data, err := os.ReadFile(readmePath)
	if err != nil {
		return "", err
	}
	subsection, _, err := extractStageSubsectionBytes(data, stage)
	return subsection, err
}

// extractStageSubsectionBytes preserves the legacy rendered subsection while
// also returning its structural raw source span from the same immutable buffer.
func extractStageSubsectionBytes(data []byte, stage string) (string, sourceSpan, error) {
	lines := splitTextLines(string(data))
	starts := textLineStarts(data)

	start := -1
	for i, line := range lines {
		if headingFirstToken(line) == stage {
			start = i
			break
		}
	}
	if start < 0 {
		for i, line := range lines {
			tokens := headingTokens(line)
			if len(tokens) > 0 && containsToken(tokens, stage) {
				stripped := strings.TrimSpace(line)
				return "", sourceSpan{}, &stageHeadingError{msg: fmt.Sprintf(
					"stage heading at line %d mentions '%s' "+
						"but does not parse as a stage heading: %s. "+
						"The stage name must be the first content token of the "+
						"heading after stripping Markdown decoration "+
						"(backticks, *, _, ~) and treating '(' and '[' as "+
						"token terminators.",
					i+1, stage, pyRepr(stripped))}
			}
		}
		return "", sourceSpan{}, nil
	}

	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		stripped := strings.TrimSpace(lines[i])
		if strings.HasPrefix(stripped, "### ") || strings.HasPrefix(stripped, "## ") {
			end = i
			break
		}
	}
	rawEnd := len(data)
	if end < len(starts) {
		rawEnd = starts[end]
	}
	renderEnd := end
	for renderEnd > start && strings.TrimSpace(lines[renderEnd-1]) == "" {
		renderEnd--
	}
	return strings.Join(lines[start:renderEnd], "\n"),
		sourceSpan{start: starts[start], end: rawEnd}, nil
}

// containsToken reports whether token is in tokens, matching Python's
// `stage_name in tokens` membership test.
func containsToken(tokens []string, token string) bool {
	for _, t := range tokens {
		if t == token {
			return true
		}
	}
	return false
}

// lineBoundary reports whether r is one of the boundaries Python str.splitlines()
// breaks on: LF, CR (and CRLF, handled by the scanner), VT, FF, FS, GS, RS, NEL
// (U+0085), LS (U+2028), PS (U+2029). The oracle reads the README in universal-
// newline text mode then calls splitlines(); the CR-family translation is
// subsumed here since splitlines() itself breaks on CR/CRLF, so a direct
// splitlines() gives the identical result.
func lineBoundary(r rune) bool {
	switch r {
	case '\n', '\r', '\v', '\f', '\x1c', '\x1d', '\x1e', '\u0085', '\u2028', '\u2029':
		return true
	}
	return false
}

// splitTextLines splits text into lines exactly as Python's str.splitlines()
// does: a boundary terminates the current line (CRLF counts as one boundary),
// and the trailing line is dropped when the text ends in a boundary (no empty
// final element). An empty input yields no lines.
func splitTextLines(text string) []string {
	if text == "" {
		return nil
	}
	var lines []string
	var cur []rune
	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if !lineBoundary(r) {
			cur = append(cur, r)
			continue
		}
		lines = append(lines, string(cur))
		cur = nil
		// CRLF is a single boundary: consume the following LF after a CR.
		if r == '\r' && i+1 < len(runes) && runes[i+1] == '\n' {
			i++
		}
	}
	if len(cur) > 0 {
		lines = append(lines, string(cur))
	}
	return lines
}

// textLineStarts maps splitTextLines' full separator set to source byte
// coordinates. The manual decoder consumes CRLF atomically, avoiding the
// double-advance bug a range-based mapper introduces.
func textLineStarts(data []byte) []int {
	if len(data) == 0 {
		return nil
	}
	starts := []int{0}
	for i := 0; i < len(data); {
		r, width := utf8.DecodeRune(data[i:])
		next := i + width
		if lineBoundary(r) {
			if r == '\r' && next < len(data) && data[next] == '\n' {
				next++
			}
			if next < len(data) {
				starts = append(starts, next)
			}
		}
		i = next
	}
	return starts
}

// resolveStageContext assembles the mandatory legacy stage subsection followed
// by declared fence-safe sections, all resolved from one README buffer.
func resolveStageContext(data []byte, stage string) (string, error) {
	subsection, stageSpan, err := extractStageSubsectionBytes(data, stage)
	if err != nil || subsection == "" {
		return subsection, err
	}
	selectors, err := status.StageContextSections(data, stage)
	if err != nil {
		return "", err
	}
	if len(selectors) == 0 {
		return subsection, nil
	}
	seen := map[string]bool{}
	for _, selector := range selectors {
		if seen[selector] {
			return "", fmt.Errorf("repeated selector %q", selector)
		}
		seen[selector] = true
	}
	sections, err := status.FindSectionSpans(data, selectors)
	if err != nil {
		return "", err
	}
	for i, section := range sections {
		span := sourceSpan{section.Start, section.End}
		if spansIntersect(stageSpan, span) {
			return "", overlapError("stage "+stage, stageSpan, section.Heading, span)
		}
		for j := 0; j < i; j++ {
			other := sourceSpan{sections[j].Start, sections[j].End}
			if spansIntersect(other, span) {
				return "", overlapError(sections[j].Heading, other, section.Heading, span)
			}
		}
	}

	parts := []string{subsection}
	for _, section := range sections {
		lines := splitTextLines(string(data[section.Start:section.End]))
		for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
			lines = lines[:len(lines)-1]
		}
		parts = append(parts, strings.Join(lines, "\n"))
	}
	return strings.Join(parts, "\n\n"), nil
}

func spansIntersect(a, b sourceSpan) bool {
	return a.start < b.end && b.start < a.end
}

func overlapError(aName string, a sourceSpan, bName string, b sourceSpan) error {
	return fmt.Errorf("overlap: %q [%d,%d) intersects %q [%d,%d)",
		aName, a.start, a.end, bName, b.start, b.end)
}

// runShowStageDef emits the README's ### {stage} subsection on stdout. Exit 0
// with stdout on success; exit 1 with a parser diagnostic on a malformed or
// missing heading or an unresolvable workflow dir / README. Matches
// cmd_show_stage_def.
func runShowStageDef(workflowDir, stage string, stdout, stderr io.Writer) int {
	if info, err := os.Stat(workflowDir); err != nil || !info.IsDir() {
		fmt.Fprintf(stderr, "error: workflow directory not found: %s\n", workflowDir)
		return 1
	}
	readmePath := filepath.Join(workflowDir, "README.md")
	if !isFile(readmePath) {
		fmt.Fprintf(stderr, "error: workflow README not found at '%s'\n", readmePath)
		return 1
	}

	data, err := os.ReadFile(readmePath)
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 1
	}
	subsection, err := resolveStageContext(data, stage)
	if err != nil {
		if she, ok := err.(*stageHeadingError); ok {
			fmt.Fprintf(stderr, "error: %s\n", she.msg)
			return 1
		}
		fmt.Fprintf(stderr, "error: workflow README %q, stage %q: %s\n", readmePath, stage, err)
		return 1
	}
	if subsection == "" {
		fmt.Fprintf(stderr, "error: stage '%s' heading not found in %s\n", stage, readmePath)
		return 1
	}
	fmt.Fprintln(stdout, subsection)
	return 0
}
