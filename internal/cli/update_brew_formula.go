package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const spacedockTap = "spacedock-dev/tap"
const compassFormula = spacedockTap + "/compass"
const spacedockTapRemote = "https://github.com/spacedock-dev/homebrew-tap"

type brewFormulaInfo struct {
	Name, Tap string
	FullName  string `json:"full_name"`
	LinkedKeg string `json:"linked_keg"`
	Installed []struct{ Version string }
	Versions  struct{ Stable string }
	URLs      struct {
		Stable struct{ Checksum string }
	}
	Bottle                                 map[string]json.RawMessage
	Pinned, Disabled, Deprecated, Outdated bool
	RecipePath                             string                  `json:"ruby_source_path"`
	RecipeHash                             struct{ SHA256 string } `json:"ruby_source_checksum"`
	UnsafePlan                             string                  `json:"-"`
}

func (r brewRepair) formulaInfo(formula string) (brewFormulaInfo, error) {
	out, err := r.brew("info", "--json=v2", "--formula", formula)
	if err != nil {
		return brewFormulaInfo{}, err
	}
	var data struct {
		Formulae []brewFormulaInfo
		Casks    []json.RawMessage
	}
	if err := json.Unmarshal([]byte(out), &data); err != nil || len(data.Formulae) != 1 || len(data.Casks) != 0 {
		return brewFormulaInfo{}, fmt.Errorf("Homebrew did not return one formula for %s", formula)
	}
	info := data.Formulae[0]
	if info.FullName != formula || info.Tap+"/"+info.Name != formula || !safeRepairSourcePath(info.Name) || strings.Contains(info.Name, "/") {
		return brewFormulaInfo{}, fmt.Errorf("Homebrew could not verify the identity of %s", formula)
	}
	var values struct{ Formulae []map[string]json.RawMessage }
	if err := json.Unmarshal([]byte(out), &values); err != nil {
		return brewFormulaInfo{}, err
	}
	for _, key := range []string{"dependencies", "build_dependencies", "test_dependencies", "recommended_dependencies",
		"optional_dependencies", "uses_from_macos", "requirements", "conflicts_with", "link_overwrite"} {
		var items []json.RawMessage
		if err := json.Unmarshal(values.Formulae[0][key], &items); err != nil || items == nil || len(items) != 0 {
			info.UnsafePlan = key
			break
		}
	}
	return info, nil
}

func (r brewRepair) ownedBrewProgram(ext extension, info brewFormulaInfo, program string) (string, error) {
	if info.LinkedKeg == "" || !safeRepairSourcePath(info.LinkedKeg) || strings.Contains(info.LinkedKeg, "/") {
		return "", fmt.Errorf("%s has no active Homebrew keg. Link its installed formula first", ext.name)
	}
	found := false
	for _, keg := range info.Installed {
		found = found || keg.Version == info.LinkedKeg
	}
	if !found {
		return "", fmt.Errorf("Homebrew did not report the active %s keg among its installations", ext.name)
	}
	cellar, err := r.brew("--cellar", ext.homebrewFormula)
	if err != nil {
		return "", err
	}
	cellar = strings.TrimSpace(cellar)
	if !filepath.IsAbs(cellar) || filepath.Base(cellar) != info.Name {
		return "", fmt.Errorf("Homebrew did not report the expected Cellar for %s", ext.name)
	}
	keg, err := filepath.EvalSymlinks(filepath.Join(cellar, info.LinkedKeg))
	if err != nil {
		return "", fmt.Errorf("could not read the active %s keg: %w", ext.name, err)
	}
	active, err := filepath.EvalSymlinks(program)
	linked, linkErr := filepath.EvalSymlinks(filepath.Join(keg, "bin", ext.program))
	if err != nil || linkErr != nil || active != linked || !strings.HasPrefix(linked, keg+string(os.PathSeparator)) {
		return "", fmt.Errorf("%s on PATH does not resolve to the linked %s executable. Correct PATH or the Homebrew link", ext.program, ext.homebrewFormula)
	}
	version := programReportedVersion(program, r.env, "")
	if version != unrevisionedKeg(info.LinkedKeg) {
		return "", fmt.Errorf("%s reports %s from the linked keg %s. Check the Homebrew installation", ext.program, version, info.LinkedKeg)
	}
	return version, nil
}

func unrevisionedKeg(version string) string {
	base, revision, found := strings.Cut(version, "_")
	if found {
		if _, err := strconv.ParseUint(revision, 10, 64); err == nil {
			return base
		}
	}
	return version
}

type brewFormulaPlan struct {
	tapDir, head, recipeHash, archiveSHA, oldVersion, linked string
}

