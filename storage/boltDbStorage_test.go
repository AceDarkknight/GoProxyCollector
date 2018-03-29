package storage

import (
	"encoding/json"
	"reflect"
	"strconv"
	"sync"
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

func TestBoltDbStorage_AddOrUpdateParallel(t *testing.T) {
	createBucket()
	var wg sync.WaitGroup
	num := 200
	for i := 0; i < num; i++ {
		wg.Add(1)
		key := strconv.Itoa(i)
		value := &testResult{t: i}
		go func(k string, v *testResult) {
			defer wg.Done()
			testDb.AddOrUpdate(k, v)
		}(key, value)
	}

	wg.Wait()
	for i := 0; i < num; i++ {
		key := strconv.Itoa(i)
		value := &testResult{t: i}
		want, _ := json.Marshal(value)
		t.Run(key, func(t *testing.T) {
			if got := testDb.Get(key); !reflect.DeepEqual(got, want) {
				t.Errorf("parallel run BoltDbStorage.AddOrUpdate() = %v,want %v ", got, want)
			}
		})
	}

	t.Run("all", func(t *testing.T) {
		if got := testDb.GetAll(); len(got) != num {
			t.Errorf("parallel run BoltDbStorage.AddOrUpdate() = %v,want %v ", got, num)
		}
	})

	deleteBucket()
}

func TestBoltDbStorage_DeleteParallel(t *testing.T) {
	createBucket()
	var wg sync.WaitGroup
	num := 200
	for i := 0; i < num; i++ {
		key := strconv.Itoa(i)
		value := &testResult{t: i}
		testDb.AddOrUpdate(key, value)
	}

	for i := 0; i < num; i++ {
		wg.Add(1)
		key := strconv.Itoa(i)
		go func(k string) {
			defer wg.Done()
			testDb.Delete(k)
		}(key)
	}

	wg.Wait()
	for i := 0; i < num; i++ {
		key := strconv.Itoa(i)
		var want []byte
		t.Run(key, func(t *testing.T) {
			if got := testDb.Get(key); !reflect.DeepEqual(got, want) {
				t.Errorf("parallel run BoltDbStorage.Delet() = %v,want %v ", got, want)
			}
		})
	}

	t.Run("all", func(t *testing.T) {
		want := make(map[string][]byte)
		if got := testDb.GetAll(); !reflect.DeepEqual(got, want) {
			t.Errorf("parallel run BoltDbStorage.AddOrUpdate() = %v,want %v ", got, want)
		}
	})

	deleteBucket()
}

func TestNewBoltDbStorage(t *testing.T) {
	a := 0.2
	b := 80

	s := a * float64(b)
}
