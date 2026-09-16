# Claude rejection-flow timeout follow-up

Run 35058663190; source ac3c06a445fcfbc773639254ff7242965248ba78. Inspection tip ba8e6eeb4e28e310b17897a0a8b12e151366b8f9. Read-only investigation; no model calls, tests, code edits, CI, or publication. Original diagnosis: /tmp/spacedock-tip-ci/claude-rejection-diagnosis.md. Exact selected transcript lines and source hashes accompany this report.

## Finding and disposition

The retained stream demonstrates a quiet-timeout chain, not a missing feedback schema or failed dispatch. At 05:19:09.400Z the FO commissioned an Explore worker to discover the feedback-context-file format. At 05:19:11.400Z that worker ran recursive grep from `/`. Its output filter for `/proc/` did not prune traversal. The FO's last public event at 05:19:21.506Z says it is waiting for that research. The child tool returned only at 05:21:11.443Z, after its 120-second timeout. The runner allows 60 seconds without stream output and reported a quiet timeout. No feedback dispatch had been attempted.

- Finding: unnecessary format research left the stream silent beyond the existing bound.
- Evidence: selected public stream lines 268, 274, 294; child lines 12–13; existing CI timeout diagnosis.
- Classification: observed agent detour plus expected harness timeout; no demonstrated deterministic instruction defect.
- Disposition/owner: retain and compare the already arranged repaired-tip Claude lane. FO owns routing; no fo-continuation product edit, new discovery mechanism, or timeout relaxation justified by this one sample.

The feedback contract already defines opaque transport (skills/first-officer/references/fo-dispatch-core.md). internal/dispatch/build.go reads the file and transports its bytes as a string. The feedback-rejection skill prescribes forwarding the authorized correction package unchanged. A context-budget lookup failed before the detour, but unavailable budget has a defined fresh-worker branch; it does not require schema discovery. These facts rule out an observed parser rejection as the trigger. They do not prove why the model chose to research.

## Scheduling attribution and merge scope

Commit 4ce49f1ea (#800) replaces common-journey scheduling with three duration-ordered slots and folds Claude substrate journeys into that pool. The shared runner avoids a second t.Parallel inside a scheduled slot. Existing parallelism remains three, but job overlap changes.

That diff does not change Claude's 60-second quiet budget, stream watcher, feedback handling, or per-scenario CLAUDE_CONFIG_DIR isolation. The failed run progressed through implementation, validation and rejection delivery before the search. No credential collision, config overwrite, scheduler wait, or resource-starvation event is established by the retained failure evidence. Therefore there is no evidence attributing this timeout to #800. This is not demonstrated independence: changed overlap could affect runtime performance, and there is no controlled comparison in this investigation. This unrelated upper-lane failure alone does not invalidate #800's separately established targeted scheduling proof or justify blanket-blocking the bottom PR.

## Smallest bounded local proof and owner

No new standing harness is needed to reproduce the observer boundary. Existing non-model tests are:

```sh
go test ./internal/ensigncycle -run '^TestDrainToExit(KillsStalledStream|ResetsDeadlineOnActivity|ReturnsFullTranscriptOnExit)$' -count=1
go test ./internal/dispatch -run '^TestBuildAdvanceGoldens$/^feedback-reflow$' -count=1
```

The first set exercises silence causing kill, continuing activity resetting the deadline, and normal exit retaining output. The second builds feedback dispatch with ordinary text (`REJECTED: do better.`), without a schema. These establish deterministic harness and transport behavior, not recurrence of the model detour. Commands are proposed, not executed here; no live build tag or model invocation is needed. Existing test owners are the ensigncycle stream-watcher harness and dispatch advance builder, respectively. No parser fix is indicated.

## Separate cleanup observation

The child transcript continued after the public-stream cutoff. At 05:21:12.778Z its tool reported that the test cwd had been deleted; later child events continue into 05:22. This establishes descendant activity after fixture cleanup. internal/ensigncycle/streamwatch_test.go cmdPoller.kill calls Process.Kill on the direct command only. Exact process ancestry and launcher signal behavior have not been reconstructed, so direct-kill behavior is a plausible cleanup mechanism, not a fully demonstrated root cause. This later activity cannot explain the initial silence.

If routed, the smallest additional non-model proof belongs to the live-runner/process-cleanup owner: use an isolated shell launcher with a child that writes a marker after a short delay; drive the existing cmdPoller/streamWatcher with the small test quiet budget; verify parent termination and whether the descendant marker appears after cleanup. Record PIDs and explicitly reap the test child. Do not scan `/`, launch Claude, or introduce a controller. Only if that local reproduction confirms the leak should the owner propose a scoped process-tree cleanup change with a normal-exit control. This investigation neither implements nor claims that fix.
