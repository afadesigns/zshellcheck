// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strconv"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1800",
		Title:    "Warn on `pg_ctl stop -m immediate` — abrupt shutdown skips checkpoint, forces WAL recovery",
		Severity: SeverityWarning,
		Description: "`pg_ctl stop -m immediate` sends `SIGQUIT` to the postmaster. Server " +
			"processes drop connections, no checkpoint is taken, and buffered changes are " +
			"left in memory. Recovery on the next start has to replay every record since the " +
			"last checkpoint; if WAL is corrupt, lost, or on different storage, committed " +
			"transactions can be lost. Use `-m smart` (default) or `-m fast` so the server " +
			"issues a shutdown checkpoint and closes cleanly; reserve `immediate` for the " +
			"\"the node is on fire\" case and pair it with a tested PITR procedure.",
		Check: checkZC1800,
	})
}

func checkZC1800(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || CommandIdentifier(cmd) != "pg_ctl" {
		return nil
	}
	if !zc1800StopOrRestart(cmd) || !zc1800ImmediateMode(cmd) {
		return nil
	}
	return []Violation{{
		KataID: "ZC1800",
		Message: "`pg_ctl stop -m immediate` kills the postmaster without a shutdown " +
			"checkpoint — WAL replay on restart can lose committed transactions " +
			"if WAL is corrupt. Use `-m smart` or `-m fast` for routine shutdowns.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func zc1800StopOrRestart(cmd *ast.SimpleCommand) bool {
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "stop" || v == "restart" {
			return true
		}
	}
	return false
}

func zc1800ImmediateMode(cmd *ast.SimpleCommand) bool {
	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--mode=immediate" {
			return true
		}
		if v == "-m" && i+1 < len(cmd.Arguments) && cmd.Arguments[i+1].String() == "immediate" {
			return true
		}
		if strings.HasPrefix(v, "-m") && len(v) > 2 && v[2:] == "immediate" {
			return true
		}
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1801",
		Title:    "Warn on `fwupdmgr update` / `install` — mid-flash interruption can brick firmware",
		Severity: SeverityWarning,
		Description: "`fwupdmgr update`, `fwupdmgr upgrade`, and `fwupdmgr install FIRMWARE` push " +
			"new firmware into BIOS / UEFI, SSD, Thunderbolt controller, NIC, or dock " +
			"microcontroller. Most of those devices have no A/B rollback — an interrupted " +
			"flash (power cut, unexpected reboot, PSU toggle) leaves the chip in an " +
			"unbootable state that needs vendor-recovery hardware. Run from a battery-backed " +
			"session, mask reboot triggers with `systemd-inhibit`, pin the power supply, and " +
			"verify the update history with `fwupdmgr get-history` once the device returns.",
		Check: checkZC1801,
	})
}

func checkZC1801(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "fwupdmgr" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	switch cmd.Arguments[0].String() {
	case "update", "upgrade", "install", "reinstall", "downgrade":
		return []Violation{{
			KataID: "ZC1801",
			Message: "`fwupdmgr " + cmd.Arguments[0].String() + "` flashes firmware — a " +
				"mid-write interruption can brick BIOS, SSD, Thunderbolt, or NIC " +
				"microcontrollers. Inhibit reboot triggers (`systemd-inhibit`) and " +
				"ensure battery / UPS before running.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1802",
		Title:    "Warn on `dnf history undo N` / `rollback N` — reverses transactions without compat check",
		Severity: SeverityWarning,
		Description: "`dnf history undo N` reverts the exact package set of transaction N — every " +
			"install turns into a remove, every remove into an install, every update into a " +
			"downgrade. `dnf history rollback N` does the same for every transaction after " +
			"N. Neither checks that the older versions still resolve cleanly against the " +
			"current package graph: dependencies that moved forward for other reasons end up " +
			"downgraded alongside, security patches get reverted, and services whose " +
			"configuration was migrated fail to start on the older binary. Review the plan " +
			"with `dnf history info N`, pin the rollback scope with `--exclude=` / `--assumeyes` " +
			"only after review, or restore from a filesystem snapshot taken before the " +
			"original transaction.",
		Check: checkZC1802,
	})
}

func checkZC1802(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dnf" && ident.Value != "yum" && ident.Value != "dnf5" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "history" {
		return nil
	}
	action := cmd.Arguments[1].String()
	if action != "undo" && action != "rollback" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1802",
		Message: "`" + ident.Value + " history " + action + "` reverses the past " +
			"transaction — deps downgrade, security patches can get reverted. " +
			"Review with `dnf history info`, or restore a filesystem snapshot.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

var zc1803MySQLClients = map[string]bool{
	"mysql":         true,
	"mysqldump":     true,
	"mysqladmin":    true,
	"mariadb":       true,
	"mariadb-dump":  true,
	"mariadb-admin": true,
}

var zc1803PgClients = map[string]bool{
	"psql":       true,
	"pg_dump":    true,
	"pgbench":    true,
	"pg_restore": true,
}

var zc1803MySQLFlags = map[string]bool{
	"--skip-ssl":          true,
	"--ssl=0":             true,
	"--ssl=false":         true,
	"--ssl-mode=disabled": true,
	"--ssl-mode=DISABLED": true,
	"--ssl-mode=disable":  true,
	"--ssl-mode=DISABLE":  true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1803",
		Title:    "Error on `mysql --skip-ssl` / `psql sslmode=disable` — plaintext credentials on the wire",
		Severity: SeverityError,
		Description: "Disabling TLS on a MySQL or PostgreSQL client pushes the login handshake " +
			"(including the password or auth challenge) and every subsequent query and " +
			"result over plaintext TCP. Anyone in the network path — the cloud VPC, the " +
			"office LAN, a compromised router — can sniff or modify the stream. The flags " +
			"vary (`--skip-ssl`, `--ssl=0`, `--ssl-mode=DISABLED` for MySQL / MariaDB; " +
			"`sslmode=disable` in the connection URI or `PGSSLMODE=disable` env var for " +
			"PostgreSQL) but the effect is the same. Prefer `--ssl-mode=VERIFY_IDENTITY` " +
			"(MySQL 8+) and `sslmode=verify-full` (psql) with a pinned CA bundle.",
		Check: checkZC1803,
	})
}

func checkZC1803(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if zc1803MySQLClients[ident.Value] {
		for _, arg := range cmd.Arguments {
			raw := strings.Trim(arg.String(), "\"'")
			v := strings.ToLower(raw)
			if v == "--skip-ssl" || v == "--ssl=0" || v == "--ssl=false" ||
				v == "--ssl-mode=disabled" || v == "--ssl-mode=disable" {
				return zc1803HitMySQL(cmd, ident.Value, raw)
			}
		}
	}

	if zc1803PgClients[ident.Value] {
		for _, arg := range cmd.Arguments {
			raw := strings.Trim(arg.String(), "\"'")
			if strings.Contains(strings.ToLower(raw), "sslmode=disable") {
				return zc1803HitPg(cmd, ident.Value, raw)
			}
		}
	}
	return nil
}

func zc1803HitMySQL(cmd *ast.SimpleCommand, tool, flag string) []Violation {
	line, col := FlagArgPosition(cmd, zc1803MySQLFlags)
	return []Violation{{
		KataID: "ZC1803",
		Message: "`" + tool + " " + flag + "` disables TLS — login handshake and " +
			"queries travel in plaintext. Use `--ssl-mode=VERIFY_IDENTITY` (MySQL) / " +
			"`sslmode=verify-full` (psql) with a pinned CA.",
		Line:   line,
		Column: col,
		Level:  SeverityError,
	}}
}

func zc1803HitPg(cmd *ast.SimpleCommand, tool, flag string) []Violation {
	return []Violation{{
		KataID: "ZC1803",
		Message: "`" + tool + " " + flag + "` disables TLS — login handshake and " +
			"queries travel in plaintext. Use `--ssl-mode=VERIFY_IDENTITY` (MySQL) / " +
			"`sslmode=verify-full` (psql) with a pinned CA.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

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

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1806",
		Title:    "Warn on `zmv 'PAT' 'REP'` without `-n` / `-i` — silent bulk rename",
		Severity: SeverityWarning,
		Description: "`zmv` (autoloaded from Zsh's functions) rewrites every filename that matches " +
			"the pattern in one shot. A small typo in the source pattern or replacement — " +
			"`*.jpg` vs `*.JPG`, a misplaced `(..)`, forgetting `**` recursion — can collide " +
			"names and silently overwrite files, since `zmv` aborts the batch only on its " +
			"own conflict check, not on semantic errors. Use `zmv -n 'PAT' 'REP'` first to " +
			"see the rename list, or `zmv -i` to prompt per file. Only drop the guard once " +
			"the preview matches what you expect.",
		Check: checkZC1806,
	})
}

func checkZC1806(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || CommandIdentifier(cmd) != "zmv" || len(cmd.Arguments) < 2 {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if zc1806HasGuardFlag(arg.String()) {
			return nil
		}
	}
	return []Violation{{
		KataID: "ZC1806",
		Message: "`zmv` without `-n` (dry-run) or `-i` (interactive) renames every " +
			"matched file in one shot — a pattern typo can collide names. Preview " +
			"with `zmv -n`, then re-run once the list looks right.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func zc1806HasGuardFlag(v string) bool {
	switch v {
	case "-n", "--dry-run", "-i", "--interactive":
		return true
	}
	if len(v) <= 1 || v[0] != '-' || v[1] == '-' {
		return false
	}
	for _, c := range v[1:] {
		if c == 'n' || c == 'i' {
			return true
		}
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1807",
		Title:    "Warn on `gh api -X DELETE` — raw GitHub DELETE bypasses `gh` command confirmations",
		Severity: SeverityWarning,
		Description: "`gh api -X DELETE /repos/OWNER/REPO` (and `--method=DELETE` variants) " +
			"sends a raw GitHub API request with the caller's token. There is no confirmation " +
			"prompt, no `--yes` guard, and no friendly dry-run — a script that builds the " +
			"path from a variable can wipe repos, releases, deploy keys, workflow runs, " +
			"issue comments, or whole organisations in one call. Use the high-level `gh` " +
			"subcommand for the target (`gh repo delete`, `gh release delete`, `gh workflow " +
			"disable`) which still at least requires `--yes`, or wrap the raw call with a " +
			"preflight `gh api -X GET /path` and an explicit confirmation in the script.",
		Check: checkZC1807,
	})
}

func checkZC1807(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gh" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "api" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		switch {
		case v == "-X" || v == "--method":
			if 1+i+1 < len(cmd.Arguments) {
				next := cmd.Arguments[1+i+1].String()
				if strings.EqualFold(next, "DELETE") {
					return zc1807Hit(cmd)
				}
			}
		case strings.EqualFold(v, "-XDELETE"):
			return zc1807Hit(cmd)
		case strings.EqualFold(v, "--method=DELETE"):
			return zc1807Hit(cmd)
		}
	}
	return nil
}

func zc1807Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1807",
		Message: "`gh api -X DELETE` sends a raw DELETE to the GitHub API with the " +
			"caller's token — no `--yes` guard, no dry-run. Use the high-level `gh` " +
			"subcommand for the target, or wrap with a preflight GET + explicit " +
			"confirmation.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1808",
		Title:    "Warn on `kubectl replace --force` — deletes + recreates resource, drops running pods",
		Severity: SeverityWarning,
		Description: "`kubectl replace --force -f FILE` is `delete` followed by `create`: the " +
			"existing resource (and every dependent pod / replicaset / endpoint) is removed " +
			"before the new manifest is applied. In-flight requests drop, PodDisruptionBudget " +
			"is ignored, and controllers that watch the object see it disappear and reappear. " +
			"Prefer `kubectl apply -f FILE` — same manifest, server-side merge that preserves " +
			"running pods — and reach for `replace --force` only when the resource schema has " +
			"changed in a way `apply` cannot patch, with traffic drained beforehand.",
		Check: checkZC1808,
	})
}

func checkZC1808(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubectl" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "replace" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--force" || v == "-f=--force" {
			return []Violation{{
				KataID: "ZC1808",
				Message: "`kubectl replace --force` is delete + create — pods die, " +
					"PDBs are ignored, in-flight requests drop. Prefer `kubectl " +
					"apply -f FILE` and reserve `replace --force` for schema changes " +
					"`apply` cannot patch, after draining traffic.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1809",
		Title:    "Error on `gsutil rm -r gs://…` / `gsutil rb -f gs://…` — bulk GCS deletion",
		Severity: SeverityError,
		Description: "`gsutil rm -r gs://bucket/prefix` and `gsutil rm -rf gs://bucket` delete " +
			"every object under the prefix — with `-m` (parallel) they do it faster than any " +
			"undo window. `gsutil rb -f gs://bucket` removes the bucket after force-deleting " +
			"the contents. Neither soft-deletes; Object Versioning can help only if it is " +
			"turned on in advance, and `gsutil rb` leaves no retention grace. Preview with " +
			"`gsutil ls`, enable Object Versioning or retention locks before the fact, and " +
			"prefer narrower `gsutil rm gs://bucket/specific-object` calls.",
		Check: checkZC1809,
	})
}

func checkZC1809(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gsutil" {
		return nil
	}

	subIdx := -1
	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "rm" || v == "rb" {
			subIdx = i
			break
		}
	}
	if subIdx == -1 {
		return nil
	}
	sub := cmd.Arguments[subIdx].String()

	hasDestFlag := false
	for _, arg := range cmd.Arguments[subIdx+1:] {
		v := arg.String()
		if sub == "rm" && (v == "-r" || v == "-R" || v == "-rf" || v == "-fr" || v == "--recursive") {
			hasDestFlag = true
			break
		}
		if sub == "rb" && (v == "-f" || v == "--force") {
			hasDestFlag = true
			break
		}
	}
	if !hasDestFlag {
		return nil
	}
	return []Violation{{
		KataID: "ZC1809",
		Message: "`gsutil " + sub + "` with recursive/force deletes every matching " +
			"GCS object (or the bucket itself). Preview with `gsutil ls`, enable " +
			"Object Versioning / retention locks ahead of time, and prefer narrower " +
			"object-level `gsutil rm` calls.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

var zc1810RecursiveFlags = map[string]bool{
	"-r":          true,
	"--recursive": true,
	"-m":          true,
	"--mirror":    true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1810",
		Title:    "Warn on `wget -r` / `--mirror` without `--level=N` — unbounded recursive download",
		Severity: SeverityWarning,
		Description: "`wget -r` and `wget --mirror` (short `-m`) follow links to arbitrary depth. " +
			"Without `--level=N` or `-l N` the crawl keeps going until `wget` hits the " +
			"remote server's limits, fills the local disk, or climbs into a parent directory " +
			"the author did not intend to mirror (add `--no-parent` to block that too). " +
			"Pin a depth (`--level=3`), restrict siblings (`--no-parent`, `--accept=` / " +
			"`--reject=`), and cap the byte budget (`--quota=1G`) before running a recursive " +
			"wget in automation.",
		Check: checkZC1810,
	})
}

func checkZC1810(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "wget" {
		return nil
	}

	hasRecursive := false
	for _, arg := range cmd.Arguments {
		if zc1810RecursiveFlags[arg.String()] {
			hasRecursive = true
			break
		}
	}
	if !hasRecursive {
		return nil
	}
	if zc1810HasLevel(cmd) {
		return nil
	}
	return zc1810Hit(cmd)
}

func zc1810HasLevel(cmd *ast.SimpleCommand) bool {
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch {
		case v == "-l" || v == "--level":
			return true
		case strings.HasPrefix(v, "--level="):
			return true
		case strings.HasPrefix(v, "-l") && len(v) > 2:
			return true
		}
	}
	return false
}

func zc1810Hit(cmd *ast.SimpleCommand) []Violation {
	line, col := FlagArgPosition(cmd, zc1810RecursiveFlags)
	return []Violation{{
		KataID: "ZC1810",
		Message: "`wget -r` / `--mirror` without `--level=N` follows links to " +
			"arbitrary depth — the crawl can exhaust disk and climb into parent " +
			"paths. Pin `--level=3`, add `--no-parent`, and cap with `--quota=1G`.",
		Line:   line,
		Column: col,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1811",
		Title:    "Error on `chown/chmod/chgrp --no-preserve-root` — disables GNU safeguard against recursive `/`",
		Severity: SeverityError,
		Description: "GNU `chown`, `chmod`, and `chgrp` refuse to recurse into `/` by default " +
			"(`--preserve-root` in coreutils). `--no-preserve-root` opts in to walking the " +
			"entire filesystem, so a stray `$PATH` expansion or wrong variable combined with " +
			"`-R` rewrites ownership or mode on every file on the host. The flag has no " +
			"legitimate script use — if a specific top-level target genuinely needs recursion, " +
			"list that path explicitly and keep the safeguard in place.",
		Check: checkZC1811,
	})
}

