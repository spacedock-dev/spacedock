package contractlint

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/release"
)

type liveGapRow struct{ kind, target, owner string }

var ownerID = regexp.MustCompile(`^[a-z0-9]{24}$`)

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
	for _, runtime := range []string{"claude", "codex"} {
		command := runtimeLiveCommonCommand(t, workflow, runtime)
		if strings.Contains(command, " -failfast") {
			t.Errorf("%s common journeys must all run before the job reports failure", runtime)
		}
	}
	if command := runtimeLiveCommonCommand(t, workflow, "pi"); !strings.Contains(command, " -failfast") {
		t.Error("Pi common journeys must retain -failfast")
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
	// A blank owner (unowned, pending analysis) is permitted only on claude-haiku.
	if !targets[values[0]] || !ownerID.MatchString(values[1]) && (values[1] != "" || values[0] != "claude-haiku") {
		return liveGapRow{}, fmt.Errorf("malformed live%s binding", kind)
	}
	return liveGapRow{strings.ToLower(kind), values[0], values[1]}, nil
}

func TestRuntimeLiveGapBindingValidation(t *testing.T) {
	targets := map[string]bool{"codex": true, "claude-haiku": true, "claude-sonnet": true}
	for source, valid := range map[string]bool{
		`liveXFail("claude-haiku", "")`:                                                          true,
		`liveXFail("claude-sonnet", "")`:                                                         false,
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

// liveBinding is one live gap binding: the journey it covers, its target and owner.
type liveBinding struct{ journey, target, owner string }

// TestRuntimeLiveRegistryMatchesBindings reads every real gap binding in
// internal/ensigncycle and checks it through parseLiveGap: a known target, and a
// 24-character owner (blank only on claude-haiku). The bindings and the rows of
// docs/runtime-live-ci-registry.md must then match in both directions.
func TestRuntimeLiveRegistryMatchesBindings(t *testing.T) {
	repo := repoRoot(t)
	code := ensigncycleLiveBindings(t, filepath.Join(repo, "internal", "ensigncycle"))
	doc := map[liveBinding]bool{}
	for _, line := range strings.Split(string(mustRead(t, filepath.Join(repo, "docs", "runtime-live-ci-registry.md"))), "\n") {
		cells := strings.Split(strings.Trim(line, "| "), " | ")
		if !strings.HasPrefix(line, "| ") || len(cells) != 4 || cells[0] == "Journey" {
			continue
		}
		owner := strings.Trim(cells[2], "`")
		if strings.HasPrefix(owner, "unowned") {
			owner = ""
		}
		doc[liveBinding{strings.Trim(strings.Fields(cells[0])[0], "`"), strings.Trim(cells[1], "`"), owner}] = true
	}
	if len(code) == 0 || len(doc) == 0 {
		t.Fatalf("found %d bindings and %d registry rows, want both non-empty", len(code), len(doc))
	}
	for binding := range code {
		if !doc[binding] {
			t.Errorf("binding %+v has no registry row", binding)
		}
	}
	for row := range doc {
		if !code[row] {
			t.Errorf("registry row %+v has no binding in internal/ensigncycle", row)
		}
	}
}

// ensigncycleLiveBindings collects each liveJourney gaps argument (inline or
// through a package var), keyed by the journey id, and each gap built in a
// non-Test function such as breakGlassSelectedBareGap, keyed by that function's
// name. Other gap literals inside Test functions are grading fixtures.
func ensigncycleLiveBindings(t *testing.T, dir string) map[liveBinding]bool {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil || len(paths) == 0 {
		t.Fatalf("glob %s: %v", dir, err)
	}
	targets := map[string]bool{"codex": true, "pi": true}
	for key := range release.LiveModels {
		if role, ok := strings.CutPrefix(key, "claude."); ok && !strings.Contains(role, ".") {
			targets["claude-"+role] = true
		}
	}
	fset, vars, funcs, found := token.NewFileSet(), map[string]ast.Expr{}, []*ast.FuncDecl{}, map[liveBinding]bool{}
	for _, path := range paths {
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				funcs = append(funcs, fn)
			} else if gen, ok := decl.(*ast.GenDecl); ok {
				for _, spec := range gen.Specs {
					if value, ok := spec.(*ast.ValueSpec); ok && len(value.Names) == 1 && len(value.Values) == 1 {
						vars[value.Names[0].Name] = value.Values[0]
					}
				}
			}
		}
	}
	collect := func(journey string, node ast.Node) {
		ast.Inspect(node, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || exprName(call.Fun) != "liveXFail" && exprName(call.Fun) != "liveTODO" {
				return true
			}
			gap, err := parseLiveGap(call, targets)
			if err != nil {
				t.Errorf("%s binding at %s: %v", journey, fset.Position(call.Pos()), err)
			}
			found[liveBinding{journey, gap.target, gap.owner}] = true
			return false
		})
	}
	for _, fn := range funcs {
		if !strings.HasPrefix(fn.Name.Name, "Test") {
			collect(fn.Name.Name, fn.Body)
			continue
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || exprName(call.Fun) != "liveJourney" || len(call.Args) < 5 {
				return true
			}
			id, _ := call.Args[1].(*ast.BasicLit)
			if id == nil {
				t.Fatalf("liveJourney at %s has no literal id", fset.Position(call.Pos()))
			}
			journey, _ := strconv.Unquote(id.Value)
			if ident, ok := call.Args[4].(*ast.Ident); ok && vars[ident.Name] != nil {
				collect(journey, vars[ident.Name])
			} else {
				collect(journey, call.Args[4])
			}
			return false
		})
	}
	return found
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

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func liveSuiteSelector(runtime string) string {
	if runtime == "pi" {
		return "-run '^TestLiveCommon'"
	}
	return "-run '^TestLiveScheduled$'"
}
