//go:build live

// ABOUTME: The pty/tmux liveDriver: drives a REAL interactive `spacedock claude` session (no -p)
// ABOUTME: so native team tools are present and the session stays resident while teammates work.
package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

// ptyLiveDriver is the second liveDriver transport (alongside the headless `-p`
// claudeLiveRunner): it launches a REAL interactive `spacedock claude` session in a
// detached tmux pane so the native TeamCreate/TeamDelete tools are present and the
// agent loop stays alive while teammates work — neither of which headless `claude
// -p` can do (the team-tool regression + the SDK end_turn race the entity records).
//
// It encodes the three spike hazards as code (AC-5):
//   - NO stdout pipe on launch: claude owns the pane tty; the driver observes via
//     the on-disk session jsonl (and capture-pane for diagnostics), never a pipe (a
//     pipe makes stdout a non-tty and the interactive TUI never renders).
//   - resolve the FO pane by #{pane_title}, not the active pane: an Agent dispatch
//     materializes the ensign as a sibling pane that becomes active, so the FO's
//     TeamDelete/marker live in the title-resolved FO pane (firstOfficerPaneIndex).
//   - gate the first send on idle: keystrokes sent mid-boot QUEUE behind the boot
//     turn and the drive stalls; waitTranscriptIdle holds the send until the pinned
//     FO transcript reads a committed turn-end (transcriptReachedIdle) for N
//     consecutive polls. The signal is the on-disk JSONL, NOT pane text — pane-text
//     scraping is render-fragile on a headless CI runner (missing terminfo → blank
//     capture-pane → it would never read idle).
//
// liveResult.stream is sourced from the session jsonl claude writes under the
// pinned CLAUDE_CONFIG_DIR — the SAME stream-json dialect the existing assertions
// and the teardown grader parse — so the shared grading runs UNCHANGED over the
// interactive transcript. finalMessage is the captured FO pane text (diagnostic;
// the team-mode tests grade the stream and the on-disk roster, not finalMessage).
type ptyLiveDriver struct {
	binary       string
	pluginDir    string
	env          []string
	modelName    string
	artifactRoot string
	homeDir      string
}

// ptyLiveDriver satisfies the liveDriver seam: the shared per-scenario
// orchestration can drive it exactly as it drives the `-p` runner, oblivious to
// the tmux transport.
var _ liveDriver = ptyLiveDriver{}

func (d ptyLiveDriver) model() string { return d.modelName }
func (d ptyLiveDriver) home() string  { return d.homeDir }

// withStubPATH returns a driver copy whose launched FO subprocess resolves a stub
// binary in dir first. It never mutates the receiver's env so parallel runs sharing
// the driver stay race-free — the mirror of claudeLiveRunner.withStubPATH.
func (d ptyLiveDriver) withStubPATH(dir string) liveDriver {
	d.env = withPATHPrefix(d.env, dir)
	return d
}

func (d ptyLiveDriver) withInvocationLedger(ledger testInvocationLedger) liveDriver {
	d.env = ledger.instrumentEnv(d.env)
	return d
}

// newPtyLiveDriver builds the pty driver against the SAME isolated auth/HOME tree
// the headless suite uses (isolatedClaudeEnv: OAuth benchmark-token / API key, else
// SKIP) so AC-6's skip-not-fatal-without-auth holds with no parallel auth path.
func newPtyLiveDriver(t *testing.T) ptyLiveDriver {
	t.Helper()
	binary := spacedockBinary(t)
	pluginDir := livePluginDir(t)
	model := envOr("SPACEDOCK_LIVE_MODEL", "sonnet")

	env := isolatedClaudeEnv(t, os.Getenv("HOME"))
	env = withBinaryOnPath(env, binary)

	homeDir, _ := envValue(env, "HOME")
	return ptyLiveDriver{
		binary:       binary,
		pluginDir:    pluginDir,
		env:          env,
		modelName:    model,
		artifactRoot: claudeLiveArtifactDir(t, "pty-team-mode"),
		homeDir:      homeDir,
	}
}

