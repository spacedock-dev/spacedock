// ABOUTME: Single source of truth for the pi-live Runtime Live E2E family pins.
// ABOUTME: Holds each package spec/version/sha512 and the readiness floors; every other copy derives from here.
package pilive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Family pins, resolved from the npm `latest` dist-tag at check time on
// 2026-10-04 with:
//
//	npm view @earendil-works/pi-coding-agent dist-tags.latest version dist.integrity --json
//	npm view pi-subagents dist-tags.latest version dist.integrity --json
//	npm view pi-intercom dist-tags.latest version dist.integrity --json
//
// `latest` is only how each version was resolved; the exact resolved version is
// what is stored below. The advisory found the captain's 1.0.1 had moved to
// 1.0.2 by fetch time, so 1.0.2 is pinned. pi-subagents and pi-intercom did not
// move. This is the ONLY place these numbers may be written; the workflow, the
// lane guard, fixtures, and docs must derive from these constants.
const (
	PiCodingAgentSpec      = "@earendil-works/pi-coding-agent"
	PiCodingAgentVersion   = "1.0.2"
	PiCodingAgentIntegrity = "sha512-3ZdIghMSELMGV3sKi5iASOb1Jwb696fLjmNu0aezaqDxTLLWWoRpqBYkGxJ1CgAMCbtfqXWEF0lcRrlVXmiEGQ=="

	PiSubagentsSpec      = "pi-subagents"
	PiSubagentsVersion   = "0.75.0"
	PiSubagentsIntegrity = "sha512-RO4DiTJM6pnK8y9PnD7Y6TLeiX2c8Kh6QhmkecFo+ou5qtnz5qiM+vUKru48UmJG1nB9eJ6rALd1M04uBOR4XQ=="

	PiIntercomSpec      = "pi-intercom"
	PiIntercomVersion   = "0.16.0"
	PiIntercomIntegrity = "sha512-ClGQuovPsz7r1iQwMRjEN+8wxywfrDrMILAkCSf/z19Nezzyfz77U8362Wb/VFNoZOwmF4PBsuGWCtB9AsEJMQ=="
)

// Readiness floors enforced by the lane's compatibility gate. pi-coding-agent's
// floor is also the spacedock ready-gate floor (s98) the launcher enforces; the
// pi-subagents floor is the last verified pairing.
const (
	NodeEngineFloor    = "22.19.0"
	PiCodingAgentFloor = "0.83.0"
	PiSubagentsFloor   = "0.53.0"
)

// Package is one pinned npm package: its install spec, exact published version,
// and published sha512 integrity.
type Package struct {
	Spec      string
	Version   string
	Integrity string
}

// Packages is the pinned family, in the lane's install order.
var Packages = []Package{
	{PiCodingAgentSpec, PiCodingAgentVersion, PiCodingAgentIntegrity},
	{PiSubagentsSpec, PiSubagentsVersion, PiSubagentsIntegrity},
	{PiIntercomSpec, PiIntercomVersion, PiIntercomIntegrity},
}

// SubstrateNpmSources are the two substrate package registrations the isolated
// Pi home's settings.json must carry so Pi's own discovery loads them.
var SubstrateNpmSources = []string{"npm:pi-subagents", "npm:pi-intercom"}

// ParsePackMetadata extracts the single tarball's filename and integrity from
// `npm pack --json` output. The JSON shape is a one-element array; anything
// else (empty, malformed, or multi-element) is rejected.
func ParsePackMetadata(data []byte) (PackMetadata, error) {
	var entries []PackMetadata
	if err := json.Unmarshal(data, &entries); err != nil {
		return PackMetadata{}, fmt.Errorf("parse npm pack --json output: %w", err)
	}
	if len(entries) != 1 {
		return PackMetadata{}, fmt.Errorf("npm pack --json returned %d entries, want exactly 1", len(entries))
	}
	if entries[0].Filename == "" || entries[0].Integrity == "" {
		return PackMetadata{}, fmt.Errorf("npm pack --json entry missing filename or integrity: %+v", entries[0])
	}
	return entries[0], nil
}

// PackMetadata is the subset of a `npm pack --json` entry the lane needs.
type PackMetadata struct {
	Filename  string `json:"filename"`
	Integrity string `json:"integrity"`
}

// VerifyIntegrity is the lane's independent registry oracle: the sha512
// re-derived by `npm pack` must equal the pinned integrity. A mismatch aborts
// before any install.
func VerifyIntegrity(actual, expected string) error {
	if actual != expected {
		return fmt.Errorf("integrity mismatch: expected %s got %s", expected, actual)
	}
	return nil
}

