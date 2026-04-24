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

var zc1767BindFlags = map[string]bool{"--bind_ip": true}

func checkZC1767(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "mongod" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if !zc1767BindFlags[arg.String()] {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			return nil
		}
		ip := cmd.Arguments[i+1].String()
		if ip != "0.0.0.0" && ip != "::" && ip != "[::]" {
			return nil
		}
		line, col := FlagArgPosition(cmd, zc1767BindFlags)
		return []Violation{{
			KataID: "ZC1767",
			Message: "`mongod --bind_ip " + ip + "` exposes MongoDB on every interface — " +
				"2017 ransomware-wave target. Bind to `127.0.0.1` or a private-network IP, " +
				"enable `--auth`, firewall port 27017.",
			Line:   line,
			Column: col,
			Level:  SeverityError,
		}}
	}
	return nil
}
