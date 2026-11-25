# Pre-Merge Whole-Stack Staff Review

**Verdict: HOLD — 1 ship-blocker, 3 non-blocking follow-ups, 1 polish item.**

The blocker concerns required runtime evidence. I found **zero demonstrated production-code blockers** in the assembled stack. The version, artifact publication, and roster mechanisms compose correctly. Configurable paths leave a documentation mismatch in the ensign bootstrap.

`pi-live` rerun **45352323710 remains in progress** at the final CI observation. `paseo-live`, `claude-live sonnet`, `codex-live`, `docs`, and `install-e2e` are green on `03a796c48`. The green Paseo job runs dispatch smoke only. It does not establish adoption on this tip.

## Scope and method

Review date: November 24, 2026, Pacific time. Final CI observation: **2026-11-25 00:40:33 UTC**.

Reviewed the complete 54-file change: **28 non-golden files and 26 golden fixtures**.

| Layer | PR | Reviewed range |
|---|---|---|
| 1 | #895, `x9m` | `3d134b255..4546505f4` |
| 2 | #897, `xdh` | `4546505f4..5210d36d0` |
| 3 | #899, `kdp` | `5210d36d0..03a796c48` |

The checkout HEAD is `03a796c4809c7e926040be9e34b17354ea789994`. Current local `origin/main` is `a85a26a35d550fc38bd46138106179fecf96cedf`. The review follows the roadmap's independent pre-merge review requirement. I read the prior Stack 2 review, its sign-off addendum, runtime support policy, and relevant entity evidence. Findings below reflect an independent source review.

No Git checkout or `.spacedock-state` files were edited. The checkout was clean before and after review. Review files and isolated helper probes are under `/tmp`. I did not start CI, launch workers, create Git fixtures, or change workflow state.

## Ship-blockers (must fix before merge)

1. **Complete the assembled-tip runtime proof before merge.** Owner: Commander, across all three PRs.

   Citations: `docs/dev/README.md:85`; `.github/workflows/paseo-live-e2e.yml:111`; `internal/ensigncycle/paseo_adoption_live_test.go:22`; `skills/first-officer/SKILL.md:35`; `skills/first-officer/SKILL.md:36`.

   The shared dispatch and instruction changes require all affected host lanes. Pi's current result is pending, not passing. The successful Paseo job selects only `^TestLivePaseoDispatchSmoke$`.

   I found recorded adoption success for earlier Layer 1 code. I did not find adoption proof tied to `03a796c48` in the inspected CI and local state records. Layer 3's custom-directory validation proves a worker read its assignment. It does not prove FO restart adoption and final archival.

   **Required action:** Obtain a green Pi rerun. Supply an assembled-tip `TestLivePaseoAdoptionSmoke` result with candidate binary and skill provenance. Retain evidence for one original worker, unchanged epoch, committed report and deliverable, terminal state, archival, and clean trees. Include the prior review's required FO reachability/adapter-selection trace and unreachable-daemon fallback evidence. Existing artifacts can satisfy this requirement if their provenance establishes the candidate.

   **Acceptance:** Link the exact-tip results and relevant traces. Close this evidence item without changing production code if all checks pass. If a check fails, route the specific defect to its owning PR. Pending CI is not itself a code defect.

## Non-blocking follow-ups (seed after merge)

1. **Use the emitted locator throughout the ensign bootstrap contract.** Owner: Layer 3.

   Citations: `skills/ensign/references/ensign-shared-core.md:106`, `:110`, `:114`, `:124`, `:125`, `:132`; `internal/dispatch/build.go:643`, `:872`, `:893`.

   The builder emits the selected absolute directory for every host and for `--advance`. The shared contract still says the first action “MUST be” a read under `/tmp/spacedock-dispatch`. Its advance pattern also names that directory. A private-`/tmp` worker using a shared directory receives inconsistent instructions.

   Replace the fixed path in normative instructions with `{dispatch_file_path}` from the intact pointer. Label any default-directory examples as examples. Preserve failure handling and appended scope instructions.

   **Why non-blocking:** The outer prompt directly supplies the correct path. The recorded Layer 3 live validation read a hashed assignment from `/tmp/spacedock-layer3-shared-dispatch`. See the Layer 3 entity at `index.md:486`, `:487`, and `:510`. The evidence does not support claiming a demonstrated launch failure. The stale MUST and pattern wording still deserve correction. A worker rewriting or ignoring the supplied path promotes this to a material defect.

   **Acceptance:** A worker reads fresh and advance assignments from a non-default directory with spaces. A missing supplied path must still produce `DISPATCH_FILE_MISSING` and stop.

