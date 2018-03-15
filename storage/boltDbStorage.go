package storage

import (
	"encoding/json"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cihub/seelog"
)

type BoltDbStorage struct {
	Db         *bolt.DB
	Keys       []string
	bucketName string
	mutex      sync.Mutex
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
		Db:         db,
		bucketName: bucketName,
		Keys:       keys,
	}

	return &storage, nil
}

// Exist will check the given key is existed in DB or not.
func (s *BoltDbStorage) Exist(key string) bool {
	exist := false
	if s.Get(key) != nil {
		exist = true
	}

	return exist
}

// Get will get the json byte value of key.
func (s *BoltDbStorage) Get(key string) []byte {
	value := make([]byte, 0)
	s.Db.View(func(tx *bolt.Tx) error {
		value = append(value, tx.Bucket([]byte(s.bucketName)).Get([]byte(key))...)
		return nil
	})

	return value
}

// Delete the value by the given key.
func (s *BoltDbStorage) Delete(key string) bool {
	isSucceed := false
	err := s.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(s.bucketName)).Delete([]byte(key))
	})

	if err == nil {
		isSucceed = true
	}

	return isSucceed
}

// AddOrUpdate will add the value into DB if key is not existed, otherwise update the existing value.
// Value will be marshal as json format.
func (s *BoltDbStorage) AddOrUpdate(key string, value interface{}) error {
	content, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = s.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(s.bucketName)).Put([]byte(key), content)
	})

	return err
}

// GetAll will return all key-value in DB.
func (s *BoltDbStorage) GetAll() map[string][]byte {
	result := make(map[string][]byte)

	s.Db.View(func(tx *bolt.Tx) error {
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
	s.Db.Close()
}

// SyncKeys will sync the DB's key to memory.
func (s *BoltDbStorage) SyncKeys() {
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
func (s *BoltDbStorage) GetRandomOne() (string, []byte) {
	s.mutex.Lock()
	if len(s.Keys) == 0 {
		s.SyncKeys()
	}

	key := s.Keys[rand.New(rand.NewSource(time.Now().Unix())).Intn(len(s.Keys))]
	s.mutex.Unlock()

	return key, s.Get(key)
}
