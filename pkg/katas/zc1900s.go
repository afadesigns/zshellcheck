// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1900",
		Title:    "Warn on `curl --location-trusted` — Authorization/cookies forwarded across redirects",
		Severity: SeverityWarning,
		Description: "`curl --location-trusted` (alias of `curl -L --location-trusted`) tells " +
			"curl to replay the `Authorization` header, cookies, and `-u user:pass` credential " +
			"on every redirect hop, even across hosts. A 302 to an attacker-controlled origin " +
			"(or a compromised CDN edge) then receives the bearer token verbatim. Drop " +
			"`--location-trusted`; if cross-origin auth is truly required, scope a short-lived " +
			"token per destination and verify the final hostname before sending secrets.",
		Check: checkZC1900,
	})
}

var zc1900LocationFlags = map[string]bool{
	"--location-trusted": true,
}

func checkZC1900(node ast.Node) []Violation {
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
		if zc1900LocationFlags[arg.String()] {
			line, col := FlagArgPosition(cmd, zc1900LocationFlags)
			return []Violation{{
				KataID: "ZC1900",
				Message: "`curl --location-trusted` replays `Authorization`, cookies, and " +
					"`-u user:pass` on every redirect — a 302 to attacker-controlled host " +
					"leaks the token. Drop the flag; verify final hostname before sending secrets.",
				Line:   line,
				Column: col,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1901",
		Title:    "Warn on `setopt POSIX_BUILTINS` — flips `command`/special-builtin semantics",
		Severity: SeverityWarning,
		Description: "`setopt POSIX_BUILTINS` switches Zsh to the POSIX rules for special " +
			"builtins: assignments before `export`, `readonly`, `eval`, `.`, `trap`, `set`, " +
			"etc. stay in the caller's scope, and `command builtin` can now resolve shell " +
			"builtins. Mid-script Zsh code written against native semantics — where those " +
			"assignments are local — silently leaks state. Leave the option off; scope any " +
			"POSIX-specific block with `emulate -LR sh` instead of toggling globally.",
		Check: checkZC1901,
	})
}

func checkZC1901(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1901Canonical(arg.String())
		switch v {
		case "POSIXBUILTINS":
			if enabling {
				return zc1901Hit(cmd, "setopt POSIX_BUILTINS")
			}
		case "NOPOSIXBUILTINS":
			if !enabling {
				return zc1901Hit(cmd, "unsetopt NO_POSIX_BUILTINS")
			}
		}
	}
	return nil
}

func zc1901Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1901Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1901",
		Message: "`" + form + "` switches Zsh to POSIX special-builtin rules — " +
			"assignments before `export`/`readonly`/`eval` stop being local, silently " +
			"leaking state. Scope any POSIX block with `emulate -LR sh` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

var zc1902SensitiveTargets = []string{
	"/var/log/",
	"/var/log/audit/",
	"/var/log/wtmp",
	"/var/log/btmp",
	"/var/log/lastlog",
	"/var/log/secure",
	"/var/log/auth.log",
	"/var/log/syslog",
	"/var/log/messages",
	"/.bash_history",
	"/.zsh_history",
	"/.ash_history",
	"/.python_history",
	"/.mysql_history",
	"/.psql_history",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1902",
		Title:    "Error on `ln -s /dev/null <logfile>` — silently discards audit or history writes",
		Severity: SeverityError,
		Description: "A symlink from an audit or shell-history path to `/dev/null` turns every " +
			"subsequent append into a no-op — `/var/log/auth.log`, `wtmp`, `~/.bash_history`, " +
			"`~/.zsh_history` all stop recording. This is the textbook way to cover tracks on " +
			"a compromised host and almost never appears in benign automation. If you really " +
			"need to stop a log, disable the writer (rsyslog rule, `set +o history`) or rotate " +
			"with `logrotate` — never redirect into `/dev/null`.",
		Check: checkZC1902,
	})
}

func checkZC1902(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ln" {
		return nil
	}

	var symbolic bool
	var source, target string
	positional := 0
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "-") {
			if strings.ContainsAny(v, "sS") {
				symbolic = true
			}
			continue
		}
		switch positional {
		case 0:
			source = v
		case 1:
			target = v
		}
		positional++
	}
	if !symbolic || source != "/dev/null" {
		return nil
	}
	if !zc1902IsSensitive(target) {
		return nil
	}

	return []Violation{{
		KataID: "ZC1902",
		Message: "`ln -s /dev/null " + target + "` redirects every write to the " +
			"bit-bucket — audit / history entries vanish silently. If the log must " +
			"stop, disable the writer or rotate with `logrotate`, never symlink to `/dev/null`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func zc1902IsSensitive(target string) bool {
	if target == "" {
		return false
	}
	for _, suffix := range zc1902SensitiveTargets {
		if strings.HasSuffix(target, suffix) || strings.Contains(target, suffix) {
			return true
		}
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1903",
		Title:    "Error on `tee /etc/sudoers*` — appends a rule that bypasses `visudo -c` validation",
		Severity: SeverityError,
		Description: "`tee /etc/sudoers` or `tee -a /etc/sudoers.d/<name>` is a common shortcut " +
			"for adding a sudoers rule, but it skips the syntax check that `visudo -c` would " +
			"perform. A malformed line (missing `ALL`, stray colon, unterminated `Cmnd_Alias`) " +
			"makes sudo refuse every invocation — you lock yourself out of root recovery. " +
			"Write the rule to a temporary file, run `visudo -cf /tmp/rule`, and only then " +
			"`install -m 0440 /tmp/rule /etc/sudoers.d/<name>`.",
		Check: checkZC1903,
	})
}

func checkZC1903(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tee" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "-") {
			continue
		}
		if zc1903IsSudoersTarget(v) {
			return []Violation{{
				KataID: "ZC1903",
				Message: "`tee " + v + "` writes a sudoers rule without `visudo -c` " +
					"validation — a syntax error locks every future `sudo` invocation. " +
					"Write to a temp file, run `visudo -cf`, then `install -m 0440` into place.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1903IsSudoersTarget(v string) bool {
	return v == "/etc/sudoers" || strings.HasPrefix(v, "/etc/sudoers.d/")
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1904",
		Title:    "Warn on `setopt KSH_GLOB` — reinterprets `*(pattern)` and breaks Zsh glob qualifiers",
		Severity: SeverityWarning,
		Description: "`setopt KSH_GLOB` turns `@(a|b)`, `*(x)`, `+(x)`, `?(x)`, `!(x)` into " +
			"Korn-shell extended glob operators. The side effect is that `*(N)`, `*(D)`, " +
			"`*(.)`, and every other Zsh glob qualifier stop working — `*(N)` becomes " +
			"\"zero or more `N` characters\", silently shattering null-glob idioms across the " +
			"script. If you need Korn-style patterns, prefer `setopt EXTENDED_GLOB` and its " +
			"`(^...)` / `(#...)` forms, which coexist with the qualifier syntax. Otherwise " +
			"scope the switch inside a function with `setopt LOCAL_OPTIONS KSH_GLOB`.",
		Check: checkZC1904,
	})
}

func checkZC1904(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1904Canonical(arg.String())
		switch v {
		case "KSHGLOB":
			if enabling {
				return zc1904Hit(cmd, "setopt KSH_GLOB")
			}
		case "NOKSHGLOB":
			if !enabling {
				return zc1904Hit(cmd, "unsetopt NO_KSH_GLOB")
			}
		}
	}
	return nil
}

func zc1904Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1904Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1904",
		Message: "`" + form + "` reinterprets `*(...)` as a ksh-style operator — every " +
			"Zsh glob qualifier (`*(N)`, `*(D)`, `*(.)`) silently stops working. Prefer " +
			"`setopt EXTENDED_GLOB`, or scope inside a function with `LOCAL_OPTIONS`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1905",
		Title:    "Warn on `ssh -g -L …` — local forward bound on all interfaces, not just loopback",
		Severity: SeverityWarning,
		Description: "`ssh -g` flips the default for `-L` (local forward) and `-D` (dynamic SOCKS) " +
			"from `127.0.0.1:port` to `0.0.0.0:port`. Any host on the same LAN/VPN/WiFi " +
			"segment can then use the tunnel without authenticating to the SSH session. " +
			"Drop `-g`, pin the bind explicitly with `-L bind_address:port:target:port`, or " +
			"use a firewall rule — never leave a forwarded port open to the network segment.",
		Check: checkZC1905,
	})
}

func checkZC1905(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh" {
		return nil
	}

	hasG := false
	hasForward := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch {
		case v == "-g":
			hasG = true
		case strings.HasPrefix(v, "-g") && len(v) > 2 && !strings.HasPrefix(v, "-go"):
			// clustered short flags, e.g. `-gNL`
			hasG = true
		case v == "-L", v == "-D":
			hasForward = true
		case strings.HasPrefix(v, "-L") && len(v) > 2:
			hasForward = true
		case strings.HasPrefix(v, "-D") && len(v) > 2:
			hasForward = true
		}
		if strings.HasPrefix(v, "-") && len(v) >= 2 && v[1] != '-' {
			for i := 1; i < len(v); i++ {
				if v[i] == 'g' {
					hasG = true
				}
				if v[i] == 'L' || v[i] == 'D' {
					hasForward = true
				}
			}
		}
	}
	if !hasG || !hasForward {
		return nil
	}

	return []Violation{{
		KataID: "ZC1905",
		Message: "`ssh -g` with `-L`/`-D` binds the forward on `0.0.0.0` — anyone on the " +
			"same LAN segment can ride the tunnel. Drop `-g` or pin `bind_address:port` " +
			"in the forward spec.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1906",
		Title:    "Warn on `setopt POSIX_CD` — changes when `cd` / `pushd` consult `CDPATH`",
		Severity: SeverityWarning,
		Description: "`setopt POSIX_CD` makes `cd`, `chdir`, and `pushd` skip `CDPATH` for any " +
			"argument that starts with `/`, `.`, or `..`. Zsh's default — consulting `CDPATH` " +
			"for anything that does not start with `/` — was exactly what made `cd foo` resolve " +
			"the \"project\" dir via `CDPATH` even when a local `./foo` existed. Flipping " +
			"the option globally makes scripts that relied on the Zsh behaviour silently enter " +
			"different directories. Keep the option off; if POSIX parity is needed, wrap a " +
			"single function with `emulate -LR sh`.",
		Check: checkZC1906,
	})
}

func checkZC1906(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1906Canonical(arg.String())
		switch v {
		case "POSIXCD":
			if enabling {
				return zc1906Hit(cmd, "setopt POSIX_CD")
			}
		case "NOPOSIXCD":
			if !enabling {
				return zc1906Hit(cmd, "unsetopt NO_POSIX_CD")
			}
		}
	}
	return nil
}

func zc1906Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1906Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1906",
		Message: "`" + form + "` changes when `cd`/`pushd` read `CDPATH` — scripts that " +
			"relied on Zsh's default silently enter different directories. Keep it off; " +
			"wrap POSIX-specific code with `emulate -LR sh`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

// zc1907Weak lists the fs.protected_* + fs.suid_dumpable values that roll back
// kernel safeguards against /tmp-race / hardlink-escalation / FIFO-owner
// symlink-attack patterns.
var zc1907Weak = map[string]string{
	"fs.protected_hardlinks=0": "hardlink following",
	"fs.protected_symlinks=0":  "symlink following in world-writable dirs",
	"fs.protected_fifos=0":     "FIFO open in world-writable dirs",
	"fs.protected_regular=0":   "regular-file open in world-writable dirs",
	"fs.suid_dumpable=1":       "SUID core-dump exposure (1 = group-only)",
	"fs.suid_dumpable=2":       "SUID core-dump exposure (2 = root-readable)",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1907",
		Title:    "Warn on `sysctl -w fs.protected_*=0` / `fs.suid_dumpable=2` — disables /tmp-race safeguards",
		Severity: SeverityWarning,
		Description: "Linux ships `fs.protected_symlinks`, `fs.protected_hardlinks`, " +
			"`fs.protected_fifos`, and `fs.protected_regular` enabled to stop classic " +
			"`/tmp`-race escalation (dangling-symlink, hardlink-pivot, FIFO-open-owner). " +
			"Setting any of them to `0`, or raising `fs.suid_dumpable` above `0`, hands " +
			"unprivileged local users back the primitives. Keep the defaults; if a legacy " +
			"tool genuinely needs them off, scope the change inside a namespace rather than " +
			"flipping the host knob.",
		Check: checkZC1907,
	})
}

func checkZC1907(node ast.Node) []Violation {
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

	var writing bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-w" || v == "--write" {
			writing = true
			continue
		}
		if strings.HasPrefix(v, "-") {
			continue
		}
		pair := strings.ReplaceAll(v, " ", "")
		if reason, ok := zc1907Weak[pair]; ok && writing {
			return []Violation{{
				KataID: "ZC1907",
				Message: "`sysctl -w " + pair + "` re-enables " + reason + " — classic " +
					"/tmp-race escalation vector. Keep the default; scope any exception in " +
					"a dedicated namespace.",
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
		ID:       "ZC1908",
		Title:    "Warn on `setopt MAGIC_EQUAL_SUBST` — enables tilde/param expansion on `key=value` args",
		Severity: SeverityWarning,
		Description: "`MAGIC_EQUAL_SUBST` tells Zsh that every unquoted argument of the form " +
			"`identifier=value` gets file expansion on the right-hand side, as if it were a " +
			"parameter assignment. Under the default (option off), `rsync host:dst=~/backup` " +
			"keeps the literal `~` — under the option on, the `~` expands to your home. " +
			"Flipping the option globally makes a whole class of literal CLI arguments silently " +
			"change meaning. Leave the option off; if a specific assignment truly needs " +
			"expansion, wrap it in quotes or use a temporary variable.",
		Check: checkZC1908,
	})
}

func checkZC1908(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1908Canonical(arg.String())
		switch v {
		case "MAGICEQUALSUBST":
			if enabling {
				return zc1908Hit(cmd, "setopt MAGIC_EQUAL_SUBST")
			}
		case "NOMAGICEQUALSUBST":
			if !enabling {
				return zc1908Hit(cmd, "unsetopt NO_MAGIC_EQUAL_SUBST")
			}
		}
	}
	return nil
}

func zc1908Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1908Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1908",
		Message: "`" + form + "` gives every `key=value` argument tilde/parameter " +
			"expansion on the RHS — literal CLI args like `rsync host:dst=~/backup` " +
			"silently change. Keep it off; quote the assignment if expansion is really wanted.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1909",
		Title:    "Warn on `kexec -l` / `-e` — jumps to an alternate kernel, bypasses bootloader and Secure Boot",
		Severity: SeverityWarning,
		Description: "`kexec -l /path/to/vmlinuz …` stages a second kernel image, and `kexec -e` " +
			"(or `kexec -f`) then transfers control to it without going through the firmware, " +
			"GRUB, or shim. On a Secure-Boot system the staged kernel is never verified against " +
			"the enrolled MOK/PK — an attacker who lands a root exec can boot a hostile kernel " +
			"while leaving /boot untouched. Reserve `kexec` for the live-patching / crash-dump " +
			"workflow it was designed for, gate the call behind `sudo` + audit, and prefer " +
			"`systemctl kexec` or a normal reboot when possible.",
		Check: checkZC1909,
	})
}

func checkZC1909(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "load" || ident.Value == "exec" || ident.Value == "unload" {
		// Parser caveat: `kexec --load X` mangles to name=`load`.
		return zc1909Hit(cmd, "kexec --"+ident.Value)
	}
	if ident.Value != "kexec" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-l", "-e", "-f", "-p":
			return zc1909Hit(cmd, "kexec "+v)
		case "--load", "--exec", "--force", "--load-panic":
			return zc1909Hit(cmd, "kexec "+v)
		}
	}
	return nil
}

func zc1909Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1909",
		Message: "`" + form + "` stages or jumps to a kernel without firmware / " +
			"bootloader verification — Secure Boot never checks the signature. Gate behind " +
			"`sudo` + audit and prefer `systemctl kexec` or a real reboot.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1910",
		Title:    "Warn on `setopt GLOB_STAR_SHORT` — makes bare `**` recurse instead of matching literal",
		Severity: SeverityWarning,
		Description: "`GLOB_STAR_SHORT` teaches Zsh to expand bare `**` (not followed by `/`) as if " +
			"it were `**/*` — suddenly `rm **` wipes every file under the current directory " +
			"instead of erroring or matching the two-star literal. Scripts that pass `**` as a " +
			"literal argument to `grep`, `sed`, or a logger call silently turn into deep " +
			"directory recursions. Keep the option off; when you really need recursive globs, " +
			"spell `**/*` explicitly so reviewers can see the intent.",
		Check: checkZC1910,
	})
}

func checkZC1910(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1910Canonical(arg.String())
		switch v {
		case "GLOBSTARSHORT":
			if enabling {
				return zc1910Hit(cmd, "setopt GLOB_STAR_SHORT")
			}
		case "NOGLOBSTARSHORT":
			if !enabling {
				return zc1910Hit(cmd, "unsetopt NO_GLOB_STAR_SHORT")
			}
		}
	}
	return nil
}

func zc1910Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1910Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1910",
		Message: "`" + form + "` turns bare `**` into `**/*` — `rm **` now wipes the tree. " +
			"Keep the option off and spell `**/*` when recursion is actually wanted.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1911",
		Title:    "Warn on `umount -l` / `--lazy` — detach now, leaves open fds pointing at a ghost mount",
		Severity: SeverityWarning,
		Description: "`umount -l` (lazy unmount) detaches the filesystem from the directory tree " +
			"immediately but defers the real cleanup until every open file descriptor on it is " +
			"closed. Any process still holding an fd keeps reading/writing into a mount that " +
			"`mount | grep` no longer lists — cron jobs drop logs into a phantom directory, a " +
			"re-mount of the same path stacks invisibly, and `lsof`/`fuser` often miss the " +
			"stale handles. Find and stop the holder (`lsof`/`fuser`/`systemd-cgls`) first, " +
			"then do a normal `umount`; reserve `-l` for break-glass recovery, not scripts.",
		Check: checkZC1911,
	})
}

var zc1911LazyFlags = map[string]bool{
	"-l":     true,
	"--lazy": true,
}

func checkZC1911(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "umount" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-l" || v == "--lazy" {
			line, col := FlagArgPosition(cmd, zc1911LazyFlags)
			return zc1911Hit(v, line, col)
		}
		if strings.HasPrefix(v, "-") && !strings.HasPrefix(v, "--") {
			// Clustered short flags, e.g. `-fl` / `-lf`.
			for i := 1; i < len(v); i++ {
				if v[i] == 'l' {
					tok := arg.TokenLiteralNode()
					return zc1911Hit("-l", tok.Line, tok.Column)
				}
			}
		}
	}
	return nil
}

