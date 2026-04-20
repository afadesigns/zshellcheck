package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1959",
		Title:    "Warn on `trivy … --skip-db-update` / `--skip-update` — scans against a stale vulnerability DB",
		Severity: SeverityWarning,
		Description: "`trivy` embeds a vulnerability database that is rehydrated on every scan " +
			"unless the operator passes `--skip-db-update` (or `--skip-update` on older " +
			"releases). In CI the flag is tempting — each build then skips a 40 MB download — " +
			"but the scan then misses every CVE disclosed since the cached DB was last " +
			"refreshed. Keep the default download, or pre-populate the cache with " +
			"`trivy image --download-db-only` once per day in a scheduled job, and only pass " +
			"`--skip-db-update` inside the same job so every scan sees the fresh data.",
		Check: checkZC1959,
	})
}

func checkZC1959(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `trivy --skip-db-update image` mangles to name=`skip-db-update`.
	if ident.Value == "skip-db-update" || ident.Value == "skip-update" {
		return zc1959Hit(cmd, "trivy --"+ident.Value)
	}
	if ident.Value != "trivy" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--skip-db-update" || v == "--skip-update" {
			return zc1959Hit(cmd, "trivy "+v)
		}
	}
	return nil
}

func zc1959Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1959",
		Message: "`" + form + "` scans against the cached DB — every CVE disclosed " +
			"since last refresh is missed. Keep the default download, or run " +
			"`trivy --download-db-only` once per day in a scheduled job.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
