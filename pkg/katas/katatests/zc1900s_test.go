// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1900(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `curl -L https://api/resource`",
			input:    `curl -L https://api/resource`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl -u user:pass https://api/resource` (no location)",
			input:    `curl -u user:pass https://api/resource`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `curl --location-trusted https://api` (leading)",
			input: `curl --location-trusted https://api`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1900",
					Message: "`curl --location-trusted` replays `Authorization`, cookies, and `-u user:pass` on every redirect — a 302 to attacker-controlled host leaks the token. Drop the flag; verify final hostname before sending secrets.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:  "invalid — `curl -u user:pass --location-trusted https://api` (trailing)",
			input: `curl -u user:pass --location-trusted https://api`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1900",
					Message: "`curl --location-trusted` replays `Authorization`, cookies, and `-u user:pass` on every redirect — a 302 to attacker-controlled host leaks the token. Drop the flag; verify final hostname before sending secrets.",
					Line:    1,
					Column:  20,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1900")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1901(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt POSIX_BUILTINS` (explicit default)",
			input:    `unsetopt POSIX_BUILTINS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt POSIX_BUILTINS`",
			input: `setopt POSIX_BUILTINS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1901",
					Message: "`setopt POSIX_BUILTINS` switches Zsh to POSIX special-builtin rules — assignments before `export`/`readonly`/`eval` stop being local, silently leaking state. Scope any POSIX block with `emulate -LR sh` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_POSIX_BUILTINS`",
			input: `unsetopt NO_POSIX_BUILTINS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1901",
					Message: "`unsetopt NO_POSIX_BUILTINS` switches Zsh to POSIX special-builtin rules — assignments before `export`/`readonly`/`eval` stop being local, silently leaking state. Scope any POSIX block with `emulate -LR sh` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1901")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1902(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ln -s /opt/app/current /opt/app/live` (app release symlink)",
			input:    `ln -s /opt/app/current /opt/app/live`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ln -s /dev/null /tmp/scratch` (non-sensitive target)",
			input:    `ln -s /dev/null /tmp/scratch`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ln -sf /dev/null /var/log/auth.log`",
			input: `ln -sf /dev/null /var/log/auth.log`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1902",
					Message: "`ln -s /dev/null /var/log/auth.log` redirects every write to the bit-bucket — audit / history entries vanish silently. If the log must stop, disable the writer or rotate with `logrotate`, never symlink to `/dev/null`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ln -s /dev/null $HOME/.bash_history`",
			input: `ln -s /dev/null $HOME/.bash_history`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1902",
					Message: "`ln -s /dev/null $HOME/.bash_history` redirects every write to the bit-bucket — audit / history entries vanish silently. If the log must stop, disable the writer or rotate with `logrotate`, never symlink to `/dev/null`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1902")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1903(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `tee /var/log/app.log`",
			input:    `tee /var/log/app.log`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `tee /etc/nginx/conf.d/site.conf`",
			input:    `tee /etc/nginx/conf.d/site.conf`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tee /etc/sudoers`",
			input: `tee /etc/sudoers`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1903",
					Message: "`tee /etc/sudoers` writes a sudoers rule without `visudo -c` validation — a syntax error locks every future `sudo` invocation. Write to a temp file, run `visudo -cf`, then `install -m 0440` into place.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `tee -a /etc/sudoers.d/app`",
			input: `tee -a /etc/sudoers.d/app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1903",
					Message: "`tee /etc/sudoers.d/app` writes a sudoers rule without `visudo -c` validation — a syntax error locks every future `sudo` invocation. Write to a temp file, run `visudo -cf`, then `install -m 0440` into place.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1903")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1904(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt KSH_GLOB` (explicit default)",
			input:    `unsetopt KSH_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_GLOB` (unrelated)",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt KSH_GLOB`",
			input: `setopt KSH_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1904",
					Message: "`setopt KSH_GLOB` reinterprets `*(...)` as a ksh-style operator — every Zsh glob qualifier (`*(N)`, `*(D)`, `*(.)`) silently stops working. Prefer `setopt EXTENDED_GLOB`, or scope inside a function with `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_KSH_GLOB`",
			input: `unsetopt NO_KSH_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1904",
					Message: "`unsetopt NO_KSH_GLOB` reinterprets `*(...)` as a ksh-style operator — every Zsh glob qualifier (`*(N)`, `*(D)`, `*(.)`) silently stops working. Prefer `setopt EXTENDED_GLOB`, or scope inside a function with `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1904")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1905(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ssh -L 8080:target:80 host` (no -g)",
			input:    `ssh -L 8080:target:80 host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ssh -g host` (no forward)",
			input:    `ssh -g host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ssh -g -L 8080:target:80 host`",
			input: `ssh -g -L 8080:target:80 host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1905",
					Message: "`ssh -g` with `-L`/`-D` binds the forward on `0.0.0.0` — anyone on the same LAN segment can ride the tunnel. Drop `-g` or pin `bind_address:port` in the forward spec.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ssh -gD 1080 host`",
			input: `ssh -gD 1080 host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1905",
					Message: "`ssh -g` with `-L`/`-D` binds the forward on `0.0.0.0` — anyone on the same LAN segment can ride the tunnel. Drop `-g` or pin `bind_address:port` in the forward spec.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1905")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1906(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt POSIX_CD` (explicit default)",
			input:    `unsetopt POSIX_CD`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt AUTO_PUSHD` (unrelated)",
			input:    `setopt AUTO_PUSHD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt POSIX_CD`",
			input: `setopt POSIX_CD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1906",
					Message: "`setopt POSIX_CD` changes when `cd`/`pushd` read `CDPATH` — scripts that relied on Zsh's default silently enter different directories. Keep it off; wrap POSIX-specific code with `emulate -LR sh`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_POSIX_CD`",
			input: `unsetopt NO_POSIX_CD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1906",
					Message: "`unsetopt NO_POSIX_CD` changes when `cd`/`pushd` read `CDPATH` — scripts that relied on Zsh's default silently enter different directories. Keep it off; wrap POSIX-specific code with `emulate -LR sh`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1906")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1907(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sysctl -w fs.protected_symlinks=1` (re-enable)",
			input:    `sysctl -w fs.protected_symlinks=1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sysctl -w vm.swappiness=10` (unrelated)",
			input:    `sysctl -w vm.swappiness=10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `sysctl -w fs.protected_hardlinks=0`",
			input: `sysctl -w fs.protected_hardlinks=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1907",
					Message: "`sysctl -w fs.protected_hardlinks=0` re-enables hardlink following — classic /tmp-race escalation vector. Keep the default; scope any exception in a dedicated namespace.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `sysctl -w fs.suid_dumpable=2`",
			input: `sysctl -w fs.suid_dumpable=2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1907",
					Message: "`sysctl -w fs.suid_dumpable=2` re-enables SUID core-dump exposure (2 = root-readable) — classic /tmp-race escalation vector. Keep the default; scope any exception in a dedicated namespace.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1907")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1908(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt MAGIC_EQUAL_SUBST` (explicit default)",
			input:    `unsetopt MAGIC_EQUAL_SUBST`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_GLOB` (unrelated)",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt MAGIC_EQUAL_SUBST`",
			input: `setopt MAGIC_EQUAL_SUBST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1908",
					Message: "`setopt MAGIC_EQUAL_SUBST` gives every `key=value` argument tilde/parameter expansion on the RHS — literal CLI args like `rsync host:dst=~/backup` silently change. Keep it off; quote the assignment if expansion is really wanted.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_MAGIC_EQUAL_SUBST`",
			input: `unsetopt NO_MAGIC_EQUAL_SUBST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1908",
					Message: "`unsetopt NO_MAGIC_EQUAL_SUBST` gives every `key=value` argument tilde/parameter expansion on the RHS — literal CLI args like `rsync host:dst=~/backup` silently change. Keep it off; quote the assignment if expansion is really wanted.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1908")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1909(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kexec --help`",
			input:    `kexec --help`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kexec -h`",
			input:    `kexec -h`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kexec -l $KERN`",
			input: `kexec -l $KERN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1909",
					Message: "`kexec -l` stages or jumps to a kernel without firmware / bootloader verification — Secure Boot never checks the signature. Gate behind `sudo` + audit and prefer `systemctl kexec` or a real reboot.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `kexec -e` (execute)",
			input: `kexec -e`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1909",
					Message: "`kexec -e` stages or jumps to a kernel without firmware / bootloader verification — Secure Boot never checks the signature. Gate behind `sudo` + audit and prefer `systemctl kexec` or a real reboot.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1909")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1910(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt GLOB_STAR_SHORT` (explicit default)",
			input:    `unsetopt GLOB_STAR_SHORT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_GLOB` (unrelated)",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt GLOB_STAR_SHORT`",
			input: `setopt GLOB_STAR_SHORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1910",
					Message: "`setopt GLOB_STAR_SHORT` turns bare `**` into `**/*` — `rm **` now wipes the tree. Keep the option off and spell `**/*` when recursion is actually wanted.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_GLOB_STAR_SHORT`",
			input: `unsetopt NO_GLOB_STAR_SHORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1910",
					Message: "`unsetopt NO_GLOB_STAR_SHORT` turns bare `**` into `**/*` — `rm **` now wipes the tree. Keep the option off and spell `**/*` when recursion is actually wanted.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1910")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1911(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `umount /mnt/scratch`",
			input:    `umount /mnt/scratch`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `umount -f /mnt/stuck` (force, not lazy)",
			input:    `umount -f /mnt/stuck`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `umount -l /mnt/scratch`",
			input: `umount -l /mnt/scratch`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1911",
					Message: "`umount -l` detaches the mount but leaves any open fd pointing at a ghost filesystem — writers keep writing, re-mounts stack invisibly. Stop the fd holder first (`lsof`/`fuser`), then do a normal `umount`.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:  "invalid — `umount --lazy /mnt/scratch`",
			input: `umount --lazy /mnt/scratch`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1911",
					Message: "`umount --lazy` detaches the mount but leaves any open fd pointing at a ghost filesystem — writers keep writing, re-mounts stack invisibly. Stop the fd holder first (`lsof`/`fuser`), then do a normal `umount`.",
					Line:    1,
					Column:  9,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1911")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1912(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `dhclient eth0` (renew)",
			input:    `dhclient eth0`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `dhcpcd --rebind eth0`",
			input:    `dhcpcd --rebind eth0`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `dhclient -r eth0`",
			input: `dhclient -r eth0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1912",
					Message: "`dhclient -r` drops the DHCP lease — SSH session cuts, VPC reachability stalls. Pair with a re-acquire (`dhclient -1`/`nmcli device reapply`), or schedule via `systemd-run --on-active=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `dhcpcd -k eth0`",
			input: `dhcpcd -k eth0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1912",
					Message: "`dhcpcd -k` drops the DHCP lease — SSH session cuts, VPC reachability stalls. Pair with a re-acquire (`dhclient -1`/`nmcli device reapply`), or schedule via `systemd-run --on-active=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1912")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1913(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt ALIAS_FUNC_DEF` (explicit default)",
			input:    `unsetopt ALIAS_FUNC_DEF`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt AUTO_CD` (unrelated)",
			input:    `setopt AUTO_CD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt ALIAS_FUNC_DEF`",
			input: `setopt ALIAS_FUNC_DEF`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1913",
					Message: "`setopt ALIAS_FUNC_DEF` lets a function silently shadow an alias — one sourced rc file replaces your function with the alias, no error surfaces. Keep it off; quote the name if the override is intentional.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_ALIAS_FUNC_DEF`",
			input: `unsetopt NO_ALIAS_FUNC_DEF`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1913",
					Message: "`unsetopt NO_ALIAS_FUNC_DEF` lets a function silently shadow an alias — one sourced rc file replaces your function with the alias, no error surfaces. Keep it off; quote the name if the override is intentional.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1913")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1914(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `curl https://api.example/resource`",
			input:    `curl https://api.example/resource`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl --resolve api:443:10.0.0.1 https://api/`",
			input:    `curl --resolve api:443:10.0.0.1 https://api/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `curl https://api/ --doh-url=https://doh/dns-query`",
			input: `curl https://api/ --doh-url=https://doh/dns-query`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1914",
					Message: "`curl --doh-url` bypasses the host's resolver chain — `/etc/hosts`, `systemd-resolved`, split-horizon DNS — so the request lands at an IP the operator did not vet. Drop the flag or pair it with `--resolve` pinning.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `curl https://api/ --dns-servers=1.1.1.1`",
			input: `curl https://api/ --dns-servers=1.1.1.1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1914",
					Message: "`curl --dns-servers` bypasses the host's resolver chain — `/etc/hosts`, `systemd-resolved`, split-horizon DNS — so the request lands at an IP the operator did not vet. Drop the flag or pair it with `--resolve` pinning.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1914")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1915(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mdadm --detail $MD` (read only)",
			input:    `mdadm --detail $MD`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mdadm --examine $DISK` (read only)",
			input:    `mdadm --examine $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mdadm --zero-superblock $DISK` (mangled)",
			input: `mdadm --zero-superblock $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1915",
					Message: "`mdadm --zero-superblock` drops RAID metadata or halts a live array — mounted root or /boot panics the host; a stale superblock scrambles data on next `--create`. Snapshot `mdadm --detail --export` first and keep behind a runbook.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:  "invalid — `mdadm -S $MD` (stop array)",
			input: `mdadm -S $MD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1915",
					Message: "`mdadm -S` drops RAID metadata or halts a live array — mounted root or /boot panics the host; a stale superblock scrambles data on next `--create`. Snapshot `mdadm --detail --export` first and keep behind a runbook.",
					Line:    1,
					Column:  7,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1915")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1916(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt NULL_GLOB` (explicit default)",
			input:    `unsetopt NULL_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_GLOB` (unrelated)",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt NULL_GLOB`",
			input: `setopt NULL_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1916",
					Message: "`setopt NULL_GLOB` makes every later unmatched glob silently empty — `cp *.log /dest` rewrites to `cp /dest`, `rm *.tmp` becomes argv-too-short. Use per-glob `*(N)`, or `setopt LOCAL_OPTIONS NULL_GLOB` in a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_NULL_GLOB`",
			input: `unsetopt NO_NULL_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1916",
					Message: "`unsetopt NO_NULL_GLOB` makes every later unmatched glob silently empty — `cp *.log /dest` rewrites to `cp /dest`, `rm *.tmp` becomes argv-too-short. Use per-glob `*(N)`, or `setopt LOCAL_OPTIONS NULL_GLOB` in a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1916")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1917(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `iw dev wlan0 link` (passive link info)",
			input:    `iw dev wlan0 link`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `iwlist wlan0 channel`",
			input:    `iwlist wlan0 channel`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `iw dev wlan0 scan`",
			input: `iw dev wlan0 scan`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1917",
					Message: "`iw dev <if> scan` runs an active probe-request sweep — interrupts the current association and broadcasts the host to every nearby AP. Use cached `iw dev $IF link` for passive queries.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `iwlist wlan0 scanning`",
			input: `iwlist wlan0 scanning`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1917",
					Message: "`iwlist <if> scan` runs an active probe-request sweep — interrupts the current association and broadcasts the host to every nearby AP. Use cached `iw dev $IF link` for passive queries.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1917")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1918(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt HIST_SUBST_PATTERN` (explicit default)",
			input:    `unsetopt HIST_SUBST_PATTERN`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_HISTORY` (unrelated)",
			input:    `setopt EXTENDED_HISTORY`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt HIST_SUBST_PATTERN`",
			input: `setopt HIST_SUBST_PATTERN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1918",
					Message: "`setopt HIST_SUBST_PATTERN` switches `:s` history/param modifiers to pattern matching — literal `*`/`?`/`^` suddenly act as glob metacharacters. Keep it off; use `${var//pat/rep}` when you actually want pattern substitution.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_HIST_SUBST_PATTERN`",
			input: `unsetopt NO_HIST_SUBST_PATTERN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1918",
					Message: "`unsetopt NO_HIST_SUBST_PATTERN` switches `:s` history/param modifiers to pattern matching — literal `*`/`?`/`^` suddenly act as glob metacharacters. Keep it off; use `${var//pat/rep}` when you actually want pattern substitution.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1918")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1919(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ss -tunlp` (read-only socket list)",
			input:    `ss -tunlp`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ss state established` (preview)",
			input:    `ss state established`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ss -K state close-wait`",
			input: `ss -K state close-wait`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1919",
					Message: "`ss -K` terminates every socket the filter matches — broad filters (`state established`, `dport 22`) kill the running SSH session. Preview with the same filter minus `-K`, and pin to a specific dst/port/state tuple.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:  "invalid — `ss --kill state established`",
			input: `ss --kill state established`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1919",
					Message: "`ss -K` terminates every socket the filter matches — broad filters (`state established`, `dport 22`) kill the running SSH session. Preview with the same filter minus `-K`, and pin to a specific dst/port/state tuple.",
					Line:    1,
					Column:  5,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1919")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1920(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt VERBOSE` (explicit default)",
			input:    `unsetopt VERBOSE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_HISTORY` (unrelated)",
			input:    `setopt EXTENDED_HISTORY`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt VERBOSE`",
			input: `setopt VERBOSE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1920",
					Message: "`setopt VERBOSE` echoes every executed command to stderr — any line that mentions a password, token, or API key leaks with the trace. Remove and use `printf` / a logger, or scope via `setopt LOCAL_OPTIONS VERBOSE` in a helper.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_VERBOSE`",
			input: `unsetopt NO_VERBOSE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1920",
					Message: "`unsetopt NO_VERBOSE` echoes every executed command to stderr — any line that mentions a password, token, or API key leaks with the trace. Remove and use `printf` / a logger, or scope via `setopt LOCAL_OPTIONS VERBOSE` in a helper.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1920")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1921(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `systemctl stop myapp`",
			input:    `systemctl stop myapp`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `systemctl kill -s HUP myapp` (graceful reload)",
			input:    `systemctl kill -s HUP myapp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `systemctl kill -s KILL myapp`",
			input: `systemctl kill -s KILL myapp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1921",
					Message: "`systemctl kill -s KILL` bypasses `ExecStop=` and `TimeoutStopSec=` — lockfiles, sockets, and shm segments survive and the next restart often fails with \"address already in use\". Use `systemctl stop` or `restart` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `systemctl kill myapp --signal=SIGKILL`",
			input: `systemctl kill myapp --signal=SIGKILL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1921",
					Message: "`systemctl kill --signal=SIGKILL` bypasses `ExecStop=` and `TimeoutStopSec=` — lockfiles, sockets, and shm segments survive and the next restart often fails with \"address already in use\". Use `systemctl stop` or `restart` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1921")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1922(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `rpm --import /tmp/key.asc` (local file)",
			input:    `rpm --import /tmp/key.asc`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rpm --import https://pinned.example/key.asc`",
			input:    `rpm --import https://pinned.example/key.asc`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `rpm --import http://repo.example/key.asc`",
			input: `rpm --import http://repo.example/key.asc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1922",
					Message: "`rpm --import http://repo.example/key.asc` fetches a GPG key over plaintext — on-path attackers swap it, every future signed package installs. Use `https://` from a pinned origin, or `gpg --verify` against a known fingerprint.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:  "invalid — `rpmkeys --import ftp://repo.example/key.asc`",
			input: `rpmkeys --import ftp://repo.example/key.asc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1922",
					Message: "`rpm --import ftp://repo.example/key.asc` fetches a GPG key over plaintext — on-path attackers swap it, every future signed package installs. Use `https://` from a pinned origin, or `gpg --verify` against a known fingerprint.",
					Line:    1,
					Column:  10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1922")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1923(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt PRINT_EXIT_VALUE` (explicit default)",
			input:    `unsetopt PRINT_EXIT_VALUE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_HISTORY` (unrelated)",
			input:    `setopt EXTENDED_HISTORY`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt PRINT_EXIT_VALUE`",
			input: `setopt PRINT_EXIT_VALUE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1923",
					Message: "`setopt PRINT_EXIT_VALUE` prints `zsh: exit N` on stderr for every non-zero exit — silent grep/test/curl probes suddenly leak status, and tools parsing stderr see interleaved shell chatter. Remove; use `|| printf …` per call.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_PRINT_EXIT_VALUE`",
			input: `unsetopt NO_PRINT_EXIT_VALUE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1923",
					Message: "`unsetopt NO_PRINT_EXIT_VALUE` prints `zsh: exit N` on stderr for every non-zero exit — silent grep/test/curl probes suddenly leak status, and tools parsing stderr see interleaved shell chatter. Remove; use `|| printf …` per call.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1923")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1924(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `virsh list --all` (read-only domain list)",
			input:    `virsh list --all`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `virt-top -d 5` (live view)",
			input:    `virt-top -d 5`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `virt-cat -d mydomain /etc/shadow`",
			input: `virt-cat -d mydomain /etc/shadow`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1924",
					Message: "`virt-cat` reads/writes the VM disk directly from the host — bypasses in-guest permissions, audit, and LUKS; a live VM risks corruption from double-mount. Snapshot first, work on the clone, prefer in-guest tooling.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `guestmount -d mydomain -i /mnt`",
			input: `guestmount -d mydomain -i /mnt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1924",
					Message: "`guestmount` reads/writes the VM disk directly from the host — bypasses in-guest permissions, audit, and LUKS; a live VM risks corruption from double-mount. Snapshot first, work on the clone, prefer in-guest tooling.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1924")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1925(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt EQUALS` (explicit default)",
			input:    `setopt EQUALS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt EQUALS`",
			input: `unsetopt EQUALS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1925",
					Message: "`unsetopt EQUALS` turns off `=cmd` path expansion and tilde-after-colon — `=python`/`=ls` become literals and `PATH=~/bin:$PATH` stops tilde-expanding. Keep on; scope with `setopt LOCAL_OPTIONS; unsetopt EQUALS` inside a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_EQUALS`",
			input: `setopt NO_EQUALS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1925",
					Message: "`setopt NO_EQUALS` turns off `=cmd` path expansion and tilde-after-colon — `=python`/`=ls` become literals and `PATH=~/bin:$PATH` stops tilde-expanding. Keep on; scope with `setopt LOCAL_OPTIONS; unsetopt EQUALS` inside a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1925")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1926(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `init 3` (multi-user, legacy but non-destructive)",
			input:    `init 3`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `systemctl reboot`",
			input:    `systemctl reboot`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `init 0` (halt)",
			input: `init 0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1926",
					Message: "`init 0` changes runlevel — `0` halts, `6` reboots, `1`/`S` drops to single-user. Use `systemctl poweroff`/`reboot`/`rescue` or `shutdown -h +N` so reviewers can read the intent.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `telinit 6` (reboot)",
			input: `telinit 6`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1926",
					Message: "`telinit 6` changes runlevel — `0` halts, `6` reboots, `1`/`S` drops to single-user. Use `systemctl poweroff`/`reboot`/`rescue` or `shutdown -h +N` so reviewers can read the intent.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1926")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1927(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `xfreerdp /u:alice /v:host.example` (no password)",
			input:    `xfreerdp /u:alice /v:host.example`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rdesktop -u alice host.example` (prompts)",
			input:    `rdesktop -u alice host.example`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `xfreerdp /u:alice /p:$PASS /v:host`",
			input: `xfreerdp /u:alice /p:$PASS /v:host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1927",
					Message: "`xfreerdp /p:$PASS` puts the RDP password in argv — visible in `ps`, `/proc`, and shell history. Pipe via `/from-stdin`, read from a protected `.rdp` file, or use NLA with a cached credential.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `rdesktop -u alice -p hunter2 host.example`",
			input: `rdesktop -u alice -p hunter2 host.example`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1927",
					Message: "`rdesktop -p hunter2` puts the RDP password in argv — visible in `ps`, `/proc`, and shell history. Pipe via `/from-stdin`, read from a protected `.rdp` file, or use NLA with a cached credential.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1927")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1928(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt SHARE_HISTORY` (explicit default)",
			input:    `unsetopt SHARE_HISTORY`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt INC_APPEND_HISTORY` (safer alternative)",
			input:    `setopt INC_APPEND_HISTORY`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt SHARE_HISTORY`",
			input: `setopt SHARE_HISTORY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1928",
					Message: "`setopt SHARE_HISTORY` flushes every command into every sibling zsh session — secrets typed in one terminal surface in `fc -l` of every other. Prefer `setopt INC_APPEND_HISTORY` plus `HIST_IGNORE_SPACE` for safer isolation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_SHARE_HISTORY`",
			input: `unsetopt NO_SHARE_HISTORY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1928",
					Message: "`unsetopt NO_SHARE_HISTORY` flushes every command into every sibling zsh session — secrets typed in one terminal surface in `fc -l` of every other. Prefer `setopt INC_APPEND_HISTORY` plus `HIST_IGNORE_SPACE` for safer isolation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1928")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1929(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cpio -o -H newc` (create, not extract)",
			input:    `cpio -o -H newc`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cpio -i --no-absolute-filenames`",
			input:    `cpio -i --no-absolute-filenames`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `cpio -i -d`",
			input: `cpio -i -d`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1929",
					Message: "`cpio -i` extracts paths verbatim — absolute and `..` entries escape the target dir. Pass `--no-absolute-filenames` and stage into a scratch dir before `mv` into place.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `cpio -idmv` (clustered)",
			input: `cpio -idmv`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1929",
					Message: "`cpio -i` extracts paths verbatim — absolute and `..` entries escape the target dir. Pass `--no-absolute-filenames` and stage into a scratch dir before `mv` into place.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1929")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1930(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt HASH_CMDS` (explicit default)",
			input:    `setopt HASH_CMDS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt HASH_CMDS`",
			input: `unsetopt HASH_CMDS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1930",
					Message: "`unsetopt HASH_CMDS` re-walks `$PATH` on every call — tens to hundreds of ms per command on slow filesystems. Keep it on; use `rehash` or `hash -r` to invalidate the cache after a targeted binary swap.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_HASH_CMDS`",
			input: `setopt NO_HASH_CMDS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1930",
					Message: "`setopt NO_HASH_CMDS` re-walks `$PATH` on every call — tens to hundreds of ms per command on slow filesystems. Keep it on; use `rehash` or `hash -r` to invalidate the cache after a targeted binary swap.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1930")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1931(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ip netns list`",
			input:    `ip netns list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ip netns exec red ping host`",
			input:    `ip netns exec red ping host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ip netns delete red`",
			input: `ip netns delete red`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1931",
					Message: "`ip netns delete` tears down every interface, veth, tunnel, and WireGuard peer inside the namespace. Stop the workloads first and verify `ip -n $NS link` is empty before deleting.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ip netns del $NS`",
			input: `ip netns del $NS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1931",
					Message: "`ip netns del` tears down every interface, veth, tunnel, and WireGuard peer inside the namespace. Stop the workloads first and verify `ip -n $NS link` is empty before deleting.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1931")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1932(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt GLOBAL_EXPORT` (explicit default)",
			input:    `setopt GLOBAL_EXPORT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt GLOBAL_EXPORT`",
			input: `unsetopt GLOBAL_EXPORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1932",
					Message: "`unsetopt GLOBAL_EXPORT` makes `typeset -x` exports function-local — helper functions that set `PATH`/`VIRTUAL_ENV`/`AWS_*` no longer propagate to callers. Keep it on; scope temporary exports in a subshell instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_GLOBAL_EXPORT`",
			input: `setopt NO_GLOBAL_EXPORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1932",
					Message: "`setopt NO_GLOBAL_EXPORT` makes `typeset -x` exports function-local — helper functions that set `PATH`/`VIRTUAL_ENV`/`AWS_*` no longer propagate to callers. Keep it on; scope temporary exports in a subshell instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1932")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1933(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ipvsadm -L -n` (list)",
			input:    `ipvsadm -L -n`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ipvsadm --save` (backup)",
			input:    `ipvsadm --save`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ipvsadm -C`",
			input: `ipvsadm -C`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1933",
					Message: "`ipvsadm -C` wipes every IPVS virtual service and real-server binding — load balancing stops, clients see 5xx. Save via `ipvsadm --save`, drain specific services with `-D`, reserve `--clear` for break-glass.",
					Line:    1,
					Column:  9,
				},
			},
		},
		{
			name:  "invalid — `ipvsadm --clear now`",
			input: `ipvsadm --clear now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1933",
					Message: "`ipvsadm --clear` wipes every IPVS virtual service and real-server binding — load balancing stops, clients see 5xx. Save via `ipvsadm --save`, drain specific services with `-D`, reserve `--clear` for break-glass.",
					Line:    1,
					Column:  10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1933")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1934(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt AUTO_NAME_DIRS` (explicit default)",
			input:    `unsetopt AUTO_NAME_DIRS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt AUTO_CD` (unrelated)",
			input:    `setopt AUTO_CD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt AUTO_NAME_DIRS`",
			input: `setopt AUTO_NAME_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1934",
					Message: "`setopt AUTO_NAME_DIRS` auto-registers every absolute-path parameter as a named dir — `foo=/srv/data` makes `~foo` expand, `%~` prompts surface names the user never picked. Keep off; use `hash -d name=/path`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_AUTO_NAME_DIRS`",
			input: `unsetopt NO_AUTO_NAME_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1934",
					Message: "`unsetopt NO_AUTO_NAME_DIRS` auto-registers every absolute-path parameter as a named dir — `foo=/srv/data` makes `~foo` expand, `%~` prompts surface names the user never picked. Keep off; use `hash -d name=/path`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1934")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1935(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `apt autoremove --dry-run` (preview)",
			input:    `apt autoremove --dry-run`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `apt autoremove` (no purge, config files kept)",
			input:    `apt autoremove`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `apt autoremove --purge -y`",
			input: `apt autoremove --purge -y`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1935",
					Message: "`apt autoremove` strips packages the resolver thinks are unused plus their configs — uproots packages installed manually but never `apt-mark manual`-ed. Dry-run first, mark keepers, drop `--purge` in CI.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `zypper rm --clean-deps foo`",
			input: `zypper rm --clean-deps foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1935",
					Message: "`zypper autoremove` strips packages the resolver thinks are unused plus their configs — uproots packages installed manually but never `apt-mark manual`-ed. Dry-run first, mark keepers, drop `--purge` in CI.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1935")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1936(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt POSIX_ALIASES` (explicit default)",
			input:    `unsetopt POSIX_ALIASES`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_GLOB` (unrelated)",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt POSIX_ALIASES`",
			input: `setopt POSIX_ALIASES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1936",
					Message: "`setopt POSIX_ALIASES` narrows alias expansion to plain identifiers — aliases on `if`/`for`/`function` silently stop firing and any library that hooked them breaks. Scope with `emulate -LR sh` instead of flipping globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_POSIX_ALIASES`",
			input: `unsetopt NO_POSIX_ALIASES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1936",
					Message: "`unsetopt NO_POSIX_ALIASES` narrows alias expansion to plain identifiers — aliases on `if`/`for`/`function` silently stop firing and any library that hooked them breaks. Scope with `emulate -LR sh` instead of flipping globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1936")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1937(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `tmux list-sessions`",
			input:    `tmux list-sessions`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `tmux kill-window -t dev:1`",
			input:    `tmux kill-window -t dev:1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tmux kill-server`",
			input: `tmux kill-server`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1937",
					Message: "`tmux kill-server` tears down every detached process inside the session — builds, log tails, port-forwards get `SIGHUP`'d with no cleanup. Use `kill-window` for surgical removal or `systemd-run --scope` for workloads.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `screen -X quit`",
			input: `screen -X quit`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1937",
					Message: "`screen -X quit` tears down every detached process inside the session — builds, log tails, port-forwards get `SIGHUP`'d with no cleanup. Use `kill-window` for surgical removal or `systemd-run --scope` for workloads.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1937")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1938(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt POSIX_JOBS` (explicit default)",
			input:    `unsetopt POSIX_JOBS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt MONITOR` (unrelated)",
			input:    `setopt MONITOR`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt POSIX_JOBS`",
			input: `setopt POSIX_JOBS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1938",
					Message: "`setopt POSIX_JOBS` scopes `%n` / `fg` / `bg` / `disown` per subshell — parent jobs become invisible inside `(…)`. Leave off; scope POSIX job semantics with `emulate -LR sh` inside a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_POSIX_JOBS`",
			input: `unsetopt NO_POSIX_JOBS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1938",
					Message: "`unsetopt NO_POSIX_JOBS` scopes `%n` / `fg` / `bg` / `disown` per subshell — parent jobs become invisible inside `(…)`. Leave off; scope POSIX job semantics with `emulate -LR sh` inside a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1938")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1939(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `systemctl reboot`",
			input:    `systemctl reboot`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `shutdown -r +5`",
			input:    `shutdown -r +5`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `reboot -f`",
			input: `reboot -f`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1939",
					Message: "`reboot -f` fires `reboot(2)` immediately — no `ExecStop=`, no filesystem sync, no clean unmount. Databases replay from last checkpoint. Use `systemctl reboot` / `shutdown -r +N`; reserve `-f` for wedged recovery.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:  "invalid — `poweroff --force now`",
			input: `poweroff --force now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1939",
					Message: "`poweroff --force` fires `reboot(2)` immediately — no `ExecStop=`, no filesystem sync, no clean unmount. Databases replay from last checkpoint. Use `systemctl reboot` / `shutdown -r +N`; reserve `-f` for wedged recovery.",
					Line:    1,
					Column:  11,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1939")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1940(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt POSIX_ARGZERO` (explicit default)",
			input:    `unsetopt POSIX_ARGZERO`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt FUNCTION_ARGZERO` (different option)",
			input:    `setopt FUNCTION_ARGZERO`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt POSIX_ARGZERO`",
			input: `setopt POSIX_ARGZERO`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1940",
					Message: "`setopt POSIX_ARGZERO` freezes `$0` to the outer script name — loggers and `case $0` dispatch inside functions lose call-site context. Scope with `emulate -LR sh` instead of flipping globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_POSIX_ARGZERO`",
			input: `unsetopt NO_POSIX_ARGZERO`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1940",
					Message: "`unsetopt NO_POSIX_ARGZERO` freezes `$0` to the outer script name — loggers and `case $0` dispatch inside functions lose call-site context. Scope with `emulate -LR sh` instead of flipping globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1940")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1941(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `restic init --password-file /etc/restic.pass`",
			input:    `restic init --password-file /etc/restic.pass`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `restic snapshots`",
			input:    `restic snapshots`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `restic init --insecure-no-password now`",
			input: `restic init --insecure-no-password now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1941",
					Message: "`restic --insecure-no-password` creates an unencrypted repo — every operator with read access to the backend can reassemble the backed-up filesystem. Use `--password-file` / `--password-command` with a real passphrase.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `restic backup /data --insecure-no-password`",
			input: `restic backup /data --insecure-no-password`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1941",
					Message: "`restic --insecure-no-password` creates an unencrypted repo — every operator with read access to the backend can reassemble the backed-up filesystem. Use `--password-file` / `--password-command` with a real passphrase.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1941")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1942(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt CLOBBER_EMPTY` (explicit default)",
			input:    `unsetopt CLOBBER_EMPTY`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_CLOBBER` (unrelated)",
			input:    `setopt NO_CLOBBER`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CLOBBER_EMPTY`",
			input: `setopt CLOBBER_EMPTY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1942",
					Message: "`setopt CLOBBER_EMPTY` lets `>file` overwrite zero-length files even under `NO_CLOBBER` — `touch`ed lock / sentinel files lose their safety net. Keep off; use explicit `>|file` to bypass `NO_CLOBBER` for a specific write.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_CLOBBER_EMPTY`",
			input: `unsetopt NO_CLOBBER_EMPTY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1942",
					Message: "`unsetopt NO_CLOBBER_EMPTY` lets `>file` overwrite zero-length files even under `NO_CLOBBER` — `touch`ed lock / sentinel files lose their safety net. Keep off; use explicit `>|file` to bypass `NO_CLOBBER` for a specific write.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1942")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1943(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `systemd-nspawn -D /srv/container /bin/sh` (not booting)",
			input:    `systemd-nspawn -D /srv/container /bin/sh`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `machinectl start web`",
			input:    `machinectl start web`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `systemd-nspawn -b -D /srv/container`",
			input: `systemd-nspawn -b -D /srv/container`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1943",
					Message: "`systemd-nspawn -b` runs the rootfs's `/sbin/init` with minimal isolation — init scripts execute first and can probe the host. Use `-U`, drop caps with `--capability=`, pair with `--private-network`, prefer `machinectl start`.",
					Line:    1,
					Column:  16,
				},
			},
		},
		{
			name:  "invalid — `systemd-nspawn --boot -D $ROOT`",
			input: `systemd-nspawn --boot -D $ROOT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1943",
					Message: "`systemd-nspawn --boot` runs the rootfs's `/sbin/init` with minimal isolation — init scripts execute first and can probe the host. Use `-U`, drop caps with `--capability=`, pair with `--private-network`, prefer `machinectl start`.",
					Line:    1,
					Column:  17,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1943")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1944(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt IGNORE_EOF` (explicit default)",
			input:    `unsetopt IGNORE_EOF`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EMACS` (unrelated)",
			input:    `setopt EMACS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt IGNORE_EOF`",
			input: `setopt IGNORE_EOF`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1944",
					Message: "`setopt IGNORE_EOF` makes `Ctrl-D` stop terminating the shell — subshells, sudo holds, SSH tunnels linger after the parent left. Keep off; use `TMOUT=NN` for a timed stale-tty exit if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_IGNORE_EOF`",
			input: `unsetopt NO_IGNORE_EOF`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1944",
					Message: "`unsetopt NO_IGNORE_EOF` makes `Ctrl-D` stop terminating the shell — subshells, sudo holds, SSH tunnels linger after the parent left. Keep off; use `TMOUT=NN` for a timed stale-tty exit if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1944")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1945(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `bpftrace -l 'tracepoint:syscalls:*'` (list only)",
			input:    `bpftrace -l tracepoint:syscalls:*`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `bpftool prog show`",
			input:    `bpftool prog show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `bpftrace -e 'tracepoint:syscalls:sys_enter_openat{printf(...)}'`",
			input: `bpftrace -e 'tracepoint:syscalls:sys_enter_openat{printf(...)}'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1945",
					Message: "`bpftrace -e` loads an in-kernel eBPF program that can read arbitrary kernel/userland memory — every syscall arg, every TCP payload. Gate behind a runbook and prefer a short-lived `bpftrace -c CMD` over a pinned trace.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `bpftool prog load prog.o /sys/fs/bpf/spy`",
			input: `bpftool prog load prog.o /sys/fs/bpf/spy`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1945",
					Message: "`bpftool prog load` loads an in-kernel eBPF program that can read arbitrary kernel/userland memory — every syscall arg, every TCP payload. Gate behind a runbook and prefer a short-lived `bpftrace -c CMD` over a pinned trace.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1945")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1946(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt HUP` (explicit default)",
			input:    `setopt HUP`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt HUP`",
			input: `unsetopt HUP`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1946",
					Message: "`unsetopt HUP` stops the shell from `SIGHUP`-ing background jobs on exit — long pipelines and spawned daemons outlive the session, orphans accumulate. Use `disown` or `systemd-run --scope` on specific commands instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_HUP`",
			input: `setopt NO_HUP`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1946",
					Message: "`setopt NO_HUP` stops the shell from `SIGHUP`-ing background jobs on exit — long pipelines and spawned daemons outlive the session, orphans accumulate. Use `disown` or `systemd-run --scope` on specific commands instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1946")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1947(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ip xfrm state list`",
			input:    `ip xfrm state list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ip xfrm policy show`",
			input:    `ip xfrm policy show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ip xfrm state flush`",
			input: `ip xfrm state flush`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1947",
					Message: "`ip xfrm state flush` tears down every IPsec SA/policy — VPN tunnels drop, kernel stops encrypting, plaintext may leak during renegotiation. Scope via `ip xfrm state deleteall src $A dst $B`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ip xfrm policy flush`",
			input: `ip xfrm policy flush`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1947",
					Message: "`ip xfrm policy flush` tears down every IPsec SA/policy — VPN tunnels drop, kernel stops encrypting, plaintext may leak during renegotiation. Scope via `ip xfrm policy deleteall src $A dst $B`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1947")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1948(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ipmitool -I lan -H bmc -U admin -f /etc/ipmi.pass chassis status`",
			input:    `ipmitool -I lan -H bmc -U admin -f /etc/ipmi.pass chassis status`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ipmitool -E -H bmc chassis status` (env password)",
			input:    `ipmitool -E -H bmc chassis status`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ipmitool -H bmc -U admin -P hunter2 chassis status`",
			input: `ipmitool -H bmc -U admin -P hunter2 chassis status`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1948",
					Message: "`ipmitool -P hunter2` leaks the BMC password into argv — visible in `ps`/`/proc`/crash dumps. Use `-f <password_file>` (mode 0400) or `IPMI_PASSWORD=… ipmitool -E`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ipmitool -Phunter2 -H bmc chassis power status` (joined)",
			input: `ipmitool -Phunter2 -H bmc chassis power status`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1948",
					Message: "`ipmitool -Phunter2` leaks the BMC password into argv — visible in `ps`/`/proc`/crash dumps. Use `-f <password_file>` (mode 0400) or `IPMI_PASSWORD=… ipmitool -E`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1948")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1949(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `rmmod nft_chain_nat` (no force)",
			input:    `rmmod nft_chain_nat`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `modprobe -r bluetooth`",
			input:    `modprobe -r bluetooth`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `rmmod -f nvidia`",
			input: `rmmod -f nvidia`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1949",
					Message: "`rmmod -f` tears down a module even when its refcount is non-zero — in-use drivers dangle, kernel oopses on the next callback. Stop holders first (`lsof`/`umount`/`ip link down`), then `rmmod` without `-f`.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:  "invalid — `rmmod --force foo`",
			input: `rmmod --force foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1949",
					Message: "`rmmod --force` tears down a module even when its refcount is non-zero — in-use drivers dangle, kernel oopses on the next callback. Stop holders first (`lsof`/`umount`/`ip link down`), then `rmmod` without `-f`.",
					Line:    1,
					Column:  8,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1949")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1950(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `tune2fs -l $DEV` (read only)",
			input:    `tune2fs -l $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `tune2fs -m 1 $DEV` (tiny but non-zero reserve)",
			input:    `tune2fs -m 1 $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tune2fs -O ^has_journal $DEV`",
			input: `tune2fs -O ^has_journal $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1950",
					Message: "`tune2fs -O ^has_journal` strips the journal — crash recovery needs a full `fsck -y` and may truncate files. Keep the default.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `tune2fs -m 0 $DEV`",
			input: `tune2fs -m 0 $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1950",
					Message: "`tune2fs -m 0` zeroes the root reserve — a full fs leaves no headroom for `journald`/`apt`/root shells. Keep the default.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1950")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1951(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ceph osd pool ls detail`",
			input:    `ceph osd pool ls detail`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ceph -s` (health)",
			input:    `ceph -s`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ceph osd pool delete rbd rbd --yes-i-really-really-mean-it now` (mangled)",
			input: `ceph osd pool delete rbd rbd --yes-i-really-really-mean-it now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1951",
					Message: "`ceph … --yes-i-really-really-mean-it` automates the double-safety phrase — a typo or stale loop silently deletes production pools. Run deletions interactively, or spell the pool name in a runbook commit.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ceph config-key rm key --yes-i-really-mean-it now`",
			input: `ceph config-key rm key --yes-i-really-mean-it now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1951",
					Message: "`ceph … --yes-i-really-really-mean-it` automates the double-safety phrase — a typo or stale loop silently deletes production pools. Run deletions interactively, or spell the pool name in a runbook commit.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1951")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1952(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `zfs set sync=standard tank/data`",
			input:    `zfs set sync=standard tank/data`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `zfs get sync tank` (read only)",
			input:    `zfs get sync tank`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `zfs set sync=disabled tank/pg`",
			input: `zfs set sync=disabled tank/pg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1952",
					Message: "`zfs set sync=disabled` turns `fsync()` into a no-op — DBs (PostgreSQL/MariaDB/etcd) lose committed transactions on crash. Leave sync at `standard`; use a SLOG vdev if latency is the concern.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `zfs set sync=disabled $POOL`",
			input: `zfs set sync=disabled $POOL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1952",
					Message: "`zfs set sync=disabled` turns `fsync()` into a no-op — DBs (PostgreSQL/MariaDB/etcd) lose committed transactions on crash. Leave sync at `standard`; use a SLOG vdev if latency is the concern.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1952")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1953(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mount --make-private /sys`",
			input:    `mount --make-private /sys`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mount -t tmpfs tmpfs /run`",
			input:    `mount -t tmpfs tmpfs /run`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mount --make-shared /data`",
			input: `mount --make-shared /data`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1953",
					Message: "`mount --make-shared` puts the mount in a shared-subtree group — later bind-mounts propagate to every peer, including containers. Classic escape stepping stone. Use `--make-private` on sensitive paths.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:  "invalid — `mount /srv --make-rshared`",
			input: `mount /srv --make-rshared`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1953",
					Message: "`mount --make-rshared` puts the mount in a shared-subtree group — later bind-mounts propagate to every peer, including containers. Classic escape stepping stone. Use `--make-private` on sensitive paths.",
					Line:    1,
					Column:  13,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1953")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1954(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setfattr -n user.comment -v 'hello' /tmp/f`",
			input:    `setfattr -n user.comment -v 'hello' /tmp/f`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `getfattr -d /tmp/f` (read-only sibling)",
			input:    `getfattr -d /tmp/f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setfattr -n security.capability -v $VAL /usr/local/bin/app`",
			input: `setfattr -n security.capability -v $VAL /usr/local/bin/app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1954",
					Message: "`setfattr -n security.capability` writes the raw kernel xattr — bypasses `setcap`/`chcon`/`evmctl` validation and audit. Use the purpose-built tool.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setfattr -n security.selinux -v $CTX /etc/app`",
			input: `setfattr -n security.selinux -v $CTX /etc/app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1954",
					Message: "`setfattr -n security.selinux` writes the raw kernel xattr — bypasses `setcap`/`chcon`/`evmctl` validation and audit. Use the purpose-built tool.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1954")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1955(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `rfkill list`",
			input:    `rfkill list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rfkill unblock all`",
			input:    `rfkill unblock all`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `rfkill block all`",
			input: `rfkill block all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1955",
					Message: "`rfkill block all` hard-downs the radio(s) — host drops off the network in one call. Scope to the radio type that really needs it and schedule an `at now + N minutes` unblock for self-recovery.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `rfkill block wifi`",
			input: `rfkill block wifi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1955",
					Message: "`rfkill block wifi` hard-downs the radio(s) — host drops off the network in one call. Scope to the radio type that really needs it and schedule an `at now + N minutes` unblock for self-recovery.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1955")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1956(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `tailscale status`",
			input:    `tailscale status`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `tailscale up --auth-key=file:/etc/ts.key` (file source)",
			input:    `tailscale up --auth-key=file:/etc/ts.key`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tailscale up host --auth-key=tskey-auth-abc123`",
			input: `tailscale up host --auth-key=tskey-auth-abc123`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1956",
					Message: "`tailscale --auth-key=tskey-auth-abc123` puts the pre-auth key in argv — visible in `ps`/`/proc`/history/crash dumps. Use `--auth-key=file:/etc/ts.key` (mode 0400) or `--authkey-env=TS_AUTHKEY`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `tailscale up host --authkey tskey-ZZZ`",
			input: `tailscale up host --authkey tskey-ZZZ`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1956",
					Message: "`tailscale --authkey tskey-ZZZ` puts the pre-auth key in argv — visible in `ps`/`/proc`/history/crash dumps. Use `--auth-key=file:/etc/ts.key` (mode 0400) or `--authkey-env=TS_AUTHKEY`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1956")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1957(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `lvchange -ay data/home` (activate)",
			input:    `lvchange -ay data/home`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `vgchange -ay data`",
			input:    `vgchange -ay data`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `lvchange -an data/home`",
			input: `lvchange -an data/home`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1957",
					Message: "`lvchange -an` deactivates the LV/VG — unflushed writes on a mounted fs may be lost, open fds see EIO. Umount and stop holders first, verify with `lsof`/`fuser`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `vgchange -an data`",
			input: `vgchange -an data`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1957",
					Message: "`vgchange -an` deactivates the LV/VG — unflushed writes on a mounted fs may be lost, open fds see EIO. Umount and stop holders first, verify with `lsof`/`fuser`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1957")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1958(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `helm upgrade myapp bitnami/nginx`",
			input:    `helm upgrade myapp bitnami/nginx`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `helm upgrade myapp ./chart --atomic --wait`",
			input:    `helm upgrade myapp ./chart --atomic --wait`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `helm upgrade myapp ./chart --force`",
			input: `helm upgrade myapp ./chart --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1958",
					Message: "`helm upgrade --force` is delete+create — pods die, PodDisruptionBudget is bypassed, Services reset their `clusterIP`. Use plain `helm upgrade` (three-way merge) or `--atomic`/`--wait` for a supervised roll.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `helm3 upgrade --install myapp ./chart --force`",
			input: `helm3 upgrade --install myapp ./chart --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1958",
					Message: "`helm upgrade --force` is delete+create — pods die, PodDisruptionBudget is bypassed, Services reset their `clusterIP`. Use plain `helm upgrade` (three-way merge) or `--atomic`/`--wait` for a supervised roll.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1958")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1959(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `trivy image alpine:3.20`",
			input:    `trivy image alpine:3.20`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `trivy image --download-db-only`",
			input:    `trivy image --download-db-only`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `trivy image alpine:3.20 --skip-db-update`",
			input: `trivy image alpine:3.20 --skip-db-update`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1959",
					Message: "`trivy --skip-db-update` scans against the cached DB — every CVE disclosed since last refresh is missed. Keep the default download, or run `trivy --download-db-only` once per day in a scheduled job.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `trivy image alpine:3.20 --skip-update`",
			input: `trivy image alpine:3.20 --skip-update`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1959",
					Message: "`trivy --skip-update` scans against the cached DB — every CVE disclosed since last refresh is missed. Keep the default download, or run `trivy --download-db-only` once per day in a scheduled job.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1959")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1960(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `az vm list`",
			input:    `az vm list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `aws ssm describe-instance-information`",
			input:    `aws ssm describe-instance-information`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `az vm run-command invoke -g rg -n vm --command-id RunShellScript --scripts $CMD`",
			input: `az vm run-command invoke -g rg -n vm --command-id RunShellScript --scripts $CMD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1960",
					Message: "`az vm run-command invoke` runs arbitrary shell on the VM via the cloud control plane — operator-composed command strings become IAM-driven RCE. Pin to a reviewed asset, template-escape input, require MFA.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `aws ssm send-command --document-name AWS-RunShellScript --parameters commands=$CMD`",
			input: `aws ssm send-command --document-name AWS-RunShellScript --parameters commands=$CMD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1960",
					Message: "`aws ssm send-command` runs arbitrary shell on the VM via the cloud control plane — operator-composed command strings become IAM-driven RCE. Pin to a reviewed asset, template-escape input, require MFA.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1960")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1961(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gcloud iam service-accounts list`",
			input:    `gcloud iam service-accounts list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gcloud auth print-access-token` (short-lived)",
			input:    `gcloud auth print-access-token`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gcloud iam service-accounts keys create key.json --iam-account=$SA`",
			input: `gcloud iam service-accounts keys create key.json --iam-account=$SA`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1961",
					Message: "`gcloud iam service-accounts keys create` mints a long-lived JSON key — no auto-rotate, no refresh. Prefer Workload Identity Federation, `--impersonate-service-account`, or the attached service account.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid — `gcloud iam service-accounts keys list` (read only)",
			input:    `gcloud iam service-accounts keys list`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1961")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1962(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kustomize build overlays/prod`",
			input:    `kustomize build overlays/prod`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kustomize build . --load-restrictor=LoadRestrictionsRootOnly`",
			input:    `kustomize build . --load-restrictor=LoadRestrictionsRootOnly`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kustomize build . --load-restrictor=LoadRestrictionsNone`",
			input: `kustomize build . --load-restrictor=LoadRestrictionsNone`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1962",
					Message: "`kustomize build --load-restrictor=LoadRestrictionsNone` drops path-root restriction — untrusted overlays can reference `../../secrets/prod.env` and pull them into the render. Keep the default; vendor sibling files into the overlay.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1962")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1963(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `npx typescript@5.4.2 --init`",
			input:    `npx typescript@5.4.2 --init`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pnpm dlx @vercel/ncc@0.38.1 build`",
			input:    `pnpm dlx @vercel/ncc@0.38.1 build`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `npx create-react-app demo`",
			input: `npx create-react-app demo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1963",
					Message: "`npx create-react-app` pulls the `latest` tag every run — a squatted or compromised package lands attacker code. Pin the version (`create-react-app@X.Y.Z`) or use a regular `npm install` + `./node_modules/.bin/`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pnpm dlx prettier`",
			input: `pnpm dlx prettier`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1963",
					Message: "`pnpm dlx prettier` pulls the `latest` tag every run — a squatted or compromised package lands attacker code. Pin the version (`prettier@X.Y.Z`) or use a regular `npm install` + `./node_modules/.bin/`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1963")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1964(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `uvx ruff@0.5.7 check`",
			input:    `uvx ruff@0.5.7 check`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pipx run 'black==24.8.0' --version`",
			input:    `pipx run 'black==24.8.0' --version`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `uvx ruff check`",
			input: `uvx ruff check`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1964",
					Message: "`uvx ruff` resolves to the PyPI `latest` release — a squatted name or compromised maintainer lands untested code. Pin `pkg==X.Y.Z` (or `pkg@X.Y.Z` for uv).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pipx run black .`",
			input: `pipx run black .`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1964",
					Message: "`pipx run black` resolves to the PyPI `latest` release — a squatted name or compromised maintainer lands untested code. Pin `pkg==X.Y.Z` (or `pkg@X.Y.Z` for uv).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1964")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1965(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `systemd-cryptenroll --wipe-slot=recovery $DEV`",
			input:    `systemd-cryptenroll --wipe-slot=recovery $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `systemd-cryptenroll --tpm2-device=auto $DEV`",
			input:    `systemd-cryptenroll --tpm2-device=auto $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `systemd-cryptenroll $DEV --wipe-slot=all`",
			input: `systemd-cryptenroll $DEV --wipe-slot=all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1965",
					Message: "`systemd-cryptenroll --wipe-slot=all` wipes every LUKS key slot (passphrase/recovery/TPM2/FIDO2) in one call. Enrol the new slot first, wipe a specific index, back up the header with `cryptsetup luksHeaderBackup`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1965")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1966(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `zpool import $POOL` (no force)",
			input:    `zpool import $POOL`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `zpool export $POOL`",
			input:    `zpool export $POOL`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `zpool import -f $POOL`",
			input: `zpool import -f $POOL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1966",
					Message: "`zpool import -f` bypasses hostid/txg safety — forced import of a pool already online elsewhere (SAN/HA) corrupts it; forced export drops in-flight txgs. `zfs unmount -a` first, then plain `zpool export`/`import`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `zpool export -f $POOL`",
			input: `zpool export -f $POOL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1966",
					Message: "`zpool export -f` bypasses hostid/txg safety — forced import of a pool already online elsewhere (SAN/HA) corrupts it; forced export drops in-flight txgs. `zfs unmount -a` first, then plain `zpool export`/`import`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1966")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1967(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt PROMPT_SUBST`",
			input:    `unsetopt PROMPT_SUBST`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_PROMPT_SUBST`",
			input:    `setopt NO_PROMPT_SUBST`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt PROMPT_SUBST`",
			input: `setopt PROMPT_SUBST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1967",
					Message: "`setopt PROMPT_SUBST` re-runs command substitution on every prompt redraw — a branch/host/dir value with `$(…)` executes each render. Prefer `%n`/`%d`/`%~`/`vcs_info`, or scope via `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_PROMPT_SUBST`",
			input: `unsetopt NO_PROMPT_SUBST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1967",
					Message: "`unsetopt NO_PROMPT_SUBST` re-runs command substitution on every prompt redraw — a branch/host/dir value with `$(…)` executes each render. Prefer `%n`/`%d`/`%~`/`vcs_info`, or scope via `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1967")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1968(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `dnf versionlock list`",
			input:    `dnf versionlock list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `dnf update --exclude=kernel`",
			input:    `dnf update --exclude=kernel`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `dnf versionlock add kernel`",
			input: `dnf versionlock add kernel`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1968",
					Message: "`dnf versionlock add` pins the rpm — blocks future CVE fixes for glibc/openssl/kernel. Prefer `--exclude` on a single transaction and schedule a `versionlock delete` review.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `yum versionlock add openssl`",
			input: `yum versionlock add openssl`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1968",
					Message: "`yum versionlock add` pins the rpm — blocks future CVE fixes for glibc/openssl/kernel. Prefer `--exclude` on a single transaction and schedule a `versionlock delete` review.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1968")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1969(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `zsh -c ':'` (no -f/-d flag)",
			input:    `zsh -c ':'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `zsh $SCRIPT`",
			input:    `zsh $SCRIPT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `zsh -f $SCRIPT`",
			input: `zsh -f $SCRIPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1969",
					Message: "`zsh -f` skips `/etc/zsh*` and `~/.zsh*` startup files — corp proxy/audit/`PATH` hardening silently dropped. For a pristine shell use `env -i zsh` with an explicit allow-list.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `zsh -d $SCRIPT`",
			input: `zsh -d $SCRIPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1969",
					Message: "`zsh -d` skips `/etc/zsh*` and `~/.zsh*` startup files — corp proxy/audit/`PATH` hardening silently dropped. For a pristine shell use `env -i zsh` with an explicit allow-list.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1969")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1970(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `losetup -r $LOOP $IMG` (readonly, no partscan)",
			input:    `losetup -r $LOOP $IMG`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sfdisk --dump $IMG` (offline parser)",
			input:    `sfdisk --dump $IMG`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `losetup -P $LOOP $IMG`",
			input: `losetup -P $LOOP $IMG`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1970",
					Message: "`losetup -P` asks the kernel to parse the partition table of the image — attacker-controlled bytes have tripped kernel CVEs. Use `fdisk -l`/`sfdisk --dump` offline first, scan only known-good images.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `kpartx -av $IMG`",
			input: `kpartx -av $IMG`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1970",
					Message: "`kpartx -a` asks the kernel to parse the partition table of the image — attacker-controlled bytes have tripped kernel CVEs. Use `fdisk -l`/`sfdisk --dump` offline first, scan only known-good images.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `partprobe $LOOP`",
			input: `partprobe $LOOP`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1970",
					Message: "`partprobe` asks the kernel to parse the partition table of the image — attacker-controlled bytes have tripped kernel CVEs. Use `fdisk -l`/`sfdisk --dump` offline first, scan only known-good images.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1970")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1971(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt GLOBAL_RCS` (keeps default on)",
			input:    `setopt GLOBAL_RCS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NO_GLOBAL_RCS` (restores default)",
			input:    `unsetopt NO_GLOBAL_RCS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt GLOBAL_RCS`",
			input: `unsetopt GLOBAL_RCS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1971",
					Message: "`unsetopt GLOBAL_RCS` tells Zsh to skip `/etc/zprofile`, `/etc/zshrc`, `/etc/zlogin`, `/etc/zlogout` — corp `PATH`/audit/umask/proxy config silently dropped. Keep on; scope pristine setup with `emulate -LR zsh` or `env -i zsh -f`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_GLOBAL_RCS`",
			input: `setopt NO_GLOBAL_RCS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1971",
					Message: "`setopt NO_GLOBAL_RCS` tells Zsh to skip `/etc/zprofile`, `/etc/zshrc`, `/etc/zlogin`, `/etc/zlogout` — corp `PATH`/audit/umask/proxy config silently dropped. Keep on; scope pristine setup with `emulate -LR zsh` or `env -i zsh -f`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1971")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1972(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `dmsetup ls`",
			input:    `dmsetup ls`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `dmsetup remove $NAME` (no force)",
			input:    `dmsetup remove $NAME`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `dmsetup remove_all`",
			input: `dmsetup remove_all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1972",
					Message: "`dmsetup remove_all` drops LVM/LUKS/multipath mappings while still in use — in-flight I/O returns `ENXIO`, metadata needs a reboot. `umount` + `vgchange -an` / `cryptsetup close` first, then `dmsetup remove`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `dmsetup remove -f $NAME`",
			input: `dmsetup remove -f $NAME`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1972",
					Message: "`dmsetup remove -f` drops LVM/LUKS/multipath mappings while still in use — in-flight I/O returns `ENXIO`, metadata needs a reboot. `umount` + `vgchange -an` / `cryptsetup close` first, then `dmsetup remove`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1972")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1973(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt POSIX_IDENTIFIERS` (restores default)",
			input:    `unsetopt POSIX_IDENTIFIERS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_POSIX_IDENTIFIERS` (restores default)",
			input:    `setopt NO_POSIX_IDENTIFIERS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt POSIX_IDENTIFIERS`",
			input: `setopt POSIX_IDENTIFIERS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1973",
					Message: "`setopt POSIX_IDENTIFIERS` restricts parameter names to ASCII; later `${café}`/`${π}` fail to parse and i18n-named libs stop loading. Scope with `emulate -LR sh` inside the helper instead of flipping globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_POSIX_IDENTIFIERS`",
			input: `unsetopt NO_POSIX_IDENTIFIERS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1973",
					Message: "`unsetopt NO_POSIX_IDENTIFIERS` restricts parameter names to ASCII; later `${café}`/`${π}` fail to parse and i18n-named libs stop loading. Scope with `emulate -LR sh` inside the helper instead of flipping globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1973")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1974(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ipset list`",
			input:    `ipset list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ipset destroy blocklist` (targeted name)",
			input:    `ipset destroy blocklist`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ipset flush`",
			input: `ipset flush`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1974",
					Message: "`ipset flush` drops named IP sets wholesale — iptables/nft rules that reference them fall through to the default policy (block-list empty, allow-list gone). Target by name; reload atomically via `ipset restore -! < snapshot`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ipset destroy` (no arg, wipes every set)",
			input: `ipset destroy`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1974",
					Message: "`ipset destroy` drops named IP sets wholesale — iptables/nft rules that reference them fall through to the default policy (block-list empty, allow-list gone). Target by name; reload atomically via `ipset restore -! < snapshot`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1974")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1975(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt EXEC` (default on, keeps running)",
			input:    `setopt EXEC`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NO_EXEC` (restores default)",
			input:    `unsetopt NO_EXEC`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt EXEC`",
			input: `unsetopt EXEC`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1975",
					Message: "`unsetopt EXEC` stops running commands but keeps parsing — every later line becomes a silent no-op. For syntax checks run `zsh -n script.zsh` from outside the script.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_EXEC`",
			input: `setopt NO_EXEC`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1975",
					Message: "`setopt NO_EXEC` stops running commands but keeps parsing — every later line becomes a silent no-op. For syntax checks run `zsh -n script.zsh` from outside the script.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1975")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1976(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `exportfs -ra` (re-sync)",
			input:    `exportfs -ra`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `exportfs -f` (flush cache after edit)",
			input:    `exportfs -f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `exportfs -au`",
			input: `exportfs -au`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1976",
					Message: "`exportfs -au` unexports live NFS shares — mounted clients see `ESTALE` on every open fd. Use `exportfs -f` after editing `/etc/exports`; reserve `-u`/`-au` for coordinated shutdowns.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `exportfs -u $HOST:$PATH`",
			input: `exportfs -u $HOST:$PATH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1976",
					Message: "`exportfs -u` unexports live NFS shares — mounted clients see `ESTALE` on every open fd. Use `exportfs -f` after editing `/etc/exports`; reserve `-u`/`-au` for coordinated shutdowns.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1976")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1977(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt CHASE_DOTS`",
			input:    `unsetopt CHASE_DOTS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_CHASE_DOTS`",
			input:    `setopt NO_CHASE_DOTS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CHASE_DOTS`",
			input: `setopt CHASE_DOTS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1977",
					Message: "`setopt CHASE_DOTS` makes `cd ..` physically resolve before walking up — blue/green `current` symlinks stop working for `../foo` lookups. Keep off; use `cd -P` one-shot when physical resolution is needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_CHASE_DOTS`",
			input: `unsetopt NO_CHASE_DOTS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1977",
					Message: "`unsetopt NO_CHASE_DOTS` makes `cd ..` physically resolve before walking up — blue/green `current` symlinks stop working for `../foo` lookups. Keep off; use `cd -P` one-shot when physical resolution is needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1977")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1978(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sftp $HOST`",
			input:    `sftp $HOST`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl -u user: https://$HOST/file`",
			input:    `curl -u user: https://$HOST/file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ftp $HOST` (owned by ZC1200)",
			input:    `ftp $HOST`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tftp $HOST`",
			input: `tftp $HOST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1978",
					Message: "`tftp` transfers over plaintext UDP/69 with no authentication — capture the payload, or push a crafted file under the expected name. Use a signed-payload `curl` over HTTPS and verify the signature before use.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1978")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1979(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt HIST_FCNTL_LOCK` (keeps default off)",
			input:    `unsetopt HIST_FCNTL_LOCK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_HIST_FCNTL_LOCK`",
			input:    `setopt NO_HIST_FCNTL_LOCK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt HIST_FCNTL_LOCK`",
			input: `setopt HIST_FCNTL_LOCK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1979",
					Message: "`setopt HIST_FCNTL_LOCK` routes `$HISTFILE` locking through POSIX `fcntl()` — on NFS home directories a hung `rpc.lockd` freezes every other shell at the next prompt. Keep off; enable only when `$HISTFILE` is on a local fs.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_HIST_FCNTL_LOCK`",
			input: `unsetopt NO_HIST_FCNTL_LOCK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1979",
					Message: "`unsetopt NO_HIST_FCNTL_LOCK` routes `$HISTFILE` locking through POSIX `fcntl()` — on NFS home directories a hung `rpc.lockd` freezes every other shell at the next prompt. Keep off; enable only when `$HISTFILE` is on a local fs.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1979")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1980(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `udevadm control --reload`",
			input:    `udevadm control --reload`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `udevadm trigger --action=change`",
			input:    `udevadm trigger --action=change`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `udevadm trigger --action=remove`",
			input: `udevadm trigger --action=remove`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1980",
					Message: "`udevadm trigger --action=remove` replays `remove` uevents across `/sys` — SATA/NIC/GPU nodes detach on a live host. Reload rules with `udevadm control --reload`; scope with `--subsystem-match=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `udevadm trigger -c remove`",
			input: `udevadm trigger -c remove`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1980",
					Message: "`udevadm trigger --action=remove` replays `remove` uevents across `/sys` — SATA/NIC/GPU nodes detach on a live host. Reload rules with `udevadm control --reload`; scope with `--subsystem-match=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1980")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1981(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `exec $BIN`",
			input:    `exec $BIN`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `exec $BIN arg1 arg2`",
			input:    `exec $BIN arg1 arg2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `exec -a login $BIN`",
			input: `exec -a login $BIN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1981",
					Message: "`exec -a NAME` sets `argv[0]` to `NAME` — `ps`/`top`/audit rules see the alias, not the real binary. Keep out of production scripts unless the alias is documented.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `exec -a $ALIAS $BIN`",
			input: `exec -a $ALIAS $BIN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1981",
					Message: "`exec -a NAME` sets `argv[0]` to `NAME` — `ps`/`top`/audit rules see the alias, not the real binary. Keep out of production scripts unless the alias is documented.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1981")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1982(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ipcs -a` (list)",
			input:    `ipcs -a`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ipcrm -m $SHMID` (scoped)",
			input:    `ipcrm -m $SHMID`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ipcrm -a`",
			input: `ipcrm -a`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1982",
					Message: "`ipcrm -a` deletes every SysV shm/sem/mqueue object — Postgres/Oracle/shm-based services lose their backing store mid-transaction. Scope with `-m`/`-s`/`-q` on the specific ID after checking `ipcs -a`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1982")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1983(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt CSH_JUNKIE_QUOTES`",
			input:    `unsetopt CSH_JUNKIE_QUOTES`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_CSH_JUNKIE_QUOTES`",
			input:    `setopt NO_CSH_JUNKIE_QUOTES`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CSH_JUNKIE_QUOTES`",
			input: `setopt CSH_JUNKIE_QUOTES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1983",
					Message: "`setopt CSH_JUNKIE_QUOTES` makes every later multi-line `\"…\"`/`'…'` an error — inlined SQL/JSON payloads and autoloaded helpers stop parsing. Scope csh-style strictness with `emulate -LR csh` in the one helper that needs it.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_CSH_JUNKIE_QUOTES`",
			input: `unsetopt NO_CSH_JUNKIE_QUOTES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1983",
					Message: "`unsetopt NO_CSH_JUNKIE_QUOTES` makes every later multi-line `\"…\"`/`'…'` an error — inlined SQL/JSON payloads and autoloaded helpers stop parsing. Scope csh-style strictness with `emulate -LR csh` in the one helper that needs it.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1983")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1984(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sgdisk -p $DISK` (print)",
			input:    `sgdisk -p $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sgdisk --backup=/root/disk.gpt $DISK`",
			input:    `sgdisk --backup=/root/disk.gpt $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `sgdisk -Z $DISK`",
			input: `sgdisk -Z $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1984",
					Message: "`sgdisk -Z` erases the GPT on the target device — a wrong `$DISK` detaches every partition/LVM/LUKS header and bricks boot. `lsblk`/`blkid` preflight, `--backup` the old table, and test with `-t`/`--pretend` first.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:  "invalid — `sgdisk -o $DISK`",
			input: `sgdisk -o $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1984",
					Message: "`sgdisk -o` erases the GPT on the target device — a wrong `$DISK` detaches every partition/LVM/LUKS header and bricks boot. `lsblk`/`blkid` preflight, `--backup` the old table, and test with `-t`/`--pretend` first.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:  "invalid — `sgdisk --zap-all $DISK`",
			input: `sgdisk --zap-all $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1984",
					Message: "`sgdisk --zap-all` erases the GPT on the target device — a wrong `$DISK` detaches every partition/LVM/LUKS header and bricks boot. `lsblk`/`blkid` preflight, `--backup` the old table, and test with `-t`/`--pretend` first.",
					Line:    1,
					Column:  9,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1984")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1985(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt SH_FILE_EXPANSION` (default)",
			input:    `unsetopt SH_FILE_EXPANSION`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_SH_FILE_EXPANSION`",
			input:    `setopt NO_SH_FILE_EXPANSION`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt SH_FILE_EXPANSION`",
			input: `setopt SH_FILE_EXPANSION`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1985",
					Message: "`setopt SH_FILE_EXPANSION` flips expansion order to POSIX — a `~` or `=cmd` sitting inside a `$VAR` value suddenly resolves, so a user-typed `~other/.cache` escapes into another home. Scope with `emulate -LR sh`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_SH_FILE_EXPANSION`",
			input: `unsetopt NO_SH_FILE_EXPANSION`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1985",
					Message: "`unsetopt NO_SH_FILE_EXPANSION` flips expansion order to POSIX — a `~` or `=cmd` sitting inside a `$VAR` value suddenly resolves, so a user-typed `~other/.cache` escapes into another home. Scope with `emulate -LR sh`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1985")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1986(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `touch $FILE` (current clock)",
			input:    `touch $FILE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `touch -c $FILE` (no create, current clock)",
			input:    `touch -c $FILE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `touch -d now $FILE`",
			input: `touch -d now $FILE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1986",
					Message: "`touch -d` writes a specific atime/mtime — also the classic \"age the dropped file\" antiforensics pattern. Derive from `$SOURCE_DATE_EPOCH` where the intent is deterministic.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `touch -t 202401011200 $FILE`",
			input: `touch -t 202401011200 $FILE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1986",
					Message: "`touch -t` writes a specific atime/mtime — also the classic \"age the dropped file\" antiforensics pattern. Derive from `$SOURCE_DATE_EPOCH` where the intent is deterministic.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `touch -r $REF $FILE`",
			input: `touch -r $REF $FILE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1986",
					Message: "`touch -r` writes a specific atime/mtime — also the classic \"age the dropped file\" antiforensics pattern. Derive from `$SOURCE_DATE_EPOCH` where the intent is deterministic.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1986")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1987(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt BRACE_CCL` (default)",
			input:    `unsetopt BRACE_CCL`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_BRACE_CCL`",
			input:    `setopt NO_BRACE_CCL`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt BRACE_CCL`",
			input: `setopt BRACE_CCL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1987",
					Message: "`setopt BRACE_CCL` promotes single-character braces to csh-style classes — `{a-z}` now expands to every letter, `{ABC}` to `A B C`, breaking regex/hex/CI-name literals. Use `{a..z}` when a real range is wanted.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_BRACE_CCL`",
			input: `unsetopt NO_BRACE_CCL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1987",
					Message: "`unsetopt NO_BRACE_CCL` promotes single-character braces to csh-style classes — `{a-z}` now expands to every letter, `{ABC}` to `A B C`, breaking regex/hex/CI-name literals. Use `{a..z}` when a real range is wanted.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1987")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1988(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `nsupdate -k $KEYFILE`",
			input:    `nsupdate -k $KEYFILE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `nsupdate -v`",
			input:    `nsupdate -v`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `nsupdate -y HMAC-SHA256:name:c2VjcmV0`",
			input: `nsupdate -y HMAC-SHA256:name:c2VjcmV0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1988",
					Message: "`nsupdate -y …` puts the TSIG key in argv — `ps`, `/proc/*/cmdline`, and `$HISTFILE` all capture it. Use `nsupdate -k $KEYFILE` with a `0600` keyfile instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1988")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1989(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt REMATCH_PCRE` (default)",
			input:    `unsetopt REMATCH_PCRE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_REMATCH_PCRE`",
			input:    `setopt NO_REMATCH_PCRE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt REMATCH_PCRE`",
			input: `setopt REMATCH_PCRE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1989",
					Message: "`setopt REMATCH_PCRE` swaps `[[ =~ ]]` from POSIX ERE to PCRE — `\\b`, `\\d`, lookahead, `(?i)` change meaning across every later match. Prefer `pcre_match` where PCRE is needed or scope with `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_REMATCH_PCRE`",
			input: `unsetopt NO_REMATCH_PCRE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1989",
					Message: "`unsetopt NO_REMATCH_PCRE` swaps `[[ =~ ]]` from POSIX ERE to PCRE — `\\b`, `\\d`, lookahead, `(?i)` change meaning across every later match. Prefer `pcre_match` where PCRE is needed or scope with `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1989")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1990(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `openssl passwd -5 $PASS` (SHA-256-crypt)",
			input:    `openssl passwd -5 $PASS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `openssl passwd -6 $PASS` (SHA-512-crypt)",
			input:    `openssl passwd -6 $PASS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `openssl passwd -crypt $PASS`",
			input: `openssl passwd -crypt $PASS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1990",
					Message: "`openssl passwd -crypt` emits a broken hash format — DES/MD5 variants crack on a laptop. Use `-5` / `-6` or a KDF-based hasher (`mkpasswd -m yescrypt`, `htpasswd -B`, `argon2`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `openssl passwd -1 $PASS`",
			input: `openssl passwd -1 $PASS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1990",
					Message: "`openssl passwd -1` emits a broken hash format — DES/MD5 variants crack on a laptop. Use `-5` / `-6` or a KDF-based hasher (`mkpasswd -m yescrypt`, `htpasswd -B`, `argon2`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `openssl passwd -apr1 $PASS`",
			input: `openssl passwd -apr1 $PASS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1990",
					Message: "`openssl passwd -apr1` emits a broken hash format — DES/MD5 variants crack on a laptop. Use `-5` / `-6` or a KDF-based hasher (`mkpasswd -m yescrypt`, `htpasswd -B`, `argon2`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1990")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1991(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt CSH_NULLCMD` (default)",
			input:    `unsetopt CSH_NULLCMD`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_CSH_NULLCMD`",
			input:    `setopt NO_CSH_NULLCMD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CSH_NULLCMD`",
			input: `setopt CSH_NULLCMD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1991",
					Message: "`setopt CSH_NULLCMD` makes `> file` / `< file` (no command) a parse error — log truncation and bare-redirect idioms stop working. Write `: > file` explicitly for truncation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_CSH_NULLCMD`",
			input: `unsetopt NO_CSH_NULLCMD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1991",
					Message: "`unsetopt NO_CSH_NULLCMD` makes `> file` / `< file` (no command) a parse error — log truncation and bare-redirect idioms stop working. Write `: > file` explicitly for truncation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1991")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1992(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sudo $CMD` (targeted sudoers drop-in)",
			input:    `sudo $CMD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `pkexec $CMD arg`",
			input: `pkexec $CMD arg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1992",
					Message: "`pkexec` elevates via PolicyKit — no agent to prompt in a script, poor CVE history (pwnkit), split audit trail. Use `sudo` with a targeted `sudoers.d` drop-in or a systemd unit with `User=`/`AmbientCapabilities=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1992")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1993(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt KSH_TYPESET` (default)",
			input:    `unsetopt KSH_TYPESET`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_KSH_TYPESET`",
			input:    `setopt NO_KSH_TYPESET`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt KSH_TYPESET`",
			input: `setopt KSH_TYPESET`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1993",
					Message: "`setopt KSH_TYPESET` re-splits the RHS of every later `typeset`/`local` — `typeset path=$HOME/My Files` now treats `Files` as a second name. Scope with `emulate -LR ksh` inside the one helper that needs it.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_KSH_TYPESET`",
			input: `unsetopt NO_KSH_TYPESET`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1993",
					Message: "`unsetopt NO_KSH_TYPESET` re-splits the RHS of every later `typeset`/`local` — `typeset path=$HOME/My Files` now treats `Files` as a second name. Scope with `emulate -LR ksh` inside the one helper that needs it.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1993")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1994(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `lvreduce -L 10G $LV` (interactive confirm)",
			input:    `lvreduce -L 10G $LV`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `lvextend -L +10G $LV` (grow)",
			input:    `lvextend -L +10G $LV`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `lvreduce -f -L 10G $LV`",
			input: `lvreduce -f -L 10G $LV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1994",
					Message: "`lvreduce -f` skips the shrink-confirmation prompt — the filesystem above still believes the tail is allocated and the next mount sees corruption. Shrink fs first, or use `lvreduce --resizefs`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `lvreduce -y -L 10G $LV`",
			input: `lvreduce -y -L 10G $LV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1994",
					Message: "`lvreduce -y` skips the shrink-confirmation prompt — the filesystem above still believes the tail is allocated and the next mount sees corruption. Shrink fs first, or use `lvreduce --resizefs`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1994")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1995(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt BG_NICE` (keeps default on)",
			input:    `setopt BG_NICE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NO_BG_NICE`",
			input:    `unsetopt NO_BG_NICE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt BG_NICE`",
			input: `unsetopt BG_NICE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1995",
					Message: "`unsetopt BG_NICE` drops the `nice +5` that bg jobs get by default — a CPU-bound `cmd &` now competes with SSH/editor work. Wrap specific jobs with `nice -n 0` or a systemd `Nice=` unit instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_BG_NICE`",
			input: `setopt NO_BG_NICE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1995",
					Message: "`setopt NO_BG_NICE` drops the `nice +5` that bg jobs get by default — a CPU-bound `cmd &` now competes with SSH/editor work. Wrap specific jobs with `nice -n 0` or a systemd `Nice=` unit instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1995")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1996(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unshare -m $CMD` (mount namespace only)",
			input:    `unshare -m $CMD`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `bwrap --unshare-all $CMD` (rootless runtime)",
			input:    `bwrap --unshare-all $CMD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unshare -Ur $CMD`",
			input: `unshare -Ur $CMD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1996",
					Message: "`unshare -Ur` opens a user namespace and maps the caller to uid 0 inside it — also the standard opening move for many kernel-LPE chains. Route legit rootless runtimes through `bwrap`/`podman --rootless` so the intent is clear.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unshare -U $CMD`",
			input: `unshare -U $CMD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1996",
					Message: "`unshare -U` opens a user namespace and maps the caller to uid 0 inside it — also the standard opening move for many kernel-LPE chains. Route legit rootless runtimes through `bwrap`/`podman --rootless` so the intent is clear.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unshare -Urm $CMD` (short bundle)",
			input: `unshare -Urm $CMD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1996",
					Message: "`unshare -Urm` opens a user namespace and maps the caller to uid 0 inside it — also the standard opening move for many kernel-LPE chains. Route legit rootless runtimes through `bwrap`/`podman --rootless` so the intent is clear.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1996")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1997(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt HIST_NO_FUNCTIONS` (default)",
			input:    `unsetopt HIST_NO_FUNCTIONS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_HIST_NO_FUNCTIONS`",
			input:    `setopt NO_HIST_NO_FUNCTIONS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt HIST_NO_FUNCTIONS`",
			input: `setopt HIST_NO_FUNCTIONS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1997",
					Message: "`setopt HIST_NO_FUNCTIONS` drops function-definition commands from `$HISTFILE` — forensic trail loses the definition while the call that used it still shows. Scope hiding via `zshaddhistory` hook instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_HIST_NO_FUNCTIONS`",
			input: `unsetopt NO_HIST_NO_FUNCTIONS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1997",
					Message: "`unsetopt NO_HIST_NO_FUNCTIONS` drops function-definition commands from `$HISTFILE` — forensic trail loses the definition while the call that used it still shows. Scope hiding via `zshaddhistory` hook instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1997")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1998(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `tpm2 getcap algorithms`",
			input:    `tpm2 getcap algorithms`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `tpm2_pcrread sha256:0,1,2`",
			input:    `tpm2_pcrread sha256:0,1,2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tpm2_clear -c p`",
			input: `tpm2_clear -c p`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1998",
					Message: "`tpm2_clear` wipes the TPM storage hierarchy — every LUKS-TPM2 keyslot, `systemd-cryptenroll --tpm2-device` slot, and TPM-sealed TLS/sshd key is destroyed. No undo. Gate behind a recovery runbook.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `tpm2 clear -c p`",
			input: `tpm2 clear -c p`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1998",
					Message: "`tpm2 clear` wipes the TPM storage hierarchy — every LUKS-TPM2 keyslot, `systemd-cryptenroll --tpm2-device` slot, and TPM-sealed TLS/sshd key is destroyed. No undo. Gate behind a recovery runbook.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1998")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1999(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt AUTO_NAME_DIRS` (canonical name; handled by ZC1934)",
			input:    `setopt AUTO_NAME_DIRS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt AUTO_NAME_DIRS`",
			input:    `unsetopt AUTO_NAME_DIRS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt AUTO_NAMED_DIRS` (typo)",
			input: `setopt AUTO_NAMED_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1999",
					Message: "`setopt AUTO_NAMED_DIRS` is a typo — the real Zsh option is `AUTO_NAME_DIRS` (no trailing `D`, see ZC1934). Fix the spelling or drop the toggle; `hash -d NAME=PATH` is the explicit alternative.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_AUTO_NAMED_DIRS` (typo)",
			input: `unsetopt NO_AUTO_NAMED_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1999",
					Message: "`unsetopt NO_AUTO_NAMED_DIRS` is a typo — the real Zsh option is `AUTO_NAME_DIRS` (no trailing `D`, see ZC1934). Fix the spelling or drop the toggle; `hash -d NAME=PATH` is the explicit alternative.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1999")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
