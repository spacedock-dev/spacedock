# 0.20.3 boot-forensics verdict

Post-`v0.20.3` boot-and-gate token-growth profile (a cold FO boot that drove to a gate presentation, measured under the 0205 sprint's `w4` entity). The verdict on the 0203 fo-efficiency boot-cost goal: **the remaining boot cost is structural, not avoidable** — the FO contracts + skill payloads (load-once, boot-resident) and the gate-review output dominate, both correct behavior.

All figures in tokens (context-window occupancy = input + cache_read + cache_creation).

**Headline:** cold baseline 22,536 tok → 59,200 tok, ×2.63 growth over 14 turns.

---

### Per-stage table

| Stage                         | Turns | What was added                                                                                              | Ctx after (tok) | Δ (tok) |
| ----------------------------- | ----- | ----------------------------------------------------------------------------------------------------------- | --------------- | ------- |
| Pre-loaded context            | —     | Handoff file, CLAUDE.md, gitStatus, system prompt, skill list                                               | 22,536          | —       |
| FO contract load              | 1–3   | first-officer-shared-core.md (Read t2), claude-first-officer-runtime.md (Read t3)                           | 32,310          | +9,774  |
| Startup — discovery           | 4–7   | `--version`, `git rev-parse`, `--discover`, README frontmatter head-60, `git pull`, `--boot --json`         | 36,010          | +3,700  |
| Status + 0205 context         | 8–9   | Full status table, `ls 0205`, roadmap README head-100, `0205/index.md`                                      | 38,910          | +2,900  |
| w4 entity inspection          | 10–11 | `spike-prelim-findings.md`, `w4 --read --json` headings, w4 AC+stage-report section (offset 180, 130 lines) | 47,784          | +8,874  |
| Gate presentation + forensics | 12–14 | present-gate skill, 2,424-tok gate review output re-entered, boot-forensics                                 | 59,200          | +11,416 |

---

### Top single-turn contributors (Δ, labeled by the prior tool/file that caused the jump)

| Rank | Turn | Δ (tok) | Prior tool / cause                                                                                        |
| ---- | ---- | ------- | --------------------------------------------------------------------------------------------------------- |
| 1    | t3   | +8,880  | Read of `first-officer-shared-core.md` (~218 lines of dense contract prose) + Bash CLAUDECODE + t2 output |
| 2    | t13  | +3,758  | Gate review text (2,424 output tokens from t12) re-entering context                                       |
| 3    | t10  | +6,877  | Read of `0205/index.md` (91 lines) + roadmap README head-100 Bash + t9 output                             |
| 4    | t14  | +3,354  | boot-forensics Bash + t13 output (gate review still live)                                                 |
| 5    | t12  | +4,304  | Read of w4 ACs+stage-reports (offset 180, 130 lines) + present-gate skill load                            |

*(ranked by raw Δ; turns 10 and 12 appear lower by rank but are second- and third-largest single-turn jumps)*

---

### Efficiency findings

**Nothing materially avoidable.** Every read was acted on:

- **shared-core.md** (+~8k) is the boot-resident contract; it loads every session. The filed token-cleanup proposal (~638 recoverable tokens) is the right lever — not session-level avoidance.
- **0205/index.md** (91 lines, full read): the handoff explicitly required it and the gate review consumed goal/DoD/dependencies densely. A section-scoped read would have saved ~1–2k but needed multiple round-trips.
- **Gate review output (2,424 tok)** re-entering at t13 is structural — a 15-25-line gate review at the correct fidelity generates this.
- **roadmap README head-100**: partial read (first 100 lines) already; a Grep for `## The two roles` + `## Sprint lifecycle checklist` would have cut this to ~500 tok, saving ~400-600 tok. Low-value optimization given the rest of the profile.
- **Pre-loaded context (22,536)**: the largest single block is the system-resident load (handoff file, CLAUDE.md, gitStatus, skill list). This is the FO's pre-boot floor and is not FO-controllable.

**The dominant cost is structural:** 38% of final context is the two FO contracts + skill payloads (boot-resident, load-once), and 15% is the gate review output. Both are correct behavior for a boot-and-gate session.
