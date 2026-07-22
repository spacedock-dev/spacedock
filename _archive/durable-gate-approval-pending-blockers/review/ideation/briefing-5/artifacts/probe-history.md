# What you asked

> Can this probe-and-reanswer mechanism operate without Spacedock, and where are the
> persistent Probe and each Briefing-bound ProbeResult serialized?

## Current answer

Yes. The mechanism operates with an ordinary Review & Gate Briefing and requires no
Spacedock repository, gate, stage, attempt, entity, or First Officer.

The review provider owns the Probe and its immutable ProbeResults in the spec-lineage
room. The concrete serialization is an append-only `probes.jsonl`: one `Probe` record
holds the stable question identity and revision, and a separate `ProbeResult` record is
appended for each exact Briefing id and digest. An equivalent durable provider store is
valid. These records are not Spacedock entity frontmatter; Spacedock may retain only an
opaque room reference.

## What changed across the evolving design

- **Briefing 1 — insufficient evidence.** It described a Git-backed Subspace room but
  did not specify standalone operation or a concrete Probe/ProbeResult serialization.
- **Briefing 2 — answered.** It defined provider-owned `probes.jsonl` and an ordinary
  Review & Gate flow with no Spacedock dependency.
- **Briefing 3 — limitation removed.** One opaque room reference now spans the whole
  spec lineage instead of being confused with attempt identity.
- **Briefing 4 — concern still holds.** Only the Spacedock prepared-application
  lifecycle changed; the Probe companion artifact is byte-identical to Briefing 3.

## Comparison friction observed

The current mechanical comparator labels Briefing 3 → Briefing 4 as `changed` because
the fresh responder used different wording, citations, and limitations. Semantically,
the concern was unaffected. This is a live counterexample to treating byte-different
answers as concern drift.

## Evidence

- `gate-review-probes.md`, “Provider-owned serialization”
- `gate-review-probes.md`, “Ordinary Review & Gate flow”
- `gate-resolution-frontmatter-contract.md`, lineage-level `room-ref` field contract

## Limitations

- This is design evidence; no provider persistence or restart behavior was executed.
- Subspace TUI does not yet discover room-level `probes.jsonl`; this artifact is a
  presentation projection so the saved concern is visible in today's TUI.

Answered by Codex · GPT-5 family · fresh runs bound to each Briefing