func zc1911Hit(flag string, line, col int) []Violation {
	return []Violation{{
		KataID: "ZC1911",
		Message: "`umount " + flag + "` detaches the mount but leaves any open fd pointing at " +
			"a ghost filesystem — writers keep writing, re-mounts stack invisibly. Stop the " +
			"fd holder first (`lsof`/`fuser`), then do a normal `umount`.",
		Line:   line,
		Column: col,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1912",
		Title:    "Warn on `dhclient -r` / `dhclient -x` / `dhcpcd -k` — drops the lease and breaks network",
		Severity: SeverityWarning,
		Description: "`dhclient -r` releases the current DHCP lease (sending a DHCPRELEASE), " +
			"`dhclient -x` terminates the daemon without releasing, and `dhcpcd -k` does the " +
			"equivalent for dhcpcd. On a remote host the very next thing that happens is the " +
			"SSH session drops, and in a VPC any automation waiting for a reply never sees " +
			"one. Stage the release together with a re-acquire (`dhclient -1 $iface` or " +
			"`nmcli device reapply $iface`) or schedule it via `systemd-run --on-active=` " +
			"so the operator is not cut off mid-session.",
		Check: checkZC1912,
	})
}

func checkZC1912(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "dhclient":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-r" || v == "-x" || v == "--release" {
				return zc1912Hit(cmd, "dhclient "+v)
			}
		}
	case "dhcpcd":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-k" || v == "--release" {
				return zc1912Hit(cmd, "dhcpcd "+v)
			}
		}
	}
	return nil
}

func zc1912Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1912",
		Message: "`" + form + "` drops the DHCP lease — SSH session cuts, VPC " +
			"reachability stalls. Pair with a re-acquire (`dhclient -1`/`nmcli device " +
			"reapply`), or schedule via `systemd-run --on-active=`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1913",
		Title:    "Warn on `setopt ALIAS_FUNC_DEF` — re-enables defining functions with aliased names",
		Severity: SeverityWarning,
		Description: "Zsh's default refuses the syntax `ls () { … }` when `ls` is aliased — " +
			"because the alias expands at definition time and the function the author meant " +
			"to write never actually exists. `setopt ALIAS_FUNC_DEF` disables that guardrail: " +
			"the alias is suppressed during definition, and the function silently shadows the " +
			"alias afterwards. The combination is almost always a bug — one alias in a sourced " +
			"rc file quietly replaces the function. Keep the option off and write " +
			"`function \\ls () { … }` (quoted) if you really need to override an aliased name.",
		Check: checkZC1913,
	})
}

func checkZC1913(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1913Canonical(arg.String())
		switch v {
		case "ALIASFUNCDEF":
			if enabling {
				return zc1913Hit(cmd, "setopt ALIAS_FUNC_DEF")
			}
		case "NOALIASFUNCDEF":
			if !enabling {
				return zc1913Hit(cmd, "unsetopt NO_ALIAS_FUNC_DEF")
			}
		}
	}
	return nil
}

func zc1913Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1913Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1913",
		Message: "`" + form + "` lets a function silently shadow an alias — one sourced " +
			"rc file replaces your function with the alias, no error surfaces. Keep it " +
			"off; quote the name if the override is intentional.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1914",
		Title:    "Warn on `curl --doh-url …` / `--dns-servers …` — overrides system resolver per-request",
		Severity: SeverityWarning,
		Description: "`curl --doh-url https://doh.example/dns-query` routes the lookup through a " +
			"caller-specified DNS-over-HTTPS endpoint; `curl --dns-servers 1.1.1.1,8.8.8.8` " +
			"forces classic UDP to the listed servers. Both detour around the host's resolver " +
			"chain — `/etc/hosts`, `systemd-resolved`, `nsswitch`, split-horizon DNS — so the " +
			"request lands at an IP the operator did not vet. In production scripts that is " +
			"usually a stray debug line left in; drop the flag or gate it behind an explicit " +
			"`--doh-insecure` + `--resolve` pinning audit so reviewers can see the intent.",
		Check: checkZC1914,
	})
}

func checkZC1914(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `curl --doh-url https://…` may mangle the name.
	switch ident.Value {
	case "doh-url":
		return zc1914Hit(cmd, "--doh-url")
	case "dns-servers":
		return zc1914Hit(cmd, "--dns-servers")
	case "curl":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			switch {
			case v == "--doh-url", strings.HasPrefix(v, "--doh-url="):
				return zc1914Hit(cmd, "--doh-url")
			case v == "--dns-servers", strings.HasPrefix(v, "--dns-servers="):
				return zc1914Hit(cmd, "--dns-servers")
			}
		}
	}
	return nil
}

func zc1914Hit(cmd *ast.SimpleCommand, flag string) []Violation {
	return []Violation{{
		KataID: "ZC1914",
		Message: "`curl " + flag + "` bypasses the host's resolver chain — `/etc/hosts`, " +
			"`systemd-resolved`, split-horizon DNS — so the request lands at an IP the " +
			"operator did not vet. Drop the flag or pair it with `--resolve` pinning.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1915",
		Title:    "Error on `mdadm --zero-superblock` / `--stop` — drops RAID metadata or live array",
		Severity: SeverityError,
		Description: "`mdadm --zero-superblock $DEV` wipes the MD superblock from a member — the " +
			"array forgets the device exists and a subsequent `--create` with the wrong layout " +
			"permanently scrambles the data. `mdadm --stop $MD` (or `-S`) halts a live array " +
			"from underneath whatever is mounted on it; if root or `/boot` lives there the host " +
			"panics on the next fsync. Run `mdadm --examine` first, snapshot the superblock " +
			"with `mdadm --detail --export`, and keep both calls behind a runbook rather than " +
			"an automated script.",
		Check: checkZC1915,
	})
}

var zc1915Flags = map[string]bool{
	"--zero-superblock": true,
	"-S":                true,
	"--stop":            true,
	"--remove":          true,
}

func checkZC1915(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "mdadm" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "--zero-superblock":
			return zc1915Hit(cmd, "mdadm --zero-superblock")
		case "-S", "--stop":
			return zc1915Hit(cmd, "mdadm "+v)
		case "--remove":
			return zc1915Hit(cmd, "mdadm --remove")
		}
	}
	return nil
}

func zc1915Hit(cmd *ast.SimpleCommand, form string) []Violation {
	line, col := FlagArgPosition(cmd, zc1915Flags)
	return []Violation{{
		KataID: "ZC1915",
		Message: "`" + form + "` drops RAID metadata or halts a live array — mounted root " +
			"or /boot panics the host; a stale superblock scrambles data on next `--create`. " +
			"Snapshot `mdadm --detail --export` first and keep behind a runbook.",
		Line:   line,
		Column: col,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1916",
		Title:    "Warn on `setopt NULL_GLOB` — every unmatched glob silently expands to nothing",
		Severity: SeverityWarning,
		Description: "`setopt NULL_GLOB` removes the Zsh default behaviour of erroring out when a " +
			"pattern matches nothing. Every later glob becomes silently empty instead — " +
			"`cp *.log /dest` when no `.log` files exist turns into `cp /dest` (wrong target), " +
			"`rm *.tmp` into `rm` (argv too short), and `for f in *.json` into a no-op. Reach " +
			"for the per-glob `*(N)` qualifier when you want a single pattern to tolerate a " +
			"zero match, or scope the switch with `setopt LOCAL_OPTIONS NULL_GLOB` inside the " +
			"one function that needs it.",
		Check: checkZC1916,
	})
}

func checkZC1916(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1916Canonical(arg.String())
		switch v {
		case "NULLGLOB":
			if enabling {
				return zc1916Hit(cmd, "setopt NULL_GLOB")
			}
		case "NONULLGLOB":
			if !enabling {
				return zc1916Hit(cmd, "unsetopt NO_NULL_GLOB")
			}
		}
	}
	return nil
}

func zc1916Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1916Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1916",
		Message: "`" + form + "` makes every later unmatched glob silently empty — " +
			"`cp *.log /dest` rewrites to `cp /dest`, `rm *.tmp` becomes argv-too-short. " +
			"Use per-glob `*(N)`, or `setopt LOCAL_OPTIONS NULL_GLOB` in a function.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1917",
		Title:    "Info on `iw dev $IF scan` / `iwlist $IF scan` — active WiFi scan from a script",
		Severity: SeverityInfo,
		Description: "`iw dev wlan0 scan` (and the older `iwlist wlan0 scan`) performs an active " +
			"probe-request sweep across every supported channel. It requires `CAP_NET_ADMIN`, " +
			"briefly interrupts the current association, and announces the host's presence to " +
			"every nearby access point — logs on the other side will show one MAC asking " +
			"about every SSID. Use the cached `iw dev $IF link` / `iwctl station $IF show` " +
			"for passive lookups, and reserve `scan` for diagnostic sessions with console " +
			"approval rather than background scripts.",
		Check: checkZC1917,
	})
}

func checkZC1917(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "iw":
		if zc1917HasSubcmd(cmd, "scan") {
			return zc1917Hit(cmd, "iw dev <if> scan")
		}
	case "iwlist":
		if zc1917HasSubcmd(cmd, "scan") || zc1917HasSubcmd(cmd, "scanning") {
			return zc1917Hit(cmd, "iwlist <if> scan")
		}
	}
	return nil
}

func zc1917HasSubcmd(cmd *ast.SimpleCommand, sub string) bool {
	for _, arg := range cmd.Arguments {
		if arg.String() == sub {
			return true
		}
	}
	return false
}

func zc1917Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1917",
		Message: "`" + form + "` runs an active probe-request sweep — interrupts the " +
			"current association and broadcasts the host to every nearby AP. Use cached " +
			"`iw dev $IF link` for passive queries.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1918",
		Title:    "Warn on `setopt HIST_SUBST_PATTERN` — `!:s/old/new/` silently switches to pattern matching",
		Severity: SeverityWarning,
		Description: "`HIST_SUBST_PATTERN` makes the `:s` and `:&` history modifiers, as well as " +
			"the identically-named parameter-expansion modifier `${foo:s/pat/rep/}`, match on " +
			"patterns rather than literal strings. Text that looked safe as a constant " +
			"(`#` comments, `^` anchors, `?`, `*`) suddenly gets interpreted as glob " +
			"metacharacters, and replacements that always returned the original string now " +
			"edit it in surprising ways. Keep the option off and use `${var//pat/rep}` " +
			"explicitly when you do want glob substitution — that form declares the intent " +
			"at the call site instead of via a shell-wide flag.",
		Check: checkZC1918,
	})
}

func checkZC1918(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1918Canonical(arg.String())
		switch v {
		case "HISTSUBSTPATTERN":
			if enabling {
				return zc1918Hit(cmd, "setopt HIST_SUBST_PATTERN")
			}
		case "NOHISTSUBSTPATTERN":
			if !enabling {
				return zc1918Hit(cmd, "unsetopt NO_HIST_SUBST_PATTERN")
			}
		}
	}
	return nil
}

func zc1918Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1918Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1918",
		Message: "`" + form + "` switches `:s` history/param modifiers to pattern " +
			"matching — literal `*`/`?`/`^` suddenly act as glob metacharacters. Keep it off; " +
			"use `${var//pat/rep}` when you actually want pattern substitution.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1919",
		Title:    "Warn on `ss -K` / `ss --kill` — terminates every socket that matches the filter",
		Severity: SeverityWarning,
		Description: "`ss -K` issues `SOCK_DESTROY` to every socket matching the filter (requires " +
			"`CAP_NET_ADMIN`). With a broad filter — `ss -K state established`, `ss -K dport 22` " +
			"— the command happily terminates the SSH session that is running it, along with " +
			"every backend keep-alive that happens to match. Spell the filter tightly " +
			"(`ss -K dst 10.0.0.5 dport 5432 state close-wait`), test it first without `-K` " +
			"to confirm only the target sockets appear, and wrap the call in a review step " +
			"rather than a scheduled job.",
		Check: checkZC1919,
	})
}

var zc1919KillFlags = map[string]bool{
	"-K":     true,
	"--kill": true,
}

func checkZC1919(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ss" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-K" || v == "--kill" {
			return zc1919Hit(cmd)
		}
		if len(v) >= 2 && v[0] == '-' && v[1] != '-' {
			for i := 1; i < len(v); i++ {
				if v[i] == 'K' {
					return zc1919Hit(cmd)
				}
			}
		}
	}
	return nil
}

func zc1919Hit(cmd *ast.SimpleCommand) []Violation {
	line, col := FlagArgPosition(cmd, zc1919KillFlags)
	return []Violation{{
		KataID: "ZC1919",
		Message: "`ss -K` terminates every socket the filter matches — broad filters " +
			"(`state established`, `dport 22`) kill the running SSH session. Preview " +
			"with the same filter minus `-K`, and pin to a specific dst/port/state tuple.",
		Line:   line,
		Column: col,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1920",
		Title:    "Warn on `setopt VERBOSE` — every executed command is echoed to stderr",
		Severity: SeverityWarning,
		Description: "`setopt VERBOSE` is Zsh's name for the POSIX `set -v` flag: the shell prints " +
			"each command line to stderr immediately after reading it. In a script that " +
			"processes secrets the stderr stream then carries every command that mentions them, " +
			"including `mysql -pSECRET`, `curl -u user:pass`, `export DB_PASS=…`. Unlike " +
			"`set -x` (which already has dedicated detectors) the `VERBOSE` flag is easy to " +
			"leave on by accident because the output looks like normal command echo. Remove " +
			"the call and rely on `printf` / a proper logger; if a debug trace is required, " +
			"scope it in a function with `setopt LOCAL_OPTIONS VERBOSE` then `unsetopt VERBOSE`.",
		Check: checkZC1920,
	})
}

func checkZC1920(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1920Canonical(arg.String())
		switch v {
		case "VERBOSE":
			if enabling {
				return zc1920Hit(cmd, "setopt VERBOSE")
			}
		case "NOVERBOSE":
			if !enabling {
				return zc1920Hit(cmd, "unsetopt NO_VERBOSE")
			}
		}
	}
	return nil
}

func zc1920Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1920Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1920",
		Message: "`" + form + "` echoes every executed command to stderr — any line that " +
			"mentions a password, token, or API key leaks with the trace. Remove and use " +
			"`printf` / a logger, or scope via `setopt LOCAL_OPTIONS VERBOSE` in a helper.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1921",
		Title:    "Warn on `systemctl kill -s KILL` / `--signal=SIGKILL` — skips `ExecStop=`, leaks resources",
		Severity: SeverityWarning,
		Description: "`systemctl kill UNIT -s KILL` (and `--signal=9` / `SIGKILL`) bypasses the " +
			"unit's `ExecStop=` sequence and the `TimeoutStopSec=` budget. Any lockfile, " +
			"socket, or shared-memory segment the service was supposed to unlink survives; the " +
			"next restart often fails with \"address already in use\" or a corrupt journal. " +
			"Default to `systemctl stop UNIT` (or `restart`) and let the stop sequence run. " +
			"Reserve `-s KILL` for a last-resort recovery path with a runbook attached.",
		Check: checkZC1921,
	})
}

func checkZC1921(node ast.Node) []Violation {
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
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "kill" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if strings.HasPrefix(v, "--signal=") {
			if zc1921IsHardSignal(strings.TrimPrefix(v, "--signal=")) {
				return zc1921Hit(cmd, v)
			}
		}
		if v == "-s" && i+2 < len(cmd.Arguments) {
			sig := cmd.Arguments[i+2].String()
			if zc1921IsHardSignal(sig) {
				return zc1921Hit(cmd, "-s "+sig)
			}
		}
	}
	return nil
}

func zc1921IsHardSignal(sig string) bool {
	switch strings.ToUpper(sig) {
	case "KILL", "SIGKILL", "9":
		return true
	}
	return false
}

func zc1921Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1921",
		Message: "`systemctl kill " + form + "` bypasses `ExecStop=` and " +
			"`TimeoutStopSec=` — lockfiles, sockets, and shm segments survive and the next " +
			"restart often fails with \"address already in use\". Use `systemctl stop` or " +
			"`restart` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1922",
		Title:    "Error on `rpm --import http://…` / `rpmkeys --import ftp://…` — plaintext GPG key fetch",
		Severity: SeverityError,
		Description: "`rpm --import` (and `rpmkeys --import`) add the supplied ASCII-armoured " +
			"key to the system RPM trust store. When the source is a plain `http://` / `ftp://` " +
			"URL an on-path attacker swaps the key, and every subsequent package they sign " +
			"installs cleanly. Serve keys over HTTPS from a TLS-authenticated origin, pin the " +
			"key's SHA-256 before import, or stage an offline copy verified out of band " +
			"(`gpg --verify` against a known-good fingerprint).",
		Check: checkZC1922,
	})
}

var zc1922ImportFlags = map[string]bool{"--import": true}

func checkZC1922(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "rpm" && ident.Value != "rpmkeys" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() == "--import" && i+1 < len(cmd.Arguments) {
			url := cmd.Arguments[i+1].String()
			if zc1922IsPlaintextURL(url) {
				return zc1922Hit(cmd, url)
			}
		}
	}
	return nil
}

func zc1922IsPlaintextURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "ftp://")
}

func zc1922Hit(cmd *ast.SimpleCommand, url string) []Violation {
	line, col := FlagArgPosition(cmd, zc1922ImportFlags)
	return []Violation{{
		KataID: "ZC1922",
		Message: "`rpm --import " + url + "` fetches a GPG key over plaintext — on-path " +
			"attackers swap it, every future signed package installs. Use `https://` from " +
			"a pinned origin, or `gpg --verify` against a known fingerprint.",
		Line:   line,
		Column: col,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1923",
		Title:    "Warn on `setopt PRINT_EXIT_VALUE` — every non-zero exit leaks a status line to stderr",
		Severity: SeverityWarning,
		Description: "`PRINT_EXIT_VALUE` makes Zsh emit `zsh: exit N` on stderr after every " +
			"foreground command that returns a non-zero status. In a script the stream is " +
			"typically captured by a supervisor or shipped to a log aggregator, and the " +
			"extra line reveals which tool returned what — including grep / test / curl " +
			"probes that were supposed to stay silent. Worse, tools that parse stderr for " +
			"diagnostics (`git`, `ssh`, `rsync`) now see interleaved shell chatter. Remove " +
			"the `setopt` call; if you actually want a per-command post-mortem, rely on " +
			"`precmd`/`preexec` hooks or an explicit `|| printf …`.",
		Check: checkZC1923,
	})
}

func checkZC1923(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1923Canonical(arg.String())
		switch v {
		case "PRINTEXITVALUE":
			if enabling {
				return zc1923Hit(cmd, "setopt PRINT_EXIT_VALUE")
			}
		case "NOPRINTEXITVALUE":
			if !enabling {
				return zc1923Hit(cmd, "unsetopt NO_PRINT_EXIT_VALUE")
			}
		}
	}
	return nil
}

func zc1923Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1923Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1923",
		Message: "`" + form + "` prints `zsh: exit N` on stderr for every non-zero " +
			"exit — silent grep/test/curl probes suddenly leak status, and tools parsing " +
			"stderr see interleaved shell chatter. Remove; use `|| printf …` per call.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

var zc1924Tools = map[string]bool{
	"virt-cat":       true,
	"virt-copy-out":  true,
	"virt-tar-out":   true,
	"virt-edit":      true,
	"virt-copy-in":   true,
	"virt-tar-in":    true,
	"guestfish":      true,
	"guestmount":     true,
	"virt-customize": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1924",
		Title:    "Warn on `virt-cat` / `virt-copy-out` / `guestfish` / `guestmount` — reads guest disk from host",
		Severity: SeverityWarning,
		Description: "libguestfs tools (`virt-cat`, `virt-copy-out`, `virt-tar-out`, `virt-edit`, " +
			"`virt-customize`, `guestfish`, `guestmount`) open a VM's disk image directly from " +
			"the hypervisor and read or mutate its contents without going through the guest " +
			"OS. That bypasses every in-guest permission, audit, and LUKS keyslot the VM was " +
			"using, and — if the VM is live — risks filesystem corruption because two writers " +
			"are now mounted on the same image. Snapshot the disk first, work on the clone, " +
			"and prefer in-guest `ssh`/`scp`/`ansible` for anything that does not need " +
			"out-of-band recovery.",
		Check: checkZC1924,
	})
}

