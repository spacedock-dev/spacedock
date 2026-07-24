---
title: Gate presentation as an overridable channel with atomic result retention
status: implementation
source: "Split from the gate-recorder task (3k), captain-approved 2026-07-21. The subspace-coupled presentation half; 3k cycles 11-12 are its banked design history."
id: xbatj4hxtxw9t83vvmfem27f
sprint: durable-decisions
group: recorder
started: 2026-07-21T01:43:36Z
worktree: .worktrees/spacedock-ensign-gate-review-presentation-command
gates:
    version: 1
    current:
        gate: gate:docs-dev:xb:validation
    records:
        - id: gate:docs-dev:xb:validation
          stage: validation
          attempts:
            - id: gate-attempt:xb-validation-1
              briefing:
                id: briefing:docs-dev:xb:validation:attempt-1:revision-1
                digest: sha256:772a856dcd3dd7d5a1bcfb589854b4b7f5b70bb26393a7e1e90aa2605daf0911
                digest-domain: canonical-bytes
                room-ref: ./review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:xb:validation:1
                briefing: briefing:docs-dev:xb:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-23T23:49:14.728783Z"
                decision: approve
                reason: Spacedock 612b72fc and provider 198f7623 have 6/6 ACs evidenced, retained-delivery and association suites green, zero binary coupling, and no material finding.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: done
                state: pending
                blockers: []
---

Gate presentation is an **overridable channel of the present-gate skill**, not a spacedock binary verb. The present-gate skill is the existing presentation authority (the captain-facing gate template + assembly rules); its DEFAULT channel is chat, already its only behavior today. A workflow or session may declare an OVERRIDE channel — a hardened Subspace float (the `review-local-zellij` lineage with the retention contract: caller-owned result path, room-resident from the first byte, probe-first ritual) — presenting the design as a briefing package in a blocking float and returning a retained result. Spacedock's binary carries **zero** Subspace knowledge; the only Subspace-coupled artifact is the opt-in override script, so the cross-repo release-train coupling (this design's biggest declared liability) dissolves to that script. The duties that genuinely need binary support — result validation, provider id-normalization, and the record handoff — are **recorder-side** (they validate what gets recorded, whatever the channel) and are proposed as a small extension of the recorder's existing parse-and-verify step via the dogfooding change protocol; the id-mapping rule specified in 3k's `gate-resolution-frontmatter-contract.md` moves its implementation home from the presentation side to the recorder (owner-tag amendment proposed below). Atomic retention lives in the override script under a testable contract whose drive suite is this task's 12 red-fixture spike assertions, run at validation. Evidence base: the 0260 shaping float findings 1-15 (blank-float EOF defect, launcher repair, probe-first ritual, retention-deleted-on-failure incidents, and finding 15 — the captain's own approval destroyed when a dead launcher unlinked its scratch result; finding 14 is a deliberate live-session numbering skip) in the shaping debrief and 3k's attempt-7 resolution provenance note. Kept from the prior attempts: the no-subspace fallback (now simply the default channel), the probe ritual, and every red fixture. Land after the recorder (3k).

## Expected surface + tolerance

Re-estimated after the attempt-2 reframe: the spacedock-side Go surface for presentation is **~zero**. The new surface is three parts, no new `spacedock` verb:

- **Present-gate skill prose** (spacedock plugin, `skills/present-gate/SKILL.md`): a new "presentation channel" section — DEFAULT chat (its current behavior, unchanged) and an OVERRIDE hook naming the channel contract (how the FO invokes an override and where the retained result lands). ~25-40 lines of skill prose.
- **One hardened override script** with a committed drive suite — the `review-local-zellij` successor implementing the retention contract (caller-owned result path, room-resident, no `trap rm -rf`, probe-first). ~80-140 script LOC + the 12 red-fixture spike assertions promoted to its committed test suite. **Home:** the subspace repo (the conforming provider ships and CI-tests its own override), so spacedock keeps zero Subspace code — see the honest-assessment section for why the home is load-bearing.
- **A recorder-surface ask** (routed to 3k via the change protocol, NOT built here): the recorder's existing parse-and-verify step also normalizes the provider envelope briefing id and validates the result it records, for any channel. Likely single-digit LOC on 3k's surface; xb proposes, 3k owns.

**Tolerance 2×**, measured against the skill-prose + script + drive-suite surface above (not a Go product surface). Hard self-check: any `spacedock` binary change that imports or shells out to Subspace (the coupling the reframe removes), any entity-frontmatter write (the recorder's, 3k), any blocker/eligibility/hold/application computation (h1's), or any subspace-tui product change (xb consumes its surface, does not build it) trips a reconfirm.

## Problem

A captain's gate decision presented through Subspace can be lost. The 0260 shaping session ran five successful Subspace floats but collected the presentation-side failure evidence in the same two days (shaping debrief float findings 1-15; finding 14 is a deliberate live-session numbering skip):

- **A dead launcher destroyed the captain's own approval (the primary red fixture).** Finding 15, recorded in 3k's attempt-7 resolution provenance note: the captain resolved a gate in a float pane whose launcher had already died; the resolution was written to a launcher-owned scratch file that was then unlinked, destroying the decision. Only a chat re-affirmation recovered it. A retention target owned by a dying launcher is exactly the failure the retention contract removes — the result lands at a caller-owned result path outside any launcher-owned scratch, so the launcher's death cannot unlink it. The interim launcher's later room-resident-result fix is the same lesson reached by hand.
- **The result is destroyed on the failure and hold paths.** The interim `review-local-zellij` launcher captures the child's stdout into a `mktemp -d` scratch and carries `trap 'rm -rf "$scratch"' EXIT` (findings 3, 12). When result validation fails, or the reviewer leaves the review open/hold, the scratch — and the only copy of the result — is deleted. The leave-open/hold result shape is exactly the branch that was never retained.
- **The blank-float EOF leaves no recoverable state.** The `subspace-tui present` child protocol opened an empty float and ended with `present-child protocol ended early: EOF`, exit 1, no recovery path emitted (finding 12). Reproduced here against the exact-tip binary: on the can't-render path the provider writes **nothing** to its `--result` path and removes the zero-entry package, so the only trace is one line of stderr — which the interim launcher then throws away.
- **Single-file concatenation.** The interim launcher accepts exactly one Markdown file, forcing the gate summary and the frozen entity snapshot to be concatenated; the multi-artifact briefing package is designed (3k) but the interim path cannot present it as separate artifacts (finding 5, captain-flagged).

