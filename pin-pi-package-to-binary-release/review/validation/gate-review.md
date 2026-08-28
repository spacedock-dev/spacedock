# Gate review — validation: pin-pi-package-to-binary-release (qc)

## Deliverable

Commit 178661933 on spacedock-ensign/pin-pi-package-to-binary-release: release-pinned Pi package source (ordered derivation: linker stamp → build-info semver tag for proxy builds → dev sentinel), one-repair-attempt front door with recheck and launch refusal, AC-5 no-repair/no-clobber suppression, AC-2..AC-5 command-behavior tests with named falsifiers, AC-1 no-override live journey wired and registry-registered.

## Independent verification (validation worker; FO re-verified durable state)

- 10/10 AC behavior tests PASS; falsifier audit confirms literal-pinning per identity (stamped, proxy-tagged, dev, plugin-dir), wrong-ref repair, one-install/one-recheck/no-launch, no-clobber
- Full suite ok and race ok with ambient markers unset (cli 181s, ensigncycle 251s); gofmt clean
- The sole failure (TestVersionAmbiguousMarkersExitZero) verified identical on origin/main — pre-existing environment-marker failure, not this branch's regression
- Repair-path code audit matches the approved design (pseudo-version guard so @next resolutions float; repoRoot suppression covers --plugin-dir and SPACEDOCK_REPO_ROOT; no clobber path)
- AC-1 live journey (pi-front-door-pinned-package, registered docs/runtime-live-ci-registry.md:336) rides the CI pi-live lane — open by design, needs CI-E2E-PI approval

## Journey record

Backlog gate approved → ideation (2 cycles, review-confirmed RESOLVED) → implementation → captain design reset on the surface breach (round implementation/1, 4 entries) → surface re-baselined +661 net / 6 files ±40/±1, reviewer-confirmed exact → validation worker independently verified. Entity nonterminal.

## FO recommendation

APPROVED — advance to done (approved-awaiting-merge); PR + merge guard next.
