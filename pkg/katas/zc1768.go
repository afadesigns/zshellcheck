package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1768",
		Title:    "Error on `sqlcmd -P PASSWORD` / `bcp -P PASSWORD` — SQL Server password in argv",
		Severity: SeverityError,
		Description: "Microsoft's SQL Server CLI tools (`sqlcmd`, `bcp`, `osql`) accept the " +
			"password via `-P PASSWORD` as a positional argument value. The password lands " +
			"in argv — visible in `ps`, `/proc/<pid>/cmdline`, shell history, CI logs, and " +
			"SQL Server's audit trace for the session. Use `-P` with no value (prompts), " +
			"or read the password from the environment variable `SQLCMDPASSWORD` (sourced " +
			"from a secrets file). On modern sqlcmd, `-G` + Azure AD integration avoids the " +
			"password altogether.",
		Check: checkZC1768,
	})
}

func checkZC1768(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "sqlcmd", "bcp", "osql":
	default:
		return nil
	}

	prevP := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevP {
			prevP = false
			if v == "" || v == "-" {
				continue
			}
			// Any other flag-looking value means `-P` ran as bare (prompt).
			if len(v) > 1 && v[0] == '-' {
				continue
			}
			return []Violation{{
				KataID: "ZC1768",
				Message: "`" + ident.Value + " -P " + v + "` puts the SQL Server password " +
					"in argv — visible in `ps`, `/proc`, history, SQL Server audit. " +
					"Use `-P` with no arg (prompt) or `SQLCMDPASSWORD` env var.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
		if v == "-P" {
			prevP = true
		}
	}
	return nil
}
