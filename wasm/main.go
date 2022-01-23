package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"syscall/js"
	"time"

	"github.com/vkuznet/ecm/crypt"
	vt "github.com/vkuznet/ecm/vault"
	//dom "honnef.co/go/js/dom/v2"
)

// Record represent map of key-valut pairs
// type Record map[string]string

// VaultRecord represents vault record subset suitable for web UI
// type VaultRecord struct {
//     ID  string // record ID
//     Map Record // record map (key-vault pairs)
// }

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

// RecordMap type defines our ECM record map
type RecordMap map[string]LoginRecord

// RecordsManager holds ECM records
type RecordsManager struct {
	Map           RecordMap
	RenewInterval int64
	Expire        int64
}

// global records manager which holds all vault records
var recordsManager *RecordsManager

// helper function to get ECM records
func (mgr *RecordsManager) update(url, cipher, password string) error {
	if recordsManager.Map == nil || recordsManager.Expire < time.Now().Unix() {
		rmap, err := getRecords(url, cipher, password)
		mgr.Map = rmap
		mgr.Expire = time.Now().Unix() + mgr.RenewInterval
		return err
	}
	return nil
}

// helper function to get ECM records from given URL
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
		var vrec vt.VaultRecord
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

// main function sets JS functions
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

	js.Global().Set("addRecord", actionWrapper("login_record"))
	js.Global().Set("addJsonRecord", actionWrapper("json_record"))
	js.Global().Set("addNote", actionWrapper("note_record"))
	js.Global().Set("addCard", actionWrapper("card_record"))
	js.Global().Set("addVault", actionWrapper("new_vault"))
	js.Global().Set("syncHosts", actionWrapper("sync_hosts"))
	js.Global().Set("uploadFile", actionWrapper("upload_file"))
	js.Global().Set("showRecords", actionWrapper("show_records"))
	js.Global().Set("generatePassword", actionWrapper("gen_password"))

	//     js.Global().Call("showRecords")

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

// wrapper function to add new vault
func actionWrapper(action string) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Create and return the Promise object
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			go ActionHandler(action, args)
			return nil
		})
		// define where we should put our data
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

// ActionHandler handles vault action
func ActionHandler(action string, args []js.Value) {
	resolve := args[0]
	reject := args[1]

	// perform action with our action
	var data []byte
	var err error
	if action == "show_records" {
		data, err = showRecords()
	} else if action == "login_record" {
		data, err = loginRecord()
	} else if action == "json_record" {
		data, err = jsonRecord()
	} else if action == "card_record" {
		data, err = cardRecord()
	} else if action == "note_record" {
		data, err = noteRecord()
	} else if action == "upload_file" {
		data, err = uploadFile()
	} else if action == "sync_host" {
		data, err = syncHosts()
	} else if action == "new_vault" {
		data, err = createVault()
	} else if action == "gen_password" {
		data, err = newPassword()
		//     } else if action == "init_mgr" {
		//         requestManager = initRequestManager()
	} else {
		data, err = defaultAction()
	}
	if err != nil {
		ErrorHandler(reject, err)
		return
	}

	// "data" is a byte slice, so we need to convert it to a JS Uint8Array object
	arrayConstructor := js.Global().Get("Uint8Array")
	dataJS := arrayConstructor.New(len(data))
	js.CopyBytesToJS(dataJS, data)

	// Create a Response object and pass the data
	responseConstructor := js.Global().Get("Response")
	response := responseConstructor.New(dataJS)

	// Resolve the Promise
	resolve.Invoke(response)
}

// wrapper function to generate new password
func genPasswordWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		const voc string = "abcdfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		const numbers string = "0123456789"
		//         const symbols string = "!@#$%&*+_-="
		size := args[0].Int()
		chars := voc + numbers
		password := generatePassword(size, chars)
		document := js.Global().Get("document")
		docRecords := document.Call("getElementById", "new-password")
		docRecords.Set("innerHTML", password)
		return ""
	})
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

// update records within DOM document
func updateRecords(url, cipher, passphrase string, extention bool) ([]string, error) {

	var rids []string
	err := recordsManager.update(url, cipher, passphrase)
	if err != nil {
		return rids, err
	}

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
		tags := lrec.Tags
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
		if extention {
			button.Call("setAttribute", "class", "label is-focus is-pointer")
		} else {
			button.Call("setAttribute", "class", "button is-secondary is-small")
		}
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
		if extention && password != "" {
			buttons.Call("appendChild", button)
		} else {
			// create grid with two columns and add button to second one
			grid := document.Call("createElement", "div")
			grid.Call("setAttribute", "class", "is-row")
			div := document.Call("createElement", "div")
			div.Call("setAttribute", "class", "is-col is-80")
			img := document.Call("createElement", "img")
			if strings.Contains(tags, "file") {
				img.Set("src", "static/images/file-32.png")
			} else {
				img.Set("src", "static/images/notes-32.png")
			}
			div.Call("appendChild", img)
			txt := document.Call("createElement", "span")
			txt.Set("innerHTML", fmt.Sprintf("&nbsp;&nbsp; %s", key))
			div.Call("appendChild", txt)
			grid.Call("appendChild", div)
			div = document.Call("createElement", "div")
			div.Call("setAttribute", "class", "is-col is-20")
			div.Call("appendChild", button)
			grid.Call("appendChild", div)
			// add grid to buttons
			buttons.Call("appendChild", grid)
		}

		// add autofill button
		if extention {
			button = document.Call("createElement", "button")
			aid := "autofill-" + key
			button.Set("id", aid)
			button.Call("setAttribute", "class", "label autofill is-bold")
			button.Set("innerHTML", "Autofill")
			button.Set("RecordID", key)
			buttons.Call("appendChild", button)
		}

		siteDiv := document.Call("createElement", "div")
		siteDiv.Set("innerHTML", "URL: "+rurl)

		if !extention {
			li.Call("append", buttons)
		}
		li.Call("append", nameDiv)
		li.Call("append", siteDiv)
		li.Call("append", loginDiv)
		li.Call("append", passDiv)
		if extention {
			li.Call("append", buttons)
		}
	}
	return rids, nil
}

// RecordsHandler handles asynchronously HTTP requests
func RecordsHandler(url, cipher, passphrase string, args []js.Value) {
	resolve := args[0]
	reject := args[1]

	// place records within DOM page
	extention := true
	rids, err := updateRecords(url, cipher, passphrase, extention)
	if err != nil {
		ErrorHandler(reject, err)
		return
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
