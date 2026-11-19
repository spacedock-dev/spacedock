# 0223 pi-dispatch-contract — sprint-wide preflight readiness review #2 (re-carve)

Independent staff review #2 of sprint `0223-pi-dispatch-contract` **as a whole**,
re-refuting the **re-carved** 3-member sprint. Scope: did the re-carve actually
close what review #1 (`staff-review.md`, run `3adf00ee`) found, and did it
introduce new problems? Sprint-level coherence, cross-member composition,
blast-radius, and Commander cold-boot readiness — NOT per-task AC quality (the
ideation gate owns that).

Evidence base: the re-carved `index.md`; review #1 (`staff-review.md`); the
three live member ideations (`pi-install-managed-skill-placement` `eq`,
`pi-dispatch-model-stamping` `bdt`, `pi-back-channel-dispatch` `b2`); the
dispatch core (`skills/first-officer/references/fo-dispatch-core.md`);
`internal/cli/pi.go`; and the two archived superseded tasks
(`pi-ensign-skill-injection` `k8t`, `pi-launcher-repo-resolution` `2m1`).

## Verdict

**Gaps to close — not yet cold-boot drivable.**

The headline finding is positive: **review #1's blocker (the child-cwd seam) is
genuinely closed** by the re-carve. The merged `pi-install-managed-skill-placement`
(`eq`) removes the seam at its root — `collectSettingsPackageSkillPaths` is a
settings/package-store scan, not cwd-keyed, and `eq`'s spike proved discovery
from a non-repo cwd. Gap 2 (Pi canonical-model-space declaration ownership) is
also substantively closed. The re-carve's *mechanism* is sound and the merge is
cleanly motivated; no material redesign is needed.

What blocks cold-boot is **incomplete follow-through**: the re-carve changed the
membership and the skill-discovery mechanism but did **not** re-validate the
cold-boot package or the cross-references that were written against the old
4-member / cwd-keyed-symlink model. Concretely:

- **Q11 is now self-contradictory and partly false** — it asserts that without
  `cwd:<repo>` on the dispatch "member 1's symlink is not discovered and ensign
  doesn't load (DoD bullet b re-breaks)," which `eq`'s spike directly disproves
  and the capstone's own gap-1 re-check explicitly retracts. The index still
  claims gap 3 (Q11–Q13) is "CLOSED," but Q11 was never re-checked against the
  re-carve.
- **A new composition gap** — the re-carve archived `2m1` (the install-record
  repo-path file), which was the documented *source* for both the launcher's
  `repoRoot` resolution and the capstone's `cwd:<repo>` forwarding argument.
  `eq` retires the cwd fallback but its design never defines what `repoRoot`
  resolves to post-install; `eq`'s own AC-3 references an "install-record
  resolution" that `eq`'s D3 does not produce. The capstone's re-check claims the
  `cwd:<repo>` source is "unchanged" — it is not.
- **Pervasive stale member references** — the `Sequencing` section, DoD "Proven
  by" bullets 1–2, Q3, Q9, Q11, Q12, Q13, the gap-2 closure note, the capstone's
  fold-in AC-2 sub-bullet, and the lifecycle checklist all still cite the old
  4-member numbering or the archived task names.

All of these are documentation/clarification-level closes; none requires
re-opening the mechanism. Once the cold-boot package is re-validated and `eq`'s
`repoRoot` source is specified, the sprint is drivable.

## Gap-1 closure verification

**The blocker is closed.** Verified end-to-end:

- Review #1's seam was: `k8t`'s `.pi/skills/ensign` symlink is cwd-keyed, and
  under `2m1`'s non-repo launch the child's cwd defaults to the parent's
  non-repo cwd (`step.cwd ?? ctx.cwd`), so the symlink is not in the child's
  discovery path and ensign silently falls back to bare `worker`.
- `eq` replaces the cwd-keyed symlink with an install-managed package
  (`package.json` `pi.skills` + `.pi/extensions/spacedock.ts`). Child discovery
  now flows through `collectSettingsPackageSkillPaths` (skills.ts:288), which
  reads `settings.json` `packages` → each package's `package.json` `pi.skills`
  via `resolveSettingsPackageRoot` — a **settings/package-store scan, not a
  `{child-cwd}/.pi/skills` scan.** It is not cwd-keyed.
