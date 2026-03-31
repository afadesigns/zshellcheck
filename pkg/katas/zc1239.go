package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1239",
		Title:    "Avoid `kubectl exec -it` in scripts",
		Severity: SeverityWarning,
		Description: "`kubectl exec -it` allocates a TTY which hangs in non-interactive scripts. " +
			"Use `kubectl exec` without `-it` or use `kubectl exec -- cmd` for scripted commands.",
		Check: checkZC1239,
	})
}

func checkZC1239(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubectl" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "exec" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "-it" || val == "-ti" {
			return []Violation{{
				KataID: "ZC1239",
				Message: "Avoid `kubectl exec -it` in scripts — TTY allocation hangs without a terminal. " +
					"Use `kubectl exec pod -- cmd` for non-interactive execution.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
