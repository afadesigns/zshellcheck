package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1703NetHardening = map[string]string{
	"rp_filter=0":                         "reverse-path filtering (anti-spoofing)",
	"accept_source_route=1":               "source-routed packet acceptance",
	"accept_redirects=1":                  "ICMP redirect acceptance (routing tampering)",
	"send_redirects=1":                    "ICMP redirect emission",
	"icmp_echo_ignore_broadcasts=0":       "ICMP broadcast ignore (enables smurf amplification)",
	"icmp_ignore_bogus_error_responses=0": "bogus ICMP error ignore",
	"log_martians=0":                      "martian-packet logging",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1703",
		Title:    "Warn on `sysctl -w` disabling network-hardening knobs",
		Severity: SeverityWarning,
		Description: "Several `net.ipv4.*` / `net.ipv6.*` sysctl knobs exist specifically to " +
			"harden the host against on-link spoofing, ICMP redirect tampering, smurf " +
			"amplification, and source-routed packets — `rp_filter=1`, " +
			"`accept_source_route=0`, `accept_redirects=0`, `send_redirects=0`, " +
			"`icmp_echo_ignore_broadcasts=1`, `log_martians=1`. Flipping any of them to " +
			"the lax value (rp_filter=0, accept_source_route=1, …) re-opens classic " +
			"layer-3 attacks. Leave the protective defaults in place; if a niche workload " +
			"really needs relaxed filtering, scope the change per-interface with a comment " +
			"explaining why.",
		Check: checkZC1703,
	})
}

func checkZC1703(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "sysctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		for suffix, note := range zc1703NetHardening {
			if !strings.HasSuffix(v, suffix) {
				continue
			}
			if !strings.HasPrefix(v, "net.") {
				continue
			}
			return []Violation{{
				KataID: "ZC1703",
				Message: "`sysctl " + v + "` disables " + note + " — classic layer-3 " +
					"attacks (spoofing / smurf / redirect tamper) reopen. Keep the " +
					"protective default.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
