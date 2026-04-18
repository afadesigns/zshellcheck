package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1682",
		Title:    "Error on `npm install --unsafe-perm` — npm lifecycle scripts keep root privileges",
		Severity: SeverityError,
		Description: "npm normally drops to the UID that owns `package.json` before running " +
			"`preinstall` / `install` / `postinstall` lifecycle scripts. `--unsafe-perm` " +
			"(or `--unsafe-perm=true`) tells npm to skip that drop and run every script as " +
			"the current UID — typically root when the install happens from a provisioning " +
			"script. Any compromised or malicious dependency then executes as root. If a " +
			"native addon truly needs privileges, scope them: drop them into a dedicated " +
			"builder container, or use `sudo -u builduser npm install` from a non-root " +
			"account that already owns `node_modules/`.",
		Check: checkZC1682,
	})
}

func checkZC1682(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "npm" && ident.Value != "yarn" && ident.Value != "pnpm" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--unsafe-perm" || v == "--unsafe-perm=true" {
			return []Violation{{
				KataID: "ZC1682",
				Message: "`" + ident.Value + " " + v + "` keeps root for every lifecycle " +
					"script — a compromised dep executes as root. Build in a dedicated " +
					"builder container or run as a non-root user.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
