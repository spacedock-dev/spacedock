package contractlint

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var liveAnnotation = regexp.MustCompile(`(?m)//spacedock:live-(journey|fixture|suite|proof)\s+([^\n]+)`)
var annotationField = regexp.MustCompile(`([a-z]+)=([a-zA-Z0-9_./,-]+)`)
var registryH3 = regexp.MustCompile("(?m)^### `([^`]+)`$")
var registryCodeToken = regexp.MustCompile("`([^`]+)`")

type liveInventory struct {
	journeys      []string
	fixtures      []string
	proofs        map[string]string
	proofEntries  map[string]string
	suites        []string
	suiteEntries  []string
	topLevelTests []string
	diagnostics   []string
}

type desiredLiveRegistry struct {
	journeys        []string
	fixtures        []string
	proofs          map[string]string
	proofEntries    map[string]string
	targets         map[string]string
	suite           string
	experimentTests []string
	missingEvidence map[string]string
	diagnostics     []string
}

func TestRuntimeLiveRegistryReconciliation(t *testing.T) {
	repo := repoRoot(t)
	desired := parseDesiredLiveRegistry(t, readContractFile(t, filepath.Join(repo, "docs", "runtime-live-ci-registry.md")))
	actual, err := collectLiveInventory(repo)
	if err != nil {
		t.Fatal(err)
	}
	workflow := readContractFile(t, filepath.Join(repo, ".github", "workflows", "runtime-live-e2e.yml"))
	scenarios := readContractFile(t, filepath.Join(repo, "internal", "ensigncycle", "shared_scenarios_test.go"))
	runner := readContractFile(t, filepath.Join(repo, "internal", "ensigncycle", "shared_live_runner_test.go"))
	if err := reconcileLiveRegistry(desired, actual, workflow, scenarios, runner); err != nil {
		t.Fatal(err)
	}
}

