#!/usr/bin/env bash
# ABOUTME: Installs a pinned, sha256-verified prebuilt gotestsum onto PATH.
# ABOUTME: Used by the live-e2e jobs to render clean logs + a -json archive.
set -euo pipefail

# Pinned release. Bump VERSION and the matching checksum entries together; keep
# this a fixed version (never a floating reference) and never skip the checksum.
VERSION="1.13.0"

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"
case "$arch" in
  x86_64) arch="amd64" ;;
  aarch64 | arm64) arch="arm64" ;;
esac
platform="${os}_${arch}"

# sha256 of each prebuilt release tarball (gotestsum publishes a checksums file
# per release: gotestsum-${VERSION}-checksums.txt). A bare case keeps this
# portable to macOS's bash 3.2, which lacks associative arrays.
case "$platform" in
  linux_amd64) sha="11ccddeaf708ef228889f9fe2f68291a75b27013ddfc3b18156e094f5f40e8ee" ;;
  linux_arm64) sha="7644a4c5cd1bb978d56245aeab25a586ac5ac62adebed20a399548867c13499d" ;;
  darwin_amd64) sha="99529350f4c7b780b1efc543ca0d9721b09f0a4228f0efa9281261f58fefa05a" ;;
  darwin_arm64) sha="509cb27aef747f48faf9bce424f59dcf79572c905204b990ee935bbfcc7fa0e9" ;;
  *)
    echo "install-gotestsum: no pinned checksum for platform $platform" >&2
    exit 1
    ;;
esac

# Install destination: a stable per-user bin dir added to PATH. On GitHub Actions
# $HOME/.local/bin is already on PATH for the live jobs; locally the caller adds
# the printed dir to PATH (the AC test reads GOTESTSUM_BIN_DIR).
bindir="${GOTESTSUM_BIN_DIR:-$HOME/.local/bin}"
mkdir -p "$bindir"

tarball="gotestsum_${VERSION}_${platform}.tar.gz"
url="https://github.com/gotestyourself/gotestsum/releases/download/v${VERSION}/${tarball}"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
curl -fsSL -o "$tmp/$tarball" "$url"

# Verify the pinned checksum before trusting the binary.
echo "${sha}  ${tmp}/${tarball}" | shasum -a 256 -c -

tar -xzf "$tmp/$tarball" -C "$tmp" gotestsum
install -m 0755 "$tmp/gotestsum" "$bindir/gotestsum"

# Surface the install for the CI step log and add to PATH when running in Actions.
"$bindir/gotestsum" --version
if [ -n "${GITHUB_PATH:-}" ]; then
  echo "$bindir" >> "$GITHUB_PATH"
fi
