package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestUnknownCommandRunsPathExtension(t *testing.T) {
	dir := t.TempDir()
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$0.args\"\nexit 7\n"
	if err := os.WriteFile(filepath.Join(dir, "spacedock-hello"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"hello", "--version", "--mode", "feedback", "file.md"}, &stdout, &stderr)

	if code != 7 {
		t.Fatalf("Run returned %d, want the extension's exit 7; stderr=%q", code, stderr.String())
	}
	got, err := os.ReadFile(filepath.Join(dir, "spacedock-hello.args"))
	if err != nil {
		t.Fatalf("extension did not run: %v", err)
	}
	if want := "--version\n--mode\nfeedback\nfile.md\n"; string(got) != want {
		t.Fatalf("extension argv = %q, want %q", got, want)
	}
}

func TestExtensionArgvLeavesBuiltinsAndFlagsToCobra(t *testing.T) {
	root := newRootCommand(context.Background(), nil, nil, "", nil, &bytes.Buffer{}, &bytes.Buffer{}, nil, nil)
	found := func(name string) (string, error) { return "/bin/" + name, nil }
	missing := func(string) (string, error) { return "", errors.New("not found") }

	for _, args := range [][]string{{"status"}, {"help"}, {"--version"}, {"../x"}, {}} {
		if argv := extensionArgv(root, args, found); argv != nil {
			t.Errorf("extensionArgv(%q) = %q, want cobra to keep it", args, argv)
		}
	}
	if argv := extensionArgv(root, []string{"review"}, missing); argv != nil {
		t.Errorf("extensionArgv without a PATH match = %q, want nil", argv)
	}
	want := []string{"/bin/spacedock-review", "doc.md"}
	if argv := extensionArgv(root, []string{"review", "doc.md"}, found); !reflect.DeepEqual(argv, want) {
		t.Errorf("extensionArgv(review doc.md) = %q, want %q", argv, want)
	}
}
