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
}

func TestBoltDbStorage_Get(t *testing.T) {
	expect, _ := json.Marshal(testResult{1})
	testDb.AddOrUpdate("2", testResult{1})

	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"get inexistent record", args{"1"}, make([]byte, 0)},
		{"get existed record", args{"2"}, expect},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testDb.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoltDbStorage.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
