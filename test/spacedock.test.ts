// Unit tests for the Spacedock pi extension module (.pi/extensions/spacedock.ts):
// the taskless-launch sendUserMessage delivery path (the single bootstrap
// delivery arm), the PI_SPACEDOCK_LAUNCH gate for the compaction boot-record
// arm, and the resources_discover manifest-skip (the intra-package
// double-registration remedy). Run: bun test test/spacedock.test.ts
import { describe, expect, test } from "bun:test";
import registerSpacedockExtension from "../.pi/extensions/spacedock.ts";

type Handler = (event?: any) => any;

function makeFakePi() {
	const handlers: Record<string, Handler> = {};
	const sent: any[] = [];
	const fake = {
		handlers,
		sent,
		on(name: string, handler: Handler) {
			handlers[name] = handler;
		},
		exec: async () => ({ code: 0, stdout: "" }),
		sendUserMessage(content: any) {
			sent.push(content);
		},
	};
	return fake;
}

// A fake whose sendUserMessage throws synchronously — the no-crash-on-throw
// harness for the single delivery path.
function makeThrowingFakePi() {
	const fake = makeFakePi();
	fake.sendUserMessage = () => {
		throw new Error("send unavailable");
	};
	return fake;
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

describe("spacedock extension FO delivery gate", () => {
	test("taskless frontdoor launch delivers the bootstrap as a real first-turn user message via sendUserMessage", async () => {
		const fake = makeFakePi();
		registerSpacedockExtension(fake as any);
		await asParentSession(() =>
			withEnv("PI_SPACEDOCK_LAUNCH", "1", () =>
				withEnv("PI_SPACEDOCK_LAUNCH_TASKLESS", "1", () => {
					fake.handlers["session_start"]?.({ type: "session_start", reason: "startup" }, { ui: {} });
					expect(fake.sent).toHaveLength(1);
					const text = typeof fake.sent[0] === "string" ? fake.sent[0] : injectedText({ messages: [{ role: "user", content: fake.sent[0] }] });
					// The extension-owned bootstrap text, resolvable trigger: the
					// <available_skills> listing, not unexpandable $spacedock: syntax
					// and not a relative SKILL.md path.
					expect(text).toContain("SPACEDOCK-FO-BOOTSTRAP-v1");
					expect(text).toContain("<available_skills>");
					expect(text).not.toContain("$spacedock:");
					expect(text).not.toContain("skills/first-officer/SKILL.md");
				}),
			),
		);
	});

	test("taskful / resume / plain launches and non-startup session_start reasons send ZERO messages", async () => {
		const cases: Array<{ name: string; env: Record<string, string | undefined>; reason?: string }> = [
			{ name: "taskful frontdoor launch (no taskless marker)", env: { PI_SPACEDOCK_LAUNCH: "1", PI_SPACEDOCK_LAUNCH_TASKLESS: "0" } },
			{ name: "resume frontdoor launch (taskless marker 0)", env: { PI_SPACEDOCK_LAUNCH: "1", PI_SPACEDOCK_LAUNCH_TASKLESS: "0" }, reason: "resume" },
			{ name: "plain pi session (no markers)", env: { PI_SPACEDOCK_LAUNCH: undefined, PI_SPACEDOCK_LAUNCH_TASKLESS: undefined } },
			{ name: "reload of a taskless session (reason gate)", env: { PI_SPACEDOCK_LAUNCH: "1", PI_SPACEDOCK_LAUNCH_TASKLESS: "1" }, reason: "reload" },
		];
		for (const tc of cases) {
			const fake = makeFakePi();
			registerSpacedockExtension(fake as any);
			await asParentSession(() =>
				withEnv("PI_SPACEDOCK_LAUNCH", tc.env.PI_SPACEDOCK_LAUNCH, () =>
					withEnv("PI_SPACEDOCK_LAUNCH_TASKLESS", tc.env.PI_SPACEDOCK_LAUNCH_TASKLESS, () => {
						fake.handlers["session_start"]?.({ type: "session_start", reason: tc.reason ?? "startup" }, { ui: {} });
						expect(fake.sent).toHaveLength(0);
					}),
				),
			);
		}
	});

	test("a throwing send is caught: no crash, named diagnostic, no bootstrap", async () => {
		const fake = makeThrowingFakePi();
		registerSpacedockExtension(fake as any);
		await asParentSession(() =>
			withEnv("PI_SPACEDOCK_LAUNCH", "1", () =>
				withEnv("PI_SPACEDOCK_LAUNCH_TASKLESS", "1", () => {
					// Must NOT throw — a throwing send leaves the launch running
					// without a bootstrap.
					fake.handlers["session_start"]?.({ type: "session_start", reason: "startup" }, { ui: {} });
					expect(fake.sent).toHaveLength(0);
				}),
			),
		);
	});

	test("NO context-hook bootstrap arm remains: the context hook injects nothing on a taskless launch's first request", async () => {
		const fake = makeFakePi();
		registerSpacedockExtension(fake as any);
		await asParentSession(() =>
			withEnv("PI_SPACEDOCK_LAUNCH", "1", () =>
				withEnv("PI_SPACEDOCK_LAUNCH_TASKLESS", "1", async () => {
					fake.handlers["session_start"]?.({ type: "session_start", reason: "startup" }, { ui: {} });
					expect(fake.sent).toHaveLength(1);
					// The bootstrap went out via sendUserMessage; the context hook
					// (no compaction) must inject nothing — request-time-only
					// bootstrap injection is gone.
					const result = await runContext(fake.handlers["context"], [userMessage("operator task")]);
					expect(result).toBeUndefined();
				}),
			),
		);
	});

	test("compaction boot-record injection carries the same gate", async () => {
		const fake = makeFakePi();
		registerSpacedockExtension(fake as any);
		await asParentSession(() =>
			withEnv("PI_SPACEDOCK_LAUNCH", undefined, async () => {
				fake.handlers["session_start"]?.({ type: "session_start", reason: "startup" }, { ui: {} });
				fake.handlers["session_compact"]?.();
				const result = await runContext(fake.handlers["context"], [userMessage("compaction summary")]);
				expect(result).toBeUndefined();
				expect(fake.sent).toHaveLength(0);
			}),
		);
		await asParentSession(() =>
			withEnv("PI_SPACEDOCK_LAUNCH", "1", async () => {
				fake.handlers["session_start"]?.({ type: "session_start", reason: "startup" }, { ui: {} });
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
