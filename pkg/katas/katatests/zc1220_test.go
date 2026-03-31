package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

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
