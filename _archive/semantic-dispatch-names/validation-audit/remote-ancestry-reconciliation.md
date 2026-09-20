# Owned remote-ancestry reconciliation

- Approved prior head: `bd7a140f23a070dc516aab2f19bb05f011c83e9b`.
- Recorded remote-tracking tip: `8d54f7b2047074b231e1f7adebbf817e47ef9418`.
- Current local parent: `72f3493c97182b607b10a4d7804a1857f2e0ea13`.
- Resulting naming head: `756b8abd81a79075521efb47e066d4cbf615d80d`.
- Exact tree of both approved prior and resulting heads: `95530042677d463abbea25ec611ce20417133397`.

The clean owned branch first recorded the remote tip with `git reset --keep`, then replayed its five naming commits onto the current local parent. The expected registry conflict was reconciled to the already-approved corresponding registry bytes. The approved promotion commit was then replayed. No other branch changed.

`git diff --exit-code bd7a140f2 HEAD` passed with no differences. `git range-diff 399f775fa..bd7a140f2 72f3493c9..HEAD` reports all six commits patch-equivalent (`=`). The current local parent is an ancestor of the result, and the worktree is clean.

The owned branch reflog explicitly retains remote tip `8d54f7b2047074b231e1f7adebbf817e47ef9418` as the reset entry preceding the rebase and cherry-pick. This preserves the remote-inclusion evidence used by bare `--force-with-lease --force-if-includes`; neither protection was bypassed. Publication remains with FO and must recheck the current remote state.

No product behavior changed. No tests, push or CI ran for the identical tree. Entity frontmatter and approved report history remain untouched.
