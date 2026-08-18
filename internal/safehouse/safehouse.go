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

// terminalHostEnvVars are the nine signals subspace's r skill probes across
// its six terminal hosts, in the probe's own resolution order.
// Source of truth: spacedock-subspace plugins/subspace/skills/r/SKILL.md,
// "Select one terminal" — duplicated by decision (see safehouse-terminal-env-passthrough.md).
var terminalHostEnvVars = []string{
	"ZELLIJ_SESSION_NAME", "ZELLIJ_PANE_ID",
	"TMUX", "TMUX_PANE",
	"HERDR_ENV", "HERDR_PANE_ID",
	"CMUX_WORKSPACE_ID", "CMUX_SURFACE_ID",
	"TERM_PROGRAM",
}

// terminalTargetingEnvArgs returns Safehouse's built-in terminal/session
// metadata allowance. Safehouse itself composes repeated --env-pass flags and
// SAFEHOUSE_ENV_PASS, so this wrapper adds its default without parsing caller
// arguments or owning an operator configuration surface.
func terminalTargetingEnvArgs() []string {
	return terminalEnvPassArgs(os.LookupEnv)
}

// terminalEnvPassArgs is the presence-filter composer: it names a variable in
// the returned --env-pass allowance only when lookup reports it set in the
// parent, in terminalHostEnvVars' probe order. An empty parent (nothing set)
// yields no allowance at all, never an empty flag.
func terminalEnvPassArgs(lookup func(string) (string, bool)) []string {
	var present []string
	for _, name := range terminalHostEnvVars {
		if _, ok := lookup(name); ok {
			present = append(present, name)
		}
	}
	if len(present) == 0 {
		return nil
	}
	return []string{"--env-pass=" + strings.Join(present, ",")}
}

// Wrap returns the inner argv wrapped as
// `safehouse --trust-workdir-config [extra...] -- <inner>`. Callers gate on
// Present (and Available) first; Wrap adds the conditional built-in terminal
// targeting allowance and remains inner-command-agnostic for every host front door.
func Wrap(inner []string, extra []string) (argv []string) {
	argv = []string{"safehouse", "--trust-workdir-config"}
	argv = append(argv, terminalTargetingEnvArgs()...)
	argv = append(argv, extra...)
	argv = append(argv, "--")
	argv = append(argv, inner...)
	return argv
}
