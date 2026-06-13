package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"

	"github.com/spacedock-dev/spacedock/internal/safehouse"
)

const piBootstrapPrompt = "Use $spacedock:first-officer for this whole Pi session."

type piRuntimeOps interface {
	LookPath(name string) (string, error)
	Stat(path string) error
	Launch(argv []string) error
}

type execPiRuntimeOps struct{}

func (execPiRuntimeOps) LookPath(name string) (string, error) { return exec.LookPath(name) }
func (execPiRuntimeOps) Stat(path string) error               { _, err := os.Stat(path); return err }
func (execPiRuntimeOps) Launch(argv []string) error           { return execHost{}.Launch(argv, os.Environ()) }

type piRuntimeConfig struct {
	repoRoot              string
	packageRoot           string
	intercomPackageRoot   string
	extensionPath         string
	subagentsSkill        string
	firstOfficer          string
	ensign                string
	authPath              string
	openAIAPIKey          string
	sessionDir            string
	pluginDirSource       string
	packageRootSource     string
	intercomPackageSource string
	authPathSource        string
	sessionDirSource      string
}

type piCheckResult struct {
	piBinOK                   bool
	piBin                     string
	authOK                    bool
	extensionOK               bool
	subagentsSkillOK          bool
	subagentsIntercomBridgeOK bool
	intercomPackageOK         bool
	intercomSkillOK           bool
	firstOfficerOK            bool
	ensignOK                  bool
	packageRoot               string
	intercomPackageRoot       string
	repoRoot                  string
	authPath                  string
	sessionDir                string
}

func runPi(ctx context.Context, args []string, dir string, env []string, ops piRuntimeOps, stdout, stderr io.Writer) int {
	_ = ctx
	fd, pluginDirs, err := parsePiFrontDoorArgs(args)
	if err != nil {
		fmt.Fprintf(stderr, "spacedock pi: %v\n", err)
		return 2
	}
	cfg := piRuntimeConfigFromEnv(env, dir, lastString(pluginDirs))
	check := checkPiRuntime(ops, cfg)
	if !piRuntimeLaunchReady(check) {
		fmt.Fprint(stderr, "spacedock pi: Pi runtime is not ready; run `spacedock doctor --host pi` or `spacedock install --host pi`\n")
		printPiDoctorReport(stdout, check)
		return 1
	}

	launchBanner("pi", dir, safehouse.Present(dir), ops.LookPath, stderr)

	argv := []string{
		"pi",
		"--extension", cfg.extensionPath,
		"--skill", cfg.subagentsSkill,
		"--skill", cfg.firstOfficerDir(),
		"--skill", cfg.ensignDir(),
	}
	argv = append(argv, fd.passthrough...)
	argv = append(argv, launchPrompt(piBootstrapPrompt, fd))
	if err := ops.Launch(argv); err != nil {
		fmt.Fprintf(stderr, "spacedock pi: launch failed: %v\n", err)
		return 1
	}
	return 0
}

func runInitWithPi(ctx context.Context, args []string, hostOps hostOps, piOps piRuntimeOps, env []string, stdout, stderr io.Writer) int {
	host, checkOnly, pluginDir, code := parsePiSetupArgs("install", args, stderr)
	if code != 0 {
		return code
	}
	if host != "pi" {
		return runInit(ctx, args, hostOps, stdout, stderr)
	}
	cfg := piRuntimeConfigFromEnv(env, cwd(), pluginDir)
	check := checkPiRuntime(piOps, cfg)
	if checkOnly {
		printPiDoctorReport(stdout, check)
		return piDoctorExit(check)
	}
	if piRuntimeLaunchReady(check) {
		fmt.Fprintf(stdout, "Pi runtime ready.\n  pi-subagents: %s\n  pi-intercom: %s\n  Spacedock skills: %s\n", check.packageRoot, check.intercomPackageRoot, check.repoRoot)
		printPiSupervisorTalkbackBoundary(stdout)
		return 0
	}
	fmt.Fprintf(stdout,
		"Pi runtime setup incomplete.\n\n"+
			"Required next steps:\n"+
			"  1. Install Pi and authenticate so %s exists.\n"+
			"  2. Install the subagent substrate, for example: pi install npm:pi-subagents\n"+
			"  3. Install the supervisor-talkback substrate, for example: pi install npm:pi-intercom or npm install pi-intercom into the Pi npm root.\n"+
			"  4. If pi-subagents or pi-intercom are installed outside the default locations, set PI_SUBAGENTS_PACKAGE_ROOT and PI_INTERCOM_PACKAGE_ROOT.\n"+
			"  5. Re-run: spacedock doctor --host pi\n\n", check.authPath)
	printPiDoctorReport(stdout, check)
	return 0
}

