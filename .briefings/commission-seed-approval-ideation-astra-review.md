## Review

Citation key: **E** = `docs/dev/.spacedock-state/commission-seed-approval-vs-result-review.md`; **D** = `docs/dev/README.md`; **C** = `skills/commission/SKILL.md`; **P** = `skills/present-gate/SKILL.md`; **L** = `skills/fo-gate-lifecycle/SKILL.md`; **H** = `skills/first-officer/references/pi-first-officer-runtime.md`; **G** = `internal/status/entered_stage.go`; **R** = `docs/dev/.spacedock-state/.briefings/fo-contract-reread-churn-astra-review.md`.

### 1. NECESSITY

**Two skill changes are defensible for the stated ACs, but not every proposed addition earns its place.**

| Mechanism | Value served | Simplest alternative; does it suffice? |
|---|---|---|
| Gate-timing guidance during commissioning | AC-1: preserve requested research | Presenter-only correction can expose the mistake before approval, but cannot ensure commissioning produces the requested route. Authoring guidance is justified. However, E’s claim that presentation is “too late to prevent” the bypass is wrong: presentation precedes decision mutation. It is too late to prevent *bad generation*, not necessarily bad advancement. [E:167–169,202; P:13] |
| Confirm Design artifact/successor summary | AC-1: make the approved route explicit | Question 2 alone does not suffice because batch mode skips directly to Confirm Design. Keep this coverage. One canonical timing rule referenced from both paths is enough. [C:15–22; E:66–87] |
| Generated README guidance and five added checklist items | AC-1/AC-2: generated declarations and evidence match intent | The guidance plus one reconciliation check—compare confirmed approvals, artifacts, and successors with generated declarations—can replace several repetitive checklist rows. The checklist is another instruction to the same author, not an independent validator. [E:97–112,203] |
| Initial-gate presenter clarification and contradiction handling | AC-2: truthful decisions, including existing malformed workflows | Commission-only edits suffice for the narrower authoring correction, **not** the stated contradictory-existing-workflow case. Keep a small presenter clarification: existing rules already identify the seed, forbid fabrication, and require concrete decision effects, but still call Gate content an authoritative override. [L:28; P:33,39–42; E:123–131,173,204] |
| Mandatory README read during presentation | AC-2: determine initial status and actual successor | Reuse available workflow declarations when sufficient; read them when unavailable. The lifecycle already needs to identify `initial: true`. E does not establish why a fresh read must be added at every presentation. This can be smaller without introducing presenter-load exemptions. [L:15,28; E:123,135] |
| Interactive/batch live drives and contradictory-seed exercise | AC-1/AC-2; AC-3 observes the same events | Static inspection cannot establish generation, spawn, or root presentation. Live observation is justified; reusing its trace for lifecycle checks avoids a second harness. [E:169–177,189–193,205] |
| Readiness probes | AC-2 regression protection | Reuse existing guard tests where they cover these failures; a positive run alone is insufficient. The design names the existing owner but does not identify its tests, so equivalence is **UNVERIFIED**. These probes do not test the new authoring behavior. [E:181,192,196,206] |
| Worked-YAML validation in the shipped commission checklist | No distinct value established beyond schema validity | **Finding — P2:** keep validation of implementation-added examples in the implementation test plan instead. E proposes no new worked YAML, yet installs a recurring disposable-example obligation in commission. Generated-workflow validation already exists in the proposed proof. Remove E:112 from the shipped addition. [E:112,163,194] |

### 2. PREVENTION OR GUIDANCE?

**Your “can no longer claim” framing is too strong. Nothing mechanical introduced here stops that claim.**

After this change:

- The author is told not to write the misleading gate.
- The presenter is told to explain the contradiction and recommend revision.
- The binary still accepts an initial gate based on tracked, clean entity bytes; that branch does not inspect the truth of its prose. [E:66–68,123–131,157; G:33–38,180–205]

This is **guidance, not prevention**. It improves the instructions at the relevant failure points rather than merely renaming the incident, but compliance remains behavioral and unproven. The design acknowledges that explicitly. [E:220]

**Guidance is a proportionate remedy for this bounded documentation task**, provided the gate approves a tested reduction of this failure mode—not a guarantee of impossibility. The existing non-initial completion guard already supplies the mechanical boundary; adding another validator is not justified by the supplied evidence. [G:33–38,84–93; E:50,161]

**Finding — P2:** replace the unqualified prevention promise with an accurate guidance claim and retain zero false claims/zero bypasses as observed acceptance thresholds. [E:37,168–173]

Your incident wording correctly describes a *prospective* bypass. The entity says the captain stopped the proposal before dispatch; an actual skipped worker is not established. The historical briefing and commit are outside the permitted reading set, so that incident provenance remains **UNVERIFIED**, not disproved. [E:35–37,215]

### 3. SOUNDNESS AND GAPS

The stage requires **at least one** independently measured value criterion—not that every AC independently measures improvement. [D:151]

