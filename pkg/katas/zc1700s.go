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
		ID:       "ZC1700",
		Title:    "Error on `ldapsearch -w PASSWORD` / `ldapmodify -w PASSWORD` — bind DN password in process list",
		Severity: SeverityError,
		Description: "OpenLDAP client tools (`ldapsearch`, `ldapmodify`, `ldapadd`, `ldapdelete`, " +
			"`ldapmodrdn`, `ldappasswd`, `ldapcompare`) accept the bind password via `-w " +
			"STRING`. Once invoked, the password sits in `/proc/PID/cmdline`, shell " +
			"history, audit records, and any `ps` output — typically granting cn=admin / " +
			"service-account bind over the whole directory. Use `-W` (prompt), `-y " +
			"FILEPATH` (read from a mode-0400 file), or `SASL` auth (`-Y GSSAPI` with " +
			"Kerberos) to keep the secret out of argv.",
		Check: checkZC1700,
	})
}

var zc1700LDAPTools = map[string]struct{}{
	"ldapsearch":  {},
	"ldapmodify":  {},
	"ldapadd":     {},
	"ldapdelete":  {},
	"ldapmodrdn":  {},
	"ldappasswd":  {},
	"ldapcompare": {},
}

func checkZC1700(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if _, ok := zc1700LDAPTools[ident.Value]; !ok {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() != "-w" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		return []Violation{{
			KataID: "ZC1700",
			Message: "`" + ident.Value + " -w PASSWORD` leaks the LDAP bind password into " +
				"`ps` / `/proc/PID/cmdline` — use `-W` to prompt, `-y FILE` for a mode-0400 " +
				"secret file, or SASL (`-Y GSSAPI`).",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1701",
		Title:    "Info: `dpkg -i FILE.deb` installs without automatic signature verification",
		Severity: SeverityInfo,
		Description: "Unlike `apt install`, which verifies package signatures against the apt " +
			"repository's `Release.gpg`, plain `dpkg -i FILE.deb` applies the package with " +
			"no integrity check beyond Debian's own `.deb` format. In a provisioning " +
			"pipeline that downloaded the file over HTTPS from a vendor, that is usually " +
			"fine — the TLS chain vouches for the bytes. In scripts that pick the file up " +
			"from `/tmp`, `/var/tmp`, `/dev/shm`, or a mutable cache, a local user could " +
			"swap the file between download and install. Verify with `sha256sum -c`, " +
			"`debsig-verify`, or `dpkg-sig --verify` before invoking `dpkg -i`.",
		Check: checkZC1701,
	})
}

func checkZC1701(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dpkg" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-i" {
			return []Violation{{
				KataID: "ZC1701",
				Message: "`dpkg -i FILE.deb` runs the package without signature verification — " +
					"`sha256sum -c` or `debsig-verify` the file first, or install via `apt " +
					"install` from a signed repo.",
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
		ID:       "ZC1702",
		Title:    "Warn on `dpkg-reconfigure` without a noninteractive frontend — hangs in CI",
		Severity: SeverityWarning,
		Description: "`dpkg-reconfigure PACKAGE` opens the package's debconf questions in " +
			"whatever frontend the caller's `DEBIAN_FRONTEND` resolves to — typically a " +
			"terminal dialog that blocks until someone presses a key. Inside a non-" +
			"interactive pipeline (Dockerfile, Ansible task, cloud-init) the call hangs " +
			"until the build times out. Pass `-f noninteractive` (or export " +
			"`DEBIAN_FRONTEND=noninteractive` at the top of the script) and accept the " +
			"debconf defaults; pre-seed any non-default answer with `debconf-set-selections`.",
		Check: checkZC1702,
	})
}

func checkZC1702(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dpkg-reconfigure" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "--frontend=") &&
			strings.Contains(v, "noninteractive") {
			return nil
		}
		if (v == "-f" || v == "--frontend") && i+1 < len(cmd.Arguments) &&
			cmd.Arguments[i+1].String() == "noninteractive" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1702",
		Message: "`dpkg-reconfigure` without `-f noninteractive` opens debconf dialogs — " +
			"non-interactive pipelines hang. Pass `-f noninteractive` or export " +
			"`DEBIAN_FRONTEND=noninteractive`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

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

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1704",
		Title:    "Error on `aws ec2 authorize-security-group-ingress --cidr 0.0.0.0/0` — port open to the internet",
		Severity: SeverityError,
		Description: "`aws ec2 authorize-security-group-ingress --cidr 0.0.0.0/0` (or `::/0` for " +
			"IPv6) adds a rule that accepts the specified protocol/port from any source — " +
			"the exact shape shodan, automated login-probers, and every exploit-as-a-" +
			"service customer scans for. Restrict the source to the office CIDR, a VPN " +
			"range, or a named security-group (`--source-group sg-…`). If the workload " +
			"genuinely needs public access, front it with an ALB / API Gateway / CloudFront " +
			"with WAF — not a raw SG rule from a shell script.",
		Check: checkZC1704,
	})
}

func checkZC1704(node ast.Node) []Violation {
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
	if cmd.Arguments[0].String() != "ec2" {
		return nil
	}
	if cmd.Arguments[1].String() != "authorize-security-group-ingress" {
		return nil
	}

	for i, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if v != "--cidr" && v != "--cidr-ip" && v != "--cidr-ipv6" {
			continue
		}
		idx := i + 3
		if idx >= len(cmd.Arguments) {
			continue
		}
		cidr := cmd.Arguments[idx].String()
		if cidr == "0.0.0.0/0" || cidr == "::/0" {
			return []Violation{{
				KataID: "ZC1704",
				Message: "`aws ec2 authorize-security-group-ingress --cidr " + cidr +
					"` opens the port to the entire internet — scope to a known source CIDR " +
					"or `--source-group sg-…`.",
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
		ID:       "ZC1705",
		Title:    "Info: `awk -i inplace` is gawk-only — script breaks on mawk / BSD awk",
		Severity: SeverityInfo,
		Description: "The `inplace` extension that powers `awk -i inplace` ships only with gawk. " +
			"On Alpine (default `mawk`), Debian-busybox, macOS, FreeBSD, NetBSD, OpenBSD, " +
			"or any container image without `gawk` installed the script aborts with " +
			"`fatal: can't open extension 'inplace'`. If portability matters, write through " +
			"a temporary file (`awk … input > tmp && mv tmp input`); if you really do need " +
			"in-place edits in scripts that target gawk only, document the requirement and " +
			"add `command -v gawk >/dev/null` at the top.",
		Check: checkZC1705,
	})
}

func checkZC1705(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "awk" && ident.Value != "gawk" && ident.Value != "mawk" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() != "-i" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		if cmd.Arguments[i+1].String() == "inplace" {
			return []Violation{{
				KataID: "ZC1705",
				Message: "`awk -i inplace` is gawk-only — fails on mawk / BSD awk / busybox " +
					"awk. For portability rewrite as `awk … input > tmp && mv tmp input`.",
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
		ID:       "ZC1706",
		Title:    "Error on `lvresize -L -SIZE` without `-r` — shrink without filesystem resize corrupts data",
		Severity: SeverityError,
		Description: "`lvresize -L -SIZE` (or `--size -SIZE`) shrinks the logical volume by " +
			"SIZE bytes/extents. The filesystem on top still thinks it owns the original " +
			"range; reads beyond the new LV end now return zeros, and the next write " +
			"corrupts metadata. The `-r` (`--resizefs`) flag tells lvresize to call " +
			"`fsadm` (which calls `resize2fs` / `xfs_growfs` / etc.) so the filesystem " +
			"shrinks first. For ext4, always shrink the FS before the LV; for XFS, online " +
			"shrink is impossible — back up, recreate, restore.",
		Check: checkZC1706,
	})
}

func checkZC1706(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "lvresize" && ident.Value != "lvreduce" {
		return nil
	}

	hasResizefs := false
	shrinking := ident.Value == "lvreduce" // lvreduce always shrinks
	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-r" || v == "--resizefs" {
			hasResizefs = true
		}
		if (v == "-L" || v == "--size") && i+1 < len(cmd.Arguments) {
			next := cmd.Arguments[i+1].String()
			if strings.HasPrefix(next, "-") {
				shrinking = true
			}
		}
	}

	if !shrinking || hasResizefs {
		return nil
	}

	return []Violation{{
		KataID: "ZC1706",
		Message: "`" + ident.Value + "` shrinks the LV without `-r` / `--resizefs` — the " +
			"filesystem on top is not shrunk first and writes past the new boundary " +
			"corrupt metadata. Add `-r` (or shrink the FS manually first).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1707",
		Title:    "Warn on `gpg --keyserver hkp://…` — plaintext keyserver fetch",
		Severity: SeverityWarning,
		Description: "`hkp://` is the unencrypted HKP keyserver protocol. A MITM on the path " +
			"(corporate proxy, hotel Wi-Fi, hostile router) can swap key bytes during the " +
			"fetch and `gpg --recv-keys` happily imports the substitute. Use `hkps://" +
			"keys.openpgp.org` (TLS) or fetch the armored key over HTTPS and verify the " +
			"fingerprint out-of-band before `gpg --import`.",
		Check: checkZC1707,
	})
}

func checkZC1707(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "keyserver" {
		// Parser-mangled form: `gpg --keyserver hkp://…` lost `gpg`.
		if len(cmd.Arguments) > 0 && strings.HasPrefix(cmd.Arguments[0].String(), "hkp://") {
			return zc1707Hit(cmd)
		}
		return nil
	}
	if ident.Value != "gpg" && ident.Value != "gpg2" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "--keyserver=hkp://") {
			return zc1707Hit(cmd)
		}
		if v == "--keyserver" && i+1 < len(cmd.Arguments) &&
			strings.HasPrefix(cmd.Arguments[i+1].String(), "hkp://") {
			return zc1707Hit(cmd)
		}
	}
	return nil
}

func zc1707Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1707",
		Message: "`gpg --keyserver hkp://…` is plaintext — a MITM swaps the key bytes. Use " +
			"`hkps://keys.openpgp.org` or fetch over HTTPS and verify the fingerprint.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1708",
		Title:    "Error on `find -L ... -delete` / `-exec rm` — symlink follow into unintended trees",
		Severity: SeverityError,
		Description: "`find -L` follows symlinks during traversal. Combined with `-delete` (or " +
			"`-exec rm`), a symlink under the start path that points outside the intended " +
			"root steers `find` into / `unlink`s files in `/etc`, `/var/lib`, or any other " +
			"directory the symlink target reaches. Drop `-L` (the default `-P` keeps " +
			"symlinks as objects), or restrict the walk with `-xdev`, `-mount`, and an " +
			"explicit `-type f` test. For log-rotation pipes, `logrotate` is safer than a " +
			"`find` one-liner.",
		Check: checkZC1708,
	})
}

func checkZC1708(node ast.Node) []Violation {
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

	hasFollow := false
	hasDestructive := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-L", "--follow":
			hasFollow = true
		case "-delete", "-exec":
			hasDestructive = true
		}
	}

	if !hasFollow || !hasDestructive {
		return nil
	}

	return []Violation{{
		KataID: "ZC1708",
		Message: "`find -L … -delete/-exec` follows symlinks into unintended trees — drop " +
			"`-L`, add `-xdev`, or scope the walk explicitly.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1709",
		Title:    "Error on `htpasswd -b USER PASSWORD` — basic-auth password in process list",
		Severity: SeverityError,
		Description: "`htpasswd -b FILE USER PASSWORD` (batch mode) takes the password as an " +
			"argv slot. The cleartext sits in `/proc/PID/cmdline`, shell history, audit " +
			"records, and any `ps` output. Use `htpasswd -i FILE USER` and pipe the " +
			"secret on stdin (`printf %s \"$pw\" | htpasswd -i FILE USER`), or omit `-b` " +
			"and `-i` so htpasswd prompts on the controlling TTY.",
		Check: checkZC1709,
	})
}

