package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1702(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dpkg-reconfigure -f noninteractive",
			input:    `dpkg-reconfigure -f noninteractive tzdata`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — different command",
			input:    `dpkg -l`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dpkg-reconfigure tzdata (no frontend)",
			input: `dpkg-reconfigure tzdata`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1702",
					Message: "`dpkg-reconfigure` without `-f noninteractive` opens debconf dialogs — non-interactive pipelines hang. Pass `-f noninteractive` or export `DEBIAN_FRONTEND=noninteractive`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — dpkg-reconfigure --priority high (still interactive)",
			input: `dpkg-reconfigure -p high tzdata`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1702",
					Message: "`dpkg-reconfigure` without `-f noninteractive` opens debconf dialogs — non-interactive pipelines hang. Pass `-f noninteractive` or export `DEBIAN_FRONTEND=noninteractive`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1702")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
