package contractlint

import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

type desiredLiveJourney struct {
	id, test string
	fixtures []string
}

type actualLiveJourney struct {
	id, test, builder, exercise, assertion string
	fixtures                               []string
	gaps                                   []liveGapRow
}

type liveGapRow struct{ kind, target, owner string }

var (
	registryHeading, registryEntry = regexp.MustCompile("^### `([^`]+)`$"), regexp.MustCompile("^- \\*\\*Entry point:\\*\\* `([^`]+)`$")
	registryFixture                = regexp.MustCompile("^  - `([^`]+)`")
	registryTestHeading            = regexp.MustCompile("(?m)^### `(Test[^`]+)`$")
	registryEntryPoint             = regexp.MustCompile("(?m)^- \\*\\*Entry point:\\*\\* `(Test[^`]+)`$")
	journeyBinding                 = regexp.MustCompile(`spacedock:live-journey id=([^ ]+) fixture=([^ ]+)`)
	fixtureBinding                 = regexp.MustCompile(`spacedock:live-fixture id=([^ ]+)`)
	ownerID                        = regexp.MustCompile(`^[a-z0-9]{24}$`)
	frontmatterID                  = regexp.MustCompile(`(?m)^id: ([a-z0-9]{24})$`)
	frontmatterStatus              = regexp.MustCompile(`(?m)^status: ([a-z-]+)$`)
)

func TestRuntimeLiveRegistryReconciliation(t *testing.T) {
	repo := repoRoot(t)
	registryPath := filepath.Join(repo, "docs", "runtime-live-ci-registry.md")
	desired := readDesiredLiveJourneys(t, registryPath)
	targets := readRegistryTargets(t, registryPath)
	registryFixtures := readRegistryFixtureUnion(t, registryPath)
	actual, fixtureOwners := readActualLiveJourneys(t, repo, targets)
	// Fatals on a malformed live-proof gap binding; the owner join re-derives these
	// under SPACEDOCK_LIVE_STATE_DIR.
	readRuntimeProofGaps(t, repo, targets)
	if len(desired) != len(actual) {
		t.Errorf("common live registry/source counts = %d/%d", len(desired), len(actual))
	}
	reconcileRegisteredLiveTests(t, repo, registryPath)
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
		for _, gap := range got.gaps {
			if seenTargets[gap.target] {
				t.Errorf("journey %q repeats gap target %q", id, gap.target)
			}
			seenTargets[gap.target] = true
			gapCounts[gap.kind+"/"+gap.target]++
		}
	}
	for id := range fixtureOwners {
		if !registryFixtures[id] {
			t.Errorf("annotated live fixture %q is orphaned", id)
		}
	}
	gapRows := make([]string, 0, len(targets))
	for target := range gapCounts {
		gapRows = append(gapRows, target)
	}
	sort.Strings(gapRows)
	for _, target := range gapRows {
		t.Logf("derived common gap %s=%s", target, strconv.Itoa(gapCounts[target]))
	}

	workflow := string(mustRead(t, filepath.Join(repo, ".github", "workflows", "runtime-live-e2e.yml")))
	docs := string(mustRead(t, filepath.Join(repo, "docs", "runtime-live-ci.md")))
	for _, runtime := range []string{"claude", "codex", "pi"} {
		if strings.Count(workflow, "SPACEDOCK_LIVE_RUNTIME="+runtime) != 1 {
			t.Errorf("workflow runtime selector %q is not unique", runtime)
		}
	}
	if strings.Count(docs, "SPACEDOCK_LIVE_RUNTIME=pi go test") != 1 {
		t.Error("Pi common selector must exist once in the local guide")
	}
}

