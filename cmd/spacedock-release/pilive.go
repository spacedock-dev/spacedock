// ABOUTME: pi-live Runtime Live E2E lane commands (pins, print-install,
// ABOUTME: install, guard, verify-manifest) folded into the release binary.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/pilive"
)

// piLivePins prints every pinned spec/version/integrity plus the readiness floors.
func piLivePins() int {
	for _, pkg := range pilive.Packages {
		fmt.Printf("%s %s %s\n", pkg.Spec, pkg.Version, pkg.Integrity)
	}
	fmt.Printf("floors node=%s pi-coding-agent=%s pi-subagents=%s\n",
		pilive.NodeEngineFloor, pilive.PiCodingAgentFloor, pilive.PiSubagentsFloor)
	return 0
}

// piLivePrintInstall prints the local-install commands for the pinned family.
// Docs reference this command instead of storing the version list.
func piLivePrintInstall() int {
	fmt.Printf("npm install -g %s@%s\n", pilive.PiCodingAgentSpec, pilive.PiCodingAgentVersion)
	fmt.Printf("npm install --prefix \"$HOME/.pi/agent/npm\" %s@%s %s@%s\n",
		pilive.PiSubagentsSpec, pilive.PiSubagentsVersion, pilive.PiIntercomSpec, pilive.PiIntercomVersion)
	return 0
}

func piLiveAgentDir() (string, error) {
	if d := os.Getenv("PI_CODING_AGENT_DIR"); d != "" {
		return d, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".pi", "agent"), nil
}

func piLiveInstall() int {
	dir, err := piLiveAgentDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacedock-release install:", err)
		return 1
	}
	opts := pilive.InstallOptions{
		PackDir:      filepath.Join(os.TempDir(), "pi-live-npm-packs"),
		PiNpmRoot:    filepath.Join(dir, "npm"),
		SettingsPath: filepath.Join(dir, "settings.json"),
	}
	if err := pilive.Install(pilive.ExecRunner{}, opts); err != nil {
		fmt.Fprintln(os.Stderr, "spacedock-release install:", err)
		return 1
	}
	return 0
}

func piLiveGuard() int {
	dir, err := piLiveAgentDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacedock-release guard:", err)
		return 1
	}
	nodeVersion, err := pilive.ExecRunner{}.Output("node", "--version")
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacedock-release guard: node --version:", err)
		return 1
	}
	globalRoot, err := pilive.ExecRunner{}.Output("npm", "root", "-g")
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacedock-release guard: npm root -g:", err)
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
		fmt.Fprintln(os.Stderr, "spacedock-release guard:", err)
		return 1
	}
	return 0
}

func piLiveVerifyManifest(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "spacedock-release verify-manifest: need exactly one <package-root>")
		return 2
	}
	resolved, err := pilive.VerifyRuntimeManifest(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacedock-release verify-manifest:", err)
		return 1
	}
	for _, path := range resolved {
		fmt.Println("verified runtime file " + path)
	}
	return 0
}