func checkZC1811(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: leading `--no-preserve-root` mangles name to `no-preserve-root`.
	if ident.Value == "no-preserve-root" {
		return []Violation{{
			KataID: "ZC1811",
			Message: "`--no-preserve-root` disables the GNU safeguard against recursing " +
				"into `/`. Remove the flag; list explicit paths instead.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}

	if ident.Value != "chown" && ident.Value != "chmod" && ident.Value != "chgrp" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if arg.String() == "--no-preserve-root" {
			return []Violation{{
				KataID: "ZC1811",
				Message: "`" + ident.Value + " --no-preserve-root` disables the GNU " +
					"safeguard against recursing into `/`. Remove the flag; list " +
					"explicit paths instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1812",
		Title:    "Error on `aws ssm put-parameter --type SecureString --value SECRET` — plaintext in argv",
		Severity: SeverityError,
		Description: "`aws ssm put-parameter` stores the value as-is under the given parameter " +
			"name; the whole point of `--type SecureString` is that the value is sensitive. " +
			"Passing the plaintext with `--value SECRET` (or `--value=SECRET`) puts the " +
			"secret in argv where `ps`, `/proc/PID/cmdline`, shell history, and AWS CLI " +
			"debug logs (`--debug`) can read it. Pipe the value in from stdin with `--cli-" +
			"input-json file://param.json` (mode 0600) or use `aws secretsmanager " +
			"create-secret --secret-string file://secret` which supports `file://` in every " +
			"code path.",
		Check: checkZC1812,
	})
}

func checkZC1812(node ast.Node) []Violation {
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
	if cmd.Arguments[0].String() != "ssm" {
		return nil
	}
	if cmd.Arguments[1].String() != "put-parameter" {
		return nil
	}

	hasSecureString := false
	hasInlineValue := false
	for i, arg := range cmd.Arguments[2:] {
		v := arg.String()
		switch {
		case v == "--type":
			if 2+i+1 < len(cmd.Arguments) && cmd.Arguments[2+i+1].String() == "SecureString" {
				hasSecureString = true
			}
		case v == "--type=SecureString":
			hasSecureString = true
		case v == "--value":
			if 2+i+1 < len(cmd.Arguments) {
				next := cmd.Arguments[2+i+1].String()
				if next != "" && !strings.HasPrefix(next, "file://") && !strings.HasPrefix(next, "-") {
					hasInlineValue = true
				}
			}
		case strings.HasPrefix(v, "--value="):
			val := strings.TrimPrefix(v, "--value=")
			if val != "" && !strings.HasPrefix(val, "file://") {
				hasInlineValue = true
			}
		}
	}
	if !hasSecureString || !hasInlineValue {
		return nil
	}
	return []Violation{{
		KataID: "ZC1812",
		Message: "`aws ssm put-parameter --type SecureString --value …` puts the " +
			"plaintext in argv — `ps` / `/proc/PID/cmdline` / history / CLI debug " +
			"logs can read it. Use `--cli-input-json file://…` (mode 0600) or the " +
			"`file://` form for `--value`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1813",
		Title:    "Warn on `cryptsetup luksFormat` / `reencrypt` — destructive LUKS header write",
		Severity: SeverityWarning,
		Description: "`cryptsetup luksFormat DEV` writes a new LUKS2 header at the start of DEV " +
			"and marks the remaining space as fresh ciphertext — any pre-existing filesystem " +
			"or LUKS metadata is gone. `cryptsetup reencrypt DEV` rewrites the entire device " +
			"in place, and an interruption mid-write leaves the volume partially re-encrypted " +
			"and dependent on the `--resume-only` recovery path. Pair `luksFormat` with " +
			"`--batch-mode` only after verifying DEV via `lsblk -o NAME,MODEL,SERIAL`, always " +
			"back up the header (`cryptsetup luksHeaderBackup`) before touching it, and run " +
			"`reencrypt` on an unmounted volume with UPS-backed power.",
		Check: checkZC1813,
	})
}

func checkZC1813(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cryptsetup" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "luksFormat", "reencrypt", "luks-format":
			return []Violation{{
				KataID: "ZC1813",
				Message: "`cryptsetup " + v + "` rewrites the LUKS header / device. " +
					"Verify the target (`lsblk`), back up with " +
					"`luksHeaderBackup`, and run on an unmounted volume with UPS.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1814",
		Title:    "Error on `dpkg --force-all` — enables every single `--force-*` option at once",
		Severity: SeverityError,
		Description: "`dpkg --force-all` is shorthand for ~18 distinct `--force-<option>` flags: " +
			"overwrite existing files, install unsigned packages, downgrade, install " +
			"depends-broken, remove essential, and more. The dpkg manual explicitly calls " +
			"this \"almost always a bad idea\". In provisioning scripts it hides the specific " +
			"constraint the author was trying to bypass, and when a later install re-triggers " +
			"the same state the underlying dependency conflict just re-surfaces on the next " +
			"unattended upgrade. Drop `--force-all` and spell out only the `--force-<option>` " +
			"you genuinely need, or fix the upstream conflict.",
		Check: checkZC1814,
	})
}

func checkZC1814(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `dpkg --force-all …` mangles to name=`force-all`.
	if ident.Value == "force-all" {
		return zc1814Hit(cmd)
	}

	if ident.Value != "dpkg" && ident.Value != "apt" && ident.Value != "apt-get" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--force-all" {
			return zc1814Hit(cmd)
		}
		if strings.Contains(v, "Dpkg::Options::=--force-all") {
			return zc1814Hit(cmd)
		}
	}
	return nil
}

func zc1814Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1814",
		Message: "`dpkg --force-all` enables every `--force-*` option at once — " +
			"overwrite, unsigned, downgrade, essential-removal, broken-deps. Drop it " +
			"and spell out only the specific `--force-<option>` you need, or fix the " +
			"upstream conflict.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

var zc1815NetUnits = map[string]bool{
	"NetworkManager":           true,
	"NetworkManager.service":   true,
	"systemd-networkd":         true,
	"systemd-networkd.service": true,
	"networking":               true,
	"networking.service":       true,
	"network":                  true,
	"network.service":          true,
	"wpa_supplicant":           true,
	"wpa_supplicant.service":   true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1815",
		Title:    "Warn on `systemctl restart NetworkManager` / `systemd-networkd` — drops the SSH session",
		Severity: SeverityWarning,
		Description: "Restarting the network manager from an SSH session tears down every active " +
			"connection the daemon supervises, including the one the script is running over. " +
			"The script freezes, the client sees a broken pipe, and recovery usually requires " +
			"console access. Route the change through `nmcli connection reload` + `nmcli " +
			"connection up <name>` (NetworkManager), `networkctl reload` (systemd-networkd), " +
			"or schedule the restart behind `systemd-run --on-active=30s` with a rollback " +
			"timer that re-enables the previous config if SSH does not reconnect.",
		Check: checkZC1815,
	})
}

func checkZC1815(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	action := cmd.Arguments[0].String()
	if action != "restart" && action != "stop" && action != "reload-or-restart" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		unit := strings.Trim(arg.String(), "\"'")
		if zc1815NetUnits[unit] {
			return []Violation{{
				KataID: "ZC1815",
				Message: "`systemctl " + action + " " + unit + "` drops every " +
					"connection the manager supervises — the SSH session freezes. " +
					"Use `nmcli connection reload` / `networkctl reload`, or a " +
					"`systemd-run --on-active=30s` rollback.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1816",
		Title:    "Warn on `docker/podman commit` — produces un-reproducible image, bakes in runtime state",
		Severity: SeverityWarning,
		Description: "`docker commit CONTAINER IMAGE` (and the podman / nerdctl equivalents) " +
			"snapshots a running container's filesystem into a new image. There is no " +
			"Dockerfile, so the build is not reproducible; the snapshot inherits whatever " +
			"`/tmp` scratch, shell history, environment variables, and — frequently — " +
			"credentials the container held at that moment; and the resulting image's layer " +
			"metadata records only the container id, not what was actually installed. Build " +
			"from a `Dockerfile` (or `docker buildx build`) so the image can be regenerated " +
			"from source, and use `docker commit` only for one-off rescue work on a local " +
			"image you are about to discard.",
		Check: checkZC1816,
	})
}

func checkZC1816(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" && ident.Value != "nerdctl" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "commit" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1816",
		Message: "`" + ident.Value + " commit` snapshots a running container — no " +
			"Dockerfile trail, runtime env / `/tmp` scratch / shell history get baked " +
			"in, and the layer metadata does not record what was installed. Build from " +
			"a `Dockerfile` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1817",
		Title:    "Warn on `git push --delete` / `git push -d` / `git push origin :branch` — remote branch removal",
		Severity: SeverityWarning,
		Description: "Deleting a branch on the remote is an irreversible server-side change the " +
			"local reflog cannot rescue. `git push --delete REMOTE BRANCH`, the short `-d`, " +
			"and the legacy `git push REMOTE :BRANCH` colon form all produce the same result: " +
			"the ref vanishes from the server, open pull requests are orphaned, CI runners " +
			"that pinned to the branch lose the target, and recovery needs the last commit " +
			"SHA which may only live in somebody else's local clone. Confirm the remote name, " +
			"check `git branch -r` / `gh pr list --head BRANCH` first, and prefer letting the " +
			"hosting platform delete the branch after a PR merge (with the auto-delete " +
			"setting) rather than scripting the push.",
		Check: checkZC1817,
	})
}

func checkZC1817(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "push" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--delete" || v == "-d" {
			return zc1817Hit(cmd, v)
		}
		// Legacy colon-form: `git push REMOTE :BRANCH` where :BRANCH starts with ":".
		if strings.HasPrefix(v, ":") && len(v) > 1 {
			return zc1817Hit(cmd, "origin "+v)
		}
	}
	return nil
}

func zc1817Hit(cmd *ast.SimpleCommand, flag string) []Violation {
	return []Violation{{
		KataID: "ZC1817",
		Message: "`git push " + flag + "` deletes the remote branch — open PRs are " +
			"orphaned, CI targets disappear, and the last commit SHA can only come " +
			"back from someone else's clone. Let the hosting platform auto-delete " +
			"after merge instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

var zc1818DeleteFlags = []string{
	"--delete",
	"--del",
	"--delete-before",
	"--delete-during",
	"--delete-delay",
	"--delete-after",
	"--delete-excluded",
	"--delete-missing-args",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1818",
		Title:    "Warn on `rsync --delete` without `--dry-run` — empty or wrong SRC wipes DST",
		Severity: SeverityWarning,
		Description: "`rsync --delete` (plus `--delete-before/-during/-after/-excluded`) removes " +
			"anything in DST that is not in SRC. If SRC is accidentally empty (typo in " +
			"path, unmounted mount point, wrong credentials pointing at an empty remote), " +
			"the destination loses every file that was there. The command has no undo. " +
			"Always preview the diff with `rsync -av --delete --dry-run SRC DST` first, " +
			"and cap the blast radius with `--max-delete=N` so the sync aborts if the plan " +
			"removes more files than expected.",
		Check: checkZC1818,
	})
}

func checkZC1818(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `rsync --delete SRC DST` mangles name to `delete` or similar.
	mangled := false
	switch ident.Value {
	case "delete", "del", "delete-before", "delete-during", "delete-delay",
		"delete-after", "delete-excluded", "delete-missing-args":
		mangled = true
	}

	if !mangled {
		if ident.Value != "rsync" {
			return nil
		}
		hasDelete := false
		for _, arg := range cmd.Arguments {
			v := arg.String()
			for _, flag := range zc1818DeleteFlags {
				if v == flag || strings.HasPrefix(v, flag+"=") {
					hasDelete = true
					break
				}
			}
			if hasDelete {
				break
			}
		}
		if !hasDelete {
			return nil
		}
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--dry-run" || v == "-n" || v == "--itemize-changes" {
			return nil
		}
		if strings.HasPrefix(v, "-") && !strings.HasPrefix(v, "--") && len(v) > 1 {
			for _, c := range v[1:] {
				if c == 'n' {
					return nil
				}
			}
		}
	}
	return []Violation{{
		KataID: "ZC1818",
		Message: "`rsync --delete` without `--dry-run` removes anything in DST that " +
			"isn't in SRC. Preview with `rsync -av --delete --dry-run SRC DST`, and " +
			"pin `--max-delete=N` so an accidentally empty SRC can't cascade.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1819",
		Title:    "Warn on `xattr -d com.apple.quarantine` / `xattr -cr` — removes macOS Gatekeeper quarantine",
		Severity: SeverityWarning,
		Description: "macOS sets the `com.apple.quarantine` extended attribute on every file " +
			"downloaded from the internet — Gatekeeper uses it to trigger the first-run " +
			"notarization / signature check. `xattr -d com.apple.quarantine FILE` strips the " +
			"attribute and lets the binary run with no prompt, and `xattr -cr DIR` does the " +
			"same recursively for every file in the tree. In a script that processes " +
			"downloaded artifacts this turns \"we vetted the binary\" into \"we trust whatever " +
			"landed in the download folder\". Verify the signature (`codesign --verify`) and " +
			"notarization (`spctl --assess --type execute`) first, or use " +
			"`xip`/`installer` packages so Gatekeeper stays in the loop.",
		Check: checkZC1819,
	})
}

func checkZC1819(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "xattr" {
		return nil
	}

	hasQuarantineDelete := false
	hasRecursiveClear := false

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-d" || v == "--delete" {
			if i+1 < len(cmd.Arguments) &&
				cmd.Arguments[i+1].String() == "com.apple.quarantine" {
				hasQuarantineDelete = true
			}
		}
		if strings.HasPrefix(v, "-") && !strings.HasPrefix(v, "--") && len(v) > 1 {
			hasC := false
			hasR := false
			for _, c := range v[1:] {
				if c == 'c' {
					hasC = true
				}
				if c == 'r' {
					hasR = true
				}
			}
			if hasC && hasR {
				hasRecursiveClear = true
			}
		}
	}

	if !hasQuarantineDelete && !hasRecursiveClear {
		return nil
	}
	return []Violation{{
		KataID: "ZC1819",
		Message: "`xattr -d com.apple.quarantine` / `-cr` strips the macOS Gatekeeper " +
			"quarantine — the binary runs with no signature / notarization check. " +
			"Verify with `codesign --verify` and `spctl --assess --type execute` " +
			"before stripping.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1820",
		Title:    "Warn on `netplan apply` — applies network config immediately with no rollback timer",
		Severity: SeverityWarning,
		Description: "`netplan apply` regenerates the rendered backend config (systemd-networkd " +
			"or NetworkManager) and brings it live right away. A mistake in the YAML — wrong " +
			"interface name, missing `dhcp4`, bad addresses, conflicting routes — drops the " +
			"admin SSH session, and recovery needs console access. Run `netplan try` first: " +
			"it applies the new config, waits for confirmation, and rolls back automatically " +
			"if no keypress arrives within the timeout. Only fall through to `netplan apply` " +
			"after the try window has elapsed successfully.",
		Check: checkZC1820,
	})
}

func checkZC1820(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "netplan" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	if cmd.Arguments[0].String() != "apply" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1820",
		Message: "`netplan apply` commits the YAML immediately — a mistake drops the " +
			"admin SSH session with no automatic rollback. Run `netplan try` first " +
			"(auto-reverts if no keypress within the timeout), then `netplan apply`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

var zc1821DiskutilDestructive = map[string]string{
	"eraseDisk":       "reformats the whole disk",
	"eraseVolume":     "reformats the volume",
	"secureErase":     "overwrites every block, no undo",
	"zeroDisk":        "writes zeros across the whole disk",
	"randomDisk":      "writes random bytes across the whole disk",
	"reformat":        "reformats the volume in place",
	"eraseCD":         "erases the optical disc",
	"erasePartitions": "removes every partition on the disk",
	"partitionDisk":   "rewrites the partition table",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1821",
		Title:    "Error on `diskutil eraseDisk` / `secureErase` / `partitionDisk` — macOS storage reformat",
		Severity: SeverityError,
		Description: "The `diskutil` subcommands `eraseDisk`, `eraseVolume`, `secureErase`, " +
			"`zeroDisk`, `randomDisk`, `reformat`, `erasePartitions`, and `partitionDisk` all " +
			"rewrite disk or volume state with no Time Machine snapshot or APFS " +
			"preservation. A wrong `/dev/diskN` (especially after a reboot that reordered " +
			"the BSD names) erases the wrong drive, and the only recovery is an offline " +
			"backup. Always pair the call with a typed confirmation, resolve the target by " +
			"`diskutil info -plist` / mount-point rather than by index, and run " +
			"`diskutil list` right before the destructive call.",
		Check: checkZC1821,
	})
}

func checkZC1821(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "diskutil" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	note, ok := zc1821DiskutilDestructive[sub]
	if !ok {
		return nil
	}
	return []Violation{{
		KataID: "ZC1821",
		Message: "`diskutil " + sub + "` " + note + ". Resolve the target by " +
			"`diskutil info -plist` / mount-point (not by index), run " +
			"`diskutil list` immediately before, and require a typed " +
			"confirmation.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1822",
		Title:    "Error on `csrutil disable` / `spctl --master-disable` — disables macOS system integrity / Gatekeeper",
		Severity: SeverityError,
		Description: "`csrutil disable` turns off System Integrity Protection: the kernel stops " +
			"blocking writes under `/System`, `/bin`, `/sbin`, runtime attachment to " +
			"protected processes becomes possible, and unsigned kexts can load. `spctl " +
			"--master-disable` (and `--global-disable`, `kext-consent disable`) removes " +
			"Gatekeeper / kext-consent enforcement, so any downloaded binary or kernel " +
			"extension runs without the user being prompted. Neither has a legitimate " +
			"provisioning use; both belong to ad-hoc developer workflows and are high-value " +
			"persistence steps for malware. Re-enable with `csrutil enable` in recovery mode " +
			"and `spctl --master-enable`.",
		Check: checkZC1822,
	})
}

func checkZC1822(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `spctl --master-disable` mangles to name=`master-disable`,
	// `spctl --global-disable` to `global-disable`.
	switch ident.Value {
	case "master-disable", "global-disable":
		return zc1822Hit(cmd, "spctl --"+ident.Value)
	}

	switch ident.Value {
	case "csrutil":
		if len(cmd.Arguments) > 0 && cmd.Arguments[0].String() == "disable" {
			return zc1822Hit(cmd, "csrutil disable")
		}
	case "spctl":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "--master-disable" || v == "--global-disable" {
				return zc1822Hit(cmd, "spctl "+v)
			}
		}
		if len(cmd.Arguments) >= 2 &&
			cmd.Arguments[0].String() == "kext-consent" &&
			cmd.Arguments[1].String() == "disable" {
			return zc1822Hit(cmd, "spctl kext-consent disable")
		}
	}
	return nil
}

func zc1822Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1822",
		Message: "`" + what + "` disables macOS SIP / Gatekeeper / kext-consent — " +
			"every malware analyst's favorite persistence primitive. Re-enable " +
			"(`csrutil enable` in recovery, `spctl --master-enable`) and keep " +
			"the default policy on.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1823",
		Title:    "Warn on `keytool -import -noprompt` — Java trust store imports without fingerprint check",
		Severity: SeverityWarning,
		Description: "`keytool -import -noprompt -trustcacerts -alias X -file CERT -keystore KS` " +
			"adds CERT to the Java trust store without showing its SHA-256 fingerprint or " +
			"asking the operator to confirm. If CERT came from an HTTP download, an attacker " +
			"wrote it in a shared temp dir, or a provisioning step fetched the wrong file, the " +
			"JVM will happily pin the attacker's CA as trusted and verify everything signed " +
			"against it. Drop `-noprompt`, or pre-verify with `keytool -printcert -file CERT` " +
			"and keep the alias+fingerprint pair in a versioned inventory before adding to any " +
			"trust store.",
		Check: checkZC1823,
	})
}

func checkZC1823(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "keytool" {
		return nil
	}

	hasImport := false
	hasNoPrompt := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-import" || v == "-importcert" || v == "-importkeystore" {
			hasImport = true
		}
		if v == "-noprompt" {
			hasNoPrompt = true
		}
	}
	if !hasImport || !hasNoPrompt {
		return nil
	}
	return []Violation{{
		KataID: "ZC1823",
		Message: "`keytool -import -noprompt` pins a cert to the Java trust store " +
			"without a fingerprint check. Drop `-noprompt`, verify with " +
			"`keytool -printcert -file CERT`, and store (alias, SHA-256) pairs in " +
			"an audited inventory.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1824",
		Title:    "Warn on `kubectl drain --disable-eviction` — bypasses PodDisruptionBudget via raw DELETE",
		Severity: SeverityWarning,
		Description: "`kubectl drain --disable-eviction` tells the client to delete pods directly " +
			"via the API instead of issuing Eviction requests. The Eviction pathway is what " +
			"honours PodDisruptionBudget — `--disable-eviction` drops pods regardless of the " +
			"minAvailable / maxUnavailable contract the workload owner defined. On a " +
			"multi-replica service this turns a rolling drain into a hard outage. Fix the " +
			"blocking PDB (raise minAvailable, wait for replicas to reschedule, or negotiate " +
			"with the owner) instead of flipping the flag off.",
		Check: checkZC1824,
	})
}

