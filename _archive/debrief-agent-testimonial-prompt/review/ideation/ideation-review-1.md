# Ideation gate review — debrief-agent-testimonial-prompt (qd)

- Entity: `debrief-agent-testimonial-prompt` (`qdb1w5r7k9nvbvkf8qetcd5m`, short-id `qd`)
- Stage: ideation (first formal gate binding; two pre-gate captain feedback rounds already incorporated)
- Reviewer: agent:first-officer (Pi host), 2026-08-01
- Verdict: **approve** — accept the cycle-2 design package and enter implementation in a worktree.

## Design under review

One mechanism: insert a testimonial-collection step at the debrief flow's Phase 3 (before draft presentation) and a `## Agent Testimonial` section with provenance fields in the produced debrief file, in `skills/debrief/SKILL.md`. The prompt requires the driving agent to self-identify (harness/runtime, model, exact version/build, `unknown` when unverifiable, never inferred) and carries the honesty clause near-verbatim ("not a request for praise"). Session scale (tasks/workers/PRs) is derived from session data, separately from the agent's identity response. A concrete before/after diff against the current skill text is recorded in the entity body.

## Necessity-first check (staff-review principle)

- Should this mechanism exist at all: yes — captain-commissioned 2026-07-02; testimonial evidence currently evaporates unless the captain remembers to ask verbally. The change edits one existing skill flow; it adds no new process, surface, or machinery.
- Could the same end be reached with less: the alternative (convention in captain memory) is what produced the gap this task closes.
- The design's own text contains no mechanism beyond the prompt insertion, the template section, and the provenance fields.

## Acceptance criteria ↔ test plan map

AC-1↔test 1 (produced debrief artifact carries exactly one testimonial section with friction content), AC-2↔test 1/2 (prompt observed in a run transcript, honesty clause near-verbatim), AC-3↔test 2 (two artifact-backed cases: verifiable identity, and version/build unknown-and-preserved; session scale derived independently), AC-4↔test 3 (split-root `_debriefs/` path + log check), AC-5↔test 5 (required live lane green at merge: `claude-live` first pass per the path→lane rule over `skills/**`), plus test 4 (go test regression gate). Every AC has behavior-level evidence required of the implementation stage; none is prose-only.

## Riskiest mechanism

Behavioral model compliance (honest `unknown`, honesty-clause effect) — exercised pre-gate by the recorded Codex prompt exercise in the entity body (four-sentence testimonial, explicit friction, honest `unknown` build). Verification of the shipped prompt is artifact-backed (run transcript + produced debrief), not a static grep of the skill file. "No spike needed" is recorded and justified by exercised mechanisms.

## Feedback cycles incorporated

1. Cycle 1 (2026-07-02): register/length guidance (person-voice, plain nouns, no internal workflow jargon in the produced section), honesty clause load-bearing, provenance label form (`date · runtime · model · N tasks, N workers, N gates, N PRs`). Guidance for implementation rendering; ACs unchanged by design.
2. Cycle 2 (2026-07-18): agent self-identification moved into the prompt itself; `unknown`-never-guess codified; the cycle-2 stage report's three DONEs map onto it.

## Scope notes

- No docs-site change in this pass (no page duplicates the debrief template); recorded in the body.
- A Pi-runtime adapter observation from cycle-1 feedback (ensign agent-type binding) is explicitly out of scope, recorded upstream; not fabricated into ACs.

## Recommendation

Approve. The package is necessity-checked, AC↔test-mapped, spike-justified, and carries its merge gate (`claude-live`) in the criteria. Approval opens implementation in worktree `.worktrees/spacedock-ensign-debrief-agent-testimonial-prompt`; expected downstream fan-out: one implementation worker, one fresh validation worker (tolerance ±1 for correction rounds).
