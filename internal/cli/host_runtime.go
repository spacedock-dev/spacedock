// ABOUTME: Per-runtime plugin-version detection for --version — claude/codex
// ABOUTME: manifest version + best-effort enabled marker and pi readiness, behind a probe.
package cli

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/contract"
)

// enabledMarker is the best-effort enabled/disabled marker rendered after a host's
// plugin version. The zero value is markerUnknown: the probe could not read the
// host's enabled-state (a failed or unparseable `plugin list`), so the version
// renders bare with no marker rather than claiming a state we did not observe.
// markerEnabled (the normal case) also renders bare — only markerDisabled appends
// the ` (disabled)` marker, since "enabled" is the unremarkable default.
type enabledMarker int

const (
	markerUnknown enabledMarker = iota
	markerEnabled
	markerDisabled
)

// runtimeStatus is a host's --version posture: whether the host binary resolves,
// the installed plugin version (from the resolved manifest; "" when no plugin), a
// best-effort enabled/disabled marker for that plugin, and pi's launch-readiness
// (pi has no marketplace version, so it renders `ready` instead). version, marker,
// and ready are meaningless when installed is false.
type runtimeStatus struct {
	installed bool
	version   string
	marker    enabledMarker
	ready     bool
}

// runtimeProbe reports a host's install + plugin-version posture. Production backs
// it with the real host CLIs (execRuntimeProbe); tests back it with a fake that
// pins each host's outcome, so --version never shells a live host CLI in the test
// path.
type runtimeProbe interface {
	ProbeRuntime(host string) runtimeStatus
}

// runtimeLine renders one host's --version line from its status, version-forward.
// An absent host binary is `not installed`. A host with a resolved plugin version
// is `spacedock <version>`, appending ` (disabled)` only when the best-effort
// marker confidently read disabled (markerEnabled and markerUnknown both render
// bare — enabled is the normal case, and an unread marker omits rather than
// invents). A host present with no plugin version is `spacedock not installed`;
// pi's launch-ready model has no version, so it renders `spacedock ready`.
func runtimeLine(host string, s runtimeStatus) string {
	if !s.installed {
		return host + ": not installed"
	}
	if s.version != "" {
		line := host + ": spacedock " + s.version
		if s.marker == markerDisabled {
			line += " (disabled)"
		}
		return line
	}
	if s.ready {
		return host + ": spacedock ready"
	}
	return host + ": spacedock not installed"
}

// claudeMarker reads a `claude plugin list --json` body and resolves the
// spacedock@spacedock entry's best-effort enabled marker from its `enabled`
// boolean: an entry with `enabled:true` is markerEnabled, `enabled:false` is
// markerDisabled. An absent spacedock entry is markerUnknown (the version comes
// from the resolved manifest, so a missing list entry must not be reported as
// disabled). A body that does not parse is an error so the caller renders the bare
// version (markerUnknown) rather than claiming a state it could not read.
func claudeMarker(body []byte) (enabledMarker, error) {
	var entries []pluginListEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return markerUnknown, fmt.Errorf("parse claude plugin list --json: %w", err)
	}
	for _, e := range entries {
		if e.ID == "spacedock@spacedock" {
			if e.Enabled {
				return markerEnabled, nil
			}
			return markerDisabled, nil
		}
	}
	return markerUnknown, nil
}

// codexEntryDisabled reports whether the `codex plugin list` text output marks the
// given plugin id as disabled. It mirrors codexEntryInstalled's field-based parse:
// the id must be a whitespace-delimited field, and a following field (within the
// row, stripped of surrounding `()` and a trailing `,`) must equal `disabled`. An
// absent `disabled` field is not a disabled plugin (an installed row reads
// `installed, enabled`), so the caller treats the unmarked case as enabled.
func codexEntryDisabled(listing, id string) bool {
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
			if strings.Trim(f, "(),") == "disabled" {
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

// ProbeRuntime resolves the host binary on PATH, then (when present) reads the
// installed plugin version and a best-effort enabled marker. A binary that does
// not resolve is not installed. The version is read ROBUSTLY from the resolved
// plugin manifest (execHost.ResolveManifest → the manifest's version, the same
// source doctor and the contract gate read), so it renders even when the
// enabled-state probe (`plugin list`) errors — in that case the marker stays
// markerUnknown and the version shows bare. pi has no marketplace version, so it
// reports launch-readiness instead.
func (execRuntimeProbe) ProbeRuntime(host string) runtimeStatus {
	if _, err := exec.LookPath(host); err != nil {
		return runtimeStatus{installed: false}
	}
	switch host {
	case "claude":
		return runtimeStatus{installed: true, version: probeVersion(host), marker: probeClaudeMarker()}
	case "codex":
		return runtimeStatus{installed: true, version: probeVersion(host), marker: probeCodexMarker()}
	case "pi":
		return runtimeStatus{installed: true, ready: probePiReady()}
	default:
		return runtimeStatus{installed: true}
	}
}

// probeVersion resolves the installed plugin manifest for host and reads its
// version, the robust version source `doctor` reads. A resolve error or a missing
// manifest yields "" (no plugin version to show), so the caller renders `spacedock
// not installed`. The marker is read separately and best-effort, so a `plugin
// list` failure never suppresses the version this returns.
func probeVersion(host string) string {
	manifestPath, err := execHost{}.ResolveManifest(host)
	if err != nil || manifestPath == "" {
		return ""
	}
	version, err := contract.ManifestVersion(manifestPath)
	if err != nil {
		return ""
	}
	return version
}

// probeClaudeMarker shells `claude plugin list --json` and reads the spacedock
// entry's `enabled` field into a best-effort marker. A failed command or
// unparseable body is markerUnknown (the version still renders from the manifest;
// only the marker is omitted).
func probeClaudeMarker() enabledMarker {
	out, err := exec.Command("claude", "plugin", "list", "--json").Output()
	if err != nil {
		return markerUnknown
	}
	marker, err := claudeMarker(out)
	if err != nil {
		return markerUnknown
	}
	return marker
}

// probeCodexMarker shells `codex plugin list` and reads the spacedock entry's
// disabled status from the text listing. A failed command (the sandbox-denied
// "Operation not permitted" mode) is markerUnknown; an explicit disabled row is
// markerDisabled; otherwise the installed-and-enabled default is markerEnabled.
func probeCodexMarker() enabledMarker {
	out, err := exec.Command("codex", "plugin", "list").CombinedOutput()
	if err != nil {
		return markerUnknown
	}
	if codexEntryDisabled(string(out), "spacedock@spacedock") {
		return markerDisabled
	}
	return markerEnabled
}

// probePiReady runs the existing pi-runtime readiness check (pi has no plugin-list
// or marketplace-version model — verified in the spike), so pi's `spacedock ready`
// reuses piRuntimeLaunchReady (skills + extension present), NOT a plugin probe.
func probePiReady() bool {
	cfg := piRuntimeConfigFromEnv(nil, cwd(), "")
	return piRuntimeLaunchReady(checkPiRuntime(execPiRuntimeOps{}, cfg))
}
