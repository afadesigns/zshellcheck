package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1535(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ip link set eth0 up",
			input:    `ip link set eth0 up`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ip link set eth0 promisc off",
			input:    `ip link set eth0 promisc off`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ip link set eth0 promisc on",
			input: `ip link set eth0 promisc on`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1535",
					Message: "Interface put into promiscuous mode — sniffer-in-place. Re-disable after capture, or grant tcpdump CAP_NET_RAW instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ifconfig eth0 promisc",
			input: `ifconfig eth0 promisc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1535",
					Message: "Interface put into promiscuous mode — sniffer-in-place. Re-disable after capture, or grant tcpdump CAP_NET_RAW instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1535")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
