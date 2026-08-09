package contractlint

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type desiredLiveJourney struct {
	id, test string
	fixtures []string
}

type actualLiveJourney struct {
	id, test, builder, exercise, assertion string
	fixtures                               []string
	todos                                  []liveTODORow
}

type liveTODORow struct{ target, owner string }

var (
	registryHeading, registryEntry = regexp.MustCompile("^### `([^`]+)`$"), regexp.MustCompile("^- \\*\\*Entry point:\\*\\* `([^`]+)`$")
	registryFixture                = regexp.MustCompile("^  - `([^`]+)`")
	journeyBinding                 = regexp.MustCompile(`spacedock:live-journey id=([^ ]+) fixture=([^ ]+)`)
	fixtureBinding                 = regexp.MustCompile(`spacedock:live-fixture id=([^ ]+)`)
	ownerID                        = regexp.MustCompile(`^[a-z0-9]{24}$`)
)

func TestRuntimeLiveRegistryReconciliation(t *testing.T) {
	repo := repoRoot(t)
	registryPath := filepath.Join(repo, "docs", "runtime-live-ci-registry.md")
	desired := readDesiredLiveJourneys(t, registryPath)
	targets := readRegistryTargets(t, registryPath)
	registryFixtures := readRegistryFixtureUnion(t, registryPath)
	actual, fixtureOwners := readActualLiveJourneys(t, repo, targets)
	if len(desired) != 16 || len(actual) != 16 {
		t.Fatalf("common live registry/source counts = %d/%d, want 16/16", len(desired), len(actual))
	}
	gateTODOs := actual["gate-guardrail"].todos
	if want := []liveTODORow{{target: "codex", owner: "3zzpdw704df1g8pg1x9thzmw"}, {target: "pi", owner: "3zzpdw704df1g8pg1x9thzmw"}}; len(gateTODOs) != len(want) || gateTODOs[0] != want[0] || gateTODOs[1] != want[1] {
		t.Fatalf("gate-guardrail TODOs = %#v, want %#v", gateTODOs, want)
	}
	defaultHeadlessTODOs := actual["default-headless-gate-stop"].todos
	if want := []liveTODORow{{target: "claude-sonnet", owner: "26nk8qd48zknqnn4kc123sez"}, {target: "codex", owner: "26nk8qd48zknqnn4kc123sez"}, {target: "pi", owner: "26nk8qd48zknqnn4kc123sez"}}; len(defaultHeadlessTODOs) != len(want) || defaultHeadlessTODOs[0] != want[0] || defaultHeadlessTODOs[1] != want[1] || defaultHeadlessTODOs[2] != want[2] {
		t.Fatalf("default-headless-gate-stop TODOs = %#v, want %#v", defaultHeadlessTODOs, want)
	}

	gapCounts := map[string]int{}
	for id, want := range desired {
		got, ok := actual[id]
		if !ok {
			t.Errorf("registered journey %q has no TestLiveCommon declaration", id)
			continue
		}
		if got.test != want.test {
			t.Errorf("journey %q entry point = %q, want %q", id, got.test, want.test)
		}
		sort.Strings(got.fixtures)
		sort.Strings(want.fixtures)
		if strings.Join(got.fixtures, ",") != strings.Join(want.fixtures, ",") {
			t.Errorf("journey %q fixtures = %v, want %v", id, got.fixtures, want.fixtures)
		}
		seenTargets := map[string]bool{}
		for _, todo := range got.todos {
			if seenTargets[todo.target] {
				t.Errorf("journey %q repeats TODO target %q", id, todo.target)
			}
			seenTargets[todo.target] = true
			gapCounts[todo.target]++
		}
	}
	for id := range fixtureOwners {
		if !registryFixtures[id] {
			t.Errorf("annotated live fixture %q is orphaned", id)
		}
	}
	gapRows := make([]string, 0, len(targets))
	for target := range targets {
		gapRows = append(gapRows, target)
	}
	sort.Strings(gapRows)
	for _, target := range gapRows {
		t.Logf("derived common TODO gap %s=%s", target, strconv.Itoa(gapCounts[target]))
	}

	workflow := string(mustRead(t, filepath.Join(repo, ".github", "workflows", "runtime-live-e2e.yml")))
	docs := string(mustRead(t, filepath.Join(repo, "docs", "runtime-live-ci.md")))
	if strings.Count(workflow, "-run '^TestLiveCommon' -failfast") != 2 || strings.Count(workflow, "-run '^TestLiveCommon'") != 2 {
		t.Errorf("workflow must contain exactly two common selectors with -failfast")
	}
	for _, runtime := range []string{"claude", "codex"} {
		if strings.Count(workflow, "SPACEDOCK_LIVE_RUNTIME="+runtime) != 1 {
			t.Errorf("workflow runtime selector %q is not unique", runtime)
		}
	}
	if strings.Contains(workflow, "SPACEDOCK_LIVE_RUNTIME=pi") || strings.Count(docs, "SPACEDOCK_LIVE_RUNTIME=pi go test") != 1 {
		t.Error("Pi common selector must exist once in the local guide and never in the workflow")
	}
}

