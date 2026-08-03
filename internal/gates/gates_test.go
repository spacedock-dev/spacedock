package gates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestWithdrawalDefinesThirdValidatedFrozenAttemptState(t *testing.T) {
	withdrawn := eligibleDocument()
	attempt := &withdrawn.Records[0].Attempts[0]
	attempt.Briefing.RequestDigest = "sha256:" + strings.Repeat("2", 64)
	attempt.Resolution = nil
	attempt.Application = nil
	attempt.Withdrawal = &Withdrawal{
		By:     "agent:first-officer",
		At:     time.Date(2026, 7, 26, 11, 30, 0, 123456000, time.UTC).Format(time.RFC3339Nano),
		Reason: "Sprint re-scope replaced the reviewed candidate.",
	}
	if err := Validate(withdrawn); err != nil {
		t.Fatalf("valid withdrawn attempt rejected: %v", err)
	}
	if got := attemptState(attempt); got != "withdrawn" {
		t.Fatalf("attempt state = %q, want withdrawn", got)
	}
	summary := CurrentSummary(withdrawn)
	if summary.State != "withdrawn" || summary.Resolution != "" || summary.Decision != "" ||
		summary.Application != "" || summary.ApplicationState != "" {
		t.Fatalf("withdrawn summary leaked closure state: %#v", summary)
	}
	stages := []ReadinessStage{{Name: "ideation", Gate: true}, {Name: "implementation"}}
	if got := CurrentStageReadiness(withdrawn, "ideation", stages); got != "withdrawn-awaiting-prepare" {
		t.Fatalf("withdrawn readiness = %q", got)
	}

	for name, mutate := range map[string]func(*Attempt){
		"wrong actor":      func(a *Attempt) { a.Withdrawal.By = "person:captain" },
		"bad time":         func(a *Attempt) { a.Withdrawal.At = "now" },
		"blank reason":     func(a *Attempt) { a.Withdrawal.Reason = " \t" },
		"no request":       func(a *Attempt) { a.Briefing.RequestDigest = "" },
		"with resolution":  func(a *Attempt) { a.Resolution = eligibleDocument().Records[0].Attempts[0].Resolution },
		"with evidence":    func(a *Attempt) { a.ProviderEvidence = &ProviderEvidence{} },
		"with application": func(a *Attempt) { a.Application = &Application{} },
	} {
		t.Run(name, func(t *testing.T) {
			doc := cloneDocument(t, withdrawn)
			mutate(&doc.Records[0].Attempts[0])
			if err := Validate(doc); err == nil {
				t.Fatalf("invalid withdrawn shape %q accepted", name)
			}
		})
	}

	oldBytes, err := yaml.Marshal(withdrawn)
	if err != nil {
		t.Fatal(err)
	}
	var oldNode yaml.Node
	if err := yaml.Unmarshal(oldBytes, &oldNode); err != nil {
		t.Fatal(err)
	}
	next := cloneDocument(t, withdrawn)
	next.Records[0].Attempts[0].Withdrawal.Reason = "rewritten"
	if err := ValidateTransition(&oldNode, next); err == nil || !strings.Contains(err.Error(), "frozen withdrawn attempt") {
		t.Fatalf("withdrawn mutation = %v, want frozen refusal", err)
	}
}

