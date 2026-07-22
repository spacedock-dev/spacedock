package gates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestTwoGateMultipleAttemptReplayPreservesApplicationsAndUnknownFields(t *testing.T) {
	fixture := filepath.Join("testdata", "two-gate-eight-history.md")
	data, err := os.ReadFile(fixture)
	if err != nil {
		t.Fatal(err)
	}
	entity := filepath.Join(t.TempDir(), "entity.md")
	if err := os.WriteFile(entity, data, 0o644); err != nil {
		t.Fatal(err)
	}
	doc, before, err := Read(entity)
	if err != nil {
		t.Fatalf("two-gate contract fixture rejected: %v", err)
	}
	if len(doc.Records) != 2 {
		t.Fatalf("logical gates = %d, want 2", len(doc.Records))
	}
	histories := 0
	for _, record := range doc.Records {
		if len(record.Attempts) < 2 {
			t.Fatalf("gate %s has no re-entry history", record.ID)
		}
		histories += len(record.Attempts)
		for _, attempt := range record.Attempts {
			if attempt.Application == nil {
				t.Fatalf("attempt %s application subtree was lost", attempt.ID)
			}
		}
	}
	if histories != 8 {
		t.Fatalf("history count = %d, want 8", histories)
	}
	if got := CurrentSummary(doc); got.Gate != "gate:docs:dev:falsifiability-ladder:validation" || got.Attempt != "gate-attempt:z7cvbvdv-validation-3" {
		t.Fatalf("current pointers disagree: %#v", got)
	}
	if err := write(entity, doc); err != nil {
		t.Fatal(err)
	}
	after, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := after.Extra["fixture-purpose"]; !ok {
		t.Fatal("unknown document field was lost")
	}
	if _, ok := after.Records[0].Attempts[3].Extra["scope-amendment"]; !ok {
		t.Fatal("unknown historical attempt field was lost")
	}
	if _, ok := after.Records[1].Attempts[2].Resolution.Extra["provider-audit"]; !ok {
		t.Fatal("unknown historical resolution field was lost")
	}
	app, ok := after.Records[1].Attempts[2].Application.(map[string]any)
	if !ok || app["state"] != "pending" {
		t.Fatalf("nested current application changed: %#v", after.Records[1].Attempts[2].Application)
	}

	t.Run("pointer fork", func(t *testing.T) {
		fork := cloneDocument(t, doc)
		fork.Current.Attempt = "gate-attempt:z7cvbvdv-validation-2"
		if err := Validate(fork); err == nil || !strings.Contains(err.Error(), "pointer conflict") {
			t.Fatalf("pointer fork = %v, want fail closed", err)
		}
	})
	t.Run("history fork", func(t *testing.T) {
		fork := cloneDocument(t, doc)
		fork.Records[0].Attempts[3].PreviousAttempt = "gate-attempt:fork"
		if err := Validate(fork); err == nil || !strings.Contains(err.Error(), "previous-attempt") {
			t.Fatalf("history fork = %v, want fail closed", err)
		}
	})
	t.Run("frozen history fork", func(t *testing.T) {
		fork := cloneDocument(t, doc)
		fork.Records[1].Attempts[0].Application = map[string]any{"state": "rewritten"}
		if err := ValidateTransition(before, fork); err == nil || !strings.Contains(err.Error(), "frozen") {
			t.Fatalf("frozen fork = %v, want fail closed", err)
		}
	})
}

