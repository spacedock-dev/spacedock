// ABOUTME: pi-live lane helper: installs and verifies the pinned Pi family and
// ABOUTME: runs the substrate compatibility guard. CI tooling, not the user binary.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/pilive"
)

// main is the pi-live helper invoked by .github/workflows/runtime-live-e2e.yml.
// All pins, integrity values, and readiness floors come from internal/pilive —
// the workflow itself carries no version, integrity, or JavaScript.
func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "pins":
		os.Exit(runPins())
	case "print-install":
		os.Exit(runPrintInstall())
	case "install":
		os.Exit(runInstall())
	case "guard":
		os.Exit(runGuard())
	case "verify-manifest":
		os.Exit(runVerifyManifest(os.Args[2:]))
	default:
		fmt.Fprintf(os.Stderr, "spacedock-pilive: unknown command %q\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

// runPins prints every pinned spec/version/integrity plus the readiness floors.
func runPins() int {
	for _, pkg := range pilive.Packages {
		fmt.Printf("%s %s %s\n", pkg.Spec, pkg.Version, pkg.Integrity)
	}
	fmt.Printf("floors node=%s pi-coding-agent=%s pi-subagents=%s\n", pilive.NodeEngineFloor, pilive.PiCodingAgentFloor, pilive.PiSubagentsFloor)
	return 0
}

// runPrintInstall prints the local-install commands for the pinned family. Docs
// reference this command instead of storing the version list.
func runPrintInstall() int {
	fmt.Printf("npm install -g %s@%s\n", pilive.PiCodingAgentSpec, pilive.PiCodingAgentVersion)
	fmt.Printf("npm install --prefix \"$HOME/.pi/agent/npm\" %s@%s %s@%s\n",
		pilive.PiSubagentsSpec, pilive.PiSubagentsVersion, pilive.PiIntercomSpec, pilive.PiIntercomVersion)
	return 0
}

func agentDir() (string, error) {
	if d := os.Getenv("PI_CODING_AGENT_DIR"); d != "" {
		return d, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".pi", "agent"), nil
}

func runInstall() int {
	dir, err := agentDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacedock-pilive install:", err)
		return 1
	}
	opts := pilive.InstallOptions{
		PackDir:      filepath.Join(os.TempDir(), "pi-live-npm-packs"),
		PiNpmRoot:    filepath.Join(dir, "npm"),
		SettingsPath: filepath.Join(dir, "settings.json"),
	}
	if err := pilive.Install(pilive.ExecRunner{}, opts); err != nil {
		fmt.Fprintln(os.Stderr, "spacedock-pilive install:", err)
		return 1
	}
	return 0
}

func runGuard() int {
	dir, err := agentDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacedock-pilive guard:", err)
		return 1
	}
	nodeVersion, err := pilive.ExecRunner{}.Output("node", "--version")
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacedock-pilive guard: node --version:", err)
		return 1
	}
	globalRoot, err := pilive.ExecRunner{}.Output("npm", "root", "-g")
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacedock-pilive guard: npm root -g:", err)
		return 1
	}
	global := strings.TrimSpace(string(globalRoot))
	opts := pilive.GuardOptions{
		AgentRoot:     filepath.Join(global, pilive.PiCodingAgentSpec),
		SubagentsRoot: filepath.Join(dir, "npm", "node_modules", pilive.PiSubagentsSpec),
		GlobalNpmRoot: global,
		NodeVersion:   strings.TrimSpace(string(nodeVersion)),
	}
	if err := pilive.Guard(pilive.ExecRunner{}, opts); err != nil {
		fmt.Fprintln(os.Stderr, "spacedock-pilive guard:", err)
		return 1
	}
	return 0
}

func runVerifyManifest(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "spacedock-pilive verify-manifest: need exactly one <package-root>")
		return 2
	}
	resolved, err := pilive.VerifyRuntimeManifest(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacedock-pilive verify-manifest:", err)
		return 1
	}
	for _, path := range resolved {
		fmt.Println("verified runtime file " + path)
	}
	return 0
}

func usage() {
	fmt.Fprint(os.Stderr, `spacedock-pilive is the pi-live Runtime Live E2E lane helper.

Usage:
  spacedock-pilive pins
  spacedock-pilive print-install
  spacedock-pilive install
  spacedock-pilive guard
  spacedock-pilive verify-manifest <package-root>
`)
}