func checkZC1824(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubectl" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "drain" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--disable-eviction" || v == "--disable-eviction=true" {
			return []Violation{{
				KataID: "ZC1824",
				Message: "`kubectl drain --disable-eviction` deletes pods via raw API " +
					"DELETE — PodDisruptionBudgets are ignored and the workload " +
					"owner's availability contract is voided. Fix the blocking PDB " +
					"instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1825",
		Title:    "Warn on `scp -O` — forces legacy SCP wire protocol exposed to filename-injection CVEs",
		Severity: SeverityWarning,
		Description: "OpenSSH 9.0 switched `scp` to use the SFTP protocol by default — SFTP performs " +
			"structured file transfer instead of piping a remote shell, and closes the " +
			"filename-injection class that the old SCP wire protocol was vulnerable to " +
			"(CVE-2020-15778 and friends). `scp -O` forces the legacy SCP protocol, putting " +
			"the connection back on the old code path where a server (or a man-in-the-middle " +
			"in the remote host's shell) can inject shell metacharacters into filenames. If a " +
			"remote endpoint genuinely needs SCP, use `sftp` instead or upgrade the remote " +
			"server. Drop `-O` unless you have a named compatibility bug that requires it.",
		Check: checkZC1825,
	})
}

func checkZC1825(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "scp" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if arg.String() == "-O" {
			return []Violation{{
				KataID: "ZC1825",
				Message: "`scp -O` forces the legacy SCP wire protocol — the one exposed " +
					"to filename-injection (CVE-2020-15778 class). Drop `-O` (default " +
					"SFTP is safer), or use `sftp` / upgrade the remote server.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1826",
		Title:    "Warn on `install -m u+s` / `g+s` — symbolic setuid/setgid bit applied at install time",
		Severity: SeverityWarning,
		Description: "`install -m u+s SRC DEST` (or `g+s` / `ug+s` / `u=rwxs` etc.) applies the " +
			"setuid / setgid bit atomically at copy time — no intermediate `chmod` " +
			"step where a tripwire would fire, no time window where the file exists " +
			"without the special bit. Symbolic forms are easy to miss in review " +
			"because they don't carry the tell-tale leading `4`/`2`/`6` digit that " +
			"numeric-mode detection (see ZC1892) keys off. If DEST is on `$PATH`, " +
			"every local user can invoke the elevated binary. Install setuid / setgid " +
			"binaries only from trusted builds you have reviewed, and prefer narrow " +
			"capabilities (`setcap cap_net_bind_service+ep`) over broad setuid.",
		Check: checkZC1826,
	})
}

func checkZC1826(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "install" {
		return nil
	}
	for i, arg := range cmd.Arguments {
		v := arg.String()
		var mode string
		switch {
		case v == "-m" || v == "--mode":
			if i+1 < len(cmd.Arguments) {
				mode = cmd.Arguments[i+1].String()
			}
		case strings.HasPrefix(v, "-m") && len(v) > 2:
			mode = v[2:]
		case strings.HasPrefix(v, "--mode="):
			mode = strings.TrimPrefix(v, "--mode=")
		}
		if mode == "" {
			continue
		}
		mode = strings.Trim(strings.TrimSpace(mode), "\"'")
		// Numeric setuid / setgid is owned by ZC1892; this kata narrows to
		// symbolic-form setuid/setgid which the numeric scan does not catch.
		if zc1826IsNumericMode(mode) {
			continue
		}
		if zc1826HasSymbolicSetuid(mode) {
			return []Violation{{
				KataID: "ZC1826",
				Message: "`install -m " + mode + "` applies a symbolic setuid/setgid " +
					"bit — easy to miss in review. Use `0755` and grant narrow " +
					"caps with `setcap` instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1826IsNumericMode(mode string) bool {
	if mode == "" {
		return false
	}
	for _, r := range mode {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func zc1826HasSymbolicSetuid(mode string) bool {
	// chmod-style symbolic modes have an `s` or `t` in the perms section.
	// Examples flagged: `u+s`, `g+s`, `ug+s`, `u=rwxs`, `+s`.
	// `s` in the user or group perm slot means setuid / setgid.
	for _, chunk := range strings.Split(mode, ",") {
		if !strings.ContainsAny(chunk, "+=") {
			continue
		}
		if !strings.Contains(chunk, "s") {
			continue
		}
		// Only trip on who-selectors that can carry setuid/setgid (`u`, `g`,
		// `a`, or default/empty `+s` / `=s`) — `o+s` is a no-op.
		idx := strings.IndexAny(chunk, "+=")
		who := chunk[:idx]
		if who == "" || strings.ContainsAny(who, "uga") {
			return true
		}
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1827",
		Title:    "Error on `npm unpublish` — breaks every downstream that pinned the version",
		Severity: SeverityError,
		Description: "`npm unpublish PKG@VERSION` removes a published version from the registry. " +
			"Every downstream that pinned to that version — directly or through a transitive " +
			"lockfile entry — fails to install on the next `npm ci` / CI run. This is the " +
			"exact mechanism behind the 2016 `left-pad` outage; npm has since limited " +
			"unpublish to within 72 hours and added the `--force` gate, but within the " +
			"window the blast radius is still the whole ecosystem that pulled the package. " +
			"Use `npm deprecate PKG@VERSION 'reason'` instead — the version stays resolvable, " +
			"but installs print a warning and users can pin forward on their own schedule.",
		Check: checkZC1827,
	})
}

func checkZC1827(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "npm" && ident.Value != "pnpm" && ident.Value != "yarn" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "unpublish" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1827",
		Message: "`" + ident.Value + " unpublish` removes a published version — every " +
			"downstream that pinned it fails to install on next CI run (the left-pad " +
			"pattern). Use `npm deprecate PKG@VERSION 'reason'` so the version stays " +
			"resolvable with a warning.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1828",
		Title:    "Warn on `gcore PID` / `strace -p PID` — live ptrace attach dumps target memory",
		Severity: SeverityWarning,
		Description: "`gcore PID` writes a core dump of the running process to disk; `strace -p " +
			"PID` streams every syscall the process makes. Both attach via ptrace and expose " +
			"the target's memory, stack, environment variables, and argument buffers — " +
			"credentials, TLS session keys, and `$AWS_SECRET_ACCESS_KEY`-style env vars are " +
			"all readable. A root-run script that attaches to another user's process extracts " +
			"whatever that user has. Keep production scripts off ptrace; reach for " +
			"`coredumpctl` with a captured core or vendor-specific `perf` counters when you " +
			"only need syscall statistics.",
		Check: checkZC1828,
	})
}

func checkZC1828(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "gcore":
		if len(cmd.Arguments) > 0 && !zc1828IsHelp(cmd.Arguments) {
			return zc1828Hit(cmd, "gcore")
		}
	case "strace":
		for _, arg := range cmd.Arguments {
			if arg.String() == "-p" {
				return zc1828Hit(cmd, ident.Value+" -p")
			}
		}
	}
	return nil
}

func zc1828IsHelp(args []ast.Expression) bool {
	for _, a := range args {
		v := a.String()
		if v == "-h" || v == "--help" || v == "-?" || v == "--version" {
			return true
		}
	}
	return false
}

func zc1828Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1828",
		Message: "`" + what + " PID` attaches via ptrace — target memory, env, and " +
			"syscall args are exposed. Production scripts should not run ptrace; " +
			"use `coredumpctl` on a captured core instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1829",
		Title:    "Warn on `tailscale down` / `wg-quick down` / `nmcli con down` — drops the VPN that may carry the SSH session",
		Severity: SeverityWarning,
		Description: "A script that closes the VPN tunnel from within a remote session cuts " +
			"itself off whenever the admin SSH rides over that tunnel. `tailscale down`, " +
			"`wg-quick down WG0`, `openvpn` teardown, and `nmcli connection down NAME` all " +
			"tear the link down in place with no grace or rollback. Schedule the teardown " +
			"behind `systemd-run --on-active=30s --unit=recover <cmd to bring it back up>` " +
			"so the VPN is back before the unit expires, or run the command from the " +
			"host's console / out-of-band path rather than over the VPN itself.",
		Check: checkZC1829,
	})
}

func checkZC1829(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "tailscale":
		if len(cmd.Arguments) > 0 && cmd.Arguments[0].String() == "down" {
			return zc1829Hit(cmd, "tailscale down")
		}
	case "wg-quick":
		if len(cmd.Arguments) > 0 && cmd.Arguments[0].String() == "down" {
			return zc1829Hit(cmd, "wg-quick down")
		}
	case "nmcli":
		// `nmcli connection down <name>` / `nmcli con down <name>`.
		if len(cmd.Arguments) >= 2 {
			first := cmd.Arguments[0].String()
			if (first == "connection" || first == "con" || first == "c") &&
				cmd.Arguments[1].String() == "down" {
				return zc1829Hit(cmd, "nmcli connection down")
			}
		}
	}
	return nil
}

func zc1829Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1829",
		Message: "`" + what + "` tears down the VPN — if the SSH session rides on it, " +
			"the script cuts itself off with no rollback. Schedule recovery via " +
			"`systemd-run --on-active=30s`, or run from console / out-of-band.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1830",
		Title:    "Warn on `unsetopt NOMATCH` — unmatched glob becomes the literal pattern, silent bugs",
		Severity: SeverityWarning,
		Description: "`NOMATCH` is on by default in Zsh — an unmatched glob (`*.log` with no matching " +
			"files) errors out instead of silently passing through. Disabling it " +
			"(`unsetopt NOMATCH` or the equivalent `setopt NO_NOMATCH`) reverts to POSIX-sh " +
			"behaviour: the pattern is handed to the command verbatim, so `rm *.log` with no " +
			"matches runs `rm '*.log'` — which fails noisily for `rm` but, for commands that " +
			"accept arbitrary strings, silently processes the literal `*.log` instead of " +
			"files. Prefer scoped `*(N)` null-glob qualifier or `setopt LOCAL_OPTIONS; setopt " +
			"NULL_GLOB` inside a function, so the rest of the script keeps the default safety.",
		Check: checkZC1830,
	})
}

func checkZC1830(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1830IsNomatch(arg.String()) {
				return zc1830Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NONOMATCH" {
				return zc1830Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1830IsNomatch(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "NOMATCH"
}

func zc1830Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1830",
		Message: "`" + where + "` silences Zsh's unmatched-glob error — typos pass " +
			"through literally. Use `*(N)` per-glob or scope inside a function " +
			"with `setopt LOCAL_OPTIONS; setopt NULL_GLOB`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

var zc1831SshUnits = map[string]bool{
	"ssh":            true,
	"sshd":           true,
	"ssh.service":    true,
	"sshd.service":   true,
	"ssh.socket":     true,
	"sshd.socket":    true,
	"openssh-server": true,
	"openssh":        true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1831",
		Title:    "Error on `systemctl stop|disable|mask ssh/sshd` — locks out the next remote login",
		Severity: SeverityError,
		Description: "Stopping, disabling, or masking the SSH daemon closes the door on the next " +
			"remote login. Existing connections survive for a while because sshd's spawned " +
			"per-session process keeps running, but any reconnect / CI follow-up step that " +
			"needs to ssh back in gets `Connection refused`. `systemctl disable ssh` and " +
			"`systemctl mask ssh` also survive reboots. Recovery requires console or out-of-" +
			"band access. If the goal is config reload, use `systemctl reload sshd`; if the " +
			"host is being retired, make sshd the last service you touch.",
		Check: checkZC1831,
	})
}

func checkZC1831(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	action := cmd.Arguments[0].String()
	if action != "stop" && action != "disable" && action != "mask" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		unit := arg.String()
		if zc1831SshUnits[unit] {
			return []Violation{{
				KataID: "ZC1831",
				Message: "`systemctl " + action + " " + unit + "` blocks SSH — " +
					"existing sessions survive but reconnects fail. `disable`/`mask` " +
					"persist across reboots. Use `reload sshd` for config changes.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1832",
		Title:    "Warn on Zsh `limit coredumpsize unlimited` — setuid memory landing in core files",
		Severity: SeverityWarning,
		Description: "Zsh's `limit` builtin is the csh-style sibling of `ulimit`; `limit " +
			"coredumpsize unlimited` is the Zsh equivalent of `ulimit -c unlimited` and has " +
			"the same consequence: a crashing setuid or key-holding process leaves its " +
			"address space on disk as a world-readable core file. Leave the coredump " +
			"ceiling at the distro default (usually 0 for non-debug sessions), or use " +
			"`systemd-coredump` with restricted permissions when you need post-mortem data. " +
			"`ulimit -c unlimited` is covered by ZC1495; this kata catches the Zsh-specific " +
			"`limit`/`unlimit coredumpsize` spelling.",
		Check: checkZC1832,
	})
}

func checkZC1832(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "limit":
		// `limit coredumpsize unlimited` (with optional -h for hard limit).
		argIdx := 0
		if argIdx < len(cmd.Arguments) && cmd.Arguments[argIdx].String() == "-h" {
			argIdx++
		}
		if argIdx+1 >= len(cmd.Arguments) {
			return nil
		}
		resource := strings.ToLower(cmd.Arguments[argIdx].String())
		value := strings.ToLower(cmd.Arguments[argIdx+1].String())
		if (resource == "coredumpsize" || resource == "coredump") && value == "unlimited" {
			return zc1832Hit(cmd, "limit coredumpsize unlimited")
		}
	case "unlimit":
		for _, arg := range cmd.Arguments {
			v := strings.ToLower(arg.String())
			if v == "coredumpsize" || v == "coredump" {
				return zc1832Hit(cmd, "unlimit coredumpsize")
			}
		}
	}
	return nil
}

func zc1832Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1832",
		Message: "`" + where + "` enables unbounded core dumps (Zsh-specific `limit` " +
			"spelling of `ulimit -c unlimited`). A setuid crash drops its memory to " +
			"disk as a world-readable file — leave the ceiling at the distro default.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1833",
		Title:    "Warn on `unsetopt WARN_CREATE_GLOBAL` — silent accidental-global bugs inside functions",
		Severity: SeverityWarning,
		Description: "`WARN_CREATE_GLOBAL` makes Zsh warn when a function assigns to a name " +
			"that is not declared `local` / `typeset` in the current scope — the single " +
			"highest-value guardrail against the classic Bash-ism where a helper function " +
			"silently stomps on a caller's variable (`tmp=`, `i=`, `result=`). Disabling it " +
			"(`unsetopt WARN_CREATE_GLOBAL` or the equivalent `setopt NO_WARN_CREATE_GLOBAL`) " +
			"reverts to permissive behaviour: every unqualified assignment inside a function " +
			"escapes to global scope with no diagnostic. Leave the option on and fix the " +
			"offending function by adding `local` / `typeset` declarations, or — if you " +
			"really must silence it for a specific block — use `setopt LOCAL_OPTIONS; " +
			"unsetopt WARN_CREATE_GLOBAL` inside a function so the rest of the script keeps " +
			"the safety.",
		Check: checkZC1833,
	})
}

func checkZC1833(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1833IsWarnCreateGlobal(arg.String()) {
				return zc1833Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOWARNCREATEGLOBAL" {
				return zc1833Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1833IsWarnCreateGlobal(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "WARNCREATEGLOBAL"
}

func zc1833Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1833",
		Message: "`" + where + "` silences Zsh's warning for assignments leaking " +
			"out of function scope — classic caller-variable stomping. Declare " +
			"`local`/`typeset`; scope with `LOCAL_OPTIONS` if you must disable.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1834",
		Title:    "Error on `tc qdisc … root netem loss 100%` — hard blackhole on a live interface",
		Severity: SeverityError,
		Description: "`tc qdisc add/replace dev IFACE root netem loss 100%` (also `corrupt 100%` " +
			"or `duplicate 100%` with no buffering) installs a Linux kernel qdisc that " +
			"drops every outbound packet on the named interface. Running this on the " +
			"interface that carries your SSH session is the canonical way to lock " +
			"yourself out of a remote host — the `tc` command returns success, the kernel " +
			"happily applies the rule, and the next TCP segment ACK never arrives. Even on " +
			"the console it halts any process that depends on the interface. Stage netem " +
			"experiments on a secondary interface, wrap them in `at now + 5 minutes` (or a " +
			"`timeout … tc qdisc del …` recovery trap) so a partial failure does not leave " +
			"the link permanently black-holed.",
		Check: checkZC1834,
	})
}

func checkZC1834(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tc" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 3 {
		return nil
	}
	if args[0].String() != "qdisc" {
		return nil
	}
	action := args[1].String()
	if action != "add" && action != "replace" && action != "change" {
		return nil
	}
	// Walk remaining args looking for `netem <mode> 100%`.
	for i := 2; i+2 < len(args); i++ {
		if args[i].String() != "netem" {
			continue
		}
		for j := i + 1; j+1 < len(args); j++ {
			mode := args[j].String()
			if mode != "loss" && mode != "corrupt" && mode != "duplicate" {
				continue
			}
			val := args[j+1].String()
			if val == "100%" || val == "100" {
				return []Violation{{
					KataID: "ZC1834",
					Message: "`tc qdisc " + action + " … netem " + mode + " 100%` " +
						"black-holes every packet on the target interface — remote SSH " +
						"dies instantly. Stage on a secondary dev or wrap in a timed " +
						"recovery (`at now + N minutes … tc qdisc del …`).",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
	}
	return nil
}

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

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1836",
		Title:    "Error on `helm uninstall --no-hooks` — skips pre-delete cleanup, orphaned state",
		Severity: SeverityError,
		Description: "`helm uninstall RELEASE --no-hooks` (also spelled `helm delete --no-hooks` on " +
			"Helm v2 / `helm3 --no-hooks` interchangeably) tears down every chart-rendered " +
			"resource but silently skips the release's `pre-delete` and `post-delete` " +
			"Jobs / ConfigMap hooks. Those hooks are where production charts flush " +
			"write-ahead logs, deregister service-discovery entries, back up PVC content " +
			"before the PVC goes away, and release external locks — skipping them on a " +
			"live release is one of the classic ways to leave the cluster in a partially " +
			"deleted state with no way to replay the cleanup. Drop `--no-hooks` and let " +
			"the chart run as designed; if a hook is genuinely wedged, disable it at the " +
			"chart level with `helm.sh/hook-delete-policy: before-hook-creation,hook-succeeded`.",
		Check: checkZC1836,
	})
}

func checkZC1836(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "helm" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 2 {
		return nil
	}
	sub := args[0].String()
	if sub != "uninstall" && sub != "delete" {
		return nil
	}
	for _, arg := range args[1:] {
		if arg.String() == "--no-hooks" {
			return []Violation{{
				KataID: "ZC1836",
				Message: "`helm " + sub + " --no-hooks` skips pre/post-delete " +
					"cleanup hooks — orphaned locks, DNS, missed PVC backups. " +
					"Drop the flag; fix stuck hooks via " +
					"`helm.sh/hook-delete-policy`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1837",
		Title:    "Error on `chmod` granting non-owner access to `/dev/kvm` / `/dev/mem` / `/dev/kmem` / `/dev/port`",
		Severity: SeverityError,
		Description: "Distros ship `/dev/mem`, `/dev/kmem`, `/dev/port`, and `/dev/kvm` with tight " +
			"owner-only or group-only permissions managed by udev rules — these nodes hand " +
			"any process that can read or write them the keys to the kingdom (physical " +
			"memory, kernel memory, raw I/O ports, full hypervisor API). Flipping the mode " +
			"from a script (`chmod 666 /dev/kvm`, `chmod a+rw /dev/mem`) is a classic local " +
			"privilege-escalation vector dressed up as a convenience fix for a permission " +
			"error. Fix the actual problem: add the user to the `kvm` group, ship a proper " +
			"udev rule (`/etc/udev/rules.d/*.rules`), or grant the specific capability the " +
			"tool needs instead of blanket-chmod-ing the device.",
		Check: checkZC1837,
	})
}

var zc1837Devices = map[string]struct{}{
	"/dev/kvm":  {},
	"/dev/mem":  {},
	"/dev/kmem": {},
	"/dev/port": {},
}

func checkZC1837(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chmod" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 2 {
		return nil
	}

	var mode, target string
	for _, arg := range args {
		v := arg.String()
		if strings.HasPrefix(v, "-") {
			continue
		}
		if mode == "" {
			mode = v
			continue
		}
		target = v
		break
	}
	if target == "" {
		return nil
	}
	if _, hit := zc1837Devices[target]; !hit {
		return nil
	}
	if !zc1837GrantsNonOwner(mode) {
		return nil
	}
	return []Violation{{
		KataID: "ZC1837",
		Message: "`chmod " + mode + " " + target + "` grants non-owner access to " +
			"a privileged kernel device — classic local-privesc vector. Use " +
			"group membership or a udev rule instead of blanket chmod.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func zc1837GrantsNonOwner(mode string) bool {
	if mode == "" {
		return false
	}
	// Symbolic: look for tokens that grant to group/other/all.
	lower := strings.ToLower(mode)
	if strings.HasPrefix(lower, "+") {
		return true
	}
	for _, frag := range []string{"o+", "o=", "a+", "a=", "ugo+", "ugo="} {
		if strings.Contains(lower, frag) {
			return true
		}
	}
	// Numeric: chmod reads the mode as octal. Parser normalises leading-zero
	// octals to decimal (e.g. "0666" -> "438"), so branch on which one we got.
	for _, r := range mode {
		if r < '0' || r > '9' {
			return false
		}
	}
	var n int64
	if strings.ContainsAny(mode, "89") {
		n, _ = strconv.ParseInt(mode, 10, 32)
	} else {
		n, _ = strconv.ParseInt(mode, 8, 32)
	}
	// Flag only if any "other" (world) bit is set — these devices are managed
	// with group-only access (e.g. /dev/kvm = 660 kvm:root); tightening to
	// 660/600 is fine, opening to the world is the privesc case.
	return (n & 0o007) != 0
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1838",
		Title:    "Warn on `setopt GLOB_DOTS` — bare `*` silently starts matching hidden files",
		Severity: SeverityWarning,
		Description: "`GLOB_DOTS` off is the Zsh default: patterns like `*`, `*.log`, and " +
			"recursive `**/*` skip filenames that begin with a dot (`.git/`, `.env`, " +
			"`.ssh/`). Setting `setopt GLOB_DOTS` script-wide reverses that quietly — every " +
			"subsequent glob now also matches hidden entries, which turns routine " +
			"maintenance lines (`rm *`, `cp -r * /backup`, `chmod 644 *`) into " +
			"repository-wiping, secret-copying, permission-flipping bugs. Leave the option " +
			"alone at the script level and request dot-inclusion per-glob with the " +
			"Zsh-native `*(D)` qualifier (or `.* *` when you explicitly want both), so the " +
			"effect is scoped to the exact line that needs it.",
		Check: checkZC1838,
	})
}

func checkZC1838(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if zc1838IsGlobDots(v) {
				return zc1838Hit(cmd, "setopt "+v)
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOGLOBDOTS" {
				return zc1838Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1838IsGlobDots(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "GLOBDOTS"
}

func zc1838Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1838",
		Message: "`" + where + "` makes every bare `*` also match hidden files — " +
			"`rm *` quietly destroys `.git/`, `cp -r *` copies `.env`. Keep the " +
			"option alone; request dotfiles per-glob with `*(D)`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

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

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1840",
		Title:    "Error on `openssl enc -k PASSWORD` — legacy flag embeds secret in argv",
		Severity: SeverityError,
		Description: "`openssl enc -k PASSWORD` (the pre-OpenSSL-3 short form of `-pass " +
			"pass:PASSWORD`) takes the password directly as the next argv element — which " +
			"makes it visible to every `ps` reader, every `/proc/<pid>/cmdline` consumer, " +
			"shell history, and anything that logs command invocations. The same leak " +
			"applies to `openssl rsa`, `openssl pkcs12`, and other subcommands that still " +
			"accept the deprecated `-k` alias. Use `-pass env:VARNAME`, `-pass file:PATH`, " +
			"or `-pass fd:N` (read from an open descriptor) so the secret never rides in " +
			"the process argument vector.",
		Check: checkZC1840,
	})
}

