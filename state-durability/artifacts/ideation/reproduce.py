import pathlib, subprocess, tempfile, json, sys
B=str(pathlib.Path(sys.argv[1]).resolve()); base=pathlib.Path(tempfile.mkdtemp(prefix='durability-repro-'))
rows=[]
def run(d,*args,stdin=None,expect=None):
 p=subprocess.run(args,cwd=d,input=stdin,text=True,capture_output=True)
 rows.append(dict(cwd=str(d),argv=args,rc=p.returncode,stdout=p.stdout,stderr=p.stderr))
 if expect is not None: assert p.returncode==expect,rows[-1]
 return p
def git(d,*a,expect=0): return run(d,'git',*a,expect=expect)
def write(p,s): p.parent.mkdir(parents=True,exist_ok=True);p.write_text(s)
def fixture(name,split=True):
 d=base/name;d.mkdir();git(d,'init','-q','-b','main');git(d,'config','user.name','Repro');git(d,'config','user.email','repro@example.test')
 w=d/'docs/dev';write(w/'README.md','---\ncommissioned-by: spacedock@1\nid-style: slug\nmerge: local\n'+('state: .spacedock-state\n' if split else '')+'stages:\n  states:\n    - name: implementation\n      initial: true\n    - name: validation\n      gate: true\n    - name: done\n      terminal: true\n---\n# Workflow\n')
 write(d/'.gitignore','docs/dev/.spacedock-state/\n');write(d/'delivery.txt','base\n');git(d,'add','.');git(d,'commit','-qm','seed')
 return d,w,w/'.spacedock-state' if split else w
def cli(d,w,*a,stdin=None,expect=None): return run(d,B,*a,'--workflow-dir',str(w),stdin=stdin,expect=expect)
# Missing checkout: table and JSON look empty, new silently creates ignored storage.
d,w,s=fixture('missing')
cli(d,w,'status',expect=0);missing_json=cli(d,w,'status','--json',expect=0).stdout
cli(d,w,'new','unsafe',stdin='---\nstatus: implementation\ntitle: Unsafe\n---\n',expect=0)
assert (s/'unsafe.md').exists();git(d,'check-ignore',str(s/'unsafe.md'));cli(d,w,'state','commit','unsafe',expect=1)
# Existing non-checkout directory equally accepts unsafe writes.
d,w,s=fixture('invalid');s.mkdir();cli(d,w,'status','--json',expect=0);cli(d,w,'new','unsafe',stdin='---\nstatus: implementation\n---\n',expect=0)
# Supported initialization restores a valid empty branch and safe filing.
d,w,s=fixture('valid');cli(d,w,'state','new',expect=0);valid_json=cli(d,w,'status','--json',expect=0).stdout
cli(d,w,'new','safe',stdin='---\nstatus: implementation\n---\n',expect=0);cli(d,w,'state','commit','safe',expect=0);git(s,'ls-files','--error-unmatch','safe.md')
# Retirement, both shapes, companion artifact, local-only and remote-backed.
for shape in ['flat','folder']:
 for remote in [False,True]:
  d,w,s=fixture('archive-'+shape+str(remote));cli(d,w,'state','new',expect=0)
  if remote:
   bare=base/(d.name+'.git');git(base,'init','-q','--bare',str(bare));git(d,'remote','add','origin',str(bare));git(s,'push','-q','-u','origin','spacedock-state/dev')
  write(s/('task.md' if shape=='flat' else 'task/index.md'),'---\nstatus: implementation\ntitle: Retire\n---\n')
  write(s/'task/evidence.txt','preserve me\n');cli(d,w,'state','commit','task',expect=0)
  before=git(s,'rev-parse','HEAD').stdout;cli(d,w,'status','--archive','task',expect=0)
  cli(d,w,'state','commit','task',expect=1);assert git(s,'rev-parse','HEAD').stdout==before
  assert (s/'_archive/task/evidence.txt').exists();git(s,'status','--porcelain')
