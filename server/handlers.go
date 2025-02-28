package main

import (
	"bytes"
	_ "embed"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dchest/captcha"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgryski/dgoogauth"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/vkuznet/ecm/crypt"
	utils "github.com/vkuznet/ecm/utils"
	vt "github.com/vkuznet/ecm/vault"
)

// we embed few html pages directly into server
// but for advanced usage users should switch to templates

//go:embed "static/tmpl/top.tmpl"
var topHTML string

//go:embed "static/tmpl/bottom.tmpl"
var bottomHTML string

// responseMsg helper function to provide response to end-user
func responseMsg(w http.ResponseWriter, r *http.Request, msg, api string, code int) int64 {
	rec := make(vt.Record)
	rec["error"] = msg
	rec["api"] = api
	rec["method"] = r.Method
	rec["exception"] = fmt.Sprintf("%d", code)
	rec["type"] = "HTTPError"
	//     data, _ := json.Marshal(rec)
	var out []vt.Record
	out = append(out, rec)
	data, _ := json.Marshal(out)
	w.WriteHeader(code)
	w.Write(data)
	return int64(len(data))
}

// helper function to generate error page with given message
func errorPage(w http.ResponseWriter, r *http.Request, msg string) {
	log.Println("ERROR:", msg)
	tmplData := make(TmplRecord)
	tmplData["Message"] = msg
	page := tmplPage("error.tmpl", tmplData)
	w.Write([]byte(page))
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
	vault := vt.Vault{Cipher: crypt.GetCipher(""), Secret: "", Directory: vdir}
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

// AuthRecord keeps vault auth attributes
type AuthRecord struct {
	Cipher string
	Secret string
}

var auth AuthRecord

// HTTPVaultRecord represents HTTP vault record
type HTTPVaultRecord struct {
	ID   string `json:"id"`
	Data []byte `json:"data"`
}

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
	// parse input parameters to identify if we need to construct id records
	var idRecord bool
	vals, ok := r.URL.Query()["id"]
	if ok && len(vals) > 0 {
		if vals[0] == "true" {
			idRecord = true
		}
	}

	vdir, err := getVault(r)
	if err != nil {
		responseMsg(w, r, fmt.Sprintf("%v", err), "VaultHandler", http.StatusBadRequest)
		return
	}
	vault := vt.Vault{Directory: vdir}
	files, err := vault.Files()
	if err != nil {
		responseMsg(w, r, err.Error(), "VaultHandler", http.StatusInternalServerError)
		return
	}
	var ids []string
	var records [][]byte
	for _, name := range files {
		fname := fmt.Sprintf("%s/%s", vdir, name)
		ids = append(ids, name)
		data, err := os.ReadFile(fname)
		if err != nil {
			responseMsg(w, r, err.Error(), "VaultHandler", http.StatusInternalServerError)
			return
		}
		records = append(records, data)
	}
	var data []byte
	if idRecord {
		var httpRecords []HTTPVaultRecord
		for idx, fid := range ids {
			rec := HTTPVaultRecord{ID: fid, Data: records[idx]}
			httpRecords = append(httpRecords, rec)
		}
		data, err = json.Marshal(httpRecords)
	} else {
		data, err = json.Marshal(records)
	}

	if err != nil {
		responseMsg(w, r, err.Error(), "VaultHandler", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// VaultRecordHandler provides basic functionality of status response
func VaultRecordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("unable to read body", err)
		}
		password := "test"
		cipher := "aes"
		data, err := crypt.Decrypt(body, password, cipher)
		if err != nil {
			log.Println("unable to decrypt", err)
		}
		var rec vt.VaultRecord
		err = json.Unmarshal(data, &rec)
		if err != nil {
			log.Println("unable to unmarshal", err, "received data:", string(data))
		}
		tags, _ := rec.Map["Tags"]
		if tags == "file" {
			// TMP: write file
			name, _ := rec.Map["Name"]
			fname := fmt.Sprintf("/tmp/%s", name)
			file, err := os.Create(fname)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			data, _ := rec.Map["Data"] // is a hex string
			raw, err := hex.DecodeString(data)
			if err != nil {
				log.Fatal(err)
			}
			err = os.WriteFile(fname, []byte(raw), 0755)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println("received", rec)
		}
		return
	}
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
	vault := vt.Vault{Cipher: crypt.GetCipher(""), Secret: "", Directory: vdir}
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
	page := tmplPage("index.tmpl", nil)
	w.Write([]byte(page))
}

// LoginHandler handles login page requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	page := tmplPage("login.tmpl", nil)
	w.Write([]byte(page))
}

// LogoutHandler handles login page requests
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	page := tmplPage("logout.tmpl", nil)
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