func TestRuntimeLiveRegistryReconciliationMutationControls(t *testing.T) {
	repo := repoRoot(t)
	registryText := readContractFile(t, filepath.Join(repo, "docs", "runtime-live-ci-registry.md"))
	desired := parseDesiredLiveRegistry(t, registryText)
	actual, err := collectLiveInventory(repo)
	if err != nil {
		t.Fatal(err)
	}
	workflow := readContractFile(t, filepath.Join(repo, ".github", "workflows", "runtime-live-e2e.yml"))
	scenarios := readContractFile(t, filepath.Join(repo, "internal", "ensigncycle", "shared_scenarios_test.go"))
	runner := readContractFile(t, filepath.Join(repo, "internal", "ensigncycle", "shared_live_runner_test.go"))

	actualMutations := map[string]func(*liveInventory){
		"deleted journey binding":         func(i *liveInventory) { i.journeys = i.journeys[1:] },
		"duplicate journey binding":       func(i *liveInventory) { i.journeys = append(i.journeys, i.journeys[0]) },
		"desired fixture without builder": func(i *liveInventory) { i.fixtures = i.fixtures[1:] },
		"missing live test":               func(i *liveInventory) { i.topLevelTests = i.topLevelTests[1:] },
		"orphan live test":                func(i *liveInventory) { i.topLevelTests = append(i.topLevelTests, "TestLiveOrphanControl") },
		"changed proof lane": func(i *liveInventory) {
			for id := range i.proofs {
				i.proofs[id] = "other-live"
				break
			}
		},
		"competing suite": func(i *liveInventory) { i.suites = append(i.suites, i.suites[0]) },
	}
	for name, mutate := range actualMutations {
		t.Run(name, func(t *testing.T) {
			candidate := cloneLiveInventory(actual)
			mutate(&candidate)
			if reconcileLiveRegistry(desired, candidate, workflow, scenarios, runner) == nil {
				t.Fatal("source binding mutation escaped reconciliation")
			}
		})
	}

	t.Run("duplicate proof annotation", func(t *testing.T) {
		candidate := liveInventory{proofs: map[string]string{}, proofEntries: map[string]string{}}
		line := "//spacedock:live-proof id=proof-control lane=claude-live\n"
		collectLiveAnnotations("control.go", line+"func TestLiveControlA(t *testing.T) {}\n"+line+"func TestLiveControlB(t *testing.T) {}\n", &candidate)
		if !containsDiagnostic(candidate.diagnostics, "DUPLICATE proof annotation proof-control") {
			t.Fatalf("duplicate proof annotation diagnostics = %v", candidate.diagnostics)
		}
	})
	t.Run("malformed proof annotation", func(t *testing.T) {
		candidate := liveInventory{proofs: map[string]string{}, proofEntries: map[string]string{}}
		collectLiveAnnotations("control.go", "//spacedock:live-proof id= lane=\n", &candidate)
		if !containsDiagnostic(candidate.diagnostics, "INVALID proof annotation control.go") {
			t.Fatalf("malformed proof annotation diagnostics = %v", candidate.diagnostics)
		}
	})
	t.Run("detached fixture annotation", func(t *testing.T) {
		candidate := liveInventory{proofs: map[string]string{}, proofEntries: map[string]string{}}
		collectLiveAnnotations("control.go", "//spacedock:live-fixture id=detached/control\n\nfunc writeDetached() {}\n", &candidate)
		if !containsDiagnostic(candidate.diagnostics, "INVALID fixture annotation control.go") {
			t.Fatalf("detached fixture diagnostics = %v", candidate.diagnostics)
		}
	})

	registryMutations := map[string]func(string) string{
		"changed desired journey": func(s string) string {
			s = strings.Replace(s, "### `full-ensign-cycle`", "### `renamed-ensign-cycle`", 1)
			return strings.Replace(s, "`TestLiveSharedScenarios/full-ensign-cycle`", "`TestLiveSharedScenarios/renamed-ensign-cycle`", 1)
		},
		"changed desired fixture": func(s string) string {
			return strings.Replace(s, "`realistic-lifecycle` —", "`renamed-lifecycle` —", 1)
		},
		"changed desired proof": func(s string) string {
			return strings.Replace(s, "### `claude-bare-dispatch`", "### `claude-other-dispatch`", 1)
		},
		"duplicate desired proof": func(s string) string {
			block := "### `claude-bare-dispatch`\n\n- **Entry point:** `TestLiveBareReachable`\n- **Lane:** `claude-live`\n- **Fixture:** `dispatch-recovery/base` — duplicate control.\n\n"
			return strings.Replace(s, "## Non-gating live experiments", block+"## Non-gating live experiments", 1)
		},
		"changed proof lane": func(s string) string {
			return strings.Replace(s, "- **Lane:** `pi-live`", "- **Lane:** `other-live`", 1)
		},
		"changed target lane": func(s string) string {
			return strings.Replace(s, "| Codex | `codex-live` |", "| Codex | `other-live` |", 1)
		},
		"duplicate target row": func(s string) string {
			return strings.Replace(s, "| Codex | `codex-live` |", "| Codex | `codex-live` |\n| Codex | `codex-live` |", 1)
		},
		"changed suite entry": func(s string) string {
			return strings.ReplaceAll(s, "`TestLiveSharedScenarios/", "`TestLiveOtherScenarios/")
		},
		"deleted TODO row": func(s string) string {
			return strings.Replace(s, "| `auto-continue-after-implementation` | `9adv48yhye5s2vkhwd7ge52d` |\n", "", 1)
		},
		"duplicate TODO row": func(s string) string {
			row := "| `auto-continue-after-implementation` | `9adv48yhye5s2vkhwd7ge52d` |"
			return strings.Replace(s, row, row+"\n"+row, 1)
		},
		"retagged TODO owner": func(s string) string {
			return strings.Replace(s, "| `keep-moving-posture` | `9adv48yhye5s2vkhwd7ge52d` |", "| `keep-moving-posture` | `otherowner` |", 1)
		},
	}
	for name, mutate := range registryMutations {
		t.Run(name, func(t *testing.T) {
			mutated := mutate(registryText)
			if mutated == registryText {
				t.Fatal("registry mutation did not apply")
			}
			candidate := parseDesiredLiveRegistry(t, mutated)
			if reconcileLiveRegistry(candidate, actual, workflow, scenarios, runner) == nil {
				t.Fatal("normative registry mutation escaped reconciliation")
			}
		})
	}

	evidence, err := extractMissingEvidence(scenarios)
	if err != nil {
		t.Fatal(err)
	}
	evidenceMutations := map[string]func(map[string]string){
		"removed TODO": func(m map[string]string) {
			for id := range m {
				delete(m, id)
				break
			}
		},
		"unowned TODO": func(m map[string]string) { m["shallow-boot"] = "" },
		"retagged TODO": func(m map[string]string) {
			for id := range m {
				m[id] = "other-owner"
				break
			}
		},
	}
	for name, mutate := range evidenceMutations {
		t.Run(name, func(t *testing.T) {
			candidate := cloneStringMap(evidence)
			mutate(candidate)
			if reconcileMissingEvidence(desired.missingEvidence, candidate) == nil {
				t.Fatal("missing-evidence mutation escaped reconciliation")
			}
		})
	}
}

