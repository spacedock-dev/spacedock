// ABOUTME: AC-2/AC-3 pagination unit tests (paginate/parsePageLimitArgs
// ABOUTME: boundaries) plus a >25-row fixture proving --page/--limit end-to-end.
package status

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

// TestPaginate is the boundary unit test for the pagination-window helper:
// a full page, a partial last page, a page past the end (a valid empty
// window, not an error), and limit=0 disabling pagination outright.
func TestPaginate(t *testing.T) {
	cases := []struct {
		name               string
		total, page, lim   int
		wantStart, wantEnd int
		wantHasNext        bool
	}{
		{"first full page of many", 100, 1, 25, 0, 25, true},
		{"middle page", 100, 2, 25, 25, 50, true},
		{"last exact page", 100, 4, 25, 75, 100, false},
		{"page past the end is a valid empty window", 100, 5, 25, 100, 100, false},
		{"exactly one page, no next", 25, 1, 25, 0, 25, false},
		{"one row over one page", 26, 1, 25, 0, 25, true},
		{"second page of the one-row overflow", 26, 2, 25, 25, 26, false},
		{"limit 0 returns everything on page 1", 100, 1, 0, 0, 100, false},
		{"empty set", 0, 1, 25, 0, 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			win := paginate(tc.total, tc.page, tc.lim)
			if win.start != tc.wantStart || win.end != tc.wantEnd || win.hasNext != tc.wantHasNext {
				t.Fatalf("paginate(%d,%d,%d) = {start:%d end:%d hasNext:%v}, want {start:%d end:%d hasNext:%v}",
					tc.total, tc.page, tc.lim, win.start, win.end, win.hasNext, tc.wantStart, tc.wantEnd, tc.wantHasNext)
			}
		})
	}
}

// TestPaginationFooter locks the footer's exact rendering and its suppression
// once there is no next page.
func TestPaginationFooter(t *testing.T) {
	got := paginationFooter(paginate(83, 1, 25))
	want := "Showing 1-25 of 83 (page 1; use --page 2 or --limit 0 for all)"
	if got != want {
		t.Fatalf("paginationFooter = %q, want %q", got, want)
	}
	if got := paginationFooter(paginate(25, 1, 25)); got != "" {
		t.Fatalf("paginationFooter on the last page = %q, want empty (no next page)", got)
	}
	if got := paginationFooter(paginate(100, 5, 25)); got != "" {
		t.Fatalf("paginationFooter on a page past the end = %q, want empty", got)
	}
}

// TestParsePageLimitArgs is the unit test for the --page/--limit argv parser:
// defaults, explicit values, invalid values (non-integer, zero page, negative),
// missing arguments, and the --page-with-limit-0 contradiction.
func TestParsePageLimitArgs(t *testing.T) {
	cases := []struct {
		name                    string
		args                    []string
		wantPage, wantLimit     int
		wantPageSet, wantLimSet bool
		wantErr                 bool
	}{
		{"no flags default to page 1 limit 25", nil, 1, 25, false, false, false},
		{"explicit page", []string{"--page", "3"}, 3, 25, true, false, false},
		{"explicit limit", []string{"--limit", "10"}, 1, 10, false, true, false},
		{"limit 0 alone disables pagination", []string{"--limit", "0"}, 1, 0, false, true, false},
		{"page and limit together", []string{"--page", "2", "--limit", "10"}, 2, 10, true, true, false},
		{"page zero is invalid", []string{"--page", "0"}, 0, 0, false, false, true},
		{"page negative is invalid", []string{"--page", "-1"}, 0, 0, false, false, true},
		{"page non-integer is invalid", []string{"--page", "abc"}, 0, 0, false, false, true},
		{"limit negative is invalid", []string{"--limit", "-1"}, 0, 0, false, false, true},
		{"limit non-integer is invalid", []string{"--limit", "abc"}, 0, 0, false, false, true},
		{"page missing argument", []string{"--page"}, 0, 0, false, false, true},
		{"limit missing argument", []string{"--limit"}, 0, 0, false, false, true},
		{"page with limit 0 is contradictory", []string{"--page", "2", "--limit", "0"}, 0, 0, false, false, true},
		{"page 1 with limit 0 is still contradictory", []string{"--page", "1", "--limit", "0"}, 0, 0, false, false, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			page, limit, pageSet, limSet, err := parsePageLimitArgs(tc.args)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("parsePageLimitArgs(%v) = nil error, want an error", tc.args)
				}
				return
			}
			if err != nil {
				t.Fatalf("parsePageLimitArgs(%v) unexpected error: %v", tc.args, err)
			}
			if page != tc.wantPage || limit != tc.wantLimit || pageSet != tc.wantPageSet || limSet != tc.wantLimSet {
				t.Fatalf("parsePageLimitArgs(%v) = (%d,%d,%v,%v), want (%d,%d,%v,%v)",
					tc.args, page, limit, pageSet, limSet, tc.wantPage, tc.wantLimit, tc.wantPageSet, tc.wantLimSet)
			}
		})
	}
}

