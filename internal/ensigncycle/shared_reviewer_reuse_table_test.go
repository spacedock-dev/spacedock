package ensigncycle

import "testing"

// Offline table tests for the host-specific reviewer-reuse assertions. They prove
// each assertion requires a REAL reuse tool call targeting the validation reviewer
// — not loose narration, not an unrelated tool, not a call to a different
// recipient. Without these the assertions are exercised only by the model-spending
// live runners; here they cost milliseconds and pin the discriminating behavior.

func TestAssertClaudeReviewerReuse(t *testing.T) {
	// A real SendMessage tool_use targeting the validation reviewer — the only shape
	// that should pass.
	realReuse := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"spacedock-ensign-validation","message":"re-review the fix"}}]}}`

	cases := []struct {
		name    string
		stream  string
		wantErr bool
	}{
		{"real SendMessage to validation", realReuse, false},
		{
			"loose narration only",
			`{"type":"assistant","message":{"content":[{"type":"text","text":"I will reuse the validation reviewer via SendMessage."}]}}`,
			true,
		},
		{
			"unrelated tool targeting validation",
			`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","input":{"to":"spacedock-ensign-validation","description":"fresh validation dispatch"}}]}}`,
			true,
		},
		{
			"SendMessage to a non-validation recipient",
			`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"spacedock-ensign-implementation","message":"apply the fix"}}]}}`,
			true,
		},
		{"empty stream", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := assertClaudeReviewerReuse(tc.stream)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for %q, got nil", tc.name)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected pass for %q, got: %v", tc.name, err)
			}
		})
	}
}

func TestAssertCodexReviewerReuse(t *testing.T) {
	// A real send_input tool call whose arguments reference the validation worker —
	// the only shape that should pass.
	realReuse := `{"type":"tool_call.started","name":"send_input","arguments":{"task":"validation-task","input":"re-review the fix"}}`

	cases := []struct {
		name    string
		jsonl   string
		wantErr bool
	}{
		{"real send_input to validation", realReuse, false},
		{
			"loose narration only",
			`{"type":"message","role":"assistant","content":"I will send_input to the validation worker."}`,
			true,
		},
		{
			"unrelated tool referencing validation",
			`{"type":"tool_call.started","name":"exec_command","arguments":{"cmd":"echo validation"}}`,
			true,
		},
		{
			"send_input to a non-validation worker",
			`{"type":"tool_call.started","name":"send_input","arguments":{"task":"implementation-task","input":"apply the fix"}}`,
			true,
		},
		{"empty transcript", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := assertCodexReviewerReuse(tc.jsonl)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for %q, got nil", tc.name)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected pass for %q, got: %v", tc.name, err)
			}
		})
	}
}
