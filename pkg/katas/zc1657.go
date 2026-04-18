package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1657",
		Title:    "Warn on `semanage permissive -a <type>` — puts SELinux domain in permissive mode",
		Severity: SeverityWarning,
		Description: "`semanage permissive -a DOMAIN` (or `--add`) marks an SELinux domain as " +
			"permissive: policy violations are logged but not blocked. It is narrower than " +
			"`setenforce 0` but still disables enforcement for whatever DOMAIN covers — often " +
			"`httpd_t`, `container_t`, or `sshd_t` — and the override persists across reboots " +
			"because it is written to policy. Fix the denial with an explicit allow rule built " +
			"from `audit2allow` or ship a custom policy module, and remove the permissive mark " +
			"with `semanage permissive -d DOMAIN` once the rule lands.",
		Check: checkZC1657,
	})
}

func checkZC1657(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "semanage" {
		return nil
	}

	hasPermissive := false
	hasAdd := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "permissive":
			hasPermissive = true
		case "-a", "--add":
			hasAdd = true
		}
	}

	if !hasPermissive || !hasAdd {
		return nil
	}

	return []Violation{{
		KataID: "ZC1657",
		Message: "`semanage permissive -a` puts an SELinux domain in permissive mode — " +
			"policy violations log but no longer block. Write a scoped allow rule with " +
			"`audit2allow` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