The decision therefore does not reliably become durable, recoverable state, and the ritual that works today depends on a skill-local shell script with a destroy-on-exit defect. This is the presentation half split from the gate recorder (3k); the recorder owns the `gates:` writes, and this task owns getting a validated, retained result into the First Officer's hands without losing it — through whatever channel the FO presents.

The prior attempts proposed a `spacedock gate review` Go verb to carry that. The captain's attempt-2 revise questioned that vehicle: baking Subspace presentation into spacedock's binary couples spacedock's release train to subspace-tui's surface — the biggest liability this design itself declared. The reframe below removes that coupling: presentation becomes an overridable channel of the present-gate skill, the binary stays Subspace-free, and the guarantees relocate to where they belong (retention in the override script, validation/id-mapping in the recorder).

## The presentation channel (present-gate overridable)

The present-gate skill selects a channel at the gate point. Neither channel is a mode or a degraded state; the gate proceeds identically and only the rendering differs.

- **DEFAULT — chat.** The First Officer renders the gate-review template in conversation (the skill's existing behavior, unchanged). The captain's chat-rendered resolution is recorded through the recorder (3k) exactly as every 0260 gate was. This is also the no-subspace fallback (below): when no override is available, chat is simply what runs.
- **OVERRIDE — hardened Subspace float.** A workflow or session declares the override. The FO invokes the override script, which: probes the presenter first (probe-first ritual); assembles the multi-artifact briefing package (gate summary, frozen design snapshot, frozen probe input/history — separate artifacts, not concatenated); runs the presenter as one **blocking** call through the caller's terminal transport, with the result written **room-resident from the first byte** to a caller-owned path; and returns a retained result. The FO keeps the gate-attempt ensign addressable for the whole blocking call and waits with `wait_agent({timeout_ms:300000})`, re-waiting on timeout — pane creation and timeout are never completion.

**Channel selection and the fallback (kept from attempt 1).** Selecting the override probes the presenter (binary via `SUBSPACE_TUI_BIN`/`PATH`, version-gated) **before any side effect**. Absent or version-mismatched → the selection falls back to the default chat channel, emitting one line that names both the install remedy and the chosen fallback, with **zero side effects and no retention directory created**. A missing presenter is an ordinary condition selecting the chat channel — not a mode, never blocking the gate — mirroring the teams-unavailable-selects-bare-dispatch ruling shape (`dispatch-failure-retry-rung.md`: an unavailable capability "is not a mode at all — it is the ordinary condition selecting [the alternate path], evaluated where dispatch happens"). The recording-identity ruling (`docs/roadmap/durable-decisions/index.md` Constraints, 2026-07-21) applies to both channels unchanged: a resolution's provenance names its channel, and the captain's identity attaches to content the captain saw. No Subspace-specific field is required to record a chat resolution.

## Splitting the duties honestly

The reframe asks which guarantees genuinely need binary support. They do not need a *spacedock* binary; each has an honest home:

| Guarantee | Home in the reframe | Why not the spacedock binary |
|---|---|---|
| Atomic retention on success AND failure (caller-owned result path, room-resident, no trap-delete, survives launcher death) | the **override script**, under a testable contract; drive suite = this task's 12 spike fixtures, run at validation | the spike proved the retention contract works in a shell script — the interim script failed on `trap rm -rf`, not because a script can't retain. No compiled code is required. |
| Result validation (result present, parseable, digest-bound to the attempt briefing, decision valid) | the **recorder** (3k) — already its contract: "parses the exact provider result, verifies authorized identity, current Briefing id/digest and log rules" | validating what gets recorded is channel-independent and already a recorder duty; the binary would duplicate it. |
| Provider id-normalization (envelope id → attempt briefing id, after digest match) | the **recorder** — a small extension of the same parse-and-verify step | the recorder is the one place that already verifies the digest; normalizing there removes the need for a presentation-side implementation. |
| Record handoff | the **recorder** verb the FO calls | unchanged — the FO hands the retained result to the recorder; no presentation binary needed. |
| Detection + version-gate + fallback | the present-gate skill's **channel selection** | prose + a probe; the skill already owns channel choice. |
| Blocking / addressable presenter / pane-not-completion | override script (blocking float) + FO followup/wait discipline (already skill prose) | a blocking script call plus existing FO discipline; no binary. |
| Never mutate frontmatter | trivially true | neither the skill nor the script writes frontmatter; the recorder owns `gates:`. |

**Owner-tag amendment (proposed, not applied here).** The contract's id-mapping rule was tagged "specified in `gate-resolution-frontmatter-contract.md` (3k), implemented by the presentation side (xb)." The reframe moves the implementation home to the recorder: **specified AND implemented recorder-side (3k)**. xb proposes this to 3k through the dogfooding change protocol (`docs/roadmap/durable-decisions/index.md`: frictions in another owner's section route through the FO to the owner; consumers re-anchor on landed text). 3k owns whether to accept it; nothing here edits 3k's surface.

## Honest assessment: is the decoupling sound?

