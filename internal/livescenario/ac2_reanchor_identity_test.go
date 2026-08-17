package livescenario

import (
	"os/exec"
	"strings"
	"testing"
)

// TestACReanchorFixtureRepoCommitsFromAnyProcess pins the fixture repo's git
// identity in its OWN config, not just on the fixture's own commit.
//
// This guards a silent-failure class rather than a cosmetic one. The fixture used
// to pass `-c user.email=… -c user.name=…` on its init commit; those flags scope to
// that single invocation, so the repo itself carried no identity and the live FO's
// spacedock process could not commit into it on a CI runner with no global config.
// Nothing surfaced, because the failing commit was a no-op until the honest state
// commit landed. Reading the config back — rather than asserting the fixture's own
// commit succeeded — is what makes this falsifiable: the old `-c` form passes any
// check that only exercises the fixture's commit.
//
// internal/testgit has the same guard for the path every other fixture uses.
func TestACReanchorFixtureRepoCommitsFromAnyProcess(t *testing.T) {
	dir := t.TempDir()
	if _, err := AuthorACReanchorScenario().Setup(dir); err != nil {
		t.Fatal(err)
	}
	for field, want := range map[string]string{
		"user.name":  "Spacedock Test",
		"user.email": "spacedock@example.invalid",
	} {
		out, err := exec.Command("git", "-C", dir, "config", "--local", field).Output()
		if err != nil {
			t.Fatalf("git config --local %s: %v — the repo carries no identity, so no other process can commit in it", field, err)
		}
		if got := strings.TrimSpace(string(out)); got != want {
			t.Errorf("git config --local %s = %q, want %q", field, got, want)
		}
	}
}
