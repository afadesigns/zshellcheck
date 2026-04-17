package katas

import (
	"strconv"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1578",
		Title:    "Warn on `ssh-keygen -b <2048` for RSA / DSA — weak SSH key",
		Severity: SeverityWarning,
		Description: "Generating an SSH RSA or DSA key shorter than 2048 bits fails current " +
			"OpenSSH baselines and is rejected by recent `ssh` versions when used for " +
			"authentication. DSA was removed from OpenSSH 9.8 outright. Use `ssh-keygen -t " +
			"ed25519` (compact, fast, modern defaults) or `ssh-keygen -t rsa -b 4096` if you " +
			"need RSA for compatibility.",
		Check: checkZC1578,
	})
}

func checkZC1578(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh-keygen" {
		return nil
	}

	var keyType string
	var bits int
	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-t" && i+1 < len(cmd.Arguments) {
			keyType = cmd.Arguments[i+1].String()
		}
		if v == "-b" && i+1 < len(cmd.Arguments) {
			n, err := strconv.Atoi(cmd.Arguments[i+1].String())
			if err == nil {
				bits = n
			}
		}
	}

	// DSA regardless of size is weak / removed.
	if keyType == "dsa" {
		return []Violation{{
			KataID:  "ZC1578",
			Message: "`ssh-keygen -t dsa` — DSA removed from OpenSSH 9.8. Use `-t ed25519`.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityWarning,
		}}
	}
	// RSA below 2048 bits is weak.
	if (keyType == "rsa" || keyType == "") && bits > 0 && bits < 2048 {
		return []Violation{{
			KataID: "ZC1578",
			Message: "`ssh-keygen -b " + strconv.Itoa(bits) + "` — RSA below 2048 bits is " +
				"rejected by modern OpenSSH. Use `-t ed25519` or `-b 4096`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}
