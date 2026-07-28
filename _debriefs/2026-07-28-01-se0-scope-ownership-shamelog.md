# Shamelog: I treated a discovered defect as permission to reshape `se0`

Date: 2026-07-28  
First Officer: Codex session `/root`  
Affected task: `se0v37bt7mhsrmhta1nyns0r` (`live-lanes-red-on-every-branch`)  
Follow-up owner: `q3vpb8hes1b3k3f1jps1kvpk` (`gate-record-stage-coherence-guard`)

## What I did wrong

The complete Sonnet run found a real recorder defect: `gate record` accepted a
validation-qualified Briefing while the task remained in `implementation`. The
worker classified the defect, stopped, and asked for direction. I authorized it
to change the shipped gate contract and then the recorder.

That authorization exceeded `se0`'s approved boundary. The test belonged to
`se0`'s evidence surface; the defect did not belong to `se0`'s deliverable.
Materiality required a durable response, but it did not assign implementation
ownership.

The ensign did not make the authority error. It stopped after each red, reported
the proposed scope, and waited. I gave it permission to proceed.

## Decision timeline

1. The minimized launcher-shim correction made the exact Sonnet
   `gate-guardrail` journey green.
2. The complete Sonnet run exposed two oracle defects. Both belonged to
   `se0`'s approved evidence repair, so the worker fixed them and made the
   focused rejection and keep-moving journeys green.
3. The next complete run left one red:
   `TestLiveDefaultHeadlessStopsAtGate`.
4. The worker proved that the First Officer had bound a validation Briefing as
   an implementation gate. It classified the failure as a material product
   conduct defect and proposed a contract guard.
5. I authorized a prose guard. Sonnet loaded and ignored it.
6. The worker then located the recorder defect: `initialIDs` silently fell back
   to the current stage when a structured Briefing ID named a different stage.
   It proposed a five-file production correction of about 109 additions and 18
   deletions, plus the prose work.
7. I authorized that correction inside `se0`.
8. The work grew to seven files, 254 additions, and 33 deletions. I invoked the
   size tripwire and asked the worker to reduce duplicated tests. I still failed
   to revisit delivery ownership.
9. The captain ruled that the test should become a narrow TODO and the recorder
   defect should receive its own task.
10. The worker removed all uncommitted semantic work. The final quarantine
    changes one test definition, affects the Sonnet and Opus executions only,
    links `q3vpb8hes1b3k3f1jps1kvpk`, and states that green CI does not prove the
    missing gate-stop capability.

## Waste caused

- One unnecessary contract experiment.
- One unnecessary recorder implementation attempt across seven files.
- A 524-second live run interrupted before verdict when the size tripwire fired.
- Review and coordination time spent reducing a patch that the active task did
  not own.
- Temporary confusion about whether a green lane would prove the quarantined
  behavior.

No recorder-semantic code reached a commit. The waste stopped before GitHub CI.

## Why the existing controls failed

### Materiality did not assign ownership

The current classifier asks whether a supported workflow triggers the defect,
whether users can observe harm, which value or safety boundary fails, and what
evidence proves the trigger. It correctly classified this defect as Material.
It does not ask whether the active task introduced the defect or owns its
remedy.

I converted “Material” into “fix here.” The correct meaning is “do not ignore.”

### The expected-surface tripwire fired too late

File and line estimates measure implementation magnitude. They do not decide
whether a semantic change belongs to the task. The tripwire caught the
seven-file expansion only after I had authorized the wrong product boundary.

### The AC described a proof suite as the deliverable boundary

`se0` required every registered journey to pass and said that a skip was not a
pass. I treated every valid failure in that suite as work owned by `se0`.
A test can discover a separate feature defect. The task that owns the test does
not automatically own every product change the test reveals.

### The First Officer contract has no external-defect disposition

The contract says a relevant check must run green and a material finding must
not be waved away. That protects proof quality. It does not define the route
for a genuine defect outside the active work unit: preserve the evidence,
assign a durable owner, state the lost claim, and seek a binding decision about
hold, reshape, or narrow quarantine.

### I misread the conn

The sprint conn authorized gate, CI, PR, and merge decisions within the sprint.
It did not authorize silent expansion of an approved member's product
semantics. I treated execution authority as design authority.

## Classification that would have prevented the incident

Every finding needs three independent classifications:

1. **Trigger and likelihood**
   - exact preconditions;
   - supported or promised path;
   - observed repeatedly, observed once, mechanically reachable,
     adversarial-only, hypothetical, or unknown;
   - measurement when frequency matters and remains unknown.
2. **Impact and release significance**
   - observable harm and reversibility;
   - affected value AC or safety, integrity, security, or compatibility
     boundary;
   - Material, Deferred risk, Polish, or Needs decision.
3. **Delivery ownership**
   - introduced or regressed by this candidate;
   - inside the approved semantic and implementation surface;
   - separate product or dependency defect;
   - unresolved attribution.

The governing rule should be:

> Material means the finding cannot be silently ignored. It does not mean the
> active work unit owns the fix.

For this incident:

- Trigger: current stage `implementation` plus a retained validation Briefing.
- Likelihood: observed once in a real Sonnet journey and mechanically
  reproducible; production frequency unknown.
- Impact: durable gate and Briefing authority disagree.
- Release significance: Material.
- Delivery ownership: pre-existing recorder feature defect outside `se0`.
- Disposition: one linked TODO, truthful loss of proof, and separate task `q3`.

## Proposed First Officer contract change

Add the following bullet under **Working Principles**, inside “Hold your own
gate, merge, and triage calls to the bar you impose on workers”:

```markdown
> - **Release significance does not assign delivery ownership.** Before
>   authorizing an edit for a newly discovered failure or finding, classify
>   both its causal relation to the active deliverable
>   (`introduced-or-regressed`, `pre-existing-or-newly-exposed`, `unknown`) and
>   its relation to the approved boundary (`owned-here`, `outside`, `unknown`).
>   “Material” means the finding cannot be silently ignored; it does not mean
>   the active entity owns the fix. An outside or unknown finding stops edits:
>   preserve the evidence and seek a binding route—reshape or hold the active
>   entity, assign a linked work unit, or narrow the blocked claim with the
>   captain's approval. The conn authorizes execution within the approved
>   boundary; it does not authorize silent reshaping.
```

Add this sentence to the existing size/surface guidance:

```markdown
File and LOC tolerance measures the cost of an already-owned correction. It
never substitutes for the delivery-ownership ruling above.
```

This belongs in the host-neutral contract because the failure can occur in any
workflow when a review, validation, pilot, or user report discovers a genuine
defect outside the active work unit.

## Proposed `docs/dev/README.md` changes

### Ideation: declare semantic surface

Replace:

```markdown
The task body declares an expected surface — the files and LOC it expects to
touch — and its tolerance...
```

with:

```markdown
The task body declares an expected implementation surface—files and LOC with a
tolerance—and the observable semantics it may change, such as command,
storage, authority, or runtime behavior. File/LOC deviation measures cost.
A remedy that changes an undeclared observable semantic is a boundary change
even when its diff is small.
```

### Implementation: add ownership before fixing

Insert before the current in-stage review-round bullet:

```markdown
- Before editing for any newly discovered failure—review finding, failing
  test, pilot result, or user report—classify three axes: defect kind, release
  scope, and delivery ownership. Delivery ownership records causal relation
  (`introduced-or-regressed`, `pre-existing-or-newly-exposed`, `unknown`) and
  approved-boundary relation (`owned-here`, `outside`, `unknown`). Fix a
  Material finding in this stage only when it is owned here. An outside or
  unknown Material finding is `Needs decision`: preserve its counterexample,
  report the exact blocked claim, and wait for the FO's binding route before
  editing. File/LOC tolerance does not grant ownership.
```

Change the current sentence:

```markdown
Fix material findings...
```

to:

```markdown
Fix Material findings owned by this task. Route outside or unknown Material
findings through the delivery-ownership rule above.
```

### Validation: make ownership a third axis

Change:

```markdown
Classify every finding on two independent axes:
```

to:

```markdown
Classify every finding on three independent axes:
```

Keep the existing defect-kind and release-scope bullets, then add:

```markdown
- **Delivery ownership:** record whether the candidate introduced or regressed
  the defect, merely exposed a pre-existing defect, or leaves attribution
  unknown; then record whether the approved task design owns the remedy.
  Materiality controls urgency, not ownership. A Material outside or unknown
  finding blocks silent edits and requires an FO/captain routing decision.
```

### Dev-only evidence quarantine

Add under validation:

```markdown
- **Evidence quarantine:** A captain may narrow a blocked proof while a linked
  task owns an external defect. Name the exact test definitions and executions
  removed from the proof, the lost AC claim, the durable owner, and the
  re-enable condition. Amend the active task's AC and test plan before claiming
  green. A skip is neither a pass nor evidence for the quarantined behavior.
  Do not disable a whole lane when one test can be isolated.
```

This paragraph belongs only in the development workflow. The host-neutral FO
contract should define ownership and routing, not test-specific mechanics.

## Required behavioral proof

Do not prove these changes with prose-grep. Add one workflow fixture:

- An active task owns a bounded deliverable.
- Validation discovers a supported, Material defect outside that deliverable.
- The worker reports release significance and delivery ownership before editing.
- The FO does not authorize a product edit under the active task.
- The durable result either holds/reshapes the task or links a separate owner
  and narrows the blocked claim with explicit captain authority.

The fixture must fail if the worker fixes the external defect, dismisses it as
non-material, or calls a skipped claim passed.

## Recommendation

Adopt both layers:

- the host-neutral FO rule that separates significance from ownership;
- the dev-workflow classification and narrow-quarantine mechanics.

Do not add a CI-specific rule. This incident concerns ownership of discovered
work, regardless of where the evidence came from.
