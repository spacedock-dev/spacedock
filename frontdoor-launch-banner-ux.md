---
id: cmxchb8y1y2m455xhx7ce87g
title: Launch-banner UX — first-officer framing, status-command overload, multi-workflow limbo
status: validation
source: "FO + captain launch-banner review (2026-06-08), following yq (frontdoor-launch-ux). The 3-line pre-launch banner (frontdoor.go:139-141) reads confusingly on the new-user / multi-workflow path."
started: 2026-06-08T23:37:46Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-frontdoor-launch-banner-ux
issue:
sprint: 0200-flip
group: frontdoor
sprint-readiness:
mod-block:
---

A followup on yq (frontdoor-launch-ux): the pre-launch banner (`launchBanner`, `internal/cli/frontdoor.go:139-141`) is clean on the common single-workflow path but wobbles on the new-user / multi-workflow path. Three issues worth fixing, worst first.

## Problem

The banner emits three lines:
1. `spacedock {Version} · first officer launching {host}`
2. `Workflow: {detectedWorkflow}`  (e.g. `5 workflows detected (run spacedock status to pick)`)
3. `{host} is starting as your first officer; run spacedock status inside the session for the queue.`

1. **The "first officer" metaphor flips within three lines.** Line 1 ("first officer **launching** claude") reads as *spacedock is the FO launching claude*; line 3 ("claude **is** your first officer") reads as *claude is the FO*. A new user can't form a stable model. Intended meaning is line 3; line 1 shouldn't call the launcher a "first officer." Candidate: `spacedock {v} · launching {host} as your first officer`.
2. **`spacedock status` is overloaded and undercuts the value prop.** The same command appears twice for two jobs — "run `spacedock status` to pick" (line 2) and "run `spacedock status` ... for the queue" (line 3). Worse: the whole pitch is "claude is your first officer" — the FO runs `status` *for* you. Telling the user to run it themselves (and "inside the session," ambiguous in a chat) contradicts that. Candidate for line 3: *"...ask your first officer for the queue."*
3. **Multi-workflow "pick" then launches anyway.** "N workflows detected (run `spacedock status` to pick)" announces an unresolved choice, then claude starts without one — leaving the user unsure which workflow the FO is on. Per the FO contract the disambiguation actually happens in-session (multiple → the FO presents the list on its first turn), so the banner instruction is redundant + misleading. Either say nothing actionable, or *"your first officer will help you pick."*

## Out of scope

The single-workflow happy path (clean today). The host's own banner ordering.

## Minor / to confirm at ideation

- The displayed `Version` is the compiled-in constant — confirm it reflects the real installed release (a dev build shows a stale `0.19.0`); a wrong version is a poor first impression and is exactly the version/channel correctness the 0.20.0 flip is about.
- `Workflow:` (singular label) + `N workflows detected` (plural content) reads slightly off.

## Proposed design

The captain wants to SEE the proposed front-door output before any code. The headline deliverable is the BEFORE/AFTER matrix below (mode × host), with the exact wording, plus behavioral ACs each pinned by an `internal/cli` test over rendered output.

### Where the banner is emitted (read from code, not guessed)

- **claude / codex.** `launchBanner(host, dir, stderr)` (`frontdoor.go:138-142`) emits the 3-line banner to **stderr**, called from `runClaude` (`:290`) and `runCodex` (`:448`) AFTER the gate passes (or after a NoPluginFound auto-install), BEFORE the host launch, and suppressed on a resume (`!resume`).
- **pi.** `runPi` (`pi.go:65-94`) emits **no banner at all** — no version line, no workflow line, no orientation line. On the success path it prints nothing of its own; it just builds the `pi --extension … --skill …` argv and execs. This is a real gap: pi launches the FO with zero pre-launch orientation. (Captured by adding a pi case to the AFTER tests — see AC-7.)

### How each MODE resolves per host (verified, not inferred)

The four modes are properties of the **contract gate**, which only claude/codex have. Pi has a different readiness model (file-presence `checkPiRuntime`, no contract-version compare, no auto-install), so two of the four modes do not exist for pi.