// buildPaginationFixture writes a 30-row single-workflow fixture, all in the
// same status so only score orders them: row-01 (score 0.99) down to row-30
// (score 0.70), monotonically, so a page's slug slice is unambiguous and the
// row order and the pagination window can be asserted independently of the
// stage-then-score comparator under test elsewhere.
func buildPaginationFixture(t *testing.T, n int) string {
	t.Helper()
	def := t.TempDir()
	writeFile(t, filepath.Join(def, "README.md"), "---\ncommissioned-by: spacedock@1\nid-style: slug\n---\n# Pagination Fixture\n")
	for i := 1; i <= n; i++ {
		slug := fmt.Sprintf("row-%02d", i)
		score := 0.99 - float64(i-1)*0.01
		body := fmt.Sprintf("---\nstatus: ideation\nscore: \"%.2f\"\n---\n", score)
		writeFile(t, filepath.Join(def, slug+".md"), body)
	}
	return def
}

// splitTableAndFooter splits a rendered status table into its data-row slugs
// and the trailing pagination footer line (empty when absent). Like
// tableSlugs, but treats the "Showing X-Y of Z ..." footer as the footer
// rather than a malformed data row.
func splitTableAndFooter(t *testing.T, table string) (slugs []string, footer string) {
	t.Helper()
	lines := strings.Split(strings.TrimRight(table, "\n"), "\n")
	if len(lines) < 2 {
		t.Fatalf("table has no header/separator rows:\n%s", table)
	}
	for _, line := range lines[2:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, "Showing ") {
			footer = line
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			t.Fatalf("malformed data row %q in table:\n%s", line, table)
		}
		slugs = append(slugs, fields[1])
	}
	return slugs, footer
}

func rowSlugs(prefix string, from, to int) []string {
	var out []string
	for i := from; i <= to; i++ {
		out = append(out, fmt.Sprintf("%s-%02d", prefix, i))
	}
	return out
}

// paginatedStatusEnvelope mirrors the default listing's {"command","entities",
// "pagination"} JSON shape for the pagination-object assertions below.
type paginatedStatusEnvelope struct {
	Command    string              `json:"command"`
	Entities   []map[string]string `json:"entities"`
	Pagination map[string]string   `json:"pagination"`
}

// TestStatusPaginationDefaultBounds proves the default (no --page/--limit)
// listing bounds a 30-row fixture to the first 25, in both text (with the
// footer) and JSON (with the pagination object), naming the exact row set
// that would fail if pagination were skipped or applied before sorting.
func TestStatusPaginationDefaultBounds(t *testing.T) {
	def := buildPaginationFixture(t, 30)
	env := pinnedEnv(t)

	out, stderr, code := runNative(t, def, env, "--workflow-dir", def)
	if code != 0 {
		t.Fatalf("default status exit=%d stderr=%q", code, stderr)
	}
	slugs, footer := splitTableAndFooter(t, out)
	want := rowSlugs("row", 1, 25)
	if !equalStrings(slugs, want) {
		t.Fatalf("default page rows = %v, want first 25 (row-01..row-25)\n%s", slugs, out)
	}
	wantFooter := "Showing 1-25 of 30 (page 1; use --page 2 or --limit 0 for all)"
	if footer != wantFooter {
		t.Fatalf("footer = %q, want %q", footer, wantFooter)
	}

	jsonOut, jsonErr, jsonCode := runNative(t, def, env, "--workflow-dir", def, "--json")
	if jsonCode != 0 {
		t.Fatalf("--json exit=%d stderr=%q", jsonCode, jsonErr)
	}
	var env1 paginatedStatusEnvelope
	if err := json.Unmarshal([]byte(jsonOut), &env1); err != nil {
		t.Fatalf("parse --json: %v\n%s", err, jsonOut)
	}
	if len(env1.Entities) != 25 {
		t.Fatalf("--json entities count = %d, want 25", len(env1.Entities))
	}
	if env1.Entities[0]["slug"] != "row-01" || env1.Entities[24]["slug"] != "row-25" {
		t.Fatalf("--json entity slugs at page boundaries = %q..%q, want row-01..row-25",
			env1.Entities[0]["slug"], env1.Entities[24]["slug"])
	}
	wantPagination := map[string]string{
		"page": "1", "limit": "25", "total": "30", "start": "1", "end": "25", "has_next": "true",
	}
	for k, v := range wantPagination {
		if env1.Pagination[k] != v {
			t.Fatalf("--json pagination[%q] = %q, want %q\n%s", k, env1.Pagination[k], v, jsonOut)
		}
	}
}

