# Sprint 019x — staff readiness review

The staff review is a **readiness gap-analysis run before the Commander drives**: it
answers "can the planned work produce the deliverable, and what's missing?" grounded in
the real entity bodies, not a guess. It is a **mechanism** — reusable per sprint — not a
bespoke write-up.

## The mechanism (reusable per sprint)

1. **Who runs it.** The Shaping FO dispatches an INDEPENDENT reviewer — not the Shaping
   FO, not any member's implementing ensign. Read-only; it produces findings, never a
   deliverable.
2. **What it is given.** This sprint's `index.md` (goal / DoD / deliverable) plus the
   full body of every `ready` member:
   ```bash
   spacedock status --workflow-dir docs/dev --where sprint=019x-pre-flip-cleanups --where 'sprint-readiness != defer'
   # then read each member's entity file
   ```
3. **What it checks (the gap-analysis prompt).**
   - **Producibility** — does the planned member work, taken together, actually produce
     the DoD? Is any DoD clause unowned by a member?
   - **Build-readiness** — is each member's design concrete enough to drive (ideation
     done, or design self-evident)? Name members that still need a gating ideation.
   - **Grounded risk** — for any member whose soundness rests on an unverified mechanism,
     was the riskiest unknown exercised first? (The workflow's spike discipline.)
   - **Coherence** — do cross-member contracts agree (shared files, overlapping surfaces,
     ordering dependencies)? Will two members collide in the same code?
   - **Integration-test reality** — can the sprint-wide DoD test actually fail today, or
     is it a tautology?
4. **What it produces.** A findings list split into **BLOCKING** (must close before the
   Commander drives) and **NON-BLOCKING** (drive anyway, track), plus a one-line
   **go / no-go** readiness verdict. Findings route back to the Shaping FO, who closes
   blockers (runs the missing ideation, resolves a coherence clash) before packaging the
   Commander dispatch.

## Findings — run 1 (2026-06-08, independent Plan reviewer)

**Readiness verdict: NO-GO** — close B1 + B2 before the Commander drives. `e30` and the
`qy` migration-check half are ready now and do not collide.

### BLOCKING (must close before the Commander drives)

- **B1 — `qy`'s gofmt half is stale/unverifiable on this base.** The two files qy names as
  gofmt-dirty (`internal/status/external_proof.go`, `no_yaml_silent_drop_test.go`) are
  **clean** on the current tree, and there is no local `next` branch to confirm the dirty
  state on qy's stated target. qy's migration-check half is real (the test genuinely
  fails — see N1), but the gofmt half may be a phantom deliverable. **Action:** re-ground
  qy's gofmt claim on the real drive branch; if clean, narrow qy to the migration-check
  fix and note DoD-3 (`gofmt -l` empty) is already satisfied — and therefore a tautology
  to "prove."
- **B2 — `jh` is not build-ready: tautological AC + unspiked mechanism.** jh's AC-3 ("Codex
  prompt stays `Read /tmp/...`, no `Skill(...)`") is **already covered** by
  `internal/dispatch/build_codex_host_test.go:45,53` — it pins nothing new. The real work
  (AC-1/AC-2: runtime instructions must not require a flag form the minimum accepted
  binary rejects) rests on an **undecided mechanism** — jh's own stage gates leave the
  approach a fork (raise the compat gate vs probe `--print-schema`/`--help` vs document a
  stdin fallback), and the failing case is the *installed 0.19.4 binary*, not in this
  tree. **Action:** a gating ideation pass + a spike against a real old binary before the
  Commander drives jh (or pull jh from this sprint).
- **B3 — the DoD is structurally cross-actor.** DoD-1 ("every `ready` member merged") and
  DoD-4 (`goreleaser check`, owned only by `78`) cannot be closed by the Commander's
  backlog members alone — `nb`/`xn`/`78` are captain gates. The drive split must sequence
  the captain clearing those gates in the same window, or the sprint cannot reach DoD even
  with qy/e30/jh done.

### NON-BLOCKING (drive anyway, track)

- **N1 — the DoD integration test is genuinely RED today (not a tautology).**
  `go test ./internal/status/` fails at `migration_check_test.go:113` (bare-YAML
  `session-date` in `_debriefs/*.md` read as string by one path, RFC3339 by the other);
  full `go test ./...` fails ONLY on `internal/status`. qy's migration-check fix is
  concrete and build-ready (exclude `_debriefs/` from the walk @ line 61, or normalize
  bare-date scalars).
- **N2 — `e30` and `jh` do NOT collide** despite both touching `internal/dispatch/build.go`:
  e30 at lines ~567–573 (path/comment), jh at `fieldsFromBuildFlags` ~124 — different
  functions ~440 lines apart. The pre-review coherence flag was over-cautious; refuted
  with evidence. Drive in any order.
- **N3 — `e30` is the most build-ready backlog member** — design fully specified (mirror
  the adjacent `derivedName` validation onto `teamName` + length cap + fix the line-567
  comment), both ACs name outside-the-body unit tests. No ideation pass needed.
- **N4 — `nb`/`xn`/`78` are genuinely done-pending-gate.** 78 PASSED (re-confirmed against
  `.goreleaser.yaml`). nb cycle-3 PASSED (one non-blocking residue: README omits the
  Codex/Pi "experimental" caveat install-journey carries — captain wave-or-fix; note the
  captain already chose to list Codex/Pi without that caveat). xn fixed the zsh-glob bug
  (`e59c0861`); its AC-3 live drive is the one open item, captain-supplied.
- **N5 — `5h0` deferral is properly recorded** (DoD-6 met): body + frontmatter both carry
  the defer + the #315 blocker.

### What the Shaping FO must do before packaging the Commander dispatch

1. **B1:** re-ground qy on the real drive branch; narrow scope if the gofmt half is moot.
2. **B2:** run jh's gating ideation + a version-skew spike (or pull jh from this sprint).
3. **B3:** confirm with the captain that the nb/xn/78 gates clear in the 0.19.7 window.

Ready to drive now without further shaping: **`e30`** (full) + **`qy`'s migration-check half**.