- `eq`'s spike PROVED this: `pi install /tmp/pi-skill-spike` (minimal
  `pi.skills` package) → `discoverAvailableSkills("/tmp")` (non-repo cwd) lists
  `spike-skill (user-package)`. "PRESENT from a non-repo cwd — no cwd
  dependency." This is the inverse of the `0637e2ed` failure mode review #1
  identified.

**Capstone reframe consistency.** The capstone's "Staff review gap-1 re-check"
section is internally consistent with — and correctly SUPERSEDES — the earlier
"Staff review fold-in" section's gap-1 stance:

- The fold-in section (written pre-re-carve) framed `cwd:<repo>` as the wiring
  that makes `k8t`'s symlink discoverable under `2m1`'s non-repo launch, and its
  AC-2 sub-bullet said "verify the ensign loads … via member 1's `.pi/skills/ensign`
  symlink."
- The re-check section marks the "required for ensign discovery" claim
  **SUPERSEDED**, re-attributes the ensign-loaded pre-check to the package-root
  scan (not `cwd:<repo>`), and reframes `cwd:<repo>` as a working-directory
  concern. This is the correct stance and it holds.

**Two residual issues with the reframe (not blockers, but must be closed):**

1. **The fold-in's AC-2 sub-bullet text is left standing un-reconciled.** The
   re-check is append-only ("no prior section was rewritten"), so the fold-in
   AC-2 sub-bullet still reads "verify the ensign loads … via member 1's
   `.pi/skills/ensign` symlink" — referencing (a) a member that no longer exists
   (`k8t` archived) and (b) a mechanism (`eq` does not ship a symlink). The
   re-check says it "re-attributes" the pre-check, but a Commander reading the
   AC text top-to-bottom hits the stale symlink claim before the re-check. This
   is a readability hazard in load-bearing AC text. (Owner: `b2` capstone — add a
   SUPERSEDED pointer on the fold-in AC-2 sub-bullet itself, or reconcile the
   text.)

2. **The reframe's `cwd:<repo>` "source is unchanged" claim is false** — see
   Cross-member composition. The re-check states "The source path
   (install-recorded / `--plugin-dir` / `SPACEDOCK_REPO_ROOT` resolution) is
   unchanged." But `2m1` (which wrote the install-record repo-path file) is
   archived, and `eq`'s design does not produce an install-record repo path. So
   the "install-recorded" source — the one that matters for the headline
   non-repo, non-dev-override launch — lost its producer.

**No scenario where the reframe reintroduces the seam for skill discovery.**
The package-root scan is genuinely cwd-independent; the reframe does not
re-key discovery on cwd. The only residual risk is the `cwd:<repo>` *argument's
source* (composition, below), not skill discovery.

## Merge cleanliness

**The merge is clean and well-motivated; one scope edge is underspecified.**

`eq` cleanly absorbs BOTH archived tasks' scope:

- **From `k8t` (ensign discovery):** `k8t` shipped a cwd-keyed `.pi/skills/ensign`
  symlink for child discovery. `eq` replaces it with `package.json` `pi.skills`
  discovered via `collectSettingsPackageSkillPaths` — a strictly better mechanism
  (no cwd dependency, no symlink, no clone). `k8t`'s AC-1 (discovery lists
  ensign), AC-2 (probe loads ensign from non-repo cwd, no `skillsWarning`), and
  AC-3 (committed artifact travels with install) all map to `eq`'s AC-2/AC-4.
  Nothing lost.
- **From `2m1` (repo resolution):** `2m1` shipped an install-record repo-path
  file + cwd-fallback demotion. `eq`'s D3 (install runs `pi install`) + D4
  (retire `--skill` flags + cwd fallback) absorb the *intent* — install is no
  longer check-only, and the cwd fallback is retired. `2m1`'s AC-1 (non-repo
  launch resolves skills), AC-3 (pi-only, claude/codex unaffected), and AC-4
  (install accepts `--plugin-dir`) all map to `eq`'s AC-1/AC-2/AC-3/AC-4.

**Nothing from `k8t`/`2m1` is lost in the merge** at the capability level.

**The merge is NOT over-scoped.** `eq`'s four deliverables (root `package.json`,
`.pi/extensions/spacedock.ts`, install rewrite, retirement) are facets of one
coherent mechanism (ship as a pi package). This is more coherent than the two
separate clone-bound workarounds it replaces.

