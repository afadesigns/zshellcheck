package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1606",
		Title:    "Warn on `mkdir -m NNN` / `install -m NNN` with world-write bit (no sticky)",
		Severity: SeverityWarning,
		Description: "`mkdir -m 777 /path` and `install -m 777 src /dest` create a path that " +
			"every local user can write and rename inside. If the script later creates files " +
			"there, classic TOCTOU symlink attacks become trivial — the attacker drops a " +
			"symlink named like the expected output file, redirecting the write wherever they " +
			"choose. A sticky-bit mode (`1777`) mitigates this for shared temp dirs. Prefer " +
			"`mkdir -m 700` (or 750), and scope access by group or ACL rather than everyone.",
		Check: checkZC1606,
	})
}

func checkZC1606(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "mkdir" && ident.Value != "install" {
		return nil
	}

	for i := 0; i+1 < len(cmd.Arguments); i++ {
		if cmd.Arguments[i].String() != "-m" {
			continue
		}
		mode := cmd.Arguments[i+1].String()
		if len(mode) != 3 {
			continue
		}
		allOctal := true
		for _, c := range mode {
			if c < '0' || c > '7' {
				allOctal = false
				break
			}
		}
		if !allOctal {
			continue
		}
		last := mode[2]
		if last != '2' && last != '3' && last != '6' && last != '7' {
			continue
		}
		return []Violation{{
			KataID: "ZC1606",
			Message: "`" + ident.Value + " -m " + mode + "` creates a world-writable path " +
				"without the sticky bit — TOCTOU symlink-attack ground. Use `-m 700` / " +
				"`-m 750`, or `-m 1777` if a shared sticky dir is actually needed.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}
