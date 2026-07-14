// ABOUTME: Darwin atomic no-replace publication for converged state checkouts.
//go:build darwin

package cli

import "golang.org/x/sys/unix"

func renameNoReplace(oldPath, newPath string) error {
	return unix.RenamexNp(oldPath, newPath, unix.RENAME_EXCL)
}
