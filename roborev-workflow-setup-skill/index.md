---
id: 008h7wr55c7fn5x3r2wk26yz
title: Ship the Roborev workflow setup skill
status: ideation
source: Captain decision after spacedock-subspace Roborev adoption pilot, 2026-07-13
started: 2026-07-13T16:11:06Z
completed:
verdict:
score:
worktree:
issue:
milestone: 0.26.0
---

Ship Roborev adoption as a first-party, user-invocable `spacedock:roborev-setup` skill in the main Spacedock plugin. The skill helps a user decide whether Roborev fits a code workflow and configures the Spacedock workflow boundary without making Roborev a dependency of ordinary Spacedock use.

## Reference artifacts

- [Approved Spacedock–Roborev integration proposal](artifacts/roborev-spacedock-integration-proposal.md)
- [Canonical `roborev-setup` skill draft](artifacts/roborev-setup-skill/SKILL.md)
- [Draft skill UI metadata](artifacts/roborev-setup-skill/agents/openai.yaml)

The artifacts preserve the reviewed proposal and the user feedback folded into it. Implementation copies the canonical draft into `skills/roborev-setup/` and evolves it only with behavior-equivalent review fixes; `/tmp/spacedock-roborev-setup` is not a source of truth.

## Problem

Spacedock users can benefit from independent Roborev code-review evidence, but adoption currently requires knowing how to combine Roborev panels, implementation-exit ownership, fresh validation, split state checkouts, daemon placement, and Safehouse access. A separate integration plugin would make the setup entry point harder to discover and would create a one-skill packaging boundary without independent runtime code or release needs.

The pilot exposed two concrete guidance failures:

1. The first skill draft led with three web-documentation links. An actual setup run went searching for guidance even though `roborev quickstart` already emitted a version-matched **Current state** and **Configuration playbook**. This was an avoidable first-contact detour.
2. The generated `review_guidelines` copied repository procedures and required test commands. Inspection of `kenn-io/roborev`'s own configuration showed the narrower useful pattern: calibrate reviewer judgment with context it cannot infer from the diff or repository, such as trust boundaries, intentional compatibility posture, false-positive suppressions, and review-focus boundaries.

Skill behavior is harder to test than command behavior. Repository policy bans prose-grep as behavioral proof, requires a live drive for a skill change, and requires a detached adversarial audit for shipped contract or scaffolding changes. The design must produce fixture-backed durable evidence rather than treating skill wording or transcript phrasing as proof.

## Mechanism spike

The riskiest external assumptions were exercised during ideation against installed `roborev v0.62.0`:

- `roborev quickstart` exited 0 from this repository and printed both Current state and Configuration playbook, including actionable missing-state commands.
- `roborev review --help` exposes `--repo`, `--branch`, `--base`, `--panel`, `--min-severity`, and `--wait`.
- `roborev wait --help` accepts a commit/ref and exits 1 for a failed verdict or missing job, allowing an exact-tip wait without enqueueing another job; `show <head> --json` must distinguish Low-only findings from blocking findings instead of using exit status as the severity oracle.
- `roborev show --help` exposes `--job` and `--json` for durable synthesis-parent evidence.
- The directly opened official panel page documents `[review] hook_review_panel`, named `[review.panels.*]`, and synthesis-parent behavior; the configuration page documents checkout-local `excluded_branches`.

No product helper spike is needed. The setup mutations are ordinary file edits and Git operations already exercised by supported hosts. The live harness needs test-only recording shims, not a new `spacedock roborev` command.

## Proposed approach

