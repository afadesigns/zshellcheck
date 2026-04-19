package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1818DeleteFlags = []string{
	"--delete",
	"--del",
	"--delete-before",
	"--delete-during",
	"--delete-delay",
	"--delete-after",
	"--delete-excluded",
	"--delete-missing-args",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1818",
		Title:    "Warn on `rsync --delete` without `--dry-run` — empty or wrong SRC wipes DST",
		Severity: SeverityWarning,
		Description: "`rsync --delete` (plus `--delete-before/-during/-after/-excluded`) removes " +
			"anything in DST that is not in SRC. If SRC is accidentally empty (typo in " +
			"path, unmounted mount point, wrong credentials pointing at an empty remote), " +
			"the destination loses every file that was there. The command has no undo. " +
			"Always preview the diff with `rsync -av --delete --dry-run SRC DST` first, " +
			"and cap the blast radius with `--max-delete=N` so the sync aborts if the plan " +
			"removes more files than expected.",
		Check: checkZC1818,
	})
}

func checkZC1818(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `rsync --delete SRC DST` mangles name to `delete` or similar.
	mangled := false
	switch ident.Value {
	case "delete", "del", "delete-before", "delete-during", "delete-delay",
		"delete-after", "delete-excluded", "delete-missing-args":
		mangled = true
	}

	if !mangled {
		if ident.Value != "rsync" {
			return nil
		}
		hasDelete := false
		for _, arg := range cmd.Arguments {
			v := arg.String()
			for _, flag := range zc1818DeleteFlags {
				if v == flag || strings.HasPrefix(v, flag+"=") {
					hasDelete = true
					break
				}
			}
			if hasDelete {
				break
			}
		}
		if !hasDelete {
			return nil
		}
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--dry-run" || v == "-n" || v == "--itemize-changes" {
			return nil
		}
		if strings.HasPrefix(v, "-") && !strings.HasPrefix(v, "--") && len(v) > 1 {
			for _, c := range v[1:] {
				if c == 'n' {
					return nil
				}
			}
		}
	}
	return []Violation{{
		KataID: "ZC1818",
		Message: "`rsync --delete` without `--dry-run` removes anything in DST that " +
			"isn't in SRC. Preview with `rsync -av --delete --dry-run SRC DST`, and " +
			"pin `--max-delete=N` so an accidentally empty SRC can't cascade.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