func TestEightProductionHistoryReplays(t *testing.T) {
	fixtures, err := filepath.Glob(filepath.Join("testdata", "production", "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(fixtures) != 8 {
		t.Fatalf("production fixture count = %d, want 8", len(fixtures))
	}
	for _, fixture := range fixtures {
		fixture := fixture
		t.Run(strings.TrimSuffix(filepath.Base(fixture), ".md"), func(t *testing.T) {
			data, err := os.ReadFile(fixture)
			if err != nil {
				t.Fatal(err)
			}
			entity := filepath.Join(t.TempDir(), "entity.md")
			if err := os.WriteFile(entity, data, 0o644); err != nil {
				t.Fatal(err)
			}
			doc, _, err := Read(entity)
			if err != nil {
				t.Fatalf("production history rejected: %v", err)
			}
			before, err := yaml.Marshal(doc)
			if err != nil {
				t.Fatal(err)
			}
			if err := write(entity, doc); err != nil {
				t.Fatal(err)
			}
			replayed, _, err := Read(entity)
			if err != nil {
				t.Fatalf("rewritten history rejected: %v", err)
			}
			after, err := yaml.Marshal(replayed)
			if err != nil {
				t.Fatal(err)
			}
			if string(after) != string(before) {
				t.Fatalf("production history changed during replay:\nbefore:\n%s\nafter:\n%s", before, after)
			}
		})
	}
}

func TestRebindCloseFreezeAndSupersedeLifecycle(t *testing.T) {
	entity := writeEntity(t, "status: ideation\n")
	briefingA := operationFile(t, `{"type":"Briefing","id":"provider:a","body":"A"}`)
	briefingB := operationFile(t, `{"type":"Briefing","id":"provider:b","body":"B"}`)
	briefingC := operationFile(t, `{"type":"Briefing","id":"provider:c","body":"C"}`)
	briefingD := operationFile(t, `{"type":"Briefing","id":"provider:d","body":"D"}`)
	digestA, _ := CanonicalDigest([]byte(readFile(t, briefingA)))
	digestB, _ := CanonicalDigest([]byte(readFile(t, briefingB)))
	digestC, _ := CanonicalDigest([]byte(readFile(t, briefingC)))
	digestD, _ := CanonicalDigest([]byte(readFile(t, briefingD)))

	open := operationFile(t, "operation: open\nexpected: {gate: '', attempt: '', briefing: '', digest: ''}\ngate-id: gate:lifecycle\nstage: ideation\nattempt-id: attempt:lifecycle-1\nbriefing: {id: briefing:a}\n")
	if err := Record(entity, open, briefingA); err != nil {
		t.Fatal(err)
	}
	assertCurrentBinding(t, entity, "attempt:lifecycle-1", "briefing:a", digestA, "open", 1)

	rebindB := operationFile(t, fmt.Sprintf("operation: rebind\nexpected: {gate: 'gate:lifecycle', attempt: 'attempt:lifecycle-1', briefing: 'briefing:a', digest: '%s'}\ngate-id: gate:lifecycle\nattempt-id: attempt:lifecycle-1\nbriefing: {id: briefing:b}\n", digestA))
	if err := Record(entity, rebindB, briefingB); err != nil {
		t.Fatal(err)
	}
	assertCurrentBinding(t, entity, "attempt:lifecycle-1", "briefing:b", digestB, "open", 1)

	rebindC := operationFile(t, fmt.Sprintf("operation: rebind\nexpected: {gate: 'gate:lifecycle', attempt: 'attempt:lifecycle-1', briefing: 'briefing:b', digest: '%s'}\ngate-id: gate:lifecycle\nattempt-id: attempt:lifecycle-1\nbriefing: {id: briefing:c}\n", digestB))
	if err := Record(entity, rebindC, briefingC); err != nil {
		t.Fatal(err)
	}
	assertCurrentBinding(t, entity, "attempt:lifecycle-1", "briefing:c", digestC, "open", 1)

	closeOp := operationFile(t, fmt.Sprintf("operation: close\nexpected: {gate: 'gate:lifecycle', attempt: 'attempt:lifecycle-1', briefing: 'briefing:c', digest: '%s'}\ngate-id: gate:lifecycle\nattempt-id: attempt:lifecycle-1\nresult:\n  briefing-digest: %s\n  authorized-by: person:captain\n  entries:\n    - type: Resolution\n      id: resolution:lifecycle-1\n      briefing: provider:c\n      by: person:captain\n      at: 2026-07-22T00:00:00Z\n      decision: approve\n", digestC, digestC))
	if err := Record(entity, closeOp, ""); err != nil {
		t.Fatal(err)
	}
	assertCurrentBinding(t, entity, "attempt:lifecycle-1", "briefing:c", digestC, "closed", 1)
	closed := readFile(t, entity)
	closedRebind := operationFile(t, fmt.Sprintf("operation: rebind\nexpected: {gate: 'gate:lifecycle', attempt: 'attempt:lifecycle-1', briefing: 'briefing:c', digest: '%s'}\ngate-id: gate:lifecycle\nattempt-id: attempt:lifecycle-1\nbriefing: {id: briefing:b}\n", digestC))
	if err := Record(entity, closedRebind, briefingB); err == nil || !strings.Contains(err.Error(), "frozen") {
		t.Fatalf("closed rebind = %v, want frozen refusal", err)
	}
	if got := readFile(t, entity); got != closed {
		t.Fatal("closed-attempt rebind mutated entity")
	}

	supersede := operationFile(t, fmt.Sprintf("operation: supersede\nexpected: {gate: 'gate:lifecycle', attempt: 'attempt:lifecycle-1', briefing: 'briefing:c', digest: '%s'}\ngate-id: gate:lifecycle\nattempt-id: attempt:lifecycle-2\nbriefing: {id: briefing:d}\n", digestC))
	if err := Record(entity, supersede, briefingD); err != nil {
		t.Fatal(err)
	}
	assertCurrentBinding(t, entity, "attempt:lifecycle-2", "briefing:d", digestD, "open", 2)
	doc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	first, second := doc.Records[0].Attempts[0], doc.Records[0].Attempts[1]
	if first.State != "closed" || first.Briefing.ID != "briefing:c" || first.Resolution == nil || second.PreviousAttempt != first.ID || second.Sequence != 2 {
		t.Fatalf("supersession did not preserve frozen lineage: first=%#v second=%#v", first, second)
	}
}

func TestMutableOpenAttemptCompatibility(t *testing.T) {
	resolution := &Resolution{ID: "resolution:a"}
	for _, tc := range []struct {
		name    string
		attempt Attempt
		want    bool
	}{
		{name: "explicit legacy open", attempt: Attempt{State: "open"}, want: true},
		{name: "minimal open", attempt: Attempt{}, want: true},
		{name: "resolution-bearing minimal", attempt: Attempt{Resolution: resolution}},
		{name: "explicit closed", attempt: Attempt{State: "closed", Resolution: resolution}},
		{name: "contradictory open with resolution", attempt: Attempt{State: "open", Resolution: resolution}},
		{name: "contradictory closed without resolution", attempt: Attempt{State: "closed"}},
		{name: "unknown explicit state", attempt: Attempt{State: "pending"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := mutableOpenAttempt(&tc.attempt); got != tc.want {
				t.Fatalf("mutableOpenAttempt(%#v) = %v, want %v", tc.attempt, got, tc.want)
			}
		})
	}
}

func TestRecordCloseNormalizesOnlyAfterDigestMatch(t *testing.T) {
	initial := "status: ideation\ngates:\n  version: 1\n  current: {gate: 'gate:design', attempt: 'gate-attempt:design-1'}\n  records:\n    - id: gate:design\n      stage: ideation\n      current-attempt: gate-attempt:design-1\n      attempts:\n        - id: gate-attempt:design-1\n          sequence: 1\n          state: open\n          briefing:\n            id: briefing:design-1\n            digest: sha256:" + strings.Repeat("a", 64) + "\n"
	entity := writeEntity(t, initial)
	badEntity := writeEntity(t, initial)
	beforeOutside := outsideGates(t, entity)
	op := operationFile(t, "operation: close\nexpected: {gate: 'gate:design', attempt: 'gate-attempt:design-1', briefing: 'briefing:design-1', digest: 'sha256:"+strings.Repeat("a", 64)+"'}\ngate-id: gate:design\nattempt-id: gate-attempt:design-1\nresult:\n  briefing-digest: sha256:"+strings.Repeat("a", 64)+"\n  authorized-by: person:captain\n  entries:\n    - type: Resolution\n      id: resolution:advisory\n      briefing: briefing:provider-envelope\n      by: agent:reviewer\n      at: 2026-07-22T00:00:00Z\n      decision: revise\n      reason: advisory\n    - type: Resolution\n      id: resolution:binding\n      briefing: briefing:provider-envelope\n      by: person:captain\n      at: 2026-07-22T00:01:00Z\n      decision: approve\n")
	if err := Record(entity, op, ""); err != nil {
		t.Fatal(err)
	}
	if afterOutside := outsideGates(t, entity); afterOutside != beforeOutside {
		t.Fatalf("recording changed bytes outside gates:\nbefore=%q\nafter=%q", beforeOutside, afterOutside)
	}
	doc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	got := doc.Records[0].Attempts[0].Resolution
	if got.ID != "resolution:binding" || got.Briefing != "briefing:design-1" {
		t.Fatalf("binding resolution not selected/normalized: %#v", got)
	}

	bad := operationFile(t, strings.ReplaceAll(readFile(t, op), strings.Repeat("a", 64), strings.Repeat("b", 64)))
	before := readFile(t, badEntity)
	if err := Record(badEntity, bad, ""); err == nil || !strings.Contains(err.Error(), "digest") {
		t.Fatalf("digest mismatch = %v, want refusal", err)
	}
	if got := readFile(t, badEntity); got != before {
		t.Fatal("failed close mutated the entity")
	}
}

func TestAdversarialWrapperFieldsStayOutsideCopiedResolution(t *testing.T) {
	digest := "sha256:" + strings.Repeat("a", 64)
	entity := writeEntity(t, "status: ideation\ngates:\n  version: 1\n  current: {gate: 'gate:design', attempt: 'attempt:design-1'}\n  records:\n    - id: gate:design\n      stage: ideation\n      current-attempt: attempt:design-1\n      attempts:\n        - id: attempt:design-1\n          sequence: 1\n          state: open\n          briefing: {id: 'briefing:design-1', digest: '"+digest+"'}\n")
	op := operationFile(t, "operation: close\nexpected: {gate: 'gate:design', attempt: 'attempt:design-1', briefing: 'briefing:design-1', digest: '"+digest+"'}\ngate-id: gate:design\nattempt-id: attempt:design-1\nresult:\n  briefing-digest: "+digest+"\n  authorized-by: person:captain\n  entries:\n    - type: Resolution\n      id: resolution:binding\n      briefing: provider:design\n      by: person:captain\n      at: 2026-07-22T00:00:00Z\n      decision: approve\n      stage: done\n      sequence: 99\n      briefing-change:\n        from: provider:old\n        to: provider:design\n      application:\n        action: advance\n        state: consumed\n      future-wrapper-field: never-durable\n")
	if err := Record(entity, op, ""); err != nil {
		t.Fatal(err)
	}
	doc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	attempt := doc.Records[0].Attempts[0]
	if attempt.Application != nil {
		t.Fatalf("provider application crossed into attempt application: %#v", attempt.Application)
	}
	for _, key := range []string{"stage", "sequence", "briefing-change", "application", "future-wrapper-field"} {
		if _, leaked := attempt.Resolution.Extra[key]; leaked {
			t.Fatalf("provider wrapper field %q leaked into durable Resolution: %#v", key, attempt.Resolution.Extra)
		}
	}
}

func TestRecordRefusesPointerConflictAndFrozenMutation(t *testing.T) {
	entity := writeEntity(t, "status: ideation\ngates:\n  version: 1\n  current: {gate: 'gate:a', attempt: 'attempt:a-1'}\n  records:\n    - id: gate:a\n      stage: ideation\n      current-attempt: attempt:a-1\n      attempts:\n        - id: attempt:a-1\n          sequence: 1\n          state: closed\n          briefing: {id: 'briefing:a-1', digest: 'sha256:"+strings.Repeat("1", 64)+"'}\n          resolution: {type: Resolution, id: 'resolution:a-1', briefing: 'briefing:a-1', by: 'person:captain', at: '2026-07-22T00:00:00Z', decision: approve}\n")
	conflict := operationFile(t, "operation: supersede\nexpected: {gate: 'gate:a', attempt: 'attempt:wrong', briefing: 'briefing:a-1', digest: 'sha256:"+strings.Repeat("1", 64)+"'}\ngate-id: gate:a\nattempt-id: attempt:a-2\nbriefing: {id: 'briefing:a-2'}\n")
	briefing := operationFile(t, `{"type":"Briefing","id":"provider:new"}`)
	if err := Record(entity, conflict, briefing); err == nil || !strings.Contains(err.Error(), "pointer conflict") {
		t.Fatalf("pointer conflict = %v", err)
	}

	var doc Document
	_, n, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	if err := n.Decode(&doc); err != nil {
		t.Fatal(err)
	}
	doc.Records[0].Attempts[0].Briefing.ID = "briefing:mutated"
	if err := ValidateTransition(n, &doc); err == nil || !strings.Contains(err.Error(), "frozen") {
		t.Fatalf("closed-attempt mutation = %v", err)
	}
}

func TestConcurrentWriterFailsClosed(t *testing.T) {
	entity := writeEntity(t, "status: ideation\n")
	if err := os.WriteFile(entity+".gates.lock", []byte("held"), 0o600); err != nil {
		t.Fatal(err)
	}
	op := operationFile(t, "operation: open\ngate-id: gate:a\nstage: ideation\nattempt-id: attempt:a-1\nbriefing: {id: 'briefing:a-1'}\n")
	briefing := operationFile(t, `{"type":"Briefing"}`)
	before := readFile(t, entity)
	if err := Record(entity, op, briefing); err == nil || !strings.Contains(err.Error(), "concurrent gate writer") {
		t.Fatalf("concurrent writer = %v, want refusal", err)
	}
	if got := readFile(t, entity); got != before {
		t.Fatal("lock contention changed entity")
	}
}

func TestDigestDomainsDivergeAndLegacyAccepted(t *testing.T) {
	pretty := []byte("{\n  \"b\": 2,\n  \"a\": 1\n}\n")
	compact := []byte(`{"a":1,"b":2}`)
	pc, err := CanonicalDigest(pretty)
	if err != nil {
		t.Fatal(err)
	}
	cc, err := CanonicalDigest(compact)
	if err != nil {
		t.Fatal(err)
	}
	if pc != cc {
		t.Fatalf("canonical digest changed with formatting: %s != %s", pc, cc)
	}
	if RawDigest(pretty) == RawDigest(compact) {
		t.Fatal("raw-file pins unexpectedly equal")
	}
	numberA, err := CanonicalDigest([]byte(`{"n":1e30}`))
	if err != nil {
		t.Fatal(err)
	}
	numberB, err := CanonicalDigest([]byte(`{"n":1e+30}`))
	if err != nil || numberA != numberB {
		t.Fatalf("RFC 8785 number normalization diverged: %s != %s (%v)", numberA, numberB, err)
	}
}

func TestPortableResolutionValidation(t *testing.T) {
	base := Result{BriefingDigest: "sha256:" + strings.Repeat("a", 64), AuthorizedBy: "person:captain"}
	tests := []struct {
		name    string
		entries []Entry
		wantErr string
	}{
		{"reasonless approve", []Entry{{Type: "Resolution", ID: "r", Briefing: "p", By: "person:captain", At: "now", Decision: "approve"}}, ""},
		{"reasonless revise", []Entry{{Type: "Resolution", ID: "r", Briefing: "p", By: "person:captain", At: "now", Decision: "revise"}}, "reason"},
		{"earlier annotation", []Entry{{Type: "Annotation", ID: "a", Briefing: "p"}, {Type: "Resolution", ID: "r", Briefing: "p", By: "person:captain", At: "now", Decision: "hold", Includes: []string{"a"}}}, ""},
		{"cross briefing include", []Entry{{Type: "Annotation", ID: "a", Briefing: "other"}, {Type: "Resolution", ID: "r", Briefing: "p", By: "person:captain", At: "now", Decision: "hold", Includes: []string{"a"}}}, "same Briefing"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := base
			got.Entries = tc.entries
			_, err := selectResolution(got, false)
			if tc.wantErr == "" && err != nil {
				t.Fatal(err)
			}
			if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
				t.Fatalf("err=%v, want %q", err, tc.wantErr)
			}
		})
	}
}

