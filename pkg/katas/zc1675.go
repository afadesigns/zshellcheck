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
func fixZC1675(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "export" {
		return nil
	}
	var flag ast.Expression
	var replace string
	for _, arg := range cmd.Arguments {
		switch arg.String() {
		case "-f":
			flag = arg
			replace = "typeset -fx"
		case "-n":
			flag = arg
			replace = "typeset +x"
		}
		if flag != nil {
			break
		}
	}
	if flag == nil {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("export") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("export")]) != "export" {
		return nil
	}
	flagTok := flag.TokenLiteralNode()
	flagOff := LineColToByteOffset(source, flagTok.Line, flagTok.Column)
	if flagOff < 0 || flagOff+2 > len(source) {
		return nil
	}
	flagLit := string(source[flagOff : flagOff+2])
	if flagLit != "-f" && flagLit != "-n" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  flagOff + 2 - nameOff,
		Replace: replace,
	}}
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
