// ABOUTME: The single source of truth for the live E2E lane model ids. The data
// ABOUTME: file live_models.txt holds one `key=id` line per lane; this package
// ABOUTME: embeds it for Go, and the live workflow reads the same file directly,
// ABOUTME: so no consumer repeats a model literal.
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

// LiveModels maps each dotted lane key to its pinned model id, parsed from
// live_models.txt: claude.sonnet, claude.opus, codex.exec, pi.oauth, and
// pi.api-key. A model change is one edit to that file.
var LiveModels = parseLiveModels(liveModelsText)

func parseLiveModels(data string) map[string]string {
	models := make(map[string]string)
	for _, line := range strings.Split(data, "\n") {
		key, id, ok := strings.Cut(strings.TrimSpace(line), "=")
		if !ok {
			continue
		}
		models[strings.TrimSpace(key)] = strings.TrimSpace(id)
	}
	return models
}