func checkZC1840(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "openssl" {
		return nil
	}
	args := cmd.Arguments
	for i, arg := range args {
		if arg.String() != "-k" {
			continue
		}
		if i+1 >= len(args) {
			continue
		}
		val := args[i+1].String()
		// Ignore empty, flag-looking, or the `-k file:` / `-k env:` style
		// that newer openssl binaries tolerate.
		if val == "" || val[0] == '-' {
			continue
		}
		return []Violation{{
			KataID: "ZC1840",
			Message: "`openssl -k " + val + "` embeds the password in argv — visible " +
				"to `ps`, `/proc/<pid>/cmdline`, and shell history. Use " +
				"`-pass env:VAR`, `-pass file:PATH`, or `-pass fd:N`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}

var zc1841ProxyFlags = map[string]bool{"--proxy-insecure": true}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1841",
		Title:    "Error on `curl --proxy-insecure` — TLS verification disabled on the proxy hop",
		Severity: SeverityError,
		Description: "`curl --proxy-insecure` (alias of `-k` but scoped to the proxy leg, " +
			"introduced alongside `--proxy-cacert` in curl 7.52) tells curl to accept any " +
			"certificate presented by the HTTPS proxy that sits between the script and the " +
			"origin server. The origin TLS handshake is still validated, which makes the " +
			"issue easy to miss in review, but any box that can intercept traffic to the " +
			"proxy — a captive portal, a rogue WPAD auto-discovery, an attacker on the same " +
			"VLAN — can present its own cert and read or rewrite the tunnel contents, " +
			"including any `Authorization:` header attached to the request. Install the " +
			"proxy's CA bundle and point `--proxy-cacert` / `CURL_CA_BUNDLE` at it instead.",
		Check: checkZC1841,
	})
}

func checkZC1841(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "curl" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if arg.String() == "--proxy-insecure" {
			return zc1841Hit(cmd)
		}
	}
	return nil
}

func zc1841Hit(cmd *ast.SimpleCommand) []Violation {
	line, col := FlagArgPosition(cmd, zc1841ProxyFlags)
	return []Violation{{
		KataID: "ZC1841",
		Message: "`curl --proxy-insecure` skips TLS verification on the proxy hop — " +
			"an on-path attacker can present any cert and decrypt the tunnel " +
			"(including `Authorization:` headers). Install the proxy CA and use " +
			"`--proxy-cacert PATH`.",
		Line:   line,
		Column: col,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1842",
		Title:    "Warn on `setopt CDABLE_VARS` — `cd NAME` silently falls back to `cd $NAME`",
		Severity: SeverityWarning,
		Description: "With `CDABLE_VARS` on, any `cd NAME` whose `NAME` does not exist as a " +
			"directory is retried as `cd ${NAME}` — if a parameter of the same name is set, " +
			"the working directory silently jumps to wherever the variable points. A typo " +
			"like `cd cinfig` (intent: `config`) suddenly lands inside `${cinfig}` when one " +
			"exists, and every later relative path in the script is computed from the wrong " +
			"root. Keep this option inside `~/.zshrc` where it is an interactive shortcut; " +
			"in scripts, always `cd \"$dir\"` explicitly and pair with `|| exit` so a missed " +
			"directory fails loudly instead of rewriting `$PWD`.",
		Check: checkZC1842,
	})
}

func checkZC1842(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1842IsCdableVars(arg.String()) {
				return zc1842Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOCDABLEVARS" {
				return zc1842Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1842IsCdableVars(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "CDABLEVARS"
}

func zc1842Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1842",
		Message: "`" + where + "` turns a failed `cd NAME` into `cd $NAME` — a typo " +
			"silently lands in whatever directory the matching variable points to. " +
			"Keep this in `~/.zshrc`; in scripts use `cd \"$dir\" || exit`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1843",
		Title:    "Warn on `docker/podman run --cgroup-parent=/system.slice|/init.scope|/` — container escapes engine limits",
		Severity: SeverityWarning,
		Description: "`--cgroup-parent=PATH` places the container under the given cgroup parent, " +
			"which is normally `/docker` (or the engine's managed slice) and inherits the " +
			"engine-wide memory/CPU/IO caps. Pointing the flag at `/`, `/system.slice`, or " +
			"any host-managed slice puts the container side-by-side with systemd services — " +
			"the engine's defaults no longer apply, and a runaway container can starve " +
			"`sshd` or the kubelet for resources. Unless a specific orchestrator is " +
			"supplying a managed cgroup path, drop the flag and let the engine choose; if " +
			"you need custom limits, use `--memory` / `--cpus` / `--pids-limit` on the run " +
			"itself.",
		Check: checkZC1843,
	})
}

