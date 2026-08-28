// ABOUTME: build assembles structured dispatch JSON from stdin + workflow README
// ABOUTME: + entity file, matching the vendored claude-team build oracle.
package dispatch

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
	"github.com/spacedock-dev/spacedock/internal/runtimehost"
	"github.com/spacedock-dev/spacedock/internal/status"
)

const (
	schemaVersion   = 2
	nameMaxLen      = 200
	modelEnumList   = "must be one of: sonnet, opus, haiku, fable"
	dispatchFileDir = "/tmp/spacedock-dispatch"
	// dispatchFileNameMaxLen caps the dispatch filename stem (the merged-mode
	// session token + derived name) so the on-disk file with its .md suffix
	// stays under the common 255-byte filesystem name limit.
	dispatchFileNameMaxLen = 251
	// nameAgentMaxLen is the Agent-tool name ceiling (^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$
	// — 64 chars). A derived name longer than this is rejected by Agent() with
	// InputValidationError, so the worker name is capped to fit.
	nameAgentMaxLen = 64
	// nameCapTarget is the budget the cap fits the worker name into. It reserves
	// nameAgentMaxLen − len("-cycle9") so any single-digit FO cycle suffix the
	// roster appends downstream of dispatch build keeps the name ≤ nameAgentMaxLen.
	nameCapTarget = nameAgentMaxLen - len("-cycle9")
	// sdB32NameIDPrefixLen is the fixed-length id-prefix substituted for the slug
	// component when an id-style: sd-b32 name overflows. A fixed length keeps the
	// embedded token stable across the workflow's lifetime (unlike the shortest-
	// unique display prefix, which lengthens as entities are added). 10 chars
	// leaves ample collision headroom while the resulting name stays ≤ nameCapTarget.
	sdB32NameIDPrefixLen = 10
)

// buildRequiredFields are the stdin keys that must be present and non-null.
var buildRequiredFields = []string{"schema_version", "entity_path", "workflow_dir", "stage", "checklist"}

