package gates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCanonicalLifecycleRebindCloseFreezeAndSupersede(t *testing.T) {
	entity := writeEntity(t, "status: ideation\ntitle: Keep me quoted\n")
	outside := outsideGates(t, entity)
	briefingA := operationFile(t, completeBriefing("briefing:local:lifecycle:ideation:attempt-1:revision-1", "A"))
	briefingB := operationFile(t, completeBriefing("briefing:local:lifecycle:ideation:attempt-1:revision-2", "B"))
	briefingC := operationFile(t, completeBriefing("briefing:local:lifecycle:ideation:attempt-2:revision-1", "C"))

	if err := RecordBriefing(entity, briefingA); err != nil {
		t.Fatal(err)
	}
	assertCurrentBinding(t, entity, "gate:local:lifecycle:ideation", "gate-attempt:lifecycle-ideation-1", "briefing:local:lifecycle:ideation:attempt-1:revision-1", "open", 1)
	if err := RecordBriefing(entity, briefingB); err != nil {
		t.Fatal(err)
	}
	assertCurrentBinding(t, entity, "gate:local:lifecycle:ideation", "gate-attempt:lifecycle-ideation-1", "briefing:local:lifecycle:ideation:attempt-1:revision-2", "open", 1)
	if err := RecordSemantic(entity, RecordInput{Decision: "approve", Actor: "person:captain"}); err != nil {
		t.Fatal(err)
	}
	doc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	frozen := marshalAttempt(t, doc.Records[0].Attempts[0])
	if err := RecordBriefing(entity, briefingC); err != nil {
		t.Fatal(err)
	}
	assertCurrentBinding(t, entity, "gate:local:lifecycle:ideation", "gate-attempt:lifecycle-ideation-2", "briefing:local:lifecycle:ideation:attempt-2:revision-1", "open", 2)
	doc, _, err = Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	if got := marshalAttempt(t, doc.Records[0].Attempts[0]); got != frozen {
		t.Fatal("successor write changed the frozen closed attempt")
	}
	if got := outsideGates(t, entity); got != outside {
		t.Fatal("canonical gates writer changed unrelated frontmatter or body bytes")
	}
}

func TestCanonicalCrossGateReentryPreservesFrozenApplication(t *testing.T) {
	entity := writeEntity(t, canonicalTwoGateFrontmatter())
	doc, oldNode, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	frozenIdeation := marshalAttempt(t, doc.Records[0].Attempts[0])
	frozenValidation := marshalAttempt(t, doc.Records[1].Attempts[0])
	outside := outsideGates(t, entity)
	briefing := operationFile(t, completeBriefing("briefing:docs-dev:3k:ideation:attempt-10:revision-18", "re-enter"))
	if err := RecordBriefing(entity, briefing); err != nil {
		t.Fatal(err)
	}
	doc, _, err = Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Current.Gate != "gate:docs-dev:3k:ideation" || len(doc.Records[0].Attempts) != 2 {
		t.Fatalf("cross-gate successor not selected: %#v", doc)
	}
	if marshalAttempt(t, doc.Records[0].Attempts[0]) != frozenIdeation || marshalAttempt(t, doc.Records[1].Attempts[0]) != frozenValidation {
		t.Fatal("cross-gate re-entry changed a frozen closure or its opaque application")
	}
	if got := outsideGates(t, entity); got != outside {
		t.Fatal("cross-gate re-entry changed bytes outside gates")
	}
	mutated := cloneDocument(t, doc)
	mutated.Records[0].Attempts[0].Application = map[string]any{"state": "rewritten"}
	if err := ValidateTransition(oldNode, mutated); err == nil || !strings.Contains(err.Error(), "frozen") {
		t.Fatalf("application mutation = %v, want frozen refusal", err)
	}
}

