package main

import (
	"fmt"
	"log"
	"testing"
)

// TestKVStore function
func TestKVStore(t *testing.T) {
	store, err := NewStore("/tmp/kvdb")
	if err != nil {
		t.Errorf("unable to create new KV store %v", err)
		return
	}
	defer store.Close()

	// store new record
	key := "test"
	val := "test"
	rec := KVRecord{Key: key, Value: []byte(val)}
	err = store.Add(rec)
	if err != nil {
		t.Errorf("fail to add new key-value pair to KV store %v", err)
		return
	}

	// fetch data from KV store
	data, err := store.Get(key)
	if err != nil {
		t.Errorf("fail to get value from KV store %v", err)
		return
	}

	if string(data) != val {
		t.Error("stored value does not match")
		return
	}

	// delete our data in KV store
	err = store.Delete(key)
	if err != nil {
		t.Errorf("fail to delete key in KV store %v", err)
		return
	}

	data, err = store.Get(key)
	if err == nil {
		t.Errorf("found key %s in KV store %v", key, err)
		return
	}
}

// BenchmarkAddRecords provides benchmark test for KV store Add API
func BenchmarkAddRecords(b *testing.B) {
	store, err := NewStore("/tmp/kvdb")
	if err != nil {
		b.Errorf("unable to create new KV store %v", err)
		return
	}
	defer store.Close()

	// perform benchmark test
	val := "test"
	for n := 0; n < b.N; n++ {
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("key-%d", i)
			rec := KVRecord{Key: key, Value: []byte(val)}
			err = store.Add(rec)
			if err != nil {
				b.Error(err.Error())
			}
		}
	}
}

// BenchmarkGetRecords provides benchmark test for KV store Get API
func BenchmarkGetRecords(b *testing.B) {
	store, err := NewStore("/tmp/kvdb")
	if err != nil {
		b.Errorf("unable to create new KV store %v", err)
		return
	}
	defer store.Close()

	// perform benchmark test
	val := "test"
	for n := 0; n < b.N; n++ {
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("key-%d", i)
			data, err := store.Get(key)
			if err != nil {
				b.Error(err.Error())
			}
			if string(data) != val {
				b.Error("returned value does not match")
			}
		}
	}
}

// TestKVStoreRecords function
func TestKVStoreRecords(t *testing.T) {
	store, err := NewStore("/tmp/kvdb")
	if err != nil {
		t.Errorf("unable to create new KV store %v", err)
		return
	}
	defer store.Close()

	// store new record
	key := "test"
	val := "test"
	rec := KVRecord{Key: key, Value: []byte(val)}
	err = store.Add(rec)
	if err != nil {
		t.Errorf("fail to add new key-value pair to KV store %v", err)
		return
	}

	// get records from our store
	records := store.Records()
	for _, r := range records {
		log.Printf("kvdb record %+v", r)
	}
}
