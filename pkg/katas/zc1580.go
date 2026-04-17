package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1580",
		Title:    "Warn on `go build -ldflags \"-X main.<SECRET>=...\"` — secret embedded in binary",
		Severity: SeverityWarning,
		Description: "`-ldflags=\"-X pkg.Var=value\"` sets a Go string variable at link time. " +
			"Putting a secret here bakes it into the resulting binary (discoverable with " +
			"`strings`, `objdump`, or simply opening the file). It also leaves the value on " +
			"the build host's shell history and in any CI transcript. Read the value at " +
			"runtime from `os.Getenv` / a mounted secret file / the cloud secret manager.",
		Check: checkZC1580,
	})
}

func checkZC1580(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "go" {
		return nil
	}

	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "build" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := strings.Trim(arg.String(), "\"'")
		// Look for -ldflags=... content or "-X ..." content
		if !strings.Contains(v, "-X ") && !strings.Contains(v, "-ldflags") {
			continue
		}
		up := strings.ToUpper(v)
		if strings.Contains(up, "PASSWORD=") || strings.Contains(up, "SECRET=") ||
			strings.Contains(up, "APIKEY=") || strings.Contains(up, "API_KEY=") ||
			strings.Contains(up, "TOKEN=") || strings.Contains(up, "PRIVATE_KEY=") {
			return []Violation{{
				KataID: "ZC1580",
				Message: "`go build -ldflags` injecting a secret bakes it into the binary. " +
					"Read from os.Getenv / mounted secret file at runtime.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
