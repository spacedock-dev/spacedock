// ABOUTME: The Claude runtime seam — owns the ~/.claude/teams reads behind a
// ABOUTME: host-supplied TeamStateProbe plus the boot hint + bare-mode advisory text.
package claudeteam

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// teamStateWindow is the lookback for "recent team-runtime evidence": a
// ~/.claude/teams/*/config.json touched inside this window means a Claude team
// session is live. The boot TEAM_STATE and the build bare-mode advisory share it.
const teamStateWindow = 30 * time.Minute

// PresentFalseHint is the boot TEAM_STATE present:false hint the Claude seam
// supplies. It names a Claude-only tool (TeamCreate), so it lives here, not in
// the generic internal/status renderer — the generic renderer reads it only when
// a probe is wired, and emits a host-neutral line when the probe is nil.
const PresentFalseHint = "dispatch a named background Agent for team-mode messaging (no TeamCreate needed)"

// TeamStateProbe reports recent local team-runtime evidence over the shared
// ~/.claude/teams read. present drives the boot TEAM_STATE present field; hint is
// the boot present:true hint line; recent drives the build bare-mode advisory
// gate. now is injected so the 30-minute window is testable. internal/status and
// internal/dispatch take this as a value (nil on a non-Claude host) so their
// source carries no ~/.claude read.
type TeamStateProbe func(home string, now time.Time) (present bool, hint string, recent bool)

// Probe is the concrete Claude implementation: it scans ~/.claude/teams/*/config.json
// mtimes under home for the newest one inside the 30-minute window. present and
// recent both report whether such a config exists; hint names the newest team
// directory on the present path. home is the resolved HOME (the caller keeps HOME
// resolution generic; only this ~/.claude read is Claude-specific). The Claude CLI
// front door wires claudeteam.Probe; Codex/bare wire nil.
func Probe(home string, now time.Time) (present bool, hint string, recent bool) {
	teamsDir := filepath.Join(home, ".claude", "teams")
	info, err := os.Stat(teamsDir)
	if err != nil || !info.IsDir() {
		return false, "", false
	}
	entries, err := os.ReadDir(teamsDir)
	if err != nil {
		return false, "", false
	}
	cutoff := now.Add(-teamStateWindow)
	var newest string
	var newestMtime time.Time
	for _, ent := range entries {
		cfg := filepath.Join(teamsDir, ent.Name(), "config.json")
		st, err := os.Stat(cfg)
		if err != nil || !st.Mode().IsRegular() {
			continue
		}
		if st.ModTime().After(newestMtime) {
			newestMtime = st.ModTime()
			newest = ent.Name()
		}
	}
	if newest != "" && !newestMtime.Before(cutoff) {
		return true, "recent team directory: " + newest, true
	}
	return false, "", false
}

// TranscriptProbe resolves session id's transcript file path under this host's
// session-transcript store, or "" when unresolvable. Used by the boot guard's
// receipt write (force-boot-at-compaction-boundary): internal/status takes this
// as a value (nil on a non-Claude host) so its source carries no ~/.claude read,
// mirroring TeamStateProbe.
type TranscriptProbe func(home, sessionID string) string

// TranscriptPath is the concrete Claude implementation: Claude Code names every
// session transcript `{session_id}.jsonl` under a per-project directory whose
// own name is a host-internal encoding of the launch cwd (captured live:
// force-boot-at-compaction-boundary/ideation-spike-evidence.md §1-3); rather
// than reconstruct that encoding, this globs on the leaf filename alone, which
// is unique to the session by construction.
func TranscriptPath(home, sessionID string) string {
	if home == "" || sessionID == "" {
		return ""
	}
	matches, err := filepath.Glob(filepath.Join(home, ".claude", "projects", "*", sessionID+".jsonl"))
	if err != nil || len(matches) == 0 {
		return ""
	}
	sort.Strings(matches)
	return matches[len(matches)-1]
}

// BareModeAdvisory writes the bare-mode dispatch warning to w. It names a
// Claude-only bootstrap tool (TeamCreate) and the ~/.claude/teams path, so the
// text lives in the Claude seam, not in the generic internal/dispatch build path.
// The generic build path calls this only when a Claude probe is wired AND reports
// no recent evidence; a nil-probe (Codex/bare) host emits no advisory at all.
func BareModeAdvisory(w io.Writer) {
	fmt.Fprintln(w,
		"WARN: bare_mode dispatch with no recent TeamCreate evidence "+
			"(no ~/.claude/teams/*/config.json modified in the last 30 minutes). "+
			"If you intend teams mode, ensure SendMessage is available (ToolSearch select:SendMessage) and dispatch a named background Agent — no TeamCreate. "+
			"If bare is intentional, this warning can be ignored.")
}

// LegacyTeamNameAdvisory writes the --team-name-on-claude warning to w. It names
// Claude-only concepts (the auto-team default, the legacy TeamCreate registry), so
// the text lives in the Claude seam beside BareModeAdvisory. The generic build path
// calls it when a non-bare claude dispatch passes a team_name — the teamName != ""
// complement of merged mode — where auto-team is the default and the explicit
// --team-name selects the sunsetting legacy shape. Unlike BareModeAdvisory it needs
// no probe: the trigger is wholly in the CLI args. Stderr-only; the dispatch
// envelope is untouched.
func LegacyTeamNameAdvisory(w io.Writer) {
	fmt.Fprintln(w,
		"WARN: --team-name selects the legacy TeamCreate-registry dispatch shape "+
			"(team_name present, run_in_background absent). On host=claude, auto-team "+
			"is the default — omit --team-name to emit the auto-team shape (name + "+
			"run_in_background, no team_name). If you mean the legacy team-registry "+
			"path, this warning can be ignored.")
}