func TestRuntimeLiveTODOOwnersAreActive(t *testing.T) {
	stateDir := os.Getenv("SPACEDOCK_LIVE_STATE_DIR")
	if stateDir == "" {
		t.Skip("set SPACEDOCK_LIVE_STATE_DIR to run the mutable state-owner join")
	}
	if !filepath.IsAbs(stateDir) {
		stateDir = filepath.Join(repoRoot(t), stateDir)
	}
	repo := repoRoot(t)
	registryPath := filepath.Join(repo, "docs", "runtime-live-ci-registry.md")
	actual, _ := readActualLiveJourneys(t, repo, readRegistryTargets(t, registryPath))
	proofGaps := readRuntimeProofGaps(t, repo, readRegistryTargets(t, registryPath))
	active := readActiveEntityIDs(t, stateDir)
	for journeyID, journey := range actual {
		for _, gap := range journey.gaps {
			if !active[gap.owner] {
				t.Errorf("journey %q target %q names inactive %s owner %q", journeyID, gap.target, gap.kind, gap.owner)
			}
		}
	}
	for proof, gaps := range proofGaps {
		for _, gap := range gaps {
			if !active[gap.owner] {
				t.Errorf("runtime proof %q target %q names inactive %s owner %q", proof, gap.target, gap.kind, gap.owner)
			}
		}
	}
}

func readRuntimeProofGaps(t *testing.T, repo string, targets map[string]bool) map[string][]liveGapRow {
	t.Helper()
	result := map[string][]liveGapRow{}
	fset := token.NewFileSet()
	dir := filepath.Join(repo, "internal", "ensigncycle")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data := mustRead(t, path)
		preamble := strings.SplitN(string(data), "package ", 2)[0]
		if !strings.Contains(preamble, "//go:build live") {
			continue
		}
		file, err := parser.ParseFile(fset, path, data, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		for _, declaration := range file.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok || function.Doc == nil || strings.HasPrefix(function.Name.Name, "TestLiveCommon") || !strings.Contains(commentText(function.Doc), "spacedock:live-proof") {
				continue
			}
			var gaps []liveGapRow
			ast.Inspect(function.Body, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok || exprName(call.Fun) != "liveXFail" {
					return true
				}
				gap, err := parseLiveGap(call, targets)
				if err != nil {
					t.Fatalf("%s: %v", function.Name.Name, err)
				}
				gaps = append(gaps, gap)
				return true
			})
			if len(gaps) > 0 {
				result[function.Name.Name] = gaps
			}
		}
	}
	return result
}

func reconcileRegisteredLiveTests(t *testing.T, repo, registryPath string) {
	registry := string(mustRead(t, registryPath))
	registered := map[string]bool{}
	for _, pattern := range []*regexp.Regexp{registryEntryPoint, registryTestHeading} {
		for _, match := range pattern.FindAllStringSubmatch(registry, -1) {
			registered[match[1]] = true
		}
	}

	found := readLiveTestNames(t, repo)
	scheduledLanes := observeLiveSchedule(t, repo, string(mustRead(t, filepath.Join(repo, "internal/ensigncycle/scheduled_live_test.go"))))
	workflow := string(mustRead(t, filepath.Join(repo, ".github", "workflows", "runtime-live-e2e.yml")))
	for name := range found {
		if err := liveTestCoverageError(name, registry, workflow, scheduledLanes); err != nil {
			t.Error(err)
		}
		if !registered[name] {
			t.Errorf("live-tagged test %q is unregistered", name)
		}
	}

	for test := range registered {
		if !found[test] {
			t.Errorf("registered live test %q has no live-tagged declaration", test)
		}
	}
}

func readLiveTestNames(t *testing.T, repo string) map[string]bool {
	t.Helper()
	fset := token.NewFileSet()
	found := map[string]bool{}
	dir := filepath.Join(repo, "internal", "ensigncycle")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data := mustRead(t, path)
		preamble := strings.SplitN(string(data), "package ", 2)[0]
		if !strings.Contains(preamble, "//go:build live") {
			continue
		}
		file, err := parser.ParseFile(fset, path, data, 0)
		if err != nil {
			t.Fatal(err)
		}

		for _, declaration := range file.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok || !strings.HasPrefix(function.Name.Name, "Test") {
				continue
			}
			found[function.Name.Name] = true
		}
	}
	return found
}

