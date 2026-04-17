package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1459",
		Title:    "Warn on `docker run --cap-add=SYS_ADMIN` / other dangerous capabilities",
		Severity: SeverityWarning,
		Description: "Granting `SYS_ADMIN`, `SYS_PTRACE`, `SYS_MODULE`, `NET_ADMIN`, or `ALL` " +
			"capabilities effectively disables the container's security boundary — most " +
			"container escapes rely on exactly these. Drop all capabilities and add back only " +
			"the specific ones the workload needs (usually none).",
		Check: checkZC1459,
	})
}

var dangerousCaps = map[string]struct{}{
	"SYS_ADMIN":       {},
	"SYS_PTRACE":      {},
	"SYS_MODULE":      {},
	"SYS_RAWIO":       {},
	"NET_ADMIN":       {},
	"DAC_READ_SEARCH": {},
	"ALL":             {},
}

func checkZC1459(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}

	var prevCap bool
	for _, arg := range cmd.Arguments {
		v := arg.String()

		// Split form: --cap-add=VALUE
		if strings.HasPrefix(v, "--cap-add=") {
			val := strings.TrimPrefix(v, "--cap-add=")
			if _, bad := dangerousCaps[strings.ToUpper(val)]; bad {
				return violateZC1459(cmd)
			}
			continue
		}

		// Space form: --cap-add VALUE
		if prevCap {
			prevCap = false
			if _, bad := dangerousCaps[strings.ToUpper(v)]; bad {
				return violateZC1459(cmd)
			}
		}
		if v == "--cap-add" {
			prevCap = true
		}
	}

	return nil
}

func violateZC1459(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1459",
		Message: "Dangerous Linux capability granted — breaks the container's security " +
			"boundary. Prefer `--cap-drop=ALL` and add back only minimum needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
