package ensigncycle

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

type testInvocation struct {
	tool string
	args []string
}

type testInvocationLedger struct {
	dir     string
	shimDir string
	real    map[string]string
}

func newTestInvocationLedger(t *testing.T, realSpacedock string) testInvocationLedger {
	t.Helper()
	root := t.TempDir()
	ledger := testInvocationLedger{
		dir:     filepath.Join(root, "ledger"),
		shimDir: filepath.Join(root, "bin"),
		real:    map[string]string{"spacedock": realSpacedock},
	}
	for _, tool := range []string{"jq", "python3", "go"} {
		if path, err := exec.LookPath(tool); err == nil {
			ledger.real[tool] = path
		}
	}
	for _, dir := range []string{ledger.dir, ledger.shimDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	shim := `#!/bin/sh
set -eu
tool=${0##*/}
record=$(mktemp "$SPACEDOCK_TEST_LEDGER_DIR/invocation.XXXXXX")
{
  printf '%s\0' "$tool"
  for arg in "$@"; do
    printf '%s\0' "$arg"
  done
} > "$record"
case "$tool" in
  spacedock) real=${SPACEDOCK_TEST_REAL_SPACEDOCK:-} ;;
  jq) real=${SPACEDOCK_TEST_REAL_JQ:-} ;;
  python3) real=${SPACEDOCK_TEST_REAL_PYTHON3:-} ;;
  go) real=${SPACEDOCK_TEST_REAL_GO:-} ;;
  *) exit 127 ;;
esac
[ -n "$real" ] || exit 127
exec "$real" "$@"
`
	for _, tool := range []string{"spacedock", "jq", "python3", "go"} {
		if err := os.WriteFile(filepath.Join(ledger.shimDir, tool), []byte(shim), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return ledger
}

func (l testInvocationLedger) instrumentEnv(env []string) []string {
	values := map[string]string{
		"SPACEDOCK_BIN":                 filepath.Join(l.shimDir, "spacedock"),
		"SPACEDOCK_TEST_LEDGER_DIR":     l.dir,
		"SPACEDOCK_TEST_REAL_SPACEDOCK": l.real["spacedock"],
		"SPACEDOCK_TEST_REAL_JQ":        l.real["jq"],
		"SPACEDOCK_TEST_REAL_PYTHON3":   l.real["python3"],
		"SPACEDOCK_TEST_REAL_GO":        l.real["go"],
	}
	out := append([]string(nil), env...)
	for key, value := range values {
		out = replaceEnvValue(out, key, value)
	}
	return prependEnvPath(out, l.shimDir)
}

func (l testInvocationLedger) read(t *testing.T) []testInvocation {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(l.dir, "invocation.*"))
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(paths)
	invocations := make([]testInvocation, 0, len(paths))
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		fields := bytes.Split(data, []byte{0})
		if len(fields) > 0 && len(fields[len(fields)-1]) == 0 {
			fields = fields[:len(fields)-1]
		}
		if len(fields) == 0 || len(fields[0]) == 0 {
			t.Fatalf("invocation ledger entry %s has no tool name", path)
		}
		invocation := testInvocation{tool: string(fields[0])}
		for _, field := range fields[1:] {
			invocation.args = append(invocation.args, string(field))
		}
		invocations = append(invocations, invocation)
	}
	return invocations
}

func writeInvocationLedgerArtifact(t *testing.T, artifactDir string, invocations []testInvocation) {
	t.Helper()
	var output strings.Builder
	for _, invocation := range invocations {
		fmt.Fprintf(&output, "%s", invocation.tool)
		for _, arg := range invocation.args {
			fmt.Fprintf(&output, "\t%q", arg)
		}
		output.WriteByte('\n')
	}
	if err := os.WriteFile(filepath.Join(artifactDir, "invocation-ledger.txt"), []byte(output.String()), 0o644); err != nil {
		t.Fatal(err)
	}
}

func replaceEnvValue(env []string, key, value string) []string {
	prefix := key + "="
	out := make([]string, 0, len(env)+1)
	for _, entry := range env {
		if !strings.HasPrefix(entry, prefix) {
			out = append(out, entry)
		}
	}
	return append(out, prefix+value)
}

func prependEnvPath(env []string, dir string) []string {
	path, _ := envValue(env, "PATH")
	if path == "" {
		return replaceEnvValue(env, "PATH", dir)
	}
	return replaceEnvValue(env, "PATH", dir+string(os.PathListSeparator)+path)
}

func invocationHasArg(invocation testInvocation, want string) bool {
	for _, arg := range invocation.args {
		if arg == want {
			return true
		}
	}
	return false
}

func invocationHasAdjacentArgs(invocation testInvocation, left, right string) bool {
	for i := 0; i+1 < len(invocation.args); i++ {
		if invocation.args[i] == left && invocation.args[i+1] == right {
			return true
		}
	}
	return false
}

func writeSuccessfulLedgerTarget(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "spacedock-real")
	if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func runLedgerShell(t *testing.T, ledger testInvocationLedger, command string) {
	t.Helper()
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Env = ledger.instrumentEnv(os.Environ())
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("ledger shell failed: %v\n%s", err, output)
	}
}

func TestInvocationLedgerRecordsExecutionNotCommandShapedNarration(t *testing.T) {
	ledger := newTestInvocationLedger(t, writeSuccessfulLedgerTarget(t))
	runLedgerShell(t, ledger, `echo 'spacedock status --boot --identify --json'; echo 'printf body | "$SPACEDOCK_BIN" new wire-the-thing'`)
	if got := ledger.read(t); len(got) != 0 {
		t.Fatalf("narrated commands reached invocation ledger: %#v", got)
	}

	runLedgerShell(t, ledger, `"$SPACEDOCK_BIN" status --boot --identify --json`)
	got := ledger.read(t)
	if len(got) != 1 || got[0].tool != "spacedock" || fmt.Sprint(got[0].args) != "[status --boot --identify --json]" {
		t.Fatalf("actual launcher invocation not recorded exactly: %#v", got)
	}
}
