## Review

**Merge verdict: BLOCK pending reconciliation and behavioral evidence.** The paragraph is a reasonable scoped rule, but it does not establish host-wide elimination of rereads.

**Citation key:** `E` = `docs/dev/.spacedock-state/fo-contract-reread-churn.md`; `S` = `skills/first-officer/references/first-officer-shared-core.md`. Existing-source line numbers refer to the supplied `/tmp/astra-main` snapshot; **candidate S:41** refers to the added paragraph in the supplied `073b7e156` patch.

### First: corrections to your premises

- **“No stage report” is false.** The implementation report exists and contains three DONE entries and a Summary (`E:232–243`). Status remains `implementation` (`E:4`).
- **59 reads / 34% is historical provenance, not the current acceptance baseline.** The entity explicitly disclaims it as measurement authority; its replacement baseline is **28/151 = 18.5%**, with a target of at most five reads and 4.0% (`E:104,135,157–158`). The underlying transcript counts are **UNVERIFIED** within this review’s bounds.
- The supplied implementation patch adds exactly the residency paragraph and its blank line; combined-boundary ordering is existing context (`/tmp/astra-artifacts/073b7e156.patch:11–17`). The supplied main snapshot contains combined-boundary ordering but not residency (`S:37–44`).
- The supplied ideation patch edits the entity itself, confirming where the design originated (`/tmp/astra-artifacts/6cb74d408.patch:5–11`).
- An empty `pr:` field is verified (`E:12`); **absence of any GitHub PR is UNVERIFIED**. Likewise, branch membership, “only implementation commit,” the nested repository’s branch, and #754’s merge identity/implementation are **UNVERIFIED**: the supplied patches establish changes, not repository ancestry or remote metadata.

### 1. Is the rule sufficient and correct?

**Correct within its explicit scope:** it changes repeated reads into residency preconditions, invalidates on compaction/replacement evidence, reloads lazily, preserves ordering, and prohibits replacement polling (**candidate S:41**).

**Finding — P1: the new interpretation leaves an explicit Claude read mandate unreconciled.**

The paragraph says:

> “`load` and `read` below mean ensure resident, not repeat a tool call.”

But Claude’s adapter says:

> “A terminal status mutation must complete the write-core `Read`, then the merge-core `Read`, before its Bash call; do not infer either core from this adapter or skip its read.”

That adapter sentence is outside “below” and explicitly names tool calls (`skills/first-officer/references/claude-first-officer-runtime.md:13`). A literal reader can still execute both reads at successive terminal mutations. This preserves precisely the ambiguity the entity identifies as the defect (`E:69`).

**Smallest fix:** make the shared residency interpretation explicitly govern references to these bodies throughout the contract and adapters; qualify concrete `Read` requirements as applying when residency must first be established or restored.

**Potential unsafe suppression:** the word **“Only”** permits invalidation solely on a harness/captain compaction cue or replacement evidence (**candidate S:41**). If a body is demonstrably absent but neither cue is available, that exclusivity conflicts with the opening requirement that it actually be resident. This is a **conditional safety risk, not a verified host failure**. Make actual availability authoritative: historical “loaded” memory or a summary is not the body. Do not introduce polling.

### 2. Does Pi’s presenter requirement contradict residency?

**Not literally. It exposes an intentional coverage hole.**

The two sentences are:

> “One successful load satisfies later triggers for that body in the same uncompacted context; `load` and `read` below mean ensure resident, not repeat a tool call.”  
> — **candidate S:41**

> “After this reread succeeds, load `spacedock:present-gate`.”  
> — `skills/first-officer/references/pi-first-officer-runtime.md:15`

`present-gate` is **not** among the named deferred bodies (`S:43–50`). The entity expressly excludes presenter loads from the paragraph and its numerator (`E:137`). Therefore claiming these sentences directly contradict each other would overstate the evidence.

**Finding — P2: presenter churn remains outside the delivered fix.** This is not uniquely Pi: the lifecycle also says “Load spacedock:present-gate” on each preparation sequence (`skills/fo-gate-lifecycle/SKILL.md:25`).

If the intended outcome includes presenter rereads, extend the shared scope and change the Pi sentence to:

> After the committed-gate reread succeeds, ensure `spacedock:present-gate` is resident under the shared residency rule, then render this gate once from its fresh evidence.

Apply the same distinction to the lifecycle instruction. **Reuse the rendering contract, never the previous gate’s evidence or presentation** (`skills/present-gate/SKILL.md:13,38–39`).

### 3. Does it survive compaction, given #754?

