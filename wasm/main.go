package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"syscall/js"

	"github.com/vkuznet/gpm/crypt"
)

// Record represent map of key-valut pairs
type Record map[string]string

// VaultRecord represents vault record subset suitable for web UI
type VaultRecord struct {
	ID  string // record ID
	Map Record // record map (key-vault pairs)
}

// helper function to get cipher name from given file
func getCipher(fname string) string {
	cipher := "aes"
	arr := strings.Split(fname, ".")
	if len(arr) > 1 {
		cipher = arr[1]
	}
	return cipher
}

// main function sets JS "gpmDecode" function to call "decodeWrapper" Go counterpart
func main() {
	// log time, filename, and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// define out decode JS function to be bound to decoreWrapper Go counterpart
	js.Global().Set("decode", decodeWrapper())
	js.Global().Set("records", recordsWrapper())
	<-make(chan bool)
}

// recordsWraper function performs business logic, i.e.
// it recordss given input obtained from JS upstream code
func recordsWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		passphrase := args[0].String()
		// TODO: we need to construct URL somehow and choose vault
		url := "http://127.0.0.1:8888/vault/Primary/records"

		// Create and return the Promise object
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			go RecordsHandler(url, passphrase, args)
			return nil
		})
		// define where we should put our data
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

// decodeWraper function performs business logic, i.e.
// it decodes given input obtained from JS upstream code
func decodeWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// here input should be vault file name, e.g.
		// "acb8a9f7-6140-42d2-bb32-f730f7ab572f.aes"
		input := args[0].String()
		passphrase := args[1].String()

		// TODO: replace how we'll accept cipher, fname, passphrase
		fname := strings.Trim(input, " ")
		cipher := getCipher(fname)
		// TODO: we need to construct URL somehow
		url := fmt.Sprintf("http://127.0.0.1:8888/vault/Primary/%s", fname)

		// Create and return the Promise object
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			go RequestHandler(url, passphrase, cipher, args)
			return nil
		})
		// define where we should put our data
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

// ErrorHandler handles JS errors
func ErrorHandler(reject js.Value, err error) {
	// Handle errors here too
	errorConstructor := js.Global().Get("Error")
	errorObject := errorConstructor.New(err.Error())
	reject.Invoke(errorObject)
}

// RequestHandler handles asynchronously HTTP requests
func RequestHandler(url, passphrase, cipher string, args []js.Value) {
	resolve := args[0]
	reject := args[1]

	// Make the HTTP request
	res, err := http.DefaultClient.Get(url)
	if err != nil {
		ErrorHandler(reject, err)
		return
	}
	defer res.Body.Close()

	// Read the response body
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ErrorHandler(reject, err)
		return
	}
	rdata, err := crypt.Decrypt(data, passphrase, cipher)
	if err != nil {
		ErrorHandler(reject, err)
		return
	}

	// "data" is a byte slice, so we need to convert it to a JS Uint8Array object
	arrayConstructor := js.Global().Get("Uint8Array")
	dataJS := arrayConstructor.New(len(rdata))
	js.CopyBytesToJS(dataJS, rdata)

	// Create a Response object and pass the data
	responseConstructor := js.Global().Get("Response")
	response := responseConstructor.New(dataJS)

	// Resolve the Promise
	resolve.Invoke(response)
}

// RecordsHandler handles asynchronously HTTP requests
func RecordsHandler(url, passphrase string, args []js.Value) {
	resolve := args[0]
	reject := args[1]

	// Make the HTTP request
	res, err := http.DefaultClient.Get(url)
	if err != nil {
		ErrorHandler(reject, err)
		return
	}
	defer res.Body.Close()

	// Read the response body
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ErrorHandler(reject, err)
		return
	}
	// records represent list of file names
	var records [][]byte
	err = json.Unmarshal(data, &records)
	if err != nil {
		ErrorHandler(reject, err)
		return
	}
	// TODO: we need to get from client cipher name of the vault
	cipher := "aes"
	rmap := make(map[string]VaultRecord)
	for _, rec := range records {
		data, err := crypt.Decrypt(rec, passphrase, cipher)
		if err != nil {
			log.Println("fail to decrypt record, error", err)
			ErrorHandler(reject, err)
			return
		}
		var vrec VaultRecord
		err = json.Unmarshal(data, &vrec)
		rmap[vrec.ID] = vrec
	}
	rdata, err := json.Marshal(rmap)
	if err != nil {
		ErrorHandler(reject, err)
		return
	}

	// "data" is a byte slice, so we need to convert it to a JS Uint8Array object
	arrayConstructor := js.Global().Get("Uint8Array")
	dataJS := arrayConstructor.New(len(rdata))
	js.CopyBytesToJS(dataJS, rdata)

	// Create a Response object and pass the data
	responseConstructor := js.Global().Get("Response")
	response := responseConstructor.New(dataJS)

	// Resolve the Promise
	resolve.Invoke(response)
}
