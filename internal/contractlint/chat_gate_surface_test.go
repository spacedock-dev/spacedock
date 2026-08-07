// ABOUTME: Pins stable v1 presentation channels to one semantic decision recorder.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStableV1GateSkillsUseOneRecorderAcrossPresentationChannels(t *testing.T) {
	paths := []string{
		filepath.Join(skillsRoot(t), "present-gate", "SKILL.md"),
		filepath.Join(skillsRoot(t), "fo-gate-lifecycle", "SKILL.md"),
	}
	for _, path := range paths {
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, forbidden := range []string{"gate record <entity> --room", "provider/result.json", "presented-inventory.json", "fall back to chat"} {
			if strings.Contains(strings.ToLower(string(body)), forbidden) {
				t.Errorf("%s exposes provider-only stable-v1 wording %q", path, forbidden)
			}
		}
	}

	presentation, err := os.ReadFile(paths[0])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(presentation), "Stable v1 permits chat or Subspace to present the committed gate.") ||
		!strings.Contains(string(presentation), "Both channels return semantic decision and reason input to the First Officer.") ||
		!strings.Contains(string(presentation), "gate record <entity> --decision") {
		t.Fatal("present-gate does not converge chat and Subspace presentation on the semantic decision recorder")
	}
	lifecycle, err := os.ReadFile(paths[1])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(lifecycle), "gate record ENTITY --decision approve|revise|hold") {
		t.Fatal("gate lifecycle no longer exposes the ordinary chat decision sequence")
	}
}