- Add `skills/roborev-setup/SKILL.md` to the main `spacedock` plugin with `name: roborev-setup` and `user-invocable: true`, exposed as `spacedock:roborev-setup`. Add `skills/roborev-setup/agents/openai.yaml` from the canonical artifact.
- Keep the skill setup-only. Do not load it from the first-officer or ensign core and do not make Roborev mandatory for normal workflows.
- After locating the canonical code root and workflow, run `roborev quickstart` as the first Roborev command. Treat its Current state and Configuration playbook as primary. Use local `roborev <command> --help` for syntax gaps. Open only a directly relevant official page after naming a remaining configuration gap; do not search or browse preemptively.
- Explain which quickstart items matter and which are intentionally left missing. In particular, do not install or invoke Roborev fix/refine skills or the agent hook in a Spacedock-managed workflow.
- Before writes, show the integration boundary, proposed files, and any proposed `review_guidelines` calibration to the user. Apply only approved changes, then rerun `roborev quickstart` to verify current state.
- Configure `quick` as the cheap post-commit cost gate, the exact-head `code_completion` panel as authoritative implementation-exit evidence, and fresh validation as the consumer of stored synthesis-parent evidence before independent behavioral validation.
- Encode panel order in generated implementation-stage text or a dispatch checklist: after the final candidate commit, wait for that exact tip's existing `quick` review; fix and recommit any Medium-or-higher finding; let Low findings remain advisory; launch `code_completion` only after `quick` clears. The panels never start independently.
- Add the split state checkout's actual branch to its checkout-local `excluded_branches` without enqueueing a probe review.
- Detect Safehouse or another sandbox. When present, advise the smallest read-only runtime access needed for the external daemon; omit sandbox advice otherwise. Never commit machine-local permissions or expose all of `~/.roborev`.
- Keep `review_guidelines` to context the reviewer cannot infer: trust boundaries, intentional compatibility posture, false-positive suppressions, and review-focus boundaries. Never duplicate `AGENTS.md`, workflow or stage instructions, component procedures, or developer/test commands.
- Put the pilot's semantic adversarial pass in generated implementation-stage instructions, not `review_guidelines`: trace changes through representations/lifecycle phases, matrix adjacent variants, validate full records atomically, inspect scale and limits, and seek cases where tests pass despite wrong observable behavior.

## Generated workflow contract

The implementation stage receives this observable sequence:

1. Before the final candidate commit, run the implementation-owned semantic adversarial pass: trace identity/cardinality/order/bytes/attribution/authority/terminal state through every representation and lifecycle phase; matrix empty/terminal, repeated/out-of-order, every input path, Unicode/EOF/size/visibility/layout; use canonical validators or atomic full-record validation; inspect hot paths/readers for multiplicative work, blocking I/O, unbounded allocation, and implicit limits; and add exact-result plus failure/cleanup/scaling/over-limit proof where relevant.
2. Resolve and record candidate `HEAD` after the final implementation commit.
3. Run `roborev wait <head>` for that exact commit's already-enqueued `quick` job, then inspect `roborev show <head> --json` and verify its panel and head. Do not start `code_completion` before both operations complete.
4. On missing/error coverage or a Medium-or-higher stored finding, stay in implementation. Repair coverage or fix, commit, and wait on the replacement tip. Because `wait` may exit nonzero for a Low-only verdict, stored finding severities—not exit status alone—allow the Low-only result to continue without another commit.
5. Once quick clears, run the complete branch through `roborev review --repo <root> --branch=<branch> --base <base> --panel code_completion --min-severity medium --wait` and preserve the synthesis parent ID even on nonzero exit.
6. Fetch `roborev show --job <parent-id> --json`; require the current head and frozen merge-base-to-head range, the `code_completion` panel, exactly one successful execution per required member, and a passing synthesis parent.
7. Any code-changing commit invalidates both quick and `code_completion` evidence and restarts at step 1.

Roborev supplies evidence only. Implementation owns fixes and commits. The First Officer owns entity state and routing. Fresh validation verifies the stored exact-head parent, then independently exercises acceptance behavior; it does not rerun an unchanged passing panel.

## Review-guideline calibration contract

The skill treats `review_guidelines` as a reviewed setup proposal, not a place to mirror repository instructions. A proposed block may include only:

- a trust boundary the reviewer cannot infer from code;
- an intentional compatibility posture that changes how a diff should be judged;
- a known false-positive suppression with its reason;
- an explicit boundary on what the review should evaluate.

