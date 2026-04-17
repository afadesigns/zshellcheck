package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1596",
		Title:    "Style: `emulate sh/bash/ksh` without `-L` — flips options for the whole shell",
		Severity: SeverityStyle,
		Description: "`emulate MODE` without the `-L` flag changes Zsh options globally. After " +
			"that line runs the shell is no longer in Zsh mode — `${(F)arr}`, 1-indexed " +
			"arrays, glob qualifiers, and other Zsh-only constructs either error or silently " +
			"behave differently. Wrap emulation in a function and use `emulate -L MODE` to " +
			"scope it to that function. A `.zsh` script that starts with `emulate sh` likely " +
			"belongs in a `.sh` file instead.",
		Check: checkZC1596,
	})
}

func checkZC1596(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "emulate" {
		return nil
	}

	var hasL bool
	var mode string
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "-") {
			if strings.Contains(v, "L") {
				hasL = true
			}
			continue
		}
		if v == "sh" || v == "bash" || v == "ksh" || v == "csh" {
			mode = v
		}
	}
	if mode == "" || hasL {
		return nil
	}

	return []Violation{{
		KataID: "ZC1596",
		Message: "`emulate " + mode + "` without `-L` flips the options for the whole " +
			"shell. Use `emulate -L " + mode + "` inside a function, or rename the script " +
			"to `.sh` if Zsh features are not needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
