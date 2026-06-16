// ABOUTME: Detects `status --read` adoption in a tool-call command string,
// ABOUTME: launcher-agnostic and flag-order independent, shared by both runtimes.
package journeymetrics

import (
	"encoding/json"
	"strings"
)

// statusReadLaunchers are the command prefixes under which `status` is the
// spacedock subcommand: the two installed binaries plus the
// ${SPACEDOCK_BIN:-spacedock} form the fetch commands emit.
var statusReadLaunchers = map[string]bool{
	"spacedock":                   true,
	"sd":                          true,
	"${SPACEDOCK_BIN:-spacedock}": true,
}

// commandInvokesStatusRead reports whether cmd invokes the spacedock `status`
// subcommand with a `--read` flag. The `status` token must be the subcommand —
// either the first token (a bare `status --read`) or immediately preceded by a
// recognized launcher — so a `status` substring quoted as an argument to a
// different command (e.g. `echo 'status --read'`) does not match. The `--read`
// flag may appear in any order among status's flags.
func commandInvokesStatusRead(cmd string) bool {
	tokens := strings.Fields(cmd)
	for i, tok := range tokens {
		if tok != "status" {
			continue
		}
		isSubcommand := i == 0 || statusReadLaunchers[tokens[i-1]]
		if !isSubcommand {
			continue
		}
		for _, rest := range tokens[i+1:] {
			if rest == "--read" {
				return true
			}
		}
	}
	return false
}

// readInputIsScoped reports whether a Read tool_use input is a scoped read — one
// carrying a non-zero offset or limit (the line-range fields). A whole-file Read
// carries only file_path, so neither field is present and the read is unscoped.
func readInputIsScoped(input json.RawMessage) bool {
	if len(input) == 0 {
		return false
	}
	var scope struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	}
	if err := json.Unmarshal(input, &scope); err != nil {
		return false
	}
	return scope.Offset != 0 || scope.Limit != 0
}

// bashCommand extracts the command string from a Bash tool_use input.
func bashCommand(input json.RawMessage) string {
	if len(input) == 0 {
		return ""
	}
	var args struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal(input, &args); err != nil {
		return ""
	}
	return args.Command
}

// codexCommand extracts the shell command string from a codex tool_call.started
// arguments object (file reads and CLI calls both go through the shell on codex).
func codexCommand(arguments json.RawMessage) string {
	if len(arguments) == 0 {
		return ""
	}
	var args struct {
		Cmd string `json:"cmd"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil {
		return ""
	}
	return args.Cmd
}
