# G1 implementation proposal and FO disposition

Candidate inspected unchanged at `49d0c7a962ae500168a3a11b6d91136a59ebdb19`.

- Released user and normal workflow: shipped Pi first officers dispatch asynchronous pi-subagents workers in shared live journeys.
- Observable harm: valid native completions grade as incomplete and mask the auto journey’s downstream missing gate.
- Authority: value-ac[AC-1] Native completion must credit only the dispatched worker at the correct parent boundary.
- Trigger evidence: CI 35058669297 native parent/child records, retained triage, and the existing replay below.
- Proposed materiality: Material.
- Proposed ownership: pi-native-completion-evidence, existing lifecycle grader.
- Proposed disposition: FIX, pending distinct FO authorization.

Before candidate edits, existing `TestPiAutoContinueReplayDoubleDispatch` with `SPACEDOCK_PI_AC_ARTIFACT_DIR=/tmp/spacedock-tip-ci/pi/live-artifacts/pi/pi-common/auto-continue-after-implementation--auto-continue/single-root` returned:

```text
validation lifecycle incomplete: spawns=1 completed=-1 validation=108 report=<nil>
FAIL github.com/spacedock-dev/spacedock/internal/ensigncycle
```

The subsequently added captured test also failed before helper implementation: default `spawns=1 completed=-1 validation=4 report=<nil>`; auto `spawns=1 completed=-1 validation=5 report=<nil>`.

Proposed surface: five existing integration/parser files, existing replay/negative test owners, one 150–200-line narrow Pi helper, and compact projected captured fixtures. Both existing lifecycle call sites receive optional artifactDir; no host interface change, synthetic events, child call concatenation, process controller, or duplicate commit parser.

FO authorization received through the addressable worker boundary before candidate implementation: “FO FIX authorized for G1 as Material and owned under approved brief. Use optional retained-artifact evidence at existing lifecycle callers, one narrow Pi correlation helper and existing test owners. Preserve native parent timeline; no synthetic completion event, child tool-call concatenation, duplicate commit parser, new runtime interface/controller or broad framework.”

Existing durable report/Git/gate checks retain ownership. The original auto gate-reference typo remains outside this correction; final tip CI is deferred to the FO after independent validation. No native/model run, push, PR, CI, or rebase was performed by this worker.

## Self-review: conflicting epoch after a valid notice

While the first race suite was active, added `TestPiNativeCompletionAttributionAndOrdering/conflicting-epoch`: valid epoch-1 notification followed by a second child locator for the same agent/run with epoch-2 metadata. The test failed `invalid completion passed`; the helper ignored the second identity as unrelated.

- Released user/workflow: same native Pi parent notification observer.
- Observable harm: conflicting assignment evidence may retain an earlier credited completion.
- Authority: value-ac[AC-1] Exact run/epoch correlation must fail closed.
- Trigger evidence: focused test above failed before the helper correction.
- Proposed classification/ownership/disposition: Material / this lifecycle correction / FIX.
- Proposed extra surface: approximately six helper lines plus the focused negative case; distinguish an unrelated run from a conflicting epoch within this exact agent/run.

FO disposition through addressable worker boundary: “FO FIX authorized for conflicting-epoch finding as Material and owned by AC-1. Reject same agent/run with conflicting epoch or multiple identities rather than accepting earlier completion; unrelated run stays unrelated. Preserve red case and focused red→green. Let current race finish and retain its result.”
