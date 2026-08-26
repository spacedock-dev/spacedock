// ABOUTME: status --read gate modes — --checklist extracts a gated stage's latest
// ABOUTME: stage-report items with line ranges; --ac-scan cites ACs within them.
package status

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// runReadGate dispatches the --read gate modes: it resolves the ref to a file,
// reads its lines, and runs --checklist or --ac-scan. The two modes are mutually
// exclusive (each emits its own envelope). The stage-report/AC parsing reads the
// raw file lines (1-based) so emitted ranges are Read(offset, lines)-sliceable.
//
// When stage is omitted, it defaults once here to the resolved file's
// frontmatter status field — the entity's current stage — before mode dispatch,
// so the default applies uniformly to --checklist, --ac-scan, text, and --json.
// Explicit --stage is passed through unchanged. A target with no status field to
// default from fails loudly rather than silently emitting.
func runReadGate(roots roots, ref, stage string, checklist, acScan, asJSON bool, stdout, stderr io.Writer) int {
	if checklist && acScan {
		return errExit(stderr, "--checklist and --ac-scan are mutually exclusive")
	}
	path, rc := resolveReadTarget(roots, ref, stderr)
	if rc != 0 {
		return rc
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return errExit(stderr, "cannot read file: "+path)
	}
	stageOmitted := stage == ""
	if stageOmitted {
		stage = parseFrontmatterContent(data)["status"]
		if stage == "" {
			return errExit(stderr, fmt.Sprintf("--stage omitted and %s has no status frontmatter to default from; pass --stage <stage>", path))
		}
	}
	lines := splitLines(string(data))
	if checklist {
		return runChecklist(lines, stage, stageOmitted, asJSON, stdout, stderr)
	}
	return runACScan(lines, stage, stageOmitted, asJSON, stdout, stderr)
}

// stageReportHeadingRe matches a `## Stage Report: <stage> [qualifier...]` line
// and captures the LEADING stage-token — the first whitespace-delimited word
// after the colon, before any `(cycle …)` or freeform qualifier. Interleaved
// real entities append append-only freeform-suffixed sections, so the gated
// stage is selected by this token, never by positional-last.
var stageReportHeadingRe = regexp.MustCompile(`^##\s+Stage Report:\s+(\S+)`)

// checklistBulletRe matches a `- DONE:` / `- SKIPPED:` / `- FAILED:` bullet,
// capturing the status verb and the bullet text after the colon. The verb set is
// the stage-report protocol's fixed three.
var checklistBulletRe = regexp.MustCompile(`^-\s+(DONE|SKIPPED|FAILED):\s*(.*)$`)

// checklistNearMissRe catches status-like bullets that the canonical parser
// intentionally ignores, such as "- DONE (annotation): ...".
var checklistNearMissRe = regexp.MustCompile(`^-\s+(DONE|SKIPPED|FAILED)\b`)

// acHeadingRe matches an `**AC-N**` acceptance-criteria heading and captures the
// AC id. The id is tokenized on the heading boundary (AC- followed by an
// alphanumeric run), never split on `-` (the spike's AC-id boundary finding). An
// asterisk-free trailing label inside the bold span is allowed and discarded, so
// `**AC-1 (VALUE)**` enumerates as `AC-1` exactly as bare `**AC-1**` — the value
// AC the README ideation policy recommends is no longer dropped. The `[^*]*`
// label excludes `*`, so it cannot span a `**` boundary: two headings on one line
// (`**AC-1** … **AC-2**`) still enumerate separately and never merge.
var acHeadingRe = regexp.MustCompile(`\*\*(AC-[0-9A-Za-z]+)[^*]*\*\*`)

// acTokenRe matches an AC-N token anywhere in a line, for citation scanning.
var acTokenRe = regexp.MustCompile(`\bAC-[0-9A-Za-z]+\b`)

// checklistItem is one extracted stage-report bullet: its status verb, the bullet
// text, and the 1-based start/end line range it owns (the bullet line through its
// trailing evidence line, stopping at the next bullet, a sub-heading, or EOF of
// the section).
type checklistItem struct {
	status string
	text   string
	start  int // 1-based
	end    int // 1-based, inclusive
}

type checklistNearMiss struct {
	status string
	line   int
	text   string
}

func firstChecklistNearMiss(lines []string, start, end int) *checklistNearMiss {
	for line := start; line <= end; line++ {
		text := lines[line-1]
		if strings.HasPrefix(text, "### ") {
			break
		}
		match := checklistNearMissRe.FindStringSubmatch(text)
		if match != nil && !checklistBulletRe.MatchString(text) {
			return &checklistNearMiss{status: match[1], line: line, text: strings.TrimSpace(text)}
		}
	}
	return nil
}

