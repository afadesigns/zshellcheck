package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1787",
		Title:    "Warn on `setopt AUTO_CD` — bare word that names a directory silently changes `$PWD`",
		Severity: SeverityWarning,
		Description: "With `AUTO_CD` on, any bare word that happens to name an existing directory " +
			"is executed as `cd <word>` — no command name, no error. This is a pleasant " +
			"interactive shortcut and an absolute footgun in scripts: a typo in a command " +
			"name (`doker` → a directory called `doker` that was left lying around) or a " +
			"user-controlled variable that expands to a path silently reshapes `$PWD` for " +
			"every later relative path. Keep `AUTO_CD` inside `~/.zshrc` where it belongs, " +
			"not in a `.zsh` script, and never turn it on inside a function that an external " +
			"caller depends on.",
		Check: checkZC1787,
	})
}

func checkZC1787(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1787IsAutoCd(arg.String()) {
				return zc1787Hit(cmd, "setopt "+arg.String())
			}
		}
	case "set":
		for i, arg := range cmd.Arguments {
			v := arg.String()
			if (v == "-o" || v == "--option") && i+1 < len(cmd.Arguments) {
				next := cmd.Arguments[i+1].String()
				if zc1787IsAutoCd(next) {
					return zc1787Hit(cmd, "set -o "+next)
				}
			}
		}
	}
	return nil
}

func zc1787IsAutoCd(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "AUTOCD"
}

func zc1787Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1787",
		Message: "`" + where + "` turns any bare directory name into a silent `cd`. " +
			"A typo or a user-controlled value reshapes `$PWD`; keep this in " +
			"`~/.zshrc`, not in scripts.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
