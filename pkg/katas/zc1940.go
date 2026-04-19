package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1940",
		Title:    "Warn on `setopt POSIX_ARGZERO` — `$0` no longer changes to the function name inside functions",
		Severity: SeverityWarning,
		Description: "Zsh's default behaviour (option off) assigns `$0` to the name of the " +
			"currently-running function, so a helper like `log() { printf '%s\\n' \"$0: $*\"; }` " +
			"prints `log: …`. `setopt POSIX_ARGZERO` keeps `$0` pointing at the outer script " +
			"name (or the interpreter when sourced) — the logger instead prints the script " +
			"path for every message and call-site context is lost. Every `case $0` dispatch " +
			"inside an auto-loaded function also stops working. Leave the option off; if you " +
			"need POSIX `$0`, scope it in a function with `emulate -LR sh`.",
		Check: checkZC1940,
	})
}

func checkZC1940(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1940Canonical(arg.String())
		switch v {
		case "POSIXARGZERO":
			if enabling {
				return zc1940Hit(cmd, "setopt POSIX_ARGZERO")
			}
		case "NOPOSIXARGZERO":
			if !enabling {
				return zc1940Hit(cmd, "unsetopt NO_POSIX_ARGZERO")
			}
		}
	}
	return nil
}

func zc1940Canonical(s string) string {
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

func zc1940Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1940",
		Message: "`" + form + "` freezes `$0` to the outer script name — loggers and " +
			"`case $0` dispatch inside functions lose call-site context. Scope with " +
			"`emulate -LR sh` instead of flipping globally.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
