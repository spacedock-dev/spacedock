// ABOUTME: AC-2 Claude extraction proof — runs the survey scan artifact against
// ABOUTME: a committed agentsview-shaped fixture DB and asserts the Claude signals surface.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// buildFixtureDB shells out to the sqlite3 CLI to materialize the committed
// fixture-sessions.sql into a temp sessions.db, returning its path. The skill's
// scan artifact uses sqlite3, so sqlite3 is the faithful executor; it is
// a standard POSIX tool present in CI, and the test skips (not fails) when it or bash
// is absent so the suite stays runnable on a minimal box without claiming a false pass.
func buildFixtureDB(t *testing.T) string {
	t.Helper()
	sqlite3, err := exec.LookPath("sqlite3")
	if err != nil {
		t.Skip("sqlite3 not on PATH; survey extraction proof needs it to run the skill's inline queries")
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

// runSurveyScan runs the survey scan artifact against the given DB, from a working dir
// named for the project key (the script derives PROJECT from the cwd basename), and
// returns the combined output. DB is normally set by the sync step; the test injects the
// fixture DB directly.
func runSurveyScan(t *testing.T, db, projectKey string) string {
	t.Helper()
	script := filepath.Join(repoRoot(t), "skills", "survey", "bin", "scan-project")
	projDir := filepath.Join(t.TempDir(), projectKey)
	if err := os.Mkdir(projDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}
	cmd := exec.Command(script)
	cmd.Dir = projDir
	cmd.Env = append(os.Environ(), "DB="+db)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run survey scan artifact %s: %v\n%s", script, err, out)
	}
	return string(out)
}

// outputSection returns the lines of the survey step-2 run output under the `## X`
// marker line (exact-match on the marker), up to but excluding the next `## ` line.
// The skill echoes section markers like `## OVERVIEW`, so this scopes an assertion to a
// single section's rows — a stray token in another section cannot satisfy the check.
func outputSection(out, marker string) string {
	lines := strings.Split(out, "\n")
	start := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == marker {
			start = i + 1
			break
		}
	}
	if start == -1 {
		return ""
	}
	end := len(lines)
	for i := start; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "## ") {
			end = i
			break
		}
	}
	return strings.Join(lines[start:end], "\n")
}

// TestSurveyExtractionSurfacesClaudeSignals is the AC-2 extraction proof. It runs the
// survey scan artifact against a committed agentsview-shaped fixture DB and asserts the
// produced output surfaces the project's Claude decisions
// (the OPEN frontier and a representative answered row), the interruption counts, and
// EXCLUDES a sibling out-of-scope Codex session under the same project key.
//
// This is behavior-fixture coverage, not a SKILL.md string-match: the expected values
// (the AskUserQuestion decisions, the OPEN-vs-done status, the veto count, the codex
// step that must NOT surface) come from the FIXTURE rows — an independent source that
// diverges from the skill text. The skill's bug was a project filter that returned "no
// history" on the real key; if a project-filter or agent-scope regression returned, a
// known row would vanish (or the codex row would leak) and this test would RED. The
// proof is the EXECUTION of the survey scan artifact against known rows, never a
// substring over the instruction file.
func TestSurveyExtractionSurfacesClaudeSignals(t *testing.T) {
	db := buildFixtureDB(t)
	got := runSurveyScan(t, db, "survey_fixture_proj")
	t.Logf("survey scan output:\n%s", got)

	// 1. OVERVIEW counts only the two Claude sessions — the Codex sibling under the same
	//    project key must NOT inflate the count. A dropped agent='claude' scope would.
	overview := outputSection(got, "## OVERVIEW")
	if !strings.Contains(overview, "2 sessions") {
		t.Errorf("OVERVIEW should count exactly the 2 Claude sessions, got: %q", overview)
	}

	// 2. The OPEN frontier and an answered Claude decision surface in DECISIONS, with
	//    OPEN first. The fixture has 20 answered decisions newer than the OPEN row, so
	//    dropping ORDER BY status ASC lets the LIMIT truncate the OPEN frontier.
	decisions := outputSection(got, "## DECISIONS  (header :: status :: question;  OPEN = still needs the human)")
	if decisions == "" {
		t.Fatalf("no DECISIONS section in output:\n%s", got)
	}
	for _, header := range []string{"Test framework", "Recent answered 20"} {
		if !strings.Contains(decisions, header) {
			t.Errorf("DECISIONS missing the AskUserQuestion header %q:\n%s", header, decisions)
		}
	}
	decisionLines := strings.Split(strings.TrimSpace(decisions), "\n")
	if len(decisionLines) == 0 || !regexp.MustCompile(`^Test framework\s+::\s+OPEN`).MatchString(decisionLines[0]) {
		t.Errorf("the OPEN frontier should lead DECISIONS so the recency LIMIT cannot hide it:\n%s", decisions)
	}
	// The session-2 rejected decision is OPEN; the newer answered decision is done.
	if !regexp.MustCompile(`Test framework\s+::\s+OPEN`).MatchString(decisions) {
		t.Errorf("the unanswered 'Test framework' decision should be OPEN:\n%s", decisions)
	}
	if !regexp.MustCompile(`Recent answered 20\s+::\s+done`).MatchString(decisions) {
		t.Errorf("the answered 'Recent answered 20' decision should be done:\n%s", decisions)
	}

	// 3. Interruption math: asks=22 (the OPEN decision + 21 answered AskUserQuestion
	//    calls), vetoes=1 (one interrupt marker).
	interruptions := outputSection(got, "## INTERRUPTIONS  (how often you had to step in)")
	if !strings.Contains(interruptions, "asks=22") {
		t.Errorf("INTERRUPTIONS should count all fixture AskUserQuestion calls (asks=22):\n%s", interruptions)
	}
	if !strings.Contains(interruptions, "vetoes=1") {
		t.Errorf("INTERRUPTIONS should count the one veto marker (vetoes=1):\n%s", interruptions)
	}

	// 4. The out-of-scope Codex session's step must NOT leak into any section — Claude
	//    scope excludes it. This is the regression guard for an over-broad query.
	if strings.Contains(got, "A codex-only step that must not surface") {
		t.Errorf("the out-of-scope Codex session leaked into the Claude-scoped survey output:\n%s", got)
	}
}
