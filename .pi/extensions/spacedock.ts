// Spacedock pi extension — parent-session skill discovery and FO bootstrap.
//
// Sessions spawned by pi-subagents carry PI_SUBAGENT_CHILD=1; those children
// are delegated workers, never first officers, so the FO-bootstrap context
// injection is skipped for them.
//
// Once `spacedock install --host pi` (or the dev `pi install ./local/path`)
// registers the Spacedock package in ~/.pi/agent/settings.json `packages`, the
// main pi session loads this extension. It implements `resources_discover` so
// the parent session discovers the Spacedock skills (first-officer, ensign, ...)
// from the package's own `skills/` directory — resolved relative to this
// extension's location, exactly like the obra/superpowers reference.
//
// It also delivers the FO contract in the parent session as a REAL first-turn
// user message via pi.sendUserMessage() — but ONLY on a taskless frontdoor
// launch (PI_SPACEDOCK_LAUNCH_TASKLESS=1, which the launcher sets exactly when
// no operator task argv and no resume passthrough is present: the only launch
// shape where nothing else triggers a first model request). The send both
// triggers the first turn and persists the bootstrap in the session transcript.
// A taskful launch or a resume already has its first request / history, so the
// extension sends nothing there; a plain `pi` session for unrelated work
// receives no bootstrap (single-owner gate, PI_SPACEDOCK_LAUNCH-gated). The
// delivery is send-only: if the send throws, the launch continues WITHOUT a
// bootstrap (caught, named diagnostic) — there is no context-hook bootstrap
// fallback arm. The launcher keeps only its argv-only duties (pass the
// operator task, suppress the launch prompt on resume); this extension owns
// the durable contract bootstrap.
//
// Compaction boundary: PR #738 (force-boot-at-compaction-boundary) established
// that at compaction the FO re-reads durable state (one «state.boot»()), NOT
// re-inject the contract.  The session_compact hook fires a boot read via
// pi.exec and injects the boot record as a context message; the contract
// survives compaction in the system prompt (rebuilt from the skill via
// resources_discover), so re-injecting FO_BOOTSTRAP_TEXT would be the thing
// #738 rejected.

import { fileURLToPath } from "node:url";
import * as fs from "node:fs";
import * as path from "node:path";

const FO_BOOTSTRAP_MARKER = "SPACEDOCK-FO-BOOTSTRAP-v1";
// The bootstrap names the RESOLVABLE trigger: the first-officer entry of the
// session's <available_skills> listing, whose location pi resolves from the
// registered package root regardless of cwd. No `$spacedock:` syntax (pi
// expands /skill: on user input only — a launch-injected message is not user
// input) and no relative SKILL.md path (ENOENT from any workflow cwd — the
// 2026-09-10 stale-skill incident's failure shape). /skill:first-officer
// remains the human-invocable form in interactive input.
const FO_BOOTSTRAP_TEXT = `<EXTREMELY_IMPORTANT>\n[${FO_BOOTSTRAP_MARKER}] You are the Spacedock first officer. Load the first-officer skill's SKILL.md at the exact location listed for it in your available skills (<available_skills>) and treat it as your operating contract: re-satisfy every load precondition at its trigger (shared core, runtime adapter, write/merge/dispatch cores), re-read durable state before the next workflow effect — the compacted summary is not authoritative. Pi tool mapping: read/write/edit/bash/grep/find/ls; load skills via read; subagent via pi-subagents when available; plans live in plan files / TODO.md.\n</EXTREMELY_IMPORTANT>`;

const FO_BOOT_RECORD_MARKER = "[SPACEDOCK-FO-BOOT-v2]";
const FO_BOOT_RECORD_DIRECTIVE = `${FO_BOOT_RECORD_MARKER} Durable state boot record — re-read before the next workflow effect. The compacted summary is not authoritative. Resume the loop where it stopped; do NOT greet or re-present a session summary. Pi tool mapping: read/write/edit/bash/grep/find/ls; load skills via read; subagent via pi-subagents when available.`;

// pi-subagents marks every spawned child session with PI_SUBAGENT_CHILD=1
// (src/runs/shared/pi-args.ts); a child is a delegated worker by definition,
// so commissioning one as first officer would leak the FO contract into
// ensign boots. Skill discovery still applies; only FO injection is exempt.
// Both gates are read at handler time so a probe/test can flip them.
function isPiSubagentChild() {
	return process.env.PI_SUBAGENT_CHILD === "1";
}

// isFrontdoorLaunch: the PI_SPACEDOCK_LAUNCH=1 marker `spacedock pi` sets on
// EVERY launch shape (wrap and non-wrap, installed and dev override). Without
// it — a plain `pi` session for unrelated work, a user's own project session —
// the FO bootstrap is NOT injected. This is the single-owner gate: the
// extension, not the launch prompt, delivers the contract.
function isFrontdoorLaunch() {
	return process.env.PI_SPACEDOCK_LAUNCH === "1";
}