func parseDesiredLiveRegistry(t *testing.T, source string) desiredLiveRegistry {
	t.Helper()
	common := registrySection(t, source, "Common journeys")
	proofs := registrySection(t, source, "Runtime-specific live proofs")
	experiments := registrySection(t, source, "Non-gating live experiments")
	missing := registrySection(t, source, "Missing live evidence")
	targets := registrySection(t, source, "Supported runtime targets")

	desired := desiredLiveRegistry{proofs: map[string]string{}, proofEntries: map[string]string{}, targets: map[string]string{}, missingEvidence: map[string]string{}}
	fixtureSet := map[string]bool{}
	for _, block := range registryBlocks(common) {
		desired.journeys = append(desired.journeys, block.id)
		entry := registryLabeledCode(block.body, "Entry point")
		suite, id, ok := strings.Cut(entry, "/")
		if !ok || id != block.id {
			t.Fatalf("common journey %q entry point = %q", block.id, entry)
		}
		if desired.suite == "" {
			desired.suite = suite
		} else if desired.suite != suite {
			t.Fatalf("common journeys name competing suites %q and %q", desired.suite, suite)
		}
		for _, id := range registryFixtureIDs(block.body) {
			fixtureSet[id] = true
		}
	}
	for _, block := range registryBlocks(proofs) {
		lane := registryLabeledCode(block.body, "Lane")
		addUniqueBinding(desired.proofs, block.id, lane, "proof", &desired.diagnostics)
		addUniqueBinding(desired.proofEntries, block.id, registryLabeledCode(block.body, "Entry point"), "proof entry point", &desired.diagnostics)
		for _, id := range registryFixtureIDs(block.body) {
			fixtureSet[id] = true
		}
	}
	for _, block := range registryBlocks(experiments) {
		desired.experimentTests = append(desired.experimentTests, block.id)
		for _, id := range registryFixtureIDs(block.body) {
			fixtureSet[id] = true
		}
	}
	for _, line := range strings.Split(missing, "\n") {
		if !strings.HasPrefix(line, "|") || strings.Contains(line, "---") || strings.Contains(line, "Common journey") {
			continue
		}
		fields := registryCodeToken.FindAllStringSubmatch(line, -1)
		if len(fields) != 2 {
			desired.diagnostics = append(desired.diagnostics, "INVALID missing-evidence row "+line)
			continue
		}
		addUniqueBinding(desired.missingEvidence, fields[0][1], fields[1][1], "missing-evidence", &desired.diagnostics)
	}
	for _, line := range strings.Split(targets, "\n") {
		if !strings.HasPrefix(line, "|") || strings.Contains(line, "---") || strings.Contains(line, "Target") {
			continue
		}
		cells := strings.Split(line, "|")
		if len(cells) < 4 {
			continue
		}
		name := strings.TrimSpace(cells[1])
		match := registryCodeToken.FindStringSubmatch(cells[2])
		if name == "" || len(match) != 2 {
			desired.diagnostics = append(desired.diagnostics, "INVALID target row "+line)
			continue
		}
		addUniqueBinding(desired.targets, name, match[1], "target", &desired.diagnostics)
	}
	for id := range fixtureSet {
		desired.fixtures = append(desired.fixtures, id)
	}
	sort.Strings(desired.fixtures)
	if desired.suite == "" || len(desired.journeys) == 0 || len(desired.proofs) == 0 || len(desired.targets) == 0 || len(desired.missingEvidence) == 0 {
		t.Fatal("normative live registry parser produced an empty required class")
	}
	return desired
}

