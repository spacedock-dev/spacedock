# Independent staff review — rq cycle 3

Verdict: **REVISE**.

## Material

### 1. The public gate entry still makes the agent reconstruct frozen authority

Cycle 3 keeps a public semantic vector containing actor, approver, entity, room,
workflow directory, and Briefing, and directs the First Officer to copy actor and
approver from `request.json` (`index.md:936-947`). That is the same model-authored
composition boundary the fixed Subspace entry is meant to remove. These values already
belong to one bound room; accepting six correlated values creates mismatch cases instead
of deriving one authority.

The current Subspace skill demonstrates the intended ownership rule: its fixed entry
owns path resolution, provider-package allocation, launch, validation, and cleanup, and
the caller does not choose private paths
(`plugins/subspace/skills/r/SKILL.md:30-36,95-109`). The rq profile should extend that
rule, not expose more internal coordinates.

**Required correction:** make the sole public form exactly
`/subspace:r gate <room>`. The agent passes no actor, approver, entity, workflow root,
Briefing, destination, package path, manifest path, or terminal command. From that room,
fixed Spacedock/Subspace code must resolve and validate:

- the active entity and definition/state workflow roots;
- the frozen request, actor, approver, Briefing locator/id/digest, and logical root map;
- the exact `ROOM/provider` package and `resolved-sources` child; and
- the child argv, recorder room, and validator inputs.

The fixed entry may call Spacedock with the room, but the agent must not read
`request.json`, copy its fields, or reconstruct the current verbose vector. Wrong room
or any derived mismatch must fail before package allocation, capability probing, host
preflight, or provider launch. If a room alone cannot establish those facts, rq and the
unlanded s4 contract must return to joint ideation; adding public arguments is not the
fallback.

### 2. The resolved-source manifest under-binds the canonical Briefing

The closed manifest currently carries only `"briefing": "<id>"`
(`index.md:181-207`). Items bind source tuples, but nothing in the manifest binds the
exact canonical Briefing bytes from which those tuples and the summary were derived.
An id-preserving Briefing substitution can therefore reach the resolved loader unless
another layer happens to catch it.

**Required correction:** bind both the canonical Briefing id and its full
`sha256:<64-lowercase-hex>` canonical digest in the closed manifest. Subspace must read
the original Briefing unchanged, recompute its canonical digest, and require both id and
digest to match before it installs any resolved source or opens the TUI. Missing,
duplicate, unknown, malformed, id-only, digest-only, and recomputed-manifest substitution
cases must fail before display.

The canonical Briefing must remain byte-identical. Resolved Artifact and Reference
payloads are separate ephemeral inputs: Subspace re-hashes them, injects their bytes
only into the in-memory Artifact/Reference models, and derives Result and presented
inventory from the unchanged Briefing. Temporary paths never enter canonical identity.
The primary Artifact summary remains the exact canonical string value; printable text
and spaces render unchanged, unsafe controls render through the reversible control-safe
form, and neither the manifest nor a derived sidecar duplicates the summary.

## Direction that stands

The underlying bridge is necessary. Current Spacedock exposes no `gate prepare` or
`gate materialize` command (`internal/cli/cli.go:163-174`), while current Subspace reads
Artifact URIs as filesystem paths (`internal/reviewv1/loader.go:13-34`) and rejects
Reference URIs containing `://` (`internal/reviewv1/log.go:556-578`). A root map alone
cannot present either source.

The fixed `ROOM/provider` rendezvous, provider-neutral materialization manifest,
unchanged Briefing, separate in-memory Artifact/Reference bytes, and local-only
Spacedock Git resolution remain the smallest coherent ownership split. Native Git in
Subspace, a rewritten Briefing, remote acquisition, q0 transport, a durable cache, and
`association.json` would all add a second authority.

Frozen actor/approver checks also remain necessary, but they must be derived from the
room. Current Spacedock checks request authority only during later room recording
(`internal/gates/operation.go:255-310`); without the pre-launch check, a wrong identity
can consume the one-shot provider location before the recorder refuses it.

The cleanup lifecycle is load-bearing because `ROOM/provider` is retained evidence.
The shell owns failures before dispatch, the TUI deletes verified payloads after
in-memory load, and the supervisor owns child/signal cleanup before publishing exit
evidence. The current supervisor only runs the child and publishes its exit
(`cmd/subspace-tui/provider_supervisor.go:40-96`), so the proposed supervisor work is
not invented. Honest SIGKILL/power-loss residue remains the correct limit.

The validator extension is likewise required. The current canonical validator accepts
exactly the nine-element ordinary child argv
(`plugins/subspace/skills/r/scripts/validate-one-file-result:4-34`); a resolved-source
child cannot become trusted merely by recomputing `argvSha256`.

The declared 10-file/1,209-LOC Spacedock and 22-file/2,508-LOC Subspace surfaces map to
real owners, and 2,084 of the 3,717 changed lines are tests, fixtures, or E2E proof.
Their size is defensible for the corrected contract, but the present counts are not
approvable: the public grammar, materializer inputs, manifest shape, validator matrix,
and E2E assertions all change under the two corrections above. Rebaseline every named
file and retain the existing hard reset triggers.

## Dependency and re-review bar

s4 is still `status: ideation` with no worktree
(`prepare-provider-neutral-gate-room/index.md:1-10`), and current `main` contains
neither its prepare command nor rq's materializer. rq must not create an implementation
worktree against state design `b739a016`.

Before another rq gate:

1. revise the contract and ACs around `/subspace:r gate <room>` and room-derived
   authority;
2. add the canonical Briefing id + digest manifest binding and its pre-display negative
   controls;
3. preserve byte-identical Briefing handling, separate in-memory resolved bytes, exact
   control-safe summary display, cleanup ownership, and canonical validator delivery;
4. rebaseline both repository surfaces and the real moved-root E2E; and
5. state that implementation begins only after s4 lands, its exact tip is recorded, and
   rq re-enters ideation if the landed room/request/root-map/summary contract differs.

This review does not certify sprint closure. The close-out audit explicitly says rq and
s4 remain unmerged and requires another detached audit plus live-lane evidence on the
actual tag candidate (`docs/roadmap/durable-decisions/staff-review-sprint-close.md:14-24,662-676`).
