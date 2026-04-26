// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1700(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ldapsearch -W (prompt)",
			input:    `ldapsearch -x -D cn=admin -W -b dc=example`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ldapsearch -y FILE",
			input:    `ldapsearch -x -D cn=admin -y /etc/ldap.password`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ldapsearch -w SECRET",
			input: `ldapsearch -x -D cn=admin -w SECRET -b dc=example`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1700",
					Message: "`ldapsearch -w PASSWORD` leaks the LDAP bind password into `ps` / `/proc/PID/cmdline` — use `-W` to prompt, `-y FILE` for a mode-0400 secret file, or SASL (`-Y GSSAPI`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ldapmodify -w SECRET",
			input: `ldapmodify -w SECRET -f change.ldif`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1700",
					Message: "`ldapmodify -w PASSWORD` leaks the LDAP bind password into `ps` / `/proc/PID/cmdline` — use `-W` to prompt, `-y FILE` for a mode-0400 secret file, or SASL (`-Y GSSAPI`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1700")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1701(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dpkg -l (list)",
			input:    `dpkg -l`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — apt install from repo",
			input:    `apt install mypkg`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dpkg -i local.deb",
			input: `dpkg -i local.deb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1701",
					Message: "`dpkg -i FILE.deb` runs the package without signature verification — `sha256sum -c` or `debsig-verify` the file first, or install via `apt install` from a signed repo.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — dpkg -i /tmp/download.deb",
			input: `dpkg -i /tmp/download.deb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1701",
					Message: "`dpkg -i FILE.deb` runs the package without signature verification — `sha256sum -c` or `debsig-verify` the file first, or install via `apt install` from a signed repo.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1701")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1702(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dpkg-reconfigure -f noninteractive",
			input:    `dpkg-reconfigure -f noninteractive tzdata`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — different command",
			input:    `dpkg -l`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dpkg-reconfigure tzdata (no frontend)",
			input: `dpkg-reconfigure tzdata`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1702",
					Message: "`dpkg-reconfigure` without `-f noninteractive` opens debconf dialogs — non-interactive pipelines hang. Pass `-f noninteractive` or export `DEBIAN_FRONTEND=noninteractive`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — dpkg-reconfigure --priority high (still interactive)",
			input: `dpkg-reconfigure -p high tzdata`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1702",
					Message: "`dpkg-reconfigure` without `-f noninteractive` opens debconf dialogs — non-interactive pipelines hang. Pass `-f noninteractive` or export `DEBIAN_FRONTEND=noninteractive`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1702")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1703(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sysctl net.ipv4.conf.all.rp_filter=1 (strict)",
			input:    `sysctl -w net.ipv4.conf.all.rp_filter=1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — sysctl unrelated knob",
			input:    `sysctl -w net.ipv4.tcp_syncookies=1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sysctl rp_filter=0",
			input: `sysctl -w net.ipv4.conf.all.rp_filter=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1703",
					Message: "`sysctl net.ipv4.conf.all.rp_filter=0` disables reverse-path filtering (anti-spoofing) — classic layer-3 attacks (spoofing / smurf / redirect tamper) reopen. Keep the protective default.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sysctl accept_source_route=1",
			input: `sysctl -w net.ipv4.conf.all.accept_source_route=1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1703",
					Message: "`sysctl net.ipv4.conf.all.accept_source_route=1` disables source-routed packet acceptance — classic layer-3 attacks (spoofing / smurf / redirect tamper) reopen. Keep the protective default.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1703")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1704(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — scoped CIDR",
			input:    `aws ec2 authorize-security-group-ingress --group-id sg-123 --protocol tcp --port 22 --cidr 10.0.0.0/8`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — describe-security-groups (different subcommand)",
			input:    `aws ec2 describe-security-groups`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — --cidr 0.0.0.0/0",
			input: `aws ec2 authorize-security-group-ingress --group-id sg-123 --protocol tcp --port 22 --cidr 0.0.0.0/0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1704",
					Message: "`aws ec2 authorize-security-group-ingress --cidr 0.0.0.0/0` opens the port to the entire internet — scope to a known source CIDR or `--source-group sg-…`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — --cidr-ipv6 ::/0",
			input: `aws ec2 authorize-security-group-ingress --group-id sg-123 --ip-permissions proto --cidr-ipv6 ::/0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1704",
					Message: "`aws ec2 authorize-security-group-ingress --cidr ::/0` opens the port to the entire internet — scope to a known source CIDR or `--source-group sg-…`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1704")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1705(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — awk without -i",
			input:    `awk '{print}' file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — awk -i with non-inplace path",
			input:    `awk -i /usr/share/awk/lib.awk '{print}' file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — awk -i inplace",
			input: `awk -i inplace '{print}' file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1705",
					Message: "`awk -i inplace` is gawk-only — fails on mawk / BSD awk / busybox awk. For portability rewrite as `awk … input > tmp && mv tmp input`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — gawk -i inplace -v",
			input: `gawk -i inplace -v x=1 '{print x}' file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1705",
					Message: "`awk -i inplace` is gawk-only — fails on mawk / BSD awk / busybox awk. For portability rewrite as `awk … input > tmp && mv tmp input`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1705")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1706(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — lvresize grow without -r",
			input:    `lvresize -L +2G vg/lv`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — lvresize shrink with -r",
			input:    `lvresize -L -2G -r vg/lv`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — lvextend (always grows)",
			input:    `lvextend -L +2G vg/lv`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — lvresize shrink without -r",
			input: `lvresize -L -2G vg/lv`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1706",
					Message: "`lvresize` shrinks the LV without `-r` / `--resizefs` — the filesystem on top is not shrunk first and writes past the new boundary corrupt metadata. Add `-r` (or shrink the FS manually first).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — lvreduce without -r",
			input: `lvreduce -L 1G vg/lv`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1706",
					Message: "`lvreduce` shrinks the LV without `-r` / `--resizefs` — the filesystem on top is not shrunk first and writes past the new boundary corrupt metadata. Add `-r` (or shrink the FS manually first).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1706")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1707(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — gpg --keyserver hkps:// trailing",
			input:    `gpg ABCD --keyserver hkps://keys.openpgp.org --recv-keys`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — gpg --recv-keys (default keyserver)",
			input:    `gpg --recv-keys ABCD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — gpg --keyserver hkp:// trailing",
			input: `gpg ABCD --keyserver hkp://keys.example.com --recv-keys`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1707",
					Message: "`gpg --keyserver hkp://…` is plaintext — a MITM swaps the key bytes. Use `hkps://keys.openpgp.org` or fetch over HTTPS and verify the fingerprint.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1707")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1708(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -L",
			input:    `find /var/log -mtime +30 -delete`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — find -L without destructive action",
			input:    `find -L /opt -name '*.bak'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -L … -delete",
			input: `find -L /var/log -mtime +30 -delete`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1708",
					Message: "`find -L … -delete/-exec` follows symlinks into unintended trees — drop `-L`, add `-xdev`, or scope the walk explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -L … -exec rm",
			input: `find -L /var/log -mtime +30 -exec rm -f`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1708",
					Message: "`find -L … -delete/-exec` follows symlinks into unintended trees — drop `-L`, add `-xdev`, or scope the walk explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1708")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1709(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — htpasswd -i (read stdin)",
			input:    `htpasswd -i /etc/nginx/.htpasswd user`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — htpasswd interactive (prompts)",
			input:    `htpasswd /etc/nginx/.htpasswd user`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — htpasswd -b user secret",
			input: `htpasswd -b /etc/nginx/.htpasswd user secret`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1709",
					Message: "`htpasswd -b USER PASSWORD` puts the password in argv — visible via `ps` / `/proc/PID/cmdline`. Use `htpasswd -i FILE USER` with the password piped on stdin instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — htpasswd -bB combined flags",
			input: `htpasswd -bB /etc/nginx/.htpasswd user secret`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1709",
					Message: "`htpasswd -b USER PASSWORD` puts the password in argv — visible via `ps` / `/proc/PID/cmdline`. Use `htpasswd -i FILE USER` with the password piped on stdin instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1709")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1710(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — journalctl --vacuum-time=2weeks (real retention)",
			input:    `journalctl -q --vacuum-time=2weeks`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — journalctl --vacuum-size=500M",
			input:    `journalctl -q --vacuum-size=500M`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — journalctl --vacuum-size=1 (wipe)",
			input: `journalctl -q --vacuum-size=1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1710",
					Message: "`journalctl --vacuum-size=1` flushes the systemd journal — classic audit-clear shape. Set retention in `/etc/systemd/journald.conf` (`SystemMaxUse=`, `MaxRetentionSec=`) instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — journalctl --vacuum-time=1s",
			input: `journalctl -m --vacuum-time=1s`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1710",
					Message: "`journalctl --vacuum-time=1s` flushes the systemd journal — classic audit-clear shape. Set retention in `/etc/systemd/journald.conf` (`SystemMaxUse=`, `MaxRetentionSec=`) instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1710")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1711(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — etcdctl del --prefix /app/",
			input:    `etcdctl del --prefix /app/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — etcdctl get --prefix \"\" (read-only)",
			input:    `etcdctl get --prefix ""`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — etcdctl del --prefix ""`,
			input: `etcdctl del --prefix ""`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1711",
					Message: "`etcdctl del --prefix \"\"` deletes the entire etcd keyspace (including kube-apiserver state) — scope to a specific namespace prefix and review with `get --prefix --keys-only` first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — etcdctl del --from-key ""`,
			input: `etcdctl del --from-key ""`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1711",
					Message: "`etcdctl del --from-key \"\"` deletes the entire etcd keyspace (including kube-apiserver state) — scope to a specific namespace prefix and review with `get --prefix --keys-only` first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1711")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1712(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — vault kv put with @file",
			input:    `vault kv put secret/app @secret.json`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — vault kv put with stdin sentinel",
			input:    `vault kv put secret/app password=-`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — vault kv put with non-secret key",
			input:    `vault kv put secret/app environment=prod`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — vault kv put password=hunter2",
			input: `vault kv put secret/app password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1712",
					Message: "`vault kv password=hunter2` puts the secret value in argv — visible to every local user. Use `password=@FILE` or `password=-` to read from disk / stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — vault write secret/app api_key=ABC",
			input: `vault write secret/app api_key=ABC123`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1712",
					Message: "`vault write api_key=ABC123` puts the secret value in argv — visible to every local user. Use `api_key=@FILE` or `api_key=-` to read from disk / stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1712")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1713(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — consul kv delete scoped",
			input:    `consul kv delete -recurse /app/staging/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — consul kv delete single key",
			input:    `consul kv delete /key`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — consul kv get (read-only)",
			input:    `consul kv get -recurse /`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — consul kv delete -recurse /",
			input: `consul kv delete -recurse /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1713",
					Message: "`consul kv delete -recurse /` removes the entire KV store — service discovery, ACL bootstrap, app config. Scope the prefix and snapshot (`consul snapshot save`) first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1713")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1714(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — gh repo delete without --yes (prompts)",
			input:    `gh repo delete owner/repo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — gh repo create",
			input:    `gh repo create owner/repo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — gh repo delete --yes",
			input: `gh repo delete owner/repo --yes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1714",
					Message: "`gh repo delete --yes` bypasses GitHub's confirmation — a typo or stale variable destroys the target with no soft-delete. Drop `--yes` so a human confirms.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — gh release delete --yes",
			input: `gh release delete v1.0 --yes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1714",
					Message: "`gh release delete --yes` bypasses GitHub's confirmation — a typo or stale variable destroys the target with no soft-delete. Drop `--yes` so a human confirms.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1714")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1715(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — Zsh-idiomatic `read \"var?prompt\"`",
			input:    `read "name?Enter your name: "`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `read -r line` (no -p)",
			input:    `read -r line`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `read var` (bare)",
			input:    `read var`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `read -p \"Prompt: \" name`",
			input: `read -p "Prompt: " name`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1715",
					Message: "`read -p` triggers Zsh's coprocess reader, not Bash's prompt — the variable stays empty. Use `read \"var?Prompt: \"` (the `?` after the variable name introduces the prompt).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `read -rp \"Prompt: \" name` (combined short flags)",
			input: `read -rp "Prompt: " name`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1715",
					Message: "`read -rp` triggers Zsh's coprocess reader, not Bash's prompt — the variable stays empty. Use `read \"var?Prompt: \"` (the `?` after the variable name introduces the prompt).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1715")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1716(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — Zsh-idiomatic `$CPUTYPE`",
			input:    `print -r -- $CPUTYPE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `uname -r` (kernel release, no Zsh equivalent)",
			input:    `uname -r`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `uname -m`",
			input: `uname -m`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1716",
					Message: "Use Zsh `$CPUTYPE` / `$MACHTYPE` instead of `uname -m` — parameter expansion avoids forking an external for an answer Zsh already cached at startup.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `uname -p`",
			input: `uname -p`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1716",
					Message: "Use Zsh `$CPUTYPE` / `$MACHTYPE` instead of `uname -p` — parameter expansion avoids forking an external for an answer Zsh already cached at startup.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1716")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1717(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `docker pull` (no bypass)",
			input:    `docker pull nginx:1.27`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `docker push` (no bypass)",
			input:    `docker push myorg/app:1.2.3`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `docker version --disable-content-trust` (not pull/push subcmd)",
			input:    `docker version --disable-content-trust`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `docker pull --disable-content-trust`",
			input: `docker pull --disable-content-trust myorg/app:1.2.3`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1717",
					Message: "`docker pull --disable-content-trust` overrides `DOCKER_CONTENT_TRUST=1` — unsigned image moves into the registry or local store. Sign the artifact (`docker trust sign`) instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `docker push --disable-content-trust`",
			input: `docker push --disable-content-trust myorg/app:1.2.3`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1717",
					Message: "`docker push --disable-content-trust` overrides `DOCKER_CONTENT_TRUST=1` — unsigned image moves into the registry or local store. Sign the artifact (`docker trust sign`) instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1717")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1718(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gh secret set NAME --body-file path`",
			input:    `gh secret set NAME --body-file /run/secrets/foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh secret set NAME --body -` (read stdin)",
			input:    `gh secret set NAME --body -`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh variable set NAME --body val` (non-secret)",
			input:    `gh variable set NAME --body publicvalue`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gh secret set NAME --body SECRET`",
			input: `gh secret set NAME --body hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1718",
					Message: "`gh secret set ... --body hunter2` puts the secret in argv — visible in `ps`, `/proc`, history. Use `--body-file PATH` or pipe via stdin (`... --body -` with `printf %s \"$SECRET\" |`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh secret set NAME --body=SECRET`",
			input: `gh secret set NAME --body=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1718",
					Message: "`gh secret set ... --body=hunter2` puts the secret in argv — visible in `ps`, `/proc`, history. Use `--body-file PATH` or pipe via stdin (`... --body -` with `printf %s \"$SECRET\" |`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh secret set NAME -b SECRET`",
			input: `gh secret set NAME -b hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1718",
					Message: "`gh secret set ... --body hunter2` puts the secret in argv — visible in `ps`, `/proc`, history. Use `--body-file PATH` or pipe via stdin (`... --body -` with `printf %s \"$SECRET\" |`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1718")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1719(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `git filter-repo`",
			input:    `git filter-repo --path secret.txt --invert-paths`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git rebase`",
			input:    `git rebase main`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `git filter-branch --tree-filter`",
			input: `git filter-branch --tree-filter rm secret.txt HEAD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1719",
					Message: "`git filter-branch` is deprecated (Git 2.24+) and its manpage redirects to `git filter-repo`. Use that instead — faster, safer defaults, no orphaned objects.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — bare `git filter-branch`",
			input: `git filter-branch`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1719",
					Message: "`git filter-branch` is deprecated (Git 2.24+) and its manpage redirects to `git filter-repo`. Use that instead — faster, safer defaults, no orphaned objects.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1719")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1720(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — Zsh-idiomatic `$COLUMNS`",
			input:    `print -r -- $COLUMNS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `tput setaf 1` (color, no $COLUMNS equivalent)",
			input:    `tput setaf 1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tput cols`",
			input: `tput cols`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1720",
					Message: "Use `$COLUMNS` instead of `tput cols` — Zsh keeps the terminal size in parameters, no fork needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `tput lines`",
			input: `tput lines`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1720",
					Message: "Use `$LINES` instead of `tput lines` — Zsh keeps the terminal size in parameters, no fork needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1720")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1721(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `chmod 660 /dev/kvm` (group, not world)",
			input:    `chmod 660 /dev/kvm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `chmod 644 /tmp/file` (not /dev/)",
			input:    `chmod 644 /tmp/file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `chmod 644 /dev/null` (read-only, ignored)",
			input:    `chmod 644 /dev/null`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `chmod 666 /dev/kvm`",
			input: `chmod 666 /dev/kvm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1721",
					Message: "`chmod 666 /dev/kvm` opens a kernel device node to every local user — privilege-escalation surface. Use a udev rule that grants the specific group access instead of world-write.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `chmod 0666 /dev/uinput` (leading-zero octal)",
			input: `chmod 0666 /dev/uinput`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1721",
					Message: "`chmod 438 /dev/uinput` opens a kernel device node to every local user — privilege-escalation surface. Use a udev rule that grants the specific group access instead of world-write.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `chmod 777 /dev/dri/card0`",
			input: `chmod 777 /dev/dri/card0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1721",
					Message: "`chmod 777 /dev/dri/card0` opens a kernel device node to every local user — privilege-escalation surface. Use a udev rule that grants the specific group access instead of world-write.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1721")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1722(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ssh-keyscan HOST` (no redirect; fingerprint check separately)",
			input:    `ssh-keyscan HOST`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ssh-keyscan HOST > /tmp/scan.tmp` (not known_hosts)",
			input:    `ssh-keyscan HOST > /tmp/scan.tmp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ssh-keyscan HOST >> ~/.ssh/known_hosts`",
			input: `ssh-keyscan HOST >> ~/.ssh/known_hosts`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1722",
					Message: "`ssh-keyscan ... >> ~/.ssh/known_hosts` accepts the first-served host key without verifying its fingerprint. Pipe to `ssh-keygen -lf -` and assert the fingerprint first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ssh-keyscan -H HOST > known_hosts`",
			input: `ssh-keyscan -H HOST > known_hosts`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1722",
					Message: "`ssh-keyscan ... > known_hosts` accepts the first-served host key without verifying its fingerprint. Pipe to `ssh-keygen -lf -` and assert the fingerprint first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1722")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1723(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gpg --list-keys`",
			input:    `gpg --list-keys`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gpg KEYID --export-secret-keys`",
			input:    `gpg KEYID --export-secret-keys`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gpg --delete-secret-keys KEYID` (leading-flag form)",
			input: `gpg --delete-secret-keys KEYID`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1723",
					Message: "`gpg --delete-secret-keys` permanently destroys keyring entries — no recovery without a separate backup. Export with `gpg --export-secret-keys --armor KEYID` first; never pair this flag with `--batch --yes`.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:  "invalid — `gpg --batch --yes --delete-secret-and-public-keys KEYID`",
			input: `gpg --batch --yes --delete-secret-and-public-keys KEYID`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1723",
					Message: "`gpg --delete-secret-and-public-keys` permanently destroys keyring entries — no recovery without a separate backup. Export with `gpg --export-secret-keys --armor KEYID` first; never pair this flag with `--batch --yes`.",
					Line:    1,
					Column:  20,
				},
			},
		},
		{
			name:  "invalid — `gpg KEYID --delete-key` (trailing-flag form)",
			input: `gpg KEYID --delete-key`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1723",
					Message: "`gpg --delete-key` permanently destroys keyring entries — no recovery without a separate backup. Export with `gpg --export-secret-keys --armor KEYID` first; never pair this flag with `--batch --yes`.",
					Line:    1,
					Column:  12,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1723")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1724(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `pacman -Syu package` (full upgrade then install)",
			input:    `pacman -Syu package`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pacman -S package` (install without DB refresh)",
			input:    `pacman -S package`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pacman -Sy` (refresh DB, no install)",
			input:    `pacman -Sy`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `pacman -Sy package`",
			input: `pacman -Sy package`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1724",
					Message: "`pacman -Sy <pkg>` is a partial-upgrade footgun — refresh the DB but install only one package against the newer metadata. Use `pacman -Syu` first, then `pacman -S <pkg>` (or `pacman -Syu --noconfirm <pkg>` to keep it atomic).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1724")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1725(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cargo publish` (no inline token)",
			input:    `cargo publish`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cargo login --token -` (stdin sentinel)",
			input:    `cargo login --token -`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cargo build --token foo` (not a publish subcmd)",
			input:    `cargo build --token foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `cargo publish --token TOKEN`",
			input: `cargo publish --token cio_abc123`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1725",
					Message: "`cargo publish --token cio_abc123` puts the credential in argv — visible in `ps`, `/proc`, history. Pipe via stdin (`--token -`) or use env vars like `CARGO_REGISTRY_TOKEN` / `NPM_TOKEN`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `cargo login --token=TOKEN`",
			input: `cargo login --token=cio_abc123`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1725",
					Message: "`cargo login --token=cio_abc123` puts the credential in argv — visible in `ps`, `/proc`, history. Pipe via stdin (`--token -`) or use env vars like `CARGO_REGISTRY_TOKEN` / `NPM_TOKEN`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `npm publish --otp 123456`",
			input: `npm publish --otp 123456`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1725",
					Message: "`npm publish --otp 123456` puts the credential in argv — visible in `ps`, `/proc`, history. Pipe via stdin (`--token -`) or use env vars like `CARGO_REGISTRY_TOKEN` / `NPM_TOKEN`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1725")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1726(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gcloud projects delete PROJECT_ID` (no --quiet)",
			input:    `gcloud projects delete PROJECT_ID`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gcloud projects list --quiet` (not delete)",
			input:    `gcloud projects list --quiet`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gcloud projects delete PROJECT_ID --quiet`",
			input: `gcloud projects delete PROJECT_ID --quiet`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1726",
					Message: "`gcloud ... delete --quiet` skips confirmation — a wrong argument wipes the resource (compute disks, secrets, BigQuery tables have no soft-delete). Drop `--quiet` or destroy through a Terraform plan with review.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gcloud sql instances delete INSTANCE -q`",
			input: `gcloud sql instances delete INSTANCE -q`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1726",
					Message: "`gcloud ... delete --quiet` skips confirmation — a wrong argument wipes the resource (compute disks, secrets, BigQuery tables have no soft-delete). Drop `--quiet` or destroy through a Terraform plan with review.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1726")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1727(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `curl URL --proxy http://PROXY:8080` (no creds in URL)",
			input:    `curl URL --proxy http://PROXY:8080`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl URL` (no proxy)",
			input:    `curl URL`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `wget URL` (no proxy creds)",
			input:    `wget URL`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `curl URL --proxy http://USER:PASS@PROXY:8080`",
			input: `curl URL --proxy http://USER:PASS@PROXY:8080`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1727",
					Message: "`curl --proxy http://USER:PASS@PROXY:8080` puts proxy credentials in argv — visible in `ps`, `/proc`, history. Move them into `~/.curlrc` / `~/.netrc` (chmod 600) or `~/.wgetrc`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `curl URL -x http://USER:PASS@PROXY:8080`",
			input: `curl URL -x http://USER:PASS@PROXY:8080`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1727",
					Message: "`curl --proxy http://USER:PASS@PROXY:8080` puts proxy credentials in argv — visible in `ps`, `/proc`, history. Move them into `~/.curlrc` / `~/.netrc` (chmod 600) or `~/.wgetrc`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `wget URL --proxy-password=hunter2`",
			input: `wget URL --proxy-password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1727",
					Message: "`wget --proxy-password=hunter2` puts proxy credentials in argv — visible in `ps`, `/proc`, history. Move them into `~/.curlrc` / `~/.netrc` (chmod 600) or `~/.wgetrc`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1727")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1728(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `pip install --index-url https://pypi.org/simple pkg`",
			input:    `pip install --index-url https://pypi.org/simple pkg`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pip install pkg` (default https index)",
			input:    `pip install pkg`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `pip install --index-url http://internal/simple pkg`",
			input: `pip install --index-url http://internal/simple pkg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1728",
					Message: "`pip install --index-url http://internal/simple` fetches packages over plaintext HTTP — any MITM swaps the wheel for code execution on the host. Use `https://`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pip install -i http://internal/simple pkg`",
			input: `pip install -i http://internal/simple pkg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1728",
					Message: "`pip install --index-url http://internal/simple` fetches packages over plaintext HTTP — any MITM swaps the wheel for code execution on the host. Use `https://`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pip install --extra-index-url=http://internal/simple pkg`",
			input: `pip install --extra-index-url=http://internal/simple pkg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1728",
					Message: "`pip install --index-url --extra-index-url=http://internal/simple` fetches packages over plaintext HTTP — any MITM swaps the wheel for code execution on the host. Use `https://`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1728")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1729(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ip route flush dev eth1`",
			input:    `ip route flush dev eth1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ip route add default via 192.168.1.1`",
			input:    `ip route add default via 192.168.1.1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ip route replace default via 192.168.1.1`",
			input:    `ip route replace default via 192.168.1.1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ip route flush all`",
			input: `ip route flush all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1729",
					Message: "`ip route flush all` removes the default gateway — the SSH session that just ran it loses connectivity. Scope the flush (`flush dev <iface>`) or use `ip route replace default via <gw>` so the new route is in place first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ip route del default`",
			input: `ip route del default`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1729",
					Message: "`ip route del default` removes the default gateway — the SSH session that just ran it loses connectivity. Scope the flush (`flush dev <iface>`) or use `ip route replace default via <gw>` so the new route is in place first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ip -6 route flush all`",
			input: `ip -6 route flush all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1729",
					Message: "`ip route flush all` removes the default gateway — the SSH session that just ran it loses connectivity. Scope the flush (`flush dev <iface>`) or use `ip route replace default via <gw>` so the new route is in place first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1729")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1730(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `brew install foo` (stable release)",
			input:    `brew install foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `brew list --HEAD` (not install/upgrade)",
			input:    `brew list --HEAD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `brew install --HEAD foo`",
			input: `brew install --HEAD foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1730",
					Message: "`brew install --HEAD` builds from upstream HEAD — every run pulls a different commit. Pin to a stable formula release or vendor a private tap with a fixed revision.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `brew reinstall --HEAD foo`",
			input: `brew reinstall --HEAD foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1730",
					Message: "`brew reinstall --HEAD` builds from upstream HEAD — every run pulls a different commit. Pin to a stable formula release or vendor a private tap with a fixed revision.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1730")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1731(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `curl URL --data 'name=John'` (non-secret key)",
			input:    `curl URL --data 'name=John'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl URL --data @secret.txt` (file reference)",
			input:    `curl URL --data @secret.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl URL --data-binary @-` (stdin sentinel)",
			input:    `curl URL --data-binary @-`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `curl URL -d 'password=hunter2'`",
			input: `curl URL -d 'password=hunter2'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1731",
					Message: "`curl -d 'password=hunter2'` puts secret-keyed POST body (`password=…`) in argv — visible in `ps`, `/proc`, history. Read the value from a file with `--data @PATH` or `--data-binary @-` piped from a secrets store.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `curl URL --data-urlencode 'token=ABC123'`",
			input: `curl URL --data-urlencode 'token=ABC123'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1731",
					Message: "`curl --data-urlencode 'token=ABC123'` puts secret-keyed POST body (`token=…`) in argv — visible in `ps`, `/proc`, history. Read the value from a file with `--data @PATH` or `--data-binary @-` piped from a secrets store.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `wget URL --post-data='api_key=ABC123'`",
			input: `wget URL --post-data='api_key=ABC123'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1731",
					Message: "`wget --post-data='api_key=ABC123'` puts secret-keyed POST body (`api_key=…`) in argv — visible in `ps`, `/proc`, history. Read the value from a file with `--data @PATH` or `--data-binary @-` piped from a secrets store.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1731")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1732(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `flatpak install --user org.example.App`",
			input:    `flatpak install --user org.example.App`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `flatpak override --filesystem=~/Documents org.example.App`",
			input:    `flatpak override --filesystem=~/Documents org.example.App`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `flatpak override --filesystem=host org.example.App`",
			input: `flatpak override --filesystem=host org.example.App`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1732",
					Message: "`flatpak override --filesystem=host` removes the Flatpak sandbox — the app gets unrestricted host-filesystem access. Grant a specific subdirectory (e.g. `--filesystem=~/Documents:ro`) or use Filesystem portals.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `flatpak override --filesystem=home org.example.App`",
			input: `flatpak override --filesystem=home org.example.App`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1732",
					Message: "`flatpak override --filesystem=home` removes the Flatpak sandbox — the app gets unrestricted host-filesystem access. Grant a specific subdirectory (e.g. `--filesystem=~/Documents:ro`) or use Filesystem portals.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `flatpak run --filesystem=host org.example.App`",
			input: `flatpak run --filesystem=host org.example.App`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1732",
					Message: "`flatpak run --filesystem=host` removes the Flatpak sandbox — the app gets unrestricted host-filesystem access. Grant a specific subdirectory (e.g. `--filesystem=~/Documents:ro`) or use Filesystem portals.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1732")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1733(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `docker plugin install vieux/sshfs` (interactive prompt kept)",
			input:    `docker plugin install vieux/sshfs`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `docker plugin ls --grant-all-permissions` (not install)",
			input:    `docker plugin ls --grant-all-permissions`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `docker plugin install --grant-all-permissions vieux/sshfs`",
			input: `docker plugin install --grant-all-permissions vieux/sshfs`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1733",
					Message: "`docker plugin install --grant-all-permissions` accepts every capability the plugin requests — root-equivalent on the host. Walk the interactive prompt manually and pin the digest once vetted.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `docker plugin install vieux/sshfs --grant-all-permissions`",
			input: `docker plugin install vieux/sshfs --grant-all-permissions`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1733",
					Message: "`docker plugin install --grant-all-permissions` accepts every capability the plugin requests — root-equivalent on the host. Walk the interactive prompt manually and pin the digest once vetted.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1733")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1734(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `useradd alice`",
			input:    `useradd alice`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cat /etc/passwd` (read-only)",
			input:    `cat /etc/passwd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `cp /tmp/passwd /etc/passwd`",
			input: `cp /tmp/passwd /etc/passwd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1734",
					Message: "`cp /etc/passwd` bypasses the lock that `vipw` / `vigr` / `useradd` use on the user-identity files. Edit through `vipw -e` / `vigr -e`, or mutate a single entry with `useradd` / `usermod` / `passwd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `tee /etc/shadow`",
			input: `tee /etc/shadow`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1734",
					Message: "`tee /etc/shadow` bypasses the lock that `vipw` / `vigr` / `useradd` use on the user-identity files. Edit through `vipw -e` / `vigr -e`, or mutate a single entry with `useradd` / `usermod` / `passwd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `echo entry >> /etc/group`",
			input: `echo entry >> /etc/group`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1734",
					Message: "`>> /etc/group` bypasses the lock that `vipw` / `vigr` / `useradd` use on the user-identity files. Edit through `vipw -e` / `vigr -e`, or mutate a single entry with `useradd` / `usermod` / `passwd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1734")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1735(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `efibootmgr -v` (inspect)",
			input:    `efibootmgr -v`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `efibootmgr -o 0001,0002` (reorder)",
			input:    `efibootmgr -o 0001,0002`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `efibootmgr -B`",
			input: `efibootmgr -B`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1735",
					Message: "`efibootmgr -B` deletes a UEFI boot entry — wrong BOOTNUM (or missing fallback) leaves the box at the UEFI shell on next reboot. Inspect `efibootmgr -v` first; demote via `-o NEW,ORDER` instead of deleting.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `efibootmgr -B -b 0001`",
			input: `efibootmgr -B -b 0001`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1735",
					Message: "`efibootmgr -B` deletes a UEFI boot entry — wrong BOOTNUM (or missing fallback) leaves the box at the UEFI shell on next reboot. Inspect `efibootmgr -v` first; demote via `-o NEW,ORDER` instead of deleting.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1735")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1736(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `pulumi preview`",
			input:    `pulumi preview`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pulumi up` (interactive prompt kept)",
			input:    `pulumi up`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pulumi stack ls --yes` (not up/destroy/refresh)",
			input:    `pulumi stack ls --yes`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `pulumi destroy --yes`",
			input: `pulumi destroy --yes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1736",
					Message: "`pulumi destroy --yes` skips the preview-and-confirm — a misresolved stack or credential wipes / mutates infrastructure with no review. Gate behind `pulumi preview` plus a manual approval step.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pulumi up -y`",
			input: `pulumi up -y`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1736",
					Message: "`pulumi up -y` skips the preview-and-confirm — a misresolved stack or credential wipes / mutates infrastructure with no review. Gate behind `pulumi preview` plus a manual approval step.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pulumi refresh --yes`",
			input: `pulumi refresh --yes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1736",
					Message: "`pulumi refresh --yes` skips the preview-and-confirm — a misresolved stack or credential wipes / mutates infrastructure with no review. Gate behind `pulumi preview` plus a manual approval step.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1736")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1737(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `wpa_passphrase MySSID` (passphrase via stdin)",
			input:    `wpa_passphrase MySSID`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `wpa_passphrase MySSID < /run/secrets/wifi`",
			input:    `wpa_passphrase MySSID < /run/secrets/wifi`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `wpa_passphrase MySSID hunter2`",
			input: `wpa_passphrase MySSID hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1737",
					Message: "`wpa_passphrase SSID PASSWORD` puts the Wi-Fi passphrase in argv — visible in `ps`, `/proc`, history. Drop the PASSWORD argument and pipe it via stdin (`wpa_passphrase SSID < /run/secrets/wifi`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1737")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1738(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `aws rds delete-db-instance` with explicit final snapshot",
			input:    `aws rds delete-db-instance --db-instance-identifier mydb --final-db-snapshot-identifier mydb-final`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `aws rds describe-db-instances`",
			input:    `aws rds describe-db-instances`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `aws rds delete-db-instance ... --skip-final-snapshot`",
			input: `aws rds delete-db-instance --db-instance-identifier mydb --skip-final-snapshot`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1738",
					Message: "`aws rds delete-db-instance --skip-final-snapshot` deletes the database with no recovery snapshot. Drop the flag or pass `--final-db-snapshot-identifier <name>` so the snapshot is explicit and verifiable.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `aws rds delete-db-cluster ... --skip-final-snapshot`",
			input: `aws rds delete-db-cluster --db-cluster-identifier mycluster --skip-final-snapshot`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1738",
					Message: "`aws rds delete-db-cluster --skip-final-snapshot` deletes the database with no recovery snapshot. Drop the flag or pass `--final-db-snapshot-identifier <name>` so the snapshot is explicit and verifiable.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1738")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1739(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `git submodule update --init --recursive` (pinned commits)",
			input:    `git submodule update --init --recursive`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git submodule add URL path`",
			input:    `git submodule add https://example.com/repo path`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `git submodule update --remote`",
			input: `git submodule update --remote`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1739",
					Message: "`git submodule update --remote` ignores the pinned commits in the parent repo and pulls each submodule's branch HEAD — non-reproducible builds, supply-chain risk. Use `--init --recursive` and bump pins via reviewed PRs.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `git submodule update --remote --recursive`",
			input: `git submodule update --remote --recursive`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1739",
					Message: "`git submodule update --remote` ignores the pinned commits in the parent repo and pulls each submodule's branch HEAD — non-reproducible builds, supply-chain risk. Use `--init --recursive` and bump pins via reviewed PRs.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1739")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1740(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gh release upload v1.0 file.tar.gz`",
			input:    `gh release upload v1.0 file.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh release create v1.0 file.tar.gz`",
			input:    `gh release create v1.0 file.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gh release upload v1.0 file.tar.gz --clobber`",
			input: `gh release upload v1.0 file.tar.gz --clobber`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1740",
					Message: "`gh release upload --clobber` silently replaces an existing asset — a re-run can downgrade the user-facing download. Drop `--clobber` or version the asset name so each upload has a unique slot.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1740")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1741(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mkpasswd -s` (read from stdin)",
			input:    `mkpasswd -s`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mkpasswd -m sha-512 -s`",
			input:    `mkpasswd -m sha-512 -s`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mkpasswd --stdin`",
			input:    `mkpasswd --stdin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mkpasswd hunter2`",
			input: `mkpasswd hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1741",
					Message: "`mkpasswd PASSWORD` puts the cleartext password in argv — visible in `ps`, `/proc`, history. Use `mkpasswd -s` and pipe the secret via stdin (`printf %s \"$PASSWORD\" | mkpasswd -s`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `mkpasswd -m sha-512 hunter2`",
			input: `mkpasswd -m sha-512 hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1741",
					Message: "`mkpasswd PASSWORD` puts the cleartext password in argv — visible in `ps`, `/proc`, history. Use `mkpasswd -s` and pipe the secret via stdin (`printf %s \"$PASSWORD\" | mkpasswd -s`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1741")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1742(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mc alias set NAME URL` (interactive prompt)",
			input:    `mc alias set myminio https://play.min.io`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mc alias list`",
			input:    `mc alias list`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mc alias set NAME URL ACCESS SECRET`",
			input: `mc alias set myminio https://play.min.io ACCESSKEY SECRETKEY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1742",
					Message: "`mc alias set ... ACCESS_KEY SECRET_KEY` puts S3 access and secret keys in argv — visible in `ps`, `/proc`, history. Drop the trailing keys (mc prompts) or use `MC_HOST_<alias>=URL` env-var form scoped to one command.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `mc config host add NAME URL ACCESS SECRET` (legacy)",
			input: `mc config host add myminio https://play.min.io ACCESSKEY SECRETKEY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1742",
					Message: "`mc config host add ... ACCESS SECRET` puts S3 access and secret keys in argv — visible in `ps`, `/proc`, history. Drop the trailing keys (mc prompts) or use `MC_HOST_<alias>=URL` env-var form scoped to one command.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1742")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1743(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `npm audit fix` (no --force)",
			input:    `npm audit fix`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `npm audit` (no fix)",
			input:    `npm audit`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `npm run build --force` (not audit fix)",
			input:    `npm run build --force`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `npm audit fix --force`",
			input: `npm audit fix --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1743",
					Message: "`npm audit ... --force` accepts every major-version bump an advisory triggers — silent breaking changes. Drop `--force` and triage advisories one by one.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pnpm audit --fix --force`",
			input: `pnpm audit --fix --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1743",
					Message: "`pnpm audit ... --force` accepts every major-version bump an advisory triggers — silent breaking changes. Drop `--force` and triage advisories one by one.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1743")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1744(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl port-forward pod/mypod 8080:8080` (loopback default)",
			input:    `kubectl port-forward pod/mypod 8080:8080`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl port-forward --address 127.0.0.1 pod/mypod 8080:8080`",
			input:    `kubectl port-forward --address 127.0.0.1 pod/mypod 8080:8080`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl port-forward pod/mypod 8080:8080 --address 0.0.0.0`",
			input: `kubectl port-forward pod/mypod 8080:8080 --address 0.0.0.0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1744",
					Message: "`kubectl port-forward --address 0.0.0.0` binds the local end of the tunnel on every interface — anyone on the LAN / VPN can reach the pod. Drop `--address` (loopback default) or pick a trusted-network interface IP.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `kubectl port-forward pod/mypod 8080:8080 --address=0.0.0.0`",
			input: `kubectl port-forward pod/mypod 8080:8080 --address=0.0.0.0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1744",
					Message: "`kubectl port-forward --address=0.0.0.0` binds the local end of the tunnel on every interface — anyone on the LAN / VPN can reach the pod. Drop `--address` (loopback default) or pick a trusted-network interface IP.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1744")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1745(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `poetry publish --repository myrepo` (no password)",
			input:    `poetry publish --repository myrepo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `twine upload dist/*` (token via env)",
			input:    `twine upload dist/*`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `poetry publish --username u --password hunter2`",
			input: `poetry publish --username u --password hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1745",
					Message: "`poetry publish --password hunter2` puts the registry password in argv — visible in `ps`, `/proc`, history. Use env vars (`POETRY_PYPI_TOKEN_<NAME>`, `TWINE_PASSWORD`) or a `0600` `~/.pypirc` file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `twine upload -u u -p hunter2 dist/*`",
			input: `twine upload -u u -p hunter2 dist/*`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1745",
					Message: "`twine upload --password hunter2` puts the registry password in argv — visible in `ps`, `/proc`, history. Use env vars (`POETRY_PYPI_TOKEN_<NAME>`, `TWINE_PASSWORD`) or a `0600` `~/.pypirc` file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `poetry publish --password=hunter2`",
			input: `poetry publish --password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1745",
					Message: "`poetry publish --password=hunter2` puts the registry password in argv — visible in `ps`, `/proc`, history. Use env vars (`POETRY_PYPI_TOKEN_<NAME>`, `TWINE_PASSWORD`) or a `0600` `~/.pypirc` file.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1745")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1746(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sysctl -w kernel.randomize_va_space=2` (default ASLR)",
			input:    `sysctl -w kernel.randomize_va_space=2`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sysctl kernel.randomize_va_space`",
			input:    `sysctl kernel.randomize_va_space`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `sysctl -w kernel.randomize_va_space=0`",
			input: `sysctl -w kernel.randomize_va_space=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1746",
					Message: "`sysctl kernel.randomize_va_space=0` weakens ASLR — absolute-address exploits become deterministic (stack overflows, ROP). Keep `kernel.randomize_va_space=2` outside a sandboxed debug context.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `sysctl kernel.randomize_va_space=1`",
			input: `sysctl kernel.randomize_va_space=1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1746",
					Message: "`sysctl kernel.randomize_va_space=1` weakens ASLR — absolute-address exploits become deterministic (stack overflows, ROP). Keep `kernel.randomize_va_space=2` outside a sandboxed debug context.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1746")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1747(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `npm install --registry=https://registry.npmjs.org/`",
			input:    `npm install --registry=https://registry.npmjs.org/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `yarn config set registry https://registry.npmjs.org/`",
			input:    `yarn config set registry https://registry.npmjs.org/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `npm install --registry=http://internal/`",
			input: `npm install --registry=http://internal/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1747",
					Message: "`npm --registry=http://internal/` uses plaintext HTTP for the package registry — any MITM swaps tarballs and runs install-time `postinstall` code. Use `https://`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pnpm install --registry http://internal/`",
			input: `pnpm install --registry http://internal/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1747",
					Message: "`pnpm --registry http://internal/` uses plaintext HTTP for the package registry — any MITM swaps tarballs and runs install-time `postinstall` code. Use `https://`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `yarn config set registry http://internal/`",
			input: `yarn config set registry http://internal/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1747",
					Message: "`yarn config set registry http://internal/` uses plaintext HTTP for the package registry — any MITM swaps tarballs and runs install-time `postinstall` code. Use `https://`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1747")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1748(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `helm repo add bitnami https://charts.bitnami.com/bitnami`",
			input:    `helm repo add bitnami https://charts.bitnami.com/bitnami`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `helm repo update`",
			input:    `helm repo update`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `helm repo add myrepo http://internal/charts`",
			input: `helm repo add myrepo http://internal/charts`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1748",
					Message: "`helm repo add myrepo http://internal/charts` fetches charts over plaintext HTTP — any MITM swaps the chart and its referenced images. Use `https://` and `helm install --verify` for provenance.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1748")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1749(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `virsh undefine mydomain` (config only)",
			input:    `virsh undefine mydomain`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `virsh list --all`",
			input:    `virsh list --all`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `virsh undefine mydomain --remove-all-storage`",
			input: `virsh undefine mydomain --remove-all-storage`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1749",
					Message: "`virsh undefine ... --remove-all-storage` deletes every disk image the domain references — no soft-delete, no recycle bin. Back up first (`qemu-img convert`), `undefine` without the flag, then `virsh vol-delete` deliberately.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `virsh undefine mydomain --wipe-storage`",
			input: `virsh undefine mydomain --wipe-storage`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1749",
					Message: "`virsh undefine ... --wipe-storage` deletes every disk image the domain references — no soft-delete, no recycle bin. Back up first (`qemu-img convert`), `undefine` without the flag, then `virsh vol-delete` deliberately.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1749")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1750(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl proxy --port 8001` (loopback default)",
			input:    `kubectl proxy --port 8001`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl proxy --address 127.0.0.1`",
			input:    `kubectl proxy --address 127.0.0.1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl proxy --address 0.0.0.0`",
			input: `kubectl proxy --address 0.0.0.0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1750",
					Message: "`kubectl proxy --address 0.0.0.0` exposes the cluster-admin API tunnel to every reachable interface. Keep the loopback default and tunnel over SSH, or restrict `--address` to a firewalled interface.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `kubectl proxy --address=0.0.0.0`",
			input: `kubectl proxy --address=0.0.0.0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1750",
					Message: "`kubectl proxy --address=0.0.0.0` exposes the cluster-admin API tunnel to every reachable interface. Keep the loopback default and tunnel over SSH, or restrict `--address` to a firewalled interface.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1750")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1751(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `rpm -e libfoo`",
			input:    `rpm -e libfoo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `dnf remove libfoo`",
			input:    `dnf remove libfoo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rpm -q libfoo --nodeps` (query, not erase)",
			input:    `rpm -q libfoo --nodeps`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `rpm -e --nodeps libfoo`",
			input: `rpm -e --nodeps libfoo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1751",
					Message: "`rpm ... --nodeps` removes the package without the dependency solver — dependents break (libc, openssl, systemd units). Resolve the conflict explicitly instead of bypassing.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `dnf remove --nodeps libfoo`",
			input: `dnf remove --nodeps libfoo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1751",
					Message: "`dnf ... --nodeps` removes the package without the dependency solver — dependents break (libc, openssl, systemd units). Resolve the conflict explicitly instead of bypassing.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1751")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1752(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `pvcreate $DISK` (prompts kept)",
			input:    `pvcreate $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `vgcreate my_vg $DISK`",
			input:    `vgcreate my_vg $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `pvcreate -ff $DISK`",
			input: `pvcreate -ff $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1752",
					Message: "`pvcreate -ff` skips the LVM confirmation — a wrong device gets its filesystem / RAID / LVM signatures wiped. Inspect with `wipefs -n` + `lsblk -f` first, drop the flag, re-run after checking the target.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pvcreate $DISK --yes`",
			input: `pvcreate $DISK --yes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1752",
					Message: "`pvcreate --yes` skips the LVM confirmation — a wrong device gets its filesystem / RAID / LVM signatures wiped. Inspect with `wipefs -n` + `lsblk -f` first, drop the flag, re-run after checking the target.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `vgcreate -y my_vg $DISK`",
			input: `vgcreate -y my_vg $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1752",
					Message: "`vgcreate -y` skips the LVM confirmation — a wrong device gets its filesystem / RAID / LVM signatures wiped. Inspect with `wipefs -n` + `lsblk -f` first, drop the flag, re-run after checking the target.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1752")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1753(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `rclone delete myremote:bucket/path`",
			input:    `rclone delete myremote:bucket/path`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rclone sync src dst`",
			input:    `rclone sync src dst`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `rclone purge myremote:bucket/path`",
			input: `rclone purge myremote:bucket/path`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1753",
					Message: "`rclone purge` removes every object under the remote path with no dry-run or soft-delete. Preview with `rclone lsf` / `rclone delete --dry-run` and prefer narrower `rclone delete`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1753")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1754(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gh auth status`",
			input:    `gh auth status`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh auth token`",
			input:    `gh auth token`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gh auth status -t`",
			input: `gh auth status -t`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1754",
					Message: "`gh auth status -t` prints the OAuth token in the status output — CI logs and scrollback become a repo-admin leak. Use `gh auth token` in automation so the secret path is explicit.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh auth status --show-token`",
			input: `gh auth status --show-token`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1754",
					Message: "`gh auth status --show-token` prints the OAuth token in the status output — CI logs and scrollback become a repo-admin leak. Use `gh auth token` in automation so the secret path is explicit.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1754")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1755(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gcloud sql users create ... --prompt-for-password`",
			input:    `gcloud sql users create myuser --instance myinst --prompt-for-password`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gcloud sql users list --instance myinst`",
			input:    `gcloud sql users list --instance myinst`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gcloud sql users create ... --password PASS`",
			input: `gcloud sql users create myuser --instance myinst --password hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1755",
					Message: "`gcloud sql users create --password hunter2` puts the Cloud SQL password in argv — visible in `ps`, `/proc`, history, and Cloud Audit Logs. Use `--prompt-for-password` or call the SQL Admin API with a body file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gcloud sql users set-password ... --password=PASS`",
			input: `gcloud sql users set-password myuser --instance myinst --password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1755",
					Message: "`gcloud sql users set-password --password=hunter2` puts the Cloud SQL password in argv — visible in `ps`, `/proc`, history, and Cloud Audit Logs. Use `--prompt-for-password` or call the SQL Admin API with a body file.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1755")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1756(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `chmod 660 /var/run/docker.sock` (group only)",
			input:    `chmod 660 /var/run/docker.sock`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `chmod 666 /tmp/file` (not a runtime socket)",
			input:    `chmod 666 /tmp/file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `chmod 666 /var/run/docker.sock`",
			input: `chmod 666 /var/run/docker.sock`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1756",
					Message: "`chmod 666 /var/run/docker.sock` grants every local user access to a root-equivalent container-runtime socket. Keep `0660` owned by the runtime group (`root:docker` etc.) and restrict membership.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `chmod 777 /run/containerd/containerd.sock`",
			input: `chmod 777 /run/containerd/containerd.sock`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1756",
					Message: "`chmod 777 /run/containerd/containerd.sock` grants every local user access to a root-equivalent container-runtime socket. Keep `0660` owned by the runtime group (`root:docker` etc.) and restrict membership.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1756")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1757(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gh auth refresh --scopes repo`",
			input:    `gh auth refresh --scopes repo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh auth login --scopes workflow,read:org`",
			input:    `gh auth login --scopes workflow,read:org`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gh auth refresh --scopes delete_repo`",
			input: `gh auth refresh --scopes delete_repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1757",
					Message: "`gh auth refresh --scopes delete_repo` escalates the token to destructive privileges that outlast the script. Request the minimum scope (`repo`, `workflow`) and rotate the token when done.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh auth refresh -s admin:org`",
			input: `gh auth refresh -s admin:org`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1757",
					Message: "`gh auth refresh --scopes admin:org` escalates the token to destructive privileges that outlast the script. Request the minimum scope (`repo`, `workflow`) and rotate the token when done.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh auth login --scopes=repo,delete_repo`",
			input: `gh auth login --scopes=repo,delete_repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1757",
					Message: "`gh auth login --scopes delete_repo` escalates the token to destructive privileges that outlast the script. Request the minimum scope (`repo`, `workflow`) and rotate the token when done.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1757")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1758(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gh codespace delete -c mycodespace` (prompt kept)",
			input:    `gh codespace delete -c mycodespace`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh codespace list`",
			input:    `gh codespace list`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gh codespace delete --force`",
			input: `gh codespace delete --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1758",
					Message: "`gh codespace delete --force` skips the prompt and drops uncommitted work along with the codespace. Let the prompt list what's about to go and verify `git status` inside first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh codespace delete -f --all`",
			input: `gh codespace delete -f --all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1758",
					Message: "`gh codespace delete -f` skips the prompt and drops uncommitted work along with the codespace. Let the prompt list what's about to go and verify `git status` inside first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1758")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1759(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `vault login -` (reads token from stdin)",
			input:    `vault login -`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `vault login -method=userpass username=alice` (no secret key)",
			input:    `vault login -method=userpass username=alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `vault login mytoken` (positional token)",
			input: `vault login mytoken`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1759",
					Message: "`vault login mytoken` puts the Vault credential in argv — visible in `ps`, `/proc`, history, Vault audit log. Use `vault login -` with stdin or source `VAULT_TOKEN` from a secrets file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `vault login -method=userpass username=alice password=hunter2`",
			input: `vault login -method=userpass username=alice password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1759",
					Message: "`vault login password=hunter2` puts the Vault credential in argv — visible in `ps`, `/proc`, history, Vault audit log. Use `vault login -` with stdin or source `VAULT_TOKEN` from a secrets file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `vault auth mytoken`",
			input: `vault auth mytoken`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1759",
					Message: "`vault auth mytoken` puts the Vault credential in argv — visible in `ps`, `/proc`, history, Vault audit log. Use `vault login -` with stdin or source `VAULT_TOKEN` from a secrets file.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1759")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1760(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `openssl rand -hex 32`",
			input:    `openssl rand -hex 32`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `openssl rand -hex 16` (borderline but accepted)",
			input:    `openssl rand -hex 16`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `openssl rand 24` (no encoding flag)",
			input:    `openssl rand 24`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `openssl rand -hex 8`",
			input: `openssl rand -hex 8`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1760",
					Message: "`openssl rand -hex 8` produces a sub-128-bit value — brute-forceable offline. Use `-hex 32` for secrets / long-lived tokens.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `openssl rand -base64 12`",
			input: `openssl rand -base64 12`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1760",
					Message: "`openssl rand -base64 12` produces a sub-128-bit value — brute-forceable offline. Use `-hex 32` for secrets / long-lived tokens.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1760")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1761(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gh gist create secret.env` (secret by default)",
			input:    `gh gist create secret.env`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh gist list`",
			input:    `gh gist list`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gh gist create --public secret.env`",
			input: `gh gist create --public secret.env`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1761",
					Message: "`gh gist create --public` publishes the file to the public discover feed — search engines crawl it within minutes. Drop the flag unless public exposure is the explicit goal.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh gist create -p note.md`",
			input: `gh gist create -p note.md`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1761",
					Message: "`gh gist create -p` publishes the file to the public discover feed — search engines crawl it within minutes. Drop the flag unless public exposure is the explicit goal.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1761")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1762(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubeadm join ... --discovery-token-ca-cert-hash sha256:xxx`",
			input:    `kubeadm join 10.0.0.1:6443 --token abc --discovery-token-ca-cert-hash sha256:xxx`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubeadm token create`",
			input:    `kubeadm token create --print-join-command`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubeadm join ... --discovery-token-unsafe-skip-ca-verification`",
			input: `kubeadm join 10.0.0.1:6443 --token abc --discovery-token-unsafe-skip-ca-verification`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1762",
					Message: "`kubeadm join --discovery-token-unsafe-skip-ca-verification` skips CA verification of the control-plane — MITM steals the bootstrap token. Pin the CA with `--discovery-token-ca-cert-hash sha256:<digest>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1762")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1763(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `docker compose down`",
			input:    `docker compose down`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `docker compose up -d`",
			input:    `docker compose up -d`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `docker compose down -v`",
			input: `docker compose down -v`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1763",
					Message: "`docker compose down -v` wipes every named volume declared in the stack — database, cache, uploaded assets go with it. Drop the flag in CI / prod scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `docker compose down --volumes`",
			input: `docker compose down --volumes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1763",
					Message: "`docker compose down --volumes` wipes every named volume declared in the stack — database, cache, uploaded assets go with it. Drop the flag in CI / prod scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `docker-compose down -v` (hyphen form)",
			input: `docker-compose down -v`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1763",
					Message: "`docker compose down -v` wipes every named volume declared in the stack — database, cache, uploaded assets go with it. Drop the flag in CI / prod scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1763")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1764(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `git commit -m \"msg\"`",
			input:    `git commit -m "msg"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git commit -S -m \"msg\"` (signed)",
			input:    `git commit -S -m "msg"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `git commit --no-verify -m \"msg\"`",
			input: `git commit --no-verify -m "msg"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1764",
					Message: "`git commit --no-verify` skips pre-commit and commit-msg hooks — the last guardrail against secret leaks and broken tests. Fix the hook or carve a narrow exemption instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `git commit -n -m \"msg\"`",
			input: `git commit -n -m "msg"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1764",
					Message: "`git commit -n` skips pre-commit and commit-msg hooks — the last guardrail against secret leaks and broken tests. Fix the hook or carve a narrow exemption instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1764")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1765(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `snap remove mysnap` (snapshot kept)",
			input:    `snap remove mysnap`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `snap list`",
			input:    `snap list`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `snap remove --purge mysnap`",
			input: `snap remove --purge mysnap`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1765",
					Message: "`snap remove --purge` skips the pre-remove data snapshot — the snap's files are gone with no rollback. Drop `--purge` or capture a `snap save` set ID first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1765")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1766(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `memcached -l 127.0.0.1`",
			input:    `memcached -l 127.0.0.1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `memcached -l 10.0.0.5 -p 11211`",
			input:    `memcached -l 10.0.0.5 -p 11211`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `memcached -l 0.0.0.0`",
			input: `memcached -l 0.0.0.0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1766",
					Message: "`memcached -l 0.0.0.0` exposes the unauthenticated cache to every interface on the host. Bind to `127.0.0.1` or a private-network IP and firewall the port.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `memcached -l0.0.0.0` (joined form)",
			input: `memcached -l0.0.0.0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1766",
					Message: "`memcached -l0.0.0.0` exposes the unauthenticated cache to every interface on the host. Bind to `127.0.0.1` or a private-network IP and firewall the port.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1766")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1767(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mongod --bind_ip 127.0.0.1`",
			input:    `mongod --bind_ip 127.0.0.1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mongod --bind_ip 0.0.0.0`",
			input: `mongod --bind_ip 0.0.0.0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1767",
					Message: "`mongod --bind_ip 0.0.0.0` exposes MongoDB on every interface — 2017 ransomware-wave target. Bind to `127.0.0.1` or a private-network IP, enable `--auth`, firewall port 27017.",
					Line:    1,
					Column:  9,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1767")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1768(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sqlcmd -S server -U user -P` (prompt)",
			input:    `sqlcmd -S server -U user -P`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sqlcmd -S server -E` (Windows auth, no password)",
			input:    `sqlcmd -S server -E`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `sqlcmd -S server -U user -P hunter2`",
			input: `sqlcmd -S server -U user -P hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1768",
					Message: "`sqlcmd -P hunter2` puts the SQL Server password in argv — visible in `ps`, `/proc`, history, SQL Server audit. Use `-P` with no arg (prompt) or `SQLCMDPASSWORD` env var.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `bcp mydb in data.csv -U user -P hunter2`",
			input: `bcp mydb in data.csv -U user -P hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1768",
					Message: "`bcp -P hunter2` puts the SQL Server password in argv — visible in `ps`, `/proc`, history, SQL Server audit. Use `-P` with no arg (prompt) or `SQLCMDPASSWORD` env var.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1768")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1769(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `vagrant destroy` (prompt kept)",
			input:    `vagrant destroy`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `vagrant halt`",
			input:    `vagrant halt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `vagrant destroy --force`",
			input: `vagrant destroy --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1769",
					Message: "`vagrant destroy --force` skips the prompt and drops the VM (and any un-exported data inside). Drop the flag, or use `vagrant halt` + `vagrant up` to cycle without destroy.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `vagrant destroy -f myvm`",
			input: `vagrant destroy -f myvm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1769",
					Message: "`vagrant destroy -f` skips the prompt and drops the VM (and any un-exported data inside). Drop the flag, or use `vagrant halt` + `vagrant up` to cycle without destroy.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1769")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1770(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gpg --verify sig.asc` (default trust model)",
			input:    `gpg --verify sig.asc`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gpg --trust-model pgp --verify sig.asc`",
			input:    `gpg --trust-model pgp --verify sig.asc`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gpg --verify --always-trust sig.asc` (trailing form)",
			input: `gpg --verify --always-trust sig.asc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1770",
					Message: "`gpg --always-trust` marks every imported key as fully trusted — a signature from an attacker-supplied key verifies cleanly. Drop the flag and pin the expected fingerprint, or assign trust via `gpg --edit-key KEYID trust`.",
					Line:    1,
					Column:  15,
				},
			},
		},
		{
			name:  "invalid — `gpg --verify --trust-model always sig.asc`",
			input: `gpg --verify --trust-model always sig.asc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1770",
					Message: "`gpg --trust-model always` marks every imported key as fully trusted — a signature from an attacker-supplied key verifies cleanly. Drop the flag and pin the expected fingerprint, or assign trust via `gpg --edit-key KEYID trust`.",
					Line:    1,
					Column:  15,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1770")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1771(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `alias ll='ls -l'` (regular alias)",
			input:    `alias ll='ls -l'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `alias` (no args, lists aliases)",
			input:    `alias`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `alias -g G='| grep'`",
			input: `alias -g G='| grep'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1771",
					Message: "`alias -g` defines a global alias that expands outside command position — a surprise for anyone reading the script later. Prefer a function, or keep global aliases in `~/.zshrc` where they are discoverable.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `alias -s log=less`",
			input: `alias -s log=less`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1771",
					Message: "`alias -s` defines a suffix alias that expands outside command position — a surprise for anyone reading the script later. Prefer a function, or keep suffix aliases in `~/.zshrc` where they are discoverable.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1771")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1772(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `hdparm -I $DISK` (info only)",
			input:    `hdparm -I $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `hdparm -tT $DISK` (benchmark)",
			input:    `hdparm -tT $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `hdparm --security-erase PASS $DISK`",
			input: `hdparm --security-erase PASS $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1772",
					Message: "`hdparm --security-erase` issues an ATA-level operation that ignores filesystems and cannot be rolled back. Pin the disk by `/dev/disk/by-id/…`, keep it behind a runbook, keep the password out of argv.",
					Line:    1,
					Column:  9,
				},
			},
		},
		{
			name:  "invalid — `hdparm --trim-sector-ranges 0:1 $DISK`",
			input: `hdparm --trim-sector-ranges 0:1 $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1772",
					Message: "`hdparm --trim-sector-ranges` issues an ATA-level operation that ignores filesystems and cannot be rolled back. Pin the disk by `/dev/disk/by-id/…`, keep it behind a runbook, keep the password out of argv.",
					Line:    1,
					Column:  9,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1772")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1773(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `xargs -r rm` (guard present)",
			input:    `xargs -r rm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `xargs --no-run-if-empty rm`",
			input:    `xargs --no-run-if-empty rm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `xargs -0r rm` (combined short flags)",
			input:    `xargs -0r rm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `xargs` alone (no command, probably a noop / listing)",
			input:    `xargs`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `xargs rm` (no guard)",
			input: `xargs rm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1773",
					Message: "`xargs` without `-r` / `--no-run-if-empty` runs the child once with no arguments when stdin is empty — a destructive surprise for `xargs rm`, `xargs kill`, etc. Add `-r` or switch to `find ... -exec cmd {} +`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `xargs -0 kill -9`",
			input: `xargs -0 kill -9`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1773",
					Message: "`xargs` without `-r` / `--no-run-if-empty` runs the child once with no arguments when stdin is empty — a destructive surprise for `xargs rm`, `xargs kill`, etc. Add `-r` or switch to `find ... -exec cmd {} +`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1773")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1774(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt EXTENDED_GLOB` (unrelated option)",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `set -e` (unrelated short option)",
			input:    `set -e`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt GLOB_SUBST`",
			input: `setopt GLOB_SUBST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1774",
					Message: "`setopt GLOB_SUBST` enables `GLOB_SUBST` — every unquoted `$var` expansion is rescanned as a glob pattern. User-controlled data becomes a filesystem query. Scope this in a subshell / function, or use explicit expansion flags.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt globsubst` (Zsh lower/underscore folded)",
			input: `setopt globsubst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1774",
					Message: "`setopt globsubst` enables `GLOB_SUBST` — every unquoted `$var` expansion is rescanned as a glob pattern. User-controlled data becomes a filesystem query. Scope this in a subshell / function, or use explicit expansion flags.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `set -o GLOB_SUBST`",
			input: `set -o GLOB_SUBST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1774",
					Message: "`set -o GLOB_SUBST` enables `GLOB_SUBST` — every unquoted `$var` expansion is rescanned as a glob pattern. User-controlled data becomes a filesystem query. Scope this in a subshell / function, or use explicit expansion flags.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1774")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1775(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `timeout -k 5 30 cmd` (escalation configured)",
			input:    `timeout -k 5 30 cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `timeout --kill-after=5 30 cmd`",
			input:    `timeout --kill-after=5 30 cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `timeout` alone (no command, probably listing help)",
			input:    `timeout`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `timeout 30 cmd` (no escalation)",
			input: `timeout 30 cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1775",
					Message: "`timeout` without `--kill-after` / `-k` only sends `SIGTERM` — a child that blocks or ignores it hangs the pipeline past the deadline. Add `--kill-after=N` so timeout escalates to `SIGKILL`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `timeout 5m long-running-job`",
			input: `timeout 5m long-running-job`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1775",
					Message: "`timeout` without `--kill-after` / `-k` only sends `SIGTERM` — a child that blocks or ignores it hangs the pipeline past the deadline. Add `--kill-after=N` so timeout escalates to `SIGKILL`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1775")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1776(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `psql postgresql://user@host/db` (no password)",
			input:    `psql postgresql://user@host/db`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `psql postgresql://host:5432/db` (port, not password)",
			input:    `psql postgresql://host:5432/db`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `psql $PG_URL`",
			input:    `psql $PG_URL`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `psql postgresql://user:hunter2@host/db`",
			input: `psql postgresql://user:hunter2@host/db`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1776",
					Message: "`postgresql://user:SECRET@…` in argv puts the password in `ps` / `/proc/PID/cmdline` / history. Use a password file (`~/.pgpass`, `~/.my.cnf`), `PGPASSWORD` / `REDISCLI_AUTH` env var, or build the URI from a secret variable.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `mongosh mongodb+srv://u:p@cluster/db`",
			input: `mongosh "mongodb+srv://u:p@cluster/db"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1776",
					Message: "`mongodb+srv://user:SECRET@…` in argv puts the password in `ps` / `/proc/PID/cmdline` / history. Use a password file (`~/.pgpass`, `~/.my.cnf`), `PGPASSWORD` / `REDISCLI_AUTH` env var, or build the URI from a secret variable.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1776")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1777(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cp lib.so /usr/local/lib/` (unrelated file)",
			input:    `cp lib.so /usr/local/lib/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cat /etc/ld.so.preload` (read only)",
			input:    `cat /etc/ld.so.preload`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tee /etc/ld.so.preload`",
			input: `tee /etc/ld.so.preload`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1777",
					Message: "`tee /etc/ld.so.preload` writes `/etc/ld.so.preload` — linker force-loads each listed library into every process. Audit for unexpected entries; for a scoped preload use `LD_PRELOAD=` on a single invocation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `cp /tmp/x.so /etc/ld.so.preload`",
			input: `cp /tmp/x.so /etc/ld.so.preload`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1777",
					Message: "`cp /etc/ld.so.preload` writes `/etc/ld.so.preload` — linker force-loads each listed library into every process. Audit for unexpected entries; for a scoped preload use `LD_PRELOAD=` on a single invocation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1777")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1778(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `systemctl start foo.service`",
			input:    `systemctl start foo.service`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `systemctl link /etc/systemd/system/foo.service` (immutable path)",
			input:    `systemctl link /etc/systemd/system/foo.service`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `systemctl link /tmp/foo.service`",
			input: `systemctl link /tmp/foo.service`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1778",
					Message: "`systemctl link /tmp/foo.service` keeps the unit in a mutable path — a tamper later changes what systemd runs. Copy the unit into `/etc/systemd/system/` with root-only perms.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `systemctl link /home/user/build/foo.service`",
			input: `systemctl link /home/user/build/foo.service`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1778",
					Message: "`systemctl link /home/user/build/foo.service` keeps the unit in a mutable path — a tamper later changes what systemd runs. Copy the unit into `/etc/systemd/system/` with root-only perms.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1778")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1779(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `az role assignment create --role Reader --assignee u --scope s`",
			input:    `az role assignment create --role Reader --assignee u --scope s`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `az role assignment list` (read-only)",
			input:    `az role assignment list`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `az role assignment create --role Owner --assignee u --scope s`",
			input: `az role assignment create --role Owner --assignee u --scope s`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1779",
					Message: "`az role assignment create --role Owner` grants a top-of-chain role. Pick a narrower built-in role (`Reader`, specific-service Contributor) or a custom role whose permission list you can enumerate.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `az role assignment create --role=\"User Access Administrator\" --assignee u --scope s`",
			input: `az role assignment create --role="User Access Administrator" --assignee u --scope s`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1779",
					Message: "`az role assignment create --role User Access Administrator` grants a top-of-chain role. Pick a narrower built-in role (`Reader`, specific-service Contributor) or a custom role whose permission list you can enumerate.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1779")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1780(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sysctl fs.protected_symlinks=1`",
			input:    `sysctl fs.protected_symlinks=1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sysctl -a` (list all)",
			input:    `sysctl -a`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `sysctl -w fs.protected_symlinks=0`",
			input: `sysctl -w fs.protected_symlinks=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1780",
					Message: "`sysctl fs.protected_symlinks=0` disables symlink follow protection in sticky dirs — re-opens a TOCTOU race in sticky dirs. Leave the default unless you have a specific reason; otherwise scope the change to a mount namespace.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `sysctl fs.protected_hardlinks=0`",
			input: `sysctl fs.protected_hardlinks=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1780",
					Message: "`sysctl fs.protected_hardlinks=0` disables hardlink creation protection in sticky dirs — re-opens a TOCTOU race in sticky dirs. Leave the default unless you have a specific reason; otherwise scope the change to a mount namespace.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1780")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1781(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `git clone https://github.com/owner/repo.git`",
			input:    `git clone https://github.com/owner/repo.git`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git clone git@github.com:owner/repo.git` (SSH)",
			input:    `git clone git@github.com:owner/repo.git`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git clone https://token@github.com/owner/repo.git` (no password segment)",
			input:    `git clone https://token@github.com/owner/repo.git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `git clone https://user:ghp_xxx@github.com/owner/repo.git`",
			input: `git clone https://user:ghp_xxx@github.com/owner/repo.git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1781",
					Message: "`git clone https://user:ghp_xxx@github.com/owner/repo.git` puts the token in argv and `.git/config`. Use a credential helper, `GIT_ASKPASS`, or `http.extraHeader` with an env-sourced bearer.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `git fetch https://u:p@host/repo`",
			input: `git fetch https://u:p@host/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1781",
					Message: "`git fetch https://u:p@host/repo` puts the token in argv and `.git/config`. Use a credential helper, `GIT_ASKPASS`, or `http.extraHeader` with an env-sourced bearer.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1781")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1782(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `flatpak remote-add flathub https://flathub.org/repo/flathub.flatpakrepo`",
			input:    `flatpak remote-add flathub https://flathub.org/repo/flathub.flatpakrepo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `flatpak install flathub org.gimp.GIMP`",
			input:    `flatpak install flathub org.gimp.GIMP`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `flatpak remote-add --no-gpg-verify local /srv/repo`",
			input: `flatpak remote-add --no-gpg-verify local /srv/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1782",
					Message: "`flatpak remote-add --no-gpg-verify` disables signature verification — updates from this remote are accepted with only HTTPS as identity. Sign the repo (`ostree gpg-sign`) and import the key with `--gpg-import=KEYFILE`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `flatpak remote-modify --gpg-verify=false myrepo`",
			input: `flatpak remote-modify --gpg-verify=false myrepo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1782",
					Message: "`flatpak remote-modify --gpg-verify=false` disables signature verification — updates from this remote are accepted with only HTTPS as identity. Sign the repo (`ostree gpg-sign`) and import the key with `--gpg-import=KEYFILE`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1782")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1783(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `podman rmi myimage:old`",
			input:    `podman rmi myimage:old`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `podman system df` (read only)",
			input:    `podman system df`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `nerdctl system prune` (no -a, no --volumes)",
			input:    `nerdctl system prune`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `podman system reset --force`",
			input: `podman system reset --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1783",
					Message: "`podman system reset` wipes every container artifact on the host — images, volumes, networks, pods. Use narrower removals (`rmi`, `volume rm`, scoped `prune`) against the specific resource you intend to delete.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `nerdctl system prune -af --volumes`",
			input: `nerdctl system prune -af --volumes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1783",
					Message: "`nerdctl system prune -a --volumes` wipes every container artifact on the host — images, volumes, networks, pods. Use narrower removals (`rmi`, `volume rm`, scoped `prune`) against the specific resource you intend to delete.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1783")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1784(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `git config core.hooksPath .githooks` (repo-relative)",
			input:    `git config core.hooksPath .githooks`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git config user.email me@example.com`",
			input:    `git config user.email me@example.com`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `git config core.hooksPath /tmp/hooks`",
			input: `git config core.hooksPath /tmp/hooks`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1784",
					Message: "`git config core.hooksPath /tmp/hooks` runs hooks from a mutable path — supply-chain primitive. Keep hooks in the repo's `.git/hooks/` (or a tracked `.githooks/`) and point `core.hooksPath` at repo-owned paths only.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `git config --global core.hooksPath /home/attacker/hooks`",
			input: `git config --global core.hooksPath /home/attacker/hooks`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1784",
					Message: "`git config core.hooksPath /home/attacker/hooks` runs hooks from a mutable path — supply-chain primitive. Keep hooks in the repo's `.git/hooks/` (or a tracked `.githooks/`) and point `core.hooksPath` at repo-owned paths only.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1784")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1785(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ufw default deny incoming`",
			input:    `ufw default deny incoming`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ufw allow 22/tcp`",
			input:    `ufw allow 22/tcp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ufw default allow incoming`",
			input: `ufw default allow incoming`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1785",
					Message: "`ufw default allow incoming` flips the firewall baseline to accept every port that is not explicitly denied. Restore with `ufw default deny incoming` and add narrow `ufw allow <port>` rules for the services that must be reachable.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ufw default allow` (direction omitted)",
			input: `ufw default allow`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1785",
					Message: "`ufw default allow incoming` flips the firewall baseline to accept every port that is not explicitly denied. Restore with `ufw default deny incoming` and add narrow `ufw allow <port>` rules for the services that must be reachable.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1785")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1786(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mount.cifs //h/s /mnt -o credentials=/etc/cifs-creds`",
			input:    `mount.cifs //h/s /mnt -o credentials=/etc/cifs-creds`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mount.cifs //h/s /mnt -o guest`",
			input:    `mount.cifs //h/s /mnt -o guest`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mount.cifs //h/s /mnt -o username=u,password=hunter2`",
			input: `mount.cifs //h/s /mnt -o username=u,password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1786",
					Message: "`mount.cifs ... password=…` leaks the SMB password into argv / `ps` / `/proc/PID/cmdline`. Use `credentials=/path/to/creds` (mode 0600) or `$PASSWD` env var instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `mount -t cifs //h/s /mnt -o user=u,password=hunter2`",
			input: `mount -t cifs //h/s /mnt -o user=u,password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1786",
					Message: "`mount ... password=…` leaks the SMB password into argv / `ps` / `/proc/PID/cmdline`. Use `credentials=/path/to/creds` (mode 0600) or `$PASSWD` env var instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1786")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1787(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt EXTENDED_GLOB` (unrelated)",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `set -e` (unrelated)",
			input:    `set -e`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt AUTO_CD`",
			input: `setopt AUTO_CD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1787",
					Message: "`setopt AUTO_CD` turns any bare directory name into a silent `cd`. A typo or a user-controlled value reshapes `$PWD`; keep this in `~/.zshrc`, not in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt autocd` (lowercase / no underscore)",
			input: `setopt autocd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1787",
					Message: "`setopt autocd` turns any bare directory name into a silent `cd`. A typo or a user-controlled value reshapes `$PWD`; keep this in `~/.zshrc`, not in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `set -o autocd`",
			input: `set -o autocd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1787",
					Message: "`set -o autocd` turns any bare directory name into a silent `cd`. A typo or a user-controlled value reshapes `$PWD`; keep this in `~/.zshrc`, not in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1787")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1788(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ssh user@host` (default config)",
			input:    `ssh user@host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ssh -F ~/.ssh/config user@host`",
			input:    `ssh -F ~/.ssh/config user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ssh -F /tmp/ssh.conf user@host`",
			input: `ssh -F /tmp/ssh.conf user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1788",
					Message: "`ssh -F /tmp/ssh.conf` loads an alternate config from a mutable path — a tamper on that file can pin `ProxyCommand` to arbitrary code. Keep the config in `~/.ssh/config` (or a repo-owned path with `0600` perms).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `scp -F/var/tmp/conf src host:dst` (attached form)",
			input: `scp -F/var/tmp/conf src host:dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1788",
					Message: "`scp -F /var/tmp/conf` loads an alternate config from a mutable path — a tamper on that file can pin `ProxyCommand` to arbitrary code. Keep the config in `~/.ssh/config` (or a repo-owned path with `0600` perms).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1788")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1789(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt EXTENDED_GLOB`",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt CORRECT` (turning off is fine)",
			input:    `unsetopt CORRECT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CORRECT`",
			input: `setopt CORRECT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1789",
					Message: "`setopt CORRECT` enables `CORRECT` — Zsh spellcheck silently rewrites tokens that look mistyped. In a script that corrupts file paths and steals stdin for the correction prompt. Keep in `~/.zshrc`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt CORRECT_ALL`",
			input: `setopt CORRECT_ALL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1789",
					Message: "`setopt CORRECT_ALL` enables `CORRECT_ALL` — Zsh spellcheck silently rewrites tokens that look mistyped. In a script that corrupts file paths and steals stdin for the correction prompt. Keep in `~/.zshrc`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `set -o correctall`",
			input: `set -o correctall`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1789",
					Message: "`set -o correctall` enables `CORRECT_ALL` — Zsh spellcheck silently rewrites tokens that look mistyped. In a script that corrupts file paths and steals stdin for the correction prompt. Keep in `~/.zshrc`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1789")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1790(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt PIPE_FAIL` (enabling)",
			input:    `setopt PIPE_FAIL`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt PIPE_FAIL`",
			input: `unsetopt PIPE_FAIL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1790",
					Message: "`unsetopt PIPE_FAIL` returns the shell to last-command-only pipeline exit — `cmd1 | cmd2` now ignores `cmd1` failures. Scope the change to a subshell or function with `emulate -L zsh` instead of flipping it globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NOPIPEFAIL`",
			input: `setopt NOPIPEFAIL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1790",
					Message: "`setopt NOPIPEFAIL` returns the shell to last-command-only pipeline exit — `cmd1 | cmd2` now ignores `cmd1` failures. Scope the change to a subshell or function with `emulate -L zsh` instead of flipping it globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1790")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1791(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `curl https://example.com`",
			input:    `curl https://example.com`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl --unix-socket /run/user/1000/bus http://localhost/` (dbus)",
			input:    `curl --unix-socket /run/user/1000/bus http://localhost/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `curl --unix-socket /var/run/docker.sock http://localhost/containers/json`",
			input: `curl --unix-socket /var/run/docker.sock http://localhost/containers/json`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1791",
					Message: "`curl --unix-socket /var/run/docker.sock` speaks the container-daemon API — a `POST /containers/create` with `Privileged=true` is a host-root primitive. Use the CLI (`docker`/`podman`) instead.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:  "invalid — `curl http://localhost/v1/services --unix-socket /run/containerd/containerd.sock` (trailing form)",
			input: `curl http://localhost/v1/services --unix-socket /run/containerd/containerd.sock`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1791",
					Message: "`curl --unix-socket /run/containerd/containerd.sock` speaks the container-daemon API — a `POST /containers/create` with `Privileged=true` is a host-root primitive. Use the CLI (`docker`/`podman`) instead.",
					Line:    1,
					Column:  36,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1791")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1792(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `btrfs subvolume list /` (read only)",
			input:    `btrfs subvolume list /`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `btrfs device usage /` (read only)",
			input:    `btrfs device usage /`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `btrfs subvolume delete /snapshots/2025-01-01`",
			input: `btrfs subvolume delete /snapshots/2025-01-01`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1792",
					Message: "`btrfs subvolume delete` drops btrfs state with no automatic rollback — snapshots vanish on `subvolume delete`, and `device remove` can leave the filesystem degraded. Confirm the target explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `btrfs device remove $DEV /mnt/pool`",
			input: `btrfs device remove $DEV /mnt/pool`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1792",
					Message: "`btrfs device remove` drops btrfs state with no automatic rollback — snapshots vanish on `subvolume delete`, and `device remove` can leave the filesystem degraded. Confirm the target explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1792")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1793(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl certificate deny CSR_NAME`",
			input:    `kubectl certificate deny CSR_NAME`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl get csr`",
			input:    `kubectl get csr`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl certificate approve CSR_NAME`",
			input: `kubectl certificate approve CSR_NAME`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1793",
					Message: "`kubectl certificate approve` signs the identity embedded in the CSR — a `system:masters` request becomes cluster admin. Decode with `openssl req -text` first; use `kubectl certificate deny` otherwise.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1793")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1794(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cosign verify registry.example.com/app:1.2.3`",
			input:    `cosign verify registry.example.com/app:1.2.3`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cosign sign --key cosign.key registry.example.com/app:1.2.3`",
			input:    `cosign sign --key cosign.key registry.example.com/app:1.2.3`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `cosign verify --insecure-ignore-tlog img`",
			input: `cosign verify --insecure-ignore-tlog img`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1794",
					Message: "`cosign --insecure-ignore-tlog` removes a rung of the signature chain (transparency log / SCT / TLS / HTTPS-only registry). Drop the flag and fix the underlying trust anchor.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `cosign sign --allow-insecure-registry img`",
			input: `cosign sign --allow-insecure-registry img`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1794",
					Message: "`cosign --allow-insecure-registry` removes a rung of the signature chain (transparency log / SCT / TLS / HTTPS-only registry). Drop the flag and fix the underlying trust anchor.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1794")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1795(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `git remote add origin git@github.com:owner/repo.git`",
			input:    `git remote add origin git@github.com:owner/repo.git`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git remote set-url origin https://github.com/owner/repo.git`",
			input:    `git remote set-url origin https://github.com/owner/repo.git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `git remote add origin https://user:ghp_xxx@github.com/owner/repo.git`",
			input: `git remote add origin https://user:ghp_xxx@github.com/owner/repo.git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1795",
					Message: "`git remote add … https://user:ghp_xxx@github.com/owner/repo.git` stores the token in `.git/config` and leaks it via argv at creation. Use a credential helper, `GIT_ASKPASS`, or an SSH deploy key instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `git remote set-url origin https://u:p@host/repo`",
			input: `git remote set-url origin https://u:p@host/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1795",
					Message: "`git remote set-url … https://u:p@host/repo` stores the token in `.git/config` and leaks it via argv at creation. Use a credential helper, `GIT_ASKPASS`, or an SSH deploy key instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1795")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1796(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `pg_restore -d mydb backup.dump`",
			input:    `pg_restore -d mydb backup.dump`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pg_restore --list backup.dump` (TOC only)",
			input:    `pg_restore --list backup.dump`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `pg_restore -c -d mydb backup.dump`",
			input: `pg_restore -c -d mydb backup.dump`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1796",
					Message: "`pg_restore -c` drops every object in the target DB before recreating from the archive — stale or wrong-target dump silently loses data. Restore into a fresh DB (`createdb new && pg_restore -d new`), or snapshot first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1796")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1797(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ip link set eth0 up`",
			input:    `ip link set eth0 up`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ip addr show` (read only)",
			input:    `ip addr show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ip link set eth0 down`",
			input: `ip link set eth0 down`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1797",
					Message: "`ip link set … down` disables a network interface — if it carries the SSH session, the script cuts itself off. Schedule a rollback via `systemd-run --on-active=30s ip link set … up` or stage via `nmcli`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ifdown eth0`",
			input: `ifdown eth0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1797",
					Message: "`ifdown eth0` disables a network interface — if it carries the SSH session, the script cuts itself off. Schedule a rollback via `systemd-run --on-active=30s ip link set … up` or stage via `nmcli`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1797")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1798(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ufw status numbered`",
			input:    `ufw status numbered`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ufw delete 3`",
			input:    `ufw delete 3`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ufw reset`",
			input: `ufw reset`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1798",
					Message: "`ufw reset` drops every user-defined firewall rule. Snapshot (`ufw status numbered > /tmp/ufw.bak`) first, and prefer `ufw delete <num>` for targeted removals.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ufw reset --force`",
			input: `ufw reset --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1798",
					Message: "`ufw reset` drops every user-defined firewall rule. Snapshot (`ufw status numbered > /tmp/ufw.bak`) first, and prefer `ufw delete <num>` for targeted removals.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1798")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1799(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `rclone sync --dry-run src dst`",
			input:    `rclone sync --dry-run src dst`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rclone copy src dst` (copy, not sync)",
			input:    `rclone copy src dst`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `rclone sync src dst`",
			input: `rclone sync src dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1799",
					Message: "`rclone sync` deletes anything in DST that's not in SRC — empty / wrong SRC silently wipes DST. Preview with `rclone sync --dry-run`, and pin guards like `--backup-dir`, `--max-delete`, or `--min-age`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `rclone sync local: remote:bucket`",
			input: `rclone sync local: remote:bucket`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1799",
					Message: "`rclone sync` deletes anything in DST that's not in SRC — empty / wrong SRC silently wipes DST. Preview with `rclone sync --dry-run`, and pin guards like `--backup-dir`, `--max-delete`, or `--min-age`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1799")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
