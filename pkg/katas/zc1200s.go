// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1200",
		Title:    "Avoid `ftp` — use `sftp` or `curl` for secure transfers",
		Severity: SeverityWarning,
		Description: "`ftp` transmits credentials and data in plain text. " +
			"Use `sftp`, `scp`, or `curl` with HTTPS/SFTP for secure file transfers.",
		Check: checkZC1200,
	})
}

func checkZC1200(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ftp" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1200",
		Message: "Avoid `ftp` — it transmits credentials in plain text. " +
			"Use `sftp`, `scp`, or `curl` with HTTPS for secure file transfers.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1201",
		Title:    "Avoid `rsh`/`rlogin`/`rcp` — use `ssh`/`scp`",
		Severity: SeverityWarning,
		Description: "`rsh`, `rlogin`, and `rcp` are insecure legacy protocols. " +
			"Use `ssh`, `scp`, or `rsync` over SSH for encrypted remote operations.",
		Check: checkZC1201,
		Fix:   fixZC1201,
	})
}

// fixZC1201 rewrites the legacy `rsh` / `rlogin` / `rcp` command
// names to `ssh` / `ssh` / `scp` respectively. Single-edit
// replacement at the violation column. Argument syntax is
// compatible (host + optional command for rsh/rlogin/ssh; src dst
// for rcp/scp). Idempotent — a re-run sees `ssh` or `scp`, not
// the legacy names.
func fixZC1201(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	var replacement string
	switch ident.Value {
	case "rsh", "rlogin":
		replacement = "ssh"
	case "rcp":
		replacement = "scp"
	default:
		return nil
	}
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 || off+len(ident.Value) > len(source) {
		return nil
	}
	if string(source[off:off+len(ident.Value)]) != ident.Value {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len(ident.Value),
		Replace: replacement,
	}}
}

func checkZC1201(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value
	if name != "rsh" && name != "rlogin" && name != "rcp" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1201",
		Message: "Avoid `" + name + "` — it is an insecure legacy protocol. " +
			"Use `ssh`/`scp`/`rsync` for encrypted remote operations.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1202",
		Title:    "Avoid `ifconfig` — use `ip` for network configuration",
		Severity: SeverityInfo,
		Description: "`ifconfig` is deprecated on modern Linux. " +
			"Use `ip addr`, `ip link`, or `ip route` from iproute2 for network operations.",
		Check: checkZC1202,
		Fix:   fixZC1202,
	})
}

// fixZC1202 rewrites `ifconfig` to `ip addr` at the command name
// position. `ip addr` is the closest single-token-equivalent iproute2
// invocation; arguments stay untouched and operators/flags must be
// adjusted manually for non-trivial cases.
func fixZC1202(node ast.Node, v Violation, _ []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ifconfig" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("ifconfig"),
		Replace: "ip addr",
	}}
}

func checkZC1202(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ifconfig" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1202",
		Message: "Avoid `ifconfig` — it is deprecated on modern Linux. " +
			"Use `ip addr`, `ip link`, or `ip route` from iproute2.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1203",
		Title:    "Avoid `netstat` — use `ss` for socket statistics",
		Severity: SeverityInfo,
		Description: "`netstat` is deprecated on modern Linux in favor of `ss` from iproute2. " +
			"`ss` is faster and provides more detailed socket information.",
		Check: checkZC1203,
		Fix:   fixZC1203,
	})
}

// fixZC1203 rewrites `netstat` to `ss` at the command name position.
// Single replacement — arguments stay untouched. The two tools share
// most short flags (`-t`, `-u`, `-l`, `-n`) so the swap is sound for
// the common cases; exotic netstat-only flags need manual review.
func fixZC1203(node ast.Node, v Violation, _ []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "netstat" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("netstat"),
		Replace: "ss",
	}}
}

func checkZC1203(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "netstat" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1203",
		Message: "Avoid `netstat` — it is deprecated on modern Linux. " +
			"Use `ss` from iproute2 for faster, more detailed socket statistics.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1204",
		Title:    "Avoid `route` — use `ip route` for routing",
		Severity: SeverityInfo,
		Description: "`route` is deprecated on modern Linux in favor of `ip route` from iproute2. " +
			"`ip route` provides consistent syntax with other `ip` subcommands.",
		Check: checkZC1204,
	})
}

func checkZC1204(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "route" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1204",
		Message: "Avoid `route` — it is deprecated on modern Linux. " +
			"Use `ip route` from iproute2 for consistent routing management.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1205",
		Title:    "Avoid `arp` — use `ip neigh` for neighbor tables",
		Severity: SeverityInfo,
		Description: "`arp` is deprecated on modern Linux in favor of `ip neigh` from iproute2. " +
			"`ip neigh` provides consistent syntax with other `ip` subcommands.",
		Check: checkZC1205,
	})
}

func checkZC1205(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "arp" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1205",
		Message: "Avoid `arp` — it is deprecated on modern Linux. " +
			"Use `ip neigh` from iproute2 for neighbor table management.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1206",
		Title:    "Avoid `crontab -e` in scripts — use `crontab file`",
		Severity: SeverityWarning,
		Description: "`crontab -e` opens an interactive editor which hangs in scripts. " +
			"Use `crontab file` or pipe content with `crontab -` for programmatic cron management.",
		Check: checkZC1206,
	})
}

func checkZC1206(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "crontab" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-e" {
			return []Violation{{
				KataID: "ZC1206",
				Message: "Avoid `crontab -e` in scripts — it opens an interactive editor. " +
					"Use `crontab file` or `echo '...' | crontab -` for programmatic cron management.",
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
		ID:       "ZC1207",
		Title:    "Avoid `passwd` in scripts — use `chpasswd`",
		Severity: SeverityWarning,
		Description: "`passwd` prompts interactively for password input. " +
			"Use `chpasswd` or `usermod --password` for non-interactive password changes.",
		Check: checkZC1207,
	})
}

func checkZC1207(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "passwd" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1207",
		Message: "Avoid `passwd` in scripts — it prompts interactively. " +
			"Use `chpasswd` or `usermod --password` for non-interactive password changes.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1208",
		Title:    "Avoid `visudo` in scripts — use sudoers.d drop-in files",
		Severity: SeverityWarning,
		Description: "`visudo` opens an interactive editor. For programmatic sudoers changes, " +
			"write to `/etc/sudoers.d/` drop-in files with `visudo -c` for validation.",
		Check: checkZC1208,
	})
}

func checkZC1208(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "visudo" {
		return nil
	}

	// visudo -c (check) is fine — it's non-interactive validation
	for _, arg := range cmd.Arguments {
		if arg.String() == "-c" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1208",
		Message: "Avoid `visudo` in scripts — it opens an interactive editor. " +
			"Write to `/etc/sudoers.d/` drop-in files and validate with `visudo -c`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1209",
		Title:    "Use `systemctl --no-pager` in scripts",
		Severity: SeverityStyle,
		Description: "`systemctl` invokes a pager by default which hangs in non-interactive scripts. " +
			"Use `--no-pager` or pipe to `cat` for reliable script output.",
		Check: checkZC1209,
		Fix:   fixZC1209,
	})
}

// fixZC1209 inserts ` --no-pager` after the `systemctl` command
// name so subcommands that emit pager output (status, list-*)
// behave predictably in scripts.
func fixZC1209(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("systemctl") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1209(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " --no-pager",
	}}
}

func offsetLineColZC1209(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1209(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "--no-pager" {
			return nil
		}
	}

	// Only flag subcommands that produce output (status, list-units, etc.)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "status" || val == "list-units" || val == "list-timers" || val == "show" {
			return []Violation{{
				KataID: "ZC1209",
				Message: "Use `systemctl --no-pager` in scripts. Without it, " +
					"systemctl invokes a pager that hangs in non-interactive execution.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1210",
		Title:    "Use `journalctl --no-pager` in scripts",
		Severity: SeverityStyle,
		Description: "`journalctl` invokes a pager by default which hangs in non-interactive scripts. " +
			"Use `--no-pager` for reliable script output.",
		Check: checkZC1210,
		Fix:   fixZC1210,
	})
}

// fixZC1210 inserts ` --no-pager` after the `journalctl` command
// name, preventing the pager from hanging in non-interactive runs.
// Mirrors ZC1209's insertion for `systemctl`.
func fixZC1210(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "journalctl" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("journalctl") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1210(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " --no-pager",
	}}
}

func offsetLineColZC1210(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1210(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "journalctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "--no-pager" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1210",
		Message: "Use `journalctl --no-pager` in scripts. Without it, " +
			"journalctl invokes a pager that hangs in non-interactive execution.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1211",
		Title:    "Use `git stash push -m` instead of bare `git stash`",
		Severity: SeverityStyle,
		Description: "Bare `git stash` creates unnamed stashes that are hard to identify later. " +
			"Use `git stash push -m 'description'` for self-documenting stashes.",
		Check: checkZC1211,
	})
}

func checkZC1211(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 1 {
		return nil
	}

	subCmd := cmd.Arguments[0].String()
	if subCmd != "stash" {
		return nil
	}

	// git stash push -m is fine, git stash pop/apply/list/drop are fine
	if len(cmd.Arguments) >= 2 {
		action := cmd.Arguments[1].String()
		if action == "push" || action == "pop" || action == "apply" ||
			action == "list" || action == "drop" || action == "show" {
			return nil
		}
	}

	// Bare "git stash" with no subcommand
	if len(cmd.Arguments) == 1 {
		return []Violation{{
			KataID: "ZC1211",
			Message: "Use `git stash push -m 'description'` instead of bare `git stash`. " +
				"Named stashes are easier to identify and manage.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1212",
		Title:    "Avoid `git add .` — use explicit paths or `git add -p`",
		Severity: SeverityInfo,
		Description: "`git add .` stages everything including unintended files. " +
			"Use explicit file paths or `git add -p` for selective staging.",
		Check: checkZC1212,
	})
}

func checkZC1212(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}

	if cmd.Arguments[0].String() != "add" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "." || val == "-A" {
			return []Violation{{
				KataID: "ZC1212",
				Message: "Avoid `git add .` or `git add -A` — they stage everything including " +
					"unintended files. Use explicit paths or `git add -p` for selective staging.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1213",
		Title:    "Use `apt-get -y` in scripts for non-interactive installs",
		Severity: SeverityWarning,
		Description: "`apt-get install` without `-y` prompts for confirmation which hangs scripts. " +
			"Use `-y` or set `DEBIAN_FRONTEND=noninteractive` for unattended installs.",
		Check: checkZC1213,
		Fix:   fixZC1213,
	})
}

// fixZC1213 inserts ` -y` after `apt-get` so install / upgrade /
// dist-upgrade run without interactive confirmation. Detector
// already guards the shape (install-class subcommand + no -y).
func fixZC1213(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "apt-get" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("apt-get") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1213(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -y",
	}}
}

