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

// Exist will
func (s *BoltDbStorage) Exist(key string) bool {
	exist := false
	if s.Get(key) != nil {
		exist = true
	}

	return exist
}

func (s *BoltDbStorage) Get(key string) []byte {
	var value []byte
	s.Db.View(func(tx *bolt.Tx) error {
		copy(value, tx.Bucket([]byte(s.bucketName)).Get([]byte(key)))
		return nil
	})

	return value
}

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

func (s *BoltDbStorage) AddOrUpdate(key string, info interface{}) error {
	content, err := json.Marshal(info)
	if err != nil {
		return err
	}

	err = s.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(s.bucketName)).Put([]byte(key), content)
	})

	return err
}

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

func (s *BoltDbStorage) Close() {
	s.Db.Close()
}
