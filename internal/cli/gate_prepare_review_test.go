package cli

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/statesync"
	"github.com/spacedock-dev/spacedock/internal/status"
	"github.com/spacedock-dev/spacedock/internal/testgit"
)

type prepareReviewFixture struct {
	definition string
	state      string
	entity     string
	slug       string
}

func newPrepareReviewFixture(t *testing.T) prepareReviewFixture {
	t.Helper()
	definition := t.TempDir()
	state := filepath.Join(definition, ".spacedock-state")
	slug := "gate-task"
	entity := filepath.Join(state, slug+".md")
	writeFixtureFile(t, filepath.Join(definition, "README.md"), `---
commissioned-by: spacedock@1
id-style: slug
state: .spacedock-state
stages:
  states:
    - name: validation
      initial: true
      gate: true
    - name: done
      terminal: true
---
# Workflow

### validation

Present the retained implementation evidence to the Captain.
`)
	writeFixtureFile(t, filepath.Join(definition, "policy.md"), "# Policy\n\nRetain exact authority.\n")
	writeFixtureFile(t, entity, `---
id: gate-task
status: validation
---
# Gate task

## Acceptance criteria

**AC-1** — The implementation is retained.

## Stage Report: validation

- DONE: Validate the implementation. AC-1
  The retained review names the implementation commit.

### Summary

Ready for Captain review.
`)
	writeFixtureFile(t, filepath.Join(state, slug, "selected", "review.md"), "# Review\n\nThe implementation is ready.\n")
	writeFixtureFile(t, filepath.Join(state, "dirty-sibling.md"), "untracked sibling\n")
	writeFixtureFile(t, filepath.Join(definition, ".gitignore"), ".spacedock-state/\n")
	testgit.InitRepo(t, definition)
	git(t, definition, "add", "README.md", "policy.md", ".gitignore")
	git(t, definition, "commit", "-m", "definition")
	testgit.InitRepo(t, state)
	stateBranch, err := status.StateBranch(definition)
	if err != nil {
		t.Fatal(err)
	}
	git(t, state, "branch", "-M", stateBranch)
	git(t, state, "add", slug+".md", filepath.Join(slug, "selected", "review.md"))
	git(t, state, "commit", "-m", "state")
	return prepareReviewFixture{definition: definition, state: state, entity: entity, slug: slug}
}

func writeFixtureFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runPrepareReview(t *testing.T, f prepareReviewFixture, args ...string) (string, string, int) {
	t.Helper()
	var out, errOut strings.Builder
	argv := append([]string{"gate", "prepare-review", f.slug, "--workflow-dir", f.definition}, args...)
	code := run(context.Background(), argv, os.Environ(), f.definition, nil, &out, &errOut, &status.NativeRunner{}, nil)
	return out.String(), errOut.String(), code
}

type prepareReviewEnvelope struct {
	Command     string              `json:"command"`
	Mode        string              `json:"mode"`
	Phase       string              `json:"phase"`
	LaunchCWD   string              `json:"launch_cwd"`
	PublishArgv []string            `json:"publish_argv"`
	Entity      map[string]string   `json:"entity"`
	Stage       map[string]string   `json:"stage"`
	Candidates  []map[string]string `json:"candidates"`
	Preparation struct {
		Room, Briefing, Digest, State string
	} `json:"preparation"`
	Sync         map[string]any  `json:"sync"`
	Checklist    json.RawMessage `json:"checklist"`
	Acceptance   json.RawMessage `json:"acceptance"`
	Presentation map[string]any  `json:"presentation"`
}

func decodePrepareReview(t *testing.T, raw string) prepareReviewEnvelope {
	t.Helper()
	var got prepareReviewEnvelope
	if err := json.Unmarshal([]byte(raw), &got); err != nil {
		t.Fatalf("parse JSON: %v\n%s", err, raw)
	}
	return got
}