func checkZC1709(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "htpasswd" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if !strings.HasPrefix(v, "-") || strings.HasPrefix(v, "--") {
			continue
		}
		body := strings.TrimPrefix(v, "-")
		if strings.ContainsRune(body, 'b') {
			return []Violation{{
				KataID: "ZC1709",
				Message: "`htpasswd -b USER PASSWORD` puts the password in argv — visible " +
					"via `ps` / `/proc/PID/cmdline`. Use `htpasswd -i FILE USER` with the " +
					"password piped on stdin instead.",
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
		ID:       "ZC1710",
		Title:    "Error on `journalctl --vacuum-size=1` / `--vacuum-time=1s` — journal-wipe pattern",
		Severity: SeverityError,
		Description: "`journalctl --vacuum-size=1` (down to 1 byte / 1K), `--vacuum-time=1s` " +
			"(retain only the last second), or `--vacuum-files=1` (keep one journal file) " +
			"effectively flushes the entire systemd journal. The classic shape after a " +
			"compromise — clear the audit trail before re-enabling logging. Real retention " +
			"belongs in `/etc/systemd/journald.conf` (`SystemMaxUse=`, `MaxRetentionSec=`), " +
			"not in an ad-hoc one-shot. If you genuinely need to bound disk use, set the " +
			"limit to a meaningful value (`--vacuum-time=2weeks`, `--vacuum-size=200M`).",
		Check: checkZC1710,
	})
}

var zc1710VacuumPrefixes = []string{
	"--vacuum-size=",
	"--vacuum-time=",
	"--vacuum-files=",
}

func checkZC1710(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "journalctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		for _, prefix := range zc1710VacuumPrefixes {
			if !strings.HasPrefix(v, prefix) {
				continue
			}
			val := strings.TrimPrefix(v, prefix)
			if !zc1710Aggressive(val) {
				continue
			}
			return []Violation{{
				KataID: "ZC1710",
				Message: "`journalctl " + v + "` flushes the systemd journal — classic " +
					"audit-clear shape. Set retention in `/etc/systemd/journald.conf` " +
					"(`SystemMaxUse=`, `MaxRetentionSec=`) instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

// zc1710Aggressive returns true for vacuum values that effectively wipe the
// journal: `1`, `1B`, `1K`, `1KB`, `1s`, `1m`, `0`, etc.
func zc1710Aggressive(val string) bool {
	switch val {
	case "0", "1":
		return true
	}
	low := strings.ToLower(val)
	switch low {
	case "1b", "1k", "1kb", "1kib", "1s", "1m", "1ms", "1µs", "0s", "0m":
		return true
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1711",
		Title:    "Error on `etcdctl del --prefix \"\"` / `--from-key \"\"` — wipes the entire keyspace",
		Severity: SeverityError,
		Description: "`etcdctl del --prefix KEY` deletes every key under KEY's range. With KEY " +
			"empty (`\"\"` or `\"\\0\"`) the range is `[\"\", \"\\xFF\")` — the whole etcd " +
			"cluster, including kube-apiserver state if etcd is the Kubernetes datastore. " +
			"`--from-key \"\"` has the same effect for the lower-bound form. Restrict the " +
			"prefix to the namespace you actually own (`/app/staging/`), or wrap the call " +
			"with an explicit `etcdctl get --prefix --keys-only` review step.",
		Check: checkZC1711,
	})
}

func checkZC1711(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "etcdctl" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "del" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v != "--prefix" && v != "--from-key" {
			continue
		}
		idx := i + 2
		if idx >= len(cmd.Arguments) {
			continue
		}
		next := cmd.Arguments[idx].String()
		if next == `""` || next == "''" {
			return []Violation{{
				KataID: "ZC1711",
				Message: "`etcdctl del " + v + " \"\"` deletes the entire etcd keyspace " +
					"(including kube-apiserver state) — scope to a specific namespace " +
					"prefix and review with `get --prefix --keys-only` first.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

var zc1712SecretKeys = []string{
	"password", "passwd", "secret", "token",
	"apikey", "api_key", "api-key",
	"accesskey", "access_key", "access-key",
	"privatekey", "private_key", "private-key",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1712",
		Title:    "Error on `vault kv put PATH password=…` — secret value in process list",
		Severity: SeverityError,
		Description: "`vault kv put PATH key=value` (and the older `vault write PATH key=value`) " +
			"put the value on the command line. When the key name screams secret " +
			"(`password`, `secret`, `token`, `apikey`, `access_key`, `private_key`), the " +
			"cleartext shows up in `ps`, `/proc/<pid>/cmdline`, shell history, and the " +
			"audit log of the calling host — exactly the surface Vault is meant to remove. " +
			"Use `key=@path/to/file` to read from disk, `key=-` to take the value on stdin, " +
			"or `vault kv put -mount=secret PATH @secret.json` for a JSON payload.",
		Check: checkZC1712,
	})
}

func checkZC1712(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "vault" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	var start int
	switch cmd.Arguments[0].String() {
	case "write":
		start = 2
	case "kv":
		if len(cmd.Arguments) < 2 || cmd.Arguments[1].String() != "put" {
			return nil
		}
		start = 3
	default:
		return nil
	}
	if start >= len(cmd.Arguments) {
		return nil
	}

	for _, arg := range cmd.Arguments[start:] {
		v := arg.String()
		eq := strings.IndexByte(v, '=')
		if eq <= 0 {
			continue
		}
		key := strings.ToLower(v[:eq])
		val := v[eq+1:]
		if val == "" || val == "-" || strings.HasPrefix(val, "@") {
			continue
		}
		for _, secret := range zc1712SecretKeys {
			if !strings.Contains(key, secret) {
				continue
			}
			return []Violation{{
				KataID: "ZC1712",
				Message: "`vault " + cmd.Arguments[0].String() + " " + v + "` puts the " +
					"secret value in argv — visible to every local user. Use " +
					"`" + key + "=@FILE` or `" + key + "=-` to read from disk / stdin.",
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
		ID:       "ZC1713",
		Title:    "Error on `consul kv delete -recurse /` — wipes the entire Consul KV store",
		Severity: SeverityError,
		Description: "`consul kv delete -recurse PREFIX` removes every key under PREFIX. With " +
			"PREFIX `/` (or an empty string) the command nukes the whole KV store, " +
			"including service-discovery payloads, ACL bootstrap tokens, and any " +
			"application-level config the cluster relies on. Scope the prefix to the app " +
			"namespace (`consul kv delete -recurse /app/staging/`), confirm the keys you " +
			"are about to lose with `consul kv get -recurse -keys`, and snapshot the " +
			"datacenter (`consul snapshot save snap.bin`) before any large delete.",
		Check: checkZC1713,
	})
}

func checkZC1713(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "consul" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "kv" || cmd.Arguments[1].String() != "delete" {
		return nil
	}

	hasRecurse := false
	rootPrefix := false
	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		switch v {
		case "-recurse", "--recurse":
			hasRecurse = true
		case "/", "", `""`, "''":
			rootPrefix = true
		}
	}
	if !hasRecurse || !rootPrefix {
		return nil
	}

	return []Violation{{
		KataID: "ZC1713",
		Message: "`consul kv delete -recurse /` removes the entire KV store — service " +
			"discovery, ACL bootstrap, app config. Scope the prefix and snapshot " +
			"(`consul snapshot save`) first.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1714",
		Title:    "Error on `gh repo delete --yes` / `gh release delete --yes` — bypassed confirmation",
		Severity: SeverityError,
		Description: "`gh repo delete OWNER/REPO --yes` (and `gh release delete TAG --yes`) " +
			"skip the interactive confirmation that protects against typos and broken " +
			"variable expansion. A repository deletion is final — issues, PRs, releases, " +
			"GitHub Actions history, and (for free accounts) any forks against it all " +
			"disappear with no soft-delete window. From a script, run without `--yes` so a " +
			"human reviews the target, or wrap deletion in a manually-triggered workflow " +
			"with explicit input prompts.",
		Check: checkZC1714,
	})
}

func checkZC1714(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "gh" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	scope := cmd.Arguments[0].String()
	if scope != "repo" && scope != "release" {
		return nil
	}
	if cmd.Arguments[1].String() != "delete" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		if arg.String() == "--yes" {
			return []Violation{{
				KataID: "ZC1714",
				Message: "`gh " + scope + " delete --yes` bypasses GitHub's confirmation — " +
					"a typo or stale variable destroys the target with no soft-delete. " +
					"Drop `--yes` so a human confirms.",
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
		ID:       "ZC1715",
		Title:    "Error on `read -p \"prompt\"` — Zsh `-p` reads from coprocess, not a prompt",
		Severity: SeverityError,
		Description: "Bash's `read -p \"Prompt: \" var` prints the prompt before reading. " +
			"Zsh's `read -p` means \"read from the coprocess set up with `coproc`\" — when " +
			"no coprocess exists, `read` errors with `no coprocess` and leaves the variable " +
			"empty, silently breaking the script. The Zsh idiom is `read \"var?Prompt: \"` " +
			"— a `?` after the variable name introduces the prompt string, with the same " +
			"behavior under `-r`, `-s`, etc.",
		Check: checkZC1715,
	})
}

func checkZC1715(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "read" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if len(v) < 2 || v[0] != '-' {
			continue
		}
		// Skip long flags (none of read's short flags need to be considered as long).
		if v[1] == '-' {
			continue
		}
		if !strings.ContainsRune(v[1:], 'p') {
			continue
		}
		return []Violation{{
			KataID: "ZC1715",
			Message: "`read " + v + "` triggers Zsh's coprocess reader, not Bash's prompt — " +
				"the variable stays empty. Use `read \"var?Prompt: \"` (the `?` after the " +
				"variable name introduces the prompt).",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1716",
		Title:    "Use Zsh `$CPUTYPE` / `$MACHTYPE` instead of `uname -m` / `-p`",
		Severity: SeverityStyle,
		Description: "Zsh maintains `$CPUTYPE` (e.g. `x86_64`, `aarch64`) and `$MACHTYPE` " +
			"(the GNU triplet) as built-in parameters. Reading them is a constant-time " +
			"parameter expansion, while `uname -m` / `uname -p` forks an external for the " +
			"same answer. The Zsh values are populated at shell start from the same `uname(2)` " +
			"call, so they stay in lockstep with what `uname` would print.",
		Check: checkZC1716,
	})
}

func checkZC1716(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "uname" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-m" || v == "-p" {
			return []Violation{{
				KataID: "ZC1716",
				Message: "Use Zsh `$CPUTYPE` / `$MACHTYPE` instead of `uname " + v + "` — " +
					"parameter expansion avoids forking an external for an answer Zsh " +
					"already cached at startup.",
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
		ID:       "ZC1717",
		Title:    "Warn on `docker pull/push --disable-content-trust` — bypasses image signature checks",
		Severity: SeverityWarning,
		Description: "When `DOCKER_CONTENT_TRUST=1` is enforced on a host (or set via `/etc/docker/" +
			"daemon.json`), Docker rejects unsigned image pulls and signs every push. The " +
			"`--disable-content-trust` flag overrides that per command: a `pull` accepts a " +
			"replaced or unsigned image into local storage, a `push` lands an unsigned tag in " +
			"the registry where downstream pulls cannot verify provenance. Drop the flag and " +
			"sign the artifact (`docker trust sign IMAGE:TAG`) instead, or scope the bypass " +
			"with a tight Notary signer policy.",
		Check: checkZC1717,
		Fix:   fixZC1717,
	})
}

// fixZC1717 strips the `--disable-content-trust` flag from a `docker
// {pull,push,build,create,run}` invocation. The argument parses as a
// ConcatenatedExpression whose token literal is just the leading `--`,
// so the whitespace-aware token-strip helper from ZC1238 can't span
// the full literal on its own. Scan the source forward from the
// argument's start offset for the literal flag bytes and delete the
// span (plus the leading whitespace, so the surrounding source stays
// byte-identical).
func fixZC1717(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}
	const flag = "--disable-content-trust"
	sawSub := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if !sawSub {
			switch v {
			case "pull", "push", "build", "create", "run":
				sawSub = true
				continue
			}
		}
		if !sawSub || v != flag {
			continue
		}
		tok := arg.TokenLiteralNode()
		anchor := LineColToByteOffset(source, tok.Line, tok.Column)
		if anchor < 0 {
			return nil
		}
		// The lexer reports `--` as a separate token, then the rest
		// of the flag is concatenated into a ConcatenatedExpression.
		// The token's column lands on the second `-` for the
		// `DASHDASH` token type, so search a 2-byte window around the
		// anchor for the full flag literal in source bytes.
		off := -1
		for _, delta := range []int{-1, 0, 1} {
			cand := anchor + delta
			if cand < 0 || cand+len(flag) > len(source) {
				continue
			}
			if string(source[cand:cand+len(flag)]) == flag {
				off = cand
				break
			}
		}
		if off < 0 {
			return nil
		}
		start := off
		for start > 0 && (source[start-1] == ' ' || source[start-1] == '\t') {
			start--
		}
		end := off + len(flag)
		startLine, startCol := offsetLineColZC1717(source, start)
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
	return nil
}

func offsetLineColZC1717(source []byte, offset int) (int, int) {
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

func checkZC1717(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	var sub string
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if sub == "" {
			switch v {
			case "pull", "push", "build", "create", "run":
				sub = v
				continue
			}
		}
		if sub != "" && v == "--disable-content-trust" {
			return []Violation{{
				KataID: "ZC1717",
				Message: "`docker " + sub + " --disable-content-trust` overrides " +
					"`DOCKER_CONTENT_TRUST=1` — unsigned image moves into the registry " +
					"or local store. Sign the artifact (`docker trust sign`) instead.",
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
		ID:       "ZC1718",
		Title:    "Error on `gh secret set --body SECRET` / `-b SECRET` — secret in process list",
		Severity: SeverityError,
		Description: "`gh secret set NAME --body VALUE` (or `-b VALUE`, `--body=VALUE`) puts the " +
			"secret on the command line. The cleartext appears in `ps`, `/proc/<pid>/cmdline`, " +
			"shell history, and the audit log of the host running `gh`. Pipe the value via " +
			"stdin (`gh secret set NAME < file`, `printf %s \"$SECRET\" | gh secret set NAME " +
			"--body -`) or use `--body-file PATH` so the value never lands in argv.",
		Check: checkZC1718,
	})
}

func checkZC1718(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gh" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "secret" || cmd.Arguments[1].String() != "set" {
		return nil
	}

	prevBody := false
	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if prevBody {
			if v == "-" {
				return nil
			}
			return zc1718Hit(cmd, "--body "+v)
		}
		switch {
		case v == "--body" || v == "-b":
			prevBody = true
		case strings.HasPrefix(v, "--body="):
			val := strings.TrimPrefix(v, "--body=")
			if val == "-" {
				return nil
			}
			return zc1718Hit(cmd, v)
		}
	}
	return nil
}

func zc1718Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1718",
		Message: "`gh secret set ... " + what + "` puts the secret in argv — visible in " +
			"`ps`, `/proc`, history. Use `--body-file PATH` or pipe via stdin " +
			"(`... --body -` with `printf %s \"$SECRET\" |`).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1719",
		Title:    "Warn on `git filter-branch` — deprecated since Git 2.24, use `git filter-repo`",
		Severity: SeverityWarning,
		Description: "`git filter-branch` is deprecated as of Git 2.24; its manpage opens with " +
			"\"WARNING: this command is deprecated\" and points users at `git filter-repo`. " +
			"`filter-branch` is single-process slow, mishandles common cases (tag rewrites, " +
			"refs/notes/*, signed commits), and leaves orphaned objects behind. The modern " +
			"replacement is `git filter-repo` (separate package; `apt/brew install git-" +
			"filter-repo`) — much faster, safer defaults, and what GitHub / GitLab guidance " +
			"recommends for secret-removal rewrites.",
		Check: checkZC1719,
	})
}

func checkZC1719(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	if cmd.Arguments[0].String() != "filter-branch" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1719",
		Message: "`git filter-branch` is deprecated (Git 2.24+) and its manpage redirects to " +
			"`git filter-repo`. Use that instead — faster, safer defaults, no orphaned " +
			"objects.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1720",
		Title:    "Use Zsh `$COLUMNS` / `$LINES` instead of `tput cols` / `tput lines`",
		Severity: SeverityStyle,
		Description: "Zsh tracks the terminal width and height in `$COLUMNS` and `$LINES`, " +
			"updated automatically on `SIGWINCH`. Reading them is a constant-time " +
			"parameter expansion, while `tput cols` / `tput lines` forks the terminfo " +
			"helper on every call. Use the parameters; reach for `tput` only for terminfo " +
			"queries Zsh does not surface as parameters.",
		Check: checkZC1720,
	})
}

func checkZC1720(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tput" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "cols" || v == "lines" {
			repl := "$COLUMNS"
			if v == "lines" {
				repl = "$LINES"
			}
			return []Violation{{
				KataID: "ZC1720",
				Message: "Use `" + repl + "` instead of `tput " + v + "` — Zsh keeps the " +
					"terminal size in parameters, no fork needed.",
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
		ID:       "ZC1721",
		Title:    "Error on `chmod NNN /dev/<node>` — world-writable device node is local privilege escalation",
		Severity: SeverityError,
		Description: "Granting world-write to a device node hands every local user a primitive: " +
			"`/dev/kvm` becomes a host-root VM-exit gadget, `/dev/uinput` lets any user inject " +
			"keystrokes into the active session, `/dev/loop-control` forges loop devices, " +
			"`/dev/dri/cardN` opens GPU shaders for code-exec, `/dev/mem` / `/dev/kmem` (where " +
			"still permitted) leak kernel state. Keep the kernel-managed default permissions; " +
			"if userspace needs access, add a udev rule that grants it to a specific group, " +
			"never `666` to the world.",
		Check: checkZC1721,
	})
}

func checkZC1721(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chmod" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	var mode, target string
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "/dev/") {
			target = v
			continue
		}
		if mode == "" && zc1671WorldWritable(v) {
			mode = v
		}
	}

	if mode == "" || target == "" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1721",
		Message: "`chmod " + mode + " " + target + "` opens a kernel device node to every " +
			"local user — privilege-escalation surface. Use a udev rule that grants the " +
			"specific group access instead of world-write.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1722",
		Title:    "Warn on `ssh-keyscan HOST >> known_hosts` — TOFU bypass, blind-trust new host key",
		Severity: SeverityWarning,
		Description: "`ssh-keyscan` fetches whatever host key the remote serves on its first reply. " +
			"Appending the result straight to `known_hosts` is the exact step the host-key " +
			"check is meant to defend against: a man-in-the-middle on first contact wins " +
			"permanently. Pin the expected fingerprint via a side channel (vendor docs, prior " +
			"verified contact) and assert it matches `ssh-keyscan HOST | ssh-keygen -lf -` " +
			"before the append.",
		Check: checkZC1722,
	})
}

func checkZC1722(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ssh-keyscan" {
		return nil
	}

	prevRedir := ""
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevRedir != "" {
			if strings.Contains(v, "known_hosts") {
				return []Violation{{
					KataID: "ZC1722",
					Message: "`ssh-keyscan ... " + prevRedir + " " + v + "` accepts the " +
						"first-served host key without verifying its fingerprint. Pipe " +
						"to `ssh-keygen -lf -` and assert the fingerprint first.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
			prevRedir = ""
			continue
		}
		if v == ">>" || v == ">" {
			prevRedir = v
		}
	}
	return nil
}

var zc1723DeleteFlags = map[string]bool{
	"--delete-secret-keys":            true,
	"--delete-secret-and-public-keys": true,
	"--delete-keys":                   true,
	"--delete-key":                    true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1723",
		Title:    "Error on `gpg --delete-secret-keys` / `--delete-key` — irreversible key destruction",
		Severity: SeverityError,
		Description: "GPG key deletion is permanent. Once `--delete-secret-keys`, " +
			"`--delete-secret-and-public-keys`, `--delete-keys`, or `--delete-key` removes " +
			"the keyring entry there is no recovery short of a separate backup or off-card " +
			"reimport. Combined with `--batch --yes`, the confirmation prompt is bypassed " +
			"and a single accidental KEYID resolves to a one-shot wipe. Export the key " +
			"first (`gpg --export-secret-keys --armor KEYID > backup.asc`, store offline) " +
			"and never pair the delete flag with `--batch --yes` in automation.",
		Check: checkZC1723,
	})
}

func checkZC1723(node ast.Node) []Violation {
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
		if zc1723DeleteFlags[v] {
			line, col := FlagArgPosition(cmd, zc1723DeleteFlags)
			return []Violation{{
				KataID: "ZC1723",
				Message: "`gpg " + v + "` permanently destroys keyring entries — no recovery " +
					"without a separate backup. Export with `gpg --export-secret-keys --armor " +
					"KEYID` first; never pair this flag with `--batch --yes`.",
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
		ID:       "ZC1724",
		Title:    "Warn on `pacman -Sy <pkg>` — partial upgrade, breaks dependency closure",
		Severity: SeverityWarning,
		Description: "Arch Linux is rolling-release on the invariant that the local package " +
			"database and the installed package set move together. `pacman -Sy <pkg>` " +
			"refreshes the database and installs ONE package against the new metadata while " +
			"every other installed package stays at its old version. The new package's " +
			"dependency closure pulls libraries newer than what the rest of the system has, " +
			"leaving a half-upgraded state that often manifests as `error while loading " +
			"shared libraries`. Run a full `pacman -Syu` first, then install (`pacman -S " +
			"<pkg>`); for CI use `pacman -Syu --noconfirm <pkg>` so the upgrade and install " +
			"are atomic.",
		Check: checkZC1724,
	})
}

func checkZC1724(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "pacman" {
		return nil
	}

	hasSyNoU := false
	hasPkg := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "-") && len(v) >= 2 {
			letters := v[1:]
			// Must contain both 'S' and 'y' but not 'u' (would be -Syu).
			if strings.Contains(letters, "S") && strings.Contains(letters, "y") && !strings.Contains(letters, "u") {
				hasSyNoU = true
			}
			continue
		}
		if v != "" {
			hasPkg = true
		}
	}

	if !hasSyNoU || !hasPkg {
		return nil
	}

	return []Violation{{
		KataID: "ZC1724",
		Message: "`pacman -Sy <pkg>` is a partial-upgrade footgun — refresh the DB but " +
			"install only one package against the newer metadata. Use `pacman -Syu` " +
			"first, then `pacman -S <pkg>` (or `pacman -Syu --noconfirm <pkg>` to keep " +
			"it atomic).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1725",
		Title:    "Error on `cargo --token TOKEN` / `npm --otp CODE` — registry credential in process list",
		Severity: SeverityError,
		Description: "`cargo publish --token TOKEN` (and `cargo login`, `cargo owner`, `cargo " +
			"yank`) puts the crates.io API token in argv — visible in `ps`, `/proc/<pid>/" +
			"cmdline`, shell history, and CI logs. `npm publish --otp CODE` leaks the " +
			"one-time code the same way. Use environment variables (`CARGO_REGISTRY_TOKEN`, " +
			"`NPM_TOKEN`) or pipe via stdin (`cargo login --token -` reads from stdin), and " +
			"source credentials from a secrets manager instead of the command line.",
		Check: checkZC1725,
	})
}

func checkZC1725(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var flag, tool string
	switch ident.Value {
	case "cargo":
		flag = "--token"
		tool = "cargo"
	case "npm", "yarn", "pnpm":
		flag = "--otp"
		tool = ident.Value
	default:
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	// First arg should be a relevant subcommand.
	sub := cmd.Arguments[0].String()
	switch tool {
	case "cargo":
		switch sub {
		case "publish", "login", "owner", "yank":
		default:
			return nil
		}
	default:
		switch sub {
		case "publish", "adduser", "login":
		default:
			return nil
		}
	}

	prevFlag := false
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if prevFlag {
			if v == "-" {
				return nil
			}
			return zc1725Hit(cmd, tool, sub, flag+" "+v)
		}
		switch {
		case v == flag:
			prevFlag = true
		case strings.HasPrefix(v, flag+"="):
			val := strings.TrimPrefix(v, flag+"=")
			if val == "-" {
				return nil
			}
			return zc1725Hit(cmd, tool, sub, v)
		}
	}
	return nil
}

func zc1725Hit(cmd *ast.SimpleCommand, tool, sub, what string) []Violation {
	return []Violation{{
		KataID: "ZC1725",
		Message: "`" + tool + " " + sub + " " + what + "` puts the credential in argv — " +
			"visible in `ps`, `/proc`, history. Pipe via stdin (`--token -`) or use env " +
			"vars like `CARGO_REGISTRY_TOKEN` / `NPM_TOKEN`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1726",
		Title:    "Error on `gcloud ... delete --quiet` — silent destruction of GCP resources",
		Severity: SeverityError,
		Description: "`gcloud` accepts `--quiet` (`-q`) globally to suppress every confirmation " +
			"prompt. Combined with `delete` on projects, SQL instances, GKE clusters, " +
			"compute VMs, secrets, or storage buckets, a single misresolved variable wipes " +
			"the resource with no human-in-the-loop. Project deletion has a 30-day soft " +
			"window but compute disks, secrets, and BigQuery tables are gone immediately. " +
			"Drop `--quiet` from delete commands or route the bulk-destroy through a " +
			"Terraform plan that surfaces the diff for review.",
		Check: checkZC1726,
	})
}

func checkZC1726(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gcloud" {
		return nil
	}

	hasDelete, hasQuiet := false, false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "delete":
			hasDelete = true
		case "--quiet", "-q":
			hasQuiet = true
		}
	}
	if !hasDelete || !hasQuiet {
		return nil
	}

	return []Violation{{
		KataID: "ZC1726",
		Message: "`gcloud ... delete --quiet` skips confirmation — a wrong argument " +
			"wipes the resource (compute disks, secrets, BigQuery tables have no soft-" +
			"delete). Drop `--quiet` or destroy through a Terraform plan with review.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1727",
		Title:    "Error on `curl/wget --proxy http://USER:PASS@HOST` — proxy credentials in argv",
		Severity: SeverityError,
		Description: "Embedding the proxy username and password in the URL passed to `--proxy` " +
			"(curl), `-x` (curl short form), or `--proxy-password=` (wget) lands the " +
			"credential in argv — visible in `ps`, `/proc/<pid>/cmdline`, shell history, " +
			"and CI logs. Configure the proxy through `~/.curlrc` / `~/.netrc` (chmod 600) " +
			"for curl, or `~/.wgetrc` for wget, so the secret never reaches the command " +
			"line.",
		Check: checkZC1727,
	})
}

func checkZC1727(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "curl":
		return zc1727Curl(cmd)
	case "wget":
		return zc1727Wget(cmd)
	}
	return nil
}

func zc1727Curl(cmd *ast.SimpleCommand) []Violation {
	prevProxy := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevProxy {
			if zc1727URLHasCreds(v) {
				return zc1727Hit(cmd, "curl --proxy "+v)
			}
			prevProxy = false
			continue
		}
		switch {
		case v == "--proxy" || v == "-x":
			prevProxy = true
		case strings.HasPrefix(v, "--proxy="):
			val := strings.TrimPrefix(v, "--proxy=")
			if zc1727URLHasCreds(val) {
				return zc1727Hit(cmd, "curl "+v)
			}
		case strings.HasPrefix(v, "-x"):
			val := v[2:]
			if zc1727URLHasCreds(val) {
				return zc1727Hit(cmd, "curl "+v)
			}
		}
	}
	return nil
}

