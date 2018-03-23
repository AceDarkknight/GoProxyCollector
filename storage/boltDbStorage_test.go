package storage

import (
	"encoding/json"
	"reflect"
	"testing"
)

var testDb *BoltDbStorage

type testResult struct {
	t int
}

func init() {
	testDb, _ = NewBoltDbStorage("test.db", "testBucket")
	testDb.AddOrUpdate("2", testResult{1})
	testDb.AddOrUpdate("0", nil)
	testDb.AddOrUpdate("3", testResult{})
}

func TestBoltDbStorage_Get(t *testing.T) {
	expect1, _ := json.Marshal(testResult{1})
	expect2, _ := json.Marshal(testResult{})

	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"test1", args{"0"}, nil},
		{"test2", args{"1"}, nil},
		{"test3", args{"2"}, expect1},
		{"test4", args{"3"}, expect2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testDb.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoltDbStorage.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoltDbStorage_Exist(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{"2"}, true},
		{"test2", args{"1"}, false},
		{"test3", args{"0"}, false},
		{"test4", args{"3"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testDb.Exist(tt.args.key); got != tt.want {
				t.Errorf("BoltDbStorage.Exist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoltDbStorage_GetAll(t *testing.T) {
	expect1, _ := json.Marshal(testResult{1})
	expect2, _ := json.Marshal(testResult{})
	want := make(map[string][]byte)
	want["2"] = expect1
	want["3"] = expect2

	tests := []struct {
		name string
		want map[string][]byte
	}{
		{"test", want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testDb.GetAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoltDbStorage.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
