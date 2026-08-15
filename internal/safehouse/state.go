// ABOUTME: The sandbox-detection registry (Inside) plus the two render surfaces:
// ABOUTME: SessionState ("is this process sandboxed?") and LaunchState ("will this launch wrap?").
package safehouse

// insideRegistry maps an environment signal to the sandbox it proves this
// process is already running inside. It is a table rather than a boolean because
// the display name must come from the SAME row as the signal — a second sandbox
// implementation adds one row and reports its own name, where a hardcoded
// "agent-safehouse" in the render string would report every sandbox as safehouse.
//
// Matching is on the VALUE, not on mere presence: APP_SANDBOX_CONTAINER_ID is a
// generic macOS app-sandbox variable that other containers also set, so presence
// alone would claim any of them for safehouse.
var insideRegistry = []struct {
	env         string
	wantValue   string
	displayName string
}{
	{env: "APP_SANDBOX_CONTAINER_ID", wantValue: "agent-safehouse", displayName: "agent-safehouse"},
}

// Inside reports whether THIS process is executing inside a sandbox, and names
// it. This is a different question from Present/Available, which ask whether a
// LAUNCH from here would be wrapped — and in the common case the two are
// anti-correlated: inside the sandbox the safehouse binary is off PATH precisely
// because the wrap already happened, so the state in which "am I sandboxed" is
// most emphatically yes is the state the old availability-only render called
// `unavailable`.
func Inside(getenv func(string) string) (name string, ok bool) {
	for _, entry := range insideRegistry {
		if getenv(entry.env) == entry.wantValue {
			return entry.displayName, true
		}
	}
	return "", false
}

// SessionState renders the answer to "is this process sandboxed?" — the question
// `--version` and `status --boot` ask, both being reports about the running
// session rather than about a launch.
//
// It deliberately says nothing about a .safehouse profile: a profile determines
// whether a launch would wrap, and a session that is already running is not
// about to launch itself. Not being sandboxed, the line reports availability and
// nothing more, because "not wrapped" would be constant across every such case
// and so would carry no information while reading as though it did.
//
// The value never restates its own label. Under a `Sandbox: ` prefix, the
// sandbox's NAME alone is the whole answer — it both says that this process is
// sandboxed and says which sandbox — and `none` is its counterpart, so a reader
// classifies on a name-or-`none` distinction rather than on a relationship word
// like "inside" or "not sandboxed" that the label already implies.
//
//   - <name> — this process is running inside the named sandbox. This dominates:
//     whether safehouse could wrap a future launch says nothing about a session
//     that is already sandboxed.
//   - none (safehouse available) — the binary resolves on PATH.
//   - none (safehouse not installed) — it does not.
func SessionState(insideName string, inside, available bool) string {
	if inside {
		return insideName
	}
	if available {
		return "none (safehouse available)"
	}
	return "none (safehouse not installed)"
}

// LaunchState renders the answer to "will this launch be wrapped?" — the
// question the pre-launch banner asks, and the only surface where the .safehouse
// profile is a load-bearing input. `selected` is whether this launch would be
// wrapped (a .safehouse profile is present, or a --safehouse* flag forced it);
// `available` is whether the safehouse binary resolves on PATH.
//
// The `inside` arm exists because a launch from within the sandbox inverted the
// same way the session surfaces did: safehouse is off PATH there, so the banner
// reported `unavailable` while already sandboxed. Being inside dominates — there
// is nothing left to wrap.
func LaunchState(insideName string, inside, selected, available bool) string {
	if inside {
		return "inside (" + insideName + ") — launching without re-wrapping"
	}
	if !selected {
		return "not wrapping this launch (no .safehouse profile)"
	}
	if !available {
		return "not wrapped (safehouse not installed; .safehouse profile present)"
	}
	return "wrapping this launch (safehouse, .safehouse profile)"
}
