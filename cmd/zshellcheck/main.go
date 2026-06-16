// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sort"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/fix"
	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
	"github.com/afadesigns/zshellcheck/pkg/reporter"
	"github.com/afadesigns/zshellcheck/pkg/version"
)

func main() {
	os.Exit(run())
}

type runFlags struct {
	format         *string
	cpuprofile     *string
	showVersion    *bool
	verbose        *bool
	noColor        *bool
	noBanner       *bool
	severityFilter *string
	fixMode        *bool
	diffMode       *bool
	dryRun         *bool
	unsafeFixes    *bool
	listRules      *bool
	explain        *string
	statistics     *bool
	baseline       *string
	baselineWrite  *string
	ruleSeverity   *string
	addNoka        *bool
	detectStale    *bool
}

func run() int {
	flags := registerRunFlags()
	flag.Usage = func() {
		printUsage(os.Stderr, flag.CommandLine, !*flags.noBanner)
	}
	flag.Parse()

	if *flags.showVersion {
		fmt.Printf("zshellcheck version %s\n", version.Version)
		return 0
	}
	if *flags.listRules {
		return printRulesList(os.Stdout, katas.Registry)
	}
	if *flags.explain != "" {
		return printRuleExplain(os.Stdout, os.Stderr, katas.Registry, *flags.explain)
	}
	stopProfile, code := startCPUProfile(*flags.cpuprofile)
	if code != 0 {
		return code
	}
	defer stopProfile()
	if len(flag.Args()) < 1 {
		printRunUsage(*flags.noBanner)
		return 1
	}
	maybeEmitBanner(*flags.format, *flags.noColor, *flags.noBanner)

	cfg, err := resolveConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
		return 1
	}
	cfg = applyFlagOverrides(cfg, *flags.noColor, *flags.verbose)

	allowedSeverities, code := parseSeverityFilter(*flags.severityFilter)
	if code != 0 {
		return code
	}
	fixOpts := buildFixOpts(*flags.fixMode, *flags.diffMode, *flags.dryRun, *flags.unsafeFixes)
	if code := configureModes(flags, &fixOpts); code != 0 {
		return code
	}
	total := scanArgs(cfg, allowedSeverities, *flags.format, fixOpts)
	emitFixSummary(fixOpts.stats)
	return runResult(flags, total, fixOpts)
}

// configureModes parses the report-shaping flags onto fixOpts, returning a
// non-zero exit code only on a malformed value.
func configureModes(flags runFlags, fixOpts *fixOptions) int {
	if *flags.statistics {
		fixOpts.statistics = map[string]int{}
	}
	regrade, err := parseRuleSeverity(*flags.ruleSeverity)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 1
	}
	fixOpts.ruleSeverity = regrade
	fixOpts.addNoka = *flags.addNoka
	if *flags.detectStale {
		fixOpts.detectStale = true
		fixOpts.staleCount = new(int)
	}
	return setupBaseline(fixOpts, *flags.baseline, *flags.baselineWrite)
}

// runResult turns the scan total and run mode into the process exit code,
// emitting any deferred output (baseline file, statistics table).
func runResult(flags runFlags, total int, fixOpts fixOptions) int {
	switch {
	case *flags.addNoka:
		return 0
	case *flags.baselineWrite != "":
		if err := fixOpts.baseline.writeBaseline(*flags.baselineWrite); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing baseline: %s\n", err)
			return 1
		}
		return 0
	case fixOpts.statistics != nil:
		emitStatistics(os.Stdout, katas.Registry, fixOpts.statistics)
		if total == 0 && !hasStale(fixOpts) {
			return 0
		}
		return 1
	}
	exit := finalExitCode(total, *flags.format, fixOpts)
	if exit == 0 && hasStale(fixOpts) {
		return 1
	}
	return exit
}

// hasStale reports whether stale-suppression detection found anything.
func hasStale(fixOpts fixOptions) bool {
	return fixOpts.staleCount != nil && *fixOpts.staleCount > 0
}

