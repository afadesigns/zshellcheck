package katas

import (
	"testing"
)

func TestKatas(t *testing.T) {
	registry := GetRegistry()
	if len(registry) == 0 {
		t.Errorf("Registry is empty")
	}
}

