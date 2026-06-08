// ABOUTME: Survey query-smoke — runs each labeled query from skills/survey/references/
// ABOUTME: queries.sql against a committed agentsview-shaped fixture and asserts the corrected shape.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// fixtureRepoRoot is the repo-root prefix the smoke binds for the prefix-scoped queries.
// The fixture seeds sessions whose cwd is AT this root, in a subdir under it, in a
// worktree-style path under it, blank, and OUTSIDE it — so the scoping/bucketing queries
// have something real to coalesce and exclude. No live git repo is needed for the smoke;
// the live drive (AC-4/AC-5) exercises the real `git rev-parse` resolution.
const fixtureRepoRoot = "/repo/proj"

// fixtureRepoProject is the git-root-basename `project` key the smoke binds for the
// codex-presence query. agentsview derives `project` from the repo's git-root basename
// (non-alphanumerics → `_`), so every in-repo checkout — root, subdir, worktree, state
// dir — shares this ONE key. The fixture's in-repo Claude sessions all carry it, and the
// blank-cwd Codex rows (cwd unrecorded for Codex) are matched by it alone.
const fixtureRepoProject = "proj"

// buildFixtureDB shells out to the sqlite3 CLI to materialize the committed
// fixture-sessions.sql into a temp sessions.db, returning its path. The skill's
// queries run via sqlite3, so sqlite3 is the faithful executor; it is a standard POSIX
// tool present in CI, and the test SKIPS (not fails) when it is absent so the suite stays
// runnable on a minimal box without claiming a false pass.
func buildFixtureDB(t *testing.T) string {
	t.Helper()
	sqlite3, err := exec.LookPath("sqlite3")
	if err != nil {
		t.Skip("sqlite3 not on PATH; survey query-smoke needs it to run the recommended queries")
	}
	sqlPath := filepath.Join("testdata", "survey", "fixture-sessions.sql")
	sql, err := os.ReadFile(sqlPath)
	if err != nil {
		t.Fatalf("read fixture SQL %s: %v", sqlPath, err)
	}
	db := filepath.Join(t.TempDir(), "sessions.db")
	cmd := exec.Command(sqlite3, db)
	cmd.Stdin = strings.NewReader(string(sql))
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build fixture DB: %v\n%s", err, out)
	}
	return db
}

// queryNameRe matches a `-- name: <id>` label line in the recommended-SQL reference.
var queryNameRe = regexp.MustCompile(`^--\s*name:\s*(\S+)\s*$`)

// loadLabeledQueries reads skills/survey/references/queries.sql and returns the SQL body
// for each `-- name: <id>` block (the lines after the label up to the next label / EOF).
// The smoke EXECUTES this reference file — it runs the artifact, never parses SKILL.md
// prose — so a missing/renamed reference file (the AC-1 RED lever) fails here, and a
// broken or dropped query reds the per-query assertion below.
func loadLabeledQueries(t *testing.T) map[string]string {
	t.Helper()
	path := filepath.Join(repoRoot(t), "skills", "survey", "references", "queries.sql")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read recommended-SQL reference %s: %v", path, err)
	}
	queries := map[string]string{}
	var name string
	var body []string
	flush := func() {
		if name != "" {
			queries[name] = strings.TrimSpace(strings.Join(body, "\n"))
		}
	}
	for _, line := range strings.Split(string(data), "\n") {
		if m := queryNameRe.FindStringSubmatch(line); m != nil {
			flush()
			name = m[1]
			body = nil
			continue
		}
		if name == "" {
			continue // preamble before the first labeled query
		}
		if strings.HasPrefix(strings.TrimSpace(line), "--") {
			continue // annotation comment inside a block — not SQL
		}
		body = append(body, line)
	}
	flush()
	return queries
}

// runQuery runs one labeled query against the fixture DB with :repo_root bound to
// fixtureRepoRoot and :repo_project bound to fixtureRepoProject, returning the rows as
// pipe-separated lines (sqlite3's list mode). Binding a param a query does not reference
// is harmless, so both are always set — the cwd-prefix queries read :repo_root, the
// codex-presence query reads :repo_project.
func runQuery(t *testing.T, db, query string) []string {
	t.Helper()
	sqlite3, err := exec.LookPath("sqlite3")
	if err != nil {
		t.Skip("sqlite3 not on PATH")
	}
	script := ".param set :repo_root '" + fixtureRepoRoot + "'\n" +
		".param set :repo_project '" + fixtureRepoProject + "'\n" + query + "\n"
	cmd := exec.Command(sqlite3, db)
	cmd.Stdin = strings.NewReader(script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run query against fixture: %v\n%s", err, out)
	}
	var rows []string
	for _, line := range strings.Split(strings.TrimRight(string(out), "\n"), "\n") {
		if line != "" {
			rows = append(rows, line)
		}
	}
	return rows
}