func offsetLineColZC1213(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1213(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "apt-get" {
		return nil
	}

	hasInstall := false
	hasYes := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "install" || val == "upgrade" || val == "dist-upgrade" {
			hasInstall = true
		}
		if val == "-y" || val == "--yes" || val == "-qq" {
			hasYes = true
		}
	}

	if hasInstall && !hasYes {
		return []Violation{{
			KataID: "ZC1213",
			Message: "Use `apt-get -y` in scripts. Without `-y`, apt-get prompts for confirmation " +
				"which hangs in non-interactive execution.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1214",
		Title:    "Avoid `su` in scripts — use `sudo -u` for user switching",
		Severity: SeverityWarning,
		Description: "`su` prompts for a password interactively which hangs scripts. " +
			"Use `sudo -u user cmd` for non-interactive privilege switching.",
		Check: checkZC1214,
	})
}

func checkZC1214(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "su" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1214",
		Message: "Avoid `su` in scripts — it prompts for a password interactively. " +
			"Use `sudo -u user cmd` for non-interactive privilege switching.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1215",
		Title:    "Source `/etc/os-release` instead of parsing with `cat`/`grep`",
		Severity: SeverityStyle,
		Description: "`/etc/os-release` is designed to be sourced directly. " +
			"Use `. /etc/os-release` to get variables like `$ID`, `$VERSION_ID` without parsing.",
		Check: checkZC1215,
		Fix:   fixZC1215,
	})
}

// fixZC1215 rewrites `cat /etc/os-release` (or `/etc/lsb-release`) to
// `. /etc/os-release`. Single-edit replacement of the `cat` command
// name with the source builtin `.`. Only fires when cat has exactly
// one argument; piped or multi-file shapes are left alone. Idempotent
// — a re-run sees `.`, not `cat`. Defensive byte-match guard.
func fixZC1215(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}
	if len(cmd.Arguments) != 1 {
		return nil
	}
	val := cmd.Arguments[0].String()
	if val != "/etc/os-release" && val != "/etc/lsb-release" {
		return nil
	}
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 || off+len("cat") > len(source) {
		return nil
	}
	if string(source[off:off+len("cat")]) != "cat" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("cat"),
		Replace: ".",
	}}
}

func checkZC1215(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/etc/os-release" || val == "/etc/lsb-release" {
			return []Violation{{
				KataID: "ZC1215",
				Message: "Source `/etc/os-release` directly with `. /etc/os-release` instead of " +
					"parsing with `cat`. It exports variables like `$ID` and `$VERSION_ID`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1216",
		Title:    "Avoid `nslookup` — use `dig` or `host` for DNS queries",
		Severity: SeverityInfo,
		Description: "`nslookup` is deprecated in many distributions. " +
			"`dig` provides more detailed output and `host` is simpler for basic lookups.",
		Check: checkZC1216,
		Fix:   fixZC1216,
	})
}

// fixZC1216 rewrites `nslookup` to `host` at the command name position.
// `host <name>` matches the most common `nslookup <name>` invocation;
// arguments stay untouched and exotic nslookup-only flags need manual
// review.
func fixZC1216(node ast.Node, v Violation, _ []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "nslookup" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("nslookup"),
		Replace: "host",
	}}
}

func checkZC1216(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "nslookup" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1216",
		Message: "Avoid `nslookup` — it is deprecated on many systems. " +
			"Use `dig` for detailed DNS queries or `host` for simple lookups.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1217",
		Title:    "Avoid `service` command — use `systemctl` on systemd",
		Severity: SeverityInfo,
		Description: "`service` is a SysVinit compatibility wrapper. " +
			"On systemd systems, use `systemctl start/stop/restart/status` directly.",
		Check: checkZC1217,
		// Reuse the `service UNIT VERB` → `systemctl VERB UNIT` rewrite
		// from ZC1512. Both detectors fire on the same shape; the
		// conflict resolver dedupes overlapping edits.
		Fix: fixZC1512,
	})
}

func checkZC1217(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "service" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1217",
		Message: "Avoid `service` — it is a SysVinit compatibility wrapper. " +
			"Use `systemctl` directly on systemd systems.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1218",
		Title:    "Avoid `useradd` without `--shell /sbin/nologin` for service accounts",
		Severity: SeverityWarning,
		Description: "Service accounts created with `useradd` should use `--shell /sbin/nologin` " +
			"and `--system` to prevent interactive login and use system UID ranges.",
		Check: checkZC1218,
	})
}

func checkZC1218(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "useradd" {
		return nil
	}

	hasSystem := false
	hasNologin := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "--system" || val == "-r" {
			hasSystem = true
		}
		if val == "/sbin/nologin" || val == "/usr/sbin/nologin" || val == "/bin/false" {
			hasNologin = true
		}
	}

	if hasSystem && !hasNologin {
		return []Violation{{
			KataID: "ZC1218",
			Message: "Add `--shell /sbin/nologin` when creating system accounts with `useradd`. " +
				"This prevents interactive login for service accounts.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1219",
		Title:    "Use `curl -fsSL` instead of `wget -O -` for piped downloads",
		Severity: SeverityStyle,
		Description: "`wget -O -` outputs to stdout but lacks `curl`'s error handling. " +
			"`curl -fsSL` fails on HTTP errors, is silent, follows redirects, and is more portable.",
		Check: checkZC1219,
		Fix:   fixZC1219,
	})
}

// fixZC1219 collapses `wget -O- URL` / `wget -qO- URL` into
// `curl -fsSL URL`. The span covers the `wget` command name and the
// `-O-`/`-qO-` flag in a single edit so the rewrite stays deterministic
// even if a separate kata also fires on the `wget` name; trailing URL
// argument(s) stay in place.
func fixZC1219(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wget" {
		return nil
	}
	var flag ast.Expression
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-O-" || val == "-qO-" {
			flag = arg
			break
		}
	}
	if flag == nil {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("wget") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("wget")]) != "wget" {
		return nil
	}
	flagTok := flag.TokenLiteralNode()
	flagOff := LineColToByteOffset(source, flagTok.Line, flagTok.Column)
	flagLen := len(flag.String())
	if flagOff < 0 || flagOff+flagLen > len(source) {
		return nil
	}
	if string(source[flagOff:flagOff+flagLen]) != flag.String() {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  flagOff + flagLen - nameOff,
		Replace: "curl -fsSL",
	}}
}

func checkZC1219(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wget" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-O-" || val == "-qO-" {
			return []Violation{{
				KataID: "ZC1219",
				Message: "Use `curl -fsSL` instead of `wget -O -` for piped downloads. " +
					"`curl` fails on HTTP errors and is available on more platforms.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1220",
		Title:    "Use `chown :group` instead of `chgrp` for group changes",
		Severity: SeverityStyle,
		Description: "`chgrp` is redundant when `chown :group file` does the same thing. " +
			"Using `chown` for both user and group changes is more consistent.",
		Check: checkZC1220,
	})
}

func checkZC1220(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chgrp" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1220",
		Message: "Use `chown :group file` instead of `chgrp group file`. " +
			"`chown` handles both user and group changes consistently.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1221",
		Title:    "Avoid `fdisk` in scripts — use `parted` or `sfdisk`",
		Severity: SeverityWarning,
		Description: "`fdisk` is interactive and not scriptable. " +
			"Use `parted -s` or `sfdisk` for non-interactive disk partitioning.",
		Check: checkZC1221,
	})
}

func checkZC1221(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "fdisk" {
		return nil
	}

	// fdisk -l (list) is non-interactive and fine
	for _, arg := range cmd.Arguments {
		if arg.String() == "-l" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1221",
		Message: "Avoid `fdisk` in scripts — it is interactive. " +
			"Use `parted -s` or `sfdisk` for scriptable disk partitioning.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1222",
		Title:    "Avoid `lsof -i` for port checks — use `ss -tlnp`",
		Severity: SeverityStyle,
		Description: "`lsof -i` is slow and requires elevated permissions on some systems. " +
			"`ss -tlnp` is faster and part of the standard iproute2 toolkit.",
		Check: checkZC1222,
	})
}

func checkZC1222(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "lsof" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-i" {
			return []Violation{{
				KataID: "ZC1222",
				Message: "Use `ss -tlnp` instead of `lsof -i` for port checks. " +
					"`ss` is faster and doesn't require elevated permissions.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1223",
		Title:    "Avoid `ip addr show` piped to `grep` — use `ip -br addr`",
		Severity: SeverityStyle,
		Description: "`ip addr show | grep` parses verbose output. " +
			"`ip -br addr` provides machine-readable brief output without needing grep.",
		Check: checkZC1223,
	})
}

func checkZC1223(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ip" {
		return nil
	}

	hasAddr := false
	hasBrief := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "addr" || val == "address" {
			hasAddr = true
		}
		if val == "-br" || val == "-brief" {
			hasBrief = true
		}
	}

	if hasAddr && !hasBrief {
		return []Violation{{
			KataID: "ZC1223",
			Message: "Use `ip -br addr` for machine-readable output instead of " +
				"parsing `ip addr show` with grep or awk.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1224",
		Title:    "Avoid parsing `free` output — read `/proc/meminfo` directly",
		Severity: SeverityStyle,
		Description: "`free` output format varies across versions and locales. " +
			"Read `/proc/meminfo` directly for reliable memory information in scripts.",
		Check: checkZC1224,
	})
}

func checkZC1224(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "free" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1224",
		Message: "Avoid parsing `free` output — its format varies across versions. " +
			"Read `/proc/meminfo` directly for reliable memory information.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1225",
		Title:    "Avoid parsing `uptime` — read `/proc/uptime` directly",
		Severity: SeverityStyle,
		Description: "`uptime` output is human-readable and varies by locale. " +
			"Read `/proc/uptime` for machine-parseable uptime in seconds.",
		Check: checkZC1225,
	})
}

func checkZC1225(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "uptime" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1225",
		Message: "Avoid parsing `uptime` — its output varies by locale. " +
			"Read `/proc/uptime` for machine-parseable seconds since boot.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1226",
		Title:    "Use `dmesg -T` or `--time-format=iso` for readable timestamps",
		Severity: SeverityStyle,
		Description: "`dmesg` without `-T` shows raw kernel timestamps in seconds since boot. " +
			"Use `-T` for human-readable timestamps or `--time-format=iso` for ISO 8601.",
		Check: checkZC1226,
		Fix:   fixZC1226,
	})
}

// fixZC1226 inserts ` -T` after the `dmesg` command name. Mirrors
// other insertion-style fixes (ZC1012 / ZC1017 / ZC1170 / ZC1209).
func fixZC1226(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "dmesg" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("dmesg") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1226(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -T",
	}}
}

func offsetLineColZC1226(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1226(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "dmesg" {
		return nil
	}

	hasTimeFlag := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-T" || val == "-t" || val == "--ctime" || val == "--reltime" {
			hasTimeFlag = true
		}
	}

	if !hasTimeFlag && len(cmd.Arguments) > 0 {
		return []Violation{{
			KataID: "ZC1226",
			Message: "Use `dmesg -T` for human-readable timestamps instead of raw " +
				"kernel boot-seconds. Or use `--time-format=iso` for ISO 8601.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1227",
		Title:    "Use `curl -f` to fail on HTTP errors",
		Severity: SeverityWarning,
		Description: "`curl` without `-f` silently returns error pages (404, 500) as success. " +
			"Use `-f` or `--fail` to return exit code 22 on HTTP errors.",
		Check: checkZC1227,
		Fix:   fixZC1227,
	})
}

