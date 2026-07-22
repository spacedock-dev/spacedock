package gates

import (
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
	after := doc
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

func TestEightProductionHistoriesSurviveTargetedSemanticWrite(t *testing.T) {
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
			beforeRecords := recordSourceSections(t, data)
			beforeOutside := outsideGates(t, entity)
			slug := strings.TrimSuffix(filepath.Base(fixture), ".md")
			briefing := operationFile(t, completeBriefing(`briefing:compat:`+slug+`:done:attempt-1:revision-1`, "compatibility"))
			if err := RecordBriefing(entity, briefing); err != nil {
				t.Fatal(err)
			}
			afterBytes, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			afterRecords := recordSourceSections(t, afterBytes)
			for id, before := range beforeRecords {
				if afterRecords[id] != before {
					t.Fatalf("legacy record %s changed during targeted write", id)
				}
			}
			if got := outsideGates(t, entity); got != beforeOutside {
				t.Fatal("targeted write changed bytes outside gates")
			}
			replayed, _, err := Read(entity)
			if err != nil {
				t.Fatalf("targeted history rejected: %v", err)
			}
			if len(replayed.Records) != len(doc.Records)+1 {
				t.Fatalf("records=%d, want one appended to %d", len(replayed.Records), len(doc.Records))
			}
			added := replayed.Records[len(replayed.Records)-1]
			attempt := added.Attempts[0]
			if added.Stage != "done" || added.CurrentAttempt != "" || attempt.Sequence != 0 || attempt.PreviousAttempt != "" || attempt.State != "" {
				t.Fatalf("new record is not minimal v1: %#v", added)
			}
		})
	}
}