// MainHandler handles login page requests
func MainHandler(w http.ResponseWriter, r *http.Request) {

	// we should be redirected to this handler from login page
	var user string
	err := r.ParseForm()
	if err == nil {
		user = r.FormValue("user")
	} else {
		log.Println("unable to parse user form data", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	agent := strings.ToLower(r.Header.Get("User-Agent"))
	mobile := false
	if strings.Contains(agent, "ipad") ||
		strings.Contains(agent, "iphone") ||
		strings.Contains(agent, "android") {
		mobile = true
	}

	tmplData := make(TmplRecord)
	tmplData["User"] = user
	tmplData["Mobile"] = mobile
	page := tmplPage("main.tmpl", tmplData)
	w.Write([]byte(page))
}

/*
 * 2fa handlers
 */

// helper function to authenticate web request to MainHandler
func authMainHandler(w http.ResponseWriter, r *http.Request, otp, user, tokenString, secret string) {
	// send POST request to /verify end-point to receive OTP token
	// JSON {"otp":"383878", "user": "UserName"}
	rec := make(map[string]string)
	rec["otp"] = otp
	rec["user"] = user
	data, err := json.Marshal(rec)
	if err != nil {
		msg := fmt.Sprintf("unable to marshal user data, error %v", err)
		errorPage(w, r, msg)
		return
	}
	url := fmt.Sprintf("%s/verify", Localhost)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		msg := fmt.Sprintf("unable to post verification request, error %v", err)
		errorPage(w, r, msg)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("unable to request verification, error %v", err)
		errorPage(w, r, msg)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("unable to read body, error %v", err)
		errorPage(w, r, msg)
		return
	}

	// get data from verify with OTP token
	var otpToken string
	err = json.Unmarshal(body, &otpToken)
	if err != nil {
		log.Println("otpToken body", string(body))
		msg := fmt.Sprintf("unable to unmarshal otp token, error %v", err)
		errorPage(w, r, msg)
		return
	}
	decodedToken, err := VerifyJwt(otpToken, secret)
	if err != nil {
		msg := fmt.Sprintf("unable to verify otp token, error %v", err)
		errorPage(w, r, msg)
		return
	}
	if decodedToken["authorized"] == true {
		context.Set(r, "decoded", decodedToken)
		// post request to MainHandler with user data and our token
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", otpToken))
		MainHandler(w, r)
		return
	}
	msg := fmt.Sprintf("2fa verification process fails")
	errorPage(w, r, msg)
}

// AuthHandler authenticate user via POST HTTP request
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// parse form parameters
	var user, password, otp string
	if err := r.ParseForm(); err == nil {
		user = r.FormValue("user")
		password = r.FormValue("password")
		otp = r.FormValue("otp")
	}

	// check if our user exist in DBStore
	if !userExist(user, password) {
		msg := "Wrong password or user does not exist"
		if r.Header.Get("Content-Type") == "application/json" {
			rec := make(map[string]string)
			rec["error"] = msg
			json.NewEncoder(w).Encode(rec)
			return
		}
		tmplData := make(TmplRecord)
		tmplData["Message"] = msg
		page := tmplPage("error.tmpl", tmplData)
		w.Write([]byte(page))
		return
	}

	secret := findUserSecret(user)
	if secret == "" {
		err := errors.New("Non existing user, please use /qr end-point to initialize and get QR code")
		json.NewEncoder(w).Encode(err)
		return
	}
	mapClaims := jwt.MapClaims{}
	mapClaims["username"] = user
	mapClaims["authorized"] = false
	mapClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	// for verification we can use either user's secret
	// or server secret
	// in latter case it should be global and available to all APIs
	tokenString, err := SignJwt(mapClaims, secret)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	if otp != "" {
		authMainHandler(w, r, otp, user, tokenString, secret)
		return
	}

	// for non web clients return JSON document with JWT token
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}

// ApiHandler represents protected end-point for our server API
// It can be only reached via 2FA method
func ApiHandler(w http.ResponseWriter, r *http.Request) {
	// so far we return content of our HTTP request context
	// but its logic can implement anything
	decoded := context.Get(r, "decoded")
	json.NewEncoder(w).Encode(decoded)
}

