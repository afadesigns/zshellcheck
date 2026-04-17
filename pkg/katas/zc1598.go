package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1598",
		Title:    "Error on `chmod` with world-write bit on a sensitive `/dev/` node",
		Severity: SeverityError,
		Description: "Device nodes under `/dev/` are kernel interfaces. Making one world-writable " +
			"( last digit `2`, `3`, `6`, or `7` ) gives every local user a direct line into the " +
			"kernel — `/dev/kvm` yields VM hypercalls, `/dev/mem` / `/dev/kmem` / `/dev/port` " +
			"read and write physical memory, `/dev/sd*` and `/dev/nvme*` give raw block access, " +
			"`/dev/input/*` sniffs keystrokes. Keep restrictive perms (600 / 660) and use udev " +
			"rules (`GROUP=`, `MODE=`) to grant access declaratively.",
		Check: checkZC1598,
	})
}

func checkZC1598(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chmod" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	mode := cmd.Arguments[0].String()
	if mode == "" {
		return nil
	}
	last := mode[len(mode)-1]
	if last != '2' && last != '3' && last != '6' && last != '7' {
		return nil
	}

	safe := map[string]bool{
		"/dev/null": true, "/dev/zero": true,
		"/dev/random": true, "/dev/urandom": true, "/dev/full": true,
		"/dev/stdin": true, "/dev/stdout": true, "/dev/stderr": true,
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if !strings.HasPrefix(v, "/dev/") {
			continue
		}
		if safe[v] || strings.HasPrefix(v, "/dev/tty") || strings.HasPrefix(v, "/dev/pts/") {
			continue
		}
		return []Violation{{
			KataID: "ZC1598",
			Message: "`chmod " + mode + " " + v + "` makes a sensitive device node world-" +
				"writable — direct kernel access for every local user. Keep restrictive " +
				"perms (600 / 660) and grant access via udev rules.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}
