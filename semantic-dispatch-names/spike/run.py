#!/usr/bin/env python3
"""Requires Python 3, Git and Go; run from a source checkout containing BASE."""
import io, os, pathlib, subprocess, tarfile, tempfile
BASE = "438053493838dc70c9478b3d991309d566783e85"
HERE = pathlib.Path(__file__).resolve().parent
with tempfile.TemporaryDirectory(prefix="semantic-names-spike-") as tmp:
    root = pathlib.Path(tmp)
    archive = subprocess.check_output(["git", "archive", BASE])
    tarfile.open(fileobj=io.BytesIO(archive)).extractall(root, filter="data")
    (root / "internal/dispatch/semantic_spike_test.go").write_bytes((HERE / "semantic_spike_test.go.txt").read_bytes())
    cmd = ["go", "test", "./internal/dispatch", "-run", "^TestSemanticSpike$", "-v", "-count=1"]
    subprocess.run(cmd, cwd=root, env=dict(os.environ, SEMANTIC_SPIKE_LEGACY="1"), check=True)
    result = subprocess.run(cmd, cwd=root, env=dict(os.environ, SEMANTIC_SPIKE_LEGACY="0"))
    if result.returncode == 0:
        raise SystemExit("baseline unexpectedly already satisfies semantic names")
    subprocess.run(["git", "apply", str(HERE / "prototype.patch")], cwd=root, check=True)
    subprocess.run(cmd, cwd=root, env=dict(os.environ, SEMANTIC_SPIKE_LEGACY="0"), check=True)
