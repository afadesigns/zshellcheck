package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1819",
		Title:    "Warn on `xattr -d com.apple.quarantine` / `xattr -cr` — removes macOS Gatekeeper quarantine",
		Severity: SeverityWarning,
		Description: "macOS sets the `com.apple.quarantine` extended attribute on every file " +
			"downloaded from the internet — Gatekeeper uses it to trigger the first-run " +
			"notarization / signature check. `xattr -d com.apple.quarantine FILE` strips the " +
			"attribute and lets the binary run with no prompt, and `xattr -cr DIR` does the " +
			"same recursively for every file in the tree. In a script that processes " +
			"downloaded artifacts this turns \"we vetted the binary\" into \"we trust whatever " +
			"landed in the download folder\". Verify the signature (`codesign --verify`) and " +
			"notarization (`spctl --assess --type execute`) first, or use " +
			"`xip`/`installer` packages so Gatekeeper stays in the loop.",
		Check: checkZC1819,
	})
}

func checkZC1819(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "xattr" {
		return nil
	}

	hasQuarantineDelete := false
	hasRecursiveClear := false

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-d" || v == "--delete" {
			if i+1 < len(cmd.Arguments) &&
				cmd.Arguments[i+1].String() == "com.apple.quarantine" {
				hasQuarantineDelete = true
			}
		}
		if strings.HasPrefix(v, "-") && !strings.HasPrefix(v, "--") && len(v) > 1 {
			hasC := false
			hasR := false
			for _, c := range v[1:] {
				if c == 'c' {
					hasC = true
				}
				if c == 'r' {
					hasR = true
				}
			}
			if hasC && hasR {
				hasRecursiveClear = true
			}
		}
	}

	if !hasQuarantineDelete && !hasRecursiveClear {
		return nil
	}
	return []Violation{{
		KataID: "ZC1819",
		Message: "`xattr -d com.apple.quarantine` / `-cr` strips the macOS Gatekeeper " +
			"quarantine — the binary runs with no signature / notarization check. " +
			"Verify with `codesign --verify` and `spctl --assess --type execute` " +
			"before stripping.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
