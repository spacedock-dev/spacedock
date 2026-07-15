// Spacedock pi extension — parent-session skill discovery and Bridge egress.
//
// Once `spacedock install --host pi` (or the dev `pi install ./local/path`)
// registers the Spacedock package in ~/.pi/agent/settings.json `packages`, the
// main pi session loads this extension. It implements `resources_discover` so
// the parent session discovers the Spacedock skills (first-officer, ensign, ...)
// from the package's own `skills/` directory — resolved relative to this
// extension's location, exactly like the obra/superpowers reference.
//
// This replaces the launcher's retired `--skill <repo>/skills/{first-officer,ensign}`
// flags. Child pi-subagents sessions discover the same skills independently via
// the package-root scan reading `package.json` `pi.skills` — no cwd dependency.

import { fileURLToPath } from "node:url";
import { spawn } from "node:child_process";
import * as path from "node:path";

export default function registerSpacedockExtension(pi: {
	on(event: "resources_discover", handler: (event: { type: "resources_discover"; cwd: string; reason: string }) => { skillPaths?: string[] } | void): void;
	on(event: string, handler: (event: Record<string, unknown>, ctx?: Record<string, unknown>) => void | Promise<void>): void;
}): void {
	pi.on("resources_discover", () => {
		const extDir = path.dirname(fileURLToPath(import.meta.url));
		// .pi/extensions/ -> ../.. -> repo root -> skills/
		const repoRoot = path.resolve(extDir, "..", "..");
		const skillsDir = path.join(repoRoot, "skills");
		return { skillPaths: [skillsDir] };
	});

	const lifecycleEvents = [
		"session_start",
		"session_shutdown",
		"agent_start",
		"agent_end",
		"turn_start",
		"turn_end",
		"tool_execution_start",
		"tool_execution_end",
		"tool_call",
		"tool_result",
	];
	for (const eventName of lifecycleEvents) {
		pi.on(eventName, (event, ctx) => {
			void emitBridgeEgress(eventName, event, ctx).catch(() => {});
		});
	}
}

async function emitBridgeEgress(eventName: string, event: Record<string, unknown>, ctx?: Record<string, unknown>): Promise<void> {
	const cwd = stringValue(event.cwd) || contextCwd(ctx) || process.cwd();
	const sessionFile = contextSessionFile(ctx);
	const payload = {
		event: eventName,
		cwd,
		session_file: sessionFile,
		session_id: sessionIdFromSessionFile(sessionFile),
		agent_id: "",
		agent_type: "",
		detail: eventDetail(eventName, event),
	};

	await invokeSpacedockEgress(payload, cwd);
}

function invokeSpacedockEgress(payload: Record<string, unknown>, cwd: string): Promise<void> {
	const bin = process.env.SPACEDOCK_BIN || "spacedock";
	return new Promise((resolve) => {
		const child = spawn(bin, ["bridge", "egress", "emit", "--host", "pi"], {
			cwd,
			stdio: ["pipe", "ignore", "ignore"],
		});
		child.on("error", () => resolve());
		child.on("close", () => resolve());
		child.stdin.on("error", () => resolve());
		try {
			child.stdin.end(JSON.stringify(payload));
		} catch {
			resolve();
		}
	});
}

function eventDetail(eventName: string, event: Record<string, unknown>): Record<string, unknown> {
	const detail: Record<string, unknown> = { source: "pi" };
	const tool = stringValue(event.toolName);
	if (tool) detail.tool = tool;
	if (eventName === "session_start" || eventName === "session_shutdown") {
		const reason = stringValue(event.reason);
		if (reason) detail.reason = reason;
	}
	const toolCallId = stringValue(event.toolCallId);
	if (toolCallId) detail.tool_call_id = toolCallId;
	return detail;
}

function contextCwd(ctx?: Record<string, unknown>): string {
	const opts = getSystemPromptOptions(ctx);
	if (!opts) return "";
	return stringValue(opts.cwd);
}

function contextSessionFile(ctx?: Record<string, unknown>): string {
	const manager = objectValue(ctx?.sessionManager);
	const getSessionFile = manager?.getSessionFile;
	if (typeof getSessionFile !== "function") return "";
	try {
		return stringValue(getSessionFile.call(manager));
	} catch {
		return "";
	}
}

function getSystemPromptOptions(ctx?: Record<string, unknown>): Record<string, unknown> | undefined {
	const fn = ctx?.getSystemPromptOptions;
	if (typeof fn !== "function") return undefined;
	try {
		return objectValue(fn.call(ctx));
	} catch {
		return undefined;
	}
}

function sessionIdFromSessionFile(sessionFile: string): string {
	if (!sessionFile) return "";
	const base = path.basename(sessionFile);
	const ext = path.extname(base);
	return ext ? base.slice(0, -ext.length) : base;
}

function objectValue(value: unknown): Record<string, unknown> | undefined {
	return value && typeof value === "object" ? value as Record<string, unknown> : undefined;
}

function stringValue(value: unknown): string {
	return typeof value === "string" ? value : "";
}
