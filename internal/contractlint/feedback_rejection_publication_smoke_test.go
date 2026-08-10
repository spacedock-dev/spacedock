// ABOUTME: Pins complete rejected-round publication before reviewer rerun and gate preparation.
package contractlint

import (
	"strings"
	"testing"
)

func TestFeedbackRejectionPublishesCompleteRoundBeforeRegating(t *testing.T) {
	body := readSkillFile(t, "skills/feedback-rejection-flow/SKILL.md")
	ordered := []string{
		"completion-signal",
		"append its authorized line",
		"gate record --round STAGE/CYCLE",
		"Re-run the reviewer",
		"Re-enter the normal gate flow",
	}
	previous := -1
	for _, token := range ordered {
		position := strings.Index(body, token)
		if position < 0 {
			t.Fatalf("feedback rejection skill is missing %q", token)
		}
		if position <= previous {
			t.Fatalf("feedback rejection skill places %q before the complete-round publication sequence", token)
		}
		previous = position
	}
	for _, token := range []string{
		"complete round summary",
		"hold the flow",
		"Do not claim that the round was recorded",
	} {
		if !strings.Contains(body, token) {
			t.Fatalf("feedback rejection skill is missing the recorder result guard %q", token)
		}
	}
}

func TestDevelopmentTemplatePublishesRejectedRoundBeforeRegating(t *testing.T) {
	body := readSkillFile(t, "skills/commission/references/templates/development.md")
	if !strings.Contains(body, "before reviewer re-run or next-gate preparation") {
		t.Fatal("development template does not place rejected-round publication before reviewer rerun or next-gate preparation")
	}
}
