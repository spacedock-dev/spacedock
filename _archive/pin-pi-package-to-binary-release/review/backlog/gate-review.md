# Gate review — backlog → ideation: pin-pi-package-to-binary-release (qc)

## Ask

Approve advancing `pin-pi-package-to-binary-release` (id qcwrzza9xkr5kfwmdmvbqhkv) to ideation.

## What the seed owns

`spacedock install --host pi` installs the Spacedock Pi package from the bare source `git:github.com/spacedock-dev/spacedock` (`internal/cli/pi.go:27`, no `@ref`). Pi's package manager floats unpinned sources to the default branch, so the released stable v0.27.2 launcher installed v0.28.0-pre1 skills from `main`, and the first-officer binary gate (shared-core requires binary minor 0.27 in that release, 0.28 on main) correctly aborted before `status --boot`. The prior sibling fix pinned Claude and Codex through the marketplace channel and explicitly excluded Pi.

## Proposed mechanism

Pin the default Pi package source to the running launcher's release identity (release-stamped binary pins to its own release ref, never floats across tags). Keep `--plugin-dir` as the explicit dev override. Give the ordinary `spacedock pi` front door one repair attempt — if the package is missing, unpinned, or incompatible, install the binary's pinned source, recheck, and refuse the launch on remaining mismatch. Proof runs the ordinary installed front door with no development override (no `--plugin-dir`, no `SPACEDOCK_REPO_ROOT`, no prose marker).

## FO verification of seed claims (code-checked this session)

1. The unpinned source const is confirmed at `internal/cli/pi.go:27`; the Pi front door's install path uses it verbatim.
2. The abort guard is real and correctly did its job: `skills/first-officer/references/first-officer-shared-core.md:9` hard-requires an exact stamped major.minor; the skew was created by the installer, not the guard.
3. The sibling-fix claim is accurate: Claude/Codex pin via the marketplace channel mechanism (`internal/cli/init.go`); Pi was excluded from that fix's scope.

## QC amendments applied (2026-08-28, captain-ordered QC)

- The unstamped development-build source behavior is now an explicit ideation obligation (float to default branch, track `next`, or refuse — with the consequence owned), because a release-shaped binary must never float and an unstamped default silently re-creates the skew for dev builds.
- The pinning tradeoff is now stated as intended semantics: a pinned source does not auto-update; stable users receive package fixes by upgrading the launcher, matching the Claude/Codex marketplace pin.

## FO assessment

The mechanism targets the correct layer (installer source + front-door repair, not the binary floor), matches the Claude/Codex architecture, and the out-of-scope list protects the floor and the tag. AC-1's no-override installed-front-door proof is the run class that would have caught this at release time. The two QC amendments are folded into the approach; no AC changes were needed.

## Recommendation

APPROVED for advancement to ideation. Ideation must: pick and own the unstamped dev-build default source behavior, name the package-lock field the repair reads to detect missing/unpinned/incompatible, and keep the floor, the tag, and the Claude/Codex inventory untouched.