func TestPrototypeAndUnknownGateShapesFailClosed(t *testing.T) {
	base := "status: ideation\n" + canonicalOpenGates()
	cases := map[string]string{
		"global attempt pointer":   strings.Replace(base, "    gate: gate:a\n", "    gate: gate:a\n    attempt: attempt:a-1\n", 1),
		"record attempt pointer":   strings.Replace(base, "      stage: ideation\n", "      stage: ideation\n      current-attempt: attempt:a-1\n", 1),
		"attempt sequence":         strings.Replace(base, "        - id: attempt:a-1\n", "        - id: attempt:a-1\n          sequence: 1\n", 1),
		"attempt lineage":          strings.Replace(base, "        - id: attempt:a-1\n", "        - id: attempt:a-1\n          previous-attempt: attempt:a-0\n", 1),
		"attempt state":            strings.Replace(base, "        - id: attempt:a-1\n", "        - id: attempt:a-1\n          state: open\n", 1),
		"unknown gates field":      strings.Replace(base, "  current:\n", "  shadow: prototype\n  current:\n", 1),
		"unknown record field":     strings.Replace(base, "      stage: ideation\n", "      stage: ideation\n      note: prototype\n", 1),
		"unknown attempt field":    strings.Replace(base, "        - id: attempt:a-1\n", "        - id: attempt:a-1\n          note: prototype\n", 1),
		"unknown briefing field":   strings.Replace(base, "            room-ref: ./room\n", "            room-ref: ./room\n            note: prototype\n", 1),
		"unknown resolution field": strings.Replace(canonicalClosedFrontmatter(), "            decision: approve\n", "            decision: approve\n            provider-audit: prototype\n", 1),
	}
	for name, frontmatter := range cases {
		t.Run(name, func(t *testing.T) {
			entity := writeEntity(t, frontmatter)
			before := readFile(t, entity)
			if _, _, err := Read(entity); err == nil || !strings.Contains(err.Error(), "field") {
				t.Fatalf("prototype shape read error = %v, want unknown-field refusal", err)
			}
			briefing := operationFile(t, completeBriefing("briefing:a:ideation:attempt-1:revision-2", "reject"))
			if err := RecordBriefing(entity, briefing); err == nil {
				t.Fatal("prototype shape was writable")
			}
			if got := readFile(t, entity); got != before {
				t.Fatal("rejected prototype write changed entity")
			}
		})
	}
}

