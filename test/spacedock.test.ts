// Unit tests for the Spacedock pi extension module (.pi/extensions/spacedock.ts):
// the PI_SPACEDOCK_LAUNCH injection gate (AC-1/AC-4) and the resources_discover
// manifest-skip (the intra-package double-registration remedy). Run: bun test
// ./.pi/extensions/spacedock.test.ts
import { describe, expect, test } from "bun:test";
import registerSpacedockExtension from "./spacedock.ts";

type Handler = (event?: any) => any;

function makeFakePi() {
	const handlers: Record<string, Handler> = {};
	return {
		handlers,
		on(name: string, handler: Handler) {
			handlers[name] = handler;
		},
		exec: async () => ({ code: 0, stdout: "" }),
	};
}

function userMessage(text: string) {
	return { role: "user", content: [{ type: "text", text }] };
}

function injectedText(result: any): string {
	const messages = result?.messages ?? [];
	return messages
		.map((message: any) =>
			(message?.content ?? [])
				.filter((part: any) => part?.type === "text")
				.map((part: any) => part?.text ?? "")
				.join(""),
		)
		.join("\n");
}

async function runContext(handler: Handler | undefined, messages: any[]) {
	if (!handler) throw new Error("no context handler registered");
	return await handler({ messages });
}

function withEnv(key: string, value: string | undefined, run: () => void | Promise<void>) {
	const prior = process.env[key];
	const restore = () => {
		if (prior === undefined) delete process.env[key];
		else process.env[key] = prior;
	};
	if (value === undefined) delete process.env[key];
	else process.env[key] = value;
	return Promise.resolve()
		.then(run)
		.finally(restore);
}

// The test runner itself may run inside a pi-subagents child (PI_SUBAGENT_CHILD=1);
// the injection gate reads it at handler time, so pin it out of the way.
function asParentSession(run: () => void | Promise<void>) {
	return withEnv("PI_SUBAGENT_CHILD", undefined, run);
}

describe("spacedock extension FO injection gate", () => {
	test("plain session (no PI_SPACEDOCK_LAUNCH) receives NO bootstrap", async () => {
		const fake = makeFakePi();
		registerSpacedockExtension(fake as any);
		await asParentSession(() =>
			withEnv("PI_SPACEDOCK_LAUNCH", undefined, async () => {
				fake.handlers["session_start"]?.();
				const result = await runContext(fake.handlers["context"], [userMessage("operator task")]);
				expect(result).toBeUndefined();
			}),
		);
	});

	test("frontdoor launch (PI_SPACEDOCK_LAUNCH=1) injects the bootstrap with the resolvable wording", async () => {
		const fake = makeFakePi();
		registerSpacedockExtension(fake as any);
		await asParentSession(() =>
			withEnv("PI_SPACEDOCK_LAUNCH", "1", async () => {
				fake.handlers["session_start"]?.();
				const result = await runContext(fake.handlers["context"], [userMessage("operator task")]);
				const text = injectedText(result);
				expect(text).toContain("SPACEDOCK-FO-BOOTSTRAP-v1");
				// Resolvable trigger: the <available_skills> listing, not unexpandable
				// $spacedock: syntax and not a relative SKILL.md path.
				expect(text).toContain("<available_skills>");
				expect(text).not.toContain("$spacedock:");
				expect(text).not.toContain("skills/first-officer/SKILL.md");
			}),
		);
	});

	test("compaction boot-record injection carries the same gate", async () => {
		const fake = makeFakePi();
		registerSpacedockExtension(fake as any);
		await asParentSession(() =>
			withEnv("PI_SPACEDOCK_LAUNCH", undefined, async () => {
				fake.handlers["session_start"]?.();
				fake.handlers["session_compact"]?.();
				const result = await runContext(fake.handlers["context"], [userMessage("compaction summary")]);
				expect(result).toBeUndefined();
			}),
		);
		await asParentSession(() =>
			withEnv("PI_SPACEDOCK_LAUNCH", "1", async () => {
				fake.handlers["session_start"]?.();
				fake.handlers["session_compact"]?.();
				const result = await runContext(fake.handlers["context"], [userMessage("compaction summary")]);
				expect(injectedText(result)).toContain("SPACEDOCK-FO-BOOT-v2");
			}),
		);
	});
});

describe("spacedock extension resources_discover manifest-skip", () => {
	test("skips the skills dir the package manifest already declares", () => {
		const fake = makeFakePi();
		registerSpacedockExtension(fake as any);
		// This test lives in <repoRoot>/.pi/extensions/, so the module's
		// repoRoot resolves to the Spacedock checkout whose package.json
		// declares pi.skills: ["./skills"] — the manifest route. The
		// extension must NOT re-register that directory.
		const result = fake.handlers["resources_discover"]?.();
		expect(result).toEqual({ skillPaths: [] });
	});
});
