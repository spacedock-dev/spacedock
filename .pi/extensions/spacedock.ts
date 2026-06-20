// Spacedock pi extension — parent-session skill discovery.
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
import * as path from "node:path";

export default function registerSpacedockExtension(pi: {
	on(event: "resources_discover", handler: (event: { type: "resources_discover"; cwd: string; reason: string }) => { skillPaths?: string[] } | void): void;
}): void {
	pi.on("resources_discover", () => {
		const extDir = path.dirname(fileURLToPath(import.meta.url));
		// .pi/extensions/ -> ../.. -> repo root -> skills/
		const repoRoot = path.resolve(extDir, "..", "..");
		const skillsDir = path.join(repoRoot, "skills");
		return { skillPaths: [skillsDir] };
	});
}
