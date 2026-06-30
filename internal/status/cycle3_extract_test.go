// ABOUTME: Cycle-3 status --read additions — structured stages array (AC-3),
// ABOUTME: --fields projection (AC-6), --checklist (AC-1), --ac-scan (AC-2), loud failure (AC-5).
package status

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// devReadmePath returns the repo's real docs/dev/README.md — the external oracle
// for AC-3/AC-4 (its actual stages: block is the ground truth, so flipping a flag
// in the README would flip the emitted field, never a frozen byte golden).
func devReadmePath(t *testing.T) string {
	t.Helper()
	// internal/status -> repo root is ../..
	p, err := filepath.Abs(filepath.Join("..", "..", "docs", "dev", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("real dev README not found at %s: %v", p, err)
	}
	return p
}

// readEnvelope is the --read --json shape carrying the cycle-3 additions: the
// existing path/total_lines/frontmatter/headings plus the new stages array.
type readEnvelope struct {
	Command     string              `json:"command"`
	Path        string              `json:"path"`
	TotalLines  string              `json:"total_lines"`
	Frontmatter map[string]string   `json:"frontmatter"`
	Headings    []map[string]string `json:"headings"`
	Stages      []map[string]string `json:"stages"`
}

// TestReadStagesArray (AC-3) asserts status --read <README> --json surfaces the
// nested stages: taxonomy as a structured array of ordered objects matching the
// REAL docs/dev/README.md field-by-field. The README is the oracle: its actual
// stages: block (backlog initial+gate, ideation gate, implementation worktree,
// validation worktree+fresh+feedback-to+gate, done terminal) is the expected set.
func TestReadStagesArray(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	readme := devReadmePath(t)

	out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", readme, "--json")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	var doc readEnvelope
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}

	// The expected stages, in declaration order, with every leaf a string. Only
	// the flags the README actually sets are present; an unset bool flag is the
	// "false" string for the typed gate/terminal/initial/worktree fields (always
	// emitted), and the optional feedback-to/fresh keys appear only when set.
	type want struct {
		name       string
		worktree   string
		gate       string
		terminal   string
		initial    string
		feedbackTo string // "" means the key should be absent
		fresh      string // "" means absent
	}
	expected := []want{
		{name: "backlog", worktree: "false", gate: "true", terminal: "false", initial: "true"},
		{name: "ideation", worktree: "false", gate: "true", terminal: "false", initial: "false"},
		{name: "implementation", worktree: "true", gate: "false", terminal: "false", initial: "false"},
		{name: "validation", worktree: "true", gate: "true", terminal: "false", initial: "false", feedbackTo: "implementation", fresh: "true"},
		{name: "done", worktree: "false", gate: "false", terminal: "true", initial: "false"},
	}

	if len(doc.Stages) != len(expected) {
		t.Fatalf("stages count = %d, want %d\n%s", len(doc.Stages), len(expected), out)
	}
	for i, w := range expected {
		s := doc.Stages[i]
		if s["name"] != w.name {
			t.Fatalf("stages[%d].name = %q, want %q", i, s["name"], w.name)
		}
		check := map[string]string{
			"worktree": w.worktree, "gate": w.gate, "terminal": w.terminal, "initial": w.initial,
		}
		for k, v := range check {
			if s[k] != v {
				t.Errorf("stages[%d=%s].%s = %q, want %q", i, w.name, k, s[k], v)
			}
		}
		// Optional keys: present-with-value when the README sets them, absent otherwise.
		if w.feedbackTo == "" {
			if _, present := s["feedback-to"]; present {
				t.Errorf("stages[%d=%s] has feedback-to=%q, want absent", i, w.name, s["feedback-to"])
			}
		} else if s["feedback-to"] != w.feedbackTo {
			t.Errorf("stages[%d=%s].feedback-to = %q, want %q", i, w.name, s["feedback-to"], w.feedbackTo)
		}
		if w.fresh == "" {
			if _, present := s["fresh"]; present {
				t.Errorf("stages[%d=%s] has fresh=%q, want absent", i, w.name, s["fresh"])
			}
		} else if s["fresh"] != w.fresh {
			t.Errorf("stages[%d=%s].fresh = %q, want %q", i, w.name, s["fresh"], w.fresh)
		}
		// Every leaf is a string (the all-strings contract): no value is a JSON
		// bool/number. json.Unmarshal into map[string]string would have failed
		// already on a non-string, so reaching here proves it.
	}
}