// namePattern is the dispatch-name regex derived worker names must match.
var namePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$`)

// modelEnum is the Agent-schema model enum declared values are validated
// against — shared by dispatch build (host=claude) and spawn-standing, the
// two sites that render Claude Agent envelopes.
var modelEnum = map[string]bool{"sonnet": true, "opus": true, "haiku": true, "fable": true}

// Pi spawn delivery defaults for `dispatch build --host pi`: pi-subagents
// resolves agents and skills by directory basename only, so the artifact
// carries the generic write-capable agent ("worker") and the basename skill
// ("ensign"); a stage agent: override replaces agent and drops skill (the
// override agent owns its own contract).
const (
	piSpawnAgent = "worker"
	piSpawnSkill = "ensign"
)

// buildOutput is the stdout JSON envelope. Field order is the emission order
// (insertion order in the oracle): schema_version, subagent_type, description,
// fetch_commands, dispatch_file_path, prompt, model, then name /
// run_in_background, then the pi-only agent / skill pair. Model is a *string so an unresolved model serializes as the
// JSON literal null; Name / RunInBackground are *T with omitempty so
// bare-mode dispatches omit the keys entirely (absent, not null). RunInBackground
// is set only on the merged Claude shape (.178+: named background teammate) —
// it carries the worker→lead inter-agent communication half of the dispatch (the
// `name` carries the lead→worker half), so a bare dispatch omits it. Agent/Skill
// are the pi-only spawn binding: set only for
// "host": "pi" fresh dispatches — the artifact itself names the pi-subagents
// agent (piSpawnAgent, or the stage's agent: override) and basename skill
// (piSpawnSkill; omitted on an override) so «worker.spawn» needs no local
// convention. Claude/codex envelopes never carry them (omitempty), so their
// bytes are unchanged.
type buildOutput struct {
	SchemaVersion   int      `json:"schema_version"`
	SubagentType    string   `json:"subagent_type"`
	Description     string   `json:"description"`
	FetchCommands   []string `json:"fetch_commands"`
	DispatchFile    string   `json:"dispatch_file_path"`
	Prompt          string   `json:"prompt"`
	Model           *string  `json:"model"`
	Name            *string  `json:"name,omitempty"`
	RunInBackground *bool    `json:"run_in_background,omitempty"`
	Agent           *string  `json:"agent,omitempty"`
	Skill           *string  `json:"skill,omitempty"`
}

// buildAdvanceOutput is the stdout JSON envelope for `--advance` mode: a pointer
// message for the reuse-advance handle, not a spawn envelope. It carries no
// subagent_type/name/run_in_background — nothing is spawned, so those
// spawn-only fields are absent from the type entirely (not merely omitempty).
// model stays so the FO's reuse-condition-4 comparator can read
// next_stage.effective_model from this output instead of a separate README read.
type buildAdvanceOutput struct {
	SchemaVersion int      `json:"schema_version"`
	Description   string   `json:"description"`
	FetchCommands []string `json:"fetch_commands"`
	DispatchFile  string   `json:"dispatch_file_path"`
	Prompt        string   `json:"prompt"`
	Model         *string  `json:"model"`
}

// buildError prints `error: {msg}` to stderr and returns code (1 by default).
func buildError(stderr io.Writer, code int, format string, a ...any) int {
	fmt.Fprintf(stderr, "error: "+format+"\n", a...)
	return code
}

// runBuild reads a dispatch request from stdin JSON or the flag/file input mode
// and assembles the dispatch envelope on stdout plus the self-contained
// dispatch-file body written to a deterministic path.
func runBuild(probe claudeteam.TeamStateProbe, workflowLauncher string, opts buildOptions, stdin io.Reader, stdout, stderr io.Writer) int {
	if opts.PrintSchema {
		return emitBuildSchema(stdout)
	}

	fields, code := loadBuildFields(opts, stdin, stderr)
	if code != 0 {
		return code
	}

	if opts.Stamp {
		if code := runStamp(opts, fields, stderr); code != 0 {
			return code
		}
	}

	return runBuildFields(probe, workflowLauncher, opts, fields, stdout, stderr)
}

func loadBuildFields(opts buildOptions, stdin io.Reader, stderr io.Writer) (map[string]json.RawMessage, int) {
	switch {
	case opts.ValidateOnly != "":
		raw, err := os.ReadFile(opts.ValidateOnly)
		if err != nil {
			return nil, buildError(stderr, 1, "failed to read validate-only file %q: %s", opts.ValidateOnly, err)
		}
		return decodeBuildFields(raw, stderr)
	case opts.hasRequestFlags():
		return fieldsFromBuildFlags(opts, stderr)
	default:
		raw, err := io.ReadAll(stdin)
		if err != nil {
			return nil, buildError(stderr, 1, "failed to read stdin: %s", err)
		}
		return decodeBuildFields(raw, stderr)
	}
}

func decodeBuildFields(raw []byte, stderr io.Writer) (map[string]json.RawMessage, int) {
	if len(raw) == 0 {
		raw = []byte{}
	}

	// Classify the top-level value the way the oracle does (json.loads then
	// isinstance(inp, dict)): invalid JSON -> "invalid JSON on stdin"; a valid
	// non-object top-level (null, array, scalar) -> "stdin must be a JSON object".
	// A bare-map decode cannot tell these apart -- decoding JSON null into a map
	// succeeds with a nil map, masking the non-object case as a missing field.
	var top interface{}
	if err := json.Unmarshal(raw, &top); err != nil {
		return nil, buildError(stderr, 1, "invalid JSON on stdin: %s", err)
	}
	if _, ok := top.(map[string]interface{}); !ok {
		return nil, buildError(stderr, 1, "stdin must be a JSON object")
	}

	// Distinguish present-but-null from absent (the required-field rule fires for
	// both), so decode into a raw-message map for typed field access.
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return nil, buildError(stderr, 1, "invalid JSON on stdin: %s", err)
	}
	return fields, 0
}

func fieldsFromBuildFlags(opts buildOptions, stderr io.Writer) (map[string]json.RawMessage, int) {
	if opts.EntityPath == "" || opts.Stage == "" || opts.ChecklistFile == "" {
		return nil, buildError(stderr, 2, "flag/file input requires --entity-path, --stage, and --checklist-file")
	}
	checklist, err := readChecklistFile(opts.ChecklistFile)
	if err != nil {
		return nil, buildError(stderr, 1, "failed to read checklist file %q: %s", opts.ChecklistFile, err)
	}
	fields := map[string]json.RawMessage{
		"schema_version": rawJSON(schemaVersion),
		"entity_path":    rawJSON(opts.EntityPath),
		"workflow_dir":   rawJSON(opts.WorkflowDir),
		"stage":          rawJSON(opts.Stage),
		"checklist":      rawJSON(checklist),
		"bare_mode":      rawJSON(opts.BareMode),
		"advance":        rawJSON(opts.Advance),
	}
	if opts.ScopeNotesFile != "" {
		scopeNotes, err := os.ReadFile(opts.ScopeNotesFile)
		if err != nil {
			return nil, buildError(stderr, 1, "failed to read scope-notes file %q: %s", opts.ScopeNotesFile, err)
		}
		fields["scope_notes"] = rawJSON(string(scopeNotes))
	}
	if opts.FeedbackContextFile != "" {
		feedbackContext, err := os.ReadFile(opts.FeedbackContextFile)
		if err != nil {
			return nil, buildError(stderr, 1, "failed to read feedback-context file %q: %s", opts.FeedbackContextFile, err)
		}
		fields["feedback_context"] = rawJSON(string(feedbackContext))
	}
	if opts.FeedbackReflow {
		fields["is_feedback_reflow"] = rawJSON(true)
	}
	return fields, 0
}

func readChecklistFile(path string) ([]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSuffix(line, "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}
		out = append(out, line)
	}
	return out, nil
}

func rawJSON(v any) json.RawMessage {
	raw, _ := json.Marshal(v)
	return raw
}

func resolveBuildHost(flagHost, jsonHost string, getenv func(string) string) (string, error) {
	if flagHost != "" && !validBuildHost(flagHost) {
		return "", fmt.Errorf("unsupported host %q (want claude, codex, or pi)", flagHost)
	}
	if jsonHost != "" && !validBuildHost(jsonHost) {
		return "", fmt.Errorf("unsupported host %q (want claude, codex, or pi)", jsonHost)
	}
	if flagHost != "" && jsonHost != "" && flagHost != jsonHost {
		return "", fmt.Errorf("conflicting explicit host sources: --host=%q, JSON host=%q", flagHost, jsonHost)
	}
	if flagHost != "" {
		return flagHost, nil
	}
	if jsonHost != "" {
		return jsonHost, nil
	}

	// The marker table lives in internal/runtimehost, shared with --version's
	// Runtime line. Detect never errors because the two callers need opposite
	// dispositions of the same facts: a build must REFUSE against an ambiguous
	// host, while --version must REPORT the ambiguity and exit 0. So the
	// detection is shared and the policy stays here.
	host, markers, ambiguous := runtimehost.Detect(getenv)
	if ambiguous {
		return "", fmt.Errorf("ambiguous runtime host sources: multiple runtime markers are set (%s); pass --host claude, codex, or pi", strings.Join(markers, ", "))
	}
	if host != "" {
		return host, nil
	}
	return "", fmt.Errorf("missing host source: pass --host, set JSON host, or run under CODEX_THREAD_ID, CLAUDECODE, PI_CODING_AGENT, or PI_CODING_AGENT_DIR")
}

func validBuildHost(host string) bool {
	return host == "claude" || host == "codex" || host == "pi"
}

func runBuildFields(probe claudeteam.TeamStateProbe, workflowLauncher string, opts buildOptions, fields map[string]json.RawMessage, stdout, stderr io.Writer) int {
	// Rule 1: Required fields present and non-null.
	for _, field := range buildRequiredFields {
		v, ok := fields[field]
		if !ok || isJSONNull(v) {
			return buildError(stderr, 1, "missing required field '%s'", field)
		}
	}

	// Rule 2: Schema version supported (CLEAN-BREAK: v1 rejected). The oracle
	// compares the parsed value against 2; a non-integer or wrong version is
	// rejected with exit 2 and the value rendered as the oracle renders it.
	if !isSchemaVersion(fields["schema_version"]) {
		return buildError(stderr, 2,
			"unsupported input schema_version %s, schema_version: %d required",
			renderSchemaVersion(fields["schema_version"]), schemaVersion)
	}

	entityPath := jsonString(fields["entity_path"])
	workflowDir := opts.WorkflowDir
	if workflowDir == "" {
		workflowDir = jsonString(fields["workflow_dir"])
	}
	stage := jsonString(fields["stage"])
	feedbackContext := optString(fields, "feedback_context")
	scopeNotes := optString(fields, "scope_notes")
	host, err := resolveBuildHost(opts.Host, optString(fields, "host"), os.Getenv)
	if err != nil {
		return buildError(stderr, 1, "%s", err)
	}
	bareMode := optBool(fields, "bare_mode")
	isFeedbackReflow := optBool(fields, "is_feedback_reflow")
	advance := optBool(fields, "advance")

	// Codex fresh dispatch has no anonymous blocking spawn shape: spawn_agent
	// requires the helper-emitted name as task_name.
	if host == "codex" && bareMode {
		return buildError(stderr, 2, "bare_mode is unsupported on host codex; Codex worker.spawn requires a named spawn_agent task")
	}

	// Rule 13: --advance excludes bare mode. A reuse advance presupposes an
	// addressable live worker to message; bare mode dispatches nothing addressable.
	if advance && bareMode {
		return buildError(stderr, 2, "--advance is incompatible with bare_mode (a reuse advance presupposes an addressable worker; bare mode has none)")
	}

	// Merged mode is the Claude .178+ team shape: every non-bare claude dispatch.
	// TeamCreate/TeamDelete are gone on supported hosts, so team membership is
	// established by a named background teammate (Agent(name=…, run_in_background=
	// true)) with no team registry name. It is distinct from bare mode (no name
	// at all, sequential). The merged shape emits a name (the lead→worker
	// SendMessage and reuse-advance handle) and run_in_background (the
	// worker→lead inter-agent communication). Claude-only: codex/pi have their
	// own non-team named-dispatch shapes handled by their host branches.
	mergedMode := !bareMode && host == "claude"

	// FO bootstrap discipline: a bare_mode dispatch with no recent TeamCreate
	// evidence on disk gets an advisory stderr warning (exit stays 0). The evidence
	// read and the warning text both live in the Claude seam; HOME resolution stays
	// generic here. A nil probe (a non-Claude host) emits no advisory — the warning
	// is Claude-specific bootstrap advice, host-neutral by absence.
	if bareMode && probe != nil {
		// HOME resolution stays generic (plain env, no ~/.claude read); the probe
		// owns the team-state read. An unset HOME yields "", on which the probe
		// reports no recent evidence and the advisory fires — the pre-seam behavior.
		if _, _, recent := probe(os.Getenv("HOME"), time.Now()); !recent {
			claudeteam.BareModeAdvisory(stderr)
		}
	}

	// Rule 12: entity_path must be project-root, not worktree-absolute.
	if strings.Contains(entityPath, "/.worktrees/") || strings.HasPrefix(entityPath, ".worktrees/") {
		return buildError(stderr, 1,
			"entity_path must be a project-root absolute path; got worktree path '%s'. "+
				"Pass the project-root location (e.g. '/repo/docs/plans/{slug}.md'), not the worktree copy. "+
				"The helper derives the worktree read target internally.", entityPath)
	}

	// Rule 9: Checklist non-empty (a non-list collapses to the same message).
	checklist, ok := jsonStringList(fields["checklist"])
	if !ok || len(checklist) == 0 {
		return buildError(stderr, 1, "checklist must not be empty")
	}

	// Rule 10: Entity file readable.
	if !isFile(entityPath) {
		return buildError(stderr, 1, "entity file not readable at '%s'", entityPath)
	}

	// Absolutize entityPath against the process cwd, mirroring the workflowDir
	// absolutization below. The dispatched ensign is a separate agent whose cwd is
	// not pinned to the workflow root, so a relative entity_path resolves to a
	// different / nonexistent file there — the "two entity files" divergence. The
	// FO supplies entity_path itself (status emits no path field) and runs with
	// cwd = workflow root, so it naturally passes a relative spelling; absolutizing
	// here makes the entity-read line and the completion signal cwd-independent.
	// After the readability error message, so that diagnostic shows the original
	// spelling.
	if abs, err := filepath.Abs(entityPath); err == nil {
		entityPath = abs
	}

	// Absolutize workflowDir against the process cwd once, so every downstream
	// join — README path, splitRootStateCheckout, the fetch line's --workflow-dir,
	// and the state-commit guidance — inherits an absolute, cwd-independent base.
	// A worktree worker runs with its cwd inside .worktrees/…, where a relative
	// emitted `git -C docs/dev/.spacedock-state` resolves nowhere; absolutizing
	// here makes both halves of the emitted state-commit command absolute.
	if abs, err := filepath.Abs(workflowDir); err == nil {
		workflowDir = abs
	}

	// Rule 11: Workflow README readable.
	readmePath := filepath.Join(workflowDir, "README.md")
	if !isFile(readmePath) {
		return buildError(stderr, 1, "workflow README not found at '%s'", readmePath)
	}
	readmeData, err := os.ReadFile(readmePath)
	if err != nil {
		return buildError(stderr, 1, "failed to read workflow README %q: %s", readmePath, err)
	}
	if _, err := status.StageContextSections(readmeData, stage); err != nil {
		return buildError(stderr, 1, "workflow README %q, stage %q: %s", readmePath, stage, err)
	}
	readmeFields := status.ParseFrontmatterData(readmeData)

	// Parse workflow stages + defaults.
	stages, stageDefaults := status.ParseStagesWithDefaultsData(readmeData)
	if stages == nil {
		return buildError(stderr, 1, "no stages block found in %s", readmePath)
	}

	stageIdx := -1
	for i, s := range stages {
		if s.Name == stage {
			stageIdx = i
			break
		}
	}
	// Rule 3: Stage exists in workflow.
	if stageIdx < 0 {
		return buildError(stderr, 1, "stage '%s' not found in %s", stage, readmePath)
	}
	stageMeta := stages[stageIdx]

	// Resolve effective_model with precedence stage > defaults > null. The
	// declared model's value space is host-scoped: on host=claude it validates
	// against the Agent-schema enum loudly (stage before defaults) and becomes
	// the effective model; on host=codex/pi it is outside that host's
	// dispatch-settable model space, so it is ignored-with-note rather than
	// validated, and the effective model is always null.
	stageModel, stageModelSet := stageMeta.Model()
	defaultsModel, defaultsModelSet := stageDefaults["model"]

	var effectiveModel *string
	modelSource := "null"

	if host == "claude" {
		if stageModelSet && !modelEnum[stageModel] {
			return buildError(stderr, 1,
				"invalid model for stages.states[%d].model: '%s' — %s",
				stageIdx, stageModel, modelEnumList)
		}
		if defaultsModelSet && !modelEnum[defaultsModel] {
			return buildError(stderr, 1,
				"invalid model for stages.defaults.model: '%s' — %s",
				defaultsModel, modelEnumList)
		}
		if stageModelSet {
			m := stageModel
			effectiveModel = &m
			modelSource = "stage"
		} else if defaultsModelSet {
			m := defaultsModel
			effectiveModel = &m
			modelSource = "defaults"
		}
		if effectiveModel != nil {
			fmt.Fprintf(stderr,
				"[build] effective_model=%s (from %s) → Agent model=%s\n",
				*effectiveModel, modelSource, *effectiveModel)
		}
	} else {
		declaredModel, declared := "", false
		if stageModelSet {
			declaredModel, declared = stageModel, true
		} else if defaultsModelSet {
			declaredModel, declared = defaultsModel, true
		}
		if declared {
			fmt.Fprintf(stderr,
				"[build] declared model '%s' ignored on host %s: outside %s's dispatch-settable model space; emitting model=null\n",
				declaredModel, host, host)
		}
	}

	// Rule 4: Stickiness — route on the entity's stamped worktree: field, not the
	// next stage's declared mode.
	entityFields := status.ParseFrontmatter(entityPath)
	entityTitle := entityFields["title"]
	entityWorktree := strings.TrimSpace(entityFields["worktree"])

	var worktreePath, gitRoot string
	if entityWorktree != "" {
		gitRoot = status.FindGitRoot(workflowDir)
		// os.path.join (status.PyJoin) lets an absolute worktree value win, matching the
		// oracle (claude-team:329). filepath.Join would graft an absolute value
		// under gitRoot and double the path, missing the existing worktree dir —
		// the FO stamps absolute worktree: values on live entities.
		worktreePath = status.PyJoin(gitRoot, entityWorktree)
		if info, err := os.Stat(worktreePath); err != nil || !info.IsDir() {
			return buildError(stderr, 1, "worktree path '%s' does not exist", worktreePath)
		}
	} else if stageMeta.Worktree {
		return buildError(stderr, 1, "worktree stage '%s' but entity has no worktree path", stage)
	}

	// Split-root: the README declares a state: checkout, so a worktree stage
	// isolates CODE only — the entity body stays at the FO-passed entity_path.
	// stateCheckout is the resolved absolute state-checkout dir (workflowDir/<state>),
	// the git repo where the entity body lives; "" when the workflow is single-root.
	var stateCheckout string
	if mode, relPath, err := status.ClassifyState(readmeFields["state"]); err == nil && mode == status.StateSplitRoot {
		stateCheckout = filepath.Join(workflowDir, relPath)
	}
	splitRoot := stateCheckout != ""

	// stateBranch is the orphan state branch peers push/pull on a split-root
	// workflow (spacedock-state/<workflow-basename>, or a README state-branch:
	// override). The state-commit guidance names it in the push reminder; an
	// underivable branch (should not happen for a split-root dir) falls back to a
	// branch-neutral reminder inside stateCommitGuidance.
	var stateBranch string
	// stateRemotePresent gates the remote-sync tail of the state-commit guidance:
	// a split-root checkout with no `origin` cannot push/pull, so the guidance
	// degrades to local-only. Probed once here so both guidance callers below agree.
	var stateRemotePresent bool
	if splitRoot {
		stateBranch = strings.TrimSpace(readmeFields["state-branch"])
		if stateBranch == "" {
			stateBranch = "spacedock-state/" + filepath.Base(filepath.Clean(workflowDir))
		}
		stateRemotePresent = stateHasOrigin(stateCheckout)
	}

	// Rule 5: Feedback context required for feedback reflow.
	if isFeedbackReflow && feedbackContext == "" {
		return buildError(stderr, 1,
			"dispatching to feedback target stage '%s' but feedback_context is missing", stage)
	}

	// Rule 8 (retired for the merged floor): a non-bare Claude dispatch is
	// always the merged .178+ team shape (mergedMode above), not an error.
	// TeamCreate is gone on supported hosts, so membership comes from a named
	// background teammate, not a team registry name.

	// Rule 6: subagent_type from the stage agent field.
	subagentType := "spacedock:ensign"
	if agent, ok := stageMeta.Agent(); ok {
		subagentType = agent
	}

	// Derive worker_key, slug, and name (slug-not-stem via EntitySlug). When the
	// readable {workerKey}-{slug}-{stage} form overflows the Agent-tool 64-char
	// ceiling, capWorkerName substitutes the entity id (id-first) or truncates the
	// slug head, preserving the workerKey prefix and -{stage} suffix so the name
	// still decomposes; short names pass through byte-identical.
	workerKey := strings.ReplaceAll(subagentType, ":", "-")
	slug := status.EntitySlug(entityPath)
	idStyle := readmeFields["id-style"]
	derivedName := capWorkerName(workerKey, slug, stage, entityFields["id"], idStyle)

	// Rule 7: Name length and safety.
	if len(derivedName) > nameMaxLen {
		return buildError(stderr, 1, "derived name '%s' exceeds %d characters", derivedName, nameMaxLen)
	}
	if !namePattern.MatchString(derivedName) {
		return buildError(stderr, 1,
			"derived name '%s' contains invalid characters: "+
				"stage name '%s' must match %s (kebab-case "+
				"lowercase letters, digits, and hyphens only). "+
				"Run `status --validate` against the workflow to surface the same "+
				"stage-name error upstream of dispatch.",
			derivedName, stage, namePattern.String())
	}

	// Preflight this immutable README version before writing a dispatch artifact.
	stageSubsection, err := resolveStageContext(readmeData, stage)
	if err != nil {
		if she, ok := err.(*stageHeadingError); ok {
			return buildError(stderr, 1, "%s", she.msg)
		}
		return buildError(stderr, 1, "workflow README %q, stage %q: %s", readmePath, stage, err)
	}
	if stageSubsection == "" {
		return buildError(stderr, 1, "stage '%s' heading not found in %s", stage, readmePath)
	}
	if opts.ValidateOnly == "" && (workflowLauncher == "" || !filepath.IsAbs(workflowLauncher)) {
		return buildError(stderr, 1, "cannot resolve the running spacedock executable; refusing to write a dispatch artifact")
	}

	// --- Prompt assembly ---
	var parts []string

	// 0-1. Operating-contract first-action directive + header. Advance mode
	// skips both: the reused worker already holds its operating contract from
	// initial dispatch, so the file opens with an advance header instead.
	if advance {
		parts = append(parts, fmt.Sprintf(
			"## Advancing to next stage: %s\n\nYou are continuing work on: %s\n", stage, entityTitle))
	} else {
		parts = append(parts, firstActionBlock(host))
		parts = append(parts, fmt.Sprintf("You are working on: %s\n\nStage: %s\n", entityTitle, stage))
	}

	// Stage loading is an exact command, not an ambient launcher policy. Resolve
	// the launcher once at generation and put its shell-quoted absolute path into
	// every generated helper invocation. Separate shell calls therefore do not
	// depend on an export persisting, and an explicitly supplied candidate binary
	// remains a legitimate launcher under test.
	fetchCommands := []string{
		fmt.Sprintf("%s dispatch show-stage-def --workflow-dir %s --stage %s",
			shlexQuote(workflowLauncher), shlexQuote(workflowDir), shlexQuote(stage)),
	}

	// 3. Worktree instructions (conditional). Under split root the state-commit
	// guidance applies to every stage, worktree or not: the worktree branch
	// prepends CODE-directory/branch lines, a non-worktree stage emits the
	// standalone guidance. The guidance carries the resolved absolute state
	// checkout (workflowDir/<state>, the git repo holding the entity body) and
	// entity path — never literal {state_checkout}/{entity_path} brace tokens.
	var worktreeEntityPath string
	if worktreePath != "" {
		branch := fmt.Sprintf("%s/%s", workerKey, slug)
		if splitRoot {
			parts = append(parts, fmt.Sprintf(
				"Your working directory for CODE is %s\n"+
					"All CODE reads and writes MUST use paths under %s.\n"+
					"Your git branch for CODE is %s. All CODE commits MUST be on "+
					"this branch in the worktree. Do NOT switch branches or commit "+
					"code to main.\n"+
					"%s",
				worktreePath, worktreePath, branch,
				stateCommitGuidance(stateCheckout, entityPath, stateBranch, stateRemotePresent)))
		} else {
			parts = append(parts, fmt.Sprintf(
				"Your working directory is %s\n"+
					"All file reads and writes MUST use paths under %s.\n"+
					"Your git branch is %s. All commits MUST be on this branch. "+
					"Do NOT switch branches or commit to main.\n",
				worktreePath, worktreePath, branch))
		}
	} else if splitRoot {
		parts = append(parts, stateCommitGuidance(stateCheckout, entityPath, stateBranch, stateRemotePresent))
	}

	// 4. Entity-read instruction. Under split root the entity lives in the state
	// checkout; a non-split worktree stage rewrites the path into the worktree.
	// Advance mode replaces the "Read ... for the spec" wording with a
	// continue-on-entity instruction — the worker already knows the entity, it is
	// resuming work on it — using the identical path resolution as fresh dispatch.
	if worktreePath != "" && !splitRoot {
		entityRel := pyRelpath(entityPath, gitRoot)
		worktreeEntityPath = status.PyJoin(worktreePath, entityRel)
		if advance {
			parts = append(parts, fmt.Sprintf(
				"Continue working on the entity at %s.\n", worktreeEntityPath))
		} else {
			parts = append(parts, fmt.Sprintf(
				"Read the entity file at %s for the full spec. It contains:\n", worktreeEntityPath))
		}
	} else {
		if advance {
			parts = append(parts, fmt.Sprintf(
				"Continue working on the entity at %s.\n", entityPath))
		} else {
			parts = append(parts, fmt.Sprintf(
				"Read the entity file at %s for the current spec.\n", entityPath))
		}
	}

	// 6. Feedback context (conditional).
	if feedbackContext != "" {
		parts = append(parts, fmt.Sprintf("### Feedback from prior review\n\n%s\n", feedbackContext))
	}

	// 7. Scope notes (conditional).
	if scopeNotes != "" {
		parts = append(parts, scopeNotes+"\n")
	}

	// 8. Completion checklist + Summary slot.
	checklistText := strings.Join(checklist, "\n")
	parts = append(parts, fmt.Sprintf(
		"### Completion checklist\n\n%s\n\n### Summary\n{brief description of what was accomplished}\n",
		checklistText))

	// 9 (retired): standing-teammate auto-injection via a legacy team_name only
	// ever fired in the deleted legacy branch — merged and bare dispatches always
	// omitted the command (documented behavior, unchanged by this removal). The
	// standing-teammate flow stays reachable directly via show-standing /
	// spawn-standing-all.

	fetchLines := []string{"### Fetch commands", ""}
	for _, command := range fetchCommands {
		fetchLines = append(fetchLines, "    "+command)
	}
	parts = append(parts, strings.Join(fetchLines, "\n"))

	// 10. Completion signal (Claude merged mode or Codex named dispatch). The
	// worker→lead completion target is pinned to the single name "team-lead"
	// (AC-6), matching the ensign runtime's completion contract. Bare mode (no
	// name) still omits it; its dispatch blocks and returns inline.
	if !bareMode && (mergedMode || host == "codex") {
		entityFileRef := entityPath
		if worktreePath != "" && !splitRoot {
			entityFileRef = worktreeEntityPath
		}
		parts = append(parts, completionSignalBlock(host, entityTitle, stage, entityFileRef))
	}

	dispatchBody := strings.Join(parts, "\n")
	if opts.ValidateOnly != "" {
		return 0
	}

	// v2 file-pointer: write the body to a collision-free path under the shared
	// dispatch dir; emit a tiny prompt the ensign Reads on first action. A bare
	// {derivedName}.md is identical across every dispatch of one slug+stage, so two
	// concurrent FOs — or back-to-back runs of one fixture — alias the same file
	// and an ensign can Read a STALE prior dispatch's entity pointer. On the
	// merged floor there is no team registry name, so the per-session auto-team
	// id ($CLAUDE_CODE_SESSION_ID) is the disambiguator instead — two concurrent
	// merged FOs run distinct sessions, so it separates their dispatch files.
	// Dispatches with neither (bare, or a merged dispatch with no session id in
	// env) keep the plain derived name. derivedName stays the readable
	// team-member name; only the on-disk path is keyed.
	dispatchFileName := derivedName
	if mergedMode {
		if sessionToken := pathSafeSessionToken(os.Getenv("CLAUDE_CODE_SESSION_ID")); sessionToken != "" {
			dispatchFileName = sessionToken + "-" + derivedName
		}
	}
	if advance {
		// -advance suffix so an advance file can never alias a fresh-dispatch file
		// for the same slug+stage: a fresh dispatch after a failed advance would
		// otherwise collide with the stale advance body at the bare derivedName path.
		dispatchFileName += "-advance"
	}
	if len(dispatchFileName) > dispatchFileNameMaxLen {
		return buildError(stderr, 1,
			"dispatch filename '%s' exceeds %d characters", dispatchFileName, dispatchFileNameMaxLen)
	}
	dispatchFilePath := filepath.Join(dispatchFileDir, dispatchFileName+".md")
	if err := os.MkdirAll(dispatchFileDir, 0o755); err != nil {
		fmt.Fprintf(stderr, "dispatch_file_write_failed: %s: %s\n", dispatchFilePath, err)
		return 1
	}
	if err := os.WriteFile(dispatchFilePath, []byte(dispatchBody), 0o644); err != nil {
		fmt.Fprintf(stderr, "dispatch_file_write_failed: %s: %s\n", dispatchFilePath, err)
		return 1
	}

	// Self-contained artifact; pointer-only transport. Assignment payload stays in
	// dispatchBody. These renderers carry only the file locator and fixed host or
	// stage routing metadata; advance deliberately retains its stage label.
	var prompt string
	if advance {
		prompt = dispatchAdvancePointerPrompt(stage, dispatchFilePath)
	} else {
		prompt = dispatchPointerPrompt(host, dispatchFilePath)
	}

	if advance {
		// Nothing is spawned, so the envelope carries no subagent_type/name/
		// run_in_background — only the pointer message and the fields
		// the FO's reuse-condition-4 comparator still needs (model).
		outAdvance := buildAdvanceOutput{
			SchemaVersion: schemaVersion,
			Description:   fmt.Sprintf("%s: %s", entityTitle, stage),
			FetchCommands: fetchCommands,
			DispatchFile:  dispatchFilePath,
			Prompt:        prompt,
			Model:         effectiveModel,
		}
		return emitBuildJSON(stdout, outAdvance)
	}

	out := buildOutput{
		SchemaVersion: schemaVersion,
		SubagentType:  subagentType,
		Description:   fmt.Sprintf("%s: %s", entityTitle, stage),
		FetchCommands: fetchCommands,
		DispatchFile:  dispatchFilePath,
		Prompt:        prompt,
		Model:         effectiveModel,
	}
	if !bareMode {
		out.Name = &derivedName
	}
	if mergedMode {
		// The merged dispatch is Agent(name=…, run_in_background=true): the named
		// background teammate whose up-channel (SendMessage to the lead) is what
		// run_in_background confers. The FO maps this field to Agent verbatim.
		runInBackground := true
		out.RunInBackground = &runInBackground
	}
	if host == "pi" {
		if agent, ok := stageMeta.Agent(); ok {
			// Override owns its contract: replace agent, omit skill.
			out.Agent = &agent
		} else {
			agent, skill := piSpawnAgent, piSpawnSkill
			out.Agent = &agent
			out.Skill = &skill
		}
	}

	return emitBuildJSON(stdout, out)
}

// capWorkerName builds the worker name {workerKey}-{slug}-{stage}, capping it to
// the Agent-tool 64-char ceiling when the readable form overflows. The cap is
// id-first: an id-style: sd-b32 entity substitutes a fixed-length prefix of its
// stored id for the slug component (a stable, decomposition-safe token); a
// sequential id substitutes the whole numeric id (short, never overflows); an
// id-less slug workflow truncates the slug head to fit. The workerKey prefix and
// -{stage} suffix are preserved verbatim in every form so reconcile's decompose()
// still peels the stage and strips the worker prefix. Short names (≤64) return
// unchanged — no id substitution, no truncation — so existing readable names and
// golden fixtures are byte-identical.
func capWorkerName(workerKey, slug, stage, id, idStyle string) string {
	full := fmt.Sprintf("%s-%s-%s", workerKey, slug, stage)
	if len(full) <= nameAgentMaxLen {
		return full
	}

	// id-first: replace the slug component with a stable id-derived token. sd-b32
	// ids are 24-char tokens from the namePattern-safe alphabet, so a fixed-length
	// prefix is hyphen-free and decomposition-safe; sequential ids are short
	// integers used whole. Either keeps the prefix + suffix intact.
	idToken := ""
	switch idStyle {
	case "sd-b32":
		if status.IsValidSDB32ID(id) {
			idToken = id[:sdB32NameIDPrefixLen]
		}
	case "sequential":
		if status.IsDigits(id) {
			idToken = id
		}
	}
	if idToken != "" {
		return fmt.Sprintf("%s-%s-%s", workerKey, idToken, stage)
	}

	// id-less (id-style: slug) fallback: truncate the slug head so the whole name
	// fits nameCapTarget (reserving cycle headroom), trimming any trailing hyphen
	// so the result matches namePattern and carries no `--`.
	fixed := len(workerKey) + 2 + len(stage) // "-{slug}-" framing minus the slug
	budget := nameCapTarget - fixed
	if budget < 1 {
		budget = 1
	}
	truncated := slug
	if len(truncated) > budget {
		truncated = truncated[:budget]
	}
	truncated = strings.TrimRight(truncated, "-")
	return fmt.Sprintf("%s-%s-%s", workerKey, truncated, stage)
}

// sessionTokenMaxLen caps the merged-mode session disambiguator prepended to the
// dispatch filename. A Claude session id is a UUID-shaped token (~36 chars); the
// cap keeps the combined {sessionToken}-{derivedName}.md well under the
// filesystem name limit while staying long enough that two concurrent sessions
// never collide (their ids differ in the high-entropy leading characters).
const sessionTokenMaxLen = 36

// pathSafeSessionToken sanitizes a session id into the kebab character class the
// dispatch filename already enforces (lowercase letters, digits, hyphens), so a
// merged-mode dispatch keyed on $CLAUDE_CODE_SESSION_ID never builds an
// unsanitized /tmp path. Uppercase is lowered; any other character becomes a
// hyphen; leading/trailing hyphens are trimmed and the result capped. An empty or
// all-unsafe id yields "" so the caller falls back to the plain derived name.
func pathSafeSessionToken(sessionID string) string {
	if sessionID == "" {
		return ""
	}
	var b strings.Builder
	for _, r := range strings.ToLower(sessionID) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteRune('-')
		}
	}
	token := strings.Trim(b.String(), "-")
	if len(token) > sessionTokenMaxLen {
		token = strings.TrimRight(token[:sessionTokenMaxLen], "-")
	}
	return token
}

func firstActionBlock(host string) string {
	if host == "codex" {
		return "## First action\n" +
			"\n" +
			"Read this dispatch file directly and treat its content as your stage-specific assignment.\n" +
			"\n" +
			"The outer fresh-worker prompt invokes `$spacedock:ensign` before this pointer. The installed " +
			"skill supplies the shared ensign discipline (stage-report format, polling, worktree " +
			"ownership, and completion signal protocol); this file supplies the stage-specific " +
			"assignment. Do not try to invoke a Claude skill wrapper; Codex dispatch uses this file " +
			"pointer after the outer bootstrap.\n"
	}
	if host == "pi" {
		return "## First action\n" +
			"\n" +
			"Before anything else, load the ensign discipline: run `/skill:ensign` " +
			"(Pi's skill-invoke slash command), or if that is unavailable, read " +
			"`skills/ensign/SKILL.md` and its `references/` directly. This loads the " +
			"shared ensign discipline (stage-report format, polling, worktree " +
			"ownership, completion signal protocol).\n" +
			"\n" +
			"Then read this dispatch file and treat its content as your " +
			"stage-specific assignment. Pi dispatch is delivered through a Pi-native " +
			"substrate such as pi-subagents; the Pi subagent completion result is the " +
			"completion signal observed by the first officer. Do not emit Claude " +
			"team-tool calls.\n"
	}
	return "## First action\n" +
		"\n" +
		"Before anything else, invoke your operating contract:\n" +
		"\n" +
		"    Skill(skill=\"spacedock:ensign\")\n" +
		"\n" +
		"This loads the shared ensign discipline (stage-report format, background-task " +
		"polling, worktree ownership, completion signal protocol). The call is safe " +
		"to call more than once; if the agent-definition preload ever starts " +
		"working, calling it again is a no-op (the skill content is re-loaded but " +
		"has no behavioral effect). Do not paraphrase; call the tool.\n"
}

func completionSignalBlock(host, entityTitle, stage, entityFileRef string) string {
	if host == "codex" {
		return fmt.Sprintf(
			"\n\n### Completion Signal\n\n"+
				"When you finish (after all commits and stage report writes are done), send one concise final "+
				"message in this Codex worker thread:\n\n"+
				"    Done: %s completed %s. Report written to %s.\n\n"+
				"The first officer observes completion through the Codex final-status notification in the FO mailbox. "+
				"Do not emit a Claude `SendMessage` call; the Codex mailbox notification is the completion signal.",
			entityTitle, stage, entityFileRef)
	}
	if host == "pi" {
		return fmt.Sprintf(
			"\n\n### Completion Signal\n\n"+
				"When you finish (after all commits and stage report writes are done), return one concise final message in this Pi worker turn:\n\n"+
				"    Done: %s completed %s. Report written to %s.\n\n"+
				"The first officer observes completion through the Pi subagent completion result. "+
				"Do not emit Claude message-tool calls; Pi dispatch completion is the worker's final result.",
			entityTitle, stage, entityFileRef)
	}
	return fmt.Sprintf(
		"\n\n### Completion Signal\n\n"+
			"When you finish (after all commits and stage report writes are done), "+
			"your last action MUST be:\n\n"+
			"    SendMessage(to=\"team-lead\", message=\"Done: %s completed "+
			"%s. Report written to %s.\")\n\n"+
			"**If you are the first officer forwarding this prompt to Agent():** copy "+
			"the entire block above into `Agent(prompt=...)` character-for-character. "+
			"Do NOT paraphrase `SendMessage(to=\"team-lead\", ...)` as \"SendMessage with "+
			"to='team-lead'\" or any other English rewrite — the parenthesis-equals "+
			"syntax is the literal call the ensign must emit, not a description of one.",
		entityTitle, stage, entityFileRef)
}

