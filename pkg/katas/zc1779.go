package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1779AdminRoles = map[string]bool{
	"owner":                     true,
	"contributor":               true,
	"user access administrator": true,
	"useraccessadministrator":   true,
	"role based access control administrator": true,
	"security admin":       true,
	"global administrator": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1779",
		Title:    "Error on `az role assignment create --role Owner|Contributor|User Access Administrator`",
		Severity: SeverityError,
		Description: "`az role assignment create --role Owner` grants full control over the " +
			"target scope (subscription, resource group, resource). `Contributor` grants " +
			"everything except role assignment, and `User Access Administrator` grants the " +
			"ability to assign any role — including Owner — elsewhere in the directory. Any " +
			"of the three is effectively top-of-chain in the assigned scope. In provisioning " +
			"automation this breaks least privilege, invites blast-radius escalations, and " +
			"sidesteps any review that would flag the permission grant. Assign a narrower " +
			"built-in role (Reader, specific-service Contributor) or a custom role whose " +
			"permission list you can enumerate.",
		Check: checkZC1779,
	})
}

func checkZC1779(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "az" {
		return nil
	}
	if len(cmd.Arguments) < 4 {
		return nil
	}
	if cmd.Arguments[0].String() != "role" ||
		cmd.Arguments[1].String() != "assignment" ||
		cmd.Arguments[2].String() != "create" {
		return nil
	}

	for i, arg := range cmd.Arguments[3:] {
		v := arg.String()
		if v == "--role" {
			if i+4 < len(cmd.Arguments) {
				role := cmd.Arguments[3+i+1].String()
				role = strings.Trim(role, "\"'")
				if zc1779IsAdminRole(role) {
					return zc1779Hit(cmd, role)
				}
			}
			continue
		}
		if strings.HasPrefix(v, "--role=") {
			role := strings.TrimPrefix(v, "--role=")
			role = strings.Trim(role, "\"'")
			if zc1779IsAdminRole(role) {
				return zc1779Hit(cmd, role)
			}
		}
	}
	return nil
}

func zc1779IsAdminRole(role string) bool {
	r := strings.ToLower(strings.TrimSpace(role))
	r = strings.ReplaceAll(r, "_", " ")
	r = strings.ReplaceAll(r, "-", " ")
	r = strings.Join(strings.Fields(r), " ")
	return zc1779AdminRoles[r]
}

func zc1779Hit(cmd *ast.SimpleCommand, role string) []Violation {
	return []Violation{{
		KataID: "ZC1779",
		Message: "`az role assignment create --role " + role + "` grants a top-of-chain " +
			"role. Pick a narrower built-in role (`Reader`, specific-service Contributor) " +
			"or a custom role whose permission list you can enumerate.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