type registryBlock struct{ id, body string }

func registryBlocks(section string) []registryBlock {
	matches := registryH3.FindAllStringSubmatchIndex(section, -1)
	blocks := make([]registryBlock, 0, len(matches))
	for i, match := range matches {
		end := len(section)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}
		blocks = append(blocks, registryBlock{id: section[match[2]:match[3]], body: section[match[1]:end]})
	}
	return blocks
}

func registrySection(t *testing.T, source, heading string) string {
	t.Helper()
	marker := "## " + heading + "\n"
	start := strings.Index(source, marker)
	if start < 0 {
		t.Fatalf("registry section %q not found", heading)
	}
	rest := source[start+len(marker):]
	if end := strings.Index(rest, "\n## "); end >= 0 {
		rest = rest[:end]
	}
	return rest
}

func registryLabeledCode(body, label string) string {
	prefix := "- **" + label + ":**"
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, prefix) {
			if match := registryCodeToken.FindStringSubmatch(line); len(match) == 2 {
				return match[1]
			}
		}
	}
	return ""
}

func registryFixtureIDs(body string) []string {
	var ids []string
	inFixtures := false
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "- **Fixture") {
			inFixtures = true
			if match := registryCodeToken.FindStringSubmatch(line); len(match) == 2 {
				ids = append(ids, match[1])
			}
			continue
		}
		if inFixtures && strings.HasPrefix(line, "  - `") {
			if match := registryCodeToken.FindStringSubmatch(line); len(match) == 2 {
				ids = append(ids, match[1])
			}
			continue
		}
		if inFixtures && strings.HasPrefix(line, "- **") {
			inFixtures = false
		}
	}
	return ids
}

func collectLiveInventory(repo string) (liveInventory, error) {
	inv := liveInventory{proofs: map[string]string{}, proofEntries: map[string]string{}}
	for _, root := range []string{filepath.Join(repo, "internal", "ensigncycle"), filepath.Join(repo, "internal", "livescenario")} {
		err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || filepath.Ext(path) != ".go" {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			collectLiveAnnotations(path, string(data), &inv)
			if strings.HasPrefix(string(data), "//go:build live\n") {
				file, err := parser.ParseFile(token.NewFileSet(), path, data, 0)
				if err != nil {
					return err
				}
				for _, decl := range file.Decls {
					fn, ok := decl.(*ast.FuncDecl)
					if ok && fn.Recv == nil && strings.HasPrefix(fn.Name.Name, "TestLive") {
						inv.topLevelTests = append(inv.topLevelTests, fn.Name.Name)
					}
				}
			}
			return nil
		})
		if err != nil {
			return liveInventory{}, err
		}
	}
	return inv, nil
}

