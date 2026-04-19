package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1948",
		Title:    "Error on `ipmitool -P PASS` / `-E` — BMC password visible in argv",
		Severity: SeverityError,
		Description: "`ipmitool -H <bmc> -U admin -P <password>` puts the BMC credential into " +
			"`ps`, `/proc/PID/cmdline`, and every process-dump crash file. The BMC is a root-" +
			"equivalent out-of-band controller (power, console, firmware update), so that " +
			"password is one of the most sensitive tokens on the host. Use `-f <password_file>` " +
			"(mode `0400`, owned by the automation user) or set `IPMI_PASSWORD` and pass `-E` — " +
			"`ipmitool` reads the env var but never echoes it.",
		Check: checkZC1948,
	})
}

func checkZC1948(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ipmitool" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-P" && i+1 < len(cmd.Arguments) {
			return zc1948Hit(cmd, "-P "+cmd.Arguments[i+1].String())
		}
		if strings.HasPrefix(v, "-P") && len(v) > 2 {
			return zc1948Hit(cmd, v)
		}
	}
	return nil
}

func zc1948Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1948",
		Message: "`ipmitool " + form + "` leaks the BMC password into argv — visible in " +
			"`ps`/`/proc`/crash dumps. Use `-f <password_file>` (mode 0400) or " +
			"`IPMI_PASSWORD=… ipmitool -E`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
