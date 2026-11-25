package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type brewHostPiece struct {
	host, title string
	before      pluginListEntry
	manifest    string
}

func (r brewRepair) readHostPiece(host, id string) (pluginListEntry, error) {
	out, err := r.run(host, "plugin", "list", "--json")
	if err != nil {
		return pluginListEntry{}, err
	}
	entries, err := parseHostPluginList(host, []byte(out))
	if err != nil {
		return pluginListEntry{}, err
	}
	var piece pluginListEntry
	for _, entry := range entries {
		if !entry.Installed || entry.ID != id {
			continue
		}
		if piece.Installed {
			return pluginListEntry{}, fmt.Errorf("%s has multiple installed scopes for %s", host, id)
		}
		piece = entry
	}
	return piece, nil
}

func hostPieceVersion(piece pluginListEntry, err error) string {
	switch {
	case err != nil:
		return "unknown"
	case !piece.Installed:
		return notInstalled
	case piece.Version == "":
		return "unknown"
	default:
		return piece.Version
	}
}

type repairMarketplace struct {
	Name, Root, Source, Repo, Ref, URL, Path, InstallLocation string
	MarketplaceSource                                         struct{ SourceType, Source, Ref string }
}

func (r brewRepair) hostMarketplace(host, source string) (repairMarketplace, error) {
	selected, err := r.marketplaceByName(host, extensionMarketplace)
	if err != nil {
		return repairMarketplace{}, err
	}
	if selected.Root == "" || !selected.matches(source) {
		return repairMarketplace{}, fmt.Errorf("%s has no matching %s registration. Correct that host's marketplace source", host, extensionMarketplace)
	}
	return selected, nil
}

func (r brewRepair) marketplaceByName(host, name string) (repairMarketplace, error) {
	out, err := r.run(host, "plugin", "marketplace", "list", "--json")
	if err != nil {
		return repairMarketplace{}, err
	}
	var markets []repairMarketplace
	if host == "claude" {
		err = json.Unmarshal([]byte(out), &markets)
	} else {
		var envelope struct{ Marketplaces []repairMarketplace }
		err = json.Unmarshal([]byte(out), &envelope)
		markets = envelope.Marketplaces
	}
	if err != nil || markets == nil {
		return repairMarketplace{}, fmt.Errorf("could not read the configured %s marketplaces", host)
	}
	var selected repairMarketplace
	for _, market := range markets {
		if market.Name != name {
			continue
		}
		if selected.Name != "" {
			return repairMarketplace{}, fmt.Errorf("%s lists multiple %s marketplaces", host, name)
		}
		selected = market
	}
	if host == "claude" {
		selected.Root = selected.InstallLocation
	}
	return selected, nil
}

func marketplaceListingPaths(host string) []string {
	const fallback = ".claude-plugin/marketplace.json"
	if host == "codex" {
		return []string{".agents/plugins/marketplace.json", fallback}
	}
	return []string{fallback}
}

func (m repairMarketplace) readListing(host string) ([]byte, error) {
	for _, file := range marketplaceListingPaths(host) {
		data, err := readRepairFile(m.Root, file)
		if !errors.Is(err, os.ErrNotExist) {
			return data, err
		}
	}
	return nil, os.ErrNotExist
}

func (m repairMarketplace) matches(source string) bool {
	switch {
	case m.Source == "directory" || m.MarketplaceSource.SourceType == "local":
		expected, err := filepath.EvalSymlinks(source)
		root, rootErr := filepath.EvalSymlinks(m.Root)
		return err == nil && rootErr == nil && filepath.IsAbs(expected) && expected == root
	case m.Source == "github":
		return m.Repo == source && m.Ref == ""
	case m.Source == "git":
		return repairMarketplaceRepo(m.URL) == source && m.Ref == ""
	case m.MarketplaceSource.SourceType == "git":
		return repairMarketplaceRepo(m.MarketplaceSource.Source) == source && m.MarketplaceSource.Ref == ""
	default:
		return false
	}
}

func repairMarketplaceRepo(remote string) string {
	return strings.TrimSuffix(strings.TrimPrefix(remote, "https://github.com/"), ".git")
}