const (
	// ptyBootBudget bounds how long the driver waits for the FO to reach a first
	// committed idle turn (boot complete, ready for the first keystroke). A team-mode
	// boot + FO contract load is multi-minute, so this is roomier than a per-step
	// budget.
	ptyBootBudget = 4 * time.Minute
	// ptyIdleStablePolls is how many CONSECUTIVE idle reads of the FO transcript gate
	// the first send — the queued-keystroke hazard (a send mid-boot queues behind the
	// boot turn). N>1 debounces the residual race where the transcript flips to
	// end_turn a beat before the TUI input box accepts keys.
	ptyIdleStablePolls = 3
	// ptyPollInterval paces the idle polls.
	ptyPollInterval = 2 * time.Second
)

// ptySession is a launched, prompted interactive session: the tmux name, the FO's
// session jsonl on disk, the proc poller over the tmux session, and the artifact
// dir. The shared-scenario run() and the team-mode tests both build one via
// launchAndSend, then either drain it whole (run) or watch it incrementally (the
// team-mode tests' expect/grade beats — the faithful resurrection of the retired
// forced-team tests).
type ptySession struct {
	tmuxName    string
	projectsDir string // CLAUDE_CONFIG_DIR/projects/<encoded-cwd>: where the FO + teammate transcripts land
	foSessionID string // the FO's OWN session UUID — the tail pins to this, never a teammate's transcript
	configDir   string // the effective CLAUDE_CONFIG_DIR this run wrote under (teams/ + projects/ live here)
	proc        *tmuxLiveProc
	artifactDir string
	started     time.Time
}

// newFileSource opens a fresh tail over the FO's OWN session jsonl from byte 0, so
// a caller's streamWatcher sees the whole FO transcript from launch. It pins to the
// FO's session UUID: once the FO creates a team and dispatches a teammate, that
// teammate writes its OWN transcript into the same projects dir (newer, also
// assistant-bearing), so a newest-file tail would FLIP to the teammate and miss the
// FO's TeamCreate/Agent/TeamDelete (the F30 hazard gym's SessionFileByID solves).
// Pinning by the FO's id keeps the tail on the FO across teammate spawns.
func (s ptySession) newFileSource() *fileLineSource {
	return newFileLineSourceByID(s.projectsDir, s.foSessionID)
}

