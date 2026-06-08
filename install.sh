#!/bin/sh
# ABOUTME: Universal curl|sh installer — fetch the checksum-verified spacedock
# ABOUTME: release tarball for this host's OS/arch and install the bare binary.
#
# Usage (Linux or macOS):
#   curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/next/install.sh | sh
#
# Behavior:
#   * Detects OS (darwin|linux) and arch (amd64|arm64) from uname.
#   * Resolves the asset to fetch from one of two sources, same extract/verify/
#     install path for both:
#       - SPACEDOCK_INSTALL_FROM unset (production): the latest GitHub Release.
#       - SPACEDOCK_INSTALL_FROM=<dir|url-base> (tests / pinned mirror): a local
#         goreleaser `dist/` directory or a URL prefix holding the same
#         `spacedock_<ver>_<os>_<arch>.tar.gz` + `checksums.txt` layout.
#   * Verifies the tarball sha256 against the matching `checksums.txt` line and
#     ABORTS (installs nothing) on any mismatch — the gate is fail-closed.
#   * Extracts the bare `spacedock` binary and installs it to SPACEDOCK_INSTALL_DIR
#     (default ~/.local/bin).
#
# SPACEDOCK_PRINT_TARGET=1 runs the detection + URL-construction path and prints
# the resolved os/arch/asset/tarball/checksums, then exits before any download —
# the inspection seam the AC-3 live-URL test asserts against.
#
# The macOS Homebrew cask remains the preferred mac path; this script is the
# universal fallback and the only Linux install path.
set -eu

REPO="spacedock-dev/spacedock"
INSTALL_DIR="${SPACEDOCK_INSTALL_DIR:-$HOME/.local/bin}"

err() { printf 'install.sh: %s\n' "$*" >&2; }
die() { err "$*"; exit 1; }

need() { command -v "$1" >/dev/null 2>&1 || die "required command not found: $1"; }

# detect_os maps `uname -s` to the goreleaser goos token. Only darwin and linux
# ship release tarballs; anything else is unsupported here (use Homebrew on mac,
# or build from source).
detect_os() {
	case "$(uname -s)" in
		Darwin) echo darwin ;;
		Linux) echo linux ;;
		*) die "unsupported OS $(uname -s); see docs/install-journey.md for source build" ;;
	esac
}

# detect_arch maps `uname -m` to the goreleaser goarch token. The release builds
# amd64 + arm64; uname reports several spellings for each.
detect_arch() {
	case "$(uname -m)" in
		x86_64 | amd64) echo amd64 ;;
		arm64 | aarch64) echo arm64 ;;
		*) die "unsupported arch $(uname -m); release ships amd64 + arm64 only" ;;
	esac
}

# sha256_of prints the lowercase hex sha256 of a file, using whichever tool the
# host carries: sha256sum (Linux), or `shasum -a 256` (macOS).
sha256_of() {
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum "$1" | awk '{print $1}'
	elif command -v shasum >/dev/null 2>&1; then
		shasum -a 256 "$1" | awk '{print $1}'
	else
		die "no sha256 tool found (need sha256sum or shasum)"
	fi
}

# fetch copies a source ref (local file path or http(s) URL) to a destination
# file. A missing local file or a non-2xx HTTP status is a hard failure so the
# caller never proceeds on a partial download.
fetch() {
	src="$1"
	dst="$2"
	case "$src" in
		http://* | https://*)
			curl -fsSL -o "$dst" "$src" || die "download failed: $src"
			;;
		*)
			[ -f "$src" ] || die "file not found: $src"
			cp "$src" "$dst"
			;;
	esac
}

# resolve_latest_tag asks the GitHub API for the repo's latest release tag (e.g.
# v0.20.0). Unauthenticated; the public releases endpoint needs no token.
resolve_latest_tag() {
	api="https://api.github.com/repos/$REPO/releases/latest"
	curl -fsSL "$api" \
		| grep '"tag_name"' \
		| head -n 1 \
		| sed -E 's/.*"tag_name"[[:space:]]*:[[:space:]]*"([^"]+)".*/\1/'
}

# asset_name builds the goreleaser archive name for a version/os/arch. The
# version carries no leading `v` (goreleaser stamps the bare semver into the
# `{{ .Version }}` template); the tag does, so callers strip it.
asset_name() {
	printf 'spacedock_%s_%s_%s.tar.gz' "$1" "$2" "$3"
}

