// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1600",
		Title:    "Warn on bare `chroot DIR CMD` — missing `--userspec=` keeps uid 0 inside the jail",
		Severity: SeverityWarning,
		Description: "`chroot` changes the filesystem root but does not drop privileges. The " +
			"caller is almost always root (the syscall needs `CAP_SYS_CHROOT`), and without " +
			"`--userspec=USER:GROUP` the command inside the chroot still runs as uid 0. It can " +
			"write anywhere inside the tree, chmod binaries, and — if proc / sys / device nodes " +
			"are bind-mounted in — escape. Pass `--userspec=` to run the command as a named " +
			"unprivileged user, or drop to a dedicated helper (bubblewrap, firejail) that also " +
			"unshares user namespaces.",
		Check: checkZC1600,
	})
}

func checkZC1600(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chroot" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1600",
		Message: "`chroot` without `--userspec=` runs the inner command as uid 0. Pass " +
			"`--userspec=USER:GROUP` to drop privileges, or use `bwrap` / `firejail` for " +
			"user-namespace isolation.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1601",
		Title:    "Warn on `ethtool -s $IF wol <g|u|m|b|a>` — enables remote Wake-on-LAN",
		Severity: SeverityWarning,
		Description: "Wake-on-LAN powers the host on from a sleep / soft-off state when a " +
			"matching packet reaches the NIC. The wake logic fires in a privileged firmware " +
			"path long before the kernel boots and firewall rules are loaded — so any packet " +
			"that reaches the interface (magic-packet, unicast, broadcast, ARP) triggers the " +
			"power-on unfiltered. On a shared or public LAN attackers on the broadcast domain " +
			"can wake hosts at will. Keep `wol d` (disable) unless a documented operational " +
			"need requires one of the wake bits.",
		Check: checkZC1601,
	})
}

func checkZC1601(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ethtool" {
		return nil
	}
	if len(cmd.Arguments) < 4 {
		return nil
	}
	if cmd.Arguments[0].String() != "-s" {
		return nil
	}

	for i := 2; i+1 < len(cmd.Arguments); i++ {
		if cmd.Arguments[i].String() != "wol" {
			continue
		}
		bits := cmd.Arguments[i+1].String()
		if bits == "d" {
			return nil
		}
		enables := false
		for _, c := range bits {
			switch c {
			case 'g', 'u', 'm', 'b', 'a', 'p', 's', 'f':
				enables = true
			}
		}
		if !enables {
			return nil
		}
		return []Violation{{
			KataID: "ZC1601",
			Message: "`ethtool -s " + cmd.Arguments[1].String() + " wol " + bits + "` " +
				"enables Wake-on-LAN — the NIC powers the host on before firewall rules " +
				"load. Keep `wol d` unless a documented operational need requires " +
				strings.TrimSpace(bits) + ".",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1602",
		Title:    "Warn on `setopt KSH_ARRAYS` / `SH_WORD_SPLIT` — flips Zsh core semantics shell-wide",
		Severity: SeverityWarning,
		Description: "`KSH_ARRAYS` makes arrays 0-indexed (the Bash / ksh convention), breaking " +
			"every Zsh access that uses `[1]` for the first element. `SH_WORD_SPLIT` makes " +
			"unquoted `$var` word-split on `IFS`, breaking the core Zsh promise that `echo " +
			"$x` passes exactly one argument. Setting either globally is a bug-magnet — pre-" +
			"existing code silently misbehaves from that line on. If you need the semantics " +
			"only inside a function, scope it with `emulate -L ksh` or `emulate -L sh`.",
		Check: checkZC1602,
	})
}

func checkZC1602(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setopt" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		raw := arg.String()
		norm := strings.ToLower(strings.ReplaceAll(raw, "_", ""))
		if norm == "ksharrays" || norm == "shwordsplit" {
			return []Violation{{
				KataID: "ZC1602",
				Message: "`setopt " + raw + "` flips Zsh core semantics for the whole shell " +
					"— pre-existing code silently misbehaves. Scope with `emulate -L ksh` / " +
					"`emulate -L sh` inside the function that needs the mode.",
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
		ID:       "ZC1603",
		Title:    "Warn on `gdb -p PID` / `ltrace -p PID` — live attach reads target memory",
		Severity: SeverityWarning,
		Description: "`gdb -p PID` and `ltrace -p PID` attach via ptrace and hand the caller " +
			"full read / write access to the target process: registers, heap, stack, open file " +
			"descriptors, and every environment variable. Credentials in `$AWS_SECRET_ACCESS_" +
			"KEY`, session tokens on the stack, TLS keys in memory — all readable. A root-run " +
			"script that attaches to another user's process extracts everything that user has. " +
			"Keep production scripts out of the debugger; if post-mortem diagnostics are " +
			"needed, use `coredumpctl` against a captured core file instead.",
		Check: checkZC1603,
	})
}

func checkZC1603(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "gdb" && ident.Value != "ltrace" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-p" {
			return []Violation{{
				KataID: "ZC1603",
				Message: "`" + ident.Value + " -p PID` attaches via ptrace — memory, " +
					"registers, env, and stack of the target are readable. Use " +
					"`coredumpctl` on a captured core, not a live attach from a script.",
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
		ID:       "ZC1604",
		Title:    "Warn on `source <glob>` / `. <glob>` — loads every match; one bad file = code exec",
		Severity: SeverityWarning,
		Description: "`source /etc/profile.d/*.sh` and similar glob-sourcing patterns load every " +
			"file that matches, in the order Zsh enumerates them. One attacker-writable file " +
			"anywhere in the glob yields arbitrary code execution as whoever is running the " +
			"script, with that caller's privileges. Prefer explicit filenames so review can " +
			"enumerate exactly what gets loaded. If a directory of drop-ins is required, audit " +
			"ownership and perms at install time and keep the directory root-owned.",
		Check: checkZC1604,
	})
}

func checkZC1604(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "source" && ident.Value != "." {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	target := cmd.Arguments[0].String()
	if !strings.ContainsAny(target, "*?[") {
		return nil
	}

	return []Violation{{
		KataID: "ZC1604",
		Message: "`" + ident.Value + " " + target + "` loads every matched file. One " +
			"attacker-writable match is arbitrary code execution. Use explicit filenames.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1605",
		Title:    "Error on `debugfs -w DEV` — write-mode filesystem debugger bypasses journal",
		Severity: SeverityError,
		Description: "`debugfs -w` opens the filesystem in write mode. It sidesteps the kernel's " +
			"normal write path — the journal doesn't see the changes, filesystem locks are " +
			"ignored, and inodes / blocks can be edited directly. On a mounted filesystem this " +
			"corrupts state silently; even on an unmounted one, the operator can repoint a " +
			"directory entry at an arbitrary inode. Scripts should never need this — keep " +
			"`debugfs -w` as an interactive last-resort from a rescue environment.",
		Check: checkZC1605,
	})
}

func checkZC1605(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "debugfs" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-w" {
			return []Violation{{
				KataID: "ZC1605",
				Message: "`debugfs -w` writes to the filesystem outside the kernel's normal " +
					"path — journal bypassed, locks ignored. Keep it as an interactive " +
					"rescue tool, not a script path.",
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
		ID:       "ZC1606",
		Title:    "Warn on `mkdir -m NNN` / `install -m NNN` with world-write bit (no sticky)",
		Severity: SeverityWarning,
		Description: "`mkdir -m 777 /path` and `install -m 777 src /dest` create a path that " +
			"every local user can write and rename inside. If the script later creates files " +
			"there, classic TOCTOU symlink attacks become trivial — the attacker drops a " +
			"symlink named like the expected output file, redirecting the write wherever they " +
			"choose. A sticky-bit mode (`1777`) mitigates this for shared temp dirs. Prefer " +
			"`mkdir -m 700` (or 750), and scope access by group or ACL rather than everyone.",
		Check: checkZC1606,
	})
}

var zc1606Names = map[string]struct{}{"mkdir": {}, "install": {}}

func checkZC1606(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	name := CommandIdentifier(cmd)
	if _, hit := zc1606Names[name]; !hit {
		return nil
	}
	mode := zc1606WorldWriteMode(cmd)
	if mode == "" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1606",
		Message: "`" + name + " -m " + mode + "` creates a world-writable path " +
			"without the sticky bit — TOCTOU symlink-attack ground. Use `-m 700` / " +
			"`-m 750`, or `-m 1777` if a shared sticky dir is actually needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func zc1606WorldWriteMode(cmd *ast.SimpleCommand) string {
	for i := 0; i+1 < len(cmd.Arguments); i++ {
		if cmd.Arguments[i].String() != "-m" {
			continue
		}
		mode := cmd.Arguments[i+1].String()
		if zc1606IsWorldWritable(mode) {
			return mode
		}
	}
	return ""
}

func zc1606IsWorldWritable(mode string) bool {
	if len(mode) != 3 {
		return false
	}
	for _, c := range mode {
		if c < '0' || c > '7' {
			return false
		}
	}
	switch mode[2] {
	case '2', '3', '6', '7':
		return true
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1607",
		Title:    "Warn on `git config safe.directory '*'` — disables CVE-2022-24765 protection",
		Severity: SeverityWarning,
		Description: "`safe.directory` is git's mitigation for CVE-2022-24765 (fake git dirs " +
			"planted by another uid). Setting it to `'*'` trusts every directory on the host " +
			"— an attacker who creates `/tmp/evil/.git` with a malicious `core.fsmonitor` hook " +
			"gets arbitrary code execution the first time any user runs `git status` near that " +
			"path. List the specific paths that need cross-owner git access instead, or fix " +
			"the underlying ownership mismatch.",
		Check: checkZC1607,
	})
}

func checkZC1607(node ast.Node) []Violation {
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

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "safe.directory=*" || strings.HasPrefix(v, "safe.directory=*") {
			return violationZC1607(cmd)
		}
		if v == "safe.directory" && i+1 < len(cmd.Arguments) {
			next := strings.Trim(cmd.Arguments[i+1].String(), "'\"")
			if next == "*" {
				return violationZC1607(cmd)
			}
		}
	}
	return nil
}

func violationZC1607(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1607",
		Message: "`git config safe.directory '*'` trusts every directory — defeats CVE-" +
			"2022-24765 protection. List specific paths, or fix the ownership mismatch.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1608",
		Title:    "Warn on `find -exec sh -c '... {} ...'` — filename in quoted script is injectable",
		Severity: SeverityWarning,
		Description: "Substituting `{}` directly into the quoted command string of `find -exec " +
			"sh -c` lets filenames with shell metacharacters break out. A file named `$(rm " +
			"-rf ~)` invokes command substitution; a file named `foo; curl evil` chains a " +
			"second command. Pass `{}` as a positional argument to `sh` so the filename " +
			"arrives as a parameter, not as source: `find -exec sh -c 'grep pat \"$1\"' _ {} " +
			"\\;`.",
		Check: checkZC1608,
	})
}

func checkZC1608(node ast.Node) []Violation {
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

	hasExec := false
	hasShellC := false
	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-exec" || v == "-execdir" {
			hasExec = true
			continue
		}
		if !hasExec {
			continue
		}
		if (v == "sh" || v == "bash" || v == "zsh" || v == "/bin/sh" || v == "/bin/bash") &&
			i+1 < len(cmd.Arguments) && cmd.Arguments[i+1].String() == "-c" {
			hasShellC = true
			continue
		}
		if !hasShellC {
			continue
		}
		if (strings.HasPrefix(v, "'") || strings.HasPrefix(v, "\"")) &&
			strings.Contains(v, "{}") {
			return []Violation{{
				KataID: "ZC1608",
				Message: "`find -exec sh -c '... {} ...'` interpolates filenames into the " +
					"shell script — metacharacters break out. Pass `{}` as a positional " +
					"arg: `sh -c '... \"$1\"' _ {} \\;`.",
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
		ID:       "ZC1609",
		Title:    "Warn on `aa-disable` / `aa-complain` / `apparmor_parser -R` — disables AppArmor enforcement",
		Severity: SeverityWarning,
		Description: "`aa-disable` fully unloads the named AppArmor profile; `aa-complain` " +
			"flips the profile from enforce to complain (violations are logged but allowed); " +
			"`apparmor_parser -R` removes a profile from the running kernel. Each one lets the " +
			"confined process run without its mandatory-access-control restrictions — if the " +
			"profile existed for a reason, that reason is now unenforced. Interactive debugging " +
			"is legitimate, but scripts that permanently disable profiles should be reviewed.",
		Check: checkZC1609,
	})
}

