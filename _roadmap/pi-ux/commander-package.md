# Commander package — pi-ux sprint

Read this file first. It is your whole assignment. It is written so that you can
run without asking the captain at every step.

## 1. Mission

Ship the remaining pi-ux sprint work as stacked pull requests. Keep the workflow
state true after every change. Present each gate once. Open pull requests when a
gate approves. Never merge.

## 2. Your grant (the conn)

The captain gave this grant in the session of 2026-10-03:

> make this work

Its scope, read plainly:

- You may **approve or reject gates** on the captain's behalf, and record the
  decision with `--actor agent:first-officer`, `--conn-quote "make this work"`,
  `--conn-source "captain, pi session 2026-10-03"`.
- You may **open pull requests** for a task whose validation gate has passed.
- You may **dispatch workers** and build dispatch artifacts.
- You may **advance entity state** and commit state changes.

Outside the grant. Do not do these. If a task needs one, present it to the
captain and stop.

- Do not change a task's scope, criteria, or surface.
- Do not merge a pull request.
- Do not cut a release or a tag.
- Do not spend a live run you were not told to spend.

## 3. Where things stand (2026-10-03)

| id | task | stage | branch | pull request |
|---|---|---|---|---|
| `mc` | pi doctor probes a stale layout | done | merged | #817, merged, archived |
| `mh` | pi-live lane pins and assertions | implementation | `spacedock-ensign/pi-live-lane-pin-refresh` | #816 open, base `main`, mergeable |
| `9w` | delegated gate approval through successor dispatch | ideation, gate prepared | none | none |
| `271` | remove the install-attempt file | validation passed | `spacedock-ensign/remove-pi-install-attempt-sentinel` | none |
| `ekw` | pi session identity | retired from the sprint | none | none |

`mh` carries three things inside it: its own lane pins and assertions, the
binding reconciliation carried over from `d52`, and the isolated-home setup
contract carried over from `3w1`. Its criteria run AC-1 through AC-9.

## 3b. The deliverable is one ready-to-merge pull-request stack

The captain reviews the pull requests and merges them. Build the stack. Never
merge it yourself.

```
main
 └── spacedock-ensign/pi-live-lane-pin-refresh            mh   PR #816, base main
      └── spacedock-ensign/remove-pi-install-attempt-sentinel   271  base mh's branch
           └── spacedock-ensign/pi-delegated-gate-continuation-reliability  9w  base 271's branch
```

Order matters. `9w` and `mh` both change files in `internal/ensigncycle`, so
`9w` must sit above `mh`. `271` changes one instruction file and touches neither
package, so it sits between them. Merge order is bottom-up: `mh`, then `271`,
then `9w`.

Each layer's pull request body must carry four things: its base branch, its
validation result, the live lane run id and conclusion for that layer, and the
deferred risks it ships.

"Fully validated" means, for each layer: its validation stage passed with the
criteria cited, and the live lane ran green at the stack tip with
`-f live_cadence=pi`, with the run id recorded and plainly stated as a tip run in
that layer's pull request body. A layer whose only evidence is offline is not
finished unless its own gate says so in words.

A stacked layer does not trigger the lane, because the lane triggers on pull
requests that target `main`. Run it by hand at the stack tip, once, and record
that single run for every layer below it: the tip contains all of their code.

The subscription credential is an OAuth snapshot, and it stays usable until its
refresh token is rotated. Three facts follow from how the lane seeds it. Each
journey builds its own isolated home from the same stored payload and converts
it without refreshing, so one run may refresh zero, one or many times, on
demand. Nothing writes a rotated token back to GitHub. And a broken OAuth secret
does not fall back to `OPENAI_API_KEY`, because OAuth takes precedence.

So reuse the snapshot for the runs you need. It fails only when the refresh
token has already been rotated elsewhere — `refresh_token_reused` means a stale
snapshot, not a code result. Do not approve a duplicate pending live deployment
for its own sake, and do not plan a run per layer: one run at the tip covers
every layer below it.

