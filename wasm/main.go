package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"syscall/js"
	"time"

	"github.com/vkuznet/gpm/crypt"
	//dom "honnef.co/go/js/dom/v2"
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
	ID       string
	Login    string
	Password string
	Name     string
	Note     string
	Tags     string
	URL      string
}

// RecordMap type defines our GPM record map
type RecordMap map[string]LoginRecord

// RecordsManager holds GPM records
type RecordsManager struct {
	Map           RecordMap
	RenewInterval int64
	Expire        int64
}

// global records manager which holds all vault records
var recordsManager *RecordsManager

// helper function to get GPM records
func (mgr *RecordsManager) update(url, cipher, password string) error {
	if recordsManager.Map == nil || recordsManager.Expire < time.Now().Unix() {
		rmap, err := getRecords(url, cipher, password)
		mgr.Map = rmap
		mgr.Expire = time.Now().Unix() + mgr.RenewInterval
		return err
	}
	return nil
}

// helper function to get GPM records from given URL
func getRecords(url, cipher, password string) (RecordMap, error) {
	rmap := make(RecordMap)

	// Make the HTTP request
	res, err := http.DefaultClient.Get(url)
	if err != nil {
		return rmap, err
	}
	defer res.Body.Close()

	// Read the response body
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return rmap, err
	}
	// records represent list of file names
	var records [][]byte
	err = json.Unmarshal(data, &records)
	if err != nil {
		return rmap, err
	}
	for _, rec := range records {
		data, err := crypt.Decrypt(rec, password, cipher)
		if err != nil {
			return rmap, err
		}
		var vrec VaultRecord
		err = json.Unmarshal(data, &vrec)
		lrec := LoginRecord{
			ID:       vrec.ID,
			Login:    vrec.Map["Login"],
			Password: vrec.Map["Password"],
			Note:     vrec.Map["Note"],
			Name:     vrec.Map["Name"],
			Tags:     vrec.Map["Tags"],
			URL:      vrec.Map["URL"],
		}
		rmap[vrec.ID] = lrec
	}
	return rmap, nil
}

// main function sets JS "gpmDecode" function to call "decodeWrapper" Go counterpart
func main() {
	// log time, filename, and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// initialize records manager with renew interval of 60 seconds
	recordsManager = &RecordsManager{RenewInterval: 60}
	// define out decode JS function to be bound to decoreWrapper Go counterpart
	//js.Global().Set("Lock", lockWrapper())
	js.Global().Set("getLogin", loginWrapper())
	js.Global().Set("getPassword", passwordWrapper())
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

// wrapper function to return password for given record id
func passwordWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		rid := args[0].String()
		if lrec, ok := recordsManager.Map[rid]; ok {
			return lrec.Password
		}
		return ""
	})
}

// wrapper function to return password for given record id
func loginWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		rid := args[0].String()
		if lrec, ok := recordsManager.Map[rid]; ok {
			return lrec.Login
		}
		return ""
	})
}

// RecordsHandler handles asynchronously HTTP requests
func RecordsHandler(url, cipher, passphrase string, args []js.Value) {
	resolve := args[0]
	reject := args[1]

	err := recordsManager.update(url, cipher, passphrase)
	if err != nil {
		ErrorHandler(reject, err)
		return
	}

	var rids []string
	document := js.Global().Get("document")
	docRecords := document.Call("getElementById", "records")
	docRecords.Set("innerHTML", "")
	ul := document.Call("createElement", "ul")
	ul.Call("setAttribute", "class", "records")
	docRecords.Call("appendChild", ul)
	for key, lrec := range recordsManager.Map {
		name := lrec.Name
		login := lrec.Login
		password := lrec.Password
		rurl := lrec.URL
		rids = append(rids, key)

		// construct frontend UI
		li := document.Call("createElement", "li")
		li.Call("setAttribute", "class", "item")
		ul.Call("appendChild", li)
		nameDiv := document.Call("createElement", "div")
		nameDiv.Set("innerHTML", "Name: "+name)
		loginDiv := document.Call("createElement", "div")
		loginDiv.Set("innerHTML", "Login: "+login)
		passDiv := document.Call("createElement", "div")
		pid := "pid-" + key
		passDiv.Set("id", pid)
		passDiv.Call("setAttribute", "class", "hide")

		// add buttons
		buttons := document.Call("createElement", "div")
		buttons.Call("setAttribute", "class", "button-right")

		// add show button
		button := document.Call("createElement", "button")
		bid := "bid-" + key
		button.Set("id", bid)
		button.Call("setAttribute", "class", "label")
		button.Set("innerHTML", "Show password")
		var callback js.Func
		callback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			doc := document.Call("getElementById", pid)
			button := document.Call("getElementById", bid)
			if button.Get("innerHTML").String() == "Show password" {
				doc.Call("setAttribute", "class", "show-inline")
				doc.Set("innerHTML", "Password: "+password)
				button.Set("innerHTML", "Hide password")
			} else {
				doc.Call("setAttribute", "class", "show-inline")
				doc.Set("innerHTML", "")
				button.Set("innerHTML", "Show password")
			}
			return nil
		})
		button.Call("addEventListener", "click", callback)
		buttons.Call("appendChild", button)

		// add autofill button
		button = document.Call("createElement", "button")
		aid := "autofill-" + key
		button.Set("id", aid)
		button.Call("setAttribute", "class", "label autofill is-bold")
		button.Set("innerHTML", "Autofill")
		button.Set("RecordID", key)
		buttons.Call("appendChild", button)

		siteDiv := document.Call("createElement", "div")
		siteDiv.Set("innerHTML", "URL: "+rurl)

		li.Call("append", nameDiv)
		li.Call("append", siteDiv)
		li.Call("append", loginDiv)
		li.Call("append", passDiv)
		li.Call("append", buttons)
	}

	adata, err := json.Marshal(rids)
	if err != nil {
		ErrorHandler(reject, err)
		return
	}
	// "data" is a byte slice, so we need to convert it to a JS Uint8Array object
	arrayConstructor := js.Global().Get("Uint8Array")
	dataJS := arrayConstructor.New(len(adata))
	js.CopyBytesToJS(dataJS, adata)

	// Create a Response object and pass the data
	responseConstructor := js.Global().Get("Response")
	response := responseConstructor.New(dataJS)

	// Resolve the Promise
	resolve.Invoke(response)
}
