package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1523",
		Title:    "Error on `tar -C /` — extracting an archive into the filesystem root",
		Severity: SeverityError,
		Description: "Extracting a tarball directly into `/` overwrites any file it carries a " +
			"matching path for. Combined with a malicious tarball that contains entries like " +
			"`etc/pam.d/sshd` or `usr/bin/ls`, this is a full system compromise disguised as a " +
			"software install. Always extract into a staging directory, inspect contents, then " +
			"copy specific files into place.",
		Check: checkZC1523,
	})
}

func checkZC1523(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tar" && ident.Value != "bsdtar" {
		return nil
	}

	var prevC bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevC {
			prevC = false
			if v == "/" {
				return []Violation{{
					KataID: "ZC1523",
					Message: "`tar -C /` extracts into the filesystem root — overwrites any " +
						"path that happens to be inside the archive. Stage, inspect, then copy.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
		if v == "-C" || v == "--directory" {
			prevC = true
		}
	}
	return nil
}
