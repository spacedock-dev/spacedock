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

Now that the doc site (`docs/site/`, shipped via #343) is the canonical browsable documentation, the repo carries redundant docs. Establish the rule **site is canonical; the GitHub README is a thin front door**, and remove the duplication.

## Problem

- `docs/install-journey.md` (168 lines: Homebrew / Linux / Codex-Pi / build-from-source / keep-in-sync / command-grammar) duplicates the site's get-started + reference pages.
- `README.md` (144 lines) carries "How it works" / "Usage" prose that duplicates the site's concepts/running-workflows.
- `internal/release/install_doc_test.go` is a **banned prose-grep**: both tests `os.ReadFile` `install-journey.md` and assert substrings the author wrote — a `curl … install.sh | sh` line (`TestInstallJourneyDocumentsLinuxPath`) and the absence of "sandboxed on linux" phrasings (`TestInstallJourneyDoesNotOverclaimLinuxSandbox`). The expected values come from the file under test, not an independent source; a valid paraphrase fails them and an inverted clause passes them. This is the instruction-file-read / prose-grep the workflow quarantine forbids outside `internal/contractlint`.

## Proposed approach

1. **Slim README** to the front-door essentials a repo visitor/contributor needs without leaving: one-line what-it-is + the "what's different" hook; the single install path (`brew install …` → `spacedock claude`) then a prominent **"full docs → spacedock.md"** link; build-from-source + run-tests; License + CONTRIBUTING pointer. Move "How it works" / "Usage" prose to a site link.
2. **Remove `docs/install-journey.md`**; fold anything not already on the site into the site's get-started pages.
3. **Remove `internal/release/install_doc_test.go`** with the doc. If a real behavior it gestured at still matters, record the owed REAL test rather than keeping the grep: (a) the Linux `install.sh` path actually works → a behavior/smoke test, not a doc-grep; (b) "spacedock ships no sandbox" honesty → already a code fact (`internal/safehouse` only detects a profile + wraps argv), guarded in code, not by scanning doc prose.

## Out of scope

- The doc site's own content edits (it is canonical; this task points at it, does not rewrite it).
- 0.20.1's behavior-doc updates (sandbox in `--version`, install-refresh, edge channel) — those RETARGET to the site/README under this task's rule, but the edits themselves belong to their 0.20.1 members.

## Acceptance criteria

Each AC names a property of the finished entity, verified outside the task body.

**AC-1 — `docs/install-journey.md` is gone and nothing references it.** Verified by: the file is absent; `grep -rl install-journey` over tracked `*.md`/`*.go`/`*.y*ml` (excluding `_archive`/`.spacedock-state`/`.claude`) returns nothing — command output.

**AC-2 — the prose-grep test is gone and the suite is green.** Verified by: `internal/release/install_doc_test.go` is absent and `go test ./...` is green from the repo root — exit code.

**AC-3 — the README is a thin front door linking the site.** Verified by: README contains the install one-liner + a `spacedock.md` link + build/test + License/CONTRIBUTING, and no longer carries the full install journey or duplicated usage prose — observed in the rendered README diff (the link target and the removed sections are the checkable change).

## Test plan

Mostly deletion + a README rewrite. AC-1/AC-2 are command/exit-code checks (cheap). AC-3 is a structural diff. Any owed real behavior test (Linux install path; no-sandbox honesty) is recorded here for follow-up rather than re-added as a grep. Note: sequences AHEAD of 0.20.1's doc edits (it removes the file five 0.20.1 members would otherwise edit) and likely targets `main` under the post-flip trunk model.
