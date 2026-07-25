// ABOUTME: Shared split-root state publication and rebase-conflict preflight.
// ABOUTME: Owns push/rebase/re-push, no-origin, and abort-with-evidence outcomes.
package statesync

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Result string

const (
	ResultReady     Result = "ready"
	ResultNoOp      Result = "no-op"
	ResultPushed    Result = "pushed"
	ResultLocalOnly Result = "local-only"
	ResultHalted    Result = "halted"
	ResultFailed    Result = "failed"
)

type Outcome struct {
	Result           Result
	IntegratedPeers  bool
	PublishedLocal   bool
	ConflictingPaths []string
	PeerCommit       string
	Detail           string
}

func Preflight(checkout, branch string) Outcome {
	if ok, detail := validBranch(branch); !ok {
		return Outcome{Result: ResultFailed, Detail: detail}
	}
	if ok, detail := exactCheckoutRoot(checkout); !ok {
		return Outcome{Result: ResultFailed, Detail: detail}
	}
	if rebaseInProgress(checkout) {
		rebaseBranch, ok := rebaseHeadBranch(checkout)
		if !ok || rebaseBranch != branch {
			return Outcome{Result: ResultFailed, Detail: "rebase already in progress for a different or unknown branch; left untouched"}
		}
		outcome := Outcome{
			Result:           ResultHalted,
			ConflictingPaths: conflictingPaths(checkout),
			PeerCommit:       peerCommit(checkout, branch),
		}
		if ok, detail := runGit(checkout, "rebase", "--abort"); !ok {
			outcome.Result = ResultFailed
			outcome.Detail = "failed to abort state rebase; checkout may remain conflicted: " + strings.TrimSpace(detail)
		}
		return outcome
	}
	ok, head := runGit(checkout, "symbolic-ref", "--quiet", "--short", "HEAD")
	actual := strings.TrimSpace(head)
	if !ok {
		actual = "detached HEAD"
	}
	if actual != branch {
		return Outcome{Result: ResultFailed, Detail: "state checkout HEAD is " + actual + "; expected branch " + branch + "; switch to it or rename the intended state branch with `git branch -M " + branch + "`"}
	}
	return Outcome{Result: ResultReady}
}

func Publish(checkout, branch string) Outcome {
	if outcome := Preflight(checkout, branch); outcome.Result != ResultReady {
		return outcome
	}
	if ok, _ := runGit(checkout, "remote", "get-url", "origin"); !ok {
		return Outcome{Result: ResultLocalOnly}
	}
	_, localHead := runGit(checkout, "rev-parse", "HEAD")
	localHead = strings.TrimSpace(localHead)
	ref := "refs/heads/" + branch
	remoteHead := remoteRefHead(checkout, ref)
	alreadyPublished := remoteHead != "" && remoteHead == localHead
	if ok, _ := runGit(checkout, "push", "origin", "HEAD:"+ref); ok {
		if alreadyPublished {
			return Outcome{Result: ResultNoOp}
		}
		return Outcome{Result: ResultPushed, PublishedLocal: true}
	}
	ok, out := runGit(checkout, "pull", "--rebase", "origin", ref)
	if !ok {
		if rebaseInProgress(checkout) {
			halted := Preflight(checkout, branch)
			if halted.Result == ResultFailed {
				halted.Detail += "\npull --rebase also reported: " + strings.TrimSpace(out)
			} else {
				halted.Detail = out
			}
			return halted
		}
		return Outcome{Result: ResultFailed, Detail: out}
	}
	if ok, out = runGit(checkout, "push", "origin", "HEAD:"+ref); !ok {
		return Outcome{Result: ResultFailed, IntegratedPeers: true, Detail: out}
	}
	publishedLocal := true
	if remoteHead != "" {
		if oldHeadWasPeerAncestor, _ := runGit(checkout, "merge-base", "--is-ancestor", localHead, remoteHead); oldHeadWasPeerAncestor {
			publishedLocal = false
		}
	}
	return Outcome{Result: ResultPushed, IntegratedPeers: true, PublishedLocal: publishedLocal}
}