// selectStageReport returns the line range [start,end] (1-based, inclusive) of
// the LATEST stage-report section whose leading stage-token equals stage, across
// interleaved sections. Within a stage's multiple cycles the append-only ordering
// makes the last-positional-of-that-stage the latest cycle. Returns ok=false when
// no stage-report heading matches — the caller fails loudly (never a silent
// positional-last fallback).
func selectStageReport(lines []string, stage string) (start, end int, ok bool) {
	// Heading line indices (0-based) of every ## Stage Report whose token == stage,
	// plus the index of every level-2 heading, so a selected section ends at the
	// next level-2 boundary.
	type hdr struct {
		idx   int // 0-based line index
		token string
	}
	var reports []hdr
	var level2 []int
	for i, line := range lines {
		if strings.HasPrefix(line, "## ") {
			level2 = append(level2, i)
		}
		if m := stageReportHeadingRe.FindStringSubmatch(line); m != nil {
			reports = append(reports, hdr{idx: i, token: m[1]})
		}
	}
	chosen := -1
	for _, r := range reports {
		if r.token == stage {
			chosen = r.idx // keep walking: last match wins (latest cycle)
		}
	}
	if chosen < 0 {
		return 0, 0, false
	}
	// The section runs from its heading to the line before the next level-2
	// heading (any `## `, including the next Stage Report or Acceptance criteria),
	// or to EOF for the last section.
	endIdx := len(lines) - 1
	for _, h := range level2 {
		if h > chosen {
			endIdx = h - 1
			break
		}
	}
	return chosen + 1, endIdx + 1, true
}

// noStageReportDiagnostic formats the "no ## Stage Report for stage" error. When
// stage was defaulted from the current status (--stage omitted), it names the
// stage as the current status so the defaulting is transparent to the caller;
// an explicit --stage keeps the plain wording.
func noStageReportDiagnostic(stage string, stageOmitted bool) string {
	if stageOmitted {
		return fmt.Sprintf("no ## Stage Report for stage %q (current status; --stage omitted) in this file", stage)
	}
	return fmt.Sprintf("no ## Stage Report for stage %q in this file", stage)
}

// extractChecklist parses the DONE/SKIPPED/FAILED items within the stage-report
// section [start,end] (1-based inclusive). Each item owns from its bullet line
// through the line before the next bullet or sub-heading (### …) within the
// section; the last item runs to the section end. A `### ` sub-heading (e.g.
// `### Summary`, `### Feedback Cycles`) terminates the checklist — items after it
// are not part of the roll-up.
func extractChecklist(lines []string, start, end int) []checklistItem {
	var items []checklistItem
	// Find the bullet lines and the first sub-heading (the checklist terminator).
	limit := end // 1-based inclusive section end
	for i := start; i <= end; i++ {
		line := lines[i-1]
		if strings.HasPrefix(line, "### ") {
			limit = i - 1 // checklist ends at the line before the sub-heading
			break
		}
	}
	var bulletLines []int // 1-based bullet line numbers
	for i := start; i <= limit; i++ {
		if checklistBulletRe.MatchString(lines[i-1]) {
			bulletLines = append(bulletLines, i)
		}
	}
	for n, bl := range bulletLines {
		m := checklistBulletRe.FindStringSubmatch(lines[bl-1])
		itemEnd := limit
		if n+1 < len(bulletLines) {
			itemEnd = bulletLines[n+1] - 1
		}
		// Trim trailing blank lines so the range owns the bullet through its last
		// non-blank evidence line, not the separator blanks before the next bullet
		// or sub-heading.
		for itemEnd > bl && strings.TrimSpace(lines[itemEnd-1]) == "" {
			itemEnd--
		}
		items = append(items, checklistItem{
			status: m[1],
			text:   strings.TrimSpace(m[2]),
			start:  bl,
			end:    itemEnd,
		})
	}
	return items
}

// checklistJSON renders the --checklist envelope: the selected stage and its
// items (status/text/1-based start/end), every leaf a string.
func checklistJSON(stage string, items []checklistItem) *jsonObj {
	arr := make(jsonArr, 0, len(items))
	for _, it := range items {
		arr = append(arr, newJSONObj().
			set("status", it.status).
			set("text", it.text).
			set("start", strconv.Itoa(it.start)).
			set("end", strconv.Itoa(it.end)))
	}
	return newJSONObj().
		set("command", "read").
		set("stage", stage).
		setValue("checklist", arr)
}

// runChecklist handles --read <entity> [--stage X] --checklist. It selects the
// latest stage-report section for X (or, when --stage was omitted, the caller's
// defaulted current-status stage) across interleaved sections and emits its
// checklist items with line ranges; a stage matching no report fails loudly
// (non-zero exit, named diagnostic), never a silent emit. When stage was
// defaulted, the diagnostic names it as the current status so the defaulting is
// transparent.
func runChecklist(lines []string, stage string, stageOmitted, asJSON bool, stdout, stderr io.Writer) int {
	start, end, ok := selectStageReport(lines, stage)
	if !ok {
		return errExit(stderr, noStageReportDiagnostic(stage, stageOmitted))
	}
	items := extractChecklist(lines, start, end)
	if asJSON {
		emitJSON(stdout, checklistJSON(stage, items))
		return 0
	}
	fmt.Fprintf(stdout, "stage=%s\n", stage)
	for _, it := range items {
		fmt.Fprintf(stdout, "status=%s start=%d end=%d text=%s\n", it.status, it.start, it.end, it.text)
	}
	return 0
}

