// ABOUTME: live_models.txt holds live model IDs and Claude cadence policies.
// ABOUTME: Go exposes model records; the workflow also reads the policies.
// ABOUTME: No consumer repeats a model ID.
package release

import (
	_ "embed"
	"strings"
)

// liveModelsText is internal/release/live_models.txt, the one place the lane
// model ids are written.
//
//go:embed live_models.txt
var liveModelsText string

// LiveModels maps dotted model keys to IDs from live_models.txt.
// Cadence-policy records are excluded. A model change needs one registry edit.
var LiveModels = parseLiveModels(liveModelsText)

func parseLiveModels(data string) map[string]string {
	models := make(map[string]string)
	for _, line := range strings.Split(data, "\n") {
		key, id, ok := strings.Cut(strings.TrimSpace(line), "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if strings.HasPrefix(key, "claude.allowed.") {
			continue
		}
		models[key] = strings.TrimSpace(id)
	}
	return models
}
