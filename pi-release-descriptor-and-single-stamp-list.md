---
title: Give the Pi extension its own release descriptor, and make the stamp list exist once
status: implementation
score: 0.8
source: "Captain constraint, 2026-10-04: the Pi extension's canonical metadata must not be the top-level package.json, because every other host has its own. Astra design review, run b3928dbb."
id: 8vahpzjd358ygqce5etw6fr2
started: 2026-10-04T05:28:21Z
worktree: .worktrees/spacedock-ensign-pi-release-descriptor-and-single-stamp-list
---

## Problem

The Pi extension has no release descriptor of its own. Claude and Codex each carry
one, `.claude-plugin/plugin.json` and `.codex-plugin/plugin.json`, and each carries
its `version`. Pi has neither, so the only candidate for its version is the root
`package.json` — the same file every other host would want, and the file Pi reads
for resource routing. Claiming it for one host is wrong, and leaving Pi unversioned
means the release cannot stamp it at all.

Two release defects sit in the same family. `StampVersion` returns a named manifest
unchanged when it has no top-level `version`, and the release CLI then writes those
unchanged bytes and prints success. And `replaceFirstVersion` rewrites the first
textual `version` occurrence rather than the top-level key, so a nested version can
capture the write.

The stamp list is written out in three places in `docs/releasing.md`, and repeated
again in the release workflow, so adding a fourth manifest would add a fifth
authority.

## Design

**1. New file `.pi/plugin.json`**, beside the extension it describes. Spacedock-owned
release metadata, not something Pi natively discovers:

    {
      "name": "spacedock",
      "version": "0.28.0-pre3",
      "description": "Turn directories of markdown files into structured workflows operated by AI agents",
      "author": { "name": "CL Kao" },
      "homepage": "https://github.com/spacedock-dev/spacedock",
      "repository": "https://github.com/spacedock-dev/spacedock",
      "license": "Apache-2.0",
      "keywords": ["workflow", "pipeline", "agents", "markdown", "automation"],
      "interface": {
        "displayName": "Spacedock",
        "shortDescription": "Plain text workflows for AI agents",
        "longDescription": "Turn directories of markdown files into structured workflows operated by AI agents.",
        "developerName": "CL Kao",
        "category": "workflow",
        "capabilities": ["Interactive", "Write"],
        "websiteURL": "https://github.com/spacedock-dev/spacedock"
      }
    }

It carries identity, description and version. It carries NO resource routing: Pi
continues to read `pi.extensions` and `pi.skills` from the root `package.json`, which
stays unchanged and unversioned.

**2. One authoritative target list** in `internal/release`:

    func StampTargets() []string {
        return []string{
            ".claude-plugin/plugin.json",
            ".codex-plugin/plugin.json",
            ".pi/plugin.json",
            "skills/first-officer/references/first-officer-shared-core.md",
        }
    }

- `stamp-version <version>` and `manifest-tag-gate <tag>` use it when given no
  explicit targets.
- A new `stamp-paths` subcommand prints it, so the release workflow's gate, stamp,
  diff and commit steps stop repeating it, and the `FO_PROSE` variable goes away.
- `docs/releasing.md` keeps one invocation per ritual step and drops its target lists.

**3. Two correctness fixes in the release CLI.** A named JSON manifest with no usable
top-level `version` must fail with a clear error and a non-zero exit instead of
writing unchanged bytes and reporting success; the low-level marketplace no-op
contract keeps its behaviour, because it is explicitly tested. And the top-level
`version` must be the field replaced, even when a nested `version` appears earlier.

**4. Root `package.json` stays as it is.** Its narrower responsibility is documented
where the release ritual is described, not encoded in the file.

## Out of scope

- Moving the Pi package into a subdirectory, and any `git:.../subdirectory` syntax.
  Pi's git source model has no subdirectory selector.
- Splitting the extension into its own repository. The extension is coupled to the
  shared skills tree, so that is a distribution migration, not a metadata move.
- Making Pi display or enforce the descriptor's version.
- Prerelease coherence, where CI skips the manifest and tag gate for prereleases.

## Expected surface and tolerance

New `.pi/plugin.json`; `internal/release/release.go` plus its test;
`cmd/spacedock-release/main.go`; `docs/releasing.md`; `.github/workflows/release.yml`.
Estimate net +120, across 6 files. Tolerance net +50 to +200, at most 8 files.

## Acceptance criteria

**AC-1 (VALUE) - One invocation stamps every host, including Pi.**
Verified by: `stamp-version X.Y.Z` with no target arguments rewrites the version in
`.claude-plugin/plugin.json`, `.codex-plugin/plugin.json`, `.pi/plugin.json` and the
shared-core prose. Falsifier: omit the Pi descriptor from the list and the stamp
leaves it at its old version.

**AC-2 - The list has one authority.**
Verified by: `stamp-paths` prints exactly the list the stamping and the tag gate use,
and the test fails when the two diverge. Falsifier: give the gate a different list.

**AC-3 - A named JSON manifest with no version fails loudly.**
Verified by: naming a JSON manifest without a usable top-level version exits
non-zero with a clear error and prints no success line. Falsifier: remove the check
and the command exits zero on unchanged bytes.

**AC-4 - The top-level field is the one replaced.**
Verified by: a manifest carrying a nested `version` before the top-level one; the
top-level value changes and the nested one does not. Falsifier: keep
first-occurrence replacement and the nested value changes instead.

**AC-5 (no-regression) -** The marketplace no-op contract, `dev-preversion`,
`manifest-tag-gate`, and the existing release tests keep their behaviour; the Pi
extension still loads through the real installed loader with the root manifest
untouched; the offline suite passes.

## Test plan

Behavioural tests only, with the version supplied independently. No test may assert a
file's own text against a copy of it, and no expected value may come from the thing
it tests. Pi discovery is checked through the real installed loader, not a copied
manifest string.
