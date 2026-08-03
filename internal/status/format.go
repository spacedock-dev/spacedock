// ABOUTME: Output formatters matching print_status_table / print_next_table /
// ABOUTME: print_boot — fixed-width columns, sort keys, and ellipsis truncation.
package status

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

var (
	defaultStatusFields = []string{"id", "slug", "status", "title", "score"}
	defaultNextFields   = []string{"id", "slug", "status"}
)

// unknownStageOrder is the stageOrder result for a status the workflow's
// stages: block does not declare.
const unknownStageOrder = 99

// stageOrder maps a status to its 1-based stage order; unknown statuses get
// unknownStageOrder. Matches stage_order.
func stageOrder(status string, stages []Stage) int {
	for i, s := range stages {
		if s.Name == status {
			return i + 1
		}
	}
	return unknownStageOrder
}

// scoreSortVal returns the comparable score component: -score for numeric
// scores, 0 for non-numeric, 1 for empty (sorts last). Matches the score
// handling shared by sort_key_default / sort_key_next.
func scoreSortVal(score string) float64 {
	if score != "" {
		if f, ok := pythonFloat(score); ok {
			return -f
		}
		return 0
	}
	return 1
}

// pythonFloat parses score with Python float() semantics for the sort key,
// returning (value, true) on success and (0, false) on what Python rejects.
// Go's strconv.ParseFloat accepts hex-float notation (0x1p4 -> 16) that Python
// float() rejects (-> ValueError); for the reachable, TrimSpace'd ASCII score
// values the two parsers otherwise agree, so rejecting the 0x/0X prefix is the
// one correction needed to order identically to the oracle.
func pythonFloat(score string) (float64, bool) {
	body := score
	if len(body) > 0 && (body[0] == '+' || body[0] == '-') {
		body = body[1:]
	}
	if strings.HasPrefix(body, "0x") || strings.HasPrefix(body, "0X") {
		return 0, false
	}
	f, err := strconv.ParseFloat(score, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

// sortDefault sorts entities by (stage order descending among known stages,
// score descending). A stage unknown to the workflow's stages: block
// (unknownStageOrder) is grouped after every known stage regardless of the
// descending direction, so a malformed/unmodeled status stays visible without
// displacing active work closer to completion. A stable sort preserves the
// discovery (slug) order for equal keys.
func sortDefault(entities []*entity, stages []Stage) []*entity {
	out := append([]*entity(nil), entities...)
	sort.SliceStable(out, func(i, j int) bool {
		oi, oj := stageOrder(out[i].fields["status"], stages), stageOrder(out[j].fields["status"], stages)
		iUnknown, jUnknown := oi == unknownStageOrder, oj == unknownStageOrder
		if iUnknown != jUnknown {
			return !iUnknown
		}
		if !iUnknown && oi != oj {
			return oi > oj
		}
		return scoreSortVal(out[i].fields["score"]) < scoreSortVal(out[j].fields["score"])
	})
	return out
}

// defaultPageLimit is the row cap the default status listing applies when the
// operator passes no explicit --limit.
const defaultPageLimit = 25

// paginationWindow is the resolved [start,end) slice (0-based, end exclusive)
// of a sorted row set for one --page/--limit request, plus the metadata the
// text footer and JSON pagination object render.
type paginationWindow struct {
	page    int
	limit   int
	total   int
	start   int
	end     int
	hasNext bool
}

// paginate computes the pagination window over total already-sorted rows.
// limit=0 disables pagination: the whole set, one page, no next. A page past
// the end (start clipped to total) is a valid empty window, not an error.
func paginate(total, page, limit int) paginationWindow {
	if limit == 0 {
		return paginationWindow{page: page, limit: 0, total: total, start: 0, end: total, hasNext: false}
	}
	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}
	return paginationWindow{page: page, limit: limit, total: total, start: start, end: end, hasNext: end < total}
}

// display returns the 1-based inclusive (start, end) row numbers the footer
// and JSON pagination object report. An empty window (e.g. a page past the
// last row) has no "row 1" to anchor to, so it reports (0, 0).
func (w paginationWindow) display() (start, end int) {
	if w.start >= w.end {
		return 0, 0
	}
	return w.start + 1, w.end
}

// paginationFooter renders the text-mode "Showing 1-25 of 83 (page 1; use
// --page 2 or --limit 0 for all)" footer, or "" when there is no next page
// (every matching row is already shown).
func paginationFooter(w paginationWindow) string {
	if !w.hasNext {
		return ""
	}
	start, end := w.display()
	return fmt.Sprintf("Showing %d-%d of %d (page %d; use --page %d or --limit 0 for all)",
		start, end, w.total, w.page, w.page+1)
}

// computeReadyGates returns only durable gate states that require scheduling,
// in the same deterministic order as the status table.
func computeReadyGates(entities []*entity, stages []Stage) []*entity {
	var out []*entity
	for _, e := range sortDefault(entities, stages) {
		switch e.fields["gate-readiness"] {
		case "awaiting-captain", "withdrawn-awaiting-prepare", "approved-awaiting-merge", "approved-awaiting-advance":
			out = append(out, e)
		}
	}
	return out
}

// sortNext sorts entities by score (desc, empty last). Stable. Matches
// sort_key_next.
func sortNext(entities []*entity) []*entity {
	out := append([]*entity(nil), entities...)
	sort.SliceStable(out, func(i, j int) bool {
		return scoreSortVal(out[i].fields["score"]) < scoreSortVal(out[j].fields["score"])
	})
	return out
}

// padRight left-justifies s in a field of width w. Width is counted in runes to
// match Python's %-Ns (which pads by code-point count). When s is at least w
// runes wide it is emitted in full.
func padRight(s string, w int) string {
	n := len([]rune(s))
	if n >= w {
		return s
	}
	return s + strings.Repeat(" ", w-n)
}

// printStatusTable writes the default table, optionally with extra columns,
// bounded to the page/limit window. Matches print_status_table plus the
// pagination slice. suppressHeader drops the header + separator rows AND the
// pagination footer for --quiet, emitting data rows only.
func printStatusTable(w io.Writer, entities []*entity, stages []Stage, extras []string, suppressHeader bool, win paginationWindow) {
	sorted := sortDefault(entities, stages)
	page := sorted[win.start:win.end]

	if len(extras) == 0 {
		row := func(a, b, c, d, e string) string {
			return padRight(a, 6) + " " + padRight(b, 30) + " " + padRight(c, 20) + " " +
				padRight(d, 30) + " " + e
		}
		if !suppressHeader {
			fmt.Fprintln(w, row("ID", "SLUG", "STATUS", "TITLE", "SCORE"))
			fmt.Fprintln(w, row("--", "----", "------", "-----", "-----"))
		}
		for _, e := range page {
			fmt.Fprintln(w, row(e.fields["id"], e.fields["slug"], e.fields["status"],
				e.fields["title"], e.fields["score"]))
		}
		printPaginationFooter(w, suppressHeader, win)
		return
	}

	base := func(a, b, c, d, e string) string {
		return padRight(a, 6) + " " + padRight(b, 30) + " " + padRight(c, 20) + " " +
			padRight(d, 30) + " " + padRight(e, 8)
	}
	headerExtras := upperAll(extras)
	sepExtras := dashSeps(headerExtras)
	if !suppressHeader {
		fmt.Fprintln(w, base("ID", "SLUG", "STATUS", "TITLE", "SCORE")+" "+joinExtras(headerExtras))
		fmt.Fprintln(w, base("--", "----", "------", "-----", "-----")+" "+joinExtras(sepExtras))
	}
	for _, e := range page {
		cells := make([]string, len(extras))
		for i, name := range extras {
			cells[i] = formatColumnCell(name, e.fields[name])
		}
		fmt.Fprintln(w, base(e.fields["id"], e.fields["slug"], e.fields["status"],
			e.fields["title"], e.fields["score"])+" "+joinExtras(cells))
	}
	printPaginationFooter(w, suppressHeader, win)
}

// printPaginationFooter writes the pagination footer line unless suppressed
// (--quiet) or there is no next page to point to.
func printPaginationFooter(w io.Writer, suppressHeader bool, win paginationWindow) {
	if suppressHeader {
		return
	}
	if footer := paginationFooter(win); footer != "" {
		fmt.Fprintln(w, footer)
	}
}

// dispatchable is an entity augmented with its computed next stage. Mirrors the
// dict the oracle builds in print_next_table.
type dispatchable struct {
	e            *entity
	next         string
	nextWorktree string
}

// computeDispatchable runs the --next dispatch rules and returns the ordered
// dispatchable list. Matches the candidate loop in print_next_table.
func computeDispatchable(entities []*entity, stages []Stage) []dispatchable {
	disp, _ := dispatchAnalysis(entities, stages)
	return disp
}

// dispatchAnalysis runs the --next candidate loop once and returns BOTH the
// ordered dispatchable list and, for every entity, why it was NOT dispatched
// (the #230 next-suppressed-by reason). The reason mirrors the loop's skip
// rules exactly so the visibility surface can never diverge from --next: a
// dispatched entity gets "" (it is not suppressed). Reasons, in the loop's
// own precedence: terminal | gate | worktree-set | concurrency-full. An entity
// whose status is not a known stage, or that sits at the last stage with no
// successor, is not attributable to one of the tracked suppression reasons and
// gets "". computeDispatchable and the surface share this one function so the
// dispatch logic is never duplicated.
func dispatchAnalysis(entities []*entity, stages []Stage) ([]dispatchable, map[*entity]string) {
	stageByName := map[string]Stage{}
	var stageNames []string
	for _, s := range stages {
		stageByName[s.Name] = s
		stageNames = append(stageNames, s.Name)
	}

	activeCounts := map[string]int{}
	for _, e := range entities {
		if e.fields["worktree"] != "" {
			activeCounts[e.fields["status"]]++
		}
	}

	candidates := sortNext(entities)
	nextCounts := map[string]int{}
	for k, v := range activeCounts {
		nextCounts[k] = v
	}

	reasons := map[*entity]string{}
	var out []dispatchable
	for _, e := range candidates {
		status := e.fields["status"]
		stg, ok := stageByName[status]
		if !ok {
			reasons[e] = ""
			continue
		}
		idx := indexOf(stageNames, status)
		if stg.terminal {
			reasons[e] = "terminal"
			continue
		}
		if stg.gate {
			reasons[e] = "gate"
			continue
		}
		if e.fields["worktree"] != "" {
			reasons[e] = "worktree-set"
			continue
		}
		nextName := status
		nextStage := stg
		if !enteredStageAwaitingCompletion(e, stg) {
			if idx+1 >= len(stageNames) {
				reasons[e] = ""
				continue
			}
			nextName = stageNames[idx+1]
			nextStage = stageByName[nextName]
		}
		if nextCounts[nextName] >= nextStage.concurrency {
			reasons[e] = "concurrency-full"
			continue
		}
		nextCounts[nextName]++
		nw := "no"
		if nextStage.Worktree {
			nw = "yes"
		}
		reasons[e] = ""
		out = append(out, dispatchable{e: e, next: nextName, nextWorktree: nw})
	}
	return out, reasons
}

// suppressedByField is the computed field name the #230 visibility surface
// exposes via --fields / --where. It is NOT a frontmatter key: it is derived
// from the --next dispatch analysis, so it is materialized only when explicitly
// named. --all-fields alone does not surface it (it is not a stored key), but
// once it is materialized because the same invocation also named it via --fields
// or --where, --all-fields includes it like any other present field.
const suppressedByField = "next-suppressed-by"

// materializeSuppressedBy writes the computed next-suppressed-by reason into
// each entity's fields ONLY when the field is explicitly referenced by --fields
// or a --where clause. When it is not referenced, e.fields is untouched, so the
// value stays out of --all-fields and the default field set and the parity-pinned
// frontmatter-keys surfaces stay byte-identical. The reason is computed from the
// shared dispatch analysis over the read's entity set, so the surface mirrors
// --next exactly. Stages may be nil (no stages block) — then every entity gets
// "" because nothing is dispatchable-or-suppressed without stage metadata.
func materializeSuppressedBy(entities []*entity, stages []Stage, explicitFields []string, filters []whereFilter) {
	referenced := false
	for _, f := range explicitFields {
		if f == suppressedByField {
			referenced = true
		}
	}
	for _, f := range filters {
		if f.field == suppressedByField {
			referenced = true
		}
	}
	if !referenced {
		return
	}
	_, reasons := dispatchAnalysis(entities, stages)
	for _, e := range entities {
		e.fields[suppressedByField] = reasons[e]
	}
}

// printNextTable writes the --next table, optionally with extras. Matches
// print_next_table. suppressHeader drops the header + separator rows for
// --quiet, emitting data rows only.
func printNextTable(w io.Writer, entities []*entity, stages []Stage, extras []string, suppressHeader bool) {
	disp := computeDispatchable(entities, stages)

	if len(extras) == 0 {
		row := func(a, b, c, d, e string) string {
			return padRight(a, 6) + " " + padRight(b, 30) + " " + padRight(c, 20) + " " +
				padRight(d, 20) + " " + e
		}
		if !suppressHeader {
			fmt.Fprintln(w, row("ID", "SLUG", "CURRENT", "NEXT", "WORKTREE"))
			fmt.Fprintln(w, row("--", "----", "-------", "----", "--------"))
		}
		for _, d := range disp {
			fmt.Fprintln(w, row(d.e.fields["id"], d.e.fields["slug"], d.e.fields["status"], d.next, d.nextWorktree))
		}
		return
	}

	base := func(a, b, c, d, e string) string {
		return padRight(a, 6) + " " + padRight(b, 30) + " " + padRight(c, 20) + " " +
			padRight(d, 20) + " " + padRight(e, 8)
	}
	headerExtras := upperAll(extras)
	sepExtras := dashSeps(headerExtras)
	if !suppressHeader {
		fmt.Fprintln(w, base("ID", "SLUG", "CURRENT", "NEXT", "WORKTREE")+" "+joinExtras(headerExtras))
		fmt.Fprintln(w, base("--", "----", "-------", "----", "--------")+" "+joinExtras(sepExtras))
	}
	for _, d := range disp {
		cells := make([]string, len(extras))
		for i, name := range extras {
			cells[i] = formatColumnCell(name, d.e.fields[name])
		}
		fmt.Fprintln(w, base(d.e.fields["id"], d.e.fields["slug"], d.e.fields["status"], d.next, d.nextWorktree)+" "+joinExtras(cells))
	}
}

// joinExtras renders the extra cells with %-20s for all but the last (plain %s).
// Matches the extra_fmt_parts layout in the oracle.
func joinExtras(cells []string) string {
	if len(cells) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, c := range cells {
		if i == len(cells)-1 {
			sb.WriteString(c)
		} else {
			sb.WriteString(padRight(c, 20))
			sb.WriteString(" ")
		}
	}
	return sb.String()
}

func upperAll(names []string) []string {
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = strings.ToUpper(n)
	}
	return out
}

