package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1448",
		Title:    "`apt-get install` / `apt install` without `-y` hangs in non-interactive scripts",
		Severity: SeverityWarning,
		Description: "In provisioning scripts, `apt-get install foo` (no `-y`) waits for " +
			"interactive confirmation and stalls CI/Dockerfiles indefinitely. Always pass `-y` " +
			"(or `--yes`), and for unattended upgrades also set " +
			"`DEBIAN_FRONTEND=noninteractive` in the environment.",
		Check: checkZC1448,
		Fix:   fixZC1448,
	})
}

// fixZC1448 inserts ` -y` after the `apt` command name so install /
// upgrade / dist-upgrade / full-upgrade run without interactive
// confirmation. Only fires for plain `apt` — for `apt-get` the legacy
// ZC1213 fix already handles the rewrite, and emitting a duplicate
// zero-length insert here would yield ` -y -y` after both edits apply.
func fixZC1448(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "apt" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("apt") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1448(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -y",
	}}
}

func offsetLineColZC1448(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1448(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "apt-get" && ident.Value != "apt" {
		return nil
	}

	hasInstall := false
	hasYes := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "install" || v == "upgrade" || v == "dist-upgrade" || v == "full-upgrade" {
			hasInstall = true
		}
		if v == "-y" || v == "--yes" || v == "--assume-yes" {
			hasYes = true
		}
	}
	if hasInstall && !hasYes {
		return []Violation{{
			KataID: "ZC1448",
			Message: "`apt-get install`/`apt install` without `-y` hangs on the interactive " +
				"prompt in scripts. Add `-y` and set DEBIAN_FRONTEND=noninteractive.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
