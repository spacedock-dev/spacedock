# Live evidence — repair-codex-rejection-round-recording (hz2ankag6fk379ssabpv4ckc)

Copied out of the agent job's tmp dir, which does not survive job cleanup. The
`.jsonl.gz` files are raw codex exec streams; read one with `gzcat FILE`. Every
sequence below is extracted from the artifact named above it, so each claim can
be re-derived rather than taken on trust.

Codex chains several verbs into a single shell call, so each entry lists EVERY
verb its command carries and shows the command itself. Reading a bracket as "one
verb ran here" would misattribute the raw-git escape in section 1.

`repair-loop.sh` is the loop driver; `ci-codex-shim-replica.sh` replicates CI's
codex PATH shim. A green requires ALL of: verdict PASS, an in-stream
`gate prepare` emitting `state=open`, an in-stream `state commit` reporting the
inline commit, and zero pre-fix no-op strings.

## 1. prefix-baseline-trap-reproduction.codex-exec.jsonl.gz

PRE-FIX binary, reached by accident: the loop ran without SPACEDOCK_BIN, so the
live test resolved a stale pre-fix spacedock from PATH. The run PASSED. The
stream shows why — the FO chained `git add && git commit` into the same shell
call as `gate prepare`, escaping the no-op verb by hand. That is the falsifying
baseline: a pre-fix pass is model improvisation, not a working chain. Counted
in-stream: 4 pre-fix "nothing to commit to a state checkout" strings and 2
successful raw git commits.

  exit=0  [state commit]
    $ '${SPACEDOCK_BIN:-spacedock} status --workflow-dir '"'/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow1971119381/002' --set rejection-task sta
    > Inline workflow — entities live beside the README; nothing to commit to a state checkout.
  exit=0  [state commit]
    $ '${SPACEDOCK_BIN:-spacedock} status --workflow-dir '"'/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow1971119381/002' --set rejection-task sta
    > Inline workflow — entities live beside the README; nothing to commit to a state checkout.
  exit=2  [gate record]
    $ '${SPACEDOCK_BIN:-spacedock} gate record --round validation/1 --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl --workflow-
    > Error: unknown gate flag: validation/1
  exit=0  [gate record]
    $ '${SPACEDOCK_BIN:-spacedock} gate record --help'
  exit=0  [gate record]
    $ '${SPACEDOCK_BIN:-spacedock} gate record rejection-task --round validation/1 --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.js
    > round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=4
  exit=0  [state commit]
    $ '${SPACEDOCK_BIN:-spacedock} status --workflow-dir '"'/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow1971119381/002' --set rejection-task sta
    > Inline workflow — entities live beside the README; nothing to commit to a state checkout.
  exit=0  [gate record]
    $ '${SPACEDOCK_BIN:-spacedock} gate record rejection-task --round validation/2 --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.js
    > round=round:rejection-task:validation:2 stage=validation cycle=2 briefing=briefing:rejection-task:validation:round-1 entries=4
  exit=0  [gate prepare+git commit]
    $ "git add -- rejection-task/index.md && git commit -m 'feedback: record rejection-task cycle 2 pass'; git status --short -- rejection-task/index.md; "'${SPACEDOCK_BIN:-spa
    > state=open
  exit=1  [state commit]
    $ '${SPACEDOCK_BIN:-spacedock} state commit rejection-task --workflow-dir '"'/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow1971119381/002'; "'
    > Inline workflow — entities live beside the README; nothing to commit to a state checkout.
    > Error: no ## Acceptance criteria section in this file
  exit=0  [git commit]
    $ "git add -- rejection-task/index.md rejection-task/review && git commit -m 'gate: bind rejection-task validation review'; git status --short"

## 2. mechanism-chain-prepare-commit-prepare.codex-exec.jsonl.gz