// emitStatistics prints a per-kata count table sorted by descending count
// then kata ID, each row tagged `[*]` when the kata ships an auto-fix.
func emitStatistics(out io.Writer, registry *katas.KatasRegistry, counts map[string]int) {
	ids := make([]string, 0, len(counts))
	for id := range counts {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		if counts[ids[i]] != counts[ids[j]] {
			return counts[ids[i]] > counts[ids[j]]
		}
		return ids[i] < ids[j]
	})
	for _, id := range ids {
		mark := "   "
		if registry.IsFixable(id) {
			mark = "[*]"
		}
		title := ""
		if k, ok := registry.GetKata(id); ok {
			title = k.Title
		}
		fmt.Fprintf(out, "%5d  %-8s  %s  %s\n", counts[id], id, mark, title)
	}
}

// setupBaseline configures the run for -baseline / -baseline-write,
// returning a non-zero exit code only when a requested baseline file
// cannot be read.
func setupBaseline(fixOpts *fixOptions, baseline, baselineWrite string) int {
	switch {
	case baselineWrite != "":
		fixOpts.baseline = &baselineState{write: true}
	case baseline != "":
		b, err := loadBaseline(baseline)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading baseline: %s\n", err)
			return 1
		}
		fixOpts.baseline = b
	}
	return 0
}

func registerRunFlags() runFlags {
	return runFlags{
		format:         flag.String("format", "text", "Output format. One of text, json, sarif."),
		cpuprofile:     flag.String("cpuprofile", "", "Write a Go pprof CPU profile to this path."),
		showVersion:    flag.Bool("version", false, "Print the version and exit."),
		verbose:        flag.Bool("verbose", false, "Include the full kata description under each violation."),
		noColor:        flag.Bool("no-color", false, "Disable ANSI colours in the report."),
		noBanner:       flag.Bool("no-banner", false, "Suppress the startup banner — useful for CI and scripted runs."),
		severityFilter: flag.String("severity", "", "Comma-separated minimum severities to surface (error, warning, info, style)."),
		fixMode:        flag.Bool("fix", false, "Apply auto-fixes in place for katas that ship a deterministic rewrite."),
		diffMode:       flag.Bool("diff", false, "Print a unified diff of the fixes instead of writing them."),
		dryRun:         flag.Bool("dry-run", false, "With -fix, report what would change without modifying files."),
		unsafeFixes:    flag.Bool("unsafe-fixes", false, "Also apply fixes that may change runtime behavior (off by default)."),
		listRules:      flag.Bool("list-rules", false, "Print every kata (ID, severity, title) and exit."),
		explain:        flag.String("explain", "", "Print the full description of a kata by ID (e.g. ZC1001) and exit."),
		statistics:     flag.Bool("statistics", false, "Print a per-kata count of findings instead of individual reports."),
		baseline:       flag.String("baseline", "", "Suppress findings recorded in this baseline file; report only new ones."),
		baselineWrite:  flag.String("baseline-write", "", "Write a baseline snapshot of current findings to this path and exit 0."),
		ruleSeverity:   flag.String("rule-severity", "", "Re-grade katas: comma-separated ZC####:level (error, warning, info, style)."),
		addNoka:        flag.Bool("add-noka", false, "Append a `# noka: ZC####` directive to every line with a finding, then exit."),
		detectStale:    flag.Bool("detect-stale-noka", false, "Report `# noka` directives that suppress no actual finding."),
	}
}

func startCPUProfile(path string) (func(), int) {
	noop := func() {}
	if path == "" {
		return noop, 0
	}
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create CPU profile: %s\n", err)
		return noop, 1
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Fprintf(os.Stderr, "Could not start CPU profile: %s\n", err)
		_ = f.Close()
		return noop, 1
	}
	return func() {
		pprof.StopCPUProfile()
		_ = f.Close()
	}, 0
}

func printRunUsage(noBanner bool) {
	if !noBanner {
		fmt.Fprint(os.Stderr, config.Banner)
	}
	fmt.Println("Usage: zshellcheck [flags] <file1.zsh> [file2.zsh]...")
	fmt.Println("Try 'zshellcheck --help' for more information.")
}

func maybeEmitBanner(format string, noColor, noBanner bool) {
	if format == "json" || format == "sarif" || noColor || noBanner {
		return
	}
	fmt.Fprint(os.Stderr, config.Banner)
}

func resolveConfig() (config.Config, error) {
	xdgPath := xdgConfigSearch("zshellcheck/config.yml")
	if xdgPath == "" {
		xdgPath = xdgConfigSearch("zshellcheck/config.yaml")
	}
	homePath := filepath.Join(homeDir(), ".zshellcheckrc")
	return loadConfig(xdgPath, homePath, ".zshellcheckrc")
}

