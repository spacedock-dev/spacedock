// ABOUTME: resolveTrunk — the single integration-trunk resolver, and the thin
// ABOUTME: `dispatch trunk` command that prints its bare value to stdout.
package dispatch

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// resolveTrunk is the ONE integration-trunk resolver. It reads the top-level
// `trunk:` key from the workflow README frontmatter (a sibling of `state:`) and
// returns its value, falling back to "main" — the post-flip trunk — when the key
// is empty, absent, or the README is unreadable. Both consumers (reconcile's
// classD/classE and the `dispatch trunk` command) call this so they cannot
// diverge. The marketplace channel stamp (`cli.devBranch`) is a separate axis and
// is never read here.
func resolveTrunk(workflowDir string) string {
	fm := status.ParseFrontmatter(filepath.Join(workflowDir, "README.md"))
	if t := fm["trunk"]; t != "" {
		return t
	}
	return "main"
}

// runDispatchTrunk prints the resolved integration trunk as a bare branch name on
// stdout (single trailing newline, nothing else) so prose consumers can capture
// it with `BASE=$(spacedock dispatch trunk --workflow-dir DIR)`. Any diagnostic
// goes to stderr; a stray stdout line would poison the `$(...)` capture. Exit 0
// always once the workflow dir resolves; exit 1 when the dir is not found.
func runDispatchTrunk(workflowDir string, stdout, stderr io.Writer) int {
	if info, err := os.Stat(workflowDir); err != nil || !info.IsDir() {
		fmt.Fprintf(stderr, "error: workflow directory not found: %s\n", workflowDir)
		return 1
	}
	fmt.Fprintln(stdout, resolveTrunk(workflowDir))
	return 0
}