// TestStatusPaginationPage2 proves --page 2 selects the remaining 5 rows and
// stops paginating (no footer, has_next=false).
func TestStatusPaginationPage2(t *testing.T) {
	def := buildPaginationFixture(t, 30)
	env := pinnedEnv(t)

	out, stderr, code := runNative(t, def, env, "--workflow-dir", def, "--page", "2")
	if code != 0 {
		t.Fatalf("--page 2 exit=%d stderr=%q", code, stderr)
	}
	slugs, footer := splitTableAndFooter(t, out)
	want := rowSlugs("row", 26, 30)
	if !equalStrings(slugs, want) {
		t.Fatalf("--page 2 rows = %v, want row-26..row-30\n%s", slugs, out)
	}
	if footer != "" {
		t.Fatalf("--page 2 (the last page) footer = %q, want empty", footer)
	}
}

// TestStatusPaginationPageOutOfRange proves a page past the last row is a
// valid, empty page (not an error) with accurate total metadata.
func TestStatusPaginationPageOutOfRange(t *testing.T) {
	def := buildPaginationFixture(t, 30)
	env := pinnedEnv(t)

	jsonOut, jsonErr, code := runNative(t, def, env, "--workflow-dir", def, "--page", "5", "--json")
	if code != 0 {
		t.Fatalf("--page 5 --json exit=%d stderr=%q", code, jsonErr)
	}
	var env1 paginatedStatusEnvelope
	if err := json.Unmarshal([]byte(jsonOut), &env1); err != nil {
		t.Fatalf("parse --json: %v\n%s", err, jsonOut)
	}
	if len(env1.Entities) != 0 {
		t.Fatalf("--page 5 entities = %v, want empty", env1.Entities)
	}
	if env1.Pagination["total"] != "30" || env1.Pagination["has_next"] != "false" || env1.Pagination["start"] != "0" || env1.Pagination["end"] != "0" {
		t.Fatalf("--page 5 pagination = %v, want total=30 has_next=false start=0 end=0", env1.Pagination)
	}
}

// TestStatusPaginationLimitZero proves --limit 0 restores every row with the
// new stage-then-score ordering intact and no footer/JSON truncation.
func TestStatusPaginationLimitZero(t *testing.T) {
	def := buildPaginationFixture(t, 30)
	env := pinnedEnv(t)

	out, stderr, code := runNative(t, def, env, "--workflow-dir", def, "--limit", "0")
	if code != 0 {
		t.Fatalf("--limit 0 exit=%d stderr=%q", code, stderr)
	}
	slugs, footer := splitTableAndFooter(t, out)
	want := rowSlugs("row", 1, 30)
	if !equalStrings(slugs, want) {
		t.Fatalf("--limit 0 rows = %v, want all 30 rows in row-01..row-30 order\n%s", slugs, out)
	}
	if footer != "" {
		t.Fatalf("--limit 0 footer = %q, want empty (no truncation)", footer)
	}

	jsonOut, jsonErr, jsonCode := runNative(t, def, env, "--workflow-dir", def, "--limit", "0", "--json")
	if jsonCode != 0 {
		t.Fatalf("--limit 0 --json exit=%d stderr=%q", jsonCode, jsonErr)
	}
	var env1 paginatedStatusEnvelope
	if err := json.Unmarshal([]byte(jsonOut), &env1); err != nil {
		t.Fatalf("parse --json: %v\n%s", err, jsonOut)
	}
	if len(env1.Entities) != 30 {
		t.Fatalf("--limit 0 --json entities count = %d, want 30", len(env1.Entities))
	}
	if env1.Pagination["limit"] != "0" || env1.Pagination["has_next"] != "false" || env1.Pagination["total"] != "30" {
		t.Fatalf("--limit 0 --json pagination = %v, want limit=0 has_next=false total=30", env1.Pagination)
	}
}
