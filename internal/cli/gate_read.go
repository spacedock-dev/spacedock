package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/status"
)

func runGateRead(args []string, dir string, stdout, stderr io.Writer) error {
	ref, workflowDir := "", ""
	references, jsonFlags, workflowFlags := 0, 0, 0
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			jsonFlags++
		case "--workflow-dir":
			workflowFlags++
			if workflowFlags > 1 {
				fmt.Fprintln(stderr, "Error: gate read accepts --workflow-dir exactly once")
				return exitCodeError{2}
			}
			if i+1 >= len(args) || args[i+1] == "" || strings.HasPrefix(args[i+1], "--") {
				fmt.Fprintln(stderr, "Error: --workflow-dir requires an argument")
				return exitCodeError{2}
			}
			i++
			workflowDir = args[i]
		default:
			if strings.HasPrefix(args[i], "-") {
				fmt.Fprintf(stderr, "Error: unknown gate read flag: %s\n", args[i])
				return exitCodeError{2}
			}
			references++
			ref = args[i]
		}
	}
	if references != 1 || ref == "" || jsonFlags != 1 {
		fmt.Fprintln(stderr, "Error: gate read requires exactly one reference and --json")
		return exitCodeError{2}
	}
	if code := status.ReadGateHistory(workflowDir, dir, ref, stdout, stderr); code != 0 {
		return exitCodeError{code}
	}
	return nil
}
