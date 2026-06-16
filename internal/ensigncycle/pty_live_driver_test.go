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
	"sort"
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
//     tmux capture-pane and the on-disk session jsonl, never a pipe (a pipe makes
//     stdout a non-tty and the interactive TUI never renders).
//   - resolve the FO pane by #{pane_title}, not the active pane: an Agent dispatch
//     materializes the ensign as a sibling pane that becomes active, so the FO's
//     TeamDelete/marker live in the title-resolved FO pane (firstOfficerPaneIndex).
//   - gate the first send on stable-idle: keystrokes sent mid-boot QUEUE behind the
//     boot turn and the drive stalls; waitStableIdle holds the send until the FO
//     pane shows no busy marker for N consecutive polls.
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
	// stable-idle (boot complete, ready for the first keystroke). A team-mode boot
	// + FO contract load is multi-minute, so this is roomier than a per-step budget.
	ptyBootBudget = 4 * time.Minute
	// ptyIdleStablePolls is how many CONSECUTIVE no-busy-marker pane reads gate the
	// first send — the spike's queued-keystroke hazard (a send mid-boot queues
	// behind the boot turn). N>1 avoids a transient between-turn blank.
	ptyIdleStablePolls = 3
	// ptyPollInterval paces the pane/idle polls.
	ptyPollInterval = 2 * time.Second
)

// ptyBusyMarkers are pane-text fragments Claude Code's TUI renders WHILE a turn is
// in flight. Their ABSENCE (plus a non-empty pane) is the stable-idle signal the
// first send gates on. "esc to interrupt" is the load-bearing one — it is present
// for the whole duration of any running turn and gone when the FO is idle.
var ptyBusyMarkers = []string{"esc to interrupt", "Press up to edit queued"}

// ptySession is a launched, prompted interactive session: the tmux name, the FO's
// session jsonl on disk, the proc poller over the tmux session, and the artifact
// dir. The shared-scenario run() and the team-mode tests both build one via
// launchAndSend, then either drain it whole (run) or watch it incrementally (the
// team-mode tests' expect/grade beats — the faithful resurrection of the retired
// forced-team tests).
type ptySession struct {
	tmuxName    string
	sessionFile string
	configDir   string // the effective CLAUDE_CONFIG_DIR this run wrote under (teams/ + projects/ live here)
	proc        *tmuxLiveProc
	artifactDir string
	started     time.Time
}

