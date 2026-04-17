package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1554",
		Title:    "Warn on `unzip -o` / `tar ... --overwrite` — silent overwrite during extract",
		Severity: SeverityWarning,
		Description: "`unzip -o` overwrites existing files without prompting; `tar --overwrite` " +
			"does the same for tarballs. In a directory that already contains user work or a " +
			"previous release, a newer archive silently wins, discarding in-flight edits and " +
			"custom config. Extract to a fresh staging directory, diff, then move specific " +
			"files into place.",
		Check: checkZC1554,
	})
}

func checkZC1554(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "unzip" {
		for _, arg := range cmd.Arguments {
			if arg.String() == "-o" {
				return []Violation{{
					KataID: "ZC1554",
					Message: "`unzip -o` overwrites existing files without prompting. Extract " +
						"to a staging directory, diff, then move.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	if ident.Value == "tar" || ident.Value == "bsdtar" {
		for _, arg := range cmd.Arguments {
			if arg.String() == "--overwrite" {
				return []Violation{{
					KataID: "ZC1554",
					Message: "`tar --overwrite` discards existing files during extract. Use a " +
						"staging directory and diff before rolling forward.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}