func collectLiveAnnotations(path, source string, inv *liveInventory) {
	lines := strings.Split(source, "\n")
	funcDecl := regexp.MustCompile(`^func ([A-Za-z0-9_]+)\(`)
	nameField := regexp.MustCompile(`name:\s*"([a-z0-9-]+)"`)
	for index, line := range lines {
		match := liveAnnotation.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		fields := map[string]string{}
		for _, field := range annotationField.FindAllStringSubmatch(match[2], -1) {
			fields[field[1]] = field[2]
		}
		next := ""
		if index+1 < len(lines) {
			next = strings.TrimSpace(lines[index+1])
		}
		switch match[1] {
		case "journey":
			if fields["id"] == "" || next != "{" {
				inv.diagnostics = append(inv.diagnostics, "INVALID journey annotation "+path)
			} else {
				bound := ""
				for scan := index + 2; scan < len(lines) && scan <= index+14 && !strings.Contains(lines[scan], "//spacedock:live-"); scan++ {
					if found := nameField.FindStringSubmatch(lines[scan]); len(found) == 2 {
						bound = found[1]
						break
					}
				}
				if bound != fields["id"] {
					inv.diagnostics = append(inv.diagnostics, "INVALID journey record "+fields["id"]+"="+bound)
				} else {
					inv.journeys = append(inv.journeys, fields["id"])
				}
			}
		case "fixture":
			decl := funcDecl.FindStringSubmatch(next)
			if fields["id"] == "" || len(decl) != 2 {
				inv.diagnostics = append(inv.diagnostics, "INVALID fixture annotation "+path)
			} else {
				inv.fixtures = append(inv.fixtures, fields["id"])
			}
		case "proof":
			decl := funcDecl.FindStringSubmatch(next)
			if len(decl) != 2 {
				inv.diagnostics = append(inv.diagnostics, "INVALID proof annotation "+path)
				continue
			}
			addUniqueBinding(inv.proofs, fields["id"], fields["lane"], "proof annotation", &inv.diagnostics)
			addUniqueBinding(inv.proofEntries, fields["id"], decl[1], "proof declaration", &inv.diagnostics)
		case "suite":
			decl := funcDecl.FindStringSubmatch(next)
			if fields["lanes"] == "" || len(decl) != 2 {
				inv.diagnostics = append(inv.diagnostics, "INVALID suite annotation "+path)
			} else {
				inv.suites = append(inv.suites, fields["lanes"])
				inv.suiteEntries = append(inv.suiteEntries, decl[1])
			}
		}
	}
}

