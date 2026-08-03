---
id: js6vwx74s0yg3vb88rekfkf8
title: Stakes declaration read-through — a declared `## Stakes` reaches every dispatch packet and reviewer context
status: backlog
source: "0260 shaping — agent-derail forensics audit, 2026-07-19."
score: "0.9"
sprint:
group:
started: 2026-07-20T01:44:22Z
gates:
    version: 1
    records:
        - id: gate:docs-dev:js6:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:js6-ideation-1
              briefing:
                id: briefing:js6-ideation-1a
                digest: sha256:6984a7e9a1809cfd9b34eafdcdf7b158c13fe23ae2088f3ce707781a7f3eaefa
                room-ref: "./review/ideation"
              resolution:
                type: Resolution
                id: resolution:captain-js6-ideation-1a-hold
                briefing: briefing:js6-ideation-1a
                by: person:captain
                at: 2026-07-20T03:05:44Z
                decision: hold
                reason: Parked at the sprint re-lock essence test — the read-through's would-be consumers are already served (validators read per-entity value ACs; roborev injects its config posture line); spike evidence banked in the entity body for revival.
sprint-readiness:
---

A workflow's declared rigor posture never reaches the agents who make rigor decisions. The workflow READMEs already grew that posture — the high-stakes surface list, the proof policy, the workflow rules — but a dispatched worker or reviewer never sees the README, so they default to maximum rigor regardless of what the project already declared. This read-through gives the EXISTING declared posture reach: the boot record exposes it and `dispatch build` injects it verbatim into every worker packet and review context, so an ensign or reviewer can cite the posture the repo already declared. It mints no new stakes concept and sets no value — it carries what is already there, and where nothing is declared it injects nothing and says so.

## Problem

Rigor is not proportional to anything a project can declare. The 0260 forensics corpus
(`_evidence/0260-agent-derail-forensics/synthesis.md`) shows agents mostly *complied*
with the contract: every shipped code mechanism prices under-verification (evidence or
momentum is enforced), every anti-over-engineering guard is prose-only, and nothing
anywhere expresses "this is a prototype." The word "prototype" has zero hits across the
skills tree, so a worker or reviewer holding a finding defaults to maximum rigor
regardless of what the project wants. Concrete derailments this leaves unpriced: a PTY
process-control harness built for a disposable zellij smoke test (zaphod
codex:019f5160-8969-75c1-a79c-244148b40a0e:472-523), a symlink-edge-case rejection on a
prototype (spacedock_subspace codex:019f63c6), and the audit's own finding that no
frontmatter field, stage property, or boot-record key can carry a stakes signal.

The fix is to give the posture the repo *already declared* reach to the agents who make rigor
decisions — not a new signal. This entity is the substrate the `triage`, `ladder`, and
`template` groups cite: the read-through that carries a workflow's existing declared-posture
section from one source into the boot record and every dispatch packet, verbatim. It delivers
the *flow*, not the value — it carries what the workflow declared and decides nothing. Where a
workflow declares no posture section, the read-through injects nothing and labels the absence
explicitly; the stage definitions and the committed finding-triage govern, exactly as they do
today. (This reduction supersedes the earlier captain ruling of a fixed prototype-grade default
plus a nag: the sprint re-lock, 2026-07-20, cut the stakes member from minting a new ontology
to a read-through of existing declared posture, so absence needs no injected directive — the
existing stage-def-plus-triage behavior already IS the default.)

Two questions had to be settled before designing the flow, and both rest on unverified
mechanisms, so they were spiked first (see **Spike results**): does an authored posture
section round-trip verbatim through the shipped README parser that boot and dispatch already
use, and does a Claude worker actually receive an AGENTS.md digest (the assumed router
channel) without it being planted in the packet?

## Spike results

Riskiest mechanisms exercised before committing to the design. Commands and outputs were
run against the shipped `spacedock 0.26.0-pre0` binary and `claude 2.1.212`; the throwaway
fixtures live under the session scratchpad.

