# Installer and first-use benchmark

Run `python3 scripts/install-benchmark.py run RECIPE.json OUTPUT` with Python
3.11 or later. Use a new output directory for every trial. The runner creates an
empty home, bin, cache, data, config, temp and working directory. Its commands
receive a minimal environment plus the recipe's explicit overrides. Use trusted
command arrays. If a tool uses system-wide state, run it in a disposable VM,
especially when Homebrew is part of the trial. Keep recipes and raw output local.

Start with a synthetic review document and known revision. Before the timed run,
record versions and expected SHA-256 values for the installer, client, release
artifact, fixture and other pinned inputs as applicable. Record the Homebrew
version and tap/recipe revision when using Brew. Stage only the inputs that the
selected cache condition calls for. Shared system and CDN caches are uncontrolled.

## Recipe

This loopback example measures acquisition and checks a binary-only install. It
assumes a server supplying a pinned installer, release archive and checksums, and
a local `verify-spacedock-version` executable that compares the reported version
with its expected value. Replace the example refs, paths, hash and version before
running. Relative command paths are resolved from the new work directory. The
fixture path is relative to the recipe file. A successful verifier, rather than
the installer's exit code, establishes the requested component.

```json
{
  "origin_boundary": "first download command for the pinned installer",
  "route": "native",
  "endpoint": "loopback",
  "consent_mode": "synthetic fixture preauthorized locally",
  "human_waits": "none",
  "sources": {"installer": "<commit and SHA-256>", "archive": "<release and SHA-256>"},
  "versions": {"expected_spacedock": "<version>"},
  "initial_state": {"claude": "absent", "codex": "absent", "pi": "absent", "brew": "absent"},
  "caches": {"runner": "fresh", "server": "uncontrolled"},
  "unmeasured": {"shell internal phases": "download, checksum and unpack are inside the shell step"},
  "fixture": {"path": "review.md", "revision": "<revision>", "sha256": "<fixture SHA-256>"},
  "requested_components": ["spacedock binary"],
  "steps": [
    {"name": "Installer download", "argv": ["curl", "-fsS", "http://127.0.0.1:8080/install.sh", "-o", "install.sh"]},
    {"name": "Archive download, checksum and binary installation", "argv": ["sh", "install.sh"]},
    {"name": "Confirm the installed version", "argv": ["/absolute/path/verify-spacedock-version", "bin/spacedock", "<version>"], "verifies": "spacedock binary"}
  ],
  "env": {
    "SPACEDOCK_INSTALL_FROM": "http://127.0.0.1:8080",
    "SPACEDOCK_INSTALL_VERSION": "<version>",
    "SPACEDOCK_INSTALL_DIR": "bin",
    "SPACEDOCK_EXTRAS": "none"
  }
}
```

The first command sets the outer clock. If trust, bootstrap or another command
can trigger the first fetch, make that the first step and record its timing
boundary accordingly. Preserve the guide's actual order, including a second
host install if the guide runs one. Use separate steps for download, release
lookup, checksum, unpack/link, extras, Brew bootstrap, trust, tap and cask
transaction wherever they are independently observable. Name indivisible
commands as such in `unmeasured`. Separate `brew fetch` trials cannot fill
in timings inside a cask transaction. Use `spacedock install --verbose` for a
host step. Its captured stderr retains the nested host clock.

Set `verifies` on a captured command that checks the expected version and
availability, or the host plugin's inventory, version, enablement and
compatibility. List every requested program and plugin that needs a check in
`requested_components`. An empty list describes a trial without a full-install
milestone, such as a standalone review reader. Use `allow_failure: true` on
an explicitly tolerated command. It records a warning without granting
verification. Give long operations an explicit timeout in seconds.

## First local use and Web readback

A recipe command can run the real reader in a pseudo-terminal or use
`"interactive": true` to inherit the current terminal. From that command or its
observer, call the following entry point as each event occurs. The observer
inherits the benchmark environment and writes evidence under the run's work
directory.

```sh
python3 /path/to/scripts/install-benchmark.py observe local_render render.txt
python3 /path/to/scripts/install-benchmark.py observe local_interaction key-response.txt
python3 /path/to/scripts/install-benchmark.py observe local_result neutral-result.json
```

Verify the fixture revision on screen before recording render, a changed view
after a real keypress before recording first interaction, and a matching neutral
result before recording it. The observer supplies application-specific checks.
A launch, exit zero or key sent to a terminal does not establish those events.
Use `observe web_link` when a real link prints. An existing, consented Web
share can be checked with a `web` step in the same recipe. It performs three
anonymous GETs against the declared origin and only records `web_readable`
after all three succeed. The expected artifact digest must equal the digest of
the verified fixture file. Choose an exact-byte fixture that the share will serve:

```json
{
  "name": "Read the synthetic Room anonymously",
  "web": {
    "origin": "http://127.0.0.1:8080",
    "room_html": {"url": "http://127.0.0.1:8080/room/example", "contains": "room-example"},
    "room_metadata": {"url": "http://127.0.0.1:8080/metadata/example", "fields": {"revision": "r1"}},
    "artifact_bytes": {"url": "http://127.0.0.1:8080/artifacts/example", "sha256": "<fixture SHA-256>"}
  }
}
```

HTML and metadata can also be verified by expected SHA-256. Keep loopback and
public readback in different trials. The readback step makes no uploads. If the
trial includes a separately authorized decision share or authenticated host
action, record `decision_share` only with pairing/auth evidence and
`host_authenticated` only after a successful authenticated action. Link output,
a preview and sign-in UI are earlier boundaries.

A step with `"prompt": "…"` and a name records when a trust or consent decision
is requested and when it was granted or declined. A timeout or interruption
leaves the answer timestamp blank. For a prompt inside a child, have the session
observer call `observe-wait brew-trust visible prompt.txt` when it sees the
prompt and `observe-wait brew-trust granted response.txt` (or `declined`) when
it sees the response. Use the same name and write both evidence files in the work
directory. A recorded decline or unanswered prompt stops later recipe steps.
If either boundary cannot be observed, set `"human_waits": "opaque"`
and name it in `unmeasured`. The runner does not report active time with opaque
or unfinished waits. A prompt decision, including a denial, remains distinct
from completing sign-in.

## Compare runs

Read `OUTPUT/run.json` alongside the referenced evidence. Offsets are relative
to the start of that run, with milliseconds rounded down. Missing milestone
times are null, including later milestones after a failed, declined or
interrupted run. Keep the terminal status and last observed boundary. A process
killed without cleanup may not leave a final record. Compare fresh sequential
profiles with the same requested components, fixture, provenance, environment
and cache/consent conditions. Report each run and its range or median. Retain
failed runs as unfinished observations.