func zc1727Wget(cmd *ast.SimpleCommand) []Violation {
	prevPwd := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevPwd {
			return zc1727Hit(cmd, "wget --proxy-password "+v)
		}
		switch {
		case v == "--proxy-password":
			prevPwd = true
		case strings.HasPrefix(v, "--proxy-password="):
			return zc1727Hit(cmd, "wget "+v)
		}
	}
	return nil
}

// zc1727URLHasCreds returns true when the URL contains a `userinfo` portion
// (text between `://` and the next `@` before any `/`).
func zc1727URLHasCreds(url string) bool {
	scheme := strings.Index(url, "://")
	if scheme < 0 {
		return false
	}
	rest := url[scheme+3:]
	at := strings.Index(rest, "@")
	if at < 0 {
		return false
	}
	if slash := strings.Index(rest, "/"); slash >= 0 && slash < at {
		return false
	}
	return true
}

func zc1727Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1727",
		Message: "`" + what + "` puts proxy credentials in argv — visible in `ps`, " +
			"`/proc`, history. Move them into `~/.curlrc` / `~/.netrc` (chmod 600) or " +
			"`~/.wgetrc`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1728",
		Title:    "Error on `pip install --index-url http://...` — plaintext index allows MITM",
		Severity: SeverityError,
		Description: "`pip install --index-url http://...`, `--extra-index-url http://...`, " +
			"and `-i http://...` tell pip to fetch packages over plaintext HTTP. Any " +
			"network-position attacker (open Wi-Fi, hostile transit, MITM proxy) can " +
			"replace package metadata or wheel contents in flight — direct code execution " +
			"on the install host. Switch to `https://`, or on internal networks terminate " +
			"TLS at the mirror and only configure the `https://` URL.",
		Check: checkZC1728,
	})
}

func checkZC1728(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "pip", "pip3", "pip2":
	default:
		return nil
	}

	prevURL := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevURL {
			if zc1728PlainHTTP(v) {
				return zc1728Hit(cmd, v)
			}
			prevURL = false
			continue
		}
		switch {
		case v == "--index-url" || v == "--extra-index-url" || v == "-i":
			prevURL = true
		case strings.HasPrefix(v, "--index-url="), strings.HasPrefix(v, "--extra-index-url="):
			eq := strings.IndexByte(v, '=')
			if zc1728PlainHTTP(v[eq+1:]) {
				return zc1728Hit(cmd, v)
			}
		}
	}
	return nil
}

func zc1728PlainHTTP(url string) bool {
	return strings.HasPrefix(url, "http://")
}

func zc1728Hit(cmd *ast.SimpleCommand, url string) []Violation {
	return []Violation{{
		KataID: "ZC1728",
		Message: "`pip install --index-url " + url + "` fetches packages over plaintext " +
			"HTTP — any MITM swaps the wheel for code execution on the host. Use " +
			"`https://`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1729",
		Title:    "Error on `ip route flush all` / `ip route del default` — script loses network connectivity",
		Severity: SeverityError,
		Description: "`ip route flush all` (or `flush table main`) wipes every routing entry, " +
			"including the default gateway. `ip route del default` removes only the default " +
			"route — same outcome. The remote SSH session that just ran the command can " +
			"no longer talk to the host, and any subsequent step that needs the network " +
			"hangs until manual console intervention. Scope the flush (`flush dev <iface>`, " +
			"`flush scope link`) or use `ip route replace default via <gw>` so the new " +
			"route is in place before the old one disappears.",
		Check: checkZC1729,
	})
}

func checkZC1729(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ip" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, arg := range cmd.Arguments {
		v := arg.String()
		// Skip leading short flags like `-4`, `-6`, `-s`, `-d`.
		if len(args) == 0 && strings.HasPrefix(v, "-") {
			continue
		}
		args = append(args, v)
	}

	if len(args) < 3 || args[0] != "route" {
		return nil
	}

	switch args[1] {
	case "flush":
		if args[2] == "all" || args[2] == "table" {
			return zc1729Hit(cmd, "ip route flush "+args[2])
		}
	case "del", "delete":
		if args[2] == "default" {
			return zc1729Hit(cmd, "ip route "+args[1]+" default")
		}
	}
	return nil
}

func zc1729Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1729",
		Message: "`" + what + "` removes the default gateway — the SSH session that " +
			"just ran it loses connectivity. Scope the flush (`flush dev <iface>`) or " +
			"use `ip route replace default via <gw>` so the new route is in place " +
			"first.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1730",
		Title:    "Warn on `brew install --HEAD <pkg>` — pulls upstream HEAD, no version stability",
		Severity: SeverityWarning,
		Description: "`brew install --HEAD <pkg>` (also `reinstall --HEAD`, `upgrade --HEAD`) " +
			"builds the formula from the upstream source repository's HEAD branch. The " +
			"build is unrepeatable — every run pulls a different commit — and any " +
			"compromised upstream commit lands directly on the install host. Pin to a " +
			"stable release of the formula, or if HEAD is genuinely required, vendor the " +
			"build into a private tap that fixes a specific revision.",
		Check: checkZC1730,
	})
}

func checkZC1730(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "brew" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	switch cmd.Arguments[0].String() {
	case "install", "reinstall", "upgrade":
	default:
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--HEAD" {
			return []Violation{{
				KataID: "ZC1730",
				Message: "`brew " + cmd.Arguments[0].String() + " --HEAD` builds from " +
					"upstream HEAD — every run pulls a different commit. Pin to a " +
					"stable formula release or vendor a private tap with a fixed " +
					"revision.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

var zc1731SecretKeys = []string{
	"password", "passwd", "pwd",
	"secret", "token",
	"apikey", "api_key", "api-key",
	"accesskey", "access_key", "access-key",
	"privatekey", "private_key", "private-key",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1731",
		Title:    "Error on `curl -d 'password=…'` / `wget --post-data='token=…'` — secret in argv",
		Severity: SeverityError,
		Description: "`curl -d` / `--data` / `--data-raw` / `--data-urlencode` and `wget " +
			"--post-data` / `--body-data` put the POST body in argv — visible in `ps`, " +
			"`/proc/<pid>/cmdline`, shell history, and CI logs. When the body contains a " +
			"credential-looking key (`password`, `secret`, `token`, `apikey`, `access_key`, " +
			"`private_key`), the secret leaks the same way an inline `-u user:pass` would. " +
			"Read the value from a file (`curl --data @secret.txt URL`, `--data-binary @-` " +
			"piped from a secrets store) so the secret never reaches the command line.",
		Check: checkZC1731,
	})
}

func checkZC1731(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var dataFlags map[string]bool
	switch ident.Value {
	case "curl":
		dataFlags = map[string]bool{
			"-d":               true,
			"--data":           true,
			"--data-raw":       true,
			"--data-urlencode": true,
			"--data-binary":    true,
		}
	case "wget":
		dataFlags = map[string]bool{
			"--post-data": true,
			"--body-data": true,
		}
	default:
		return nil
	}

	prevFlag := ""
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevFlag != "" {
			if hit := zc1731MatchSecret(v); hit != "" {
				return zc1731Hit(cmd, ident.Value, prevFlag+" "+v, hit)
			}
			prevFlag = ""
			continue
		}
		if dataFlags[v] {
			prevFlag = v
			continue
		}
		// Joined `--data=key=value` form (curl long flags).
		if eq := strings.IndexByte(v, '='); eq > 0 {
			flag := v[:eq]
			if dataFlags[flag] {
				if hit := zc1731MatchSecret(v[eq+1:]); hit != "" {
					return zc1731Hit(cmd, ident.Value, v, hit)
				}
			}
		}
	}
	return nil
}

func zc1731MatchSecret(value string) string {
	body := strings.Trim(value, "'\"")
	if body == "" {
		return ""
	}
	// File reference (curl `@FILE`) or stdin sentinel — safe.
	if body[0] == '@' || body == "-" {
		return ""
	}
	for _, pair := range strings.Split(body, "&") {
		eq := strings.IndexByte(pair, '=')
		if eq <= 0 {
			continue
		}
		key := strings.ToLower(pair[:eq])
		val := pair[eq+1:]
		if val == "" {
			continue
		}
		for _, secret := range zc1731SecretKeys {
			if strings.Contains(key, secret) {
				return key
			}
		}
	}
	return ""
}

func zc1731Hit(cmd *ast.SimpleCommand, tool, flagPart, key string) []Violation {
	return []Violation{{
		KataID: "ZC1731",
		Message: "`" + tool + " " + flagPart + "` puts secret-keyed POST body (`" + key +
			"=…`) in argv — visible in `ps`, `/proc`, history. Read the value from a " +
			"file with `--data @PATH` or `--data-binary @-` piped from a secrets " +
			"store.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

var zc1732BroadFilesystems = map[string]bool{
	"--filesystem=host":     true,
	"--filesystem=host:rw":  true,
	"--filesystem=home":     true,
	"--filesystem=home:rw":  true,
	"--filesystem=/":        true,
	"--filesystem=/:rw":     true,
	"--filesystem=host-os":  true,
	"--filesystem=host-etc": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1732",
		Title:    "Warn on `flatpak override --filesystem=host` — removes Flatpak sandbox isolation",
		Severity: SeverityWarning,
		Description: "Flatpak's primary security guarantee is filesystem sandboxing — apps see " +
			"only their own data plus paths the user explicitly grants via portals. " +
			"`flatpak override --filesystem=host` (also `host-os`, `host-etc`, `home`, `/`) " +
			"persistently grants the app unrestricted read/write to the host filesystem at " +
			"every subsequent run. Same risk applies to `flatpak run --filesystem=host`. " +
			"Grant the specific subdirectory the app actually needs (`--filesystem=" +
			"~/Documents:ro`) or rely on Filesystem portals so the user picks paths " +
			"interactively per session.",
		Check: checkZC1732,
	})
}

func checkZC1732(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "flatpak" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	switch cmd.Arguments[0].String() {
	case "override", "run":
	default:
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if zc1732BroadFilesystems[v] {
			return []Violation{{
				KataID: "ZC1732",
				Message: "`flatpak " + cmd.Arguments[0].String() + " " + v + "` removes " +
					"the Flatpak sandbox — the app gets unrestricted host-filesystem " +
					"access. Grant a specific subdirectory (e.g. " +
					"`--filesystem=~/Documents:ro`) or use Filesystem portals.",
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
		ID:       "ZC1733",
		Title:    "Error on `docker plugin install --grant-all-permissions` — accepts every requested cap",
		Severity: SeverityError,
		Description: "Docker plugins run as root with whatever privileges they ask for at install " +
			"time — host networking, `/dev/*` mounts, arbitrary capability grants. The " +
			"interactive prompt enumerates each request so the operator can refuse anything " +
			"unexpected. `--grant-all-permissions` skips the prompt and accepts the whole " +
			"list, so a compromised plugin author or a typo-squatted name owns the host " +
			"on first install. Install plugins by name, walk the prompt manually, then pin " +
			"the tag (`@sha256:...`) once vetted.",
		Check: checkZC1733,
	})
}

func checkZC1733(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "plugin" || cmd.Arguments[1].String() != "install" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		if arg.String() == "--grant-all-permissions" {
			return []Violation{{
				KataID: "ZC1733",
				Message: "`docker plugin install --grant-all-permissions` accepts every " +
					"capability the plugin requests — root-equivalent on the host. Walk " +
					"the interactive prompt manually and pin the digest once vetted.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

var zc1734IdentityFiles = map[string]bool{
	"/etc/passwd":  true,
	"/etc/shadow":  true,
	"/etc/group":   true,
	"/etc/gshadow": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1734",
		Title:    "Error on `cp/mv/tee` overwriting `/etc/passwd|shadow|group|gshadow`",
		Severity: SeverityError,
		Description: "The user-identity files are managed by `useradd` / `usermod` / `vipw` / " +
			"`vigr`, which take a file lock and keep `passwd` / `shadow` (and `group` / " +
			"`gshadow`) in sync. Replacing them with `cp`, `mv`, `tee`, or a redirect " +
			"(`echo … > /etc/passwd`) bypasses the lock: concurrent edits race, malformed " +
			"entries lock the whole system out, and the shadow file ends up pointing at " +
			"users that no longer exist. Use `vipw -e` / `vigr -e` to edit, or `useradd` " +
			"/ `usermod` / `passwd` to mutate one entry at a time.",
		Check: checkZC1734,
	})
}

func checkZC1734(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "cp", "mv", "tee", "install", "dd":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if zc1734IdentityFiles[v] {
				return zc1734Hit(cmd, ident.Value+" "+v)
			}
		}
	}

	// Redirect form: any command whose args contain `>` or `>>` followed by an
	// identity file path.
	prevRedir := ""
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevRedir != "" {
			if zc1734IdentityFiles[v] {
				return zc1734Hit(cmd, prevRedir+" "+v)
			}
			prevRedir = ""
			continue
		}
		if v == ">" || v == ">>" {
			prevRedir = v
		}
	}
	return nil
}

func zc1734Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1734",
		Message: "`" + what + "` bypasses the lock that `vipw` / `vigr` / `useradd` use " +
			"on the user-identity files. Edit through `vipw -e` / `vigr -e`, or mutate " +
			"a single entry with `useradd` / `usermod` / `passwd`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1735",
		Title:    "Error on `efibootmgr -B` — deletes UEFI boot entry, may brick boot",
		Severity: SeverityError,
		Description: "`efibootmgr -B` deletes the currently-selected UEFI boot entry; combined " +
			"with `-b BOOTNUM` it removes the specific entry instead. If that entry was " +
			"the only viable bootloader (or the firmware's removable-media fallback is " +
			"not present), the next reboot drops into the UEFI shell or picks an " +
			"unexpected device — recovery needs console access. Run `efibootmgr -v` first " +
			"to inspect `BootOrder`, ensure a fallback (`/EFI/BOOT/BOOTX64.EFI`) is in " +
			"place, and prefer `efibootmgr -o NEW,ORDER` to demote rather than delete.",
		Check: checkZC1735,
	})
}

func checkZC1735(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "efibootmgr" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-B" {
			return []Violation{{
				KataID: "ZC1735",
				Message: "`efibootmgr -B` deletes a UEFI boot entry — wrong BOOTNUM (or " +
					"missing fallback) leaves the box at the UEFI shell on next reboot. " +
					"Inspect `efibootmgr -v` first; demote via `-o NEW,ORDER` instead " +
					"of deleting.",
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
		ID:       "ZC1736",
		Title:    "Error on `pulumi destroy --yes` / `up --yes` — silent infra mutation in CI",
		Severity: SeverityError,
		Description: "`pulumi destroy --yes` (or `-y`) skips the preview-and-confirm step that " +
			"normally surfaces every resource scheduled for deletion. A single misresolved " +
			"stack name or wrong AWS credential resolves to a one-shot wipe of cloud " +
			"infrastructure. `pulumi up --yes` and `pulumi refresh --yes` carry the same " +
			"footgun for resource creation/replacement. Pipe `pulumi preview` output into " +
			"a review step (manual approval, GitHub Actions environment protection rule) " +
			"before applying, and never combine `--yes` with the `destroy` verb in " +
			"automation.",
		Check: checkZC1736,
	})
}

func checkZC1736(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "pulumi" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	sub := cmd.Arguments[0].String()
	switch sub {
	case "destroy", "up", "refresh":
	default:
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--yes" || v == "-y" {
			return []Violation{{
				KataID: "ZC1736",
				Message: "`pulumi " + sub + " " + v + "` skips the preview-and-confirm — " +
					"a misresolved stack or credential wipes / mutates infrastructure " +
					"with no review. Gate behind `pulumi preview` plus a manual " +
					"approval step.",
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
		ID:       "ZC1737",
		Title:    "Error on `wpa_passphrase SSID PASSWORD` — Wi-Fi passphrase in process list",
		Severity: SeverityError,
		Description: "`wpa_passphrase SSID PASSPHRASE` generates `wpa_supplicant.conf` content " +
			"on stdout. Putting PASSPHRASE on the command line lands it in `ps`, `/proc/<" +
			"pid>/cmdline`, shell history, and the audit log of every local user that can " +
			"list processes. Drop the second positional argument and let `wpa_passphrase " +
			"SSID < /run/secrets/wifi` (or piped via stdin from a secrets store) read the " +
			"passphrase from a file descriptor instead.",
		Check: checkZC1737,
	})
}

func checkZC1737(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wpa_passphrase" {
		return nil
	}

	positionals := 0
	for _, arg := range cmd.Arguments {
		v := arg.String()
		// Skip redirection markers — they aren't true positional args.
		if v == "<" || v == ">" || v == ">>" || v == "<<" {
			break
		}
		if v == "" {
			continue
		}
		positionals++
	}

	if positionals < 2 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1737",
		Message: "`wpa_passphrase SSID PASSWORD` puts the Wi-Fi passphrase in argv — " +
			"visible in `ps`, `/proc`, history. Drop the PASSWORD argument and pipe it " +
			"via stdin (`wpa_passphrase SSID < /run/secrets/wifi`).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1738",
		Title:    "Error on `aws rds delete-db-instance --skip-final-snapshot` — DB destroyed unrecoverable",
		Severity: SeverityError,
		Description: "RDS keeps a final snapshot when an instance or cluster is deleted — the only " +
			"path back from a typo'd identifier or wrong account. `--skip-final-snapshot` " +
			"opts out of that snapshot, so the database is gone the moment the API call " +
			"returns; same applies to `aws rds delete-db-cluster --skip-final-snapshot`. " +
			"Drop the flag (or pass `--final-db-snapshot-identifier <name>` so the snapshot " +
			"name is explicit) and verify the snapshot lands before reusing the identifier.",
		Check: checkZC1738,
	})
}

