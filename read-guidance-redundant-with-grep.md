---
title: Trim the redundant status --read section-read guidance — grep already covers it
status: backlog
source: "FO+captain analysis (2026-06-16), measured on real files: for entity bodies (the actual read target) `status --read --json` and `grep -nE '^#{1,4} '` produce IDENTICAL heading maps (m4 entity 19=19; FO shared-core 18=18). grep over-counts ONLY on fenced markdown-like content (dev README 23 vs 18 — the task-template block). The contract's FO line `first-officer-shared-core.md:214` pairs 'prefer grep, anchor on the heading' WITH 'status --read for offset/lines' — but grep's heading list already yields a section's offset AND its span (the next heading's line); `wc -l` yields the append-point total_lines (ensign `ensign-shared-core.md:92`); `status --resolve`/`--where --json` yield frontmatter. The ensign sites (`:18`, `:92`) name ONLY --read, dropping the grep alternative entirely. So the --read adoption guidance largely re-states the grep-anchor rule the contract already mandates; the sole non-redundant residue is fence-safe heading detection (situational). hf's four FO captures read 0/0 — consistent with re-selling a tool grep already covers, and explains why the trimmed site-6 (4x) was the wrong lever (instruct harder)."
sprint: 0204-structured-reads
sprint-readiness: ready
issue:
id: 82kzghcy3j3cet3hynwa4165
---

## Problem
The `status --read` section-read guidance spans FO `first-officer-shared-core.md:214` and ensign `ensign-shared-core.md:18,:92`. Measured against grep on real files, for the primary read target (entity bodies) the two produce identical heading maps; grep over-counts only on fenced markdown-like content (the README's task-template block). Everything the --read guidance instructs is already available: grep's heading list gives a section's offset + span; `wc -l` gives the append-point total_lines; `status --resolve`/`--where --json` give frontmatter. The FO bullet (214) at least names grep as primary and tacks --read on; the ensign sites (18, 92) name ONLY --read with no grep alternative — the redundancy is concentrated there. Net: the adoption push re-sells what "prefer grep" already covers, which is why adoption read 0/0 in hf's FO captures and why instructing harder (the trimmed site-6) missed.

## What's needed (evidence-first; composes with hf + f5)
- Reduce the --read adoption instruction to its non-redundant residue: fence-safe heading detection where grep over-matches, plus one-call frontmatter+spans as a *convenience* — not a mandate over grep.
- Make FO and ensign consistent: name grep as the primary section-locator in the ensign sites (18, 92) as FO 214 does, or collapse to one shared rule.
- Keep the `status --read` TOOL (correctness-by-construction is a thin but real rationale); trim the INSTRUCTION, not the binary.
- Gate the trim on measured adoption (hf's metric, ensign-transcript-aware per f5): prove --read+scoped-Read usage does not regress after the trim, so this is evidence-driven rather than another assertion.

## Acceptance criteria
- **AC-1** — The --read adoption guidance across the FO + ensign contracts is reduced to its residue and grep is named as the primary section-locator consistently in both, with the dispatch goldens regenerated to the trimmed prompt and a journeymetrics before/after (hf's metric) showing --read+scoped-Read adoption does not regress. Verified behaviorally — golden diff + measured counts — not by a prose-grep over the instruction files.
