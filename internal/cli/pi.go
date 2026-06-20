package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/pflag"

	"github.com/spacedock-dev/spacedock/internal/safehouse"
)

const piBootstrapPrompt = "Use $spacedock:first-officer for this whole Pi session."

// piSpacedockPackageSource is the published install source for the Spacedock
// package. `spacedock install --host pi` runs `pi install <source>`, which
// registers the package in ~/.pi/agent/settings.json `packages` and places the
// repo in pi's package store. The dev override (--plugin-dir) replaces this with
// a local checkout path so in-tree edits are picked up without reinstall.
const piSpacedockPackageSource = "git:github.com/spacedock-dev/spacedock"

type piRuntimeOps interface {
	LookPath(name string) (string, error)
	Stat(path string) error
	Launch(argv []string) error
	// PiInstall runs `pi install <source>` and returns its combined output. The
	// real implementation execs `pi`; tests record the source and return canned
	// output. This is the install seam that retires Pi's check-only status.
	PiInstall(source string) (string, error)
	// SpacedockPackageStatus reports whether the Spacedock package is registered
	// in ~/.pi/agent/settings.json `packages` and whether the ensign (and
	// first-officer) skills are discoverable via the package-root skill scan —
	// the same mechanism pi-subagents' collectSettingsPackageSkillPaths uses.
	// agentDir is the pi agent directory (~/.pi/agent or PI_CODING_AGENT_DIR);
	// home is used to resolve ~ entries.
	SpacedockPackageStatus(agentDir, home string) piPackageStatus
}

// piPackageStatus is the result of the package-registration + skill-discovery
// check that replaces the retired repo-path Stat skill checks.
type piPackageStatus struct {
	registered               bool
	ensignDiscoverable       bool
	firstOfficerDiscoverable bool
	source                   string // the settings.json packages entry for spacedock
	packageRoot              string // the resolved package root
}

type execPiRuntimeOps struct{}

func (execPiRuntimeOps) LookPath(name string) (string, error) { return exec.LookPath(name) }
func (execPiRuntimeOps) Stat(path string) error               { _, err := os.Stat(path); return err }
func (execPiRuntimeOps) Launch(argv []string) error           { return execHost{}.Launch(argv, os.Environ()) }

func (execPiRuntimeOps) PiInstall(source string) (string, error) {
	out, err := exec.Command("pi", "install", source).CombinedOutput()
	return string(out), err
}

func (execPiRuntimeOps) SpacedockPackageStatus(agentDir, home string) piPackageStatus {
	return piSpacedockPackageStatus(agentDir, home)
}

