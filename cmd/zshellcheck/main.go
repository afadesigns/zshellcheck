package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"

	"github.com/adrg/xdg"
	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/fix"
	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
	"github.com/afadesigns/zshellcheck/pkg/reporter"
	"github.com/afadesigns/zshellcheck/pkg/version"
	"gopkg.in/yaml.v3"
)

func main() {
	os.Exit(run())
}

func run() int {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, config.Banner)
		fmt.Fprintf(os.Stderr, "ZShellCheck - The Zsh Static Analysis Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage: zshellcheck [flags] <file1.zsh> [file2.zsh]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  zshellcheck script.zsh\n")
		fmt.Fprintf(os.Stderr, "  zshellcheck -format json script.zsh\n")
		fmt.Fprintf(os.Stderr, "  zshellcheck ./scripts/\n")
	}

	format := flag.String("format", "text", "The output format (text, json, or sarif)")
	cpuprofile := flag.String("cpuprofile", "", "Write CPU profile to file")
	showVersion := flag.Bool("version", false, "Show version and exit")
	verbose := flag.Bool("verbose", false, "Show detailed Kata descriptions in text output")
	noColor := flag.Bool("no-color", false, "Disable colored output")
	severityFilter := flag.String("severity", "", "Comma-separated list of severities to show (error,warning,info,style)")
	fixMode := flag.Bool("fix", false, "Apply auto-fixes in place for katas that declare a Fix")
	diffMode := flag.Bool("diff", false, "Print a unified diff of the fixes instead of writing them")
	dryRun := flag.Bool("dry-run", false, "When used with -fix, do not modify files; just list what would change")
	flag.Parse()

	if *showVersion {
		fmt.Printf("zshellcheck version %s\n", version.Version)
		return 0
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create CPU profile: %s\n", err)
			return 1
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Could not start CPU profile: %s\n", err)
			return 1
		}
		defer pprof.StopCPUProfile()
	}

	if len(flag.Args()) < 1 {
		fmt.Fprint(os.Stderr, config.Banner)
		fmt.Println("Usage: zshellcheck [flags] <file1.zsh> [file2.zsh]...")
		fmt.Println("Try 'zshellcheck --help' for more information.")
		return 1
	}

	// Print banner on successful run too, as per original request
	// But suppress it for JSON/SARIF output to keep it clean for parsing
	if *format != "json" && *format != "sarif" && !*noColor {
		fmt.Fprint(os.Stderr, config.Banner)
	}

	// Pick either config.yml or config.yaml, but not both
	configFileXdgConfigPath, err := xdg.SearchConfigFile("zshellcheck/config.yml")
	if err != nil {
		configFileXdgConfigPath, _ = xdg.SearchConfigFile("zshellcheck/config.yaml")
	}

	configFileHomePath := filepath.Join(xdg.Home, ".zshellcheckrc")
	config, err := loadConfig(configFileXdgConfigPath, configFileHomePath, ".zshellcheckrc")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
		return 1
	}

	if *noColor {
		config.NoColor = true
	}

	if *verbose {
		config.Verbose = true
	}

	var allowedSeverities []katas.Severity
	if *severityFilter != "" {
		for _, s := range strings.Split(*severityFilter, ",") {
			sTrimmed := katas.Severity(strings.TrimSpace(s))
			switch sTrimmed {
			case katas.SeverityError, katas.SeverityWarning, katas.SeverityInfo, katas.SeverityStyle:
				allowedSeverities = append(allowedSeverities, sTrimmed)
			default:
				fmt.Fprintf(os.Stderr, "Invalid severity level: %s. Must be one of error, warning, info, style.\n", s)
				return 1
			}
		}
	}

	kataRegistry := katas.Registry

	fixOpts := fixOptions{
		enabled:   *fixMode || *diffMode,
		diff:      *diffMode,
		dryRun:    *dryRun,
		maxPasses: 5,
	}
	if fixOpts.diff {
		// -diff is equivalent to -fix -dry-run with diff rendering.
		fixOpts.enabled = true
		fixOpts.dryRun = true
	}
	if fixOpts.enabled {
		fixOpts.stats = &fixStats{}
	}

	totalViolations := 0
	for _, filename := range flag.Args() {
		totalViolations += processPath(filename, os.Stdout, os.Stderr, config, kataRegistry, *format, allowedSeverities, fixOpts)
	}

	if fixOpts.stats != nil && fixOpts.stats.filesScanned > 1 {
		fmt.Fprintf(os.Stderr, "\nfix summary: %d edit(s) across %d file(s) (scanned %d)\n",
			fixOpts.stats.totalEdits, fixOpts.stats.filesModified, fixOpts.stats.filesScanned)
	}

	if totalViolations > 0 {
		if *format == "text" {
			fmt.Fprintf(os.Stderr, "\nFound %d violations.\n", totalViolations)
		}
		return 1
	}
	return 0
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

		var fileConfig config.Config
		err = yaml.Unmarshal(data, &fileConfig)
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
func applyFixesUntilStable(src string, initialEdits []katas.FixEdit, registry *katas.KatasRegistry, disabled []string, cfg config.Config, allowedSeverities []katas.Severity, maxPasses int) (string, int, error) {
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
		totalEdits += len(edits)
		current = next
		// Re-collect edits from the new source.
		edits = collectEdits(current, registry, disabled, cfg, allowedSeverities)
	}
	return current, totalEdits, nil
}