## 4. Do these, in this order

### Step 1 — `mh`: finish the live proof

1. The grade fix landed as `32049bdbd`. It makes the smoke require both the
   `subagent` tool and the `intercom` tool in the run's recorded tool calls, and
   require that neither `PI_SUBAGENTS_PACKAGE_ROOT` nor
   `PI_INTERCOM_PACKAGE_ROOT` is present in the child environment, including as
   an empty assignment.
2. A live run is already in flight on the branch. Confirm its result:

   ```bash
   gh run list --workflow "Runtime Live E2E" --limit 3 \
     --json databaseId,status,headBranch,conclusion
   gh run view <id> --repo spacedock-dev/spacedock --log-failed
   ```

   To start one if none is running:

   ```bash
   gh workflow run "Runtime Live E2E" \
     --ref spacedock-ensign/pi-live-lane-pin-refresh -f live_cadence=pi
   ```

3. The lane does not trigger on a branch whose base is not `main`, so use the
   manual command above. Never approve a *pending deployment* run on `main`; it
   duplicates the branch run across three live lanes.
4. When the lane is green, advance the task and dispatch validation:

   ```bash
   spacedock status --workflow-dir docs/dev \
     --set pi-live-lane-pin-refresh status=validation started
   spacedock state commit pi-live-lane-pin-refresh
   spacedock dispatch build --workflow-dir docs/dev \
     --entity-path docs/dev/.spacedock-state/pi-live-lane-pin-refresh/index.md \
     --stamp --stage validation --checklist-file - <<'CL'
   <one checklist line per item>
   CL
   ```

5. A previous validation of this task recommended REJECTED for two reasons: the
   live proof was missing, and our own test did not observe what the criterion
   promised. The second reason is now fixed in code. The first reason closes
   when the lane is green and the grade artifact is attached as evidence.
6. Present the validation gate to the captain. Include the lane run id, its
   conclusion, and the grade artifact path.

### Step 2 — `271`: present the gate, then open the pull request

1. Its validation passed at instruction-only scope. One file changed:
   `skills/first-officer/references/fo-install.md`, four insertions and five
   deletions. The install-attempt file is gone.
2. The gate is prepared at
   `briefing:271f46crset81jwf2rgast0c:validation:attempt-1:revision-1`. Present
   it once, with the deferred risk stated: no check on HEAD exercises the
   install-offer flow, so the cross-session offer behaviour is unobserved. The
   shippable claim is one sentence — the file is gone and nothing consults it.
3. Rebase this layer onto `mh`'s branch, so the stack is linear:

   ```bash
   W=.worktrees/spacedock-ensign-remove-pi-install-attempt-sentinel
   git -C "$W" fetch origin
   git -C "$W" rebase origin/spacedock-ensign/pi-live-lane-pin-refresh
   git -C "$W" push --force-with-lease origin \
     HEAD:refs/heads/spacedock-ensign/remove-pi-install-attempt-sentinel
   ```

   If the rebase conflicts, stop and report the exact paths. Do not resolve a
   conflict by editing the other layer.
4. On approval, open the pull request with the stacked base:

   ```bash
   gh pr create --repo spacedock-dev/spacedock \
     --base spacedock-ensign/pi-live-lane-pin-refresh \
     --head spacedock-ensign/remove-pi-install-attempt-sentinel \
     --title "fix(install): remove the machine-wide install-attempt file" \
     --body "<base branch, what changed, checks run, deferred risk, lane run id>"
   ```

5. Record the pull request on the entity with
   `spacedock status --workflow-dir docs/dev --set remove-pi-install-attempt-sentinel pr='#<N>'`,
   then `spacedock state commit remove-pi-install-attempt-sentinel`.
6. Verify the base after creating it, because GitHub retargets a stacked base
   when the parent merges:

   ```bash
   gh pr view <N> --repo spacedock-dev/spacedock --json baseRefName,mergeable,state
   ```

7. Do not merge it. The merge guards live with the captain.

### Step 3 — `9w`: present its ideation gate