func dispatchPointerPrompt(host, dispatchFilePath string) string {
	if host == "codex" {
		return fmt.Sprintf("$spacedock:ensign; then Read %s and treat its content as your assignment.", dispatchFilePath)
	}
	if host == "pi" {
		return fmt.Sprintf("Read %s and treat its content as your assignment.", dispatchFilePath)
	}
	return fmt.Sprintf(
		"Skill(skill=\"spacedock:ensign\"); then Read %s and treat its content as your assignment.",
		dispatchFilePath)
}

// dispatchAdvancePointerPrompt is the reuse-advance pointer message sent to a
// live worker in place of the hand-assembled verbatim-stage-section template.
// Unlike dispatchPointerPrompt, the wording is host-uniform: a reused worker
// already holds its operating contract from initial dispatch (Claude's
// Skill(...) invocation included), so no host branches on a skill-wrapper
// clause here.
func dispatchAdvancePointerPrompt(stage, dispatchFilePath string) string {
	return fmt.Sprintf(
		"Advancing to next stage: %s.\n\nRead %s and treat its content as your next-stage assignment.",
		stage, dispatchFilePath)
}

// stateHasOrigin reports whether the state checkout has a named `origin` remote,
// the named-remote question the split-root sync contract pushes/pulls against —
// true iff `git remote get-url origin` exits 0. Network-free (unlike ls-remote)
// and discriminating (unlike a bare `git remote`, which exits 0 with no output).
// A non-repo dir or any other git failure reports false, degrading the checkout
// to local-only. This is the dispatch-package local exec the design pairs with
// the status package's runGitCmd-backed probe; both ask the identical question.
func stateHasOrigin(checkout string) bool {
	cmd := exec.Command("git", "-C", checkout, "remote", "get-url", "origin")
	return cmd.Run() == nil
}

