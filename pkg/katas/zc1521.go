package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1521",
		Title:    "Style: `strace` without `-e` filter — captures every syscall (incl. secrets, huge output)",
		Severity: SeverityStyle,
		Description: "Unfiltered `strace` records every syscall the process makes: every " +
			"`read()`/`write()` buffer, every `connect()` sockaddr, every `open()` path. That " +
			"includes passwords read from stdin, session tokens written to TLS sockets, and " +
			"any memory a `write()` buffer happens to point at. Scope with `-e trace=<set>` " +
			"(e.g. `trace=openat,connect`) and strip sensitive content with `-e abbrev=all`.",
		Check: checkZC1521,
	})
}

func checkZC1521(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "strace" && ident.Value != "ltrace" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	// Any filter flag present → skip.
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-e" || v == "--trace" || v == "--trace-path" || v == "-P" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1521",
		Message: "`" + ident.Value + "` without `-e` captures every syscall including secrets " +
			"in read/write buffers. Scope with `-e trace=<set>`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