// ReadInstalledIdentity reads name and version from an installed package's
// package.json.
func ReadInstalledIdentity(root string) (name, version string, err error) {
	data, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return "", "", err
	}
	var manifest struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return "", "", fmt.Errorf("parse %s/package.json: %w", root, err)
	}
	return manifest.Name, manifest.Version, nil
}

// VerifyInstalled asserts the installed package at root reports wantName and
// wantVersion.
func VerifyInstalled(root, wantName, wantVersion string) error {
	name, version, err := ReadInstalledIdentity(root)
	if err != nil {
		return err
	}
	if name != wantName {
		return fmt.Errorf("installed package name %q at %s, want %q", name, root, wantName)
	}
	if version != wantVersion {
		return fmt.Errorf("installed %s version %q, want %q", wantName, version, wantVersion)
	}
	return nil
}

// VerifyRuntimeManifest asserts an installed substrate package declares usable
// runtime entrypoints: a nonempty pi.extensions array and the
// exports["./intercom-bridge"].default runtime string, with every declared
// extension and the bridge resolving relative to root to a regular file. It
// returns the resolved paths.
func VerifyRuntimeManifest(root string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return nil, err
	}
	var manifest struct {
		Pi struct {
			Extensions []string `json:"extensions"`
		} `json:"pi"`
		Exports map[string]json.RawMessage `json:"exports"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse %s/package.json: %w", root, err)
	}
	extensions := manifest.Pi.Extensions
	if len(extensions) == 0 {
		return nil, fmt.Errorf("%s/package.json declares no pi.extensions", root)
	}
	bridgeRaw, ok := manifest.Exports["./intercom-bridge"]
	var bridge struct {
		Default string `json:"default"`
	}
	if !ok || json.Unmarshal(bridgeRaw, &bridge) != nil || bridge.Default == "" {
		return nil, fmt.Errorf("%s/package.json declares no exports[\"./intercom-bridge\"].default runtime export", root)
	}
	var resolved []string
	for _, entry := range append(append([]string{}, extensions...), bridge.Default) {
		if entry == "" {
			return nil, fmt.Errorf("invalid runtime entry in %s/package.json", root)
		}
		target := filepath.Join(root, entry)
		info, err := os.Stat(target)
		if err != nil {
			return nil, fmt.Errorf("runtime entry %q does not resolve to a file: %w", entry, err)
		}
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("runtime entry %q resolves to %s, not a regular file", entry, target)
		}
		resolved = append(resolved, target)
	}
	return resolved, nil
}

// MergeSettings create-or-merges the substrate npm sources into the packages
// array at settingsPath, preserving unrelated keys and existing package entries.
// A missing file starts empty; malformed JSON is an error. It is idempotent.
func MergeSettings(settingsPath string, sources []string) ([]string, error) {
	settings := map[string]json.RawMessage{}
	data, err := os.ReadFile(settingsPath)
	switch {
	case err == nil:
		if len(strings.TrimSpace(string(data))) > 0 {
			if err := json.Unmarshal(data, &settings); err != nil {
				return nil, fmt.Errorf("parse %s: %w", settingsPath, err)
			}
		}
	case os.IsNotExist(err):
	default:
		return nil, err
	}
	var packages []string
	if raw, ok := settings["packages"]; ok {
		var entries []json.RawMessage
		if err := json.Unmarshal(raw, &entries); err == nil {
			for _, entry := range entries {
				var s string
				if json.Unmarshal(entry, &s) == nil {
					packages = append(packages, s)
				}
			}
		}
	}
	for _, source := range sources {
		found := false
		for _, existing := range packages {
			if existing == source {
				found = true
				break
			}
		}
		if !found {
			packages = append(packages, source)
		}
	}
	encoded, err := json.Marshal(packages)
	if err != nil {
		return nil, err
	}
	settings["packages"] = encoded
	out, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(settingsPath, append(out, '\n'), 0o644); err != nil {
		return nil, err
	}
	return packages, nil
}

// VersionAtLeast reports whether version parses at or above floor. A component
// that is not a number counts as 0 (the npm output is a plain semver triple).
func VersionAtLeast(version, floor string) bool {
	v, f := parseVersion(version), parseVersion(floor)
	for i := 0; i < 3; i++ {
		if v[i] > f[i] {
			return true
		}
		if v[i] < f[i] {
			return false
		}
	}
	return true
}

func parseVersion(version string) [3]int {
	core := strings.TrimPrefix(strings.TrimSpace(version), "v")
	if i := strings.IndexByte(core, '-'); i >= 0 {
		core = core[:i]
	}
	parts := strings.Split(core, ".")
	var out [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		n, err := strconv.Atoi(strings.TrimSpace(parts[i]))
		if err != nil {
			n = 0
		}
		out[i] = n
	}
	return out
}

// NodeEngineAtLeast reports whether the running Node satisfies the
// pi-coding-agent engine floor.
func NodeEngineAtLeast(nodeVersion string) bool {
	return VersionAtLeast(nodeVersion, NodeEngineFloor)
}

// CompatExportPath resolves @earendil-works/pi-ai's "./compat" import export to a
// file, checking the nested copy under pi-coding-agent first (npm usually does
// not hoist it to the global root), then the top-level global root. It returns
// the pi-ai package root and the resolved compat file path.
func CompatExportPath(agentRoot, globalNpmRoot string) (root, compatPath string, err error) {
	candidates := []string{
		filepath.Join(agentRoot, "node_modules", "@earendil-works", "pi-ai"),
		filepath.Join(globalNpmRoot, "@earendil-works", "pi-ai"),
	}
	for _, candidate := range candidates {
		data, readErr := os.ReadFile(filepath.Join(candidate, "package.json"))
		if readErr != nil {
			continue
		}
		var manifest struct {
			Version string `json:"version"`
			Exports map[string]struct {
				Import string `json:"import"`
			} `json:"exports"`
		}
		if err := json.Unmarshal(data, &manifest); err != nil {
			return "", "", fmt.Errorf("parse %s/package.json: %w", candidate, err)
		}
		compat, ok := manifest.Exports["./compat"]
		if !ok || compat.Import == "" {
			continue
		}
		path := filepath.Join(candidate, compat.Import)
		if _, statErr := os.Stat(path); statErr != nil {
			return "", "", fmt.Errorf("@earendil-works/pi-ai/compat resolves to missing file %s", path)
		}
		return candidate, path, nil
	}
	return "", "", fmt.Errorf("cannot locate @earendil-works/pi-ai with a ./compat import export (checked %s)", strings.Join(candidates, ", "))
}

// Runner runs an external command and returns its stdout. The lane supplies an
// exec-based runner; tests supply a fake.
type Runner interface {
	Output(name string, args ...string) ([]byte, error)
}

// ExecRunner is the production Runner. It returns stdout only; a failing
// command's error carries the captured stderr, so an npm warning emitted while
// resolving a path (e.g. `npm root -g`) cannot pollute the result.
type ExecRunner struct{}

func (ExecRunner) Output(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}

// PackTarball runs `npm pack` for pkg and returns the downloaded tarball path
// after verifying its published integrity against the pin — the lane's
// independent registry oracle. A mismatch aborts before any install.
func PackTarball(r Runner, pkg Package, packDir string) (string, error) {
	out, err := r.Output("npm", "pack", pkg.Spec+"@"+pkg.Version, "--pack-destination", packDir, "--ignore-scripts", "--json")
	if err != nil {
		return "", fmt.Errorf("npm pack %s@%s: %w: %s", pkg.Spec, pkg.Version, err, strings.TrimSpace(string(out)))
	}
	meta, err := ParsePackMetadata(out)
	if err != nil {
		return "", err
	}
	if err := VerifyIntegrity(meta.Integrity, pkg.Integrity); err != nil {
		return "", fmt.Errorf("%s@%s: %w", pkg.Spec, pkg.Version, err)
	}
	return filepath.Join(packDir, meta.Filename), nil
}

// CompatExportLoads evaluates the resolved compat module in Node, so a module
// that resolves but fails to load is caught. This is the one place the lane
// still needs a JS runtime: verifying an ESM load is inherently a Node action.
func CompatExportLoads(r Runner, compatPath string) error {
	fileURL := (&url.URL{Scheme: "file", Path: compatPath}).String()
	const script = `import(process.argv[1]).catch((error) => { console.error(error.message); process.exit(1); })`
	out, err := r.Output("node", "-e", script, fileURL)
	if err != nil {
		return fmt.Errorf("@earendil-works/pi-ai/compat exists but failed to load: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// InstallOptions names the directories the lane install writes into.
type InstallOptions struct {
	PackDir      string // where npm pack downloads tarballs
	PiNpmRoot    string // <agentDir>/npm, the substrate install prefix
	SettingsPath string // <agentDir>/settings.json
}

// Install performs the whole "Install Pi CLI and substrates" step: verify each
// pin's integrity, install pi-coding-agent globally and both substrates into the
// agent npm prefix, verify every installed name+version, assert the substrate
// runtime manifests, check the packaged skill files, and register both npm
// substrates in the agent settings so Pi's own discovery loads them. It never
// exports either package-root override.
func Install(r Runner, opts InstallOptions) error {
	if err := os.MkdirAll(opts.PackDir, 0o755); err != nil {
		return err
	}
	tarballs := map[string]string{}
	for _, pkg := range Packages {
		tarball, err := PackTarball(r, pkg, opts.PackDir)
		if err != nil {
			return err
		}
		tarballs[pkg.Spec] = tarball
	}

	if out, err := r.Output("npm", "install", "-g", tarballs[PiCodingAgentSpec], "--ignore-scripts", "--no-audit", "--no-fund", "--omit=dev"); err != nil {
		return fmt.Errorf("install %s: %w: %s", PiCodingAgentSpec, err, strings.TrimSpace(string(out)))
	}
	globalRoot, err := npmGlobalRoot(r)
	if err != nil {
		return err
	}
	if err := VerifyInstalled(filepath.Join(globalRoot, PiCodingAgentSpec), PiCodingAgentSpec, PiCodingAgentVersion); err != nil {
		return err
	}

	if err := os.MkdirAll(opts.PiNpmRoot, 0o755); err != nil {
		return err
	}
	if out, err := r.Output("npm", "install", "--prefix", opts.PiNpmRoot,
		tarballs[PiSubagentsSpec], tarballs[PiIntercomSpec],
		"--ignore-scripts", "--no-audit", "--no-fund", "--omit=dev"); err != nil {
		return fmt.Errorf("install substrates: %w: %s", err, strings.TrimSpace(string(out)))
	}
	subagentsRoot := filepath.Join(opts.PiNpmRoot, "node_modules", PiSubagentsSpec)
	intercomRoot := filepath.Join(opts.PiNpmRoot, "node_modules", PiIntercomSpec)
	if err := VerifyInstalled(subagentsRoot, PiSubagentsSpec, PiSubagentsVersion); err != nil {
		return err
	}
	if err := VerifyInstalled(intercomRoot, PiIntercomSpec, PiIntercomVersion); err != nil {
		return err
	}
	if _, err := VerifyRuntimeManifest(subagentsRoot); err != nil {
		return err
	}
	for _, skill := range []string{
		filepath.Join(subagentsRoot, "skills", "pi-subagents", "SKILL.md"),
		filepath.Join(intercomRoot, "skills", "pi-intercom", "SKILL.md"),
	} {
		if info, err := os.Stat(skill); err != nil || !info.Mode().IsRegular() {
			return fmt.Errorf("missing packaged skill file %s", skill)
		}
	}
	if _, err := MergeSettings(opts.SettingsPath, SubstrateNpmSources); err != nil {
		return err
	}
	fmt.Printf("installed %s@%s, %s@%s, %s@%s; registered %s\n",
		PiCodingAgentSpec, PiCodingAgentVersion, PiSubagentsSpec, PiSubagentsVersion,
		PiIntercomSpec, PiIntercomVersion, strings.Join(SubstrateNpmSources, ", "))
	return nil
}

func npmGlobalRoot(r Runner) (string, error) {
	out, err := r.Output("npm", "root", "-g")
	if err != nil {
		return "", fmt.Errorf("npm root -g: %w: %s", err, strings.TrimSpace(string(out)))
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", fmt.Errorf("npm root -g returned no path")
	}
	return root, nil
}

// GuardOptions names the installed roots and running Node version the
// compatibility guard reads.
type GuardOptions struct {
	AgentRoot     string
	SubagentsRoot string
	GlobalNpmRoot string
	NodeVersion   string
}

// Guard is the lane's compatibility guard: it enforces the Node engine floor,
// the pi-coding-agent ready-gate floor, and the pi-subagents pairing floor, then
// requires @earendil-works/pi-ai/compat to resolve AND load.
func Guard(r Runner, opts GuardOptions) error {
	if !NodeEngineAtLeast(opts.NodeVersion) {
		return fmt.Errorf("Node %s does not satisfy the pi-coding-agent >=%s engine requirement", opts.NodeVersion, NodeEngineFloor)
	}
	agentName, agentVersion, err := ReadInstalledIdentity(opts.AgentRoot)
	if err != nil {
		return err
	}
	if !VersionAtLeast(agentVersion, PiCodingAgentFloor) {
		return fmt.Errorf("%s %s is below the spacedock ready-gate floor %s", agentName, agentVersion, PiCodingAgentFloor)
	}
	subName, subVersion, err := ReadInstalledIdentity(opts.SubagentsRoot)
	if err != nil {
		return err
	}
	if !VersionAtLeast(subVersion, PiSubagentsFloor) {
		return fmt.Errorf("%s %s is below the last verified pairing floor %s", subName, subVersion, PiSubagentsFloor)
	}
	_, compatPath, err := CompatExportPath(opts.AgentRoot, opts.GlobalNpmRoot)
	if err != nil {
		return err
	}
	if err := CompatExportLoads(r, compatPath); err != nil {
		return err
	}
	fmt.Printf("verified %s %s pairs with @earendil-works/pi-ai/compat at %s\n", subName, subVersion, compatPath)
	return nil
}
