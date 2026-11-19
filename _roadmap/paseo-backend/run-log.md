# Paseo dispatch run log

18 November 2026, local macOS daemon. One first officer (FO) ran a stage-completion probe and a child return-channel probe, each with a fresh ensign. Commands below use the recorded paths and IDs. Outputs are selected fields or excerpts. A blank Git status means the checkout was clean.

## Runtime and access

```console
$ git rev-parse --short HEAD
db4728199
$ paseo --version
0.16.1
$ spacedock --version
spacedock 0.29.0-pre0
OS: darwin/arm64
Channel: edge (spacedock@spacedock-edge)
Runtime: codex (CODEX_THREAD_ID)
Sandbox: agent-safehouse
contract 3
$ paseo status --json  # selected fields
{"daemonVersion":"0.16.0","localDaemon":"running","connectedDaemon":"reachable"}
$ df -k /tmp  # Available, 1024-byte blocks
4987860
```

MCP calls and selected returns:

```text
list_workspaces({})
  title=paseo-backend probe  workspaceId=wks_b6ea4c89878703f3
  projectId=prj_d8edf163a62d3477  isolation=local  kind=local_checkout
  cwd=/Users/clkao/git/spacedock-research/spacedock-v1
get_agent_status({agentId:"7aaba008-a8cd-4f7d-af28-7932445591ec"})
  provider=safehouse-codex  model=gpt-6.2-astra
  role=worker  repo=spacedock-v1  mandate=paseo-backend
  paseo-safehouse.paseo-access=daemon
list_providers({})
  providers_count=15  (includes safehouse-claude, safehouse-codex and safehouse-pi)
inspect_provider({provider:"safehouse-codex"})
  status=available  modes_ids=auto,auto-review,full-access
list_models({provider:"safehouse-codex"})
  models_count=10  (includes gpt-6.2-astra, label GPT-6.2-Astra)
```

`paseo plugin ls --json` listed `paseo-safehouse` as enabled. The FO used its existing daemon access throughout.

## Fixture and dispatch

`mktemp -d /tmp/spacedock-paseo.XXXXXX` returned `/tmp/spacedock-paseo.qjyrQl`. On this host its physical path starts with `/private/tmp`. Fixture setup:

```sh
set -eu
cd /private/tmp/spacedock-paseo.qjyrQl
git init -q -b main fixture
cd fixture
git config user.name 'Paseo Probe'
git config user.email 'probe@example.invalid'
mkdir -p docs/probe
cat > .gitignore <<'EOF'
.worktrees/
docs/probe/.spacedock-state/
EOF
cat > docs/probe/README.md <<'EOF'
---
entity-type: task
id-style: slug
state: .spacedock-state
stages:
  defaults:
    worktree: false
    concurrency: 1
  states:
    - name: backlog
      initial: true
    - name: implementation
      worktree: true
    - name: done
      terminal: true
---
# Paseo receipt fixture

### backlog

Seed one task.

### implementation

Write receipt.txt in the assigned code worktree with the marker and the branch name. Commit only that file. Append an implementation stage report to the canonical task in the state checkout, including each checklist item, the code commit and a Summary. Commit only the entity. Keep frontmatter intact.

- **Outputs:** A code receipt and a committed state report.

### done

The first officer verifies the two commits.
EOF
git add .gitignore docs/probe/README.md
git commit -qm 'seed workflow'
git update-ref refs/remotes/origin/main HEAD
git commit --allow-empty -qm 'local trunk ahead of synthetic origin/main'
git worktree add -qb spacedock-state/probe docs/probe/.spacedock-state
s=docs/probe/.spacedock-state
mkdir "$s/paseo-qjyrql"
cat > "$s/paseo-qjyrql/index.md" <<'EOF'
---
id: paseo-qjyrql
title: Paseo receipt
status: implementation
started:
worktree:
---
# Paseo receipt

Record evidence of a Codex dispatch through Paseo.
EOF
git -C "$s" add -- paseo-qjyrql/index.md
git -C "$s" commit -qm 'seed task' -- paseo-qjyrql/index.md
git remote -v
git rev-parse main origin/main
```

