// ABOUTME: Durable under-lock outcome for concurrent split-root resume waiters.
package cli

import (
	"crypto/rand"
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
	Path       string `json:"path"`
	Result     string `json:"result"`
	Generation string `json:"generation,omitempty"`
}

const stateResumeGenerationFile = "spacedock-state-resume.generation"

func unsupportedStateResumeLock(_ string, _ func() int) (int, error) {
	return 0, errStateResumeLockUnsupported
}

func stateResumePathKey(statePath string) string {
	key := sha256.Sum256([]byte(status.RealpathOf(statePath)))
	return fmt.Sprintf("%x", key)
}

func stateResumeOutcomePath(workflowDir, statePath string) (string, error) {
	placement, err := status.ResolveRepositoryPlacement(workflowDir)
	if err != nil {
		return "", err
	}
	if !placement.InGit {
		return "", fmt.Errorf("not a git repository at %s", workflowDir)
	}
	return filepath.Join(placement.CommonGitDir, "spacedock-state-resume."+stateResumePathKey(statePath)+".outcome.json"), nil
}

func writeStateResumeOutcome(workflowDir, statePath, result string) error {
	path, err := stateResumeOutcomePath(workflowDir, statePath)
	if err != nil {
		return err
	}
	outcome := stateResumeOutcome{Path: status.RealpathOf(statePath), Result: result}
	if result == "ready" {
		generationPath, err := stateResumeGenerationPath(statePath)
		if err != nil {
			return err
		}
		generation := make([]byte, 32)
		if _, err := rand.Read(generation); err != nil {
			return fmt.Errorf("generate state checkout identity: %w", err)
		}
		outcome.Generation = fmt.Sprintf("%x", generation)
		if err := os.WriteFile(generationPath, []byte(outcome.Generation), 0o600); err != nil {
			return fmt.Errorf("record state checkout identity: %w", err)
		}
	}
	body, err := json.Marshal(outcome)
	if err != nil {
		return err
	}
	return writeStateResumeOutcomeFile(path, body)
}

func readStateResumeOutcome(workflowDir, statePath string) (string, error) {
	outcome, err := readStateResumeOutcomeRecord(workflowDir, statePath)
	if err != nil {
		return "", err
	}
	return outcome.Result, nil
}

func readStateResumeOutcomeRecord(workflowDir, statePath string) (stateResumeOutcome, error) {
	path, err := stateResumeOutcomePath(workflowDir, statePath)
	if err != nil {
		return stateResumeOutcome{}, err
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return stateResumeOutcome{}, err
	}
	var outcome stateResumeOutcome
	if err := json.Unmarshal(body, &outcome); err != nil {
		return stateResumeOutcome{}, err
	}
	if outcome.Path != status.RealpathOf(statePath) {
		return stateResumeOutcome{}, fmt.Errorf("resume outcome belongs to a different state checkout")
	}
	if outcome.Result == "ready" {
		if outcome.Generation == "" {
			return stateResumeOutcome{}, fmt.Errorf("ready resume outcome has no checkout identity")
		}
		generationPath, err := stateResumeGenerationPath(statePath)
		if err != nil {
			return stateResumeOutcome{}, err
		}
		generation, err := os.ReadFile(generationPath)
		if err != nil {
			return stateResumeOutcome{}, fmt.Errorf("read state checkout identity: %w", err)
		}
		if string(generation) != outcome.Generation {
			return stateResumeOutcome{}, fmt.Errorf("resume outcome belongs to a previous state checkout generation")
		}
	}
	return outcome, nil
}

func stateResumeGenerationPath(statePath string) (string, error) {
	ok, out := runGit(statePath, "rev-parse", "--path-format=absolute", "--git-dir")
	if !ok {
		return "", fmt.Errorf("resolve state checkout identity directory:\n%s", out)
	}
	gitDir := status.TrimGitLineTerminator(out)
	if !filepath.IsAbs(gitDir) {
		return "", fmt.Errorf("state checkout identity directory is not absolute: %s", gitDir)
	}
	return filepath.Join(gitDir, stateResumeGenerationFile), nil
}

func writeStateResumeOutcomeFile(path string, body []byte) (err error) {
	temp, err := os.CreateTemp(filepath.Dir(path), ".spacedock-state-resume-outcome-")
	if err != nil {
		return err
	}
	tempPath := temp.Name()
	defer func() {
		temp.Close()
		if err != nil {
			os.Remove(tempPath)
		}
	}()
	if err = temp.Chmod(0o600); err != nil {
		return err
	}
	if _, err = temp.Write(body); err != nil {
		return err
	}
	if err = temp.Close(); err != nil {
		return err
	}
	return os.Rename(tempPath, path)
}
