// ABOUTME: Survey sync→codex-presence end-to-end test — drives a real `agentsview sync`
// ABOUTME: over synthesized Codex sources, then runs the codex-presence query against it.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The fixture-only query-smoke (survey_queries_test.go) injects Codex rows straight into
// the DB, so it proves the codex-presence QUERY counts them but never exercises the SYNC
// that must put them there. This test closes that gap: it runs the real `agentsview sync`
// the SKILL.md §1 recipe runs, then runs the codex-presence query from queries.sql against
// the synced DB. It reds if the sync path yields no Codex rows for the repo project — the
// failure mode behind a survey that renders "0 Codex sessions" against a Codex-blind DB.

// writeCodexSession writes a minimal agentsview-ingestible Codex rollout jsonl under
// codexDir for a session whose source cwd is repoCwd. agentsview derives `project` from
// the git-root basename of that cwd and blanks the stored cwd, so the session lands with
// project = basename(repoCwd-normalized) and cwd = ''.
func writeCodexSession(t *testing.T, codexDir, id, repoCwd string) {
	t.Helper()
	day := filepath.Join(codexDir, "2026", "06", "02")
	if err := os.MkdirAll(day, 0o755); err != nil {
		t.Fatalf("mkdir codex day dir: %v", err)
	}
	// A git repo at repoCwd so agentsview resolves a git-root basename for `project`.
	if err := os.MkdirAll(repoCwd, 0o755); err != nil {
		t.Fatalf("mkdir repo cwd: %v", err)
	}
	if out, err := exec.Command("git", "-C", repoCwd, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init %s: %v\n%s", repoCwd, err, out)
	}
	jsonl := `{"timestamp":"2026-06-02T10:00:00.000Z","type":"session_meta","payload":{"id":"` + id +
		`","timestamp":"2026-06-02T10:00:00.000Z","cwd":"` + repoCwd +
		`","originator":"codex-tui","cli_version":"0.136.0","source":"cli"}}` + "\n" +
		`{"timestamp":"2026-06-02T10:00:02.000Z","type":"response_item","payload":{"type":"message","role":"user","content":[{"type":"input_text","text":"hi"}]}}` + "\n"
	if err := os.WriteFile(filepath.Join(day, "rollout-2026-06-02T10-00-00-"+id+".jsonl"), []byte(jsonl), 0o644); err != nil {
		t.Fatalf("write codex session: %v", err)
	}
}

// TestSurveyCodexPresenceThroughSync drives the §1 sync recipe over a HOME-isolated,
// synthesized Codex source (one session in the repo-project, one in a same-basename
// sibling root), then runs the codex-presence query from queries.sql against the synced
// DB. The sync is the real `agentsview` binary — the artifact under test end-to-end. It
// asserts codex-presence counts BOTH same-basename sessions (the documented collision)
// with blank_cwd > 0. A sync that fails to ingest Codex reds the count assertion.
func TestSurveyCodexPresenceThroughSync(t *testing.T) {
	agentsview, err := exec.LookPath("agentsview")
	if err != nil {
		t.Skip("agentsview not on PATH; the sync→codex-presence e2e needs the real binary")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not on PATH; needed to give the synthesized Codex cwd a git-root basename")
	}
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 not on PATH; needed to run the codex-presence query")
	}

	root := t.TempDir()
	home := filepath.Join(root, "home")           // empty HOME isolates ALL default sources
	dataDir := filepath.Join(root, "data")         // AGENTSVIEW_DATA_DIR
	codexDir := filepath.Join(root, "codex")       // CODEX_SESSIONS_DIR
	claudeDir := filepath.Join(root, "claude")     // CLAUDE_PROJECTS_DIR (empty — Codex-only test)
	for _, d := range []string{home, dataDir, codexDir, claudeDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}

	// Two Codex sessions whose source cwds are DIFFERENT roots sharing the basename `proj`.
	// agentsview keys both to project='proj', blanks both cwds — so codex-presence by
	// project name counts both (the same-basename collision the skill warns about).
	writeCodexSession(t, codexDir, "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", filepath.Join(root, "repoA", "proj"))
	writeCodexSession(t, codexDir, "bbbbbbbb-cccc-dddd-eeee-ffffffffffff", filepath.Join(root, "repoB", "proj"))

	// Drive the §1 sync recipe: HOME isolated, only CODEX/CLAUDE sources overridden.
	sync := exec.Command(agentsview, "sync")
	sync.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + home,
		"AGENTSVIEW_DATA_DIR=" + dataDir,
		"CODEX_SESSIONS_DIR=" + codexDir,
		"CLAUDE_PROJECTS_DIR=" + claudeDir,
	}
	if out, err := sync.CombinedOutput(); err != nil {
		t.Fatalf("agentsview sync: %v\n%s", err, out)
	}

	db := filepath.Join(dataDir, "sessions.db")
	if _, err := os.Stat(db); err != nil {
		t.Fatalf("sync produced no sessions.db at %s: %v", db, err)
	}

	queries := loadLabeledQueries(t)
	codexPresence, ok := queries["codex-presence"]
	if !ok {
		t.Fatalf("queries.sql is missing the codex-presence query (have: %v)", sortedQueryNames(queries))
	}

	// Run codex-presence bound to :repo_project='proj' against the SYNCED db.
	script := ".param set :repo_project 'proj'\n" + codexPresence + "\n"
	cmd := exec.Command("sqlite3", db)
	cmd.Stdin = strings.NewReader(script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run codex-presence against synced db: %v\n%s", err, out)
	}
	row := strings.TrimSpace(string(out))
	fields := strings.Split(row, "|")
	if len(fields) != 2 {
		t.Fatalf("codex-presence row should be codex_sessions|blank_cwd, got %q", row)
	}
	if fields[0] != "2" {
		t.Errorf("codex-presence should count both same-basename Codex sessions ingested by the sync, got %q (a Codex-blind sync would yield 0)", fields[0])
	}
	if fields[1] == "0" {
		t.Errorf("agentsview blanks Codex cwd, so blank_cwd must be > 0, got %q", fields[1])
	}
}
