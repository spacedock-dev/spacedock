---
id: ywjfm1fjmqbpqe7dwjcwt841
title: README slim + docs-architecture cleanup — site is canonical, remove install-journey.md + its prose-grep test
status: ideation
source: "captain (2026-06-13) — now that the mkdocs doc site shipped (#343), the repo's standalone install/usage docs duplicate it. README should be a thin front door; install-journey.md is redundant; its test is a banned prose-grep."
started: 2026-06-13T05:03:08Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0201-post-flip-release-model
group: docs
sprint-readiness: ready
---

Now that the doc site (`docs/site/`, shipped via #343) is the canonical browsable documentation, the repo carries redundant docs. Establish the rule **site is canonical; the GitHub README is a thin front door**, remove the duplication, and — by captain decision — make this task the **sole owner of the whole 0.20.1 sprint's user-visible doc surface** so the five behavior members touch NO docs.

> **Branch grounding (read first — verified during ideation, the entity's original framing was off).** This task targets **`main`**, the post-flip trunk (roadmap 0201; the doc site shipped via #343 lives on `main`, NOT on `next`). `main` and `next` have **diverged** on exactly this surface — the true per-branch state (`git cat-file -e`):
>
> | file | on `main` | on `next` |
> |---|---|---|
> | `docs/site/` (the canonical site) | **present** | absent |
> | `docs/install-journey.md` | **absent** (already removed when the site landed) | present (169 lines) |
> | `internal/release/install_doc_test.go` | **present, repointed** to read `docs/site/get-started/install.md` | present, reads `docs/install-journey.md` |
> | `README.md` | unslimmed, but Install already links `docs/site/get-started/install.md` | links `docs/install-journey.md`, has How-it-works/Usage |
>
> **What this changes for the task.** On `main` the file-deletion is already done — `install-journey.md` is gone and its content lives at `docs/site/get-started/install.md` (verified: the site page IS that content). The **prose-grep test still exists on `main`**, just repointed: `readInstallJourney` now `os.ReadFile`s `docs/site/get-started/install.md` and the two tests still assert author-written substrings against it. **It is still the banned prose-grep** — repointing it at the site page did not fix the tautology, it moved it. So on `main` the live deliverables are: (1) remove the repointed prose-grep test, (2) slim the README, (3) own the five behaviors' site docs. AC-1's `install-journey.md`-absence is already true on `main` (the grep-clean is the live half there); the deletion is only a real edit if work lands on `next` first and merges.

## Problem

- `internal/release/install_doc_test.go` is a **banned prose-grep**, and repointing it at the site page (as `main` did) did not cure it. Both tests `os.ReadFile` the install doc and assert substrings the author wrote — a `curl … install.sh | sh` line (`TestInstallJourneyDocumentsLinuxPath`) and the absence of "sandboxed on linux" phrasings (`TestInstallJourneyDoesNotOverclaimLinuxSandbox`). The expected values come from the file under test, not an independent source; a valid paraphrase fails them and an inverted clause passes them. This is the instruction-file-read / prose-grep the workflow quarantine forbids outside `internal/contractlint`. (The test names still say "InstallJourney" though `main` repointed the read to the site page — a stale name, but the removal makes that moot.)
- `README.md` carries "How it works" / "Usage" / "Quick start" prose that duplicates the site's `concepts/` + `running-workflows/` + `reference/command-reference.md`.
- `docs/install-journey.md` (169 lines on `next`) is a duplicate of `docs/site/get-started/install.md` (the canonical page). On `main` it is already gone; on `next` it is dead weight. The grep-clean (nothing references it in live docs/code) is the durable invariant.

## Anti-collision: yw owns the whole sprint's doc surface

By captain decision, this task is the **single writer** of every user-visible doc change in 0.20.1. The five behavior members (`gj`, `te`/`tes`, `gp`, `zrc`, `8p`) each authored a doc-diff in their ideation bodies; **those doc-diffs are SUPERSEDED by this task** and the members ship code-and-tests only, touching no `*.md`. The reason is concrete: every member's doc-diff edits `docs/install-journey.md` — the exact file this task **deletes**. If each member applied its own diff, they would collide on a file that no longer exists, and the README's launch/version lines would be edited by three hands. Consolidating the doc surface here removes the collision and lands each behavior's documentation on its **canonical site page** instead of the dead standalone doc.

### Canonical home per member (the retarget map)

For each behavior member, the table names where its now-superseded doc-diff content lands. The member's original diff (cited) is the intended wording; the canonical home is the site page or README section this task edits instead.

| Member | User-visible change | Member's original diff target (SUPERSEDED) | **Canonical home (yw edits)** |
|---|---|---|---|
| **gj** (`gja5h…`, startup-sandbox-status) | `--version` gains a `Sandbox:` line + per-runtime install/enablement block; startup banner + `status --boot` gain sandbox state | `docs/install-journey.md` two `--version` spots; `README.md:139` | `docs/site/reference/command-reference.md` (the `--version` section — its example `spacedock 0.20.0 (contract 1)` gets the `Sandbox:` line + per-runtime block); `docs/site/get-started/install.md` two "Confirm it" steps (the `spacedock --version` example output); README `--version` usage line if README retains a usage block (see AC-3 — it does not) |
| **te**/**tes** (`tes9t…`, install-refresh + upgrade hint) | doctor + front-door gate print an opt-in upgrade hint on a contract-compatible-but-behind plugin; `spacedock install --host codex` refreshes a present plugin | `docs/install-journey.md` "Keep things in sync" | `docs/site/get-started/install.md` "Keep things in sync" section |
| **gp** (`gpvg3…`, marketplace + pinned channels) | edge channel (`spacedock-edge`) is a first-class marketplace channel tracking `next`; stable serves a pinned tag | `docs/install-journey.md` "Build from source" edge note; `docs/releasing.md` (the manifest-in-plugin-branch lines) | `docs/site/get-started/install.md` "Build from source" / a new "Edge channel" note; `docs/site/contributing/releasing.md` (the releasing page on the site). NOTE: the deep release-ritual rewrite (stamp-then-tag) is `ezn`, **deferred** out of 0.20.1 — yw makes only the **decouple-accurate** edits gp's diff named (remove now-false plugin-branch-manifest references; add the edge-channel pointer). |
| **zrc** (`zrcmx…`, non-sandboxed auto-mode) | unsandboxed `spacedock claude` starts in `--permission-mode auto`; unsandboxed `spacedock codex` gets the captain-chosen equivalent (`--ask-for-approval on-request` under option A) | `README.md` launch lines; any `docs/` enumerating the inner host argv | `docs/site/reference/command-reference.md` (the Launch section, where the inner argv flags are described) + `docs/site/get-started/first-launch.md` "command grammar" if it enumerates the permission posture. README launch lines only if README keeps a launch block (AC-3 says no). **Pending the captain's codex option A-vs-B choice** (zrc's open gate item) — yw writes the claude half (fixed) and the codex half per the captain's pick. |
| **8p** (`8pwjd…`, brew agentsview dep) | `brew install spacedock` now also installs `agentsview`; safehouse stays a separate caveat | `docs/install-journey.md` Homebrew install step | `docs/site/get-started/install.md` Homebrew install step (the one-sentence "also installs agentsview" note) |

**Sequencing inside the Commander drive (recorded, not an AC):** yw's *structural* cleanup (delete `install-journey.md` + its test, slim the README) leads and can land immediately — it depends on nothing from the five members. yw's *behavior-doc* updates (the table above) follow each member **landing its code**, because the doc describes shipped behavior; a doc that describes `--version`'s sandbox block before `gj` ships is a lie. So the drive order is: (1) yw structural cleanup → (2) the five behaviors land code+tests → (3) yw's behavior-doc pass updates the canonical site pages. The whole surface targets `main` under the post-flip trunk model.

## Proposed approach

1. **Slim README** to the front-door essentials a repo visitor/contributor needs without leaving: one-line what-it-is + the "what's different" hook; the single install path (`brew install …` → `spacedock claude`) then a prominent **"full docs"** link to the published doc site; build-from-source + run-tests; License + CONTRIBUTING pointer. Move "How it works" / "Usage" prose to the site (it already lives on `docs/site/concepts/` + `docs/site/reference/command-reference.md`). See the README skeleton below for the exact target shape.
2. **Ensure `docs/install-journey.md` is absent and unreferenced.** On `main` it is already gone (its content is the canonical `docs/site/get-started/install.md`); the live work there is the grep-clean. If work lands on `next` first, the deletion is real there — no content fold needed (the site page already carries it, verified). The durable end-state is: file absent, nothing references it.
3. **Remove `internal/release/install_doc_test.go`** — the banned prose-grep, whichever doc it reads (`main` repointed it to the site page; that did not cure the tautology). The two real behaviors it gestured at are recorded as owed-REAL-tests below, NOT re-added as greps: (a) the Linux `install.sh` path actually works → a behavior/smoke test, not a doc-grep; (b) "spacedock ships no sandbox" honesty → already a code fact (`internal/safehouse` only detects a profile + wraps argv), guarded in code, not by scanning doc prose.
4. **Apply the five behaviors' doc content to the canonical homes** in the retarget map above, after each member's code lands.

### README target skeleton (AC-3's checkable shape)

The slimmed README keeps these sections and drops the rest:

- `# Spacedock` + the one-line what-it-is + "Why?" hook (kept — the front-door pitch).
- `## What's different` (kept — the differentiator, short).
- `## Install` — the `brew tap … && brew install spacedock` block + the one-line `spacedock claude "/spacedock:survey"` launch + the **full-docs link** (see decision below). **Dropped from here:** nothing; this stays the install front door.
- `## License`.
- A `CONTRIBUTING` pointer line (the site's `contributing/` is canonical; README links it).
- **Removed entirely:** `## Quick start` (the commission examples — now `docs/site/running-workflows/commission.md`), `## How it works` (now `docs/site/concepts/operating-model.md`), `## Usage` (the command grammar table — now `docs/site/reference/command-reference.md`).

The removed `## Quick start` / `## How it works` / `## Usage` sections and the presence of the full-docs link are AC-3's checkable change (a structural README diff).

### Open decision for the gate — the "full docs" link target

The assignment names a "prominent spacedock.md link," but no `spacedock.md` file exists anywhere in the repo or its history; it is shorthand for "the canonical doc-site link." Two readings:
- **(A) Link the published doc site URL** (`https://spacedock-dev.github.io/spacedock/`, the mkdocs `site_url`) — the browsable canonical home a README reader expects from "full docs." **Recommended.**
- **(B) Link the in-repo `docs/site/index.md`** path (what `main`'s README does today for install, pointing at `docs/site/get-started/install.md`).

**Recommendation: (A)** as the headline "full docs" link, keeping a secondary in-repo `docs/site/get-started/install.md` path link for readers browsing the repo on GitHub. AC-3 is written to accept either a published-site URL or an in-repo `docs/site/` link as the checkable target; the captain picks the headline at the gate.

## Out of scope

- The doc site's own *structural* redesign (nav, new pages). This task points the README at the site and lands the five behaviors' content on existing site pages; it does not restructure the site.
- The deep **stamp-then-tag** release-ritual rewrite of `docs/site/contributing/releasing.md` — that is `ezn`, deferred out of 0.20.1. yw makes only gp's decouple-accurate edits.
- The behavior **code/tests** themselves — those are the five members' deliverables. yw owns only their *docs*.

## Owed REAL tests (recorded, not re-added as greps)

The deleted `install_doc_test.go` gestured at two real behaviors. Neither is re-added as a doc-grep; each is recorded here for follow-up as a genuine behavioral check:

- **Linux `install.sh` path works.** A smoke test that runs `install.sh` (in a container or with a stubbed download) and asserts the `spacedock` binary lands on `PATH` and `spacedock --version` exits 0 — proves the install path, not that a doc mentions a URL. Follow-up; not blocking this task.
- **No-sandbox honesty.** Already a **code fact**, not a doc claim: `internal/safehouse` only detects a `.safehouse` profile and wraps argv when a `safehouse` binary resolves; it ships no sandbox. The honesty is guarded by the safehouse package's own behavior tests (a profile-present-but-binary-absent launch proceeds unsandboxed), not by scanning doc prose. No new test owed — the code already can't over-claim.

## Acceptance criteria

Each AC names a property of the finished entity, verified outside the task body.

**AC-1 — `docs/install-journey.md` is gone and nothing references it.** Verified by: the file is absent on the target branch (already true on `main`); `git grep -l install-journey -- '*.md' '*.go' '*.y*ml'` (excluding `_archive`/`.spacedock-state`/`.claude`/`docs/roadmap` historical debriefs) returns nothing in live docs/code — command output. (Roadmap debriefs under `docs/roadmap/019x*` are historical records of past sprints and are NOT rewritten; the grep scope excludes them. The grep MUST also come back clean of `install-journey` in `README.md` and in any live `*.go` — on `main` the only live `*.go` hit today is `install_doc_test.go`'s error strings, which go away with AC-2.)

**AC-2 — the prose-grep test is gone and the suite is green.** Verified by: `internal/release/install_doc_test.go` is absent (it is the banned prose-grep regardless of which doc it reads — `main` repointed it to `docs/site/get-started/install.md` without curing the tautology) and `go test ./...` is green from the repo root — exit code. Removing it is clean (verified during ideation): `readInstallJourney` lives only in `install_doc_test.go` (goes with the file), while `executableShellCommands` is defined in `workflow_exec_guard_test.go:441` and used by three other guard tests (`channel_agreement_guard_test.go`, `e2egate_workflow_test.go`, `workflow_exec_guard_test.go`) — it STAYS. So deleting the whole `install_doc_test.go` file leaves no dangling references; `go test ./internal/release/` green is the proof.

**AC-3 — the README is a thin front door linking the site.** Verified by: README contains the install one-liner + a canonical full-docs link (a published-site URL or an in-repo `docs/site/` link) + build/test + License/CONTRIBUTING pointer, and no longer carries the `## Quick start`, `## How it works`, or `## Usage` sections — observed in the rendered README diff (the link target and the three removed section headings are the checkable change). A structural check can assert the three section headings are absent and a `docs/site/` (or `spacedock-dev.github.io`) link is present.

**AC-4 — the five behaviors are documented on their named canonical site pages, and no member ships a `*.md` edit.** Verified by: after the behaviors land, the canonical homes in the retarget map carry the behavior's wording — `docs/site/reference/command-reference.md`'s `--version` example shows the `Sandbox:` line + per-runtime block (gj); its Launch section names the unsandboxed `--permission-mode auto` / codex-equivalent posture (zrc); `docs/site/get-started/install.md` "Keep things in sync" describes the upgrade hint (te), its Homebrew step names the agentsview brew dep (8p), and its Build-from-source / edge note names the `spacedock-edge` channel (gp). The checkable change is the site-page diffs landing under yw's authorship; the proof that NO member edited docs is that each member's PR/commit touches no `*.md` (a reviewer or a `git show --stat` over each member's change shows zero doc files). This AC is the anti-collision contract made checkable: the doc surface moved as one diff, authored once. (Note: AC-4 lands in drive step 3, after the behaviors; AC-1/AC-2/AC-3 land in step 1 and do not wait on the members.)

## Test plan

Mostly deletion + a README rewrite + a site-page doc pass. AC-1/AC-2 are command/exit-code checks (cheap, land first with the structural cleanup). AC-3 is a structural README diff (land first). AC-4 is a set of site-page diffs that land after the five behaviors ship their code — its proof is the diffs plus a per-member `git show --stat` showing zero `*.md` edits. Owed real behavior tests (Linux install path; no-sandbox honesty already code-guarded) are recorded above for follow-up rather than re-added as a grep.

**Sequencing:** this task is split into two waves inside the Commander drive — wave 1 (structural cleanup, AC-1/2/3) leads and depends on nothing; wave 2 (behavior-doc pass, AC-4) follows the five members landing their code. The whole surface targets `main` under the post-flip trunk model. No spike needed: the mechanism is file deletion + Markdown edits over an already-existing canonical site; no parser round-trip, runtime handoff, or on-disk format is at risk. The one verified fact the plan rests on — that `docs/site/get-started/install.md` (on `main`) already carries `install-journey.md`'s content (on `next`) so the deletion loses nothing — was confirmed during ideation by reading both files: same Homebrew/Linux/Codex-Pi/build-from-source/keep-in-sync/command-grammar sections, the site page IS that content. (A direct `git diff` across branches is not possible because each file lives on only one branch; the confirmation is the side-by-side read.)

## Stage Report: ideation

- DONE: Firm yw as the SOLE owner of the 0.20.1 sprint's doc surface — the structural cleanup AND the canonical site/README documentation for the five behavior members' user-visible changes, so those members touch NO docs. Read each member's ideated doc-diff and name the exact canonical home for each (gj, te, gp, zrc, 8p). The five members' own doc-diffs are SUPERSEDED by yw — record that.
  "Anti-collision" section declares yw the single doc writer and SUPERSEDES all five members' doc-diffs; the retarget map table names each member's canonical home (gj → `command-reference.md` `--version` + `install.md` Confirm-it; te → `install.md` Keep-things-in-sync; gp → `install.md` edge note + `contributing/releasing.md`; zrc → `command-reference.md` Launch + `first-launch.md` grammar; 8p → `install.md` Homebrew step), each citing the member's original (now-superseded) diff.
- DONE: Confirm `internal/release/install_doc_test.go` is removed (banned prose-grep); verify nothing else references install-journey.md (grep); record any owed REAL behavior test rather than re-adding a grep.
  AC-2 removes the test; verified the removal is clean (`readInstallJourney` is file-local, `executableShellCommands` is shared from `workflow_exec_guard_test.go:441` and stays). Exercised the AC-1 grep on both branches: `next` live refs = README + the test; `main` live refs = the test's stale error strings only. Owed-REAL-tests section records the Linux `install.sh` smoke (follow-up) and notes the no-sandbox honesty is already code-guarded in `internal/safehouse`, no new test owed.
- DONE: Produce the consolidated doc plan + ACs — external proof only. Note the implementation sequencing inside the Commander drive and the `main` post-flip target.
  Four ACs, all externally checkable: AC-1 (file absent + grep-empty, command output), AC-2 (test absent + `go test ./...` green, exit code), AC-3 (README structural diff — three section headings absent + a `docs/site/` or published-site link present), AC-4 (the five behaviors documented on named site pages + per-member `git show --stat` showing zero `*.md` edits). Two-wave sequencing recorded (wave 1 structural cleanup leads; wave 2 behavior-docs follow the members landing code); whole surface targets `main`.
- DONE: Riskiest-mechanism determination.
  No spike needed: the work is file deletion + Markdown edits over an already-existing canonical site — no parser round-trip, runtime handoff, or on-disk format at risk. The one fact the plan rests on (the site page already carries install-journey's content) was confirmed by a side-by-side read of both branch files.

### Summary
Reframed the task from a single-branch cleanup to the sprint-wide anti-collision doc-owner: yw is now the sole writer of every 0.20.1 user-visible doc change, and all five behavior members' doc-diffs are SUPERSEDED with their content retargeted to canonical `docs/site/` pages (mapped per member). Corrected the entity's branch grounding — `main` and `next` have diverged: `install-journey.md` is already gone on `main` (the site carries its content) and the prose-grep test was *repointed* to the site page on `main` (which did NOT cure the tautology, so its removal is still the live deliverable). Four externally-checkable ACs and a two-wave Commander-drive sequencing (structural cleanup leads, behavior-docs follow the members) recorded; the whole surface targets `main` under the post-flip trunk model. One open gate decision surfaced: the "full docs" link target (no `spacedock.md` exists — recommend the published doc-site URL).
