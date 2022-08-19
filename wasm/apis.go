package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"syscall/js"
	"time"

	uuid "github.com/google/uuid"
	crypt "github.com/vkuznet/ecm/crypt"
	vt "github.com/vkuznet/ecm/vault"
)

// credentials provides vault credentials
func credentials() (string, string, string) {
	document := js.Global().Get("document")
	vault := document.Call("getElementById", "vault-name").Get("value").String()
	cipher := document.Call("getElementById", "vault-cipher").Get("value").String()
	password := document.Call("getElementById", "vault-password").Get("value").String()
	return vault, cipher, password
}

// helper function to post data to our server in secure manner
// first we encrypt the data, and then send it over HTTP
func postData(api string, rec interface{}) error {
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	// get user credentials
	vault, cipher, password := credentials()

	data, err = crypt.Encrypt(data, password, cipher)
	if err != nil {
		return err
	}
	host := fmt.Sprintf("/vault/%s/%s", vault, api)
	body := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPost, host, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := httpClient(RootCA)
	_, err = client.Do(req)
	return err
}

func defaultAction() ([]byte, error) {
	document := js.Global().Get("document")
	records := document.Call("getElementById", "records")
	records.Set("innerHTML", "Not implemented")
	doc := document.Call("getElementById", "results")
	doc.Set("innerHTML", "Not implemented")
	result := make(map[string]string)
	result["Result"] = "show records"
	data, err := json.Marshal(result)
	return data, err
}
func showRecords() ([]byte, error) {
	var empty []byte
	document := js.Global().Get("document")
	records := document.Call("getElementById", "records")
	if records.IsNull() || records.IsUndefined() {
		return empty, nil
	}
	records.Set("innerHTML", "")
	vault, cipher, password := credentials()
	url := fmt.Sprintf("/vault/%s/records", vault)
	extention := false
	pattern := ""
	rids, err := updateRecords(url, cipher, password, pattern, extention)
	if err != nil {
		return []byte{}, err
	}
	data, err := json.Marshal(rids)
	return data, err
}

// helper function to create new vault record
func newRecord(rmap vt.Record) ([]byte, error) {
	var data []byte
	var err error

	document := js.Global().Get("document")
	records := document.Call("getElementById", "records")
	uid := uuid.NewString()

	vrec := vt.VaultRecord{ID: uid, Map: rmap, ModificationTime: time.Now()}
	err = postData("record", vrec)
	msg := fmt.Sprintf("New record created with UUID: %s", uid)
	if err != nil {
		msg = fmt.Sprintf("Failt to create new record %s, error %v", uid, err)
		msg = fmt.Sprintf("<div class=\"alert is-error is-shadow-2\">%s</div>", msg)
	}
	records.Set("innerHTML", msg)

	return data, err
}

func loginRecord() ([]byte, error) {
	rmap := make(vt.Record)
	document := js.Global().Get("document")
	rmap["Name"] = document.Call("getElementById", "new-record-name").Get("value").String()
	rmap["Login"] = document.Call("getElementById", "new-record-login").Get("value").String()
	rmap["Password"] = document.Call("getElementById", "new-record-password").Get("value").String()
	rmap["Tags"] = document.Call("getElementById", "new-record-tags").Get("value").String()
	rmap["Url"] = document.Call("getElementById", "new-record-url").Get("value").String()
	return newRecord(rmap)

}
func cardRecord() ([]byte, error) {
	rmap := make(vt.Record)
	document := js.Global().Get("document")
	rmap["Name"] = document.Call("getElementById", "new-card-name").Get("value").String()
	rmap["Number"] = document.Call("getElementById", "new-card-number").Get("value").String()
	rmap["Code"] = document.Call("getElementById", "new-card-code").Get("value").String()
	rmap["Tags"] = document.Call("getElementById", "new-card-tags").Get("value").String()
	rmap["Date"] = document.Call("getElementById", "new-card-date").Get("value").String()
	rmap["Phone"] = document.Call("getElementById", "new-card-phone").Get("value").String()
	return newRecord(rmap)
}
func jsonRecord() ([]byte, error) {
	rmap := make(vt.Record)
	// get all inputs for login record
	document := js.Global().Get("document")
	rmap["name"] = document.Call("getElementById", "new-json-name").Get("value").String()
	json := document.Call("getElementById", "new-json-record").Get("value").String()
	// TODO: I need to parse and decompose JSON into individual key-value pairs
	rmap["JSON"] = json
	return newRecord(rmap)
}
func noteRecord() ([]byte, error) {
	rmap := make(vt.Record)
	document := js.Global().Get("document")
	rmap["Name"] = document.Call("getElementById", "new-note-name").Get("value").String()
	rmap["Note"] = document.Call("getElementById", "new-note-record").Get("value").String()
	return newRecord(rmap)
}
func uploadFile(fname string, size int, ftype, content string) ([]byte, error) {
	rmap := make(vt.Record)
	rmap["Name"] = fname
	rmap["Size"] = fmt.Sprintf("%d", size)
	rmap["Type"] = ftype
	rmap["Tags"] = "file"
	rmap["Data"] = content
	return newRecord(rmap)
}

func syncHosts() ([]byte, error) {
	var data []byte
	var err error
	return data, err
}
func createVault() ([]byte, error) {
	var data []byte
	var err error
	return data, err
}

func newPassword() ([]byte, error) {
	const voc string = "abcdfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const numbers string = "0123456789"
	const symbols string = "!@#$%&*+_-="
	document := js.Global().Get("document")
	docRecords := document.Call("getElementById", "new-password")
	size := document.Call("getElementById", "password-size").Get("value").String()
	chars := document.Call("getElementById", "characters").Get("value").String()
	var password string
	if chars == "chars+numbers" {
		password = generatePassword(size, voc+numbers)
	} else if chars == "chars+numbers+symbols" {
		password = generatePassword(size, voc+numbers+symbols)
	} else {
		password = generatePassword(size, voc)
	}
	docRecords.Set("innerHTML", fmt.Sprintf("New password: <b>%s</b>", password))
	result := make(map[string]string)
	result["Result"] = "new password was generated"
	data, err := json.Marshal(result)
	return data, err
}

// helper function to generate password of certain length and chars
func generatePassword(size interface{}, chars string) string {
	rand.Seed(time.Now().UnixNano())
	length := 16
	switch v := size.(type) {
	case string:
		s, _ := strconv.Atoi(v)
		length = s
	case int, int32, int64:
		length = v.(int)
	}
	password := ""
	for i := 0; i < length; i++ {
		password += string([]rune(chars)[rand.Intn(len(chars))])
	}
	return password
}