It must exclude text copied from `AGENTS.md`, workflow/stage procedures, component runbooks, command sequences, and required test commands. The user sees and approves the proposed calibration before `.roborev.toml` changes. The semantic adversarial pass is intentionally stored with implementer/process guidance instead.

## Out of scope

- A general `spacedock-integrations` plugin or dedicated Roborev plugin.
- A new binary setup command, independently versioned adapter, lifecycle mod, or stage-exit hook.
- Roborev daemon installation, service management, or agent authentication without explicit user action.
- Roborev `fix`, `refine`, agent-hook, commit, or routing ownership.
- A duplicated Roborev schema or a general web-documentation navigator.
- Tests that claim behavioral coverage by matching prose in `SKILL.md` or assistant transcript wording.

## Acceptance criteria

**AC-1 (VALUE) - With a complete quickstart fixture, setup reaches an approved, review-ready workflow with `roborev quickstart` as the first Roborev command and zero web-search or documentation-fetch events.**
Verified by: a network-denied live Codex drive using the installed local plugin, a recording fake `roborev`, and the host's structured tool-call log. Assertions require first Roborev argv `quickstart`, no search/browser/HTTP/curl/wget event, exit 0, approved workflow/config changes, state-checkout commit history, and clean status. The independent baseline is the prior draft's observed pre-setup documentation detour; any external-doc event moves the measured count from 0 and fails.

**AC-2 - When quickstart leaves a named panel-schema gap, setup consults local help next and any web fallback is a direct official page opened only after those local events, never a search.**
Verified by: a second live fixture mode whose quickstart output deliberately omits panel schema. The fake-command journal and structured tool log must order `quickstart` before `review --help`, then allowlist only the direct Roborev panel URL; search events, unrelated pages, or docs-before-help fail.

**AC-3 - The installed skill is discoverable as `spacedock:roborev-setup` without loading Roborev instructions into normal first-officer or ensign operation.**
Verified by: contractlint for valid frontmatter, directory/name agreement, `user-invocable: true`, manifest skill-surface discovery, and absence from FO/ensign reference closure; current-checkout install smoke resolves the skill in Codex, Claude, and Pi package surfaces. Codex carries the behavioral drive; the other hosts need discovery smoke because the skill body has no host adapter.

**AC-4 - The resulting implementation contract performs an implementation-owned semantic adversarial pass, then enforces exact-tip quick before `code_completion`, blocks Medium-or-higher findings, allows Low findings, and preserves fresh behavioral validation.**
Verified by: a live-gated implementation journey with a seeded passing-but-observably-wrong boundary case, deterministic fake jobs, and append-only command/proof journals. The worker must expose and repair the boundary case with exact-result and failure/cleanup proof before candidate HEAD; then `wait <head>` and `show <head> --json` complete before the panel starts, a Medium result produces a fixing commit and replacement wait, a Low-only stored result may follow a nonzero wait yet starts the panel without another commit, `code_completion` evidence matches the replacement head, and fresh validation consumes the stored parent before exercising its independent fixture check.

**AC-5 - Split workflow state is excluded using its actual non-default branch, and setup never queues a state-only review or Roborev-owned fix/routing action.**
Verified by: the isolated split-root fixture and fake-command journal. Assertions parse the state checkout's `.roborev.toml`, confirm its actual branch in `excluded_branches`, confirm its Git commit, and reject any state-path `review`/`post-commit` call or any `fix`, `refine`, agent-hook, commit, or entity-routing event.

**AC-6 - Sandbox advice is conditional, minimal, and never becomes broad or machine-local project configuration.**
Verified by: sandbox-present and sandbox-absent live fixture cases with a test-only `setup-recorder` decision sink. The present journal must contain one structured read-only runtime-path proposal for `~/.roborev/runtime`; the absent journal must contain no sandbox-access proposal. Parsed project files and before/after snapshots of the isolated home must show no `~/.roborev` mount, `ROBOREV_DATA_DIR`, or machine-local Safehouse change.

