package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1469(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dnf install curl",
			input:    `dnf install curl`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — rpm -i package.rpm",
			input:    `rpm -i package.rpm`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dnf install --nogpgcheck foo",
			input: `dnf install --nogpgcheck foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1469",
					Message: "Package signature verification disabled (dnf --nogpgcheck) — any mirror / MITM becomes immediate root.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — yum install --nogpgcheck foo",
			input: `yum install --nogpgcheck foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1469",
					Message: "Package signature verification disabled (yum --nogpgcheck) — any mirror / MITM becomes immediate root.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — rpm -i --nosignature foo.rpm",
			input: `rpm -i --nosignature foo.rpm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1469",
					Message: "Package signature verification disabled (rpm --nosignature) — any mirror / MITM becomes immediate root.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — zypper install --no-gpg-checks foo",
			input: `zypper install --no-gpg-checks foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1469",
					Message: "Package signature verification disabled (zypper --no-gpg-checks) — any mirror / MITM becomes immediate root.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1469")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
