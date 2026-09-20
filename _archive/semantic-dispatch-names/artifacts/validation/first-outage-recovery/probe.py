import hashlib,json,os,pathlib,subprocess,tempfile,sys,shlex
B=pathlib.Path(sys.argv[1]).resolve()
root=pathlib.Path(tempfile.mkdtemp(prefix='semantic-outage-audit-'));state=root/'state';state.mkdir()
readme='---\nentity-type: task\nid-style: slug\nstate: state\nstages:\n  states:\n    - name: backlog\n      initial: true\n    - name: done\n      terminal: true\n---\n\n### backlog\n\nWork.\n\n### done\n\nFinished.\n'
(root/'README.md').write_text(readme)
short='validation-first-outage';long='validation-first-outage-'+('readable-long-task-'*5)+'alpha'
for slug in [short,long]: (state/(slug+'.md')).write_text('---\ntitle: Outage audit\nstatus: backlog\n---\n\nBody.\n')
(state/'ignored-sentinel').write_bytes(b'unchanged\x00bytes')
for repo in [root,state]:
 for args in [['init','-q','-b','main'],['add','.'],['-c','user.name=Audit','-c','user.email=audit@example.test','commit','-qm','fixture']]: subprocess.run(['git','-C',str(repo),*args],check=True,capture_output=True)
shim=root/'outage';shim.write_text('#!/bin/sh\nif [ "$1" = dispatch ] && [ "$2" = build ]; then echo first-build-outage >&2; exit 73; fi\nexec '+shlex.quote(str(B))+' "$@"\n');shim.chmod(0o755)
env=dict(os.environ,CLAUDE_CODE_SESSION_ID='')
def snap(): return {str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for p in root.rglob('*') if p.is_file()}
def call(route,ep,stage='backlog',extra=(),binary=shim):
 return subprocess.run([str(binary),'dispatch',route,'--workflow-dir','.','--entity-path',ep,'--stage',stage,*extra],cwd=root,env=env,text=True,input='- probe\n' if route=='build' else '',capture_output=True)
results=[];names={}
for slug in [short,long]:
 ep='state/'+slug+'.md';before=snap()
 first=call('build',ep,extra=['--checklist-file','-']);assert first.returncode==73
 q=call('name',ep);assert q.returncode==0,(q.stdout,q.stderr)
 again=call('build',ep,extra=['--checklist-file','-']);assert again.returncode==73
 assert before==snap(),'valid query or failing build mutated fixture bytes/files'
 names[slug]=q.stdout.strip();results.append({'case':'uncached-split-root-'+('short' if slug==short else 'long'),'build_exits':[first.returncode,again.returncode],'name':q.stdout.strip(),'snapshot_files':len(before),'unchanged':True})
assert names[short]==short+'-backlog';assert len(names[long])<=56
for label,ep,stage,extra in [('missing','state/missing.md','backlog',()),('undeclared','state/'+short+'.md','absent',()),('positional','state/'+short+'.md','backlog',('unexpected',)),('malformed-option','state/'+short+'.md','backlog',('--stamp',))]:
 before=snap();q=call('name',ep,stage,extra);assert q.returncode!=0 and not q.stdout,(label,q.stdout);assert before==snap();results.append({'case':label,'exit':q.returncode,'no_stdout':True,'unchanged':True})
alias='spacedock-ensign-'+short;(state/(alias+'.md')).write_text('---\ntitle: Alias\nstatus: backlog\n---\n')
before=snap();q=call('name','state/'+alias+'.md');assert q.returncode!=0 and not q.stdout and 'ambiguous' in q.stderr;assert before==snap();results.append({'case':'split-root-cross-generation-collision','exit':q.returncode,'unchanged':True})
(root/'README.md').write_text('not a workflow\n');before=snap();q=call('name','state/'+short+'.md');assert q.returncode!=0 and not q.stdout;assert before==snap();results.append({'case':'malformed-workflow','exit':q.returncode,'unchanged':True});(root/'README.md').write_text(readme)
# Only now permit successful build, after all uncached-query checks.
for slug in [short,long]:
 q=call('build','state/'+slug+'.md',extra=['--checklist-file','-'],binary=B)
 assert q.returncode==0,(q.stdout,q.stderr)
 out=json.loads(q.stdout);assert out['name']==names[slug]
 results.append({'case':'post-outage-build-equivalence','name':out['name'],'exact':True})
print(json.dumps({'fixture':str(root),'results':results},indent=2))
