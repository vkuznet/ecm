package main

import (
	"io/ioutil"
	"os"
	"testing"
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