**AC-7 - Written `review_guidelines` equal the user-approved calibration and do not duplicate repository procedures or commands.**
Verified by: the fixture provides four structured calibration facts plus an `AGENTS.md` and workflow containing decoy procedures and commands. After the live approval turn, parse the resulting TOML and compare its semantic guideline content with the approved facts; fail if any decoy procedure, stage instruction, component command, or test command appears. A rejection turn must leave `.roborev.toml` unchanged.

**AC-8 - The test suite can refute materially broken setup behavior rather than merely certify shipped text.**
Verified by: a detached adversarial audit on a throwaway checkout mutates at least two independent boundaries—quickstart-first ordering and either state exclusion or exact-head gating—and demonstrates the relevant live/fixture oracle turns red, then restores the checkout and records clean status.

## Test architecture and plan

### Hermetic fixture

Create a temporary commissioned code workflow with implementation and fresh-validation stages, a split state checkout on a deliberately non-default branch, isolated `HOME`, and Git remotes entirely under the temp root. Seed approved agent/panel choices and calibration facts so setup does not depend on network access or interactive agent discovery.

Place recording fakes first on `PATH`:

- `roborev` appends timestamp, cwd, argv, resolved HEAD, and fixture response ID to a JSONL journal. It implements `quickstart`, relevant `--help`, `wait`, `review`, and `show` responses and rejects `fix`, `refine`, agent-hook, and unexpected state-path review calls.
- `setup-recorder` is a test-only evidence sink for proposed actions. The live prompt generically asks the installed skill to record every proposal before seeking approval; it does not name the expected sandbox decision. Production skill text never mentions this helper.

Capture the host's structured tool events separately from assistant prose. Classify browser/search tools and shell network clients as external-documentation events. This makes quickstart-first and search-negative claims independently queryable without matching transcript sentences.

### Live setup matrix

Use Codex for the behavioral drive because the canonical installed plugin and qualified `$spacedock:roborev-setup` surface are directly exercised there. Run two fixture modes in the same live-gated test family:

1. **complete quickstart:** network denied, all required playbook/panel facts seeded, zero external-doc events expected;
2. **named schema gap:** quickstart omits panel schema, local help is required next, and only a direct official panel-page open is allowed after help.

Both modes assert process exit, resulting README and TOML semantics, state-checkout git log, clean main and state statuses, and the command/decision journals. Claude and Pi reuse their existing current-checkout install lanes for lightweight discovery of the shipped user skill; they do not duplicate the expensive behavior matrix unless implementation finds a host-specific difference.

### Panel-order journey

Drive the generated implementation text once through a supported live worker with fake synthesis records:

- the baseline suite passes while a seeded Unicode/over-limit record produces wrong bytes or terminal cleanup; the pre-review pass exposes it, uses the canonical full-record validator, fixes it, and records exact/scaling/failure proof before candidate HEAD;
- tip A quick returns Medium; no `code_completion` event may exist;
- implementation fixes and commits tip B;
- tip B quick returns Low-only with nonzero wait; `show` supplies the stored severities and no further commit occurs;
- `code_completion` starts only after tip B's wait and show completion and returns a parent frozen to tip B;
- fresh validation reads that parent and writes an independent acceptance marker;
- a negative fixture that starts both panels independently fails from event order alone.

### Sandbox and guideline oracles

Run the setup drive with and without a `.safehouse`/launcher sandbox signal. Compare structured decision journals, parsed project files, and isolated-home snapshots. For guidelines, compare parsed TOML value against the user-approved fact set, then run a user-rejection case that must preserve the original file hash. These are resulting-state oracles, not assertions over the skill text.

### Structural, install, and repository gates

- Extend `internal/contractlint` only for frontmatter, user-command discovery, name/path agreement, plugin-manifest reachability, and FO/ensign load-boundary invariants. Structural tests do not claim setup behavior.
- Add current-checkout installed-skill discovery smoke to the existing Codex, Claude, and Pi runtime lanes. A missing host in a local developer run may skip its host-specific smoke; the corresponding release CI lane is required.
- Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
- Before merge, run the proof-policy detached adversarial audit because this adds shipped skill/scaffolding.

