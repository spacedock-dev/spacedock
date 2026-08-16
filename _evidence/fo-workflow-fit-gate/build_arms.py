import pathlib
base = pathlib.Path("/Users/clkao/.claude/jobs/4e49247e/tmp/fitgate/arms")
src = (base / "A-main.md").read_text()

SEED = """## Workflow Fit Gate

Before creating or materially reclassifying an entity, the FO must verify the work fits the commissioned workflow's subject and value model. Write authorization is not workflow-fit authorization.

A new entity belongs only when it produces or validates a deliverable the workflow exists to track, using that workflow's own `entity-type`, README purpose, stage outputs, and acceptance/proof policy.

Do not file FO/process maintenance, workflow refits, split-root migration, debriefing, status reporting, cleanup of agent/session state, or operating-ledger work into a product/dev workflow unless the workflow README explicitly names that class as an executable deliverable. Record those in a debrief, reconciliation ledger, runbook, roadmap/planning doc, or a separate workflow/process track instead.

"""

ANTIRESHAPE = """A fit failure is not repaired by adding a shippable mechanism. If the work does not belong in this workflow, it does not belong at any shape: reshaping it until it satisfies the workflow's proof policy buys admission with machinery the workflow never needed. Re-shape only after fit passes.

"""

TAIL = """If fit is ambiguous, stop before `spacedock new` or `status --set` and ask the captain where the work should live.

"""

ANCHOR = "## FO Write Scope\n"
ONELINER_ANCHOR = "the FO runs `«state.commit»(slug)` after `new` to commit and sync it."
ONELINER_ADD = " `spacedock new` is only the atomic creation mechanism after the Workflow Fit Gate passes; it does not decide whether the work belongs in this workflow."

def build(name, section):
    assert src.count(ANCHOR) == 1, "anchor not unique"
    assert src.count(ONELINER_ANCHOR) == 1, "oneliner anchor not unique"
    out = src.replace(ANCHOR, section + ANCHOR)
    out = out.replace(ONELINER_ANCHOR, ONELINER_ANCHOR + ONELINER_ADD)
    (base / name).write_text(out)
    print(name, len(out.encode()), "bytes", "(+%d over main)" % (len(out.encode()) - len(src.encode())))

build("B-seed.md", SEED + TAIL)
build("C-proposed.md", SEED + ANTIRESHAPE + TAIL)
print("A-main.md", len(src.encode()), "bytes")
