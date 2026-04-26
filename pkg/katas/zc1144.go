// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1144",
		Title:    "Avoid `trap` with signal numbers — use names",
		Severity: SeverityInfo,
		Description: "Signal numbers vary across platforms. Use signal names like " +
			"`SIGTERM`, `SIGINT`, `EXIT` instead of numeric values for portability.",
		Check: checkZC1144,
		Fix:   fixZC1144,
	})
}

// zc1144SignalNames maps POSIX signal numbers to their canonical
// names. Numbers above 31 (realtime signals) aren't included because
// their names vary across platforms and are rarely used with `trap`.
var zc1144SignalNames = map[string]string{
	"1":  "HUP",
	"2":  "INT",
	"3":  "QUIT",
	"4":  "ILL",
	"5":  "TRAP",
	"6":  "ABRT",
	"7":  "BUS",
	"8":  "FPE",
	"9":  "KILL",
	"10": "USR1",
	"11": "SEGV",
	"12": "USR2",
	"13": "PIPE",
	"14": "ALRM",
	"15": "TERM",
	"17": "CHLD",
	"18": "CONT",
	"19": "STOP",
	"20": "TSTP",
	"21": "TTIN",
	"22": "TTOU",
	"23": "URG",
	"24": "XCPU",
	"25": "XFSZ",
	"26": "VTALRM",
	"27": "PROF",
	"28": "WINCH",
	"29": "IO",
	"30": "PWR",
	"31": "SYS",
}

// fixZC1144 replaces numeric signal arguments in a `trap` call with
// their canonical names. Each numeric arg becomes a separate edit at
// that arg's position. Unknown numbers stay untouched.
func fixZC1144(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "trap" {
		return nil
	}
	var edits []FixEdit
	for i := 1; i < len(cmd.Arguments); i++ {
		arg := cmd.Arguments[i]
		val := arg.String()
		name, ok := zc1144SignalNames[val]
		if !ok {
			continue
		}
		tok := arg.TokenLiteralNode()
		off := LineColToByteOffset(source, tok.Line, tok.Column)
		if off < 0 || off+len(val) > len(source) {
			continue
		}
		if string(source[off:off+len(val)]) != val {
			continue
		}
		edits = append(edits, FixEdit{
			Line:    tok.Line,
			Column:  tok.Column,
			Length:  len(val),
			Replace: name,
		})
	}
	return edits
}

func checkZC1144(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "trap" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}

	// Check last arguments for numeric signal values
	for i := 1; i < len(cmd.Arguments); i++ {
		val := cmd.Arguments[i].String()
		// Numeric signals: 1-31
		isNumeric := len(val) > 0
		for _, ch := range val {
			if ch < '0' || ch > '9' {
				isNumeric = false
				break
			}
		}
		if isNumeric && val != "0" {
			return []Violation{{
				KataID: "ZC1144",
				Message: "Use signal names (`SIGTERM`, `SIGINT`, `EXIT`) instead of numbers in `trap`. " +
					"Signal numbers vary across platforms.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}
