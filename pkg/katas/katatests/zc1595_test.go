package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

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
