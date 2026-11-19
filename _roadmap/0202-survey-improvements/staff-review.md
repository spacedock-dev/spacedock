# 0202 survey-improvements — preflight staff review

Independent staff review run BEFORE the sprint drive. Mandate: refute the ideation
designs against the real repo, not ratify the entities' self-reports. Read-only over the
five entity bodies plus the code/config/fixtures/schema each claim rests on. The verdicts
below are mine, not the ensigns'.

**Method.** Read each entity body in `docs/dev/.spacedock-state/` (the four flat seeds +
`state-sync-no-origin-local-mode/index.md`). For every load-bearing claim I checked the
real artifact: `skills/survey/SKILL.md` + `references/queries.sql`, the survey fixture
`skills/integration/testdata/survey/fixture-sessions.sql`, `internal/status/` (native_runner,
new, validate, boot, json_commands, handlers), `internal/dispatch/build.go`, `internal/cli/`,
`docs/schema/entity.mdschema.yml`, `install.sh`, `.github/workflows/*.yml`, the source seeds
`9h`/`5x`/`za`/`zb`, the captain-locked mock `docs/roadmap/0202-survey-improvements/index.md`
(it lives on `next`, read via `git show next:…`), and ran `gofmt -d` / `gofmt -l` under the
host go1.26.1. The survey skill + queries.sql are byte-identical between `drive/0201` and
`next` (empty diff), so the `5wv` verification holds regardless of base.

**Verdict: NOT-READY** — one material hole on `5wv` (the no-follow-up figure), the rest READY.

## Verdict per member

- **`5wv` survey-output-redesign** — **MATERIAL (1).** The riskiest-mechanism table
  silently redefines the no-follow-up lede figure as a count of the existing BACKLOG
  cross-check to claim "no new query" — contradicting its own source seed `9h`, which
  spiked that figure as a `message_id → ordinal` chronological join needing a NEW query
  AND a fixture extension the current fixture lacks. The fixture-ordinal requirement is
  absent from `5wv`'s test plan. Everything else (zb#1 shipped, the other two frontier
  figures as derived counts, the vocab rename, the dispatch-fact-is-the-only-new-query
  finding) verifies clean.
- **`nd` prefer-new-over-next-id** — **READY (0 material).** AC-2's proof is a genuine
  live FO drive grading the tool-call stream, not a prose-grep. Every code seam it names
  is real. Refuted nothing material.
- **`td` mdschema-conformance-validator** — **READY (0 material).** The spike's claimed
  schema patterns/enums are present in the real schema; the warn-tier design correctly
  honors the documented read-path lockout hazard; ACs feed bad fixtures with independent
  oracles. One blast-radius polish item. Refuted nothing material.
- **`5ar` pre-cut-audit-cleanups-0199** — **READY (0 material).** The checksum tamper
  test goes genuinely RED when the gate is stripped (load-bearing assertion is real); the
  node-action target majors and flip-cut-path attribution match the actual workflow files
  exactly; the go1.26.1 gofmt curly-quote trap is real and correctly dodged. Refuted
  nothing material.
- **`gf` state-sync-no-origin-local-mode** — **READY (0 material).** Every seam (boot
  STATE_BACKEND, bootJSON key-order, stateCommitGuidance + its two callers, the git
  runners, the remote-tolerant create path) is real at the cited locations; the
  detection spike is a sound network-free discriminator. Refuted nothing material.

Per-member Material count — **5wv: 1 · nd: 0 · td: 0 · 5ar: 0 · gf: 0.**

---

## Material findings (fold before the gate)

### M1 (`5wv`, AC-1 / R1) — the no-follow-up figure is redefined to dodge the `9h` spike's new-query + fixture-ordinal requirement

