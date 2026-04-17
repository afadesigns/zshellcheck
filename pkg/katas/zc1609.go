package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1609",
		Title:    "Warn on `aa-disable` / `aa-complain` / `apparmor_parser -R` — disables AppArmor enforcement",
		Severity: SeverityWarning,
		Description: "`aa-disable` fully unloads the named AppArmor profile; `aa-complain` " +
			"flips the profile from enforce to complain (violations are logged but allowed); " +
			"`apparmor_parser -R` removes a profile from the running kernel. Each one lets the " +
			"confined process run without its mandatory-access-control restrictions — if the " +
			"profile existed for a reason, that reason is now unenforced. Interactive debugging " +
			"is legitimate, but scripts that permanently disable profiles should be reviewed.",
		Check: checkZC1609,
	})
}

func checkZC1609(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "aa-disable", "aa-complain":
		if len(cmd.Arguments) == 0 {
			return nil
		}
		return []Violation{{
			KataID: "ZC1609",
			Message: "`" + ident.Value + "` disables or softens the AppArmor profile — the " +
				"confined process loses MAC restrictions. Review the profile's intent " +
				"before disabling in automation.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	case "apparmor_parser":
		for _, arg := range cmd.Arguments {
			if arg.String() == "-R" || arg.String() == "--remove" {
				return []Violation{{
					KataID: "ZC1609",
					Message: "`apparmor_parser -R` removes the AppArmor profile from the " +
						"kernel — the confined process loses MAC restrictions. Review " +
						"the profile's intent before removing in automation.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}
