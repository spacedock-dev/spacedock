// ABOUTME: Proves Claude live detector evidence is silent on success and secondary on failure.
// ABOUTME: Covers diagnostic selection priority and failure-only cleanup reporting offline.
package ensigncycle

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

type fakeClaudeLiveDiagnosticReporter struct {
	failed   bool
	cleanups []func()
	events   []string
}

func (r *fakeClaudeLiveDiagnosticReporter) Cleanup(cleanup func()) {
	r.cleanups = append(r.cleanups, cleanup)
}

func (r *fakeClaudeLiveDiagnosticReporter) Failed() bool {
	return r.failed
}

func (r *fakeClaudeLiveDiagnosticReporter) Logf(format string, args ...any) {
	r.events = append(r.events, fmt.Sprintf(format, args...))
}

func (r *fakeClaudeLiveDiagnosticReporter) injectPrimary(message string) {
	r.failed = true
	r.events = append(r.events, message)
}

func (r *fakeClaudeLiveDiagnosticReporter) runCleanups() {
	for i := len(r.cleanups) - 1; i >= 0; i-- {
		r.cleanups[i]()
	}
}

func TestClaudeLiveFailureDiagnosticIsSecondaryOnly(t *testing.T) {
	const fixtureRoot = "/tmp/TestLiveClaudeSharedScenariosself-evidence-merge-triage1234567890/001"
	stream := strings.Join([]string{
		streamLine(`spacedock status --boot --workflow-dir ` + fixtureRoot),
		streamLine(`find / -maxdepth 6 -iname "fo-merge-core.md"`),
	}, "\n")
	diagnostic := detectClaudeLiveFailureDiagnostic(stream, fixtureRoot)
	if diagnostic == nil {
		t.Fatal("test setup did not produce broad-search diagnostic evidence")
	}

	tests := []struct {
		name         string
		primary      string
		wantFailed   bool
		wantTimeline []string
	}{
		{
			name:         "passing scenario stays silent",
			wantTimeline: nil,
		},
		{
			name:       "runner failure stays primary",
			primary:    "runner failed: exact bytes\nwith context",
			wantFailed: true,
			wantTimeline: []string{
				"runner failed: exact bytes\nwith context",
				"Additional diagnostic: " + diagnostic.Error(),
			},
		},
		{
			name:       "scenario assertion stays primary",
			primary:    "scenario assertion: exact bytes",
			wantFailed: true,
			wantTimeline: []string{
				"scenario assertion: exact bytes",
				"Additional diagnostic: " + diagnostic.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reporter := &fakeClaudeLiveDiagnosticReporter{}
			registerClaudeLiveFailureDiagnostic(reporter, diagnostic)
			if tt.primary != "" {
				reporter.injectPrimary(tt.primary)
			}
			reporter.runCleanups()

			if reporter.failed != tt.wantFailed {
				t.Fatalf("failed = %v, want %v", reporter.failed, tt.wantFailed)
			}
			if !reflect.DeepEqual(reporter.events, tt.wantTimeline) {
				t.Fatalf("timeline = %#v, want %#v", reporter.events, tt.wantTimeline)
			}
			if tt.primary != "" && reporter.events[0] != tt.primary {
				t.Fatalf("primary bytes changed: got %q, want %q", reporter.events[0], tt.primary)
			}
		})
	}
}

func TestRegisterClaudeLiveFailureDiagnostic(t *testing.T) {
	const diagnosticText = "full transcript contains a filesystem hunt"
	diagnostic := fmt.Errorf("%s", diagnosticText)

	tests := []struct {
		name         string
		diagnostic   error
		primary      string
		wantCleanups int
		wantTimeline []string
	}{
		{
			name: "no observation registers nothing",
		},
		{
			name:         "passing run emits nothing",
			diagnostic:   diagnostic,
			wantCleanups: 1,
		},
		{
			name:         "runner failure precedes diagnostic",
			diagnostic:   diagnostic,
			primary:      "runner primary",
			wantCleanups: 1,
			wantTimeline: []string{"runner primary", "Additional diagnostic: " + diagnosticText},
		},
		{
			name:         "scenario failure precedes diagnostic",
			diagnostic:   diagnostic,
			primary:      "scenario primary",
			wantCleanups: 1,
			wantTimeline: []string{"scenario primary", "Additional diagnostic: " + diagnosticText},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reporter := &fakeClaudeLiveDiagnosticReporter{}
			registerClaudeLiveFailureDiagnostic(reporter, tt.diagnostic)
			if got := len(reporter.cleanups); got != tt.wantCleanups {
				t.Fatalf("registered cleanups = %d, want %d", got, tt.wantCleanups)
			}
			if tt.primary != "" {
				reporter.injectPrimary(tt.primary)
			}
			reporter.runCleanups()
			if !reflect.DeepEqual(reporter.events, tt.wantTimeline) {
				t.Fatalf("timeline = %#v, want %#v", reporter.events, tt.wantTimeline)
			}
		})
	}
}