func checkZC1609(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "aa-disable", "aa-complain":
		if len(cmd.Arguments) == 0 {
			return nil
		}
		return []Violation{{
			KataID: "ZC1609",
			Message: "`" + ident.Value + "` disables or softens the AppArmor profile — the " +
				"confined process loses MAC restrictions. Review the profile's intent " +
				"before disabling in automation.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	case "apparmor_parser":
		for _, arg := range cmd.Arguments {
			if arg.String() == "-R" || arg.String() == "--remove" {
				return []Violation{{
					KataID: "ZC1609",
					Message: "`apparmor_parser -R` removes the AppArmor profile from the " +
						"kernel — the confined process loses MAC restrictions. Review " +
						"the profile's intent before removing in automation.",
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
		ID:       "ZC1610",
		Title:    "Warn on `curl -o /etc/...` / `wget -O /etc/...` — direct download to a system path",
		Severity: SeverityWarning,
		Description: "Writing the body of an HTTP response straight into `/etc/`, `/usr/`, " +
			"`/bin/`, `/sbin/`, or `/lib/` skips every integrity check the system usually " +
			"applies. If the URL is compromised or MITM'd, the attacker's content replaces a " +
			"system config or binary the next command over. Download to a temp file, verify " +
			"signature / checksum, and `install -m 0644` the final file into place. Package " +
			"managers exist for a reason — prefer them for system files.",
		Check: checkZC1610,
	})
}

func checkZC1610(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "curl" && ident.Value != "wget" {
		return nil
	}

	systemPrefixes := []string{"/etc/", "/usr/", "/bin/", "/sbin/", "/lib/", "/lib64/", "/opt/"}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		isOutFlag := v == "-o" || v == "-O" || v == "--output" || v == "--output-document"
		if isOutFlag && i+1 < len(cmd.Arguments) {
			next := cmd.Arguments[i+1].String()
			for _, p := range systemPrefixes {
				if strings.HasPrefix(next, p) {
					return []Violation{{
						KataID: "ZC1610",
						Message: "`" + ident.Value + " " + v + " " + next + "` writes " +
							"an HTTP response straight into a system path — a compromised " +
							"URL replaces the target. Download to a temp file, verify, " +
							"then `install` into place.",
						Line:   cmd.Token.Line,
						Column: cmd.Token.Column,
						Level:  SeverityWarning,
					}}
				}
			}
		}
		// Handle --output=/path and --output-document=/path joined forms.
		for _, prefix := range []string{"--output=", "--output-document="} {
			if strings.HasPrefix(v, prefix) {
				path := strings.TrimPrefix(v, prefix)
				for _, p := range systemPrefixes {
					if strings.HasPrefix(path, p) {
						return []Violation{{
							KataID: "ZC1610",
							Message: "`" + ident.Value + " " + v + "` writes an HTTP " +
								"response straight into a system path — a compromised " +
								"URL replaces the target. Download to a temp file, " +
								"verify, then `install` into place.",
							Line:   cmd.Token.Line,
							Column: cmd.Token.Column,
							Level:  SeverityWarning,
						}}
					}
				}
			}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1611",
		Title:    "Style: `${var^^}` / `${var,,}` — prefer Zsh `${(U)var}` / `${(L)var}` for case change",
		Severity: SeverityStyle,
		Description: "`${var^^}` (uppercase) and `${var,,}` (lowercase) came from Bash 4. Zsh " +
			"accepts them for compatibility but the idiomatic form is the parameter-expansion " +
			"flag: `${(U)var}` / `${(L)var}`. The flag is also available per-element in " +
			"arrays (`${(U)array}`) and composes with other flags (`${(UL)array}` doesn't " +
			"make sense, but `${(U)${(f)str}}` does). Prefer the Zsh-native form in a `.zsh` " +
			"script; it keeps the codebase consistent with other `(X)var` patterns.",
		Check: checkZC1611,
	})
}

func checkZC1611(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if !strings.Contains(v, "${") {
			continue
		}
		if strings.Contains(v, "^^}") || strings.Contains(v, ",,}") {
			return []Violation{{
				KataID: "ZC1611",
				Message: "`${var^^}` / `${var,,}` — prefer Zsh `${(U)var}` / `${(L)var}` " +
					"for case conversion.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}
	return nil
}

var zc1612HardeningDisables = map[string]string{
	"kernel.yama.ptrace_scope=0":         "YAMA ptrace scope (lets any process attach)",
	"kernel.kptr_restrict=0":             "kernel pointer restriction (leaks kptrs to /proc)",
	"kernel.dmesg_restrict=0":            "dmesg restriction (unprivileged users read ring buffer)",
	"kernel.unprivileged_bpf_disabled=0": "unprivileged BPF gate (any user loads BPF)",
	"net.core.bpf_jit_harden=0":          "BPF JIT hardening (JIT-spray mitigations off)",
	"kernel.perf_event_paranoid=-1":      "perf_event paranoid (unprivileged perf access)",
	"kernel.perf_event_paranoid=0":       "perf_event paranoid (unprivileged perf access)",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1612",
		Title:    "Warn on `sysctl -w` disabling kernel hardening knobs",
		Severity: SeverityWarning,
		Description: "Several sysctl knobs exist specifically to constrain what unprivileged " +
			"users can do — `kernel.yama.ptrace_scope`, `kernel.kptr_restrict`, " +
			"`kernel.dmesg_restrict`, `kernel.unprivileged_bpf_disabled`, " +
			"`net.core.bpf_jit_harden`, and `kernel.perf_event_paranoid`. Setting any of them " +
			"to the lowest-restriction value removes a distinct defense-in-depth layer: " +
			"unrelated processes can ptrace each other, kernel pointers leak to `/proc`, " +
			"unprivileged users read kernel ring buffers, BPF JIT-spray mitigations disappear. " +
			"Leave these defaults alone unless a measured performance or debugging need justifies it.",
		Check: checkZC1612,
	})
}

func checkZC1612(node ast.Node) []Violation {
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
		if note, ok := zc1612HardeningDisables[v]; ok {
			return []Violation{{
				KataID: "ZC1612",
				Message: "`sysctl ... " + v + "` disables " + note + " — defense-in-depth " +
					"loss. Leave the default unless a measured need justifies it.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

var zc1613Readers = map[string]bool{
	"cat": true, "less": true, "more": true,
	"head": true, "tail": true, "wc": true,
	"grep": true, "awk": true, "sed": true, "cut": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1613",
		Title:    "Warn on reading SSH private-key files with `cat` / `less` / `grep` / `head`",
		Severity: SeverityWarning,
		Description: "Piping an SSH private key through a generic text tool copies the raw " +
			"key material into the process and — if stdout is redirected or piped — often " +
			"into logs, backup files, or a terminal scrollback buffer. Host keys under " +
			"`/etc/ssh/ssh_host_*_key` impersonate the server; user keys under `~/.ssh/id_*` " +
			"impersonate the user. Use `ssh-keygen -l -f KEY` for fingerprint / metadata, or " +
			"pass the key path to the consumer directly (`ssh -i`, `git -c core.sshCommand`) " +
			"without staging it through a shell tool.",
		Check: checkZC1613,
	})
}

func checkZC1613(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if !zc1613Readers[ident.Value] {
		return nil
	}

	userSuffixes := []string{
		"/.ssh/id_rsa", "/.ssh/id_ed25519", "/.ssh/id_ecdsa", "/.ssh/id_dsa",
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "/etc/ssh/ssh_host_") && strings.HasSuffix(v, "_key") {
			return zc1613Hit(cmd, v)
		}
		for _, s := range userSuffixes {
			if strings.HasSuffix(v, s) {
				return zc1613Hit(cmd, v)
			}
		}
	}
	return nil
}

func zc1613Hit(cmd *ast.SimpleCommand, path string) []Violation {
	return []Violation{{
		KataID: "ZC1613",
		Message: "Reading `" + path + "` through a text tool copies private-key material " +
			"into the process and often into logs / scrollback. Use `ssh-keygen -l -f` " +
			"for metadata, or pass the path directly to the consumer.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1614",
		Title:    "Error on `expect` script containing `password` / `passphrase`",
		Severity: SeverityError,
		Description: "`expect -c '... password ... send \"...\"'` puts the entire scripted " +
			"dialog on the command line. Anything there — including the password or passphrase " +
			"— is visible in `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs. Use " +
			"key-based authentication (SSH keys, GSSAPI) where possible. If password feeding is " +
			"truly unavoidable, read it from a protected file with `spawn -o`, or source it " +
			"from an environment variable the script does not print.",
		Check: checkZC1614,
	})
}

func checkZC1614(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "expect" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		low := strings.ToLower(arg.String())
		if strings.Contains(low, "password") || strings.Contains(low, "passphrase") {
			return []Violation{{
				KataID: "ZC1614",
				Message: "`expect` script contains `password` / `passphrase` — the full " +
					"argv lands in `ps` and audit logs. Switch to key-based auth, or " +
					"read the credential from a protected file the expect script opens.",
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
		ID:       "ZC1615",
		Title:    "Style: use Zsh `$EPOCHREALTIME` / `$epochtime` instead of `date \"+%s.%N\"`",
		Severity: SeverityStyle,
		Description: "Zsh's `zsh/datetime` module exposes `$EPOCHREALTIME` (scalar with " +
			"fractional seconds) and `$epochtime` (two-element array of seconds and " +
			"nanoseconds). Both read straight from `clock_gettime(CLOCK_REALTIME)` without " +
			"forking `date`. On a hot path the builtin is dramatically faster and avoids " +
			"subshell process-startup overhead. Autoload the module once with `zmodload " +
			"zsh/datetime`.",
		Check: checkZC1615,
	})
}

func checkZC1615(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "date" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if !strings.Contains(v, "%s") {
			continue
		}
		if strings.Contains(v, "%N") {
			return []Violation{{
				KataID: "ZC1615",
				Message: "`date " + v + "` forks for sub-second time. Use Zsh " +
					"`$EPOCHREALTIME` / `$epochtime` from `zmodload zsh/datetime`.",
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
		ID:       "ZC1616",
		Title:    "Warn on `fsfreeze -f MOUNTPOINT` — filesystem stays frozen until `-u` runs",
		Severity: SeverityWarning,
		Description: "`fsfreeze -f` blocks every write on the mountpoint until `fsfreeze -u` " +
			"thaws it. The intended use is a short window around a hypervisor or LVM snapshot. " +
			"If the script errors between the freeze and the unfreeze (or is killed), the " +
			"filesystem stays frozen — every subsequent write hangs forever until the admin " +
			"manually thaws it, and a reboot may be the only way out on the root fs. Pair " +
			"every freeze with `trap 'fsfreeze -u MOUNTPOINT' EXIT` and keep the window under " +
			"a few seconds.",
		Check: checkZC1616,
	})
}

func checkZC1616(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "fsfreeze" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-f" {
			return []Violation{{
				KataID: "ZC1616",
				Message: "`fsfreeze -f` freezes the mountpoint — every write hangs until " +
					"`fsfreeze -u` runs. Wrap the call in `trap 'fsfreeze -u PATH' EXIT` " +
					"so the thaw fires even on failure.",
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
		ID:       "ZC1617",
		Title:    "Warn on `xargs -P 0` — unbounded parallelism risks CPU / fd / memory exhaustion",
		Severity: SeverityWarning,
		Description: "`xargs -P 0` tells xargs to spawn as many concurrent children as input " +
			"lines. On any non-trivial input that number can blow past `RLIMIT_NPROC`, " +
			"saturate the downstream tool's file-descriptor limit, or drive the host OOM. " +
			"Pick an explicit cap — `xargs -P $(nproc)` for CPU-bound work, `-P 4..8` for " +
			"I/O-bound — so the failure mode is bounded and predictable.",
		Check: checkZC1617,
	})
}

func checkZC1617(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "xargs" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-P0" {
			return zc1617Hit(cmd)
		}
		if v == "-P" && i+1 < len(cmd.Arguments) && cmd.Arguments[i+1].String() == "0" {
			return zc1617Hit(cmd)
		}
	}
	return nil
}

func zc1617Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1617",
		Message: "`xargs -P 0` spawns one child per input line — CPU / FD / memory " +
			"exhaustion risk. Use `-P $(nproc)` or an explicit cap.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1618",
		Title:    "Warn on `git commit --no-verify` / `git push --no-verify` — bypasses hooks",
		Severity: SeverityWarning,
		Description: "`--no-verify` skips pre-commit, commit-msg, and pre-push hooks. Those " +
			"hooks are where projects run linting, type-checking, unit tests, and secret " +
			"scanning before code lands. A commit or push with `--no-verify` ships code the " +
			"project's own automation would have rejected. Reserve the flag for emergencies " +
			"with a follow-up commit that passes the hooks; scripts should not use it " +
			"routinely.",
		Check: checkZC1618,
	})
}

func checkZC1618(node ast.Node) []Violation {
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
	if len(cmd.Arguments) < 2 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "commit" && sub != "push" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--no-verify" || (sub == "commit" && v == "-n") {
			return []Violation{{
				KataID: "ZC1618",
				Message: "`git " + sub + " " + v + "` skips pre-" + sub + " / commit-msg " +
					"hooks — lint, test, and secret-scan checks do not run. Reserve for " +
					"emergencies; scripts should pass the hooks.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

var zc1619NetworkFS = map[string]bool{
	"nfs": true, "nfs4": true,
	"cifs": true, "smbfs": true, "smb3": true,
	"sshfs": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1619",
		Title:    "Warn on `mount -t nfs/cifs/smb/sshfs` missing `nosuid` or `nodev`",
		Severity: SeverityWarning,
		Description: "Network filesystems present files whose mode bits are controlled by a " +
			"remote server. Without `nosuid` in the mount options, a compromised or hostile " +
			"server can plant a setuid-root binary on the share; the client kernel honors the " +
			"suid bit and the binary runs as root on the mounting host. Without `nodev`, the " +
			"server can plant device nodes the kernel treats as real. Always mount network " +
			"shares with `nosuid,nodev`; add `noexec` unless the export is intended to hold " +
			"executables.",
		Check: checkZC1619,
	})
}

func checkZC1619(node ast.Node) []Violation {
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

	var fsType, opts string
	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-t" && i+1 < len(cmd.Arguments) {
			fsType = cmd.Arguments[i+1].String()
		}
		if v == "-o" && i+1 < len(cmd.Arguments) {
			opts = cmd.Arguments[i+1].String()
		}
	}

	if !zc1619NetworkFS[fsType] {
		return nil
	}
	if strings.Contains(opts, "nosuid") && strings.Contains(opts, "nodev") {
		return nil
	}

	missing := []string{}
	if !strings.Contains(opts, "nosuid") {
		missing = append(missing, "nosuid")
	}
	if !strings.Contains(opts, "nodev") {
		missing = append(missing, "nodev")
	}
	return []Violation{{
		KataID: "ZC1619",
		Message: "`mount -t " + fsType + "` without " + strings.Join(missing, ",") +
			" — a hostile server can plant setuid binaries or device nodes that the " +
			"client kernel honors. Add `nosuid,nodev` to the `-o` options.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1620",
		Title:    "Error on `tee /etc/sudoers` / `/etc/sudoers.d/*` — writes without `visudo -cf`",
		Severity: SeverityError,
		Description: "`tee` copies stdin to the file with no syntax check. A typo in a sudoers " +
			"rule — a stray comma, a missing `ALL`, an unclosed alias — leaves the file " +
			"unparseable. The next sudo call refuses to load it and on most systems nobody " +
			"can become root until someone boots from rescue media. Pipe the content through " +
			"`visudo -cf /dev/stdin` first, or write to a temp file, validate with " +
			"`visudo -cf`, then atomically `mv` into `/etc/sudoers.d/`.",
		Check: checkZC1620,
	})
}

func checkZC1620(node ast.Node) []Violation {
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
		if v == "/etc/sudoers" || strings.HasPrefix(v, "/etc/sudoers.d/") {
			return []Violation{{
				KataID: "ZC1620",
				Message: "`tee " + v + "` writes without syntax validation — a typo locks " +
					"everyone out of sudo. Pipe through `visudo -cf /dev/stdin` or stage " +
					"in a temp file and `visudo -cf` before `mv`.",
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
		ID:       "ZC1621",
		Title:    "Warn on `tmux -S /tmp/SOCKET` — shared-path socket invites session hijack",
		Severity: SeverityWarning,
		Description: "`tmux -S PATH` overrides the default socket location (normally under " +
			"`$XDG_RUNTIME_DIR/tmux-$UID/`, a 0700-mode directory). Paths under `/tmp/` or " +
			"`/var/tmp/` are world-traversable; if the socket is created with loose " +
			"permissions, any local user who can read it can `tmux -S /tmp/PATH attach` and " +
			"see / drive the session — keystrokes, output, arbitrary commands in the attached " +
			"pane. Keep the socket in `$XDG_RUNTIME_DIR` or another 0700-scoped directory.",
		Check: checkZC1621,
	})
}

func checkZC1621(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tmux" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() != "-S" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		path := cmd.Arguments[i+1].String()
		if strings.HasPrefix(path, "/tmp/") || strings.HasPrefix(path, "/var/tmp/") {
			return []Violation{{
				KataID: "ZC1621",
				Message: "`tmux -S " + path + "` places the socket in a world-traversable " +
					"directory — any local user who can read the socket can attach the " +
					"session. Use `$XDG_RUNTIME_DIR` or a 0700-scoped parent dir.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

var zc1622BashOps = []string{"@U}", "@L}", "@Q}", "@E}", "@A}", "@K}", "@a}"}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1622",
		Title:    "Style: `${var@U/L/Q/...}` — prefer Zsh `${(U)var}` / `${(L)var}` / `${(Q)var}` flags",
		Severity: SeverityStyle,
		Description: "The `@<op>` suffix came from Bash 5. Zsh 5.9+ compiles in compatibility " +
			"for the common ones, but the idiomatic Zsh form is the `(X)var` parameter-" +
			"expansion flag — `${(U)var}` uppercase, `${(L)var}` lowercase, `${(Q)var}` " +
			"unquote, `${(k)var}` keys, `${(t)var}` type, `${(e)var}` re-evaluate. The flag " +
			"form composes (`${(Uf)str}` works) and reads consistently across the Zsh " +
			"documentation. Prefer the native flag over the Bash-compat form.",
		Check: checkZC1622,
	})
}

func checkZC1622(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if !strings.Contains(v, "${") {
			continue
		}
		for _, op := range zc1622BashOps {
			if strings.Contains(v, op) {
				return []Violation{{
					KataID: "ZC1622",
					Message: "`${var" + strings.TrimSuffix(op, "}") + "}` — prefer Zsh " +
						"`${(X)var}` parameter-expansion flags (e.g. `${(U)var}` for " +
						"uppercase).",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityStyle,
				}}
			}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1623",
		Title:    "Warn on `kill -STOP PID` / `pkill -STOP` — target halts until `kill -CONT` runs",
		Severity: SeverityWarning,
		Description: "Sending SIGSTOP halts the target process until SIGCONT arrives. If the " +
			"script fails, is killed, or exits before the resume, the target stays paused " +
			"indefinitely — consuming memory, holding locks, blocking its dependents. Wrap " +
			"every `kill -STOP $PID` with `trap \"kill -CONT $PID\" EXIT` (or an explicit " +
			"cleanup path) so the resume fires even on failure. Prefer `kill -TSTP` if the " +
			"target can handle it (the user-space tstop that the process can ignore).",
		Check: checkZC1623,
	})
}

func checkZC1623(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kill" && ident.Value != "pkill" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-STOP" || v == "-SIGSTOP" || v == "-19" {
			return zc1623Hit(cmd)
		}
		if v == "-s" && i+1 < len(cmd.Arguments) {
			sig := cmd.Arguments[i+1].String()
			if sig == "STOP" || sig == "SIGSTOP" || sig == "19" {
				return zc1623Hit(cmd)
			}
		}
	}
	return nil
}

func zc1623Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1623",
		Message: "`kill -STOP` halts the target until SIGCONT arrives. Pair every STOP " +
			"with `trap \"kill -CONT PID\" EXIT` so the resume fires even on failure.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1624",
		Title:    "Error on `az login -p` / `--password` — service-principal secret in process list",
		Severity: SeverityError,
		Description: "`az login -p SECRET` passes the service-principal password as an argv " +
			"element. The expanded value shows up in `ps`, `/proc/<pid>/cmdline`, shell " +
			"history, and audit logs — readable by any local user who can list processes. " +
			"Prefer federated-token OIDC (`--federated-token`), managed identity on the host, " +
			"or interactive device-code flow. If a password is unavoidable, export it as " +
			"`AZURE_PASSWORD` via a protected env var and call plain `az login --service-" +
			"principal` (which reads from env).",
		Check: checkZC1624,
	})
}

func checkZC1624(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "az" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "login" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-p" || v == "--password" {
			return []Violation{{
				KataID: "ZC1624",
				Message: "`az login " + v + "` puts the SP password in argv — visible in " +
					"`ps` / `/proc/<pid>/cmdline`. Use federated-token OIDC, managed " +
					"identity, or `AZURE_PASSWORD` via a protected env var.",
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
		ID:       "ZC1625",
		Title:    "Error on `rm --no-preserve-root` — disables GNU rm safeguard against `rm -rf /`",
		Severity: SeverityError,
		Description: "GNU `rm` refuses to remove `/` by default — the `--preserve-root` " +
			"safeguard added in coreutils 8.4. `--no-preserve-root` explicitly disables that " +
			"check so `rm -rf /` actually recurses and wipes the filesystem. Scripts that pass " +
			"the flag are asking `rm` to go ahead if the argument happens to evaluate to `/`. " +
			"Remove the flag; if a specific path genuinely needs deletion, list it explicitly " +
			"and leave the safeguard in place.",
		Check: checkZC1625,
	})
}

