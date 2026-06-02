---
id: n3md630s32x3x8zey698xkn6
title: status --set silently drops leading-# frontmatter values — breaks pr-mod merge-detection
status: backlog
source: FO (2026-06-02) — hit LIVE this session: pr-mod recorded pr=#260; status --set wrote it bare; s0's comment-strip read it empty; jg terminalize guard misfired + #258/#259 were invisible to merge-detection
started:
completed:
verdict:
score: "0.32"
worktree:
issue:
---

A 0.19.2 regression in the interaction between two shipped features. `s0` (#254) added the inline-comment strip (a value after `#` is treated as a YAML comment). The pr-merge mod records the PR as `pr=#{N}`. `status --set pr=#260` writes the value **bare** (`pr: #260`) — the writer fails to auto-quote a value that *starts* with `#`. On the next read, the comment-strip eats `#260` → `pr` reads **empty**.

Consequences, all observed live this session:
- `status --where "pr !="` returns `[]` even when entities have `pr` set → the **pr-merge mod merge-detection (startup/idle/event-loop) silently misses every PR-pending entity** → they never auto-terminalize.
- The terminal-transition guard sees `pr` empty → refuses terminalization of a genuinely-merged entity (jg #260 was merged + green, but `status --set completed` was rejected with "pr field is empty").

The AC-4 design (debrief 04) said "the writer auto-quotes `#`-bearing values to round-trip" — but it has a gap: a value that *begins* with `#` (the whole value is a comment token) is not quoted. A value like `foo #bar` (internal space-`#`) may be handled; `#260` (leading) is not.

**Workaround applied to live state (not a fix):** store `pr="#N"` quoted — verified it round-trips (`status --set pr='"#260"'` → `pr: "#260"` → reads `#260`, and `--where pr !=` sees it). jg terminalized, #258/#259 re-quoted.

## Fix direction (ideation hardens)
- Preferred: `status --set`'s frontmatter writer auto-quotes any value YAML would misparse — a leading `#`, a leading `!`/`&`/`*`/`>`/`|`/`%`/`@`, a leading-space, etc. — not only internal space-`#`. The pr-mod's `pr=#{N}` then round-trips unchanged.
- Alternatives to weigh at ideation: have the pr-merge mod record `pr="#{N}"` (quote in the mod); or make the reader tolerate a leading-`#` in known scalar fields. Preferred is the writer fix (it makes ALL `#`-leading values safe, not just `pr`).

## Acceptance criteria (provisional — harden at ideation)

**AC-1 — `status --set field=#value` round-trips.** Setting a field to a leading-`#` value (e.g. `pr=#260`) writes a form that reads back the same value — `status --where "field !="` sees it and a subsequent read returns `#260`, not empty.
Verified by: a status round-trip test setting `pr=#260` then reading it back non-empty; plus a `--where "pr !="` assertion. Must fail against the current bare-write behavior and pass after.

**AC-2 — the pr-mod lifecycle works end-to-end without manual re-quoting.** A merge-boundary `pr=#{N}` record is detected by `--where "pr !="` and does not trip the terminal guard.
Verified by: a guard/merge-detection test over a `pr=#N` entity.

## Notes
- `internal/status` lane (the frontmatter writer/quoting — mutate.go / frontmatter.go). Coordinate with the other status-lane items (#251 merged; terminal-guard-rejected-consistency, 2a).
- 0.19.4-class per FO recommendation (don't delay the 0.19.3 codex/--plugin-dir patch; the workaround holds) — unless the captain folds it into 0.19.3.