// launchAndSend launches the interactive session WITHOUT a stdout pipe (claude
// owns the pane tty — the no-pipe hazard), gates the first send on stable-idle (the
// queued-keystroke hazard), sends prompt to the title-resolved FO pane (the
// FO-pane-capture hazard), and resolves the FO's session jsonl on disk. It returns
// a ptySession ready to drain or watch. label names the run for artifacts/errors.
// The caller MUST defer session.proc.kill() to reap the resident tmux session.
func (d ptyLiveDriver) launchAndSend(t *testing.T, label, workflowRoot, prompt string) ptySession {
	t.Helper()
	artifactDir := filepath.Join(d.artifactRoot, label)
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Per-run CLAUDE_CONFIG_DIR so concurrent sessions never share claude's
	// session/config state — the same isolation the `-p` runner applies. A COPY of
	// the env (never a mutation of the shared d.env) keeps parallel runs race-free.
	env := d.env
	configDir, _ := envValue(d.env, "CLAUDE_CONFIG_DIR")
	if configDir != "" {
		configDir = filepath.Join(configDir, label)
		env = withClaudeConfigDir(d.env, configDir)
	}
	effectiveConfigDir := configDirOrDefault(configDir, d.homeDir)

	// Claude Code encodes the symlink-resolved real cwd (e.g. /tmp -> /private/tmp
	// on macOS) both in the trust-dialog project key and in the projects-dir name,
	// so resolve it ONCE and reuse it for the config seed and the jsonl path.
	resolvedCwd := workflowRoot
	if r, err := filepath.EvalSymlinks(workflowRoot); err == nil {
		resolvedCwd = r
	}

	// Seed the config-dir .claude.json so the INTERACTIVE launch skips every
	// first-run dialog (theme/login pickers + the per-project trust-folder dialog).
	// Without this seed a fresh isolated HOME stalls at a blocking picker BEFORE the
	// FO loads, so boot never reaches the committed end_turn the idle gate waits for.
	// Pre-clearing onboarding also means boot ends at a clean end_turn (no picker
	// open), which is why the transcript stop_reason signal is sufficient here without
	// a separate picker-footer check. The headless `-p` suite never hits this (non-
	// interactive `-p` consumes the OAuth token directly); the interactive transport
	// must pre-clear onboarding itself. This is the implementation hazard the spike
	// missed (it ran on already-onboarded operator auth).
	if err := seedInteractiveClaudeConfig(effectiveConfigDir, resolvedCwd); err != nil {
		t.Fatalf("seed interactive claude config for %s: %v", label, err)
	}

	// Seed the stored-login credential so the INTERACTIVE launch authenticates.
	// Interactive Claude Code does NOT consume CLAUDE_CODE_OAUTH_TOKEN the way
	// headless `-p` does (the env token 401s every interactive turn, banner
	// "Claude API"); it reads the stored-login OAuth credential from the config
	// dir's .credentials.json. isolatedClaudeEnv supplies only the env token, so
	// without this seed the isolated interactive child has no credential it honors.
	// The stored login lives in the operator's OS credential store (macOS keychain
	// `Claude Code-credentials`), so this is best-effort + operator-local: when no
	// stored login is reachable (CI Linux, no keychain) it is a no-op and the child
	// falls back to the env-token / ANTHROPIC_API_KEY path unchanged — so AC-6's
	// skip-not-fatal-without-auth (decided upstream in isolatedClaudeEnv) is intact.
	//
	// When the seed lands, DROP CLAUDE_CODE_OAUTH_TOKEN from the child env: with both
	// present, interactive claude prefers the env token (banner "Claude API") and
	// 401s, so the seeded stored login must be the ONLY credential for it to be
	// authoritative (banner "Claude Max", a real authenticated turn). The token is
	// dropped only on the seeded path; the no-keychain fallback keeps it.
	if err := seedStoredLoginCredential(effectiveConfigDir); err != nil {
		t.Logf("[pty %s] stored-login seed skipped (%v) — falling back to the env credential", label, err)
	} else {
		env = withoutEnvKey(env, "CLAUDE_CODE_OAUTH_TOKEN")
	}

	session := fmt.Sprintf("sdpty-%s-%d", label, time.Now().UnixNano())
	proc := newTmuxLiveProc(session)

	// The launch line is the live runner's exact front door minus `-p`:
	// spacedock-owned flags BEFORE `--`, host flags AFTER `--`. tmux runs it in the
	// pane; the env is applied via the tmux invocation's own environment.
	//
	// `env -u` unsets the NESTED-SESSION markers a parent Claude Code session exports
	// before exec'ing claude. When the harness itself runs inside a Claude session
	// (the common case: a teammate/CI agent drives `go test -tags live`), the child
	// claude inherits CLAUDECODE / CLAUDE_CODE_CHILD_SESSION / CLAUDE_CODE_SESSION_ID /
	// CLAUDE_CODE_AGENT / CLAUDE_CODE_ENTRYPOINT / CLAUDE_CODE_EXECPATH and
	// self-identifies as a NESTED session — which, since CC 2.1.170, SUPPRESSES the
	// on-disk conversation transcript (the "sessions launched from a shell that
	// inherited Claude Code env vars don't save transcripts" fix). With the markers
	// present the FO writes teams/<name>/config.json but NO projects/<cwd>/<id>.jsonl,
	// so the idle gate, the stream drain, and the teardown grade have nothing to read.
	// Unsetting them restores transcript persistence. tmux -e cannot do this (it ADDS
	// per-session vars, it cannot REMOVE a var the tmux server inherited), so the unset
	// rides on the launch command itself. The team flag
	// (CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS) and the auth/config vars are NOT unset —
	// they must reach the child for team tools + isolated transcript path.
	launch := shellJoin(append([]string{"env"}, unsetNestedSessionArgs(d.binary, "claude",
		"--plugin-dir", d.pluginDir,
		"--skip-compat-check",
		"--",
		"--permission-mode", "bypassPermissions",
		"--model", d.modelName,
	)...))
	// The child env MUST ride on per-session `-e KEY=VAL` flags, NOT the exec.Cmd's
	// own Env: `tmux new-session` against a PRE-EXISTING tmux server inherits the
	// SERVER's global environment and SILENTLY drops the command process's env — so
	// HOME/CLAUDE_CONFIG_DIR/the OAuth token would leak to the operator's real values
	// and claude would write its transcript outside the isolated tree (the harness
	// would then never find the session jsonl). `-e` sets per-session environment
	// that overrides the server env and is robust to an already-running server.
	args := []string{"new-session", "-d", "-s", session, "-x", "220", "-y", "50", "-c", workflowRoot}
	for _, kv := range env {
		args = append(args, "-e", kv)
	}
	args = append(args, launch)
	newSession := exec.Command("tmux", args...)
	if out, err := newSession.CombinedOutput(); err != nil {
		proc.kill()
		t.Fatalf("tmux new-session for %s failed: %v\n%s", label, err, out)
	}

	started := time.Now()

	// The session jsonl appears under CLAUDE_CONFIG_DIR/projects/<encoded-cwd> once
	// the FO's session is created — keyed by the same symlink-resolved cwd the seed
	// used. Resolve it BEFORE the idle gate: the gate reads the on-disk transcript for
	// a turn-end signal, not pane text (the render-fragile capture-pane scrape that
	// stalled 0/3 on a headless CI runner with missing terminfo).
	projectsDir := filepath.Join(effectiveConfigDir, "projects", encodeProjectDir(resolvedCwd))
	t.Logf("[pty %s] polling for session jsonl under: %s", label, projectsDir)
	sessionFile, err := waitForSessionFile(projectsDir, proc, ptyBootBudget)
	if err != nil {
		dumpPane := d.captureFOPane(session)
		_ = os.WriteFile(filepath.Join(artifactDir, "pane-no-jsonl.txt"), []byte(dumpPane), 0o644)
		proc.kill()
		t.Fatalf("session jsonl never appeared under %s for %s: %v\nFO pane:\n%s\nArtifacts: %s",
			projectsDir, label, err, dumpPane, artifactDir)
	}

	// Capture the FO's OWN session UUID and pin to it, so a teammate transcript
	// spawned later never lures the read off the FO (F30 — observed: the tail read the
	// comm-officer transcript and missed TeamCreate). The FO/root transcript is the
	// flat {FO-session-id}.jsonl whose FIRST sessionId-bearing entry carries no
	// top-level `agentName` (teammate/subagent transcripts live in the sibling
	// {FO-session-id}/subagents/ subdir; the FO root MAY itself carry a later
	// `agent-name` entry naming spacedock:first-officer, which the first-entry oracle
	// does not mistake for a teammate). At boot only the FO root file exists.
	foSessionID := waitForFOSessionID(projectsDir, proc, ptyBootBudget)
	if foSessionID == "" {
		dumpPane := d.captureFOPane(session)
		_ = os.WriteFile(filepath.Join(artifactDir, "pane-no-session-id.txt"), []byte(dumpPane), 0o644)
		proc.kill()
		t.Fatalf("could not read the FO root session id under %s for %s (first transcript seen: %s): the tail cannot pin to the FO without it\nFO pane:\n%s\nArtifacts: %s",
			projectsDir, label, sessionFile, dumpPane, artifactDir)
	}
	t.Logf("[pty %s] pinned FO session id: %s", label, foSessionID)

	// Gate the first send on the FO's boot turn reaching a committed idle (the
	// queued-keystroke hazard: a send mid-boot queues behind the boot turn and the
	// drive stalls). transcriptReachedIdle reads the pinned FO transcript for a
	// turn-end signal — render-INDEPENDENT, unlike the old pane-text scrape. The
	// 3-poll debounce covers the residual race where the transcript flips to end_turn
	// a beat before the TUI input box accepts keys.
	if err := d.waitTranscriptIdle(projectsDir, foSessionID, proc, ptyBootBudget); err != nil {
		dumpPane := d.captureFOPane(session)
		_ = os.WriteFile(filepath.Join(artifactDir, "pane-at-stall.txt"), []byte(dumpPane), 0o644)
		proc.kill()
		t.Fatalf("FO never reached a committed idle turn for %s within %s: %v\nFO pane:\n%s\nArtifacts: %s",
			label, ptyBootBudget, err, dumpPane, artifactDir)
	}

	if err := d.sendToFO(session, prompt); err != nil {
		proc.kill()
		t.Fatalf("send prompt to FO for %s: %v\nArtifacts: %s", label, err, artifactDir)
	}

	return ptySession{
		tmuxName:    session,
		projectsDir: projectsDir,
		foSessionID: foSessionID,
		configDir:   effectiveConfigDir,
		proc:        proc,
		artifactDir: artifactDir,
		started:     started,
	}
}