`git remote -v` printed nothing. The first hash below is local `main`. The second is a deliberately older local *remote-tracking ref*, created without contacting a remote.

```text
f98e81a10a5c397629231fa6db16565d8c617aac
ff2b5607126c4e954f650dbd736d21287d891483
```

From the fixture root:

```sh
spacedock dispatch build --host codex --workflow-dir "$PWD/docs/probe" \
  --entity-path "$PWD/$s/paseo-qjyrql/index.md" --stage implementation \
  --checklist-file - --stamp > ../envelope.json <<'EOF'
Commit receipt.txt with PASEO-RECEIPT-QJYRQL and the actual branch in the assigned code worktree.
Commit the exact-stage report with the code commit and the steering word supplied by the first officer.
EOF
cat ../envelope.json
```

```json
{
  "schema_version": 2,
  "subagent_type": "spacedock:ensign",
  "description": "Paseo receipt: implementation",
  "fetch_commands": [
    "/opt/homebrew/Caskroom/spacedock@next/0.29.0-pre0/spacedock dispatch show-stage-def --workflow-dir /private/tmp/spacedock-paseo.qjyrQl/fixture/docs/probe --stage implementation"
  ],
  "dispatch_file_path": "/tmp/spacedock-dispatch/paseo-qjyrql-implementation.md",
  "prompt": "$spacedock:ensign; then Read /tmp/spacedock-dispatch/paseo-qjyrql-implementation.md and treat its content as your assignment.",
  "model": null,
  "name": "paseo-qjyrql-implementation"
}
```

```console
$ git worktree list
/private/tmp/spacedock-paseo.qjyrQl/fixture                                           f98e81a [main]
/private/tmp/spacedock-paseo.qjyrQl/fixture/.worktrees/spacedock-ensign-paseo-qjyrql  f98e81a [paseo-qjyrql]
/private/tmp/spacedock-paseo.qjyrQl/fixture/docs/probe/.spacedock-state               764719a [spacedock-state/probe]
$ git -C "$s" log -1 --format='%h %s'
764719a dispatch: paseo-qjyrql entering implementation
$ cat "$s/paseo-qjyrql/index.md"  # stamped fields
started: 2026-11-18T20:59:29Z
worktree: .worktrees/spacedock-ensign-paseo-qjyrql
$ rg -n 'checkout|worktree|entity' /tmp/spacedock-dispatch/paseo-qjyrql-implementation.md  # excerpts
11:Your working directory for CODE is /private/tmp/spacedock-paseo.qjyrQl/fixture/.worktrees/spacedock-ensign-paseo-qjyrql
13:Your git branch for CODE is paseo-qjyrql. All CODE commits MUST be on this branch in the worktree. Do NOT switch branches or commit code to main.
16:Read the entity file at /private/tmp/spacedock-paseo.qjyrQl/fixture/docs/probe/.spacedock-state/paseo-qjyrql/index.md for the current spec.
$ git worktree add '/private/tmp/spacedock-paseo.qjyrQl/second' paseo-qjyrql
Preparing worktree (checking out 'paseo-qjyrql')
fatal: 'paseo-qjyrql' is already used by worktree at '/private/tmp/spacedock-paseo.qjyrQl/fixture/.worktrees/spacedock-ensign-paseo-qjyrql'
```

The last command exited 128. This exercised Git's same-branch conflict in the fixture. No second Paseo-managed worktree was requested.

## Workspace and spawn

```text
create_workspace({
  isolation:"local", projectId:"prj_d8edf163a62d3477",
  path:"/private/tmp/spacedock-paseo.qjyrQl/fixture/.worktrees/spacedock-ensign-paseo-qjyrql", title:"Paseo receipt fixture qjyrql"
})
  workspaceId=wks_8aec77dc0eda8e94  parentWorkspaceId=wks_b6ea4c89878703f3
  cwd=/private/tmp/spacedock-paseo.qjyrQl/fixture/.worktrees/spacedock-ensign-paseo-qjyrql
  isolation=worktree  kind=worktree
```

