package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
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
	Filename             string        // vault filename
	Cipher               string        // vault cipher
	Secret               string        // vault secret
	Verbose              int           // verbose mode
	Records              []VaultRecord // vault records
	LastModificationTime time.Time     // vault last modification time
	LastBackup           string        // vault last backup
	Size                 int64         // vault size
	Mode                 string        // vault mode
}

// Update vault records
func (v *Vault) Update(rec VaultRecord) {
	for i, r := range v.Records {
		if r.ID == rec.ID {
			if v.Verbose > 0 {
				log.Printf("update record %+v", rec)
			}
			v.Records[i] = rec
			v.LastModificationTime = time.Now()
		}
	}
}

// Write provides write capability to vault records
func (v *Vault) Write() {
	var err error
	// make backup vault first
	v.LastBackup, _, err = backup(v.Filename)

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
	finfo, err := os.Stat(v.Filename)
	if err == nil {
		v.Size = finfo.Size()
		v.LastModificationTime = finfo.ModTime()
		v.Mode = finfo.Mode().String()
	} else {
		log.Printf("unable to get stat for %s, error", v.Filename, err)
	}
}

// helper function to read vault and return list of records
func (v *Vault) Read() error {

	// check first if file exsist
	if _, err := os.Stat(v.Filename); os.IsNotExist(err) {
		log.Println("vault %s does not exists", v.Filename)
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
		v.LastModificationTime = finfo.ModTime()
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
		if v.Verbose > 0 {
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
		if matched, err := regexp.MatchString(pat, rec.Name); err == nil && matched {
			if v.Verbose > 0 {
				log.Println("matched record name")
			}
			out = append(out, rec)
		} else if strings.Contains(rec.Name, pat) {
			out = append(out, rec)
			if v.Verbose > 0 {
				log.Println("substring of record name")
			}
		} else if matched, err := regexp.MatchString(pat, rec.Note); err == nil && matched {
			if v.Verbose > 0 {
				log.Println("match record note")
			}
			out = append(out, rec)
		} else if InList(pat, rec.Tags) {
			if v.Verbose > 0 {
				log.Println("found in record tags")
			}
			out = append(out, rec)
		} else if strings.Contains(rec.URL, pat) {
			out = append(out, rec)
			if v.Verbose > 0 {
				log.Println("substring of record URL")
			}
		}
		for _, item := range rec.Items {
			if matched, err := regexp.MatchString(pat, item.Name); err == nil && matched {
				if v.Verbose > 0 {
					log.Println("match record item namej")
				}
				out = append(out, rec)
			}
		}
	}
	return out
}

// Info provides information about the vault
func (v *Vault) Info() string {
	tstamp := v.LastModificationTime.String()
	size := SizeFormat(v.Size)
	mode := v.Mode
	cipher := v.Cipher
	nrec := len(v.Records)
	info := fmt.Sprintf("vault %s\nLast modified: %s\nSize %s, mode %s\n %d records, encrypted with %s cipher", v.Filename, tstamp, size, mode, nrec, cipher)
	if v.Verbose > 0 {
		log.Printf("vault %+v", v)
		log.Println(info)
	}
	return info
}
