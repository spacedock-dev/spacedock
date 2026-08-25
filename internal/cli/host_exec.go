// ABOUTME: Production hostOps — resolves the installed plugin manifest via
// ABOUTME: `claude/codex plugin list --json`, spawns the host, and shells installs.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// execHost backs hostOps with the real host CLIs and process exec.
type execHost struct{}

var _ hostOps = execHost{}

// pluginListEntry is the subset of `<host> plugin list --json` this binary
// reads: the `plugin@marketplace` id, the resolved install path, and whether the
// plugin is enabled in the host. (Observed schema: the entry carries `id`, not
// separate name/marketplace fields; `enabled` is distinct from `installPath` —
// an installed plugin can be present-but-disabled.)
type pluginListEntry struct {
	ID          string `json:"id"`
	InstallPath string `json:"installPath"`
	Enabled     bool   `json:"enabled"`
}

// ResolveManifest returns the installed spacedock@spacedock plugin manifest path
// for host, or "" (no error) when no plugin is installed. The two hosts resolve
// differently: Claude reports an installPath in `claude plugin list --json`;
// Codex's `plugin list --json` carries no install path (only an installed flag,
// different schema), so the Codex path confirms the install via the text listing
// and resolves the manifest under the deterministic Codex plugin cache.
func (e execHost) ResolveManifest(host string) (string, error) {
	if host == "codex" {
		return e.resolveCodexManifest()
	}
	return e.resolveClaudeManifest(host)
}

// resolveClaudeManifest shells `claude plugin list --json`, finds the channel's
// plugin id (`spacedock@spacedock` stable / `spacedock@spacedock-edge` edge — the
// binary's devBranch selects it), and returns <installPath>/.claude-plugin/
// plugin.json. Returns "" (no error) when the host reports no matching install or
// no installPath. Matching the channel id (not the hardcoded stable id) is what
// lets an edge binary recognize/refresh an installed edge plugin.
func (execHost) resolveClaudeManifest(host string) (string, error) {
	out, err := exec.Command(host, "plugin", "list", "--json").Output()
	if err != nil {
		return "", fmt.Errorf("%s plugin list --json: %w", host, err)
	}
	var entries []pluginListEntry
	if err := json.Unmarshal(out, &entries); err != nil {
		return "", fmt.Errorf("parse %s plugin list --json: %w", host, err)
	}
	id := channelPluginID(devBranch)
	for _, e := range entries {
		if e.ID == id {
			if e.InstallPath == "" {
				return "", nil
			}
			return filepath.Join(e.InstallPath, manifestSubpath(host)), nil
		}
	}
	return "", nil
}

