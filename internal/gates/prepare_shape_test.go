package gates

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPreparedBriefingBytesArePinned is AC-2. Every byte of the published
// Briefing is literal here except the Git coordinates, which move with the
// fixture commits. A changed field, a changed key order, or a changed indent
// reds. The amended q0 preflight recomputes its digest from these exact bytes.
func TestPreparedBriefingBytesArePinned(t *testing.T) {
	workflow, _, entity, artifact, reference := prepareFixture(t, "flat")
	result, err := Prepare(entity, PrepareInput{
		WorkflowDir: workflow,
		Question:    "Should this gate advance?",
		Artifact:    artifact,
		Summary:     "Exact summary.",
		References:  []string{reference},
	})
	if err != nil {
		t.Fatal(err)
	}
	got := readPreparedFile(t, filepath.Join(result.Room, preparedBriefingLocator))

	var parsed struct {
		Artifacts []struct{ URI, Rev string } `json:"artifacts"`
		Context   []struct{ URI, Rev string } `json:"context"`
	}
	if err := json.Unmarshal(got, &parsed); err != nil {
		t.Fatal(err)
	}
	if len(parsed.Artifacts) != 1 || len(parsed.Context) != 1 {
		t.Fatalf("unexpected item counts: %#v", parsed)
	}
	want := fmt.Sprintf("{\n"+
		"  \"type\": \"Briefing\",\n"+
		"  \"version\": \"1\",\n"+
		"  \"id\": \"briefing:task:validation:attempt-1:revision-1\",\n"+
		"  \"question\": \"Should this gate advance?\",\n"+
		"  \"artifacts\": [\n"+
		"    {\n"+
		"      \"id\": \"artifact:task:validation:attempt-1:revision-1:item-1\",\n"+
		"      \"uri\": %q,\n"+
		"      \"rev\": %q,\n"+
		"      \"mediaType\": \"text/markdown\",\n"+
		"      \"summary\": \"Exact summary.\"\n"+
		"    }\n"+
		"  ],\n"+
		"  \"context\": [\n"+
		"    {\n"+
		"      \"type\": \"Reference\",\n"+
		"      \"id\": \"reference:task:validation:attempt-1:revision-1:item-2\",\n"+
		"      \"uri\": %q,\n"+
		"      \"mediaType\": \"application/json\",\n"+
		"      \"rev\": %q\n"+
		"    }\n"+
		"  ]\n"+
		"}\n",
		parsed.Artifacts[0].URI, parsed.Artifacts[0].Rev,
		parsed.Context[0].URI, parsed.Context[0].Rev)
	if string(got) != want {
		t.Fatalf("published Briefing bytes changed:\n--- got ---\n%s--- want ---\n%s", got, want)
	}

	// The digest the entity binds is the recomputation of these exact bytes.
	// That equality is what lets the q0 preflight drop its request read.
	digest, err := CanonicalDigest(got)
	if err != nil {
		t.Fatal(err)
	}
	doc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	if binding := doc.Records[0].Attempts[0].Briefing; binding.Digest != digest || digest != result.Digest {
		t.Fatalf("bound digest=%q recomputed=%q result=%q", binding.Digest, digest, result.Digest)
	}
}

