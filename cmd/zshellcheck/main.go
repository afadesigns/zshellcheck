package main

import (
	"flag"
	"fmt"
	"io"
	"os"

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
	banner := "\n" +
		"\033[38;5;217m      ____\033[0m\033[38;5;229m______\033[0m\033[38;5;245m__  __   ____   __              __  \033[0m\n" +
		"\033[38;5;217m     /_  /\033[0m\033[38;5;229m  ___/\033[0m\033[38;5;245m /_  /  / __/  / /  ___   _____/ /__\033[0m\n" +
		"\033[38;5;217m      / /_\033[0m\033[38;5;229m___ \\ \033[0m\033[38;5;245m/ __ \\ / _ \\ / /  / _ \\ / ___/ //_/\033[0m\n" +
		"\033[38;5;217m     / /_\033[0m\033[38;5;229m___/ /\033[0m\033[38;5;245m / / //  __// /__/  __// /__/ ,<   \033[0m\n" +
		"\033[38;5;217m    /___/\033[0m\033[38;5;229m____/\033[0m\033[38;5;245m_/_/ /_/\\___//____/\\___/ \\___/_/|_|   \033[0m\n" +
		"\n"

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, banner)
		fmt.Fprintf(os.Stderr, "ZShellCheck - The Zsh Static Analysis Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage: zshellcheck [flags] <file1.zsh> [file2.zsh]...\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  zshellcheck script.zsh\n")
		fmt.Fprintf(os.Stderr, "  zshellcheck -format json script.zsh\n")
		fmt.Fprintf(os.Stderr, "  zshellcheck ./scripts/\n")
	}

	format := flag.String("format", "text", "The output format (text or json)")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Fprint(os.Stderr, banner)
		fmt.Println("Usage: zshellcheck [flags] <file1.zsh> [file2.zsh]...")
		fmt.Println("Try 'zshellcheck --help' for more information.")
		os.Exit(1)
	}

	// Print banner on successful run too, as per original request
	fmt.Fprint(os.Stderr, banner)

	config, err := loadConfig(".zshellcheckrc")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
		os.Exit(1)
	}

	kataRegistry := &katas.Registry

	for _, filename := range flag.Args() {
		processFile(filename, os.Stdout, os.Stderr, config, kataRegistry, *format)
	}
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
		default:
			fmt.Fprintf(out, "Violations in %s:\n", filename)
			r = reporter.NewTextReporter(out)
		}
		if err := r.Report(violations); err != nil {
			fmt.Fprintf(errOut, "Error reporting violations: %s\n", err)
		}
	}
}
