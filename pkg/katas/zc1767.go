package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1767",
		Title:    "Error on `mongod --bind_ip 0.0.0.0` — MongoDB exposed on every interface",
		Severity: SeverityError,
		Description: "`mongod --bind_ip 0.0.0.0` (or `::`) binds MongoDB's listener to every " +
			"interface on the host. Combined with no-auth defaults (pre-3.4) or a wildcard " +
			"database user, this was the source of the 2017 ransomware wave that wiped " +
			"tens of thousands of public MongoDB instances. Bind to `127.0.0.1` or a " +
			"private-network IP, enable authentication with `--auth`, and firewall port " +
			"`27017`.",
		Check: checkZC1767,
	})
}

func checkZC1767(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	// Parser caveat: `mongod --bind_ip 0.0.0.0` mangles to SimpleCommand name
	// "bind_ip" with the IP as the first arg.
	if ident.Value != "bind_ip" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	first := cmd.Arguments[0].String()
	if first != "0.0.0.0" && first != "::" && first != "[::]" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1767",
		Message: "`mongod --bind_ip " + first + "` exposes MongoDB on every interface — " +
			"2017 ransomware-wave target. Bind to `127.0.0.1` or a private-network IP, " +
			"enable `--auth`, firewall port 27017.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
