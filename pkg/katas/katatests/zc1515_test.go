package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1515(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sha256sum file",
			input:    `sha256sum file.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — b2sum file",
			input:    `b2sum file.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — md5sum file",
			input: `md5sum file.tar.gz`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1515",
					Message: "`md5sum` is collision-vulnerable — don't use for integrity checks. Use `sha256sum` / `sha512sum` / `b2sum` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sha1sum file",
			input: `sha1sum file.tar.gz`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1515",
					Message: "`sha1sum` is collision-vulnerable — don't use for integrity checks. Use `sha256sum` / `sha512sum` / `b2sum` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1515")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
