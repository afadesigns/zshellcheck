package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1609(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — aa-enforce reapplies enforcement",
			input:    `aa-enforce /etc/apparmor.d/usr.bin.firefox`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — apparmor_parser -r (reload, not remove)",
			input:    `apparmor_parser -r /etc/apparmor.d/profile`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — aa-disable",
			input: `aa-disable /etc/apparmor.d/usr.bin.firefox`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1609",
					Message: "`aa-disable` disables or softens the AppArmor profile — the confined process loses MAC restrictions. Review the profile's intent before disabling in automation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — aa-complain",
			input: `aa-complain /etc/apparmor.d/usr.bin.firefox`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1609",
					Message: "`aa-complain` disables or softens the AppArmor profile — the confined process loses MAC restrictions. Review the profile's intent before disabling in automation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — apparmor_parser -R",
			input: `apparmor_parser -R /etc/apparmor.d/profile`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1609",
					Message: "`apparmor_parser -R` removes the AppArmor profile from the kernel — the confined process loses MAC restrictions. Review the profile's intent before removing in automation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1609")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
