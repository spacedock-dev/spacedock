// ABOUTME: Counter test for the one-run AC — appends a byte to FIXTURE_COUNTER_FILE
// ABOUTME: each execution so the AC test can assert the suite ran once, not twice.
package cleanoutputfixture

import (
	"os"
	"testing"
)

// TestCounter records one byte per execution to the file named by
// FIXTURE_COUNTER_FILE. The AC test reads the file size after running the
// changed command's shape: a "render clean, then re-run for json" regression
// would record two bytes, proving the double run; the one-run shape records one.
func TestCounter(t *testing.T) {
	path := os.Getenv("FIXTURE_COUNTER_FILE")
	if path == "" {
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("open counter file: %v", err)
	}
	defer f.Close()
	if _, err := f.Write([]byte{'x'}); err != nil {
		t.Fatalf("write counter: %v", err)
	}
}