// run satisfies liveDriver: it launches+sends, then drains the session jsonl to
// session-death-or-quiet-stall via the EXISTING streamWatcher (one liveness
// mechanism, no second impl), captures the FO pane as finalMessage, and tears down
// the tmux session. The shared assertions consume the result oblivious to the
// transport.
func (d ptyLiveDriver) run(t *testing.T, scenario sharedRuntimeScenario, workflowRoot, prompt string) liveResult {
	t.Helper()
	session := d.launchAndSend(t, scenario.name, workflowRoot, prompt+" "+antiShutdownOverride)
	defer session.proc.kill()

	watcher := newStreamWatcher(session.newFileSource(), session.proc, discardStreamLine)
	stream, _ := watcher.drainToExit(quietBudgetDispatchClose, "pty "+scenario.name)
	duration := time.Since(session.started)

	finalMessage := d.captureFOPane(session.tmuxName)

	if writeErr := os.WriteFile(filepath.Join(session.artifactDir, "pty-stream.jsonl"), []byte(stream), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	_ = os.WriteFile(filepath.Join(session.artifactDir, "fo-pane-final.txt"), []byte(finalMessage), 0o644)

	return liveResult{
		finalMessage: finalMessage,
		stream:       stream,
		artifactDir:  session.artifactDir,
		duration:     duration,
	}
}

// foTarget resolves the `<session>.<index>` send/capture target for the FO pane by
// title. It errors rather than mis-route to a dispatched-ensign sibling pane.
func (d ptyLiveDriver) foTarget(session string) (string, error) {
	out, err := exec.Command("tmux", "list-panes", "-t", session,
		"-F", "#{pane_index}\t#{pane_active}\t#{pane_title}").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("tmux list-panes: %w (%s)", err, out)
	}
	idx, err := firstOfficerPaneIndex(string(out))
	if err != nil {
		return "", err
	}
	return session + "." + idx, nil
}

