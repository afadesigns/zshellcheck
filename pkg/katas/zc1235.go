package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1235",
		Title:    "Use `git push --force-with-lease` instead of `--force`",
		Severity: SeverityWarning,
		Description: "`git push --force` overwrites remote history unconditionally. " +
			"`--force-with-lease` is safer as it fails if the remote has changed.",
		Check: checkZC1235,
		Fix:   fixZC1235,
	})
}

// fixZC1235 rewrites `git push -f` to `git push --force-with-lease`.
// Single-edit replacement of the `-f` flag at its argument position;
// surrounding subcommand and refspec arguments stay in place.
func fixZC1235(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "push" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		if arg.String() != "-f" {
			continue
		}
		tok := arg.TokenLiteralNode()
		off := LineColToByteOffset(source, tok.Line, tok.Column)
		if off < 0 || off+2 > len(source) {
			return nil
		}
		if string(source[off:off+2]) != "-f" {
			return nil
		}
		return []FixEdit{{
			Line:    tok.Line,
			Column:  tok.Column,
			Length:  2,
			Replace: "--force-with-lease",
		}}
	}
	return nil
}

func checkZC1235(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "push" {
		return nil
	}

	hasForce := false
	hasFWL := false

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "-f" {
			hasForce = true
		}
	}

	if hasForce && !hasFWL {
		return []Violation{{
			KataID: "ZC1235",
			Message: "Use `git push --force-with-lease` instead of `-f`/`--force`. " +
				"It prevents overwriting remote changes made by others.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
