package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1562",
		Title:    "Warn on `env -u PATH` / `-u LD_LIBRARY_PATH` — clears security-relevant env",
		Severity: SeverityWarning,
		Description: "`env -u PATH` unsets the caller's `PATH` before running the child, forcing " +
			"the child to fall back to the hard-coded search list (`/bin:/usr/bin` on glibc). " +
			"That bypasses PATH hardening done by the parent shell (e.g. a sanitised PATH " +
			"under `sudo`). Unsetting `LD_PRELOAD` / `LD_LIBRARY_PATH` mid-stream is also " +
			"usually the caller trying to shake off an earlier `export`. Either use `env -i` " +
			"to sanitise completely, or explicitly set the variables the child should see.",
		Check: checkZC1562,
	})
}

func checkZC1562(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "env" {
		return nil
	}

	var prevU bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevU {
			prevU = false
			switch v {
			case "PATH", "LD_PRELOAD", "LD_LIBRARY_PATH", "LD_AUDIT":
				return []Violation{{
					KataID: "ZC1562",
					Message: "`env -u " + v + "` clears a security-relevant variable mid-run. " +
						"Use `env -i` to sanitise, or set the right value explicitly.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
		if v == "-u" || v == "--unset" {
			prevU = true
		}
	}
	return nil
}