# Actual approved gate then premature finalize, followed by real merge conflict/success.
for conflict in [True,False]:
 d,w,s=fixture('delivery-'+str(conflict));cli(d,w,'state','new',expect=0)
 wt=base/(d.name+'-worktree');git(d,'worktree','add','-q','-b','task',str(wt))
 write(wt/'delivery.txt','task\n');git(wt,'add','delivery.txt');git(wt,'commit','-qm','deliver task');taskhead=git(wt,'rev-parse','HEAD').stdout.strip()
 if conflict: write(d/'delivery.txt','trunk\n');git(d,'add','delivery.txt');git(d,'commit','-qm','conflict')
 write(s/'task.md','---\nid: task\nstatus: validation\ntitle: Task\nworktree: '+str(wt)+'\n---\n# Task\n');write(w/'gate-review.md','# Review\n');git(d,'add','docs/dev/gate-review.md');git(d,'commit','-qm','review')
 cli(d,w,'state','commit','task',expect=0)
 cli(d,w,'gate','prepare','task','--question','Advance?','--artifact',str(w/'gate-review.md'),'--summary','Reproduction.',expect=0)
 cli(d,w,'gate','record','task','--decision','approve','--actor','person:captain',expect=0)
 cli(d,w,'gate','consume','task',expect=0);assert 'state: pending' in (s/'task.md').read_text()
 cli(d,w,'merge','guard','task','--verdict','passed',expect=0)
 assert not (s/'task.md').exists();archive=s/'_archive/task.md';assert 'state: consumed' in archive.read_text();assert 'status: done' in archive.read_text()
 git(d,'merge','--no-ff','task','-m','deliver',expect=1 if conflict else 0)
 git(d,'merge-base','--is-ancestor',taskhead,'HEAD',expect=1 if conflict else 0)
 assert archive.exists()
 if conflict: git(d,'merge','--abort')
# Existing sentinel path supports delivery-first recovery without a new writer.
d,w,s=fixture('delivery-first');cli(d,w,'state','new',expect=0)
wt=base/'delivery-first-worktree';git(d,'worktree','add','-q','-b','task',str(wt));write(wt/'delivery.txt','delivered\n');git(wt,'add','delivery.txt');git(wt,'commit','-qm','task')
write(s/'task.md','---\nid: task\nstatus: validation\ntitle: Task\nworktree: '+str(wt)+'\n---\n# Task\n');write(w/'gate-review.md','# Review\n');git(d,'add','docs/dev/gate-review.md');git(d,'commit','-qm','review')
cli(d,w,'state','commit','task',expect=0);cli(d,w,'gate','prepare','task','--question','Advance?','--artifact',str(w/'gate-review.md'),'--summary','Ordered recovery.',expect=0);cli(d,w,'gate','record','task','--decision','approve','--actor','person:captain',expect=0)
git(d,'merge','--no-ff','task','-m','delivery',expect=0);mergehead=git(d,'rev-parse','HEAD').stdout.strip();taskhead=git(wt,'rev-parse','HEAD').stdout.strip();git(d,'merge-base','--is-ancestor',taskhead,'main');git(d,'merge-base','--is-ancestor',mergehead,'main')
cli(d,w,'status','--set','task','pr=local-merge:'+mergehead,expect=0);cli(d,w,'merge','guard','task','--verdict','passed',expect=0);assert 'state: consumed' in (s/'_archive/task.md').read_text()
head=git(s,'rev-parse','HEAD').stdout;cli(d,w,'merge','guard','task','--verdict','passed',expect=1);assert head==git(s,'rev-parse','HEAD').stdout
(base/'transcript.json').write_text(json.dumps(rows,indent=2));print(json.dumps({'base':str(base),'commands':len(rows),'missing_json':missing_json,'valid_empty_json':valid_json,'result':'All historical behaviors reproduced; archive refusal is already rc=1.'},indent=2))
