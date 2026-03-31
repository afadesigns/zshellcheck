package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1201",
		Title:    "Avoid `rsh`/`rlogin`/`rcp` — use `ssh`/`scp`",
		Severity: SeverityWarning,
		Description: "`rsh`, `rlogin`, and `rcp` are insecure legacy protocols. " +
			"Use `ssh`, `scp`, or `rsync` over SSH for encrypted remote operations.",
		Check: checkZC1201,
	})
}

func checkZC1201(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value
	if name != "rsh" && name != "rlogin" && name != "rcp" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1201",
		Message: "Avoid `" + name + "` — it is an insecure legacy protocol. " +
			"Use `ssh`/`scp`/`rsync` for encrypted remote operations.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
