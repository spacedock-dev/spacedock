# Validation gate: canonical v1 gate recorder

Recommendation: **APPROVE** commit `024a2c56` for landing.

## Capability

The worktree binary owns one canonical `gates:` projection and exposes four agent-facing forms:

- `spacedock gate record ENTITY --briefing PATH/briefing.json` opens or replaces the current-stage Briefing, or appends a successor after closure.
- `spacedock gate record ENTITY --result FILE --association FILE --actor ID` verifies an exact provider Result against the complete artifact inventory derived from the frozen canonical Briefing, then records its Resolution.
- `spacedock gate record ENTITY --decision approve|revise|hold --actor ID ...` records a chat-rendered Resolution.
- `spacedock gate validate ENTITY` validates and reports the selected record without writing.

Record operations rebuild only the canonical `gates:` subtree under an entity lock, compare the source subtree, validate the rebuilt entity, sync a temporary file, and atomically replace the entity. They do not transition workflow status, dispatch workers, apply decisions, manage the opaque `application` subtree, or copy externally referenced artifact payloads.

## Validation evidence

- The prototype attempt pointers, sequence/lineage/state fields, arbitrary gates extensions, source-position scalar editor, compatibility inference, and prototype replay fixtures are removed. Unsupported pilot encodings fail closed.
- A retained manifest must be named exactly `briefing.json`; another basename fails before lock acquisition or entity mutation.
- Complete Result association is checked against the artifact inventory independently derived from the frozen Briefing JCS digest. A consistently truncated association fails without mutation.
- Canonical open/replace/close/successor behavior, cross-logical-gate re-entry, frozen Resolution/application state, CAS, lock contention, invalid rebuild, atomic replacement, mixed-line-ending preservation outside `gates`, status projection, and unrelated `status --set` coexistence passed.
- Focused gates/CLI/status tests passed. Uncached `go test ./... -count=1` and `go test ./... -race -count=1` passed all 18 packages. The detached validation found no material, deferred, or polish findings.

## Decision

- `approve`: authorize landing and terminal completion.
- `revise`: return concrete material findings to implementation.
- `hold`: retain the validation gate for a named prerequisite.

The entity and final recorder specification are bound as immutable git artifacts in this Briefing rather than copied snapshots. Their artifact revisions are raw SHA-256 pins; the Briefing itself is bound by its RFC 8785/JCS digest when recorded.
