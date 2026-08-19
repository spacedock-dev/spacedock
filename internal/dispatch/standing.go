// ABOUTME: standing-subcommand orchestration — list/show/spawn-standing-all wire
// ABOUTME: the runtime-neutral _mods parsing to the Claude SendMessage render.
package dispatch

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
)

// runListStanding writes the absolute paths of standing: true mods under
// {workflowDir}/_mods/*.md, one per line, sorted by filename. Exit 0 on success
// including zero matches (empty stdout); exit 1 only when the workflow dir is
// unresolvable or a mod is unreadable. The _mods parsing is runtime-neutral, so
// this lives in the generic dispatch package. Mirrors cmd_list_standing.
func runListStanding(workflowDir string, stdout, stderr io.Writer) int {
	if !isDir(workflowDir) {
		fmt.Fprintf(stderr, "error: workflow directory not found: %s\n", workflowDir)
		return 1
	}
	modsDir := filepath.Join(workflowDir, "_mods")
	if !isDir(modsDir) {
		return 0
	}

	var standingPaths []string
	for _, modPath := range sortedModPaths(modsDir) {
		meta, err := ParseModMetadata(modPath)
		if err != nil {
			fmt.Fprintf(stderr, "error: unreadable mod %s: %s\n", modPath, err)
			return 1
		}
		if meta.Standing {
			abs, e := filepath.Abs(modPath)
			if e != nil {
				abs = modPath
			}
			standingPaths = append(standingPaths, abs)
		}
	}
	for _, p := range standingPaths {
		fmt.Fprintln(stdout, p)
	}
	return 0
}

// runShowStanding writes the rendered `### Standing teammates available in your
// team` markdown for the workflow's declared standing teammates. The enumeration
// is runtime-neutral (here); the rendered body is Claude SendMessage routing prose
// (claudeteam.RenderStandingTeammatesSection). Exit 0 on success including the
// degenerate empty case. Mirrors cmd_show_standing.
func runShowStanding(workflowDir string, stdout, stderr io.Writer) int {
	if !isDir(workflowDir) {
		fmt.Fprintf(stderr, "error: workflow directory not found: %s\n", workflowDir)
		return 1
	}
	// A truthy sentinel team name: the filesystem scan is team-agnostic, and the
	// bare-mode short-circuit lives in build upstream of this call.
	teammates := EnumerateDeclaredStandingTeammates(workflowDir, "_show_standing_")
	rendered := claudeteam.RenderStandingTeammatesSection(teammates)
	if rendered != "" {
		fmt.Fprintln(stdout, rendered)
	}
	return 0
}

// runSpawnStandingAll drives the full standing-teammate inject loop in one call:
// enumerate the workflow's declared standing mods, build a spawn spec for each,
// and emit a JSON ARRAY of the specs to spawn. Empty input (no _mods / no standing
// mods) emits `[]`. Loud failure (exit 1, stderr naming the mod) on: missing mod
// file, missing standing: true, missing ## Agent Prompt, missing subagent_type /
// name / model, or a model outside the enum — validation is shared via
// buildSpawnSpec, not re-implemented. Composes the _mods scan (discovery) +
// buildSpawnSpec (validate + spec), mirroring build.go's generic show-standing
// injection. The .178+ host has no TeamCreate name, so every spec is emitted in
// the merged background shape (name present, run_in_background true); there is no
// team config keyed by name to dedup against, so every declared standing teammate
// is emitted — idempotency is the FO's own-roster concern, not a config probe.
func runSpawnStandingAll(workflowDir string, stdout, stderr io.Writer) int {
	if !isDir(workflowDir) {
		fmt.Fprintf(stderr, "error: workflow directory not found: %s\n", workflowDir)
		return 1
	}

	specs := []spawnSpec{}
	modsDir := filepath.Join(workflowDir, "_mods")
	if isDir(modsDir) {
		for _, modPath := range sortedModPaths(modsDir) {
			meta, err := ParseModMetadata(modPath)
			if err != nil {
				fmt.Fprintf(stderr, "error: unreadable mod %s: %s\n", modPath, err)
				return 1
			}
			if !meta.Standing {
				continue
			}
			spec, errMsg := buildSpawnSpec(modPath)
			if errMsg != "" {
				fmt.Fprintf(stderr, "error: %s\n", errMsg)
				return 1
			}
			specs = append(specs, spec)
		}
	}
	return claudeteam.EmitPythonJSON(stdout, specs)
}

// buildSpawnSpec parses a standing-teammate mod and returns the Agent() spec to
// spawn. errMsg is non-empty on any validation failure (missing mod, missing
// standing: true, missing ## Agent Prompt, missing subagent_type / name / model,
// or a model outside the enum); callers turn it into the exit-1 stderr line.
func buildSpawnSpec(modPath string) (spec spawnSpec, errMsg string) {
	if !isFile(modPath) {
		return spec, fmt.Sprintf("mod file not found: %s", modPath)
	}

	meta, err := ParseModMetadata(modPath)
	if err != nil {
		return spec, err.Error()
	}
	if !meta.Standing {
		return spec, fmt.Sprintf(
			"mod %s is missing 'standing: true' in frontmatter — not a standing-teammate mod", modPath)
	}
	if meta.AgentPrompt == nil {
		return spec, fmt.Sprintf(
			"mod %s has no '## Agent Prompt' section — required for standing-teammate spawn", modPath)
	}

	hookConfig := ParseHookStartupSpawnConfig(modPath)
	subagentType := hookConfig["subagent_type"]
	if subagentType == "" {
		subagentType = meta.Frontmatter["subagent_type"]
	}
	declaredName := hookConfig["name"]
	if declaredName == "" {
		declaredName = meta.Name
	}
	model := hookConfig["model"]

	if subagentType == "" {
		return spec, fmt.Sprintf(
			"mod %s has no 'subagent_type' in '## Hook: startup' or frontmatter", modPath)
	}
	if declaredName == "" {
		return spec, fmt.Sprintf(
			"mod %s has no 'name' in '## Hook: startup' or frontmatter", modPath)
	}
	if model == "" {
		return spec, fmt.Sprintf(
			"mod %s has no 'model' in '## Hook: startup' — %s", modPath, modelEnumList)
	}
	if !modelEnum[model] {
		return spec, fmt.Sprintf(
			"invalid model '%s' in '## Hook: startup' of %s — %s", model, modPath, modelEnumList)
	}

	// There is no team config keyed by name to dedup against on the merged
	// floor, so every declared standing teammate emits the merged background
	// shape — name present, run_in_background true. Mirrors build.go's merged
	// dispatch emission for an ensign.
	runInBackground := true
	return spawnSpec{
		SubagentType:    subagentType,
		Description:     fmt.Sprintf("standing teammate: %s", declaredName),
		Name:            declaredName,
		Model:           model,
		Prompt:          *meta.AgentPrompt,
		RunInBackground: &runInBackground,
	}, ""
}

// spawnSpec is the Agent() call spec spawn-standing-all emits for each standing
// teammate. The Agent tool REQUIRES description, so the spec carries one (else
// the forwarded Agent() call fails InputValidationError and the teammate never
// spawns). Field order matches dispatch build's envelope: subagent_type,
// description, name, model, prompt, run_in_background.
type spawnSpec struct {
	SubagentType    string `json:"subagent_type"`
	Description     string `json:"description"`
	Name            string `json:"name"`
	Model           string `json:"model"`
	Prompt          string `json:"prompt"`
	RunInBackground *bool  `json:"run_in_background,omitempty"`
}