// xdgConfigSearch returns the first existing path for rel under the XDG
// configuration search path: `$XDG_CONFIG_HOME` (or `~/.config`), then the
// colon-separated `$XDG_CONFIG_DIRS` (or `/etc/xdg`). It returns "" when
// none exist. This replaces github.com/adrg/xdg so the binary carries no
// third-party dependencies.
func xdgConfigSearch(rel string) string {
	var dirs []string
	if home := os.Getenv("XDG_CONFIG_HOME"); home != "" {
		dirs = append(dirs, home)
	} else if home := homeDir(); home != "" {
		dirs = append(dirs, filepath.Join(home, ".config"))
	}
	sys := os.Getenv("XDG_CONFIG_DIRS")
	if sys == "" {
		sys = "/etc/xdg"
	}
	dirs = append(dirs, filepath.SplitList(sys)...)
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		path := filepath.Join(dir, rel)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// homeDir returns the user's home directory, or "" when it cannot be
// resolved.
func homeDir() string {
	if home, err := os.UserHomeDir(); err == nil {
		return home
	}
	return ""
}

func applyFlagOverrides(cfg config.Config, noColor, verbose bool) config.Config {
	if noColor {
		cfg.NoColor = true
	}
	if verbose {
		cfg.Verbose = true
	}
	return cfg
}

func parseSeverityFilter(filter string) ([]katas.Severity, int) {
	if filter == "" {
		return nil, 0
	}
	var out []katas.Severity
	for _, s := range strings.Split(filter, ",") {
		sev := katas.Severity(strings.TrimSpace(s))
		switch sev {
		case katas.SeverityError, katas.SeverityWarning, katas.SeverityInfo, katas.SeverityStyle:
			out = append(out, sev)
		default:
			fmt.Fprintf(os.Stderr, "Invalid severity level: %s. Must be one of error, warning, info, style.\n", s)
			return nil, 1
		}
	}
	return out, 0
}

func buildFixOpts(fixMode, diffMode, dryRun, unsafeFixes bool) fixOptions {
	opts := fixOptions{
		enabled:   fixMode || diffMode,
		diff:      diffMode,
		dryRun:    dryRun,
		unsafe:    unsafeFixes,
		maxPasses: 5,
	}
	if opts.diff {
		opts.enabled = true
		opts.dryRun = true
	}
	if opts.enabled {
		opts.stats = &fixStats{}
	}
	opts.fixable = new(int)
	opts.unsafeFixable = new(int)
	return opts
}

func scanArgs(cfg config.Config, allowed []katas.Severity, format string, fixOpts fixOptions) int {
	// The machine-readable formats accumulate findings across every file
	// and are emitted once as a single document; without this each file
	// printed its own array, so multi-file output was not valid JSON or
	// SARIF. Text streams per file.
	var collector *[]reporter.FileViolations
	if format == "json" || format == "sarif" {
		collector = &[]reporter.FileViolations{}
		fixOpts.collector = collector
	}
	total := 0
	for _, filename := range flag.Args() {
		total += processPath(filename, os.Stdout, os.Stderr, cfg, katas.Registry, format, allowed, fixOpts)
	}
	if collector != nil {
		emitAggregate(os.Stdout, os.Stderr, format, *collector)
	}
	return total
}

// emitAggregate writes the collected findings for the machine-readable
// formats as a single document.
func emitAggregate(out, errOut io.Writer, format string, files []reporter.FileViolations) {
	var err error
	switch format {
	case "json":
		err = reporter.ReportJSON(out, files)
	case "sarif":
		err = reporter.ReportSARIF(out, files, version.Version, sarifRuleMeta)
	}
	if err != nil {
		fmt.Fprintf(errOut, "Error reporting violations: %s\n", err)
	}
}

// sarifRuleMeta supplies SARIF rule metadata for a kata ID from the
// registry: its title, full description, and a link to the kata catalog.
func sarifRuleMeta(id string) reporter.RuleMeta {
	k, ok := katas.Registry.GetKata(id)
	if !ok {
		return reporter.RuleMeta{}
	}
	return reporter.RuleMeta{
		Name:        k.Title,
		Title:       k.Title,
		Description: k.Description,
		HelpURI:     "https://github.com/afadesigns/zshellcheck/blob/main/KATAS.md",
	}
}

func emitFixSummary(stats *fixStats) {
	if stats == nil || stats.filesScanned <= 1 {
		return
	}
	fmt.Fprintf(os.Stderr, "\nfix summary: %d edit(s) across %d file(s) (scanned %d)\n",
		stats.totalEdits, stats.filesModified, stats.filesScanned)
}

func finalExitCode(total int, format string, fixOpts fixOptions) int {
	if total == 0 {
		return 0
	}
	if format == "text" {
		fmt.Fprintf(os.Stderr, "\nFound %d violations.\n", total)
		// Point at the auto-fixable subsets, unless fixes are already
		// being applied or previewed this run.
		if !fixOpts.enabled {
			if fixOpts.fixable != nil && *fixOpts.fixable > 0 {
				fmt.Fprintf(os.Stderr, "[*] %d fixable with the `-fix` option.\n", *fixOpts.fixable)
			}
			if fixOpts.unsafeFixable != nil && *fixOpts.unsafeFixable > 0 {
				fmt.Fprintf(os.Stderr, "    %d more fixable with `-fix -unsafe-fixes`.\n", *fixOpts.unsafeFixable)
			}
		}
	}
	return 1
}

func loadConfig(paths ...string) (config.Config, error) {
	cfg := config.DefaultConfig()

	for _, path := range paths {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return cfg, err
		}

		fileConfig, err := config.Parse(data)
		if err != nil {
			return cfg, err
		}

		cfg = config.MergeConfig(cfg, fileConfig)
	}

	return cfg, nil
}

