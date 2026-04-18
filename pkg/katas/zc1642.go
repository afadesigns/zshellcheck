package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1642",
		Title:    "Warn on `tshark -w FILE` / `dumpcap -w FILE` without `-Z user` — capture file owned by root",
		Severity: SeverityWarning,
		Description: "Packet captures routinely need `CAP_NET_RAW`, so the capture process " +
			"typically runs as root. Without `-Z USER` the resulting pcap is root-owned — a " +
			"subsequent analyst who opens it with Wireshark (which can run helper scripts from " +
			"the file) operates on a root-owned file and may unintentionally invoke things as " +
			"root. `-Z USER` tells `tshark` / `dumpcap` to drop privileges for the actual " +
			"capture and write the file as `USER`.",
		Check: checkZC1642,
	})
}

func checkZC1642(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tshark" && ident.Value != "dumpcap" {
		return nil
	}

	var hasWrite, hasDrop bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-w" {
			hasWrite = true
		}
		if v == "-Z" {
			hasDrop = true
		}
	}
	if !hasWrite || hasDrop {
		return nil
	}

	return []Violation{{
		KataID: "ZC1642",
		Message: "`" + ident.Value + " -w FILE` without `-Z USER` leaves the pcap root-" +
			"owned. Add `-Z USER` to drop privileges for the capture.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
