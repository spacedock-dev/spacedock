#!/usr/bin/env python3
"""Record one isolated installer and first-use run from an explicit recipe."""

import argparse
import datetime
import hashlib
import json
import os
import pathlib
import platform
import select
import signal
import subprocess
import sys
import time
import urllib.parse
import urllib.request
import uuid


OBSERVATIONS = ("local_render", "local_interaction", "local_result", "web_link",
                "host_authenticated", "decision_share")
WEB_CHECKS = ("room_html", "room_metadata", "artifact_bytes")
MILESTONES = (*OBSERVATIONS, *WEB_CHECKS, "web_readable", "full_install")
PREREQUISITES = {"local_interaction": "local_render", "local_result": "local_interaction",
                 "decision_share": "web_readable"}


def digest(path):
    with path.open("rb") as source:
        return hashlib.file_digest(source, "sha256").hexdigest()


def text_field(obj, key):
    value = obj[key]
    if not isinstance(value, str) or not value.strip():
        raise ValueError(f"{key} must be a nonempty string")
    return value


def load_recipe(path):
    recipe = json.loads(path.read_text())
    for field in ("origin_boundary", "route", "endpoint", "consent_mode"):
        text_field(recipe, field)
    if recipe["human_waits"] not in ("observed", "none", "opaque"):
        raise ValueError("human_waits must be observed, none or opaque")
    for field in ("sources", "versions", "initial_state", "caches", "unmeasured"):
        if not isinstance(recipe[field], dict) or not recipe[field]:
            raise ValueError(f"{field} must be a nonempty object")
    fixture = recipe["fixture"]
    text_field(fixture, "revision")
    if digest(path.parent / text_field(fixture, "path")) != text_field(fixture, "sha256"):
        raise ValueError("fixture SHA-256 does not match")
    requested, steps = recipe["requested_components"], recipe["steps"]
    if (not isinstance(requested, list) or any(not isinstance(s, str) or not s for s in requested)
            or len(set(requested)) != len(requested)):
        raise ValueError("requested_components must contain unique names")
    if not isinstance(steps, list) or not steps or "argv" not in steps[0]:
        raise ValueError("the first step must be the command at the origin boundary")
    if sum("web" in step for step in steps) > 1:
        raise ValueError("use one web readback per run")
    for step in steps:
        text_field(step, "name")
        kinds = set(step) & {"argv", "prompt", "web"}
        if len(kinds) != 1:
            raise ValueError("each step needs exactly one of argv, prompt or web")
        if "argv" in step and (not isinstance(step["argv"], list) or not step["argv"]
                               or any(not isinstance(arg, str) for arg in step["argv"])):
            raise ValueError("argv must be a nonempty list of strings")
        if "verifies" in step and ("argv" not in step or step["verifies"] not in requested or step.get("interactive")):
            raise ValueError("verifies must name a requested component and use a captured command")
        if "web" in step and text_field(step["web"]["artifact_bytes"], "sha256") != fixture["sha256"]:
            raise ValueError("Web artifact digest must match the verified fixture digest")
        if "prompt" in step:
            text_field(step, "prompt")
            if recipe["human_waits"] == "none":
                raise ValueError("a prompt requires human_waits observed or opaque")
        if step.get("timeout", 300) <= 0:
            raise ValueError("timeouts must be positive")
    return recipe


def record_observation(entry, evidence):
    entry["ns"] = time.monotonic_ns()
    evidence = evidence.resolve(strict=True)
    work = pathlib.Path(os.environ["INSTALL_BENCH_WORK"])
    if not evidence.is_relative_to(work) or not evidence.is_file():
        raise ValueError("observation evidence must be a file inside the run's work directory")
    entry["evidence"] = str(evidence.relative_to(work))
    with open(os.environ["INSTALL_BENCH_EVENTS"], "a", encoding="utf-8") as out:
        out.write(json.dumps(entry) + "\n")


def read_observations(event_log, start, offset, mark, waits, human_waits):
    if not event_log.exists():
        return None
    pending = {}
    declined = False
    for line in event_log.read_text().splitlines():
        event = json.loads(line)
        at = offset(event["ns"])
        if not start <= at <= offset():
            raise ValueError("observation is outside the command interval")
        name, evidence = event["name"], "work/" + event["evidence"]
        boundary = event.get("boundary")
        if boundary is None:
            if not declined:
                mark(name, at, evidence)
        elif human_waits == "none":
            raise ValueError("observed a human wait in a run declared to have none")
        elif boundary == "visible":
            if name in pending:
                raise ValueError("prompt is already visible")
            wait = {"name": name, "visible_ms": at, "answer_ms": None, "visible_evidence": evidence}
            waits.append(wait)
            pending[name] = wait
        else:
            wait = pending.pop(name)
            if at < wait["visible_ms"]:
                raise ValueError("response precedes the prompt")
            wait.update(answer_ms=at, result=boundary, answer_evidence=evidence)
            if boundary == "declined":
                declined = True
    if declined:
        return "declined"
    if pending:
        return "prompt unanswered"
    return None