type piRuntimeConfig struct {
	repoRoot              string // dev-override only: --plugin-dir / SPACEDOCK_REPO_ROOT
	packageRoot           string
	intercomPackageRoot   string
	extensionPath         string
	subagentsSkill        string
	authPath              string
	openAIAPIKey          string
	sessionDir            string
	agentDir              string
	home                  string
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
	spacedockPackageOK        bool
	packageStatus             piPackageStatus
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

	// The Spacedock first-officer/ensign skills are no longer passed as --skill
	// flags: the installed package's .pi/extensions/spacedock.ts extension
	// discovers them for the parent session via resources_discover, and
	// pi-subagents children discover them via the package-root scan. Only the
	// pi-subagents extension + skill are passed explicitly here.
	argv := []string{
		"pi",
		"--extension", cfg.extensionPath,
		"--skill", cfg.subagentsSkill,
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
		// --plugin-dir is the pi dev-override install source; it is not supported
		// for claude/codex install (which use the host plugin marketplace).
		if pluginDir != "" {
			fmt.Fprintln(stderr, "spacedock install: --plugin-dir is not supported; use SPACEDOCK_REPO_ROOT or run from the Spacedock checkout")
			return 2
		}
		return runInit(ctx, args, hostOps, stdout, stderr)
	}
	if !checkOnly {
		source := piSpacedockPackageSource
		if pluginDir != "" {
			source = pluginDir
		}
		out, err := piOps.PiInstall(source)
		if strings.TrimSpace(out) != "" {
			fmt.Fprint(stdout, out)
		}
		if err != nil {
			fmt.Fprintf(stderr, "spacedock install: pi install %q failed: %v\n", source, err)
			return 1
		}
	}
	cfg := piRuntimeConfigFromEnv(env, cwd(), pluginDir)
	check := checkPiRuntime(piOps, cfg)
	if checkOnly {
		printPiDoctorReport(stdout, check)
		return piDoctorExit(check)
	}
	printPiDoctorReport(stdout, check)
	if piRuntimeLaunchReady(check) {
		fmt.Fprintf(stdout, "Pi runtime ready.\n  pi-subagents: %s\n  pi-intercom: %s\n  Spacedock package: %s\n", check.packageRoot, check.intercomPackageRoot, check.packageStatus.source)
		printPiSupervisorTalkbackBoundary(stdout)
		return 0
	}
	fmt.Fprintf(stdout, "Pi runtime setup incomplete.\n\n"+
		"Required next steps:\n"+
		"  1. Install Pi and authenticate so %s exists.\n"+
		"  2. Install the subagent substrate, for example: pi install npm:pi-subagents\n"+
		"  3. Install the supervisor-talkback substrate, for example: pi install npm:pi-intercom or npm install pi-intercom into the Pi npm root.\n"+
		"  4. Install the Spacedock package: spacedock install --host pi\n"+
		"  5. If pi-subagents or pi-intercom are installed outside the default locations, set PI_SUBAGENTS_PACKAGE_ROOT and PI_INTERCOM_PACKAGE_ROOT.\n"+
		"  6. Re-run: spacedock doctor --host pi\n\n", check.authPath)
	printPiSupervisorTalkbackBoundary(stdout)
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
			// Accepted for both install (pi dev-override source) and doctor.
			// Non-pi install rejects it in runInitWithPi; non-pi doctor rejects
			// it via the re-parse in runDoctor.
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
	_ = dir
	envMap := envMap(env)
	home := envMap["HOME"]
	if home == "" {
		home = os.Getenv("HOME")
	}
	// repoRoot is the dev-override path only (--plugin-dir / SPACEDOCK_REPO_ROOT).
	// The cwd fallback is removed (D5c): the installed package is discovered via
	// the package-root scan regardless of cwd, and repoRoot is no longer needed
	// for skill discovery (the retired --skill flags were its only consumer).
	repo := pluginDir
	pluginDirSource := "--plugin-dir"
	if repo == "" {
		repo = envMap["SPACEDOCK_REPO_ROOT"]
		pluginDirSource = "SPACEDOCK_REPO_ROOT"
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
	agentDir := envMap["PI_CODING_AGENT_DIR"]
	authPathSource := "PI_CODING_AGENT_DIR"
	if agentDir == "" {
		agentDir = filepath.Join(home, ".pi", "agent")
		authPathSource = "default ~/.pi/agent"
	}
	authPath := filepath.Join(agentDir, "auth.json")
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
		authPath:              authPath,
		openAIAPIKey:          envMap["OPENAI_API_KEY"],
		sessionDir:            sessionDir,
		agentDir:              agentDir,
		home:                  home,
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
	// The retired repo-path Stat checks (firstOfficerOK/ensignOK) are replaced by
	// spacedockPackageOK: the package is registered AND ensign is discoverable via
	// the package-root skill scan — the real discovery contract, not a filesystem
	// coincidence at a cwd-derived path.
	status := ops.SpacedockPackageStatus(cfg.agentDir, cfg.home)
	res.packageStatus = status
	res.spacedockPackageOK = status.registered && status.ensignDiscoverable
	return res
}

func piRuntimeLaunchReady(c piCheckResult) bool {
	return c.piBinOK && c.extensionOK && c.subagentsSkillOK && c.subagentsIntercomBridgeOK && c.intercomPackageOK && c.intercomSkillOK && c.spacedockPackageOK
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
	printPiCheck(w, c.spacedockPackageOK, "Spacedock package", piPackageReportPath(c.packageStatus), "run `spacedock install --host pi` to install the Spacedock package (or `spacedock install --host pi --plugin-dir <checkout>` for a dev override)")
	printPiSupervisorTalkbackBoundary(w)
}

func piPackageReportPath(s piPackageStatus) string {
	if s.packageRoot != "" {
		return s.packageRoot
	}
	if s.source != "" {
		return s.source
	}
	return ""
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

// piSpacedockPackageStatus replicates the package-registration + skill-discovery
// contract that pi-subagents' collectSettingsPackageSkillPaths uses: it reads
// ~/.pi/agent/settings.json `packages`, resolves each entry's package root (via
// resolveSettingsPackageRoot), reads the package's package.json, and confirms a
// package named "spacedock" is registered with ensign discoverable under its
// pi.skills paths. This is the real discovery check — not a Stat of a cwd-derived
// path — so it holds from a non-repo cwd once the package is installed.
func piSpacedockPackageStatus(agentDir, home string) piPackageStatus {
	if agentDir == "" {
		return piPackageStatus{}
	}
	data, err := os.ReadFile(filepath.Join(agentDir, "settings.json"))
	if err != nil {
		return piPackageStatus{}
	}
	var settings struct {
		Packages []json.RawMessage `json:"packages"`
	}
	if json.Unmarshal(data, &settings) != nil {
		return piPackageStatus{}
	}
	for _, raw := range settings.Packages {
		src := piPackageSourceFromEntry(raw)
		if src == "" {
			continue
		}
		root := resolveSettingsPackageRoot(src, agentDir, home)
		if root == "" {
			continue
		}
		name, skillPaths := readPackagePiSkills(root)
		if name != "spacedock" {
			continue
		}
		st := piPackageStatus{registered: true, source: src, packageRoot: root}
		for _, sp := range skillPaths {
			dir := filepath.Join(root, sp)
			if piSkillFileExists(dir, "ensign") {
				st.ensignDiscoverable = true
			}
			if piSkillFileExists(dir, "first-officer") {
				st.firstOfficerDiscoverable = true
			}
		}
		return st
	}
	return piPackageStatus{}
}

func piPackageSourceFromEntry(raw json.RawMessage) string {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	var obj struct {
		Source string `json:"source"`
	}
	if json.Unmarshal(raw, &obj) == nil {
		return obj.Source
	}
	return ""
}

func readPackagePiSkills(root string) (name string, skills []string) {
	data, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return "", nil
	}
	var pkg struct {
		Name string `json:"name"`
		Pi   struct {
			Skills []string `json:"skills"`
		} `json:"pi"`
	}
	if json.Unmarshal(data, &pkg) != nil {
		return "", nil
	}
	return pkg.Name, pkg.Pi.Skills
}

func piSkillFileExists(dir, skill string) bool {
	_, err := os.Stat(filepath.Join(dir, skill, "SKILL.md"))
	return err == nil
}

// resolveSettingsPackageRoot replicates pi-subagents' resolveSettingsPackageRoot:
// it resolves a settings.json `packages` entry to a filesystem package root,
// handling git:, npm:, file:, ~, absolute, and relative path sources.
func resolveSettingsPackageRoot(source, baseDir, home string) string {
	s := strings.TrimSpace(source)
	if s == "" {
		return ""
	}
	if strings.HasPrefix(s, "git:") {
		host, repoPath := parseGitPackagePath(strings.TrimSpace(strings.TrimPrefix(s, "git:")))
		if host == "" || repoPath == "" {
			return ""
		}
		return filepath.Join(baseDir, "git", host, repoPath)
	}
	if strings.HasPrefix(s, "npm:") {
		name := parseNpmPackageName(strings.TrimSpace(strings.TrimPrefix(s, "npm:")))
		if name == "" {
			return ""
		}
		return filepath.Join(baseDir, "npm", "node_modules", name)
	}
	norm := strings.TrimPrefix(s, "file:")
	if norm == "~" {
		return home
	}
	if strings.HasPrefix(norm, "~/") {
		return filepath.Join(home, norm[2:])
	}
	if filepath.IsAbs(norm) {
		return norm
	}
	if norm == "." || norm == ".." || strings.HasPrefix(norm, "./") || strings.HasPrefix(norm, "../") {
		return filepath.Join(baseDir, norm)
	}
	return ""
}

var scpGitRe = regexp.MustCompile(`^git@([^:]+):(.+)$`)

// parseGitPackagePath replicates pi-subagents' parseGitPackagePath, returning the
// host and normalized repo path for a git: package source.
func parseGitPackagePath(spec string) (host, repoPath string) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return "", ""
	}
	if m := scpGitRe.FindStringSubmatch(spec); m != nil {
		host = m[1]
		repoPath = m[2]
	} else if u, err := url.Parse(spec); err == nil && u.IsAbs() && u.Host != "" {
		host = u.Hostname()
		repoPath = strings.TrimPrefix(u.Path, "/")
	} else if i := strings.Index(spec, "/"); i > 0 {
		host = spec[:i]
		repoPath = spec[i+1:]
	} else {
		return "", ""
	}
	repoPath = stripGitRef(repoPath)
	repoPath = strings.TrimSuffix(repoPath, ".git")
	repoPath = strings.TrimPrefix(repoPath, "/")
	if !isSafePackagePath(host) || !isSafePackagePath(repoPath) || len(strings.Split(repoPath, "/")) < 2 {
		return "", ""
	}
	return host, repoPath
}

