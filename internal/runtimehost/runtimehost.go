// ABOUTME: Env-only detection of the agent runtime host this process is running
// ABOUTME: under — the marker table and the ambiguity rule.
package runtimehost

// markerTable is the single marker table for runtime-host detection. Order is
// load-bearing: it fixes the order markers are reported in, which both the
// `--version` Runtime line and dispatch's ambiguity error render.
var markerTable = []struct {
	host   string
	marker string
}{
	{host: "codex", marker: "CODEX_THREAD_ID"},
	{host: "claude", marker: "CLAUDECODE"},
	{host: "pi", marker: "PI_CODING_AGENT"},
	{host: "pi", marker: "PI_CODING_AGENT_DIR"},
}

// Detect resolves the runtime host from env alone. It returns the host name, the
// markers that were set (in table order), and whether the set markers belong to
// more than one host.
//
// It NEVER returns an error, because its two callers need opposite dispositions
// of the same facts: dispatch must refuse to build against an ambiguous host,
// while `--version` must report the ambiguity and still exit 0 (refusing there
// would break the version gate and therefore every boot). So the shared unit is
// the detection and the policy stays with each caller.
//
// Two markers belonging to the SAME host — pi sets both PI_CODING_AGENT and
// PI_CODING_AGENT_DIR — are not ambiguity. A naive "more than one marker set"
// rule gets that wrong.
//
// When markers are ambiguous, host is empty: naming one would assert a host was
// resolved, which is precisely what this branch declines to do. With no markers
// set at all, host and markers are empty and ambiguous is false — being outside
// every runtime is a normal state, not a fault.
func Detect(getenv func(string) string) (host string, markers []string, ambiguous bool) {
	for _, entry := range markerTable {
		if getenv(entry.marker) == "" {
			continue
		}
		markers = append(markers, entry.marker)
		if host == "" {
			host = entry.host
			continue
		}
		if host != entry.host {
			ambiguous = true
		}
	}
	if ambiguous {
		return "", markers, true
	}
	return host, markers, false
}