func writeEntity(t *testing.T, frontmatter string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "entity.md")
	if err := os.WriteFile(p, []byte("---\n"+frontmatter+"---\n# Entity\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func operationFile(t *testing.T, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "input")
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func outsideGates(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	root, start, end, err := frontmatterNode(b)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value != "gates" {
			continue
		}
		blockStart, blockEnd := start+root.Content[i].Line, end
		if i+2 < len(root.Content) {
			blockEnd = start + root.Content[i+2].Line
		}
		return strings.Join(append(append([]string{}, lines[:blockStart]...), lines[blockEnd:]...), "\n")
	}
	return string(b)
}

func cloneDocument(t *testing.T, doc *Document) *Document {
	t.Helper()
	b, err := yaml.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	var clone Document
	if err := yaml.Unmarshal(b, &clone); err != nil {
		t.Fatal(err)
	}
	return &clone
}

func assertCurrentBinding(t *testing.T, entity, attemptID, briefingID, digest, state string, histories int) {
	t.Helper()
	doc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Records) != 1 || len(doc.Records[0].Attempts) != histories {
		t.Fatalf("history count = %#v, want %d", doc.Records, histories)
	}
	a, err := currentAttempt(doc, "gate:lifecycle", attemptID)
	if err != nil {
		t.Fatal(err)
	}
	if a.Briefing.ID != briefingID || a.Briefing.Digest != digest || a.State != state {
		t.Fatalf("current binding = %#v, want attempt=%s briefing=%s digest=%s state=%s", a, attemptID, briefingID, digest, state)
	}
}
