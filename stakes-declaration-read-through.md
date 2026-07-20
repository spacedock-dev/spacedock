---
id: js6vwx74s0yg3vb88rekfkf8
title: Stakes declaration read-through — a declared `## Stakes` reaches every dispatch packet and reviewer context
status: ideation
source: "0260 shaping — agent-derail forensics audit, 2026-07-19."
score: "0.9"
sprint: 0260-proportionality
group: stakes
started: 2026-07-20T01:44:22Z
gates:
  version: 1
  current:
    gate: gate:docs-dev:js6:ideation
    attempt: gate-attempt:js6-ideation-1
  records:
    - id: gate:docs-dev:js6:ideation
      stage: ideation
      current-attempt: gate-attempt:js6-ideation-1
      attempts:
        - id: gate-attempt:js6-ideation-1
          sequence: 1
          state: open
          briefing:
            id: briefing:js6-ideation-1a
            digest: sha256:6984a7e9a1809cfd9b34eafdcdf7b158c13fe23ae2088f3ce707781a7f3eaefa
            room-ref: "_reviews/js6-ideation"
          note: "Captain reviewed via Subspace advisory float 2026-07-20 and took the leave-open (hold) action: decision held for captain/FO discussion of three fundamentals — declaration governance/evolution/default semantics, read cost (answered: zero added roundtrips), and stage-differential stakes. No resolution recorded; attempt remains open; revised briefing to follow the discussion."
---

No rigor dial exists anywhere in the system: "prototype" has zero hits across the skills tree, so workers and reviewers default to maximum rigor regardless of what the project wants. The workflow README declares stakes (who depends on this, what a defect costs, derived test-depth/infra/materiality policies); the boot record exposes it; dispatch build injects it verbatim into worker packets and review context so an ensign or reviewer can cite it. The contract requires the declaration exists and flows — it never sets the value.

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

The fix is a single declared stakes signal that *reaches* the agents who make rigor
decisions. This entity is the substrate the rest of the `stakes` group and the `triage`,
`ladder`, and `template` groups cite: the read-through that carries a workflow's declared
`## Stakes` from one source into the boot record and every dispatch packet, verbatim. It
delivers the *flow*, not the value — the contract requires the declaration exists and is
carried; it never decides what a project's stakes are. A corollary the design must honor:
the *absence* of a declaration must itself resolve to an explicit, bounded default — an
undefined default silently reinstates the max-rigor inference this field exists to stop.

Two questions had to be settled before designing the flow, and both rest on unverified
mechanisms, so they were spiked first (see **Spike results**): does a `## Stakes` section
round-trip verbatim through the shipped README parser that boot and dispatch already use,
and does a Claude worker actually receive an AGENTS.md digest (the assumed router channel)
without it being planted in the packet?

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

**Placement — a prose `## Stakes` section in the workflow README (not frontmatter, not
both).** A stakes declaration is inherently multi-paragraph prose (who depends on this,
what a defect costs, the derived test-depth/infra/materiality policy). A section carries
that natively and is carried *verbatim* into a packet by the already-proven section
extractor. The boot record exposes the section body as a JSON string, so it is machine-
readable without a second representation.

- *Losing alternative — README frontmatter `stakes:` scalar.* It is trivially machine-read
  (`ParseFrontmatter` already parses README frontmatter), but a YAML scalar mangles multi-
  paragraph prose (escaping, block-scalar clutter, no markdown), and "carried verbatim into
  the packet" is section-native, not scalar-native. The sprint's own artifacts name a
  `## Stakes` *section* throughout. It loses on fidelity.
- *Losing alternative — both (a machine flag in frontmatter plus the prose section).* It
  loses to YAGNI and divergence risk: no consumer needs a machine-parseable stakes *level*
  yet (the field is cited as text by triage/ladder/template), and two representations of one
  fact in one README is exactly the divergence the sprint's "one source, never a fourth
  copy" rule fights. Adopt frontmatter only if and when a consumer needs a typed level; the
  section does not preclude adding one later.

**Rigor — an injected prose hint (a verbatim `## Stakes` block carried into boot + every
packet), NOT defined stakes levels with per-level behaviors.** The contract must "require
the declaration exists and flows, never set the value" (sprint goal, verbatim). A prose
hint that the worker/reviewer reads and applies with judgment satisfies that exactly.