main() {
	need uname
	need curl
	need tar
	need mktemp

	os="$(detect_os)"
	arch="$(detect_arch)"

	# Resolve the asset name + the source refs (tarball + checksums) for either
	# source. The SAME verify/extract/install path runs below for both.
	if [ -n "${SPACEDOCK_INSTALL_FROM:-}" ]; then
		from="${SPACEDOCK_INSTALL_FROM%/}"
		# A local dist directory embeds the snapshot version in the filename, not
		# known a priori, so glob for the os/arch tarball. A URL base must carry a
		# resolvable version, so SPACEDOCK_INSTALL_VERSION pins it.
		case "$from" in
			http://* | https://*)
				ver="${SPACEDOCK_INSTALL_VERSION:?SPACEDOCK_INSTALL_VERSION required when SPACEDOCK_INSTALL_FROM is a URL}"
				asset="$(asset_name "$ver" "$os" "$arch")"
				tarball_src="$from/$asset"
				checksums_src="$from/checksums.txt"
				;;
			*)
				[ -d "$from" ] || die "SPACEDOCK_INSTALL_FROM is not a directory or URL: $from"
				asset="$(cd "$from" && ls spacedock_*_"${os}"_"${arch}".tar.gz 2>/dev/null | head -n 1)"
				[ -n "$asset" ] || die "no spacedock_*_${os}_${arch}.tar.gz in $from"
				tarball_src="$from/$asset"
				checksums_src="$from/checksums.txt"
				;;
		esac
	else
		tag="$(resolve_latest_tag)"
		[ -n "$tag" ] || die "could not resolve the latest release tag for $REPO"
		ver="${tag#v}"
		asset="$(asset_name "$ver" "$os" "$arch")"
		base="https://github.com/$REPO/releases/download/$tag"
		tarball_src="$base/$asset"
		checksums_src="$base/checksums.txt"
	fi

	# Inspection mode: print the resolved target and stop before any download or
	# install. This runs the EXACT production detection + URL-construction path
	# above (no divergent branch) so a test can assert the asset name + URL the
	# real installer would fetch against the live release.
	if [ -n "${SPACEDOCK_PRINT_TARGET:-}" ]; then
		printf 'os=%s\narch=%s\nasset=%s\ntarball=%s\nchecksums=%s\n' \
			"$os" "$arch" "$asset" "$tarball_src" "$checksums_src"
		return 0
	fi

	tmp="$(mktemp -d)"
	trap 'rm -rf "$tmp"' EXIT

	fetch "$tarball_src" "$tmp/$asset"
	fetch "$checksums_src" "$tmp/checksums.txt"

	# Checksum gate (fail-closed). Pull THIS asset's expected hash from
	# checksums.txt by exact filename, compute the downloaded tarball's hash, and
	# abort installing anything on any mismatch or a missing checksum line.
	expected="$(awk -v f="$asset" '$2 == f {print $1}' "$tmp/checksums.txt" | head -n 1)"
	[ -n "$expected" ] || die "no checksum line for $asset in checksums.txt — refusing to install"
	actual="$(sha256_of "$tmp/$asset")"
	if [ "$expected" != "$actual" ]; then
		die "checksum mismatch for $asset (expected $expected, got $actual) — refusing to install"
	fi

	# Extract the bare `spacedock` binary (archive root, no wrapping dir) and
	# install it. Only after the checksum passes do we touch the install dir.
	tar -xzf "$tmp/$asset" -C "$tmp" spacedock || die "tarball did not contain a spacedock binary"
	mkdir -p "$INSTALL_DIR"
	install -m 0755 "$tmp/spacedock" "$INSTALL_DIR/spacedock" 2>/dev/null \
		|| { cp "$tmp/spacedock" "$INSTALL_DIR/spacedock" && chmod 0755 "$INSTALL_DIR/spacedock"; }

	printf 'install.sh: installed spacedock %s to %s/spacedock\n' "$asset" "$INSTALL_DIR" >&2
	case ":$PATH:" in
		*":$INSTALL_DIR:"*) ;;
		*) err "note: $INSTALL_DIR is not on PATH; add it to run 'spacedock' directly" ;;
	esac
}

main "$@"