// fixZC1227 inserts ` -f` after the `curl` command name so HTTP
// errors translate into a non-zero exit code. Detector guards the
// shape (URL arg present, no existing -f/-fsSL/etc.).
func fixZC1227(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "curl" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("curl") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1227(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -f",
	}}
}

func offsetLineColZC1227(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1227(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "curl" {
		return nil
	}

	hasFail := false
	hasURL := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-f" || val == "-fsSL" || val == "-fSL" || val == "-fsS" {
			hasFail = true
		}
		if len(val) > 4 && (val[:4] == "http" || val[:5] == "https") {
			hasURL = true
		}
	}

	if hasURL && !hasFail {
		return []Violation{{
			KataID: "ZC1227",
			Message: "Use `curl -f` to fail on HTTP errors. Without `-f`, curl silently " +
				"returns error pages (404, 500) as if they were successful.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1228",
		Title:    "Avoid `ssh` without host key policy in scripts",
		Severity: SeverityWarning,
		Description: "`ssh` without `-o BatchMode=yes` or `-o StrictHostKeyChecking` prompts " +
			"interactively for host key verification, hanging non-interactive scripts.",
		Check: checkZC1228,
	})
}

func checkZC1228(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ssh" {
		return nil
	}

	hasBatchOrStrict := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "BatchMode=yes" || val == "StrictHostKeyChecking=accept-new" ||
			val == "StrictHostKeyChecking=no" || val == "StrictHostKeyChecking=yes" {
			hasBatchOrStrict = true
		}
	}

	if !hasBatchOrStrict {
		return []Violation{{
			KataID: "ZC1228",
			Message: "Use `ssh -o BatchMode=yes` or `-o StrictHostKeyChecking=accept-new` in scripts. " +
				"Without these, ssh may prompt interactively and hang.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1229",
		Title:    "Prefer `rsync` over `scp` for file transfers",
		Severity: SeverityStyle,
		Description: "`scp` uses a deprecated protocol and lacks delta transfer, resume, " +
			"and progress features. `rsync` is more efficient and reliable for scripts.",
		Check: checkZC1229,
	})
}

func checkZC1229(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "scp" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1229",
		Message: "Prefer `rsync -az` over `scp` for file transfers. " +
			"`rsync` supports delta transfers, resume, and is more efficient.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1230",
		Title:    "Use `ping -c N` in scripts to limit ping count",
		Severity: SeverityWarning,
		Description: "`ping` without `-c` runs indefinitely on Linux, hanging scripts. " +
			"Always specify `-c N` to limit the number of packets.",
		Check: checkZC1230,
		Fix:   fixZC1230,
	})
}

// fixZC1230 inserts ` -c 4` after the `ping` command name. Detector
// already guards against an existing `-c` / `-W` flag, so the
// insertion is safe and idempotent on a re-run.
func fixZC1230(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ping" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	if IdentLenAt(source, nameOff) != len("ping") {
		return nil
	}
	insertAt := nameOff + len("ping")
	insLine, insCol := offsetLineColZC1230(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -c 4",
	}}
}

func offsetLineColZC1230(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1230(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ping" {
		return nil
	}

	hasCount := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-c" || val == "-W" {
			hasCount = true
		}
	}

	if !hasCount {
		return []Violation{{
			KataID: "ZC1230",
			Message: "Use `ping -c N` in scripts. Without `-c`, ping runs " +
				"indefinitely on Linux and will hang the script.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1231",
		Title:    "Use `git clone --depth 1` for CI and build scripts",
		Severity: SeverityStyle,
		Description: "`git clone` without `--depth` downloads the entire history. " +
			"Use `--depth 1` in CI/build scripts where only the latest commit is needed.",
		Check: checkZC1231,
		Fix:   fixZC1231,
	})
}

// fixZC1231 inserts ` --depth 1` after the `clone` subcommand in
// `git clone …`. Mirrors ZC1234's subcommand-level insertion for
// docker run --rm.
func fixZC1231(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	cloneArg := cmd.Arguments[0]
	if cloneArg.String() != "clone" {
		return nil
	}
	tok := cloneArg.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 || off+5 > len(source) {
		return nil
	}
	if string(source[off:off+5]) != "clone" {
		return nil
	}
	insertAt := off + 5
	insLine, insCol := offsetLineColZC1231(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " --depth 1",
	}}
}

func offsetLineColZC1231(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1231(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}

	if cmd.Arguments[0].String() != "clone" {
		return nil
	}

	hasDepth := false
	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "--depth" || val == "--shallow-since" || val == "--single-branch" {
			hasDepth = true
		}
	}

	if !hasDepth {
		return []Violation{{
			KataID: "ZC1231",
			Message: "Consider `git clone --depth 1` in scripts. Full clones download " +
				"entire history which is unnecessary for builds and CI.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1232",
		Title:    "Avoid bare `pip install` — use `--user` or virtualenv",
		Severity: SeverityWarning,
		Description: "Bare `pip install` may modify system Python packages. " +
			"Use `pip install --user`, `pipx`, or a virtualenv to isolate dependencies.",
		Check: checkZC1232,
	})
}

func checkZC1232(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "pip" && ident.Value != "pip3" {
		return nil
	}

	hasInstall := false
	hasSafe := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "install" {
			hasInstall = true
		}
		if val == "--user" || val == "-t" || val == "--target" || val == "--prefix" {
			hasSafe = true
		}
	}

	if hasInstall && !hasSafe {
		return []Violation{{
			KataID: "ZC1232",
			Message: "Use `pip install --user` or a virtualenv instead of bare `pip install`. " +
				"System-wide pip installs can break OS package managers.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1233",
		Title:    "Avoid `npm install -g` — use `npx` for one-off tools",
		Severity: SeverityStyle,
		Description: "Global npm installs pollute the system. Use `npx` to run tools " +
			"without installing, or `npm install --save-dev` for project dependencies.",
		Check: checkZC1233,
	})
}

func checkZC1233(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "npm" {
		return nil
	}

	hasInstall := false
	hasGlobal := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "install" || val == "i" {
			hasInstall = true
		}
		if val == "-g" {
			hasGlobal = true
		}
	}

	if hasInstall && hasGlobal {
		return []Violation{{
			KataID: "ZC1233",
			Message: "Avoid `npm install -g`. Use `npx` for one-off tool execution " +
				"or `npm install --save-dev` for project-scoped dependencies.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1234",
		Title:    "Use `docker run --rm` to auto-remove containers",
		Severity: SeverityStyle,
		Description: "`docker run` without `--rm` leaves stopped containers behind. " +
			"Use `--rm` in scripts to automatically clean up after execution.",
		Check: checkZC1234,
		Fix:   fixZC1234,
	})
}

// fixZC1234 inserts ` --rm` after the `run` subcommand in a
// `docker run …` invocation. Detector has already verified the shape
// (docker + run + no --rm + no -d).
func fixZC1234(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	runArg := cmd.Arguments[0]
	if runArg.String() != "run" {
		return nil
	}
	tok := runArg.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 || off+3 > len(source) {
		return nil
	}
	if string(source[off:off+3]) != "run" {
		return nil
	}
	insertAt := off + 3
	insLine, insCol := offsetLineColZC1234(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " --rm",
	}}
}

func offsetLineColZC1234(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1234(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "run" {
		return nil
	}

	hasRM := false
	hasDetach := false

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "--rm" {
			hasRM = true
		}
		if val == "-d" {
			hasDetach = true
		}
	}

	if !hasRM && !hasDetach {
		return []Violation{{
			KataID: "ZC1234",
			Message: "Use `docker run --rm` to auto-remove containers after exit. " +
				"Without `--rm`, stopped containers accumulate on disk.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1235",
		Title:    "Use `git push --force-with-lease` instead of `--force`",
		Severity: SeverityWarning,
		Description: "`git push --force` overwrites remote history unconditionally. " +
			"`--force-with-lease` is safer as it fails if the remote has changed.",
		Check: checkZC1235,
		Fix:   fixZC1235,
	})
}

// fixZC1235 rewrites `git push -f` to `git push --force-with-lease`.
// Single-edit replacement of the `-f` flag at its argument position;
// surrounding subcommand and refspec arguments stay in place.
func fixZC1235(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "push" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		if arg.String() != "-f" {
			continue
		}
		tok := arg.TokenLiteralNode()
		off := LineColToByteOffset(source, tok.Line, tok.Column)
		if off < 0 || off+2 > len(source) {
			return nil
		}
		if string(source[off:off+2]) != "-f" {
			return nil
		}
		return []FixEdit{{
			Line:    tok.Line,
			Column:  tok.Column,
			Length:  2,
			Replace: "--force-with-lease",
		}}
	}
	return nil
}

func checkZC1235(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "push" {
		return nil
	}

	hasForce := false
	hasFWL := false

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "-f" {
			hasForce = true
		}
	}

	if hasForce && !hasFWL {
		return []Violation{{
			KataID: "ZC1235",
			Message: "Use `git push --force-with-lease` instead of `-f`/`--force`. " +
				"It prevents overwriting remote changes made by others.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1236",
		Title:    "Avoid `git reset --hard` — irreversible data loss risk",
		Severity: SeverityWarning,
		Description: "`git reset --hard` discards all uncommitted changes irreversibly. " +
			"Use `git stash` to save changes first, or `git reset --soft` to keep them staged.",
		Check: checkZC1236,
	})
}

func checkZC1236(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "reset" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "--hard" {
			return []Violation{{
				KataID: "ZC1236",
				Message: "Avoid `git reset --hard` — it permanently discards uncommitted changes. " +
					"Use `git stash` first, or `git reset --soft` to keep changes staged.",
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
		ID:       "ZC1237",
		Title:    "Use `git clean -n` before `git clean -fd`",
		Severity: SeverityWarning,
		Description: "`git clean -fd` permanently deletes untracked files and directories. " +
			"Use `-n` (dry run) first to preview what will be removed.",
		Check: checkZC1237,
	})
}

func checkZC1237(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "clean" {
		return nil
	}

	hasForce := false
	hasDryRun := false

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "-f" || val == "-fd" || val == "-df" || val == "-fdx" {
			hasForce = true
		}
		if val == "-n" {
			hasDryRun = true
		}
	}

	if hasForce && !hasDryRun {
		return []Violation{{
			KataID: "ZC1237",
			Message: "Use `git clean -n` first to preview removals before `git clean -fd`. " +
				"Forced clean permanently deletes untracked files.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1238",
		Title:    "Avoid `docker exec -it` in scripts — drop `-it` for non-interactive",
		Severity: SeverityWarning,
		Description: "`docker exec -it` allocates a TTY and attaches stdin, which hangs " +
			"in non-interactive scripts. Use `docker exec` without `-it` for scripted commands.",
		Check: checkZC1238,
		Fix:   fixZC1238,
	})
}

// fixZC1238 strips the `-it` (or `-ti`) flag from a `docker exec`
// invocation. The span covers the leading whitespace plus the flag
// token so the surrounding source stays byte-identical.
func fixZC1238(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}
	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "exec" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		if v := arg.String(); v == "-it" || v == "-ti" {
			return zc1238StripFlag(source, arg, v)
		}
	}
	return nil
}

