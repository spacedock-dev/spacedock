package cli

import (
	"fmt"
	"strconv"
)

func validatePassthrough(fd frontDoorArgs, pluginDir bool) error {
	if fd.passthroughMode {
		if fd.hasTask {
			return fmt.Errorf("--passthrough rejects task text before --. Place native prompts after --")
		}
		if pluginDir {
			return fmt.Errorf("--passthrough rejects --plugin-dir before --. Place native plugin arguments after --")
		}
	}
	return nil
}

// When a corresponding explicit flag is absent, applySandboxEnv uses its environment value.
func applySandboxEnv(fd *frontDoorArgs, env []string) error {
	getenv := envGetenv(env)
	if value := getenv("SPACEDOCK_SAFEHOUSE"); !fd.safehouseExplicit && value != "" {
		force, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid SPACEDOCK_SAFEHOUSE value %q. Use a Boolean value or an explicit --safehouse flag", value)
		}
		fd.forceSafehouse = force
	}
	if value := getenv("SPACEDOCK_SAFEHOUSE_ENABLE"); !fd.enableExplicit && value != "" {
		fd.safehouseFlags = append([]string{"enable=" + value}, fd.safehouseFlags...)
	}
	return nil
}

// canonicalHostEnvPassFlags preserves only the selected host's explicit provider settings.
func canonicalHostEnvPassFlags(host string, provider bool, env []string) []string {
	if !provider {
		return nil
	}
	var names, flags []string
	switch host {
	case "claude":
		names = []string{"CLAUDE_CONFIG_DIR", "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS"}
	case "codex":
		names = []string{"CODEX_HOME"}
	case "pi":
		names = []string{"PI_CODING_AGENT_DIR"}
	}
	for _, name := range names {
		if hasEnv(env, name) {
			flags = append(flags, "--env-pass", name)
		}
	}
	return flags
}
