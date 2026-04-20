package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1999",
		Title:    "Error on `setopt AUTO_NAMED_DIRS` — unknown option, typo of `AUTO_NAME_DIRS`",
		Severity: SeverityError,
		Description: "`AUTO_NAMED_DIRS` (with the trailing `D`) is not a real Zsh option — " +
			"`setopt AUTO_NAMED_DIRS` fails with `no such option` and the dir-to-" +
			"`~name` auto-registration the author likely wanted is never enabled. " +
			"The canonical spelling is `AUTO_NAME_DIRS` (see ZC1934 for its " +
			"semantics and why flipping it on is usually wrong). Drop the typo and, " +
			"if you actually want the behaviour, reach for `hash -d NAME=PATH` " +
			"explicitly or scope `setopt LOCAL_OPTIONS AUTO_NAME_DIRS` inside the " +
			"single helper that needs named-directory expansion.",
		Check: checkZC1999,
	})
}

func checkZC1999(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setopt" && ident.Value != "unsetopt" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := zc1999Canonical(arg.String())
		switch v {
		case "AUTONAMEDDIRS", "NOAUTONAMEDDIRS":
			return []Violation{{
				KataID: "ZC1999",
				Message: "`" + ident.Value + " " + arg.String() + "` is a typo — the real " +
					"Zsh option is `AUTO_NAME_DIRS` (no trailing `D`, see ZC1934). " +
					"Fix the spelling or drop the toggle; `hash -d NAME=PATH` is " +
					"the explicit alternative.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1999Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}