// captureFOPane returns the FO pane's visible text, or "" if it cannot resolve the
// pane (used for diagnostics, never as a grade signal).
func (d ptyLiveDriver) captureFOPane(session string) string {
	target, err := d.foTarget(session)
	if err != nil {
		return ""
	}
	out, err := exec.Command("tmux", "capture-pane", "-p", "-t", target).CombinedOutput()
	if err != nil {
		return ""
	}
	return string(out)
}

// waitTranscriptIdle blocks until the pinned FO transcript reads a committed idle
// turn (transcriptReachedIdle) for ptyIdleStablePolls consecutive reads, bounded by
// budget. The signal is the on-disk JSONL turn-end — render-INDEPENDENT, replacing
// the pane-text scrape that stalled on a headless CI runner with missing terminfo
// (blank capture-pane → never idle). It errors if the session dies or the budget
// elapses without a stable idle. The consecutive-read debounce covers the residual
// race where the transcript flips to end_turn a beat before the TUI accepts keys.
func (d ptyLiveDriver) waitTranscriptIdle(projectsDir, sessionID string, proc procPoller, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	consecutive := 0
	for {
		idle := false
		if path := sessionFileByID(projectsDir, sessionID); path != "" {
			if data, err := os.ReadFile(path); err == nil {
				idle = transcriptReachedIdle(strings.Split(string(data), "\n"))
			}
		}
		if idle {
			consecutive++
			if consecutive >= ptyIdleStablePolls {
				return nil
			}
		} else {
			consecutive = 0
		}
		if _, exited := proc.poll(); exited {
			return fmt.Errorf("tmux session died before the FO reached a committed idle turn")
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("no committed idle turn within %s (consecutive idle polls=%d/%d)", budget, consecutive, ptyIdleStablePolls)
		}
		time.Sleep(ptyPollInterval)
	}
}