1. The gate has been prepared and committed since earlier in the session. It is
   bound to the design that restores the Pi recorded-gate grade.
2. Present it once. State plainly that the task **measures** a fault it does not
   repair: the conduct repair belongs to `gcmfwfjd9735b58sbzw7xsb8`
   (`repair-pi-recorded-gate-lifecycle`), outside this sprint. Ask the captain
   whether they are buying the measurement or the fix.
3. On approval, base the work on `271`'s branch, which already sits on `mh`'s
   branch, so the stack stays linear:

   ```bash
   git branch spacedock-ensign/pi-delegated-gate-continuation-reliability \
     spacedock-ensign/remove-pi-install-attempt-sentinel
   ```

   Then dispatch, and open its pull request with
   `--base spacedock-ensign/remove-pi-install-attempt-sentinel`.
4. Run the live lane at this layer, which is the stack tip:

   ```bash
   gh workflow run "Runtime Live E2E" \
     --ref spacedock-ensign/pi-delegated-gate-continuation-reliability \
     -f live_cadence=pi
   ```

   Record the run id and the conclusion in the layer's pull request body.
5. When all three layers are open, validated and green, hand the stack to the
   captain for review. State the merge order: `mh`, then `271`, then `9w`.

## 5. Commands you will need

```bash
B=${SPACEDOCK_BIN:-spacedock}
D=docs/dev

# state
$B status --workflow-dir $D --where 'sprint=pi-ux' --fields id,slug,status,pr,worktree
$B status --read <slug> --checklist --json --workflow-dir $D
$B status --read <slug> --ac-scan --json --workflow-dir $D
$B status --read <slug> --workflow-dir $D          # list the body's sections
$B status --workflow-dir $D --set <slug> status=<stage> started
$B state commit <slug>

# gates
$B gate prepare <slug> --workflow-dir $D --question TEXT --artifact <entity path> --summary TEXT
$B gate record <slug> --workflow-dir $D --decision approve --actor agent:first-officer \
  --conn-quote "make this work" --conn-source "captain, pi session 2026-10-03" --consume
$B gate record <slug> --workflow-dir $D --decision revise --actor person:captain --reason TEXT
$B state commit <slug>

# dispatch
$B dispatch build --workflow-dir $D --entity-path <entity path> --stage <stage> \
  [--stamp] --checklist-file - --scope-notes-file <file>

# pull requests
gh pr create --repo spacedock-dev/spacedock --base <base> --head <branch> --title T --body B
gh pr view <N> --repo spacedock-dev/spacedock --json state,baseRefName,mergeable
gh workflow run "Runtime Live E2E" --ref <branch> -f live_cadence=pi
gh run list --workflow "Runtime Live E2E" --limit 3 --json databaseId,status,conclusion
```

## 6. How to dispatch a worker

Use the `worker` agent, fresh context, the `ensign` skill, and an explicit model.
Then set a deadline, because a long suite will eat the default.

| stage | model |
|---|---|
| implementation | `lunaroute/deepseek-4.1-flash` |
| ideation, validation | `openai-codex/gpt-6-astra` |

Shape:

```
subagent({
  agent: "worker", skill: "ensign", context: "fresh",
  cwd: "/Users/clkao/git/spacedock-research/spacedock-v1",
  model: "<from the table>", async: true, timeoutMs: 3600000,
  task: "Read <dispatch artifact path> and treat its content as your assignment.\n\n<scope notes>"
})
```

Set `timeoutMs` explicitly. A previous round died at the default 30 minutes.

## 7. Rules that were learned the hard way

1. **Never run a wide search, and never let a worker run one.** Keep `find`,
   recursive `grep` and recursive reads inside the repo, the named directory, or
   one session directory. A worker searched your whole home directory and the
   captain killed it by hand.
2. **Name the file and the mechanism in every dispatch.** A worker told to make
   something observable, without being told where the evidence lives, will go
   hunting. Put the path in the checklist.
