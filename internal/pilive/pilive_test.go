// ABOUTME: Behavioural tests for the pi-live lane commands the release binary
// ABOUTME: exposes — install, guard, verify-manifest — one negative per check.
package pilive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeRunner struct {
	run func(name string, args ...string) ([]byte, error)
}

func (f fakeRunner) Output(name string, args ...string) ([]byte, error) {
	if f.run == nil {
		return nil, fmt.Errorf("unexpected command: %s %s", name, strings.Join(args, " "))
	}
	return f.run(name, args...)
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// substrateRoot writes an installed pi-subagents package.json and its declared
// runtime files: two extensions and the intercom bridge.
func substrateRoot(t *testing.T, root string) {
	t.Helper()
	writeFile(t, filepath.Join(root, "package.json"),
		`{"name":"pi-subagents","version":"`+PiSubagentsVersion+`",`+
			`"pi":{"extensions":["./index.js"]},`+
			`"exports":{"./intercom-bridge":{"default":"./src/api/intercom-bridge.js"}}}`)
	writeFile(t, filepath.Join(root, "index.js"), "x")
	writeFile(t, filepath.Join(root, "src/api/intercom-bridge.js"), "x")
}

func TestVerifyManifestResolvesDeclaredRuntimeFiles(t *testing.T) {
	root := t.TempDir()
	substrateRoot(t, root)
	resolved, err := verifyManifest(root)
	if err != nil || len(resolved) != 2 {
		t.Fatalf("resolved=%v err=%v, want the extension and bridge", resolved, err)
	}
	if os.Remove(filepath.Join(root, "src/api/intercom-bridge.js")); verifyManifestErr(t, root) == nil {
		t.Fatal("a missing declared runtime file must fail setup")
	}
	writeFile(t, filepath.Join(root, "package.json"), `{"name":"pi-subagents","version":"0.0.0"}`)
	if verifyManifestErr(t, root) == nil {
		t.Fatal("a manifest with no pi.extensions must fail setup")
	}
}

func verifyManifestErr(t *testing.T, root string) error {
	t.Helper()
	_, err := verifyManifest(root)
	return err
}

func TestPackRejectsIntegrityMismatch(t *testing.T) {
	packOut := `[{"filename":"a.tgz","integrity":"sha512-corrupted"}]`
	r := fakeRunner{run: func(name string, args ...string) ([]byte, error) { return []byte(packOut), nil }}
	if _, err := pack(r, packages[1], t.TempDir()); err == nil {
		t.Fatal("a tarball whose integrity does not match the pin must abort before install")
	}
}

func TestMergeSettingsPreservesEntriesAndIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if pkgs, err := mergeSettings(path, substrateNpmSources); err != nil || strings.Join(pkgs, ",") != "npm:pi-subagents,npm:pi-intercom" {
		t.Fatalf("fresh packages=%v err=%v", pkgs, err)
	}
	writeFile(t, path, `{"theme":"dark","packages":["npm:other",{"source":"npm:other-object"}]}`)
	pkgs, err := mergeSettings(path, substrateNpmSources)
	if err != nil || strings.Join(pkgs, ",") != "npm:other,npm:other-object,npm:pi-subagents,npm:pi-intercom" {
		t.Fatalf("merged packages=%v err=%v", pkgs, err)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), `"theme": "dark"`) {
		t.Fatalf("merge dropped an unrelated key: %s", data)
	}
	if !strings.Contains(string(data), `"source": "npm:other-object"`) {
		t.Fatalf("merge dropped an existing object package entry: %s", data)
	}
	if _, err := mergeSettings(path, substrateNpmSources); err != nil {
		t.Fatalf("merge must be idempotent: %v", err)
	}
	again, _ := os.ReadFile(path)
	if string(data) != string(again) {
		t.Fatalf("a second merge duplicated registrations:\n%s\n---\n%s", data, again)
	}
	writeFile(t, path, "{not json")
	if _, err := mergeSettings(path, substrateNpmSources); err == nil {
		t.Fatal("malformed settings JSON must be an error")
	}
}