// TestReadStagesArrayNoRegression (AC-3) asserts adding the stages array leaves
// the pre-existing default --read output untouched: the flat frontmatter still
// carries "stages":"" (the flattened scalar), and headings/total_lines for the
// README are intact.
func TestReadStagesArrayNoRegression(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	readme := devReadmePath(t)

	out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", readme, "--json")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	var doc readEnvelope
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	// The flat frontmatter object STILL flattens the nested stages: block to "".
	if v, ok := doc.Frontmatter["stages"]; !ok || v != "" {
		t.Fatalf("flat frontmatter[stages] = %q (present=%v), want \"\" present — the flat map is unchanged", v, ok)
	}
	// The README has headings and a positive line count; their presence is the
	// default-read contract, which the stages addition must not disturb.
	if len(doc.Headings) == 0 {
		t.Fatal("headings empty — the default heading map regressed")
	}
	if doc.TotalLines == "" || doc.TotalLines == "0" {
		t.Fatalf("total_lines = %q, want the README's real line count", doc.TotalLines)
	}
}

// TestReadStagesArrayAbsentForPlainFile (AC-3) asserts a markdown file with NO
// stages: block emits no stages array at all (the array is keyed on the block
// existing, not always-present) — the section-reader fixture has no stages:.
func TestReadStagesArrayAbsentForPlainFile(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixturePath(t), "--json")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	if strings.Contains(out, "\"stages\":[") {
		t.Fatalf("plain fixture (no stages: block) emitted a stages array: %s", out)
	}
}

// fieldsFixturePath returns the committed entity whose known frontmatter
// (id/title/status/verdict/score/source) is the AC-6 projection oracle.
func fieldsFixturePath(t *testing.T) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join("testdata", "section-reader", "fields-fixture.md"))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// frontmatterKeyOrder returns the frontmatter object's keys in emission order by
// walking the raw JSON with a token decoder, so a projection's key ORDER (not
// just membership) is assertable.
func frontmatterKeyOrder(t *testing.T, out string) []string {
	t.Helper()
	var top map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &top); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	raw, ok := top["frontmatter"]
	if !ok {
		t.Fatalf("no frontmatter object in %s", out)
	}
	dec := json.NewDecoder(strings.NewReader(string(raw)))
	if _, err := dec.Token(); err != nil { // opening '{'
		t.Fatalf("decode frontmatter: %v", err)
	}
	var keys []string
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("decode key: %v", err)
		}
		keys = append(keys, tok.(string))
		var v json.RawMessage
		if err := dec.Decode(&v); err != nil {
			t.Fatalf("decode value: %v", err)
		}
	}
	return keys
}

// TestReadFieldsProjection (AC-6) asserts status --read <entity> --fields k1,k2
// projects the frontmatter to exactly the named keys in named order, omitting the
// rest, reusing the existing --fields flag's semantics. With no --fields, the
// whole frontmatter map is returned (unchanged from today). A requested-but-absent
// key projects an empty string, matching entityJSONObj's behavior.
func TestReadFieldsProjection(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	fixture := fieldsFixturePath(t)

	t.Run("projects exactly the named keys in order", func(t *testing.T) {
		out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixture, "--fields", "status,verdict", "--json")
		if code != 0 {
			t.Fatalf("exit=%d stderr=%q", code, stderr)
		}
		keys := frontmatterKeyOrder(t, out)
		want := []string{"status", "verdict"}
		if len(keys) != len(want) {
			t.Fatalf("frontmatter keys = %v, want %v", keys, want)
		}
		for i, k := range want {
			if keys[i] != k {
				t.Fatalf("frontmatter key[%d] = %q, want %q (order matters)", i, keys[i], k)
			}
		}
		var doc readEnvelope
		json.Unmarshal([]byte(out), &doc)
		if _, present := doc.Frontmatter["id"]; present {
			t.Fatalf("unrequested key id present in projection: %v", doc.Frontmatter)
		}
		if doc.Frontmatter["status"] != "validation" || doc.Frontmatter["verdict"] != "approve" {
			t.Fatalf("projected values wrong: %v", doc.Frontmatter)
		}
	})

	t.Run("no --fields returns the whole map", func(t *testing.T) {
		out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixture, "--json")
		if code != 0 {
			t.Fatalf("exit=%d stderr=%q", code, stderr)
		}
		var doc readEnvelope
		if err := json.Unmarshal([]byte(out), &doc); err != nil {
			t.Fatalf("not JSON: %v\n%s", err, out)
		}
		for _, k := range []string{"id", "title", "status", "verdict", "score", "source"} {
			if _, present := doc.Frontmatter[k]; !present {
				t.Fatalf("no-fields read dropped key %q: %v", k, doc.Frontmatter)
			}
		}
	})

	t.Run("requested-but-absent key projects empty string", func(t *testing.T) {
		out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixture, "--fields", "status,nonesuch", "--json")
		if code != 0 {
			t.Fatalf("exit=%d stderr=%q", code, stderr)
		}
		keys := frontmatterKeyOrder(t, out)
		want := []string{"status", "nonesuch"}
		if len(keys) != len(want) || keys[0] != want[0] || keys[1] != want[1] {
			t.Fatalf("frontmatter keys = %v, want %v (missing key projects empty, not dropped)", keys, want)
		}
		var doc readEnvelope
		json.Unmarshal([]byte(out), &doc)
		if v := doc.Frontmatter["nonesuch"]; v != "" {
			t.Fatalf("absent key nonesuch = %q, want empty string (matching entityJSONObj)", v)
		}
	})
}

