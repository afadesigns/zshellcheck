package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1725",
		Title:    "Error on `cargo --token TOKEN` / `npm --otp CODE` — registry credential in process list",
		Severity: SeverityError,
		Description: "`cargo publish --token TOKEN` (and `cargo login`, `cargo owner`, `cargo " +
			"yank`) puts the crates.io API token in argv — visible in `ps`, `/proc/<pid>/" +
			"cmdline`, shell history, and CI logs. `npm publish --otp CODE` leaks the " +
			"one-time code the same way. Use environment variables (`CARGO_REGISTRY_TOKEN`, " +
			"`NPM_TOKEN`) or pipe via stdin (`cargo login --token -` reads from stdin), and " +
			"source credentials from a secrets manager instead of the command line.",
		Check: checkZC1725,
	})
}

func checkZC1725(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var flag, tool string
	switch ident.Value {
	case "cargo":
		flag = "--token"
		tool = "cargo"
	case "npm", "yarn", "pnpm":
		flag = "--otp"
		tool = ident.Value
	default:
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	// First arg should be a relevant subcommand.
	sub := cmd.Arguments[0].String()
	switch tool {
	case "cargo":
		switch sub {
		case "publish", "login", "owner", "yank":
		default:
			return nil
		}
	default:
		switch sub {
		case "publish", "adduser", "login":
		default:
			return nil
		}
	}

	prevFlag := false
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if prevFlag {
			if v == "-" {
				return nil
			}
			return zc1725Hit(cmd, tool, sub, flag+" "+v)
		}
		switch {
		case v == flag:
			prevFlag = true
		case strings.HasPrefix(v, flag+"="):
			val := strings.TrimPrefix(v, flag+"=")
			if val == "-" {
				return nil
			}
			return zc1725Hit(cmd, tool, sub, v)
		}
	}
	return nil
}

func zc1725Hit(cmd *ast.SimpleCommand, tool, sub, what string) []Violation {
	return []Violation{{
		KataID: "ZC1725",
		Message: "`" + tool + " " + sub + " " + what + "` puts the credential in argv — " +
			"visible in `ps`, `/proc`, history. Pipe via stdin (`--token -`) or use env " +
			"vars like `CARGO_REGISTRY_TOKEN` / `NPM_TOKEN`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
