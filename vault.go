package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	uuid "github.com/google/uuid"
)

// Record represent map of key-valut pairs
type Record map[string]string

// ValutRecord represents full vault record
type VaultRecord struct {
	ID               string    // record ID
	Map              Record    // record map (key-valut pairs)
	Attachments      []string  // record attachment files
	ModificationTime time.Time // record modification time
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
func (r *VaultRecord) Keys() []string {
	// predefined keys order
	keys := []string{"Name", "Login", "Password"}
	// output keys
	var out []string
	// map keys
	var mapKeys []string
	for k, _ := range r.Map {
		mapKeys = append(mapKeys, k)
	}
	sort.Sort(StringList(mapKeys))
	for _, k := range mapKeys {
		if k == "Name" || k == "Login" || k == "Password" {
			continue
		}
		out = append(out, k)
	}
	keys = append(keys, out...)
	return keys
}

// NewVaultRecord creates new VaultRecord
func NewVaultRecord(kind string) *VaultRecord {
	uid := uuid.NewString()
	rmap := make(Record)
	var attributes []string
	switch kind {
	case "note":
		attributes = []string{"Name", "Tags", "Note"}
	case "file":
		attributes = []string{"Name", "File", "Tags", "Note"}
	default: // default login record
		attributes = []string{"Name", "Login", "Password", "URL", "Tags", "Note"}
	}
	for _, attr := range attributes {
		rmap[attr] = ""
	}
	return &VaultRecord{ID: uid, Map: rmap, ModificationTime: time.Now()}
}

// Vault represent our vault
type Vault struct {
	Filename         string        // vault filename
	Cipher           string        // vault cipher
	Secret           string        // vault secret
	Verbose          int           // verbose mode
	Records          []VaultRecord // vault records
	ModificationTime time.Time     // vault last modification time
	LastBackup       string        // vault last backup
	Size             int64         // vault size
	Mode             string        // vault mode
}

// AddRecord vault record and return its index
func (v *Vault) AddRecord(kind string) int {
	rec := NewVaultRecord(kind)
	v.Records = append(v.Records, *rec)
	return len(v.Records) - 1
}

// Update vault records
func (v *Vault) Update(rec VaultRecord) {
	updated := false
	for i, r := range v.Records {
		if r.ID == rec.ID {
			if v.Verbose > 0 {
				log.Printf("update record %+v", rec)
			}
			rec.ModificationTime = time.Now()
			v.Records[i] = rec
			v.ModificationTime = time.Now()
			updated = true
		}
	}
	if !updated {
		// insert new record
		v.Records = append(v.Records, rec)
	}
}

// Write provides write capability to vault records
func (v *Vault) Write() {
	var err error
	// TODO: fix backup once we move to directory based vault
	// make backup vault first
	//     v.LastBackup, _, err = backup(v.Filename)
	//     if err != nil {
	//         log.Fatal(err)
	//     }

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
		if v.Verbose > 1 {
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
		w.Write([]byte(separator))
		w.Flush()
	}
	finfo, err := os.Stat(v.Filename)
	if err == nil {
		v.Size = finfo.Size()
		v.ModificationTime = finfo.ModTime()
		v.Mode = finfo.Mode().String()
	} else {
		log.Printf("unable to get stat for %s, error", v.Filename, err)
	}
}

// helper function to read vault and return list of records
func (v *Vault) Read() error {

	// check first if file exsist
	if _, err := os.Stat(v.Filename); os.IsNotExist(err) {
		log.Printf("vault %s does not exists, will create one", v.Filename)
		_, err := os.Create(v.Filename)
		if err != nil {
			log.Fatal(err)
		}
		return err
	}

	// alwasy keep file safe
	err := os.Chmod(v.Filename, 0600)
	if err != nil {
		log.Println("unable to change file permission of", v.Filename)
	}

	// get vault file info
	finfo, err := os.Stat(v.Filename)
	if err == nil {
		v.Size = finfo.Size()
		v.ModificationTime = finfo.ModTime()
		v.Mode = finfo.Mode().String()
	} else {
		log.Printf("unable to get stat for %s, error", v.Filename, err)
	}

	// open file
	file, err := os.Open(v.Filename)
	if err != nil {
		log.Println("unable to open a vault", err)
		return err
	}
	// remember to close the file at the end of the program
	defer file.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(file)
	scanner.Split(pwmSplitFunc)
	for scanner.Scan() {
		text := scanner.Text()
		textData := []byte(text)
		if v.Verbose > 1 {
			log.Printf("read record\n%v\n", textData)
		}

		data := textData
		if v.Cipher != "" {
			data, err = decrypt(textData, v.Secret, v.Cipher)
			if err != nil {
				log.Printf("unable to decrypt data\n%v\nerror %v", textData, err)
				return err
			}
		}

		var rec VaultRecord
		err = json.Unmarshal(data, &rec)
		if err != nil {
			log.Println("ERROR: unable to unmarshal the data", err)
			return err
		}
		v.Records = append(v.Records, rec)
	}
	return nil
}

// Find method finds given pattern in our vault and return its index
func (v *Vault) Find(pat string) []VaultRecord {
	var out []VaultRecord
	for _, rec := range v.Records {
		for key, val := range rec.Map {
			if strings.Contains(key, pat) {
				out = append(out, rec)
				if v.Verbose > 0 {
					log.Println("match record key")
				}
			} else if matched, err := regexp.MatchString(pat, val); err == nil && matched {
				if v.Verbose > 0 {
					log.Println("matched record value")
				}
				out = append(out, rec)
			}
		}
	}
	return out
}

// Info provides information about the vault
func (v *Vault) Info() string {
	tstamp := v.ModificationTime.String()
	size := SizeFormat(v.Size)
	mode := v.Mode
	cipher := v.Cipher
	nrec := len(v.Records)
	info := fmt.Sprintf("vault %s\nLast modified: %s\nSize %s, mode %s\n%d records, encrypted with %s cipher", v.Filename, tstamp, size, mode, nrec, cipher)
	if v.Verbose > 0 {
		log.Println(info)
	}
	return info
}
