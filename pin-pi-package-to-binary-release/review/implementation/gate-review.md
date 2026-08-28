# Gate review — implementation (returned): pin-pi-package-to-binary-release (qc)

## Delivery summary

Commit 178661933 on spacedock-ensign/pin-pi-package-to-binary-release: pin-ref derivation helper (`pi_package_source.go`: linker stamp → build-info semver tag → dev sentinel), pinned install call site in `pi.go`, one-repair-attempt front door with recheck and launch refusal, AC-5 suppression. Behavior tests for AC-2..AC-5 with named falsifiers (`pi_package_repair_test.go` + `pi_frontdoor_test.go`); AC-1 no-override live journey wired and registry-registered. FO independently verified: full suite green except the pre-existing env-marker failure (reproduced on main), race clean, gofmt clean.

## Surface deviation (captain ruled: design reset)

Approved surface: +100 net / 4 files, ±60/±2. Actual: ~664 net insertions across 6 files (pi.go 54, helper 150, repair tests 335, frontdoor tests 9, live wiring 96, registry 18). File count within tolerance; net LOC ~4x over ceiling. Growth driven by cycle-2 reviewer-demanded branches (proxy identity, AC-5, wrong-ref tests) after the estimate was approved; the surface was never re-baselined at the ideation gate.

## Captain ruling (2026-08-28 chat)

Option 2 — design reset: return through ideation to re-baseline the expected surface against the implemented reality, re-approve, then re-confirm implementation against the re-baselined surface.

## FO recommendation

REVISE: reset to ideation with the sole scope of re-baselining the expected surface (net LOC/files) to the as-implemented figures, keeping the design, ACs, and code unchanged. No code changes requested; the deliverable stands.