func checkZC1738(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "aws" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "rds" {
		return nil
	}
	sub := cmd.Arguments[1].String()
	if sub != "delete-db-instance" && sub != "delete-db-cluster" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		if arg.String() == "--skip-final-snapshot" {
			return []Violation{{
				KataID: "ZC1738",
				Message: "`aws rds " + sub + " --skip-final-snapshot` deletes the database " +
					"with no recovery snapshot. Drop the flag or pass `--final-db-" +
					"snapshot-identifier <name>` so the snapshot is explicit and " +
					"verifiable.",
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
		ID:       "ZC1739",
		Title:    "Warn on `git submodule update --remote` — pulls upstream HEAD, breaks reproducibility",
		Severity: SeverityWarning,
		Description: "`git submodule update --remote` fetches each submodule's tracked branch HEAD " +
			"instead of the commit pinned in the parent repo's index. Builds become " +
			"non-reproducible — every CI run pulls a different commit — and any compromised " +
			"upstream commit lands directly in the build. Use `git submodule update --init " +
			"--recursive` (defaults to the pinned commit) and bump submodule pins through " +
			"reviewed PRs.",
		Check: checkZC1739,
	})
}

func checkZC1739(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "submodule" || cmd.Arguments[1].String() != "update" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		if arg.String() == "--remote" {
			return []Violation{{
				KataID: "ZC1739",
				Message: "`git submodule update --remote` ignores the pinned commits in the " +
					"parent repo and pulls each submodule's branch HEAD — non-" +
					"reproducible builds, supply-chain risk. Use `--init --recursive` " +
					"and bump pins via reviewed PRs.",
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
		ID:       "ZC1740",
		Title:    "Warn on `gh release upload --clobber` — silent overwrite of release asset",
		Severity: SeverityWarning,
		Description: "`gh release upload TAG FILE --clobber` replaces an existing asset with the " +
			"same name without prompting. In production this is how a release artifact " +
			"gets silently downgraded — a CI job re-runs with a stale build and the user-" +
			"facing download moves backward without anyone noticing. Drop `--clobber` so " +
			"the second upload errors out, or version the asset name (`mytool-1.2.3-linux." +
			"tar.gz` instead of `mytool-linux.tar.gz`) so each upload has a unique slot.",
		Check: checkZC1740,
	})
}

func checkZC1740(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gh" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "release" || cmd.Arguments[1].String() != "upload" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		if arg.String() == "--clobber" {
			return []Violation{{
				KataID: "ZC1740",
				Message: "`gh release upload --clobber` silently replaces an existing " +
					"asset — a re-run can downgrade the user-facing download. Drop " +
					"`--clobber` or version the asset name so each upload has a " +
					"unique slot.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

// mkpasswd flags that take a value as the next argument.
var zc1741ValueFlags = map[string]bool{
	"-m": true, "-S": true, "-R": true, "-P": true,
	"--method": true, "--salt": true, "--rounds": true, "--password-fd": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1741",
		Title:    "Error on `mkpasswd PASSWORD` — clear-text password in process list",
		Severity: SeverityError,
		Description: "`mkpasswd PASSWORD` (whatwg/Debian `whois`-package version) and `mkpasswd " +
			"-m METHOD PASSWORD` hash the password and print the crypt(3) string on stdout. " +
			"Putting PASSWORD on the command line lands it in `ps`, `/proc/<pid>/cmdline`, " +
			"shell history, and the host audit log. Drop the positional password and read " +
			"from stdin (`mkpasswd -s` reads the password from stdin) — pipe the secret " +
			"from a credentials file or vault: `printf %s \"$PASSWORD\" | mkpasswd -s`.",
		Check: checkZC1741,
	})
}

func checkZC1741(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mkpasswd" {
		return nil
	}

	skipNext := false
	hasPositional := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if skipNext {
			skipNext = false
			continue
		}
		if v == "-s" || v == "--stdin" {
			return nil
		}
		if zc1741ValueFlags[v] {
			skipNext = true
			continue
		}
		if v == "" || v[0] == '-' {
			continue
		}
		hasPositional = true
	}

	if !hasPositional {
		return nil
	}

	return []Violation{{
		KataID: "ZC1741",
		Message: "`mkpasswd PASSWORD` puts the cleartext password in argv — visible in " +
			"`ps`, `/proc`, history. Use `mkpasswd -s` and pipe the secret via stdin " +
			"(`printf %s \"$PASSWORD\" | mkpasswd -s`).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1742",
		Title:    "Error on `mc alias set NAME URL ACCESS_KEY SECRET_KEY` — S3 keys in process list",
		Severity: SeverityError,
		Description: "MinIO's `mc alias set NAME URL ACCESS_KEY SECRET_KEY` (also `mc config " +
			"host add ALIAS URL ACCESS SECRET` on legacy versions) accepts the S3 access " +
			"and secret keys as positional arguments. Both land in argv — visible in " +
			"`ps`, `/proc/<pid>/cmdline`, shell history, and CI logs. Drop the trailing " +
			"keys and let `mc alias set NAME URL` prompt for them, or use the `MC_HOST_<" +
			"alias>=https://ACCESS:SECRET@host` env-var form scoped to a single command " +
			"and unset immediately after.",
		Check: checkZC1742,
	})
}

func checkZC1742(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "mc" && ident.Value != "mcli" {
		return nil
	}
	if len(cmd.Arguments) < 4 {
		return nil
	}

	sub0 := cmd.Arguments[0].String()
	sub1 := cmd.Arguments[1].String()
	rest := cmd.Arguments[2:]

	if sub0 == "alias" && sub1 == "set" {
		// Need NAME URL ACCESS_KEY SECRET_KEY (4 positionals after `alias set`).
		if zc1742PositionalCount(rest) >= 4 {
			return zc1742Hit(cmd, "mc alias set ... ACCESS_KEY SECRET_KEY")
		}
	}
	if sub0 == "config" && sub1 == "host" && len(cmd.Arguments) >= 5 {
		// Legacy: `mc config host add ALIAS URL ACCESS SECRET` (5 args).
		if cmd.Arguments[2].String() == "add" && zc1742PositionalCount(cmd.Arguments[3:]) >= 4 {
			return zc1742Hit(cmd, "mc config host add ... ACCESS SECRET")
		}
	}
	return nil
}

func zc1742PositionalCount(args []ast.Expression) int {
	count := 0
	for _, arg := range args {
		v := arg.String()
		if v == "" || v[0] == '-' {
			continue
		}
		count++
	}
	return count
}

func zc1742Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1742",
		Message: "`" + what + "` puts S3 access and secret keys in argv — visible in " +
			"`ps`, `/proc`, history. Drop the trailing keys (mc prompts) or use " +
			"`MC_HOST_<alias>=URL` env-var form scoped to one command.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1743",
		Title:    "Warn on `npm audit fix --force` — accepts major-version dependency bumps silently",
		Severity: SeverityWarning,
		Description: "`npm audit fix --force` (and `pnpm audit --fix --force`) resolves advisories " +
			"by upgrading dependencies past semver-major boundaries when no backward-" +
			"compatible patch exists. The flag accepts every upgrade without surfacing the " +
			"breaking changes — a build can silently move to a new major of a transitive " +
			"dependency that removes APIs your code calls. Drop `--force` and triage each " +
			"advisory individually; `npm audit fix` handles compatible patches, and the " +
			"remaining advisory targets need a pin or a vendored fork.",
		Check: checkZC1743,
	})
}

func checkZC1743(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var matches bool
	switch ident.Value {
	case "npm":
		if len(cmd.Arguments) >= 2 && cmd.Arguments[0].String() == "audit" && cmd.Arguments[1].String() == "fix" {
			matches = true
		}
	case "pnpm":
		if len(cmd.Arguments) >= 2 && cmd.Arguments[0].String() == "audit" {
			hasFix := false
			for _, a := range cmd.Arguments[1:] {
				if a.String() == "--fix" {
					hasFix = true
					break
				}
			}
			if hasFix {
				matches = true
			}
		}
	default:
		return nil
	}
	if !matches {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "--force" || arg.String() == "-f" {
			return []Violation{{
				KataID: "ZC1743",
				Message: "`" + ident.Value + " audit ... --force` accepts every major-" +
					"version bump an advisory triggers — silent breaking changes. Drop " +
					"`--force` and triage advisories one by one.",
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
		ID:       "ZC1744",
		Title:    "Warn on `kubectl port-forward --address 0.0.0.0` — cluster port exposed to every interface",
		Severity: SeverityWarning,
		Description: "`kubectl port-forward` defaults to binding the local end of the tunnel on " +
			"`127.0.0.1`. `--address 0.0.0.0` (or a specific non-loopback IP) exposes the " +
			"target pod's port to every interface on the developer's workstation or the " +
			"bastion host running the command. Anyone on the LAN / VPN can reach internal " +
			"cluster services that never meant to be externally reachable. Drop the flag " +
			"(loopback default), or pick a specific interface that is already scoped to a " +
			"trusted network.",
		Check: checkZC1744,
	})
}

func checkZC1744(node ast.Node) []Violation {
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
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "port-forward" {
		return nil
	}

	prevAddress := false
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if prevAddress {
			if v == "0.0.0.0" || v == "::" {
				return zc1744Hit(cmd, "--address "+v)
			}
			prevAddress = false
			continue
		}
		switch {
		case v == "--address":
			prevAddress = true
		case strings.HasPrefix(v, "--address="):
			val := strings.TrimPrefix(v, "--address=")
			if val == "0.0.0.0" || val == "::" {
				return zc1744Hit(cmd, v)
			}
		}
	}
	return nil
}

func zc1744Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1744",
		Message: "`kubectl port-forward " + what + "` binds the local end of the tunnel " +
			"on every interface — anyone on the LAN / VPN can reach the pod. Drop " +
			"`--address` (loopback default) or pick a trusted-network interface IP.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1745",
		Title:    "Error on `poetry publish --password PASS` / `twine upload -p PASS` — registry secret in argv",
		Severity: SeverityError,
		Description: "Poetry's `publish --username USER --password PASS` and Twine's `upload " +
			"--username USER --password PASS` (or the short `-u`/`-p` forms) put the PyPI / " +
			"private-index password in argv — visible in `ps`, `/proc/<pid>/cmdline`, shell " +
			"history, and CI logs. Use the `POETRY_PYPI_TOKEN_<NAME>` / `TWINE_USERNAME` + " +
			"`TWINE_PASSWORD` environment variables (sourced from a secrets manager) or a " +
			"`~/.pypirc` file with `0600` perms so the credential never reaches the command " +
			"line.",
		Check: checkZC1745,
	})
}

func checkZC1745(node ast.Node) []Violation {
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
	case "poetry":
		if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "publish" {
			return nil
		}
		tool = "poetry publish"
	case "twine":
		if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "upload" {
			return nil
		}
		tool = "twine upload"
	default:
		return nil
	}

	prevPwd := false
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if prevPwd {
			return zc1745Hit(cmd, tool, "--password "+v)
		}
		switch {
		case v == "--password" || v == "-p":
			prevPwd = true
		case strings.HasPrefix(v, "--password="):
			return zc1745Hit(cmd, tool, v)
		}
	}
	return nil
}

func zc1745Hit(cmd *ast.SimpleCommand, tool, what string) []Violation {
	return []Violation{{
		KataID: "ZC1745",
		Message: "`" + tool + " " + what + "` puts the registry password in argv — " +
			"visible in `ps`, `/proc`, history. Use env vars (`POETRY_PYPI_TOKEN_<" +
			"NAME>`, `TWINE_PASSWORD`) or a `0600` `~/.pypirc` file.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1746",
		Title:    "Error on `sysctl -w kernel.randomize_va_space=0|1` — weakens or disables ASLR",
		Severity: SeverityError,
		Description: "`kernel.randomize_va_space` controls Address Space Layout Randomization. " +
			"Value `2` (default) randomizes stack, heap, VDSO, and mmap regions; value `1` " +
			"omits the heap; value `0` disables ASLR entirely, making every memory layout " +
			"deterministic. Exploits that rely on absolute addresses — stack overflows, ROP " +
			"chains, kernel gadgets — become one-shot instead of brute-forceable. Never " +
			"lower this below `2` outside a sandboxed kernel-debug context.",
		Check: checkZC1746,
	})
}

func checkZC1746(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sysctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "kernel.randomize_va_space=0" || v == "kernel.randomize_va_space=1" {
			return []Violation{{
				KataID: "ZC1746",
				Message: "`sysctl " + v + "` weakens ASLR — absolute-address exploits " +
					"become deterministic (stack overflows, ROP). Keep " +
					"`kernel.randomize_va_space=2` outside a sandboxed debug context.",
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
		ID:       "ZC1747",
		Title:    "Error on `npm/yarn/pnpm --registry http://...` — plaintext registry allows MITM",
		Severity: SeverityError,
		Description: "`npm install --registry http://...`, `pnpm --registry http://...`, and " +
			"`yarn config set registry http://...` configure a plaintext HTTP package " +
			"registry. Any network-position attacker (open Wi-Fi, hostile transit, MITM " +
			"proxy) can replace tarball metadata or content in flight; npm install-time " +
			"`postinstall` scripts then execute the swapped code on the build host. Switch " +
			"the registry URL to `https://` (or terminate TLS at the internal mirror) and " +
			"pair it with a lockfile to pin tarball integrity hashes.",
		Check: checkZC1747,
	})
}

func checkZC1747(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "npm", "pnpm":
		return zc1747ScanFlags(cmd, ident.Value)
	case "yarn":
		if hit := zc1747ScanFlags(cmd, "yarn"); hit != nil {
			return hit
		}
		return zc1747YarnConfigSet(cmd)
	}
	return nil
}

func zc1747ScanFlags(cmd *ast.SimpleCommand, tool string) []Violation {
	prevRegistry := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevRegistry {
			if strings.HasPrefix(v, "http://") {
				return zc1747Hit(cmd, tool, "--registry "+v)
			}
			prevRegistry = false
			continue
		}
		switch {
		case v == "--registry":
			prevRegistry = true
		case strings.HasPrefix(v, "--registry="):
			if strings.HasPrefix(strings.TrimPrefix(v, "--registry="), "http://") {
				return zc1747Hit(cmd, tool, v)
			}
		}
	}
	return nil
}

func zc1747YarnConfigSet(cmd *ast.SimpleCommand) []Violation {
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
	if strings.HasPrefix(url, "http://") {
		return zc1747Hit(cmd, "yarn", "config set registry "+url)
	}
	return nil
}

func zc1747Hit(cmd *ast.SimpleCommand, tool, what string) []Violation {
	return []Violation{{
		KataID: "ZC1747",
		Message: "`" + tool + " " + what + "` uses plaintext HTTP for the package registry — " +
			"any MITM swaps tarballs and runs install-time `postinstall` code. Use " +
			"`https://`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1748",
		Title:    "Error on `helm repo add NAME http://...` — plaintext chart repo allows MITM",
		Severity: SeverityError,
		Description: "`helm repo add NAME http://URL` registers a chart repository reached over " +
			"plaintext HTTP. Any network-position attacker can swap `index.yaml` or a " +
			"chart tarball in flight, and subsequent `helm install` pulls container images " +
			"and Kubernetes manifests straight from the substituted content — fast path to " +
			"cluster-wide code execution. Use `https://`, and pair it with chart provenance " +
			"(`helm install --verify` or OCI signatures) to pin the digest.",
		Check: checkZC1748,
	})
}

func checkZC1748(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "helm" {
		return nil
	}
	if len(cmd.Arguments) < 4 {
		return nil
	}
	if cmd.Arguments[0].String() != "repo" || cmd.Arguments[1].String() != "add" {
		return nil
	}

	url := cmd.Arguments[3].String()
	if !strings.HasPrefix(url, "http://") {
		return nil
	}

	return []Violation{{
		KataID: "ZC1748",
		Message: "`helm repo add " + cmd.Arguments[2].String() + " " + url + "` fetches " +
			"charts over plaintext HTTP — any MITM swaps the chart and its referenced " +
			"images. Use `https://` and `helm install --verify` for provenance.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1749",
		Title:    "Error on `virsh undefine DOMAIN --remove-all-storage` — wipes VM disk images",
		Severity: SeverityError,
		Description: "`virsh undefine DOMAIN --remove-all-storage` (also `--wipe-storage` and the " +
			"newer `--storage <vol,vol>`) removes the VM's configuration AND deletes every " +
			"disk image the domain references. There is no soft-delete and no recycle bin — " +
			"a misresolved DOMAIN or a shared storage pool turns one typo into data loss " +
			"across VMs that happened to share a snapshot chain. Split the operation: back " +
			"up the qcow2 images (`virsh vol-clone` or `qemu-img convert`), then `virsh " +
			"undefine` without the storage flags, then delete volumes deliberately with " +
			"`virsh vol-delete` after a review.",
		Check: checkZC1749,
	})
}

func checkZC1749(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "virsh" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "undefine" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--remove-all-storage" || v == "--wipe-storage" {
			return []Violation{{
				KataID: "ZC1749",
				Message: "`virsh undefine ... " + v + "` deletes every disk image the " +
					"domain references — no soft-delete, no recycle bin. Back up " +
					"first (`qemu-img convert`), `undefine` without the flag, then " +
					"`virsh vol-delete` deliberately.",
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
		ID:       "ZC1750",
		Title:    "Error on `kubectl proxy --address 0.0.0.0` — cluster API proxy on every interface",
		Severity: SeverityError,
		Description: "`kubectl proxy` tunnels Kubernetes API requests authenticated with the " +
			"local kubeconfig's credentials. Defaults bind to `127.0.0.1` and accept only " +
			"`localhost` hosts. `--address 0.0.0.0` (or a specific non-loopback IP) exposes " +
			"that tunnel to every interface on the workstation / bastion, so anyone on the " +
			"LAN or VPN gets the cluster admin the kubeconfig holds. Same risk applies to " +
			"`--accept-hosts '.*'`. Keep the loopback default and scope with SSH port " +
			"forwarding, or restrict `--address` to an interface behind a tight firewall.",
		Check: checkZC1750,
	})
}

func checkZC1750(node ast.Node) []Violation {
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
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "proxy" {
		return nil
	}

	prevAddress := false
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if prevAddress {
			if v == "0.0.0.0" || v == "::" {
				return zc1750Hit(cmd, "--address "+v)
			}
			prevAddress = false
			continue
		}
		switch {
		case v == "--address":
			prevAddress = true
		case strings.HasPrefix(v, "--address="):
			val := strings.TrimPrefix(v, "--address=")
			if val == "0.0.0.0" || val == "::" {
				return zc1750Hit(cmd, v)
			}
		}
	}
	return nil
}

func zc1750Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1750",
		Message: "`kubectl proxy " + what + "` exposes the cluster-admin API tunnel to every " +
			"reachable interface. Keep the loopback default and tunnel over SSH, or " +
			"restrict `--address` to a firewalled interface.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1751",
		Title:    "Error on `rpm/dnf/yum remove --nodeps` — bypasses dependency check, breaks dependents",
		Severity: SeverityError,
		Description: "`rpm -e --nodeps PKG` (also `dnf remove --nodeps`, `yum remove --nodeps`, " +
			"`zypper remove --force`) removes the package while skipping the dependency " +
			"solver. Anything transitively depending on the target immediately breaks — " +
			"`libc`, `openssl`, `systemd` units, even `dnf` itself can get pulled out, " +
			"leaving the host unbootable or unpackageable. Resolve the dependency conflict " +
			"explicitly (`dnf swap`, `rpm -e --rebuilddb` never, pin the conflicting package) " +
			"instead of bypassing the check.",
		Check: checkZC1751,
	})
}

