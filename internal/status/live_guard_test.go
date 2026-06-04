// ABOUTME: Terminal --set live-run guard tests — refuse an uncited runtime-
// ABOUTME: observable AC, pass a cited/offline one, inert when opt-in is off.
package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// stageLiveGuardFixture builds a fresh git-initialized workflow opting into
// `require-external-proof: true` with one entity at `implementation` carrying
// the supplied AC body. Returns the absolute root and the entity path. The AC
// proof uses session: citations (offline-resolvable) so the guard is exercised
// end to end with no network.
func stageLiveGuardFixture(t *testing.T, optIn, acBody string) (string, string) {
	t.Helper()
	root := t.TempDir()
	readme := "---\n" +
		"entity-type: task\n" +
		"entity-label: task\n" +
		"entity-label-plural: tasks\n" +
		"id-style: sequential\n" +
		optIn +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: backlog\n" +
		"      initial: true\n" +
		"      gate: true\n" +
		"    - name: implementation\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n\n# Live-Guard Fixture\n"
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	entityPath := filepath.Join(root, "010-live-ac.md")
	entity := "---\n" +
		"id: \"010\"\n" +
		"title: Live AC entity\n" +
		"status: implementation\n" +
		"score: \"0.50\"\n" +
		"source: roadmap\n" +
		"---\n\n# Live AC entity\n\n## Acceptance criteria\n\n" + acBody
	if err := os.WriteFile(entityPath, []byte(entity), 0o644); err != nil {
		t.Fatal(err)
	}
	gitInit(t, root)
	return root, entityPath
}

// TestLiveRunGuardRefusesUncitedRuntimeAC locks AC-1's refusal half + AC-4: a
// runtime-observable AC declared `Verified by: live …` with no resolvable
// live-run citation (a placeholder — yy's shape) is REFUSED at terminal --set
// under require-external-proof: true: exit 1, frontmatter untouched, diagnostic
// names the live-run requirement.
func TestLiveRunGuardRefusesUncitedRuntimeAC(t *testing.T) {
	env := pinnedEnv(t)
	acBody := "**AC-1 — The model exits on the contract at runtime.**\n" +
		"Verified by: live <none — passed offline + 3 audit cycles, merged pending-live-run>\n"
	root, entityPath := stageLiveGuardFixture(t, "require-external-proof: true\n", acBody)

	_, nErr, nCode := runNative(t, root, env, "--workflow-dir", root, "--set", "010-live-ac", "status=done")
	if nCode != 1 {
		t.Fatalf("uncited runtime-observable AC must be refused (exit 1), got %d (%q)", nCode, nErr)
	}
	if !strings.Contains(nErr, "live run") && !strings.Contains(nErr, "live-run") {
		t.Fatalf("diagnostic must name the live-run requirement, got %q", nErr)
	}
	// The masked-404 case (a private/unscoped-token repo where the run exists but
	// gh returns 404) also lands as definitivelyAbsent → refuse. The diagnostic
	// must name BOTH remediations so a token-scope failure is not sent chasing a
	// wrong fix: cite a resolvable ref OR check the token's repo scope.
	if !strings.Contains(strings.ToLower(nErr), "scope") {
		t.Fatalf("refusal diagnostic must name the token-scope remediation (masked-404), got %q", nErr)
	}
	fm := readFrontmatter(t, entityPath)
	if !strings.Contains(fm, "status: implementation") {
		t.Fatalf("frontmatter must be untouched on refusal, got %s", fm)
	}
}