func TestRebindCloseFreezeAndSupersedeLifecycle(t *testing.T) {
	entity := writeEntity(t, "status: ideation\n")
	briefingA := operationFile(t, completeBriefing("briefing:local:lifecycle:ideation:attempt-1:revision-1", "A"))
	briefingB := operationFile(t, completeBriefing("briefing:local:lifecycle:ideation:attempt-1:revision-2", "B"))
	briefingC := operationFile(t, completeBriefing("briefing:local:lifecycle:ideation:attempt-1:revision-3", "C"))
	briefingD := operationFile(t, completeBriefing("briefing:local:lifecycle:ideation:attempt-2:revision-1", "D"))
	digestA, _ := CanonicalDigest([]byte(readFile(t, briefingA)))
	digestB, _ := CanonicalDigest([]byte(readFile(t, briefingB)))
	digestC, _ := CanonicalDigest([]byte(readFile(t, briefingC)))
	digestD, _ := CanonicalDigest([]byte(readFile(t, briefingD)))

	if err := RecordBriefing(entity, briefingA); err != nil {
		t.Fatal(err)
	}
	assertCurrentBinding(t, entity, "gate:local:lifecycle:ideation", "gate-attempt:lifecycle-ideation-1", "briefing:local:lifecycle:ideation:attempt-1:revision-1", digestA, "open", 1)

	if err := RecordBriefing(entity, briefingB); err != nil {
		t.Fatal(err)
	}
	assertCurrentBinding(t, entity, "gate:local:lifecycle:ideation", "gate-attempt:lifecycle-ideation-1", "briefing:local:lifecycle:ideation:attempt-1:revision-2", digestB, "open", 1)

	if err := RecordBriefing(entity, briefingC); err != nil {
		t.Fatal(err)
	}
	assertCurrentBinding(t, entity, "gate:local:lifecycle:ideation", "gate-attempt:lifecycle-ideation-1", "briefing:local:lifecycle:ideation:attempt-1:revision-3", digestC, "open", 1)

	if err := RecordSemantic(entity, RecordInput{Decision: "approve", Actor: "person:captain"}); err != nil {
		t.Fatal(err)
	}
	assertCurrentBinding(t, entity, "gate:local:lifecycle:ideation", "gate-attempt:lifecycle-ideation-1", "briefing:local:lifecycle:ideation:attempt-1:revision-3", digestC, "closed", 1)
	closed := readFile(t, entity)
	if err := RecordBriefing(entity, briefingD); err != nil {
		t.Fatal(err)
	}
	assertCurrentBinding(t, entity, "gate:local:lifecycle:ideation", "gate-attempt:lifecycle-ideation-2", "briefing:local:lifecycle:ideation:attempt-2:revision-1", digestD, "open", 2)
	doc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	first, second := doc.Records[0].Attempts[0], doc.Records[0].Attempts[1]
	if first.State != "" || first.Resolution == nil || second.PreviousAttempt != "" || second.Sequence != 0 || second.State != "" {
		t.Fatalf("supersession did not preserve frozen lineage: first=%#v second=%#v", first, second)
	}
	if !strings.HasPrefix(readFile(t, entity), strings.TrimSuffix(closed, "---\n# Entity\n")) {
		t.Fatal("supersession rewrote the frozen closed attempt")
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

func TestFrozenMutationIsRejected(t *testing.T) {
	entity := writeEntity(t, "status: ideation\ngates:\n  version: 1\n  current: {gate: 'gate:a', attempt: 'attempt:a-1'}\n  records:\n    - id: gate:a\n      stage: ideation\n      current-attempt: attempt:a-1\n      attempts:\n        - id: attempt:a-1\n          sequence: 1\n          state: closed\n          briefing: {id: 'briefing:a-1', digest: 'sha256:"+strings.Repeat("1", 64)+"'}\n          resolution: {type: Resolution, id: 'resolution:a-1', briefing: 'briefing:a-1', by: 'person:captain', at: '2026-07-22T00:00:00Z', decision: approve}\n")
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
	briefing := operationFile(t, completeBriefing("briefing:local:a:ideation:attempt-1:revision-1", "lock"))
	before := readFile(t, entity)
	if err := RecordBriefing(entity, briefing); err == nil || !strings.Contains(err.Error(), "concurrent gate writer") {
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
	tests := []struct {
		name       string
		resolution Resolution
		wantErr    string
	}{
		{"reasonless approve", Resolution{Type: "Resolution", ID: "r", Briefing: "p", By: "person:captain", At: "now", Decision: "approve"}, ""},
		{"reasonless revise", Resolution{Type: "Resolution", ID: "r", Briefing: "p", By: "person:captain", At: "now", Decision: "revise"}, "reason"},
		{"included rationale", Resolution{Type: "Resolution", ID: "r", Briefing: "p", By: "person:captain", At: "now", Decision: "hold", Includes: []string{"annotation:a"}}, ""},
		{"unknown decision", Resolution{Type: "Resolution", ID: "r", Briefing: "p", By: "person:captain", At: "now", Decision: "later"}, "decision"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateResolution(&tc.resolution, "p")
			if tc.wantErr == "" && err != nil {
				t.Fatal(err)
			}
			if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
				t.Fatalf("err=%v, want %q", err, tc.wantErr)
			}
		})
	}
}

func TestProviderResolutionIncludesRequireSameBriefingAnnotation(t *testing.T) {
	result := providerResult{Briefing: "briefing:provider"}
	result.Resolution = Resolution{Type: "Resolution", ID: "resolution:r", Briefing: result.Briefing, By: "person:reviewer", At: "now", Decision: "hold", Includes: []string{"annotation:a"}}
	result.Annotations = append(result.Annotations, struct {
		Type     string `json:"type"`
		ID       string `json:"id"`
		Briefing string `json:"briefing"`
	}{Type: "Annotation", ID: "annotation:a", Briefing: "briefing:other"})
	if err := verifyProviderResolution(&result); err == nil || !strings.Contains(err.Error(), "same Briefing") {
		t.Fatalf("cross-Briefing include = %v, want refusal", err)
	}
	result.Annotations[0].Briefing = result.Briefing
	if err := verifyProviderResolution(&result); err != nil {
		t.Fatalf("same-Briefing Annotation rejected: %v", err)
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

func completeBriefing(id, question string) string {
	return `{"type":"Briefing","version":"1","id":"` + id + `","question":"` + question + `","artifacts":[{"id":"artifact:primary","uri":"artifact.md","rev":"sha256:` + strings.Repeat("a", 64) + `"}]}`
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

func recordSourceSections(t *testing.T, data []byte) map[string]string {
	t.Helper()
	root, fmStart, fmEnd, err := frontmatterNode(data)
	if err != nil {
		t.Fatal(err)
	}
	_, gatesEnd, gatesNode, err := mappingPairRange(root, "gates", fmStart, fmEnd)
	if err != nil {
		t.Fatal(err)
	}
	_, recordsEnd, records, err := mappingPairRange(gatesNode, "records", fmStart, gatesEnd)
	if err != nil {
		t.Fatal(err)
	}
	lines := normalizedLines(data)
	sections := map[string]string{}
	for i, record := range records.Content {
		id := mappingValue(record, "id")
		if id == nil {
			t.Fatal("record fixture has no id")
		}
		start, end := fmStart+record.Line, recordsEnd
		if i+1 < len(records.Content) {
			end = fmStart + records.Content[i+1].Line
		}
		sections[id.Value] = strings.Join(lines[start:end], "\n")
	}
	return sections
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

func assertCurrentBinding(t *testing.T, entity, gateID, attemptID, briefingID, digest, state string, histories int) {
	t.Helper()
	doc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Records) != 1 || len(doc.Records[0].Attempts) != histories {
		t.Fatalf("history count = %#v, want %d", doc.Records, histories)
	}
	record := findRecord(doc, gateID)
	if record == nil {
		t.Fatalf("gate %s not found", gateID)
	}
	a, err := recordAttempt(record, recordCurrentAttempt(record))
	if err != nil {
		t.Fatal(err)
	}
	if a.ID != attemptID || a.Briefing.ID != briefingID || a.Briefing.Digest != digest || attemptState(a) != state {
		t.Fatalf("current binding = %#v, want attempt=%s briefing=%s digest=%s state=%s", a, attemptID, briefingID, digest, state)
	}
}