// zc1238StripFlag deletes the flag arg plus the run of horizontal
// whitespace immediately preceding it; the leading space the user
// typed disappears with the flag, leaving `docker exec CMD`.
func zc1238StripFlag(source []byte, arg ast.Expression, lit string) []FixEdit {
	tok := arg.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 || off+len(lit) > len(source) {
		return nil
	}
	if string(source[off:off+len(lit)]) != lit {
		return nil
	}
	start := off
	for start > 0 && (source[start-1] == ' ' || source[start-1] == '\t') {
		start--
	}
	end := off + len(lit)
	startLine, startCol := offsetLineColZC1238(source, start)
	if startLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    startLine,
		Column:  startCol,
		Length:  end - start,
		Replace: "",
	}}
}

func offsetLineColZC1238(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1238(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "exec" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "-it" || val == "-ti" {
			return []Violation{{
				KataID: "ZC1238",
				Message: "Avoid `docker exec -it` in scripts — TTY allocation hangs without a terminal. " +
					"Use `docker exec` without `-it` for non-interactive commands.",
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
		ID:       "ZC1239",
		Title:    "Avoid `kubectl exec -it` in scripts",
		Severity: SeverityWarning,
		Description: "`kubectl exec -it` allocates a TTY which hangs in non-interactive scripts. " +
			"Use `kubectl exec` without `-it` or use `kubectl exec -- cmd` for scripted commands.",
		Check: checkZC1239,
		Fix:   fixZC1239,
	})
}

// fixZC1239 strips the `-it` (or `-ti`) flag from a `kubectl exec`
// invocation. Reuses the token-strip helper introduced for ZC1238.
func fixZC1239(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubectl" {
		return nil
	}
	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "exec" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		if v := arg.String(); v == "-it" || v == "-ti" {
			return zc1238StripFlag(source, arg, v)
		}
	}
	return nil
}

func checkZC1239(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubectl" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "exec" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "-it" || val == "-ti" {
			return []Violation{{
				KataID: "ZC1239",
				Message: "Avoid `kubectl exec -it` in scripts — TTY allocation hangs without a terminal. " +
					"Use `kubectl exec pod -- cmd` for non-interactive execution.",
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
		ID:       "ZC1240",
		Title:    "Use `find -maxdepth` with `-delete` to limit scope",
		Severity: SeverityWarning,
		Description: "`find -delete` without `-maxdepth` recurses infinitely and may " +
			"delete more than intended. Always limit the search depth.",
		Check: checkZC1240,
	})
}

func checkZC1240(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	hasDelete := false
	hasMaxdepth := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-delete" {
			hasDelete = true
		}
		if val == "-maxdepth" {
			hasMaxdepth = true
		}
	}

	if hasDelete && !hasMaxdepth {
		return []Violation{{
			KataID: "ZC1240",
			Message: "Use `find -maxdepth N` with `-delete` to limit deletion scope. " +
				"Without depth limits, find recurses infinitely.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1241",
		Title:    "Use `xargs -0` with null separators for safe argument passing",
		Severity: SeverityWarning,
		Description: "`xargs` without `-0` splits on whitespace, breaking on filenames with spaces. " +
			"Use `xargs -0` paired with `find -print0` for safe handling.",
		Check: checkZC1241,
		Fix:   fixZC1241,
	})
}

// fixZC1241 inserts ` -0` after the `xargs` command name so
// null-terminated input from `find -print0` is consumed safely.
// Detector gates on `xargs rm` without `-0`.
func fixZC1241(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "xargs" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("xargs") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1241(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -0",
	}}
}

func offsetLineColZC1241(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1241(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "xargs" {
		return nil
	}

	hasNull := false
	hasRM := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-0" {
			hasNull = true
		}
		if val == "rm" {
			hasRM = true
		}
	}

	if hasRM && !hasNull {
		return []Violation{{
			KataID: "ZC1241",
			Message: "Use `xargs -0 rm` with `find -print0` for safe deletion. " +
				"Without `-0`, filenames with spaces or special characters break.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1242",
		Title:    "Use `tar -C dir` to extract into a specific directory",
		Severity: SeverityInfo,
		Description: "`tar xf` without `-C` extracts into the current directory which may " +
			"overwrite files unexpectedly. Use `-C dir` to control the extraction target.",
		Check: checkZC1242,
	})
}

func checkZC1242(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tar" {
		return nil
	}

	hasExtract := false
	hasTarget := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-x" || val == "xf" || val == "xzf" || val == "xjf" || val == "xJf" {
			hasExtract = true
		}
		if val == "-C" {
			hasTarget = true
		}
	}

	if hasExtract && !hasTarget {
		return []Violation{{
			KataID: "ZC1242",
			Message: "Use `tar -C dir` to specify extraction directory. " +
				"Without `-C`, tar extracts into the current directory which may overwrite files.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1243",
		Title:    "Use `grep -lZ` with `xargs -0` for safe file lists",
		Severity: SeverityWarning,
		Description: "`grep -l` outputs one filename per line, breaking on names with newlines. " +
			"Use `grep -lZ` (null-terminated) paired with `xargs -0` for safe processing.",
		Check: checkZC1243,
	})
}

func checkZC1243(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	hasListFiles := false
	hasNullTerm := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-l" || val == "-rl" || val == "-lr" {
			hasListFiles = true
		}
		if val == "-Z" || val == "-lZ" || val == "-Zl" {
			hasNullTerm = true
		}
	}

	if hasListFiles && !hasNullTerm {
		return []Violation{{
			KataID: "ZC1243",
			Message: "Use `grep -lZ` instead of `grep -l` for null-terminated file lists. " +
				"Pair with `xargs -0` to safely handle filenames with special characters.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1244",
		Title:    "Consider `mv -n` to prevent overwriting existing files",
		Severity: SeverityInfo,
		Description: "`mv` overwrites existing files without warning by default. " +
			"Use `-n` (no-clobber) to prevent accidental overwrites in scripts.",
		Check: checkZC1244,
	})
}

func checkZC1244(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mv" {
		return nil
	}

	hasSafe := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-n" || val == "-i" || val == "-f" {
			hasSafe = true
		}
	}

	if !hasSafe && len(cmd.Arguments) >= 2 {
		return []Violation{{
			KataID: "ZC1244",
			Message: "Consider `mv -n` to prevent overwriting existing files. " +
				"Without `-n`, `mv` silently overwrites the target.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1245",
		Title:    "Avoid disabling TLS certificate verification",
		Severity: SeverityError,
		Description: "Flags like `--no-check-certificate` (wget) or `-k`/`--insecure` (curl) " +
			"disable TLS verification, making connections vulnerable to MITM attacks.",
		Check: checkZC1245,
	})
}

func checkZC1245(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value

	for _, arg := range cmd.Arguments {
		val := arg.String()

		if name == "curl" && (val == "-k" || val == "--insecure") {
			return []Violation{{
				KataID: "ZC1245",
				Message: "Avoid `curl -k`/`--insecure` — it disables TLS certificate verification. " +
					"Fix the certificate chain or use `--cacert` to specify a CA bundle.",
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
		ID:       "ZC1246",
		Title:    "Avoid hardcoded passwords in command arguments",
		Severity: SeverityError,
		Description: "Passing passwords as command arguments exposes them in process lists " +
			"and shell history. Use environment variables or credential files instead.",
		Check: checkZC1246,
	})
}

func checkZC1246(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value
	if name != "mysql" && name != "psql" && name != "mongosh" && name != "redis-cli" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if strings.HasPrefix(val, "-p") && len(val) > 2 && val != "-p" {
			return []Violation{{
				KataID: "ZC1246",
				Message: "Avoid passing passwords as command arguments — they appear in process lists. " +
					"Use environment variables (e.g., `MYSQL_PWD`) or credential files instead.",
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
		ID:       "ZC1247",
		Title:    "Avoid `chmod +s` — setuid/setgid bits are security risks",
		Severity: SeverityError,
		Description: "Setting the setuid or setgid bit (`chmod +s` or `chmod u+s`) allows " +
			"files to execute with the owner's privileges, creating privilege escalation risks.",
		Check: checkZC1247,
	})
}

func checkZC1247(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chmod" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "+s" || val == "u+s" || val == "g+s" || val == "4755" || val == "2755" {
			return []Violation{{
				KataID: "ZC1247",
				Message: "Avoid `chmod +s` — setuid/setgid bits create privilege escalation risks. " +
					"Use `sudo`, capabilities (`setcap`), or dedicated service accounts instead.",
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
		ID:       "ZC1248",
		Title:    "Prefer `ufw`/`firewalld` over raw `iptables`",
		Severity: SeverityInfo,
		Description: "Raw `iptables` rules are complex and non-persistent by default. " +
			"Use `ufw` (Ubuntu) or `firewalld` (RHEL) for manageable, persistent firewall rules.",
		Check: checkZC1248,
	})
}

func checkZC1248(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "iptables" && ident.Value != "ip6tables" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1248",
		Message: "Prefer `ufw` or `firewalld` over raw `iptables`. " +
			"Firewall frontends provide persistent, manageable rules.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1249",
		Title:    "Use `ssh-keygen -f` to specify key file in scripts",
		Severity: SeverityWarning,
		Description: "`ssh-keygen` without `-f` prompts for a file path interactively. " +
			"Use `-f /path/to/key` and `-N ''` for non-interactive key generation.",
		Check: checkZC1249,
	})
}

func checkZC1249(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ssh-keygen" {
		return nil
	}

	hasFile := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-f" {
			hasFile = true
		}
	}

	if !hasFile {
		return []Violation{{
			KataID: "ZC1249",
			Message: "Use `ssh-keygen -f /path/to/key -N ''` in scripts. " +
				"Without `-f`, ssh-keygen prompts interactively for the file path.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1250",
		Title:    "Use `gpg --batch` in scripts for non-interactive operation",
		Severity: SeverityWarning,
		Description: "`gpg` without `--batch` may prompt for passphrases or confirmations. " +
			"Use `--batch` and `--yes` for fully non-interactive GPG operations in scripts.",
		Check: checkZC1250,
	})
}

func checkZC1250(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gpg" {
		return nil
	}

	hasBatch := false
	hasOperation := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-b" {
			hasBatch = true
		}
		if val == "-d" || val == "-s" || val == "-e" || val == "-c" {
			hasOperation = true
		}
	}

	if hasOperation && !hasBatch {
		return []Violation{{
			KataID: "ZC1250",
			Message: "Use `gpg --batch` in scripts for non-interactive operation. " +
				"Without `--batch`, gpg may prompt for passphrases or confirmations.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1251",
		Title:    "Use `mount -o noexec,nosuid` for untrusted media",
		Severity: SeverityWarning,
		Description: "Mounting untrusted filesystems without `noexec,nosuid` allows execution " +
			"of malicious binaries and setuid exploits. Always restrict mount options.",
		Check: checkZC1251,
	})
}

