---
id: ekw79nn8z9829d77dw7y9353
title: "Pi session identity via PI_SUBAGENT_PARENT_SESSION — runtimehost identity column + install-gate sentinel key"
status: ideation
source: "Live env evidence, 2026-07-31: both the captain's shell and the FO's own root-session tool shell carried PI_SUBAGENT_PARENT_SESSION equal to the running pi session's own id (019fb5d1-85af-73f6-bb07-20bfc04004db). The runtimehost marker table (internal/runtimehost/runtimehost.go:23-24) claims pi exposes no identity env var — the code is stale about pi's actual env surface."
started: 2026-10-03T04:04:53Z
completed:
verdict:
score:
worktree:
issue:
sprint: pi-ux
group: tooling
sprint-readiness: ready
gates:
    version: 1
    records:
        - id: gate:ekw79nn8z9829d77dw7y9353:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:ekw79nn8z9829d77dw7y9353-backlog-1
              briefing:
                id: briefing:ekw79nn8z9829d77dw7y9353:backlog:attempt-1:revision-1
                digest: sha256:6c48049df2276ebd29429675eec4a44c1c35107ff9a5a198fa0def74cb6d0e7a
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:ekw79nn8z9829d77dw7y9353:backlog:1
                briefing: briefing:ekw79nn8z9829d77dw7y9353:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-10-03T04:04:17.666431Z"
                decision: approve
                reason: Covers scoping the Pi install-gate sentinel to the real session identity; live env evidence contradicts the current marker table.
                conn:
                    quote: i already said dispatch to ideation, but don't present the ideation gate until staff review finishes
                    source: captain instruction, this session, 2026-10-03
              application:
                target-stage: ideation
                state: consumed
        - id: gate:ekw79nn8z9829d77dw7y9353:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:ekw79nn8z9829d77dw7y9353-ideation-1
              briefing:
                id: briefing:ekw79nn8z9829d77dw7y9353:ideation:attempt-1:revision-1
                digest: sha256:bddefa5d368fe3663efb7ab60add2e51f80670e99294b164eb261bd194dd74b0
                room-ref: '@review/ideation/briefing-1'
---

## Problem

A failed first-officer install currently leaves `${TMPDIR:-/tmp}/spacedock-install-attempted`. Every later session sharing that temp directory sees the same sentinel and loses its install offer, even in a different project. The desired value is an install retry opportunity in a new Pi session without permitting repeated attempts in the same session.

**Baseline correction and proposed restatement (captain approval required).** This design targets main `cdfa462d1`, not a predecessor's future implementation. `internal/runtimehost/runtimehost.go` has no identity column: `Detect` returns `(host, markers, ambiguous)`. The real deferred reference is `skills/first-officer/references/fo-install.md`, not `fo-install-gate.md`; there is no cwd-hash fallback. The title/frontmatter preserve the original hypothesis for provenance, not a verified implementation requirement. The original AC-1 mapping and AC-3 root-env assertion are superseded by the proposed observable criteria below.

The seed's root identity assumption is false in the measured environment. The first officer's root shell reported `PI_CODING_AGENT=true` and no `PI_SUBAGENT_PARENT_SESSION`; the dispatched worker reported `PI_SUBAGENT_PARENT_SESSION=01a0ffab-86e4-762a-996b-1f13eec90317`, the **parent's** session ID. Supervisor evidence identifies Pi 1.0.0, pi-subagents 0.75.0, pi-intercom 0.16.0. Installed `pi-subagents/docs/watchdog.md:185` states that detached runners retain the exact launch parent identity, while root and in-process foreground hosts deliberately publish no global parent identity because multiple sessions can share a host. Child env is not evidence of the child's own or a root's identity.

## Proposed approach

Use Pi's current-session API in the **existing** Spacedock extension, which is loaded independently of the launcher binary. Add one read-only tool, `spacedock_install_scope`, with no arguments. At each invocation, read `ctx.sessionManager.getSessionId()`; for a nonempty string return structured `sessionId` and `sentinelPath` plus matching JSON text. The path is `${process.env.TMPDIR || '/tmp'}/spacedock-install-attempted-pi-${sha256(sessionId)}` (join as a path; lowercase full SHA-256). Read the context on every invocation, never cache session identity at startup and never publish it in `process.env`. Missing/invalid identity yields an explicit tool error and no path. Hashing accepts Pi's custom session IDs safely without putting arbitrary IDs into filenames.

