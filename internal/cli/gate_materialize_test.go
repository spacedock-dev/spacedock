package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// TestGateMaterializeAcceptsOnlyTheRoomOperand is AC-1's caller-coordinate
// falsifier at the CLI boundary. Every coordinate the public provider grammar
// forbids must be refused as a usage error, and the refusal must happen before
// the room is touched at all.
func TestGateMaterializeAcceptsOnlyTheRoomOperand(t *testing.T) {
	forbidden := [][]string{
		{"task", "--room", "ROOM"},
		{"--room", "ROOM", "--workflow-dir", "WORKFLOW"},
		{"--room", "ROOM", "--briefing", "BRIEFING"},
		{"--room", "ROOM", "--actor", "person:captain"},
		{"--room", "ROOM", "--approver", "person:captain"},
		{"--room", "ROOM", "--destination", "DEST"},
		{"--room", "ROOM", "--provider-package", "PKG"},
		{"--room", "ROOM", "--resolved-sources", "MANIFEST"},
		{"--room", "ROOM", "--terminal", "TERM"},
		{"--room", "ROOM", "--entity", "task"},
	}
	for _, args := range forbidden {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			root := materializeCLIRoom(t)
			before := roomSnapshot(t, root)

			var out, errOut bytes.Buffer
			full := append([]string{"gate", "materialize"}, args...)
			code := run(context.Background(), full, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)

			if code != 2 {
				t.Fatalf("exit=%d, want usage error 2 (stdout=%q stderr=%q)", code, out.String(), errOut.String())
			}
			if out.Len() != 0 {
				t.Fatalf("a refused invocation printed a launch tuple: %q", out.String())
			}
			if got := roomSnapshot(t, root); got != before {
				t.Fatal("a refused invocation changed room bytes")
			}
		})
	}
}

// TestGateMaterializeRequiresExactlyOneRoom keeps the sole operand singular, so
// a second room cannot silently select a different gate.
func TestGateMaterializeRequiresExactlyOneRoom(t *testing.T) {
	for _, args := range [][]string{
		{},
		{"--room", "ROOM", "--room", "OTHER"},
	} {
		root := materializeCLIRoom(t)
		var out, errOut bytes.Buffer
		full := append([]string{"gate", "materialize"}, args...)
		code := run(context.Background(), full, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
		if code != 2 || out.Len() != 0 {
			t.Fatalf("args %v exit=%d stdout=%q stderr=%q, want usage error 2", args, code, out.String(), errOut.String())
		}
		if !strings.Contains(errOut.String(), "exactly once") && !strings.Contains(errOut.String(), "requires an argument") {
			t.Fatalf("args %v stderr=%q, want a named singular-room refusal", args, errOut.String())
		}
	}
}

// TestGateMaterializeRefusesARoomOutsideACommissionedWorkflow proves the
// workflow root is derived from the room rather than assumed, so a room with no
// enclosing workflow fails instead of resolving against the process directory.
func TestGateMaterializeRefusesARoomOutsideACommissionedWorkflow(t *testing.T) {
	orphan := t.TempDir()
	room := filepath.Join(orphan, "task", "review", "validation", "briefing-1")
	if err := os.MkdirAll(filepath.Join(room, "provider"), 0o700); err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "materialize", "--room", room}, nil, orphan, nil, &out, &errOut, &status.NativeRunner{}, nil)

	if code != 1 || out.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q, want exit 1 with no launch tuple", code, out.String())
	}
	if !strings.Contains(errOut.String(), "commissioned workflow") {
		t.Fatalf("stderr=%q, want it to name the missing commissioned workflow", errOut.String())
	}
}

// materializeCLIRoom builds a commissioned workflow with a prepared-shaped room.
// It carries no valid gate authority: these tests assert the CLI flag closure,
// which must refuse before any authority read.
func materializeCLIRoom(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\ncommissioned-by: spacedock@0.27.0\nid-style: slug\nstages:\n  states:\n    - name: validation\n      initial: true\n      gate: true\n    - name: done\n      terminal: true\n---\n# Workflow\n")
	writeFile(t, filepath.Join(root, "task.md"), "---\nstatus: validation\ntitle: Task\n---\n# Task\n")
	room := filepath.Join(root, "task", "review", "validation", "briefing-1")
	if err := os.MkdirAll(filepath.Join(room, "provider"), 0o700); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(room, "request.json"), "{}\n")
	return root
}

func roomSnapshot(t *testing.T, root string) string {
	t.Helper()
	var sb strings.Builder
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		sb.WriteString(rel + "\x00" + info.Mode().String() + "\x00")
		if info.Mode().IsRegular() {
			body, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			sb.Write(body)
		}
		sb.WriteString("\x00")
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return sb.String()
}
