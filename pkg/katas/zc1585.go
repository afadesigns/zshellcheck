package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1585",
		Title:    "Warn on `ufw allow from any` / `firewall-cmd --add-source=0.0.0.0/0`",
		Severity: SeverityWarning,
		Description: "`ufw allow from any to any port …` (and its firewall-cmd sibling " +
			"`--add-source=0.0.0.0/0`) opens the port to the whole internet. That is " +
			"sometimes the point (public HTTP / HTTPS), but on management ports (22, 3306, " +
			"5432, 6379, 9200, 27017) it is a routine foot-gun when the script author " +
			"assumed the host would only ever be reached via VPN. Scope the rule to a " +
			"specific source CIDR.",
		Check: checkZC1585,
	})
}

func checkZC1585(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ufw" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	if len(args) < 3 || args[0] != "allow" {
		return nil
	}
	// ufw allow from any ...
	for i := 1; i+1 < len(args); i++ {
		if args[i] == "from" && args[i+1] == "any" {
			return []Violation{{
				KataID: "ZC1585",
				Message: "`ufw allow from any …` opens the port to the whole internet. " +
					"Scope to a specific source CIDR.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
