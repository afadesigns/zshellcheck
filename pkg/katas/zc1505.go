package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1505",
		Title:    "Warn on `dpkg --force-confnew` / `--force-confold` — silently overrides /etc changes",
		Severity: SeverityWarning,
		Description: "`--force-confnew` replaces any locally-modified config file with the " +
			"maintainer version; `--force-confold` keeps the local file and drops the new " +
			"defaults on the floor. Either way dpkg silently picks a side without prompting, " +
			"so a legitimate /etc tweak (hardening, compliance override) can vanish or a " +
			"security-relevant config update can be ignored. Review the conffile diff per " +
			"upgrade (`ucf` / `etckeeper`) rather than hard-coding the decision.",
		Check: checkZC1505,
	})
}

func checkZC1505(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dpkg" && ident.Value != "apt" && ident.Value != "apt-get" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "--force-conf") ||
			strings.HasPrefix(v, "-oDpkg::Options::=--force-conf") ||
			strings.HasPrefix(v, "-o=Dpkg::Options::=--force-conf") {
			return []Violation{{
				KataID: "ZC1505",
				Message: "`" + v + "` silently picks maintainer or local conffile — legit /etc " +
					"changes disappear or new defaults are ignored. Use ucf/etckeeper.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