func Pull(checkout, branch string) Outcome {
	if outcome := Preflight(checkout, branch); outcome.Result != ResultReady {
		return outcome
	}
	if ok, _ := runGit(checkout, "remote", "get-url", "origin"); !ok {
		return Outcome{Result: ResultLocalOnly}
	}
	ok, out := runGit(checkout, "pull", "--rebase", "origin", "refs/heads/"+branch)
	if ok {
		return Outcome{Result: ResultReady, IntegratedPeers: true}
	}
	if rebaseInProgress(checkout) {
		halted := Preflight(checkout, branch)
		if halted.Result == ResultFailed {
			halted.Detail += "\npull --rebase also reported: " + strings.TrimSpace(out)
		} else {
			halted.Detail = out
		}
		return halted
	}
	return Outcome{Result: ResultFailed, Detail: out}
}

func runGit(checkout string, args ...string) (bool, string) {
	cmd := exec.Command("git", append([]string{"-C", checkout}, args...)...)
	out, err := cmd.CombinedOutput()
	return err == nil, string(out)
}

func validBranch(branch string) (bool, string) {
	if branch == "" || strings.HasPrefix(branch, "-") {
		return false, "invalid state branch " + branch
	}
	cmd := exec.Command("git", "check-ref-format", "refs/heads/"+branch)
	if out, err := cmd.CombinedOutput(); err != nil {
		return false, "state branch validation failed for " + branch + ": " + strings.TrimSpace(string(out))
	}
	return true, ""
}

func exactCheckoutRoot(checkout string) (bool, string) {
	ok, out := runGit(checkout, "rev-parse", "--show-toplevel")
	if !ok {
		return false, "state checkout is not a Git worktree: " + strings.TrimSpace(out)
	}
	want, err := filepath.EvalSymlinks(checkout)
	if err != nil {
		return false, "cannot resolve state checkout: " + err.Error()
	}
	got, err := filepath.EvalSymlinks(strings.TrimSpace(out))
	if err != nil || filepath.Clean(got) != filepath.Clean(want) {
		return false, "state checkout must be the exact Git toplevel: configured " + want + ", observed " + strings.TrimSpace(out) + "; point the workflow state: field at its dedicated linked worktree"
	}
	return true, ""
}

func remoteRefHead(checkout, ref string) string {
	ok, out := runGit(checkout, "ls-remote", "--refs", "origin", ref)
	if !ok {
		return ""
	}
	fields := strings.Fields(out)
	if len(fields) < 2 || fields[1] != ref {
		return ""
	}
	return fields[0]
}

func rebaseInProgress(checkout string) bool {
	for _, name := range []string{"rebase-merge", "rebase-apply"} {
		ok, out := runGit(checkout, "rev-parse", "--git-path", name)
		if !ok {
			continue
		}
		path := strings.TrimSpace(out)
		if !filepath.IsAbs(path) {
			path = filepath.Join(checkout, path)
		}
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	return false
}

func rebaseHeadBranch(checkout string) (string, bool) {
	for _, kind := range []string{"rebase-merge", "rebase-apply"} {
		ok, out := runGit(checkout, "rev-parse", "--git-path", filepath.Join(kind, "head-name"))
		if !ok {
			continue
		}
		path := strings.TrimSpace(out)
		if !filepath.IsAbs(path) {
			path = filepath.Join(checkout, path)
		}
		body, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if branch := strings.TrimPrefix(strings.TrimSpace(string(body)), "refs/heads/"); branch != "" {
			return branch, true
		}
	}
	return "", false
}

func conflictingPaths(checkout string) []string {
	ok, out := runGit(checkout, "diff", "--name-only", "--diff-filter=U")
	if !ok {
		return nil
	}
	var paths []string
	for _, path := range strings.Split(strings.TrimSpace(out), "\n") {
		if path = strings.TrimSpace(path); path != "" {
			paths = append(paths, path)
		}
	}
	return paths
}

func peerCommit(checkout, branch string) string {
	ok, out := runGit(checkout, "rev-parse", "--short", "origin/"+branch)
	if !ok {
		return ""
	}
	return strings.TrimSpace(out)
}
