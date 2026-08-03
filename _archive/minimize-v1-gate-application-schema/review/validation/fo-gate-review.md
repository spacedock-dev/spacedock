# Validation gate review

The implementation candidate is not ready to enter `done`.

## Recommendation

REJECTED for correction round 2. The validator’s focused tests and full/race reruns pass, but three material acceptance-criterion findings remain:

- AC-1: a detached extra `Policy` YAML leaf remains accepted by the recorder/replay evidence, so the exact `{target-stage,state}` key boundary is not proven.
- AC-2: the promised checked-in 31-path active/archive manifest and iterator are absent; current state also reintroduced `application.action: advance`, while the candidate reader rejects newer concurrent gate fields.
- AC-3: the consumed CLI projection `application=advance/consumed` is not asserted; a detached wrong-action mutation stays green.

AC-4 passes for the implementation commit’s 14-file, +138/-177 surface and the 27-path normalization budget.

## Science Officer advisory

REVISE/SEND BACK: add exact YAML-node assertions, add the consumed-action CLI assertion, check in and iterate the 31-path manifest, and reconcile the current pilot paths under First Officer state authority. Do not add compatibility decoding; escalate for design reset if the concurrent fields cannot be reconciled inside the approved boundary.
