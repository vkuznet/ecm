package main

// VaultItem holds individual key-value pair for given item
type VaultItem struct {
	Name  string // name of the key (label)
	Value string // value of the data
}

// ValutRecord represents full vault record
type VaultRecord struct {
	ID    string      // record ID
	Name  string      // record Name
	URL   string      // url value record represent
	Tags  []string    // record tags
	Items []VaultItem // list of record items
	Note  string      // record note
}
