package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1470",
		Title:    "Error on `git config http.sslVerify false` / `git -c http.sslVerify=false`",
		Severity: SeverityError,
		Description: "Disabling `http.sslVerify` in git means every subsequent fetch / clone " +
			"accepts any TLS certificate — MITM trivially replaces the tree you are cloning with " +
			"attacker-controlled code. Fix the broken CA instead: install the certificate, " +
			"point at the right store with `GIT_SSL_CAINFO`, or use an SSH transport.",
		Check: checkZC1470,
	})
}

func checkZC1470(node ast.Node) []Violation {
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

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// Form A: `git config [--global|--local|--system] http.sslVerify false`
	for i := 0; i < len(args); i++ {
		if args[i] == "config" {
			// Walk following args, skip scope flag if present.
			j := i + 1
			for j < len(args) && strings.HasPrefix(args[j], "--") && args[j] != "--" {
				j++
			}
			if j+1 < len(args) && strings.EqualFold(args[j], "http.sslVerify") {
				if strings.EqualFold(args[j+1], "false") || args[j+1] == "0" {
					return zc1470Violation(cmd)
				}
			}
		}
	}

	// Form B: `git -c http.sslVerify=false ...`
	for i := 0; i < len(args); i++ {
		if args[i] == "-c" && i+1 < len(args) {
			kv := args[i+1]
			if k, v, ok := strings.Cut(kv, "="); ok {
				if strings.EqualFold(k, "http.sslVerify") &&
					(strings.EqualFold(v, "false") || v == "0") {
					return zc1470Violation(cmd)
				}
			}
		}
	}

	return nil
}

func zc1470Violation(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1470",
		Message: "`http.sslVerify=false` disables TLS verification — any MITM swaps the " +
			"clone for attacker code. Fix the CA instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