func TestGatePrepareReviewInspectIsReadOnlyAndAuthoritative(t *testing.T) {
	f := newPrepareReviewFixture(t)
	mainBefore := git(t, f.definition, "status", "--porcelain=v1")
	stateBefore := git(t, f.state, "status", "--porcelain=v1")
	headBefore := strings.TrimSpace(git(t, f.state, "rev-parse", "HEAD"))

	out, errOut, code := runPrepareReview(t, f, "--json")
	if code != 0 {
		t.Fatalf("inspect exit=%d stderr=%q", code, errOut)
	}
	got := decodePrepareReview(t, out)
	if got.Command != "gate prepare-review" || got.Mode != "inspect" || got.LaunchCWD != f.definition {
		t.Fatalf("wrong inspect envelope: %+v", got)
	}
	if got.Stage["name"] != "validation" || !strings.Contains(got.Stage["bytes"], "retained implementation evidence") {
		t.Fatalf("stage prose missing: %#v", got.Stage)
	}
	if len(got.Candidates) != 2 || len(got.PublishArgv) == 0 || got.PublishArgv[0] != "spacedock" {
		t.Fatalf("candidate/argv contract missing: candidates=%v argv=%v", got.Candidates, got.PublishArgv)
	}
	if git(t, f.definition, "status", "--porcelain=v1") != mainBefore || git(t, f.state, "status", "--porcelain=v1") != stateBefore || strings.TrimSpace(git(t, f.state, "rev-parse", "HEAD")) != headBefore {
		t.Fatal("inspect mutated a repository")
	}
}

func TestGatePrepareReviewPublishCommitsBindingAndRoomOnce(t *testing.T) {
	f := newPrepareReviewFixture(t)
	out, errOut, code := runPrepareReview(t, f, "--publish", "--question", "Approve this implementation?", "--artifact", ".spacedock-state/gate-task/selected/review.md", "--reference", "policy.md", "--summary", "The implementation satisfies AC-1.", "--recommendation", "Approve the retained implementation.", "--json")
	if code != 0 {
		t.Fatalf("publish exit=%d stderr=%q", code, errOut)
	}
	got := decodePrepareReview(t, out)
	if got.Mode != "publish" || got.Phase != "complete" || got.Preparation.State != "open" || got.Preparation.Briefing == "" || got.Preparation.Digest == "" {
		t.Fatalf("incomplete publication: %+v", got)
	}
	if len(got.Checklist) == 0 || len(got.Acceptance) == 0 || got.Presentation["recommendation"] != "Approve the retained implementation." {
		t.Fatalf("projection/presentation missing: %+v", got)
	}
	if entity := string(readPrepareReviewFile(t, f.entity)); strings.Contains(entity, "Approve the retained") || strings.Contains(entity, "resolution:") || strings.Contains(entity, "application:") {
		t.Fatalf("composite crossed recommendation/decision authority: %s", entity)
	}
	names := strings.Fields(git(t, f.state, "show", "--name-only", "--pretty=format:", "HEAD"))
	if len(names) != 3 || names[0] != "gate-task.md" || !strings.Contains(strings.Join(names, "\n"), "gate-briefing.json") || !strings.Contains(strings.Join(names, "\n"), "request.json") {
		t.Fatalf("binding and room were not one exact commit: %q", names)
	}
	if porcelain := git(t, f.state, "status", "--porcelain=v1"); !strings.Contains(porcelain, "dirty-sibling.md") || strings.Contains(porcelain, "gate-task.md") {
		t.Fatalf("dirty sibling or entity unit mishandled:\n%s", porcelain)
	}
	head := strings.TrimSpace(git(t, f.state, "rev-parse", "HEAD"))
	out, errOut, code = runPrepareReview(t, f, "--publish", "--question", "Approve this implementation?", "--artifact", ".spacedock-state/gate-task/selected/review.md", "--reference", "policy.md", "--summary", "The implementation satisfies AC-1.", "--recommendation", "Approve the retained implementation.", "--json")
	if code != 0 {
		t.Fatalf("replay exit=%d stderr=%q", code, errOut)
	}
	if strings.TrimSpace(git(t, f.state, "rev-parse", "HEAD")) != head {
		t.Fatal("exact replay created a duplicate attempt commit")
	}
	if got = decodePrepareReview(t, out); got.Preparation.Briefing == "" || got.Phase != "complete" {
		t.Fatalf("replay did not re-emit projection: %+v", got)
	}
}

