---
id: p1tgy61tbhxj9apvpswqbhcy
title: The skill-first install hint is channel-blind - an edge plugin tells the user to install the stable binary its own version gate then rejects
status: validation
source: "Captain CL in chat, 2026-08-24, reviewing install-sh-edge-prerelease-parity (#756): 'does it work for the skill-first journey that tells the user to run it?' - it does not; follow-up scoped to the skill hint path #756 could not reach"
started: 2026-08-24T19:20:28Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-install-gate-channel-aware-hint
issue:
gates:
    version: 1
    records:
        - id: gate:p1tgy61tbhxj9apvpswqbhcy:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:p1tgy61tbhxj9apvpswqbhcy-backlog-1
              briefing:
                id: briefing:p1tgy61tbhxj9apvpswqbhcy:backlog:attempt-1:revision-1
                digest: sha256:20c0bb050069025cdf5f2b5627e1e28a9b2898b5ba07ead3e4dfc696107d6424
                request-digest: sha256:e389ddd2fbb0e491876aca632fed54018124271aefbb8b98726f2b7afae32f95
                room-ref: ./install-gate-channel-aware-hint/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:p1tgy61tbhxj9apvpswqbhcy:backlog:1
                briefing: briefing:p1tgy61tbhxj9apvpswqbhcy:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-24T19:20:03.982043Z"
                decision: approve
                reason: 'Captain CL in chat 2026-08-24: ''let''s dispatch p1t on to the channel install script PR stack'' - accepts the seed into ideation with stacked delivery on the #756 branch'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:p1tgy61tbhxj9apvpswqbhcy:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:p1tgy61tbhxj9apvpswqbhcy-ideation-1
              briefing:
                id: briefing:p1tgy61tbhxj9apvpswqbhcy:ideation:attempt-1:revision-1
                digest: sha256:d66d1ab902025ff59a0bfff073ddb52edf75d310d2db6af8dbc5f039da8bca81
                request-digest: sha256:f927de54dd5186eda198920e8432c2baf7331a22506f5aef705dd51233ddfb25
                room-ref: ./install-gate-channel-aware-hint/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:p1tgy61tbhxj9apvpswqbhcy:ideation:1
                briefing: briefing:p1tgy61tbhxj9apvpswqbhcy:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-24T19:39:29.493268Z"
                decision: approve
                reason: 'Captain CL in chat 2026-08-24: ''ok approve. keep it lean.'' at ideation attempt-1 (digest d66d1ab9) - accepts the classifier design and the +150/4 surface with a lean directive on the boot-resident prose'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:p1tgy61tbhxj9apvpswqbhcy:validation
          stage: validation
          attempts:
            - id: gate-attempt:p1tgy61tbhxj9apvpswqbhcy-validation-1
              briefing:
                id: briefing:p1tgy61tbhxj9apvpswqbhcy:validation:attempt-1:revision-1
                digest: sha256:2104a34a524cf1538441ef18d01eb9851e0a1b2b7d2badb1391bfb2b45aa0152
                request-digest: sha256:caebbac1327f72c6d968f6c39cc6a746c1747ce0ba0be7bea85a512fcbe4bcf5
                room-ref: ./install-gate-channel-aware-hint/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:p1tgy61tbhxj9apvpswqbhcy:validation:1
                briefing: briefing:p1tgy61tbhxj9apvpswqbhcy:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-24T21:34:30.504004Z"
                decision: approve
                reason: 'Captain CL in chat 2026-08-24: ''p1t - approve'' at validation attempt-1 (digest 2104a34a) - accepts PASSED with the three deferred risks and their promote-conditions'
              application:
                target-stage: done
                state: pending
mod-block: merge:pr-merge
pr: pr-merge:757
---