func checkZC1251(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mount" {
		return nil
	}

	hasOptions := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-o" {
			hasOptions = true
		}
	}

	// Only flag mount with device arguments but no -o options
	hasDevice := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 4 && val[:4] == "/dev" {
			hasDevice = true
		}
	}

	if hasDevice && !hasOptions {
		return []Violation{{
			KataID: "ZC1251",
			Message: "Use `mount -o noexec,nosuid,nodev` when mounting external media. " +
				"Without restrictions, mounted filesystems can contain executable exploits.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1252",
		Title:    "Use `getent passwd` instead of `cat /etc/passwd`",
		Severity: SeverityStyle,
		Description: "`cat /etc/passwd` misses users from LDAP, NIS, or SSSD sources. " +
			"`getent passwd` queries NSS and returns all configured user databases.",
		Check: checkZC1252,
		Fix:   fixZC1252,
	})
}

// fixZC1252 rewrites `cat /etc/{passwd,group,shadow}` to
// `getent {passwd,group,shadow}`. Two edits per fire: the command
// name and the file argument. Only fires when the cat command has
// exactly one argument; piped or multi-file shapes are left alone
// (ZC1146 handles `cat FILE | tool`, and multi-file `cat` doesn't
// translate cleanly to `getent`). Idempotent: a re-run sees `getent`,
// not `cat`, so the detector won't fire.
func fixZC1252(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}
	if len(cmd.Arguments) != 1 {
		return nil
	}
	arg := cmd.Arguments[0]
	val := arg.String()
	var dbName string
	switch val {
	case "/etc/passwd":
		dbName = "passwd"
	case "/etc/group":
		dbName = "group"
	case "/etc/shadow":
		dbName = "shadow"
	default:
		return nil
	}
	catOff := LineColToByteOffset(source, v.Line, v.Column)
	if catOff < 0 || catOff+len("cat") > len(source) {
		return nil
	}
	if string(source[catOff:catOff+len("cat")]) != "cat" {
		return nil
	}
	argTok := arg.TokenLiteralNode()
	argOff := LineColToByteOffset(source, argTok.Line, argTok.Column)
	if argOff < 0 || argOff+len(val) > len(source) {
		return nil
	}
	if string(source[argOff:argOff+len(val)]) != val {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: len("cat"), Replace: "getent"},
		{Line: argTok.Line, Column: argTok.Column, Length: len(val), Replace: dbName},
	}
}

func checkZC1252(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/etc/passwd" || val == "/etc/group" || val == "/etc/shadow" {
			return []Violation{{
				KataID: "ZC1252",
				Message: "Use `getent` instead of `cat " + val + "`. " +
					"`getent` queries all NSS sources including LDAP and SSSD.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1253",
		Title:    "Use `docker build --no-cache` in CI for reproducible builds",
		Severity: SeverityStyle,
		Description: "`docker build` uses layer caching which can mask dependency changes. " +
			"Use `--no-cache` in CI pipelines to ensure fully reproducible builds.",
		Check: checkZC1253,
		Fix:   fixZC1253,
	})
}

// fixZC1253 inserts ` --no-cache` after the `build` subcommand in
// `docker build …`. Mirrors ZC1234's subcommand-level insertion.
func fixZC1253(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	subArg := cmd.Arguments[0]
	if subArg.String() != "build" {
		return nil
	}
	tok := subArg.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 || off+5 > len(source) {
		return nil
	}
	if string(source[off:off+5]) != "build" {
		return nil
	}
	insertAt := off + 5
	insLine, insCol := offsetLineColZC1253(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " --no-cache",
	}}
}

func offsetLineColZC1253(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1253(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "build" {
		return nil
	}

	hasNoCache := false
	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "--no-cache" {
			hasNoCache = true
		}
	}

	if !hasNoCache {
		return []Violation{{
			KataID: "ZC1253",
			Message: "Consider `docker build --no-cache` in CI for reproducible builds. " +
				"Layer caching can mask changed dependencies.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1254",
		Title:    "Avoid `git commit --amend` in shared branches",
		Severity: SeverityWarning,
		Description: "`git commit --amend` rewrites the last commit which causes problems " +
			"if already pushed. Use `git commit --fixup` or a new commit instead.",
		Check: checkZC1254,
	})
}

func checkZC1254(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "commit" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "--amend" {
			return []Violation{{
				KataID: "ZC1254",
				Message: "Avoid `git commit --amend` on shared branches — it rewrites history. " +
					"Use `git commit --fixup` or create a new commit instead.",
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
		ID:       "ZC1255",
		Title:    "Use `curl -L` to follow HTTP redirects",
		Severity: SeverityInfo,
		Description: "`curl` without `-L` does not follow redirects, returning 301/302 responses " +
			"instead of the actual content. Use `-L` to follow redirects automatically.",
		Check: checkZC1255,
		Fix:   fixZC1255,
	})
}

// fixZC1255 inserts ` -L` after the `curl` command name. Detector
// already guards against any existing follow-redirect flag, so the
// insertion is idempotent on a re-run.
func fixZC1255(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "curl" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	if IdentLenAt(source, nameOff) != len("curl") {
		return nil
	}
	insertAt := nameOff + len("curl")
	insLine, insCol := offsetLineColZC1255(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -L",
	}}
}

func offsetLineColZC1255(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1255(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "curl" {
		return nil
	}

	hasFollow := false
	hasURL := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-L" || val == "-fsSL" || val == "-fSL" || val == "-sL" {
			hasFollow = true
		}
		if len(val) > 7 && val[:5] == "https" {
			hasURL = true
		}
	}

	if hasURL && !hasFollow {
		return []Violation{{
			KataID: "ZC1255",
			Message: "Use `curl -L` to follow HTTP redirects. Without `-L`, curl returns " +
				"redirect responses (301/302) instead of the actual content.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1256",
		Title:    "Clean up `mkfifo` pipes with a trap on EXIT",
		Severity: SeverityInfo,
		Description: "`mkfifo` creates named pipes that persist on the filesystem. " +
			"Set up a `trap` to remove them on EXIT to prevent leftover files.",
		Check: checkZC1256,
	})
}

func checkZC1256(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mkfifo" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1256",
		Message: "Set up `trap 'rm -f pipe' EXIT` after `mkfifo`. " +
			"Named pipes persist on the filesystem and need explicit cleanup.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1257",
		Title:    "Use `docker stop -t` to set graceful shutdown timeout",
		Severity: SeverityStyle,
		Description: "`docker stop` defaults to 10s before SIGKILL. In CI scripts, " +
			"set an explicit timeout with `-t` to control shutdown behavior.",
		Check: checkZC1257,
		Fix:   fixZC1257,
	})
}

// fixZC1257 inserts ` -t 10` after the `stop` subcommand of a
// `docker stop …` invocation. Mirrors the subcommand-level pattern
// used by ZC1265 (`systemctl enable --now`).
func fixZC1257(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}
	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "stop" {
		return nil
	}
	stopArg := cmd.Arguments[0]
	tok := stopArg.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 || off+len("stop") > len(source) {
		return nil
	}
	if string(source[off:off+len("stop")]) != "stop" {
		return nil
	}
	insertAt := off + len("stop")
	insLine, insCol := offsetLineColZC1257(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -t 10",
	}}
}

func offsetLineColZC1257(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1257(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "stop" {
		return nil
	}

	hasTimeout := false
	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "-t" {
			hasTimeout = true
		}
	}

	if !hasTimeout {
		return []Violation{{
			KataID: "ZC1257",
			Message: "Use `docker stop -t N` to set an explicit shutdown timeout. " +
				"The default 10s may be too long or too short for your use case.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1258",
		Title:    "Consider `rsync --delete` for directory sync",
		Severity: SeverityWarning,
		Description: "`rsync` without `--delete` keeps files on the destination that were " +
			"removed from the source. Use `--delete` for true directory mirroring.",
		Check: checkZC1258,
	})
}

func checkZC1258(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "rsync" {
		return nil
	}

	hasDelete := false
	hasTrailingSlash := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "--delete" {
			hasDelete = true
		}
		if strings.HasSuffix(val, "/") && !strings.HasPrefix(val, "-") {
			hasTrailingSlash = true
		}
	}

	if hasTrailingSlash && !hasDelete {
		return []Violation{{
			KataID: "ZC1258",
			Message: "Consider `rsync --delete` for directory sync. Without `--delete`, " +
				"files removed from source remain on the destination.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1259",
		Title:    "Avoid `docker pull` without explicit tag — pin image versions",
		Severity: SeverityWarning,
		Description: "`docker pull image` without a tag defaults to `:latest` which is " +
			"mutable and non-reproducible. Always pin to a specific version tag or digest.",
		Check: checkZC1259,
	})
}

func checkZC1259(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	if len(cmd.Arguments) < 2 || cmd.Arguments[0].String() != "pull" {
		return nil
	}

	image := cmd.Arguments[1].String()
	if !strings.Contains(image, ":") && !strings.Contains(image, "@sha256") {
		return []Violation{{
			KataID: "ZC1259",
			Message: "Pin Docker image to a specific tag instead of defaulting to `:latest`. " +
				"Untagged pulls are non-reproducible and may break unexpectedly.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1260",
		Title:    "Use `git branch -d` instead of `-D` for safe deletion",
		Severity: SeverityWarning,
		Description: "`git branch -D` force-deletes branches even if unmerged. " +
			"Use `-d` which refuses to delete unmerged branches, preventing data loss.",
		Check: checkZC1260,
		Fix:   fixZC1260,
	})
}

// fixZC1260 rewrites `git branch -D` to `git branch -d`. Single-character
// flag swap at the `-D` argument position; surrounding subcommand and
// branch-name arguments stay in place.
func fixZC1260(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "branch" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		if arg.String() != "-D" {
			continue
		}
		tok := arg.TokenLiteralNode()
		off := LineColToByteOffset(source, tok.Line, tok.Column)
		if off < 0 || off+2 > len(source) {
			return nil
		}
		if string(source[off:off+2]) != "-D" {
			return nil
		}
		return []FixEdit{{
			Line:    tok.Line,
			Column:  tok.Column,
			Length:  2,
			Replace: "-d",
		}}
	}
	return nil
}

func checkZC1260(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "branch" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "-D" {
			return []Violation{{
				KataID: "ZC1260",
				Message: "Use `git branch -d` instead of `-D`. The lowercase `-d` refuses to " +
					"delete unmerged branches, preventing accidental data loss.",
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
		ID:       "ZC1261",
		Title:    "Avoid piping `base64 -d` output to shell execution",
		Severity: SeverityError,
		Description: "Decoding base64 and piping to `sh`/`zsh`/`eval` is a code injection risk. " +
			"Always inspect decoded content before execution.",
		Check: checkZC1261,
	})
}

func checkZC1261(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "base64" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-d" || val == "-D" {
			return []Violation{{
				KataID: "ZC1261",
				Message: "Inspect `base64 -d` output before piping to execution. " +
					"Blindly executing decoded content is a code injection vector.",
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
		ID:       "ZC1262",
		Title:    "Avoid `chmod -R 777` — recursive world-writable is critical",
		Severity: SeverityError,
		Description: "`chmod -R 777` makes every file and directory world-writable and executable. " +
			"Use specific permissions like `755` for directories and `644` for files.",
		Check: checkZC1262,
	})
}

func checkZC1262(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chmod" {
		return nil
	}

	hasRecursive := false
	has777 := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-R" {
			hasRecursive = true
		}
		if val == "777" {
			has777 = true
		}
	}

	if hasRecursive && has777 {
		return []Violation{{
			KataID: "ZC1262",
			Message: "Never use `chmod -R 777` — it makes everything world-writable. " +
				"Use `find -type d -exec chmod 755` and `find -type f -exec chmod 644` instead.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1263",
		Title:    "Use `apt-get` instead of `apt` in scripts",
		Severity: SeverityStyle,
		Description: "`apt` is designed for interactive use and its output format may change. " +
			"`apt-get` has a stable interface suitable for scripts and CI.",
		Check: checkZC1263,
		Fix:   fixZC1263,
	})
}

