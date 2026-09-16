package cli

import (
 "bytes"
 "os"
 "path/filepath"
 "testing"
)

func TestValidationRetirementRollbackExactCompanions(t *testing.T) {
 for _, folder := range []bool{false,true} {
  name := "flat"; if folder { name = "folder" }
  t.Run(name,func(t *testing.T) {
   _,wf,_,_ := twoHostStateWorkflow(t)
   state := filepath.Join(wf,".spacedock-state")
   git(t,state,"remote","remove","origin")
   rel := "retire.md"; if folder { rel = "retire/index.md" }
   original := "---\nstatus: ideation\n---\nEOF without newline"
   writeFileWithDirs(t,filepath.Join(state,rel),original)
   artifact := filepath.Join(state,"retire","artifact.bin")
   writeFileWithDirs(t,artifact,"\x00\xff\nretained")
   deleted := filepath.Join(state,"retire","deleted.txt")
   writeFileWithDirs(t,deleted,"tracked deletion")
   git(t,state,"add","."); git(t,state,"commit","-qm","audit seed")
   if err := os.Remove(deleted); err != nil {t.Fatal(err)}
   writeFile(t,filepath.Join(state,"sibling.md"),"staged")
   git(t,state,"add","sibling.md")
   writeFile(t,filepath.Join(state,"sibling.md"),"unstaged")
   sibling := git(t,state,"ls-files","--stage","sibling.md")
   refs := git(t,state,"show-ref")
   hooks := t.TempDir()
   writeFile(t,filepath.Join(hooks,"pre-commit"),"#!/bin/sh\nexit 1\n")
   os.Chmod(filepath.Join(hooks,"pre-commit"),0755)
   git(t,state,"config","core.hooksPath",hooks)
   c,o,e := terminalInvoke(t,wf,"status","--workflow-dir",wf,"--archive","retire","--json")
   if c != 1 { t.Fatalf("commit refusal: %d %s %s",c,o,e) }
   b,err := os.ReadFile(filepath.Join(state,rel)); if err != nil || string(b)!=original {t.Fatalf("entity restored: %q %v",b,err)}
   b,err = os.ReadFile(artifact); if err != nil || !bytes.Equal(b,[]byte("\x00\xff\nretained")) {t.Fatalf("artifact restored: %q %v",b,err)}
   if git(t,state,"show-ref")!=refs || git(t,state,"ls-files","--stage","sibling.md")!=sibling {t.Fatal("refs or sibling index changed")}
   b,_=os.ReadFile(filepath.Join(state,"sibling.md")); if string(b)!="unstaged" {t.Fatal("sibling worktree changed")}
   if _,err=os.Stat(deleted); !os.IsNotExist(err) {t.Fatal("deleted companion resurrected")}
   git(t,state,"config","--unset","core.hooksPath")
   c,o,e = terminalInvoke(t,wf,"status","--workflow-dir",wf,"--archive","retire","--json")
   if c!=0 {t.Fatalf("recovery: %d %s %s",c,o,e)}
   if git(t,state,"ls-files","--stage","sibling.md")!=sibling {t.Fatal("recovery changed sibling")}
   if got:=git(t,state,"ls-tree","-r","--name-only","HEAD","retire","retire.md","_archive/retire/deleted.txt"); got!="" {t.Fatalf("incomplete committed move: %s",got)}
  })
 }
}
