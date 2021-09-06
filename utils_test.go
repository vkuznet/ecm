package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/atotto/clipboard"
)

// TestDecryptInputToFile function
func TestDecryptInputToFile(t *testing.T) {
	password := "test"
	data := []byte(`{"attr": "test"}`)
	cipher := "aes"
	attr := ""
	edata, err := encrypt(data, password, cipher)
	tmpFile, err := ioutil.TempFile(os.TempDir(), "input-")
	if err != nil {
		t.Error(err.Error())
	}
	defer os.Remove(tmpFile.Name())
	// write our data to temp file
	if _, err = tmpFile.Write(edata); err != nil {
		t.Error("Failed to write to temporary file", err)
	}
	tmpFile.Close()
	// create output file
	outTmpFile, err := ioutil.TempFile(os.TempDir(), "ouput-")
	if err != nil {
		t.Error(err.Error())
	}
	defer os.Remove(outTmpFile.Name())
	// compare reading from file and writing to a file
	decryptInput(tmpFile.Name(), password, cipher, outTmpFile.Name(), attr)
	// read data from output file
	res, err := os.ReadFile(outTmpFile.Name())
	if err != nil {
		t.Error(err.Error())
	}
	if string(res) != string(data) {
		t.Errorf("wrong encrypted data written to out file, expect '%s' result '%s'", string(data), string(res))
	}
}

// TestDecryptInputToClipboard function
func TestDecryptInputToClipboard(t *testing.T) {
	if os.Getenv("SKIP_CLIPBOARD_TEST") == "1" {
		log.Println("skip clipboard test")
		return
	}
	password := "test"
	rec := make(Record)
	rec["attr"] = "test"
	vrec := VaultRecord{ID: "123", Map: rec}
	cipher := "aes"
	data, err := json.Marshal(vrec)
	if err != nil {
		t.Error(err.Error())
	}
	edata, err := encrypt(data, password, cipher)
	if err != nil {
		t.Error(err.Error())
	}
	tmpFile, err := ioutil.TempFile(os.TempDir(), "input-")
	if err != nil {
		t.Error(err.Error())
	}
	defer os.Remove(tmpFile.Name())
	// write our data to temp file
	if _, err = tmpFile.Write(edata); err != nil {
		t.Error("Failed to write to temporary file", err)
	}
	tmpFile.Close()
	// now we can test clipboard read/write operation
	decryptInput(tmpFile.Name(), password, cipher, "clipboard", "")
	cdata, err := clipboard.ReadAll()
	if err != nil {
		log.Fatal("unable to copy to clipboard, error ", err)
	}
	if string(cdata) != string(data) {
		t.Errorf("wrong data written to clipboard, expect '%s' result '%s'", string(data), string(cdata))
	}

	// read certain attribute from clipboard
	decryptInput(tmpFile.Name(), password, cipher, "clipboard", "attr")
	cdata, err = clipboard.ReadAll()
	if err != nil {
		log.Fatal("unable to copy to clipboard, error ", err)
	}
	// clipboard should now only contain valut of {"attr":"test"} record
	if string(cdata) != "test" {
		t.Errorf("wrong data written to clipboard, expect '%s' result '%s'", string(data), string(cdata))
	}

}

// TestGenToken function
func TestGenToken(t *testing.T) {
	salt := "test"
	cipher := "aes"
	token, err := encryptToken(salt, cipher)
	if err != nil {
		t.Error(err.Error())
	}
	dtoken, err := decryptToken(token, salt, cipher)
	if err != nil {
		t.Error(err.Error())
	}
	if token != dtoken {
		t.Errorf("wrong token, expect '%s' got '%s'", token, dtoken)
	}
}