| Criterion | Assessment |
|---|---|
| **AC-1** | **Real value criterion.** Requested order is fixed independently of generated output; worker and prototype spawn order can contradict it. This satisfies D:151 in design. The historical “one bypass” is a proposed bad route, not an observed baseline run, and must remain labelled that way. [E:35–37,167–169,183] |
| **AC-2** | **Real behavioral criterion.** Root text can falsely claim completed research, and readiness probes can wrongly succeed. Report absence/durability and the fixed requested successor supply independent observations. It mixes changed presentation behavior with unchanged guard regression coverage, but does not merely assert prose shipped. [E:171–173,191–192; G:33–38] |
| **AC-3** | **Supporting compatibility criterion**, not an independent improvement metric. It checks observed sequencing, presentation cardinality, and authority—not merely the existence of a presenter edit. Its schema-only portion proves less, which the entity correctly admits. [E:175–177; H:15,33] |

**No whole AC merely asserts its mechanism shipped.** Their outcomes remain **UNVERIFIED** until exercised; a sound criterion is not a passing criterion. [E:169,173,177,220]

**Finding — P1: the primary proof plan assumes an unreconciled Pi commissioning path.**

The plan requires the same commission-to-FO drives to supply Pi root evidence. But commission explicitly loads the **Claude runtime** and mandates a Claude-shaped boot probe/dispatch path; Pi specifies a different spawn binding. Startup redesign is expressly excluded. [E:169,177,185,193,161; C:665,674–678; H:7]

This does **not** establish that Pi commissioning is impossible. It establishes that the proposed test path is **UNVERIFIED**, and the plan has not said how its instructions compose.

**Smallest fix:** specify and preflight the actual entry path. Either demonstrate the existing Pi commission handoff, or explicitly separate commission generation from Pi consumption of the unchanged generated workflow and stop calling that a continuous commissioning handoff. Any required product startup change needs a scope decision. [E:149,161,220]

### 4. OVERLAP AND SEQUENCING

**No present conflict or dependency is established.**

The earlier review explicitly says presenter loads were excluded from residency coverage and its numerator; extending that coverage was conditional advice, not an already-required change. This design preserves the current per-gate presenter rule. [R:47–55; E:135,157; H:15; L:25]

**Recommended order: land this seed/result correction first once its evidence gaps are resolved.** It can use the existing lifecycle unchanged; waiting for an unrelated residency expansion buys no correctness here. This is a recommendation, not a technical ordering requirement. [E:135,157,161; R:47–55]

If the residency work later expands to cover the presenter, reconcile this design’s load-specific AC-3 checks with that approved change. Reusing the rendering instructions must still preserve fresh committed evidence and one presentation per gate. [E:176–177,193; R:53–55]

Whether the residency implementation has subsequently changed or landed is **UNVERIFIED**: the supplied earlier review identifies a particular patch/snapshot, not current repository ancestry. [R:5,13]

### 5. RISKIEST UNVERIFIED MECHANISM

**The riskiest mechanism is instruction efficacy:** whether an agent preserves the requested work during generation and rejects misleading seed Gate content during presentation. It is only proposed, not exercised. [E:66–131,220]

**Finding — P1: “no spike needed” does not discharge the stage’s risk requirement.**

The entity does record an auditable-looking disposition with source paths, commit identifiers, and reported command results. But its rationale excludes parser/format/flag changes while explicitly leaving the principal mechanism unproven until implementation. The stage requires a spike when the design rests on an unverified mechanism, and says such a design is not gate-ready. [E:210–220,242–243; D:153,163]

Evidence strength differs:

- **Verified from source:** initial readiness uses clean entity durability; later readiness checks report structure and entity durability. This is not proof that substantive research or separate artifacts exist. [G:33–45,84–142,184–205]
- **Verified as instructions:** Pi requires commit, reread, presenter load, and root review; lifecycle specifies successor consumption and dispatch. Written contracts are not execution traces. [H:15,33; L:67]
- **Verified as current entity contents:** backlog approval is recorded as consumed into ideation. The cited commit history and actual spawn remain **UNVERIFIED**. [E:3,17–28,216]
- **UNVERIFIED reported executions:** pre3 version and validator success; they do not exercise the proposed instructions. [E:217,220]

A focused throwaway exercise is enough; demanding the entire implementation validation matrix before ideation would be excessive.

### 6. MINIMAL CHANGES BEFORE PRESENTATION

1. **Correct the claim:** distinguish authoring/presentation guidance from mechanical prevention. Keep the observable zero-bypass and zero-false-claim thresholds. [E:37,168–173]
2. **Exercise the risky instruction path now:** retain a small candidate-instruction trial of malformed seed generation/presentation, including the actual response to contradictory Gate content. Replace the unsupported blanket “no spike needed” with its result and remaining limits. [E:191,210–220; D:153,163]
3. **Resolve the test entry path:** document and preflight how commission reaches Pi without silently substituting Claude instructions or changing the generated route. [C:665,674–678; E:185]
4. **Trim recurring instructions:** remove the worked-example validation checklist addition, consolidate redundant generation checks, and avoid requiring an extra README read when sufficient declarations are already available. Keep both skill files for the currently stated AC-2 scope. [E:108–112,123,173]
5. **State the residency boundary explicitly:** this task neither fixes presenter reload churn nor depends on that fix; future residency expansion must preserve fresh gate evidence. [E:135,161; R:47–55]

**Correct:** the design preserves legitimate initial gates, separates proposed from observed incident harm, and supplies genuinely falsifiable value criteria. [D:14–16,136; E:35–41,167–177]

**Fixed:** none; this was read-only.

**Merge verdict: BLOCK — not yet ready for the ideation gate because the principal mechanism is unexercised and the primary Pi proof path is unresolved.** No binary rewrite, third deliverable file, or residency dependency is justified by these findings.