Bind `«install.scope»` in `pi-first-officer-runtime.md` to that tool. On the Pi branch of the existing deferred install offer, resolve the scope **before** inspecting a sentinel. Use the returned path for the existing existence check, pre-install touch, and manual-removal hint. A legacy global sentinel must not suppress a supported Pi session's offer. Keep touch-before-install and the single post-install version re-check unchanged. Query again on a later invocation, so `/new` or a session switch does not reuse the prior scope. Same Pi session resumed retains the same key; a new/forked session with a different ID gets a new key. Identity deliberately does not include cwd.

If the extension/tool/API is unavailable or errors, Pi must print the existing channel-correct manual command and abort the automatic offer with the explicit reason that session identity could not be obtained. Do not silently substitute the global sentinel, parent ID, newest file, project hash, PID, or a fresh random token. This is a bounded unsupported lane, not a claim that all Pi setups now have session-safe offers. Non-Pi behavior remains unchanged. The tool itself neither installs nor touches a sentinel; captain approval and the FO's existing shell actions remain the authority boundary.

### Necessity and alternatives

All added mechanisms serve value AC-2:

- **Current-context tool:** the smallest root-observable real identity bridge available before a working binary exists. The simpler env-table change fails in root shells and misidentifies children. Extending the existing extension needs no new package, broker, CLI verb, global environment mutation, or bootstrap injection.
- **Session artifact alternative:** Pi documents JSONL session headers and `/session`; an explicitly known active file could supply its ID. But default agent/session directories contain concurrent sessions and allow `--session-dir`, resume, custom IDs, and `--no-session`. Picking newest/by-cwd cannot identify the caller. Asking the captain to copy `/session` output for every failed install adds a manual prerequisite and still does not solve automatic root observation. The current-context API already owns this association.
- **Other root-observable values:** random per-invocation tokens break same-session retry suppression; a retained token introduces persistence/compaction ownership solely to recreate an existing ID. PID and agent-directory values do not rotate on `/new`. `intercom({action:'status'})` returns a broker identity that can be overridden by configured or claimed stable IDs (`pi-intercom/index.ts` session-start registration), so it is not a reliable Pi-session identity contract. No new dependency on intercom is needed.
- **Safe path encoding:** full SHA-256 using Node's built-in crypto, rather than interpolating arbitrary custom IDs into a pathname. This buys a bounded filename, not secrecy or a new identity registry.
- **Unsupported fallback:** hint-and-abort when identity cannot be obtained avoids pretending session isolation exists. Making all Pi automatic offers unsupported is simpler but forfeits the value in the measured, already supported extension lane; retain it only for missing capability. A cwd/global fallback preserves the original false-negative and cannot satisfy AC-2.

## Out of scope

No changes to Pi, pi-subagents, pi-intercom, runtimehost's detection/signature/marker order, Claude/Codex identity or install behavior, sandbox policy, approval semantics, installers, channel selection, launch grammar, state formats, extension packaging, or FO bootstrap/compaction behavior. No registry, session-file discovery, automatic sentinel cleanup/migration, concurrency lock, or broad one-attempt redesign. Same-session simultaneous FO offers remain outside this existing sequential-offer contract. Do not depend on `fo-boot-install-hint-linux-direct-sandbox` landing first.

## Expected surface

Estimate net LOC change: +170, across 6 files. Estimated insertions: 190; deletions: 20. Tolerance: net +90 to +250 and at most 7 files; anything outside this or the semantic boundaries requires a revised gate.

