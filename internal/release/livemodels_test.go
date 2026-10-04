// ABOUTME: Guards that the live workflow resolves every lane model through
// ABOUTME: `spacedock live-models`, and checks the pinned Pi ids against the
// ABOUTME: installed provider catalog; lanes with no independent oracle are named.
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

// TestRuntimeLiveWorkflowModelSitesResolveThroughPrinter pins AC-3 without
// enumerating model ids. It fails if a model value appears as a literal at any
// workflow site that carries one — a `--model` argument, a matrix `model:` key,
// or a step-summary `Model:` label. Because the rule keys off the site's shape
// and not a fixed id list, a model id that did not exist when this test was
// written still fails if it is inlined at one of these sites.
//
// Deferred risk: this rule only sees the three site shapes the workflow uses
// today. Its exact trigger is a model literal inlined at a different site shape
// (for example a new `SPACEDOCK_*_MODEL:` env assignment or a `model=` flag on
// a different command); that literal would pass. Extend the site list when such
// a site is added rather than adding ids to a forbidden list.
func TestRuntimeLiveWorkflowModelSitesResolveThroughPrinter(t *testing.T) {
	workflow := readWorkflow(t, "runtime-live-e2e.yml")

	flagSites := 0
	for _, m := range modelFlagArgPattern.FindAllStringSubmatch(workflow, -1) {
		flagSites++
		if arg := m[1]; !strings.HasPrefix(arg, `"$`) {
			t.Errorf("`--model %s` is a literal; lanes must resolve models via `spacedock live-models`", arg)
		}
	}
	if flagSites == 0 {
		t.Fatal("runtime-live-e2e.yml has no `--model` site; the Codex shim's model forwarding was removed")
	}

	keySites := 0
	for _, line := range strings.Split(workflow, "\n") {
		m := modelKeyLinePattern.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		keySites++
		if !strings.Contains(m[1], "${{") {
			t.Errorf("workflow line %q sets a lane model to a literal; it must resolve through `spacedock live-models`", strings.TrimSpace(line))
		}
	}
	if keySites == 0 {
		t.Fatal("runtime-live-e2e.yml has no `model:`/`Model:` site; the lane model wiring was removed")
	}
}

var (
	modelFlagArgPattern = regexp.MustCompile(`--model\s+(\S+)`)
	modelKeyLinePattern = regexp.MustCompile(`(?:^|\s)[Mm]odel:\s*(.+)$`)
	resolverLinePattern = regexp.MustCompile(`echo "([A-Za-z0-9_]+)=\$\(.*live-models --get ([^)]+)\)"`)
	liveModelGetPattern = regexp.MustCompile(`live-models --get ([^ )"']+)`)
)

// TestRuntimeLiveWorkflowResolverBindsEachOutputToItsKey proves the offline
// resolver fills each published output from its OWN `--get` key. It fails if
// the resolver's keys are exchanged — e.g. `claude_sonnet` filled from
// `--get claude.opus` — which leaves every output name and every `--get` key
// present, so the presence-only wiring guard below cannot see it. The expected
// output→key map is authored here, independent of the workflow.
func TestRuntimeLiveWorkflowResolverBindsEachOutputToItsKey(t *testing.T) {
	workflow := readWorkflow(t, "runtime-live-e2e.yml")
	step, ok := stepNamed(parseWorkflowSteps(workflow), "Resolve live lane models")
	if !ok {
		t.Fatal("workflow lacks the `Resolve live lane models` step")
	}
	got := map[string]string{}
	for _, m := range resolverLinePattern.FindAllStringSubmatch(step.run, -1) {
		got[m[1]] = m[2]
	}
	want := map[string]string{
		"claude_sonnet": "claude.sonnet",
		"claude_opus":   "claude.opus",
		"codex_exec":    "codex.exec",
		"pi_oauth":      "pi.oauth",
		"pi_api_key":    "pi.api-key",
	}
	if len(got) != len(want) {
		t.Fatalf("resolver emits %d lane outputs, want %d: %v", len(got), len(want), got)
	}
	for output, key := range want {
		if got[output] != key {
			t.Errorf("resolver output %q resolves key %q, want %q", output, got[output], key)
		}
	}
}

// TestRuntimeLiveWorkflowResolvedKeysExist proves every `spacedock live-models
// --get <key>` the workflow issues names a real lane. It fails if a key is
// renamed in the workflow without a matching LiveModels() entry, which would
// make the workflow resolve an empty model at run time.
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

// TestRuntimeLiveWorkflowLaneModelWiring proves the offline job publishes each
// lane output and a live job consumes it, so a model resolved by the print
// command actually reaches the matrix, the Codex shim, and the Pi summary. It
// fails if an output is renamed on one side only (the consumer would read an
// empty string) or if the matrix/shim/summary are reverted to literals.
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
	want := map[string]string{
		"claude_sonnet": "claude.sonnet",
		"claude_opus":   "claude.opus",
		"codex_exec":    "codex.exec",
		"pi_oauth":      "pi.oauth",
		"pi_api_key":    "pi.api-key",
	}
	for output, key := range want {
		if got := offline.Outputs[output]; got != "${{ steps.live_models.outputs."+output+" }}" {
			t.Errorf("offline outputs[%q] = %q, want the live_models step output", output, got)
		}
		if !strings.Contains(workflow, "--get "+key) {
			t.Errorf("workflow resolves no lane model for key %q", key)
		}
		// A trailing word character means the output name is a prefix of a
		// different name (a renamed consumer), which would read an empty model.
		consumer := regexp.MustCompile(`needs\.offline\.outputs\.` + regexp.QuoteMeta(output) + `(?:[^A-Za-z0-9_]|$)`)
		if !consumer.MatchString(workflow) {
			t.Errorf("no live job consumes needs.offline.outputs.%s", output)
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
//   - pi.oauth, pi.api-key: this test (installed catalog).
//   - codex.exec: none in reach. No installed registry declares the Codex
//     `exec --model` ids, so its value rests on the authored exact-output
//     oracle in internal/cli/live_models_test.go and on the live smoke.
//   - claude.sonnet, claude.opus: none in reach. The installed Claude CLI
//     validates a model only after auth; with an isolated home it short-circuits
//     with "Not logged in" before the model is checked, and with real
//     credentials it would spend an API call and depend on the network. So the
//     Claude ids rest on the authored exact-output oracle and the recorded
//     rejection evidence, not on a CLI probe. This is stated rather than claimed.
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
