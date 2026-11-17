//go:build unix

package cli

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/creack/pty"
)

// Only a real terminal counts, and the production client sends nothing but the
// User-Agent (no Accept-Encoding) to a local TLS server standing in for the route.
func TestSuiteCheckTerminalAndProductionRequest(t *testing.T) {
	master, slave, err := pty.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer master.Close()
	defer slave.Close()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()
	file, err := os.Create(filepath.Join(t.TempDir(), "out"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	for name, f := range map[string]*os.File{"pipe": w, "regular file": file} {
		if suiteCheckTerminal(f) {
			t.Errorf("%s counted as a terminal", name)
		}
	}
	headers := make(chan http.Header, 1)
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers <- r.Header.Clone()
		fmt.Fprint(w, `{"tag_name":"v0.0.0"}`)
	}))
	defer srv.Close()
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.Proxy = nil
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // DialContext only reaches the local test server.
	tr.DialContext = func(ctx context.Context, network, _ string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, network, strings.TrimPrefix(srv.URL, "https://"))
	}
	old := http.DefaultTransport
	http.DefaultTransport = tr
	t.Cleanup(func() { http.DefaultTransport = old; tr.CloseIdleConnections() })
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	t.Setenv("SPACEDOCK_NO_UPDATE_CHECK", "")
	suiteVersionHint(context.Background(), frontDoorArgs{}, os.Environ(), slave)
	select {
	case got := <-headers:
		if len(got) != 1 || got.Get("User-Agent") == "" {
			t.Fatalf("production wire headers=%v", got)
		}
	default:
		t.Fatal("a real terminal made no lookup")
	}
}