The following is the actual input to `create_agent`. The worker name and description came from the envelope. The provider selection, workspace, labels, callback and text after the pointer were supplied by the FO.

```json
{
  "provider": "safehouse-codex/gpt-6.2-astra",
  "workspaceId": "wks_8aec77dc0eda8e94",
  "title": "Paseo receipt: implementation",
  "labels": {
    "role": "ensign",
    "repo": "spacedock-v1",
    "mandate": "paseo-backend",
    "spacedock.worker": "paseo-qjyrql-implementation",
    "spacedock.epoch": "qjyrql-1"
  },
  "settings": {
    "modeId": "auto"
  },
  "notifyOnFinish": true,
  "initialPrompt": "$spacedock:ensign; then Read /tmp/spacedock-dispatch/paseo-qjyrql-implementation.md and treat its content as your assignment.\nThe two allowed write targets are the assigned receipt.txt and the canonical task body in the fixture state checkout, plus their path-scoped Git commits. Report pwd -P and branch. I will send the steering word during this run. If the word has not arrived when needed, ask me through the Paseo parent channel. No push, PR, external comment, agent launch, or other project writes. If a permission blocks a required write, report the exact diagnostic."
}
```

```text
agentId=0cb685d1-8e49-4795-9705-d0b7759a5045  type=safehouse-codex  status=running
cwd=/private/tmp/spacedock-paseo.qjyrQl/fixture/.worktrees/spacedock-ensign-paseo-qjyrql
workspaceId=wks_8aec77dc0eda8e94  currentModeId=full-access
guidance: You will get notified when the created agent finishes, errors, or needs permission.
```

## Steering and roster

```text
send_agent_prompt({agentId:"0cb685d1-8e49-4795-9705-d0b7759a5045",activeTurnBehavior:"steer",
  notifyOnFinish:true,prompt:"The steering word is KINGFISHER. Include it in the state report. If your existing agent-scoped Paseo MCP offers send_agent_prompt, send a brief progress acknowledgement with the word to the mandate parent 7aaba008-a8cd-4f7d-af28-7932445591ec, then finish your assignment."})
  success=true  deliveryStatus=queued  status=running
  messageId=b6676533-58b5-4488-be87-7ea71830e47d
list_agents({cwd:"/private/tmp/spacedock-paseo.qjyrQl/fixture/.worktrees/spacedock-ensign-paseo-qjyrql",limit:12,includeBackground:true})
  agents_count=1  agents_ids=0cb685d1-8e49-4795-9705-d0b7759a5045
  status=running  provider=safehouse-codex  model=gpt-6.2-astra
  labels: role=ensign, repo=spacedock-v1, mandate=paseo-backend,
    spacedock.worker=paseo-qjyrql-implementation, spacedock.epoch=qjyrql-1,
    paseo.parent-server-id=srv_j3VeUq-Ba722, paseo.parent-agent-id=7aaba008-a8cd-4f7d-af28-7932445591ec
```

```console
$ paseo ls --global --label 'spacedock.epoch=qjyrql-1' --json | jq 'map({id,status,cwd})'
[{"id":"0cb685d1-8e49-4795-9705-d0b7759a5045","status":"running","cwd":"/private/tmp/spacedock-paseo.qjyrQl/fixture/.worktrees/spacedock-ensign-paseo-qjyrql"}]
$ paseo wait 0cb685d1-8e49-4795-9705-d0b7759a5045 --timeout 1 --json
```

Exit 0. The JSON `status` was `timeout`. The returned activity included the worker's message:

> Received KINGFISHER. `pwd -P` returned `/private/tmp/spacedock-paseo.qjyrQl/fixture/.worktrees/spacedock-ensign-paseo-qjyrql`, and the branch is `paseo-qjyrql`. This session has no agent-scoped Paseo `send_agent_prompt` tool.

`get_agent_status({agentId})` during the run returned `status=running`, `currentModeId=full-access`, `pendingPermissions=[]` and `requiresAttention=false`. Worker-to-parent addressed delivery was not observed.

## Completion and durable proof

The FO received a pushed `paseo-system` notification without another wait:

