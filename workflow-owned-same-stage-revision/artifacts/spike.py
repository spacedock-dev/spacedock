import pathlib,subprocess,json,os
ROOT=pathlib.Path('/tmp/spacedock-same-stage-evidence-zz1y')
BIN='/opt/homebrew/Caskroom/spacedock@next/0.28.0-pre3/spacedock'
log=[]
def run(args,cwd,expect=0):
 p=subprocess.run(args,cwd=cwd,text=True,capture_output=True); log.append(dict(args=args,cwd=str(cwd),exit=p.returncode,stdout=p.stdout,stderr=p.stderr)); (ROOT/'commands.json').write_text(json.dumps(log,indent=2)); print(p.returncode,' '.join(args),p.stdout,p.stderr,flush=True)
 if expect is not None: assert p.returncode==expect
 return p
w=ROOT/'self';w.mkdir(exist_ok=True)
def write(p,s):p.parent.mkdir(parents=True,exist_ok=True);p.write_text(s)
write(w/'README.md','''---
commissioned-by: spacedock@0.28.0-pre3
id-style: slug
stages:
  states:
    - name: plan
      initial: true
      gate: true
      feedback-to: plan
    - name: done
      terminal: true
---
# Triage-like plan workflow
### `plan`
Correct the local plan against frozen-input.txt, commit it, and present for captain review. No separate reviewer or correction-round publication is declared.
''')
write(w/'task/index.md','---\nid: task\ntitle: Same-stage plan\nstatus: plan\n---\n# Task\n\n## Stage Report: plan\n- DONE: Initial plan\n### Summary\nInitial candidate.\n')
write(w/'task/frozen-input.txt','KEEP message A; DELETE message B\n');write(w/'task/plan.md','KEEP message B; DELETE message A\n')
for a in [['git','init','-q'],['git','config','user.email','spike@example.invalid'],['git','config','user.name','Ideation spike'],['git','add','.'],['git','commit','-qm','seed frozen input and incorrect plan']]:run(a,w)
def sd(*a,expect=0):return run([BIN,*a,'--workflow-dir',str(w)],w,expect)
sd('gate','prepare','task','--question','Approve local plan?','--artifact',str(w/'task/plan.md'),'--summary','initial incorrect plan')
sd('state','commit','task')
sd('gate','record','task','--decision','revise','--actor','person:captain','--reason','Synthetic spike authority: align plan with frozen input')
sd('gate','consume','task',expect=None)
write(ROOT/'checklist.txt','Correct plan against frozen input and commit stage report\n');write(ROOT/'feedback.txt','Synthetic captain revise: KEEP message A; DELETE message B. Correct plan only.\n')
sd('dispatch','build','--entity-path',str(w/'task/index.md'),'--stage','plan','--host','codex','--checklist-file',str(ROOT/'checklist.txt'),'--feedback-context-file',str(ROOT/'feedback.txt'),'--feedback-reflow')