1. `.pi/extensions/spacedock.ts`: small read-only tool registration, current-context ID and built-in hash/path handling.
2. `test/spacedock.test.ts`: existing fake gains registerTool support; scope isolation/error/lifecycle tests. Never put tests under `.pi/extensions/`.
3. `skills/first-officer/references/fo-install.md`: only Pi sentinel selection/failure lane and selected-path hints.
4. `skills/first-officer/references/pi-first-officer-runtime.md`: the install-scope tool binding.
5. `internal/contractlint/version_gate_smoke_test.go`: smoke guards for the new instruction binding and preservation of install sequencing; text claims only.
6. `docs/site/get-started/install.md`: Pi failure/retry and capability-limit documentation.

Observable changes: one Pi read-only tool surface; Pi's temporary install sentinel filename; new-session offer eligibility; explicit manual-only behavior when Pi identity is unavailable. No new CLI grammar or persisted workflow format. Existing install approval and mutation authority remain unchanged. This ideation edits only this state entity; code changes await captain approval.

## Acceptance criteria

These criteria are a **proposed restatement**, not approval to implement a root-env identity claim.

**AC-1 — The install scope identifies the currently executing Pi session, not its parent or process.** With the parent env absent or conflicting, the tool's `sessionId` equals the caller's Pi context ID and the path encodes that ID; `/new`/switch to B changes it and returning/resuming A restores A's path. Custom IDs cannot escape the temp directory. Verified by `test/spacedock.test.ts` tests `install scope follows current session context` and `install scope bounds custom IDs`, plus the root RPC spike/live recheck against `get_state.sessionId`. Named falsifying edits: **use-parent-env** (replace context ID with `PI_SUBAGENT_PARENT_SESSION`); **cache-first-session** (reuse the first ID); **raw-id-path** (remove hashing).

**AC-2 — After a failed install in Pi session A, session B in the same project and temp directory still receives an install offer, while A cannot retry automatically.** Sentinel keys must derive from each session's actual context identity; a legacy global sentinel must not suppress B. Independent baseline: current main's global sentinel yields B offer count 0 after A fails; the proposed lane yields B offer count 1, A retry count 0, and distinct on-disk sentinel paths. Verified by `install scope isolates failed attempts` executing the registered tool for A/B/A with real temp files and a failing installer stub, and a bounded FO smoke that observes the actual offer/approval/sentinel-touch behavior using the shipped instructions. Named falsifying edit: **restore-global-sentinel** (make the Pi gate check the legacy global path); this must fail the FO smoke even if tool-only tests remain green.

**AC-3 — Pi cannot silently fall back to a coarser or fabricated identity when session scope is unavailable.** Missing tool/API, empty ID, or tool error results in manual command plus explicit identity-unavailable reason and no automatic install; no sentinel is created. Verified by `install scope rejects missing identity`, instruction smoke, and the FO failure-lane trace. Named falsifying edits: **fallback-global-on-error** (permit the legacy offer path); **invent-id-on-error** (generate a random key).

**AC-4 — The existing install safety and non-Pi contracts are preserved.** Approval is still required, sentinel creation still precedes the installer, decline/sandbox cause no installation, and failure never loops; non-Pi retains its existing sentinel. `runtimehost.Detect` still uses the same host markers, ordering and ambiguity rule with no new parent-env inference. Verified by existing Go owners, new instruction smoke and FO trace. Named falsifying edits: **touch-after-install** (move touch after a failing installer); **parent-is-host-marker** (add parent env to detection); **replace-non-pi-sentinel** (apply Pi policy globally).

## Test plan

### Primary proof owners and implementation order

