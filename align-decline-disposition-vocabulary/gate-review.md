# sh — corrected ideation direction

## End value

Generic gate recording must accept a structurally valid review record without
owning the development workflow's finding taxonomy. Workflow-owned labels must
remain available to the workflow that declares them, without a generic parser
allowlist blocking a valid round.

## Recommendation

Move `align-decline-disposition-vocabulary` from backlog to ideation only after
reframing it as a workflow-ownership task.

The ideation worker must:

- remove the proposed allowlist for `deferred-risk`, `polish`, and
  `correct-but-disproportionate`;
- make the generic round path preserve class/disposition content as opaque
  workflow data and validate only record shape, identity, Briefing association,
  includes, and immutable bytes;
- coordinate with WJ, which owns removal of the generic development-policy
  parser; sh owns the regression and workflow-level contract gap that proves an
  arbitrary well-formed workflow class is accepted;
- place taxonomy and semantic disposition rules in workflow-owned policy, not
  in generic gate code.

## Decision ask

Approve to enter ideation and rewrite the task body, acceptance criteria, and
test plan around the opaque-class invariant. A later ideation gate must approve
the exact implementation surface before any code dispatch.

## Evidence

- `docs/dev/.spacedock-state/align-decline-disposition-vocabulary.md`
- `docs/dev/README.md` — Review-finding disposition
- `docs/dev/.spacedock-state/cut-workflow-specific-round-recorder-from-v1/index.md`
- `docs/specs/gate-resolution-frontmatter-contract.md`
