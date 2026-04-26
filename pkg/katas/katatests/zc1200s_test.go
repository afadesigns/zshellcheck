// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1200(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sftp",
			input:    `sftp user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ftp",
			input: `ftp server.example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1200",
					Message: "Avoid `ftp` — it transmits credentials in plain text. Use `sftp`, `scp`, or `curl` with HTTPS for secure file transfers.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1200")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1201(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ssh",
			input:    `ssh user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid rsh",
			input: `rsh host command`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1201",
					Message: "Avoid `rsh` — it is an insecure legacy protocol. Use `ssh`/`scp`/`rsync` for encrypted remote operations.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid rcp",
			input: `rcp file host:/path`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1201",
					Message: "Avoid `rcp` — it is an insecure legacy protocol. Use `ssh`/`scp`/`rsync` for encrypted remote operations.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1201")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1202(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ip addr",
			input:    `ip addr show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ifconfig",
			input: `ifconfig eth0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1202",
					Message: "Avoid `ifconfig` — it is deprecated on modern Linux. Use `ip addr`, `ip link`, or `ip route` from iproute2.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1202")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1203(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ss",
			input:    `ss -tulnp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid netstat",
			input: `netstat -tulnp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1203",
					Message: "Avoid `netstat` — it is deprecated on modern Linux. Use `ss` from iproute2 for faster, more detailed socket statistics.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1203")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1204(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ip route",
			input:    `ip route show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid route",
			input: `route -n`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1204",
					Message: "Avoid `route` — it is deprecated on modern Linux. Use `ip route` from iproute2 for consistent routing management.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1204")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1205(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ip neigh",
			input:    `ip neigh show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid arp",
			input: `arp -a`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1205",
					Message: "Avoid `arp` — it is deprecated on modern Linux. Use `ip neigh` from iproute2 for neighbor table management.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1205")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1206(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid crontab file",
			input:    `crontab /tmp/cron.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid crontab -l",
			input:    `crontab -l`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid crontab -e",
			input: `crontab -e`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1206",
					Message: "Avoid `crontab -e` in scripts — it opens an interactive editor. Use `crontab file` or `echo '...' | crontab -` for programmatic cron management.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1206")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1207(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid chpasswd",
			input:    `chpasswd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid passwd",
			input: `passwd user`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1207",
					Message: "Avoid `passwd` in scripts — it prompts interactively. Use `chpasswd` or `usermod --password` for non-interactive password changes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1207")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1208(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid visudo -c",
			input:    `visudo -c -f /etc/sudoers.d/myconfig`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid visudo",
			input: `visudo -f /etc/sudoers`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1208",
					Message: "Avoid `visudo` in scripts — it opens an interactive editor. Write to `/etc/sudoers.d/` drop-in files and validate with `visudo -c`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1208")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1209(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid systemctl with --no-pager",
			input:    `systemctl --no-pager status nginx`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid systemctl start",
			input:    `systemctl start nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid systemctl status",
			input: `systemctl status nginx`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1209",
					Message: "Use `systemctl --no-pager` in scripts. Without it, systemctl invokes a pager that hangs in non-interactive execution.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1209")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1210(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid journalctl with --no-pager",
			input:    `journalctl --no-pager -u nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid journalctl without --no-pager",
			input: `journalctl -u nginx`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1210",
					Message: "Use `journalctl --no-pager` in scripts. Without it, journalctl invokes a pager that hangs in non-interactive execution.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1210")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1211(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git stash push -m",
			input:    `git stash push -m "wip"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid git stash pop",
			input:    `git stash pop`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bare git stash",
			input: `git stash`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1211",
					Message: "Use `git stash push -m 'description'` instead of bare `git stash`. Named stashes are easier to identify and manage.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1211")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1212(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git add file",
			input:    `git add main.go`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid git add dot",
			input: `git add .`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1212",
					Message: "Avoid `git add .` or `git add -A` — they stage everything including unintended files. Use explicit paths or `git add -p` for selective staging.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid git add -A",
			input: `git add -A`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1212",
					Message: "Avoid `git add .` or `git add -A` — they stage everything including unintended files. Use explicit paths or `git add -p` for selective staging.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1212")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1213(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid apt-get -y",
			input:    `apt-get -y install curl`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid apt-get update",
			input:    `apt-get update`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid apt-get install without -y",
			input: `apt-get install curl`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1213",
					Message: "Use `apt-get -y` in scripts. Without `-y`, apt-get prompts for confirmation which hangs in non-interactive execution.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1213")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1214(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sudo -u",
			input:    `sudo -u postgres psql`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid su",
			input: `su -c "service nginx restart" www-data`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1214",
					Message: "Avoid `su` in scripts — it prompts for a password interactively. Use `sudo -u user cmd` for non-interactive privilege switching.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1214")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1215(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid source os-release",
			input:    `source /etc/os-release`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat os-release",
			input: `cat /etc/os-release`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1215",
					Message: "Source `/etc/os-release` directly with `. /etc/os-release` instead of parsing with `cat`. It exports variables like `$ID` and `$VERSION_ID`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1215")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1216(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid dig",
			input:    `dig example.com`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid nslookup",
			input: `nslookup example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1216",
					Message: "Avoid `nslookup` — it is deprecated on many systems. Use `dig` for detailed DNS queries or `host` for simple lookups.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1216")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1217(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid systemctl",
			input:    `systemctl restart nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid service command",
			input: `service nginx restart`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1217",
					Message: "Avoid `service` — it is a SysVinit compatibility wrapper. Use `systemctl` directly on systemd systems.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1217")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1218(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid useradd with nologin",
			input:    `useradd --system --shell /sbin/nologin myservice`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid regular useradd",
			input:    `useradd newuser`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid system useradd without nologin",
			input: `useradd -r myservice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1218",
					Message: "Add `--shell /sbin/nologin` when creating system accounts with `useradd`. This prevents interactive login for service accounts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1218")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1219(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid wget to file",
			input:    `wget -O file.tar.gz https://example.com/file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid wget -qO-",
			input: `wget -qO- https://example.com/script.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1219",
					Message: "Use `curl -fsSL` instead of `wget -O -` for piped downloads. `curl` fails on HTTP errors and is available on more platforms.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1219")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1220(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid chown :group",
			input:    `chown :www-data /var/www`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid chgrp",
			input: `chgrp www-data /var/www`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1220",
					Message: "Use `chown :group file` instead of `chgrp group file`. `chown` handles both user and group changes consistently.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1220")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1221(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid fdisk -l",
			input:    `fdisk -l /dev/sda`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid fdisk interactive",
			input: `fdisk /dev/sda`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1221",
					Message: "Avoid `fdisk` in scripts — it is interactive. Use `parted -s` or `sfdisk` for scriptable disk partitioning.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1221")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1222(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid lsof for files",
			input:    `lsof /var/log/syslog`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid lsof -i",
			input: `lsof -i :8080`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1222",
					Message: "Use `ss -tlnp` instead of `lsof -i` for port checks. `ss` is faster and doesn't require elevated permissions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1222")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1223(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ip -br addr",
			input:    `ip -br addr`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ip addr show",
			input: `ip addr show`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1223",
					Message: "Use `ip -br addr` for machine-readable output instead of parsing `ip addr show` with grep or awk.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1223")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1224(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid /proc/meminfo",
			input:    `cat /proc/meminfo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid free",
			input: `free -m`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1224",
					Message: "Avoid parsing `free` output — its format varies across versions. Read `/proc/meminfo` directly for reliable memory information.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1224")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1225(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid /proc/uptime",
			input:    `cat /proc/uptime`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid uptime",
			input: `uptime -p`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1225",
					Message: "Avoid parsing `uptime` — its output varies by locale. Read `/proc/uptime` for machine-parseable seconds since boot.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1225")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1226(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid dmesg -T",
			input:    `dmesg -T`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid dmesg without -T",
			input: `dmesg -l err`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1226",
					Message: "Use `dmesg -T` for human-readable timestamps instead of raw kernel boot-seconds. Or use `--time-format=iso` for ISO 8601.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1226")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1227(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid curl -fsSL",
			input:    `curl -fsSL https://example.com/file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid curl without URL",
			input:    `curl -s localhost`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid curl without -f",
			input: `curl -sL https://example.com/install.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1227",
					Message: "Use `curl -f` to fail on HTTP errors. Without `-f`, curl silently returns error pages (404, 500) as if they were successful.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1227")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1228(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ssh with BatchMode",
			input:    `ssh -o BatchMode=yes user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ssh without policy",
			input: `ssh user@host ls`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1228",
					Message: "Use `ssh -o BatchMode=yes` or `-o StrictHostKeyChecking=accept-new` in scripts. Without these, ssh may prompt interactively and hang.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1228")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1229(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid rsync",
			input:    `rsync -az src/ user@host:dst/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid scp",
			input: `scp file.tar.gz user@host:/tmp/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1229",
					Message: "Prefer `rsync -az` over `scp` for file transfers. `rsync` supports delta transfers, resume, and is more efficient.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1229")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1230(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ping -c",
			input:    `ping -c 3 localhost`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ping without -c",
			input: `ping localhost`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1230",
					Message: "Use `ping -c N` in scripts. Without `-c`, ping runs indefinitely on Linux and will hang the script.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1230")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1231(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git clone --depth 1",
			input:    `git clone --depth 1 https://github.com/user/repo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid git clone full",
			input: `git clone https://github.com/user/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1231",
					Message: "Consider `git clone --depth 1` in scripts. Full clones download entire history which is unnecessary for builds and CI.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1231")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1232(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid pip install --user",
			input:    `pip install --user requests`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid pip3 list",
			input:    `pip3 list`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bare pip install",
			input: `pip install requests`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1232",
					Message: "Use `pip install --user` or a virtualenv instead of bare `pip install`. System-wide pip installs can break OS package managers.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1232")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1233(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid npx",
			input:    `npx create-react-app myapp`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid npm install local",
			input:    `npm install express`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid npm install -g",
			input: `npm install -g typescript`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1233",
					Message: "Avoid `npm install -g`. Use `npx` for one-off tool execution or `npm install --save-dev` for project-scoped dependencies.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1233")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1234(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid docker run --rm",
			input:    `docker run --rm alpine echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid docker run -d",
			input:    `docker run -d nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid docker run without --rm",
			input: `docker run alpine echo hello`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1234",
					Message: "Use `docker run --rm` to auto-remove containers after exit. Without `--rm`, stopped containers accumulate on disk.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1234")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1235(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git push",
			input:    `git push origin main`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid git push -f",
			input: `git push -f origin main`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1235",
					Message: "Use `git push --force-with-lease` instead of `-f`/`--force`. It prevents overwriting remote changes made by others.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1235")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1236(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git reset --soft",
			input:    `git reset --soft HEAD~1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid git reset --hard",
			input: `git reset --hard HEAD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1236",
					Message: "Avoid `git reset --hard` — it permanently discards uncommitted changes. Use `git stash` first, or `git reset --soft` to keep changes staged.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1236")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1237(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git clean -n",
			input:    `git clean -n`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid git clean -fd",
			input: `git clean -fd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1237",
					Message: "Use `git clean -n` first to preview removals before `git clean -fd`. Forced clean permanently deletes untracked files.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1237")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1238(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid docker exec without -it",
			input:    `docker exec mycontainer ls`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid docker exec -it",
			input: `docker exec -it mycontainer bash`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1238",
					Message: "Avoid `docker exec -it` in scripts — TTY allocation hangs without a terminal. Use `docker exec` without `-it` for non-interactive commands.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1238")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1239(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid kubectl exec without -it",
			input:    `kubectl exec mypod -- ls`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid kubectl exec -it",
			input: `kubectl exec -it mypod -- bash`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1239",
					Message: "Avoid `kubectl exec -it` in scripts — TTY allocation hangs without a terminal. Use `kubectl exec pod -- cmd` for non-interactive execution.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1239")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1240(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid find with maxdepth and delete",
			input:    `find /tmp -maxdepth 1 -name "*.tmp" -delete`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid find without delete",
			input:    `find . -name "*.log"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid find -delete without maxdepth",
			input: `find /tmp -name "*.tmp" -delete`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1240",
					Message: "Use `find -maxdepth N` with `-delete` to limit deletion scope. Without depth limits, find recurses infinitely.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1240")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1241(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid xargs -0 rm",
			input:    `xargs -0 rm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid xargs without rm",
			input:    `xargs grep pattern`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid xargs rm without -0",
			input: `xargs rm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1241",
					Message: "Use `xargs -0 rm` with `find -print0` for safe deletion. Without `-0`, filenames with spaces or special characters break.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1241")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1242(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid tar with -C",
			input:    `tar xzf archive.tar.gz -C /opt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid tar create",
			input:    `tar czf archive.tar.gz dir`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid tar extract without -C",
			input: `tar xzf archive.tar.gz`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1242",
					Message: "Use `tar -C dir` to specify extraction directory. Without `-C`, tar extracts into the current directory which may overwrite files.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1242")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1243(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -lZ",
			input:    `grep -lZ pattern dir`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid grep without -l",
			input:    `grep pattern file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep -l without -Z",
			input: `grep -l pattern dir`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1243",
					Message: "Use `grep -lZ` instead of `grep -l` for null-terminated file lists. Pair with `xargs -0` to safely handle filenames with special characters.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1243")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1244(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid mv -n",
			input:    `mv -n src dst`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid mv -f explicit",
			input:    `mv -f old new`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bare mv",
			input: `mv file.txt backup.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1244",
					Message: "Consider `mv -n` to prevent overwriting existing files. Without `-n`, `mv` silently overwrites the target.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1244")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1245(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid curl with TLS",
			input:    `curl -fsSL https://example.com`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid curl -k",
			input: `curl -k https://example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1245",
					Message: "Avoid `curl -k`/`--insecure` — it disables TLS certificate verification. Fix the certificate chain or use `--cacert` to specify a CA bundle.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1245")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1246(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid mysql with -p prompt",
			input:    `mysql -u root -p mydb`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid mysql with inline password",
			input: `mysql -u root -pMySecret mydb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1246",
					Message: "Avoid passing passwords as command arguments — they appear in process lists. Use environment variables (e.g., `MYSQL_PWD`) or credential files instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1246")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1247(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid chmod 755",
			input:    `chmod 755 script.sh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid chmod 2755",
			input: `chmod 2755 /usr/local/bin/tool`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1247",
					Message: "Avoid `chmod +s` — setuid/setgid bits create privilege escalation risks. Use `sudo`, capabilities (`setcap`), or dedicated service accounts instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid chmod 4755",
			input: `chmod 4755 /usr/local/bin/tool`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1247",
					Message: "Avoid `chmod +s` — setuid/setgid bits create privilege escalation risks. Use `sudo`, capabilities (`setcap`), or dedicated service accounts instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1247")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1248(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ufw",
			input:    `ufw allow 22`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid iptables",
			input: `iptables -A INPUT -p tcp --dport 22 -j ACCEPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1248",
					Message: "Prefer `ufw` or `firewalld` over raw `iptables`. Firewall frontends provide persistent, manageable rules.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1248")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1249(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ssh-keygen -f",
			input:    `ssh-keygen -t ed25519 -f /tmp/key`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ssh-keygen without -f",
			input: `ssh-keygen -t rsa`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1249",
					Message: "Use `ssh-keygen -f /path/to/key -N ''` in scripts. Without `-f`, ssh-keygen prompts interactively for the file path.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1249")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1250(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid gpg -b -d",
			input:    `gpg -b -d file.gpg`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid gpg without operation",
			input:    `gpg -k`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid gpg -d without -b",
			input: `gpg -d file.gpg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1250",
					Message: "Use `gpg --batch` in scripts for non-interactive operation. Without `--batch`, gpg may prompt for passphrases or confirmations.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1250")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1251(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid mount with -o",
			input:    `mount -o noexec,nosuid /dev/sdb1 /mnt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid mount without device",
			input:    `mount -a`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid mount device without -o",
			input: `mount /dev/sdb1 /mnt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1251",
					Message: "Use `mount -o noexec,nosuid,nodev` when mounting external media. Without restrictions, mounted filesystems can contain executable exploits.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1251")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1252(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid getent",
			input:    `getent passwd root`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat /etc/passwd",
			input: `cat /etc/passwd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1252",
					Message: "Use `getent` instead of `cat /etc/passwd`. `getent` queries all NSS sources including LDAP and SSSD.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1252")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1253(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid docker build --no-cache",
			input:    `docker build --no-cache -t myapp .`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid docker run (not build)",
			input:    `docker run --rm alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid docker build without --no-cache",
			input: `docker build -t myapp .`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1253",
					Message: "Consider `docker build --no-cache` in CI for reproducible builds. Layer caching can mask changed dependencies.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1253")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1254(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git commit",
			input:    `git commit -m "feat: add feature"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid git commit --amend",
			input: `git commit --amend -m "fix"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1254",
					Message: "Avoid `git commit --amend` on shared branches — it rewrites history. Use `git commit --fixup` or create a new commit instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1254")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1255(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid curl -fsSL",
			input:    `curl -fsSL https://example.com/file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid curl without -L",
			input: `curl -s https://example.com/file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1255",
					Message: "Use `curl -L` to follow HTTP redirects. Without `-L`, curl returns redirect responses (301/302) instead of the actual content.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1255")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1256(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid other command",
			input:    `mkdir /tmp/dir`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid mkfifo without trap",
			input: `mkfifo /tmp/mypipe`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1256",
					Message: "Set up `trap 'rm -f pipe' EXIT` after `mkfifo`. Named pipes persist on the filesystem and need explicit cleanup.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1256")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1257(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid docker stop -t",
			input:    `docker stop -t 5 mycontainer`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid docker stop without -t",
			input: `docker stop mycontainer`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1257",
					Message: "Use `docker stop -t N` to set an explicit shutdown timeout. The default 10s may be too long or too short for your use case.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1257")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1258(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid rsync --delete",
			input:    `rsync -az --delete src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid rsync single file",
			input:    `rsync -az file host:path`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid rsync dir without --delete",
			input: `rsync -az src/ dst/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1258",
					Message: "Consider `rsync --delete` for directory sync. Without `--delete`, files removed from source remain on the destination.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1258")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1259(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid docker pull with tag",
			input:    `docker pull alpine:3.19`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid docker pull without tag",
			input: `docker pull nginx`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1259",
					Message: "Pin Docker image to a specific tag instead of defaulting to `:latest`. Untagged pulls are non-reproducible and may break unexpectedly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1259")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1260(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git branch -d",
			input:    `git branch -d feature`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid git branch -D",
			input: `git branch -D feature`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1260",
					Message: "Use `git branch -d` instead of `-D`. The lowercase `-d` refuses to delete unmerged branches, preventing accidental data loss.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1260")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1261(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid base64 encode",
			input:    `base64 file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid base64 -d",
			input: `base64 -d encoded.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1261",
					Message: "Inspect `base64 -d` output before piping to execution. Blindly executing decoded content is a code injection vector.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1261")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1262(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid chmod -R 755",
			input:    `chmod -R 755 /var/www`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid chmod 777 non-recursive",
			input:    `chmod 777 /tmp/test`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid chmod -R 777",
			input: `chmod -R 777 /var/www`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1262",
					Message: "Never use `chmod -R 777` — it makes everything world-writable. Use `find -type d -exec chmod 755` and `find -type f -exec chmod 644` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1262")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1263(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid apt-get",
			input:    `apt-get install curl`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid apt in script",
			input: `apt install curl`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1263",
					Message: "Use `apt-get` instead of `apt` in scripts. `apt` is for interactive use; `apt-get` has a stable scripting interface.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1263")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1264(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid dnf",
			input:    `dnf install -y curl`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid yum",
			input: `yum install -y curl`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1264",
					Message: "Use `dnf` instead of `yum`. `yum` is deprecated on modern Fedora and RHEL; `dnf` has better dependency resolution.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1264")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1265(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid systemctl enable --now",
			input:    `systemctl enable --now nginx`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid systemctl start",
			input:    `systemctl start nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid systemctl enable without --now",
			input: `systemctl enable nginx`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1265",
					Message: "Use `systemctl enable --now` to enable and start the service immediately. Without `--now`, the service only starts on next boot.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1265")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1266(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid nproc",
			input:    `nproc`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat /proc/cpuinfo",
			input: `cat /proc/cpuinfo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1266",
					Message: "Use `nproc` instead of parsing `/proc/cpuinfo` for CPU count. `nproc` is portable and available on Linux and macOS (via coreutils).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1266")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1267(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid df -P",
			input:    `df -P /`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid df -h without -P",
			input: `df -h /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1267",
					Message: "Use `df -P` for script-safe output. `df -h` format varies across systems and may split long device names across lines.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1267")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1268(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid du with specific path",
			input:    `du -sh /tmp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid du -sh *",
			input: `du -sh *`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1268",
					Message: "Use `du -sh -- *` instead of `du -sh *`. The `--` prevents filenames starting with `-` from being interpreted as options.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1268")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1269(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid pgrep usage",
			input:    `pgrep -f myprocess`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid ps with no grep-related args",
			input:    `ps -p 1234`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ps aux for process search",
			input: `ps aux`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1269",
					Message: "Use `pgrep` instead of `ps aux | grep`. `pgrep` is purpose-built for process searching and doesn't match itself.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid ps -ef for process search",
			input: `ps -ef`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1269",
					Message: "Use `pgrep` instead of `ps -ef | grep`. `pgrep` is purpose-built for process searching and doesn't match itself.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid ps -e for process search",
			input: `ps -e`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1269",
					Message: "Use `pgrep` instead of `ps -e | grep`. `pgrep` is purpose-built for process searching and doesn't match itself.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1269")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1270(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid mktemp usage",
			input:    `local tmpfile=$(mktemp)`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid touch with non-tmp path",
			input:    `touch /var/log/myapp.log`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid dynamic tmp path",
			input:    `touch /tmp/$USER-cache`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid hardcoded tmp touch",
			input: `touch /tmp/myfile.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1270",
					Message: "Use `mktemp` instead of hardcoded `/tmp/myfile.txt`. Hardcoded `/tmp` paths are vulnerable to symlink attacks and race conditions.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid hardcoded tmp cat",
			input: `cat /tmp/output.log`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1270",
					Message: "Use `mktemp` instead of hardcoded `/tmp/output.log`. Hardcoded `/tmp` paths are vulnerable to symlink attacks and race conditions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1270")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1271(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid command -v usage",
			input:    `command -v git`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid type command",
			input:    `type git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid which usage",
			input: `which git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1271",
					Message: "Use `command -v` instead of `which`. `command -v` is POSIX-compliant and built into Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1271")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1272(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid install usage",
			input:    `install -m 0755 mybin /usr/local/bin`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cp to non-system dir",
			input:    `cp file.txt /home/user/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cp to /usr/local/bin",
			input: `cp mybin /usr/local/bin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1272",
					Message: "Use `install -m 0755` instead of `cp` to system directories. `install` sets permissions atomically.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid cp to /usr/bin",
			input: `cp mybin /usr/bin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1272",
					Message: "Use `install -m 0755` instead of `cp` to system directories. `install` sets permissions atomically.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1272")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1273(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -q usage",
			input:    `grep -q pattern file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid grep without redirect",
			input:    `grep pattern file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep redirected to /dev/null",
			input: `grep pattern file.txt /dev/null`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1273",
					Message: "Use `grep -q` instead of redirecting to `/dev/null`. It is faster and more idiomatic.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1273")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1274(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid Zsh :t modifier",
			input:    `echo ${filepath:t}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid non-basename command",
			input:    `dirname /path/to/file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid basename usage",
			input: `basename /path/to/file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1274",
					Message: "Use Zsh parameter expansion `${var:t}` instead of `basename`. The `:t` modifier extracts the filename without forking a process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1274")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1275(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid Zsh :h modifier",
			input:    `echo ${filepath:h}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid non-dirname command",
			input:    `basename /path/to/file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid dirname usage",
			input: `dirname /path/to/file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1275",
					Message: "Use Zsh parameter expansion `${var:h}` instead of `dirname`. The `:h` modifier extracts the directory without forking a process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1275")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1276(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid brace expansion",
			input:    `echo {1..10}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid other command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid seq usage",
			input: `seq 1 10`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1276",
					Message: "Use Zsh brace expansion `{start..end}` instead of `seq`. Brace expansion is built-in and avoids forking.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid seq with single arg",
			input: `seq 5`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1276",
					Message: "Use Zsh brace expansion `{start..end}` instead of `seq`. Brace expansion is built-in and avoids forking.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1276")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