// stateCommitGuidance is the split-root state-commit instruction, shared by the
// worktree and non-worktree branches so the wording lives in one place. It
// substitutes the resolved absolute state checkout and entity paths into the
// path-scoped commit command — never literal {state_checkout}/{entity_path}
// brace tokens — and carries the concurrency-safe "never a bare git add -A"
// rule that governs every split-root stage. After the commit the remote-sync
// tail diverges on hasOrigin: with an origin it reminds the worker to push the
// orphan state branch peers share and `pull --rebase` on a rejection
// (stateBranch named verbatim when resolved, else a branch-neutral reminder);
// without one it tells the worker the checkout is local-only and to skip
// push/pull — the path-scoped commit instruction is unchanged either way.
func stateCommitGuidance(stateCheckout, entityPath, stateBranch string, hasOrigin bool) string {
	var pushReminder string
	if hasOrigin {
		pushReminder = "Then push the state branch so peers see your entity/report: " +
			"`git -C " + stateCheckout + " push origin "
		if stateBranch != "" {
			pushReminder += stateBranch + "`"
		} else {
			pushReminder += "<state-branch>`"
		}
		pushReminder += "; on a non-fast-forward rejection, " +
			"`git -C " + stateCheckout + " pull --rebase origin "
		if stateBranch != "" {
			pushReminder += stateBranch + "`"
		} else {
			pushReminder += "<state-branch>`"
		}
		pushReminder += " then re-push.\n"
	} else {
		pushReminder = "This state checkout has no `origin` remote — commit " +
			"path-scoped locally as above; do NOT run `git push`/`git pull` " +
			"(there is no remote to sync). State is local-only and will not " +
			"survive on a shared remote until an `origin` is configured.\n"
	}

	return fmt.Sprintf(
		"This workflow is split-root: the entity body and your stage report "+
			"live in the shared state checkout, not alongside the code. Write and "+
			"commit them to the state checkout at the entity path below. Commit "+
			"state path-scoped — "+
			"`git -C %s add %s && "+
			"git -C %s commit -m \"...\" -- %s` — "+
			"never a bare `git add -A` or bare `git commit` in the state "+
			"checkout (a bare commit cross-attributes a concurrent writer's "+
			"staged entity). Retry on index.lock contention after a short wait. %s",
		stateCheckout, entityPath, stateCheckout, entityPath, pushReminder)
}

