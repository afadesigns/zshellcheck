package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1279",
		Title:    "Use `realpath` instead of `readlink -f` for canonical paths",
		Severity: SeverityInfo,
		Description: "`readlink -f` is not portable across all platforms (notably macOS). " +
			"Use `realpath` which is POSIX-standard and available on modern systems.",
		Check: checkZC1279,
		Fix:   fixZC1279,
	})
}

// fixZC1279 collapses `readlink -f` to `realpath` when `-f` is the
// first argument. Single span replacement from the start of
// `readlink` through the end of `-f`. Only fires when `-f` is the
// first argument; other shapes (`readlink -n -f`, `readlink path -f`)
// are left alone to avoid clobbering unrelated flags. Idempotent —
// a re-run sees `realpath`, not `readlink`. Defensive byte-match
// guards on both anchors.
func fixZC1279(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "readlink" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "-f" {
		return nil
	}
	cmdOff := LineColToByteOffset(source, v.Line, v.Column)
	if cmdOff < 0 || cmdOff+len("readlink") > len(source) {
		return nil
	}
	if string(source[cmdOff:cmdOff+len("readlink")]) != "readlink" {
		return nil
	}
	fTok := cmd.Arguments[0].TokenLiteralNode()
	fOff := LineColToByteOffset(source, fTok.Line, fTok.Column)
	if fOff < 0 || fOff+len("-f") > len(source) {
		return nil
	}
	if string(source[fOff:fOff+len("-f")]) != "-f" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  fOff + len("-f") - cmdOff,
		Replace: "realpath",
	}}
}

func checkZC1279(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "readlink" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-f" {
			return []Violation{{
				KataID:  "ZC1279",
				Message: "Use `realpath` instead of `readlink -f`. `realpath` is more portable, especially on macOS.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityInfo,
			}}
		}
	}

	return nil
}