func checkZC1625(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "rm" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "--no-preserve-root" {
			return []Violation{{
				KataID: "ZC1625",
				Message: "`rm --no-preserve-root` disables the GNU safeguard against `rm -rf " +
					"/`. Remove the flag; if a specific path needs deletion, list it " +
					"explicitly.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

var zc1626SecretKeys = []string{
	"password", "passwd", "secret", "token",
	"apikey", "api_key", "api-key",
	"accesskey", "access_key", "access-key",
	"privatekey", "private_key", "private-key",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1626",
		Title:    "Error on `helm install/upgrade --set KEY=VALUE` with secret-shaped key",
		Severity: SeverityError,
		Description: "`--set` and `--set-string` put the full `KEY=VALUE` pair on the helm " +
			"command line. When the key name looks like a secret (`password`, `secret`, " +
			"`token`, `apikey`, `access_key`, `private_key`), the expanded VALUE appears in " +
			"`ps`, `/proc/<pid>/cmdline`, shell history, and audit logs — readable by any " +
			"local user who can list processes. Put secrets in a protected values file " +
			"(`helm install -f /secure/values.yaml`), or use `--set-file KEY=PATH` so helm " +
			"reads the content from PATH at apply time.",
		Check: checkZC1626,
	})
}

func checkZC1626(node ast.Node) []Violation {
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
	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "install" && sub != "upgrade" && sub != "template" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v != "--set" && v != "--set-string" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		pair := cmd.Arguments[i+1].String()
		eq := strings.Index(pair, "=")
		if eq < 0 {
			continue
		}
		key := strings.ToLower(pair[:eq])
		for _, s := range zc1626SecretKeys {
			if strings.Contains(key, s) {
				return []Violation{{
					KataID: "ZC1626",
					Message: "`helm " + sub + " " + v + " " + pair + "` places a secret " +
						"value in argv — readable via `ps`. Use `-f values.yaml` or " +
						"`--set-file " + key + "=PATH`.",
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
		ID:       "ZC1627",
		Title:    "Warn on `crontab /tmp/FILE` — attacker-writable path installed as a user's cron",
		Severity: SeverityWarning,
		Description: "`crontab PATH` replaces the user's cron with whatever PATH currently " +
			"contains. A path under `/tmp/` or `/var/tmp/` is world-traversable; a concurrent " +
			"local user can replace the file between the moment the script writes it and the " +
			"moment `crontab` reads it, substituting their own cron rules. Keep the staging " +
			"file in a 0700-scoped directory (e.g. `$XDG_RUNTIME_DIR/` or `mktemp -d`), or " +
			"pipe the content via `crontab -` after generating it in-memory.",
		Check: checkZC1627,
	})
}

func checkZC1627(node ast.Node) []Violation {
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

	var skipNext bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if skipNext {
			skipNext = false
			continue
		}
		if v == "-u" || v == "-s" {
			skipNext = true
			continue
		}
		if strings.HasPrefix(v, "-") {
			continue
		}
		if strings.HasPrefix(v, "/tmp/") || strings.HasPrefix(v, "/var/tmp/") {
			return []Violation{{
				KataID: "ZC1627",
				Message: "`crontab " + v + "` reads cron rules from a world-traversable " +
					"path — a concurrent local user can substitute the file between write " +
					"and read. Stage the file in `$XDG_RUNTIME_DIR/` or `mktemp -d`, or " +
					"pipe via `crontab -`.",
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
		ID:       "ZC1628",
		Title:    "Warn on `insmod` / `modprobe -f` — loads modules bypassing blacklist / signature checks",
		Severity: SeverityWarning,
		Description: "`insmod PATH.ko` loads a kernel module from a file, skipping the depmod-" +
			"built dependency graph and the `/etc/modprobe.d/*.conf` blacklist. `modprobe " +
			"-f` instructs modprobe to ignore version-magic and kernel-mismatch checks. " +
			"Either path lets a module enter the kernel that the administrator explicitly " +
			"disabled, or one compiled against a different kernel — crash, privesc, or full " +
			"kernel compromise. Use plain `modprobe MODNAME` so the system's policy and " +
			"signature verification run.",
		Check: checkZC1628,
	})
}

func checkZC1628(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value == "insmod" {
		if len(cmd.Arguments) == 0 {
			return nil
		}
		return []Violation{{
			KataID: "ZC1628",
			Message: "`insmod` loads a kernel module bypassing depmod / blacklist — prefer " +
				"`modprobe MODNAME` so system policy and signature checks apply.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	if ident.Value == "modprobe" {
		for _, arg := range cmd.Arguments {
			if arg.String() == "-f" {
				return []Violation{{
					KataID: "ZC1628",
					Message: "`modprobe -f` ignores version-magic and kernel-mismatch " +
						"checks — a mismatched module can crash or compromise the kernel. " +
						"Drop the flag and fix the underlying version mismatch.",
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
		ID:       "ZC1629",
		Title:    "Warn on `rsync --rsync-path='sudo rsync'` — hidden remote privilege escalation",
		Severity: SeverityWarning,
		Description: "`--rsync-path` normally overrides the path to the remote rsync binary. " +
			"Setting it to `sudo rsync` (or `doas rsync` / `pkexec rsync`) instead makes the " +
			"remote side run rsync as root. That is sometimes legitimate — copying into " +
			"`/etc/` from a CI job — but the flag is easy to miss in review because it looks " +
			"like a path override. Provision a scoped sudoers rule that names exactly which " +
			"rsync invocation the remote user may run, and keep the path explicit (`--rsync-" +
			"path=/usr/bin/rsync`).",
		Check: checkZC1629,
	})
}

func checkZC1629(node ast.Node) []Violation {
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

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if !strings.HasPrefix(v, "--rsync-path=") {
			continue
		}
		val := strings.TrimPrefix(v, "--rsync-path=")
		val = strings.Trim(val, "\"'")
		if strings.Contains(val, "sudo") ||
			strings.Contains(val, "doas") ||
			strings.Contains(val, "pkexec") {
			return []Violation{{
				KataID: "ZC1629",
				Message: "`rsync --rsync-path='" + val + "'` runs remote rsync under " +
					"privilege escalation. Use a scoped sudoers rule on the remote host " +
					"and keep the path explicit (`/usr/bin/rsync`).",
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
		ID:       "ZC1630",
		Title:    "Warn on `php -S 0.0.0.0:PORT` — PHP dev server exposes CWD to all interfaces",
		Severity: SeverityWarning,
		Description: "`php -S 0.0.0.0:PORT` starts PHP's built-in dev server listening on every " +
			"interface the host has. It serves files from the working directory (or the " +
			"docroot named after the bind) with no auth, no TLS, and minimal access logging. " +
			"The PHP docs explicitly say not to use it in production. Bind to `127.0.0.1:PORT` " +
			"for local testing and put nginx / caddy in front for anything externally exposed.",
		Check: checkZC1630,
	})
}

func checkZC1630(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "php" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() != "-S" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			return nil
		}
		bind := cmd.Arguments[i+1].String()
		if strings.HasPrefix(bind, "0.0.0.0:") ||
			strings.HasPrefix(bind, "*:") ||
			strings.HasPrefix(bind, "[::]:") {
			return []Violation{{
				KataID: "ZC1630",
				Message: "`php -S " + bind + "` binds the dev server to every interface — " +
					"unauthenticated access to the working directory. Use `127.0.0.1:PORT` " +
					"locally, nginx / caddy for external exposure.",
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
		ID:       "ZC1631",
		Title:    "Error on `openssl ... -passin pass:SECRET` / `-passout pass:SECRET`",
		Severity: SeverityError,
		Description: "OpenSSL's `-passin` / `-passout` accept a password source selector. The " +
			"`pass:LITERAL` form embeds the password as an argv element — visible in `ps`, " +
			"`/proc/<pid>/cmdline`, shell history, and audit logs. Use one of the safer " +
			"sources: `env:VARNAME` reads from an env var, `file:PATH` reads the first line " +
			"of PATH, `fd:N` reads from an open descriptor, `stdin` reads a line from stdin.",
		Check: checkZC1631,
	})
}

func checkZC1631(node ast.Node) []Violation {
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

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v != "-passin" && v != "-passout" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		val := cmd.Arguments[i+1].String()
		if !strings.HasPrefix(val, "pass:") {
			continue
		}
		return []Violation{{
			KataID: "ZC1631",
			Message: "`openssl " + v + " " + val + "` puts the password in argv — visible " +
				"via `ps`. Use `env:VARNAME`, `file:PATH`, `fd:N`, or `stdin`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1632",
		Title:    "Warn on `shred` — unreliable on journaled / CoW filesystems (ext4, btrfs, zfs)",
		Severity: SeverityWarning,
		Description: "`shred` assumes in-place overwrites, which is how ext2 worked. On a " +
			"journaled ext4 the overwrite passes go through the journal and may not hit the " +
			"original data blocks. On CoW filesystems (btrfs, zfs, xfs with reflink) the " +
			"overwrite lands in fresh blocks and leaves the old content intact until garbage " +
			"collection decides otherwise. `shred`'s own man page warns about this. For modern " +
			"secure deletion, use full-disk encryption with key destruction, or retire the " +
			"device with `blkdiscard` on an SSD.",
		Check: checkZC1632,
	})
}

func checkZC1632(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "shred" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1632",
		Message: "`shred` may not overwrite original blocks on ext4/btrfs/zfs. For " +
			"guaranteed erasure, use full-disk encryption with key destruction, or " +
			"`blkdiscard` when retiring an SSD.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1633",
		Title:    "Error on `gpg --passphrase SECRET` — passphrase on cmdline",
		Severity: SeverityError,
		Description: "`gpg --passphrase VALUE` passes the key passphrase as an argv element. " +
			"Visible in `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs for every " +
			"local user who can list processes. Use `--passphrase-file PATH` (reads the first " +
			"line of PATH), `--passphrase-fd N` (reads from file descriptor N), or " +
			"`--pinentry-mode=loopback` with the passphrase piped on stdin. Pair with " +
			"`--batch` for non-interactive runs.",
		Check: checkZC1633,
	})
}

func checkZC1633(node ast.Node) []Violation {
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
		if arg.String() == "--passphrase" {
			return []Violation{{
				KataID: "ZC1633",
				Message: "`gpg --passphrase` puts the passphrase in argv — visible via `ps`. " +
					"Use `--passphrase-file`, `--passphrase-fd`, or `--pinentry-" +
					"mode=loopback` with the value on stdin.",
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
		ID:       "ZC1634",
		Title:    "Warn on `umask NNN` that fails to mask world-write — mask-inversion footgun",
		Severity: SeverityWarning,
		Description: "`umask` is a mask: bits that are set are removed from the default " +
			"permission. The classic pitfall is reading it as \"permissions I want\" — " +
			"`umask 111` feels tight (\"no execute for anyone\") but it does not mask the write " +
			"bit, so every new file is `666` (rw-rw-rw-). The \"other\" digit must be one of " +
			"`2/3/6/7` to strip world-write. Use `022` for publicly readable files, `077` for " +
			"secrets-handling.",
		Check: checkZC1634,
	})
}

func checkZC1634(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || CommandIdentifier(cmd) != "umask" || len(cmd.Arguments) != 1 {
		return nil
	}
	v := cmd.Arguments[0].String()
	if !zc1634UmaskMissesWorldWrite(v) {
		return nil
	}
	return []Violation{{
		KataID: "ZC1634",
		Message: "`umask " + v + "` leaves world-write on new files — the \"other\" digit " +
			"must be `2`/`3`/`6`/`7` to mask the write bit. Use `022` for public, `077` " +
			"for secrets.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

var zc1634AllZero = map[string]struct{}{"0": {}, "00": {}, "000": {}, "0000": {}}

func zc1634UmaskMissesWorldWrite(v string) bool {
	if _, hit := zc1634AllZero[v]; hit {
		return false
	}
	if len(v) < 3 || len(v) > 4 {
		return false
	}
	for _, c := range v {
		if c < '0' || c > '7' {
			return false
		}
	}
	switch v[len(v)-1] {
	case '0', '1', '4', '5':
		return true
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1635",
		Title:    "Error on `mysql -pSECRET` / `--password=SECRET` — password in process list",
		Severity: SeverityError,
		Description: "MySQL / MariaDB clients accept the password concatenated with the `-p` " +
			"flag (`-pSECRET`) or via `--password=SECRET`. Both forms put the secret in argv " +
			"— visible in `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs for every " +
			"local user who can list processes. Use `-p` with no argument for an interactive " +
			"prompt, `--login-path` for the credentials helper file, or a `~/.my.cnf` with " +
			"`0600` perms and `[client] password=...` so the client reads it at startup.",
		Check: checkZC1635,
	})
}

func checkZC1635(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "mysql", "mysqldump", "mysqladmin", "mariadb", "mariadb-dump", "mariadb-admin":
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "-p") && len(v) > 2 {
			return []Violation{{
				KataID: "ZC1635",
				Message: "`" + ident.Value + " " + v + "` puts the MySQL password in argv. " +
					"Use `-p` with no arg (prompt), `--login-path`, or a 0600 `~/.my.cnf`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
		if strings.HasPrefix(v, "--password=") {
			return []Violation{{
				KataID: "ZC1635",
				Message: "`" + ident.Value + " " + v + "` puts the MySQL password in argv. " +
					"Use `-p` with no arg (prompt), `--login-path`, or a 0600 `~/.my.cnf`.",
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
		ID:       "ZC1636",
		Title:    "Warn on `virsh destroy DOMAIN` — force-stops VM (no graceful shutdown)",
		Severity: SeverityWarning,
		Description: "`virsh destroy DOM` is the libvirt equivalent of pulling the plug on a " +
			"running VM. The guest OS gets no chance to flush filesystems, close network " +
			"connections, or run its own shutdown services — data corruption risk on any " +
			"open file in the guest. For graceful shutdown use `virsh shutdown DOM` (ACPI " +
			"event), wait for completion, and only fall back to `destroy` for a genuinely " +
			"unresponsive guest. `virsh destroy --graceful DOM` attempts a timed graceful " +
			"first, then forces — that variant is not flagged.",
		Check: checkZC1636,
	})
}

func checkZC1636(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "virsh" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "destroy" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "--graceful" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1636",
		Message: "`virsh destroy` yanks power from the VM — filesystem corruption risk. Use " +
			"`virsh shutdown` for graceful stop, or `virsh destroy --graceful` as a timed " +
			"fallback.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1637",
		Title:    "Style: prefer Zsh `typeset -r NAME=value` over POSIX `readonly NAME=value`",
		Severity: SeverityStyle,
		Description: "Both `readonly NAME` and `typeset -r NAME` create a read-only parameter. " +
			"In Zsh the idiomatic form is `typeset -r` — it composes with other typeset flags " +
			"(`-ir` for readonly integer, `-xr` for readonly export, `-gr` to pin a readonly " +
			"global from inside a function). `readonly` works but reads as a Bash / POSIX-ism " +
			"in a Zsh codebase.",
		Check: checkZC1637,
		Fix:   fixZC1637,
	})
}

// fixZC1637 rewrites the `readonly` command name to `typeset -r`.
// Single-edit replacement at the violation column. Detector gates on
// the bare command name match, so the rewrite is idempotent on a
// re-run (the new line starts with `typeset` not `readonly`).
func fixZC1637(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "readonly" {
		return nil
	}
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 || off+len("readonly") > len(source) {
		return nil
	}
	if string(source[off:off+len("readonly")]) != "readonly" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("readonly"),
		Replace: "typeset -r",
	}}
}

func checkZC1637(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "readonly" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1637",
		Message: "`readonly` works but `typeset -r NAME=value` is the Zsh-native form and " +
			"composes with other typeset flags.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

var zc1638SecretArgs = []string{
	"password", "passwd", "secret", "token",
	"apikey", "api_key", "api-key",
	"accesskey", "access_key", "access-key",
	"privatekey", "private_key", "private-key",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1638",
		Title:    "Error on `docker/podman build --build-arg SECRET=VALUE` — secret baked into image layer",
		Severity: SeverityError,
		Description: "`--build-arg KEY=VALUE` values land in the image metadata that `docker " +
			"history` (and the analogous podman / buildah tooling) read back from the layer. " +
			"Even if the Dockerfile only uses the arg to export as a build-time env var, the " +
			"literal value is cached in the layer forever. A key-shaped name (`password`, " +
			"`secret`, `token`, `apikey`, `access_key`, `private_key`) with a concrete value " +
			"embeds that secret in every image pulled. Use BuildKit secrets " +
			"(`--secret id=mysecret,src=path`) or a multi-stage build where the secret stays " +
			"in a discarded stage.",
		Check: checkZC1638,
	})
}

func checkZC1638(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "docker", "podman", "buildah":
	default:
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "build" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() != "--build-arg" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		pair := cmd.Arguments[i+1].String()
		eq := strings.Index(pair, "=")
		if eq < 0 {
			continue
		}
		key := strings.ToLower(pair[:eq])
		for _, s := range zc1638SecretArgs {
			if strings.Contains(key, s) {
				return []Violation{{
					KataID: "ZC1638",
					Message: "`" + ident.Value + " build --build-arg " + pair + "` bakes " +
						"the secret into the image layer metadata. Use `--secret " +
						"id=NAME,src=PATH` (BuildKit) or a multi-stage build.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
	}
	return nil
}

var zc1639AuthHeaders = []string{
	"authorization:",
	"proxy-authorization:",
	"x-api-key:",
	"api-key:",
	"apikey:",
	"x-auth-token:",
	"x-access-token:",
	"cookie:",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1639",
		Title:    "Error on `curl -H 'Authorization: ...'` — credential header in process list",
		Severity: SeverityError,
		Description: "`-H \"Authorization: Bearer $TOKEN\"` (and similar credential-bearing " +
			"headers like `X-Api-Key`, `X-Auth-Token`, `Proxy-Authorization`, `Cookie`) put " +
			"the expanded value in argv. It shows up in `ps`, `/proc/<pid>/cmdline`, shell " +
			"history, and audit logs — every local user who can list processes reads the " +
			"secret. Pass the header via a file with `-H @FILE` or use `--config FILE` so the " +
			"value stays on disk (with 0600 perms), never on the command line.",
		Check: checkZC1639,
	})
}

func checkZC1639(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "curl" && ident.Value != "wget" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v != "-H" && v != "--header" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		header := strings.ToLower(cmd.Arguments[i+1].String())
		for _, h := range zc1639AuthHeaders {
			if strings.Contains(header, h) {
				return []Violation{{
					KataID: "ZC1639",
					Message: "`" + ident.Value + " -H " + cmd.Arguments[i+1].String() +
						"` places the credential in argv — visible via `ps`. Use `-H @FILE`" +
						" or `--config FILE` with 0600 perms.",
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
		ID:       "ZC1640",
		Title:    "Style: `${!var}` Bash indirect expansion — prefer Zsh `${(P)var}`",
		Severity: SeverityStyle,
		Description: "`${!var}` is Bash indirect expansion — it reads the value of the " +
			"parameter whose name is stored in `$var`. Zsh has the native flag form " +
			"`${(P)var}` which does the same and composes with other parameter-expansion " +
			"flags (`${(Pf)var}` to split the indirect value on newlines, for example). " +
			"`${!prefix*}` / `${!array[@]}` have Zsh equivalents via the `$parameters` hash " +
			"or `(k)` subscript flags. Prefer the native Zsh form in a Zsh codebase.",
		Check: checkZC1640,
	})
}

func checkZC1640(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "${!") {
			return []Violation{{
				KataID: "ZC1640",
				Message: "`${!var}` Bash indirect — prefer Zsh `${(P)var}` for the same " +
					"semantics with flag composability.",
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
		ID:       "ZC1641",
		Title:    "Error on `kubectl create secret --from-literal=...` / `--docker-password=...`",
		Severity: SeverityError,
		Description: "`kubectl create secret generic --from-literal=KEY=VALUE` and " +
			"`kubectl create secret docker-registry --docker-password=VALUE` put the secret " +
			"content in argv. The expanded value shows up in `ps`, `/proc/<pid>/cmdline`, " +
			"shell history, and audit logs — readable by any local user who can list " +
			"processes. Use `--from-file=KEY=PATH` (reads from a 0600-protected file), " +
			"`--from-env-file=PATH` (reads KEY=VALUE lines), or pipe a manifest into " +
			"`kubectl apply -f -` with base64-encoded `data:` values staged on disk.",
		Check: checkZC1641,
	})
}

func checkZC1641(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "create" || cmd.Arguments[1].String() != "secret" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if strings.HasPrefix(v, "--from-literal=") || strings.HasPrefix(v, "--docker-password=") {
			return []Violation{{
				KataID: "ZC1641",
				Message: "`kubectl create secret " + v + "` puts the secret in argv — " +
					"visible via `ps`. Use `--from-file=KEY=PATH` / `--from-env-file=PATH`.",
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
		ID:       "ZC1642",
		Title:    "Warn on `tshark -w FILE` / `dumpcap -w FILE` without `-Z user` — capture file owned by root",
		Severity: SeverityWarning,
		Description: "Packet captures routinely need `CAP_NET_RAW`, so the capture process " +
			"typically runs as root. Without `-Z USER` the resulting pcap is root-owned — a " +
			"subsequent analyst who opens it with Wireshark (which can run helper scripts from " +
			"the file) operates on a root-owned file and may unintentionally invoke things as " +
			"root. `-Z USER` tells `tshark` / `dumpcap` to drop privileges for the actual " +
			"capture and write the file as `USER`.",
		Check: checkZC1642,
	})
}

func checkZC1642(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tshark" && ident.Value != "dumpcap" {
		return nil
	}

	var hasWrite, hasDrop bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-w" {
			hasWrite = true
		}
		if v == "-Z" {
			hasDrop = true
		}
	}
	if !hasWrite || hasDrop {
		return nil
	}

	return []Violation{{
		KataID: "ZC1642",
		Message: "`" + ident.Value + " -w FILE` without `-Z USER` leaves the pcap root-" +
			"owned. Add `-Z USER` to drop privileges for the capture.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1643",
		Title:    "Style: `$(cat file)` — use `$(<file)` to skip the fork / exec",
		Severity: SeverityStyle,
		Description: "`$(cat FILE)` forks, execs `/usr/bin/cat`, reads FILE, writes the bytes " +
			"to the pipe, waits for the child. `$(<FILE)` is a shell builtin — it reads FILE " +
			"directly into the command-substitution buffer with no fork and no exec. In a hot " +
			"path the speedup is dramatic, and even in cold paths it avoids one of the most " +
			"common useless-use-of-cat patterns in review feedback.",
		Check: checkZC1643,
		Fix:   fixZC1643,
	})
}

// fixZC1643 rewrites `$(cat FILE)` into `$(<FILE)`. The detector
// matches on the literal `$(cat ` substring inside a command argument;
// the Fix scans each argument's source span for that prefix and
// replaces `cat ` with `<` (4-byte → 1-byte collapse). Each occurrence
// becomes its own FixEdit, so multiple `$(cat …)` substitutions on the
// same line are all handled in one pass.
//
// Idempotent because the literal `$(cat ` no longer appears in the
// rewritten text — the second-pass detector sees `$(<…)` and stays
// silent.
func fixZC1643(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	const needle = "$(cat "
	var edits []FixEdit
	for _, arg := range cmd.Arguments {
		tok := arg.TokenLiteralNode()
		argOff := LineColToByteOffset(source, tok.Line, tok.Column)
		if argOff < 0 {
			continue
		}
		// Constrain the search window to the argument's literal text
		// span. arg.String() may include AST decorations (parens for
		// expressions) so use the byte slice directly: walk forward
		// from the start of the argument until a delimiter that ends
		// the unquoted argument context.
		end := argOff
		// Track quoting so the search does not stop inside a quoted
		// substring. The argument span runs to the next unquoted
		// whitespace, `;`, `&`, `|`, `)` (when paren depth = 0), or
		// newline.
		inSingle := false
		inDouble := false
		parenDepth := 0
		for end < len(source) {
			c := source[end]
			if c == '\\' && end+1 < len(source) {
				end += 2
				continue
			}
			switch {
			case inSingle:
				if c == '\'' {
					inSingle = false
				}
			case inDouble:
				if c == '"' {
					inDouble = false
				}
			default:
				switch c {
				case '\'':
					inSingle = true
				case '"':
					inDouble = true
				case '(':
					parenDepth++
				case ')':
					if parenDepth == 0 {
						goto done
					}
					parenDepth--
				case ' ', '\t', '\n', ';', '&', '|':
					if parenDepth == 0 {
						goto done
					}
				}
			}
			end++
		}
	done:
		// Find every `$(cat ` occurrence inside [argOff, end).
		i := argOff
		for i+len(needle) <= end {
			if string(source[i:i+len(needle)]) != needle {
				i++
				continue
			}
			// Replace the `cat ` portion (4 bytes after `$(`) with `<`.
			line, col := offsetLineColZC1643(source, i+2)
			if line < 0 {
				i += len(needle)
				continue
			}
			edits = append(edits, FixEdit{
				Line:    line,
				Column:  col,
				Length:  4,
				Replace: "<",
			})
			i += len(needle)
		}
	}
	return edits
}

func offsetLineColZC1643(source []byte, offset int) (int, int) {
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

func checkZC1643(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "$(cat ") {
			return []Violation{{
				KataID: "ZC1643",
				Message: "`$(cat FILE)` forks cat just to read a file — use `$(<FILE)` " +
					"(shell builtin, no fork).",
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
		ID:       "ZC1644",
		Title:    "Error on `unzip -P SECRET` / `zip -P SECRET` — archive password in process list",
		Severity: SeverityError,
		Description: "`unzip -P PASSWORD` / `zip -P PASSWORD` (or the concatenated `-PPASSWORD` " +
			"form) places the archive password in argv. The expanded value shows up in `ps`, " +
			"`/proc/<pid>/cmdline`, shell history, and audit logs for every local user who " +
			"can list processes. Both tools prompt interactively if `-P` is absent — use that " +
			"for human workflows. For automation prefer an archive format with a real key-" +
			"derivation story (for example `7z -p` piped over stdin, or `age` / `gpg` " +
			"envelope encryption that reads keys from a protected file).",
		Check: checkZC1644,
	})
}

func checkZC1644(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "unzip" && ident.Value != "zip" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-P" && i+1 < len(cmd.Arguments) {
			return zc1644Hit(cmd, ident.Value)
		}
		if strings.HasPrefix(v, "-P") && len(v) > 2 {
			return zc1644Hit(cmd, ident.Value)
		}
	}
	return nil
}

func zc1644Hit(cmd *ast.SimpleCommand, name string) []Violation {
	return []Violation{{
		KataID: "ZC1644",
		Message: "`" + name + " -P` places the archive password in argv — visible via " +
			"`ps`. Drop `-P` for interactive prompt, or switch to `7z -p` (reads from " +
			"stdin) / `age` / `gpg` with keys in a protected file.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1645",
		Title:    "Style: `lsb_release` — prefer sourcing `/etc/os-release` (no dependency, no fork)",
		Severity: SeverityStyle,
		Description: "`lsb_release` is provided by the `lsb-release` / `redhat-lsb-core` " +
			"package, which is missing on most minimal / container images (Alpine does not " +
			"ship it at all). Scripts that depend on `lsb_release` fail the moment they hit " +
			"a stripped image. `/etc/os-release` is standardized by systemd and always " +
			"present on modern Linux — `source /etc/os-release; print -r -- $ID $VERSION_ID` " +
			"gives the same distribution info without the extra package, and without forking.",
		Check: checkZC1645,
	})
}

func checkZC1645(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "lsb_release" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1645",
		Message: "`lsb_release` needs an optional package. Use `source /etc/os-release` and " +
			"read `$ID` / `$VERSION_ID` / `$PRETTY_NAME` instead — always present, no fork.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1646",
		Title:    "Warn on `btrfs check --repair` / `xfs_repair -L` — last-resort recovery, may worsen damage",
		Severity: SeverityWarning,
		Description: "Both commands are destructive last-resort recovery. `btrfs check " +
			"--repair` explicitly warns in its man page that it \"may cause additional " +
			"filesystem damage\" and the btrfs developers ask users to try `btrfs scrub` and " +
			"read-only `btrfs check` first. `xfs_repair -L` zeroes the log, dropping any " +
			"uncommitted transactions and the data they held. In both cases snapshot the " +
			"underlying block device before running, so the attempt is reversible.",
		Check: checkZC1646,
	})
}

func checkZC1646(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value == "btrfs" {
		if len(cmd.Arguments) < 2 || cmd.Arguments[0].String() != "check" {
			return nil
		}
		for _, arg := range cmd.Arguments[1:] {
			if arg.String() == "--repair" {
				return []Violation{{
					KataID: "ZC1646",
					Message: "`btrfs check --repair` may worsen damage — try `btrfs scrub` " +
						"and read-only `btrfs check` first, and snapshot the block device " +
						"before running.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
		return nil
	}
	if ident.Value == "xfs_repair" {
		for _, arg := range cmd.Arguments {
			if arg.String() == "-L" {
				return []Violation{{
					KataID: "ZC1646",
					Message: "`xfs_repair -L` zeroes the log — uncommitted transactions are " +
						"lost. Snapshot the block device first; mount read-only and read " +
						"the log if possible.",
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
		ID:       "ZC1647",
		Title:    "Warn on `kubectl apply -f URL` — remote manifest applied without digest verification",
		Severity: SeverityWarning,
		Description: "`kubectl apply -f https://...` fetches the manifest over the network and " +
			"applies it to the cluster. TLS (when present) verifies transport but not " +
			"authorship — if the URL is compromised or the content changes between reviews, " +
			"the cluster picks up the new definition. Pin the content: download to disk, " +
			"verify a known SHA256, then `kubectl apply -f local.yaml`. For plain HTTP the " +
			"attacker controls the response directly — never acceptable.",
		Check: checkZC1647,
	})
}

func checkZC1647(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "apply" && sub != "create" && sub != "replace" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		if arg.String() != "-f" {
			continue
		}
		idx := i + 2
		if idx >= len(cmd.Arguments) {
			continue
		}
		target := cmd.Arguments[idx].String()
		if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
			return []Violation{{
				KataID: "ZC1647",
				Message: "`kubectl " + sub + " -f " + target + "` applies a remote " +
					"manifest — verify digest first. Download, check SHA256, then apply " +
					"the local file.",
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
		ID:       "ZC1648",
		Title:    "Error on `cp /dev/null /var/log/...` / `truncate -s 0 /var/log/...` — audit-log wipe",
		Severity: SeverityError,
		Description: "Replacing a file under `/var/log/` with `/dev/null` or truncating it to " +
			"size zero erases audit evidence: failed login attempts from `auth.log`, sudo " +
			"usage from `sudo.log`, kernel audit trail from `audit/audit.log`, console " +
			"history from `wtmp` / `btmp`. Scripts that do this during \"cleanup\" are almost " +
			"always misusing logrotate (which handles rotation safely via a `create` stage) " +
			"or deliberately covering tracks. Use `logrotate -f /etc/logrotate.d/...` for " +
			"rotation, `journalctl --vacuum-time=...` for journald.",
		Check: checkZC1648,
	})
}

func checkZC1648(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "cp":
		if len(cmd.Arguments) < 2 {
			return nil
		}
		if cmd.Arguments[0].String() != "/dev/null" {
			return nil
		}
		dest := cmd.Arguments[1].String()
		if strings.HasPrefix(dest, "/var/log/") {
			return zc1648Hit(cmd, "cp /dev/null "+dest)
		}
	case "truncate":
		var zeroSize bool
		var target string
		for i, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-s" && i+1 < len(cmd.Arguments) && cmd.Arguments[i+1].String() == "0" {
				zeroSize = true
			}
			if strings.HasPrefix(v, "/var/log/") {
				target = v
			}
		}
		if zeroSize && target != "" {
			return zc1648Hit(cmd, "truncate -s 0 "+target)
		}
	}
	return nil
}

func zc1648Hit(cmd *ast.SimpleCommand, desc string) []Violation {
	return []Violation{{
		KataID: "ZC1648",
		Message: "`" + desc + "` wipes an audit log — use `logrotate -f` or " +
			"`journalctl --vacuum-time=...` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1649",
		Title:    "Warn on `openssl req -days N` with N > 825 — long-validity certificate",
		Severity: SeverityWarning,
		Description: "CA/Browser Forum capped public TLS cert validity at 825 days in 2018 and " +
			"major browsers tightened it to 398 days in 2020. A cert issued for 3650 days " +
			"(10 years) can not be revoked effectively — once the private key leaks, the " +
			"attacker keeps access until the cert expires naturally. For an internal root CA " +
			"the long validity is defensible; for leaf / server certs keep it under 398 " +
			"days and automate rotation. `-days` over 825 almost always means \"I don't want " +
			"to deal with renewal,\" which is a maintenance smell dressed up as security.",
		Check: checkZC1649,
	})
}

func checkZC1649(node ast.Node) []Violation {
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
	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "req" && sub != "x509" && sub != "ca" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() != "-days" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		days, err := strconv.Atoi(cmd.Arguments[i+1].String())
		if err != nil {
			continue
		}
		if days > 825 {
			return []Violation{{
				KataID: "ZC1649",
				Message: "`openssl " + sub + " -days " + cmd.Arguments[i+1].String() +
					"` issues a cert with a long validity. Keep leaf certs under 398 " +
					"days and automate rotation.",
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
		ID:       "ZC1650",
		Title:    "Warn on `setopt RM_STAR_SILENT` / `unsetopt RM_STAR_WAIT` — removes `rm *` prompt",
		Severity: SeverityWarning,
		Description: "Zsh's default behaviour on an interactive `rm *` (or `rm /path/*`) is to " +
			"pause for 10 seconds and ask \"do you really want to delete N files?\" — the " +
			"`RM_STAR_WAIT` option. `setopt RM_STAR_SILENT` or `unsetopt RM_STAR_WAIT` both " +
			"disable the prompt. In a profile / dot file the option leaks to every future " +
			"interactive shell and removes a safety net that has saved countless home " +
			"directories.",
		Check: checkZC1650,
	})
}

func checkZC1650(node ast.Node) []Violation {
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
			norm := strings.ToLower(strings.ReplaceAll(arg.String(), "_", ""))
			if norm == "rmstarsilent" {
				return zc1650Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			norm := strings.ToLower(strings.ReplaceAll(arg.String(), "_", ""))
			if norm == "rmstarwait" {
				return zc1650Hit(cmd, "unsetopt "+arg.String())
			}
		}
	}
	return nil
}

func zc1650Hit(cmd *ast.SimpleCommand, desc string) []Violation {
	return []Violation{{
		KataID: "ZC1650",
		Message: "`" + desc + "` removes the `rm *` confirmation prompt — keep the default " +
			"`RM_STAR_WAIT` so accidental deletions pause before they happen.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1651",
		Title:    "Warn on `docker/podman run -p 0.0.0.0:PORT:PORT` — explicit all-interfaces publish",
		Severity: SeverityWarning,
		Description: "A port spec of `0.0.0.0:HOST:CONT`, `[::]:HOST:CONT`, or `*:HOST:CONT` " +
			"publishes the container port to every interface the host has. On a multi-" +
			"tenant LAN or a cloud host with a public IP the service is immediately reachable " +
			"from anywhere. If the service needs only local reverse-proxy access, bind to " +
			"`127.0.0.1:HOST:CONT` and let nginx / caddy handle external exposure.",
		Check: checkZC1651,
	})
}

func checkZC1651(node ast.Node) []Violation {
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
	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "run" && sub != "create" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v != "-p" && v != "--publish" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		spec := strings.Trim(cmd.Arguments[i+1].String(), "\"'")
		if strings.HasPrefix(spec, "0.0.0.0:") ||
			strings.HasPrefix(spec, "[::]:") ||
			strings.HasPrefix(spec, "*:") {
			return []Violation{{
				KataID: "ZC1651",
				Message: "`" + ident.Value + " " + sub + " -p " + spec + "` publishes to " +
					"every interface. Bind to `127.0.0.1:HOST:CONT` and put nginx / caddy " +
					"in front for external access.",
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
		ID:       "ZC1652",
		Title:    "Warn on `ssh -Y` — trusted X11 forwarding grants full X-server access to remote clients",
		Severity: SeverityWarning,
		Description: "`ssh -Y` enables trusted X11 forwarding. Remote X clients can read every " +
			"keystroke on the local display, take screenshots, inject synthetic events, and " +
			"otherwise drive the local session with no sandbox. `ssh -X` enables the " +
			"untrusted variant, which routes X traffic through the X SECURITY extension so " +
			"those capabilities are limited (some GUI features break, which is why people " +
			"reach for `-Y` — usually at far higher risk than they realised). Prefer `-X` " +
			"when X11 forwarding is genuinely needed; better yet drop it for Wayland tools " +
			"or VNC-over-SSH with its own auth.",
		Check: checkZC1652,
	})
}

func checkZC1652(node ast.Node) []Violation {
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

	for _, arg := range cmd.Arguments {
		if arg.String() == "-Y" {
			return []Violation{{
				KataID: "ZC1652",
				Message: "`ssh -Y` enables trusted X11 forwarding — remote clients get full " +
					"access to the local X server. Use `-X` (untrusted) or drop X11 " +
					"forwarding entirely.",
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
		ID:       "ZC1653",
		Title:    "Avoid `$BASHPID` — Bash-only; Zsh uses `$sysparams[pid]` from `zsh/system`",
		Severity: SeverityWarning,
		Description: "`$BASHPID` returns the PID of the current subshell (while `$$` returns " +
			"the parent shell's PID). In Zsh this parameter is not set — scripts that rely " +
			"on `$BASHPID` silently get an empty string and misbehave. After `zmodload " +
			"zsh/system`, Zsh exposes the current process PID as `$sysparams[pid]`, which " +
			"updates inside subshells just like Bash's `$BASHPID`.",
		Check: checkZC1653,
	})
}

func checkZC1653(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "$BASHPID") || strings.Contains(v, "${BASHPID}") {
			return []Violation{{
				KataID: "ZC1653",
				Message: "`$BASHPID` is Bash-only. Use `$sysparams[pid]` after " +
					"`zmodload zsh/system`.",
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
		ID:       "ZC1654",
		Title:    "Warn on `sysctl -p /tmp/...` — loading kernel tunables from attacker-writable path",
		Severity: SeverityWarning,
		Description: "`sysctl -p PATH` reads `key=value` lines from PATH and applies them as " +
			"kernel tunables. A PATH under `/tmp/` or `/var/tmp/` is world-traversable; a " +
			"concurrent local user can substitute the file between write and read, " +
			"injecting `kernel.core_pattern=|/tmp/evil`, `kernel.modprobe=/tmp/evil`, or " +
			"disabling hardening knobs (`kernel.kptr_restrict=0`, `kernel.yama.ptrace_scope=" +
			"0`). Keep sysctl configs under `/etc/sysctl.d/` with root ownership.",
		Check: checkZC1654,
	})
}

func checkZC1654(node ast.Node) []Violation {
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

	for i, arg := range cmd.Arguments {
		if arg.String() != "-p" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		next := cmd.Arguments[i+1].String()
		if strings.HasPrefix(next, "/tmp/") || strings.HasPrefix(next, "/var/tmp/") {
			return []Violation{{
				KataID: "ZC1654",
				Message: "`sysctl -p " + next + "` reads tunables from a world-traversable " +
					"path — a concurrent local user can substitute the file. Keep configs " +
					"under `/etc/sysctl.d/`.",
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
		ID:       "ZC1655",
		Title:    "Warn on `read -n N` — Bash reads N chars; Zsh's `-n` means \"drop newline\"",
		Severity: SeverityWarning,
		Description: "In Bash, `read -n N var` reads exactly N characters (handy for single-" +
			"keypress prompts). In Zsh, `-n` is the \"don't append newline to the reply " +
			"string\" flag and doesn't take a count — `read -n 1 var` sets `var` to the " +
			"whole line, not a single character. Use `read -k N var` in Zsh for N-character " +
			"reads.",
		Check: checkZC1655,
	})
}

func checkZC1655(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "read" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-n" && i+1 < len(cmd.Arguments) {
			if _, err := strconv.Atoi(cmd.Arguments[i+1].String()); err == nil {
				return zc1655Hit(cmd)
			}
		}
		if strings.HasPrefix(v, "-n") && len(v) > 2 {
			if _, err := strconv.Atoi(v[2:]); err == nil {
				return zc1655Hit(cmd)
			}
		}
	}
	return nil
}

func zc1655Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1655",
		Message: "`read -n N` is Bash syntax for \"read N characters\". Zsh's `-n` means " +
			"\"drop trailing newline\" with no count. Use `read -k N var` in Zsh.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1656",
		Title:    "Error on `rsync -e 'ssh -o StrictHostKeyChecking=no'` — host-key verify disabled",
		Severity: SeverityError,
		Description: "Disabling host-key verification through rsync's `-e` transport is the " +
			"same attack surface as ZC1479 but easier to miss in review because the ssh flags " +
			"sit inside a quoted string. A MITM on the network path can impersonate the " +
			"remote host and the rsync stream goes straight through. Use `ssh-keyscan` or " +
			"pre-provisioned `~/.ssh/known_hosts` to trust hosts deliberately, and keep " +
			"`StrictHostKeyChecking=yes`.",
		Check: checkZC1656,
	})
}

func checkZC1656(node ast.Node) []Violation {
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

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "StrictHostKeyChecking=no") ||
			strings.Contains(v, "UserKnownHostsFile=/dev/null") {
			return []Violation{{
				KataID: "ZC1656",
				Message: "`rsync -e 'ssh -o StrictHostKeyChecking=no'` disables host-key " +
					"verification — MITM risk. Pre-provision `known_hosts` and keep " +
					"`StrictHostKeyChecking=yes`.",
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
		ID:       "ZC1657",
		Title:    "Warn on `semanage permissive -a <type>` — puts SELinux domain in permissive mode",
		Severity: SeverityWarning,
		Description: "`semanage permissive -a DOMAIN` (or `--add`) marks an SELinux domain as " +
			"permissive: policy violations are logged but not blocked. It is narrower than " +
			"`setenforce 0` but still disables enforcement for whatever DOMAIN covers — often " +
			"`httpd_t`, `container_t`, or `sshd_t` — and the override persists across reboots " +
			"because it is written to policy. Fix the denial with an explicit allow rule built " +
			"from `audit2allow` or ship a custom policy module, and remove the permissive mark " +
			"with `semanage permissive -d DOMAIN` once the rule lands.",
		Check: checkZC1657,
	})
}

func checkZC1657(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "semanage" {
		return nil
	}

	hasPermissive := false
	hasAdd := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "permissive":
			hasPermissive = true
		case "-a", "--add":
			hasAdd = true
		}
	}

	if !hasPermissive || !hasAdd {
		return nil
	}

	return []Violation{{
		KataID: "ZC1657",
		Message: "`semanage permissive -a` puts an SELinux domain in permissive mode — " +
			"policy violations log but no longer block. Write a scoped allow rule with " +
			"`audit2allow` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1658",
		Title:    "Warn on `curl -OJ` / `-J -O` — server-controlled output filename",
		Severity: SeverityWarning,
		Description: "`curl -J` (`--remote-header-name`) combined with `-O` (`--remote-name`) " +
			"saves the response using the filename the server puts in the `Content-Disposition` " +
			"header. The server — or anything on the path that can set headers, including a " +
			"compromised CDN or an HTTP-serving reverse proxy — chooses the destination name. " +
			"Paths like `../../etc/cron.d/evil` are rejected by curl's sanitizer, but benign-" +
			"looking names still overwrite files in the current directory. Use `-o NAME` with " +
			"a filename you control, and validate the payload before you act on it.",
		Check: checkZC1658,
	})
}

func checkZC1658(node ast.Node) []Violation {
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

	hasJ := false
	hasO := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--remote-header-name" {
			hasJ = true
			continue
		}
		if v == "--remote-name" {
			hasO = true
			continue
		}
		if !strings.HasPrefix(v, "-") || strings.HasPrefix(v, "--") {
			continue
		}
		body := strings.TrimPrefix(v, "-")
		if strings.ContainsRune(body, 'J') {
			hasJ = true
		}
		if strings.ContainsRune(body, 'O') {
			hasO = true
		}
	}

	if !hasJ || !hasO {
		return nil
	}

	return []Violation{{
		KataID: "ZC1658",
		Message: "`curl -OJ` saves the response under the name the server picks in " +
			"`Content-Disposition` — path traversal is blocked but arbitrary same-dir " +
			"overwrites are not. Pass `-o NAME` with a filename you control.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1659",
		Title:    "Warn on `fuser -k <path>` — kills every process holding the subtree open",
		Severity: SeverityWarning,
		Description: "`fuser -k PATH` sends a signal (SIGKILL by default) to every process that " +
			"has any file under PATH open — not just the one you expected. On `/`, `/var`, " +
			"`/tmp`, or any mount-root this reaches sshd, cron, dbus, and the caller's own " +
			"shell; on a bind-mount it kills workloads that share the host inode. Target " +
			"specific PIDs (`kill $(pidof app)`) or ports (`fuser -k PORT/tcp`), or use " +
			"`systemctl stop UNIT` for services. `fuser -k` against a filesystem path is " +
			"blast-radius that the caller rarely owns.",
		Check: checkZC1659,
	})
}

func checkZC1659(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "fuser" {
		return nil
	}

	hasKill := false
	pathTarget := ""
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "-") && !strings.HasPrefix(v, "--") {
			if strings.ContainsRune(strings.TrimPrefix(v, "-"), 'k') {
				hasKill = true
			}
			continue
		}
		if strings.HasPrefix(v, "/") {
			if strings.HasSuffix(v, "/tcp") || strings.HasSuffix(v, "/udp") ||
				strings.HasSuffix(v, "/sctp") {
				continue
			}
			pathTarget = v
		}
	}

	if !hasKill || pathTarget == "" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1659",
		Message: "`fuser -k " + pathTarget + "` signals every process with a file open " +
			"anywhere under the path — use PID / port targets or `systemctl stop` for " +
			"services.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

var zc1660ZeroPad = regexp.MustCompile(`%0[1-9][0-9]*d`)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1660",
		Title:    "Style: `printf '%0Nd' $n` — prefer Zsh `${(l:N::0:)n}` left-zero-pad",
		Severity: SeverityStyle,
		Description: "Zero-padding an integer through `printf '%0Nd'` forks a tiny sub-process " +
			"and relies on printf's format-string parser — both things Zsh can avoid. " +
			"`${(l:N::0:)n}` left-pads `$n` with `0` to width N using Zsh parameter " +
			"expansion, no fork, and composes cleanly with other `(q)` / `(L)` / `(U)` " +
			"flags. For right-pad use `${(r:N::0:)n}`; for space padding swap the fill " +
			"character: `${(l:N:)n}` or `${(r:N:)n}`.",
		Check: checkZC1660,
	})
}

func checkZC1660(node ast.Node) []Violation {
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

	if len(cmd.Arguments) == 0 {
		return nil
	}

	fmtArg := cmd.Arguments[0].String()
	if !zc1660ZeroPad.MatchString(fmtArg) {
		return nil
	}

	return []Violation{{
		KataID: "ZC1660",
		Message: "`printf '%0Nd'` forks for zero-padding — prefer Zsh `${(l:N::0:)n}` " +
			"parameter-expansion pad (same for `(r:N::0:)` on the right).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1661",
		Title:    "Error on `curl --cacert /dev/null` — empty trust store, any cert passes",
		Severity: SeverityError,
		Description: "Pointing `--cacert` (or `--capath`) at `/dev/null` hands curl an empty " +
			"trust anchor set. Counter-intuitively, curl treats the peer certificate as " +
			"valid when no issuers are configured for the selected TLS backend (OpenSSL, " +
			"wolfSSL, Schannel all accept any cert chain against an empty CA bundle). This is " +
			"the TLS equivalent of `--insecure` with one more keystroke of plausible " +
			"deniability. Use a real bundle (`/etc/ssl/certs/ca-certificates.crt`) or " +
			"`--pinnedpubkey sha256//…` for known endpoints.",
		Check: checkZC1661,
	})
}

func checkZC1661(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "cacert", "capath":
		if len(cmd.Arguments) > 0 && cmd.Arguments[0].String() == "/dev/null" {
			return zc1661Hit(cmd)
		}
		return nil
	case "curl":
	default:
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v != "--cacert" && v != "--capath" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		if cmd.Arguments[i+1].String() == "/dev/null" {
			return zc1661Hit(cmd)
		}
	}
	return nil
}

func zc1661Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1661",
		Message: "`curl --cacert /dev/null` feeds curl an empty trust store — most TLS " +
			"backends then accept any peer cert. Use a real bundle or `--pinnedpubkey`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1662",
		Title:    "Error on `pkexec env VAR=VAL CMD` — controlled env crossed into the root session",
		Severity: SeverityError,
		Description: "`pkexec env VAR=VALUE CMD` invokes `/usr/bin/env` as the target user (root " +
			"by default) with a caller-controlled environment. Polkit sanitizes a short " +
			"allow-list on its own, but once `env` takes over the remaining variables " +
			"(`LD_PRELOAD`, `GCONV_PATH`, `PYTHONPATH`, `XDG_RUNTIME_DIR`, `LANGUAGE`) ride " +
			"straight into root. CVE-2021-4034 (pwnkit) demonstrated the same primitive by " +
			"abusing argv[0]; the `env` wrapper makes the bypass trivial. If the child needs " +
			"specific variables, set them in a polkit rule or via `systemd-run --user` " +
			"instead, not through `env`.",
		Check: checkZC1662,
	})
}

