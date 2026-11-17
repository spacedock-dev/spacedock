// ABOUTME: The once-per-UTC-day claim: an exclusive marker file under
// ABOUTME: ${XDG_STATE_HOME:-$HOME/.local/state}/spacedock/version-checks/.
package updatecheck

import (
	"io"
	"path/filepath"
	"time"
)

// ClaimDay records today's UTC date and reports whether this launch is the
// first to do so. Exclusive daily files also serialize concurrent launchers
// without a stale lock after a crash. A failed lookup leaves its claim in place.
func ClaimDay(getenv func(string) string, now time.Time) bool {
	base := getenv("XDG_STATE_HOME")
	if base == "" {
		if getenv("HOME") == "" {
			return false
		}
		base = filepath.Join(getenv("HOME"), ".local", "state")
	}
	if !filepath.IsAbs(base) {
		return false
	}
	day := now.UTC().Format("2006-01-02")
	f, err := openMarker(base, day)
	if err != nil {
		return false
	}
	_, err = io.WriteString(f, day+"\n")
	closeErr := f.Close()
	return err == nil && closeErr == nil
}
