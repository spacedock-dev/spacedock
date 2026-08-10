// ABOUTME: Offline test for the broad-search-at-boot detector — proves it reds a
// ABOUTME: zero-discover boot that filesystem-sweeps to hunt a workflow, passes report-and-stop.
package ensigncycle

import (
	"strings"
	"testing"
)

// globLine builds one assistant-with-Glob-tool_use stream-json line carrying the
// given pattern, the shape detectBroadSearchAtBoot scans for a recursive sweep.
func globLine(pattern string) string {
	return `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Glob","input":{"pattern":` +
		mustJSONString(pattern) + `}}]}}`
}

// grepLine builds one assistant-with-Grep-tool_use stream-json line carrying the
// given pattern and search path (path "" → repo-wide).
func grepLine(pattern, path string) string {
	return `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Grep","input":{"pattern":` +
		mustJSONString(pattern) + `,"path":` + mustJSONString(path) + `}}]}}`
}

// TestDetectBroadSearchAtBoot covers the pure detector both ways: it reds a
// zero-discover boot that broad-searches the filesystem to hunt a workflow (a
// repo-rooted find / grep -r / ls -R Bash, or a recursive Glob/Grep over the
// project root), and passes a clean report-and-stop boot, a scoped search under an
// already-resolved workflow dir, and an empty stream.
func TestDetectBroadSearchAtBoot(t *testing.T) {
	const fixtureRoot = "/tmp/TestLiveZeroDiscover/001"

	cases := []struct {
		name      string
		lines     []string
		wantRed   bool
		wantNames string // substring the error must name when wantRed
	}{
		{
			name: "repo_rooted_find_reds",
			lines: []string{
				streamLine(`spacedock --version`),
				streamLine(`spacedock status --discover`),
				streamLine(`find ` + fixtureRoot + ` -name README.md`),
			},
			wantRed:   true,
			wantNames: "find",
		},
		{
			name: "grep_r_repo_root_reds",
			lines: []string{
				streamLine(`spacedock status --discover`),
				streamLine(`grep -r "commissioned-by: spacedock@" ` + fixtureRoot),
			},
			wantRed:   true,
			wantNames: "grep -r",
		},
		{
			name: "ls_recursive_repo_root_reds",
			lines: []string{
				streamLine(`spacedock status --discover`),
				streamLine(`ls -R ` + fixtureRoot),
			},
			wantRed:   true,
			wantNames: "ls -R",
		},
		{
			name: "ls_non_recursive_repo_root_passes",
			lines: []string{
				streamLine(`spacedock status --discover`),
				streamLine(`ls ` + fixtureRoot),
			},
			wantRed: false,
		},
		{
			name: "bare_ls_default_cwd_passes",
			lines: []string{
				streamLine(`spacedock status --discover`),
				streamLine(`ls -la`),
			},
			wantRed: false,
		},
		{
			name: "ls_la_repo_root_passes",
			lines: []string{
				streamLine(`spacedock status --discover`),
				streamLine(`ls -la ` + fixtureRoot),
			},
			wantRed: false,
		},
		{
			name: "ls_ltr_repo_root_passes",
			lines: []string{
				streamLine(`spacedock status --discover`),
				// -r on ls is reverse-sort, not recursion -- must not be read as -R.
				streamLine(`ls -ltr ` + fixtureRoot),
			},
			wantRed: false,
		},
		{
			name: "ls_scoped_under_resolved_workflow_passes",
			lines: []string{
				streamLine(`spacedock status --discover`),
				streamLine(`ls ` + fixtureRoot + `/docs/dev`),
			},
			wantRed: false,
		},
		{
			name: "find_path_arg_less_reds",
			lines: []string{
				streamLine(`spacedock status --discover`),
				streamLine(`find`),
			},
			wantRed:   true,
			wantNames: "find",
		},
		{
			name: "artifact_9074747236_adapter_find_reds",
			lines: []string{
				streamLine(`find /tmp/spacedock-live-plugin-3439009114/skills/first-officer/references -iname "*claude*"`),
			},
			wantRed:   true,
			wantNames: "find",
		},
		{
			name: "find_scoped_under_resolved_workflow_passes",
			lines: []string{
				streamLine(`spacedock status --discover`),
				streamLine(`find ` + fixtureRoot + `/docs/dev -name README.md`),
			},
			wantRed: false,
		},
		{
			name: "recursive_glob_readme_reds",
			lines: []string{
				streamLine(`spacedock status --discover`),
				globLine(`**/README.md`),
			},
			wantRed:   true,
			wantNames: "**/README.md",
		},
		{
			name: "repo_wide_grep_tool_reds",
			lines: []string{
				streamLine(`spacedock status --discover`),
				grepLine(`commissioned-by`, ""),
			},
			wantRed:   true,
			wantNames: "commissioned-by",
		},
		{
			name: "clean_report_and_stop_passes",
			lines: []string{
				streamLine(`spacedock --version`),
				streamLine(`git rev-parse --show-toplevel`),
				streamLine(`spacedock status --discover`),
			},
			wantRed: false,
		},
		{
			name: "scoped_grep_under_resolved_workflow_passes",
			lines: []string{
				streamLine(`spacedock status --discover`),
				streamLine(`grep -r "stages:" ` + fixtureRoot + `/docs/dev/README.md`),
				grepLine(`feedback-to`, fixtureRoot+`/docs/dev`),
			},
			wantRed: false,
		},
		{
			name:    "empty_stream_passes",
			lines:   nil,
			wantRed: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stream := strings.Join(tc.lines, "\n")
			err := detectBroadSearchAtBoot(stream, fixtureRoot)
			if tc.wantRed {
				if err == nil {
					t.Fatalf("detector passed a boot that broad-searched the filesystem — want a red\nstream:\n%s", stream)
				}
				if tc.wantNames != "" && !strings.Contains(err.Error(), tc.wantNames) {
					t.Errorf("error must name the offending search %q: %v", tc.wantNames, err)
				}
			} else if err != nil {
				t.Errorf("detector red a legitimate boot: %v\nstream:\n%s", err, stream)
			}
		})
	}
}

