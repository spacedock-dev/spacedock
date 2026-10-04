// ABOUTME: Guards that the live workflow resolves every lane model through
// ABOUTME: `spacedock live-models`, and checks the Codex and Pi pins against
// ABOUTME: installed provider metadata; lanes without an independent oracle are named.
package release

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestRuntimeLiveWorkflowModelSitesResolveThroughPrinter is AC-3's structural
// guard. At every model-bearing site the workflow uses today the model must be a
// reference, never an inline value, so a literal model id cannot satisfy the
// required shape:
//
//   - `--model` argument: must be exactly the `"$live_model"` shell variable,
//     which the Codex shim fills from `SPACEDOCK_LIVE_CODEX_MODEL`.
//   - `live_model=` assignment: must read `${SPACEDOCK_LIVE_CODEX_MODEL...}`.
//   - matrix `model:` / summary `Model:` lines: must contain a `${{ ... }}`
//     reference and no model-id token before it.
//   - `SPACEDOCK_LIVE_CODEX_MODEL=` / `SPACEDOCK_LIVE_MODEL:` env lines: must
//     reference the offline resolver output or `matrix.model`.
//   - resolver `echo` lines: the exported value must come from
//     `live-models --get`, so a duplicated literal write is caught.
//
// This rule does NOT cover, and cannot without re-implementing GitHub's
// expression language or the shell:
//   - a model literal nested inside a `${{ ... }}` operand (e.g.
//     `... || 'gpt-7-luna'`) still contains a `${{` reference;
//   - a model literal appended to a derived string that already interpolates
//     `${{ needs.offline.outputs.* }}` (e.g. the journey-delta artifact name);
//   - a model-bearing site shape the workflow does not use today (a
//     `--model=<id>` `=` form, a workflow input, another host's flag).
//
// Extend the site list when such a site appears; do not add a forbidden-id list.
func TestRuntimeLiveWorkflowModelSitesResolveThroughPrinter(t *testing.T) {
	workflow := readWorkflow(t, "runtime-live-e2e.yml")

	flagSites := 0
	for _, m := range modelFlagArgPattern.FindAllStringSubmatch(workflow, -1) {
		flagSites++
		if arg := m[1]; arg != `"$live_model"` {
			t.Errorf("`--model %s` is not the resolved shell variable; lanes must resolve models via `spacedock live-models`", arg)
		}
	}
	if flagSites == 0 {
		t.Fatal("runtime-live-e2e.yml has no `--model` site; the Codex shim's model forwarding was removed")
	}

	assignSites := 0
	for _, line := range strings.Split(workflow, "\n") {
		m := codexLiveModelAssignPattern.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		assignSites++
		if !strings.Contains(m[1], "${SPACEDOCK_LIVE_CODEX_MODEL") {
			t.Errorf("workflow line %q sets live_model to a value; it must read SPACEDOCK_LIVE_CODEX_MODEL", strings.TrimSpace(line))
		}
	}
	if assignSites == 0 {
		t.Fatal("runtime-live-e2e.yml has no `live_model=` assignment; the Codex shim no longer reads SPACEDOCK_LIVE_CODEX_MODEL")
	}

	keySites := 0
	for _, line := range strings.Split(workflow, "\n") {
		m := modelKeyLinePattern.FindStringSubmatch(stripLineComment(line))
		if m == nil {
			continue
		}
		keySites++
		value := m[1]
		ref := strings.Index(value, "${{")
		if ref < 0 {
			t.Errorf("workflow line %q sets a lane model without a `${{ ... }}` reference", strings.TrimSpace(line))
			continue
		}
		if lhs := value[:ref]; modelIDTokenPattern.MatchString(lhs) {
			t.Errorf("workflow line %q carries a model literal before its reference", strings.TrimSpace(line))
		}
	}
	if keySites == 0 {
		t.Fatal("runtime-live-e2e.yml has no `model:`/`Model:` site; the lane model wiring was removed")
	}

	envSites := 0
	for _, line := range strings.Split(workflow, "\n") {
		switch {
		case strings.Contains(line, "SPACEDOCK_LIVE_CODEX_MODEL="):
			envSites++
			if !strings.Contains(line, "${{ needs.offline.outputs.") {
				t.Errorf("workflow line %q sets SPACEDOCK_LIVE_CODEX_MODEL without the offline resolver output", strings.TrimSpace(line))
			}
		case strings.Contains(line, "SPACEDOCK_LIVE_MODEL:"):
			envSites++
			if !strings.Contains(line, "${{ matrix.model }}") {
				t.Errorf("workflow line %q sets SPACEDOCK_LIVE_MODEL without matrix.model", strings.TrimSpace(line))
			}
		}
	}
	if envSites == 0 {
		t.Fatal("runtime-live-e2e.yml sets no SPACEDOCK_LIVE_* model env; the model wiring was removed")
	}

	step, ok := stepNamed(parseWorkflowSteps(workflow), "Resolve live lane models")
	if !ok {
		t.Fatal("workflow lacks the `Resolve live lane models` step")
	}
	resolverSites := 0
	for _, line := range strings.Split(step.run, "\n") {
		m := resolverEchoPattern.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		resolverSites++
		if !strings.Contains(m[2], "live-models --get") {
			t.Errorf("resolver line %q exports a literal; it must resolve through `spacedock live-models --get`", strings.TrimSpace(line))
		}
	}
	if resolverSites == 0 {
		t.Fatal("the `Resolve live lane models` step exports no lane model; the resolver was removed")
	}
}

