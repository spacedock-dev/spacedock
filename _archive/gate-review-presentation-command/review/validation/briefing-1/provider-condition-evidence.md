# Pinned sibling provider evidence

- Repository: `/Users/clkao/git/spacedock-research/spacedock-subspace`
- Pinned revision: `20694cd3bdf0a7d43da630adb58812ce9ef96468`
- Parent: `0de94dd14ee0f973afca30cc8fb67afab5868bb5`
- Subject: `roadmap: pull briefing presentation into skill wiring`
- Author date: `2026-07-22T22:25:30+08:00`

The evidence reads the committed tree at the pinned revision, not unrelated untracked working-copy files.

`git ls-tree -r --name-only 20694cd3bdf0a7d43da630adb58812ce9ef96468 -- scripts .github` returned exactly:

```text
.github/workflows/publish-beta-readme.yml
scripts/tests/local-plugin-marketplace-test.sh
scripts/tests/publish-beta-readme-test.sh
scripts/tests/subspace-r-async-helper-test.sh
scripts/tests/subspace-r-cmux-atomic-real-smoke.sh
scripts/tests/subspace-r-contract-test.sh
scripts/tests/subspace-r-tmux-recipe-test.sh
scripts/tests/subspace-r-zellij-helper-test.sh
scripts/tests/subspace-r-zellij-real-smoke.sh
scripts/tests/subspace-tui-agent-contract-fixture
scripts/tests/subspace-tui-agent-contract-fixture-test.sh
```

`git grep` at the pinned revision over `scripts` and `.github` found zero matches for `review-local-zellij`, `room-resident`, `probe-first`, `caller-owned result`, `12-fixture`, or `drive suite`.

Therefore the pinned sibling revision does not carry the named hardened override script or a committed 12-fixture CI drive suite. It cannot prove xb AC-1, AC-2, AC-3, or AC-5 and does not satisfy release eligibility.
