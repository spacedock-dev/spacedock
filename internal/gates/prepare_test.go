package gates

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
	"gopkg.in/yaml.v3"
)

func TestPrepareRequiresActionableCurrentStage(t *testing.T) {
	for stage, allow := range map[string]bool{"validation": true, "implementation": false, "done": false, "contradictory": false} {
		t.Run(stage, func(t *testing.T) {
			workflow, state, entity, artifact, _ := prepareFixture(t, "flat")
			body, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			body = bytes.Replace(body, []byte("status: validation"), []byte("status: "+stage), 1)
			if err := os.WriteFile(entity, body, 0o644); err != nil {
				t.Fatal(err)
			}
			review := filepath.Join(state, "task", "review")
			beforeReview := prepareTreeSnapshot(t, review)

			input := PrepareInput{WorkflowDir: workflow, Question: "Review?", Artifact: artifact, Summary: "summary"}
			result, err := Prepare(entity, input)
			if allow {
				if err != nil || result.State != "open" {
					t.Fatalf("actionable stage result=%#v error=%v", result, err)
				}
				if err := RecordSemantic(entity, RecordInput{Decision: "hold", Actor: "person:captain", Reason: "wait", WorkflowDir: workflow}); err != nil {
					t.Fatal(err)
				}
				successor, err := Prepare(entity, input)
				if err != nil {
					t.Fatal(err)
				}
				doc, _, err := Read(entity)
				if err != nil {
					t.Fatal(err)
				}
				if successor.Briefing != "briefing:task:validation:attempt-2:revision-1" || filepath.Base(successor.Room) != "briefing-2" || len(doc.Records) != 1 || len(doc.Records[0].Attempts) != 2 || doc.Records[0].Attempts[0].Briefing.ID != result.Briefing || doc.Records[0].Attempts[0].Resolution == nil || doc.Records[0].Attempts[1].ID != "gate-attempt:task-validation-2" {
					t.Fatalf("invalid retained successor: result=%#v records=%#v", successor, doc.Records)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), stage) || !strings.Contains(err.Error(), "is not an actionable gate") {
				t.Fatalf("stage %q error=%v", stage, err)
			}
			afterEntity, readErr := os.ReadFile(entity)
			if readErr != nil {
				t.Fatal(readErr)
			}
			if !bytes.Equal(afterEntity, body) || prepareTreeSnapshot(t, review) != beforeReview {
				t.Fatalf("stage %q refusal changed entity bytes or review tree", stage)
			}
		})
	}
}

func prepareTreeSnapshot(t *testing.T, root string) string {
	t.Helper()
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return ""
	}
	var entryTypes strings.Builder
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		entryTypes.WriteString(rel + "\x00" + entry.Type().String() + "\x00")
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return entryTypes.String() + treeDigest(t, root)
}

func TestPrepareCreatesOneTwoFileRecorderRoomForFolderAndFlatEntities(t *testing.T) {
	for _, form := range []string{"folder", "flat"} {
		t.Run(form, func(t *testing.T) {
			workflow, state, entity, artifact, reference := prepareFixture(t, form)
			summary := "  Résumé — validates Git-root presentation exactly.  "
			result, err := Prepare(entity, PrepareInput{
				WorkflowDir: workflow,
				Question:    "Should this gate advance?",
				Artifact:    artifact,
				Summary:     summary,
				References:  []string{reference},
			})
			if err != nil {
				t.Fatal(err)
			}
			slug := "task"
			wantRoom := filepath.Join(state, slug, "review", "validation", "briefing-1")
			if result.Room != wantRoom || result.State != "open" ||
				result.Briefing != "briefing:task:validation:attempt-1:revision-1" ||
				!digestRE.MatchString(result.Digest) {
				t.Fatalf("unexpected result: %#v want room=%s", result, wantRoom)
			}
			entries, err := os.ReadDir(result.Room)
			if err != nil {
				t.Fatal(err)
			}
			var names []string
			for _, entry := range entries {
				if !entry.Type().IsRegular() {
					t.Fatalf("prepare-time room contains non-regular entry %s", entry.Name())
				}
				names = append(names, entry.Name())
			}
			if want := []string{"gate-briefing.json", "request.json"}; !reflect.DeepEqual(names, want) {
				t.Fatalf("room files=%v want %v", names, want)
			}
			for _, copied := range []string{filepath.Base(artifact), filepath.Base(reference), "association.json", "provider"} {
				if _, err := os.Stat(filepath.Join(result.Room, copied)); !os.IsNotExist(err) {
					t.Fatalf("prepare copied or invented %s: %v", copied, err)
				}
			}

			briefingBytes, err := os.ReadFile(filepath.Join(result.Room, "gate-briefing.json"))
			if err != nil {
				t.Fatal(err)
			}
			var briefing struct {
				ID        string `json:"id"`
				Artifacts []struct {
					ID        string  `json:"id"`
					URI       string  `json:"uri"`
					Rev       string  `json:"rev"`
					MediaType string  `json:"mediaType"`
					Summary   *string `json:"summary"`
				} `json:"artifacts"`
				Context []struct {
					Type      string  `json:"type"`
					ID        string  `json:"id"`
					URI       string  `json:"uri"`
					Rev       string  `json:"rev"`
					MediaType string  `json:"mediaType"`
					Summary   *string `json:"summary"`
				} `json:"context"`
			}
			if err := json.Unmarshal(briefingBytes, &briefing); err != nil {
				t.Fatal(err)
			}
			if briefing.ID != result.Briefing || len(briefing.Artifacts) != 1 || len(briefing.Context) != 1 {
				t.Fatalf("unexpected Briefing: %#v", briefing)
			}
			if briefing.Artifacts[0].Summary == nil || *briefing.Artifacts[0].Summary != summary {
				t.Fatalf("primary summary=%v want exact %q", briefing.Artifacts[0].Summary, summary)
			}
			if briefing.Context[0].Summary != nil || briefing.Context[0].Type != "Reference" {
				t.Fatalf("Reference acquired summary or wrong type: %#v", briefing.Context[0])
			}
			if !strings.HasPrefix(briefing.Artifacts[0].URI, "git-root://main/") ||
				!strings.HasPrefix(briefing.Context[0].URI, "git-root://state/") {
				t.Fatalf("logical roots not classified: artifact=%q reference=%q", briefing.Artifacts[0].URI, briefing.Context[0].URI)
			}
			if briefing.Artifacts[0].MediaType != "text/markdown" || briefing.Context[0].MediaType != "application/json" {
				t.Fatalf("media types artifact=%q reference=%q", briefing.Artifacts[0].MediaType, briefing.Context[0].MediaType)
			}

			doc, _, err := Read(entity)
			if err != nil {
				t.Fatal(err)
			}
			current := CurrentSummary(doc)
			if current.State != "open" || current.Briefing != result.Briefing {
				t.Fatalf("prepared binding not open: %#v", current)
			}
			attempt := doc.Records[0].Attempts[0]
			if attempt.Briefing.Digest != result.Digest || attempt.Briefing.RequestDigest == "" {
				t.Fatalf("binding pins incomplete: %#v", attempt.Briefing)
			}
			wantRef := "./task/review/validation/briefing-1"
			if form == "folder" {
				wantRef = "./review/validation/briefing-1"
			}
			if attempt.Briefing.RoomRef != wantRef {
				t.Fatalf("room-ref=%q want %q", attempt.Briefing.RoomRef, wantRef)
			}

			beforeEntity, _ := os.ReadFile(entity)
			beforeBriefing := append([]byte(nil), briefingBytes...)
			replayed, err := Prepare(entity, PrepareInput{
				WorkflowDir: workflow,
				Question:    "Should this gate advance?",
				Artifact:    artifact,
				Summary:     summary,
				References:  []string{reference},
			})
			if err != nil || replayed != result {
				t.Fatalf("replay=%#v err=%v want %#v", replayed, err, result)
			}
			afterEntity, _ := os.ReadFile(entity)
			afterBriefing, _ := os.ReadFile(filepath.Join(result.Room, "gate-briefing.json"))
			if !bytes.Equal(beforeEntity, afterEntity) || !bytes.Equal(beforeBriefing, afterBriefing) {
				t.Fatal("exact replay changed durable bytes")
			}
		})
	}
}