func checkZC1662(node ast.Node) []Violation {
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
	if cmd.Arguments[0].String() != "env" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1662",
		Message: "`pkexec env VAR=VAL CMD` hands the root session a caller-controlled " +
			"environment — use a polkit rule or `systemd-run --user` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1663",
		Title:    "Warn on `tune2fs -c 0` / `-i 0` — disables periodic filesystem checks",
		Severity: SeverityWarning,
		Description: "`tune2fs -c 0` (mount count) and `tune2fs -i 0` (time interval) disable " +
			"the ext2/3/4 periodic-fsck machinery so the filesystem only gets checked after a " +
			"dirty unmount or a manual `fsck -f`. For desktops the nag is annoying; for " +
			"long-lived servers it is the last line of defence against silent metadata " +
			"corruption. Lower the cadence if the default is too aggressive (`tune2fs -c 30`, " +
			"`-i 3m`) rather than turning it off, and schedule an offline `fsck` on a cadence " +
			"you can defend.",
		Check: checkZC1663,
	})
}

func checkZC1663(node ast.Node) []Violation {
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
		if v != "-c" && v != "-i" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		next := cmd.Arguments[i+1].String()
		if next != "0" {
			continue
		}
		return []Violation{{
			KataID: "ZC1663",
			Message: "`tune2fs " + v + " 0` disables periodic fsck on the filesystem — " +
				"lower the cadence (e.g. `-c 30` / `-i 3m`) instead of turning it off.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1664",
		Title:    "Error on `systemctl set-default rescue.target|emergency.target` — persistent single-user boot",
		Severity: SeverityError,
		Description: "`systemctl set-default` rewrites `/etc/systemd/system/default.target` as a " +
			"symlink to the named target. Pointing it at `rescue.target` or " +
			"`emergency.target` means every subsequent boot drops to single-user mode " +
			"before networking, sshd, or any normal unit starts — you lose remote access to " +
			"the box unless you have serial console / out-of-band management. Unlike " +
			"`systemctl isolate` (one-shot, caught by ZC1561) this persists across reboots. " +
			"Revert with `systemctl set-default multi-user.target` (servers) or `graphical." +
			"target` (desktops).",
		Check: checkZC1664,
	})
}

