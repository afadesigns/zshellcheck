package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1825",
		Title:    "Warn on `scp -O` — forces legacy SCP wire protocol exposed to filename-injection CVEs",
		Severity: SeverityWarning,
		Description: "OpenSSH 9.0 switched `scp` to use the SFTP protocol by default — SFTP performs " +
			"structured file transfer instead of piping a remote shell, and closes the " +
			"filename-injection class that the old SCP wire protocol was vulnerable to " +
			"(CVE-2020-15778 and friends). `scp -O` forces the legacy SCP protocol, putting " +
			"the connection back on the old code path where a server (or a man-in-the-middle " +
			"in the remote host's shell) can inject shell metacharacters into filenames. If a " +
			"remote endpoint genuinely needs SCP, use `sftp` instead or upgrade the remote " +
			"server. Drop `-O` unless you have a named compatibility bug that requires it.",
		Check: checkZC1825,
	})
}

func checkZC1825(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "scp" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if arg.String() == "-O" {
			return []Violation{{
				KataID: "ZC1825",
				Message: "`scp -O` forces the legacy SCP wire protocol — the one exposed " +
					"to filename-injection (CVE-2020-15778 class). Drop `-O` (default " +
					"SFTP is safer), or use `sftp` / upgrade the remote server.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
