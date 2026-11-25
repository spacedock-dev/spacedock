package cli

import (
	"fmt"
	"strings"
)

func (r brewRepair) cleanBrewTap(tap string) (head, branch string, err error) {
	remote, err := r.git(tap, "remote", "get-url", "--all", "origin")
	if err != nil || (strings.TrimSpace(remote) != spacedockTapRemote && strings.TrimSpace(remote) != spacedockTapRemote+".git") {
		return "", "", fmt.Errorf("the installed tap does not track the official Spacedock source. Verify its origin")
	}
	actualRoot, err := r.git(tap, "rev-parse", "--show-toplevel")
	if err != nil || strings.TrimSpace(actualRoot) != tap {
		return "", "", fmt.Errorf("Homebrew's Spacedock tap is not its own Git checkout")
	}
	dirty, err := r.git(tap, "status", "--porcelain", "--untracked-files=all")
	if err != nil || strings.TrimSpace(dirty) != "" {
		return "", "", fmt.Errorf("the Spacedock tap is dirty or unreadable. Finish the local tap work first")
	}
	head, err = r.git(tap, "rev-parse", "HEAD")
	if err != nil || !repairCommit.MatchString(strings.TrimSpace(head)) {
		return "", "", fmt.Errorf("the Spacedock tap has no readable Git commit")
	}
	branch, err = r.git(tap, "symbolic-ref", "--short", "HEAD")
	if err != nil {
		return "", "", fmt.Errorf("the Spacedock tap is detached. Switch to its tracked default branch")
	}
	branch = strings.TrimSpace(branch)
	upstream, upErr := r.git(tap, "rev-parse", "--symbolic-full-name", "@{upstream}")
	defaultRef, defaultErr := r.git(tap, "symbolic-ref", "refs/remotes/origin/HEAD")
	ref := "refs/remotes/origin/" + branch
	if upErr != nil || defaultErr != nil || strings.TrimSpace(upstream) != ref || strings.TrimSpace(defaultRef) != ref {
		return "", "", fmt.Errorf("the Spacedock tap is not on its tracked default branch")
	}
	if _, err := r.git(tap, "merge-base", "--is-ancestor", strings.TrimSpace(head), ref); err != nil {
		return "", "", fmt.Errorf("the Spacedock tap is ahead of its tracked branch or has diverged. Reconcile it first")
	}
	return strings.TrimSpace(head), branch, nil
}

func (r brewRepair) advanceBrewTap(tap, branch, head string) (string, error) {
	refs, err := r.git(tap, "ls-remote", "--symref", spacedockTapRemote, "HEAD")
	lines := strings.Split(strings.TrimSpace(refs), "\n")
	if err != nil || len(lines) != 2 || lines[0] != "ref: refs/heads/"+branch+"\tHEAD" {
		return "", fmt.Errorf("could not verify the official tap default branch")
	}
	candidate, ok := strings.CutSuffix(lines[1], "\tHEAD")
	if !ok || !repairCommit.MatchString(candidate) {
		return "", fmt.Errorf("could not identify the official Spacedock tap commit")
	}
	tracking := "refs/remotes/origin/" + branch
	if _, err := r.git(tap, "fetch", "--no-tags", "--no-recurse-submodules", "--no-write-fetch-head",
		spacedockTapRemote, "refs/heads/"+branch+":"+tracking); err != nil {
		return "", err
	}
	if _, err := r.git(tap, "merge-base", "--is-ancestor", candidate, tracking); err != nil {
		return "", fmt.Errorf("the official tap default branch moved during the fetch")
	}
	if _, err := r.git(tap, "merge-base", "--is-ancestor", head, candidate); err != nil {
		return "", fmt.Errorf("the Spacedock tap has diverged. Reconcile its tracked branch first")
	}
	now, currentBranch, err := r.cleanBrewTap(tap)
	if err != nil || now != head || currentBranch != branch {
		return "", fmt.Errorf("the Spacedock tap changed during preflight")
	}
	if _, err := r.git(tap, "merge", "--ff-only", "--no-edit", candidate); err != nil {
		return "", err
	}
	return candidate, nil
}
