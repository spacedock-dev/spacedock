package release

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

func TestClaudeLiveInstallsExactCandidateOutsideCheckout(t *testing.T) {
	workflow := readWorkflow(t, "runtime-live-e2e.yml")
	script := extractStepRun(t, workflow, "Install spacedock candidate")
	checkout, revision := candidateFixture(t)
	runnerTemp := t.TempDir()
	artifactDir := filepath.Join(t.TempDir(), "artifacts")
	githubEnv := filepath.Join(t.TempDir(), "github-env")
	githubPath := filepath.Join(t.TempDir(), "github-path")

	output, err := runCandidateInstall(t, script, checkout, runnerTemp, artifactDir, githubEnv, githubPath, "")
	if err != nil {
		t.Fatalf("candidate install failed: %v\n%s", err, output)
	}

	candidate := filepath.Join(runnerTemp, "spacedock-live-bin", executableName("spacedock"))
	info, err := os.Lstat(candidate)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() || info.Mode()&0o111 == 0 {
		t.Fatalf("candidate is not a physical executable: mode=%s", info.Mode())
	}
	realCandidate, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		t.Fatal(err)
	}
	realCheckout, err := filepath.EvalSymlinks(checkout)
	if err != nil {
		t.Fatal(err)
	}
	if realCandidate == realCheckout || strings.HasPrefix(realCandidate, realCheckout+string(os.PathSeparator)) {
		t.Fatalf("candidate %q resolves inside checkout %q", realCandidate, realCheckout)
	}

	envBytes, err := os.ReadFile(githubEnv)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(envBytes), "SPACEDOCK_BIN="+candidate+"\n"; got != want {
		t.Fatalf("GITHUB_ENV = %q, want %q", got, want)
	}
	pathBytes, err := os.ReadFile(githubPath)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(pathBytes), filepath.Dir(candidate)+"\n"; got != want {
		t.Fatalf("GITHUB_PATH = %q, want %q", got, want)
	}

	provenanceBytes, err := os.ReadFile(filepath.Join(artifactDir, "candidate-binary-provenance.txt"))
	if err != nil {
		t.Fatal(err)
	}
	provenance := string(provenanceBytes)
	for _, want := range []string{
		"candidate_path=" + candidate,
		"candidate_canonical_path=" + realCandidate,
		"checkout_revision=" + revision,
		"embedded_revision=" + revision,
		"vcs_modified=false",
	} {
		if !strings.Contains(provenance, want+"\n") {
			t.Errorf("provenance missing %q:\n%s", want, provenance)
		}
	}

	claudeJob := workflow[strings.Index(workflow, "  claude-live:"):strings.Index(workflow, "  codex-live:")]
	if strings.Contains(claudeJob, "            ./spacedock\n") {
		t.Error("Claude artifact upload still retains the obsolete checkout-root binary")
	}
}

