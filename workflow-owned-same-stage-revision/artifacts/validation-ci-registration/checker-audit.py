from pathlib import Path
import subprocess,json
root=Path.cwd(); registry=root/'docs/runtime-live-ci-registry.md';workflow=root/'.github/workflows/runtime-live-e2e.yml';scheduler=root/'internal/ensigncycle/scheduled_live_test.go'
original={p:p.read_text() for p in [registry,workflow,scheduler]}
out=root/'.audit-results';out.mkdir(exist_ok=True)
name='TestLiveCommonSameStageRevision';row='\t\t{"'+name+'", 1400, 1800, '+name+'},\n'
variants={
'baseline':{},
'missing-selector':{workflow:original[workflow].replace("-run '^TestLiveScheduled$'","-run '^Other$'",1)},
'missing-membership':{scheduler:original[scheduler].replace(row,'')},
'wrong-runtime-filter':{scheduler:original[scheduler].replace('for _, row := range rows {','for _, row := range rows {\n if runtime == "codex" && row.name == "'+name+'" { continue }',1)},
'commented-selector':{workflow:original[workflow].replace("-- -tags live -count=1 -timeout=90m -parallel=3 -run '^TestLiveScheduled$'", "-- -tags live -count=1 -timeout=90m -parallel=3 -run '^TestLiveScheduled$'",1)},
'unclassified':{registry:original[registry].replace('### `same-stage-revision`','## Targeted implementation proofs\n### `same-stage-revision`',1)},
'blank-exemption':{registry:original[registry].replace('- **Reason unselected:**', '- **No reason:**',1)},
}
# Comment out the actual first live scheduler invocation, retaining its text.
lines=original[workflow].splitlines(True)
for i,l in enumerate(lines):
 if '-tags live' in l and "-run '^TestLiveScheduled$'" in l:
  lines[i]=l[:len(l)-len(l.lstrip())]+'# '+l.lstrip();break
variants['commented-selector']={workflow:''.join(lines)}
results=[]
for label,changes in variants.items():
 for p,data in original.items():p.write_text(changes.get(p,data))
 r=subprocess.run(['go','test','./internal/contractlint','-run','^TestRuntimeLiveRegistryReconciliation$','-count=1'],text=True,capture_output=True)
 (out/(label+'.log')).write_text(r.stdout+r.stderr)
 results.append({'case':label,'exit':r.returncode,'expected':0 if label=='baseline' else 1})
 print(label,r.returncode,flush=True)
for p,data in original.items():p.write_text(data)
(out/'results.json').write_text(json.dumps(results,indent=2)+'\n')
