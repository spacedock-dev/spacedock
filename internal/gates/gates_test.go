package gates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestEightHistoryReplayPreservesApplicationsAndUnknownFields(t *testing.T) {
	var attempts []string
	for i := 1; i <= 8; i++ {
		previous := ""
		extra := ""
		if i > 1 {
			previous = "\n          previous-attempt: gate-attempt:design-" + itoa(i-1)
		}
		if i == 8 {
			extra = "          scope-amendment:\n            decision: keep the historical amendment\n"
		}
		attempts = append(attempts, ""+
			"        - id: gate-attempt:design-"+itoa(i)+"\n"+
			"          sequence: "+itoa(i)+previous+"\n"+
			"          state: closed\n"+
			"          briefing:\n"+
			"            id: briefing:design-"+itoa(i)+"\n"+
			"            digest: sha256:"+strings.Repeat(itoa(i), 64)+"\n"+
			"            note: RAW-FILE PIN legacy shaping record\n"+
			"          resolution:\n"+
			"            type: Resolution\n"+
			"            id: resolution:design-"+itoa(i)+"\n"+
			"            briefing: briefing:design-"+itoa(i)+"\n"+
			"            by: person:captain\n"+
			"            at: 2026-07-2"+itoa(i)+"T00:00:00Z\n"+
			"            decision: approve\n"+
			"          application:\n"+
			"            action: advance\n"+
			"            target-stage: implementation\n"+
			"            state: consumed\n")
		attempts[len(attempts)-1] += extra
	}
	body := "version: 1\ncurrent:\n  gate: gate:design\n  attempt: gate-attempt:design-8\nrecords:\n  - id: gate:design\n    stage: ideation\n    current-attempt: gate-attempt:design-8\n    attempts:\n" + strings.Join(attempts, "")
	var before yaml.Node
	if err := yaml.Unmarshal([]byte(body), &before); err != nil {
		t.Fatal(err)
	}
	var doc Document
	if err := before.Decode(&doc); err != nil {
		t.Fatal(err)
	}
	if err := Validate(&doc); err != nil {
		t.Fatalf("eight-attempt production-shaped history rejected: %v", err)
	}
	out, err := yaml.Marshal(&doc)
	if err != nil {
		t.Fatal(err)
	}
	var after Document
	if err := yaml.Unmarshal(out, &after); err != nil {
		t.Fatal(err)
	}
	for i, a := range after.Records[0].Attempts {
		if a.Application == nil {
			t.Fatalf("attempt %d application subtree was lost", i+1)
		}
		app, ok := a.Application.(map[string]any)
		if !ok || app["state"] != "consumed" {
			t.Fatalf("attempt %d application changed: %#v", i+1, a.Application)
		}
	}
	if _, ok := after.Records[0].Attempts[7].Extra["scope-amendment"]; !ok {
		t.Fatal("unknown historical attempt field was lost")
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

func itoa(i int) string { return string(rune('0' + i)) }