// interleavedFixturePath returns the committed fixture with INTERLEAVED stage
// reports — ideation, implementation, validation, then a LATER implementation
// (cycle 2) so the positional-last section is NOT the gated validation stage.
func interleavedFixturePath(t *testing.T) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join("testdata", "section-reader", "interleaved-fixture.md"))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// interleavedLines returns the fixture's splitLines view, so an emitted 1-based
// range slices these directly: lines[start-1:end].
func interleavedLines(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile(interleavedFixturePath(t))
	if err != nil {
		t.Fatal(err)
	}
	return splitLines(string(data))
}

// checklistEnvelope is the --checklist --json shape: the selected stage, the
// section's heading line, and each item's status/text/1-based range.
type checklistEnvelope struct {
	Command string `json:"command"`
	Stage   string `json:"stage"`
	Items   []struct {
		Status string `json:"status"`
		Text   string `json:"text"`
		Start  string `json:"start"`
		End    string `json:"end"`
	} `json:"checklist"`
}

// TestChecklistSelectsGatedStage (AC-1) asserts --stage validation --checklist
// selects the validation report (NOT the positional-last implementation cycle 2)
// and emits exactly its DONE/SKIPPED/FAILED items with 1-based ranges that slice
// the fixture to the item's known bullet+evidence text. The fixture's cat -n line
// numbers are the external oracle.
func TestChecklistSelectsGatedStage(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	lines := interleavedLines(t)
	fixture := interleavedFixturePath(t)

	out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixture, "--stage", "validation", "--checklist", "--json")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	var doc checklistEnvelope
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	if doc.Stage != "validation" {
		t.Fatalf("stage = %q, want validation", doc.Stage)
	}
	// The validation section's known items (oracle: cat -n above).
	type wantItem struct {
		status     string
		start, end int
		// bullet+evidence text, recomputed from the fixture lines for byte equality.
	}
	want := []wantItem{
		{"DONE", 37, 38},
		{"SKIPPED", 39, 40},
		{"FAILED", 41, 42},
	}
	if len(doc.Items) != len(want) {
		t.Fatalf("item count = %d, want %d\n%s", len(doc.Items), len(want), out)
	}
	for i, w := range want {
		got := doc.Items[i]
		if got.Status != w.status {
			t.Errorf("item[%d].status = %q, want %q", i, got.Status, w.status)
		}
		gs, ge := atoiT(t, got.Start), atoiT(t, got.End)
		if gs != w.start || ge != w.end {
			t.Fatalf("item[%d] range = [%d,%d], want [%d,%d]", i, gs, ge, w.start, w.end)
		}
		// Slice the fixture by the emitted range and assert it equals the section's
		// known bytes (the bullet line through its trailing evidence line).
		gotSlice := strings.Join(lines[gs-1:ge], "\n")
		wantSlice := strings.Join(lines[w.start-1:w.end], "\n")
		if gotSlice != wantSlice {
			t.Fatalf("item[%d] slice mismatch\n--- got ---\n%s\n--- want ---\n%s", i, gotSlice, wantSlice)
		}
	}
	// The positional-last section is implementation (cycle 2); a positional-last
	// bug would emit its "rework after the gate" item. Assert it is absent.
	for _, it := range doc.Items {
		if strings.Contains(it.Text, "rework after the gate") {
			t.Fatalf("validation checklist leaked the positional-last implementation cycle-2 item: %+v", it)
		}
	}
}

