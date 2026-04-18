package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1670Dangerous = map[string]struct{}{
	"allow_execstack":           {},
	"allow_execmod":             {},
	"allow_execmem":             {},
	"httpd_execmem":             {},
	"httpd_unified":             {},
	"selinuxuser_execmod":       {},
	"selinuxuser_execstack":     {},
	"selinuxuser_execheap":      {},
	"domain_kernel_load_modules": {},
	"deny_ptrace":               {}, // flipping this OFF re-enables ptrace globally
	"mmap_low_allowed":          {},
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1670",
		Title:    "Warn on `setsebool -P` enabling memory-protection-relaxing SELinux boolean",
		Severity: SeverityWarning,
		Description: "Specific SELinux policy booleans (`allow_execstack`, `allow_execmem`, " +
			"`httpd_execmem`, `selinuxuser_execstack`, `domain_kernel_load_modules`, " +
			"`mmap_low_allowed`, etc.) relax per-domain memory protections that the policy " +
			"puts in place precisely because those domains should not need writable-and-" +
			"executable pages. Persisting the flip with `-P` carries the regression across " +
			"reboots. Fix the underlying binary (`execstack -c`, `chcon`, stop generating " +
			"runtime-JIT code in the wrong domain) instead of loosening policy.",
		Check: checkZC1670,
	})
}

func checkZC1670(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setsebool" {
		return nil
	}

	hasPersist := false
	boolName := ""
	boolValue := ""
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch {
		case v == "-P":
			hasPersist = true
		case boolName == "":
			boolName = v
		case boolValue == "":
			boolValue = v
		}
	}

	if !hasPersist || boolName == "" || boolValue == "" {
		return nil
	}
	if _, dangerous := zc1670Dangerous[boolName]; !dangerous {
		return nil
	}
	if boolValue != "1" && boolValue != "on" && boolValue != "true" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1670",
		Message: "`setsebool -P " + boolName + " " + boolValue + "` persistently relaxes " +
			"SELinux memory-protection policy — fix the binary instead (`execstack -c`, " +
			"relabel with `chcon`, or change the domain).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
