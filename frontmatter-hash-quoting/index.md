---
id: n3md630s32x3x8zey698xkn6
title: Frontmatter writer doesn't quote YAML-ambiguous values (leading-#, internal colon) — breaks pr-mod + GitHub render
status: backlog
source: "FO (2026-06-02) — pr=#N eaten by s0's comment-strip (live: jg #260, #258, #259) AND captain-observed GitHub YAML parse errors on colon-bearing source values"
started:
completed:
verdict:
score: "0.34"
worktree:
issue:
---

The frontmatter writer (`status --set`, and the seed/Write path) does not quote values that strict YAML would misparse. The line-hack status parser tolerates them, but strict YAML readers — GitHub's frontmatter renderer AND s0's own inline-comment strip — do not. Two classes observed live this session:

1. **Leading `#` → comment.** The pr-merge mod records `pr=#{N}`; `status --set pr=#260` writes it bare (`pr: #260`); s0's (#254) comment-strip then reads `#260` as a YAML comment → `pr` reads **empty**. Consequences: `status --where "pr !="` returns `[]` so the **pr-mod merge-detection silently misses every PR-pending entity**, and the terminal guard refuses a genuinely-merged entity ("pr field is empty"). Hit on jg #260 (terminalize blocked) and #258/#259 (invisible to merge-detection until re-quoted).
2. **Internal `: ` (colon-space) → mapping.** A free-text value like `source: captain — reproduced LIVE: codex ...` parses, under strict YAML, as a nested mapping (`LIVE:` becomes a key) → GitHub renders the frontmatter as **"object not allowed"** (captain-observed on the state-branch entity views / PR audit links). The `source:` field is the prime victim (it routinely contains colons).

The AC-4 design (debrief 04) claimed "the writer auto-quotes `#`-bearing values to round-trip," but the gap is broader: it covers neither *leading*-`#` nor internal colon-space, and the same applies to other YAML indicators.

**Workaround applied to live state (not a fix):** store quoted values (`pr="#N"`, `source="..."`). This session: re-quoted #258/#259 `pr`, terminalized jg, and quoted the `source:` of the three active colon-bearing entities (codex-plugin-list-parse-drift, orphan-status-drift, this one).

## Fix direction (ideation hardens)
- Preferred: the frontmatter writer auto-quotes any value strict YAML would misparse — a leading `#`/`!`/`&`/`*`/`>`/`|`/`%`/`@`/`` ` ``, a leading or trailing space, AND an internal `: ` (colon-space) — for every field, not only `pr`.
- Plus a one-off **re-quote migration** of existing affected entities (active AND archived) — many archived `source:` values carry colons and render broken on GitHub today.
- Alternatives to weigh: have the pr-merge mod record `pr="#{N}"` (narrow); reader-tolerate (doesn't fix GitHub render). Preferred is the writer fix — it makes ALL fields safe and is the only option that fixes the GitHub-render half.

## Acceptance criteria (provisional — harden at ideation)

**AC-1 — leading-`#` values round-trip.** `status --set field=#value` (e.g. `pr=#260`) writes a form that reads back the same value; `status --where "field !="` sees it; the pr-mod merge-detection + terminal guard work without manual re-quoting.
Verified by: a round-trip + merge-detection test over a `pr=#N` entity, failing against the current bare-write and passing after.

**AC-2 — colon-bearing values are GitHub-YAML-parseable.** A `source:` (or any field) value containing `: ` is written quoted so a strict YAML parser reads it as a single scalar (no "object not allowed").
Verified by: a strict-YAML parse test over a written entity with a colon-bearing field; assert it parses as a scalar.

**AC-3 — existing affected entities are migrated.** A one-off pass re-quotes the colon/`#`-bearing frontmatter values across active + archived entities so the state branch renders cleanly on GitHub.
Verified by: a scan asserting no active/archived entity has an unquoted colon-space or leading-`#` frontmatter value.

## Notes
- `internal/status` lane (the frontmatter writer/quoting — mutate.go / frontmatter.go). Coordinate with the other status-lane items (terminal-guard-rejected-consistency, 2a).
- 0.19.4-class per FO recommendation (don't delay the 0.19.3 codex/--plugin-dir patch; the workaround holds) — unless the captain folds it in.