func checkZC1924(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if !zc1924Tools[ident.Value] {
		return nil
	}
	return []Violation{{
		KataID: "ZC1924",
		Message: "`" + ident.Value + "` reads/writes the VM disk directly from the host — " +
			"bypasses in-guest permissions, audit, and LUKS; a live VM risks corruption " +
			"from double-mount. Snapshot first, work on the clone, prefer in-guest tooling.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1925",
		Title:    "Warn on `unsetopt EQUALS` — disables `=cmd` path expansion and tilde-after-colon",
		Severity: SeverityWarning,
		Description: "Zsh's `EQUALS` option (on by default) is what makes `=python`, `=ls`, and " +
			"`=vim` expand to the absolute path of the command via `$PATH` lookup. It also " +
			"drives the `PATH=~/bin:$PATH` tilde-after-colon expansion. `unsetopt EQUALS` " +
			"turns both off: `=cmd` becomes a literal argument (breaking any idiom that " +
			"relies on the short-path), and `PATH=~/bin:$PATH` stops expanding the tilde " +
			"inside the colon-separated list. Keep the option on; if one function needs " +
			"literal `=` arguments, scope via `setopt LOCAL_OPTIONS; unsetopt EQUALS` inside it.",
		Check: checkZC1925,
	})
}

func checkZC1925(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var disabling bool
	switch ident.Value {
	case "unsetopt":
		disabling = true
	case "setopt":
		disabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1925Canonical(arg.String())
		switch v {
		case "EQUALS":
			if disabling {
				return zc1925Hit(cmd, "unsetopt EQUALS")
			}
		case "NOEQUALS":
			if !disabling {
				return zc1925Hit(cmd, "setopt NO_EQUALS")
			}
		}
	}
	return nil
}

func zc1925Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1925Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1925",
		Message: "`" + form + "` turns off `=cmd` path expansion and tilde-after-colon — " +
			"`=python`/`=ls` become literals and `PATH=~/bin:$PATH` stops tilde-expanding. " +
			"Keep on; scope with `setopt LOCAL_OPTIONS; unsetopt EQUALS` inside a function.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1926",
		Title:    "Warn on `telinit 0/1/6` / `init 0/1/6` — SysV runlevel change halts, reboots, or isolates the host",
		Severity: SeverityWarning,
		Description: "`init 0`, `init 6`, `init 1`, and their `telinit` aliases ask systemd (or " +
			"SysV) to switch runlevel: `0` → `poweroff.target`, `6` → `reboot.target`, " +
			"`1`/`S` → `rescue.target`. From a script the side effect is a remote SSH " +
			"disconnect, an immediate service teardown for every other session on the host, " +
			"and — in the `1`/`S` case — dropping to single-user mode without a console to " +
			"recover. Use `systemctl poweroff`/`reboot`/`rescue` (which are clearer in " +
			"reviews) or schedule via `shutdown -h +N` so the operator has a cancel window.",
		Check: checkZC1926,
	})
}

func checkZC1926(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "init" && ident.Value != "telinit" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	lvl := cmd.Arguments[0].String()
	switch lvl {
	case "0", "1", "6", "S", "s":
		return []Violation{{
			KataID: "ZC1926",
			Message: "`" + ident.Value + " " + lvl + "` changes runlevel — `0` halts, `6` " +
				"reboots, `1`/`S` drops to single-user. Use `systemctl poweroff`/`reboot`/" +
				"`rescue` or `shutdown -h +N` so reviewers can read the intent.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1927",
		Title:    "Error on `xfreerdp /p:SECRET` / `rdesktop -p SECRET` — RDP password visible in argv",
		Severity: SeverityError,
		Description: "`xfreerdp /p:<password>` and `rdesktop -p <password>` (plus the `-p -` " +
			"stdin form when followed by an argv password) put the Windows credential into " +
			"`ps`, `/proc/PID/cmdline`, shell history, and every `ps aux` captured by " +
			"monitoring. Use `xfreerdp /from-stdin` + a piped credential, " +
			"`freerdp-shadow-cli /sec:nla` with a cached credential, or drop the password " +
			"into a protected `.rdp` file passed via `/load-config-file`. Never inline the " +
			"password on the command line.",
		Check: checkZC1927,
	})
}

func checkZC1927(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "xfreerdp", "xfreerdp3", "wlfreerdp":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if strings.HasPrefix(v, "/p:") && len(v) > 3 && v != "/p:-" {
				return zc1927Hit(cmd, ident.Value+" "+v)
			}
		}
	case "rdesktop":
		for i, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-p" && i+1 < len(cmd.Arguments) {
				next := cmd.Arguments[i+1].String()
				if next != "-" {
					return zc1927Hit(cmd, "rdesktop -p "+next)
				}
			}
		}
	}
	return nil
}

func zc1927Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1927",
		Message: "`" + form + "` puts the RDP password in argv — visible in `ps`, " +
			"`/proc`, and shell history. Pipe via `/from-stdin`, read from a protected " +
			"`.rdp` file, or use NLA with a cached credential.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1928",
		Title:    "Warn on `setopt SHARE_HISTORY` — every session writes its history into every sibling session",
		Severity: SeverityWarning,
		Description: "`SHARE_HISTORY` flushes each command to `$HISTFILE` immediately and tells " +
			"all other running zsh sessions to re-read the file. A secret typed in a one-off " +
			"\"private\" terminal — `ssh user@host \"$PASS\"`, `aws sts ... --output text`, " +
			"`git push https://user:token@…` — shows up in every other terminal's `fc -l` list " +
			"seconds later. Prefer `setopt INC_APPEND_HISTORY` (append-only, per-session " +
			"isolation) and `setopt HIST_IGNORE_SPACE` so a leading space keeps the line out " +
			"of history altogether.",
		Check: checkZC1928,
	})
}

func checkZC1928(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1928Canonical(arg.String())
		switch v {
		case "SHAREHISTORY":
			if enabling {
				return zc1928Hit(cmd, "setopt SHARE_HISTORY")
			}
		case "NOSHAREHISTORY":
			if !enabling {
				return zc1928Hit(cmd, "unsetopt NO_SHARE_HISTORY")
			}
		}
	}
	return nil
}

func zc1928Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1928Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1928",
		Message: "`" + form + "` flushes every command into every sibling zsh session — " +
			"secrets typed in one terminal surface in `fc -l` of every other. Prefer " +
			"`setopt INC_APPEND_HISTORY` plus `HIST_IGNORE_SPACE` for safer isolation.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1929",
		Title:    "Warn on `cpio -i` / `--extract` without `--no-absolute-filenames` — archive writes outside CWD",
		Severity: SeverityWarning,
		Description: "`cpio -i` (and `--extract`) is the default copy-in mode: it materialises " +
			"every path stored in the archive verbatim. Paths starting with `/` land where the " +
			"archive told them to, and relative paths containing `..` slip out of the " +
			"extraction directory entirely — so a rogue initramfs or firmware bundle can drop " +
			"files into `/etc/cron.d/`, `/usr/lib/systemd/system/`, or the operator's " +
			"`~/.ssh/authorized_keys`. Always pass `--no-absolute-filenames` and extract into a " +
			"fresh scratch directory reviewed before `mv`.",
		Check: checkZC1929,
	})
}

func checkZC1929(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "cpio" {
		return nil
	}

	extract := false
	safe := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-i", "--extract":
			extract = true
		case "--no-absolute-filenames":
			safe = true
		}
		if len(v) >= 2 && v[0] == '-' && v[1] != '-' {
			for i := 1; i < len(v); i++ {
				if v[i] == 'i' {
					extract = true
				}
			}
		}
	}
	if !extract || safe {
		return nil
	}
	return []Violation{{
		KataID: "ZC1929",
		Message: "`cpio -i` extracts paths verbatim — absolute and `..` entries escape the " +
			"target dir. Pass `--no-absolute-filenames` and stage into a scratch dir before " +
			"`mv` into place.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1930",
		Title:    "Warn on `unsetopt HASH_CMDS` — every command invocation re-walks `$PATH`",
		Severity: SeverityWarning,
		Description: "`HASH_CMDS` (on by default) caches the resolved absolute path of every " +
			"command after its first successful lookup. `unsetopt HASH_CMDS` disables the " +
			"cache, so each invocation re-walks every `$PATH` entry and re-runs `stat()` on " +
			"every candidate. On a slow filesystem (NFS home, encrypted volume, large `$PATH`) " +
			"this adds tens to hundreds of milliseconds per command and can double the runtime " +
			"of a long pipeline. Keep the option on; if you are changing a binary and want the " +
			"cache invalidated, `rehash` (one-shot) or `hash -r` is the scoped fix.",
		Check: checkZC1930,
	})
}

func checkZC1930(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var disabling bool
	switch ident.Value {
	case "unsetopt":
		disabling = true
	case "setopt":
		disabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1930Canonical(arg.String())
		switch v {
		case "HASHCMDS":
			if disabling {
				return zc1930Hit(cmd, "unsetopt HASH_CMDS")
			}
		case "NOHASHCMDS":
			if !disabling {
				return zc1930Hit(cmd, "setopt NO_HASH_CMDS")
			}
		}
	}
	return nil
}

func zc1930Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1930Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1930",
		Message: "`" + form + "` re-walks `$PATH` on every call — tens to hundreds of ms " +
			"per command on slow filesystems. Keep it on; use `rehash` or `hash -r` to " +
			"invalidate the cache after a targeted binary swap.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1931",
		Title:    "Warn on `ip netns delete $NS` / `netns del` — drops the whole network namespace",
		Severity: SeverityWarning,
		Description: "`ip netns delete NAME` / `ip netns del NAME` unmounts the namespace and " +
			"tears down every interface, veth pair, VXLAN, and WireGuard peer living inside. " +
			"Processes still attached lose their network abruptly — container health checks " +
			"fail, BGP sessions drop, and any other process using `ip netns exec NAME …` " +
			"errors out with \"No such file or directory\". Stop the workloads first " +
			"(`systemctl stop`, `pkill -SIGTERM -n $NS`), confirm `ip -n $NS link` is empty, " +
			"then `delete` deliberately — or leave the namespace alone if it is managed by " +
			"Docker/containerd/systemd-nspawn.",
		Check: checkZC1931,
	})
}

func checkZC1931(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ip" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "netns" {
		return nil
	}
	sub := cmd.Arguments[1].String()
	if sub != "delete" && sub != "del" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1931",
		Message: "`ip netns " + sub + "` tears down every interface, veth, tunnel, and " +
			"WireGuard peer inside the namespace. Stop the workloads first and verify " +
			"`ip -n $NS link` is empty before deleting.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1932",
		Title:    "Warn on `unsetopt GLOBAL_EXPORT` — `typeset -x` in a function stops leaking to outer scope",
		Severity: SeverityWarning,
		Description: "`GLOBAL_EXPORT` (on by default) makes `typeset -x VAR=val` inside a function " +
			"not only export `VAR` but also promote it to the outer scope, so callers and " +
			"subsequent functions see the same value. Turning it off changes the meaning of " +
			"every such assignment across the script: exports become function-local and " +
			"vanish the moment the function returns. Scripts that rely on a helper to set up " +
			"`PATH`, `VIRTUAL_ENV`, or `AWS_*` variables suddenly run commands under the old " +
			"environment. Keep the option on; if you want a temporary export, scope it with a " +
			"subshell instead of a shell-wide flip.",
		Check: checkZC1932,
	})
}

func checkZC1932(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var disabling bool
	switch ident.Value {
	case "unsetopt":
		disabling = true
	case "setopt":
		disabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1932Canonical(arg.String())
		switch v {
		case "GLOBALEXPORT":
			if disabling {
				return zc1932Hit(cmd, "unsetopt GLOBAL_EXPORT")
			}
		case "NOGLOBALEXPORT":
			if !disabling {
				return zc1932Hit(cmd, "setopt NO_GLOBAL_EXPORT")
			}
		}
	}
	return nil
}

func zc1932Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1932Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1932",
		Message: "`" + form + "` makes `typeset -x` exports function-local — helper " +
			"functions that set `PATH`/`VIRTUAL_ENV`/`AWS_*` no longer propagate to callers. " +
			"Keep it on; scope temporary exports in a subshell instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1933",
		Title:    "Error on `ipvsadm -C` / `--clear` — wipes every IPVS virtual service, drops load balancer",
		Severity: SeverityError,
		Description: "`ipvsadm -C` (and the long form `--clear`) removes every virtual service, " +
			"real server, and connection entry from the in-kernel IPVS table. Traffic that was " +
			"being load-balanced to a backend farm now falls through to the host's local " +
			"listen sockets (or drops), active keepalived/`ldirectord` states invert, and " +
			"clients see 5xx until an operator replays the config. Save the current table first " +
			"(`ipvsadm --save -n > /run/ipvs.bak`), drain specific services with `ipvsadm -D`, " +
			"and keep `--clear` in break-glass-only runbooks.",
		Check: checkZC1933,
	})
}

var zc1933ClearFlags = map[string]bool{
	"-C":      true,
	"--clear": true,
}

func checkZC1933(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ipvsadm" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-C" || v == "--clear" {
			return zc1933Hit(cmd, "ipvsadm "+v)
		}
	}
	return nil
}

func zc1933Hit(cmd *ast.SimpleCommand, form string) []Violation {
	line, col := FlagArgPosition(cmd, zc1933ClearFlags)
	return []Violation{{
		KataID: "ZC1933",
		Message: "`" + form + "` wipes every IPVS virtual service and real-server binding — " +
			"load balancing stops, clients see 5xx. Save via `ipvsadm --save`, drain " +
			"specific services with `-D`, reserve `--clear` for break-glass.",
		Line:   line,
		Column: col,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1934",
		Title:    "Warn on `setopt AUTO_NAME_DIRS` — any absolute-path parameter becomes a `~name` alias",
		Severity: SeverityWarning,
		Description: "`AUTO_NAME_DIRS` (off by default) auto-registers any parameter whose value is " +
			"an absolute directory path as a named directory — so `foo=/srv/data` immediately " +
			"makes `~foo` resolve to `/srv/data` in later expansions and in `%~` prompt " +
			"sequences. The option silently changes the meaning of `ls ~foo` across the " +
			"script and surfaces directory names in `%~` prompts that the user never opted " +
			"into. Keep the option off and call `hash -d name=/path` explicitly when a named " +
			"directory is actually wanted.",
		Check: checkZC1934,
	})
}

func checkZC1934(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1934Canonical(arg.String())
		switch v {
		case "AUTONAMEDIRS":
			if enabling {
				return zc1934Hit(cmd, "setopt AUTO_NAME_DIRS")
			}
		case "NOAUTONAMEDIRS":
			if !enabling {
				return zc1934Hit(cmd, "unsetopt NO_AUTO_NAME_DIRS")
			}
		}
	}
	return nil
}

func zc1934Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1934Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1934",
		Message: "`" + form + "` auto-registers every absolute-path parameter as a " +
			"named dir — `foo=/srv/data` makes `~foo` expand, `%~` prompts surface names " +
			"the user never picked. Keep off; use `hash -d name=/path`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1935",
		Title:    "Warn on `apt autoremove --purge` / `dnf autoremove` — deletes auto-installed deps and their config",
		Severity: SeverityWarning,
		Description: "`apt autoremove --purge` (and `apt-get autoremove --purge`, `dnf autoremove`, " +
			"`zypper rm --clean-deps`) remove every package the resolver thinks is no longer " +
			"required, plus — with `--purge` — their `/etc` config and data dirs. In CI this " +
			"quietly uproots packages someone else installed manually but never `apt-mark " +
			"manual`-ed, and `--purge` makes the removal irreversible. Run a plain `apt " +
			"autoremove --dry-run` in review, mark the keepers with `apt-mark manual`, and " +
			"drop `--purge` from unattended jobs.",
		Check: checkZC1935,
	})
}

func checkZC1935(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var tool string
	switch ident.Value {
	case "apt", "apt-get":
		tool = ident.Value
	case "dnf", "yum":
		tool = ident.Value
	case "zypper":
		tool = "zypper"
	default:
		return nil
	}

	hasAutoremove := false
	hasPurge := false
	hasCleanDeps := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "autoremove":
			hasAutoremove = true
		case "--purge", "--purge-unused":
			hasPurge = true
		case "--clean-deps":
			hasCleanDeps = true
		}
		if v == "rm" && tool == "zypper" {
			hasAutoremove = true
		}
	}
	if !hasAutoremove {
		return nil
	}

	// apt/apt-get autoremove --purge, or dnf/yum autoremove (always purges
	// configs on RPM distros), or zypper rm --clean-deps.
	if (tool == "apt" || tool == "apt-get") && !hasPurge {
		return nil
	}
	if tool == "zypper" && !hasCleanDeps {
		return nil
	}

	return []Violation{{
		KataID: "ZC1935",
		Message: "`" + tool + " autoremove` strips packages the resolver thinks are " +
			"unused plus their configs — uproots packages installed manually but never " +
			"`apt-mark manual`-ed. Dry-run first, mark keepers, drop `--purge` in CI.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1936",
		Title:    "Warn on `setopt POSIX_ALIASES` — aliases on reserved words (`if`, `for`, …) stop expanding",
		Severity: SeverityWarning,
		Description: "Zsh by default lets `alias if='…'`, `alias function='…'`, etc. expand when " +
			"the reserved word appears in command position — the feature that makes oh-my-zsh " +
			"plugins able to hook `if` into their `preexec` chain. `setopt POSIX_ALIASES` " +
			"narrows alias expansion to plain identifiers, so any library that aliased a " +
			"reserved word silently stops being picked up. Keep the option off for " +
			"interactive Zsh; if you need POSIX parity for a specific block, wrap it with " +
			"`emulate -LR sh` instead of flipping the flag script-wide.",
		Check: checkZC1936,
	})
}

func checkZC1936(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1936Canonical(arg.String())
		switch v {
		case "POSIXALIASES":
			if enabling {
				return zc1936Hit(cmd, "setopt POSIX_ALIASES")
			}
		case "NOPOSIXALIASES":
			if !enabling {
				return zc1936Hit(cmd, "unsetopt NO_POSIX_ALIASES")
			}
		}
	}
	return nil
}

func zc1936Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1936Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1936",
		Message: "`" + form + "` narrows alias expansion to plain identifiers — aliases " +
			"on `if`/`for`/`function` silently stop firing and any library that hooked " +
			"them breaks. Scope with `emulate -LR sh` instead of flipping globally.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1937",
		Title:    "Warn on `tmux kill-server` / `tmux kill-session` — tears down every detached process inside",
		Severity: SeverityWarning,
		Description: "`tmux kill-server` terminates the whole tmux daemon, `tmux kill-session -t NAME` " +
			"drops one named session, and `screen -X quit` does the screen equivalent. Anything " +
			"the operator parked inside — a long-running build, a `tail -F` on production " +
			"logs, a held `sudo` token, a port-forward — dies with the session, and the " +
			"detached processes get `SIGHUP`'d with no cleanup. Use `tmux kill-window -t …` for " +
			"surgical removal, send `SIGTERM` to the specific backend, or rely on `systemd-run " +
			"--scope` for workloads that should survive terminal churn.",
		Check: checkZC1937,
	})
}

