# Sprint: pi-ux

Pi operator and harness experience. Filing and scoping work that makes Pi
first-contact and live runs usable without per-harness hand-wiring.

**Sprint:** the tasks matching `sprint: pi-ux`:

```bash
spacedock status --workflow-dir docs/dev --where sprint=pi-ux
```

The drivable set excludes deliberately deferred candidates:

```bash
spacedock status --workflow-dir docs/dev \
  --where sprint=pi-ux --where 'sprint-readiness != defer'
```

**Target train:** `next` development line. No stable-release tag from this sprint.

## Goal

Make Pi support honest and usable by default: an isolated-home Pi run discovers
its extensions without the harness hard-coding paths; the doctor reports the
truth about an installed Pi; a live lane that stalls leaves enough evidence to
classify the stall; gate presentation and delegated-approval continuation are
reliable; and Pi session identity is scoped to its session.

Membership is owned by the workflow query. This index records direction and
clusters. It never lists members and never tracks their state.

## Carved work

Carved work is discovered from workflow state, never enumerated here. The sprint
groups work in these outcome areas:

- **First-contact defaults** — isolated-home extension discovery, so an operator
  does not export `PI_SUBAGENTS_PACKAGE_ROOT` or hand-wire `--extension` to get a
  Pi run going, and a doctor that reports installed prerequisites truthfully.
- **Live-run operability** — a stalled live run leaves an artifact that classifies
  the stall, with the limits of that classification stated.
- **Gate reliability** — delegated-conn gate continuation completes the
  presentation-to-application boundary on Pi.
- **Identity plumbing** — a session-scoped identity, so one session's install
  state cannot suppress the next session's offer.

## Held directions (not yet backlog entities)

These are recorded as direction only. Each needs to decompose into concrete,
measurable items before it earns a backlog entity. Do not file these as single
unbounded blobs.

- **General usability fixes.** A bucket, not an entity. Decompose into concrete
  friction items as they are found (each its own entity with its own AC). Collect
  candidates here until each has a movable-baseline value statement.
- **Harness-/state-level activity and decision log — deferred from this sprint's
  promise.** No `spacedock` command produces an activity feed or decision log
  today; the `gates:` frontmatter records, `_debriefs/`, and `_evidence/` are the
  closest existing surfaces. A feed or decision-log view needs an end-value, a
  movable baseline, a decision on whether it reads the `gates:` records or a new
  event stream, and a decision on whether it is a new command, a `status`
  projection, or a view over existing state. This sprint does not promise it.

## Inter-sprint boundaries

- Gate work here coordinates with the gate lifecycle and dispatch-after-consume
  owners in other sprints.
- Isolated-home setup is shared with the live-harness owner; coordinate the two so
  their setup does not diverge.
- The Pi journey repairs in `pi-live-completeness` stay there. This sprint does not
  absorb them, and it does not weaken their assertions.

## Definition of done

- An isolated Pi home discovers the `pi-subagents` and `pi-intercom` extensions
  with no package-root variable exported and no hand-wired extension path.
- `spacedock doctor --host pi` reports correct verdicts for the installed Pi
  family, and its remedy text can clear the line it appears under.
- A stalled live Pi run leaves an artifact that classifies the stall, with the
  classification limits stated.
- A repeated Pi gate journey completes presentation through successor dispatch.
- Pi session identity is session-scoped, so a failed install in one session does
  not suppress the next session's offer.
- The `pi-live` lane installs and validates the published Pi family, and its
  substrate assertions derive from the installed package rather than a hard-coded
  source layout.
- Every live Pi XFAIL binding names an active owner and is backed by classified
  evidence.
- `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` pass.
  Changed live surfaces also pass their exact required runtime lanes.

## Notes

- Filed by the FO on 2026-08-13 after a Pi live-test session surfaced the
  extension-discovery friction, the custom-slow-model harness gaps, and
  undiagnosable stalls as concrete operator-experience gaps.
- Rewritten by the FO on 2026-10-03: the member table was removed, because the
  index never lists members; the definition of done was tightened from stages and
  task names to outcomes; and the activity/decision log moved to held directions
  as explicitly deferred.
