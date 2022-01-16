package main

// KV store based on badger DB
// https://dgraph.io/docs/badger/get-started/

import (
	"log"
	"path"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

// KVRecord represents key-value record in our Store
type KVRecord struct {
	Key   string
	Value []byte
}

// Store represents Badger DB
type Store struct {
	DB *badger.DB
}

// NewStore create new Store object
func NewStore(dir string) (*Store, error) {
	opt := badger.DefaultOptions(dir)
	opt.ValueDir = path.Join(dir, "data")
	opt.Logger = nil

	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	return &Store{DB: db}, nil
}

// StoreCleaner may be run as goroutine to perform Badger GC
func (s *Store) Cleaner() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
	again:
		err := s.DB.RunValueLogGC(0.7)
		if err == nil {
			goto again
		}
	}
}

// Delete deletes key entry in our store
func (s *Store) Delete(key string) error {
	return s.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// Add adds give record to the store
func (s *Store) Add(rec KVRecord) error {
	// commit new key-value records into our store
	txn := s.DB.NewTransaction(true)
	defer txn.Discard()
	err := txn.Set([]byte(rec.Key), rec.Value)
	if err != nil {
		msg := "unable to set new key-value pair"
		log.Println(msg)
		return err
	}
	err = txn.Commit()
	if err != nil {
		msg := "unable to commit new key-value pair"
		log.Println(msg)
		return err
	}
	return nil
}

// Get finds record in our store for given key
func (s *Store) Get(key string) ([]byte, error) {
	var val []byte
	err := s.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		val, err = item.ValueCopy(nil)
		return nil
	})
	return val, err
}

// Close closes KV Store
func (s *Store) Close() {
	s.DB.Close()
}

// Records returns full list of records in our store
func (s *Store) Records() []KVRecord {
	var records []KVRecord
	err := s.DB.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			rec := KVRecord{Key: string(key), Value: val}
			records = append(records, rec)
			return nil
		}
		return nil
	})
	if err != nil {
		log.Println("fail during store iteration", err)
	}
	return records
}
