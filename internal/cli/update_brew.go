package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spacedock-dev/spacedock/internal/release"
)

type brewRepair struct {
	ctx context.Context
	env []string
	dir string
}

func (r brewRepair) run(name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 3*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env, cmd.Dir = r.env, r.dir
	cmd.WaitDelay = 2 * time.Second
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = strings.TrimSpace(string(out))
		}
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		return string(out), fmt.Errorf("%s %s: %w (%s)", name, strings.Join(args, " "), err, detail)
	}
	return string(out), nil
}

func (r brewRepair) brew(args ...string) (string, error) {
	r.env = append(append([]string(nil), r.env...),
		"HOMEBREW_NO_AUTO_UPDATE=1", "HOMEBREW_NO_INSTALLED_DEPENDENTS_CHECK=1",
		"HOMEBREW_NO_INSTALL_CLEANUP=1", "HOMEBREW_NO_ANALYTICS=1", "HOMEBREW_NO_ENV_HINTS=1",
		"HOMEBREW_NO_COLOR=1", "LC_ALL=C")
	return r.run("brew", args...)
}

func (r brewRepair) git(dir string, args ...string) (string, error) {
	if dir != "" {
		r.dir = dir
	}
	r.env = append(append([]string(nil), r.env...), "GIT_TERMINAL_PROMPT=0")
	return r.run("git", append([]string{"--no-optional-locks", "-c", "core.hooksPath=/dev/null", "-c", "credential.helper="}, args...)...)
}

func (r brewRepair) rawGitHub(repo, ref, file string) ([]byte, error) {
	address := (&url.URL{Scheme: "https", Host: "raw.githubusercontent.com", Path: "/" + repo + "/" + ref + "/" + file}).String()
	req, err := http.NewRequestWithContext(r.ctx, http.MethodGet, address, nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, os.ErrNotExist
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub returned HTTP %d", resp.StatusCode)
	}
	const maxListingBytes = 1 << 20
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxListingBytes+1))
	if len(data) > maxListingBytes {
		return nil, fmt.Errorf("the GitHub document exceeds the supported size")
	}
	return data, err
}

type brewExtensionPlan struct {
	ext       extension
	path      string
	program   string
	entry     repairEntry
	hosts     []brewHostPiece
	formula   *brewFormulaPlan
	preflight error
}

func (r brewRepair) planExtensions(before versionReport) (string, []brewExtensionPlan) {
	source := r.releaseSource()
	if info, err := os.Stat(source); err == nil && info.IsDir() {
		source, _ = filepath.Abs(source)
	}
	var listing []byte
	var listingErr error
	listingRead := false
	var plans []brewExtensionPlan
	for _, ext := range extensions {
		if ext.homebrewFormula != compassFormula {
			continue
		}
		p := brewExtensionPlan{ext: ext}
		for _, row := range before.rows {
			if row.label == ext.name+" program" {
				p.path, p.program = row.where, row.version
			}
		}
		if p.program == notInstalled && before.programsOptedOut {
			continue
		}
		for _, host := range before.hosts {
			if host.err == nil && (!host.listed() || host.lookup(channelPluginID(devBranch)) == notInstalled) {
				continue
			}
			piece, err := r.readHostPiece(host.name, extensionPluginID(ext.name))
			if err != nil && p.preflight == nil {
				p.preflight = fmt.Errorf("%s installation could not be read: %w", host.title, err)
				piece.Version = "unknown"
			}
			p.hosts = append(p.hosts, brewHostPiece{host: host.name, title: host.title, before: piece})
		}
		if len(p.hosts) == 0 {
			continue
		}
		if p.preflight == nil {
			if p.program == notInstalled {
				p.preflight = fmt.Errorf("the program is missing. Install it with: brew install %s", ext.homebrewFormula)
			} else {
				if !listingRead {
					listing, listingErr = r.readReleaseListing(source)
					listingRead = true
				}
				p.preflight = r.prepareExtension(&p, listing, listingErr, source)
			}
		}
		plans = append(plans, p)
	}
	return source, plans
}

