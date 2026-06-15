package ensigncycle

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// seedScenarioBlock captures the IDs declared between the machine-readable
// `<!-- seed-scenarios -->` markers in scenario-testing-principles.md.
var seedScenarioBlock = regexp.MustCompile(`(?s)<!-- seed-scenarios -->(.*?)<!-- /seed-scenarios -->`)

// seedScenarioID matches one backtick-quoted scenario ID at the start of a seed
// bullet line: `- ` + "`id`".
var seedScenarioID = regexp.MustCompile("(?m)^- `([^`]+)`")

// TestSeedScenariosDocLock is the AC-5 doc-lock: it binds the machine-readable
// `<!-- seed-scenarios -->` block in docs/specs/scenario-testing-principles.md to
// the code-backed sharedRuntimeScenarios() table. The doc declares the seed IDs as
// the human-readable face of a code-backed truth, so the two sets must be equal;
// this test reds on drift in EITHER direction — a scenario added, dropped, or
// renamed on one side without the other. It is offline (default tag): it reads the
// doc and the table, spending no model.
func TestSeedScenariosDocLock(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	docPath := filepath.Join(wd, "..", "..", "docs", "specs", "scenario-testing-principles.md")
	b, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("read docs/specs/scenario-testing-principles.md: %v", err)
	}

	block := seedScenarioBlock.FindStringSubmatch(string(b))
	if block == nil {
		t.Fatal("scenario-testing-principles.md is missing the `<!-- seed-scenarios -->` machine-readable block")
	}
	var docIDs []string
	for _, m := range seedScenarioID.FindAllStringSubmatch(block[1], -1) {
		docIDs = append(docIDs, m[1])
	}
	if len(docIDs) == 0 {
		t.Fatal("the `<!-- seed-scenarios -->` block declared no scenario IDs")
	}

	var codeIDs []string
	for _, scenario := range sharedRuntimeScenarios() {
		codeIDs = append(codeIDs, scenario.name)
	}

	sort.Strings(docIDs)
	sort.Strings(codeIDs)
	if strings.Join(docIDs, ",") != strings.Join(codeIDs, ",") {
		t.Fatalf("seed-scenarios doc block drifted from sharedRuntimeScenarios():\n  doc:  %v\n  code: %v", docIDs, codeIDs)
	}
}
