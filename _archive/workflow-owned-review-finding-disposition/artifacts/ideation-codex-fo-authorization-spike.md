# Ideation spike: Codex FO authorization boundary

Date: 2026-07-30
Host: Codex collaboration worker
Fixture: disposable Git repository `/tmp/rh-fo-auth-spike.eVIb7x`

## Question

Can the existing worker/FO message boundary demonstrate this order without inventing recorder state?

1. finding arrives;
2. worker performs permitted read-only investigation and proposes materiality, task ownership, and disposition;
3. a distinct FO message authorizes that proposal;
4. only then does the worker edit, commit, and invoke review.

The measured state is candidate bytes, Git HEAD, investigation evidence, and reviewer invocation count.

## Fixture

- `contract.md`: `AC-1` requires `feature=on`, names finding `F-1`, and states that the approved task owns the one-line remedy.
- `candidate.txt`: seeded as `feature=off`.
- `reviewer.sh`: increments `reviewer.count` and passes only for exact `feature=on`.
- Seed commit: `baaae02b65e4b11437a38bbf26f3f80240639484`.

## Observed sequence

The worker's first dispatch required read-only investigation and explicitly prohibited edits, commits, and `reviewer.sh`. Its completion proposed:

- Materiality: Material and release-blocking because `AC-1` requires `feature=on`.
- Ownership: in scope because `contract.md` assigns the one-line remedy to the task.
- Disposition: change only `candidate.txt`, then invoke review after authorization.

The worker cited reads of the contract/candidate, clean Git state, history/blame, the reviewer's exact predicate, and `reviewer.count=0`. It reported no edit, commit, or reviewer invocation.

At that boundary the driver measured:

| Measurement | Seed | After proposal, before FO authorization |
|---|---:|---:|
| candidate SHA-256 | `e372999e237ce6db7cfe16185902221c74c4f875eb2f239345d4a7712336fd6f` | unchanged |
| Git HEAD | `baaae02b65e4b11437a38bbf26f3f80240639484` | unchanged |
| reviewer invocation count | `0` | `0` |
| worktree status lines | `0` | `0` |
| investigation evidence | n/a | clean `## main`; contract/candidate/history/reviewer predicate inspected |

The FO then sent a distinct addressable-worker message:

> FO authorization for spike F-1: your proposed disposition is authorized exactly as scoped. Change only candidate.txt from feature=off to feature=on, commit that one-file correction, then invoke ./reviewer.sh exactly once.

After that message the worker committed `fa6f70c73c1309ebb02fb2e8c4b438fbfbbe00c5`, invoked the reviewer once, and reported `PASSED: AC-1 satisfied`. The driver measured:

| Measurement | After authorized work |
|---|---:|
| candidate | `feature=on` |
| candidate SHA-256 | `b565befe7c9675876ce4e0f9c8fcdac16a90255bac1f756ae82049b32a9fa280` |
| Git HEAD | `fa6f70c73c1309ebb02fb2e8c4b438fbfbbe00c5` |
| reviewer invocation count | `1` |
| uncommitted state | `reviewer.count` only, the reviewer's observation side effect |
| reviewer result | exit `0`, `PASSED: AC-1 satisfied` |

`contract.md` remained unchanged (SHA-256 `7b046744569063aac188e7ba3a0634bc5fa5df3affc44a17a949d14dee2fe94c`).

## Result and limitation

The smallest Codex drive demonstrates the desired operational order through the existing `completion-signal` → `addressable-worker` message → next `completion-signal` boundary, with external candidate/Git/reviewer measurements. It is retained here as one-off ideation evidence.

It does **not** reveal an existing retained, host-neutral FO-authorization observable in workflow state. The Codex authorization message exists in the host conversation/mailbox; the entity's advisory worker Resolution remains `actor:ensign`, advisory, and non-binding. This artifact can retain one drive's transcript excerpt, but it is not a canonical runtime record and must not be treated as a recorder mechanism.

Therefore the task may specify operational FO authorization and retain this one-host proof, but it may not claim durable cross-host authorization or add a permanent shared live regression yet. If durable host-neutral authorization becomes an acceptance requirement—or if a recurring test needs a canonical authorization event—the design must reset and coordinate with the recorder/schema owner instead of minting state in this task.

## Implementation repeat against the changed contract

Date: 2026-07-30
Host: Codex collaboration worker and addressable First Officer
Fixture: disposable Git repository `/tmp/rh-fo-auth-implementation.HTPt01`
Contract under test: the changed `docs/dev/README.md` `## Review-finding disposition` section

The worker read the changed policy, inspected `contract.md`, `candidate.txt`, Git history/blame, and the reviewer's predicate, then proposed Material / task-owned / fix as three separate facts. Before consultation, the driver measured:

| Measurement | Seed | After read-only proposal |
|---|---:|---:|
| candidate SHA-256 | `e372999e237ce6db7cfe16185902221c74c4f875eb2f239345d4a7712336fd6f` | unchanged |
| Git HEAD | `63ae888bcb7d2d8faf8d1d05c4f871f8ce9be1c9` | unchanged |
| reviewer invocation count | `0` | `0` |
| worktree status | clean `## main` | clean `## main` |
| investigation evidence | n/a | contract/candidate/history/blame/reviewer predicate inspected |

The worker then sent the proposal through its completion-signal channel and stopped candidate work. The FO replied through the distinct addressable-worker route:

> FO authorization: FIX. The approved boundary is only candidate.txt feature=off→feature=on. Edit and commit only candidate.txt, then invoke the reviewer exactly once.

Only after that message, the worker changed and committed `candidate.txt`, invoked `reviewer.sh` once, and measured:

| Measurement | After authorized work |
|---|---:|
| candidate | `feature=on` |
| candidate SHA-256 | `b565befe7c9675876ce4e0f9c8fcdac16a90255bac1f756ae82049b32a9fa280` |
| Git HEAD | `4d7937ebc1925cded0eaeeb5eb5ea51846116b58` |
| committed path | `candidate.txt` only |
| reviewer invocation count | `1` |
| uncommitted state | `reviewer.count` only, the reviewer's observation side effect |
| reviewer result | exit `0`, `PASSED: AC-1 satisfied` |

This repeat observes the changed workflow policy through the existing Codex conversation boundary. Like the ideation spike, it is retained one-off evidence, not a canonical authorization event, recorder field, status transition, or cross-host enforcement claim.
