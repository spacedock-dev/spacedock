// ABOUTME: Durable under-lock outcome for concurrent split-root resume waiters.
package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spacedock-dev/spacedock/internal/status"
)

var errStateResumeLockUnsupported = errors.New("atomic state resume is unsupported on this platform")

type stateResumeOutcome struct {
	Path   string `json:"path"`
	Result string `json:"result"`
}

func unsupportedStateResumeLock(_ string, _ func() int) (int, error) {
	return 0, errStateResumeLockUnsupported
}

func stateResumeOutcomePath(workflowDir string) (string, error) {
	placement, err := status.ResolveRepositoryPlacement(workflowDir)
	if err != nil {
		return "", err
	}
	if !placement.InGit {
		return "", fmt.Errorf("not a git repository at %s", workflowDir)
	}
	return filepath.Join(placement.CommonGitDir, "spacedock-state-resume.outcome.json"), nil
}

func writeStateResumeOutcome(workflowDir, statePath, result string) error {
	path, err := stateResumeOutcomePath(workflowDir)
	if err != nil {
		return err
	}
	body, err := json.Marshal(stateResumeOutcome{Path: status.RealpathOf(statePath), Result: result})
	if err != nil {
		return err
	}
	return os.WriteFile(path, body, 0o600)
}

func readStateResumeOutcome(workflowDir, statePath string) (string, error) {
	path, err := stateResumeOutcomePath(workflowDir)
	if err != nil {
		return "", err
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var outcome stateResumeOutcome
	if err := json.Unmarshal(body, &outcome); err != nil {
		return "", err
	}
	if outcome.Path != status.RealpathOf(statePath) {
		return "", fmt.Errorf("resume outcome belongs to a different state checkout")
	}
	return outcome.Result, nil
}

func cleanupFailedStateResume(workflowDir, statePath string) error {
	if dirExists(statePath) {
		runGit(workflowDir, "worktree", "remove", "--force", statePath)
	}
	if dirExists(statePath) {
		if err := os.RemoveAll(statePath); err != nil {
			return err
		}
	}
	if err := repairStaleWorktreeRegistration(workflowDir, statePath, io.Discard); err != nil {
		return err
	}
	if dirExists(statePath) {
		return fmt.Errorf("state checkout still exists at %s", statePath)
	}
	return nil
}