// TestChecklistSelectsLatestCycle (AC-1) asserts --stage implementation selects
// the LATEST implementation cycle (cycle 2) across the interleaved sections, not
// the earlier cycle-1 section.
func TestChecklistSelectsLatestCycle(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	fixture := interleavedFixturePath(t)

	out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixture, "--stage", "implementation", "--checklist", "--json")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	var doc checklistEnvelope
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	// Cycle-2 items present; cycle-1's "write the first cut" absent.
	var sawCycle2, sawCycle1 bool
	for _, it := range doc.Items {
		if strings.Contains(it.Text, "rework after the gate") {
			sawCycle2 = true
		}
		if strings.Contains(it.Text, "write the first cut") {
			sawCycle1 = true
		}
	}
	if !sawCycle2 {
		t.Fatalf("latest implementation cycle (cycle 2) item absent: %+v", doc.Items)
	}
	if sawCycle1 {
		t.Fatalf("earlier implementation cycle-1 item leaked into the latest-cycle selection: %+v", doc.Items)
	}
}

// acScanEnvelope is the --ac-scan --json shape: the stage and per-AC citation map.
type acScanEnvelope struct {
	Command string `json:"command"`
	Stage   string `json:"stage"`
	ACs     []struct {
		ID          string `json:"id"`
		Line        string `json:"line"`
		Unevidenced string `json:"unevidenced"` // all-strings contract: "true"/"false"
		Citations   []struct {
			Line string `json:"line"`
			Text string `json:"text"`
		} `json:"citations"`
	} `json:"acs"`
}

// TestACScanScopedToChecklist (AC-2) asserts --ac-scan cites AC evidence from the
// gated stage's checklist line ranges ONLY: AC-1 (cited at the validation DONE
// evidence line 38) is unevidenced=false; AC-2 (mentioned only at the Summary
// line 46, outside any checklist item) is unevidenced=true with no citations.
// No satisfied / no natural_place key is emitted.
func TestACScanScopedToChecklist(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	fixture := interleavedFixturePath(t)

	out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixture, "--stage", "validation", "--ac-scan", "--json")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	var doc acScanEnvelope
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	if len(doc.ACs) != 2 {
		t.Fatalf("AC count = %d, want 2\n%s", len(doc.ACs), out)
	}
	byID := map[string]int{}
	for i, ac := range doc.ACs {
		byID[ac.ID] = i
	}
	ac1i, ok1 := byID["AC-1"]
	ac2i, ok2 := byID["AC-2"]
	if !ok1 || !ok2 {
		t.Fatalf("expected AC-1 and AC-2 ids, got %v", byID)
	}
	ac1 := doc.ACs[ac1i]
	if ac1.Unevidenced != "false" {
		t.Errorf("AC-1 unevidenced = %q, want \"false\" (cited in the validation DONE evidence)", ac1.Unevidenced)
	}
	if len(ac1.Citations) == 0 {
		t.Fatalf("AC-1 has no citations, want one at line 38")
	}
	if atoiT(t, ac1.Citations[0].Line) != 38 {
		t.Errorf("AC-1 citation line = %s, want 38 (the DONE evidence line)", ac1.Citations[0].Line)
	}
	ac2 := doc.ACs[ac2i]
	if ac2.Unevidenced != "true" {
		t.Errorf("AC-2 unevidenced = %q, want \"true\" (mentioned only in the Summary prose, outside checklist ranges)", ac2.Unevidenced)
	}
	if len(ac2.Citations) != 0 {
		t.Errorf("AC-2 has %d citations, want 0 (the Summary mention is NOT a checklist citation)", len(ac2.Citations))
	}
	// No satisfied / no natural_place key anywhere in the emitted JSON.
	if strings.Contains(out, "satisfied") {
		t.Errorf("emitted JSON contains a 'satisfied' key (the verdict is L3's, not extracted): %s", out)
	}
	if strings.Contains(out, "natural_place") {
		t.Errorf("emitted JSON contains a 'natural_place' key (cut, routes to L3): %s", out)
	}
}

