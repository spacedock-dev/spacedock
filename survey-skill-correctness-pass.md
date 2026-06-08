---
id: vh5rzxexn9wc7dedex7cepzy
title: Survey skill correctness pass — agentsview git-root model fix + codex-cwd + scaffold-fact + sandbox probe (consolidates 69/1p/4t)
status: validation
source: "captain (2026-06-08) — preflight B2 + consolidation directive. The survey members 69/1p/4t all edit skills/survey/SKILL.md and share one agentsview model; 69's spike disproved the cwd-basename keying the survey teaches. Consolidate into ONE coherent survey-skill pass on the corrected git-root-basename model."
score: "0.3"
started: 2026-06-08T17:02:54Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-survey-skill-correctness-pass
issue:
sprint: 0198-pre-flip-hardening
group: survey
sprint-readiness: ready
mod-block: merge:pr-merge
---

One coherent correctness pass over the survey skill (`skills/survey/SKILL.md` + `references/queries.sql`), on a single corrected agentsview model. Consolidates the three coupled survey members (69 codex-cwd, 1p scaffold-fact, 4t sandbox-probe) plus the preflight's B2 model fix — they all edit one file, share the agentsview model, and bottom out on one live drive, so they ship as one task.

## Deliverables

1. **agentsview model fix (preflight B2 — the new piece).** 69's spike proved (and this ideation re-confirmed against real data) that agentsview v0.32.1 keys `project` by GIT-ROOT basename (this repo's root, every `.worktrees/*`, and `.spacedock-state` all key to ONE project), contradicting `SKILL.md:64` + the `queries.sql` `scoping`/`folded_keys` rationale (which assume per-cwd-basename divergence — the model `xn` shipped in 0.19.7). Ideation resolves the two open questions: `folded_keys` is structurally dead (always 1) and is DROPPED; the cwd-prefix-union scoping STAYS but its rationale inverts (it excludes same-basename sibling repos, not merges divergent keys). See "The corrected agentsview model" below.
2. **codex-cwd presence + hint (absorbs 69).** Codex sessions land `cwd=''` (agentsview doesn't persist Codex cwd) → the cwd-scoped query misses them. A separate flagged codex-presence count + report hint (NOT a silent project union — 69's collision evidence: workspace→6 roots, spacedock→3 roots). See 69's ideation report.
3. **SCAFFOLD state-the-fact (absorbs 1p).** Drop the recovered-vs-installed taxonomy (`SKILL.md` §3); state the observed scaffold fact, keeping the "recovered from behavior, not files" fact.
4. **sandbox install-detect probe swap (absorbs 4t).** Replace `command -v agentsview` (`SKILL.md:27`, a stat/access builtin Seatbelt denies) with an `agentsview --version` exit-code probe (an execve Seatbelt allows). See 4t's ideation report (root-caused).

## Absorbed entities

Supersedes `survey-codex-cwd-workaround` (69), `survey-scaffold-state-the-fact` (1p), `survey-agentsview-detect-under-sandbox` (4t) — each deferred with `superseded-by` = this. Their ideation reports (designs + spikes) are the inputs for deliverables 2/3/4; this task's ideation adds deliverable 1 (the §64 model fix) and weaves all four into one coherent SKILL.md/queries plan.

## Proof

Survey skill (`skills/survey`) — built by a worker under test; proof = the query-smoke (`queries.sql` executed against the production-shaped fixture) + a LIVE DRIVE exercising all observable changes (codex-presence hint, the corrected scoping with `folded_keys` dropped, the reshaped SCAFFOLD, the sandbox probe non-prompt) in ONE pass. Not a high-stakes surface → normal validation. Per the survey discipline, a grep over SKILL.md never satisfies a behavioral AC.

## The corrected agentsview model (deliverable 1 — the spine the other three hang off)

**The fact.** agentsview v0.32.1 derives `sessions.project` from the **git-root basename** (normalized: non-alphanumerics → `_`), NOT the cwd basename. Every checkout of one repo — the root, every worktree (`.worktrees/*` and `.claude/worktrees/*`), and the `docs/dev/.spacedock-state` split-root checkout — keys to the SAME `project`. Confirmed against real data (see Spike), reusing 69's git-root-basename finding plus a confirming probe in THIS ideation: 29 distinct spacedock-v1 cwds → ONE project key `spacedock_v1`; no spacedock-v1 subpath ever keyed to anything else; no single cwd ever mapped to >1 project.

**What this breaks.** Three places teach the disproven cwd-basename-divergence model:
- `SKILL.md:64` — "agentsview keys each session by the basename of its working directory, so a subdir checkout, a worktree, or the split-root state dir each get a DIFFERENT `project` key." FALSE: they get the SAME key.
- `queries.sql:18-21` (the `:repo_root` rationale) — "a subdir checkout and a worktree-style path coalesce to ONE repo identity even though agentsview keyed them under divergent basename-derived `project` columns." The coalescing is real; the *reason given* (divergent keys to merge) is false.
- `queries.sql:30-34` (the `scoping` comment) — "regardless of the divergent basename-derived `project` agentsview assigned a subdir/worktree checkout … the prefix union recovers the subdir + worktree sessions the basename key drops."

