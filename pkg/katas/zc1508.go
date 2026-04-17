package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1508",
		Title:    "Style: `ldd <binary>` may execute the binary — use `objdump -p` / `readelf -d` for untrusted files",
		Severity: SeverityStyle,
		Description: "On glibc, `ldd` is implemented by setting `LD_TRACE_LOADED_OBJECTS=1` and " +
			"invoking the binary. A malicious ELF with a custom interpreter (`PT_INTERP`) or " +
			"constructors can therefore run code when `ldd` is pointed at it. `objdump -p " +
			"<file> | grep NEEDED` or `readelf -d <file>` give the same shared-library list " +
			"without executing the binary.",
		Check: checkZC1508,
	})
}

func checkZC1508(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ldd" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1508",
		Message: "`ldd` on glibc can execute the target binary. Use `objdump -p` or " +
			"`readelf -d` to inspect ELF dependencies safely.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
