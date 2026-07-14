// ABOUTME: Repository-scoped advisory lock for atomic split-root resume.
//go:build unix

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spacedock-dev/spacedock/internal/status"
)

func withStateResumeLock(workflowDir, statePath string, fn func(waitedForState bool) int) (int, error) {
	placement, err := status.ResolveRepositoryPlacement(workflowDir)
	if err != nil {
		return 0, err
	}
	if !placement.InGit {
		return 0, fmt.Errorf("not a git repository at %s", workflowDir)
	}
	stateLockPath := filepath.Join(placement.CommonGitDir, "spacedock-state-resume."+stateResumePathKey(statePath)+".lock")
	stateLock, err := os.OpenFile(stateLockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return 0, fmt.Errorf("open checkout resume lock: %w", err)
	}
	defer stateLock.Close()
	waitedForState := false
	if err := syscall.Flock(int(stateLock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		if err != syscall.EWOULDBLOCK && err != syscall.EAGAIN {
			return 0, fmt.Errorf("acquire checkout resume lock: %w", err)
		}
		waitedForState = true
		if stateResumeWaitHook != nil {
			stateResumeWaitHook(statePath)
		}
		if err := syscall.Flock(int(stateLock.Fd()), syscall.LOCK_EX); err != nil {
			return 0, fmt.Errorf("wait for checkout resume lock: %w", err)
		}
	}
	defer syscall.Flock(int(stateLock.Fd()), syscall.LOCK_UN)

	// Git worktree administration is repository-wide, so a second lock retains
	// the original cross-workflow serialization. The per-checkout lock above
	// separately proves whether this caller waited for work on its own state path.
	repositoryLockPath := filepath.Join(placement.CommonGitDir, "spacedock-state-resume.lock")
	repositoryLock, err := os.OpenFile(repositoryLockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return 0, fmt.Errorf("open repository resume lock: %w", err)
	}
	defer repositoryLock.Close()
	if err := syscall.Flock(int(repositoryLock.Fd()), syscall.LOCK_EX); err != nil {
		return 0, fmt.Errorf("acquire repository resume lock: %w", err)
	}
	defer syscall.Flock(int(repositoryLock.Fd()), syscall.LOCK_UN)
	return fn(waitedForState), nil
}
