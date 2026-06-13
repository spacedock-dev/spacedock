// ABOUTME: Per-runtime install/enablement detection for --version — claude/codex
// ABOUTME: plugin-list enablement and pi readiness, behind an injectable probe seam.
package cli

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// enablement is a host's spacedock-plugin enablement posture. The zero value is
// notEnabled (the safe default: a host we could read but found no enabled
// spacedock entry is not enabled). enablementUnknown is reserved for the probe
// that could not determine enablement — a binary that resolves but whose
// enablement read errored (the sandboxed `codex plugin list` "Operation not
// permitted" mode), distinct from a confidently not-enabled host.
type enablement int

const (
	enablementNotEnabled enablement = iota
	enablementEnabled
	enablementUnknown
)

// runtimeStatus is a host's --version posture: whether the host binary resolves,
// and (when it does) its spacedock-plugin enablement. enablement is meaningless
// when installed is false.
type runtimeStatus struct {
	installed  bool
	enablement enablement
}

// runtimeProbe reports a host's install + enablement posture. Production backs it
// with the real host CLIs (execRuntimeProbe); tests back it with a fake that pins
// each host's outcome, so --version never shells a live host CLI in the test path.
type runtimeProbe interface {
	ProbeRuntime(host string) runtimeStatus
}

// runtimeLine renders one host's --version line from its status. An absent binary
// is `not installed`; a probe that could not read enablement is `installed,
// enablement unknown` (never silently not installed); an enabled plugin is
// `installed, spacedock enabled`; otherwise the bare `installed`.
func runtimeLine(host string, s runtimeStatus) string {
	if !s.installed {
		return host + ": not installed"
	}
	switch s.enablement {
	case enablementEnabled:
		return host + ": installed, spacedock enabled"
	case enablementUnknown:
		return host + ": installed, enablement unknown"
	default:
		return host + ": installed"
	}
}

// claudeEnablement reads a `claude plugin list --json` body and resolves the
// spacedock@spacedock entry's enablement from its `enabled` boolean (AC-4): an
// entry with `enabled:true` is enabled, `enabled:false` (or no spacedock entry at
// all) is not enabled. A body that does not parse is an error so the caller renders
// `enablement unknown` rather than silently downgrading to not-enabled.
func claudeEnablement(body []byte) (enablement, error) {
	var entries []pluginListEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return enablementUnknown, fmt.Errorf("parse claude plugin list --json: %w", err)
	}
	for _, e := range entries {
		if e.ID == "spacedock@spacedock" {
			if e.Enabled {
				return enablementEnabled, nil
			}
			return enablementNotEnabled, nil
		}
	}
	return enablementNotEnabled, nil
}

// codexEntryEnabled reports whether the `codex plugin list` text output marks the
// given plugin id as enabled. It mirrors codexEntryInstalled's field-based parse:
// the id must be a whitespace-delimited field, and a following field (within the
// row, stripped of surrounding `()` and a trailing `,`) must equal `enabled`. The
// codex status renders as `installed, enabled` (table form) or `(installed,
// enabled)` (legacy paren form), so `enabled` is the field after `installed`.
func codexEntryEnabled(listing, id string) bool {
	for _, line := range strings.Split(listing, "\n") {
		fields := strings.Fields(line)
		idIdx := -1
		for i, f := range fields {
			if f == id {
				idIdx = i
				break
			}
		}
		if idIdx < 0 {
			continue
		}
		for _, f := range fields[idIdx+1:] {
			if strings.Trim(f, "(),") == "enabled" {
				return true
			}
		}
	}
	return false
}

// execRuntimeProbe backs runtimeProbe with the real host CLIs and exec. It is the
// production seam --version uses; the test path injects a fake instead.
type execRuntimeProbe struct{}

var _ runtimeProbe = execRuntimeProbe{}

// ProbeRuntime resolves the host binary on PATH, then (when present) reads its
// spacedock-plugin enablement. A binary that does not resolve is not installed; a
// binary that resolves but whose enablement read errors is `enablement unknown`
// (the sandbox-denied probe), never silently downgraded to not-installed.
func (execRuntimeProbe) ProbeRuntime(host string) runtimeStatus {
	if _, err := exec.LookPath(host); err != nil {
		return runtimeStatus{installed: false}
	}
	switch host {
	case "claude":
		return runtimeStatus{installed: true, enablement: probeClaudeEnablement()}
	case "codex":
		return runtimeStatus{installed: true, enablement: probeCodexEnablement()}
	case "pi":
		return runtimeStatus{installed: true, enablement: probePiEnablement()}
	default:
		return runtimeStatus{installed: true, enablement: enablementUnknown}
	}
}

// probeClaudeEnablement shells `claude plugin list --json` and reads the spacedock
// entry's `enabled` field. A failed command or unparseable body is `enablement
// unknown` (the host resolved, but enablement could not be determined).
func probeClaudeEnablement() enablement {
	out, err := exec.Command("claude", "plugin", "list", "--json").Output()
	if err != nil {
		return enablementUnknown
	}
	state, err := claudeEnablement(out)
	if err != nil {
		return enablementUnknown
	}
	return state
}

// probeCodexEnablement shells `codex plugin list` and reads the spacedock entry's
// `enabled` status from the text listing. A failed command (the sandbox-denied
// "Operation not permitted" mode) is `enablement unknown`; an installed-but-not-
// enabled entry is not enabled.
func probeCodexEnablement() enablement {
	out, err := exec.Command("codex", "plugin", "list").CombinedOutput()
	if err != nil {
		return enablementUnknown
	}
	if codexEntryEnabled(string(out), "spacedock@spacedock") {
		return enablementEnabled
	}
	return enablementNotEnabled
}

// probePiEnablement runs the existing pi-runtime readiness check (pi has no
// plugin-list model — verified in the spike), so pi's "spacedock enabled" reuses
// piRuntimeLaunchReady (skills + extension present), NOT a plugin probe. Pi
// readiness is a boolean, so it never resolves to `enablement unknown`.
func probePiEnablement() enablement {
	cfg := piRuntimeConfigFromEnv(nil, cwd(), "")
	if piRuntimeLaunchReady(checkPiRuntime(execPiRuntimeOps{}, cfg)) {
		return enablementEnabled
	}
	return enablementNotEnabled
}
