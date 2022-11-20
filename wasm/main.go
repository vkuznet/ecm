package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
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
		if err != nil {
			rmap := make(RecordMap)
			rid := "12345"
			lrec := LoginRecord{
				ID:    rid,
				Login: "Error",
				Name:  fmt.Sprintf("url: %s, error: %v", url, err),
			}
			rmap[rid] = lrec
		}
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
	client, err := httpClient()
	if err != nil {
		return rmap, err
	}

	// get results from our url
	res, err := client.Get(url)
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
	js.Global().Set("uploadFile", uploadFileWrapper())

	js.Global().Set("addRecord", actionWrapper("login_record"))
	js.Global().Set("addJsonRecord", actionWrapper("json_record"))
	js.Global().Set("addNote", actionWrapper("note_record"))
	js.Global().Set("addCard", actionWrapper("card_record"))
	js.Global().Set("addVault", actionWrapper("new_vault"))
	js.Global().Set("syncHosts", actionWrapper("sync_hosts"))
	js.Global().Set("showRecords", actionWrapper("show_records"))
	js.Global().Set("generatePassword", actionWrapper("gen_password"))

	//     js.Global().Call("showRecords")

	<-make(chan bool)
}

// VaultPassword keeps copy of vault password
var VaultPassword string

// recordsWraper function performs business logic, i.e.
// it recordss given input obtained from JS upstream code
func recordsWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		server := args[0].String()
		vault := args[1].String()
		cipher := args[2].String()
		password := args[3].String()
		pattern := args[4].String()
		pageUrl := args[5].String()

		// construct URL, e.g.
		// http://127.0.0.1:8888/vault/Primary/records
		// for exact IP:PORT values please consult
		// extension/manifest.json and extension/index.html
		// and/or set it up in settings menu of wasm extension
		url := fmt.Sprintf("%s/vault/%s/records", server, vault)

		if VaultPassword == "" {
			VaultPassword = password
		}

		// Create and return the Promise object
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			go RecordsHandler(url, cipher, VaultPassword, pattern, pageUrl, args)
			return nil
		})
		// define where we should put our data
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

// uploadFileWraper function performs file uploadFile
func uploadFileWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fname := args[0].String()
		size := args[1].Int()
		ftype := args[2].String()
		content := args[3].String() // adjust accordingly with utils.js:readFile function
		// Create and return the Promise object
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			go uploadFile(fname, size, ftype, content)
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
	} else if action == "sync_host" {
		data, err = syncHosts()
	} else if action == "new_vault" {
		data, err = createVault()
	} else if action == "gen_password" {
		data, err = newPassword()
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

// helper function to match url with pattern
func urlMatch(url, pattern string) bool {
	purl := strings.Replace(url, "http://", "", -1)
	purl = strings.Replace(purl, "https://", "", -1)
	arr := strings.Split(purl, "/")
	for _, v := range strings.Split(arr[0], ".") {
		if matched, err := regexp.MatchString(v, pattern); err == nil && matched {
			return true
		}
	}
	return false
}

// helper function to match url parts with LoginRecord
func urlMatchRecord(url string, rec LoginRecord) bool {
	if urlMatch(url, rec.Name) ||
		urlMatch(url, rec.URL) ||
		urlMatch(url, rec.Tags) ||
		urlMatch(url, rec.Note) {
		return true
	}
	return false
}

// update records within DOM document
func updateRecords(url, cipher, passphrase, pattern, pageUrl string, extention bool) ([]string, error) {

	var rids []string
	err := recordsManager.update(url, cipher, passphrase)
	if err != nil {
		return rids, err
	}

	// get pageUrl
	document := js.Global().Get("document")
	// TODO: I should find a way to obtain page URL
	// in wasm there is no way to access page url
	//     pageUrl := ""

	docRecords := document.Call("getElementById", "records")
	docRecords.Set("innerHTML", "")
	ul := document.Call("createElement", "ul")
	ul.Call("setAttribute", "class", "records")
	docRecords.Call("appendChild", ul)
	count := 0
	nrec := 5 // total number of records to show
	var mkeys []string
	for k, _ := range recordsManager.Map {
		mkeys = append(mkeys, k)
	}
	sort.Strings(mkeys)

	//         for key, lrec := range recordsManager.Map {
	for _, key := range mkeys {
		lrec, ok := recordsManager.Map[key]
		if !ok {
			continue
		}
		count += 1
		// skip records which does not match page url
		if pageUrl != "" && pattern == "" {
			if !urlMatchRecord(pageUrl, lrec) {
				continue
			}
		}
		name := lrec.Name
		login := lrec.Login
		password := lrec.Password
		rurl := lrec.URL
		tags := lrec.Tags

		if pattern != "" {
			// TODO: we may need to fetch only appropriate records from server
			// instead of filtering here
			if !(strings.Contains(name, pattern) ||
				strings.Contains(login, pattern) ||
				strings.Contains(password, pattern) ||
				strings.Contains(rurl, pattern) ||
				strings.Contains(tags, pattern)) {
				continue
			}
		}

		if count > nrec {
			continue
		}
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
	if count > nrec {
		moreDiv := document.Call("createElement", "div")
		if extention {
			moreDiv.Set("innerHTML", fmt.Sprintf("Total vault records: %d<br>URL %s<br>Pattern: %s", count, pageUrl, pattern))
		} else {
			moreDiv.Set("innerHTML", fmt.Sprintf("<b>Total vault records: %d</b><br>", count))
		}
		docRecords.Call("append", moreDiv)
	}
	// if we do not have any records matched above we'll use first nrec ones
	if len(rids) == 0 {
		count := 0
		for key, _ := range recordsManager.Map {
			count += 1
			if count > nrec {
				break
			}
			rids = append(rids, key)
		}
	}
	return rids, nil
}

// RecordsHandler handles asynchronously HTTP requests
func RecordsHandler(url, cipher, passphrase, pattern, pageUrl string, args []js.Value) {
	resolve := args[0]
	reject := args[1]

	// place records within DOM page
	extention := true
	rids, err := updateRecords(url, cipher, passphrase, pattern, pageUrl, extention)
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
