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
		ID:       "ZC1500",
		Title:    "Warn on `systemctl edit <unit>` in scripts — requires interactive editor",
		Severity: SeverityWarning,
		Description: "`systemctl edit <unit>` (without `--no-edit` and without a piped `EDITOR`) " +
			"opens `$EDITOR` on a tmpfile and waits for the user. In a non-interactive script " +
			"this either hangs until timeout or silently succeeds with no change, depending on " +
			"how the editor handles a closed stdin. For scripted unit tweaks, drop a `.conf` " +
			"drop-in under `/etc/systemd/system/<unit>.d/` and call `systemctl daemon-reload`.",
		Check: checkZC1500,
	})
}

func checkZC1500(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "systemctl" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	if cmd.Arguments[0].String() != "edit" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--no-edit" || v == "--runtime" {
			// Still odd, but at least doesn't spin on an editor. Let it pass.
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1500",
		Message: "`systemctl edit` opens $EDITOR and waits for the user. Use a drop-in " +
			"`/etc/systemd/system/<unit>.d/*.conf` + `daemon-reload` in scripts.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1501",
		Title:    "Style: `docker-compose` (hyphen) — use `docker compose` (space, built-in plugin)",
		Severity: SeverityStyle,
		Description: "`docker-compose` is the Python Compose V1 binary. Docker stopped shipping " +
			"it with Docker Desktop in 2023 and Compose V2 is now the first-class `docker " +
			"compose` subcommand. Scripts that invoke `docker-compose` silently degrade on " +
			"fresh installs and miss V2-only options (`--profile`, `--wait`, richer env " +
			"interpolation). Call `docker compose` (space) or pin the V2 binary explicitly.",
		Check: checkZC1501,
		Fix:   fixZC1501,
	})
}

// fixZC1501 rewrites the hyphenated `docker-compose` command name into
// the space-separated `docker compose` subcommand form. Arguments stay
// untouched — the V2 plugin accepts the same shape.
func fixZC1501(node ast.Node, v Violation, _ []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker-compose" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("docker-compose"),
		Replace: "docker compose",
	}}
}

func checkZC1501(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker-compose" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1501",
		Message: "`docker-compose` is the deprecated Python V1 binary. Use `docker compose` " +
			"(space-separated subcommand) for the bundled V2 plugin.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1502",
		Title:    "Warn on `grep \"$var\" file` without `--` — flag injection when `$var` starts with `-`",
		Severity: SeverityWarning,
		Description: "Without a `--` end-of-flags marker, `grep` (and most POSIX tools) treats " +
			"any argument that starts with `-` as a flag. If `$var` comes from user input or a " +
			"fuzzed filename, an attacker can pass `--include=*secret*` or `-f /etc/shadow` " +
			"and get grep to read paths the script author never intended. Always write " +
			"`grep -- \"$var\" file` or use a grep-compatible library with explicit pattern API.",
		Check: checkZC1502,
		Fix:   fixZC1502,
	})
}

// fixZC1502 inserts `-- ` before the first variable-shaped argument
// of a grep / egrep / fgrep / rg / ag invocation that lacks the
// `--` end-of-options marker. Idempotent — the detector gates on
// the absence of `--`, so once `-- ` is present a re-run won't
// re-insert.
func fixZC1502(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || !zc1502IsGrepFamily(cmd) {
		return nil
	}
	firstVar := zc1502FirstVarArg(cmd)
	if firstVar == nil {
		return nil
	}
	tok := firstVar.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 {
		return nil
	}
	insLine, insCol := offsetLineColZC1502(source, off)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{Line: insLine, Column: insCol, Length: 0, Replace: "-- "}}
}

var zc1502GrepFamily = map[string]struct{}{
	"grep": {}, "egrep": {}, "fgrep": {}, "rg": {}, "ag": {},
}

func zc1502IsGrepFamily(cmd *ast.SimpleCommand) bool {
	_, hit := zc1502GrepFamily[CommandIdentifier(cmd)]
	return hit
}

func zc1502FirstVarArg(cmd *ast.SimpleCommand) ast.Expression {
	var first ast.Expression
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--" {
			return nil
		}
		if first == nil && (strings.HasPrefix(v, "\"$") || strings.HasPrefix(v, "$")) {
			first = arg
		}
	}
	return first
}

func offsetLineColZC1502(source []byte, offset int) (int, int) {
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

func checkZC1502(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || !zc1502IsGrepFamily(cmd) {
		return nil
	}
	firstVar := zc1502FirstVarArg(cmd)
	if firstVar == nil {
		return nil
	}
	return []Violation{{
		KataID: "ZC1502",
		Message: "Variable `" + firstVar.String() + "` used as pattern without `--` end-of-flags " +
			"marker — attacker-controlled leading `-` becomes a flag. Write `grep -- \"$var\"`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1503",
		Title:    "Error on `groupadd -g 0` / `groupmod -g 0` — creates duplicate root group",
		Severity: SeverityError,
		Description: "Creating or renaming a group to GID 0 gives its members the same privileges " +
			"as members of `root` for every file that grants permissions to GID 0. Combined " +
			"with `usermod -G 0 <user>` it becomes an invisible privilege escalation path. " +
			"Distro tooling already reserves GID 0 for `root`; pick a sensible unused GID " +
			"(`getent group` gives the list) and scope access via sudoers or polkit.",
		Check: checkZC1503,
	})
}

func checkZC1503(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "groupadd" && ident.Value != "groupmod" && ident.Value != "addgroup" {
		return nil
	}

	var prevG bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevG {
			prevG = false
			if v == "0" {
				return zc1503Violation(cmd)
			}
		}
		if v == "-g" || v == "--gid" {
			prevG = true
		}
		if v == "-g0" || v == "--gid=0" {
			return zc1503Violation(cmd)
		}
	}
	return nil
}

func zc1503Violation(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1503",
		Message: "Creating a group with GID 0 duplicates the `root` group — hidden privesc. " +
			"Pick an unused GID (see `getent group`) and scope via sudoers/polkit.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1504",
		Title:    "Warn on `git push --mirror` — overwrites every remote ref",
		Severity: SeverityWarning,
		Description: "`git push --mirror` pushes every ref under `refs/` and deletes any remote " +
			"ref that is not present locally. Running it against a shared origin instantly " +
			"wipes everyone else's branches and tags. Legitimate uses are mirror-to-mirror " +
			"replication where the source is the authoritative tree; for everyday pushes use " +
			"an explicit refspec or `git push --all`.",
		Check: checkZC1504,
	})
}

func checkZC1504(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "push" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--mirror" {
			return []Violation{{
				KataID: "ZC1504",
				Message: "`git push --mirror` overwrites every remote ref and deletes ones " +
					"missing locally. Use an explicit refspec or `--all` for everyday pushes.",
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
		ID:       "ZC1505",
		Title:    "Warn on `dpkg --force-confnew` / `--force-confold` — silently overrides /etc changes",
		Severity: SeverityWarning,
		Description: "`--force-confnew` replaces any locally-modified config file with the " +
			"maintainer version; `--force-confold` keeps the local file and drops the new " +
			"defaults on the floor. Either way dpkg silently picks a side without prompting, " +
			"so a legitimate /etc tweak (hardening, compliance override) can vanish or a " +
			"security-relevant config update can be ignored. Review the conffile diff per " +
			"upgrade (`ucf` / `etckeeper`) rather than hard-coding the decision.",
		Check: checkZC1505,
	})
}

func checkZC1505(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dpkg" && ident.Value != "apt" && ident.Value != "apt-get" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "--force-conf") ||
			strings.HasPrefix(v, "-oDpkg::Options::=--force-conf") ||
			strings.HasPrefix(v, "-o=Dpkg::Options::=--force-conf") {
			return []Violation{{
				KataID: "ZC1505",
				Message: "`" + v + "` silently picks maintainer or local conffile — legit /etc " +
					"changes disappear or new defaults are ignored. Use ucf/etckeeper.",
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
		ID:       "ZC1506",
		Title:    "Warn on `newgrp <group>` in scripts — spawns a new shell, breaks control flow",
		Severity: SeverityWarning,
		Description: "`newgrp` starts a new login shell with the requested primary group. Inside " +
			"a non-interactive script that shell inherits no commands, so the script either " +
			"hangs waiting for the user or exits immediately depending on stdin. If the script " +
			"genuinely needs temporarily-augmented group access, call `sg <group> -c <cmd>` " +
			"or, in a service context, use `SupplementaryGroups=` in the unit file.",
		Check: checkZC1506,
	})
}

func checkZC1506(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "newgrp" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1506",
		Message: "`newgrp` starts a new shell — script either hangs or exits. Use " +
			"`sg <group> -c <cmd>` or systemd `SupplementaryGroups=`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1507",
		Title:    "Warn on `rsync -l` / default symlink handling — follows escaping symlinks",
		Severity: SeverityWarning,
		Description: "By default rsync copies symlinks as-is but does not prevent one from " +
			"pointing outside the source tree. When the destination is rooted elsewhere (or " +
			"the receiver creates a file at the symlink's resolved path) this becomes a path " +
			"traversal primitive. Use `--safe-links` to skip symlinks pointing outside the " +
			"transfer set, or `--copy-unsafe-links` to materialise them as regular files.",
		Check: checkZC1507,
	})
}

func checkZC1507(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "rsync" {
		return nil
	}

	// Trigger only when rsync is actually asked to handle symlinks:
	// -l (preserve symlinks) or -a (archive, which includes -l).
	var hasSymlinkMode, hasSafeHandling bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-a" || v == "-l" || v == "--archive" || v == "--links" ||
			strings.HasPrefix(v, "-a") && strings.ContainsAny(v[1:], "lavx") {
			hasSymlinkMode = true
		}
		if v == "--safe-links" || v == "--copy-unsafe-links" || v == "--no-links" ||
			v == "--munge-links" {
			hasSafeHandling = true
		}
	}
	if !hasSymlinkMode || hasSafeHandling {
		return nil
	}

	return []Violation{{
		KataID: "ZC1507",
		Message: "`rsync` preserving symlinks without `--safe-links` follows ones pointing " +
			"outside the source tree — path traversal vector. Add `--safe-links` or " +
			"`--copy-unsafe-links`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1508",
		Title:    "Style: `ldd <binary>` may execute the binary — use `objdump -p` / `readelf -d` for untrusted files",
		Severity: SeverityStyle,
		Description: "On glibc, `ldd` is implemented by setting `LD_TRACE_LOADED_OBJECTS=1` and " +
			"invoking the binary. A malicious ELF with a custom interpreter (`PT_INTERP`) or " +
			"constructors can therefore run code when `ldd` is pointed at it. `objdump -p " +
			"<file> | grep NEEDED` or `readelf -d <file>` give the same shared-library list " +
			"without executing the binary.",
		Check: checkZC1508,
	})
}

func checkZC1508(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ldd" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1508",
		Message: "`ldd` on glibc can execute the target binary. Use `objdump -p` or " +
			"`readelf -d` to inspect ELF dependencies safely.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1509",
		Title:    "Warn on `trap '' TERM` / `trap - TERM` — ignores/resets fatal signal",
		Severity: SeverityWarning,
		Description: "`trap '' <signal>` makes the signal uninterruptible. `trap - <signal>` " +
			"restores the default disposition, which on `TERM`/`INT`/`HUP` means the script " +
			"exits without running any cleanup handler. Both forms are routinely used to " +
			"harden long-running scripts against accidental `Ctrl-C`, but also to hide from " +
			"`kill` during incident response. Keep the explicit cleanup handler on at least " +
			"`EXIT` so state is always unwound.",
		Check: checkZC1509,
	})
}

func checkZC1509(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "trap" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}
	handler := cmd.Arguments[0].String()
	if handler != "''" && handler != `""` && handler != "-" {
		return nil
	}
	// Signals after the handler.
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		switch v {
		case "TERM", "SIGTERM", "INT", "SIGINT", "HUP", "SIGHUP",
			"QUIT", "SIGQUIT":
			return []Violation{{
				KataID: "ZC1509",
				Message: "`trap " + handler + " " + v + "` silences a fatal signal — cleanup " +
					"handlers never run. Keep at least a cleanup trap on EXIT.",
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
		ID:       "ZC1510",
		Title:    "Error on `auditctl -e 0` / `auditctl -D` — disables kernel audit logging",
		Severity: SeverityError,
		Description: "`auditctl -e 0` switches the Linux audit subsystem off, and `auditctl -D` " +
			"deletes every audit rule, including the ones that monitor `/etc/shadow`, `execve`, " +
			"and privilege escalations. Both are textbook anti-forensics steps. If you need to " +
			"temporarily quiet audit for a maintenance window, use `-e 2` (lock enabled + " +
			"immutable) to require a reboot for any further change and document the action.",
		Check: checkZC1510,
	})
}

func checkZC1510(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "auditctl" {
		return nil
	}

	var prevE bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-D" {
			return zc1510Violation(cmd, "-D", "deletes every audit rule")
		}
		if prevE {
			prevE = false
			if v == "0" {
				return zc1510Violation(cmd, "-e 0", "disables audit subsystem")
			}
		}
		if v == "-e" {
			prevE = true
		}
	}
	return nil
}

func zc1510Violation(cmd *ast.SimpleCommand, flag, what string) []Violation {
	return []Violation{{
		KataID: "ZC1510",
		Message: "`auditctl " + flag + "` " + what + " — anti-forensics tactic. Use `-e 2` " +
			"for a reboot-locked maintenance window instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1511",
		Title:    "Error on `nmcli ... <wireless/vpn secret>` on command line",
		Severity: SeverityError,
		Description: "Passing Wi-Fi pre-shared keys or VPN secrets as positional `nmcli` args " +
			"puts them in `ps`, shell history, and `/proc/<pid>/cmdline`. Let NetworkManager " +
			"store the secret for you via `--ask` (interactive prompt, no TTY echo) or use " +
			"`keyfile` connection profiles under `/etc/NetworkManager/system-connections/` " +
			"with mode 0600.",
		Check: checkZC1511,
	})
}

var nmcliSecretKeys = []string{
	"802-11-wireless-security.psk",
	"wifi-sec.psk",
	"wifi.psk",
	"vpn.secrets.password",
	"ipsec-secret",
	"openvpn-password",
	"802-1x.password",
}

