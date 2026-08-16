import pathlib
base = pathlib.Path("/Users/clkao/.claude/jobs/4e49247e/tmp/fitgate/arms")
src = (base / "A-main.md").read_text()

HEAD = """## Workflow Fit Gate

Before creating or materially reclassifying an entity, verify the work fits the commissioned workflow's subject and value model. Write authorization is not workflow-fit authorization. The write classifier is not evidence of fit either: a path's class says who may write it, never whether this workflow should be tracking this work.

A new entity belongs only when it produces or validates a deliverable the workflow exists to track, using that workflow's own `entity-type`, README purpose, stage outputs, and acceptance/proof policy.

Name the output's existing home before filing. If a documented process already owns this output — a release ritual, a debrief, a reconciliation ledger, a runbook, a roadmap or planning doc, a registry, or another workflow — the work belongs there, and filing it here duplicates its owner instead of adding a deliverable. """

OVERFIT = """Release narratives, status summaries, reports, and standalone decisions have such homes. """

TAIL_LIST = """So do FO/process maintenance, workflow refits, split-root migration, debriefing, status reporting, cleanup of agent/session state, and operating-ledger work; none belong in a product/dev workflow unless its README names that class as an executable deliverable.

A fit failure is not repaired by adding a shippable mechanism. If the work does not belong here it does not belong at any shape: reshaping it until it satisfies the proof policy buys admission with machinery the workflow never needed. "It can carry a real value AC" answers the proof policy, not the fit question.

If fit is ambiguous, stop before `spacedock new` or `status --set` and ask the captain where the work should live.

"""

ANCHOR = "## FO Write Scope\n"
OL_A = "the FO runs `«state.commit»(slug)` after `new` to commit and sync it."
OL_ADD = " `spacedock new` is only the atomic creation mechanism after the Workflow Fit Gate passes; it does not decide whether the work belongs in this workflow."

def build(name, section):
    out = src.replace(ANCHOR, section + ANCHOR).replace(OL_A, OL_A + OL_ADD)
    (base / name).write_text(out)
    print(name, len(out.encode()) - len(src.encode()), "bytes over main;", out.count("\n") - src.count("\n"), "lines added")

build("D-home.md", HEAD + TAIL_LIST)
build("E-home-named.md", HEAD + OVERFIT + TAIL_LIST)
