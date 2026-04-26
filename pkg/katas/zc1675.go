// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1675",
		Title:    "Avoid Bash-only `export -f` / `export -n` — use Zsh `typeset -fx` / `typeset +x`",
		Severity: SeverityInfo,
		Description: "`export -f FUNC` (export a function to child processes) and `export -n " +
			"VAR` (strip the export flag while keeping the value) are Bash-only. Zsh's " +
			"`export` ignores `-f` entirely and prints usage for `-n`, so scripts that " +
			"depend on either silently break under Zsh. The Zsh equivalents are `typeset " +
			"-fx FUNC` for function export (parameter-passing via `$FUNCTIONS` in a " +
			"subshell) and `typeset +x VAR` to drop the export flag. Functions that must " +
			"cross a subshell are usually better handled by `autoload -Uz` from an `fpath` " +
			"directory than by serialisation.",
		Check: checkZC1675,
		Fix:   fixZC1675,
	})
}

// fixZC1675 collapses `export -f` and `export -n` into the Zsh
// equivalents `typeset -fx` and `typeset +x`. Single edit spans the
// command name + flag together, mirroring fixZC1283's `set -o OPT`
// → `setopt OPT` collapse.
var zc1675FlagReplace = map[string]string{
	"-f": "typeset -fx",
	"-n": "typeset +x",
}

func fixZC1675(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || CommandIdentifier(cmd) != "export" {
		return nil
	}
	flag, replace := zc1675FindFlag(cmd)
	if flag == nil {
		return nil
	}
	nameOff, ok := zc1675ExportOffset(source, v)
	if !ok {
		return nil
	}
	flagOff, ok := zc1675FlagOffset(source, flag)
	if !ok {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  flagOff + 2 - nameOff,
		Replace: replace,
	}}
}

func zc1675FindFlag(cmd *ast.SimpleCommand) (ast.Expression, string) {
	for _, arg := range cmd.Arguments {
		if r, hit := zc1675FlagReplace[arg.String()]; hit {
			return arg, r
		}
	}
	return nil, ""
}

func zc1675ExportOffset(source []byte, v Violation) (int, bool) {
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 || off+len("export") > len(source) {
		return 0, false
	}
	if string(source[off:off+len("export")]) != "export" {
		return 0, false
	}
	return off, true
}

func zc1675FlagOffset(source []byte, flag ast.Expression) (int, bool) {
	tok := flag.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 || off+2 > len(source) {
		return 0, false
	}
	lit := string(source[off : off+2])
	if lit != "-f" && lit != "-n" {
		return 0, false
	}
	return off, true
}

func checkZC1675(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "export" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-f":
			return zc1675Hit(cmd, "export -f", "typeset -fx")
		case "-n":
			return zc1675Hit(cmd, "export -n", "typeset +x")
		}
	}
	return nil
}

func zc1675Hit(cmd *ast.SimpleCommand, bad, good string) []Violation {
	return []Violation{{
		KataID:  "ZC1675",
		Message: "`" + bad + "` is Bash-only — use `" + good + "` in Zsh.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityInfo,
	}}
}
