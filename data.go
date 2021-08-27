package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

// Vault represent our vault
type Vault struct {
	Filename string
	Cipher   string
	Secret   string
	Verbose  int
	Records  []VaultRecord
}

// Update vault records
func (v *Vault) Update(rec VaultRecord) {
	for i, r := range v.Records {
		if r.ID == rec.ID {
			if v.Verbose > 0 {
				log.Printf("update record %+v", rec)
			}
			v.Records[i] = rec
		}
	}
}

// Write vault records
func (v *Vault) Write() {
	// make backup vault first
	backup(v.Filename)

	file, err := os.Create(v.Filename)
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(file)
	for _, rec := range v.Records {
		// marshall single record
		data, err := json.Marshal(rec)
		if err != nil {
			log.Fatal(err)
		}

		// encrypt our record
		if v.Verbose > 0 {
			log.Printf("record '%s'\n", string(data))
		}
		edata := data
		if v.Cipher != "" {
			edata, err = encrypt(data, v.Secret, v.Cipher)
			if err != nil {
				log.Fatal(err)
			}
		}
		if v.Verbose > 1 {
			log.Printf("write data record\n%v\nsecret '%v'", edata, v.Secret)
		}
		w.Write(edata)
		w.Write([]byte("---\n"))
		w.Flush()
	}
}
