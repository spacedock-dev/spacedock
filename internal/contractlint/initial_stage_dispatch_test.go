package contractlint

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestInitialStageDispatchUsesOneSelectedTarget(t *testing.T) {
	body := readRepoFile(t, filepath.FromSlash("skills/first-officer/references/fo-dispatch-core.md"))
	for _, want := range []string{
		"`current` is initial and `next` is terminal, set `dispatch_stage = current`",
		"`status={dispatch_stage}`",
		"`dispatch: {slug} entering {dispatch_stage}`",
		"`«dispatch.build» --stage {dispatch_stage}`",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("dispatch contract missing %q", want)
		}
	}
	if strings.Contains(body, "Initial-stage successor rows retain legacy meaning.") {
		t.Error("dispatch contract retains the legacy initial-stage successor rule")
	}
}
