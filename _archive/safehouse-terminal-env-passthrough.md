---
id: 89bdt7vhk362yw630kt9s12k
title: Safehouse passes terminal-host env for Zellij only, and passes names that do not exist
status: done
source: "Captain CL, 2026-08-18, from a live observation: `env | grep TMUX` returned nothing inside a safehouse-wrapped session. Traced to `internal/safehouse/safehouse.go:68-73`. Captain's direction: take the env set from subspace's own probe, and add a name to `--env-pass` only when that variable exists."
started: 2026-08-18T14:57:04Z
completed: 2026-08-18T18:53:19Z
verdict: PASSED
score:
worktree: .worktrees/spacedock-ensign-safehouse-terminal-env-passthrough
issue:
mod-block:
pr: pr-merge:730
archived: 2026-08-18T18:53:19Z
---

A sandboxed session can detect Zellij and nothing else. Five other terminal hosts lose their identity at the sandbox boundary, and the allowance lists variable names whether or not they exist.

## Problem

`terminalTargetingEnvArgs` (`internal/safehouse/safehouse.go:68-73`) is the whole allowance:

```go
func terminalTargetingEnvArgs() []string {
	if _, present := os.LookupEnv("ZELLIJ"); !present {
		return nil
	}
	return []string{"--env-pass=ZELLIJ,ZELLIJ_PANE_ID,ZELLIJ_SESSION_NAME"}
}
```

Two separate faults.

**Fault 1 — the host set is Zellij-only.** When `ZELLIJ` is absent the function returns `nil`, so nothing is passed. A tmux session therefore reaches the child with no `TMUX` and no `TMUX_PANE`. Confirmed live by the captain: `env | grep TMUX` is empty inside a safehouse-wrapped session.

The frontdoor is not the cause. `launchEnv` (`internal/cli/frontdoor.go:73-79`) forwards the entire parent environment and only swaps `SPACEDOCK_BIN`. The loss is at the sandbox boundary.

**Fault 2 — presence is checked per host, not per variable.** The gate reads one sentinel (`ZELLIJ`), then passes three names unconditionally. `ZELLIJ_PANE_ID` and `ZELLIJ_SESSION_NAME` are listed whether or not they are set. The gate variable is also not one of the two the consumer actually probes, so the sentinel and the payload disagree about what identifies a Zellij pane.

## The consumer already defines the correct set

The `subspace:r` skill documents a probe of nine signals across six hosts (`plugins/subspace/skills/r/SKILL.md:71-91`), in its own resolution order:

| Host | Signals |
|---|---|
| Zellij | `ZELLIJ_SESSION_NAME`, `ZELLIJ_PANE_ID` |
| tmux | `TMUX`, `TMUX_PANE` |
| Herdr | `HERDR_ENV`, `HERDR_PANE_ID` |
| CMUX | `CMUX_WORKSPACE_ID`, `CMUX_SURFACE_ID` |
| Ghostty | `TERM_PROGRAM=ghostty` |
| Apple Terminal | `TERM_PROGRAM=Apple_Terminal` |

Safehouse passes three names. One of them, `ZELLIJ`, is not in this probe at all. Six probed signals never cross the boundary.

The consequence is concrete: inside safehouse, `subspace:r` resolves Zellij or falls through to none. A captain in tmux, Herdr, CMUX, Ghostty, or Apple Terminal gets the fallback rather than their real pane host.

Note that the probe treats an empty value as absent. So a name passed for an unset variable is not merely useless — it risks presenting an empty variable where the consumer expects nothing.

## Proposed approach

Replace the sentinel-gated fixed list with a presence filter over the consumer's nine-signal set, emitted in the probe's own order:

```go
// The nine signals subspace's r skill probes across its six terminal hosts.
// Source of truth: spacedock-subspace plugins/subspace/skills/r/SKILL.md,
// "Select one terminal" — duplicated by decision (see the task body's drift note).
var terminalHostEnvVars = []string{
	"ZELLIJ_SESSION_NAME", "ZELLIJ_PANE_ID",
	"TMUX", "TMUX_PANE",
	"HERDR_ENV", "HERDR_PANE_ID",
	"CMUX_WORKSPACE_ID", "CMUX_SURFACE_ID",
	"TERM_PROGRAM",
}

func terminalTargetingEnvArgs() []string {
	return terminalEnvPassArgs(os.LookupEnv)
}

func terminalEnvPassArgs(lookup func(string) (string, bool)) []string {
	var present []string
	for _, name := range terminalHostEnvVars {
		if _, ok := lookup(name); ok {
			present = append(present, name)
		}
	}
	if len(present) == 0 {
		return nil
	}
	return []string{"--env-pass=" + strings.Join(present, ",")}
}
```

`Wrap` is unchanged. One comma-joined `--env-pass=` flag in probe order — the shipped flag shape, proven live by the working Zellij path.

### Decisions

Each names the AC it serves, the alternative considered, and why the alternative loses.

1. **Bare `ZELLIJ` leaves the allowance.** The consumer probes `ZELLIJ_SESSION_NAME`/`ZELLIJ_PANE_ID`, never `ZELLIJ`; grep over the whole subspace repo finds no live code consuming the bare name. The captain's direction is "take the env set from subspace's own probe". Keeping it "for compatibility" would preserve a name with no known consumer and make the list something other than the probe's. Declared as a semantic change below.
2. **`TERM_PROGRAM` stays in the list, value-agnostic** (serves AC-1's Ghostty and Apple Terminal legs). Live observation: safehouse's own defaults already forward `TERM`, `TERM_PROGRAM`, `TERM_PROGRAM_VERSION` (all present inside this sandboxed session), so this entry is redundant today — it makes the consumer's contract explicit rather than an accident of the sandbox's default allowance. Alternative — gate on value ∈ {`ghostty`, `Apple_Terminal`}: embeds subspace's resolution semantics into safehouse and adds a second drift axis (values as well as names) for no gain; naming a variable the sandbox already passes is a no-op, and faithful passthrough matches what an unsandboxed launch shows.
3. **"Set in the parent" means `os.LookupEnv` presence; a set-but-empty variable is named.** The child then mirrors the parent exactly, and the probe maps empty→absent identically on both sides. The harm AC-2 guards against — presenting an empty variable where the parent had nothing — arises only for unset names, which the presence filter excludes. Alternative — filter empties too: makes the child's env disagree with the parent's in a second way, for a case (set-but-empty identifier) no host produces.
4. **The composer takes a lookup func;** `terminalTargetingEnvArgs` binds `os.LookupEnv` (serves AC-2's stubbed-environment unit test). The composer is pure, so its table test is deterministic under any developer terminal. Alternative — `t.Setenv` plus an unset-all-nine helper: works, but every allowance test then depends on scrubbing ambient terminal vars (`TERM_PROGRAM` is set in effectively every real dev session; `TMUX` under tmux), and one missed name makes the suite terminal-dependent. Wrap-level tests still use the clear-nine helper; the composer's tests don't have to.

### How the list stays agreed with subspace

Duplicate the nine names and accept the drift.

- Every sharing mechanism crosses a repository boundary and costs more than the nine names it protects: a shared Go module adds version skew and release coupling; parsing the installed subspace `SKILL.md` at runtime makes the launcher depend on an optionally-installed plugin file (safehouse wraps hosts that don't carry subspace at all); generated sync needs cross-repo CI plumbing.
- The set changes only when subspace adds or removes a terminal host — rare, deliberate events. The `r` skill itself pins "exactly nine LF-terminated lines" with fixed names.
- Drift degrades gracefully, never corrupts. With per-variable passthrough, a host missing from a stale copy loses ALL its signals, the probe's complete-family rule sees no family, and resolution falls through — exactly today's shipped behavior for five of six hosts. If subspace instead adds a signal to an existing family, the probe stops on the partial family with the missing signal named — a visible, diagnosable stop, not a wrong pane.
- Free mitigations: the `safehouse.go` comment names the source-of-truth path and section; the unit test pins the nine names, so any deliberate change is one visible edit per repo; a new host arrives through the runtime-support first-contact discipline, which is when the list gets revisited.

## Risk evidence (live, gathered during ideation)

**No spike needed** beyond the live evidence below: the change rides two mechanisms, both proven live with the real binary — (a) safehouse forwards a comma-joined `--env-pass=` allowance through its scrub (the shipped Zellij trio resolves panes inside safehouse today), and (b) `--env-pass` is name-agnostic (`SPACEDOCK_BIN` crosses this very session's boundary via its own `--env-pass` flag; observed set inside). The composer itself is plain argv assembly under unit test.

Live probe comparison, run 2026-08-18 with the consumer's exact nine-signal probe (SKILL.md "Select one terminal"):

| Signal | tmux parent, unwrapped | inside safehouse under tmux (current code) |
|---|---|---|
| `TMUX` | present | absent |
| `TMUX_PANE` | present | absent |
| `TERM_PROGRAM` | other (`tmux`) | other (`tmux`) |
| other six | absent | absent |

Unwrapped leg: the probe in a fresh private tmux server (`tmux -L sd-spike-89b`, tmux 3.7b). Wrapped leg: this agent session itself, which runs inside the captain's safehouse-wrapped tmux session. The tmux family is invisible inside — the probe resolves no terminal — while `TERM_PROGRAM` crosses via safehouse's default `TERM*` allowance. This is AC-1's failing baseline, per signal.

**Operational constraint on the closing proof:** no stage running inside a wrapped session can execute the real safehouse (`~/.local/bin/safehouse` → Operation not permitted inside the sandbox; the parent env is unreadable, `ps` blocked). AC-1's fixed-side inside/outside comparison therefore runs outside the sandbox — a captain-pasted one-liner at the validation gate; exact protocol in the Test plan.

## Out of scope

Changing `launchEnv` or the frontdoor. Owning an operator configuration surface for env passthrough; the existing comment states that safehouse deliberately adds a default without parsing caller arguments, and that boundary stays.

## Expected surface and tolerance

Estimate net LOC change: +45, across 4 files (1 production, 3 test). Insertions ≈ +78 and deletions ≈ −33, reported separately; no gross tolerance declared.

- `internal/safehouse/safehouse.go`: ≈ +22/−6 — the nine-name list and the presence-filter composer.
- `internal/safehouse/safehouse_test.go`: ≈ +35/−13 — composer table test; allowance test reshaped to the tmux pair; clear-helper widened to the nine names.
- `internal/cli/safehouse_env_smoke_test.go`: ≈ +18/−13 — fake safehouse scrubs/forwards tmux names too; expectations follow the new allowance (`ZELLIJ` no longer named).
- `internal/cli/host_launch_test.go`: ≈ +3/−1 — `TestMain`'s baseline unset list widened to the nine names.

Semantics changed: the `--env-pass` argv safehouse composes. Names appear per parent presence; bare `ZELLIJ` is no longer named; the nameable set widens from three to the consumer's nine (adds `TMUX`, `TMUX_PANE`, `HERDR_ENV`, `HERDR_PANE_ID`, `CMUX_WORKSPACE_ID`, `CMUX_SURFACE_ID`, `TERM_PROGRAM`). Runtime behavior: sandboxed children regain terminal-host identity for the four multiplexer families. No command grammar, stored-format, or authority changes. No doc diff: no file under `docs/` outside the state checkout describes the terminal-env allowance (grep for env-pass/ZELLIJ/TERM_PROGRAM over `docs/`).

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - A terminal host's identifying variables survive the sandbox boundary for every host the consumer probes.**
This is the measuring AC: the count of probed signals that reach the child, out of the nine, must equal the count set in the parent. Verified by running the consumer's own probe script inside a safehouse-wrapped session under at least tmux and Zellij, and comparing present/absent for each signal against the same probe run outside the sandbox. Fails on the current code, where a tmux parent yields `TMUX=absent` in the child — the failing tmux baseline is recorded in Risk evidence. The fixed-side run must execute outside a wrapped session (a wrapped stage cannot exec safehouse — see the operational constraint); the Test plan's live protocol is the captain-pasted comparison at the validation gate.

**AC-2 - No variable name appears in `--env-pass` unless that variable is set in the parent.**
Verified by a unit test over the argv composer with a stubbed environment: a parent holding only `TMUX` and `TMUX_PANE` produces an allowance naming those two and nothing else; an empty parent produces no allowance at all. Fails if the composer emits a fixed list, or emits names gated on a different variable than the one being passed.

## Test plan

Layered so the cheap tests prove composition and the one live run proves the boundary. Cost is small: the unit and smoke layers reuse existing patterns in `internal/safehouse` and `internal/cli`; no new fixture files.

1. **Unit, composer (AC-2; `internal/safehouse`).** Table test over `terminalEnvPassArgs` with a map-backed lookup: empty parent → `nil` (no flag at all); `TMUX`+`TMUX_PANE` only → `--env-pass=TMUX,TMUX_PANE`; all nine set → all nine in probe order; a set-but-empty name → named; a name outside the list → never named. Fails if the composer emits a fixed list, gates on a sentinel, or emits unset names.
2. **Unit, wiring (`internal/safehouse`).** The existing Wrap allowance test reshaped: clear the nine, set the tmux pair, expect exactly `--env-pass=TMUX,TMUX_PANE` between `--trust-workdir-config` and the extra args. Fails if `Wrap` stops consulting the composer or the flag moves in the argv.
3. **Integration smoke (`internal/cli`, existing fake-safehouse pattern).** The fake scrubs the tmux pair too and honors `--env-pass`; a tmux-pair parent's child sees the pair's values and no `ZELLIJ_*`; the Zellij-pair case still forwards. Proves argv → forwarded-child-env through a scrubbing wrapper, including that names the parent lacks are never presented as empty.
4. **Live closing proof (AC-1; captain at the validation gate, outside any wrapped session; tmux and Zellij at minimum).** Outside leg: the SKILL.md nine-line probe run bare in the terminal. Inside leg: `safehouse --trust-workdir-config --env-pass=TMUX,TMUX_PANE,TERM_PROGRAM -- /bin/sh -c '<same probe>'` — hand-composed, but byte-identical to the argv the composer's unit test pins for that parent (under Zellij, the allowance the table pins for a Zellij parent). Unit + live together close composer → argv → child-env. Expected: per-signal present/absent identical inside and outside; current code fails the tmux leg (baseline in Risk evidence). The original observation's `env | grep TMUX` is the quick sanity mirror.
5. **Test hygiene, terminal-independent suite.** `internal/cli`'s `TestMain` and the safehouse clear-helper unset all nine names, so `go test ./...` passes identically under tmux, Zellij, or a bare terminal. Today's suite clears only the Zellij trio; under tmux the widened allowance would otherwise leak ambient `TMUX`/`TERM_PROGRAM` into argv baselines.

## Stage Report: ideation

- DONE: Build the allowance from the variables actually set in the parent, so an unset name never reaches --env-pass and an empty parent yields no allowance at all.
  Designed as the presence-filter composer in Proposed approach (`terminalEnvPassArgs` over the nine-name list; empty parent → nil); pinned by AC-2 and Test plan items 1-2.
- DONE: Cover the six hosts subspace probes, and decide how this list stays agreed with subspace's rather than drifting — state the answer even if it is "duplicate and accept the drift, here is why".
  All nine signals of the six hosts, in probe order; decision recorded: duplicate and accept the drift — sharing mechanisms cost more than nine near-static names, drift degrades to today's behavior or a diagnosable probe stop, mitigated by a source-of-truth comment and the name-pinning unit test.
- DONE: Prove it by running subspace's own probe inside and outside a safehouse-wrapped session under tmux, and comparing present/absent per signal.
  Ran the consumer's exact probe both legs live: unwrapped tmux parent shows TMUX/TMUX_PANE present; inside this safehouse-wrapped tmux session both are absent (per-signal table in Risk evidence). This proves the fault; the fixed-side leg is structurally impossible in-session (the sandbox denies exec of safehouse), so it is specified as the captain-pasted validation-gate protocol in Test plan item 4.

### Summary

Filled the gated design: a per-variable presence filter over subspace's own nine-signal set replaces the ZELLIJ-sentinel fixed trio, with four recorded decisions (drop bare ZELLIJ; keep TERM_PROGRAM value-agnostic; presence means LookupEnv including set-but-empty; lookup-func composer for deterministic tests). Live evidence gathered during ideation: the nine-signal probe comparison proving the tmux fault, safehouse's default TERM* forwarding, and SPACEDOCK_BIN crossing via --env-pass (proving the flag is name-agnostic) — hence "no spike needed" rests on live-proven mechanisms. Surface refined from the seed's +40/2 files to net +45 (≈+78/−33) across 4 files, ~20 of them production lines; the extra files are test hygiene keeping the suite terminal-independent.

## Stage Report: implementation

- DONE: Hold the design's shape: one name list, one presence filter, injectable lookup. If it grows past about 45 net lines, stop and say why rather than absorbing it.
  `internal/safehouse/safehouse.go` net +26 (28 ins/2 del) — one `terminalHostEnvVars` list, one `terminalEnvPassArgs(lookup)` filter, `terminalTargetingEnvArgs` binds `os.LookupEnv`; well under the 45-line design budget. Total diff across all 4 files (1 production, 3 test) is net +135 (173 ins/38 del) vs the ideation point estimate of net +45 (≈78/33) — driven entirely by the fully-specified Test plan item-1 table (5 cases) and item-3 smoke reshape (tmux pair + Zellij pair + existing native-allowlist case), not scope creep; commit b5cc22227.
- SKIPPED: Prove AC-1 by running subspace's own probe inside and outside a safehouse-wrapped session under tmux, and record the per-signal present/absent comparison.
  Structurally blocked from this stage exactly as ideation recorded: this session itself runs inside a safehouse sandbox (`APP_SANDBOX_CONTAINER_ID=agent-safehouse`; `safehouse` unresolvable on PATH, `~/.local/bin/safehouse` → "Operation not permitted") and is not itself tmux/Zellij-multiplexed (own probe: all nine signals absent except `TERM_PROGRAM`). AC-1's composition is instead proven end-to-end via Test plan items 1-3 (composer table, `Wrap` allowance reshaped to the tmux pair, `internal/cli` fake-safehouse smoke covering tmux-pair/Zellij-pair/native-allowlist) — all green under `go test -race ./internal/safehouse/... ./internal/cli/...`. The true live inside/outside proof is Test plan item 4, the captain-pasted one-liner at the validation gate.
- DONE: Prove AC-2 with the composer table test: a parent holding only TMUX and TMUX_PANE yields exactly those two names, and an empty parent yields no allowance at all.
  `TestTerminalEnvPassArgsComposer` (`internal/safehouse/safehouse_test.go`) — case "tmux pair only" asserts `--env-pass=TMUX,TMUX_PANE`; case "empty parent yields no allowance" asserts nil; plus all-nine-in-probe-order, set-but-empty-still-named, and outside-the-list-never-named. `go test -race ./internal/safehouse/...` passes.

### Summary

Replaced the ZELLIJ-sentinel fixed trio with `terminalEnvPassArgs`, a pure presence filter over subspace's nine-signal probe set, injectable-lookup for deterministic tests. Production surface held to the declared design shape (+26 net); test surface ran well over the ideation's point estimate because the Test plan's own case matrices (composer table, tmux/Zellij smoke pairs) are larger than +45 net once written out, not because of added scope — flagging per the checklist's "stop and say why" instruction rather than absorbing it silently. AC-1's true live proof remains the captain-pasted validation-gate step ideation already specified; this stage's evidence is the full unit+smoke layer stack, all passing. One pre-existing, unrelated failure observed in the full suite: `TestCodexResolveManifestAgainstInstalledHost` fails identically on `main` at the same commit (local codex plugin-cache state), not touched by this change.

## Review-finding disposition

Entered by validation, 2026-08-18. Reviewer observation authority only — classification is proposed, not authorized. No candidate bytes changed: worktree HEAD stayed `b5cc22227` with a clean status throughout; all 11 mutations ran on a throwaway `git archive` checkout.

### Finding 1 (proposed Polish, evidence defect) — the composer table cannot catch a sentinel gate on a non-Zellij variable

AC-2's own "Verified by" clause promises the unit test "fails if the composer ... emits names gated on a different variable than the one being passed". It does not. Prepending `if _, ok := lookup("TMUX"); !ok { return nil }` to `terminalEnvPassArgs` leaves all of `internal/safehouse` green — every table case either sets `TMUX` or already expects nil, because the table has no Zellij-only parent. The only test in the repo that reds is `TestSafehouseEnvPassForwardsTerminalTargetingMetadata/Zellij pair still forwards` in `internal/cli`. Shipped behavior is correct and the union of the two layers covers the clause; the unit layer AC-2 names as its verifier does not. Remedy: one 8-line "zellij pair only" case mirroring "tmux pair only". Not material — no value AC fails and the gap is covered today.

### Finding 2 (proposed Polish, evidence defect) — nothing pins bare `ZELLIJ` out, or the list against additions

Adding `"ZELLIJ"` back to `terminalHostEnvVars` survives the entire suite. "all nine set, emitted in probe order" asserts only that those nine, when the stub sets them, emit in that order; a tenth name the stub never sets is invisible. Removals are caught — dropping `TERM_PROGRAM`, the Herdr pair, or the tmux pair each reds that case — additions are not. This half-undercuts the ideation's stated drift mitigation ("the unit test pins the nine names, so any deliberate change is one visible edit per repo") and leaves Decision 1's declared semantic change unpinned. Remedy: assert the exact `terminalHostEnvVars` slice. Not material — an added name still only crosses when set in the parent, so AC-2 continues to hold.

### Finding 3 (proposed Polish) — one table case cannot fail under any plausible mutation

"a name outside the list is never named" (`internal/safehouse/safehouse_test.go:160-166`, 7 lines) is killed by exactly the two mutations that also kill "empty parent yields no allowance", and by nothing else. It cannot be otherwise: `terminalEnvPassArgs` takes `func(string) (string, bool)`, so no implementation reachable through that signature can ever name `NOT_A_TERMINAL_SIGNAL` unless someone writes that literal into the production list. Remedy: cut the case.

### Deferred risk — AC-1 measures signal presence, not host reachability

AC-1 counts probed signals crossing the boundary; the delivered value is that `subspace:r` reaches the captain's real pane. Presence would be insufficient if the sandbox denied the multiplexer socket. Probed live from inside this safehouse session: `/tmp/tmux-501/default` is listable, and `TMUX=/tmp/tmux-501/default,1,0 tmux display-message -p '#{pane_id}'` returns `%19` — a real pane on the captain's tmux server, from inside the sandbox. Presence is therefore sufficient on this host. Promote to material if any host or safehouse profile denies the multiplexer socket inside the sandbox; the free check is the `tmux display-message` line already folded into the AC-1 protocol below.

## Stage Report: validation

- DONE: Judge the surface honestly: production is +26 against a declared +22, but tests are +109 against a declared +23. Say whether each test case earns its lines or whether the table is padded, and name any you would cut.
  Measured `f100e6be0..b5cc22227`: production +26 net (28/2) against a declared +16 net (22/6); tests +109 net (145/36) against a declared +29 net (56/27); whole change +135 against the declared +45 — 3.0x, with no tolerance declared and no AC narrowed. Judged case by case with 11 mutations: not padded, one exception. Each of these kills a mutation nothing else kills — "set-but-empty name is still named" (empty-treated-as-absent), "all nine set, in probe order" (any name dropped from or reordered in the list), the smoke "Zellij pair still forwards" leg (Finding 1's sentinel gate), and the +8 test-hygiene lines (reverting them reds 3 `internal/safehouse` and 11 `internal/cli` tests under this session's real `TERM_PROGRAM=tmux`). I would cut "a name outside the list is never named" — 7 lines, Finding 3.
- DONE: AC-1 cannot be machine-proven from inside the sandbox, so it rests on a captain-run comparison. State plainly what is and is not established without it, and do not treat the unit tests as a substitute for the measuring AC.
  **AC-1 is NOT VERIFIED.** The blocker is independently re-confirmed: `APP_SANDBOX_CONTAINER_ID=agent-safehouse`, `safehouse` absent from PATH, `~/.local/bin/safehouse` → "Operation not permitted", no other copy on disk. Established without the captain: (a) composer→argv, exercised with the **real compiled production code** under real process env (a throwaway `main` calling `safehouse.Wrap`) — `--env-pass=TMUX,TMUX_PANE,TERM_PROGRAM` for a tmux parent, `--env-pass=ZELLIJ_SESSION_NAME,ZELLIJ_PANE_ID` for a Zellij parent, and no built-in flag at all for a bare parent; (b) argv→child-env through a **fake** scrubbing wrapper only (the smoke test); (c) child-env→working host, live inside this sandbox (deferred risk above). NOT established: that the real `safehouse` binary honors `--env-pass` for the seven newly named variables. That one link is AC-1's entire residual, and no unit or smoke test touches it — they exercise a wrapper the implementation itself wrote. The failing baseline reproduces here: this session runs under tmux (`TERM_PROGRAM=tmux` crosses via safehouse's default `TERM*` allowance) yet `TMUX`/`TMUX_PANE` are absent inside.
- DONE: Attack the composer directly: confirm an unset name never reaches --env-pass, an empty parent yields no allowance, and dropping bare ZELLIJ breaks no live consumer.
  All three confirmed by running behavior, not by reading it. Unset names: `env -i HERDR_ENV=1` yields `--env-pass=HERDR_ENV` alone; removing the presence filter reds 4 composer cases, both new smoke legs, `TestWrapAddsTerminalTargetingEnvArgument` and `TestWrapWithExtra`. Empty parent: `env -i` emits no built-in flag; returning `[]string{"--env-pass="}` instead of nil reds "empty parent yields no allowance" and `TestWrapWithExtra`. Bare `ZELLIJ`: `env -i ZELLIJ=0 ZELLIJ_SESSION_NAME=… ZELLIJ_PANE_ID=…` names only the two probed signals even though `ZELLIJ` is set; the live consumer `plugins/subspace/skills/r/scripts/review-zellij:51-54` reads only `ZELLIJ_SESSION_NAME`/`ZELLIJ_PANE_ID`; the only bare-`ZELLIJ` hits anywhere in spacedock-subspace are a `SMOKE_ZELLIJ` shell variable in a test harness; and `zellij action list-panes` returns byte-identical output with and without `ZELLIJ` set. Caveat recorded as Finding 2: no test pins its absence.

### Recommendation: PASSED, contingent on the captain's AC-1 comparison at this gate

AC-2 is verified and falsifiable. AC-1 is not, and must not be signed off on the unit evidence; if the captain declines the run or the comparison fails a leg, this is REJECTED. `go test ./...` and `go test ./... -race` are clean except `TestCodexResolveManifestAgainstInstalledHost`, which fails identically on `main` at the merge-base `f100e6be0` (local codex plugin-cache state). Test-plan item 5 holds: the full suite gives an identical result with all nine signals forged as with none. `gofmt -l` and `go vet` clean. Three Polish findings and one deferred risk above; none blocks.

**Captain's AC-1 protocol** — outside any wrapped session, in a tmux pane, then repeat under Zellij. Set `P` first, then compare the two legs:

    P='^(ZELLIJ_SESSION_NAME|ZELLIJ_PANE_ID|TMUX|TMUX_PANE|HERDR_ENV|HERDR_PANE_ID|CMUX_WORKSPACE_ID|CMUX_SURFACE_ID|TERM_PROGRAM)='
    env | grep -E "$P" | sort                                                              # outside leg
    safehouse --trust-workdir-config --env-pass=TMUX,TMUX_PANE,TERM_PROGRAM -- /bin/sh -c "env | grep -E '$P' | sort; tmux display-message -p '#{pane_id}'"

Under Zellij the allowance is `--env-pass=ZELLIJ_SESSION_NAME,ZELLIJ_PANE_ID`, plus `,TERM_PROGRAM` whenever the outside leg shows `TERM_PROGRAM` set, and the last command becomes `zellij action list-panes >/dev/null && echo zellij-ok`. The rule for composing the string is exactly the outside leg's own output: name every one of the nine the outside leg printed, in the order listed in `P`. That is what the built binary composes for that parent, per evidence (a) — under a tmux pane it comes out as the `TMUX,TMUX_PANE,TERM_PROGRAM` written above, since tmux always sets `TERM_PROGRAM=tmux`. PASS = identical present/absent per signal on both legs, plus a real pane id from inside the wrapped session.