// parseNpmPackageName replicates pi-subagents' parseNpmPackageName: it extracts
// the package name (without @version) from an npm: source.
func parseNpmPackageName(spec string) string {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return ""
	}
	re := regexp.MustCompile(`^(@?[^@]+(?:/[^@]+)?)(?:@(.+))?$`)
	m := re.FindStringSubmatch(spec)
	name := spec
	if m != nil && m[1] != "" {
		name = m[1]
	}
	if !isSafePackagePath(name) {
		return ""
	}
	return name
}

// stripGitRef replicates pi-subagents' stripGitRef: it strips a @ref or #ref
// suffix (the first @ or #, whichever comes first) from a repo path.
func stripGitRef(repoPath string) string {
	at := strings.Index(repoPath, "@")
	hash := strings.Index(repoPath, "#")
	var idx int = -1
	if at >= 0 && (hash < 0 || at < hash) {
		idx = at
	} else if hash >= 0 {
		idx = hash
	}
	if idx < 0 {
		return repoPath
	}
	return repoPath[:idx]
}

// isSafePackagePath replicates pi-subagents' isSafePackagePath: a path is safe
// when it is non-empty, not absolute, and has no "." or ".." segments.
func isSafePackagePath(value string) bool {
	if value == "" || filepath.IsAbs(value) {
		return false
	}
	for _, part := range strings.Split(value, string(filepath.Separator)) {
		if part == "" || part == "." || part == ".." {
			return false
		}
	}
	return true
}
