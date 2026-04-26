// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1824",
		Title:    "Warn on `kubectl drain --disable-eviction` — bypasses PodDisruptionBudget via raw DELETE",
		Severity: SeverityWarning,
		Description: "`kubectl drain --disable-eviction` tells the client to delete pods directly " +
			"via the API instead of issuing Eviction requests. The Eviction pathway is what " +
			"honours PodDisruptionBudget — `--disable-eviction` drops pods regardless of the " +
			"minAvailable / maxUnavailable contract the workload owner defined. On a " +
			"multi-replica service this turns a rolling drain into a hard outage. Fix the " +
			"blocking PDB (raise minAvailable, wait for replicas to reschedule, or negotiate " +
			"with the owner) instead of flipping the flag off.",
		Check: checkZC1824,
	})
}

func checkZC1824(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubectl" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "drain" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--disable-eviction" || v == "--disable-eviction=true" {
			return []Violation{{
				KataID: "ZC1824",
				Message: "`kubectl drain --disable-eviction` deletes pods via raw API " +
					"DELETE — PodDisruptionBudgets are ignored and the workload " +
					"owner's availability contract is voided. Fix the blocking PDB " +
					"instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