func TestWithdrawPreparedAttemptThenPrepareAppendsSuccessorWithoutRewritingOldRoom(t *testing.T) {
	workflow, _, entity, artifact, reference := prepareFixture(t, "folder")
	first, err := Prepare(entity, PrepareInput{
		WorkflowDir: workflow,
		Question:    "Advance?",
		Artifact:    artifact,
		Summary:     "first candidate",
		References:  []string{reference},
	})
	if err != nil {
		t.Fatal(err)
	}
	firstBriefing := readPreparedFile(t, filepath.Join(first.Room, preparedBriefingLocator))
	firstRequest := readPreparedFile(t, filepath.Join(first.Room, "request.json"))

	summary, err := Withdraw(entity, WithdrawInput{
		WorkflowDir: workflow,
		Reason:      "Sprint re-scope replaced the reviewed candidate.",
	})
	if err != nil {
		t.Fatal(err)
	}
	if summary.State != "withdrawn" || summary.Attempt != "gate-attempt:task-validation-1" ||
		summary.Resolution != "" || summary.Decision != "" || summary.Application != "" {
		t.Fatalf("withdraw summary = %#v", summary)
	}

	second, err := Prepare(entity, PrepareInput{
		WorkflowDir: workflow,
		Question:    "Advance the replacement?",
		Artifact:    artifact,
		Summary:     "replacement candidate",
		References:  []string{reference},
	})
	if err != nil {
		t.Fatal(err)
	}
	if second.Room == first.Room || second.Briefing != "briefing:task:validation:attempt-2:revision-1" {
		t.Fatalf("successor result = %#v, first=%#v", second, first)
	}
	if got := readPreparedFile(t, filepath.Join(first.Room, preparedBriefingLocator)); !bytes.Equal(got, firstBriefing) {
		t.Fatal("successor prepare rewrote withdrawn Briefing bytes")
	}
	if got := readPreparedFile(t, filepath.Join(first.Room, "request.json")); !bytes.Equal(got, firstRequest) {
		t.Fatal("successor prepare rewrote withdrawn request bytes")
	}
	doc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Records[0].Attempts) != 2 || attemptState(&doc.Records[0].Attempts[0]) != "withdrawn" ||
		attemptState(&doc.Records[0].Attempts[1]) != "open" {
		t.Fatalf("replacement history = %#v", doc.Records[0].Attempts)
	}
}

