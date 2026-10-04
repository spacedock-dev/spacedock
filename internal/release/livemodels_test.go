// ABOUTME: Guards that the live workflow carries no lane model literal and that
// ABOUTME: every model it runs resolves through `spacedock live-models`.
package release

import (
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestRuntimeLiveWorkflowCarriesNoLaneModelLiteral pins AC-3: the workflow must
// resolve every lane model from `spacedock live-models`, so no pinned id may
// appear as a literal. The ids are authored here, independent of LiveModels(),
// so this fails if a literal is re-added (including the old gpt-5.6-luna) or a
// new model is inlined instead of added to the single source.
func TestRuntimeLiveWorkflowCarriesNoLaneModelLiteral(t *testing.T) {
	workflow := readWorkflow(t, "runtime-live-e2e.yml")
	for _, id := range []string{"claude-sonnet-5", "claude-opus-4-8", "gpt-6-luna", "gpt-5.6-luna"} {
		if strings.Contains(workflow, id) {
			t.Errorf("runtime-live-e2e.yml carries the model literal %q; lanes must resolve models via `spacedock live-models`", id)
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

var liveModelGetPattern = regexp.MustCompile(`live-models --get ([^ )"']+)`)

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
