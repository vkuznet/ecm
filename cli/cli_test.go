package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/vkuznet/ecm/crypt"
	vt "github.com/vkuznet/ecm/vault"
)

// helper function to generate CSV records
func csvRecords(nrec int) []byte {
	var out []string
	out = append(out, "Name,Login,Password,Tag")
	for i := 0; i < nrec; i++ {
		data := fmt.Sprintf("name-%d,login-%d,pass-%d,tag-%d", i, i, i, i)
		out = append(out, data)
	}
	return []byte(strings.Join(out, "\n"))
}

// TestCLI function
func TestCLI(t *testing.T) {
	cipher := "aes"
	secret := "test"
	nrec := 3
	verbose := 2
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// recover os.Exit(0) panics
	/*
		defer func() {
			if err := recover(); err != nil {
				if strings.Contains(fmt.Sprintf("%v", err), "os.Exit(0)") {
					// we will skip os.Exit(0)
				} else {
					panic(err)
				}
			}
		}()
	*/

	// generate test.csv file with random records
	// we should use test.csv file as it is checked in vault.go
	csvFile, err := ioutil.TempFile(os.TempDir(), "test.csv-")
	if err != nil {
		t.Error(err.Error())
	}
	log.Println("create cvs file", csvFile.Name())
	defer os.Remove(csvFile.Name())
	// write our data to temp file
	csvData := csvRecords(nrec)
	if _, err = csvFile.Write(csvData); err != nil {
		t.Error("Failed to write to temporary file", err)
	}
	if info, err := os.Stat(csvFile.Name()); err == nil {
		log.Printf("CSV file %+v\n", info)
	} else {
		log.Fatalf("Unable to stat file %s, error %v", csvFile.Name(), err)
	}
	// close csv file we will need to read from it
	csvFile.Close()

	// read records from csvFile
	if file, e := os.Open(csvFile.Name()); e == nil {
		data, err := io.ReadAll(file)
		if err == nil {
			log.Println("csv data:\n", string(data))
		} else {
			log.Fatalf("unable to read csvFile, %v", err)
		}
	}

	// generate ecm file which will hold converted csv records
	// we should use ecm.json name as it is checked in vault.go
	ecmFile, err := ioutil.TempFile(os.TempDir(), "ecm.json-")
	if err != nil {
		t.Error(err.Error())
	}
	log.Println("create json file", ecmFile.Name())
	defer os.Remove(ecmFile.Name())
	// close the temp file we will need to write to it
	ecmFile.Close()

	// initialize our vault
	vname, err := os.MkdirTemp(os.TempDir(), ".ecm")
	if err != nil {
		log.Fatalf("unable to create temp dir, error %v", err)
	}
	log.Println("create vault at", vname)
	vault := vt.Vault{
		Secret:    secret,
		Directory: vname,
		Cipher:    crypt.GetCipher(cipher),
		Verbose:   verbose,
		Start:     time.Now()}
	err = vault.Create(vname)
	if err != nil {
		log.Fatalf("unable to create vault, error %v", err)
	}

	var efile, dfile, add, pat, rid, edit, pcopy, export, vimport, sync string
	var recreate, info bool
	log.Println("emulate `ecm -import test.csv -export ecm.json`")
	vimport = csvFile.Name()
	export = ecmFile.Name()
	log.Printf("will export records from %s (csv file) to %s (ecm json file)", vimport, export)
	cli(&vault,
		efile, dfile, add, pat, rid, edit, pcopy, export, vimport, sync,
		recreate, info,
		verbose,
	)

	// read records from csvFile
	if file, e := os.Open(ecmFile.Name()); e == nil {
		data, err := io.ReadAll(file)
		if err == nil {
			log.Println("ecm json data:\n", string(data))
		} else {
			log.Fatalf("unable to read csvFile, %v", err)
		}
	}

	log.Println("emulate `ecm -import ecm.json -export <vault>`")
	vimport = ecmFile.Name()
	export = vname
	log.Printf("will export records from %s (csv file) to %s (vault.Directory %s)", vimport, export, vault.Directory)
	cli(&vault,
		efile, dfile, add, pat, rid, edit, pcopy, export, vimport, sync,
		recreate, info,
		verbose,
	)

	// list vault directory
	files, err := ioutil.ReadDir(vname)
	if err == nil {
		log.Printf("vault %s content\n", vname)
		for _, file := range files {
			fmt.Println(file.Name())
		}
	}

	// re-read vault records
	err = vault.Read()
	if err != nil {
		log.Fatalf("unable to read vault, error %v", err)
	}

	// get vault info
	log.Println(vault.Info())

	// run cli -pat name-327
	vimport = ""
	export = ""
	pat = "name-1"
	cli(&vault,
		efile, dfile, add, pat, rid, edit, pcopy, export, vimport, sync,
		recreate, info,
		verbose,
	)
}
