// ABOUTME: In-process fallback for platforms without Unix advisory locks.
//go:build !unix

package cli

func withStateResumeLock(_ string, _ string, fn func(bool) int) (int, error) {
	return unsupportedStateResumeLock("", func() int { return fn(false) })
}
