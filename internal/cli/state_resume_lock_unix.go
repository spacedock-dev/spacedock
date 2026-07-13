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

func withStateResumeLock(workflowDir string, fn func() int) (int, error) {
	placement, err := status.ResolveRepositoryPlacement(workflowDir)
	if err != nil {
		return 0, err
	}
	if !placement.InGit {
		return 0, fmt.Errorf("not a git repository at %s", workflowDir)
	}
	lockPath := filepath.Join(placement.CommonGitDir, "spacedock-state-resume.lock")
	lock, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return 0, fmt.Errorf("open resume lock: %w", err)
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return 0, fmt.Errorf("acquire resume lock: %w", err)
	}
	defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
	return fn(), nil
}
