package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1621",
		Title:    "Warn on `tmux -S /tmp/SOCKET` — shared-path socket invites session hijack",
		Severity: SeverityWarning,
		Description: "`tmux -S PATH` overrides the default socket location (normally under " +
			"`$XDG_RUNTIME_DIR/tmux-$UID/`, a 0700-mode directory). Paths under `/tmp/` or " +
			"`/var/tmp/` are world-traversable; if the socket is created with loose " +
			"permissions, any local user who can read it can `tmux -S /tmp/PATH attach` and " +
			"see / drive the session — keystrokes, output, arbitrary commands in the attached " +
			"pane. Keep the socket in `$XDG_RUNTIME_DIR` or another 0700-scoped directory.",
		Check: checkZC1621,
	})
}

func checkZC1621(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tmux" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() != "-S" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		path := cmd.Arguments[i+1].String()
		if strings.HasPrefix(path, "/tmp/") || strings.HasPrefix(path, "/var/tmp/") {
			return []Violation{{
				KataID: "ZC1621",
				Message: "`tmux -S " + path + "` places the socket in a world-traversable " +
					"directory — any local user who can read the socket can attach the " +
					"session. Use `$XDG_RUNTIME_DIR` or a 0700-scoped parent dir.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
