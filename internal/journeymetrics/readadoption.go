// ABOUTME: Detects `status --read` adoption in a tool-call command string,
// ABOUTME: launcher-agnostic and flag-order independent, shared by both runtimes.
package journeymetrics

import (
	"encoding/json"
	"strings"
)

// statusReadLaunchers are the literal command names under which `status` is the
// spacedock subcommand: the two installed binaries, the `./spacedock` dev build,
// and the `spacedock_launcher` shell function the dispatch fetch commands define
// and call. The ${SPACEDOCK_BIN:-spacedock} variable form is recognized separately
// (launcherIsSpacedock) since it is a value reference, not a fixed name.
var statusReadLaunchers = map[string]bool{
	"spacedock":          true,
	"sd":                 true,
	"./spacedock":        true,
	"spacedock_launcher": true,
}

// launcherIsSpacedock reports whether the token immediately preceding `status`
// launches the spacedock binary. It normalizes the token first — stripping one
// layer of surrounding quotes, so the contract-canonical quoted
// "${SPACEDOCK_BIN:-spacedock}" form (shared-core invariant) is recognized — then
// accepts either a literal launcher name or any reference to the SPACEDOCK_BIN
// variable family (${SPACEDOCK_BIN:-spacedock}, quoted or not). $SPACEDOCK is a
// legacy variable, NOT the SPACEDOCK_BIN family, so it stays unrecognized.
func launcherIsSpacedock(tok string) bool {
	tok = stripSurroundingQuotes(tok)
	if statusReadLaunchers[tok] {
		return true
	}
	return strings.Contains(tok, "SPACEDOCK_BIN")
}

// stripSurroundingQuotes removes one matching pair of surrounding single or double
// quotes from s, leaving an unquoted or unbalanced token untouched.
func stripSurroundingQuotes(s string) string {
	if len(s) >= 2 {
		q := s[0]
		if (q == '"' || q == '\'') && s[len(s)-1] == q {
			return s[1 : len(s)-1]
		}
	}
	return s
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
		isSubcommand := i == 0 || launcherIsSpacedock(tokens[i-1])
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

// FoldEnsignReadAdoption sums each dispatched-ensign sub-agent transcript's
// `status --read` and scoped-`Read` counts onto obs's two --read adoption
// counters. `status --read` adoption is principally an ensign behavior, and the
// ensign runs as a separate sub-agent session whose transcript never reaches the
// FO front-door stream — so the metric must fold those transcripts to observe the
// surface the contract sites actually steer.
//
// Each transcript is parsed independently and its counts SUMMED (not concatenated
// then parsed): tool IDs are unique per session, so summing cannot cross-count and
// the per-stream multi-delta tool-ID dedup stays correct within each transcript.
// Concatenating would risk tool-ID collision across sessions. Only the two
// adoption counters fold; turns/tokens/tool_calls keep their FO-front-door
// semantics.
func FoldEnsignReadAdoption(obs Observation, transcripts [][]byte) (Observation, error) {
	for _, transcript := range transcripts {
		parsed, err := ParseClaudeJSONL(transcript)
		if err != nil {
			return obs, err
		}
		obs.StatusReadCalls += parsed.Observation.StatusReadCalls
		obs.ScopedReadCalls += parsed.Observation.ScopedReadCalls
	}
	return obs, nil
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