func TestWithdrawRefusalsLeaveEntityRoomAndLockBytesClean(t *testing.T) {
	t.Run("blank reason", func(t *testing.T) {
		workflow, state, entity, artifact, _ := prepareFixture(t, "folder")
		if _, err := Prepare(entity, PrepareInput{WorkflowDir: workflow, Question: "Advance?", Artifact: artifact, Summary: "candidate"}); err != nil {
			t.Fatal(err)
		}
		before := treeDigest(t, state)
		if _, err := Withdraw(entity, WithdrawInput{WorkflowDir: workflow, Reason: " \t"}); err == nil {
			t.Fatal("blank reason was accepted")
		}
		if got := treeDigest(t, state); got != before {
			t.Fatal("blank reason changed state tree")
		}
	})

	t.Run("provider output", func(t *testing.T) {
		workflow, state, entity, artifact, _ := prepareFixture(t, "folder")
		prepared, err := Prepare(entity, PrepareInput{WorkflowDir: workflow, Question: "Advance?", Artifact: artifact, Summary: "candidate"})
		if err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(filepath.Join(prepared.Room, "provider"), 0o755); err != nil {
			t.Fatal(err)
		}
		before := treeDigest(t, state)
		if _, err := Withdraw(entity, WithdrawInput{WorkflowDir: workflow, Reason: "stale"}); err == nil ||
			!strings.Contains(err.Error(), "exactly two regular files") {
			t.Fatalf("provider-output withdrawal = %v", err)
		}
		if got := treeDigest(t, state); got != before {
			t.Fatal("provider-output refusal changed state tree")
		}
	})

	t.Run("corrupt retained request", func(t *testing.T) {
		workflow, state, entity, artifact, _ := prepareFixture(t, "folder")
		prepared, err := Prepare(entity, PrepareInput{WorkflowDir: workflow, Question: "Advance?", Artifact: artifact, Summary: "candidate"})
		if err != nil {
			t.Fatal(err)
		}
		request := filepath.Join(prepared.Room, "request.json")
		body := readPreparedFile(t, request)
		if err := os.WriteFile(request, bytes.Replace(body, []byte(`"actor": "person:captain"`), []byte(`"actor": "agent:other"`), 1), 0o644); err != nil {
			t.Fatal(err)
		}
		before := treeDigest(t, state)
		if _, err := Withdraw(entity, WithdrawInput{WorkflowDir: workflow, Reason: "stale"}); err == nil ||
			!strings.Contains(err.Error(), "frozen digest") {
			t.Fatalf("corrupt-authority withdrawal = %v", err)
		}
		if got := treeDigest(t, state); got != before {
			t.Fatal("corrupt-authority refusal changed state tree")
		}
	})

	t.Run("repeat and closed", func(t *testing.T) {
		for _, terminal := range []string{"withdrawn", "closed"} {
			t.Run(terminal, func(t *testing.T) {
				workflow, state, entity, artifact, _ := prepareFixture(t, "folder")
				if _, err := Prepare(entity, PrepareInput{WorkflowDir: workflow, Question: "Advance?", Artifact: artifact, Summary: "candidate"}); err != nil {
					t.Fatal(err)
				}
				if terminal == "withdrawn" {
					if _, err := Withdraw(entity, WithdrawInput{WorkflowDir: workflow, Reason: "stale"}); err != nil {
						t.Fatal(err)
					}
				} else if err := RecordSemantic(entity, RecordInput{
					Actor: "person:captain", Decision: "hold", Reason: "wait", WorkflowDir: workflow,
				}); err != nil {
					t.Fatal(err)
				}
				before := treeDigest(t, state)
				if _, err := Withdraw(entity, WithdrawInput{WorkflowDir: workflow, Reason: "stale"}); err == nil ||
					!strings.Contains(err.Error(), "frozen "+terminal) {
					t.Fatalf("%s withdrawal = %v", terminal, err)
				}
				if got := treeDigest(t, state); got != before {
					t.Fatalf("%s refusal changed state tree", terminal)
				}
			})
		}
	})

	t.Run("lock contention", func(t *testing.T) {
		workflow, state, entity, artifact, _ := prepareFixture(t, "folder")
		if _, err := Prepare(entity, PrepareInput{WorkflowDir: workflow, Question: "Advance?", Artifact: artifact, Summary: "candidate"}); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(entity+".gates.lock", []byte("held"), 0o600); err != nil {
			t.Fatal(err)
		}
		before := treeDigest(t, state)
		if _, err := Withdraw(entity, WithdrawInput{WorkflowDir: workflow, Reason: "stale"}); err == nil ||
			!strings.Contains(err.Error(), "concurrent gate writer") {
			t.Fatalf("lock contention = %v", err)
		}
		if got := treeDigest(t, state); got != before {
			t.Fatal("lock contention changed state tree")
		}
	})

	t.Run("chat-only request-less attempt", func(t *testing.T) {
		workflow, entity := recordStageFixture(t, "validation",
			"briefing:task:validation:attempt-2:revision-1", "      gate: true\n")
		before := treeDigest(t, workflow)
		if _, err := Withdraw(entity, WithdrawInput{WorkflowDir: workflow, Reason: "stale"}); err == nil ||
			!strings.Contains(err.Error(), "not request-backed") {
			t.Fatalf("request-less withdrawal = %v", err)
		}
		if got := treeDigest(t, workflow); got != before {
			t.Fatal("request-less refusal changed workflow tree")
		}
	})

	t.Run("stale current selection", func(t *testing.T) {
		workflow, _, entity, artifact, _ := prepareFixture(t, "folder")
		if _, err := Prepare(entity, PrepareInput{WorkflowDir: workflow, Question: "Advance?", Artifact: artifact, Summary: "candidate"}); err != nil {
			t.Fatal(err)
		}
		doc, oldNode, err := Read(entity)
		if err != nil {
			t.Fatal(err)
		}
		doc.Records = append(doc.Records, GateRecord{
			ID: "gate:task:other", Stage: "other", Attempts: []Attempt{{
				ID: "gate-attempt:task-other-1",
				Briefing: Briefing{
					ID: "briefing:task:other:attempt-1:revision-1", Digest: "sha256:" + strings.Repeat("3", 64),
					RoomRef: "./other",
				},
			}},
		})
		if err := writeDocument(entity, oldNode, doc); err != nil {
			t.Fatal(err)
		}
		if _, err := Withdraw(entity, WithdrawInput{WorkflowDir: workflow, Reason: "stale"}); err != nil {
			t.Fatalf("status-derived withdrawal = %v", err)
		}
	})
}