// VerifyHandler authorizes user based on provided token and OTP code
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// this API expect JSON OtpToken payload
	var otpToken OtpToken
	err := json.NewDecoder(r.Body).Decode(&otpToken)
	if err != nil {
		log.Println("error", err)
		json.NewEncoder(w).Encode(err)
		return
	}
	if ServerConfig.Verbose > 0 {
		log.Println("otp token", otpToken)
	}
	user := otpToken.User
	secret := findUserSecret(user)
	if secret == "" {
		err := errors.New("Non existing user, please use /qr end-point to initialize and get QR code")
		json.NewEncoder(w).Encode(err)
		return
	}

	bearerToken, err := getBearerToken(r.Header.Get("authorization"))
	if ServerConfig.Verbose > 0 {
		log.Println("verify otp", bearerToken, "error", err)
	}
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	decodedToken, err := VerifyJwt(bearerToken, secret)
	if err != nil {
		log.Println("unable to verify jwt token", err)
		json.NewEncoder(w).Encode(err)
		return
	}
	otpc := &dgoogauth.OTPConfig{
		Secret:      secret,
		WindowSize:  3,
		HotpCounter: 0,
	}
	decodedToken["authorized"], err = otpc.Authenticate(otpToken.Token)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	if decodedToken["authorized"] == false {
		json.NewEncoder(w).Encode("Invalid one-time password")
		return
	}
	if ServerConfig.Verbose > 0 {
		log.Println("otp authorized", otpToken)
	}
	// for verification we can use either user's secret
	// or server secret
	// in latter case it should be global and available to all APIs
	jwToken, _ := SignJwt(decodedToken, secret)
	json.NewEncoder(w).Encode(jwToken)
}

// UserHandler handles sign-up HTTP requests
func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// parse form parameters
	var user, email, password string
	err := r.ParseForm()
	if err == nil {
		user = r.FormValue("user")
		email = r.FormValue("email")
		password = r.FormValue("password")
	} else {
		log.Println("unable to parse user form data", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	passwordHash, err := getPasswordHash(password)
	if err != nil {
		log.Println("unable to get password hash", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check if user provide the captcha
	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		tmplData := make(TmplRecord)
		tmplData["Message"] = "Wrong captcha match, robots are not allowed"
		page := tmplPage("error.tmpl", tmplData)
		w.Write([]byte(page))
		return
	}

	// check if user exists, otherwise create new user entry in DB
	if !userExist(user, password) {
		userRecord := User{
			Name:     user,
			Password: passwordHash,
			Email:    email,
			Secret:   "",
		}
		addUser(userRecord)
	}

	// redirect request to qrcode end-point
	if ServerConfig.Verbose > 0 {
		log.Printf("redirect %+v", r)
	}
	// to preserve the same HTTP method we should use
	// 307 StatusTemporaryRedirect code
	// https://softwareengineering.stackexchange.com/questions/99894/why-doesnt-http-have-post-redirect
	http.Redirect(w, r, "/qrcode", http.StatusTemporaryRedirect)
}

// QRHandler represents handler for /qr end-point to present our QR code
// to the client
func QRHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("call QRHandler %+v", r)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// parse form parameters
	var user string
	err := r.ParseForm()
	if err == nil {
		user = r.FormValue("user")
	} else {
		log.Println("unable to parse form data", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check if our user exists in DB
	if !userExist(user, "do not check") {
		msg := fmt.Sprintf("Unknown user %s", user)
		w.Write([]byte(msg))
		return
	}

	// proceed and either create or retrieve QR code for our user
	udir := fmt.Sprintf("static/data/%s", user)
	qrImgFile := fmt.Sprintf("%s/QRImage.png", udir)
	err = os.MkdirAll(udir, 0755)
	if err != nil {
		if ServerConfig.Verbose > 0 {
			log.Printf("unable to create directory %s, error %v", udir, err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userSecret := findUserSecret(user)
	if userSecret == "" {
		// no user exists in DB
		if ServerConfig.Verbose > 0 {
			log.Println("generate user secret")
		}

		// generate a random string: 10 characters will produce non = signs in base32 encoding
		randomStr := utils.RandomString(10, "alphanum")

		// For Google Authenticator purpose
		// for more details see
		// https://github.com/google/google-authenticator/wiki/Key-Uri-Format
		secret := base32.StdEncoding.EncodeToString([]byte(randomStr))
		jwtSecret = secret

		// update user secret in DB
		updateUser(user, jwtSecret)
	} else {
		if ServerConfig.Verbose > 0 {
			log.Println("read user secret from DB")
		}
		jwtSecret = userSecret
	}

	// authentication link.
	// for more details see
	// https://github.com/google/google-authenticator/wiki/Key-Uri-Format
	authLink := fmt.Sprintf("otpauth://totp/ECM:%s?secret=%s&issuer=ECM", user, jwtSecret)

	// generate QR image
	// Remember to clean up the file afterwards
	//     defer os.Remove(qrImgFile)
	err = generateQRImage(authLink, qrImgFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generate page content
	tmplData := make(TmplRecord)
	tmplData["User"] = user
	tmplData["ImageFile"] = qrImgFile
	page := tmplPage("qrcode.tmpl", tmplData)
	w.Write([]byte(page))
}
