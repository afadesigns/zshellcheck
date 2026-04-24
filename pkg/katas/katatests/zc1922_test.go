package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1922(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `rpm --import /tmp/key.asc` (local file)",
			input:    `rpm --import /tmp/key.asc`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rpm --import https://pinned.example/key.asc`",
			input:    `rpm --import https://pinned.example/key.asc`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `rpm --import http://repo.example/key.asc`",
			input: `rpm --import http://repo.example/key.asc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1922",
					Message: "`rpm --import http://repo.example/key.asc` fetches a GPG key over plaintext — on-path attackers swap it, every future signed package installs. Use `https://` from a pinned origin, or `gpg --verify` against a known fingerprint.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:  "invalid — `rpmkeys --import ftp://repo.example/key.asc`",
			input: `rpmkeys --import ftp://repo.example/key.asc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1922",
					Message: "`rpm --import ftp://repo.example/key.asc` fetches a GPG key over plaintext — on-path attackers swap it, every future signed package installs. Use `https://` from a pinned origin, or `gpg --verify` against a known fingerprint.",
					Line:    1,
					Column:  10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1922")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