func runDoctorWithPi(ctx context.Context, args []string, hostOps hostOps, piOps piRuntimeOps, env []string, stdout, stderr io.Writer) int {
	host, _, pluginDir, code := parsePiSetupArgs("doctor", args, stderr)
	if code != 0 {
		return code
	}
	if host != "pi" {
		return runDoctor(ctx, args, hostOps, stdout, stderr)
	}
	cfg := piRuntimeConfigFromEnv(env, cwd(), pluginDir)
	check := checkPiRuntime(piOps, cfg)
	printPiDoctorReport(stdout, check)
	return piDoctorExit(check)
}

func parsePiFrontDoorArgs(args []string) (fd frontDoorArgs, pluginDirs []string, err error) {
	fs := pflag.NewFlagSet("spacedock-pi", pflag.ContinueOnError)
	fs.SetOutput(io.Discard)
	pluginDir := fs.StringArray("plugin-dir", nil, "Load local Spacedock skill checkout")
	if err := fs.Parse(args); err != nil {
		return frontDoorArgs{}, nil, err
	}
	positionals := fs.Args()
	dash := fs.ArgsLenAtDash()
	var taskTokens []string
	if dash < 0 {
		taskTokens = positionals
	} else {
		taskTokens = positionals[:dash]
		fd.passthrough = positionals[dash:]
	}
	if len(taskTokens) > 0 {
		fd.task = strings.Join(taskTokens, " ")
		fd.hasTask = true
	}
	return fd, *pluginDir, nil
}

func parsePiSetupArgs(command string, args []string, stderr io.Writer) (host string, check bool, pluginDir string, code int) {
	host = "claude"
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--host":
			if i+1 >= len(args) {
				fmt.Fprintf(stderr, "spacedock %s: --host requires a value (claude, codex, or pi)\n", command)
				return "", false, "", 2
			}
			host = args[i+1]
			i++
		case "--check":
			if command != "install" {
				fmt.Fprintf(stderr, "spacedock %s: unknown argument %q\n", command, args[i])
				return "", false, "", 2
			}
			check = true
		case "--plugin-manifest":
			if command != "doctor" {
				fmt.Fprintf(stderr, "spacedock %s: unknown argument %q\n", command, args[i])
				return "", false, "", 2
			}
			if i+1 >= len(args) {
				fmt.Fprintln(stderr, "spacedock doctor: --plugin-manifest requires a path")
				return "", false, "", 2
			}
			i++
		case "--plugin-dir":
			if command == "install" {
				fmt.Fprintln(stderr, "spacedock install: --plugin-dir is not supported; use SPACEDOCK_REPO_ROOT or run from the Spacedock checkout")
				return "", false, "", 2
			}
			if i+1 >= len(args) {
				fmt.Fprintf(stderr, "spacedock %s: --plugin-dir requires a path\n", command)
				return "", false, "", 2
			}
			pluginDir = args[i+1]
			i++
		default:
			fmt.Fprintf(stderr, "spacedock %s: unknown argument %q\n", command, args[i])
			return "", false, "", 2
		}
	}
	return host, check, pluginDir, 0
}