**This is the only material hole in the sprint, and it sits on the highest-risk member.**
`5wv`'s load-bearing claim — the one the assignment flagged for verification — is that
"the 3 frontier lede figures are derived COUNTS of what the existing decision-open query +
step-4 cross-check already produce — no new query." For two of the three figures
(hand-steering interruptions, hanging threads) this is TRUE and I confirmed it:
`decision-open` exists at `skills/survey/references/queries.sql:275` (AskUserQuestion/
ExitPlanMode at :300), the veto markers `[Request interrupted` / `doesn't want to proceed`
exist at queries.sql:314/339, and there is no `dispatch-fact` query today — so dispatch-fact
is genuinely the only new query for those.

But the **third** figure — `decisions you made with no follow-up action` — does not hold.

- **The source seed `9h` spiked it as a NEW query over a chronological join.**
  `9h` (survey-value-prop-legibility, the R1 seed) records (entity lines 121, 135-141):
  the figure is the count of `done` decisions with no *later* Edit/Write, where "later"
  is `tool_calls.message_id → messages.ordinal` — NOT `tool_calls.id` insertion order. It
  names a new `decision-no-followup` query and states the load-bearing build risk
  explicitly: *"the committed fixture lacks both columns — so the AC-2 fixture MUST be
  extended with message ordinals (and `tool_calls.message_id`)."* Live figure: 2/102.

- **`5wv` redefined the figure to a BACKLOG count to claim "no new query."**
  `5wv`'s riskiest-mechanism table (entity line 49) sources no-follow-up from "the **count**
  of cross-check `decided-not-shipped` forks (the `BACKLOG` set)" → verdict "PROVEN — derived
  count of the existing cross-check, no new query." Its render spine (line 69) and doc diff
  (line 87) repeat the BACKLOG framing. That is a *decided-but-not-shipped* metric — a
  different thing from `9h`'s *decision-with-no-subsequent-edit* metric. The substitution is
  what lets `5wv` declare "no new query."

- **The captain-locked mock means `9h`'s figure, not the BACKLOG count.** The mock
  (`docs/roadmap/0202-survey-improvements/index.md:33-35`, on `next`) lists `hanging
  threads` and `decisions you made with no follow-up action` as two *separate* lede
  figures (values 6 and 3), and the prose immediately under it (line 48) cites it as
  "`9h`#4". So the locked structure expects `9h`'s no-subsequent-edit figure as a distinct
  number — `5wv` cannot satisfy it by reusing the BACKLOG count, which is already its own
  separate `BACKLOG` section.

- **The fixture genuinely lacks the columns `9h` flagged.** I read
  `skills/integration/testdata/survey/fixture-sessions.sql`: `tool_calls` (lines 64-72) has
  no `message_id` column; `messages` (lines 74-79) is `(id, session_id, role, content)` —
  no `ordinal`. So the join `9h` needs cannot be exercised against the current fixture.

- **`5wv`'s test plan never adds them.** The R1 entries in `5wv`'s test plan (entity lines
  121-129) cover the dispatch-fact columns (`parent_session_id`/`relationship_type`), the
  knowledge-work track, the vocab rename, the all-unlabeled Codex set, and `5x`'s date/cwd
  fixtures — but never the `message_id`/`ordinal` extension or a `decision-no-followup`
  query. AC-1's "non-vacuous: a fixture mutation that adds an interruption flips the count"
  exercises only the *interruptions* figure, not no-follow-up.

**Net:** as written, R1's no-follow-up figure either (a) silently ships a different metric
than the seed and the locked mock ask for, or (b) ships the right metric with no query, no
fixture columns, and no AC to prove it — an unbuilt, untested arm hiding behind a
"no new query" claim.

**Fix (owner: `5wv`).** Restore `9h`'s definition for the no-follow-up figure: adopt the
`decision-no-followup` query (the `message_id → ordinal` chronological join), add the
`tool_calls.message_id` column to the fixture's tool_calls rows and the `ordinal` column to
its messages rows, seed a no-follow-up decision plus a follow-up decision (a later-ordinal
Edit), and give AC-1 an independent oracle for *this* figure (non-vacuous: inserting an Edit
at a higher ordinal than a no-follow-up decision's message decrements the count — `9h`'s own
AC-2). Correct the riskiest-mechanism table row from "no new query" to "one new query
(`decision-no-followup`) + a fixture-ordinal extension, per `9h`'s spike." This does not
disturb the genuinely-correct finding that the other two frontier figures and dispatch-fact
are the only new query beyond it — only the third figure's accounting is wrong.

