---
id: mp2jx24h3c92ef1yz9w1tjpz
title: Fix four confirmed tautological tests in internal/ensigncycle and internal/status
status: backlog
source: "Found via a per-package audit sweep applying the 'testing-without-tautologies' checklist (github.com/kenn-io/middleman skills/testing-without-tautologies/SKILL.md), each confirmed by actually applying a suggested production-code mutation and observing the test stay green (not just static reasoning). Four findings: (1) internal/ensigncycle/pty_session_test.go TestFOSessionPinning/activeSessionFile_would_flip_to_teammate — the only failure-path action is t.Logf, no t.Error/t.Fatal/t.Fail anywhere in the subtest, so it passes regardless of activeSessionFile's actual behavior. (2) internal/status/native_new_test.go TestNewFolderForm — computes its 'expected' bytes via the exact same unexported stampID function the code under test also calls internally (mirror assertion); a real splice-offset regression in stampID left it green. (3) internal/status/zz_independent_parity_test.go TestIndNewAtomicCreate — identical mirror-assertion defect against the same stampID function, independently reconfirmed. (4) internal/status/boot_probe_parity_test.go TestBootTeamStateProbeConfinement — builds its 'expected' string directly from the same production constant (teamStateNeutralHint) the code renders from; emptying that constant left it green. Scope for ideation: for each, replace the mirror/no-op assertion with one that has an independent oracle (a literal/hand-checked expected value, or an assertion with a real failure path) while preserving whatever real coverage the test does provide elsewhere in the same function."
started:
completed:
verdict:
score:
worktree:
issue:
---

Four tests across internal/ensigncycle and internal/status assert nothing that can actually fail — confirmed by mutation-testing each one. Fix them so they have an independent oracle instead of mirroring the production logic or dropping the assertion entirely.
