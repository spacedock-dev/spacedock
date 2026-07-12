// ABOUTME: Shared safehouse seam — detect a workdir profile, gate the binary,
// ABOUTME: and wrap an inner command argv for `safehouse --trust-workdir-config`.
package safehouse

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// installHint is the pinned, actionable stderr message emitted when a workdir
// carries a .safehouse profile but the safehouse binary is not resolvable.
const installHint = "Spacedock: this directory has a .safehouse profile but the `safehouse` binary was not found on PATH. " +
	"Install safehouse (brew install eugene1g/safehouse/agent-safehouse; https://agent-safehouse.dev) or remove .safehouse to launch without it."

// terminalTargetingEnvSet describes one terminal/session identity source. Future
// multiplexers extend this data, leaving Wrap and its env-pass composition intact.
type terminalTargetingEnvSet struct {
	trigger string
	names   []string
}

var terminalTargetingEnvSets = []terminalTargetingEnvSet{
	{
		trigger: "ZELLIJ",
		names:   []string{"ZELLIJ", "ZELLIJ_PANE_ID", "ZELLIJ_SESSION_NAME"},
	},
}

// Present reports whether a .safehouse profile exists in workdir. A regular file
// or a directory both count (os.Stat truthiness) — the profile may be either.
func Present(workdir string) bool {
	_, err := os.Stat(filepath.Join(workdir, ".safehouse"))
	return err == nil
}

// Available reports whether the safehouse binary is resolvable via lookPath
// (production passes exec.LookPath; tests pin not-found). When the binary is
// absent it returns ok=false and a pinned install-hint string for stderr.
func Available(lookPath func(string) (string, error)) (ok bool, hint string) {
	if _, err := lookPath("safehouse"); err != nil {
		return false, installHint
	}
	return true, ""
}

// TranslateFlags turns the de-prefixed `--safehouse-*` knob tokens (the namespace
// prefix already stripped by the internal/cli dispatcher) into the safehouse
// `extra` argv fed verbatim into Wrap's pre-`--` slot. It owns only safehouse's
// flag vocabulary: `enable=ssh,docker` comma-splits into repeated `--enable=KEY`;
// `add-dirs=P` / `add-dirs-ro=P` map to `--add-dirs=P` / `--add-dirs-ro=P`. An
// unrecognized key is a hard error so a typo never silently reaches the host. It
// holds no `--safehouse-` namespace knowledge — a future sandbox's translator
// sits beside this one, each owning its own host's vocabulary.
func TranslateFlags(deprefixed []string) (extra []string, err error) {
	for _, tok := range deprefixed {
		key, value, ok := strings.Cut(tok, "=")
		if !ok {
			return nil, fmt.Errorf("safehouse: malformed flag %q (expected key=value)", tok)
		}
		switch key {
		case "enable":
			for _, v := range strings.Split(value, ",") {
				extra = append(extra, "--enable="+v)
			}
		case "add-dirs":
			extra = append(extra, "--add-dirs="+value)
		case "add-dirs-ro":
			extra = append(extra, "--add-dirs-ro="+value)
		default:
			return nil, fmt.Errorf("safehouse: unknown flag --safehouse-%s", key)
		}
	}
	return extra, nil
}

// withTerminalTargetingEnvPass conditionally adds the terminal/session metadata
// names to Safehouse's named env allowlist. The values stay in parent; Safehouse
// copies only explicitly named values across its scrub boundary. It merges with
// an existing --env-pass flag so the wrapper emits one allowlist for the current
// launch rather than relying on repeat-flag semantics.
func withTerminalTargetingEnvPass(extra, parent []string) []string {
	targetingNames := terminalTargetingEnvPassNames(parent)
	if targetingNames == "" {
		return extra
	}

	for i, arg := range extra {
		if arg == "--env-pass" && i+1 < len(extra) {
			out := append([]string(nil), extra...)
			out[i+1] = joinEnvPassNames(out[i+1], targetingNames)
			return out
		}
		if value, ok := strings.CutPrefix(arg, "--env-pass="); ok {
			out := append([]string(nil), extra...)
			out[i] = "--env-pass=" + joinEnvPassNames(value, targetingNames)
			return out
		}
	}

	out := append([]string(nil), extra...)
	return append(out, "--env-pass", targetingNames)
}

func terminalTargetingEnvPassNames(parent []string) string {
	var names []string
	for _, set := range terminalTargetingEnvSets {
		if hasEnv(parent, set.trigger) {
			names = append(names, set.names...)
		}
	}
	return strings.Join(names, ",")
}

func hasEnv(env []string, key string) bool {
	prefix := key + "="
	for _, entry := range env {
		if strings.HasPrefix(entry, prefix) {
			return true
		}
	}
	return false
}

func joinEnvPassNames(current, additions string) string {
	if current == "" {
		return additions
	}
	return current + "," + additions
}

// Wrap returns the inner argv wrapped as
// `safehouse --trust-workdir-config [extra...] -- <inner>`. Callers gate on
// Present (and Available) first; Wrap adds conditional terminal-targeting names
// to the allowlist and remains inner-command-agnostic for every host front door.
func Wrap(inner []string, extra []string) (argv []string) {
	extra = withTerminalTargetingEnvPass(extra, os.Environ())
	argv = []string{"safehouse", "--trust-workdir-config"}
	argv = append(argv, extra...)
	argv = append(argv, "--")
	argv = append(argv, inner...)
	return argv
}
