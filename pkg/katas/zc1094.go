package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1094",
		Title: "Use parameter expansion instead of `sed` for simple substitutions",
		Description: "For simple string substitutions on variables, use Zsh parameter expansion " +
			"`${var//pattern/replacement}` instead of piping through `sed`. It avoids spawning an external process.",
		Check: checkZC1094,
	})
}

func checkZC1094(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sed" {
		return nil
	}

	// Only flag simple sed invocations: sed 's/pattern/replacement/' or sed 's/pattern/replacement/g'
	// Skip sed with flags like -i (in-place), -n, -e (multiple expressions), -f (script file)
	hasComplexFlags := false
	hasSubstitution := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			hasComplexFlags = true
			break
		}
		if len(val) >= 4 && val[0] == 's' && (val[1] == '/' || val[1] == '|' || val[1] == '#') {
			hasSubstitution = true
		}
		// Also check quoted forms
		if len(val) >= 6 && (val[0] == '\'' || val[0] == '"') {
			inner := val[1 : len(val)-1]
			if len(inner) >= 4 && inner[0] == 's' && (inner[1] == '/' || inner[1] == '|' || inner[1] == '#') {
				hasSubstitution = true
			}
		}
	}

	if hasComplexFlags || !hasSubstitution {
		return nil
	}

	// Only flag if sed has exactly one argument (the substitution pattern)
	if len(cmd.Arguments) != 1 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1094",
		Message: "Use `${var//pattern/replacement}` instead of piping through `sed` for simple substitutions. " +
			"Parameter expansion avoids spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
	}}
}
