package status
import (
 "os"
 "os/exec"
 "path/filepath"
 "testing"
)
func TestValidationDeliveryAncestryConflictAndRecovery(t *testing.T) {
 root := stageFixture(t,"merge-local-workflow")
 gitOutput(t,root,"branch","-M","main")
 os.WriteFile(filepath.Join(root,"proof.txt"),[]byte("base"),0644)
 gitOutput(t,root,"add","."); gitOutput(t,root,"commit","-qm","proof base")
 old := gitOutput(t,root,"rev-parse","HEAD")
 wt := filepath.Join(root,"candidate")
 gitOutput(t,root,"worktree","add","-qb","audit-task",wt)
 os.WriteFile(filepath.Join(wt,"proof.txt"),[]byte("task"),0644)
 gitOutput(t,wt,"commit","-qam","task proof")
 task := gitOutput(t,wt,"rev-parse","HEAD")
 os.WriteFile(filepath.Join(root,"proof.txt"),[]byte("trunk"),0644)
 gitOutput(t,root,"commit","-qam","conflicting trunk")
 if err:=exec.Command("git","-C",root,"merge","--no-ff","audit-task").Run(); err==nil {t.Fatal("expected actual conflict")}
 for _,sentinel := range []string{"", "local-merge:"+old,"local-merge:"+task} {
  before:=gitOutput(t,root,"ls-files","--stage")
  if err:=verifyLocalDelivery(root,sentinel,"candidate");err==nil {t.Fatalf("accepted undelivered proof %s",sentinel)}
  if gitOutput(t,root,"ls-files","--stage")!=before {t.Fatal("proof changed conflicting index")}
 }
 gitOutput(t,root,"merge","--abort")
 gitOutput(t,root,"merge","--no-ff","-X","ours","-m","delivered ancestry","audit-task")
 sentinel := "local-merge:"+gitOutput(t,root,"rev-parse","HEAD")
 if err:=verifyLocalDelivery(root,sentinel,"candidate");err!=nil {t.Fatal(err)}
 if err:=verifyLocalDelivery(root,sentinel,"missing-worktree");err==nil {t.Fatal("missing worktree accepted")}
 gitOutput(t,root,"branch","-m","main","missing-trunk")
 if err:=verifyLocalDelivery(root,sentinel,"candidate");err==nil {t.Fatal("missing trunk accepted")}
}