// TestACScanScopeIsLoadBearing (AC-2 mutation test) proves the citation-scope
// boundary is non-tautological: widening the scan to the WHOLE stage report
// (counting the Summary prose at line 46) would flip AC-2 to evidenced. The
// production --ac-scan must NOT do this, so AC-2 stays unevidenced. We assert the
// boundary directly: AC-2's only on-disk mention (line 46) lies OUTSIDE every
// emitted validation checklist range, so a correctly-scoped scan cannot cite it.
func TestACScanScopeIsLoadBearing(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	lines := interleavedLines(t)
	fixture := interleavedFixturePath(t)

	// The checklist ranges the gated stage owns.
	clOut, _, clCode := runNative(t, root, env, "--workflow-dir", root, "--read", fixture, "--stage", "validation", "--checklist", "--json")
	if clCode != 0 {
		t.Fatalf("checklist exit=%d", clCode)
	}
	var cl checklistEnvelope
	json.Unmarshal([]byte(clOut), &cl)

	// Find AC-2's sole on-disk mention line (the Summary at 46) and prove it is
	// outside every checklist range — so the scope boundary is what keeps AC-2
	// unevidenced (a whole-report scan would include line 46 and flip it).
	ac2Line := 0
	for i, l := range lines {
		if strings.Contains(l, "AC-2 is not yet evidenced") {
			ac2Line = i + 1
		}
	}
	if ac2Line == 0 {
		t.Fatal("fixture invariant broken: AC-2's Summary mention not found")
	}
	inRange := false
	for _, it := range cl.Items {
		s, e := atoiT(t, it.Start), atoiT(t, it.End)
		if ac2Line >= s && ac2Line <= e {
			inRange = true
		}
	}
	if inRange {
		t.Fatalf("AC-2's Summary mention (line %d) is INSIDE a checklist range — the scope boundary would not be load-bearing", ac2Line)
	}

	// And the production scan, correctly scoped, reports AC-2 unevidenced.
	out, _, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixture, "--stage", "validation", "--ac-scan", "--json")
	if code != 0 {
		t.Fatalf("ac-scan exit=%d", code)
	}
	var doc acScanEnvelope
	json.Unmarshal([]byte(out), &doc)
	for _, ac := range doc.ACs {
		if ac.ID == "AC-2" && ac.Unevidenced != "true" {
			t.Fatalf("AC-2 evidenced — the scan counted an out-of-checklist mention (scope not enforced)")
		}
	}
}

// annotatedACFixturePath returns the committed fixture whose ## Acceptance
// criteria section annotates the value AC inside the bold (**AC-1 (VALUE)**),
// alongside a bare **AC-2** and a prose-only "see AC-3 above" mention — the
// over-match discriminator for the broadened heading matcher.
func annotatedACFixturePath(t *testing.T) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join("testdata", "section-reader", "annotated-ac-fixture.md"))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// TestACScanEnumeratesAnnotatedAC (AC-1 + AC-2) asserts --ac-scan enumerates an
// AC whose bold span carries an annotation: **AC-1 (VALUE)** lists with id AC-1,
// exactly as the bare **AC-2** does, so the value AC the README ideation policy
// recommends is no longer silently dropped from the gate cross-check. The
// enumerated id set is EXACTLY {AC-1, AC-2}: AC-1 once, AC-2 once (no
// over-match), AC-3 absent (the prose "see AC-3 above" is not a heading, proving
// the **…** delimiter requirement survived the broadening). RED on the bare
// matcher (AC-1 absent), GREEN after acHeadingRe gains the trailing-label form.
func TestACScanEnumeratesAnnotatedAC(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	fixture := annotatedACFixturePath(t)

	out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixture, "--stage", "ideation", "--ac-scan", "--json")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	var doc acScanEnvelope
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	// Count each enumerated id's occurrences: the set must be exactly {AC-1, AC-2}.
	counts := map[string]int{}
	for _, ac := range doc.ACs {
		counts[ac.ID]++
	}
	if counts["AC-1"] != 1 {
		t.Errorf("AC-1 enumerated %d time(s), want exactly 1 (the annotated value AC must list as AC-1)", counts["AC-1"])
	}
	if counts["AC-2"] != 1 {
		t.Errorf("AC-2 enumerated %d time(s), want exactly 1 (bare AC, no over-match)", counts["AC-2"])
	}
	if counts["AC-3"] != 0 {
		t.Errorf("AC-3 enumerated %d time(s), want 0 (the prose 'see AC-3 above' is not a heading)", counts["AC-3"])
	}
	if len(doc.ACs) != 2 {
		t.Fatalf("enumerated %d ACs, want exactly 2 ({AC-1, AC-2}): %v\n%s", len(doc.ACs), counts, out)
	}
	// AC-1's evidence flag is exercised: its DONE evidence line cites AC-1, so a
	// correctly-enumerated AC-1 reports unevidenced=false (not merely "present").
	for _, ac := range doc.ACs {
		if ac.ID == "AC-1" && ac.Unevidenced != "false" {
			t.Errorf("AC-1 unevidenced = %q, want \"false\" (cited in the ideation DONE evidence)", ac.Unevidenced)
		}
	}
}

