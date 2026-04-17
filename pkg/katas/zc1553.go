package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1553",
		Title:    "Style: use Zsh `${(U)var}` / `${(L)var}` instead of `tr '[:lower:]' '[:upper:]'`",
		Severity: SeverityStyle,
		Description: "Zsh provides `${(U)var}` and `${(L)var}` parameter-expansion flags for " +
			"case conversion in-process. Spawning `tr` for this forks/execs per call (noticeable " +
			"in a hot loop), relies on the external `tr` being POSIX-compliant (BusyBox and old " +
			"macOS differ), and round-trips the data through a pipe. Drop `tr` for the " +
			"built-in: `upper=${(U)lower}` / `lower=${(L)upper}`.",
		Check: checkZC1553,
	})
}

func checkZC1553(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tr" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, strings.Trim(a.String(), "'\""))
	}

	// Need at least the from/to sets.
	if len(args) < 2 {
		return nil
	}
	var from, to string
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			continue
		}
		if from == "" {
			from = a
		} else if to == "" {
			to = a
			break
		}
	}

	isUpper := func(s string) bool { return s == "[:upper:]" || s == "A-Z" }
	isLower := func(s string) bool { return s == "[:lower:]" || s == "a-z" }
	if (isLower(from) && isUpper(to)) || (isUpper(from) && isLower(to)) {
		return []Violation{{
			KataID: "ZC1553",
			Message: "`tr` for case conversion — use Zsh `${(U)var}` / `${(L)var}` to avoid " +
				"the fork/exec and portability hazard.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}
	return nil
}
