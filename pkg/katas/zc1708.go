package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1708",
		Title:    "Error on `find -L ... -delete` / `-exec rm` — symlink follow into unintended trees",
		Severity: SeverityError,
		Description: "`find -L` follows symlinks during traversal. Combined with `-delete` (or " +
			"`-exec rm`), a symlink under the start path that points outside the intended " +
			"root steers `find` into / `unlink`s files in `/etc`, `/var/lib`, or any other " +
			"directory the symlink target reaches. Drop `-L` (the default `-P` keeps " +
			"symlinks as objects), or restrict the walk with `-xdev`, `-mount`, and an " +
			"explicit `-type f` test. For log-rotation pipes, `logrotate` is safer than a " +
			"`find` one-liner.",
		Check: checkZC1708,
	})
}

func checkZC1708(node ast.Node) []Violation {
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

	hasFollow := false
	hasDestructive := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-L", "--follow":
			hasFollow = true
		case "-delete", "-exec":
			hasDestructive = true
		}
	}

	if !hasFollow || !hasDestructive {
		return nil
	}

	return []Violation{{
		KataID: "ZC1708",
		Message: "`find -L … -delete/-exec` follows symlinks into unintended trees — drop " +
			"`-L`, add `-xdev`, or scope the walk explicitly.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