func checkZC1664(node ast.Node) []Violation {
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
	if cmd.Arguments[0].String() != "set-default" {
		return nil
	}
	target := cmd.Arguments[1].String()
	if target != "rescue.target" && target != "emergency.target" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1664",
		Message: "`systemctl set-default " + target + "` makes every subsequent boot land " +
			"in single-user mode — revert with `set-default multi-user.target` or " +
			"`graphical.target`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1665",
		Title:    "Warn on `chrt -r` / `-f` — real-time scheduling class from a shell script",
		Severity: SeverityWarning,
		Description: "`chrt -r PRIO CMD` (SCHED_RR) and `chrt -f PRIO CMD` (SCHED_FIFO) launch " +
			"the child under a POSIX real-time scheduling class. An RT thread preempts " +
			"every normal-priority task until it voluntarily yields; a busy-loop or a " +
			"deadlock leaves the kernel with kworker, ksoftirqd, and sshd starved, often " +
			"forcing a hard reboot. Unless the binary is known-bounded (audio glitch-free " +
			"path, protocol timing loop), keep scripts on SCHED_OTHER — use `nice -n -5` or " +
			"a systemd unit with `CPUWeight=` / `IOWeight=` instead of `chrt -r`.",
		Check: checkZC1665,
	})
}

func checkZC1665(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chrt" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-r" || v == "-f" || v == "--rr" || v == "--fifo" {
			return []Violation{{
				KataID: "ZC1665",
				Message: "`chrt " + v + "` puts the child on a real-time scheduling class — a " +
					"busy-loop or deadlock then starves kworker / sshd. Prefer `nice -n -5` " +
					"or a systemd unit with `CPUWeight=`.",
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
		ID:       "ZC1666",
		Title:    "Warn on `kubectl patch --type=json` — bypasses strategic-merge defaults",
		Severity: SeverityWarning,
		Description: "`kubectl patch --type=json` applies a raw RFC-6902 JSON patch: `remove`, " +
			"`replace`, `add /spec/containers/0`, and `move` land verbatim on the resource. " +
			"Unlike strategic-merge or merge-patch, Kubernetes does not reconcile the " +
			"patch against field ownership or default values — so a mistyped `path` or an " +
			"index that no longer exists fails silently or drops the wrong field. From a " +
			"script this is a foot-gun for drift and supply-chain compromise: an attacker " +
			"with write access to the patch file can slip `privileged: true` or `hostPath` " +
			"mounts in. Prefer `--type=strategic` (the default) and hold JSON patches " +
			"behind code review.",
		Check: checkZC1666,
	})
}

func checkZC1666(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" {
		return nil
	}

	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "patch" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--type=json" {
			return zc1666Hit(cmd)
		}
		if v == "--type" && i+1 < len(cmd.Arguments) &&
			cmd.Arguments[i+1].String() == "json" {
			return zc1666Hit(cmd)
		}
	}
	return nil
}