func checkZC1843(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 2 {
		return nil
	}
	sub := args[0].String()
	if sub != "run" && sub != "create" {
		return nil
	}

	for i := 1; i < len(args); i++ {
		v := args[i].String()
		var parent string
		switch {
		case strings.HasPrefix(v, "--cgroup-parent="):
			parent = strings.TrimPrefix(v, "--cgroup-parent=")
		case v == "--cgroup-parent" && i+1 < len(args):
			parent = args[i+1].String()
		default:
			continue
		}
		if zc1843IsHostSlice(parent) {
			return []Violation{{
				KataID: "ZC1843",
				Message: "`" + ident.Value + " " + sub + " --cgroup-parent=" + parent +
					"` puts the container under a host-managed slice — the engine's " +
					"memory/CPU caps no longer apply. Drop the flag or pass " +
					"`--memory`/`--cpus`/`--pids-limit` directly.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1843IsHostSlice(v string) bool {
	if v == "" {
		return false
	}
	trimmed := strings.Trim(v, "\"'")
	switch trimmed {
	case "/", "/system.slice", "/user.slice", "/init.scope", "/machine.slice":
		return true
	}
	// Anything under /system.slice or /init.scope qualifies too.
	for _, prefix := range []string{"/system.slice/", "/init.scope/", "/machine.slice/"} {
		if strings.HasPrefix(trimmed, prefix) {
			return true
		}
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1844",
		Title:    "Warn on `logger -p local0.info|local7.notice|…` — unreserved facility often uncollected",
		Severity: SeverityWarning,
		Description: "The eight `local0`–`local7` syslog facilities are reserved for site-specific " +
			"use. Most distro `rsyslog` and `systemd-journald` defaults do not route them " +
			"anywhere — they drop on the floor unless someone dropped a matching rule into " +
			"`/etc/rsyslog.d/*.conf`. Scripts that call `logger -p local0.info 'audit: user " +
			"added to wheel'` therefore log to nothing in the audit trail on a stock " +
			"machine. For portable audit-style logging use the POSIX-reserved `auth.notice` " +
			"or `authpriv.info` facility; for application events, pass `-t TAG` and use " +
			"`user.notice` (the default) or a dedicated journald unit.",
		Check: checkZC1844,
	})
}

func checkZC1844(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "logger" {
		return nil
	}
	args := cmd.Arguments
	for i, arg := range args {
		v := arg.String()
		var facPrio string
		switch {
		case v == "-p" || v == "--priority":
			if i+1 < len(args) {
				facPrio = args[i+1].String()
			}
		case strings.HasPrefix(v, "-p"):
			facPrio = strings.TrimPrefix(v, "-p")
		case strings.HasPrefix(v, "--priority="):
			facPrio = strings.TrimPrefix(v, "--priority=")
		}
		if facPrio == "" {
			continue
		}
		facility := facPrio
		if idx := strings.Index(facPrio, "."); idx >= 0 {
			facility = facPrio[:idx]
		}
		facility = strings.ToLower(strings.Trim(facility, "\"'"))
		if zc1844IsLocalFacility(facility) {
			return []Violation{{
				KataID: "ZC1844",
				Message: "`logger -p " + facPrio + "` writes to a `local*` facility — " +
					"stock `rsyslog`/`journald` rarely collects these. Use " +
					"`auth.notice`/`authpriv.info` for audit events, or " +
					"`user.notice` + `-t TAG` for app logs.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1844IsLocalFacility(f string) bool {
	if len(f) != len("local0") || !strings.HasPrefix(f, "local") {
		return false
	}
	c := f[len("local")]
	return c >= '0' && c <= '7'
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1845",
		Title:    "Warn on `setopt PATH_DIRS` — slash-bearing command names fall back to `$PATH` lookup",
		Severity: SeverityWarning,
		Description: "`PATH_DIRS` (off by default) changes how Zsh resolves a command that " +
			"contains a `/`: instead of treating `./foo/bar` or `subdir/cmd` as a direct " +
			"path, Zsh walks `$path` and retries `${path[i]}/subdir/cmd` until one is " +
			"executable. The surface intent — run a local binary — is silently replaced by " +
			"`/usr/local/bin/subdir/cmd` or any other same-shaped subtree that exists on " +
			"`$PATH`. This gets even worse on shared build hosts where `$PATH` contains " +
			"user-owned directories. Leave the option off and call local binaries with an " +
			"explicit leading `./`, or hand the full absolute path to the shell.",
		Check: checkZC1845,
	})
}

func checkZC1845(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if zc1845IsPathDirs(v) {
				return zc1845Hit(cmd, "setopt "+v)
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOPATHDIRS" {
				return zc1845Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1845IsPathDirs(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "PATHDIRS"
}

func zc1845Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1845",
		Message: "`" + where + "` lets `subdir/cmd` fall back to a `$PATH` lookup — " +
			"a missing local binary silently runs a same-named subtree elsewhere on " +
			"`$PATH`. Leave the option off; call locals as `./subdir/cmd`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1846",
		Title:    "Warn on `certbot … --force-renewal` — bypasses ACME rate-limit safety",
		Severity: SeverityWarning,
		Description: "`certbot renew --force-renewal` and `certbot certonly --force-renewal` reissue " +
			"a certificate regardless of remaining validity. Placed in a daily cron, the " +
			"same hostname burns through Let's Encrypt's per-domain rate limits (50 " +
			"certificates per registered domain per 7 days, 5 duplicate certificates per " +
			"domain per 7 days); once the limit trips, no cert for that host — fresh or " +
			"renewal — can be issued until the rolling window expires, which often happens " +
			"during an outage when you need it least. Drop `--force-renewal` and let " +
			"certbot's default 30-days-before-expiry gate do its job, or if you really need " +
			"a specific reissue, run it once manually.",
		Check: checkZC1846,
	})
}

func checkZC1846(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "certbot" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 2 {
		return nil
	}
	sub := args[0].String()
	if sub != "renew" && sub != "certonly" && sub != "run" {
		return nil
	}
	for _, arg := range args[1:] {
		if arg.String() == "--force-renewal" {
			return []Violation{{
				KataID: "ZC1846",
				Message: "`certbot " + sub + " --force-renewal` reissues regardless of " +
					"expiry — in a cron it burns Let's Encrypt rate limits (50 certs " +
					"per domain / 7 days). Drop the flag and let the 30-day gate work.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1847",
		Title:    "Warn on `setopt CHASE_LINKS` — every `cd` silently swaps symlink paths for the real inode",
		Severity: SeverityWarning,
		Description: "`CHASE_LINKS` off is the Zsh default: `cd releases/current` leaves `$PWD` " +
			"as the logical path the user typed, and `cd ..` steps back up through the " +
			"symlink to where they came from. Turning the option on globally makes every " +
			"`cd` resolve the target to its physical inode — so `cd releases/current` lands " +
			"in `/srv/app/releases/20260415-deadbeef`, and the next `cd ../config` looks " +
			"for `/srv/app/releases/config` instead of the `/srv/app/config` that the user " +
			"expected. Scripts that rely on blue/green-style `current` symlinks break " +
			"silently. Keep the option off at the script level and request one-shot " +
			"physical resolution with `cd -P target` when a specific call needs it.",
		Check: checkZC1847,
	})
}

func checkZC1847(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if zc1847IsChaseLinks(v) {
				return zc1847Hit(cmd, "setopt "+v)
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOCHASELINKS" {
				return zc1847Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1847IsChaseLinks(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "CHASELINKS"
}

func zc1847Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1847",
		Message: "`" + where + "` makes every `cd` resolve symlinks to the physical " +
			"inode — `cd releases/current` lands in the release dir, breaking `..` " +
			"navigation. Keep it off; use `cd -P target` one-shot when needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1848",
		Title:    "Warn on `ssh -o CheckHostIP=no` — DNS-spoof warning for known hosts silenced",
		Severity: SeverityWarning,
		Description: "`CheckHostIP` (on by default) stores the host's IP address alongside its " +
			"host key in `~/.ssh/known_hosts`; if DNS later resolves the same name to a " +
			"different IP but the key still matches, ssh warns you. Turning the check off " +
			"with `-o CheckHostIP=no` keeps the host-key comparison but silences the " +
			"IP-mismatch warning — which means a DNS-poisoning attacker who already holds " +
			"the previously-seen host key (stolen, misplaced backup, leaked by a " +
			"decommissioned box) can route the session through their box without a peep. " +
			"Leave the default, and if you really need to skip the IP record (load-balanced " +
			"pool with shared keys) document the risk and prefer `HostKeyAlias` instead.",
		Check: checkZC1848,
	})
}

func checkZC1848(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh" && ident.Value != "scp" && ident.Value != "sftp" {
		return nil
	}
	args := cmd.Arguments
	for i, arg := range args {
		v := arg.String()
		var kv string
		switch {
		case v == "-o" && i+1 < len(args):
			kv = args[i+1].String()
		case strings.HasPrefix(v, "-o"):
			kv = strings.TrimPrefix(v, "-o")
		default:
			continue
		}
		if zc1848IsCheckHostIPNo(kv) {
			return []Violation{{
				KataID: "ZC1848",
				Message: "`" + ident.Value + " -o CheckHostIP=no` silences the " +
					"IP-mismatch warning for known hosts — a DNS-spoof + leaked " +
					"host-key attack goes undetected. Leave the default, or use " +
					"`HostKeyAlias` for load-balanced pools.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1848IsCheckHostIPNo(kv string) bool {
	norm := strings.ToLower(strings.Trim(kv, "\"' \t"))
	for _, frag := range []string{"checkhostip=no", "checkhostip = no", "checkhostip=false", "checkhostip=0", "checkhostip=off"} {
		if norm == frag {
			return true
		}
	}
	// Tolerate stray spaces around `=`.
	if strings.HasPrefix(norm, "checkhostip") {
		rest := strings.TrimPrefix(norm, "checkhostip")
		rest = strings.TrimSpace(rest)
		if strings.HasPrefix(rest, "=") {
			val := strings.TrimSpace(strings.TrimPrefix(rest, "="))
			return val == "no" || val == "false" || val == "0" || val == "off"
		}
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1849",
		Title:    "Warn on `setopt ALL_EXPORT` — every later `var=value` silently becomes `export var=value`",
		Severity: SeverityWarning,
		Description: "`ALL_EXPORT` (POSIX `set -a` equivalent, off by default) tells Zsh to mark " +
			"every parameter assignment for export as soon as it is created, so " +
			"`password=$(cat secret)` immediately rides into the environment of every " +
			"child process the script spawns — the `ps e`, `/proc/<pid>/environ`, and " +
			"journal of any later `| tee`, `| mail`, or `logger` call. Enabling it " +
			"script-wide to avoid a few `export` keywords leaks credentials and private " +
			"config by default. Drop the `setopt`, scope exports explicitly with " +
			"`export VAR=value`, or wrap a narrow section in `setopt LOCAL_OPTIONS; setopt " +
			"ALL_EXPORT` inside a function so the effect cannot leak past the closing " +
			"brace.",
		Check: checkZC1849,
	})
}

func checkZC1849(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if zc1849IsAllExport(v) {
				return zc1849Hit(cmd, "setopt "+v)
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOALLEXPORT" {
				return zc1849Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1849IsAllExport(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "ALLEXPORT"
}

func zc1849Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1849",
		Message: "`" + where + "` marks every later assignment for export — secrets " +
			"like `password=...` leak into every child's env. Drop it; use " +
			"explicit `export`, or scope inside a `LOCAL_OPTIONS` function.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1850",
		Title:    "Warn on `ssh -o LogLevel=QUIET` — silences security-relevant ssh diagnostics",
		Severity: SeverityWarning,
		Description: "`LogLevel=QUIET` (aliased to the `-q` short flag) suppresses every " +
			"informational or warning message ssh would otherwise print: host-key " +
			"changes, key-exchange downgrades, agent-forwarding permission denials, " +
			"canonical-hostname rewrites. In a script, that means the output looks clean " +
			"even when ssh is shouting about a MITM on the other end. Keep the default " +
			"`INFO` level (or raise to `VERBOSE` during debugging), capture stderr to a " +
			"log if the noise bothers you, and never pair `LogLevel=QUIET` with " +
			"`StrictHostKeyChecking=no` in the same call — that combination actively " +
			"hides known-bad-key events.",
		Check: checkZC1850,
	})
}

func checkZC1850(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh" && ident.Value != "scp" && ident.Value != "sftp" {
		return nil
	}
	args := cmd.Arguments
	for i, arg := range args {
		v := arg.String()
		var kv string
		switch {
		case v == "-o" && i+1 < len(args):
			kv = args[i+1].String()
		case strings.HasPrefix(v, "-o"):
			kv = strings.TrimPrefix(v, "-o")
		default:
			continue
		}
		if zc1850IsLogLevelQuiet(kv) {
			return []Violation{{
				KataID: "ZC1850",
				Message: "`" + ident.Value + " -o LogLevel=QUIET` silences host-key, " +
					"agent-forward, and canonical-hostname warnings — a MITM " +
					"event produces no stderr. Keep the default level; capture " +
					"stderr to a log if you need it clean.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1850IsLogLevelQuiet(kv string) bool {
	norm := strings.ToLower(strings.Trim(kv, "\"' \t"))
	if !strings.HasPrefix(norm, "loglevel") {
		return false
	}
	rest := strings.TrimSpace(strings.TrimPrefix(norm, "loglevel"))
	if !strings.HasPrefix(rest, "=") {
		return false
	}
	val := strings.TrimSpace(strings.TrimPrefix(rest, "="))
	return val == "quiet" || val == "fatal" || val == "error"
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1851",
		Title:    "Warn on `unsetopt FUNCTION_ARGZERO` — `$0` inside a function stops reporting the function name",
		Severity: SeverityWarning,
		Description: "`FUNCTION_ARGZERO` is Zsh's default: inside a function or `source`d file, " +
			"`$0` holds the function/file name, which is what every `log_error \"$0: ...\"` " +
			"helper, every self-reflecting `$funcfiletrace` fallback, and every `case $0` " +
			"dispatcher expects. Turning it off reverts to POSIX-sh behaviour where `$0` " +
			"always points at the outer script — so `my_func() { echo \"${0}: bad input\" }` " +
			"silently starts logging `myscript.sh: bad input` for every function, which " +
			"makes stack-trace logs unreadable and breaks dispatchers that branch on `$0`. " +
			"Keep the option on at the script level and, if one specific helper needs the " +
			"POSIX name, reach it explicitly with `$ZSH_ARGZERO` or `$ZSH_SCRIPT`.",
		Check: checkZC1851,
	})
}

func checkZC1851(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1851IsFunctionArgzero(arg.String()) {
				return zc1851Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOFUNCTIONARGZERO" {
				return zc1851Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1851IsFunctionArgzero(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "FUNCTIONARGZERO"
}

func zc1851Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1851",
		Message: "`" + where + "` makes `$0` inside functions point at the outer " +
			"script — breaks `log \"$0: ...\"` helpers and `case $0` dispatchers. " +
			"Keep the option on; reach the script name explicitly via " +
			"`$ZSH_ARGZERO` when needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

var zc1852PanicFlags = map[string]bool{"--panic-on": true}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1852",
		Title:    "Error on `firewall-cmd --panic-on` — firewalld drops every packet, kills the SSH session",
		Severity: SeverityError,
		Description: "`firewall-cmd --panic-on` puts firewalld into panic mode, which drops every " +
			"inbound and outbound packet regardless of zone or rule. Running this over a " +
			"remote SSH session is the textbook way to lock yourself out: the command " +
			"returns success, the TCP ACK for that reply never arrives, and nobody can " +
			"reach the host until someone visits the console to `--panic-off`. Stage " +
			"panic-mode experiments on a machine you can power-cycle, gate the call behind " +
			"`at now + 5 minutes` with an auto-disable, or use targeted zone rules instead " +
			"of the blanket switch.",
		Check: checkZC1852,
	})
}

func checkZC1852(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "firewall-cmd" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if arg.String() == "--panic-on" {
			return zc1852Hit(cmd)
		}
	}
	return nil
}

func zc1852Hit(cmd *ast.SimpleCommand) []Violation {
	line, col := FlagArgPosition(cmd, zc1852PanicFlags)
	return []Violation{{
		KataID: "ZC1852",
		Message: "`firewall-cmd --panic-on` drops every packet regardless of zone — " +
			"an SSH-run call loses the session instantly. Use targeted zone rules; " +
			"if you really need panic mode, gate behind `at now + N minutes … " +
			"firewall-cmd --panic-off`.",
		Line:   line,
		Column: col,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1853",
		Title:    "Warn on `setopt MARK_DIRS` — glob-matched directories gain a silent trailing `/`",
		Severity: SeverityWarning,
		Description: "With `MARK_DIRS` on, every filename produced by a glob that resolves to a " +
			"directory picks up a trailing `/`. Inside a shell it looks harmless, but " +
			"scripts that pass the glob result to other tools break in quiet ways: " +
			"`[[ -f \"$f\" ]]` rejects `dir/` because it is not a regular file, `rm -f *` " +
			"sees `dir/` and silently skips it (GNU rm refuses to remove directories " +
			"without `-r`), and downstream hash maps indexed on basenames suddenly carry " +
			"two keys for what the user thinks is one entry. Keep the option off at the " +
			"script level and request the trailing slash per-glob with the `(/)` qualifier " +
			"(`dirs=( *(/) )`) when you really need directories only.",
		Check: checkZC1853,
	})
}

func checkZC1853(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1853IsMarkDirs(arg.String()) {
				return zc1853Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOMARKDIRS" {
				return zc1853Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1853IsMarkDirs(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "MARKDIRS"
}

func zc1853Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1853",
		Message: "`" + where + "` appends a trailing `/` to every glob-matched " +
			"directory — `[[ -f \"$f\" ]]` and `rm -f *` start skipping, hash " +
			"maps keyed on basenames double up. Keep the option off; use " +
			"`*(/)` when you need dirs only.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1854",
		Title:    "Error on `yum-config-manager --add-repo http://…` / `zypper addrepo http://…` — plaintext repo allows MITM",
		Severity: SeverityError,
		Description: "Adding a package repository over plain HTTP (`yum-config-manager " +
			"--add-repo http://…`, `dnf config-manager --add-repo http://…`, `zypper " +
			"addrepo http://…`) tells the package manager to fetch metadata and RPMs " +
			"without TLS — any on-path attacker can substitute packages, and even GPG " +
			"signature checks do not help because the attacker can simply strip the " +
			"`repo_gpgcheck=1` line from the unsigned `.repo` file. Use the `https://` " +
			"mirror (every major distro now publishes one), or pin to a local mirror over " +
			"TLS and drop a `gpgkey=file:///etc/pki/...` entry in the same `.repo` so " +
			"signatures cannot be disabled mid-install.",
		Check: checkZC1854,
	})
}

var zc1854Flags = map[string]bool{
	"--add-repo": true,
	"addrepo":    true,
	"ar":         true,
}

func checkZC1854(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "yum-config-manager":
		for _, arg := range cmd.Arguments {
			if zc1854IsHTTPURL(arg.String()) {
				return zc1854Hit(cmd, "yum-config-manager --add-repo "+arg.String())
			}
		}
	case "dnf":
		if len(cmd.Arguments) >= 3 &&
			cmd.Arguments[0].String() == "config-manager" &&
			cmd.Arguments[1].String() == "--add-repo" {
			for _, arg := range cmd.Arguments[2:] {
				if zc1854IsHTTPURL(arg.String()) {
					return zc1854Hit(cmd, "dnf config-manager --add-repo "+arg.String())
				}
			}
		}
	case "zypper":
		if len(cmd.Arguments) >= 2 {
			sub := cmd.Arguments[0].String()
			if sub == "addrepo" || sub == "ar" {
				for _, arg := range cmd.Arguments[1:] {
					if zc1854IsHTTPURL(arg.String()) {
						return zc1854Hit(cmd, "zypper addrepo "+arg.String())
					}
				}
			}
		}
	}
	return nil
}

func zc1854IsHTTPURL(v string) bool {
	return strings.HasPrefix(v, "http://")
}

func zc1854Hit(cmd *ast.SimpleCommand, where string) []Violation {
	line, col := FlagArgPosition(cmd, zc1854Flags)
	return []Violation{{
		KataID: "ZC1854",
		Message: "`" + where + "` registers a plaintext repo — on-path attacker can " +
			"substitute packages and strip GPG-check directives. Use `https://` and " +
			"pin `gpgkey=file://` in the `.repo`.",
		Line:   line,
		Column: col,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1855",
		Title:    "Avoid `$GROUPS` — Bash-only array; Zsh exposes supplementary groups as `$groups`",
		Severity: SeverityWarning,
		Description: "`$GROUPS` is a Bash magic parameter that holds the caller's supplementary " +
			"GIDs as a numeric array. Zsh does not populate `$GROUPS`; it has " +
			"`$groups`, a lowercase associative array keyed by group *name* with the GID " +
			"as value (`${(k)groups}` for names, `${(v)groups}` for IDs). Scripts ported " +
			"from Bash that iterate `${GROUPS[@]}` therefore see an empty list under " +
			"Zsh and silently skip group-membership checks. Use `${(k)groups}` for names " +
			"or `${(v)groups}` for numeric GIDs; the Zsh `id -Gn` fallback keeps the " +
			"script portable across shells.",
		Check: checkZC1855,
	})
}

func checkZC1855(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if cmd.Name == nil {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if zc1855RefersToGROUPS(arg.String()) {
			return []Violation{{
				KataID: "ZC1855",
				Message: "`$GROUPS` is a Bash-only array — Zsh populates `$groups` " +
					"(associative name→GID) instead. Iterate `${(k)groups}` for " +
					"names or `${(v)groups}` for GIDs, or fall back to `id -Gn`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1855RefersToGROUPS(v string) bool {
	// Walk the arg looking for `$GROUPS` or `${GROUPS...}` as a distinct
	// token. Accept trailing `[`, `}`, or end-of-string so callers like
	// `${GROUPS[@]}` still match but `$GROUPSIZE` does not.
	i := 0
	for {
		idx := strings.Index(v[i:], "GROUPS")
		if idx < 0 {
			return false
		}
		idx += i
		// Require `$` or `${` immediately before.
		prefixOK := false
		switch {
		case idx >= 2 && v[idx-2:idx] == "${":
			prefixOK = true
		case idx >= 1 && v[idx-1] == '$':
			prefixOK = true
		}
		if prefixOK {
			end := idx + len("GROUPS")
			if end == len(v) {
				return true
			}
			next := v[end]
			if next == '[' || next == '}' || next == '"' || next == ' ' || next == '\t' {
				return true
			}
		}
		i = idx + 1
		if i >= len(v) {
			return false
		}
	}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1856",
		Title:    "Warn on `unset arr[N]` — Zsh does not delete the array element, the array keeps its length",
		Severity: SeverityWarning,
		Description: "In Bash, `unset arr[N]` removes the N-th element of the array (leaving a " +
			"sparse hole). In Zsh the same invocation passes the literal string `arr[N]` " +
			"to the `unset` builtin, which looks for a parameter with that name — finds " +
			"nothing — and returns success. The array is left untouched, `${#arr[@]}` " +
			"does not budge, and every downstream `for x in \"${arr[@]}\"` keeps iterating " +
			"the element the script thought it had removed. Use Zsh's native assignment " +
			"form `arr[N]=()` to delete an index, or `arr=(\"${(@)arr:#pattern}\")` to " +
			"filter by value.",
		Check: checkZC1856,
	})
}

func checkZC1856(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "unset" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1856IsArraySubscript(v) {
			return []Violation{{
				KataID: "ZC1856",
				Message: "`unset " + v + "` is a Bash idiom — in Zsh it tries to " +
					"unset a parameter literally named `" + v + "` and leaves the " +
					"array untouched. Use `arr[N]=()` or rebuild with " +
					"`arr=(\"${(@)arr:#pattern}\")`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1856IsArraySubscript(v string) bool {
	// Strip surrounding quotes and parser-applied `(...)` wrapping.
	trimmed := strings.TrimSpace(v)
	trimmed = strings.Trim(trimmed, "\"'")
	if strings.HasPrefix(trimmed, "(") && strings.HasSuffix(trimmed, ")") && len(trimmed) >= 2 {
		trimmed = trimmed[1 : len(trimmed)-1]
	}
	open := strings.Index(trimmed, "[")
	close := strings.LastIndex(trimmed, "]")
	if open <= 0 || close <= open+1 {
		return false
	}
	// The name portion must look like a shell identifier.
	return zc1856IsIdentifier(trimmed[:open])
}

func zc1856IsIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 && !zc1856IsIdentStart(r) {
			return false
		}
		if i > 0 && !zc1856IsIdentTail(r) {
			return false
		}
	}
	return true
}

func zc1856IsIdentStart(r rune) bool {
	return r == '_' || isAsciiLetter(r)
}

func zc1856IsIdentTail(r rune) bool {
	return r == '_' || isAsciiLetter(r) || (r >= '0' && r <= '9')
}

func isAsciiLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1857",
		Title:    "Error on `cloud-init clean` — wipes boot state, next reboot re-provisions the host",
		Severity: SeverityError,
		Description: "`cloud-init clean` (and variants `--logs`, `--reboot`, `--machine-id`) " +
			"removes every marker under `/var/lib/cloud/` and `/var/log/cloud-init*`, " +
			"which tells cloud-init to re-run from scratch on the next boot. That run " +
			"re-imports the image-builder's user-data: regenerates SSH host keys, resets " +
			"the hostname, replaces `/etc/fstab` entries the operator may have edited, " +
			"and (with `--reboot`) triggers the replay immediately. In a maintenance " +
			"script this silently erases everything the operator configured after " +
			"first-boot. Keep the command out of automation; if you truly need to " +
			"re-seed an instance, snapshot state first and run the command interactively.",
		Check: checkZC1857,
	})
}

func checkZC1857(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cloud-init" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	if cmd.Arguments[0].String() != "clean" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1857",
		Message: "`cloud-init clean` wipes `/var/lib/cloud/` boot state — the next " +
			"reboot re-runs the user-data and overwrites operator changes " +
			"(SSH host keys, hostname, `/etc/fstab`). Run interactively only.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1858",
		Title:    "Error on `ssh -c 3des-cbc|arcfour|blowfish-cbc` — weak cipher forced on the tunnel",
		Severity: SeverityError,
		Description: "OpenSSH disables legacy ciphers by default; a script that explicitly forces " +
			"one with `-c 3des-cbc`, `-c arcfour`, `-c blowfish-cbc`, or a matching entry " +
			"in `-o Ciphers=...` downgrades the tunnel to an algorithm with known plaintext " +
			"recovery, IV-reuse, or birthday-bound attacks. Typically this is done to reach " +
			"an old appliance — but it drags every other session on the same invocation " +
			"down with it. Leave cipher selection to OpenSSH's default; if a legacy device " +
			"absolutely requires a weak cipher, isolate it in a `Host ...` block in " +
			"`~/.ssh/config` with explicit `HostKeyAlgorithms` and keep the rest of the " +
			"fleet on strong defaults.",
		Check: checkZC1858,
	})
}

var zc1858Weak = []string{
	"3des-cbc",
	"arcfour",
	"arcfour128",
	"arcfour256",
	"blowfish-cbc",
	"cast128-cbc",
	"des-cbc",
	"rijndael-cbc@lysator.liu.se",
}

func checkZC1858(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh" && ident.Value != "scp" && ident.Value != "sftp" {
		return nil
	}
	args := cmd.Arguments
	for i, arg := range args {
		v := arg.String()
		var candidate string
		switch {
		case v == "-c" && i+1 < len(args):
			candidate = args[i+1].String()
		case v == "-o" && i+1 < len(args):
			candidate = zc1858ExtractCiphers(args[i+1].String())
		case strings.HasPrefix(v, "-o"):
			candidate = zc1858ExtractCiphers(strings.TrimPrefix(v, "-o"))
		}
		if candidate == "" {
			continue
		}
		if weak := zc1858FirstWeakCipher(candidate); weak != "" {
			return []Violation{{
				KataID: "ZC1858",
				Message: "`" + ident.Value + " ... " + weak + "` forces a weak cipher " +
					"with known plaintext-recovery / IV-reuse attacks. Leave " +
					"cipher selection to OpenSSH defaults; if a legacy peer needs " +
					"it, scope inside a `Host` block in `~/.ssh/config`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1858ExtractCiphers(kv string) string {
	kv = strings.TrimSpace(strings.Trim(kv, "\"'"))
	lower := strings.ToLower(kv)
	if !strings.HasPrefix(lower, "ciphers") {
		return ""
	}
	rest := strings.TrimSpace(strings.TrimPrefix(lower, "ciphers"))
	if !strings.HasPrefix(rest, "=") {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(rest, "="))
}

func zc1858FirstWeakCipher(list string) string {
	lower := strings.ToLower(list)
	for _, entry := range strings.Split(lower, ",") {
		entry = strings.TrimSpace(strings.Trim(entry, "+^-"))
		for _, weak := range zc1858Weak {
			if entry == weak {
				return entry
			}
		}
	}
	return ""
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1859",
		Title:    "Warn on `unsetopt MULTIOS` — `cmd >a >b` silently keeps only the last redirection",
		Severity: SeverityWarning,
		Description: "`MULTIOS` is on by default in Zsh: `cmd >out.log >>archive.log` sends stdout " +
			"to both files via an implicit `tee`, and `cmd <a <b` concatenates the two " +
			"inputs in order. Disabling it reverts to POSIX-sh semantics — Zsh opens each " +
			"earlier redirection, closes it immediately, and only the last one in the " +
			"direction wins. Any script that was written for Zsh suddenly starts dropping " +
			"the `archive.log` tail, and log collectors that opened `archive.log` keep " +
			"the fd but never receive new lines. Keep the option on at the script level; " +
			"if one specific line really needs POSIX behaviour, wrap it in a function with " +
			"`setopt LOCAL_OPTIONS; unsetopt MULTIOS`.",
		Check: checkZC1859,
	})
}

func checkZC1859(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1859IsMultios(arg.String()) {
				return zc1859Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOMULTIOS" {
				return zc1859Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1859IsMultios(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "MULTIOS"
}

func zc1859Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1859",
		Message: "`" + where + "` reverts to POSIX single-output redirection — " +
			"`cmd >a >b` silently drops `a`, log collectors stop receiving new " +
			"lines. Keep the option on; scope inside a `LOCAL_OPTIONS` function " +
			"if one line really needs POSIX.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1860",
		Title:    "Warn on `hostnamectl set-hostname NEW` — caches and certs still reference the old name",
		Severity: SeverityWarning,
		Description: "`hostnamectl set-hostname NEW` (and the new-style `hostnamectl hostname NEW` " +
			"and `hostname NEW`) updates `/etc/hostname` and `kernel.hostname` atomically, " +
			"but every process that called `gethostname()` at startup keeps the old " +
			"value until it restarts: syslog tags, Prometheus scrape labels, Docker " +
			"daemons, and anything that populated a TLS `subjectAltName` with `$(hostname)` " +
			"still speak as the previous host. Change the hostname interactively, then " +
			"plan a restart window — in automation, prefer shipping the new hostname via " +
			"cloud-init / Ignition so every service starts with it from boot.",
		Check: checkZC1860,
	})
}

func checkZC1860(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value == "hostnamectl" {
		if len(cmd.Arguments) < 2 {
			return nil
		}
		sub := cmd.Arguments[0].String()
		if sub != "set-hostname" && sub != "hostname" {
			return nil
		}
		return zc1860Hit(cmd, "hostnamectl "+sub+" "+cmd.Arguments[1].String())
	}
	if ident.Value == "hostname" {
		if len(cmd.Arguments) != 1 {
			return nil
		}
		v := cmd.Arguments[0].String()
		// `hostname` with no args just prints; `hostname -f` is read-only.
		if len(v) == 0 || v[0] == '-' {
			return nil
		}
		return zc1860Hit(cmd, "hostname "+v)
	}
	return nil
}

func zc1860Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1860",
		Message: "`" + where + "` updates the kernel hostname live, but running " +
			"services keep the old `gethostname()` — syslog tags, Prometheus " +
			"labels, TLS SANs stay stale. Apply at provisioning or reboot.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1861",
		Title:    "Warn on `setopt OCTAL_ZEROES` — leading-zero integers silently reinterpret as octal",
		Severity: SeverityWarning,
		Description: "`OCTAL_ZEROES` is off in Zsh by default: arithmetic treats `0100` as the " +
			"decimal integer one hundred, matching what every other scripting language " +
			"does. Setting it on reverts to POSIX-shell semantics where the leading `0` " +
			"flags the literal as octal — `(( n = 0100 ))` assigns 64, not 100. Scripts " +
			"that read timestamps padded to `00:59`, CSVs of phone-number prefixes " +
			"(`0049`), or file modes formatted as `0700` silently return the wrong " +
			"integer. Keep the option off at script level; if you really want C-style " +
			"octal literals, stay explicit with `(( n = 8#100 ))` or `$(( 8#$val ))` " +
			"so the intent is obvious.",
		Check: checkZC1861,
	})
}