def run_command(step, work, env, logs):
    interactive = step.get("interactive", False)
    with logs.with_suffix(".stdout").open("wb") as stdout, logs.with_suffix(".stderr").open("wb") as stderr:
        child = subprocess.Popen(
            step["argv"], cwd=work, env=env, start_new_session=not interactive,
            stdin=None if interactive else subprocess.DEVNULL,
            stdout=None if interactive else stdout, stderr=None if interactive else stderr)
        try:
            return child.wait(timeout=step.get("timeout", 300))
        except (subprocess.TimeoutExpired, KeyboardInterrupt):
            if child.poll() is None:
                try:
                    if interactive:
                        child.kill()
                    else:
                        os.killpg(child.pid, signal.SIGKILL)
                except ProcessLookupError:
                    pass
            child.wait()
            raise


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, *args, **kwargs):
        raise ValueError("readback redirected away from the exact intended endpoint")


def read_web(step, offset, mark, destination):
    web = step["web"]
    scheme, netloc = urllib.parse.urlsplit(text_field(web, "origin"))[:2]
    if scheme not in ("http", "https") or not netloc or "@" in netloc:
        raise ValueError("web origin must be an HTTP(S) origin without credentials")
    client = urllib.request.build_opener(urllib.request.ProxyHandler({}), NoRedirect())
    for name in WEB_CHECKS:
        check = web[name]
        url = urllib.parse.urlsplit(text_field(check, "url"))
        if (url.scheme, url.netloc) != (scheme, netloc):
            raise ValueError(f"{name} is outside the intended origin")
        with client.open(url.geturl(), timeout=step.get("timeout", 15)) as response:
            if response.status != 200:
                raise ValueError(f"{name} did not return HTTP 200")
            content_type = response.headers.get_content_type()
            data = response.read(32 * 1024 * 1024 + 1)
        sha256 = hashlib.sha256(data).hexdigest()
        if len(data) > 32 * 1024 * 1024:
            raise ValueError(f"{name} response exceeded the size limit")
        if name == "artifact_bytes" or "sha256" in check:
            if sha256 != text_field(check, "sha256"):
                raise ValueError(f"{name} did not return the expected bytes")
        if name == "room_html":
            if content_type != "text/html":
                raise ValueError("Room response is not HTML")
            if "sha256" not in check and text_field(check, "contains").encode() not in data:
                raise ValueError("HTML did not identify the expected Room")
        if name == "room_metadata":
            metadata = json.loads(data)
            if "sha256" not in check:
                fields = check["fields"]
                if not isinstance(fields, dict) or not fields or not isinstance(metadata, dict):
                    raise ValueError("metadata field verification needs nonempty expected fields")
                if any(key not in metadata or metadata[key] != value for key, value in fields.items()):
                    raise ValueError("metadata did not match the expected record")
        receipt = destination / f"{name}.json"
        receipt.write_text(json.dumps({"status": 200, "sha256": sha256, "bytes": len(data)}))
        mark(name, offset(), receipt.name)
    mark("web_readable", offset(), "all three anonymous GETs verified")


