package main

import (
	"bytes"
	_ "embed"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dchest/captcha"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgryski/dgoogauth"
	"github.com/gorilla/context"
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

	tmplData := make(TmplRecord)
	tmplData["User"] = user
	page := tmplPage("main.tmpl", tmplData)
	w.Write([]byte(page))
}

/*
 * 2fa handlers
 */

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
		// send POST request to /verify end-point to receive OTP token
		// JSON {"otp":"383878", "user": "UserName"}
		rec := make(map[string]string)
		rec["otp"] = otp
		rec["user"] = user
		data, err := json.Marshal(rec)
		if err != nil {
			log.Println("unable to marshal user data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		host := fmt.Sprintf("http://localhost:%d/verify", ServerConfig.Port)
		req, err := http.NewRequest("POST", host, bytes.NewBuffer(data))
		if err != nil {
			log.Println("unable to post request to /verify", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
		req.Header.Set("Content-Type", "application/json")
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("unable to post request to /verify", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("unable to read body from /verify", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// get data from verify with OTP token
		var otpToken string
		err = json.Unmarshal(body, &otpToken)
		if err != nil {
			log.Println("unable to unmarshal otpToken", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// post request to MainHandler with user data
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", otpToken))
		MainHandler(w, r)
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
		json.NewEncoder(w).Encode(err)
		return
	}
	otpc := &dgoogauth.OTPConfig{
		Secret:      secret,
		WindowSize:  3,
		HotpCounter: 1, // originally was 0
	}
	decodedToken["authorized"], err = otpc.Authenticate(otpToken.Token)
	if err != nil {
		log.Println("error", err)
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

		// generate a random string - preferbly 6 or 8 characters
		randomStr := randStr(6, "alphanum")

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
	authLink := fmt.Sprintf("otpauth://totp/GPM:%s?secret=%s&issuer=GPM", user, jwtSecret)

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
