package cli

import (
	"fmt"
	"regexp"
	"strings"
)

const compassRecipeTemplate = `class Compass < Formula
  desc "Agnostic agent cartography and visualization tool"
  homepage "https://github.com/spacedock-dev/compass"
  version "{version}"
  license all_of: ["Apache-2.0", "OFL-1.1"]

  on_macos do
    on_arm do
      url "https://github.com/spacedock-dev/compass/releases/download/v{version}/compass_{version}_darwin_arm64.tar.gz"
      sha256 "{sha256}"
    end
    on_intel do
      url "https://github.com/spacedock-dev/compass/releases/download/v{version}/compass_{version}_darwin_amd64.tar.gz"
      sha256 "{sha256}"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/spacedock-dev/compass/releases/download/v{version}/compass_{version}_linux_arm64.tar.gz"
      sha256 "{sha256}"
    end
    on_intel do
      url "https://github.com/spacedock-dev/compass/releases/download/v{version}/compass_{version}_linux_amd64.tar.gz"
      sha256 "{sha256}"
    end
  end

  def install
    libexec.install Dir["*"]
    bin.install_symlink libexec/"compass" => "compass"
    bin.install_symlink libexec/"compass" => "spacedock-compass"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/compass --version")
  end
end
`

var compassVersionLine = regexp.MustCompile(`(?m)^  version "([^"]+)"$`)
var platformSHA256 = regexp.MustCompile(`(?m)^      sha256 "([0-9a-f]{64})"$`)

func compassRecipe(data string) (string, map[string]string, error) {
	version := compassVersionLine.FindStringSubmatch(data)
	hashes := platformSHA256.FindAllStringSubmatch(data, -1)
	platforms := []string{"darwin_arm64", "darwin_amd64", "linux_arm64", "linux_amd64"}
	if len(version) != 2 || !stableRepairVersion.MatchString(version[1]) || len(hashes) != len(platforms) {
		return "", nil, fmt.Errorf("the Compass formula is not the supported stable binary recipe")
	}
	expected := strings.ReplaceAll(compassRecipeTemplate, "{version}", version[1])
	assets := make(map[string]string, len(platforms))
	for i, platform := range platforms {
		assets[platform] = hashes[i][1]
		expected = strings.Replace(expected, "{sha256}", hashes[i][1], 1)
	}
	if data != expected {
		return "", nil, fmt.Errorf("the Compass formula changed its install recipe. Verify it before upgrading")
	}
	return version[1], assets, nil
}

func compassAsset(version, goos, goarch string) string {
	return fmt.Sprintf("https://github.com/spacedock-dev/compass/releases/download/v%s/compass_%s_%s_%s.tar.gz", version, version, goos, goarch)
}

func (r brewRepair) cleanBrewTap(tap string) (head, branch string, err error) {
	remote, err := r.git(tap, "remote", "get-url", "--all", "origin")
	if err != nil || (strings.TrimSpace(remote) != spacedockTapRemote && strings.TrimSpace(remote) != spacedockTapRemote+".git") {
		return "", "", fmt.Errorf("the installed tap does not track the official Spacedock source. Verify its origin")
	}
	actualRoot, err := r.git(tap, "rev-parse", "--show-toplevel")
	if err != nil || strings.TrimSpace(actualRoot) != tap {
		return "", "", fmt.Errorf("Homebrew's Spacedock tap is not its own Git checkout")
	}
	dirty, err := r.git(tap, "status", "--porcelain", "--untracked-files=all")
	if err != nil || strings.TrimSpace(dirty) != "" {
		return "", "", fmt.Errorf("the Spacedock tap is dirty or unreadable. Finish the local tap work first")
	}
	head, err = r.git(tap, "rev-parse", "HEAD")
	if err != nil || !repairCommit.MatchString(strings.TrimSpace(head)) {
		return "", "", fmt.Errorf("the Spacedock tap has no readable Git commit")
	}
	branch, err = r.git(tap, "symbolic-ref", "--short", "HEAD")
	if err != nil {
		return "", "", fmt.Errorf("the Spacedock tap is detached. Switch to its tracked default branch")
	}
	branch = strings.TrimSpace(branch)
	upstream, upErr := r.git(tap, "rev-parse", "--symbolic-full-name", "@{upstream}")
	defaultRef, defaultErr := r.git(tap, "symbolic-ref", "refs/remotes/origin/HEAD")
	ref := "refs/remotes/origin/" + branch
	if upErr != nil || defaultErr != nil || strings.TrimSpace(upstream) != ref || strings.TrimSpace(defaultRef) != ref {
		return "", "", fmt.Errorf("the Spacedock tap is not on its tracked default branch")
	}
	if _, err := r.git(tap, "merge-base", "--is-ancestor", strings.TrimSpace(head), ref); err != nil {
		return "", "", fmt.Errorf("the Spacedock tap is ahead of its tracked branch or has diverged. Reconcile it first")
	}
	return strings.TrimSpace(head), branch, nil
}

func (r brewRepair) advanceBrewTap(tap, branch, head, recipe, target string) (string, string, error) {
	refs, err := r.git(tap, "ls-remote", "--symref", spacedockTapRemote, "HEAD")
	lines := strings.Split(strings.TrimSpace(refs), "\n")
	if err != nil || len(lines) != 2 || lines[0] != "ref: refs/heads/"+branch+"\tHEAD" {
		return "", "", fmt.Errorf("could not verify the official tap default branch")
	}
	candidate, ok := strings.CutSuffix(lines[1], "\tHEAD")
	if !ok || !repairCommit.MatchString(candidate) {
		return "", "", fmt.Errorf("could not identify the official Spacedock tap commit")
	}
	tracking := "refs/remotes/origin/" + branch
	if _, err := r.git(tap, "fetch", "--no-tags", "--no-recurse-submodules", "--no-write-fetch-head",
		spacedockTapRemote, "refs/heads/"+branch+":"+tracking); err != nil {
		return "", "", err
	}
	if _, err := r.git(tap, "merge-base", "--is-ancestor", candidate, tracking); err != nil {
		return "", "", fmt.Errorf("the official tap default branch moved during the fetch")
	}
	if _, err := r.git(tap, "merge-base", "--is-ancestor", head, candidate); err != nil {
		return "", "", fmt.Errorf("the Spacedock tap has diverged. Reconcile its tracked branch first")
	}
	data, err := r.git(tap, "show", candidate+":"+recipe)
	if err != nil {
		return "", "", err
	}
	version, _, err := compassRecipe(data)
	if err != nil {
		return "", "", err
	}
	if version != target {
		return "", "", fmt.Errorf("the official tap offers Compass %s. Retry after it offers %s", version, target)
	}
	now, currentBranch, err := r.cleanBrewTap(tap)
	if err != nil || now != head || currentBranch != branch {
		return "", "", fmt.Errorf("the Spacedock tap changed during preflight")
	}
	if _, err := r.git(tap, "merge", "--ff-only", "--no-edit", candidate); err != nil {
		return "", "", err
	}
	return candidate, data, nil
}