- **Runtimehost marker table:** existing `internal/runtimehost/runtimehost_test.go::TestDetectMarkerMatrix` owns host/marker ordering and ambiguity. No identity-column tests or production edits are planned because no identity column exists and the value does not require one. If a negative parent-only regression case is needed, keep it in that owner, not a new test suite (the seventh-file tolerance). This exercises **parent-is-host-marker**.
- **Pi identity/path behavior:** existing `test/spacedock.test.ts` is the extension behavior owner; extend its fake with tool registration and call the real registered execute function with context A, B, A and isolated TMPDIR. Assert known independent SHA-256 expectations, not an expectation computed by the production helper. Parent env deliberately conflicts. Include missing method, thrown API, empty ID, custom ID, same-process session change, and no writes during the read-only tool call. Add focused tests before implementing the tool; run `bun test test/spacedock.test.ts`.
- **Install sentinel:** there is currently no behavioral sentinel owner: `internal/contractlint/version_gate_smoke_test.go` covers deferred-load/sandbox **text**, and `install_hint_drift_test.go` covers command wording. The new `install scope isolates failed attempts` in the extension test is primary deterministic path/on-disk proof: seed the legacy sentinel, claim A's returned path, run a stub exiting 1, observe A suppressed and B eligible. This proves path isolation, not that an LLM obeyed the instructions. The actual FO offer policy is separately owned by a bounded live/manual trace below; do not misreport a duplicated test-side gate as end-to-end proof.
- **Instruction wiring/safety:** extend the existing contractlint smoke before editing skills. Assert the Pi runtime binding, required unavailable behavior, selected-path hints, touch-before-install ordering, and unchanged non-Pi branch. These falsify stale text/wiring independently of tool hashing; they are not proof of runtime offers. Run `go test ./internal/contractlint ./internal/runtimehost`.
- **FO behavior trace (required before validation acceptance):** in disposable directories, use the installed Spacedock extension and real Pi IDs in two root sessions with missing/incompatible binary and an explicitly approved harmless installer stub returning 1. No network/package install. Keep TMPDIR/cwd fixed and preseed the old global sentinel. Record session ID, displayed offer, approval, touch-before-stub event, exit 1, sentinel path/existence, A's repeated startup suppression, and B's new offer. Repeat the missing-tool/API lane to observe manual-only output and no stub invocation; decline/sandbox remain no-effect checks. This is the proof that kills **restore-global-sentinel** in the actual gate. Preserve redacted trace and temp-file assertions as validation evidence; no success claim from prose search alone.

Estimated cost: small TS unit/mocked-context addition and Go text smoke, seconds after tooling warm-up; isolated no-model RPC probe under 10 seconds; bounded FO smoke a few model turns per lane. Live credentials may be copied into an isolated Pi home only when needed; never mutate global configuration or actually install software. Run repo-required Go suite/race/format checks during implementation. No new test framework or production dependency is required.

### Spike result and limits

**Spike executed, not “no spike needed.”** Root Pi 1.0.0 in RPC mode, isolated HOME/agent/session directories, stripped PI/SPACEDOCK env, discovery disabled, only the throwaway explicit extension, offline, no credentials or model requests. A read-only registered tool's execute callback was exercised from the root `session_start` context. Independent `get_state` RPC output supplied the expected real session IDs and files.

Observed session A `01a0fff3-4917-76f0-983e-dd2742cedc9c`, key `dde108e22df0a649ee0755674ae41980a72a163ae0f0c8ce3475f46a279ad5cb`; after `new_session`, B `01a0fff3-49ca-76f0-983e-dd29146df446`, key `46fc295063c68ae7c2b67d4cb0458083c5cd7a7e871e338f19be512f99809deb`. Both matched RPC identity/file exactly and had parent env null; repeat `get_state` retained B. A real temp-file failure simulation gave legacy B offer=false, scoped B offer=true, A retry suppressed=true.

First run failed a harness assumption of exactly two `session_start` callbacks: Pi fired three callbacks for two distinct IDs. Corrected the probe to compare distinct identities, then passed. This supports lookup at tool invocation rather than relying on startup delivery counts. This is an API/root-identity/path spike, **not** a model-issued tool/FO-install end-to-end proof, persisted transcript proof, or a compatibility matrix. TUI, older Pi versions, disabled-extension sessions and resumed/forked sessions still require the named later tests/trace; no unsupported environment is presented as proven. The current root/worker env captures above additionally demonstrate why the seed mechanism is rejected.

The reproducible spike below requires only Python 3 and Pi 1.0.0 on PATH, creates/removes all probe files under a temp directory, and makes no model calls. Run it as `python3 probe.py`. Exact absolute machine package paths are not prerequisites. Package documentation inspected for semantic context: Pi `docs/sessions.md` (`/session`, overrides and `--no-session`), `docs/session-format.md` (custom IDs), `docs/rpc-commands.md` (`get_state`), pi-subagents `docs/watchdog.md:185`, and pi-intercom `index.ts` session-start identity registration.

