// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// titleSeverity renders a lowercase severity ("error") in the prose form
// the docs use ("Error"). It avoids the deprecated strings.Title.
func titleSeverity(s katas.Severity) string {
	str := string(s)
	if str == "" {
		return ""
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

// printRulesList writes one line per kata (ID, severity, title), sorted
// by ID, to out and returns the process exit code. It backs `--list-rules`.
func printRulesList(out io.Writer, registry *katas.KatasRegistry) int {
	all := registry.AllKatas()
	for _, k := range all {
		fmt.Fprintf(out, "%s  %-7s  %s\n", k.ID, titleSeverity(k.Severity), k.Title)
	}
	fmt.Fprintf(out, "\n%d katas.\n", len(all))
	return 0
}

// printRuleExplain writes the full description of the kata identified by
// id to out, or an error to errOut when the ID is unknown. It backs
// `--explain ZC####` and returns the process exit code.
func printRuleExplain(out, errOut io.Writer, registry *katas.KatasRegistry, id string) int {
	id = strings.ToUpper(strings.TrimSpace(id))
	kata, ok := registry.GetKata(id)
	if !ok {
		fmt.Fprintf(errOut, "Unknown kata %q. Run 'zshellcheck --list-rules' to see every kata.\n", id)
		return 1
	}
	fmt.Fprintf(out, "%s — %s\n", kata.ID, kata.Title)
	fmt.Fprintf(out, "Severity: %s\n", titleSeverity(kata.Severity))
	if kata.Description != "" {
		fmt.Fprintf(out, "\n%s\n", kata.Description)
	}
	return 0
}
