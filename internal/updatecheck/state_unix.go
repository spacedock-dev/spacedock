// ABOUTME: Unix daily-marker creation: checked owner-only directories and an
// ABOUTME: exclusive, no-follow marker open.
//go:build unix

package updatecheck

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
)

var errStateNotPrivate = errors.New("update-check state directory is not private")

func openMarker(base, day string) (*os.File, error) {
	if err := os.MkdirAll(base, 0o700); err != nil {
		return nil, err
	}
	if info, err := os.Lstat(base); err != nil || !privateDirectory(info, false) {
		return nil, errStateNotPrivate
	}
	dir := base
	for _, component := range []string{"spacedock", "version-checks"} {
		dir = filepath.Join(dir, component)
		if err := os.Mkdir(dir, 0o700); err != nil && !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		if info, err := os.Lstat(dir); err != nil || !privateDirectory(info, true) {
			return nil, errStateNotPrivate
		}
	}
	return os.OpenFile(filepath.Join(dir, day), os.O_WRONLY|os.O_CREATE|os.O_EXCL|syscall.O_NOFOLLOW, 0o600)
}

// privateDirectory accepts a real directory owned by this user. The state root
// may be group-writable (umask 002 makes ~/.local/state 0775, and refusing it
// would silently disable the check) but not world-writable; Spacedock's own
// directories must be exactly 0700.
func privateDirectory(info os.FileInfo, ownerOnly bool) bool {
	if info == nil || !info.IsDir() { // from Lstat, so a symlink is not a directory
		return false
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok || uint64(stat.Uid) != uint64(os.Geteuid()) {
		return false
	}
	if ownerOnly {
		return info.Mode().Perm() == 0o700
	}
	return info.Mode().Perm()&0o002 == 0
}
