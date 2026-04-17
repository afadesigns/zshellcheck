package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1570(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — smbclient -U user //server/share",
			input:    `smbclient -U user //server/share`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — smbclient -N //server/share",
			input: `smbclient -N //server/share`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1570",
					Message: "`smbclient -N` is anonymous SMB access — any host on-net can read the share. Use credentials=<file> 0600 or -k.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — mount.cifs -N //server/share /mnt",
			input: `mount.cifs -N //server/share /mnt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1570",
					Message: "`mount.cifs -N` is anonymous SMB access — any host on-net can read the share. Use credentials=<file> 0600 or -k.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1570")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
