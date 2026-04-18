package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1657(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — semanage permissive -d removes domain",
			input:    `semanage permissive -d httpd_t`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — semanage boolean -l listing",
			input:    `semanage boolean -l`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — semanage permissive -a httpd_t",
			input: `semanage permissive -a httpd_t`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1657",
					Message: "`semanage permissive -a` puts an SELinux domain in permissive mode — policy violations log but no longer block. Write a scoped allow rule with `audit2allow` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — semanage permissive --add sshd_t",
			input: `semanage permissive --add sshd_t`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1657",
					Message: "`semanage permissive -a` puts an SELinux domain in permissive mode — policy violations log but no longer block. Write a scoped allow rule with `audit2allow` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1657")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