func checkZC1937(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "tmux" && len(cmd.Arguments) >= 1 {
		sub := cmd.Arguments[0].String()
		if sub == "kill-server" || sub == "kill-session" {
			return zc1937Hit(cmd, "tmux "+sub)
		}
	}
	if ident.Value == "screen" {
		for i, arg := range cmd.Arguments {
			if arg.String() == "-X" && i+1 < len(cmd.Arguments) &&
				cmd.Arguments[i+1].String() == "quit" {
				return zc1937Hit(cmd, "screen -X quit")
			}
		}
	}
	return nil
}

func zc1937Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1937",
		Message: "`" + form + "` tears down every detached process inside the session — " +
			"builds, log tails, port-forwards get `SIGHUP`'d with no cleanup. Use " +
			"`kill-window` for surgical removal or `systemd-run --scope` for workloads.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1938",
		Title:    "Warn on `setopt POSIX_JOBS` — flips job-control semantics and `%n` scope",
		Severity: SeverityWarning,
		Description: "`POSIX_JOBS` makes Zsh's job-control spec follow POSIX: `%1` / `%n` refer " +
			"only to jobs of the current shell (forked subshells get their own job table), " +
			"`fg`/`bg` no longer accept a job ID from an outer shell, and `disown` on a " +
			"subshell's job is a no-op. Scripts that launched a background job in the parent " +
			"and then `wait %1`-ed from a `( subshell )` suddenly fail with \"no such job\". " +
			"Leave the option off in Zsh; if POSIX job semantics are required, scope them via " +
			"`emulate -LR sh` inside the single function that needs them.",
		Check: checkZC1938,
	})
}

func checkZC1938(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1938Canonical(arg.String())
		switch v {
		case "POSIXJOBS":
			if enabling {
				return zc1938Hit(cmd, "setopt POSIX_JOBS")
			}
		case "NOPOSIXJOBS":
			if !enabling {
				return zc1938Hit(cmd, "unsetopt NO_POSIX_JOBS")
			}
		}
	}
	return nil
}

func zc1938Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1938Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1938",
		Message: "`" + form + "` scopes `%n` / `fg` / `bg` / `disown` per subshell — parent " +
			"jobs become invisible inside `(…)`. Leave off; scope POSIX job semantics with " +
			"`emulate -LR sh` inside a function.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1939",
		Title:    "Error on `reboot -f` / `halt -f` / `poweroff -f` — skips shutdown sequence, no graceful service stop",
		Severity: SeverityError,
		Description: "`reboot -f`, `halt -f`, and `poweroff -f` short-circuit the systemd " +
			"shutdown graph — no `ExecStop=`, no `DefaultDependencies=`, no filesystem sync, " +
			"no Before/After ordering. The kernel's `reboot(2)` fires immediately and every " +
			"dirty buffer that was not yet flushed is lost. Journal writes stop mid-line, " +
			"databases on the host replay from the last checkpoint, and anything that needed a " +
			"clean unmount (LUKS, NFS, cephfs) logs a dirty state. Use plain `systemctl " +
			"reboot` / `shutdown -r +N`, and reserve `-f` for recovery when the normal path is " +
			"already wedged.",
		Check: checkZC1939,
	})
}

var zc1939ForceFlags = map[string]bool{
	"-f":      true,
	"--force": true,
}

func checkZC1939(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "reboot" && ident.Value != "halt" && ident.Value != "poweroff" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-f" || v == "--force" {
			return zc1939Hit(cmd, ident.Value+" "+v)
		}
		if len(v) >= 2 && v[0] == '-' && v[1] != '-' {
			for i := 1; i < len(v); i++ {
				if v[i] == 'f' {
					return zc1939Hit(cmd, ident.Value+" -f")
				}
			}
		}
	}
	return nil
}

func zc1939Hit(cmd *ast.SimpleCommand, form string) []Violation {
	line, col := FlagArgPosition(cmd, zc1939ForceFlags)
	return []Violation{{
		KataID: "ZC1939",
		Message: "`" + form + "` fires `reboot(2)` immediately — no `ExecStop=`, no " +
			"filesystem sync, no clean unmount. Databases replay from last checkpoint. " +
			"Use `systemctl reboot` / `shutdown -r +N`; reserve `-f` for wedged recovery.",
		Line:   line,
		Column: col,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1940",
		Title:    "Warn on `setopt POSIX_ARGZERO` — `$0` no longer changes to the function name inside functions",
		Severity: SeverityWarning,
		Description: "Zsh's default behaviour (option off) assigns `$0` to the name of the " +
			"currently-running function, so a helper like `log() { printf '%s\\n' \"$0: $*\"; }` " +
			"prints `log: …`. `setopt POSIX_ARGZERO` keeps `$0` pointing at the outer script " +
			"name (or the interpreter when sourced) — the logger instead prints the script " +
			"path for every message and call-site context is lost. Every `case $0` dispatch " +
			"inside an auto-loaded function also stops working. Leave the option off; if you " +
			"need POSIX `$0`, scope it in a function with `emulate -LR sh`.",
		Check: checkZC1940,
	})
}

func checkZC1940(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1940Canonical(arg.String())
		switch v {
		case "POSIXARGZERO":
			if enabling {
				return zc1940Hit(cmd, "setopt POSIX_ARGZERO")
			}
		case "NOPOSIXARGZERO":
			if !enabling {
				return zc1940Hit(cmd, "unsetopt NO_POSIX_ARGZERO")
			}
		}
	}
	return nil
}

func zc1940Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1940Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1940",
		Message: "`" + form + "` freezes `$0` to the outer script name — loggers and " +
			"`case $0` dispatch inside functions lose call-site context. Scope with " +
			"`emulate -LR sh` instead of flipping globally.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1941",
		Title:    "Error on `restic init --insecure-no-password` — creates an unencrypted backup repository",
		Severity: SeverityError,
		Description: "`restic init --insecure-no-password` creates a repo whose data chunks are " +
			"reachable without a key. Every later `backup` and `restore` round-trips " +
			"plaintext blocks to the storage backend, so any operator with read access to the " +
			"bucket / NFS share / SFTP directory can assemble the backed-up filesystem — " +
			"including shell history, SSH keys, and database dumps. Pass a real passphrase via " +
			"`--password-file` (mode `0400`, readable only by the backup user) or " +
			"`--password-command`, and never use the `--insecure-*` family outside a local " +
			"test repo.",
		Check: checkZC1941,
	})
}

func checkZC1941(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `restic --insecure-no-password …` mangles the command
	// name to `insecure-no-password`.
	if ident.Value == "insecure-no-password" {
		return zc1941Hit(cmd)
	}
	if ident.Value != "restic" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if arg.String() == "--insecure-no-password" {
			return zc1941Hit(cmd)
		}
	}
	return nil
}

func zc1941Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1941",
		Message: "`restic --insecure-no-password` creates an unencrypted repo — every " +
			"operator with read access to the backend can reassemble the backed-up " +
			"filesystem. Use `--password-file` / `--password-command` with a real passphrase.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1942",
		Title:    "Warn on `setopt CLOBBER_EMPTY` — `>file` still overwrites zero-length files under `NO_CLOBBER`",
		Severity: SeverityWarning,
		Description: "`setopt CLOBBER_EMPTY` relaxes `NO_CLOBBER`: a bare `>file` redirect still " +
			"succeeds when the target is zero bytes. Scripts that rely on `setopt NO_CLOBBER` " +
			"as a guard against accidental overwrite lose their safety net for every " +
			"freshly-`touch`ed lock file, sentinel, or `install -D`-created placeholder — the " +
			"next stray `>sentinel` quietly overwrites it. Keep the option off; use `>|file` " +
			"explicitly when you do want to bypass the `NO_CLOBBER` guard for a specific write.",
		Check: checkZC1942,
	})
}

func checkZC1942(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1942Canonical(arg.String())
		switch v {
		case "CLOBBEREMPTY":
			if enabling {
				return zc1942Hit(cmd, "setopt CLOBBER_EMPTY")
			}
		case "NOCLOBBEREMPTY":
			if !enabling {
				return zc1942Hit(cmd, "unsetopt NO_CLOBBER_EMPTY")
			}
		}
	}
	return nil
}

func zc1942Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1942Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1942",
		Message: "`" + form + "` lets `>file` overwrite zero-length files even under " +
			"`NO_CLOBBER` — `touch`ed lock / sentinel files lose their safety net. Keep " +
			"off; use explicit `>|file` to bypass `NO_CLOBBER` for a specific write.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1943",
		Title:    "Warn on `systemd-nspawn -b` / `--boot` — runs a full init inside a possibly untrusted rootfs",
		Severity: SeverityWarning,
		Description: "`systemd-nspawn -b -D $ROOT` (and `--boot -D $ROOT`) launches the rootfs's " +
			"`/sbin/init` inside a minimally-isolated namespace — by default the container " +
			"inherits `CAP_AUDIT_CONTROL`, `CAP_NET_ADMIN`, and read-write access to the " +
			"host's `/dev` nodes that match the container's cgroup. If `$ROOT` is an " +
			"operator-supplied tarball, any init script it ships runs first and can probe the " +
			"host. Use `-U` for user-namespace isolation, drop capabilities with " +
			"`--capability=`, pair with `--private-network`, and prefer `machinectl start` on " +
			"a reviewed image instead of ad-hoc boots.",
		Check: checkZC1943,
	})
}

var zc1943BootFlags = map[string]bool{
	"-b":     true,
	"--boot": true,
}

func checkZC1943(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "systemd-nspawn" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1943BootFlags[v] {
			line, col := FlagArgPosition(cmd, zc1943BootFlags)
			return []Violation{{
				KataID: "ZC1943",
				Message: "`systemd-nspawn " + v + "` runs the rootfs's `/sbin/init` with minimal isolation — " +
					"init scripts execute first and can probe the host. Use `-U`, drop caps with " +
					"`--capability=`, pair with `--private-network`, prefer `machinectl start`.",
				Line:   line,
				Column: col,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1944",
		Title:    "Warn on `setopt IGNORE_EOF` — Ctrl-D no longer exits the shell, masking runaway pipelines",
		Severity: SeverityWarning,
		Description: "`IGNORE_EOF` tells the interactive shell to treat an end-of-file on stdin as " +
			"if it were nothing, so `Ctrl-D` stops terminating a login. In an unattended `zsh " +
			"-i -c` launch, or a sourced rc, this keeps a subshell alive that was supposed to " +
			"wind down when the controlling terminal went away — sudo sessions, SSH tunnels, " +
			"port-forwards, and build supervisors then linger long after the parent left. Keep " +
			"the option off; if a stale-tty guard is truly wanted, set `TMOUT=NN` for a timed " +
			"exit instead.",
		Check: checkZC1944,
	})
}

func checkZC1944(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1944Canonical(arg.String())
		switch v {
		case "IGNOREEOF":
			if enabling {
				return zc1944Hit(cmd, "setopt IGNORE_EOF")
			}
		case "NOIGNOREEOF":
			if !enabling {
				return zc1944Hit(cmd, "unsetopt NO_IGNORE_EOF")
			}
		}
	}
	return nil
}

func zc1944Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1944Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1944",
		Message: "`" + form + "` makes `Ctrl-D` stop terminating the shell — subshells, " +
			"sudo holds, SSH tunnels linger after the parent left. Keep off; use " +
			"`TMOUT=NN` for a timed stale-tty exit if needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1945",
		Title:    "Warn on `bpftrace -e` / `bpftool prog load` — loads in-kernel eBPF from a script",
		Severity: SeverityWarning,
		Description: "`bpftrace -e '…'` compiles an inline script into an eBPF program and attaches " +
			"to kprobes, tracepoints, or uprobes; `bpftool prog load FILE pinned /sys/fs/bpf/…` " +
			"installs a pre-built program. Both require `CAP_BPF`/`CAP_SYS_ADMIN` and can read " +
			"arbitrary kernel/userland memory — every command a sibling process runs, every " +
			"syscall argument, every TCP payload. Pin the loaded program to a directory the " +
			"operator owns, gate invocation behind a runbook, and prefer a short-lived " +
			"`bpftrace -c CMD` window over long-running traces left on the host.",
		Check: checkZC1945,
	})
}

func checkZC1945(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "bpftrace":
		for _, arg := range cmd.Arguments {
			if arg.String() == "-e" {
				return zc1945Hit(cmd, "bpftrace -e")
			}
		}
	case "bpftool":
		if len(cmd.Arguments) >= 2 &&
			cmd.Arguments[0].String() == "prog" &&
			cmd.Arguments[1].String() == "load" {
			return zc1945Hit(cmd, "bpftool prog load")
		}
	}
	return nil
}

func zc1945Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1945",
		Message: "`" + form + "` loads an in-kernel eBPF program that can read arbitrary " +
			"kernel/userland memory — every syscall arg, every TCP payload. Gate behind a " +
			"runbook and prefer a short-lived `bpftrace -c CMD` over a pinned trace.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1946",
		Title:    "Warn on `unsetopt HUP` — background jobs keep running after shell exit",
		Severity: SeverityWarning,
		Description: "Zsh's `HUP` option (on by default) sends `SIGHUP` to each running child job " +
			"when the shell exits, letting them wind down cleanly. `unsetopt HUP` / " +
			"`setopt NO_HUP` disables that, so long pipelines, `sleep` loops, and " +
			"user-spawned daemons live on — `ps aux` accumulates orphaned workers across " +
			"logouts and resource consumption creeps up. If a specific job really needs to " +
			"outlive the shell, use `disown` or `systemd-run --scope` on that one invocation; " +
			"leave `HUP` on globally.",
		Check: checkZC1946,
	})
}

func checkZC1946(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var disabling bool
	switch ident.Value {
	case "unsetopt":
		disabling = true
	case "setopt":
		disabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1946Canonical(arg.String())
		switch v {
		case "HUP":
			if disabling {
				return zc1946Hit(cmd, "unsetopt HUP")
			}
		case "NOHUP":
			if !disabling {
				return zc1946Hit(cmd, "setopt NO_HUP")
			}
		}
	}
	return nil
}

func zc1946Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1946Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1946",
		Message: "`" + form + "` stops the shell from `SIGHUP`-ing background jobs on " +
			"exit — long pipelines and spawned daemons outlive the session, orphans " +
			"accumulate. Use `disown` or `systemd-run --scope` on specific commands instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1947",
		Title:    "Error on `ip xfrm state flush` / `ip xfrm policy flush` — tears down every IPsec SA and policy",
		Severity: SeverityError,
		Description: "`ip xfrm state flush` removes every IPsec Security Association; " +
			"`ip xfrm policy flush` removes every policy that would have driven them. " +
			"Strongswan, libreswan, FRR, and WireGuard-over-xfrm all lose their tunnels " +
			"instantly — site-to-site VPNs drop, kernel packet paths stop encrypting, and " +
			"peers renegotiate from scratch (with traffic leaking in plaintext during the gap " +
			"on misconfigured hosts). Use `ip xfrm state deleteall src $A dst $B` to scope " +
			"the change to a single tunnel, and pair flushes with a maintenance window.",
		Check: checkZC1947,
	})
}

func checkZC1947(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ip" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "xfrm" {
		return nil
	}
	tbl := cmd.Arguments[1].String()
	if tbl != "state" && tbl != "policy" {
		return nil
	}
	if cmd.Arguments[2].String() != "flush" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1947",
		Message: "`ip xfrm " + tbl + " flush` tears down every IPsec SA/policy — " +
			"VPN tunnels drop, kernel stops encrypting, plaintext may leak during renegotiation. " +
			"Scope via `ip xfrm " + tbl + " deleteall src $A dst $B`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

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

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1949",
		Title:    "Error on `rmmod -f` / `rmmod --force` — bypasses refcount, can panic the kernel",
		Severity: SeverityError,
		Description: "`rmmod -f` asks the kernel to tear down a module even if its reference count " +
			"is non-zero. Any live `open(\"/dev/…\")`, mounted filesystem, or in-flight network " +
			"device driven by that module becomes a dangling pointer — the kernel oopses or " +
			"outright panics as soon as the next callback fires. The feature is compiled out " +
			"on most distros (`CONFIG_MODULE_FORCE_UNLOAD=n`), but when present it is strictly " +
			"a break-glass recovery tool. Stop the holders first (`lsof /dev/FOO`, `umount`, " +
			"`ip link set dev … down`), then use plain `rmmod` or `modprobe -r`.",
		Check: checkZC1949,
	})
}

var zc1949ForceFlags = map[string]bool{
	"-f":      true,
	"--force": true,
}

func checkZC1949(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "rmmod" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1949ForceFlags[v] {
			line, col := FlagArgPosition(cmd, zc1949ForceFlags)
			return []Violation{{
				KataID: "ZC1949",
				Message: "`rmmod " + v + "` tears down a module even when its refcount is non-zero — " +
					"in-use drivers dangle, kernel oopses on the next callback. Stop holders first " +
					"(`lsof`/`umount`/`ip link down`), then `rmmod` without `-f`.",
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
		ID:       "ZC1950",
		Title:    "Error on `tune2fs -O ^has_journal` / `-m 0` — removes journal or root reserve",
		Severity: SeverityError,
		Description: "`tune2fs -O ^has_journal $DEV` strips the ext3/4 journal from the " +
			"filesystem. Crash recovery drops from \"replay the journal\" to \"scan the whole " +
			"block device with `fsck -y`\", which frequently truncates partially-written files. " +
			"`tune2fs -m 0 $DEV` takes the reserved-for-root space down to zero; when the " +
			"filesystem fills up there is no headroom for `journald`, `apt`, or even a root " +
			"shell to clean up — recovery needs rescue media. Keep the journal on and leave " +
			"`-m` at the distro default (5% is overkill on large disks, but `-m 1` is still " +
			"safe).",
		Check: checkZC1950,
	})
}

func checkZC1950(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tune2fs" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-O" && i+1 < len(cmd.Arguments) {
			spec := cmd.Arguments[i+1].String()
			if zc1950StripsJournal(spec) {
				return zc1950Hit(cmd, "-O "+spec, "strips the journal — crash recovery needs a full `fsck -y` and may truncate files")
			}
		}
		if v == "-m" && i+1 < len(cmd.Arguments) {
			if cmd.Arguments[i+1].String() == "0" {
				return zc1950Hit(cmd, "-m 0", "zeroes the root reserve — a full fs leaves no headroom for `journald`/`apt`/root shells")
			}
		}
	}
	return nil
}

func zc1950StripsJournal(spec string) bool {
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "^has_journal" {
			return true
		}
	}
	return false
}