func reconcileLiveRegistry(desired desiredLiveRegistry, actual liveInventory, workflow, scenarios, runner string) error {
	diagnostics := append(append([]string(nil), desired.diagnostics...), actual.diagnostics...)
	diagnostics = append(diagnostics, compareIDs("journey", actual.journeys, desired.journeys)...)
	diagnostics = append(diagnostics, compareFixtureBuilders(actual.fixtures, desired.fixtures)...)
	diagnostics = append(diagnostics, compareMap("proof", actual.proofs, desired.proofs)...)
	diagnostics = append(diagnostics, compareMap("proof entry", actual.proofEntries, desired.proofEntries)...)
	if len(actual.suiteEntries) != 1 || first(actual.suiteEntries) != desired.suite {
		diagnostics = append(diagnostics, fmt.Sprintf("UNACCOUNTED-TEST suite declarations=%v want=%s", actual.suiteEntries, desired.suite))
	}
	wantTests := []string{desired.suite}
	for _, entry := range desired.proofEntries {
		wantTests = append(wantTests, entry)
	}
	wantTests = append(wantTests, desired.experimentTests...)
	diagnostics = append(diagnostics, compareLiveTests(actual.topLevelTests, wantTests)...)

	wantLanes := uniqueMapValues(desired.targets)
	if len(actual.suites) != 1 || !sameSet(strings.Split(first(actual.suites), ","), wantLanes) {
		diagnostics = append(diagnostics, fmt.Sprintf("UNSELECTED suite lanes=%v want=%v", actual.suites, wantLanes))
	}
	for _, id := range desired.journeys {
		if strings.Count(runner, `"`+id+`":`) != 1 {
			diagnostics = append(diagnostics, "MISSING runner "+id)
		}
	}
	if !strings.Contains(runner, "liveDurableJourneyTODO(scenario.name)") || strings.Count(runner, "t.Skip(reason)") != 1 {
		diagnostics = append(diagnostics, "UNOWNED-EVIDENCE common TODO accounting")
	}
	evidence, err := extractMissingEvidence(scenarios)
	if err != nil {
		diagnostics = append(diagnostics, err.Error())
	} else if err := reconcileMissingEvidence(desired.missingEvidence, evidence); err != nil {
		diagnostics = append(diagnostics, err.Error())
	}
	for _, lane := range wantLanes {
		runtime := strings.TrimSuffix(lane, "-live")
		if strings.Count(workflow, "SPACEDOCK_LIVE_RUNTIME: "+runtime) != 1 {
			diagnostics = append(diagnostics, "UNSELECTED runtime "+runtime)
		}
	}
	selector := "-run '^" + desired.suite + "$' ./internal/ensigncycle"
	if strings.Count(workflow, selector) != len(wantLanes) {
		diagnostics = append(diagnostics, fmt.Sprintf("UNSELECTED suite %s count=%d", desired.suite, strings.Count(workflow, selector)))
	}
	if len(diagnostics) != 0 {
		sort.Strings(diagnostics)
		return fmt.Errorf("live registry reconciliation: %s", strings.Join(diagnostics, "; "))
	}
	return nil
}

func extractMissingEvidence(source string) (map[string]string, error) {
	ownerMatch := regexp.MustCompile(`const liveDurableJourneyDefectID = "([a-z0-9]+)"`).FindStringSubmatch(source)
	if len(ownerMatch) != 2 {
		return nil, fmt.Errorf("MISSING-EVIDENCE repair owner")
	}
	start := strings.Index(source, "func liveDurableJourneyTODO")
	if start < 0 {
		return nil, fmt.Errorf("MISSING-EVIDENCE TODO function")
	}
	rest := source[start:]
	if end := strings.Index(rest[1:], "\nfunc "); end >= 0 {
		rest = rest[:end+1]
	}
	out := map[string]string{}
	for _, match := range regexp.MustCompile(`case "([a-z0-9-]+)"`).FindAllStringSubmatch(rest, -1) {
		out[match[1]] = ownerMatch[1]
	}
	return out, nil
}

func reconcileMissingEvidence(want, got map[string]string) error {
	var diagnostics []string
	for id, owner := range want {
		gotOwner, ok := got[id]
		if !ok {
			diagnostics = append(diagnostics, "MISSING-EVIDENCE TODO "+id)
		} else if gotOwner != owner {
			diagnostics = append(diagnostics, "MISSING-EVIDENCE owner "+id+"="+gotOwner)
		}
	}
	for id, owner := range got {
		if _, ok := want[id]; !ok {
			diagnostics = append(diagnostics, "UNOWNED-EVIDENCE "+id+"="+owner)
		}
	}
	if len(diagnostics) != 0 {
		sort.Strings(diagnostics)
		return fmt.Errorf("live evidence reconciliation: %s", strings.Join(diagnostics, "; "))
	}
	return nil
}

func compareIDs(kind string, got, want []string) []string {
	gotSet, duplicates := setAndDuplicates(got)
	wantSet, wantDuplicates := setAndDuplicates(want)
	var diagnostics []string
	for _, id := range append(duplicates, wantDuplicates...) {
		diagnostics = append(diagnostics, "DUPLICATE "+kind+" "+id)
	}
	for id := range wantSet {
		if !gotSet[id] {
			diagnostics = append(diagnostics, "MISSING "+kind+" "+id)
		}
	}
	for id := range gotSet {
		if !wantSet[id] {
			diagnostics = append(diagnostics, "ORPHAN "+kind+" "+id)
		}
	}
	return diagnostics
}

