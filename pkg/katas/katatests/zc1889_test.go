package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1889(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `skopeo copy docker://a docker://b`",
			input:    `skopeo copy docker://a docker://b`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `skopeo copy --src-tls-verify=true docker://a docker://b`",
			input:    `skopeo copy --src-tls-verify=true docker://a docker://b`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `skopeo copy --src-tls-verify=false docker://a docker://b`",
			input: `skopeo copy --src-tls-verify=false docker://a docker://b`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1889",
					Message: "`skopeo --src-tls-verify=false` disables TLS verification on image copy — on-path attacker can substitute a malicious manifest. Pin a private CA with `--src-cert-dir`/`--dest-cert-dir` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `skopeo copy --dest-tls-verify=false docker://a docker://b`",
			input: `skopeo copy --dest-tls-verify=false docker://a docker://b`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1889",
					Message: "`skopeo --dest-tls-verify=false` disables TLS verification on image copy — on-path attacker can substitute a malicious manifest. Pin a private CA with `--src-cert-dir`/`--dest-cert-dir` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1889")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