func piRuntimeConfigFromEnv(env []string, dir, pluginDir string) piRuntimeConfig {
	envMap := envMap(env)
	home := envMap["HOME"]
	if home == "" {
		home = os.Getenv("HOME")
	}
	repo := pluginDir
	pluginDirSource := "--plugin-dir"
	if repo == "" {
		repo = envMap["SPACEDOCK_REPO_ROOT"]
		pluginDirSource = "SPACEDOCK_REPO_ROOT"
	}
	if repo == "" {
		repo = dir
		pluginDirSource = "working directory"
	}
	pkg := envMap["PI_SUBAGENTS_PACKAGE_ROOT"]
	pkgSource := "PI_SUBAGENTS_PACKAGE_ROOT"
	if pkg == "" {
		pkg = filepath.Join(home, ".pi", "agent", "npm", "node_modules", "pi-subagents")
		pkgSource = "default ~/.pi/agent/npm/node_modules/pi-subagents"
	}
	intercomPkg := envMap["PI_INTERCOM_PACKAGE_ROOT"]
	intercomPkgSource := "PI_INTERCOM_PACKAGE_ROOT"
	if intercomPkg == "" {
		intercomPkg = filepath.Join(home, ".pi", "agent", "npm", "node_modules", "pi-intercom")
		intercomPkgSource = "default ~/.pi/agent/npm/node_modules/pi-intercom"
	}
	authRoot := envMap["PI_CODING_AGENT_DIR"]
	authPath := ""
	authPathSource := "PI_CODING_AGENT_DIR"
	if authRoot != "" {
		authPath = filepath.Join(authRoot, "auth.json")
	} else {
		authPath = filepath.Join(home, ".pi", "agent", "auth.json")
		authPathSource = "default ~/.pi/agent/auth.json"
	}
	sessionDir := envMap["PI_CODING_AGENT_SESSION_DIR"]
	sessionDirSource := "PI_CODING_AGENT_SESSION_DIR"
	if sessionDir == "" {
		sessionDir = filepath.Join(home, ".pi", "agent", "sessions")
		sessionDirSource = "default ~/.pi/agent/sessions"
	}
	return piRuntimeConfig{
		repoRoot:              repo,
		packageRoot:           pkg,
		intercomPackageRoot:   intercomPkg,
		extensionPath:         filepath.Join(pkg, "src", "extension", "index.ts"),
		subagentsSkill:        filepath.Join(pkg, "skills", "pi-subagents"),
		firstOfficer:          filepath.Join(repo, "skills", "first-officer", "SKILL.md"),
		ensign:                filepath.Join(repo, "skills", "ensign", "SKILL.md"),
		authPath:              authPath,
		openAIAPIKey:          envMap["OPENAI_API_KEY"],
		sessionDir:            sessionDir,
		pluginDirSource:       pluginDirSource,
		packageRootSource:     pkgSource,
		intercomPackageSource: intercomPkgSource,
		authPathSource:        authPathSource,
		sessionDirSource:      sessionDirSource,
	}
}

func checkPiRuntime(ops piRuntimeOps, cfg piRuntimeConfig) piCheckResult {
	bin, err := ops.LookPath("pi")
	res := piCheckResult{
		piBinOK:             err == nil,
		piBin:               bin,
		packageRoot:         cfg.packageRoot,
		intercomPackageRoot: cfg.intercomPackageRoot,
		repoRoot:            cfg.repoRoot,
		authPath:            cfg.authPath,
		sessionDir:          cfg.sessionDir,
	}
	res.authOK = ops.Stat(cfg.authPath) == nil || strings.TrimSpace(cfg.openAIAPIKey) != ""
	res.extensionOK = ops.Stat(cfg.extensionPath) == nil
	res.subagentsSkillOK = ops.Stat(filepath.Join(cfg.subagentsSkill, "SKILL.md")) == nil
	res.subagentsIntercomBridgeOK = ops.Stat(filepath.Join(cfg.packageRoot, "src", "intercom", "intercom-bridge.ts")) == nil
	res.intercomPackageOK = ops.Stat(cfg.intercomPackageRoot) == nil
	res.intercomSkillOK = ops.Stat(filepath.Join(cfg.intercomPackageRoot, "skills", "pi-intercom", "SKILL.md")) == nil
	res.firstOfficerOK = ops.Stat(cfg.firstOfficer) == nil
	res.ensignOK = ops.Stat(cfg.ensign) == nil
	return res
}

func piRuntimeLaunchReady(c piCheckResult) bool {
	return c.piBinOK && c.extensionOK && c.subagentsSkillOK && c.subagentsIntercomBridgeOK && c.intercomPackageOK && c.intercomSkillOK && c.firstOfficerOK && c.ensignOK
}

