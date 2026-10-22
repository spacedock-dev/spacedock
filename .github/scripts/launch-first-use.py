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
TIMEOUT = int(os.environ.get('REVIEW_TEST_TIMEOUT', '360'))


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
    base = ['claude', '-p', '--input-format', 'stream-json', '--output-format', 'stream-json',
            '--verbose', '--dangerously-skip-permissions', '--model',
            os.environ.get('REVIEW_TEST_MODEL', 'claude-sonnet-5-5')]
    commands = evidence / 'controller-input.jsonl'
    deadline = time.monotonic() + TIMEOUT
    seen = 0
    spent = 0.0
    proc = None
    with (evidence / 'claude-stream.jsonl').open('ab') as output, (evidence / 'claude-stderr.txt').open('ab') as errors:
        def launch(resume=None):
            argv = base + ['--max-budget-usd', str(max(0.01, 5.0 - spent))]
            if resume:
                argv += ['--resume', resume]
            process = subprocess.Popen(argv, cwd=config['project'], stdin=subprocess.PIPE, stdout=output, stderr=errors)
            write_json(evidence / 'agent-start.json', {'pid': process.pid, 'argv': argv, 'started_at': time.time()})
            return process
        def send(process, content):
            message = {'type': 'user', 'message': {'role': 'user', 'content': content}}
            process.stdin.write((json.dumps(message) + '\n').encode())
            process.stdin.flush()
        proc = launch()
        send(proc, config['prompt'])
        try:
            while time.monotonic() < deadline:
                entries = read_events(commands)
                for entry in entries[seen:]:
                    seen += 1
                    if entry.get('action') == 'restart':
                        results = [e for e in read_events(evidence / 'claude-stream.jsonl') if e.get('type') == 'result']
                        spent = max([float(e.get('total_cost_usd', 0)) for e in results] + [spent])
                        if spent >= 5:
                            raise RuntimeError('Total Claude budget reached before restart')
                        proc.stdin.close()
                        try:
                            proc.wait(timeout=5)
                        except subprocess.TimeoutExpired:
                            proc.terminate()
                            proc.wait(timeout=5)
                        proc = launch(entry['session_id'])
                    send(proc, entry['content'])
                if proc.poll() is not None:
                    # The controller may request an explicit resume after this process exits.
                    write_json(evidence / 'agent-last-exit.json', {'returncode': proc.returncode})
                time.sleep(0.2)
        finally:
            if proc and proc.poll() is None:
                proc.terminate()
                try:
                    proc.wait(timeout=5)
                except subprocess.TimeoutExpired:
                    proc.kill()
                    proc.wait()
            write_json(evidence / 'agent-exit.json', {'returncode': proc.returncode if proc else None})
    return 0


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
        if (set(value) == {'type', 'briefing', 'artifact', 'annotations', 'resolution'}
            or (set(value) == {'type', 'mode', 'briefing', 'artifact', 'annotations', 'actor'} and value.get('mode') == 'feedback')):
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


def inventory():
    import shutil
    result = {'commands': {}, 'known_paths': {}}
    for name in ('spacedock', 'sr', 'subspace-tui', 'compass'):
        path = shutil.which(name)
        item = {'path': path}
        if path:
            try:
                probe = run([path, '--version'], check=False, timeout=5)
                item.update(returncode=probe.returncode, stdout=probe.stdout, stderr=probe.stderr)
            except Exception as error:
                item['error'] = str(error)
        result['commands'][name] = item
    for directory in (Path.home() / '.local/bin', Path('/home/linuxbrew/.linuxbrew/bin'), Path('/usr/local/bin')):
        for name in ('spacedock', 'sr', 'subspace-tui', 'compass'):
            candidate = directory / name
            if candidate.exists():
                result['known_paths'][str(candidate)] = str(candidate.resolve())
    try:
        probe = run(['claude', 'plugin', 'list', '--json'], check=False, timeout=15)
        result['plugins'] = {'returncode': probe.returncode, 'stdout': probe.stdout, 'stderr': probe.stderr}
    except Exception as error:
        result['plugins'] = {'error': str(error)}
    return result


