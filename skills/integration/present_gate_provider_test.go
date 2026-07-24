// ABOUTME: Fixture-backed smoke for the present-gate provider-channel contract.
// ABOUTME: It probes before launch, blocks through retention, and never chat-falls back after launch.
package integration

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

type fakeGatePresenter struct {
	t         *testing.T
	mode      string
	events    *[]string
	fixture   string
	failProbe bool
}

func (p *fakeGatePresenter) probe() error {
	*p.events = append(*p.events, "probe")
	if p.failProbe {
		return errors.New("presenter unavailable")
	}
	return nil
}

func (p *fakeGatePresenter) launch(room string) error {
	*p.events = append(*p.events, "/subspace:r gate "+room)
	providerDir := filepath.Join(room, "provider")
	if err := os.MkdirAll(providerDir, 0o755); err != nil {
		return err
	}
	copyProviderFixture(p.t, filepath.Join(providerDir, "result.json"), "provider", "result.json")
	copyProviderFixture(p.t, filepath.Join(providerDir, "presented-inventory.json"), "provider", "presented-inventory.json")
	if p.mode == "advisory" {
		resultPath := filepath.Join(providerDir, "result.json")
		var result map[string]any
		body, err := os.ReadFile(resultPath)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return err
		}
		result["status"] = "advisory"
		body, err = json.Marshal(result)
		if err != nil {
			return err
		}
		if err := os.WriteFile(resultPath, body, 0o644); err != nil {
			return err
		}
	}
	*p.events = append(*p.events, "presenter-finished")
	return nil
}

func runProviderChannel(t *testing.T, root, room string, presenter *fakeGatePresenter) (fallback bool, output string, code int) {
	t.Helper()
	if err := presenter.probe(); err != nil {
		return true, err.Error(), 0
	}
	if err := presenter.launch(room); err != nil {
		return false, err.Error(), 1
	}
	*presenter.events = append(*presenter.events, "spacedock gate record task --room "+room)
	cmd := exec.Command(spacedockBinary(t), "gate", "record", "task", "--workflow-dir", root, "--room", room)
	body, err := cmd.CombinedOutput()
	if err == nil {
		return false, string(body), 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return false, string(body), exitErr.ExitCode()
	}
	t.Fatalf("run provider room recorder: %v\n%s", err, body)
	return false, "", 1
}

func stageProviderGate(t *testing.T) (root, entity, room string) {
	t.Helper()
	root = t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("---\nid-style: slug\nstages:\n  states:\n    - name: validation\n      initial: true\n    - name: done\n      terminal: true\n---\n# Workflow\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	entity = filepath.Join(root, "task.md")
	if err := os.WriteFile(entity, []byte("---\nstatus: validation\ntitle: Task\n---\n# Task\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	room = filepath.Join(root, "review", "validation", "briefing-1")
	if err := os.MkdirAll(room, 0o755); err != nil {
		t.Fatal(err)
	}
	copyProviderFixture(t, filepath.Join(room, "briefing.json"), "briefing.json")
	copyProviderFixture(t, filepath.Join(room, "request.json"), "request.json")
	cmd := exec.Command(spacedockBinary(t), "gate", "record", "task", "--workflow-dir", root, "--briefing", filepath.Join(room, "briefing.json"))
	if body, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("bind provider room: %v\n%s", err, body)
	}
	return root, entity, room
}

func copyProviderFixture(t *testing.T, destination string, parts ...string) {
	t.Helper()
	sourceParts := append([]string{repoRoot(t), "internal", "gates", "testdata", "gate-room"}, parts...)
	body, err := os.ReadFile(filepath.Join(sourceParts...))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(destination, body, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestPresentGateProviderChannelSmoke(t *testing.T) {
	skillPath := filepath.Join(repoRoot(t), "skills", "present-gate", "SKILL.md")
	skill, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatal(err)
	}
	contract := string(skill)
	probeAt := strings.Index(contract, "**Probe before side effects.**")
	launchAt := strings.Index(contract, "/subspace:r gate <gate-room>")
	for _, required := range []string{
		"spacedock} gate record <entity> --room <gate-room>",
		"Do not fall back to chat after launch",
		"complete only after the presenter exits, the Result validates, and retention finishes",
	} {
		if !strings.Contains(contract, required) {
			t.Fatalf("present-gate provider contract missing %q", required)
		}
	}
	if probeAt < 0 || launchAt < 0 || probeAt >= launchAt {
		t.Fatal("present-gate provider contract no longer requires probing before launch")
	}

	t.Run("successful retained binding", func(t *testing.T) {
		root, entity, room := stageProviderGate(t)
		var events []string
		presenter := &fakeGatePresenter{t: t, events: &events}
		fallback, output, code := runProviderChannel(t, root, room, presenter)
		if fallback || code != 0 {
			t.Fatalf("successful provider fallback=%v code=%d output=%q", fallback, code, output)
		}
		wantEvents := []string{
			"probe",
			"/subspace:r gate " + room,
			"presenter-finished",
			"spacedock gate record task --room " + room,
		}
		if !slices.Equal(events, wantEvents) {
			t.Fatalf("provider events = %q, want %q", events, wantEvents)
		}
		body, err := os.ReadFile(entity)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), "decision: approve") || !strings.Contains(string(body), "by: person:captain") {
			t.Fatalf("provider Result did not close the bound gate:\n%s", body)
		}
		for _, name := range []string{"result.json", "presented-inventory.json"} {
			if info, err := os.Stat(filepath.Join(room, "provider", name)); err != nil || info.Size() == 0 {
				t.Fatalf("retained provider output %s: info=%v err=%v", name, info, err)
			}
		}
	})

	t.Run("advisory result stays retained and open", func(t *testing.T) {
		root, entity, room := stageProviderGate(t)
		before, err := os.ReadFile(entity)
		if err != nil {
			t.Fatal(err)
		}
		var events []string
		presenter := &fakeGatePresenter{t: t, mode: "advisory", events: &events}
		fallback, output, code := runProviderChannel(t, root, room, presenter)
		if fallback || code != 1 || !strings.Contains(output, "advisory Result remains evidence") {
			t.Fatalf("advisory provider fallback=%v code=%d output=%q", fallback, code, output)
		}
		after, err := os.ReadFile(entity)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(before, after) {
			t.Fatal("invalid provider Result mutated or chat-closed the gate")
		}
		if len(events) != 4 || events[2] != "presenter-finished" {
			t.Fatalf("advisory provider did not block through retention: %q", events)
		}
		if _, err := os.Stat(filepath.Join(room, "provider", "result.json")); err != nil {
			t.Fatalf("advisory Result was not retained: %v", err)
		}
	})

	t.Run("failed probe falls back without launch", func(t *testing.T) {
		root, entity, room := stageProviderGate(t)
		before, err := os.ReadFile(entity)
		if err != nil {
			t.Fatal(err)
		}
		var events []string
		presenter := &fakeGatePresenter{t: t, events: &events, failProbe: true}
		fallback, _, code := runProviderChannel(t, root, room, presenter)
		if !fallback || code != 0 || !slices.Equal(events, []string{"probe"}) {
			t.Fatalf("failed probe fallback=%v code=%d events=%q", fallback, code, events)
		}
		after, err := os.ReadFile(entity)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(before, after) {
			t.Fatal("failed probe mutated the gate")
		}
		if _, err := os.Stat(filepath.Join(room, "provider")); !os.IsNotExist(err) {
			t.Fatalf("failed probe launched or created provider outputs: %v", err)
		}
	})
}