// Execute the shipped selection/slot loop; only model-driving callables become stubs.
func observeLiveSchedule(t *testing.T, repo, source string) map[string]map[string]int {
	t.Helper()
	root := t.TempDir()
	callbacks := "package ensigncycle\nimport \"testing\"\n"
	for name := range readLiveTestNames(t, repo) {
		if name != "TestLiveScheduled" {
			callbacks += fmt.Sprintf("func %s(t *testing.T) { t.Log(\"SCHEDULED %s\") }\n", name, name)
		}
	}
	for name, data := range map[string]string{
		"scheduled_live_test.go": source,
		"live_schedule_test.go":  string(mustRead(t, filepath.Join(repo, "internal/ensigncycle/live_schedule_test.go"))),
		"callbacks_test.go":      callbacks,
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	observed := map[string]map[string]int{}
	for _, runtime := range []string{"claude", "codex"} {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		cmd := exec.CommandContext(ctx, "go", "test", "-json", "-tags", "live", "-run", "^TestLiveScheduled$", "-count=1", filepath.Join(root, "scheduled_live_test.go"), filepath.Join(root, "live_schedule_test.go"), filepath.Join(root, "callbacks_test.go"))
		cmd.Env = append(os.Environ(), "SPACEDOCK_LIVE_RUNTIME="+runtime)
		start := time.Now()
		output, err := cmd.CombinedOutput()
		cancel()
		if err != nil {
			t.Fatalf("offline %s scheduler: %v\n%s", runtime, err, output)
		}
		t.Logf("offline %s scheduler: %s", runtime, time.Since(start))
		for _, line := range strings.Split(string(output), "\n") {
			var event struct{ Action, Output string }
			if json.Unmarshal([]byte(line), &event) != nil || event.Action != "output" {
				continue
			}
			_, name, ok := strings.Cut(event.Output, "SCHEDULED ")
			if !ok {
				continue
			}
			name = strings.TrimSpace(name)
			if observed[name] == nil {
				observed[name] = map[string]int{}
			}
			observed[name][runtime+"-live"]++
		}
	}
	return observed
}

func readActiveEntityIDs(t *testing.T, stateDir string) map[string]bool {
	t.Helper()
	active := map[string]bool{}
	err := filepath.WalkDir(stateDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		idMatch, statusMatch := frontmatterID.FindSubmatch(data), frontmatterStatus.FindSubmatch(data)
		if idMatch == nil || statusMatch == nil {
			return nil
		}
		status := string(statusMatch[1])
		archived := strings.Contains(filepath.ToSlash(path), "/_archive/")
		active[string(idMatch[1])] = !archived && status != "done" && status != "rejected"
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return active
}

// commonRunShape is the run-defining subset of a common-suite command that the
// workflow and the local guide must agree on per runtime: timeout, parallelism,
// and fail-fast. The reporting wrapper (gotestsum vs `go test -v`) legitimately
// differs between the two sources and is not part of the shape.
type commonRunShape struct {
	timeout, parallel string
	failfast          bool
}

func commonCommandRunShape(t *testing.T, command string) commonRunShape {
	t.Helper()
	fields := strings.Fields(command)
	var shape commonRunShape
	for i, field := range fields {
		switch field {
		case "-timeout":
			if i+1 < len(fields) {
				shape.timeout = fields[i+1]
			}
		case "-parallel":
			if i+1 < len(fields) {
				shape.parallel = fields[i+1]
			}
		case "-failfast":
			shape.failfast = true
		}
	}
	if shape.timeout == "" {
		t.Fatalf("common-suite command carries no -timeout value: %s", command)
	}
	return shape
}

func docsLiveCommonCommand(t *testing.T, docs, runtime string) string {
	t.Helper()
	prefix := "SPACEDOCK_LIVE_RUNTIME=" + runtime + " go test "
	var matches []string
	for _, line := range strings.Split(docs, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefix) && strings.Contains(line, liveSuiteSelector(runtime)) {
			matches = append(matches, line)
		}
	}
	if len(matches) != 1 {
		t.Fatalf("local guide has %d %s common-journey commands, want exactly 1", len(matches), runtime)
	}
	return matches[0]
}

// TestRuntimeLiveCommonSuiteTimeouts binds the local guide's per-runtime
// common-suite command to the workflow's: each command is extracted from its
// own file and the two are compared on run shape. No expected command text
// lives in this test, so the sources can only agree with each other, not with
// a third hand-authored copy.
func TestRuntimeLiveCommonSuiteTimeouts(t *testing.T) {
	repo := repoRoot(t)
	workflow := string(mustRead(t, filepath.Join(repo, ".github", "workflows", "runtime-live-e2e.yml")))
	docs := string(mustRead(t, filepath.Join(repo, "docs", "runtime-live-ci.md")))
	for _, runtime := range []string{"claude", "codex", "pi"} {
		workflowShape := commonCommandRunShape(t, runtimeLiveCommonCommand(t, workflow, runtime))
		docsShape := commonCommandRunShape(t, docsLiveCommonCommand(t, docs, runtime))
		if workflowShape != docsShape {
			t.Errorf("%s common-suite run shape drift: workflow %+v, docs %+v", runtime, workflowShape, docsShape)
		}
	}
}

func TestRuntimeLiveCommonFailFastPolicy(t *testing.T) {
	workflow := string(mustRead(t, filepath.Join(repoRoot(t), ".github", "workflows", "runtime-live-e2e.yml")))
	for _, runtime := range []string{"claude", "codex", "pi"} {
		command := runtimeLiveCommonCommand(t, workflow, runtime)
		if strings.Contains(command, " -failfast") {
			t.Errorf("%s common journeys must all run before the job reports failure", runtime)
		}
	}
}

func runtimeLiveCommonCommand(t *testing.T, workflow, runtime string) string {
	t.Helper()
	prefix := "SPACEDOCK_LIVE_RUNTIME=" + runtime + " gotestsum "
	for _, line := range strings.Split(workflow, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefix) && strings.Contains(line, liveSuiteSelector(runtime)) {
			return line
		}
	}
	t.Fatalf("workflow has no %s common-journey command", runtime)
	return ""
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
		t.Fatalf("%s gap metadata must be nil or a composite literal", fn.Name.Name)
	}
	for _, element := range list.Elts {
		gap, err := parseLiveGap(element, targets)
		if err != nil {
			t.Fatalf("%s: %v", fn.Name.Name, err)
		}
		got.gaps = append(got.gaps, gap)
	}
	return got
}

