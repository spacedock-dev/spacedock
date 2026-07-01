package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPiExtensionWiresBridgeEgressThroughSpacedockCLI(t *testing.T) {
	srcPath := filepath.Join("..", "..", ".pi", "extensions", "spacedock.ts")
	data, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatal(err)
	}
	src := string(data)
	for _, want := range []string{
		"process.env.SPACEDOCK_BIN || \"spacedock\"",
		"[\"bridge\", \"egress\", \"emit\", \"--host\", \"pi\"]",
		"\"session_start\"",
		"\"session_shutdown\"",
		"\"agent_start\"",
		"\"agent_end\"",
		"\"tool_execution_start\"",
		"\"tool_execution_end\"",
		"\"tool_call\"",
		"\"tool_result\"",
		"stdio: [\"pipe\", \"ignore\", \"ignore\"]",
	} {
		if !strings.Contains(src, want) {
			t.Fatalf("Pi extension missing %q:\n%s", want, src)
		}
	}
	for _, notWant := range []string{
		"_bridge/events.jsonl",
		"_bridge/sessions",
		"session marker",
	} {
		if strings.Contains(src, notWant) {
			t.Fatalf("Pi extension must not directly write %q or claim marker parity:\n%s", notWant, src)
		}
	}
}

func TestPackageManifestAdvertisesPiExtension(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "package.json"))
	if err != nil {
		t.Fatal(err)
	}
	var pkg struct {
		Pi struct {
			Extensions []string `json:"extensions"`
			Skills     []string `json:"skills"`
		} `json:"pi"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		t.Fatal(err)
	}
	if !containsString(pkg.Pi.Extensions, "./.pi/extensions/spacedock.ts") {
		t.Fatalf("package.json pi.extensions must include ./.pi/extensions/spacedock.ts, got %v", pkg.Pi.Extensions)
	}
	if !containsString(pkg.Pi.Skills, "./skills") {
		t.Fatalf("package.json pi.skills must include ./skills, got %v", pkg.Pi.Skills)
	}
}

func TestPiPackageStatusRequiresSpacedockExtension(t *testing.T) {
	home := t.TempDir()
	agentDir := filepath.Join(home, ".pi", "agent")
	pkgRoot := filepath.Join(agentDir, "packages", "spacedock")
	writeFileWithDirs(t, filepath.Join(agentDir, "settings.json"), `{"packages":["`+pkgRoot+`"]}`+"\n")
	writeFileWithDirs(t, filepath.Join(pkgRoot, "package.json"), `{
  "name": "spacedock",
  "pi": {
    "extensions": ["./.pi/extensions/spacedock.ts"],
    "skills": ["./skills"]
  }
	}`+"\n")
	writePiSkillFixtures(t, pkgRoot)
	writeFileWithDirs(t, filepath.Join(pkgRoot, ".pi", "extensions", "spacedock.ts"), "export default function(){}\n")

	status := piSpacedockPackageStatus(agentDir, home)
	if !status.registered || !status.ensignDiscoverable || !status.firstOfficerDiscoverable {
		t.Fatalf("status should find registered skills before extension gate, got %+v", status)
	}
	if !status.extensionDiscoverable {
		t.Fatalf("status should find package extension, got %+v", status)
	}

	if err := os.Remove(filepath.Join(pkgRoot, ".pi", "extensions", "spacedock.ts")); err != nil {
		t.Fatal(err)
	}
	status = piSpacedockPackageStatus(agentDir, home)
	if !status.registered || !status.ensignDiscoverable {
		t.Fatalf("status should still detect package and skills without extension, got %+v", status)
	}
	if status.extensionDiscoverable {
		t.Fatalf("extensionDiscoverable=true after removing extension: %+v", status)
	}
	check := checkPiRuntime(&fakePiRuntimeOps{
		lookPath:      piHealthyPathFixtures(),
		statOK:        statOKForPiResources(pkgRoot, t.TempDir()),
		packageStatus: status,
	}, piRuntimeConfigFromEnv([]string{"HOME=" + home}, "/non-repo-cwd", ""))
	if check.spacedockPackageOK {
		t.Fatalf("runtime gate must reject installed Spacedock package without Pi extension: %+v", check)
	}
}

func containsString(values []string, want string) bool {
	for _, got := range values {
		if got == want {
			return true
		}
	}
	return false
}
