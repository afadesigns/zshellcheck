package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1437",
		Title:    "`dmesg -c` / `-C` clears the kernel ring buffer — destroys evidence",
		Severity: SeverityWarning,
		Description: "`dmesg -c` prints the ring buffer and then **clears** it. `dmesg -C` clears " +
			"without printing. Any later debugging loses the earlier messages. Prefer plain " +
			"`dmesg` for read-only inspection, or `journalctl -k` with a time filter.",
		Check: checkZC1437,
	})
}

func checkZC1437(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "dmesg" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c" || v == "-C" || v == "--clear" || v == "--read-clear" {
			return []Violation{{
				KataID: "ZC1437",
				Message: "`dmesg -c`/`-C` clears the kernel ring buffer — subsequent debugging " +
					"loses earlier messages. Use plain `dmesg` or `journalctl -k`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