type fixOptions struct {
	enabled   bool
	diff      bool
	dryRun    bool
	maxPasses int
	// stats tracks per-run aggregate fix activity so processPath
	// can print a one-line summary footer when -fix runs over a
	// directory tree. nil when -fix is disabled.
	stats *fixStats
	// collector accumulates per-file findings for the machine-readable
	// formats so they are emitted once as a single JSON / SARIF document.
	// nil for the text format, which streams per file.
	collector *[]reporter.FileViolations
	// unsafe applies fixes that may change runtime behavior. When false,
	// only value-preserving (safe) fixes are applied.
	unsafe bool
	// fixable counts findings with a safe auto-fix; unsafeFixable counts
	// those whose only fix may change behavior. Both feed the post-report
	// `[*] N fixable` / `M with -unsafe-fixes` hints.
	fixable       *int
	unsafeFixable *int
	// statistics, when non-nil, switches output to a per-kata count table:
	// individual findings are suppressed and tallied here instead.
	statistics map[string]int
	// baseline, when non-nil, records or filters findings against a saved
	// snapshot so a run reports only findings new since the baseline.
	baseline *baselineState
	// ruleSeverity re-grades specific katas to a chosen severity.
	ruleSeverity map[string]katas.Severity
	// addNoka rewrites each file in place, appending a `# noka` directive
	// to every line with a finding, then stops short of reporting.
	addNoka bool
	// detectStale reports `# noka` directives that suppress no finding;
	// staleCount tallies them for the exit code.
	detectStale bool
	staleCount  *int
}

// fixStats accumulates fix activity across all files visited in one
// run. Updated by processFile when an in-place rewrite lands.
type fixStats struct {
	filesScanned  int
	filesModified int
	totalEdits    int
}

// applyFixesUntilStable runs fix.Apply repeatedly, re-parsing and
// re-collecting edits between passes, until no new edits are produced
// or maxPasses is reached. Returns the final source, the total
// number of edits applied across all passes, and any apply error.
//
// Multi-pass is needed because some fixes expose other fixes:
// `result=`which git“ first becomes `result=$(which git)` (ZC1002),
// which a second pass then rewrites to `result=$(whence git)`
// (ZC1005). A single pass would leave the inner stale.
func applyFixesUntilStable(src string, initialEdits []katas.FixEdit, registry *katas.KatasRegistry, disabled []string, cfg config.Config, allowedSeverities []katas.Severity, maxPasses int, unsafe bool) (string, int, error) {
	if maxPasses < 1 {
		maxPasses = 5
	}
	current := src
	totalEdits := 0
	edits := initialEdits
	for pass := 0; pass < maxPasses; pass++ {
		if len(edits) == 0 {
			break
		}
		next, err := fix.Apply(current, edits)
		if err != nil {
			return current, totalEdits, err
		}
		if next == current {
			break
		}
		applied := len(edits)
		// Safety net: never accept a pass that introduces parse errors.
		// Two fixes can collide on adjacent spans (for example ZC1073
		// deleting the `$` while another rewrites the same token) and
		// produce broken source. Rather than write it, keep only the
		// edits that are individually safe.
		if parseErrorCount(next) > parseErrorCount(current) {
			next, applied = applySafeEdits(current, edits)
			if next == current {
				break
			}
		}
		totalEdits += applied
		current = next
		// Re-collect edits from the new source.
		edits = collectEdits(current, registry, disabled, cfg, allowedSeverities, unsafe)
	}
	return current, totalEdits, nil
}

