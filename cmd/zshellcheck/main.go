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
	format := flag.String("format", "text", "The output format (text or json)")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: zshellcheck [flags] <file1.zsh> [file2.zsh]...")
		os.Exit(1)
	}

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