func zc1950Hit(cmd *ast.SimpleCommand, form, why string) []Violation {
	return []Violation{{
		KataID:  "ZC1950",
		Message: "`tune2fs " + form + "` " + why + ". Keep the default.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1951",
		Title:    "Error on `ceph osd pool delete … --yes-i-really-really-mean-it` — automates Ceph's double-safety phrase",
		Severity: SeverityError,
		Description: "Ceph intentionally requires both the pool name twice and the flag " +
			"`--yes-i-really-really-mean-it` before it will delete a pool, so a typo during a " +
			"live operation cannot drop production data. Baking the phrase into a script " +
			"defeats the friction — a rebase of the wrong variable, a typo in the pool name, " +
			"or a stale `for pool in $(…)` loop then silently deletes real pools. Remove the " +
			"flag from scripts. Do the deletion interactively, or wrap it in a runbook that " +
			"spells out the pool name in the commit message the operator acknowledges.",
		Check: checkZC1951,
	})
}

func checkZC1951(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `ceph osd pool delete … --yes-i-really-really-mean-it`
	// mangles the command name to `yes-i-really-really-mean-it`.
	if ident.Value == "yes-i-really-really-mean-it" ||
		ident.Value == "yes-i-really-mean-it" {
		return zc1951Hit(cmd)
	}
	if ident.Value != "ceph" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--yes-i-really-really-mean-it" || v == "--yes-i-really-mean-it" {
			return zc1951Hit(cmd)
		}
	}
	return nil
}

func zc1951Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1951",
		Message: "`ceph … --yes-i-really-really-mean-it` automates the double-safety " +
			"phrase — a typo or stale loop silently deletes production pools. Run " +
			"deletions interactively, or spell the pool name in a runbook commit.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1952",
		Title:    "Error on `zfs set sync=disabled` — `fsync()` becomes a no-op, crash loses unflushed writes",
		Severity: SeverityError,
		Description: "`zfs set sync=disabled POOL/DATASET` turns `fsync()`, `O_SYNC`, and `O_DSYNC` " +
			"into no-ops on that dataset. PostgreSQL, MariaDB, etcd, and every application that " +
			"relies on fsync for durability will report success for writes that are still in the " +
			"ARC, so a panic or power cut loses minutes of committed transactions. The flag is " +
			"a benchmarking knob, not a production setting. Leave sync at `standard` and, if " +
			"latency is the concern, add a `log` vdev (SLOG) or tune " +
			"`zfs_txg_timeout` instead.",
		Check: checkZC1952,
	})
}

func checkZC1952(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "zfs" {
		return nil
	}
	if len(cmd.Arguments) < 2 || cmd.Arguments[0].String() != "set" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if strings.HasPrefix(v, "sync=") {
			val := strings.TrimPrefix(v, "sync=")
			if val == "disabled" {
				return []Violation{{
					KataID: "ZC1952",
					Message: "`zfs set sync=disabled` turns `fsync()` into a no-op — DBs " +
						"(PostgreSQL/MariaDB/etcd) lose committed transactions on crash. Leave " +
						"sync at `standard`; use a SLOG vdev if latency is the concern.",
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
		ID:       "ZC1953",
		Title:    "Warn on `mount --make-shared` / `--make-rshared` — flips propagation, container-escape vector",
		Severity: SeverityWarning,
		Description: "`mount --make-shared /path` (and the recursive `--make-rshared`) turns the " +
			"mount point into a peer in a shared-subtree group. Any later bind-mount that " +
			"lands inside it propagates to every other peer, including containers and other " +
			"namespaces. Combined with `CAP_SYS_ADMIN` inside a pod, that is one of the " +
			"classic container-escape stepping stones — a hostile workload can mount into the " +
			"host's `/` via the propagation group. Use `--make-private` on sensitive paths and " +
			"mount containers with `--mount-propagation=private` / `slave` unless the app " +
			"genuinely requires `shared`.",
		Check: checkZC1953,
	})
}

var zc1953ShareFlags = map[string]bool{
	"--make-shared":  true,
	"--make-rshared": true,
}

func checkZC1953(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "mount" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1953ShareFlags[v] {
			line, col := FlagArgPosition(cmd, zc1953ShareFlags)
			return []Violation{{
				KataID: "ZC1953",
				Message: "`mount " + v + "` puts the mount in a shared-subtree group — later " +
					"bind-mounts propagate to every peer, including containers. Classic escape " +
					"stepping stone. Use `--make-private` on sensitive paths.",
				Line:   line,
				Column: col,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1954",
		Title:    "Warn on `setfattr -n security.capability|security.selinux|security.ima` — bypasses `setcap`/`chcon`",
		Severity: SeverityWarning,
		Description: "`setfattr -n security.capability -v …` writes the raw file-capability xattr " +
			"that the kernel consults when a binary `execve()`s, bypassing the `setcap` " +
			"wrapper's validation and audit trail. Similarly, `security.selinux` replaces the " +
			"SELinux label without going through `chcon` / `semanage`, and `security.ima` " +
			"overwrites the IMA hash that integrity-measurement trusts. These attributes are " +
			"the raw kernel knobs behind purpose-built tools; script usage is almost always " +
			"wrong. Use `setcap`, `chcon`/`semanage fcontext`, and `evmctl` instead.",
		Check: checkZC1954,
	})
}

func checkZC1954(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setfattr" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-n" && i+1 < len(cmd.Arguments) {
			name := cmd.Arguments[i+1].String()
			if zc1954SecurityAttr(name) {
				return zc1954Hit(cmd, name)
			}
		}
		if strings.HasPrefix(v, "-n") && len(v) > 2 {
			if zc1954SecurityAttr(v[2:]) {
				return zc1954Hit(cmd, v[2:])
			}
		}
		if strings.HasPrefix(v, "--name=") {
			if zc1954SecurityAttr(strings.TrimPrefix(v, "--name=")) {
				return zc1954Hit(cmd, strings.TrimPrefix(v, "--name="))
			}
		}
	}
	return nil
}

func zc1954SecurityAttr(name string) bool {
	switch {
	case name == "security.capability",
		name == "security.selinux",
		name == "security.ima",
		name == "security.evm":
		return true
	case strings.HasPrefix(name, "security.apparmor"):
		return true
	}
	return false
}

func zc1954Hit(cmd *ast.SimpleCommand, attr string) []Violation {
	return []Violation{{
		KataID: "ZC1954",
		Message: "`setfattr -n " + attr + "` writes the raw kernel xattr — bypasses " +
			"`setcap`/`chcon`/`evmctl` validation and audit. Use the purpose-built tool.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1955",
		Title:    "Warn on `rfkill block all` / `block wifi|bluetooth|wwan` — disables every radio, cuts wireless",
		Severity: SeverityWarning,
		Description: "`rfkill block all` toggles the soft-kill switch on every radio the kernel " +
			"registered — WiFi, Bluetooth, WWAN, NFC, GPS, UWB — so the host drops off the " +
			"network in one call. A follow-up `rfkill unblock all` takes seconds to a minute " +
			"on some drivers and requires the operator to be physically present or have a " +
			"cellular fallback. Scope the block to a specific type (e.g. `rfkill block " +
			"bluetooth`) and schedule via `at now + 5 minutes ... rfkill unblock all` so the " +
			"host recovers on its own.",
		Check: checkZC1955,
	})
}

func checkZC1955(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "rfkill" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "block" {
		return nil
	}
	target := cmd.Arguments[1].String()
	if target != "all" && target != "wifi" && target != "wlan" &&
		target != "bluetooth" && target != "wwan" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1955",
		Message: "`rfkill block " + target + "` hard-downs the radio(s) — host drops " +
			"off the network in one call. Scope to the radio type that really needs it " +
			"and schedule an `at now + N minutes` unblock for self-recovery.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1956",
		Title:    "Error on `tailscale up --auth-key=SECRET` — single-use join key visible in argv",
		Severity: SeverityError,
		Description: "`tailscale up --auth-key tskey-auth-…` (and the joined `--auth-key=…` form) " +
			"passes the Tailscale pre-auth key as a command-line argument. Pre-auth keys grant " +
			"full tailnet membership, and short-lived or not, the value ends up in `ps`, " +
			"`/proc/PID/cmdline`, shell history, and any process dump taken before the join " +
			"completes. Read the key from `TS_AUTHKEY` with `tailscale up --authkey-env=TS_AUTHKEY` " +
			"(newer tailscaled), or from a file with `tailscale up --auth-key=file:/etc/ts.key` " +
			"(mode `0400` owned by the provisioning user).",
		Check: checkZC1956,
	})
}

func checkZC1956(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tailscale" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "--auth-key=") || strings.HasPrefix(v, "--authkey=") {
			val := v[strings.IndexByte(v, '=')+1:]
			if zc1956IsLiteralKey(val) {
				return zc1956Hit(cmd, v)
			}
		}
		if (v == "--auth-key" || v == "--authkey") && i+1 < len(cmd.Arguments) {
			val := cmd.Arguments[i+1].String()
			if zc1956IsLiteralKey(val) {
				return zc1956Hit(cmd, v+" "+val)
			}
		}
	}
	return nil
}

func zc1956IsLiteralKey(val string) bool {
	if val == "" {
		return false
	}
	if strings.HasPrefix(val, "file:") {
		return false
	}
	if strings.HasPrefix(val, "$") {
		return false
	}
	return true
}

func zc1956Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1956",
		Message: "`tailscale " + form + "` puts the pre-auth key in argv — visible in " +
			"`ps`/`/proc`/history/crash dumps. Use `--auth-key=file:/etc/ts.key` (mode 0400) " +
			"or `--authkey-env=TS_AUTHKEY`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1957",
		Title:    "Warn on `lvchange -an` / `vgchange -an` — deactivates a live LV/VG, risks mounted-fs corruption",
		Severity: SeverityWarning,
		Description: "`lvchange -an VG/LV` (and `vgchange -an VG` for the whole group) deactivates " +
			"a logical volume by removing its device-mapper entry. If the LV is mounted, " +
			"writes that the kernel has buffered but not yet flushed may be lost, and any " +
			"process holding an open fd on the filesystem gets EIO on the next syscall. " +
			"`umount` the mount first, stop any service keeping files open, verify with " +
			"`lsof` / `fuser`, and only then `lvchange -an`. For a scripted teardown, prefer " +
			"`umount` + `lvremove` with a recovery snapshot in hand.",
		Check: checkZC1957,
	})
}

func checkZC1957(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "lvchange" && ident.Value != "vgchange" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-an", "-aN":
			return zc1957Hit(cmd, ident.Value+" -an")
		case "--activate=n", "--activate=N":
			return zc1957Hit(cmd, ident.Value+" "+v)
		case "--activate":
			if i+1 < len(cmd.Arguments) {
				next := cmd.Arguments[i+1].String()
				if next == "n" || next == "N" {
					return zc1957Hit(cmd, ident.Value+" --activate n")
				}
			}
		}
	}
	return nil
}

func zc1957Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1957",
		Message: "`" + form + "` deactivates the LV/VG — unflushed writes on a mounted " +
			"fs may be lost, open fds see EIO. Umount and stop holders first, verify with " +
			"`lsof`/`fuser`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1958",
		Title:    "Warn on `helm upgrade --force` — delete-and-recreate resources, drops running pods",
		Severity: SeverityWarning,
		Description: "`helm upgrade RELEASE CHART --force` flips the upgrade strategy from " +
			"three-way-merge to `delete + create` for every resource Helm owns. Deployments " +
			"become new objects, Services lose their `clusterIP` for a beat, and any " +
			"`PodDisruptionBudget` is bypassed because the resource is deleted, not rolled " +
			"out. Use plain `helm upgrade` (three-way merge) or `--atomic` / `--wait` for a " +
			"supervised roll. Reserve `--force` for recovery after a failed upgrade with a " +
			"stuck resource, not routine deploys.",
		Check: checkZC1958,
	})
}

func checkZC1958(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "helm" && ident.Value != "helm3" {
		return nil
	}
	if len(cmd.Arguments) < 2 || cmd.Arguments[0].String() != "upgrade" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "--force" {
			return zc1958Hit(cmd)
		}
	}
	return nil
}

func zc1958Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1958",
		Message: "`helm upgrade --force` is delete+create — pods die, PodDisruptionBudget " +
			"is bypassed, Services reset their `clusterIP`. Use plain `helm upgrade` " +
			"(three-way merge) or `--atomic`/`--wait` for a supervised roll.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1959",
		Title:    "Warn on `trivy … --skip-db-update` / `--skip-update` — scans against a stale vulnerability DB",
		Severity: SeverityWarning,
		Description: "`trivy` embeds a vulnerability database that is rehydrated on every scan " +
			"unless the operator passes `--skip-db-update` (or `--skip-update` on older " +
			"releases). In CI the flag is tempting — each build then skips a 40 MB download — " +
			"but the scan then misses every CVE disclosed since the cached DB was last " +
			"refreshed. Keep the default download, or pre-populate the cache with " +
			"`trivy image --download-db-only` once per day in a scheduled job, and only pass " +
			"`--skip-db-update` inside the same job so every scan sees the fresh data.",
		Check: checkZC1959,
	})
}

func checkZC1959(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `trivy --skip-db-update image` mangles to name=`skip-db-update`.
	if ident.Value == "skip-db-update" || ident.Value == "skip-update" {
		return zc1959Hit(cmd, "trivy --"+ident.Value)
	}
	if ident.Value != "trivy" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--skip-db-update" || v == "--skip-update" {
			return zc1959Hit(cmd, "trivy "+v)
		}
	}
	return nil
}

func zc1959Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1959",
		Message: "`" + form + "` scans against the cached DB — every CVE disclosed " +
			"since last refresh is missed. Keep the default download, or run " +
			"`trivy --download-db-only` once per day in a scheduled job.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1960",
		Title:    "Warn on `az vm run-command invoke` / `aws ssm send-command` — arbitrary commands on remote VM",
		Severity: SeverityWarning,
		Description: "`az vm run-command invoke --command-id RunShellScript --scripts \"$CMD\"` " +
			"(and the AWS equivalent `aws ssm send-command --document-name AWS-RunShellScript " +
			"--parameters \"commands=['$CMD']\"`) runs arbitrary shell on the target instance " +
			"via the cloud control plane. The identity making the call is whatever role the " +
			"script's credentials carry; if `$CMD` is composed from any operator or attacker " +
			"input, the result is remote code execution through IAM. Gate the call behind " +
			"a shell-escape-safe templater, pin the document version / script to a reviewed " +
			"asset in blob / S3, and require MFA on the invoking role.",
		Check: checkZC1960,
	})
}

func checkZC1960(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "az":
		// Look for `az vm run-command invoke` subcommand sequence.
		if len(cmd.Arguments) >= 3 &&
			cmd.Arguments[0].String() == "vm" &&
			cmd.Arguments[1].String() == "run-command" &&
			(cmd.Arguments[2].String() == "invoke" || cmd.Arguments[2].String() == "create") {
			return zc1960Hit(cmd, "az vm run-command "+cmd.Arguments[2].String())
		}
	case "aws":
		if len(cmd.Arguments) >= 2 &&
			cmd.Arguments[0].String() == "ssm" &&
			cmd.Arguments[1].String() == "send-command" {
			return zc1960Hit(cmd, "aws ssm send-command")
		}
	case "gcloud":
		if len(cmd.Arguments) >= 2 &&
			cmd.Arguments[0].String() == "compute" &&
			cmd.Arguments[1].String() == "ssh" {
			for _, arg := range cmd.Arguments[2:] {
				v := arg.String()
				if v == "--command" || len(v) > 10 && v[:10] == "--command=" {
					return zc1960Hit(cmd, "gcloud compute ssh --command")
				}
			}
		}
	}
	return nil
}

func zc1960Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1960",
		Message: "`" + form + "` runs arbitrary shell on the VM via the cloud control " +
			"plane — operator-composed command strings become IAM-driven RCE. Pin to a " +
			"reviewed asset, template-escape input, require MFA.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1961",
		Title:    "Warn on `gcloud iam service-accounts keys create` — mints a long-lived service-account JSON key",
		Severity: SeverityWarning,
		Description: "`gcloud iam service-accounts keys create key.json --iam-account=SA@PROJECT` " +
			"exports an RSA key pair wrapped in a JSON file. Once written it is effectively a " +
			"forever-valid bearer credential: no automatic rotation, no refresh, and a single " +
			"\"leaked by a `cat key.json`\" is game-over. Prefer Workload Identity Federation " +
			"(`gcloud iam workload-identity-pools …`), short-lived impersonation via " +
			"`gcloud auth print-access-token --impersonate-service-account=SA`, or the " +
			"key-less GCE/GKE attached service account. Reserve static JSON keys for provably " +
			"off-platform callers.",
		Check: checkZC1961,
	})
}

func checkZC1961(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "gcloud" {
		return nil
	}
	if len(cmd.Arguments) < 4 {
		return nil
	}
	if cmd.Arguments[0].String() != "iam" ||
		cmd.Arguments[1].String() != "service-accounts" ||
		cmd.Arguments[2].String() != "keys" ||
		cmd.Arguments[3].String() != "create" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1961",
		Message: "`gcloud iam service-accounts keys create` mints a long-lived JSON key — " +
			"no auto-rotate, no refresh. Prefer Workload Identity Federation, " +
			"`--impersonate-service-account`, or the attached service account.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1962",
		Title:    "Warn on `kustomize build --load-restrictor=LoadRestrictionsNone` — path-traversal in overlays",
		Severity: SeverityWarning,
		Description: "Kustomize's default `LoadRestrictionsRootOnly` limits every base, patch, " +
			"configMapGenerator, and secretGenerator to paths under the current kustomization " +
			"root. `kustomize build … --load-restrictor=LoadRestrictionsNone` (also the legacy " +
			"spelling `--load_restrictor none` / `--load-restrictor=LoadRestrictionsNone_WarnForAll`) " +
			"drops that guard, so an overlay from an untrusted remote base can reference " +
			"`../../secrets/prod.env` or absolute paths and pull them into the render. Keep " +
			"the default; if a legitimate overlay needs a sibling file, vendor it in.",
		Check: checkZC1962,
	})
}

func checkZC1962(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kustomize" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "build" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if strings.HasPrefix(v, "--load-restrictor=") ||
			strings.HasPrefix(v, "--load_restrictor=") {
			val := v[strings.IndexByte(v, '=')+1:]
			if zc1962IsNoneVariant(val) {
				return zc1962Hit(cmd, v)
			}
		}
		if (v == "--load-restrictor" || v == "--load_restrictor") && i+2 <= len(cmd.Arguments)-1 {
			val := cmd.Arguments[i+2].String()
			if zc1962IsNoneVariant(val) {
				return zc1962Hit(cmd, v+" "+val)
			}
		}
	}
	return nil
}

func zc1962IsNoneVariant(val string) bool {
	switch val {
	case "none", "None", "LoadRestrictionsNone":
		return true
	}
	return false
}

func zc1962Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1962",
		Message: "`kustomize build " + form + "` drops path-root restriction — untrusted " +
			"overlays can reference `../../secrets/prod.env` and pull them into the render. " +
			"Keep the default; vendor sibling files into the overlay.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1963",
		Title:    "Warn on `npx pkg` / `pnpm dlx pkg` / `bunx pkg` without a version pin — runs latest registry code",
		Severity: SeverityWarning,
		Description: "`npx PKG`, `pnpm dlx PKG`, `bunx PKG`, and `bun x PKG` fetch the named " +
			"package from the npm registry and execute its `bin` entry. Without a version " +
			"pin (`pkg@1.2.3`), each run resolves to the registry's `latest` tag — a " +
			"compromised maintainer, squatted name, or even a mistyped package is enough to " +
			"land attacker code in the build. Pin the exact version (`npx pkg@1.2.3`), cache " +
			"the binary under `./node_modules/.bin/` via a regular `npm install`, or verify " +
			"the tarball signature before execution.",
		Check: checkZC1963,
	})
}

