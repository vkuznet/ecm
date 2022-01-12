package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"log"
	"path"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

// KVRecord represents key-value record in our Store
type KVRecord struct {
	Key   string
	Value []byte
	Sha   string
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

func (s *Store) Add(rec KVRecord, sha string) error {
	// create hash value for given key
	var h hash.Hash
	if sha == "sha256" || rec.Sha == "sha256" {
		h = sha256.New()
		rec.Sha = "sha256"
	} else if sha == "sha512" || rec.Sha == "sha512" {
		h = sha512.New()
		rec.Sha = "sha512"
	} else {
		h = sha1.New()
		rec.Sha = "sha1"
	}
	h.Write([]byte(rec.Key))
	// if record value is not provided we'll create a hash for it
	// this will allow to anonimise the data
	if len(rec.Value) == 0 {
		rec.Value = []byte(hex.EncodeToString(h.Sum(nil)))
	}

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
