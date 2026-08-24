---
id: p1tgy61tbhxj9apvpswqbhcy
title: The skill-first install hint is channel-blind - an edge plugin tells the user to install the stable binary its own version gate then rejects
status: ideation
source: "Captain CL in chat, 2026-08-24, reviewing install-sh-edge-prerelease-parity (#756): 'does it work for the skill-first journey that tells the user to run it?' - it does not; follow-up scoped to the skill hint path #756 could not reach"
started: 2026-08-24T19:20:28Z
completed:
verdict:
score:
worktree:
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

**AC-2 — The shipped classifier maps every observed real install-path shape to its true channel.**
Verified by a table test that extracts the classifier one-liner from the shared core and executes it (never greps it) against fixture paths recorded from the real plugin cache and `installed_plugins.json`, plus a `--plugin-dir` checkout and an adversarial `/Users/pre-release-tester/…` home. Expected values come from the observed install layout, not from the file under test. Falsifying changes, both confirmed to produce wrong answers on real paths: removing the marketplace-segment arm calls `cache/spacedock-edge/spacedock/0.27.0` stable; removing the prerelease-suffix arm calls `cache/spacedock/spacedock/0.27.0-pre7-dev` stable.

**AC-3 — Stable installs' hint bytes are unchanged, and the contract's stable command cannot drift from the published one.**
Verified by binding two independently maintained values: the Linux stable command extracted from the shared core must equal the command published in `docs/site/get-started/install.md`, and the two cask names the contract hints (`spacedock`, `spacedock@next`) must both resolve in `spacedock-dev/homebrew-tap`. Falsifying change: edit either the contract's or the doc's curl line alone, or rename a cask — the pair diverges and the test fails.

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