---

## Polish (fold opportunistically; none gate)

### P1 (`5wv`) — "ALREADY FIRMED" overstates the state of `5x` and `za`

`5wv` says R2 "adopts `5x`'s firmed AC-1/AC-2 and doc diff verbatim" and R3 "adopts `za`'s
firmed query + ACs verbatim" (entity lines 35-36). Both `5x` (survey-lens-honesty) and `za`
(survey-report-subagent-dispatch-fact) are `status: backlog` — not complete/firmed. The
verbatim source content does exist in those bodies (27/28 matching lines for the
dispatch-fact query and the WHAT-THIS-CAN'T-SEE doc diff respectively), so this is not a
missing-dependency hole — but `5wv`'s AC-2/AC-3 inherit a dependency on two unbuilt seeds,
and "firmed" is the wrong word for a backlog entity. Soften to "the verbatim query/ACs are
authored in `5x`/`za` and folded here," and own that R2/R3's fixtures land as part of this
build (not pre-existing). No structural change needed.

### P2 (`5wv`) — R6#2 is the one band whose AC is fully conditional

Of R1-R6, R6#2 (branch-aware work-by-area) is the only ask whose AC (AC-6 clause b) is
explicitly "[conditional]" with a "degrades to caveat-only" fallback. That is an honest
scoping, not a hole — every other band (R1-R5, R6#1) has a concrete fixture-bound AC. Flagging
only so the gate knows R6#2 may ship as detect-and-caveat rather than full re-attribution;
the entity already records this as the lowest-priority arm.

### P3 (`td`) — warn-tier stderr could perturb the exact-match `stdout == "VALID"` assertions

`td` extends `status --validate` with a warn tier. The existing
`internal/status/native_validate_test.go` asserts `stdout` is exactly `"VALID"`
(`TestValidateFlagsSelfRefACs`, :156) and that defects exit 1 (`TestNativeValidationGatesReads`).
`td`'s AC-2/AC-3 correctly preserve the exit-1 structural contract and keep warns off the
read-path gate, so the *contract* is right — but at implementation, field warnings printing
to stderr on the existing fixtures (if any of their frontmatter trips a pattern/enum) could
break those golden assertions. Low risk (the fixtures are narrow), but worth a green-suite
check on the existing validate tests as part of the build, not just the new ones. Owner: `td`.

### P4 (`5ar`) — AC-2's header-vs-config text check is self-flagged as possibly brittle

`5ar`'s AC-2 (release.yml header darwin→darwin+linux) honestly records a fallback: if the
header-vs-`.goreleaser.yaml`-goos text check "proves too brittle to bind cleanly," it
downgrades to a reviewer-confirmed prose fix backed by the existing
`TestGoreleaserBuildsLinuxAndDarwin` (which I confirmed exists at
`internal/release/goreleaser_guard_test.go:53`). That existing test already proves linux is
built, so the fallback has a real backstop. No action — recording that AC-2's binding oracle
is the one "implementation decides" call in the package, and it has a safe floor.

---

## What I verified clean (so the gate knows these are NOT open)

- **`5wv` claim (a) — zb#1 shipped in #335.** TRUE. `skills/survey/SKILL.md:162-163`
  leads the CODEX block with `{codex_scoped_sessions} … attributed to this repo by
  exec_command working dir` and demotes name-match to a caveat; #335 is the merge
  (d6b27eb9). R5 correctly narrows to just the knowledge-work archetype — and
  `mode-classification` (queries.sql:399-401) really has only mechanical/exploration/
  unlabeled, so the archetype is genuinely new. The `'mechanical' → 'manual'` rename
  target is the literal at queries.sql:399.

- **`nd` AC-2 is a live drive, not a grep.** AC-2 (entity line 41) grades a new `filing`
  shared scenario in `internal/ensigncycle` on the FO's recorded tool-call stream (`new`
  invoked, `--next-id`+`Write` absent) — explicitly "not a grep of the contract prose, and
  not just the end-state file." The code seams are real: stdout-only id write at
  `internal/status/native_runner.go:224`, the JSON branch separate at :220-221, `runNew`
  atomic-create at `internal/status/new.go`, and the `new` verb aliasing `status --new` at
  `internal/cli/cli.go:308`. The contract files it edits are under `skills/first-officer/
  references/` — which is exactly why a live drive (not self-certification) is the right proof.

- **`td` spike patterns are real.** `docs/schema/entity.mdschema.yml` carries `mod-block`
  pattern `^[^:]+:[^:]+$` (severity `warn`, line 96-97) and `verdict` enum `[PASSED,
  REJECTED]` (severity `warn`, line 92) — exactly what the spike claims to have extracted.
  `yaml.v3` is in go.sum. The read-path lockout hazard the warn-tier design protects is the
  real documented comment at `internal/status/validate.go:138-145`.

- **`5ar` #1 checksum tamper goes RED when the gate is removed.** The gate is
  `install.sh:160-169` (awk the expected hash, sha256 the tarball, `die` on mismatch/missing
  line). The `SPACEDOCK_INSTALL_FROM=<dir>` local-dist path (install.sh:115-134) consumes a
  plain `dist/` with `spacedock_*_<os>_<arch>.tar.gz` + `checksums.txt`, drivable from Go via
  `sh install.sh` exactly like `internal/release/install_url_test.go` (which exists, with
  `scrubInstallEnv` at :133). AC-1's load-bearing half — strip lines 164-169 to a temp copy,
  run the tamper case, assert it now wrongly exits 0 — is a genuine "the test would have
  reded" proof, modeled on the real `TestGoreleaserBuildGuardRejectsDroppedLinux`
  (goreleaser_guard_test.go:79).

- **`5ar` #5 node-action targets match the real workflows.** Confirmed every pin against
  `.github/workflows/`: `checkout@v4` (→v5), `setup-go@v5` (→v6), `goreleaser-action@v6`
  (→v7; the entity is right that "v6 stays node20") all present, with `goreleaser-action@v6`
  on the flip-cut path at release.yml:175 AND install-e2e.yml:39. `setup-node@v4`,
  `setup-python@v5`, `upload-artifact@v4` present as described. `deploy-pages@v4` /
  `upload-pages-artifact@v3` appear ONLY in docs.yml — so the entity's "no node24 release
  yet, not on the flip-cut path" carve-out is accurate.

- **`5ar` #3 gofmt trap is real.** `gofmt -l skills/integration/survey_sync_codex_test.go`
  lists the file (drift real); `gofmt -d` under host go1.26.1 rewrites the line-23 comment
  `cwd = ''.` into `cwd = ”.` (U+201D) — the documentation-content change the entity warns
  about. The reword-then-format fix is correct.

- **`gf` seams are all real at the cited lines.** `STATE_BACKEND:` text at
  `internal/status/boot.go:283-284` with `stateBackend` set at :196-198; the bootJSON
  key-order seam (`state_backend` then `entity_dir_present`) at
  `internal/status/json_commands.go:195-198` with the append-after comment at :193;
  `stateCommitGuidance(stateCheckout, entityPath, stateBranch)` at
  `internal/dispatch/build.go:719` with callers at :498/:508 and the push/pull-rebase
  reminder at :721/:728; `runGitCmd` at handlers.go:575; the remote-tolerant create-path
  warn at `internal/cli/state.go:178`; `build_statecommit_test.go` present. The detection
  spike (`git remote get-url origin` exit 2 vs 0) is a sound network-free named-remote
  discriminator, and the contract push/pull clause is correctly flagged as a prose delta
  proven indirectly by the dispatch fixture's generated-prompt-bytes assertion.
