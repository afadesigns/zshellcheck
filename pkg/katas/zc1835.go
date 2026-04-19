package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1835",
		Title:    "Warn on `smartctl -s off` — drive self-monitoring (SMART) disabled, silent failure",
		Severity: SeverityWarning,
		Description: "`smartctl -s off DEV` tells the drive firmware to stop recording the SMART " +
			"attribute counters that warn operators about pending failure — reallocated " +
			"sectors, pending sectors, uncorrectable errors, temperature excursions. " +
			"Rotating disks and SSDs both ship with the monitoring on; disabling it keeps " +
			"`smartctl -H` reporting PASSED right up until the drive falls off the bus, so " +
			"the periodic fleet health scan never escalates until data loss is already " +
			"happening. Use `smartctl -s on DEV` (default) and configure `smartd.conf` for " +
			"proactive alerts instead of muting the source.",
		Check: checkZC1835,
	})
}

func checkZC1835(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "smartctl" {
		return nil
	}
	args := cmd.Arguments
	for i := 0; i+1 < len(args); i++ {
		flag := args[i].String()
		if flag != "-s" && flag != "--smart" {
			continue
		}
		val := args[i+1].String()
		if val == "off" {
			return []Violation{{
				KataID: "ZC1835",
				Message: "`smartctl -s off` disables the drive's SMART attribute " +
					"collection — `smartctl -H` keeps reporting PASSED until the " +
					"disk falls off the bus. Leave it `on` and configure " +
					"`smartd.conf` for proactive alerts.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
