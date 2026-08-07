// ABOUTME: Pins stable v1 gate skills to the single chat decision path.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStableV1GateSkillsExposeOnlyChatClosure(t *testing.T) {
	paths := []string{
		filepath.Join(skillsRoot(t), "present-gate", "SKILL.md"),
		filepath.Join(skillsRoot(t), "fo-gate-lifecycle", "SKILL.md"),
	}
	for _, path := range paths {
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, forbidden := range []string{"gate record <entity> --room", "gate <room>", "provider package", "fall back to chat"} {
			if strings.Contains(strings.ToLower(string(body)), forbidden) {
				t.Errorf("%s exposes provider-only stable-v1 wording %q", path, forbidden)
			}
		}
	}

	presentation, err := os.ReadFile(paths[0])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(presentation), "Stable v1 presents gates in chat only.") ||
		!strings.Contains(string(presentation), "gate record <entity> --decision") {
		t.Fatal("present-gate does not state the stable-v1 chat presentation and decision-record contract")
	}
	lifecycle, err := os.ReadFile(paths[1])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(lifecycle), "gate record ENTITY --decision approve|revise|hold") {
		t.Fatal("gate lifecycle no longer exposes the ordinary chat decision sequence")
	}
}
