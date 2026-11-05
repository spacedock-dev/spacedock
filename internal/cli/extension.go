// ABOUTME: git-style external subcommands — `spacedock <name>` runs `spacedock-<name>`
// ABOUTME: from PATH when no built-in command claims <name>.
package cli

import (
	"strings"

	"github.com/spf13/cobra"
)

const extensionPrefix = "spacedock-"

// extensionArgv returns the argv for the PATH executable an unclaimed command
// resolves to, or nil to leave args to cobra. Built-ins always win, and the
// lookup runs before cobra parses so the extension receives its argv verbatim:
// the root's --version flag and unknown-flag whitelist would otherwise rewrite it.
func extensionArgv(root *cobra.Command, args []string, lookPath func(string) (string, error)) []string {
	if len(args) == 0 || strings.HasPrefix(args[0], "-") || strings.ContainsRune(args[0], '/') {
		return nil
	}
	root.InitDefaultHelpCmd()
	if cmd, _, err := root.Find(args[:1]); err != nil || cmd != root {
		return nil
	}
	bin, err := lookPath(extensionPrefix + args[0])
	if err != nil {
		return nil
	}
	return append([]string{bin}, args[1:]...)
}
