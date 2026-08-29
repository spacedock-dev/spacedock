Gate review: Pin the installed Pi package to the launcher release — validation
Recommend approve.
Reviewed snapshot: validation cycle-3 report at entity `pin-pi-package-to-binary-release` (qcwrzza9xkr5kfwmdmvbqhkv); deliverable tip `284299e3b` on branch `spacedock-ensign/pin-pi-package-to-binary-release`, PR #782 stacked on #780.

Evidence for validation:
- AC-1 (VALUE): hand-executed front-door runs against a fixture pi-home seeded with the incident state — unpinned and wrong-line starts each perform exactly one install of `git:github.com/spacedock-dev/spacedock@v0.27.2`, rewrite settings to the pinned source, exit 0, and launch once. Retained run artifacts under the entity.
- AC-2: release-shaped identities (stamped, proxy-tagged via build-info semver) expose the literal pinned source; dev-sentinel and `--plugin-dir`/`SPACEDOCK_REPO_ROOT` overrides avoid the release source. Pure-function derivation/classifier/repair-decision tables pass.
- AC-3: stable (0.27.2), prerelease (0.28-pre), proxy-tagged, checkout-dev, clean pseudo, and dirty pseudo (`v0.28.0-pre0.0.20260828165724-81e3386e8234+dirty`) identities retain their literal expected derivations; the dirty-pseudo row was added by the correction commit.
- AC-4: missing and wrong-line install failures, plus reported-success remaining mismatch, each spend one attempt, exit 1, make zero launches, preserve settings bytes, and print actionable refusal naming `spacedock doctor --host pi`; successful repair still launches exactly once.
- AC-5: `file:`, dirty dev, `--plugin-dir`, and `SPACEDOCK_REPO_ROOT` cases make zero installs and preserve settings bytes; healthy pinned remains a zero-install launch.
- Checks: focused pure-function tables PASS; isolated full suite PASS (CLI 112s, ensigncycle 168s); isolated race PASS (CLI 145s, ensigncycle 205s); gofmt clean. PR #782 standard CI green (docs, install-e2e, offline, build). Live-lane failures on #782 reproduce identically on stack base #780 four hours earlier (codex `default-headless-gate-stop` gate-hold flake; claude `TestLiveCommon*` model-variance) — neither is a pi-frontdoor-pin test, neither is this branch's regression.
- Surface: +310/−1 = +309 net across `pi.go` +64/−1, `pi_package_source.go` +150, `pi_package_source_test.go` +96 — within the +295±80 net / 3±1 re-baselined tolerance. No removed scaffolding (flow suite, live journey, registry) restored.

Reviewer findings:
- None material. Both validation/2 Material outcome defects (spent repair could launch on remaining mismatch; dirty pseudo-version misclassified) are RESOLVED by observed terminal behavior at tip `284299e3b`.
- Deferred risks (accepted, non-blocking): (1) older Pi package managers might not honor `@ref` — trigger is a Pi version outside the locally spiked supported setup; promote to material if a supported Pi version fails to reconcile the pinned source. (2) concurrent `spacedock pi` processes can race the settings rewrite — trigger is simultaneous launch in a single-user path; promote if concurrent use becomes supported or an observed race leaves persistent wrong-line state.

Decision: approve to enter done/approved-awaiting-merge and drive the merge guard to merge PR #782 and archive the entity.
