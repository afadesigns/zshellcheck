package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1608",
		Title:    "Warn on `find -exec sh -c '... {} ...'` — filename in quoted script is injectable",
		Severity: SeverityWarning,
		Description: "Substituting `{}` directly into the quoted command string of `find -exec " +
			"sh -c` lets filenames with shell metacharacters break out. A file named `$(rm " +
			"-rf ~)` invokes command substitution; a file named `foo; curl evil` chains a " +
			"second command. Pass `{}` as a positional argument to `sh` so the filename " +
			"arrives as a parameter, not as source: `find -exec sh -c 'grep pat \"$1\"' _ {} " +
			"\\;`.",
		Check: checkZC1608,
	})
}

func checkZC1608(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "find" {
		return nil
	}

	hasExec := false
	hasShellC := false
	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-exec" || v == "-execdir" {
			hasExec = true
			continue
		}
		if !hasExec {
			continue
		}
		if (v == "sh" || v == "bash" || v == "zsh" || v == "/bin/sh" || v == "/bin/bash") &&
			i+1 < len(cmd.Arguments) && cmd.Arguments[i+1].String() == "-c" {
			hasShellC = true
			continue
		}
		if !hasShellC {
			continue
		}
		if (strings.HasPrefix(v, "'") || strings.HasPrefix(v, "\"")) &&
			strings.Contains(v, "{}") {
			return []Violation{{
				KataID: "ZC1608",
				Message: "`find -exec sh -c '... {} ...'` interpolates filenames into the " +
					"shell script — metacharacters break out. Pass `{}` as a positional " +
					"arg: `sh -c '... \"$1\"' _ {} \\;`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