func parseLiveGap(expr ast.Expr, targets map[string]bool) (liveGapRow, error) {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return liveGapRow{}, fmt.Errorf("gap element must be liveTODO or liveXFail")
	}
	kind := strings.TrimPrefix(exprName(call.Fun), "live")
	wantArgs := map[string]int{"TODO": 2, "XFail": 2}[kind]
	if wantArgs == 0 || len(call.Args) != wantArgs {
		return liveGapRow{}, fmt.Errorf("malformed live%s", kind)
	}
	values := make([]string, wantArgs)
	for i, arg := range call.Args {
		literal, ok := arg.(*ast.BasicLit)
		if !ok || literal.Kind != token.STRING {
			return liveGapRow{}, fmt.Errorf("gap argument is not a string literal")
		}
		values[i], _ = strconv.Unquote(literal.Value)
	}
	if !targets[values[0]] || !ownerID.MatchString(values[1]) {
		return liveGapRow{}, fmt.Errorf("malformed live%s binding", kind)
	}
	return liveGapRow{strings.ToLower(kind), values[0], values[1]}, nil
}

func TestRuntimeLiveGapBindingValidation(t *testing.T) {
	targets := map[string]bool{"codex": true}
	for source, valid := range map[string]bool{
		`liveTODO("codex", "98aa776adg66gn823a8gamdq")`:                                          true,
		`liveXFail("codex", "98aa776adg66gn823a8gamdq")`:                                         true,
		`liveTODO("codex", "98aa776adg66gn823a8gamdq", "code")`:                                  false,
		`liveXFail("codex", "bad-owner")`:                                                        false,
		`liveXFail("unknown", "98aa776adg66gn823a8gamdq")`:                                       false,
		`liveXFail("codex", "98aa776adg66gn823a8gamdq", "implementation-worker-not-dispatched")`: false,
	} {
		expr, err := parser.ParseExpr(source)
		if err != nil {
			t.Fatal(err)
		}
		_, err = parseLiveGap(expr, targets)
		if (err == nil) != valid {
			t.Errorf("parseLiveGap(%s) error = %v, valid=%t", source, err, valid)
		}
	}
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

func liveSuiteSelector(runtime string) string {
	if runtime == "pi" {
		return "-run '^TestLiveCommon'"
	}
	return "-run '^TestLiveScheduled$'"
}

func TestLiveCoverageRequiresSelectionOrExplicitExemption(t *testing.T) {
	registry := "## Common journeys\n### `sample`\n- **Entry point:** `TestLiveCommonSample`\n"
	workflow := "jobs:\n"
	for _, runtime := range []string{"claude", "codex", "pi"} {
		workflow += fmt.Sprintf("  %s-live:\n    steps:\n      - run: go test -tags live %s ./internal/ensigncycle\n", runtime, liveSuiteSelector(runtime))
	}
	scheduled := map[string]map[string]int{"TestLiveCommonSample": {"claude-live": 1, "codex-live": 1}}
	for _, tc := range []struct {
		name, registry, workflow string
		valid                    bool
	}{
		{"selected", registry, workflow, true},
		{"missing-membership", registry, workflow, false},
		{"duplicate-selection", registry, workflow + "      - run: go test -tags live -run '^TestLiveCommon' ./internal/ensigncycle\n", false},
		{"commented-selector", registry, strings.ReplaceAll(workflow, "- run: go test", "- run: |\n          # go test"), false},
		{"missing-selector", registry, strings.ReplaceAll(workflow, "^TestLiveScheduled$", "^Other$"), false},
		{"unclassified", strings.ReplaceAll(registry, "Common journeys", "Targeted implementation proofs"), workflow, false},
		{"explicit-exemption", "## Non-gating live experiments\n### `TestLiveCommonSample`\n- **Reason unselected:** Measures an experimental host.\n", "", true},
		{"blank-reason", "## Non-gating live experiments\n### `TestLiveCommonSample`\n- **Reason unselected:**   \n", "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			membership := scheduled
			if tc.name == "missing-membership" {
				membership = nil
			}
			if err := liveTestCoverageError("TestLiveCommonSample", tc.registry, tc.workflow, membership); (err == nil) != tc.valid {
				t.Fatalf("coverage error = %v, valid=%t", err, tc.valid)
			}
		})
	}
}

