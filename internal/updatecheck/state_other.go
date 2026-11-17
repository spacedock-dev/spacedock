// ABOUTME: Non-unix daily-marker shim: without owner checks and no-follow opens
// ABOUTME: the claim fails closed, so no lookup is made.
//go:build !unix

package updatecheck

import "os"

func openMarker(base, day string) (*os.File, error) { return nil, os.ErrPermission }
