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
// It also commissions the parent session as Spacedock first officer through
// Pi's context hook. The launcher keeps only its minimal skill-trigger prompt;
// this extension owns the durable warm/contract bootstrap because context hooks
// can re-inject after compaction while launch argv prompts cannot.
//
// Compaction boundary: PR #738 (force-boot-at-compaction-boundary) established
// that at compaction the FO re-reads durable state (one «state.boot»()), NOT
// re-inject the contract.  The session_compact hook fires a boot read via
// pi.exec and injects the boot record as a context message; the contract
// survives compaction in the system prompt (rebuilt from the skill via
// resources_discover), so re-injecting FO_BOOTSTRAP_TEXT would be the thing
// #738 rejected.

import { fileURLToPath } from "node:url";
import * as path from "node:path";

const FO_BOOTSTRAP_MARKER = "SPACEDOCK-FO-BOOTSTRAP-v1";
const FO_BOOTSTRAP_TEXT = `<EXTREMELY_IMPORTANT>\n[${FO_BOOTSTRAP_MARKER}] You are the Spacedock first officer. Load the $spacedock:first-officer skill (skills/first-officer/SKILL.md) and treat it as your operating contract: re-satisfy every load precondition at its trigger (shared core, runtime adapter, write/merge/dispatch cores), re-read durable state before the next workflow effect — the compacted summary is not authoritative. Pi tool mapping: read/write/edit/bash/grep/find/ls; load skills via read; subagent via pi-subagents when available; plans live in plan files / TODO.md.\n</EXTREMELY_IMPORTANT>`;

const FO_BOOT_RECORD_MARKER = "[SPACEDOCK-FO-BOOT-v2]";
const FO_BOOT_RECORD_DIRECTIVE = `${FO_BOOT_RECORD_MARKER} Durable state boot record — re-read before the next workflow effect. The compacted summary is not authoritative. Resume the loop where it stopped; do NOT greet or re-present a session summary. Pi tool mapping: read/write/edit/bash/grep/find/ls; load skills via read; subagent via pi-subagents when available.`;

// pi-subagents marks every spawned child session with PI_SUBAGENT_CHILD=1
// (src/runs/shared/pi-args.ts); a child is a delegated worker by definition,
// so commissioning one as first officer would leak the FO contract into
// ensign boots. Skill discovery still applies; only FO injection is exempt.
const isPiSubagentChild = process.env.PI_SUBAGENT_CHILD === "1";

function messageText(message) {
	if (Array.isArray(message?.content)) {
		return message.content
			.filter((part) => part?.type === "text")
			.map((part) => String(part?.text ?? ""))
			.join("");
	}
	return String(message?.content ?? "");
}

function hasStructuralBootstrap(message) {
	if (message?.role !== "user") return false;
	const text = messageText(message);
	return text.startsWith("<EXTREMELY_IMPORTANT>") && text.includes(FO_BOOTSTRAP_MARKER);
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
	let injectBootstrap = false;
	let injectBootRecord = false;

	pi.on("resources_discover", () => {
		const extDir = path.dirname(fileURLToPath(import.meta.url));
		// .pi/extensions/ -> ../.. -> repo root -> skills/
		const repoRoot = path.resolve(extDir, "..", "..");
		const skillsDir = path.join(repoRoot, "skills");
		return { skillPaths: [skillsDir] };
	});

	pi.on("session_start", () => {
		injectBootstrap = true;
	});

	pi.on("session_compact", () => {
		injectBootstrap = false;
		injectBootRecord = true;
	});

	pi.on("agent_end", () => {
		injectBootstrap = false;
		injectBootRecord = false;
	});

	pi.on("context", async (event) => {
		if (isPiSubagentChild) return;

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

		if (injectBootstrap) {
			if (event.messages.some(hasStructuralBootstrap)) return;

			const bootstrapMessage = {
				role: "user",
				content: [{ type: "text", text: FO_BOOTSTRAP_TEXT }],
				timestamp: Date.now(),
			};

			let insertAt = 0;
			while (insertAt < event.messages.length && isLeadingCompactionSummary(event.messages[insertAt])) {
				insertAt++;
			}

			return {
				messages: [
					...event.messages.slice(0, insertAt),
					bootstrapMessage,
					...event.messages.slice(insertAt),
				],
			};
		}

		return;
	});
}
