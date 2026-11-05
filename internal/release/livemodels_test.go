package release

import (
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseLiveModelsExcludesCadencePolicy(t *testing.T) {
	got := parseLiveModels("claude.future=fixture-model\nclaude.allowed.sonnet=future\ncodex.exec=fixture-codex\npi.oauth=fixture-pi:max\n")
	want := map[string]string{
		"claude.future": "fixture-model",
		"codex.exec":    "fixture-codex",
		"pi.oauth":      "fixture-pi:max",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseLiveModels() = %v, want %v", got, want)
	}
}

func TestClaudeLiveKeepsCadenceBoundaryAndSelectedModel(t *testing.T) {
	var doc struct {
		Jobs map[string]struct {
			Needs    any
			Outputs  map[string]string
			Strategy struct {
				Matrix map[string][]map[string]string
			}
		}
	}
	if err := yaml.Unmarshal([]byte(readWorkflow(t, "runtime-live-e2e.yml")), &doc); err != nil || doc.Jobs["claude-live"].Needs != "offline" {
		t.Errorf("claude-live needs = %v (parse error: %v); want offline so cadence refusal precedes approval", doc.Jobs["claude-live"].Needs, err)
	}
	if m := doc.Jobs["claude-live"].Strategy.Matrix["include"]; len(m) != 1 || doc.Jobs["offline"].Outputs[strings.Trim(strings.TrimPrefix(m[0]["model"], "${{ needs.offline.outputs."), " }")] != "${{ steps.live_models.outputs.claude_selected }}" {
		t.Errorf("claude-live matrix = %v; want model from the offline output bound to the resolver's claude_selected", m)
	}
	if _, rest, ok := strings.Cut(liveModelsText, "\nclaude.allowed.sonnet="); !ok || strings.Contains(","+strings.SplitN(rest, "\n", 2)[0]+",", ",opus,") {
		t.Errorf("routine sonnet cadence policy is missing or permits opus: %q", strings.SplitN(rest, "\n", 2)[0])
	}
}