// resolveCodexManifest confirms an installed spacedock plugin id via the text
// `codex plugin list` and resolves the manifest under the Codex plugin cache.
// TWO ids are checked, in this order: `codexLocalPluginID`
// (`spacedock@spacedock-local`, the --plugin-dir dev install, devBranch-
// independent by design — see installCodexLocalPluginDir) first, then the
// channel's own id (`spacedock@spacedock` stable / `spacedock@spacedock-edge`
// edge — the binary's devBranch selects it). Without checking the local id, a
// --plugin-dir install would resolve to NoPluginFound here (the gate reads only
// the channel id), and resolveHealableGate would then try to auto-install the
// real channel plugin on top of a working dev install. Codex's listing carries
// no install path (its `--json` form has one only for the marketplace root, not
// the cached plugin), so the deterministic cache layout is the resolver. Codex
// installs land at
// <CODEX_HOME>/plugins/cache/<marketplace>/<plugin>/<version>/.codex-plugin/plugin.json.
// Returns "" (no error) when neither id is installed or no cached manifest
// exists for the installed one yet.
func (execHost) resolveCodexManifest() (string, error) {
	out, err := exec.Command("codex", "plugin", "list").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("codex plugin list: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	listing := string(out)
	if codexEntryInstalled(listing, codexLocalPluginID) {
		return codexLocalCacheManifest()
	}
	if !codexEntryInstalled(listing, channelPluginID(devBranch)) {
		return "", nil
	}
	return codexCacheManifest()
}

// codexCacheManifest resolves the channel's cached manifest under the Codex plugin
// cache: <CODEX_HOME>/plugins/cache/<marketplace>/spacedock/<version>/.codex-plugin/
// plugin.json, where <marketplace> is the channel marketplace name (`spacedock`
// stable / `spacedock-edge` edge) and the entry dir is always `spacedock`. It picks
// the semver-greatest version dir and returns that manifest path, or "" (no error)
// when no cached manifest exists yet (absent cache root, no version dir, or the
// manifest file is missing).
func codexCacheManifest() (string, error) {
	return codexCacheManifestAt(channelMarketplace(devBranch), channelEntry(devBranch))
}

// codexLocalCacheManifest is codexCacheManifest's --plugin-dir counterpart: the
// dedicated local marketplace name (devBranch-independent), entry always
// `spacedock`.
func codexLocalCacheManifest() (string, error) {
	return codexCacheManifestAt(codexLocalMarketplaceName, "spacedock")
}

// codexCacheManifestAt resolves the cached manifest under
// <CODEX_HOME>/plugins/cache/<marketplace>/<entry>/, picking the semver-greatest
// version dir. Returns "" (no error) when no cached manifest exists yet (absent
// cache root, no version dir, or the manifest file is missing).
func codexCacheManifestAt(marketplace, entry string) (string, error) {
	cacheRoot := filepath.Join(codexHome(), "plugins", "cache", marketplace, entry)
	versionDir, err := latestVersionDir(cacheRoot)
	if err != nil || versionDir == "" {
		return "", nil
	}
	manifest := filepath.Join(versionDir, manifestSubpath("codex"))
	if _, statErr := os.Stat(manifest); statErr != nil {
		return "", nil
	}
	return manifest, nil
}

// codexEntryInstalled reports whether the `codex plugin list` text output marks
// the given plugin id as installed. The listing renders a column-aligned table
// (PLUGIN STATUS VERSION PATH header) with one space-padded row per plugin as
// `<id>  <status>  <ver>  <path>`, where an installed row's status is
// `installed,` (or `installed`) and a not-installed row's is `not installed`.
// (Older codex rendered the status in parens — `<id> (installed[, enabled])` —
// which this still tolerates.) The match is field-based: the id must be a
// whitespace-delimited field, and the next field, stripped of surrounding `()`
// and a trailing `,`, must equal `installed` exactly. Field equality (not a
// substring scan) rejects a marketplace PATH cell that contains the bare
// marketplace word; reading the next field rejects `not installed`.
func codexEntryInstalled(listing, id string) bool {
	for _, line := range strings.Split(listing, "\n") {
		fields := strings.Fields(line)
		for i, f := range fields {
			if f != id {
				continue
			}
			if i+1 < len(fields) && strings.Trim(fields[i+1], "(),") == "installed" {
				return true
			}
		}
	}
	return false
}

// codexHome returns the Codex config/cache root: $CODEX_HOME when set, else
// ~/.codex (matching the Codex CLI's own resolution).
func codexHome() string {
	if h := os.Getenv("CODEX_HOME"); h != "" {
		return h
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".codex"
	}
	return filepath.Join(home, ".codex")
}

// latestVersionDir returns the semver-greatest immediate subdirectory of root
// (the installed plugin's version dir). Returns "" (no error) when root is absent
// or has no subdirectories. Codex installs a single version, but a stale cache
// may hold several; a semver compare picks the most recent install — a lexical
// compare would wrongly order `0.10.0` before `0.9.0`.
func latestVersionDir(root string) (string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	latest := ""
	for _, e := range entries {
		if e.IsDir() && (latest == "" || compareVersion(e.Name(), latest) > 0) {
			latest = e.Name()
		}
	}
	if latest == "" {
		return "", nil
	}
	return filepath.Join(root, latest), nil
}

// compareVersion orders dotted plugin version names (e.g. `0.12.1`) numerically
// per component, so `0.10.0` sorts after `0.9.0`. It returns -1, 0, or 1. A
// component that does not parse as an integer falls back to a lexical compare of
// that component, so non-numeric names still order deterministically.
func compareVersion(a, b string) int {
	as, bs := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(as) || i < len(bs); i++ {
		var av, bv string
		if i < len(as) {
			av = as[i]
		}
		if i < len(bs) {
			bv = bs[i]
		}
		an, aerr := strconv.Atoi(av)
		bn, berr := strconv.Atoi(bv)
		if aerr == nil && berr == nil {
			if an != bn {
				if an < bn {
					return -1
				}
				return 1
			}
			continue
		}
		if av != bv {
			if av < bv {
				return -1
			}
			return 1
		}
	}
	return 0
}

// manifestSubpath returns the per-host manifest location under an install root.
func manifestSubpath(host string) string {
	if host == "codex" {
		return filepath.Join(".codex-plugin", "plugin.json")
	}
	return filepath.Join(".claude-plugin", "plugin.json")
}

// channelEntry is the marketplace ENTRY name the binary installs. It is always
// `spacedock` — equal to the plugin's own plugin.json `name` — on every channel,
// because the host (codex confirmed, claude by construction) rejects a marketplace
// entry whose name differs from the manifest name. The channel is NOT the entry
// name; it is the marketplace name (see channelMarketplace).
func channelEntry(devBranch string) string {
	return "spacedock"
}

// channelMarketplace maps the binary's devBranch stamp to the marketplace NAME the
// channel resolves from: a stable binary (devBranch=main) resolves the `spacedock`
// marketplace; an edge binary (any other devBranch, e.g. next) resolves
// `spacedock-edge`. The standalone marketplace repo exposes both names as distinct
// marketplace.json sources, each with the single `spacedock` entry pinned to that
// channel's version (stable to a release tag, edge tracking next HEAD). Encoding
// the channel here — not in the entry name — keeps the entry equal to the manifest
// `name`, so the host name-match passes on both channels.
func channelMarketplace(devBranch string) string {
	if devBranch == "main" {
		return "spacedock"
	}
	return "spacedock-edge"
}

// channelPluginID is the `<entry>@<marketplace>` plugin id for the channel devBranch
// selects: `spacedock@spacedock` (stable) or `spacedock@spacedock-edge` (edge). The
// entry (before the `@`) is always `spacedock`; the suffix is the channel
// marketplace NAME (the marketplace.json `name`).
func channelPluginID(devBranch string) string {
	return channelEntry(devBranch) + "@" + channelMarketplace(devBranch)
}

// codexLocalMarketplaceName is the FIXED marketplace name a codex `--plugin-dir`
// dev install registers under — distinct from either real channel's marketplace
// name (`spacedock` / `spacedock-edge`) and independent of the binary's
// devBranch. A dedicated name removes, by construction, the collision a shared
// channel name used to create: codex refuses `plugin marketplace add` when a
// name is already registered from a DIFFERENT source (measured live: exit 1,
// "already added from a different source; remove it before adding this
// source"), which a real channel install (git source) and a local dev install
// (a staged local path) would otherwise both trip over the SAME channel name.
// Captain-ruled fix, option (c) — see the entity's "Design change" note.
const codexLocalMarketplaceName = "spacedock-local"

// codexLocalPluginID is the plugin id a codex `--plugin-dir` dev install always
// answers to, on every channel — the entry stays `spacedock` (equal to the
// manifest name, same as every channel); only the marketplace name is
// dedicated.
const codexLocalPluginID = "spacedock@" + codexLocalMarketplaceName

// otherChannelMarketplace returns the marketplace NAME of the Spacedock channel
// devBranch does NOT select — the sibling channel a codex install must also stay
// clear of (Codex's skill namespace is global, so both channels cannot resolve
// at once; see codexInstallArgvSequence).
func otherChannelMarketplace(devBranch string) string {
	if devBranch == "main" {
		return "spacedock-edge"
	}
	return "spacedock"
}

// channelMarketplaceSource is the marketplace add source the channel installs from:
// a stable binary (devBranch=main) installs from the bare repo `spacedock-dev/marketplace`,
// whose root marketplace.json is named `spacedock`; an edge binary (any other devBranch)
// installs from `spacedock-dev/marketplace@edge`, the `edge` branch whose root
// marketplace.json is named `spacedock-edge`. The `@edge` ref is what makes the host
// register a marketplace NAMED `spacedock-edge`, so the channel id `spacedock@spacedock-edge`
// then resolves — a bare-source add registers `spacedock` and the edge id lookup fails.
// An explicit SPACEDOCK_MARKETPLACE_SOURCE override (marketplaceSource != the production
// default) replaces the source verbatim on every channel: the operator's chosen
// local/alternate marketplace is already the complete marketplace they want, so no
// channel ref is appended.
func channelMarketplaceSource(devBranch string) string {
	if marketplaceSource != defaultMarketplaceSource {
		return marketplaceSource
	}
	if devBranch == "main" {
		return marketplaceSource
	}
	return marketplaceSource + "@edge"
}

// Launch spawns argv as a child and stays resident as its parent: the launcher
// inherits the real terminal (interactive `claude --agent …`), forwards
// externally-targeted signals while letting terminal signals reach the host
// through the shared foreground process group, waits, and returns the host's
// propagated exit code. The host is NOT placed in its own process group (no
// Setpgid), so it owns the controlling TTY and the kernel delivers
// terminal-generated signals to it directly. The returned int is the host's exit
// code (signal-death rendered as 128+signum on unix); a non-nil error is a
// *launch* failure (host binary not found, fork failure), never a non-zero host
// exit.
func (execHost) Launch(argv []string, env []string) (int, error) {
	bin, err := exec.LookPath(argv[0])
	if err != nil {
		return 1, err
	}
	cmd := exec.Command(bin, argv[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	// Arm signal handling BEFORE Start so no terminal signal races the launcher's
	// default disposition (a default SIGINT would kill the launcher before the
	// host is even up); a signal arriving during fork/exec queues on the buffered
	// channel. Begin forwarding only AFTER Start returns, so the go-statement that
	// spawns the pump carries a happens-before edge from Start's write of
	// cmd.Process to the goroutine's read of it (forwarding before that edge would
	// race Start's write). Unix-only; both are no-ops on other platforms.
	forward, stop := forwardHostSignals(cmd)
	defer stop()
	if err := cmd.Start(); err != nil {
		return 1, err
	}
	forward()
	waitErr := cmd.Wait()
	return hostExitCode(cmd.ProcessState, waitErr), nil
}

// installStep is one entry in the install upgrade sequence: an argv to pass to
// the host CLI and a per-step tolerateExit flag. When tolerateExit is true,
// Install treats a non-zero exit as recoverable (appends output, continues);
// otherwise the non-zero exit aborts Install with a wrapped error. The flag is
// co-located with the argv so the tolerance decision is visible to tests via
// installArgvSequence rather than hidden in Install's control flow.
type installStep struct {
	argv         []string
	tolerateExit bool
}

// retiredRouteAID is the retired route-A plugin id from round 1's collapse (the
// entry-name shape `spacedock-edge@spacedock`, superseded by the channel-in-
// marketplace-name shape `spacedock@spacedock-edge` / `spacedock@spacedock`). Its
// migration step below is unconditional (not gated by devBranch) because a
// route-A holder needs migrating on any subsequent claude install/heal run,
// stable or edge — round 1's AC-4.
const retiredRouteAID = "spacedock-edge@spacedock"

// installArgvSequence gives the 6 commands that `Install` runs for claude. No step
// uses `plugin marketplace remove`. On stable, that command names the `spacedock`
// marketplace, which also holds other plugins (subspace, cargento). Probe 1
// measured that this command uninstalls every plugin from the removed marketplace.
// This sequence exists to prevent that damage.
//
// The commands are:
//
//  1. Uninstall the plugin of the selected channel (tolerated). Claude tracks an
//     installed plugin through its marketplace record. If the record goes first,
//     the uninstall has no record to act on.
//  2. Uninstall the plugin of the sibling channel (tolerated). This keeps one
//     channel on the host, as codexInstallArgvSequence does. Both channels use the
//     entry name `spacedock` and the same skill names. If both stay installed, the
//     plugin that serves is unpredictable. `plugin uninstall` acts on one plugin,
//     so it cannot cascade to a plugin on the sibling marketplace. This sequence
//     never touches the marketplace record of the sibling.
//  3. Uninstall the retired route-A plugin (tolerated). This is the round-1
//     migration.
//  4. Add the marketplace source of the channel (fail-fast). Claude re-pins a
//     changed source at the same name (probe 3). This step alone makes the re-pin
//     that the old remove step forced.
//  5. Update the marketplace snapshot of the channel (tolerated). Probe 4 shows
//     that this refresh is not destructive. It covers the case where step 4 made
//     no change because the source was already current.
//  6. Install the channel id that devBranch selects (fail-fast). Probe 9 shows
//     that `plugin install` always clones the content again.
//
// The id is `spacedock@spacedock` for stable and `spacedock@spacedock-edge` for
// edge. The entry name is always `spacedock`, and the channel is the marketplace
// name. The marketplace manifest holds the version pin.
//
// The tolerance is asymmetric. The four cleanup commands are best-effort. Claude
// exits 1 on a fresh box and on an already-current source. No stable stderr shape
// separates these cases from a true failure. An absent sibling alone gives two
// different non-zero shapes, which depend on the registration of its marketplace.
// The two pinning commands stay fail-fast, and they report a broken install from
// the network, a contract incompatibility, or a missing source.
func installArgvSequence(source, devBranch string) []installStep {
	id := channelPluginID(devBranch)
	return []installStep{
		{argv: []string{"plugin", "uninstall", id}, tolerateExit: true},
		{argv: []string{"plugin", "uninstall", "spacedock@" + otherChannelMarketplace(devBranch)}, tolerateExit: true},
		{argv: []string{"plugin", "uninstall", retiredRouteAID}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "add", source}},
		{argv: []string{"plugin", "marketplace", "update", channelMarketplace(devBranch)}, tolerateExit: true},
		{argv: []string{"plugin", "install", id}},
	}
}

// codexInstallArgvSequence is the codex analog of installArgvSequence: no step
// ever spells `plugin marketplace remove` either (codex cascades the same
// dependent-uninstall harm as claude on that command, probe 1 — measured against
// the captain's real codex config, which has `subspace@spacedock` installed
// alongside Spacedock). Codex's global skill namespace still means every
// Spacedock provider must be exclusive, so exclusivity is carried entirely by
// `plugin remove` (content-level, not marketplace-level): remove the sibling
// channel's plugin (tolerated — channel exclusivity), the selected channel's own
// plugin (tolerated — cache hygiene: without this, stale version dirs accumulate
// under the cache and latestVersionDir's fallback can misorder them), AND the
// dedicated `--plugin-dir` dev id (tolerated — the round-trip case: a captain who
// used `--plugin-dir` and now runs a normal channel install must not end up with
// both `spacedock@spacedock-local` and the channel id installed at once, or
// $spacedock:* resolves ambiguously; see codexPluginDirInstallArgvSequence for
// the reverse direction). The marketplace add carries the source
// channelMarketplaceSource resolved (no `--ref` flag — any channel ref is part of
// the source string): the bare repo for stable (root marketplace.json named
// `spacedock`), `<repo>@edge` for edge (the `edge` branch whose root
// marketplace.json is named `spacedock-edge`, so the add registers
// `spacedock-edge` and the channel id `spacedock@spacedock-edge` resolves) —
// fail-fast, the real-failure backstop. `plugin marketplace upgrade` re-pins the
// snapshot non-destructively (probe 4; tolerated — a local-path source errors
// harmlessly at exit 0, still tolerated for forward safety). The entry is always
// `spacedock`; the version pin lives in the marketplace manifest. `plugin add`
// unconditionally re-clones content (probe 8) and stays fail-fast.
func codexInstallArgvSequence(source, devBranch string) []installStep {
	id := channelPluginID(devBranch)
	return []installStep{
		{argv: []string{"plugin", "remove", "spacedock@" + otherChannelMarketplace(devBranch)}, tolerateExit: true},
		{argv: []string{"plugin", "remove", id}, tolerateExit: true},
		{argv: []string{"plugin", "remove", codexLocalPluginID}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "add", source}},
		{argv: []string{"plugin", "marketplace", "upgrade", channelMarketplace(devBranch)}, tolerateExit: true},
		{argv: []string{"plugin", "add", id}},
	}
}

// codexPluginDirInstallArgvSequence is the `--plugin-dir` dev-install analog of
// codexInstallArgvSequence, targeting the dedicated codexLocalMarketplaceName
// instead of a real channel. Because the local marketplace's name never
// collides with either real channel's name (unlike the pre-fix shape, which
// borrowed channelMarketplace(devBranch)), no step here — or in
// codexInstallArgvSequence above — ever needs a conditional `marketplace
// remove`: a real channel install and a local dev install can never contend for
// the same marketplace registration, so the "no install step ever removes a
// marketplace" rule holds with NO exception. Plugin-level exclusivity still
// applies (codex's skill namespace is global): remove BOTH real channel ids
// (tolerated) plus any prior local install (tolerated — cache hygiene, same
// reasoning as codexInstallArgvSequence's own-channel remove) before
// add/upgrade/install under the local name.
func codexPluginDirInstallArgvSequence(source string) []installStep {
	return []installStep{
		{argv: []string{"plugin", "remove", "spacedock@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "remove", "spacedock@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "remove", codexLocalPluginID}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "add", source}},
		{argv: []string{"plugin", "marketplace", "upgrade", codexLocalMarketplaceName}, tolerateExit: true},
		{argv: []string{"plugin", "add", codexLocalPluginID}},
	}
}

