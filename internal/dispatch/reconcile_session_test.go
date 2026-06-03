// ABOUTME: AC-1..AC-5 session-scoping regression — bare reconcile follows
// ABOUTME: leadSessionId, never poisons A/B/C from a foreign team, degrades to D/E.
package dispatch

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
)

// writeTeamConfigFile writes a config.json under the fixture home with the given
// team name, lead session id, and roster, then stamps its mtime to modAge ago so
// a test can make a decoy NEWER than the session-matched config (proving the
// loader no longer consults mtime).
func (f *reconcileFixture) writeTeamConfigFile(team, leadSessionID string, members []claudeteam.ReconcileMember, modAge time.Duration) string {
	f.t.Helper()
	cfgPath := filepath.Join(f.home, ".claude", "teams", team, "config.json")
	writeFile(f.t, cfgPath, teamConfigJSONWithSession(team, leadSessionID, members))
	mod := time.Now().Add(-modAge)
	if err := os.Chtimes(cfgPath, mod, mod); err != nil {
		f.t.Fatal(err)
	}
	return cfgPath
}

// runBare drives Reconcile with an empty team name and the given session id,
// returning the parsed envelope plus the raw stderr (so degrade tests can assert
// the note and exit code). It does NOT t.Fatal on a non-zero code — the caller
// asserts the code.
func (f *reconcileFixture) runBare(sessionID string) (reconcileResult, string, int) {
	f.t.Helper()
	var stdout, stderr bytes.Buffer
	opts := reconcileOpts{
		workflowDir: f.workflowDir,
		teamName:    "", // bare: no --team-name
		sessionID:   sessionID,
		repoRoot:    f.repoRoot,
		include:     map[string]bool{"A": true, "B": true, "C": true, "D": true, "E": true},
		home:        f.home,
		roster:      claudeteam.LoadReconcileTeam,
		gh: func(pr string) (string, error) {
			if state, ok := f.ghResponses[pr]; ok {
				return state, nil
			}
			return "OPEN", nil
		},
		git: gitRunnerExec,
	}
	code := Reconcile(opts, &stdout, &stderr)
	var result reconcileResult
	if code == 0 {
		if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
			f.t.Fatalf("decode: %v\nstdout=%s", err, stdout.String())
		}
	}
	return result, stderr.String(), code
}

// the fixture's own config carries leadSessionId "session-fixture" and a
// Class-A-tripping archived ensign roster. A current session id that differs
// from it is, by construction, the foreign/stale-team case.
const currentSessionID = "aec06dd4-current-session-86011"

// TestReconcileForeignTeamNeverPoisonsRoster is AC-1: the only on-disk ensign
// team belongs to a foreign session (leadSessionId != the injected current
// session), and a NEWER-mtime decoy is also foreign — so the retired
// newest-mtime loader would resolve a team and emit Class A. The fixed loader
// must degrade: zero A/B/C entries. This locks both the repeated-session (stale
// prior team) and parallel-session (another live session's team) cases.
func TestReconcileForeignTeamNeverPoisonsRoster(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	f := newReconcileFixture(t)
	// Decoy with a NEWER mtime than the fixture config, still foreign. If the
	// loader consulted mtime it would pick this; session-id must ignore it.
	f.writeTeamConfigFile("team-newer-decoy", "decoy-session-zzz",
		[]claudeteam.ReconcileMember{
			{Name: "team-lead", AgentType: "team-lead", Model: "m"},
			{Name: "spacedock-ensign-release-notes-local-summary-implementation",
				AgentType: "spacedock:ensign", Model: "m"},
		}, 0) // mtime = now (newest)

	result, _, code := f.runBare(currentSessionID)
	if code != 0 {
		t.Fatalf("bare reconcile exit=%d, want 0 (degrade is not an error)", code)
	}
	for _, d := range result.Drift {
		if d.Class == "A" || d.Class == "B" || d.Class == "C" {
			t.Errorf("foreign team poisoned roster class %s: %+v", d.Class, d)
		}
	}
	if result.TeamName != "" {
		t.Errorf("degrade must report empty team_name; got %q", result.TeamName)
	}
}

