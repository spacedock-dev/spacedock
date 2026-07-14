package ensigncycle

import (
	"fmt"
	"strings"
	"testing"
)

type mergeModRecoveryObservation struct {
	activeAfter   string
	archivedAfter string
	gitClean      bool
}

func assertMergeModRecovery(o mergeModRecoveryObservation) error {
	if o.activeAfter != "" {
		return fmt.Errorf("merge-recovery remains active instead of archived")
	}
	if o.archivedAfter == "" {
		return fmt.Errorf("merge-recovery has no archived durable record")
	}
	frontmatter, err := parseOpeningFrontmatter(o.archivedAfter)
	if err != nil {
		return fmt.Errorf("archived merge-recovery frontmatter: %w", err)
	}
	if !strings.EqualFold(frontmatter["status"], "done") {
		return fmt.Errorf("archived merge-recovery is not status: done")
	}
	if !strings.EqualFold(frontmatter["verdict"], "PASSED") {
		return fmt.Errorf("archived merge-recovery is not verdict: PASSED")
	}
	modBlock, ok := frontmatter["mod-block"]
	if !ok || modBlock != "" {
		return fmt.Errorf("archived merge-recovery retains its mod-block")
	}
	if !o.gitClean {
		return fmt.Errorf("merge-recovery left the workflow git worktree dirty")
	}
	return nil
}

func parseOpeningFrontmatter(document string) (map[string]string, error) {
	normalized := strings.ReplaceAll(document, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return nil, fmt.Errorf("opening YAML fence is missing")
	}
	fields := make(map[string]string)
	for _, line := range lines[1:] {
		if line == "---" {
			return fields, nil
		}
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		if line[0] == ' ' || line[0] == '\t' {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok || strings.TrimSpace(key) == "" {
			return nil, fmt.Errorf("invalid top-level YAML line %q", line)
		}
		fields[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return nil, fmt.Errorf("closing YAML fence is missing")
}

func TestMergeModRecoveryDurableOracle(t *testing.T) {
	good := mergeModRecoveryObservation{
		archivedAfter: "---\nstatus: done\nverdict: PASSED\nmod-block:\n---\n",
		gitClean:      true,
	}
	if err := assertMergeModRecovery(good); err != nil {
		t.Fatalf("positive recovery: %v", err)
	}
	controls := []mergeModRecoveryObservation{
		{activeAfter: "status: implementation", archivedAfter: good.archivedAfter, gitClean: true},
		{gitClean: true},
		{archivedAfter: "status: implementation\nverdict: PASSED\nmod-block:\n", gitClean: true},
		{archivedAfter: "status: done\nverdict: PASSED\nmod-block: merge:pr-merge\n", gitClean: true},
		{archivedAfter: "---\nstatus: implementation\nverdict: REJECTED\nmod-block: merge:pr-merge\n---\n\nStale body example:\nstatus: done\nverdict: PASSED\nmod-block:\n", gitClean: true},
		{archivedAfter: good.archivedAfter, gitClean: false},
	}
	for i, control := range controls {
		if err := assertMergeModRecovery(control); err == nil {
			t.Errorf("broken recovery control %d passed", i)
		}
	}
}
