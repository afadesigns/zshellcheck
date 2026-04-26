// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1839",
		Title:    "Warn on `timedatectl set-ntp false` / disabling `systemd-timesyncd` / `chronyd`",
		Severity: SeverityWarning,
		Description: "`timedatectl set-ntp false` (also spelled `set-ntp no` / `set-ntp 0`) tells " +
			"systemd to stop the network time client; `systemctl disable systemd-timesyncd` " +
			"and `systemctl disable chronyd` / `ntpd` have the same effect. With no time " +
			"source the hardware clock drifts, and within days TLS handshakes begin failing " +
			"`notBefore`/`notAfter` checks, Kerberos tickets refuse to validate, time-based " +
			"one-time passwords go out of sync, and log entries arrive in the wrong order — " +
			"all silently, because the original command succeeded. Keep NTP enabled in " +
			"production; if you really need a frozen clock for reproducibility, isolate it " +
			"to a namespace or CI container rather than the host.",
		Check: checkZC1839,
	})
}

var zc1839DisableServices = map[string]struct{}{
	"systemd-timesyncd":         {},
	"systemd-timesyncd.service": {},
	"chronyd":                   {},
	"chronyd.service":           {},
	"chrony":                    {},
	"chrony.service":            {},
	"ntpd":                      {},
	"ntpd.service":              {},
	"ntp":                       {},
	"ntp.service":               {},
}

func checkZC1839(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "timedatectl":
		if len(cmd.Arguments) < 2 {
			return nil
		}
		if cmd.Arguments[0].String() != "set-ntp" {
			return nil
		}
		val := strings.ToLower(cmd.Arguments[1].String())
		if val == "false" || val == "no" || val == "0" || val == "off" {
			return zc1839Hit(cmd, "timedatectl set-ntp "+cmd.Arguments[1].String())
		}
	case "systemctl":
		if len(cmd.Arguments) < 2 {
			return nil
		}
		action := cmd.Arguments[0].String()
		if action != "disable" && action != "mask" && action != "stop" {
			return nil
		}
		for _, arg := range cmd.Arguments[1:] {
			svc := strings.ToLower(arg.String())
			if _, hit := zc1839DisableServices[svc]; hit {
				return zc1839Hit(cmd, "systemctl "+action+" "+arg.String())
			}
		}
	}
	return nil
}

func zc1839Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1839",
		Message: "`" + where + "` turns off network time sync — clock drift " +
			"breaks TLS `notBefore`/`notAfter`, Kerberos, and TOTP. Leave NTP " +
			"enabled; isolate frozen clocks to namespaces/CI.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