// parseErrorCount reports how many parser errors src produces. Used by
// the auto-fix safety net to reject any rewrite that breaks the parse.
func parseErrorCount(src string) int {
	p := parser.New(lexer.New(src))
	p.ParseProgram()
	return len(p.Errors())
}

// applySafeEdits is the fallback when the full edit batch breaks the
// parse. It applies edits one at a time, highest source offset first, and
// keeps an edit only when it does not raise the parser-error count of the
// accumulated result. Highest-offset-first ordering means accepting an
// edit never shifts the line/column of the edits not yet tried, so each
// stays valid against the growing source. A failed apply or a fix that
// would break the parse is simply skipped, so a broken file is never
// written.
func applySafeEdits(base string, edits []katas.FixEdit) (string, int) {
	baseErrs := parseErrorCount(base)
	ordered := append([]katas.FixEdit(nil), edits...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].Line != ordered[j].Line {
			return ordered[i].Line > ordered[j].Line
		}
		return ordered[i].Column > ordered[j].Column
	})
	acc := base
	applied := 0
	for _, e := range ordered {
		trial, err := fix.Apply(acc, []katas.FixEdit{e})
		if err != nil || parseErrorCount(trial) > baseErrs {
			continue
		}
		acc = trial
		applied++
	}
	return acc, applied
}

// collectEdits parses src and returns the auto-fix edits the registry
// would emit for it under the given disabled / severity filters.
// Used by the multi-pass loop in applyFixesUntilStable.
func collectEdits(src string, registry *katas.KatasRegistry, disabled []string, cfg config.Config, allowedSeverities []katas.Severity, unsafe bool) []katas.FixEdit {
	l := lexer.New(src)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return nil
	}
	directives := config.ParseDirectives(src)
	allDisabled := disabled
	if len(directives.File) > 0 {
		allDisabled = append(append([]string(nil), disabled...), directives.File...)
	}
	var violations []katas.Violation
	var edits []katas.FixEdit
	ast.Walk(program, func(node ast.Node) bool {
		vs, es := registry.CheckAndFix(node, allDisabled, []byte(src))
		violations = append(violations, vs...)
		edits = append(edits, es...)
		return true
	})
	if len(directives.PerLine) > 0 {
		keptV := violations[:0]
		keptE := edits[:0]
		for i, v := range violations {
			if directives.IsDisabledOn(v.KataID, v.Line) {
				continue
			}
			keptV = append(keptV, v)
			if i < len(edits) {
				keptE = append(keptE, edits[i])
			}
		}
		violations = keptV
		edits = keptE
	}
	if len(allowedSeverities) > 0 {
		filtered := edits[:0]
		filteredV := violations[:0]
		for i, v := range violations {
			for _, s := range allowedSeverities {
				if v.Level == s {
					filteredV = append(filteredV, v)
					if i < len(edits) {
						filtered = append(filtered, edits[i])
					}
					break
				}
			}
		}
		edits = filtered
	}
	return applicableEdits(edits, registry, unsafe)
}

func processPath(path string, out, errOut io.Writer, cfg config.Config, registry *katas.KatasRegistry, format string, allowedSeverities []katas.Severity, fixOpts fixOptions) int {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(errOut, "Error stating path %s: %s\n", path, err)
		return 0
	}

	count := 0
	if info.IsDir() {
		err := filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if strings.HasPrefix(d.Name(), ".") && d.Name() != "." && d.Name() != ".." {
					return filepath.SkipDir // Skip hidden directories like .git
				}
				return nil
			}

			// Skip non-shell files to avoid parsing errors on Go source code etc.
			ext := filepath.Ext(d.Name())
			if ext == ".go" || ext == ".md" || ext == ".json" || ext == ".yml" || ext == ".yaml" || ext == ".txt" {
				return nil
			}

			// Process only files that look like shell scripts?
			// For now, let's try to parse everything, or maybe filter by extension/shebang if it gets too noisy.
			// Shellcheck defaults to checking all files passed, but for recursive it might filter.
			// Let's assume user wants to check all files in the dir if they passed the dir.
			count += processFile(p, out, errOut, cfg, registry, format, allowedSeverities, fixOpts)
			return nil
		})
		if err != nil {
			fmt.Fprintf(errOut, "Error walking directory %s: %s\n", path, err)
		}
	} else {
		count += processFile(path, out, errOut, cfg, registry, format, allowedSeverities, fixOpts)
	}
	return count
}