- **fresh-install (no plugin).** claude/codex: `gateHost` returns `NoPluginFound` → caller prints `Installing the {host} plugin…` (stderr), runs `ops.Install`, then proceeds to banner+launch (`frontdoor.go:263-277`, `:422-435`). pi: not a contract state — if the runtime files are absent, `runPi` prints `spacedock pi: Pi runtime is not ready…` + a doctor report and exits 1 (`pi.go:74-78`); there is no auto-install.
- **upgrade / version-skew (mismatch).** claude/codex: `gateHost` prints the `mismatchMessage` (`contract.go:155-161` — `Spacedock version mismatch: binary {bin}, plugin {plugin}. …` + per-class remedy) and returns non-Compatible → caller fails fast, exit 1, **no banner** (`frontdoor.go:278-282`). pi: no contract-version axis, so "skew" has no analog; pi only checks whether the skill/extension files are present.
- **compatible-launch.** claude/codex: the 3-line banner then launch. pi: nothing today.
- **`--no-install` refuse.** claude/codex: `gateHost` returns `NoPluginFound`, `fd.noInstall` true → caller prints `noPluginRemedy(host)` (`frontdoor.go:127-132`) and exits 1, **no banner**. pi: `--no-install` is not a pi flag — `parsePiFrontDoorArgs` rejects it as `unknown flag: --no-install`, exit 2 (verified by a throwaway probe). N/A for pi.

This task scopes the **banner wording** (the `launchBanner` output). The gate/install/refuse messages (mismatch remedy, `Installing…`, `noPluginRemedy`) are out of scope for rewording — they are not the banner and were settled by the contract/yq work. The banner only appears on the compatible-launch and post-auto-install paths.

### BEFORE / AFTER matrix (mode × host)

Banner = the `launchBanner` 3-line block on stderr. "—" = no banner emitted on that path.

| Mode | Host | BEFORE (today) | AFTER (proposed) |
|------|------|----------------|------------------|
| compatible, **single** workflow | claude | `spacedock 0.19.0 · first officer launching claude`<br>`Workflow: docs/dev`<br>`claude is starting as your first officer; run \`spacedock status\` inside the session for the queue.` | `spacedock {v} · launching claude as your first officer`<br>`Workflow: docs/dev`<br>`claude is your first officer — ask it for the queue and next steps.` |
| compatible, **single** workflow | codex | same 3 lines, `codex` | `spacedock {v} · launching codex as your first officer`<br>`Workflow: docs/dev`<br>`codex is your first officer — ask it for the queue and next steps.` |
| compatible, **multi** workflow | claude | `spacedock 0.19.0 · first officer launching claude`<br>`Workflow: 2 workflows detected (run \`spacedock status\` to pick)`<br>`claude is starting as your first officer; run \`spacedock status\` inside the session for the queue.` | `spacedock {v} · launching claude as your first officer`<br>`Workflows: docs/dev docs/user-testing`<br>`claude is your first officer — ask it for the queue and next steps.` |
| compatible, **none** detected | codex | `spacedock 0.19.0 · first officer launching codex`<br>`Workflow: none detected (launching anyway)`<br>`codex is starting as your first officer; run \`spacedock status\` inside the session for the queue.` | `spacedock {v} · launching codex as your first officer`<br>`Workflow: none detected`<br>`codex is your first officer — ask it for the queue and next steps.` |
| compatible | pi | *(no banner)* | the 3-line banner, `pi` (closes the pi gap) |
| fresh-install (no plugin) | claude/codex | `Installing the {host} plugin…` then the BEFORE banner | `Installing the {host} plugin…` then the AFTER banner (banner change only; install line unchanged) |
| fresh-install (runtime not ready) | pi | `spacedock pi: Pi runtime is not ready…` + doctor report, exit 1 | unchanged (out of scope — not the banner) |
| upgrade / version-skew | claude/codex | mismatch message, exit 1, no banner | unchanged (out of scope — not the banner) |
| upgrade / version-skew | pi | N/A (no contract axis) | N/A |
| `--no-install` refuse | claude/codex | `noPluginRemedy(host)`, exit 1, no banner | unchanged (out of scope — not the banner) |
| `--no-install` refuse | pi | `unknown flag: --no-install`, exit 2 | N/A (pi has no `--no-install`) |
| resume | all | *(banner suppressed)* | unchanged (still suppressed) |

