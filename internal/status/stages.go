// ABOUTME: Stage parser — typed-decodes the README frontmatter's stages: block
// ABOUTME: through gopkg.in/yaml.v3 (yaml.Node preserves declaration order).
package status

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// stageNameRe is the dispatch-name regex stage names must match.
var stageNameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$`)

// Stage is a resolved workflow stage with defaults applied. Optional carried
// fields (feedback-to, agent, fresh, model) are kept verbatim when present.
type Stage struct {
	Name        string
	Worktree    bool
	concurrency int
	gate        bool
	terminal    bool
	initial     bool
	optional    map[string]string
}

// Model returns the stage's declared model, with ok=false when the stage
// carries no model field (distinct from an empty-string value).
func (s Stage) Model() (string, bool) {
	v, ok := s.optional["model"]
	return v, ok
}

// Agent returns the stage's declared agent (subagent_type), with ok=false when
// the stage carries no agent field.
func (s Stage) Agent() (string, bool) {
	v, ok := s.optional["agent"]
	return v, ok
}

// parseStagesBlock parses the stages: block from README frontmatter, returning
// ordered stages with resolved defaults, or nil when there is no stages: block
// or no states entries.
func parseStagesBlock(path string) []Stage {
	stages, _ := ParseStagesWithDefaults(path)
	return stages
}

// ParseStagesWithDefaults returns the ordered stages and the raw stages.defaults
// map. The README frontmatter is parsed by yaml.v3 — top-level `stages:` is a
// mapping with optional `defaults:` (mapping) and `states:` (sequence of
// mappings); yaml.Node preserves declaration order and natively strips
// inline `# comments` on scalar values.
func ParseStagesWithDefaults(path string) ([]Stage, map[string]string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, map[string]string{}
	}
	stagesNode := stagesNodeFromFrontmatter(data)
	if stagesNode == nil {
		return nil, map[string]string{}
	}

	defaults := scalarMap(mappingValue(stagesNode, "defaults"))
	statesNode := mappingValue(stagesNode, "states")
	if statesNode == nil || statesNode.Kind != yaml.SequenceNode || len(statesNode.Content) == 0 {
		return nil, defaults
	}

	defaultWorktree := strings.EqualFold(getOr(defaults, "worktree", "false"), "true")
	defaultConcurrency := atoiOr(getOr(defaults, "concurrency", "2"), 2)

	result := make([]Stage, 0, len(statesNode.Content))
	for _, item := range statesNode.Content {
		if item.Kind != yaml.MappingNode {
			continue
		}
		fields := scalarMap(item)
		s := Stage{
			Name:        fields["name"],
			Worktree:    strings.EqualFold(getOr(fields, "worktree", boolStr(defaultWorktree)), "true"),
			concurrency: atoiOr(getOr(fields, "concurrency", strconv.Itoa(defaultConcurrency)), defaultConcurrency),
			gate:        strings.EqualFold(getOr(fields, "gate", "false"), "true"),
			terminal:    strings.EqualFold(getOr(fields, "terminal", "false"), "true"),
			initial:     strings.EqualFold(getOr(fields, "initial", "false"), "true"),
			optional:    map[string]string{},
		}
		for _, of := range []string{"feedback-to", "agent", "fresh", "model"} {
			if v, ok := fields[of]; ok {
				s.optional[of] = v
			}
		}
		result = append(result, s)
	}
	if len(result) == 0 {
		return nil, defaults
	}
	return result, defaults
}

// stagesNodeFromFrontmatter parses the in-fence frontmatter slice and returns
// the `stages:` mapping value-node, or nil when absent. Shares the hand-rolled
// fence finder with the reader so the YAML body boundary is identical.
func stagesNodeFromFrontmatter(data []byte) *yaml.Node {
	slice := frontmatterSlice(data)
	if len(slice) == 0 {
		return nil
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(slice, &doc); err != nil {
		return nil
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return nil
	}
	return mappingValue(doc.Content[0], "stages")
}

// mappingValue returns the value node for key in a yaml.v3 MappingNode, or nil
// when the key is absent or the input is not a mapping.
func mappingValue(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		k := mapping.Content[i]
		if k.Kind == yaml.ScalarNode && k.Value == key {
			return mapping.Content[i+1]
		}
	}
	return nil
}

// scalarMap collects the top-level scalar key/value pairs from a yaml.v3
// MappingNode. Nested mappings/sequences render as "" (matches the
// indented-lines-ignored semantic of the frontmatter reader). yaml.v3's
// Node.Value is the literal scalar text with comments + quotes stripped, so
// `concurrency: 5  # debate` yields "5" and `worktree: false` yields "false".
func scalarMap(mapping *yaml.Node) map[string]string {
	out := map[string]string{}
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return out
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		k := mapping.Content[i]
		v := mapping.Content[i+1]
		if k.Kind != yaml.ScalarNode {
			continue
		}
		if v.Kind == yaml.ScalarNode {
			out[k.Value] = v.Value
		} else {
			out[k.Value] = ""
		}
	}
	return out
}

func getOr(m map[string]string, k, def string) string {
	if v, ok := m[k]; ok {
		return v
	}
	return def
}

func atoiOr(s string, def int) int {
	if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
		return n
	}
	return def
}

func boolStr(b bool) string {
	if b {
		return "True"
	}
	return "False"
}