**One scope edge is underspecified (the merge's one rough seam):** `2m1` owned
explicit `repoRoot` *resolution* (the launcher reading an install-recorded repo
path ahead of cwd). `eq` retires the cwd fallback (D4) but its design never
says what `repoRoot` resolves TO post-install. The current `piRuntimeConfigFromEnv`
(pi.go:214) resolves `repo` as `--plugin-dir` → `SPACEDOCK_REPO_ROOT` → cwd, and
`cfg.firstOfficer`/`cfg.ensign` (used by the doctor's `firstOfficerOK`/`ensignOK`
checks) derive from it. `eq`'s D4 retires the cwd fallback but D3 (run `pi install`)
does not write a repo-path install record. `eq`'s AC-3 references "the
install-record resolution" — a producer its own design does not define. This is
the composition gap detailed under Cross-member composition. It is a
specification gap in `eq`, not a lost-scope gap.

## DoD coverage

**Every DoD bullet still has an owning member after the re-carve. No orphaned
DoD bullet.** The (a)/(b)/(c)/(d)/core bullets map: (a)→`bdt`, (b)→`eq`,
(c)→`eq`, (d)→`b2`, core→`b2`.

**But two DoD "Proven by" items name archived tasks, and bullet (c)'s mechanism
phrasing is stale:**

- DoD "Proven by" item 1 reads `pi-ensign-skill-injection:` (archived `k8t`) —
  its *content* (subagents-doctor lists ensign; probe loads ensign from non-repo
  cwd) is now owned by `eq`'s AC-2. The name is stale.
- DoD "Proven by" item 2 reads `pi-launcher-repo-resolution:` (archived `2m1`) —
  its *content* (non-repo launch resolves skills; repo path install-recorded not
  cwd-derived) is now partly owned by `eq`. The name is stale.
- DoD bullet (c) says "resolves the Spacedock repo explicitly via
  **install-recorded path** or `--plugin-dir`/`SPACEDOCK_REPO_ROOT`." The
  "install-recorded path" was `2m1`'s mechanism (a repo-path file). `eq`'s
  mechanism is a pi package registration (`settings.json` `packages`), not a
  repo-path file. The bullet is still *satisfiable in spirit* (cwd-luck is
  retired), but the phrasing references a replaced mechanism and, per the
  composition gap below, the "install-recorded path" source is currently
  undefined in `eq`'s design.

These are freshness gaps in the DoD text, not coverage gaps. (Owner: Shaping FO
— refresh DoD "Proven by" items 1–2 to name `eq`, and reconcile bullet (c)'s
"install-recorded path" with `eq`'s package-install mechanism once `eq`'s
`repoRoot` source is specified.)

## Sequencing

**The re-carved sequencing is structurally correct, but the `index.md`
"Sequencing" section itself is stale (old 4-member layout).**

- **Members 1 (`eq`) + 2 (`bdt`) parallel** — correct. Independent code surfaces
  (`eq`: `pi.go` + `package.json` + `.pi/extensions/`; `bdt`:
  `pi-first-officer-runtime.md` prose). No inter-member dependency. ✓
- **Member 3 (`b2` capstone) `pi-live` drive requires 1+2** — correct. AC-2
  needs ensign discoverable (`eq`) + parent model stamped (`bdt`). The capstone's
  gap-1 re-check explicitly updates the dependency from `k8t`+`2m1` to `eq`. ✓
- **Capstone Deliverable A (core rewrite) starts in parallel** — correct.
  Prose-structural reorganization of `fo-dispatch-core.md`; does not depend on
  1+2. The `claude-live`/`codex-live` regression (AC-6) can run as soon as the
  core rewrite is draft. ✓

**Stale section:** the dedicated `## Sequencing` section (index, after the DoD)
still reads "Parallel start: members 1, 2, 3 — independent… After 1–3 land…:
member 4 capstone" and "a dispatched ensign that runs on the parent model
(member 3) with the ensign contract loaded (member 1) from the explicitly-resolved
repo (member 2)." This is the old 4-member layout (three parallel non-capstone
members + a 4th capstone). After the re-carve there are 3 members total and
member 3 IS the capstone; model-stamping is member 2. The re-carve note at the
top of the Members table gives the correct sequencing, but this section was not
updated. A Commander reading it sees "member 4 capstone," which does not exist.
(Owner: Shaping FO — rewrite the `## Sequencing` section to the 3-member
layout.)

## Cross-member composition

**Two seams here; one is clean, one is a new gap the re-carve introduced.**

