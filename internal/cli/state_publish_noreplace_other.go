// ABOUTME: Fail-closed fallback where atomic no-replace rename is unavailable.
//go:build !darwin && !linux

package cli

import "fmt"

func renameNoReplace(_, _ string) error {
	return fmt.Errorf("atomic no-replace state publication is unsupported on this platform")
}
