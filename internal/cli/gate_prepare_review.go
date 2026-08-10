package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/spacedock-dev/spacedock/internal/gates"
	"github.com/spacedock-dev/spacedock/internal/statesync"
	"github.com/spacedock-dev/spacedock/internal/status"
)

var prepareReviewSync = syncActiveEntity

type prepareReviewArgs struct {
	publish        bool
	json           bool
	workflowDir    string
	question       string
	artifact       string
	references     []string
	summary        string
	recommendation string
}
type prepareReviewEvidence struct {
	Entity     map[string]string   `json:"entity"`
	Stage      map[string]string   `json:"stage"`
	Candidates []map[string]string `json:"candidates"`
}
type prepareReviewOutput struct {
	Command      string               `json:"command"`
	Mode         string               `json:"mode"`
	Phase        string               `json:"phase,omitempty"`
	LaunchCWD    string               `json:"launch_cwd"`
	PublishArgv  []string             `json:"publish_argv,omitempty"`
	Entity       map[string]string    `json:"entity"`
	Stage        map[string]string    `json:"stage"`
	Candidates   []map[string]string  `json:"candidates"`
	Preparation  *gates.PrepareResult `json:"preparation,omitempty"`
	Sync         map[string]any       `json:"sync,omitempty"`
	Checklist    json.RawMessage      `json:"checklist,omitempty"`
	Acceptance   json.RawMessage      `json:"acceptance,omitempty"`
	Presentation map[string]any       `json:"presentation,omitempty"`
}

