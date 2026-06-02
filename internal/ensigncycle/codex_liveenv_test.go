package ensigncycle

import "testing"

func TestDecideCodexLiveAuth(t *testing.T) {
	t.Run("api_key_runs", func(t *testing.T) {
		d := decideCodexLiveAuth("sk-test", "")
		if d.mode != codexAuthRun {
			t.Fatalf("mode = %d, want codexAuthRun", d.mode)
		}
		if d.message != "" {
			t.Fatalf("message = %q, want empty", d.message)
		}
	})

	t.Run("missing_local_auth_skips", func(t *testing.T) {
		d := decideCodexLiveAuth("", "")
		if d.mode != codexAuthSkip {
			t.Fatalf("mode = %d, want codexAuthSkip", d.mode)
		}
		if d.message == "" {
			t.Fatal("skip decision must carry a message")
		}
	})

	t.Run("missing_required_auth_fails", func(t *testing.T) {
		d := decideCodexLiveAuth("", "1")
		if d.mode != codexAuthFatal {
			t.Fatalf("mode = %d, want codexAuthFatal", d.mode)
		}
		if d.message == "" {
			t.Fatal("fatal decision must carry a clear message")
		}
	})
}