func checkZC1751(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var tool string
	var verbOK bool
	switch ident.Value {
	case "rpm":
		tool = "rpm"
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-e" || v == "--erase" {
				verbOK = true
				break
			}
		}
	case "dnf", "yum", "microdnf":
		tool = ident.Value
		if len(cmd.Arguments) == 0 {
			return nil
		}
		sub := cmd.Arguments[0].String()
		if sub == "remove" || sub == "erase" {
			verbOK = true
		}
	case "zypper":
		tool = "zypper"
		if len(cmd.Arguments) == 0 {
			return nil
		}
		sub := cmd.Arguments[0].String()
		if sub == "remove" || sub == "rm" {
			verbOK = true
		}
	default:
		return nil
	}
	if !verbOK {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--nodeps" || v == "--no-deps" {
			return []Violation{{
				KataID: "ZC1751",
				Message: "`" + tool + " ... " + v + "` removes the package without the " +
					"dependency solver — dependents break (libc, openssl, systemd " +
					"units). Resolve the conflict explicitly instead of bypassing.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

var zc1752ForceFlags = map[string]bool{
	"-f": true, "-ff": true, "-fff": true,
	"--force": true,
	"-y":      true, "--yes": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1752",
		Title:    "Error on `pvcreate/vgcreate/lvcreate -ff|--yes` — force-init LVM over existing data",
		Severity: SeverityError,
		Description: "LVM prompts before overwriting existing filesystem, RAID, or LVM signatures " +
			"on a device — that prompt is the only thing saving you from a typo'd target " +
			"destroying someone else's data. `pvcreate -ff`, `pvcreate --yes`, and the same " +
			"flags on `vgcreate` / `lvcreate` skip the prompt. Drop the flag, inspect with " +
			"`wipefs -n` and `lsblk -f` first, then confirm the target before re-running " +
			"the create command.",
		Check: checkZC1752,
	})
}

func checkZC1752(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "pvcreate", "vgcreate", "lvcreate":
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1752ForceFlags[v] {
			return []Violation{{
				KataID: "ZC1752",
				Message: "`" + ident.Value + " " + v + "` skips the LVM confirmation — a " +
					"wrong device gets its filesystem / RAID / LVM signatures wiped. " +
					"Inspect with `wipefs -n` + `lsblk -f` first, drop the flag, re-" +
					"run after checking the target.",
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
		ID:       "ZC1753",
		Title:    "Error on `rclone purge REMOTE:PATH` — bulk delete of every object under the remote path",
		Severity: SeverityError,
		Description: "`rclone purge REMOTE:PATH` removes every object and empty directory under " +
			"PATH on the remote — no dry-run gate, no confirmation, no soft-delete unless " +
			"the backend happens to version. A typo'd path or a stale variable turns one " +
			"line into a bucket-wide wipe (S3, GCS, Azure, Swift all honour the same API " +
			"call). Preview with `rclone lsf REMOTE:PATH` or `rclone delete --dry-run`, " +
			"then use `rclone delete` scoped narrower; enable object versioning on the " +
			"backend so a bad run can roll back.",
		Check: checkZC1753,
	})
}

func checkZC1753(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "rclone" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	if cmd.Arguments[0].String() != "purge" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1753",
		Message: "`rclone purge` removes every object under the remote path with no dry-run " +
			"or soft-delete. Preview with `rclone lsf` / `rclone delete --dry-run` and " +
			"prefer narrower `rclone delete`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1754",
		Title:    "Error on `gh auth status -t` / `--show-token` — prints OAuth token to stdout",
		Severity: SeverityError,
		Description: "`gh auth status -t` (alias `--show-token`) prints the stored GitHub OAuth " +
			"token alongside the status summary. In CI logs, shared terminals, piped to " +
			"`less`/`tee`, or captured via `script`, the token ends up on disk or in " +
			"scrollback where anyone with log access becomes repo-admin. Never combine " +
			"`-t` with `auth status` in automation; if a machine-readable token is needed, " +
			"`gh auth token` prints only the token and makes the secret-handling path " +
			"explicit.",
		Check: checkZC1754,
	})
}

func checkZC1754(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gh" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "auth" || cmd.Arguments[1].String() != "status" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if v == "-t" || v == "--show-token" {
			return []Violation{{
				KataID: "ZC1754",
				Message: "`gh auth status " + v + "` prints the OAuth token in the status " +
					"output — CI logs and scrollback become a repo-admin leak. Use " +
					"`gh auth token` in automation so the secret path is explicit.",
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
		ID:       "ZC1755",
		Title:    "Error on `gcloud sql users {create,set-password} --password PASS` — DB password in argv",
		Severity: SeverityError,
		Description: "`gcloud sql users create USER --instance INST --password PASS` (and the " +
			"`set-password` variant) place the Cloud SQL user password on the command " +
			"line — visible in `ps`, `/proc/<pid>/cmdline`, shell history, and CI logs, " +
			"and stored in Cloud Audit Logs' request payload. Use `--prompt-for-password` " +
			"(interactive) or generate the password server-side in Secret Manager and post " +
			"to the SQL Admin API via `gcloud auth print-access-token` piped to `curl` with " +
			"the body sourced from a file.",
		Check: checkZC1755,
	})
}

func checkZC1755(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gcloud" {
		return nil
	}
	if len(cmd.Arguments) < 4 {
		return nil
	}
	if cmd.Arguments[0].String() != "sql" || cmd.Arguments[1].String() != "users" {
		return nil
	}
	sub := cmd.Arguments[2].String()
	if sub != "create" && sub != "set-password" {
		return nil
	}

	prevPwd := false
	for _, arg := range cmd.Arguments[3:] {
		v := arg.String()
		if prevPwd {
			return zc1755Hit(cmd, sub, "--password "+v)
		}
		switch {
		case v == "--password":
			prevPwd = true
		case strings.HasPrefix(v, "--password="):
			return zc1755Hit(cmd, sub, v)
		}
	}
	return nil
}

func zc1755Hit(cmd *ast.SimpleCommand, sub, what string) []Violation {
	return []Violation{{
		KataID: "ZC1755",
		Message: "`gcloud sql users " + sub + " " + what + "` puts the Cloud SQL password " +
			"in argv — visible in `ps`, `/proc`, history, and Cloud Audit Logs. Use " +
			"`--prompt-for-password` or call the SQL Admin API with a body file.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

var zc1756DockerSocketPaths = map[string]bool{
	"/var/run/docker.sock":              true,
	"/run/docker.sock":                  true,
	"/run/containerd/containerd.sock":   true,
	"/var/run/crio/crio.sock":           true,
	"/run/crio/crio.sock":               true,
	"/var/run/podman/podman.sock":       true,
	"/run/podman/podman.sock":           true,
	"/run/user/1000/podman/podman.sock": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1756",
		Title:    "Error on `chmod NNN /run/docker.sock` — world access is root-equivalent privesc",
		Severity: SeverityError,
		Description: "Container-runtime sockets (`/var/run/docker.sock`, `/run/containerd/" +
			"containerd.sock`, `/run/crio/crio.sock`, `/run/podman/podman.sock`) accept " +
			"commands that run on the host with root privilege — starting privileged " +
			"containers, mounting the host filesystem, reading every file on disk. " +
			"Making the socket world-readable or world-writable (`chmod 644/660/666/777`) " +
			"hands every local user that root-escalation primitive. Keep the socket " +
			"`0660 root:docker` (or the equivalent runtime group) and add only trusted " +
			"accounts to that group.",
		Check: checkZC1756,
	})
}

func checkZC1756(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chmod" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	var mode, target string
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1756DockerSocketPaths[v] {
			target = v
			continue
		}
		if mode == "" && zc1756WorldAccess(v) {
			mode = v
		}
	}
	if mode == "" || target == "" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1756",
		Message: "`chmod " + mode + " " + target + "` grants every local user access to a " +
			"root-equivalent container-runtime socket. Keep `0660` owned by the " +
			"runtime group (`root:docker` etc.) and restrict membership.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

// zc1756WorldAccess: true when MODE grants world-read or world-write on the
// target. Handles both octal literal and the parser-normalised decimal form.
func zc1756WorldAccess(mode string) bool {
	if mode == "" {
		return false
	}
	if zc1671WorldWritable(mode) {
		return true
	}
	// zc1671WorldWritable catches world-write; also catch world-read (bit 4).
	// Parse octal, fallback to decimal.
	for _, c := range mode {
		if c < '0' || c > '7' {
			// Decimal fallback (parser-normalised leading-zero octal).
			return false
		}
	}
	// Simple octal parse: last char is world triad.
	last := mode[len(mode)-1]
	return last == '4' || last == '5' || last == '6' || last == '7'
}

var zc1757DangerousScopes = []string{
	"delete_repo",
	"admin:org",
	"admin:enterprise",
	"admin:public_key",
	"admin:org_hook",
	"site_admin",
	"admin:repo_hook",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1757",
		Title:    "Warn on `gh auth refresh --scopes delete_repo|admin:*` — token escalated to destructive perms",
		Severity: SeverityWarning,
		Description: "`gh auth refresh --scopes <list>` (also `gh auth login --scopes`) rotates " +
			"the stored OAuth token with additional scopes. `delete_repo`, `admin:org`, " +
			"`admin:enterprise`, `admin:public_key`, and `admin:*_hook` give the token " +
			"permanent destructive perms that outlast the script that asked for them — a " +
			"compromised token now carries repo-deletion, org-membership, and SSH-key " +
			"manipulation rights. Request the minimum scope the task needs (`repo`, " +
			"`workflow`) and rotate the token off when the elevated operation completes.",
		Check: checkZC1757,
	})
}

func checkZC1757(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gh" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "auth" {
		return nil
	}
	sub := cmd.Arguments[1].String()
	if sub != "refresh" && sub != "login" {
		return nil
	}

	prevScopes := false
	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if prevScopes {
			if hit := zc1757MatchScopes(v); hit != "" {
				return zc1757Hit(cmd, sub, hit)
			}
			prevScopes = false
			continue
		}
		switch {
		case v == "--scopes" || v == "-s":
			prevScopes = true
		case strings.HasPrefix(v, "--scopes="):
			if hit := zc1757MatchScopes(strings.TrimPrefix(v, "--scopes=")); hit != "" {
				return zc1757Hit(cmd, sub, hit)
			}
		}
	}
	return nil
}

func zc1757MatchScopes(list string) string {
	for _, scope := range strings.Split(list, ",") {
		scope = strings.TrimSpace(scope)
		for _, danger := range zc1757DangerousScopes {
			if scope == danger {
				return scope
			}
		}
	}
	return ""
}

func zc1757Hit(cmd *ast.SimpleCommand, sub, scope string) []Violation {
	return []Violation{{
		KataID: "ZC1757",
		Message: "`gh auth " + sub + " --scopes " + scope + "` escalates the token to " +
			"destructive privileges that outlast the script. Request the minimum " +
			"scope (`repo`, `workflow`) and rotate the token when done.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1758",
		Title:    "Warn on `gh codespace delete --force` — destroys codespace with uncommitted work",
		Severity: SeverityWarning,
		Description: "`gh codespace delete --force` (alias `-f`) skips the confirmation prompt " +
			"and deletes the target codespace along with any uncommitted, unpushed, or " +
			"unstaged work inside it. Combined with `--all`, one line wipes every codespace " +
			"on the account. Drop the flag, let the prompt enumerate what is about to go, " +
			"and only confirm after verifying no local state would be lost — `git status` " +
			"/ `git stash list` inside the codespace first.",
		Check: checkZC1758,
	})
}

func checkZC1758(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gh" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "codespace" || cmd.Arguments[1].String() != "delete" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if v == "--force" || v == "-f" {
			return []Violation{{
				KataID: "ZC1758",
				Message: "`gh codespace delete " + v + "` skips the prompt and drops " +
					"uncommitted work along with the codespace. Let the prompt list " +
					"what's about to go and verify `git status` inside first.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

var zc1759SecretKeys = []string{
	"password", "passwd",
	"token", "secret",
	"apikey", "api_key", "api-key",
	"accesskey", "access_key", "access-key",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1759",
		Title:    "Error on `vault login TOKEN` / `login -method=… password=…` — credential in process list",
		Severity: SeverityError,
		Description: "Vault accepts credentials on its `login` / `auth` subcommands in two " +
			"argv-leaking shapes: a positional token (`vault login <TOKEN>`) and KEY=VALUE " +
			"pairs for non-token methods (`vault login -method=userpass username=U " +
			"password=P`). Both land the secret in `ps`, `/proc/<pid>/cmdline`, shell " +
			"history, and Vault's audit log request payload. Read the token from stdin " +
			"(`vault login -` with `printf %s \"$TOKEN\" |`) or source `VAULT_TOKEN` from " +
			"a secrets file and run `vault login -method=token`.",
		Check: checkZC1759,
	})
}

func checkZC1759(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "vault" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "login" && sub != "auth" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		// Stdin sentinel, file sentinel, flag-forms are safe.
		if v == "-" || strings.HasPrefix(v, "-") || strings.HasPrefix(v, "@") {
			continue
		}
		// KEY=VALUE pair: flag secret-named keys.
		if eq := strings.IndexByte(v, '='); eq > 0 {
			key := strings.ToLower(v[:eq])
			for _, secret := range zc1759SecretKeys {
				if strings.Contains(key, secret) {
					return zc1759Hit(cmd, sub, v)
				}
			}
			continue
		}
		// Bare positional token.
		return zc1759Hit(cmd, sub, v)
	}
	return nil
}

func zc1759Hit(cmd *ast.SimpleCommand, sub, what string) []Violation {
	return []Violation{{
		KataID: "ZC1759",
		Message: "`vault " + sub + " " + what + "` puts the Vault credential in argv — " +
			"visible in `ps`, `/proc`, history, Vault audit log. Use `vault login -` " +
			"with stdin or source `VAULT_TOKEN` from a secrets file.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1760",
		Title:    "Warn on `openssl rand -hex|-base64 N` with N < 16 — generated value too short",
		Severity: SeverityWarning,
		Description: "`openssl rand -hex N` (and `-base64 N`) outputs N random bytes encoded into " +
			"the requested form. N below 16 (128 bits) produces a value short enough that " +
			"an attacker with modest GPU resources can brute-force it offline — too weak " +
			"for passwords, API tokens, reset URLs, or any other secret that sits at rest. " +
			"Use `-hex 32` (256-bit) for secrets and long-lived tokens; `-hex 16` is " +
			"acceptable only for short-validity nonces paired with rate-limited consumers.",
		Check: checkZC1760,
	})
}

func checkZC1760(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "openssl" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "rand" {
		return nil
	}

	prevEnc := ""
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if prevEnc != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 && n < 16 {
				return zc1760Hit(cmd, prevEnc+" "+v)
			}
			prevEnc = ""
			continue
		}
		if v == "-hex" || v == "-base64" {
			prevEnc = v
		}
	}
	return nil
}

func zc1760Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1760",
		Message: "`openssl rand " + what + "` produces a sub-128-bit value — brute-forceable " +
			"offline. Use `-hex 32` for secrets / long-lived tokens.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1761",
		Title:    "Warn on `gh gist create --public` — file becomes world-visible and indexed on GitHub",
		Severity: SeverityWarning,
		Description: "`gh gist create --public FILE` (alias `-p`) creates the gist with `public: " +
			"true`. Public gists are listed on `gist.github.com/discover`, crawled by " +
			"search engines, and archived by secondary scrapers — a leaked secret, private " +
			"company snippet, or unreleased note is effectively permanent the moment it " +
			"lands. The default (`public: false`) keeps the gist unlisted and reachable " +
			"only via its URL. Drop `--public` unless public exposure is the explicit " +
			"goal.",
		Check: checkZC1761,
	})
}

func checkZC1761(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gh" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "gist" || cmd.Arguments[1].String() != "create" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if v == "--public" || v == "-p" {
			return []Violation{{
				KataID: "ZC1761",
				Message: "`gh gist create " + v + "` publishes the file to the public " +
					"discover feed — search engines crawl it within minutes. Drop the " +
					"flag unless public exposure is the explicit goal.",
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
		ID:       "ZC1762",
		Title:    "Error on `kubeadm join --discovery-token-unsafe-skip-ca-verification` — cluster CA not checked",
		Severity: SeverityError,
		Description: "`kubeadm join` verifies the control-plane API server's CA before accepting " +
			"the kubelet bootstrap token. `--discovery-token-unsafe-skip-ca-verification` " +
			"skips that check, so a network-position attacker can impersonate the API " +
			"server, harvest the bootstrap token, and seed malicious workloads onto the " +
			"joining node. Always pin the CA with `--discovery-token-ca-cert-hash sha256:" +
			"<digest>` (emitted by `kubeadm token create --print-join-command`) or supply " +
			"a kubeconfig discovery file that has the CA baked in.",
		Check: checkZC1762,
	})
}

func checkZC1762(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubeadm" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "join" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "--discovery-token-unsafe-skip-ca-verification" {
			return []Violation{{
				KataID: "ZC1762",
				Message: "`kubeadm join --discovery-token-unsafe-skip-ca-verification` " +
					"skips CA verification of the control-plane — MITM steals the " +
					"bootstrap token. Pin the CA with `--discovery-token-ca-cert-hash " +
					"sha256:<digest>`.",
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
		ID:       "ZC1763",
		Title:    "Error on `docker compose down -v` / `--volumes` — wipes named volumes (data loss)",
		Severity: SeverityError,
		Description: "`docker compose down -v` (alias `--volumes`, equivalent in `docker-compose " +
			"down -v`) tears the stack down AND deletes every named volume declared in the " +
			"compose file. Database contents, cache state, uploaded assets, and any other " +
			"volume-backed data goes with them — there is no soft-delete. Drop the flag in " +
			"CI and production scripts; keep it only for throwaway local testbeds where " +
			"losing volume state is intentional.",
		Check: checkZC1763,
	})
}

func checkZC1763(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var argsAfterDown []ast.Expression
	switch ident.Value {
	case "docker":
		if len(cmd.Arguments) < 2 {
			return nil
		}
		if cmd.Arguments[0].String() != "compose" || cmd.Arguments[1].String() != "down" {
			return nil
		}
		argsAfterDown = cmd.Arguments[2:]
	case "docker-compose":
		if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "down" {
			return nil
		}
		argsAfterDown = cmd.Arguments[1:]
	default:
		return nil
	}

	for _, arg := range argsAfterDown {
		v := arg.String()
		if v == "-v" || v == "--volumes" {
			return []Violation{{
				KataID: "ZC1763",
				Message: "`docker compose down " + v + "` wipes every named volume declared " +
					"in the stack — database, cache, uploaded assets go with it. Drop " +
					"the flag in CI / prod scripts.",
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
		ID:       "ZC1764",
		Title:    "Warn on `git commit --no-verify` / `-n` — skips pre-commit and commit-msg hooks",
		Severity: SeverityWarning,
		Description: "`git commit --no-verify` (alias `-n`) bypasses both the pre-commit and " +
			"commit-msg hooks, which are often the last guardrail against leaked secrets, " +
			"formatting drift, or failing tests. The flag is usually a symptom of a hook " +
			"that needs fixing rather than silencing — the exception quickly becomes the " +
			"rule. Fix the blocking hook, carve out a narrow per-file exemption in the " +
			"hook itself, or file a tracked issue, instead of adding `--no-verify` to " +
			"every commit in a script.",
		Check: checkZC1764,
	})
}

func checkZC1764(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "commit" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--no-verify" || v == "-n" {
			return []Violation{{
				KataID: "ZC1764",
				Message: "`git commit " + v + "` skips pre-commit and commit-msg hooks " +
					"— the last guardrail against secret leaks and broken tests. Fix " +
					"the hook or carve a narrow exemption instead.",
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
		ID:       "ZC1765",
		Title:    "Error on `snap remove --purge SNAP` — skips the automatic data snapshot",
		Severity: SeverityError,
		Description: "`snap remove SNAP` takes a snapshot of every writable area (`$SNAP_DATA`, " +
			"`$SNAP_USER_DATA`, `$SNAP_COMMON`) before uninstalling, so the data can later " +
			"be restored with `snap restore`. `--purge` skips that snapshot: the snap is " +
			"gone along with every file it owned, and snapd has no record to roll back. " +
			"Drop `--purge` unless the snap's data is genuinely disposable; otherwise " +
			"`snap save SNAP` first, capture the set ID, and only then remove.",
		Check: checkZC1765,
	})
}

func checkZC1765(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "snap" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "remove" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "--purge" {
			return []Violation{{
				KataID: "ZC1765",
				Message: "`snap remove --purge` skips the pre-remove data snapshot — the " +
					"snap's files are gone with no rollback. Drop `--purge` or capture " +
					"a `snap save` set ID first.",
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
		ID:       "ZC1766",
		Title:    "Error on `memcached -l 0.0.0.0` — memcached exposed on every interface",
		Severity: SeverityError,
		Description: "`memcached -l 0.0.0.0` (or `::`, `--listen=0.0.0.0`) binds memcached's TCP " +
			"listener to every interface on the host. Memcached has no authentication and, " +
			"before `-U 0` became default, its UDP handler was the largest DDoS-" +
			"amplification vector on the internet. Bind to `127.0.0.1` or a private-" +
			"network IP only, and put memcached behind a firewall / security group scoped " +
			"to the application that consumes it.",
		Check: checkZC1766,
	})
}

func checkZC1766(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "memcached" {
		return nil
	}

	prevListen := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevListen {
			if zc1766IsUnrestrictedBind(v) {
				return zc1766Hit(cmd, "-l "+v)
			}
			prevListen = false
			continue
		}
		switch {
		case v == "-l":
			prevListen = true
		case strings.HasPrefix(v, "-l") && len(v) > 2:
			if zc1766IsUnrestrictedBind(v[2:]) {
				return zc1766Hit(cmd, v)
			}
		case strings.HasPrefix(v, "--listen="):
			if zc1766IsUnrestrictedBind(strings.TrimPrefix(v, "--listen=")) {
				return zc1766Hit(cmd, v)
			}
		}
	}
	return nil
}

func zc1766IsUnrestrictedBind(s string) bool {
	return s == "0.0.0.0" || s == "::" || s == "[::]"
}

func zc1766Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1766",
		Message: "`memcached " + what + "` exposes the unauthenticated cache to every " +
			"interface on the host. Bind to `127.0.0.1` or a private-network IP and " +
			"firewall the port.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1767",
		Title:    "Error on `mongod --bind_ip 0.0.0.0` — MongoDB exposed on every interface",
		Severity: SeverityError,
		Description: "`mongod --bind_ip 0.0.0.0` (or `::`) binds MongoDB's listener to every " +
			"interface on the host. Combined with no-auth defaults (pre-3.4) or a wildcard " +
			"database user, this was the source of the 2017 ransomware wave that wiped " +
			"tens of thousands of public MongoDB instances. Bind to `127.0.0.1` or a " +
			"private-network IP, enable authentication with `--auth`, and firewall port " +
			"`27017`.",
		Check: checkZC1767,
	})
}

