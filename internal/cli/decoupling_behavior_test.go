// ABOUTME: AC-1 decoupling behavior — a real isolated-CLAUDE_CONFIG_DIR install of a
// ABOUTME: tag-pinned stable + branch-HEAD edge channel, advancing HEAD to prove the freeze.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// TestStableChannelDecoupledFromBranchHead is the load-bearing AC-1 proof: a
// tag-pinned stable channel stays frozen on the tag's tree while an edge channel
// on a branch HEAD advances, both resolved from ONE marketplace as two entries of
// one source. It reproduces the recorded decoupling spike against the real claude
// host with LOCAL FIXTURE git repos (the spike used /tmp throwaway repos):
//
//  1. A plugin git repo with tag v0.0.1 (version 0.0.1) on an early commit and a
//     `next` HEAD ahead at 0.0.2.
//  2. A marketplace dir holding ONE marketplace.json with two entries of that one
//     {source:url,url:file://<plugin>,ref:…} source — `spacedock` (stable) ref
//     v0.0.1, `spacedock-edge` (edge) ref next.
//  3. Install both channels into an isolated CLAUDE_CONFIG_DIR + plugin cache.
//     Observe two distinct cache dirs: stable 0.0.1, edge 0.0.2 — byte-verified
//     from the installed skill body, not a command's self-claim.
//  4. Advance plugin `next` HEAD to 0.0.3, bump ONLY the edge entry's version,
//     refresh the marketplace, and update both channels.
//  5. Assert stable is FROZEN (its cache still holds only 0.0.1, body v0.0.1)
//     while edge ADVANCED (a new 0.0.3 cache dir, body v0.0.3).
//
// Skips when `claude` is not on PATH; a real install kept hermetic by env
// isolation + local-path fixtures (file:// git urls → offline), not a mock.
func TestStableChannelDecoupledFromBranchHead(t *testing.T) {
	claudeBin, err := exec.LookPath("claude")
	if err != nil {
		t.Skip("claude not on PATH; decoupling behavior test requires the host CLI")
	}

	tmp := t.TempDir()
	plugin := buildPluginGitRepo(t, filepath.Join(tmp, "plugin"))
	marketplace := buildChannelMarketplace(t, filepath.Join(tmp, "marketplace"), plugin, "v0.0.1", "next")
	configDir := filepath.Join(tmp, "config")
	cacheDir := filepath.Join(tmp, "cache")
	mustMkdir(t, configDir)
	mustMkdir(t, cacheDir)

	env := append(os.Environ(),
		"CLAUDE_CONFIG_DIR="+configDir,
		"CLAUDE_CODE_PLUGIN_CACHE_DIR="+cacheDir,
	)

	// Install both channels from the one marketplace.
	runHost(t, claudeBin, env, "plugin", "marketplace", "add", marketplace)
	runHost(t, claudeBin, env, "plugin", "install", "spacedock@spacedock")
	runHost(t, claudeBin, env, "plugin", "install", "spacedock-edge@spacedock")

	// Two distinct cache dirs resolve from the one marketplace: stable on the tag
	// commit (v0.0.1), edge on next HEAD (0.0.2). Byte-verified from the installed
	// skill body — the tag pin checks out the tag's tree, not HEAD's.
	if body := installedSkillBody(t, cacheDir, "spacedock", "0.0.1"); body != "body v0.0.1\n" {
		t.Fatalf("stable cache skill body = %q, want %q (tag pin must serve the tag's tree)", body, "body v0.0.1\n")
	}
	if body := installedSkillBody(t, cacheDir, "spacedock-edge", "0.0.2"); body != "body v0.0.2\n" {
		t.Fatalf("edge cache skill body = %q, want %q (edge serves next HEAD)", body, "body v0.0.2\n")
	}

	// New work lands post-release: advance plugin `next` HEAD to 0.0.3 and bump
	// ONLY the edge entry's version, then refresh the marketplace.
	advancePluginHead(t, plugin, "0.0.3")
	bumpEntryVersion(t, marketplace, "spacedock-edge", "0.0.3")
	runHost(t, claudeBin, env, "plugin", "marketplace", "update", "spacedock")
	// Updates are channel-scoped; stable is a no-op (already at the latest 0.0.1),
	// edge advances. The stable update may report "already at the latest" (exit 0)
	// — tolerate it and read the on-disk cache as the source of truth either way.
	runHostTolerant(t, claudeBin, env, "plugin", "update", "spacedock@spacedock")
	runHost(t, claudeBin, env, "plugin", "update", "spacedock-edge@spacedock")

	// The decoupling holds under advance: stable is FROZEN — its cache still holds
	// only 0.0.1 (no 0.0.3 dir appeared) and the body is unchanged at v0.0.1,
	// despite plugin HEAD moving to 0.0.3.
	if dirs := cacheVersionDirs(t, cacheDir, "spacedock"); len(dirs) != 1 || dirs[0] != "0.0.1" {
		t.Fatalf("stable cache version dirs = %v, want [0.0.1] only (tag-pinned stable must not advance with HEAD)", dirs)
	}
	if body := installedSkillBody(t, cacheDir, "spacedock", "0.0.1"); body != "body v0.0.1\n" {
		t.Fatalf("after advance, stable cache body = %q, want %q (frozen)", body, "body v0.0.1\n")
	}

	// Edge ADVANCED: a new 0.0.3 cache dir, body v0.0.3 (from next HEAD).
	if body := installedSkillBody(t, cacheDir, "spacedock-edge", "0.0.3"); body != "body v0.0.3\n" {
		t.Fatalf("after advance, edge 0.0.3 cache body = %q, want %q (edge must track HEAD)", body, "body v0.0.3\n")
	}
}

