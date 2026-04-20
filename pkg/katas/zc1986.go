package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1986",
		Title:    "Warn on `touch -d` / `-t` / `-r` — explicit timestamp write is a common antiforensics pattern",
		Severity: SeverityWarning,
		Description: "`touch -d \"2 years ago\" $F`, `touch -t YYYYMMDDhhmm $F`, and `touch -r " +
			"$REF $F` all write the atime/mtime to a specific value rather than the " +
			"current clock. Legitimate uses exist — re-stamping a mirror to match " +
			"upstream, generating deterministic tarballs for reproducible-build " +
			"pipelines, `rsync --archive` edge cases — but the pattern also matches the " +
			"classic \"age the dropped file\" antiforensics trick where an attacker " +
			"normalises a new binary to look as old as its neighbours so `find -mtime`- " +
			"based triage misses it. Audit rules should flag these forms in production " +
			"scripts; in reproducible-build contexts, keep the timestamp derived from " +
			"`SOURCE_DATE_EPOCH` via `touch -d @$SOURCE_DATE_EPOCH` so operators can " +
			"recognise the intent at a glance.",
		Check: checkZC1986,
	})
}

func checkZC1986(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "touch" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-d", "-t", "-r":
			return []Violation{{
				KataID: "ZC1986",
				Message: "`touch " + v + "` writes a specific atime/mtime — also the " +
					"classic \"age the dropped file\" antiforensics pattern. Derive " +
					"from `$SOURCE_DATE_EPOCH` where the intent is deterministic.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
