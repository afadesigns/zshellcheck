package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1809(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gsutil ls gs://bucket`",
			input:    `gsutil ls gs://bucket`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gsutil rm gs://bucket/specific-object` (single object)",
			input:    `gsutil rm gs://bucket/specific-object`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gsutil -m rm -r gs://bucket/prefix`",
			input: `gsutil -m rm -r gs://bucket/prefix`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1809",
					Message: "`gsutil rm` with recursive/force deletes every matching GCS object (or the bucket itself). Preview with `gsutil ls`, enable Object Versioning / retention locks ahead of time, and prefer narrower object-level `gsutil rm` calls.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gsutil rb -f gs://bucket`",
			input: `gsutil rb -f gs://bucket`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1809",
					Message: "`gsutil rb` with recursive/force deletes every matching GCS object (or the bucket itself). Preview with `gsutil ls`, enable Object Versioning / retention locks ahead of time, and prefer narrower object-level `gsutil rm` calls.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1809")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
