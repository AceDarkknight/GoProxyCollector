package storage

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"strconv"
	"testing"

	"github.com/boltdb/bolt"
)

var testDb *BoltDbStorage

type testResult struct {
	t int
}

func init() {
	testDb, _ = NewBoltDbStorage("test.db", "testBucket")
}

func createBucket() {
	testDb.Db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("testBucket"))
		return err
	})
}

func deleteBucket() {
	testDb.Db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte("testBucket"))
	})
}

func TestBoltDbStorage_Get(t *testing.T) {
	createBucket()
	testDb.AddOrUpdate("2", &testResult{1})
	testDb.AddOrUpdate("0", nil)
	testDb.AddOrUpdate("3", &testResult{})

	expect1, _ := json.Marshal(&testResult{1})
	expect2, _ := json.Marshal(&testResult{})

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

	deleteBucket()
}

func TestBoltDbStorage_Exist(t *testing.T) {
	createBucket()
	testDb.AddOrUpdate("2", &testResult{1})
	testDb.AddOrUpdate("0", nil)
	testDb.AddOrUpdate("3", &testResult{})

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

	deleteBucket()
}

func TestBoltDbStorage_GetAll(t *testing.T) {
	createBucket()
	testDb.AddOrUpdate("2", &testResult{1})
	testDb.AddOrUpdate("0", nil)
	testDb.AddOrUpdate("3", &testResult{})

	expect1, _ := json.Marshal(&testResult{1})
	expect2, _ := json.Marshal(&testResult{})
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

	deleteBucket()
}

func BenchmarkBoltDbStorage_AddOrUpdate(b *testing.B) {
	createBucket()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		testDb.AddOrUpdate(strconv.Itoa(i), &testResult{})
	}

	deleteBucket()
}

func BenchmarkBoltDbStorage_AddOrUpdateParallel(b *testing.B) {
	createBucket()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := rand.Intn(100000)
			testDb.AddOrUpdate(strconv.Itoa(key), &testResult{})
		}
	})

	deleteBucket()
}

func BenchmarkBoltDbStorage_Get(b *testing.B) {
	createBucket()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		testDb.Get(strconv.Itoa(i))
	}

	deleteBucket()
}

func BenchmarkBoltDbStorage_GetParallel(b *testing.B) {
	createBucket()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testDb.Get(strconv.Itoa(rand.Intn(100)))
		}
	})

	deleteBucket()
}