var zc1767BindFlags = map[string]bool{"--bind_ip": true}

func checkZC1767(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "mongod" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if !zc1767BindFlags[arg.String()] {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			return nil
		}
		ip := cmd.Arguments[i+1].String()
		if ip != "0.0.0.0" && ip != "::" && ip != "[::]" {
			return nil
		}
		line, col := FlagArgPosition(cmd, zc1767BindFlags)
		return []Violation{{
			KataID: "ZC1767",
			Message: "`mongod --bind_ip " + ip + "` exposes MongoDB on every interface — " +
				"2017 ransomware-wave target. Bind to `127.0.0.1` or a private-network IP, " +
				"enable `--auth`, firewall port 27017.",
			Line:   line,
			Column: col,
			Level:  SeverityError,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1768",
		Title:    "Error on `sqlcmd -P PASSWORD` / `bcp -P PASSWORD` — SQL Server password in argv",
		Severity: SeverityError,
		Description: "Microsoft's SQL Server CLI tools (`sqlcmd`, `bcp`, `osql`) accept the " +
			"password via `-P PASSWORD` as a positional argument value. The password lands " +
			"in argv — visible in `ps`, `/proc/<pid>/cmdline`, shell history, CI logs, and " +
			"SQL Server's audit trace for the session. Use `-P` with no value (prompts), " +
			"or read the password from the environment variable `SQLCMDPASSWORD` (sourced " +
			"from a secrets file). On modern sqlcmd, `-G` + Azure AD integration avoids the " +
			"password altogether.",
		Check: checkZC1768,
	})
}

func checkZC1768(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "sqlcmd", "bcp", "osql":
	default:
		return nil
	}

	prevP := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevP {
			prevP = false
			if v == "" || v == "-" {
				continue
			}
			// Any other flag-looking value means `-P` ran as bare (prompt).
			if len(v) > 1 && v[0] == '-' {
				continue
			}
			return []Violation{{
				KataID: "ZC1768",
				Message: "`" + ident.Value + " -P " + v + "` puts the SQL Server password " +
					"in argv — visible in `ps`, `/proc`, history, SQL Server audit. " +
					"Use `-P` with no arg (prompt) or `SQLCMDPASSWORD` env var.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
		if v == "-P" {
			prevP = true
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1769",
		Title:    "Warn on `vagrant destroy --force` — VM destroyed without confirmation",
		Severity: SeverityWarning,
		Description: "`vagrant destroy --force` (alias `-f`) tears every VM in the Vagrantfile " +
			"down — and their ephemeral filesystem state — without prompting. Any data " +
			"provisioned into the VM that was never exported back to the host (database " +
			"seeds, build caches, local-only test fixtures) goes with it. In unattended " +
			"scripts, drop the flag so the prompt still gates the destroy; for CI cycles, " +
			"`vagrant halt` + `vagrant up` reuses the same box without losing state.",
		Check: checkZC1769,
	})
}

func checkZC1769(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "vagrant" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "destroy" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--force" || v == "-f" {
			return []Violation{{
				KataID: "ZC1769",
				Message: "`vagrant destroy " + v + "` skips the prompt and drops the VM " +
					"(and any un-exported data inside). Drop the flag, or use `vagrant " +
					"halt` + `vagrant up` to cycle without destroy.",
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
		ID:       "ZC1770",
		Title:    "Warn on `gpg --always-trust` / `--trust-model always` — bypasses Web-of-Trust",
		Severity: SeverityWarning,
		Description: "`gpg --always-trust` (equivalent to `--trust-model always`) accepts every key " +
			"in the keyring as fully trusted, regardless of signatures from the owner or any " +
			"introducer. A signature made by an attacker-controlled key pair that was imported " +
			"with no further vetting will verify cleanly. In automation this turns signature " +
			"verification into a presence check — any key bundled with the payload satisfies " +
			"`gpg --verify`. Remove the flag and build a proper trust path: either mark the " +
			"expected signer key trusted once (`gpg --edit-key KEYID trust`), or pin the " +
			"expected fingerprint and match it against the signer after `gpg --verify --status-fd 1`.",
		Check: checkZC1770,
	})
}

func checkZC1770(node ast.Node) []Violation {
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

	for i, arg := range cmd.Arguments {
		v := arg.String()
		switch {
		case v == "--always-trust":
			return zc1770Hit(cmd, "--always-trust", map[string]bool{"--always-trust": true})
		case v == "--trust-model" && i+1 < len(cmd.Arguments) && cmd.Arguments[i+1].String() == "always":
			return zc1770Hit(cmd, "--trust-model always", map[string]bool{"--trust-model": true})
		case v == "--trust-model=always":
			return zc1770Hit(cmd, "--trust-model=always", map[string]bool{"--trust-model=always": true})
		}
	}
	return nil
}

func zc1770Hit(cmd *ast.SimpleCommand, flag string, needle map[string]bool) []Violation {
	line, col := FlagArgPosition(cmd, needle)
	return []Violation{{
		KataID: "ZC1770",
		Message: "`gpg " + flag + "` marks every imported key as fully trusted — a " +
			"signature from an attacker-supplied key verifies cleanly. Drop the flag " +
			"and pin the expected fingerprint, or assign trust via `gpg --edit-key KEYID trust`.",
		Line:   line,
		Column: col,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1771",
		Title:    "Warn on `alias -g` / `alias -s` — global and suffix aliases surprise script readers",
		Severity: SeverityWarning,
		Description: "`alias -g NAME=value` defines a global alias that expands anywhere on the " +
			"command line, not just in command position. `alias -s ext=cmd` (suffix alias) runs " +
			"`cmd file.ext` whenever a bare `file.ext` appears as a command. Both forms are " +
			"Zsh-idiomatic interactive conveniences; in scripts they produce surprising " +
			"substitutions that a reader cannot infer from local context — a bare word like " +
			"`G` or `foo.log` stops meaning what it looks like. Use a function or a regular " +
			"alias instead, and keep `alias -g` / `alias -s` in your `~/.zshrc` where the " +
			"definition is discoverable.",
		Check: checkZC1771,
	})
}

func checkZC1771(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "alias" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	first := cmd.Arguments[0].String()
	switch {
	case first == "-g":
		return zc1771Hit(cmd, "-g", "global")
	case first == "-s":
		return zc1771Hit(cmd, "-s", "suffix")
	case strings.HasPrefix(first, "-") && !strings.HasPrefix(first, "--"):
		if strings.ContainsRune(first, 'g') {
			return zc1771Hit(cmd, first, "global")
		}
		if strings.ContainsRune(first, 's') {
			return zc1771Hit(cmd, first, "suffix")
		}
	}
	return nil
}

func zc1771Hit(cmd *ast.SimpleCommand, flag, kind string) []Violation {
	return []Violation{{
		KataID: "ZC1771",
		Message: "`alias " + flag + "` defines a " + kind + " alias that expands outside " +
			"command position — a surprise for anyone reading the script later. Prefer a " +
			"function, or keep " + kind + " aliases in `~/.zshrc` where they are discoverable.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

var zc1772EraseFlags = map[string]bool{
	"--security-erase":          true,
	"--security-erase-enhanced": true,
	"--security-set-pass":       true,
	"--security-unlock":         true,
	"--security-disable":        true,
	"--security-freeze":         true,
	"--trim-sector-ranges":      true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1772",
		Title:    "Error on `hdparm --security-erase` / `--trim-sector-ranges` — ATA-level data destruction",
		Severity: SeverityError,
		Description: "`hdparm --security-erase PASS $DISK` issues the ATA `SECURITY ERASE UNIT` " +
			"command: the drive firmware wipes every block, ignoring filesystem or partition " +
			"boundaries, and the operation cannot be interrupted or rolled back. " +
			"`--security-erase-enhanced` is the same but also clears reallocated sectors, and " +
			"`--trim-sector-ranges` discards the listed LBAs on any TRIM-capable device. " +
			"`--security-set-pass`, `--security-disable`, `--security-unlock`, and " +
			"`--security-freeze` alter the drive-level password state and, if misused in a " +
			"script, lock the device out of future access. Keep these calls behind a guarded " +
			"runbook with the exact disk pinned by `/dev/disk/by-id/…` and the password stored " +
			"outside argv.",
		Check: checkZC1772,
	})
}

func checkZC1772(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "hdparm" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1772EraseFlags[v] {
			return zc1772Hit(cmd, v)
		}
	}
	return nil
}

func zc1772Hit(cmd *ast.SimpleCommand, flag string) []Violation {
	line, col := FlagArgPosition(cmd, zc1772EraseFlags)
	return []Violation{{
		KataID: "ZC1772",
		Message: "`hdparm " + flag + "` issues an ATA-level operation that ignores " +
			"filesystems and cannot be rolled back. Pin the disk by " +
			"`/dev/disk/by-id/…`, keep it behind a runbook, keep the password out of argv.",
		Line:   line,
		Column: col,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1773",
		Title:    "Warn on `xargs` without `-r` / `--no-run-if-empty` — runs once on empty input",
		Severity: SeverityWarning,
		Description: "GNU `xargs` (the common default on Linux) invokes the child command once " +
			"with no arguments when its stdin is empty. Paired with a destructive child " +
			"(`xargs rm`, `xargs kill`, `xargs docker stop`) a pipeline that produces zero " +
			"hits silently runs the command with no operand — usually an error at best and a " +
			"footgun at worst. The flag `-r` (GNU) / `--no-run-if-empty` tells xargs to skip " +
			"the call when no items arrive. Add `-r` to every `xargs` pipeline whose producer " +
			"can return no results, or switch to `find ... -exec cmd {} +` which never runs " +
			"the child on empty input. BSD xargs defaults to this behavior, but the portable " +
			"and explicit choice is to pass `-r` and document the intent.",
		Check: checkZC1773,
		Fix:   fixZC1773,
	})
}

// fixZC1773 inserts ` -r` after the `xargs` command name. Detector
// already guards against any existing `-r` / `--no-run-if-empty` /
// combined-short-flag form so the insertion is idempotent.
func fixZC1773(node ast.Node, v Violation, source []byte) []FixEdit {
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
	if IdentLenAt(source, nameOff) != len("xargs") {
		return nil
	}
	insertAt := nameOff + len("xargs")
	insLine, insCol := offsetLineColZC1773(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -r",
	}}
}

func offsetLineColZC1773(source []byte, offset int) (int, int) {
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

func checkZC1773(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "xargs" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-r" || v == "--no-run-if-empty" {
			return nil
		}
		// Combined short-flag form like `-rt` or `-0r`.
		if len(v) > 1 && v[0] == '-' && v[1] != '-' {
			for _, c := range v[1:] {
				if c == 'r' {
					return nil
				}
			}
		}
	}
	return []Violation{{
		KataID: "ZC1773",
		Message: "`xargs` without `-r` / `--no-run-if-empty` runs the child once with no " +
			"arguments when stdin is empty — a destructive surprise for `xargs rm`, " +
			"`xargs kill`, etc. Add `-r` or switch to `find ... -exec cmd {} +`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1774",
		Title:    "Warn on `setopt GLOB_SUBST` — `$var` starts glob-expanding, user data becomes a pattern",
		Severity: SeverityWarning,
		Description: "With `GLOB_SUBST` enabled, the result of any parameter expansion is " +
			"rescanned for filename-generation metacharacters (`*`, `?`, `[`, `^`, `~`, " +
			"brace ranges, qualifiers). Zsh's default — `NO_GLOB_SUBST` — keeps `$var` literal " +
			"and matches the behavior most script authors expect after moving from Bash or " +
			"POSIX sh. Turning `GLOB_SUBST` on globally means any unquoted `$var` that " +
			"contains a metacharacter (environment, argv, file contents, user prompt) is " +
			"expanded against the filesystem — an injection vector, and a subtle source of " +
			"`no matches found` failures on empty variables. Keep `setopt GLOB_SUBST` inside a " +
			"narrow subshell or function body, or use explicit `~` / `(e)` / `(P)` flags where " +
			"you actually want the rescan.",
		Check: checkZC1774,
	})
}

func checkZC1774(node ast.Node) []Violation {
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
			v := strings.ToUpper(strings.ReplaceAll(arg.String(), "_", ""))
			if v == "GLOBSUBST" {
				return zc1774Hit(cmd, "setopt "+arg.String())
			}
		}
	case "set":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-G" {
				return zc1774Hit(cmd, "set -G")
			}
			if v == "-o" || v == "--option" {
				continue
			}
			if strings.EqualFold(strings.ReplaceAll(v, "_", ""), "globsubst") {
				return zc1774Hit(cmd, "set -o "+v)
			}
		}
	}
	return nil
}

func zc1774Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1774",
		Message: "`" + where + "` enables `GLOB_SUBST` — every unquoted `$var` expansion " +
			"is rescanned as a glob pattern. User-controlled data becomes a filesystem " +
			"query. Scope this in a subshell / function, or use explicit expansion flags.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1775",
		Title:    "Warn on `timeout DURATION cmd` without `--kill-after` / `-k` — hang on SIGTERM-resistant child",
		Severity: SeverityWarning,
		Description: "`timeout DURATION cmd` sends `SIGTERM` once the duration elapses and then " +
			"waits for the child to exit. A child that blocks or ignores `SIGTERM` (long-running " +
			"daemons, processes stuck in `D` state, a trapped / reset signal handler) never " +
			"dies, so the entire pipeline hangs past the intended bound. Add `--kill-after=N` " +
			"(`-k N`) so timeout escalates to `SIGKILL` after N seconds, guaranteeing exit. " +
			"Typical choice: a few seconds shorter than your CI step budget, so the overall " +
			"wait remains bounded.",
		Check: checkZC1775,
	})
}

func checkZC1775(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "timeout" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-k" || v == "--kill-after" {
			return nil
		}
		if strings.HasPrefix(v, "--kill-after=") {
			return nil
		}
	}
	return []Violation{{
		KataID: "ZC1775",
		Message: "`timeout` without `--kill-after` / `-k` only sends `SIGTERM` — a child " +
			"that blocks or ignores it hangs the pipeline past the deadline. Add " +
			"`--kill-after=N` so timeout escalates to `SIGKILL`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

// Schemes that commonly embed credentials in a connection URI and are
// passed to a CLI client that keeps the URI in argv.
var zc1776CredSchemes = []string{
	"postgres://",
	"postgresql://",
	"mysql://",
	"mariadb://",
	"mongodb://",
	"mongodb+srv://",
	"redis://",
	"rediss://",
	"amqp://",
	"amqps://",
	"kafka://",
	"nats://",
	"clickhouse://",
	"cockroachdb://",
	"db2://",
	"mssql://",
	"oracle://",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1776",
		Title:    "Error on `psql postgresql://user:secret@host/db` — password in argv via connection URI",
		Severity: SeverityError,
		Description: "Database and message-broker CLIs accept a single connection URI " +
			"(`postgresql://`, `mysql://`, `mongodb://`, `redis://`, `amqp://`, `kafka://`, " +
			"and friends). When the URI embeds a password — `scheme://user:secret@host/db` — " +
			"the secret lands in argv, visible to every user via `ps`, `/proc/PID/cmdline`, " +
			"process accounting, and audit trails, and it often survives in shell history. " +
			"Keep the password out of argv: use the client's password-file / `.pgpass` / " +
			"`PGPASSWORD` / `REDISCLI_AUTH` equivalent, or interpolate the URI from an " +
			"environment variable so the secret is not on the command line that other users " +
			"can see.",
		Check: checkZC1776,
	})
}

func checkZC1776(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if _, ok := cmd.Name.(*ast.Identifier); !ok {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		v = strings.Trim(v, "\"'")
		if leak, scheme := zc1776UriHasPassword(v); leak {
			return []Violation{{
				KataID: "ZC1776",
				Message: "`" + scheme + "user:SECRET@…` in argv puts the password in `ps` / " +
					"`/proc/PID/cmdline` / history. Use a password file (`~/.pgpass`, " +
					"`~/.my.cnf`), `PGPASSWORD` / `REDISCLI_AUTH` env var, or build the URI " +
					"from a secret variable.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1776UriHasPassword(v string) (bool, string) {
	for _, scheme := range zc1776CredSchemes {
		if !strings.HasPrefix(v, scheme) {
			continue
		}
		rest := v[len(scheme):]
		at := strings.Index(rest, "@")
		if at <= 0 {
			return false, scheme
		}
		userinfo := rest[:at]
		colon := strings.Index(userinfo, ":")
		if colon <= 0 || colon == len(userinfo)-1 {
			// No password segment, or empty password.
			return false, scheme
		}
		return true, scheme
	}
	return false, ""
}

const zc1777PreloadPath = "/etc/ld.so.preload"

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1777",
		Title:    "Error on `tee/cp/mv/install/dd` writing `/etc/ld.so.preload` — classic rootkit persistence",
		Severity: SeverityError,
		Description: "`/etc/ld.so.preload` lists shared libraries that the dynamic linker " +
			"forcibly loads into every dynamically-linked binary, root processes included. " +
			"The file is almost never needed on a modern distribution — package managers " +
			"do not touch it, and `LD_PRELOAD` handles the per-invocation case without " +
			"persisting the change. A script that pipes content into `/etc/ld.so.preload` " +
			"with `tee` / `cp` / `mv` / `install` / `dd` is a textbook rootkit persistence " +
			"primitive (`libprocesshider`, `Azazel`, `Jynx`). Remove the line, audit " +
			"`/etc/ld.so.preload` for unexpected entries (`sha256sum`, `diff` against a " +
			"known-good backup), and if preloading is legitimately required, use a scoped " +
			"`LD_PRELOAD=` on the specific invocation.",
		Check: checkZC1777,
	})
}

func checkZC1777(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "cp", "mv", "tee", "install", "dd":
		for _, arg := range cmd.Arguments {
			if arg.String() == zc1777PreloadPath {
				return zc1777Hit(cmd, ident.Value+" "+zc1777PreloadPath)
			}
		}
	}

	prevRedir := ""
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevRedir != "" {
			if v == zc1777PreloadPath {
				return zc1777Hit(cmd, prevRedir+" "+zc1777PreloadPath)
			}
			prevRedir = ""
			continue
		}
		if v == ">" || v == ">>" {
			prevRedir = v
		}
	}
	return nil
}

