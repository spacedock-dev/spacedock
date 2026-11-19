# Sprint 0198 — staff readiness review

Same mechanism as `019x` (see [`../019x-pre-flip-cleanups/staff-review.md`](../019x-pre-flip-cleanups/staff-review.md)):
before the Commander drives, INDEPENDENT reviewers do a readiness gap-analysis over the
ideated sprint — producibility, build-readiness, grounded-risk, coherence,
integration-test reality — output as BLOCKING / NON-BLOCKING + a go/no-go.

## Findings — run 1 (2026-06-08, two-lens preflight: build-readiness + coherence)

**Verdict: NO-GO → after the DoD#4 fix below, ONE blocker remains (the survey model
contradiction).** Every member's cited `file:line` claims were verified against real repo
bytes; every ideated member spiked its riskiest unknown against real bytes/CLI.

### BLOCKING

- **B1 — DoD#4 (live drive) was unowned [producibility].** The DoD demanded qa's behavior
  be proven by a live drive, but qa's ACs/test-plan are entirely offline (RunDoctor fixture
  tests + a frontdoor gate-parity seam + an AC-3 text check). **Closed** by reassigning DoD#4
  to a **captain-run sprint-acceptance live drive** (index.md updated, mirroring 019x AC-3) —
  qa's offline ACs prove the mechanism; the captain live drive proves the messages.
- **B2 — the survey group does NOT converge on one agentsview `project` model [coherence].**
  `69`'s spike proved agentsview v0.32.1 keys `project` by **git-root basename** (this repo's
  root, every `.worktrees/*`, and `.spacedock-state` all key to ONE `project`). But the
  shipped survey still teaches the opposite — `SKILL.md:64` ("agentsview keys each session by
  the basename of its working directory, so a subdir / worktree / split-root state dir each
  get a DIFFERENT key") and `queries.sql:27-43` (the `scoping`/`folded_keys` prefix-union
  rationale) rest on the per-cwd-basename premise. So `folded_keys` is ~always 1 for one repo
  and the prefix-union does far less than the prose claims. `69` flagged this and scoped the
  correction OUT; meanwhile `1p` reshapes the same SKILL.md and `4t` edits it too — three
  members polishing around a now-known-false core claim (a model `xn` shipped in 0.19.7).
  **Open — needs the captain's call:** correct §64 + the scoping/folded_keys prose to the
  git-root-basename reality (recommend folding it into `69`, the survey-scoping member that
  found it) **or** defer with a tracked follow-up — before any survey member edits around it.

### Per-member readiness

| Member | Readiness | Note |
|--------|-----------|------|
| `nzb` gate-release-on-e2e | **drive-now** | spiked live; require its `--dry-run` vs real `gh run list` as the not-just-prose bar (true e2e is the flip's) |
| `z9` codex-auto-install | **drive-now** | spiked live; must **fix** the now-false comments/error-strings it builds around (`host_exec.go:271-273`, `:32-34`; `frontdoor.go:314-316`), not add around them |
| `4t` agentsview-detect | **drive-now** | one-line probe swap; reproduced on a real sibling binary |
| `kb` migration-check | **drive-now** | skip-ideation-ready; add a POSITIVE walk-skip assertion AND verify the `checked == 0` Fatal guard stays non-vacuous after the prune |
| `qa` binary/version UX | needs-fix → **resolved** | mechanism sound; DoD test fails-today as required; B1 closed via the captain live-drive step |
| `69` codex-cwd | **needs-fix** | blocked on B2 (the §64 model correction) |
| `1p` scaffold-fact | **needs-fix** | thinner than the kb/78 bar — firm the exact SKILL.md edit (how to drop the taxonomy while keeping the "recovered" fact); also entangled in B2 |

### Coordination (baked into the Commander dispatch)

- **qa BEFORE z9** — both edit `frontdoor.go` + `host_exec.go` on DISJOINT functions (no
  semantic collision); sequence to avoid merge-adjacent textual overlap; qa is the headline +
  shared-file owner.
- **Survey order:** `4t` (line 27, isolated) → **resolve B2** → `69` → `1p` (`69`+`1p` both
  edit the step-4 report-template fence ~141-169 → land sequentially or hand-merge).
- **`kb`** independent — land any time, but before the survey live-drives (orphaned
  `scaffolds/` fixtures gone). It is the de-facto precondition for DoD#2.
- **`nzb`** shares no files — parallelizable, but merge before the captain-gated release.
- **Validation budget:** `69`/`4t`/`1p` all bottom out on a live drive — the DoD live-drive
  pass must exercise all three observable changes in one pass (or run thrice).