// collectEdits parses src and returns the auto-fix edits the registry
// would emit for it under the given disabled / severity filters.
// Used by the multi-pass loop in applyFixesUntilStable.
func collectEdits(src string, registry *katas.KatasRegistry, disabled []string, cfg config.Config, allowedSeverities []katas.Severity) []katas.FixEdit {
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
	return edits
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

	l := lexer.New(string(data))
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Fprintf(errOut, "Parser Error in %s: %s\n", filename, msg)
		}
		// Parser errors should technically count as failures too?
		// But for now let's just count violations.
		// Actually, if parser fails, we probably want to fail build.
		return 1
	}

	directives := config.ParseDirectives(string(data))
	disabled := cfg.DisabledKatas
	if len(directives.File) > 0 {
		disabled = append(append([]string(nil), disabled...), directives.File...)
	}

	violations := []katas.Violation{}
	var edits []katas.FixEdit
	ast.Walk(program, func(node ast.Node) bool {
		if fixOpts.enabled {
			vs, es := registry.CheckAndFix(node, disabled, data)
			violations = append(violations, vs...)
			edits = append(edits, es...)
		} else {
			violations = append(violations, registry.Check(node, disabled)...)
		}
		return true // Continue walking
	})

	// Drop violations silenced by a per-line `# zshellcheck disable=…` directive.
	if len(directives.PerLine) > 0 {
		kept := violations[:0]
		for _, v := range violations {
			if directives.IsDisabledOn(v.KataID, v.Line) {
				continue
			}
			kept = append(kept, v)
		}
		violations = kept
	}

	// Filter violations by severity
	var filteredViolations []katas.Violation
	if len(allowedSeverities) > 0 {
		for _, v := range violations {
			for _, s := range allowedSeverities {
				if v.Level == s {
					filteredViolations = append(filteredViolations, v)
					break
				}
			}
		}
		violations = filteredViolations
	}

	// Apply auto-fixes when requested. Edits land only when the file has
	// surviving violations (silenced / severity-filtered violations drop
	// their associated fixes too — silence should silence the fix).
	if fixOpts.enabled && len(edits) > 0 && len(violations) > 0 {
		if fixOpts.diff {
			diff, derr := fix.Diff(filename, string(data), edits)
			if derr != nil {
				fmt.Fprintf(errOut, "fix: diff failed for %s: %s\n", filename, derr)
			} else if diff != "" {
				fmt.Fprint(out, diff)
			}
		} else if !fixOpts.dryRun {
			// Multi-pass fix: apply edits, re-parse, re-collect, re-apply
			// until the source stops changing or a small iteration cap is
			// hit. Many fixes resolve nested patterns (e.g. backtick →
			// $() exposes a `which` inside that ZC1005 then rewrites to
			// `whence`); a single Apply leaves them stranded.
			fixed, totalEdits, perr := applyFixesUntilStable(string(data), edits, registry, disabled, cfg, allowedSeverities, fixOpts.maxPasses)
			if perr != nil {
				fmt.Fprintf(errOut, "fix: apply failed for %s: %s\n", filename, perr)
			} else if fixed != string(data) {
				mode := os.FileMode(0o600)
				if info, statErr := os.Stat(filename); statErr == nil {
					mode = info.Mode().Perm()
				}
				if werr := os.WriteFile(filename, []byte(fixed), mode); werr != nil {
					fmt.Fprintf(errOut, "fix: write failed for %s: %s\n", filename, werr)
				} else {
					fmt.Fprintf(errOut, "fixed %d edit(s) in %s\n", totalEdits, filename)
					if fixOpts.stats != nil {
						fixOpts.stats.filesModified++
						fixOpts.stats.totalEdits += totalEdits
					}
				}
			}
		}
		if fixOpts.stats != nil {
			fixOpts.stats.filesScanned++
		}
	}

	if len(violations) > 0 {
		var r reporter.Reporter
		switch format {
		case "json":
			r = reporter.NewJSONReporter(out)
		case "sarif":
			r = reporter.NewSarifReporter(out, filename)
		default:
			r = reporter.NewTextReporter(out, filename, string(data), cfg)
		}
		if err := r.Report(violations); err != nil {
			fmt.Fprintf(errOut, "Error reporting violations: %s\n", err)
		}
	}
	return len(violations)
}