func TestRuntimeLiveCommonSuiteTimeouts(t *testing.T) {
	repo := repoRoot(t)
	live := string(mustRead(t, filepath.Join(repo, ".github", "workflows", "runtime-live-e2e.yml")))
	docs := string(mustRead(t, filepath.Join(repo, "docs", "runtime-live-ci.md")))
	for _, command := range []struct {
		name, text, want string
	}{
		{"workflow Claude", live, `SPACEDOCK_LIVE_RUNTIME=claude gotestsum --jsonfile live-e2e-detail.jsonl --format pkgname -- -tags live -count=1 -timeout 90m -run '^TestLiveCommon' -failfast ./internal/ensigncycle/`},
		{"workflow Codex", live, `SPACEDOCK_LIVE_RUNTIME=codex gotestsum --jsonfile codex-shared-scenarios-detail.jsonl --format pkgname -- -tags live -count=1 -timeout 40m -run '^TestLiveCommon' -failfast ./internal/ensigncycle`},
		{"docs Claude", docs, `SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 90m -run '^TestLiveCommon' -failfast ./internal/ensigncycle -v`},
		{"docs Codex", docs, `SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run '^TestLiveCommon' -failfast ./internal/ensigncycle -v`},
		{"docs Pi", docs, `SPACEDOCK_LIVE_RUNTIME=pi go test -tags live -count=1 -timeout 40m -run '^TestLiveCommon' -failfast ./internal/ensigncycle -v`},
	} {
		if count := strings.Count(command.text, command.want); count != 1 {
			t.Errorf("%s common-suite command count = %d, want 1", command.name, count)
		}
	}
}

func readDesiredLiveJourneys(t *testing.T, path string) map[string]desiredLiveJourney {
	text := string(mustRead(t, path))
	start := strings.Index(text, "## Common journeys\n")
	end := strings.Index(text, "## Runtime-specific live proofs")
	if start < 0 || end <= start {
		t.Fatal("registry common-journey boundaries are missing or reversed")
	}
	common := text[start+len("## Common journeys\n") : end]
	out := map[string]desiredLiveJourney{}
	var current desiredLiveJourney
	flush := func() {
		if current.id == "" {
			return
		}
		if current.test == "" || len(current.fixtures) == 0 || out[current.id].id != "" {
			t.Fatalf("malformed or duplicate registry journey %#v", current)
		}
		out[current.id] = current
	}
	for _, line := range strings.Split(common, "\n") {
		if match := registryHeading.FindStringSubmatch(line); match != nil {
			flush()
			current = desiredLiveJourney{id: match[1]}
		} else if match := registryEntry.FindStringSubmatch(line); match != nil {
			current.test = match[1]
		} else if match := registryFixture.FindStringSubmatch(line); match != nil {
			current.fixtures = append(current.fixtures, match[1])
		}
	}
	flush()
	return out
}

func readActualLiveJourneys(t *testing.T, repo string, targets map[string]bool) (map[string]actualLiveJourney, map[string]string) {
	fset := token.NewFileSet()
	var files []*ast.File
	for _, dir := range []string{filepath.Join(repo, "internal", "ensigncycle"), filepath.Join(repo, "internal", "livescenario")} {
		pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		for _, pkg := range pkgs {
			for name, file := range pkg.Files {
				if strings.Contains(dir, "ensigncycle") {
					data := string(mustRead(t, name))
					for _, forbidden := range []string{"sharedRuntimeScenarios", "claudeScenarioRunners", "codexScenarioRunners", "auditedMissingEvidence", "TestLivePiAutoContinueAfterImplementation"} {
						if strings.Contains(data, forbidden) {
							t.Fatalf("common live source retains superseded symbol %q", forbidden)
						}
					}
				}
				files = append(files, file)
			}
		}
	}
	fixtures := map[string]string{}
	actual := map[string]actualLiveJourney{}
	functions := map[string]*ast.FuncDecl{}
	for _, file := range files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if file.Name.Name == "ensigncycle" {
				functions[fn.Name.Name] = fn
			}
			if match := fixtureBinding.FindStringSubmatch(commentText(fn.Doc)); match != nil {
				for _, id := range strings.Split(match[1], ",") {
					if fixtures[id] != "" {
						t.Fatalf("duplicate live fixture annotation %q", id)
					}
					fixtures[id] = fn.Name.Name
				}
			}
			if !strings.HasPrefix(fn.Name.Name, "TestLiveCommon") {
				continue
			}
			match := journeyBinding.FindStringSubmatch(commentText(fn.Doc))
			if match == nil {
				t.Fatalf("%s lacks adjacent live-journey annotation", fn.Name.Name)
			}
			got := parseLiveJourneyCall(t, fn, targets)
			got.test = fn.Name.Name
			if got.id != match[1] || strings.Join(got.fixtures, ",") != match[2] {
				t.Fatalf("%s annotation/call drift: annotation=%v call=%#v", fn.Name.Name, match[1:], got)
			}
			if actual[got.id].id != "" {
				t.Fatalf("duplicate common journey ID %q", got.id)
			}
			actual[got.id] = got
		}
	}
	for id, journey := range actual {
		for _, fixture := range journey.fixtures {
			if fixtures[fixture] != journey.builder {
				t.Fatalf("journey %q builder %q does not own annotated fixture %q (owner %q)", id, journey.builder, fixture, fixtures[fixture])
			}
		}
		if !callsLiveBindings(functions[journey.exercise]) {
			t.Fatalf("journey %q exercise %q does not directly call its bound builder and assertion", id, journey.exercise)
		}
	}
	return actual, fixtures
}