func TestVersionAtLeastComparesFloors(t *testing.T) {
	for _, tc := range []struct {
		version, floor string
		want           bool
	}{
		{PiCodingAgentFloor, PiCodingAgentFloor, true}, {"0.82.9", PiCodingAgentFloor, false},
		{"1.0.2", PiCodingAgentFloor, true}, {"0.0.0", PiSubagentsFloor, false},
		{PiSubagentsVersion, PiSubagentsFloor, true}, {"22.18.0", NodeEngineFloor, false},
		{"24.13.1", NodeEngineFloor, true},
	} {
		if got := versionAtLeast(tc.version, tc.floor); got != tc.want {
			t.Errorf("versionAtLeast(%q, %q) = %v, want %v", tc.version, tc.floor, got, tc.want)
		}
	}
}

func TestGuardRejectsBelowFloorNodeAndSubstrate(t *testing.T) {
	agentDir := t.TempDir()
	nodeBelow := fakeRunner{run: func(name string, args ...string) ([]byte, error) { return []byte("22.18.0\n"), nil }}
	if err := guard(nodeBelow, agentDir); err == nil || !strings.Contains(err.Error(), NodeEngineFloor) {
		t.Fatalf("err=%v, want the Node engine floor %s", err, NodeEngineFloor)
	}
	writeFile(t, filepath.Join(agentDir, "npm", "node_modules", PiSubagentsSpec, "package.json"),
		`{"name":"pi-subagents","version":"0.52.0"}`)
	globalRoot := filepath.Join(agentDir, "global")
	writeFile(t, filepath.Join(globalRoot, PiCodingAgentSpec, "package.json"),
		`{"name":"`+PiCodingAgentSpec+`","version":"`+PiCodingAgentVersion+`"}`)
	r := fakeRunner{run: func(name string, args ...string) ([]byte, error) {
		if name == "node" {
			return []byte("v24.13.1\n"), nil
		}
		return []byte(globalRoot + "\n"), nil
	}}
	if err := guard(r, agentDir); err == nil || !strings.Contains(err.Error(), PiSubagentsFloor) {
		t.Fatalf("err=%v, want the pi-subagents floor %s", err, PiSubagentsFloor)
	}
}

func TestCompatExportPathPrefersInstalledAndRejectsIncompatible(t *testing.T) {
	globalRoot := t.TempDir()
	global := filepath.Join(globalRoot, "@earendil-works", "pi-ai")
	writeFile(t, filepath.Join(global, "package.json"), `{"exports":{"./compat":{"import":"compat.js"}}}`)
	writeFile(t, filepath.Join(global, "compat.js"), "x")
	agentRoot := filepath.Join(t.TempDir(), "pi-coding-agent")
	nested := filepath.Join(agentRoot, "node_modules", "@earendil-works", "pi-ai")
	writeFile(t, filepath.Join(nested, "package.json"), `{"exports":{"./compat":{"import":"compat.js"}}}`)
	writeFile(t, filepath.Join(nested, "compat.js"), "x")
	if path, err := compatExportPath(agentRoot, globalRoot); err != nil || path != filepath.Join(nested, "compat.js") {
		t.Fatalf("installed pi-ai must win: path=%s err=%v", path, err)
	}
	if path, err := compatExportPath(filepath.Join(t.TempDir(), "pi-coding-agent"), globalRoot); err != nil || path != filepath.Join(global, "compat.js") {
		t.Fatalf("global fallback when nothing is installed: path=%s err=%v", path, err)
	}
	writeFile(t, filepath.Join(nested, "package.json"), `{"exports":{}}`)
	if path, err := compatExportPath(agentRoot, globalRoot); err == nil {
		t.Fatalf("installed pi-ai without ./compat must fail, not use the global copy (%s)", path)
	}
}

func TestCompatExportLoadsPropagatesNodeFailure(t *testing.T) {
	r := fakeRunner{run: func(name string, args ...string) ([]byte, error) {
		return []byte("boom"), fmt.Errorf("exit status 1")
	}}
	if err := compatExportLoads(r, "/tmp/compat.js"); err == nil {
		t.Fatal("a resolved compat module that fails to load must fail the guard")
	}
}