// dashSeps builds separator dashes of min(len(name),20) for each header extra.
// Matches sep_extras.
func dashSeps(headers []string) []string {
	out := make([]string, len(headers))
	for i, h := range headers {
		n := len(h)
		if n > 20 {
			n = 20
		}
		out[i] = strings.Repeat("-", n)
	}
	return out
}

// localMergeSentinelPrefix is the pr-field prefix the no-PR fallback writes
// (`pr=local-merge:{short-sha}`) to honestly record that a local merge — not a
// remote PR — shipped the entity. The table renders such a value distinctly so a
// reader can tell a local merge from a real PR reference.
const localMergeSentinelPrefix = "local-merge:"

// formatPRCell renders the pr column: a local-merge sentinel
// (`local-merge:{short-sha}`) shows as `{short-sha} (local)`, distinguishing the
// no-PR fallback from a real PR reference; any other value renders unchanged.
// Truncation is applied after the transform via formatExtraCell.
func formatPRCell(value string) string {
	if sha := strings.TrimPrefix(value, localMergeSentinelPrefix); sha != value {
		value = sha + " (local)"
	}
	return formatExtraCell(value)
}

// formatColumnCell renders an extra-column cell by field name, routing the pr
// field through the local-merge-sentinel transform and every other field through
// the plain extra-cell formatter.
func formatColumnCell(name, value string) string {
	if name == "pr" {
		return formatPRCell(value)
	}
	return formatExtraCell(value)
}

