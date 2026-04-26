// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1252",
		Title:    "Use `getent passwd` instead of `cat /etc/passwd`",
		Severity: SeverityStyle,
		Description: "`cat /etc/passwd` misses users from LDAP, NIS, or SSSD sources. " +
			"`getent passwd` queries NSS and returns all configured user databases.",
		Check: checkZC1252,
		Fix:   fixZC1252,
	})
}

// fixZC1252 rewrites `cat /etc/{passwd,group,shadow}` to
// `getent {passwd,group,shadow}`. Two edits per fire: the command
// name and the file argument. Only fires when the cat command has
// exactly one argument; piped or multi-file shapes are left alone
// (ZC1146 handles `cat FILE | tool`, and multi-file `cat` doesn't
// translate cleanly to `getent`). Idempotent: a re-run sees `getent`,
// not `cat`, so the detector won't fire.
func fixZC1252(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}
	if len(cmd.Arguments) != 1 {
		return nil
	}
	arg := cmd.Arguments[0]
	val := arg.String()
	var dbName string
	switch val {
	case "/etc/passwd":
		dbName = "passwd"
	case "/etc/group":
		dbName = "group"
	case "/etc/shadow":
		dbName = "shadow"
	default:
		return nil
	}
	catOff := LineColToByteOffset(source, v.Line, v.Column)
	if catOff < 0 || catOff+len("cat") > len(source) {
		return nil
	}
	if string(source[catOff:catOff+len("cat")]) != "cat" {
		return nil
	}
	argTok := arg.TokenLiteralNode()
	argOff := LineColToByteOffset(source, argTok.Line, argTok.Column)
	if argOff < 0 || argOff+len(val) > len(source) {
		return nil
	}
	if string(source[argOff:argOff+len(val)]) != val {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: len("cat"), Replace: "getent"},
		{Line: argTok.Line, Column: argTok.Column, Length: len(val), Replace: dbName},
	}
}

func checkZC1252(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/etc/passwd" || val == "/etc/group" || val == "/etc/shadow" {
			return []Violation{{
				KataID: "ZC1252",
				Message: "Use `getent` instead of `cat " + val + "`. " +
					"`getent` queries all NSS sources including LDAP and SSSD.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