func TestWriterCASValidationAtomicityAndLock(t *testing.T) {
	entity := writeEntity(t, "status: ideation\n"+canonicalOpenGates())
	doc, expected, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	external := strings.Replace(readFile(t, entity), "room-ref: ./room", "room-ref: ./other", 1)
	if err := os.WriteFile(entity, []byte(external), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeDocument(entity, expected, doc); err == nil || !strings.Contains(err.Error(), "changed during") {
		t.Fatalf("stale write = %v, want CAS refusal", err)
	}
	if got := readFile(t, entity); got != external {
		t.Fatal("CAS refusal replaced external state")
	}

	doc, expected, err = Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	doc.Current.Gate = "gate:missing"
	before := readFile(t, entity)
	if err := writeDocument(entity, expected, doc); err == nil {
		t.Fatal("invalid rebuilt document was accepted")
	}
	if got := readFile(t, entity); got != before {
		t.Fatal("invalid rebuild changed entity")
	}

	if err := os.WriteFile(entity+".gates.lock", []byte("held"), 0o600); err != nil {
		t.Fatal(err)
	}
	briefing := operationFile(t, completeBriefing("briefing:a:ideation:attempt-1:revision-2", "lock"))
	if err := RecordBriefing(entity, briefing); err == nil || !strings.Contains(err.Error(), "concurrent gate writer") {
		t.Fatalf("concurrent writer = %v, want refusal", err)
	}
}

func TestWriterPreservesMixedLineEndingsOutsideGates(t *testing.T) {
	entity := filepath.Join(t.TempDir(), "entity.md")
	original := "---\r\nstatus: ideation\r\ntitle: Mixed\r\n---\r\n# Entity\nBody keeps LF.\n"
	if err := os.WriteFile(entity, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	briefing := operationFile(t, completeBriefing("briefing:mixed:ideation:attempt-1:revision-1", "mixed"))
	if err := RecordBriefing(entity, briefing); err != nil {
		t.Fatal(err)
	}
	if got := outsideGates(t, entity); got != original {
		t.Fatalf("mixed line endings outside gates changed:\nwant=%q\n got=%q", original, got)
	}
}

func TestExactCanonicalBriefingIsIndependentAssociationInventory(t *testing.T) {
	room := filepath.Join(t.TempDir(), "room")
	if err := os.MkdirAll(room, 0o755); err != nil {
		t.Fatal(err)
	}
	bytes, err := os.ReadFile(filepath.Join("testdata", "exact-validation-briefing.json"))
	if err != nil {
		t.Fatal(err)
	}
	if got, err := CanonicalDigest(bytes); err != nil || got != "sha256:0a54f1baec0120c1c93523e6900a6ce28e025c570289e5dfa9835e28099042ac" {
		t.Fatalf("canonical fixture digest = %s (%v)", got, err)
	}
	if err := os.WriteFile(filepath.Join(room, "briefing.json"), bytes, 0o644); err != nil {
		t.Fatal(err)
	}
	entity := filepath.Join(filepath.Dir(room), "entity.md")
	binding := Briefing{ID: "briefing:docs-dev:3k:validation:attempt-1:revision-1", Digest: "sha256:0a54f1baec0120c1c93523e6900a6ce28e025c570289e5dfa9835e28099042ac", DigestDomain: "canonical-bytes", RoomRef: "./room"}
	manifest, err := boundBriefingManifest(entity, binding)
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.Artifacts) != 3 {
		t.Fatalf("independent artifact inventory = %d, want 3", len(manifest.Artifacts))
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
	if err := os.WriteFile(p, []byte("---\n"+frontmatter+"---\n# Entity\nBody keeps   spaces.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func operationFile(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "briefing.json")
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func completeBriefing(id, question string) string {
	return `{"type":"Briefing","version":"1","id":"` + id + `","question":"` + question + `","artifacts":[{"id":"artifact:primary","uri":"artifact.md","rev":"sha256:` + strings.Repeat("a", 64) + `"}]}`
}

func canonicalOpenGates() string {
	return "gates:\n  version: 1\n  current:\n    gate: gate:a\n  records:\n    - id: gate:a\n      stage: ideation\n      attempts:\n        - id: attempt:a-1\n          briefing:\n            id: briefing:a-1\n            digest: sha256:" + strings.Repeat("1", 64) + "\n            digest-domain: canonical-bytes\n            room-ref: ./room\n"
}

func canonicalClosedFrontmatter() string {
	return "status: ideation\n" + canonicalOpenGates() + "          resolution:\n            type: Resolution\n            id: resolution:a-1\n            briefing: briefing:a-1\n            by: person:captain\n            at: 2026-07-22T00:00:00Z\n            decision: approve\n"
}

func canonicalTwoGateFrontmatter() string {
	return "status: ideation\ntitle: Task\ngates:\n  version: 1\n  current:\n    gate: gate:docs-dev:3k:validation\n  records:\n    - id: gate:docs-dev:3k:ideation\n      stage: ideation\n      attempts:\n        - id: gate-attempt:3k-ideation-9\n          briefing:\n            id: briefing:ideation:9\n            digest: sha256:" + strings.Repeat("1", 64) + "\n            digest-domain: raw-file-pin\n            room-ref: ./review/ideation/9\n          resolution:\n            type: Resolution\n            id: resolution:ideation:9\n            briefing: briefing:ideation:9\n            by: person:captain\n            at: 2026-07-22T00:00:00Z\n            decision: approve\n          application:\n            state: pending\n            blockers:\n              - preserve me\n    - id: gate:docs-dev:3k:validation\n      stage: validation\n      attempts:\n        - id: gate-attempt:3k-validation-1\n          briefing:\n            id: briefing:validation:1\n            digest: sha256:" + strings.Repeat("2", 64) + "\n            digest-domain: raw-file-pin\n            room-ref: ./review/validation/1\n          resolution:\n            type: Resolution\n            id: resolution:validation:1\n            briefing: briefing:validation:1\n            by: person:captain\n            at: 2026-07-22T00:00:00Z\n            decision: revise\n            reason: Re-enter ideation.\n"
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
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value != "gates" {
			continue
		}
		blockStart, blockEnd := start+root.Content[i].Line, end
		if i+2 < len(root.Content) {
			blockEnd = start + root.Content[i+2].Line
		}
		startByte, endByte := lineOffset(b, blockStart), lineOffset(b, blockEnd)
		return string(append(append([]byte{}, b[:startByte]...), b[endByte:]...))
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

func marshalAttempt(t *testing.T, attempt Attempt) string {
	t.Helper()
	b, err := yaml.Marshal(attempt)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func assertCurrentBinding(t *testing.T, entity, gateID, attemptID, briefingID, state string, histories int) {
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
	a := &record.Attempts[len(record.Attempts)-1]
	if a.ID != attemptID || a.Briefing.ID != briefingID || attemptState(a) != state {
		t.Fatalf("current binding = %#v, want attempt=%s briefing=%s state=%s", a, attemptID, briefingID, state)
	}
}