install.sh now takes SPACEDOCK_CHANNEL (#756), but the skill-first journey never uses it. The FO binary-absent gate (first-officer-shared-core.md:10) hints the channel-less `curl ... | sh` on Linux and plain brew on macOS; fo-install-gate.md has zero channel awareness. On an edge-channel plugin install the journey self-defeats: the skills pin "binary minor 0.27" (prerelease-only until 0.27 goes stable), the hinted command installs stable v0.26.0, and the same gate then aborts on the binary it just told the user to install - the fresh-VM trap of 2026-08-24, surviving in the path that tells humans what to run.

## Problem

The binary-absent arm of the FO version gate emits one hardcoded command per OS. Both are stable-channel commands. On an edge plugin install the skills pin binary minor 0.27, which no stable release satisfies (stable's latest is 0.26.0), so the hinted command lands 0.26.0 and the same gate rejects it on the next boot. The gate holds the evidence of its own channel and never consults it.

Two seed premises turned out to be wrong and are corrected here:

- **macOS does have an edge cask.** `spacedock@next` is published in `spacedock-dev/homebrew-tap` (`Casks/spacedock@next.rb`, 1827 bytes) pinned at `0.27.0-pre8`, against the stable `spacedock` cask at `0.26.0`. `.goreleaser.yaml` publishes it on every tag including prereleases (`skip_upload: false`, vs the stable cask's `auto`). So the fix does not force macOS onto the curl script; each OS keeps its existing install style and only picks the channel-correct variant.
- **A per-channel skill stamp is impossible.** One `goreleaser release` emits both channels from an identical payload; `.goreleaser.yaml` states the two builds "differ ONLY in the devBranch ldflag." The skill bytes are byte-identical across channels, so nothing can be stamped into them at release time. This retires seed candidate 2.

## Risk evidence

The riskiest claim was: **with no binary present, can the FO know its own channel?** Probed before designing.

**1. Env-var detection is dead.** `env` in the Bash tool of a live Claude Code session exposes `CLAUDECODE`, `CLAUDE_CODE_SESSION_ID`, `CLAUDE_CODE_AGENT`, and friends, but **no `CLAUDE_PLUGIN_ROOT` and no `PLUGIN_ROOT`**. Those are injected into hook commands (`hooks.json`), not into the agent's shell. Any design reading the plugin root from the environment would have failed at runtime.

**2. The FO already holds its own absolute path.** `skills/first-officer/SKILL.md` declares: "The skill loader supplies the absolute base directory for this `first-officer` skill when it opens this file. Retain that exact directory as `{first_officer_base}` for the session," and forbids deriving it any other way. It is already load-bearing — the FO cannot read `references/fo-install-gate.md` without it. **Channel detection therefore adds no new mechanism; it classifies a string the FO is already required to hold.**

**3. The marketplace name IS the channel, by documented design.** `docs/site/get-started/install.md`: "The channel is the marketplace name" — stable `spacedock`, edge `spacedock-edge`. Install paths are `~/.claude/plugins/cache/<marketplace>/<plugin>/<version>/`, confirmed in `installed_plugins.json`:

    spacedock@spacedock-edge | 0.27.0-pre8     | …/cache/spacedock-edge/spacedock/0.27.0-pre8
    spacedock@spacedock      | 0.27.0-pre7+dev | …/cache/spacedock/spacedock/0.27.0-pre7-dev

**4. Neither single signal is sufficient — measured, not assumed.** Enumerating all 29 real spacedock plugin dirs in the local cache produced a counterexample for each candidate discriminator:

- marketplace-segment alone misclassifies `cache/spacedock/spacedock/0.27.0-pre7-dev` (a `next` dev build recorded under the stable marketplace name) as stable;
- prerelease-suffix alone misclassifies `cache/spacedock-edge/spacedock/0.27.0` (an edge install on an unsuffixed version) as stable.

Their **union** classifies all observed shapes correctly. A third arm is required for source checkouts (`--plugin-dir`), which carry no marketplace segment at all.

**5. The classifier was exercised, and its failure modes confirmed.** Ten fixture paths (real shapes plus adversarial) all classify correctly, including a `/Users/pre-release-tester/…` home that a naive substring match on `-pre` would have falsely called edge. Two falsifying edits were run and both produce wrong answers on real recorded paths: dropping the marketplace arm calls `spacedock-edge/spacedock/0.27.0` stable; dropping the suffix arm calls `spacedock/spacedock/0.27.0-pre7-dev` stable and the dev checkout stable instead of local. The table can fail.

**Seed candidate 3 (minor-pin-as-proxy) is retired:** deciding whether the pinned minor has a stable release needs a network round trip, which is what `install.sh` already does, and it cannot distinguish "pin not yet released" from "user is behind."

## Proposed approach

Classify the channel from `{first_officer_base}` into three arms — `edge`, `stable`, `local` — and hint that arm's command for the detected OS.

    R="${B%/skills/first-officer}"; V="${R##*/}"; case "$R" in */spacedock-edge/*) echo edge;; *) case "$V" in [0-9]*-*) echo edge;; [0-9]*) echo stable;; *) echo local;; esac;; esac

`local` (no version-shaped directory — a source checkout) routes to the source build the contract already names for the unsupported-OS arm, so it introduces no new command.

**Why the rule lives in the shared core, not only in `fo-install-gate.md`.** The seed asked to keep detection detail in the deferred file. That is only half possible: `fo-install-gate.md` is explicitly *not* loaded inside a sandbox, yet the sandbox arm must still print "the exact OS-aware install command." The classifier must therefore be reachable without the deferred file. The split adopted is: the **rule** (one sub-bullet) in the shared core; the **rationale, the two counterexamples, and the `local` arm's handling** in `fo-install-gate.md`.

### Mechanism justification

**The path classifier** serves AC-1. Simplest alternative considered: *outcome-driven retry* — hint stable, and on a version-too-low re-check re-run with `SPACEDOCK_CHANNEL=edge`. Rejected because `fo-install-gate.md` step 4 is an explicit one-attempt loop bound ("no second install attempt"); the retry would weaken the guard that exists to stop install loops, and would spend a network install to rediscover what the base path already states. Second alternative: *print both and let the human choose* — rejected because the automated install offer must select exactly one command to run, and the human has strictly less evidence than the gate does.

**The extraction test** serves AC-1/AC-2. Simplest alternative: a live-lane journey. Insufficient — the live lanes never reach the binary-absent class (that is exactly the open finding in `install-gate-sentinel-behavior-journey` / xw), so no live path exercises this today.

### Coordination

`xw` (sentinel behavior journey) and `d2k` (binary-present wrong-version upgrade hint) are both still `backlog`, so neither offers a harness to reuse and this task must not block on them. The captive-PATH fixture here is a candidate seed for xw's broader harness. Boundary unchanged: this task owns binary-absent only.

### Before/after wording

`skills/first-officer/references/first-officer-shared-core.md`, the **Binary absent** bullet (currently line 10).

Before:

    - **Binary absent:** retry bare `spacedock` once if `SPACEDOCK_BIN` is unusable. Use `uname -s`, not `doctor`/`OS:`. Linux: `curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh`; macOS: `brew tap spacedock-dev/homebrew-tap`, then `brew install spacedock-dev/homebrew-tap/spacedock`; other OS: unsupported OS, hint `go build -o spacedock ./cmd/spacedock`, ABORT. Outside sandbox read `references/fo-install-gate.md`.

After:

    - **Binary absent:** retry bare `spacedock` once if `SPACEDOCK_BIN` is unusable. Use `uname -s`, not `doctor`/`OS:`. Classify the channel (next bullet), then hint that channel's command. Linux stable: `curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh`; Linux edge: that command with `SPACEDOCK_CHANNEL=edge ` prefixed. macOS: `brew tap spacedock-dev/homebrew-tap`, then stable `brew install spacedock-dev/homebrew-tap/spacedock`, edge `brew install spacedock-dev/homebrew-tap/spacedock@next`. `local` channel or other OS: hint `go build -o spacedock ./cmd/spacedock`, ABORT. Outside sandbox read `references/fo-install-gate.md`.
      - **Channel.** Substitute the retained `{first_officer_base}` for `B` and run: `R="${B%/skills/first-officer}"; V="${R##*/}"; case "$R" in */spacedock-edge/*) echo edge;; *) case "$V" in [0-9]*-*) echo edge;; [0-9]*) echo stable;; *) echo local;; esac;; esac`

`skills/first-officer/references/fo-install-gate.md`, new section appended after "Install-and-resume offer":

    ## Channel selection

    The shared core's Binary-absent bullet carries the classifier itself, because the sandbox arm must name a channel-correct command without loading this file. This section is the rationale it defers.

    - **Why two signals.** The marketplace segment is the documented channel name, but a `next`-branch dev build can be installed under the stable marketplace name (`cache/spacedock/spacedock/0.27.0-pre7-dev`); a prerelease suffix catches that. Conversely an edge install can carry an unsuffixed version (`cache/spacedock-edge/spacedock/0.27.0`); the marketplace segment catches that. Either signal alone misclassifies a real observed install.
    - **`local`.** A base with no version-shaped directory is a `--plugin-dir` source checkout. Do not hint a package install: the repo is present, so hint `go build -o spacedock ./cmd/spacedock` and ABORT. Guessing a channel here would install over a tree the human is editing.
    - **Never widen the match.** Test the version segment, never the whole path: a home directory such as `/Users/pre-release-tester` would otherwise force every install on that machine onto the edge channel.

## Documentation diff

`docs/site/get-started/install.md` tells users to install the **edge plugin** (`claude plugin install spacedock@spacedock-edge`) but shows only **stable binary** commands — the same channel blindness, in the docs. Append to the two install tabs:

    === "macOS (Homebrew)"

         ```bash
         brew tap spacedock-dev/homebrew-tap
         brew install spacedock
    +    # Edge channel (tracks prereleases; conflicts with the stable cask):
    +    # brew install spacedock-dev/homebrew-tap/spacedock@next
         ```

    === "Binary (macOS / Linux)"

         ```bash
         curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh
    +    # Edge channel:
    +    # SPACEDOCK_CHANNEL=edge curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh
         ```

         Installs a checksum-verified binary to `~/.local/bin`.
    +
    +    Match the binary channel to the plugin channel you install below: the edge
    +    plugin pins a binary minor that no stable release satisfies.

## Declared semantic changes

- **Runtime behavior (user-visible):** the binary-absent hint text becomes channel-dependent. Edge installs see different command bytes; stable installs see today's bytes unchanged.
- **New abort path:** a source checkout now classifies `local` and is sent to the source build instead of receiving a package-install hint it should not run.
- **Sandbox arm:** same message shape, channel-correct command bytes.
- **Unchanged:** command grammar, stored formats, write authority, the sentinel loop bound, and `install.sh` itself.

## Out of scope

install.sh behavior (shipped in #756). The wrong-version upgrade journey (d2k). An edge Homebrew cask (its own release-machinery decision).

## Expected surface and tolerance

Estimate net LOC change: **+150, across 4 files** (insertions ~+153, deletions ~-3, reported separately per the workflow's net-not-gross rule).

| File | ins | del | net |
|---|---|---|---|
| `skills/first-officer/references/first-officer-shared-core.md` | +2 | -1 | +1 |
| `skills/first-officer/references/fo-install-gate.md` | +8 | 0 | +8 |
| `skills/integration/install_hint_channel_test.go` (new) | +130 | 0 | +130 |
| `docs/site/get-started/install.md` | +13 | -2 | +11 |

Tolerance: **±40 net lines, ±1 file.**

This is well above the seeded ~+25/2 files. The growth is the test file: the seed's estimate covered only the prose edit. The shared-core prose delta is +1 line but a large character delta (~490 → ~1000 chars on that bullet); flagged because a line count understates the boot-resident cost.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified. Seeded; ideation refines.

**AC-1 (value) — A binary-absent boot on an edge-channel plugin ends with an edge-parity binary, not an abort loop.**
Measured end-to-end against a baseline that can move the wrong way: the stable channel's latest minor (0.26) sits *below* the skills' pin (0.27), so any regression to a channel-blind hint is observable as a landed binary that fails the pin. Verified by executing the exact command the shipped classifier selects for an edge fixture path, in a captive `HOME`/`PATH`, then asserting on the landed binary: `--version` line 1 parses to a minor ≥ the contract's pin, and the `Channel:` line matches `^Channel: edge`. Falsifying change: drop the `SPACEDOCK_CHANNEL=edge` prefix from the hint — the run then lands 0.26.0 reporting `Channel: stable` and both assertions fail.

**AC-2 — The shipped classifier maps every observed real install-path shape to its true channel (`edge` or `stable`).**
Verified by a table test that extracts the classifier one-liner from `references/fo-install.md` (relocated there by the captain's 2026-08-24 design correction; originally the shared core) and executes it (never greps it) against fixture paths recorded from the real plugin cache and `installed_plugins.json`, plus a `--plugin-dir` checkout and an adversarial `/Users/pre-release-tester/…` home. Expected values come from the observed install layout, not from the file under test. The channel space is `edge`/`stable` only: the third `local` arm was struck by captain ruling 2026-08-24 (a skill base with no binary is not a real class), so source-checkout rows now expect `stable` — the version test requires a digit-leading segment, so a non-version directory name falls through to stable rather than reading as a prerelease. Falsifying changes, both confirmed to produce wrong answers on real paths: removing the marketplace-segment arm calls `cache/spacedock-edge/spacedock/0.27.0` stable; removing the prerelease-suffix arm calls `cache/spacedock/spacedock/0.27.0-pre7-dev` stable.

**AC-3 — Stable installs' hint bytes are unchanged, and the contract's stable command cannot drift from the published one.**
Verified by binding two independently maintained values: the Linux stable command extracted from `references/fo-install.md` (relocated there by the captain's 2026-08-24 design correction; originally the shared core) must equal the command published in `docs/site/get-started/install.md`, and the two cask names the contract hints (`spacedock`, `spacedock@next`) must both resolve in the tap (spelled `spacedock-dev/tap` since the same date's sweep). Falsifying change: edit either the contract's or the doc's curl line alone, or rename a cask — the pair diverges and the test fails.

**AC-4 — The macOS edge remedy is a real, pin-satisfying cask, not the curl fallback.**
Verified by resolving `spacedock-dev/homebrew-tap/spacedock@next` and asserting its pinned version's minor ≥ the contract's pin (observed today: cask `0.27.0-pre8` vs stable cask `0.26.0`). Falsifying change: if the edge cask ever stops tracking prereleases, its version falls back to the stable line and the minor assertion fails — which is exactly when the macOS hint would need to revert to the curl form.

## Test plan

**Primary: `skills/integration/install_hint_channel_test.go`** (new, ~130 lines, Go, offline). Follows the blessed extraction shape of `skills/integration/survey_probe_test.go` — extract the shipped runnable one-liner via an anchored regex and `sh -c` it, so the contract and the test cannot drift and the oracle is observed output, never file wording. This is the proof policy's second instruction-file exception, not prose-grep: the expected values originate in the recorded plugin-cache layout, outside the file under test.

- `TestChannelClassifierTable` — AC-2. Table over the fixture paths above; extraction failure (a regex miss) is itself a guard against the one-liner being reshaped.
- `TestStableHintMatchesPublishedDoc` — AC-3. Extracts the curl command from both the contract and `install.md`; asserts byte equality.
- `TestCaskNamesResolve` / `TestEdgeCaskSatisfiesPin` — AC-3/AC-4. Network-gated behind `testing.Short()` and a tap-reachability skip so the offline lane stays deterministic.

**Value run: AC-1, one live captive-PATH install.** Run the classifier-selected edge command with `HOME` and `PATH` redirected to a temp dir, then exercise the landed binary's `--version`. Cost: one network install, ~30s. Recorded in the implementation report as evidence with the observed `Channel:` and version lines pasted. Not committed as a standing test — it downloads a release asset, so it belongs in the report, not the offline lane. The macOS brew arm is covered by AC-4's cask resolution rather than a full `brew install`, which is too heavy for CI and proves nothing the pinned cask version does not.

**Required CI lanes.** The diff touches `skills/first-officer/references/**`, which is `skills/**/references/**`. Per the workflow's Proof policy, that makes the touched hosts' **live lanes REQUIRED green at merge** — the shared core is host-neutral contract, so **every host live lane** (`claude-live`, `codex-live`, `pi-live`) is required, not just the deterministic build/install/offline lanes. A flake there is grounds to re-run to green, never to skip. This is the dominant cost of the task and was weighed in scoping: the prose delta is deliberately confined to one bullet plus one deferred section to keep that cost bounded.

**Adversarial audit.** The shipped contract is one of the four high-stakes surfaces, so a detached read-only audit on a throwaway checkout applies before merge.

## Delivery

Stacked layer per the captain's order and `docs/dev/_mods/pr-merge.md` **Stacked mode**. Implementation branches off `spacedock-ensign/install-sh-edge-prerelease-parity` (PR #756, commit `b00886d3a`), which ships `SPACEDOCK_CHANNEL` — do **not** assume `main` carries the env var. Create the PR with `gh pr create` (never `gh stack submit`, which discards approved title bytes), then join with `gh stack link` using PR numbers confirmed via `gh pr view`. Checks are approved at the **tip** of the stack; a green tip is evidence for every layer beneath it.

## Stage Report: ideation

- DONE: The channel-detection mechanism is chosen with the riskiest claim exercised first - how a binary-absent session actually knows its channel (host plugin-install record vs per-channel skill stamp vs minor-pin-as-proxy) - with the probe result in the new Risk evidence section, or an auditable "no spike needed" naming the proven mechanisms.
  Probed all three seed candidates; see `## Risk evidence`. Candidate 1 (plugin-install record) selected in its cheapest form — classify `{first_officer_base}`, a variable the FO contract already requires it to retain. Candidate 2 refuted from `.goreleaser.yaml` (both channels ship byte-identical payloads; only the devBranch ldflag differs). Candidate 3 refuted (needs a network round trip). A fourth candidate I tried first, reading the plugin root from the environment, was killed by probe: `CLAUDE_PLUGIN_ROOT`/`PLUGIN_ROOT` are absent from the agent's Bash env.
- DONE: The task body carries the value AC measured end-to-end (a binary-absent boot on an edge plugin ends with an edge-parity binary, observed via --version Channel/minor), per-OS hint text for both channels (macOS has no edge cask), a net-LOC surface with tolerance, and declared semantic changes.
  AC-1 measures the landed binary's `Channel:`/minor against the stable channel's 0.26 vs the 0.27 pin — a baseline that moves the wrong way under regression. Surface table gives ins/del/net (+150 net, 4 files, ±40 tolerance); `## Declared semantic changes` lists four. **The item's parenthetical is wrong and is corrected in the body:** macOS *does* have an edge cask — `spacedock@next` is live in the tap at `0.27.0-pre8`, so each OS keeps its own install style instead of macOS being forced onto curl.
- DONE: Specific before/after wording for both skill files (shared-core install line stays boot-lean; detection detail lives in fo-install-gate.md), plus the stacked-delivery plan: implementation branches off spacedock-ensign/install-sh-edge-prerelease-parity (#756 ships SPACEDOCK_CHANNEL), and the live-lane requirement for skills/**/references/** changes is named in the test plan.
  Verbatim before/after for both files under `### Before/after wording`; full doc diff for `install.md` under `## Documentation diff`. `## Delivery` names the #756 base branch and the Stacked-mode ceremony. The test plan's "Required CI lanes" names all three host live lanes as REQUIRED, since the shared core is host-neutral.

### Summary

The riskiest claim — that a binary-absent FO can know its channel — resolved better than the seed assumed: the FO is already contractually required to retain its own absolute skill base directory in order to read its reference files, so channel detection classifies a string it already holds and adds no new mechanism. That is the answer to the staff review's necessity question. The design is a three-arm classifier (`edge`/`stable`/`local`) shipped as one runnable shell line, exercised here against ten fixture paths with two falsifying edits confirmed to break it.

Two seed premises were refuted with evidence rather than carried forward. macOS has a published edge cask (`spacedock@next`, `0.27.0-pre8`), so the fix is "pick the channel-correct variant of the command you already emit" on both platforms, not "force everyone onto curl." And a per-channel skill stamp is impossible by construction, since one goreleaser run emits both channels from an identical payload.

Two design points worth the gate's attention. First, the classifier cannot live only in `fo-install-gate.md` as the seed proposed: that file is never loaded inside a sandbox, yet the sandbox arm must still print an exact install command, so the rule belongs in the shared core and only its rationale defers. Second, the surface estimate is +150 net against the seed's +25 — the growth is the test file, which the seed did not price. Neither single path signal was sufficient; I found a real counterexample to each in the local cache, which is why the shipped rule is their union plus a third arm for source checkouts.

## Stage Report: implementation

- DONE: The deliverable is committed on the worktree branch REBASED ONTO the parent layer origin/spacedock-ensign/install-sh-edge-prerelease-parity (stacked mode - never main, which lacks SPACEDOCK_CHANNEL): the shared-core Binary-absent bullet and Channel sub-bullet land exactly per the approved before/after wording with NO growth beyond it (captain directive: keep it lean - the +1 line / ~2x-chars bullet is the approved ceiling, trim if possible), the fo-install-gate Channel selection section, and the install.md doc diff.
  `88cc3e794` on `spacedock-ensign/install-gate-channel-aware-hint`, parent `b00886d3a` (verified: remote main install.sh greps 0 `SPACEDOCK_CHANNEL`, the parent's greps 4). Lean cap MET and trimmed: bullet+sub-bullet 925 chars vs the approved 932-char ceiling (before 457, so +1 line at ~2.02x) — "same with `SPACEDOCK_CHANNEL=edge ` before `sh`" is shorter than the approved "that command with … prefixed", and "`local` or other OS" drops "channel". **One intentional deviation from the approved bytes, escalated to the FO before commit — see Summary.**
- DONE: The extraction-style tests are green offline (classifier table over the recorded fixture paths, contract-vs-doc curl byte binding, network-gated cask legs) and go test ./... plus -race are green; the two falsifying classifier edits from ideation red the table test.
  Four tests in `skills/integration/install_hint_channel_test.go`, all green; falsification run for each claim:
  · AC-2 `TestChannelClassifierTable` — 9 fixture paths from the real cache. Dropping the marketplace arm reds the `spacedock-edge/spacedock/0.27.0` row ("stable", want "edge"); dropping the suffix arm reds the `spacedock/spacedock/0.27.0-pre7-dev` row. Both observed, control restored green.
  · AC-3 `TestStableHintMatchesPublishedDoc` — reds if either the contract's or the doc's curl line is edited alone.
  · AC-1 shape guard `TestEdgeHintDeliversChannelToScript` — executes the published edge command with only the fetch stubbed; reds with `channel=UNSET` when the assignment is moved onto `curl`. Observed.
  · AC-3/AC-4 `TestContractCasks` — tap-resolved; logged `edge cask 0.27.0-pre8 vs stable cask 0.26.0, contract pin 0.27`. Reds if a cask is renamed or the edge cask stops tracking prereleases.
  **`go test ./...` and `-race` are NOT fully green:** `internal/cli TestCodexResolveManifestAgainstInstalledHost` fails, PROVEN PRE-EXISTING — it fails identically at parent `b00886d3a` with my changes stashed. It is machine-state, not diff-caused (a stray `spacedock-local` codex plugin in this host's cache). Everything else green, both lanes. `gofmt -l` clean.
- DONE: The AC-1 live captive-PATH edge install evidence is recorded in the stage report (observed --version Channel/minor lines pasted), and the surface lands within the approved +150 net / 4 files (tolerance +/-40, +/-1).
  Evidence recorded (below). The surface bound WAS BREACHED — final **+329 net / 4 files** vs the approved +150 ±40 (ceiling +190), 73% past it; files match (4, ±1 ok) and the entire overage is the test file at 313 lines vs the estimated 130. **`captain-ruling[2026-08-24]` accepted the breach as an estimate correction, superseding the bound:** +329 net / 4 files is now the correction-round baseline, and the overage is accepted test coverage. The same ruling ratified the FO-authorized sh-side wording fix. Item DONE on the superseding bound, not on the original one; the deviation reasoning is retained in the Summary and in `### Finding` as the durable record.

  AC-1 live captive-PATH run, `HOME` and `PATH` redirected to a temp dir, commands extracted verbatim from the shipped doc. **The published URL points at `main`, which does NOT yet carry `SPACEDOCK_CHANNEL`** (it ships in unmerged parent #756), so against the live URL BOTH arms land 0.26.0 — the hint is inert until #756 reaches main. Re-run against the parent's actual `install.sh` bytes (what main will hold post-merge):

      SPACEDOCK_CHANNEL=edge  ->  install.sh: resolved edge channel to v0.27.0-pre8
                                  spacedock 0.27.0-pre8
                                  Channel: edge (spacedock@spacedock-edge)
      stable (default)        ->  install.sh: resolved stable channel to v0.26.0
                                  spacedock 0.26.0 (contract 3)          [no Channel: line]

  Edge minor 0.27 >= pin 0.27 and `^Channel: edge` both satisfied; the stable baseline moves the wrong way (0.26 < 0.27) exactly as AC-1 requires.

### Summary

The channel classifier ships as designed and is proven by execution rather than by reading the contract: 9 recorded fixture paths, with each of ideation's two counterexamples confirmed to red the table when its arm is removed. The macOS arm resolves against the real tap (edge `0.27.0-pre8` vs stable `0.26.0` vs pin `0.27`), so both cask names in the contract are real.

**One deviation from the approved wording, and it is the load-bearing finding.** The approved bytes said the Linux edge hint is the stable command "with `SPACEDOCK_CHANNEL=edge ` prefixed", i.e. `SPACEDOCK_CHANNEL=edge curl … | sh`. In POSIX shell a variable prefix binds to the FIRST command of a pipeline, so the script runs with the variable UNSET and install.sh falls back to stable — reproducing precisely the defect this task exists to remove, and precisely AC-1's declared falsifying change. Measured, not reasoned: `channel=UNSET` vs `channel=edge`. I shipped the assignment on `sh`, which is byte-neutral, keeps the lean cap, and matches the form the parent layer already publishes at `install.md:31` (the string AC-3 binds against). Raised with the FO before committing; `TestEdgeHintDeliversChannelToScript` now guards it behaviorally, and `fo-install-gate.md` records why.

**Two items for the gate.** First, the surface is 73% past its tolerance ceiling (+328 vs +190). The ideation estimate priced only ~130 lines of test for what turned out to be four ACs plus the new shape guard; I consolidated the two cask tests into subtests and trimmed comment prose, which recovered only 14 lines, and cutting further means dropping AC coverage — I chose falsifiability over the line estimate rather than deciding that trade silently. Second, AC-1's value is gated on parent #756 merging to main: until then the hint text is correct but the script it points at ignores the variable. Stacked delivery resolves this by construction (the parent merges first), but validation must not expect a green live AC-1 against main before then.

Live-lane requirement for the PR ceremony: this diff touches `skills/first-officer/references/**`, so per the workflow's Proof policy ALL THREE host live lanes (`claude-live`, `codex-live`, `pi-live`) are REQUIRED green at merge — the shared core is host-neutral contract. I did not run them. A flake there is grounds to re-run to green, never to skip. The shipped contract is also one of the four high-stakes surfaces, so the detached read-only adversarial audit applies before merge.

### Finding: the approved edge-hint wording cannot satisfy AC-1

Recorded per `## Review-finding disposition`. Raised by the implementing ensign before any contract bytes were committed; candidate unchanged pending authorization.

**Four evidence fields.**

1. *Released user and normal workflow:* an edge-channel plugin user hits the binary-absent first-officer boot gate and runs the hinted Linux install command — this task's primary journey.
2. *Observable harm:* the run lands `spacedock 0.26.0` reporting `Channel: stable`; the next boot aborts against the 0.27 pin. The abort loop the task exists to remove survives unchanged.
3. *Affected value AC:* `value-ac[AC-1]` — "A binary-absent boot on an edge-channel plugin ends with an edge-parity binary, not an abort loop." The approved bytes are also AC-1's own declared falsifying change, so the approved wording reds the approved AC.
4. *Trigger evidence:* measured in POSIX shell, both arms —

       SPACEDOCK_CHANNEL=edge printf '…' | sh   ->  channel=[UNSET]
       printf '…' | SPACEDOCK_CHANNEL=edge sh   ->  channel=[edge]

   A variable prefix binds to the FIRST command of a pipeline, so prefixing `curl` never exports the variable to the `sh` that runs the script; `install.sh` then takes its `stable` default.

**Classification (worker proposal):** Material · owned by this task · disposition `fix`.

**FO authorization:** `fix`, authorized as proposed, 2026-08-24. The first officer reproduced the measurement independently (`channel=[UNSET]` on the prefix form, `channel=[edge]` on the sh-side form) and confirmed the parent layer already publishes the sh-side form at `install.md:31`, so the remedy is byte-neutral, satisfies AC-1, and additionally removes a contract-vs-doc divergence AC-3 would have caught. Captain's lean cap unaffected (925 chars vs the 932-char approved ceiling).

**Disposition applied.** The Linux edge hint ships with the assignment on `sh` in every shipped surface: the shared-core Binary-absent bullet, the `fo-install-gate.md` "Where the env var goes" rationale bullet, and the tests (which execute the published doc form rather than a copy). Verified no prefix-on-curl form survives in any shipped command — the single remaining textual occurrence is a test comment naming the defect being guarded against. `TestEdgeHintDeliversChannelToScript` reds with `channel=UNSET` if the assignment is moved back onto `curl`.

**Scope note.** The ideation-approved wording is superseded on this one point only. Design intent — hint the channel-correct command — is unchanged; the approved bytes were a means that failed the approved AC-1. No AC was narrowed and no approved surface boundary moved.

**Also authorized this round.** The `install.md` delta shrinking to +5 (the parent layer already covers the Binary tab; only the macOS Homebrew tab needed the edge line) — landing under estimate, no ruling required. Brew's tap-name normalization is now recorded in `fo-install-gate.md` (commit `84699ca23`): the contract spells the tap `spacedock-dev/homebrew-tap`, brew reports `spacedock-dev/tap`, both resolve, and the cask test encodes the mapping.

## Stage Report: validation

- DONE: Every AC verified with REPRODUCED evidence (run the tests/commands, not re-read the report): the classifier table with both arm-removal falsifications, the contract-vs-doc curl byte binding, the shape guard proving the published edge command delivers the variable to the script, the tap-resolved cask legs, and the AC-1 captive-PATH live run against the PARENT branch's install.sh bytes (never live main, which lacks SPACEDOCK_CHANNEL until #756 merges) - plus one semantic adversarial pass over the classifier (adjacent variants: the ten fixture paths, an adversarial home dir, a --plugin-dir checkout).
  · AC-2 `TestChannelClassifierTable` green at `84699ca23`; both arm-removal falsifications reproduced on the throwaway checkout: dropping the marketplace arm reds exactly `spacedock-edge/spacedock/0.27.0` ("stable", want "edge"); dropping the suffix arm reds exactly `spacedock/spacedock/0.27.0-pre7-dev`. Bytes restored, green again.
  · AC-3 `TestStableHintMatchesPublishedDoc` green; editing the doc's curl line alone reds it (the anchored extraction fails — the binding cannot pass on divergence).
  · AC-1 shape guard `TestEdgeHintDeliversChannelToScript` green; moving the assignment onto `curl` reds with `channel=UNSET` — the standing proof of the captain-ratified sh-side wording (ruling 2), so the divergence from the ideation before/after is not re-flagged.
  · AC-3/AC-4 `TestContractCasks` green tap-resolved: both hinted cask tokens publish in the tap; logged `edge cask 0.27.0-pre8 vs stable cask 0.26.0, contract pin 0.27`. Caveat: under this session's sandbox, brew hits bootsnap EPERM and the leg silently SKIPs; `HOMEBREW_NO_BOOTSNAP=1` resolves it.
  · AC-1 value run, captive `env -i HOME/PATH`, parent install.sh bytes (the candidate diff touches install.sh on 0 lines, so HEAD carries them): edge arm resolves `v0.27.0-pre8` and the landed binary reports `spacedock 0.27.0-pre8` / `Channel: edge (spacedock@spacedock-edge)` — minor 27 ≥ pin 27 and `^Channel: edge` both hold; the stable baseline lands `spacedock 0.26.0` with no Channel line — below the pin, so a channel-blind regression is observable. Live main greps 0 `SPACEDOCK_CHANNEL`: AC-1 stays inert against main until #756 merges (recorded gating condition, ruling 4 — not a red).
  · Adversarial pass: the shipped table is 9 rows (both counterexamples, the adversarial `-pre` home, two source-checkout rows — the ten-vs-nine count noted under polish) plus 8 extra variants probed live: empty base, base without the skills suffix, v-prefixed version dir, and a checkout directory literally named `spacedock-edge` all classify `local` (the safe no-install arm); a trailing slash on a real cache path classifies `local`; a checkout nested under a `spacedock-edge/` parent classifies `edge` (deferred risk 3).
- DONE: The detached adversarial audit applies (shipped contract is a high-stakes surface): run it on a throwaway checkout at 84699ca23, never the implementation worktree, and note the result in the reviewer-findings block.
  Fresh clone in the session scratchpad, checked out at `84699ca23`; all falsifying mutations ran there and were reverted (final status clean). The implementation worktree was never mutated. Audit result: no material finding — reviewer-findings block below.
- DONE: A PASSED/REJECTED recommendation with findings classified on both axes (outcome vs evidence defect; material vs deferred risk), deferred risks listed separately with promote-conditions, and only material findings blocking.
  PASSED. Zero material findings; three deferred risks and two polish notes below.

### Reviewer findings (validation)

Detached audit: clean of material findings. Full suite at `84699ca23`: `go test ./...` and `-race` green except `internal/cli TestCodexResolveManifestAgainstInstalledHost`, reproduced byte-identical at parent `b00886d3a` with the candidate absent — machine state (a stray `spacedock-local` codex cache entry), per ruling 5. `gofmt` clean. Surface `git diff --numstat b00886d3a..HEAD` = 330 ins / 1 del = +329 net, 4 files — exactly the ratified correction-round baseline (ruling 1, not re-flagged). Lean cap: bullet + Channel sub-bullet 925 chars with newlines (923 without) ≤ the 932 ceiling — no growth (ruling 3).

Deferred risks (all evidence-defect axis; none blocks):

1. The contract's Linux-edge clause has no test binding: the shape guard executes the doc's edge line and AC-3 binds the stable command only, so a future contract-side rewording back to the prefix form would red nothing. Current bytes are correct ("before `sh`"). Promote to material: any edit to the Binary-absent bullet's Linux-edge clause without adding a contract-side extraction leg.
2. The cask legs skip (never fail) when brew or the tap is unavailable — by the approved test plan, to keep the offline lane deterministic — so AC-4 has no standing automated proof on a machine without the tap; observed concretely as the silent sandbox-EPERM skip above. Promote to material: a lane inventory showing no CI lane resolves the tap while the edge cask version can drift below the pin.
3. Classifier misroutes on two unobserved layouts: a source checkout nested under a parent directory literally named `spacedock-edge/` classifies `edge` (a package hint against a human-edited tree), and a trailing slash on the loader-supplied base classifies `local`. Every other probed failure lands on `local`, the no-install arm. Promote to material: a real loader-supplied base observed with either shape.

Polish (no user-visible loss): the shipped fixture table is 9 rows against ideation's "ten fixture paths" claim — all load-bearing rows present and AC-2 pins no count. The doc's macOS edge line hints bare `spacedock@next` where the approved diff spelled the fully-qualified token — it resolves in-context directly under the tab's `brew tap` line and matches the tab's existing bare style; no test binds it.

### Summary

PASSED. All four ACs reproduce by execution at `84699ca23` on the #756 stack: the classifier table reds on exactly the documented counterexample row when either arm is removed, the shape guard reds with `channel=UNSET` on the prefix form, the tap-resolved cask legs log edge `0.27.0-pre8` vs stable `0.26.0` vs pin `0.27`, and the captive-PATH run against the parent's install.sh bytes lands an edge-parity `0.27.0-pre8` while the stable baseline lands `0.26.0` — the wrong way, as AC-1 requires. The one full-suite red is the reproduced pre-existing machine-state failure named in ruling 5. No material findings; three deferred risks with promote conditions and two polish notes are recorded above for FO disposition. Merge ceremony still owes the PR, the stack link, and all three host live lanes green — not run here, per the dispatch boundary.

### Addendum: post-approval delta-check of the captain-ordered sweep (d4e1c5cbc..9583a4dfb)

Scope: sweep commits `f6542e088` + `9583a4dfb` on the branch restacked onto newer main (parent now `035934968`); every SHA above is stale by rebase, expected. All verified by execution at `9583a4dfb`:

1. Generalization holds: grep over `fo-install-gate.md` and `first-officer-shared-core.md` for literal versions, cache paths, task identifiers, and issue refs finds none in the swept sections — the only hits are hyphenated English words in untouched pre-existing prose ("needs-preparation", "context-pressure", "wording-present", the sandbox paragraph), plus the deliberate `minor 0.27` release stamp, excluded by design.
2. Full unfiltered suites green: `skills/integration` (all tests; `TestContractCasks` tap-resolved, names + edge-satisfies-pin), `internal/contractlint` including `TestInstallHintNoDrift` — the two-arm contract-vs-doc binding (curl token equality AND brew tap+formula tokens) that the implementer's filtered run missed — and `internal/release`. Tap short form `spacedock-dev/tap` is consistent across contract, doc, and test, and the pair is now mechanically enforced by the drift lint.
3. AC-3 ruling: its tested binding — the Linux stable curl byte-equality (`TestStableHintMatchesPublishedDoc` plus contractlint arm 1) — is untouched by the sweep and green; both hinted cask tokens resolve in the real tap under the short spelling. The macOS stable hint's spelling change (`homebrew-tap` → `tap`) is the captain-ordered standardization superseding AC-3's "bytes unchanged" headline on that token only; AC-3's substance holds.
4. Lean cap: bullet + Channel sub-bullet now 898 chars incl. newlines ≤ 932.
5. Superseded figures: surface `git diff --numstat 035934968..HEAD` = 330 ins / 2 del = +328 net, 4 files (was +329/925-char bullet; now +328/898). Deferred risk 1's promote-condition NOT triggered: the Linux-edge clause is byte-identical between `d4e1c5cbc` and `9583a4dfb` (verified by extracting the clause from both blobs). The polish note on the doc's bare `spacedock@next` is resolved — the sweep ships the fully-qualified `spacedock-dev/tap/spacedock@next`.

PASSED stands. No new findings; deferred risks 1-3 carry over unchanged (risk 2's tap-skip caveat now reads `spacedock-dev/tap` in the skip messages, same behavior).

### Addendum: round-3 delta-check of the fo-install.md relocation (9583a4dfb..0fea3d5ec)

Scope: commit `0fea3d5ec` — `fo-install-gate.md` renamed to `fo-install.md`, the classifier and per-OS command table moved out of the boot-resident core into it, the load rule made sandbox-independent with an explicit two-branch sandbox arm. All verified by execution at `0fea3d5ec`:

1. Repointed tests genuinely bind the new file — re-falsified independently on the throwaway checkout, not re-run: dropping the suffix arm in `fo-install.md` reds `TestChannelClassifierTable` on the `0.27.0-pre7-dev` row; a byte-changing edit to `fo-install.md`'s stable curl line (`sh` → `bash`) reds BOTH `TestStableHintMatchesPublishedDoc` (which now cites fo-install.md in its extraction failure) and `TestInstallHintNoDrift`; controls restored green each time. Evidence note, non-blocking: a suffix-APPEND edit (`| sh` → `| sh -s -- --quiet`) reds only `TestStableHintMatchesPublishedDoc` — the drift lint's arm 1 is substring containment and the doc's command survives as a prefix — so the "reds both" claim holds for byte-changing edits only; the pair stays fully bound either way because the anchored byte-equality test alone closes every mutation class tried.
2. Old-name closure: `grep -r fo-install-gate` over the tree hits ONLY the frozen live recording `internal/ensigncycle/testdata/claude_live_auto_continue_run31915540750_sonnet.stream.jsonl` — historical transcript bytes, referenced by no Go code as a path; `TestBootResidentDeferredLoadPointsResolve` and `TestUserSkillReferenceClosureResolves` pass, which would fail on a dangling contract pointer.
3. Boot-resident bullet: 232 chars incl. newline (232 reported, ≤ 932 ceiling); the shared core greps 0 for `curl|brew|case "$R"|SPACEDOCK_CHANNEL` — no install command or one-liner remains; the `require binary minor 0.27` pin did NOT move (1 hit in shared core, 0 in fo-install.md).
4. Sandbox arm shape verified by read: the load rule is sandbox-independent ("reading this file is always safe; only running an install inside a sandbox is forbidden"); inside = classify, select the channel-correct row, print, defer to the human — no install offer, no sentinel touch ("Skip `## Install-and-resume offer` entirely; touch no sentinel"); outside = the existing offer/sentinel/converge machinery, which itself opens with the outside-only guard.
5. Full unfiltered `go test ./...` at `0fea3d5ec`: green except `internal/cli TestCodexResolveManifestAgainstInstalledHost`, reproduced byte-identical at the current parent `035934968` — the same machine-local codex-cache state, not this diff.
6. Surface: `git diff --numstat 035934968..0fea3d5ec` = 401 ins / 42 del = **+359 net, 7 files** vs the accepted +328/4. Two of the seven files are the rename's delete+add split (`fo-install-gate.md` −22 / `fo-install.md` +52 — the same contract, relocated); the real growth is the captain-ordered explicit sandbox arm and relocation framing in `fo-install.md` plus test growth (`install_hint_channel_test.go` 313 → 322 and the two contractlint repoints). Read as a captain-ordered correction on the ratified baseline, not implementer drift.

PASSED stands. No new material findings; the substring-containment nuance in `TestInstallHintNoDrift` arm 1 is recorded as an evidence note (the anchored pair-binding test covers it); deferred risks 1-3 carry over. Risk 1's subject is now `fo-install.md`'s `- **Linux, edge:**` row and it carries over UNCHANGED in substance — probed by execution: reverting that row to the broken prefix form (`SPACEDOCK_CHANNEL=edge curl … | sh`) reds nothing in skills/integration or internal/contractlint (the shape guard runs install.md's line; the byte-equality leg binds the stable row only). Current bytes are correct and the "Where the env var goes" rationale sits directly under the table. Promote to material: any edit to that row without adding a contract-side shape or equality leg.

### Addendum: round-4 delta-check of the captain PR-review cuts (0fea3d5ec..ef43c7ca8)

Scope: commit `ef43c7ca8` — fo-install.md cut 52→32 lines, classifier rewritten in the readable basename/dirname form, `local` arm struck by captain ruling (source checkouts classify stable), rationale prose deleted, sentinel reduced to one fixed path, one-sentence sandbox arm, `TestVersionGateSandboxRegistry` narrowed to the core alone, the binary-present `^Sandbox:` corroboration rule dropped. All verified by execution at `ef43c7ca8`:

1. Verdict identity: the old and new classifier forms agree on **29/29 real plugin-cache paths** (every `…/cache/<marketplace>/<plugin>/<version>` dir on this machine — a superset of the implementer's 16). One observed row worth a future fixture: `spacedock/spacedock-edge/0.19.9` classifies `edge` via the marketplace arm matching the plugin-name segment — identical in both forms, pre-existing semantics, and plausibly the true channel for that legacy naming. The re-expected source-checkout rows are load-bearing: widening the version arm `[0-9]*-*` to `*-*` reds exactly the two source-checkout rows ("edge", want "stable"); restored green.
2. All four falsification legs re-red on the right assertion and controls restore green: marketplace-arm drop reds the unsuffixed-edge row; suffix-arm drop reds the `0.27.0-pre7-dev` row; renaming `APP_SANDBOX_CONTAINER_ID` in the shared core reds the narrowed `TestVersionGateSandboxRegistry` ("core prose does not check the sandbox env var"); a byte edit to fo-install.md's stable curl line reds BOTH `TestStableHintMatchesPublishedDoc` and `TestInstallHintNoDrift`.
3. fo-install.md is 32 lines: no load narration, no why-bullets, no sentinel key derivation (one fixed path `${TMPDIR:-/tmp}/spacedock-install-attempted`), a one-sentence sandbox arm; the one-attempt-ever and touch-sentinel-BEFORE-the-install rules are both still stated as rules (offer steps 1 and 3).
4. Nothing dangles on the dropped corroboration rule: `grep -rn '\^Sandbox|corroborat'` over skills/, internal/contractlint/, docs/site/ returns zero live references (the FO surfaced the rule's homelessness to the captain separately; not re-flagged here).
5. Full unfiltered `go test ./...` and `-race` at `ef43c7ca8`: green except the machine-local `internal/cli TestCodexResolveManifestAgainstInstalledHost` in each lane — same bytes as the round-3 reproduction at the unchanged parent `035934968`.
6. Surface: `git diff --numstat 035934968..ef43c7ca8` = 395 ins / 51 del = **+344 net, 7 files** (rename delete+add split still counted as two). The entity's AC-2 restatement (`edge`/`stable` only, classifier relocated to fo-install.md, `local` struck by captain ruling) is consistent with the shipped classifier and the updated table rows.

Deferred-risk movement: risk 3's trailing-slash half is FIXED by the basename/dirname rewrite (a trailing-slash cache path now classifies stable, probed); its nested-`spacedock-edge/`-parent half carries over with reduced consequence (source checkouts now legitimately receive package hints — the residual error is channel choice only). Risk 1 re-confirmed unchanged at this HEAD: the prefix-form regression on fo-install.md's edge row still reds nothing. Risk 2 unchanged.

PASSED stands.

### Captain-ordered post-approval correction: de-specify the contract, standardize the tap spelling

Ordered by the captain in chat 2026-08-24, after the implementation gate approval, scoped to the contract text this layer added. Principle applied: the contract carries durable rules only — no specific versions, no task identifiers, no GitHub issue refs, no point-in-time observations. This layer's added text was the only offender in the contract. Delivered as an ordinary commit on top; no force-push, no parent rebase.

- **`fo-install-gate.md` "Why two signals"** — dropped the literal cache paths and release numbers; the bullet now states the shape (a dev build can be installed under the stable marketplace name while carrying a prerelease-suffixed version; an edge-marketplace install can carry an unsuffixed version). The literal observed paths remain where they belong: the test fixtures and this entity body.
- **`fo-install-gate.md` "Never widen the match"** — dropped the specific home directory; now names the property (a home directory whose name contains a prerelease-like token).
- **Tap spelling** — standardized on brew's canonical short form `spacedock-dev/tap` in the shared-core Binary-absent bullet, `fo-install-gate.md`, and this layer's `install.md` macOS addition. The repo is NOT renamed: brew requires the `homebrew-` prefix on the repository, and the short form is the derived canonical name brew itself reports. The "Cask names" both-spellings note was DELETED — with one spelling it has no job.

**Verification.** Short form confirmed against live brew before shipping: `brew tap spacedock-dev/tap` exits 0 and `brew tap-info --json spacedock-dev/tap` resolves both cask tokens as `spacedock-dev/tap/spacedock` and `spacedock-dev/tap/spacedock@next`. `TestContractCasks` now binds the contract's short-form tokens directly against brew's, with no normalization step in between — a cleaner binding than before; re-falsified by pointing the contract at a nonexistent `spacedock@edge` cask, which reds with "the tap publishes no cask". All four tests plus `internal/contractlint` (which reads `fo-install-gate.md`) green.

**Lean cap.** The short form SHRANK the boot-resident bullet: **898 chars, down from 925**, against the 932 approved ceiling — now 34 under. Sweep is net −1 line, so the surface moves to **+328 net / 4 files** (from +329, within the captain-accepted correction-round baseline).

**AC-3 note.** The macOS stable hint bytes changed (`spacedock-dev/homebrew-tap/spacedock` → `spacedock-dev/tap/spacedock`), which touches AC-3's "stable installs' hint bytes are unchanged" clause. AC-3's substance is intact: the cask NAMES it binds (`spacedock`, `spacedock@next`) are unchanged, and its tested binding is the Linux curl command, which this sweep does not touch. Flagged rather than assumed harmless — the captain's post-approval order supersedes the clause on this point.

**Scoping call corrected mid-sweep — the doc's tap line had to move too.** I first left `docs/site/get-started/install.md:9` (`brew tap spacedock-dev/homebrew-tap`) alone as pre-existing and not this layer's addition, and recorded the resulting mixed spelling as cosmetic. That was wrong. `internal/contractlint` `TestInstallHintNoDrift` binds the FO version-gate prose's `brew tap` and `brew install` tokens to the ones install.md's Homebrew tab publishes, so standardizing the contract without moving the doc broke a real, enforced pair. Fixed in `9583a4dfb`: the doc's tap line now reads `brew tap spacedock-dev/tap`, and no `homebrew-tap` spelling remains anywhere in `skills/` or `docs/site/`.

**Process note, recorded because it nearly shipped.** The breakage was masked by my own test run: after the sweep I ran `internal/contractlint` with a `-run` filter that excluded `TestInstallHintNoDrift`, saw green, and only the unfiltered full-suite run caught it. The lesson is the obvious one — after a contract-text change, run the contract-lint package unfiltered, since its whole job is binding prose to other files.

**Rebase note.** The stack was rebased onto a newer `main` (now carrying the #754 merge) by the stack machinery between the gate approval and this sweep. Commit SHAs recorded earlier in this report — and in the validation report — are stale. Current layer state, parent `035934968`:

    149932b42  implementation
    d4e1c5cbc  tap-normalization note (superseded: the note itself is deleted by the sweep)
    f6542e088  the sweep
    9583a4dfb  doc tap line, completing the sweep

Local and remote were in sync throughout, so both sweep commits pushed as ordinary fast-forwards. No force-push, no parent rebase.

**Relationship to the validation PASS.** Validation passed at pre-rebase `84699ca23` (= post-rebase `d4e1c5cbc`), which PRE-DATES this sweep; the captain ordered the sweep afterward. Three consequences for anyone re-reading that report:

- Its recorded surface (+329 net) and lean-cap figure (925 chars) are superseded by **+328 net / 4 files** and **898 chars**. Both move in the safe direction.
- Its polish note that "the doc's macOS edge line hints bare `spacedock@next` where the approved diff spelled the fully-qualified token" is RESOLVED by the sweep — that line now reads `brew install spacedock-dev/tap/spacedock@next`.
- Its **deferred risk 1** promote-condition is NOT triggered. That condition is an edit to the Binary-absent bullet's *Linux-edge* clause; the sweep edited only the *macOS* clause's tap spelling and left `Linux edge: same with SPACEDOCK_CHANNEL=edge before sh` byte-identical. The risk stands open, unchanged, still unpromoted.

Deferred risks 2 and 3 are untouched by the sweep. Re-validation, if the FO wants one, is a re-run of the same four tests plus unfiltered `internal/contractlint` at `9583a4dfb`; the full suite there is green except the same pre-existing `TestCodexResolveManifestAgainstInstalledHost` machine-state failure.

### Captain-ordered design correction: the channel machinery moves out of the boot-resident core

Ordered by the captain in chat 2026-08-24. This **reverses the placement point the ideation gate approved**, on captain authority. It supersedes the ideation rationale under `## Proposed approach` ("Why the rule lives in the shared core, not only in `fo-install-gate.md`") and the paragraph in the deferred file that argued the same case — that paragraph is DELETED, because it argues for the design the captain rejected.

**The premise was wrong.** The ideation argument was: the deferred file is not loaded inside a sandbox, yet the sandbox arm must still print an exact install command, therefore the classifier must be boot-resident. The captain's correction: "never install inside a sandbox" never meant "never READ the reference inside a sandbox." Reading is harmless. Once the file loads in both cases, the reason to keep the machinery in the boot-resident core evaporates.

**Three changes.**

- **Renamed** `references/fo-install-gate.md` → `references/fo-install.md` (captain: do not call it a gate). Every live reference updated: shared-core lines 10 and 48, `internal/contractlint/version_gate_smoke_test.go`, and the reference-closure check in `boot_resident_closure_test.go`, which resolves `references/…md` mentions against the filesystem and would have caught a miss.
- **Moved** the per-OS/per-channel command table AND the classifier out of the shared core into `fo-install.md`, which now owns channel classification, the install commands, the sandbox arm, and the install-and-resume offer.
- **Changed the load rule:** `fo-install.md` loads at the binary-absent trigger REGARDLESS of sandbox. Shared-core line 48 drops "outside a sandbox"; line 10 drops the "Outside sandbox read …" phrasing. `fo-install.md` gained an explicit two-branch sandbox arm — inside: classify, print the exact channel-correct command, tell the human to run it outside, never offer or run an install, never touch the sentinel; outside: the existing install-and-resume machinery. The shared core KEEPS its never-offer/run-installation-inside sentence: that is behavior, not load policy.

**Lean cap, now trivially met.** The Binary-absent bullet collapses from **898 chars to 232** against the 932 ceiling — roughly a quarter of the budget, and the boot-resident core no longer carries any install command or shell one-liner at all.

**Tests repointed and re-falsified.** The classifier and stable-curl extraction in `skills/integration/install_hint_channel_test.go`, and the brew/curl token binding in `internal/contractlint/install_hint_drift_test.go`, now read `fo-install.md`. Repointing a test risks it silently binding nothing, so each was re-falsified against the new file: dropping the marketplace arm reds the classifier table on exactly the `spacedock-edge/spacedock/0.27.0` row; editing the stable curl command in `fo-install.md` alone reds BOTH `TestStableHintMatchesPublishedDoc` and `TestInstallHintNoDrift`. Controls restored green. `contractPin` still reads the shared core, correctly — the `require binary minor` declaration did not move.

**AC citations updated, AC properties unchanged.** AC-2 and AC-3 named the shared core as the extraction source; they now name `references/fo-install.md`. Only the mechanism's location moved — the properties each AC asserts are untouched, so the validation PASS's reproduction of them still stands on substance.

**One deliberate non-change.** `internal/ensigncycle/testdata/claude_live_auto_continue_run31915540750_sonnet.stream.jsonl` still contains the old `fo-install-gate.md` name. It is a frozen recording of a past live session's tool output; rewriting it would falsify historical evidence. Left as-is by intent, not by omission.

### Captain PR review on #757: the reference was overbuilt

Six inline comments on `fo-install.md`, 2026-08-24. Ruling: **the file is overbuilt.** The principle for this round — the reference carries RULES and COMMANDS; rationale, narration, and admitted-tradeoff prose belong in this entity body and in the tests, which already hold them. The captain explicitly reversed the earlier FO direction to "note" and "record" things in-file, which is where part of the bloat came from.

**Result: 52 lines to 32, and it reads in well under a minute.** The shared-core bullet is untouched at 232 chars.

- **Load narration deleted.** The shared core owns the trigger; the file opens on its title.
- **Classifier made readable.** `${R##*/}` indirection replaced with named commands: `V=$(basename "$(dirname "$(dirname "$B")")")`. Still one extractable, executable line. Verified verdict-identical to the previous form on all 16 real cache paths plus the adversarial home.
- **The `local` arm is struck** — captain: a source-checkout skill base with no binary is not a real class. The classifier returns `edge`/`stable` only; the `local` bullet and the `local` install row are gone; the other-OS `go build` hint stays. Because the version arm requires a digit-leading segment, a source checkout falls through to `stable`, so the two source-checkout fixture rows re-expect `stable` rather than being dropped — the behavior stays pinned instead of going untested.
- **"Why two signals" and "Never widen the match" deleted** as editor-rationale, not rule. The two-signal behavior remains pinned by the table test's counterexample rows, and the no-widening behavior by the adversarial `/Users/pre-release-tester/…` row — the contract states the rule and the tests hold the argument.
- **Sentinel dumbed down** to one fixed path, `${TMPDIR:-/tmp}/spacedock-install-attempted`. The per-runtime session-key derivation, the pi working-directory hash fallback, and the over-reach admission prose are all gone. Retained properties: one attempt ever, `touch` before running, and sentinel-present ⇒ skip the offer and print the manual command plus `rm <path>`.
- **Sandbox arm reduced to one sentence,** and the `## Fallback-message grammar` section folded into the sentinel-present bullet (step 5 now points back at step 1 rather than restating it).

**Two consequences that required changing a test I do not own — flagged, not done silently.** `internal/contractlint` `TestVersionGateSandboxRegistry` asserted the sandbox registry env name+value in BOTH the core and this file, and separately required this file to carry the `^Sandbox: ` corroboration line. The ordered cuts delete exactly that content.

1. *Registry mirror.* Narrowed to the CORE ALONE. This is principled under the new split: the core PERFORMS the sandbox check in every gate class, `fo-install.md` only consumes the verdict, so the duplicate was a second place to forget to update. Re-falsified: renaming the env var in the core reds the test.
2. *`^Sandbox: ` corroboration.* **This rule now has no home in the repo** — the core does not carry it (checked: zero occurrences). It governs binary-PRESENT classes, where `--version`'s `^Sandbox: ` line corroborates the env verdict and the env check wins on disagreement. Deleting it from a binary-ABSENT reference is correct scoping, but the rule itself was dropped rather than relocated. **Recommend the captain decide** whether it moves to the core's binary-present/wrong-version arm; it is out of this task's scope and I did not invent a home for it.

**Falsification after the cuts** — every changed test re-red, then restored: dropping the marketplace arm reds the `spacedock-edge/spacedock/0.27.0` row; dropping the prerelease-suffix arm reds the `spacedock/spacedock/0.27.0-pre7-dev` row; renaming the core's sandbox env var reds the registry test; editing the stable curl command in `fo-install.md` alone reds both the doc-binding test and `TestInstallHintNoDrift`.