def run_recipe(recipe_path, destination):
    recipe = load_recipe(recipe_path)
    destination.mkdir(mode=0o700)
    work = destination / "work"
    for directory in (work, work / "home", work / "bin", work / "cache",
                      work / "data", work / "config", work / "tmp"):
        directory.mkdir()
    env = dict(
        PATH=f"{work / 'bin'}:/usr/bin:/bin:/usr/sbin:/sbin", HOME=str(work / "home"),
        XDG_CACHE_HOME=str(work / "cache"), XDG_DATA_HOME=str(work / "data"),
        XDG_CONFIG_HOME=str(work / "config"), TMPDIR=str(work / "tmp"),
        TERM=os.getenv("TERM", "dumb"), LC_ALL="C")
    env.update(recipe.get("env", {}))
    env["INSTALL_BENCH_WORK"] = str(work)
    requested = set(recipe["requested_components"])
    milestones = {name: {"elapsed_ms": None, "reason": "no verified evidence observed"} for name in MILESTONES}
    if not requested:
        milestones["full_install"]["reason"] = "no installation requested"
    operations, verified, waits = [], {}, []

    def offset(ns=None):
        return (time.monotonic_ns() - zero if ns is None else ns - zero) // 1_000_000

    def mark(name, at, evidence):
        prerequisite = PREREQUISITES.get(name)
        if prerequisite:
            prior = milestones[prerequisite]["elapsed_ms"]
            if prior is None or at < prior:
                raise ValueError(f"{name} requires earlier verified {prerequisite}")
        if milestones[name]["elapsed_ms"] is None:
            milestones[name] = {"elapsed_ms": at, "evidence": evidence}

    record = {
        "schema": 1, "run_id": str(uuid.uuid4()), "platform": f"{platform.system()} {platform.release()} {platform.machine()}",
        "clock": "monotonic elapsed milliseconds", "recipe_sha256": digest(recipe_path),
        "operations": operations, "milestones": milestones, "waits": waits,
        "verified_components": verified,
    }
    for key in ("origin_boundary", "route", "endpoint", "consent_mode", "sources", "versions",
                "initial_state", "caches", "unmeasured", "requested_components", "human_waits"):
        record[key] = recipe[key]
    record["fixture"] = {key: recipe["fixture"][key] for key in ("revision", "sha256")}
    record["started_utc"] = datetime.datetime.now(datetime.timezone.utc).isoformat()
    zero = time.monotonic_ns()
    exit_code = 0
    for index, step in enumerate(recipe["steps"], 1):
        logs = destination / f"{index:02d}"
        operation = {"name": step["name"], "start_ms": offset(), "exit": None}
        operations.append(operation)
        try:
            if "prompt" in step:
                print(step["prompt"] + " [y/N] ", end="", flush=True)
                wait = {"name": step["name"], "visible_ms": offset(), "answer_ms": None}
                waits.append(wait)
                if not select.select([sys.stdin], [], [], step.get("timeout", 300))[0]:
                    raise TimeoutError("no human response before deadline")
                answer = sys.stdin.readline()
                if not answer:
                    raise EOFError("prompt input closed without a response")
                granted = answer.strip().lower() in ("y", "yes")
                wait.update(answer_ms=offset(), result="granted" if granted else "declined")
                operation["result"] = wait["result"]
                exit_code = 0 if granted else 1
            elif "web" in step:
                read_web(step, offset, mark, destination)
                operation["result"] = "done"
            else:
                event_log = logs.with_suffix(".events")
                env["INSTALL_BENCH_EVENTS"] = str(event_log)
                operation["evidence"] = logs.name + ".{stdout,stderr}"
                if step.get("interactive"):
                    operation["evidence"] = "interactive terminal; observer supplies evidence"
                try:
                    code = run_command(step, work, env, logs)
                    operation["exit"] = code
                finally:
                    wait_stop = read_observations(
                        event_log, operation["start_ms"], offset, mark, waits, recipe["human_waits"])
                operation["result"] = wait_stop or ("done" if code == 0 else "interrupted" if code < 0 else "failed")
                if code == 0 and not wait_stop and "verifies" in step:
                    verified[step["verifies"]] = logs.name + ".stdout"
                    if requested == verified.keys():
                        mark("full_install", offset(), "all requested component checks passed")
                exit_code = 128 - code if code < 0 else code
                if wait_stop and exit_code == 0:
                    exit_code = 1
        except KeyboardInterrupt:
            operation["result"], exit_code = "interrupted", 130
        except EOFError:
            operation["result"], exit_code = "prompt unanswered", 1
        except (subprocess.TimeoutExpired, TimeoutError):
            operation["result"], exit_code = "timed out", 124
        except (OSError, KeyError, TypeError, ValueError) as exc:
            operation.update(result="failed", detail=str(exc))
            exit_code = 1
        operation["end_ms"] = offset()
        if exit_code and step.get("allow_failure") and operation["result"] == "failed":
            operation["result"], exit_code = "warning", 0
        if exit_code:
            break
    intervals = sorted((item["visible_ms"], item["answer_ms"]) for item in waits if item["answer_ms"] is not None)
    observed_wait = last = 0
    for start, end in intervals:
        observed_wait += max(0, end - max(start, last))
        last = max(last, end)
    record.update(exit=exit_code, total_ms=offset(), observed_wait_ms=observed_wait)
    if recipe["human_waits"] == "opaque" or any(item["answer_ms"] is None for item in waits):
        record["active_ms"] = None
    else:
        record["active_ms"] = record["total_ms"] - observed_wait
    if exit_code:
        for entry in milestones.values():
            if entry.get("reason") == "no verified evidence observed":
                entry["reason"] = "no verified evidence before run stopped"
    (destination / "run.json").write_text(json.dumps(record, indent=2) + "\n")
    print(f"Run record: {destination / 'run.json'}")
    return exit_code


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest="command", required=True)
    run = sub.add_parser("run", help="execute a recipe into a new private directory")
    run.add_argument("recipe", type=pathlib.Path)
    run.add_argument("output", type=pathlib.Path)
    observe = sub.add_parser("observe", help="timestamp verified evidence from a recipe child")
    observe.add_argument("event", choices=OBSERVATIONS)
    observe.add_argument("evidence", type=pathlib.Path)
    wait = sub.add_parser("observe-wait", help="timestamp a human prompt or response inside a recipe child")
    wait.add_argument("name")
    wait.add_argument("boundary", choices=("visible", "granted", "declined"))
    wait.add_argument("evidence", type=pathlib.Path)
    args = parser.parse_args()
    os.umask(0o077)
    try:
        if args.command == "run":
            return run_recipe(args.recipe.resolve(), args.output.resolve())
        entry = ({"name": args.event} if args.command == "observe"
                 else {"name": args.name, "boundary": args.boundary})
        record_observation(entry, args.evidence)
        return 0
    except (OSError, KeyError, TypeError, ValueError) as exc:
        print(f"install-benchmark: {exc}", file=sys.stderr)
        return 2


if __name__ == "__main__":
    sys.exit(main())