```text
kind=agent-notification  event=finished  source-server-id=srv_j3VeUq-Ba722
source-agent-id=0cb685d1-8e49-4795-9705-d0b7759a5045
message-id=ef9f76e1-3e66-4874-b585-619af798e5c0
Done: Paseo receipt completed implementation. Report written to /private/tmp/spacedock-paseo.qjyrQl/fixture/docs/probe/.spacedock-state/paseo-qjyrql/index.md.
```

```text
get_agent_status({agentId:"0cb685d1-8e49-4795-9705-d0b7759a5045"})
  status=idle  activeTurn=null  currentModeId=full-access
  updatedAt=2026-11-18T21:00:40.405Z  requiresAttention=true  attentionReason=finished
  pendingPermissions=[]
  lastUsage: inputTokens=30045 cachedInputTokens=29568 outputTokens=48
             contextWindowUsedTokens=30093 contextWindowMaxTokens=258400
paseo wait 0cb685d1-8e49-4795-9705-d0b7759a5045 --timeout 5 --json
  status=idle
paseo inspect 0cb685d1-8e49-4795-9705-d0b7759a5045 --json
  Status=idle  Mode=full-access  ParentAgentId=7aaba008-a8cd-4f7d-af28-7932445591ec
  LastUsage: InputTokens=30045 OutputTokens=48 CachedTokens=29568 CostUsd=0
```

```console
$ w=.worktrees/spacedock-ensign-paseo-qjyrql
$ s=docs/probe/.spacedock-state
$ git -C "$w" log -1 --format='%H %s' --name-only
a5a82f41505306dd902d89f22b470d43895ff6d7 Record Paseo dispatch receipt
receipt.txt
$ git -C "$s" log -1 --format='%H %s' --name-only
e6524c9f45ba12895411adcf88135ca50897db7e Record implementation report for Paseo receipt
paseo-qjyrql/index.md
$ cat "$w/receipt.txt"
PASEO-RECEIPT-QJYRQL
paseo-qjyrql
$ git -C "$s" diff HEAD^ HEAD -- paseo-qjyrql/index.md  # excerpt
 worktree: .worktrees/spacedock-ensign-paseo-qjyrql
 # Paseo receipt
 
 Record evidence of a Codex dispatch through Paseo.
+
+## Stage Report: implementation
+
+- DONE: Commit receipt.txt with PASEO-RECEIPT-QJYRQL and the actual branch in the assigned code worktree.
+  Code commit `a5a82f41505306dd902d89f22b470d43895ff6d7` contains only `receipt.txt`. Its two lines record the marker and the observed branch, `paseo-qjyrql`.
+- DONE: Commit the exact-stage report with the code commit and the steering word supplied by the first officer.
+  The first officer supplied `KINGFISHER` during the run.
+
+### Summary
+
+The worker recorded the receipt on branch `paseo-qjyrql`. The physical working directory reported by `pwd -P` was `/private/tmp/spacedock-paseo.qjyrQl/fixture/.worktrees/spacedock-ensign-paseo-qjyrql`.
$ for tree in . "$w" "$s"; do echo "status $tree"; git -C "$tree" status --porcelain; done
status .
status .worktrees/spacedock-ensign-paseo-qjyrql
status docs/probe/.spacedock-state
```

## Cross-host envelope

After the live run, a CLI subprocess used a Claude runtime marker and the same checklist. This was a dispatch-build probe. The real parent and worker both ran Codex.

```sh
set -eu
cd '/private/tmp/spacedock-paseo.qjyrQl/fixture'
env -u CODEX_THREAD_ID CLAUDECODE=1 spacedock --version | rg '^Runtime:'
env -u CODEX_THREAD_ID CLAUDECODE=1 spacedock dispatch build --host codex --workflow-dir "$PWD/docs/probe" --entity-path "$PWD/docs/probe/.spacedock-state/paseo-qjyrql/index.md" --stage implementation --checklist-file - > ../cross-host.json <<'EOF'
Commit receipt.txt with PASEO-RECEIPT-QJYRQL and the actual branch in the assigned code worktree.
Commit the exact-stage report with the code commit and the steering word supplied by the first officer.
EOF
cmp ../envelope.json ../cross-host.json && echo 'Codex envelopes byte-identical'
```

