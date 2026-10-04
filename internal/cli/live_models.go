// ABOUTME: `spacedock live-models` prints the live E2E lane model ids from the
// ABOUTME: single source in internal/release, so the workflow carries no literal.
package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/spacedock-dev/spacedock/internal/release"
)

// newLiveModelsCommand is the workflow-facing model printer: it emits one
// `key=id` line per lane/auth entry, or a single id for `--get key`. The live
// workflow resolves each lane's model through this command instead of repeating
// a literal, mirroring how the release pipeline derives its targets.
func newLiveModelsCommand(stdout, stderr io.Writer) *cobra.Command {
	var get string
	cmd := &cobra.Command{
		Use:    "live-models",
		Short:  "Print the live E2E lane model ids",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if get != "" {
				id, ok := release.LiveModel(get)
				if !ok {
					fmt.Fprintf(stderr, "spacedock live-models: unknown key %q\n", get)
					return exitCodeError{2}
				}
				fmt.Fprintln(stdout, id)
				return nil
			}
			for _, m := range release.LiveModels() {
				fmt.Fprintf(stdout, "%s=%s\n", m.Key, m.ID)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&get, "get", "", "print only the id for this lane key")
	return cmd
}