// fixZC1263 rewrites `apt` to `apt-get` at the command name position.
// Arguments stay intact — the two tools accept the same shape for
// install / upgrade / update / remove.
func fixZC1263(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "apt" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("apt"),
		Replace: "apt-get",
	}}
}

func checkZC1263(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "apt" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1263",
		Message: "Use `apt-get` instead of `apt` in scripts. " +
			"`apt` is for interactive use; `apt-get` has a stable scripting interface.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1264",
		Title:    "Use `dnf` instead of `yum` on modern Fedora/RHEL",
		Severity: SeverityStyle,
		Description: "`yum` is deprecated on Fedora 22+ and RHEL 8+. " +
			"`dnf` is the modern replacement with better dependency resolution.",
		Check: checkZC1264,
		Fix:   fixZC1264,
	})
}

// fixZC1264 rewrites `yum` to `dnf`. dnf is broadly compatible with
// yum's CLI surface so arguments carry over unchanged.
func fixZC1264(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "yum" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("yum"),
		Replace: "dnf",
	}}
}

func checkZC1264(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "yum" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1264",
		Message: "Use `dnf` instead of `yum`. `yum` is deprecated on modern " +
			"Fedora and RHEL; `dnf` has better dependency resolution.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1265",
		Title:    "Use `systemctl enable --now` to enable and start together",
		Severity: SeverityStyle,
		Description: "`systemctl enable` without `--now` only enables on next boot. " +
			"Use `--now` to enable and immediately start the service.",
		Check: checkZC1265,
		Fix:   fixZC1265,
	})
}

// fixZC1265 inserts ` --now` after the `enable` subcommand in a
// `systemctl enable …` invocation. Same subcommand-level insertion
// pattern as ZC1234's docker-run --rm.
func fixZC1265(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}
	var enableArg ast.Expression
	for _, arg := range cmd.Arguments {
		if arg.String() == "enable" {
			enableArg = arg
			break
		}
	}
	if enableArg == nil {
		return nil
	}
	tok := enableArg.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 || off+len("enable") > len(source) {
		return nil
	}
	if string(source[off:off+len("enable")]) != "enable" {
		return nil
	}
	insertAt := off + len("enable")
	insLine, insCol := offsetLineColZC1265(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " --now",
	}}
}

func offsetLineColZC1265(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1265(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}

	hasEnable := false
	hasNow := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "enable" {
			hasEnable = true
		}
		if val == "--now" {
			hasNow = true
		}
	}

	if hasEnable && !hasNow {
		return []Violation{{
			KataID: "ZC1265",
			Message: "Use `systemctl enable --now` to enable and start the service immediately. " +
				"Without `--now`, the service only starts on next boot.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1266",
		Title:    "Use `nproc` instead of parsing `/proc/cpuinfo`",
		Severity: SeverityStyle,
		Description: "Parsing `/proc/cpuinfo` for CPU count is fragile and platform-specific. " +
			"`nproc` is a portable, dedicated tool for this purpose.",
		Check: checkZC1266,
	})
}

func checkZC1266(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/proc/cpuinfo" {
			return []Violation{{
				KataID: "ZC1266",
				Message: "Use `nproc` instead of parsing `/proc/cpuinfo` for CPU count. " +
					"`nproc` is portable and available on Linux and macOS (via coreutils).",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1267",
		Title:    "Use `df -P` for POSIX-portable disk usage output",
		Severity: SeverityStyle,
		Description: "`df -h` output format varies across systems and locales. " +
			"Use `df -P` for single-line, fixed-format output safe for script parsing.",
		Check: checkZC1267,
		Fix:   fixZC1267,
	})
}

// fixZC1267 inserts ` -P` after the `df` command name. Detector
// narrows to `df -h` (script-unsafe), so only that shape is
// rewritten.
func fixZC1267(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "df" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("df") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1267(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -P",
	}}
}

func offsetLineColZC1267(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1267(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "df" {
		return nil
	}

	hasPortable := false
	hasHuman := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-P" {
			hasPortable = true
		}
		if val == "-h" {
			hasHuman = true
		}
	}

	if hasHuman && !hasPortable {
		return []Violation{{
			KataID: "ZC1267",
			Message: "Use `df -P` for script-safe output. `df -h` format varies across " +
				"systems and may split long device names across lines.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1268",
		Title:    "Use `du -sh --` to handle filenames starting with dash",
		Severity: SeverityInfo,
		Description: "`du -sh *` breaks if a filename starts with `-`. " +
			"Use `--` to signal end of options and safely handle all filenames.",
		Check: checkZC1268,
		Fix:   fixZC1268,
	})
}

// fixZC1268 inserts `-- ` before the first positional argument of a
// `du …` invocation that lacks the `--` end-of-options marker. The
// detector already gates on a glob (`*` / `.`) being present, and on
// the absence of `--`, so the insertion is idempotent.
func fixZC1268(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "du" {
		return nil
	}
	// Find the first positional (non-flag) argument.
	var positional ast.Expression
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if len(v) > 0 && v[0] != '-' {
			positional = arg
			break
		}
	}
	if positional == nil {
		return nil
	}
	tok := positional.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 {
		return nil
	}
	insLine, insCol := offsetLineColZC1268(source, off)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: "-- ",
	}}
}

func offsetLineColZC1268(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1268(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "du" {
		return nil
	}

	hasEndOfOpts := false
	hasGlob := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "--" {
			hasEndOfOpts = true
		}
		if val == "*" || val == "." {
			hasGlob = true
		}
	}

	if hasGlob && !hasEndOfOpts {
		return []Violation{{
			KataID: "ZC1268",
			Message: "Use `du -sh -- *` instead of `du -sh *`. The `--` prevents " +
				"filenames starting with `-` from being interpreted as options.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1269",
		Title:    "Use `pgrep` instead of `ps aux | grep` for process search",
		Severity: SeverityStyle,
		Description: "`ps aux | grep` matches itself in the process list requiring workarounds. " +
			"Use `pgrep` which is designed for process searching without self-matching.",
		Check: checkZC1269,
	})
}

func checkZC1269(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ps" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "aux" || val == "-ef" || val == "-e" {
			return []Violation{{
				KataID: "ZC1269",
				Message: "Use `pgrep` instead of `ps " + val + " | grep`. `pgrep` is purpose-built " +
					"for process searching and doesn't match itself.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1270",
		Title:    "Use `mktemp` instead of hardcoded `/tmp` paths",
		Severity: SeverityWarning,
		Description: "Hardcoding `/tmp/filename` is vulnerable to symlink attacks and race conditions. " +
			"Use `mktemp` to create unique temporary files safely.",
		Check: checkZC1270,
	})
}

func checkZC1270(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value
	if name != "touch" && name != "cat" && name != "echo" && name != "cp" && name != "mv" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if strings.HasPrefix(val, "/tmp/") && !strings.Contains(val, "$") {
			return []Violation{{
				KataID:  "ZC1270",
				Message: "Use `mktemp` instead of hardcoded `" + val + "`. Hardcoded `/tmp` paths are vulnerable to symlink attacks and race conditions.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityWarning,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1271",
		Title:    "Use `command -v` instead of `which` for command existence checks",
		Severity: SeverityStyle,
		Description: "`which` is not POSIX-standard and behaves inconsistently across systems. " +
			"Use `command -v` which is portable and built into Zsh.",
		Check: checkZC1271,
		Fix:   fixZC1271,
	})
}

// fixZC1271 rewrites `which` to `command -v` at the command name
// position. Single replacement — arguments stay untouched.
func fixZC1271(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "which" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("which"),
		Replace: "command -v",
	}}
}

func checkZC1271(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "which" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1271",
		Message: "Use `command -v` instead of `which`. `command -v` is POSIX-compliant and built into Zsh.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1272",
		Title:    "Use `install -m` instead of separate `cp` and `chmod`",
		Severity: SeverityStyle,
		Description: "`install` atomically copies a file and sets permissions in one step. " +
			"Using separate `cp` and `chmod` creates a window where the file has wrong permissions.",
		Check: checkZC1272,
	})
}

func checkZC1272(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cp" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/usr/local/bin" || val == "/usr/bin" || val == "/opt/bin" ||
			val == "/usr/local/sbin" || val == "/usr/sbin" {
			return []Violation{{
				KataID:  "ZC1272",
				Message: "Use `install -m 0755` instead of `cp` to system directories. `install` sets permissions atomically.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1273",
		Title:    "Use `grep -q` instead of redirecting grep output to `/dev/null`",
		Severity: SeverityStyle,
		Description: "`grep -q` suppresses output and exits on first match, which is faster and more " +
			"idiomatic than piping or redirecting to `/dev/null`.",
		Check: checkZC1273,
		Fix:   fixZC1273,
	})
}

// fixZC1273 inserts ` -q` after `grep` and strips the trailing
// `/dev/null` argument (including its leading whitespace). Two edits;
// the detector already gates on the absence of `-q`, so the rewrite
// is idempotent on a re-run.
func fixZC1273(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}
	var devNull ast.Expression
	for _, arg := range cmd.Arguments {
		if arg.String() == "/dev/null" {
			devNull = arg
			break
		}
	}
	if devNull == nil {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || IdentLenAt(source, nameOff) != len("grep") {
		return nil
	}
	insertAt := nameOff + len("grep")
	insLine, insCol := offsetLineColZC1273(source, insertAt)
	if insLine < 0 {
		return nil
	}
	stripEdits := zc1238StripFlag(source, devNull, "/dev/null")
	if stripEdits == nil {
		return nil
	}
	return append([]FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -q",
	}}, stripEdits...)
}

