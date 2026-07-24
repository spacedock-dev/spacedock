# Cycle 23 counterexample evidence

This package retains the surviving structured host artifacts from the cycle-23
implementation checkpoint `ce4365053158ef80c1a4dc34c63256fd83da62d0`.
No artifact is a published review round.

## Results

- Claude `gate-guardrail`: PASS in 108.49s.
- Claude `recorded-gate-lifecycle`: PASS in 261.48s.
- Codex `gate-guardrail`: PASS in 99.83s.
- Codex `recorded-gate-lifecycle`: PASS in 219.01s.
- Pi round 1: FAILED before workflow action because the copied
  `openai-codex` refresh token was already used and no API-key fallback existed.
- Pi round 2: FAILED the shipped lifecycle after 84.98s. The root skipped the
  bind commit before decision and did not produce a valid successor worker
  transport. This exposed an oracle flaw that could attribute the later close
  commit to the bind barrier.
- Pi round 3: FAILED the unchanged successor-dispatch grade after 271.08s.
  The corrected bind barrier passed and a fresh `delegate` subagent exited 0,
  used `openrouter/openai/gpt-4.1-mini`, and reported a path-scoped commit. The root
  nevertheless crossed the same reported successor-dispatch boundary and also
  persisted `BEGIN_CONN`/`END_CONN` around the opaque directive token. Resolving
  the remaining discrepancy requires the prohibited harness/host-lifecycle
  expansion or another unchanged retry, so the declared hard reset fired.

Pi rounds 2 and 3 selected the existing compatible read-only package at
`/Users/clkao/.local/share/fnm/node-versions/v24.13.1/installation/lib/node_modules/pi-subagents`;
no global package was changed. Round 2 used `openai/gpt-4.1`; round 3 used
`openai/gpt-5.3-codex`.

## Retention failure

The approved host retention checklist required the wrapper command log, entity
before/after snapshots, and state Git history for every pass or failure. The Pi
fixture repositories were temporary and were cleaned when each test returned.
Neither this package nor the surviving
`/tmp/spacedock-cycle23-gate-lifecycle-evidence` tree contains those three Pi
artifact classes. They cannot be recreated from the retained sessions, process
output, or successor-subagent files, and this package must not imply otherwise.
The Pi evidence is therefore incomplete: it preserves model-visible actions and
the grader result, but not the independent durable-state history required to
reconstruct that grade.

## Inventory

Claude and Codex directories retain each root structured stream plus the final
message; Codex also retains the process result. Pi retains every root session
and surviving process output. Round 3 additionally retains successor subagent
metadata, output, and child session. The inventory does not include a Pi wrapper
command log, entity before/after pair, or state Git history.

| Artifact | SHA-256 |
| --- | --- |
| `evidence/claude/gate-guardrail-stream.jsonl` | `ceb616dec93df3eb70b73badff5ee974bfb7459005e79dc5d9aa3332a53ddc62` |
| `evidence/claude/recorded-lifecycle-stream.jsonl` | `6baf6465799ced8f7b95081abb5b543e29702c73f0764a4a3d3f740da8e79d49` |
| `evidence/codex/gate-guardrail-exec.jsonl` | `ab296f1251fa75c4d32d363a63311b7ac9b090312df10fe88ebbb74cca0c459c` |
| `evidence/codex/recorded-gate-lifecycle-exec.jsonl` | `75b92d83baffe56b732bb35649fe0e49f961529fe546ef5ade4f5ac9d8de2505` |
| `evidence/pi-round-1/root-session.jsonl` | `ffc29d704e0194c8a93722dff6dc7c80d22228fac6baa684ffd131f30bf000eb` |
| `evidence/pi-round-1/stderr.txt` | `73f2aee7737b9edee47c8ae1b31e4c90c6c79f3d6e9c5e683a5b1edfc46b1640` |
| `evidence/pi-round-2/root-session.jsonl` | `61ee8f0e67d03585f81b0bd48ba726e9045bda849a8fd30f56b143de0eddabeb` |
| `evidence/pi-round-2/stdout.txt` | `6033f8c1047cdd8d854b683ed874ff6c2d9691a2fee0046d6797c9bc2df7a9d1` |
| `evidence/pi-round-3/root-session.jsonl` | `3f79a6ad17630d84b20c63efd543646682230ee6c89192f794e8594fefa944f2` |
| `evidence/pi-round-3/stdout.txt` | `eff1c1a83c6cc37098c5a497c46343bba037a0540871018acaeaa6faee62b241` |
| `evidence/pi-round-3/successor-meta.json` | `8a07cc0939c5372a6e71e9dfbaf4d1017cc7ab3d86ea96858b8dd3104d94f53f` |
| `evidence/pi-round-3/successor-output.md` | `b1cfcd9fe78f3fecac4a33546975e3fce38428b3fccf70c1e4285a3e11e21107` |
| `evidence/pi-round-3/successor-session.jsonl` | `92bd5fd3cf1e294566a816bd73b40b0e1602863875b467cf360140513f91c728` |