### The three problems, and how the AFTER column fixes each

1. **First-officer metaphor flip.** BEFORE line 1 "first officer launching claude" reads as *spacedock-is-FO*; line 3 "claude is your first officer" reads as *claude-is-FO*. AFTER line 1 is `launching {host} as your first officer` — the launcher is never called a first officer; only the host is, consistently with line 3. One stable model: the host is your first officer; spacedock launched it.
2. **`status` overloaded + self-serve contradiction.** BEFORE prints `spacedock status` twice and tells the user to run it themselves — contradicting "the host is your first officer (it runs status for you)." AFTER removes BOTH `spacedock status` instructions. Line 3 becomes `ask it for the queue and next steps.` The pitch holds: you ask your first officer; it runs status.
3. **Multi-workflow "pick then launches anyway" limbo.** BEFORE announces "N workflows detected (run `spacedock status` to pick)" then launches with no pick. AFTER (captain wording call) just LISTS the detected workflow paths, space-separated, with no count, no helper sentence, and no command: `Workflows: docs/dev docs/user-testing`. This still resolves the problem — no actionable-but-ignored instruction — and the user sees exactly what is there; the FO disambiguates in-session on its first turn (per the FO contract).

Minors:
- **Singular label / plural content.** BEFORE always says `Workflow:` even for N workflows. AFTER uses `Workflow:` for the single + none cases and `Workflows:` for the multi case, so label and content agree.
- **Stale Version display.** `Version` (`cli.go:28`) defaults to the literal `"0.19.0"` and is only overwritten by the goreleaser ldflags stamp (`.goreleaser.yaml:37`, `git describe`). A dev/`go install` build (no ldflags) prints `0.19.0` regardless of the actual commit — impersonating a stale release. Fix: default `Version` to a non-release sentinel (`"dev"`) so an unstamped build reads honestly as a dev build. A stamped release overwrites it as before. This is a real, checkable change (AC-6).

### Wording notes

- Line 1 separator stays ` · ` (matches today). `{v}` is the `Version` value (a real release on a stamped build, `dev` on an unstamped one).
- Line 2 single keeps `Workflow: {path}` (e.g. `Workflow: docs/dev`), unchanged.
- Line 2 multi (captain wording call): `Workflows: {space-separated workflow paths}` — just LIST the detected workflow paths (e.g. `Workflows: docs/dev docs/user-testing`). No count, no helper sentence, no command. The captain chose this terse path-list form over `{N} found — your first officer will help you pick`: just the facts, the user sees what is there, the FO picks in-session. Still resolves Problem 3 (no actionable-but-ignored instruction).
- Line 2 none (captain wording call): `Workflow: none detected` — the `(launching anyway)` parenthetical is dropped for terseness.
- comm-officer polish (Strunk) had been applied earlier (line 3 `what to do next` → `next steps`, kept); its line-2 suggestions are superseded by the captain's terse path-list call above. Line 3 stays `{host} is your first officer — ask it for the queue and next steps.` The result is terse, each line stands alone, and satisfies the three constraints (consistent metaphor, no self-serve `status`, clean single-workflow path).

## Acceptance criteria

Each AC is verified by an `internal/cli` test over the **rendered** banner output (the bytes `launchBanner` / `runPi` writes), never a source-grep. The independent expected value is the proposed wording above (a renamed-but-equivalent branch must still fail the assertion).

