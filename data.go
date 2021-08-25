package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

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

// String provides string representation of vault record
func (r *VaultRecord) String() string {
	data, err := json.Marshal(r)
	if err == nil {
		return string(data)
	}
	return fmt.Sprintf("%+v", r)
}

// Details provides record details
func (r *VaultRecord) Details() (string, string, string, string, string) {
	name := r.Name
	note := r.Note
	rurl := r.URL
	var login, password string
	for _, item := range r.Items {
		if strings.ToLower(item.Name) == "login" {
			login = item.Value
		}
		if strings.ToLower(item.Name) == "password" {
			password = item.Value
		}
	}
	return name, rurl, login, password, note
}