## Proposed instruction and documentation diff

These are exact replacement/addition targets for implementation, subject to the gate.

**`skills/first-officer/references/pi-first-officer-runtime.md`, add binding:**

> - `«install.scope»`: Call `spacedock_install_scope` with no arguments at each install-offer invocation. Its `sessionId` is this Pi context's identity and `sentinelPath` is the selected attempt marker. Missing tool, tool error, or absent path means identity unavailable: manual install hint and ABORT, never env/file guessing. Do not use `PI_SUBAGENT_PARENT_SESSION`, which names a detached child's parent, not the root/current session.

**`skills/first-officer/references/fo-install.md`, replace step 1** (currently “One attempt ever — sentinel `${TMPDIR:-/tmp}/spacedock-install-attempted`” and a literal removal hint):

> 1. **Select scope, then check the one-attempt sentinel.** On Pi, resolve `«install.scope»` through the Pi runtime binding and use its `sentinelPath`. If unavailable, print the channel-correct manual command, say “Pi session identity unavailable; automatic install is disabled”, and ABORT without creating a marker or offering an install. On other hosts, retain `${TMPDIR:-/tmp}/spacedock-install-attempted`. If the selected sentinel exists, skip the offer and hint-and-abort: print the channel-correct command, name the selected sentinel, and give a shell-quoted `rm` command for that exact path. If `SPACEDOCK_BIN` was set-but-stale, name its session-scoped replacement. Same message on the step-4 fallback.

Replace step 3's “`touch` the sentinel BEFORE the install runs” with “`touch` the selected sentinel path (shell-quoted) BEFORE the install runs”; retain the rest. Step 2's offer, step 4's convergence, step 5's no-loop policy, and the sandbox/channel sections remain unchanged.

**`docs/site/get-started/install.md`, under `## Troubleshooting`, before “Run `spacedock doctor`.”, add:**

> With the Spacedock Pi extension loaded, the first officer's automatic install offer is limited to one attempt per Pi session. A failed attempt does not suppress the offer in a new session; resuming the same session retains its attempt marker. The marker lives under `${TMPDIR:-/tmp}` as `spacedock-install-attempted-pi-<session-id-sha256>`. The first officer names the exact marker in its manual-retry hint. Older global attempt markers do not block this Pi lane.
>
> If Pi cannot supply the current session identity (for example, the extension is disabled or outdated), the first officer prints a manual install command instead of offering an automatic install. `PI_SUBAGENT_PARENT_SESSION` is not a root-session identity source. Other runtimes and sandbox restrictions are unchanged.

## Reproducible spike