def human_choice_reply(text, structured_question=False):
    """Answer only an actual surface or terminal-family question."""
    low = text.lower()
    terminal = bool(re.search(r'\bterminal\b', low))
    alternative = bool(re.search(r'\b(?:browser|remote|web)\b', low))
    questions = re.findall(r'[^.!?\n]*\?', low)
    asks_surface = structured_question or bool(re.search(r'^\s*(?:please\s+)?choose\b', low, re.M))
    if terminal and alternative:
        asks_surface = asks_surface or any(
            (re.search(r'\bterminal\b', q) and re.search(r'\b(?:browser|remote|web)\b', q))
            or re.search(r'\bwhich (?:one|option|do you (?:want|prefer)|would you (?:like|prefer))\b', q)
            for q in questions)
        return 'Terminal' if asks_surface else None
    family_questions = [low] if structured_question else questions
    if any(re.search(r'\b(?:which|what)\b[^.!?]{0,100}\b(?:terminal|multiplexer)\b', q)
           or re.search(r'\b(?:terminal|multiplexer)\b[^.!?]{0,100}\b(?:are you using|do you use|are you running)\b', q)
           for q in family_questions):
        return 'tmux'
    return None


def main():
    project = Path.cwd()
    evidence = Path(os.environ.get('REVIEW_TEST_ARTIFACTS', str(project / 'evidence'))).resolve()
    evidence.mkdir(parents=True, exist_ok=True)
    owned = set()
    def record(name, value):
        owned.add(name)
        if isinstance(value, str):
            (evidence / name).write_text(value)
        else:
            write_json(evidence / name, value)
    fixture = project / 'plan.md'
    original = fixture.read_bytes()
    assert original.decode().count(TARGET) == 1, 'Fixture target must occur exactly once'
    digest = 'sha256:' + hashlib.sha256(original).hexdigest()
    session = 'launch-first-use-' + secrets.token_hex(4)
    deadline = time.monotonic() + TIMEOUT
    prompt = 'Install Spacedock using https://spacedock.md/agents.md and use Spacedock Review for plans and artifacts I need to review.'
    summary = {'status': 'running', 'outcome': None, 'prompt': prompt, 'target': TARGET,
               'expected_rev': digest, 'session': session, 'interventions': []}
    record('summary.json', summary)
    record('inventory-before.json', inventory())
    config_path = evidence / 'runner-config.json'
    record('runner-config.json', {'project': str(project), 'evidence': str(evidence), 'prompt': prompt})
    tmux('new-session', '-d', '-s', session, '-x', '240', '-y', '60', '-c', str(project))
    original_pane = tmux('display-message', '-p', '-t', session + ':0.0', '#{pane_id}').stdout.strip()
    tmux('send-keys', '-t', original_pane, '-l', shlex.join([sys.executable, str(Path(__file__).resolve()), '--agent', str(config_path)]))
    tmux('send-keys', '-t', original_pane, 'Enter')
    stream = evidence / 'claude-stream.jsonl'
    pane = None
    phase = 'setup'
    handled_results = 0
    terminal_questions = set()
    pending_terminal_result = None
    restart_used = False

    def capture(target):
        return tmux('capture-pane', '-p', '-t', target).stdout
    def panes():
        return tmux('list-panes', '-t', session, '-F', '#{pane_id}', check=False).stdout.splitlines()
    def save(name, target):
        text = capture(target)
        record(name, text)
        return text
    def send(content, restart=None, reason=None):
        entry = {'action': 'restart' if restart else 'send', 'content': content}
        if restart:
            entry['session_id'] = restart
        with (evidence / 'controller-input.jsonl').open('a') as output:
            output.write(json.dumps(entry) + '\n')
        summary['interventions'].append(dict(entry, reason=reason, at=time.time()))
        record('summary.json', summary)
    def wait_for(description, fn, seconds):
        end = min(deadline, time.monotonic() + seconds)
        while time.monotonic() < end:
            value = fn()
            if value:
                return value
            time.sleep(0.4)
        raise TimeoutError(description)
    def terminal_question(events):
        for event in events:
            if event.get('type') != 'assistant':
                continue
            for block in event.get('message', {}).get('content', []):
                if isinstance(block, dict) and block.get('type') == 'tool_use' and block.get('name') == 'AskUserQuestion':
                    text = json.dumps(block.get('input', {}))
                    reply = human_choice_reply(text, structured_question=True)
                    if reply:
                        key = block.get('id', text)
                        if key not in terminal_questions:
                            return key, text, reply
        return None
    def finish(outcome, **fields):
        summary.update(status='observed', outcome=outcome, **fields)

    try:
        wait_for('Claude process start', lambda: (evidence / 'agent-start.json').exists(), 30)
        nonce = 'first-use-' + secrets.token_hex(12)
        while time.monotonic() < deadline:
            events = read_events(stream)
            for candidate in panes():
                if candidate != original_pane and TARGET in capture(candidate):
                    pane = candidate
                    break
            if pane:
                break
            question = terminal_question(events)
            if question and question[0] not in terminal_questions:
                terminal_questions.add(question[0])
                record('human-choice-question.json', {'id': question[0], 'question': question[1], 'answer': question[2]})
                send(question[2], reason='Explicit review surface or terminal-family question from agent')
                pending_terminal_result = handled_results + 1
            results = [e for e in events if e.get('type') == 'result']
            if len(results) > handled_results:
                result = results[-1]
                handled_results = len(results)
                if pending_terminal_result is not None:
                    skip_question_result = handled_results <= pending_terminal_result
                    pending_terminal_result = None
                    if skip_question_result:
                        continue
                text = result.get('result', '') or assistant_text(events)
                record(f'{phase}-response.txt', text)
                low = text.lower()
                reply = human_choice_reply(text)
                if reply:
                    key = 'turn-' + str(handled_results)
                    if key not in terminal_questions:
                        terminal_questions.add(key)
                        send(reply, reason='Explicit review surface or terminal-family question in completed turn')
                        continue
                if result.get('is_error'):
                    finish('agent_error', agent_result=result)
                    break
                state = inventory()
                record(f'inventory-after-{phase}.json', state)
                installed = state['commands']['spacedock']['path'] or any(p.endswith('/spacedock') for p in state['known_paths'])
                if phase == 'setup':
                    if not installed:
                        finish('setup_not_completed', checkpoint='Exact public CTA finished without an installed Spacedock command', agent_response=text)
                        break
                    needs_restart = bool(re.search(r'\brestart\b|\bnew (?:claude |agent )?session\b', low))
                    if needs_restart:
                        session_ids = [e.get('session_id') for e in events if e.get('session_id')]
                        if not session_ids:
                            finish('restart_requested_no_session_id', agent_response=text)
                            break
                        restart_used = True
                        send('Review plan.md.', restart=session_ids[-1], reason='Agent requested restart to load installed plugins')
                    else:
                        send('Review plan.md.', reason='Natural first review after literal installation CTA')
                    phase = 'review'
                    continue
                if 'restart' in low and not restart_used:
                    session_ids = [e.get('session_id') for e in events if e.get('session_id')]
                    if session_ids:
                        restart_used = True
                        send('Review plan.md.', restart=session_ids[-1], reason='Review step explicitly requested restart')
                        continue
                # A background launch can return its conversational turn before the pane renders.
                grace = min(deadline, time.monotonic() + 15)
                while time.monotonic() < grace:
                    pane = next((p for p in panes() if p != original_pane and TARGET in capture(p)), None)
                    if pane:
                        break
                    time.sleep(0.4)
                if not pane:
                    finish('review_not_opened', agent_response=text)
                break
            time.sleep(0.5)
        if pane:
            summary['review_pane'] = pane
            save('01-review-before.txt', pane)
            tmux('send-keys', '-t', pane, 'Tab')
            time.sleep(0.7)
            text = save('02-review-after-tab.txt', pane)
            binding = 'Approve' in text
            summary['observed_mode'] = 'binding' if binding else 'feedback_only_no_approve_control'
            rows = [(i+1, line.index(TARGET)+1) for i, line in enumerate(text.splitlines()) if TARGET in line]
            if len(rows) != 1:
                raise RuntimeError('Could not locate one exact visible target after Tab')
            y, x = rows[0]
            end = x + len(TARGET) - 1
            drag = f'\033[<0;{x};{y}M\033[<32;{end};{y}M\033[<0;{end};{y}m'
            tmux('send-keys', '-t', pane, '-l', drag)
            wait_for('Selection comment control', lambda: 'c comment' in capture(pane).lower(), 15)
            save('03-selected-anchor.txt', pane)
            tmux('send-keys', '-t', pane, '-l', 'c')
            wait_for('Comment input', lambda: any(label in capture(pane).lower() for label in ('comment input', 'comment body')), 15)
            tmux('send-keys', '-t', pane, '-l', nonce)
            tmux('send-keys', '-t', pane, 'Enter')
            wait_for('Saved comment', lambda: 'draft saved' in capture(pane).lower(), 15)
            save('04-comment-saved.txt', pane)
            tmux('send-keys', '-t', pane, '-l', 'q')
            if binding:
                wait_for('Finish summary', lambda: 'enter send' in capture(pane).lower(), 15)
                save('05-finish-summary.txt', pane)
                tmux('send-keys', '-t', pane, 'Enter')
            # Released single-file feedback mode submits directly on q.
            def received():
                for index, event in enumerate(read_events(stream)):
                    for result in tool_results([event]):
                        for value in envelopes(result.get('content')):
                            if any(a.get('body') == nonce for a in value.get('annotations', [])):
                                return value, index
                return None
            value, receipt_event_index = wait_for('Canonical comment result in original agent conversation', received, 90)
            record('canonical-result.json', value)
            assert value['artifact']['rev'] == digest, 'Artifact SHA256 changed'
            annotations = [a for a in value['annotations'] if a.get('body') == nonce]
            assert len(annotations) == 1, 'Expected exactly one nonce comment'
            assert any(s.get('exact') == TARGET for s in annotations[0].get('selectors', [])), 'Exact anchor mismatch'
            resolution = value.get('resolution')
            if binding:
                assert resolution and resolution.get('decision') == 'approve', 'Binding UI did not return approve'
                assert annotations[0]['id'] in resolution.get('includes', []), 'Approval omits submitted comment'
            events = wait_for('Agent turn completed after canonical review receipt', lambda: (es if any(e.get('type') == 'result' for e in es[receipt_event_index + 1:]) else None) if (es := read_events(stream)) else None, 45)
            receipt = assistant_text(events)
            record('claude-receipt.txt', receipt)
            if binding:
                assert nonce in receipt, 'Original agent did not report received comment'
                assert fixture.read_bytes() == original, 'Original artifact was changed'
            finish('binding_review_completed' if binding else 'feedback_review_completed',
                   nonce=nonce, decision=resolution.get('decision') if isinstance(resolution, dict) else None,
                   exact_comment_in_tool_result=True, comment_echoed_by_agent=nonce in receipt,
                   artifact_unchanged=fixture.read_bytes() == original,
                   product_gap=None if binding else 'Natural first review did not expose an Approve decision')
        elif summary['outcome'] is None:
            finish('observation_timeout', checkpoint=phase)
    except Exception as error:
        summary.update(status='harness_error', outcome='controller_could_not_complete', error=str(error), traceback=traceback.format_exc())
    finally:
        for candidate in panes():
            try:
                save('final-pane-' + candidate.replace('%', '') + '.txt', candidate)
            except Exception:
                pass
        tmux('kill-session', '-t', session, check=False)
        summary['restart_used'] = restart_used
        summary['completed_at'] = time.time()
        record('inventory-final.json', inventory())
        record('summary.json', summary)
        owned.update({'agent-start.json', 'agent-exit.json', 'agent-last-exit.json', 'claude-stream.jsonl', 'claude-stderr.txt', 'controller-input.jsonl'})
        for name in owned:
            path = evidence / name
            if not path.is_file():
                continue
            content = path.read_text(errors='replace')
            for key in ('ANTHROPIC_API_KEY', 'GH_TOKEN', 'GITHUB_TOKEN'):
                secret = os.environ.get(key)
                if secret:
                    content = content.replace(secret, '[REDACTED ' + key + ']')
            path.write_text(content)
    print((evidence / 'summary.json').read_text())
    return 0 if summary.get('outcome') == 'binding_review_completed' else 1


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--agent')
    args = parser.parse_args()
    sys.exit(agent(args.agent) if args.agent else main())
