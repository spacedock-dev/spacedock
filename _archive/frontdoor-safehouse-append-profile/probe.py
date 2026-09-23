"""Harmless Safehouse policy probe; no model, sandbox execution, or config writes.

Usage: python3 probe.py /path/to/agent-safehouse/dist/safehouse.sh
Source: https://github.com/eugene1g/agent-safehouse
Recorded source commit: e376993ee8e15c4e4b3aa3a2ee282f15f6e3c680 (0.12.0).
Requires Python 3, /bin/bash, and the upstream standalone Safehouse script.
"""
import os
import pathlib
import subprocess
import sys
import tempfile

script = str(pathlib.Path(sys.argv[1]).resolve())
env = {k: v for k, v in os.environ.items() if not k.startswith("SAFEHOUSE_")}
with tempfile.TemporaryDirectory(prefix="append-profile-probe-") as tmp:
    base = pathlib.Path(tmp)
    project = base / "project"
    project.mkdir()
    (project / ".safehouse").write_text("append-profile=project.sb\n")
    (project / "project.sb").write_text(";; PROBE_PROJECT\n")
    first = "a space,colon:=literal.sb"
    (base / first).write_text(";; PROBE_A\n")
    (base / "b.sb").write_text(";; PROBE_B\n")

    def run(*args):
        return subprocess.run(
            ["/bin/bash", script, *args], cwd=base, env=env,
            text=True, capture_output=True, timeout=30,
        )

    help_result = run("--help")
    assert help_result.returncode == 0, help_result.stderr
    assert "--append-profile=PATH" in help_result.stdout
    print(help_result.stdout.splitlines()[0])
    print("help: --append-profile PATH and --append-profile=PATH available")
    result = run(
        "--stdout", "--trust-workdir-config", "--workdir=" + str(project),
        "--append-profile=" + first, "--append-profile", "b.sb",
        "--append-profile=b.sb",
    )
    assert result.returncode == 0, result.stderr
    markers = [line for line in result.stdout.splitlines() if "PROBE_" in line]
    assert markers == [";; PROBE_PROJECT", ";; PROBE_A", ";; PROBE_B", ";; PROBE_B"], markers
    print("policy: project,A,B,B; CLI paths use cwd; punctuation survives")
    for args, error in [
        (["--append-profile="], "Appended profile path cannot be empty."),
        (["--append-profile=missing.sb"], "Appended profile path does not exist: missing.sb"),
        (["--append-profile=project"], "Appended profile path is not a regular file: project"),
        (["--append-profile"], "Missing value for --append-profile"),
    ]:
        result = run("--stdout", *args)
        assert result.returncode == 1, (args, result.returncode)
        assert error in result.stderr, result.stderr
        print(f"error: {args!r}: exit 1: {error}")