func checkZC1861(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1861IsOctalZeroes(arg.String()) {
				return zc1861Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOOCTALZEROES" {
				return zc1861Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1861IsOctalZeroes(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "OCTALZEROES"
}

func zc1861Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1861",
		Message: "`" + where + "` reinterprets leading-zero integers as octal — " +
			"`(( n = 0100 ))` assigns 64 instead of 100, breaking timestamp, " +
			"phone-prefix, and mode parsing. Keep the option off; use `8#100` " +
			"when you want explicit octal.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1862",
		Title:    "Warn on `ssh-keygen -R HOST` — deletes a known-hosts entry, next `ssh` re-trusts silently",
		Severity: SeverityWarning,
		Description: "`ssh-keygen -R HOST` scrubs the entry for `HOST` from `~/.ssh/known_hosts`. " +
			"The legitimate trigger is a real key rotation (server reinstall, HSM " +
			"replacement), but the flag is frequently dropped into automation to " +
			"silence the REMOTE HOST IDENTIFICATION HAS CHANGED banner without ever " +
			"confirming the new fingerprint. The very next `ssh` call then prompts " +
			"once (or not at all under `StrictHostKeyChecking=no`) and blindly accepts " +
			"whatever the network hands back — a MITM attacker who was waiting for a " +
			"rebuild slips in without a trace. Fetch the new key out-of-band and " +
			"`ssh-keyscan -t rsa,ed25519 HOST | ssh-keygen -lf -` before adding it, or " +
			"pin fingerprints in a managed `known_hosts` file.",
		Check: checkZC1862,
	})
}

func checkZC1862(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ssh-keygen" {
		return nil
	}
	args := cmd.Arguments
	for i, arg := range args {
		v := arg.String()
		if v != "-R" {
			continue
		}
		if i+1 >= len(args) {
			return nil
		}
		host := args[i+1].String()
		return []Violation{{
			KataID: "ZC1862",
			Message: "`ssh-keygen -R " + host + "` deletes a known-hosts entry — the " +
				"next `ssh` silently re-trusts whatever key the network returns. " +
				"Fetch the new fingerprint out-of-band and verify before " +
				"re-adding.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1863",
		Title:    "Warn on `unsetopt CASE_GLOB` — globs silently go case-insensitive across the script",
		Severity: SeverityWarning,
		Description: "`CASE_GLOB` on is the Zsh default: `*.log` matches `app.log` but not " +
			"`APP.LOG`, `[A-Z]*` is a real case-sensitive range, and `[[ $f == Foo* ]]` " +
			"keeps the distinction between `Foo1` and `foo1`. Turning it off (or " +
			"equivalently `setopt NO_CASE_GLOB`) silently re-evaluates every subsequent " +
			"pattern case-insensitively — `rm *.log` now sweeps `APP.LOG` up, pattern " +
			"dispatchers that used to distinguish `README` from `readme` stop doing so, " +
			"and hash maps keyed on glob-built labels start colliding. Keep the option " +
			"on at script level; request case-folding per-pattern with the Zsh qualifier " +
			"`(#i)*.log`.",
		Check: checkZC1863,
	})
}

func checkZC1863(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1863IsCaseGlob(arg.String()) {
				return zc1863Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOCASEGLOB" {
				return zc1863Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1863IsCaseGlob(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "CASEGLOB"
}

func zc1863Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1863",
		Message: "`" + where + "` flips every later glob to case-insensitive — " +
			"`rm *.log` sweeps `APP.LOG`, dispatchers keyed on case collisions. " +
			"Keep the option on; use `(#i)pattern` per-glob when you need " +
			"case-folding.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1864",
		Title:    "Error on `mount -o remount,exec` — re-enables exec on a previously `noexec` mount",
		Severity: SeverityError,
		Description: "Hardened systems mount `/tmp`, `/var/tmp`, `/dev/shm`, and `/home` with " +
			"`noexec` so a dropper cannot chmod and launch a payload out of a world-writable " +
			"directory. `mount -o remount,exec /tmp` (or the narrower `remount,suid`) " +
			"removes that guardrail for the live kernel, and every shell that already had " +
			"`cd /tmp` open picks it up immediately. Most legitimate uses come from install " +
			"scripts that briefly relax `noexec`; those scripts should restore the flag in " +
			"a `trap 'mount -o remount,noexec /tmp' EXIT`. Blanket `remount,exec` without a " +
			"restore path leaves the system in a permanently weakened state until reboot.",
		Check: checkZC1864,
	})
}

func checkZC1864(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mount" {
		return nil
	}
	args := cmd.Arguments
	for i, arg := range args {
		v := arg.String()
		var opts string
		switch {
		case v == "-o" && i+1 < len(args):
			opts = args[i+1].String()
		case strings.HasPrefix(v, "-o"):
			opts = strings.TrimPrefix(v, "-o")
		default:
			continue
		}
		opts = strings.ToLower(strings.Trim(opts, "\"'"))
		if !strings.Contains(opts, "remount") {
			continue
		}
		if weak := zc1864FirstWeakenedFlag(opts); weak != "" {
			return []Violation{{
				KataID: "ZC1864",
				Message: "`mount -o " + opts + "` re-enables `" + weak + "` on a " +
					"`noexec`/`nosuid`/`nodev`-hardened mount — dropped payloads " +
					"suddenly execute. Pair with a `trap ... EXIT` that restores " +
					"the original flags or skip the remount.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1864FirstWeakenedFlag(opts string) string {
	for _, entry := range strings.Split(opts, ",") {
		entry = strings.TrimSpace(entry)
		switch entry {
		case "exec", "suid", "dev":
			return entry
		}
	}
	return ""
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1865",
		Title:    "Warn on `unsetopt CASE_MATCH` — `[[ =~ ]]` and pattern tests quietly fold case",
		Severity: SeverityWarning,
		Description: "`CASE_MATCH` on is Zsh's default: `[[ $x =~ ^FOO ]]`, `[[ $x == Foo* ]]`, " +
			"and the subst-in-conditional forms honour letter case exactly as written. " +
			"Turning the option off flips every later test to case-insensitive — " +
			"`[[ $user == Admin ]]` also matches `admin`/`ADMIN`, regex dispatchers stop " +
			"distinguishing `README` from `readme`, and log-pattern filters over-collect. " +
			"Keep the option on at script level; if one specific regex really needs " +
			"case-folding, request it per-pattern with the Zsh `(#i)` flag " +
			"(e.g. `[[ $x =~ (#i)foo ]]`).",
		Check: checkZC1865,
	})
}

func checkZC1865(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1865IsCaseMatch(arg.String()) {
				return zc1865Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOCASEMATCH" {
				return zc1865Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1865IsCaseMatch(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "CASEMATCH"
}

func zc1865Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1865",
		Message: "`" + where + "` flips every `[[ =~ ]]` / `[[ == pat ]]` to " +
			"case-insensitive — `Admin` matches `ADMIN`, dispatchers collide. " +
			"Keep it on; scope per-line with `(#i)pattern`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1866",
		Title:    "Warn on `docker exec -u 0` — bypasses the image's non-root `USER` directive",
		Severity: SeverityWarning,
		Description: "A hardened image runs with a non-root `USER` set in its Dockerfile so " +
			"exploited processes inside the container are contained by the Linux " +
			"user-namespace mapping. `docker exec -u 0` (and `-u root`, `--user=0`, the " +
			"podman equivalent) overrides that choice on a per-exec basis and drops a " +
			"shell back into uid 0 — every subsequent file write, cap check, and namespace " +
			"test now runs as root inside the container, which on a default Docker setup " +
			"is also root on the host via the shared mount namespace. Keep exec sessions " +
			"as the container's configured user; if you genuinely need root for a one-off " +
			"fix, document it in the ticket and consider rebuilding the image with the " +
			"capability baked in so `-u 0` is never required.",
		Check: checkZC1866,
	})
}

func checkZC1866(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}
	args := cmd.Arguments
	if len(args) == 0 || args[0].String() != "exec" {
		return nil
	}
	for i := 1; i < len(args); i++ {
		v := args[i].String()
		var user string
		switch {
		case (v == "-u" || v == "--user") && i+1 < len(args):
			user = args[i+1].String()
		case strings.HasPrefix(v, "-u") && v != "-u":
			user = strings.TrimPrefix(v, "-u")
		case strings.HasPrefix(v, "--user="):
			user = strings.TrimPrefix(v, "--user=")
		default:
			continue
		}
		user = strings.Trim(user, "\"'")
		if zc1866IsRoot(user) {
			return []Violation{{
				KataID: "ZC1866",
				Message: "`" + ident.Value + " exec -u " + user + "` drops a root " +
					"shell — bypasses the image's non-root `USER` and, without " +
					"userns remap, equals host root. Keep execs as the container " +
					"user.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1866IsRoot(v string) bool {
	if v == "0" || v == "root" || v == "0:0" {
		return true
	}
	// `0:gid` or `0:groupname` — still uid 0.
	if strings.HasPrefix(v, "0:") {
		return true
	}
	if strings.HasPrefix(v, "root:") {
		return true
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1867",
		Title:    "Warn on `unsetopt GLOB` — pattern expansion turned off, `rm *.log` tries the literal filename",
		Severity: SeverityWarning,
		Description: "`GLOB` is on by default in Zsh: `*`, `?`, `[...]`, and `**/` expand against " +
			"the filesystem before the command runs. Turning the option off script-wide " +
			"(via `unsetopt GLOB` or the equivalent `setopt NO_GLOB`, same as POSIX " +
			"`set -f`) means every later pattern is handed to the command verbatim, so " +
			"`rm *.log` tries to remove a file literally named `*.log`, `for f in *.txt` " +
			"iterates over the single literal string, and expected-array-length checks " +
			"always return 1. Keep the option on at the script level; if one specific " +
			"line needs the pattern as a literal, quote the argument (`'*.log'`) or scope " +
			"with `setopt LOCAL_OPTIONS; setopt NO_GLOB` inside a function.",
		Check: checkZC1867,
	})
}

func checkZC1867(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1867IsGlob(arg.String()) {
				return zc1867Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOGLOB" {
				return zc1867Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1867IsGlob(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "GLOB"
}

func zc1867Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1867",
		Message: "`" + where + "` disables glob expansion — `rm *.log` chases the " +
			"literal `*.log`, `for f in *.txt` loops once. Quote specific args or " +
			"scope with `LOCAL_OPTIONS` inside a function.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1868",
		Title:    "Error on `gcloud config set auth/disable_ssl_validation true` — disables TLS on every later gcloud call",
		Severity: SeverityError,
		Description: "`gcloud config set auth/disable_ssl_validation true` writes the flag into " +
			"the active configuration file, so every subsequent `gcloud` invocation on " +
			"that machine stops verifying the Google API certificate until someone " +
			"reverses it. A MITM holding a self-signed cert can then intercept service " +
			"account tokens, project-level credentials, and every deploy that runs under " +
			"the same user. Remove the setting (`gcloud config unset " +
			"auth/disable_ssl_validation`), and if a corporate proxy really needs a custom " +
			"CA use `core/custom_ca_certs_file` to pin it rather than disabling the check.",
		Check: checkZC1868,
	})
}

func checkZC1868(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gcloud" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 4 {
		return nil
	}
	if args[0].String() != "config" || args[1].String() != "set" {
		return nil
	}
	key := args[2].String()
	val := strings.ToLower(args[3].String())
	if key == "auth/disable_ssl_validation" && (val == "true" || val == "1" || val == "on") {
		return []Violation{{
			KataID: "ZC1868",
			Message: "`gcloud config set auth/disable_ssl_validation " + args[3].String() +
				"` turns off TLS for every later `gcloud` call — service-account " +
				"tokens and deploys become interceptable. Unset it; pin custom " +
				"CAs via `core/custom_ca_certs_file`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1869",
		Title:    "Warn on `setopt RC_EXPAND_PARAM` — brace-adjacent array expansion silently distributes",
		Severity: SeverityWarning,
		Description: "`RC_EXPAND_PARAM` is off in Zsh by default: `echo x${arr[@]}y` concatenates " +
			"once, producing `xay xby xcy` only if you wrote the template carefully. " +
			"Turning it on changes the rule — every adjacent literal is distributed " +
			"across each array element, so `cp src/${files[@]}.bak /tmp` suddenly " +
			"rewrites as `cp src/a.bak src/b.bak src/c.bak /tmp`. That is exactly what " +
			"you want when you want it, and a nasty surprise anywhere else because the " +
			"same syntax keeps working silently. Leave the option off at script level; " +
			"if one specific line needs distributive expansion, request it per-use with " +
			"`${^arr}` (the `^` flag scopes the behaviour to that parameter only).",
		Check: checkZC1869,
	})
}

func checkZC1869(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1869IsRcExpandParam(arg.String()) {
				return zc1869Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NORCEXPANDPARAM" {
				return zc1869Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1869IsRcExpandParam(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "RCEXPANDPARAM"
}

func zc1869Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1869",
		Message: "`" + where + "` distributes literal prefix/suffix across every " +
			"array element — `cp src/${arr[@]}.bak dst` silently rewrites as " +
			"`cp src/a.bak src/b.bak dst`. Keep it off; opt in per-use with " +
			"`${^arr}`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1870",
		Title:    "Warn on `setopt GLOB_ASSIGN` — RHS of `var=pattern` silently glob-expands",
		Severity: SeverityWarning,
		Description: "`GLOB_ASSIGN` is off by default in Zsh: `logs=*.log` sets `$logs` to the " +
			"literal string `*.log`, just like every other shell. Turning it on expands the " +
			"right-hand side of unquoted assignments — `logs=*.log` silently becomes the " +
			"first matching filename, `latest=backup-*` captures whatever sort-order the " +
			"filesystem returns, and any empty-match case assigns an empty string. Scripts " +
			"that port cleanly between Bash and Zsh suddenly diverge, and sensitive " +
			"assignments like `cert=~/secrets/*` can grab attacker-dropped files. Keep the " +
			"option off; use `set -A arr *.log` or explicit `arr=( *.log )` when you really " +
			"want the expansion.",
		Check: checkZC1870,
	})
}

func checkZC1870(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1870IsGlobAssign(arg.String()) {
				return zc1870Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOGLOBASSIGN" {
				return zc1870Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1870IsGlobAssign(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "GLOBASSIGN"
}

func zc1870Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1870",
		Message: "`" + where + "` expands glob patterns on the RHS of `var=` — " +
			"`logs=*.log` silently captures the first match, `cert=~/secrets/*` " +
			"picks up attacker drops. Keep it off; use explicit `arr=( *.log )`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1871",
		Title:    "Warn on `setopt IGNORE_BRACES` — brace expansion stops working script-wide",
		Severity: SeverityWarning,
		Description: "`IGNORE_BRACES` is off by default in Zsh, which means `{1..10}`, " +
			"`file.{log,bak}`, and nested combinations like `{a..z}{1..3}` all expand " +
			"exactly as they do in Bash with `brace_expand` on. Turning it on disables " +
			"every one of those — `for i in {1..10}` iterates over the single literal " +
			"token `{1..10}`, and `cp app.{conf,conf.bak}` tries to copy a file literally " +
			"called `app.{conf,conf.bak}`. Scripts that depend on either numeric or " +
			"comma-list expansion silently become no-ops or fail with ENOENT. Keep the " +
			"option off; if you really need a literal brace string, quote the specific " +
			"argument (`'app.{conf,bak}'`) instead of flipping the shell-wide behaviour.",
		Check: checkZC1871,
	})
}

func checkZC1871(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1871IsIgnoreBraces(arg.String()) {
				return zc1871Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOIGNOREBRACES" {
				return zc1871Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1871IsIgnoreBraces(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "IGNOREBRACES"
}

func zc1871Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1871",
		Message: "`" + where + "` disables brace expansion — `for i in {1..10}` " +
			"loops once over the literal token, `cp app.{conf,bak}` fails ENOENT. " +
			"Keep the option off; quote the specific argument if you need a " +
			"literal brace string.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1872",
		Title:    "Error on `badblocks -w` — destructive write-mode pattern test wipes the device",
		Severity: SeverityError,
		Description: "`badblocks -w` (alias `--write-mode`) runs the write-mode bad-block check, " +
			"which overwrites every sector of the target device with a test pattern and " +
			"reads it back. On a fresh drive about to be formatted that is exactly what " +
			"you want; on an already-populated disk it is a silent data-wipe — the " +
			"command returns success even as it bulldozes the filesystem. If only " +
			"non-destructive checking is needed, use `badblocks -n` (read-test-restore) " +
			"or `badblocks` without any mode flag (read-only). When a true destructive " +
			"test is intended, gate the call behind a confirmation prompt and a freshly " +
			"partitioned device.",
		Check: checkZC1872,
	})
}

func checkZC1872(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "badblocks" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-w" || v == "--write-mode" {
			return zc1872Hit(cmd)
		}
		// Also catch combined short-option clusters like `-wsv`.
		if len(v) > 1 && v[0] == '-' && v[1] != '-' && strings.ContainsRune(v[1:], 'w') {
			return zc1872Hit(cmd)
		}
	}
	return nil
}

func zc1872Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1872",
		Message: "`badblocks -w` overwrites every sector of the target device — " +
			"silent data wipe on a populated disk. Use `-n` (non-destructive) " +
			"or gate destructive runs behind a confirmation and a fresh partition.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1873",
		Title:    "Warn on `setopt ERR_RETURN` — functions silently bail out on the first non-zero exit",
		Severity: SeverityWarning,
		Description: "`ERR_RETURN` is the function-scoped cousin of `ERR_EXIT` and is off by " +
			"default in Zsh. Turning it on script-wide makes every function `return` at " +
			"the first command whose status is non-zero, which in practice means helpers " +
			"that deliberately probe the environment (`test -f /some/file`, `grep -q " +
			"PATTERN`, `id -u user`) will bail before they reach the branch that was meant " +
			"to run when the probe failed. Callers see a success-or-nothing return and " +
			"no stderr. Keep the option off at script level; inside one function that " +
			"really wants fail-fast semantics, scope with `setopt LOCAL_OPTIONS; setopt " +
			"ERR_RETURN` so the behaviour cannot leak to the rest of the shell.",
		Check: checkZC1873,
	})
}

func checkZC1873(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1873IsErrReturn(arg.String()) {
				return zc1873Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOERRRETURN" {
				return zc1873Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1873IsErrReturn(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "ERRRETURN"
}

func zc1873Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1873",
		Message: "`" + where + "` returns from every function on first non-zero " +
			"exit — probing helpers (`test -f`, `grep -q`) bail before the " +
			"fallback branch. Scope inside a `LOCAL_OPTIONS` function if " +
			"needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1874",
		Title:    "Warn on `sshuttle -r HOST 0/0` — every outbound packet tunneled through the jump host",
		Severity: SeverityWarning,
		Description: "`sshuttle -r user@host 0/0` (or `0.0.0.0/0`, `::/0`) installs a VPN-like " +
			"catch-all route: every TCP connection and DNS lookup on the local machine " +
			"egresses through `user@host`, including traffic to corporate VPN endpoints, " +
			"cloud APIs, and package mirrors that had been explicitly split-tunnel. If the " +
			"jump host is compromised, misconfigured, or simply overloaded, every session " +
			"on the workstation silently degrades or leaks to the wrong peer. Scope the " +
			"subnet list to the networks you actually need (`10.0.0.0/8 172.16.0.0/12 " +
			"192.168.0.0/16`), or prefer `ssh -D` with `--exclude` rules for a single " +
			"browser profile.",
		Check: checkZC1874,
	})
}

