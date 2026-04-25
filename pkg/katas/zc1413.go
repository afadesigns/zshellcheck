package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1413",
		Title:    "Use Zsh `whence -p cmd` instead of `hash -t cmd` for resolved path",
		Severity: SeverityStyle,
		Description: "Bash's `hash -t cmd` prints the hashed path for `cmd` (or fails if not " +
			"hashed). Zsh's `whence -p cmd` prints the PATH-resolved absolute path, whether " +
			"hashed or not — more reliable and the native Zsh idiom.",
		Check: checkZC1413,
		Fix:   fixZC1413,
	})
}

// fixZC1413 rewrites `hash -t cmd` to `whence -p cmd`. Two edits per
// fire: the command name and the `-t` flag. Idempotent — a re-run
// sees `whence`, not `hash`, so the detector won't fire. Defensive
// byte-match guards on both edits refuse to insert if the source
// at the offset doesn't match.
func fixZC1413(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "hash" {
		return nil
	}
	var dashT ast.Node
	for _, arg := range cmd.Arguments {
		if arg.String() == "-t" {
			dashT = arg
			break
		}
	}
	if dashT == nil {
		return nil
	}
	cmdOff := LineColToByteOffset(source, v.Line, v.Column)
	if cmdOff < 0 || cmdOff+len("hash") > len(source) {
		return nil
	}
	if string(source[cmdOff:cmdOff+len("hash")]) != "hash" {
		return nil
	}
	tTok := dashT.TokenLiteralNode()
	tOff := LineColToByteOffset(source, tTok.Line, tTok.Column)
	if tOff < 0 || tOff+len("-t") > len(source) {
		return nil
	}
	if string(source[tOff:tOff+len("-t")]) != "-t" {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: len("hash"), Replace: "whence"},
		{Line: tTok.Line, Column: tTok.Column, Length: len("-t"), Replace: "-p"},
	}
}

func checkZC1413(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "hash" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-t" {
			return []Violation{{
				KataID: "ZC1413",
				Message: "Use `whence -p cmd` (Zsh) instead of `hash -t cmd`. " +
					"`whence -p` always returns the absolute path, regardless of hash state.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
