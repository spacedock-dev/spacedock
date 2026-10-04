# Validation gate — pi-live-lane-pin-refresh (mh)

## Outcome

The layer is validated after a criteria amendment, and it is ready for the merge ceremony.

## Candidate

- Tip `5a2ad16e9` on `spacedock-ensign/pi-live-lane-pin-refresh`, pushed and level with origin.
- Stack: `#822` (271, base `main`) -> `#820` (9w, base 271) -> `#816` (mh, base 9w).
- Size: 19 files, **net +1078**. Port `internal/pilive/pilive.go` 427 + `pilive_test.go` 251 = 678 lines.

## Evidence

- **Lane run `37189752489` at the frozen tip: success.** `offline` job success, which runs
  `go test ./...` on the runner; `pi-live` job success, `DONE 17 tests` for the 17 common
  journeys plus the front-door smoke.
- **Race, local, uncached:** `go test ./... -race -count=1` -> 21/21 packages ok, no failures.
- **Exception register**, read from the run's `pi-coverage-detail.jsonl`: 5 XFAIL engagements
  and 2 XPASS alerts. All 17 journeys passed.
- **AC-6** carries the independent live artifact: the doctor and durable smoke grade
  (`boot_contract` true, isolated discovery exposing both tools, package-root env absent).

## Amendment

Validation returned **AC-1 to AC-4 as not met**, because each named an artifact the captain's
simplification rounds removed — the registry-comparison Go guard, the Go structural guard with
its falsifiers, the `verified_pack` wiring with `PI_SUBAGENTS_INTEGRITY`, and the version-naming
workflow comment and doc lines — and AC-1's literal named `1.0.0` where the shipped pin is
`1.0.2`. The captain directed the amendment; each criterion is restated against the shipped
mechanism and records the guarantee it loses.

This layer therefore claims: the pinned family, the manifest-resolved checks, the integrity
comparison, and the single pin source. It does **not** claim a standalone pin-revert falsifier,
an independent structural guard, an independent integrity falsifier, or a standing
doc-versus-pin check.

## Residuals recorded, not hidden

- **AC-2, AC-3, AC-7, AC-8 and AC-10 rest on demonstrations authored by the same change.**
- **The `default-headless-gate-stop` binding engaged on the wrong fault.** It observed
  `implementation-worker-not-dispatched`, which its own comment says it does not cover, and
  `gradeLive` never checks an observed code against the bound owner.
- **Two bindings are stale:** `smallest-sufficient-mechanism` and `keep-moving-posture` pass
  with `observed=[]`.
- **One local offline flake** in a test this layer does not touch, under a full disk; the same
  code is green on the runner.

## Question

Approve this layer's validation, so `merge guard` may proceed and the train merges in order
`#822` -> `#820` -> `#816`?