- *Losing alternative — defined stakes levels (e.g. `prototype|production|critical`) with a
  behavior spec per level.* It loses on three counts: it puts the system in the business of
  *setting* rigor (the one thing the sprint says the contract must not do); a fixed enum
  cannot carry the who-depends / what-a-defect-costs nuance a prose declaration can, so it
  couples every downstream group to a lossy taxonomy; and it is speculative boot-resident
  machinery against the leanness constraint. A project that wants a crisp level simply
  writes it as the first line of its prose `## Stakes`.

### Mechanism

One source, read by two code paths that must never diverge:

1. **Source of truth:** the workflow README's `## Stakes` section (`docs/dev/README.md` for
   this workflow). Declared once.
2. **Boot channel:** `status --boot` reads the README's `## Stakes` section via the existing
   heading-map extractor and exposes its body under a `stakes` key in `--boot --json` (and a
   `STAKES:` line in the text form). Absent section → an explicit `none declared` value, so
   absence is visible, never silent. Read from the same `definitionDir/README.md` boot
   already opens.
3. **Packet channel:** `dispatch build` extracts the same `## Stakes` section and injects it
   **inline verbatim** as a block in every worker packet (`internal/dispatch/build.go`
   prompt-`parts` assembly), placed with the assignment header so it frames the work.
   Absent → the same `none declared` marker. This reaches ensigns and every dispatched
   reviewer (staff-review, validation) uniformly.

**Read cost — no extra roundtrip for any consumer.** The packet block is assembled inline at
`dispatch build` time (build already reads the workflow README for the stage/frontmatter
data), so an ensign receives stakes in the single dispatch-file Read it already performs —
there is deliberately no `show-stakes` fetch line. Boot rides the existing `--boot` read (the
`--boot` path already parses the same `definitionDir/README.md` for the stage taxonomy), so
the FO pays no extra command. Every dispatched reviewer gets it inline in its own packet. Net
added tool calls across all consumers: zero.

**Governance — the declaration is a captain-declared fact.** By the sprint's anti-inference
principle the `## Stakes` *value* is set and changed only by the captain (or with captain
sign-off), recorded like a durable gate decision. The README is FO-writable process doc, so
an FO may scaffold the empty heading at commission/refit, but no worker, reviewer, or FO
infers, escalates, or rewrites the declared value. This is a documented convention, not a
code-enforced author gate — identity-gating the edit needs authorship tracking this repo does
not have and is out of scope; the convention is the contract and ships as a line in the README
note (below).

**The `none declared` default is explicit, not silent max-rigor.** An undefined default leaves
a worker to infer rigor — the disease itself. So the marker is not a bare "none": it carries
its own default behavior verbatim — *apply the rigor the stage definitions specify and no
more; do not infer additional project-level stakes; if this workflow needs a different
baseline, ask the captain to declare a `## Stakes` section.* This caps inference at the written
stage-def rigor (which still mandates, e.g., the detached adversarial audit for high-stakes
surfaces, so safety is not lowered) and turns a missing declaration into a visible prompt to
declare rather than a licence to over-build.

**Stage-differential stakes — uniform declaration; per-stage depth stays in the stage
definitions.** The declaration is injected uniformly into every stage's packet, byte-identical
at ideation and validation, because who-depends / what-a-defect-costs is a project fact that
does not change with the stage. Per-*stage* rigor already lives in the stage definitions
(ideation "spike the riskiest path" cheaply and throwaway; validation's detached adversarial
audit + semantic adversarial pass) and is already carried into every packet via the existing
`show-stage-def` fetch line. The declaration's *derived-policy* prose is written to COMPOSE
with them — it states the project's weight and explicitly carves out spikes — so a worker sees
`high project stakes` (declaration) + `spike cheaply` (ideation stage def) + `spikes exempt
from full depth` (declaration carve-out) and does not over-build a throwaway.

- *Losing alternative — per-stage-aware injection (the mechanism selects different stakes text
  per stage).* Loses: it duplicates the stage-appropriate rigor the stage definitions already
  express and already deliver, couples the stakes mechanism to the stage taxonomy, and splits
  the single source into per-stage copies — against leanness and "one source".
