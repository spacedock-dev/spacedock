# Gate review: Gate recorder (3k) — ideation, attempt 5 (multi-artifact presentation)

Same design as the attempt you were reviewing — re-presented as the briefing package you asked for. Three artifacts: this summary, the contract spec (open it for the ownership mermaid, right after its first section — the TUI shows the diagram source; it renders graphically on GitHub/the docs site), and the frozen entity snapshot.

**The design, unchanged from attempt 4:**
1. **Resolution-first split (your direction):** 3k records what the decision IS — gate → attempt → briefing binding → exact resolution, record invariants, snapshot digests, resolution-state surfacing. The application layer (what the decision DOES) is h1's, per the aligned responsibility boundary. The recorder round-trips h1's `application` sub-object unchanged, so the eight-entity replay stays green.
2. **Surface:** ~400-650 production LOC + equal tests, tolerance 2×.
3. **Spec owner-tagged section by section** (3k / h1 / xb) with the boundary mermaid; one honest deferral — the provider id-normalization prose elaboration is xb's first designated spec edit under the change protocol.
4. **Retained ACs:** 1, 4, 6 record-subset, 10, 12, 13, 14 — scanner-verified, all evidenced. PR-510 alignment and its four fork defaults stand (fork 1 defaulted: recorder ids stay internal).

**Recommend approve.**

**Decision:** approve = attempt 5 closes with the pending advance to implementation; h1 immediately dispatches with the application scope. Revise = annotate. Hold = discuss.