Yes — and it is better than the binary. The captain's question is well-founded: presentation is inherently a provider concern, and compiling it into spacedock's release train is the coupling this design flagged as its biggest liability. Every guarantee relocates without loss (table above): retention is a script contract the spike already demonstrated in bash; validation and id-mapping were always recorder duties (3k's contract states result validation explicitly); detection/fallback is channel selection the skill already owns. The spacedock binary ends up Subspace-free, which is a checkable property (AC-6), and the release-train coupling dissolves to an opt-in script.

**The one load-bearing condition — not a defeater, a constraint.** The binary's real, non-decorative value was bringing the retention contract under committed `go test` with the 12 fixtures as first-class tests. `review-local-zellij` shipped the destroy-on-failure defect precisely because it was an untested "working copy … must not be published as-is" script. If the override lands as another untested skill-local script, the reframe reintroduces exactly the class of defect this task exists to kill. So the decoupling is sound **only if** the override script carries a committed, CI-run drive suite (the 12 red fixtures). That condition also fixes the script's home: the subspace repo, which already runs `go test` over `reviewonefile` and can host a shell/bats or Go-harness drive suite, so spacedock keeps zero Subspace code while the retention contract stays tested. This is the single point I put in front of the captain; it constrains the reframe, it does not defeat it. I found no piece that genuinely cannot be a skill override without losing a guarantee.

## Provider id-mapping (specified in the contract; implementation home moves to the recorder)

The provider mints its own briefing id: subspace-tui's one-file mode derives `briefing:single-file:<hex16>` from the result path (`BriefingIDForInvocation`), and stamps it into `result.briefing` and `result.resolution.briefing`. The recorder verifies the current attempt Briefing id/digest, so a provider-minted id would never bind. Per `gate-resolution-frontmatter-contract.md` (the frozen Briefing binds an immutable snapshot; the provider result is keyed by Briefing id and joined after digest validation), the recorder's parse-and-verify step:

1. validates the artifact digest in the returned result against the digest bound by the attempt briefing;
2. **only on digest match**, rewrites the provider envelope briefing id in the recorded result to the attempt briefing id;
3. on digest mismatch, rejects without normalizing (an unverified result is never laundered into the attempt id).

Under the reframe this rides the recorder's existing result-verification, not a separate presentation command (owner-tag amendment above). It concerns only what the recorder records; it never touches frontmatter outside the `gates:` write the recorder already owns.

## Spike: the riskiest unverified mechanism (atomic retention on success AND failure)

Exercised end-to-end before design lock; harness at `scratchpad/retention-spike.sh`, plus a probe of the exact-tip binary.

- **Real-binary grounding.** Built the tip `subspace-tui` (`go build ./cmd/subspace-tui`, version `dev`) and ran the advisory one-file mode with `--result` on the can't-render path (no TTY, the blank-float EOF analogue). Result: **exit 1, nothing written to `--result`, zero-entry package removed, one line of stderr.** Confirms the command cannot rely on the provider to retain anything on the failure path — it must own retention.
- **Retention-contract spike (12/12 assertions pass).** A fake TUI child reproduces the three red fixtures — (A) blank-float EOF: exit non-zero, empty result; (B) leave-open/hold: writes a `"status":"open"` result, exit 0; (C) launcher/controller death after the result is written (finding 15): a real launcher process writes the approval then is `kill -TERM`ed before returning — under two wrappers: the interim `trap rm -rf` scratch pattern (launcher owns the result path), and the command's contract (caller-owned dir never deleted, `--result` binding, stderr→diagnostics, argv+exit recorded). The interim baseline retains **nothing** on all three (reproducing finding 3, finding 12, and finding 15 — the dying launcher's trap unlinks the captain's approval); the command's contract retains the result-or-diagnostics, the log when present, and the launch record on **all three**, and the approval written to the caller-owned `--result` path survives the launcher's death because the launcher owns no cleanup over it.
- **No-subspace fallback fixture (D).** Detection runs before any side effect: an absent presenter exits non-zero, emits one line naming both the install remedy and the chat fallback, and creates no retention directory — the checkable form of AC-5's fallback clause.

Conclusion: the retention mechanism is proven at the contract level with real files, real exit codes, and a real killed launcher process, reproducing the exact debrief failure signatures — including finding 15, 3k's named primary red fixture. Under the reframe these 12 assertions become the **override script's committed drive suite**, run at validation — the condition that keeps the script from repeating `review-local-zellij`'s untested-script defect. That a shell script passes the whole retention contract is itself the evidence that retention needs no spacedock binary.

## Acceptance criteria

**AC-1 (VALUE) — No presented decision is lost on any exit path.** Driving the override script through each retention fixture — resolved approve, resolved revise, leave-open/hold, blank-float EOF (child exits non-zero, empty result), mid-run child crash, and post-exit result-validation failure — leaves a caller-owned retention directory that still holds the result (or, when the child produced none, a diagnostics record naming the launch argv, exit code, and captured stderr) plus the review log when one exists. Retained count is N/N across the fixtures; the same fixtures through the interim `review-local-zellij` baseline retain 0/N on the failure and hold paths. *Test:* the override script's committed drive suite runs a fake TUI child per fixture (the spike harness promoted to committed tests in the script's home repo), asserting on-disk retention-directory contents after each run and contrasting the baseline's `trap rm -rf` destruction. This measures the end-value (result-retention rate) against an independent baseline that can regress to 0.

