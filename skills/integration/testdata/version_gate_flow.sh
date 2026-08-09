#!/bin/sh
# version_gate_flow.sh — deterministic shell mirror of the FO version gate's
# Startup step 1 flow (skills/first-officer/references/first-officer-shared-core.md):
# sandbox session-env registry check (every class) -> --version line-1 parse ->
# binary-absent class: OS-aware hint -> sentinel loop bound -> install once ->
# SPACEDOCK_BIN session-scoped repoint -> --version re-check -> converge or
# hint-and-abort fallback. The Go driver (version_gate_fixture_test.go) puts
# captive `curl`, `install.sh`, and `spacedock` commands on a minimal PATH and
# asserts on exit codes and observable side effects, never prose wording.
set -u

REQUIRED_MINOR="${REQUIRED_MINOR:-0.27}"
REQUIRED_CAPABILITY="${REQUIRED_CAPABILITY:-spacedock gate withdraw <entity> --reason TEXT}"
LINUX_CMD='curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh'
DARWIN_CMD='brew tap spacedock-dev/homebrew-tap && brew install spacedock'
key="${CLAUDE_CODE_SESSION_ID:-${CODEX_THREAD_ID:-fixture}}"
SENTINEL="${TMPDIR:-/tmp}/spacedock-install-attempted-$key"
STDERR_CAPTURE="${TMPDIR:-/tmp}/install.stderr"

say() { printf '%s\n' "$*"; }

major_minor() { printf '%s' "$1" | sed -n 's/^\([0-9][0-9]*\)\.\([0-9][0-9]*\).*/\1.\2/p'; }

# --- Sandbox check (every class; session-env registry read, name+VALUE) ---
INSIDE=0
if [ "${APP_SANDBOX_CONTAINER_ID:-}" = "agent-safehouse" ]; then
	INSIDE=1
fi

# --- Gate: --version on the resolved launcher, line 1 only ---
present=0
vout=""
launcher="${SPACEDOCK_BIN:-spacedock}"
if vout="$("$launcher" --version 2>/dev/null)"; then
	case "$(printf '%s\n' "$vout" | sed -n '1p')" in
	"spacedock "*) present=1 ;;
	esac
fi

# --- ^Sandbox: corroboration (binary-present classes only; env wins) ---
if [ "$present" = 1 ] && [ "$INSIDE" = 1 ]; then
	sbline="$(printf '%s\n' "$vout" | sed -n 's/^Sandbox: \(.*\)$/\1/p')"
	case "$sbline" in
	inside*) say "corroboration: ^Sandbox: line agrees with the env verdict (inside)" ;;
	*) say "corroboration: DISAGREEMENT — env marker set but ^Sandbox: says '${sbline:-<absent>}'; env check wins" ;;
	esac
fi

# --- Binary-present classes ---
if [ "$present" = 1 ]; then
	ver="$(printf '%s\n' "$vout" | sed -n '1s/^spacedock //p')"
	if [ "$(major_minor "$ver")" = "$REQUIRED_MINOR" ]; then
		helpout=""
		if ! helpout="$("$launcher" gate --help 2>/dev/null)" || ! printf '%s\n' "$helpout" | awk -v required="$REQUIRED_CAPABILITY" 'index($0, required) { found=1 } END { exit !found }'; then
			os="${FIXTURE_OS:-$(uname -s)}"
			case "$os" in
			Darwin) REMEDY='brew upgrade spacedock' ;;
			Linux) REMEDY="$LINUX_CMD" ;;
			*) REMEDY='install the current Spacedock binary for this OS' ;;
			esac
			say "selected launcher: $launcher"
			say "observed version: spacedock $ver"
			say "missing capability: $REQUIRED_CAPABILITY"
			say "upgrade the installed launcher: $REMEDY"
			say "then relaunch"
			exit 3
		fi
		say "gate passed: spacedock $ver"
		"$launcher" status --boot --identify --json
		exit $?
	fi
	# Wrong-version class: unchanged by this task — abort + doctor pointer.
	say "version mismatch: binary $ver, skills require $REQUIRED_MINOR — run \${SPACEDOCK_BIN:-spacedock} doctor for the per-class remedy"
	[ "$INSIDE" = 1 ] && say "sandbox: env verdict is inside — named in the human-facing message"
	exit 3
fi

# --- Binary-absent class: OS-aware hint (FIXTURE_OS pins the uname -s source) ---
os="${FIXTURE_OS:-$(uname -s)}"
case "$os" in
Linux) INSTALL_CMD="$LINUX_CMD" ;;
Darwin) INSTALL_CMD="$DARWIN_CMD" ;;
*)
	say "hint-and-abort: unsupported OS $os — install.sh does not support it; build from source: go build -o spacedock ./cmd/spacedock"
	exit 3
	;;
esac

# Inside a sandbox the install is never OFFERED or run (silent no-op for the host).
if [ "$INSIDE" = 1 ]; then
	say "You're running inside a sandbox. Run this command yourself, outside the sandbox: $INSTALL_CMD"
	exit 3
fi

hint_and_abort() {
	say "hint-and-abort: spacedock is not on PATH — run: $INSTALL_CMD"
	if [ -f "$SENTINEL" ]; then
		say "an install attempt was already made: sentinel $SENTINEL exists; rm $SENTINEL re-enables the install offer"
	fi
	exit 3
}

# --- One-attempt loop bound: durable sentinel file ---
if [ -f "$SENTINEL" ]; then
	say "sentinel-blocked: install already attempted"
	hint_and_abort
fi

# --- Operator-approved install offer: create BEFORE run, run once ---
touch "$SENTINEL"
eval "$INSTALL_CMD" 2>"$STDERR_CAPTURE" || say "install command exited $?"
inst_path="$(sed -n 's/^install\.sh: installed spacedock [^ ]* to \(\/.*\)$/\1/p' "$STDERR_CAPTURE")"
if [ -z "$inst_path" ] && [ -x "$HOME/.local/bin/spacedock" ]; then
	inst_path="$HOME/.local/bin/spacedock"
fi
if [ -z "$inst_path" ]; then
	say "install produced no installed-path line and the \$HOME/.local/bin probe missed"
	hint_and_abort
fi

# --- Session-scoped repoint + re-check ---
export SPACEDOCK_BIN="$inst_path"
say "SPACEDOCK_BIN repointed to $SPACEDOCK_BIN (this session only; nothing persisted)"
vout2="$("$SPACEDOCK_BIN" --version 2>/dev/null)" || true
line1="$(printf '%s\n' "$vout2" | sed -n '1p')"
case "$line1" in
"spacedock "*)
	if [ "$(major_minor "${line1#spacedock }")" = "$REQUIRED_MINOR" ]; then
		say "gate passed after install: $line1"
		exit 0
	fi
	;;
esac
say "--version re-check failed after install ($line1)"
hint_and_abort