*Clean seam — `bdt` model stamp → `b2` worker-identity-capture.* The capstone's
capability table binds `worker-identity-capture`'s "stamped model" to
"(`pi-dispatch-model-stamping`)." Review #1's gap 2 (Pi canonical-model-space
declaration ownership) is closed: the capstone's fold-in Gap-2 statement
declares the Pi canonical-model-space declaration OWNED BY the model-stamping
task and REFERENCES (not re-declares) it. (Minor: the fold-in says "member 3"
where it should now say "member 2" — stale number, correct task. Owner: `b2`,
cosmetic.)

*Clean seam — `eq` install → `b2` pi-live drive prereq.* The capstone's gap-1
re-check updates the pi-live drive dependency from `k8t`+`2m1` to `eq`. Correct
and explicit.

*NEW composition gap — the `repoRoot` / `cwd:<repo>` source lost its owner.*

Review #1's gap-1 close wired `cwd:<repo>` into the capstone's `async-dispatch`
binding, "sourced from the same install-recorded / explicitly-resolved repo path
that member 2 records." That source was `2m1`'s install-record repo-path file.
The re-carve archived `2m1`. `eq`'s design:

- D3 makes `spacedock install --host pi` run `pi install` — registering the
  package in `settings.json` `packages` and placing the repo in pi's package
  store. This is a *package registration*, not a *repo-path record*. No
  repo-path file is written.
- D4 retires the cwd fallback but says only "`--plugin-dir` / `SPACEDOCK_REPO_ROOT`
  still resolve the repo for the extension to find skills; the cwd fallback is
  demoted to a last resort… or removed entirely." It does not add a new
  repoRoot resolution source for the installed case.
- AC-3 says the cwd fallback is "removed or demoted below the install-record
  resolution" — but `eq`'s design defines no "install-record resolution." This
  is an internal inconsistency in `eq`'s own AC.

Consequence: under the headline scenario (parent launched from a non-repo cwd,
package installed, no `--plugin-dir`/`SPACEDOCK_REPO_ROOT`), there is no defined
producer of the repo path that (a) the launcher's doctor check
(`firstOfficerOK`/`ensignOK` at `cfg.firstOfficer`/`cfg.ensign`, still computed
from `repoRoot`) needs, and (b) the capstone's `cwd:<repo>` working-directory
argument needs. The capstone's re-check claims the source is "unchanged" — it is
not; `2m1` is gone.

This does not re-break skill discovery (the package-root scan is cwd-independent
— confirmed). It breaks the *working-directory* `cwd:<repo>` forwarding and the
doctor's repo-path skill checks, both of which still want a `repoRoot` that no
member now produces for the installed-from-non-repo case. (Owner: `eq` — specify
how `repoRoot` resolves post-install: e.g. the launcher reads the package store
path from `settings.json` `packages`, OR the doctor's repo-path-based
`firstOfficerOK`/`ensignOK` checks are retired in favor of the package-root
scan. Then `b2` capstone sources `cwd:<repo>` from that same resolved path, and
the re-check's "unchanged" claim becomes true.)

This is closeable without redesign; it is a specification gap, not a mechanism
flaw.

## Blast-radius

**The path→lane mapping is still correct after the re-carve. No gap.**

- **Member 1 (`eq`):** `internal/cli/pi.go` + root `package.json` +
  `.pi/extensions/spacedock.ts` — pi-only surfaces → `pi-live` required. ✓
  (The root `package.json` is a new repo-root artifact carried by all clones,
  but it is inert for Go tooling and claude/codex do not read `package.json`
  for skill discovery — they use their native plugin systems. No cross-host
  blast. Noting only.)
- **Member 2 (`bdt`):** `skills/first-officer/references/pi-first-officer-runtime.md`
  — pi-only → `pi-live` required. The null-model stamp override is Pi-local;
  the core's "OMIT on null" still governs claude/codex. ✓
- **Member 3 (`b2` capstone):** touches every host adapter
  (`claude-first-officer-runtime.md`, `codex-first-officer-runtime.md`,
  `pi-first-officer-runtime.md`, ensign adapters) → `claude-live` AND
  `codex-live` AND `pi-live` required (the dogfood). AC-1 contractlint binds
  capability→tool across all adapters; AC-6 is the claude/codex regression
  gate. Q8 states this correctly. ✓

The "live lane is unrelated" proof-policy burden (Q8) is respected: each
member's lane claim is tied to a concrete diff surface. No path→lane mapping gap
introduced by the re-carve.

## Commander cold-boot readiness

**Q1–Q10 (async, explicit model, skill injection, state branch, path-scoped
commits, back-channel service, completion+file-verify, live lanes, in-flight
state, sandbox) remain correct and re-carve-neutral.** The re-carve does not
touch the operational surface they cover.

**Q11–Q13 — the quirks review #1 demanded — were ADDED but NOT re-validated
against the re-carve. This is the core cold-boot gap.** The index claims gap 3
(Q11–Q13) is "CLOSED (Q11–Q13 added to this index)," but "added" ≠ "re-validated
after the re-carve." All three carry the old 4-member numbering; Q11 is
substantively contradictory.

- **Q11 — STALE AND CONTRADICTORY (must rewrite).** Q11 asserts the capstone's
  AC-2 drive "MUST launch from a non-repo cwd AND pass `cwd:<repo>`" and that
  this "exercises the three-way composition (member 1's cwd-keyed
  `.pi/skills/ensign` symlink + member 2's non-repo-cwd launch + member 4's
  dispatch wiring)," and "if from a non-repo cwd WITHOUT `cwd:<repo>` on the
  dispatch, member 1's symlink is not discovered and ensign doesn't load (DoD
  bullet b re-breaks)." Three problems: (a) "member 1's symlink" — `eq` ships no
  symlink; (b) "member 2's non-repo-cwd launch" and "member 4" — old numbering,
  no such members; (c) the claim that ensign doesn't load without `cwd:<repo>`
  is **directly disproven** by `eq`'s spike (skill discovered from `/tmp`) and
  **explicitly retracted** by the capstone's own gap-1 re-check ("The failure
  mode the fold-in was wired to prevent … is removed at its root"). A Commander
  following Q11 would believe `cwd:<repo>` is load-bearing for ensign discovery
  and could misdiagnose a clean non-repo drive as broken. The non-repo launch
  pin is still valuable (it exercises install-managed discovery and proves no
  cwd dependency) — but Q11's *rationale* and failure claim must be rewritten
  for the re-carve. (Owner: Shaping FO, cold-boot package.)

