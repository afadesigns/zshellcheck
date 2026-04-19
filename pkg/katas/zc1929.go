package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1929",
		Title:    "Warn on `cpio -i` / `--extract` without `--no-absolute-filenames` — archive writes outside CWD",
		Severity: SeverityWarning,
		Description: "`cpio -i` (and `--extract`) is the default copy-in mode: it materialises " +
			"every path stored in the archive verbatim. Paths starting with `/` land where the " +
			"archive told them to, and relative paths containing `..` slip out of the " +
			"extraction directory entirely — so a rogue initramfs or firmware bundle can drop " +
			"files into `/etc/cron.d/`, `/usr/lib/systemd/system/`, or the operator's " +
			"`~/.ssh/authorized_keys`. Always pass `--no-absolute-filenames` and extract into a " +
			"fresh scratch directory reviewed before `mv`.",
		Check: checkZC1929,
	})
}

func checkZC1929(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "cpio" {
		return nil
	}

	extract := false
	safe := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-i", "--extract":
			extract = true
		case "--no-absolute-filenames":
			safe = true
		}
		if len(v) >= 2 && v[0] == '-' && v[1] != '-' {
			for i := 1; i < len(v); i++ {
				if v[i] == 'i' {
					extract = true
				}
			}
		}
	}
	if !extract || safe {
		return nil
	}
	return []Violation{{
		KataID: "ZC1929",
		Message: "`cpio -i` extracts paths verbatim — absolute and `..` entries escape the " +
			"target dir. Pass `--no-absolute-filenames` and stage into a scratch dir before " +
			"`mv` into place.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
