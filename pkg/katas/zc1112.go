package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1112",
		Title: "Avoid `grep -c` — use Zsh pattern matching for counting",
		Description: "For counting matches in a variable, use Zsh `${#${(f)...}}` or array filtering " +
			"with `${(M)array:#pattern}` instead of piping through `grep -c`.",
		Severity: SeverityStyle,
		Check:    checkZC1112,
	})
}

func checkZC1112(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	// Only flag grep -c without file arguments (pipeline use)
	hasCountFlag := false
	hasFileAfterPattern := false
	patternSeen := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			if val == "-c" || val == "--count" {
				hasCountFlag = true
			}
		} else {
			if patternSeen {
				hasFileAfterPattern = true
				break
			}
			patternSeen = true
		}
	}

	if !hasCountFlag || hasFileAfterPattern {
		return nil
	}

	return []Violation{{
		KataID: "ZC1112",
		Message: "Use Zsh array filtering `${(M)array:#pattern}` or `${#${(f)...}}` for counting " +
			"instead of `grep -c`. Avoids spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
