// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1500(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — systemctl status sshd",
			input:    `systemctl status sshd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — systemctl edit --no-edit sshd",
			input:    `systemctl edit --no-edit sshd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — systemctl edit sshd",
			input: `systemctl edit sshd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1500",
					Message: "`systemctl edit` opens $EDITOR and waits for the user. Use a drop-in `/etc/systemd/system/<unit>.d/*.conf` + `daemon-reload` in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — systemctl edit --full myapp.service",
			input: `systemctl edit --full myapp.service`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1500",
					Message: "`systemctl edit` opens $EDITOR and waits for the user. Use a drop-in `/etc/systemd/system/<unit>.d/*.conf` + `daemon-reload` in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1500")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1501(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker compose up",
			input:    `docker compose up -d`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker-compose up",
			input: `docker-compose up -d`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1501",
					Message: "`docker-compose` is the deprecated Python V1 binary. Use `docker compose` (space-separated subcommand) for the bundled V2 plugin.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker-compose down",
			input: `docker-compose down --volumes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1501",
					Message: "`docker-compose` is the deprecated Python V1 binary. Use `docker compose` (space-separated subcommand) for the bundled V2 plugin.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1501")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1502(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — grep literal pattern",
			input:    `grep "hello" file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — grep -- \"$var\" (end-of-flags marker)",
			input:    `grep -- "$pattern" file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — grep \"$var\" file (no --)",
			input: `grep "$pattern" file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1502",
					Message: "Variable `\"$pattern\"` used as pattern without `--` end-of-flags marker — attacker-controlled leading `-` becomes a flag. Write `grep -- \"$var\"`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — rg $pattern files (unquoted)",
			input: `rg $pattern files`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1502",
					Message: "Variable `$pattern` used as pattern without `--` end-of-flags marker — attacker-controlled leading `-` becomes a flag. Write `grep -- \"$var\"`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1502")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1503(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — groupadd mygroup",
			input:    `groupadd mygroup`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — groupadd -g 2000 mygroup",
			input:    `groupadd -g 2000 mygroup`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — groupadd -g 0 fakeroot",
			input: `groupadd -g 0 fakeroot`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1503",
					Message: "Creating a group with GID 0 duplicates the `root` group — hidden privesc. Pick an unused GID (see `getent group`) and scope via sudoers/polkit.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — groupmod -g0 service",
			input: `groupmod -g0 service`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1503",
					Message: "Creating a group with GID 0 duplicates the `root` group — hidden privesc. Pick an unused GID (see `getent group`) and scope via sudoers/polkit.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1503")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1504(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — git push origin main",
			input:    `git push origin main`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — git push --all origin",
			input:    `git push --all origin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — git push --mirror origin",
			input: `git push --mirror origin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1504",
					Message: "`git push --mirror` overwrites every remote ref and deletes ones missing locally. Use an explicit refspec or `--all` for everyday pushes.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — git push origin --mirror",
			input: `git push origin --mirror`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1504",
					Message: "`git push --mirror` overwrites every remote ref and deletes ones missing locally. Use an explicit refspec or `--all` for everyday pushes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1504")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1505(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dpkg -i pkg.deb",
			input:    `dpkg -i pkg.deb`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dpkg -i --force-confnew pkg.deb",
			input: `dpkg -i --force-confnew pkg.deb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1505",
					Message: "`--force-confnew` silently picks maintainer or local conffile — legit /etc changes disappear or new defaults are ignored. Use ucf/etckeeper.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — dpkg -i pkg.deb --force-confold",
			input: `dpkg -i pkg.deb --force-confold`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1505",
					Message: "`--force-confold` silently picks maintainer or local conffile — legit /etc changes disappear or new defaults are ignored. Use ucf/etckeeper.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1505")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1506(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sg audio -c cmd",
			input:    `sg audio -c 'ls /var/log'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — newgrp audio",
			input: `newgrp audio`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1506",
					Message: "`newgrp` starts a new shell — script either hangs or exits. Use `sg <group> -c <cmd>` or systemd `SupplementaryGroups=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1506")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1507(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — rsync without archive / -l",
			input:    `rsync -rv src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — rsync -a --safe-links src/ dst/",
			input:    `rsync -a --safe-links src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — rsync --no-links src/ dst/",
			input:    `rsync -a --no-links src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — rsync -a src/ dst/",
			input: `rsync -a src/ dst/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1507",
					Message: "`rsync` preserving symlinks without `--safe-links` follows ones pointing outside the source tree — path traversal vector. Add `--safe-links` or `--copy-unsafe-links`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — rsync -al src/ dst/",
			input: `rsync -al src/ dst/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1507",
					Message: "`rsync` preserving symlinks without `--safe-links` follows ones pointing outside the source tree — path traversal vector. Add `--safe-links` or `--copy-unsafe-links`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1507")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1508(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — objdump -p",
			input:    `objdump -p /bin/ls`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — readelf -d",
			input:    `readelf -d /bin/ls`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ldd /bin/ls",
			input: `ldd /bin/ls`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1508",
					Message: "`ldd` on glibc can execute the target binary. Use `objdump -p` or `readelf -d` to inspect ELF dependencies safely.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ldd /tmp/downloaded.bin",
			input: `ldd /tmp/downloaded.bin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1508",
					Message: "`ldd` on glibc can execute the target binary. Use `objdump -p` or `readelf -d` to inspect ELF dependencies safely.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1508")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1509(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — trap 'cleanup' EXIT",
			input:    `trap 'cleanup' EXIT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — trap 'cleanup' TERM",
			input:    `trap 'cleanup' TERM`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — trap '' TERM",
			input: `trap '' TERM`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1509",
					Message: "`trap '' TERM` silences a fatal signal — cleanup handlers never run. Keep at least a cleanup trap on EXIT.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — trap - SIGINT",
			input: `trap - SIGINT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1509",
					Message: "`trap - SIGINT` silences a fatal signal — cleanup handlers never run. Keep at least a cleanup trap on EXIT.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1509")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1510(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — auditctl -w /etc/passwd",
			input:    `auditctl -w /etc/passwd -p wa`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — auditctl -e 1",
			input:    `auditctl -e 1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — auditctl -e 0",
			input: `auditctl -e 0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1510",
					Message: "`auditctl -e 0` disables audit subsystem — anti-forensics tactic. Use `-e 2` for a reboot-locked maintenance window instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — auditctl -D",
			input: `auditctl -D`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1510",
					Message: "`auditctl -D` deletes every audit rule — anti-forensics tactic. Use `-e 2` for a reboot-locked maintenance window instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1510")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1511(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — nmcli con up myssid",
			input:    `nmcli con up myssid`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — nmcli --ask",
			input:    `nmcli con up myssid --ask`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — nmcli con mod myssid 802-11-wireless-security.psk mypassword",
			input: `nmcli con mod myssid 802-11-wireless-security.psk mypassword`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1511",
					Message: "`nmcli` passed `802-11-wireless-security.psk <secret>` on the command line — ends up in ps/history. Use `--ask` or a keyfile profile.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — nmcli con mod vpn vpn.secrets.password pw",
			input: `nmcli con mod myvpn vpn.secrets.password vpnpass`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1511",
					Message: "`nmcli` passed `vpn.secrets.password <secret>` on the command line — ends up in ps/history. Use `--ask` or a keyfile profile.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1511")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1512(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — systemctl restart sshd",
			input:    `systemctl restart sshd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — service with unrecognized verb",
			input:    `service --help`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — service sshd restart",
			input: `service sshd restart`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1512",
					Message: "`service sshd restart` — prefer `systemctl restart sshd` for consistency with other systemd commands.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — service nginx reload",
			input: `service nginx reload`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1512",
					Message: "`service nginx reload` — prefer `systemctl reload nginx` for consistency with other systemd commands.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1512")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1513(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — make",
			input:    `make`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — make DESTDIR=/tmp/pkg install",
			input:    `make DESTDIR=/tmp/pkg install`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — make install",
			input: `make install`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1513",
					Message: "`make install` without `DESTDIR=` leaves no package-manager record. Set `DESTDIR=/tmp/pkgroot` and wrap in checkinstall / fpm, or use stow.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — gmake install",
			input: `gmake install`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1513",
					Message: "`make install` without `DESTDIR=` leaves no package-manager record. Set `DESTDIR=/tmp/pkgroot` and wrap in checkinstall / fpm, or use stow.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1513")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1514(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — useradd alice",
			input:    `useradd alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — useradd -p hash alice",
			input: `useradd -p $6$salt$hashhashhash alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1514",
					Message: "`useradd -p <hash>` puts the hashed password in ps / /proc / history. Use `chpasswd --crypt-method=SHA512` from stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — usermod -p hash bob",
			input: `usermod -p $6$salt$hashhash bob`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1514",
					Message: "`usermod -p <hash>` puts the hashed password in ps / /proc / history. Use `chpasswd --crypt-method=SHA512` from stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1514")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1515(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sha256sum file",
			input:    `sha256sum file.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — b2sum file",
			input:    `b2sum file.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — md5sum file",
			input: `md5sum file.tar.gz`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1515",
					Message: "`md5sum` is collision-vulnerable — don't use for integrity checks. Use `sha256sum` / `sha512sum` / `b2sum` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sha1sum file",
			input: `sha1sum file.tar.gz`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1515",
					Message: "`sha1sum` is collision-vulnerable — don't use for integrity checks. Use `sha256sum` / `sha512sum` / `b2sum` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1515")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1516(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — umask 022",
			input:    `umask 022`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — umask 077",
			input:    `umask 077`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — umask 000 (parser normalizes to 0)",
			input: `umask 000`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1516",
					Message: "`umask 0` leaves new files world-readable and world-writable. Use `022` for public software, `077` for secrets handling.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — umask 0",
			input: `umask 0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1516",
					Message: "`umask 0` leaves new files world-readable and world-writable. Use `022` for public software, `077` for secrets handling.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1516")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1517(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — print -P with literal",
			input:    `print -P "%F{red}hello%f"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — print without -P",
			input:    `print "$var"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — print -P with single-quoted var (no interpolation)",
			input:    `print -P '$var'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — print -P \"$var\"",
			input: `print -P "$var"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1517",
					Message: "`print -P \"$var\"` expands prompt escapes inside the variable — use `${(V)var}` / `${(q-)var}` to neutralize metacharacters, or drop -P.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — print -P $msg (unquoted)",
			input: `print -P $msg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1517",
					Message: "`print -P $msg` expands prompt escapes inside the variable — use `${(V)var}` / `${(q-)var}` to neutralize metacharacters, or drop -P.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1517")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1518(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — bash -c 'cmd'",
			input:    `bash -c 'true'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — bash -p -c 'cmd'",
			input: `bash -p -c 'true'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1518",
					Message: "`bash -p` keeps the privileged environment on a setuid wrapper — almost never needed, audit and remove.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1518")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1519(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ulimit -u 4096",
			input:    `ulimit -u 4096`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ulimit -n unlimited (different limit)",
			input:    `ulimit -n unlimited`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ulimit -u unlimited",
			input: `ulimit -u unlimited`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1519",
					Message: "`ulimit -u unlimited` removes the user process cap — fork bomb surface. Pick a realistic number or set it via /etc/security/limits.d/.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1519")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1520(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — read varname",
			input:    `read varname`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — vared myvar",
			input: `vared myvar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1520",
					Message: "`vared` requires a TTY — in a non-interactive script it errors or hangs. Use `read`, stdin, or environment variables for scripted input.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1520")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1521(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — strace -e trace=openat cmd",
			input:    `strace -e trace=openat ls`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — strace -f cmd",
			input: `strace -f ls`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1521",
					Message: "`strace` without `-e` captures every syscall including secrets in read/write buffers. Scope with `-e trace=<set>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — strace cmd (bare)",
			input: `strace ls`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1521",
					Message: "`strace` without `-e` captures every syscall including secrets in read/write buffers. Scope with `-e trace=<set>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1521")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1522(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ip route add 10.0.0.0/24 dev eth1",
			input:    `ip route add 10.0.0.0/24 dev eth1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ip route show default",
			input:    `ip route show default`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ip route add default via 1.2.3.4",
			input: `ip route add default via 1.2.3.4`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1522",
					Message: "`ip route add default` silently reroutes every non-local packet through the new gateway — MITM surface or CI footgun. Use NetworkManager/networkd config.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — route add default gw 1.2.3.4",
			input: `route add default gw 1.2.3.4`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1522",
					Message: "`route add default` silently reroutes every non-local packet through the new gateway — MITM surface or CI footgun. Use NetworkManager/networkd config.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1522")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1523(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tar xf foo.tar -C /tmp/stage",
			input:    `tar xf foo.tar -C /tmp/stage`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tar xf foo.tar -C /",
			input: `tar xf foo.tar -C /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1523",
					Message: "`tar -C /` extracts into the filesystem root — overwrites any path that happens to be inside the archive. Stage, inspect, then copy.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1523")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1524(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sysctl -p",
			input:    `sysctl -p /etc/sysctl.d/99-hardening.conf`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sysctl -e -p",
			input: `sysctl -e -p /etc/sysctl.d/99-hardening.conf`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1524",
					Message: "`sysctl -e` suppresses error output — typos in sysctl.d/ conffiles silently skip. Remove and surface the real error.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sysctl -q",
			input: `sysctl -q -p /etc/sysctl.d/99-hardening.conf`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1524",
					Message: "`sysctl -q` suppresses error output — typos in sysctl.d/ conffiles silently skip. Remove and surface the real error.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1524")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1525(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ping -c 4 host",
			input:    `ping -c 4 example.com`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ping -f host",
			input: `ping -f example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1525",
					Message: "`ping -f` (flood) bypasses the rate limit — saturates slow links. Scope tightly and document.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ping6 -f host",
			input: `ping6 -f 2001:db8::1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1525",
					Message: "`ping6 -f` (flood) bypasses the rate limit — saturates slow links. Scope tightly and document.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1525")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1526(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — wipefs --no-act",
			input:    `wipefs --no-act $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — wipefs -a disk",
			input: `wipefs -a $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1526",
					Message: "`wipefs -a` erases every filesystem signature — unrecoverable. Run with `--no-act` first, or use `sgdisk --zap-all` for scoped deletion.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — wipefs -af disk",
			input: `wipefs -af $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1526",
					Message: "`wipefs -a` erases every filesystem signature — unrecoverable. Run with `--no-act` first, or use `sgdisk --zap-all` for scoped deletion.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1526")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1527(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — crontab -l",
			input:    `crontab -l`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — crontab file",
			input:    `crontab /etc/cron.d/myfile`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — crontab -",
			input: `crontab -`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1527",
					Message: "`crontab -` overwrites the user's crontab from stdin — silently drops existing rows. Use /etc/cron.d/ files or a diff/merge workflow.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — crontab -u svc -",
			input: `crontab -u svc -`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1527",
					Message: "`crontab -` overwrites the user's crontab from stdin — silently drops existing rows. Use /etc/cron.d/ files or a diff/merge workflow.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1527")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1528(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chage -M 90 alice",
			input:    `chage -M 90 alice`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — chage -l alice",
			input:    `chage -l alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chage -M 99999 alice",
			input: `chage -M 99999 alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1528",
					Message: "`chage -M 99999` disables password aging — removes automatic lockout. Use a PAM profile instead of per-user chage.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chage -E -1 alice",
			input: `chage -E -1 alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1528",
					Message: "`chage -E -1` disables password aging — removes automatic lockout. Use a PAM profile instead of per-user chage.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1528")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1529(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — fsck -n $disk (dry run)",
			input:    `fsck -n $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — fsck -p $disk (preen)",
			input:    `fsck -p $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — fsck -y $disk",
			input: `fsck -y $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1529",
					Message: "`fsck -y` answers yes to every repair prompt — can destroy salvageable data. Prefer `-n` (dry-run) or `-p` (preen).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — fsck.ext4 -y $disk",
			input: `fsck.ext4 -y $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1529",
					Message: "`fsck.ext4 -y` answers yes to every repair prompt — can destroy salvageable data. Prefer `-n` (dry-run) or `-p` (preen).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1529")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1530(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pkill sshd (process name)",
			input:    `pkill sshd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — pkill -U 1000 java",
			input:    `pkill -U 1000 java`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pkill -f server",
			input: `pkill -f server`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1530",
					Message: "`pkill -f` matches the full command line — easy to over-kill. Drop `-f`, scope with `-U/-G/-P`, or anchor the pattern with ^/$.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1530")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1531(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — wget -t 5 https://host",
			input:    `wget -t 5 https://host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — wget -t 0 https://host",
			input: `wget -t 0 https://host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1531",
					Message: "`wget -t 0` retries forever — script hangs on dead endpoint. Use finite `-t 5` plus `--timeout=<seconds>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1531")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1532(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — screen -S mysession",
			input:    `screen -S mysession`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tmux attach",
			input:    `tmux attach-session -t work`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — screen -S name -dm cmd",
			input: `screen -S work -dm /usr/local/bin/worker`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1532",
					Message: "`screen -dm` backgrounds work outside systemd — no journal, no cgroup, common persistence technique. Use a systemd unit instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tmux new-session -d -s name cmd",
			input: `tmux new-session -d -s work /usr/local/bin/worker`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1532",
					Message: "`tmux new-session -d` backgrounds work outside systemd — no journal, no cgroup, common persistence technique. Use a systemd unit instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1532")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1533(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid — setsid /usr/local/bin/worker",
			input: `setsid /usr/local/bin/worker`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1533",
					Message: "`setsid` detaches the child from the TTY / session — escapes supervision. Prefer a systemd unit; document a detach if one is genuinely needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — setsid -f /usr/local/bin/worker",
			input: `setsid -f /usr/local/bin/worker`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1533",
					Message: "`setsid` detaches the child from the TTY / session — escapes supervision. Prefer a systemd unit; document a detach if one is genuinely needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1533")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1534(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dmesg -T",
			input:    `dmesg -T`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dmesg -c",
			input: `dmesg -c`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1534",
					Message: "`dmesg -c` wipes the kernel ring buffer — subsequent readers see no OOM/panic/audit messages. Read without clearing.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — dmesg -C",
			input: `dmesg -C`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1534",
					Message: "`dmesg -C` wipes the kernel ring buffer — subsequent readers see no OOM/panic/audit messages. Read without clearing.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1534")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1535(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ip link set eth0 up",
			input:    `ip link set eth0 up`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ip link set eth0 promisc off",
			input:    `ip link set eth0 promisc off`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ip link set eth0 promisc on",
			input: `ip link set eth0 promisc on`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1535",
					Message: "Interface put into promiscuous mode — sniffer-in-place. Re-disable after capture, or grant tcpdump CAP_NET_RAW instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ifconfig eth0 promisc",
			input: `ifconfig eth0 promisc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1535",
					Message: "Interface put into promiscuous mode — sniffer-in-place. Re-disable after capture, or grant tcpdump CAP_NET_RAW instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1535")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1536(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — iptables -L (list)",
			input:    `iptables -L`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — iptables -A INPUT -j DROP",
			input:    `iptables -A INPUT -j DROP`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — iptables -I PREROUTING ... -j DNAT",
			input: `iptables -t nat -I PREROUTING -p tcp -j DNAT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1536",
					Message: "`iptables -j DNAT` rewrites packet destination — silent redirect surface. Use declarative nftables/firewalld config.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — iptables -A OUTPUT ... -j REDIRECT",
			input: `iptables -t nat -A OUTPUT -p tcp -j REDIRECT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1536",
					Message: "`iptables -j REDIRECT` rewrites packet destination — silent redirect surface. Use declarative nftables/firewalld config.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1536")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1537(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — lvremove vg/lv",
			input:    `lvremove vg0/lv0`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — lvremove -f vg/lv",
			input: `lvremove -f vg0/lv0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1537",
					Message: "`lvremove -f` skips the confirmation — a typo in the volume name destroys every filesystem on top of it.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — vgremove -f vg",
			input: `vgremove -f vg0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1537",
					Message: "`vgremove -f` skips the confirmation — a typo in the volume name destroys every filesystem on top of it.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — pvremove -ff pv",
			input: `pvremove -ff /devicenode`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1537",
					Message: "`pvremove -ff` skips the confirmation — a typo in the volume name destroys every filesystem on top of it.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1537")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1538(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — zpool list",
			input:    `zpool list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — zfs list",
			input:    `zfs list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — zfs destroy mydataset (no -r)",
			input:    `zfs destroy tank/data/old`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — zpool destroy -f tank",
			input: `zpool destroy -f tank`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1538",
					Message: "`zpool destroy -f` irrecoverably destroys the ZFS pool/dataset and every snapshot on it. Require explicit target confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — zfs destroy -rR tank/data",
			input: `zfs destroy -rR tank/data`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1538",
					Message: "`zfs destroy -rR` irrecoverably destroys the ZFS pool/dataset and every snapshot on it. Require explicit target confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1538")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1539(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — parted -s DISK print",
			input:    `parted -s $DISK print`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — parted DISK mklabel (interactive)",
			input:    `parted $DISK mklabel gpt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — parted -s DISK mklabel gpt",
			input: `parted -s $DISK mklabel gpt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1539",
					Message: "`parted -s <disk> mklabel` bypasses the confirmation prompt — a typo in the disk variable silently repartitions the wrong device.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — parted -s DISK rm 1",
			input: `parted -s $DISK rm 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1539",
					Message: "`parted -s <disk> rm` bypasses the confirmation prompt — a typo in the disk variable silently repartitions the wrong device.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1539")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1540(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — cryptsetup luksOpen",
			input:    `cryptsetup luksOpen $DEV mapname`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — cryptsetup luksRemoveKey",
			input:    `cryptsetup luksRemoveKey $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — cryptsetup erase $DEV",
			input: `cryptsetup erase $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1540",
					Message: "`cryptsetup erase` wipes the LUKS header — ciphertext becomes unrecoverable. Back up the header first, or use luksRemoveKey/luksKillSlot for single-slot rotation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — cryptsetup luksErase $DEV",
			input: `cryptsetup luksErase $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1540",
					Message: "`cryptsetup luksErase` wipes the LUKS header — ciphertext becomes unrecoverable. Back up the header first, or use luksRemoveKey/luksKillSlot for single-slot rotation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1540")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1541(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — apk add curl",
			input:    `apk add curl`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — apk add --allow-untrusted local.apk",
			input: `apk add --allow-untrusted ./local.apk`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1541",
					Message: "`apk --allow-untrusted` skips signature verification on the package — MITM-to-root on Alpine. Sign and place key in /etc/apk/keys/.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1541")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1542(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — snap install firefox",
			input:    `snap install firefox`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — snap install --dangerous local.snap",
			input: `snap install --dangerous ./local.snap`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1542",
					Message: "`snap install --dangerous` installs an assertion-unverified snap — any .snap on disk can register system services. Use --devmode or the store.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1542")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1543(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — go install pkg@v1.2.3",
			input:    `go install github.com/foo/bar@v1.2.3`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — cargo install --git url --rev sha",
			input:    `cargo install --git https://example.com/foo --rev abc123 foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — cargo install foo (crates.io pin via crate version)",
			input:    `cargo install foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — go install pkg@latest",
			input: `go install github.com/foo/bar@latest`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1543",
					Message: "`go install github.com/foo/bar@latest` is unpinned — HEAD-of-default can change between runs. Pin to a version tag or commit hash for reproducibility.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — cargo install --git url (no rev)",
			input: `cargo install --git https://example.com/foo foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1543",
					Message: "`cargo install --git (no --rev/--tag/--branch)` is unpinned — HEAD-of-default can change between runs. Pin to a version tag or commit hash for reproducibility.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1543")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1544(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dnf install curl",
			input:    `dnf install curl`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dnf copr enable user/repo",
			input: `dnf copr enable user/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1544",
					Message: "`dnf copr enable` pulls from a single-contributor repo — no distro security team. Pin the build, verify key fingerprint, mirror internally.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — add-apt-repository ppa:user/repo",
			input: `add-apt-repository ppa:user/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1544",
					Message: "`add-apt-repository` pulls from a single-contributor repo — no distro security team. Pin the build, verify key fingerprint, mirror internally.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1544")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1545(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker image prune",
			input:    `docker image prune`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker system prune (no -a / --volumes)",
			input:    `docker system prune -f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker system prune -af --volumes",
			input: `docker system prune -af --volumes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1545",
					Message: "`docker system prune` with `-a`/`--volumes` drops unused volumes — stopped stacks lose their databases. Scope the prune.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker volume prune -a",
			input: `docker volume prune -a`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1545",
					Message: "`docker volume prune` with `-a`/`--volumes` drops unused volumes — stopped stacks lose their databases. Scope the prune.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1545")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1546(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kubectl delete pod foo",
			input:    `kubectl delete pod foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kubectl delete pod foo --force (grace-period not 0)",
			input:    `kubectl delete pod foo --force`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubectl delete pod foo --force --grace-period=0",
			input: `kubectl delete pod foo --force --grace-period=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1546",
					Message: "`kubectl delete --force --grace-period=0` skips PreStop hooks and kubelet drain — corrupts StatefulSet state. Use standard delete.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — oc delete pod foo --force --grace-period=0",
			input: `oc delete pod foo --force --grace-period=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1546",
					Message: "`kubectl delete --force --grace-period=0` skips PreStop hooks and kubelet drain — corrupts StatefulSet state. Use standard delete.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1546")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1547(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kubectl apply -f manifests/",
			input:    `kubectl apply -f manifests/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kubectl apply --prune -l app=x -f manifests/",
			input:    `kubectl apply --prune -l app=x -f manifests/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubectl apply --prune --all -f m/",
			input: `kubectl apply --prune --all -f m/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1547",
					Message: "`kubectl apply --prune --all` deletes every matching resource not in the manifest — manifest typo wipes other teams' resources. Scope with a narrow `-l <selector>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — kubectl apply --prune -A -f m/",
			input: `kubectl apply --prune -A -f m/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1547",
					Message: "`kubectl apply --prune --all` deletes every matching resource not in the manifest — manifest typo wipes other teams' resources. Scope with a narrow `-l <selector>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1547")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1548(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — helm install foo chart",
			input:    `helm install foo bitnami/nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — helm install foo chart --disable-openapi-validation",
			input: `helm install foo bitnami/nginx --disable-openapi-validation`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1548",
					Message: "`helm --disable-openapi-validation` hides bad manifests until the controller crashes. Fix the schema deviation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — helm upgrade foo chart --disable-openapi-validation",
			input: `helm upgrade foo bitnami/nginx --disable-openapi-validation`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1548",
					Message: "`helm --disable-openapi-validation` hides bad manifests until the controller crashes. Fix the schema deviation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1548")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1549(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unzip to /tmp/stage",
			input:    `unzip foo.zip -d /tmp/stage`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — unzip -o to /opt/app",
			input:    `unzip -o foo.zip -d /opt/app`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — unzip -d /",
			input: `unzip foo.zip -d /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1549",
					Message: "`unzip -d /` extracts into a system path — any archive entry overwrites matching system file. Stage, inspect, copy.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — unzip -o file -d /boot",
			input: `unzip -o foo.zip -d /boot`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1549",
					Message: "`unzip -d /boot` extracts into a system path — any archive entry overwrites matching system file. Stage, inspect, copy.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1549")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1550(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — apt-mark unhold pkg",
			input:    `apt-mark unhold openssh-server`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — apt-mark showhold",
			input:    `apt-mark showhold`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — apt-mark hold pkg",
			input: `apt-mark hold openssh-server`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1550",
					Message: "`apt-mark hold` pins the package — blocks future CVE fixes. Document the reason and schedule an unhold review.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — apt-mark hold multiple pkgs",
			input: `apt-mark hold openssh-server libc6`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1550",
					Message: "`apt-mark hold` pins the package — blocks future CVE fixes. Document the reason and schedule an unhold review.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1550")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1551(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — helm install chart",
			input:    `helm install foo bitnami/nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — helm install chart --skip-crds",
			input: `helm install foo bitnami/nginx --skip-crds`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1551",
					Message: "`helm --skip-crds` installs .Release objects without their CRDs — custom resources fail validation. Install CRDs first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — helm upgrade chart --skip-crds",
			input: `helm upgrade foo bitnami/nginx --skip-crds`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1551",
					Message: "`helm --skip-crds` installs .Release objects without their CRDs — custom resources fail validation. Install CRDs first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1551")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1552(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — openssl genrsa 2048",
			input:    `openssl genrsa 2048`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — openssl dhparam 4096",
			input:    `openssl dhparam 4096`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — openssl x509 (not key-producing)",
			input:    `openssl x509 -in cert.pem -noout`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — openssl genrsa 1024",
			input: `openssl genrsa 1024`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1552",
					Message: "`openssl genrsa 1024` uses a weak key/param size — modern baselines require 2048+. Use 2048 or 3072/4096 for long-lived keys.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl dhparam 512",
			input: `openssl dhparam 512`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1552",
					Message: "`openssl dhparam 512` uses a weak key/param size — modern baselines require 2048+. Use 2048 or 3072/4096 for long-lived keys.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1552")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1553(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tr -d '[:space:]'",
			input:    `tr -d '[:space:]'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tr -s ' '",
			input:    `tr -s ' '`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tr '[:lower:]' '[:upper:]'",
			input: `tr '[:lower:]' '[:upper:]'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1553",
					Message: "`tr` for case conversion — use Zsh `${(U)var}` / `${(L)var}` to avoid the fork/exec and portability hazard.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tr a-z A-Z",
			input: `tr a-z A-Z`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1553",
					Message: "`tr` for case conversion — use Zsh `${(U)var}` / `${(L)var}` to avoid the fork/exec and portability hazard.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tr '[:upper:]' '[:lower:]'",
			input: `tr '[:upper:]' '[:lower:]'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1553",
					Message: "`tr` for case conversion — use Zsh `${(U)var}` / `${(L)var}` to avoid the fork/exec and portability hazard.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1553")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1554(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unzip without -o",
			input:    `unzip file.zip`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tar xf without --overwrite",
			input:    `tar xf foo.tar`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — unzip -o file.zip",
			input: `unzip -o file.zip`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1554",
					Message: "`unzip -o` overwrites existing files without prompting. Extract to a staging directory, diff, then move.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tar xf foo.tar --overwrite",
			input: `tar xf foo.tar --overwrite`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1554",
					Message: "`tar --overwrite` discards existing files during extract. Use a staging directory and diff before rolling forward.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1554")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1555(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chmod on own file",
			input:    `chmod 600 /tmp/myfile`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — chown on /etc/nginx",
			input:    `chown root:root /etc/nginx/nginx.conf`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chmod 666 /etc/shadow",
			input: `chmod 666 /etc/shadow`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1555",
					Message: "`chmod ... /etc/shadow` races the distro-managed tool — use passwd/chage/visudo or a config-management drop-in.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chown root:root /etc/sudoers",
			input: `chown root:root /etc/sudoers`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1555",
					Message: "`chown ... /etc/sudoers` races the distro-managed tool — use passwd/chage/visudo or a config-management drop-in.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chgrp shadow /etc/gshadow",
			input: `chgrp shadow /etc/gshadow`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1555",
					Message: "`chgrp ... /etc/gshadow` races the distro-managed tool — use passwd/chage/visudo or a config-management drop-in.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1555")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1556(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — openssl enc -aes-256-gcm",
			input:    `openssl enc -aes-256-gcm -in file -out enc`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — openssl enc -chacha20-poly1305",
			input:    `openssl enc -chacha20-poly1305 -in file -out enc`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — openssl enc -des",
			input: `openssl enc -des -in file -out enc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1556",
					Message: "`openssl enc -des` is a broken or deprecated cipher. Use `-aes-256-gcm` / `-chacha20-poly1305`, or `age` / `gpg` for files.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl enc -rc4",
			input: `openssl enc -rc4 -in file -out enc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1556",
					Message: "`openssl enc -rc4` is a broken or deprecated cipher. Use `-aes-256-gcm` / `-chacha20-poly1305`, or `age` / `gpg` for files.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl enc -des-ede3-cbc",
			input: `openssl enc -des-ede3-cbc -in file -out enc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1556",
					Message: "`openssl enc -des-ede3-cbc` is a broken or deprecated cipher. Use `-aes-256-gcm` / `-chacha20-poly1305`, or `age` / `gpg` for files.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1556")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1557(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kubeadm init",
			input:    `kubeadm init`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kubeadm reset (no -f)",
			input:    `kubeadm reset`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubeadm reset -f",
			input: `kubeadm reset -f`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1557",
					Message: "`kubeadm reset -f` skips the confirmation and wipes /etc/kubernetes / kubelet state. Drain and remove the node first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — kubeadm reset --force",
			input: `kubeadm reset --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1557",
					Message: "`kubeadm reset -f` skips the confirmation and wipes /etc/kubernetes / kubelet state. Drain and remove the node first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1557")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1558(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — usermod -aG audio alice",
			input:    `usermod -aG audio alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — usermod -aG wheel alice",
			input: `usermod -aG wheel alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1558",
					Message: "Adding user to `wheel` grants persistent admin-level access — use a scoped sudoers.d drop-in via configuration management.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — usermod -aG docker alice",
			input: `usermod -aG docker alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1558",
					Message: "Adding user to `docker` grants persistent admin-level access — use a scoped sudoers.d drop-in via configuration management.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — gpasswd -a alice sudo",
			input: `gpasswd -a alice sudo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1558",
					Message: "Adding user to `sudo` grants persistent admin-level access — use a scoped sudoers.d drop-in via configuration management.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — usermod -aG audio,wheel alice (mixed)",
			input: `usermod -aG audio,wheel alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1558",
					Message: "Adding user to `wheel` grants persistent admin-level access — use a scoped sudoers.d drop-in via configuration management.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1558")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1559(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh-copy-id user@host",
			input:    `ssh-copy-id user@host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ssh-copy-id -i ~/.ssh/id_ed25519.pub user@host",
			input:    `ssh-copy-id -i ~/.ssh/id_ed25519.pub user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ssh-copy-id -f user@host",
			input: `ssh-copy-id -f user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1559",
					Message: "`ssh-copy-id -f` pushes a long-term credential without host-key verification. Verify the fingerprint out of band first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh-copy-id -o StrictHostKeyChecking=no user@host",
			input: `ssh-copy-id -o StrictHostKeyChecking=no user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1559",
					Message: "`ssh-copy-id -o StrictHostKeyChecking=no` pushes a long-term credential without host-key verification. Verify the fingerprint out of band first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1559")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1560(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pip install foo",
			input:    `pip install foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — pip install --index-url https://x foo",
			input:    `pip install --index-url https://pypi.example.com foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pip install --trusted-host pypi.example.com foo",
			input: `pip install --trusted-host pypi.example.com foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1560",
					Message: "`pip --trusted-host` skips TLS verification and allows plain-HTTP for that index. Fix the CA trust and keep --index-url on https://.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — pip3 install foo --trusted-host pypi.org",
			input: `pip3 install foo --trusted-host pypi.org`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1560",
					Message: "`pip --trusted-host` skips TLS verification and allows plain-HTTP for that index. Fix the CA trust and keep --index-url on https://.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1560")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1561(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — systemctl isolate multi-user.target",
			input:    `systemctl isolate multi-user.target`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — systemctl start nginx.service",
			input:    `systemctl start nginx.service`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — systemctl isolate rescue.target",
			input: `systemctl isolate rescue.target`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1561",
					Message: "`systemctl isolate rescue.target` kills SSH and most services — console-only recovery. Do not run from a script.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — systemctl isolate emergency.target",
			input: `systemctl isolate emergency.target`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1561",
					Message: "`systemctl isolate emergency.target` kills SSH and most services — console-only recovery. Do not run from a script.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1561")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1562(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — env cmd",
			input:    `env cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — env -u TMPDIR cmd",
			input:    `env -u TMPDIR cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — env -u PATH cmd",
			input: `env -u PATH cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1562",
					Message: "`env -u PATH` clears a security-relevant variable mid-run. Use `env -i` to sanitise, or set the right value explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — env -u LD_PRELOAD cmd",
			input: `env -u LD_PRELOAD cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1562",
					Message: "`env -u LD_PRELOAD` clears a security-relevant variable mid-run. Use `env -i` to sanitise, or set the right value explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1562")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1563(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — swapon -a",
			input:    `swapon -a`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — swapoff specific file",
			input:    `swapoff /swapfile`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — swapoff -a",
			input: `swapoff -a`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1563",
					Message: "`swapoff -a` disables all swap devices — next memory-hungry process hits OOM. Document the trade-off if kubelet requires it.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1563")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1564(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — date (read)",
			input:    `date`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — timedatectl status",
			input:    `timedatectl status`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — date -s 2025-01-01",
			input: `date -s 2025-01-01`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1564",
					Message: "`date -s` sets the wall clock manually — breaks TLS certs, cron catch-up, and systemd timer math. Use timesyncd/chrony/ntpd.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — timedatectl set-time",
			input: `timedatectl set-time 2025-01-01`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1564",
					Message: "`timedatectl set-time` sets the wall clock manually — breaks TLS certs, cron catch-up, and systemd timer math. Use timesyncd/chrony/ntpd.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — hwclock -w",
			input: `hwclock -w`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1564",
					Message: "`hwclock -w` sets the wall clock manually — breaks TLS certs, cron catch-up, and systemd timer math. Use timesyncd/chrony/ntpd.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1564")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1565(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — command -v cmd",
			input:    `command -v git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — whereis git",
			input: `whereis git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1565",
					Message: "`whereis` is index-based and stale-prone. Use `command -v <cmd>` for runtime existence checks.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — locate foo",
			input: `locate foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1565",
					Message: "`locate` is index-based and stale-prone. Use `command -v <cmd>` for runtime existence checks.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1565")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1566(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — gem install rails",
			input:    `gem install rails`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — gem install -P HighSecurity",
			input:    `gem install -P HighSecurity rails`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — gem install -P NoSecurity",
			input: `gem install -P NoSecurity rails`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1566",
					Message: "`gem -P NoSecurity` skips signature verification — MITM or account compromise becomes RCE at install. Use HighSecurity.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — gem install --trust-policy LowSecurity",
			input: `gem install --trust-policy LowSecurity rails`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1566",
					Message: "`gem -P LowSecurity` skips signature verification — MITM or account compromise becomes RCE at install. Use HighSecurity.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1566")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1567(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — python -m http.server --bind 127.0.0.1",
			input:    `python -m http.server --bind 127.0.0.1 8080`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — python -m http.server -b 127.0.0.1",
			input:    `python -m http.server -b 127.0.0.1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — python -m venv myenv",
			input:    `python -m venv myenv`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — python -m http.server",
			input: `python -m http.server`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1567",
					Message: "`python -m http.server` without `--bind` defaults to 0.0.0.0 — exposes the cwd to every network the host sees. Add `--bind 127.0.0.1`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — python3 -m http.server 8080",
			input: `python3 -m http.server 8080`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1567",
					Message: "`python -m http.server` without `--bind` defaults to 0.0.0.0 — exposes the cwd to every network the host sees. Add `--bind 127.0.0.1`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — python2 -m SimpleHTTPServer",
			input: `python2 -m SimpleHTTPServer`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1567",
					Message: "`python -m http.server` without `--bind` defaults to 0.0.0.0 — exposes the cwd to every network the host sees. Add `--bind 127.0.0.1`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1567")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1568(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — useradd -u 1000 alice",
			input:    `useradd -u 1000 alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — useradd -o -u 1000 alice",
			input: `useradd -o -u 1000 alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1568",
					Message: "`useradd -o` assigns a non-unique UID — the two accounts share kernel identity, indistinguishable in audit. Use a fresh UID.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — usermod -o -u 500 bob",
			input: `usermod -o -u 500 bob`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1568",
					Message: "`usermod -o` assigns a non-unique UID — the two accounts share kernel identity, indistinguishable in audit. Use a fresh UID.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1568")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1569(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — nvme list",
			input:    `nvme list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — nvme id-ctrl $DISK",
			input:    `nvme id-ctrl $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — nvme format -s1 $DISK",
			input: `nvme format -s1 $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1569",
					Message: "`nvme format -s1` unrecoverably erases the namespace in seconds. Do not run from automation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — nvme format -s2 $DISK",
			input: `nvme format -s2 $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1569",
					Message: "`nvme format -s2` unrecoverably erases the namespace in seconds. Do not run from automation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — nvme sanitize -a 4 $DISK",
			input: `nvme sanitize -a 4 $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1569",
					Message: "`nvme sanitize -a` unrecoverably erases the namespace in seconds. Do not run from automation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1569")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1570(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — smbclient -U user //server/share",
			input:    `smbclient -U user //server/share`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — smbclient -N //server/share",
			input: `smbclient -N //server/share`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1570",
					Message: "`smbclient -N` is anonymous SMB access — any host on-net can read the share. Use credentials=<file> 0600 or -k.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — mount.cifs -N //server/share /mnt",
			input: `mount.cifs -N //server/share /mnt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1570",
					Message: "`mount.cifs -N` is anonymous SMB access — any host on-net can read the share. Use credentials=<file> 0600 or -k.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1570")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1571(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chronyc makestep",
			input:    `chronyc makestep`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ntpdate pool.ntp.org",
			input: `ntpdate pool.ntp.org`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1571",
					Message: "`ntpdate` is deprecated and races any running chrony/timesyncd. Use `chronyc makestep` or `systemctl restart systemd-timesyncd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sntp -sS pool.ntp.org",
			input: `sntp -sS pool.ntp.org`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1571",
					Message: "`sntp` is deprecated and races any running chrony/timesyncd. Use `chronyc makestep` or `systemctl restart systemd-timesyncd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1571")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1572(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run -e LOG_LEVEL=info",
			input:    `docker run -e LOG_LEVEL=info alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker run -e PASSWORD (no value, inherits)",
			input:    `docker run -e PASSWORD alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker run --env-file secrets alpine",
			input:    `docker run --env-file secrets alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run -e PASSWORD=hunter2",
			input: `docker run -e PASSWORD=hunter2 alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1572",
					Message: "`-e PASSWORD=<value>` writes the secret into `docker inspect` and `/proc/1/environ`. Use `--env-file` 0600 or `--secret`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman run -e API_KEY=abc123",
			input: `podman run -e API_KEY=abc123 alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1572",
					Message: "`-e API_KEY=<value>` writes the secret into `docker inspect` and `/proc/1/environ`. Use `--env-file` 0600 or `--secret`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1572")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1573(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chattr =i (set exclusive)",
			input:    `chattr =i /etc/shadow`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chattr -i /etc/shadow",
			input: `chattr -i /etc/shadow`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1573",
					Message: "`chattr -i` removes the tamper-evident attribute. If this is a one-shot upgrade, re-set the attribute at the end of the block.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chattr -a /var/log/auth.log",
			input: `chattr -a /var/log/auth.log`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1573",
					Message: "`chattr -a` removes the tamper-evident attribute. If this is a one-shot upgrade, re-set the attribute at the end of the block.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1573")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1574(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — git config credential.helper libsecret",
			input:    `git config credential.helper libsecret`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — git config credential.helper 'cache --timeout=3600'",
			input:    `git config credential.helper 'cache --timeout=3600'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — git config credential.helper store",
			input: `git config credential.helper store`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1574",
					Message: "`git credential.helper store` saves credentials in plaintext — backups leak the token. Use platform helper (manager-core / libsecret) or `cache --timeout=<sec>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — git config --global credential.helper store",
			input: `git config --global credential.helper store`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1574",
					Message: "`git credential.helper store` saves credentials in plaintext — backups leak the token. Use platform helper (manager-core / libsecret) or `cache --timeout=<sec>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1574")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1575(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — aws configure set region us-east-1",
			input:    `aws configure set region us-east-1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — aws configure set aws_secret_access_key VALUE",
			input: `aws configure set aws_secret_access_key AKIAEXAMPLEKEYXYZ`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1575",
					Message: "`aws configure set aws_secret_access_key …` puts the secret in ps / history. Use IAM-role auth or import from stdin / 0600 file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — aws configure set aws_session_token",
			input: `aws configure set aws_session_token FwoGZXIvYXdzEXAMPLE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1575",
					Message: "`aws configure set aws_session_token …` puts the secret in ps / history. Use IAM-role auth or import from stdin / 0600 file.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1575")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1576(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — terraform apply",
			input:    `terraform apply`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — terraform apply -target=module.foo",
			input: `terraform apply -target=module.foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1576",
					Message: "`terraform -target=module.foo` bypasses dependency order — documented as incident response tool only. Re-run without -target or split root modules.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — terraform apply -target module.foo",
			input: `terraform apply -target module.foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1576",
					Message: "`terraform -target module.foo` bypasses dependency order — documented as incident response tool only. Re-run without -target or split root modules.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tofu destroy -target=aws_instance.web",
			input: `tofu destroy -target=aws_instance.web`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1576",
					Message: "`terraform -target=aws_instance.web` bypasses dependency order — documented as incident response tool only. Re-run without -target or split root modules.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1576")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1577(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dig A example.com",
			input:    `dig A example.com`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — dig MX example.com",
			input:    `dig MX example.com`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dig ANY example.com",
			input: `dig ANY example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1577",
					Message: "`dig ... ANY` is RFC 8482-deprecated — filtered by recursors. Query specific types (A / MX / NS / …) and combine.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — dig example.com ANY",
			input: `dig example.com ANY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1577",
					Message: "`dig ... ANY` is RFC 8482-deprecated — filtered by recursors. Query specific types (A / MX / NS / …) and combine.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1577")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1578(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh-keygen -t ed25519 -f key",
			input:    `ssh-keygen -t ed25519 -f key`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ssh-keygen -t rsa -b 4096",
			input:    `ssh-keygen -t rsa -b 4096`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ssh-keygen -t rsa -b 1024",
			input: `ssh-keygen -t rsa -b 1024`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1578",
					Message: "`ssh-keygen -b 1024` — RSA below 2048 bits is rejected by modern OpenSSH. Use `-t ed25519` or `-b 4096`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh-keygen -t dsa",
			input: `ssh-keygen -t dsa -f key`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1578",
					Message: "`ssh-keygen -t dsa` — DSA removed from OpenSSH 9.8. Use `-t ed25519`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1578")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1579(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — curl https://host",
			input:    `curl https://host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — curl --retry-all-errors --max-time 30 URL",
			input:    `curl https://host --retry-all-errors --max-time 30`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — curl --retry-all-errors -m 30 URL",
			input:    `curl https://host --retry-all-errors -m 30`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — curl URL --retry-all-errors (no max-time)",
			input: `curl https://host --retry-all-errors`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1579",
					Message: "`curl --retry-all-errors` with no `--max-time` hammers the upstream on failure. Pair with `-m <seconds>` or use `--retry-connrefused`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1579")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1580(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — go build",
			input:    `go build -o app ./cmd/app`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — go build -ldflags with version",
			input:    `go build -ldflags "-X main.Version=1.2.3" ./cmd/app`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — go build -ldflags with PASSWORD",
			input: `go build -ldflags "-X main.PASSWORD=hunter2" ./cmd/app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1580",
					Message: "`go build -ldflags` injecting a secret bakes it into the binary. Read from os.Getenv / mounted secret file at runtime.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — go build -ldflags with API_KEY",
			input: `go build -ldflags "-X main.API_KEY=xyz" ./cmd/app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1580",
					Message: "`go build -ldflags` injecting a secret bakes it into the binary. Read from os.Getenv / mounted secret file at runtime.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1580")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1581(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh user@host",
			input:    `ssh user@host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ssh -o PubkeyAuthentication=yes",
			input:    `ssh -o PubkeyAuthentication=yes user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ssh -o PubkeyAuthentication=no",
			input: `ssh -o PubkeyAuthentication=no user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1581",
					Message: "`ssh -o PubkeyAuthentication=no` forces password auth — weaker than key auth. Let the default preference pick.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh -o PasswordAuthentication=yes",
			input: `ssh -o PasswordAuthentication=yes user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1581",
					Message: "`ssh -o PasswordAuthentication=yes` forces password auth — weaker than key auth. Let the default preference pick.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh -o PreferredAuthentications=password",
			input: `ssh -o PreferredAuthentications=password user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1581",
					Message: "`ssh -o PreferredAuthentications=password` forces password auth — weaker than key auth. Let the default preference pick.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1581")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1582(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — bash script.sh",
			input:    `bash script.sh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — bash -x script.sh",
			input: `bash -x script.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1582",
					Message: "`bash -x` traces every expanded command — CI logs leak secrets verbatim. Scope with `set -x; …; set +x`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sh -x script.sh",
			input: `sh -x script.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1582",
					Message: "`sh -x` traces every expanded command — CI logs leak secrets verbatim. Scope with `set -x; …; set +x`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — zsh -xv script.zsh",
			input: `zsh -xv script.zsh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1582",
					Message: "`zsh -xv` traces every expanded command — CI logs leak secrets verbatim. Scope with `set -x; …; set +x`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1582")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1583(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -delete",
			input:    `find /tmp -name '*.log'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — find -maxdepth 2 -delete",
			input:    `find /tmp -maxdepth 2 -name '*.log' -delete`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — find -xdev -delete",
			input:    `find /var -xdev -name '*.log' -delete`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find /tmp -name '*.log' -delete",
			input: `find /tmp -name '*.log' -delete`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1583",
					Message: "`find -delete` without `-maxdepth` / `-xdev` / `-prune` walks the whole tree. Scope the depth (e.g. `-maxdepth 2`) and dry-run first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1583")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1584(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sudo -u root cmd",
			input:    `sudo -u root cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — sudo env VAR=1 cmd",
			input:    `sudo env VAR=1 cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sudo -E cmd",
			input: `sudo -E cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1584",
					Message: "`sudo -E` carries the caller's PATH / LD_PRELOAD / … into the privileged process. Use `env_keep` in sudoers or explicit `sudo env VAR=… cmd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sudo -E -u svc cmd",
			input: `sudo -E -u svc cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1584",
					Message: "`sudo -E` carries the caller's PATH / LD_PRELOAD / … into the privileged process. Use `env_keep` in sudoers or explicit `sudo env VAR=… cmd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1584")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1585(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ufw allow from 10.0.0.0/8 to any port 22",
			input:    `ufw allow from 10.0.0.0/8 to any port 22`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ufw status",
			input:    `ufw status`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ufw allow from any to any port 22",
			input: `ufw allow from any to any port 22`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1585",
					Message: "`ufw allow from any …` opens the port to the whole internet. Scope to a specific source CIDR.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1585")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1586(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — systemctl enable sshd",
			input:    `systemctl enable sshd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chkconfig sshd on",
			input: `chkconfig sshd on`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1586",
					Message: "`chkconfig` is a SysV-init relic. Use `systemctl enable|disable <unit>` directly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — update-rc.d sshd defaults",
			input: `update-rc.d sshd defaults`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1586",
					Message: "`update-rc.d` is a SysV-init relic. Use `systemctl enable|disable <unit>` directly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — insserv sshd",
			input: `insserv sshd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1586",
					Message: "`insserv` is a SysV-init relic. Use `systemctl enable|disable <unit>` directly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1586")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1587(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — modprobe nvme",
			input:    `modprobe nvme`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — lsmod",
			input:    `lsmod`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — modprobe -r nvme",
			input: `modprobe -r nvme`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1587",
					Message: "`modprobe -r` unloads an in-use module — the backing subsystem goes offline. Use `systemctl stop` if you meant to stop a service.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — rmmod nvidia",
			input: `rmmod nvidia`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1587",
					Message: "`rmmod` unloads a kernel module — the backing subsystem goes offline. Use `systemctl stop` if you meant to stop a service.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1587")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1588(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — nsenter without target 1",
			input:    `nsenter -t 4242 -m sh`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — nsenter on arbitrary pid",
			input:    `nsenter -t 8123 -m sh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — nsenter -t 1 -m -u -i -n -p sh",
			input: `nsenter -t 1 -m -u -i -n -p sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1588",
					Message: "`nsenter --target 1` joins the host init namespaces — classic container-escape primitive. Use `docker exec` / `kubectl exec` for legitimate debugging.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — nsenter -t1 -m sh",
			input: `nsenter -t1 -m sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1588",
					Message: "`nsenter --target 1` joins the host init namespaces — classic container-escape primitive. Use `docker exec` / `kubectl exec` for legitimate debugging.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1588")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1589(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — trap cleanup EXIT",
			input:    `trap cleanup EXIT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — trap with safe dump",
			input:    `trap 'echo failed' ERR`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — trap 'set -x' ERR",
			input: `trap 'set -x' ERR`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1589",
					Message: "`trap 'set -x' ERR` enables shell trace from a trap — expansions leak env vars (tokens, passwords) to stderr / CI logs. Use a scoped `set -x ... set +x`, not a trap.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — trap 'set -o xtrace' EXIT",
			input: `trap 'set -o xtrace' EXIT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1589",
					Message: "`trap 'set -x' EXIT` enables shell trace from a trap — expansions leak env vars (tokens, passwords) to stderr / CI logs. Use a scoped `set -x ... set +x`, not a trap.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — trap 'set -x' RETURN",
			input: `trap 'set -x' RETURN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1589",
					Message: "`trap 'set -x' RETURN` enables shell trace from a trap — expansions leak env vars (tokens, passwords) to stderr / CI logs. Use a scoped `set -x ... set +x`, not a trap.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1589")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1590(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sshpass -e (env var)",
			input:    `sshpass -e ssh user@host cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — sshpass -f FILE",
			input:    `sshpass -f /run/secrets/pw ssh user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sshpass -p 'secret' ssh ...",
			input: `sshpass -p 'secret' ssh user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1590",
					Message: "`sshpass -p` places the password in argv — visible in `ps` / `/proc/<pid>/cmdline`. Switch to key-based auth, or at least use `sshpass -f FILE` / `sshpass -e`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sshpass -psecret ssh ...",
			input: `sshpass -psecret ssh user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1590",
					Message: "`sshpass -p` places the password in argv — visible in `ps` / `/proc/<pid>/cmdline`. Switch to key-based auth, or at least use `sshpass -f FILE` / `sshpass -e`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1590")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1591(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — printf with specific format",
			input:    `printf '%-20s %d\n' "${pairs[@]}"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — printf on scalar",
			input:    `printf '%s\n' "$msg"`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — printf '%s\n' "${array[@]}"`,
			input: `printf '%s\n' "${array[@]}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1591",
					Message: "`printf '%s\\n' \"${array[@]}\"` — use Zsh `print -l -r -- \"${array[@]}\"` or `${(F)array}` for newline-joined output.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — printf '%s' "${a[*]}"`,
			input: `printf '%s' "${a[*]}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1591",
					Message: "`printf '%s\\n' \"${array[@]}\"` — use Zsh `print -l -r -- \"${array[@]}\"` or `${(F)array}` for newline-joined output.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1591")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1592(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — faillock status",
			input:    `faillock -u bob`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — pam_tally2 status",
			input:    `pam_tally2 -u bob`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — faillock -u bob --reset",
			input: `faillock -u bob --reset`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1592",
					Message: "`faillock --reset` clears the PAM failed-auth counter — masks ongoing brute force. Log the prior count and alert before resetting.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — pam_tally2 -r -u bob",
			input: `pam_tally2 -r -u bob`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1592",
					Message: "`pam_tally2 -r` clears the PAM failed-auth counter — masks ongoing brute force. Log the prior count and alert before resetting.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1592")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1593(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — different command",
			input:    `lsblk $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — blkdiscard $DISK",
			input: `blkdiscard $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1593",
					Message: "`blkdiscard` issues TRIM/DISCARD across the full device — data is unrecoverable once the controller acknowledges. Require operator confirmation before running.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — blkdiscard -z $DISK",
			input: `blkdiscard -z $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1593",
					Message: "`blkdiscard` issues TRIM/DISCARD across the full device — data is unrecoverable once the controller acknowledges. Require operator confirmation before running.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1593")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1594(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — plain docker run",
			input:    `docker run --rm alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — podman run with different security-opt",
			input:    `podman run --security-opt=no-new-privileges alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --security-opt=systempaths=unconfined",
			input: `docker run --security-opt=systempaths=unconfined alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1594",
					Message: "`docker run --security-opt=systempaths=unconfined` unhides `/proc/sys`, `/proc/sysrq-trigger`, and other kernel knobs. A compromise in the container can then panic or re-tune the host.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman create systempaths=unconfined",
			input: `podman create --security-opt=systempaths=unconfined alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1594",
					Message: "`podman run --security-opt=systempaths=unconfined` unhides `/proc/sys`, `/proc/sysrq-trigger`, and other kernel knobs. A compromise in the container can then panic or re-tune the host.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1594")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1595(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — setfacl -m u:alice:r file",
			input:    `setfacl -m u:alice:r /tmp/report`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — setfacl -x",
			input:    `setfacl -x u:alice /tmp/report`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — setfacl -m u:nobody:rwx file",
			input: `setfacl -m u:nobody:rwx /etc/app.conf`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1595",
					Message: "`setfacl -m u:nobody:rwx` grants perms via ACL, bypassing `chmod` / `stat -c %a` checks. Prefer chmod for world perms, and for specific users name the real account with minimum perms.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — setfacl -m o::rwx file",
			input: `setfacl -m o::rwx /etc/app.conf`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1595",
					Message: "`setfacl -m o::rwx` grants perms via ACL, bypassing `chmod` / `stat -c %a` checks. Prefer chmod for world perms, and for specific users name the real account with minimum perms.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1595")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1596(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — emulate -L sh",
			input:    `emulate -L sh`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — emulate -LR bash",
			input:    `emulate -LR bash`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — emulate zsh (reset to zsh)",
			input:    `emulate zsh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — emulate sh",
			input: `emulate sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1596",
					Message: "`emulate sh` without `-L` flips the options for the whole shell. Use `emulate -L sh` inside a function, or rename the script to `.sh` if Zsh features are not needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — emulate -R bash",
			input: `emulate -R bash`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1596",
					Message: "`emulate bash` without `-L` flips the options for the whole shell. Use `emulate -L bash` inside a function, or rename the script to `.sh` if Zsh features are not needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1596")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1597(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — systemd-run as non-root user",
			input:    `systemd-run -p User=www-data /usr/bin/cleanup`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — systemd-run without user property",
			input:    `systemd-run /usr/bin/cleanup`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — systemd-run -p User=root",
			input: `systemd-run -p User=root sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1597",
					Message: "`systemd-run -p User=root` runs arbitrary commands as root via systemd — bypasses the `sudo` audit path. Prefer explicit `sudo` or a fixed systemd unit.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — systemd-run -p User=0",
			input: `systemd-run -p User=0 sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1597",
					Message: "`systemd-run -p User=0` runs arbitrary commands as root via systemd — bypasses the `sudo` audit path. Prefer explicit `sudo` or a fixed systemd unit.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1597")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1598(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chmod 600 on /dev/kvm",
			input:    `chmod 600 /dev/kvm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — chmod 666 on /dev/null (safe)",
			input:    `chmod 666 /dev/null`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — chmod 666 on regular file (not /dev/)",
			input:    `chmod 666 /tmp/log`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chmod 666 /dev/kvm",
			input: `chmod 666 /dev/kvm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1598",
					Message: "`chmod 666 /dev/kvm` makes a sensitive device node world-writable — direct kernel access for every local user. Keep restrictive perms (600 / 660) and grant access via udev rules.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chmod 777 /dev/mem",
			input: `chmod 777 /dev/mem`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1598",
					Message: "`chmod 777 /dev/mem` makes a sensitive device node world-writable — direct kernel access for every local user. Keep restrictive perms (600 / 660) and grant access via udev rules.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1598")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1599(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — plain ldconfig",
			input:    `ldconfig`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ldconfig -f /etc/ld.so.conf.d/custom.conf",
			input:    `ldconfig -f /etc/ld.so.conf.d/custom.conf`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ldconfig -f $LDCONF (variable)",
			input:    `ldconfig -f $LDCONF`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ldconfig -f /tmp/fake.conf",
			input: `ldconfig -f /tmp/fake.conf`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1599",
					Message: "`ldconfig -f /tmp/fake.conf` uses a config outside `/etc/`. If the file is attacker-writable, every binary on the host loads the attacker's library. Keep config under `/etc/ld.so.conf.d/`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ldconfig -f /var/tmp/x.conf",
			input: `ldconfig -f /var/tmp/x.conf`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1599",
					Message: "`ldconfig -f /var/tmp/x.conf` uses a config outside `/etc/`. If the file is attacker-writable, every binary on the host loads the attacker's library. Keep config under `/etc/ld.so.conf.d/`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1599")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
