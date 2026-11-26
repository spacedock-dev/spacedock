package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func (r brewRepair) refreshRegisteredSubspace(host string, before pluginListEntry) (bool, error) {
	source := r.releaseSource()
	if !repairGitHubRepo.MatchString(source) {
		return false, nil
	}
	market, err := r.marketplaceByName(host, extensionMarketplace)
	if err != nil || !market.matches(source) || !market.isGit() {
		return false, nil
	}
	const name = "subspace"
	listing, err := market.readListing(host)
	if err != nil {
		return true, fmt.Errorf("could not read the %s marketplace snapshot: %w", host, err)
	}
	cached, err := resolveRepairEntry(listing, name)
	if err != nil || !isTaggedRepairEntry(cached) {
		return false, nil
	}
	target, err := r.refreshRegisteredSnapshot(host, source, market, name, cached, before.Version)
	if err != nil {
		return true, err
	}
	now, err := r.readHostPiece(host, extensionPluginID(name))
	if err != nil {
		return true, err
	}
	if now != before {
		return true, fmt.Errorf("%s's installed %s plugin changed during the marketplace refresh", host, name)
	}
	if now.Version == target.Version {
		return true, nil
	}
	if err := r.changeHostPiece(host, now.ID, now); err != nil {
		return true, err
	}
	after, err := r.readHostPiece(host, now.ID)
	if err != nil {
		return true, err
	}
	if hostPieceVersion(after, nil) != target.Version || after.Scope != before.Scope || after.Enabled != before.Enabled {
		return true, fmt.Errorf("%s did not report %s at %s with its original scope and enabled state", host, now.ID, target.Version)
	}
	return true, nil
}

func isTaggedRepairEntry(entry repairEntry) bool {
	return entry.Repo != "" && (entry.Ref == entry.Version || entry.Ref == "v"+entry.Version)
}

func (r brewRepair) refreshRegisteredSnapshot(host, source string, market repairMarketplace, name string, cached repairEntry, installed string) (repairEntry, error) {
	target := cached
	err := r.fastForwardRegisteredSnapshot(host, source, market, func(published []byte) (bool, error) {
		if published != nil {
			next, err := resolveRepairEntry(published, name)
			if err != nil || !isTaggedRepairEntry(next) || next.Repo != cached.Repo || next.Path != cached.Path {
				return false, fmt.Errorf("the published listing does not give %s a matching release tag from its registered plugin source", name)
			}
			if err := repairCanAdvance("the marketplace snapshot", cached.Version, next.Version); err != nil {
				return false, err
			}
			target = next
		}
		if err := repairCanAdvance("the installed plugin", installed, target.Version); err != nil {
			return false, err
		}
		return cached != target, nil
	})
	if err != nil {
		return repairEntry{}, err
	}
	data, err := market.readListing(host)
	if err != nil {
		return repairEntry{}, err
	}
	actual, err := resolveRepairEntry(data, name)
	if err != nil || actual != target {
		return repairEntry{}, fmt.Errorf("%s's marketplace snapshot did not select %s %s", host, name, target.Version)
	}
	return target, nil
}

// syncHostSnapshot fast-forwards a host's registered Git snapshot of the extension
// marketplace when it cannot yet select entry, the release the published listing names.
// Only the catalog moves, so installed plugins keep their versions until the repair
// updates the one it selected. A directory registration is left to checkHostSource.
func (r brewRepair) syncHostSnapshot(host, name, source string, entry repairEntry) error {
	if !repairGitHubRepo.MatchString(source) {
		return nil
	}
	market, err := r.hostMarketplace(host, source)
	if err != nil || !market.isGit() {
		return nil
	}
	if data, err := market.readListing(host); err == nil {
		if actual, err := resolveRepairEntry(data, name); err == nil && actual == entry {
			return nil
		}
	}
	err = r.fastForwardRegisteredSnapshot(host, source, market, func(published []byte) (bool, error) {
		if published == nil {
			return false, nil
		}
		if next, err := resolveRepairEntry(published, name); err != nil || next != entry {
			return false, fmt.Errorf("the published listing does not select %s %s yet. Retry after it does", name, entry.Version)
		}
		return true, nil
	})
	if err != nil {
		return fmt.Errorf("could not refresh %s's %s snapshot: %w", host, extensionMarketplace, err)
	}
	return nil
}