func TestGatePrepareReviewPeerConflictHaltsAfterOneDurableAttempt(t *testing.T) {
	f := newPrepareReviewFixture(t)
	branch, _ := status.StateBranch(f.definition)
	bare := filepath.Join(t.TempDir(), "state.git")
	git(t, t.TempDir(), "init", "--bare", bare)
	git(t, f.state, "remote", "add", "origin", bare)
	git(t, f.state, "push", "-u", "origin", branch)
	peer := filepath.Join(t.TempDir(), "peer")
	git(t, t.TempDir(), "clone", "--branch", branch, bare, peer)
	peerEntity := filepath.Join(peer, f.slug+".md")
	peerBytes := strings.Replace(string(readPrepareReviewFile(t, peerEntity)), "status: validation", "status: validation\nowner: peer", 1)
	writeFixtureFile(t, peerEntity, peerBytes)
	git(t, peer, "add", f.slug+".md")
	git(t, peer, "commit", "-m", "peer edit")
	git(t, peer, "push", "origin", branch)

	_, errOut, code := runPrepareReview(t, f, "--publish", "--question", "Approve?", "--artifact", ".spacedock-state/gate-task/selected/review.md", "--summary", "Ready.", "--recommendation", "Approve.", "--json")
	if code != 3 || !strings.Contains(errOut, "publication halted") {
		t.Fatalf("peer conflict must halt after local durability: exit=%d stderr=%q", code, errOut)
	}
	head := strings.TrimSpace(git(t, f.state, "rev-parse", "HEAD"))
	if names := strings.Fields(git(t, f.state, "show", "--name-only", "--pretty=format:", "HEAD")); len(names) != 3 {
		t.Fatalf("halted local commit lacks atomic binding-room unit: %v", names)
	}
	_, _, _ = runPrepareReview(t, f, "--publish", "--question", "Approve?", "--artifact", ".spacedock-state/gate-task/selected/review.md", "--summary", "Ready.", "--recommendation", "Approve.", "--json")
	if strings.TrimSpace(git(t, f.state, "rev-parse", "HEAD")) != head {
		t.Fatal("peer-conflict restart created another attempt commit")
	}
}

func TestGatePrepareReviewRejectsIdentityAndRestoresPreCommitFailure(t *testing.T) {
	f := newPrepareReviewFixture(t)
	before := append([]byte(nil), readPrepareReviewFile(t, f.entity)...)
	head := strings.TrimSpace(git(t, f.state, "rev-parse", "HEAD"))
	oldSync := prepareReviewSync
	prepareReviewSync = func(string, string, string, string) (bool, statesync.Outcome, error) {
		return false, statesync.Outcome{}, os.ErrPermission
	}
	t.Cleanup(func() { prepareReviewSync = oldSync })
	_, errOut, code := runPrepareReview(t, f, "--publish", "--question", "Approve?", "--artifact", ".spacedock-state/gate-task/selected/review.md", "--summary", "Ready.", "--recommendation", "Approve.", "--json")
	if code == 0 || !strings.Contains(errOut, "permission denied") {
		t.Fatalf("injected failure exit=%d stderr=%q", code, errOut)
	}
	if string(readPrepareReviewFile(t, f.entity)) != string(before) || strings.TrimSpace(git(t, f.state, "rev-parse", "HEAD")) != head {
		t.Fatal("pre-commit failure did not restore entity bytes and HEAD")
	}
	if matches, _ := filepath.Glob(filepath.Join(f.state, f.slug, "review", "validation", "briefing-*")); len(matches) != 0 {
		t.Fatalf("pre-commit failure retained room: %v", matches)
	}
	if _, errOut, code = runPrepareReview(t, f, "--publish", "--question", "Approve?", "--artifact", "README.md", "--summary", "Ready.", "--recommendation", "Approve.", "--json"); code == 0 || !strings.Contains(errOut, "not an inspect candidate") {
		t.Fatalf("non-candidate identity accepted exit=%d stderr=%q", code, errOut)
	}
}

func readPrepareReviewFile(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