func processFile(filename string, out, errOut io.Writer, cfg config.Config, registry *katas.KatasRegistry, format string, allowedSeverities []katas.Severity, fixOpts fixOptions) int {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(errOut, "Error reading file %s: %s\n", filename, err)
		return 0
	}
	program, errs := parseSource(data)
	if len(errs) != 0 {
		for _, msg := range errs {
			fmt.Fprintf(errOut, "Parser Error in %s: %s\n", filename, msg)
		}
		return 1
	}
	directives := config.ParseDirectives(string(data))
	disabled := mergeDisabled(cfg.DisabledKatas, directives.File)

	violations, edits := collectViolations(program, registry, disabled, data, fixOpts.enabled)
	regradeSeverity(violations, fixOpts.ruleSeverity)
	// Stale-suppression detection compares the raw findings against the
	// `# noka` directives before any are silenced.
	if fixOpts.detectStale {
		reportStaleNoka(out, filename, violations, directives, fixOpts.staleCount)
	}
	violations, edits = applyDirectiveSilences(violations, edits, directives)

	// add-noka rewrites the file in place, silencing every current finding,
	// and stops short of severity filtering, fixing, or reporting.
	if fixOpts.addNoka {
		if err := addNokaDirectives(filename, data, violations); err != nil {
			fmt.Fprintf(errOut, "add-noka: %s\n", err)
		}
		return len(violations)
	}

	violations, edits = applySeverityFilter(violations, edits, allowedSeverities)

	// The baseline ratchet records or suppresses findings against a saved
	// snapshot. Write mode collects them and stops short of fixing or
	// reporting; filter mode leaves only findings new since the baseline.
	if fixOpts.baseline != nil {
		violations = fixOpts.baseline.applyBaseline(filename, data, violations)
		if fixOpts.baseline.write {
			return len(violations)
		}
	}

	applyFixIfEnabled(filename, data, registry, disabled, cfg, allowedSeverities, edits, violations, fixOpts, out, errOut)
	emitReport(filename, out, errOut, format, cfg, violations, data, registry, fixOpts)
	return len(violations)
}

func parseSource(data []byte) (*ast.Program, []string) {
	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()
	return program, p.Errors()
}

func mergeDisabled(base, extra []string) []string {
	if len(extra) == 0 {
		return base
	}
	return append(append([]string(nil), base...), extra...)
}

func collectViolations(program *ast.Program, registry *katas.KatasRegistry, disabled []string, data []byte, withFix bool) ([]katas.Violation, []katas.FixEdit) {
	violations := []katas.Violation{}
	var edits []katas.FixEdit
	ast.Walk(program, func(node ast.Node) bool {
		if withFix {
			vs, es := registry.CheckAndFix(node, disabled, data)
			violations = append(violations, vs...)
			edits = append(edits, es...)
		} else {
			violations = append(violations, registry.Check(node, disabled)...)
		}
		return true
	})
	return violations, edits
}

func applyDirectiveSilences(violations []katas.Violation, edits []katas.FixEdit, directives config.Directives) ([]katas.Violation, []katas.FixEdit) {
	if len(directives.PerLine) == 0 {
		return violations, edits
	}
	kept := violations[:0]
	for _, v := range violations {
		if directives.IsDisabledOn(v.KataID, v.Line) {
			continue
		}
		kept = append(kept, v)
	}
	return kept, edits
}

func applySeverityFilter(violations []katas.Violation, edits []katas.FixEdit, allowed []katas.Severity) ([]katas.Violation, []katas.FixEdit) {
	if len(allowed) == 0 {
		return violations, edits
	}
	var filtered []katas.Violation
	for _, v := range violations {
		for _, s := range allowed {
			if v.Level == s {
				filtered = append(filtered, v)
				break
			}
		}
	}
	return filtered, edits
}