func checkZC1963(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var form string
	var pkgs []ast.Expression
	switch ident.Value {
	case "npx":
		form = "npx"
		pkgs = cmd.Arguments
	case "bunx":
		form = "bunx"
		pkgs = cmd.Arguments
	case "pnpm", "bun":
		if len(cmd.Arguments) < 2 {
			return nil
		}
		sub := cmd.Arguments[0].String()
		if ident.Value == "pnpm" && sub != "dlx" {
			return nil
		}
		if ident.Value == "bun" && sub != "x" {
			return nil
		}
		form = ident.Value + " " + sub
		pkgs = cmd.Arguments[1:]
	default:
		return nil
	}

	for _, arg := range pkgs {
		v := arg.String()
		if strings.HasPrefix(v, "-") {
			continue
		}
		if strings.Contains(v, "@") && !strings.HasPrefix(v, "@") {
			// pkg@version — pinned.
			return nil
		}
		if strings.HasPrefix(v, "@") {
			// scoped name like @scope/pkg — check for second @ (version).
			rest := v[1:]
			if strings.Contains(rest, "@") {
				return nil
			}
		}
		if strings.HasPrefix(v, "$") {
			return nil
		}
		return []Violation{{
			KataID: "ZC1963",
			Message: "`" + form + " " + v + "` pulls the `latest` tag every run — " +
				"a squatted or compromised package lands attacker code. Pin the version " +
				"(`" + v + "@X.Y.Z`) or use a regular `npm install` + `./node_modules/.bin/`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1964",
		Title:    "Warn on `uvx pkg` / `uv tool run pkg` / `pipx run pkg` without a version pin — runs latest PyPI release",
		Severity: SeverityWarning,
		Description: "`uvx PKG`, `uv tool run PKG`, and `pipx run PKG` each resolve the package " +
			"against PyPI and execute its entry point. Without a version constraint " +
			"(`pkg==1.2.3` or `pkg@1.2.3` for uv), every run takes whatever the registry " +
			"currently serves — a typosquatted lookalike, a compromised maintainer " +
			"release, or a sudden major-version bump lands untested code in the " +
			"pipeline. Pin the version at the call site or use `uv tool install pkg==X.Y.Z` + " +
			"`uv tool run pkg` so the lockfile is the source of truth.",
		Check: checkZC1964,
	})
}

func checkZC1964(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var form string
	var pkgs []ast.Expression
	switch ident.Value {
	case "uvx":
		form = "uvx"
		pkgs = cmd.Arguments
	case "uv":
		if len(cmd.Arguments) < 3 {
			return nil
		}
		if cmd.Arguments[0].String() != "tool" || cmd.Arguments[1].String() != "run" {
			return nil
		}
		form = "uv tool run"
		pkgs = cmd.Arguments[2:]
	case "pipx":
		if len(cmd.Arguments) < 2 {
			return nil
		}
		if cmd.Arguments[0].String() != "run" {
			return nil
		}
		form = "pipx run"
		pkgs = cmd.Arguments[1:]
	default:
		return nil
	}

	for _, arg := range pkgs {
		v := arg.String()
		if strings.HasPrefix(v, "-") {
			continue
		}
		if zc1964IsPinned(v) {
			return nil
		}
		if strings.HasPrefix(v, "$") {
			return nil
		}
		return []Violation{{
			KataID: "ZC1964",
			Message: "`" + form + " " + v + "` resolves to the PyPI `latest` release — " +
				"a squatted name or compromised maintainer lands untested code. Pin " +
				"`pkg==X.Y.Z` (or `pkg@X.Y.Z` for uv).",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func zc1964IsPinned(v string) bool {
	return strings.Contains(v, "==") || strings.Contains(v, "@") ||
		strings.Contains(v, ">=") || strings.Contains(v, "~=")
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1965",
		Title:    "Error on `systemd-cryptenroll --wipe-slot=all` — wipes every LUKS key slot",
		Severity: SeverityError,
		Description: "`systemd-cryptenroll --wipe-slot=all $DEV` removes every key slot on the " +
			"LUKS volume — passphrase, recovery key, TPM2, FIDO2, PKCS#11 — in one call. " +
			"`--wipe-slot=recovery` / `--wipe-slot=empty` are scoped; the `all` form is a " +
			"one-shot brick with no confirmation. Either enrol the new slot first and then " +
			"wipe the specific index you are retiring (`--wipe-slot=<n>`), or back up the " +
			"header with `cryptsetup luksHeaderBackup` before the call so recovery is " +
			"possible.",
		Check: checkZC1965,
	})
}

func checkZC1965(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `systemd-cryptenroll --wipe-slot=all $DEV` mangles the
	// command name to `wipe-slot=all`.
	if strings.HasPrefix(ident.Value, "wipe-slot=") {
		if ident.Value == "wipe-slot=all" {
			return zc1965Hit(cmd)
		}
	}
	if ident.Value != "systemd-cryptenroll" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--wipe-slot=all" || v == "--wipe-slot" {
			if v == "--wipe-slot=all" {
				return zc1965Hit(cmd)
			}
		}
	}
	return nil
}

func zc1965Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1965",
		Message: "`systemd-cryptenroll --wipe-slot=all` wipes every LUKS key slot " +
			"(passphrase/recovery/TPM2/FIDO2) in one call. Enrol the new slot first, " +
			"wipe a specific index, back up the header with `cryptsetup luksHeaderBackup`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1966",
		Title:    "Error on `zpool import -f` / `zpool export -f` — forced ZFS pool op bypasses hostid/txg checks",
		Severity: SeverityError,
		Description: "`zpool import -f $POOL` force-imports a pool even when the on-disk " +
			"hostid differs — i.e. the pool is already imported on another host " +
			"(multipath/SAN, shared JBOD, HA cluster). The second import writes to the " +
			"same vdevs and silently corrupts the pool. `zpool export -f` skips the " +
			"graceful-flush path and detaches vdevs with in-flight txgs, which can lose " +
			"the tail of the ZIL. Export without `-f` after `zfs unmount -a`; import " +
			"without `-f` after verifying `zpool import` (no target) reports the pool " +
			"as `ONLINE` and the hostid matches.",
		Check: checkZC1966,
	})
}

func checkZC1966(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "zpool" || len(cmd.Arguments) < 2 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "import" && sub != "export" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-f" || v == "--force" {
			return zc1966Hit(cmd, "zpool "+sub+" -f")
		}
	}
	return nil
}

func zc1966Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1966",
		Message: "`" + form + "` bypasses hostid/txg safety — forced import of a pool " +
			"already online elsewhere (SAN/HA) corrupts it; forced export drops in-flight " +
			"txgs. `zfs unmount -a` first, then plain `zpool export`/`import`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1967",
		Title:    "Warn on `setopt PROMPT_SUBST` — expansions inside `$PROMPT` evaluate command substitution every redraw",
		Severity: SeverityWarning,
		Description: "`setopt PROMPT_SUBST` turns on parameter, command, and arithmetic " +
			"substitution inside `$PS1`/`$PROMPT`/`$RPROMPT`. Any value that lands in the " +
			"prompt from an untrusted source — a git branch name, a checkout path, a " +
			"hostname in `/etc/hostname`, an env var set by a spawned tool — is reparsed " +
			"as shell code on every redraw, so a branch like `$(id>/tmp/p)` runs each time " +
			"the cursor returns. Prefer Zsh prompt escapes (`%n`, `%d`, `%~`, `%m`, " +
			"`vcs_info`) which already sanitise their inputs, or scope with `setopt " +
			"LOCAL_OPTIONS PROMPT_SUBST` inside the prompt-building function instead of " +
			"flipping the option globally.",
		Check: checkZC1967,
	})
}

func checkZC1967(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1967Canonical(arg.String())
		switch v {
		case "PROMPTSUBST":
			if enabling {
				return zc1967Hit(cmd, "setopt PROMPT_SUBST")
			}
		case "NOPROMPTSUBST":
			if !enabling {
				return zc1967Hit(cmd, "unsetopt NO_PROMPT_SUBST")
			}
		}
	}
	return nil
}

func zc1967Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1967Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1967",
		Message: "`" + form + "` re-runs command substitution on every prompt " +
			"redraw — a branch/host/dir value with `$(…)` executes each render. Prefer " +
			"`%n`/`%d`/`%~`/`vcs_info`, or scope via `LOCAL_OPTIONS`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1968",
		Title:    "Warn on `dnf versionlock add` / `yum versionlock add` — pins RPM, blocks CVE updates",
		Severity: SeverityWarning,
		Description: "`dnf versionlock add pkg` (and the legacy `yum versionlock add pkg`) " +
			"write an entry to `/etc/dnf/plugins/versionlock.list` that excludes the " +
			"package from future `dnf update` / `dnf upgrade` runs. Mirrors `apt-mark " +
			"hold` on Debian (ZC1550): the lock survives reboots and unattended-upgrades " +
			"never sees the newer rpm, so kernel, openssl, or glibc CVEs pile up unseen. " +
			"Document the exact reason in the commit, pair the lock with a scheduled " +
			"`dnf versionlock delete` date, and prefer excluding the problematic " +
			"transaction via `--exclude` or a one-shot `dnf update --setopt=exclude=pkg` " +
			"rather than a persistent pin.",
		Check: checkZC1968,
	})
}

func checkZC1968(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dnf" && ident.Value != "yum" && ident.Value != "microdnf" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "versionlock" {
		return nil
	}
	sub := cmd.Arguments[1].String()
	if sub != "add" && sub != "exclude" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1968",
		Message: "`" + ident.Value + " versionlock " + sub + "` pins the rpm — blocks " +
			"future CVE fixes for glibc/openssl/kernel. Prefer `--exclude` on a single " +
			"transaction and schedule a `versionlock delete` review.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1969",
		Title:    "Warn on `zsh -f` / `zsh -d` — skips `/etc/zsh*` and `~/.zsh*` startup files",
		Severity: SeverityWarning,
		Description: "`zsh -f` is the short form of `--no-rcs`, which skips every personal " +
			"and system-wide startup file: `/etc/zshenv`, `/etc/zprofile`, `/etc/zshrc`, " +
			"`/etc/zlogin`, `~/.zshenv`, `~/.zshrc`, `~/.zlogin`. `zsh -d` (`--no-" +
			"globalrcs`) drops only the `/etc/zsh*` set but keeps per-user ones. Either " +
			"form strips corp-mandated settings — proxy/hosts overrides, audit hooks, " +
			"umask, `HISTFILE` redirection, `PATH` hardening — silently. Use it " +
			"deliberately only for a pristine test harness or a minimal repro; never as " +
			"the shebang of a production script. When isolation is required, prefer " +
			"`env -i zsh` with an explicit allow-list of variables.",
		Check: checkZC1969,
	})
}

func checkZC1969(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "zsh" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-f" || v == "-d" {
			return zc1969Hit(cmd, "zsh "+v)
		}
	}
	return nil
}

func zc1969Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1969",
		Message: "`" + form + "` skips `/etc/zsh*` and `~/.zsh*` startup files — " +
			"corp proxy/audit/`PATH` hardening silently dropped. For a pristine " +
			"shell use `env -i zsh` with an explicit allow-list.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1970",
		Title:    "Warn on `losetup -P` / `kpartx -a` / `partprobe` on untrusted image — runs kernel partition parser",
		Severity: SeverityWarning,
		Description: "`losetup -P $LOOP $IMG`, `kpartx -av $IMG`, and `partprobe $LOOP` all " +
			"tell the kernel to rescan a block device's partition table and emit `/dev/" +
			"loopNpX` (or dm-N) entries. When the image comes from an untrusted source " +
			"— a customer-supplied VM disk, a downloaded installer, a forensic capture — " +
			"the scan runs MBR/GPT/LVM parsers over attacker-controlled bytes and has " +
			"historically triggered kernel CVEs (fsconfig heap overflow, ext4 mount " +
			"bugs). Do the inspection in a throwaway VM or an offline parser like " +
			"`fdisk -l $IMG` / `sfdisk --dump $IMG` that reads without kernel scan, and " +
			"only attach partitions with `losetup -P` after the layout is known-good.",
		Check: checkZC1970,
	})
}

var (
	zc1970LosetupFlags = map[string]struct{}{
		"-P": {}, "--partscan": {}, "-Pf": {}, "-fP": {}, "-rP": {}, "-Pr": {},
	}
	zc1970KpartxFlags = map[string]struct{}{
		"-a": {}, "-av": {}, "-va": {}, "-as": {}, "-sa": {},
	}
)

func checkZC1970(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	switch CommandIdentifier(cmd) {
	case "losetup":
		if HasArgFlag(cmd, zc1970LosetupFlags) {
			return zc1970Hit(cmd, "losetup -P")
		}
	case "kpartx":
		if HasArgFlag(cmd, zc1970KpartxFlags) {
			return zc1970Hit(cmd, "kpartx -a")
		}
	case "partprobe":
		return zc1970Hit(cmd, "partprobe")
	}
	return nil
}

func zc1970Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1970",
		Message: "`" + form + "` asks the kernel to parse the partition table of the " +
			"image — attacker-controlled bytes have tripped kernel CVEs. Use `fdisk " +
			"-l`/`sfdisk --dump` offline first, scan only known-good images.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1971",
		Title:    "Warn on `unsetopt GLOBAL_RCS` / `setopt NO_GLOBAL_RCS` — skips `/etc/zprofile`, `/etc/zshrc`, `/etc/zlogin`, `/etc/zlogout`",
		Severity: SeverityWarning,
		Description: "`GLOBAL_RCS` is on by default; only `/etc/zshenv` is sourced before it " +
			"can be toggled. Flipping the option off (either `unsetopt GLOBAL_RCS` or " +
			"`setopt NO_GLOBAL_RCS`) tells Zsh to skip `/etc/zprofile`, `/etc/zshrc`, " +
			"`/etc/zlogin`, and `/etc/zlogout` — which is exactly where admins put " +
			"corp-wide `PATH` hardening, audit hooks, umask, `HISTFILE` redirection, " +
			"and proxy variables. A login-shell script that disables the option in " +
			"`/etc/zshenv` neutralises every downstream system rc without a trace. " +
			"Keep the option on; if a specific helper needs pristine setup use " +
			"`emulate -LR zsh` inside a function or spawn `env -i zsh -f` scoped to " +
			"that helper.",
		Check: checkZC1971,
	})
}

func checkZC1971(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1971Canonical(arg.String())
		switch v {
		case "GLOBALRCS":
			if !enabling {
				return zc1971Hit(cmd, "unsetopt GLOBAL_RCS")
			}
		case "NOGLOBALRCS":
			if enabling {
				return zc1971Hit(cmd, "setopt NO_GLOBAL_RCS")
			}
		}
	}
	return nil
}

func zc1971Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1971Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1971",
		Message: "`" + form + "` tells Zsh to skip `/etc/zprofile`, `/etc/zshrc`, " +
			"`/etc/zlogin`, `/etc/zlogout` — corp `PATH`/audit/umask/proxy config " +
			"silently dropped. Keep on; scope pristine setup with `emulate -LR zsh` " +
			"or `env -i zsh -f`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1972",
		Title:    "Error on `dmsetup remove_all` / `dmsetup remove -f` — tears down live LVM/LUKS/multipath mappings",
		Severity: SeverityError,
		Description: "`dmsetup remove_all` iterates every device-mapper node on the host — " +
			"LVM logical volumes, LUKS containers, multipath aggregates, `cryptsetup` " +
			"mappings — and asks the kernel to drop each one. `dmsetup remove --force " +
			"$NAME` targets a single mapping but still evicts it with in-flight I/O. " +
			"When any of those devices is mounted or backing a running VM, new I/O to " +
			"it returns `ENXIO`, `fsck` is no longer possible, and LVM metadata needs a " +
			"cold reboot to reappear. Use `dmsetup remove $NAME` without `--force` " +
			"after `umount`/`vgchange -an`/`cryptsetup close`, and never `remove_all` " +
			"on a host you care about.",
		Check: checkZC1972,
	})
}

func checkZC1972(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dmsetup" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub == "remove_all" {
		return zc1972Hit(cmd, "dmsetup remove_all")
	}
	if sub == "remove" {
		for _, arg := range cmd.Arguments[1:] {
			v := arg.String()
			if v == "-f" || v == "--force" {
				return zc1972Hit(cmd, "dmsetup remove -f")
			}
		}
	}
	return nil
}

func zc1972Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1972",
		Message: "`" + form + "` drops LVM/LUKS/multipath mappings while still in " +
			"use — in-flight I/O returns `ENXIO`, metadata needs a reboot. `umount` " +
			"+ `vgchange -an` / `cryptsetup close` first, then `dmsetup remove`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1973",
		Title:    "Warn on `setopt POSIX_IDENTIFIERS` — restricts parameter names to ASCII, breaks Unicode `$var`",
		Severity: SeverityWarning,
		Description: "Zsh accepts Unicode parameter names by default: `$café`, `$π`, `$данные` " +
			"all parse. `setopt POSIX_IDENTIFIERS` tightens that to the POSIX subset — " +
			"ASCII letters, digits, underscore, not starting with a digit. Once the " +
			"option is on, every later `${café}` or `café=1` is a parse error, and " +
			"scripts/libraries that expose i18n-named vars stop loading. If you need " +
			"POSIX identifiers for a specific helper, scope it inside a function with " +
			"`emulate -LR sh`; leave the global option off so the rest of the shell " +
			"keeps the Zsh behaviour the user expects.",
		Check: checkZC1973,
	})
}

func checkZC1973(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1973Canonical(arg.String())
		switch v {
		case "POSIXIDENTIFIERS":
			if enabling {
				return zc1973Hit(cmd, "setopt POSIX_IDENTIFIERS")
			}
		case "NOPOSIXIDENTIFIERS":
			if !enabling {
				return zc1973Hit(cmd, "unsetopt NO_POSIX_IDENTIFIERS")
			}
		}
	}
	return nil
}

func zc1973Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1973Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1973",
		Message: "`" + form + "` restricts parameter names to ASCII; later " +
			"`${café}`/`${π}` fail to parse and i18n-named libs stop loading. " +
			"Scope with `emulate -LR sh` inside the helper instead of flipping " +
			"globally.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1974",
		Title:    "Error on `ipset flush` / `ipset destroy` — nukes named sets referenced by iptables/nft rules",
		Severity: SeverityError,
		Description: "`ipset flush` empties every entry from a named IP set; `ipset destroy` " +
			"(no args) removes every set on the host. iptables/nft rules of the form " +
			"`-m set --match-set $NAME src` then reference a set that is either empty " +
			"or gone, so block-lists disappear instantly and allow-lists stop " +
			"whitelisting — the ruleset falls through to its default policy. Target a " +
			"specific set by name (`ipset destroy $NAME` after confirming no rule " +
			"references it), or add new entries with `ipset add` instead of rebuilding " +
			"from scratch. Reload atomically with `ipset restore -! < snapshot` if a " +
			"full replace is genuinely needed.",
		Check: checkZC1974,
	})
}