func checkZC1874(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sshuttle" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1874IsDefaultRoute(v) {
			return []Violation{{
				KataID: "ZC1874",
				Message: "`sshuttle ... " + v + "` routes every outbound packet through " +
					"the jump host — a compromise of `user@host` sees the whole " +
					"fleet's traffic. Scope to the subnets you actually need.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1874IsDefaultRoute(v string) bool {
	switch v {
	case "0/0", "0.0.0.0/0", "::/0":
		return true
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1875",
		Title:    "Warn on `setopt RC_QUOTES` — `''` inside single quotes flips from empty-concat to literal apostrophe",
		Severity: SeverityWarning,
		Description: "`RC_QUOTES` is off by default in Zsh: inside a single-quoted string `'it''s'` " +
			"parses as two adjacent single-quoted regions with an empty middle, producing " +
			"the literal `its`. Turning the option on reinterprets the doubled apostrophe " +
			"as one escaped quote, so `'it''s'` suddenly becomes `it's`. That is a " +
			"source-level change to every already-written string literal in the file — " +
			"password strings, SQL fragments, display text — so log lines, stored tokens, " +
			"and API payloads silently diverge. Keep the option off; write a literal " +
			"apostrophe with `\\'` outside the quotes or with double-quoted wrapping.",
		Check: checkZC1875,
	})
}

func checkZC1875(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1875IsRCQuotes(arg.String()) {
				return zc1875Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NORCQUOTES" {
				return zc1875Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1875IsRCQuotes(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "RCQUOTES"
}

func zc1875Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1875",
		Message: "`" + where + "` reinterprets `''` inside single quotes as a " +
			"literal apostrophe — `'it''s'` flips from `its` to `it's`, " +
			"breaking tokens and SQL. Use double quotes or `\\'` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1876",
		Title:    "Warn on `cargo publish --allow-dirty` — publishes the crate with uncommitted local changes",
		Severity: SeverityWarning,
		Description: "`cargo publish` by default refuses to upload when the working tree is dirty, " +
			"because the published tarball is a snapshot of whatever is on disk — not " +
			"whatever is committed. `--allow-dirty` skips that check, so a `println!` " +
			"dropped in for debugging, an uncommitted `Cargo.toml` dep bump, or a " +
			"`patch.crates-io` override that only exists locally ends up on crates.io " +
			"under the same version users see on GitHub. This is irreversible — once a " +
			"version is uploaded it cannot be replaced, only yanked. Commit first and " +
			"publish from a clean checkout; if you truly must publish from a dirty tree, " +
			"scope the flag to a one-off manual call with a `--dry-run` rehearsal first.",
		Check: checkZC1876,
	})
}

func checkZC1876(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cargo" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 2 {
		return nil
	}
	if args[0].String() != "publish" {
		return nil
	}
	for _, arg := range args[1:] {
		if arg.String() == "--allow-dirty" {
			return []Violation{{
				KataID: "ZC1876",
				Message: "`cargo publish --allow-dirty` uploads a tarball snapshot of " +
					"the dirty working tree — debug prints and local-only patches " +
					"end up on crates.io for a version that cannot be replaced. " +
					"Commit first; `--dry-run` to rehearse.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1877",
		Title:    "Warn on `unsetopt SHORT_LOOPS` — short-form `for`/`while` bodies stop parsing",
		Severity: SeverityWarning,
		Description: "`SHORT_LOOPS` is on in Zsh by default: the compact forms `for x in *.log; " +
			"print $x`, `while true; print .`, and `repeat 3 sleep 1` parse with an implicit " +
			"single-command body. Turning the option off reverts to POSIX-shell parsing, " +
			"which demands an explicit `do ... done` or `{ ... }` block. Every subsequent " +
			"short-form loop raises a parse error (`parse error near '\\n'`), and the " +
			"behaviour is global so even helper files sourced later fall over. Keep the " +
			"option on; if you genuinely need POSIX-strict parsing, scope inside a function " +
			"with `setopt LOCAL_OPTIONS; unsetopt SHORT_LOOPS`.",
		Check: checkZC1877,
	})
}

func checkZC1877(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1877IsShortLoops(arg.String()) {
				return zc1877Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOSHORTLOOPS" {
				return zc1877Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1877IsShortLoops(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "SHORTLOOPS"
}

func zc1877Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1877",
		Message: "`" + where + "` disables short-form loops — `for f in *.log; " +
			"print $f` raises a parse error. Keep the option on; scope inside " +
			"a function with `LOCAL_OPTIONS` if POSIX-strict parsing is " +
			"really needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1878",
		Title:    "Warn on `kubectl apply --force-conflicts` — steals ownership of fields managed by other controllers",
		Severity: SeverityWarning,
		Description: "Server-side apply tracks every field of a resource by the applier that " +
			"last set it (`metadata.managedFields`). When two appliers disagree, the " +
			"default behaviour is to abort with `conflict` so you can reconcile " +
			"deliberately. `kubectl apply --server-side --force-conflicts` overrides " +
			"that: the current caller snatches ownership of every conflicting field — " +
			"including fields set by operators, HPA, cert-manager, and webhook-injected " +
			"sidecars — and those controllers will silently lose their reconcile " +
			"pressure until their next write. Resolve the conflict instead: either " +
			"drop the disputed fields from your manifest so the other owner can keep " +
			"them, or coordinate a hand-off by first removing the managed-field entry " +
			"(`kubectl apply --field-manager=... --subresource=...`).",
		Check: checkZC1878,
	})
}

func checkZC1878(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubectl" {
		return nil
	}
	args := cmd.Arguments
	if len(args) == 0 || args[0].String() != "apply" {
		return nil
	}
	for _, arg := range args[1:] {
		if arg.String() == "--force-conflicts" {
			return []Violation{{
				KataID: "ZC1878",
				Message: "`kubectl apply --force-conflicts` grabs ownership of every " +
					"conflicting field from other controllers (HPA, cert-manager, " +
					"sidecar injectors). Resolve the conflict instead — drop the " +
					"disputed fields or hand off via managed-field edit.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1879",
		Title:    "Warn on `unsetopt BAD_PATTERN` — malformed glob patterns silently pass through as literals",
		Severity: SeverityWarning,
		Description: "`BAD_PATTERN` is on in Zsh by default: a syntactically broken glob (unbalanced " +
			"`[`, stray `^` outside extended-glob context, runaway `(alt|…`) produces a " +
			"`zsh: bad pattern` error so the script knows the filename filter is wrong. " +
			"Turning the option off reverts to POSIX behaviour — the pattern is handed to " +
			"the command verbatim, and `rm [abc` silently tries to remove a file literally " +
			"called `[abc`. Malformed patterns routed to `find -name` or passed to `case` " +
			"blocks likewise stop firing. Keep the option on at script level; if one " +
			"particular line really needs POSIX pass-through, quote the pattern or scope " +
			"with `setopt LOCAL_OPTIONS; unsetopt BAD_PATTERN` inside a function.",
		Check: checkZC1879,
	})
}

func checkZC1879(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1879IsBadPattern(arg.String()) {
				return zc1879Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOBADPATTERN" {
				return zc1879Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1879IsBadPattern(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "BADPATTERN"
}

func zc1879Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1879",
		Message: "`" + where + "` silences `bad pattern` errors — `rm [abc` tries " +
			"to remove a literal `[abc`, broken `case` arms stop firing. Keep " +
			"the option on; quote one-off patterns or scope with `LOCAL_OPTIONS`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1880",
		Title:    "Warn on `kubectl annotate|label --overwrite` — silently rewrites controller signals",
		Severity: SeverityWarning,
		Description: "Kubernetes annotations and labels are not plain metadata — they are the " +
			"protocol by which cert-manager, external-dns, ingress-nginx, the " +
			"HorizontalPodAutoscaler, and most Helm-managed controllers decide what to " +
			"do with a resource. `kubectl annotate --overwrite` and `kubectl label " +
			"--overwrite` suppress the conflict check and replace whatever value was " +
			"there, so the script silently rewrites `kubectl.kubernetes.io/last-applied-" +
			"configuration`, `cert-manager.io/cluster-issuer`, or " +
			"`prometheus.io/scrape`, triggering reissue / reconfiguration or breaking " +
			"the next apply. Inspect the existing annotation with `kubectl get -o " +
			"jsonpath='{.metadata.annotations}'` first, and drop `--overwrite` so a " +
			"conflict surfaces as an error.",
		Check: checkZC1880,
	})
}

func checkZC1880(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubectl" {
		return nil
	}
	args := cmd.Arguments
	if len(args) == 0 {
		return nil
	}
	sub := args[0].String()
	if sub != "annotate" && sub != "label" {
		return nil
	}
	for _, arg := range args[1:] {
		if arg.String() == "--overwrite" {
			return []Violation{{
				KataID: "ZC1880",
				Message: "`kubectl " + sub + " --overwrite` silently replaces an " +
					"existing controller signal — cert-manager, external-dns, " +
					"HPA watchers reconcile on the new value. Inspect first; " +
					"drop `--overwrite` so conflicts error.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1881",
		Title:    "Warn on `unsetopt MULTIBYTE` — `${#str}`, substring, and `[[ =~ ]]` stop counting characters",
		Severity: SeverityWarning,
		Description: "`MULTIBYTE` is on in Zsh by default: `${#str}` returns character count, " +
			"`${str:0:3}` extracts the first three characters, and `[[ $str =~ ... ]]` " +
			"matches whole UTF-8 codepoints. Turning it off reverts every string " +
			"operation to per-byte math, so an emoji that encodes to four bytes counts " +
			"as four, a substring spanning a multi-byte character slices mid-codepoint " +
			"and produces invalid UTF-8, and `[[ =~ ]]` regex ranges no longer cover " +
			"Unicode blocks. Filenames containing non-ASCII, i18n log strings, and JSON " +
			"snippets silently drift from their assumed layout. Keep the option on; if " +
			"you truly need byte-level counting, use `${#${(%)str}}` or `wc -c <<< $str`.",
		Check: checkZC1881,
	})
}

func checkZC1881(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1881IsMultibyte(arg.String()) {
				return zc1881Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOMULTIBYTE" {
				return zc1881Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1881IsMultibyte(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "MULTIBYTE"
}

func zc1881Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1881",
		Message: "`" + where + "` flips every string op to per-byte math — `${#str}` " +
			"on an emoji returns 4, substrings slice mid-codepoint, `[[ =~ ]]` " +
			"Unicode ranges break. Keep the option on; byte-count with " +
			"`wc -c <<< $str` when truly needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1882",
		Title:    "Warn on `sudo -s` / `sudo su` / `sudo bash` — spawns an interactive root shell from a script",
		Severity: SeverityWarning,
		Description: "`sudo -s`, `sudo -i`, `sudo su [-]`, and `sudo bash` (or `zsh`/`sh`/`ksh`) " +
			"with no trailing command hand you an interactive root shell. That is fine " +
			"at a prompt, but in a non-interactive script the shell either hangs " +
			"waiting for stdin or drains stdin into root's shell as if those lines were " +
			"the shell's commands — neither is what the script author meant. Pass the " +
			"actual command to sudo (`sudo /usr/local/bin/provision.sh`) so the " +
			"elevation is scoped and audit logs capture the real work.",
		Check: checkZC1882,
	})
}

func checkZC1882(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sudo" {
		return nil
	}
	args := cmd.Arguments
	if len(args) == 0 {
		return nil
	}

	first := args[0].String()
	rest := args[1:]

	// sudo -s / sudo -i with no trailing positional command.
	if (first == "-s" || first == "-i") && !zc1882HasPositional(rest) {
		return zc1882Hit(cmd, "sudo "+first)
	}

	// sudo su, sudo su -, sudo su -l, sudo su --login (no -c).
	if first == "su" {
		if !zc1882HasArg(rest, "-c", "--command") {
			return zc1882Hit(cmd, "sudo su")
		}
	}

	// sudo bash / sudo zsh / sudo sh / sudo ksh without -c.
	switch first {
	case "bash", "zsh", "sh", "ksh", "dash", "ash":
		if !zc1882HasArg(rest, "-c") {
			return zc1882Hit(cmd, "sudo "+first)
		}
	}
	return nil
}

func zc1882HasPositional(args []ast.Expression) bool {
	for _, a := range args {
		v := a.String()
		if v == "" || v[0] == '-' {
			continue
		}
		return true
	}
	return false
}

func zc1882HasArg(args []ast.Expression, names ...string) bool {
	for _, a := range args {
		v := a.String()
		for _, n := range names {
			if v == n {
				return true
			}
		}
	}
	return false
}

func zc1882Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1882",
		Message: "`" + where + "` spawns an interactive root shell — in a script " +
			"either hangs on stdin or drains the rest of the file into root's " +
			"shell. Pass the command to sudo: `sudo /path/to/cmd arg …`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1883",
		Title:    "Warn on `setopt PATH_SCRIPT` — `. ./script.sh` silently falls back to `$PATH` lookup",
		Severity: SeverityWarning,
		Description: "`PATH_SCRIPT` (off by default) lets the `.` builtin and `source` fall back to " +
			"a `$PATH` walk when the literal path resolves to no file. With it on, " +
			"`. helper.sh` looks for `helper.sh` in every `$path` entry — including " +
			"user-owned directories like `~/bin` or `./` — and silently sources whichever " +
			"matches first. An attacker who can drop `helper.sh` into any `$PATH` " +
			"component runs their code inside the current shell's process, with every " +
			"parent env var and exported secret available. Keep the option off; always " +
			"source scripts with an explicit path (`./helper.sh`, `/opt/…/helper.sh`) so " +
			"the source cannot be redirected.",
		Check: checkZC1883,
	})
}

func checkZC1883(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1883IsPathScript(arg.String()) {
				return zc1883Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOPATHSCRIPT" {
				return zc1883Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1883IsPathScript(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "PATHSCRIPT"
}

func zc1883Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1883",
		Message: "`" + where + "` lets `.`/`source` fall back to `$PATH` when a " +
			"literal path misses — a dropper in `~/bin` or `./` runs inside the " +
			"current shell with every exported secret. Keep the option off; " +
			"always use explicit paths.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1884",
		Title:    "Error on `curl/wget https://...?apikey=...` — credential in URL query string",
		Severity: SeverityError,
		Description: "Anything passed as an HTTP query parameter is logged by every intermediary: " +
			"the server's access log, the transparent proxy, the CDN request-id trail, " +
			"browser referrer headers, and any client-side observability tooling. A URL " +
			"like `https://api.example/widgets?apikey=SECRET&token=xyz` therefore " +
			"tattoos the credential into logs that live forever and are often shared " +
			"with downstream teams. Move the secret into an HTTP header " +
			"(`curl -H \"Authorization: Bearer $TOKEN\"`), a POST body with " +
			"`--data-urlencode` + TLS, or an `-u user:` basic-auth combo — never the " +
			"query string.",
		Check: checkZC1884,
	})
}

var zc1884SecretKeys = []string{
	"apikey=",
	"api_key=",
	"api-key=",
	"token=",
	"access_token=",
	"id_token=",
	"auth_token=",
	"access-token=",
	"password=",
	"passwd=",
	"secret=",
	"client_secret=",
	"sig=",
	"signature=",
}

func checkZC1884(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "curl" && ident.Value != "wget" && ident.Value != "http" && ident.Value != "httpie" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if match := zc1884FirstSecretKey(v); match != "" {
			return []Violation{{
				KataID: "ZC1884",
				Message: "`" + ident.Value + " " + v + "` carries `" + match +
					"...` in the URL query — logged by every proxy, CDN, and " +
					"server access log along the path. Move credentials to " +
					"`-H \"Authorization: Bearer \"$TOKEN\"` or a POST body.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1884FirstSecretKey(v string) string {
	lower := strings.ToLower(strings.Trim(v, "\"'"))
	// Only flag when this looks like a URL; drop anything without `://` or `?`.
	if !strings.Contains(lower, "://") || !strings.Contains(lower, "?") {
		return ""
	}
	// Scan after the `?` boundary.
	idx := strings.Index(lower, "?")
	if idx < 0 {
		return ""
	}
	query := lower[idx:]
	for _, key := range zc1884SecretKeys {
		if strings.Contains(query, key) {
			return strings.TrimSuffix(key, "=")
		}
	}
	return ""
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1885",
		Title:    "Warn on `setopt CSH_NULL_GLOB` — unmatched globs drop instead of erroring when any sibling matches",
		Severity: SeverityWarning,
		Description: "`CSH_NULL_GLOB` (off by default) mimics csh's rule: in a list like " +
			"`rm *.log *.bak *.tmp`, if at least one pattern produces matches the " +
			"remaining unmatched patterns are silently discarded, and only if every " +
			"pattern produces nothing does the shell raise `no match`. That is a " +
			"partial-failure concealer — a genuine typo `rm *.lg *.bak` can still " +
			"delete the `.bak` files while hiding the `.lg` mismatch, and maintenance " +
			"loops that relied on `NOMATCH` to stop on typos pass right through. Keep " +
			"the option off at script level; use `*(N)` per-glob when you want " +
			"null-glob behaviour.",
		Check: checkZC1885,
	})
}

func checkZC1885(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1885IsCshNullGlob(arg.String()) {
				return zc1885Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOCSHNULLGLOB" {
				return zc1885Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1885IsCshNullGlob(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "CSHNULLGLOB"
}

func zc1885Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1885",
		Message: "`" + where + "` silently discards unmatched globs in a list when " +
			"any sibling matches — `rm *.lg *.bak` deletes the `.bak` files " +
			"and hides the typo. Keep the option off; use `*(N)` per-glob.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1886",
		Title:    "Error on `tee/cp/mv/install/dd` writing system shell-init files — persistent privesc surface",
		Severity: SeverityError,
		Description: "`/etc/profile`, `/etc/bash.bashrc`, `/etc/zshrc`, `/etc/zsh/zshenv`, " +
			"`/etc/environment`, and every drop-in under `/etc/profile.d/` are sourced " +
			"by every interactive shell (and `/etc/zshenv` by every Zsh invocation). A " +
			"script that `tee`s, `cp`s, `mv`s, or `dd`s arbitrary content into any of " +
			"those paths becomes a persistent foothold — the next root login runs the " +
			"injected code. These files belong to the packaging system; hand-edit " +
			"carefully, stage a temp file, validate it with a dry-run login, and move " +
			"it into place with an atomic `install -m 644`.",
		Check: checkZC1886,
	})
}

var zc1886SensitivePaths = []string{
	"/etc/profile",
	"/etc/bash.bashrc",
	"/etc/bashrc",
	"/etc/zshrc",
	"/etc/zshenv",
	"/etc/zsh/zshrc",
	"/etc/zsh/zshenv",
	"/etc/zsh/zprofile",
	"/etc/zsh/zlogin",
	"/etc/zprofile",
	"/etc/zlogin",
	"/etc/environment",
}

