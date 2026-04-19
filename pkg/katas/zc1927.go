package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1927",
		Title:    "Error on `xfreerdp /p:SECRET` / `rdesktop -p SECRET` — RDP password visible in argv",
		Severity: SeverityError,
		Description: "`xfreerdp /p:<password>` and `rdesktop -p <password>` (plus the `-p -` " +
			"stdin form when followed by an argv password) put the Windows credential into " +
			"`ps`, `/proc/PID/cmdline`, shell history, and every `ps aux` captured by " +
			"monitoring. Use `xfreerdp /from-stdin` + a piped credential, " +
			"`freerdp-shadow-cli /sec:nla` with a cached credential, or drop the password " +
			"into a protected `.rdp` file passed via `/load-config-file`. Never inline the " +
			"password on the command line.",
		Check: checkZC1927,
	})
}

func checkZC1927(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "xfreerdp", "xfreerdp3", "wlfreerdp":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if strings.HasPrefix(v, "/p:") && len(v) > 3 && v != "/p:-" {
				return zc1927Hit(cmd, ident.Value+" "+v)
			}
		}
	case "rdesktop":
		for i, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-p" && i+1 < len(cmd.Arguments) {
				next := cmd.Arguments[i+1].String()
				if next != "-" {
					return zc1927Hit(cmd, "rdesktop -p "+next)
				}
			}
		}
	}
	return nil
}

func zc1927Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1927",
		Message: "`" + form + "` puts the RDP password in argv — visible in `ps`, " +
			"`/proc`, and shell history. Pipe via `/from-stdin`, read from a protected " +
			"`.rdp` file, or use NLA with a cached credential.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