func checkZC1974(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ipset" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	switch sub {
	case "flush", "-F":
		return zc1974Hit(cmd, "ipset flush")
	case "destroy", "-X":
		if len(cmd.Arguments) == 1 {
			return zc1974Hit(cmd, "ipset destroy")
		}
	}
	return nil
}

func zc1974Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1974",
		Message: "`" + form + "` drops named IP sets wholesale — iptables/nft rules " +
			"that reference them fall through to the default policy (block-list " +
			"empty, allow-list gone). Target by name; reload atomically via " +
			"`ipset restore -! < snapshot`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1975",
		Title:    "Warn on `unsetopt EXEC` / `setopt NO_EXEC` — parser keeps scanning, commands stop running",
		Severity: SeverityWarning,
		Description: "`EXEC` is on by default; the shell both parses and runs each command. " +
			"Turning it off (`unsetopt EXEC` or `setopt NO_EXEC`) tells Zsh to parse " +
			"everything but silently skip the execution step — nothing fires, yet " +
			"parameter assignments on the same line don't either, `$?` stays frozen, " +
			"and functions that follow look defined but never run. That is the " +
			"semantics behind `zsh -n script.zsh` for a pure syntax check; flipping " +
			"the option in the middle of a production script converts every later " +
			"line into a no-op without a visible error. Run syntax checks via `zsh " +
			"-n` from the outside, never by flipping `EXEC` in-line.",
		Check: checkZC1975,
	})
}

func checkZC1975(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1975Canonical(arg.String())
		switch v {
		case "EXEC":
			if !enabling {
				return zc1975Hit(cmd, "unsetopt EXEC")
			}
		case "NOEXEC":
			if enabling {
				return zc1975Hit(cmd, "setopt NO_EXEC")
			}
		}
	}
	return nil
}

func zc1975Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1975Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1975",
		Message: "`" + form + "` stops running commands but keeps parsing — every " +
			"later line becomes a silent no-op. For syntax checks run `zsh -n " +
			"script.zsh` from outside the script.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1976",
		Title:    "Error on `exportfs -au` / `exportfs -u` — unexports live NFS shares, clients get `ESTALE`",
		Severity: SeverityError,
		Description: "`exportfs -au` unexports every NFS share on the server; `exportfs -u " +
			"HOST:/PATH` removes a single share. Any client that currently has the " +
			"export mounted is not notified — the next read/write returns `ESTALE`, " +
			"the mount looks live but every open fd fails, and the only recovery is a " +
			"client-side `umount -l` + remount. `exportfs -f` (flush) is almost always " +
			"what you actually want after an `/etc/exports` edit; keep `-u`/`-au` for " +
			"planned shutdowns with a coordinated client `umount` first.",
		Check: checkZC1976,
	})
}

func checkZC1976(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "exportfs" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-au" || v == "-ua" || v == "-u" {
			return zc1976Hit(cmd, "exportfs "+v)
		}
	}
	return nil
}

func zc1976Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1976",
		Message: "`" + form + "` unexports live NFS shares — mounted clients see " +
			"`ESTALE` on every open fd. Use `exportfs -f` after editing " +
			"`/etc/exports`; reserve `-u`/`-au` for coordinated shutdowns.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1977",
		Title:    "Warn on `setopt CHASE_DOTS` — `cd ..` physically resolves before walking up, breaking logical paths",
		Severity: SeverityWarning,
		Description: "Default Zsh keeps `..` logical: from `/app/current/lib` (where " +
			"`/app/current` → `/app/releases/v5`), `cd ..` goes back to `/app/current`, " +
			"matching the user's mental model and blue/green deployment symlinks. " +
			"`setopt CHASE_DOTS` flips that — `..` first resolves the current directory " +
			"to its physical inode, so the same `cd ..` lands in `/app/releases/v5` " +
			"and the next `cd config` looks for `/app/releases/config` instead of " +
			"`/app/config`. Scripts that rely on `${PWD}` staying logical or on " +
			"`cd ../foo` matching the typed path break silently. Leave the option off; " +
			"use `cd -P` one-shot when a specific call really needs physical resolution.",
		Check: checkZC1977,
	})
}

func checkZC1977(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1977Canonical(arg.String())
		switch v {
		case "CHASEDOTS":
			if enabling {
				return zc1977Hit(cmd, "setopt CHASE_DOTS")
			}
		case "NOCHASEDOTS":
			if !enabling {
				return zc1977Hit(cmd, "unsetopt NO_CHASE_DOTS")
			}
		}
	}
	return nil
}

func zc1977Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1977Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1977",
		Message: "`" + form + "` makes `cd ..` physically resolve before walking up — " +
			"blue/green `current` symlinks stop working for `../foo` lookups. " +
			"Keep off; use `cd -P` one-shot when physical resolution is needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1978",
		Title:    "Warn on `tftp` — cleartext, unauthenticated UDP transfer",
		Severity: SeverityWarning,
		Description: "`tftp` has no authentication at all and moves the payload in plaintext " +
			"over UDP/69 — any packet capture on the path recovers the full transfer " +
			"and an attacker at the server can push an arbitrary file under the " +
			"expected name without noticing a lack of credentials. The dual-channel " +
			"design is also routinely mishandled by NAT/firewall gear. For PXE-style " +
			"provisioning that historically used `tftp`, fetch a signed payload over " +
			"HTTPS with `curl` and verify the signature locally before use. (See " +
			"ZC1200 for `ftp`, the authenticated-but-plaintext sibling.)",
		Check: checkZC1978,
	})
}

func checkZC1978(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	// `ftp` is owned by ZC1200; ZC1978 narrows to tftp (no auth, UDP).
	if ident.Value != "tftp" {
		return nil
	}
	// Require at least one arg so bare `tftp` at a prompt isn't flagged.
	if len(cmd.Arguments) == 0 {
		return nil
	}
	return []Violation{{
		KataID: "ZC1978",
		Message: "`tftp` transfers over plaintext UDP/69 with no authentication — " +
			"capture the payload, or push a crafted file under the expected " +
			"name. Use a signed-payload `curl` over HTTPS and verify the " +
			"signature before use.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1979",
		Title:    "Warn on `setopt HIST_FCNTL_LOCK` — `fcntl()` lock on NFS `$HISTFILE` stalls or deadlocks",
		Severity: SeverityWarning,
		Description: "Off by default, Zsh serialises writes to `$HISTFILE` with its own " +
			"lock-file dance next to the history. `setopt HIST_FCNTL_LOCK` switches to " +
			"POSIX `fcntl()` advisory locking — which is the safer primitive on local " +
			"filesystems, but on NFS homes the lock is proxied through `rpc.lockd` and " +
			"a single hung client or rebooted NFS server leaves every other shell " +
			"blocked the next time it tries to write history. The interactive shell " +
			"appears frozen on prompt return, and scripts that source user rc files " +
			"hang in `zshaddhistory`. Keep the option off on NFS homes; only turn it on " +
			"when `$HISTFILE` lives on a local filesystem (ext4, xfs, btrfs, zfs local " +
			"pool) that implements `fcntl()` without network round-trips.",
		Check: checkZC1979,
	})
}

func checkZC1979(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1979Canonical(arg.String())
		switch v {
		case "HISTFCNTLLOCK":
			if enabling {
				return zc1979Hit(cmd, "setopt HIST_FCNTL_LOCK")
			}
		case "NOHISTFCNTLLOCK":
			if !enabling {
				return zc1979Hit(cmd, "unsetopt NO_HIST_FCNTL_LOCK")
			}
		}
	}
	return nil
}

func zc1979Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1979Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1979",
		Message: "`" + form + "` routes `$HISTFILE` locking through POSIX `fcntl()` — " +
			"on NFS home directories a hung `rpc.lockd` freezes every other shell " +
			"at the next prompt. Keep off; enable only when `$HISTFILE` is on a " +
			"local fs.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1980",
		Title:    "Error on `udevadm trigger --action=remove` — replays `remove` uevents, detaches live devices",
		Severity: SeverityError,
		Description: "`udevadm trigger --action=remove` (also spelled `-c remove`) walks " +
			"`/sys` and synthesises a `remove` uevent for every matching device. The " +
			"kernel reacts as if every matched disk, NIC, GPU, or USB node was " +
			"physically yanked — SATA controllers detach drives that back mounted " +
			"filesystems, netdevs disappear mid-session, and `systemd-udevd` fires " +
			"per-device cleanup rules it was never meant to run on a live host. The " +
			"normal way to replay `add`/`change` events after a rules edit is " +
			"`udevadm control --reload` followed by `udevadm trigger` with the default " +
			"action (`change`); scope any `--action=remove` to a specific device " +
			"subsystem with `--subsystem-match=` + `--attr-match=` and test on a " +
			"non-production box first.",
		Check: checkZC1980,
	})
}

func checkZC1980(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "udevadm" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "trigger" {
		return nil
	}
	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if strings.HasPrefix(v, "--action=") {
			if strings.TrimPrefix(v, "--action=") == "remove" {
				return zc1980Hit(cmd)
			}
		}
		if v == "-c" || v == "--action" {
			if i+2 < len(cmd.Arguments) {
				next := cmd.Arguments[i+2].String()
				if next == "remove" {
					return zc1980Hit(cmd)
				}
			}
		}
	}
	return nil
}

func zc1980Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1980",
		Message: "`udevadm trigger --action=remove` replays `remove` uevents across " +
			"`/sys` — SATA/NIC/GPU nodes detach on a live host. Reload rules " +
			"with `udevadm control --reload`; scope with `--subsystem-match=`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1981",
		Title:    "Warn on `exec -a NAME cmd` — replaces `argv[0]`, hides the real binary from `ps`",
		Severity: SeverityWarning,
		Description: "`exec -a NAME $BIN` tells Zsh to set `argv[0]` of the `exec`'d process " +
			"to `NAME` instead of the actual program path. `ps`, `top`, `proc`-based " +
			"audit tools, and systemd's unit accounting all see `NAME` — the real " +
			"binary on disk is only discoverable from `/proc/PID/exe`, which most " +
			"monitoring does not read. The feature has legitimate uses (login shells " +
			"spelling themselves `-zsh` so tty/shell detection works) but also makes a " +
			"great disguise for a reverse shell or a cron-triggered helper. Keep " +
			"`exec -a` out of production scripts unless the intent is documented; " +
			"prefer running the binary at its real path so operators can match process " +
			"name to on-disk file.",
		Check: checkZC1981,
	})
}

func checkZC1981(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "exec" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if arg.String() == "-a" {
			return []Violation{{
				KataID: "ZC1981",
				Message: "`exec -a NAME` sets `argv[0]` to `NAME` — `ps`/`top`/audit " +
					"rules see the alias, not the real binary. Keep out of production " +
					"scripts unless the alias is documented.",
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
		ID:       "ZC1982",
		Title:    "Error on `ipcrm -a` — removes every SysV IPC object, breaks Postgres/Oracle/shm apps",
		Severity: SeverityError,
		Description: "`ipcrm -a` deletes every System V shared-memory segment, semaphore set, " +
			"and message queue owned by the caller (or, as root, every object on the " +
			"host). Long-running services that rely on SysV IPC — PostgreSQL's shared " +
			"buffers, Oracle's SGA, the `sysv` session store in several RDBMS test " +
			"suites, shm-based mutexes in batch pipelines — lose their backing store " +
			"mid-transaction and either SIGSEGV or return `EINVAL` on the next access. " +
			"Scope the removal: `ipcrm -m ID`/`-s ID`/`-q ID` against the specific " +
			"identifier reported by `ipcs -a`, after confirming no running process " +
			"attached to it.",
		Check: checkZC1982,
	})
}

func checkZC1982(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ipcrm" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if arg.String() == "-a" {
			return []Violation{{
				KataID: "ZC1982",
				Message: "`ipcrm -a` deletes every SysV shm/sem/mqueue object — " +
					"Postgres/Oracle/shm-based services lose their backing store " +
					"mid-transaction. Scope with `-m`/`-s`/`-q` on the specific ID " +
					"after checking `ipcs -a`.",
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
		ID:       "ZC1983",
		Title:    "Warn on `setopt CSH_JUNKIE_QUOTES` — single/double-quoted strings that span lines become errors",
		Severity: SeverityWarning,
		Description: "With `CSH_JUNKIE_QUOTES` off (the default), Zsh lets `\"foo\\nbar\"` and " +
			"`'line1\\nline2'` span physical lines. Setting the option on makes the " +
			"parser emit an error on the first newline inside a quoted string — which " +
			"breaks any existing multi-line SQL, JSON, or here-style payload that the " +
			"script has been inlining up to this point. Functions that are autoloaded " +
			"later or sourced from third-party helpers fail to parse, and the " +
			"diagnostic points at the closing quote, not at the option toggle. Leave " +
			"the option off; if csh-style strictness is genuinely required, scope " +
			"with `emulate -LR csh` inside the single helper that needs it.",
		Check: checkZC1983,
	})
}

func checkZC1983(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1983Canonical(arg.String())
		switch v {
		case "CSHJUNKIEQUOTES":
			if enabling {
				return zc1983Hit(cmd, "setopt CSH_JUNKIE_QUOTES")
			}
		case "NOCSHJUNKIEQUOTES":
			if !enabling {
				return zc1983Hit(cmd, "unsetopt NO_CSH_JUNKIE_QUOTES")
			}
		}
	}
	return nil
}

func zc1983Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1983Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1983",
		Message: "`" + form + "` makes every later multi-line `\"…\"`/`'…'` an error — " +
			"inlined SQL/JSON payloads and autoloaded helpers stop parsing. Scope " +
			"csh-style strictness with `emulate -LR csh` in the one helper that " +
			"needs it.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1984",
		Title:    "Error on `sgdisk -Z` / `sgdisk -o` — erases the GPT partition table on the target disk",
		Severity: SeverityError,
		Description: "`sgdisk -Z $DISK` (`--zap-all`) wipes the primary GPT, the protective " +
			"MBR, and the backup GPT at the end of the device. `sgdisk -o $DISK` " +
			"(`--clear`) replaces the existing partition table with a fresh empty GPT. " +
			"Either command detaches every partition, LVM PV, LUKS container, and " +
			"filesystem header on the device — when the target variable resolves to a " +
			"wrong path (tab completion, `$DISK` defaulted to `/dev/sda`), the host " +
			"becomes unbootable. Require an `lsblk $DISK` + `blkid $DISK` preflight in " +
			"the script, route the action through `--pretend` (`-t`) first, and keep a " +
			"`sgdisk --backup=/root/$DISK.gpt $DISK` image before any zap.",
		Check: checkZC1984,
	})
}

var zc1984ZapFlags = map[string]bool{
	"-Z":        true,
	"-o":        true,
	"--zap-all": true,
}

func checkZC1984(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "sgdisk" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1984ZapFlags[v] {
			line, col := FlagArgPosition(cmd, zc1984ZapFlags)
			return []Violation{{
				KataID: "ZC1984",
				Message: "`sgdisk " + v + "` erases the GPT on the target device — a wrong " +
					"`$DISK` detaches every partition/LVM/LUKS header and bricks boot. " +
					"`lsblk`/`blkid` preflight, `--backup` the old table, and test with " +
					"`-t`/`--pretend` first.",
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
		ID:       "ZC1985",
		Title:    "Warn on `setopt SH_FILE_EXPANSION` — expansion order flips from Zsh-native to sh/bash, `~` leaks",
		Severity: SeverityWarning,
		Description: "Default Zsh runs parameter expansion first, then filename/`~` " +
			"expansion — so a `VAR='~/cache'` keeps the tilde literal when you do " +
			"`mkdir -p -- $VAR` because the `~` never leaves the value. `setopt " +
			"SH_FILE_EXPANSION` (POSIX/sh ordering) flips the pass: filename expansion " +
			"runs first on the raw text, then parameter expansion happens, so the " +
			"same line suddenly makes the tilde resolve to `$HOME`, paths pointing at " +
			"`~evil/.cache` resolve into another user's home, and `=cmd` spellings " +
			"look up `$PATH` silently. Keep the option off; when a specific helper " +
			"needs POSIX ordering use `emulate -LR sh` inside that function.",
		Check: checkZC1985,
	})
}

func checkZC1985(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1985Canonical(arg.String())
		switch v {
		case "SHFILEEXPANSION":
			if enabling {
				return zc1985Hit(cmd, "setopt SH_FILE_EXPANSION")
			}
		case "NOSHFILEEXPANSION":
			if !enabling {
				return zc1985Hit(cmd, "unsetopt NO_SH_FILE_EXPANSION")
			}
		}
	}
	return nil
}

func zc1985Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1985Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1985",
		Message: "`" + form + "` flips expansion order to POSIX — a `~` or `=cmd` " +
			"sitting inside a `$VAR` value suddenly resolves, so a user-typed " +
			"`~other/.cache` escapes into another home. Scope with `emulate -LR " +
			"sh`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1986",
		Title:    "Warn on `touch -d` / `-t` / `-r` — explicit timestamp write is a common antiforensics pattern",
		Severity: SeverityWarning,
		Description: "`touch -d \"2 years ago\" $F`, `touch -t YYYYMMDDhhmm $F`, and `touch -r " +
			"$REF $F` all write the atime/mtime to a specific value rather than the " +
			"current clock. Legitimate uses exist — re-stamping a mirror to match " +
			"upstream, generating deterministic tarballs for reproducible-build " +
			"pipelines, `rsync --archive` edge cases — but the pattern also matches the " +
			"classic \"age the dropped file\" antiforensics trick where an attacker " +
			"normalises a new binary to look as old as its neighbours so `find -mtime`- " +
			"based triage misses it. Audit rules should flag these forms in production " +
			"scripts; in reproducible-build contexts, keep the timestamp derived from " +
			"`SOURCE_DATE_EPOCH` via `touch -d @$SOURCE_DATE_EPOCH` so operators can " +
			"recognise the intent at a glance.",
		Check: checkZC1986,
	})
}

func checkZC1986(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "touch" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-d", "-t", "-r":
			return []Violation{{
				KataID: "ZC1986",
				Message: "`touch " + v + "` writes a specific atime/mtime — also the " +
					"classic \"age the dropped file\" antiforensics pattern. Derive " +
					"from `$SOURCE_DATE_EPOCH` where the intent is deterministic.",
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
		ID:       "ZC1987",
		Title:    "Warn on `setopt BRACE_CCL` — `{a-z}` expands to each character instead of staying literal",
		Severity: SeverityWarning,
		Description: "`BRACE_CCL` is off by default: `echo {a-z}` stays literal `a-z` in Zsh, " +
			"which is what most scripts that only want the numeric range form " +
			"`{1..10}` actually expect. `setopt BRACE_CCL` promotes single-character " +
			"ranges and enumerations inside braces to csh-style character-class " +
			"expansion, so `echo {a-z}` suddenly prints every letter from `a` to `z` " +
			"and `echo {ABC}` becomes `A B C`. Any later command line that embeds " +
			"single-character ranges — regex fragments, hex masks, CI job names with " +
			"stage suffixes — expands unexpectedly. Leave the option off; use `{a..z}` " +
			"when a real range is wanted and quote literals that contain `{…}`.",
		Check: checkZC1987,
	})
}