// TestSurveyQuerySmoke is the AC-2 query-smoke. It runs each labeled query from
// skills/survey/references/queries.sql against a committed production-shaped fixture DB
// and asserts the CORRECTED shape. Expected values come from the FIXTURE rows — an
// independent source that diverges from the skill prose — so a broken query, a dropped
// self-exclusion, a wrong prefix bound, or schema drift reds a per-query check. This is
// the deliberately light bar: it pins SQL-extraction correctness, not the full
// orchestration (the prose join / cross-check / report assembly are AC-3's live drive).
func TestSurveyQuerySmoke(t *testing.T) {
	db := buildFixtureDB(t)
	queries := loadLabeledQueries(t)

	for _, name := range []string{"scoping", "codex-presence", "scaffold-usage", "work-by-area", "decision-open"} {
		if _, ok := queries[name]; !ok {
			t.Fatalf("recommended-SQL reference is missing the %q query (have: %v)", name, sortedQueryNames(queries))
		}
	}

	// scoping (#318): under the corrected git-root-basename model every in-repo checkout
	// shares ONE `project` key, so COUNT(DISTINCT project) is structurally always 1 and
	// `folded_keys` is gone. The cwd-prefix-union still does the load-bearing work — it
	// counts the cwd-AT-root + subdir + worktree sessions (3) and EXCLUDES the same-basename
	// sibling, the blank-cwd session, the out-of-repo session, and the codex rows. The row
	// is the corrected 3-field shape: sessions|blank_cwd|span.
	t.Run("scoping", func(t *testing.T) {
		rows := runQuery(t, db, queries["scoping"])
		if len(rows) != 1 {
			t.Fatalf("scoping should return one summary row, got %d: %v", len(rows), rows)
		}
		fields := strings.Split(rows[0], "|")
		if len(fields) != 3 {
			t.Fatalf("scoping row should have 3 fields (sessions|blank_cwd|span) — folded_keys is dropped, got: %q", rows[0])
		}
		if fields[0] != "3" {
			t.Errorf("the cwd-prefix should count 3 in-repo Claude sessions, got sessions=%q", fields[0])
		}
		if fields[1] != "0" {
			t.Errorf("the blank-cwd Claude session is outside the prefix and must not count, got blank_cwd=%q", fields[1])
		}
	})

	// codex-presence (#69): Codex sessions land cwd='' (agentsview does not persist Codex
	// cwd), so the cwd-prefix scope misses them. This separate flagged count matches by
	// `project = :repo_project` ALONE — which means it also catches a same-basename SIBLING
	// repo's Codex sessions (the documented collision). The fixture has two such rows (one
	// in-repo, one same-basename sibling shape), both blank-cwd, so the count is 2 and
	// blank_cwd > 0. This is a presence flag, NOT a union — the scoping count below is
	// asserted UNCHANGED by these rows.
	t.Run("codex-presence", func(t *testing.T) {
		rows := runQuery(t, db, queries["codex-presence"])
		if len(rows) != 1 {
			t.Fatalf("codex-presence should return one summary row, got %d: %v", len(rows), rows)
		}
		fields := strings.Split(rows[0], "|")
		if len(fields) != 2 {
			t.Fatalf("codex-presence row should have 2 fields (codex_sessions|blank_cwd), got: %q", rows[0])
		}
		if fields[0] != "2" {
			t.Errorf("codex-presence should count 2 Codex sessions matching the repo project name (in-repo + same-basename sibling), got %q", fields[0])
		}
		if fields[1] == "0" {
			t.Errorf("Codex cwd is unrecorded so blank_cwd must be > 0, got blank_cwd=%q", fields[1])
		}
	})

	// no-union (AC-2c): the added Codex rows must NOT inflate the Claude scope. The scoping
	// query is asserted to 3 above (the same value the pre-Codex fixture yielded), proving
	// Codex stays out of the Claude `sessions` count — a flagged presence, never a silent
	// project union.
	t.Run("codex-not-folded-into-scope", func(t *testing.T) {
		rows := runQuery(t, db, queries["scoping"])
		if len(rows) != 1 {
			t.Fatalf("scoping should return one summary row, got %d: %v", len(rows), rows)
		}
		if sessions := strings.Split(rows[0], "|")[0]; sessions != "3" {
			t.Errorf("the Codex rows must not be folded into the Claude scope; scoping.sessions should stay 3, got %q", sessions)
		}
	})

	// scaffold-usage (#319): the behavioral tally GROUPs skill_name by FAMILY over the
	// repo-scoped set, EXCLUDES the dominant spacedock self rows, and folds the
	// namespaced + bare superpowers rows into ONE `superpowers` family. The out-of-repo
	// session's superpowers row must NOT inflate the count.
	t.Run("scaffold-usage", func(t *testing.T) {
		rows := runQuery(t, db, queries["scaffold-usage"])
		got := map[string]string{}
		for _, r := range rows {
			f := strings.Split(r, "|")
			if len(f) != 2 {
				t.Fatalf("scaffold-usage row should be family|invocations, got: %q", r)
			}
			got[f[0]] = f[1]
		}
		if _, present := got["spacedock"]; present {
			t.Errorf("the spacedock self-invocation family MUST be excluded from the scaffold tally, got: %v", got)
		}
		if got["superpowers"] != "3" {
			t.Errorf("superpowers should tally 3 in-repo invocations (2 namespaced + 1 bare), got %q in %v", got["superpowers"], got)
		}
		if len(got) != 1 {
			t.Errorf("only the superpowers family should remain after excluding spacedock, got %v", got)
		}
	})

	// work-by-area (#317.2): Edit/Write file_paths bucket by first package segment under
	// the repo root; a path OUTSIDE the prefix buckets as <external> (a reference, not
	// this project's identity).
	t.Run("work-by-area", func(t *testing.T) {
		rows := runQuery(t, db, queries["work-by-area"])
		got := map[string]string{}
		for _, r := range rows {
			f := strings.Split(r, "|")
			if len(f) != 2 {
				t.Fatalf("work-by-area row should be area|edits, got: %q", r)
			}
			got[f[0]] = f[1]
		}
		if got["internal"] != "2" {
			t.Errorf("two edits under internal/ should bucket as internal=2, got %q in %v", got["internal"], got)
		}
		if got["skills"] != "1" {
			t.Errorf("one write under skills/ should bucket as skills=1, got %q in %v", got["skills"], got)
		}
		if got["<external>"] != "1" {
			t.Errorf("the edit to a sibling repo outside the prefix should bucket as <external>=1, got %q in %v", got["<external>"], got)
		}
	})

	// decision-open (#320): the rejected AskUserQuestion is OPEN; the ExitPlanMode
	// "User has approved your plan" row is `done` via the new done-prefix (the cheap fix —
	// before it, an approved plan fell to the OPEN frontier); OPEN sorts first.
	t.Run("decision-open", func(t *testing.T) {
		rows := runQuery(t, db, queries["decision-open"])
		if len(rows) == 0 {
			t.Fatal("decision-open returned no rows")
		}
		status := map[string]string{} // header -> status
		for _, r := range rows {
			f := strings.Split(r, "|")
			if len(f) < 2 {
				t.Fatalf("decision-open row should be header|status|question, got: %q", r)
			}
			status[f[0]] = f[1]
		}
		if status["Test framework"] != "OPEN" {
			t.Errorf("the rejected 'Test framework' decision should be OPEN, got %q in %v", status["Test framework"], status)
		}
		if status["Refactor scope"] != "done" {
			t.Errorf("the answered 'Refactor scope' decision should be done, got %q in %v", status["Refactor scope"], status)
		}
		if status["PLAN"] != "done" {
			t.Errorf("the ExitPlanMode 'User has approved your plan' row should be done (the #320 done-prefix fix), got %q in %v", status["PLAN"], status)
		}
		// OPEN must lead so the recency LIMIT cannot truncate the frontier.
		if first := strings.Split(rows[0], "|"); len(first) >= 2 && first[1] != "OPEN" {
			t.Errorf("the OPEN frontier must sort first so the LIMIT cannot hide it, got leading row: %q", rows[0])
		}
	})
}

// sortedQueryNames returns the labeled query names for a diagnostic message.
func sortedQueryNames(q map[string]string) []string {
	names := make([]string, 0, len(q))
	for n := range q {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
