package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1599(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — plain ldconfig",
			input:    `ldconfig`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ldconfig -f /etc/ld.so.conf.d/custom.conf",
			input:    `ldconfig -f /etc/ld.so.conf.d/custom.conf`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ldconfig -f $LDCONF (variable)",
			input:    `ldconfig -f $LDCONF`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ldconfig -f /tmp/fake.conf",
			input: `ldconfig -f /tmp/fake.conf`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1599",
					Message: "`ldconfig -f /tmp/fake.conf` uses a config outside `/etc/`. If the file is attacker-writable, every binary on the host loads the attacker's library. Keep config under `/etc/ld.so.conf.d/`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ldconfig -f /var/tmp/x.conf",
			input: `ldconfig -f /var/tmp/x.conf`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1599",
					Message: "`ldconfig -f /var/tmp/x.conf` uses a config outside `/etc/`. If the file is attacker-writable, every binary on the host loads the attacker's library. Keep config under `/etc/ld.so.conf.d/`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1599")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
