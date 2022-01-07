package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
)

// we embed few html pages directly into server
// but for advanced usage users should switch to templates

//go:embed "static/tmpl/top.tmpl"
var topHTML string

//go:embed "static/tmpl/bottom.tmpl"
var bottomHTML string

// responseMsg helper function to provide response to end-user
func responseMsg(w http.ResponseWriter, r *http.Request, msg, api string, code int) int64 {
	rec := make(Record)
	rec["error"] = msg
	rec["api"] = api
	rec["method"] = r.Method
	rec["exception"] = fmt.Sprintf("%d", code)
	rec["type"] = "HTTPError"
	//     data, _ := json.Marshal(rec)
	var out []Record
	out = append(out, rec)
	data, _ := json.Marshal(out)
	w.WriteHeader(code)
	w.Write(data)
	return int64(len(data))
}

// helper function to get vault parameters from the HTTP request
func getVault(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	var vault string
	vault, ok := vars["vault"]
	if !ok {
		return "", errors.New("there is no vault parameter in HTTP request")
	}
	vdir := filepath.Join(ServerConfig.VaultArea, vault)
	_, err := os.Stat(vdir)
	if err != nil {
		msg := fmt.Sprintf("unable to access vault %s, error %v", vdir, err)
		return "", errors.New(msg)
	}
	return vdir, nil
}

// helper function to get vault parameters from the HTTP request
func getVaultRecord(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	rid, ok := vars["rid"]
	if !ok {
		return "", errors.New("there is no rid parameter in HTTP request")
	}
	return rid, nil
}

// VaultHandler provides basic functionality of status response
func VaultHandler(w http.ResponseWriter, r *http.Request) {
	vdir, err := getVault(r)
	if err != nil {
		responseMsg(w, r, fmt.Sprintf("%v", err), "VaultHandler", http.StatusBadRequest)
		return
	}
	vault := Vault{Cipher: getCipher(""), Secret: "", Directory: vdir}
	files, err := vault.Files()
	if err != nil {
		responseMsg(w, r, err.Error(), "VaultHandler", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(files)
	if err != nil {
		responseMsg(w, r, err.Error(), "VaultHandler", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// fill them out at VaultAuthHandler
type AuthRecord struct {
	Cipher string
	Secret string
}

var auth AuthRecord

// VaultAuthHandler provides authentication with our vault
func VaultAuthHandler(w http.ResponseWriter, r *http.Request) {
	// it should be POST request which will ready vault credentials
	if r.Method == "POST" {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&auth)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// VaultRecordsHandler provides basic functionality of status response
func VaultRecordsHandler(w http.ResponseWriter, r *http.Request) {
	vdir, err := getVault(r)
	if err != nil {
		responseMsg(w, r, fmt.Sprintf("%v", err), "VaultHandler", http.StatusBadRequest)
		return
	}
	vault := Vault{Directory: vdir}
	files, err := vault.Files()
	if err != nil {
		responseMsg(w, r, err.Error(), "VaultHandler", http.StatusInternalServerError)
		return
	}
	var records [][]byte
	for _, name := range files {
		fname := fmt.Sprintf("%s/%s", vdir, name)
		data, err := os.ReadFile(fname)
		if err != nil {
			responseMsg(w, r, err.Error(), "VaultHandler", http.StatusInternalServerError)
			return
		}
		records = append(records, data)
	}
	data, err := json.Marshal(records)
	if err != nil {
		responseMsg(w, r, err.Error(), "VaultHandler", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// VaultRecordHandler provides basic functionality of status response
func VaultRecordHandler(w http.ResponseWriter, r *http.Request) {
	vdir, err := getVault(r)
	if err != nil {
		responseMsg(w, r, fmt.Sprintf("%v", err), "VaultRecordHandler", http.StatusBadRequest)
		return
	}
	log.Println("vault", vdir)
	rid, err := getVaultRecord(r)
	if err != nil {
		responseMsg(w, r, fmt.Sprintf("%v", err), "VaultRecordHandler", http.StatusBadRequest)
		return
	}
	fname := filepath.Join(vdir, rid)
	_, err = os.Stat(fname)
	if err != nil {
		responseMsg(w, r, fmt.Sprintf("%v", err), "VaultRecordHandler", http.StatusBadRequest)
		return
	}
	data, err := os.ReadFile(fname)
	if err != nil {
		responseMsg(w, r, fmt.Sprintf("%v", err), "VaultRecordHandler", http.StatusBadRequest)
		return
	}
	w.Write(data)
}

// VaultAddHandler provides basic functionality of status response
func VaultAddHandler(w http.ResponseWriter, r *http.Request) {
	vdir, err := getVault(r)
	if err != nil {
		responseMsg(w, r, fmt.Sprintf("%v", err), "VaultAddHandler", http.StatusBadRequest)
		return
	}
	log.Println("vault", vdir)
	w.WriteHeader(http.StatusOK)
}

// VaultDeleteHandler provides basic functionality of status response
func VaultDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vdir, err := getVault(r)
	if err != nil {
		responseMsg(w, r, fmt.Sprintf("%v", err), "VaultDeleteHandler", http.StatusBadRequest)
		return
	}
	rid, err := getVaultRecord(r)
	if err != nil {
		responseMsg(w, r, fmt.Sprintf("%v", err), "VaultDeleteHandler", http.StatusBadRequest)
		return
	}
	vault := Vault{Cipher: getCipher(""), Secret: "", Directory: vdir}
	err = vault.DeleteRecord(rid)
	if err != nil {
		responseMsg(w, r, fmt.Sprintf("%v", err), "VaultDeleteHandler", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// TokenHandler provides basic functionality of status response
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	passphrase := Config.TokenSecret
	cipher := "aes"
	token, err := encryptToken(passphrase, cipher)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(token))
}

// FaviconHandler provides favicon icon file
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	//     http.ServeFile(w, r, "relative/path/to/favicon.ico")
	w.WriteHeader(http.StatusOK)
}

// helper function to parse given template and return HTML page
func tmplPage(tmpl string, tmplData TmplRecord) string {
	if tmplData == nil {
		tmplData = make(TmplRecord)
	}
	var templates Templates
	page := templates.Tmpl(ServerConfig.Templates, tmpl, tmplData)
	return topHTML + page + bottomHTML
}

// HomeHandler handles home page requests
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmplData := make(TmplRecord)
	captchaStr := captcha.New()
	tmplData["CaptchaId"] = captchaStr
	page := tmplPage("index.tmpl", tmplData)
	w.Write([]byte(page))
}

// SignUpHandler handles sign-up page requests
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	tmplData := make(TmplRecord)
	captchaStr := captcha.New()
	tmplData["CaptchaId"] = captchaStr
	page := tmplPage("signup.tmpl", tmplData)
	w.Write([]byte(page))
}