**1. `## Stakes` section round-trips verbatim through the shipped parser (the substrate
boot + dispatch build reuse).** A fixture README with a multi-paragraph `## Stakes`
section was read by the shipped `status --read <readme> --json`; the parser
(`internal/status/section_read.go` `scanHeadings`, the same fence-safe heading map that
backs `--read` and `extractStageSubsection`) located `## Stakes` at offset 19, lines 11.
Slicing exactly that span and applying the trailing-blank trim the existing stage
extractor already applies (`internal/dispatch/showstagedef.go:111-113`) produced a byte-
identical copy of the authored section (`sliced == authored` True). An inline `` `## note` ``
inside the body was correctly NOT mis-detected as a heading (fence/inline-safe). Result:
the "declaration → machine-addressable section → verbatim slice" path both channels need
is already proven code; no new parser is required.

**2. Dispatch-packet baseline is zero.** `dispatch build --workflow-dir docs/dev --stage
ideation` for a real entity produced a 1018-byte packet whose only "stakes" substring is
the entity slug in the dispatch-file path; it carries none of a stakes declaration
("who depends", "what a defect costs", "derived policy" all absent). So the value AC's
baseline is a measured 0 declaration-bytes-in-packet on `main`.

**3. AGENTS.md-ingestion canary — Claude does NOT auto-ingest AGENTS.md.** A throwaway
directory carried distinct codewords in a project `AGENTS.md` (the claim) and a project
`CLAUDE.md` (positive control). A fresh non-interactive `claude -p` session in that dir,
told to answer only from already-loaded context with no file reads, returned the
`CLAUDE.md` codeword and `NONE` for the `AGENTS.md` codeword — reproduced across three
runs, including an AGENTS.md-only dir (no CLAUDE.md fallback fires). The positive control
(CLAUDE.md seen) proves the method can detect ingestion, so the AGENTS.md null is real,
not a blind test. **Design consequence:** the dispatch packet — this entity's mechanism —
is the reliable stakes channel for Claude workers. AGENTS.md serves codex (the majority
worker runtime, which does read it); Claude ad-hoc sessions would need a `CLAUDE.md`
import, not AGENTS.md. That per-runtime split is the router-layer/`template` group's
concern; it is recorded here because it is why the packet channel, not a shared AGENTS.md
digest, is the substrate.

## Proposed approach

### Recommendation on the two design axes

**Placement — a prose declared-posture *section* in the workflow README (not frontmatter).**
A rigor posture is inherently multi-paragraph prose (which surfaces are promised, what a
defect costs, what is outside the promise). A section carries that natively and is carried
*verbatim* into a packet by the already-proven section extractor (Spike result 1). The boot
record exposes the section body as a JSON string, so it is machine-readable without a second
representation.

- *Losing alternative — a README frontmatter scalar.* It is trivially machine-read
  (`ParseFrontmatter` already parses README frontmatter), but a YAML scalar mangles multi-
  paragraph prose (escaping, block-scalar clutter, no markdown), and "carried verbatim into
  the packet" is section-native, not scalar-native. It loses on fidelity, and nothing needs a
  typed field yet; a section does not preclude adding one later.

**Content — a consolidation of the workflow's EXISTING declared posture, not a new stakes
ontology.** The section carries what the repo already grew (the Proof policy's high-stakes
surface list and defect-cost framing, the workflow rules), distilled into one readable place
so it can travel. It is a restatement with a pointer to the authority, not an invented fact:
the read-through gives existing declarations reach, it does not mint a rigor taxonomy. This is
the sprint re-lock's reduction (2026-07-20) — "reach for declarations the repos already grew,
not a new stakes ontology" — and it is why the cycle-2 four-part worked example survives here
only as the *shape of the consolidation*, no longer as a new concept the contract introduces.

