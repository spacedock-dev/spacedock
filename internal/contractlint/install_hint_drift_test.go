// ABOUTME: No-drift check between the FO version-gate's OS-aware install hint
// ABOUTME: and docs/site/get-started/install.md: the curl|sh command and the
// ABOUTME: Homebrew tap+formula tokens in the prose must equal install.md's
// ABOUTME: documented forms, so doc drift cannot strand the gate on a stale hint.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// installMDSection returns the lines of install.md's named mkdocs tab section
// (`=== "<header>"`) up to the next tab header or top-level heading.
func installMDSection(t *testing.T, lines []string, header string) []string {
	t.Helper()
	start := -1
	for i, l := range lines {
		if strings.HasPrefix(l, `=== "`+header+`"`) {
			start = i + 1
			break
		}
	}
	if start < 0 {
		t.Fatalf("install.md has no %q tab — the documented install surface moved", header)
	}
	end := len(lines)
	for i := start; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], `=== "`) || strings.HasPrefix(lines[i], "## ") {
			end = i
			break
		}
	}
	return lines[start:end]
}

// TestInstallHintNoDrift asserts token equality (not formatting equality)
// between the FO install reference's hint and install.md. The commands live in
// references/fo-install.md, not the boot-resident core: the core defers the
// whole binary-absent arm to that file, so it is the prose that must not drift.
//
//  1. The curl|sh token in the FO prose equals install.md's documented Linux
//     command — extraction: the single line inside the
//     `=== "Binary (macOS / Linux)"` fence matching `^curl ` and carrying no
//     channel selector, trimmed of the tab's 4-space indentation. The tab also
//     documents the `SPACEDOCK_CHANNEL=edge` variant; the hint tracks the
//     DEFAULT command, so that line is excluded and the default stays unique.
//  2. The FO prose's brew-install hint refers to the same tap and formula as
//     install.md's Homebrew tab (`spacedock-dev/tap` + `spacedock`);
//     both the two-line `brew tap` + `brew install` form and the one-line
//     full-token form are checked.
//
// The brew UPGRADE form in internal/contract is deliberately excluded — the
// upgrade task (fo-boot-upgrade-hint-latest-release) owns that text.
func TestInstallHintNoDrift(t *testing.T) {
	root := repoRoot(t)
	rawInstall, err := os.ReadFile(filepath.Join(root, "docs", "site", "get-started", "install.md"))
	if err != nil {
		t.Fatalf("read install.md: %v", err)
	}
	lines := strings.Split(string(rawInstall), "\n")
	prose := readSkillFile(t, installRefPath)

	// Arm 1: curl|sh token equality.
	var curlLines []string
	for _, l := range installMDSection(t, lines, "Binary (macOS / Linux)") {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "curl ") && !strings.Contains(trimmed, "SPACEDOCK_CHANNEL=") {
			curlLines = append(curlLines, trimmed)
		}
	}
	if len(curlLines) != 1 {
		t.Fatalf("install.md Binary tab carries %d default (channel-free) `^curl ` lines, want exactly 1", len(curlLines))
	}
	documentedCurl := curlLines[0]
	if !strings.Contains(prose, documentedCurl) {
		t.Fatalf("FO install-reference prose's curl|sh hint drifts from install.md's documented command %q", documentedCurl)
	}

	// Arm 2: Homebrew tap+formula token equality.
	tap, formula := brewTokens(installMDSection(t, lines, "macOS (Homebrew)"))
	if tap == "" || formula == "" {
		t.Fatalf("install.md Homebrew tab lacks the two-line brew tap/install form (tap=%q formula=%q)", tap, formula)
	}
	if !strings.Contains(prose, "brew tap "+tap) || !strings.Contains(prose, "brew install "+formula) {
		t.Fatalf("FO install-reference prose's brew-install hint drifts from install.md (tap=%q formula=%q)", tap, formula)
	}
	if oneLine := "brew install " + tap + "/" + formula; !strings.Contains(prose, oneLine) {
		t.Fatalf("FO install-reference prose must also carry the one-line form %q", oneLine)
	}
}

// TestInstallHintProductReadmeNoDrift binds the Homebrew commands of the product
// README to the `macOS (Homebrew)` tab of install.md. It uses the same token
// equality as arm 2 above, with the README in place of the FO prose.
//
// The README holds the last copy of the install commands that no test protects.
// mkdocs.yml does not include the README, so the docs strict build never reads it.
// The two files diverged before this test: the README named
// `spacedock-dev/homebrew-tap` and install.md named `spacedock-dev/tap`. Homebrew
// removes the `homebrew-` prefix, so both names resolved and no reader found the
// difference.
//
// This test compares two independent files that can diverge. It asserts nothing
// about the prose of the README.
func TestInstallHintProductReadmeNoDrift(t *testing.T) {
	root := repoRoot(t)
	rawInstall, err := os.ReadFile(filepath.Join(root, "docs", "site", "get-started", "install.md"))
	if err != nil {
		t.Fatalf("read install.md: %v", err)
	}
	tap, formula := brewTokens(installMDSection(t, strings.Split(string(rawInstall), "\n"), "macOS (Homebrew)"))
	if tap == "" || formula == "" {
		t.Fatalf("install.md Homebrew tab lacks the two-line brew tap/install form (tap=%q formula=%q)", tap, formula)
	}

	rawReadme, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	readmeInstall := markdownSectionFromText(t, string(rawReadme), "## Install")
	readmeTap, readmeFormula := brewTokens(strings.Split(readmeInstall, "\n"))
	if readmeTap != tap {
		t.Errorf("README.md's Install section names tap %q; install.md's Homebrew tab names %q", readmeTap, tap)
	}
	if readmeFormula != formula {
		t.Errorf("README.md's Install section names formula %q; install.md's Homebrew tab names %q", readmeFormula, formula)
	}
}

// brewTokens returns the tap and the formula from the last `brew tap` line and the
// last `brew install` line of a section. It removes the indentation and it ignores
// comment lines. The Homebrew tab gives the edge cask as a comment line that starts
// with `# brew install`. An empty string shows that the section has no such command.
func brewTokens(section []string) (tap, formula string) {
	for _, l := range section {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "brew tap ") {
			tap = strings.TrimPrefix(trimmed, "brew tap ")
		}
		if strings.HasPrefix(trimmed, "brew install ") {
			formula = strings.TrimPrefix(trimmed, "brew install ")
		}
	}
	return tap, formula
}
