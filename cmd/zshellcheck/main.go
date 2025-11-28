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

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
	"github.com/afadesigns/zshellcheck/pkg/reporter"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DisabledKatas []string `yaml:"disabled_katas"`
}

func main() {
	os.Exit(run())
}

func run() int {
	banner := "\n" +
		"\033[38;5;51m███████╗███████╗██╗  ██╗███████╗██╗     ██╗      ██████╗██╗  ██╗███████╗ ██████╗██╗  ██╗\033[0m\n" +
		"\033[38;5;45m╚══███╔╝██╔════╝██║  ██║██╔════╝██║     ██║     ██╔════╝██║  ██║██╔════╝██╔════╝██║ ██╔╝\033[0m\n" +
		"\033[38;5;39m  ███╔╝ ███████╗███████║█████╗  ██║     ██║     ██║     ███████║█████╗  ██║     █████╔╝\033[0m\n" +
		"\033[38;5;33m ███╔╝  ╚════██║██╔══██║██╔══╝  ██║     ██║     ██║     ██╔══██║██╔══╝  ██║     ██╔═██╗\033[0m\n" +
		"\033[38;5;27m███████╗███████║██║  ██║███████╗███████╗███████╗╚██████╗██║  ██║███████╗╚██████╗██║  ██╗\033[0m\n" +
		"\033[38;5;21m╚══════╝╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝ ╚═════╝╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝\033[0m\n" +
		"\n"

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, banner)
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
	flag.Parse()

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
		fmt.Fprint(os.Stderr, banner)
		fmt.Println("Usage: zshellcheck [flags] <file1.zsh> [file2.zsh]...")
		fmt.Println("Try 'zshellcheck --help' for more information.")
		return 1
	}

	// Print banner on successful run too, as per original request
	// But suppress it for JSON/SARIF output to keep it clean for parsing
	if *format != "json" && *format != "sarif" {
		fmt.Fprint(os.Stderr, banner)
	}

	config, err := loadConfig(".zshellcheckrc")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
		return 1
	}

	kataRegistry := katas.Registry

	for _, filename := range flag.Args() {
		processPath(filename, os.Stdout, os.Stderr, config, kataRegistry, *format)
	}
	return 0
}

func loadConfig(path string) (Config, error) {
	var config Config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func processPath(path string, out, errOut io.Writer, config Config, registry *katas.KatasRegistry, format string) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(errOut, "Error stating path %s: %s\n", path, err)
		return
	}

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
			processFile(p, out, errOut, config, registry, format)
			return nil
		})
		if err != nil {
			fmt.Fprintf(errOut, "Error walking directory %s: %s\n", path, err)
		}
	} else {
		processFile(path, out, errOut, config, registry, format)
	}
}

func processFile(filename string, out, errOut io.Writer, config Config, registry *katas.KatasRegistry, format string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(errOut, "Error reading file %s: %s\n", filename, err)
		return
	}

	l := lexer.New(string(data))
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Fprintf(errOut, "Parser Error in %s: %s\n", filename, msg)
		}
		return
	}

	violations := []katas.Violation{}
	ast.Walk(program, func(node ast.Node) bool {
		violations = append(violations, registry.Check(node, config.DisabledKatas)...)
		return true // Continue walking
	})

	if len(violations) > 0 {
		var r reporter.Reporter
		switch format {
		case "json":
			r = reporter.NewJSONReporter(out)
		case "sarif":
			r = reporter.NewSarifReporter(out, filename)
		default:
			r = reporter.NewTextReporter(out, filename)
		}
		if err := r.Report(violations); err != nil {
			fmt.Fprintf(errOut, "Error reporting violations: %s\n", err)
		}
	}
}
