package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1245",
		Title:    "Avoid disabling TLS certificate verification",
		Severity: SeverityError,
		Description: "Flags like `--no-check-certificate` (wget) or `-k`/`--insecure` (curl) " +
			"disable TLS verification, making connections vulnerable to MITM attacks.",
		Check: checkZC1245,
	})
}

func checkZC1245(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value

	for _, arg := range cmd.Arguments {
		val := arg.String()

		if name == "curl" && (val == "-k" || val == "--insecure") {
			return []Violation{{
				KataID: "ZC1245",
				Message: "Avoid `curl -k`/`--insecure` — it disables TLS certificate verification. " +
					"Fix the certificate chain or use `--cacert` to specify a CA bundle.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}
