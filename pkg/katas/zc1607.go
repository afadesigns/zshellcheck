package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1607",
		Title:    "Warn on `git config safe.directory '*'` — disables CVE-2022-24765 protection",
		Severity: SeverityWarning,
		Description: "`safe.directory` is git's mitigation for CVE-2022-24765 (fake git dirs " +
			"planted by another uid). Setting it to `'*'` trusts every directory on the host " +
			"— an attacker who creates `/tmp/evil/.git` with a malicious `core.fsmonitor` hook " +
			"gets arbitrary code execution the first time any user runs `git status` near that " +
			"path. List the specific paths that need cross-owner git access instead, or fix " +
			"the underlying ownership mismatch.",
		Check: checkZC1607,
	})
}

func checkZC1607(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "git" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "safe.directory=*" || strings.HasPrefix(v, "safe.directory=*") {
			return violationZC1607(cmd)
		}
		if v == "safe.directory" && i+1 < len(cmd.Arguments) {
			next := strings.Trim(cmd.Arguments[i+1].String(), "'\"")
			if next == "*" {
				return violationZC1607(cmd)
			}
		}
	}
	return nil
}

func violationZC1607(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1607",
		Message: "`git config safe.directory '*'` trusts every directory — defeats CVE-" +
			"2022-24765 protection. List specific paths, or fix the ownership mismatch.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
