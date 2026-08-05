package ensigncycle

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type loopEvent struct {
	action string
	value  string
}

func TestFirstOfficerEventLoopCommandLog(t *testing.T) {
	logs := readLoopCommandLog(t)

	mixed := logs["mixed"]
	assertBefore(t, mixed, "mod-hold", "next")
	for _, action := range []string{"gate-present", "gate-advance", "gate-merge"} {
		assertBefore(t, mixed, action, "next")
	}
	assertValues(t, mixed, "dispatch", "independent-a", "independent-b")
	assertAbsent(t, mixed, "idle", "stop")

	retry := logs["retry"]
	assertActions(t, retry, "next", "idle", "reconcile", "next", "dispatch")
	assertValues(t, retry, "idle", "-")
	assertValues(t, retry, "reconcile", "active-1")
	assertValues(t, retry, "next", "-", "released")
	assertValues(t, retry, "dispatch", "released")
	assertAbsent(t, retry, "stop")

	empty := logs["empty"]
	assertActions(t, empty, "next", "idle", "reconcile", "next", "stop")
	assertValues(t, empty, "idle", "-")
	assertValues(t, empty, "reconcile", "-")
	assertValues(t, empty, "next", "-", "-")
	assertValues(t, empty, "stop", "no-dispatchable")

	if err := validateNoFalseStop(logs["false-stop"]); err == nil {
		t.Fatal("legacy status --next=[] -> idle trace was accepted despite pending mod/gate work")
	}
}

func TestCodexWaitPredicate(t *testing.T) {
	tests := []struct {
		name      string
		worker    string
		otherWork bool
		wantWait  bool
	}{
		{"active unresolved and empty", "active", false, true},
		{"active unresolved with dispatch", "active", true, false},
		{"active unresolved with gate", "active", true, false},
		{"active unresolved with mod action", "active", true, false},
		{"active unresolved with state work", "active", true, false},
		{"completed", "completed", false, false},
		{"errored", "errored", false, false},
		{"absent", "absent", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldWaitForWorker(tt.worker, tt.otherWork); got != tt.wantWait {
				t.Fatalf("shouldWaitForWorker(%q, %v) = %v, want %v", tt.worker, tt.otherWork, got, tt.wantWait)
			}
		})
	}
}

func shouldWaitForWorker(worker string, otherWork bool) bool {
	return worker == "active" && !otherWork
}

func validateNoFalseStop(events []loopEvent) error {
	pending := false
	for _, event := range events {
		switch event.action {
		case "mod-hold", "gate-present", "gate-advance", "gate-merge":
			pending = true
		case "dispatch":
			pending = false
		case "stop":
			if pending || event.value != "no-dispatchable" {
				return os.ErrInvalid
			}
		}
	}
	return nil
}

func readLoopCommandLog(t *testing.T) map[string][]loopEvent {
	t.Helper()
	f, err := os.Open(filepath.Join("testdata", "fo_event_loop_command.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	logs := map[string][]loopEvent{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\t")
		if len(parts) != 3 {
			t.Fatalf("invalid command-log row %q", scanner.Text())
		}
		logs[parts[0]] = append(logs[parts[0]], loopEvent{action: parts[1], value: parts[2]})
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	return logs
}

func assertBefore(t *testing.T, events []loopEvent, before, after string) {
	t.Helper()
	indices := map[string]int{before: -1, after: -1}
	for i, event := range events {
		if _, ok := indices[event.action]; ok && indices[event.action] < 0 {
			indices[event.action] = i
		}
	}
	if indices[before] < 0 || indices[after] < 0 || indices[before] >= indices[after] {
		t.Fatalf("%q must precede %q: %#v", before, after, events)
	}
}

func assertActions(t *testing.T, events []loopEvent, want ...string) {
	t.Helper()
	got := make([]string, len(events))
	for i, event := range events {
		got[i] = event.action
	}
	if strings.Join(got, "\x00") != strings.Join(want, "\x00") {
		t.Fatalf("actions = %v, want %v", got, want)
	}
}

func assertValues(t *testing.T, events []loopEvent, action string, want ...string) {
	t.Helper()
	var got []string
	for _, event := range events {
		if event.action == action {
			got = append(got, event.value)
		}
	}
	if strings.Join(got, "\x00") != strings.Join(want, "\x00") {
		t.Fatalf("%s values = %v, want %v", action, got, want)
	}
}

func assertAbsent(t *testing.T, events []loopEvent, actions ...string) {
	t.Helper()
	for _, event := range events {
		for _, action := range actions {
			if event.action == action {
				t.Fatalf("unexpected %s event: %#v", action, events)
			}
		}
	}
}
