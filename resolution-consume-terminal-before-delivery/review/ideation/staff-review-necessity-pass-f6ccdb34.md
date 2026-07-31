CL,

## Review

### 1. Fidelity — endorse with revisions

- **Correct:** The proposed mechanism faithfully states all six required semantics: pending-preserving terminal consume, uniform routing across all delivery classes, merge-guard-only spend, supersede-and-feedback failure routing, mutation-free retry, and no new state/lock/scanner machinery (`resolution-consume-terminal-before-delivery.md:84-93`).
- **Correct:** The sole admitted addition is narrowly scoped to non-forced terminal status updates while a terminal-target application is pending; `--force` remains available (`...md:99-101`).
- **Material finding:** The atomic-envelope implementation is described as more proven than it is. `writeDocumentAndStatus` only replaces `gates` and `status` (`internal/gates/application.go:281-283`; `internal/gates/io.go:188-226`). Current merge finalization separately writes `status`, `verdict`, and `completed` through `runSet` (`internal/status/merge.go:392-399,702-709`). Therefore the four-field application/status/verdict/completed writer in `...md:90,199` is a necessary new composition, not the existing helper unchanged. Revise the mechanism claim and explicitly preserve the current merge-guard status validations when introducing that transaction. This does not justify a spike or broader architecture.
- **Note:** Expected-surface arithmetic omits the CLI row: 80+120+25+60 mechanism/plumbing LOC is 285, and 285+300+155 is approximately 740, not 680 (`...md:174-186`). The semantic core remains 225 only if CLI wiring is explicitly excluded.

### 2. Necessity of the status refusal — endorse

The refusal is transaction-necessary, not a ratchet click. Today a bare terminal status update is deliberately allowed (`internal/status/handlers.go:210-218`), and `merge: local` exempts it from the merge-proof guard (`internal/status/handlers.go:253-261`). The shipping test proves such an update succeeds without a sentinel (`internal/status/merge_policy_guard_test.go:30-38`).

Without the proposed refusal, a hookless/manual-local entity with a pending terminal application could enter terminal status outside merge guard. The proposed predicate is minimal: non-forced, terminal-target pending application, terminal status target only (`...md:101`).

### 3. Truth of proven mechanisms — endorse with correction

1. **Application + status co-write exists:** `internal/gates/application.go:176-181,281-283`; the replacement and atomic rename are at `internal/gates/io.go:188-226,384-410`.
2. **Guarded `pending→superseded` exists:** `internal/gates/application.go:157-165`, guarded byte-exactly by `internal/gates/application.go:245-278`.
3. **Revise reads `feedback-to`:** `internal/gates/operation.go:486-510`.
   - **Material finding:** The claim that existing revise routing applies “no silent fallback” is false (`...md:155,199`). Missing `feedback-to` currently falls back to the same stage (`internal/gates/operation.go:505-510`), while parsing does not validate that the target is defined or nonterminal (`internal/gates/operation.go:552-570`). The proposed `--failed` refusals are necessary new validation around a proven lookup, not unchanged shipping semantics.
4. **Merge guard phases exist:** armed/blocked/finalize classification is at `internal/status/merge.go:121-175`; finalization is at `internal/status/merge.go:374-423`. Existing tests cover both policies (`internal/status/merge_guard_test.go:92-133`).

### 4. Accepted prepare residual — accept

“Fail-closed at merge guard” is honest, although “shadow” understates the mutation. Preparing a successor explicitly changes the old pending application to `superseded` and appends an open attempt (`internal/gates/prepare.go:251-261`). Latest-attempt readiness then sees the open successor (`internal/gates/model.go:174-180`), so a merge guard requiring a current pending terminal application must refuse without spending authority or terminalizing.

A mistaken prepare can delay recognition of an already-landed delivery until a fresh approval, but it cannot spend the old authority or erase the external delivery. Normal routing does not prepare: a pending terminal approval projects as `approved-awaiting-merge` (`internal/gates/model.go:214-221`; `internal/status/boot_identify_test.go:321-334`).

### 5. AC sufficiency at the value — revisions required

- **Correct:** AC-1 can fail against the old behavior. Current consume unconditionally changes an eligible application to consumed and advances status (`internal/gates/application.go:176-185`); the existing boot fixture confirms the terminal pending gate disappears after consume (`internal/status/boot_identify_test.go:329-338`). The old CLI also lacks `merge guard --failed`.
- **Material finding:** AC-1 says recovery needs “no second gate cycle,” but its fixture explicitly creates attempt-2 and obtains a fresh approval (`...md:145-146`). Replace that claim with “no ad hoc replacement gate or surgery; normal successor attempt after rework.”
- **Material finding:** AC-1’s final-state assertions cannot prove “all four in one replacement.” A split application write followed by a status/verdict/completed write reaches the same final file and passes the described CLI assertions (`...md:146`). Add a narrow package-level replacement/failure-injection assertion that observes one candidate replacement containing all four changes and proves replacement failure leaves the original bytes intact. This directly measures the required atomic envelope without adding production machinery.
- **Material finding:** AC-2 names manual local merge but does not require the fixture to prove the manual merge occurred before merge guard (`...md:148-149`). Current behavior finalizes first and only then instructs the operator to merge (`internal/status/merge_guard_test.go:631-644`), so the old ordering could survive a final-state-only matrix. The manual leg should record/assert the merge commit before invoking merge guard.

## OVERALL: endorse-with-revisions

**Single strongest reason:** the design mechanism is sound and minimal, but its central atomic-envelope claim currently relies on a helper that writes only application/status and on a CLI fixture that cannot distinguish one replacement from split writes.

No blocker requires rejecting or broadening the mechanism.