- **Q12 — STALE NUMBERING + STALE CHECK (must refresh).** Q12's three pre-flight
  checks are numbered to the old layout: "(a) member 1 — subagents-doctor lists
  ensign (project source); (b) member 2 — `spacedock doctor --host pi` shows the
  install-recorded repo source (not 'working directory'); (c) member 3 — a
  null-model probe dispatch stamps the parent's live model." After the re-carve:
  ensign-listing AND repo-source are BOTH `eq` (member 1); model-stamping is
  `bdt` (member 2); the capstone is `b2` (member 3). So (b)'s "member 2" is now
  the model-stamping task, and (c)'s "member 3" is now the capstone. Worse, (b)'s
  specific check ("doctor shows the install-recorded repo source") references
  `2m1`'s mechanism — `eq`'s install is a package registration, and per the
  composition gap `eq`'s `repoRoot` source is currently undefined, so this check
  may not even be satisfiable as worded. (Owner: Shaping FO — re-number to the
  3-member layout and re-base (b) on `eq`'s actual install verification, e.g.
  `settings.json` `packages` contains the Spacedock entry + `subagents-doctor`
  lists ensign as `user-package`.)

- **Q13 — STALE NUMBERING (cosmetic).** "Between member 3 landing and the
  capstone landing" — model-stamping is now member 2, capstone is member 3. The
  contradiction window it describes (core says "OMIT on null," Pi adapter says
  "stamp on null") is real and intentional; only the numbering is stale. (Owner:
  Shaping FO — re-number.)

- **Q3 — STALE NAME (cosmetic).** "Until `pi-ensign-skill-injection` lands" and
  "member 1" — member 1 is now `pi-install-managed-skill-placement`; the archived
  name will confuse a `status --where` query. The quirk's *content* (skill
  injection broken until install-managed lands; verify via `subagents-doctor`)
  is still valid. (Owner: Shaping FO — rename to `pi-install-managed-skill-placement`.)

- **Q9 — STALE NUMBERING (cosmetic).** References "member 4" and
  "`pi-ensign-skill-injection` (member 1)." The in-flight-state guidance
  (preserve `0637e2ed` spike evidence; redispatch `b929622e` clean) is still
  valid; the numbering/naming is stale. (Owner: Shaping FO — re-number/rename.)

