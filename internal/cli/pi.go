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
	repoRoot        string
	packageRoot     string
	extensionPath   string
	subagentsSkill  string
	firstOfficer    string
	ensign          string
	authPath        string
	openAIAPIKey    string
	pluginDirSource string
}

type piCheckResult struct {
	piBinOK          bool
	piBin            string
	authOK           bool
	extensionOK      bool
	subagentsSkillOK bool
	firstOfficerOK   bool
	ensignOK         bool
	packageRoot      string
	repoRoot         string
	authPath         string
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
	if piRuntimeLaunchReady(check) {
		fmt.Fprintf(stdout, "Pi runtime ready.\n  pi-subagents: %s\n  Spacedock skills: %s\n", check.packageRoot, check.repoRoot)
		return 0
	}
	if checkOnly {
		printPiDoctorReport(stdout, check)
		return piDoctorExit(check)
	}
	fmt.Fprintf(stdout,
		"Pi runtime setup incomplete.\n\n"+
			"Required next steps:\n"+
			"  1. Install Pi and authenticate so %s exists.\n"+
			"  2. Install the subagent substrate, for example: pi install npm:pi-subagents\n"+
			"  3. If pi-subagents is installed outside the default location, set PI_SUBAGENTS_PACKAGE_ROOT.\n"+
			"  4. Re-run: spacedock doctor --host pi\n\n", check.authPath)
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
	if pkg == "" {
		home := envMap["HOME"]
		if home == "" {
			home = os.Getenv("HOME")
		}
		pkg = filepath.Join(home, ".pi", "agent", "npm", "node_modules", "pi-subagents")
	}
	authRoot := envMap["PI_CODING_AGENT_DIR"]
	authPath := ""
	if authRoot != "" {
		authPath = filepath.Join(authRoot, "auth.json")
	} else {
		home := envMap["HOME"]
		if home == "" {
			home = os.Getenv("HOME")
		}
		authPath = filepath.Join(home, ".pi", "agent", "auth.json")
	}
	return piRuntimeConfig{
		repoRoot:        repo,
		packageRoot:     pkg,
		extensionPath:   filepath.Join(pkg, "src", "extension", "index.ts"),
		subagentsSkill:  filepath.Join(pkg, "skills", "pi-subagents"),
		firstOfficer:    filepath.Join(repo, "skills", "first-officer", "SKILL.md"),
		ensign:          filepath.Join(repo, "skills", "ensign", "SKILL.md"),
		authPath:        authPath,
		openAIAPIKey:    envMap["OPENAI_API_KEY"],
		pluginDirSource: pluginDirSource,
	}
}

func checkPiRuntime(ops piRuntimeOps, cfg piRuntimeConfig) piCheckResult {
	bin, err := ops.LookPath("pi")
	res := piCheckResult{piBinOK: err == nil, piBin: bin, packageRoot: cfg.packageRoot, repoRoot: cfg.repoRoot, authPath: cfg.authPath}
	res.authOK = ops.Stat(cfg.authPath) == nil || strings.TrimSpace(cfg.openAIAPIKey) != ""
	res.extensionOK = ops.Stat(cfg.extensionPath) == nil
	res.subagentsSkillOK = ops.Stat(filepath.Join(cfg.subagentsSkill, "SKILL.md")) == nil
	res.firstOfficerOK = ops.Stat(cfg.firstOfficer) == nil
	res.ensignOK = ops.Stat(cfg.ensign) == nil
	return res
}

func piRuntimeLaunchReady(c piCheckResult) bool {
	return c.piBinOK && c.extensionOK && c.subagentsSkillOK && c.firstOfficerOK && c.ensignOK
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
	printPiCheck(w, c.firstOfficerOK, "Spacedock first-officer skill", filepath.Join(c.repoRoot, "skills", "first-officer"), "pass --plugin-dir <spacedock checkout> or set SPACEDOCK_REPO_ROOT")
	printPiCheck(w, c.ensignOK, "Spacedock ensign skill", filepath.Join(c.repoRoot, "skills", "ensign"), "pass --plugin-dir <spacedock checkout> or set SPACEDOCK_REPO_ROOT")
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
