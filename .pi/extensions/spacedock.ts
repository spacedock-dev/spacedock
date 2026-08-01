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

import { fileURLToPath } from "node:url";
import * as path from "node:path";

const FO_BOOTSTRAP_MARKER = "SPACEDOCK-FO-BOOTSTRAP-v1";
const FO_BOOTSTRAP_TEXT = `<EXTREMELY_IMPORTANT>\n[${FO_BOOTSTRAP_MARKER}] You are the Spacedock first officer. Load the $spacedock:first-officer skill (skills/first-officer/SKILL.md) and treat it as your operating contract: re-satisfy every load precondition at its trigger (shared core, runtime adapter, write/merge/dispatch cores), re-read durable state before the next workflow effect — the compacted summary is not authoritative. Pi tool mapping: read/write/edit/bash/grep/find/ls; load skills via read; subagent via pi-subagents when available; plans live in plan files / TODO.md.\n</EXTREMELY_IMPORTANT>`;

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

function isLeadingCompactionSummary(message) {
	const text = messageText(message).toLowerCase();
	return text.includes("compaction") && text.includes("summary");
}

export default function registerSpacedockExtension(pi) {
	let injectBootstrap = false;

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
		injectBootstrap = true;
	});

	pi.on("agent_end", () => {
		injectBootstrap = false;
	});

	pi.on("context", (event) => {
		if (isPiSubagentChild || !injectBootstrap) return;
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
	});
}
