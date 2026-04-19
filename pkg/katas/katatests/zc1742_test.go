package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1742(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mc alias set NAME URL` (interactive prompt)",
			input:    `mc alias set myminio https://play.min.io`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mc alias list`",
			input:    `mc alias list`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mc alias set NAME URL ACCESS SECRET`",
			input: `mc alias set myminio https://play.min.io ACCESSKEY SECRETKEY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1742",
					Message: "`mc alias set ... ACCESS_KEY SECRET_KEY` puts S3 access and secret keys in argv — visible in `ps`, `/proc`, history. Drop the trailing keys (mc prompts) or use `MC_HOST_<alias>=URL` env-var form scoped to one command.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `mc config host add NAME URL ACCESS SECRET` (legacy)",
			input: `mc config host add myminio https://play.min.io ACCESSKEY SECRETKEY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1742",
					Message: "`mc config host add ... ACCESS SECRET` puts S3 access and secret keys in argv — visible in `ps`, `/proc`, history. Drop the trailing keys (mc prompts) or use `MC_HOST_<alias>=URL` env-var form scoped to one command.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1742")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
