// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1804Ec2Destructive = map[string]string{
	"terminate-instances":      "tears down EC2 instance(s) and their instance-store volumes",
	"delete-volume":            "deletes the EBS volume and its data",
	"delete-snapshot":          "deletes the EBS / RDS snapshot",
	"delete-vpc":               "removes the VPC along with its routing / dependencies",
	"delete-internet-gateway":  "detaches / removes the IGW",
	"delete-network-interface": "removes the ENI",
	"delete-security-group":    "removes the security group",
	"delete-launch-template":   "removes the launch template",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1804",
		Title:    "Warn on `aws ec2 terminate-instances` / `delete-volume` / `delete-snapshot` — destructive cloud state change",
		Severity: SeverityWarning,
		Description: "AWS EC2 destructive actions (`terminate-instances`, `delete-volume`, " +
			"`delete-snapshot`, `delete-vpc`, and friends) drop cloud state without any " +
			"automatic backup: instance-store volumes vanish on terminate, EBS volumes and " +
			"snapshots cannot be restored from the AWS side once deleted, and a wrong " +
			"VPC / ENI / security-group ID can take down workloads in the same account. " +
			"Review the target list with `aws ec2 describe-…`, pair destructive commands " +
			"with `--dry-run`, and keep the IDs pinned in a file that `aws ... --cli-input-" +
			"json` can consume rather than passing them inline.",
		Check: checkZC1804,
	})
}

func checkZC1804(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "aws" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "ec2" {
		return nil
	}
	action := cmd.Arguments[1].String()
	note, ok := zc1804Ec2Destructive[action]
	if !ok {
		return nil
	}

	// `--dry-run` makes the command a no-op. Allow it.
	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if v == "--dry-run" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1804",
		Message: "`aws ec2 " + action + "` " + note + " with no automatic backup. " +
			"Review with `aws ec2 describe-…`, add `--dry-run` to verify the target, " +
			"and pin IDs through `--cli-input-json`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