2. **Extend the existing adoption trace check to cover boot selection.** Owner: Layers 1 and 2.

   Citations: `internal/ensigncycle/paseo_adoption_live_test.go:43`, `:67`, `:73`, `:200`, `:218`; `skills/first-officer/SKILL.md:35`, `:36`, `:37`.

   Both FO prompts explicitly request Paseo. The trace predicate accepts an inspect/wait call for W1. It does not establish the FO's reachability probe or reject an inner-host FO adapter load. Ambient-marker selection and sandboxed native fallback remain outside this smoke's assertions.

   Reuse the existing trace reader for an ambient-only case and the probe/load assertions. Exercise unreachable ambient Paseo with a native adapter. Keep committed state as the completion oracle.

   **Why non-blocking:** The current smoke establishes adoption outcomes. Direct candidate traces can satisfy the immediate proof gate. This follow-up prevents later routing regressions from escaping the automated checks.

3. **Document the shared directory's write and hard-link requirements.** Owner: Layer 3.

   Citations: `docs/site/reference/command-reference.md:63`, `:66`; `skills/first-officer/references/paseo-fo-dispatch.md:14`; `internal/dispatch/build.go:125`, `:130`, `:985`, `:1002`, `:1017`.

   The operator notes require shared readability. Publication also requires writer access and hard-link support on the destination filesystem. An existing unwritable directory can pass `MkdirAll`, then fail publication after `--stamp`.

   State those storage requirements in the command reference. Explain that later publication failure does not roll back the stamp. Keep the short adapter note linked or concise.

   **Why non-blocking:** The approved Layer 3 design explicitly accepts this post-stamp boundary. The entity records that disposition at `index.md:496`. Writable, hard-link-capable storage is the supported path. This finding requests clearer documentation, not rollback or another preflight mechanism. Revisit the implementation only if broader filesystem support or transactional stamping becomes promised.

## Polish

1. **Split the remaining changed instructions and comments to meet STE structure rules.** Owners: Layers 1, 2, and 3.

   Citations: `skills/first-officer/SKILL.md:33`, `:35`; `skills/first-officer/references/paseo-fo-dispatch.md:5`, `:142`; `internal/cli/dev_version.go:17`; `internal/cli/dev_version_test.go:57`.

   Representative violations and rewrites:

   - Rule 5.4, condition before command: “Skip this check … when `PASEO_AGENT_ID` is unset and Paseo was not requested.” Rewrite: “If `PASEO_AGENT_ID` is unset and Paseo was not requested, do not run a Paseo command.”
   - Rule 8.1, no semicolons: “load only one transport; never mix native subagent tools with Paseo.” Rewrite: “Load one transport. Do not mix native subagent tools with Paseo.”
   - Rule 6.3, 25-word descriptive limit: the unstamped-build explanation at `internal/cli/dev_version.go:17` spans a long sentence through line 21. Rewrite: “An unstamped source build reports the embedded manifest version with available VCS metadata. Compatibility checks compare the binary and plugin minor versions.”
   - Rules 5.1 and 5.2, short instructions and one action per sentence: the recovery row at `paseo-fo-dispatch.md:142` combines archival, retry construction, stamp selection, and failure routing. Split its actions into short sentences while retaining the current `aK` guard.

   The changed prose under `docs/` and the three new operational notes are concise and structurally sound. The main STE debt is in touched skill paragraphs and comments. This is a structural audit using the supplied skill. Full vocabulary certification requires the official ASD-STE100 dictionary.

## Cross-layer composition and seam audit

### Version identity, session output, and releases

**Assessment: composes cleanly.**

`displayVersion` returns a linker-stamped `Version` before consulting build information (`internal/cli/dev_version.go:24`). Unstamped builds use the embedded manifest, then append available SHA metadata. Missing or invalid revision information retains `+dev`. A valid dirty revision adds `.dirty`.

The display result remains one whitespace-free version token. `printVersion` uses that token independently of runtime selection (`internal/cli/cli.go:901`, `:914`). Native markers retain precedence. `PASEO_AGENT_ID` supplies a display fallback only when native detection is empty and unambiguous. The UUID value is not printed.

The marker-free case still emits three lines. Session cases still emit six lines, with `contract 3` last. The exception at `skills/first-officer/SKILL.md:46` compares the required minor rather than a literal complete version. It therefore accepts `0.29.0-pre0+dev.<sha7>[.dirty]` without weakening wrong-minor or explicit-contract rejection.

