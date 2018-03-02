package storage

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

type BoltDbStorage struct {
	Db     *bolt.DB
	bucket *bolt.Bucket
}

func NewBoltDbStorage(fileName string) (*BoltDbStorage, error) {
	if fileName == "" {
		return nil, errors.New("open boltdb fileName is empty")
	}

	db, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		return nil, err
	}

	var bucket *bolt.Bucket
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err = tx.CreateBucketIfNotExists([]byte("ip"))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil || bucket == nil {
		return nil, err
	}

	storage := BoltDbStorage{
		Db:     db,
		bucket: bucket,
	}

	return &storage, nil
}

func (s *BoltDbStorage) Exist(ip string) bool {
	exist := false
	if s.Get(ip) != nil {
		exist = true
	}

	return exist
}

func (s *BoltDbStorage) Get(ip string) []byte {
	var value []byte
	s.Db.View(func(tx *bolt.Tx) error {
		value = s.bucket.Get([]byte(ip))
		return nil
	})

	return value
}

func (s *BoltDbStorage) Delete(ip string) bool {
	isSucceed := false
	err := s.Db.Update(func(tx *bolt.Tx) error {
		return s.bucket.Delete([]byte(ip))
	})

	if err == nil {
		isSucceed = true
	}

	return isSucceed
}

func (s *BoltDbStorage) Update(ip string, info interface{}) error {
	content, err := json.Marshal(info)
	if err != nil {
		return err
	}

	err = s.Db.Update(func(tx *bolt.Tx) error {
		return s.bucket.Put([]byte(ip), content)
	})

	return err
}