func offsetLineColZC1273(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1273(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	hasQuiet := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-q" || val == "--quiet" || val == "--silent" {
			hasQuiet = true
			break
		}
	}

	if hasQuiet {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/dev/null" {
			return []Violation{{
				KataID:  "ZC1273",
				Message: "Use `grep -q` instead of redirecting to `/dev/null`. It is faster and more idiomatic.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1274",
		Title:    "Use Zsh `${var:t}` instead of `basename`",
		Severity: SeverityStyle,
		Description: "Zsh provides the `:t` (tail) modifier for parameter expansion which extracts " +
			"the filename component, avoiding the overhead of forking `basename`.",
		Check: checkZC1274,
	})
}

func checkZC1274(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "basename" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1274",
		Message: "Use Zsh parameter expansion `${var:t}` instead of `basename`. The `:t` modifier extracts the filename without forking a process.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1275",
		Title:    "Use Zsh `${var:h}` instead of `dirname`",
		Severity: SeverityStyle,
		Description: "Zsh provides the `:h` (head) modifier for parameter expansion which extracts " +
			"the directory component, avoiding the overhead of forking `dirname`.",
		Check: checkZC1275,
	})
}

func checkZC1275(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "dirname" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1275",
		Message: "Use Zsh parameter expansion `${var:h}` instead of `dirname`. The `:h` modifier extracts the directory without forking a process.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1276",
		Title:    "Use Zsh `{start..end}` instead of `seq`",
		Severity: SeverityStyle,
		Description: "Zsh natively supports `{start..end}` brace expansion for generating number " +
			"sequences, avoiding the overhead of forking the external `seq` command.",
		Check: checkZC1276,
		Fix:   fixZC1061,
	})
}

func checkZC1276(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "seq" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1276",
		Message: "Use Zsh brace expansion `{start..end}` instead of `seq`. Brace expansion is built-in and avoids forking.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

// Issue #344: ZC1277 fires on the same input as the canonical
// ZC1108 with overlapping advice. Stubbed to a no-op so legacy
// `disabled_katas` lists keep parsing; the detection now lives in
// ZC1108.

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1277",
		Title:       "Superseded by ZC1108 — retired duplicate",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. See https://github.com/afadesigns/zshellcheck/issues/344 for context; the canonical detection lives in ZC1108.",
		Check:       checkZC1277,
	})
}

func checkZC1277(ast.Node) []Violation {
	return nil
}

// Issue #343: ZC1278 fires on the same input as the canonical
// ZC1009 with overlapping advice. Stubbed to a no-op so legacy
// `disabled_katas` lists keep parsing; the detection now lives in
// ZC1009.

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1278",
		Title:       "Superseded by ZC1009 — retired duplicate",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. See https://github.com/afadesigns/zshellcheck/issues/343 for context; the canonical detection lives in ZC1009.",
		Check:       checkZC1278,
	})
}

func checkZC1278(ast.Node) []Violation {
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1279",
		Title:    "Use `realpath` instead of `readlink -f` for canonical paths",
		Severity: SeverityInfo,
		Description: "`readlink -f` is not portable across all platforms (notably macOS). " +
			"Use `realpath` which is POSIX-standard and available on modern systems.",
		Check: checkZC1279,
		Fix:   fixZC1279,
	})
}

// fixZC1279 collapses `readlink -f` to `realpath` when `-f` is the
// first argument. Single span replacement from the start of
// `readlink` through the end of `-f`. Only fires when `-f` is the
// first argument; other shapes (`readlink -n -f`, `readlink path -f`)
// are left alone to avoid clobbering unrelated flags. Idempotent —
// a re-run sees `realpath`, not `readlink`. Defensive byte-match
// guards on both anchors.
func fixZC1279(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "readlink" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "-f" {
		return nil
	}
	cmdOff := LineColToByteOffset(source, v.Line, v.Column)
	if cmdOff < 0 || cmdOff+len("readlink") > len(source) {
		return nil
	}
	if string(source[cmdOff:cmdOff+len("readlink")]) != "readlink" {
		return nil
	}
	fTok := cmd.Arguments[0].TokenLiteralNode()
	fOff := LineColToByteOffset(source, fTok.Line, fTok.Column)
	if fOff < 0 || fOff+len("-f") > len(source) {
		return nil
	}
	if string(source[fOff:fOff+len("-f")]) != "-f" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  fOff + len("-f") - cmdOff,
		Replace: "realpath",
	}}
}

func checkZC1279(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "readlink" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-f" {
			return []Violation{{
				KataID:  "ZC1279",
				Message: "Use `realpath` instead of `readlink -f`. `realpath` is more portable, especially on macOS.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityInfo,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1280",
		Title:    "Use `Zsh ${var:e}` instead of shell expansion to extract file extension",
		Severity: SeverityStyle,
		Description: "Zsh provides the `:e` (extension) modifier for parameter expansion which " +
			"extracts the file extension, avoiding complex shell patterns or external tools.",
		Check: checkZC1280,
	})
}

func checkZC1280(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cut" {
		return nil
	}

	hasDot := false
	hasField := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-d." || val == "-d" {
			hasDot = true
		}
		if val == "-f2" {
			hasField = true
		}
		if val == "." {
			hasDot = true
		}
	}

	if hasDot && hasField {
		return []Violation{{
			KataID:  "ZC1280",
			Message: "Use Zsh parameter expansion `${var:e}` to extract the file extension instead of `cut -d. -f2`.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1281",
		Title:    "Use `sort -u` instead of `sort | uniq` for deduplication",
		Severity: SeverityStyle,
		Description: "`sort -u` combines sorting and deduplication in a single pass, " +
			"which is more efficient than piping `sort` into `uniq` as a separate process.",
		Check: checkZC1281,
	})
}

func checkZC1281(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "uniq" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-c" || val == "-d" || val == "-D" || val == "-u" {
			return nil
		}
	}

	return []Violation{{
		KataID:  "ZC1281",
		Message: "Use `sort -u` instead of `sort | uniq`. The `-u` flag deduplicates in a single pass.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1282",
		Title:    "Use Zsh `${var:r}` instead of `sed` to remove file extension",
		Severity: SeverityStyle,
		Description: "Zsh provides the `:r` modifier to remove a filename extension. " +
			"Using `sed` or `cut` to strip the extension is unnecessary when the built-in " +
			"parameter expansion handles it directly.",
		Check: checkZC1282,
	})
}

func checkZC1282(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sed" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "'s/\\.[^.]*$//'" || val == "s/\\.[^.]*$//" ||
			val == "'s/\\.[^.]*$//g'" || val == "s/\\.[^.]*$//g" {
			return []Violation{{
				KataID:  "ZC1282",
				Message: "Use Zsh parameter expansion `${var:r}` to remove the file extension instead of `sed`.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1283",
		Title:    "Use `setopt` instead of `set -o` for Zsh options",
		Severity: SeverityStyle,
		Description: "Zsh provides `setopt` and `unsetopt` as native builtins for managing shell " +
			"options. Using `set -o` / `set +o` is a POSIX compatibility form that is less " +
			"idiomatic in Zsh scripts.",
		Check: checkZC1283,
		Fix:   fixZC1283,
	})
}

// fixZC1283 rewrites `set -o OPTION` into `setopt OPTION`. The span
// covers the `set` command name and the `-o` flag in a single edit;
// trailing option arguments stay in place.
func fixZC1283(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "set" {
		return nil
	}
	var dashO ast.Expression
	for _, arg := range cmd.Arguments {
		if arg.String() == "-o" {
			dashO = arg
			break
		}
	}
	if dashO == nil {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("set") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("set")]) != "set" {
		return nil
	}
	dashTok := dashO.TokenLiteralNode()
	dashOff := LineColToByteOffset(source, dashTok.Line, dashTok.Column)
	if dashOff < 0 || dashOff+2 > len(source) {
		return nil
	}
	if string(source[dashOff:dashOff+2]) != "-o" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  dashOff + 2 - nameOff,
		Replace: "setopt",
	}}
}

func checkZC1283(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "set" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-o" {
			return []Violation{{
				KataID:  "ZC1283",
				Message: "Use `setopt` instead of `set -o` in Zsh scripts. `setopt` is the native Zsh idiom.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1284",
		Title:    "Use Zsh `${(s:sep:)var}` instead of `cut -d` for field splitting",
		Severity: SeverityStyle,
		Description: "Zsh provides the `(s:separator:)` parameter expansion flag to split strings " +
			"into arrays by a delimiter. This is more idiomatic than invoking `cut -d` and " +
			"avoids spawning an external process.",
		Check: checkZC1284,
	})
}

func checkZC1284(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cut" {
		return nil
	}

	hasDelim := false
	hasField := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if strings.HasPrefix(val, "-d") && val != "-d." {
			hasDelim = true
		}
		if strings.HasPrefix(val, "-f") {
			hasField = true
		}
	}

	if hasDelim && hasField {
		return []Violation{{
			KataID:  "ZC1284",
			Message: "Use Zsh parameter expansion `${(s:sep:)var}` for field splitting instead of `cut -d -f`.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1285",
		Title:    "Use Zsh `${(o)array}` for sorting instead of piping to `sort`",
		Severity: SeverityStyle,
		Description: "Zsh provides the `(o)` parameter expansion flag to sort array elements " +
			"in ascending order and `(O)` for descending order. This avoids spawning " +
			"an external `sort` process for simple array sorting.",
		Check: checkZC1285,
	})
}

func checkZC1285(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sort" {
		return nil
	}

	// sort with complex flags has legitimate uses beyond simple array sorting
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-t" || val == "-k" || val == "-n" || val == "-r" ||
			val == "-u" || val == "-h" || val == "-V" || val == "-g" ||
			val == "-c" || val == "-m" || val == "-s" {
			return nil
		}
	}

	// sort with only a filename — suggest Zsh native sorting
	if len(cmd.Arguments) == 1 {
		return []Violation{{
			KataID:  "ZC1285",
			Message: "Use Zsh `${(o)array}` for sorting instead of piping to `sort`. The `(o)` flag sorts in-shell.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1286",
		Title:    "Use Zsh `${array:#pattern}` instead of `grep -v` for filtering",
		Severity: SeverityStyle,
		Description: "Zsh provides `${array:#pattern}` to remove matching elements from an array " +
			"and `${(M)array:#pattern}` to keep only matching elements. This avoids " +
			"spawning an external `grep` process for simple filtering tasks.",
		Check: checkZC1286,
	})
}

func checkZC1286(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-v" {
			return []Violation{{
				KataID:  "ZC1286",
				Message: "Use Zsh `${array:#pattern}` for filtering instead of `grep -v`. Parameter expansion avoids a subprocess.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1287",
		Title:    "Use `cat -v` alternative: Zsh `${(V)var}` for visible control characters",
		Severity: SeverityStyle,
		Description: "Zsh provides the `(V)` parameter expansion flag to make control characters " +
			"visible in a variable. This avoids piping through `cat -v` for simple " +
			"visibility of non-printable characters.",
		Check: checkZC1287,
	})
}

func checkZC1287(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-v" || val == "-A" {
			return []Violation{{
				KataID:  "ZC1287",
				Message: "Use Zsh `${(V)var}` to make control characters visible instead of `cat -v`.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.DeclarationStatementNode, Kata{
		ID:       "ZC1288",
		Title:    "Use `typeset` instead of `declare` in Zsh scripts",
		Severity: SeverityStyle,
		Description: "`typeset` is the native Zsh builtin for variable declarations. " +
			"`declare` is a Bash compatibility alias. Using `typeset` is more idiomatic " +
			"and signals that the script is Zsh-native.",
		Check: checkZC1288,
		Fix:   fixZC1288,
	})
}

