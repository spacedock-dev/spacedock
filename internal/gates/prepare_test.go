package gates

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gitsource"
)

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

func TestPreparedCandidateFailureLeavesNoRoomParents(t *testing.T) {
	root := t.TempDir()
	room := filepath.Join(root, "task", "review", "validation", "briefing-1")
	err := validatePreparedCandidate(
		filepath.Join(root, "task.md"),
		gitsource.Roots{},
		room,
		Briefing{},
		"gate:task:validation",
		"gate-attempt:task-validation-1",
		[]byte(`{"type":"Briefing","type":"duplicate"}`),
		[]byte(`{}`),
	)
	if err == nil {
		t.Fatal("invalid candidate unexpectedly passed")
	}
	if _, err := os.Stat(filepath.Join(root, "task")); !os.IsNotExist(err) {
		t.Fatalf("candidate validation changed prepared tree: %v", err)
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
			if _, err := SummaryFileAt(entity, workflow); err == nil {
				t.Fatal("read-only validation accepted drifted prepared authority")
			}
		})
	}
}

func prepareFixture(t *testing.T, form string) (workflow, state, entity, artifact, reference string) {
	t.Helper()
	root := t.TempDir()
	workflow = filepath.Join(root, "main", "docs", "dev")
	if err := os.MkdirAll(workflow, 0o755); err != nil {
		t.Fatal(err)
	}
	mainRoot := filepath.Dir(filepath.Dir(workflow))
	prepareGitRun(t, mainRoot, "init", "-q")
	prepareGitIdentity(t, mainRoot)
	if err := os.WriteFile(filepath.Join(workflow, "README.md"), []byte("---\nid-style: slug\nstate: ../../../state\nstages:\n  states:\n    - name: validation\n      initial: true\n      gate: true\n    - name: done\n      terminal: true\n---\n# Workflow\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	artifact = filepath.Join(mainRoot, "gate-review.md")
	if err := os.WriteFile(artifact, []byte("# Review\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	prepareGitRun(t, mainRoot, "add", ".")
	prepareGitRun(t, mainRoot, "commit", "-q", "-m", "main fixture")

	state = filepath.Join(root, "state")
	if err := os.MkdirAll(state, 0o755); err != nil {
		t.Fatal(err)
	}
	prepareGitRun(t, state, "init", "-q")
	prepareGitIdentity(t, state)
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

func prepareGitIdentity(t *testing.T, dir string) {
	t.Helper()
	prepareGitRun(t, dir, "config", "user.name", "Spacedock Test")
	prepareGitRun(t, dir, "config", "user.email", "spacedock@example.invalid")
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
