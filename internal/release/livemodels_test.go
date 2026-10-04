// ABOUTME: Independent provider oracles for the pinned lane model values: the
// ABOUTME: Pi ids are checked against the installed pi-ai catalog and the Codex
// ABOUTME: id against the installed Codex model cache. The placement of those
// ABOUTME: ids in the live workflow has no enforcement here.
package release

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestPiLaneModelsExistInInstalledCatalog is the independent value oracle for
// the Pi lane. It parses each pinned Pi id into provider/model/thinking and
// requires the installed pi-ai provider catalog to declare that model with that
// thinking level. It fails when a pinned Pi id names a provider/model the
// catalog does not have, or a thinking level it does not allow — a wrong value
// that no name-presence guard can catch. It skips when the host has no pi
// install with a catalog, because the catalog is host state, not repository
// content, so this oracle cannot run in the offline CI image.
//
// Independent-value coverage by lane:
//   - pi.oauth, pi.api-key: this test (installed pi-ai catalog).
//   - codex.exec: TestCodexLaneModelExistsInInstalledCache (installed Codex
//     model cache).
//   - claude.sonnet, claude.opus: no independent oracle. The installed Claude
//     CLI validates a model only after auth; with an isolated home it
//     short-circuits with "Not logged in" before the model is checked, and with
//     real credentials it would spend an API call and depend on the network.
//     The recorded rejection establishes only that one candidate id
//     (claude-sonnet-5.5) is refused by the installed CLI; it does not validate
//     the retained ids. The retained claude.sonnet and claude.opus ids rest only
//     on the authored print test, with no independent oracle.
//
// Where the values land has no enforcement today: which workflow lane output
// feeds which cadence lane is unchecked, so a swap of the two Claude output
// operands passes every remaining check. A comparison between the workflow and
// the source is not independent here, because this change authors both sides;
// that is an absence of enforcement today, not proof that none is possible.
func TestPiLaneModelsExistInInstalledCatalog(t *testing.T) {
	catalogDir, ok := installedPiCatalogDir()
	if !ok {
		t.Skip("no installed pi-ai provider catalog on this host; the catalog is host state, not repository content")
	}
	for _, tc := range []struct{ lane, id string }{
		{"pi.oauth", PiOAuthModel},
		{"pi.api-key", PiAPIKeyModel},
	} {
		provider, model, thinking := splitPiModelID(t, tc.lane, tc.id)
		catalogPath := filepath.Join(catalogDir, provider+".json")
		data, err := os.ReadFile(catalogPath)
		if err != nil {
			t.Errorf("%s: read catalog %s: %v", tc.lane, catalogPath, err)
			continue
		}
		declared, err := catalogDeclaresModel(data, model, thinking)
		if err != nil {
			t.Errorf("%s: parse catalog %s: %v", tc.lane, catalogPath, err)
			continue
		}
		if !declared {
			t.Errorf("%s: pinned id %q is not declared by %s with thinking %q", tc.lane, tc.id, provider+".json", thinking)
		}
	}
}

// installedPiCatalogDir finds the pi-ai provider catalog beside the installed
// pi package. It returns ("", false) when pi is absent or the catalog layout
// is not present, so a caller skips rather than fails on a machine without pi.
func installedPiCatalogDir() (string, bool) {
	bin, err := exec.LookPath("pi")
	if err != nil {
		return "", false
	}
	resolved, err := filepath.EvalSymlinks(bin)
	if err != nil {
		return "", false
	}
	for dir := filepath.Dir(resolved); ; {
		catalog := filepath.Join(dir, "node_modules", "@earendil-works", "pi-ai", "dist", "providers", "data")
		if info, err := os.Stat(catalog); err == nil && info.IsDir() {
			return catalog, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

// splitPiModelID splits a pinned Pi id of the form provider/model:thinking.
func splitPiModelID(t *testing.T, lane, id string) (provider, model, thinking string) {
	t.Helper()
	provider, rest, ok := strings.Cut(id, "/")
	if !ok {
		t.Fatalf("%s: pinned id %q is not provider-qualified", lane, id)
	}
	model, thinking, ok = strings.Cut(rest, ":")
	if !ok {
		t.Fatalf("%s: pinned id %q carries no :thinking suffix", lane, id)
	}
	return provider, model, thinking
}

// catalogDeclaresModel reports whether the pi-ai catalog declares model with a
// non-null thinkingLevelMap entry for thinking.
func catalogDeclaresModel(data []byte, model, thinking string) (bool, error) {
	var catalog map[string]map[string]struct {
		ID               string             `json:"id"`
		ThinkingLevelMap map[string]*string `json:"thinkingLevelMap"`
	}
	if err := json.Unmarshal(data, &catalog); err != nil {
		return false, err
	}
	for _, models := range catalog {
		for _, entry := range models {
			if entry.ID != model {
				continue
			}
			level, ok := entry.ThinkingLevelMap[thinking]
			return ok && level != nil, nil
		}
	}
	return false, nil
}

// TestCodexLaneModelExistsInInstalledCache is the independent value oracle for
// the Codex lane. Codex ships an installed model cache
// (${CODEX_HOME:-$HOME/.codex}/models_cache.json) that declares each supported
// slug and its reasoning levels. The test requires CodexExecModel to be a
// declared slug that supports the `max` reasoning level the shim sets, so a typo
// or an unsupported id fails against provider metadata rather than against the
// authored pin. It fails when the pinned id or level is absent; it skips when
// the host has no Codex model cache (host state, not repository content). The
// cache proves the provider supports the spelling and level, not account
// entitlement or a successful `codex exec`.
func TestCodexLaneModelExistsInInstalledCache(t *testing.T) {
	path, ok := codexModelsCachePath()
	if !ok {
		t.Skip("no installed Codex model cache on this host; the cache is host state, not repository content")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read Codex model cache %s: %v", path, err)
	}
	var cache struct {
		Models []struct {
			Slug                     string `json:"slug"`
			SupportedReasoningLevels []struct {
				Effort string `json:"effort"`
			} `json:"supported_reasoning_levels"`
		} `json:"models"`
	}
	if err := json.Unmarshal(data, &cache); err != nil {
		t.Fatalf("parse Codex model cache %s: %v", path, err)
	}
	for _, m := range cache.Models {
		if m.Slug != CodexExecModel {
			continue
		}
		for _, level := range m.SupportedReasoningLevels {
			if level.Effort == "max" {
				return
			}
		}
		t.Errorf("pinned Codex id %q does not declare the max reasoning level the shim sets", CodexExecModel)
		return
	}
	t.Errorf("pinned Codex id %q is not declared by the installed Codex model cache", CodexExecModel)
}

// codexModelsCachePath returns the installed Codex model cache path, honoring
// CODEX_HOME, or ("", false) when no cache file exists.
func codexModelsCachePath() (string, bool) {
	home := os.Getenv("CODEX_HOME")
	if home == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return "", false
		}
		home = filepath.Join(userHome, ".codex")
	}
	path := filepath.Join(home, "models_cache.json")
	if _, err := os.Stat(path); err != nil {
		return "", false
	}
	return path, true
}
