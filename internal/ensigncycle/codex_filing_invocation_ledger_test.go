package ensigncycle

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type codexFilingInvocation struct {
	argv     []string
	exitCode int
}

type codexFilingInvocationLedger struct{ dir, shimDir, real string }

func newCodexFilingInvocationLedger(t testing.TB, real string) codexFilingInvocationLedger {
	root := t.TempDir()
	l := codexFilingInvocationLedger{filepath.Join(root, "ledger"), filepath.Join(root, "bin"), real}
	for _, dir := range []string{l.dir, l.shimDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	const shim = `#!/bin/sh
set -u
r=$(mktemp "$SPACEDOCK_CODEX_FILING_LEDGER_DIR/.pending.XXXXXX") || exit 125
trap 'rm -f "$r"' EXIT HUP INT TERM
printf 'spacedock\0argc\0%s\0' "$#" > "$r" || exit 125
for arg in "$@"; do printf '%s\0' "$arg" >> "$r" || exit 125; done
"$SPACEDOCK_CODEX_FILING_REAL" "$@"
status=$?
printf 'exit\0%s\0' "$status" >> "$r" || exit 125
mv "$r" "$SPACEDOCK_CODEX_FILING_LEDGER_DIR/invocation.$$.${r##*.}" || exit 125
trap - EXIT HUP INT TERM
exit "$status"
`
	if err := os.WriteFile(filepath.Join(l.shimDir, "spacedock"), []byte(shim), 0o755); err != nil {
		t.Fatal(err)
	}
	return l
}

func (l codexFilingInvocationLedger) instrumentEnv(env []string) []string {
	env = replaceEnvValue(env, "SPACEDOCK_BIN", filepath.Join(l.shimDir, "spacedock"))
	env = replaceEnvValue(env, "SPACEDOCK_CODEX_FILING_LEDGER_DIR", l.dir)
	env = replaceEnvValue(env, "SPACEDOCK_CODEX_FILING_REAL", l.real)
	path, _ := envValue(env, "PATH")
	return replaceEnvValue(env, "PATH", l.shimDir+string(os.PathListSeparator)+path)
}

func (l codexFilingInvocationLedger) read() ([]codexFilingInvocation, error) {
	paths, _ := filepath.Glob(filepath.Join(l.dir, "invocation.*"))
	var out []codexFilingInvocation
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		fields := bytes.Split(data, []byte{0})
		if len(fields) < 6 || len(fields[len(fields)-1]) != 0 || string(fields[0]) != "spacedock" || string(fields[1]) != "argc" {
			return nil, fmt.Errorf("invalid Codex filing ledger record %s", path)
		}
		argc, argcErr := strconv.Atoi(string(fields[2]))
		exitCode, exitErr := strconv.Atoi(string(fields[len(fields)-2]))
		if argcErr != nil || exitErr != nil || argc < 0 || exitCode < 0 || exitCode > 255 || len(fields) != argc+6 || string(fields[3+argc]) != "exit" {
			return nil, fmt.Errorf("invalid Codex filing ledger fields in %s", path)
		}
		invocation := codexFilingInvocation{exitCode: exitCode}
		for _, field := range fields[3 : 3+argc] {
			invocation.argv = append(invocation.argv, string(field))
		}
		out = append(out, invocation)
	}
	return out, nil
}

func assertCodexFilingInvocations(invocations []codexFilingInvocation, slug string) error {
	filed := false
	for _, invocation := range invocations {
		args := invocation.argv
		if invocation.exitCode != 0 || len(args) == 0 {
			continue
		}
		if args[0] == "status" && slices.Contains(args[1:], "--next-id") {
			return fmt.Errorf("filing previewed --next-id instead of using the atomic new path")
		}
		alias := slices.Index(args[1:], "--new")
		filed = filed || len(args) > 1 && args[0] == "new" && args[1] == slug || args[0] == "status" && alias >= 0 && alias+2 < len(args) && args[alias+2] == slug
	}
	if !filed {
		return fmt.Errorf("filing ledger has no successful spacedock new %s invocation", slug)
	}
	return nil
}

func replaceEnvValue(env []string, key, value string) []string {
	out := slices.DeleteFunc(append([]string(nil), env...), func(entry string) bool { return strings.HasPrefix(entry, key+"=") })
	return append(out, key+"="+value)
}