- *Losing alternative — a per-stage derivation table inside the declaration prose.* Loses: it
  bloats a workflow-level fact with stage mechanics that belong in (and duplicate) the stage
  definitions, and it is carried verbatim into stages it does not apply to. Keep the
  declaration to the project fact plus a one-line spike carve-out; the stage defs own per-stage
  depth. The dev template's stage definitions are the natural home for the per-stage half.

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

**Who depends on this:** every agent session in every repo that loads the shipped
Spacedock contract — the first-officer/ensign skills, the `status`/`dispatch` launcher,
and the workflow scaffolding. A change here reaches codex, Claude, and pi workers across
spacedock_v1, zaphod, and spacedock_subspace.

**What a defect costs:** a wrong contract clause or a launcher regression derails real
multi-hour sessions at scale — the 0260 forensics corpus records 15 confirmed
derailments, 4 of them multi-hour runaway loops. Silent state or dispatch corruption can
lose a worker's committed work.

**Derived policy:** high rigor for the four high-stakes surfaces (front-door launcher,
`status` mutation/guard paths, shipped contract/scaffolding, CI/release machinery):
behavior tests over prose, the detached adversarial audit before merge, live-drive proof
for contract claims. This is a HIGH-stakes workflow and the default rigor is correct
here. A throwaway spike entity may record a lower local stakes in its own body and
decline disproportionate findings; this declaration sets the project baseline, not an
entity-by-entity floor. Per-stage depth is set by the stage definitions (ideation
spikes cheaply and throwaway; validation applies full depth including the detached
adversarial audit), and the two compose — a spike at ideation is not held to
validation-grade rigor even in this high-stakes workflow.
```

*Add one mechanism note under `docs/dev/README.md` `### Reading sections` (kept OUTSIDE the
`## Stakes` section so it is not itself carried into packets):*

```markdown
The workflow's `## Stakes` section is read by `status --boot` and injected verbatim into
every dispatch packet, so a declared stakes reaches every worker and reviewer. Declare it
once here; never copy it into entity bodies. Its content is a captain-declared fact: an FO
may scaffold the heading, but the declared value changes only with captain sign-off — no
agent infers, escalates, or rewrites it. A workflow with no `## Stakes` section is treated
as: apply the rigor the stage definitions specify and no more; do not infer additional
project stakes; ask the captain to declare one if a different baseline is needed.
```

## Out of scope

- **Setting or grading stakes values.** The contract carries the declaration; it never
  decides a project's stakes or enforces a rigor level from them.
- **The AGENTS.md one-line digest + pointer** (router-layer / `template` group) and **the
  roborev config alignment** (repo-local group). This entity delivers the boot + dispatch-
  packet channels (which reach ensigns and every dispatched reviewer). The canary result
  that shapes those channels is recorded above; wiring them is theirs.
- **Per-entity stakes override.** Stakes is workflow-level here (one `## Stakes` per README).
  An entity-level override is a possible future extension, not this entity.
- **`status --validate` hard-requiring a `## Stakes` section.** Presence is made *visible*
  (the `none declared` marker); hard-fail validation would break every existing workflow
  that predates the field and belongs with the `template`/validate work. Keeping this
  surface minimal and stable is deliberate — downstream groups cite the field.

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

**AC-5 — An undeclared workflow gets an explicit default directive, not silent rigor
inference.**
When a workflow has no `## Stakes` section, both the packet and `status --boot` emit a
`none declared` marker carrying the default behavior (apply the stage definitions' rigor as
written; do not infer additional project stakes; ask the captain to declare). Verified by: a
toggle test that runs `dispatch build` and `--boot --json` against one fixture with and
without the section and asserts the output SWAPS between the verbatim section and the
default-directive marker on the section's presence — the marker appears only in the absent
case and the section only in the present case. Fails if absence yields a bare/empty value, if
the marker leaks into the declared case, or if the directive text is missing.

## Test plan

- **Fixture/CLI (primary, cheap):** the AC-1/AC-2/AC-3/AC-5 Go behavior fixtures in
  `internal/dispatch` and `internal/status` — packet-carries-verbatim, boot-exposes,
  edit-follows, cross-channel identity, stage-invariance (same block at ideation and
  validation, AC-1d), and the presence-toggle that swaps section↔default-directive marker
  (AC-5). These drive the built binary / package functions and assert on generated bytes, not
  on instruction-file text. No new harness, no new infra: they reuse the spike-proven
  `scanHeadings`/section-extractor substrate.
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
