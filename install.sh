#!/bin/sh
# ABOUTME: Universal curl|sh installer — fetch the checksum-verified spacedock
# ABOUTME: release tarball for this host's OS/arch and install the bare binary.
#
# Usage (Linux or macOS):
#   curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh
#
# SPACEDOCK_CHANNEL selects the release channel: `stable` (default) resolves the
#   latest stable release and the unsuffixed asset; `edge` resolves the newest
#   release INCLUDING prereleases and the `_edge` asset. Any other value aborts —
#   a typo must not silently install the other channel. The resolved tag is
#   printed to stderr before any download.
#
# Behavior:
#   * Detects OS (darwin|linux) and arch (amd64|arm64) from uname.
#   * Resolves the asset to fetch from one of two sources, same extract/verify/
#     install path for both:
#       - SPACEDOCK_INSTALL_FROM unset (production): the newest GitHub Release on
#         the selected channel.
#       - SPACEDOCK_INSTALL_FROM=<dir|url-base> (tests / pinned mirror): a local
#         goreleaser `dist/` directory or a URL prefix holding the same
#         `spacedock_<ver>_<os>_<arch>[_edge].tar.gz` + `checksums.txt` layout.
#   * Verifies the tarball sha256 against the matching `checksums.txt` line and
#     ABORTS (installs nothing) on any mismatch — the gate is fail-closed.
#   * Extracts the bare `spacedock` binary and installs it to SPACEDOCK_INSTALL_DIR
#     (default ~/.local/bin).
#
# SPACEDOCK_PLUGINS=none skips the plugin step; `spacedock update` sets it and
# refreshes the plugins itself, only on the hosts that already have Spacedock's.
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
CHANNEL="${SPACEDOCK_CHANNEL:-stable}"

err() { printf 'install.sh: %s\n' "$*" >&2; }
die() { err "$*"; exit 1; }

need() { command -v "$1" >/dev/null 2>&1 || die "required command not found: $1"; }

# validate_channel rejects any channel value other than the two accepted ones.
# Falling through to stable would silently install the OTHER channel on a typo
# (`SPACEDOCK_CHANNEL=egde`) — the silent skew this script exists to close.
validate_channel() {
	case "$CHANNEL" in
		stable | edge) ;;
		*) die "unknown SPACEDOCK_CHANNEL '$CHANNEL'; accepted values: stable, edge" ;;
	esac
}

# detect_os maps `uname -s` to the goreleaser goos token. Only darwin and linux
# ship release tarballs; anything else is unsupported here (use Homebrew on mac,
# or build from source).
detect_os() {
	case "$(uname -s)" in
		Darwin) echo darwin ;;
		Linux) echo linux ;;
		*) die "unsupported OS $(uname -s); see docs/site/get-started/install.md for source build" ;;
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

# resolve_latest_tag asks the GitHub API for the newest release tag on the
# selected channel (e.g. v0.20.0). `/releases/latest` EXCLUDES prereleases, which
# is exactly the stable channel; the edge line ships as `-pre` tags, so edge
# reads the newest entry of the created_at-descending releases list instead.
# Unauthenticated and one request either way, and `tag_name` appears only at
# release level in both responses, so the same first-match parse serves both.
resolve_latest_tag() {
	if [ "$CHANNEL" = edge ]; then
		api="https://api.github.com/repos/$REPO/releases?per_page=1"
	else
		api="https://api.github.com/repos/$REPO/releases/latest"
	fi
	curl -fsSL "$api" \
		| grep '"tag_name"' \
		| head -n 1 \
		| sed -E 's/.*"tag_name"[[:space:]]*:[[:space:]]*"([^"]+)".*/\1/'
}

# asset_name builds the goreleaser archive name for a version/os/arch on the
# selected channel. The version carries no leading `v` (goreleaser stamps the
# bare semver into the `{{ .Version }}` template); the tag does, so callers strip
# it. Every release publishes a PAIR of per-arch archives: the edge-stamped one
# always carries `_edge`, and on a prerelease the stable one carries `_stable`
# (.goreleaser.yaml), so the unsuffixed name exists only on stable tags — the
# only tags the stable channel ever resolves.
asset_name() {
	if [ "$CHANNEL" = edge ]; then
		printf 'spacedock_%s_%s_%s_edge.tar.gz' "$1" "$2" "$3"
	else
		printf 'spacedock_%s_%s_%s.tar.gz' "$1" "$2" "$3"
	fi
}

