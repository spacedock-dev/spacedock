// ABOUTME: Integration-private room-only route for gate source materialization.
// ABOUTME: Accepts one room operand and derives every other coordinate from it.
package cli

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/spacedock-dev/spacedock/internal/gates"
	"github.com/spacedock-dev/spacedock/internal/status"
)

// runGateMaterialize is the integration-private provider entry. The caller
// supplies exactly one room and no authority, source, destination, or terminal
// coordinate; fixed code derives the entity, workflow root, provider root, and
// manifest path from that room before any Git read or provider-visible write.
func runGateMaterialize(dir string, args []string, stdout, stderr io.Writer) error {
	room := ""
	roomCount := 0
	for i := 0; i < len(args); i++ {
		if args[i] != "--room" {
			fmt.Fprintf(stderr, "Error: gate materialize accepts only --room: %s\n", args[i])
			return exitCodeError{2}
		}
		if i+1 >= len(args) {
			fmt.Fprintln(stderr, "Error: --room requires an argument")
			return exitCodeError{2}
		}
		room = args[i+1]
		roomCount++
		i++
	}
	if roomCount != 1 {
		fmt.Fprintln(stderr, "Error: gate materialize requires --room exactly once")
		return exitCodeError{2}
	}
	if !filepath.IsAbs(room) {
		room = filepath.Join(dir, room)
	}
	// The workflow definition root is derived by walking up from the room, not
	// accepted from the caller, so the public provider grammar stays room-only.
	workflowDir, ok := status.DiscoverWorkflowDir(room)
	if !ok {
		fmt.Fprintln(stderr, "Error: gate room is not inside a commissioned workflow")
		return exitCodeError{1}
	}
	result, err := gates.Materialize(gates.MaterializeInput{Room: room, WorkflowDir: workflowDir})
	if err != nil {
		fmt.Fprintln(stderr, "Error:", err)
		return exitCodeError{1}
	}
	fmt.Fprintf(stdout, "manifest=%s\n", result.Manifest)
	fmt.Fprintf(stdout, "sources=%d\n", result.Sources)
	fmt.Fprintf(stdout, "briefing=%s\n", result.Briefing)
	fmt.Fprintf(stdout, "provider=%s\n", result.ProviderRoot)
	fmt.Fprintf(stdout, "actor=%s\n", result.Actor)
	fmt.Fprintf(stdout, "approver=%s\n", result.Approver)
	return nil
}
