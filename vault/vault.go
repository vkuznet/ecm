package vault

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	uuid "github.com/google/uuid"
	"github.com/vkuznet/ecm/crypt"
)

// recordAttribute performs conversion from one record attribute
// name to another, e.g. when we import 1Password records to ECM format
func recordAttribute(key string) string {
	if key == "Username" { // 1Password convention
		key = "Login"
	} else if key == "Title" {
		key = "Name"
	}
	return key
}

// Record represent map of key-valut pairs
type Record map[string]string

// VaultRecord represents full vault record
type VaultRecord struct {
	ID               string    // record ID
	Map              Record    // record map (key-vault pairs)
	Attachments      []string  // record attachment files
	ModificationTime time.Time // record modification time
}

// String provides string representation of vault record
func (r *VaultRecord) String() string {
	data, err := json.MarshalIndent(r, "", "   ")
	if err == nil {
		return string(data)
	}
	return ""
}

// Keys provides vault record keys
func (r *VaultRecord) Keys() []string {
	// predefined keys order
	keys := []string{"Name", "Login", "Password"}
	// output keys
	var out []string
	// map keys
	var mapKeys []string
	for k := range r.Map {
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

// WriteRecord writes single record to the vault area
func (r *VaultRecord) WriteRecord(vdir, secret, cipher string, verbose int) error {
	var err error
	tstamp := time.Now().Format(time.RFC3339)
	if r.ID == "" {
		log.Fatalf("unable to write record without ID, record %v", r)
	}
	fname := fmt.Sprintf("%s.%s", filepath.Join(vdir, r.ID), cipher)
	bdir := filepath.Join(vdir, "backups")
	err = os.MkdirAll(bdir, 0755)
	if err != nil {
		log.Fatalf("unable to create %s, error %v", bdir, err)
	}
	bname := filepath.Join(bdir, fmt.Sprintf("%s.%s-%s", r.ID, cipher, tstamp))
	// make backup of our record
	_, err = backup(fname, bname)
	if err != nil {
		if verbose > 0 {
			log.Println("unable to make backup for record", r.ID, " error ", err)
		}
	}

	file, err := os.Create(fname)
	if err != nil {
		log.Println("unable to create file name", fname, " error ", err)
		return err
	}
	w := bufio.NewWriter(file)
	// marshall single record
	data, err := json.Marshal(r)
	if err != nil {
		log.Println("unable to Marshal record, error ", err)
		return err
	}

	// encrypt our record
	if verbose > 1 {
		log.Printf("record '%s' using cipher %s\n", string(data), cipher)
	} else if verbose > 0 {
		log.Printf("record '%s' using cipher %s\n", r.ID, cipher)
	}
	edata := data
	if cipher != "" {
		edata, err = crypt.Encrypt(data, secret, cipher)
		if err != nil {
			log.Println("unable to encrypt record, error ", err)
			return err
		}
	}
	if verbose > 1 {
		log.Printf("write data record\n%v\nsecret '%v'", edata, secret)
	}
	w.Write(edata)
	w.Flush()
	return nil

}

// NewVaultRecord creates new VaultRecord
func NewVaultRecord(kind string) *VaultRecord {
	uid := uuid.NewString()
	rmap := make(Record)
	var attributes []string
	switch kind {
	case "note":
		attributes = []string{"Name", "Tags"}
	case "file":
		attributes = []string{"Name", "File", "Tags"}
	default: // default login record
		attributes = []string{"Name", "Login", "Password", "URL", "Tags"}
	}
	for _, attr := range attributes {
		rmap[attr] = ""
	}
	return &VaultRecord{ID: uid, Map: rmap, ModificationTime: time.Now()}
}

// Vault represent our vault
type Vault struct {
	Directory        string        // vault directory
	Cipher           string        // vault cipher
	Secret           string        // vault secret
	Verbose          int           // verbose mode
	Records          []VaultRecord // vault records
	ModificationTime time.Time     // vault last modification time
	LastBackup       string        // vault last backup
	Size             int64         // vault size
	Mode             string        // vault mode
	Start            time.Time     // vault expire
}

// AddRecord vault record and return its index
func (v *Vault) AddRecord(kind string) int {
	rec := NewVaultRecord(kind)
	v.Records = append(v.Records, *rec)
	return len(v.Records) - 1
}

// Delete deletes given vault record file from the vault directory
func (v *Vault) DeleteRecordFile(rid string) error {
	// physically delete vault record file
	fname := fmt.Sprintf("%s.%s", filepath.Join(v.Directory, rid), v.Cipher)
	err := os.Remove(fname)
	if err != nil {
		return err
	}
	return nil
}

// DeleteRecord vault record
func (v *Vault) DeleteRecord(rid string) error {
	idx := -1
	for i, rec := range v.Records {
		if rec.ID == rid {
			idx = i
		}
	}
	if idx > -1 {
		remove(v.Records, idx)
	} else {
		msg := fmt.Sprintf("no record %s found in a vault", rid)
		return errors.New(msg)
	}
	return nil
}

// helper function to remove specific entry in vault record list
func remove(s []VaultRecord, i int) []VaultRecord {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// EncryptFile provides ability to encrypt given file name and place into the vault
func (v *Vault) EncryptFile(efile string) {
	data, err := ioutil.ReadFile(efile)
	if err != nil {
		log.Fatalf("unable to read file %s, error %v", efile, err)
	}
	uid := uuid.NewString()
	attachments := []string{efile}
	rmap := make(Record)
	rmap["Data"] = string(data)
	rmap["Name"] = filepath.Base(efile)
	rmap["Tags"] = "file"
	rec := VaultRecord{ID: uid, Map: rmap, Attachments: attachments}
	rec.WriteRecord(v.Directory, v.Secret, v.Cipher, v.Verbose)
	log.Printf("created new vault record %s", rec.ID)
}

// Update vault records
func (v *Vault) Update(rec VaultRecord) error {
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
	err := v.WriteRecord(rec)
	return err
}

// Create provides vault creation functionality
func (v *Vault) Create(vname string) error {
	// setup defaults
	if vname == "" {
		vname = "Primary"
	}

	var vaultDir string
	// construct proper full path
	if v.Directory != "" {
		abs, err := filepath.Abs(v.Directory)
		if err != nil {
			log.Fatal(err)
		}
		v.Directory = abs
	}

	// determine vault location and if it is not provided or does not exists
	// creat $HOME/.ecm area and assign new vault area there
	_, err := os.Stat(v.Directory)
	if v.Directory == "" || os.IsNotExist(err) {
		udir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		vdir := filepath.Join(udir, ".ecm")
		v.Directory = vdir
		err = os.MkdirAll(vdir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// procceed with vault
	vaultDir = filepath.Join(v.Directory, vname)
	v.Directory = vaultDir
	_, err = os.Stat(vaultDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(vaultDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

// Files returns list of vault files
func (v *Vault) Files() ([]string, error) {
	files, err := ioutil.ReadDir(v.Directory)
	if err != nil {
		return []string{}, err
	}
	var out []string
	for _, f := range files {
		if f.Name() != "backups" {
			out = append(out, f.Name())
		}
	}
	return out, nil
}

// Read reads vault records
func (v *Vault) Read() error {
	files, err := ioutil.ReadDir(v.Directory)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: we can parallelize the read from vault area via goroutine pool
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), v.Cipher) {
			if strings.HasPrefix(file.Name(), "backups") {
				continue
			}
			log.Fatalf("File '%s' and Vault cipher '%s' do not match\n", file.Name(), v.Cipher)
		}
		fname := filepath.Join(v.Directory, file.Name())
		rec, err := v.ReadRecord(fname)
		if err != nil {
			log.Fatal("unable to read vault record", fname, " error ", err)
		}
		v.Records = append(v.Records, rec)
	}

	// get vault file info
	finfo, err := os.Stat(v.Directory)
	if err == nil {
		v.Size = finfo.Size()
		v.ModificationTime = finfo.ModTime()
		v.Mode = finfo.Mode().String()
	} else {
		log.Printf("unable to get stat for %s, error %v", v.Directory, err)
	}
	return nil
}

// helper function to read vault and return list of records
func (v *Vault) Write() error {
	// TODO: we can parallelize the read from vault area via goroutine pool
	for _, rec := range v.Records {
		err := rec.WriteRecord(v.Directory, v.Secret, v.Cipher, v.Verbose)
		if err != nil {
			log.Fatalf("unable to write vault record %s, error %v", rec.ID, err)
		}
	}
	return nil
}

// WriteRecord provides write record functionality of vault
func (v *Vault) WriteRecord(rec VaultRecord) error {
	err := rec.WriteRecord(v.Directory, v.Secret, v.Cipher, v.Verbose)
	if err != nil {
		log.Fatalf("unable to write vault record %s, error %v", rec.ID, err)
		return err
	}
	return nil
}

// ReadRecord provides read record functionality of our vault
func (v *Vault) ReadRecord(fname string) (VaultRecord, error) {
	var rec VaultRecord
	// check first if file exsist
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		log.Printf("vault record %s does not exists, will create one", fname)
		_, err := os.Create(fname)
		if err != nil {
			log.Fatal(err)
		}
		return rec, err
	}

	// always keep file safe
	err := os.Chmod(fname, 0600)
	if err != nil {
		log.Println("unable to change file permission of", fname)
	}

	// open file
	file, err := os.Open(fname)
	if err != nil {
		log.Println("unable to open a vault", err)
		return rec, err
	}
	// remember to close the file at the end of the program
	defer file.Close()

	// read data from the record file
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	if v.Cipher != "" {
		data, err = crypt.Decrypt(data, v.Secret, v.Cipher)
		if err != nil {
			if v.Verbose > 0 {
				log.Printf("unable to decrypt data, error %v", err)
			}
			return rec, err
		}
	}

	err = json.Unmarshal(data, &rec)
	if err != nil {
		log.Println("ERROR: unable to unmarshal the data", err)
		return rec, err
	}
	return rec, nil
}

// Find method finds given pattern in our vault and return its index
func (v *Vault) Find(pat string) []VaultRecord {
	var ids []string
	var out []VaultRecord
	for _, rec := range v.Records {
		for key, val := range rec.Map {
			if strings.Contains(key, pat) {
				if !crypt.InList(rec.ID, ids) {
					ids = append(ids, rec.ID)
					out = append(out, rec)
				}
				if v.Verbose > 0 {
					log.Println("match record key")
				}
			} else if matched, err := regexp.MatchString(pat, val); err == nil && matched {
				if v.Verbose > 0 {
					log.Println("matched record value")
				}
				if !crypt.InList(rec.ID, ids) {
					ids = append(ids, rec.ID)
					out = append(out, rec)
				}
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
	info := fmt.Sprintf("vault %s\nLast modified: %s\nSize %s, mode %s\n%d records, encrypted with %s cipher", v.Directory, tstamp, size, mode, nrec, cipher)
	if v.Verbose > 0 {
		log.Println(info)
	}
	return info
}

// Recreate re-creates vault records with new password and cipher
func (v *Vault) Recreate(secret, cipher string) error {
	// make copy of existing vault directory
	tstamp := time.Now().Format(time.RFC3339)
	dstDir := fmt.Sprintf("%s.%s", v.Directory, tstamp)
	err := CopyDir(v.Directory, dstDir)
	if err != nil {
		return err
	}
	log.Printf("Original vault records are saved in %s", dstDir)
	// get all existing records
	for _, rec := range v.Records {
		err := rec.WriteRecord(v.Directory, secret, cipher, v.Verbose)
		if err != nil {
			return err
		}
		// delete record from the vault
		err = v.DeleteRecord(rec.ID)
		if err != nil {
			return err
		}
		// delete existing record file from vault directory
		err = v.DeleteRecordFile(rec.ID)
		if err != nil {
			return err
		}
	}
	// change vault secret and cipher
	v.Secret = secret
	v.Cipher = cipher
	log.Printf("Vault changed and re-encrypted all records in %s using cipher %s", v.Directory, v.Cipher)
	return nil
}

// Import allows to import vault records to a given file
// CSV, JSON or ECM-JSON data-format are supported
func (v *Vault) Import(fname, oname string) error {
	// open file
	f, err := os.Open(fname)
	if err != nil {
		return err
	}

	// remember to close the file at the end of the program
	defer f.Close()

	var records []VaultRecord
	if fname == "ecm.json" {
		data, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &records)
		if err != nil {
			return err
		}

	} else if strings.HasSuffix(fname, "csv") {

		// read csv values using csv.Reader
		reader := csv.NewReader(f)
		var headers []string
		for {
			values, err := reader.Read()
			if len(headers) == 0 {
				headers = values
				continue
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			vRecord := NewVaultRecord("login")
			for idx := range values {
				if values[idx] == "" {
					continue
				}
				key := recordAttribute(headers[idx])
				vRecord.Map[key] = values[idx]
			}
			if v.Verbose > 0 {
				log.Println("Import VaultRecord\n", vRecord.String())
			}
			records = append(records, *vRecord)
		}

	} else if strings.HasSuffix(fname, "json") {

		data, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		var jsonRecords []map[string]any
		err = json.Unmarshal(data, &jsonRecords)
		if err != nil {
			return err
		}
		for _, rec := range jsonRecords {
			vRecord := NewVaultRecord("login")
			for key, val := range rec {
				vRecord.Map[key] = fmt.Sprintf("%svs", val)
			}
			if v.Verbose > 0 {
				log.Println("Import VaultRecord\n", vRecord.String())
			}
			records = append(records, *vRecord)
		}
	}

	if oname != "" {
		// check if our destination is a vault
		if oname == v.Directory {
			for _, rec := range records {
				err := rec.WriteRecord(v.Directory, v.Secret, v.Cipher, v.Verbose)
				if err != nil {
					log.Fatalf("unable to write vault record %s, error %v", rec.ID, err)
				}
			}
			return nil
		}

		// otherwise write records to destination file
		file, err := os.Open(oname)
		if err != nil {
			return err
		}
		// remember to close the file at the end of the program
		defer file.Close()
		data, err := json.MarshalIndent(records, "", "   ")
		if err != nil {
			return err
		}
		err = os.WriteFile(oname, data, 0755)
		return err
	}
	return nil
}

// Export allows to export vault records in JSON data format to a given file
func (v *Vault) Export(fname string) error {
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.MarshalIndent(v.Records, "", "   ")
	if err != nil {
		return err
	}
	err = os.WriteFile(fname, data, 0755)
	return err
}