func (r brewRepair) prepareBrewFormula(ext extension, program, target string) (brewFormulaPlan, error) {
	if ext.homebrewFormula != compassFormula {
		return brewFormulaPlan{}, fmt.Errorf("automatic formula upgrades are unavailable for %s. Verify the formula manually", ext.homebrewFormula)
	}
	tap, err := r.brew("--repo", spacedockTap)
	if err != nil {
		return brewFormulaPlan{}, err
	}
	tap = strings.TrimSpace(tap)
	if !filepath.IsAbs(tap) {
		return brewFormulaPlan{}, fmt.Errorf("Homebrew did not report the official tap directory")
	}
	current, branch, err := r.cleanBrewTap(tap)
	if err != nil {
		return brewFormulaPlan{}, err
	}
	info, err := r.formulaInfo(ext.homebrewFormula)
	if err != nil {
		return brewFormulaPlan{}, err
	}
	oldVersion, err := r.ownedBrewProgram(ext, info, program)
	if err != nil {
		return brewFormulaPlan{}, err
	}
	linked := info.LinkedKeg
	if info.Versions.Stable != target {
		if err := repairCanAdvance("the Homebrew Compass catalog", info.Versions.Stable, target); err != nil {
			return brewFormulaPlan{}, err
		}
		current, err = r.advanceBrewTap(tap, branch, current)
		if err != nil {
			return brewFormulaPlan{}, err
		}
		info, err = r.formulaInfo(ext.homebrewFormula)
		if err != nil {
			return brewFormulaPlan{}, err
		}
	}
	plan := brewFormulaPlan{
		tapDir: tap, head: current, recipeHash: info.RecipeHash.SHA256,
		archiveSHA: info.URLs.Stable.Checksum, oldVersion: oldVersion, linked: linked,
	}
	if err := plan.checkInfo(info, target); err != nil {
		return brewFormulaPlan{}, err
	}
	if _, err := r.brew("fetch", "--formula", ext.homebrewFormula); err != nil {
		return brewFormulaPlan{}, fmt.Errorf("could not fetch the %s %s archive: %w", ext.name, target, err)
	}
	return plan, nil
}

func (p brewFormulaPlan) checkInfo(info brewFormulaInfo, target string) error {
	if info.UnsafePlan != "" {
		return fmt.Errorf("Homebrew cannot scope %s to one package (%s)", info.FullName, info.UnsafePlan)
	}
	if info.Pinned || info.Disabled || info.Deprecated {
		return fmt.Errorf("%s is pinned, disabled or deprecated in Homebrew. Check its formula before retrying", info.FullName)
	}
	if info.Versions.Stable != target {
		return fmt.Errorf("Homebrew offers %s %s. Retry after it offers %s", info.FullName, info.Versions.Stable, target)
	}
	if info.RecipePath != "Formula/compass.rb" ||
		p.recipeHash == "" || info.RecipeHash.SHA256 != p.recipeHash ||
		info.LinkedKeg != p.linked || !info.Outdated || len(info.Bottle) != 0 ||
		p.archiveSHA == "" || info.URLs.Stable.Checksum != p.archiveSHA {
		return fmt.Errorf("Homebrew's selected formula or archive does not match the prepared %s %s upgrade", info.FullName, target)
	}
	return nil
}

func (r brewRepair) applyBrewFormula(ext extension, program, target string, plan brewFormulaPlan) error {
	index, err := r.git(plan.tapDir, "rev-parse", "--path-format=absolute", "--git-path", "index")
	if err != nil {
		return err
	}
	lock := strings.TrimSpace(index) + ".lock"
	if !filepath.IsAbs(lock) {
		return fmt.Errorf("could not locate the tap index to protect the reviewed formula")
	}
	// Keep Git from switching the tap before Homebrew loads the checked recipe.
	unlock, err := holdTapIndex(lock)
	if err != nil {
		return fmt.Errorf("could not lock the tap against another Git update: %w", err)
	}
	defer unlock()

	head, _, err := r.cleanBrewTap(plan.tapDir)
	if err != nil || head != plan.head {
		return fmt.Errorf("the tap or formula changed after preflight. Finish the tap update first")
	}
	info, err := r.formulaInfo(ext.homebrewFormula)
	if err != nil {
		return err
	}
	if err := plan.checkInfo(info, target); err != nil {
		return err
	}
	if version, err := r.ownedBrewProgram(ext, info, program); err != nil || version != plan.oldVersion {
		return fmt.Errorf("the active %s program changed after preflight", ext.name)
	}
	if _, err := r.brew("upgrade", "--formula", ext.homebrewFormula); err != nil {
		return err
	}
	info, err = r.formulaInfo(ext.homebrewFormula)
	if err != nil {
		return err
	}
	version, err := r.ownedBrewProgram(ext, info, program)
	if err != nil {
		return err
	}
	if version != target {
		return fmt.Errorf("Homebrew did not activate %s %s (active %s)", ext.name, target, version)
	}
	return nil
}

func holdTapIndex(path string) (func(), error) {
	// The helper owns the lock until its stdin closes, including when update is interrupted.
	cmd := exec.Command("/bin/sh", "-c", `
set -C
: > "$1" || exit 1
trap '/bin/rm -f "$1"' EXIT
trap 'exit 1' HUP INT TERM
printf 'locked\n'
/bin/cat >/dev/null
`, "spacedock", path)
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		in.Close()
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		in.Close()
		out.Close()
		return nil, err
	}
	unlock := func() {
		in.Close()
		cmd.Wait()
	}
	if line, err := bufio.NewReader(out).ReadString('\n'); err != nil || line != "locked\n" {
		unlock()
		return nil, fmt.Errorf("the Git index is busy or the tap is unwritable")
	}
	return unlock, nil
}
