package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1591",
		Title:    "Style: use Zsh `print -l` / `${(F)array}` instead of `printf '%s\\n' \"${array[@]}\"`",
		Severity: SeverityStyle,
		Description: "`printf '%s\\n' \"${array[@]}\"` is the Bash-idiomatic way to print one " +
			"element per line. Zsh has `print -l -r -- \"${array[@]}\"` (one element per line, " +
			"raw, sentinel-safe) and the parameter-expansion flag `${(F)array}` (newline-join, " +
			"fine for `$(...)`). Both are shorter than the printf incantation and avoid format-" +
			"string surprises if the array ever contains a literal `%`.",
		Check: checkZC1591,
	})
}

func checkZC1591(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "printf" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	format := strings.Trim(cmd.Arguments[0].String(), "'\"")
	if format != `%s\n` && format != "%s" {
		return nil
	}

	second := cmd.Arguments[1].String()
	if !strings.Contains(second, "[@]") && !strings.Contains(second, "[*]") {
		return nil
	}

	return []Violation{{
		KataID: "ZC1591",
		Message: "`printf '%s\\n' \"${array[@]}\"` — use Zsh `print -l -r -- \"${array[@]}\"` " +
			"or `${(F)array}` for newline-joined output.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