EXTRAS='subspace|spacedock-dev/subspace|subspace|subspace-tui subspace-relay-client|subspace-tui:spacedock-review|bin'

tree_version() {
	tlink="$INSTALL_DIR/$(basename "$1")"
	[ -e "$tlink" ] || return 0
	ttarget="$(readlink "$tlink" 2>/dev/null)" || return 0
	troot="${XDG_DATA_HOME:-$HOME/.local/share}/$2/"
	case "$ttarget" in
		"$troot"*) trest="${ttarget#"$troot"}"; printf '%s\n' "${trest%%/*}" ;;
	esac
}

installed_version() {
	iver=""
	for iprog in $2; do
		if [ "$3" = tree ]; then
			v="$(tree_version "$iprog" "$1")"
		else
			[ -x "$INSTALL_DIR/$iprog" ] || return 0
			v="$("$INSTALL_DIR/$iprog" --version 2>/dev/null | head -n 1)"
		fi
		[ -n "$v" ] && { [ -z "$iver" ] || [ "$iver" = "$v" ]; } || return 0
		iver="$v"
	done
	for l in $4; do
		[ -e "$INSTALL_DIR/${l#*:}" ] || return 0
	done
	printf '%s\n' "$iver"
}

install_extra() {
	name="$1" repo="$2" prefix="$3" programs="$4" links="$5" kind="$6" os="$7" arch="$8"
	dir="$(mktemp -d)"
	local_from="${SPACEDOCK_EXTRAS_FROM:-}"
	asset=""
	if [ -n "$local_from" ]; then
		asset="$(cd "$local_from" 2>/dev/null && ls ${prefix}_*_${os}_${arch}.tar.gz 2>/dev/null | head -n 1)" || asset=""
	fi
	if [ -n "$asset" ]; then
		ver="${asset#${prefix}_}"; ver="${ver%_${os}_${arch}.tar.gz}"; tag="v$ver (from $local_from)"
		[ "$(installed_version "$name" "$programs" "$kind" "$links")" = "$ver" ] && { err "$name $ver unchanged"; rm -rf "$dir"; return 0; }
		cp "$local_from/$asset" "$dir/$asset" && cp "$local_from/checksums.txt" "$dir/checksums.txt" || return 1
	else
		if [ "$CHANNEL" = edge ]; then
			api="https://api.github.com/repos/$repo/releases?per_page=1"
		else
			api="https://api.github.com/repos/$repo/releases/latest"
		fi
		tag="$(curl -fsSL "$api" \
			| grep '"tag_name"' | head -n 1 \
			| sed -E 's/.*"tag_name"[[:space:]]*:[[:space:]]*"([^"]+)".*/\1/')" || tag=""
		[ -n "$tag" ] || { err "$name: could not resolve the latest release of $repo"; return 1; }
		ver="${tag#v}"
		[ "$(installed_version "$name" "$programs" "$kind" "$links")" = "$ver" ] && { err "$name $ver unchanged"; rm -rf "$dir"; return 0; }
		asset="${prefix}_${ver}_${os}_${arch}.tar.gz"
		base="https://github.com/$repo/releases/download/$tag"
		curl -fsSL -o "$dir/$asset" "$base/$asset" && curl -fsSL -o "$dir/checksums.txt" "$base/checksums.txt" \
			|| { err "$name: download failed — skipped"; return 1; }
	fi
	want_sum="$(awk -v f="$asset" '$2 == f {print $1}' "$dir/checksums.txt" | head -n 1)"
	[ -n "$want_sum" ] && [ "$want_sum" = "$(sha256_of "$dir/$asset")" ] \
		|| { err "$name: checksum mismatch for $asset — skipped"; return 1; }
	if [ "$kind" = tree ]; then
		prev="$(tree_version "${programs%% *}" "$name")"
		root="${XDG_DATA_HOME:-$HOME/.local/share}/$name"
		dest="$root/$ver"
		rm -rf "$dest" && mkdir -p "$dest" && tar -xzf "$dir/$asset" -C "$dest" || return 1
		for p in $programs; do
			[ -f "$dest/$p" ] && ln -sf "$dest/$p" "$INSTALL_DIR/$(basename "$p")"
		done
		for old in "$root"/*; do
			[ -n "$prev" ] && [ -d "$old" ] && [ "$(basename "$old")" != "$ver" ] && [ "$(basename "$old")" != "$prev" ] && rm -rf "$old"
		done
	else
		tar -xzf "$dir/$asset" -C "$dir" || return 1
		for p in $programs; do
			[ -f "$dir/$p" ] || continue
			install -m 0755 "$dir/$p" "$INSTALL_DIR/$p" 2>/dev/null \
				|| { cp "$dir/$p" "$INSTALL_DIR/$p" && chmod 0755 "$INSTALL_DIR/$p"; } || return 1
		done
	fi
	for l in $links; do
		[ -f "$INSTALL_DIR/${l%%:*}" ] && ln -sf "${l%%:*}" "$INSTALL_DIR/${l#*:}"
	done
	rm -rf "$dir"
	err "installed $name $tag"
}

# install_plugins uses the spacedock this script just installed, not an older copy
# that comes first on PATH.
install_plugins() {
	[ "${SPACEDOCK_PLUGINS:-}" = none ] && return 0
	if command -v claude >/dev/null 2>&1; then
		err "adding the Claude Code plugins (spacedock install --host claude)"
		"$INSTALL_DIR/spacedock" install --host claude \
			|| err "the plugin install failed; run: spacedock install --host claude"
	else
		err "claude is not on PATH; run \`spacedock install --host claude\` after you install Claude Code"
	fi
}

# SPACEDOCK_INSTALL_FROM runs skip the extensions unless SPACEDOCK_EXTRAS asks, so
# tests and CI stay offline.
install_extras() {
	want="${SPACEDOCK_EXTRAS:-}"
	[ "$want" = none ] && return 0
	[ -n "${SPACEDOCK_INSTALL_FROM:-}" ] && [ -z "$want" ] && return 0
	printf '%s\n' "$EXTRAS" | while IFS='|' read -r name repo prefix programs links kind; do
		[ -n "$name" ] || continue
		install_extra "$name" "$repo" "$prefix" "$programs" "$links" "${kind:-bin}" "$1" "$2" \
			|| err "$name failed to install; spacedock works without it"
	done
	install_plugins
}

main() {
	validate_channel
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
				# A goreleaser dist holds BOTH channels' archives, so the glob
				# carries the channel too. The two are disjoint: the unsuffixed
				# pattern cannot match a name ending `_edge.tar.gz`.
				if [ "$CHANNEL" = edge ]; then
					glob="spacedock_*_${os}_${arch}_edge.tar.gz"
				else
					glob="spacedock_*_${os}_${arch}.tar.gz"
				fi
				asset="$(cd "$from" && ls $glob 2>/dev/null | head -n 1)"
				[ -n "$asset" ] || die "no $glob in $from"
				tarball_src="$from/$asset"
				checksums_src="$from/checksums.txt"
				;;
		esac
	else
		tag="$(resolve_latest_tag)"
		[ -n "$tag" ] || die "could not resolve the latest release tag for $REPO"
		# Say which version this run resolved, before anything is downloaded, so a
		# channel skew is visible at install time instead of surfacing later as a
		# first-officer boot abort. stderr and `=`-free, so the print-target
		# stdout below stays machine-parseable.
		err "resolved $CHANNEL channel to $tag"
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

	first="$(command -v spacedock 2>/dev/null)" || first=""
	if [ -n "$first" ] && ! [ "$first" -ef "$INSTALL_DIR/spacedock" ]; then
		err "warning: $first ($("$first" --version 2>/dev/null | head -n 1)) comes first on PATH and will run instead"
	fi

	install_extras "$os" "$arch"
}

main "$@"
