package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1725(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cargo publish` (no inline token)",
			input:    `cargo publish`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cargo login --token -` (stdin sentinel)",
			input:    `cargo login --token -`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cargo build --token foo` (not a publish subcmd)",
			input:    `cargo build --token foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `cargo publish --token TOKEN`",
			input: `cargo publish --token cio_abc123`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1725",
					Message: "`cargo publish --token cio_abc123` puts the credential in argv — visible in `ps`, `/proc`, history. Pipe via stdin (`--token -`) or use env vars like `CARGO_REGISTRY_TOKEN` / `NPM_TOKEN`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `cargo login --token=TOKEN`",
			input: `cargo login --token=cio_abc123`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1725",
					Message: "`cargo login --token=cio_abc123` puts the credential in argv — visible in `ps`, `/proc`, history. Pipe via stdin (`--token -`) or use env vars like `CARGO_REGISTRY_TOKEN` / `NPM_TOKEN`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `npm publish --otp 123456`",
			input: `npm publish --otp 123456`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1725",
					Message: "`npm publish --otp 123456` puts the credential in argv — visible in `ps`, `/proc`, history. Pipe via stdin (`--token -`) or use env vars like `CARGO_REGISTRY_TOKEN` / `NPM_TOKEN`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1725")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