func TestClaudeCandidateInstallFailsClosed(t *testing.T) {
	workflow := readWorkflow(t, "runtime-live-e2e.yml")
	script := extractStepRun(t, workflow, "Install spacedock candidate")

	t.Run("dangling symlink", func(t *testing.T) {
		checkout, _ := candidateFixture(t)
		runnerTemp := t.TempDir()
		candidateDir := filepath.Join(runnerTemp, "spacedock-live-bin")
		if err := os.MkdirAll(candidateDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(filepath.Join(runnerTemp, "missing"), filepath.Join(candidateDir, executableName("spacedock"))); err != nil {
			t.Fatal(err)
		}
		output, err := runCandidateInstall(t, script, checkout, runnerTemp, filepath.Join(t.TempDir(), "artifacts"), filepath.Join(t.TempDir(), "env"), filepath.Join(t.TempDir(), "path"), "")
		if err == nil || !strings.Contains(output, "candidate output already exists") {
			t.Fatalf("dangling symlink was not rejected: err=%v\n%s", err, output)
		}
	})

	t.Run("inside checkout", func(t *testing.T) {
		checkout, _ := candidateFixture(t)
		runnerTemp := filepath.Join(checkout, "runner-temp")
		if err := os.Mkdir(runnerTemp, 0o755); err != nil {
			t.Fatal(err)
		}
		output, err := runCandidateInstall(t, script, checkout, runnerTemp, filepath.Join(t.TempDir(), "artifacts"), filepath.Join(t.TempDir(), "env"), filepath.Join(t.TempDir(), "path"), "")
		if err == nil || !strings.Contains(output, "candidate resolves inside checkout") {
			t.Fatalf("inside-checkout candidate was not rejected: err=%v\n%s", err, output)
		}
	})

	t.Run("modified checkout", func(t *testing.T) {
		checkout, _ := candidateFixture(t)
		mainFile := filepath.Join(checkout, "cmd", "spacedock", "main.go")
		if err := os.WriteFile(mainFile, []byte("package main\nfunc main() { println(\"modified\") }\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		output, err := runCandidateInstall(t, script, checkout, t.TempDir(), filepath.Join(t.TempDir(), "artifacts"), filepath.Join(t.TempDir(), "env"), filepath.Join(t.TempDir(), "path"), "")
		if err == nil || !strings.Contains(output, "candidate build reports vcs.modified=true") {
			t.Fatalf("modified checkout was not rejected: err=%v\n%s", err, output)
		}
	})

	t.Run("revision mismatch", func(t *testing.T) {
		checkout, _ := candidateFixture(t)
		realGit, err := exec.LookPath("git")
		if err != nil {
			t.Fatal(err)
		}
		fakeBin := t.TempDir()
		wrapper := fmt.Sprintf("#!/bin/sh\nif [ \"$*\" = \"rev-parse HEAD\" ]; then\n  printf '0000000000000000000000000000000000000000\\n'\nelse\n  exec %q \"$@\"\nfi\n", realGit)
		if err := os.WriteFile(filepath.Join(fakeBin, executableName("git")), []byte(wrapper), 0o755); err != nil {
			t.Fatal(err)
		}
		output, err := runCandidateInstall(t, script, checkout, t.TempDir(), filepath.Join(t.TempDir(), "artifacts"), filepath.Join(t.TempDir(), "env"), filepath.Join(t.TempDir(), "path"), fakeBin)
		if err == nil || !strings.Contains(output, "embedded revision does not match checkout") {
			t.Fatalf("revision mismatch was not rejected: err=%v\n%s", err, output)
		}
	})
}

func candidateFixture(t *testing.T) (string, string) {
	t.Helper()
	checkout := t.TempDir()
	mainDir := filepath.Join(checkout, "cmd", "spacedock")
	if err := os.MkdirAll(mainDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(checkout, "go.mod"), []byte("module example.com/candidate\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(mainDir, "main.go"), []byte("package main\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = checkout
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.com",
			"GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@example.com",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	testgit.InitRepo(t, checkout, "-q")
	git("add", ".")
	git("commit", "-qm", "fixture")
	return checkout, git("rev-parse", "HEAD")
}

func runCandidateInstall(t *testing.T, script, checkout, runnerTemp, artifactDir, githubEnv, githubPath, pathPrefix string) (string, error) {
	t.Helper()
	cmd := exec.Command("bash", "-eu", "-o", "pipefail", "-c", script)
	cmd.Dir = checkout
	path := os.Getenv("PATH")
	if pathPrefix != "" {
		path = pathPrefix + string(os.PathListSeparator) + path
	}
	cmd.Env = append(os.Environ(),
		"PATH="+path,
		"RUNNER_TEMP="+runnerTemp,
		"GITHUB_WORKSPACE="+checkout,
		"SPACEDOCK_LIVE_ARTIFACT_DIR="+artifactDir,
		"GITHUB_ENV="+githubEnv,
		"GITHUB_PATH="+githubPath,
	)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func executableName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}
