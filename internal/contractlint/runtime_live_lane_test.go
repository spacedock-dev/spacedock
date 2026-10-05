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
