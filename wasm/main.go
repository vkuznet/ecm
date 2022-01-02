package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
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

// LoginRecord represent login credentials
type LoginRecord struct {
	Login    string
	Password string
	Name     string
	Note     string
	Tags     string
	URL      string
}

// main function sets JS "gpmDecode" function to call "decodeWrapper" Go counterpart
func main() {
	// log time, filename, and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// define out decode JS function to be bound to decoreWrapper Go counterpart
	js.Global().Set("records", recordsWrapper())
	<-make(chan bool)
}

// recordsWraper function performs business logic, i.e.
// it recordss given input obtained from JS upstream code
func recordsWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		server := args[0].String()
		vault := args[1].String()
		cipher := args[2].String()
		password := args[3].String()
		// construct URL, e.g.
		// http://127.0.0.1:8888/vault/Primary/records
		url := fmt.Sprintf("%s/vault/%s/records", server, vault)

		// Create and return the Promise object
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			go RecordsHandler(url, cipher, password, args)
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
func RecordsHandler(url, cipher, passphrase string, args []js.Value) {
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
	rmap := make(map[string]LoginRecord)
	for _, rec := range records {
		data, err := crypt.Decrypt(rec, passphrase, cipher)
		if err != nil {
			log.Println("fail to decrypt record, error", err)
			ErrorHandler(reject, err)
			return
		}
		var vrec VaultRecord
		err = json.Unmarshal(data, &vrec)
		rmap[vrec.ID] = LoginRecord{
			Login:    vrec.Map["Login"],
			Password: vrec.Map["Password"],
			Note:     vrec.Map["Note"],
			Name:     vrec.Map["Name"],
			Tags:     vrec.Map["Tags"],
			URL:      vrec.Map["URL"],
		}
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

// return MAC address of network interface
// to convert to string use
// fmt.Sprintf("%16.16X", macAddress())
// https://gist.github.com/tsilvers/085c5f39430ced605d970094edf167ba
func macAddress() uint64 {
	interfaces, err := net.Interfaces()
	if err != nil {
		return uint64(0)
	}

	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {

			// Skip locally administered addresses
			if i.HardwareAddr[0]&2 == 2 {
				continue
			}

			var mac uint64
			for j, b := range i.HardwareAddr {
				if j >= 8 {
					break
				}
				mac <<= 8
				mac += uint64(b)
			}

			return mac
		}
	}

	return uint64(0)
}