var registryLane = regexp.MustCompile("(?m)^- \\*\\*Lane:\\*\\* `([^`]+)`$")

func liveRegistryEntry(name, registry string) (string, string) {
	section, blockSection, block, current := "", "", "", ""
	for _, line := range strings.Split(registry, "\n") {
		if strings.HasPrefix(line, "## ") {
			section = strings.TrimPrefix(line, "## ")
			current = ""
		}
		if registryHeading.MatchString(line) {
			current = line + "\n"
		} else {
			current += line + "\n"
		}
		heading := registryTestHeading.FindStringSubmatch(current)
		entry := registryEntryPoint.FindStringSubmatch(current)
		if (heading != nil && heading[1] == name) || (entry != nil && entry[1] == name) {
			block, blockSection = current, section
		}
	}
	return blockSection, block
}

func liveTestCoverageError(name, registry, workflow string, scheduled map[string]map[string]int) error {
	blockSection, block := liveRegistryEntry(name, registry)
	lanes := []string{}
	switch blockSection {
	case "Common journeys":
		lanes = []string{"claude-live", "codex-live", "pi-live"}
	case "Suite orchestration":
		lanes = []string{"claude-live", "codex-live"}
	case "Runtime-specific live proofs":
		match := registryLane.FindStringSubmatch(block)
		if match != nil {
			lanes = []string{match[1]}
		}
	case "Non-gating live experiments":
		if regexp.MustCompile(`(?m)^- \*\*Reason unselected:\*\*[^\S\n]*\S[^\n]*$`).MatchString(block) {
			return nil
		}
		return fmt.Errorf("%s has no explicit exemption reason", name)
	}
	if len(lanes) == 0 {
		return fmt.Errorf("%s has no recognized live coverage classification", name)
	}
	var document struct {
		Jobs map[string]struct{ Steps []struct{ Run string } }
	}
	if err := yaml.Unmarshal([]byte(workflow), &document); err != nil {
		return err
	}
	selections := map[string]int{}
	for lane, job := range document.Jobs {
		for _, step := range job.Steps {
			for _, line := range strings.Split(step.Run, "\n") {
				if strings.HasPrefix(strings.TrimSpace(line), "#") || !strings.Contains(line, "-tags live") || !strings.Contains(line, "./internal/ensigncycle") {
					continue
				}
				match := regexp.MustCompile(`-run\s+['"]?([^'"\s]+)`).FindStringSubmatch(line)
				if match == nil {
					return fmt.Errorf("live command has no explicit selector: %s", line)
				}
				selector, err := regexp.Compile(match[1])
				if err != nil {
					return err
				}
				if selector.MatchString(name) {
					selections[lane]++
				}
				if name != "TestLiveScheduled" && selector.MatchString("TestLiveScheduled") {
					selections[lane] += scheduled[name][lane]
				}
			}
		}
	}

	for _, lane := range lanes {
		if selections[lane] != 1 {
			return fmt.Errorf("%s selected %d times in %s, want 1", name, selections[lane], lane)
		}
	}
	return nil
}