func runGatePrepareReview(slug string, args []string, dir string, stdout, stderr io.Writer) int {
	parsed, err := parsePrepareReviewArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, "Error:", err)
		return 2
	}
	if !parsed.json {
		fmt.Fprintln(stderr, "Error: gate prepare-review requires --json")
		return 2
	}
	definitionDir := parsed.workflowDir
	if definitionDir == "" {
		var code int
		definitionDir, code = status.ResolveWorkflowDir(dir, stderr)
		if code != 0 {
			return code
		}
	} else if !filepath.IsAbs(definitionDir) {
		definitionDir = filepath.Join(dir, definitionDir)
	}
	definitionDir = filepath.Clean(definitionDir)
	evidence, err := loadPrepareReviewEvidence(definitionDir, slug)
	if err != nil {
		fmt.Fprintln(stderr, "Error:", err)
		return 1
	}
	out := prepareReviewOutput{
		Command: "gate prepare-review", Mode: "inspect", LaunchCWD: definitionDir,
		Entity: evidence.Entity, Stage: evidence.Stage, Candidates: evidence.Candidates,
	}
	if !parsed.publish {
		out.PublishArgv = prepareReviewArgv(slug, definitionDir)
		return emitPrepareReviewJSON(stdout, out)
	}
	out.Mode = "publish"
	if err := validatePrepareReviewInputs(parsed, evidence.Candidates); err != nil {
		fmt.Fprintln(stderr, "Error:", err)
		return 1
	}
	entityPath := filepath.Join(definitionDir, filepath.FromSlash(evidence.Entity["input_path"]))
	checkout, branch, mode, code := resolveStateCheckout("gate prepare-review", definitionDir, stderr)
	if code != 0 {
		return code
	}
	if mode != status.StateSplitRoot {
		fmt.Fprintln(stderr, "Error: gate prepare-review requires a split-root workflow")
		return 1
	}
	target, found, err := resolveEntityCommitTarget(checkout, slug)
	if err != nil || !found || target.scope != entityScopeActive {
		if err == nil {
			err = fmt.Errorf("active entity commit unit not found")
		}
		fmt.Fprintln(stderr, "Error:", err)
		return 1
	}
	if err := requireCleanPrepareReviewUnit(checkout, target.entityPaths); err != nil {
		fmt.Fprintln(stderr, "Error:", err)
		return 1
	}

	prepareInput := gates.PrepareInput{
		WorkflowDir: definitionDir,
		Question:    parsed.question,
		Artifact:    filepath.Join(definitionDir, filepath.FromSlash(parsed.artifact)),
		Summary:     parsed.summary,
	}
	for _, ref := range parsed.references {
		prepareInput.References = append(prepareInput.References, filepath.Join(definitionDir, filepath.FromSlash(ref)))
	}
	var syncOutcome statesync.Outcome
	var committed bool
	result, err := gates.PrepareTransactional(entityPath, prepareInput, func(gates.PrepareResult) (bool, error) {
		committed, syncOutcome, err = prepareReviewSync(checkout, branch, slug, "gate: prepare-review "+slug)
		if err != nil {
			restorePrepareReviewIndex(checkout, target.entityPaths)
			return false, err
		}
		switch syncOutcome.Result {
		case statesync.ResultHalted:
			return true, fmt.Errorf("state publication halted: peer=%s conflicts=%s", syncOutcome.PeerCommit, strings.Join(syncOutcome.ConflictingPaths, ","))
		case statesync.ResultFailed:
			return true, fmt.Errorf("state publication failed: %s", strings.TrimSpace(syncOutcome.Detail))
		default:
			return committed, nil
		}
	})
	if err != nil {
		out.Phase = "pre-commit"
		if syncOutcome.Result == statesync.ResultHalted || syncOutcome.Result == statesync.ResultFailed {
			out.Phase = "sync-pending"
		}
		out.Preparation = &result
		out.Sync = prepareReviewSyncJSON(branch, committed, syncOutcome)
		emitPrepareReviewJSON(stdout, out)
		fmt.Fprintln(stderr, "Error:", err)
		if syncOutcome.Result == statesync.ResultHalted {
			return 3
		}
		return 1
	}
	out.Preparation = &result
	out.Sync = prepareReviewSyncJSON(branch, committed, syncOutcome)
	checklist, err := loadPrepareReviewProjection(definitionDir, slug, "--checklist")
	if err != nil {
		out.Phase = "projection-pending"
		emitPrepareReviewJSON(stdout, out)
		fmt.Fprintln(stderr, "Error:", err)
		return 1
	}
	acceptance, err := loadPrepareReviewProjection(definitionDir, slug, "--ac-scan")
	if err != nil {
		out.Phase = "projection-pending"
		emitPrepareReviewJSON(stdout, out)
		fmt.Fprintln(stderr, "Error:", err)
		return 1
	}
	out.Phase = "complete"
	out.Checklist = checklist
	out.Acceptance = acceptance
	out.Presentation = map[string]any{
		"room": result.Room, "stage": evidence.Stage["name"], "stage_prose": evidence.Stage["bytes"],
		"question": parsed.question, "artifact": parsed.artifact, "references": parsed.references,
		"summary": parsed.summary, "recommendation": parsed.recommendation,
	}
	return emitPrepareReviewJSON(stdout, out)
}
func parsePrepareReviewArgs(args []string) (prepareReviewArgs, error) {
	var out prepareReviewArgs
	counts := map[string]int{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--publish":
			out.publish = true
		case "--json":
			out.json = true
		case "--workflow-dir", "--question", "--artifact", "--reference", "--summary", "--recommendation":
			if i+1 >= len(args) {
				return out, fmt.Errorf("%s requires an argument", args[i])
			}
			flag, value := args[i], args[i+1]
			i++
			counts[flag]++
			switch flag {
			case "--workflow-dir":
				out.workflowDir = value
			case "--question":
				out.question = value
			case "--artifact":
				out.artifact = filepath.ToSlash(filepath.Clean(value))
			case "--reference":
				out.references = append(out.references, filepath.ToSlash(filepath.Clean(value)))
			case "--summary":
				out.summary = value
			case "--recommendation":
				out.recommendation = value
			}
		default:
			return out, fmt.Errorf("unknown gate prepare-review flag %q", args[i])
		}
	}
	for _, flag := range []string{"--workflow-dir", "--question", "--artifact", "--summary", "--recommendation"} {
		if counts[flag] > 1 {
			return out, fmt.Errorf("gate prepare-review accepts %s exactly once", flag)
		}
	}
	if !out.publish && (out.question != "" || out.artifact != "" || len(out.references) != 0 || out.summary != "" || out.recommendation != "") {
		return out, fmt.Errorf("inspect mode accepts only --workflow-dir and --json")
	}
	return out, nil
}
func validatePrepareReviewInputs(args prepareReviewArgs, candidates []map[string]string) error {
	for name, value := range map[string]string{"--question": args.question, "--artifact": args.artifact, "--summary": args.summary, "--recommendation": args.recommendation} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("gate prepare-review --publish requires nonblank %s", name)
		}
		if !utf8.ValidString(value) {
			return fmt.Errorf("%s must be valid UTF-8", name)
		}
	}
	allowed := map[string]bool{}
	for _, candidate := range candidates {
		allowed[candidate["input_path"]] = true
	}
	seen := map[string]bool{}
	for _, selected := range append([]string{args.artifact}, args.references...) {
		if !allowed[selected] {
			return fmt.Errorf("selected path %q is not an inspect candidate", selected)
		}
		if seen[selected] {
			return fmt.Errorf("selected inspect candidate %q appears more than once", selected)
		}
		seen[selected] = true
	}
	return nil
}
func loadPrepareReviewEvidence(definitionDir, slug string) (prepareReviewEvidence, error) {
	raw, err := runPrepareReviewStatus(definitionDir, "--read", slug, "--gate-evidence", "--json")
	if err != nil {
		return prepareReviewEvidence{}, err
	}
	var evidence prepareReviewEvidence
	if err := json.Unmarshal(raw, &evidence); err != nil {
		return evidence, fmt.Errorf("decode gate evidence: %w", err)
	}
	return evidence, nil
}
func loadPrepareReviewProjection(definitionDir, slug, mode string) (json.RawMessage, error) {
	return runPrepareReviewStatus(definitionDir, "--read", slug, mode, "--json")
}
func runPrepareReviewStatus(definitionDir string, args ...string) (json.RawMessage, error) {
	var stdout, stderr bytes.Buffer
	argv := append([]string{"--workflow-dir", definitionDir}, args...)
	code, err := (&status.NativeRunner{}).Run(context.Background(), status.Request{Args: argv, Dir: definitionDir, Env: os.Environ(), Stdout: &stdout, Stderr: &stderr})
	if err != nil {
		return nil, err
	}
	if code != 0 {
		return nil, fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
	}
	return append(json.RawMessage(nil), stdout.Bytes()...), nil
}
func requireCleanPrepareReviewUnit(checkout string, paths []string) error {
	args := []string{"status", "--porcelain=v1", "--untracked-files=all", "--"}
	for _, path := range paths {
		args = append(args, literalGitPathspec(relToCheckout(checkout, path)))
	}
	ok, detail := runGit(checkout, args...)
	if !ok {
		return fmt.Errorf("inspect entity commit unit: %s", strings.TrimSpace(detail))
	}
	if strings.TrimSpace(detail) != "" {
		return fmt.Errorf("entity commit unit is dirty; commit or restore it before publish:\n%s", strings.TrimSpace(detail))
	}
	return nil
}
func restorePrepareReviewIndex(checkout string, paths []string) {
	args := []string{"reset", "--mixed", "HEAD", "--"}
	for _, path := range paths {
		args = append(args, literalGitPathspec(relToCheckout(checkout, path)))
	}
	runGit(checkout, args...)
}
func prepareReviewArgv(slug, definitionDir string) []string {
	return []string{"spacedock", "gate", "prepare-review", slug, "--workflow-dir", definitionDir, "--publish", "--question", "<TEXT>", "--artifact", "<INPUT_PATH>", "--summary", "<TEXT>", "--recommendation", "<TEXT>", "--reference", "<INPUT_PATH>", "--json"}
}
func prepareReviewSyncJSON(branch string, committed bool, outcome statesync.Outcome) map[string]any {
	return map[string]any{"result": outcome.Result, "state_branch": branch, "committed": committed, "integrated_peers": outcome.IntegratedPeers, "published_local": outcome.PublishedLocal, "peer_commit": outcome.PeerCommit, "conflicting_paths": outcome.ConflictingPaths}
}
func emitPrepareReviewJSON(w io.Writer, out prepareReviewOutput) int {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return 1
	}
	return 0
}