// readPinnedTranscript returns the FO transcript lines (resolved by the pinned
// session id), or nil if not yet present.
func (s ptySession) readPinnedTranscript() []string {
	path := sessionFileByID(s.projectsDir, s.foSessionID)
	if path == "" {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return strings.Split(string(data), "\n")
}

// nudgePastGreet drives the interactive FO past its CONTRACTUAL greet-stop (item 9:
// an interactive FO presents the summary then STOPS for input — correct behavior,
// not a bug), exactly as a captain would: it sends a "go + conn" message and waits
// for the next committed idle, re-checking the pinned transcript for signalSeen,
// bounded by maxNudges. It returns nil once signalSeen is true (the FO acted), or an
// error if the cap is hit or the session dies. Each nudge waits for a FRESH committed
// turn (the transcript must grow past the pre-nudge length AND read idle), so a
// stale pre-nudge idle never counts as the nudge's response.
func (d ptyLiveDriver) nudgePastGreet(session ptySession, nudgeText string, signalSeen func([]string) bool, maxNudges int, perNudgeBudget time.Duration) error {
	for attempt := 1; attempt <= maxNudges; attempt++ {
		if signalSeen(session.readPinnedTranscript()) {
			return nil
		}
		preLen := len(session.readPinnedTranscript())
		if err := d.sendToFO(session.tmuxName, nudgeText); err != nil {
			return fmt.Errorf("nudge %d/%d send failed: %w", attempt, maxNudges, err)
		}
		// Wait for a FRESH committed turn: the transcript grows past preLen AND reads
		// idle again. A bare idle check would pass immediately on the stale pre-nudge
		// idle without the FO having processed the nudge.
		deadline := time.Now().Add(perNudgeBudget)
		for {
			lines := session.readPinnedTranscript()
			if len(lines) > preLen && transcriptReachedIdle(lines) {
				break
			}
			if signalSeen(lines) {
				return nil
			}
			if _, exited := session.proc.poll(); exited {
				return fmt.Errorf("tmux session died during nudge %d/%d", attempt, maxNudges)
			}
			if time.Now().After(deadline) {
				break // this nudge produced no fresh idle; the outer loop re-nudges
			}
			time.Sleep(ptyPollInterval)
		}
	}
	if signalSeen(session.readPinnedTranscript()) {
		return nil
	}
	return fmt.Errorf("the FO did not produce the awaited signal after %d captain nudges", maxNudges)
}

// sendToFO types text into the title-resolved FO pane via send-keys, then Enter.
// The text is sent with `-l` (literal) so claude's input box receives it verbatim;
// a brief pause before Enter avoids the send-keys race where a fast Enter is
// swallowed before the input widget commits the paste.
func (d ptyLiveDriver) sendToFO(session, text string) error {
	target, err := d.foTarget(session)
	if err != nil {
		return err
	}
	if out, err := exec.Command("tmux", "send-keys", "-t", target, "-l", text).CombinedOutput(); err != nil {
		return fmt.Errorf("tmux send-keys (literal): %w (%s)", err, out)
	}
	time.Sleep(300 * time.Millisecond)
	if out, err := exec.Command("tmux", "send-keys", "-t", target, "Enter").CombinedOutput(); err != nil {
		return fmt.Errorf("tmux send-keys (enter): %w (%s)", err, out)
	}
	return nil
}

// --- tmux session as a procPoller ----------------------------------------

// tmuxLiveProc adapts a tmux session to the streamWatcher's procPoller surface:
// poll() reports the session gone (interactive claude does not stream-json to a
// pipe, so "exited" means the tmux session died), kill() tears it down. The
// deferred kill in run() is the launcher reaping the resident session.
type tmuxLiveProc struct {
	session string
	mu      sync.Mutex
	killed  bool
}

func newTmuxLiveProc(session string) *tmuxLiveProc {
	return &tmuxLiveProc{session: session}
}

func (p *tmuxLiveProc) poll() (int, bool) {
	if tmuxHasSession(p.session) {
		return 0, false
	}
	return 0, true
}

func (p *tmuxLiveProc) kill() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.killed {
		return
	}
	p.killed = true
	_ = exec.Command("tmux", "kill-session", "-t", p.session).Run()
}

