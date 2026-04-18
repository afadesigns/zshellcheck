package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1672(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — different command",
			input:    `semanage fcontext -a -t httpd_sys_content_t /var/www/app`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — chcon with no args",
			input:    `chcon`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chcon -t unconfined_t path",
			input: `chcon -t unconfined_t /usr/local/bin/script`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1672",
					Message: "`chcon` writes an ephemeral SELinux label — `restorecon` / policy rebuild reverts it. Persist via `semanage fcontext -a -t TYPE 'REGEX'` + `restorecon`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chcon -R -t bin_t dir",
			input: `chcon -R -t bin_t /srv/app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1672",
					Message: "`chcon` writes an ephemeral SELinux label — `restorecon` / policy rebuild reverts it. Persist via `semanage fcontext -a -t TYPE 'REGEX'` + `restorecon`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1672")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
