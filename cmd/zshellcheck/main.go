package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/config"
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
	severityFilter := flag.String("severity", "", "Comma-separated list of severities to show (Error,Warning,Info)")
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

	config, err := loadConfig(".zshellcheckrc")
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
			case katas.Error, katas.Warning, katas.Info:
				allowedSeverities = append(allowedSeverities, sTrimmed)
			default:
				fmt.Fprintf(os.Stderr, "Invalid severity level: %s. Must be one of Error, Warning, Info.\n", s)
				return 1
			}
		}
	}

	kataRegistry := katas.Registry

	totalViolations := 0
	for _, filename := range flag.Args() {
		totalViolations += processPath(filename, os.Stdout, os.Stderr, config, kataRegistry, *format, allowedSeverities)
	}

	if totalViolations > 0 {
		if *format == "text" {
			fmt.Fprintf(os.Stderr, "\nFound %d violations.\n", totalViolations)
		}
		return 1
	}
	return 0
}

func loadConfig(path string) (config.Config, error) {
	cfg := config.DefaultConfig()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
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

	return config.MergeConfig(cfg, fileConfig), nil
}

func processPath(path string, out, errOut io.Writer, config config.Config, registry *katas.KatasRegistry, format string, allowedSeverities []katas.Severity) int {
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
			count += processFile(p, out, errOut, config, registry, format, allowedSeverities)
			return nil
		})
		if err != nil {
			fmt.Fprintf(errOut, "Error walking directory %s: %s\n", path, err)
		}
	} else {
		count += processFile(path, out, errOut, config, registry, format, allowedSeverities)
	}
	return count
}

func processFile(filename string, out, errOut io.Writer, config config.Config, registry *katas.KatasRegistry, format string, allowedSeverities []katas.Severity) int {
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

	violations := []katas.Violation{}
	ast.Walk(program, func(node ast.Node) bool {
		violations = append(violations, registry.Check(node, config.DisabledKatas)...)
		return true // Continue walking
	})

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

	if len(violations) > 0 {
		var r reporter.Reporter
		switch format {
		case "json":
			r = reporter.NewJSONReporter(out)
		case "sarif":
			r = reporter.NewSarifReporter(out, filename)
		default:
			r = reporter.NewTextReporter(out, filename, string(data), config)
		}
		if err := r.Report(violations); err != nil {
			fmt.Fprintf(errOut, "Error reporting violations: %s\n", err)
		}
	}
	return len(violations)
}
