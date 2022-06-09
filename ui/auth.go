package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	fyne "fyne.io/fyne/v2"
	"github.com/joho/godotenv"
)

// global var for dropbox object
var dropboxClient *Dropbox

// Auth interface declares how to access cloud providers
type Auth interface {
	OAuth()
	GetToken(code string) ([]byte, error)
}

// Dropbox structure represent Dropbox auth object, for more information see
// https://developers.dropbox.com/oauth-guide
// https://www.dropbox.com/developers/documentation/http/documentation#oauth2-token
type Dropbox struct {
	ClientID     string // dropbox client id
	ClientSecret string // dropbox client secret
	Port         string // redirect port
	RedirectURI  string // redirect URI
	AuthURL      string // dropbox authentication URL
	TokenURL     string // dropbox token URL
}

// OAuth implements Auth.OAuth method for Dropbox
func (d *Dropbox) OAuth() {
	rurl := url.QueryEscape(d.RedirectURI)
	auri := fmt.Sprintf(
		"%s?client_id=%s&response_type=code&redirect_uri=%s",
		d.AuthURL,
		d.ClientID,
		rurl,
	)
	fmt.Println("auth url", auri)
	openURL(auri)
}

// GetToken implements Auth.GetToken method for Dropbox
func (d *Dropbox) GetToken(code string) ([]byte, error) {
	vals := url.Values{}
	vals.Set("code", code)
	vals.Set("grant_type", "authorization_code")
	vals.Set("redirect_uri", d.RedirectURI)

	log.Printf("values %+v", vals)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", d.TokenURL, strings.NewReader(vals.Encode()))
	if err != nil {
		return []byte{}, err
	}
	req.SetBasicAuth(d.ClientID, d.ClientSecret)
	resp, err := client.Do(req)

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	return data, err

}

// DropboxToken represents structure of dropbox token response
type DropboxToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     int64  `json:"expires_in"`
	Scope       string `json:"scope"`
	UID         string `json:"uid"`
	AccountID   string `json:"account_id"`
}

// helper function to perform dropbox authentication
func authDropbox() {
	err := godotenv.Load("credentials.env")
	if err != nil {
		appLog("INFO", "uunable to load credentials.env file", err)
		return
	}
	cid := os.Getenv("DROPBOX_CLIENT_ID")
	sec := os.Getenv("DROPBOX_CLIENT_SECRET")
	port := os.Getenv("DROPBOX_PORT")
	dropboxClient = &Dropbox{
		ClientID:     cid,
		ClientSecret: sec,
		Port:         port,
		RedirectURI:  fmt.Sprintf("http://localhost:%s", port),
		TokenURL:     "https://api.dropbox.com/oauth2/token",
		AuthURL:      "https://www.dropbox.com/oauth2/authorize",
	}
	dropboxClient.OAuth()
}

// authServer provides internal web server which handles access token HTTP requests
func authServer(app fyne.App, ctx context.Context) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")
		if dropboxClient != nil && code != "" {
			data, err := dropboxClient.GetToken(code)
			appLog("INFO", string(data), err)
			if err == nil {
				updateSyncConfig(app, "dropbox", data)
			}
			msg := "Your ECM confiugration is updated with Dropbox credentials"
			appLog("INFO", msg, nil)
			msg += "Please restart the ECM app to proceed"
			htmlMsg := fmt.Sprintf("<html><body><h3>%s</h3></body></html>", msg)
			w.Write([]byte(htmlMsg))
		}
	})
	http.ListenAndServe(fmt.Sprintf(":%s", dropboxClient.Port), nil)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
	}
}
