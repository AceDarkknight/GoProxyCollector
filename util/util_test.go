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
		//{"usable ip", args{"61.135.217.7", 80}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyProxyIp(tt.args.ip, tt.args.port); got != tt.want {
				t.Errorf("VerifierHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIp(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{"0.0.0.0"}, true},
		{"test2", args{"2555.255.255.255"}, false},
		{"test3", args{"-=0:125.125.125"}, false},
		{"test4", args{"-1.0.0.0"}, false},
		{"test5", args{"127.0.0.-2"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIp(tt.args.ip); got != tt.want {
				t.Errorf("IsIp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInputMatchRegex(t *testing.T) {
	type args struct {
		input string
		regex string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{input: "", regex: ""}, true},
		{"test2", args{input: "123456", regex: ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInputMatchRegex(tt.args.input, tt.args.regex); got != tt.want {
				t.Errorf("IsInputMatchRegex() = %v, want %v", got, tt.want)
			}
		})
	}
}
