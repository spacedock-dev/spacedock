// ABOUTME: Env-only detection of the agent runtime host this process is running
// ABOUTME: under — the marker table, the per-host session identity, and ShortID.
package runtimehost

// markerTable is the single marker table for runtime-host detection. Order is
// load-bearing: it fixes the order markers are reported in, which both the
// `--version` Runtime line and dispatch's ambiguity error render.
//
// `identity` is a SEPARATE column from `marker` and is not interchangeable with
// it. A detection marker only has to be set; an identity variable has to name
// THIS session. They coincide for codex alone (CODEX_THREAD_ID is both), and
// reusing the marker's value elsewhere would render `session 1` for every claude
// session, since CLAUDECODE is the literal flag `1`. An empty identity means the
// host exposes none (pi) or its variable is unset; both take the same path and
// the caller omits the segment.
var markerTable = []struct {
	host     string
	marker   string
	identity string
}{
	{host: "codex", marker: "CODEX_THREAD_ID", identity: "CODEX_THREAD_ID"},
	{host: "claude", marker: "CLAUDECODE", identity: "CLAUDE_CODE_SESSION_ID"},
	{host: "pi", marker: "PI_CODING_AGENT", identity: ""},
	{host: "pi", marker: "PI_CODING_AGENT_DIR", identity: ""},
}

// Detect resolves the runtime host from env alone. It returns the host name, the
// markers that were set (in table order), the session's own identifier where the
// host exposes one, and whether the set markers belong to more than one host.
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
// When markers are ambiguous, host and identity are both empty: naming either
// would assert a host was resolved, which is precisely what this branch declines
// to do. With no markers set at all, host and markers are empty and ambiguous is
// false — being outside every runtime is a normal state, not a fault.
func Detect(getenv func(string) string) (host string, markers []string, identity string, ambiguous bool) {
	for _, entry := range markerTable {
		if getenv(entry.marker) == "" {
			continue
		}
		markers = append(markers, entry.marker)
		if host == "" {
			host = entry.host
			if entry.identity != "" {
				identity = getenv(entry.identity)
			}
			continue
		}
		if host != entry.host {
			ambiguous = true
		}
	}
	if ambiguous {
		return "", markers, "", true
	}
	return host, markers, identity, false
}

// shortIDLen is the session-identifier prefix length. It is not arbitrary: it
// matches the convention already on disk, since Claude Code names its own team
// directory `~/.claude/teams/session-<first 8 chars>` of the session id. Printing
// the same eight characters makes the token directly greppable against session
// state on the filesystem, and eight hex characters distinguish the handful of
// sessions a human runs concurrently with room to spare.
const shortIDLen = 8

// ShortID truncates a session identifier to the greppable prefix. These ids are
// UUID-shaped and would otherwise dominate a four-line output. An id of
// shortIDLen characters or fewer is returned whole rather than padded, and an
// empty id stays empty so the caller omits the segment entirely.
func ShortID(id string) string {
	if len(id) <= shortIDLen {
		return id
	}
	return id[:shortIDLen]
}
