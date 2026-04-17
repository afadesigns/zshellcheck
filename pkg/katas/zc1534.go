package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1534",
		Title:    "Warn on `dmesg -c` / `--clear` — wipes kernel ring buffer",
		Severity: SeverityWarning,
		Description: "`dmesg -c` reads and then clears the kernel ring buffer. Any subsequent " +
			"reader sees an empty log, so OOM kills, driver panics, and audit messages that " +
			"landed between the wipe and the incident response are gone. It is also an " +
			"anti-forensics step in post-exploitation playbooks. Use `dmesg` (no flags) for a " +
			"read, and let the journal retention policy handle rotation.",
		Check: checkZC1534,
	})
}

func checkZC1534(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dmesg" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c" || v == "-C" || v == "--clear" || v == "--read-clear" {
			return []Violation{{
				KataID: "ZC1534",
				Message: "`dmesg " + v + "` wipes the kernel ring buffer — subsequent " +
					"readers see no OOM/panic/audit messages. Read without clearing.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
