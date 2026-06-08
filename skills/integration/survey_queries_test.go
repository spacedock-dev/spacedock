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

// execSQLite runs a non-query SQL statement (an UPDATE) against the fixture DB. The
// non-vacuousness sub-tests mutate a fresh fixture copy and re-run a query to prove an
// expected value FLIPS under the mutation — so the query is load-bearing, not a constant.
func execSQLite(t *testing.T, db, stmt string) {
	t.Helper()
	sqlite3, err := exec.LookPath("sqlite3")
	if err != nil {
		t.Skip("sqlite3 not on PATH")
	}
	cmd := exec.Command(sqlite3, db)
	cmd.Stdin = strings.NewReader(stmt + "\n")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("exec mutation against fixture: %v\n%s", err, out)
	}
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

	for _, name := range []string{
		"scoping", "codex-presence", "codex-scoped", "codex-workstreams", "codex-activity",
		"scaffold-usage", "work-by-area", "decision-open", "mode-classification",
	} {
		if _, ok := queries[name]; !ok {
			t.Fatalf("recommended-SQL reference is missing the %q query (have: %v)", name, sortedQueryNames(queries))
		}
	}

	// scoping (#318): under the corrected git-root-basename model every in-repo checkout
	// shares ONE `project` key; the cwd-prefix-union does the load-bearing work — it counts
	// the in-repo Claude sessions (cwd AT root, subdir, worktree, plus the F/G worktree-shape
	// + the mode-classification track sessions, all under the prefix) and EXCLUDES the
	// blank-cwd session, the out-of-repo session, and ALL the codex rows. The fixture has 9
	// in-repo Claude sessions: A,B,C + WT + issue-feed×2 + landing-copy×2 + mixed-bag.
	t.Run("scoping", func(t *testing.T) {
		rows := runQuery(t, db, queries["scoping"])
		if len(rows) != 1 {
			t.Fatalf("scoping should return one summary row, got %d: %v", len(rows), rows)
		}
		fields := strings.Split(rows[0], "|")
		if len(fields) != 3 {
			t.Fatalf("scoping row should have 3 fields (sessions|blank_cwd|span) — folded_keys is dropped, got: %q", rows[0])
		}
		if fields[0] != "9" {
			t.Errorf("the cwd-prefix should count 9 in-repo Claude sessions, got sessions=%q", fields[0])
		}
		if fields[1] != "0" {
			t.Errorf("the blank-cwd Claude session is outside the prefix and must not count, got blank_cwd=%q", fields[1])
		}
	})

	// codex-presence (#69): Codex sessions land cwd='' (agentsview does not persist Codex
	// cwd), so the cwd-prefix scope misses them. This separate flagged count matches by
	// `project = :repo_project` ALONE — which means it also catches a same-basename SIBLING
	// repo's Codex sessions (the documented collision). The fixture has five such rows (four
	// in-repo F* + one same-basename sibling G), all blank-cwd, so the count is 5 and
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
		if fields[0] != "5" {
			t.Errorf("codex-presence should count 5 Codex sessions matching the repo project name (4 in-repo F* + same-basename sibling G), got %q", fields[0])
		}
		if fields[1] == "0" {
			t.Errorf("Codex cwd is unrecorded so blank_cwd must be > 0, got blank_cwd=%q", fields[1])
		}
	})

	// codex-scoped (#321, AC-1): attributes Codex to THIS repo by exec_command.$.workdir
	// prefix — DISTINCT from codex-presence's name-only match. The four F* sessions have an
	// exec_command whose $.workdir is under /repo/proj (one is a worktree path), so they are
	// IN scope; the sibling G's workdir is under /sibling/proj, so it is EXCLUDED. The count
	// is 4 (the four F*), strictly fewer than codex-presence's 5 — proving the two signals
	// MEASURE DIFFERENT THINGS (scoped ⊂ presence, sibling-free). (AC-1 illustrates the
	// mechanism at "1 vs 2"; the clustering AC needs 4 attributed sessions, so the fixture
	// scales to 4 vs 5 — the binding asserts, sibling-exclusion + the prefix-load-bearing
	// flip, hold identically.) Non-vacuous: re-pointing G's workdir under /repo/proj flips
	// the count 4→5, proving the prefix is load-bearing, not a constant.
	t.Run("codex-scoped", func(t *testing.T) {
		rows := runQuery(t, db, queries["codex-scoped"])
		if len(rows) != 1 {
			t.Fatalf("codex-scoped should return one count row, got %d: %v", len(rows), rows)
		}
		if rows[0] != "4" {
			t.Errorf("codex-scoped should count 4 workdir-attributed Codex sessions (F* in-repo, sibling G excluded), got %q", rows[0])
		}
		// distinct from codex-presence (5) — the two signals differ on the same fixture.
		pres := runQuery(t, db, queries["codex-presence"])
		if presCount := strings.Split(pres[0], "|")[0]; presCount == rows[0] {
			t.Errorf("codex-scoped (%q) must differ from codex-presence (%q) — scoped is the sibling-free subset", rows[0], presCount)
		}
		// non-vacuous: re-point sibling G's exec_command workdir UNDER the repo prefix → 4 becomes 5.
		db2 := buildFixtureDB(t)
		execSQLite(t, db2, `UPDATE tool_calls SET input_json='{"command":"go build","workdir":"/repo/proj"}' WHERE id=46;`)
		flipped := runQuery(t, db2, queries["codex-scoped"])
		if flipped[0] != "5" {
			t.Errorf("re-pointing the sibling's workdir under the repo prefix must flip codex-scoped 4→5 (prefix is load-bearing), got %q", flipped[0])
		}
	})

	// codex-workstreams (#322, AC-3): clusters the codex-scoped sessions by the 3-case rule —
	// dispatch-pattern → {TASK} (stage stripped), task/entity backtick → {TASK}, else
	// (unlabeled). The expected labels are SUBSTRINGS of the fixture first_messages (an
	// independent source — never written in SKILL.md), so a broken extractor reds. Non-vacuous:
	// the stage suffix must be STRIPPED (journey-cost-ledger, NOT journey-cost-ledger-implementation),
	// the two distinct dispatch tasks must NOT merge (codex-live-ci separate), the backtick
	// task name must anchor past the leading reviewer-label backtick (orient-workflow-discovery,
	// not 142-validation/Ensign), and (unlabeled) must sort LAST.
	t.Run("codex-workstreams", func(t *testing.T) {
		rows := runQuery(t, db, queries["codex-workstreams"])
		got := map[string]string{}
		for _, r := range rows {
			f := strings.Split(r, "|")
			if len(f) != 2 {
				t.Fatalf("codex-workstreams row should be workstream|sessions, got: %q", r)
			}
			got[f[0]] = f[1]
		}
		for _, want := range []string{"journey-cost-ledger", "orient-workflow-discovery", "codex-live-ci", "(unlabeled)"} {
			if got[want] != "1" {
				t.Errorf("workstream %q should cluster 1 session, got %q in %v", want, got[want], got)
			}
		}
		if _, leaked := got["journey-cost-ledger-implementation"]; leaked {
			t.Errorf("the dispatch stage suffix must be stripped — saw an un-stripped label in %v", got)
		}
		if _, leaked := got["142-validation/Ensign"]; leaked {
			t.Errorf("the task/entity label must anchor past the leading reviewer-label backtick, got %v", got)
		}
		if len(got) != 4 {
			t.Errorf("expected exactly 4 workstream buckets (3 named + unlabeled), got %v", got)
		}
		// (unlabeled) sorts last so the named tracks lead the rendered list.
		if last := strings.Split(rows[len(rows)-1], "|")[0]; last != "(unlabeled)" {
			t.Errorf("(unlabeled) must sort last, got trailing row %q", rows[len(rows)-1])
		}
	})

	// codex-activity (#323): per-tool tally over the codex-scoped set — exec_command (4, one
	// per F* session), update_plan (1), spawn_agent (1). The sibling G's exec_command must NOT
	// count (it is outside the workdir prefix), proving the activity tally honors the same scope.
	t.Run("codex-activity", func(t *testing.T) {
		rows := runQuery(t, db, queries["codex-activity"])
		got := map[string]string{}
		for _, r := range rows {
			f := strings.Split(r, "|")
			if len(f) != 2 {
				t.Fatalf("codex-activity row should be tool|calls, got: %q", r)
			}
			got[f[0]] = f[1]
		}
		if got["exec_command"] != "4" {
			t.Errorf("exec_command should tally 4 over the codex-scoped set (sibling G excluded), got %q in %v", got["exec_command"], got)
		}
		if got["update_plan"] != "1" {
			t.Errorf("update_plan should tally 1, got %q in %v", got["update_plan"], got)
		}
		if got["spawn_agent"] != "1" {
			t.Errorf("spawn_agent should tally 1, got %q in %v", got["spawn_agent"], got)
		}
	})

	// no-union (AC-2c): the added Codex rows must NOT inflate the Claude scope. The scoping
	// query is asserted to 9 above (the Claude-only in-repo count), proving Codex stays out
	// of the Claude `sessions` count — a flagged presence, never a silent project union.
	t.Run("codex-not-folded-into-scope", func(t *testing.T) {
		rows := runQuery(t, db, queries["scoping"])
		if len(rows) != 1 {
			t.Fatalf("scoping should return one summary row, got %d: %v", len(rows), rows)
		}
		if sessions := strings.Split(rows[0], "|")[0]; sessions != "9" {
			t.Errorf("the Codex rows must not be folded into the Claude scope; scoping.sessions should stay 9, got %q", sessions)
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

	// work-by-area (#317.2, F-corrected / AC-7a): Edit/Write file_paths bucket by LOGICAL
	// area after stripping any `.worktrees/<wt>/` (or `.claude/worktrees/<wt>/`) physical
	// prefix — so a worktree `src/` edit and a main-checkout `src/` edit BOTH bucket as `src`
	// (NOT `.worktrees`/`<external>`). A `kind` partition demotes genuine config
	// (`.claude`/`.beads`/`.git`/`<external>`) WITHOUT filtering it (still counted), and the
	// ORDER puts product areas FIRST. The fixture's `src` bucket has 4 edits: 2 worktree
	// (render.ts, palette.ts) + main.ts + feed.ts — the worktree strip is what folds them.
	t.Run("work-by-area", func(t *testing.T) {
		rows := runQuery(t, db, queries["work-by-area"])
		kind := map[string]string{}
		edits := map[string]string{}
		var order []string
		for _, r := range rows {
			f := strings.Split(r, "|")
			if len(f) != 3 {
				t.Fatalf("work-by-area row should be area|kind|edits, got: %q", r)
			}
			kind[f[0]] = f[1]
			edits[f[0]] = f[2]
			order = append(order, f[0])
		}
		// worktree src/ edits attribute to `src` ALONGSIDE the main-checkout src/ edit.
		if edits["src"] != "4" {
			t.Errorf("the 2 worktree src/ edits + 2 main-checkout src/ edits should all bucket as src=4 (the strip folds them), got %q in %v", edits["src"], edits)
		}
		// a worktree src/ edit must NEVER leak into a `.worktrees` bucket (the strip is load-bearing).
		if _, leaked := edits[".worktrees"]; leaked {
			t.Errorf("a worktree edit must NOT bucket as `.worktrees` — the physical prefix must be stripped; got %v", edits)
		}
		// `.claude/worktrees/<wt>/internal/codex.go` strips to `internal` (the second worktree layout).
		if edits["internal"] != "4" {
			t.Errorf("internal should count 4 (build.go, parse.go, index.go, the .claude/worktrees-stripped codex.go), got %q in %v", edits["internal"], edits)
		}
		// genuine config demotes to kind=config (still counted), NOT filtered.
		for _, c := range []string{".claude", ".beads", "<external>"} {
			if kind[c] != "config" {
				t.Errorf("%s should be tagged kind=config (demoted, still counted), got %q in %v", c, kind[c], kind)
			}
		}
		if kind["src"] != "product" || kind["docs"] != "product" || kind["internal"] != "product" {
			t.Errorf("product areas (src/docs/internal) should be tagged kind=product, got %v", kind)
		}
		// product leads: the first row must be a product area, never a config one.
		if len(order) > 0 && kind[order[0]] != "product" {
			t.Errorf("a product area must lead the work-by-area ordering, got leading %q (kind=%q)", order[0], kind[order[0]])
		}
		// non-vacuous: re-point a worktree src/ edit to `.claude/` → it leaves `src` for the config footnote.
		db2 := buildFixtureDB(t)
		execSQLite(t, db2, `UPDATE tool_calls SET input_json='{"file_path":"/repo/proj/.claude/render.ts"}' WHERE id=50;`)
		rerows := runQuery(t, db2, queries["work-by-area"])
		reEdits := map[string]string{}
		for _, r := range rerows {
			f := strings.Split(r, "|")
			reEdits[f[0]] = f[2]
		}
		if reEdits["src"] != "3" {
			t.Errorf("re-pointing one worktree src/ edit to .claude/ must drop src 4→3, got %q in %v", reEdits["src"], reEdits)
		}
		if reEdits[".claude"] != "2" {
			t.Errorf("the re-pointed edit must move to the .claude config bucket (1→2), got %q in %v", reEdits[".claude"], reEdits)
		}
	})

	// mode-classification (#324, G / AC-8a): classify each TRACK (keyed by git_branch) into a
	// work MODE from the per-track signal tallies (veto density, gate-pass ratio, loop markers,
	// edit-kind). The fixture carries a MECHANICAL track (issue-feed: gate-pass, worktree loop,
	// code edits, no veto), an EXPLORATION track (landing-copy: vetoes, a rejected path, .md
	// edits), and a NEITHER-DOMINANT track (mixed-bag → unlabeled). The labels DERIVE from the
	// signal rows (the independent oracle), never from SKILL.md text. Non-vacuous: (i) swapping
	// the mechanical track's rows to carry high vetoes + a rejected path + prose flips its label
	// to exploration; (ii) the neither-dominant track stays unlabeled (no guessed automation).
	t.Run("mode-classification", func(t *testing.T) {
		rows := runQuery(t, db, queries["mode-classification"])
		mode := map[string]string{}
		for _, r := range rows {
			f := strings.Split(r, "|")
			if len(f) != 2 {
				t.Fatalf("mode-classification row should be track|mode, got: %q", r)
			}
			mode[f[0]] = f[1]
		}
		if mode["issue-feed"] != "mechanical" {
			t.Errorf("the gate-pass/worktree-loop/code track should classify mechanical, got %q in %v", mode["issue-feed"], mode)
		}
		if mode["landing-copy"] != "exploration" {
			t.Errorf("the high-veto/rejected/prose track should classify exploration, got %q in %v", mode["landing-copy"], mode)
		}
		if mode["mixed-bag"] != "unlabeled" {
			t.Errorf("a neither-dominant track must stay unlabeled (generic book-keeping, never a guessed automation pitch), got %q in %v", mode["mixed-bag"], mode)
		}
		// non-vacuous (i): swap issue-feed's signals (high veto + rejected path + prose) → flips to exploration.
		db2 := buildFixtureDB(t)
		execSQLite(t, db2, `UPDATE messages SET content='[Request interrupted by user]' WHERE session_id='claude:91111111-1111-1111-1111-111111111111';`)
		execSQLite(t, db2, `UPDATE messages SET content='doesn''t want to proceed' WHERE session_id='claude:92222222-2222-2222-2222-222222222222' AND id=8;`)
		execSQLite(t, db2, `UPDATE tool_calls SET result_content='The user doesn''t want to proceed with this tool use.' WHERE id=60;`)
		execSQLite(t, db2, `UPDATE tool_calls SET input_json='{"file_path":"/repo/proj/content/a.md"}' WHERE id=61;`)
		execSQLite(t, db2, `UPDATE tool_calls SET input_json='{"file_path":"/repo/proj/content/b.md"}' WHERE id=62;`)
		flipped := map[string]string{}
		for _, r := range runQuery(t, db2, queries["mode-classification"]) {
			f := strings.Split(r, "|")
			flipped[f[0]] = f[1]
		}
		if flipped["issue-feed"] != "exploration" {
			t.Errorf("swapping the mechanical track's signals to the exploration signature must flip its label, got %q in %v", flipped["issue-feed"], flipped)
		}
		if flipped["mixed-bag"] != "unlabeled" {
			t.Errorf("the neither-dominant track must stay unlabeled under the signal swap, got %q in %v", flipped["mixed-bag"], flipped)
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
