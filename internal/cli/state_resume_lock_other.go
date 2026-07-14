// ABOUTME: In-process fallback for platforms without Unix advisory locks.
//go:build !unix

package cli

func withStateResumeLock(_ string, fn func() int) (int, error) {
	return unsupportedStateResumeLock("", fn)
}
