package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1567",
		Title:    "Warn on `python -m http.server` without `--bind 127.0.0.1` — serves to all interfaces",
		Severity: SeverityWarning,
		Description: "`python -m http.server` (and the legacy `SimpleHTTPServer`) default to " +
			"`0.0.0.0`, exposing the current directory's contents to every network the host " +
			"is on. Tmp scratch files, `.env`, SSH keys, or a `node_modules` tree with private " +
			"config all become reachable from anywhere on the LAN (or the internet, on a VPS). " +
			"Pass `--bind 127.0.0.1` (or `--bind ::1`) unless you really need external access " +
			"and know what is in the cwd.",
		Check: checkZC1567,
	})
}

func checkZC1567(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "python" && ident.Value != "python2" && ident.Value != "python3" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// python -m http.server  or  python -m SimpleHTTPServer
	var isServer bool
	for i := 0; i+1 < len(args); i++ {
		if args[i] == "-m" && (args[i+1] == "http.server" || args[i+1] == "SimpleHTTPServer") {
			isServer = true
			break
		}
	}
	if !isServer {
		return nil
	}

	for _, a := range args {
		if a == "--bind" || a == "-b" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1567",
		Message: "`python -m http.server` without `--bind` defaults to 0.0.0.0 — exposes the " +
			"cwd to every network the host sees. Add `--bind 127.0.0.1`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