**AC-2 — Retention survives every failure class, including launcher/controller death.** Launch, controller, child-exit, validation, and retention-write failures each leave the briefing package, diagnostics, and any partial result recoverable; none deletes the retention directory, and each returns non-zero. In particular, a launcher/controller that dies after the result is written (finding 15 — the red fixture where a dying launcher unlinked the captain's approval) leaves the result intact, because it lands at the caller-owned result path the override script allocates, not in a launcher-owned scratch the launcher's cleanup can unlink. *Test:* the drive suite injects each failure class, including a launcher process that writes the result then is killed before returning; asserts directory survival, the result surviving the launcher's death, diagnostics presence, and exit code.

**AC-3 — Pane/session creation and wait-timeout are never completion.** The override run returns success only after the child exits, the result validates, and retention completes; a launched-but-unexited child yields neither a success exit nor a retained validated result. The FO keeps the gate-attempt ensign addressable for the whole blocking call and re-waits on `wait_agent` timeout. A mutant that returns on pane creation, and a mutant that lets the presenter resolve before child exit, both fail. *Test:* drive-suite fixture with a child that emits a pane marker then blocks; assert no success and no validated result until the child exits.

**AC-4 — The recorded result is keyed to the attempt briefing id, only after digest validation (recorder-homed, proposed).** Given a result whose `briefing`/`resolution.briefing` carry the provider-minted `briefing:single-file:<hex>` id: on a digest matching the attempt briefing, the recorded result's briefing ids equal the attempt briefing id; on a mismatch, it is rejected with the ids un-rewritten. Under the reframe this is the recorder's parse-and-verify duty (owner-tag amendment proposed to 3k); if 3k declines the move it stays presentation-side. *Test:* a recorder fixture feeding a provider result with the minted id plus a matching and a mismatching artifact digest; assert normalized ids on match and rejection (no normalization) on mismatch. Lands on whichever surface owns id-mapping after the change-protocol proposal resolves.

**AC-5 — The override channel validates the briefing and derives the title before any launch; an absent or version-mismatched presenter falls back to chat with zero side effects.** A valid multi-artifact briefing yields a canonical pane title in the provider's grammar and confirms each Artifact and supporting Reference revision is present. An absent `subspace-tui`/`SUBSPACE_TUI_BIN` or a version mismatch makes channel selection fall back to chat, emitting one line naming **both** the install remedy **and** the chat fallback, launching nothing and creating **no retention directory** — the entity and working tree byte-unchanged. *Test:* a valid-briefing fixture (assert title derivation + Reference presence); a missing-binary and a version-mismatch fixture (assert fallback-to-chat, a message naming both remedy and fallback, zero launch, no retention directory).

**AC-6 (VALUE) — Presentation adds zero Subspace coupling to the spacedock binary, and no channel mutates entity frontmatter.** After this task, the `spacedock` binary depends on no Subspace package and exposes no presentation verb (no `spacedock gate review`); presentation lives entirely in the present-gate skill (prose + channel selection) and the override script. Neither channel writes entity frontmatter. This is the reframe's headline value — the release-train coupling this design flagged as its biggest liability measures to zero. *Test:* a dependency/surface check on the built binary (`go list -deps ./cmd/spacedock` names no subspace package; the CLI verb list contains no `gate review`) — a build-artifact/dependency assertion, not a prose-grep; plus a byte-compare of the entity file across every override fixture. Baseline that can regress the wrong way: any reintroduced import or verb makes the count non-zero.

## New mechanisms (value AC served / simplest alternative / why insufficient)

- **Presentation as an overridable present-gate channel, not a spacedock verb** — serves AC-6 (VALUE). *Alternative:* a `spacedock gate review` Go verb (the prior attempts' design). *Insufficient:* it compiles Subspace presentation into spacedock's release train — the coupling this design declared as its biggest liability and the captain's attempt-2 revise. Every guarantee relocates without loss (duty-split table), so the coupling buys nothing.
- **Caller-owned retention directory the override script never auto-deletes** — serves AC-1, AC-2. *Alternative:* rely on subspace-tui's own `--result` write and package retention. *Insufficient:* the real-binary probe shows the provider writes nothing and removes the zero-entry package on the blank-float EOF path, and the interim skill compounds it with `trap rm -rf`. The script must own a directory it never deletes.
- **Diagnostics capture (argv + exit code + stderr) when the child yields no result** — serves AC-1, AC-2. *Alternative:* capture only the result/stdout. *Insufficient:* the EOF path produces no result and only one stderr line; without argv/exit/stderr there is nothing to retry from ("no recovery path emitted", finding 12).
- **Provider id normalization gated on digest validation, homed in the recorder** — serves AC-4. *Alternative:* a presentation-side normalizer, or normalizing before validating the digest. *Insufficient:* the recorder already verifies the attempt Briefing id/digest, so normalizing there removes a duplicate implementation; normalizing before validation would launder an unverified result into the attempt id.

## Behavioral test plan

Behavior-first fixtures driving the override script against a fake subspace-tui child; no live TUI, no host smoke, no committed prose-greps. The retention items are the override script's drive suite (its home repo); the id-mapping item lands on the recorder surface; the decoupling item checks the spacedock build. Estimated cost: medium.

1. **Retention drive suite (AC-1, AC-2).** One fixture per exit path (approve, revise, hold/leave-open, EOF, crash, validation-failure, retention-write-failure, and launcher/controller death after the result is written — finding 15); assert retention-directory contents, the launcher-death approval surviving, and the N/N-vs-0/N baseline contrast. This is `scratchpad/retention-spike.sh` (12/12) promoted to the override script's committed tests.
2. **Completion boundary (AC-3).** Child emits a pane marker then blocks; assert no success/validated result until exit. Mutants (return-on-pane, resolve-before-exit) fail.
3. **Id normalization (AC-4), recorder fixture (proposed to 3k).** Provider-minted id with matching/mismatching digest; assert normalize-on-match, reject-on-mismatch, on whichever surface owns id-mapping after the proposal resolves.
4. **Channel selection + no-subspace fallback (AC-5).** Valid multi-artifact briefing → canonical title + Reference presence; missing binary and version-mismatch → fall back to chat, one line naming both remedy and fallback, zero launch, no retention directory.
5. **Decoupling check (AC-6).** `go list -deps ./cmd/spacedock` names no subspace package; the CLI verb surface has no `gate review`; entity frontmatter byte-unchanged across override fixtures.

Riskiest mechanism first: item 1 (atomic retention on success AND failure) — already exercised end-to-end in the spike (12/12).

## Documentation change proposal

Under the reframe there is no `spacedock gate review` verb to document. The doc change is to the present-gate channel description in `docs/site` (the presentation-side entry; the recorder's frontmatter/gates doc changes are 3k's):

```diff
+### Gate presentation channels
+
+The First Officer presents a gate through the present-gate skill. By default it renders the
+gate-review template in chat; the captain's decision is recorded through the recorder either way.
+A workflow or session may declare an override channel — a hardened Subspace float that presents
+the design as a briefing package in a blocking review and returns a retained result (the review
+log, the result, and diagnostics survive on success **and** on failure). Selecting the override
+first probes for the presenter; if Subspace is absent or version-mismatched, presentation simply
+falls back to chat, naming the install remedy and the fallback, with no side effects. A missing
+presenter changes only the channel; it never blocks the gate. Spacedock's binary carries no
+Subspace code — the override is an opt-in script the provider supplies.
```

## Out of scope

- Writing the resolution into `gates:` frontmatter, or any entity-frontmatter mutation (the recorder, 3k).
- Computing blockers, execution holds, eligibility, or applying the gate action (h1, `gate-blockers-and-eligibility`).
- A `spacedock gate review` binary verb, or any Subspace import/shell-out in the spacedock binary (removed by the reframe; the decoupling is AC-6).
- Building or repairing the subspace-tui transport, the `--review-v1` + `--result` coexistence, or any subspace-tui product change (the override script consumes the surface; the provider ships and tests the override).
- Implementing the recorder-side id-mapping/result-validation itself (proposed to 3k via the change protocol; 3k owns whether to accept the owner-tag move).
- Rendering ProbeResult/comparison UI (a recorded Subspace product gap; the override carries the frozen semantic-delta summary as a supporting Reference, it does not render it).
- The probes companion (`gate-review-probes.md`) as a provider-owned convention — it rides along, not built here.

## Stage Report: ideation

- DONE: A concrete command design against the banked 3k cycles 11-12 lifecycle: blocking presenter, atomic result/log/diagnostics retention on success AND failure (the destroyed hold-path result and the blank-float EOF are the red fixtures), and the provider id-mapping implemented as SPECIFIED in 3k's gate-resolution-frontmatter-contract.md.
  `## Required capability` designs the six-step lifecycle (validate → announce → blocking child → owned retention dir → completion-is-exit+validated+retained → no frontmatter write) matching cycle 11-12's addressable-blocking-presenter and non-terminal-pane/timeout discipline; `## Provider id-mapping` implements the contract's normalize-after-digest-validation rule; AC-1/AC-2 carry both red fixtures as fixtures.
- DONE: The cross-repo dependency declared honestly: what needs subspace-tui surfaces vs what ships spacedock-side now; the working-copy-skill ritual as the measured interim baseline it must beat.
  `## What ships spacedock-side now vs. the subspace-tui dependency` splits the two: the full validate/launch/retain/id-map lifecycle ships spacedock-side (testable against a stub child); the briefing-package + `--result` coexistence and the non-EOF transport are the declared subspace-tui gaps. The `review-local-zellij` interim is named as the baseline, and AC-1 measures the command beating it (N/N vs 0/N retention).
- DONE: Expected surface + tolerance declared; riskiest unverified mechanism spiked end-to-end first or an auditable no-spike-needed recorded.
  `## Expected surface + tolerance` declares ~350-550 prod LOC ≈ equal test LOC, ~2-3 `internal/` files + 1 `cmd/` verb, tolerance 2×, with a hard self-check fencing off 3k/h1/subspace scope. `## Spike` records the end-to-end retention spike (7/7, `scratchpad/retention-spike.sh`) plus the exact-tip binary probe proving the provider retains nothing on the blank-float EOF path.

### Summary

Designed `spacedock gate review` as the presentation half split from 3k: a blocking presenter that owns a never-deleted retention directory and atomically retains result/log/diagnostics on every exit path, reproducing and beating the interim ritual's two destruction defects. The riskiest mechanism (atomic retention on success AND failure) was spiked end-to-end against the real tip binary before design lock — the provider writes nothing to `--result` on the blank-float EOF path, so the command owns retention. Scope holds strictly to the resolution/briefing side: no frontmatter writes (3k), no application/eligibility (h1), no subspace-tui product work.

## Stage Report: ideation (cycle 2)

Preflight fold applied (the first decline: the destroyed-approval fixture).

- DONE: Extend the evidence base to findings 1-15.
  Bumped both evidence lines (overview + Problem) to findings 1-15, cited 3k's attempt-7 resolution provenance note (verified at the source: index.md line 164), and marked finding 14 as the deliberate live-session numbering skip.
- DONE: Add finding 15 as the controller/launcher-death red fixture, with the AC-2 clause.
  Added a Problem bullet naming finding 15 (a dying launcher unlinked the captain's own approval) as the primary red fixture; extended AC-2 with the launcher-death clause (the caller-owned `--result` path survives because the launcher owns no cleanup over it) and its test. Exercised it: spike fixture C kills a real launcher process after it writes the result — interim scratch is unlinked (finding 15 reproduced), caller-owned result survives. Spike now 9/9.

### Summary

The fold strengthens the evidence base with finding 15 — the strongest red fixture for this exact mechanism, since it is the caller-owned retention directory's whole reason to exist. Verified the provenance at 3k's attempt-7 note rather than taking it on faith, then proved the fixture by exercising a real killed launcher process (spike 9/9), not by asserting it.

## Stage Report: ideation (cycle 3)

Captain gate feedback (attempt 1, revise): "what happens when the user does not have subspace installed? what's the fallback?"

- DONE: Design the no-subspace fallback — detection, fallback, shape.
  Added `## No-subspace fallback`: detection resolves + version-gates the TUI binary before any side effect and, on absence/mismatch, exits non-zero naming both the install remedy and the chat fallback with zero side effects; fallback is chat presentation recorded through the recorder exactly as every 0260 gate was, presentation-agnostic, under the recording-identity ruling (verified at `docs/roadmap/durable-decisions/index.md` Constraints line 25); shape mirrors the teams-unavailable-selects-bare-dispatch ruling (verified at `dispatch-failure-retry-rung.md` — an unavailable capability is not a mode, just the ordinary condition selecting the alternate path). Made detection step 1 of Required capability, before any side effect.
- DONE: Make it checkable — one AC clause + test-plan line, exercised.
  Strengthened AC-5 to require non-zero exit + a message naming both remedy and fallback + zero launch + no retention directory created (entity byte-unchanged), with the matching test-plan line; exercised it as spike fixture D (absent presenter → exit 3, message names remedy+fallback, no dir created). Spike now 12/12. Added the fallback sentence to the doc diff.

### Summary

The no-subspace answer names an existing practice rather than inventing one: detection fails clean and names the fallback, chat presentation records through the same recorder under the recording-identity ruling, and a missing presenter is an ordinary channel selection — not a mode, never blocking the gate. Grounded both cited rulings at their sources and proved the checkable clause by exercising fixture D (12/12).

## Stage Report: ideation (cycle 4)

Captain gate feedback (attempt 2, revise): "i am wondering if this should be left as an overridable of the current present-gate skill, so that they are not coupled."

- DONE: Assess the decoupling honestly, then reframe the architecture around it.
  Read the present-gate skill (its default and only behavior is chat presentation) and confirmed the reframe fits: presentation becomes an overridable present-gate channel (default chat, override the hardened float script), the spacedock binary carries zero Subspace code, and the cross-repo release-train coupling — this design's biggest declared liability — dissolves to the opt-in override script. Rewrote the overview, Expected surface, Problem close, the channel section (replacing the binary lifecycle), and the id-mapping section.
- DONE: Split the duties by their honest home; move the id-mapping implementation and amend the owner tag.
  Added `## Splitting the duties honestly` (table): retention → override script + committed drive suite (the 12 spike fixtures); result validation + id-normalization + record handoff → recorder-side (already 3k's parse-and-verify duty); detection/fallback → the skill's channel selection. Proposed the owner-tag amendment (id-mapping specified AND implemented recorder-side) to 3k via the change protocol; did not edit 3k's surface.
- DONE: Re-estimate the surface with tolerance; state the honest counter-case.
  Spacedock Go surface for presentation → ~zero. New surface declared: present-gate skill prose (~25-40 lines) + one hardened override script (~80-140 LOC + the 12-fixture drive suite, homed in the subspace repo) + a small recorder-verb ask; tolerance 2× against that surface. `## Honest assessment` records the one load-bearing condition — the override script MUST carry a committed CI-run drive suite or it repeats `review-local-zellij`'s untested-script defect (the exact class this task exists to kill); that condition fixes the script's home (subspace repo, already `go test`-covered). Found no piece that genuinely cannot be a skill override without losing a guarantee.
- DONE: Keep the no-subspace fallback, probe ritual, and every red fixture; make the decoupling itself checkable.
  Fallback is now the default channel (unchanged behavior); probe ritual and fixtures A-D retained. Added AC-6 (VALUE): the spacedock binary depends on no Subspace package and exposes no presentation verb — a build/dependency assertion (not a prose-grep) whose count regresses the wrong way if a coupling returns. Reframed AC-1/2/3/5, the test plan, and the doc diff (present-gate channel, no `gate review` verb).

### Summary

The captain's question is well-founded and the answer is yes: the binary was the wrong vehicle. Presentation moves to an overridable present-gate channel, the spacedock binary ends Subspace-free (AC-6, measurable), and every guarantee relocates without loss — retention to a testable override script (the spike already proved the contract in bash), validation and id-mapping to the recorder where 3k's contract already puts result verification. The single load-bearing condition, surfaced for the captain rather than buried: the override script must carry a committed drive suite (the 12 fixtures), or the reframe reintroduces the untested-script defect it exists to remove.

### Feedback Cycles

- Cycle 1: REJECTED — Roborev job 541; surface 4 files/49 added, 3 removed vs estimate skill 25-40 lines + tests 20-60 lines + docs 8-16 lines (within 2×); AC unchanged.
- Cycle 2: PASSED — Roborev job 542; surface 4 files/49 added, 3 removed vs estimate skill 25-40 lines + tests 20-60 lines + docs 8-16 lines (within 2×); AC unchanged.
- Cycle 3: REJECTED — fresh validation and Roborev job 1955; surface 4 files/49 added, 3 removed vs estimate skill 25-40 lines + tests 20-60 lines + docs 8-16 lines (within 2×); AC unchanged.
- Cycle 4: REJECTED — Captain design reset after validation cycle 3: replace advisory-adoption plumbing with the gate-room/minimal-binding Result boundary and derive provider association behind the scaffold; surface 4 files/49 added, 3 removed vs estimate skill 25-40 lines + tests 20-60 lines + docs 8-16 lines (within 2×); AC unchanged.
- Cycle 5: REJECTED — Roborev job 1974; surface 15 files/708 LOC vs pre-edit declaration 9 source/docs files plus fixture set/~585 changed lines (121%); AC unchanged.

## Stage Report: implementation

- DONE: Ship the present-gate overridable-channel contract with chat fallback, complete canonical Briefing presentation, exact retained Result handoff, and honest association to what the reviewer actually saw.
  Commits `cf6008fd` and `612b72fc` add default chat, the six-part override contract, exact retained Result handoff to `gate record`, and one-to-one provider Artifact/Reference mapping to the complete canonical inventory.
- DONE: Prove missing/mismatched presenter falls back without side effects and the full-package journey never promotes a single-file float into a complete-package association.
  The skill makes availability/version probing precede package and retention creation; the provider's 12-fixture drive suite remains the required pinned cross-repo proof and is not claimed as a local run, while `TestGateRecordConsumesExactResultOnlyWithCompleteAssociation` rejects the primary-only fixture and proves the entity stays byte-unchanged.
- DONE: Keep the Spacedock binary Subspace-free and preserve the provider-owned hardened override script plus committed drive suite as the named pinned cross-repo release condition.
  `go list -deps ./cmd/spacedock` returns no Subspace dependency; `TestGatePresentationRemainsOutsideBinary` fails if a `gate review` verb appears or if rejecting it changes the working directory; the skill names the pinned provider script and CI suite as release eligibility.
- DONE: Reconcile implementation surface against the pre-edit declaration.
  Actual surface is 15 skill lines, 8 docs lines, and 26 test-line additions with 3 replacements across 4 files, below the declared prose estimate and within the 2× test/docs tolerance; no Go production, recorder, provider, gate-state, or frontmatter surface changed.
- DONE: Run repository and documentation verification.
  PASS (exit 0): `go test ./...`; PASS (exit 0): `go test ./... -race`; PASS (exit 0): `mkdocs build --strict`; PASS (exit 0): `gofmt -w ./cmd ./internal`; PASS (exit 0): `git diff --check`; final code-worktree `git status --short` output was empty (clean).
- DONE: Request and triage Roborev on the completed commits.
  Job 541 found one material AC-5 ambiguity about provider References; `612b72fc` made every Artifact/Reference mapping explicit without expanding recorder scope, and job 542 returned `No issues found.`

### Summary

The present-gate skill now owns an overridable rendering channel: chat remains the side-effect-free default, and a provider override must present the complete canonical package, retain exact results and diagnostics, and associate only content the reviewer saw. The implementation leaves the binary Subspace-free and keeps the hardened provider script plus its committed drive suite as the pinned cross-repo release condition.

## Stage Report: validation

- DONE: Reproduce the complete canonical Briefing journey: exact question and every Artifact/Reference presented, exact retained Result handed off, and association limited to content the reviewer actually saw.
  `TestExactCanonicalBriefingIsIndependentAssociationInventory` binds the exact question and three-Artifact inventory at digest `sha256:0a54f1ba...`; `TestGateRecordConsumesExactResultOnlyWithCompleteAssociation` maps `artifact:primary`, `reference:entity-snapshot`, and `reference:recorder-contract` one-to-one, then records exact Result digest `sha256:46096103...` and decision `revise`.
- SKIPPED: Attack missing/mismatched presenter fallback and primary-only/single-file promotion paths, proving zero side effects on fallback and fail-closed no-mutation association rejection.
  The provider fallback fixtures are absent from this repository and were not claimed; the in-repo primary-only association attack returned exit 1 with `complete presentation mapping` and left the entity byte-identical.
- DONE: Verify the binary stays Subspace-free, the provider script/committed-drive-suite remains an honestly named pinned release condition, and all ACs/tests/surface claims survive detached adversarial review.
  `go list -deps ./cmd/spacedock` found zero Subspace dependencies; `gate review` remains absent; the skill makes the provider script plus committed CI suite a release condition, and the detached audit found no in-repo material defect.
- SKIPPED: AC-1 (VALUE) — No presented decision is lost on any exit path.
  This proof belongs to the provider-owned 12-fixture drive suite; no provider repository or pinned revision was supplied, so validation did not claim a local run.
- SKIPPED: AC-2 — Retention survives every failure class, including launcher/controller death.
  The launcher-death, retention-write, crash, EOF, hold, and validation-failure fixtures remain the same provider-owned release condition.
- SKIPPED: AC-3 — Pane/session creation and wait-timeout are never completion.
  The blocking-child and return-on-pane mutants require the provider script and committed drive suite, which are outside this checkout.
- DONE: AC-4 — The recorded result is keyed to the attempt briefing id, only after digest validation (recorder-homed, proposed).
  The matching fixture normalized both provider ids to the attempt id; an independent mismatch run returned exit 1 and preserved the entity and Result bytes. Removing the revision check made that audit fail by closing the gate.
- SKIPPED: AC-5 — The override channel validates the briefing and derives the title before any launch; an absent or version-mismatched presenter falls back to chat with zero side effects.
  The shipped skill states the probe-first, zero-side-effect fallback contract, but only the provider-owned suite can prove launch count, title derivation, and retention-directory absence.
- DONE: AC-6 (VALUE) — Presentation adds zero Subspace coupling to the spacedock binary, and no channel mutates entity frontmatter.
  Dependency count was zero; the absent-verb test returned exit 2 and left its working directory unchanged; all recorder rejection controls preserved the entity bytes.
- DONE: Audit the declared 4-file/15-skill-line/8-doc-line/26-test-line surface.
  Diff from `fa240a76` is exactly 4 files: 15 skill additions, 8 docs additions, 21 CLI-test additions, and 5 additions/3 replacements in contractlint; no Go production, provider, recorder, or frontmatter file changed.
- DONE: Run detached adversarial controls.
  Mutants accepting `gate review`, trusting association-declared inventory, skipping exact-Result digest binding, changing the exact question, and skipping canonical-revision binding each broke the named focused test or independent audit.
- DONE: Verify Roborev jobs 541/542 closed all material findings without ownership crossover.
  Job 541's sole Reference-association finding is corrected by `612b72fc`; job 542 reports `No issues found.` No provider, recorder, or gate-state ownership crossed. Roborev metadata still marks both review records `closed:false`, an administrative state rather than an unresolved code finding.
- DONE: Run normal, race, documentation, formatting, and cleanliness checks.
  PASS: `go test ./...`, `go test ./... -race`, pinned-env `mkdocs build --strict`, `gofmt -w ./cmd ./internal` with no diff, `git diff --check`, and clean committed implementation worktree at `612b72fc`.
- DONE: Recommend PASSED for the in-repo deliverable, with the cross-repo release condition explicit.
  No material in-repo finding remains; release eligibility still requires a pinned provider revision carrying the hardened override script and its committed 12-fixture CI drive suite.

### Summary

Fresh detached validation passed the four-file Spacedock deliverable and five claim-breaking controls. The provider transport remains deliberately outside this repository: its pinned script and 12-fixture CI suite are an unmet cross-repo release condition, not local test evidence.

## Stage Report: validation (cycle 2)

- DONE: Run and pin the provider-owned retained-delivery drive suite at sibling Subspace revision 198f76238aeb74ff38900e17b751f0460d0c55ee, mapping results explicitly to previously skipped AC-1, AC-2, AC-3, and AC-5.
  `scripts/tests/subspace-r-provider-retained-delivery-test.sh` passed at exact `198f7623`; its approve/revise/hold/open, blank/EOF, crash, invalid-result, retention-write, launcher-death, alive-child, missing/mismatched-presenter, complete-package, and title rows fail on deletion, early delivery, relaunch, or inventory drift.
- DONE: Replay the complete canonical Briefing, Result, and presented-inventory association at Spacedock candidate 612b72fc; distinguish provider defects from the current Codex/Safehouse Zellij transport limitation and make no sibling-repository edits.
  Subspace's frozen Result `sha256:46096103...` and association `sha256:95ca15ab...` are byte-identical to Spacedock fixtures; recursive inventory checks and the exact recorder test passed, while the primary-only map failed without mutation. No provider defect appeared; no headed captain float is claimed because this Codex turn does not expose the agreed `/subspace:r gate <gate-room>` surface, and the private 4p vector is not a substitute.
- DONE: Reissue a fresh exact-tip PASSED or REJECTED recommendation for all six ACs, preserving the zero-Subspace binary boundary and treating the verbose internal 4p vector as implementation plumbing rather than xb's agent-facing interface.
  **PASSED** at Spacedock `612b72fc` with provider `198f7623`: all six ACs have executable or retained-state evidence, zero material findings remain, and `/subspace:r <file.md>` plus `/subspace:r gate <gate-room>` remains the public shape.
- DONE: AC-1 (VALUE) — No presented decision is lost on any exit path.
  The provider matrix retained complete bundles for approve, revise, hold, and open; blank/EOF, crash, validation failure, retention failure, and launcher death retained every produced result/log/inventory/diagnostic byte, so removing any case artifact fails the suite.
- DONE: AC-2 — Retention survives every failure class, including launcher/controller death.
  Child exit `42` and launcher exit `43` propagated once; launcher-death kept the non-empty Result, log, inventory, argv/stderr diagnostics, and death marker, while retention-write failure kept the Result plus its error.
- DONE: AC-3 — Pane/session creation and wait-timeout are never completion.
  The alive-child fixture published a pane marker and Result while holding the exact child alive; entry return, validation, child-exit publication, and delivery stayed absent until release, then occurred exactly once.
- DONE: AC-4 — The recorded result is keyed to the attempt briefing id, only after digest validation (recorder-homed, proposed).
  The exact Spacedock recorder replay normalized the provider envelope only through the complete digest/revision-bound association; missing/changed nested Reference and canonical id/revision mutations all failed closed.
- DONE: AC-5 — The override channel validates the briefing and derives the title before any launch; an absent or version-mismatched presenter falls back to chat with zero side effects.
  The complete provider row derived `Subspace — Ship the complete package?` and ordered 2 Artifacts plus 2 recursively reached References; missing/mismatched presenter rows returned `127`/`2` before host preflight or launch, leaving chat selection to the probe-first Spacedock channel contract.
- DONE: AC-6 (VALUE) — Presentation adds zero Subspace coupling to the spacedock binary, and no channel mutates entity frontmatter.
  `go list -deps ./cmd/spacedock` found zero Subspace dependencies; the absent `gate review` command and recorder rejection controls changed no working directory or entity bytes.
- DONE: Run the required repository verification at the exact Spacedock tip.
  PASS: focused cross-repo tests, `go test ./...`, `go test ./... -race`, `uv run --with-requirements docs/requirements.txt mkdocs build --strict`, `gofmt -w ./cmd ./internal` with no diff, `git diff --check`, and clean code worktree.

### Summary

The previously skipped provider-owned evidence is now pinned and reproduced, and the exact provider Result/inventory/association crosses into Spacedock unchanged and fails closed under independent mutations. Validation recommends **PASSED** for all six ACs at the exact tips; the provider's verbose deterministic vector remains private plumbing, not a new public Spacedock or agent-facing interface.

## Stage Report: implementation (rework)

- DONE: Rebase `spacedock-ensign/gate-review-presentation-command` onto current `origin/main` and resolve only the two known semantic conflicts in `internal/cli/gate_test.go` and `internal/contractlint/fo_function_reference_invariant_test.go`.
  Rebasing `612b72fc` onto `dd6bd114` produced only those conflicts; candidate `4779fff8` preserves main's gate tests and adds the presentation refusal beside them, while host baselines equal main plus xb's measured 2,872-byte load.
- DONE: Preserve xb's approved behavior and intended four-file delta: presentation stays an overridable `present-gate` channel; complete association remains required; the binary remains Subspace-free.
  The final diff against `origin/main` is exactly the original four files and `+49/-3`; complete provider Artifact/Reference mapping remains in the skill/docs, the recorder rejection test passes, and `go list -deps ./cmd/spacedock` names no Subspace dependency.
- DONE: Reconcile with merged gate-lifecycle/advisory-round behavior from commits `03b1a7fc`, `f06cce04`, `c355fbe4`, and `c9dfc491`; do not invent compatibility constraints for prototype behavior.
  All four commits are ancestors of `origin/main`; round record/validate, owning-workflow discovery, eligibility/consume-once, and presentation-refusal tests pass, with the latter expecting main's real `record|validate|eligibility|consume` command surface.
- DONE: Before committing, report actual changed files/LOC versus xb's pre-rebase 4 files, +49/-3, and treat unexplained product-surface drift as a blocker.
  Pre-commit audit reported docs `+8`, CLI test `+21`, load-ratchet test `+5/-3`, and present-gate skill `+15`; no extra file, production Go surface, provider/recorder change, or unexplained drift appeared.
- DONE: Run focused affected tests, gofmt as applicable, `go test ./...`, `go test ./... -race`, mkdocs strict, and `git diff --check`.
  Focused CLI/contractlint suites, `gofmt -w ./cmd ./internal`, normal/race suites, `uv run --with-requirements docs/requirements.txt mkdocs build --strict`, and diff checks all passed; removing a merged verb, complete association, or per-host load increment breaks its named test.
- DONE: Commit the rebased integration candidate and write an implementation/rework report, but do not self-validate or mutate workflow frontmatter/gates.
  Rebased commits are `f8fd5c2a` and `425812bc`; integration reconciliation is `4779fff8`. This append-only body report makes no validation recommendation and leaves frontmatter/gates untouched for a fresh validator.

### Summary

The rebased xb candidate preserves the approved four-file `+49/-3` boundary while composing with merged gate lifecycle and folder-form advisory rounds. Candidate `4779fff8` is clean and fully implementation-tested; independent validation remains the next gate.

## Stage Report: validation (cycle 3)

- DONE: Independently verify rebased candidate 4779fff8 preserves xb’s six ACs, exact four-file +49/-3 boundary, and composes semantically with the merged gate lifecycle/advisory-round tests without compatibility inventions.
  The diff from `dd6bd114` is exactly docs `+8`, CLI test `+21`, contractlint `+5/-3`, and skill `+15`; round record/validate, owning-workflow, eligibility/consume-once, absent-presentation-verb, association, and host-load tests pass against the real `record|validate|eligibility|consume` surface.
- DONE: Reproduce the pinned Subspace 198f762 retained-delivery suite and complete Briefing/Result/association path against candidate 4779fff8; explicitly establish that the old 612b72fc approval cannot authorize this new candidate and identify the replacement Briefing inputs.
  Exact provider commit `198f76238aeb74ff38900e17b751f0460d0c55ee` passed its 12-fixture suite; the exact Result/full-association recorder path passes, while an advisory Result with that association but no adoption note fails with `advisory Result requires --adoption-note`.
- DONE: Establish approval freshness and replacement Briefing inputs.
  The closed attempt’s immutable question and approval reason name only `612b72fc`; although `gate eligibility` reports its application pending/eligible, it cannot authorize `4779fff8`. A successor Briefing must freeze the `4779fff8`/`dd6bd114` four-file diff, current entity plus this validation, merged lifecycle/advisory evidence, provider `198f7623` suite evidence, and the corrected handoff; only its fresh exact Result and complete association can authorize landing.
- DONE: AC-1 (VALUE) — No presented decision is lost on any exit path.
  Provider rows for approve, revise, hold/open, blank/EOF, crash, invalid result, retention failure, and launcher death retain their promised bytes; deleting any retained artifact breaks the pinned suite.
- DONE: AC-2 — Retention survives every failure class, including launcher/controller death.
  The pinned suite preserves Result/log/inventory/diagnostics across child and launcher failures and proves nonzero status propagation without relaunch.
- DONE: AC-3 — Pane/session creation and wait-timeout are never completion.
  The alive-child row withholds delivery and validation until the blocking child exits, then publishes exactly once.
- FAILED: AC-4 — The recorded result is keyed to the attempt briefing id, only after digest validation (recorder-homed, proposed).
  The recorder’s complete association and digest normalization work, but the shipped skill’s prescribed command omits the required `--adoption-note`; Subspace’s advisory Result therefore fails at the supported handoff boundary. This is a material outcome defect and a narrow same-layer fix, not a design reset.
- DONE: AC-5 — The override channel validates the briefing and derives the title before any launch; an absent or version-mismatched presenter falls back to chat with zero side effects.
  Complete-package/title and missing/mismatched-presenter rows pass at the pinned provider revision before host launch or retention creation.
- DONE: AC-6 (VALUE) — Presentation adds zero Subspace coupling to the spacedock binary, and no channel mutates entity frontmatter.
  `go list -deps ./cmd/spacedock` reports zero Subspace dependencies; absent-verb and recorder rejection controls preserve the working directory/entity.
- DONE: Run focused tests, gofmt, go test ./..., go test ./... -race, mkdocs strict, diff/cleanliness checks, detached adversarial review and Roborev as applicable; classify every finding by materiality and issue a fresh exact-tip PASSED or REJECTED report.
  Final uncontaminated normal/race suites, strict MkDocs, gofmt-no-diff, diff check, and clean code worktree pass; Roborev panel job 1955 is REJECTED at exact tip.
- FAILED: Material outcome finding — advisory Result handoff omits adoption authority.
  Supported trigger: provider returns `status:"advisory"`; `gate record` exits 1 without `--adoption-note`. Fix AC-4 by branching on Result authority, supplying captain-authorized adoption text, and behaviorally testing the exact prescribed invocation.
- SKIPPED: Deferred evidence risk — provider revision is not pinned by a code-repository release check.
  Current gate evidence pins and reproduces `198f7623`, so no present AC lacks proof; promote if a stable release consumes the skill without a fresh exact provider pin and suite result.
- SKIPPED: Deferred evidence risk — absent-verb test compares directory entry count, and host-load ratchets enforce ceilings rather than exact byte equality.
  Current early-return path creates nothing and all hosts equal main plus 2,872 bytes; promote the first if the early-return guard moves, and the second if ratchets become exact accounting rather than upper bounds.
- FAILED: Fresh exact-tip recommendation.
  **REJECTED** at Spacedock `4779fff8` with provider `198f7623`: AC-4’s documented supported handoff cannot record the provider’s advisory Result; no old `612b72fc` Resolution authorizes this candidate.

### Summary

The rebase, four-file boundary, provider retention suite, merged gate behavior, and five ACs validate cleanly. One material handoff defect remains: the present-gate override contract must carry captain-authorized adoption for advisory Results and prove that exact invocation before a fresh candidate-bound Briefing can be approved.
