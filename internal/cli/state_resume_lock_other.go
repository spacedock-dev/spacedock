// ABOUTME: In-process fallback for platforms without Unix advisory locks.
//go:build !unix

package cli

import "sync"

var stateResumeMu sync.Mutex

func withStateResumeLock(_ string, fn func() int) (int, error) {
	stateResumeMu.Lock()
	defer stateResumeMu.Unlock()
	return fn(), nil
}