// TestLiveRunGuardPassesCitedRuntimeAC locks AC-1's pass half: a runtime-
// observable AC citing a RESOLVABLE live-run ref (a present session .jsonl) is
// NOT refused — terminal --set succeeds and the frontmatter advances.
func TestLiveRunGuardPassesCitedRuntimeAC(t *testing.T) {
	env := pinnedEnv(t)
	dir := t.TempDir()
	sessionPath := filepath.Join(dir, "run.jsonl")
	if err := os.WriteFile(sessionPath, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	acBody := "**AC-1 — The model exits on the contract at runtime.**\n" +
		"Verified by: live session:" + sessionPath + "\n"
	root, entityPath := stageLiveGuardFixture(t, "require-external-proof: true\n", acBody)

	_, nErr, nCode := runNative(t, root, env, "--workflow-dir", root, "--set", "010-live-ac",
		"status=done", "completed", "verdict=accepted")
	if nCode != 0 {
		t.Fatalf("cited runtime-observable AC must pass, got %d (%q)", nCode, nErr)
	}
	fm := readFrontmatter(t, entityPath)
	if !strings.Contains(fm, "status: done") {
		t.Fatalf("entity should advance to done, got %s", fm)
	}
}

// TestLiveRunGuardNeverRefusesOfflineAC locks AC-1's no-over-gate half: an
// offline-checkable AC (any non-`live` proof clause) is NEVER refused, even
// under require-external-proof: true.
func TestLiveRunGuardNeverRefusesOfflineAC(t *testing.T) {
	env := pinnedEnv(t)
	acBody := "**AC-1 — The parser round-trips.**\n" +
		"Verified by: a Go unit test `TestRoundTrip` in `internal/status/parse_test.go` asserts the exit code.\n"
	root, entityPath := stageLiveGuardFixture(t, "require-external-proof: true\n", acBody)

	_, nErr, nCode := runNative(t, root, env, "--workflow-dir", root, "--set", "010-live-ac",
		"status=done", "completed", "verdict=accepted")
	if nCode != 0 {
		t.Fatalf("offline AC must never be refused by the live-run guard, got %d (%q)", nCode, nErr)
	}
	fm := readFrontmatter(t, entityPath)
	if !strings.Contains(fm, "status: done") {
		t.Fatalf("offline-AC entity should advance to done, got %s", fm)
	}
}

// TestLiveRunGuardInertWhenOptInAbsentOrFalse locks AC-1's inert half: with
// require-external-proof absent or false, the live-run guard never fires — an
// uncited runtime-observable AC advances to terminal (no over-gate).
func TestLiveRunGuardInertWhenOptInAbsentOrFalse(t *testing.T) {
	env := pinnedEnv(t)
	acBody := "**AC-1 — The model exits on the contract at runtime.**\n" +
		"Verified by: live <no artifact yet>\n"
	for _, optIn := range []string{"", "require-external-proof: false\n"} {
		name := optIn
		if name == "" {
			name = "absent"
		}
		t.Run(strings.TrimSpace(name), func(t *testing.T) {
			root, entityPath := stageLiveGuardFixture(t, optIn, acBody)
			_, nErr, nCode := runNative(t, root, env, "--workflow-dir", root, "--set", "010-live-ac",
				"status=done", "completed", "verdict=accepted")
			if nCode != 0 {
				t.Fatalf("guard must be inert when opt-in absent/false, got %d (%q)", nCode, nErr)
			}
			fm := readFrontmatter(t, entityPath)
			if !strings.Contains(fm, "status: done") {
				t.Fatalf("entity should advance to done when guard inert, got %s", fm)
			}
		})
	}
}

// TestLiveRunGuardForceBypassWarnsLoudly locks AC-1's --force escape: --force
// bypasses the refusal but emits a loud, risk-naming warning (naming the
// runtime-observable AC being terminalized without a cited live run — the yy
// slip), so a default merge cannot silently skip it.
func TestLiveRunGuardForceBypassWarnsLoudly(t *testing.T) {
	env := pinnedEnv(t)
	acBody := "**AC-1 — The model exits on the contract at runtime.**\n" +
		"Verified by: live <no artifact yet>\n"
	root, entityPath := stageLiveGuardFixture(t, "require-external-proof: true\n", acBody)

	_, nErr, nCode := runNative(t, root, env, "--workflow-dir", root, "--set", "010-live-ac",
		"status=done", "completed", "verdict=accepted", "--force")
	if nCode != 0 {
		t.Fatalf("--force must bypass the live-run guard, got %d (%q)", nCode, nErr)
	}
	if !strings.Contains(nErr, "live run") {
		t.Fatalf("--force bypass must warn loudly naming the live-run risk, got %q", nErr)
	}
	fm := readFrontmatter(t, entityPath)
	if !strings.Contains(fm, "status: done") {
		t.Fatalf("entity should advance under --force, got %s", fm)
	}
}
