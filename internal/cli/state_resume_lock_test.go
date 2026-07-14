// ABOUTME: Separate-process proof for resume-lock serialization and fail-closed fallback.
//go:build unix

package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestStateResumeLockProcessHelper(t *testing.T) {
	mode := os.Getenv("SPACEDOCK_LOCK_HELPER")
	if mode == "" {
		return
	}
	callbackPath := os.Getenv("SPACEDOCK_LOCK_CALLBACK")
	callback := func() int {
		activePath := os.Getenv("SPACEDOCK_LOCK_ACTIVE")
		if activePath != "" {
			active, err := os.OpenFile(activePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
			if err != nil {
				t.Fatalf("cross-process callbacks overlapped: %v", err)
			}
			active.Close()
			defer os.Remove(activePath)
		}
		f, err := os.OpenFile(callbackPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		actor := os.Getenv("SPACEDOCK_LOCK_ACTOR")
		fmt.Fprintln(f, actor+"-start")
		time.Sleep(120 * time.Millisecond)
		fmt.Fprintln(f, actor+"-end")
		return 0
	}

	switch mode {
	case "unsupported":
		if _, err := unsupportedStateResumeLock("", callback); !errors.Is(err, errStateResumeLockUnsupported) {
			t.Fatalf("unsupported lock error=%v", err)
		}
	case "unix":
		if err := os.WriteFile(os.Getenv("SPACEDOCK_LOCK_READY"), nil, 0o600); err != nil {
			t.Fatal(err)
		}
		gate := os.Getenv("SPACEDOCK_LOCK_GATE")
		deadline := time.Now().Add(5 * time.Second)
		for {
			if _, err := os.Stat(gate); err == nil {
				break
			}
			if time.Now().After(deadline) {
				t.Fatal("timed out waiting for process gate")
			}
			time.Sleep(5 * time.Millisecond)
		}
		if code, err := withStateResumeLock(os.Getenv("SPACEDOCK_LOCK_REPO"), os.Getenv("SPACEDOCK_LOCK_STATE"), func(bool) int { return callback() }); err != nil || code != 0 {
			t.Fatalf("withStateResumeLock code=%d err=%v", code, err)
		}
	default:
		t.Fatalf("unknown helper mode %q", mode)
	}
}

func startResumeLockHelper(t *testing.T, env ...string) *exec.Cmd {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestStateResumeLockProcessHelper$")
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	return cmd
}

func waitForFiles(t *testing.T, paths ...string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		all := true
		for _, path := range paths {
			if _, err := os.Stat(path); err != nil {
				all = false
			}
		}
		if all {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for helper readiness: %v", paths)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func TestStateResumeLockSerializesSeparateProcesses(t *testing.T) {
	root := t.TempDir()
	git(t, root, "init", "-q")
	callback := filepath.Join(t.TempDir(), "callbacks")
	active := filepath.Join(t.TempDir(), "active")
	gate := filepath.Join(t.TempDir(), "gate")
	readyA := filepath.Join(t.TempDir(), "ready-a")
	readyB := filepath.Join(t.TempDir(), "ready-b")
	common := []string{
		"SPACEDOCK_LOCK_HELPER=unix",
		"SPACEDOCK_LOCK_REPO=" + root,
		"SPACEDOCK_LOCK_STATE=" + filepath.Join(root, "state"),
		"SPACEDOCK_LOCK_CALLBACK=" + callback,
		"SPACEDOCK_LOCK_ACTIVE=" + active,
		"SPACEDOCK_LOCK_GATE=" + gate,
	}
	cmdA := startResumeLockHelper(t, append(common, "SPACEDOCK_LOCK_ACTOR=a", "SPACEDOCK_LOCK_READY="+readyA)...)
	cmdB := startResumeLockHelper(t, append(common, "SPACEDOCK_LOCK_ACTOR=b", "SPACEDOCK_LOCK_READY="+readyB)...)
	waitForFiles(t, readyA, readyB)
	if err := os.WriteFile(gate, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := cmdA.Wait(); err != nil {
		t.Fatal(err)
	}
	if err := cmdB.Wait(); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(callback)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Fields(string(body))
	wantA := []string{"a-start", "a-end", "b-start", "b-end"}
	wantB := []string{"b-start", "b-end", "a-start", "a-end"}
	if fmt.Sprint(got) != fmt.Sprint(wantA) && fmt.Sprint(got) != fmt.Sprint(wantB) {
		t.Fatalf("separate-process callbacks overlapped: %v", got)
	}
}

func TestUnsupportedStateResumeLockFailsClosedAcrossProcesses(t *testing.T) {
	callback := filepath.Join(t.TempDir(), "must-not-run")
	cmdA := startResumeLockHelper(t, "SPACEDOCK_LOCK_HELPER=unsupported", "SPACEDOCK_LOCK_CALLBACK="+callback, "SPACEDOCK_LOCK_ACTOR=a")
	cmdB := startResumeLockHelper(t, "SPACEDOCK_LOCK_HELPER=unsupported", "SPACEDOCK_LOCK_CALLBACK="+callback, "SPACEDOCK_LOCK_ACTOR=b")
	if err := cmdA.Wait(); err != nil {
		t.Fatal(err)
	}
	if err := cmdB.Wait(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(callback); !os.IsNotExist(err) {
		t.Fatalf("unsupported-platform callback ran: %v", err)
	}
}
