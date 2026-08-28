# Gate review — ideation: pin-pi-package-to-binary-release (qc)

## Selected approach (cycle 2, review-confirmed)

**Pin ref derivation** (one new helper file, ordered): (1) linker-stamped `internal/cli.Version` (goreleaser artifacts, stable and edge); (2) `debug.ReadBuildInfo().Main.Version` when it is a semver tag — covers `go install …@vX.Y.Z` proxy builds, which are release-shaped and previously floated (the cycle-1 blocker); (3) otherwise the `dev` sentinel keeps the unpinned floating source, consequence stated (first-officer binary gate is the loud backstop). A release-shaped binary never floats. Tradeoff stated as intended semantics: pinned sources do not auto-update; stable users get fixes by upgrading the launcher (matches the Claude/Codex marketplace pin). `--plugin-dir` remains the explicit dev override.

## Risk evidence

Spike verified against the real shipped Pi package-manager dist (`core/package-manager.js` + `utils/git.js`): `pinned: Boolean(ref)`; unpinned git sources update from origin HEAD; `pi install git:<repo>@<ref>` clones+checkouts the ref; `addSourceToSettings` rewrites the existing entry (match key ignores ref). The wrong-ref repair trigger is detected from the settings entry string (already surfaced as `packageStatus.source`), so no package-manifest version field is needed — the reviewer confirmed the seam exists (`internal/cli/pi.go:689-750`).

## Expected surface and tolerance (unchanged from seed baseline)

+100 net LOC across 4 files; tolerance ±60 net / ±2 files. `internal/cli/pi.go`, one new source-derivation helper file, `pi_frontdoor_test.go`, Pi live test wiring. Semantics declared: stored format (`packages` entry gains `@<ref>` for release-shaped binaries) and front-door runtime behavior (one repair attempt, launch refusal on remaining mismatch). Command grammar, stored schema, and authority unchanged.

## Acceptance criteria and proposed proof

- **AC-1 (VALUE):** live Pi journey through the ordinary installed front door, no dev override, started with package absent/unpinned/wrong-line; exactly one pinned `pi install`, then `status --boot` launch-ready. Independent baseline (moves the wrong way without the fix): repair removed → package stays missing/unpinned → FO binary gate aborts (the observed v0.27.2 failure). Retained artifacts are the anti-tautology baseline for AC-2..AC-5.
- **AC-2:** literal-source command-behavior assertions per identity (linker-stamped → own stamp; proxy-tagged → own tag; dev → no install; `--plugin-dir` → plugin dir, never the release source); falsifier flips any literal.
- **AC-3:** table tests over stable/prerelease/unstamped/proxy-tagged identities with literal expected sources.
- **AC-4:** failed/ineffective repair → exactly one install attempt, one recheck, no launch, actionable error; falsifiers named.
- **AC-5:** no-repair/no-clobber for non-git entries, dev builds, and declared dev overrides, each with a named falsifier.

## Risk evidence (spike)

Spike verified against the real shipped Pi package manager on this machine (`pi-coding-agent` dist `core/package-manager.js` + `utils/git.js`): `pinned: Boolean(ref)`; unpinned git sources float via origin-HEAD update; `@ref` install clones+checkouts and `addSourceToSettings` rewrites the existing entry (match key ignores ref). Residual risks stated in the body, not silently: older Pi package managers may not honor `@ref` (live run + remedy wording is the fallback); concurrent front-door launches could race the settings rewrite (accepted, single-user launch path).

## Journey record

- Cycle 1: fresh-context independent review — NEEDS REVISION (proxy-build floating gap; unimplementable manifest-version trigger; unspecified non-git/override repair semantics). Routed to the ideation worker.
- Cycle 2: revision commit `4a018730e` resolved all three; reviewer re-run verdict **RESOLVED — ready for gate**.

## FO recommendation

APPROVED for advancement to implementation.
