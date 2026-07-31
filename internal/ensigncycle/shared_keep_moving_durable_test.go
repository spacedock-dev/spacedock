package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	kmApprovedGate = "approved-gate"
	kmReadyOne     = "ready-one"
	kmReadyTwo     = "ready-two"
	kmQuestioned   = "questioned"
	kmNextStage    = "implementation"
)

func kmExpected() map[string]string {
	return map[string]string{
		kmApprovedGate: kmNextStage, kmReadyOne: kmNextStage, kmReadyTwo: kmNextStage,
	}
}

type durableCommit struct {
	message, blob    string
	entityFileScoped bool
	entityOwned      bool
}

func gradeDurableTaskJourneys(t *testing.T, root string, expected map[string]string) (int, map[string]string) {
	completed, failures := 0, map[string]string{}
	for slug, stage := range expected {
		if reason := durableTaskJourneyFailure(t, root, slug, stage); reason != "" {
			failures[slug] = reason
		} else {
			completed++
		}
	}
	return completed, failures
}
func durableTaskJourneyFailure(t *testing.T, root, slug, stage string) string {
	active := slug + ".md"
	archive := filepath.Join("_archive", slug+".md")
	logPath := active
	if _, err := os.Stat(filepath.Join(root, archive)); err == nil {
		logPath = archive
	}
	history := durableEntityHistory(t, root, slug, logPath)
	dispatch, report, terminal := -1, -1, -1
	reportBefore, unscopedReport, archived := false, false, false
	for i, c := range history {
		hasReport := strings.Contains(c.blob, "\n## Stage Report: "+stage+"\n")
		if c.message == "dispatch: "+slug+" entering "+stage && c.entityFileScoped &&
			durableField(c.blob, "status") == stage && durableField(c.blob, "started") != "" {
			dispatch = i
		}
		if hasReport {
			if dispatch < 0 || i <= dispatch {
				reportBefore = true
			} else if !c.entityFileScoped {
				unscopedReport = true
			} else if report < 0 && !unscopedReport && !reportBefore {
				report = i
			}
		}
		if report >= 0 && i > report && c.entityFileScoped && durableField(c.blob, "status") == "done" &&
			durableField(c.blob, "completed") != "" && durableField(c.blob, "verdict") != "" {
			terminal = i
		}
		if terminal >= 0 && i >= terminal && c.entityOwned && i == len(history)-1 {
			archived = true
		}
	}
	switch {
	case dispatch < 0:
		return "missing path-scoped dispatch entry with stage and started"
	case report < 0 && unscopedReport:
		return "missing path-scoped worker report after dispatch"
	case report < 0 && reportBefore:
		return "missing worker report after dispatch; stale report precedes dispatch"
	case report < 0:
		return "missing worker report after dispatch"
	case terminal < 0:
		return "missing later terminal fields completed and verdict"
	case !archived:
		return "missing canonical archive after terminalization"
	}
	if _, err := os.Stat(filepath.Join(root, active)); !os.IsNotExist(err) {
		return "active entity remains beside canonical archive"
	}
	return ""
}
func durableEntityHistory(t *testing.T, root, slug, logPath string) []durableCommit {
	out := git(t, root, "log", "--follow", "--format=%H%x09%s", "--", logPath)
	var history []durableCommit
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		fields := strings.SplitN(line, "\t", 2)
		if len(fields) != 2 {
			continue
		}
		blob := durableBlobAt(root, fields[0], slug)
		files := strings.Fields(git(t, root, "diff-tree", "--root", "--no-commit-id", "--name-only", "-r", "-M", fields[0]))
		scoped, entityOwned := len(files) > 0, len(files) > 0
		for _, file := range files {
			if file != slug+".md" && file != filepath.Join("_archive", slug+".md") {
				scoped = false
			}
			if !durablePathOwnedBySlug(slug, file) {
				entityOwned = false
			}
		}
		history = append(history, durableCommit{fields[1], blob, scoped, entityOwned})
	}
	return history
}
func durablePathOwnedBySlug(slug, path string) bool {
	return path == slug+".md" || path == filepath.Join("_archive", slug+".md") ||
		strings.HasPrefix(path, slug+"/") ||
		strings.HasPrefix(path, filepath.Join("_archive", slug)+"/")
}
func durableBlobAt(root, hash, slug string) string {
	for _, path := range []string{filepath.Join("_archive", slug+".md"), slug + ".md"} {
		cmd := exec.Command("git", "-C", root, "show", hash+":"+path)
		if out, err := cmd.Output(); err == nil {
			return string(out)
		}
	}
	return ""
}
func durableField(content, name string) string {
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, name+":") {
			return strings.TrimSpace(strings.TrimPrefix(line, name+":"))
		}
	}
	return ""
}
func assertDurableKeepMoving(t *testing.T, root string) error {
	completed, failures := gradeDurableTaskJourneys(t, root, kmExpected())
	if completed != 3 {
		return fmt.Errorf("durable keep-moving journeys = %d/3: %v", completed, failures)
	}
	content, where, found := locateEntity(root, kmQuestioned)
	if !found || where != filepath.Join(root, kmQuestioned+".md") ||
		durableField(content, "status") == "done" || durableField(content, "completed") != "" || durableField(content, "verdict") != "" {
		return fmt.Errorf("%s must remain active and nonterminal", kmQuestioned)
	}
	history := durableEntityHistory(t, root, kmQuestioned, kmQuestioned+".md")
	if len(history) < 2 || history[len(history)-1].blob == history[0].blob || !history[len(history)-1].entityOwned {
		return fmt.Errorf("%s has no durable entity-owned re-shape", kmQuestioned)
	}
	return nil
}
func assertDurableSmallestMechanism(t *testing.T, root string, tr mechanismTrace, edits, commissioned []string) error {
	expected := map[string]string{}
	for _, slug := range commissioned {
		expected[slug] = "ready"
	}
	completed, failures := gradeDurableTaskJourneys(t, root, expected)
	if completed != len(expected) {
		return fmt.Errorf("durable commissioned journeys = %d/%d: %v", completed, len(expected), failures)
	}
	for _, slug := range commissioned {
		tr.engaged[slug] = true
	}
	return gradeSmallestSufficientMechanism(tr, edits, commissioned)
}
func codexCommandOutput(command, output string, exit int, status string) string {
	b, _ := json.Marshal(map[string]any{"type": "item.completed", "item": map[string]any{
		"type": "command_execution", "command": command, "aggregated_output": output,
		"exit_code": exit, "status": status,
	}})
	return string(b)
}
func TestDurableTaskJourneys(t *testing.T) {
	tests := []struct {
		name, mutation, slug, reason string
	}{
		{"three independent journeys", "", "", ""},
		{"missing dispatch", "missing-dispatch", "ready-one", "dispatch entry"},
		{"missing report", "missing-report", "ready-one", "worker report"},
		{"missing terminal fields", "missing-terminal", "ready-one", "terminal fields"},
		{"missing archive", "missing-archive", "ready-one", "canonical archive"},
		{"report before dispatch", "report-before-dispatch", "ready-one", "worker report after dispatch"},
		{"cross-attributed report", "cross-attributed-report", "ready-one", "path-scoped worker report"},
		{"cross-attributed archive", "cross-attributed-archive", "ready-one", "canonical archive"},
		{"slug-prefix archive", "slug-prefix-archive", "ready-one", "canonical archive"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			completed, failures := gradeDurableTaskJourneys(t, durableJourneyFixture(t, tt.mutation), kmExpected())
			wantFailures := 1
			if tt.slug == "" {
				wantFailures = 0
			}
			if completed != 3-wantFailures || len(failures) != wantFailures {
				t.Fatalf("completed journeys = %d/3, failures=%v; want %d/3 and %d failures", completed, failures, 3-wantFailures, wantFailures)
			}
			if tt.slug != "" && !strings.Contains(failures[tt.slug], tt.reason) {
				t.Fatalf("failure for %q = %q, want reason containing %q", tt.slug, failures[tt.slug], tt.reason)
			}
		})
	}
}
func durableJourneyFixture(t *testing.T, mutation string) string {
	root := t.TempDir()
	for _, slug := range []string{"approved-gate", "ready-one", "ready-two"} {
		writeFile(t, filepath.Join(root, slug+".md"), durableEntity(slug, "implementation", "", ""))
	}
	gitInit(t, root)
	for _, slug := range []string{"approved-gate", "ready-one", "ready-two"} {
		if mutation == "report-before-dispatch" && slug == "ready-one" {
			durableAppendReport(t, root, slug)
			gitCommitPathScoped(t, root, slug+".md", "worker: stale "+slug)
		}
		if mutation != "missing-dispatch" || slug != "ready-one" {
			writeFile(t, filepath.Join(root, slug+".md"), strings.Replace(readFile(t, filepath.Join(root, slug+".md")), "started: ", "started: 2026-07-31T00:00:00Z", 1))
			gitCommitPathScoped(t, root, slug+".md", "dispatch: "+slug+" entering implementation")
		}
		if (mutation != "missing-report" || slug != "ready-one") && !(mutation == "report-before-dispatch" && slug == "ready-one") {
			durableAppendReport(t, root, slug)
			if mutation == "cross-attributed-report" && slug == "ready-one" {
				writeFile(t, filepath.Join(root, "ready-two.md"), readFile(t, filepath.Join(root, "ready-two.md"))+"\ncross-attributed\n")
				git(t, root, "add", "--", "ready-one.md", "ready-two.md")
				git(t, root, "commit", "-q", "-m", "worker: wrong scope", "--", "ready-one.md", "ready-two.md")
			} else {
				gitCommitPathScoped(t, root, slug+".md", "worker: "+slug)
			}
		}
		completed, verdict := "2026-07-31T00:01:00Z", "passed"
		if mutation == "missing-terminal" && slug == "ready-one" {
			completed, verdict = "", ""
		}
		content := readFile(t, filepath.Join(root, slug+".md"))
		content = strings.Replace(content, "status: implementation", "status: done", 1)
		content = strings.Replace(content, "completed:", "completed: "+completed, 1)
		content = strings.Replace(content, "verdict:", "verdict: "+verdict, 1)
		writeFile(t, filepath.Join(root, slug+".md"), content)
		gitCommitPathScoped(t, root, slug+".md", "terminalize: "+slug)
		if mutation != "missing-archive" || slug != "ready-one" {
			writeFile(t, filepath.Join(root, "_archive", slug+".md"), content)
			git(t, root, "rm", "-q", "--", slug+".md")
			git(t, root, "add", "--", "_archive/"+slug+".md")
			commitPaths := []string{slug + ".md", "_archive/" + slug + ".md"}
			if slug == "approved-gate" {
				sidecar := "_archive/approved-gate/review/briefing.json"
				writeFile(t, filepath.Join(root, sidecar), "{}\n")
				git(t, root, "add", "--", sidecar)
				commitPaths = append(commitPaths, sidecar)
			}
			if mutation == "cross-attributed-archive" && slug == "ready-one" {
				writeFile(t, filepath.Join(root, "ready-two.md"), readFile(t, filepath.Join(root, "ready-two.md"))+"\nforeign archive\n")
				git(t, root, "add", "--", "ready-two.md")
				commitPaths = append(commitPaths, "ready-two.md")
			}
			if mutation == "slug-prefix-archive" && slug == "ready-one" {
				foreign := "ready-one-other/review/briefing.json"
				writeFile(t, filepath.Join(root, foreign), "{}\n")
				git(t, root, "add", "--", foreign)
				commitPaths = append(commitPaths, foreign)
			}
			git(t, root, append([]string{"commit", "-q", "-m", "archive: " + slug, "--"}, commitPaths...)...)
		}
	}
	return root
}
func durableEntity(slug, stage, started, report string) string {
	return "---\nid: " + slug + "\ntitle: " + slug + "\nstatus: " + stage +
		"\nstarted: " + started + "\ncompleted:\nverdict:\n---\n# " + slug + "\n" + report
}
func durableAppendReport(t *testing.T, root, slug string) {
	path := filepath.Join(root, slug+".md")
	writeFile(t, path, readFile(t, path)+"\n## Stage Report: implementation\n\n- DONE: complete\n  durable evidence\n\n### Summary\n\nDone.\n")
}