func checkZC1511(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "nmcli" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	for i, a := range args {
		low := strings.ToLower(a)
		for _, key := range nmcliSecretKeys {
			if low == key && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				return []Violation{{
					KataID: "ZC1511",
					Message: "`nmcli` passed `" + key + " <secret>` on the command line — " +
						"ends up in ps/history. Use `--ask` or a keyfile profile.",
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
		ID:       "ZC1512",
		Title:    "Style: `service <unit> <verb>` — use `systemctl <verb> <unit>` on systemd hosts",
		Severity: SeverityStyle,
		Description: "`service` is the SysV init compatibility wrapper. On a systemd-managed " +
			"host (every mainstream distro since ~2016) it translates to `systemctl` anyway, " +
			"but reverses argument order, loses `--user` scope, ignores unit templating, and " +
			"can't restart sockets or timers. Prefer `systemctl start|stop|restart|reload " +
			"<unit>` for consistency across scripts and interactive shells.",
		Check: checkZC1512,
		Fix:   fixZC1512,
	})
}

// fixZC1512 rewrites `service UNIT VERB` into `systemctl VERB UNIT`.
// Three edits per match: rename `service` → `systemctl`, swap the
// textual contents of the UNIT and VERB positions. Gated to simple
// Identifier args so the swap stays byte-exact; concat-form units
// (rare in practice) stay detection-only.
func fixZC1512(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "service" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	unitIdent, ok := cmd.Arguments[0].(*ast.Identifier)
	if !ok {
		return nil
	}
	verbIdent, ok := cmd.Arguments[1].(*ast.Identifier)
	if !ok {
		return nil
	}
	unitTok := unitIdent.TokenLiteralNode()
	verbTok := verbIdent.TokenLiteralNode()
	unitOff := LineColToByteOffset(source, unitTok.Line, unitTok.Column)
	verbOff := LineColToByteOffset(source, verbTok.Line, verbTok.Column)
	if unitOff < 0 || verbOff < 0 {
		return nil
	}
	if unitOff+len(unitIdent.Value) > len(source) ||
		string(source[unitOff:unitOff+len(unitIdent.Value)]) != unitIdent.Value {
		return nil
	}
	if verbOff+len(verbIdent.Value) > len(source) ||
		string(source[verbOff:verbOff+len(verbIdent.Value)]) != verbIdent.Value {
		return nil
	}
	return []FixEdit{
		{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("service"),
			Replace: "systemctl",
		},
		{
			Line:    unitTok.Line,
			Column:  unitTok.Column,
			Length:  len(unitIdent.Value),
			Replace: verbIdent.Value,
		},
		{
			Line:    verbTok.Line,
			Column:  verbTok.Column,
			Length:  len(verbIdent.Value),
			Replace: unitIdent.Value,
		},
	}
}

func checkZC1512(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "service" {
		return nil
	}

	// Needs at least <unit> <verb>.
	if len(cmd.Arguments) < 2 {
		return nil
	}
	verb := cmd.Arguments[1].String()
	switch verb {
	case "start", "stop", "restart", "reload", "status", "force-reload", "try-restart":
	default:
		return nil
	}

	unit := cmd.Arguments[0].String()
	return []Violation{{
		KataID: "ZC1512",
		Message: "`service " + unit + " " + verb + "` — prefer `systemctl " + verb + " " +
			unit + "` for consistency with other systemd commands.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1513",
		Title:    "Style: `make install` without `DESTDIR=` — unmanaged system-wide install",
		Severity: SeverityStyle,
		Description: "`make install` drops files directly into `$(prefix)` with no package " +
			"manager tracking. Upgrades can leave stale files behind, uninstalls rely on " +
			"`make uninstall` being accurate, and the operation typically needs `sudo`. For " +
			"local builds, set `DESTDIR=/tmp/pkgroot` + wrap in `checkinstall` / `fpm` / " +
			"distro packaging, or use `stow` / `xstow` to manage symlinks under `/usr/local`.",
		Check: checkZC1513,
	})
}

func checkZC1513(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "make" && ident.Value != "gmake" {
		return nil
	}

	hasInstall := false
	hasDestdir := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "install" {
			hasInstall = true
		}
		if len(v) >= 8 && v[:8] == "DESTDIR=" {
			hasDestdir = true
		}
	}
	if !hasInstall || hasDestdir {
		return nil
	}
	return []Violation{{
		KataID: "ZC1513",
		Message: "`make install` without `DESTDIR=` leaves no package-manager record. Set " +
			"`DESTDIR=/tmp/pkgroot` and wrap in checkinstall / fpm, or use stow.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1514",
		Title:    "Error on `useradd -p <hash>` / `usermod -p <hash>` — password hash on cmdline",
		Severity: SeverityError,
		Description: "`-p` takes an already-hashed password (crypt(3) format) and writes it " +
			"to `/etc/shadow`. That hash is in `ps`, `/proc/<pid>/cmdline`, and history for as " +
			"long as the process runs — enough time for a co-tenant to grab it and start an " +
			"offline crack. Use `chpasswd` with `--crypt-method=SHA512` reading from stdin, " +
			"or write `/etc/shadow` via a configuration-management tool with proper file " +
			"permissions.",
		Check: checkZC1514,
	})
}

func checkZC1514(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "useradd" && ident.Value != "usermod" && ident.Value != "adduser" {
		return nil
	}

	var prevP bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevP && v != "" && v[0] != '-' {
			return []Violation{{
				KataID: "ZC1514",
				Message: "`" + ident.Value + " -p <hash>` puts the hashed password in ps / " +
					"/proc / history. Use `chpasswd --crypt-method=SHA512` from stdin.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
		prevP = (v == "-p" || v == "--password")
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1515",
		Title:    "Warn on `md5sum` / `sha1sum` for integrity check — collision-vulnerable",
		Severity: SeverityWarning,
		Description: "MD5 and SHA-1 are broken for collision resistance: public attacks cheaply " +
			"craft two different files with the same hash. For verifying a download against a " +
			"published checksum, or for comparing archives against a manifest, use " +
			"`sha256sum` / `sha512sum` / `b2sum` instead. MD5 is still fine for non-adversarial " +
			"cache keys but almost every invocation in scripts is the integrity case.",
		Check: checkZC1515,
	})
}

func checkZC1515(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "md5sum" && ident.Value != "sha1sum" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1515",
		Message: "`" + ident.Value + "` is collision-vulnerable — don't use for integrity " +
			"checks. Use `sha256sum` / `sha512sum` / `b2sum` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1516",
		Title:    "Error on `umask 000` / `umask 0` — new files / directories world-writable",
		Severity: SeverityError,
		Description: "`umask 000` means every file created after this line inherits mode 0666 " +
			"and every directory inherits 0777 — world-readable, world-writable, no " +
			"authorization layer. On a multi-user host (build runner, shared workstation) this " +
			"leaks secrets through the filesystem and invites tampering. Pick a sensible " +
			"umask (`022` for public software, `077` for secrets handling).",
		Check: checkZC1516,
	})
}

func checkZC1516(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "umask" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	v := cmd.Arguments[0].String()
	if v == "0" || v == "00" || v == "000" || v == "0000" {
		return []Violation{{
			KataID: "ZC1516",
			Message: "`umask " + v + "` leaves new files world-readable and world-writable. " +
				"Use `022` for public software, `077` for secrets handling.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1517",
		Title:    "Warn on `print -P \"$var\"` — prompt-escape injection via user-controlled string",
		Severity: SeverityWarning,
		Description: "`print -P` enables prompt-escape expansion (`%F`, `%K`, `%B`, `%S`, plus " +
			"arbitrary command substitution via `%{...%}`). Interpolating a shell variable " +
			"means any of those sequences inside the variable are expanded — at best messing " +
			"up terminal state, at worst running the attacker's command via `%(e:...)` or " +
			"similar. Either drop `-P` or wrap the variable with `${(q-)var}` / `${(V)var}` " +
			"to neutralize metacharacters before printing.",
		Check: checkZC1517,
	})
}

func checkZC1517(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "print" {
		return nil
	}

	hasP := false
	var varArg string
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-P" {
			hasP = true
			continue
		}
		if !hasP || varArg != "" {
			continue
		}
		// Look for interpolation:  "$x" or $x (double-quoted or unquoted).
		raw := v
		if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
			raw = raw[1 : len(raw)-1]
		}
		if strings.Contains(raw, "$") && !(len(v) >= 2 && v[0] == '\'' && v[len(v)-1] == '\'') {
			varArg = v
		}
	}
	if !hasP || varArg == "" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1517",
		Message: "`print -P " + varArg + "` expands prompt escapes inside the variable — use " +
			"`${(V)var}` / `${(q-)var}` to neutralize metacharacters, or drop -P.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1518",
		Title:    "Warn on `bash -p` — privileged mode (skips env sanitisation on setuid)",
		Severity: SeverityWarning,
		Description: "`bash -p` (and `-o privileged`) tells bash not to drop its effective UID/GID " +
			"and not to sanitize the environment when started on a setuid wrapper. It is " +
			"explicitly the flag you use to keep `BASH_ENV`, `SHELLOPTS`, and similar " +
			"attacker-controlled variables active while running as a more privileged user. " +
			"Almost no legitimate script needs `-p`; audit and remove.",
		Check: checkZC1518,
	})
}

func checkZC1518(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "bash" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-p" {
			return []Violation{{
				KataID: "ZC1518",
				Message: "`bash -p` keeps the privileged environment on a setuid wrapper — " +
					"almost never needed, audit and remove.",
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
		ID:       "ZC1519",
		Title:    "Warn on `ulimit -u unlimited` — removes user process cap, enables fork bombs",
		Severity: SeverityWarning,
		Description: "`ulimit -u` caps the number of processes a UID can run; `unlimited` removes " +
			"that cap. Combined with a bug in a background loop (or a literal fork bomb via " +
			"`:(){ :|:& };:`) it pegs the scheduler until the machine has to be cold-booted. " +
			"Pick a realistic number (distro defaults around 4096 for interactive sessions) or " +
			"set it in `/etc/security/limits.d/` so it is persistent and visible.",
		Check: checkZC1519,
	})
}

func checkZC1519(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ulimit" {
		return nil
	}

	var prevU bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevU {
			prevU = false
			if v == "unlimited" {
				return []Violation{{
					KataID: "ZC1519",
					Message: "`ulimit -u unlimited` removes the user process cap — fork bomb " +
						"surface. Pick a realistic number or set it via /etc/security/limits.d/.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
		if v == "-u" {
			prevU = true
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1520",
		Title:    "Warn on `vared <var>` in scripts — reads interactively, hangs non-interactive",
		Severity: SeverityWarning,
		Description: "`vared` is the Zsh interactive line-editor builtin that lets the user edit " +
			"the value of a variable in place. In a non-interactive script (cron job, CI " +
			"runner, ssh-with-command) `vared` has no TTY, so the script either errors out or " +
			"hangs waiting for input that never arrives. For scripted input, read the value " +
			"from stdin (`read varname`), a file, or an environment variable.",
		Check: checkZC1520,
	})
}

func checkZC1520(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "vared" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1520",
		Message: "`vared` requires a TTY — in a non-interactive script it errors or hangs. " +
			"Use `read`, stdin, or environment variables for scripted input.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1521",
		Title:    "Style: `strace` without `-e` filter — captures every syscall (incl. secrets, huge output)",
		Severity: SeverityStyle,
		Description: "Unfiltered `strace` records every syscall the process makes: every " +
			"`read()`/`write()` buffer, every `connect()` sockaddr, every `open()` path. That " +
			"includes passwords read from stdin, session tokens written to TLS sockets, and " +
			"any memory a `write()` buffer happens to point at. Scope with `-e trace=<set>` " +
			"(e.g. `trace=openat,connect`) and strip sensitive content with `-e abbrev=all`.",
		Check: checkZC1521,
	})
}

func checkZC1521(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "strace" && ident.Value != "ltrace" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	// Any filter flag present → skip.
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-e" || v == "--trace" || v == "--trace-path" || v == "-P" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1521",
		Message: "`" + ident.Value + "` without `-e` captures every syscall including secrets " +
			"in read/write buffers. Scope with `-e trace=<set>`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1522",
		Title:    "Warn on `ip route add default` / `route add default` — changes default gateway",
		Severity: SeverityWarning,
		Description: "Setting a new default route in a script silently redirects every non-local " +
			"packet through the specified gateway. That is exactly the knob an attacker turns " +
			"to MITM a whole host after a foothold, and it is also a common accidental foot- " +
			"gun in CI runners (gateway in the runner network ≠ gateway in production). Use " +
			"NetworkManager / systemd-networkd config files for persistent routes, and " +
			"document any runtime change with a comment explaining why.",
		Check: checkZC1522,
	})
}

func checkZC1522(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// `ip route add default ...`
	if ident.Value == "ip" && len(args) >= 3 &&
		args[0] == "route" && args[1] == "add" && args[2] == "default" {
		return zc1522Violation(cmd, "ip route add default")
	}
	// `route add default ...`
	if ident.Value == "route" && len(args) >= 2 &&
		args[0] == "add" && args[1] == "default" {
		return zc1522Violation(cmd, "route add default")
	}
	return nil
}

func zc1522Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1522",
		Message: "`" + what + "` silently reroutes every non-local packet through the new " +
			"gateway — MITM surface or CI footgun. Use NetworkManager/networkd config.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1523",
		Title:    "Error on `tar -C /` — extracting an archive into the filesystem root",
		Severity: SeverityError,
		Description: "Extracting a tarball directly into `/` overwrites any file it carries a " +
			"matching path for. Combined with a malicious tarball that contains entries like " +
			"`etc/pam.d/sshd` or `usr/bin/ls`, this is a full system compromise disguised as a " +
			"software install. Always extract into a staging directory, inspect contents, then " +
			"copy specific files into place.",
		Check: checkZC1523,
	})
}

