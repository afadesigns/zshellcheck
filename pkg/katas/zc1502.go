package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1502",
		Title:    "Warn on `grep \"$var\" file` without `--` — flag injection when `$var` starts with `-`",
		Severity: SeverityWarning,
		Description: "Without a `--` end-of-flags marker, `grep` (and most POSIX tools) treats " +
			"any argument that starts with `-` as a flag. If `$var` comes from user input or a " +
			"fuzzed filename, an attacker can pass `--include=*secret*` or `-f /etc/shadow` " +
			"and get grep to read paths the script author never intended. Always write " +
			"`grep -- \"$var\" file` or use a grep-compatible library with explicit pattern API.",
		Check: checkZC1502,
	})
}

func checkZC1502(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "grep" && ident.Value != "egrep" && ident.Value != "fgrep" &&
		ident.Value != "rg" && ident.Value != "ag" {
		return nil
	}

	hasDashDash := false
	firstVarArg := ""
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--" {
			hasDashDash = true
			break
		}
		// First variable-like pattern argument (no surrounding single quote means interpolation)
		if firstVarArg == "" && (strings.HasPrefix(v, "\"$") || strings.HasPrefix(v, "$")) {
			firstVarArg = v
		}
	}

	if firstVarArg == "" || hasDashDash {
		return nil
	}
	return []Violation{{
		KataID: "ZC1502",
		Message: "Variable `" + firstVarArg + "` used as pattern without `--` end-of-flags " +
			"marker — attacker-controlled leading `-` becomes a flag. Write `grep -- \"$var\"`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
