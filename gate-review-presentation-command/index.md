---
title: Gate presentation as an overridable channel with atomic result retention
status: implementation
source: "Split from the gate-recorder task (3k), captain-approved 2026-07-21. The subspace-coupled presentation half; 3k cycles 11-12 are its banked design history."
id: xbatj4hxtxw9t83vvmfem27f
gates:
    version: 1
    current:
        gate: gate:docs-dev:xb:ideation
        attempt: gate-attempt:xb-ideation-3
    records:
        - id: gate:docs-dev:xb:ideation
          stage: ideation
          current-attempt: gate-attempt:xb-ideation-3
          attempts:
            - id: gate-attempt:xb-ideation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:docs-dev:xb:ideation:briefing-1
                digest: sha256:bfb87cff3c021c17af9a9d8a999cb682cde972d73d292a5838f21baec7a240de
                room-ref: "./review/ideation/briefing-1"
                note: "Multi-artifact package: gate summary, frozen entity snapshot, frozen recorder-contract snapshot — each digest-pinned inside the briefing; the digest above binds briefing.json itself."
              resolution:
                type: Resolution
                id: resolution:actor-1784637808799451000
                briefing: briefing:docs-dev:xb:ideation:briefing-1
                by: person:reviewer
                at: 2026-07-21T12:43:28Z
                decision: revise
                reason: "Annotation on the command definition: what happens when the user does not have subspace installed? what's the fallback?"
              application:
                action: feedback
                target-stage: ideation
                state: consumed
              note: "Provider result and log retained in-room. Routed to the live worker; attempt 2 opens at re-presentation."
            - id: gate-attempt:xb-ideation-2
              sequence: 2
              previous-attempt: gate-attempt:xb-ideation-1
              state: closed
              briefing:
                id: briefing:docs-dev:xb:ideation:briefing-2
                digest: sha256:98b247e79ce5e88285a3a903a1c339dab6c658d8e89f44d28f1a612f3faccc33
                room-ref: "./review/ideation/briefing-2"
                note: "Same package shape; the design now carries the no-subspace chat fallback (captain's attempt-1 question), exercised as spike fixture D."
              resolution:
                type: Resolution
                id: resolution:actor-1784638768658049000
                briefing: briefing:docs-dev:xb:ideation:briefing-2
                by: person:reviewer
                at: 2026-07-21T12:59:28Z
                decision: revise
                reason: "Annotation on the Go surface estimate: i am wondering if this should be left as an overridable of the current present-gate skill, so that they are not coupled."
              application:
                action: feedback
                target-stage: ideation
                state: consumed
              note: "Architecture reframe routed to the live worker: presentation as an overridable channel of the present-gate skill (default chat, override the hardened float script); validation/id-mapping/record-handoff duties move to the recorder surface where they belong. Attempt 3 at re-presentation."
            - id: gate-attempt:xb-ideation-3
              sequence: 3
              previous-attempt: gate-attempt:xb-ideation-2
              state: closed
              briefing:
                id: briefing:docs-dev:xb:ideation:briefing-3
                digest: sha256:7fb094109c3ad7c873dcd66cf77295c555395394c20da9823af4fed7fd9abc37
                room-ref: "./review/ideation/briefing-3"
                note: "The reframed design: overridable present-gate channel, zero-Go surface, recorder-side validation ask, the committed-drive-suite condition homed subspace-side."
              resolution:
                type: Resolution
                id: resolution:actor-1784640771337164000
                briefing: briefing:docs-dev:xb:ideation:briefing-3
                by: person:reviewer
                at: 2026-07-21T13:32:51Z
                decision: approve
              application:
                action: advance
                target-stage: implementation
                state: consumed
              note: "Subspace multi-artifact float, no annotations; provider result and review log retained in-room. The decoupled architecture is approved: presentation as an overridable present-gate channel, the binary subspace-free by checkable criterion. ADOPTION PROVENANCE (captain-confirmed in chat, 2026-07-21): person:reviewer at this float was the captain personally; the FO adopted the provider's advisory result as the captain's binding approval on that basis. The provider envelope (binding:false) is evidence, not the binding record — this note is the durable authorization the promotion previously lacked."
sprint: durable-decisions
group: recorder
started: 2026-07-21T01:43:36Z
worktree: .worktrees/spacedock-ensign-gate-review-presentation-command
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
