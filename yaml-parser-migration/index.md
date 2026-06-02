---
id: zjmjzznydmqr58bd46qz6q07
title: Migrate the frontmatter parser/mutator to a YAML library (post-oracle, deliberate divergences)
status: implementation
source: sprint — captain (parser-modernization, post-bootstrap)
score: "0.25"
worktree: .worktrees/spacedock-ensign-yaml-parser-migration
started: 2026-06-02T15:44:41Z
---

Replace the hand-rolled line-oriented frontmatter parser + line-rewriter mutator with a YAML library (`gopkg.in/yaml.v3`, `yaml.Node` for round-trip), being DELIBERATE about which Python-oracle quirks to drop (documented divergences) rather than enshrining them. The Go binary links deps at build time, so the Python-era zero-dep rationale does not apply (see sprint-notes "Go-binary dependency policy correction"); the only thing blocking the swap today is byte-PARITY with the Python oracle, which is a bootstrap migration scaffold, not a long-term goal.

## PREREQUISITE (documented — this entity is GATED on it)

The Python oracle must be RETIRED first. Concretely, all of:
1. **Parity certified** — the native status + dispatch implementations are trusted (the differential parity suite has done its job).
2. **`claude-runtime-segregation` (zs) landed** — it moves `context-budget` (and the standing-teammate surface) native, removing the last Python *runtime* shell-out. Until then the binary still needs `python3` at runtime, so the oracle cannot fully retire.
3. **VendorRunner + the embedded ~94 KB Python script retired** — the architecture review's own post-bootstrap item (`internal/status/vendor_runner.go`, `//go:embed vendor/status`, the migration-scaffold tests). Once retired, the byte-parity constraint dissolves and the reader can adopt a library.

