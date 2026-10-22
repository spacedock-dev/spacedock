#!/usr/bin/env python3
"""Drive a real released Review TUI, then verify its result reaches Claude."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import secrets
import shlex
import subprocess
import sys
import threading
import time
import traceback

TARGET = os.environ.get('REVIEW_TEST_TARGET', 'Keep the launch goal visible.')
TIMEOUT = int(os.environ.get('REVIEW_TEST_TIMEOUT', '480'))


def run(argv, *, check=True, timeout=15):
    p = subprocess.run(argv, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, timeout=timeout)
    if check and p.returncode:
        raise RuntimeError(f'{argv[0]} failed ({p.returncode}): {p.stderr.strip()}')
    return p


def tmux(*args, **kwargs):
    return run(['tmux', *args], **kwargs)


def write_json(path, value):
    Path(path).write_text(json.dumps(value, indent=2) + '\n')


def agent(config_path):
    config = json.loads(Path(config_path).read_text())
    evidence = Path(config['evidence'])
    argv = ['claude', '-p', '--input-format', 'stream-json', '--output-format', 'stream-json',
            '--verbose', '--dangerously-skip-permissions', '--model',
            os.environ.get('REVIEW_TEST_MODEL', 'claude-sonnet-5-5'), '--max-budget-usd', '5']
    with (evidence / 'claude-stream.jsonl').open('wb') as output, (evidence / 'claude-stderr.txt').open('wb') as errors:
        proc = subprocess.Popen(argv, cwd=config['project'], stdin=subprocess.PIPE, stdout=output, stderr=errors)
        write_json(evidence / 'agent-start.json', {'pid': proc.pid, 'argv': argv, 'started_at': time.time()})
        message = {'type': 'user', 'message': {'role': 'user', 'content': config['prompt']}}
        proc.stdin.write((json.dumps(message) + '\n').encode())
        proc.stdin.flush()
        # Keep input open so this is a genuine persistent stream-json conversation.
        try:
            code = proc.wait(timeout=TIMEOUT)
        except subprocess.TimeoutExpired:
            proc.terminate()
            try:
                code = proc.wait(timeout=10)
            except subprocess.TimeoutExpired:
                proc.kill()
                code = proc.wait()
        finally:
            proc.stdin.close()
        write_json(evidence / 'agent-exit.json', {'returncode': code, 'finished_at': time.time()})
    return code


def read_events(path):
    if not path.exists():
        return []
    events = []
    for line in path.read_text(errors='replace').splitlines():
        try:
            events.append(json.loads(line))
        except json.JSONDecodeError:
            pass  # The logger can be midway through writing its last event.
    return events


def envelopes(value):
    """Inspect serialized tool result bodies, never assistant claims."""
    if isinstance(value, dict):
        if set(value) == {'type', 'briefing', 'artifact', 'annotations', 'resolution'}:
            yield value
        for child in value.values():
            yield from envelopes(child)
    elif isinstance(value, list):
        for child in value:
            yield from envelopes(child)
    elif isinstance(value, str):
        try:
            child = json.loads(value)
        except (json.JSONDecodeError, ValueError):
            # Tool output can wrap a JSON document in explanatory output.
            decoder = json.JSONDecoder()
            for match in re.finditer(r'\{', value):
                try:
                    child, _ = decoder.raw_decode(value[match.start():])
                except json.JSONDecodeError:
                    continue
                yield from envelopes(child)
            return
        if not isinstance(child, str):
            yield from envelopes(child)


def tool_results(events):
    for event in events:
        message = event.get('message', {})
        for block in message.get('content', []) if isinstance(message.get('content'), list) else []:
            if isinstance(block, dict) and block.get('type') == 'tool_result':
                yield block


def assistant_text(events):
    parts = []
    for event in events:
        if event.get('type') == 'assistant':
            for block in event.get('message', {}).get('content', []):
                if isinstance(block, dict) and block.get('type') == 'text':
                    parts.append(block.get('text', ''))
        elif event.get('type') == 'result':
            parts.append(event.get('result', ''))
    return '\n'.join(parts)


def main():
    project = Path.cwd()
    evidence = Path(os.environ.get('REVIEW_TEST_ARTIFACTS', str(project / 'evidence'))).resolve()
    evidence.mkdir(parents=True, exist_ok=True)
    fixture = project / 'plan.md'
    original = fixture.read_bytes()
    assert original.decode().count(TARGET) == 1, 'Fixture must contain REVIEW_TEST_TARGET exactly once'
    digest = 'sha256:' + hashlib.sha256(original).hexdigest()
    session = 'launch-clean-review-' + secrets.token_hex(4)
    deadline = time.monotonic() + TIMEOUT
    summary = {'status': 'running', 'session': session, 'artifact': str(fixture), 'expected_rev': digest,
               'target': TARGET, 'model': os.environ.get('REVIEW_TEST_MODEL', 'claude-sonnet-5-5')}
    write_json(evidence / 'summary.json', summary)
    prompt = (
        'Use the installed /spacedock:review skill to review the existing plan.md. '
        'Use local Terminal review in a tmux right split, with binding Approve / Send back decisions enabled. '
        'Do not use remote review and do not rewrite the file. Follow the installed skill and use the real '
        'released Spacedock Review tool; do not emulate any part of the review. The human will operate the '
        'review pane. Wait for the actual tool result. After the human finishes, report the exact decision '
        'and every comment body and its anchored quotation back in this same conversation. '
        'Do not decide or finish the review for the human. Do not inspect environment variables, credentials, or authentication files.'
    )
    config = {'project': str(project), 'evidence': str(evidence), 'prompt': prompt}
    config_path = evidence / 'runner-config.json'
    write_json(config_path, config)
    tmux('new-session', '-d', '-s', session, '-x', '240', '-y', '60', '-c', str(project))
    original_pane = tmux('display-message', '-p', '-t', session + ':0.0', '#{pane_id}').stdout.strip()
    # Explicit shell quoting is required only by tmux's pane command interface.
    command = shlex.join([sys.executable, str(Path(__file__).resolve()), '--agent', str(config_path)])
    tmux('send-keys', '-t', original_pane, '-l', command)
    tmux('send-keys', '-t', original_pane, 'Enter')
    pane = None

    def wait_for(description, condition, seconds=90):
        end = min(deadline, time.monotonic() + seconds)
        while time.monotonic() < end:
            result = condition()
            if result:
                return result
            time.sleep(0.5)
        raise TimeoutError('Timed out waiting for ' + description)

    def capture(target):
        return tmux('capture-pane', '-p', '-t', target).stdout

    def save(name, target=None):
        text = capture(target or pane)
        (evidence / name).write_text(text)
        return text

    def pane_ids():
        return tmux('list-panes', '-t', session, '-F', '#{pane_id}').stdout.splitlines()

    try:
        wait_for('Claude process start', lambda: (evidence / 'agent-start.json').exists(), 30)
        # This value is never part of the model's prompt or any user input message.
        nonce = 'clean-review-' + secrets.token_hex(12)
        summary['nonce'] = nonce
        write_json(evidence / 'summary.json', summary)
        pane = wait_for('real Review split pane', lambda: next((p for p in pane_ids() if p != original_pane), None), 180)
        summary['review_pane'] = pane
        wait_for('artifact in initially uncluttered Review pane', lambda: TARGET in capture(pane), 90)
        save('01-review-before.txt')
        tmux('send-keys', '-t', pane, 'Tab')
        wait_for('Approve decision selected', lambda: 'Approve' in capture(pane), 15)
        text = save('02-approve-selected.txt')
        lines = text.splitlines()
        matches = [(i + 1, line.index(TARGET) + 1) for i, line in enumerate(lines) if TARGET in line]
        assert len(matches) == 1, f'Expected exactly one visible target, found {len(matches)}'
        y, x = matches[0]
        end_x = x + len(TARGET) - 1
        drag = f'\033[<0;{x};{y}M\033[<32;{end_x};{y}M\033[<0;{end_x};{y}m'
        tmux('send-keys', '-t', pane, '-l', drag)
        wait_for('selection comment shortcut', lambda: 'c comment' in capture(pane).lower(), 15)
        save('03-selected-anchor.txt')
        tmux('send-keys', '-t', pane, '-l', 'c')
        wait_for('comment editor', lambda: 'comment input' in capture(pane).lower(), 15)
        save('04-comment-editor.txt')
        tmux('send-keys', '-t', pane, '-l', nonce)
        tmux('send-keys', '-t', pane, 'Enter')
        wait_for('saved comment', lambda: 'draft saved' in capture(pane).lower(), 20)
        save('05-saved-comment.txt')
        tmux('send-keys', '-t', pane, '-l', 'q')
        wait_for('finish review summary', lambda: 'approve' in capture(pane).lower() and
                 'enter send' in capture(pane).lower(), 15)
        save('06-finish-summary.txt')
        tmux('send-keys', '-t', pane, 'Enter')
        stream = evidence / 'claude-stream.jsonl'

        def received_result():
            for result in tool_results(read_events(stream)):
                for envelope in envelopes(result.get('content')):
                    resolution = envelope.get('resolution')
                    if isinstance(resolution, dict) and resolution.get('decision') == 'approve':
                        return envelope
            return None

        envelope = wait_for('canonical approved tool result in original Claude stream', received_result, 120)
        write_json(evidence / 'canonical-result.json', envelope)
        assert envelope['artifact']['rev'] == digest, 'Artifact revision differs from original file SHA256'
        assert fixture.read_bytes() == original, 'Agent changed the original artifact'
        annotations = [a for a in envelope['annotations'] if a.get('body') == nonce]
        assert len(annotations) == 1, 'Expected exactly one controller nonce comment'
        annotation = annotations[0]
        assert any(s.get('exact') == TARGET for s in annotation.get('selectors', [])), 'Anchor exact quotation mismatch'
        assert annotation['id'] in envelope['resolution']['includes'], 'Comment not included in approved resolution'
        final_text = wait_for('Claude receipt of the controller-only comment',
                              lambda: (t if nonce in t else None) if (t := assistant_text(read_events(stream))) else None, 90)
        (evidence / 'claude-receipt.txt').write_text(final_text)
        # A stream-json result event establishes completion even while stdin stays open.
        wait_for('completed Claude turn', lambda: any(e.get('type') == 'result' for e in read_events(stream)), 30)
        save('07-original-agent-pane.txt', original_pane)
        summary.update(status='passed', decision='approve', annotation_id=annotation['id'],
                       comment_received_by_same_agent=True, completed_at=time.time())
    except Exception as error:
        summary.update(status='failed', error=str(error), traceback=traceback.format_exc(), completed_at=time.time())
        for p in pane_ids():
            try:
                save('failure-pane-' + p.replace('%', '') + '.txt', p)
            except Exception:
                pass
        raise
    finally:
        write_json(evidence / 'summary.json', summary)
        # Closing the test-owned tmux session also ends any still-open stream input.
        tmux('kill-session', '-t', session, check=False)
        # Redact only this runner's known evidence files, never unrelated project files.
        for path in evidence.iterdir():
            if path.is_file() and (path.name in {'runner-config.json', 'summary.json', 'agent-start.json', 'agent-exit.json', 'claude-stream.jsonl', 'claude-stderr.txt', 'canonical-result.json', 'claude-receipt.txt'} or re.fullmatch(r'(?:[0-9]{2}-[a-z-]+|failure-pane-[0-9]+)\.txt', path.name)):
                content = path.read_text(errors='replace')
                for key in ('ANTHROPIC_API_KEY', 'GH_TOKEN', 'GITHUB_TOKEN'):
                    secret = os.environ.get(key)
                    if secret:
                        content = content.replace(secret, '[REDACTED ' + key + ']')
                path.write_text(content)
    print(json.dumps(summary, indent=2))
    return 0


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--agent')
    args = parser.parse_args()
    sys.exit(agent(args.agent) if args.agent else main())