- *Losing alternative — minting a new stakes ontology (a `## Stakes` fact with its own
  taxonomy, cycle-2's direction).* Rejected at the re-lock: the cycle-2 worked example proved
  the ontology was nearly all restatement of declarations the repo already had, so a new
  concept added a fourth divergent copy to maintain for no new information.
- *Losing alternative — reading the raw `## Proof policy` section verbatim (no consolidation).*
  Loses on packet weight and mixing: Proof policy is long and interleaves posture with CI-lane
  mechanics, so carrying it whole bloats every packet and buries the posture. A short
  consolidation that points back to Proof policy travels lean and stays readable.

### Mechanism

One source, read by two code paths that must never diverge:

1. **Source of truth:** the workflow README's declared-posture section — heading `## Rigor
   posture` (`docs/dev/README.md` for this workflow). Declared once; the extractor is heading-
   agnostic (Spike result 1 proved it on an arbitrary section), so the heading is the one
   small convention.
2. **Boot channel:** `status --boot` reads the README's `## Rigor posture` section via the
   existing heading-map extractor and exposes its body under a `rigor_posture` key in
   `--boot --json` (and a `RIGOR_POSTURE:` line in the text form). No section → `none declared`
   plus a one-line label (`stage definitions and committed finding-triage govern`), so absence
   is explicit, never silent. Read from the same `definitionDir/README.md` boot already opens.
3. **Packet channel:** `dispatch build` extracts the same `## Rigor posture` section and
   injects it **inline verbatim** as a block in every worker packet (`internal/dispatch/build.go`
   prompt-`parts` assembly), placed with the assignment header so it frames the work. No
   section → nothing is injected except the same one-line explicit-absence label. This reaches
   ensigns and every dispatched reviewer (staff-review, validation) uniformly.

**Read cost — no extra roundtrip for any consumer.** The packet block is assembled inline at
`dispatch build` time (build already reads the workflow README for the stage/frontmatter
data), so an ensign receives stakes in the single dispatch-file Read it already performs —
there is deliberately no `show-stakes` fetch line. Boot rides the existing `--boot` read (the
`--boot` path already parses the same `definitionDir/README.md` for the stage taxonomy), so
the FO pays no extra command. Every dispatched reviewer gets it inline in its own packet. Net
added tool calls across all consumers: zero.

**Governance — existing README ownership, plus a sanctioned drift flag.** The posture section
lives in the workflow README and is owned and changed exactly as the rest of that authored
process doc is — no new sign-off apparatus, write-classifier carve-out, or ceremony
re-affirmation. (The reduction dropped all of that: the earlier captain-only-plus-ceremony
design was tied to a *minted* `## Stakes` fact; a read-through of declarations the repo already
owns does not need a new authority layer.) The one governance mechanism that survives is the
**stakes-drift flag**: when a worker or reviewer hits evidence that contradicts the declared
posture, it does NOT re-grade on its own — it triages the finding under the DECLARED posture
anyway and raises a drift flag carrying the contradicting evidence, which the FO surfaces at
the gate for a human decision. That keeps the anti-inference guarantee (no silent agent
re-grading) while giving contradicting evidence a sanctioned path instead of a private
judgment call.

**The undeclared default injects nothing and says so.** A workflow with no `## Rigor posture`
section resolves to the behavior that *already happens today*: nothing is injected, the boot
`rigor_posture` is `none declared` with the one-line label, and the stage definitions plus the
committed finding-triage govern. There is no injected rigor directive and no nag — the earlier
prototype-grade-default-plus-nag ruling was superseded by the re-lock, because a read-through
has nothing to assert where the repo declared nothing, and the existing stage-def-plus-triage
behavior IS the default. Absence is still *explicit* (the label), so a reader can tell "no
posture declared" from "posture declared as empty".

**Stage-differential stakes — uniform injection because stakes *parameterizes* stage rules,
it does not replace them.** The declaration is injected uniformly into every stage's packet,
byte-identical at ideation and validation, because who-depends / promised-surfaces / costliest-
defects / outside-the-promise is a project fact that does not change with the stage. The way
stakes affects a stage is by parameterizing that stage's existing rules: high stakes shifts
which risk a stage's work *selects* — at ideation, which path the "spike the riskiest path"
rule aims at; at validation, which surfaces earn the detached adversarial audit — but it never
changes the *polish* a given unit of work gets. A spike stays throwaway at any stakes level;
high stakes makes you spike the *promised* surface first, not gold-plate the spike. The stage
definitions keep all the per-stage rigor mechanics (already carried into every packet via the
existing `show-stage-def` fetch line); the declaration is the parameter they read. So the
worker composes `promised surfaces = status/dispatch/state` (declaration part b) + `spike the
riskiest path` (ideation stage def) into "spike the riskiest status/dispatch path", and does
not over-build a throwaway.

The "stakes shifts spike selection, never polish" *behavior* is a worker-consumption rule; its
falsifiable AC (a high-stakes context changes which path a worker spikes, not how much) belongs
to the `triage`/`ladder` groups that own finding-triage and own the ensign's consumption of the
field — they blocker-depend on this entity. This entity's job is to make that composition
*possible* by delivering the four-part declaration to every stage uniformly and provably (AC-1,
AC-1d); it does not re-implement the consumption behavior here.

- *Losing alternative — per-stage-aware injection (the mechanism selects different stakes text
  per stage).* Loses: it duplicates the stage-appropriate rigor the stage definitions already
  express and already deliver, couples the stakes mechanism to the stage taxonomy, and splits
  the single source into per-stage copies — against leanness and "one source".
- *Losing alternative — a per-stage derivation table inside the declaration prose.* Loses: it
  bloats a workflow-level fact with stage mechanics that belong in (and duplicate) the stage
  definitions, and it is carried verbatim into stages it does not apply to. Keep the
  declaration to the four project-level parts; the stage defs own per-stage depth. The dev
  template's stage definitions are the natural home for the per-stage half.

Each new mechanism, the value AC it serves, the simplest alternative, and why that
alternative is insufficient:

1. *Reusing the README section extractor in boot + dispatch build* (serves AC-1/AC-4).
   Simplest alternative: a frontmatter scalar read by `ParseFrontmatter`. Insufficient — see
   the placement axis: a scalar mangles multi-paragraph prose and cannot be carried verbatim.
2. *Inline verbatim injection into the packet, not a `dispatch show-stakes` fetch line*
   (serves AC-1). Simplest alternative: a fetch line like `show-stage-def`. Insufficient — a
   dispatched reviewer's context and roborev do not run fetch commands, so a fetch line would
   not reach reviewers; and stakes is short, so fetching saves none of the bytes that justify
   fetching the large stage def, while a skipped or failed fetch would silently drop a load-
   bearing signal.
3. *A `stakes` key in the boot record* (serves AC-2). Simplest alternative: the FO reads the
   README `## Stakes` by hand at boot. Insufficient — the boot record is the machine-read
   startup surface the triage/ladder groups cite; a hand-read is neither exposed to them nor
   testable as "the field exists and tracks the README."

### Documentation changes (applied by implementation)

This entity changes user-visible surfaces (the boot record gains a `STAKES` line/key; every
packet gains a stakes block; the dev README gains a section). The concrete doc diff:

*Add to `docs/dev/README.md`, a new top-level section (the workflow's own declaration):*

```markdown
## Stakes

**Who depends / what breaks:** every agent session in every repo that loads the shipped
Spacedock contract — the first-officer/ensign skills and the `status`/`dispatch` launcher
— across spacedock_v1, zaphod, and spacedock_subspace. A defect derails real multi-hour
sessions at scale (the 0260 forensics corpus: 15 confirmed derailments, 4 multi-hour
runaway loops) and can lose a worker's committed state.

**Promised surfaces (high-stakes):** the front-door launcher; the `status` mutation and
guard paths; the `.spacedock-state` checkout and the dispatch/launch path; the CI and
release machinery. Changes to these earn the full depth — the falsifiability rule,
behavior-over-prose tests, and the detached adversarial audit before merge.

**Most-expensive defect classes:** a false green — a test that passes while the behavior is
wrong (the presence/prose-grep tautologies this sprint is retiring) — because it ships a
regression undetected; and silent state or dispatch corruption that loses committed work.
Rigor is spent to make these two catchable, not spread evenly.

**Explicitly outside the promise:** contract and skill *prose* wording (proven by a live
drive, never a grep over the file); throwaway spikes (the smallest sufficient exercise,
not gold-plated); and research/analysis entities whose product is a decision recorded in
the roadmap. A correct-but-disproportionate finding on one of these is a candidate decline,
not a dutiful fix.
```

*Add one mechanism note under `docs/dev/README.md` `### Reading sections` (kept OUTSIDE the
`## Stakes` section so it is not itself carried into packets):*

```markdown
The workflow's `## Stakes` section is read by `status --boot` and injected verbatim into
every dispatch packet, so a declared stakes reaches every worker and reviewer. Declare it
once here as four parts (who depends / what breaks; promised surfaces; costliest defect
classes; explicitly outside the promise); never copy it into entity bodies. Its content is
captain-only: an FO may scaffold the heading, but the value changes only by the captain, and
it is re-affirmed at ceremonies (sprint scope-lock, release cut), not on a timer. A worker or
reviewer whose evidence contradicts the declaration triages under the DECLARED stakes anyway
and raises a stakes-drift flag with that evidence for the captain at the gate — it never
re-grades stakes on its own. A workflow with no `## Stakes` section is treated as
prototype-grade (smallest sufficient test surface, no new enforcement infra without consent,
lean triage), and boot/gate presentations nag until one is declared.
```

## Out of scope

- **Grading a *declared* project's stakes.** The contract carries the four-part declaration
  and sets exactly one value — the fixed prototype-grade default when nothing is declared. It
  never grades or escalates a declared project's stakes per diff.
- **The worker-consumption behavior of the field** — the ensign's finding-triage against
  stakes, the recorded decline of a disproportionate-but-correct finding, and the falsifiable
  "high stakes shifts spike selection, not polish" AC. These are the `triage` and `ladder`
  groups (which blocker-depend on this entity). This entity delivers the field to every stage;
  it does not re-implement how a worker acts on it.
- **The AGENTS.md one-line digest + pointer** (router-layer / `template` group) and **the
  roborev config alignment** (repo-local group). This entity delivers the boot + dispatch-
  packet channels (which reach ensigns and every dispatched reviewer). The canary result that
  shapes those channels is recorded above; wiring them is theirs.
- **The four-part `## Stakes` template scaffold + commission rigor question** (`template`
  group). This entity ships the field, the read-through, and the dev workflow's own
  declaration; the reusable scaffold future commissions start from is the template's.
- **The write-classifier carve-out that enforces captain-only edits.** It is the recommended
  enforcement (see Governance) but is an FO-contract surface; whether the code change lands
  with this entity or adjacent is the gate's call (probe 1). The captain-only *convention*
  ships here in the README note regardless.
- **Per-entity stakes override.** Stakes is workflow-level here (one `## Stakes` per README).
  An entity-level override is a possible future extension, not this entity.
- **`status --validate` hard-requiring a `## Stakes` section.** Absence is made *loud* (the
  boot/gate nag) and *safe* (the fixed prototype-grade default); hard-fail validation would
  break every existing workflow that predates the field and belongs with the `template`/
  validate work. Keeping this surface minimal and stable is deliberate — downstream groups
  cite the field.

## Acceptance criteria

Each AC names a property of the finished entity and how it is verified. Verified-by clauses
compare generated output against an independent source (the README) or a live run; none is
a substring match over an instruction file the implementer wrote.

**AC-1 (VALUE) — A dispatch packet reproduces its workflow's declared `## Stakes` section,
byte-for-byte and derived (not hardcoded); a stakes-less workflow gets an explicit
`none declared` marker.**
Verified by: a Go behavior fixture in `internal/dispatch` that (a) runs `dispatch build`
against a fixture README with a known `## Stakes` section and asserts the packet's stakes
block equals that section (trailing-blank-trimmed) by exact string compare; (b) mutates the
README section, rebuilds, and asserts the packet block changes to the new text — a hardcoded
copy passes (a) but fails (b); (c) runs against a `## Stakes`-less README and asserts the
packet carries the `none declared` marker; (d) runs the stakes-bearing fixture at two
different stages (ideation and validation) and asserts the stakes block is byte-identical
across them, pinning the uniform workflow-level injection. Baseline that can move the wrong
way: on `main` today the packet carries 0 stakes-declaration bytes (Spike result 2); the AC
fails if the block is absent, stale, truncated, varies by stage, or fails to track a README
edit.

**AC-2 — `status --boot --json` exposes the workflow's `## Stakes` body under a `stakes`
key, tracking the README; absent → an explicit none value.**
Verified by: a Go behavior/golden fixture asserting the `--boot --json` `stakes` field equals
the fixture README's `## Stakes` body, equals the none value for a stakes-less fixture, and
changes when the README section is edited and boot re-run. Fails if boot omits the field,
hardcodes it, or diverges from the README.

**AC-3 — One source: the packet block and the boot field derive from the same README
section and never diverge.**
Verified by: a test that, for one fixture, asserts `dispatch build`'s packet stakes block and
`status --boot --json`'s stakes field are byte-identical after the shared trailing-blank
trim, and that a single README edit moves BOTH. Fails if the two read paths diverge (e.g.
one reads a section, the other frontmatter).

**AC-4 — The dev workflow declares `## Stakes` once, and a live `dispatch build` +
`status --boot` for `docs/dev` carry it (the sprint DoD line).**
Verified by: recorded output of `spacedock dispatch build --workflow-dir docs/dev` for a real
entity/stage showing the dev README's `## Stakes` verbatim in the packet, and
`spacedock status --boot --json --workflow-dir docs/dev` showing the same body under
`stakes`. Fails if the dev README lacks the section or the live carry drops or mangles it.

**AC-5 — An undeclared workflow resolves to the fixed prototype-grade default and nags,
rather than a bare or empty value.**
When a workflow has no `## Stakes` section, both the packet and `status --boot` emit a
`none declared` marker carrying the fixed default (prototype-grade: smallest sufficient test
surface; no new enforcement infra without consent; lean triage), and the `--boot` text form
emits a one-line `STAKES: none declared` nag. Verified by: a toggle test that runs
`dispatch build` and `--boot` (text and `--json`) against one fixture with and without the
section and asserts the output SWAPS on the section's presence — the verbatim section in the
present case, the prototype-grade default marker + the boot nag line in the absent case, and
neither leaking into the other. The default text is a single system constant, so the test also
asserts it is byte-identical across two different fixtures that both lack a declaration (it is
not per-workflow inference). Fails if absence yields a bare/empty value, if the marker leaks
into the declared case, if the nag line is missing, or if the default text varies by workflow.

## Test plan

- **Fixture/CLI (primary, cheap):** the AC-1/AC-2/AC-3/AC-5 Go behavior fixtures in
  `internal/dispatch` and `internal/status` — packet-carries-verbatim, boot-exposes,
  edit-follows, cross-channel identity, stage-invariance (same block at ideation and
  validation, AC-1d), the presence-toggle that swaps section↔prototype-default marker with the
  `STAKES: none declared` boot nag (AC-5), and the default-constant invariance across two
  undeclared fixtures (AC-5, proving the default is system-fixed not per-workflow). These drive
  the built binary / package functions and assert on generated bytes, not on instruction-file
  text. No new harness, no new infra: they reuse the spike-proven `scanHeadings`/section-
  extractor substrate.
- **Live smoke (one run, AC-4):** after the dev README gains `## Stakes`, run `dispatch
  build` and `status --boot --json` against `docs/dev` and record that both carry the
  section. Single-command each; this is the DoD's live proof.
- **No spike outstanding:** the two riskiest mechanisms (section round-trip; AGENTS.md
  ingestion) were exercised in **Spike results**; the round-trip fixture there seeds AC-1's
  first test.
- **Cost/complexity:** low. All fixture/CLI except one live smoke; the change reuses proven
  code, adds one boot field, one packet block, and one README section. High-stakes surface
  (shipped contract/scaffolding + `status`/dispatch paths), so implementation owes the
  detached adversarial audit before merge per the Proof policy.

## Stage Report: ideation

- DONE: Riskiest mechanism exercised first — live spike of declaration → boot → packet, plus the AGENTS.md-ingestion canary, evidence recorded not asserted
  See `## Spike results`: shipped `status --read` addresses `## Stakes` (offset 19, lines 11) and the parser-directed slice is byte-identical to the authored section; `dispatch build` for docs/dev carries 0 declaration-bytes on `main`; `claude 2.1.212` returns NONE for an AGENTS.md codeword across 3 runs while the CLAUDE.md positive control is seen.
- DONE: One recommended most-practical design on both axes, each losing alternative named and why it lost
  See `## Proposed approach`: prose `## Stakes` section (frontmatter loses on prose fidelity; both loses on YAGNI + divergence) and injected prose hint (defined levels lose because the contract must never set the value, an enum is lossy and speculative).
- DONE: At least one AC measures end-value against a baseline that can move the wrong way; every Verified-by can fail; no prose-grep
  AC-1 (VALUE) measures stakes-declaration bytes reaching a packet: 0 on `main` (spike-measured) → the verbatim section, with an edit-follows clause that a hardcoded copy fails. AC-1..4 all compare generated packet/boot output against the independent README or a live run.

### Summary

Designed the stakes read-through as one source (a prose `## Stakes` section in the workflow
README) carried by two non-diverging code paths — a `stakes` key in `status --boot --json`
and an inline verbatim block in every `dispatch build` packet — reusing the already-proven
fence-safe section extractor, so no new parser or infra. Two unverified mechanisms were
spiked first: the section round-trips byte-identically through the shipped parser, and the
AGENTS.md-ingestion canary proved Claude 2.1.212 does NOT auto-ingest AGENTS.md (CLAUDE.md
control confirms the test), which is why the packet — not a shared AGENTS.md digest — is the
Claude-worker channel. Notable: the placement and rigor axes both resolve *against* the more
mechanical option because the sprint requires the contract to carry the declaration without
ever setting its value; the frontmatter/levels alternatives are named but recorded as
losing. A concrete doc diff for the dev README (`## Stakes` section + a Reading-sections
mechanism note) is in the body for implementation to apply.

## Stage Report: ideation (cycle 2)

Revised against captain gate feedback (three design questions). Problem, approach, ACs, and
test plan updated together per the ideation stage rules.

- DONE: Q1 — lifecycle/governance and the `none declared` default designed explicitly
  Added `### Mechanism` "Governance" (the declaration is a captain-declared fact — FO may scaffold the heading, value changes only with captain sign-off, no agent infers/rewrites; convention not a code author-gate, which is out of scope) and "The `none declared` default is explicit" (the marker carries its own directive: apply stage-def rigor and no more, do not infer project stakes, ask the captain to declare — caps inference at the written stage-def rigor instead of reinstating silent max-rigor). Reflected in the Problem corollary, the README note doc-diff, and new AC-5 (presence-toggle swaps section↔directive marker).
- DONE: Q2 — read cost confirmed as a one-liner in Mechanism
  Added "Read cost — no extra roundtrip for any consumer": inline at build time (no `show-stakes` fetch line), boot rides the existing `--boot` read (verified: the `--boot` path already parses `definitionDir/README.md` for the stage taxonomy — `internal/status/handlers.go:438`), reviewers inline. Net added tool calls: zero.
- DONE: Q3 — stage-differential decided, losing alternatives named
  Decision: uniform workflow-level declaration; per-stage depth stays in the stage definitions (already carried via the `show-stage-def` fetch line) and the declaration's derived-policy prose composes with them (spike carve-out). Losing: per-stage-aware injection (duplicates + couples to the stage taxonomy + splits the source) and a per-stage table in the declaration (bloats a workflow fact, carried into stages it does not apply to). Pinned by AC-1(d): the block is byte-identical across ideation and validation. Dev `## Stakes` derived-policy updated to state the composition.

### Summary

Cycle-2 revision addresses all three captain questions without enlarging the surface: the
governance and default semantics are convention + a self-documenting `none declared` marker
(no new code gate), read cost is confirmed zero-extra-roundtrip against the actual boot code
path, and the stage-differential resolves to uniform injection with per-stage rigor left in
its existing home (the stage definitions). ACs grew by one (AC-5, the presence-toggle) and
AC-1 gained a stage-invariance clause; both are behavioral (branch selection driven by the
README's section presence), not prose-greps. No mechanism, spike result, or the two-axis
recommendation changed — the additions are default/governance/stage-composition semantics the
original design implied but had not stated.

## Stage Report: ideation (cycle 3)

Revised against the captain's three rulings (probes 1 and 3 were insufficient-evidence against
briefing-1a; this revision makes them answerable). Problem, approach, ACs, and test plan
updated together. Briefing-2a written; ProbeResults appended for all three probes.

- DONE: Ruling 1 — undeclared default is now fixed prototype-grade + a visible nag
  Reversed cycle-2's "apply stage-def rigor". `### Mechanism` "The undeclared default is prototype-grade, and it nags": smallest sufficient test surface / no new enforcement infra without consent / lean triage, a system constant (not per-diff inference), with a boot `STAKES: none declared` nag and a gate-presentation nag (present-gate reads the boot field). Problem corollary, README note, and AC-5 rewritten to match; AC-5 now asserts the default text is byte-identical across two undeclared fixtures (proving system-fixed).
- DONE: Ruling 2 — staleness ergonomics: event-triggered re-affirmation + sanctioned stakes-drift flag
  `### Mechanism` Governance now has three parts: edit authority captain-only (write-classifier carve-out recommended, FO-proposes-captain-approves the named alternative — gate decides which via probe 1); evolution at ceremonies only (sprint scope-lock, release cut; no timers); and the stakes-drift flag (worker/reviewer triages under DECLARED stakes, flags drift + evidence, FO surfaces at the gate, only the captain edits).
- DONE: Ruling 3 — declaration is the opinionated four-part shape, not a level label
  Second design axis retitled "Shape"; the dev `## Stakes` doc-diff rewritten into four labeled parts (who depends/what breaks; promised surfaces; costliest defect classes; outside the promise) with the spacedock-v1 worked example (status/dispatch/state high-stakes; false greens the costliest class; contract prose + spikes outside the promise). Rationale recorded: a bare "high" label endorses the two observed HIGH failure modes (fabricated rigor, rigor on unpromised surfaces); the four parts are what let materiality triage and the falsifiability rule bite AT high stakes. Losing alternatives named: bare level label, and cycle-2's free-form hint.
- DONE: Stage-differential reframed to parameterization; "shifts spike selection not polish" AC located in triage/ladder
  `### Mechanism` stage-differential now: stakes parameterizes stage rules (shifts spike SELECTION, never polish); stage defs keep the mechanics; the falsifiable selection-behavior AC belongs to the `triage`/`ladder` consumers (out of scope here, recorded). Uniform injection unchanged; still pinned by AC-1d.

### Summary

Cycle-3 folds the captain's three rulings without growing the shipped surface: the flow
(README→boot+packet), the fixed prototype-grade default marker + nag, and the dev four-part
declaration are what ships; governance (captain-only, ceremony re-affirmation, drift-flag),
the write-classifier enforcement, the four-part template scaffold, and the worker-consumption
"selection-not-polish" AC are recorded as conventions or sibling-owned with explicit pointers.
The two-axis recommendation is unchanged on placement (prose section) and sharpened on shape
(four-part, replacing the free-form hint). One open item is surfaced honestly for the gate: whether
the captain-only write-classifier carve-out ships here or adjacent (probe 1). Briefing-2a and
three fresh ProbeResults bound to its digest are in `review/ideation/`.

## Parked (0260 re-lock, 2026-07-20)

Captain decision: hold, parked. The direction trail, kept honest: filed as a new
`## Stakes` ontology; reduced to an existing-posture read-through when the worked example
proved the ontology was restatement; parked when the bare kill-test ("which observed
failure does this kill, and is the killer already in place?") showed its consumers already
served — validators anchor materiality on per-entity value ACs and the release promise
(committed triage), and roborev's config already injects its posture line into every
review. Banked for revival: the Spike results above (byte-identical section round-trip
through the shipped parser; measured zero-baseline packet; AGENTS.md-ingestion canary
null with positive control). If a later member's ideation surfaces a genuinely unreached
consumer, this entity revives with its riskiest mechanisms pre-proven. The one residue —
consolidating the dev README's scattered posture declarations into a single section — is
an ordinary process-doc edit needing no entity; it rides template-group work or the next
README touch.
