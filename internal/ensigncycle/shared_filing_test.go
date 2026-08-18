package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

// launcherCapture matches the three contract-blessed assignments of the
// resolved launcher: unquoted, balanced-double-quoted, or balanced-single-quoted.
// It also accepts the exact balanced command-v fallback observed in PR #496.
// Each alternative captures its var name separately; independent optional quote
// classes are forbidden because they accept mismatched or one-sided quotes.
var launcherCapture = regexp.MustCompile(`(?:([A-Za-z_][A-Za-z0-9_]*)=\$\{SPACEDOCK_BIN:-spacedock\}|([A-Za-z_][A-Za-z0-9_]*)="\$\{SPACEDOCK_BIN:-spacedock\}"|([A-Za-z_][A-Za-z0-9_]*)='\$\{SPACEDOCK_BIN:-spacedock\}'|([A-Za-z_][A-Za-z0-9_]*)="\$\{SPACEDOCK_BIN:-\$\(command -v spacedock\)\}")(?:[;\s]|$)`)

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
	if regexp.MustCompile(newInvocation.String() + `[ \t]+` + regexp.QuoteMeta(slug) + `(?:[ \t'\n]|$)`).MatchString(command) {
		return true
	}
	return capturedLauncherFilesViaNew(command, slug)
}

// capturedLauncherFilesViaNew reports whether the command captured the resolved
// launcher into a var (`B=${SPACEDOCK_BIN:-spacedock}`) and then ran `$B new` /
// `$B --new` with that exact var. Tying recognition to the captured var name keeps
// it from matching an unrelated `$X new` that never resolved a spacedock launcher.
func capturedLauncherFilesViaNew(command, slug string) bool {
	m := launcherCapture.FindStringSubmatch(command)
	if m == nil {
		return false
	}
	varName := ""
	for _, candidate := range m[1:] {
		if candidate != "" {
			varName = candidate
			break
		}
	}
	if varName == "" {
		return false
	}
	// This is deliberately a bounded simple-command recognizer, not a shell
	// parser. Shell command separators end the candidate; the captured variable
	// must then be the balanced executable at the start of one segment, with the
	// create verb and requested slug in that same segment.
	captureEnd := strings.Index(command, m[0]) + len(m[0])
	segments := regexp.MustCompile(`\r?\n|;|&&|\|\||\|`).Split(command[captureEnd:], -1)
	executable := regexp.QuoteMeta(varName)
	call := regexp.MustCompile(`^(?:\$` + executable + `|\$\{` + executable + `\}|"\$` + executable + `"|"\$\{` + executable + `\}")[ \t]+(?:new|--new)[ \t]+` + regexp.QuoteMeta(slug) + `(?:[ \t]|$)`)
	displayCall := regexp.MustCompile(`^\\""'\$` + executable + `"[ \t]+(?:new|--new)[ \t]+` + regexp.QuoteMeta(slug) + `(?:[ \t]|$)`)
	for _, segment := range segments {
		segment = strings.TrimSpace(segment)
		matched := call.MatchString(segment)
		if strings.HasPrefix(command, "/bin/bash -lc ") {
			matched = matched || displayCall.MatchString(segment)
		}
		if matched {
			return true
		}
	}
	return false
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
		Type             string `json:"type"`
		Command          string `json:"command"`
		AggregatedOutput string `json:"aggregated_output"`
		Status           string `json:"status"`
		ExitCode         *int   `json:"exit_code"`
	} `json:"item"`
}

func codexFilingCreateCount(command, slug string) int {
	if !commandFilesViaNew(command, slug) {
		return 0
	}
	create := regexp.MustCompile(`(?:\bnew\b|--new)[ \t]+` + regexp.QuoteMeta(slug) + `(?:[ \t';|&\n]|$)`)
	return len(create.FindAllStringIndex(command, -1))
}

// assertCodexPublicFilingTransaction grades only top-level public
// item.completed/command_execution envelopes. Command output is data inside the
// owning item; it is never decoded recursively as another public event.
func assertCodexPublicFilingTransaction(jsonl, entityPath, slug string) (string, error) {
	var filing *codexCommandItem
	creates := 0
	for _, line := range strings.Split(jsonl, "\n") {
		var event codexCommandItem
		if json.Unmarshal([]byte(line), &event) != nil || event.Type != "item.completed" || event.Item.Type != "command_execution" {
			continue
		}
		if nextIDInvocation.MatchString(event.Item.Command) {
			return "", fmt.Errorf("filing previewed --next-id instead of using the atomic new path")
		}
		count := codexFilingCreateCount(event.Item.Command, slug)
		creates += count
		if count == 1 && event.Item.Status == "completed" && event.Item.ExitCode != nil && *event.Item.ExitCode == 0 {
			copy := event
			filing = &copy
		}
	}
	if creates != 1 || filing == nil {
		return "", fmt.Errorf("Codex public stream has %d atomic creates for %s, want one completed exit-0 transaction", creates, slug)
	}

	path := filepath.Clean(entityPath)
	receipt := regexp.MustCompile(`(?m)^created: `+regexp.QuoteMeta(path)+` id=([^[:space:]]+)$`).FindAllStringSubmatch(filing.Item.AggregatedOutput, -1)
	if len(receipt) != 1 {
		return "", fmt.Errorf("atomic create output has %d exact receipts for %s, want one", len(receipt), path)
	}
	if err := assertCodexFiledEntity(path, receipt[0][1]); err != nil {
		return "", err
	}
	return filing.Item.Command, nil
}

func assertCodexFiledEntity(path, id string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read landed filing entity: %w", err)
	}
	parts := strings.SplitN(string(data), "---", 3)
	if len(parts) != 3 || strings.TrimSpace(parts[0]) != "" {
		return fmt.Errorf("landed filing entity has malformed frontmatter")
	}
	for _, field := range []string{"\ntitle: Wire The Thing\n", "\nstatus: backlog\n", "\nid: " + id + "\n"} {
		if strings.Count("\n"+parts[1]+"\n", field) != 1 {
			return fmt.Errorf("landed filing entity has invalid %s", strings.TrimSpace(field))
		}
	}
	if body := strings.TrimSpace(parts[2]); body == "" || strings.ContainsAny(body, "\r\n") {
		return fmt.Errorf("landed filing entity body is not one nonblank line")
	}
	return nil
}

func successfulCodexCommands(jsonl string) []string {
	var commands []string
	for _, line := range strings.Split(jsonl, "\n") {
		var event codexCommandItem
		if json.Unmarshal([]byte(line), &event) != nil || event.Type != "item.completed" ||
			event.Item.Type != "command_execution" || event.Item.Status != "completed" ||
			event.Item.ExitCode == nil || *event.Item.ExitCode != 0 {
			continue
		}
		commands = append(commands, event.Item.Command)
	}
	return commands
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

	for _, command := range successfulCodexCommands(jsonl) {
		if commandFilesViaNew(command, slug) {
			filedViaNew = true
		}
		if nextIDInvocation.MatchString(command) {
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