// acCitation is one AC-token mention found within a checklist item's line range.
type acCitation struct {
	line int // 1-based
	text string
}

// acEvidence is an acceptance criterion with its heading line and the citations
// found within the gated stage's checklist line ranges. Unevidenced is true when
// the citation list is empty.
type acEvidence struct {
	id        string
	line      int // 1-based heading line
	citations []acCitation
}

// findAcceptanceCriteria returns the 1-based [start,end] range of the
// `## Acceptance criteria` section, or ok=false when absent. The section runs to
// the next level-2 heading or EOF.
func findAcceptanceCriteria(lines []string) (start, end int, ok bool) {
	headingIdx := -1
	var level2 []int
	for i, line := range lines {
		if strings.HasPrefix(line, "## ") {
			level2 = append(level2, i)
		}
		if headingIdx < 0 && strings.EqualFold(strings.TrimSpace(line), "## Acceptance criteria") {
			headingIdx = i
		}
	}
	if headingIdx < 0 {
		return 0, 0, false
	}
	endIdx := len(lines) - 1
	for _, h := range level2 {
		if h > headingIdx {
			endIdx = h - 1
			break
		}
	}
	return headingIdx + 1, endIdx + 1, true
}

// scanACs collects the AC headings in the acceptance-criteria section and, for
// each, every AC-token citation found within the checklist items' line ranges
// (the scope boundary the spike settled: the citation search is confined to
// checklist evidence lines, never the Summary / Feedback Cycles / reviewer
// prose, so a reviewer complaining "AC-N has no evidence" is not counted as AC-N
// being evidenced).
func scanACs(lines []string, acStart, acEnd int, items []checklistItem) []acEvidence {
	var acs []acEvidence
	seen := map[string]bool{}
	for i := acStart; i <= acEnd; i++ {
		for _, m := range acHeadingRe.FindAllStringSubmatch(lines[i-1], -1) {
			id := m[1]
			if seen[id] {
				continue
			}
			seen[id] = true
			acs = append(acs, acEvidence{id: id, line: i})
		}
	}
	for idx := range acs {
		id := acs[idx].id
		for _, it := range items {
			for ln := it.start; ln <= it.end; ln++ {
				for _, tok := range acTokenRe.FindAllString(lines[ln-1], -1) {
					if tok == id {
						acs[idx].citations = append(acs[idx].citations, acCitation{line: ln, text: lines[ln-1]})
					}
				}
			}
		}
	}
	return acs
}

// acScanJSON renders the --ac-scan envelope: the stage and per-AC citation map
// with the unevidenced flag. It emits NO satisfied and NO natural_place key —
// whether a citation is evidence-of-satisfaction and whether an unevidenced AC's
// natural place is here are judgment that routes to level-3.
func acScanJSON(stage string, acs []acEvidence) *jsonObj {
	arr := make(jsonArr, 0, len(acs))
	for _, ac := range acs {
		cites := make(jsonArr, 0, len(ac.citations))
		for _, c := range ac.citations {
			cites = append(cites, newJSONObj().
				set("line", strconv.Itoa(c.line)).
				set("text", c.text))
		}
		o := newJSONObj().
			set("id", ac.id).
			set("line", strconv.Itoa(ac.line)).
			set("unevidenced", strconv.FormatBool(len(ac.citations) == 0))
		o.setValue("citations", cites)
		arr = append(arr, o)
	}
	return newJSONObj().
		set("command", "read").
		set("stage", stage).
		setValue("acs", arr)
}

// runACScan handles --read <entity> [--stage X] --ac-scan. It cites each AC's
// evidence from the gated stage's checklist line ranges only and flags the
// unevidenced ACs. A stage matching no report, or an absent ## Acceptance
// criteria section, fails loudly, never a silent emit. When stage was omitted
// and defaulted by the caller, the no-report diagnostic names it as the current
// status.
func runACScan(lines []string, stage string, stageOmitted, asJSON bool, stdout, stderr io.Writer) int {
	start, end, ok := selectStageReport(lines, stage)
	if !ok {
		return errExit(stderr, noStageReportDiagnostic(stage, stageOmitted))
	}
	acStart, acEnd, acOK := findAcceptanceCriteria(lines)
	if !acOK {
		return errExit(stderr, "no ## Acceptance criteria section in this file")
	}
	items := extractChecklist(lines, start, end)
	acs := scanACs(lines, acStart, acEnd, items)
	if asJSON {
		emitJSON(stdout, acScanJSON(stage, acs))
		return 0
	}
	fmt.Fprintf(stdout, "stage=%s\n", stage)
	for _, ac := range acs {
		fmt.Fprintf(stdout, "ac=%s line=%d unevidenced=%t citations=%d\n", ac.id, ac.line, len(ac.citations) == 0, len(ac.citations))
	}
	return 0
}
