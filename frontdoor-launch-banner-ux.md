---
id: cmxchb8y1y2m455xhx7ce87g
title: Launch-banner UX — first-officer framing, status-command overload, multi-workflow limbo
status: ideation
source: "FO + captain launch-banner review (2026-06-08), following yq (frontdoor-launch-ux). The 3-line pre-launch banner (frontdoor.go:139-141) reads confusingly on the new-user / multi-workflow path."
started: 2026-06-08T23:37:46Z
completed:
verdict:
score:
worktree:
issue:
sprint:
group:
sprint-readiness:
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
| compatible, **single** workflow | claude | `spacedock 0.19.0 · first officer launching claude`<br>`Workflow: docs/dev`<br>`claude is starting as your first officer; run \`spacedock status\` inside the session for the queue.` | `spacedock {v} · launching claude as your first officer`<br>`Workflow: docs/dev`<br>`claude is your first officer — ask it for the queue and what to do next.` |
| compatible, **single** workflow | codex | same 3 lines, `codex` | `spacedock {v} · launching codex as your first officer`<br>`Workflow: docs/dev`<br>`codex is your first officer — ask it for the queue and what to do next.` |
| compatible, **multi** workflow | claude | `spacedock 0.19.0 · first officer launching claude`<br>`Workflow: 2 workflows detected (run \`spacedock status\` to pick)`<br>`claude is starting as your first officer; run \`spacedock status\` inside the session for the queue.` | `spacedock {v} · launching claude as your first officer`<br>`Workflows: 2 detected — your first officer will help you pick`<br>`claude is your first officer — ask it for the queue and what to do next.` |
| compatible, **none** detected | codex | `spacedock 0.19.0 · first officer launching codex`<br>`Workflow: none detected (launching anyway)`<br>`codex is starting as your first officer; run \`spacedock status\` inside the session for the queue.` | `spacedock {v} · launching codex as your first officer`<br>`Workflow: none detected (launching anyway)`<br>`codex is your first officer — ask it for the queue and what to do next.` |
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
2. **`status` overloaded + self-serve contradiction.** BEFORE prints `spacedock status` twice and tells the user to run it themselves — contradicting "the host is your first officer (it runs status for you)." AFTER removes BOTH `spacedock status` instructions. Line 3 becomes `ask it for the queue and what to do next.` The pitch holds: you ask your first officer; it runs status.
3. **Multi-workflow "pick then launches anyway" limbo.** BEFORE announces "N workflows detected (run `spacedock status` to pick)" then launches with no pick. AFTER states no actionable command; it tells the truth about what happens next: `your first officer will help you pick` (the FO disambiguates in-session on its first turn, per the FO contract).

Minors:
- **Singular label / plural content.** BEFORE always says `Workflow:` even for N workflows. AFTER uses `Workflow:` for the single + none cases and `Workflows:` for the multi case, so label and content agree.
- **Stale Version display.** `Version` (`cli.go:28`) defaults to the literal `"0.19.0"` and is only overwritten by the goreleaser ldflags stamp (`.goreleaser.yaml:37`, `git describe`). A dev/`go install` build (no ldflags) prints `0.19.0` regardless of the actual commit — impersonating a stale release. Fix: default `Version` to a non-release sentinel (`"dev"`) so an unstamped build reads honestly as a dev build. A stamped release overwrites it as before. This is a real, checkable change (AC-6).

### Wording notes

- Line 1 separator stays ` · ` (matches today). `{v}` is the `Version` value (a real release on a stamped build, `dev` on an unstamped one).
- Line 2 single/none keep the existing `detectedWorkflow` text verbatim (`docs/dev`, `none detected (launching anyway)`) — those are not problems; only the multi case changes.
- Line 2 multi: `Workflows: {N} detected — your first officer will help you pick`. No backticked command, so nothing to mis-run.
- comm-officer polish was requested on the wording (best-effort); no reply within the 2-minute budget, so the wording above is the ensign's. It is terse, each line stands alone, and it satisfies the three constraints (consistent metaphor, no self-serve `status`, clean single-workflow path).

## Acceptance criteria

Each AC is verified by an `internal/cli` test over the **rendered** banner output (the bytes `launchBanner` / `runPi` writes), never a source-grep. The independent expected value is the proposed wording above (a renamed-but-equivalent branch must still fail the assertion).

- **AC-1 — consistent metaphor, no launcher-as-FO.** The rendered banner never contains the phrase `first officer launching` (the flipped line-1 framing); line 1 reads `launching {host} as your first officer`.
  *Verified by:* a test asserting the banner output (claude and codex) contains `launching {host} as your first officer` and does NOT contain `first officer launching`.
- **AC-2 — no self-serve `spacedock status` in the banner.** The rendered banner contains no `spacedock status` instruction on any of the single / multi / none paths; line 3 reads `{host} is your first officer — ask it for the queue and what to do next.`
  *Verified by:* a test asserting the banner output for all three workflow cases does NOT contain `` `spacedock status` `` and DOES contain `ask it for the queue`.
- **AC-3 — multi-workflow states no pick instruction.** On the multi-workflow path, the rendered banner reads `Workflows: {N} detected — your first officer will help you pick` and contains neither `run \`spacedock status\` to pick` nor any bare-command pick instruction.
  *Verified by:* a test that builds a 2-workflow git-repo fixture, renders the banner, and asserts it contains `2 detected — your first officer will help you pick` and NOT `to pick)` / `spacedock status`.