func TestDetectClaudeLiveFailureDiagnosticBroadSearch(t *testing.T) {
	const fixtureRoot = "/tmp/TestLiveClaudeSharedScenariosself-evidence-merge-triage1234567890/001"
	stream := strings.Join([]string{
		streamLine(`spacedock status --boot --workflow-dir ` + fixtureRoot),
		streamLine(`find / -maxdepth 6 -iname "fo-merge-core.md"`),
	}, "\n")

	diagnostic := detectClaudeLiveFailureDiagnostic(stream, fixtureRoot)
	if diagnostic == nil {
		t.Fatal("broad-search evidence was not selected")
	}
	if !strings.Contains(diagnostic.Error(), "FO broad-searched the filesystem at boot") {
		t.Fatalf("diagnostic must name broad-search evidence, got: %v", diagnostic)
	}
}

func TestDetectClaudeLiveFailureDiagnosticWrongRootTakesPriority(t *testing.T) {
	const fixtureRoot = "/tmp/TestLiveClaudeSharedScenariosfiling1234567890/001"
	const realRepo = "/home/runner/work/spacedock/spacedock"
	stream := strings.Join([]string{
		streamLine(`cd ` + realRepo + ` && spacedock status --discover`),
		streamLine(`find / -maxdepth 6 -iname "fo-merge-core.md"`),
	}, "\n")

	diagnostic := detectClaudeLiveFailureDiagnostic(stream, fixtureRoot)
	if diagnostic == nil {
		t.Fatal("wrong-root evidence was not selected")
	}
	if !strings.Contains(diagnostic.Error(), "FO booted the wrong root") {
		t.Fatalf("wrong-root evidence must take priority, got: %v", diagnostic)
	}
}

func TestDetectClaudeLiveFailureDiagnosticCleanStream(t *testing.T) {
	const fixtureRoot = "/tmp/TestLiveClaudeSharedScenariosgate-guardrail1234567890/001"
	stream := strings.Join([]string{
		streamLine(`spacedock --version`),
		streamLine(`spacedock status --boot --workflow-dir ` + fixtureRoot),
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Read","input":{"file_path":"` + fixtureRoot + `/README.md"}}]}}`,
	}, "\n")

	if diagnostic := detectClaudeLiveFailureDiagnostic(stream, fixtureRoot); diagnostic != nil {
		t.Fatalf("clean stream selected diagnostic evidence: %v", diagnostic)
	}

	reporter := &fakeClaudeLiveDiagnosticReporter{}
	registerClaudeLiveFailureDiagnostic(reporter, nil)
	if len(reporter.cleanups) != 0 {
		t.Fatalf("clean stream registered %d cleanups, want 0", len(reporter.cleanups))
	}
}

func TestClaudeSharedRunnerRoutesDetectorEvidenceThroughFailureOnlyReporter(t *testing.T) {
	sourcePath := "claude_live_runner_test.go"
	source, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatalf("read %s: %v", sourcePath, err)
	}
	body := string(source)
	want := "registerClaudeLiveFailureDiagnostic(t, detectClaudeLiveFailureDiagnostic(stream, workflowRoot))"
	if !strings.Contains(body, want) {
		t.Fatalf("live gate-stop does not route its captured stream through the failure-only reporter: missing %q", want)
	}
	for _, direct := range []string{"detectWrongRootBoot(stream, workflowRoot)", "detectBroadSearchAtBoot(stream, workflowRoot)"} {
		if strings.Contains(body, direct) {
			t.Errorf("live gate-stop directly invokes fatal detector path %q outside the failure-only reporter", direct)
		}
	}
}

func TestLiveGateStopCapturedDiagnosticStaysSecondary(t *testing.T) {
	const fixtureRoot = "/tmp/TestLiveDefaultHeadlessStopsAtGate1234567890/001"
	stream := streamLine(`cd /home/runner/work/spacedock/spacedock && spacedock status --discover`)
	reporter := &fakeClaudeLiveDiagnosticReporter{}
	registerClaudeLiveFailureDiagnostic(reporter, detectClaudeLiveFailureDiagnostic(stream, fixtureRoot))
	reporter.injectPrimary("durable gate assertion failed")
	reporter.runCleanups()
	if len(reporter.events) != 2 {
		t.Fatalf("timeline = %#v, want durable failure plus one detector diagnostic", reporter.events)
	}
	if reporter.events[0] != "durable gate assertion failed" {
		t.Fatalf("primary failure changed or moved: %#v", reporter.events)
	}
	if !strings.Contains(reporter.events[1], "Additional diagnostic: FO booted the wrong root") {
		t.Fatalf("detector evidence was not secondary: %#v", reporter.events)
	}
}