func zc1666Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1666",
		Message: "`kubectl patch --type=json` applies a raw RFC-6902 patch that bypasses " +
			"strategic-merge reconciliation — prefer `--type=strategic` and hold JSON " +
			"patches behind code review.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1667",
		Title:    "Warn on `openssl enc` without `-pbkdf2` — legacy MD5-based key derivation",
		Severity: SeverityWarning,
		Description: "Without `-pbkdf2`, `openssl enc` derives the symmetric key through " +
			"EVP_BytesToKey, which is a single MD5 round over `password || salt`. A modern " +
			"GPU cracks that at billions of guesses per second. Add `-pbkdf2 -iter 100000` " +
			"(OpenSSL 1.1.1+) to switch to PBKDF2-HMAC-SHA256 with a real iteration count. " +
			"Even better, stop using `openssl enc` for new code — it has no AEAD support and " +
			"`-aes-256-gcm` silently drops the auth tag — and reach for `age`, " +
			"`gpg --symmetric --cipher-algo AES256`, or `openssl smime` instead.",
		Check: checkZC1667,
	})
}

func checkZC1667(node ast.Node) []Violation {
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
		if arg.String() == "-pbkdf2" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1667",
		Message: "`openssl enc` without `-pbkdf2` uses single-round EVP_BytesToKey (MD5) — " +
			"add `-pbkdf2 -iter 100000`, or prefer `age` / `gpg --symmetric`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1668",
		Title:    "Error on `aws iam attach-*-policy ... AdministratorAccess` — grants full AWS admin",
		Severity: SeverityError,
		Description: "Attaching the AWS-managed `AdministratorAccess` (or `PowerUserAccess`) " +
			"policy gives the target principal `*:*` — create/delete IAM users, mutate KMS " +
			"keys, rotate root passwords, exfiltrate every S3 bucket. Scripts rarely need " +
			"full admin; the pattern usually means someone hit a permissions error and " +
			"replaced the scoped policy with the blanket one. Write a least-privilege inline " +
			"policy (`iam put-user-policy --policy-document`), or reference a customer-" +
			"managed policy with only the `Action`/`Resource` pairs the workload needs. Admin " +
			"attachment should land via change-reviewed Terraform, not a shell loop.",
		Check: checkZC1668,
	})
}

func checkZC1668(node ast.Node) []Violation {
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

	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "iam" {
		return nil
	}
	sub := cmd.Arguments[1].String()
	if sub != "attach-user-policy" && sub != "attach-role-policy" &&
		sub != "attach-group-policy" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if strings.HasSuffix(v, "/AdministratorAccess") ||
			strings.HasSuffix(v, "/PowerUserAccess") ||
			strings.HasSuffix(v, "/IAMFullAccess") {
			return []Violation{{
				KataID: "ZC1668",
				Message: "`aws iam " + sub + " ... " + v + "` grants sweeping admin — " +
					"use a scoped inline policy (`put-user-policy`) or a customer-managed " +
					"policy with the minimum `Action`/`Resource` set.",
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
		ID:       "ZC1669",
		Title:    "Warn on `git gc --prune=now` / `git reflog expire --expire=now` — deletes recovery window",
		Severity: SeverityWarning,
		Description: "Git keeps dropped commits and orphaned objects for `gc.reflogExpire` " +
			"(default 90 days) and `gc.pruneExpire` (default two weeks) so a `git reflog` + " +
			"`git reset` can still recover work you thought you threw away. `git gc " +
			"--prune=now` and `git reflog expire --expire=now --all` bulldoze both windows " +
			"in one go — a stray interactive rebase no longer has a safety net. Use the " +
			"default cadence (`git gc`, no `--prune=now`) unless you are actively purging " +
			"leaked secrets or proof-of-concept code; pair the destructive form with a " +
			"stale mirror push so at least one copy of the dropped history remains.",
		Check: checkZC1669,
	})
}

func checkZC1669(node ast.Node) []Violation {
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

	if len(cmd.Arguments) == 0 {
		return nil
	}

	switch cmd.Arguments[0].String() {
	case "gc":
		for _, arg := range cmd.Arguments[1:] {
			if arg.String() == "--prune=now" || arg.String() == "--prune=0" {
				return zc1669Hit(cmd, "git gc --prune=now")
			}
		}
	case "reflog":
		if len(cmd.Arguments) < 2 || cmd.Arguments[1].String() != "expire" {
			return nil
		}
		for _, arg := range cmd.Arguments[2:] {
			if arg.String() == "--expire=now" {
				return zc1669Hit(cmd, "git reflog expire --expire=now")
			}
		}
	}
	return nil
}

func zc1669Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1669",
		Message: "`" + form + "` bulldozes the reflog / prune recovery window — keep the " +
			"default cadence unless you are actively purging leaked secrets, and mirror " +
			"the dropped history off-box first.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

var zc1670Dangerous = map[string]struct{}{
	"allow_execstack":            {},
	"allow_execmod":              {},
	"allow_execmem":              {},
	"httpd_execmem":              {},
	"httpd_unified":              {},
	"selinuxuser_execmod":        {},
	"selinuxuser_execstack":      {},
	"selinuxuser_execheap":       {},
	"domain_kernel_load_modules": {},
	"deny_ptrace":                {},
	"mmap_low_allowed":           {},
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1670",
		Title:    "Warn on `setsebool -P` enabling memory-protection-relaxing SELinux boolean",
		Severity: SeverityWarning,
		Description: "Specific SELinux policy booleans (`allow_execstack`, `allow_execmem`, " +
			"`httpd_execmem`, `selinuxuser_execstack`, `domain_kernel_load_modules`, " +
			"`mmap_low_allowed`, etc.) relax per-domain memory protections that the policy " +
			"puts in place precisely because those domains should not need writable-and-" +
			"executable pages. Persisting the flip with `-P` carries the regression across " +
			"reboots. Fix the underlying binary (`execstack -c`, `chcon`, stop generating " +
			"runtime-JIT code in the wrong domain) instead of loosening policy.",
		Check: checkZC1670,
	})
}

func checkZC1670(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setsebool" {
		return nil
	}

	hasPersist := false
	boolName := ""
	boolValue := ""
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch {
		case v == "-P":
			hasPersist = true
		case boolName == "":
			boolName = v
		case boolValue == "":
			boolValue = v
		}
	}

	if !hasPersist || boolName == "" || boolValue == "" {
		return nil
	}
	if _, dangerous := zc1670Dangerous[boolName]; !dangerous {
		return nil
	}
	if boolValue != "1" && boolValue != "on" && boolValue != "true" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1670",
		Message: "`setsebool -P " + boolName + " " + boolValue + "` persistently relaxes " +
			"SELinux memory-protection policy — fix the binary instead (`execstack -c`, " +
			"relabel with `chcon`, or change the domain).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1671",
		Title:    "Error on `install -m 777` / `mkdir -m 777` — creates world-writable target",
		Severity: SeverityError,
		Description: "`install -m MODE` / `mkdir -m MODE` applies MODE atomically at file or " +
			"directory creation, so the world-writable window from a later `chmod 777` is " +
			"not even needed — the path is wide-open from the moment it exists. Any local " +
			"user can swap binaries under `/usr/local/bin`, write shell-completion hooks " +
			"into `/etc/bash_completion.d`, or turn a shared directory into an LPE staging " +
			"ground. Drop the world-write bit: `0755` for binaries, `0644` for files, `2770` " +
			"with `chgrp` for shared directories.",
		Check: checkZC1671,
	})
}

func checkZC1671(node ast.Node) []Violation {
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

	for i, arg := range cmd.Arguments {
		if arg.String() != "-m" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		mode := cmd.Arguments[i+1].String()
		if !zc1671WorldWritable(mode) {
			continue
		}
		return []Violation{{
			KataID: "ZC1671",
			Message: "`" + ident.Value + " -m " + mode + "` creates a world-writable " +
				"target — drop the world-write bit (e.g. `0755` / `0644`).",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}

// zc1671WorldWritable returns true if MODE has the world-write (o+w) bit set.
// Users spell modes in octal. If the literal parses as octal, trust that
// reading. Otherwise (a digit 8/9 appears — that only happens because the
// parser normalized a leading-zero octal like `0666` to decimal `438`), parse
// as decimal and still check the o+w bit.
func zc1671WorldWritable(mode string) bool {
	if n, err := strconv.ParseInt(mode, 8, 32); err == nil {
		return n > 0 && n&0o002 != 0
	}
	if n, err := strconv.ParseInt(mode, 10, 32); err == nil && n > 0 && n&0o002 != 0 {
		return true
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1672",
		Title:    "Info: `chcon` writes an ephemeral SELinux label — next `restorecon` wipes it",
		Severity: SeverityInfo,
		Description: "`chcon -t TYPE PATH` sets the file context out-of-band; it does not update " +
			"the `file_contexts` policy database. As soon as `restorecon`, `semodule -n`, or " +
			"a policy rebuild runs, the label snaps back to whatever the compiled policy " +
			"says — often `default_t`, which can break a deployed workload or silently " +
			"re-introduce a denial the script tried to fix. For anything long-lived use " +
			"`semanage fcontext -a -t TYPE '<regex>'` then `restorecon -F <path>` so the " +
			"mapping lives in policy.",
		Check: checkZC1672,
	})
}

func checkZC1672(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chcon" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1672",
		Message: "`chcon` writes an ephemeral SELinux label — `restorecon` / policy rebuild " +
			"reverts it. Persist via `semanage fcontext -a -t TYPE 'REGEX'` + `restorecon`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1673",
		Title:    "Style: `stty -echo` around `read` — prefer Zsh `read -s`",
		Severity: SeverityStyle,
		Description: "The classic `stty -echo; IFS= read -r password; stty echo` pattern has a " +
			"serious failure mode: a crash or SIGINT between the two `stty` calls leaves " +
			"the user's terminal stuck in echo-off, which is silent and confusing. Zsh's " +
			"`read -s VAR` (also Bash 4+) disables echo only for that one `read`, restores " +
			"it on return even if the read is interrupted, and avoids two external forks. " +
			"Switch the prompt to `read -s` (or `read -ks` for single-key password) and " +
			"drop the `stty` bracketing.",
		Check: checkZC1673,
	})
}

func checkZC1673(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "stty" {
		return nil
	}

	if len(cmd.Arguments) != 1 {
		return nil
	}
	if cmd.Arguments[0].String() != "-echo" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1673",
		Message: "`stty -echo` to mask password entry is fragile — a crash leaves the " +
			"terminal echo-off. Use `read -s VAR` (Zsh / Bash 4+) instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1674",
		Title:    "Warn on `docker/podman run --oom-kill-disable` or `--oom-score-adj <= -500`",
		Severity: SeverityWarning,
		Description: "`--oom-kill-disable` tells the kernel OOM killer to never touch the " +
			"container's memory cgroup — a leak inside then drives the whole host into OOM " +
			"reclaim until `sshd`, `systemd-journald`, or the init daemon itself gets " +
			"killed. `--oom-score-adj <= -500` stops short of full immunity but still " +
			"preferentially kills unrelated host processes under pressure. If the workload " +
			"genuinely needs resilience, cap memory with `--memory=<limit>` and accept the " +
			"container being killed on overrun; shift the heavy workload to a dedicated " +
			"node instead of rigging OOM scores.",
		Check: checkZC1674,
	})
}

func checkZC1674(node ast.Node) []Violation {
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
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "run" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--oom-kill-disable" {
			return zc1674Hit(cmd, v)
		}
		if strings.HasPrefix(v, "--oom-score-adj=") {
			adj := strings.TrimPrefix(v, "--oom-score-adj=")
			if zc1674Harsh(adj) {
				return zc1674Hit(cmd, v)
			}
			continue
		}
		if v == "--oom-score-adj" {
			idx := i + 2
			if idx >= len(cmd.Arguments) {
				continue
			}
			adj := cmd.Arguments[idx].String()
			if zc1674Harsh(adj) {
				return zc1674Hit(cmd, v+" "+adj)
			}
		}
	}
	return nil
}

func zc1674Harsh(adj string) bool {
	n, err := strconv.Atoi(adj)
	if err != nil {
		return false
	}
	return n <= -500
}

func zc1674Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1674",
		Message: "`" + form + "` shifts OOM pressure onto the rest of the host — cap " +
			"memory with `--memory=<limit>` instead of rigging the OOM score.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1675",
		Title:    "Avoid Bash-only `export -f` / `export -n` — use Zsh `typeset -fx` / `typeset +x`",
		Severity: SeverityInfo,
		Description: "`export -f FUNC` (export a function to child processes) and `export -n " +
			"VAR` (strip the export flag while keeping the value) are Bash-only. Zsh's " +
			"`export` ignores `-f` entirely and prints usage for `-n`, so scripts that " +
			"depend on either silently break under Zsh. The Zsh equivalents are `typeset " +
			"-fx FUNC` for function export (parameter-passing via `$FUNCTIONS` in a " +
			"subshell) and `typeset +x VAR` to drop the export flag. Functions that must " +
			"cross a subshell are usually better handled by `autoload -Uz` from an `fpath` " +
			"directory than by serialisation.",
		Check: checkZC1675,
		Fix:   fixZC1675,
	})
}

// fixZC1675 collapses `export -f` and `export -n` into the Zsh
// equivalents `typeset -fx` and `typeset +x`. Single edit spans the
// command name + flag together, mirroring fixZC1283's `set -o OPT`
// → `setopt OPT` collapse.
var zc1675FlagReplace = map[string]string{
	"-f": "typeset -fx",
	"-n": "typeset +x",
}

func fixZC1675(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || CommandIdentifier(cmd) != "export" {
		return nil
	}
	flag, replace := zc1675FindFlag(cmd)
	if flag == nil {
		return nil
	}
	nameOff, ok := zc1675ExportOffset(source, v)
	if !ok {
		return nil
	}
	flagOff, ok := zc1675FlagOffset(source, flag)
	if !ok {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  flagOff + 2 - nameOff,
		Replace: replace,
	}}
}