func checkZC1523(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tar" && ident.Value != "bsdtar" {
		return nil
	}

	var prevC bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevC {
			prevC = false
			if v == "/" {
				return []Violation{{
					KataID: "ZC1523",
					Message: "`tar -C /` extracts into the filesystem root — overwrites any " +
						"path that happens to be inside the archive. Stage, inspect, then copy.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
		if v == "-C" || v == "--directory" {
			prevC = true
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1524",
		Title:    "Warn on `sysctl -e` / `sysctl -q` — silently skip unknown keys, hide config drift",
		Severity: SeverityWarning,
		Description: "`sysctl -e` and `-q` suppress error output for unknown keys or failed " +
			"writes. That is how a typo in `/etc/sysctl.d/99-hardening.conf` goes unnoticed " +
			"for months — the hardening didn't actually take effect because the key name was " +
			"wrong. Drop `-e`/`-q` in scripts and let errors bubble up; fix the offending " +
			"conffile instead.",
		Check: checkZC1524,
	})
}

func checkZC1524(node ast.Node) []Violation {
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
		if v == "-e" || v == "-q" || v == "-eq" || v == "-qe" {
			return []Violation{{
				KataID: "ZC1524",
				Message: "`sysctl " + v + "` suppresses error output — typos in sysctl.d/ " +
					"conffiles silently skip. Remove and surface the real error.",
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
		ID:       "ZC1525",
		Title:    "Warn on `ping -f` — flood ping sends packets as fast as possible",
		Severity: SeverityWarning,
		Description: "`ping -f` (flood mode) removes the one-per-second rate limit and sends " +
			"ICMP echo requests in a tight loop. It's a root-only builtin specifically because " +
			"it can saturate a slow link or overload a low-end host. Legitimate uses exist " +
			"(latency benchmarking, stress testing known-internal targets), but in a script " +
			"aimed at arbitrary hosts it is a noisy traffic generator. Scope tightly and " +
			"document.",
		Check: checkZC1525,
	})
}

func checkZC1525(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ping" && ident.Value != "ping6" && ident.Value != "fping" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-f" {
			return []Violation{{
				KataID: "ZC1525",
				Message: "`" + ident.Value + " -f` (flood) bypasses the rate limit — saturates " +
					"slow links. Scope tightly and document.",
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
		ID:       "ZC1526",
		Title:    "Error on `wipefs -a` / `wipefs -af` — erases filesystem signatures (unrecoverable)",
		Severity: SeverityError,
		Description: "`wipefs -a` overwrites every filesystem, partition table, and RAID signature " +
			"it finds on the target. Unlike `rm`, there is no retention anywhere — the only " +
			"recovery path is a disk image backup taken beforehand. If the target variable is " +
			"wrong (typo, empty, resolves to the wrong `/dev/sdX`), this bricks the disk. " +
			"Always run with `--no-act` first or prefer `sgdisk --zap-all` for partition-table " +
			"scope.",
		Check: checkZC1526,
	})
}

func checkZC1526(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "wipefs" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-a" || v == "-af" || v == "-fa" || v == "--all" {
			return []Violation{{
				KataID: "ZC1526",
				Message: "`wipefs -a` erases every filesystem signature — unrecoverable. Run " +
					"with `--no-act` first, or use `sgdisk --zap-all` for scoped deletion.",
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
		ID:       "ZC1527",
		Title:    "Warn on `crontab -` — replaces cron from stdin, overwrites without diff",
		Severity: SeverityWarning,
		Description: "`crontab -` (or `crontab -u <user> -`) reads a full crontab from stdin and " +
			"replaces the user's existing entries wholesale. Any manual tweak, oncall " +
			"override, or colleague's row is silently deleted. Paired with `curl | crontab -` " +
			"it is a common persistence one-liner. Use `crontab -l > /tmp/old && ... " +
			"crontab -e` with an explicit diff/merge, or ship cron entries via " +
			"`/etc/cron.d/*` managed by config tooling.",
		Check: checkZC1527,
	})
}

func checkZC1527(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "crontab" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-" {
			return []Violation{{
				KataID: "ZC1527",
				Message: "`crontab -` overwrites the user's crontab from stdin — silently " +
					"drops existing rows. Use /etc/cron.d/ files or a diff/merge workflow.",
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
		ID:       "ZC1528",
		Title:    "Warn on `chage -M 99999` / `-E -1` — disables password aging / expiry",
		Severity: SeverityWarning,
		Description: "`chage -M 99999` sets the max password age to roughly 273 years (effectively " +
			"never). `chage -E -1` clears the account expiration date. Both silently remove an " +
			"automatic lockout mechanism a compromised credential would otherwise hit. If " +
			"passwords genuinely should not expire (SSO, cert-based auth), encode that in a " +
			"PAM profile rather than per-user `chage`.",
		Check: checkZC1528,
	})
}

func checkZC1528(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chage" {
		return nil
	}

	var prevM, prevE bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevM {
			prevM = false
			if v == "99999" || v == "-1" || v == "0" {
				return zc1528Violation(cmd, "-M "+v)
			}
		}
		if prevE {
			prevE = false
			if v == "-1" {
				return zc1528Violation(cmd, "-E -1")
			}
		}
		switch v {
		case "-M", "--maxdays":
			prevM = true
		case "-E", "--expiredate":
			prevE = true
		}
	}
	return nil
}

func zc1528Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1528",
		Message: "`chage " + what + "` disables password aging — removes automatic lockout. " +
			"Use a PAM profile instead of per-user chage.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1529",
		Title:    "Warn on `fsck -y` / `fsck.<fs> -y` — auto-answer yes can corrupt",
		Severity: SeverityWarning,
		Description: "`fsck -y` answers `yes` to every repair prompt. For the happy case it is a " +
			"timesaver, but on a filesystem with unusual corruption (bad sector storm, mangled " +
			"journal after power loss) the automatic answer can turn salvageable data into " +
			"`lost+found` entries or zero it outright. In scripts, prefer `fsck -n` for a " +
			"dry-run and let a human adjudicate a real repair, or run with `-p` (preen: only " +
			"safe automatic fixes).",
		Check: checkZC1529,
	})
}

func checkZC1529(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "fsck" && !strings.HasPrefix(ident.Value, "fsck.") &&
		ident.Value != "e2fsck" && ident.Value != "xfs_repair" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-y" {
			return []Violation{{
				KataID: "ZC1529",
				Message: "`" + ident.Value + " -y` answers yes to every repair prompt — can " +
					"destroy salvageable data. Prefer `-n` (dry-run) or `-p` (preen).",
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
		ID:       "ZC1530",
		Title:    "Warn on `pkill -f <pattern>` — matches full command line, easy to over-kill",
		Severity: SeverityWarning,
		Description: "`pkill -f` matches the pattern against the full command line, not just " +
			"the process name. A pattern like `-f server` also matches the `grep -- server` " +
			"in a user's shell history or any backup tool named `server-backup`. For routine " +
			"use, drop `-f` (matches process name only) or scope with `-U <uid>` / `-G " +
			"<gid>` / `-P <ppid>`. When you must match the command line, pin it with `^` / `$` " +
			"anchors in the pattern.",
		Check: checkZC1530,
	})
}

func checkZC1530(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "pkill" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-f" {
			return []Violation{{
				KataID: "ZC1530",
				Message: "`pkill -f` matches the full command line — easy to over-kill. Drop " +
					"`-f`, scope with `-U/-G/-P`, or anchor the pattern with ^/$.",
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
		ID:       "ZC1531",
		Title:    "Warn on `wget -t 0` — infinite retries, hangs on a dead endpoint",
		Severity: SeverityWarning,
		Description: "`wget -t 0` (or `--tries=0`) means retry forever. Paired with `-w` (wait " +
			"between retries) and a dead endpoint, the script hangs until killed — in a cron " +
			"job, every subsequent invocation piles up and eventually the UID's process limit " +
			"trips. Use a finite retry count (`-t 5`) plus `--timeout=<seconds>` to cap total " +
			"wall time.",
		Check: checkZC1531,
	})
}

func checkZC1531(node ast.Node) []Violation {
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

	var prevT bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevT {
			prevT = false
			if v == "0" {
				return []Violation{{
					KataID: "ZC1531",
					Message: "`wget -t 0` retries forever — script hangs on dead endpoint. " +
						"Use finite `-t 5` plus `--timeout=<seconds>`.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
		if v == "-t" {
			prevT = true
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1532",
		Title:    "Warn on `screen -dm` / `tmux new-session -d` — detached long-running session",
		Severity: SeverityWarning,
		Description: "Starting a detached screen/tmux session from a script puts a long-running " +
			"process outside the systemd supervisory tree: no logs in the journal, no cgroup " +
			"accounting, no restart-on-failure, no OOM scoring. It is also a common post- " +
			"compromise persistence technique because the session survives the initial shell " +
			"exit and hides in `ps -ef` as a short tmux/screen helper. For real long-running " +
			"work, write a systemd unit (user or system) and start it with `systemctl " +
			"[--user] start`.",
		Check: checkZC1532,
	})
}

func checkZC1532(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "screen" {
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-dm" || v == "-dmS" {
				return zc1532Violation(cmd, "screen "+v)
			}
		}
	}
	if ident.Value == "tmux" && len(cmd.Arguments) >= 2 && cmd.Arguments[0].String() == "new-session" {
		for _, arg := range cmd.Arguments[1:] {
			if arg.String() == "-d" {
				return zc1532Violation(cmd, "tmux new-session -d")
			}
		}
	}
	return nil
}

func zc1532Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1532",
		Message: "`" + what + "` backgrounds work outside systemd — no journal, no cgroup, " +
			"common persistence technique. Use a systemd unit instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1533",
		Title:    "Warn on `setsid <cmd>` — detaches from controlling TTY, escapes supervision",
		Severity: SeverityWarning,
		Description: "`setsid` starts a new session and process group. Combined with `-f` " +
			"(`--fork`) the child is fully detached from the invoking shell: `SIGHUP` from " +
			"logout does not reach it, the tty hang-up no longer terminates it, and it falls " +
			"off the script's job table. That is legitimate for daemonising a long-running " +
			"helper (though systemd does this better) and is also a standard persistence " +
			"mechanism. Prefer a systemd unit; if you must detach, document why.",
		Check: checkZC1533,
	})
}

func checkZC1533(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setsid" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1533",
		Message: "`setsid` detaches the child from the TTY / session — escapes supervision. " +
			"Prefer a systemd unit; document a detach if one is genuinely needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1534",
		Title:    "Warn on `dmesg -c` / `--clear` — wipes kernel ring buffer",
		Severity: SeverityWarning,
		Description: "`dmesg -c` reads and then clears the kernel ring buffer. Any subsequent " +
			"reader sees an empty log, so OOM kills, driver panics, and audit messages that " +
			"landed between the wipe and the incident response are gone. It is also an " +
			"anti-forensics step in post-exploitation playbooks. Use `dmesg` (no flags) for a " +
			"read, and let the journal retention policy handle rotation.",
		Check: checkZC1534,
	})
}

func checkZC1534(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dmesg" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c" || v == "-C" || v == "--clear" || v == "--read-clear" {
			return []Violation{{
				KataID: "ZC1534",
				Message: "`dmesg " + v + "` wipes the kernel ring buffer — subsequent " +
					"readers see no OOM/panic/audit messages. Read without clearing.",
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
		ID:       "ZC1535",
		Title:    "Warn on `ip link set <iface> promisc on` — enables packet capture",
		Severity: SeverityWarning,
		Description: "Putting an interface into promiscuous mode tells the NIC to deliver every " +
			"frame to userspace, not just frames addressed to this host. Legitimate for tools " +
			"like tcpdump/tshark (which turn it on themselves) but running it from a script " +
			"and leaving it on is a sniffer-in-place — traffic from other hosts on the same " +
			"broadcast domain lands in anyone's `tshark -i`. Re-disable as soon as capture is " +
			"done, and prefer giving tcpdump `CAP_NET_RAW` so the mode is scoped to a single " +
			"invocation.",
		Check: checkZC1535,
	})
}

func checkZC1535(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ip" && ident.Value != "ifconfig" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// ip link set <iface> promisc on
	if ident.Value == "ip" {
		for i := 0; i+4 < len(args); i++ {
			if args[i] == "link" && args[i+1] == "set" && args[i+3] == "promisc" && args[i+4] == "on" {
				return zc1535Violation(cmd)
			}
		}
	}
	// ifconfig <iface> promisc
	if ident.Value == "ifconfig" {
		for _, a := range args {
			if a == "promisc" {
				return zc1535Violation(cmd)
			}
		}
	}
	return nil
}

func zc1535Violation(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1535",
		Message: "Interface put into promiscuous mode — sniffer-in-place. Re-disable after " +
			"capture, or grant tcpdump CAP_NET_RAW instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1536",
		Title:    "Warn on `iptables -j DNAT` / `-j REDIRECT` — rewrites traffic destination",
		Severity: SeverityWarning,
		Description: "`-j DNAT` and `-j REDIRECT` in an iptables rule rewrite the destination " +
			"address/port of matching packets. That is how you transparently proxy, but also " +
			"how you silently redirect a victim's connections to an attacker-controlled port. " +
			"Scripts that touch NAT rules should be carefully reviewed; prefer declarative " +
			"network config (nftables ruleset, NetworkManager connection, firewalld service) " +
			"and store rule provenance.",
		Check: checkZC1536,
	})
}

var (
	zc1536AddVerbs   = map[string]struct{}{
		"-A": {}, "-I": {}, "-R": {},
		"--append": {}, "--insert": {}, "--replace": {},
	}
	zc1536DnatTargets = map[string]struct{}{"DNAT": {}, "REDIRECT": {}, "NETMAP": {}}
	zc1536Tools       = map[string]struct{}{"iptables": {}, "ip6tables": {}}
)

func checkZC1536(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if _, hit := zc1536Tools[CommandIdentifier(cmd)]; !hit {
		return nil
	}
	args := zc1464StringArgs(cmd)
	if !zc1536HasAddVerb(args) {
		return nil
	}
	tgt := zc1536FirstDnatTarget(args)
	if tgt == "" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1536",
		Message: "`iptables -j " + tgt + "` rewrites packet destination — " +
			"silent redirect surface. Use declarative nftables/firewalld config.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func zc1536HasAddVerb(args []string) bool {
	for _, a := range args {
		if _, hit := zc1536AddVerbs[a]; hit {
			return true
		}
	}
	return false
}

func zc1536FirstDnatTarget(args []string) string {
	for i, a := range args {
		if a != "-j" || i+1 >= len(args) {
			continue
		}
		if _, hit := zc1536DnatTargets[args[i+1]]; hit {
			return args[i+1]
		}
	}
	return ""
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1537",
		Title:    "Error on `lvremove -f` / `vgremove -f` / `pvremove -f` — force-destroys LVM metadata",
		Severity: SeverityError,
		Description: "The `-f`/`--force` flag on the LVM destructive commands skips the " +
			"confirmation prompt that protects against a typo in the volume name. If the " +
			"target variable resolves to the wrong VG/LV/PV (empty, unset, different host), " +
			"a single line destroys every filesystem on top of that LVM stack. Leave the " +
			"prompt in and pipe `yes` to it only when you have explicitly confirmed the " +
			"target immediately beforehand.",
		Check: checkZC1537,
	})
}

func checkZC1537(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "lvremove" && ident.Value != "vgremove" && ident.Value != "pvremove" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-f" || v == "-ff" || v == "--force" {
			return []Violation{{
				KataID: "ZC1537",
				Message: "`" + ident.Value + " " + v + "` skips the confirmation — a typo in " +
					"the volume name destroys every filesystem on top of it.",
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
		ID:       "ZC1538",
		Title:    "Error on `zpool destroy -f` / `zfs destroy -rR` — recursive ZFS destruction",
		Severity: SeverityError,
		Description: "`zpool destroy -f` nukes a whole ZFS pool including every dataset, " +
			"snapshot, and clone on it. `zfs destroy -r` recurses into descendant datasets; " +
			"`-R` additionally drops descendant clones. Unlike `rm`, the space is freed " +
			"immediately and there is no recycle bin. Always require `zfs list`/`zpool list` " +
			"+ explicit target confirmation in the same script block, and prefer snapshot-" +
			"based rollback for recoverable workflows.",
		Check: checkZC1538,
	})
}

var zc1538ZfsRecursiveFlags = map[string]struct{}{
	"-r": {}, "-R": {}, "-rR": {}, "-Rr": {}, "-rf": {}, "-fr": {},
}

func checkZC1538(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	switch CommandIdentifier(cmd) {
	case "zpool":
		if zc1538DestroyFlag(cmd, map[string]struct{}{"-f": {}}) != "" {
			return zc1538Violation(cmd, "zpool destroy -f")
		}
	case "zfs":
		if hit := zc1538DestroyFlag(cmd, zc1538ZfsRecursiveFlags); hit != "" {
			return zc1538Violation(cmd, "zfs destroy "+hit)
		}
	}
	return nil
}

func zc1538DestroyFlag(cmd *ast.SimpleCommand, flags map[string]struct{}) string {
	if len(cmd.Arguments) < 2 || cmd.Arguments[0].String() != "destroy" {
		return ""
	}
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if _, hit := flags[v]; hit {
			return v
		}
	}
	return ""
}

func zc1538Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1538",
		Message: "`" + what + "` irrecoverably destroys the ZFS pool/dataset and every " +
			"snapshot on it. Require explicit target confirmation.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1539",
		Title:    "Warn on `parted -s <disk> <destructive-op>` — script mode bypasses confirmation",
		Severity: SeverityWarning,
		Description: "`parted -s` (script mode) answers the `data will be destroyed` prompt " +
			"with `yes`. Combined with `mklabel`, `mkpart`, `rm`, or `resizepart` on the " +
			"wrong device variable it silently repartitions or zeros the partition table on a " +
			"disk the author never intended. Require an explicit `parted <disk> print` check " +
			"plus an out-of-band confirmation before the destructive call.",
		Check: checkZC1539,
	})
}

