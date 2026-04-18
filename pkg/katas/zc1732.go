package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1732BroadFilesystems = map[string]bool{
	"--filesystem=host":     true,
	"--filesystem=host:rw":  true,
	"--filesystem=home":     true,
	"--filesystem=home:rw":  true,
	"--filesystem=/":        true,
	"--filesystem=/:rw":     true,
	"--filesystem=host-os":  true,
	"--filesystem=host-etc": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1732",
		Title:    "Warn on `flatpak override --filesystem=host` — removes Flatpak sandbox isolation",
		Severity: SeverityWarning,
		Description: "Flatpak's primary security guarantee is filesystem sandboxing — apps see " +
			"only their own data plus paths the user explicitly grants via portals. " +
			"`flatpak override --filesystem=host` (also `host-os`, `host-etc`, `home`, `/`) " +
			"persistently grants the app unrestricted read/write to the host filesystem at " +
			"every subsequent run. Same risk applies to `flatpak run --filesystem=host`. " +
			"Grant the specific subdirectory the app actually needs (`--filesystem=" +
			"~/Documents:ro`) or rely on Filesystem portals so the user picks paths " +
			"interactively per session.",
		Check: checkZC1732,
	})
}

func checkZC1732(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "flatpak" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	switch cmd.Arguments[0].String() {
	case "override", "run":
	default:
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if zc1732BroadFilesystems[v] {
			return []Violation{{
				KataID: "ZC1732",
				Message: "`flatpak " + cmd.Arguments[0].String() + " " + v + "` removes " +
					"the Flatpak sandbox — the app gets unrestricted host-filesystem " +
					"access. Grant a specific subdirectory (e.g. " +
					"`--filesystem=~/Documents:ro`) or use Filesystem portals.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