Estimated complexity: structural/install checks are low; the recording fakes and semantic TOML assertions are medium; the two-mode setup drive and implementation-order journey are high but reuse one fixture package.

## Documentation diff proposed for implementation

Add `docs/site/advanced/roborev.md`:

```markdown
# Add Roborev review evidence

Run `/spacedock:roborev-setup` from an existing commissioned code project when you want an independent review record before fresh validation.

The setup starts with `roborev quickstart`, applies only the missing steps you approve, and keeps Roborev optional. It configures a cheap exact-tip quick review before the required `code_completion` panel; Spacedock still owns fixes, routing, and behavioral validation.

For split-root workflows, setup excludes the actual state branch without creating a probe review. In a sandbox, it proposes only the read-only Roborev runtime access the client needs and never commits machine-local Safehouse settings.

Review the proposed `review_guidelines` before setup writes them. They should calibrate trust, compatibility, false positives, and review focus—not copy your agent instructions, workflow procedures, or test commands.
```

Add one navigation entry to `mkdocs.yml` after Split-root state:

```diff
       - Split-root state: advanced/split-root-state.md
+      - Add Roborev review evidence: advanced/roborev.md
       - Bridge an external tracker: advanced/external-tracker.md
```

## Resolved ideation decisions

- The user skill drives setup directly; no product setup command is justified. A test-only command and structured host event log provide durable evidence without expanding the public binary.
- Conditional sandbox advice is proven by the decision journal plus negative filesystem/home snapshots, not transcript wording.
- Codex carries the live behavior matrix. Claude and Pi carry installed-skill discovery smoke because there is no host-specific adapter in this skill.
- Quickstart and local help are primary. Direct official pages are allowed only for a named residual gap, after the local evidence; broad search is never part of setup.
- `review_guidelines` are an explicit user-reviewed calibration artifact, not a summary of repository instructions.
- The semantic adversarial pass is implementer/process guidance in the generated stage, never Roborev reviewer calibration.

## Stage Report: ideation

- DONE: Revise the canonical bundled skill draft so `roborev quickstart` is the first Roborev command, its Current state/playbook is primary, and web documentation is a targeted fallback only when quickstart leaves a specific gap; record the observed search detour as the motivating failure.
  The canonical artifact and proposal now encode quickstart → local help → named direct-page fallback, and the entity records the pilot detour.
- DONE: Complete the behavior-first design and canonical artifacts while preserving the captain's exact-tip `quick` cost gate before `code_completion`, the split-state exclusion boundary, and the ban on Roborev-owned fixes or routing.
  The skill, proposal, UI metadata, generated workflow contract, and AC-4/AC-5 specify the ordered cost gate and ownership boundary.
- DONE: Resolve the open test-design questions enough for a staff review: define durable quickstart-first/search-negative, sandbox, fake-command, installed-skill discovery, and live-drive evidence without changing the active workflow README or shipping product code in ideation.
  The test architecture selects Codex for behavior, all hosts for install discovery, JSONL fakes and host tool events for command/search proof, and decision/filesystem journals for sandbox proof.
- DONE: Narrow generated `review_guidelines` to user-reviewed calibration that the reviewer cannot infer from the repository.
  The skill and AC-7 exclude duplicated agent/workflow procedures and commands, and compare resulting TOML with approved fixture facts.
- DONE: Preserve the pilot's semantic adversarial pass as implementer-facing pre-review guidance rather than `review_guidelines`.
  The canonical skill and AC-4 cover representation/lifecycle tracing, variant matrices, atomic validation, scaling/limits, and false-green exact-result proof before quick.

### Summary

Ideation revised the canonical task artifacts around two observed failures: a documentation-search detour and over-broad generated review guidelines. The design now has independent durable oracles for command order, targeted web fallback, panel order, split-state exclusion, sandbox advice, guideline calibration, installed discovery, and adversarial refutation; no active workflow definition or product code changed. `go test ./...` and `go test ./... -race` passed after the required format gate.