- **AC-1 — consistent metaphor, no launcher-as-FO.** The rendered banner never contains the phrase `first officer launching` (the flipped line-1 framing); line 1 reads `launching {host} as your first officer`.
  *Verified by:* a test asserting the banner output (claude and codex) contains `launching {host} as your first officer` and does NOT contain `first officer launching`.
- **AC-2 — no self-serve `spacedock status` in the banner.** The rendered banner contains no `spacedock status` instruction on any of the single / multi / none paths; line 3 reads `{host} is your first officer — ask it for the queue and next steps.`
  *Verified by:* a test asserting the banner output for all three workflow cases does NOT contain `` `spacedock status` `` and DOES contain `ask it for the queue`.
- **AC-3 — multi-workflow lists the detected paths, with no count and no pick instruction.** On the multi-workflow path, the rendered banner line 2 is `Workflows: {space-separated workflow paths}` (the detected workflow paths, space-separated, in discovery order) and contains no count, no `will help you pick`, no `run \`spacedock status\` to pick`, and no bare-command pick instruction.
  *Verified by:* a test that builds a 2-workflow git-repo fixture (`docs/dev`, `docs/user-testing`), renders the banner, and asserts line 2 equals `Workflows: docs/dev docs/user-testing` and that the banner does NOT contain `found` / `will help you pick` / `to pick)` / `spacedock status`. The fixture's two workflow paths are the independent expected value.
- **AC-4 — singular/plural label agreement.** The single-workflow and none paths render `Workflow:`; the multi path renders `Workflows:`.
  *Verified by:* a test asserting line 2 is `Workflow: docs/dev` for the single fixture, `Workflow: none detected` for the bare fixture, and starts with `Workflows: ` (plural) for the 2-workflow fixture.
- **AC-5 — single-workflow happy path stays clean and terse.** The single-workflow banner is exactly the three proposed lines, in order, with no extra lines.
  *Verified by:* a test asserting the full rendered single-workflow output equals the three expected lines joined by `\n` (a byte-exact golden assertion, the strongest form, so any drift fails).
- **AC-6 — unstamped build does not impersonate a release.** An unstamped (`go build` / `go install`) binary's `Version` is not a release-shaped semver; it reads as a dev build.
  *Verified by:* a test asserting the package default `Version` (the value before any ldflags stamp) equals the dev sentinel and does not match a `^\d+\.\d+\.\d+$` release pattern. (A stamped release still overwrites it; the ldflags path is exercised by the existing release pipeline.)
- **AC-7 — pi emits the banner.** `spacedock pi` renders the same 3-line banner before launch (closing the pi no-orientation gap), suppressed on the pi resume path if one exists.
  *Verified by:* a test driving `runPi` with a ready stub runtime and asserting the rendered stderr contains `launching pi as your first officer` and `pi is your first officer — ask it for the queue`, and that the launch seam is still reached.

### Risk / spike

No spike needed. The changes: (1) reword `launchBanner`'s line 1 and line 3; (2) change `detectedWorkflow`'s multi branch to return the space-joined workflow paths instead of the `"{N} workflows detected (…)"` count string, and its none branch to return `none detected` (drop the `(launching anyway)` suffix) — the discovery already returns the workflow dir list (`status.DiscoverWorkflows`), so the multi branch joins `workflowLabel(repoRoot, wf)` over that list rather than counting it; no new discovery mechanism; (3) add the existing `launchBanner` call site to `runPi`; (4) default `Version` to the `dev` sentinel. The workflow-detection mechanism (`DiscoverWorkflows`/`workflowLabel`, repo discovery), the gate, and the launch seam are all already proven by the existing `launch_banner_test.go` / `frontdoor_test.go` / `launch_parity_test.go` suites — and `workflowLabel` is already exercised per-workflow for the single case, so the multi list reuses a proven renderer. The proposed-output BEFORE column was captured by exercising the real `launchBanner` (single / multi / none) and the pi `--no-install` path was probed live, both via throwaway tests — recorded above, not inferred from source. The `Version` sentinel (AC-6) composes the proven goreleaser stamp path (the ldflags `-X` already works for releases); only the default literal changes.

### Test plan

