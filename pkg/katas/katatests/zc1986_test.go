package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1986(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `touch $FILE` (current clock)",
			input:    `touch $FILE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `touch -c $FILE` (no create, current clock)",
			input:    `touch -c $FILE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `touch -d now $FILE`",
			input: `touch -d now $FILE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1986",
					Message: "`touch -d` writes a specific atime/mtime — also the classic \"age the dropped file\" antiforensics pattern. Derive from `$SOURCE_DATE_EPOCH` where the intent is deterministic.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `touch -t 202401011200 $FILE`",
			input: `touch -t 202401011200 $FILE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1986",
					Message: "`touch -t` writes a specific atime/mtime — also the classic \"age the dropped file\" antiforensics pattern. Derive from `$SOURCE_DATE_EPOCH` where the intent is deterministic.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `touch -r $REF $FILE`",
			input: `touch -r $REF $FILE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1986",
					Message: "`touch -r` writes a specific atime/mtime — also the classic \"age the dropped file\" antiforensics pattern. Derive from `$SOURCE_DATE_EPOCH` where the intent is deterministic.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1986")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
