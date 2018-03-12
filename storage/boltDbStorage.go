package storage

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

type BoltDbStorage struct {
	Db         *bolt.DB
	bucketName string
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
		_, err := tx.CreateBucketIfNotExists([]byte("ip"))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	storage := BoltDbStorage{
		Db:         db,
		bucketName: bucketName,
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
	var value []byte
	s.Db.View(func(tx *bolt.Tx) error {
		copy(value, tx.Bucket([]byte(s.bucketName)).Get([]byte(key)))
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
			var key, value []byte
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
