// ABOUTME: AC-1..AC-4 oracles for the unsandboxed-launch permission posture:
// ABOUTME: claude --permission-mode auto / codex --ask-for-approval on-request injection.
package cli

import (
	"bytes"
	"context"
	"testing"
)

// argvHasFlagValue reports whether argv contains the space-form flag/value pair
// (a `flag` token immediately followed by `value`), and how many times `flag`
// appears in total. The count lets an oracle assert a single occurrence (operator
// override must not produce a duplicate).
func argvHasFlagValue(argv []string, flag, value string) (pair bool, count int) {
	for i, tok := range argv {
		if tok == flag {
			count++
			if i+1 < len(argv) && argv[i+1] == value {
				pair = true
			}
		}
	}
	return pair, count
}

// AC-1: an unsandboxed `spacedock claude` launch injects `--permission-mode auto`
// and carries NO `--dangerously-skip-permissions`; the sandboxed launch is
// unchanged (`--dangerously-skip-permissions`, NO injected `--permission-mode`).
func TestClaudeUnsandboxedInjectsAutoPermissionMode(t *testing.T) {
	t.Run("unsandboxed-injects-auto", func(t *testing.T) {
		dir := t.TempDir() // no .safehouse
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runClaude(context.Background(), nil, dir, fake, lookFound, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		if pair, _ := argvHasFlagValue(fake.launchedArg, "--permission-mode", "auto"); !pair {
			t.Fatalf("unsandboxed claude launch missing --permission-mode auto: %v", fake.launchedArg)
		}
		for _, tok := range fake.launchedArg {
			if tok == "--dangerously-skip-permissions" {
				t.Fatalf("unsandboxed claude launch carried --dangerously-skip-permissions: %v", fake.launchedArg)
			}
		}
	})
	t.Run("sandboxed-unchanged", func(t *testing.T) {
		dir := safehouseFixtureDir(t)
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runClaude(context.Background(), nil, dir, fake, lookFound, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		sawSkip := false
		for _, tok := range fake.launchedArg {
			if tok == "--dangerously-skip-permissions" {
				sawSkip = true
			}
			if tok == "--permission-mode" {
				t.Fatalf("sandboxed claude launch carried injected --permission-mode: %v", fake.launchedArg)
			}
		}
		if !sawSkip {
			t.Fatalf("sandboxed claude launch missing --dangerously-skip-permissions: %v", fake.launchedArg)
		}
	})
}

// AC-2 (captain option A): an unsandboxed `spacedock codex` launch injects
// `--ask-for-approval on-request` and carries NO bypass flag; the sandboxed
// launch is unchanged (`--dangerously-bypass-approvals-and-sandbox`, NO injected
// approval flag).
func TestCodexUnsandboxedInjectsOnRequestApproval(t *testing.T) {
	t.Run("unsandboxed-injects-on-request", func(t *testing.T) {
		dir := t.TempDir() // no .safehouse
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runCodex(context.Background(), nil, dir, fake, lookFound, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		if pair, _ := argvHasFlagValue(fake.launchedArg, "--ask-for-approval", "on-request"); !pair {
			t.Fatalf("unsandboxed codex launch missing --ask-for-approval on-request: %v", fake.launchedArg)
		}
		for _, tok := range fake.launchedArg {
			if tok == "--dangerously-bypass-approvals-and-sandbox" {
				t.Fatalf("unsandboxed codex launch carried the bypass flag: %v", fake.launchedArg)
			}
		}
	})
	t.Run("sandboxed-unchanged", func(t *testing.T) {
		dir := safehouseFixtureDir(t)
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runCodex(context.Background(), nil, dir, fake, lookFound, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		sawBypass := false
		for _, tok := range fake.launchedArg {
			if tok == "--dangerously-bypass-approvals-and-sandbox" {
				sawBypass = true
			}
			if tok == "--ask-for-approval" {
				t.Fatalf("sandboxed codex launch carried injected --ask-for-approval: %v", fake.launchedArg)
			}
		}
		if !sawBypass {
			t.Fatalf("sandboxed codex launch missing --dangerously-bypass-approvals-and-sandbox: %v", fake.launchedArg)
		}
	})
}

// AC-3: an operator-supplied permission/approval flag in the passthrough
// suppresses the spacedock-injected one — exactly one occurrence, operator value
// wins (no duplicate). Covered for both space form and equals form.
func TestOperatorPermissionFlagSuppressesInjection(t *testing.T) {
	t.Run("claude-space-form", func(t *testing.T) {
		dir := t.TempDir() // no .safehouse
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runClaude(context.Background(), []string{"--", "--permission-mode", "plan"}, dir, fake, lookFound, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		pair, count := argvHasFlagValue(fake.launchedArg, "--permission-mode", "plan")
		if !pair {
			t.Fatalf("operator --permission-mode plan not preserved: %v", fake.launchedArg)
		}
		if count != 1 {
			t.Fatalf("--permission-mode appears %d times, want 1 (no injected duplicate): %v", count, fake.launchedArg)
		}
	})
	t.Run("claude-equals-form", func(t *testing.T) {
		dir := t.TempDir() // no .safehouse
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runClaude(context.Background(), []string{"--", "--permission-mode=plan"}, dir, fake, lookFound, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		for _, tok := range fake.launchedArg {
			if tok == "--permission-mode" {
				t.Fatalf("injected space-form --permission-mode despite operator equals-form: %v", fake.launchedArg)
			}
		}
	})
	t.Run("codex-space-form", func(t *testing.T) {
		dir := t.TempDir() // no .safehouse
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runCodex(context.Background(), []string{"--", "--ask-for-approval", "untrusted"}, dir, fake, lookFound, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		pair, count := argvHasFlagValue(fake.launchedArg, "--ask-for-approval", "untrusted")
		if !pair {
			t.Fatalf("operator --ask-for-approval untrusted not preserved: %v", fake.launchedArg)
		}
		if count != 1 {
			t.Fatalf("--ask-for-approval appears %d times, want 1 (no injected duplicate): %v", count, fake.launchedArg)
		}
	})
	t.Run("codex-short-form", func(t *testing.T) {
		dir := t.TempDir() // no .safehouse
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runCodex(context.Background(), []string{"--", "-a", "untrusted"}, dir, fake, lookFound, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		for _, tok := range fake.launchedArg {
			if tok == "--ask-for-approval" {
				t.Fatalf("injected --ask-for-approval despite operator short-form -a: %v", fake.launchedArg)
			}
		}
	})
}

// AC-4: the injected flag rides the non-resume gate — a resumed unsandboxed
// launch is NOT forced into the auto/approval mode. Mirrors the resume-suppression
// oracle: the bootstrap prompt and the injected flag share the same gate.
func TestResumeUnsandboxedSuppressesInjection(t *testing.T) {
	t.Run("claude-resume", func(t *testing.T) {
		dir := t.TempDir() // no .safehouse
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runClaude(context.Background(), []string{"--", "--resume"}, dir, fake, lookFound, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		for _, tok := range fake.launchedArg {
			if tok == "--permission-mode" {
				t.Fatalf("resumed claude launch carried injected --permission-mode: %v", fake.launchedArg)
			}
		}
	})
	t.Run("codex-resume", func(t *testing.T) {
		dir := t.TempDir() // no .safehouse
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runCodex(context.Background(), []string{"--", "resume", "abc123"}, dir, fake, lookFound, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		for _, tok := range fake.launchedArg {
			if tok == "--ask-for-approval" {
				t.Fatalf("resumed codex launch carried injected --ask-for-approval: %v", fake.launchedArg)
			}
		}
	})
}
