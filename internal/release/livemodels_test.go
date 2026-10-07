package release

import (
	"reflect"
	"testing"
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