**Conditionally, not demonstrated.**

The exact contract-level signal is:

> “harness notice or captain cue”

That is already the shared continuity condition (`S:143`) and becomes the residency invalidator (**candidate S:41**).

| Host | Exact automatic, model-visible signal established by permitted sources? |
|---|---|
| Claude | **UNVERIFIED.** The runtime adapter defines no compaction-notice binding (`skills/first-officer/references/claude-first-officer-runtime.md:1–37`). |
| Codex | **UNVERIFIED.** The adapter defines runtime capabilities but no compaction-notice binding (`skills/first-officer/references/codex-first-officer-runtime.md:5–30`). The entity records transcript compaction events, not proof of what the resumed FO sees (`E:106–108`). |
| Pi | **UNVERIFIED.** The runtime implementation does not identify a compaction-notice binding (`skills/first-officer/references/pi-first-officer-runtime.md:5–15`). |

These are **not findings that the hosts lack signals**; the bounded evidence does not establish them.

A captain’s explicit compaction message qualifies by contract on all three. Automatic delivery, retention of the residency rule itself, and post-compaction access to `{first_officer_base}` remain unproven (`S:143`; `skills/first-officer/SKILL.md:25`).

Even accepting the described #754 behavior, **rereading durable workflow state is not proof that a deferred contract body remains in context**. The shared rule separately requires re-satisfying both load preconditions and state reads (`S:143`). #754’s actual implementation is outside the supplied two commit artifacts, so its claimed protection remains **UNVERIFIED**.

### 4. Which approved boundary is missing?

**None of the four is omitted from the composed contract.**

| Boundary | Preserved by |
|---|---|
| Gate | Lifecycle residency before gate effects (`S:44`). |
| Write | Own completed host event, after lifecycle when both apply, before mutation (`S:48`). |
| Merge | Terminal/recovery load plus write-before-merge ordering (`S:41,49`). |
| Post-compaction | Lazy invalidation/reload in **candidate S:41**, plus continuity (`S:143`). |

The paragraph does not individually spell out “gate,” “write,” and “merge”; it preserves them by reference to existing triggers. That is not silence on one boundary. The unresolved issue is host interpretation and evidence, not a missing fourth boundary (`E:32,90`).

### 5. What must change before landing?

1. **Resolve the Claude ambiguity and explicitly settle presenter scope.** Expanding beyond the approved one-file surface requires reapproval under the entity’s own constraint (`E:151`).
2. **Provide inspectable behavioral evidence.**  
   **Finding — P1:** the implementation report claims a replay that “retained five first loads,” but supplies no trace locator or fresh-host invocation proving the changed FO—not a filtered historical stream or hand-built oracle—produced those events (`E:236–239`). AC-1 explicitly requires replay against a fresh FO (`E:158`). The report exists; its behavioral claims are **UNVERIFIED**, not disproved.
3. **Demonstrate real compaction on each host.** Injected cues test response to a cue, not automatic cue delivery. The current plan explicitly waives Claude/Pi live coverage (`E:164,181`); that waiver cannot substantiate a per-host compaction claim.

**Falsifiable acceptance I would demand:**

- **All three hosts:** actual FO transcripts with repeated gate, write, dispatch, status, and terminal triggers. Each unchanged covered body loads once per context window; later triggers cause no reads or replacement probes. Record underlying body deliveries, not merely skill invocation counts (`E:160–161,181`).
- **Claude:** two terminal mutations demonstrate no repeated write/merge `Read` calls while resident, while initial and invalidated loads retain separate successful host events (`claude-first-officer-runtime.md:13`; `S:41,48–49`).
- **Codex:** fresh-FO historical replay meets **≤5 reads and ≤4.0%**, preserves workflow effects, and does not simply delete redundant historical events (`E:157–158,170`).
- **Pi:** two gates reuse any newly covered presenter body but independently commit/reread and present each gate before decision mutation (`pi-first-officer-runtime.md:15`; `present-gate/SKILL.md:13`).
- **Each host’s real compaction:** capture the exact model-visible boundary signal; show zero eager deferred reads, then exactly one required reload at each next trigger. Separately test direct replacement, including a changed instruction whose execution proves the new body was used (`E:163–167`).
- Reject traces with repeated reads, omitted reloads, premature effects, or reordered initial/invalidation loads; attach durable-state comparisons and raw trace locations (`E:167,170,179`).

No tests were run in this read-only review. The report records a default race-test timeout followed by a passing extended-timeout run; those executions remain **UNVERIFIED** here (`E:239`).