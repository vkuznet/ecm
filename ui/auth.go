package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	fyne "fyne.io/fyne/v2"
)

// Auth interface declares how to access cloud providers
type Auth interface {
	OAuth()
	GetToken(code string) ([]byte, error)
}

var dropboxClient *Dropbox

// Dropbox client
type Dropbox struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	AuthURL      string
	TokenURL     string
}

func (d *Dropbox) OAuth() {
	//     ruri := url.QueryEscape("http://localhost:5151")
	//     auri := fmt.Sprintf("https://www.dropbox.com/oauth2/authorize?client_id=%s&response_type=token&redirect_uri=%s", d.ClientID, ruri)
	//     auri := fmt.Sprintf("https://www.dropbox.com/oauth2/authorize?client_id=%s&response_type=code", d.ClientID)
	//     auri := fmt.Sprintf("https://www.dropbox.com/oauth2/authorize?client_id=%s&response_type=code&redirect_uri=%s", d.ClientID, ruri)
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
	dropboxClient = &Dropbox{
		ClientID:     "v1yr4i1mk5bi8fw",
		ClientSecret: "m6w57h8tsandqx4",
		RedirectURI:  "http://localhost:5151",
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
				//                 var token DropboxToken
				//                 err = json.Unmarshal(data, &token)
			}
			msg := "Your ECM confiugration is updated with Dropbox credentials"
			appLog("INFO", msg, nil)
			msg += "<p>Please return to ECM app</p>"
			w.Write([]byte(msg))
		}
	})
	http.ListenAndServe(":5151", nil)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
	}
}