// TestReconcileSessionMatchedDiscovery is AC-2: a config whose leadSessionId
// EQUALS the injected current session id (plus a newer-mtime foreign decoy)
// resolves via session identity — the archived-entity ensign yields a Class A
// entry. The decoy is newer-mtime to make the distinction load-bearing: mtime
// would pick the decoy, session-id picks the match.
func TestReconcileSessionMatchedDiscovery(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	f := newReconcileFixture(t)
	// Match: leadSessionId == current session, with the Class-A archived ensign.
	f.writeTeamConfigFile("team-session-match", currentSessionID,
		[]claudeteam.ReconcileMember{
			{Name: "team-lead", AgentType: "team-lead", Model: "m"},
			{Name: "spacedock-ensign-release-notes-local-summary-implementation",
				AgentType: "spacedock:ensign", Model: "m"},
		}, 10*time.Minute) // older mtime than the decoy below
	// Decoy: foreign session, NEWER mtime. mtime would pick this; session-id must not.
	f.writeTeamConfigFile("team-newer-decoy", "decoy-session-zzz",
		[]claudeteam.ReconcileMember{
			{Name: "team-lead", AgentType: "team-lead", Model: "m"},
			{Name: "spacedock-ensign-alive-implementation",
				AgentType: "spacedock:ensign", Model: "m"},
		}, 0)

	result, _, code := f.runBare(currentSessionID)
	if code != 0 {
		t.Fatalf("bare reconcile exit=%d, want 0", code)
	}
	if result.TeamName != "team-session-match" {
		t.Errorf("resolved team_name=%q, want team-session-match (session id, not mtime)", result.TeamName)
	}
	aFor := 0
	for _, d := range result.Drift {
		if d.Class == "A" && d.Slug == "release-notes-local-summary" {
			aFor++
		}
	}
	if aFor != 1 {
		t.Errorf("session-matched roster: want exactly 1 Class A for the archived ensign; got %d\n%s",
			aFor, formatDrift(result.Drift))
	}
}

// TestReconcileExplicitTeamNameIgnoresSession is AC-3: passing --team-name
// directly computes A/B/C against that roster regardless of session id — the
// explicit path never consults leadSessionId. Reuses the full-sweep fixture run
// (which passes f.teamName) to confirm no regression; here we assert the
// explicit path produces roster classes even when the injected session id does
// not match the config's leadSessionId.
func TestReconcileExplicitTeamNameIgnoresSession(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	f := newReconcileFixture(t)
	var stdout, stderr bytes.Buffer
	opts := reconcileOpts{
		workflowDir: f.workflowDir,
		teamName:    f.teamName, // explicit
		sessionID:   "totally-unrelated-session",
		repoRoot:    f.repoRoot,
		include:     map[string]bool{"A": true, "B": true, "C": true, "D": true, "E": true},
		home:        f.home,
		roster:      claudeteam.LoadReconcileTeam,
		gh: func(pr string) (string, error) {
			if state, ok := f.ghResponses[pr]; ok {
				return state, nil
			}
			return "OPEN", nil
		},
		git: gitRunnerExec,
	}
	if code := Reconcile(opts, &stdout, &stderr); code != 0 {
		t.Fatalf("explicit reconcile exit=%d stderr=%s", code, stderr.String())
	}
	var result reconcileResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if result.TeamName != f.teamName {
		t.Errorf("explicit team_name=%q, want %q", result.TeamName, f.teamName)
	}
	byClass := groupDriftByClass(result.Drift)
	for _, class := range []string{"A", "B", "C"} {
		if len(byClass[class]) == 0 {
			t.Errorf("explicit path with non-matching session id must still emit class %s; got none", class)
		}
	}
}

// TestReconcileDegradeEmitsGitClasses is AC-4: bare reconcile with no
// session-matched team still emits the git/filesystem classes D and E (the
// fixture's stale worktree and ahead-main). Confirms the degrade path ran the
// sweep — it suppressed only the roster classes.
func TestReconcileDegradeEmitsGitClasses(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	f := newReconcileFixture(t)
	result, _, code := f.runBare(currentSessionID)
	if code != 0 {
		t.Fatalf("bare reconcile exit=%d, want 0", code)
	}
	byClass := groupDriftByClass(result.Drift)
	if len(byClass["D"]) == 0 {
		t.Errorf("degrade must still emit Class D (stale branch); got none\n%s", formatDrift(result.Drift))
	}
	if len(byClass["E"]) == 0 {
		t.Errorf("degrade must still emit Class E (stale local main); got none\n%s", formatDrift(result.Drift))
	}
	for _, class := range []string{"A", "B", "C"} {
		if len(byClass[class]) != 0 {
			t.Errorf("degrade must suppress class %s; got %s", class, formatDrift(byClass[class]))
		}
	}
}