```text
Runtime: claude (CLAUDECODE)
Codex envelopes byte-identical
```

## Retry and budget contracts

`sed -n '79,95p' skills/first-officer/references/claude-fo-dispatch.md` returned the retry rules: a spawn error, or an ended session missing both signal and report, consults the per-entity `### Dispatch Retries` ledger. The first retry nudges a live worker or dispatches a fresh `-retry` worker. A second failure holds that entity. The ledger format printed was:

```text
- Retry 1: {stage} — {agent-error | no-completion-signal}; {nudged | re-dispatched -retry}
```

`sed -n '68,82p;147,164p' skills/first-officer/references/fo-dispatch-core.md` returned the current-stage report verification gate and the reuse rule that an over-budget or unavailable reading from a bound probe forces fresh dispatch. `rg -n 'threshold|reuse_ok' internal/claudeteam/contextbudget.go` showed `thresholdPct = 60` and `reuseOK := usagePct <= thresholdPct`. These commands read the contract. This run exercised a wait timeout and a successful assignment, with no provider error, stranded worker or replacement.

## Shutdown and cleanup

```text
kill_agent({agentId:"0cb685d1-8e49-4795-9705-d0b7759a5045"}) → {"success":true}
get_agent_status({agentId:"0cb685d1-8e49-4795-9705-d0b7759a5045"}) → status=closed, snapshot.status=closed
paseo ls --global --label 'spacedock.epoch=qjyrql-1' --json | jq 'map({id,status})'
  [{"id":"0cb685d1-8e49-4795-9705-d0b7759a5045","status":"closed"}]
archive_workspace({workspaceId:"wks_8aec77dc0eda8e94"})
  archivedAgentIds_count=1  archivedAgentIds=["0cb685d1-8e49-4795-9705-d0b7759a5045"]
  removedDirectory=false
paseo ls --all --global --label 'spacedock.epoch=qjyrql-1' --json
  id=0cb685d1-8e49-4795-9705-d0b7759a5045  status=closed
```

From the fixture root, `git worktree remove .worktrees/spacedock-ensign-paseo-qjyrql` and `git worktree remove docs/probe/.spacedock-state` succeeded. `git worktree list` then contained only the fixture's main checkout. Cleanup checked the fixture HEAD, clean Git status, generated file inventory and pointer contents before unlinking them. Final output:

```text
Fixture and assignment removed: True
df -k /tmp: 5019064 available 1024-byte blocks
```

## Child return channel

A follow-up used a second disposable ensign to test a child with the same Safehouse daemon label as the FO. The fixture used the split-root setup above with slug `paseo-return-7uqzwp`. Its implementation stage asked the worker to send `RETURN-7UQZWP` to this FO and commit the delivery result in the canonical task body. The two probe ensigns ran sequentially.

```console
$ mktemp -d /tmp/spacedock-return.XXXXXX
/tmp/spacedock-return.7uqZWP
$ paseo status --json | jq '{daemonVersion, connectedDaemon}'
{"daemonVersion":"0.16.0","connectedDaemon":"reachable"}
$ workflow=/private/tmp/spacedock-return.7uqZWP/fixture/docs/probe
$ entity="$workflow/.spacedock-state/paseo-return-7uqzwp/index.md"
$ spacedock dispatch build --host codex --stamp --workflow-dir "$workflow" --entity-path "$entity" --stage implementation --checklist-file -  # checklist on stdin, selected stdout
"name": "paseo-return-7uqzwp-implementation"
"model": null
"prompt": "$spacedock:ensign; then Read /tmp/spacedock-dispatch/paseo-return-7uqzwp-implementation.md and treat its content as your assignment."
$ spacedock status --workflow-dir "$workflow" --read "$entity" --json  # selected fields
path=/private/tmp/spacedock-return.7uqZWP/fixture/docs/probe/.spacedock-state/paseo-return-7uqzwp/index.md
frontmatter.worktree=.worktrees/spacedock-ensign-paseo-return-7uqzwp
```