func (r brewRepair) checkHostSource(host, name, source string, entry repairEntry) (string, error) {
	market, err := r.hostMarketplace(host, source)
	if err != nil {
		return "", err
	}
	data, err := market.readListing(host)
	if err != nil {
		return "", fmt.Errorf("could not read the installed %s marketplace snapshot: %w", host, err)
	}
	actual, err := resolveRepairEntry(data, name)
	if err != nil || actual != entry {
		return "", fmt.Errorf("%s's %s snapshot cannot select %s %s. Make the matching listing available in this host", host, extensionMarketplace, name, entry.Version)
	}
	if entry.Local != "" && !filepath.IsAbs(source) {
		return "", fmt.Errorf("%s resolves %s from a local path in a remote marketplace", host, name)
	}
	manifest, err := r.verifyRepairManifest(entry, market.Root, host, name)
	if err != nil {
		return "", err
	}
	return market.Root + ":" + manifest, nil
}

func (r brewRepair) applyHostRepair(p brewHostPiece, name, source string, entry repairEntry) error {
	id := extensionPluginID(name)
	now, err := r.readHostPiece(p.host, id)
	if err != nil {
		return err
	}
	if now != p.before {
		return fmt.Errorf("%s inventory changed while the repair was being prepared", p.title)
	}
	if hostPieceVersion(now, nil) == entry.Version {
		return nil
	}
	manifest, err := r.checkHostSource(p.host, name, source, entry)
	if err != nil {
		return err
	}
	if manifest != p.manifest {
		return fmt.Errorf("%s source changed while the repair was being prepared", p.title)
	}
	if err := r.changeHostPiece(p.host, id, now); err != nil {
		return err
	}
	after, err := r.readHostPiece(p.host, id)
	if err != nil {
		return err
	}
	if hostPieceVersion(after, nil) != entry.Version || now.Installed && (after.Scope != now.Scope || after.Enabled != now.Enabled) {
		return fmt.Errorf("%s did not report %s %s with its original scope and enabled state", p.title, name, entry.Version)
	}
	return nil
}

func (r brewRepair) changeHostPiece(host, id string, before pluginListEntry) error {
	verb := "add"
	var scope []string
	if host == "claude" {
		verb = "install"
		if before.Installed {
			verb = "update"
		}
		scope = []string{"--scope", "user"}
	}
	_, err := r.run(host, append([]string{"plugin", verb, id}, scope...)...)
	return err
}

func (r brewRepair) refreshCorePlugin(host string) error {
	id := channelPluginID(devBranch)
	piece, err := r.readHostPiece(host, id)
	if err != nil {
		return err
	}
	if !piece.Installed || host == "claude" && piece.Scope != "user" || !piece.Enabled {
		return nil
	}
	if r.localCoreIsCurrent(host, piece.Version) {
		return nil
	}
	return r.changeHostPiece(host, id, piece)
}

func (r brewRepair) localCoreIsCurrent(host, installedVersion string) bool {
	market, err := r.marketplaceByName(host, channelMarketplace(devBranch))
	if err != nil || market.Root == "" || (market.Source != "directory" && market.MarketplaceSource.SourceType != "local") {
		return false
	}
	data, err := market.readListing(host)
	var listing struct {
		Plugins []struct {
			Name   string
			Source json.RawMessage
		}
	}
	if err != nil || json.Unmarshal(data, &listing) != nil {
		return false
	}
	name := channelEntry(devBranch)
	for _, plugin := range listing.Plugins {
		if plugin.Name != name {
			continue
		}
		var source string
		if json.Unmarshal(plugin.Source, &source) != nil || !strings.HasPrefix(source, "./") {
			return false
		}
		dir := strings.TrimPrefix(source, "./")
		if !safeRepairSourcePath(dir) {
			return false
		}
		data, err := readRepairFile(market.Root, filepath.Join(dir, manifestSubpath(host)))
		if host == "codex" && errors.Is(err, os.ErrNotExist) {
			data, err = readRepairFile(market.Root, filepath.Join(dir, manifestSubpath("claude")))
		}
		var manifest struct{ Name, Version string }
		return err == nil && json.Unmarshal(data, &manifest) == nil &&
			manifest.Name == name && manifest.Version != "" && manifest.Version == installedVersion
	}
	return false
}

func (r brewRepair) refreshExtensionPlugin(host, name string) error {
	id := extensionPluginID(name)
	piece, err := r.readHostPiece(host, id)
	if err != nil {
		return err
	}
	if piece.Installed && (host == "claude" && piece.Scope != "user" || host == "codex" && !piece.Enabled) {
		return nil
	}
	if name == "subspace" && piece.Installed {
		if handled, err := r.refreshRegisteredSubspace(host, piece); handled {
			return err
		}
	}
	return r.changeHostPiece(host, id, piece)
}
