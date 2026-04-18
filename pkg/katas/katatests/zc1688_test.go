package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1688(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — aws s3 sync without --delete",
			input:    `aws s3 sync ./build s3://my-bucket/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — aws s3 cp",
			input:    `aws s3 cp file s3://bucket/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — aws s3 sync --delete",
			input: `aws s3 sync ./build s3://my-bucket/ --delete`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1688",
					Message: "`aws s3 sync --delete` wipes DST objects that are missing from SRC — a mistyped SRC bulk-deletes the bucket. Scope to the prefix, dry-run first, or enable versioning + MFA-delete.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — aws s3 sync between buckets with --delete",
			input: `aws s3 sync s3://src/ s3://dst/ --delete`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1688",
					Message: "`aws s3 sync --delete` wipes DST objects that are missing from SRC — a mistyped SRC bulk-deletes the bucket. Scope to the prefix, dry-run first, or enable versioning + MFA-delete.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1688")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
