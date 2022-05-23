package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/vkuznet/ecm/crypt"
	utils "github.com/vkuznet/ecm/utils"
	vt "github.com/vkuznet/ecm/vault"
	// clone of "code.google.com/p/rsc/qr" which no longer available
	// "github.com/vkuznet/rsc/qr"
	// imaging library
)

const (
	separator = "---\n" // used in ecm data format
)

// helper function to extract cipher name from file extension
func fileCipher(fname string) string {
	arr := strings.Split(fname, ".")
	cipher := strings.Split(arr[len(arr)-1], "-")[0]
	if !utils.InList(cipher, crypt.SupportedCiphers) {
		log.Fatalf("given cipher %s is not supported, please use one from the following %v", cipher, crypt.SupportedCiphers)
	}
	return cipher
}

// helper function to decrypt given input (file or stdin)
func decryptInput(fname, password, cipher, write, attr string) {
	var err error
	if cipher == "" {
		cipher = fileCipher(fname)
	}
	if !utils.InList(cipher, crypt.SupportedCiphers) {
		log.Fatalf("given cipher %s is not supported, please use one from the following %v", cipher, crypt.SupportedCiphers)
	}
	var data []byte
	if fname == "-" { // stdin
		var input string
		reader := bufio.NewReader(os.Stdin)
		input, err = reader.ReadString('\n')
		data = []byte(input)
	} else {
		data, err = os.ReadFile(fname)
	}
	if err != nil {
		panic(err)
	}
	data, err = crypt.Decrypt(data, password, cipher)
	if err != nil {
		log.Fatal(err)
	}
	if attr != "" {
		var rec vt.VaultRecord
		err := json.Unmarshal(data, &rec)
		if err != nil {
			log.Fatal("unable to unarmashal vault record", err)
		}
		val, ok := rec.Map[attr]
		if ok {
			data = []byte(val)
		} else {
			log.Fatalf("unable to extract attribute '%s' from the record %+v", attr, rec)
		}
	}
	if write == "stdout" {
		fmt.Println(string(data))
	} else if write == "clipboard" {
		if err := clipboard.WriteAll(string(data)); err != nil {
			log.Fatal("unable to copy to clipboard, error", err)
		}
	} else {
		err := os.WriteFile(write, data, 0755)
		if err != nil {
			log.Fatal("unable to data to output file", err)
		}
	}
}
