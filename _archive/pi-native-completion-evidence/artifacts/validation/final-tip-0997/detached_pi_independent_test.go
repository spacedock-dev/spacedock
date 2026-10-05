package ensigncycle

import (
 "encoding/json"
 "path/filepath"
 "strings"
 "testing"
)

func TestDetachedPiFullRetainedObservers(t *testing.T) {
 for _, scenario := range []string{"default-headless-gate-stop", "self-feedback/plain"} {
  t.Run(scenario,func(t *testing.T){
   root := filepath.Join("/tmp/spacedock-tip35149242496-pi/live-artifacts/pi/pi-common",scenario)
   paths,err := filepath.Glob(filepath.Join(root,"sessions","*.jsonl")); if err!=nil || len(paths)!=1 {t.Fatal("missing assigned retained input")}
   stream := readFile(t,filepath.Join(root,"pi-stdout.txt"))+"\n"+readFile(t,filepath.Join(root,"pi-stderr.txt"))+"\n"+readFile(t,paths[0])
   routes,_,err := piRejectionRoutes(stream,root); if err!=nil {t.Fatal(err)}
   if scenario=="default-headless-gate-stop" {
    if len(routes)!=4 || routes[0].index!=65 || routes[1].index!=66 || routes[2].index!=80 || routes[3].index!=81 {t.Fatalf("full sync boundaries %v",routes)}
    if err:=assertImplementationWorkerLifecycle(stream,"## Stage Report: implementation\n- DONE: report\n",root);err!=nil {t.Fatal(err)}
    if err:=assertImplementationWorkerLifecycle(stream,"",root);err==nil {t.Fatal("missing report accepted")}
   } else {
    if len(routes)!=2 || routes[0].index!=112 || routes[1].index!=116 || routes[0].target!=routes[1].target {t.Fatalf("full async boundaries %v",routes)}
    if err:=assertSameStageWorkers(routes,false);err!=nil {t.Fatal(err)}
    missing,_,err:=piRejectionRoutes(stream); if err==nil && assertSameStageWorkers(missing,false)==nil {t.Fatal("lost artifact propagation silently accepted")}
    _,_,err=piRejectionRoutes(stream,t.TempDir());if err==nil {t.Fatal("wrong retained artifact root accepted")}
   }
  })
 }
}

func TestDetachedPiSynchronousCrossEvidence(t *testing.T) {
 for _, name:=range []string{"retained-parent-cwd","unicode-assignment","terminal-not-last","stale-second-dispatch","cross-worker-call","cross-worker-run","after-boundary-with-false-clock","eof-no-newline"} {
  t.Run(name,func(t *testing.T){
   stream,root,path:=capturedPiTipCompletion(t,"sync"); rows:=strings.Split(stream,"\n");child:=readFile(t,path)
   var result piSessionRecord; if err:=json.Unmarshal([]byte(rows[66]),&result);err!=nil {t.Fatal(err)}
   oldRun:=result.Message.Details.RunID; oldCall:=result.Message.ToolCallID
   switch name {
   case "retained-parent-cwd":
    p:=filepath.Join(root,"sessions",filepath.Base(result.Message.Details.Mission.OwnerSessionID));writeFile(t,p,strings.Replace(readFile(t,p),"2958864088/004","2958864088/005",1))
   case "unicode-assignment":child=strings.Replace(child,"Read /tmp/","Read\u200b /tmp/",1)
   case "terminal-not-last":child+="{\"type\":\"message\",\"message\":{\"role\":\"user\",\"content\":[]}}\n"
   case "stale-second-dispatch","cross-worker-call","cross-worker-run":
    // A second real dispatch/result pair with fresh IDs and a separate retained child.
    secondSpawn:=strings.ReplaceAll(rows[65],oldCall,"independent-second-call")
    secondSpawn=strings.ReplaceAll(secondSpawn,"22:07:12.007Z","22:08:00.000Z")
    secondResult:=strings.ReplaceAll(strings.ReplaceAll(rows[66],oldCall,"independent-second-call"),oldRun,"independent-second-run")
    secondResult=strings.ReplaceAll(secondResult,"22:07:58.319Z","22:08:10.000Z")
    secondChild:=strings.ReplaceAll(child,oldRun,"independent-second-run")
    if name!="stale-second-dispatch" {secondChild=strings.ReplaceAll(strings.ReplaceAll(secondChild,"22:07:12.070Z","22:08:01.000Z"),"22:07:58.294Z","22:08:09.000Z")}
    if name=="cross-worker-call" {secondResult=strings.ReplaceAll(secondResult,"independent-second-call",oldCall)}
    if name=="cross-worker-run" {secondResult=strings.ReplaceAll(secondResult,"independent-second-run",oldRun)}
    writeFile(t,strings.ReplaceAll(path,oldRun,"independent-second-run"),secondChild)
    rows[70],rows[71]=secondSpawn,secondResult
   case "after-boundary-with-false-clock":rows=append(rows,rows[66]);rows[66]=""
   case "eof-no-newline":child=strings.TrimRight(child,"\n")
   }
   writeFile(t,path,child);stream=strings.Join(rows,"\n")
   err:=assertImplementationWorkerLifecycle(stream,"## Stage Report: implementation\n- DONE: report\n",root)
   if name=="eof-no-newline" {if err!=nil {t.Fatal(err)};return}
   if err==nil {t.Fatal("invalid evidence accepted at lifecycle boundary")}
   routes,_,err:=piRejectionRoutes(stream,root);if err==nil && assertSameStageWorkers(routes,false)==nil {t.Fatal("invalid evidence accepted at route boundary")}
  })
 }
}
