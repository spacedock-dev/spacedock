# Validation gate — repair-pi-recorded-gate-lifecycle

## Outcome

Validated after a criterion amendment, with one criterion restated and its undelivered part recorded.

## Candidate

Layer 3. One commit: the `fo-gate-lifecycle` Prepare step anchors its selected-source path to the observed
root, and the `default-headless-gate-stop` Pi binding is cleared.

## Evidence

- The diff changes only that Prepare prose and that binding plus its comment.
- Shared assertions are byte-unchanged; no masking binding is added.
- Focused gate-prepare tests pass; `TestFOInstructionComponentCaps` passes; the skill file is 7695 of 7700
  bytes.

## Restated criterion, and what is not delivered

Its value criterion named the Pi **delegated-authority** journey `recorded-gate-lifecycle`. The work addresses
`default-headless-gate-stop`, a different journey whose assertion stops open at the gate and forbids
record/consume/successor dispatch. So this layer claims the prepare-path fix and the headless binding clear,
and **does not** claim the delegated-authority journey. That journey needs its own work if it is still wanted.

## Pending

- The live lane clause. There is no green lane at this tip.

## Question

Do you accept the validation of this layer? This approves the validation only, not any merge.
