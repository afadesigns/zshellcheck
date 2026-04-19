package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1899",
		Title:    "Error on `mokutil --disable-validation` — turns UEFI Secure Boot off at the shim",
		Severity: SeverityError,
		Description: "`mokutil --disable-validation` queues a request for the shim to stop " +
			"validating the kernel and modules against the enrolled MOK/PK certificates at " +
			"next boot — Secure Boot silently becomes advisory. Any unsigned kernel or " +
			"rootkit module then loads without prompt. Leave Secure Boot validation on; " +
			"if you must load a custom module, enrol its key with `mokutil --import` and " +
			"approve via the `MokManager` prompt at reboot.",
		Check: checkZC1899,
	})
}

func checkZC1899(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `mokutil --disable-validation` mangles the command
	// name to `disable-validation`.
	switch ident.Value {
	case "disable-validation":
		return zc1899Hit(cmd)
	case "mokutil":
		for _, arg := range cmd.Arguments {
			if arg.String() == "--disable-validation" {
				return zc1899Hit(cmd)
			}
		}
	}
	return nil
}

func zc1899Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1899",
		Message: "`mokutil --disable-validation` stops the shim from validating " +
			"kernel/modules against enrolled keys — Secure Boot becomes advisory. " +
			"Leave validation on; enrol specific keys with `mokutil --import`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
