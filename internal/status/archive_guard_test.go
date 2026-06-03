// ABOUTME: AC-3 --archive dest-spelling parity (relative vs absolute --workflow-dir)
// ABOUTME: and the terminal-transition-under-mod-block guard (exit 1, current text).
package status

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestArchiveRelativeDest locks that --archive's dest tracks the --workflow-dir
// spelling: a relative `--workflow-dir .` (run with cwd=root) yields
// `archived: ./_archive/{slug}.md`. Compared launcher-vs-oracle, both run from
// the same relative spelling, and the moved file lands under _archive.
func TestArchiveRelativeDest(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := []string{"--workflow-dir", ".", "--archive", "001-design-seam"}
	nOut, nErr, nCode := runNative(t, root, env, args...)

	if nCode != 0 {
		t.Fatalf("exit: native=%d (%q)", nCode, nErr)
	}
	want := "archived: ./_archive/001-design-seam.md\n"
	if nOut != want {
		t.Fatalf("native narration = %q, want %q (relative dest spelling)", nOut, want)
	}
	// The entity actually moved under _archive and left the active dir.
	if _, err := os.Stat(filepath.Join(root, "_archive", "001-design-seam.md")); err != nil {
		t.Fatalf("archived file missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "001-design-seam.md")); !os.IsNotExist(err) {
		t.Fatalf("source file should be gone after archive, stat err=%v", err)
	}
}

// TestArchiveAbsoluteDest locks the absolute-spelling case: an absolute
// --workflow-dir yields an absolute archived: dest. Compared launcher-vs-oracle
// with the root prefix normalized so no machine path is asserted literally.
func TestArchiveAbsoluteDest(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := []string{"--workflow-dir", root, "--archive", "001-design-seam"}
	nOut, nErr, nCode := runNative(t, root, env, args...)

	if nCode != 0 {
		t.Fatalf("exit: native=%d (%q)", nCode, nErr)
	}
	// Native emits an absolute dest (not realpath'd) under its own root.
	wantNative := "archived: " + filepath.Join(root, "_archive", "001-design-seam.md") + "\n"
	if nOut != wantNative {
		t.Fatalf("native narration = %q, want %q (absolute dest spelling)", nOut, wantNative)
	}
}

// TestTerminalSetUnderSelfRefRejected locks AC-1: under `require-external-proof:
// true`, a terminal --set on an entity whose ACs are only self-referentially
// proven exits 1, leaves the frontmatter byte-identical, prints the guard error
// in the mod-block idiom, and --force bypasses with the standard warning. A
// real-proof entity terminalizes cleanly (exit 0). A flip test re-runs --set
// after editing the self-ref entity's Verified by clause to cite a Go test —
// the guard now lets it through, confirming it keys on the proof clause not
// the entity slug.
func TestTerminalSetUnderSelfRefRejected(t *testing.T) {
	t.Run("self-ref-rejected", func(t *testing.T) {
		env := pinnedEnv(t)
		root := stageExternalProofFixture(t, "require-external-proof: true\n")

		args := []string{"--workflow-dir", root, "--set", "010-self-ref-only", "status=done"}
		nOut, nErr, nCode := runNative(t, root, env, args...)

		if nCode != 1 {
			t.Fatalf("guard must reject self-ref terminal --set, got exit=%d", nCode)
		}
		if !strings.Contains(nErr, "self-referential proof") {
			t.Fatalf("stderr should name the rejection cause, got %q", nErr)
		}
		if !strings.Contains(nErr, "AC-1") {
			t.Fatalf("stderr should name the flagged AC, got %q", nErr)
		}
		if nOut != "" {
			t.Fatalf("stdout must be empty on rejection: %q", nOut)
		}
		fm := readFrontmatter(t, filepath.Join(root, "010-self-ref-only.md"))
		if !strings.Contains(fm, "status: implementation") {
			t.Fatalf("entity was mutated despite guard rejection:\n%s", fm)
		}
		if strings.Contains(fm, "status: done") {
			t.Fatalf("entity advanced to done despite guard rejection:\n%s", fm)
		}
	})

	t.Run("real-proof-passes", func(t *testing.T) {
		env := pinnedEnv(t)
		root := stageExternalProofFixture(t, "require-external-proof: true\n")

		args := []string{"--workflow-dir", root, "--set", "020-real-proof", "status=done"}
		_, nErr, nCode := runNative(t, root, env, args...)
		if nCode != 0 {
			t.Fatalf("real-proof entity should pass guard, got exit=%d stderr=%q", nCode, nErr)
		}
		fm := readFrontmatter(t, filepath.Join(root, "020-real-proof.md"))
		if !strings.Contains(fm, "status: done") {
			t.Fatalf("real-proof entity should advance to done, fm=%s", fm)
		}
	})

	t.Run("force-bypasses-with-warning", func(t *testing.T) {
		env := pinnedEnv(t)
		root := stageExternalProofFixture(t, "require-external-proof: true\n")

		args := []string{"--workflow-dir", root, "--set", "030-force-bypass", "status=done", "--force"}
		_, nErr, nCode := runNative(t, root, env, args...)
		if nCode != 0 {
			t.Fatalf("--force must bypass guard, got exit=%d", nCode)
		}
		if !strings.Contains(nErr, "Warning: --force overriding require-external-proof") {
			t.Fatalf("stderr should carry the bypass warning, got %q", nErr)
		}
		fm := readFrontmatter(t, filepath.Join(root, "030-force-bypass.md"))
		if !strings.Contains(fm, "status: done") {
			t.Fatalf("force-bypass entity should advance to done, fm=%s", fm)
		}
	})

	t.Run("flip-cite-go-test", func(t *testing.T) {
		env := pinnedEnv(t)
		root := stageExternalProofFixture(t, "require-external-proof: true\n")

		// Rewrite the self-ref entity's Verified by clause to cite a Go test.
		path := filepath.Join(root, "010-self-ref-only.md")
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		replaced := strings.Replace(string(body),
			"Verified by: review of this entity's decision section, which states the\nintent and cites the design rationale.",
			"Verified by: a Go test `TestSomething` asserts the intent holds.",
			1)
		if replaced == string(body) {
			t.Fatalf("flip rewrite did not match the fixture body")
		}
		if err := os.WriteFile(path, []byte(replaced), 0o644); err != nil {
			t.Fatal(err)
		}
		// Re-stage the git index so the entity is committed cleanly.
		gitCmd := func(args ...string) {
			cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("git %v: %v\n%s", args, err, out)
			}
		}
		gitCmd("-c", "user.email=t@t", "-c", "user.name=t", "add", "-A")
		gitCmd("-c", "user.email=t@t", "-c", "user.name=t", "commit", "-q", "-m", "flip")

		args := []string{"--workflow-dir", root, "--set", "010-self-ref-only", "status=done"}
		_, nErr, nCode := runNative(t, root, env, args...)
		if nCode != 0 {
			t.Fatalf("guard should let through a real-proof AC, got exit=%d stderr=%q", nCode, nErr)
		}
		fm := readFrontmatter(t, path)
		if !strings.Contains(fm, "status: done") {
			t.Fatalf("flipped entity should advance to done, fm=%s", fm)
		}
	})
}