func zc1777Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1777",
		Message: "`" + what + "` writes `/etc/ld.so.preload` — linker force-loads " +
			"each listed library into every process. Audit for unexpected entries; " +
			"for a scoped preload use `LD_PRELOAD=` on a single invocation.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1778",
		Title:    "Warn on `systemctl link /path/to/unit` — persistence from a mutable source path",
		Severity: SeverityWarning,
		Description: "`systemctl link` symlinks the given unit file into `/etc/systemd/system/` " +
			"so it can be `enable`d and `start`ed by name, but the unit definition lives at " +
			"the original path forever. If that path is writable by any non-root user " +
			"(`/tmp/*`, `/var/tmp/*`, `/home/*`, `/opt/` with wide perms, a build output " +
			"directory), a later tamper of the source file silently changes what systemd " +
			"runs the next time the unit starts. Copy the unit into `/etc/systemd/system/` " +
			"with root-only permissions, or install a package that ships it under " +
			"`/lib/systemd/system/`, rather than linking from a mutable location.",
		Check: checkZC1778,
	})
}

var zc1778MutablePrefixes = []string{
	"/tmp/",
	"/var/tmp/",
	"/home/",
	"/root/",
	"/opt/",
	"/srv/",
	"/mnt/",
	"/media/",
	"/var/lib/",
}

func checkZC1778(node ast.Node) []Violation {
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

	linkIdx := -1
	for i, arg := range cmd.Arguments {
		if arg.String() == "link" {
			linkIdx = i
			break
		}
	}
	if linkIdx == -1 {
		return nil
	}
	for _, arg := range cmd.Arguments[linkIdx+1:] {
		v := arg.String()
		if strings.HasPrefix(v, "-") {
			continue
		}
		if !strings.HasPrefix(v, "/") {
			continue
		}
		for _, prefix := range zc1778MutablePrefixes {
			if strings.HasPrefix(v, prefix) {
				return []Violation{{
					KataID: "ZC1778",
					Message: "`systemctl link " + v + "` keeps the unit in a mutable " +
						"path — a tamper later changes what systemd runs. Copy the " +
						"unit into `/etc/systemd/system/` with root-only perms.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}

var zc1779AdminRoles = map[string]bool{
	"owner":                     true,
	"contributor":               true,
	"user access administrator": true,
	"useraccessadministrator":   true,
	"role based access control administrator": true,
	"security admin":       true,
	"global administrator": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1779",
		Title:    "Error on `az role assignment create --role Owner|Contributor|User Access Administrator`",
		Severity: SeverityError,
		Description: "`az role assignment create --role Owner` grants full control over the " +
			"target scope (subscription, resource group, resource). `Contributor` grants " +
			"everything except role assignment, and `User Access Administrator` grants the " +
			"ability to assign any role — including Owner — elsewhere in the directory. Any " +
			"of the three is effectively top-of-chain in the assigned scope. In provisioning " +
			"automation this breaks least privilege, invites blast-radius escalations, and " +
			"sidesteps any review that would flag the permission grant. Assign a narrower " +
			"built-in role (Reader, specific-service Contributor) or a custom role whose " +
			"permission list you can enumerate.",
		Check: checkZC1779,
	})
}

func checkZC1779(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "az" {
		return nil
	}
	if len(cmd.Arguments) < 4 {
		return nil
	}
	if cmd.Arguments[0].String() != "role" ||
		cmd.Arguments[1].String() != "assignment" ||
		cmd.Arguments[2].String() != "create" {
		return nil
	}

	for i, arg := range cmd.Arguments[3:] {
		v := arg.String()
		if v == "--role" {
			if i+4 < len(cmd.Arguments) {
				role := cmd.Arguments[3+i+1].String()
				role = strings.Trim(role, "\"'")
				if zc1779IsAdminRole(role) {
					return zc1779Hit(cmd, role)
				}
			}
			continue
		}
		if strings.HasPrefix(v, "--role=") {
			role := strings.TrimPrefix(v, "--role=")
			role = strings.Trim(role, "\"'")
			if zc1779IsAdminRole(role) {
				return zc1779Hit(cmd, role)
			}
		}
	}
	return nil
}

func zc1779IsAdminRole(role string) bool {
	r := strings.ToLower(strings.TrimSpace(role))
	r = strings.ReplaceAll(r, "_", " ")
	r = strings.ReplaceAll(r, "-", " ")
	r = strings.Join(strings.Fields(r), " ")
	return zc1779AdminRoles[r]
}

func zc1779Hit(cmd *ast.SimpleCommand, role string) []Violation {
	return []Violation{{
		KataID: "ZC1779",
		Message: "`az role assignment create --role " + role + "` grants a top-of-chain " +
			"role. Pick a narrower built-in role (`Reader`, specific-service Contributor) " +
			"or a custom role whose permission list you can enumerate.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

var zc1780FsProtections = map[string]string{
	"fs.protected_symlinks=0":  "symlink follow protection in sticky dirs",
	"fs.protected_hardlinks=0": "hardlink creation protection in sticky dirs",
	"fs.protected_fifos=0":     "FIFO open protection in sticky dirs",
	"fs.protected_regular=0":   "regular-file open protection in sticky dirs",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1780",
		Title:    "Warn on `sysctl -w fs.protected_symlinks=0|protected_hardlinks=0|…` — TOCTOU guard disabled",
		Severity: SeverityWarning,
		Description: "The `fs.protected_*` sysctls close a classic race: in a sticky directory " +
			"(`/tmp`, `/var/tmp`, `/dev/shm`), a non-owner cannot follow a symlink, create a " +
			"hardlink to a file they don't own, or open a FIFO / regular file they didn't " +
			"create. Those four gates block the shape of attack where a privileged program " +
			"predictably opens a `/tmp/NAME` that an attacker has already placed as a " +
			"symlink to `/etc/shadow`. Setting any of them to `0` re-enables the race across " +
			"the whole host. Leave the defaults (`1` / `2`) in place; if a specific " +
			"application legitimately needs the old behavior, run it in a mount namespace " +
			"where `/tmp` is not sticky-shared.",
		Check: checkZC1780,
	})
}

func checkZC1780(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sysctl" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if note, ok := zc1780FsProtections[v]; ok {
			return []Violation{{
				KataID: "ZC1780",
				Message: "`sysctl " + v + "` disables " + note + " — re-opens a TOCTOU " +
					"race in sticky dirs. Leave the default unless you have a specific " +
					"reason; otherwise scope the change to a mount namespace.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

var zc1781GitSubcommands = map[string]bool{
	"clone":     true,
	"fetch":     true,
	"pull":      true,
	"push":      true,
	"ls-remote": true,
	"archive":   true,
}

var zc1781GitUrlSchemes = []string{
	"https://",
	"http://",
	"git+https://",
	"git+http://",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1781",
		Title:    "Error on `git clone https://user:token@host/...` — PAT in argv and git config",
		Severity: SeverityError,
		Description: "A git remote URL in the form `https://user:token@host/path` puts the " +
			"personal access token directly in argv — visible via `ps`, `/proc/PID/cmdline`, " +
			"shell history, and process accounting. `git clone` additionally records the URL " +
			"(including the credentials) in `.git/config` as the `origin` remote, so every " +
			"later `git fetch` / `pull` re-exposes the same token to every user who can read " +
			"that file. Use a credential helper (`git credential-store`, `git credential-" +
			"osxkeychain`), `GIT_ASKPASS` with a secret pulled from an env var, HTTPS + an " +
			"SSH deploy key, or set the token via the `Authorization: Bearer` header with " +
			"`http.extraHeader` from an env-sourced value.",
		Check: checkZC1781,
	})
}

func checkZC1781(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	subIdx := -1
	for i, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "-") {
			continue
		}
		if zc1781GitSubcommands[v] {
			subIdx = i
		}
		break
	}
	if subIdx == -1 {
		return nil
	}

	for _, arg := range cmd.Arguments[subIdx+1:] {
		v := strings.Trim(arg.String(), "\"'")
		if leak := zc1781HasCredsInURL(v); leak {
			return []Violation{{
				KataID: "ZC1781",
				Message: "`git " + cmd.Arguments[subIdx].String() + " " + v + "` puts " +
					"the token in argv and `.git/config`. Use a credential helper, " +
					"`GIT_ASKPASS`, or `http.extraHeader` with an env-sourced bearer.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1781HasCredsInURL(v string) bool {
	for _, scheme := range zc1781GitUrlSchemes {
		if !strings.HasPrefix(v, scheme) {
			continue
		}
		rest := v[len(scheme):]
		at := strings.Index(rest, "@")
		if at <= 0 {
			return false
		}
		userinfo := rest[:at]
		colon := strings.Index(userinfo, ":")
		if colon <= 0 || colon == len(userinfo)-1 {
			return false
		}
		return true
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1782",
		Title:    "Error on `flatpak remote-add --no-gpg-verify` — trust chain disabled for the repo",
		Severity: SeverityError,
		Description: "A Flatpak remote without GPG verification accepts any OSTree update that " +
			"the server (or anyone on the path) cares to send. Signatures are what connect " +
			"`flatpak install FOO` to the operator that actually built `FOO` — strip them and " +
			"the install reduces to a plain HTTPS download with no identity attached. If you " +
			"genuinely need a local / air-gapped repo, sign it yourself with `ostree gpg-sign` " +
			"and add the key via `--gpg-import=KEYFILE`. Never leave `--no-gpg-verify` in " +
			"provisioning scripts for production systems.",
		Check: checkZC1782,
	})
}

func checkZC1782(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "flatpak" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "remote-add" && cmd.Arguments[0].String() != "remote-modify" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--no-gpg-verify" ||
			v == "--gpg-verify=false" ||
			v == "--no-gpg-verify=true" {
			return []Violation{{
				KataID: "ZC1782",
				Message: "`flatpak " + cmd.Arguments[0].String() + " " + v + "` disables " +
					"signature verification — updates from this remote are accepted with " +
					"only HTTPS as identity. Sign the repo (`ostree gpg-sign`) and import " +
					"the key with `--gpg-import=KEYFILE`.",
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
		ID:       "ZC1783",
		Title:    "Error on `podman system reset` / `nerdctl system prune -af --volumes` — wipes every container artifact",
		Severity: SeverityError,
		Description: "`podman system reset` removes every podman container, image, volume, " +
			"network, pod, secret, and storage driver scratch area — a full factory reset " +
			"of the local engine. `nerdctl system prune -af --volumes` achieves the same for " +
			"containerd. On a developer workstation this wipes cached images for unrelated " +
			"projects; on a CI runner or build host it invalidates every warm artifact the " +
			"job relies on; on a prod host it drops the volumes the workload stores data in. " +
			"Use narrower commands (`podman rmi`, `podman volume rm`, scoped `podman prune`) " +
			"that only touch the resource you intend to remove, and never pair the reset with " +
			"`--force`.",
		Check: checkZC1783,
	})
}

func checkZC1783(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "podman":
		if len(cmd.Arguments) >= 2 &&
			cmd.Arguments[0].String() == "system" &&
			cmd.Arguments[1].String() == "reset" {
			return zc1783Hit(cmd, "podman system reset")
		}
	case "nerdctl":
		if len(cmd.Arguments) >= 2 &&
			cmd.Arguments[0].String() == "system" &&
			cmd.Arguments[1].String() == "prune" {
			hasAll := false
			hasVolumes := false
			for _, arg := range cmd.Arguments[2:] {
				v := arg.String()
				if v == "-af" || v == "-fa" || v == "-a" || v == "--all" {
					hasAll = true
				}
				if v == "--volumes" {
					hasVolumes = true
				}
			}
			if hasAll && hasVolumes {
				return zc1783Hit(cmd, "nerdctl system prune -a --volumes")
			}
		}
	}
	return nil
}

func zc1783Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1783",
		Message: "`" + what + "` wipes every container artifact on the host — images, " +
			"volumes, networks, pods. Use narrower removals (`rmi`, `volume rm`, scoped " +
			"`prune`) against the specific resource you intend to delete.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

var zc1784HooksMutablePrefixes = []string{
	"/tmp/",
	"/var/tmp/",
	"/dev/shm/",
	"/home/",
	"/root/",
	"/opt/",
	"/srv/",
	"/mnt/",
	"/media/",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1784",
		Title:    "Warn on `git config core.hooksPath /tmp/...` — hook execution from a mutable path",
		Severity: SeverityWarning,
		Description: "`core.hooksPath` tells git which directory to run repository hooks from. " +
			"Any file named `pre-commit`, `post-checkout`, `post-merge`, etc. under that " +
			"directory becomes executable code invoked by routine git operations. Pointing " +
			"`core.hooksPath` at `/tmp`, `/var/tmp`, `/dev/shm`, `/home/<other>`, `/opt`, " +
			"`/srv`, or `/mnt` hands the git CLI an execution primitive from a path that a " +
			"non-root (or another) user can write at will — a classic supply-chain entry " +
			"point on shared hosts and CI runners. Keep hooks inside the repo's `.git/hooks/` " +
			"(or a repo-owned `.githooks/` directory) and configure `core.hooksPath` only to " +
			"paths that share the repo's owner and permissions.",
		Check: checkZC1784,
	})
}

func checkZC1784(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "config" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v != "core.hooksPath" {
			continue
		}
		// Path follows.
		if 1+i+1 >= len(cmd.Arguments) {
			return nil
		}
		path := cmd.Arguments[1+i+1].String()
		path = strings.Trim(path, "\"'")
		if !strings.HasPrefix(path, "/") {
			return nil
		}
		for _, prefix := range zc1784HooksMutablePrefixes {
			if strings.HasPrefix(path, prefix) {
				return []Violation{{
					KataID: "ZC1784",
					Message: "`git config core.hooksPath " + path + "` runs hooks from " +
						"a mutable path — supply-chain primitive. Keep hooks in the " +
						"repo's `.git/hooks/` (or a tracked `.githooks/`) and point " +
						"`core.hooksPath` at repo-owned paths only.",
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
		ID:       "ZC1785",
		Title:    "Error on `ufw default allow` — flips host firewall from deny-by-default to allow-by-default",
		Severity: SeverityError,
		Description: "`ufw default allow incoming` (or `allow outgoing`, `allow routed`) changes " +
			"the chain's baseline verdict — instead of only what you explicitly opened, every " +
			"port that does not have a matching `deny` rule is accepted. On an internet-facing " +
			"host this is effectively \"turn the firewall off\", and the effect survives reboots " +
			"because the default is persisted to `/etc/default/ufw`. Restore with `ufw default " +
			"deny incoming` and add narrow `ufw allow <port>` rules for the services that " +
			"actually need to be reachable.",
		Check: checkZC1785,
	})
}

func checkZC1785(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ufw" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "default" {
		return nil
	}

	verdict := cmd.Arguments[1].String()
	if verdict != "allow" {
		return nil
	}

	direction := "incoming"
	if len(cmd.Arguments) >= 3 {
		d := cmd.Arguments[2].String()
		if d == "incoming" || d == "outgoing" || d == "routed" {
			direction = d
		}
	}

	return []Violation{{
		KataID: "ZC1785",
		Message: "`ufw default allow " + direction + "` flips the firewall baseline to " +
			"accept every port that is not explicitly denied. Restore with `ufw default " +
			"deny incoming` and add narrow `ufw allow <port>` rules for the services that " +
			"must be reachable.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1786",
		Title:    "Error on `mount.cifs ... -o password=SECRET` — SMB password in argv",
		Severity: SeverityError,
		Description: "Passing `password=` (or `pass=`) inside `mount.cifs` / `mount -t cifs` " +
			"options puts the SMB password in argv. Any local user who can read `ps`, " +
			"`/proc/PID/cmdline`, or process-accounting records gets the cleartext, and the " +
			"line also ends up in shell history and — if captured — in CI logs. Use a " +
			"`credentials=/etc/cifs-creds` file (`0600`, `username=` and `password=` lines), " +
			"the `$USER`/`$PASSWD` env vars `mount.cifs` reads when those options are " +
			"missing, or `pam_mount` for login-time mounts.",
		Check: checkZC1786,
	})
}

func checkZC1786(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "mount.cifs":
		// direct form
	case "mount":
		isCifs := false
		for i, arg := range cmd.Arguments {
			v := arg.String()
			if (v == "-t" || v == "--types") && i+1 < len(cmd.Arguments) {
				t := cmd.Arguments[i+1].String()
				if t == "cifs" || t == "smb3" {
					isCifs = true
					break
				}
			}
		}
		if !isCifs {
			return nil
		}
	default:
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		var opts string
		if strings.HasPrefix(v, "-o") && len(v) > 2 {
			opts = v[2:]
		} else if v == "-o" && i+1 < len(cmd.Arguments) {
			opts = cmd.Arguments[i+1].String()
		}
		if opts == "" {
			continue
		}
		opts = strings.Trim(opts, "\"'")
		if zc1786OptsHavePassword(opts) {
			return []Violation{{
				KataID: "ZC1786",
				Message: "`" + ident.Value + " ... password=…` leaks the SMB password " +
					"into argv / `ps` / `/proc/PID/cmdline`. Use `credentials=/path/" +
					"to/creds` (mode 0600) or `$PASSWD` env var instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1786OptsHavePassword(opts string) bool {
	for _, field := range strings.Split(opts, ",") {
		key, value, ok := strings.Cut(strings.TrimSpace(field), "=")
		if !ok {
			continue
		}
		switch strings.ToLower(key) {
		case "password", "pass", "password2":
			if value != "" {
				return true
			}
		}
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1787",
		Title:    "Warn on `setopt AUTO_CD` — bare word that names a directory silently changes `$PWD`",
		Severity: SeverityWarning,
		Description: "With `AUTO_CD` on, any bare word that happens to name an existing directory " +
			"is executed as `cd <word>` — no command name, no error. This is a pleasant " +
			"interactive shortcut and an absolute footgun in scripts: a typo in a command " +
			"name (`dockr` → a directory called `dockr` that was left lying around) or a " +
			"user-controlled variable that expands to a path silently reshapes `$PWD` for " +
			"every later relative path. Keep `AUTO_CD` inside `~/.zshrc` where it belongs, " +
			"not in a `.zsh` script, and never turn it on inside a function that an external " +
			"caller depends on.",
		Check: checkZC1787,
	})
}

func checkZC1787(node ast.Node) []Violation {
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
			if zc1787IsAutoCd(arg.String()) {
				return zc1787Hit(cmd, "setopt "+arg.String())
			}
		}
	case "set":
		for i, arg := range cmd.Arguments {
			v := arg.String()
			if (v == "-o" || v == "--option") && i+1 < len(cmd.Arguments) {
				next := cmd.Arguments[i+1].String()
				if zc1787IsAutoCd(next) {
					return zc1787Hit(cmd, "set -o "+next)
				}
			}
		}
	}
	return nil
}

func zc1787IsAutoCd(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "AUTOCD"
}

func zc1787Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1787",
		Message: "`" + where + "` turns any bare directory name into a silent `cd`. " +
			"A typo or a user-controlled value reshapes `$PWD`; keep this in " +
			"`~/.zshrc`, not in scripts.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

var zc1788SSHMutablePrefixes = []string{
	"/tmp/",
	"/var/tmp/",
	"/dev/shm/",
	"/home/",
	"/root/",
	"/opt/",
	"/srv/",
	"/mnt/",
	"/media/",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1788",
		Title:    "Warn on `ssh -F /tmp/config` — config from a mutable path can pin `ProxyCommand` to arbitrary code",
		Severity: SeverityWarning,
		Description: "`ssh -F PATH` (and `scp -F PATH`, `sftp -F PATH`) loads a user-supplied " +
			"config file. Anything in `/etc/ssh/ssh_config` can be overridden — notably " +
			"`ProxyCommand`, `LocalCommand`, `PermitLocalCommand`, and `Include` — which means " +
			"a mutable source path is an execution primitive: another local user flips " +
			"`ProxyCommand` to `/tmp/pwn`, and the next `ssh` run launches it with the " +
			"caller's credentials and forwarded agent. Keep the config in `~/.ssh/config` (or " +
			"a repo-owned path with the same owner and `0600` perms) and never pass `-F` to " +
			"`/tmp`, `/var/tmp`, `/dev/shm`, another user's `/home`, `/opt`, `/srv`, or `/mnt`.",
		Check: checkZC1788,
	})
}

func checkZC1788(node ast.Node) []Violation {
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

	for i, arg := range cmd.Arguments {
		v := arg.String()
		var path string
		switch {
		case v == "-F":
			if i+1 >= len(cmd.Arguments) {
				return nil
			}
			path = cmd.Arguments[i+1].String()
		case strings.HasPrefix(v, "-F"):
			path = v[2:]
		default:
			continue
		}
		path = strings.Trim(path, "\"'")
		if !strings.HasPrefix(path, "/") {
			continue
		}
		for _, prefix := range zc1788SSHMutablePrefixes {
			if strings.HasPrefix(path, prefix) {
				return []Violation{{
					KataID: "ZC1788",
					Message: "`" + ident.Value + " -F " + path + "` loads an alternate " +
						"config from a mutable path — a tamper on that file can pin " +
						"`ProxyCommand` to arbitrary code. Keep the config in " +
						"`~/.ssh/config` (or a repo-owned path with `0600` perms).",
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
		ID:       "ZC1789",
		Title:    "Warn on `setopt CORRECT` / `CORRECT_ALL` — Zsh spellcheck silently rewrites script tokens",
		Severity: SeverityWarning,
		Description: "`setopt CORRECT` prompts to rewrite command names that look mistyped; " +
			"`CORRECT_ALL` extends the check to every argument on the line. In an interactive " +
			"shell this is a friendly nudge. In a script it becomes a footgun: a filename " +
			"that is *close enough* to an existing file gets silently replaced with that " +
			"other file, and the \"nlh?\" prompt reads from stdin — which may be the input " +
			"the script was supposed to process. Keep `CORRECT` / `CORRECT_ALL` in " +
			"`~/.zshrc` only and never toggle them inside a function a script calls.",
		Check: checkZC1789,
	})
}

func checkZC1789(node ast.Node) []Violation {
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
			if name, hit := zc1789Matches(arg.String()); hit {
				return zc1789Hit(cmd, "setopt "+arg.String(), name)
			}
		}
	case "set":
		for i, arg := range cmd.Arguments {
			v := arg.String()
			if (v == "-o" || v == "--option") && i+1 < len(cmd.Arguments) {
				next := cmd.Arguments[i+1].String()
				if name, hit := zc1789Matches(next); hit {
					return zc1789Hit(cmd, "set -o "+next, name)
				}
			}
		}
	}
	return nil
}

func zc1789Matches(v string) (string, bool) {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	switch norm {
	case "CORRECT":
		return "CORRECT", true
	case "CORRECTALL":
		return "CORRECT_ALL", true
	}
	return "", false
}

func zc1789Hit(cmd *ast.SimpleCommand, where, canonical string) []Violation {
	return []Violation{{
		KataID: "ZC1789",
		Message: "`" + where + "` enables `" + canonical + "` — Zsh spellcheck " +
			"silently rewrites tokens that look mistyped. In a script that corrupts " +
			"file paths and steals stdin for the correction prompt. Keep in `~/.zshrc`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1790",
		Title:    "Warn on `unsetopt PIPE_FAIL` — pipeline exit status reverts to last-command-only",
		Severity: SeverityWarning,
		Description: "With `PIPE_FAIL` off (the shell default), `cmd1 | cmd2 | cmd3` exits with " +
			"`cmd3`'s status; failures in `cmd1` and `cmd2` are silently dropped. " +
			"`unsetopt PIPE_FAIL` (or the equivalent `setopt NOPIPEFAIL`) mid-script turns a " +
			"previously-enabled error check back off — typically because a known-flaky pipe " +
			"stage was tripping `set -e`, and the author reached for the global off-switch. " +
			"Undo the change in a subshell (`( unsetopt pipefail; …; )`) or a function with " +
			"`emulate -L zsh; unsetopt pipefail` so the rest of the script keeps strict-pipe " +
			"error propagation.",
		Check: checkZC1790,
	})
}

func checkZC1790(node ast.Node) []Violation {
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
			if zc1790IsPipeFail(arg.String()) {
				return zc1790Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOPIPEFAIL" {
				return zc1790Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1790IsPipeFail(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "PIPEFAIL"
}

func zc1790Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1790",
		Message: "`" + where + "` returns the shell to last-command-only pipeline exit — " +
			"`cmd1 | cmd2` now ignores `cmd1` failures. Scope the change to a subshell " +
			"or function with `emulate -L zsh` instead of flipping it globally.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

var zc1791DaemonSockets = []string{
	"/var/run/docker.sock",
	"/run/docker.sock",
	"/var/run/podman/podman.sock",
	"/run/podman/podman.sock",
	"/run/containerd/containerd.sock",
	"/run/crio/crio.sock",
	"/var/run/docker/containerd/containerd.sock",
}

var zc1791UnixSocketFlags = map[string]bool{"--unix-socket": true}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1791",
		Title:    "Error on `curl --unix-socket /var/run/docker.sock` — direct container-daemon API access",
		Severity: SeverityError,
		Description: "A curl request to `docker.sock` / `containerd.sock` / `crio.sock` speaks " +
			"the container-daemon HTTP API with no authentication beyond the socket's " +
			"filesystem permissions. Anyone who can invoke curl as that uid can `POST " +
			"/containers/create` with `HostConfig.Privileged=true` and a bind mount of `/` " +
			"and land a root shell on the host — the primitive every \"docker socket " +
			"escape\" write-up leans on. Use the real CLI (`docker`, `podman`, `nerdctl`) " +
			"which enforces its own policy, or access the daemon over a TLS-protected TCP " +
			"endpoint with mutual auth.",
		Check: checkZC1791,
	})
}

func checkZC1791(node ast.Node) []Violation {
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

	for i, arg := range cmd.Arguments {
		v := arg.String()
		var path string
		switch {
		case v == "--unix-socket":
			if i+1 >= len(cmd.Arguments) {
				return nil
			}
			path = cmd.Arguments[i+1].String()
		case strings.HasPrefix(v, "--unix-socket="):
			path = strings.TrimPrefix(v, "--unix-socket=")
		default:
			continue
		}
		path = strings.Trim(path, "\"'")
		if hit := zc1791MatchSocket(cmd, path); hit != nil {
			return hit
		}
	}
	return nil
}

func zc1791MatchSocket(cmd *ast.SimpleCommand, path string) []Violation {
	for _, sock := range zc1791DaemonSockets {
		if path == sock {
			line, col := FlagArgPosition(cmd, zc1791UnixSocketFlags)
			return []Violation{{
				KataID: "ZC1791",
				Message: "`curl --unix-socket " + path + "` speaks the container-daemon " +
					"API — a `POST /containers/create` with `Privileged=true` is a " +
					"host-root primitive. Use the CLI (`docker`/`podman`) instead.",
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
		ID:       "ZC1792",
		Title:    "Warn on `btrfs subvolume delete` / `btrfs device remove` — unrecoverable btrfs data loss",
		Severity: SeverityWarning,
		Description: "`btrfs subvolume delete PATH` unlinks the subvolume and drops all of its " +
			"extents once cleanup completes — on Snapper / Timeshift systems the argument is " +
			"often a snapshot that is the only remaining copy of pre-incident state. " +
			"`btrfs device remove DEV POOL` moves the stored chunks off DEV before detaching " +
			"it; wrong device, mid-rebalance failure, or insufficient free space across the " +
			"remaining members puts the filesystem into degraded mode with no automatic " +
			"rollback. Keep a fresh `btrfs subvolume list`/`btrfs device usage` snapshot and " +
			"confirm the target explicitly before running either command in automation.",
		Check: checkZC1792,
	})
}

func checkZC1792(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "btrfs" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	sub0 := cmd.Arguments[0].String()
	sub1 := cmd.Arguments[1].String()

	switch {
	case sub0 == "subvolume" && sub1 == "delete":
		return zc1792Hit(cmd, "btrfs subvolume delete")
	case sub0 == "device" && sub1 == "remove":
		return zc1792Hit(cmd, "btrfs device remove")
	}
	return nil
}

func zc1792Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1792",
		Message: "`" + what + "` drops btrfs state with no automatic rollback — " +
			"snapshots vanish on `subvolume delete`, and `device remove` can leave " +
			"the filesystem degraded. Confirm the target explicitly.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1793",
		Title:    "Warn on `kubectl certificate approve CSR` — signs the identity baked into the CSR",
		Severity: SeverityWarning,
		Description: "`kubectl certificate approve NAME` tells the cluster signer to sign the " +
			"pending CSR unchanged. The signer respects the Subject (CN, O) and the " +
			"SubjectAltName extensions the caller put in the CSR — approve one that requests " +
			"`system:masters` and you have handed the requester full admin on the cluster. " +
			"In automation, review the CSR body first (`kubectl get csr NAME -o " +
			"jsonpath='{.spec.request}' | base64 -d | openssl req -text`) and reject (`kubectl " +
			"certificate deny`) any request that names a privileged group, kube-system service " +
			"account, or hostname outside the intended scope.",
		Check: checkZC1793,
	})
}

func checkZC1793(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubectl" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "certificate" {
		return nil
	}
	if cmd.Arguments[1].String() != "approve" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1793",
		Message: "`kubectl certificate approve` signs the identity embedded in the CSR " +
			"— a `system:masters` request becomes cluster admin. Decode with " +
			"`openssl req -text` first; use `kubectl certificate deny` otherwise.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

var zc1794CosignSkipFlags = map[string]bool{
	"--insecure-ignore-tlog":    true,
	"--insecure-ignore-sct":     true,
	"--insecure-skip-verify":    true,
	"--allow-insecure-registry": true,
	"--allow-http-registry":     true,
	"--allow-insecure-bundle":   true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1794",
		Title:    "Error on `cosign verify --insecure-ignore-tlog` / `--allow-insecure-registry` — signature chain disabled",
		Severity: SeverityError,
		Description: "`cosign verify` with `--insecure-ignore-tlog` skips Rekor transparency-log " +
			"verification, `--insecure-ignore-sct` skips Fulcio SCT verification, and " +
			"`--insecure-skip-verify` turns off TLS certificate validation for the registry / " +
			"Rekor / Fulcio endpoints. `cosign sign --allow-insecure-registry` and " +
			"`--allow-http-registry` push signatures over plain HTTP. Each flag removes a " +
			"distinct rung of the signature chain that `cosign` was built to enforce — a " +
			"malicious registry or on-path attacker now passes verification without detection. " +
			"Drop the flag, fix the underlying trust anchor (CA bundle, Rekor URL, " +
			"Fulcio OIDC), and keep signature verification strict.",
		Check: checkZC1794,
	})
}

func checkZC1794(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cosign" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		flag := v
		if idx := strings.Index(flag, "="); idx >= 0 {
			flag = flag[:idx]
		}
		if zc1794CosignSkipFlags[flag] {
			return []Violation{{
				KataID: "ZC1794",
				Message: "`cosign " + v + "` removes a rung of the signature chain " +
					"(transparency log / SCT / TLS / HTTPS-only registry). Drop " +
					"the flag and fix the underlying trust anchor.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

var zc1795RemoteActions = map[string]bool{
	"add":     true,
	"set-url": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1795",
		Title:    "Error on `git remote add NAME https://user:token@host/repo` — credentials persisted in `.git/config`",
		Severity: SeverityError,
		Description: "`git remote add NAME URL` and `git remote set-url NAME URL` write the URL " +
			"into `.git/config` verbatim. When the URL embeds a `user:token@host` credential " +
			"segment, every reader of the repo — other local users, a compromised backup, a " +
			"CI cache, or anyone who runs `git config --list` — picks up the secret. It also " +
			"shows up in argv at the moment of creation (visible via `ps` / " +
			"`/proc/PID/cmdline`). Use a credential helper (`git credential-store`, " +
			"`credential-osxkeychain`), `GIT_ASKPASS` sourced from an env var, or HTTPS + a " +
			"deploy SSH key.",
		Check: checkZC1795,
	})
}

func checkZC1795(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) < 4 {
		return nil
	}
	if cmd.Arguments[0].String() != "remote" {
		return nil
	}
	if !zc1795RemoteActions[cmd.Arguments[1].String()] {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		v := strings.Trim(arg.String(), "\"'")
		if zc1795UrlHasCreds(v) {
			return []Violation{{
				KataID: "ZC1795",
				Message: "`git remote " + cmd.Arguments[1].String() + " … " + v + "` " +
					"stores the token in `.git/config` and leaks it via argv at " +
					"creation. Use a credential helper, `GIT_ASKPASS`, or an SSH " +
					"deploy key instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1795UrlHasCreds(v string) bool {
	for _, scheme := range []string{"https://", "http://", "git+https://", "git+http://"} {
		if !strings.HasPrefix(v, scheme) {
			continue
		}
		rest := v[len(scheme):]
		at := strings.Index(rest, "@")
		if at <= 0 {
			return false
		}
		userinfo := rest[:at]
		colon := strings.Index(userinfo, ":")
		if colon <= 0 || colon == len(userinfo)-1 {
			return false
		}
		return true
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1796",
		Title:    "Warn on `pg_restore --clean` / `-c` — drops existing DB objects before recreating",
		Severity: SeverityWarning,
		Description: "`pg_restore -c` (also `--clean`) issues `DROP` for every table, index, " +
			"function, and sequence in the target database before recreating them from the " +
			"archive. If the backup is stale, incomplete, or points at the wrong database, " +
			"the destination loses any object that isn't in the dump — including data added " +
			"after the backup ran. Restore into a fresh empty database (`createdb new && " +
			"pg_restore -d new`) or snapshot the target (`pg_dump -Fc > pre.dump`) before " +
			"running `--clean`, and never pair it with `--if-exists` on a live production DB " +
			"without a tested rollback path.",
		Check: checkZC1796,
	})
}

func checkZC1796(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `pg_restore --clean …` mangles to name=`clean`.
	// To avoid false positives on unrelated commands with `clean` as name,
	// require another pg_restore-ish argument to be present.
	if ident.Value == "clean" {
		if zc1796HasPgArg(cmd) {
			return zc1796Hit(cmd, "pg_restore --clean")
		}
		return nil
	}

	if ident.Value != "pg_restore" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c" || v == "--clean" {
			return zc1796Hit(cmd, "pg_restore "+v)
		}
	}
	return nil
}

func zc1796HasPgArg(cmd *ast.SimpleCommand) bool {
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-d", "--dbname", "-F", "--format", "-U", "--username",
			"--if-exists", "--no-owner", "--no-acl":
			return true
		}
	}
	return false
}

func zc1796Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1796",
		Message: "`" + what + "` drops every object in the target DB before recreating " +
			"from the archive — stale or wrong-target dump silently loses data. Restore " +
			"into a fresh DB (`createdb new && pg_restore -d new`), or snapshot first.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1797",
		Title:    "Warn on `ip link set <iface> down` / `ifdown <iface>` — locks out remote admin on that path",
		Severity: SeverityWarning,
		Description: "Taking a network interface down from an SSH session that rides on the same " +
			"interface cuts the script off mid-run: the TCP connection freezes, any later " +
			"step silently fails, and recovery requires console / out-of-band access. " +
			"Common bugs are typos (`eth1` instead of `eth0`), scripts that target the only " +
			"uplink on a cloud VM, or running the command without first confirming that the " +
			"interface is not the one carrying the admin session. Wrap the `down` in a " +
			"`systemd-run --on-active=30s --unit=recover ip link set <iface> up` rollback, " +
			"or stage both `down` and `up` through `nmcli connection up/down` with a pinned " +
			"fallback profile.",
		Check: checkZC1797,
	})
}

