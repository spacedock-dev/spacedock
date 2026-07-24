// ABOUTME: --read section helper — parses a markdown file's frontmatter plus an
// ABOUTME: ordered ATX-heading map (text/level/1-based offset/section line count).
package status

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// heading is one ATX heading found in a file's body: its text, level, 1-based
// line offset (identical to Read(offset, …) semantics), and the section's line
// count (offset through the line before the next heading of level <= its own;
// the final heading runs to EOF). text/level/offset/lines map directly to the
// --read JSON object's four string fields.
type heading struct {
	text   string
	level  int
	offset int // 1-based line number of the heading line
	lines  int // section line count, so Read(offset, lines) returns exactly it
	start  int // raw half-open source span, populated by scanHeadingSpans
	end    int
}

// SectionSpan is one selected fence-safe heading section in raw README coordinates.
type SectionSpan struct {
	Heading    string
	Start, End int
}

// sectionRead is the parsed result of reading a markdown file: its frontmatter,
// the ordered heading map, and the file's total line count (the exact append
// offset a caller adds a new trailing section after).
type sectionRead struct {
	frontmatter map[string]string
	headings    []heading
	totalLines  int
}

// readSections parses path into its frontmatter plus an ordered heading map.
// Frontmatter reuses ParseFrontmatter verbatim. Headings are every ATX heading
// (`#`..`######` followed by a space) in the body, with fenced-code-block lines
// skipped. Returns (sectionRead, true) on a readable file, (zero, false) when
// the file cannot be read.
func readSections(path string) (sectionRead, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return sectionRead{}, false
	}
	lines := splitLines(string(data))
	headings := scanHeadings(lines)
	return sectionRead{
		frontmatter: ParseFrontmatter(path),
		headings:    headings,
		totalLines:  len(lines),
	}, true
}

// scanHeadings returns the ordered ATX headings in lines, with each heading's
// section line count resolved by the level-aware ownership rule: a heading owns
// from its line to the line before the next heading of level <= its own; the
// final heading runs to the last line. Lines inside a fenced code block
// (``` ``` ``` / `~~~`) are not eligible to be headings. lines is the splitLines
// view, so a heading's `offset` is its 1-based index (i+1).
func scanHeadings(lines []string) []heading {
	var headings []heading
	fence := ""
	for i, line := range lines {
		if marker := codeFenceMarker(line); marker != "" {
			switch {
			case fence == "":
				fence = marker
			case strings.HasPrefix(strings.TrimSpace(line), fence):
				fence = ""
			}
			continue
		}
		if fence != "" {
			continue
		}
		level, text, ok := parseATXHeading(line)
		if !ok {
			continue
		}
		headings = append(headings, heading{text: text, level: level, offset: i + 1})
	}
	assignSectionLines(headings, len(lines))
	return headings
}

// FindSectionSpans resolves selectors in declaration order against the existing
// fence-safe heading scanner. Missing and ambiguous exact heading text fail.
func FindSectionSpans(data []byte, selectors []string) ([]SectionSpan, error) {
	headings := scanHeadingSpans(data)
	result := make([]SectionSpan, 0, len(selectors))
	for _, selector := range selectors {
		var matches []heading
		for _, h := range headings {
			if h.text == selector {
				matches = append(matches, h)
			}
		}
		if len(matches) != 1 {
			return nil, fmt.Errorf("selector %q matches %d headings", selector, len(matches))
		}
		h := matches[0]
		result = append(result, SectionSpan{Heading: h.text, Start: h.start, End: h.end})
	}
	return result, nil
}

func scanHeadingSpans(data []byte) []heading {
	headings := scanHeadings(splitLines(string(data)))
	starts := rawLineStarts(data)
	for i := range headings {
		startLine := headings[i].offset - 1
		endLine := startLine + headings[i].lines
		headings[i].start = starts[startLine]
		headings[i].end = len(data)
		if endLine < len(starts) {
			headings[i].end = starts[endLine]
		}
	}
	return headings
}

// rawLineStarts maps splitLines' CR/LF line model back to source bytes. CRLF is
// one boundary and an EOF boundary does not create a trailing logical line.
func rawLineStarts(data []byte) []int {
	if len(data) == 0 {
		return nil
	}
	starts := []int{0}
	for i := 0; i < len(data); i++ {
		if data[i] != '\r' && data[i] != '\n' {
			continue
		}
		if data[i] == '\r' && i+1 < len(data) && data[i+1] == '\n' {
			i++
		}
		if i+1 < len(data) {
			starts = append(starts, i+1)
		}
	}
	return starts
}

