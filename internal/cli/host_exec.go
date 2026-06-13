// ABOUTME: Production hostOps — resolves the installed plugin manifest via
// ABOUTME: `claude/codex plugin list --json`, execs the host, and shells installs.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
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

// resolveClaudeManifest shells `claude plugin list --json`, finds the spacedock@
// spacedock entry, and returns <installPath>/.claude-plugin/plugin.json. Returns
// "" (no error) when the host reports no matching install or no installPath.
func (execHost) resolveClaudeManifest(host string) (string, error) {
	out, err := exec.Command(host, "plugin", "list", "--json").Output()
	if err != nil {
		return "", fmt.Errorf("%s plugin list --json: %w", host, err)
	}
	var entries []pluginListEntry
	if err := json.Unmarshal(out, &entries); err != nil {
		return "", fmt.Errorf("parse %s plugin list --json: %w", host, err)
	}
	for _, e := range entries {
		if e.ID == "spacedock@spacedock" {
			if e.InstallPath == "" {
				return "", nil
			}
			return filepath.Join(e.InstallPath, manifestSubpath(host)), nil
		}
	}
	return "", nil
}

// resolveCodexManifest confirms spacedock@spacedock is installed via the text
// `codex plugin list` and resolves the manifest under the Codex plugin cache.
// Codex's listing carries no install path (its `--json` form has one only for
// the marketplace root, not the cached plugin), so the deterministic cache
// layout is the resolver. Codex installs land at
// <CODEX_HOME>/plugins/cache/<marketplace>/<plugin>/<version>/.codex-plugin/plugin.json.
// Returns "" (no error) when the plugin is not installed or no cached manifest
// exists for it yet.
func (execHost) resolveCodexManifest() (string, error) {
	out, err := exec.Command("codex", "plugin", "list").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("codex plugin list: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	if !codexEntryInstalled(string(out), "spacedock@spacedock") {
		return "", nil
	}
	return codexCacheManifest()
}

// codexCacheManifest resolves the cached spacedock@spacedock manifest under the
// Codex plugin cache: <CODEX_HOME>/plugins/cache/spacedock/spacedock/<version>/
// .codex-plugin/plugin.json. It picks the semver-greatest version dir and
// returns that manifest path, or "" (no error) when no cached manifest exists
// yet (absent cache root, no version dir, or the manifest file is missing).
func codexCacheManifest() (string, error) {
	cacheRoot := filepath.Join(codexHome(), "plugins", "cache", "spacedock", "spacedock")
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

// channelEntry maps the binary's devBranch stamp to the marketplace ENTRY name
// it installs: a stable binary (devBranch=main) installs the `spacedock` entry; an
// edge binary (any other devBranch, e.g. next) installs `spacedock-edge`. The two
// entries live in the one marketplace repo and resolve distinct pinned versions
// (stable pinned to a release tag, edge tracking next HEAD), so the channel is the
// entry name — the version pin lives in the manifest, not an @ref on the install
// command.
func channelEntry(devBranch string) string {
	if devBranch == "main" {
		return "spacedock"
	}
	return "spacedock-edge"
}

// channelPluginID is the `<entry>@spacedock` plugin id for the channel devBranch
// selects: `spacedock@spacedock` (stable) or `spacedock-edge@spacedock` (edge).
// The `@spacedock` suffix is the marketplace NAME (the marketplace.json `name`);
// the prefix is the channel entry.
func channelPluginID(devBranch string) string {
	return channelEntry(devBranch) + "@spacedock"
}

// Launch replaces the current process with argv via execve, so the host CLI
// owns the terminal (interactive `claude --agent …`). It returns only when exec
// itself fails.
func (execHost) Launch(argv []string, env []string) error {
	bin, err := exec.LookPath(argv[0])
	if err != nil {
		return err
	}
	return syscall.Exec(bin, argv, env)
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

// installArgvSequence is the 4-command upgrade shape `Install` issues:
// uninstall any existing channel plugin first (cause-and-effect — claude
// tracks an installed plugin via its marketplace record, so the marketplace
// remove later would orphan a live uninstall), drop the existing marketplace
// declaration for spacedock (so the next add re-pins the source instead of
// no-op'ing on "already on disk"), add the marketplace-repo source, then install
// the channel entry the devBranch selects (`spacedock` stable / `spacedock-edge`
// edge). The marketplace add carries the BARE marketplace-repo source — the
// channel and version pin live in the manifest, not an @ref shorthand. The
// tolerance asymmetry: BOTH cleanup steps (plugin uninstall + marketplace remove)
// are tolerated as best-effort, because claude exits 1 on the fresh-box cases
// ("Plugin not found in installed plugins" for uninstall, "Marketplace 'spacedock'
// not found" for remove) with no way to distinguish those from real failures via
// stable stderr matching. BOTH pinning steps (marketplace add + plugin install)
// stay fail-fast — they are the real-failure backstops that surface a broken
// install (network, contract incompatibility, missing source).
func installArgvSequence(source, devBranch string) []installStep {
	id := channelPluginID(devBranch)
	return []installStep{
		{argv: []string{"plugin", "uninstall", id}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "remove", "spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "add", source}},
		{argv: []string{"plugin", "install", id}},
	}
}

// codexInstallArgvSequence is the codex analog of installArgvSequence: the same
// 4-command cleanup-then-pin shape, but in codex's verb vocabulary (`plugin
// remove` / `plugin add`, not claude's `uninstall` / `install`). The marketplace
// add carries the BARE marketplace-repo source with NO `--ref` — the channel is
// the entry name (`spacedock` stable / `spacedock-edge` edge) the devBranch
// selects, and the version pin lives in the manifest, not a branch ref. The
// tolerance asymmetry matches claude: BOTH cleanup steps (plugin remove +
// marketplace remove) are tolerated — on a fresh box `plugin remove` exits 0
// (idempotent) but `marketplace remove` exits 1 ("marketplace `spacedock` is not
// configured or installed"), and neither is a real failure. BOTH pinning steps
// (marketplace add + plugin add) stay fail-fast — they are the real-failure
// backstops.
func codexInstallArgvSequence(source, devBranch string) []installStep {
	id := channelPluginID(devBranch)
	return []installStep{
		{argv: []string{"plugin", "remove", id}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "remove", "spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "add", source}},
		{argv: []string{"plugin", "add", id}},
	}
}

// Install shells the host plugin upgrade sequence (cleanup-then-pin) for claude
// or codex, returning combined output. Each host uses its own verb vocabulary
// (claude `uninstall`/`install`; codex `remove`/`add`), supplied by
// installArgvSequence / codexInstallArgvSequence; devBranch selects the channel
// entry both install. The two cleanup steps are tolerated — their non-zero exits
// on a fresh-box ("not installed" / "not found") are appended to combined output
// and the loop continues. The two pinning steps (marketplace add, plugin
// install/add) are fail-fast and surface real install failures.
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