Joining that field to `git -C "$workflow" rev-parse --show-toplevel` and calling Python `Path.resolve(strict=True)` returned the physical `path` used below. No worker assignment text was parsed.

```text
create_workspace({isolation:"local",projectId:"prj_d8edf163a62d3477",
  path:"/private/tmp/spacedock-return.7uqZWP/fixture/.worktrees/spacedock-ensign-paseo-return-7uqzwp",title:"Paseo return channel fixture 7uqzwp"})
  workspaceId=wks_4e53fa26c7fa08f4  isolation=worktree  kind=worktree
```

Exact follow-up `create_agent` input:

```json
{
  "provider": "safehouse-codex/gpt-6.2-astra",
  "workspaceId": "wks_4e53fa26c7fa08f4",
  "title": "Paseo return channel: implementation",
  "labels": {
    "role": "ensign",
    "repo": "spacedock-v1",
    "mandate": "paseo-backend",
    "spacedock.worker": "paseo-return-7uqzwp-implementation",
    "spacedock.epoch": "7uqzwp-1",
    "paseo-safehouse.paseo-access": "daemon"
  },
  "settings": {
    "modeId": "auto"
  },
  "notifyOnFinish": true,
  "initialPrompt": "$spacedock:ensign; then Read /tmp/spacedock-dispatch/paseo-return-7uqzwp-implementation.md and treat its content as your assignment.\nThis disposable probe authorizes one mid-run message to the specified FO only. Discover your runtime's available Paseo MCP tools, and record the receipt or missing tool in your stage report. If MCP is missing, use the documented paseo send CLI targeting only the supplied FO, with --behavior steer --no-wait, and record which transport worked. Do not inspect other agents, projects or ACLs. No push, PR, comment, additional launch or writes outside the assigned state path and its scoped commit."
}
```

```text
create_agent → agentId=f8645a84-9f41-428f-ba12-4e309cfa024a  status=running  currentModeId=full-access
get_agent_status({agentId:"f8645a84-9f41-428f-ba12-4e309cfa024a"})
  status=running  currentModeId=full-access
  labels: paseo-safehouse.paseo-access=daemon, paseo.parent-server-id=srv_j3VeUq-Ba722,
          paseo.parent-agent-id=7aaba008-a8cd-4f7d-af28-7932445591ec
```

The FO received these two distinct events in order:

```text
kind=agent-message  source-agent-id=f8645a84-9f41-428f-ba12-4e309cfa024a
message-id=059cc77b-0f8b-4769-9c93-600311597bc2  body=RETURN-7UQZWP
kind=agent-notification  event=finished  source-agent-id=f8645a84-9f41-428f-ba12-4e309cfa024a
message-id=83892e20-27a1-4070-a2ee-a14afebf232f
```

```console
$ git -C docs/probe/.spacedock-state show -1 --format='%H %s' -- paseo-return-7uqzwp/index.md
b67ff477c9eba38076801a8b60f096386ebe1567 Record Paseo return message receipt
[…]
+  Paseo accepted the steer request and returned messageId `059cc77b-0f8b-4769-9c93-600311597bc2` with deliveryStatus `queued`.
[…]
+  The runtime exposed `mcp__paseo__send_agent_prompt` and it returned a receipt on the first call.
$ git -C docs/probe/.spacedock-state status --porcelain
```

The test exercised one allowed recipient. It did not test daemon restrictions on other recipients or methods.

```text
kill_agent({agentId:"f8645a84-9f41-428f-ba12-4e309cfa024a"}) → {"success":true}
get_agent_status({agentId:"f8645a84-9f41-428f-ba12-4e309cfa024a"}) → status=closed
archive_workspace({workspaceId:"wks_4e53fa26c7fa08f4"})
  archivedAgentIds=["f8645a84-9f41-428f-ba12-4e309cfa024a"]  removedDirectory=false
Fixture and assignment removed: True
df -k /tmp: 4974828 available 1024-byte blocks
```
