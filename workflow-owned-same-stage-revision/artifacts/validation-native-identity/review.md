# Independent native identity and recorder evidence review

Candidate67c68536fd6a49461624f1c6d86081d22bef4811; existing V3a/V3b and adjacent AC-2 recorder evidence FO FIX authorizations. Detached checkout only. No candidate edits, model/native workflow, broad suite, rebase, push or CI.

## Native identity

Prior tip investigation established two distinct Claude task IDs under one name and one Codex task with two completed turns. Producer captured-regrade.log executes both retained full streams against the corrected extractor/grader and gets exactly those identities. Reused that positive evidence rather than running it a second time. Inspected metadata correlation: Claude Agent tool ID is correlated to returned task ID; matching completion must satisfy opened tool and available task owner; routing label is used for followup lookup, not reviewer identity. Stage association persists on the opened tool. The grader still parses ordered dispatch/completion pairs, requires nonempty identities, rejects repeated fresh identity, and allows reuse only of the current worker at its original stage. Independent review requires a second distinct fresh identity. The caller detects reviewer dispatch through fresh-spawn count instead of number of turns.

Independent falsifiers mutate the actual retained streams: replace the second Claude native ID with the correction ID; remove only the second native completion; retarget Codex's completed followup to an unspawned owner. Each fails the actual grader. Retained owned tests additionally cover wrong tool/task completion owner, changed-label self-review, unfinished reuse, extra fresh worker, self-review via reuse and gate-selected revision negatives. No guard over durable outcome or authority was removed.

## Actual recorder evidence

The conventional Claude runner now enters the same existing logging shim path as Codex. Its actual withStubPATH method sets shell startup interception; this matches the producer's real launcher→non-model Claude stand-in test. The shim invokes the real binary with "$@", captures that process's exit, appends its command observation and returns the same exit. Publication parsing accepts actual exit0 gate-record observations for the expected entity/round. It does not use enclosing shell success or an echo supplied by the model.

Reused the producer's executed launcher check: incomplete canonical log produces recorder exit1 under a successful shell; complete log produces exit0 under another successful shell; only one publication counts. The original execution matrix also rejects echo-only, masked recorder failure, duplicate/second-round true successes and missing room. Independent parser boundary matrix passes: absent log, exit1 plus successful echo, wrong command, wrong entity and duplicate successful observations all reject; one failed observation followed by one true success passes. These pure negative checks do not replace the producer's real interception proof.

The retained premature failed recorder call remains observed conduct. This repair changes evidence attribution and does not claim flawless model behavior. Pi uses its existing separate path; this review does not claim Pi completion or fix it.

## Outcome and scope

All frozen-input/gate-selected bytes, old attempt immutability, round evidence cardinality, review hold, cycle-limit and no-successor checks remain in place. No skill, CLI, gate/recorder product code or workflow policy changed in this correction. Exact scope is161 additions/27 deletions =+134net/five existing files; cumulative above4ce49f1ea is+571net/ten files, matching the FO-authorized report.

Independent command: `go test ./internal/ensigncycle -run '^TestDetached(CapturedIdentityFalsifiers|RecorderEvidenceMatrix)$' -count=1 -v`, pass0.317s. Source and output retained. Reused producer focused normal7.003s/race8.092s and registry1.005s with command metadata and source hashes. No duplicate green suites were run.

Recommendation PASSED for the deterministic correction; V3a/V3b satisfied, no new material finding. Final combined normal/race and final host CI remain FO-owned pending Pi correction. Prior broad results and prior tip CI are not new-head full/runtime verification; known installed-manifest resolver failure remains limited to earlier broad runs.
