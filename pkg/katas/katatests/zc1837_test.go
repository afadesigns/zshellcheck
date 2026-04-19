package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

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