func zc1675FindFlag(cmd *ast.SimpleCommand) (ast.Expression, string) {
	for _, arg := range cmd.Arguments {
		if r, hit := zc1675FlagReplace[arg.String()]; hit {
			return arg, r
		}
	}
	return nil, ""
}

func zc1675ExportOffset(source []byte, v Violation) (int, bool) {
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 || off+len("export") > len(source) {
		return 0, false
	}
	if string(source[off:off+len("export")]) != "export" {
		return 0, false
	}
	return off, true
}

func zc1675FlagOffset(source []byte, flag ast.Expression) (int, bool) {
	tok := flag.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 || off+2 > len(source) {
		return 0, false
	}
	lit := string(source[off : off+2])
	if lit != "-f" && lit != "-n" {
		return 0, false
	}
	return off, true
}

func checkZC1675(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "export" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-f":
			return zc1675Hit(cmd, "export -f", "typeset -fx")
		case "-n":
			return zc1675Hit(cmd, "export -n", "typeset +x")
		}
	}
	return nil
}

func zc1675Hit(cmd *ast.SimpleCommand, bad, good string) []Violation {
	return []Violation{{
		KataID:  "ZC1675",
		Message: "`" + bad + "` is Bash-only — use `" + good + "` in Zsh.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1676",
		Title:    "Warn on `helm rollback --force` — recreates in-flight resources, corrupts rolling updates",
		Severity: SeverityWarning,
		Description: "`helm rollback RELEASE N --force` asks Helm to delete and recreate any " +
			"resource that it cannot patch cleanly. If a deployment is mid-rollout, the " +
			"`--force` flag takes out both the old and new ReplicaSets, kicks the pods, " +
			"and forces a cold start — losing in-flight requests and any `PodDisruptionBudget` " +
			"protections. Worse, rolling back to revision N brings back whatever CVEs or " +
			"config regressions the later revisions had already fixed. Pin the target " +
			"revision explicitly, omit `--force`, and gate the rollback behind a change-" +
			"review ticket rather than a shell one-liner.",
		Check: checkZC1676,
	})
}

func checkZC1676(node ast.Node) []Violation {
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
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "rollback" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "--force" {
			return []Violation{{
				KataID: "ZC1676",
				Message: "`helm rollback --force` deletes and recreates unpatched resources — " +
					"loses in-flight traffic and bypasses PodDisruptionBudget. Drop `--force` " +
					"and gate the rollback via change review.",
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
		ID:       "ZC1677",
		Title:    "Warn on `trap 'set -x' DEBUG` — xtrace on every command leaks secrets",
		Severity: SeverityWarning,
		Description: "`trap 'set -x' DEBUG` runs the trap handler before every simple command, " +
			"turning on xtrace for the remainder of the shell. Every subsequent `curl " +
			"-H 'Authorization: Bearer …'`, `mysql -p<password>`, or `aws configure set " +
			"…` then prints its full argv to stderr — commonly into a log file or CI " +
			"artifact. The same antipattern shows up as `set -o xtrace` inside a DEBUG " +
			"trap. Instrument selectively with `typeset -ft FUNC` (Zsh function-level " +
			"xtrace), or add `exec 2>>\"$log\"; set -x` only around the part of the " +
			"script you want traced.",
		Check: checkZC1677,
	})
}

func checkZC1677(node ast.Node) []Violation {
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

	last := cmd.Arguments[len(cmd.Arguments)-1].String()
	if last != "DEBUG" {
		return nil
	}

	handler := strings.Trim(cmd.Arguments[0].String(), "'\"")
	if !strings.Contains(handler, "set -x") && !strings.Contains(handler, "set -o xtrace") {
		return nil
	}

	return []Violation{{
		KataID: "ZC1677",
		Message: "`trap 'set -x' DEBUG` keeps xtrace on after the first command — " +
			"every subsequent argv (passwords, bearer tokens) lands in the log. Trace a " +
			"narrow block with `set -x … set +x` or use `typeset -ft FUNC` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1678",
		Title:    "Error on `borg init --encryption=none` — unencrypted backup repository",
		Severity: SeverityError,
		Description: "`borg init --encryption=none REPO` creates a backup repository without " +
			"client-side encryption or authentication. Anyone with read access to the repo " +
			"gets every file in every archive, and no one can detect silent tampering — " +
			"borg will happily extract a modified chunk. Even for local-only repos the cost " +
			"of authenticated-encryption is tiny; use `--encryption=repokey-blake2` (or " +
			"`--encryption=keyfile-blake2` when you want the key off the server), and store " +
			"the passphrase in `BORG_PASSPHRASE_FILE` pointing at a mode-0400 file.",
		Check: checkZC1678,
	})
}

func checkZC1678(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "borg" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "init" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--encryption=none" {
			return zc1678Hit(cmd)
		}
		if (v == "--encryption" || v == "-e") && i+2 < len(cmd.Arguments) &&
			cmd.Arguments[i+2].String() == "none" {
			return zc1678Hit(cmd)
		}
	}
	return nil
}

func zc1678Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1678",
		Message: "`borg init --encryption=none` leaves archives unauthenticated and " +
			"readable — use `--encryption=repokey-blake2` (or `keyfile-blake2`) and store " +
			"the passphrase in `BORG_PASSPHRASE_FILE`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

var zc1679BroadRoles = map[string]struct{}{
	"roles/owner":                             {},
	"roles/editor":                            {},
	"roles/iam.securityAdmin":                 {},
	"roles/iam.serviceAccountTokenCreator":    {},
	"roles/iam.serviceAccountKeyAdmin":        {},
	"roles/iam.workloadIdentityUser":          {},
	"roles/resourcemanager.organizationAdmin": {},
	"roles/resourcemanager.projectIamAdmin":   {},
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1679",
		Title:    "Error on `gcloud ... add-iam-policy-binding ... --role=roles/owner` — GCP primitive admin",
		Severity: SeverityError,
		Description: "`gcloud projects|folders|organizations add-iam-policy-binding` with the " +
			"primitive roles `roles/owner` or `roles/editor`, or with the IAM-escalation " +
			"roles (`roles/iam.securityAdmin`, `roles/iam.serviceAccountTokenCreator`, " +
			"`roles/iam.serviceAccountKeyAdmin`, `roles/resourcemanager.organizationAdmin`), " +
			"hands the principal the ability to grant themselves any other permission. " +
			"Scripts rarely need that scope; the pattern signals someone papering over a " +
			"permissions error. Grant a specific predefined role (e.g. `roles/compute." +
			"viewer`) or build a custom role with only the `Action`s the workload needs, " +
			"and apply admin changes via Terraform under change review.",
		Check: checkZC1679,
	})
}

func checkZC1679(node ast.Node) []Violation {
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

	hasAdd := false
	var hit string
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "add-iam-policy-binding" {
			hasAdd = true
			continue
		}
		if strings.HasPrefix(v, "--role=") {
			if _, broad := zc1679BroadRoles[strings.TrimPrefix(v, "--role=")]; broad {
				hit = v
			}
		}
	}
	if !hasAdd || hit == "" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1679",
		Message: "`gcloud ... add-iam-policy-binding " + hit + "` grants primitive / IAM-" +
			"admin — use a predefined role with the minimum scope or a custom role, and " +
			"apply admin changes via Terraform.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1680",
		Title:    "Error on `ansible-playbook --vault-password-file=/tmp/...` — world-traversable vault key",
		Severity: SeverityError,
		Description: "The Ansible Vault decryption key lives in the `--vault-password-file` " +
			"path. `/tmp`, `/var/tmp`, and `/dev/shm` are world-traversable: a concurrent " +
			"local user who guesses (or `inotifywait`s for) the filename opens it during " +
			"the playbook run and dumps every secret the vault protects. Keep vault keys " +
			"in a root-owned mode-0400 file under `/etc/ansible/` or `$HOME/.ansible/`, or " +
			"supply the passphrase via a no-echo helper script (`vault-password-client`) " +
			"that fetches from `pass` / `vault kv get`.",
		Check: checkZC1680,
	})
}

func checkZC1680(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ansible-playbook" && ident.Value != "ansible" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "--vault-password-file=") {
			if zc1680Unsafe(strings.TrimPrefix(v, "--vault-password-file=")) {
				return zc1680Hit(cmd)
			}
			continue
		}
		if v == "--vault-password-file" && i+1 < len(cmd.Arguments) {
			if zc1680Unsafe(cmd.Arguments[i+1].String()) {
				return zc1680Hit(cmd)
			}
		}
	}
	return nil
}

func zc1680Unsafe(path string) bool {
	return strings.HasPrefix(path, "/tmp/") ||
		strings.HasPrefix(path, "/var/tmp/") ||
		strings.HasPrefix(path, "/dev/shm/")
}

func zc1680Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1680",
		Message: "`ansible-playbook --vault-password-file` under `/tmp/` / `/var/tmp/` / " +
			"`/dev/shm/` — world-traversable, any local user can race-read it. Store the " +
			"key mode-0400 under `/etc/ansible/` or supply via a `vault-password-client` helper.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1681",
		Title:    "Error on `tar -P` / `--absolute-names` — archive absolute paths, can overwrite host files",
		Severity: SeverityError,
		Description: "By default GNU tar strips the leading `/` from archive member paths so " +
			"that extraction stays under the current directory. `-P` (or the long form " +
			"`--absolute-names`) disables that strip: `tar -xPf evil.tar` happily writes to " +
			"`/etc/cron.d/evil`, `/usr/local/bin/sshd`, or any other absolute path the " +
			"archive mentions. Archives from untrusted sources should never be unpacked " +
			"with `-P`. Drop the flag, extract with `-C <scratch-dir>`, audit the tree, " +
			"then copy files into place with `install` or `cp`.",
		Check: checkZC1681,
	})
}

func checkZC1681(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "tar", "bsdtar", "gtar":
	case "absolute-names":
		// `tar --absolute-names …` parses with `tar` consumed — the trailing
		// name alone is unambiguous evidence of the flag.
		return zc1681Hit(cmd, "--absolute-names")
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-P" || v == "--absolute-names" {
			return zc1681Hit(cmd, v)
		}
		if !strings.HasPrefix(v, "-") || strings.HasPrefix(v, "--") {
			continue
		}
		body := strings.TrimPrefix(v, "-")
		if strings.ContainsRune(body, 'P') {
			return zc1681Hit(cmd, v)
		}
	}
	return nil
}

func zc1681Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1681",
		Message: "`tar " + form + "` keeps absolute paths during extraction — an " +
			"untrusted archive can overwrite `/etc/cron.d`, `/usr/local/bin`, etc. Drop " +
			"the flag and extract with `-C <scratch-dir>` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1682",
		Title:    "Error on `npm install --unsafe-perm` — npm lifecycle scripts keep root privileges",
		Severity: SeverityError,
		Description: "npm normally drops to the UID that owns `package.json` before running " +
			"`preinstall` / `install` / `postinstall` lifecycle scripts. `--unsafe-perm` " +
			"(or `--unsafe-perm=true`) tells npm to skip that drop and run every script as " +
			"the current UID — typically root when the install happens from a provisioning " +
			"script. Any compromised or malicious dependency then executes as root. If a " +
			"native addon truly needs privileges, scope them: drop them into a dedicated " +
			"builder container, or use `sudo -u builduser npm install` from a non-root " +
			"account that already owns `node_modules/`.",
		Check: checkZC1682,
	})
}

func checkZC1682(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "npm" && ident.Value != "yarn" && ident.Value != "pnpm" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--unsafe-perm" || v == "--unsafe-perm=true" {
			return []Violation{{
				KataID: "ZC1682",
				Message: "`" + ident.Value + " " + v + "` keeps root for every lifecycle " +
					"script — a compromised dep executes as root. Build in a dedicated " +
					"builder container or run as a non-root user.",
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
		ID:       "ZC1683",
		Title:    "Error on `npm/yarn/pnpm config set registry http://...` — plaintext package index",
		Severity: SeverityError,
		Description: "Pointing a JavaScript package manager at an `http://` registry disables " +
			"TLS during fetch. Any host on the path (corporate proxy, hotel Wi-Fi, " +
			"compromised CDN) can rewrite tarballs mid-flight; lockfile hashes catch the " +
			"rewrite only if the user locks every dependency before the swap. Even on " +
			"internal networks, pin to `https://` — reach for your own CA via " +
			"`NODE_EXTRA_CA_CERTS` or `registry.cafile` rather than falling back to HTTP.",
		Check: checkZC1683,
	})
}

func checkZC1683(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "npm" && ident.Value != "yarn" && ident.Value != "pnpm" {
		return nil
	}

	if len(cmd.Arguments) < 4 {
		return nil
	}
	if cmd.Arguments[0].String() != "config" || cmd.Arguments[1].String() != "set" {
		return nil
	}
	if cmd.Arguments[2].String() != "registry" {
		return nil
	}
	url := cmd.Arguments[3].String()
	if !strings.HasPrefix(url, "http://") {
		return nil
	}

	return []Violation{{
		KataID: "ZC1683",
		Message: "`" + ident.Value + " config set registry " + url + "` uses plaintext " +
			"HTTP — any proxy / CDN can rewrite tarballs. Use `https://` and a custom " +
			"CA via `NODE_EXTRA_CA_CERTS` if needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1684",
		Title:    "Error on `redis-cli -a PASSWORD` — authentication password in process list",
		Severity: SeverityError,
		Description: "`redis-cli -a <password>` (and the joined form `-aPASSWORD`) puts the " +
			"authentication password in the command line — visible to every user on the " +
			"host through `ps`, `/proc/PID/cmdline`, audit logs, and shell history. redis-" +
			"cli 6.0+ prints a warning to stderr but still connects. Use the " +
			"`REDISCLI_AUTH` environment variable (read automatically by redis-cli), or " +
			"`-askpass` to prompt from TTY; both keep the secret out of the argv tail.",
		Check: checkZC1684,
	})
}

func checkZC1684(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "redis-cli" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-a" && i+1 < len(cmd.Arguments) {
			return zc1684Hit(cmd)
		}
		if strings.HasPrefix(v, "-a") && v != "-a" && !strings.HasPrefix(v, "--") {
			return zc1684Hit(cmd)
		}
	}
	return nil
}

func zc1684Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1684",
		Message: "`redis-cli -a PASSWORD` leaks the password into `ps` / `/proc/PID/cmdline` — " +
			"use `REDISCLI_AUTH` env var or `-askpass` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1685",
		Title:    "Info: `sleep infinity` — container keep-alive pattern that ignores SIGTERM",
		Severity: SeverityInfo,
		Description: "`sleep infinity` is most often used as a container or systemd-unit keep-" +
			"alive. Problem: GNU `sleep` does not install a SIGTERM handler, so when " +
			"`docker stop` / `systemctl stop` sends SIGTERM the process sits unresponsive " +
			"until the grace period expires and SIGKILL lands. The orchestrator reports a " +
			"hung stop, logs look wrong, and any cleanup registered on signal handlers in " +
			"a wrapping shell never runs. Replace with `exec tail -f /dev/null` (signal-" +
			"handles cleanly) or front with `tini` / `dumb-init` when PID 1 must stay.",
		Check: checkZC1685,
		Fix:   fixZC1685,
	})
}

// fixZC1685 rewrites `sleep infinity` to `exec tail -f /dev/null`.
// Single span replacement covers both tokens. Idempotent — a re-run
// sees `exec`, not `sleep`, so the detector won't fire. Defensive
// byte-match guards on both anchors.
func fixZC1685(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sleep" {
		return nil
	}
	if len(cmd.Arguments) != 1 || cmd.Arguments[0].String() != "infinity" {
		return nil
	}
	cmdOff := LineColToByteOffset(source, v.Line, v.Column)
	if cmdOff < 0 || cmdOff+len("sleep") > len(source) {
		return nil
	}
	if string(source[cmdOff:cmdOff+len("sleep")]) != "sleep" {
		return nil
	}
	argTok := cmd.Arguments[0].TokenLiteralNode()
	argOff := LineColToByteOffset(source, argTok.Line, argTok.Column)
	if argOff < 0 || argOff+len("infinity") > len(source) {
		return nil
	}
	if string(source[argOff:argOff+len("infinity")]) != "infinity" {
		return nil
	}
	end := argOff + len("infinity")
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - cmdOff,
		Replace: "exec tail -f /dev/null",
	}}
}

func checkZC1685(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "sleep" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "infinity" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1685",
		Message: "`sleep infinity` does not trap SIGTERM — the orchestrator hangs until " +
			"SIGKILL. Use `exec tail -f /dev/null` or front with `tini` / `dumb-init`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1686",
		Title:    "Warn on `compinit -C` / `compinit -u` — skips / ignores `$fpath` integrity checks",
		Severity: SeverityWarning,
		Description: "Zsh's completion system loads every file from `$fpath` as shell code. " +
			"`compinit` normally warns when an `$fpath` directory (or a file in one) is " +
			"writable by someone other than the current user or root, and skips loading. " +
			"`compinit -C` skips the security check entirely for speed; `compinit -u` " +
			"acknowledges the warning and loads the insecure files anyway. Either way, a " +
			"world-writable entry in `$fpath` becomes an execution primitive for any user " +
			"on the host. Audit `$fpath` with `compaudit`, fix ownership / permissions, " +
			"then run plain `compinit`.",
		Check: checkZC1686,
	})
}