func piDoctorHealthy(c piCheckResult) bool {
	return piRuntimeLaunchReady(c) && c.authOK
}

func piDoctorExit(c piCheckResult) int {
	if piDoctorHealthy(c) {
		return 0
	}
	return 1
}

func printPiDoctorReport(w io.Writer, c piCheckResult) {
	fmt.Fprintln(w, "Pi runtime check")
	printPiCheck(w, c.piBinOK, "pi CLI", c.piBin, "install Pi and ensure `pi` is on PATH")
	printPiCheck(w, c.authOK, "Pi auth", c.authPath, "run `pi` login/auth flow; live tests copy this file into an isolated PI_CODING_AGENT_DIR")
	printPiCheck(w, c.extensionOK, "pi-subagents extension", filepath.Join(c.packageRoot, "src", "extension", "index.ts"), "run `pi install npm:pi-subagents` or set PI_SUBAGENTS_PACKAGE_ROOT")
	printPiCheck(w, c.subagentsSkillOK, "pi-subagents skill", filepath.Join(c.packageRoot, "skills", "pi-subagents"), "run `pi install npm:pi-subagents` or set PI_SUBAGENTS_PACKAGE_ROOT")
	fmt.Fprintf(w, "INFO Pi auth/session dirs: auth=%s session=%s\n", c.authPath, c.sessionDir)
	fmt.Fprintln(w, "Supervisor-talkback setup prerequisites")
	printPiCheck(w, c.subagentsIntercomBridgeOK, "pi-subagents intercom bridge", filepath.Join(c.packageRoot, "src", "intercom", "intercom-bridge.ts"), "install/update pi-subagents or set PI_SUBAGENTS_PACKAGE_ROOT to a package root containing the intercom bridge")
	printPiCheck(w, c.intercomPackageOK, "pi-intercom package root", c.intercomPackageRoot, "set PI_INTERCOM_PACKAGE_ROOT to the installed pi-intercom package root")
	printPiCheck(w, c.intercomSkillOK, "pi-intercom skill", filepath.Join(c.intercomPackageRoot, "skills", "pi-intercom"), "install pi-intercom or set PI_INTERCOM_PACKAGE_ROOT to a package root containing skills/pi-intercom/SKILL.md")
	printPiCheck(w, c.firstOfficerOK, "Spacedock first-officer skill", filepath.Join(c.repoRoot, "skills", "first-officer"), "pass --plugin-dir <spacedock checkout> or set SPACEDOCK_REPO_ROOT")
	printPiCheck(w, c.ensignOK, "Spacedock ensign skill", filepath.Join(c.repoRoot, "skills", "ensign"), "pass --plugin-dir <spacedock checkout> or set SPACEDOCK_REPO_ROOT")
	printPiSupervisorTalkbackBoundary(w)
}

func printPiSupervisorTalkbackBoundary(w io.Writer) {
	fmt.Fprintln(w, "NOTE: These checks verify necessary supervisor-talkback setup prerequisites only; they are insufficient to prove live child talkback.")
	fmt.Fprintln(w, "NOTE: Live proof still requires the cq-style progress -> decision -> supervisor reply -> child resume -> durable marker probe for pi-intercom-supervisor-talkback.")
}

func printPiCheck(w io.Writer, ok bool, label, path, remedy string) {
	status := "OK"
	if !ok {
		status = "MISSING"
	}
	if path != "" {
		fmt.Fprintf(w, "%s %s: %s\n", status, label, path)
	} else {
		fmt.Fprintf(w, "%s %s\n", status, label)
	}
	if !ok {
		fmt.Fprintf(w, "  remedy: %s\n", remedy)
	}
}

func (c piRuntimeConfig) firstOfficerDir() string { return filepath.Dir(c.firstOfficer) }
func (c piRuntimeConfig) ensignDir() string       { return filepath.Dir(c.ensign) }

func lastString(v []string) string {
	if len(v) == 0 {
		return ""
	}
	return v[len(v)-1]
}

func envMap(env []string) map[string]string {
	m := map[string]string{}
	for _, kv := range env {
		k, v, ok := strings.Cut(kv, "=")
		if ok {
			m[k] = v
		}
	}
	return m
}
