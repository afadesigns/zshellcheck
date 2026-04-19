package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1805AwsDestructive = map[string]map[string]string{
	"cloudformation": {
		"delete-stack":     "removes every resource the stack manages, no rollback",
		"delete-stack-set": "deletes the stack set and all its instances",
	},
	"dynamodb": {
		"delete-table":  "drops the table and its data",
		"delete-backup": "drops the backup record",
	},
	"logs": {
		"delete-log-group":  "loses the audit trail in that group",
		"delete-log-stream": "drops the stream's events",
	},
	"kms": {
		"schedule-key-deletion": "queues CMK deletion — ciphertext becomes unreadable after the grace window",
	},
	"lambda": {
		"delete-function":             "removes the function and its versions",
		"delete-event-source-mapping": "drops the trigger wiring",
	},
	"ecr": {
		"delete-repository":  "deletes the image repository and every tag",
		"batch-delete-image": "drops tagged images in bulk",
	},
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1805",
		Title:    "Warn on `aws cloudformation delete-stack` / `dynamodb delete-table` / `logs delete-log-group` / `kms schedule-key-deletion` — destructive AWS state change",
		Severity: SeverityWarning,
		Description: "Each of these AWS actions drops state that AWS cannot restore: " +
			"`cloudformation delete-stack` tears down every resource the stack manages in " +
			"dependency order and has no rollback, `dynamodb delete-table` removes a table " +
			"and its items, `logs delete-log-group` erases the CloudWatch audit trail, and " +
			"`kms schedule-key-deletion` makes every ciphertext encrypted with the CMK " +
			"unreadable after the grace window. Add `--dry-run` where supported, stage the " +
			"call behind a typed confirmation, pin IDs through `--cli-input-json`, and " +
			"export backups (`dynamodb export-table-to-point-in-time`, `logs " +
			"create-export-task`) before pulling the trigger.",
		Check: checkZC1805,
	})
}

func checkZC1805(node ast.Node) []Violation {
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
	service := cmd.Arguments[0].String()
	action := cmd.Arguments[1].String()

	actions, ok := zc1805AwsDestructive[service]
	if !ok {
		return nil
	}
	note, ok := actions[action]
	if !ok {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		if arg.String() == "--dry-run" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1805",
		Message: "`aws " + service + " " + action + "` " + note + ". Stage a " +
			"confirmation, pin IDs via `--cli-input-json`, and export a backup " +
			"first where the service supports one.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