func TestInstallVerifiesAndRegistersThePinnedFamily(t *testing.T) {
	agentDir, globalRoot, packDir := t.TempDir(), t.TempDir(), filepath.Join(os.TempDir(), "pi-live-npm-packs")
	writeFile(t, filepath.Join(globalRoot, PiCodingAgentSpec, "package.json"),
		`{"name":"`+PiCodingAgentSpec+`","version":"`+PiCodingAgentVersion+`"}`)
	subRoot := filepath.Join(agentDir, "npm", "node_modules", PiSubagentsSpec)
	substrateRoot(t, subRoot)
	writeFile(t, filepath.Join(subRoot, "skills/pi-subagents/SKILL.md"), "x")
	icRoot := filepath.Join(agentDir, "npm", "node_modules", PiIntercomSpec)
	writeFile(t, filepath.Join(icRoot, "package.json"), `{"name":"pi-intercom","version":"`+PiIntercomVersion+`"}`)
	writeFile(t, filepath.Join(icRoot, "skills/pi-intercom/SKILL.md"), "x")

	var installs [][]string
	r := fakeRunner{run: func(name string, args ...string) ([]byte, error) {
		switch {
		case name == "npm" && args[0] == "pack":
			spec := args[1][:strings.LastIndex(args[1], "@")]
			for _, p := range packages {
				if p.spec == spec && args[1] == p.spec+"@"+p.version {
					return []byte(fmt.Sprintf(`[{"filename":%q,"integrity":%q}]`, args[1]+".tgz", p.integrity)), nil
				}
			}
			return nil, fmt.Errorf("unexpected pack %q", args[1])
		case name == "npm" && args[0] == "root":
			return []byte(globalRoot + "\n"), nil
		case name == "npm" && args[0] == "install":
			installs = append(installs, args)
			return nil, nil
		}
		return nil, fmt.Errorf("unexpected command %s %v", name, args)
	}}
	if err := install(r, agentDir); err != nil {
		t.Fatalf("install: %v", err)
	}
	if len(installs) != 2 || installs[0][1] != "-g" || installs[0][2] != filepath.Join(packDir, PiCodingAgentSpec+"@"+PiCodingAgentVersion+".tgz") {
		t.Fatalf("global install = %v, want exactly the pinned pi-coding-agent tarball", installs)
	}
	if installs[1][1] != "--prefix" || installs[1][2] != filepath.Join(agentDir, "npm") ||
		installs[1][3] != filepath.Join(packDir, PiSubagentsSpec+"@"+PiSubagentsVersion+".tgz") ||
		installs[1][4] != filepath.Join(packDir, PiIntercomSpec+"@"+PiIntercomVersion+".tgz") {
		t.Fatalf("substrate install = %v, want both pinned substrate tarballs", installs[1])
	}
	settings, err := os.ReadFile(filepath.Join(agentDir, "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, source := range substrateNpmSources {
		if !strings.Contains(string(settings), source) {
			t.Fatalf("settings %s missing %s", settings, source)
		}
	}
	writeFile(t, filepath.Join(subRoot, "package.json"), `{"name":"pi-subagents","version":"0.0.0"}`)
	if err := install(r, agentDir); err == nil {
		t.Fatal("install must fail when an installed version is not the pin")
	}
	substrateRoot(t, subRoot)
	if err := os.Remove(filepath.Join(icRoot, "skills/pi-intercom/SKILL.md")); err != nil {
		t.Fatal(err)
	}
	if err := install(r, agentDir); err == nil {
		t.Fatal("install must fail when a packaged skill file is missing")
	}
	writeFile(t, filepath.Join(icRoot, "skills/pi-intercom/SKILL.md"), "x")
	if err := os.Remove(filepath.Join(subRoot, "src/api/intercom-bridge.js")); err != nil {
		t.Fatal(err)
	}
	if err := install(r, agentDir); err == nil {
		t.Fatal("install must fail when a declared runtime file is missing")
	}
}
