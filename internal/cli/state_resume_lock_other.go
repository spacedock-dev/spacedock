// ABOUTME: In-process fallback for platforms without Unix advisory locks.
//go:build !unix

package cli

import "sync"

var stateResumeMu sync.Mutex

func withStateResumeLock(_ string, fn func() (int, bool)) (int, bool, error) {
	stateResumeMu.Lock()
	defer stateResumeMu.Unlock()
	code, reached := fn()
	return code, reached, nil
}