CI-config. This layer's proof, and no raw git anywhere in it: `gate prepare`
REFUSED the uncommitted artifact, `state commit` made it durable, the retried
`gate prepare` emitted state=open. On the pre-fix binary that middle step is an
exit-0 no-op and the retry fails identically.

  exit=2  [gate record]
    $ '${SPACEDOCK_BIN:-spacedock} gate record --round validation/1 --briefing /var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow685325272/002/rejecti
    > Error: unknown gate flag: validation/1
  exit=0  [gate record]
    $ '${SPACEDOCK_BIN:-spacedock} gate record rejection-task --round validation/1 --briefing /var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow685325
    > round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=4
  exit=0  [gate record]
    $ '${SPACEDOCK_BIN:-spacedock} gate record rejection-task --round validation/2 --briefing /var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow685325
    > round=round:rejection-task:validation:2 stage=validation cycle=2 briefing=briefing:rejection-task:validation:round-1 entries=4
  exit=1  [gate prepare]
    $ '${SPACEDOCK_BIN:-spacedock} gate prepare rejection-task --question '"'Should the corrected rejection-task validation be approved after cycle 2 passed?' --artifact /var/f
    > Error: --artifact /var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow685325272/002/rejection-task/index.md: selected source d
  exit=0  [state commit]
    $ '${SPACEDOCK_BIN:-spacedock} state commit rejection-task --workflow-dir /var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow685325272/002'
    > Committed rejection-task in the inline workflow repository; nothing pushed.
  exit=0  [gate prepare]
    $ '${SPACEDOCK_BIN:-spacedock} gate prepare rejection-task --question '"'Should the corrected rejection-task validation be approved after cycle 2 passed?' --artifact /var/f
    > state=open
  exit=0  [state commit]
    $ '${SPACEDOCK_BIN:-spacedock} state commit rejection-task --workflow-dir /var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow685325272/002'
    > Committed rejection-task in the inline workflow repository; nothing pushed.

## 3. composed-tree-green.codex-exec.jsonl.gz

Composed tree: skill rewrite, recognizer quote runs, and inline state commit all
present. The round records once with the entity operand — no usage error — and
the gate is prepared. Journey-level green, not just mechanism confirmation.

  exit=0  [state commit]
    $ '${SPACEDOCK_BIN:-spacedock} status --workflow-dir '"'/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow1006978280/002' --set rejection-task sta
    > Committed rejection-task in the inline workflow repository; nothing pushed.
  exit=0  [state commit]
    $ '${SPACEDOCK_BIN:-spacedock} status --workflow-dir '"'/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow1006978280/002' --set rejection-task sta
    > Committed rejection-task in the inline workflow repository; nothing pushed.
  exit=0  [gate record]
    $ '${SPACEDOCK_BIN:-spacedock} gate record rejection-task --round validation/1 --briefing '"'/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow100
    > round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=4
  exit=0  [state commit]
    $ '${SPACEDOCK_BIN:-spacedock} state commit rejection-task'
    > Committed rejection-task in the inline workflow repository; nothing pushed.
  exit=0  [state commit]
    $ '${SPACEDOCK_BIN:-spacedock} status --workflow-dir '"'/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow1006978280/002' --set rejection-task sta
    > Committed rejection-task in the inline workflow repository; nothing pushed.
  exit=0  [gate prepare]
    $ '${SPACEDOCK_BIN:-spacedock} gate prepare rejection-task --question '"'Approve the corrected implementation after cycle-2 validation passed?' --artifact '/var/folders/h1/
    > state=open
  exit=0  [state commit]
    $ '${SPACEDOCK_BIN:-spacedock} state commit rejection-task --workflow-dir '"'/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow1006978280/002'"
    > Committed rejection-task in the inline workflow repository; nothing pushed.

## Ledgers

| file | tree under test | runs | greens |
|---|---|---|---|
| ledger-bare-default-model.tsv | pre-rebase, WRONG codex config | 2 | 0 |
| ledger-ci-config-prequoting.tsv | pre-rebase, CI config | 3 | 1 |
| ledger-ci-config-postquoting.tsv | pre-rebase + quoting fix, CI config | 4 | 2 |
| ledger-composed-tree.tsv | composed tree, CI config | 4 | 2 |

The bare-default ledger is retained as a labeled dataset, NOT as gating evidence:
those runs drove codex's own defaults rather than CI's pinned model and reasoning
effort, so they measure a different configuration.
