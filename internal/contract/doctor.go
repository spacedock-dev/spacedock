// ABOUTME: spacedock doctor — read a plugin manifest's requires-contract, compare
// ABOUTME: against CONTRACT_VERSION, and report one of five verdicts with an exit code.
package contract

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// errNoManifest is returned by readRequiresContract when the manifest file does
// not exist — distinguishing the no-plugin-found report from a malformed field.
var errNoManifest = errors.New("manifest not found")

// readManifest reads a plugin manifest JSON and returns its display version and
// requires-contract string. A missing file yields errNoManifest; an absent
// requires-contract field yields an empty string (which Compare classifies as
// plugin-predates-contract — the installed plugin predates the contract
// mechanism — and routes to the actionable `spacedock install` upgrade remedy).
// The display version is the user-facing semver shown in the verdict message.
func readManifest(manifestPath string) (version, requiresContract string, err error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", errNoManifest
		}
		return "", "", err
	}
	var m struct {
		Version          string `json:"version"`
		RequiresContract string `json:"requires-contract"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return "", "", fmt.Errorf("parse manifest %s: %w", manifestPath, err)
	}
	return m.Version, m.RequiresContract, nil
}

// readRequiresContract reads a plugin manifest JSON and returns its
// requires-contract string, discarding the display version. Used by the
// fixture-bracketing test that only needs the range.
func readRequiresContract(manifestPath string) (string, error) {
	_, raw, err := readManifest(manifestPath)
	return raw, err
}

// ManifestVerdict resolves the compatibility verdict for the manifest at
// manifestPath against this binary's CONTRACT_VERSION, for the named host and
// (pre-release) dev branch. A missing manifest file yields NoPluginFound; an
// unparseable manifest JSON yields a MalformedRange-shaped Result naming the
// parse error. The plugin's display version (from the manifest) and the binary's
// display version (threaded in by the cli caller, which owns cli.Version) are
// woven into the user-facing message. The front door inspects the verdict
// directly (a non-empty path to a missing file is NoPluginFound, NOT compatible);
// RunDoctor maps the same verdict to an exit code and stream.
func ManifestVerdict(manifestPath, host, branch, binaryVersion string) Result {
	pluginVersion, raw, err := readManifest(manifestPath)
	if errors.Is(err, errNoManifest) {
		return Result{Verdict: NoPluginFound, Message: noPluginMessage(host)}
	}
	if err != nil {
		return Result{Verdict: MalformedRange, Message: fmt.Sprintf("error: %s", err)}
	}
	return compareWithManifest(CONTRACT_VERSION, raw, host, branch, manifestPath, pluginVersion, binaryVersion)
}

// RunDoctor reports the compatibility verdict for the manifest at manifestPath
// against this binary's CONTRACT_VERSION, for the named host and (pre-release)
// dev branch. binaryVersion is the binary's display version threaded in for the
// user-facing message. A compatible verdict and a no-plugin-found report exit 0
// (the report is non-fatal-by-default); every mismatch (too-old-binary,
// too-old-plugin, malformed-range) exits 1 with the pinned remedy on stderr.
func RunDoctor(manifestPath, host, branch, binaryVersion string, stdout, stderr io.Writer) int {
	res := ManifestVerdict(manifestPath, host, branch, binaryVersion)
	switch res.Verdict {
	case Compatible:
		fmt.Fprintln(stdout, res.Message)
		return 0
	case NoPluginFound:
		fmt.Fprintln(stdout, res.Message)
		return 0
	default:
		fmt.Fprintln(stderr, res.Message)
		return 1
	}
}