3. **Run the offline suite locally before CI. Never let CI be the first run.**
   The gate CI runs is `go test ./...`, about a minute. Run it in the worktree
   before you push or call a layer validated. Do not ask a worker to prove the
   suite fails at the *base commit* — that doubles the cost. For a base
   comparison, name the specific failing tests instead.
4. **A flaky gate is not a passing gate.** If the offline job fails on a test
   that asserts a timing budget, run that one test several times locally. If it
   passes locally and passed in an earlier CI run at the same commit, say so in
   words, with both run ids. Do not call it green.
4. **One writer per worktree.** Do not rebuild a dispatch artifact while its
   worker is reading it.
5. **Keep the completion gate.** After a worker finishes a stage, read its
   checklist and its criterion scan, then advance and dispatch in the same turn.
6. **A rejected gate is a record, not a silence.** Record `revise` with the
   reason, so the next reader sees why.
7. **Do not use the word "substrate".** Name the packages `pi-subagents` and
   `pi-intercom`, and the tools `subagent`, `bg_wait`, `intercom`. The captain
   rejected that word.
8. **Do not invent labels.** Write "the live proof is missing", not "V-1".
9. **Tautological tests are forbidden.** A test may not assert an artifact's own
text against a copy of that text, and its expected value may not come from the
thing it tests. If its only failure mode is someone editing a copy, delete it.
Every test must be able to fail because a real behaviour broke, or because an
independent source diverged.
10. **The science officer is astra.** You have a standing advisor: a fresh
read-only child on `openai-codex/gpt-6-astra`, agent `reviewer`. Before you bring
the captain a question you are unsure about, get the science officer's opinion and
present it beside your own recommendation. Bring the captain a decision, not a
question. Ask it for a verdict, the cost of the option it rejects, and any fact of
yours it finds wrong. Give it the evidence pack you would have given the captain —
written to a file it can read, because it has no shell and cannot reach GitHub.
11. **Name the shape and the bound in every dispatch.** A dispatch that names a
destination and not a shape buys a library. Say which files, how many functions,
and what line ceiling. A worker with an unbounded brief builds the correct thing
properly, which is how 68 lines of script became 1001 lines of Go.
12. **A value change is a value change.** Do not turn "set this value" into a
restructure. A restructure is its own task with its own justification. Five model
ids cost 402 lines because a value change was answered with a migration.
13. **Never demand proof of a property whose only check compares two artifacts your
own change authored.** If you cannot name an independent source, record the gap
instead of building a guard, and prefer deleting a guard over strengthening it.
14. **Ask who else reads the artifact before you change it.** That question is one
search, and it decides the scope. Do not inherit a task's own framing of its scope.
15. **Verify before you repeat.** A worker's phrase, a document's line, a merge
flag, a token count — run the command. Never present an unverified number.
16. **A correction is not an instruction.** When the captain corrects a fact, fix
the fact. Do not launch work from it.
17. **Serialize heavy work, and re-check the base after a restack.** One suite at a
time; check the machine's load and any orphaned daemons before a live run. After a
restack, confirm that some pull request still faces the base branch, and read every
base back.

## 8. Escalate to the captain, and stop

- A task needs work outside its declared surface.
- A gate decision would change what a task promises.
- A merge, a release, or a tag is needed.
- A live run would be spent and the grant does not already cover it.
- Two tasks in the sprint want the same file and both must change it.

## 9. Keep the record true

After every state change, update the goal mission's next action:

```
~/.pi/agent/missions/projects/e46951e4cd86ce8cf6f77d20f379afffe036c09128b36ac2a216e6f533f6687c/\
119999bf-449e-4955-8d81-02049aceb07d/state.json
```

Write the field `nextReadyAction`. Keep it under 1,200 characters. Name what is
running, what is blocked, and what the captain must decide.

## 10. Stop rules

Stop, and report, when any of these is true.

- Every sprint task is merged, archived, or explicitly deferred by the captain.
- A worker fails twice on the same stage for the same reason.
- The state validator reports an error you cannot clear.
- You need something outside the grant.