```python
import hashlib, json, os, pathlib, select, shutil, subprocess, tempfile

pi = shutil.which('pi')
assert pi, 'Install Pi 1.0.0 and put pi on PATH'
with tempfile.TemporaryDirectory(prefix='spacedock-pi-identity-') as tmp:
    root = pathlib.Path(tmp)
    extension = root / 'probe.ts'
    evidence = root / 'evidence.jsonl'
    extension.write_text('''import { appendFileSync } from "node:fs";
import { createHash } from "node:crypto";
export default function(pi) {
  const tool = {
    name: "spacedock_install_scope_probe", label: "Install scope probe",
    description: "Read the current Pi session scope",
    parameters: {type: "object", properties: {}, additionalProperties: false},
    async execute(_id, _params, _signal, _update, ctx) {
      const sessionId = ctx.sessionManager.getSessionId();
      const key = createHash("sha256").update(sessionId).digest("hex");
      return {content: [{type: "text", text: key}], details: {sessionId, key}};
    }
  };
  pi.registerTool(tool);
  pi.on("session_start", async (_event, ctx) => {
    const result = await tool.execute("probe", {}, undefined, undefined, ctx);
    appendFileSync(process.env.PROBE_OUTPUT, JSON.stringify({
      ...result.details, sessionFile: ctx.sessionManager.getSessionFile(),
      parentEnv: process.env.PI_SUBAGENT_PARENT_SESSION ?? null
    }) + "\\n");
  });
}
''')
    env = {k: v for k, v in os.environ.items() if not k.startswith(('PI_', 'SPACEDOCK_'))}
    env.update(HOME=tmp, PI_CODING_AGENT_DIR=str(root/'agent'), PROBE_OUTPUT=str(evidence))
    command = [pi, '--offline', '--mode', 'rpc', '--no-extensions', '--no-skills',
               '--no-prompt-templates', '--no-context-files', '--no-themes',
               '--session-dir', str(root/'sessions'), '--extension', str(extension)]
    p = subprocess.Popen(command, cwd=tmp, env=env, stdin=subprocess.PIPE,
                         stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    def rpc(command):
        p.stdin.write(json.dumps({'type': command})+'\n'); p.stdin.flush()
        for _ in range(20):
            ready, _, _ = select.select([p.stdout], [], [], 2)
            if not ready: continue
            line = p.stdout.readline()
            if not line: raise RuntimeError(p.stderr.read())
            value = json.loads(line)
            if value.get('command') == command:
                assert value.get('success'), value
                return value.get('data')
        raise TimeoutError(command)
    try:
        first = rpc('get_state')
        rpc('new_session')
        second = rpc('get_state')
        repeated = rpc('get_state')
        callbacks = [json.loads(line) for line in evidence.read_text().splitlines()]
        rows = list({row['sessionId']: row for row in callbacks}.values())
        assert len(rows) == 2, rows
        for row, state in zip(rows, (first, second)):
            assert row['sessionId'] == state['sessionId']
            assert row['sessionFile'] == state['sessionFile']
            assert row['parentEnv'] is None
            assert row['key'] == hashlib.sha256(state['sessionId'].encode()).hexdigest()
        assert first['sessionId'] != second['sessionId'] == repeated['sessionId']
        # Failed-install stand-in: claim before a process returning exit 1.
        legacy = root/'spacedock-install-attempted'
        legacy.touch()
        sentinel_a = root/('spacedock-install-attempted-pi-'+rows[0]['key'])
        sentinel_b = root/('spacedock-install-attempted-pi-'+rows[1]['key'])
        sentinel_a.touch()
        assert subprocess.run(['sh', '-c', 'exit 1']).returncode == 1
        assert sentinel_a.exists() and not sentinel_b.exists() and legacy.exists()
        print(json.dumps({'pi_version': subprocess.check_output([pi, '--version'], text=True).strip(),
                          'sessions': rows, 'session_start_callbacks': len(callbacks), 'root_parent_env_absent': True,
                          'rpc_identity_matches': True, 'new_session_rotates': True,
                          'same_session_stable': True, 'legacy_offer_B': False,
                          'session_key_offer_B': True, 'same_session_retry_suppressed': True}, indent=2))
    finally:
        p.terminate()
        p.communicate(timeout=10)
```

## Stage Report: ideation

- DONE: Task body states the problem, the chosen approach, and criteria under the exact `## Acceptance criteria` heading, with every AC paired to a named falsifying edit.
  Main cdfa462d1 baseline corrected; four proposed observable ACs pair scope/isolation/fallback/safety with named mutations; six-file +170 net LOC estimate and concrete skill/doc replacements recorded.
- DONE: A value AC measures the end outcome - a failed install in one session does not suppress the next session's install offer, with the sentinel key derived from Pi's real session identity - against an independent baseline.
  AC-2 contrasts legacy B offer count 0 with scoped B offer count 1 and A retry count 0; root RPC context identity, not parent env, owns the key. Restatement is PROPOSED for captain approval.
