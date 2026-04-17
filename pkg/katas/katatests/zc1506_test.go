package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1506(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sg audio -c cmd",
			input:    `sg audio -c 'ls /var/log'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — newgrp audio",
			input: `newgrp audio`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1506",
					Message: "`newgrp` starts a new shell — script either hangs or exits. Use `sg <group> -c <cmd>` or systemd `SupplementaryGroups=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1506")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
