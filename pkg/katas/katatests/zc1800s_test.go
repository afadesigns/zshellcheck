// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1800(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `pg_ctl stop -m fast`",
			input:    `pg_ctl stop -m fast`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pg_ctl start` (no stop)",
			input:    `pg_ctl start`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `pg_ctl stop -m immediate -D /var/lib/pg`",
			input: `pg_ctl stop -m immediate -D /var/lib/pg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1800",
					Message: "`pg_ctl stop -m immediate` kills the postmaster without a shutdown checkpoint — WAL replay on restart can lose committed transactions if WAL is corrupt. Use `-m smart` or `-m fast` for routine shutdowns.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pg_ctl restart --mode=immediate`",
			input: `pg_ctl restart --mode=immediate`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1800",
					Message: "`pg_ctl stop -m immediate` kills the postmaster without a shutdown checkpoint — WAL replay on restart can lose committed transactions if WAL is corrupt. Use `-m smart` or `-m fast` for routine shutdowns.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1800")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1801(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `fwupdmgr get-devices` (read only)",
			input:    `fwupdmgr get-devices`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `fwupdmgr refresh` (metadata, not flash)",
			input:    `fwupdmgr refresh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `fwupdmgr update` (all devices)",
			input: `fwupdmgr update`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1801",
					Message: "`fwupdmgr update` flashes firmware — a mid-write interruption can brick BIOS, SSD, Thunderbolt, or NIC microcontrollers. Inhibit reboot triggers (`systemd-inhibit`) and ensure battery / UPS before running.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `fwupdmgr install firmware.cab`",
			input: `fwupdmgr install firmware.cab`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1801",
					Message: "`fwupdmgr install` flashes firmware — a mid-write interruption can brick BIOS, SSD, Thunderbolt, or NIC microcontrollers. Inhibit reboot triggers (`systemd-inhibit`) and ensure battery / UPS before running.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1801")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1802(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `dnf history list`",
			input:    `dnf history list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `dnf history info 5`",
			input:    `dnf history info 5`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `dnf history undo 5`",
			input: `dnf history undo 5`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1802",
					Message: "`dnf history undo` reverses the past transaction — deps downgrade, security patches can get reverted. Review with `dnf history info`, or restore a filesystem snapshot.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `yum history rollback 3`",
			input: `yum history rollback 3`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1802",
					Message: "`yum history rollback` reverses the past transaction — deps downgrade, security patches can get reverted. Review with `dnf history info`, or restore a filesystem snapshot.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1802")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1803(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mysql --ssl-mode=VERIFY_IDENTITY -h db -u u`",
			input:    `mysql --ssl-mode=VERIFY_IDENTITY -h db -u u`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `psql postgresql://u@db/mydb`",
			input:    `psql postgresql://u@db/mydb`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mysql --skip-ssl -h db -u u`",
			input: `mysql --skip-ssl -h db -u u`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1803",
					Message: "`mysql --skip-ssl` disables TLS — login handshake and queries travel in plaintext. Use `--ssl-mode=VERIFY_IDENTITY` (MySQL) / `sslmode=verify-full` (psql) with a pinned CA.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:  "invalid — `psql \"host=db sslmode=disable user=u\"`",
			input: `psql "host=db sslmode=disable user=u"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1803",
					Message: "`psql host=db sslmode=disable user=u` disables TLS — login handshake and queries travel in plaintext. Use `--ssl-mode=VERIFY_IDENTITY` (MySQL) / `sslmode=verify-full` (psql) with a pinned CA.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1803")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1804(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `aws ec2 describe-instances`",
			input:    `aws ec2 describe-instances`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `aws ec2 terminate-instances --instance-ids i-1 --dry-run`",
			input:    `aws ec2 terminate-instances --instance-ids i-1 --dry-run`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `aws ec2 terminate-instances --instance-ids i-1 i-2`",
			input: `aws ec2 terminate-instances --instance-ids i-1 i-2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1804",
					Message: "`aws ec2 terminate-instances` tears down EC2 instance(s) and their instance-store volumes with no automatic backup. Review with `aws ec2 describe-…`, add `--dry-run` to verify the target, and pin IDs through `--cli-input-json`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `aws ec2 delete-snapshot --snapshot-id snap-abc`",
			input: `aws ec2 delete-snapshot --snapshot-id snap-abc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1804",
					Message: "`aws ec2 delete-snapshot` deletes the EBS / RDS snapshot with no automatic backup. Review with `aws ec2 describe-…`, add `--dry-run` to verify the target, and pin IDs through `--cli-input-json`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1804")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1805(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `aws dynamodb describe-table --table-name mytbl`",
			input:    `aws dynamodb describe-table --table-name mytbl`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `aws cloudformation list-stacks`",
			input:    `aws cloudformation list-stacks`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `aws cloudformation delete-stack --stack-name prod`",
			input: `aws cloudformation delete-stack --stack-name prod`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1805",
					Message: "`aws cloudformation delete-stack` removes every resource the stack manages, no rollback. Stage a confirmation, pin IDs via `--cli-input-json`, and export a backup first where the service supports one.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `aws kms schedule-key-deletion --key-id k`",
			input: `aws kms schedule-key-deletion --key-id k`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1805",
					Message: "`aws kms schedule-key-deletion` queues CMK deletion — ciphertext becomes unreadable after the grace window. Stage a confirmation, pin IDs via `--cli-input-json`, and export a backup first where the service supports one.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1805")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1806(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `zmv -n '*.JPG' '*.jpg'` (dry-run)",
			input:    `zmv -n '*.JPG' '*.jpg'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `zmv -i '(*).txt' '$1.md'` (interactive)",
			input:    `zmv -i '(*).txt' '$1.md'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `zmv` alone (help)",
			input:    `zmv`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `zmv '*.txt' '*.md'`",
			input: `zmv '*.txt' '*.md'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1806",
					Message: "`zmv` without `-n` (dry-run) or `-i` (interactive) renames every matched file in one shot — a pattern typo can collide names. Preview with `zmv -n`, then re-run once the list looks right.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `zmv -W '(*).jpg' 'archive/$1.jpg'`",
			input: `zmv -W '(*).jpg' 'archive/$1.jpg'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1806",
					Message: "`zmv` without `-n` (dry-run) or `-i` (interactive) renames every matched file in one shot — a pattern typo can collide names. Preview with `zmv -n`, then re-run once the list looks right.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1806")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1807(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gh api /repos/owner/repo`",
			input:    `gh api /repos/owner/repo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh api -X GET /repos/owner/repo`",
			input:    `gh api -X GET /repos/owner/repo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gh api -X DELETE /repos/owner/repo`",
			input: `gh api -X DELETE /repos/owner/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1807",
					Message: "`gh api -X DELETE` sends a raw DELETE to the GitHub API with the caller's token — no `--yes` guard, no dry-run. Use the high-level `gh` subcommand for the target, or wrap with a preflight GET + explicit confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh api --method=DELETE /repos/owner/repo/releases/123`",
			input: `gh api --method=DELETE /repos/owner/repo/releases/123`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1807",
					Message: "`gh api -X DELETE` sends a raw DELETE to the GitHub API with the caller's token — no `--yes` guard, no dry-run. Use the high-level `gh` subcommand for the target, or wrap with a preflight GET + explicit confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1807")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1808(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl apply -f deploy.yaml`",
			input:    `kubectl apply -f deploy.yaml`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl replace -f deploy.yaml` (no --force)",
			input:    `kubectl replace -f deploy.yaml`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl replace --force -f deploy.yaml`",
			input: `kubectl replace --force -f deploy.yaml`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1808",
					Message: "`kubectl replace --force` is delete + create — pods die, PDBs are ignored, in-flight requests drop. Prefer `kubectl apply -f FILE` and reserve `replace --force` for schema changes `apply` cannot patch, after draining traffic.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1808")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1809(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gsutil ls gs://bucket`",
			input:    `gsutil ls gs://bucket`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gsutil rm gs://bucket/specific-object` (single object)",
			input:    `gsutil rm gs://bucket/specific-object`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gsutil -m rm -r gs://bucket/prefix`",
			input: `gsutil -m rm -r gs://bucket/prefix`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1809",
					Message: "`gsutil rm` with recursive/force deletes every matching GCS object (or the bucket itself). Preview with `gsutil ls`, enable Object Versioning / retention locks ahead of time, and prefer narrower object-level `gsutil rm` calls.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gsutil rb -f gs://bucket`",
			input: `gsutil rb -f gs://bucket`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1809",
					Message: "`gsutil rb` with recursive/force deletes every matching GCS object (or the bucket itself). Preview with `gsutil ls`, enable Object Versioning / retention locks ahead of time, and prefer narrower object-level `gsutil rm` calls.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1809")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1810(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `wget https://example.com/file.tar.gz`",
			input:    `wget https://example.com/file.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `wget -r --level=2 https://example.com/`",
			input:    `wget -r --level=2 https://example.com/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `wget -r -l3 https://example.com/`",
			input:    `wget -r -l3 https://example.com/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `wget -r https://example.com/`",
			input: `wget -r https://example.com/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1810",
					Message: "`wget -r` / `--mirror` without `--level=N` follows links to arbitrary depth — the crawl can exhaust disk and climb into parent paths. Pin `--level=3`, add `--no-parent`, and cap with `--quota=1G`.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:  "invalid — `wget --mirror https://example.com/`",
			input: `wget --mirror https://example.com/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1810",
					Message: "`wget -r` / `--mirror` without `--level=N` follows links to arbitrary depth — the crawl can exhaust disk and climb into parent paths. Pin `--level=3`, add `--no-parent`, and cap with `--quota=1G`.",
					Line:    1,
					Column:  7,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1810")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1811(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `chown -R user:group /srv/app`",
			input:    `chown -R user:group /srv/app`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `chmod -R 0750 /srv/app`",
			input:    `chmod -R 0750 /srv/app`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `chown -R --no-preserve-root user /target`",
			input: `chown -R --no-preserve-root user /target`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1811",
					Message: "`chown --no-preserve-root` disables the GNU safeguard against recursing into `/`. Remove the flag; list explicit paths instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `chmod -R --no-preserve-root 0755 /target`",
			input: `chmod -R --no-preserve-root 0755 /target`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1811",
					Message: "`chmod --no-preserve-root` disables the GNU safeguard against recursing into `/`. Remove the flag; list explicit paths instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1811")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1812(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `aws ssm put-parameter --type String --value plain --name /app/region`",
			input:    `aws ssm put-parameter --type String --value plain --name /app/region`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `aws ssm put-parameter --type SecureString --value file://secret --name /app/token`",
			input:    `aws ssm put-parameter --type SecureString --value file://secret --name /app/token`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `aws ssm put-parameter --type SecureString --value hunter2 --name /app/token`",
			input: `aws ssm put-parameter --type SecureString --value hunter2 --name /app/token`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1812",
					Message: "`aws ssm put-parameter --type SecureString --value …` puts the plaintext in argv — `ps` / `/proc/PID/cmdline` / history / CLI debug logs can read it. Use `--cli-input-json file://…` (mode 0600) or the `file://` form for `--value`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `aws ssm put-parameter --type=SecureString --value=hunter2 --name /app/token`",
			input: `aws ssm put-parameter --type=SecureString --value=hunter2 --name /app/token`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1812",
					Message: "`aws ssm put-parameter --type SecureString --value …` puts the plaintext in argv — `ps` / `/proc/PID/cmdline` / history / CLI debug logs can read it. Use `--cli-input-json file://…` (mode 0600) or the `file://` form for `--value`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1812")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1813(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cryptsetup status cryptroot` (read only)",
			input:    `cryptsetup status cryptroot`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cryptsetup open $DEV cryptroot`",
			input:    `cryptsetup open $DEV cryptroot`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `cryptsetup luksFormat $DEV`",
			input: `cryptsetup luksFormat $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1813",
					Message: "`cryptsetup luksFormat` rewrites the LUKS header / device. Verify the target (`lsblk`), back up with `luksHeaderBackup`, and run on an unmounted volume with UPS.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `cryptsetup reencrypt $DEV`",
			input: `cryptsetup reencrypt $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1813",
					Message: "`cryptsetup reencrypt` rewrites the LUKS header / device. Verify the target (`lsblk`), back up with `luksHeaderBackup`, and run on an unmounted volume with UPS.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1813")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1814(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `dpkg -i pkg.deb` (no force)",
			input:    `dpkg -i pkg.deb`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `dpkg --force-overwrite -i pkg.deb` (specific force)",
			input:    `dpkg --force-overwrite -i pkg.deb`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `dpkg -i --force-all pkg.deb`",
			input: `dpkg -i --force-all pkg.deb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1814",
					Message: "`dpkg --force-all` enables every `--force-*` option at once — overwrite, unsigned, downgrade, essential-removal, broken-deps. Drop it and spell out only the specific `--force-<option>` you need, or fix the upstream conflict.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `apt-get -o Dpkg::Options::=--force-all install pkg`",
			input: `apt-get -o Dpkg::Options::=--force-all install pkg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1814",
					Message: "`dpkg --force-all` enables every `--force-*` option at once — overwrite, unsigned, downgrade, essential-removal, broken-deps. Drop it and spell out only the specific `--force-<option>` you need, or fix the upstream conflict.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1814")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1815(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `systemctl status NetworkManager` (read only)",
			input:    `systemctl status NetworkManager`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `systemctl restart nginx`",
			input:    `systemctl restart nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `systemctl restart NetworkManager`",
			input: `systemctl restart NetworkManager`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1815",
					Message: "`systemctl restart NetworkManager` drops every connection the manager supervises — the SSH session freezes. Use `nmcli connection reload` / `networkctl reload`, or a `systemd-run --on-active=30s` rollback.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `systemctl restart systemd-networkd.service`",
			input: `systemctl restart systemd-networkd.service`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1815",
					Message: "`systemctl restart systemd-networkd.service` drops every connection the manager supervises — the SSH session freezes. Use `nmcli connection reload` / `networkctl reload`, or a `systemd-run --on-active=30s` rollback.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1815")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1816(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `docker build -t myimage .`",
			input:    `docker build -t myimage .`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `podman ps`",
			input:    `podman ps`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `docker commit mycontainer myimage:latest`",
			input: `docker commit mycontainer myimage:latest`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1816",
					Message: "`docker commit` snapshots a running container — no Dockerfile trail, runtime env / `/tmp` scratch / shell history get baked in, and the layer metadata does not record what was installed. Build from a `Dockerfile` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `podman commit web web:snap`",
			input: `podman commit web web:snap`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1816",
					Message: "`podman commit` snapshots a running container — no Dockerfile trail, runtime env / `/tmp` scratch / shell history get baked in, and the layer metadata does not record what was installed. Build from a `Dockerfile` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1816")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1817(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `git push origin main`",
			input:    `git push origin main`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git push -u origin feature-x`",
			input:    `git push -u origin feature-x`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `git push --delete origin mybranch`",
			input: `git push --delete origin mybranch`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1817",
					Message: "`git push --delete` deletes the remote branch — open PRs are orphaned, CI targets disappear, and the last commit SHA can only come back from someone else's clone. Let the hosting platform auto-delete after merge instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `git push origin :mybranch`",
			input: `git push origin :mybranch`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1817",
					Message: "`git push origin :mybranch` deletes the remote branch — open PRs are orphaned, CI targets disappear, and the last commit SHA can only come back from someone else's clone. Let the hosting platform auto-delete after merge instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1817")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1818(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `rsync -avn --delete src/ dst/` (dry-run short)",
			input:    `rsync -avn --delete src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rsync -av --delete --dry-run src/ dst/`",
			input:    `rsync -av --delete --dry-run src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rsync -av src/ dst/` (no delete)",
			input:    `rsync -av src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `rsync -av --delete src/ dst/`",
			input: `rsync -av --delete src/ dst/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1818",
					Message: "`rsync --delete` without `--dry-run` removes anything in DST that isn't in SRC. Preview with `rsync -av --delete --dry-run SRC DST`, and pin `--max-delete=N` so an accidentally empty SRC can't cascade.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1818")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1819(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `xattr file` (read only)",
			input:    `xattr file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `xattr -d com.apple.metadata:kMDLabel_xxx file` (unrelated xattr)",
			input:    `xattr -d com.apple.metadata:kMDLabel_xxx file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `xattr -d com.apple.quarantine /Applications/MyApp.app`",
			input: `xattr -d com.apple.quarantine /Applications/MyApp.app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1819",
					Message: "`xattr -d com.apple.quarantine` / `-cr` strips the macOS Gatekeeper quarantine — the binary runs with no signature / notarization check. Verify with `codesign --verify` and `spctl --assess --type execute` before stripping.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `xattr -cr $HOME/Downloads`",
			input: `xattr -cr $HOME/Downloads`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1819",
					Message: "`xattr -d com.apple.quarantine` / `-cr` strips the macOS Gatekeeper quarantine — the binary runs with no signature / notarization check. Verify with `codesign --verify` and `spctl --assess --type execute` before stripping.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1819")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1820(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `netplan try` (auto-reverting try)",
			input:    `netplan try`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `netplan get`",
			input:    `netplan get`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `netplan apply`",
			input: `netplan apply`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1820",
					Message: "`netplan apply` commits the YAML immediately — a mistake drops the admin SSH session with no automatic rollback. Run `netplan try` first (auto-reverts if no keypress within the timeout), then `netplan apply`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `netplan apply --debug`",
			input: `netplan apply --debug`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1820",
					Message: "`netplan apply` commits the YAML immediately — a mistake drops the admin SSH session with no automatic rollback. Run `netplan try` first (auto-reverts if no keypress within the timeout), then `netplan apply`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1820")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1821(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `diskutil list` (read only)",
			input:    `diskutil list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `diskutil info $DISK`",
			input:    `diskutil info $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `diskutil eraseDisk JHFS+ NewVol $DISK`",
			input: `diskutil eraseDisk JHFS+ NewVol $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1821",
					Message: "`diskutil eraseDisk` reformats the whole disk. Resolve the target by `diskutil info -plist` / mount-point (not by index), run `diskutil list` immediately before, and require a typed confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `diskutil secureErase 0 $DISK`",
			input: `diskutil secureErase 0 $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1821",
					Message: "`diskutil secureErase` overwrites every block, no undo. Resolve the target by `diskutil info -plist` / mount-point (not by index), run `diskutil list` immediately before, and require a typed confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1821")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1822(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `csrutil status` (read only)",
			input:    `csrutil status`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `spctl --status`",
			input:    `spctl --status`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `csrutil disable`",
			input: `csrutil disable`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1822",
					Message: "`csrutil disable` disables macOS SIP / Gatekeeper / kext-consent — every malware analyst's favorite persistence primitive. Re-enable (`csrutil enable` in recovery, `spctl --master-enable`) and keep the default policy on.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `spctl kext-consent disable`",
			input: `spctl kext-consent disable`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1822",
					Message: "`spctl kext-consent disable` disables macOS SIP / Gatekeeper / kext-consent — every malware analyst's favorite persistence primitive. Re-enable (`csrutil enable` in recovery, `spctl --master-enable`) and keep the default policy on.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1822")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1823(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `keytool -list -keystore trust.jks`",
			input:    `keytool -list -keystore trust.jks`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `keytool -import -alias ca -file ca.pem -keystore trust.jks` (prompt)",
			input:    `keytool -import -alias ca -file ca.pem -keystore trust.jks`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `keytool -import -noprompt -alias ca -file ca.pem -keystore trust.jks`",
			input: `keytool -import -noprompt -alias ca -file ca.pem -keystore trust.jks`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1823",
					Message: "`keytool -import -noprompt` pins a cert to the Java trust store without a fingerprint check. Drop `-noprompt`, verify with `keytool -printcert -file CERT`, and store (alias, SHA-256) pairs in an audited inventory.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `keytool -importcert -noprompt -file ca.pem -keystore cacerts`",
			input: `keytool -importcert -noprompt -file ca.pem -keystore cacerts`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1823",
					Message: "`keytool -import -noprompt` pins a cert to the Java trust store without a fingerprint check. Drop `-noprompt`, verify with `keytool -printcert -file CERT`, and store (alias, SHA-256) pairs in an audited inventory.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1823")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1824(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl drain node-1 --ignore-daemonsets`",
			input:    `kubectl drain node-1 --ignore-daemonsets`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl cordon node-1`",
			input:    `kubectl cordon node-1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl drain node-1 --disable-eviction`",
			input: `kubectl drain node-1 --disable-eviction`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1824",
					Message: "`kubectl drain --disable-eviction` deletes pods via raw API DELETE — PodDisruptionBudgets are ignored and the workload owner's availability contract is voided. Fix the blocking PDB instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1824")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1825(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `scp src user@host:dst` (default SFTP on OpenSSH 9+)",
			input:    `scp src user@host:dst`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `scp -r dir user@host:/path`",
			input:    `scp -r dir user@host:/path`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `scp -O src user@host:dst`",
			input: `scp -O src user@host:dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1825",
					Message: "`scp -O` forces the legacy SCP wire protocol — the one exposed to filename-injection (CVE-2020-15778 class). Drop `-O` (default SFTP is safer), or use `sftp` / upgrade the remote server.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1825")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1826(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `install -m 0755 src /usr/local/bin/app`",
			input:    `install -m 0755 src /usr/local/bin/app`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `install -d /opt/app` (no mode)",
			input:    `install -d /opt/app`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `install -m 4755 …` (numeric setuid is owned by ZC1892)",
			input:    `install -m 4755 src /usr/local/bin/app`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `install -m u+s src /usr/local/bin/app`",
			input: `install -m u+s src /usr/local/bin/app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1826",
					Message: "`install -m u+s` applies a symbolic setuid/setgid bit — easy to miss in review. Use `0755` and grant narrow caps with `setcap` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `install -m ug+s src /usr/local/bin/app`",
			input: `install -m ug+s src /usr/local/bin/app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1826",
					Message: "`install -m ug+s` applies a symbolic setuid/setgid bit — easy to miss in review. Use `0755` and grant narrow caps with `setcap` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1826")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1827(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `npm deprecate mypkg@1.2.3 'use 1.2.4'`",
			input:    `npm deprecate mypkg@1.2.3 'use 1.2.4'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `npm publish`",
			input:    `npm publish`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `npm unpublish mypkg@1.2.3`",
			input: `npm unpublish mypkg@1.2.3`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1827",
					Message: "`npm unpublish` removes a published version — every downstream that pinned it fails to install on next CI run (the left-pad pattern). Use `npm deprecate PKG@VERSION 'reason'` so the version stays resolvable with a warning.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1827")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1828(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gcore --help`",
			input:    `gcore --help`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `strace ls` (trace a child, not ptrace-attach)",
			input:    `strace ls`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gcore 1234`",
			input: `gcore 1234`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1828",
					Message: "`gcore PID` attaches via ptrace — target memory, env, and syscall args are exposed. Production scripts should not run ptrace; use `coredumpctl` on a captured core instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `strace -f -p 1234`",
			input: `strace -f -p 1234`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1828",
					Message: "`strace -p PID` attaches via ptrace — target memory, env, and syscall args are exposed. Production scripts should not run ptrace; use `coredumpctl` on a captured core instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1828")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1829(t *testing.T) {
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
			name:     "valid — `nmcli connection show`",
			input:    `nmcli connection show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tailscale down`",
			input: `tailscale down`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1829",
					Message: "`tailscale down` tears down the VPN — if the SSH session rides on it, the script cuts itself off with no rollback. Schedule recovery via `systemd-run --on-active=30s`, or run from console / out-of-band.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `wg-quick down wg0`",
			input: `wg-quick down wg0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1829",
					Message: "`wg-quick down` tears down the VPN — if the SSH session rides on it, the script cuts itself off with no rollback. Schedule recovery via `systemd-run --on-active=30s`, or run from console / out-of-band.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1829")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1830(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt NOMATCH`",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt HIST_IGNORE_DUPS` (unrelated)",
			input:    `unsetopt HIST_IGNORE_DUPS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt NOMATCH`",
			input: `unsetopt NOMATCH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1830",
					Message: "`unsetopt NOMATCH` silences Zsh's unmatched-glob error — typos pass through literally. Use `*(N)` per-glob or scope inside a function with `setopt LOCAL_OPTIONS; setopt NULL_GLOB`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_NOMATCH`",
			input: `setopt NO_NOMATCH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1830",
					Message: "`setopt NO_NOMATCH` silences Zsh's unmatched-glob error — typos pass through literally. Use `*(N)` per-glob or scope inside a function with `setopt LOCAL_OPTIONS; setopt NULL_GLOB`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1830")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1831(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `systemctl reload sshd`",
			input:    `systemctl reload sshd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `systemctl status sshd`",
			input:    `systemctl status sshd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `systemctl stop sshd`",
			input: `systemctl stop sshd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1831",
					Message: "`systemctl stop sshd` blocks SSH — existing sessions survive but reconnects fail. `disable`/`mask` persist across reboots. Use `reload sshd` for config changes.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `systemctl mask ssh`",
			input: `systemctl mask ssh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1831",
					Message: "`systemctl mask ssh` blocks SSH — existing sessions survive but reconnects fail. `disable`/`mask` persist across reboots. Use `reload sshd` for config changes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1831")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1832(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `limit coredumpsize 0` (disable cores)",
			input:    `limit coredumpsize 0`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `limit stacksize unlimited` (unrelated resource)",
			input:    `limit stacksize unlimited`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `limit coredumpsize unlimited`",
			input: `limit coredumpsize unlimited`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1832",
					Message: "`limit coredumpsize unlimited` enables unbounded core dumps (Zsh-specific `limit` spelling of `ulimit -c unlimited`). A setuid crash drops its memory to disk as a world-readable file — leave the ceiling at the distro default.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unlimit coredumpsize`",
			input: `unlimit coredumpsize`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1832",
					Message: "`unlimit coredumpsize` enables unbounded core dumps (Zsh-specific `limit` spelling of `ulimit -c unlimited`). A setuid crash drops its memory to disk as a world-readable file — leave the ceiling at the distro default.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1832")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1833(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt WARN_CREATE_GLOBAL`",
			input:    `setopt WARN_CREATE_GLOBAL`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt WARN_CREATE_GLOBAL`",
			input: `unsetopt WARN_CREATE_GLOBAL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1833",
					Message: "`unsetopt WARN_CREATE_GLOBAL` silences Zsh's warning for assignments leaking out of function scope — classic caller-variable stomping. Declare `local`/`typeset`; scope with `LOCAL_OPTIONS` if you must disable.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_WARN_CREATE_GLOBAL`",
			input: `setopt NO_WARN_CREATE_GLOBAL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1833",
					Message: "`setopt NO_WARN_CREATE_GLOBAL` silences Zsh's warning for assignments leaking out of function scope — classic caller-variable stomping. Declare `local`/`typeset`; scope with `LOCAL_OPTIONS` if you must disable.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1833")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1834(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `tc qdisc add dev eth0 root netem loss 5%` (partial)",
			input:    `tc qdisc add dev eth0 root netem loss 5%`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `tc qdisc del dev eth0 root` (cleanup)",
			input:    `tc qdisc del dev eth0 root`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tc qdisc add … netem loss 100%`",
			input: `tc qdisc add dev eth0 root netem loss 100%`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1834",
					Message: "`tc qdisc add … netem loss 100%` black-holes every packet on the target interface — remote SSH dies instantly. Stage on a secondary dev or wrap in a timed recovery (`at now + N minutes … tc qdisc del …`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `tc qdisc replace … netem corrupt 100%`",
			input: `tc qdisc replace dev eth0 root netem corrupt 100%`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1834",
					Message: "`tc qdisc replace … netem corrupt 100%` black-holes every packet on the target interface — remote SSH dies instantly. Stage on a secondary dev or wrap in a timed recovery (`at now + N minutes … tc qdisc del …`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1834")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1835(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `smartctl -s on $DISK` (default)",
			input:    `smartctl -s on $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `smartctl -a $DISK` (just report)",
			input:    `smartctl -a $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `smartctl -s off $DISK`",
			input: `smartctl -s off $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1835",
					Message: "`smartctl -s off` disables the drive's SMART attribute collection — `smartctl -H` keeps reporting PASSED until the disk falls off the bus. Leave it `on` and configure `smartd.conf` for proactive alerts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1835")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1836(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `helm uninstall mychart`",
			input:    `helm uninstall mychart`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `helm uninstall mychart --keep-history` (unrelated)",
			input:    `helm uninstall mychart --keep-history`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `helm uninstall mychart --no-hooks`",
			input: `helm uninstall mychart --no-hooks`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1836",
					Message: "`helm uninstall --no-hooks` skips pre/post-delete cleanup hooks — orphaned locks, DNS, missed PVC backups. Drop the flag; fix stuck hooks via `helm.sh/hook-delete-policy`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `helm delete mychart --no-hooks` (Helm v2 spelling)",
			input: `helm delete mychart --no-hooks`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1836",
					Message: "`helm delete --no-hooks` skips pre/post-delete cleanup hooks — orphaned locks, DNS, missed PVC backups. Drop the flag; fix stuck hooks via `helm.sh/hook-delete-policy`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1836")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1837(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `chmod 660 /dev/kvm` (distro default)",
			input:    `chmod 660 /dev/kvm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `chmod 600 /dev/mem`",
			input:    `chmod 600 /dev/mem`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `chmod 666 /tmp/x` (unrelated file)",
			input:    `chmod 666 /tmp/x`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `chmod 666 /dev/kvm`",
			input: `chmod 666 /dev/kvm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1837",
					Message: "`chmod 666 /dev/kvm` grants non-owner access to a privileged kernel device — classic local-privesc vector. Use group membership or a udev rule instead of blanket chmod.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `chmod 644 /dev/mem` (world-read)",
			input: `chmod 644 /dev/mem`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1837",
					Message: "`chmod 644 /dev/mem` grants non-owner access to a privileged kernel device — classic local-privesc vector. Use group membership or a udev rule instead of blanket chmod.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `chmod a+rw /dev/port` (symbolic)",
			input: `chmod a+rw /dev/port`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1837",
					Message: "`chmod a+rw /dev/port` grants non-owner access to a privileged kernel device — classic local-privesc vector. Use group membership or a udev rule instead of blanket chmod.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1837")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1838(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt GLOB_DOTS` (explicit default)",
			input:    `unsetopt GLOB_DOTS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt GLOB_DOTS`",
			input: `setopt GLOB_DOTS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1838",
					Message: "`setopt GLOB_DOTS` makes every bare `*` also match hidden files — `rm *` quietly destroys `.git/`, `cp -r *` copies `.env`. Keep the option alone; request dotfiles per-glob with `*(D)`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_GLOB_DOTS`",
			input: `unsetopt NO_GLOB_DOTS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1838",
					Message: "`unsetopt NO_GLOB_DOTS` makes every bare `*` also match hidden files — `rm *` quietly destroys `.git/`, `cp -r *` copies `.env`. Keep the option alone; request dotfiles per-glob with `*(D)`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1838")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1839(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `timedatectl set-ntp true`",
			input:    `timedatectl set-ntp true`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `systemctl enable chronyd`",
			input:    `systemctl enable chronyd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `timedatectl set-ntp false`",
			input: `timedatectl set-ntp false`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1839",
					Message: "`timedatectl set-ntp false` turns off network time sync — clock drift breaks TLS `notBefore`/`notAfter`, Kerberos, and TOTP. Leave NTP enabled; isolate frozen clocks to namespaces/CI.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `systemctl disable systemd-timesyncd`",
			input: `systemctl disable systemd-timesyncd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1839",
					Message: "`systemctl disable systemd-timesyncd` turns off network time sync — clock drift breaks TLS `notBefore`/`notAfter`, Kerberos, and TOTP. Leave NTP enabled; isolate frozen clocks to namespaces/CI.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1839")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1840(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `openssl enc -pass env:MYPASS`",
			input:    `openssl enc -aes-256-cbc -pass env:MYPASS -in in.txt -out out.bin`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `openssl enc` without `-k`",
			input:    `openssl enc -aes-256-cbc -in in.txt -out out.bin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `openssl enc -k SECRET`",
			input: `openssl enc -aes-256-cbc -k hunter2 -in in.txt -out out.bin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1840",
					Message: "`openssl -k hunter2` embeds the password in argv — visible to `ps`, `/proc/<pid>/cmdline`, and shell history. Use `-pass env:VAR`, `-pass file:PATH`, or `-pass fd:N`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1840")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1841(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `curl --proxy-cacert /etc/ssl/proxy.pem https://api`",
			input:    `curl --proxy-cacert /etc/ssl/proxy.pem https://api`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl https://api` (no proxy flags)",
			input:    `curl https://api`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `curl --proxy-insecure https://api` (flag first)",
			input: `curl --proxy-insecure https://api`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1841",
					Message: "`curl --proxy-insecure` skips TLS verification on the proxy hop — an on-path attacker can present any cert and decrypt the tunnel (including `Authorization:` headers). Install the proxy CA and use `--proxy-cacert PATH`.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:  "invalid — `curl https://api --proxy-insecure` (flag trailing)",
			input: `curl https://api --proxy-insecure`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1841",
					Message: "`curl --proxy-insecure` skips TLS verification on the proxy hop — an on-path attacker can present any cert and decrypt the tunnel (including `Authorization:` headers). Install the proxy CA and use `--proxy-cacert PATH`.",
					Line:    1,
					Column:  19,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1841")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1842(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt CDABLE_VARS` (explicit default)",
			input:    `unsetopt CDABLE_VARS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CDABLE_VARS`",
			input: `setopt CDABLE_VARS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1842",
					Message: "`setopt CDABLE_VARS` turns a failed `cd NAME` into `cd $NAME` — a typo silently lands in whatever directory the matching variable points to. Keep this in `~/.zshrc`; in scripts use `cd \"$dir\" || exit`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_CDABLE_VARS`",
			input: `unsetopt NO_CDABLE_VARS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1842",
					Message: "`unsetopt NO_CDABLE_VARS` turns a failed `cd NAME` into `cd $NAME` — a typo silently lands in whatever directory the matching variable points to. Keep this in `~/.zshrc`; in scripts use `cd \"$dir\" || exit`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1842")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1843(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `docker run ubuntu` (no cgroup-parent)",
			input:    `docker run ubuntu`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `docker run --cgroup-parent=custom app` (non-host slice)",
			input:    `docker run --cgroup-parent=custom app`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `docker run --cgroup-parent=/ ubuntu`",
			input: `docker run --cgroup-parent=/ ubuntu`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1843",
					Message: "`docker run --cgroup-parent=/` puts the container under a host-managed slice — the engine's memory/CPU caps no longer apply. Drop the flag or pass `--memory`/`--cpus`/`--pids-limit` directly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `podman run --cgroup-parent /system.slice alpine`",
			input: `podman run --cgroup-parent /system.slice alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1843",
					Message: "`podman run --cgroup-parent=/system.slice` puts the container under a host-managed slice — the engine's memory/CPU caps no longer apply. Drop the flag or pass `--memory`/`--cpus`/`--pids-limit` directly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1843")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1844(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `logger -p auth.notice` (audit)",
			input:    `logger -p auth.notice -t scriptaudit "user added"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `logger message` (default)",
			input:    `logger "hello"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `logger -p local0.info`",
			input: `logger -p local0.info "audit: user added to wheel"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1844",
					Message: "`logger -p local0.info` writes to a `local*` facility — stock `rsyslog`/`journald` rarely collects these. Use `auth.notice`/`authpriv.info` for audit events, or `user.notice` + `-t TAG` for app logs.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `logger msg --priority=local7.notice` (trailing)",
			input: `logger "site event" --priority=local7.notice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1844",
					Message: "`logger -p local7.notice` writes to a `local*` facility — stock `rsyslog`/`journald` rarely collects these. Use `auth.notice`/`authpriv.info` for audit events, or `user.notice` + `-t TAG` for app logs.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1844")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1845(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt PATH_DIRS` (explicit default)",
			input:    `unsetopt PATH_DIRS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt PATH_DIRS`",
			input: `setopt PATH_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1845",
					Message: "`setopt PATH_DIRS` lets `subdir/cmd` fall back to a `$PATH` lookup — a missing local binary silently runs a same-named subtree elsewhere on `$PATH`. Leave the option off; call locals as `./subdir/cmd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_PATH_DIRS`",
			input: `unsetopt NO_PATH_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1845",
					Message: "`unsetopt NO_PATH_DIRS` lets `subdir/cmd` fall back to a `$PATH` lookup — a missing local binary silently runs a same-named subtree elsewhere on `$PATH`. Leave the option off; call locals as `./subdir/cmd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1845")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1846(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `certbot renew` (default)",
			input:    `certbot renew`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `certbot certificates`",
			input:    `certbot certificates`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `certbot renew --force-renewal`",
			input: `certbot renew --force-renewal`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1846",
					Message: "`certbot renew --force-renewal` reissues regardless of expiry — in a cron it burns Let's Encrypt rate limits (50 certs per domain / 7 days). Drop the flag and let the 30-day gate work.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `certbot certonly --force-renewal -d example.com`",
			input: `certbot certonly --force-renewal -d example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1846",
					Message: "`certbot certonly --force-renewal` reissues regardless of expiry — in a cron it burns Let's Encrypt rate limits (50 certs per domain / 7 days). Drop the flag and let the 30-day gate work.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1846")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1847(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt CHASE_LINKS` (explicit default)",
			input:    `unsetopt CHASE_LINKS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CHASE_LINKS`",
			input: `setopt CHASE_LINKS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1847",
					Message: "`setopt CHASE_LINKS` makes every `cd` resolve symlinks to the physical inode — `cd releases/current` lands in the release dir, breaking `..` navigation. Keep it off; use `cd -P target` one-shot when needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_CHASE_LINKS`",
			input: `unsetopt NO_CHASE_LINKS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1847",
					Message: "`unsetopt NO_CHASE_LINKS` makes every `cd` resolve symlinks to the physical inode — `cd releases/current` lands in the release dir, breaking `..` navigation. Keep it off; use `cd -P target` one-shot when needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1847")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1848(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ssh -o CheckHostIP=yes host`",
			input:    `ssh -o CheckHostIP=yes host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ssh host` (default)",
			input:    `ssh host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ssh -o CheckHostIP=no host` (split form)",
			input: `ssh -o CheckHostIP=no host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1848",
					Message: "`ssh -o CheckHostIP=no` silences the IP-mismatch warning for known hosts — a DNS-spoof + leaked host-key attack goes undetected. Leave the default, or use `HostKeyAlias` for load-balanced pools.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ssh -oCheckHostIP=no host` (attached form)",
			input: `ssh -oCheckHostIP=no host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1848",
					Message: "`ssh -o CheckHostIP=no` silences the IP-mismatch warning for known hosts — a DNS-spoof + leaked host-key attack goes undetected. Leave the default, or use `HostKeyAlias` for load-balanced pools.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1848")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1849(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt ALL_EXPORT` (explicit default)",
			input:    `unsetopt ALL_EXPORT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt ALL_EXPORT`",
			input: `setopt ALL_EXPORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1849",
					Message: "`setopt ALL_EXPORT` marks every later assignment for export — secrets like `password=...` leak into every child's env. Drop it; use explicit `export`, or scope inside a `LOCAL_OPTIONS` function.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_ALL_EXPORT`",
			input: `unsetopt NO_ALL_EXPORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1849",
					Message: "`unsetopt NO_ALL_EXPORT` marks every later assignment for export — secrets like `password=...` leak into every child's env. Drop it; use explicit `export`, or scope inside a `LOCAL_OPTIONS` function.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1849")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1850(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ssh -o LogLevel=INFO host`",
			input:    `ssh -o LogLevel=INFO host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ssh host` (default)",
			input:    `ssh host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ssh -o LogLevel=QUIET host`",
			input: `ssh -o LogLevel=QUIET host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1850",
					Message: "`ssh -o LogLevel=QUIET` silences host-key, agent-forward, and canonical-hostname warnings — a MITM event produces no stderr. Keep the default level; capture stderr to a log if you need it clean.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ssh -oLogLevel=fatal host` (attached)",
			input: `ssh -oLogLevel=fatal host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1850",
					Message: "`ssh -o LogLevel=QUIET` silences host-key, agent-forward, and canonical-hostname warnings — a MITM event produces no stderr. Keep the default level; capture stderr to a log if you need it clean.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1850")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1851(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt FUNCTION_ARGZERO` (explicit default)",
			input:    `setopt FUNCTION_ARGZERO`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt FUNCTION_ARGZERO`",
			input: `unsetopt FUNCTION_ARGZERO`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1851",
					Message: "`unsetopt FUNCTION_ARGZERO` makes `$0` inside functions point at the outer script — breaks `log \"$0: ...\"` helpers and `case $0` dispatchers. Keep the option on; reach the script name explicitly via `$ZSH_ARGZERO` when needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_FUNCTION_ARGZERO`",
			input: `setopt NO_FUNCTION_ARGZERO`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1851",
					Message: "`setopt NO_FUNCTION_ARGZERO` makes `$0` inside functions point at the outer script — breaks `log \"$0: ...\"` helpers and `case $0` dispatchers. Keep the option on; reach the script name explicitly via `$ZSH_ARGZERO` when needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1851")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1852(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `firewall-cmd --panic-off foo`",
			input:    `firewall-cmd --panic-off foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `firewall-cmd --reload`",
			input:    `firewall-cmd --reload foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `firewall-cmd --panic-on >/dev/null` (mangled name)",
			input: `firewall-cmd --panic-on foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1852",
					Message: "`firewall-cmd --panic-on` drops every packet regardless of zone — an SSH-run call loses the session instantly. Use targeted zone rules; if you really need panic mode, gate behind `at now + N minutes … firewall-cmd --panic-off`.",
					Line:    1,
					Column:  15,
				},
			},
		},
		{
			name:  "invalid — `firewall-cmd \"\" --panic-on` (trailing flag)",
			input: `firewall-cmd "" --panic-on`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1852",
					Message: "`firewall-cmd --panic-on` drops every packet regardless of zone — an SSH-run call loses the session instantly. Use targeted zone rules; if you really need panic mode, gate behind `at now + N minutes … firewall-cmd --panic-off`.",
					Line:    1,
					Column:  18,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1852")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1853(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt MARK_DIRS` (explicit default)",
			input:    `unsetopt MARK_DIRS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt MARK_DIRS`",
			input: `setopt MARK_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1853",
					Message: "`setopt MARK_DIRS` appends a trailing `/` to every glob-matched directory — `[[ -f \"$f\" ]]` and `rm -f *` start skipping, hash maps keyed on basenames double up. Keep the option off; use `*(/)` when you need dirs only.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_MARK_DIRS`",
			input: `unsetopt NO_MARK_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1853",
					Message: "`unsetopt NO_MARK_DIRS` appends a trailing `/` to every glob-matched directory — `[[ -f \"$f\" ]]` and `rm -f *` start skipping, hash maps keyed on basenames double up. Keep the option off; use `*(/)` when you need dirs only.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1853")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1854(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `yum-config-manager --add-repo https://…` (TLS)",
			input:    `yum-config-manager --add-repo https://mirror.example/app.repo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `zypper addrepo https://…` (TLS)",
			input:    `zypper addrepo https://mirror.example/app app`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `yum-config-manager --add-repo http://…`",
			input: `yum-config-manager --add-repo http://mirror.example/app.repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1854",
					Message: "`yum-config-manager --add-repo http://mirror.example/app.repo` registers a plaintext repo — on-path attacker can substitute packages and strip GPG-check directives. Use `https://` and pin `gpgkey=file://` in the `.repo`.",
					Line:    1,
					Column:  21,
				},
			},
		},
		{
			name:  "invalid — `dnf config-manager --add-repo http://…`",
			input: `dnf config-manager --add-repo http://mirror.example/app.repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1854",
					Message: "`dnf config-manager --add-repo http://mirror.example/app.repo` registers a plaintext repo — on-path attacker can substitute packages and strip GPG-check directives. Use `https://` and pin `gpgkey=file://` in the `.repo`.",
					Line:    1,
					Column:  21,
				},
			},
		},
		{
			name:  "invalid — `zypper addrepo http://…`",
			input: `zypper addrepo http://mirror.example/app app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1854",
					Message: "`zypper addrepo http://mirror.example/app` registers a plaintext repo — on-path attacker can substitute packages and strip GPG-check directives. Use `https://` and pin `gpgkey=file://` in the `.repo`.",
					Line:    1,
					Column:  8,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1854")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1855(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `echo ${(k)groups}` (Zsh-native)",
			input:    `echo ${(k)groups}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `echo GROUPSIZE` (unrelated literal)",
			input:    `echo GROUPSIZE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `echo $GROUPS`",
			input: `echo $GROUPS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1855",
					Message: "`$GROUPS` is a Bash-only array — Zsh populates `$groups` (associative name→GID) instead. Iterate `${(k)groups}` for names or `${(v)groups}` for GIDs, or fall back to `id -Gn`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `printf '%s\\n' \"${GROUPS[@]}\"`",
			input: `printf '%s\n' "${GROUPS[@]}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1855",
					Message: "`$GROUPS` is a Bash-only array — Zsh populates `$groups` (associative name→GID) instead. Iterate `${(k)groups}` for names or `${(v)groups}` for GIDs, or fall back to `id -Gn`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1855")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1856(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unset arr` (delete whole variable)",
			input:    `unset arr`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unset FOO` (scalar)",
			input:    `unset FOO`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unset arr[0]`",
			input: `unset arr[0]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1856",
					Message: "`unset (arr[0])` is a Bash idiom — in Zsh it tries to unset a parameter literally named `(arr[0])` and leaves the array untouched. Use `arr[N]=()` or rebuild with `arr=(\"${(@)arr:#pattern}\")`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unset myarray[3]`",
			input: `unset myarray[3]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1856",
					Message: "`unset (myarray[3])` is a Bash idiom — in Zsh it tries to unset a parameter literally named `(myarray[3])` and leaves the array untouched. Use `arr[N]=()` or rebuild with `arr=(\"${(@)arr:#pattern}\")`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1856")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1857(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cloud-init init` (boot-time init)",
			input:    `cloud-init init`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cloud-init status`",
			input:    `cloud-init status`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `cloud-init clean`",
			input: `cloud-init clean`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1857",
					Message: "`cloud-init clean` wipes `/var/lib/cloud/` boot state — the next reboot re-runs the user-data and overwrites operator changes (SSH host keys, hostname, `/etc/fstab`). Run interactively only.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `cloud-init clean --logs --reboot`",
			input: `cloud-init clean --logs --reboot`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1857",
					Message: "`cloud-init clean` wipes `/var/lib/cloud/` boot state — the next reboot re-runs the user-data and overwrites operator changes (SSH host keys, hostname, `/etc/fstab`). Run interactively only.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1857")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1858(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ssh -c aes256-gcm@openssh.com host`",
			input:    `ssh -c aes256-gcm@openssh.com host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ssh host` (default ciphers)",
			input:    `ssh host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ssh -c 3des-cbc host`",
			input: `ssh -c 3des-cbc host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1858",
					Message: "`ssh ... 3des-cbc` forces a weak cipher with known plaintext-recovery / IV-reuse attacks. Leave cipher selection to OpenSSH defaults; if a legacy peer needs it, scope inside a `Host` block in `~/.ssh/config`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ssh -o Ciphers=arcfour,aes256-ctr host`",
			input: `ssh -o Ciphers=arcfour,aes256-ctr host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1858",
					Message: "`ssh ... arcfour` forces a weak cipher with known plaintext-recovery / IV-reuse attacks. Leave cipher selection to OpenSSH defaults; if a legacy peer needs it, scope inside a `Host` block in `~/.ssh/config`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1858")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1859(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt MULTIOS` (explicit default)",
			input:    `setopt MULTIOS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt MULTIOS`",
			input: `unsetopt MULTIOS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1859",
					Message: "`unsetopt MULTIOS` reverts to POSIX single-output redirection — `cmd >a >b` silently drops `a`, log collectors stop receiving new lines. Keep the option on; scope inside a `LOCAL_OPTIONS` function if one line really needs POSIX.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_MULTIOS`",
			input: `setopt NO_MULTIOS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1859",
					Message: "`setopt NO_MULTIOS` reverts to POSIX single-output redirection — `cmd >a >b` silently drops `a`, log collectors stop receiving new lines. Keep the option on; scope inside a `LOCAL_OPTIONS` function if one line really needs POSIX.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1859")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1860(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `hostnamectl status`",
			input:    `hostnamectl status`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `hostname -f` (read-only query)",
			input:    `hostname -f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `hostnamectl set-hostname worker-42`",
			input: `hostnamectl set-hostname worker-42`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1860",
					Message: "`hostnamectl set-hostname worker-42` updates the kernel hostname live, but running services keep the old `gethostname()` — syslog tags, Prometheus labels, TLS SANs stay stale. Apply at provisioning or reboot.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `hostname worker-42`",
			input: `hostname worker-42`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1860",
					Message: "`hostname worker-42` updates the kernel hostname live, but running services keep the old `gethostname()` — syslog tags, Prometheus labels, TLS SANs stay stale. Apply at provisioning or reboot.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1860")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1861(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt OCTAL_ZEROES` (explicit default)",
			input:    `unsetopt OCTAL_ZEROES`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt OCTAL_ZEROES`",
			input: `setopt OCTAL_ZEROES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1861",
					Message: "`setopt OCTAL_ZEROES` reinterprets leading-zero integers as octal — `(( n = 0100 ))` assigns 64 instead of 100, breaking timestamp, phone-prefix, and mode parsing. Keep the option off; use `8#100` when you want explicit octal.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_OCTAL_ZEROES`",
			input: `unsetopt NO_OCTAL_ZEROES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1861",
					Message: "`unsetopt NO_OCTAL_ZEROES` reinterprets leading-zero integers as octal — `(( n = 0100 ))` assigns 64 instead of 100, breaking timestamp, phone-prefix, and mode parsing. Keep the option off; use `8#100` when you want explicit octal.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1861")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1862(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ssh-keygen -t ed25519 -f id_host`",
			input:    `ssh-keygen -t ed25519 -f id_host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ssh-keygen -lf id_host.pub`",
			input:    `ssh-keygen -lf id_host.pub`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ssh-keygen -R server.example`",
			input: `ssh-keygen -R server.example`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1862",
					Message: "`ssh-keygen -R server.example` deletes a known-hosts entry — the next `ssh` silently re-trusts whatever key the network returns. Fetch the new fingerprint out-of-band and verify before re-adding.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1862")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1863(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt CASE_GLOB` (explicit default)",
			input:    `setopt CASE_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt CASE_GLOB`",
			input: `unsetopt CASE_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1863",
					Message: "`unsetopt CASE_GLOB` flips every later glob to case-insensitive — `rm *.log` sweeps `APP.LOG`, dispatchers keyed on case collisions. Keep the option on; use `(#i)pattern` per-glob when you need case-folding.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_CASE_GLOB`",
			input: `setopt NO_CASE_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1863",
					Message: "`setopt NO_CASE_GLOB` flips every later glob to case-insensitive — `rm *.log` sweeps `APP.LOG`, dispatchers keyed on case collisions. Keep the option on; use `(#i)pattern` per-glob when you need case-folding.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1863")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1864(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mount -o remount,noexec /tmp` (tightening)",
			input:    `mount -o remount,noexec /tmp`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mount -o remount,rw /` (unrelated)",
			input:    `mount -o remount,rw /`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mount -o remount,exec /tmp`",
			input: `mount -o remount,exec /tmp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1864",
					Message: "`mount -o remount,exec` re-enables `exec` on a `noexec`/`nosuid`/`nodev`-hardened mount — dropped payloads suddenly execute. Pair with a `trap ... EXIT` that restores the original flags or skip the remount.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `mount -o remount,rw,suid /var`",
			input: `mount -o remount,rw,suid /var`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1864",
					Message: "`mount -o remount,rw,suid` re-enables `suid` on a `noexec`/`nosuid`/`nodev`-hardened mount — dropped payloads suddenly execute. Pair with a `trap ... EXIT` that restores the original flags or skip the remount.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1864")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1865(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt CASE_MATCH` (explicit default)",
			input:    `setopt CASE_MATCH`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt CASE_MATCH`",
			input: `unsetopt CASE_MATCH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1865",
					Message: "`unsetopt CASE_MATCH` flips every `[[ =~ ]]` / `[[ == pat ]]` to case-insensitive — `Admin` matches `ADMIN`, dispatchers collide. Keep it on; scope per-line with `(#i)pattern`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_CASE_MATCH`",
			input: `setopt NO_CASE_MATCH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1865",
					Message: "`setopt NO_CASE_MATCH` flips every `[[ =~ ]]` / `[[ == pat ]]` to case-insensitive — `Admin` matches `ADMIN`, dispatchers collide. Keep it on; scope per-line with `(#i)pattern`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1865")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1866(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `docker exec web bash`",
			input:    `docker exec web bash`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `docker exec -u app web bash`",
			input:    `docker exec -u app web bash`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `docker exec -u 0 web bash`",
			input: `docker exec -u 0 web bash`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1866",
					Message: "`docker exec -u 0` drops a root shell — bypasses the image's non-root `USER` and, without userns remap, equals host root. Keep execs as the container user.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `podman exec --user=root app sh`",
			input: `podman exec --user=root app sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1866",
					Message: "`podman exec -u root` drops a root shell — bypasses the image's non-root `USER` and, without userns remap, equals host root. Keep execs as the container user.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1866")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1867(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt GLOB` (explicit default)",
			input:    `setopt GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt GLOB`",
			input: `unsetopt GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1867",
					Message: "`unsetopt GLOB` disables glob expansion — `rm *.log` chases the literal `*.log`, `for f in *.txt` loops once. Quote specific args or scope with `LOCAL_OPTIONS` inside a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_GLOB`",
			input: `setopt NO_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1867",
					Message: "`setopt NO_GLOB` disables glob expansion — `rm *.log` chases the literal `*.log`, `for f in *.txt` loops once. Quote specific args or scope with `LOCAL_OPTIONS` inside a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1867")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1868(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gcloud config set compute/zone us-central1-a`",
			input:    `gcloud config set compute/zone us-central1-a`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gcloud config set auth/disable_ssl_validation false`",
			input:    `gcloud config set auth/disable_ssl_validation false`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gcloud config set auth/disable_ssl_validation true`",
			input: `gcloud config set auth/disable_ssl_validation true`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1868",
					Message: "`gcloud config set auth/disable_ssl_validation true` turns off TLS for every later `gcloud` call — service-account tokens and deploys become interceptable. Unset it; pin custom CAs via `core/custom_ca_certs_file`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1868")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1869(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt RC_EXPAND_PARAM` (explicit default)",
			input:    `unsetopt RC_EXPAND_PARAM`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt RC_EXPAND_PARAM`",
			input: `setopt RC_EXPAND_PARAM`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1869",
					Message: "`setopt RC_EXPAND_PARAM` distributes literal prefix/suffix across every array element — `cp src/${arr[@]}.bak dst` silently rewrites as `cp src/a.bak src/b.bak dst`. Keep it off; opt in per-use with `${^arr}`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_RC_EXPAND_PARAM`",
			input: `unsetopt NO_RC_EXPAND_PARAM`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1869",
					Message: "`unsetopt NO_RC_EXPAND_PARAM` distributes literal prefix/suffix across every array element — `cp src/${arr[@]}.bak dst` silently rewrites as `cp src/a.bak src/b.bak dst`. Keep it off; opt in per-use with `${^arr}`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1869")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1870(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt GLOB_ASSIGN` (explicit default)",
			input:    `unsetopt GLOB_ASSIGN`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt GLOB_ASSIGN`",
			input: `setopt GLOB_ASSIGN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1870",
					Message: "`setopt GLOB_ASSIGN` expands glob patterns on the RHS of `var=` — `logs=*.log` silently captures the first match, `cert=~/secrets/*` picks up attacker drops. Keep it off; use explicit `arr=( *.log )`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_GLOB_ASSIGN`",
			input: `unsetopt NO_GLOB_ASSIGN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1870",
					Message: "`unsetopt NO_GLOB_ASSIGN` expands glob patterns on the RHS of `var=` — `logs=*.log` silently captures the first match, `cert=~/secrets/*` picks up attacker drops. Keep it off; use explicit `arr=( *.log )`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1870")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1871(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt IGNORE_BRACES` (explicit default)",
			input:    `unsetopt IGNORE_BRACES`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt IGNORE_BRACES`",
			input: `setopt IGNORE_BRACES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1871",
					Message: "`setopt IGNORE_BRACES` disables brace expansion — `for i in {1..10}` loops once over the literal token, `cp app.{conf,bak}` fails ENOENT. Keep the option off; quote the specific argument if you need a literal brace string.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_IGNORE_BRACES`",
			input: `unsetopt NO_IGNORE_BRACES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1871",
					Message: "`unsetopt NO_IGNORE_BRACES` disables brace expansion — `for i in {1..10}` loops once over the literal token, `cp app.{conf,bak}` fails ENOENT. Keep the option off; quote the specific argument if you need a literal brace string.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1871")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1872(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `badblocks -n $DISK`",
			input:    `badblocks -n $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `badblocks $DISK` (read-only)",
			input:    `badblocks $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `badblocks -w $DISK`",
			input: `badblocks -w $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1872",
					Message: "`badblocks -w` overwrites every sector of the target device — silent data wipe on a populated disk. Use `-n` (non-destructive) or gate destructive runs behind a confirmation and a fresh partition.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `badblocks -wsv $DISK` (combined)",
			input: `badblocks -wsv $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1872",
					Message: "`badblocks -w` overwrites every sector of the target device — silent data wipe on a populated disk. Use `-n` (non-destructive) or gate destructive runs behind a confirmation and a fresh partition.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1872")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1873(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt ERR_RETURN` (explicit default)",
			input:    `unsetopt ERR_RETURN`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt ERR_RETURN`",
			input: `setopt ERR_RETURN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1873",
					Message: "`setopt ERR_RETURN` returns from every function on first non-zero exit — probing helpers (`test -f`, `grep -q`) bail before the fallback branch. Scope inside a `LOCAL_OPTIONS` function if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_ERR_RETURN`",
			input: `unsetopt NO_ERR_RETURN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1873",
					Message: "`unsetopt NO_ERR_RETURN` returns from every function on first non-zero exit — probing helpers (`test -f`, `grep -q`) bail before the fallback branch. Scope inside a `LOCAL_OPTIONS` function if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1873")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1874(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sshuttle -r user@host 10.0.0.0/8`",
			input:    `sshuttle -r user@host 10.0.0.0/8`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sshuttle -r user@host 192.168.1.0/24`",
			input:    `sshuttle -r user@host 192.168.1.0/24`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `sshuttle -r user@host 0/0`",
			input: `sshuttle -r user@host 0/0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1874",
					Message: "`sshuttle ... 0/0` routes every outbound packet through the jump host — a compromise of `user@host` sees the whole fleet's traffic. Scope to the subnets you actually need.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `sshuttle -r user@host 0.0.0.0/0`",
			input: `sshuttle -r user@host 0.0.0.0/0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1874",
					Message: "`sshuttle ... 0.0.0.0/0` routes every outbound packet through the jump host — a compromise of `user@host` sees the whole fleet's traffic. Scope to the subnets you actually need.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1874")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1875(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt RC_QUOTES` (explicit default)",
			input:    `unsetopt RC_QUOTES`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt RC_QUOTES`",
			input: `setopt RC_QUOTES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1875",
					Message: "`setopt RC_QUOTES` reinterprets `''` inside single quotes as a literal apostrophe — `'it''s'` flips from `its` to `it's`, breaking tokens and SQL. Use double quotes or `\\'` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_RC_QUOTES`",
			input: `unsetopt NO_RC_QUOTES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1875",
					Message: "`unsetopt NO_RC_QUOTES` reinterprets `''` inside single quotes as a literal apostrophe — `'it''s'` flips from `its` to `it's`, breaking tokens and SQL. Use double quotes or `\\'` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1875")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1876(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cargo publish`",
			input:    `cargo publish`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cargo publish --dry-run`",
			input:    `cargo publish --dry-run`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `cargo publish --allow-dirty`",
			input: `cargo publish --allow-dirty`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1876",
					Message: "`cargo publish --allow-dirty` uploads a tarball snapshot of the dirty working tree — debug prints and local-only patches end up on crates.io for a version that cannot be replaced. Commit first; `--dry-run` to rehearse.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1876")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1877(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt SHORT_LOOPS` (explicit default)",
			input:    `setopt SHORT_LOOPS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt SHORT_LOOPS`",
			input: `unsetopt SHORT_LOOPS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1877",
					Message: "`unsetopt SHORT_LOOPS` disables short-form loops — `for f in *.log; print $f` raises a parse error. Keep the option on; scope inside a function with `LOCAL_OPTIONS` if POSIX-strict parsing is really needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_SHORT_LOOPS`",
			input: `setopt NO_SHORT_LOOPS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1877",
					Message: "`setopt NO_SHORT_LOOPS` disables short-form loops — `for f in *.log; print $f` raises a parse error. Keep the option on; scope inside a function with `LOCAL_OPTIONS` if POSIX-strict parsing is really needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1877")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1878(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl apply -f manifest.yaml`",
			input:    `kubectl apply -f manifest.yaml`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl apply --server-side -f manifest.yaml`",
			input:    `kubectl apply --server-side -f manifest.yaml`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl apply --server-side --force-conflicts -f manifest.yaml`",
			input: `kubectl apply --server-side --force-conflicts -f manifest.yaml`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1878",
					Message: "`kubectl apply --force-conflicts` grabs ownership of every conflicting field from other controllers (HPA, cert-manager, sidecar injectors). Resolve the conflict instead — drop the disputed fields or hand off via managed-field edit.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1878")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1879(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt BAD_PATTERN` (explicit default)",
			input:    `setopt BAD_PATTERN`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt BAD_PATTERN`",
			input: `unsetopt BAD_PATTERN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1879",
					Message: "`unsetopt BAD_PATTERN` silences `bad pattern` errors — `rm [abc` tries to remove a literal `[abc`, broken `case` arms stop firing. Keep the option on; quote one-off patterns or scope with `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_BAD_PATTERN`",
			input: `setopt NO_BAD_PATTERN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1879",
					Message: "`setopt NO_BAD_PATTERN` silences `bad pattern` errors — `rm [abc` tries to remove a literal `[abc`, broken `case` arms stop firing. Keep the option on; quote one-off patterns or scope with `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1879")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1880(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl annotate pod/foo key=val`",
			input:    `kubectl annotate pod/foo key=val`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl label pod/foo role=app`",
			input:    `kubectl label pod/foo role=app`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl annotate pod/foo --overwrite key=val`",
			input: `kubectl annotate pod/foo --overwrite key=val`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1880",
					Message: "`kubectl annotate --overwrite` silently replaces an existing controller signal — cert-manager, external-dns, HPA watchers reconcile on the new value. Inspect first; drop `--overwrite` so conflicts error.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `kubectl label node/bar --overwrite role=worker`",
			input: `kubectl label node/bar --overwrite role=worker`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1880",
					Message: "`kubectl label --overwrite` silently replaces an existing controller signal — cert-manager, external-dns, HPA watchers reconcile on the new value. Inspect first; drop `--overwrite` so conflicts error.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1880")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1881(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt MULTIBYTE` (explicit default)",
			input:    `setopt MULTIBYTE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt MULTIBYTE`",
			input: `unsetopt MULTIBYTE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1881",
					Message: "`unsetopt MULTIBYTE` flips every string op to per-byte math — `${#str}` on an emoji returns 4, substrings slice mid-codepoint, `[[ =~ ]]` Unicode ranges break. Keep the option on; byte-count with `wc -c <<< $str` when truly needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_MULTIBYTE`",
			input: `setopt NO_MULTIBYTE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1881",
					Message: "`setopt NO_MULTIBYTE` flips every string op to per-byte math — `${#str}` on an emoji returns 4, substrings slice mid-codepoint, `[[ =~ ]]` Unicode ranges break. Keep the option on; byte-count with `wc -c <<< $str` when truly needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1881")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1882(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sudo /usr/local/bin/setup.sh`",
			input:    `sudo /usr/local/bin/setup.sh`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sudo -i /usr/local/bin/setup.sh`",
			input:    `sudo -i /usr/local/bin/setup.sh`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sudo su -c \"cmd\"`",
			input:    `sudo su -c "cmd"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `sudo -s`",
			input: `sudo -s`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1882",
					Message: "`sudo -s` spawns an interactive root shell — in a script either hangs on stdin or drains the rest of the file into root's shell. Pass the command to sudo: `sudo /path/to/cmd arg …`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `sudo su -`",
			input: `sudo su -`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1882",
					Message: "`sudo su` spawns an interactive root shell — in a script either hangs on stdin or drains the rest of the file into root's shell. Pass the command to sudo: `sudo /path/to/cmd arg …`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `sudo bash`",
			input: `sudo bash`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1882",
					Message: "`sudo bash` spawns an interactive root shell — in a script either hangs on stdin or drains the rest of the file into root's shell. Pass the command to sudo: `sudo /path/to/cmd arg …`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1882")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1883(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt PATH_SCRIPT` (explicit default)",
			input:    `unsetopt PATH_SCRIPT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt PATH_SCRIPT`",
			input: `setopt PATH_SCRIPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1883",
					Message: "`setopt PATH_SCRIPT` lets `.`/`source` fall back to `$PATH` when a literal path misses — a dropper in `~/bin` or `./` runs inside the current shell with every exported secret. Keep the option off; always use explicit paths.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_PATH_SCRIPT`",
			input: `unsetopt NO_PATH_SCRIPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1883",
					Message: "`unsetopt NO_PATH_SCRIPT` lets `.`/`source` fall back to `$PATH` when a literal path misses — a dropper in `~/bin` or `./` runs inside the current shell with every exported secret. Keep the option off; always use explicit paths.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1883")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1884(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `curl https://api.example/public`",
			input:    `curl https://api.example/public`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl -H \"Authorization: Bearer $T\" https://api.example/private`",
			input:    `curl -H "Authorization: Bearer $T" https://api.example/private`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `curl https://api/thing?apikey=abc`",
			input: `curl https://api.example/thing?apikey=abc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1884",
					Message: "`curl https://api.example/thing?apikey=abc` carries `apikey...` in the URL query — logged by every proxy, CDN, and server access log along the path. Move credentials to `-H \"Authorization: Bearer \"$TOKEN\"` or a POST body.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `curl -X POST https://api.example/auth?token=xyz`",
			input: `curl -X POST https://api.example/auth?token=xyz`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1884",
					Message: "`curl https://api.example/auth?token=xyz` carries `token...` in the URL query — logged by every proxy, CDN, and server access log along the path. Move credentials to `-H \"Authorization: Bearer \"$TOKEN\"` or a POST body.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1884")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1885(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt CSH_NULL_GLOB` (explicit default)",
			input:    `unsetopt CSH_NULL_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CSH_NULL_GLOB`",
			input: `setopt CSH_NULL_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1885",
					Message: "`setopt CSH_NULL_GLOB` silently discards unmatched globs in a list when any sibling matches — `rm *.lg *.bak` deletes the `.bak` files and hides the typo. Keep the option off; use `*(N)` per-glob.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_CSH_NULL_GLOB`",
			input: `unsetopt NO_CSH_NULL_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1885",
					Message: "`unsetopt NO_CSH_NULL_GLOB` silently discards unmatched globs in a list when any sibling matches — `rm *.lg *.bak` deletes the `.bak` files and hides the typo. Keep the option off; use `*(N)` per-glob.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1885")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1886(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cp /tmp/app.tar /opt/app/`",
			input:    `cp /tmp/app.tar /opt/app/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `tee /var/log/install.log`",
			input:    `tee /var/log/install.log`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tee /etc/profile`",
			input: `tee /etc/profile`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1886",
					Message: "`tee ... /etc/profile` writes a shell-init file sourced by every interactive shell — persistent foothold for the next root login. Stage a temp file, validate, and `install -m 644` atomically.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `cp new.sh /etc/profile.d/custom.sh`",
			input: `cp new.sh /etc/profile.d/custom.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1886",
					Message: "`cp ... /etc/profile.d/custom.sh` writes a shell-init file sourced by every interactive shell — persistent foothold for the next root login. Stage a temp file, validate, and `install -m 644` atomically.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1886")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1887(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt POSIX_TRAPS` (explicit default)",
			input:    `unsetopt POSIX_TRAPS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt POSIX_TRAPS`",
			input: `setopt POSIX_TRAPS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1887",
					Message: "`setopt POSIX_TRAPS` flips `trap ... EXIT` inside functions from function-return to shell-exit scope — per-call cleanup leaks across the whole shell, TRAPZERR helpers stop firing. Keep the option off.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_POSIX_TRAPS`",
			input: `unsetopt NO_POSIX_TRAPS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1887",
					Message: "`unsetopt NO_POSIX_TRAPS` flips `trap ... EXIT` inside functions from function-return to shell-exit scope — per-call cleanup leaks across the whole shell, TRAPZERR helpers stop firing. Keep the option off.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1887")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1888(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `aws iam list-access-keys`",
			input:    `aws iam list-access-keys`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `aws iam get-role --role-name foo`",
			input:    `aws iam get-role --role-name foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `aws iam create-access-key --user-name ci-bot`",
			input: `aws iam create-access-key --user-name ci-bot`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1888",
					Message: "`aws iam create-access-key` mints a long-lived `AKIA.../secret` — prefer short-lived creds via instance profiles, IRSA, Lambda roles, or OIDC federation. If static keys are unavoidable, store in Secrets Manager and rotate.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1888")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1889(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `skopeo copy docker://a docker://b`",
			input:    `skopeo copy docker://a docker://b`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `skopeo copy --src-tls-verify=true docker://a docker://b`",
			input:    `skopeo copy --src-tls-verify=true docker://a docker://b`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `skopeo copy --src-tls-verify=false docker://a docker://b`",
			input: `skopeo copy --src-tls-verify=false docker://a docker://b`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1889",
					Message: "`skopeo --src-tls-verify=false` disables TLS verification on image copy — on-path attacker can substitute a malicious manifest. Pin a private CA with `--src-cert-dir`/`--dest-cert-dir` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `skopeo copy --dest-tls-verify=false docker://a docker://b`",
			input: `skopeo copy --dest-tls-verify=false docker://a docker://b`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1889",
					Message: "`skopeo --dest-tls-verify=false` disables TLS verification on image copy — on-path attacker can substitute a malicious manifest. Pin a private CA with `--src-cert-dir`/`--dest-cert-dir` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1889")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1890(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kadmin -p admin/admin -k -t /etc/krb5.keytab`",
			input:    `kadmin -p admin/admin -k -t /etc/krb5.keytab`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kinit admin/admin`",
			input:    `kinit admin/admin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kadmin -p admin/admin -w hunter2`",
			input: `kadmin -p admin/admin -w hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1890",
					Message: "`kadmin -w hunter2` embeds the Kerberos admin password in argv — visible to `ps`, `/proc`, shell history. Use `-k -t /etc/krb5.keytab` (keytab root-only) or pipe the password on stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `kadmin.local -w hunter2 addprinc user`",
			input: `kadmin.local -w hunter2 addprinc user`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1890",
					Message: "`kadmin.local -w hunter2` embeds the Kerberos admin password in argv — visible to `ps`, `/proc`, shell history. Use `-k -t /etc/krb5.keytab` (keytab root-only) or pipe the password on stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1890")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1891(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl config view`",
			input:    `kubectl config view`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl config view -o jsonpath='{.current-context}'`",
			input:    `kubectl config view -o jsonpath='{.current-context}'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl config view --raw`",
			input: `kubectl config view --raw`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1891",
					Message: "`kubectl config view --raw` prints the full kubeconfig including client-certificate/key-data and bearer tokens — any script-captured stdout exfiltrates the creds. Emit the specific field with `-o jsonpath='…'`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `kubectl config view -R`",
			input: `kubectl config view -R`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1891",
					Message: "`kubectl config view --raw` prints the full kubeconfig including client-certificate/key-data and bearer tokens — any script-captured stdout exfiltrates the creds. Emit the specific field with `-o jsonpath='…'`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1891")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1892(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `install -m 0755 foo /usr/local/bin/foo`",
			input:    `install -m 0755 foo /usr/local/bin/foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `install -m 0644 foo.conf /etc/foo.conf`",
			input:    `install -m 0644 foo.conf /etc/foo.conf`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `install -m 4755 foo /usr/local/bin/foo` (setuid)",
			input: `install -m 4755 foo /usr/local/bin/foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1892",
					Message: "`install -m 4755` sets setuid/setgid bits at install time — every execution becomes a privesc vector. Use `0755` and grant narrow caps with `setcap` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `install -m 2755 foo /usr/local/bin/foo` (setgid)",
			input: `install -m 2755 foo /usr/local/bin/foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1892",
					Message: "`install -m 2755` sets setuid/setgid bits at install time — every execution becomes a privesc vector. Use `0755` and grant narrow caps with `setcap` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1892")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1893(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt BARE_GLOB_QUAL` (explicit default)",
			input:    `setopt BARE_GLOB_QUAL`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt BARE_GLOB_QUAL`",
			input: `unsetopt BARE_GLOB_QUAL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1893",
					Message: "`unsetopt BARE_GLOB_QUAL` disables `*(qualifier)` syntax — `*(N)` stops being null-glob and becomes an alternation, so null-glob idioms silently break. Keep the option on; scope inside a `LOCAL_OPTIONS` function.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_BARE_GLOB_QUAL`",
			input: `setopt NO_BARE_GLOB_QUAL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1893",
					Message: "`setopt NO_BARE_GLOB_QUAL` disables `*(qualifier)` syntax — `*(N)` stops being null-glob and becomes an alternation, so null-glob idioms silently break. Keep the option on; scope inside a `LOCAL_OPTIONS` function.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1893")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1894(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `conntrack -L` (list)",
			input:    `conntrack -L`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `conntrack -D -s 10.0.0.5` (narrow delete)",
			input:    `conntrack -D -s 10.0.0.5`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `conntrack -F`",
			input: `conntrack -F`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1894",
					Message: "`conntrack -F` wipes every tracked flow — stateful `ctstate ESTABLISHED` allowances drop, running SSH sessions lose their entry. Gate with `at now + N min` or narrow to one flow with `conntrack -D -s <ip>`.",
					Line:    1,
					Column:  11,
				},
			},
		},
		{
			name:  "invalid — `conntrack --flush conntrack`",
			input: `conntrack --flush conntrack`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1894",
					Message: "`conntrack -F` wipes every tracked flow — stateful `ctstate ESTABLISHED` allowances drop, running SSH sessions lose their entry. Gate with `at now + N min` or narrow to one flow with `conntrack -D -s <ip>`.",
					Line:    1,
					Column:  12,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1894")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1895(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt NUMERIC_GLOB_SORT` (explicit default)",
			input:    `unsetopt NUMERIC_GLOB_SORT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt NUMERIC_GLOB_SORT`",
			input: `setopt NUMERIC_GLOB_SORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1895",
					Message: "`setopt NUMERIC_GLOB_SORT` switches every later glob to numeric sort — log rotations sorted on numeric suffixes silently shuffle. Keep it off; use the per-glob `*(n)` qualifier when needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_NUMERIC_GLOB_SORT`",
			input: `unsetopt NO_NUMERIC_GLOB_SORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1895",
					Message: "`unsetopt NO_NUMERIC_GLOB_SORT` switches every later glob to numeric sort — log rotations sorted on numeric suffixes silently shuffle. Keep it off; use the per-glob `*(n)` qualifier when needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1895")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1896(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `docker run -v /etc/app:/app/etc ubuntu`",
			input:    `docker run -v /etc/app:/app/etc ubuntu`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `docker run -v /home/user:/work ubuntu`",
			input:    `docker run -v /home/user:/work ubuntu`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `docker run -v /proc:/host/proc:ro ubuntu`",
			input: `docker run -v /proc:/host/proc:ro ubuntu`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1896",
					Message: "`docker ... -v /proc:/host/proc:ro` bind-mounts host /proc into the container — every process's `environ`/`cmdline` and `/proc/1/ns/` breakout handles become readable. Use `--cap-add=SYS_PTRACE` or host-side monitoring instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `podman run --volume=/sys:/host/sys alpine`",
			input: `podman run --volume=/sys:/host/sys alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1896",
					Message: "`podman ... -v /sys:/host/sys` bind-mounts host /sys into the container — every process's `environ`/`cmdline` and `/proc/1/ns/` breakout handles become readable. Use `--cap-add=SYS_PTRACE` or host-side monitoring instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1896")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1897(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt SH_GLOB` (explicit default)",
			input:    `unsetopt SH_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt SH_GLOB`",
			input: `setopt SH_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1897",
					Message: "`setopt SH_GLOB` disables Zsh-extended glob patterns — `*(N)` qualifiers, `<1-10>` ranges, and `(alt1|alt2)` alternation raise parse errors. Keep the option off; scope with `LOCAL_OPTIONS` if strict POSIX is needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_SH_GLOB`",
			input: `unsetopt NO_SH_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1897",
					Message: "`unsetopt NO_SH_GLOB` disables Zsh-extended glob patterns — `*(N)` qualifiers, `<1-10>` ranges, and `(alt1|alt2)` alternation raise parse errors. Keep the option off; scope with `LOCAL_OPTIONS` if strict POSIX is needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1897")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1898(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gpg --export KEYID`",
			input:    `gpg --export KEYID`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gpg --list-keys`",
			input:    `gpg --list-keys`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gpg --export-secret-keys KEYID` (leading)",
			input: `gpg --export-secret-keys KEYID`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1898",
					Message: "`gpg --export-secret-keys` writes the private key to stdout — one CI-log or wrong-tty redirect leaks it. Back up interactively on an air-gapped host, or write to a `umask 077` path and re-encrypt.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:  "invalid — `gpg KEYID --export-secret-subkeys` (trailing)",
			input: `gpg KEYID --export-secret-subkeys`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1898",
					Message: "`gpg --export-secret-subkeys` writes the private key to stdout — one CI-log or wrong-tty redirect leaks it. Back up interactively on an air-gapped host, or write to a `umask 077` path and re-encrypt.",
					Line:    1,
					Column:  12,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1898")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1899(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mokutil --list-enrolled`",
			input:    `mokutil --list-enrolled`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mokutil --import /root/MOK.der`",
			input:    `mokutil --import /root/MOK.der`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mokutil --disable-validation now` (leading)",
			input: `mokutil --disable-validation now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1899",
					Message: "`mokutil --disable-validation` stops the shim from validating kernel/modules against enrolled keys — Secure Boot becomes advisory. Leave validation on; enrol specific keys with `mokutil --import`.",
					Line:    1,
					Column:  10,
				},
			},
		},
		{
			name:  "invalid — `mokutil -l --disable-validation` (trailing)",
			input: `mokutil -l --disable-validation`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1899",
					Message: "`mokutil --disable-validation` stops the shim from validating kernel/modules against enrolled keys — Secure Boot becomes advisory. Leave validation on; enrol specific keys with `mokutil --import`.",
					Line:    1,
					Column:  13,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1899")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