// emitBuildJSON writes out as two-space-indented JSON with a trailing newline,
// matching Python json.dumps(indent=2) followed by print() byte-for-byte,
// including its ensure_ascii escaping of any non-ASCII entity title / prompt.
// out is buildOutput for a spawn envelope or buildAdvanceOutput for --advance.
func emitBuildJSON(stdout io.Writer, out any) int {
	return claudeteam.EmitPythonJSON(stdout, out)
}

func emitBuildSchema(stdout io.Writer) int {
	schema := map[string]any{
		"$schema":              "https://json-schema.org/draft/2020-12/schema",
		"title":                "spacedock dispatch build request",
		"type":                 "object",
		"additionalProperties": true,
		"required":             buildRequiredFields,
		"properties": map[string]any{
			"schema_version": map[string]any{
				"const": schemaVersion,
			},
			"entity_path": map[string]any{
				"type": "string",
			},
			"workflow_dir": map[string]any{
				"type": "string",
			},
			"stage": map[string]any{
				"type": "string",
			},
			"checklist": map[string]any{
				"type":  "array",
				"items": map[string]any{"type": "string"},
			},
			"feedback_context": map[string]any{
				"type": []string{"string", "null"},
			},
			"scope_notes": map[string]any{
				"type": []string{"string", "null"},
			},
			"bare_mode": map[string]any{
				"type": "boolean",
			},
			"is_feedback_reflow": map[string]any{
				"type": "boolean",
			},
			"advance": map[string]any{
				"type": "boolean",
			},
			"host": map[string]any{
				"type": "string",
				"enum": []string{"claude", "codex", "pi"},
			},
		},
	}
	return claudeteam.EmitPythonJSON(stdout, schema)
}