// fixZC1288 rewrites the `declare` keyword to `typeset`. Arguments,
// flags and assignments carry over unchanged because the two
// builtins share the same Zsh interface.
func fixZC1288(node ast.Node, v Violation, source []byte) []FixEdit {
	decl, ok := node.(*ast.DeclarationStatement)
	if !ok || decl.Command != "declare" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("declare"),
		Replace: "typeset",
	}}
}

func checkZC1288(node ast.Node) []Violation {
	decl, ok := node.(*ast.DeclarationStatement)
	if !ok {
		return nil
	}

	if decl.Command != "declare" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1288",
		Message: "Use `typeset` instead of `declare` in Zsh scripts. `typeset` is the native Zsh idiom.",
		Line:    decl.Token.Line,
		Column:  decl.Token.Column,
		Level:   SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1289",
		Title:    "Use Zsh `${(u)array}` for unique elements instead of `sort -u`",
		Severity: SeverityStyle,
		Description: "Zsh provides the `(u)` parameter expansion flag to remove duplicate " +
			"elements from an array. This preserves original order and avoids spawning " +
			"an external `sort -u` process.",
		Check: checkZC1289,
	})
}

func checkZC1289(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sort" {
		return nil
	}

	hasUnique := false
	hasOtherFlags := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-u" {
			hasUnique = true
		} else if len(val) > 1 && val[0] == '-' {
			hasOtherFlags = true
		}
	}

	if hasUnique && !hasOtherFlags {
		return []Violation{{
			KataID:  "ZC1289",
			Message: "Use Zsh `${(u)array}` for unique elements instead of `sort -u`. The `(u)` flag preserves order.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1290",
		Title:    "Use Zsh `${(n)array}` for numeric sorting instead of `sort -n`",
		Severity: SeverityStyle,
		Description: "Zsh provides the `(n)` parameter expansion flag to sort array elements " +
			"numerically. This avoids spawning an external `sort -n` process for " +
			"simple numeric sorting of array data.",
		Check: checkZC1290,
	})
}

func checkZC1290(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sort" {
		return nil
	}

	hasNumeric := false
	hasOtherFlags := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-n" {
			hasNumeric = true
		} else if len(val) > 1 && val[0] == '-' {
			hasOtherFlags = true
		}
	}

	if hasNumeric && !hasOtherFlags {
		return []Violation{{
			KataID:  "ZC1290",
			Message: "Use Zsh `${(n)array}` for numeric sorting instead of `sort -n`. The `(n)` flag sorts numerically in-shell.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1291",
		Title:    "Use Zsh `${(O)array}` for reverse sorting instead of `sort -r`",
		Severity: SeverityStyle,
		Description: "Zsh provides the `(O)` parameter expansion flag to sort array elements " +
			"in descending (reverse) order. This avoids spawning an external `sort -r` " +
			"process for simple reverse sorting of array data.",
		Check: checkZC1291,
	})
}

func checkZC1291(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sort" {
		return nil
	}

	hasReverse := false
	hasOtherFlags := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-r" {
			hasReverse = true
		} else if len(val) > 1 && val[0] == '-' {
			hasOtherFlags = true
		}
	}

	if hasReverse && !hasOtherFlags {
		return []Violation{{
			KataID:  "ZC1291",
			Message: "Use Zsh `${(O)array}` for reverse sorting instead of `sort -r`. The `(O)` flag sorts descending in-shell.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1292",
		Title:    "Use Zsh `${var//old/new}` instead of `tr` for character translation",
		Severity: SeverityStyle,
		Description: "Zsh provides `${var//old/new}` for global substitution within a variable. " +
			"For simple single-character translation, this avoids spawning `tr` as an " +
			"external process.",
		Check: checkZC1292,
	})
}

func checkZC1292(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tr" {
		return nil
	}

	if len(cmd.Arguments) != 2 {
		return nil
	}

	first := strings.Trim(cmd.Arguments[0].String(), "'\"")
	second := strings.Trim(cmd.Arguments[1].String(), "'\"")

	// Only flag simple single-character translations, not ranges or classes
	if len(first) == 1 && len(second) == 1 {
		return []Violation{{
			KataID:  "ZC1292",
			Message: "Use Zsh `${var//" + first + "/" + second + "}` for character substitution instead of `tr`.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1293",
		Title:    "Use `[[ ]]` instead of `test` command in Zsh",
		Severity: SeverityStyle,
		Description: "Zsh `[[ ]]` provides a more powerful conditional expression syntax than " +
			"the `test` command. It supports pattern matching, regex, and does not require " +
			"quoting of variable expansions to prevent word splitting.",
		Check: checkZC1293,
		Fix:   fixZC1293,
	})
}

// fixZC1293 rewrites `test EXPR…` into `[[ EXPR… ]]`. Two edits:
// the `test` command name becomes `[[`; ` ]]` is appended after the
// last argument's source span. Bails when there are no arguments
// (a bare `test` is invalid anyway). Idempotent — a re-run sees
// `[[ ... ]]`, not `test`.
func fixZC1293(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "test" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	cmdOff := LineColToByteOffset(source, v.Line, v.Column)
	if cmdOff < 0 || cmdOff+len("test") > len(source) {
		return nil
	}
	if string(source[cmdOff:cmdOff+len("test")]) != "test" {
		return nil
	}
	lastArg := cmd.Arguments[len(cmd.Arguments)-1]
	lastTok := lastArg.TokenLiteralNode()
	lastOff := LineColToByteOffset(source, lastTok.Line, lastTok.Column)
	if lastOff < 0 {
		return nil
	}
	lastVal := lastArg.String()
	lastEnd := lastOff + len(lastVal)
	if lastEnd > len(source) {
		return nil
	}
	endLine, endCol := offsetLineColZC1293(source, lastEnd)
	if endLine < 0 {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: len("test"), Replace: "[["},
		{Line: endLine, Column: endCol, Length: 0, Replace: " ]]"},
	}
}

func offsetLineColZC1293(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1293(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "test" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1293",
		Message: "Use `[[ ]]` instead of the `test` command in Zsh. `[[ ]]` is more powerful and does not require variable quoting.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1294",
		Title:    "Use `bindkey` instead of `bind` for key bindings in Zsh",
		Severity: SeverityWarning,
		Description: "`bind` is a Bash builtin for key bindings. Zsh uses `bindkey` for " +
			"ZLE (Zsh Line Editor) key bindings. Using `bind` in a Zsh script will " +
			"fail unless Bash compatibility is loaded.",
		Check: checkZC1294,
	})
}

func checkZC1294(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "bind" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1294",
		Message: "Use `bindkey` instead of `bind` in Zsh. `bind` is a Bash builtin; Zsh uses `bindkey` for ZLE key bindings.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1295",
		Title:    "Use `vared` instead of `read -e` for interactive editing in Zsh",
		Severity: SeverityStyle,
		Description: "Zsh provides `vared` for interactive editing of variables with full " +
			"ZLE support (tab completion, history, cursor movement). The `read -e` flag " +
			"is a Bash extension; Zsh `vared` is the native equivalent.",
		Check: checkZC1295,
	})
}

func checkZC1295(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "read" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-e" {
			return []Violation{{
				KataID:  "ZC1295",
				Message: "Use `vared` instead of `read -e` in Zsh. `vared` provides full ZLE editing support natively.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1296",
		Title:    "Avoid `shopt` in Zsh — use `setopt`/`unsetopt` instead",
		Severity: SeverityWarning,
		Description: "`shopt` is a Bash builtin that does not exist in Zsh. Use `setopt` " +
			"or `unsetopt` to control Zsh shell options. Common Bash `shopt` options " +
			"have Zsh equivalents via `setopt`.",
		Check: checkZC1296,
	})
}

func checkZC1296(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "shopt" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1296",
		Message: "Avoid `shopt` in Zsh — it is a Bash builtin. Use `setopt`/`unsetopt` for Zsh shell options.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1297",
		Title:    "Avoid `$BASH_SOURCE` — use `$0` or `${(%):-%x}` in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_SOURCE` is a Bash-specific variable that does not exist in Zsh. " +
			"In Zsh, use `$0` inside a sourced file to get the script path, or " +
			"`${(%):-%x}` for the current file regardless of sourcing context.",
		Check: checkZC1297,
		Fix:   fixZC1297,
	})
}

// fixZC1297 renames the Bash `$BASH_SOURCE` identifier to the Zsh
// `${(%):-%x}` prompt-flag expansion that resolves to the current file
// regardless of sourcing context. Only the dollar-prefixed form is
// rewritten — the bare `BASH_SOURCE` form (inside `${...}` or as an
// assignment target) is left to manual review because the surrounding
// braces would need adjusting too.
func fixZC1297(node ast.Node, v Violation, _ []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok || ident == nil {
		return nil
	}
	if ident.Value != "$BASH_SOURCE" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("$BASH_SOURCE"),
		Replace: "${(%):-%x}",
	}}
}

func checkZC1297(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_SOURCE" && ident.Value != "BASH_SOURCE" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1297",
		Message: "Avoid `$BASH_SOURCE` in Zsh — use `$0` or `${(%):-%x}` instead. `BASH_SOURCE` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1298",
		Title:    "Avoid `$FUNCNAME` — use `$funcstack` in Zsh",
		Severity: SeverityWarning,
		Description: "`$FUNCNAME` is a Bash-specific array that does not exist in Zsh. " +
			"Zsh provides `$funcstack` as the equivalent, containing the call stack " +
			"of function names with the current function at index 1.",
		Check: checkZC1298,
		Fix:   fixZC1298,
	})
}

// fixZC1298 renames the Bash `$FUNCNAME` identifier to the Zsh
// `$funcstack` equivalent. Handles both the dollar-prefixed and
// bare forms.
func fixZC1298(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}
	switch ident.Value {
	case "$FUNCNAME":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$FUNCNAME"),
			Replace: "$funcstack",
		}}
	case "FUNCNAME":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("FUNCNAME"),
			Replace: "funcstack",
		}}
	}
	return nil
}

func checkZC1298(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$FUNCNAME" && ident.Value != "FUNCNAME" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1298",
		Message: "Avoid `$FUNCNAME` in Zsh — use `$funcstack` instead. `FUNCNAME` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1299",
		Title:    "Avoid `$BASH_LINENO` — use `$funcfiletrace` in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_LINENO` is a Bash-specific array that does not exist in Zsh. " +
			"Zsh provides `$funcfiletrace` as the equivalent, containing file:line " +
			"pairs for each call in the function stack.",
		Check: checkZC1299,
	})
}

func checkZC1299(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_LINENO" && ident.Value != "BASH_LINENO" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1299",
		Message: "Avoid `$BASH_LINENO` in Zsh — use `$funcfiletrace` instead. `BASH_LINENO` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