`internal/contract/contract.go:74` removes prerelease/build suffixes before parsing major and minor. The focused tests exercise SHA metadata, dirty state, missing metadata, release bypass, native precedence, ambiguity, and token placement. The no-host-command witness now includes `git` and `paseo`.

I also checked the existing clean candidate binary. Its SHA-256 matches the recorded validation hash, and build information identifies `03a796c4809c7e926040be9e34b17354ea789994` with `vcs.modified=false`. This supports the recorded clean-build provenance.

### Immutable publication across all four hosts

**Assessment: the binary path composes cleanly. The bootstrap documentation follow-up remains.**

All hosts assemble `dispatchBody`, then use the same naming and publication functions (`internal/dispatch/build.go:620`, `:639`, `:644`). Claude and Paseo retain their existing parent-token prefixes. Codex and Pi receive the same body-hash protection without requiring a new parent identifier.

The digest changes only the artifact locator. `envelope.name`, `dispatch name`, worker labels, and retry epochs remain unchanged. Consequently, Layer 1 does not mistake a filename hash for worker identity.

The suffix is eight hexadecimal SHA-256 characters, calculated over raw assignment bytes. The final stem includes the suffix before the 251-byte limit check (`internal/dispatch/build.go:991`). A 242-byte old stem succeeds. A 243-byte old stem fails. The `.md` extension makes the accepted maximum 254 bytes. There is no silent truncation.

`--advance` adds its distinguishing stem suffix before hashing (`internal/dispatch/build.go:633`). Its pointer still uses the same published absolute path. Paseo's runtime remains one-shot because its addressability binding is ABSENT. Binary support for an advance envelope does not grant Paseo worker reuse.

Publication creates a temporary file in the destination directory, writes complete bytes, sets mode, closes, then calls `os.Link`. Same-directory creation avoids cross-device linking. An existing regular file succeeds only when bytes match. Different content, a symlink, or a directory fails without replacement. Hash collisions therefore cause an error rather than silent assignment corruption.

The mechanism supplies atomic visibility and no replacement by the builder. It does not make a `0644` file externally tamper-proof or promise power-loss durability. Neither property is required by the accepted design. A lease registry, lock service, or content database would add unnecessary machinery.

Directory selection preserves explicit presence: flag, then environment, then default (`internal/dispatch/dispatch.go:153`; `internal/dispatch/build.go:969`). Empty selected values fail. Invalid unused environment values do not override a valid flag. Relative paths resolve from caller cwd. Both the JSON locator and prompt receive the absolute result. Directory resolution and creation precede stamping. Later publication errors preserve the accepted post-stamp boundary.

The adoption smoke does not reconstruct artifact filenames. It discovers worker identity through labels and inspect results, then observes committed state. Hashing therefore leaves its identity oracle intact. The custom-directory integration fixture proves path emission and fresh-process readability, but `cat` is not an agent bootstrap test (`skills/integration/dispatch_test.go:223`).

### Recovery, ownership, and operational notes

**Assessment: Layer 3 preserves Layer 1's recovery contract.**

The Layer 3 change adds only the shared-directory, workspace-location, and thinking-setting notes at `paseo-fo-dispatch.md:14`, `:15`, and `:16`. It does not alter the matrix, labels, reconciliation filter, or state ownership.

The module retains exactly eight recovery rows at lines 138–145. The paths remain:

| Observation | Result |
|---|---|
| Running worker times out | Continue waiting without retry |
| Pending permission | Hold the entity and preserve the worker |
| One live worker | Adopt it |
| No live worker and valid new report | Archive and advance |
| Failed `a1`, no `a2` | One replacement |
| Failed `a2` | Stop the entity |
| Multiple live workers | Hold for ambiguity |
| Failed create with current attempt absent | Retry `a1` once, stop after `a2` |

The failed-create check compares the exact current `aK` before ordinary recovery (`paseo-fo-dispatch.md:119`). Archived `a1` plus absent failed `a2` reaches Row 8 and stops. A present `a2` from a lost receipt remains eligible for adoption or report verification. Layer 3 does not reopen the prior Row 5 loop.

The three authoritative inputs remain entity/history, code Git state, and daemon labels/inspect. No production lease files or entity-body worker registry are added. Normal FO frontmatter and report operations remain governed by the existing workflow contract. The zero-body-write claim concerns roster bookkeeping, not all workflow activity.

The FO module contains no ensign `references/*.md` citations. The spawn and roster bindings refer to the opening invariant without repeating its complete prohibition. The shared core is **23,464 bytes**, leaving **36 bytes** under its unchanged **23,500-byte cap**.

### Runtime selection and unsupported harnesses

**Assessment: coherent source contract, with the runtime proof boundary stated above.**

