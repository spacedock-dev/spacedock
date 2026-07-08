// ABOUTME: standing-subcommand orchestration — list/show/spawn-standing wire the
// ABOUTME: runtime-neutral _mods parsing to the Claude render + member-exists probe.
package dispatch

import (
	"encoding/json"
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

// runSpawnStanding emits an Agent() spec JSON for a standing-teammate mod, or
// reports already-alive when the target team config already lists the declared
// member. It fails loudly (exit 1, stderr) on: a placeholder team name, missing
// mod file, missing standing: true, missing ## Agent Prompt, missing subagent_type
// / name / model, or a model outside the enum. home is the resolved HOME (the
// member-exists probe's only ~/.claude read, owned by claudeteam). Mirrors
// cmd_spawn_standing.
func runSpawnStanding(home, modPath, teamName string, stdout, stderr io.Writer) int {
	if teamName == "" || teamName == "none" || teamName == "None" {
		fmt.Fprintf(stderr,
			"error: spawn-standing requires a real team name; got '%s'. "+
				"Call TeamCreate first and pass the returned team_name via --team.\n", teamName)
		return 1
	}
	spec, alreadyAlive, errMsg := buildSpawnSpec(home, modPath, teamName)
	if errMsg != "" {
		fmt.Fprintf(stderr, "error: %s\n", errMsg)
		return 1
	}
	if alreadyAlive {
		emitAlreadyAlive(stdout, spec.Name)
		return 0
	}
	return emitSpawnJSON(stdout, spec)
}

// runSpawnStandingAll drives the full standing-teammate inject loop in one call:
// enumerate the workflow's declared standing mods, build a spawn spec for each,
// and emit a JSON ARRAY of the specs to spawn. Empty input (no _mods / no standing
// mods) emits `[]`. Loud failure (exit 1, stderr naming the mod) on the same
// conditions runSpawnStanding fails on — validation is shared via buildSpawnSpec,
// not re-implemented. Composes the _mods scan (discovery) + buildSpawnSpec
// (validate + spec), mirroring build.go's generic show-standing injection.
//
// teamName is the mode discriminator, mirroring build.go's team_name:
//   - legacy (teamName non-empty): each spec carries team_name and is deduped
//     against the team config (MemberExists), so an already-alive member is
//     omitted from the array.
//   - merged (teamName == ""): the .178+ host has no TeamCreate name, so each
//     spec is the merged background shape (name present, team_name absent,
//     run_in_background true). There is no team config keyed by name to dedup
//     against, so every declared standing teammate is emitted; idempotency is the
//     FO's own-roster concern on the merged floor, not a config probe.
func runSpawnStandingAll(home, workflowDir, teamName string, stdout, stderr io.Writer) int {
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
			spec, alreadyAlive, errMsg := buildSpawnSpec(home, modPath, teamName)
			if errMsg != "" {
				fmt.Fprintf(stderr, "error: %s\n", errMsg)
				return 1
			}
			if alreadyAlive {
				continue
			}
			specs = append(specs, spec)
		}
	}
	return claudeteam.EmitPythonJSON(stdout, specs)
}

