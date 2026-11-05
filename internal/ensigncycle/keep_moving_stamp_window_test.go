package ensigncycle

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestKeepMovingStrayCommitKeepsDispatchProof replays the dispatch window from
// Sonnet run 31894934349 attempt 2 with the real binary: the non-gated advance,
// a stray `git add -A` commit by the FO, then `dispatch build --stamp`. When the
// advance leaves `started` to --stamp, the exact path-scoped
// `dispatch: ready-one entering implementation` commit still exists. When the
// advance writes `started` itself, the stray commit takes the only change, so
// --stamp has nothing to commit and the proof is gone.
func TestKeepMovingStrayCommitKeepsDispatchProof(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	for _, tc := range []struct {
		name      string
		advance   []string
		wantProof bool
	}{
		{"advance-without-started", []string{"status=" + kmNextStage}, true},
		{"advance-with-started", []string{"status=" + kmNextStage, "started"}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, filepath.Join(root, "README.md"), keepMovingReadme())
			entity := filepath.Join(root, kmReadyOne+".md")
			writeFile(t, entity, keepMovingReadyEntity(kmReadyOne, "Ready One"))
			gitInit(t, root)

			mustRecordedGate(t, binary, root, append([]string{"status", "--workflow-dir", root, "--set", kmReadyOne}, tc.advance...)...)
			// The stray hand commit also sweeps a sibling's pending state.
			writeFile(t, filepath.Join(root, kmQuestioned+".md"), keepMovingQuestionedEntity())
			git(t, root, "add", "-A")
			git(t, root, "commit", "-q", "-m", "dispatch: batch of entities entering their stages")

			checklist := filepath.Join(t.TempDir(), "checklist")
			writeFile(t, checklist, "Write the implementation stage report.\n")
			cmd := exec.Command(binary, "dispatch", "build", "--stamp", "--host", "codex", "--workflow-dir", root,
				"--entity-path", entity, "--stage", kmNextStage, "--checklist-file", checklist)
			cmd.Env = append(os.Environ(), "HOME="+t.TempDir())
			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("dispatch build --stamp: %v\n%s", err, out)
			}
			var envelope struct {
				DispatchFile string `json:"dispatch_file_path"`
			}
			if json.Unmarshal(out, &envelope) == nil && envelope.DispatchFile != "" {
				t.Cleanup(func() { os.Remove(envelope.DispatchFile) })
			}

			history := durableEntityHistory(t, root, kmReadyOne, kmReadyOne+".md")
			proof := false
			for i := range history {
				proof = proof || durableDispatchProof(history, i, kmReadyOne, kmNextStage)
			}
			if proof != tc.wantProof {
				t.Fatalf("path-scoped dispatch proof = %v, want %v; history %+v", proof, tc.wantProof, history)
			}
			if durableField(readFile(t, entity), "started") == "" {
				t.Fatal("started is empty after --stamp")
			}
		})
	}
}