A reachable selected Paseo daemon wins FO selection and defers its dispatch module until dispatch/resume. Ambient Paseo failure falls through to the native adapter. Explicit selection with an unreachable daemon stops. Sessions without a marker or explicit selection skip the Paseo probe.

Ensign first-match precedence gives an explicit Paseo pointer priority over native markers (`skills/ensign/SKILL.md:13`). An inherited ambient marker alone does not override a recognized native worker. The unsupported-harness fallback reuses `--host paseo` for assignment shape while binding lifecycle capabilities to actual native features. It does not require a Paseo daemon merely to build an artifact.

Addressability and context budget have separate requirements. A follow-up handle does not imply a usage probe. The host-neutral skill-load edits preserve existing targets and trigger boundaries. They remove concrete `Skill(skill=...)` syntax from the changed shared paths without adding a new loader.

### Current main and scope discipline

A read-only three-way `git merge-tree` check reports **no conflict markers** against `a85a26a35`. Only `docs/site/reference/command-reference.md` changed on both branches. Main adds update-progress documentation in a separate paragraph.

Main's install opt-out and update-progress work does not modify the stack's publication, runtime selection, or version-construction functions. Its version inventory continues to consume the display string. I found no new dependency on an exact `+dev` suffix in those changes. This is source-level composition evidence, not a test run of a merged tree.

The stack adds no dependency, daemon, lease subsystem, or production recovery simulator. The filesystem helpers and test seam are proportional to their purpose. The large fixture diff is mechanical. Remaining prose debt is covered under Polish.

## Necessity inventory

Every non-golden modified file is included. Line references identify the final-tip implementation or test anchor.

| File and anchor | Purpose and necessity assessment |
|---|---|
| `docs/runtime-live-ci.md:105` | Necessary local entry point for dispatch and adoption smoke, plus provider override. Correctly distinguishes dispatch-only CI. |
| `docs/site/reference/command-reference.md:7` | Necessary version and directory documentation. Storage prerequisites need the listed follow-up. |
| `internal/cli/cli.go:914` | Minimal display-only Paseo fallback. Keeps native host detection unchanged. |
| `internal/cli/dev_version.go:24` | Necessary embedded VCS identity with release bypass and metadata fallback. No runtime Git process. |
| `internal/cli/dev_version_test.go:62` | Real source-build proof plus malformed/missing metadata and release tests. Injection avoids dependence on local checkout state. |
| `internal/cli/version_session_test.go:48` | Necessary composed SHA/session/line-count, privacy, precedence, and no-command-execution coverage. |
| `internal/dispatch/build.go:620` | Necessary common publication path, hash naming, directory resolution, and immutable conflict handling. No alternate host-specific writer. |
| `internal/dispatch/build_advance_test.go:117` | Necessary filename expectation update while retaining fresh/advance separation. |
| `internal/dispatch/build_artifact_test.go:28` | Necessary overlap, retry, conflict, directory, pre-stamp, all-host, and boundary coverage. Relative-cwd isolation corrects the prior evidence gap. |
| `internal/dispatch/build_codex_host_test.go:74` | Replaces a stale fixed filename probe with an isolated artifact-directory assertion. Necessary rejection proof. |
| `internal/dispatch/build_errors_test.go:197` | Obtains the actual hashed path before inducing a directory conflict. Preserves the real write-error scenario. |
| `internal/dispatch/build_host_test.go:29` | Proves display fallback does not change native dispatch-host resolution. Necessary seam coverage. |
| `internal/dispatch/build_input_mode_test.go:274` | Checks the whole isolated output directory for rejected input. Necessary after deterministic names change. |
| `internal/dispatch/build_paseo_host_test.go:20` | Reuses the shape oracle for scrubbed markers, hashed names, parent prefixes, and bare rejection. No second harness. |
| `internal/dispatch/build_pi_host_test.go:137` | Necessary hash-aware name expectation while preserving canonical stage assertions. |
| `internal/dispatch/dispatch.go:115` | Necessary option storage, explicit-presence tracking, both flag syntaxes, and help text. |
| `internal/dispatch/golden_harness_test.go:64` | Necessary locator normalization. Raw digest assertion at line 190 protects the behavior that normalization hides. |
| `internal/dispatch/parity_harness_test.go:27` | Necessary environment isolation and removal of the obsolete body-reader implementation. |
| `internal/dispatch/self_contained_assignment_test.go:86` | Necessary isolated no-artifact assertion for unresolved launcher failures. |
| `internal/ensigncycle/paseo_adoption_live_test.go:22` | Necessary restart/adoption outcome proof. Uses real daemon lifecycle and committed Git evidence. The release barrier is test synchronization, not production lease state. |
| `skills/ensign/SKILL.md:12` | Necessary first-match selection for explicit Paseo and inherited-marker fallback. |
| `skills/ensign/references/ensign-shared-core.md:112` | Necessary Paseo bootstrap example. The directory-generalization follow-up applies to this shared owner. |
| `skills/feedback-rejection-flow/SKILL.md:27` | Necessary host-neutral gate load. Existing correction ordering remains intact. |
| `skills/first-officer/SKILL.md:33` | Necessary transport selection owner, reachability decision, generic fallback, and marker-free version exception. |
| `skills/first-officer/references/first-officer-shared-core.md:45` | Necessary host-neutral deferred loads and addressability-aware gate retention. Fits the cap. |
| `skills/first-officer/references/fo-dispatch-core.md:61` | Necessary host-neutral gate invocation and adapter-authorized explicit host at line 188. |
| `skills/first-officer/references/paseo-fo-dispatch.md:107` | Necessary reconciliation and single recovery decision owner. The three operational notes address actual dispatch frictions. |
| `skills/integration/dispatch_test.go:223` | Necessary shared-directory consumer check and preservation of split-root paths. Does not replace live bootstrap evidence. |

