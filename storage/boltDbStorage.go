package storage

import (
	"encoding/json"
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cihub/seelog"
)

type BoltDbStorage struct {
	db              *bolt.DB
	Keys            []string
	bucketName      string
	mutex           sync.Mutex
	indexContentMap sync.Map
	keyIndexMap     sync.Map
	length          int32
}

type content struct {
	key        string
	resultByte []byte
}

// NewBoltDbStorage will return a boltdb object and error.
func NewBoltDbStorage(fileName string, bucketName string) (*BoltDbStorage, error) {

	if fileName == "" {
		return nil, errors.New("open boltdb whose fileName is empty")
	}

	if bucketName == "" {
		return nil, errors.New("create a bucket whose name is empty")
	}

	db, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, 10000)

	storage := BoltDbStorage{
		db:              db,
		bucketName:      bucketName,
		Keys:            keys,
		indexContentMap: sync.Map{},
		keyIndexMap:     sync.Map{},
	}

	return &storage, nil
}

// Exist will check the given key is existed in DB or not.
func (s *BoltDbStorage) Exist(key string) bool {
	return s.Get(key) != nil
}

// Get will get the json byte value of key.
func (s *BoltDbStorage) Get(key string) []byte {
	var value []byte

	if index, ok := s.keyIndexMap.Load(key); ok {
		if c, ok := s.indexContentMap.Load(index); ok {
			value = append(value, c.(content).resultByte...)
		}
	}

	// s.db.View(func(tx *bolt.Tx) error {
	// 	value = append(value, tx.Bucket([]byte(s.bucketName)).Get([]byte(key))...)
	// 	return nil
	// })

	return value
}

// Delete the value by the given key.
func (s *BoltDbStorage) Delete(key string) bool {
	isSucceed := false

	if index, ok := s.keyIndexMap.Load(key); ok {
		if _, ok = s.indexContentMap.Load(index); ok {
			s.keyIndexMap.Delete(key)
			s.indexContentMap.Delete(index)
			atomic.AddInt32(&s.length, -1)
		}
	}

	err := s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(s.bucketName)).Delete([]byte(key))
	})

	if err == nil {
		isSucceed = true
	}

	return isSucceed
}

// AddOrUpdate will add the value into DB if key is not existed, otherwise update the existing value.
// Null value will be ignored and the value will be marshal as json format.
func (s *BoltDbStorage) AddOrUpdate(key string, value interface{}) error {
	if value == nil {
		return errors.New("value is null")
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if index, ok := s.keyIndexMap.Load(key); ok {
		if _, ok = s.indexContentMap.Load(index); ok {

		}
	}

	s.indexContentMap.Store(key, bytes)

	err = s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(s.bucketName)).Put([]byte(key), bytes)
	})

	return err
}

// GetAll will return all key-value in DB.
func (s *BoltDbStorage) GetAll() map[string][]byte {
	result := make(map[string][]byte)

	s.db.View(func(tx *bolt.Tx) error {
		tx.Bucket([]byte(s.bucketName)).ForEach(func(k, v []byte) error {
			key, value := make([]byte, len(k)), make([]byte, len(v))
			copy(key, k)
			copy(value, v)
			result[string(key)] = value
			return nil
		})

		return nil
	})

	return result
}

// Close will close the DB.
func (s *BoltDbStorage) Close() {
	s.db.Close()
}

// Sync will sync the DB's key to memory.
func (s *BoltDbStorage) Sync() {
	result := s.GetAll()
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Keys = s.Keys[0:0:cap(s.Keys)]
	for k := range result {
		s.Keys = append(s.Keys, k)
	}

	seelog.Debug(s.Keys)
}

// Get one random record.
func (s *BoltDbStorage) GetRandomOne() []byte {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.Keys) == 0 {
		return nil
	}

	key := s.Keys[rand.New(rand.NewSource(time.Now().Unix())).Intn(len(s.Keys))]

	return s.Get(key)
}