// TestGateModeLoudFailures (AC-5) asserts the gate modes fail loudly with a
// non-zero exit and a named diagnostic, never a partial/silent emit: missing
// --stage, a --stage matching no report (no silent positional-last fallback),
// and --ac-scan over a file with no ## Acceptance criteria section.
func TestGateModeLoudFailures(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	fixture := interleavedFixturePath(t)
	plain := fixturePath(t) // section-reader/fixture.md has no ## Acceptance criteria

	cases := []struct {
		name     string
		args     []string
		wantWord string
	}{
		{"checklist without --stage", []string{"--read", fixture, "--checklist", "--json"}, "--stage"},
		{"ac-scan without --stage", []string{"--read", fixture, "--ac-scan", "--json"}, "--stage"},
		{"stage matches no report", []string{"--read", fixture, "--stage", "done", "--checklist", "--json"}, "done"},
		{"ac-scan no acceptance criteria", []string{"--read", plain, "--stage", "ideation", "--ac-scan", "--json"}, "acceptance criteria"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"--workflow-dir", root}, tc.args...)
			out, stderr, code := runNative(t, root, env, args...)
			if code == 0 {
				t.Fatalf("exit=0 (want non-zero) — silent/partial emit\nstdout=%q", out)
			}
			if !strings.Contains(strings.ToLower(stderr), strings.ToLower(tc.wantWord)) {
				t.Fatalf("stderr = %q, want a diagnostic naming %q", stderr, tc.wantWord)
			}
			// A loud failure emits nothing on stdout (no partial JSON).
			if strings.TrimSpace(out) != "" {
				t.Fatalf("stdout = %q, want empty on a loud failure", out)
			}
		})
	}
}

// atoiT parses a string field to int, failing the test on a non-int (the
// all-strings contract means every numeric leaf is a string).
func atoiT(t *testing.T, s string) int {
	t.Helper()
	n, err := strconv.Atoi(s)
	if err != nil {
		t.Fatalf("non-int string %q: %v", s, err)
	}
	return n
}

// TestBootTaxonomySourceSufficient (AC-4 source-sufficiency) proves the rewritten
// FO Startup step 4 has a real source: a boot-shaped status --read docs/dev/README.md
// --json yields, in its stages array, every per-stage flag step 4 enumerates
// (initial/terminal/gate/worktree/feedback-to) for the five real stages — so the
// contract sentence that points step 4 at status --read is sufficient, proven by
// EXERCISING the real output, not a prose-grep over the contract.
//
// This AC proves only source-sufficiency. The behavioral adoption — a booting FO
// actually CALLING status --read and not the Read tool — is observed by
// haiku-drive-validation's live drive, not this member.
func TestBootTaxonomySourceSufficient(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	readme := devReadmePath(t)

	out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", readme, "--json")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	var doc readEnvelope
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	byName := map[string]map[string]string{}
	for _, s := range doc.Stages {
		byName[s["name"]] = s
	}
	// The taxonomy flags step 4 enumerates, per real stage that sets them. Every
	// stage object MUST carry the four typed flags (initial/terminal/gate/worktree)
	// as string leaves; feedback-to is present only where the README sets it.
	for _, name := range []string{"backlog", "ideation", "implementation", "validation", "done"} {
		s, ok := byName[name]
		if !ok {
			t.Fatalf("stage %q absent from the boot taxonomy source: %v", name, byName)
		}
		for _, flag := range []string{"initial", "terminal", "gate", "worktree"} {
			v, present := s[flag]
			if !present {
				t.Errorf("stage %q missing flag %q step 4 consumes: %v", name, flag, s)
				continue
			}
			if v != "true" && v != "false" {
				t.Errorf("stage %q flag %q = %q, want a string bool", name, flag, v)
			}
		}
	}
	// feedback-to is the one routing flag the gate needs; validation carries it.
	if byName["validation"]["feedback-to"] != "implementation" {
		t.Errorf("validation feedback-to = %q, want implementation (the gate's routing source)", byName["validation"]["feedback-to"])
	}
}