- **Fixtures, fast.** All ACs are `internal/cli` Go tests using the existing `commissionWorkflowAt` / `gitRepoFixture` helpers in `launch_banner_test.go` (single / multi / none repos) plus the `fakeHost` (claude/codex launch) and a pi stub runtime (AC-7). No live workflow, no network.
- **Cost:** low — extends the existing `launch_banner_test.go` patterns; sub-second.
- **Golden vs substring:** AC-5 is a byte-exact golden on the single-workflow happy path (the strongest assertion). AC-1/2/3/4/7 are presence/absence assertions over rendered output where a partial match is the right granularity. AC-6 is a property assertion over the package default plus a regex non-match. None is a source-grep; every expected value is the independent proposed wording, so a renamed branch or a valid paraphrase fails.
- **Regression:** the existing `TestLaunchBannerNamesDetectedWorkflow` / `…ReachesStderrBeforeLaunch` / `…SuppressedOnResume` will need their expected strings updated to the AFTER wording (they currently assert the BEFORE lines — e.g. the multi sub-test asserts `2 workflows detected (run \`spacedock status\` to pick)`, which becomes the space-joined path list, and the none sub-test's `none detected (launching anyway)` becomes `none detected`); that update IS part of the implementation, not extra scope.

## Stage Report: ideation

- DONE: Produce the PROPOSED launch-output matrix the captain asked for (redesigned front-door output for each MODE × HOST, before/after, addressing the three filed banner problems; ideation/design, no implementation).
  BEFORE/AFTER table under `## Proposed design` covers compatible single/multi/none, fresh-install, upgrade-skew, --no-install, resume × {claude, codex, pi}; deliverable is matrix + ACs, no code changed.
- DONE: Map the CURRENT front-door output for each cell first (read the real code, don't guess).
  Read frontdoor.go (`launchBanner`:138, gate arms :259-283/:418-441, `noPluginRemedy`:127), init.go, host_exec.go, pi.go (`runPi`:65 — no banner), contract.go (`mismatchMessage`:155); BEFORE column captured by exercising real `launchBanner` (single/multi/none) + a live `pi --no-install` probe (both throwaway tests, removed).
- DONE: Propose the REVISED output per cell fixing the three filed problems + minors; present as a BEFORE/AFTER table (mode × host); give ACs each verified by an `internal/cli` test over rendered output, never a source-grep.
  Problem 1 (metaphor flip) → line 1 `launching {host} as your first officer`; Problem 2 (status overload/self-serve) → line 3 `ask it for the queue and next steps.`, both `spacedock status` removed; Problem 3 (multi limbo) → `Workflows: N found — your first officer will help you pick`; minors: `Workflow:`/`Workflows:` agreement + `Version` dev sentinel. AC-1..AC-7 each cite a rendered-output test (AC-5 byte-exact golden, AC-6 property+regex); discovered the pi no-banner gap (AC-7).

### Summary

Mapped today's front-door output for every MODE × HOST cell from the real code (claude/codex via `launchBanner` after the gate; pi has no banner at all — a found gap) and captured the BEFORE column by exercising `launchBanner` rather than reading it. Proposed a 3-line AFTER banner that fixes all three filed problems (metaphor flip, `spacedock status` overload/self-serve contradiction, multi-workflow pick-then-launch limbo) plus the two minors (singular/plural label, stale unstamped `Version`). The gate/install/refuse messages are scoped OUT (not the banner). Seven ACs, each pinned by an `internal/cli` test over rendered output with the proposed wording as the independent expectation; no spike needed (text reorder + an existing call site added to `runPi`, all over already-proven mechanisms).

### Captain wording fold

Captain wording call — "keep it simple." Line 1 and line 3 unchanged from the prior design. Line 2 made terser — just the facts, no count, no helper sentence, no command:

- single: `Workflow: {path}` (unchanged, e.g. `Workflow: docs/dev`).
- multi: `Workflows: {space-separated workflow paths}` (e.g. `Workflows: docs/dev docs/user-testing`) — LIST the detected paths, space-separated. The `{N} found — your first officer will help you pick` wording is dropped entirely. Still resolves Problem 3: no actionable-but-ignored instruction; the user sees what is there and the FO picks in-session.
- none: `Workflow: none detected` — the `(launching anyway)` parenthetical dropped.

`detectedWorkflow`'s multi branch now returns the space-joined paths (mapping `workflowLabel` over the discovered list it already has) and its none branch returns `none detected`; the `Workflows:`/`Workflow:` label agreement (plural for multi, singular otherwise) stays. Folded into: the AFTER matrix rows, Problem 3 prose, the minors + wording notes, AC-3 (multi expected → space-joined path list; assert NO count / NO `will help you pick` / NO `spacedock status`), AC-4 (label agreement holds), AC-5 (single golden unchanged), and the risk/regression notes (the multi branch is now a list-render, not a reword). The wording notes record the captain's terse path-list choice.

## Stage Report: implementation

- DONE: Rewrite `launchBanner` (internal/cli/frontdoor.go) to the captain-APPROVED AFTER wording: line 1 `launching {host} as your first officer`; line 3 `{host} is your first officer — ask it for the queue and next steps.`; multi branch → space-joined discovered workflow paths; none branch → `none detected`; `Workflow:`/`Workflows:` label agreement — AC-1/2/3/4.
  `detectedWorkflow` now returns a `(label, value)` pair; multi maps `workflowLabel` over the sorted `DiscoverWorkflows` list. Worktree commit `f953498a`. Probed rendered output for all three shapes — byte-exact to the matrix.
- DONE: Add the `launchBanner` call site to `runPi` so pi emits the same 3-line banner before launch, and default `Version` in `cli.go` to the `dev` sentinel — AC-7/AC-6.
  `runPi` calls `launchBanner("pi", dir, stderr)` before building argv; pi has no recognized resume seam today so the banner is unconditional (the "if one exists" gate has nothing to attach to). `Version = "dev"`.
- DONE: Cover every AC with an `internal/cli` test over RENDERED banner output (AC-5 byte-exact golden; AC-1/2/3/4/7 presence/absence; AC-6 property + `^\d+\.\d+\.\d+$` non-match) AND update the existing `launch_banner_test.go` expected strings to the AFTER wording. `go test ./internal/cli/` green.
  New `launch_banner_wording_test.go` (AC-1..AC-7). Adversarially verified: reverting line 1 fails AC-1+AC-5; reverting the multi branch fails AC-3; reverting `Version` to `0.19.0` fails AC-6. `go test ./internal/cli/` green (full package, ~11s); `go vet`/`gofmt` clean.

### Summary

Rewrote the pre-launch banner to the captain-approved wording across all three workflow shapes and all three hosts, fixing the metaphor flip (line 1), the self-serve `spacedock status` overload (line 3), and the multi-workflow pick-then-launch limbo (line 2 now lists the detected paths space-joined). Closed the pi no-orientation gap by adding the banner to `runPi`, and made the unstamped `Version` read as `dev` so a `go build`/`go install` binary cannot impersonate a release. Every AC is pinned by a rendered-output test (AC-5 a byte-exact golden, AC-6 a property+regex), each independently verified to FAIL against the BEFORE wording. One note for the validator: a pre-existing, unrelated failure in `internal/dispatch` (`TestLauncherCommandFallsBackToPathWhenSpacedockBinUnsetEmptyOrUnusable/unset`) reproduces on the clean base with my changes stashed — it is the dispatch launcher-token test resolving a real `spacedock` on `$PATH`, disjoint from this banner work.

## Stage Report: validation

- DONE: Reproduce each AC-1..AC-7's cited evidence independently over RENDERED banner output: `go test ./internal/cli/` GREEN; AC-5 byte-exact golden on the single-workflow path holds; AC-1/2/3/4/7 presence/absence over the rendered bytes; AC-6 asserts the package-default `Version` == the dev sentinel AND does not match `^\d+\.\d+\.\d+$`.
  `go test ./internal/cli/` GREEN (272 tests). Drove the real `launchBanner`/`runPi` via a throwaway probe and dumped the exact bytes for all four cells — single/claude+codex, multi/claude, none/codex — each matches the captain-approved AFTER matrix character-for-character (line 1 `launching {host} as your first officer`, never `first officer launching`; line 2 single `Workflow: docs/dev` / multi `Workflows: docs/dev docs/user-testing` / none `Workflow: none detected`; line 3 `{host} is your first officer — ask it for the queue and next steps.`, no `spacedock status`). AC-5 golden GREEN, AC-6 `Version == "dev"` and regex non-match GREEN, AC-7 pi banner GREEN. Probe removed; impl worktree clean. Every AC's expected value is the proposed wording (independent source), not a source-grep of the constant.
- DONE: Also run `go test ./...` from root and account for every failure — confirm the disclosed `internal/dispatch` launcher failure is PRE-EXISTING on the clean baseline (env-dependent `spacedock`-on-PATH) and not introduced here.
  Only `internal/dispatch` `TestLauncherCommandFallsBackToPathWhenSpacedockBinUnsetEmptyOrUnusable/unset` fails. The impl commit `f953498a` touches ZERO files under `internal/dispatch`. Reproduced the identical `/unset` failure on a detached checkout of the merge-base `432cda52` (zero banner changes) — pre-existing, env-dependent (a real `/opt/homebrew/bin/spacedock` exits 2 on bare invocation). Not introduced here; not grounds for rejection. `gofmt -l` clean, `go vet ./internal/cli/` clean.
- DONE: Detached adversarial audit (MANDATORY — front-door is high-stakes): on a THROWAWAY checkout of the merge result, try to REFUTE the wording guards and confirm the deliverable's own tests catch EACH.
  Ran on a detached worktree at the branch HEAD (= merge result; branch is 4 ahead / 0 behind main, fast-forwards clean), never the impl worktree. All five adversarial reverts caught: (1) line 1 → `first officer launching` ⇒ AC-1+AC-5+AC-7 RED; (2) line 3 → self-serve `spacedock status` form ⇒ AC-2+AC-3+AC-5+AC-7 RED; (3) multi branch → `{N} workflows detected (run \`spacedock status\` to pick)` ⇒ AC-2+AC-3+AC-4 + existing `TestLaunchBannerNamesDetectedWorkflow` RED; (4) `Version` → `0.19.0` ⇒ AC-6 RED; (5) drop `runPi` banner call ⇒ AC-7 RED (compiles — `dir` still used at `piRuntimeConfigFromEnv`, so a real assertion failure, not a compile error). Each edit reverted between runs; checkout clean + GREEN after all reverts; both throwaway worktrees torn down.
- DONE: File a PASSED/REJECTED recommendation citing reproduced evidence per AC; record the audit outcome; route any Material finding back through feedback.
  Recommendation: PASSED. Detached audit: no material findings. No `### Feedback Cycles` entry needed.

### Summary

PASSED. Reproduced every AC independently over the RENDERED banner bytes (not a source-grep): the rendered output for all four mode×host cells matches the captain-approved AFTER matrix exactly — AC-5's byte-exact single-workflow golden holds, AC-6's `Version == "dev"` + `^\d+\.\d+\.\d+$` non-match holds, AC-7's pi banner reaches stderr and the launch seam. The MANDATORY detached adversarial audit ran on a throwaway merge-result checkout: all five reverted-wording edits (line 1 framing, line 3 self-serve status, multi count string, release-shaped `Version`, dropped pi banner call) are caught by the deliverable's own tests — no green-stays-green hole, refuted nothing material. The lone full-suite failure (`internal/dispatch/…/unset`) is confirmed pre-existing on the clean merge-base `432cda52` (zero banner changes), env-dependent on a real `spacedock` on `$PATH`, and disjoint from this work — accounted for, not a rejection ground. `gofmt`/`go vet` clean.
