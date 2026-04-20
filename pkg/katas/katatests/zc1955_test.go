package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1955(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `rfkill list`",
			input:    `rfkill list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rfkill unblock all`",
			input:    `rfkill unblock all`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `rfkill block all`",
			input: `rfkill block all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1955",
					Message: "`rfkill block all` hard-downs the radio(s) — host drops off the network in one call. Scope to the radio type that really needs it and schedule an `at now + N minutes` unblock for self-recovery.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `rfkill block wifi`",
			input: `rfkill block wifi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1955",
					Message: "`rfkill block wifi` hard-downs the radio(s) — host drops off the network in one call. Scope to the radio type that really needs it and schedule an `at now + N minutes` unblock for self-recovery.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1955")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
