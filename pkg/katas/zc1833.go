package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1833",
		Title:    "Warn on `unsetopt WARN_CREATE_GLOBAL` — silent accidental-global bugs inside functions",
		Severity: SeverityWarning,
		Description: "`WARN_CREATE_GLOBAL` makes Zsh warn when a function assigns to a name " +
			"that is not declared `local` / `typeset` in the current scope — the single " +
			"highest-value guardrail against the classic Bash-ism where a helper function " +
			"silently stomps on a caller's variable (`tmp=`, `i=`, `result=`). Disabling it " +
			"(`unsetopt WARN_CREATE_GLOBAL` or the equivalent `setopt NO_WARN_CREATE_GLOBAL`) " +
			"reverts to permissive behaviour: every unqualified assignment inside a function " +
			"escapes to global scope with no diagnostic. Leave the option on and fix the " +
			"offending function by adding `local` / `typeset` declarations, or — if you " +
			"really must silence it for a specific block — use `setopt LOCAL_OPTIONS; " +
			"unsetopt WARN_CREATE_GLOBAL` inside a function so the rest of the script keeps " +
			"the safety.",
		Check: checkZC1833,
	})
}

func checkZC1833(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1833IsWarnCreateGlobal(arg.String()) {
				return zc1833Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOWARNCREATEGLOBAL" {
				return zc1833Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1833IsWarnCreateGlobal(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "WARNCREATEGLOBAL"
}

func zc1833Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1833",
		Message: "`" + where + "` silences Zsh's warning for assignments leaking " +
			"out of function scope — classic caller-variable stomping. Declare " +
			"`local`/`typeset`; scope with `LOCAL_OPTIONS` if you must disable.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