func callsLiveBindings(fn *ast.FuncDecl) bool {
	if fn == nil || len(fn.Type.Params.List) != 5 || len(fn.Type.Params.List[3].Names) != 1 || len(fn.Type.Params.List[4].Names) != 1 {
		return false
	}
	called := map[string]bool{}
	ast.Inspect(fn.Body, func(node ast.Node) bool {
		if call, ok := node.(*ast.CallExpr); ok {
			if callee, direct := call.Fun.(*ast.Ident); direct {
				called[callee.Name] = true
			}
		}
		return true
	})
	return called[fn.Type.Params.List[3].Names[0].Name] && called[fn.Type.Params.List[4].Names[0].Name]
}

func parseLiveJourneyCall(t *testing.T, fn *ast.FuncDecl, targets map[string]bool) actualLiveJourney {
	var calls []*ast.CallExpr
	ast.Inspect(fn.Body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if ok && exprName(call.Fun) == "liveJourney" {
			calls = append(calls, call)
		}
		return true
	})
	if len(calls) != 1 || len(calls[0].Args) != 7 {
		t.Fatalf("%s has %d liveJourney calls; canonical call must have seven arguments", fn.Name.Name, len(calls))
	}
	args := calls[0].Args
	got := actualLiveJourney{id: stringLiteral(t, args[1]), fixtures: strings.Split(stringLiteral(t, args[2]), ","), builder: exprName(args[3]), exercise: exprName(args[5]), assertion: exprName(args[6])}
	if got.builder == "" || got.exercise == "" || got.assertion == "" {
		t.Fatalf("%s binds a string or non-symbol builder/exercise/assertion", fn.Name.Name)
	}
	ident, isNil := args[4].(*ast.Ident)
	if isNil && ident.Name == "nil" {
		return got
	}
	list, ok := args[4].(*ast.CompositeLit)
	if !ok {
		t.Fatalf("%s TODO metadata must be nil or a composite literal", fn.Name.Name)
	}
	for _, element := range list.Elts {
		call, ok := element.(*ast.CallExpr)
		if !ok || exprName(call.Fun) != "liveTODO" {
			t.Fatalf("%s TODO element must be liveTODO(two string literals)", fn.Name.Name)
		}
		if len(call.Args) != 2 {
			t.Fatalf("%s has malformed liveTODO", fn.Name.Name)
		}
		target, owner := stringLiteral(t, call.Args[0]), stringLiteral(t, call.Args[1])
		if !ownerID.MatchString(owner) || !targets[target] {
			t.Fatalf("%s has malformed TODO(%q, %q)", fn.Name.Name, target, owner)
		}
		got.todos = append(got.todos, liveTODORow{target: target, owner: owner})
	}
	return got
}

func exprName(expr ast.Expr) string {
	switch value := expr.(type) {
	case *ast.Ident:
		return value.Name
	case *ast.SelectorExpr:
		if prefix := exprName(value.X); prefix != "" {
			return prefix + "." + value.Sel.Name
		}
	}
	return ""
}

func stringLiteral(t *testing.T, expr ast.Expr) string {
	literal, ok := expr.(*ast.BasicLit)
	if !ok || literal.Kind != token.STRING {
		t.Fatalf("live metadata argument is not a string literal")
	}
	value, err := strconv.Unquote(literal.Value)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func commentText(group *ast.CommentGroup) string {
	if group == nil {
		return ""
	}
	var comments []string
	for _, comment := range group.List {
		comments = append(comments, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")))
	}
	return strings.Join(comments, "\n")
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
func readRegistryTargets(t *testing.T, path string) map[string]bool {
	out := map[string]bool{}
	for _, match := range regexp.MustCompile(`(?m)^\| ([A-Za-z ]+) \|`).FindAllStringSubmatch(string(mustRead(t, path)), -1) {
		name := strings.ToLower(strings.TrimSpace(match[1]))
		if name != "target" {
			out[strings.ReplaceAll(name, " ", "-")] = true
		}
	}
	if len(out) == 0 {
		t.Fatal("registry target table is empty")
	}
	return out
}

func readRegistryFixtureUnion(t *testing.T, path string) map[string]bool {
	text := strings.Split(string(mustRead(t, path)), "## Source binding convention")[0]
	out := map[string]bool{}
	for _, pattern := range []string{"(?m)^  - `([^`]+)`", `(?m)^- \*\*Fixture:\*\* ` + "`([^`]+)`"} {
		for _, match := range regexp.MustCompile(pattern).FindAllStringSubmatch(text, -1) {
			out[match[1]] = true
		}
	}
	return out
}
