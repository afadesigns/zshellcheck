package katas

import (
	"testing"
)

func TestKatas(t *testing.T) {
	if len(Registry.KatasByID) == 0 {
		t.Errorf("Registry is empty")
	}
}