// applicableEdits drops behavior-changing (unsafe) fixes unless the run
// opted into them with -unsafe-fixes.
func applicableEdits(edits []katas.FixEdit, registry *katas.KatasRegistry, unsafe bool) []katas.FixEdit {
	if unsafe {
		return edits
	}
	kept := edits[:0]
	for _, e := range edits {
		if registry.IsSafeFix(e.KataID) {
			kept = append(kept, e)
		}
	}
	return kept
}

func applyFixIfEnabled(filename string, data []byte, registry *katas.KatasRegistry, disabled []string, cfg config.Config, allowed []katas.Severity, edits []katas.FixEdit, violations []katas.Violation, fixOpts fixOptions, out, errOut io.Writer) {
	edits = applicableEdits(edits, registry, fixOpts.unsafe)
	if !fixOpts.enabled || len(edits) == 0 || len(violations) == 0 {
		return
	}
	if fixOpts.diff {
		emitFixDiff(filename, data, edits, out, errOut)
	} else if !fixOpts.dryRun {
		applyFixInPlace(filename, data, registry, disabled, cfg, allowed, edits, fixOpts, errOut)
	}
	if fixOpts.stats != nil {
		fixOpts.stats.filesScanned++
	}
}

func emitFixDiff(filename string, data []byte, edits []katas.FixEdit, out, errOut io.Writer) {
	diff, derr := fix.Diff(filename, string(data), edits)
	if derr != nil {
		fmt.Fprintf(errOut, "fix: diff failed for %s: %s\n", filename, derr)
		return
	}
	if diff != "" {
		fmt.Fprint(out, diff)
	}
}

func applyFixInPlace(filename string, data []byte, registry *katas.KatasRegistry, disabled []string, cfg config.Config, allowed []katas.Severity, edits []katas.FixEdit, fixOpts fixOptions, errOut io.Writer) {
	fixed, totalEdits, perr := applyFixesUntilStable(string(data), edits, registry, disabled, cfg, allowed, fixOpts.maxPasses, fixOpts.unsafe)
	if perr != nil {
		fmt.Fprintf(errOut, "fix: apply failed for %s: %s\n", filename, perr)
		return
	}
	if fixed == string(data) {
		return
	}
	mode := os.FileMode(0o600)
	if info, statErr := os.Stat(filename); statErr == nil {
		mode = info.Mode().Perm()
	}
	if werr := os.WriteFile(filename, []byte(fixed), mode); werr != nil {
		fmt.Fprintf(errOut, "fix: write failed for %s: %s\n", filename, werr)
		return
	}
	fmt.Fprintf(errOut, "fixed %d edit(s) in %s\n", totalEdits, filename)
	if fixOpts.stats != nil {
		fixOpts.stats.filesModified++
		fixOpts.stats.totalEdits += totalEdits
	}
}

func emitReport(filename string, out, errOut io.Writer, format string, cfg config.Config, violations []katas.Violation, data []byte, registry *katas.KatasRegistry, fixOpts fixOptions) {
	if len(violations) == 0 {
		return
	}
	// Statistics mode tallies findings per kata and suppresses the
	// individual reports.
	if fixOpts.statistics != nil {
		for _, v := range violations {
			fixOpts.statistics[v.KataID]++
		}
		return
	}
	// The machine-readable formats are aggregated and emitted once by the
	// caller; collect this file's findings and return. Text reports inline.
	if fixOpts.collector != nil {
		*fixOpts.collector = append(*fixOpts.collector, reporter.FileViolations{Filename: filename, Violations: violations})
		return
	}
	// `[*]` marks the fixes the current command would apply: safe-only, or
	// every fix under -unsafe-fixes. Findings whose only fix is unsafe are
	// counted separately so the footer can point at -unsafe-fixes.
	marked := registry.IsSafeFix
	if fixOpts.unsafe {
		marked = registry.IsFixable
	}
	if fixOpts.fixable != nil {
		for _, v := range violations {
			switch {
			case marked(v.KataID):
				*fixOpts.fixable++
			case registry.IsFixable(v.KataID) && fixOpts.unsafeFixable != nil:
				*fixOpts.unsafeFixable++
			}
		}
	}
	r := reporter.NewTextReporter(out, filename, string(data), cfg)
	r.MarkFixable(marked)
	if err := r.Report(violations); err != nil {
		fmt.Fprintf(errOut, "Error reporting violations: %s\n", err)
	}
}
