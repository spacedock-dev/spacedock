// ABOUTME: Extension repositories as Pi packages: the listing's ref, installed beside
// ABOUTME: Spacedock's own package, and read back from Pi's settings.json.
package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type extensionListing struct {
	Plugins []struct {
		Name   string `json:"name"`
		Source struct {
			Source string `json:"source"`
			URL    string `json:"url"`
			Ref    string `json:"ref"`
		} `json:"source"`
	} `json:"plugins"`
}

// piExtensionSource returns extension name's Pi package source pinned at the listing's
// ref, and the same source without it, which is what `pi remove` matches.
func piExtensionSource(listing []byte, name string) (pinned, bare string, err error) {
	var l extensionListing
	if err := json.Unmarshal(listing, &l); err != nil {
		return "", "", fmt.Errorf("the extension listing is not valid JSON: %w", err)
	}
	for _, p := range l.Plugins {
		if p.Name != name {
			continue
		}
		host, repo := parseGitPackagePath(p.Source.URL)
		if p.Source.Source != "url" || p.Source.Ref == "" || host == "" {
			return "", "", fmt.Errorf("the extension listing does not give %s as a git url with a ref", name)
		}
		bare = "git:" + host + "/" + repo
		return bare + "@" + p.Source.Ref, bare, nil
	}
	return "", "", fmt.Errorf("the extension listing has no %s entry", name)
}

// installPiExtensions installs each extension that has a Pi package. A failure prints
// a `spacedock install:` line and leaves the exit code alone, as for a Claude Code
// companion, so `update` can show it.
func installPiExtensions(piOps piRuntimeOps, agentDir string, stdout io.Writer) {
	var wanted []extension
	for _, ext := range extensions {
		if ext.piTitle != "" {
			wanted = append(wanted, ext)
		}
	}
	if len(wanted) == 0 {
		return
	}
	listingSource := extensionMarketplaceSource()
	listing, err := piOps.ReadExtensionListing(listingSource)
	if err != nil {
		fmt.Fprintf(stdout, "spacedock install: could not read the extension listing %s, so the Pi packages were not installed: %v\n", listingSource, err)
		return
	}
	for _, ext := range wanted {
		pinned, bare, err := piExtensionSource(listing, ext.name)
		if err != nil {
			fmt.Fprintf(stdout, "spacedock install: %v, so the %s Pi package was not installed\n", err, ext.name)
			continue
		}
		if piSettingsHold(agentDir, pinned) {
			continue
		}
		piOps.PiRemove(bare)
		out, err := piOps.PiInstall(pinned)
		if strings.TrimSpace(out) != "" {
			fmt.Fprint(stdout, out)
		}
		if err != nil {
			fmt.Fprintf(stdout, "spacedock install: pi install %q failed: %v\n", pinned, err)
		}
	}
}

func (execPiRuntimeOps) ReadExtensionListing(source string) ([]byte, error) {
	if info, err := os.Stat(source); err == nil && info.IsDir() {
		return os.ReadFile(filepath.Join(source, ".claude-plugin", "marketplace.json"))
	}
	if strings.Count(source, "/") != 1 {
		return nil, fmt.Errorf("%q is neither a directory nor owner/repo", source)
	}
	client := http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get("https://raw.githubusercontent.com/" + source + "/HEAD/.claude-plugin/marketplace.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub answered %s", resp.Status)
	}
	return io.ReadAll(io.LimitReader(resp.Body, 1<<20))
}

type piExtensionPackage struct {
	title, where string
	ok           bool
}

// piSettingsHold reports whether settings.json already holds the pinned source. That
// is evidence only for a tag ref: a branch ref keeps its name while it moves, so
// Spacedock's own package, pinned to a channel branch, is always replaced.
func piSettingsHold(agentDir, pinned string) bool {
	for _, src := range piSettingsSources(agentDir) {
		if strings.TrimSpace(src) == pinned {
			return true
		}
	}
	return false
}

func piSettingsSources(agentDir string) []string {
	data, err := os.ReadFile(filepath.Join(agentDir, "settings.json"))
	if err != nil {
		return nil
	}
	var settings struct {
		Packages []json.RawMessage `json:"packages"`
	}
	if json.Unmarshal(data, &settings) != nil {
		return nil
	}
	var sources []string
	for _, raw := range settings.Packages {
		sources = append(sources, piPackageSourceFromEntry(raw))
	}
	return sources
}

// piRegisteredPackage finds the git package of repository repoName in settings.json by
// its repository name: the root of a package without a package.json has no name to
// read, and doctor cannot ask the listing.
func piRegisteredPackage(agentDir, home, repoName string) (root, ref string, ok bool) {
	for _, src := range piSettingsSources(agentDir) {
		spec, isGit := strings.CutPrefix(strings.TrimSpace(src), "git:")
		if !isGit {
			continue
		}
		if _, repoPath := parseGitPackagePath(spec); path.Base(repoPath) == repoName {
			if i := strings.LastIndex(spec, "@"); i > strings.Index(spec, "/") {
				ref = spec[i+1:]
			}
			return resolveSettingsPackageRoot(src, agentDir, home), ref, true
		}
	}
	return "", "", false
}

func piExtensionPackages(agentDir, home string) []piExtensionPackage {
	var packages []piExtensionPackage
	for _, ext := range extensions {
		if ext.piTitle == "" {
			continue
		}
		root, ref, ok := piRegisteredPackage(agentDir, home, ext.name)
		where := ""
		if ok {
			where = root + " (" + refOrUnpinned(ref) + ")"
		}
		packages = append(packages, piExtensionPackage{title: ext.piTitle, where: where, ok: ok})
	}
	return packages
}

func refOrUnpinned(ref string) string {
	if ref == "" {
		return "unpinned"
	}
	return ref
}
