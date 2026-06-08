---
id: xpgjtms50m8705hbw6wthzyx
title: Cross-session FO↔Commander comms — a real channel to replace inbox-file injection
status: backlog
source: "captain + FO probe (2026-06-08, this session). The shaping FO and the driving Commander run as SEPARATE Claude Code sessions with separate teams, and there is no supported cross-session message path. The handoff's open thread: 'Commander handoff is mechanically underspecified — a candidate spacedock affordance.' A live probe proved inbox-file injection works one-way (see prior-art below); this task is the real channel."
started:
completed:
verdict:
score:
worktree:
issue:
---

Give the FO (shaping seat) and the Commander (drive seat) — which run as separate Claude Code sessions — a real, supported way to talk: scope changes, caveats, gate/handoff signals, without poking another session's private state files.

## Problem

The FO/Commander split (FO shapes; Commander drives a separate top-level session) has no message path. Each session is its own team-lead; `SendMessage` only reaches members of the SAME team. So FO→Commander coordination today has no sanctioned transport.

## Prior-art — the proven probe (this session)

Writing a message object directly into the Commander's team inbox file **works**, one-way:

- **Target:** `~/.claude/teams/{team}/inboxes/team-lead.json` (a JSON array).
- **Message schema:** `{from, text, summary, timestamp:"…Z", color, read:false}`. The receiver renders it as `@{from}❯ {summary}` + body, **identical to a native `SendMessage`** — and the receiving Commander correctly classified it as parallel-session intel rather than a team-lead instruction.
- **Write atomically** (temp + rename) to avoid corrupting a concurrent read.

**Why it is a probe, not infrastructure (the requirements this task must satisfy):**
1. **One-way.** The receiver can't reply — the sender is not in the receiver's team roster, so there is no return inbox. A real channel is **bidirectional**.
2. **Race-prone.** Concurrent member-sends can clobber a naive append; there is no lock. A real channel is **atomic / locked**.
3. **Latency.** Surfaces only on the receiver's next inbox poll. Acceptable, but should be defined.
4. **Brittle / unsupported.** It writes another session's private state file; the format can change without notice.

## Proposed approach

Ideation fills this in. A real channel: bidirectional, locked, addressable by session/role (FO ↔ Commander), so cross-session coordination is a first-class affordance rather than a file poke. Prior-art for a bidirectional transport already in this repo: the `telegram` channel skill. Decide whether the channel is a spacedock-owned construct (a `spacedock channel`-style surface) or a thin adapter over an existing transport.

## Out of scope

- Intra-session multi-agent messaging within ONE team — `SendMessage` already covers that.
- The inbox-injection probe as a shipped feature — it is recorded here as prior-art, not the deliverable.

## Acceptance criteria

Ideation/implementation fills in. Sketch:

- The FO and a separate-session Commander exchange a message **both directions** (verified by a two-session live exercise: FO sends, Commander receives and replies, FO receives — observed in durable channel state / both transcripts, not by reading one inbox file).
- Concurrent sends do not lose a message (verified by a test that interleaves writes and asserts all are delivered).

## Test plan

Ideation/implementation fills in. The riskiest unknown is the bidirectional + locked transport — exercise the smallest round-trip first (FO↔Commander send/reply across two real sessions) before building the surface around it.
