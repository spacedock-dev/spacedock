package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
	hash, message, blob string
	entityFileScoped    bool
	entityOwned         bool
	files               []string
	added               map[string]bool
}

type durableJourneyProof struct {
	engaged, terminal string
}

func gradeDurableTaskJourneys(t *testing.T, root string, expected, batchExpected map[string]string) (int, map[string]string) {
	completed, failures := 0, map[string]string{}
	for slug, stage := range expected {
		if _, reason := durableTaskJourney(t, root, slug, stage, batchExpected); reason != "" {
			failures[slug] = reason
		} else {
			completed++
		}
	}
	return completed, failures
}
func durableTaskJourney(t *testing.T, root, slug, stage string, batchExpected map[string]string) (durableJourneyProof, string) {
	proof := durableJourneyProof{}
	active := slug + ".md"
	archive := filepath.Join("_archive", slug+".md")
	logPath := active
	if _, err := os.Stat(filepath.Join(root, archive)); err == nil {
		logPath = archive
	}
	history := durableEntityHistory(t, root, slug, logPath)
	dispatch, report, firstTerminal, terminal := -1, -1, -1, -1
	reportBefore, unscopedReport, archived := false, false, false
	for i, c := range history {
		hasReport := strings.Contains(c.blob, "\n## Stage Report: "+stage+"\n")
		atomicWorker := durableAtomicWorkerProof(history, i, slug, stage)
		if firstTerminal < 0 && durableTerminalState(c.blob) {
			firstTerminal = i
		}
		if durableDispatchProof(history, i, slug, stage) ||
			durableBatchDispatchProof(t, root, c, slug, batchExpected) {
			dispatch = i
		}
		if atomicWorker {
			dispatch, report = i, i
		} else if hasReport {
			if dispatch < 0 || i <= dispatch {
				reportBefore = true
			} else if !c.entityFileScoped {
				unscopedReport = true
			} else if report < 0 && !unscopedReport && !reportBefore {
				report = i
			}
		}
		terminalScoped := c.entityFileScoped || (c.entityOwned && i == len(history)-1) ||
			durableBatchTerminalProof(root, c, batchExpected)
		if firstTerminal == i && report >= 0 && i > report && terminalScoped && durableField(c.blob, "status") == "done" &&
			durableField(c.blob, "completed") != "" && durableField(c.blob, "verdict") != "" {
			terminal = i
		}
		if terminal >= 0 && i >= terminal && c.entityOwned && i == len(history)-1 {
			archived = true
		}
	}
	switch {
	case dispatch < 0:
		return proof, "missing path-scoped dispatch entry with stage and started"
	case report < 0 && unscopedReport:
		return proof, "missing path-scoped worker report after dispatch"
	case report < 0 && reportBefore:
		return proof, "missing worker report after dispatch; stale report precedes dispatch"
	case report < 0:
		return proof, "missing worker report after dispatch"
	case firstTerminal >= 0 && firstTerminal <= report:
		return proof, "first terminal transition must follow worker report"
	case terminal < 0:
		return proof, "missing later terminal fields completed and verdict"
	case durableField(history[len(history)-1].blob, "status") != "done" ||
		durableField(history[len(history)-1].blob, "completed") == "" ||
		durableField(history[len(history)-1].blob, "verdict") == "":
		return proof, "canonical archive missing terminal fields completed and verdict"
	case !archived:
		return proof, "missing canonical archive after terminalization"
	}
	if _, err := os.Stat(filepath.Join(root, active)); !os.IsNotExist(err) {
		return proof, "active entity remains beside canonical archive"
	}
	return durableJourneyProof{history[dispatch].hash, history[terminal].hash}, ""
}
func durableDispatchProof(history []durableCommit, i int, slug, stage string) bool {
	c := history[i]
	if c.message != "dispatch: "+slug+" entering "+stage ||
		durableField(c.blob, "status") != stage || durableField(c.blob, "started") == "" {
		return false
	}
	if c.entityFileScoped {
		return true
	}
	if i == 0 {
		return false
	}
	return durableNewRoomScope(history[i-1].blob, c, slug)
}
func durableNewRoomScope(parent string, c durableCommit, slug string) bool {
	room := durableNewRoomRef(parent, c.blob, slug)
	if room == "" {
		return false
	}
	for _, path := range c.files {
		if path != slug+".md" && (!c.added[path] || !strings.HasPrefix(path, room+"/")) {
			return false
		}
	}
	return true
}
func durableNewRoomRef(parent, child, slug string) string {
	room := durableRoomRef(child)
	if room == "" || durableRoomRef(parent) != "" || filepath.IsAbs(room) ||
		filepath.Clean(room) != room || !strings.HasPrefix(room, slug+"/") {
		return ""
	}
	return room
}
func durableBatchDispatchProof(t *testing.T, root string, c durableCommit, target string, expected map[string]string) bool {
	if len(expected) < 2 {
		return false
	}
	changed, allowed, rooms := map[string]bool{}, map[string]bool{}, map[string]string{}
	for _, path := range c.files {
		changed[path] = true
	}
	started := 0
	for slug, stage := range expected {
		path := slug + ".md"
		if !changed[path] {
			continue
		}
		parent, child := durableBlobAt(root, c.hash+"^", slug), durableBlobAt(root, c.hash, slug)
		if durableField(parent, "started") != "" ||
			durableField(child, "status") != stage || durableField(child, "started") == "" {
			return false
		}
		started++
		allowed[path] = true
		if durableRoomRef(parent) != durableRoomRef(child) {
			if rooms[slug] = durableNewRoomRef(parent, child, slug); rooms[slug] == "" {
				return false
			}
		}
	}
	if started < 2 || !allowed[target+".md"] {
		return false
	}
	if changed[kmQuestioned+".md"] {
		parent := durableBlobAt(root, c.hash+"^", kmQuestioned)
		child := durableBlobAt(root, c.hash, kmQuestioned)
		rooms[kmQuestioned] = durableNewRoomRef(parent, child, kmQuestioned)
		reworked := durableQuestionedBatchRework(parent, child)
		if started < len(expected) && !reworked || rooms[kmQuestioned] == "" && !reworked {
			return false
		}
		allowed[kmQuestioned+".md"] = true
	}
	for _, path := range c.files {
		if allowed[path] {
			continue
		}
		owned := false
		for _, room := range rooms {
			owned = owned || (c.added[path] && strings.HasPrefix(path, room+"/"))
		}
		if !owned {
			return false
		}
	}
	return true
}
func durableQuestionedBatchRework(parent, child string) bool {
	return durableField(parent, "status") == "review" && durableField(parent, "started") == "" &&
		durableField(child, "status") == "ideation" && durableField(child, "started") != "" &&
		!durableTerminalState(child)
}
func durableBatchTerminalProof(root string, c durableCommit, expected map[string]string) bool {
	if len(expected) < 2 || len(c.files) != len(expected) {
		return false
	}
	changed := map[string]bool{}
	for _, path := range c.files {
		changed[path] = true
	}
	for slug := range expected {
		blob := durableBlobAt(root, c.hash, slug)
		if !changed[slug+".md"] || durableField(blob, "status") != "done" ||
			durableField(blob, "completed") == "" || durableField(blob, "verdict") == "" {
			return false
		}
	}
	return true
}
func durableRoomRef(content string) string {
	room := ""
	for _, line := range strings.Split(content, "\n") {
		if line = strings.TrimSpace(line); strings.HasPrefix(line, "room-ref:") {
			room = strings.TrimSpace(strings.TrimPrefix(line, "room-ref:"))
		}
	}
	return strings.TrimPrefix(room, "./")
}
func durableAtomicWorkerProof(history []durableCommit, i int, slug, stage string) bool {
	if i == 0 || (!history[i].entityFileScoped && !durableNewRoomScope(history[i-1].blob, history[i], slug)) {
		return false
	}
	report := "\n## Stage Report: " + stage + "\n"
	parent, child := history[i-1].blob, history[i].blob
	return durableField(parent, "started") == "" && !strings.Contains(parent, report) &&
		durableField(child, "status") == stage && durableField(child, "started") != "" &&
		strings.Contains(child, report)
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
		changes := strings.Split(strings.TrimSpace(git(t, root, "diff-tree", "--root", "--no-commit-id", "--name-status", "-r", "-M", fields[0])), "\n")
		files, added := []string{}, map[string]bool{}
		for _, change := range changes {
			parts := strings.Split(change, "\t")
			if len(parts) < 2 {
				continue
			}
			path := parts[len(parts)-1]
			files = append(files, path)
			added[path] = parts[0] == "A"
		}
		scoped, entityOwned := len(files) > 0, len(files) > 0
		for _, file := range files {
			if file != slug+".md" && file != filepath.Join("_archive", slug+".md") {
				scoped = false
			}
			if !durablePathOwnedBySlug(slug, file) {
				entityOwned = false
			}
		}
		history = append(history, durableCommit{fields[0], fields[1], blob, scoped, entityOwned, files, added})
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
func durableTerminalState(content string) bool {
	return durableField(content, "status") == "done" ||
		durableField(content, "completed") != "" || durableField(content, "verdict") != ""
}
func assertDurableKeepMoving(t *testing.T, root string) error {
	completed, failures := gradeDurableTaskJourneys(t, root, kmExpected(), kmExpected())
	if completed != 3 {
		return fmt.Errorf("durable keep-moving journeys = %d/3: %v", completed, failures)
	}
	if err := durableIndependentOverlap(t, root, kmExpected()); err != nil {
		return err
	}
	content, where, found := locateEntity(root, kmQuestioned)
	if !found || where != filepath.Join(root, kmQuestioned+".md") || durableTerminalState(content) {
		return fmt.Errorf("%s must remain active and nonterminal", kmQuestioned)
	}
	history := durableEntityHistory(t, root, kmQuestioned, kmQuestioned+".md")
	transitioned, reported, terminal := false, false, durableTerminalState(history[0].blob)
	for i := 1; i < len(history); i++ {
		transitioned = transitioned || durableField(history[i].blob, "status") != durableField(history[0].blob, "status")
		reported = reported || (history[i].entityFileScoped &&
			strings.Count(history[i].blob, "\n## Stage Report: ") > strings.Count(history[i-1].blob, "\n## Stage Report: "))
		terminal = terminal || durableTerminalState(history[i].blob)
	}
	if terminal {
		return fmt.Errorf("%s must remain historically active and nonterminal", kmQuestioned)
	}
	if !transitioned || !reported || !history[len(history)-1].entityOwned {
		return fmt.Errorf("%s has no meaningful entity-owned stage re-shape", kmQuestioned)
	}
	return nil
}
func durableIndependentOverlap(t *testing.T, root string, expected map[string]string) error {
	rank := map[string]int{}
	for i, hash := range strings.Fields(git(t, root, "log", "--reverse", "--format=%H")) {
		rank[hash] = i
	}
	latestEngagement, firstTerminal := -1, len(rank)
	proofs := map[string]durableJourneyProof{}
	for slug, stage := range expected {
		proof, reason := durableTaskJourney(t, root, slug, stage, expected)
		if reason != "" {
			return fmt.Errorf("%s: %s", slug, reason)
		}
		proofs[slug] = proof
		if rank[proof.engaged] > latestEngagement {
			latestEngagement = rank[proof.engaged]
		}
		if rank[proof.terminal] < firstTerminal {
			firstTerminal = rank[proof.terminal]
		}
	}
	if latestEngagement >= firstTerminal &&
		!durableDelayedPersistenceOverlap(root, proofs, rank, firstTerminal) {
		return fmt.Errorf("independent task journeys do not overlap before terminalization")
	}
	return nil
}
func durableDelayedPersistenceOverlap(root string, proofs map[string]durableJourneyProof, rank map[string]int, firstTerminal int) bool {
	starts := map[string]time.Time{}
	var frontier, earliestTerminal time.Time
	engagedBeforeTerminal := 0
	for slug, proof := range proofs {
		started, err := time.Parse(time.RFC3339Nano, durableField(durableBlobAt(root, proof.engaged, slug), "started"))
		if err != nil {
			return false
		}
		completed, err := time.Parse(time.RFC3339Nano, durableField(durableBlobAt(root, proof.terminal, slug), "completed"))
		if err != nil {
			return false
		}
		starts[slug] = started
		if earliestTerminal.IsZero() || completed.Before(earliestTerminal) {
			earliestTerminal = completed
		}
		if rank[proof.engaged] < firstTerminal {
			engagedBeforeTerminal++
			if frontier.IsZero() || started.After(frontier) {
				frontier = started
			}
		}
	}
	if engagedBeforeTerminal < 2 {
		return false
	}
	for _, started := range starts {
		if started.After(frontier) || !started.Before(earliestTerminal) {
			return false
		}
	}
	return true
}
func assertDurableSmallestMechanism(t *testing.T, root string, tr mechanismTrace, edits, commissioned []string) error {
	expected := map[string]string{}
	for _, slug := range commissioned {
		expected[slug] = "ready"
	}
	completed, failures := gradeDurableTaskJourneys(t, root, expected, nil)
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
		{"archive reverts terminal fields", "archive-reverts-terminal", "ready-one", "terminal fields"},
		{"missing archive", "missing-archive", "ready-one", "canonical archive"},
		{"report before dispatch", "report-before-dispatch", "ready-one", "worker report after dispatch"},
		{"terminal before report", "terminal-before-report", "ready-one", "worker report"},
		{"terminal status before report", "terminal-status-before-report", "ready-one", "worker report"},
		{"completed before report", "terminal-completed-before-report", "ready-one", "worker report"},
		{"verdict before report", "terminal-verdict-before-report", "ready-one", "worker report"},
		{"cross-attributed report", "cross-attributed-report", "ready-one", "path-scoped worker report"},
		{"cross-attributed archive", "cross-attributed-archive", "ready-one", "canonical archive"},
		{"slug-prefix archive", "slug-prefix-archive", "ready-one", "canonical archive"},
		{"atomic first worker", "atomic-worker", "", ""},
		{"atomic preexisting started", "atomic-preexisting-started", "ready-one", "dispatch entry"},
		{"atomic preexisting report", "atomic-preexisting-report", "ready-one", "dispatch entry"},
		{"atomic report only", "atomic-report-only", "ready-one", "dispatch entry"},
		{"atomic foreign path", "atomic-foreign-path", "ready-one", "dispatch entry"},
		{"atomic slug prefix", "atomic-slug-prefix", "ready-one", "dispatch entry"},
		{"atomic worker gate room", "atomic-room", "", ""},
		{"atomic worker preexisting room", "atomic-room-preexisting", "ready-one", "dispatch entry"},
		{"atomic worker replaced room", "atomic-room-replaced", "ready-one", "dispatch entry"},
		{"atomic worker modified room file", "atomic-room-modified", "ready-one", "dispatch entry"},
		{"atomic worker path outside room", "atomic-room-outside", "ready-one", "dispatch entry"},
		{"atomic worker slug prefix room", "atomic-room-prefix", "ready-one", "dispatch entry"},
		{"dispatch gate room", "dispatch-room", "", ""},
		{"dispatch preexisting room", "dispatch-room-preexisting", "ready-one", "dispatch entry"},
		{"dispatch preexisting different room", "dispatch-room-preexisting-different", "ready-one", "dispatch entry"},
		{"dispatch outside slug room", "dispatch-room-outside-slug", "ready-one", "dispatch entry"},
		{"dispatch sibling outside room", "dispatch-room-sibling", "ready-one", "dispatch entry"},
		{"dispatch modified room file", "dispatch-room-modified", "ready-one", "dispatch entry"},
		{"dispatch slug prefix room", "dispatch-room-prefix", "ready-one", "dispatch entry"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			completed, failures := gradeDurableTaskJourneys(t, durableJourneyFixture(t, tt.mutation), kmExpected(), nil)
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
	if mutation == "atomic-preexisting-started" {
		path := filepath.Join(root, "ready-one.md")
		writeFile(t, path, strings.Replace(readFile(t, path), "started: ", "started: 2026-07-30T23:59:00Z", 1))
	}
	if mutation == "atomic-preexisting-report" {
		durableAppendReport(t, root, "ready-one")
	}
	room := "ready-one/review/review/briefing-1"
	roomFile := room + "/index.json"
	if mutation == "atomic-room-preexisting" {
		path := filepath.Join(root, "ready-one.md")
		writeFile(t, path, durableBindRoom(readFile(t, path), room))
	}
	if mutation == "atomic-room-replaced" {
		path := filepath.Join(root, "ready-one.md")
		writeFile(t, path, durableBindRoom(readFile(t, path), "ready-one/review/review/old-briefing"))
	}
	if mutation == "atomic-room-modified" {
		writeFile(t, filepath.Join(root, roomFile), "old\n")
	}
	if mutation == "dispatch-room-preexisting" {
		path := filepath.Join(root, "ready-one.md")
		writeFile(t, path, durableBindRoom(readFile(t, path), room))
	}
	if mutation == "dispatch-room-preexisting-different" {
		path := filepath.Join(root, "ready-one.md")
		writeFile(t, path, durableBindRoom(readFile(t, path), "ready-one/review/review/old-briefing"))
	}
	if mutation == "dispatch-room-modified" {
		writeFile(t, filepath.Join(root, roomFile), "old\n")
	}
	gitInit(t, root)
	if mutation == "parallel-dispatch" {
		for _, slug := range []string{kmApprovedGate, kmReadyOne, kmReadyTwo} {
			path := filepath.Join(root, slug+".md")
			writeFile(t, path, strings.Replace(readFile(t, path), "started: ", "started: now", 1))
			gitCommitPathScoped(t, root, slug+".md", "dispatch: "+slug+" entering "+kmNextStage)
		}
	}
	for _, slug := range []string{"approved-gate", "ready-one", "ready-two"} {
		atomic := strings.HasPrefix(mutation, "atomic-") && (mutation == "atomic-worker" || slug == "ready-one")
		roomDispatch := strings.HasPrefix(mutation, "dispatch-room") && slug == "ready-one"
		if mutation == "report-before-dispatch" && slug == "ready-one" {
			durableAppendReport(t, root, slug)
			gitCommitPathScoped(t, root, slug+".md", "worker: stale "+slug)
		}
		if mutation != "parallel-dispatch" && !atomic && (mutation != "missing-dispatch" || slug != "ready-one") {
			path := filepath.Join(root, slug+".md")
			content := strings.Replace(readFile(t, path), "started: ", "started: 2026-07-31T00:00:00Z", 1)
			if roomDispatch && mutation != "dispatch-room-preexisting" {
				if mutation == "dispatch-room-outside-slug" {
					room = "ready-two/review/review/briefing-1"
				} else if mutation == "dispatch-room-prefix" {
					room = "ready-one-other/review/review/briefing-1"
				}
				content = durableBindRoom(content, room)
			}
			writeFile(t, path, content)
			if roomDispatch {
				if mutation == "dispatch-room-sibling" {
					roomFile = "ready-one/review/review/briefing-2/index.json"
				} else {
					roomFile = room + "/index.json"
				}
				writeFile(t, filepath.Join(root, roomFile), "new\n")
				git(t, root, "add", "--", slug+".md", roomFile)
				git(t, root, "commit", "-q", "-m", "dispatch: "+slug+" entering implementation", "--", slug+".md", roomFile)
			} else {
				gitCommitPathScoped(t, root, slug+".md", "dispatch: "+slug+" entering implementation")
			}
		}
		if strings.Contains(mutation, "before-report") && slug == "ready-one" {
			path := filepath.Join(root, slug+".md")
			content := readFile(t, path)
			if mutation == "terminal-before-report" || mutation == "terminal-status-before-report" {
				content = strings.Replace(content, "status: implementation", "status: done", 1)
			}
			if mutation == "terminal-before-report" || mutation == "terminal-completed-before-report" {
				content = strings.Replace(content, "completed:", "completed: premature", 1)
			}
			if mutation == "terminal-before-report" || mutation == "terminal-verdict-before-report" {
				content = strings.Replace(content, "verdict:", "verdict: passed", 1)
			}
			writeFile(t, path, content)
			gitCommitPathScoped(t, root, slug+".md", "premature terminal: "+slug)
		}
		if (mutation != "missing-report" || slug != "ready-one") && !(mutation == "report-before-dispatch" && slug == "ready-one") {
			path := filepath.Join(root, slug+".md")
			if atomic && mutation != "atomic-report-only" {
				writeFile(t, path, strings.Replace(readFile(t, path), "started: ", "started: 2026-07-31T00:00:00Z", 1))
			}
			if !(atomic && mutation == "atomic-preexisting-report") {
				durableAppendReport(t, root, slug)
			}
			if strings.HasPrefix(mutation, "atomic-room") && slug == "ready-one" {
				path := filepath.Join(root, slug+".md")
				workerRoom := room
				if mutation == "atomic-room-prefix" {
					workerRoom = "ready-one-other/review/review/briefing-1"
				}
				if mutation != "atomic-room-preexisting" {
					writeFile(t, path, durableBindRoom(readFile(t, path), workerRoom))
				}
				workerRoomFile := workerRoom + "/index.json"
				if mutation == "atomic-room-outside" {
					workerRoomFile = "ready-one/review/review/briefing-2/index.json"
				}
				writeFile(t, filepath.Join(root, workerRoomFile), "new\n")
				git(t, root, "add", "--", slug+".md", workerRoomFile)
				git(t, root, "commit", "-q", "-m", "worker: "+slug, "--", slug+".md", workerRoomFile)
			} else if (mutation == "cross-attributed-report" || mutation == "atomic-foreign-path") && slug == "ready-one" {
				writeFile(t, filepath.Join(root, "ready-two.md"), readFile(t, filepath.Join(root, "ready-two.md"))+"\ncross-attributed\n")
				git(t, root, "add", "--", "ready-one.md", "ready-two.md")
				git(t, root, "commit", "-q", "-m", "worker: wrong scope", "--", "ready-one.md", "ready-two.md")
			} else if mutation == "atomic-slug-prefix" && slug == "ready-one" {
				foreign := "ready-one-other/review/briefing.json"
				writeFile(t, filepath.Join(root, foreign), "{}\n")
				git(t, root, "add", "--", "ready-one.md", foreign)
				git(t, root, "commit", "-q", "-m", "worker: wrong prefix", "--", "ready-one.md", foreign)
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
		if mutation != "atomic-worker" {
			writeFile(t, filepath.Join(root, slug+".md"), content)
			gitCommitPathScoped(t, root, slug+".md", "terminalize: "+slug)
		}
		if mutation != "missing-archive" || slug != "ready-one" {
			if mutation == "archive-reverts-terminal" && slug == "ready-one" {
				content = strings.Replace(content, "status: done", "status: implementation", 1)
				content = strings.Replace(content, "completed: "+completed, "completed:", 1)
				content = strings.Replace(content, "verdict: "+verdict, "verdict:", 1)
			}
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
func TestDurableKeepMovingRequiresOverlappingJourneys(t *testing.T) {
	if err := assertDurableKeepMoving(t, durableAddQuestioned(t, durableJourneyFixture(t, "parallel-dispatch"), true)); err != nil {
		t.Fatalf("parallel journey: %v", err)
	}
	if err := assertDurableKeepMoving(t, durableAddQuestioned(t, durableJourneyFixture(t, ""), true)); err == nil ||
		!strings.Contains(err.Error(), "do not overlap") {
		t.Fatalf("serialized journey error = %v, want overlap failure", err)
	}
}
func TestDurableKeepMovingDelayedPersistence(t *testing.T) {
	tests := []struct {
		name, mutation string
		wantPass       bool
	}{
		{"corroborated frontier", "", true},
		{"absent start", "absent-start", false},
		{"unparseable start", "unparseable-start", false},
		{"start at first terminal", "start-at-terminal", false},
		{"start after first terminal", "start-after-terminal", false},
		{"no corroborating frontier", "no-frontier", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := assertDurableKeepMoving(t, durableDelayedPersistenceFixture(t, tt.mutation))
			if tt.wantPass && err != nil {
				t.Fatal(err)
			}
			if !tt.wantPass && err == nil {
				t.Fatal("delayed persistence passed, want failure")
			}
		})
	}
}
func durableDelayedPersistenceFixture(t *testing.T, mutation string) string {
	root := t.TempDir()
	slugs := []string{kmApprovedGate, kmReadyOne, kmReadyTwo}
	for _, slug := range slugs {
		writeFile(t, filepath.Join(root, slug+".md"), durableEntity(slug, kmNextStage, "", ""))
	}
	gitInit(t, root)
	report := func(slug, started string) {
		path := filepath.Join(root, slug+".md")
		if started != "" {
			writeFile(t, path, strings.Replace(readFile(t, path), "started: ", "started: "+started, 1))
		}
		durableAppendReport(t, root, slug)
		gitCommitPathScoped(t, root, slug+".md", "worker: "+slug)
	}
	archive := func(slug, completed string) {
		path := filepath.Join(root, slug+".md")
		content := strings.Replace(readFile(t, path), "status: implementation", "status: done", 1)
		content = strings.Replace(content, "completed:", "completed: "+completed, 1)
		content = strings.Replace(content, "verdict:", "verdict: passed", 1)
		writeFile(t, filepath.Join(root, "_archive", slug+".md"), content)
		git(t, root, "rm", "-q", "--", slug+".md")
		git(t, root, "add", "--", "_archive/"+slug+".md")
		git(t, root, "commit", "-q", "-m", "archive: "+slug)
	}
	report(kmApprovedGate, "2026-07-31T00:00:00Z")
	if mutation != "no-frontier" {
		report(kmReadyOne, "2026-07-31T00:00:00Z")
	}
	archive(kmApprovedGate, "2026-07-31T00:01:00Z")
	if mutation == "no-frontier" {
		report(kmReadyOne, "2026-07-31T00:00:00Z")
	}
	started := "2026-07-31T00:00:00Z"
	if mutation == "absent-start" {
		started = ""
	} else if mutation == "unparseable-start" {
		started = "not-a-time"
	} else if mutation == "start-at-terminal" {
		started = "2026-07-31T00:01:00Z"
	} else if mutation == "start-after-terminal" {
		started = "2026-07-31T00:01:01Z"
	}
	report(kmReadyTwo, started)
	archive(kmReadyOne, "2026-07-31T00:02:00Z")
	archive(kmReadyTwo, "2026-07-31T00:02:00Z")
	return durableAddQuestioned(t, root, true)
}
func TestDurableQuestionedRejectsUnrelatedEdit(t *testing.T) {
	if err := assertDurableKeepMoving(t, durableAddQuestioned(t, durableJourneyFixture(t, "parallel-dispatch"), false)); err == nil ||
		!strings.Contains(err.Error(), "meaningful") {
		t.Fatalf("unrelated questioned edit error = %v, want meaningful re-shape failure", err)
	}
}
func TestDurableQuestionedRejectsTerminalHistory(t *testing.T) {
	for _, field := range []string{"all", "status", "completed", "verdict"} {
		t.Run(field, func(t *testing.T) {
			root := durableAddQuestionedTerminalHistory(t, durableJourneyFixture(t, "parallel-dispatch"), field)
			if err := assertDurableKeepMoving(t, root); err == nil || !strings.Contains(err.Error(), "nonterminal") {
				t.Fatalf("questioned terminal history error = %v, want nonterminal failure", err)
			}
		})
	}
}
func durableAddQuestionedTerminalHistory(t *testing.T, root, field string) string {
	path := filepath.Join(root, kmQuestioned+".md")
	writeFile(t, path, durableEntity(kmQuestioned, "review", "", ""))
	gitCommitPathScoped(t, root, kmQuestioned+".md", "seed questioned")
	content := readFile(t, path)
	if field == "all" || field == "status" {
		content = strings.Replace(content, "status: review", "status: done", 1)
	}
	if field == "all" || field == "completed" {
		content = strings.Replace(content, "completed:", "completed: premature", 1)
	}
	if field == "all" || field == "verdict" {
		content = strings.Replace(content, "verdict:", "verdict: passed", 1)
	}
	writeFile(t, path, content)
	gitCommitPathScoped(t, root, kmQuestioned+".md", "premature terminal questioned")
	content = strings.Replace(content, "status: done", "status: ideation", 1)
	content = strings.Replace(content, "status: review", "status: ideation", 1)
	content = strings.Replace(content, "completed: premature", "completed:", 1)
	content = strings.Replace(content, "verdict: passed", "verdict:", 1)
	writeFile(t, path, content)
	durableAppendStageReport(t, root, kmQuestioned, "ideation")
	gitCommitPathScoped(t, root, kmQuestioned+".md", "reopen questioned")
	return root
}
func TestDurableKeepMovingBatchMotion(t *testing.T) {
	tests := []struct {
		name, mutation, reason string
	}{
		{"full set", "", ""},
		{"split batch with questioned rework", "split", ""},
		{"split batch unrelated questioned", "split-questioned-unrelated", "dispatch"},
		{"partial set", "partial", "dispatch"},
		{"foreign path", "foreign", "dispatch"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := durableBatchKeepMovingFixture(t, tt.mutation)
			err := assertDurableKeepMoving(t, root)
			if tt.reason == "" && err != nil {
				t.Fatal(err)
			}
			if tt.reason != "" && (err == nil || !strings.Contains(err.Error(), tt.reason)) {
				t.Fatalf("batch motion error = %v, want reason containing %q", err, tt.reason)
			}
			if tt.mutation == "partial" {
				completed, failures := gradeDurableTaskJourneys(t, root, kmExpected(), kmExpected())
				if completed != 2 || len(failures) != 1 || !strings.Contains(failures[kmReadyTwo], "dispatch") {
					t.Fatalf("partial batch = %d/3, failures=%v; want only omitted %s red", completed, failures, kmReadyTwo)
				}
			}
		})
	}
}
func durableBatchKeepMovingFixture(t *testing.T, mutation string) string {
	root := t.TempDir()
	slugs := []string{kmApprovedGate, kmReadyOne, kmReadyTwo}
	for _, slug := range slugs {
		writeFile(t, filepath.Join(root, slug+".md"), durableEntity(slug, kmNextStage, "", ""))
	}
	split := strings.HasPrefix(mutation, "split")
	if split {
		writeFile(t, filepath.Join(root, kmQuestioned+".md"), durableEntity(kmQuestioned, "review", "", ""))
	}
	gitInit(t, root)
	if split {
		path := filepath.Join(root, kmApprovedGate+".md")
		writeFile(t, path, strings.Replace(readFile(t, path), "started: ", "started: now", 1))
		gitCommitPathScoped(t, root, kmApprovedGate+".md", "dispatch: "+kmApprovedGate+" entering "+kmNextStage)
	}
	var paths []string
	for i, slug := range slugs {
		if mutation == "partial" && i == len(slugs)-1 {
			continue
		}
		if split && slug == kmApprovedGate {
			continue
		}
		path := slug + ".md"
		writeFile(t, filepath.Join(root, path), strings.Replace(readFile(t, filepath.Join(root, path)), "started: ", "started: now", 1))
		paths = append(paths, path)
	}
	if split {
		path := filepath.Join(root, kmQuestioned+".md")
		content := readFile(t, path)
		if mutation == "split" {
			content = strings.Replace(content, "status: review", "status: ideation", 1)
			content = strings.Replace(content, "started: ", "started: now", 1)
		} else {
			content += "\nunrelated note\n"
		}
		writeFile(t, path, content)
		paths = append(paths, kmQuestioned+".md")
	}
	if mutation == "foreign" {
		writeFile(t, filepath.Join(root, "foreign.md"), "foreign\n")
		paths = append(paths, "foreign.md")
	}
	git(t, root, append([]string{"add", "--"}, paths...)...)
	git(t, root, append([]string{"commit", "-q", "-m", "dispatch: implementation batch", "--"}, paths...)...)
	for _, slug := range slugs {
		durableAppendReport(t, root, slug)
		gitCommitPathScoped(t, root, slug+".md", "worker: "+slug)
	}
	if split {
		durableAppendStageReport(t, root, kmQuestioned, "ideation")
		gitCommitPathScoped(t, root, kmQuestioned+".md", "reshape questioned")
	}
	paths = nil
	for _, slug := range slugs {
		path := filepath.Join(root, slug+".md")
		content := strings.Replace(readFile(t, path), "status: "+kmNextStage, "status: done", 1)
		content = strings.Replace(content, "completed:", "completed: now", 1)
		content = strings.Replace(content, "verdict:", "verdict: passed", 1)
		writeFile(t, path, content)
		paths = append(paths, slug+".md")
	}
	git(t, root, append([]string{"add", "--"}, paths...)...)
	git(t, root, append([]string{"commit", "-q", "-m", "terminalize: implementation batch", "--"}, paths...)...)
	for _, slug := range slugs {
		content := readFile(t, filepath.Join(root, slug+".md"))
		writeFile(t, filepath.Join(root, "_archive", slug+".md"), content)
		git(t, root, "rm", "-q", "--", slug+".md")
		git(t, root, "add", "--", "_archive/"+slug+".md")
		git(t, root, "commit", "-q", "-m", "archive: "+slug)
	}
	if split {
		return root
	}
	return durableAddQuestioned(t, root, true)
}
func durableAddQuestioned(t *testing.T, root string, reshape bool) string {
	writeFile(t, filepath.Join(root, kmQuestioned+".md"), durableEntity(kmQuestioned, "review", "", ""))
	gitCommitPathScoped(t, root, kmQuestioned+".md", "seed questioned")
	path := filepath.Join(root, kmQuestioned+".md")
	if reshape {
		writeFile(t, path, strings.Replace(readFile(t, path), "status: review", "status: ideation", 1))
		durableAppendStageReport(t, root, kmQuestioned, "ideation")
	} else {
		writeFile(t, path, readFile(t, path)+"\nunrelated note\n")
	}
	gitCommitPathScoped(t, root, kmQuestioned+".md", "reshape questioned")
	return root
}
func TestRetainedAtomicWorkerJourney(t *testing.T) {
	root := os.Getenv("SPACEDOCK_KEEP_MOVING_RETAIN_ROOT")
	if root == "" {
		t.Skip("SPACEDOCK_KEEP_MOVING_RETAIN_ROOT is not set")
	}
	if err := assertDurableKeepMoving(t, root); err != nil {
		t.Fatal(err)
	}
}
func durableEntity(slug, stage, started, report string) string {
	return "---\nid: " + slug + "\ntitle: " + slug + "\nstatus: " + stage +
		"\nstarted: " + started + "\ncompleted:\nverdict:\n---\n# " + slug + "\n" + report
}
func durableBindRoom(content, room string) string {
	return strings.Replace(content, "\n---\n# ", "\nroom-ref: ./"+room+"\n---\n# ", 1)
}
func durableAppendReport(t *testing.T, root, slug string) {
	durableAppendStageReport(t, root, slug, "implementation")
}
func durableAppendStageReport(t *testing.T, root, slug, stage string) {
	path := filepath.Join(root, slug+".md")
	writeFile(t, path, readFile(t, path)+"\n## Stage Report: "+stage+"\n\n- DONE: complete\n  durable evidence\n\n### Summary\n\nDone.\n")
}
