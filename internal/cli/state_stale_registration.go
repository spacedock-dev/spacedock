// ABOUTME: Exact stale state-worktree administrative record cleanup.
// ABOUTME: Claims obsolete metadata without deleting a recreated public checkout.
package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spacedock-dev/spacedock/internal/status"
)

type staleStateRegistrationAdmin struct {
	path      string
	info      os.FileInfo
	gitdir    []byte
	head      []byte
	commondir []byte
}

func findStaleStateRegistrationAdmin(workflowDir, statePath, branch string) (staleStateRegistrationAdmin, string, error) {
	placement, err := status.ResolveRepositoryPlacement(workflowDir)
	if err != nil {
		return staleStateRegistrationAdmin{}, "", err
	}
	if !placement.InGit {
		return staleStateRegistrationAdmin{}, "", fmt.Errorf("state checkout %s has no repository administration directory", statePath)
	}
	worktreesDir := filepath.Join(placement.CommonGitDir, "worktrees")
	entries, err := os.ReadDir(worktreesDir)
	if err != nil {
		return staleStateRegistrationAdmin{}, "", fmt.Errorf("reading state worktree administration directory: %w", err)
	}
	target := status.RealpathOf(statePath)
	wantHead := "ref: refs/heads/" + branch
	var matches []staleStateRegistrationAdmin
	for _, entry := range entries {
		if entry.Type()&os.ModeSymlink != 0 || !entry.IsDir() {
			continue
		}
		adminPath := filepath.Join(worktreesDir, entry.Name())
		gitdir, err := os.ReadFile(filepath.Join(adminPath, "gitdir"))
		if err != nil {
			continue
		}
		gitdirPath := status.TrimGitLineTerminator(string(gitdir))
		if !filepath.IsAbs(gitdirPath) || filepath.Base(gitdirPath) != ".git" || status.RealpathOf(filepath.Dir(gitdirPath)) != target {
			continue
		}
		head, err := os.ReadFile(filepath.Join(adminPath, "HEAD"))
		if err != nil {
			return staleStateRegistrationAdmin{}, "", fmt.Errorf("reading stale state registration HEAD: %w", err)
		}
		if status.TrimGitLineTerminator(string(head)) != wantHead {
			return staleStateRegistrationAdmin{}, "", fmt.Errorf("stale state registration administrative HEAD is not %q", wantHead)
		}
		commondir, err := os.ReadFile(filepath.Join(adminPath, "commondir"))
		if err != nil {
			return staleStateRegistrationAdmin{}, "", fmt.Errorf("reading stale state registration common directory: %w", err)
		}
		commonPath := status.TrimGitLineTerminator(string(commondir))
		if !filepath.IsAbs(commonPath) {
			commonPath = filepath.Join(adminPath, commonPath)
		}
		if status.RealpathOf(commonPath) != status.RealpathOf(placement.CommonGitDir) {
			return staleStateRegistrationAdmin{}, "", fmt.Errorf("stale state registration belongs to a different Git repository")
		}
		info, err := os.Stat(adminPath)
		if err != nil {
			return staleStateRegistrationAdmin{}, "", fmt.Errorf("stating stale state registration administration directory: %w", err)
		}
		matches = append(matches, staleStateRegistrationAdmin{
			path: adminPath, info: info, gitdir: gitdir, head: head, commondir: commondir,
		})
	}
	if len(matches) != 1 {
		return staleStateRegistrationAdmin{}, "", fmt.Errorf("state checkout %s has %d matching administrative registrations, expected 1", statePath, len(matches))
	}
	return matches[0], placement.CommonGitDir, nil
}

func (admin staleStateRegistrationAdmin) matchesAt(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if !os.SameFile(admin.info, info) {
		return false, nil
	}
	for name, want := range map[string][]byte{
		"gitdir": admin.gitdir, "HEAD": admin.head, "commondir": admin.commondir,
	} {
		got, err := os.ReadFile(filepath.Join(path, name))
		if err != nil {
			return false, err
		}
		if !bytes.Equal(got, want) {
			return false, nil
		}
	}
	return true, nil
}

func restoreClaimedRegistration(quarantine string, admin staleStateRegistrationAdmin) error {
	if err := renameNoReplace(quarantine, admin.path); err != nil {
		return fmt.Errorf("restoring concurrently replaced state registration: %w", err)
	}
	return nil
}

func removeStaleStateRegistration(workflowDir, statePath, branch string) error {
	if _, err := os.Lstat(statePath); err == nil {
		return fmt.Errorf("state checkout appeared concurrently at %s; left it untouched", statePath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("checking state checkout before registration repair %s: %w", statePath, err)
	}
	ok, out := runGit(workflowDir, "worktree", "list", "--porcelain", "-z")
	if !ok {
		return fmt.Errorf("listing worktree registrations before repair failed:\n%s", out)
	}
	records, err := status.ParseWorktreePorcelainZ([]byte(out))
	if err != nil {
		return fmt.Errorf("parsing worktree registrations before repair: %w", err)
	}
	stale, err := classifyStaleStateRegistration(records, statePath, branch)
	if err != nil {
		return err
	}
	if !stale {
		return fmt.Errorf("state checkout %s no longer has the expected prunable registration; refusing stale repair", statePath)
	}
	admin, commonGitDir, err := findStaleStateRegistrationAdmin(workflowDir, statePath, branch)
	if err != nil {
		return err
	}
	if stateResumeAfterStaleRegistrationVerificationHook != nil {
		stateResumeAfterStaleRegistrationVerificationHook(statePath)
	}

	quarantine, err := os.MkdirTemp(commonGitDir, ".spacedock-stale-registration-")
	if err != nil {
		return fmt.Errorf("allocating stale state registration quarantine: %w", err)
	}
	if err := os.Remove(quarantine); err != nil {
		return fmt.Errorf("preparing stale state registration quarantine: %w", err)
	}
	if err := renameNoReplace(admin.path, quarantine); err != nil {
		return fmt.Errorf("claiming exact stale state registration: %w", err)
	}
	matches, matchErr := admin.matchesAt(quarantine)
	if matchErr != nil || !matches {
		if restoreErr := restoreClaimedRegistration(quarantine, admin); restoreErr != nil {
			return fmt.Errorf("state checkout %s registration changed concurrently and could not be restored: %v (identity check: %v)", statePath, restoreErr, matchErr)
		}
		return fmt.Errorf("state checkout %s registration changed concurrently; restored it and left it untouched", statePath)
	}
	if err := os.RemoveAll(quarantine); err != nil {
		return fmt.Errorf("removing exact stale state administrative registration: %w", err)
	}
	return nil
}