func checkZC1686(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "compinit" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-C" {
			return zc1686Hit(cmd, "-C", "skip-security-check")
		}
		if v == "-u" {
			return zc1686Hit(cmd, "-u", "load-insecure-files")
		}
	}
	return nil
}

func zc1686Hit(cmd *ast.SimpleCommand, flag, what string) []Violation {
	return []Violation{{
		KataID: "ZC1686",
		Message: "`compinit " + flag + "` (" + what + ") loads `$fpath` files that are " +
			"writable by others — any user on the host can inject shell code. Run " +
			"`compaudit`, fix permissions, then `compinit` without the flag.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1687",
		Title:    "Warn on `snap install --classic` / `--devmode` — weakens snap confinement",
		Severity: SeverityWarning,
		Description: "`snap install --classic` drops the AppArmor / cgroup / seccomp sandbox " +
			"entirely — the snap behaves like a normal Debian package with full system " +
			"access. `--devmode` keeps the sandbox wired up but logs violations instead of " +
			"blocking them. Both modes are documented escape hatches for snaps that cannot " +
			"yet fit the strict confinement (IDEs, compilers, some network tooling), but in " +
			"provisioning scripts they usually mean \"I could not be bothered to pick a " +
			"strict snap.\" Find a strict alternative, or install from the distro repository " +
			"with proper AppArmor profiles; if `--classic` is truly required, document the " +
			"specific snap and the interface that needed elevation.",
		Check: checkZC1687,
	})
}

func checkZC1687(node ast.Node) []Violation {
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
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "install" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--classic" {
			return zc1687Hit(cmd, "--classic", "drops AppArmor / cgroup / seccomp sandbox")
		}
		if v == "--devmode" {
			return zc1687Hit(cmd, "--devmode", "logs confinement violations instead of blocking")
		}
	}
	return nil
}

func zc1687Hit(cmd *ast.SimpleCommand, flag, what string) []Violation {
	return []Violation{{
		KataID: "ZC1687",
		Message: "`snap install " + flag + "` " + what + " — find a strict snap or a " +
			"distro-package alternative, or document why this specific snap needs it.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1688",
		Title:    "Warn on `aws s3 sync --delete` — destination objects deleted when source diverges",
		Severity: SeverityWarning,
		Description: "`aws s3 sync SRC DST --delete` removes every object in DST that does not " +
			"exist under SRC. A misspelled SRC, an empty build directory, or a stale " +
			"`cd` turns the sync into a full-bucket wipe with no second confirmation and " +
			"no recovery unless the bucket had versioning enabled. Restrict deletion to " +
			"the prefix that really changed (`aws s3 sync ./build s3://bucket/app/ " +
			"--delete`), add `--dryrun` behind a gate, or enable versioning and MFA-delete " +
			"before running the command from a pipeline.",
		Check: checkZC1688,
	})
}

func checkZC1688(node ast.Node) []Violation {
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

	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "s3" || cmd.Arguments[1].String() != "sync" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		if arg.String() == "--delete" {
			return []Violation{{
				KataID: "ZC1688",
				Message: "`aws s3 sync --delete` wipes DST objects that are missing from " +
					"SRC — a mistyped SRC bulk-deletes the bucket. Scope to the prefix, " +
					"dry-run first, or enable versioning + MFA-delete.",
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
		ID:       "ZC1689",
		Title:    "Error on `borg delete --force` — forced deletion of backup archives or repository",
		Severity: SeverityError,
		Description: "`borg delete --force REPO[::ARCHIVE]` bypasses the confirmation prompt " +
			"and removes the archive (or the whole repository, if ARCHIVE is omitted) in " +
			"one go. Unlike `borg prune`, which keeps a retention ladder and logs what it " +
			"would drop, `--force` deletion leaves nothing to restore from if the target " +
			"was typed wrong. Keep scripts to `borg prune --keep-daily` / `--keep-within` " +
			"with an explicit retention policy and gate any outright `borg delete` behind " +
			"a human `--checkpoint-interval` review.",
		Check: checkZC1689,
	})
}

func checkZC1689(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "borg" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "delete" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "--force" {
			return []Violation{{
				KataID: "ZC1689",
				Message: "`borg delete --force` skips confirmation and can nuke the whole " +
					"repository on a typo — use `borg prune --keep-*` with a retention " +
					"policy, or gate outright deletion behind a manual review.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

var (
	zc1690GitURL = regexp.MustCompile(`^git\+(https?|ssh|file|git)://`)
	zc1690Hash   = regexp.MustCompile(`^[0-9a-f]{7,40}$`)
	zc1690Tag    = regexp.MustCompile(`^v?\d+\.\d+`)
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1690",
		Title:    "Warn on `pip install git+<URL>` without a commit / tag pin",
		Severity: SeverityWarning,
		Description: "`pip install git+https://host/repo[@main]` checks out a moving ref (the " +
			"repository's default branch when no `@` suffix is given, otherwise a branch " +
			"name the attacker can rewrite). Every subsequent install pulls whatever HEAD " +
			"the branch currently points at — no lockfile, no checksum, no reproducibility. " +
			"Pin to a specific commit SHA (`@abc1234…`) or a signed tag (`@v1.2.3`). If a " +
			"proper PyPI release is available, drop the `git+` form entirely and install " +
			"the versioned package.",
		Check: checkZC1690,
	})
}

func checkZC1690(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "pip", "pip3", "pipx", "uv":
	default:
		return nil
	}

	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "install" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if !zc1690GitURL.MatchString(v) {
			continue
		}
		at := strings.LastIndex(v, "@")
		// Skip the `@` that's part of `git+ssh://git@host/...` — locate only
		// the refspec `@` that follows the path.
		if at > 0 && at > strings.LastIndex(v, "/") {
			ref := v[at+1:]
			if zc1690Hash.MatchString(ref) || zc1690Tag.MatchString(ref) {
				continue
			}
		}
		return []Violation{{
			KataID: "ZC1690",
			Message: "`" + ident.Value + " install " + v + "` tracks a moving git ref — " +
				"pin to a commit SHA (`@abc1234…`) or signed tag (`@v1.2.3`), or use the " +
				"PyPI release.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1691",
		Title:    "Warn on `rsync --remove-source-files` — SRC deletion tied to optimistic success",
		Severity: SeverityWarning,
		Description: "`rsync --remove-source-files` deletes each source file once rsync has " +
			"transferred it. The delete is gated on rsync's per-file success, which is " +
			"generous: a remote out-of-disk error after the partial write, a `--chmod` " +
			"rejection, or a flaky network that drops after the data bytes but before " +
			"metadata can still look like success. Couple that with a wrong DST path and " +
			"the source is gone with nothing to recover. Prefer a two-step flow: `rsync " +
			"-a SRC DST` first, verify DST (checksums / file count), then `rm` the source " +
			"explicitly, or use `mv` for local moves.",
		Check: checkZC1691,
	})
}

func checkZC1691(node ast.Node) []Violation {
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

	for _, arg := range cmd.Arguments {
		if arg.String() == "--remove-source-files" {
			return []Violation{{
				KataID: "ZC1691",
				Message: "`rsync --remove-source-files` deletes SRC on optimistic per-file " +
					"success — verify DST after the transfer and `rm` explicitly instead.",
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
		ID:       "ZC1692",
		Title:    "Error on `kexec -e` — jumps into a new kernel without reboot, no audit trail",
		Severity: SeverityError,
		Description: "`kexec -e` transfers control to whatever kernel image is currently " +
			"loaded via `kexec -l` — there is no firmware reboot, no init re-run, no " +
			"chance for PAM / auditd / systemd hooks to record the transition. Malware " +
			"uses it to pivot into a rootkit kernel while the audit log shows no reboot. " +
			"If the intent is a fast reboot, prefer `systemctl kexec` (writes a wtmp entry " +
			"and flushes filesystems), or just `reboot` / `systemctl reboot` and take the " +
			"firmware cost for the audit trail.",
		Check: checkZC1692,
	})
}

func checkZC1692(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kexec" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-e" {
			return []Violation{{
				KataID: "ZC1692",
				Message: "`kexec -e` jumps to a preloaded kernel without firmware reboot " +
					"— wtmp / auditd see nothing. Use `systemctl kexec` or a real " +
					"`systemctl reboot` to keep the audit trail.",
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
		ID:       "ZC1693",
		Title:    "Warn on `ionice -c 1` — real-time I/O class starves every other disk consumer",
		Severity: SeverityWarning,
		Description: "`ionice -c 1` (real-time I/O scheduling class) promotes the child above " +
			"every best-effort (class 2) and idle (class 3) task queued against the same " +
			"device. A busy workload — `rsync`, `dd`, database backup — then blocks sshd " +
			"reads, systemd journal writes, and every other process until it yields, which " +
			"for sequential I/O is effectively never. If the intent is \"fast I/O\", stay on " +
			"class 2 and let CFQ / BFQ handle it; reserve class 1 for latency-critical " +
			"paths launched by a scheduler that knows how to cap duration.",
		Check: checkZC1693,
	})
}

func checkZC1693(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ionice" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c1" {
			return zc1693Hit(cmd)
		}
		if v == "-c" && i+1 < len(cmd.Arguments) && cmd.Arguments[i+1].String() == "1" {
			return zc1693Hit(cmd)
		}
	}
	return nil
}

func zc1693Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1693",
		Message: "`ionice -c 1` puts the child in the real-time I/O class — a long-running " +
			"workload starves sshd / journald / the rest of the host. Stay on class 2.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1694",
		Title:    "Warn on `ssh -A` / `-o ForwardAgent=yes` — remote host can reuse local keys",
		Severity: SeverityWarning,
		Description: "`ssh -A` (and `-o ForwardAgent=yes`) forwards the caller's `SSH_AUTH_SOCK` " +
			"into the remote session. Anyone with root on the remote (and any process " +
			"that shares its uid) can read the socket and impersonate the caller against " +
			"every host the caller's keys unlock. Prefer `ssh -J JUMP HOST` (ProxyJump) " +
			"for multi-hop access — it keeps the keys on the local side — or configure a " +
			"scoped key for the remote task and copy it in with `ssh-copy-id`. Save key-" +
			"forwarding for interactive use on trusted hosts.",
		Check: checkZC1694,
	})
}

func checkZC1694(node ast.Node) []Violation {
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

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-A" {
			return zc1694Hit(cmd, "-A")
		}
		if v == "-oForwardAgent=yes" {
			return zc1694Hit(cmd, "-o ForwardAgent=yes")
		}
		if v == "-o" && i+1 < len(cmd.Arguments) &&
			cmd.Arguments[i+1].String() == "ForwardAgent=yes" {
			return zc1694Hit(cmd, "-o ForwardAgent=yes")
		}
	}
	return nil
}

func zc1694Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1694",
		Message: "`ssh " + form + "` forwards the caller's `SSH_AUTH_SOCK` into the " +
			"remote — any root on that host can reuse the keys. Use `ssh -J jumphost` " +
			"instead, or a scoped key for the remote task.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1695",
		Title:    "Warn on `terraform state rm` / `state push` — surgery on shared state outside plan/apply",
		Severity: SeverityWarning,
		Description: "`terraform state rm RESOURCE` drops the resource from Terraform's " +
			"tracking without touching the real cloud object — the next `terraform apply` " +
			"sees it as newly-created and tries to re-provision, often hitting name-" +
			"collision errors. `terraform state push FILE` replaces the entire remote " +
			"state with a local file, bypassing locking and overwriting any concurrent " +
			"changes. Both commands skirt the usual plan/apply audit trail. Reach for " +
			"`terraform import` / `terraform apply -replace=ADDR` instead, and only run " +
			"`state rm|push` from a reviewed fix-up PR with state backup in hand.",
		Check: checkZC1695,
	})
}

func checkZC1695(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "terraform" && ident.Value != "tofu" && ident.Value != "terragrunt" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "state" {
		return nil
	}
	sub := cmd.Arguments[1].String()
	if sub != "rm" && sub != "push" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1695",
		Message: "`" + ident.Value + " state " + sub + "` mutates shared state outside " +
			"plan/apply — use `terraform import` or `apply -replace=ADDR` instead, and " +
			"review / back up first.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1696",
		Title:    "Warn on `pnpm install --no-frozen-lockfile` / `yarn install --no-immutable` — CI lockfile drift",
		Severity: SeverityWarning,
		Description: "`pnpm install --no-frozen-lockfile` (pnpm) and `yarn install " +
			"--no-immutable` (yarn 4+) tell the package manager that the lockfile is " +
			"merely a suggestion — any dep resolution change since the lockfile was " +
			"written gets picked up silently. Run that from CI and the artifact no longer " +
			"matches the pinned dependency graph reviewers signed off on. Use `pnpm " +
			"install --frozen-lockfile` (the CI default) or `yarn install --immutable`, " +
			"and let lockfile regen happen only from a dev workstation PR.",
		Check: checkZC1696,
	})
}

func checkZC1696(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "pnpm" && ident.Value != "yarn" && ident.Value != "npm" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "--no-frozen-lockfile":
			return zc1696Hit(cmd, v)
		case "--no-immutable":
			return zc1696Hit(cmd, v)
		}
	}
	return nil
}

func zc1696Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1696",
		Message: "`" + form + "` allows the lockfile to drift — the CI artifact no " +
			"longer matches the reviewed dependency graph. Use `--frozen-lockfile` / " +
			"`--immutable` in CI.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1697",
		Title:    "Info: `cryptsetup open --allow-discards` — TRIM pass-through leaks free-sector map",
		Severity: SeverityInfo,
		Description: "`--allow-discards` tells dm-crypt to forward TRIM/DISCARD commands from " +
			"the filesystem to the underlying SSD. The performance and wear-levelling gains " +
			"are real, but so is the side effect: an attacker with raw-device access can " +
			"read the free-sector map and see which blocks are empty — enough to fingerprint " +
			"partition layouts, distinguish encrypted-full-volume from encrypted-sparse-" +
			"content cases, and defeat plausible-deniability scenarios. If the threat model " +
			"includes offline-disk inspection, drop `--allow-discards` and accept the perf " +
			"hit; otherwise keep the flag but state the trade-off in the runbook.",
		Check: checkZC1697,
	})
}

func checkZC1697(node ast.Node) []Violation {
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
		if arg.String() == "--allow-discards" {
			return []Violation{{
				KataID: "ZC1697",
				Message: "`cryptsetup --allow-discards` leaks free-sector layout to anyone " +
					"with raw-device access — drop it if offline-disk inspection is in " +
					"scope, or document the trade-off.",
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
		ID:       "ZC1698",
		Title:    "Warn on `fail2ban-client unban --all` / `stop` — wipes the active brute-force ban list",
		Severity: SeverityWarning,
		Description: "`fail2ban-client unban --all` clears every active ban across every jail; " +
			"`fail2ban-client stop` shuts the service down and flushes its rules. Either " +
			"command restores network access for the exact attacker IPs `fail2ban` has " +
			"already flagged as hostile — usually hundreds of known bots. Target a single " +
			"IP with `fail2ban-client set <jail> unbanip <ip>` or reload a jail with " +
			"`reload <jail>` when you only need to pick up new filter rules.",
		Check: checkZC1698,
	})
}

func checkZC1698(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "fail2ban-client" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	switch cmd.Arguments[0].String() {
	case "stop":
		return zc1698Hit(cmd, "fail2ban-client stop")
	case "unban":
		for _, arg := range cmd.Arguments[1:] {
			if arg.String() == "--all" {
				return zc1698Hit(cmd, "fail2ban-client unban --all")
			}
		}
	}
	return nil
}

func zc1698Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1698",
		Message: "`" + form + "` wipes every active brute-force ban — attacker IPs " +
			"regain access. Target individual IPs with `set <jail> unbanip <ip>` or " +
			"reload a single jail.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1699",
		Title:    "Warn on `kubectl drain --delete-emptydir-data` — pod-local scratch data lost",
		Severity: SeverityWarning,
		Description: "`kubectl drain NODE --delete-emptydir-data` (older alias `--delete-local-" +
			"data`) lets drain evict pods that mount an `emptyDir` volume — the volume is " +
			"deleted along with the pod, destroying any scratch data it held. Production " +
			"clusters use `emptyDir` for caches, write-ahead logs, and scratch state that " +
			"takes hours to rebuild. Confirm the pods on the node tolerate the loss (or " +
			"migrate to a `persistentVolumeClaim`) before adding the flag; otherwise plan " +
			"a controlled drain without it and accept the stuck-drain warning for the " +
			"affected pods.",
		Check: checkZC1699,
	})
}

func checkZC1699(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "drain" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--delete-emptydir-data" || v == "--delete-local-data" {
			return []Violation{{
				KataID: "ZC1699",
				Message: "`kubectl drain " + v + "` deletes `emptyDir` volumes along with the " +
					"evicted pods — caches / WAL / scratch state are lost. Verify tolerance " +
					"or migrate to a PersistentVolumeClaim first.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
