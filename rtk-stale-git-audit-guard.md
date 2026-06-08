---
id: m1whe67x7nnh1gjwrnsksgak
title: Detached audits / validation must catch rtk-stale-git via an un-proxied SHA verify
status: ideation
source: "Commander (2026-06-08) qa validation: the validation ensign caught an rtk-stale-git discrepancy on its own — the worktree contract.go blob ≠ the expected e050996d blob — and only a raw, un-proxied blob-SHA compare exposed it, forcing a re-pin to 32ceb73e before the adversarial audit re-ran (verdict held PASSED). FO boot (this session) independently hit rtk mangling/blocking ls/find/ps. A SHA-pinned audit that trusts proxied git is a silent hole in exactly the discipline meant to catch silent holes."
started: 2026-06-08T18:17:44Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0199-pre-flip-mechanics
group: dev-quality
sprint-readiness: ready
---

A SHA-pinned detached adversarial audit (or validation) must not be fooled by a git proxy that serves stale state. The defense the qa validator improvised — verify the resolved commit/blob SHA via an un-proxied git read before trusting a pin — should be a standard, checkable guard, not a one-off save.

## Problem

A token-saving git proxy (rtk) can return transformed or stale git output. The detached-adversarial-audit discipline pins to a specific merge-result SHA and asserts that the deliverable's tests would catch a broken edit. If the audit resolves that SHA (or reads blobs) through the proxy and the proxy disagrees with the real object, the audit pins to the WRONG tree and validates something other than the merge result — and reports green. That is a silent correctness hole inside the discipline whose entire job is to catch silent holes. It already happened once (qa, contract.go — the contract gate) and was caught only because the validator fell back to a raw blob-SHA compare.

## How the proxy actually diverges (spike result)

The riskiest unknown was "can a proxied read deterministically disagree with the real git object?" Exercised first, before designing the guard. Findings (rtk 0.40.0, this machine; commands and outputs reproduced in a throwaway repo):

- **The proxy is the Claude Code Bash `PreToolUse` hook, not a PATH shim.** `git` on PATH is real `/usr/bin/git`. rtk only rewrites the *command string a model emits as a Bash tool call* — `rtk hook check` confirms `git show|diff|status|log` → `rtk git …` (proxied), while `git rev-parse`, `git cat-file`, `git ls-tree`, `git rev-list`, `git hash-object` are left RAW. The rewrite is a bare-`git <subcommand>` string match, so a wrapped form (`bash -c 'git show …'`) or any binary subcommand that runs git internally also escapes it.
- **The proxied read is lossy and not byte-reconstructable.** `rtk git show <commit>` returns a compacted commit summary + truncated diff: 8839 bytes → 2415 bytes, 414 lines → 109 lines on a 400-line file. A blob hashed from that output is a DIFFERENT SHA than the real object. (`git show HEAD:path` blob-ref reads happen to pass through verbatim today, so the divergence lives in commit/diff/status reads — exactly the reads an auditor uses to inspect "the merge result.")
- **Divergence is deterministic and demonstrable end-to-end.** With a stand-in stale proxy (a fake `git show` returning prior content), the proxied-derived blob SHA was `3643b652…` while the real un-proxied object SHA was `a1ba12ce…`. A SHA-pinned audit trusting the proxy would have validated the wrong blob and stayed green. An un-proxied `git cat-file blob | git hash-object --stdin` recovered the true `a1ba12ce…`, and when the working object was then mutated to `a6da3b01…` the un-proxied compare against the pin went red — the qa save, mechanized.

Consequence for the deliverable form: a Go subcommand calling `exec.Command("git", "cat-file"/"rev-parse", …)` is **structurally** un-proxied — the hook never sees a subprocess exec, and the plumbing verbs are not on the rewrite list even if it did. The guard is that subcommand; the audit/validation discipline text points at it instead of hand-rolling a raw compare.

## Proposed approach

A checkable guard the audit + validation flow invokes — `spacedock audit verify-pin` (a cobra subcommand under the existing `internal/cli` tree, using the established `runGitCmd`/`exec.Command("git", …)` idiom):