func checkZC1539(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "parted" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	var hasScript bool
	for _, a := range args {
		if a == "-s" {
			hasScript = true
		}
	}
	if !hasScript {
		return nil
	}
	for _, a := range args {
		switch a {
		case "mklabel", "mkpart", "rm", "resizepart", "mkpartfs":
			return []Violation{{
				KataID: "ZC1539",
				Message: "`parted -s <disk> " + a + "` bypasses the confirmation prompt — a " +
					"typo in the disk variable silently repartitions the wrong device.",
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
		ID:       "ZC1540",
		Title:    "Error on `cryptsetup erase` / `luksErase` — destroys LUKS header, data unrecoverable",
		Severity: SeverityError,
		Description: "`cryptsetup erase` (alias `luksErase`) overwrites the LUKS header and " +
			"every key slot. Without the header the ciphertext on the device is unrecoverable " +
			"— even the original passphrase cannot unlock it. Keep a `cryptsetup " +
			"luksHeaderBackup` image somewhere safe before running erase, and prefer " +
			"`luksRemoveKey`/`luksKillSlot` when only rotating one slot.",
		Check: checkZC1540,
	})
}

func checkZC1540(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "cryptsetup" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "erase" || v == "luksErase" {
			return []Violation{{
				KataID: "ZC1540",
				Message: "`cryptsetup " + v + "` wipes the LUKS header — ciphertext becomes " +
					"unrecoverable. Back up the header first, or use luksRemoveKey/" +
					"luksKillSlot for single-slot rotation.",
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
		ID:       "ZC1541",
		Title:    "Error on `apk add --allow-untrusted` — installs unsigned Alpine packages",
		Severity: SeverityError,
		Description: "`apk add --allow-untrusted` skips signature verification on the package " +
			"being installed. On Alpine that is a direct MITM-to-root path: any mirror, " +
			"cache, or typo-squat can slip a replacement `.apk` and the daemon starts running " +
			"attacker code on next restart. Sign internal packages with your own key in " +
			"`/etc/apk/keys/` and keep verification on.",
		Check: checkZC1541,
	})
}

func checkZC1541(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "apk" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "--allow-untrusted" {
			return []Violation{{
				KataID: "ZC1541",
				Message: "`apk --allow-untrusted` skips signature verification on the " +
					"package — MITM-to-root on Alpine. Sign and place key in /etc/apk/keys/.",
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
		ID:       "ZC1542",
		Title:    "Error on `snap install --dangerous` — installs unsigned snap",
		Severity: SeverityError,
		Description: "`snap install --dangerous` tells snapd to install a snap that is not " +
			"assertion-verified. That bypass is named after the risk: any `.snap` file on disk " +
			"can register system services, confinement profiles, and hooks, running as whatever " +
			"user the snap declares. Use `--devmode` for developer work (still verified) or " +
			"ship the snap through the store / a private brand store for production rollouts.",
		Check: checkZC1542,
	})
}

func checkZC1542(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "snap" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "--dangerous" {
			return []Violation{{
				KataID: "ZC1542",
				Message: "`snap install --dangerous` installs an assertion-unverified snap — " +
					"any .snap on disk can register system services. Use --devmode or the " +
					"store.",
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
		ID:       "ZC1543",
		Title:    "Warn on `go install pkg@latest` / `cargo install --git <url>` without rev pin",
		Severity: SeverityWarning,
		Description: "`go install pkg@latest` and `cargo install --git <url>` without `--rev` / " +
			"`--tag` / `--branch` resolve to whatever HEAD is at install time. The next CI " +
			"run can pull a different commit — great for supply-chain attackers to inject " +
			"post-breach, bad for reproducibility. Pin to a specific version tag (`pkg@v1.2.3`) " +
			"or a commit hash (`cargo install --rev abc123 --git ...`).",
		Check: checkZC1543,
	})
}

var (
	zc1543MovingTags = []string{"@latest", "@master", "@main"}
	zc1543CargoPins  = map[string]struct{}{"--rev": {}, "--tag": {}, "--branch": {}, "--locked": {}}
)

func checkZC1543(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	args := zc1464StringArgs(cmd)
	switch CommandIdentifier(cmd) {
	case "go":
		if where := zc1543GoInstall(args); where != "" {
			return zc1543Violation(cmd, where)
		}
	case "cargo":
		if zc1543CargoInstallUnpinned(args) {
			return zc1543Violation(cmd, "cargo install --git (no --rev/--tag/--branch)")
		}
	}
	return nil
}

func zc1543GoInstall(args []string) string {
	if len(args) < 2 || args[0] != "install" {
		return ""
	}
	for _, a := range args[1:] {
		if strings.HasPrefix(a, "-") {
			continue
		}
		for _, suffix := range zc1543MovingTags {
			if strings.HasSuffix(a, suffix) {
				return "go install " + a
			}
		}
		if !strings.Contains(a, "@") && strings.Contains(a, "/") {
			return "go install " + a + " (no @version)"
		}
	}
	return ""
}

func zc1543CargoInstallUnpinned(args []string) bool {
	if len(args) < 2 || args[0] != "install" {
		return false
	}
	hasGit, hasPin := false, false
	for _, a := range args[1:] {
		if a == "--git" {
			hasGit = true
		}
		if _, hit := zc1543CargoPins[a]; hit {
			hasPin = true
		}
	}
	return hasGit && !hasPin
}

func zc1543Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1543",
		Message: "`" + what + "` is unpinned — HEAD-of-default can change between runs. Pin " +
			"to a version tag or commit hash for reproducibility.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1544",
		Title:    "Warn on `dnf copr enable` / `add-apt-repository ppa:` — unvetted third-party repo",
		Severity: SeverityWarning,
		Description: "Enabling a COPR project or an Ubuntu PPA pulls packages signed by a single " +
			"community contributor — there is no distro security team or reproducible-builds " +
			"guarantee behind that key. Any future compromise of that contributor's account " +
			"ships a rootkit to every box that ran this line. If you need the package badly " +
			"enough, pin to a specific `build-id`, verify the key fingerprint out of band, " +
			"and mirror to an internal repository.",
		Check: checkZC1544,
	})
}

func checkZC1544(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "dnf" && len(cmd.Arguments) >= 2 &&
		cmd.Arguments[0].String() == "copr" && cmd.Arguments[1].String() == "enable" {
		return zc1544Violation(cmd, "dnf copr enable")
	}
	if ident.Value == "add-apt-repository" {
		return zc1544Violation(cmd, "add-apt-repository")
	}
	return nil
}

func zc1544Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1544",
		Message: "`" + what + "` pulls from a single-contributor repo — no distro security " +
			"team. Pin the build, verify key fingerprint, mirror internally.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1545",
		Title:    "Warn on `docker system prune -af --volumes` — drops unused volumes too",
		Severity: SeverityWarning,
		Description: "`docker system prune -af --volumes` removes stopped containers, unused " +
			"networks, dangling images — and every volume not currently attached to a running " +
			"container. On a host where `docker-compose down` is used casually (shutdown " +
			"before a laptop close, for example), the matching database volume looks " +
			"\"unused\" to prune and goes with it. Drop `--volumes` from routine cleanup, or " +
			"target specific prune scopes (`docker image prune`, `docker container prune`).",
		Check: checkZC1545,
	})
}

var (
	zc1545Runtimes      = map[string]struct{}{"docker": {}, "podman": {}, "nerdctl": {}}
	zc1545PruneSubcmds  = map[string]struct{}{"system": {}, "volume": {}}
	zc1545AllVolFlags   = map[string]struct{}{
		"--volumes": {}, "-a": {}, "--all": {},
		"-af": {}, "-fa": {}, "--all --volumes": {},
	}
)

func checkZC1545(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	tool := CommandIdentifier(cmd)
	if _, hit := zc1545Runtimes[tool]; !hit {
		return nil
	}
	args := zc1464StringArgs(cmd)
	pruneTarget, ok := zc1545PruneTarget(args)
	if !ok || !zc1545HasAllOrVolumes(args[2:]) {
		return nil
	}
	return []Violation{{
		KataID: "ZC1545",
		Message: "`" + tool + " " + pruneTarget + " prune` with `-a`/`--volumes` drops " +
			"unused volumes — stopped stacks lose their databases. Scope the prune.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func zc1545PruneTarget(args []string) (string, bool) {
	if len(args) < 2 || args[1] != "prune" {
		return "", false
	}
	if _, hit := zc1545PruneSubcmds[args[0]]; !hit {
		return "", false
	}
	return args[0], true
}

func zc1545HasAllOrVolumes(args []string) bool {
	for _, a := range args {
		if _, hit := zc1545AllVolFlags[a]; hit {
			return true
		}
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1546",
		Title:    "Warn on `kubectl delete --force --grace-period=0` — skips PreStop, corrupts state",
		Severity: SeverityWarning,
		Description: "`kubectl delete --force --grace-period=0` tells the API server to remove " +
			"the resource from etcd without waiting for the kubelet to run PreStop hooks or " +
			"drain the pod. For a StatefulSet pod this routinely corrupts the backing PV " +
			"(database mid-flush, file lock left held) and the replacement pod refuses to " +
			"start. Use standard delete and let the graceful shutdown run; only reach for " +
			"`--force` when the node itself is gone.",
		Check: checkZC1546,
	})
}

func checkZC1546(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" && ident.Value != "oc" {
		return nil
	}

	var sawDelete, hasForce, hasGrace0 bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "delete" {
			sawDelete = true
		}
		if !sawDelete {
			continue
		}
		if v == "--force" {
			hasForce = true
		}
		if v == "--grace-period=0" {
			hasGrace0 = true
		}
	}
	if sawDelete && hasForce && hasGrace0 {
		return []Violation{{
			KataID: "ZC1546",
			Message: "`kubectl delete --force --grace-period=0` skips PreStop hooks and kubelet " +
				"drain — corrupts StatefulSet state. Use standard delete.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1547",
		Title:    "Warn on `kubectl apply --prune --all` — deletes resources missing from manifest",
		Severity: SeverityWarning,
		Description: "`kubectl apply --prune --all` (or `--prune -l <selector>`) deletes every " +
			"cluster resource whose label matches but which is not in the manifest you just " +
			"applied. In a partial-repo deploy or a manifest typo, that can delete production " +
			"Deployments, Services, or Secrets another team owns. Pair `--prune` with a " +
			"narrow `-l` selector unique to your stack, or use a GitOps controller (Argo CD, " +
			"Flux) that scopes prune to its own Application.",
		Check: checkZC1547,
	})
}

func checkZC1547(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" && ident.Value != "oc" {
		return nil
	}

	var sawApply, hasPrune, hasAll bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "apply" {
			sawApply = true
		}
		if !sawApply {
			continue
		}
		if v == "--prune" {
			hasPrune = true
		}
		if v == "--all" || v == "-A" || v == "--all-namespaces" {
			hasAll = true
		}
	}
	if sawApply && hasPrune && hasAll {
		return []Violation{{
			KataID: "ZC1547",
			Message: "`kubectl apply --prune --all` deletes every matching resource not in the " +
				"manifest — manifest typo wipes other teams' resources. Scope with a " +
				"narrow `-l <selector>`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1548",
		Title:    "Warn on `helm install/upgrade --disable-openapi-validation` — skips schema check",
		Severity: SeverityWarning,
		Description: "`--disable-openapi-validation` tells Helm to skip the OpenAPI schema check " +
			"the API server would apply. Malformed CRD instances or Deployments with " +
			"invalid spec fields then silently land in etcd, only failing when the " +
			"controller tries to reconcile — usually 3am, usually in prod. Keep the " +
			"validation on; fix the schema deviation instead.",
		Check: checkZC1548,
	})
}

