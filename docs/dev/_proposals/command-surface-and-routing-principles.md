# Command Surface & Routing Principles (anti-sprawl)

> **Decision record, drafted 2026-06-17 (Shaping FO + captain).** Status: **ACCEPTED** (captain-ratified 2026-06-18). Stands as the **reviewer rubric** applied at the ideation gate. The binary-enforcement task (`command-surface-gate` / dwd) was deliberately **shelved** — the principles are enforced by reviewers, not (yet) by a contractlint allowlist.
> Prompted by the 0205 question "should `merge guard` / `state` verbs be standalone `spacedock-merge` / `spacedock-state` binaries routed through the frontdoor?" (merge is a mod, state is pluggable). The answer (keep core; pluggability already lives at the workflow layer) generalized into the surface principles below so the command surface stays lean as 0205+ adds verbs.

## Why this exists

0205 (layered-FO) deliberately ADDS binary verbs (state, merge-guard, gate-extract, and later next-action) to push mechanical FO work into the binary. Without an explicit rubric, "add a command" is the path of least resistance and the surface sprawls. We have already been applying an implicit rubric by feel; this records it so the next verb-author and the gate reviewer apply the same bar.

## The routing mechanism (decided)

**The frontdoor is compiled-in cobra (`root.AddCommand()`, internal/cli/cli.go). There is NO external `spacedock-<verb>` PATH-discovery / plugin-routing layer, and we are not building one now.** Verified: no dispatch seam, registry, or fallback exists; an unknown verb exits 2.

- **Pluggability lives at the WORKFLOW-DEFINITION layer**, not the binary layer: the README `merge:` / `state:` keys + registered `_mods/*` hooks encode each workflow's choice. The binary READS that choice generically (`resolveMergePolicy`, `scanMods["merge"]`, `ClassifyState`).
- A core verb is **backend-AWARE** (it reads the workflow's policy/mods) **without being backend-routable** (the implementation does not swap out of process). Awareness is what "merge is a mod / state is pluggable" actually needs; an external-binary routing layer is a different, speculative thing.
- **If a second backend ever ships** (a non-git state backend, a second merge mod), prefer **INTERNAL dispatch-on-policy** — a policy-selected handler or a sibling subcommand inside the binary — NOT an external routing layer. External routing is the expensive option that would jeopardize the non-optional exit-code property the layered-FO safety thesis rests on (exit-3 HALT, guard-refusal exit 1 surviving a process boundary). Open it only when a concrete second backend lands, and spike it as its own effort.

## The surface-placement decision tree (apply top-down; stop at the first match)

1. **Is it a READ / projection over data an existing command already reads?**
   → a **flag/mode on that command**, never a new command.
   *Precedent:* `gate-extract-verbs` (6re) dropped the proposed `spacedock gate` command and became `status --read --stage X --checklist` / `--ac-scan` modes.

2. **Does it select/narrow output the way an existing flag already does?**
   → **reuse the existing flag**, don't invent a variant.
   *Precedent:* frontmatter projection on `--read` reused the existing `--fields` flag (shared with `--where`/`--next`) rather than a `--read`-specific one.

3. **Does an existing command already compute it, read-only, with no novel risk?**
   → **point at it** (a guillemet prose-function or a doc reference), don't wrap it in a new verb for namespace tidiness.
   *Precedent (proposed):* `«state.sweep-merged»` stays a guillemet pointing at `dispatch reconcile --include un-advanced-pr` rather than a thin `state sweep` verb.

4. **Is it a NEW mechanical INVARIANT the binary must OWN** so a weak (Haiku) FO cannot reorder / skip / collapse it?
   → a **core subcommand under the relevant noun-group**, compiled-in. This is the layered-FO safety lever: the verb is the first-line driver; the existing status guards become the backstop.
   *Examples:* `state commit <slug>` (enforces the rebase-conflict HALT via exit 3), `merge guard <slug>` (atomic mod-block set→clear→terminalize).

5. **Is the behavior WORKFLOW-VARIABLE** (a per-workflow choice of WHAT happens)?
   → a **mod** (a registered `_mods/*` hook the binary reads generically), NOT baked logic. The verb stays backend-AWARE (reads the mod/policy); it never bakes a specific mod's behavior in.

6. **Only then: a new top-level command** — and only when it is a genuinely distinct phase/noun (e.g. `merge` is a distinct loop phase from `status` / `dispatch`). Default to a subcommand/flag/mod on an existing noun before minting a new top-level verb.

## The anti-sprawl gate (the rule reviewers apply)

- A new **top-level command** must justify why it is not a flag, a mode, a subcommand of an existing noun, or a mod. The burden is on adding surface.
- **Reuse a flag** before adding one; **point at** an existing read-only capability before wrapping it.
- A verb earns its place by owning a **mechanical invariant** (safety) or a **deterministic extraction** a weak model re-derives unreliably — not by namespace tidiness.
- The ideation gate's AC cross-check includes: "could this be a flag/mode/mod instead of a command?" If yes and the answer isn't a safety-owned invariant, it routes back.

## Application to the 0205 verb-core (the case that prompted this)

- **`merge guard`** → core subcommand (rule 4): owns the atomic terminal-ceremony invariant; reads `merge:` policy + `scanMods["merge"]` generically (rule 5 satisfied — never bakes in pr-merge). Gate-ready as designed.
- **`state ready` / `state commit`** → core subcommands under the existing `state` noun (rule 4): own the path-scoped-commit + rebase-HALT invariant; reuse the single git interpreter (`runGit`/`ClassifyState`). Gate-ready as designed.
- **`state sweep`** → guillemet pointing at `dispatch reconcile` (rule 3): read-only, no invariant to own. (The one open notation call — verb vs guillemet — is recorded in `state-verbs.md`.)
- **Standalone `spacedock-merge` / `spacedock-state` routing** → NOT pursued: net-new mechanism, YAGNI (one merge mod, one git backend exist), and it weakens the safety property. Located at the workflow layer instead.

Each verb ideation's "Out of scope" gains one line: *pluggability lives at the workflow-definition layer; the verb is backend-AWARE not backend-routable; a future second backend arrives as internal dispatch-on-policy, not external-binary routing.*