All **26 golden files** differ only by adding `-<HASH8>` to artifact locators and their pointer occurrences. Removing that suffix reproduces each old file byte-for-byte. Spot checks included advance, Codex, Pi, HTML escaping, spaces, and non-ASCII titles. Bodies remain unchanged. `readDispatchBody` checks the raw file's digest before normalization. The normalizer operates on the decoded locator and its encoded occurrences, not arbitrary assignment text.

## Validation and evidence record

Independent checks executed successfully:

1. Focused race tests across `internal/cli`, `internal/dispatch`, `internal/runtimehost`, and `internal/contract`. These cover source-build metadata, session output, release bypass, marker precedence, dispatch-host resolution, and minor parsing.
2. Contract lint with `-race`, excluding `TestBreakGlassCommitLandsInEntityRepo` because it creates and changes Git fixtures.
3. `go vet ./...`.
4. Live-tag compilation through `go test -tags live -run '^$' ./internal/ensigncycle`. This is compilation evidence, not a live pass.
5. An isolated `/tmp` test of verbatim extracted publication functions. Thirty-two concurrent identical publishers succeed without replacement or leaked temporary files. Conflicting bytes, symlinks, and directories fail. A known SHA prefix and the 242/243-byte boundary pass. This exercises the filesystem helpers, not the full CLI.
6. Read-only formatting, diff whitespace, golden comparison, shared-core byte count, and three-way merge inspection.

Changed Go files are formatted. Whole-tree `gofmt -l` reports `internal/release/runtime_live_evidence_workflow_test.go`. The same formatting diff exists on `a85a26a35`, so it is baseline debt rather than a stack regression.

I did not run local full normal/race suites or live fixtures because they write Git checkouts. I did not run `gofmt -w`. Those omissions obey the explicit read-only request.

The exact-tip CI offline job passed its full normal suite. Layer 3's entity records full normal/race package coverage with a documented PATH split at `index.md:471` and `:472`. Existing local normal/race logs contain passing package results. I treated those logs as recorded validation, not as an independently repeated full-suite run.

| Exact-tip CI evidence | Result |
|---|---|
| [Runtime run 45352320212](https://github.com/spacedock-dev/spacedock/actions/runs/45352320212) | Offline, Codex, and Claude sonnet passed. Pi skipped in this run. |
| [Pi rerun 45352323710](https://github.com/spacedock-dev/spacedock/actions/runs/45352323710) | Pi in progress. Offline passed. |
| [Paseo run 45352327681](https://github.com/spacedock-dev/spacedock/actions/runs/45352327681) | Dispatch smoke passed. Adoption is not selected. |
| [Docs run 45352212924](https://github.com/spacedock-dev/spacedock/actions/runs/45352212924) | Passed. |
| [Install run 45352212863](https://github.com/spacedock-dev/spacedock/actions/runs/45352212863) | Passed. |

State citations for the Layer 3 entity refer to `/Users/clkao/git/spacedock-research/spacedock-v1/docs/dev/.spacedock-state/paseo-vm-dispatch-and-version-frictions/index.md`. Other relative citations refer to the reviewed checkout at `03a796c48`.

Review tracking is complete: scope, file inventory, composition audit, permitted checks, evidence review, classification, and report. No findings were silently dropped. No follow-up entities were created under the read-only restriction.
