package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1513(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — make",
			input:    `make`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — make DESTDIR=/tmp/pkg install",
			input:    `make DESTDIR=/tmp/pkg install`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — make install",
			input: `make install`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1513",
					Message: "`make install` without `DESTDIR=` leaves no package-manager record. Set `DESTDIR=/tmp/pkgroot` and wrap in checkinstall / fpm, or use stow.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — gmake install",
			input: `gmake install`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1513",
					Message: "`make install` without `DESTDIR=` leaves no package-manager record. Set `DESTDIR=/tmp/pkgroot` and wrap in checkinstall / fpm, or use stow.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1513")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
