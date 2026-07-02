---
title: "contractlint dangler-oracle does not resolve the Skill() cross-file pointer idiom (referenceProsePointerDanglers lags the nt ref->skill move)"
source: "0240 pre-cut antipattern audit (2026-07-01, lens 2). internal/contractlint/boot_resident_closure_test.go referenceProsePointerDanglers resolves a cross-file prose pointer to a watched section ONLY via an intra-file heading or a references/*.md path token (bodyReferenceRe). nt made the correct cross-file idiom for fo-write-core/fo-status-viewer the Skill(skill=spacedock:...) form and deleted the references/*.md paths (now dead per deadDeferredReferencePaths). A future contract line pointing at one of these sections via the Skill() form while naming the section phrase would FALSE-FAIL as dangling. Latent (no such line today); the guard's resolution model lags the idiom it introduced. Non-blocking; captain fast-tracked."
status: validation
group: tooling
id: fezpp9ddsaz6bhek6ra2se69
started: 2026-07-02T01:37:56Z
worktree: .worktrees/spacedock-ensign-contractlint-dangler-oracle-skill-idiom
mod-block: merge:pr-merge
pr: pr-merge:459
---

## Problem
`referenceProsePointerDanglers` (internal/contractlint/boot_resident_closure_test.go ~L288) resolves a cross-file prose pointer to a watched section (e.g. "FO Write Scope", "Status Viewer") only via an intra-file heading or a `references/*.md` path token. nt's ref->skill move made the correct cross-file idiom `Skill(skill="spacedock:fo-write-core")` and deleted the reference paths — so a future line pointing at a watched section via `Skill()` while naming the section phrase would false-fail as dangling.

## Desired fix
Teach the oracle to ALSO resolve a watched-section prose pointer when the same context names the owning skill via the `Skill(skill="spacedock:<name>")` token — mapping the skill to its SKILL.md section anchors (the deferredSkillCores map already keys skill path -> anchors). The guard's resolution model then matches the Skill() idiom nt introduced.

## Proposed approach

Add a THIRD resolution path to `referenceProsePointerDanglers`, alongside the two it already honors (an intra-file heading, a `references/*.md` path token on the line). The new path: a non-heading line that names a watched section resolves if the SAME line also names the section's **owning skill** via the `spacedock:<owner>` token (the `Skill(skill="spacedock:<owner>")` idiom nt introduced). "Owning skill" is derived from `deferredSkillCores`, which already keys `skills/<owner>/SKILL.md → [anchors]`; the reverse (section-name → owner) is `filepath.Base(filepath.Dir(skillPath))` — the exact idiom already used at `TestDeferredSkillCoresResolveAndCarryCeremony` (~L322).

Resolution stays **owner-specific and same-line**, matching the existing path-token check's same-line scope (and the entity source's framing: "a ... line ... while naming the section phrase"). Owner-specificity is the non-weakening property: a Skill token for the WRONG skill does NOT suppress the dangler, so the fix cannot degrade into "any `spacedock:` token silences danglers." A cross-line idiom, if one ever emerges, is a deferred follow-up (YAGNI now).

The change is **additive to the resolution set**, so it can only REDUCE the dangler count for a given body, never increase it — real-file greenness (`TestDeferredSkillProsePointersResolve`) cannot regress from this change. The only regression risk is over-broadening (silencing a genuine dangler), which AC-2's wrong-owner row guards.

The scanner currently takes `(body string, watched []string)`. Since the watched-name set is exactly the keys of the section→owner map, replace the `watched` param with `owners map[string]string` (one source of truth, no redundancy); `watchedSectionNames()` becomes `sectionOwners()`. Both callers (`TestDeferredSkillProsePointersResolve`, the control) pass `sectionOwners()`, preserving the "single scanner drives real guard + control" property. Map-iteration order makes multi-dangler-per-line message order nondeterministic; no test depends on order (all check count/presence), so this is cosmetic only.

### Before

```go
func watchedSectionNames() []string {
	var names []string
	seen := map[string]bool{}
	for _, anchors := range deferredSkillCores {
		for _, anchor := range anchors {
			name := strings.TrimLeft(anchor, "# ")
			if !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	return names
}

func referenceProsePointerDanglers(body string, watched []string) []string {
	lines := strings.Split(body, "\n")
	defined := map[string]bool{}
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") {
			continue
		}
		heading := strings.TrimLeft(line, "# ")
		for _, name := range watched {
			if heading == name {
				defined[name] = true
			}
		}
	}
	var danglers []string
	for i, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}
		hasPath := bodyReferenceRe.MatchString(line)
		for _, name := range watched {
			if !strings.Contains(line, name) || defined[name] || hasPath {
				continue
			}
			danglers = append(danglers, fmt.Sprintf("line %d names section %q but resolves neither intra-file (no `## %s` heading here) nor via a references/*.md token: %q", i+1, name, name, strings.TrimSpace(line)))
		}
	}
	return danglers
}
```

### After

```go
// sectionOwners maps each watched section name to the skill that owns it
// (skills/<owner>/SKILL.md), derived from deferredSkillCores so a newly-registered
// anchor is watched AND owner-resolved from one place. The section-name set the
// prose-pointer scanner watches is exactly this map's keys.
func sectionOwners() map[string]string {
	owners := map[string]string{}
	for skillPath, anchors := range deferredSkillCores {
		owner := filepath.Base(filepath.Dir(skillPath)) // skills/<owner>/SKILL.md -> <owner>
		for _, anchor := range anchors {
			owners[strings.TrimLeft(anchor, "# ")] = owner
		}
	}
	return owners
}