func TestCurrentStageReadinessFailClosedTable(t *testing.T) {
	stages := []ReadinessStage{
		{Name: "ideation", Gate: true},
		{Name: "implementation"},
		{Name: "validation", Gate: true},
		{Name: "done", Terminal: true},
	}
	open := eligibleDocument()
	open.Records[0].Attempts[0].Resolution = nil
	open.Records[0].Attempts[0].Application = nil

	tests := []struct {
		name   string
		status string
		doc    *Document
		mutate func(*Document)
		want   string
	}{
		{name: "gate without selected attempt", status: "ideation", want: "validating"},
		{name: "open selected current attempt", status: "ideation", doc: open, want: "awaiting-captain"},
		{name: "approved nonterminal target", status: "ideation", doc: eligibleDocument(), want: "approved-awaiting-advance"},
		{name: "approved terminal target", status: "ideation", doc: eligibleDocument(), mutate: func(d *Document) {
			d.Records[0].Attempts[0].Application.TargetStage = "done"
		}, want: "approved-awaiting-merge"},
		{name: "feedback pending without application", status: "ideation", doc: eligibleDocument(), mutate: func(d *Document) {
			a := &d.Records[0].Attempts[0]
			a.Resolution.Decision, a.Resolution.Reason = "revise", "changes requested"
			a.Application = nil
		}, want: "feedback-pending"},
		{name: "older terminal approval cannot override newer rejection", status: "ideation", doc: eligibleDocument(), mutate: func(d *Document) {
			d.Records[0].Attempts[0].Application.TargetStage = "done"
			newer := cloneDocument(t, d).Records[0].Attempts[0]
			newer.ID, newer.Briefing.ID = "attempt:a-2", "briefing:a-2"
			newer.Resolution.ID, newer.Resolution.Briefing = "resolution:a-2", newer.Briefing.ID
			newer.Resolution.Decision, newer.Resolution.Reason = "revise", "changes requested"
			newer.Application = nil
			d.Records[0].Attempts = append(d.Records[0].Attempts, newer)
		}, want: "feedback-pending"},
		{name: "consumed approval", status: "ideation", doc: eligibleDocument(), mutate: setApplicationState("consumed"), want: "consumed"},
		{name: "superseded approval", status: "ideation", doc: eligibleDocument(), mutate: setApplicationState("superseded"), want: "superseded"},
		{name: "not applicable hold without application", status: "ideation", doc: eligibleDocument(), mutate: func(d *Document) {
			a := &d.Records[0].Attempts[0]
			a.Resolution.Decision, a.Resolution.Reason = "hold", "wait"
			a.Application = nil
		}, want: "not-applicable"},
		{name: "unknown target", status: "ideation", doc: eligibleDocument(), mutate: func(d *Document) {
			d.Records[0].Attempts[0].Application.TargetStage = "missing"
		}, want: "invalid"},
		{name: "stale selected old stage", status: "validation", doc: open, want: "validating"},
		{name: "ordinary stage", status: "implementation", doc: open, want: ""},
		{name: "terminal stage", status: "done", doc: open, want: ""},
		{name: "unknown stage", status: "missing", doc: open, want: ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var doc *Document
			if tc.doc != nil {
				doc = cloneDocument(t, tc.doc)
			}
			if tc.mutate != nil {
				tc.mutate(doc)
			}
			if got := CurrentStageReadiness(doc, tc.status, stages); got != tc.want {
				t.Fatalf("readiness = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDuplicateStageRecordsFailClosedAndPreserveBytes(t *testing.T) {
	doc := eligibleDocument()
	duplicate := cloneDocument(t, doc).Records[0]
	duplicate.ID = "gate:duplicate"
	doc.Records = append(doc.Records, duplicate)
	if got := CurrentStageReadiness(doc, "ideation", []ReadinessStage{{Name: "ideation", Gate: true}}); got != "invalid" {
		t.Fatalf("duplicate-stage readiness = %q, want invalid", got)
	}
	body, err := yaml.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	gatesBlock := "  " + strings.ReplaceAll(strings.TrimSuffix(string(body), "\n"), "\n", "\n  ") + "\n"
	entity := writeEntity(t, "status: ideation\ngates:\n"+gatesBlock)
	before := readFile(t, entity)
	if _, _, err := Read(entity); err == nil || !strings.Contains(err.Error(), "multiple logical gates claim workflow stage ideation") {
		t.Fatalf("duplicate-stage read error = %v, want ambiguity refusal", err)
	}
	if got := readFile(t, entity); got != before {
		t.Fatal("duplicate-stage refusal changed entity bytes")
	}
}

func TestPrototypeAndUnknownGateShapesFailClosed(t *testing.T) {
	base := "status: ideation\n" + canonicalOpenGates()
	cases := map[string]string{
		"record attempt pointer":   strings.Replace(base, "      stage: ideation\n", "      stage: ideation\n      current-attempt: attempt:a-1\n", 1),
		"attempt sequence":         strings.Replace(base, "        - id: attempt:a-1\n", "        - id: attempt:a-1\n          sequence: 1\n", 1),
		"attempt lineage":          strings.Replace(base, "        - id: attempt:a-1\n", "        - id: attempt:a-1\n          previous-attempt: attempt:a-0\n", 1),
		"attempt state":            strings.Replace(base, "        - id: attempt:a-1\n", "        - id: attempt:a-1\n          state: open\n", 1),
		"global attempt pointer":   strings.Replace(base, "  records:\n", "  current:\n    gate: gate:a\n  records:\n", 1),
		"unknown gates field":      strings.Replace(base, "  records:\n", "  shadow: prototype\n  records:\n", 1),
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
			if got := readFile(t, entity); got != before {
				t.Fatal("rejected prototype read changed entity")
			}
		})
	}
}

func TestDigestDomainFieldFailsClosed(t *testing.T) {
	entity := writeEntity(t, strings.Replace(canonicalOpenGates(), "            room-ref: ./room", "            room-ref: ./room\n            digest-domain: canonical-bytes", 1))
	before := readFile(t, entity)
	if _, _, err := Read(entity); err == nil || !strings.Contains(err.Error(), "field") {
		t.Fatalf("digest-domain read error = %v, want canonical-v1 refusal", err)
	}
	if got := readFile(t, entity); got != before {
		t.Fatal("rejected digest-domain read changed entity")
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
	doc.Records[0].Stage = ""
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
	if err := RecordSemantic(entity, RecordInput{Decision: "hold", Actor: "person:captain", Reason: "wait"}); err == nil || !strings.Contains(err.Error(), "concurrent gate writer") {
		t.Fatalf("concurrent writer = %v, want refusal", err)
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
	binding := Briefing{ID: "briefing:docs-dev:3k:validation:attempt-1:revision-1", Digest: "sha256:0a54f1baec0120c1c93523e6900a6ce28e025c570289e5dfa9835e28099042ac", RoomRef: "./room/briefing.json"}
	manifest, err := boundBriefingManifest(entity, binding)
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.Artifacts) != 3 {
		t.Fatalf("independent artifact inventory = %d, want 3", len(manifest.Artifacts))
	}
}

func TestAuthorityDocumentDecodersRejectRecursiveDuplicateMembers(t *testing.T) {
	duplicateBriefing := []byte(`{"type":"Briefing","version":"1","id":"briefing:a","question":"Review?","artifacts":[{"id":"artifact:a","uri":"a.md","rev":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","extra":{"owner":"one","owner":"two"}}]}`)
	if _, err := parseBriefingManifest(duplicateBriefing); err == nil || !strings.Contains(err.Error(), "duplicate JSON object member") {
		t.Fatalf("Briefing duplicate error=%v", err)
	}

	duplicateRequest := []byte(`{"type":"spacedock-gate-presentation-request","version":"1","gate":"gate:a","attempt":"attempt:a-1","briefing":{"locator":"gate.json","id":"briefing:a","digest":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","extra":{"pin":"one","pin":"two"}},"actor":"person:captain","approver":"person:captain"}`)
	if _, err := decodeGateRoomRequest(duplicateRequest); err == nil || !strings.Contains(err.Error(), "duplicate JSON object member") {
		t.Fatalf("request duplicate error=%v", err)
	}

	duplicateResult := []byte(`{"type":"review-v1-result","briefing":"briefing:a","artifact":{"id":"artifact:a","rev":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"resolution":{"type":"Resolution","id":"resolution:a","briefing":"briefing:a","by":"person:captain","at":"now","decision":"approve","extra":{"authority":"one","authority":"two"}},"annotations":[]}`)
	if _, err := decodeProviderResult(duplicateResult); err == nil || !strings.Contains(err.Error(), "duplicate JSON object member") {
		t.Fatalf("Result duplicate error=%v", err)
	}

	duplicateInventory := []byte(`{"items":[{"type":"Artifact","id":"artifact:a","rev":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","extra":{"source":"one","source":"two"}}]}`)
	if _, err := decodePresentedInventory(duplicateInventory); err == nil || !strings.Contains(err.Error(), "duplicate JSON object member") {
		t.Fatalf("inventory duplicate error=%v", err)
	}

	deep := strings.Repeat(`{"nested":`, maxAuthorityJSONDepth+1) + `null` + strings.Repeat(`}`, maxAuthorityJSONDepth+1)
	if err := rejectDuplicateMembers([]byte(deep)); err == nil || !strings.Contains(err.Error(), "nesting exceeds") {
		t.Fatalf("deep authority JSON error=%v", err)
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
	result.Annotations = append(result.Annotations, Annotation{Type: "Annotation", ID: "annotation:a", Briefing: result.Briefing})
	if err := verifyProviderResolution(&result); err != nil {
		t.Fatalf("compatible provider Annotation without by/at = %v", err)
	}
	result.Annotations[0].Briefing = "briefing:other"
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
	dir := t.TempDir()
	readme := "---\nstages:\n  states:\n    - name: ideation\n      initial: true\n    - name: implementation\n---\n# Workflow\n"
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(dir, "entity.md")
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
	return "gates:\n  version: 1\n  records:\n    - id: gate:a\n      stage: ideation\n      attempts:\n        - id: attempt:a-1\n          briefing:\n            id: briefing:a-1\n            digest: sha256:" + strings.Repeat("1", 64) + "\n            room-ref: ./room\n"
}

func canonicalClosedFrontmatter() string {
	return "status: ideation\n" + canonicalOpenGates() + "          resolution:\n            type: Resolution\n            id: resolution:a-1\n            briefing: briefing:a-1\n            by: person:captain\n            at: 2026-07-22T00:00:00Z\n            decision: approve\n"
}

func canonicalTwoGateFrontmatter() string {
	return "status: ideation\ntitle: Task\ngates:\n  version: 1\n  records:\n    - id: gate:docs-dev:3k:ideation\n      stage: ideation\n      attempts:\n        - id: gate-attempt:3k-ideation-9\n          briefing:\n            id: briefing:ideation:9\n            digest: sha256:" + strings.Repeat("1", 64) + "\n            room-ref: ./review/ideation/9\n          resolution:\n            type: Resolution\n            id: resolution:ideation:9\n            briefing: briefing:ideation:9\n            by: person:captain\n            at: 2026-07-22T00:00:00Z\n            decision: approve\n          application:\n            action: advance\n            target-stage: implementation\n            state: pending\n            blockers:\n              - id: blocker:preserve-me\n                state: unsatisfied\n    - id: gate:docs-dev:3k:validation\n      stage: validation\n      attempts:\n        - id: gate-attempt:3k-validation-1\n          briefing:\n            id: briefing:validation:1\n            digest: sha256:" + strings.Repeat("2", 64) + "\n            room-ref: ./review/validation/1\n          resolution:\n            type: Resolution\n            id: resolution:validation:1\n            briefing: briefing:validation:1\n            by: person:captain\n            at: 2026-07-22T00:00:00Z\n            decision: revise\n            reason: Re-enter ideation.\n"
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