**What is STILL correct (do NOT change the mechanism).** The cwd-prefix-union (`cwd = :repo_root OR cwd LIKE :repo_root || '/%'`) stays. Under the corrected model its job inverts but is still load-bearing: `project` is git-root-*basename*, so it COLLIDES across unrelated repos that share a basename (69's real evidence: `spacedock` → 3 different roots, `workspace` → 6). A `project = :repo_project` scope would therefore fold a *sibling* repo's sessions into this one. The cwd-prefix is the ONLY column that distinguishes "this repo's worktrees/subdirs" (all under the absolute repo-root prefix) from "a same-basename sibling elsewhere on disk." So the scoping is right; only its *rationale prose* is wrong. This same collision is exactly why deliverable 2 (Codex, blank cwd) must be a flagged presence count, not a project union.

**`folded_keys` is now structurally dead.** `COUNT(DISTINCT project)` over a cwd-prefix-scoped Claude set is **always 1** for a git repo (every in-prefix cwd shares the git-root key) and trivially ≤1 for the non-git fallback (`REPO_ROOT=$(pwd)`, one cwd). It can never be >1, so the `folded_keys > 1` branch (`SKILL.md:96`, report line `SKILL.md:143`) is unreachable. **Decision: drop `folded_keys` from the `scoping` SELECT and remove the `{if folded_keys>1: coalesced from {folded_keys} keys}` report line + the §96 prose that interprets it.** Keep `blank_cwd` (still meaningful for Claude — 0.5% of Claude sessions have blank cwd per 69's spike — and the honest-accounting hook §98/§208 still wants it). This is a real shape change the query-smoke pins (AC-1).

## Deliverable details (exact before/after)

**D1 — model fix.** Rewrite `SKILL.md:64` to the git-root-basename fact: agentsview keys every checkout of one repo (root, worktrees, split-root state dir) to ONE `project` = the git-root basename, so the risk is the INVERSE — that key COLLIDES with same-basename sibling repos; resolve the repo root and scope by the absolute cwd-prefix so siblings stay out. Rewrite the `queries.sql` `:repo_root` rationale (`:18-21`) and the `scoping` comment (`:30-34`) the same way: the prefix excludes same-basename siblings, it does not merge divergent keys. Drop `folded_keys` (above). Net: prose tells the true story; mechanism unchanged.

**D2 — codex-presence count + hint (absorbs 69).** Add a labeled `-- name: codex-presence` query to `queries.sql`, bound to a new `:repo_project` param, returning `COUNT(*)` of `agent='codex'` sessions whose `project = :repo_project` and `SUM(cwd='')` among them. SKILL.md step 2 derives `:repo_project` from `REPO_ROOT`'s basename with the `-`/non-alphanumeric → `_` normalization agentsview applies (e.g. `spacedock-v1` → `spacedock_v1`), runs the query, and step 4 renders a hint line ONLY when the count > 0: "N Codex sessions match this repo by project NAME only (agentsview does not record Codex cwd) — may include a same-named sibling repo; the Claude-scoped body above does not cover them." NOT a union into the Claude `scoping`/scaffold/work-by-area sets (collision risk above). Codex decision/scaffold/work-by-area surfacing stays the documented deferred follow-up (`SKILL.md:13,96,209`).

**D3 — SCAFFOLD state-the-fact (absorbs 1p).** Drop the recovered-vs-installed-vs-active 3-bucket *taxonomy* (`SKILL.md` §3 step, and the §3-driven SCAFFOLD report block). State the observed fact directly: per scaffold family, its invocation count and whether it is checked in on disk. The "recovered from behavior, not files" phrasing IS the fact to state (1p's desired-output example), not taxonomy overhead. The file-probe + behavioral-tally JOIN stays (both signals are still read); only the bucket-label classification framing is removed in favor of the plain factual statement. Desired report shape (1p's example):
> `superpowers was recovered from behavior, not files: 186 skill invocations, but no checked-in .claude/skills/superpowers. Other recovered one-offs: plan-writing, using-git-worktrees, systematic-debug, simplify, debugging.`

**D4 — sandbox install-detect probe swap (absorbs 4t).** `SKILL.md:27`:
```
- if ! command -v agentsview >/dev/null; then echo "AGENTSVIEW MISSING"; fi
+ if ! agentsview --version >/dev/null 2>&1; then echo "AGENTSVIEW MISSING"; fi
```
`command -v` is an FS-access PATH-walk (`access()`/`stat()`) Seatbelt can deny while allowing `execve`; `agentsview --version` is an `execve` that survives whatever the survey's real through-the-binary reads survive. `2>&1` suppresses both the present banner and the absent "not found" so only the `AGENTSVIEW MISSING` sentinel reaches the agent. Contract preserved (silent ⇒ present; sentinel ⇒ absent). Install fallback (`SKILL.md:30`) unchanged.

## Step-4 report-fence composition (where 69 and 1p overlap)

Both touch the step-4 fence (`SKILL.md:141-169`). Resolved together so the build worker has ONE target shape:
- The PROJECT header line loses its `{if folded_keys>1: coalesced from {folded_keys} keys}` clause (D1) and KEEPS the `{if blank_cwd>0: · {blank_cwd} uncaptured-cwd sessions}` clause.
- A NEW conditional hint line under PROJECT renders only when `codex-presence > 0` (D2).
- The SCAFFOLD block becomes the plain state-the-fact statement (D3), not the bucket list.

## Spike: no new spike needed

**Mechanism de-risked.** The riskiest unverified path (D1, the model the other three hang off) was exercised in THIS ideation against real agentsview v0.32.1, driven through the binary (raw `~/.agentsview` reads are TCC-denied, per `SKILL.md:21`): synced THIS repo's Claude sessions across root + a worktree + the `.spacedock-state` checkout into a readable `AGENTSVIEW_DATA_DIR`, then queried. Result: 29 distinct spacedock-v1 cwds → ONE `project='spacedock_v1'`; `folded_keys=1` for the real scoping query; no in-prefix cwd ever keyed to another project; cwd→project is a function. D2's premise (Codex `cwd=''`, populated git-root-basename `project`, cross-repo basename collision) and D4's premise (FS-access vs execve asymmetry under Seatbelt) were each de-risked by 69's and 4t's spikes respectively (recorded in their now-superseded reports) and the D4 probe contract was re-confirmed here (present → silent exit 0; absent → `AGENTSVIEW MISSING`, underlying exit 127). **No spike needed: D1 rests on the git-root-basename finding confirmed here (and originally by 69); D2 on 69's spike; D3 composes already-proven query reads; D4 on 4t's Seatbelt-asymmetry spike.** The build's first tests are the fixture rows AC-1/AC-2 seed from these shapes.

## Acceptance criteria

- **AC-1 — corrected `scoping` shape: `folded_keys` dropped, `blank_cwd` kept, cwd-prefix scope unchanged.** `queries.sql`'s `scoping` query returns `sessions | blank_cwd | span` (no `folded_keys` column); the cwd-prefix-union (`cwd = :repo_root OR cwd LIKE :repo_root || '/%'`) is unchanged. Verified by: the query-smoke (AC-2) asserts the corrected column shape and that the prefix still coalesces the in-repo cwds while excluding the out-of-prefix and blank-cwd sessions — fails if `folded_keys` survives in the SELECT, if the column count is wrong, or if the prefix scope regresses. Expected values come from the fixture rows (independent of skill prose).

- **AC-2 — query-smoke over a production-accurate fixture proves Codex-by-`project` is counted, blank-cwd Codex rows are surfaced, and no Codex row is silently folded into the Claude scope.** `skills/integration/testdata/survey/fixture-sessions.sql` is corrected to the git-root-basename model (the in-repo Claude sessions all carry ONE `project` matching the repo basename, NOT the disproven divergent keys `_spacedock_state`/`feature_x`) and gains `agent='codex'` rows with the production shape (`project` set, `cwd = ''`): at least one matching `:repo_project` (in-repo Codex) and one with a same-`project`-different-implied-root sibling shape. `survey_queries_test.go` asserts: (a) the new `codex-presence` query returns the matching Codex count with `blank_cwd > 0`; (b) `scoping.sessions` is UNCHANGED by the added Codex rows (Codex stays out of the Claude scope — no union); (c) the corrected scoping column shape (AC-1). Verified by: `go test ./skills/integration/ -run TestSurveyQuerySmoke` passes. Expected values come from the fixture rows.

- **AC-3 — the sandbox install-detect probe uses invocation exit code, and emits only the sentinel.** `SKILL.md` step 1 probes presence with `agentsview --version >/dev/null 2>&1` (exit-code semantics), never `command -v`/`which`/`test -x`/`stat`/PATH-walk. A deterministic probe-behavior test runs the exact step-1 one-liner twice — once with a stub `agentsview` on a synthesized PATH (expect empty output, exit 0), once with the name absent from PATH (expect sole line `AGENTSVIEW MISSING`) — asserting captured stdout+stderr against those two independent fixture conditions. Verified by: that probe-behavior test passes (the gate-able half); the behavioral no-false-negative claim is AC-5's sandboxed live drive. Fails if the probe reverts to an FS-access form or leaks `--version` output.

- **AC-4 — the live survey report renders the corrected observable surface.** Driving `/spacedock:survey` end-to-end on a repo whose survey DB carries project-matched blank-cwd Codex sessions produces a report that: (a) shows the Codex-presence hint line stating the count and that the match is by project name only (cwd unrecorded); (b) shows the SCAFFOLD section as the state-the-fact statement (family + invocation count + checked-in presence/absence) with NO recovered-vs-installed-vs-active taxonomy labels; (c) shows the PROJECT line WITHOUT a "coalesced from N keys" clause (folded_keys gone). A second drive against a Codex-free DB omits the hint line. Verified by: the live-drive rendered report (the survey's proof bar); a grep over SKILL.md does NOT satisfy this — only the rendered output does.

- **AC-5 — under sandbox with agentsview installed, the survey detects it present and does NOT prompt to install.** Driving `/spacedock:survey` under the sandbox where the pre-fix `command -v` false-negatives, the run proceeds past step 1 to the sync/scan without emitting the `brew install --cask agentsview` consent prompt. Verified by: a sandboxed live drive observing the survey continue without the install prompt, with the pre-fix `command -v` probe confirmed to false-negative in the same sandbox as the negative control (establishing the drive exercises the failing path, not a vacuous pass).

## Test plan

- **Query-smoke (AC-1, AC-2) — fixture, cheap (~seconds), reuses the existing sqlite3-driven harness.** Correct `fixture-sessions.sql` to the git-root-basename model: the three in-repo Claude sessions all key to one `project` matching the repo basename (replacing the divergent `_spacedock_state`/`feature_x` keys), and drop the fixture comment lines that assert divergent keys. Add blank-cwd Codex rows (`project` set, `cwd=''`): one in-repo (`project = :repo_project`) and one same-basename sibling shape. Update `survey_queries_test.go`: the `scoping` assertions move from `sessions|folded_keys|blank_cwd|span` (4 fields, folded_keys=3) to the corrected `sessions|blank_cwd|span` (3 fields) and assert the prefix still counts 3 in-repo sessions; add the `codex-presence` assertions (matched count, blank_cwd>0) and the `scoping.sessions` unchanged-by-Codex assertion. Cost: LOW.
- **Probe-behavior test (AC-3) — deterministic, the gate-able half.** Extract/run the exact step-1 probe line against two synthesized PATH conditions (stub binary present; name absent), assert captured output. If the survey harness is shell-fixture-based, add it there; otherwise a focused shell/Go test under the survey integration testdata. Expected values (`""`, `AGENTSVIEW MISSING`) come from the fixture conditions. Cost: LOW.
- **Live drive (AC-4, AC-5) — minutes; the ONLY proof for the rendered report and the sandbox non-prompt.** ONE end-to-end pass exercises all observable changes at once: seed a survey DB (via the spike's `agentsview session sync` / scoped `agentsview sync` into a readable `AGENTSVIEW_DATA_DIR`) with project-matched blank-cwd Codex sessions, drive `/spacedock:survey` under sandbox, observe (a) the Codex hint, (b) the state-the-fact SCAFFOLD, (c) no folded-keys line, (d) no install prompt. Then a second short drive against a Codex-free DB to confirm the hint is absent (AC-4 negative), and the pre-fix `command -v` confirmed false-negative in the same sandbox (AC-5 negative control). Needs the agentsview binary (present, v0.32.1) and the sandboxed harness.
- **Estimated total cost/complexity: LOW–MEDIUM.** Prose edits across `SKILL.md` (§1 probe, §2 scoping/model, §3 SCAFFOLD, §4 fence) + one labeled query + a `folded_keys` drop + fixture correction + Codex fixture rows + test assertions + one probe-behavior test + one combined live drive. No agentsview/binary/Go-runtime changes beyond the test files (agentsview changes are explicitly out of scope).

## Out of scope

- Fixing agentsview to persist Codex `cwd` (upstream `kenn-io/agentsview`).
- Surfacing Codex decision / scaffold / work-by-area signals (broadening the `agent='claude'` filter across the body) — the documented deferred follow-up; this pass adds only the codex-presence count + hint.
- The install-instruction path (`SKILL.md:30`) — unchanged; D4 fixes detection only.

## Stage Report: ideation

- DONE: Absorb + integrate the three superseded designs (69 codex-cwd, 1p scaffold-fact, 4t sandbox-probe) into ONE coherent skills/survey/SKILL.md + references/queries.sql plan — including the step-4 report-template fence where 69 and 1p overlap.
  Read all three ideation reports; wove D2 (codex-presence count + hint), D3 (SCAFFOLD state-the-fact), D4 (`agentsview --version` probe) into "Deliverable details" + the "Step-4 report-fence composition" section resolving the 69/1p overlap (folded_keys clause dropped, codex hint added, SCAFFOLD reshaped) into ONE target fence shape.
- DONE: Firm the B2 model fix (the new piece): correct SKILL.md:64 + the queries.sql scoping/folded_keys rationale to the git-root-basename reality; decide whether folded_keys stays meaningful and whether the cwd-prefix-union is still right.
  Recorded in "The corrected agentsview model" + D1. Decisions: `folded_keys` is structurally always 1 → DROP it (column + report line + §96 prose); the cwd-prefix-union STAYS but its rationale inverts (excludes same-basename sibling repos rather than merging divergent keys, since the keys are NOT divergent). Reused 69's git-root-basename finding AND re-confirmed with a fresh probe (Spike section) — 29 distinct spacedock-v1 cwds → ONE `project='spacedock_v1'`, real scoping query `folded_keys=1`.
- DONE: Produce the unified build-ready ACs + test plan: query-smoke (queries.sql vs production fixture incl. cwd='' Codex rows) + the deterministic sandbox-probe test + ONE live-drive pass exercising all observable changes. No tautological SKILL.md grep as a behavioral AC.
  Five ACs (AC-1 corrected scoping shape, AC-2 query-smoke incl. blank-cwd Codex rows + no-union, AC-3 deterministic probe-behavior test, AC-4 live-drive rendered report, AC-5 sandboxed non-prompt with `command -v` negative control) + a 4-bullet test plan. Each AC's proof is a test/command/rendered-output outside the task body; the two text-half changes (probe wording, model prose) are explicitly NOT standalone behavioral ACs — the behavioral half is the live drive.

### Summary

Consolidated the three coupled survey members plus the preflight's B2 model fix into one coherent ideation on the corrected git-root-basename agentsview model. The B2 fix is the spine: re-confirmed against real agentsview v0.32.1 data (29 spacedock-v1 cwds → one `project` key) that `folded_keys` is structurally always 1 (DROP it) while the cwd-prefix scoping is still load-bearing for the inverse reason (it excludes same-basename sibling repos, which is also why Codex's blank-cwd sessions get a flagged presence count, not a silent project union). No new spike needed — D1 re-confirmed here, D2/D4 rest on 69's and 4t's recorded spikes, D3 composes proven reads. ACs pin the corrected query shape and Codex no-union via query-smoke, the probe via a deterministic present/absent test, and the rendered report + sandbox non-prompt via ONE combined live drive.

## Stage Report: ideation (cycle 2)

- DONE: Lock the corrected agentsview model + record the proven mechanism.
  Re-confirmed by EXERCISING, not re-reading: synced this repo's Claude sessions through `agentsview v0.32.1` into a readable `AGENTSVIEW_DATA_DIR` (narrow farm captured root + 7 `.worktrees/*` + `.spacedock-state`) and queried the live DB. Result is STRONGER than the body's recorded 29: now 34 distinct in-prefix cwds → ONE `project='spacedock_v1'`; real `scoping` query `folded_keys=1`; the `cwd→project` HAVING-COUNT(DISTINCT project)>1 check returns empty (cwd→project is a function). Basename `spacedock-v1` → `spacedock_v1` (non-alnum→`_`) confirms D2's normalization. The body's "29" is the honest record of the original spike; the conclusion is unchanged (drift over time, same key). No spike needed (re-affirmed): D1 re-confirmed here; D2 — verified Codex `cwd=''` directly (591/591 Codex sessions blank-cwd in the live DB; gemini/antigravity/pi also blank, which is why broadening past Codex is correctly out-of-scope); D4 — re-ran the probe contract (present stub→exit 0, silent; absent→exit 127, sole line `AGENTSVIEW MISSING`, `2>&1` suppresses the leak); D3 composes already-proven reads.
- DONE: Consolidate 69/1p/4t into ONE survey-skill pass + record disposition.
  Verified all three superseded entities carry `superseded-by: survey-skill-correctness-pass` in frontmatter (69 `survey-codex-cwd-workaround`, 1p `survey-scaffold-state-the-fact`, 4t `survey-agentsview-detect-under-sandbox`). The single body already weaves D2/D3/D4 plus the B2 D1 spine; consolidation stands.
- DONE: AC + test plan proven by EXERCISING survey behavior over real sessions.db keyed by git-root basename — never a prose-grep.
  The five ACs and 4-bullet test plan stand as-authored and were validated against the real artifacts: confirmed `fixture-sessions.sql` currently encodes the DISPROVEN divergent-key model (`_spacedock_state`/`feature_x` keys, with comments asserting divergence) and `survey_queries_test.go` currently asserts the 4-field `folded_keys=3` scoping row with no codex query — so AC-1/AC-2's fixture-correction + folded_keys-drop + codex-no-union assertions are real, failing-today, checkable changes against existing code, not prose. AC-3 probe-behavior shape exercised here (passes); AC-4/AC-5 remain the combined live drive (the survey's proof bar).

### Summary

Cycle-2 re-dispatch: verification pass, not a rewrite. Independently exercised every load-bearing premise of the existing ideation against real agentsview v0.32.1 data and the real fixture/test artifacts (per proof-policy: prove by exercising, not re-reading). The git-root-basename spine holds and is stronger than recorded (34 cwds → one key, folded_keys=1, cwd→project a function); Codex blank-cwd is 591/591; the D4 probe contract reproduces exactly; all three superseded entities are dispositioned; and the AC fixture/test targets are confirmed to be failing-today changes against existing code. The body is sound and build-ready as authored — no factual correction required, only the on-record note that the cwd count drifted 29→34 with the conclusion unchanged.

## Stage Report: implementation

- DONE: `queries.sql` + `fixture-sessions.sql` + `survey_queries_test.go` corrected to the git-root-basename model: `scoping` DROPS `folded_keys` (keeps `blank_cwd`, cwd-prefix-union UNCHANGED); add a `-- name: codex-presence` query bound to `:repo_project` returning matched `agent='codex'` count + blank-cwd sum; correct the fixture (in-repo Claude sessions share ONE git-root-basename `project`, NOT the divergent `_spacedock_state`/`feature_x`) and add blank-cwd Codex rows; tests assert the 3-field scoping shape, the prefix still counts the in-repo sessions, codex-presence count + blank_cwd>0, AND `scoping.sessions` UNCHANGED by the Codex rows. (AC-1, AC-2)
  Commit 916b2ecf. TDD: corrected the test RED first (`codex-presence` missing → FAIL), then queries+fixture GREEN. `go test ./skills/integration/ -run TestSurveyQuerySmoke` → 6/6 subtests pass (scoping 3-field `3|0|span`, codex-presence `2|2`, codex-not-folded-into-scope sessions=3). EXERCISED against LIVE agentsview v0.32.1 (this repo synced through the binary): in-repo Claude → ONE key `spacedock_v1` (folded_keys=1); codex-presence(`spacedock_v1`)=`61|61`; scoping=`48|0|...` — the 61 Codex are NOT folded in (no union).
- DONE: `SKILL.md` corrected: §2 model prose rewritten to git-root-basename (key COLLIDES with same-basename siblings; absolute cwd-prefix excludes them) and the `:repo_root`/`scoping` rationale inverted; step-4 fence drops the `folded_keys` clause, adds the conditional codex-presence hint (renders only when count>0, "match by project NAME only, cwd unrecorded"), and reshapes SCAFFOLD to the state-the-fact statement (taxonomy dropped); §1 probe swapped to `agentsview --version >/dev/null 2>&1`. (D1/D2/D3/D4)
  Commit 23316398. `:repo_project` derivation `basename | tr -c '[:alnum:]' '_'` verified against the real key (`spacedock-v1`→`spacedock_v1`, matches the 61 live Codex rows). Newline-bug caught and fixed (`printf '%s' "$(basename …)"` strips the trailing `\n` that `tr` would otherwise convert to a stray `_`). No stale `folded_keys`/`command -v`/3-bucket-taxonomy references remain (grep-verified). queries.sql preamble + scoping comment + codex-presence comment all rewritten to the true model.
- DONE: AC-3 deterministic probe-behavior test: run the exact step-1 one-liner twice — stub `agentsview` on a synthesized PATH (empty output, exit 0) and name absent (sole line `AGENTSVIEW MISSING`) — oracle from the two fixture conditions, never a SKILL.md grep. AC-4/AC-5 are the LIVE DRIVE for validation.
  Commit 23316398, `skills/integration/survey_probe_test.go`. Extracts the shipped probe line from SKILL.md (executes the artifact, not a copy — no drift) via a regex that matches ONLY the execve form; verified it REJECTS `command -v`/`which`/`test -x` and the explanatory prose line, so a revert reds via extraction failure. `go test -run TestSurveyInstallProbe` → 2/2 pass (present silent/exit-0; absent sole sentinel). Full `go test ./skills/integration/` GREEN; `go build ./...` clean.

### Summary

Implemented the consolidated survey-skill correctness pass on the corrected git-root-basename agentsview model across the four files. The gate-able halves are landed and proven: AC-1/AC-2 query-smoke (6 subtests, including codex-presence and the no-union assertion) and AC-3 probe-behavior test (2 subtests, FS-access-regression-proof via the extraction regex). The model fix was re-confirmed by EXERCISING against live agentsview v0.32.1 — synced this repo through the binary and ran the corrected queries: one Claude `project` key, 61 blank-cwd Codex matched by name, scoping excludes them. SKILL.md prose (D1-D4) rewritten with no stale-model leftovers. Two commits on `spacedock-ensign/survey-skill-correctness-pass` in the worktree.

**For validation — the live drive (AC-4/AC-5) is your proof bar, not landed here.** AC-4: drive `/spacedock:survey` end-to-end on a repo whose synced DB carries project-matched blank-cwd Codex sessions — assert the rendered report shows (a) the Codex-presence hint (count + "by project name only"), (b) SCAFFOLD as the state-the-fact statement with NO active/installed/recovered taxonomy labels, (c) the PROJECT line with NO "coalesced from N keys" clause; a second drive against a Codex-free DB omits the hint. AC-5: drive under the sandbox, observe the survey proceeds past step 1 without the `brew install` consent prompt, with the pre-fix `command -v` confirmed to false-negative in the SAME sandbox as the negative control. agentsview v0.32.1 is present on this box; the sync recipe in SKILL.md §1 seeds the DB. Per ideation, the live drive may escalate to a captain-run drive. A grep over SKILL.md does NOT satisfy AC-4/AC-5 — only the rendered output / observed non-prompt does.

## Stage Report: validation

- DONE: Reproduce AC-1/AC-2 at HEAD 159f19f2 — `TestSurveyQuerySmoke` green; corrected `scoping` returns the 3-field shape (NO `folded_keys`), cwd-prefix-union counts the in-repo sessions, `codex-presence` returns matched count + blank_cwd>0, `scoping.sessions` UNCHANGED by Codex rows.
  `go test ./skills/integration/ -run TestSurveyQuerySmoke` → 6/6 subtests pass. Exact oracle from the materialized fixture DB: `scoping` → `3|0|2026-06-05 .. 2026-06-06` (3 fields, folded_keys gone); `codex-presence` → `2|2`; `codex-not-folded-into-scope` → sessions=3.
- DONE: Confirm `fixture-sessions.sql` encodes the git-root-basename model — ONE in-repo Claude `project` key, NOT the disproven `_spacedock_state`/`feature_x`.
  Queried the materialized DB on-disk (not a re-read): the 3 in-repo Claude checkouts (root, `.spacedock-state`, `.worktrees/feature-x`) all carry ONE `project='proj'`; the disproven divergent keys count = 0; the 2 Codex rows carry `project='proj'`, BLANK cwd.
- DONE: Independently confirm AC-1/AC-2 non-vacuousness — a regression must red the relevant assertion.
  AC-1: re-added `COUNT(DISTINCT project) AS folded_keys` to the `scoping` SELECT → `scoping` reds (`got "3|1|0|…"`, 4 fields; folded_keys=1 confirms the structural-always-1 finding). AC-2: unioned the `project=:repo_project` Codex rows into the Claude scope → `scoping.sessions` 3→5, both the count-3 and the dedicated `codex-not-folded-into-scope` assertions red. Both reverted via `git checkout`; re-verified green; worktree clean at 159f19f2.
- DONE: Reproduce AC-3 — `TestSurveyInstallProbe` green; the probe-behavior test EXTRACTS the shipped §1 line from SKILL.md (executes the artifact) and its regex REJECTS FS-access forms; shipped §1 probe is `agentsview --version >/dev/null 2>&1`.
  `go test -run TestSurveyInstallProbe` → 2/2 pass (present: silent/exit-0; absent: sole line `AGENTSVIEW MISSING`). SKILL.md:29 confirmed = `if ! agentsview --version >/dev/null 2>&1; then echo "AGENTSVIEW MISSING"; fi`.
- DONE: Independently confirm AC-3 non-vacuousness — reverting the probe to an FS-access form reds via extraction failure.
  Reverted SKILL.md:29 to `command -v` → test reds `expected exactly one runnable … probe line … found 0`. Standalone regex check: the extraction regex matches ONLY the execve form and rejects `command -v`/`which`/`test -x`/`stat`, the line-23 prose mention of `command -v`, and a near-miss missing `2>&1`. Reverted; re-verified green.
- DONE: Full `go test ./skills/integration/` + `go build ./...` green at 159f19f2.
  Full integration package `ok` (1.185s); `go build ./...` exit 0. Worktree clean at HEAD 159f19f2 after all non-vacuousness probes reverted.
- DONE: Prep the CAPTAIN AC-4/AC-5 live-drive handoff (do NOT block PASSED on running it).
  Ready-to-run handoff assembled below; prerequisites verified present (agentsview v0.32.1, `spacedock` at /opt/homebrew/bin, sqlite3); SKILL.md step-2 key derivation spot-checked against the real repo (`spacedock_v1`, resolves to the main checkout even from a worktree).
- SKIPPED: Detached adversarial audit.
  Per the assignment and the stage definition's high-stakes list, survey (skill prose + SQL) is NOT a high-stakes surface → normal validation, no detached audit.

### Captain handoff — AC-4 / AC-5 live drive (sprint-acceptance, not gate-blocking)

Verified SHA: **159f19f2**. Branch `spacedock-ensign/survey-skill-correctness-pass`, worktree `/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-survey-skill-correctness-pass`. Prereqs confirmed on this box: agentsview v0.32.1, `spacedock` launcher, sqlite3.

**Seed recipe (per SKILL.md §1).** Drive the read through the binary into a readable data dir; seed a DB that carries project-matched blank-cwd Codex sessions:
```
export SPACEDOCK_SURVEY_DB_DIR="$(mktemp -d)/spacedock-survey"; mkdir -p "$SPACEDOCK_SURVEY_DB_DIR"
# narrow CLAUDE_PROJECTS_DIR to THIS repo's session dirs (SKILL.md §1 symlink-farm block), then:
AGENTSVIEW_DATA_DIR="$SPACEDOCK_SURVEY_DB_DIR" CLAUDE_PROJECTS_DIR="$NARROW" timeout 300 agentsview session sync
# this repo's Codex history lands cwd='' / project=spacedock_v1 → the Codex-presence hint fires (61 live rows observed during implementation).
```
For the AC-4 negative (hint-absent) drive, seed a DB whose Codex rows do NOT match `spacedock_v1` (or omit Codex sync).

**Invocation.** `spacedock claude --plugin-dir /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-survey-skill-correctness-pass`, then run `/spacedock:survey`.

**AC-4 — observe in the rendered report (SKILL.md §4 fence, lines 152-157):**
- (a) Codex hint line under PROJECT: "N Codex sessions match this repo by project NAME only (agentsview does not record Codex cwd) — may include a same-named sibling repo …". Codex-free DB → line ABSENT.
- (b) SCAFFOLD as the state-the-fact statement (family + invocation count + checked-in presence/absence; a family invoked-but-not-on-disk reads "recovered from behavior, not files") with NO recovered/installed/active taxonomy LABELS.
- (c) PROJECT line with NO "coalesced from N keys" clause (folded_keys is gone; only the `{if blank_cwd>0}` clause may appear).

**AC-5 — under sandbox:** the survey proceeds past step 1 to sync/scan WITHOUT emitting the `brew install --cask agentsview` consent prompt (agentsview present). Negative control in the SAME sandbox: confirm `command -v agentsview >/dev/null` false-negatives there, establishing the drive exercises the failing path rather than a vacuous pass.

A grep over SKILL.md does NOT satisfy AC-4/AC-5 — only the rendered output / observed non-prompt does.

### Recommendation

**PASSED on the gate-able halves (AC-1, AC-2, AC-3).** All three reproduce green at 159f19f2, the fixture encodes the corrected git-root-basename model on-disk, and each is independently non-vacuous (folded_keys re-add reds AC-1; Codex union reds AC-2; `command -v` revert reds AC-3 via extraction failure). Full integration suite + `go build ./...` green. **AC-4/AC-5 are handed to the captain as sprint-acceptance** (the rendered-report and sandbox-non-prompt live drive) — not gate-blocking per the assignment; the ready-to-run handoff is above.

### Summary

Validated the consolidated survey-skill correctness pass at HEAD 159f19f2 by exercising behavior, not re-reading. The three gate-able ACs pass and are proven non-vacuous via the adversarial edits the assignment named (re-add folded_keys → AC-1 reds; union Codex into Claude scope → AC-2 reds; revert probe to `command -v` → AC-3 reds via extraction failure). The fixture's git-root-basename model is confirmed on-disk (one in-repo `project` key, zero disproven divergent keys), and SKILL.md's step-2 key derivation spot-checks correctly to `spacedock_v1` even from a worktree. Full integration suite and build are green; the worktree is clean at HEAD after every probe was reverted. AC-4/AC-5 (live-drive rendered report + sandbox non-prompt) are handed to the captain as sprint-acceptance with a verified ready-to-run recipe. Survey is not a high-stakes surface, so no detached audit. Recommend PASSED.

## Stage Report: implementation (feedback cycle 1)

- DONE: Root-cause the captain's AC-4 "0 Codex sessions" (build saw 61) — reproduce end-to-end, confirm/refute the Codex-blind-sync hypothesis.
  REFUTED the hypothesis. Reproduced on this box (agentsview v0.32.1): the exact §1 recipe sync into a fresh SURVEY_DB_DIR ingests claude 1346 + codex 592; the EXACT SKILL.md run_query codex-presence returns 61|61. CLAUDE_PROJECTS_DIR narrows ONLY Claude; Codex syncs UNSCOPED from the default ~/.codex/sessions (separate CODEX_SESSIONS_DIR), so the sync is NOT Codex-blind. ACTUAL root cause: the persisted default SURVEY_DB_DIR ($TMPDIR/spacedock-survey) was a stale Claude-only DB (4856 claude, 0 codex, last synced 2026-06-07 pre-fix) → codex-presence = 0|0; the recipe reuses it and its "incremental/persists" framing invited skipping the sync. One incremental sync backfills it 0→592 (no --full, no watermark-skip — verified with an old-dated synthesized session 0→1).
- DONE: Fix the real end-to-end gap (recipe robustness + the test gap AC-2 missed), minimal, within survey skill + tests.
  SKILL.md §1: state the narrowing is Claude-only and Codex syncs unscoped from its default dir (what populates codex-presence); make the sync MANDATORY every survey — never query a pre-existing SURVEY_DB_DIR without re-syncing (skipping the sync is what yields the stale 0). Query/mechanism UNCHANGED. Commit 396534cc.
- DONE: Add the sync→codex-presence end-to-end test (the coverage AC-2's direct-DB-injection fixture lacked).
  skills/integration/survey_sync_codex_test.go: drives the real `agentsview sync` over a HOME-isolated synthesized Codex source (two same-basename roots → the documented collision), then runs the codex-presence query from queries.sql against the SYNCED DB and asserts 2|2 (blank_cwd>0). Verified load-bearing: reds (0|0) on a Codex-blind sync. `go test ./skills/integration/` green (probe 2, query-smoke 6, sync-e2e 1); go vet + go build clean.

### Summary

Feedback cycle 1: the captain's "0 Codex" was NOT a query/sync-mechanism defect — root-caused (by exercising, not re-reading) to a stale persisted Claude-only SURVEY_DB_DIR queried without a Codex-inclusive re-sync. The recipe DOES ingest Codex (unscoped from ~/.codex/sessions). Fixed the genuine gap: §1 now mandates the sync every run and documents the Codex source, and a new HOME-isolated sync→codex-presence e2e test exercises the path the fixture-only AC-2 never did (reds on a Codex-blind sync). HEAD 396534cc on spacedock-ensign/survey-skill-correctness-pass. For the re-drive: start from a FRESH SURVEY_DB_DIR, or rely on the now-mandatory sync to backfill the persisted one (I synced this box's persisted dir to current — codex-presence there is now 61|61).

### Feedback Cycles

**Cycle 1 (validation → implementation) — RESOLVED.** Trigger: the captain's AC-4 live drive rendered "0 Codex sessions match this repo" against a phantom zero (the build had observed 61). Root cause (refuted the sandbox/all-agents hypothesis): the survey queried a STALE persisted Claude-only `SURVEY_DB_DIR`; §1 reused the dir and framed it as "incremental/persists," inviting a skipped sync. The sync itself is not Codex-blind (Codex ingests unscoped from `~/.codex/sessions`). Fix (HEAD `396534cc`, query/mechanism unchanged): §1 mandates a fresh sync every survey + documents the Codex source; new `survey_sync_codex_test.go` exercises sync→codex-presence end-to-end (the gap AC-2's fixture-injection missed), non-vacuous (reds 0|0 on a Codex-blind sync). Captain re-drove AC-4 with `~/.codex` allowed → correct 61-Codex hint. Re-validation dispatched to the kept-alive reviewer. Not escalated (cycle 1 of 3).