func checkZC1548(node ast.Node) []Violation {
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

	var sawVerb bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "install" || v == "upgrade" || v == "template" {
			sawVerb = true
			continue
		}
		if !sawVerb {
			continue
		}
		if v == "--disable-openapi-validation" {
			return []Violation{{
				KataID: "ZC1548",
				Message: "`helm --disable-openapi-validation` hides bad manifests until the " +
					"controller crashes. Fix the schema deviation.",
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
		ID:       "ZC1549",
		Title:    "Error on `unzip -d /` / `unzip -o ... -d /` — extract archive into filesystem root",
		Severity: SeverityError,
		Description: "Unzipping directly into `/` (or `/root`, `/boot`) overwrites any system file " +
			"whose path matches an entry in the archive. A malicious zip that carries " +
			"`etc/passwd`, `usr/bin/ls`, or `root/.ssh/authorized_keys` turns a seemingly " +
			"harmless extract into full system compromise. Stage to a scratch directory, " +
			"inspect contents, then copy or install specific files.",
		Check: checkZC1549,
	})
}

func checkZC1549(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "unzip" && ident.Value != "busybox" {
		return nil
	}

	var prevD bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevD {
			prevD = false
			if v == "/" || v == "/root" || v == "/boot" {
				return []Violation{{
					KataID: "ZC1549",
					Message: "`unzip -d " + v + "` extracts into a system path — any archive " +
						"entry overwrites matching system file. Stage, inspect, copy.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
		if v == "-d" {
			prevD = true
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1550",
		Title:    "Warn on `apt-mark hold <pkg>` — pins a package, blocks security updates",
		Severity: SeverityWarning,
		Description: "`apt-mark hold` tells apt to leave the package at its current version on " +
			"`apt upgrade` and `unattended-upgrades`. That is occasionally correct (pinning a " +
			"kernel variant for a driver, or a broken-upstream version) but silently keeps the " +
			"package vulnerable to every subsequent CVE. Document the reason in a comment, " +
			"schedule a review, and prefer `apt-mark unhold` + `apt upgrade <pkg>` over leaving " +
			"the pin in place indefinitely.",
		Check: checkZC1550,
	})
}

func checkZC1550(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "apt-mark" && ident.Value != "dpkg" {
		return nil
	}

	if ident.Value == "apt-mark" && len(cmd.Arguments) >= 2 &&
		cmd.Arguments[0].String() == "hold" {
		return []Violation{{
			KataID: "ZC1550",
			Message: "`apt-mark hold` pins the package — blocks future CVE fixes. Document " +
				"the reason and schedule an unhold review.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	// `echo "<pkg> hold" | dpkg --set-selections` is the legacy equivalent; flag when
	// dpkg is called with --set-selections.
	if ident.Value == "dpkg" {
		for _, arg := range cmd.Arguments {
			if arg.String() == "--set-selections" {
				return []Violation{{
					KataID: "ZC1550",
					Message: "`dpkg --set-selections` with a `hold` entry pins a package — " +
						"blocks future CVE fixes. Use apt-mark hold and document.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1551",
		Title:    "Warn on `helm install/upgrade --skip-crds` — chart CRs land before their CRDs",
		Severity: SeverityWarning,
		Description: "`--skip-crds` tells Helm to install only the `.Release` objects and skip " +
			"the CustomResourceDefinition manifests under `crds/`. Without the CRDs present, " +
			"any `.Release` object that references a custom resource is rejected by the API " +
			"server at validation time, or — worse — fails later when a reconciler tries to " +
			"watch a type that does not exist. Use the default (install CRDs) on first roll- " +
			"out; if you need split lifecycle, install CRDs manually (`kubectl apply -f " +
			"chart/crds/`) before the `helm install`.",
		Check: checkZC1551,
	})
}

func checkZC1551(node ast.Node) []Violation {
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

	var sawVerb bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "install" || v == "upgrade" {
			sawVerb = true
			continue
		}
		if !sawVerb {
			continue
		}
		if v == "--skip-crds" {
			return []Violation{{
				KataID: "ZC1551",
				Message: "`helm --skip-crds` installs .Release objects without their CRDs — " +
					"custom resources fail validation. Install CRDs first.",
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
		ID:       "ZC1552",
		Title:    "Warn on `openssl dhparam <2048` / `genrsa <2048` — weak key/parameter size",
		Severity: SeverityWarning,
		Description: "Generating DH parameters or RSA keys shorter than 2048 bits is below every " +
			"modern compliance baseline (NIST SP 800-57, BSI TR-02102, Mozilla Server Side TLS). " +
			"A 1024-bit RSA modulus or DH group is within reach of academic precomputation " +
			"(Logjam) and a 512-bit one was broken on commodity hardware in the 1990s. Use " +
			"2048 as a floor and 3072 / 4096 for long-lived keys.",
		Check: checkZC1552,
	})
}

func checkZC1552(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "openssl" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "dhparam" && sub != "genrsa" && sub != "gendsa" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n < 2048 {
			return []Violation{{
				KataID: "ZC1552",
				Message: "`openssl " + sub + " " + v + "` uses a weak key/param size — " +
					"modern baselines require 2048+. Use 2048 or 3072/4096 for long-lived keys.",
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
		ID:       "ZC1553",
		Title:    "Style: use Zsh `${(U)var}` / `${(L)var}` instead of `tr '[:lower:]' '[:upper:]'`",
		Severity: SeverityStyle,
		Description: "Zsh provides `${(U)var}` and `${(L)var}` parameter-expansion flags for " +
			"case conversion in-process. Spawning `tr` for this forks/execs per call (noticeable " +
			"in a hot loop), relies on the external `tr` being POSIX-compliant (BusyBox and old " +
			"macOS differ), and round-trips the data through a pipe. Drop `tr` for the " +
			"built-in: `upper=${(U)lower}` / `lower=${(L)upper}`.",
		Check: checkZC1553,
	})
}

func checkZC1553(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || CommandIdentifier(cmd) != "tr" {
		return nil
	}
	from, to := zc1553TrSets(cmd)
	if !zc1553IsCasePair(from, to) {
		return nil
	}
	return []Violation{{
		KataID: "ZC1553",
		Message: "`tr` for case conversion — use Zsh `${(U)var}` / `${(L)var}` to avoid " +
			"the fork/exec and portability hazard.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func zc1553TrSets(cmd *ast.SimpleCommand) (from, to string) {
	for _, a := range cmd.Arguments {
		v := strings.Trim(a.String(), "'\"")
		if strings.HasPrefix(v, "-") {
			continue
		}
		if from == "" {
			from = v
			continue
		}
		if to == "" {
			to = v
			return
		}
	}
	return
}

func zc1553IsCasePair(from, to string) bool {
	upper := from == "[:upper:]" || from == "A-Z"
	lower := from == "[:lower:]" || from == "a-z"
	if !upper && !lower {
		return false
	}
	other := to == "[:upper:]" || to == "A-Z" || to == "[:lower:]" || to == "a-z"
	return other && upper != (to == "[:upper:]" || to == "A-Z")
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1554",
		Title:    "Warn on `unzip -o` / `tar ... --overwrite` — silent overwrite during extract",
		Severity: SeverityWarning,
		Description: "`unzip -o` overwrites existing files without prompting; `tar --overwrite` " +
			"does the same for tarballs. In a directory that already contains user work or a " +
			"previous release, a newer archive silently wins, discarding in-flight edits and " +
			"custom config. Extract to a fresh staging directory, diff, then move specific " +
			"files into place.",
		Check: checkZC1554,
	})
}

func checkZC1554(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "unzip" {
		for _, arg := range cmd.Arguments {
			if arg.String() == "-o" {
				return []Violation{{
					KataID: "ZC1554",
					Message: "`unzip -o` overwrites existing files without prompting. Extract " +
						"to a staging directory, diff, then move.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	if ident.Value == "tar" || ident.Value == "bsdtar" {
		for _, arg := range cmd.Arguments {
			if arg.String() == "--overwrite" {
				return []Violation{{
					KataID: "ZC1554",
					Message: "`tar --overwrite` discards existing files during extract. Use a " +
						"staging directory and diff before rolling forward.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1555",
		Title:    "Error on `chmod` / `chown` on `/etc/shadow` or `/etc/sudoers` (managed files)",
		Severity: SeverityError,
		Description: "`/etc/shadow`, `/etc/gshadow`, `/etc/sudoers`, and `/etc/passwd` have " +
			"specific ownership and mode invariants that the distro `passwd`, `chage`, and " +
			"`visudo` tools maintain atomically with file locking. Direct `chmod`/`chown` races " +
			"those tools, can leave the file world-readable mid-modification (leaking the " +
			"shadow file), and will be clobbered on the next `shadow -p` run. Use the proper " +
			"wrapper, or ship a configuration-management drop-in.",
		Check: checkZC1555,
	})
}

func checkZC1555(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chmod" && ident.Value != "chown" && ident.Value != "chgrp" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "/etc/shadow", "/etc/gshadow", "/etc/sudoers", "/etc/passwd",
			"/etc/shadow-", "/etc/gshadow-", "/etc/passwd-", "/etc/sudoers-":
			return []Violation{{
				KataID: "ZC1555",
				Message: "`" + ident.Value + " ... " + v + "` races the distro-managed tool — " +
					"use passwd/chage/visudo or a config-management drop-in.",
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
		ID:       "ZC1556",
		Title:    "Error on `openssl enc -des` / `-rc4` / `-3des` — broken symmetric cipher",
		Severity: SeverityError,
		Description: "DES, RC4, and 3DES are all broken or on-deprecation-path: DES's 56-bit key " +
			"fell to commodity brute-force decades ago, RC4 has practical biased-output attacks, " +
			"and 3DES suffers the Sweet32 birthday collision when reused for more than ~32GB. " +
			"None of them provide authenticity either. Use `-aes-256-gcm` or `-chacha20-poly1305`, " +
			"or move up to a dedicated tool (`age`, `gpg`, `libsodium`).",
		Check: checkZC1556,
	})
}

func checkZC1556(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "openssl" {
		return nil
	}

	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "enc" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := strings.ToLower(arg.String())
		switch v {
		case "-des", "-des-cbc", "-des-ecb",
			"-des3", "-des-ede", "-des-ede-cbc", "-des-ede3", "-des-ede3-cbc",
			"-rc4", "-rc4-40",
			"-bf", "-bf-cbc",
			"-rc2", "-rc2-cbc",
			"-cast", "-cast5-cbc":
			return []Violation{{
				KataID: "ZC1556",
				Message: "`openssl enc " + v + "` is a broken or deprecated cipher. Use " +
					"`-aes-256-gcm` / `-chacha20-poly1305`, or `age` / `gpg` for files.",
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
		ID:       "ZC1557",
		Title:    "Error on `kubeadm reset -f` / `--force` — wipes Kubernetes control-plane state",
		Severity: SeverityError,
		Description: "`kubeadm reset` stops kubelet, tears down static-pod manifests, clears " +
			"`/etc/kubernetes`, and (with `-f`) skips the confirmation that protects a mistyped " +
			"target. On a control-plane node it also breaks every tenant that relied on that " +
			"etcd quorum. Drain first, remove the node from the cluster, then run reset " +
			"interactively to confirm.",
		Check: checkZC1557,
	})
}

func checkZC1557(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubeadm" {
		return nil
	}

	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "reset" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-f" || v == "--force" {
			return []Violation{{
				KataID: "ZC1557",
				Message: "`kubeadm reset -f` skips the confirmation and wipes " +
					"/etc/kubernetes / kubelet state. Drain and remove the node first.",
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
		ID:       "ZC1558",
		Title:    "Warn on `usermod -aG wheel|sudo|root|adm` — silent privilege group escalation",
		Severity: SeverityWarning,
		Description: "Adding a user to `wheel`, `sudo`, `root`, `adm`, `docker`, or `libvirt` " +
			"from a script grants persistent admin-level access without the review a sudoers " +
			"drop-in or PAM profile would get. `docker` and `libvirt` in particular are " +
			"equivalent to root (spawn privileged containers / raw disk access). Use a " +
			"sudoers.d file scoped to specific commands and audit changes in configuration " +
			"management.",
		Check: checkZC1558,
	})
}

var privGroups = map[string]struct{}{
	"wheel":   {},
	"sudo":    {},
	"root":    {},
	"adm":     {},
	"docker":  {},
	"libvirt": {},
	"lxd":     {},
	"kvm":     {},
	"disk":    {},
}

var (
	zc1558UserModTools = map[string]struct{}{"usermod": {}, "gpasswd": {}, "adduser": {}}
	zc1558GroupFlags   = map[string]struct{}{
		"-aG": {}, "-Ga": {}, "-G": {}, "--groups": {}, "--append": {},
	}
)

func checkZC1558(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	tool := CommandIdentifier(cmd)
	if _, hit := zc1558UserModTools[tool]; !hit {
		return nil
	}
	args := zc1464StringArgs(cmd)
	if g := zc1558PrivGroupViaFlag(args); g != "" {
		return zc1558Violation(cmd, g)
	}
	if tool == "gpasswd" {
		if g := zc1558GpasswdAddPrivGroup(args); g != "" {
			return zc1558Violation(cmd, g)
		}
	}
	return nil
}

func zc1558PrivGroupViaFlag(args []string) string {
	for i, a := range args {
		if _, hit := zc1558GroupFlags[a]; !hit || i+1 >= len(args) {
			continue
		}
		for _, g := range strings.Split(args[i+1], ",") {
			g = strings.TrimSpace(g)
			if _, bad := privGroups[g]; bad {
				return g
			}
		}
	}
	return ""
}

func zc1558GpasswdAddPrivGroup(args []string) string {
	if len(args) < 3 || args[0] != "-a" {
		return ""
	}
	g := args[2]
	if _, bad := privGroups[g]; bad {
		return g
	}
	return ""
}

func zc1558Violation(cmd *ast.SimpleCommand, group string) []Violation {
	return []Violation{{
		KataID: "ZC1558",
		Message: "Adding user to `" + group + "` grants persistent admin-level access — use a " +
			"scoped sudoers.d drop-in via configuration management.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1559",
		Title:    "Warn on `ssh-copy-id -f` / `-o StrictHostKeyChecking=no` — trust-on-first-use key push",
		Severity: SeverityWarning,
		Description: "`ssh-copy-id` opens an SSH connection to deposit the caller's public key. " +
			"With `-f` it overwrites existing `authorized_keys` without prompting; with " +
			"`-o StrictHostKeyChecking=no` it does not verify the host key. Together they " +
			"push a long-term credential at a host the script has never authenticated — a " +
			"network MITM lands a permanent backdoor. Verify the target host's fingerprint " +
			"out of band before pushing keys.",
		Check: checkZC1559,
	})
}

func checkZC1559(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh-copy-id" {
		return nil
	}

	var prevO bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-f" {
			return zc1559Violation(cmd, "-f")
		}
		if prevO {
			prevO = false
			s := strings.TrimSpace(strings.ToLower(v))
			if s == "stricthostkeychecking=no" || s == "userknownhostsfile=/dev/null" {
				return zc1559Violation(cmd, "-o "+v)
			}
		}
		if v == "-o" {
			prevO = true
		}
	}
	return nil
}

func zc1559Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1559",
		Message: "`ssh-copy-id " + what + "` pushes a long-term credential without host-key " +
			"verification. Verify the fingerprint out of band first.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1560",
		Title:    "Error on `pip install --trusted-host` — accepts MITM / plain-HTTP PyPI index",
		Severity: SeverityError,
		Description: "`--trusted-host` tells pip to skip TLS certificate verification for the " +
			"specified host and to allow plain-HTTP URLs from that host. Any MITM on the path " +
			"can substitute packages on install, and a typo in the host name means every " +
			"subsequent `install` from the misspelled host is unauthenticated. Fix the CA " +
			"trust (install the real corporate CA) instead of silencing pip, and keep the " +
			"default `--index-url https://...` over the TLS-verified endpoint.",
		Check: checkZC1560,
	})
}

func checkZC1560(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "pip" && ident.Value != "pip3" && ident.Value != "pipx" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "--trusted-host" {
			return []Violation{{
				KataID: "ZC1560",
				Message: "`pip --trusted-host` skips TLS verification and allows plain-HTTP " +
					"for that index. Fix the CA trust and keep --index-url on https://.",
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
		ID:       "ZC1561",
		Title:    "Error on `systemctl isolate rescue.target` / `emergency.target` from a script",
		Severity: SeverityError,
		Description: "`systemctl isolate rescue.target` drops the host into single-user rescue " +
			"mode; `emergency.target` goes even further, leaving only the root shell on the " +
			"console. Both terminate networking, SSH sessions, and most services. On a remote " +
			"host the script loses its own connection mid-run, and anyone relying on the box " +
			"is cut off without warning. Reserve these for console recovery, not script flow.",
		Check: checkZC1561,
	})
}

func checkZC1561(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "systemctl" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "isolate" {
		return nil
	}
	target := cmd.Arguments[1].String()
	switch target {
	case "rescue.target", "emergency.target", "rescue", "emergency":
		return []Violation{{
			KataID: "ZC1561",
			Message: "`systemctl isolate " + target + "` kills SSH and most services — " +
				"console-only recovery. Do not run from a script.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1562",
		Title:    "Warn on `env -u PATH` / `-u LD_LIBRARY_PATH` — clears security-relevant env",
		Severity: SeverityWarning,
		Description: "`env -u PATH` unsets the caller's `PATH` before running the child, forcing " +
			"the child to fall back to the hard-coded search list (`/bin:/usr/bin` on glibc). " +
			"That bypasses PATH hardening done by the parent shell (e.g. a sanitised PATH " +
			"under `sudo`). Unsetting `LD_PRELOAD` / `LD_LIBRARY_PATH` mid-stream is also " +
			"usually the caller trying to shake off an earlier `export`. Either use `env -i` " +
			"to sanitise completely, or explicitly set the variables the child should see.",
		Check: checkZC1562,
	})
}

func checkZC1562(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "env" {
		return nil
	}

	var prevU bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevU {
			prevU = false
			switch v {
			case "PATH", "LD_PRELOAD", "LD_LIBRARY_PATH", "LD_AUDIT":
				return []Violation{{
					KataID: "ZC1562",
					Message: "`env -u " + v + "` clears a security-relevant variable mid-run. " +
						"Use `env -i` to sanitise, or set the right value explicitly.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
		if v == "-u" || v == "--unset" {
			prevU = true
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1563",
		Title:    "Warn on `swapoff -a` — disables swap (memory pressure, potential OOM)",
		Severity: SeverityWarning,
		Description: "`swapoff -a` turns off every active swap device. Kubelet installers do " +
			"this because kubelet refuses to run with swap, but leaving it in a general-purpose " +
			"script means the next memory-hungry process on the host hits the OOM killer " +
			"instead of paging. If the goal is kubelet-friendly, also remove the swap entry " +
			"from `/etc/fstab` and document the trade-off; otherwise keep swap on.",
		Check: checkZC1563,
	})
}

func checkZC1563(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "swapoff" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-a" || arg.String() == "--all" {
			return []Violation{{
				KataID: "ZC1563",
				Message: "`swapoff -a` disables all swap devices — next memory-hungry process " +
					"hits OOM. Document the trade-off if kubelet requires it.",
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
		ID:       "ZC1564",
		Title:    "Warn on `date -s` / `timedatectl set-time` — manual clock change breaks TLS / cron",
		Severity: SeverityWarning,
		Description: "Setting the system clock by hand (`date -s`, `timedatectl set-time`, " +
			"`hwclock --set`) moves wall-clock time enough to invalidate short-lived TLS " +
			"certificates, reset `cron`'s missed-job catch-up, and confuse `systemd.timer` " +
			"units that depend on monotonic math. Use `systemd-timesyncd` / `chrony` / `ntpd` " +
			"for routine correction; reserve manual set for first-boot bootstrap or air-gapped " +
			"recovery and document the action.",
		Check: checkZC1564,
	})
}

func checkZC1564(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "date" {
		for _, arg := range cmd.Arguments {
			if arg.String() == "-s" || arg.String() == "--set" {
				return zc1564Violation(cmd, "date -s")
			}
		}
	}
	if ident.Value == "timedatectl" && len(cmd.Arguments) >= 1 &&
		cmd.Arguments[0].String() == "set-time" {
		return zc1564Violation(cmd, "timedatectl set-time")
	}
	if ident.Value == "hwclock" {
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "--set" || v == "-w" || v == "--systohc" {
				return zc1564Violation(cmd, "hwclock "+v)
			}
		}
	}
	return nil
}

func zc1564Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1564",
		Message: "`" + what + "` sets the wall clock manually — breaks TLS certs, cron " +
			"catch-up, and systemd timer math. Use timesyncd/chrony/ntpd.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1565",
		Title:    "Style: use `command -v` instead of `whereis` / `locate` for command existence",
		Severity: SeverityStyle,
		Description: "`whereis` searches a hard-coded list of binary/manual/source directories " +
			"and returns everything it finds, including stale paths on custom `$PATH` layouts. " +
			"`locate` relies on a cron-maintained index that may be hours or days stale. For " +
			"a scripted \"does this command exist?\" check, `command -v <cmd>` respects the " +
			"current `$PATH`, returns the selected resolution, and has no index-refresh " +
			"coupling.",
		Check: checkZC1565,
		Fix:   fixZC1565,
	})
}

// fixZC1565 rewrites a `whereis` / `locate` / `mlocate` / `plocate`
// command-name lookup into `command -v`. The detector restricts to the
// four index-based forms so the swap is safe; arguments stay untouched.
func fixZC1565(node ast.Node, v Violation, _ []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "whereis", "locate", "mlocate", "plocate":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len(ident.Value),
			Replace: "command -v",
		}}
	}
	return nil
}

func checkZC1565(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "whereis" && ident.Value != "locate" && ident.Value != "mlocate" &&
		ident.Value != "plocate" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1565",
		Message: "`" + ident.Value + "` is index-based and stale-prone. Use `command -v " +
			"<cmd>` for runtime existence checks.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1566",
		Title:    "Error on `gem install -P NoSecurity|LowSecurity` / `--trust-policy NoSecurity`",
		Severity: SeverityError,
		Description: "RubyGems' trust policy decides what signatures the installer accepts. " +
			"`NoSecurity` skips signature verification entirely; `LowSecurity` warns but still " +
			"installs unsigned gems. On a registry MITM or a hijacked maintainer account those " +
			"policies turn into arbitrary code execution at gem-install time. Use `HighSecurity` " +
			"(reject all but fully-signed) or `MediumSecurity` for hybrid repos.",
		Check: checkZC1566,
	})
}

func checkZC1566(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "gem" {
		return nil
	}

	var prevP bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevP {
			prevP = false
			if v == "NoSecurity" || v == "LowSecurity" {
				return []Violation{{
					KataID: "ZC1566",
					Message: "`gem -P " + v + "` skips signature verification — MITM or " +
						"account compromise becomes RCE at install. Use HighSecurity.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
		if v == "-P" || v == "--trust-policy" {
			prevP = true
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1567",
		Title:    "Warn on `python -m http.server` without `--bind 127.0.0.1` — serves to all interfaces",
		Severity: SeverityWarning,
		Description: "`python -m http.server` (and the legacy `SimpleHTTPServer`) default to " +
			"`0.0.0.0`, exposing the current directory's contents to every network the host " +
			"is on. Tmp scratch files, `.env`, SSH keys, or a `node_modules` tree with private " +
			"config all become reachable from anywhere on the LAN (or the internet, on a VPS). " +
			"Pass `--bind 127.0.0.1` (or `--bind ::1`) unless you really need external access " +
			"and know what is in the cwd.",
		Check: checkZC1567,
	})
}

func checkZC1567(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "python" && ident.Value != "python2" && ident.Value != "python3" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// python -m http.server  or  python -m SimpleHTTPServer
	var isServer bool
	for i := 0; i+1 < len(args); i++ {
		if args[i] == "-m" && (args[i+1] == "http.server" || args[i+1] == "SimpleHTTPServer") {
			isServer = true
			break
		}
	}
	if !isServer {
		return nil
	}

	for _, a := range args {
		if a == "--bind" || a == "-b" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1567",
		Message: "`python -m http.server` without `--bind` defaults to 0.0.0.0 — exposes the " +
			"cwd to every network the host sees. Add `--bind 127.0.0.1`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1568",
		Title:    "Error on `useradd -o` / `usermod -o` — allows non-unique UID (alias user)",
		Severity: SeverityError,
		Description: "`-o` (or `--non-unique`) lets `useradd` / `usermod` assign a UID that is " +
			"already in use. The new account has the same kernel identity as the existing one " +
			"but its own login name, password, shell, and home dir. It is indistinguishable in " +
			"`ps` / audit / file ACLs, so a compromise of either account is a compromise of " +
			"both. Pick a fresh UID instead.",
		Check: checkZC1568,
	})
}

func checkZC1568(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "useradd" && ident.Value != "usermod" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-o" || v == "--non-unique" {
			return []Violation{{
				KataID: "ZC1568",
				Message: "`" + ident.Value + " -o` assigns a non-unique UID — the two " +
					"accounts share kernel identity, indistinguishable in audit. Use a " +
					"fresh UID.",
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
		ID:       "ZC1569",
		Title:    "Error on `nvme format -s1` / `-s2` — cryptographic or full-block SSD erase",
		Severity: SeverityError,
		Description: "`nvme format -s1` does a cryptographic erase of the target namespace; " +
			"`-s2` (or the full-NVMe sanitize) rewrites every block. Both are unrecoverable " +
			"in seconds. On a typo in the device variable — or a script that iterates over " +
			"`/dev/nvme*n*` and catches the wrong namespace — the wrong disk is gone by the " +
			"time the operator notices. Run interactively on verified targets, or not at all " +
			"from automation.",
		Check: checkZC1569,
	})
}

func checkZC1569(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "nvme" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "format" && sub != "sanitize" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-s1" || v == "-s2" || v == "--ses=1" || v == "--ses=2" ||
			v == "-a" || v == "--sanact" {
			return []Violation{{
				KataID: "ZC1569",
				Message: "`nvme " + sub + " " + v + "` unrecoverably erases the namespace in " +
					"seconds. Do not run from automation.",
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
		ID:       "ZC1570",
		Title:    "Warn on `smbclient -N` / `mount.cifs guest` — anonymous SMB share access",
		Severity: SeverityWarning,
		Description: "`smbclient -N` skips authentication entirely (anonymous / null session); " +
			"`mount.cifs` with `guest,username=` or `-o guest` does the same at the mount " +
			"layer. Any host on the network segment can then read the share. If the share is " +
			"truly public (software mirror, build cache) wrap in a read-only filesystem and " +
			"document it; otherwise require Kerberos (`-k`) or pass credentials via " +
			"`credentials=<file>` with 0600 perms.",
		Check: checkZC1570,
	})
}

func checkZC1570(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "smbclient" && ident.Value != "mount.cifs" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-N" || v == "--no-pass" {
			return []Violation{{
				KataID: "ZC1570",
				Message: "`" + ident.Value + " " + v + "` is anonymous SMB access — any " +
					"host on-net can read the share. Use credentials=<file> 0600 or -k.",
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
		ID:       "ZC1571",
		Title:    "Style: `ntpdate` is deprecated — use `chronyc makestep` / `systemd-timesyncd`",
		Severity: SeverityStyle,
		Description: "`ntpdate` was retired by the ntp.org project around 4.2.6. Distros " +
			"increasingly ship without it, and packaging it breaks the invariant that only " +
			"one program writes the clock at a time (if `chrony` or `timesyncd` is also " +
			"running the two fight). Use `chronyc makestep` (if chrony is active) or " +
			"`systemctl restart systemd-timesyncd` (if timesyncd is active) for a one-shot " +
			"step, and leave the daemon to keep it synchronised.",
		Check: checkZC1571,
	})
}

func checkZC1571(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ntpdate" && ident.Value != "sntp" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1571",
		Message: "`" + ident.Value + "` is deprecated and races any running chrony/timesyncd. " +
			"Use `chronyc makestep` or `systemctl restart systemd-timesyncd`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1572",
		Title:    "Warn on `docker run -e PASSWORD=<value>` — secret in container env / inspect",
		Severity: SeverityWarning,
		Description: "Passing a secret through `docker run -e NAME=value` puts it in the output " +
			"of `docker inspect`, the container's `/proc/1/environ` (readable by anything that " +
			"shares the PID namespace), and the shell history of whoever launched the " +
			"container. Use `--env-file` with 0600 perms, a secret-mount `--secret` via " +
			"BuildKit / Swarm, or mount a tmpfs file the container reads at runtime.",
		Check: checkZC1572,
	})
}

var secretEnvPrefixes = []string{
	"PASSWORD", "PASSWD", "PASS",
	"SECRET", "SECRET_KEY", "API_KEY",
	"TOKEN", "AUTH_TOKEN", "ACCESS_TOKEN",
	"PRIVATE_KEY", "DB_PASSWORD", "DB_PASS",
	"AWS_SECRET", "AWS_SECRET_ACCESS_KEY",
	"GITHUB_TOKEN", "GH_TOKEN", "NPM_TOKEN",
}

var zc1572Runtimes = map[string]struct{}{"docker": {}, "podman": {}, "nerdctl": {}}

func checkZC1572(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if _, hit := zc1572Runtimes[CommandIdentifier(cmd)]; !hit {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "run" {
		return nil
	}
	if form, name := zc1572FirstSecretEnv(cmd.Arguments); name != "" {
		return []Violation{{
			KataID: "ZC1572",
			Message: "`" + form + name + "=<value>` writes the secret into `docker " +
				"inspect` and `/proc/1/environ`. Use `--env-file` 0600 or " +
				"`--secret`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func zc1572FirstSecretEnv(args []ast.Expression) (form, name string) {
	prevE := false
	for _, arg := range args[1:] {
		v := arg.String()
		if prevE {
			prevE = false
			if n := zc1572SecretAssignment(v); n != "" {
				return "-e ", n
			}
			continue
		}
		if v == "-e" || v == "--env" {
			prevE = true
			continue
		}
		if strings.HasPrefix(v, "--env=") {
			if n := zc1572SecretAssignment(v[len("--env="):]); n != "" {
				return "--env=", n
			}
		}
	}
	return "", ""
}

func zc1572SecretAssignment(v string) string {
	name, value, ok := strings.Cut(v, "=")
	if !ok || value == "" {
		return ""
	}
	if !looksLikeSecret(name) {
		return ""
	}
	return name
}

func looksLikeSecret(name string) bool {
	up := strings.ToUpper(name)
	for _, p := range secretEnvPrefixes {
		if up == p || strings.Contains(up, p) {
			return true
		}
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1573",
		Title:    "Warn on `chattr -i` / `chattr -a` — removes immutable / append-only attribute",
		Severity: SeverityWarning,
		Description: "Removing the immutable (`-i`) or append-only (`-a`) attribute lets the " +
			"file be overwritten or truncated again. When the target is a log file, shadow " +
			"file, or hardened system binary, that flag was explicitly set to make tampering " +
			"noisy. Removing it mid-script is either a one-shot upgrade (follow with the " +
			"`chattr +i` restore) or an anti-forensics step. If it is the former, wrap the " +
			"change in a function and re-set the attribute at the end.",
		Check: checkZC1573,
	})
}

func checkZC1573(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chattr" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-i" || v == "-a" || v == "-ia" || v == "-ai" {
			return []Violation{{
				KataID: "ZC1573",
				Message: "`chattr " + v + "` removes the tamper-evident attribute. If this " +
					"is a one-shot upgrade, re-set the attribute at the end of the block.",
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
		ID:       "ZC1574",
		Title:    "Warn on `git config credential.helper store` — plaintext credentials on disk",
		Severity: SeverityWarning,
		Description: "`credential.helper store` writes the username and password to " +
			"`~/.git-credentials` in plaintext. Anything that backs up that file (rsync, " +
			"imaging, cloud sync) then carries the credential around. Use a platform helper " +
			"instead: `manager` / `manager-core` on Windows / Mac, `libsecret` on Linux, or " +
			"`cache --timeout=3600` for short-lived in-memory caching.",
		Check: checkZC1574,
	})
}

func checkZC1574(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "git" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// git config [scope] credential.helper store
	for i, a := range args {
		if a != "config" {
			continue
		}
		// Walk past optional scope flag.
		j := i + 1
		for j < len(args) && strings.HasPrefix(args[j], "--") && args[j] != "--" {
			j++
		}
		if j+1 < len(args) && args[j] == "credential.helper" {
			v := args[j+1]
			// Accept bare `store` or quoted `store --file=...`.
			if v == "store" || strings.HasPrefix(v, "store ") ||
				strings.Contains(v, "store --file=") ||
				strings.Contains(v, "'store") || strings.Contains(v, "\"store") {
				return []Violation{{
					KataID: "ZC1574",
					Message: "`git credential.helper store` saves credentials in plaintext — " +
						"backups leak the token. Use platform helper (manager-core / " +
						"libsecret) or `cache --timeout=<sec>`.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1575",
		Title:    "Error on `aws configure set aws_secret_access_key <value>` — secret on cmdline",
		Severity: SeverityError,
		Description: "`aws configure set aws_secret_access_key …` writes the secret access key " +
			"into `~/.aws/credentials` and leaves the raw value in `ps` / shell history until " +
			"the process exits. On a shared CI runner or a multi-user host, that window is " +
			"long enough for a co-tenant to snapshot the key. Use IAM-role-based auth (EC2 " +
			"instance profile, IRSA on EKS, OIDC from GitHub / GitLab) or read the value from " +
			"stdin / a 0600 file and let `aws configure` import it interactively.",
		Check: checkZC1575,
	})
}

func checkZC1575(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "aws" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// aws configure set aws_secret_access_key VALUE
	for i := 0; i+3 < len(args); i++ {
		if args[i] == "configure" && args[i+1] == "set" {
			key := strings.ToLower(args[i+2])
			if key == "aws_secret_access_key" || key == "aws_session_token" ||
				strings.Contains(key, "secret") {
				return []Violation{{
					KataID: "ZC1575",
					Message: "`aws configure set " + args[i+2] + " …` puts the secret in " +
						"ps / history. Use IAM-role auth or import from stdin / 0600 file.",
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
		ID:       "ZC1576",
		Title:    "Warn on `terraform apply -target=...` — cherry-pick apply bypasses dependencies",
		Severity: SeverityWarning,
		Description: "`-target` restricts `terraform apply` to a specific resource / module and " +
			"everything it depends on. In theory that is a surgical fix; in practice it " +
			"routinely skips changes the targeted resource actually depends on, leading to " +
			"drift between state and configuration. HashiCorp documents `-target` as a tool " +
			"for incident response, not routine operations. Re-run without `-target` or " +
			"split the configuration into separate root modules.",
		Check: checkZC1576,
	})
}

func checkZC1576(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "terraform" && ident.Value != "terragrunt" && ident.Value != "tofu" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "apply" && sub != "destroy" && sub != "plan" {
		return nil
	}

	var prevTarget bool
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if prevTarget {
			return zc1576Violation(cmd, "-target "+v)
		}
		if v == "-target" {
			prevTarget = true
			continue
		}
		if strings.HasPrefix(v, "-target=") {
			return zc1576Violation(cmd, v)
		}
	}
	return nil
}

func zc1576Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1576",
		Message: "`terraform " + what + "` bypasses dependency order — documented as incident " +
			"response tool only. Re-run without -target or split root modules.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1577",
		Title:    "Warn on `dig <name> ANY` — deprecated query type (RFC 8482)",
		Severity: SeverityWarning,
		Description: "ANY queries return whatever the authoritative server feels like sending " +
			"back — or just the HINFO placeholder mandated by RFC 8482. Modern recursors " +
			"filter ANY to avoid reflection-amplification abuse, so scripts that rely on ANY " +
			"for enumeration get inconsistent or empty results. Query the specific record " +
			"types you want (`dig A name`, `dig MX name`, `dig NS name`) and combine them.",
		Check: checkZC1577,
	})
}

func checkZC1577(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dig" && ident.Value != "drill" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "ANY" || v == "any" {
			return []Violation{{
				KataID: "ZC1577",
				Message: "`" + ident.Value + " ... ANY` is RFC 8482-deprecated — filtered by " +
					"recursors. Query specific types (A / MX / NS / …) and combine.",
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
		ID:       "ZC1578",
		Title:    "Warn on `ssh-keygen -b <2048` for RSA / DSA — weak SSH key",
		Severity: SeverityWarning,
		Description: "Generating an SSH RSA or DSA key shorter than 2048 bits fails current " +
			"OpenSSH baselines and is rejected by recent `ssh` versions when used for " +
			"authentication. DSA was removed from OpenSSH 9.8 outright. Use `ssh-keygen -t " +
			"ed25519` (compact, fast, modern defaults) or `ssh-keygen -t rsa -b 4096` if you " +
			"need RSA for compatibility.",
		Check: checkZC1578,
	})
}

func checkZC1578(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh-keygen" {
		return nil
	}

	var keyType string
	var bits int
	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-t" && i+1 < len(cmd.Arguments) {
			keyType = cmd.Arguments[i+1].String()
		}
		if v == "-b" && i+1 < len(cmd.Arguments) {
			n, err := strconv.Atoi(cmd.Arguments[i+1].String())
			if err == nil {
				bits = n
			}
		}
	}

	// DSA regardless of size is weak / removed.
	if keyType == "dsa" {
		return []Violation{{
			KataID:  "ZC1578",
			Message: "`ssh-keygen -t dsa` — DSA removed from OpenSSH 9.8. Use `-t ed25519`.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityWarning,
		}}
	}
	// RSA below 2048 bits is weak.
	if (keyType == "rsa" || keyType == "") && bits > 0 && bits < 2048 {
		return []Violation{{
			KataID: "ZC1578",
			Message: "`ssh-keygen -b " + strconv.Itoa(bits) + "` — RSA below 2048 bits is " +
				"rejected by modern OpenSSH. Use `-t ed25519` or `-b 4096`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1579",
		Title:    "Warn on `curl --retry-all-errors` without `--max-time` — hammers endpoint on failure",
		Severity: SeverityWarning,
		Description: "`--retry-all-errors` (curl 7.71+) treats every HTTP error as retryable. " +
			"Without `--max-time` capping total wall clock, a server that responds `500` " +
			"quickly gets hit back-to-back until `--retry` exhausts — a mini-DoS against your " +
			"own upstream, especially if the script itself is scheduled on many nodes. Pair " +
			"with `--max-time <seconds>` or prefer `--retry-connrefused` (only retries " +
			"connection-level failures).",
		Check: checkZC1579,
	})
}

func checkZC1579(node ast.Node) []Violation {
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

	var hasRetryAll, hasMaxTime bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--retry-all-errors" {
			hasRetryAll = true
		}
		if v == "--max-time" || v == "-m" {
			hasMaxTime = true
		}
	}
	if !hasRetryAll || hasMaxTime {
		return nil
	}
	return []Violation{{
		KataID: "ZC1579",
		Message: "`curl --retry-all-errors` with no `--max-time` hammers the upstream on " +
			"failure. Pair with `-m <seconds>` or use `--retry-connrefused`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1580",
		Title:    "Warn on `go build -ldflags \"-X main.<SECRET>=...\"` — secret embedded in binary",
		Severity: SeverityWarning,
		Description: "`-ldflags=\"-X pkg.Var=value\"` sets a Go string variable at link time. " +
			"Putting a secret here bakes it into the resulting binary (discoverable with " +
			"`strings`, `objdump`, or simply opening the file). It also leaves the value on " +
			"the build host's shell history and in any CI transcript. Read the value at " +
			"runtime from `os.Getenv` / a mounted secret file / the cloud secret manager.",
		Check: checkZC1580,
	})
}

func checkZC1580(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "go" {
		return nil
	}

	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "build" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := strings.Trim(arg.String(), "\"'")
		// Look for -ldflags=... content or "-X ..." content
		if !strings.Contains(v, "-X ") && !strings.Contains(v, "-ldflags") {
			continue
		}
		up := strings.ToUpper(v)
		if strings.Contains(up, "PASSWORD=") || strings.Contains(up, "SECRET=") ||
			strings.Contains(up, "APIKEY=") || strings.Contains(up, "API_KEY=") ||
			strings.Contains(up, "TOKEN=") || strings.Contains(up, "PRIVATE_KEY=") {
			return []Violation{{
				KataID: "ZC1580",
				Message: "`go build -ldflags` injecting a secret bakes it into the binary. " +
					"Read from os.Getenv / mounted secret file at runtime.",
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
		ID:       "ZC1581",
		Title:    "Warn on `ssh -o PubkeyAuthentication=no` / `-o PasswordAuthentication=yes`",
		Severity: SeverityWarning,
		Description: "Forcing password authentication on a connection that has a working key " +
			"turns a strong (challenge-response, no password leaves the client) into a weak " +
			"(password-in-the-clear-on-disk-or-prompt) authentication path. Similarly " +
			"disabling pubkey skips the good path entirely. Leave the defaults, let the " +
			"server's `PubkeyAuthentication yes` pick the key, and document any exception.",
		Check: checkZC1581,
	})
}

func checkZC1581(node ast.Node) []Violation {
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

	var prevO bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevO {
			prevO = false
			s := strings.TrimSpace(strings.ToLower(v))
			if s == "pubkeyauthentication=no" || s == "passwordauthentication=yes" ||
				s == "preferredauthentications=password" {
				return []Violation{{
					KataID: "ZC1581",
					Message: "`" + ident.Value + " -o " + v + "` forces password auth — " +
						"weaker than key auth. Let the default preference pick.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
		if v == "-o" {
			prevO = true
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1582",
		Title:    "Warn on `bash -x` / `sh -x` / `zsh -x` — traces every command, leaks secrets",
		Severity: SeverityWarning,
		Description: "`-x` turns on xtrace, printing every command (expanded) to stderr before " +
			"it runs. In a CI log that is indexed / shared / archived, any line that touches " +
			"a secret leaks it verbatim — `curl` with a `Bearer` header, `psql` with a " +
			"password, `echo $API_TOKEN > ...`. If you really need tracing, wrap the non-" +
			"secret block with `set -x; ...; set +x` and exclude the secret-handling parts.",
		Check: checkZC1582,
	})
}

func checkZC1582(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "bash" && ident.Value != "sh" && ident.Value != "zsh" &&
		ident.Value != "dash" && ident.Value != "ksh" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-x" || v == "-xv" || v == "-vx" {
			return []Violation{{
				KataID: "ZC1582",
				Message: "`" + ident.Value + " " + v + "` traces every expanded command — CI logs " +
					"leak secrets verbatim. Scope with `set -x; …; set +x`.",
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
		ID:       "ZC1583",
		Title:    "Warn on `find ... -delete` without `-maxdepth` — unbounded recursive delete",
		Severity: SeverityWarning,
		Description: "`find PATH -delete` walks the tree recursively and removes every match. " +
			"Without `-maxdepth N` the walk crosses into every subtree, including symlinks " +
			"that point outside the intended scope and mount points that expand the blast " +
			"radius. Scope the depth (`-maxdepth 2`) and prefer a dry-run first " +
			"(`find ... -print | head`).",
		Check: checkZC1583,
	})
}

func checkZC1583(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "find" {
		return nil
	}

	var hasDelete, hasMaxdepth, hasPrune, hasXdev bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-delete":
			hasDelete = true
		case "-maxdepth":
			hasMaxdepth = true
		case "-prune":
			hasPrune = true
		case "-xdev", "-mount":
			hasXdev = true
		}
	}
	if !hasDelete || hasMaxdepth || hasPrune || hasXdev {
		return nil
	}
	return []Violation{{
		KataID: "ZC1583",
		Message: "`find -delete` without `-maxdepth` / `-xdev` / `-prune` walks the whole " +
			"tree. Scope the depth (e.g. `-maxdepth 2`) and dry-run first.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1584",
		Title:    "Warn on `sudo -E` / `--preserve-env` — carries caller env into root shell",
		Severity: SeverityWarning,
		Description: "`sudo -E` preserves the invoking user's environment — `PATH`, " +
			"`LD_PRELOAD`, `PYTHONPATH`, etc. On a workstation where the user has a personal " +
			"`~/bin` early in `$PATH`, any wrapper named like a system binary gets executed " +
			"by the privileged process. That is exactly the sudoers `secure_path` mechanic " +
			"fails to protect against. Whitelist specific variables with `env_keep` in " +
			"sudoers, or call `sudo env VAR=value cmd` with the minimum.",
		Check: checkZC1584,
	})
}

func checkZC1584(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "sudo" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-E" || v == "--preserve-env" {
			return []Violation{{
				KataID: "ZC1584",
				Message: "`sudo " + v + "` carries the caller's PATH / LD_PRELOAD / … into " +
					"the privileged process. Use `env_keep` in sudoers or explicit `sudo " +
					"env VAR=… cmd`.",
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
		ID:       "ZC1585",
		Title:    "Warn on `ufw allow from any` / `firewall-cmd --add-source=0.0.0.0/0`",
		Severity: SeverityWarning,
		Description: "`ufw allow from any to any port …` (and its firewall-cmd sibling " +
			"`--add-source=0.0.0.0/0`) opens the port to the whole internet. That is " +
			"sometimes the point (public HTTP / HTTPS), but on management ports (22, 3306, " +
			"5432, 6379, 9200, 27017) it is a routine foot-gun when the script author " +
			"assumed the host would only ever be reached via VPN. Scope the rule to a " +
			"specific source CIDR.",
		Check: checkZC1585,
	})
}

func checkZC1585(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ufw" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	if len(args) < 3 || args[0] != "allow" {
		return nil
	}
	// ufw allow from any ...
	for i := 1; i+1 < len(args); i++ {
		if args[i] == "from" && args[i+1] == "any" {
			return []Violation{{
				KataID: "ZC1585",
				Message: "`ufw allow from any …` opens the port to the whole internet. " +
					"Scope to a specific source CIDR.",
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
		ID:       "ZC1586",
		Title:    "Style: `chkconfig` / `update-rc.d` / `insserv` — SysV init relics, use `systemctl`",
		Severity: SeverityStyle,
		Description: "`chkconfig` (Red Hat), `update-rc.d` (Debian), and `insserv` (SUSE) are " +
			"SysV-init compatibility wrappers for enabling/disabling services at boot. On any " +
			"distro that has used systemd for the last decade they are translated to " +
			"`systemctl enable|disable`, but silently lose unit-template arguments, " +
			"`[Install]` alias handling, and socket-activated services. Call `systemctl " +
			"enable <unit>` directly.",
		Check: checkZC1586,
	})
}

func checkZC1586(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chkconfig" && ident.Value != "update-rc.d" &&
		ident.Value != "insserv" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1586",
		Message: "`" + ident.Value + "` is a SysV-init relic. Use `systemctl enable|disable " +
			"<unit>` directly.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1587",
		Title:    "Warn on `modprobe -r` / `rmmod` from scripts — unloading active kernel modules",
		Severity: SeverityWarning,
		Description: "Unloading a kernel module that is in use — `nvme` (storage), `nvidia` " +
			"(GPU), `e1000`/`ixgbe` (network), `kvm` (virt) — instantly takes the backing " +
			"subsystem offline. On a remote host the script loses its storage or network " +
			"mid-run. Reserve `modprobe -r` / `rmmod` for console maintenance, and consider " +
			"`systemctl stop <unit>` if you are trying to stop a service.",
		Check: checkZC1587,
	})
}

func checkZC1587(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "modprobe" {
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-r" || v == "--remove" {
				return []Violation{{
					KataID: "ZC1587",
					Message: "`modprobe -r` unloads an in-use module — the backing subsystem " +
						"goes offline. Use `systemctl stop` if you meant to stop a service.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	if ident.Value == "rmmod" && len(cmd.Arguments) > 0 {
		return []Violation{{
			KataID: "ZC1587",
			Message: "`rmmod` unloads a kernel module — the backing subsystem goes offline. " +
				"Use `systemctl stop` if you meant to stop a service.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1588",
		Title:    "Error on `nsenter --target 1` — joins host init namespaces (container escape)",
		Severity: SeverityError,
		Description: "`nsenter -t 1` attaches to the namespaces of pid 1. Inside a privileged " +
			"container or one with `CAP_SYS_ADMIN`, pid 1 is the host init — joining its " +
			"mount / pid / net / uts / ipc namespaces is the canonical escape primitive. " +
			"From that new shell the caller sees and writes the host filesystem, kills host " +
			"processes, and hijacks host network. Legit debugging runs from the host, not from " +
			"inside the container. If you need to exec into a container, use `docker exec` / " +
			"`kubectl exec`.",
		Check: checkZC1588,
	})
}

func checkZC1588(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "nsenter" {
		return nil
	}

	var expectTarget, hit bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--target=1" || v == "-t1" {
			hit = true
			break
		}
		if expectTarget && v == "1" {
			hit = true
			break
		}
		expectTarget = v == "-t" || v == "--target"
	}
	if !hit {
		return nil
	}
	return []Violation{{
		KataID: "ZC1588",
		Message: "`nsenter --target 1` joins the host init namespaces — classic container-" +
			"escape primitive. Use `docker exec` / `kubectl exec` for legitimate debugging.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1589",
		Title:    "Warn on `trap 'set -x' ERR/RETURN/EXIT/ZERR` — trace hook leaks env to stderr",
		Severity: SeverityWarning,
		Description: "Installing a trap that enables `set -x` (or `set -o xtrace` / `set -v`) " +
			"causes every subsequent expanded command to hit stderr. Expansions embed " +
			"environment variables — API tokens, passwords, signed URLs — directly into " +
			"the trace. In CI, that stderr lands in build logs and gets shipped to long-term " +
			"log retention. Scope `set -x` to a `set -x ... set +x` block around the suspect " +
			"code, or replace the trap with `trap 'safe_dump' ERR` that prints only non-" +
			"sensitive state.",
		Check: checkZC1589,
	})
}

func checkZC1589(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "trap" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	action := strings.Trim(cmd.Arguments[0].String(), "'\"")
	if !strings.Contains(action, "set -x") &&
		!strings.Contains(action, "set -o xtrace") &&
		!strings.Contains(action, "set -v") &&
		!strings.Contains(action, "set -o verbose") {
		return nil
	}

	sig := cmd.Arguments[1].String()
	switch sig {
	case "ERR", "RETURN", "EXIT", "ZERR":
	default:
		return nil
	}

	return []Violation{{
		KataID: "ZC1589",
		Message: "`trap 'set -x' " + sig + "` enables shell trace from a trap — expansions " +
			"leak env vars (tokens, passwords) to stderr / CI logs. Use a scoped `set -x " +
			"... set +x`, not a trap.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1590",
		Title:    "Error on `sshpass -p SECRET` — password in process list and history",
		Severity: SeverityError,
		Description: "`sshpass -p SECRET` places the password in argv. It leaks into `ps`, " +
			"`/proc/<pid>/cmdline`, shell history, and audit logs for every process on the box " +
			"that can list processes. The `-f FILE` and `-e` (SSHPASS env) variants keep it off " +
			"argv, but key-based auth is the real fix. Generate an SSH key, authorize it on the " +
			"remote, and drop the password tool entirely.",
		Check: checkZC1590,
	})
}

func checkZC1590(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "sshpass" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-p" || (strings.HasPrefix(v, "-p") && len(v) > 2 && v[2] != '=') {
			return []Violation{{
				KataID: "ZC1590",
				Message: "`sshpass -p` places the password in argv — visible in `ps` / " +
					"`/proc/<pid>/cmdline`. Switch to key-based auth, or at least use " +
					"`sshpass -f FILE` / `sshpass -e`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
		if strings.HasPrefix(v, "-p=") {
			return []Violation{{
				KataID: "ZC1590",
				Message: "`sshpass -p` places the password in argv — visible in `ps` / " +
					"`/proc/<pid>/cmdline`. Switch to key-based auth, or at least use " +
					"`sshpass -f FILE` / `sshpass -e`.",
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
		ID:       "ZC1591",
		Title:    "Style: use Zsh `print -l` / `${(F)array}` instead of `printf '%s\\n' \"${array[@]}\"`",
		Severity: SeverityStyle,
		Description: "`printf '%s\\n' \"${array[@]}\"` is the Bash-idiomatic way to print one " +
			"element per line. Zsh has `print -l -r -- \"${array[@]}\"` (one element per line, " +
			"raw, sentinel-safe) and the parameter-expansion flag `${(F)array}` (newline-join, " +
			"fine for `$(...)`). Both are shorter than the printf incantation and avoid format-" +
			"string surprises if the array ever contains a literal `%`.",
		Check: checkZC1591,
		Fix:   fixZC1591,
	})
}

// fixZC1591 rewrites `printf '%s\n' "${array[@]}"` (or `printf '%s'
// "${array[@]}"`) to `print -l -r -- "${array[@]}"`. Single span
// replacement covers the `printf` command name and the format
// argument; subsequent args (the array expansion) pass through
// unchanged. Idempotent — a re-run sees `print`, not `printf`.
func fixZC1591(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "printf" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	formatArg := cmd.Arguments[0]
	formatVal := formatArg.String()
	trimmed := strings.Trim(formatVal, "'\"")
	if trimmed != `%s\n` && trimmed != "%s" {
		return nil
	}
	cmdOff := LineColToByteOffset(source, v.Line, v.Column)
	if cmdOff < 0 || cmdOff+len("printf") > len(source) {
		return nil
	}
	if string(source[cmdOff:cmdOff+len("printf")]) != "printf" {
		return nil
	}
	formatTok := formatArg.TokenLiteralNode()
	formatOff := LineColToByteOffset(source, formatTok.Line, formatTok.Column)
	if formatOff < 0 || formatOff+len(formatVal) > len(source) {
		return nil
	}
	if string(source[formatOff:formatOff+len(formatVal)]) != formatVal {
		return nil
	}
	end := formatOff + len(formatVal)
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - cmdOff,
		Replace: "print -l -r --",
	}}
}

func checkZC1591(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "printf" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	format := strings.Trim(cmd.Arguments[0].String(), "'\"")
	if format != `%s\n` && format != "%s" {
		return nil
	}

	second := cmd.Arguments[1].String()
	if !strings.Contains(second, "[@]") && !strings.Contains(second, "[*]") {
		return nil
	}

	return []Violation{{
		KataID: "ZC1591",
		Message: "`printf '%s\\n' \"${array[@]}\"` — use Zsh `print -l -r -- \"${array[@]}\"` " +
			"or `${(F)array}` for newline-joined output.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1592",
		Title:    "Warn on `faillock --reset` / `pam_tally2 -r` — clears failed-auth counter",
		Severity: SeverityWarning,
		Description: "Both tools zero the PAM counter that triggers account lockout after too " +
			"many failed logins. A script that resets lockouts — even legitimately, to recover " +
			"locked users — also erases evidence of an ongoing brute-force attempt. Intrusion " +
			"detection relies on those counters for alerting. Do not automate resets; if you " +
			"must, log the prior count and page security on every invocation.",
		Check: checkZC1592,
	})
}

func checkZC1592(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "faillock" && ident.Value != "pam_tally2" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--reset" || v == "-r" {
			return []Violation{{
				KataID: "ZC1592",
				Message: "`" + ident.Value + " " + v + "` clears the PAM failed-auth counter — " +
					"masks ongoing brute force. Log the prior count and alert before resetting.",
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
		ID:       "ZC1593",
		Title:    "Error on `blkdiscard` — issues TRIM/DISCARD across the whole device (data loss)",
		Severity: SeverityError,
		Description: "`blkdiscard $DEV` tells the underlying SSD controller to invalidate every " +
			"block in the range. On most modern drives the data is unrecoverable the moment the " +
			"controller acknowledges — even forensic recovery cannot pull it back. Scripts that " +
			"reach this command from any codepath an attacker or typo can trigger destroy the " +
			"drive. Gate it behind interactive confirmation, not shell flow control.",
		Check: checkZC1593,
	})
}

func checkZC1593(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "blkdiscard" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1593",
		Message: "`blkdiscard` issues TRIM/DISCARD across the full device — data is " +
			"unrecoverable once the controller acknowledges. Require operator confirmation " +
			"before running.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1594",
		Title:    "Warn on `docker/podman run --security-opt=systempaths=unconfined` — unhides host kernel knobs",
		Severity: SeverityWarning,
		Description: "`systempaths=unconfined` removes the container runtime's masking of " +
			"`/proc/sys`, `/proc/sysrq-trigger`, `/sys/firmware`, and related kernel surfaces. " +
			"Without the default shield a compromised process inside the container can write " +
			"`/proc/sysrq-trigger` to panic the host, or edit `/proc/sys/kernel/*` to change " +
			"kernel policy on the fly. Keep the default `systempaths=all` (masked) unless you " +
			"have a specific kernel tunable you need, then mount only that path.",
		Check: checkZC1594,
	})
}

func checkZC1594(node ast.Node) []Violation {
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

	runSeen := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if !runSeen {
			if v == "run" || v == "create" {
				runSeen = true
			}
			continue
		}
		if strings.Contains(v, "systempaths=unconfined") {
			return []Violation{{
				KataID: "ZC1594",
				Message: "`" + ident.Value + " run --security-opt=systempaths=unconfined` " +
					"unhides `/proc/sys`, `/proc/sysrq-trigger`, and other kernel knobs. A " +
					"compromise in the container can then panic or re-tune the host.",
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
		ID:       "ZC1595",
		Title:    "Warn on `setfacl -m u:nobody:... / o::rwx` — ACL grants that bypass `chmod` scrutiny",
		Severity: SeverityWarning,
		Description: "Filesystem ACLs live outside the mode bits that `chmod` / `ls -l` / " +
			"`stat -c %a` surface. Granting `u:nobody:rwx` gives the daemon-fallback account " +
			"write access to a file; `o::rwx` / `o::rw` world-writes via ACL even when the mode " +
			"bits still look safe. Review scripts that check `stat -c %a` miss both. Prefer " +
			"`chmod` for world perms, and for specific users name the real account with the " +
			"minimum perm set.",
		Check: checkZC1595,
	})
}

func checkZC1595(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setfacl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "u:nobody:") ||
			strings.HasPrefix(v, "g:nobody:") ||
			strings.HasPrefix(v, "u:nogroup:") ||
			strings.HasPrefix(v, "g:nogroup:") ||
			v == "o::rwx" || v == "o::rw" || v == "o::rwX" {
			return []Violation{{
				KataID: "ZC1595",
				Message: "`setfacl -m " + v + "` grants perms via ACL, bypassing `chmod` / " +
					"`stat -c %a` checks. Prefer chmod for world perms, and for specific " +
					"users name the real account with minimum perms.",
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
		ID:       "ZC1596",
		Title:    "Style: `emulate sh/bash/ksh` without `-L` — flips options for the whole shell",
		Severity: SeverityStyle,
		Description: "`emulate MODE` without the `-L` flag changes Zsh options globally. After " +
			"that line runs the shell is no longer in Zsh mode — `${(F)arr}`, 1-indexed " +
			"arrays, glob qualifiers, and other Zsh-only constructs either error or silently " +
			"behave differently. Wrap emulation in a function and use `emulate -L MODE` to " +
			"scope it to that function. A `.zsh` script that starts with `emulate sh` likely " +
			"belongs in a `.sh` file instead.",
		Check: checkZC1596,
	})
}

func checkZC1596(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "emulate" {
		return nil
	}

	var hasL bool
	var mode string
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "-") {
			if strings.Contains(v, "L") {
				hasL = true
			}
			continue
		}
		if v == "sh" || v == "bash" || v == "ksh" || v == "csh" {
			mode = v
		}
	}
	if mode == "" || hasL {
		return nil
	}

	return []Violation{{
		KataID: "ZC1596",
		Message: "`emulate " + mode + "` without `-L` flips the options for the whole " +
			"shell. Use `emulate -L " + mode + "` inside a function, or rename the script " +
			"to `.sh` if Zsh features are not needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1597",
		Title:    "Warn on `systemd-run -p User=root` — launches arbitrary command with root privileges",
		Severity: SeverityWarning,
		Description: "`systemd-run` submits a transient unit to systemd. With `-p User=root` " +
			"(or `User=0`) the unit runs as root — bypassing the usual `sudo` audit path in " +
			"`/var/log/auth.log`. On hosts where the caller's polkit / dbus rules allow the " +
			"operation, this is effectively privilege escalation by a different name. Prefer " +
			"explicit `sudo` so the invocation is logged, or pre-provision a dedicated systemd " +
			"unit that names the exact command it can run.",
		Check: checkZC1597,
	})
}

func checkZC1597(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "systemd-run" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "User=root" || v == "User=0" {
			return []Violation{{
				KataID: "ZC1597",
				Message: "`systemd-run -p " + v + "` runs arbitrary commands as root via " +
					"systemd — bypasses the `sudo` audit path. Prefer explicit `sudo` or a " +
					"fixed systemd unit.",
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
		ID:       "ZC1598",
		Title:    "Error on `chmod` with world-write bit on a sensitive `/dev/` node",
		Severity: SeverityError,
		Description: "Device nodes under `/dev/` are kernel interfaces. Making one world-writable " +
			"( last digit `2`, `3`, `6`, or `7` ) gives every local user a direct line into the " +
			"kernel — `/dev/kvm` yields VM hypercalls, `/dev/mem` / `/dev/kmem` / `/dev/port` " +
			"read and write physical memory, `/dev/sd*` and `/dev/nvme*` give raw block access, " +
			"`/dev/input/*` sniffs keystrokes. Keep restrictive perms (600 / 660) and use udev " +
			"rules (`GROUP=`, `MODE=`) to grant access declaratively.",
		Check: checkZC1598,
	})
}

func checkZC1598(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chmod" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	mode := cmd.Arguments[0].String()
	if mode == "" {
		return nil
	}
	last := mode[len(mode)-1]
	if last != '2' && last != '3' && last != '6' && last != '7' {
		return nil
	}

	safe := map[string]bool{
		"/dev/null": true, "/dev/zero": true,
		"/dev/random": true, "/dev/urandom": true, "/dev/full": true,
		"/dev/stdin": true, "/dev/stdout": true, "/dev/stderr": true,
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if !strings.HasPrefix(v, "/dev/") {
			continue
		}
		if safe[v] || strings.HasPrefix(v, "/dev/tty") || strings.HasPrefix(v, "/dev/pts/") {
			continue
		}
		return []Violation{{
			KataID: "ZC1598",
			Message: "`chmod " + mode + " " + v + "` makes a sensitive device node world-" +
				"writable — direct kernel access for every local user. Keep restrictive " +
				"perms (600 / 660) and grant access via udev rules.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1599",
		Title:    "Warn on `ldconfig -f PATH` outside `/etc/` — attacker-writable loader cache",
		Severity: SeverityWarning,
		Description: "`ldconfig -f PATH` rebuilds `/etc/ld.so.cache` using PATH instead of the " +
			"system `/etc/ld.so.conf`. If PATH sits in `/tmp`, `/var/tmp`, `$HOME`, or any " +
			"directory an attacker can create, they can inject an `include` line that points " +
			"at their directory of malicious shared objects. After the cache rebuild, every " +
			"subsequent executable on the host loads their library first. Keep the config " +
			"under `/etc/ld.so.conf.d/` with root ownership and run `ldconfig` with no `-f`.",
		Check: checkZC1599,
	})
}

func checkZC1599(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ldconfig" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() != "-f" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		next := cmd.Arguments[i+1].String()
		if strings.HasPrefix(next, "/etc/") || strings.HasPrefix(next, "$") {
			continue
		}
		return []Violation{{
			KataID: "ZC1599",
			Message: "`ldconfig -f " + next + "` uses a config outside `/etc/`. If the " +
				"file is attacker-writable, every binary on the host loads the attacker's " +
				"library. Keep config under `/etc/ld.so.conf.d/`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}