// assignSectionLines fills each heading's `lines` by the level-aware ownership
// rule: a heading runs to the line before the next heading of level <= its own,
// and the final-by-ownership heading runs to totalLines. Operates in place over
// the offset-ordered slice.
func assignSectionLines(headings []heading, totalLines int) {
	for i := range headings {
		end := totalLines // default: to EOF
		for j := i + 1; j < len(headings); j++ {
			if headings[j].level <= headings[i].level {
				end = headings[j].offset - 1
				break
			}
		}
		headings[i].lines = end - headings[i].offset + 1
	}
}

// parseATXHeading reports whether line is an ATX heading (1-6 leading `#`
// followed by a space), returning its level and trimmed heading text. The text
// is the content after the `# ` marker with surrounding whitespace and any
// trailing closing `#` run stripped, matching the rendered heading text.
func parseATXHeading(line string) (level int, text string, ok bool) {
	hashes := 0
	for hashes < len(line) && line[hashes] == '#' {
		hashes++
	}
	if hashes == 0 || hashes > 6 {
		return 0, "", false
	}
	if hashes >= len(line) || line[hashes] != ' ' {
		return 0, "", false
	}
	body := strings.TrimSpace(line[hashes+1:])
	body = strings.TrimRight(body, "#")
	body = strings.TrimRight(body, " ")
	return hashes, body, true
}

// runReadSection handles --read <ref-or-path>: it resolves the argument to a
// file (an existing filesystem path is read directly; otherwise the argument is
// resolved as an entity reference the same way --resolve resolves it), parses
// the file into frontmatter + heading map, and emits the JSON envelope (--json)
// or the key=value text mirror. fields, when non-nil, projects the frontmatter
// object to exactly those keys in order (the same --fields semantics --where /
// --next apply). Returns the exit code.
func runReadSection(roots roots, ref string, asJSON bool, fields []string, stdout, stderr io.Writer) int {
	path, rc := resolveReadTarget(roots, ref, stderr)
	if rc != 0 {
		return rc
	}
	sr, ok := readSections(path)
	if !ok {
		return errExit(stderr, "cannot read file: "+path)
	}
	if asJSON {
		emitJSON(stdout, readJSON(path, sr, fields))
		return 0
	}
	fmt.Fprint(stdout, formatReadText(path, sr))
	return 0
}

// resolveReadTarget turns --read's argument into a file path: an existing
// regular file is used verbatim, so the README and any report-shaped markdown
// are addressable without being a tracked entity; otherwise the argument is
// resolved as an entity reference (the --resolve resolver), and an
// unknown/ambiguous ref fails with that resolver's error shape. Returns
// (path, 0) on success or ("", 1) after printing the error.
func resolveReadTarget(roots roots, ref string, stderr io.Writer) (string, int) {
	if isRegularFile(ref) {
		return ref, 0
	}
	idStyle, err := workflowIDStyle(roots.definitionDir)
	if err != nil {
		return "", errExit(stderr, err.Error())
	}
	res := resolveReferenceCandidates(roots.definitionDir, roots.entityDir, ref, true, idStyle, stderr)
	if res.status != "ok" {
		for _, e := range res.errors {
			fmt.Fprintln(stderr, e)
		}
		return "", 1
	}
	return res.matches[0].path, 0
}

// formatReadText renders the key=value text mirror of a --read result: a header
// line carrying the realpath'd path and total line count, then one line per
// heading with its level, offset, lines, and text.
func formatReadText(path string, sr sectionRead) string {
	var b strings.Builder
	fmt.Fprintf(&b, "path=%s total_lines=%d\n", realpathOf(path), sr.totalLines)
	for _, h := range sr.headings {
		fmt.Fprintf(&b, "level=%d offset=%d lines=%d text=%s\n", h.level, h.offset, h.lines, h.text)
	}
	for _, s := range parseStagesBlock(path) {
		fmt.Fprintf(&b, "stage=%s worktree=%t gate=%t terminal=%t initial=%t", s.Name, s.Worktree, s.gate, s.terminal, s.initial)
		for _, of := range []string{"feedback-to", "agent", "fresh", "model"} {
			if v, ok := s.optional[of]; ok {
				fmt.Fprintf(&b, " %s=%s", of, v)
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

// codeFenceMarker returns the fence marker (``` ``` ``` or `~~~`) when line opens
// or closes a fenced code block, else "". A fence is a line whose first
// non-whitespace run is three or more backticks or `~`; the returned marker is the
// three-char prefix so a closing fence is matched by the same opening marker.
func codeFenceMarker(line string) string {
	trimmed := strings.TrimSpace(line)
	for _, ch := range []byte{'`', '~'} {
		run := 0
		for run < len(trimmed) && trimmed[run] == ch {
			run++
		}
		if run >= 3 {
			return strings.Repeat(string(ch), 3)
		}
	}
	return ""
}