func (r brewRepair) prepareExtension(p *brewExtensionPlan, listing []byte, listingErr error, source string) error {
	if listingErr != nil {
		return fmt.Errorf("could not read the stable marketplace listing: %w", listingErr)
	}
	var err error
	p.entry, err = resolveRepairEntry(listing, p.ext.name)
	if err != nil {
		return err
	}
	target := p.entry.Version
	if err := repairCanAdvance("program", p.program, target); err != nil {
		return err
	}
	for i := range p.hosts {
		host := &p.hosts[i]
		piece := host.before
		version := hostPieceVersion(piece, nil)
		if err := repairCanAdvance(host.title+" plugin", version, target); err != nil {
			return err
		}
		edgeCopy, err := r.readHostPiece(host.host, edgeCopyID(p.ext.name))
		if err != nil {
			return err
		}
		if edgeCopy.Installed {
			return fmt.Errorf("%s also holds %s. Resolve that duplicate before repairing %s", host.title, edgeCopy.ID, p.ext.name)
		}
		if version == target && p.program == target {
			continue
		}
		if version != target && piece.Installed && (host.host == "claude" && piece.Scope != "user" || host.host == "codex" && !piece.Enabled) {
			return fmt.Errorf("%s cannot safely update the existing scope or disabled state for %s", host.title, p.ext.name)
		}
		if err := r.syncHostSnapshot(host.host, p.ext.name, source, p.entry); err != nil {
			return err
		}
		host.manifest, err = r.checkHostSource(host.host, p.ext.name, source, p.entry)
		if err != nil {
			return err
		}
	}
	if p.program == target {
		info, err := r.formulaInfo(p.ext.homebrewFormula)
		if err != nil {
			return err
		}
		version, err := r.ownedBrewProgram(p.ext, info, p.path)
		if err != nil {
			return err
		}
		if version != target {
			return fmt.Errorf("the active program changed during preflight")
		}
		return nil
	}
	formula, err := r.prepareBrewFormula(p.ext, p.path, target)
	if err != nil {
		return err
	}
	if p.program != formula.oldVersion {
		return fmt.Errorf("the active program changed while preparing its formula")
	}
	p.formula = &formula
	return nil
}

func repairCanAdvance(piece, installed, target string) error {
	if installed == notInstalled {
		return nil
	}
	if !stableRepairVersion.MatchString(installed) {
		return fmt.Errorf("%s reports %s, which cannot be compared to stable release %s", piece, installed, target)
	}
	if cmp, err := release.ComparePreVersion(installed, target); err != nil || cmp > 0 {
		return fmt.Errorf("%s is at %s. A repair to %s would require a downgrade", piece, installed, target)
	}
	return nil
}

func (r brewRepair) applyExtension(p brewExtensionPlan, source string, progress *updateProgress) error {
	if p.preflight != nil {
		return p.preflight
	}
	for _, h := range p.hosts {
		current, err := r.readHostPiece(h.host, extensionPluginID(p.ext.name))
		if err != nil {
			return err
		}
		if current != h.before {
			return fmt.Errorf("%s installation changed during preflight", h.title)
		}
		if h.manifest == "" {
			continue
		}
		fingerprint, err := r.checkHostSource(h.host, p.ext.name, source, p.entry)
		if err != nil {
			return err
		}
		if fingerprint != h.manifest {
			return fmt.Errorf("%s source changed during preflight", h.title)
		}
	}
	if p.formula != nil {
		end := progress.step("Upgrading the Homebrew " + p.ext.name + " formula")
		err := r.applyBrewFormula(p.ext, p.path, p.entry.Version, *p.formula)
		end()
		if err != nil {
			return err
		}
	}
	for _, h := range p.hosts {
		if h.manifest == "" || hostPieceVersion(h.before, nil) == p.entry.Version {
			continue
		}
		end := progress.step("Repairing the " + p.ext.name + " plugin on " + h.title)
		err := r.applyHostRepair(h, p.ext.name, source, p.entry)
		end()
		if err != nil {
			return fmt.Errorf("%s: %w", h.title, err)
		}
	}

	return nil
}

func (r brewRepair) reportExtension(p brewExtensionPlan, stepErr error, stdout, stderr io.Writer) bool {
	path, version := extensionProgram(p.ext, r.env, envMap(r.env)["PATH"])
	pieces := []string{"program " + version}
	matched := version == p.entry.Version
	kept := p.program == version && stableRepairVersion.MatchString(version)
	for _, h := range p.hosts {
		piece, err := r.readHostPiece(h.host, extensionPluginID(p.ext.name))
		got := hostPieceVersion(piece, err)
		pieces = append(pieces, h.title+" plugin "+got)
		matched = matched && got == p.entry.Version &&
			(!h.before.Installed || h.before.Scope == piece.Scope && h.before.Enabled == piece.Enabled)
		kept = kept && got == version && h.before.Version == got
	}

	if stepErr == nil && matched {
		info, err := r.formulaInfo(p.ext.homebrewFormula)
		if err == nil {
			var ownedVersion string
			ownedVersion, err = r.ownedBrewProgram(p.ext, info, path)
			if err == nil && ownedVersion != p.entry.Version {
				err = fmt.Errorf("the active program changed while checking its Homebrew ownership")
			}
		}
		if err == nil {
			fmt.Fprintf(stdout, "%s %s verified (%s).\n", p.ext.name, p.entry.Version, strings.Join(pieces, ", "))
			return true
		}
		stepErr = err
	}
	if stepErr == nil {
		stepErr = fmt.Errorf("the installed pieces did not all report %s", p.entry.Version)
	}
	if kept {
		fmt.Fprintf(stderr, "spacedock update: Kept %s (%s). %v.\n", p.ext.name, strings.Join(pieces, ", "), stepErr)
	} else {
		fmt.Fprintf(stderr, "spacedock update: %s repair incomplete (%s). %v.\n", p.ext.name, strings.Join(pieces, ", "), stepErr)
	}
	return false
}