func checkZC1886(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "tee", "cp", "mv", "install", "dd":
	default:
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1886IsSensitivePath(v) {
			return []Violation{{
				KataID: "ZC1886",
				Message: "`" + ident.Value + " ... " + v + "` writes a shell-init " +
					"file sourced by every interactive shell — persistent " +
					"foothold for the next root login. Stage a temp file, " +
					"validate, and `install -m 644` atomically.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1886IsSensitivePath(v string) bool {
	trimmed := strings.Trim(v, "\"'")
	for _, p := range zc1886SensitivePaths {
		if trimmed == p {
			return true
		}
	}
	return strings.HasPrefix(trimmed, "/etc/profile.d/")
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1887",
		Title:    "Warn on `setopt POSIX_TRAPS` — EXIT/ZERR traps change scope and no longer fire on function return",
		Severity: SeverityWarning,
		Description: "`POSIX_TRAPS` is off by default in Zsh. With it off, `trap cleanup EXIT` " +
			"inside a function fires when that function returns — the idiomatic Zsh way " +
			"to scope cleanup to a scope. Turning the option on reverts to POSIX-sh " +
			"semantics, where the EXIT trap only fires when the whole shell exits and " +
			"is shared across the entire process. Scripts that installed a cleanup trap " +
			"inside `do_work()` expecting it to run at each invocation now leak the " +
			"first trap's handler into everything after, and helpers that counted on " +
			"TRAPZERR / TRAPEXIT function-scoped behaviour silently skip. Keep the " +
			"option off at script level; if a specific line really needs POSIX-scope, " +
			"use `trap … EXIT` at top level and document it.",
		Check: checkZC1887,
	})
}

func checkZC1887(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1887IsPosixTraps(arg.String()) {
				return zc1887Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOPOSIXTRAPS" {
				return zc1887Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1887IsPosixTraps(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "POSIXTRAPS"
}

func zc1887Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1887",
		Message: "`" + where + "` flips `trap ... EXIT` inside functions from " +
			"function-return to shell-exit scope — per-call cleanup leaks across " +
			"the whole shell, TRAPZERR helpers stop firing. Keep the option off.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1888",
		Title:    "Warn on `aws iam create-access-key` — mints long-lived static AWS credentials",
		Severity: SeverityWarning,
		Description: "`aws iam create-access-key` hands out a static `AKIA.../secret` pair that is " +
			"valid forever until someone rotates it; whoever gets the pair speaks for " +
			"the IAM user on every API call AWS accepts. Most modern deploys no longer " +
			"need these: EC2 instance profiles, EKS/IRSA, Lambda roles, GitHub OIDC, " +
			"and IAM Identity Center all hand out short-lived session credentials on " +
			"demand. Prefer those; if a static key is genuinely required (legacy third-" +
			"party tooling), store it in AWS Secrets Manager, scope the user to the " +
			"narrowest policy possible, and rotate on a schedule with `aws iam update-" +
			"access-key --status Inactive` / `delete-access-key`.",
		Check: checkZC1888,
	})
}

func checkZC1888(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "aws" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 2 {
		return nil
	}
	if args[0].String() != "iam" || args[1].String() != "create-access-key" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1888",
		Message: "`aws iam create-access-key` mints a long-lived `AKIA.../secret` — " +
			"prefer short-lived creds via instance profiles, IRSA, Lambda roles, " +
			"or OIDC federation. If static keys are unavoidable, store in Secrets " +
			"Manager and rotate.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1889",
		Title:    "Error on `skopeo copy --src-tls-verify=false` / `--dest-tls-verify=false` — MITM on image copy",
		Severity: SeverityError,
		Description: "`skopeo copy` is the glue for promoting container images between registries in " +
			"CI, mirroring upstream images into internal caches, and rehydrating images " +
			"to an air-gapped registry. `--src-tls-verify=false` and " +
			"`--dest-tls-verify=false` drop certificate verification on the respective " +
			"leg, which means any on-path attacker can substitute a malicious manifest or " +
			"layer and the copy completes without a warning. Use `--src-cert-dir`/" +
			"`--dest-cert-dir` to pin a private CA if you are mirroring to or from an " +
			"internal registry with self-signed certs, or fix the upstream's cert.",
		Check: checkZC1889,
	})
}

func checkZC1889(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "skopeo" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		lower := strings.ToLower(v)
		for _, prefix := range []string{"--src-tls-verify=", "--dest-tls-verify=", "--tls-verify="} {
			if !strings.HasPrefix(lower, prefix) {
				continue
			}
			val := strings.TrimPrefix(lower, prefix)
			val = strings.Trim(val, "\"'")
			if val == "false" || val == "0" || val == "no" || val == "off" {
				return []Violation{{
					KataID: "ZC1889",
					Message: "`skopeo " + v + "` disables TLS verification on image " +
						"copy — on-path attacker can substitute a malicious manifest. " +
						"Pin a private CA with `--src-cert-dir`/`--dest-cert-dir` " +
						"instead.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1890",
		Title:    "Error on `kadmin -w PASS` / `kinit` with password arg — Kerberos password in argv",
		Severity: SeverityError,
		Description: "`kadmin -w PASS` and `kadmin.local -w PASS` pass the Kerberos admin " +
			"principal's password directly as an argv element. Every `ps`, `/proc/<pid>/" +
			"cmdline`, history file, and CI-pipeline log therefore sees it in plaintext, " +
			"which is catastrophic for an account that can edit the realm's KDC. Use " +
			"`-k -t /etc/krb5.keytab` for non-interactive auth (keytab permissioned to " +
			"root only), or pipe the password through stdin with the `-q` batch form so " +
			"it never rides in argv.",
		Check: checkZC1890,
	})
}

func checkZC1890(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kadmin" && ident.Value != "kadmin.local" && ident.Value != "kpasswd" {
		return nil
	}
	args := cmd.Arguments
	for i := 0; i+1 < len(args); i++ {
		if args[i].String() != "-w" {
			continue
		}
		val := args[i+1].String()
		if val == "" || val[0] == '-' {
			continue
		}
		return []Violation{{
			KataID: "ZC1890",
			Message: "`" + ident.Value + " -w " + val + "` embeds the Kerberos " +
				"admin password in argv — visible to `ps`, `/proc`, shell history. " +
				"Use `-k -t /etc/krb5.keytab` (keytab root-only) or pipe the " +
				"password on stdin.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1891",
		Title:    "Error on `kubectl config view --raw` — prints the full kubeconfig with client keys",
		Severity: SeverityError,
		Description: "`kubectl config view` by default redacts secrets: `client-certificate-data`, " +
			"`client-key-data`, `token`, and `password` fields are replaced with `REDACTED`. " +
			"Adding `--raw` (or the synonym `-R`) undoes every redaction and prints the " +
			"client's base64-encoded private key, bearer tokens, and any embedded user " +
			"password to stdout. In a script where stdout lands in CI log storage, a " +
			"`journalctl` ring buffer, or a Slack paste, the entire kubeconfig walks out. " +
			"Emit only the specific field you need (e.g. `kubectl config view -o " +
			"jsonpath='{.current-context}'`) or decrypt once into a temp file and `shred` it.",
		Check: checkZC1891,
	})
}

func checkZC1891(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubectl" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 2 {
		return nil
	}
	if args[0].String() != "config" || args[1].String() != "view" {
		return nil
	}
	for _, arg := range args[2:] {
		v := arg.String()
		if v == "--raw" || v == "-R" {
			return []Violation{{
				KataID: "ZC1891",
				Message: "`kubectl config view --raw` prints the full kubeconfig " +
					"including client-certificate/key-data and bearer tokens — " +
					"any script-captured stdout exfiltrates the creds. Emit " +
					"the specific field with `-o jsonpath='…'`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1892",
		Title:    "Error on `install -m 4755|6755|2755` — sets setuid/setgid bit at install time",
		Severity: SeverityError,
		Description: "`install -m <mode>` with the setuid (`4xxx`), setgid (`2xxx`), or combined " +
			"(`6xxx`) octal prefix creates the target with those special bits set, which " +
			"turns every execution into a privilege-elevation vector. An uninspected " +
			"binary installed this way — especially from a build script or package " +
			"post-install — becomes a persistent local-privesc primitive if the binary " +
			"is writable, has command-injection, or links against attacker-influenced " +
			"libraries. Drop the setuid/setgid bits from the mode (`install -m 0755`) and " +
			"grant the narrow capability the program actually needs with `setcap " +
			"cap_net_bind_service+ep`; audit the remaining setuid binaries with " +
			"`find / -perm -4000`.",
		Check: checkZC1892,
	})
}

func checkZC1892(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "install" && ident.Value != "mkdir" {
		return nil
	}
	args := cmd.Arguments
	for i, arg := range args {
		v := arg.String()
		var mode string
		switch {
		case v == "-m" && i+1 < len(args):
			mode = args[i+1].String()
		case v == "--mode" && i+1 < len(args):
			mode = args[i+1].String()
		case strings.HasPrefix(v, "-m") && v != "-m":
			mode = strings.TrimPrefix(v, "-m")
		case strings.HasPrefix(v, "--mode="):
			mode = strings.TrimPrefix(v, "--mode=")
		default:
			continue
		}
		if zc1892HasSetuidBits(mode) {
			return []Violation{{
				KataID: "ZC1892",
				Message: "`" + ident.Value + " -m " + mode + "` sets setuid/setgid " +
					"bits at install time — every execution becomes a privesc " +
					"vector. Use `0755` and grant narrow caps with `setcap` instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1892HasSetuidBits(mode string) bool {
	mode = strings.Trim(mode, "\"'")
	if mode == "" {
		return false
	}
	for _, r := range mode {
		if r < '0' || r > '9' {
			return false
		}
	}
	var n int64
	if strings.ContainsAny(mode, "89") {
		n, _ = strconv.ParseInt(mode, 10, 32)
	} else {
		n, _ = strconv.ParseInt(mode, 8, 32)
	}
	// setuid (0o4000), setgid (0o2000)
	return (n & 0o6000) != 0
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1893",
		Title:    "Warn on `unsetopt BARE_GLOB_QUAL` — `*(N)` null-glob qualifier stops being special",
		Severity: SeverityWarning,
		Description: "`BARE_GLOB_QUAL` is on by default in Zsh — that is what makes the " +
			"per-glob qualifier syntax (`*(N)` for null-glob, `*(.x)` for " +
			"executable, `*(Om)` for sort-by-mtime) work. Turning it off reverts " +
			"to ksh-style parsing where `(...)` inside a glob is a pattern " +
			"alternation, so `*(N)` stops being a null-glob and turns into " +
			"\"match zero-or-one N\" — a completely different pattern. Scripts that " +
			"relied on `for f in *.log(N)` to cope with empty directories then " +
			"silently iterate the literal string or fail under NOMATCH. Keep the " +
			"option on; if you really want ksh-style qualifiers, use " +
			"`setopt LOCAL_OPTIONS; unsetopt BARE_GLOB_QUAL` inside a function.",
		Check: checkZC1893,
	})
}

func checkZC1893(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1893IsBareGlobQual(arg.String()) {
				return zc1893Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOBAREGLOBQUAL" {
				return zc1893Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1893IsBareGlobQual(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "BAREGLOBQUAL"
}

func zc1893Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1893",
		Message: "`" + where + "` disables `*(qualifier)` syntax — `*(N)` stops being " +
			"null-glob and becomes an alternation, so null-glob idioms silently " +
			"break. Keep the option on; scope inside a `LOCAL_OPTIONS` function.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1894",
		Title:    "Error on `conntrack -F` / `--flush` — every tracked connection (including SSH) is reset",
		Severity: SeverityError,
		Description: "`conntrack -F` (alias `--flush`) wipes the netfilter connection-tracking " +
			"table. Every established TCP flow that depended on conntrack (every " +
			"stateful-NAT connection, every `-m conntrack --ctstate RELATED,ESTABLISHED` " +
			"allowance, every MASQUERADE session) loses its entry and the next packet is " +
			"matched from scratch; most firewall rulesets drop it as \"new\" and the " +
			"session dies. Over SSH, that means the shell running the very command drops. " +
			"Stage the flush behind `at now + 5 minutes` so the session can re-enter the " +
			"table via a preceding rule, or narrow the scope with `conntrack -D -s " +
			"<client-IP>` for a specific hung flow.",
		Check: checkZC1894,
	})
}

var zc1894FlushFlags = map[string]bool{
	"-F":      true,
	"--flush": true,
}

func checkZC1894(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "conntrack" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if zc1894FlushFlags[arg.String()] {
			line, col := FlagArgPosition(cmd, zc1894FlushFlags)
			return []Violation{{
				KataID: "ZC1894",
				Message: "`conntrack -F` wipes every tracked flow — stateful " +
					"`ctstate ESTABLISHED` allowances drop, running SSH sessions " +
					"lose their entry. Gate with `at now + N min` or narrow to " +
					"one flow with `conntrack -D -s <ip>`.",
				Line:   line,
				Column: col,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1895",
		Title:    "Warn on `setopt NUMERIC_GLOB_SORT` — glob output switches from lexicographic to numeric order",
		Severity: SeverityWarning,
		Description: "`NUMERIC_GLOB_SORT` is off by default: `ls *.log` returns filenames in the " +
			"collation order the filesystem-iteration/sort step produces (lexicographic " +
			"in the C locale, so `app-1.log`, `app-10.log`, `app-2.log`). Turning it on " +
			"makes every subsequent glob and array expansion sort numeric runs " +
			"numerically — the same glob now returns `app-1.log`, `app-2.log`, " +
			"`app-10.log`. Scripts that tail the \"latest\" file by taking the last array " +
			"element, pipelines that expect a specific stable order, and backup rotations " +
			"built on `*[0-9].tar` silently shuffle. Keep the option off script-wide; " +
			"request numeric sort per-glob with the `*(n)` qualifier when needed.",
		Check: checkZC1895,
	})
}

func checkZC1895(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1895IsNumericGlobSort(arg.String()) {
				return zc1895Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NONUMERICGLOBSORT" {
				return zc1895Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1895IsNumericGlobSort(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "NUMERICGLOBSORT"
}

func zc1895Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1895",
		Message: "`" + where + "` switches every later glob to numeric sort — log " +
			"rotations sorted on numeric suffixes silently shuffle. Keep it off; " +
			"use the per-glob `*(n)` qualifier when needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1896",
		Title:    "Error on `docker/podman run -v /proc:…|/sys:…` — bind-mounts host kernel interfaces into container",
		Severity: SeverityError,
		Description: "`docker run -v /proc:/host/proc` (or `-v /sys:…`) bind-mounts the host's " +
			"procfs / sysfs hierarchy into the container's mount namespace. From inside, " +
			"the container can read every host process's `environ` (secrets passed via " +
			"env), every `cmdline`, every `/proc/1/ns/` to open namespace fds for a " +
			"breakout, and `/sys/fs/cgroup` to modify resource limits that affect host " +
			"services. `:ro` does not help — `/proc/<pid>/ns/...` handles remain usable. " +
			"If the container genuinely needs process / kernel visibility, grant the " +
			"narrowest capability instead (`--cap-add=SYS_PTRACE`) or run the monitoring " +
			"agent on the host rather than inside an untrusted image.",
		Check: checkZC1896,
	})
}

func checkZC1896(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}
	args := cmd.Arguments
	if len(args) == 0 || (args[0].String() != "run" && args[0].String() != "create") {
		return nil
	}
	for i := 1; i < len(args); i++ {
		v := args[i].String()
		var mount string
		switch {
		case (v == "-v" || v == "--volume" || v == "--mount") && i+1 < len(args):
			mount = args[i+1].String()
		case strings.HasPrefix(v, "--volume=") || strings.HasPrefix(v, "--mount=") || strings.HasPrefix(v, "-v="):
			mount = v[strings.Index(v, "=")+1:]
		default:
			continue
		}
		if src := zc1896HostKernelSource(mount); src != "" {
			return []Violation{{
				KataID: "ZC1896",
				Message: "`" + ident.Value + " ... -v " + mount + "` bind-mounts host " +
					src + " into the container — every process's `environ`/`cmdline` " +
					"and `/proc/1/ns/` breakout handles become readable. Use " +
					"`--cap-add=SYS_PTRACE` or host-side monitoring instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1896HostKernelSource(v string) string {
	trimmed := strings.Trim(v, "\"'")
	// Accept `source:target[:opts]` bind form and `source=/path,…` mount form.
	source := trimmed
	if idx := strings.Index(trimmed, ":"); idx > 0 {
		source = trimmed[:idx]
	}
	// `--mount type=bind,source=/proc,…`
	if strings.Contains(trimmed, "source=") {
		for _, entry := range strings.Split(trimmed, ",") {
			if strings.HasPrefix(entry, "source=") {
				source = strings.TrimPrefix(entry, "source=")
			} else if strings.HasPrefix(entry, "src=") {
				source = strings.TrimPrefix(entry, "src=")
			}
		}
	}
	switch source {
	case "/proc", "/sys":
		return source
	}
	if strings.HasPrefix(source, "/proc/") || strings.HasPrefix(source, "/sys/") {
		return source
	}
	return ""
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1897",
		Title:    "Warn on `setopt SH_GLOB` — Zsh-specific glob patterns (`*(N)`, `<1-10>`, alternation) stop parsing",
		Severity: SeverityWarning,
		Description: "`SH_GLOB` is off by default in Zsh. With it off, the shell recognises Zsh's " +
			"extended patterns: `*(N)` null-glob qualifier, `<1-10>` numeric range globs, " +
			"`(alt1|alt2)` in-glob alternation, and the whole `(#i)`/`(#c,m)` flag " +
			"family. Turning the option on forces strict POSIX-sh parsing, so the parser " +
			"re-interprets `(...)` as command grouping and the null-glob / range idioms " +
			"raise parse errors. Every kata recommending `*(N)` (see ZC1830, ZC1893) " +
			"silently breaks, and downstream helpers sourced after the setopt inherit the " +
			"restricted pattern syntax. Keep the option off; scope inside a function " +
			"with `setopt LOCAL_OPTIONS; setopt SH_GLOB` if a specific block genuinely " +
			"needs POSIX patterns.",
		Check: checkZC1897,
	})
}

func checkZC1897(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1897IsShGlob(arg.String()) {
				return zc1897Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOSHGLOB" {
				return zc1897Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1897IsShGlob(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "SHGLOB"
}

func zc1897Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1897",
		Message: "`" + where + "` disables Zsh-extended glob patterns — `*(N)` " +
			"qualifiers, `<1-10>` ranges, and `(alt1|alt2)` alternation raise " +
			"parse errors. Keep the option off; scope with `LOCAL_OPTIONS` if " +
			"strict POSIX is needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1898",
		Title:    "Error on `gpg --export-secret-keys` — private-key material leaks to stdout",
		Severity: SeverityError,
		Description: "`gpg --export-secret-keys KEYID` and `--export-secret-subkeys` write the " +
			"ASCII-armoured private key to stdout. In a script, that stream usually lands " +
			"in a file the operator plans to move off-box — and any misstep (wrong " +
			"`cd`, script-wide stdout captured by CI, tee to a world-readable log, " +
			"piped into a remote unencrypted channel) permanently leaks the key. Backup " +
			"the key interactively on an air-gapped machine; if automation is required, " +
			"write the output to a `umask 077`-protected path and immediately encrypt " +
			"with a second symmetric passphrase.",
		Check: checkZC1898,
	})
}

var zc1898ExportFlags = map[string]bool{
	"--export-secret-keys":    true,
	"--export-secret-subkeys": true,
}

func checkZC1898(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "gpg" && ident.Value != "gpg2" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1898ExportFlags[v] {
			line, col := FlagArgPosition(cmd, zc1898ExportFlags)
			return []Violation{{
				KataID: "ZC1898",
				Message: "`gpg " + v + "` writes the private key to stdout — one " +
					"CI-log or wrong-tty redirect leaks it. Back up interactively on an " +
					"air-gapped host, or write to a `umask 077` path and re-encrypt.",
				Line:   line,
				Column: col,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1899",
		Title:    "Error on `mokutil --disable-validation` — turns UEFI Secure Boot off at the shim",
		Severity: SeverityError,
		Description: "`mokutil --disable-validation` queues a request for the shim to stop " +
			"validating the kernel and modules against the enrolled MOK/PK certificates at " +
			"next boot — Secure Boot silently becomes advisory. Any unsigned kernel or " +
			"rootkit module then loads without prompt. Leave Secure Boot validation on; " +
			"if you must load a custom module, enrol its key with `mokutil --import` and " +
			"approve via the `MokManager` prompt at reboot.",
		Check: checkZC1899,
	})
}

var zc1899DisableFlags = map[string]bool{
	"--disable-validation": true,
}

func checkZC1899(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "mokutil" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if zc1899DisableFlags[arg.String()] {
			line, col := FlagArgPosition(cmd, zc1899DisableFlags)
			return []Violation{{
				KataID: "ZC1899",
				Message: "`mokutil --disable-validation` stops the shim from validating " +
					"kernel/modules against enrolled keys — Secure Boot becomes advisory. " +
					"Leave validation on; enrol specific keys with `mokutil --import`.",
				Line:   line,
				Column: col,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
