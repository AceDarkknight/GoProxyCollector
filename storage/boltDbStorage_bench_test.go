package storage

import (
	"math/rand"
	"strconv"
	"testing"
)

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

func BenchmarkBoltDbStorage_GetRandomOneParallel(b *testing.B) {
	createBucket()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testDb.GetRandomOne()
		}
	})

	deleteBucket()
}
