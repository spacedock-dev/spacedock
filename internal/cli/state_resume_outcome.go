// ABOUTME: Durable under-lock outcome for concurrent split-root resume waiters.
package cli

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
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

func stateResumeOutcomePath(workflowDir, statePath string) (string, error) {
	placement, err := status.ResolveRepositoryPlacement(workflowDir)
	if err != nil {
		return "", err
	}
	if !placement.InGit {
		return "", fmt.Errorf("not a git repository at %s", workflowDir)
	}
	key := sha256.Sum256([]byte(status.RealpathOf(statePath)))
	return filepath.Join(placement.CommonGitDir, fmt.Sprintf("spacedock-state-resume.%x.outcome.json", key)), nil
}

func writeStateResumeOutcome(workflowDir, statePath, result string) error {
	path, err := stateResumeOutcomePath(workflowDir, statePath)
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
	path, err := stateResumeOutcomePath(workflowDir, statePath)
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