// TestReconcileDegradeExitAndNote is AC-5: the git-only degrade is exit-0 with a
// one-line stderr note that roster reconciliation needs a team identity, and an
// explicit --team-name pointing at a missing config still exits 1 (the degrade
// must not mask a genuine setup failure).
func TestReconcileDegradeExitAndNote(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	f := newReconcileFixture(t)

	// Degrade: exit 0 + the roster-suppressed note on stderr.
	_, stderr, code := f.runBare(currentSessionID)
	if code != 0 {
		t.Errorf("degrade exit=%d, want 0", code)
	}
	if !strings.Contains(stderr, "roster reconciliation needs a team identity") {
		t.Errorf("degrade stderr missing the roster-suppressed note; got %q", stderr)
	}

	// Setup failure: explicit --team-name at a missing config still exits 1.
	var stdout, stderr2 bytes.Buffer
	opts := reconcileOpts{
		workflowDir: f.workflowDir,
		teamName:    "team-does-not-exist",
		sessionID:   currentSessionID,
		repoRoot:    f.repoRoot,
		include:     map[string]bool{"A": true, "B": true, "C": true, "D": true, "E": true},
		home:        f.home,
		roster:      claudeteam.LoadReconcileTeam,
		gh:          func(string) (string, error) { return "OPEN", nil },
		git:         gitRunnerExec,
	}
	if code := Reconcile(opts, &stdout, &stderr2); code != 1 {
		t.Errorf("missing explicit team config exit=%d, want 1 (degrade must not mask setup failure); stderr=%s",
			code, stderr2.String())
	}
}

// TestReconcileGateSuppressesEvenWithPopulatedRoster pins the assembly-level
// rosterTrusted gate (reconcile.go) directly via a stub rosterLoader, decoupled
// from the loader's own degrade logic. The stub returns the degrade sentinel
// (empty TeamName) but deliberately carries a Class-A-tripping ensign in
// Members — so if the A/B/C emit blocks were not guarded by `rosterTrusted`,
// this populated roster would leak a Class A entry. The gate must suppress it on
// the strength of the empty TeamName alone, regardless of what Members holds.
func TestReconcileGateSuppressesEvenWithPopulatedRoster(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	f := newReconcileFixture(t)
	var stdout, stderr bytes.Buffer
	opts := reconcileOpts{
		workflowDir: f.workflowDir,
		teamName:    "",
		sessionID:   currentSessionID,
		repoRoot:    f.repoRoot,
		include:     map[string]bool{"A": true, "B": true, "C": true, "D": true, "E": true},
		// Degrade sentinel (empty TeamName) that nonetheless carries a roster.
		roster: func(home, teamName, sessionID string) (claudeteam.ReconcileTeamState, error) {
			return claudeteam.ReconcileTeamState{
				TeamName: "", // sentinel → rosterTrusted == false
				Members: []claudeteam.ReconcileMember{
					{Name: "spacedock-ensign-release-notes-local-summary-implementation",
						AgentType: "spacedock:ensign", Model: "m"},
				},
			}, nil
		},
		home: f.home,
		gh:   func(string) (string, error) { return "OPEN", nil },
		git:  gitRunnerExec,
	}
	if code := Reconcile(opts, &stdout, &stderr); code != 0 {
		t.Fatalf("degrade exit=%d, want 0; stderr=%s", code, stderr.String())
	}
	var result reconcileResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	for _, d := range result.Drift {
		if d.Class == "A" || d.Class == "B" || d.Class == "C" {
			t.Errorf("rosterTrusted gate leaked class %s despite empty TeamName sentinel: %+v", d.Class, d)
		}
	}
	// The git/filesystem classes must still emit — the gate suppresses roster
	// classes only, not the whole sweep.
	byClass := groupDriftByClass(result.Drift)
	if len(byClass["D"]) == 0 && len(byClass["E"]) == 0 {
		t.Errorf("gate must still emit git classes D/E on degrade; got none\n%s", formatDrift(result.Drift))
	}
}
