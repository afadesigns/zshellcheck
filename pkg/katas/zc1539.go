package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1539",
		Title:    "Warn on `parted -s <disk> <destructive-op>` — script mode bypasses confirmation",
		Severity: SeverityWarning,
		Description: "`parted -s` (script mode) answers the `data will be destroyed` prompt " +
			"with `yes`. Combined with `mklabel`, `mkpart`, `rm`, or `resizepart` on the " +
			"wrong device variable it silently repartitions or zeros the partition table on a " +
			"disk the author never intended. Require an explicit `parted <disk> print` check " +
			"plus an out-of-band confirmation before the destructive call.",
		Check: checkZC1539,
	})
}

func checkZC1539(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "parted" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	var hasScript bool
	for _, a := range args {
		if a == "-s" {
			hasScript = true
		}
	}
	if !hasScript {
		return nil
	}
	for _, a := range args {
		switch a {
		case "mklabel", "mkpart", "rm", "resizepart", "mkpartfs":
			return []Violation{{
				KataID: "ZC1539",
				Message: "`parted -s <disk> " + a + "` bypasses the confirmation prompt — a " +
					"typo in the disk variable silently repartitions the wrong device.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