func TestDetectBroadSearchAtBootPR495CapturedFilingHunt(t *testing.T) {
	const fixtureRoot = "/tmp/TestLiveClaudeSharedScenariosfiling180063050/001"
	const captured = `find / -path /proc -prune -o -iname "fo-write-core*" -print 2>/dev/null`
	err := detectBroadSearchAtBoot(streamLine(captured), fixtureRoot)
	if err == nil {
		t.Fatal("PR #495 captured filing hunt was not detected")
	}
	if !strings.Contains(err.Error(), "find /") || !strings.Contains(err.Error(), "fo-write-core") {
		t.Fatalf("diagnostic does not identify captured filing hunt %q: %v", captured, err)
	}
}

// TestDetectBroadSearchAtBootSecondBlock proves the detector iterates ALL tool_use
// blocks of a multi-tool assistant turn, not just the first — a broad-search Bash
// riding as a second block must still red (it cannot be evaded by block ordering).
func TestDetectBroadSearchAtBootSecondBlock(t *testing.T) {
	const fixtureRoot = "/tmp/TestLiveZeroDiscover/001"
	// One assistant turn with a benign first tool_use and a repo-rooted find second.
	line := `{"type":"assistant","message":{"content":[` +
		`{"type":"tool_use","name":"Bash","input":{"command":` + mustJSONString(`spacedock status --discover`) + `}},` +
		`{"type":"tool_use","name":"Bash","input":{"command":` + mustJSONString(`find `+fixtureRoot+` -name README.md`) + `}}]}}`
	err := detectBroadSearchAtBoot(line, fixtureRoot)
	if err == nil {
		t.Fatal("detector missed a broad-search Bash in the SECOND tool_use block of a multi-tool turn")
	}
	if !strings.Contains(err.Error(), "find") {
		t.Errorf("error must name the offending find command: %v", err)
	}
}