func readPreparedFile(t *testing.T, path string) []byte {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func TestPrepareUsesSlugIdentityWhenEntityHasNoStoredID(t *testing.T) {
	for _, form := range []string{"folder", "flat"} {
		t.Run(form, func(t *testing.T) {
			workflow, state, entity, artifact, _ := prepareFixture(t, form)
			body, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			body = bytes.Replace(body, []byte("id: task\n"), nil, 1)
			if err := os.WriteFile(entity, body, 0o644); err != nil {
				t.Fatal(err)
			}
			prepareGitRun(t, state, "add", ".")
			prepareGitRun(t, state, "commit", "-q", "-m", "slug identity fixture")

			result, err := Prepare(entity, PrepareInput{
				WorkflowDir: workflow,
				Question:    "Should this gate advance?",
				Artifact:    artifact,
				Summary:     "Exact summary.",
			})
			if err != nil {
				t.Fatal(err)
			}
			if result.Briefing != "briefing:task:validation:attempt-1:revision-1" {
				t.Fatalf("briefing=%q", result.Briefing)
			}
			doc, _, err := Read(entity)
			if err != nil {
				t.Fatal(err)
			}
			record, err := recordForStage(doc, "validation")
			if err != nil || record.ID != "gate:task:validation" {
				t.Fatalf("status-derived gate=%v/%q", err, record.ID)
			}
		})
	}
}

func TestPrepareRejectsDivergentOccupancyAndLeavesEntityUnchanged(t *testing.T) {
	workflow, state, entity, artifact, _ := prepareFixture(t, "flat")
	room := filepath.Join(state, "task", "review", "validation", "briefing-1")
	if err := os.MkdirAll(room, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(room, "foreign.txt"), []byte("occupied"), 0o644); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(entity)
	if _, err := Prepare(entity, PrepareInput{
		WorkflowDir: workflow,
		Question:    "Review?",
		Artifact:    artifact,
		Summary:     "summary",
	}); err == nil || !strings.Contains(err.Error(), "occupied") {
		t.Fatalf("divergent occupancy error=%v", err)
	}
	after, _ := os.ReadFile(entity)
	if !bytes.Equal(before, after) {
		t.Fatal("occupied room changed entity")
	}
	if _, err := os.Stat(entity + ".gates.lock"); !os.IsNotExist(err) {
		t.Fatalf("prepare left lock residue: %v", err)
	}
}

func TestPrepareAttributesRejectedSelectedSourceToItsFlag(t *testing.T) {
	for _, flag := range []string{"--reference", "--artifact"} {
		t.Run(strings.TrimPrefix(flag, "--"), func(t *testing.T) {
			workflow, state, entity, artifact, reference := prepareFixture(t, "flat")
			input := PrepareInput{
				WorkflowDir: workflow,
				Question:    "Should this gate advance?",
				Artifact:    artifact,
				Summary:     "Exact summary.",
				References:  []string{reference},
			}
			uncommitted := filepath.Join(state, "uncommitted-reference.json")
			if flag == "--artifact" {
				uncommitted = filepath.Join(filepath.Dir(filepath.Dir(workflow)), "uncommitted-artifact.md")
				input.Artifact = uncommitted
			} else {
				input.References = append(input.References, uncommitted)
			}
			if err := os.WriteFile(uncommitted, []byte("stub\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			_, err := Prepare(entity, input)
			want := flag + " " + uncommitted + ":"
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Fatalf("error=%v want prefix %q", err, want)
			}
		})
	}
}

func TestPrepareReplaySurvivesRequiredStateCommit(t *testing.T) {
	workflow, state, entity, artifact, reference := prepareFixture(t, "flat")
	input := PrepareInput{
		WorkflowDir: workflow,
		Question:    "Should this gate advance?",
		Artifact:    artifact,
		Summary:     "Exact summary.",
		References:  []string{reference},
	}
	result, err := Prepare(entity, input)
	if err != nil {
		t.Fatal(err)
	}
	prepareGitRun(t, state, "add", "task.md", "task")
	prepareGitRun(t, state, "commit", "-q", "-m", "bind prepared room")
	beforeHead := strings.TrimSpace(prepareGitOutput(t, state, "rev-parse", "HEAD"))
	beforeEntity, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	beforeBriefing, err := os.ReadFile(filepath.Join(result.Room, preparedBriefingLocator))
	if err != nil {
		t.Fatal(err)
	}
	replayed, err := Prepare(entity, input)
	if err != nil {
		t.Fatal(err)
	}
	if replayed != result {
		t.Fatalf("replay=%#v want %#v", replayed, result)
	}
	afterHead := strings.TrimSpace(prepareGitOutput(t, state, "rev-parse", "HEAD"))
	afterEntity, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	afterBriefing, err := os.ReadFile(filepath.Join(result.Room, preparedBriefingLocator))
	if err != nil {
		t.Fatal(err)
	}
	if afterHead != beforeHead || !bytes.Equal(afterEntity, beforeEntity) || !bytes.Equal(afterBriefing, beforeBriefing) {
		t.Fatal("post-commit replay changed state HEAD or durable bytes")
	}
	if status := prepareGitOutput(t, state, "status", "--porcelain"); strings.TrimSpace(status) != "" {
		t.Fatalf("post-commit replay dirtied state repository: %q", status)
	}
}

func TestPrepareSuccessorRejectsCorruptedRetainedAuthorityWithoutChangingTree(t *testing.T) {
	workflow, state, entity, artifact, _ := prepareFixture(t, "flat")
	input := PrepareInput{
		WorkflowDir: workflow,
		Question:    "Should this gate advance?",
		Artifact:    artifact,
		Summary:     "Exact summary.",
	}
	first, err := Prepare(entity, input)
	if err != nil {
		t.Fatal(err)
	}
	if err := RecordSemantic(entity, RecordInput{
		Decision:    "approve",
		Actor:       "person:captain",
		WorkflowDir: workflow,
	}); err != nil {
		t.Fatal(err)
	}
	requestPath := filepath.Join(first.Room, "request.json")
	request, err := os.ReadFile(requestPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(requestPath, bytes.Replace(request, []byte(`"actor": "person:captain"`), []byte(`"actor": "agent:other"`), 1), 0o644); err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, state)

	if _, err := Prepare(entity, input); err == nil || !strings.Contains(err.Error(), "retained request.json") {
		t.Fatalf("successor prepare retained-authority error=%v", err)
	}
	if after := treeDigest(t, state); after != before {
		t.Fatal("rejected successor preparation changed the entity or room tree")
	}
}

func TestPrepareReplayRejectsSymlinkedAuthorityEntriesWithoutChangingBytes(t *testing.T) {
	for _, name := range []string{preparedBriefingLocator, "request.json"} {
		t.Run(name, func(t *testing.T) {
			workflow, state, entity, artifact, _ := prepareFixture(t, "flat")
			input := PrepareInput{
				WorkflowDir: workflow,
				Question:    "Should this gate advance?",
				Artifact:    artifact,
				Summary:     "Exact summary.",
			}
			prepared, err := Prepare(entity, input)
			if err != nil {
				t.Fatal(err)
			}
			roomEntry := filepath.Join(prepared.Room, name)
			body, err := os.ReadFile(roomEntry)
			if err != nil {
				t.Fatal(err)
			}
			external := filepath.Join(t.TempDir(), name)
			if err := os.WriteFile(external, body, 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.Remove(roomEntry); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(external, roomEntry); err != nil {
				t.Fatal(err)
			}
			beforeState := treeDigest(t, state)
			beforeExternal := append([]byte(nil), body...)

			if _, err := Prepare(entity, input); err == nil {
				t.Fatalf("symlinked %s replay error=%v", name, err)
			}
			if after := treeDigest(t, state); after != beforeState {
				t.Fatalf("symlinked %s replay changed the state tree", name)
			}
			afterExternal, err := os.ReadFile(external)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(afterExternal, beforeExternal) {
				t.Fatalf("symlinked %s replay changed external bytes", name)
			}
			if info, err := os.Lstat(roomEntry); err != nil || info.Mode()&os.ModeSymlink == 0 {
				t.Fatalf("symlinked %s replay changed the room entry: info=%v err=%v", name, info, err)
			}
		})
	}
}

func TestPrepareReplayRejectsUnavailablePreparedGitSourceWithoutChangingTree(t *testing.T) {
	workflow, state, entity, artifact, _ := prepareFixture(t, "flat")
	originalEntity, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Prepare(entity, PrepareInput{
		WorkflowDir: workflow,
		Question:    "Should this gate advance?",
		Artifact:    artifact,
		Summary:     "Exact summary.",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(entity, originalEntity, 0o644); err != nil {
		t.Fatal(err)
	}
	blob := strings.TrimSpace(prepareGitOutput(t, filepath.Dir(artifact), "hash-object", artifact))
	object := filepath.Join(filepath.Dir(artifact), ".git", "objects", blob[:2], blob[2:])
	if err := os.Remove(object); err != nil {
		t.Fatalf("remove selected local object: %v", err)
	}
	before := treeDigest(t, state)

	if _, err := Prepare(entity, PrepareInput{
		WorkflowDir: workflow,
		Question:    "Should this gate advance?",
		Artifact:    artifact,
		Summary:     "Exact summary.",
	}); err == nil || !strings.Contains(err.Error(), "selected source") {
		t.Fatalf("prepared-room replay unavailable source error=%v", err)
	}
	if after := treeDigest(t, state); after != before {
		t.Fatal("rejected request-backed record changed the entity or room tree")
	}
}

func TestPrepareReplayAcceptsEntitySelectedAsArtifactOrReference(t *testing.T) {
	for _, entityRole := range []string{"artifact", "reference"} {
		t.Run(entityRole, func(t *testing.T) {
			workflow, state, entity, artifact, _ := prepareFixture(t, "flat")
			input := PrepareInput{
				WorkflowDir: workflow,
				Question:    "Should this gate advance?",
				Artifact:    artifact,
				Summary:     "Exact summary.",
			}
			if entityRole == "artifact" {
				input.Artifact = entity
			} else {
				input.References = []string{entity}
			}
			result, err := Prepare(entity, input)
			if err != nil {
				t.Fatal(err)
			}
			beforeEntity, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			replayed, err := Prepare(entity, input)
			if err != nil {
				t.Fatalf("pre-commit replay: %v", err)
			}
			if replayed != result {
				t.Fatalf("pre-commit replay=%#v want %#v", replayed, result)
			}
			afterEntity, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(afterEntity, beforeEntity) {
				t.Fatal("pre-commit replay changed entity bytes")
			}

			prepareGitRun(t, state, "add", "task.md", "task")
			prepareGitRun(t, state, "commit", "-q", "-m", "bind entity-selected room")
			beforeHead := strings.TrimSpace(prepareGitOutput(t, state, "rev-parse", "HEAD"))
			replayed, err = Prepare(entity, input)
			if err != nil {
				t.Fatalf("post-commit replay: %v", err)
			}
			if replayed != result {
				t.Fatalf("post-commit replay=%#v want %#v", replayed, result)
			}
			if afterHead := strings.TrimSpace(prepareGitOutput(t, state, "rev-parse", "HEAD")); afterHead != beforeHead {
				t.Fatal("post-commit replay changed state HEAD")
			}
			if status := prepareGitOutput(t, state, "status", "--porcelain"); strings.TrimSpace(status) != "" {
				t.Fatalf("post-commit replay dirtied state repository: %q", status)
			}

			changed := bytes.Replace(beforeEntity, []byte("title: Task"), []byte("title: Changed"), 1)
			if err := os.WriteFile(entity, changed, 0o644); err != nil {
				t.Fatal(err)
			}
			if _, err := Prepare(entity, input); err == nil {
				t.Fatal("replay accepted an entity change outside binary-owned gate state")
			}
			afterRejected, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(afterRejected, changed) {
				t.Fatal("rejected replay changed entity bytes")
			}
		})
	}
}

func TestPrepareRejectsSymlinkedFlatCompanionWithoutChangingBytes(t *testing.T) {
	workflow, state, entity, artifact, _ := prepareFixture(t, "flat")
	external := filepath.Join(filepath.Dir(state), "external")
	if err := os.MkdirAll(filepath.Join(external, "review"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(external, "sentinel"), []byte("outside\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	companion := filepath.Join(state, "task")
	if err := os.Symlink(external, companion); err != nil {
		t.Fatal(err)
	}
	beforeEntity, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	beforeExternal := treeDigest(t, external)
	beforeStatus := prepareGitOutput(t, state, "status", "--porcelain")

	if _, err := Prepare(entity, PrepareInput{
		WorkflowDir: workflow,
		Question:    "Review?",
		Artifact:    artifact,
		Summary:     "summary",
	}); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("symlinked flat companion error=%v", err)
	}
	afterEntity, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(afterEntity, beforeEntity) {
		t.Fatal("symlink rejection changed entity bytes")
	}
	if got := treeDigest(t, external); got != beforeExternal {
		t.Fatal("symlink rejection changed external bytes")
	}
	if got := prepareGitOutput(t, state, "status", "--porcelain"); got != beforeStatus {
		t.Fatalf("symlink rejection changed state tree: before=%q after=%q", beforeStatus, got)
	}
	if target, err := os.Readlink(companion); err != nil || target != external {
		t.Fatalf("flat companion symlink changed: target=%q err=%v", target, err)
	}
	if _, err := os.Stat(entity + ".gates.lock"); !os.IsNotExist(err) {
		t.Fatalf("prepare left lock residue: %v", err)
	}
}

func TestPrepareRejectsUnsafeOrUndefinedStatusBeforePublishing(t *testing.T) {
	for _, status := range []string{"../../../outside", "rogue"} {
		t.Run(status, func(t *testing.T) {
			workflow, state, entity, artifact, _ := prepareFixture(t, "flat")
			before, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			mutated := bytes.Replace(before, []byte("status: validation"), []byte("status: "+status), 1)
			if err := os.WriteFile(entity, mutated, 0o644); err != nil {
				t.Fatal(err)
			}
			if _, err := Prepare(entity, PrepareInput{
				WorkflowDir: workflow,
				Question:    "Review?",
				Artifact:    artifact,
				Summary:     "summary",
			}); err == nil || !strings.Contains(err.Error(), "workflow stage") {
				t.Fatalf("status %q error=%v", status, err)
			}
			if _, err := os.Stat(filepath.Join(state, "task")); !os.IsNotExist(err) {
				t.Fatalf("status %q published a flat companion: %v", status, err)
			}
			if _, err := os.Stat(filepath.Join(filepath.Dir(state), "outside")); !os.IsNotExist(err) {
				t.Fatalf("status %q escaped the state root: %v", status, err)
			}
			after, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(after, mutated) {
				t.Fatalf("status %q changed entity bytes", status)
			}
		})
	}
}

func TestPrepareRollsBackPublishedRoomAfterBindingWriteFailure(t *testing.T) {
	original := prepareWriteBinding
	prepareWriteBinding = func(string, *yaml.Node, *Document) error {
		return os.ErrPermission
	}
	t.Cleanup(func() { prepareWriteBinding = original })

	for _, preexistingHome := range []bool{false, true} {
		t.Run(map[bool]string{false: "new-home", true: "preexisting-home"}[preexistingHome], func(t *testing.T) {
			workflow, state, entity, artifact, _ := prepareFixture(t, "flat")
			home := filepath.Join(state, "task")
			if preexistingHome {
				if err := os.Mkdir(home, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			before, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := Prepare(entity, PrepareInput{
				WorkflowDir: workflow,
				Question:    "Review?",
				Artifact:    artifact,
				Summary:     "summary",
			}); !os.IsPermission(err) {
				t.Fatalf("post-publication error=%v, want permission failure", err)
			}
			if _, err := os.Stat(filepath.Join(home, "review")); !os.IsNotExist(err) {
				t.Fatalf("rollback retained review parents: %v", err)
			}
			_, homeErr := os.Stat(home)
			if preexistingHome && homeErr != nil {
				t.Fatalf("rollback removed pre-existing home: %v", homeErr)
			}
			if !preexistingHome && !os.IsNotExist(homeErr) {
				t.Fatalf("rollback retained newly created home: %v", homeErr)
			}
			after, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(before, after) {
				t.Fatal("post-publication failure changed entity bytes")
			}
		})
	}
}

func TestPreparedAuthorityIsRecomputedDuringReadOnlyValidation(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(t *testing.T, workflow, entity, artifact, room string)
	}{
		{
			name: "request drift",
			mutate: func(t *testing.T, _, _, _, room string) {
				path := filepath.Join(room, "request.json")
				body, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, bytes.Replace(body, []byte(`"actor": "person:captain"`), []byte(`"actor": "agent:other"`), 1), 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "Briefing summary drift",
			mutate: func(t *testing.T, _, _, _, room string) {
				path := filepath.Join(room, preparedBriefingLocator)
				body, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, bytes.Replace(body, []byte(`"summary": "summary"`), []byte(`"summary": "changed"`), 1), 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "missing local source object",
			mutate: func(t *testing.T, _, _, artifact, _ string) {
				blob := strings.TrimSpace(prepareGitOutput(t, filepath.Dir(artifact), "hash-object", artifact))
				object := filepath.Join(filepath.Dir(artifact), ".git", "objects", blob[:2], blob[2:])
				if err := os.Remove(object); err != nil {
					t.Fatalf("remove selected local object: %v", err)
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			workflow, _, entity, artifact, _ := prepareFixture(t, "flat")
			result, err := Prepare(entity, PrepareInput{
				WorkflowDir: workflow,
				Question:    "Review?",
				Artifact:    artifact,
				Summary:     "summary",
			})
			if err != nil {
				t.Fatal(err)
			}
			tc.mutate(t, workflow, entity, artifact, result.Room)
			if _, err := EligibilityFileAt(entity, workflow); err == nil {
				t.Fatal("read-only validation accepted drifted prepared authority")
			}
		})
	}
}

func TestChatDecisionValidatesPreparedAuthorityBeforeMutation(t *testing.T) {
	workflow, _, entity, artifact, _ := prepareFixture(t, "flat")
	result, err := Prepare(entity, PrepareInput{
		WorkflowDir: workflow,
		Question:    "Review?",
		Artifact:    artifact,
		Summary:     "summary",
	})
	if err != nil {
		t.Fatal(err)
	}
	requestPath := filepath.Join(result.Room, "request.json")
	request, err := os.ReadFile(requestPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(requestPath, bytes.Replace(request, []byte(`"actor": "person:captain"`), []byte(`"actor": "agent:other"`), 1), 0o644); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if err := RecordSemantic(entity, RecordInput{
		Decision:    "approve",
		Actor:       "person:captain",
		WorkflowDir: workflow,
	}); err == nil || !strings.Contains(err.Error(), "retained request.json") {
		t.Fatalf("drifted authority chat decision error=%v", err)
	}
	after, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(after, before) {
		t.Fatal("drifted retained authority chat decision changed entity bytes")
	}
}

// TestPrepareResolvesStateRelativeArtifactAgainstEntityRoot is the value AC-1:
// under split-root, a state-relative --artifact path (as the FO would pass it)
// resolves against the state-checkout entity root, not the workflow root, and
// reaches state=open. Reverting resolution to cwd (the workflow root) fails
// because the file does not exist there.
func TestPrepareResolvesStateRelativeArtifactAgainstEntityRoot(t *testing.T) {
	workflow, state, entity, _, _ := prepareFixture(t, "flat")
	selected := filepath.Join(state, "selected", "gate-review.md")
	if err := os.MkdirAll(filepath.Dir(selected), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(selected, []byte("# Selected review\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	prepareGitRun(t, state, "add", "selected")
	prepareGitRun(t, state, "commit", "-q", "-m", "committed selected artifact")

	rel := filepath.ToSlash(filepath.Join("selected", "gate-review.md"))
	result, err := Prepare(entity, PrepareInput{
		WorkflowDir: workflow,
		Question:    "Advance?",
		Artifact:    rel,
		Summary:     "committed selected review",
	})
	if err != nil {
		t.Fatalf("state-relative artifact error=%v", err)
	}
	if result.State != "open" {
		t.Fatalf("state-relative artifact state=%q want open", result.State)
	}
}

// TestPrepareResolvesStateRelativeReferenceAgainstEntityRoot verifies AC-2(d):
// a relative --reference resolves against the same entity root as --artifact.
func TestPrepareResolvesStateRelativeReferenceAgainstEntityRoot(t *testing.T) {
	workflow, state, entity, artifact, _ := prepareFixture(t, "flat")
	selectedRef := filepath.Join(state, "selected", "evidence.json")
	if err := os.MkdirAll(filepath.Dir(selectedRef), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(selectedRef, []byte("{\"ok\":true}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	prepareGitRun(t, state, "add", "selected")
	prepareGitRun(t, state, "commit", "-q", "-m", "committed selected reference")

	relRef := filepath.ToSlash(filepath.Join("selected", "evidence.json"))
	result, err := Prepare(entity, PrepareInput{
		WorkflowDir: workflow,
		Question:    "Advance?",
		Artifact:    artifact,
		Summary:     "with state-relative reference",
		References:  []string{relRef},
	})
	if err != nil {
		t.Fatalf("state-relative reference error=%v", err)
	}
	if result.State != "open" {
		t.Fatalf("state-relative reference state=%q want open", result.State)
	}
}

// TestPrepareAbsoluteArtifactPassesThroughUnchanged verifies AC-2(c): an absolute
// --artifact path passes through unchanged in both split-root and single-root.
func TestPrepareAbsoluteArtifactPassesThroughUnchanged(t *testing.T) {
	t.Run("split-root", func(t *testing.T) {
		workflow, _, entity, artifact, _ := prepareFixture(t, "flat")
		result, err := Prepare(entity, PrepareInput{
			WorkflowDir: workflow,
			Question:    "Advance?",
			Artifact:    artifact,
			Summary:     "absolute path",
		})
		if err != nil {
			t.Fatalf("absolute artifact error=%v", err)
		}
		if result.State != "open" {
			t.Fatalf("absolute artifact state=%q want open", result.State)
		}
	})
	t.Run("single-root", func(t *testing.T) {
		workflow, entity, artifact := prepareSingleRootFixture(t)
		result, err := Prepare(entity, PrepareInput{
			WorkflowDir: workflow,
			Question:    "Advance?",
			Artifact:    artifact,
			Summary:     "absolute path",
		})
		if err != nil {
			t.Fatalf("absolute artifact error=%v", err)
		}
		if result.State != "open" {
			t.Fatalf("absolute artifact state=%q want open", result.State)
		}
	})
}

// TestPrepareSingleRootResolvesRelativeArtifactAgainstWorkflowDir verifies
// AC-2(b): in single-root (no state: field), a relative --artifact resolves
// against the workflow dir unchanged (entity root = workflow dir).
func TestPrepareSingleRootResolvesRelativeArtifactAgainstWorkflowDir(t *testing.T) {
	workflow, entity, _ := prepareSingleRootFixture(t)
	result, err := Prepare(entity, PrepareInput{
		WorkflowDir: workflow,
		Question:    "Advance?",
		Artifact:    "gate-review.md",
		Summary:     "single-root relative",
	})
	if err != nil {
		t.Fatalf("single-root relative artifact error=%v", err)
	}
	if result.State != "open" {
		t.Fatalf("single-root relative artifact state=%q want open", result.State)
	}
}

// TestPrepareWrongRootRelativeArtifactFails is the falsifying change: a relative
// path that resolves to an existing committed file under the workflow root but
// NOT under the state checkout fails. Under the old cwd-based resolution
// (cwd = workflow root in a live run) this would have succeeded against the
// wrong root.
func TestPrepareWrongRootRelativeArtifactFails(t *testing.T) {
	workflow, state, entity, _, _ := prepareFixture(t, "flat")
	wrongRoot := filepath.Join(workflow, "gate-review.md")
	if err := os.WriteFile(wrongRoot, []byte("# Wrong root\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	prepareGitRun(t, filepath.Dir(filepath.Dir(workflow)), "add", filepath.ToSlash(filepath.Join("docs", "dev", "gate-review.md")))
	prepareGitRun(t, filepath.Dir(filepath.Dir(workflow)), "commit", "-q", "-m", "wrong-root artifact")

	resolved := filepath.Join(state, "gate-review.md")
	if _, err := os.Stat(resolved); !os.IsNotExist(err) {
		t.Fatalf("precondition: %s must not exist under state checkout", resolved)
	}
	_, err := Prepare(entity, PrepareInput{
		WorkflowDir: workflow,
		Question:    "Advance?",
		Artifact:    "gate-review.md",
		Summary:     "wrong root",
	})
	if err == nil {
		t.Fatal("wrong-root relative artifact succeeded; resolution reverted to the workflow root")
	}
	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Fatalf("wrong-root error=%v want no such file or directory", err)
	}
}

// TestPrepareWorkflowRootReferenceResolves verifies that a --reference pointing
// to a file at the workflow root (not under the state checkout) resolves
// correctly in a split-root workflow. The FO legitimately references a
// workflow-root file (e.g. recorder-contract.md) alongside a state-rooted
// artifact; the reference must fall back to the workflow directory when the
// entity-root join does not exist. The artifact must stay strict (state-root
// only) — TestPrepareWrongRootRelativeArtifactFails guards that separately.
func TestPrepareWorkflowRootReferenceResolves(t *testing.T) {
	workflow, state, entity, artifact, _ := prepareFixture(t, "flat")
	// Place a workflow-root reference file (NOT under the state checkout).
	workflowRef := filepath.Join(workflow, "recorder-contract.md")
	if err := os.WriteFile(workflowRef, []byte("# Recorder contract\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	prepareGitRun(t, filepath.Dir(filepath.Dir(workflow)), "add", filepath.ToSlash(filepath.Join("docs", "dev", "recorder-contract.md")))
	prepareGitRun(t, filepath.Dir(filepath.Dir(workflow)), "commit", "-q", "-m", "workflow-root reference")

	// Precondition: the reference must NOT exist under the state checkout.
	stateRef := filepath.Join(state, "recorder-contract.md")
	if _, err := os.Stat(stateRef); !os.IsNotExist(err) {
		t.Fatalf("precondition: %s must not exist under state checkout", stateRef)
	}

	result, err := Prepare(entity, PrepareInput{
		WorkflowDir: workflow,
		Question:    "Advance?",
		Artifact:    artifact,
		Summary:     "Split-root with workflow-root reference.",
		References:  []string{"recorder-contract.md"},
	})
	if err != nil {
		t.Fatalf("prepare with workflow-root reference failed: %v", err)
	}
	if result.State != "open" {
		t.Fatalf("prepare state = %q, want open", result.State)
	}
}

// stage where the entity, artifact, and workflow dir all live in the same Git
// repo — the entity root equals the workflow dir.
func prepareSingleRootFixture(t *testing.T) (workflow, entity, artifact string) {
	t.Helper()
	workflow = t.TempDir()
	testgit.InitRepo(t, workflow, "-q")
	if err := os.WriteFile(filepath.Join(workflow, "README.md"), []byte("---\nid-style: slug\nstages:\n  states:\n    - name: validation\n      initial: true\n      gate: true\n    - name: done\n      terminal: true\n---\n# Workflow\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	artifact = filepath.Join(workflow, "gate-review.md")
	if err := os.WriteFile(artifact, []byte("# Review\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	entity = filepath.Join(workflow, "task.md")
	if err := os.WriteFile(entity, []byte("---\nid: task\nstatus: validation\ntitle: Task\n---\n# Task\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	prepareGitRun(t, workflow, "add", ".")
	prepareGitRun(t, workflow, "commit", "-q", "-m", "single-root fixture")
	return workflow, entity, artifact
}

func prepareFixture(t *testing.T, form string) (workflow, state, entity, artifact, reference string) {
	t.Helper()
	root := t.TempDir()
	workflow = filepath.Join(root, "main", "docs", "dev")
	if err := os.MkdirAll(workflow, 0o755); err != nil {
		t.Fatal(err)
	}
	mainRoot := filepath.Dir(filepath.Dir(workflow))
	testgit.InitRepo(t, mainRoot, "-q")
	if err := os.WriteFile(filepath.Join(workflow, "README.md"), []byte("---\nid-style: slug\nstate: .state\nstages:\n  states:\n    - name: validation\n      initial: true\n      gate: true\n    - name: implementation\n    - name: done\n      terminal: true\n    - name: contradictory\n      gate: true\n      terminal: true\n---\n# Workflow\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	artifact = filepath.Join(mainRoot, "gate-review.md")
	if err := os.WriteFile(artifact, []byte("# Review\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	prepareGitRun(t, mainRoot, "add", ".")
	prepareGitRun(t, mainRoot, "commit", "-q", "-m", "main fixture")

	state = filepath.Join(workflow, ".state")
	if err := os.MkdirAll(state, 0o755); err != nil {
		t.Fatal(err)
	}
	testgit.InitRepo(t, state, "-q")
	switch form {
	case "folder":
		entity = filepath.Join(state, "task", "index.md")
	case "flat":
		entity = filepath.Join(state, "task.md")
	default:
		t.Fatalf("unknown form %s", form)
	}
	if err := os.MkdirAll(filepath.Dir(entity), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(entity, []byte("---\nid: task\nstatus: validation\ntitle: Task\n---\n# Task\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	reference = filepath.Join(state, "reference.json")
	if err := os.WriteFile(reference, []byte("{\"evidence\":true}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	prepareGitRun(t, state, "add", ".")
	prepareGitRun(t, state, "commit", "-q", "-m", "state fixture")
	return workflow, state, entity, artifact, reference
}

func prepareGitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func prepareGitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return string(out)
}

// declareFolderForm adds the workflow declaration the prepare guard reads.
// Commissioning and the filing-time default own writing this key; preparation
// only reads it, so the fixture writes it by hand.
func declareFolderForm(t *testing.T, workflow string) {
	t.Helper()
	readme := filepath.Join(workflow, "README.md")
	body, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	declared := bytes.Replace(body, []byte("id-style: slug\n"), []byte("id-style: slug\nentity-form: folder\n"), 1)
	if bytes.Equal(declared, body) {
		t.Fatal("fixture README lost the id-style anchor this helper writes beside")
	}
	if err := os.WriteFile(readme, declared, 0o644); err != nil {
		t.Fatal(err)
	}
	mainRoot := filepath.Dir(filepath.Dir(workflow))
	prepareGitRun(t, mainRoot, "add", "-A")
	prepareGitRun(t, mainRoot, "commit", "-q", "-m", "declare folder form")
}

// TestPrepareRefusesFlatEntityOnlyWhereTheWorkflowDeclaresFolderForm covers the
// defect and its bound. A room beside a flat entity writes a
// ./<slug>/review/... ref that breaks if the entity ever becomes
// <slug>/index.md, so a workflow that has declared folder form refuses to mint
// one. A workflow that declares nothing has promised no such thing and keeps
// preparing on the flat entity that filing still mints by default — the shape
// every live workflow is in. An entity already holding rooms is grandfathered
// under the declaration, because its refs are right for the form it is in.
func TestPrepareRefusesFlatEntityOnlyWhereTheWorkflowDeclaresFolderForm(t *testing.T) {
	for _, tc := range []struct {
		name     string
		declared bool
		rooms    bool
	}{
		{name: "undeclared-workflow-prepares-flat", declared: false, rooms: false},
		{name: "declared-folder-form-refuses-flat", declared: true, rooms: false},
		{name: "declared-folder-form-grandfathers-rooms", declared: true, rooms: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			workflow, state, entity, artifact, _ := prepareFixture(t, "flat")
			if tc.declared {
				declareFolderForm(t, workflow)
			}
			if tc.rooms {
				if err := os.MkdirAll(filepath.Join(state, "task", "review"), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			before := prepareTreeSnapshot(t, state)

			result, err := Prepare(entity, PrepareInput{
				WorkflowDir: workflow, Question: "Review?", Artifact: artifact, Summary: "summary",
			})

			if tc.declared && !tc.rooms {
				if err == nil || !strings.Contains(err.Error(), "requires folder-form entity task/index.md") ||
					!strings.Contains(err.Error(), "declares `entity-form: folder`") {
					t.Fatalf("declared-form refusal error=%v", err)
				}
				if got := prepareTreeSnapshot(t, state); got != before {
					t.Fatal("refusal changed the state tree")
				}
				if _, statErr := os.Stat(filepath.Join(state, "task")); !os.IsNotExist(statErr) {
					t.Fatalf("refusal minted a companion: %v", statErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("prepare refused a shape it may not refuse: %v", err)
			}
			wantRoom := filepath.Join(state, "task", "review", "validation", "briefing-1")
			if result.Room != wantRoom {
				t.Fatalf("room=%q want %q", result.Room, wantRoom)
			}
			doc, _, err := Read(entity)
			if err != nil {
				t.Fatal(err)
			}
			// Flat form binds a slug-prefixed ref. It is correct while the entity
			// stays flat and is exactly what a later conversion must rewrite; the
			// validator warning carries that instruction.
			if got := doc.Records[0].Attempts[0].Briefing.RoomRef; got != "./task/review/validation/briefing-1" {
				t.Fatalf("flat room-ref=%q", got)
			}
		})
	}
}