// newFileSource opens a fresh tail over the FO's session jsonl from byte 0, so a
// caller's streamWatcher sees the whole transcript from launch.
func (s ptySession) newFileSource() *fileLineSource { return newFileLineSource(s.sessionFile) }

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
	// FO loads — and the stable-idle gate would falsely fire into that picker, since
	// a picker shows no busy marker. The headless `-p` suite never hits this (non-
	// interactive `-p` consumes the OAuth token directly); the interactive transport
	// must pre-clear onboarding itself. This is the implementation hazard the spike
	// missed (it ran on already-onboarded operator auth).
	if err := seedInteractiveClaudeConfig(effectiveConfigDir, resolvedCwd); err != nil {
		t.Fatalf("seed interactive claude config for %s: %v", label, err)
	}

	session := fmt.Sprintf("sdpty-%s-%d", label, time.Now().UnixNano())
	proc := newTmuxLiveProc(session)

	// The launch line is the live runner's exact front door minus `-p`:
	// spacedock-owned flags BEFORE `--`, host flags AFTER `--`. tmux runs it in the
	// pane; the env is applied via the tmux invocation's own environment.
	launch := shellJoin([]string{d.binary, "claude",
		"--plugin-dir", d.pluginDir,
		"--skip-contract-check",
		"--",
		"--permission-mode", "bypassPermissions",
		"--model", d.modelName,
	})
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

	if err := d.waitStableIdle(session, ptyBootBudget); err != nil {
		dumpPane := d.captureFOPane(session)
		_ = os.WriteFile(filepath.Join(artifactDir, "pane-at-stall.txt"), []byte(dumpPane), 0o644)
		proc.kill()
		t.Fatalf("FO never reached stable-idle for %s within %s: %v\nFO pane:\n%s\nArtifacts: %s",
			label, ptyBootBudget, err, dumpPane, artifactDir)
	}

	if err := d.sendToFO(session, prompt); err != nil {
		proc.kill()
		t.Fatalf("send prompt to FO for %s: %v\nArtifacts: %s", label, err, artifactDir)
	}

	// The session jsonl appears under CLAUDE_CONFIG_DIR/projects/<encoded-cwd> once
	// the FO's session is created — keyed by the same symlink-resolved cwd the seed
	// used.
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

	return ptySession{
		tmuxName:    session,
		sessionFile: sessionFile,
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

// waitStableIdle blocks until the FO pane shows no busy marker for
// ptyIdleStablePolls consecutive reads (the stable-idle gate), bounded by budget.
// It errors if the session dies or the budget elapses without a stable idle.
func (d ptyLiveDriver) waitStableIdle(session string, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	consecutive := 0
	for {
		if !tmuxHasSession(session) {
			return fmt.Errorf("tmux session %q died before reaching stable-idle", session)
		}
		pane, err := d.foTarget(session)
		idle := false
		if err == nil {
			if out, capErr := exec.Command("tmux", "capture-pane", "-p", "-t", pane).CombinedOutput(); capErr == nil {
				idle = paneIsStableIdle(string(out))
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
		if time.Now().After(deadline) {
			return fmt.Errorf("no stable-idle within %s (consecutive idle polls=%d/%d)", budget, consecutive, ptyIdleStablePolls)
		}
		time.Sleep(ptyPollInterval)
	}
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

// paneIsStableIdle reports whether a captured pane is at a ready-for-input idle: a
// non-empty pane carrying NO busy marker. The busy markers ("esc to interrupt",
// the queued-message banner) are present for the whole duration of any running
// turn; their absence over a non-empty pane is the stable-idle signal.
func paneIsStableIdle(pane string) bool {
	if strings.TrimSpace(pane) == "" {
		return false
	}
	for _, m := range ptyBusyMarkers {
		if strings.Contains(pane, m) {
			return false
		}
	}
	return true
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

// --- session jsonl tail as a lineSource -----------------------------------

// fileLineSource is the pty lineSource: it tails a session jsonl file, returning
// the complete lines appended since the last drain. It is the file analog of
// pipeLineSource (which tails a stdout pipe): the interactive session writes its
// stream-json transcript to disk, not a pipe, so the driver tails the file for
// liveness while feeding the SAME streamWatcher. It advances a byte offset only
// past COMPLETE (newline-terminated) lines, so a half-written JSONL record split
// across two appends is held until its newline arrives and never parsed early.
type fileLineSource struct {
	path   string
	offset int64
}

func newFileLineSource(path string) *fileLineSource {
	return &fileLineSource{path: path}
}

func (s *fileLineSource) drain() []string {
	data, err := os.ReadFile(s.path)
	if err != nil || int64(len(data)) <= s.offset {
		return nil
	}
	tail := data[s.offset:]
	// Only consume up to the last newline; bytes after it are a partial trailing
	// line held for the next drain.
	last := strings.LastIndexByte(string(tail), '\n')
	if last < 0 {
		return nil
	}
	consumed := tail[:last+1]
	s.offset += int64(len(consumed))
	var out []string
	for _, line := range strings.Split(strings.TrimRight(string(consumed), "\n"), "\n") {
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

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

// activeSessionFile returns the live conversation transcript in dir: the newest
// *.jsonl that already carries an assistant entry, else the newest overall (right
// after launch, before the first assistant turn). Ported from spacedock-gym's
// ActiveSessionFile (reference-only, not imported).
func activeSessionFile(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	type cand struct {
		path    string
		mod     int64
		hasAsst bool
	}
	var cands []cand
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		path := filepath.Join(dir, e.Name())
		body, _ := os.ReadFile(path)
		cands = append(cands, cand{path, info.ModTime().UnixNano(), strings.Contains(string(body), `"type":"assistant"`)})
	}
	if len(cands) == 0 {
		return ""
	}
	sort.Slice(cands, func(i, j int) bool {
		if cands[i].hasAsst != cands[j].hasAsst {
			return cands[i].hasAsst
		}
		return cands[i].mod > cands[j].mod
	})
	return cands[0].path
}

// encodeProjectDir converts an absolute cwd into the directory name Claude Code
// uses under <config>/projects: each of '/', '.', and '_' becomes '-'. Ported from
// spacedock-gym's EncodeProjectDir (reference-only, not imported).
func encodeProjectDir(cwd string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '/', '.', '_':
			return '-'
		default:
			return r
		}
	}, cwd)
}

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
