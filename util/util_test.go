package util

import "testing"

func TestRandomUA(t *testing.T) {
	t.Run("RandomUA Test", func(t *testing.T) {
		if got := RandomUA(); got == "" {
			t.Errorf("RandomUA() = %v, but expected string", got)
		}
	})
}
