#!/bin/sh
# ABOUTME: Full-set curl|sh installer — runs each product's own install.sh in turn.
# ABOUTME: Holds only the list; how each product installs stays in its own script.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install-all.sh | sh
#
# A failed spacedock install fails the run; a failed extra is reported and skipped.
set -u

INSTALL_DIR="${SPACEDOCK_INSTALL_DIR:-$HOME/.local/bin}"
SPACEDOCK_INSTALLER="${SPACEDOCK_INSTALLER:-https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh}"
SUBSPACE_INSTALLER="${SUBSPACE_INSTALLER:-https://raw.githubusercontent.com/spacedock-dev/subspace/main/install.sh}"

export SPACEDOCK_INSTALL_DIR="$INSTALL_DIR" SUBSPACE_INSTALL_DIR="$INSTALL_DIR"

installed=""
failed=""

install_one() {
	printf '==> %s\n' "$1" >&2
	script=$(mktemp)
	if curl -fsSL "$2" -o "$script" && sh "$script"; then
		installed="$installed $1"
	else
		failed="$failed $1"
	fi
	rm -f "$script"
}

install_one spacedock "$SPACEDOCK_INSTALLER"
case " $failed " in
	*" spacedock "*) printf 'install-all.sh: spacedock failed to install\n' >&2; exit 1 ;;
esac
install_one subspace "$SUBSPACE_INSTALLER"

printf 'install-all.sh: installed:%s\n' "$installed" >&2
[ -z "$failed" ] || printf 'install-all.sh: failed:%s (spacedock still works without it)\n' "$failed" >&2
printf 'install-all.sh: next, run `spacedock claude`; it installs the plugins on first launch\n' >&2
