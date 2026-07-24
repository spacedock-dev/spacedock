# Gate review: One-command gate review presentation with atomic result retention (xb) — ideation, attempt 2

**What you are looking at:** the presentation command's design — the machinery that obtains your gate decisions without ever losing one. Three artifacts: this summary, the full design, and the recorder contract it implements against.

**Chosen direction:** one binary command presents a briefing package: validates it, derives the title, launches the Subspace TUI as a blocking child on the caller's terminal transport, and atomically retains the result, review log, and diagnostics on success AND on failure. The presenter stays addressable until TUI exit plus retention — pane creation and timeout are never completion. It implements the provider id-mapping exactly as the recorder contract specifies, hands the validated resolution to the recorder, and never writes gate state itself.

**The proof that matters:** the retention-contract spike runs 9/9 against three real red fixtures — the blank-float EOF, the leave-open result, and the destroyed-approval case: a real launcher process writes an approval and is killed before returning. The interim launcher pattern retains NOTHING on all three (reproducing the live incidents, including your own approval destroyed two nights ago); the command's contract — caller-owned retention directory it never deletes — retains everything on all three.

**New since attempt 1 (your question):** the no-subspace fallback is designed and exercised — detection version-gates the presenter before any side effect; absence exits non-zero naming the install remedy AND the chat fallback (present in conversation, record through the recorder, per the recording-identity ruling); a missing presenter is an ordinary channel selection, never a mode. Spike fixture: absent binary leaves zero side effects, 12/12.

**Honest bounds:** the cross-repo dependency is declared to the argv level — what ships spacedock-side now versus what needs subspace-tui's surfaces; the working-copy ritual is the measured interim baseline it must beat. Preflight: no material findings against this design; one free fold (the destroyed-approval fixture citation) applied and exercised.

**Expected surface:** declared in the design with 2× tolerance; no gate-state writes, no second writer.

**Recommend approve.** Checklist: 3 done, 0 skipped, 0 failed; value ACs measured against the interim baseline that moves the wrong way.

**Decision:** approve = a pending advance to implementation (the sprint Commander consumes it). Revise = annotate. Hold = discuss.

---
*Companion artifacts: the full design (entity snapshot), the recorder contract.*
