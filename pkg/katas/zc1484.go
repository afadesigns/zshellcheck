package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1484",
		Title:    "Error on `npm/yarn/pnpm config set strict-ssl false` — disables registry TLS verification",
		Severity: SeverityError,
		Description: "Turning off `strict-ssl` for npm, yarn, or pnpm makes the client accept any " +
			"TLS certificate from the registry — a MITM (corporate proxy, compromised WiFi, rogue " +
			"BGP) can substitute any package, including new versions of `react` or `lodash`. If " +
			"the registry uses a private CA, point `cafile` / `NODE_EXTRA_CA_CERTS` at the right " +
			"bundle instead.",
		Check: checkZC1484,
	})
}

func checkZC1484(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "npm" && ident.Value != "yarn" && ident.Value != "pnpm" &&
		ident.Value != "bun" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// Form A: `npm config set strict-ssl false`
	for i := 0; i+2 < len(args); i++ {
		if args[i] == "config" && args[i+1] == "set" {
			j := i + 2
			// Skip optional scope flags like `--global` / `-g`.
			for j < len(args) && strings.HasPrefix(args[j], "-") {
				j++
			}
			if j+1 < len(args) && args[j] == "strict-ssl" {
				val := strings.ToLower(args[j+1])
				if val == "false" || val == "0" || val == "no" {
					return zc1484Violation(cmd)
				}
			}
			if j < len(args) && strings.HasPrefix(strings.ToLower(args[j]), "strict-ssl=") {
				val := strings.ToLower(strings.TrimPrefix(args[j], "strict-ssl="))
				if val == "false" || val == "0" || val == "no" {
					return zc1484Violation(cmd)
				}
			}
		}
	}

	// Form B: `npm install --strict-ssl=false` (one-shot)
	for _, v := range args {
		if strings.EqualFold(v, "--strict-ssl=false") ||
			strings.EqualFold(v, "--no-strict-ssl") {
			return zc1484Violation(cmd)
		}
	}
	return nil
}

func zc1484Violation(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1484",
		Message: "`strict-ssl=false` disables npm/yarn/pnpm registry TLS verification — any " +
			"MITM swaps packages. Point `cafile` at the right CA bundle instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