- **Inputs:** a target ref/worktree and the expected pinned SHA (commit and/or `path` for a blob pin).
- **Behavior:** resolves the current object SHA via an explicitly un-proxied git invocation (`git rev-parse <commit>` / `git rev-parse <commit>:<path>` or `git cat-file blob … | hash-object`), compares it to the expected pin, and exits non-zero with a LOUD message naming both SHAs on any mismatch. On match it exits 0 and prints the confirmed SHA. Because the resolution is an `exec.Command("git", …)` plumbing read inside the binary, it cannot be rewritten to `rtk git …` by the hook — the proxy is bypassed by construction, not by convention.
- **Discipline wiring:** the `## Detached adversarial audit` block in `docs/dev/README.md` and the validation stage text gain a concrete step — "resolve and confirm the merge-result pin with `spacedock audit verify-pin` before trusting any proxied `git show`/`diff`/`status` reading of it" — pointing at the guard rather than restating a prose discipline. (The prose edit is authoring work, not an AC on its own; the AC is the guard's behavior under test.)

## Out of scope

- rtk itself (the operator's global tooling) — this is the spacedock-side defense, not an rtk fix.
- General command-output mangling by rtk (`ls`/`find`/`ps`) — annoyance, not a correctness hole; separate if ever worth it.
- Auto-detecting "is rtk active right now" — the guard is unconditionally correct because it always reads the real object un-proxied; it does not need to know whether a proxy is in play.

## Acceptance criteria

- **The pin-verify guard ships as a `spacedock` subcommand that resolves the object SHA via an un-proxied git read and exits non-zero on a pin mismatch.** Verified by a Go test (`internal/cli`) that builds a throwaway repo, pins the real blob SHA from `git cat-file blob … | git hash-object --stdin`, then (a) runs the guard against the matching pin and asserts exit 0 + the confirmed SHA, and (b) mutates the object so it diverges and asserts the guard exits non-zero naming both SHAs. The expected SHA comes from the real git object the test independently created — an independent source that diverges from the pin in case (b) — so the check is not a tautology.
- **A proxy-shaped read that disagrees with the real object cannot make the guard pass.** Verified by a Go test that feeds the guard a divergent/transformed object value (the stale-proxy stand-in from the spike) as the "observed" side and asserts the guard rejects it, where a guard that trusted the proxied value would have passed. The pass/fail pivot is the real un-proxied object SHA, not any string in an instruction file.
- **The audit/validation discipline points at the guard.** This is the prose edit in `docs/dev/README.md`; it is authoring work that rides along, NOT a standalone acceptance criterion — per proof-policy, a substring match over the README proves nothing. The behavioral proof is the two guard tests above.

## Test plan

- **What verifies it:** Go unit tests in `internal/cli` driving the real `spacedock audit verify-pin` subcommand against throwaway git repos (real `git`, no mocks). Match case → exit 0; divergent case → exit non-zero with both SHAs. A negative test feeds a proxy-shaped divergent value and asserts the guard still goes red.
- **Cost/complexity:** Low. One small subcommand (~40-60 lines) reusing `runGitCmd`; 2-3 table-style unit tests building fixture repos with `git init`/`add`/`commit`. No live workflow run needed — the behavior is fully exercised at the unit level against real git objects.
- **Fixture vs CLI vs live:** CLI/unit (drive the binary subcommand against real on-disk repos). No live host run required.

### Spike record

Spike performed (this ideation), result above under "How the proxy actually diverges": the proxied-vs-un-proxied blob-SHA divergence was reproduced deterministically (`3643b652…` proxied-derived vs `a1ba12ce…` real object), and a prototype un-proxied compare both confirmed the matching pin and went red on a mutated object. The guard relies only on already-proven mechanisms: `exec.Command("git", …)` plumbing reads (used throughout `internal/status`, `internal/cli`) and the fact that the rtk hook rewrites only bare-`git <porcelain>` strings, never a subprocess exec.

## Stage Report: ideation

- DONE: Exercise the riskiest mechanism FIRST: reproduce an rtk proxied-vs-un-proxied git blob-SHA divergence deterministically — the staleness must be reproducible to test the guard.
  Reproduced with `rtk hook check` + a throwaway repo: hook rewrites `git show/diff/status/log`→`rtk git …` (lossy: 8839→2415 bytes, 414→109 lines on commit-show) while leaving `rev-parse`/`cat-file`/`ls-tree` raw; stale-proxy stand-in yielded blob `3643b652…` vs real object `a1ba12ce…`. Recorded under "How the proxy actually diverges (spike result)".
- DONE: Produce AC + test plan: the guard catches the divergence and refuses the stale pin, where the expected value comes from the un-proxied git object (an independent source that can diverge from the proxy) — proven by a test that goes red on a stale/divergent blob.
  Two behavioral ACs (guard exits non-zero on pin mismatch; proxy-shaped divergent value cannot make it pass) plus a CLI/unit test plan; expected SHA comes from the test's own real `git cat-file`/`hash-object` object, not from prose. Prose README edit explicitly marked NOT a standalone AC.
- DONE: Decide the deliverable form: a checkable guard/helper (not pure prose — proof-policy bars a prose-only discipline rule as the AC).
  `spacedock audit verify-pin` cobra subcommand under `internal/cli`, reusing the `runGitCmd`/`exec.Command("git", …)` idiom; structurally un-proxied because the rtk hook never sees a subprocess exec and the plumbing verbs are not on its rewrite list.

### Summary

The rtk "proxy" is the Claude Code Bash PreToolUse hook (`rtk hook claude`), not a PATH shim: it rewrites only bare-`git <porcelain>` strings a model emits, so a Go subcommand's `exec.Command("git", "cat-file"/"rev-parse", …)` is un-proxied by construction. The spike reproduced a deterministic proxied-vs-real blob-SHA divergence and a prototype guard that confirms a matching pin and goes red on a mutated object. Deliverable is a `spacedock audit verify-pin` subcommand the audit/validation discipline calls; ACs are proven by Go unit tests against real throwaway repos with the expected SHA sourced from the real git object, not from any instruction file.
