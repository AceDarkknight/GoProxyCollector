package util

import (
	"testing"
)

func TestRandomUA(t *testing.T) {
	t.Run("RandomUA Test", func(t *testing.T) {
		if got := RandomUA(); got == "" {
			t.Errorf("RandomUA() = %v, but expected string", got)
		}
	})
}

func TestVerifyHTTP(t *testing.T) {
	type args struct {
		ip   string
		port int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"unusable ip", args{"0.0.0.0", 80}, false},
		{"length of ip = 0", args{"", 80}, false},
		{"port is less than 0", args{"0.0.0.0", -20}, false},
		// Just for temporary test.
		// {"usable ip", args{"183.136.218.253", 80}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyHTTP(tt.args.ip, tt.args.port); got != tt.want {
				t.Errorf("VerifierHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}