func referenceProsePointerDanglers(body string, owners map[string]string) []string {
	lines := strings.Split(body, "\n")
	defined := map[string]bool{}
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") {
			continue
		}
		heading := strings.TrimLeft(line, "# ")
		if _, watched := owners[heading]; watched {
			defined[heading] = true
		}
	}
	var danglers []string
	for i, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}
		hasPath := bodyReferenceRe.MatchString(line)
		// skills named on this line via the spacedock:<name> token (the Skill(skill="…") idiom)
		skillsNamed := map[string]bool{}
		for _, m := range bodySkillRe.FindAllStringSubmatch(line, -1) {
			skillsNamed[m[1]] = true
		}
		for name, owner := range owners {
			if !strings.Contains(line, name) || defined[name] || hasPath {
				continue
			}
			if skillsNamed[owner] {
				continue // resolves cross-file via Skill(skill="spacedock:<owner>")
			}
			danglers = append(danglers, fmt.Sprintf("line %d names section %q but resolves neither intra-file (no `## %s` heading here), via a references/*.md token, nor via Skill(skill=\"spacedock:%s\"): %q", i+1, name, name, owner, strings.TrimSpace(line)))
		}
	}
	return danglers
}
```

Callers change `watchedSectionNames()` → `sectionOwners()` and the vacuity guard `if len(watched) == 0` → `if len(owners) == 0`.

## Acceptance criteria

- **AC-1 — a Skill-form owner-matched pointer resolves (the reason the task exists).** For a non-heading line that names a watched section (e.g. `FO Write Scope`) and names its owning skill via the `spacedock:fo-write-core` token, `referenceProsePointerDanglers` returns zero danglers.
  Test: a new row in `TestDeferredSkillProsePointerGuardFailsOnDanglingTarget` asserting `len(danglers) == 0`. RED-then-GREEN: the row is confirmed to FAIL against the current scanner (returns ≥1 — the latent false-fail is real) before the fix, and pass after. The pre-fix count (≥1) is the independent baseline; the AC asserts it moves to 0.

- **AC-2 — value measure: the guard classifies a 4-case battery 4/4 correctly and stays owner-specific (non-weakening).** Over the fixture battery, the count of correctly-classified cases is 4/4: bare-name pointer → dangles (≥1); wrong-owner Skill token (`Skill(skill="spacedock:fo-status-viewer")` on a line naming `FO Write Scope`) → dangles (≥1); intra-file heading → resolves (0); Skill-form owner-matched → resolves (0). The two SHOULD-dangle counts are baselines that a too-broad fix (any `spacedock:` token suppresses) moves the wrong way to 0; this AC forbids that. This is the measured end-value the entity exists for — it brackets correct behavior from both sides (AC-1 catches "did nothing"; the wrong-owner row catches "did too much").
  Test: table-driven `TestDeferredSkillProsePointerGuardFailsOnDanglingTarget`, one assertion per row (`> 0` or `== 0`).

- **AC-3 — real deferred-skill bodies stay green; no regression.** `go test ./internal/contractlint/` passes: `TestDeferredSkillProsePointersResolve` reports zero danglers over the actual `fo-status-viewer/SKILL.md` and `fo-write-core/SKILL.md`, and `TestDeferredSkillCoresResolveAndCarryCeremony` + `TestNoSurvivingContractFileNamesDeadDeferredReferencePath` remain green. (Safe by construction: the change is additive to the resolution set, so it can only reduce danglers on the real bodies.)
  Test: `go test ./internal/contractlint/` → `ok`.

## Test plan

- **What verifies it:** Go unit tests in `internal/contractlint/boot_resident_closure_test.go`. The scanner is a pure function over in-memory string fixtures, so no on-disk fixtures, no CLI driver, and no live-workflow test are needed.
- **Extend the control** `TestDeferredSkillProsePointerGuardFailsOnDanglingTarget` into a table-driven battery with the existing three rows (bare-name RED, intra-file heading GREEN, path-token GREEN) plus two new rows: (4) Skill-form owner-matched `"... FO Write Scope ... `+"`Skill(skill=\"spacedock:fo-write-core\")`"+` ..."` → expect 0; (5) wrong-owner `"... FO Write Scope ... `+"`Skill(skill=\"spacedock:fo-status-viewer\")`"+` ..."` → expect ≥1.
- **RED-then-GREEN order:** add rows 4 and 5 first. Row 4 must FAIL against current code (confirms the latent false-fail); row 5 already passes against current code and must STAY passing after the fix (confirms owner-specificity — the fix is not a blanket suppressor). Apply the fix; both rows plus the three originals go green.
- **Cost / complexity:** trivial. Pure-function unit test; full package runs in ~0.3s (`ok internal/contractlint 0.343s` observed at baseline). Low complexity.
- **Fixture / CLI / live:** none. In-memory string fixtures suffice.

## Spike (riskiest path — exercised during ideation)

Riskiest unverified bit: does the resolver's mechanism exist without new parsing — (a) does `deferredSkillCores` already key skill-path → section-anchors so the section→owner reverse map is derivable, and (b) can the resolver key on the parsed `Skill()` token with an existing regex? Exercised via a throwaway replicating `bodySkillRe` + `deferredSkillCores` (scratchpad `spike_regex.go`, `go run`):
- `bodySkillRe` (`spacedock:([a-z0-9-]+)`, already defined at ~L67 and already used by `extractDeferredLoadPoints`) captures `fo-write-core` from `` `Skill(skill="spacedock:fo-write-core")` `` — the `Skill(skill="…")` wrapper is transparent to the token regex. It also captures a bare `spacedock:fo-status-viewer`. No new regex or parser is needed.
- The section→owner reverse map derives cleanly: all five anchors map to their owning skill (`FO Write Scope`/`ID Styles` → `fo-write-core`; `Status Viewer`/`Captain-Facing State Display`/`Issue Filing` → `fo-status-viewer`) via `filepath.Base(filepath.Dir(skillPath))`, the idiom already at ~L322.
- Baseline `go test ./internal/contractlint/` is green before any change.

Finding: the mechanism is fully proven from existing code — the fix is a small additive branch plus a derived reverse map, with no unverified parser/round-trip risk remaining. Assumption on record: section NAMES are unique across the two deferred skills (no anchor phrase collides), so `sectionOwners()` never overwrites an owner; the current anchor set satisfies this.

## Docs

No user-visible behavior change — this is an internal contractlint test-guard; no CLI output, banner, host integration, or docs-site surface changes. No doc diff required.

## Stage Report: ideation

- DONE: Flesh the design: teach referenceProsePointerDanglers (boot_resident_closure_test.go) to ALSO resolve a watched-section prose pointer when the context names the owning skill via Skill(skill="spacedock:<name>"), mapping the skill to its SKILL.md anchors via the existing deferredSkillCores map — concrete before/after of the resolution logic
  Proposed approach with full Before/After code of `sectionOwners()` + `referenceProsePointerDanglers` in the entity body: additive same-line owner-specific branch reusing `bodySkillRe`.
- DONE: Tighten into measured ACs each naming its test: RED-then-GREEN (a Skill-form pointer at a watched section resolves; a genuine dangler with no heading/path/Skill token still reds; existing dangling-target + prose-pointer controls stay RED-capable)
  AC-1 (Skill-form resolves, RED→GREEN), AC-2 (value measure: 4/4 battery classification incl. wrong-owner non-weakening row), AC-3 (`go test ./internal/contractlint/` green, no regression) — each names its test.
- DONE: Spike the riskiest bit: confirm deferredSkillCores already keys skill-path->section-anchors and the resolver can key on the parsed Skill token (quick read of the helper); record the finding
  Exercised via scratchpad `spike_regex.go` (`go run`): `bodySkillRe` captures the name inside `Skill(skill="spacedock:…")`; reverse map derives via `filepath.Base(filepath.Dir(skillPath))`; baseline suite green. Recorded in the Spike section.

### Summary

Fleshed the latent-dangler fix into a concrete additive resolution path: a watched-section pointer resolves if the same line names the section's owning skill via the `spacedock:<owner>` token, with owner-specificity as the non-weakening property. Key decision: replace the scanner's `watched []string` param with `owners map[string]string` (keys = watched set) for one source of truth, keeping the single-scanner-drives-control property. Spike proved the mechanism needs no new regex/parser — `bodySkillRe` and `deferredSkillCores` already supply both halves — so no unverified round-trip risk remains; the change is monotonic (can only reduce danglers), making real-file greenness safe by construction.

## Stage Report: implementation

- DONE: Add the resolver branch in referenceProsePointerDanglers — build a section->owner reverse map from deferredSkillCores, key on the existing bodySkillRe token; ADDITIVE, owner-specific (a wrong-owner token must NOT suppress)
  `sectionOwners()` replaces `watchedSectionNames()` (`filepath.Base(filepath.Dir(skillPath))` reverse map); resolver skips a dangler only when `skillsNamed[owner]` matches the same-line `spacedock:<owner>` token. Worktree commit da111ffc.
- DONE: Extend TestDeferredSkillProsePointerGuardFailsOnDanglingTarget into the table-driven battery: row 4 (Skill-form owner-matched -> 0) and row 5 (wrong-owner -> >=1); RED-then-GREEN — row 4 FAIL pre-fix, pass post-fix; row 5 STAYS passing
  5-row battery; pre-fix run showed only row 4 failing (`got 1, expected 0`) with rows 1-3+5 passing (RED confirmed the latent false-fail); post-fix all 5 green, row 5 owner-specific.
- DONE: go test ./internal/contractlint/ green — including TestDeferredSkillProsePointersResolve, TestDeferredSkillCoresResolveAndCarryCeremony, TestNoSurvivingContractFileNamesDeadDeferredReferencePath; commit in the worktree
  Full package `ok` (0.338s); the three named tests + the battery all PASS; `go vet` clean; committed da111ffc on branch spacedock-ensign/contractlint-dangler-oracle-skill-idiom.

### Summary

Added a third, owner-specific resolution path to `referenceProsePointerDanglers`: a watched-section prose pointer now resolves when the same line names its owning skill via the `Skill(skill="spacedock:<owner>")` token nt introduced, closing the latent false-fail. Replaced `watchedSectionNames() []string` with `sectionOwners() map[string]string` as the single source of truth (keys = watched set, values = owner), keeping the single-scanner-drives-control property; both callers updated. The change is additive and monotonic, so it can only reduce danglers — real deferred-skill bodies stay green — while owner-specificity (verified by the wrong-owner row 5) prevents any `spacedock:` token from blanket-silencing a genuine dangler.

## Stage Report: validation

- DONE: MEASURE AC-1 + AC-2 (the 4/4 battery, reproduce not assert): row 4 (Skill-form owner-matched pointer -> 0 danglers) AND confirm it FAILED pre-fix on the current-main scanner (the latent false-fail is real = baseline moves the wrong way); row 5 (wrong-owner Skill token on a line naming a different section -> >=1) STAYS red (owner-specific, non-weakening); bare-name -> >=1, intra-file heading -> 0
  Reproduced raw counts by driving the REAL scanners on throwaway checkouts (never the impl worktree). Pre-fix (main 59251f18): row4=**1 danglers (RED — latent false-fail confirmed)**, row1=1, row2=0, row3=0, row5=1. Post-fix (da111ffc): row4=**0 (RED→GREEN)**, row5=1 (owner-specific), row1=1, row2=0, row3=0. 4/4 classified correctly both sides.
- DONE: AC-3: go test ./internal/contractlint/ green — TestDeferredSkillProsePointersResolve reports 0 danglers over the real fo-status-viewer/fo-write-core SKILL.md bodies, TestDeferredSkillCoresResolveAndCarryCeremony + TestNoSurvivingContractFileNamesDeadDeferredReferencePath green
  `go test ./internal/contractlint/ -count=1` → `ok 0.386s`; all four named tests PASS (`-v`); `go vet` exit 0. Change scoped to one `_test.go` file (+98/-55), no production code.
- DONE: Confirm the fix is ADDITIVE + owner-specific: a blanket "any spacedock: token suppresses" over-fix would flip row 5 to 0 — verify row 5 stays >=1 (the built-in over-fix guard bites)
  Constructed the blanket over-fix and ran the battery on the throwaway: it flips row4→0 AND row5→**0** (silences the genuine dangler), so the wrong-owner row would FAIL under it — the guard is non-vacuous. Real fix keeps row5=1. `sectionOwners()` derives 5 distinct anchors → 2 owners with no collision (Spike uniqueness assumption holds).

### Summary

PASSED. Independently reproduced the RED→GREEN baseline on throwaway checkouts of main and the fix commit (not the impl worktree): the Skill-form owner-matched pointer (row 4) genuinely false-fails on the current-main scanner (1 dangler) and resolves to 0 under the fix, while the wrong-owner row (row 5) stays at 1 on both sides — the owner-specific, non-weakening property. Proved the over-fix guard is non-vacuous by constructing the blanket "any spacedock: token suppresses" variant and observing it flips row 5 to 0, a regression the wrong-owner row would catch. AC-3 real-body greenness holds: `go test ./internal/contractlint/` is `ok`, all four named tests pass, `go vet` clean; the change is test-file-only (additive, monotonic). Recommendation: PASSED.
