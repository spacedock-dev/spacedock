package ensigncycle

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// The host-specific filing assertions for the `filing` scenario. They grade the
// FO's recorded tool-call STREAM — not the end-state file, which looks identical
// whether filed via `spacedock new` or hand-assembled, and not a grep of the
// contract prose. The producer signal is: the FO ran a `spacedock … new <slug>`
// invocation (the atomic-create path) and did NOT fall back to the manual
// `--next-id` + file-write pair. Claude runs commands via the `Bash` tool and
// writes files via the `Write` tool; Codex runs everything (including file
// writes) as `command_execution` items — so the manual-pair shape differs per
// host and the assertions live behind host adapters, like reviewer-reuse. They
// sit under the DEFAULT build tags (stdlib JSON only) so the offline negative
// tests exercise them without spending a model.

// newInvocation matches a spacedock atomic-create invocation in a command string:
// either the `new` subcommand or the `--new` flag (its alias), in a `spacedock`
// or `${SPACEDOCK_BIN…}` launcher call. The slug is matched separately so the
// command must carry BOTH the create verb and the requested slug. Single-line by
// design (`[^\n]*?` never crosses a newline) so it does not pair a launcher token
// on one line with a `new` verb on an unrelated later line.
var newInvocation = regexp.MustCompile(`(?:spacedock|SPACEDOCK_BIN)[^\n]*?(?:\bnew\b|--new)`)

// launcherCapture matches the contract-blessed var-capture of the resolved
// launcher — `B=${SPACEDOCK_BIN:-spacedock}` or
// `B="${SPACEDOCK_BIN:-spacedock}"` — anywhere in a command string. The captured
// var name is recorded so the create call below can require THAT var.
var launcherCapture = regexp.MustCompile(`([A-Za-z_][A-Za-z0-9_]*)=(?:\$\{SPACEDOCK_BIN:-spacedock\}|"\$\{SPACEDOCK_BIN:-spacedock\}")`)

// nextIDInvocation matches a `status --next-id` candidate-preview command — the
// first half of the manual filing pair the atomic path replaces.
var nextIDInvocation = regexp.MustCompile(`--next-id\b`)

// commandFilesViaNew reports whether a command string is the atomic-create call
// for the requested slug: a `new`/`--new` invocation that names the slug. It
// accepts two launcher shapes: a direct `spacedock … new` / `${SPACEDOCK_BIN…} new`
// call, and the var-capture idiom `B=${SPACEDOCK_BIN:-spacedock}; $B new` where the
// create call invokes the captured var (so the `$B new` segment carries no literal
// launcher token). The slug must always appear.
func commandFilesViaNew(command, slug string) bool {
	if !strings.Contains(command, slug) {
		return false
	}
	if newInvocation.MatchString(command) {
		return true
	}
	return capturedLauncherFilesViaNew(command)
}

// capturedLauncherFilesViaNew reports whether the command captured the resolved
// launcher into a var (`B=${SPACEDOCK_BIN:-spacedock}`) and then ran `$B new` /
// `$B --new` with that exact var. Tying recognition to the captured var name keeps
// it from matching an unrelated `$X new` that never resolved a spacedock launcher.
func capturedLauncherFilesViaNew(command string) bool {
	m := launcherCapture.FindStringSubmatch(command)
	if m == nil {
		return false
	}
	varName := regexp.QuoteMeta(m[1])
	// Require the captured name to end at a shell-variable boundary. Without it,
	// a capture of BIN could falsely recognize `$BINN new` as `$BIN new`.
	call := regexp.MustCompile(`"?\$(?:\{` + varName + `\}|` + varName + `)(?:[^\nA-Za-z0-9_]|$)[^\n]*?(?:\bnew\b|--new)`)
	return call.MatchString(command)
}

// assertClaudeFilingViaNew scans the stream-json transcript for the FO filing the
// seed via `spacedock … new <slug>` (a Bash tool call) and NOT via the manual
// `--next-id` + `Write` pair. It enforces both halves, because either alone
// false-passes the manual flow:
//
//  1. The FO ran a `spacedock … new <slug>` Bash command (the atomic-create path).
//  2. The FO did NOT emit the manual pair: a `--next-id` Bash command AND a
//     `Write` tool_use creating the entity `.md`. A `Write` of the entity file
//     after a `--next-id` preview is exactly the drift-prone flow `new` replaces,
//     so its presence FAILS even if `new` was also run.
func assertClaudeFilingViaNew(stream, slug string) error {
	filedViaNew := false
	sawNextID := false
	wroteEntityFile := false

	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var entry struct {
			Message *struct {
				Content []struct {
					Type  string `json:"type"`
					Name  string `json:"name"`
					Input struct {
						Command  string `json:"command"`
						FilePath string `json:"file_path"`
					} `json:"input"`
				} `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &entry); err != nil || entry.Message == nil {
			continue
		}
		for _, block := range entry.Message.Content {
			if block.Type != "tool_use" {
				continue
			}
			switch block.Name {
			case "Bash":
				if commandFilesViaNew(block.Input.Command, slug) {
					filedViaNew = true
				}
				if nextIDInvocation.MatchString(block.Input.Command) {
					sawNextID = true
				}
			case "Write":
				if strings.Contains(block.Input.FilePath, slug) && strings.HasSuffix(block.Input.FilePath, ".md") {
					wroteEntityFile = true
				}
			}
		}
	}

	if !filedViaNew {
		return fmt.Errorf("the FO did not file the seed via a `spacedock … new %s` command — it never used the atomic-create path", slug)
	}
	if sawNextID && wroteEntityFile {
		return fmt.Errorf("the FO emitted the manual `--next-id` + `Write %s.md` pair — it hand-assembled the entity instead of letting `new` write it atomically", slug)
	}
	return nil
}

// codexCommandItem is one `codex exec --json` command_execution item: Codex runs
// every shell action — including writing a file via heredoc/apply_patch — as a
// command_execution, so both the `new` invocation and any manual file-write land
// here.
type codexCommandItem struct {
	Type string `json:"type"`
	Item struct {
		Type    string `json:"type"`
		Command string `json:"command"`
	} `json:"item"`
}

// assertCodexFilingViaNew scans the `codex exec --json` transcript for the FO
// filing the seed via `spacedock … new <slug>` and NOT via the manual flow. On
// Codex there is no `Write` tool — the manual pair would be a `--next-id` command
// followed by a shell file-write — so the discriminator is the `--next-id`
// candidate-preview command: the atomic path needs none. It enforces both halves:
//
//  1. The FO ran a `spacedock … new <slug>` command_execution.
//  2. The FO did NOT run a `--next-id` filing command (the manual pair's id
//     source). `new` mints the id itself, so a `--next-id` here means the FO
//     reached for the manual flow.
func assertCodexFilingViaNew(jsonl, slug string) error {
	filedViaNew := false
	sawNextID := false

	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var ev codexCommandItem
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}
		if ev.Item.Type != "command_execution" {
			continue
		}
		if commandFilesViaNew(ev.Item.Command, slug) {
			filedViaNew = true
		}
		if nextIDInvocation.MatchString(ev.Item.Command) {
			sawNextID = true
		}
	}

	if !filedViaNew {
		return fmt.Errorf("the FO did not file the seed via a `spacedock … new %s` command — it never used the atomic-create path", slug)
	}
	if sawNextID {
		return fmt.Errorf("the FO ran a `--next-id` filing command — it reached for the manual preview-then-write flow instead of the atomic `new` path")
	}
	return nil
}