// formatExtraCell renders an extra-column cell: blank for empty, truncate to 20
// with a U+2026 ellipsis. Matches format_extra_cell. Width counts characters
// (runes), matching Python len() on str.
func formatExtraCell(value string) string {
	const width = 20
	runes := []rune(value)
	if len(runes) > width {
		return string(runes[:width-1]) + "…"
	}
	return value
}

// resolveExtraFields returns the extra column names: explicit fields de-duped
// against the table's displayed columns, or every non-empty non-default
// frontmatter key (sorted) under --all-fields, else none. The de-dupe drops an
// explicit name that already names a displayed column so it is not rendered a
// second time as an extra. displayedColumns are the columns the table already
// shows (defaultStatusFields for the default table, the computed five for
// --next); defaultFields is the --all-fields exclusion set.
func resolveExtraFields(entities []*entity, explicitFields []string, allFields bool, defaultFields, displayedColumns []string) []string {
	if explicitFields != nil {
		displayed := map[string]bool{}
		for _, f := range displayedColumns {
			displayed[f] = true
		}
		var out []string
		for _, f := range explicitFields {
			if displayed[f] {
				continue
			}
			out = append(out, f)
		}
		return out
	}
	if allFields {
		defaults := map[string]bool{}
		for _, f := range defaultFields {
			defaults[f] = true
		}
		seen := map[string]bool{}
		for _, e := range entities {
			for key, val := range e.fields {
				if strings.HasPrefix(key, "_") || defaults[key] || val == "" {
					continue
				}
				seen[key] = true
			}
		}
		out := make([]string, 0, len(seen))
		for k := range seen {
			out = append(out, k)
		}
		sort.Strings(out)
		return out
	}
	return nil
}

func indexOf(s []string, v string) int {
	for i, x := range s {
		if x == v {
			return i
		}
	}
	return -1
}
