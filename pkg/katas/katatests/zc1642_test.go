package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1642(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tshark -w with -Z user",
			input:    `tshark -i eth0 -w /tmp/cap.pcap -Z analyst`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tshark without -w (display only)",
			input:    `tshark -i any`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tshark -w without -Z",
			input: `tshark -i any -w /tmp/cap.pcap`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1642",
					Message: "`tshark -w FILE` without `-Z USER` leaves the pcap root-owned. Add `-Z USER` to drop privileges for the capture.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — dumpcap -w without -Z",
			input: `dumpcap -i eth0 -w /tmp/cap.pcap`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1642",
					Message: "`dumpcap -w FILE` without `-Z USER` leaves the pcap root-owned. Add `-Z USER` to drop privileges for the capture.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1642")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
