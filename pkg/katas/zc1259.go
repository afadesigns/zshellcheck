package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1259",
		Title:    "Avoid `docker pull` without explicit tag — pin image versions",
		Severity: SeverityWarning,
		Description: "`docker pull image` without a tag defaults to `:latest` which is " +
			"mutable and non-reproducible. Always pin to a specific version tag or digest.",
		Check: checkZC1259,
	})
}

func checkZC1259(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	if len(cmd.Arguments) < 2 || cmd.Arguments[0].String() != "pull" {
		return nil
	}

	image := cmd.Arguments[1].String()
	if !strings.Contains(image, ":") && !strings.Contains(image, "@sha256") {
		return []Violation{{
			KataID: "ZC1259",
			Message: "Pin Docker image to a specific tag instead of defaulting to `:latest`. " +
				"Untagged pulls are non-reproducible and may break unexpectedly.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