// isTasklessLaunch: the PI_SPACEDOCK_LAUNCH_TASKLESS=1 marker the frontdoor
// sets ONLY when the launch carries no operator task argv and no resume
// passthrough. Env-only — the frontdoor never appends launch text (AC-3), and
// at session_start the argv prompt is not yet in the session, so the extension
// cannot distinguish taskless from taskful by content alone.
function isTasklessLaunch() {
	return process.env.PI_SPACEDOCK_LAUNCH_TASKLESS === "1";
}

// manifestDeclaredSkillDirs resolves the package manifest's `pi.skills` entries
// against repoRoot. Each declared dir is a skill-registration route pi already
// scans — resources_discover returning the same dir again doubles every skill
// registration (the packages x routes duplicate condition).
function manifestDeclaredSkillDirs(repoRoot) {
	try {
		const manifest = JSON.parse(fs.readFileSync(path.join(repoRoot, "package.json"), "utf8"));
		return (manifest?.pi?.skills ?? []).map((entry) => path.resolve(repoRoot, entry));
	} catch {
		return [];
	}
}

function messageText(message) {
	if (Array.isArray(message?.content)) {
		return message.content
			.filter((part) => part?.type === "text")
			.map((part) => String(part?.text ?? ""))
			.join("");
	}
	return String(message?.content ?? "");
}

function hasBootRecord(message) {
	if (message?.role !== "user") return false;
	return messageText(message).includes(FO_BOOT_RECORD_MARKER);
}

function isLeadingCompactionSummary(message) {
	const text = messageText(message).toLowerCase();
	return text.includes("compaction") && text.includes("summary");
}

export default function registerSpacedockExtension(pi) {
	let injectBootRecord = false;

	pi.on("resources_discover", () => {
		const extDir = path.dirname(fileURLToPath(import.meta.url));
		// .pi/extensions/ -> ../.. -> repo root -> skills/
		const repoRoot = path.resolve(extDir, "..", "..");
		const skillsDir = path.join(repoRoot, "skills");
		if (manifestDeclaredSkillDirs(repoRoot).includes(skillsDir)) {
			// The package manifest's pi.skills scan already registers this
			// directory; re-registering it here doubles every skill
			// registration. Skip it — the manifest route is the one that
			// counts, and the doctor's route-aware duplicate count shrinks
			// accordingly.
			return { skillPaths: [] };
		}
		return { skillPaths: [skillsDir] };
	});

	// Taskless frontdoor launch → deliver the FO bootstrap as a REAL first-turn
	// user message. This is the SINGLE delivery path: the send triggers the
	// first model turn (which a taskless launch otherwise never has) and
	// persists the bootstrap as a normal transcript message. No context-hook
	// bootstrap fallback exists — if the send throws, the launch gets no
	// bootstrap and must not crash (caught, named diagnostic below).
	// Gated on reason "startup" only: a reload/new/resume/fork of an existing
	// session already has turns and must not be re-greeted.
	pi.on("session_start", (event, ctx) => {
		if (event?.reason !== "startup") return;
		if (isPiSubagentChild() || !isFrontdoorLaunch() || !isTasklessLaunch()) return;
		try {
			pi.sendUserMessage(FO_BOOTSTRAP_TEXT);
		} catch (err) {
			const detail = err instanceof Error ? err.message : String(err);
			const diagnostic = `spacedock: FO bootstrap send failed — launch continues without greet (${detail})`;
			try {
				ctx?.ui?.notify?.(diagnostic, "warning");
			} catch {
				console.error(diagnostic);
			}
		}
	});

	pi.on("session_compact", () => {
		injectBootRecord = true;
	});

	pi.on("agent_end", () => {
		injectBootRecord = false;
	});

	pi.on("context", async (event) => {
		if (isPiSubagentChild() || !isFrontdoorLaunch()) return;

		if (injectBootRecord) {
			if (event.messages.some(hasBootRecord)) return;

			let bootRecordJson = "";
			try {
				const result = await pi.exec("spacedock", ["status", "--boot", "--identify", "--json"]);
				if (result.code === 0 && result.stdout) {
					bootRecordJson = result.stdout;
				}
			} catch {
				// Boot read failed — fall back to the directive without the
				// boot record so a compaction boundary never blocks the
				// session.  The directive alone still tells the FO to re-read
				// durable state; the extension just couldn't pre-read it.
			}

			const text = bootRecordJson
				? `${FO_BOOT_RECORD_DIRECTIVE}\n\n${bootRecordJson}`
				: FO_BOOT_RECORD_DIRECTIVE;

			const bootMessage = {
				role: "user",
				content: [{ type: "text", text }],
				timestamp: Date.now(),
			};

			let insertAt = 0;
			while (insertAt < event.messages.length && isLeadingCompactionSummary(event.messages[insertAt])) {
				insertAt++;
			}

			return {
				messages: [
					...event.messages.slice(0, insertAt),
					bootMessage,
					...event.messages.slice(insertAt),
				],
			};
		}

		// (The once-per-turn FO bootstrap injection arm was removed: the
		// bootstrap now ships as a real first-turn user message on taskless
		// launches via sendUserMessage — see session_start. Only the compaction
		// boot record remains request-time-only, by design: it rides requests
		// and must never persist as a transcript message.)
		return;
	});
}
