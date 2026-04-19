package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1811",
		Title:    "Error on `chown/chmod/chgrp --no-preserve-root` — disables GNU safeguard against recursive `/`",
		Severity: SeverityError,
		Description: "GNU `chown`, `chmod`, and `chgrp` refuse to recurse into `/` by default " +
			"(`--preserve-root` in coreutils). `--no-preserve-root` opts in to walking the " +
			"entire filesystem, so a stray `$PATH` expansion or wrong variable combined with " +
			"`-R` rewrites ownership or mode on every file on the host. The flag has no " +
			"legitimate script use — if a specific top-level target genuinely needs recursion, " +
			"list that path explicitly and keep the safeguard in place.",
		Check: checkZC1811,
	})
}

func checkZC1811(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: leading `--no-preserve-root` mangles name to `no-preserve-root`.
	if ident.Value == "no-preserve-root" {
		return []Violation{{
			KataID: "ZC1811",
			Message: "`--no-preserve-root` disables the GNU safeguard against recursing " +
				"into `/`. Remove the flag; list explicit paths instead.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}

	if ident.Value != "chown" && ident.Value != "chmod" && ident.Value != "chgrp" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if arg.String() == "--no-preserve-root" {
			return []Violation{{
				KataID: "ZC1811",
				Message: "`" + ident.Value + " --no-preserve-root` disables the GNU " +
					"safeguard against recursing into `/`. Remove the flag; list " +
					"explicit paths instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
