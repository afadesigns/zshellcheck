package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1107(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		want     int
	}{
		{
			name: "Valid arithmetic",
			src:  "if (( a > b )); then echo yes; fi",
			want: 0,
		},
		{
			name: "String comparison",
			src:  "if [[ $a == $b ]]; then echo yes; fi",
			want: 0,
		},
		{
			name: "File check",
			src:  "if [[ -f file ]]; then echo yes; fi",
			want: 0,
		},
		{
			name: "Invalid -eq",
			src:  "if [ $a -eq $b ]; then echo yes; fi",
			want: 1,
		},
		{
			name: "Invalid -gt in [[ ]]",
			src:  "if [[ $a -gt 5 ]]; then echo yes; fi",
			want: 1,
		},
		{
			name: "Invalid -le",
			src:  "while [ $count -le 10 ]; do ((count++)); done",
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.src, "ZC1107")
			if len(violations) != tt.want {
				t.Errorf("Test %q failed: want %d violations, got %d", tt.name, tt.want, len(violations))
			}
		})
	}
}