// fastForwardRegisteredSnapshot moves a host's managed Git checkout of the extension
// marketplace to the published HEAD. decide receives the listing at that commit, or nil
// when the checkout is already there, and reports whether to merge. Git refreshes only
// the catalog; the caller runs any single-plugin host command after its own readback.
func (r brewRepair) fastForwardRegisteredSnapshot(host, source string, market repairMarketplace, decide func(published []byte) (bool, error)) error {
	root, err := r.managedMarketplaceRoot(host, market)
	if err != nil {
		return err
	}
	for _, key := range []string{"GIT_DIR", "GIT_COMMON_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE", "GIT_OBJECT_DIRECTORY"} {
		r.env = withoutEnv(r.env, key)
	}
	origin := "https://github.com/" + source
	remote, err := r.git(root, "config", "--local", "--get-all", "remote.origin.url")
	registrations := []string{origin, origin + ".git", "git@github.com:" + source + ".git", "ssh://git@github.com/" + source + ".git"}
	if err != nil || !slices.Contains(registrations, strings.TrimSpace(remote)) {
		return fmt.Errorf("%s's marketplace checkout does not track its registered GitHub source", host)
	}
	head, err := r.cleanMarketplaceHead(root)
	if err != nil {
		return err
	}
	ref, err := r.git(root, "ls-remote", "--exit-code", origin+".git", "HEAD")
	fields := strings.Fields(ref)
	if err != nil || len(fields) != 2 || fields[1] != "HEAD" || !repairCommit.MatchString(fields[0]) {
		return fmt.Errorf("could not resolve the published %s marketplace HEAD for %s", extensionMarketplace, host)
	}
	tip := fields[0]
	var published []byte
	if tip != head {
		if _, err := r.git(root, "fetch", "--no-tags", "--no-recurse-submodules", "--no-write-fetch-head", origin+".git", tip); err != nil {
			return err
		}
		if _, err := r.git(root, "merge-base", "--is-ancestor", head, tip); err != nil {
			return fmt.Errorf("%s's marketplace snapshot cannot fast-forward to the published listing", host)
		}
		if published, err = r.marketplaceListingAt(host, root, tip); err != nil {
			return err
		}
	}
	merge, err := decide(published)
	if err != nil {
		return err
	}
	if merge && tip != head {
		if current, err := r.cleanMarketplaceHead(root); err != nil || current != head {
			return fmt.Errorf("%s's marketplace snapshot changed while preparing its refresh", host)
		}
		if _, err := r.git(root, "-c", "submodule.recurse=false", "merge", "--ff-only", tip); err != nil {
			return err
		}
	}
	selected, err := r.marketplaceByName(host, extensionMarketplace)
	if err != nil || selected != market {
		return fmt.Errorf("%s's marketplace registration changed during the refresh", host)
	}
	return nil
}

func (r brewRepair) managedMarketplaceRoot(host string, market repairMarketplace) (string, error) {
	env := envMap(r.env)
	profile, relative := env["CLAUDE_CONFIG_DIR"], filepath.Join("plugins", "marketplaces")
	defaultDir := ".claude"
	if host == "codex" {
		profile, relative, defaultDir = env["CODEX_HOME"], filepath.Join(".tmp", "marketplaces"), ".codex"
	}
	if profile == "" && env["HOME"] != "" {
		profile = filepath.Join(env["HOME"], defaultDir)
	}
	parent, pathErr := filepath.EvalSymlinks(filepath.Join(profile, relative))
	root, rootErr := filepath.EvalSymlinks(market.Root)
	if !filepath.IsAbs(profile) || pathErr != nil || rootErr != nil || root != filepath.Join(parent, extensionMarketplace) {
		return "", fmt.Errorf("%s's marketplace is outside its supported managed Git snapshot", host)
	}
	gitDir, err := os.Lstat(filepath.Join(root, ".git"))
	if err != nil || !gitDir.IsDir() {
		return "", fmt.Errorf("%s's managed marketplace snapshot is not a Git checkout", host)
	}
	return root, nil
}

func (r brewRepair) cleanMarketplaceHead(root string) (string, error) {
	dir, err := r.git(root, "rev-parse", "--show-toplevel")
	if err != nil || strings.TrimSpace(dir) != root {
		return "", fmt.Errorf("marketplace snapshot is not its own Git checkout")
	}
	status, err := r.git(root, "status", "--porcelain", "--untracked-files=all")
	if err != nil || strings.TrimSpace(status) != "" {
		return "", fmt.Errorf("marketplace snapshot has local changes or cannot be read")
	}
	head, err := r.git(root, "rev-parse", "HEAD")
	if err != nil || !repairCommit.MatchString(strings.TrimSpace(head)) {
		return "", fmt.Errorf("could not identify the marketplace snapshot HEAD")
	}
	return strings.TrimSpace(head), nil
}

func (r brewRepair) marketplaceListingAt(host, root, revision string) ([]byte, error) {
	for _, file := range marketplaceListingPaths(host) {
		item, err := r.git(root, "ls-tree", revision, "--", file)
		if err != nil {
			return nil, err
		}
		if item == "" {
			continue
		}
		if !strings.HasPrefix(item, "100644 blob ") && !strings.HasPrefix(item, "100755 blob ") {
			return nil, fmt.Errorf("published marketplace listing is not a regular file")
		}
		data, err := r.git(root, "show", revision+":"+file)
		return []byte(data), err
	}
	return nil, fmt.Errorf("published marketplace listing is missing")
}