- **AC-4 — singular/plural label agreement.** The single-workflow and none paths render `Workflow:`; the multi path renders `Workflows:`.
  *Verified by:* a test asserting `Workflow: docs/dev` for the single fixture, `Workflow: none detected` for the bare fixture, and `Workflows: 2 detected` for the 2-workflow fixture.
- **AC-5 — single-workflow happy path stays clean and terse.** The single-workflow banner is exactly the three proposed lines, in order, with no extra lines.
  *Verified by:* a test asserting the full rendered single-workflow output equals the three expected lines joined by `\n` (a byte-exact golden assertion, the strongest form, so any drift fails).
- **AC-6 — unstamped build does not impersonate a release.** An unstamped (`go build` / `go install`) binary's `Version` is not a release-shaped semver; it reads as a dev build.
  *Verified by:* a test asserting the package default `Version` (the value before any ldflags stamp) equals the dev sentinel and does not match a `^\d+\.\d+\.\d+$` release pattern. (A stamped release still overwrites it; the ldflags path is exercised by the existing release pipeline.)
- **AC-7 — pi emits the banner.** `spacedock pi` renders the same 3-line banner before launch (closing the pi no-orientation gap), suppressed on the pi resume path if one exists.
  *Verified by:* a test driving `runPi` with a ready stub runtime and asserting the rendered stderr contains `launching pi as your first officer` and `pi is your first officer — ask it for the queue`, and that the launch seam is still reached.

### Risk / spike

No spike needed. The design only reorders/rewords text already produced by `launchBanner` and adds an existing call site to `runPi`; the workflow-detection mechanism (`detectedWorkflow`, repo discovery), the gate, and the launch seam are all already proven by the existing `launch_banner_test.go` / `frontdoor_test.go` / `launch_parity_test.go` suites. The proposed-output BEFORE column was captured by exercising the real `launchBanner` (single / multi / none) and the pi `--no-install` path was probed live, both via throwaway tests — recorded above, not inferred from source. The `Version` sentinel (AC-6) composes the proven goreleaser stamp path (the ldflags `-X` already works for releases); only the default literal changes.

### Test plan

- **Fixtures, fast.** All ACs are `internal/cli` Go tests using the existing `commissionWorkflowAt` / `gitRepoFixture` helpers in `launch_banner_test.go` (single / multi / none repos) plus the `fakeHost` (claude/codex launch) and a pi stub runtime (AC-7). No live workflow, no network.
- **Cost:** low — extends the existing `launch_banner_test.go` patterns; sub-second.
- **Golden vs substring:** AC-5 is a byte-exact golden on the single-workflow happy path (the strongest assertion). AC-1/2/3/4/7 are presence/absence assertions over rendered output where a partial match is the right granularity. AC-6 is a property assertion over the package default plus a regex non-match. None is a source-grep; every expected value is the independent proposed wording, so a renamed branch or a valid paraphrase fails.
- **Regression:** the existing `TestLaunchBannerNamesDetectedWorkflow` / `…ReachesStderrBeforeLaunch` / `…SuppressedOnResume` will need their expected strings updated to the AFTER wording (they currently assert the BEFORE lines, e.g. `2 workflows detected (run \`spacedock status\` to pick)`); that update IS part of the implementation, not extra scope.

## Stage Report: ideation

- DONE: Produce the PROPOSED launch-output matrix the captain asked for (redesigned front-door output for each MODE × HOST, before/after, addressing the three filed banner problems; ideation/design, no implementation).
  BEFORE/AFTER table under `## Proposed design` covers compatible single/multi/none, fresh-install, upgrade-skew, --no-install, resume × {claude, codex, pi}; deliverable is matrix + ACs, no code changed.
- DONE: Map the CURRENT front-door output for each cell first (read the real code, don't guess).
  Read frontdoor.go (`launchBanner`:138, gate arms :259-283/:418-441, `noPluginRemedy`:127), init.go, host_exec.go, pi.go (`runPi`:65 — no banner), contract.go (`mismatchMessage`:155); BEFORE column captured by exercising real `launchBanner` (single/multi/none) + a live `pi --no-install` probe (both throwaway tests, removed).
- DONE: Propose the REVISED output per cell fixing the three filed problems + minors; present as a BEFORE/AFTER table (mode × host); give ACs each verified by an `internal/cli` test over rendered output, never a source-grep.
  Problem 1 (metaphor flip) → line 1 `launching {host} as your first officer`; Problem 2 (status overload/self-serve) → line 3 `ask it for the queue and what to do next.`, both `spacedock status` removed; Problem 3 (multi limbo) → `Workflows: N detected — your first officer will help you pick`; minors: `Workflow:`/`Workflows:` agreement + `Version` dev sentinel. AC-1..AC-7 each cite a rendered-output test (AC-5 byte-exact golden, AC-6 property+regex); discovered the pi no-banner gap (AC-7).

### Summary

Mapped today's front-door output for every MODE × HOST cell from the real code (claude/codex via `launchBanner` after the gate; pi has no banner at all — a found gap) and captured the BEFORE column by exercising `launchBanner` rather than reading it. Proposed a 3-line AFTER banner that fixes all three filed problems (metaphor flip, `spacedock status` overload/self-serve contradiction, multi-workflow pick-then-launch limbo) plus the two minors (singular/plural label, stale unstamped `Version`). The gate/install/refuse messages are scoped OUT (not the banner). Seven ACs, each pinned by an `internal/cli` test over rendered output with the proposed wording as the independent expectation; no spike needed (text reorder + an existing call site added to `runPi`, all over already-proven mechanisms).