var (
	modelFlagArgPattern         = regexp.MustCompile(`--model\s+(\S+)`)
	modelKeyLinePattern         = regexp.MustCompile(`(?:^|\s)[Mm]odel:\s*(.+)$`)
	modelIDTokenPattern         = regexp.MustCompile(`(?:gpt-|claude-|openai/|openai-codex/)`)
	codexLiveModelAssignPattern = regexp.MustCompile(`(?:^|\s)live_model=(.*)$`)
	resolverEchoPattern         = regexp.MustCompile(`echo "([A-Za-z0-9_]+)=(.*)"`)
	liveModelGetPattern         = regexp.MustCompile(`live-models --get ([^ )"']+)`)
)

// stripLineComment drops a trailing ` # ...` comment so a model literal cannot
// sit in a comment while a `${{ ... }}` reference satisfies the shape check.
func stripLineComment(line string) string {
	if i := strings.Index(line, " #"); i >= 0 {
		return line[:i]
	}
	return line
}

// TestRuntimeLiveWorkflowLaneModelWiring proves each lane output the offline job
// declares is filled from the `live_models` step and consumed by a later job, so
// a model the print command resolves actually reaches a lane. It derives the
// output names from the workflow itself (no copied output→key map) and fails if
// an output is wired to anything but its own resolver output, or if no job reads
// `needs.offline.outputs.<name>` (a renamed consumer would read an empty model).
func TestRuntimeLiveWorkflowLaneModelWiring(t *testing.T) {
	workflow := readWorkflow(t, "runtime-live-e2e.yml")
	var parsed struct {
		Jobs map[string]struct {
			Outputs map[string]string `yaml:"outputs"`
		} `yaml:"jobs"`
	}
	if err := yaml.Unmarshal([]byte(workflow), &parsed); err != nil {
		t.Fatal(err)
	}
	offline, ok := parsed.Jobs["offline"]
	if !ok {
		t.Fatal("workflow lacks the offline job that resolves lane models")
	}
	if len(offline.Outputs) == 0 {
		t.Fatal("offline job declares no lane model outputs")
	}
	for output, value := range offline.Outputs {
		if want := "${{ steps.live_models.outputs." + output + " }}"; value != want {
			t.Errorf("offline outputs[%q] = %q, want %q", output, value, want)
		}
		consumer := regexp.MustCompile(`needs\.offline\.outputs\.` + regexp.QuoteMeta(output) + `(?:[^A-Za-z0-9_]|$)`)
		if !consumer.MatchString(workflow) {
			t.Errorf("no live job consumes needs.offline.outputs.%s", output)
		}
	}
}

// TestRuntimeLiveWorkflowResolvedKeysExist proves every `spacedock live-models
// --get <key>` the workflow issues names a real lane in the Go source. It fails
// if a key is renamed in the workflow without a matching LiveModels() entry,
// which would make the workflow resolve an empty model at run time.
func TestRuntimeLiveWorkflowResolvedKeysExist(t *testing.T) {
	workflow := readWorkflow(t, "runtime-live-e2e.yml")
	matches := liveModelGetPattern.FindAllStringSubmatch(workflow, -1)
	if len(matches) == 0 {
		t.Fatal("runtime-live-e2e.yml never resolves a lane model via `spacedock live-models --get`")
	}
	for _, m := range matches {
		if _, ok := LiveModel(m[1]); !ok {
			t.Errorf("workflow resolves live-models key %q, which LiveModels() does not define", m[1])
		}
	}
}

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
//   - claude.sonnet, claude.opus: no implemented provider-backed oracle. The
//     installed Claude CLI validates a model only after auth; with an isolated
//     home it short-circuits with "Not logged in" before the model is checked,
//     and with real credentials it would spend an API call and depend on the
//     network. So the Claude ids rest on the authored exact-output oracle and
//     the recorded rejection evidence, not on a CLI probe. This is stated
//     rather than claimed.
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