- DONE: The test plan names the primary proof owner for the runtimehost marker table and the install-gate sentinel, and records the spike result or an auditable "no spike needed" with the proven environment mechanisms.
  DetectMarkerMatrix remains detection owner; existing extension tests own new path isolation, contractlint owns text, and bounded FO trace owns actual offers; exact root RPC reproducer and passing observations are in the body.
  AC-1 proof plan: `install scope follows current session context` and `install scope bounds custom IDs` in `test/spacedock.test.ts`, plus root RPC/live comparison with `get_state.sessionId`, falsify use-parent-env, cache-first-session, and raw-id-path; the recorded root RPC spike matched current IDs before/after new_session.
  AC-3 proof plan: `install scope rejects missing identity`, instruction smoke, and the FO failure-lane trace require manual command plus explicit identity-unavailable reason, no automatic install and no sentinel; fallback-global-on-error and invent-id-on-error must fail these checks.
  AC-4 proof plan: existing Go owners including `TestDetectMarkerMatrix`, new instruction smoke, and the FO trace preserve approval, touch-before-install, decline/sandbox no-effect, no-loop and non-Pi sentinel behavior; touch-after-install, parent-is-host-marker, and replace-non-pi-sentinel must fail these checks.
- DONE: Riskiest identity mechanism exercised.
  Isolated Pi 1.0.0 root callback matched get_state IDs before/after new_session with no parent env; temp sentinel simulation discriminated legacy/scoped behavior. Initial callback-count assertion failed (3 events/2 IDs), corrected to identity comparison and reran successfully.
- DONE: Baseline focused checks and artifact integrity.
  `go test ./internal/runtimehost ./internal/contractlint` passed; no production edits made. `gofmt -w ./cmd ./internal` completed; its unrelated pre-existing formatting change was restored to keep this stage entity-only.
- SKIPPED: Shipped FO offer end-to-end proof and compatibility matrix.
  Ideation proposes code/instructions, not their implementation; spike does not claim an LLM-issued tool call or actual FO offer. Test plan makes the bounded real FO trace and missing-capability lane mandatory for acceptance.
- FAILED: Full Go suite and race baseline completion within probe budgets.
  `go test ./...` exceeded 180 seconds; `go test ./... -race` exceeded 45 seconds. No full-suite success claimed; implementation must rerun with an appropriate budget.

### Summary

The seed's identity-column/file assumptions and root parent-env claim were disproved, and the task body now explicitly proposes a captain-gated restatement preserving the actual cross-session install value. The chosen minimal bridge uses the existing Pi extension's current-session context, with no runtimehost expansion or process-global identity; its root identity mechanism was exercised against independent RPC state. Missing capability stays manual-only, and exact bounds, alternatives, proof ownership, documentation changes and a reproducible spike are recorded for staff review before the captain gate.

### Report repair verification

`spacedock status --read docs/dev/.spacedock-state/pi-session-identity-subagent-parent-session.md --ac-scan --json --workflow-dir docs/dev` (exit 0; all four criteria have citations):

```json
{"command":"read","stage":"ideation","acs":[{"id":"AC-1","line":"89","unevidenced":"false","citations":[{"line":"232","text":"  AC-1 proof plan: `install scope follows current session context` and `install scope bounds custom IDs` in `test/spacedock.test.ts`, plus root RPC/live comparison with `get_state.sessionId`, falsify use-parent-env, cache-first-session, and raw-id-path; the recorded root RPC spike matched current IDs before/after new_session."}]},{"id":"AC-2","line":"91","unevidenced":"false","citations":[{"line":"229","text":"  AC-2 contrasts legacy B offer count 0 with scoped B offer count 1 and A retry count 0; root RPC context identity, not parent env, owns the key. Restatement is PROPOSED for captain approval."}]},{"id":"AC-3","line":"93","unevidenced":"false","citations":[{"line":"233","text":"  AC-3 proof plan: `install scope rejects missing identity`, instruction smoke, and the FO failure-lane trace require manual command plus explicit identity-unavailable reason, no automatic install and no sentinel; fallback-global-on-error and invent-id-on-error must fail these checks."}]},{"id":"AC-4","line":"95","unevidenced":"false","citations":[{"line":"234","text":"  AC-4 proof plan: existing Go owners including `TestDetectMarkerMatrix`, new instruction smoke, and the FO trace preserve approval, touch-before-install, decline/sandbox no-effect, no-loop and non-Pi sentinel behavior; touch-after-install, parent-is-host-marker, and replace-non-pi-sentinel must fail these checks."}]}]}
```
