//go:build unix

package updatecheck

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

type foreignOwner struct{ os.FileInfo }

func (info foreignOwner) Sys() any {
	stat := *info.FileInfo.Sys().(*syscall.Stat_t)
	stat.Uid++
	return &stat
}

func TestClaimDayStateDirectories(t *testing.T) {
	now := time.Date(2026, 11, 17, 12, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name  string
		setup func(t *testing.T, root string)
		want  bool
	}{
		{"group-writable root", func(t *testing.T, root string) { chmod(t, root, 0o775) }, true},
		{"world-writable root", func(t *testing.T, root string) { chmod(t, root, 0o777) }, false},
		{"symlinked spacedock", func(t *testing.T, root string) { symlinkPrivate(t, filepath.Join(root, "spacedock")) }, false},
		{"symlinked version-checks", func(t *testing.T, root string) {
			mkdir(t, filepath.Join(root, "spacedock"))
			symlinkPrivate(t, filepath.Join(root, "spacedock", "version-checks"))
		}, false},
		{"shared version-checks", func(t *testing.T, root string) {
			mkdir(t, filepath.Join(root, "spacedock", "version-checks"))
			chmod(t, filepath.Join(root, "spacedock", "version-checks"), 0o777)
		}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			chmod(t, root, 0o700)
			tc.setup(t, root)
			got := ClaimDay(func(k string) string {
				return map[string]string{"XDG_STATE_HOME": root}[k]
			}, now)
			if got != tc.want {
				t.Fatalf("claimed=%v, want %v", got, tc.want)
			}
		})
	}
	info, err := os.Lstat(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if !privateDirectory(info, false) || privateDirectory(foreignOwner{info}, false) {
		t.Fatal("owner check: own directory refused or foreign owner accepted")
	}
}

func chmod(t *testing.T, path string, mode os.FileMode) {
	if err := os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
}

func mkdir(t *testing.T, path string) {
	if err := os.MkdirAll(path, 0o700); err != nil {
		t.Fatal(err)
	}
}

// symlinkPrivate links path to a fresh 0700 directory, so only the symlink
// check can refuse it, and fails the test if a marker lands in the target.
func symlinkPrivate(t *testing.T, path string) {
	target := t.TempDir()
	chmod(t, target, 0o700)
	if err := os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if entries, _ := os.ReadDir(target); len(entries) != 0 {
			t.Errorf("wrote through the symlink: %v", entries)
		}
	})
}