func checkZC1987(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1987Canonical(arg.String())
		switch v {
		case "BRACECCL":
			if enabling {
				return zc1987Hit(cmd, "setopt BRACE_CCL")
			}
		case "NOBRACECCL":
			if !enabling {
				return zc1987Hit(cmd, "unsetopt NO_BRACE_CCL")
			}
		}
	}
	return nil
}

func zc1987Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1987Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1987",
		Message: "`" + form + "` promotes single-character braces to csh-style classes " +
			"— `{a-z}` now expands to every letter, `{ABC}` to `A B C`, breaking " +
			"regex/hex/CI-name literals. Use `{a..z}` when a real range is wanted.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1988",
		Title:    "Error on `nsupdate -y HMAC:NAME:SECRET` — TSIG key visible in argv and shell history",
		Severity: SeverityError,
		Description: "`nsupdate -y [alg:]name:base64secret` hands the TSIG shared secret " +
			"directly on the command line, so `ps auxf`, `/proc/PID/cmdline`, and " +
			"`$HISTFILE` all capture the key — and whoever owns the key can rewrite " +
			"any zone that trusts it (DNS hijack, MX hijack, ACME domain-validation " +
			"bypass). `nsupdate -k /etc/named/KEY` (or `-k $KEYFILE` with `0600` " +
			"perms) reads the key from disk without exposing it. If the secret must " +
			"come from a secret store, pipe it through `nsupdate -k /dev/stdin <<<\"$KEYFILE_CONTENTS\"` " +
			"so the raw material never lands in argv.",
		Check: checkZC1988,
	})
}

func checkZC1988(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "nsupdate" {
		return nil
	}
	for i, arg := range cmd.Arguments {
		if arg.String() == "-y" && i+1 < len(cmd.Arguments) {
			return []Violation{{
				KataID: "ZC1988",
				Message: "`nsupdate -y …` puts the TSIG key in argv — `ps`, " +
					"`/proc/*/cmdline`, and `$HISTFILE` all capture it. Use " +
					"`nsupdate -k $KEYFILE` with a `0600` keyfile instead.",
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
		ID:       "ZC1989",
		Title:    "Warn on `setopt REMATCH_PCRE` — `[[ =~ ]]` regex flips from POSIX ERE to PCRE, changes semantics",
		Severity: SeverityWarning,
		Description: "By default Zsh's `[[ $str =~ pattern ]]` uses POSIX extended regex " +
			"(ERE). `setopt REMATCH_PCRE` (after `zmodload zsh/pcre`) swaps the engine " +
			"to PCRE for every later match. Patterns that pass through both engines " +
			"change meaning subtly: `\\b` is a word boundary in PCRE but a literal `b` " +
			"in ERE, `\\d`/`\\s`/`\\w` work in PCRE but not ERE, lookahead/lookbehind " +
			"(`(?=…)`) parse in PCRE but error in ERE, and inline flags `(?i)` only " +
			"exist in PCRE. Flipping the option globally silently rewrites the " +
			"meaning of every existing regex — prefer an explicit `pcre_match`/`pcre_compile` " +
			"call when PCRE is needed, or a `setopt LOCAL_OPTIONS REMATCH_PCRE` inside " +
			"the single function that uses PCRE syntax.",
		Check: checkZC1989,
	})
}

func checkZC1989(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1989Canonical(arg.String())
		switch v {
		case "REMATCHPCRE":
			if enabling {
				return zc1989Hit(cmd, "setopt REMATCH_PCRE")
			}
		case "NOREMATCHPCRE":
			if !enabling {
				return zc1989Hit(cmd, "unsetopt NO_REMATCH_PCRE")
			}
		}
	}
	return nil
}

func zc1989Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1989Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1989",
		Message: "`" + form + "` swaps `[[ =~ ]]` from POSIX ERE to PCRE — `\\b`, " +
			"`\\d`, lookahead, `(?i)` change meaning across every later match. " +
			"Prefer `pcre_match` where PCRE is needed or scope with `LOCAL_OPTIONS`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1990",
		Title:    "Warn on `openssl passwd -crypt` / `-1` / `-apr1` — obsolete password hash formats",
		Severity: SeverityWarning,
		Description: "`openssl passwd -crypt` emits DES-crypt, 8-char truncated and crackable in " +
			"seconds on modern hardware. `-1` is FreeBSD-style MD5, unsuitable for " +
			"storage, long broken. `-apr1` is Apache's MD5-based variant with the same " +
			"weakness. Any hash produced by these flags lands in `/etc/shadow`, an " +
			"htpasswd file, or a database row where an attacker can offline-crack the " +
			"whole batch with a single GPU. Use `-5` (SHA-256-crypt), `-6` (SHA-512-" +
			"crypt), or prefer a dedicated KDF-based hasher — `mkpasswd -m yescrypt`, " +
			"`htpasswd -B` (bcrypt), or `argon2` — so brute-force cost scales with " +
			"hardware.",
		Check: checkZC1990,
	})
}

func checkZC1990(node ast.Node) []Violation {
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
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "passwd" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		switch v {
		case "-crypt", "-1", "-apr1":
			return []Violation{{
				KataID: "ZC1990",
				Message: "`openssl passwd " + v + "` emits a broken hash format — " +
					"DES/MD5 variants crack on a laptop. Use `-5` / `-6` or a " +
					"KDF-based hasher (`mkpasswd -m yescrypt`, `htpasswd -B`, " +
					"`argon2`).",
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
		ID:       "ZC1991",
		Title:    "Warn on `setopt CSH_NULLCMD` — bare `> file` raises an error instead of running `$NULLCMD`",
		Severity: SeverityWarning,
		Description: "Default Zsh executes `$NULLCMD` (initially `cat`) when a line has " +
			"redirections but no command, so `> file < input` copies input to file " +
			"and `< file` pages through it with `$READNULLCMD` (initially `more`). " +
			"`setopt CSH_NULLCMD` drops the Zsh convention and follows csh — any " +
			"command line without an explicit command is a parse error, regardless of " +
			"redirections. Scripts that rely on the bare-redirect idiom (log " +
			"truncation via `> $LOG`, drop-in includes via `< file`, piped filters " +
			"built from aliases) stop working with a confusing `parse error near '<'`. " +
			"Keep the option off; write `: > file` (or `true > file`) explicitly when " +
			"you mean to truncate.",
		Check: checkZC1991,
	})
}

func checkZC1991(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1991Canonical(arg.String())
		switch v {
		case "CSHNULLCMD":
			if enabling {
				return zc1991Hit(cmd, "setopt CSH_NULLCMD")
			}
		case "NOCSHNULLCMD":
			if !enabling {
				return zc1991Hit(cmd, "unsetopt NO_CSH_NULLCMD")
			}
		}
	}
	return nil
}

func zc1991Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1991Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1991",
		Message: "`" + form + "` makes `> file` / `< file` (no command) a parse error " +
			"— log truncation and bare-redirect idioms stop working. Write `: > " +
			"file` explicitly for truncation.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1992",
		Title:    "Warn on `pkexec cmd` — PolicyKit privilege elevation is historically bug-prone and hard to audit from scripts",
		Severity: SeverityWarning,
		Description: "`pkexec` lifts a command to the UID configured in a PolicyKit `.policy` " +
			"file — typically root — after consulting an authorisation agent. From a " +
			"non-interactive script the agent has no way to prompt, so the call " +
			"either depends on a pre-authorised `.policy` override or fails in a " +
			"confusing manner. The binary also has a poor CVE track record (CVE-2021-" +
			"4034 pwnkit, CVE-2017-16089, envvar handling bugs) and its audit trail is " +
			"split across journald and `/var/log/auth.log`. Use `sudo` with a targeted " +
			"`sudoers` drop-in for scripted privilege elevation, or run the script " +
			"under a systemd unit with `User=` / `AmbientCapabilities=` when specific " +
			"capabilities are needed.",
		Check: checkZC1992,
	})
}

func checkZC1992(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "pkexec" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	return []Violation{{
		KataID: "ZC1992",
		Message: "`pkexec` elevates via PolicyKit — no agent to prompt in a script, " +
			"poor CVE history (pwnkit), split audit trail. Use `sudo` with a " +
			"targeted `sudoers.d` drop-in or a systemd unit with " +
			"`User=`/`AmbientCapabilities=`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1993",
		Title:    "Warn on `setopt KSH_TYPESET` — `typeset var=$val` starts word-splitting the RHS",
		Severity: SeverityWarning,
		Description: "Off by default, Zsh treats every `typeset`/`declare` assignment like a " +
			"shell assignment: the whole RHS after `=` is one token, so `typeset " +
			"msg=\"a b c\"` produces a single-element string. `setopt KSH_TYPESET` " +
			"follows ksh instead — each word on the `typeset` line is its own " +
			"assignment or name, and the shell re-splits the RHS on whitespace. " +
			"Functions that used to accept `typeset path=$HOME/My Files` suddenly " +
			"treat `Files` as a second variable name, and `local` (an alias for " +
			"`typeset` inside functions) inherits the same change. Keep the option " +
			"off; if ksh compatibility is genuinely needed, scope with `emulate -LR " +
			"ksh` inside the helper that needs it.",
		Check: checkZC1993,
	})
}

func checkZC1993(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1993Canonical(arg.String())
		switch v {
		case "KSHTYPESET":
			if enabling {
				return zc1993Hit(cmd, "setopt KSH_TYPESET")
			}
		case "NOKSHTYPESET":
			if !enabling {
				return zc1993Hit(cmd, "unsetopt NO_KSH_TYPESET")
			}
		}
	}
	return nil
}

func zc1993Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1993Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1993",
		Message: "`" + form + "` re-splits the RHS of every later `typeset`/`local` — " +
			"`typeset path=$HOME/My Files` now treats `Files` as a second name. " +
			"Scope with `emulate -LR ksh` inside the one helper that needs it.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1994",
		Title:    "Error on `lvreduce -f` / `lvreduce -y` — shrinks the LV without checking the filesystem above",
		Severity: SeverityError,
		Description: "`lvreduce -L SIZE $LV` cuts the block device below an existing filesystem. " +
			"The confirmation prompt exists precisely because ext4/xfs/btrfs do not " +
			"shrink themselves — LVM happily lops off the tail even though the " +
			"filesystem still believes those blocks are allocated. `-f` / `-y` / " +
			"`--force` / `--yes` skip the prompt, and the next mount returns " +
			"corruption or missing files. Shrink the filesystem first with " +
			"`resize2fs $LV $NEWSIZE` (or `xfs_growfs` equivalent — xfs cannot shrink, " +
			"so offline backup + recreate), verify `df` / `fsck`, then `lvreduce " +
			"--resizefs` (which performs both steps atomically) instead of bypassing " +
			"the prompt.",
		Check: checkZC1994,
	})
}

func checkZC1994(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "lvreduce" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-f", "-y", "--force", "--yes":
			return []Violation{{
				KataID: "ZC1994",
				Message: "`lvreduce " + v + "` skips the shrink-confirmation prompt — " +
					"the filesystem above still believes the tail is allocated and " +
					"the next mount sees corruption. Shrink fs first, or use " +
					"`lvreduce --resizefs`.",
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
		ID:       "ZC1995",
		Title:    "Warn on `unsetopt BGNICE` — background jobs run at full interactive priority, starve the foreground",
		Severity: SeverityWarning,
		Description: "Default Zsh applies `nice +5` to every backgrounded job so long-running " +
			"work does not starve the interactive session. `unsetopt BGNICE` (or " +
			"`setopt NO_BGNICE`) turns that off and bg jobs compete at the same " +
			"priority as the foreground shell — SSH keystroke handling, editor " +
			"redraws, and `cmd &` batch fan-out all feel laggy, and a single CPU-" +
			"bound bg job can peg every core of a container it shares with a human " +
			"operator. Keep the option on; when a background job legitimately needs " +
			"full priority (audio pipeline, realtime simulator), wrap just that one " +
			"with `nice -n 0 -- cmd &` or a systemd unit with `Nice=` instead of " +
			"flipping globally.",
		Check: checkZC1995,
	})
}

func checkZC1995(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1995Canonical(arg.String())
		switch v {
		case "BGNICE":
			if !enabling {
				return zc1995Hit(cmd, "unsetopt BG_NICE")
			}
		case "NOBGNICE":
			if enabling {
				return zc1995Hit(cmd, "setopt NO_BG_NICE")
			}
		}
	}
	return nil
}

func zc1995Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1995Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1995",
		Message: "`" + form + "` drops the `nice +5` that bg jobs get by default — a " +
			"CPU-bound `cmd &` now competes with SSH/editor work. Wrap specific " +
			"jobs with `nice -n 0` or a systemd `Nice=` unit instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1996",
		Title:    "Warn on `unshare -U` / `-r` — unprivileged user namespace maps caller to root inside the NS",
		Severity: SeverityWarning,
		Description: "`unshare -U` opens a new user namespace and `-r` / `--map-root-user` " +
			"maps the caller's UID to `0` inside it. That's the foundation of rootless " +
			"containers (bubblewrap, podman rootless, flatpak) and is legitimate in " +
			"that context. It is also the standard opening move for a long list of " +
			"LPE chains — once you are uid `0` in a user namespace you can create " +
			"additional mount/net/cgroup namespaces, run `mount -t overlay` against " +
			"attacker-controlled dirs, and probe kernel attack surface that is " +
			"normally gated on `CAP_SYS_ADMIN`. Audit rules should flag the pattern " +
			"in production scripts; if a rootless runtime really needs it, route " +
			"through the runtime binary (`bwrap`, `podman --rootless`) so the invocation " +
			"is recognisable.",
		Check: checkZC1996,
	})
}

var zc1996ExactFlags = map[string]struct{}{
	"-U": {}, "-r": {}, "-Ur": {}, "-rU": {},
	"--user": {}, "--map-root-user": {},
}

func checkZC1996(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || CommandIdentifier(cmd) != "unshare" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1996FlagHit(v) {
			return zc1996Hit(cmd, "unshare "+v)
		}
	}
	return nil
}

func zc1996FlagHit(v string) bool {
	if _, hit := zc1996ExactFlags[v]; hit {
		return true
	}
	// Short-flag bundles like `-Urm` carry user / root mapping via U or r.
	if !strings.HasPrefix(v, "-") || strings.HasPrefix(v, "--") || len(v) <= 1 || strings.Contains(v, "=") {
		return false
	}
	return strings.ContainsAny(v[1:], "Ur")
}

func zc1996Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1996",
		Message: "`" + form + "` opens a user namespace and maps the caller to uid 0 " +
			"inside it — also the standard opening move for many kernel-LPE " +
			"chains. Route legit rootless runtimes through `bwrap`/`podman " +
			"--rootless` so the intent is clear.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1997",
		Title:    "Warn on `setopt HIST_NO_FUNCTIONS` — function definitions skipped from `$HISTFILE`, breaks forensic trail",
		Severity: SeverityWarning,
		Description: "Default Zsh writes every command you type, including function " +
			"definitions, to `$HISTFILE`. `setopt HIST_NO_FUNCTIONS` suppresses " +
			"storage of commands that define a function. On a multi-admin box or a " +
			"shared root account this breaks the forensic trail — the function the " +
			"attacker just defined (or that an operator typed before running the " +
			"destructive bit) vanishes from history while the invocation that used " +
			"it still shows, leaving responders with a command that references a " +
			"name that no longer exists on disk or in any log. Keep the option off " +
			"and scope any hiding needs with the Zsh hook `zshaddhistory { return " +
			"1 }` inside a function where the secret actually lives.",
		Check: checkZC1997,
	})
}

func checkZC1997(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1997Canonical(arg.String())
		switch v {
		case "HISTNOFUNCTIONS":
			if enabling {
				return zc1997Hit(cmd, "setopt HIST_NO_FUNCTIONS")
			}
		case "NOHISTNOFUNCTIONS":
			if !enabling {
				return zc1997Hit(cmd, "unsetopt NO_HIST_NO_FUNCTIONS")
			}
		}
	}
	return nil
}

func zc1997Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1997Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1997",
		Message: "`" + form + "` drops function-definition commands from `$HISTFILE` " +
			"— forensic trail loses the definition while the call that used it " +
			"still shows. Scope hiding via `zshaddhistory` hook instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1998",
		Title:    "Error on `tpm2_clear` / `tpm2 clear` — wipes TPM storage hierarchy, kills every sealed key",
		Severity: SeverityError,
		Description: "`tpm2_clear -c p` (or `tpm2 clear -c p`) invokes the TPM 2.0 `TPM2_Clear` " +
			"command, which invalidates every object sealed against the storage " +
			"hierarchy — LUKS-TPM2 keyslots, systemd-cryptenroll's `--tpm2-device` " +
			"slot, sshd TPM-backed host keys, and SecureBoot measured-boot state. " +
			"The machine can still boot but any disk that unlocked through the TPM " +
			"now needs a recovery passphrase, and every TLS cert issued from a " +
			"TPM-sealed CA loses its anchor. There is no undo. Run `tpm2_clear` only " +
			"under a documented recovery runbook with the recovery material in hand; " +
			"never put it in an automated scheduled script.",
		Check: checkZC1998,
	})
}

func checkZC1998(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value == "tpm2_clear" {
		return zc1998Hit(cmd, "tpm2_clear")
	}
	if ident.Value == "tpm2" && len(cmd.Arguments) > 0 &&
		cmd.Arguments[0].String() == "clear" {
		return zc1998Hit(cmd, "tpm2 clear")
	}
	return nil
}

func zc1998Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1998",
		Message: "`" + form + "` wipes the TPM storage hierarchy — every LUKS-TPM2 " +
			"keyslot, `systemd-cryptenroll --tpm2-device` slot, and TPM-sealed " +
			"TLS/sshd key is destroyed. No undo. Gate behind a recovery runbook.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1999",
		Title:    "Error on `setopt AUTO_NAMED_DIRS` — unknown option, typo of `AUTO_NAME_DIRS`",
		Severity: SeverityError,
		Description: "`AUTO_NAMED_DIRS` (with the trailing `D`) is not a real Zsh option — " +
			"`setopt AUTO_NAMED_DIRS` fails with `no such option` and the dir-to-" +
			"`~name` auto-registration the author likely wanted is never enabled. " +
			"The canonical spelling is `AUTO_NAME_DIRS` (see ZC1934 for its " +
			"semantics and why flipping it on is usually wrong). Drop the typo and, " +
			"if you actually want the behaviour, reach for `hash -d NAME=PATH` " +
			"explicitly or scope `setopt LOCAL_OPTIONS AUTO_NAME_DIRS` inside the " +
			"single helper that needs named-directory expansion.",
		Check: checkZC1999,
	})
}

func checkZC1999(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setopt" && ident.Value != "unsetopt" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := zc1999Canonical(arg.String())
		switch v {
		case "AUTONAMEDDIRS", "NOAUTONAMEDDIRS":
			return []Violation{{
				KataID: "ZC1999",
				Message: "`" + ident.Value + " " + arg.String() + "` is a typo — the real " +
					"Zsh option is `AUTO_NAME_DIRS` (no trailing `D`, see ZC1934). " +
					"Fix the spelling or drop the toggle; `hash -d NAME=PATH` is " +
					"the explicit alternative.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1999Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}