func tmuxHasSession(session string) bool {
	return exec.Command("tmux", "has-session", "-t", session).Run() == nil
}

// --- session jsonl resolution waiters (live-only: budget + procPoller) -----

// waitForSessionFile blocks until a session jsonl appears in dir (the FO's own
// transcript), bounded by budget. Returns the newest assistant-bearing file, or
// the newest file overall before the first assistant turn. Errors if the tmux
// session dies before any transcript appears.
func waitForSessionFile(dir string, proc procPoller, budget time.Duration) (string, error) {
	deadline := time.Now().Add(budget)
	for {
		if path := activeSessionFile(dir); path != "" {
			return path, nil
		}
		if _, exited := proc.poll(); exited {
			return "", fmt.Errorf("tmux session died before a session jsonl appeared under %s", dir)
		}
		if time.Now().After(deadline) {
			return "", fmt.Errorf("timed out after %s waiting for a session jsonl under %s", budget, dir)
		}
		time.Sleep(ptyPollInterval)
	}
}

// waitForFOSessionID polls dir for the FO's own session UUID, bounded by budget.
// The FO is the team ROOT, not a spawned teammate: foRootSessionID identifies it by
// the first-entry oracle (the flat transcript whose FIRST sessionId-bearing entry
// carries no top-level `agentName`). Keying on the root transcript (not the
// newest-assistant-bearing file, which can already be a teammate's when the FO
// dispatched during boot) pins the id to the FO. The id can lag the file by a tick,
// so this retries; "" if the session dies or the budget elapses.
func waitForFOSessionID(dir string, proc procPoller, budget time.Duration) string {
	deadline := time.Now().Add(budget)
	for {
		if id := foRootSessionID(dir); id != "" {
			return id
		}
		if _, exited := proc.poll(); exited {
			return ""
		}
		if time.Now().After(deadline) {
			return ""
		}
		time.Sleep(ptyPollInterval)
	}
}

// (session-jsonl resolution helpers — fileLineSource, sessionFileByID,
// sessionIDFromFile, foRootSessionID, activeSessionFile, encodeProjectDir — are pure
// and live in pty_session_test.go so they compile + are unit-tested OFFLINE.)

// configDirOrDefault returns the resolved CLAUDE_CONFIG_DIR, or the isolated-HOME
// default (home/.claude) when the env did not set one — the mirror of
// resolveClaudeConfigDir's local default.
func configDirOrDefault(configDir, home string) string {
	if configDir != "" {
		return configDir
	}
	return filepath.Join(home, ".claude")
}

