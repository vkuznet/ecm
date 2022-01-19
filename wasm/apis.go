package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"syscall/js"
	"time"

	uuid "github.com/google/uuid"
)

// credentials provides vault credentials
func credentials() (string, string) {
	// TODO: implement proper logic to get it from HTML
	cipher := "aes"
	password := "test"
	return cipher, password
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
	document := js.Global().Get("document")
	records := document.Call("getElementById", "records")
	records.Set("innerHTML", "")
	doc := document.Call("getElementById", "vault")
	vault := doc.Get("value").String()
	url := fmt.Sprintf("/vault/%s/records", vault)
	cipher, password := credentials()
	extention := false
	rids, err := updateRecords(url, cipher, password, extention)
	if err != nil {
		return []byte{}, err
	}
	data, err := json.Marshal(rids)
	return data, err
}

func loginRecord() ([]byte, error) {
	var data []byte
	var err error

	document := js.Global().Get("document")
	records := document.Call("getElementById", "records")
	// get all inputs for login record
	name := document.Call("getElementById", "new-record-name").String()
	login := document.Call("getElementById", "new-record-login").String()
	password := document.Call("getElementById", "new-record-password").String()
	tags := document.Call("getElementById", "new-record-tags").String()
	url := document.Call("getElementById", "new-record-url").String()

	uid := uuid.NewString()
	rmap := make(Record)
	rmap["Name"] = name
	rmap["Login"] = login
	rmap["Password"] = password
	rmap["URL"] = url
	rmap["Tags"] = tags
	//     vrec := vault.VaultRecord{ID: uid, Map: rmap, ModificationTime: time.Now()}
	//     data, err := json.Marshal(rmap)
	//     err = vault.WriteRecord(vrec)
	data, err = json.Marshal(rmap)
	msg := fmt.Sprintf("New record created with UUID: %s", uid)
	records.Set("innerHTML", msg)

	return data, err
}
func cardRecord() ([]byte, error) {
	var data []byte
	var err error
	return data, err
}
func jsonRecord() ([]byte, error) {
	var data []byte
	var err error
	return data, err
}
func noteRecord() ([]byte, error) {
	var data []byte
	var err error
	return data, err
}
func uploadFile() ([]byte, error) {
	var data []byte
	var err error
	return data, err
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
	docRecords.Set("innerHTML", password)
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