// ZC1277 was retired as a duplicate of ZC1108. Kept as a no-op stub so
// legacy `disabled_katas` lists keep parsing; the detection runs under
// ZC1108 now.

func TestZC1277Stub(t *testing.T) {
	cases := []string{
		"echo hi",
		"ls",
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			v := testutil.Check(in, "ZC1277")
			testutil.AssertViolations(t, in, v, []katas.Violation{})
		})
	}
}

// ZC1278 was retired as a duplicate of ZC1009. Kept as a no-op stub so
// legacy `disabled_katas` lists keep parsing; the detection runs under
// ZC1009 now.

func TestZC1278Stub(t *testing.T) {
	cases := []string{
		"echo hi",
		"ls",
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			v := testutil.Check(in, "ZC1278")
			testutil.AssertViolations(t, in, v, []katas.Violation{})
		})
	}
}

func TestZC1279(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid realpath usage",
			input:    `realpath /some/path`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid readlink without -f",
			input:    `readlink /some/symlink`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid readlink -f",
			input: `readlink -f /some/path`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1279",
					Message: "Use `realpath` instead of `readlink -f`. `realpath` is more portable, especially on macOS.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1279")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1280(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid Zsh :e modifier",
			input:    `echo ${file:e}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cut with different delimiter",
			input:    `cut -d: -f1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cut to extract extension",
			input: `cut -d. -f2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1280",
					Message: "Use Zsh parameter expansion `${var:e}` to extract the file extension instead of `cut -d. -f2`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1280")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1281(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sort -u usage",
			input:    `sort -u file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid uniq with flags",
			input:    `uniq -c file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid uniq with file",
			input: `uniq file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1281",
					Message: "Use `sort -u` instead of `sort | uniq`. The `-u` flag deduplicates in a single pass.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1281")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1282(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid Zsh :r modifier",
			input:    `echo ${file:r}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sed with different pattern",
			input:    `sed s/foo/bar/g`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sed to strip extension",
			input: `sed 's/\.[^.]*$//'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1282",
					Message: "Use Zsh parameter expansion `${var:r}` to remove the file extension instead of `sed`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1282")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1283(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid setopt usage",
			input:    `setopt noglob`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid unsetopt usage",
			input:    `unsetopt noglob`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid set without -o",
			input:    `set -e`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid set -o",
			input: `set -o noglob`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1283",
					Message: "Use `setopt` instead of `set -o` in Zsh scripts. `setopt` is the native Zsh idiom.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1283")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1284(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid Zsh split expansion",
			input:    `echo ${(s/:/)PATH}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cut with dot delimiter (covered by ZC1280)",
			input:    `cut -d. -f2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cut with colon delimiter",
			input: `cut -d: -f1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1284",
					Message: "Use Zsh parameter expansion `${(s:sep:)var}` for field splitting instead of `cut -d -f`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid cut with comma delimiter",
			input: `cut -d, -f3`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1284",
					Message: "Use Zsh parameter expansion `${(s:sep:)var}` for field splitting instead of `cut -d -f`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1284")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1285(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sort with flags",
			input:    `sort -n file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sort with key",
			input:    `sort -k 2`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sort with reverse flag",
			input:    `sort -r file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "sort with single file argument",
			input: `sort data.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1285",
					Message: "Use Zsh `${(o)array}` for sorting instead of piping to `sort`. The `(o)` flag sorts in-shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1285")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1286(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep without -v",
			input:    `grep pattern file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid Zsh array filter",
			input:    `echo ${array:#pattern}`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep -v for filtering",
			input: `grep -v pattern file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1286",
					Message: "Use Zsh `${array:#pattern}` for filtering instead of `grep -v`. Parameter expansion avoids a subprocess.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1286")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1287(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid cat with file",
			input:    `cat file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cat with -n flag",
			input:    `cat -n file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat -v for visible chars",
			input: `cat -v file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1287",
					Message: "Use Zsh `${(V)var}` to make control characters visible instead of `cat -v`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid cat -A for all visible",
			input: `cat -A file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1287",
					Message: "Use Zsh `${(V)var}` to make control characters visible instead of `cat -v`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1287")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1288(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid typeset usage",
			input:    `typeset -A mymap`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid local usage",
			input:    `local myvar=hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid declare usage",
			input: `declare -A mymap`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1288",
					Message: "Use `typeset` instead of `declare` in Zsh scripts. `typeset` is the native Zsh idiom.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid declare with -i flag",
			input: `declare -i count`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1288",
					Message: "Use `typeset` instead of `declare` in Zsh scripts. `typeset` is the native Zsh idiom.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1288")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1289(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sort with numeric and unique",
			input:    `sort -n -u file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sort without unique",
			input:    `sort file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sort -u alone",
			input: `sort -u file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1289",
					Message: "Use Zsh `${(u)array}` for unique elements instead of `sort -u`. The `(u)` flag preserves order.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1289")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1290(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sort -n with -r flag",
			input:    `sort -n -r file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sort without -n",
			input:    `sort file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sort -n alone",
			input: `sort -n file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1290",
					Message: "Use Zsh `${(n)array}` for numeric sorting instead of `sort -n`. The `(n)` flag sorts numerically in-shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1290")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1291(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sort -r with -n",
			input:    `sort -r -n file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sort without -r",
			input:    `sort file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sort -r alone",
			input: `sort -r file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1291",
					Message: "Use Zsh `${(O)array}` for reverse sorting instead of `sort -r`. The `(O)` flag sorts descending in-shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1291")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1292(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid tr with character ranges",
			input:    `tr 'a-z' 'A-Z'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid tr with POSIX classes",
			input:    `tr '[:upper:]' '[:lower:]'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid tr with single char translation",
			input: `tr '/' '_'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1292",
					Message: "Use Zsh `${var////_}` for character substitution instead of `tr`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1292")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1293(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid [[ ]] usage",
			input:    `[[ -f file.txt ]]`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid test command",
			input: `test -f file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1293",
					Message: "Use `[[ ]]` instead of the `test` command in Zsh. `[[ ]]` is more powerful and does not require variable quoting.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid test with -z flag",
			input: `test -z "$var"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1293",
					Message: "Use `[[ ]]` instead of the `test` command in Zsh. `[[ ]]` is more powerful and does not require variable quoting.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1293")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1294(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid bindkey usage",
			input:    `bindkey -e`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bind usage",
			input: `bind -x '"\C-r": history-search'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1294",
					Message: "Use `bindkey` instead of `bind` in Zsh. `bind` is a Bash builtin; Zsh uses `bindkey` for ZLE key bindings.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1294")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1295(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid vared usage",
			input:    `vared myvar`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid read without -e",
			input:    `read -r myvar`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid read -e for editing",
			input: `read -e myvar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1295",
					Message: "Use `vared` instead of `read -e` in Zsh. `vared` provides full ZLE editing support natively.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1295")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1296(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid setopt usage",
			input:    `setopt extendedglob`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid shopt usage",
			input: `shopt -s extglob`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1296",
					Message: "Avoid `shopt` in Zsh — it is a Bash builtin. Use `setopt`/`unsetopt` for Zsh shell options.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1296")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1297(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid $0 usage",
			input:    `echo $0`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_SOURCE usage",
			input: `echo $BASH_SOURCE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1297",
					Message: "Avoid `$BASH_SOURCE` in Zsh — use `$0` or `${(%):-%x}` instead. `BASH_SOURCE` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1297")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1298(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid funcstack usage",
			input:    `echo $funcstack`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid FUNCNAME usage",
			input: `echo $FUNCNAME`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1298",
					Message: "Avoid `$FUNCNAME` in Zsh — use `$funcstack` instead. `FUNCNAME` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1298")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1299(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid funcfiletrace usage",
			input:    `echo $funcfiletrace`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_LINENO usage",
			input: `echo $BASH_LINENO`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1299",
					Message: "Avoid `$BASH_LINENO` in Zsh — use `$funcfiletrace` instead. `BASH_LINENO` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1299")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