// seedInteractiveClaudeConfig writes <configDir>/.claude.json with the keys that
// make an INTERACTIVE Claude Code launch skip its first-run dialogs:
//   - hasCompletedOnboarding/lastOnboardingVersion/theme: skip the theme + login
//     pickers a fresh isolated HOME otherwise stalls at;
//   - bypassPermissionsModeAccepted: pre-accept the bypass-permissions mode prompt;
//   - projects[resolvedCwd].hasTrustDialogAccepted: skip the per-project
//     "trust this folder?" dialog (keyed by the symlink-resolved cwd).
//
// It MERGES onto any existing config (a re-launch into a dir claude already wrote)
// so it never clobbers session state, only ensures the dialog-skip keys are set.
// When the file is absent (the normal per-run case) it creates it. The credential
// itself rides in the env (CLAUDE_CODE_OAUTH_TOKEN / ANTHROPIC_API_KEY via
// isolatedClaudeEnv); this only clears the interactive onboarding gate.
func seedInteractiveClaudeConfig(configDir, resolvedCwd string) error {
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(configDir, ".claude.json")
	cfg := map[string]any{}
	if raw, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(raw, &cfg) // tolerate a partial/empty file; the keys below are forced
	}
	cfg["hasCompletedOnboarding"] = true
	cfg["lastOnboardingVersion"] = "1.0.53"
	cfg["theme"] = "dark"
	cfg["bypassPermissionsModeAccepted"] = true

	projects, _ := cfg["projects"].(map[string]any)
	if projects == nil {
		projects = map[string]any{}
	}
	entry, _ := projects[resolvedCwd].(map[string]any)
	if entry == nil {
		entry = map[string]any{}
	}
	entry["hasTrustDialogAccepted"] = true
	entry["projectOnboardingSeenCount"] = 1
	projects[resolvedCwd] = entry
	cfg["projects"] = projects

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o600)
}

// seedStoredLoginCredential writes <configDir>/.credentials.json from the
// operator's stored-login OAuth credential so the INTERACTIVE child authenticates.
// The interactive TUI reads this file (the `{"claudeAiOauth": {...}}` object),
// unlike headless `-p`, which consumes CLAUDE_CODE_OAUTH_TOKEN directly. The
// credential lives in the OS credential store (macOS keychain item
// `Claude Code-credentials`), whose `-w` output is ALREADY the full
// `{"claudeAiOauth": {...}}` shape claude expects, so it is written verbatim.
//
// It is best-effort + operator-local: it returns an error (the caller logs and
// proceeds) when the credential store is unreachable — a non-macOS host (CI Linux),
// no `security` binary, no stored login, or an empty/malformed payload. In every
// such case the seed is a no-op and the child falls back to the env-token /
// ANTHROPIC_API_KEY path, so AC-6's skip-not-fatal-without-auth (decided upstream in
// isolatedClaudeEnv) is unaffected.
func seedStoredLoginCredential(configDir string) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("stored-login seed is macOS-keychain only (GOOS=%s)", runtime.GOOS)
	}
	raw, err := exec.Command("security", "find-generic-password",
		"-s", "Claude Code-credentials", "-w").Output()
	if err != nil {
		return fmt.Errorf("read keychain credential: %w", err)
	}
	cred := strings.TrimSpace(string(raw))
	// Confirm the payload is the expected stored-login shape before seeding, so a
	// malformed/empty keychain entry falls through to the env-credential path rather
	// than writing a .credentials.json that 401s the same way the env token does.
	var probe struct {
		ClaudeAIOauth json.RawMessage `json:"claudeAiOauth"`
	}
	if cred == "" || json.Unmarshal([]byte(cred), &probe) != nil || len(probe.ClaudeAIOauth) == 0 {
		return fmt.Errorf("keychain credential is empty or not the {claudeAiOauth:{...}} shape")
	}
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(configDir, ".credentials.json"), []byte(cred), 0o600)
}

// shellJoin renders argv as a single shell command line, single-quoting any arg
// that needs it, so tmux new-session runs the exact launch the `-p` runner builds.
func shellJoin(argv []string) string {
	parts := make([]string, 0, len(argv))
	for _, a := range argv {
		parts = append(parts, shellQuote(a))
	}
	return strings.Join(parts, " ")
}

func shellQuote(s string) string {
	if s != "" && !strings.ContainsAny(s, " \t\n'\"\\$`*?[]{}()&;|<>~#") {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