// Install shells the host plugin upgrade sequence (cleanup/refresh-then-pin) for
// claude or codex, returning combined output. Each host uses its own verb
// vocabulary (claude `uninstall`/`install`; codex `remove`/`add`), supplied by
// installArgvSequence / codexInstallArgvSequence; devBranch selects the channel
// entry both install. Neither sequence ever spells `plugin marketplace remove` —
// the cleanup/refresh steps are tolerated (their non-zero exits on a fresh-box or
// already-current case are appended to combined output and the loop continues);
// the pinning steps (marketplace add, plugin install/add) are fail-fast and
// surface real install failures.
func (execHost) Install(host, source, devBranch string) (string, error) {
	var steps []installStep
	switch host {
	case "claude":
		steps = installArgvSequence(source, devBranch)
	case "codex":
		steps = codexInstallArgvSequence(source, devBranch)
	default:
		return "", fmt.Errorf("programmatic install is supported for claude and codex, not %q", host)
	}
	return runInstallSteps(host, steps)
}

// InstallCodexLocalPluginDir shells the codex `--plugin-dir` dev-install
// sequence (codexPluginDirInstallArgvSequence) against source — the persistent
// local marketplace WriteCodexLocalMarketplace builds — returning combined
// output.
func (execHost) InstallCodexLocalPluginDir(source string) (string, error) {
	return runInstallSteps("codex", codexPluginDirInstallArgvSequence(source))
}

// runInstallSteps shells argv steps against host in order, tolerating a
// tolerateExit step's non-zero exit (appended to combined output, loop
// continues) and returning a wrapped error on a fail-fast step's non-zero exit.
// Shared by Install and InstallCodexLocalPluginDir so the loop is written once.
func runInstallSteps(host string, steps []installStep) (string, error) {
	var sb strings.Builder
	for _, step := range steps {
		cmd := exec.Command(host, step.argv...)
		out, err := cmd.CombinedOutput()
		sb.Write(out)
		if err != nil {
			if step.tolerateExit {
				continue
			}
			return sb.String(), fmt.Errorf("%s %s: %w", host, strings.Join(step.argv, " "), err)
		}
	}
	return sb.String(), nil
}