func compareFixtureBuilders(got, want []string) []string {
	gotSet, duplicates := setAndDuplicates(got)
	wantSet, wantDuplicates := setAndDuplicates(want)
	var diagnostics []string
	for _, id := range append(duplicates, wantDuplicates...) {
		diagnostics = append(diagnostics, "DUPLICATE fixture "+id)
	}
	for id := range wantSet {
		if !gotSet[id] {
			diagnostics = append(diagnostics, "UNACCOUNTED-BUILDER "+id)
		}
	}
	for id := range gotSet {
		if !wantSet[id] {
			diagnostics = append(diagnostics, "ORPHAN fixture "+id)
		}
	}
	return diagnostics
}

func compareLiveTests(got, want []string) []string {
	gotSet, duplicates := setAndDuplicates(got)
	wantSet, wantDuplicates := setAndDuplicates(want)
	var diagnostics []string
	for _, id := range append(duplicates, wantDuplicates...) {
		diagnostics = append(diagnostics, "DUPLICATE live test "+id)
	}
	for id := range wantSet {
		if !gotSet[id] {
			diagnostics = append(diagnostics, "MISSING live test "+id)
		}
	}
	for id := range gotSet {
		if !wantSet[id] {
			diagnostics = append(diagnostics, "UNACCOUNTED-TEST "+id)
		}
	}
	return diagnostics
}

func compareMap(kind string, got, want map[string]string) []string {
	var diagnostics []string
	for id, value := range want {
		if gotValue, ok := got[id]; !ok {
			diagnostics = append(diagnostics, "MISSING "+kind+" "+id)
		} else if gotValue != value {
			diagnostics = append(diagnostics, "INVALID "+kind+" "+id+"="+gotValue)
		}
	}
	for id := range got {
		if _, ok := want[id]; !ok {
			diagnostics = append(diagnostics, "ORPHAN "+kind+" "+id)
		}
	}
	return diagnostics
}

func setAndDuplicates(values []string) (map[string]bool, []string) {
	set := map[string]bool{}
	var duplicates []string
	for _, value := range values {
		if value == "" || set[value] {
			duplicates = append(duplicates, value)
		}
		set[value] = true
	}
	return set, duplicates
}
func sameSet(a, b []string) bool {
	as, _ := setAndDuplicates(a)
	bs, _ := setAndDuplicates(b)
	if len(as) != len(bs) {
		return false
	}
	for value := range as {
		if !bs[value] {
			return false
		}
	}
	return true
}
func uniqueMapValues(values map[string]string) []string {
	set := map[string]bool{}
	for _, value := range values {
		set[value] = true
	}
	var out []string
	for value := range set {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}
func first(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
func cloneStringMap(values map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range values {
		out[k] = v
	}
	return out
}
func cloneLiveInventory(in liveInventory) liveInventory {
	return liveInventory{journeys: append([]string(nil), in.journeys...), fixtures: append([]string(nil), in.fixtures...), proofs: cloneStringMap(in.proofs), proofEntries: cloneStringMap(in.proofEntries), suites: append([]string(nil), in.suites...), suiteEntries: append([]string(nil), in.suiteEntries...), topLevelTests: append([]string(nil), in.topLevelTests...), diagnostics: append([]string(nil), in.diagnostics...)}
}

func addUniqueBinding(values map[string]string, id, value, kind string, diagnostics *[]string) {
	if id == "" || value == "" {
		*diagnostics = append(*diagnostics, "INVALID "+kind)
		return
	}
	if _, exists := values[id]; exists {
		*diagnostics = append(*diagnostics, "DUPLICATE "+kind+" "+id)
		return
	}
	values[id] = value
}

func containsDiagnostic(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
func readContractFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