// buildPluginGitRepo writes a plugin git repo under root with a `next` branch
// carrying tag v0.0.1 (version 0.0.1, skill body "body v0.0.1") on an early commit
// and a `next` HEAD advanced to 0.0.2 ("body v0.0.2"). Returns the repo path. The
// plugin manifest carries the display version the doctor verdict reads.
func buildPluginGitRepo(t *testing.T, root string) string {
	t.Helper()
	mustMkdir(t, filepath.Join(root, ".claude-plugin"))
	mustMkdir(t, filepath.Join(root, "skills", "demo"))

	writePluginVersion(t, root, "0.0.1")
	testgit.InitRepo(t, root, "-q", "-b", "next")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "v0.0.1")
	git(t, root, "tag", "v0.0.1")

	writePluginVersion(t, root, "0.0.2")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "v0.0.2 HEAD")
	return root
}

// writePluginVersion (re)writes the plugin's .claude-plugin/plugin.json and demo
// skill body to the given version (the skill body is "body v<version>", the
// decoupling proof's per-version observable).
func writePluginVersion(t *testing.T, root, version string) {
	t.Helper()
	mustWrite(t, filepath.Join(root, ".claude-plugin", "plugin.json"),
		fmt.Sprintf(`{ "name": "spacedock", "version": "%s", "skills": "./skills/" }`+"\n", version))
	mustWrite(t, filepath.Join(root, "skills", "demo", "SKILL.md"),
		fmt.Sprintf("---\nname: demo\ndescription: demo\n---\nbody v%s\n", version))
}

// advancePluginHead advances the plugin repo's `next` HEAD to a new version,
// committing the bump (new work landing post-release).
func advancePluginHead(t *testing.T, root, version string) {
	t.Helper()
	writePluginVersion(t, root, version)
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "v"+version+" HEAD")
}

// buildChannelMarketplace writes a marketplace dir under root holding ONE
// marketplace.json with two entries of one {source:url,url:file://<plugin>,ref:…}
// source: `spacedock` (stable) pinned to stableRef, `spacedock-edge` (edge) on
// edgeRef. Returns the marketplace dir (added to claude by path; claude reads
// .claude-plugin/marketplace.json directly and clones the plugin source url@ref).
func buildChannelMarketplace(t *testing.T, root, pluginRepo, stableRef, edgeRef string) string {
	t.Helper()
	mustMkdir(t, filepath.Join(root, ".claude-plugin"))
	url := "file://" + pluginRepo
	manifest := fmt.Sprintf(`{
  "name": "spacedock",
  "owner": { "name": "test" },
  "plugins": [
    { "name": "spacedock", "source": { "source": "url", "url": "%s", "ref": "%s" }, "description": "stable", "version": "0.0.1", "category": "workflow" },
    { "name": "spacedock-edge", "source": { "source": "url", "url": "%s", "ref": "%s" }, "description": "edge", "version": "0.0.2", "category": "workflow" }
  ]
}
`, url, stableRef, url, edgeRef)
	mustWrite(t, filepath.Join(root, ".claude-plugin", "marketplace.json"), manifest)
	return root
}

// bumpEntryVersion rewrites the named plugin entry's version in the marketplace
// manifest (leaving the other entry untouched) — the edge-only version bump the
// `claude plugin update` re-pull keys on.
func bumpEntryVersion(t *testing.T, marketplace, entry, version string) {
	t.Helper()
	path := filepath.Join(marketplace, ".claude-plugin", "marketplace.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read marketplace %s: %v", path, err)
	}
	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse marketplace: %v", err)
	}
	plugins, _ := doc["plugins"].([]any)
	for _, p := range plugins {
		pm, _ := p.(map[string]any)
		if pm["name"] == entry {
			pm["version"] = version
		}
	}
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("marshal marketplace: %v", err)
	}
	mustWrite(t, path, string(out)+"\n")
}

// installedSkillBody reads the installed demo skill body for the given entry at
// the given version from the plugin cache:
// <CACHE>/cache/<marketplace>/<entry>/<version>/skills/demo/SKILL.md, returning
// the body line after the YAML frontmatter (the per-version observable).
func installedSkillBody(t *testing.T, cacheDir, entry, version string) string {
	t.Helper()
	path := filepath.Join(cacheDir, "cache", "spacedock", entry, version, "skills", "demo", "SKILL.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read installed skill %s: %v", path, err)
	}
	return skillBody(string(data))
}

// skillBody returns the content after the closing `---` of the YAML frontmatter
// (the SKILL.md shape is `---\n<frontmatter>---\n<body>`), or the whole content
// when there is no frontmatter fence.
func skillBody(content string) string {
	parts := strings.SplitN(content, "---\n", 3)
	if len(parts) < 3 {
		return content
	}
	return parts[2]
}

// cacheVersionDirs lists the version subdirectories the host cached for the given
// entry under <CACHE>/cache/spacedock/<entry>/. A tag-pinned channel must hold
// exactly one (frozen); an advancing channel grows a new dir.
func cacheVersionDirs(t *testing.T, cacheDir, entry string) []string {
	t.Helper()
	root := filepath.Join(cacheDir, "cache", "spacedock", entry)
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("read cache root %s: %v", root, err)
	}
	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}
	return dirs
}
