package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

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