// TestBindingShapesClassifyAndResolve is AC-5. Every binding shape that reaches
// the gate readers appears once here, with what the predicate makes of it and
// where the resolver sends it. Only the two live room names are prepared rooms.
// Drop the briefing.json test, so that the predicate matches any directory, and
// the archived-room case reds.
func TestBindingShapesClassifyAndResolve(t *testing.T) {
	root := t.TempDir()
	entity := filepath.Join(root, "task.md")
	if err := os.WriteFile(entity, []byte("---\nid: task\nstatus: validation\n---\n# Task\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dir := func(name string, holds ...string) {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.MkdirAll(path, 0o755); err != nil {
			t.Fatal(err)
		}
		for _, held := range holds {
			if err := os.WriteFile(filepath.Join(path, held), []byte("{}\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	dir("reserved", preparedBriefingLocator)
	dir("earlier", legacyPreparedLocator)
	dir("archived", archivedBriefingLocator)
	dir("emptied")
	if err := os.WriteFile(filepath.Join(root, "legacyfile"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name, ref string
		prepared  bool
		// resolves is the path boundBriefingPath must return, relative to root.
		// An empty value means the resolver must report an error.
		resolves string
	}{
		{"reserved name", "./reserved", true, "reserved/" + preparedBriefingLocator},
		{"earlier name", "./earlier", true, "earlier/" + legacyPreparedLocator},
		{"emptied prepared room", "./emptied", true, ""},
		{"archived room directory", "./archived", false, "archived"},
		{"legacy ref names the Briefing file", "./legacyfile", false, "legacyfile"},
		{"opaque provider ref", "subspace-room:3k-gate-design", false, "subspace-room:3k-gate-design"},
		{"absent room", "./missing", false, "missing"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			binding := Briefing{
				ID:      "briefing:task:validation:attempt-1:revision-1",
				Digest:  "sha256:" + strings.Repeat("1", 64),
				RoomRef: tc.ref,
			}
			if got := preparedRoomBinding(entity, binding); got != tc.prepared {
				t.Fatalf("preparedRoomBinding(%q) = %v, want %v", tc.ref, got, tc.prepared)
			}
			path, err := boundBriefingPath(entity, binding)
			if tc.resolves == "" {
				if err == nil || !strings.Contains(err.Error(), "resolve canonical Briefing locator") {
					t.Fatalf("emptied prepared room resolved %q err=%v", path, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolve %q: %v", tc.ref, err)
			}
			// A prepared room resolves through resolveBriefingLocator, which
			// resolves symlinks. An archived binding returns the ref as it
			// stands, exactly as it does today.
			want := filepath.Join(root, filepath.FromSlash(tc.resolves))
			if tc.prepared {
				resolved, symErr := filepath.EvalSymlinks(want)
				if symErr != nil {
					t.Fatal(symErr)
				}
				want = resolved
			}
			if path != want {
				t.Fatalf("%q resolved %q want %q", tc.ref, path, want)
			}
		})
	}
}

// TestRoomThatLostItsBriefingStaysAPreparedRoom pins why preparedRoomBinding
// tests for the archived name and not the reserved one. A room whose Briefing is
// deleted must keep failing loudly. The archived read path instead gives it the
// retained-authority skip, and gate record then closes the captain's approval
// over a Briefing that no longer exists.
func TestRoomThatLostItsBriefingStaysAPreparedRoom(t *testing.T) {
	workflow, _, entity, artifact, _ := prepareFixture(t, "folder")
	result, err := Prepare(entity, PrepareInput{
		WorkflowDir: workflow, Question: "Advance?", Artifact: artifact, Summary: "candidate",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(result.Room, preparedBriefingLocator)); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if err := RecordSemantic(entity, RecordInput{
		Decision: "approve", Actor: "person:captain", WorkflowDir: workflow,
	}); err == nil || !strings.Contains(err.Error(), "resolve canonical Briefing locator") {
		t.Fatalf("record over a deleted Briefing = %v", err)
	}
	after, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Fatal("refused record over a deleted Briefing changed entity bytes")
	}
}

// TestRetainedTwoFileRoomKeepsItsFullValidation is the other half of AC-5. The
// rooms prepared before this change hold two files and bind a request digest.
// They keep the read path, the entry set, and the request validation they have
// today, and nothing migrates them.
func TestRetainedTwoFileRoomKeepsItsFullValidation(t *testing.T) {
	workflow, state, entity, artifact, _ := prepareFixture(t, "folder")
	result, err := Prepare(entity, PrepareInput{
		WorkflowDir: workflow, Question: "Advance?", Artifact: artifact, Summary: "candidate",
	})
	if err != nil {
		t.Fatal(err)
	}
	requestPath := retainTwoFileRoom(t, entity, result)

	doc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	binding := doc.Records[0].Attempts[0].Briefing
	if binding.RequestDigest == "" || !preparedRoomBinding(entity, binding) {
		t.Fatalf("retained room binding = %#v", binding)
	}
	if err := validateRetainedAuthority(entity, workflow, doc); err != nil {
		t.Fatalf("retained two-file room lost its validation: %v", err)
	}
	// The resolver reads it through the frozen request locator, not by name.
	path, err := boundBriefingPath(entity, binding)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(path) != legacyPreparedLocator {
		t.Fatalf("retained room resolved to %q", path)
	}

	// A drifted request still refuses, and the refusal is byte-clean.
	body := readPreparedFile(t, requestPath)
	before := treeDigest(t, state)
	if err := os.WriteFile(requestPath, []byte(strings.Replace(string(body),
		`"actor": "person:captain"`, `"actor": "agent:other"`, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Withdraw(entity, WithdrawInput{WorkflowDir: workflow, Reason: "stale"}); err == nil ||
		!strings.Contains(err.Error(), "retained request.json") {
		t.Fatalf("drifted retained request withdrawal = %v", err)
	}
	if err := os.WriteFile(requestPath, body, 0o644); err != nil {
		t.Fatal(err)
	}
	if got := treeDigest(t, state); got != before {
		t.Fatal("drifted-request refusal changed the state tree")
	}

	// Restored, the two-file room still withdraws: its entry set is two files.
	if _, err := Withdraw(entity, WithdrawInput{WorkflowDir: workflow, Reason: "stale"}); err != nil {
		t.Fatalf("retained two-file room refused withdraw: %v", err)
	}
}

// retainTwoFileRoom rewrites a freshly prepared room into the shape preparation
// published before this change. The Briefing takes the earlier name, beside the
// request that names it, and the request digest is bound. The function returns
// the request path. Preparation never writes this shape again, and no room
// migrates, so the fixture builds it by hand.
func retainTwoFileRoom(t *testing.T, entity string, result PrepareResult) string {
	t.Helper()
	briefing := readPreparedFile(t, filepath.Join(result.Room, preparedBriefingLocator))
	if err := os.Remove(filepath.Join(result.Room, preparedBriefingLocator)); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(result.Room, legacyPreparedLocator), briefing, 0o644); err != nil {
		t.Fatal(err)
	}
	doc, oldNode, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	record := &doc.Records[0]
	attempt := &record.Attempts[0]
	request := gateRoomRequest{
		Type:     "spacedock-gate-presentation-request",
		Version:  "1",
		Gate:     record.ID,
		Attempt:  attempt.ID,
		Actor:    "person:captain",
		Approver: "person:captain",
	}
	request.Briefing.Locator = legacyPreparedLocator
	request.Briefing.ID = attempt.Briefing.ID
	request.Briefing.Digest = attempt.Briefing.Digest
	requestBytes, err := indentedJSON(request)
	if err != nil {
		t.Fatal(err)
	}
	requestPath := filepath.Join(result.Room, "request.json")
	if err := os.WriteFile(requestPath, requestBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	digest, err := CanonicalDigest(requestBytes)
	if err != nil {
		t.Fatal(err)
	}
	attempt.Briefing.RequestDigest = digest
	if err := writeDocument(entity, oldNode, doc); err != nil {
		t.Fatal(err)
	}
	return requestPath
}
