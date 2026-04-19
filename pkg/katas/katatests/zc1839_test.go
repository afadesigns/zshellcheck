package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1839(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `timedatectl set-ntp true`",
			input:    `timedatectl set-ntp true`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `systemctl enable chronyd`",
			input:    `systemctl enable chronyd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `timedatectl set-ntp false`",
			input: `timedatectl set-ntp false`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1839",
					Message: "`timedatectl set-ntp false` turns off network time sync — clock drift breaks TLS `notBefore`/`notAfter`, Kerberos, and TOTP. Leave NTP enabled; isolate frozen clocks to namespaces/CI.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `systemctl disable systemd-timesyncd`",
			input: `systemctl disable systemd-timesyncd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1839",
					Message: "`systemctl disable systemd-timesyncd` turns off network time sync — clock drift breaks TLS `notBefore`/`notAfter`, Kerberos, and TOTP. Leave NTP enabled; isolate frozen clocks to namespaces/CI.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1839")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