// TestTerminalSetUnderModBlockRejected locks the guard: a terminal --set
// (status -> terminal stage) on an entity with an active mod-block exits 1 with
// the current error text, and the entity is NOT mutated. Compared launcher vs
// oracle for exit code, stderr text, and unchanged frontmatter.
func TestTerminalSetUnderModBlockRejected(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "guard-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "010-blocked", "status=done")
	nOut, nErr, nCode := runNative(t, root, env, args...)

	if nCode != 1 {
		t.Fatalf("native exit=%d, want 1 (guard must reject)", nCode)
	}
	wantErr := "Error: entity 010-blocked has pending mod-block (merge:pr-merge). Clear mod-block in a separate --set call, or use --force."
	if !strings.Contains(nErr, wantErr) {
		t.Fatalf("native stderr = %q, want it to contain %q", nErr, wantErr)
	}
	if nOut != "" {
		t.Fatalf("stdout must be empty on rejection: native=%q", nOut)
	}
	assertEnvelopeGolden(t, "archive-guard-terminal-modblock", goldenEnvelope{
		stdout: normalize(nOut, root), stderr: normalize(nErr, root), exit: nCode,
	})
	// The entity status must be unchanged (still implementation, not done).
	fm := readFrontmatter(t, filepath.Join(root, "010-blocked.md"))
	if !strings.Contains(fm, "status: implementation") {
		t.Fatalf("entity was mutated despite guard rejection:\n%s", fm)
	}
	if strings.Contains(fm, "status: done") {
		t.Fatalf("entity advanced to done despite guard rejection:\n%s", fm)
	}
}
