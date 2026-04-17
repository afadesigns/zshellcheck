package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1465",
		Title:    "Warn on `setenforce 0` — disables SELinux enforcement",
		Severity: SeverityWarning,
		Description: "`setenforce 0` switches SELinux to permissive mode, silencing every policy " +
			"decision into an audit log line instead of a deny. It is the textbook post-" +
			"compromise persistence step and also a common \"fix\" that papers over an actual " +
			"policy bug. Address the specific AVC with `audit2allow` instead, and leave " +
			"`setenforce 1` (enforcing) in production.",
		Check: checkZC1465,
	})
}

func checkZC1465(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setenforce" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	v := cmd.Arguments[0].String()
	if v == "0" || v == "Permissive" || v == "permissive" {
		return []Violation{{
			KataID: "ZC1465",
			Message: "`setenforce 0` disables SELinux enforcement host-wide. Fix the AVC with " +
				"`audit2allow` instead and keep enforcing mode on.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}