func TestLiveSelectionMutationsAreObservable(t *testing.T) {
	repo := repoRoot(t)
	source := string(mustRead(t, filepath.Join(repo, "internal/ensigncycle/scheduled_live_test.go")))
	source = strings.Replace(source, "for _, row := range rows {", "for _, row := range rows {\n if runtime == \"codex\" && row.name == \"TestLiveCommonSameStageRevision\" { continue }", 1)
	observed := observeLiveSchedule(t, repo, source)
	registry := string(mustRead(t, filepath.Join(repo, "docs/runtime-live-ci-registry.md")))
	workflow := string(mustRead(t, filepath.Join(repo, ".github/workflows/runtime-live-e2e.yml")))
	if got := observed["TestLiveCommonSameStageRevision"]; got["claude-live"] != 1 || got["codex-live"] != 0 {
		t.Fatalf("observed calls: %v", got)
	}
	if err := liveTestCoverageError("TestLiveCommonSameStageRevision", registry, workflow, observed); err == nil {
		t.Fatal("runtime omission passed coverage")
	}
	if err := liveTestCoverageError("TestLiveCommonFiling", registry, workflow, observed); err != nil {
		t.Fatal(err)
	}
	command := "SPACEDOCK_LIVE_RUNTIME=claude gotestsum"
	if strings.Count(workflow, command) != 1 {
		t.Fatal("expected one actual Claude command")
	}
	workflow = strings.Replace(workflow, command, "# "+command, 1)
	if err := liveTestCoverageError("TestLiveCommonFiling", registry, workflow, observed); err == nil || !strings.Contains(err.Error(), "selected 0 times in claude-live") {
		t.Fatalf("commented actual command: %v", err)
	}
}