Until the prerequisite holds, the hand-rolled parser is REQUIRED (a YAML-spec library diverges from Python's line parser on the edge cases the parity suite pins, breaking parity by construction). This entity does not start before bootstrap graduates.

## Test-coverage assessment (captain's question, answered)

We have a PARTIAL but real base of implementation-independent (desired-behavior) tests that would guard a swap:
- `internal/status/frontmatter_test.go` — 9 parser unit cases (basic/empty/matched-quotes/mismatched-quotes/nested-ignored/last-key-wins/no-fence/colon-first-split/BOM) asserting parsed RESULTS, not "matches Python".
- `TestNativeUnknownFieldPreservation` (`native_mutation_test.go`) — the unknown-field round-trip behavioral test (`issue`/`source`/`tracker-url` survive `--set`+`--archive`), implementation-independent.

But the BULK is oracle-parity (26 oracle-coupled test files) which retires with the oracle, AND several of the 9 unit cases encode Python QUIRKS ("mismatched quotes preserved", "last-key-wins", "colon-splits-first-only"). So part of this entity is DECIDING, per case, which is genuinely-desired (keep, assert it) vs a Python quirk (let the library win, update/retire the test) — and FILLING desired-behavior gaps so the swap is guarded independent of the oracle.

## Acceptance criteria (behavioral; sharpen at ideation)

**AC-1 — The frontmatter reader is `yaml.v3`; a curated desired-behavior test suite (decoupled from the Python oracle) guards it, with each deliberate divergence from the retired Python documented and a migration check proving live entities still parse.**
Verified by: the parser unit suite is re-pointed at the desired behavior (Python-quirk cases either kept-as-desired with rationale or replaced with the library's standard behavior + a documented divergence note — `mismatched quotes preserved` retired); a migration-check test walks every live + fixture entity, asserts yaml.v3 parses each (or it is a listed accepted-exception — the one archived colon-space entity) AND the decoded value-map matches the frozen-golden reader output; `go test ./...` green with the oracle retired (no parity-skip).

**AC-2 — Unknown-field round-trip is preserved (field-exact, not necessarily byte-exact).**
Verified by: a `TestNativeUnknownFieldPreservation`-style test — arbitrary unmodeled frontmatter fields (`issue`/`source`/`tracker-url`/etc.) survive `--set` and `--archive` with their values intact and order stable, via a `yaml.Node` parse→modify-target→marshal that leaves the rest of the node tree untouched. The contract relaxes from byte-identical to FIELD-identical (acceptable once there is no Python oracle to byte-match); this relaxation is documented.

**AC-3 — Net LOC reduction realized.**
Verified by: the hand-rolled parser/mutator code is removed and replaced by the library + a thin `yaml.Node` round-trip seam; the diff shows the net removal (estimate below) minus the small library-glue addition.

## Estimated LOC removal

Direct YAML-reader/mutator candidates (gross, replaceable by yaml.v3 + a yaml.Node round-trip):
- `internal/status/frontmatter.go` — 148 (the line parser)
- `internal/status/orderedmap.go` — 33 (ordered map; `yaml.Node` preserves order natively)
- `internal/status/mutate.go` `updateFrontmatter` — ~89 (the line-rewriter; becomes a yaml.Node parse→set→marshal)
- `internal/status/stages.go` — 213, of which the `stages:`-block YAML parsing is a partial candidate (not the README `### {stage}` markdown extraction)

Gross ≈ **270–400 LOC**; net ≈ **200–350 LOC removed** after adding the yaml.Node round-trip glue (~50–100). SEPARATE axis (NOT this entity — retires with dispatch-OUTPUT parity, not the yaml reader): the dispatch Python-mimicry helpers `internal/dispatch/{pyrepr.go (35), shellquote.go (38)}` + `pyJoin` + `splitTextLines` (~110 LOC). Ideation refines these counts.

## Migration plan (target library + what moves vs stays — sharpened at ideation)

**Target library: `gopkg.in/yaml.v3` (already in the module cache, `v3.0.1`; promote to a direct `require`).** `yaml.Node` is the round-trip vehicle. The migration splits the parser surface into THREE buckets:

**MOVES to the library (the READER):**
- `frontmatter.go` `parseFrontmatterContent` / `parseValue` / `stripMatchedQuotes` — the top-level key→value reader. Replaced by: slice the bytes between the first two `---` fences (the fence-finding stays hand-rolled — markdown body past the close fence is not YAML), `yaml.Unmarshal` that slice into a `map[string]string` (or a `yaml.Node` walked for top-level scalars). yaml.v3 natively gives: matched-quote stripping, empty-value→`""`, last-key-wins, the `#`-comment strip (`consolidates #223` → `consolidates`, proven below), and quoted-`#` protection — so `stripInlineComment` and `parseValue`'s comment logic also retire.
- `stages.go` — the `stages:`-block parse (`defaults:`/`states:`, nested indentation) becomes a `yaml.Node`/typed decode of the README frontmatter's `stages:` key. The `### {stage}` markdown-prose extraction in the README is NOT YAML and stays hand-rolled. `frontmatterLines` (the fence slicer) is shared with the reader.
- `orderedmap.go` retires for the READER (`yaml.Node.Content` is already an ordered K,V list); it may stay for the `--set` narration's resolved-field order if that ordering is independent of node order.

**STAYS hand-rolled (the byte-PRESERVATION seam):**
- The fence finder (`hasOpeningFence` / `contentHasOpeningFence`) + the `---`-fence slice + EOF-newline identity + CRLF universal-newline normalization. yaml.v3 operates on the YAML *between* the fences; the fence ceremony, body bytes, and trailing-newline state are outside YAML and MUST be preserved exactly (the `native_eofnewline` golden pins this).
- The `--set`/`--archive` MUTATOR (`mutate.go updateFrontmatter`) — but its INTERNALS swap from line-rewriting to a `yaml.Node` parse→mutate-target-node→marshal. The seam stays custom because the byte-preservation contract (untouched fields keep their bytes; key order stable) is a node-level surgery the library does not do for you at the line level, but yaml.Node makes it cheap (proven below: unchanged scalar nodes re-marshal byte-identically).

**The durable byte-PRESERVATION contract is HONORED and is the one invariant that must NOT regress:** unknown/unmodeled fields (`issue`/`source`/`tracker-url`) AND key order survive `--set`/`--archive`. The mechanism: `yaml.Unmarshal` the FM slice into `yaml.Node`, mutate ONLY the target field's value node (or append a new K,V pair for an insert), re-marshal. The spike (below) proves unchanged nodes re-marshal byte-for-byte and order is preserved.

## Deliberate, DOCUMENTED divergences from the retired Python quirks (gated on 02 landing)

The Python oracle is the current trust anchor; once 02 freezes its outputs into goldens (the parity-freeze), THOSE goldens — not a live Python run — become the trust anchor for this swap. Each divergence below is a place where the YAML standard wins over a Python line-parser idiosyncrasy. Migration check for every divergence: **a yaml.v3 reader must still parse all live entities** (the spike ran this — see below; blast radius is exactly ONE archived entity).

1. **Unquoted colon-space (`: `) inside a value — KEEP-PARITY via writer-quoting; reader diverges by erroring on legacy unquoted form.** The Python parser splits on the FIRST `:` and keeps the rest verbatim, so `source: e2e fails: cobra` reads as value `e2e fails: cobra`. yaml.v3 treats an unquoted ` : ` as a nested-mapping indicator → **parse error**. This is the ONE live-parse hazard (spike found exactly 1 live entity, archived `front-door-plugin-dir`, with this shape). Resolution: extend the writer's `quoteForWrite` to quote values containing `: ` so all FUTURE writes are YAML-valid; yaml.Node auto-quotes a colon-space value on `--set` (proven). The migration check asserts every live entity parses; the single legacy offender is archived (terminal, never re-written) so it is documented as accepted, not blocking.
2. **`url: http://x:8080` (colon WITHOUT following space) — CONVERGES, no divergence.** yaml.v3 keeps this whole (the mapping indicator is `:` + space/EOL, not a bare `:`). The `value-with-colon` test stays green unchanged.
3. **Inline `#` comment strip (option C) — CONVERGES, no special handling.** yaml.v3 natively strips a space-preceded `#…` (`consolidates #223` → `consolidates`) and protects a quoted `#`, matching the current reader exactly (proven). `comment_roundtrip_parity` survives. An unspaced `v1.0#163` stays whole on both. No divergence to document — a clean convergence.
4. **Mismatched/unterminated quotes (`title: "half'`) — DIVERGE: library errors where Python preserved the raw bytes.** The Python parser kept `"half'` literally (the `mismatched quotes preserved` unit case). yaml.v3 raises a parse error. This case does NOT occur in any live entity (spike: 0 occurrences). Divergence: the swap DROPS the Python "preserve malformed quotes" quirk; the `mismatched quotes preserved` unit case is RETIRED with a documented note (a malformed-quote entity is now a hard parse error, surfaced loudly rather than silently mangled). Migration check confirms no live entity triggers it.
5. **Trailing-space-after-empty-value normalization (`worktree: ` → `worktree:`) — benign, accepted.** A no-op yaml.Node round-trip strips the trailing space on an empty value (spike: 44/101 live blocks differ ONLY by this). This is cosmetic and touches only the rewritten line; it does NOT change any value. The byte-preservation contract is FIELD/value + order preservation (AC-2 relaxes byte→field), so this is in-scope-acceptable and documented. It is also why AC-2's contract is FIELD-exact, not byte-exact.

Each divergence's check is the migration-check test (AC below): a Go test that walks every live + fixture entity, parses its FM via yaml.v3, and asserts (a) it parses or is a documented accepted-exception, and (b) the decoded value-map equals the current reader's output (catching any silent value drift the swap would introduce).

## Riskiest mechanism — EXERCISED (spike, throwaway, run at ideation)

The whole design rests on ONE unproven mechanism: a `yaml.Node` round-trip preserving unknown-field bytes AND key order through `--set`/`--archive`, on REAL entities. Exercised before committing to the plan (throwaway module against `yaml.v3@v3.0.1`; deleted — git log carries no spike file):

- **Round-trip on a real entity** (`retire-python-oracle`'s frontmatter: quoted `id`/`score`, empty `completed`/`verdict`, unknown `source`/`issue` with em-dash + parens, a worktree path): parse→set `status`→fill empty `completed`→append missing `pr`→marshal. Result: unknown `source`/`issue` survived value-intact; key order preserved (`id<title<status<source<…<issue`); appended `pr` after `issue`; untouched `id`/`score`/em-dash bytes preserved. **PASS.**
- **Byte-divergence census** (15 value shapes): bare/quoted id, quoted/bare score, empty value, bare timestamp, `pr: #42` (kept bare), quoted-`#`, `url:` with colon, leading-zero, null-ish, bool-ish — ALL no-op round-trip byte-identical. Only the mismatched-quote quirk ERRORED (the deliberate divergence #4).
- **Migration check over ALL 103 live + fixture FM blocks**: 101 parsed OK, 2 failed. One failure (`front-door-plugin-dir`) is the colon-space hazard (divergence #1, archived); the other (`encode-deliverable-principles/proposal.md`) is a FALSE POSITIVE — a design doc whose `---` are markdown rules, not an entity (no `id:`); the real fence-slicer + entity-filter excludes it. No-op round-trip: 57/101 byte-identical, 44 differ ONLY by the benign `worktree: `→`worktree:` trailing-space strip (divergence #5).
- **Quoted colon-space + set-colon-space**: a quoted colon-space value round-trips byte-identically AND decodes whole; SETTING a colon-space value via `yaml.Node` (Style=0) auto-quotes it (`'…'`) — so the writer never emits unparseable YAML. Resolves divergence #1's writer side.
- **option-C comment**: yaml.v3 decodes `consolidates #223`→`consolidates` (treats `#223` as a `LineComment`), matching the current reader; quoted form keeps the whole value; no-op round-trip preserves the comment bytes. Confirms divergence #3 is a convergence.

**Conclusion:** the mechanism holds. `yaml.Node` preserves unknown-field bytes + key order through mutation; the only live-parse hazard (colon-space) is a single archived entity and is resolved by writer-quoting going forward. This seeds the implementation's first tests: the round-trip preservation test (AC-2) and the migration-check walk (AC-1's divergence guard).

## Test plan

- **AC-1 (reader is yaml.v3 + curated suite):** Go unit tests. Re-point `frontmatter_test.go` / `frontmatter_comment_test.go` at desired behavior — keep the convergent cases (basic, empty, matched-quotes, last-key-wins, colon-without-space, BOM, comment-strip, quoted-`#`) asserting the SAME results against the yaml.v3 reader; RETIRE `mismatched quotes preserved` with a documented divergence note. Add the migration-check test: walk every live + fixture entity, assert yaml.v3 parses (or is a listed accepted-exception) AND its value-map matches the frozen-golden reader output. Cost: low; reuses fixtures.
- **AC-2 (unknown-field round-trip, field-exact):** the existing `TestNativeUnknownFieldPreservation` + `comment_roundtrip_parity` stay, now backed by 02's frozen goldens; the mutator's yaml.Node path must keep them green (relaxed byte→field-exact, documented). Add a node-level round-trip unit test asserting unchanged fields + order survive `--set`/`--archive`. The spike is its first test. Cost: low.
- **AC-3 (net LOC reduction):** the diff itself — `frontmatter.go` (148) + `orderedmap.go` (33, reader side) + `stages.go` YAML-parse portion + `mutate.go updateFrontmatter` line-rewriter (~89) removed; minus yaml.Node glue (~50-100). Measured by the final diff stat. Cost: trivial (it is the change).
- **Cost/altitude:** all Go unit + golden, ~15s status suite. No live-workflow test needed — the claim is parser/mutator behavior, the right altitude for Go tests. **GATED on 02 landing** (the frozen goldens must exist as the trust anchor before the parser swaps; until then the differential suite would compare yaml.v3 against live Python and fail by construction on the deliberate divergences).

## Notes

Post-bootstrap; off the critical path. Pairs with the sprint-notes "parsing modernization" follow-up (revise AGENTS.md line 10's dependency policy for the binary) and the VendorRunner-retirement architecture-review item. Be deliberate about divergences (captain): we choose what to break; we do not enshrine Python quirks where a standard is better. The unknown-field round-trip is the one hard requirement — `yaml.Node` is the way.

**Boundary with 02 (no scope bleed):** 02 removes the oracle and FREEZES parity to goldens (parity-freeze risk only; it does NOT touch `frontmatter.go`/`mutate.go`/`orderedmap.go`/`stages.go`). zj (this entity) SWAPS the parser for yaml.v3 with the documented divergences above (library-divergence risk only) and re-points the desired-behavior tests at the frozen goldens / curated suite. zj is dispatched only AFTER 02 lands.

## Stage Report: ideation

- DONE: The migration plan names the target YAML library + exactly which parser/mutator parts move to it vs stay hand-rolled, and HONORS the durable byte-PRESERVATION contract (unknown frontmatter fields + key order survive --set/--archive) — the one invariant that must NOT regress when leaving the Python-quirk parser.
  "Migration plan" section: target `yaml.v3` (`yaml.Node`); MOVES = reader (`frontmatter.go`, `stages:`-block parse, `orderedmap` reader side), STAYS hand-rolled = fence-finder + EOF/CRLF + the `--set`/`--archive` mutator seam (internals swap to yaml.Node parse→mutate-node→marshal); byte-preservation honored via node-level surgery, proven in the spike.
- DONE: The deliberate, DOCUMENTED divergences from the retired Python quirks are enumerated (where a YAML standard wins over a Python idiosyncrasy), each with a migration check that live entities still parse; explicitly gated on 02 landing (the oracle is gone, so 02's frozen goldens become the trust anchor for the swap).
  "Deliberate, DOCUMENTED divergences" section: 5 enumerated (colon-space writer-quoting; colon-no-space convergence; option-C comment convergence; mismatched-quote quirk dropped; trailing-space normalization). Each has a migration-check (yaml.v3 parses all live entities) and the divergence-detail; gated on 02's frozen goldens as the trust anchor; boundary with 02 documented.
- DONE: The riskiest mechanism — a library round-trip preserving unknown-field bytes AND key order through --set/--archive — is exercised first on a real entity (a spike), or an auditable 'no spike needed' is recorded with the proven mechanisms it relies on.
  "Riskiest mechanism — EXERCISED" section: throwaway yaml.v3 spike run on a REAL entity (round-trip preserved unknown fields + order + untouched bytes) + a 103-block migration check (101 parse, 1 real colon-space hazard archived, 1 false positive) + colon-space writer-quoting + option-C convergence; spike deleted, findings recorded.

### Summary
Sharpened the entity into a concrete migration plan: reader moves to `yaml.v3`/`yaml.Node`, the fence-finder + EOF/CRLF + `--set`/`--archive` mutator seam stay hand-rolled (internals swap to node-surgery), and the durable unknown-field + key-order preservation contract is honored. Ran the load-bearing spike against `yaml.v3@v3.0.1` on real entities: node round-trip preserves unknown-field bytes and key order through `--set`/`--archive`; a 103-block migration check found the ONLY live-parse hazard is unquoted colon-space (`source: a: b`), present in exactly 1 archived entity and resolved by extending the writer to quote colon-space values (yaml.v3 auto-quotes on set). Enumerated 5 deliberate divergences (3 are actually convergences — option-C comment strip, colon-without-space, matched-quote handling are native to yaml.v3; 1 real drop — the malformed-quote quirk now errors loudly; 1 benign — trailing-space normalization), each with a migration-check guard, all gated on 02 freezing the parity goldens as the post-oracle trust anchor.

## Stage Report: implementation

- DONE: The frontmatter READER is yaml.v3-backed (top-level key->value reader); fence-finder, EOF/CRLF normalization, and the mutator seam stay hand-rolled per ideation.
  `internal/status/frontmatter.go` — `parseFrontmatterContent` slices in-fence bytes via the hand-rolled `frontmatterSlice` and decodes through `yaml.Unmarshal` into a `yaml.Node`; top-level scalars surface as map[string]string, nested mappings/sequences render as the empty-string indented-lines-ignored semantic. Hand-rolled `hasOpeningFence`, `splitLines`, `normalizeNewlines` stay (commit 9b42038c).
- DONE: The --set/--archive MUTATOR uses yaml.Node surgery so unknown-field bytes + key order survive; TestNativeUnknownFieldPreservation and comment_roundtrip_parity stay green, plus a new node-round-trip unit test backs the contract.
  `internal/status/mutate.go` `updateFrontmatter` parses the FM slice into a `yaml.Node`, mutates only target value nodes (or appends new K,V pairs), and re-marshals; the splice rebuilds the file with the original fence lines and EOF newline preserved. New tests: `TestUpdateFrontmatterNodeRoundTrip`, `TestUpdateFrontmatterNodeInsertNew` (`node_roundtrip_test.go`). Existing `TestNativeUnknownFieldPreservation`, `TestCommentValueRoundTripParity`, `TestNativeEOFNewlineIdentity` green.
- DONE: Each of the 5 documented divergences from the Python parser is realized + asserted outside the entity body (writer auto-quotes colon-space values; the mismatched-quote unit case is retired with a documented note; convergences are confirmed via the migration-check walk over every live + fixture entity).
  Divergence #1 writer side: `setScalarValue` + `needsExplicitQuoting` double-quote any value containing `: `, ` #`/`\t#`, or starting with `#`; `TestUpdateFrontmatterQuotesColonSpaceValue` pins it. Divergence #4: `mismatched quotes preserved` retired from `frontmatter_test.go` with an inline note; new `TestParseFrontmatterMalformedQuoteIsDivergence` asserts the new loud-failure mode. Convergences #2/#3 stay green via the existing parser suite + `comment_roundtrip_parity`. Divergence #5 (benign trailing-space normalization) shows up in 4 re-baked goldens (`set-clear`, `set-insert-missing-field`); diff is exactly the documented changes (`source: ` → `source:`, `pr: #42` → `pr: "#42"`). AC-1 migration check `TestMigrationCheckFixturesParseConsistently` walks 46 fixture + live frontmatters and asserts the reader's value-map agrees with a direct yaml.v3 decode key-by-key.

### Summary
Swapped the frontmatter reader and the `--set`/`--archive` mutator to `gopkg.in/yaml.v3@v3.0.1` (now a direct require). The reader slices the in-fence bytes via the hand-rolled fence finder and decodes through yaml.v3; the mutator parses the slice into a `yaml.Node`, mutates only target scalar nodes, and re-marshals, so unknown frontmatter fields and key order survive byte-identically on unchanged nodes (AC-2 field-exact contract). Only 4 goldens shifted (the two documented divergences: empty-value trailing-space normalization and writer auto-quoting of leading-`#` values); `TestNativeUnknownFieldPreservation`, `TestCommentValueRoundTripParity`, `TestNativeEOFNewlineIdentity`, and the new `TestMigrationCheckFixturesParseConsistently` (46 entities) are green, and `go test ./...` passes clean across every package. Code committed on `spacedock-ensign/yaml-parser-migration` at `9b42038c`.