**Net:** Q11 is a cold-boot hazard (contradictory + partly false); Q12 is stale
and references a possibly-unsatisfiable check; Q3/Q9/Q13 are stale-name/number
cosmetics. The index's claim that gap 3 is "CLOSED" is not yet earned — the
quirks exist but were not re-validated against the re-carve.

## Gaps to close

1. **[Gap — owner: `eq` (`pi-install-managed-skill-placement`)] Specify the
   `repoRoot` source post-install.** `eq`'s D4 retires the cwd fallback but D3
   (run `pi install`) writes no repo-path install record; `eq`'s AC-3 references
   an "install-record resolution" its own design does not produce. The launcher's
   doctor check (`firstOfficerOK`/`ensignOK` from `cfg.firstOfficer`/`cfg.ensign`
   ← `repoRoot`) and the capstone's `cwd:<repo>` working-directory argument both
   need a defined `repoRoot` for the installed-from-non-repo case. **Close:**
   `eq`'s design specifies how `repoRoot` resolves post-install (e.g. read the
   package store path from `settings.json` `packages`, OR retire the doctor's
   repo-path skill checks in favor of the package-root scan). This also resolves
   `eq`'s AC-3 internal inconsistency. (Composition consumer: `b2` capstone
   sources `cwd:<repo>` from this same resolved path; the re-check's "source is
   unchanged" claim then becomes true.)

2. **[Gap — owner: Shaping FO, cold-boot package] Rewrite Q11 for the re-carve.**
   Q11's claim that ensign doesn't load without `cwd:<repo>` is disproven by
   `eq`'s spike and retracted by the capstone's own re-check; its "member 1
   symlink + member 2 + member 4" composition is the old layout. **Close:**
   rewrite Q11 to reframe `cwd:<repo>` as a working-directory concern (not skill
   discovery), drop the "DoD bullet b re-breaks" failure claim, keep the
   non-repo launch pin (it exercises install-managed discovery), and re-number
   to the 3-member layout. The index's "gap 3 CLOSED" claim depends on this.

3. **[Gap — owner: Shaping FO, cold-boot package] Refresh Q12 for the re-carve.**
   Re-number (a)→`eq`, (b)→`eq` (ensign-listing + install-registered), (c)→`bdt`
   (model-stamping), and re-base (b)'s "doctor shows install-recorded repo source"
   check on `eq`'s actual install verification (`settings.json` `packages` +
   `subagents-doctor` `user-package`), since `2m1`'s install-record repo source
   no longer exists.

4. **[Gap — owner: Shaping FO, index] Rewrite the `## Sequencing` section to the
   3-member layout.** It currently describes the old 4-member layout ("member 4
   capstone"; "parent model (member 3) … ensign (member 1) … repo (member 2)").
   The re-carve note at the top of the Members table has the correct sequencing;
   propagate it into the `## Sequencing` section.

5. **[Gap — owner: Shaping FO, index] Refresh DoD "Proven by" items 1–2 and
   bullet (c) to name `eq` and `eq`'s package-install mechanism** (not the
   archived `k8t`/`2m1` names or `2m1`'s "install-recorded path" phrasing).
   Bullet (c)'s refresh depends on gap 1 landing.

6. **[Gap — owner: `b2` capstone] Reconcile the fold-in AC-2 sub-bullet with the
   gap-1 re-check.** The fold-in AC-2 sub-bullet still reads "verify the ensign
   loads … via member 1's `.pi/skills/ensign` symlink" (archived member, absent
   mechanism). **Close:** add a SUPERSEDED pointer on the sub-bullet itself (or
   reconcile the text to attribute the ensign-loaded pre-check to the
   package-root scan), so a Commander reading AC-2 top-to-bottom is not misled
   before reaching the re-check. (Cosmetic companion: fix the fold-in Gap-2
   "member 3" → "member 2.")

7. **[Cosmetic — owner: Shaping FO, cold-boot package] Refresh stale member
   names/numbers in Q3, Q9, Q13** (`pi-ensign-skill-injection`→`eq`; "member 4"
   / "member 3" → re-carved numbering). Content unchanged; names/numbers only.

No material redesign needed. Gaps 1–3 are the cold-boot-relevant ones (one
composition specification, one contradictory quirk, one stale pre-flight check);
gaps 4–7 are documentation freshness that should land with the
`dispatch-sprint-execution.md` package. Once gap 1 is specified and gaps 2–4 are
rewritten, the sprint is cold-boot drivable.
