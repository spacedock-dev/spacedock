# Gate review: Git-root Review v1 source materialization — backlog

## Capability

Make a recorder-ready Git-addressed Briefing actually presentable by a real provider
after main and state checkouts move independently, without duplicating selected source
payloads in durable gate state.

## Artifact summary

- **Gate review:** this scope, ownership boundary, evidence requirement, and
  recommendation.
- **Entity:** the complete seed problem, four acceptance criteria, initial test gate,
  and cross-repository options that ideation must compare.

## Why this is separate

s4 can identify and validate
`git-root://<root>/<full-commit>/<repository-relative-path>` plus raw SHA through local
Git objects. Current Subspace cannot render that URI. The missing work is a real
consumer/materialization API and its cross-repository proof—not more gate preparation
or terminal transport.

## Recommendation

Approve entry into ideation. The worker must first exercise the current
Spacedock-room-to-Subspace failure boundary, choose the smallest provider-neutral byte
delivery mechanism, assign each repository's ownership, and declare files/LOC before
implementation.
