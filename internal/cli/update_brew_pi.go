package cli

import (
	"fmt"
	"strings"
)

type brewPiPiece struct {
	agentDir, home, name, before string
	source, commit               string
	pending                      bool
}

func (p brewPiPiece) version() string {
	_, ref, ok := piRegisteredPackage(p.agentDir, p.home, p.name)
	if !ok {
		return notInstalled
	}
	return strings.TrimPrefix(refOrUnpinned(ref), "v")
}

func (p *brewPiPiece) prepare(r brewRepair, entry repairEntry) error {
	if entry.Repo == "" || entry.Path != "." || (entry.Ref != entry.Version && entry.Ref != "v"+entry.Version) {
		return fmt.Errorf("Pi's %s source does not use a versioned tag. Verify its published Pi package", p.name)
	}
	source := "git:github.com/" + entry.Repo
	if err := repairCanAdvance("Pi package", p.before, entry.Version); err != nil {
		return err
	}
	sources := 0
	for _, found := range piSettingsSources(p.agentDir) {
		bare, _, _ := strings.Cut(found, "@")
		if bare == source {
			sources++
		}
	}
	if sources > 1 || sources == 0 && p.before != notInstalled {
		return fmt.Errorf("Pi holds an unexpected or duplicate %s source. Correct the package registration", p.name)
	}
	commit, err := r.resolveRepairRef(entry.Repo, entry.Ref, true)
	if err != nil {
		return err
	}
	p.source, p.commit, p.pending = source, commit, p.before != entry.Version
	return nil
}

func (p brewPiPiece) apply(r brewRepair, entry repairEntry, ops piRuntimeOps) error {
	commit, err := r.resolveRepairRef(entry.Repo, entry.Ref, true)
	if err != nil {
		return err
	}
	if commit != p.commit || p.version() != p.before {
		return fmt.Errorf("Pi's %s source or package changed while the repair was being prepared", p.name)
	}
	if p.before != notInstalled {
		if out, err := ops.PiRemove(p.source); err != nil {
			return fmt.Errorf("Pi could not remove the older %s package: %w (%s)", p.name, err, strings.TrimSpace(out))
		}
	}
	if out, err := ops.PiInstall(p.source + "@" + entry.Ref); err != nil {
		return fmt.Errorf("Pi could not install %s %s: %w (%s)", p.name, entry.Version, err, strings.TrimSpace(out))
	}
	if p.version() != entry.Version {
		return fmt.Errorf("Pi did not register %s %s", p.name, entry.Version)
	}
	return nil
}