// buildSpawnSpec parses a standing-teammate mod and either reports it
// already-alive (the team config lists the declared member) or returns the
// Agent() spec to spawn. errMsg is non-empty on any validation failure (missing
// mod, missing standing: true, missing ## Agent Prompt, missing subagent_type /
// name / model, or a model outside the enum); callers turn it into the exit-1
// stderr line. Shared by runSpawnStanding (single) and runSpawnStandingAll (loop)
// so validation lives in one place. On already-alive, only Name is populated.
func buildSpawnSpec(home, modPath, teamName string) (spec spawnSpec, alreadyAlive bool, errMsg string) {
	if !isFile(modPath) {
		return spec, false, fmt.Sprintf("mod file not found: %s", modPath)
	}

	meta, err := ParseModMetadata(modPath)
	if err != nil {
		return spec, false, err.Error()
	}
	if !meta.Standing {
		return spec, false, fmt.Sprintf(
			"mod %s is missing 'standing: true' in frontmatter — not a standing-teammate mod", modPath)
	}
	if meta.AgentPrompt == nil {
		return spec, false, fmt.Sprintf(
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
		return spec, false, fmt.Sprintf(
			"mod %s has no 'subagent_type' in '## Hook: startup' or frontmatter", modPath)
	}
	if declaredName == "" {
		return spec, false, fmt.Sprintf(
			"mod %s has no 'name' in '## Hook: startup' or frontmatter", modPath)
	}
	if model == "" {
		return spec, false, fmt.Sprintf(
			"mod %s has no 'model' in '## Hook: startup' — %s", modPath, modelEnumList)
	}
	if !modelEnum[model] {
		return spec, false, fmt.Sprintf(
			"invalid model '%s' in '## Hook: startup' of %s — %s", model, modPath, modelEnumList)
	}

	// Merged mode (no team name): there is no team config keyed by name to dedup
	// against, so skip the MemberExists probe and emit the merged background shape
	// — name present, team_name absent (nil), run_in_background true. Mirrors
	// build.go's merged dispatch emission for an ensign.
	if teamName == "" {
		runInBackground := true
		return spawnSpec{
			SubagentType:    subagentType,
			Description:     fmt.Sprintf("standing teammate: %s", declaredName),
			Name:            declaredName,
			Model:           model,
			Prompt:          *meta.AgentPrompt,
			RunInBackground: &runInBackground,
		}, false, ""
	}

	// Legacy mode (real team name): dedup the already-alive member against the
	// team config, and carry team_name on the emitted spec.
	if claudeteam.MemberExists(home, teamName, declaredName) {
		return spawnSpec{Name: declaredName}, true, ""
	}

	tn := teamName
	return spawnSpec{
		SubagentType: subagentType,
		Description:  fmt.Sprintf("standing teammate: %s", declaredName),
		Name:         declaredName,
		TeamName:     &tn,
		Model:        model,
		Prompt:       *meta.AgentPrompt,
	}, false, ""
}

// emitAlreadyAlive writes the compact already-alive JSON for a single member:
// a one-line object matching the oracle's json.dumps without indent — Python's
// default ", " / ": " separators, not Go's space-free Marshal. The two keys are
// fixed, so format directly. Escape the name through the same ensure_ascii routine
// the emitters use, so a non-ASCII declared name stays byte-identical to json.dumps.
func emitAlreadyAlive(stdout io.Writer, name string) {
	nameJSON, _ := json.Marshal(name)
	nameJSON = claudeteam.EscapeNonASCII(nameJSON)
	fmt.Fprintf(stdout, "{\"status\": \"already-alive\", \"name\": %s}\n", nameJSON)
}

// spawnSpec is the Agent() call spec spawn-standing emits when the member is not
// yet alive. The Agent tool REQUIRES description, so the spec carries one (else
// the forwarded Agent() call fails InputValidationError and the teammate never
// spawns). Field order matches dispatch build's envelope: subagent_type,
// description, name, team_name, model, prompt. TeamName is a *string with
// omitempty so the legacy spec emits team_name in place while the merged spec
// omits it (absent, not null). RunInBackground is a *bool with omitempty, set
// only on the merged spec (the named-background-teammate shape); the legacy spec
// leaves it nil so the legacy emission stays byte-identical.
type spawnSpec struct {
	SubagentType    string  `json:"subagent_type"`
	Description     string  `json:"description"`
	Name            string  `json:"name"`
	TeamName        *string `json:"team_name,omitempty"`
	Model           string  `json:"model"`
	Prompt          string  `json:"prompt"`
	RunInBackground *bool   `json:"run_in_background,omitempty"`
}

// emitSpawnJSON writes the spec as two-space-indented JSON with a trailing
// newline, matching Python json.dumps(spec, indent=2) + print() byte-for-byte,
// including its ensure_ascii escaping of a non-ASCII Agent Prompt.
func emitSpawnJSON(stdout io.Writer, spec spawnSpec) int {
	return claudeteam.EmitPythonJSON(stdout, spec)
}