func checkZC1797(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "ip":
		// Expect: ip link set IFACE down
		if len(cmd.Arguments) < 4 {
			return nil
		}
		if cmd.Arguments[0].String() != "link" || cmd.Arguments[1].String() != "set" {
			return nil
		}
		for _, arg := range cmd.Arguments[2:] {
			if arg.String() == "down" {
				return zc1797Hit(cmd, "ip link set … down")
			}
		}
	case "ifdown":
		if len(cmd.Arguments) == 0 {
			return nil
		}
		first := cmd.Arguments[0].String()
		if first == "--help" || first == "-h" {
			return nil
		}
		return zc1797Hit(cmd, "ifdown "+first)
	}
	return nil
}

func zc1797Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1797",
		Message: "`" + what + "` disables a network interface — if it carries the " +
			"SSH session, the script cuts itself off. Schedule a rollback via " +
			"`systemd-run --on-active=30s ip link set … up` or stage via `nmcli`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1798",
		Title:    "Warn on `ufw reset` — wipes every configured firewall rule",
		Severity: SeverityWarning,
		Description: "`ufw reset` returns the firewall to the distro default: every user-defined " +
			"rule is removed, default incoming policy reverts (usually to `deny`, but the net " +
			"effect is the loss of every allow-list entry the host relied on). Paired with " +
			"`--force`, no prompt is issued. In a provisioning script the operation is " +
			"sometimes desired to start from a clean slate, but running it mid-session or on " +
			"a host that currently serves traffic drops connections without warning. Snapshot " +
			"the rules first (`ufw status numbered > /tmp/ufw.bak`), and prefer removing " +
			"specific rules with `ufw delete <num>` over a full reset.",
		Check: checkZC1798,
	})
}

func checkZC1798(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `ufw --force reset` mangles to name=`force` with `reset` as arg[0].
	if ident.Value == "force" {
		if len(cmd.Arguments) > 0 && cmd.Arguments[0].String() == "reset" {
			return zc1798Hit(cmd, "ufw --force reset")
		}
		return nil
	}

	if ident.Value != "ufw" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	if cmd.Arguments[0].String() != "reset" {
		return nil
	}
	return zc1798Hit(cmd, "ufw reset")
}

func zc1798Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1798",
		Message: "`" + what + "` drops every user-defined firewall rule. Snapshot " +
			"(`ufw status numbered > /tmp/ufw.bak`) first, and prefer " +
			"`ufw delete <num>` for targeted removals.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1799",
		Title:    "Warn on `rclone sync SRC DST` without `--dry-run` — one-way mirror can wipe DST",
		Severity: SeverityWarning,
		Description: "`rclone sync` makes DST look exactly like SRC: anything in DST that isn't in " +
			"SRC is deleted, including object versions on providers that support them. If SRC " +
			"is accidentally empty (typo in path, unmounted drive, wrong credentials " +
			"pointing at an empty bucket), the command silently wipes every object under DST " +
			"without a confirmation prompt. Always preview the diff with `rclone sync " +
			"--dry-run SRC DST` first; when you commit to the sync, keep `--backup-dir`, " +
			"`--max-delete`, or `--min-age` guards so a bad SRC cannot cascade into " +
			"unbounded deletion.",
		Check: checkZC1799,
	})
}

func checkZC1799(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "rclone" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "sync" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--dry-run" || v == "-n" || v == "--interactive" || v == "-i" {
			return nil
		}
	}
	return []Violation{{
		KataID: "ZC1799",
		Message: "`rclone sync` deletes anything in DST that's not in SRC — empty / " +
			"wrong SRC silently wipes DST. Preview with `rclone sync --dry-run`, and " +
			"pin guards like `--backup-dir`, `--max-delete`, or `--min-age`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
