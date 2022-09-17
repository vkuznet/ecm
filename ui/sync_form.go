package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	binding "fyne.io/fyne/v2/data/binding"
	theme "fyne.io/fyne/v2/theme"
	widget "fyne.io/fyne/v2/widget"
	ecmsync "github.com/vkuznet/ecm/sync"
	"golang.org/x/exp/errors"
)

// supportedProviders contains a list of supported providers
var supportedProviders string

// global variable we'll use to update sync status
var syncStatus, localURI, cloudURI binding.String

// helper function to get list of supported providers
func supportedProvider(provider string) bool {
	var providers []string
	for _, p := range strings.Split(supportedProviders, ",") {
		providers = append(providers, strings.ToLower(strings.Trim(p, " ")))
	}
	if len(providers) == 0 {
		providers = append(providers, "dropbox")
	}
	for _, p := range providers {
		if p == provider {
			return true
		}
	}
	return false

}

// helper function to read sync config
func syncConfigMap(app fyne.App) (map[string]string, error) {
	out := make(map[string]string)
	// get our sync config file
	sconf := syncPath(app)

	// open our config file in rw mode
	file, err := os.Open(sconf)
	defer file.Close()
	if err != nil {
		return out, err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return out, err
	}
	var key, oldKey, section string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			oldKey = key
			if key != "" {
				out[key] = section
				section = ""
			}
			key = line[1 : len(line)-1]
		}
		if key != oldKey {
			if section == "" {
				section = line
			} else {
				section = fmt.Sprintf("%s\n%s\n", section, line)
			}
		}
	}
	return out, nil
}

// helper function to read and log sync config
func logSyncConfig(app fyne.App) error {
	// get our sync config file
	sconf := syncPath(app)

	// open our config file in rw mode
	file, err := os.Open(sconf)
	defer file.Close()
	if err != nil {
		return err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	appLog("INFO", "sync config", nil)
	appLog("INFO\n", string(data), nil)
	return nil
}

// helper function to update sync config file with token data
func updateSyncConfig(app fyne.App, provider, cid, secret string, data []byte) error {
	// get our sync config file
	sconf := syncPath(app)

	// open our config file in rw mode
	file, err := os.OpenFile(sconf, os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		return err
	}
	sdata, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	// process our data
	var out []string
	lines := strings.Split(string(sdata), "\n")
	match := fmt.Sprintf("[%s]", provider)
	var found bool
	for _, line := range lines {
		if line == match {
			out = append(out, line)
			// write our token data afterwards
			row := fmt.Sprintf("token = %s", string(data))
			out = append(out, row)
			found = true
			// write client id and secret
			row = fmt.Sprintf("client_id = %s", cid)
			out = append(out, row)
			row = fmt.Sprintf("client_secret = %s", secret)
			out = append(out, row)
		} else {
			if found && strings.HasPrefix(line, "token") {
				continue // skip previous token if it existed for our provider
			} else if found && strings.HasPrefix(line, "client") {
				continue // skip previous token if it existed for our provider
			} else {
				out = append(out, line)
			}
		}
	}

	// inject token info into appropriate section of config file
	content := strings.Join(out, "\n")
	_, err = file.WriteAt([]byte(content), 0)
	if err != nil {
		return err
	}
	appLog("INFO", "wrote new sync config", nil)
	appLog("INFO\n", content, nil)
	return nil
}

// helper function to provide sync path configuration and set RCLONE environment
func syncPath(app fyne.App) string {
	dir := app.Storage().RootURI().Path()
	spath := fmt.Sprintf("%s/rclone.conf", dir)
	sconf := os.Getenv("RCLONE_CONFIG")
	if sconf != "" {
		spath = sconf
	} else {
		err := os.Setenv("RCLONE_CONFIG", spath)
		if err != nil {
			msg := fmt.Sprintf("unable to setup RCLONE_CONFIG path %s", spath)
			appLog("ERROR", msg, err)
		}
	}
	return spath
}

// helper function to create sync config for rclone library
func syncConfig(app fyne.App) {

	// NOTE: write sync config for rclone usage
	// if we will not use rclone then I need to comment this out
	//     WriteSyncConfig(app)

	sconf := syncPath(app)
	if _, err := os.Stat(sconf); errors.Is(err, os.ErrNotExist) {
		msg := fmt.Sprintf("create %s", sconf)
		log.Println("INFO: ", msg)
		err := ecmsync.EcmCreateConfig(sconf)
		if err != nil {
			msg := fmt.Sprintf("untable to create %s", sconf)
			log.Println("ERROR: ", msg, " error: ", err)
		}
	}
}

// helper function to authenticate sync provider
func authSyncProvider(app fyne.App, provider string) {
	sconf := syncPath(app)
	msg := fmt.Sprintf("update %s token", provider)
	appLog("INFO", msg, nil)
	err := ecmsync.EcmUpdateConfig(sconf, provider)
	if err != nil {
		msg := fmt.Sprintf("untable to update %s", sconf)
		appLog("ERROR", msg, err)
	}
}

// SyncUI represents UI SyncUI
type SyncUI struct {
	preferences  fyne.Preferences
	window       fyne.Window
	app          fyne.App
	vaultRecords *vaultRecords
}

func newUISync(a fyne.App, w fyne.Window, v *vaultRecords) *SyncUI {
	return &SyncUI{
		app:          a,
		window:       w,
		vaultRecords: v,
		preferences:  a.Preferences(),
	}
}

func (r *SyncUI) onDropboxPathChanged(v string) {
	r.preferences.SetString("dropbox", v)
}

func (r *SyncUI) onCloudPathChanged(v string) {
	cloudURI.Set(v)
	r.preferences.SetString("cloud", v)
}

func (r *SyncUI) onLocalPathChanged(v string) {
	localURI.Set(v)
	r.preferences.SetString("local", v)
}

// helper function to provide sync button to given destination
func (r *SyncUI) authButton(provider string) *widget.Button {
	if provider == "noauth" {
		return &widget.Button{Text: ""}
	}
	btn := &widget.Button{
		Text: "OAuth",
		Icon: theme.ConfirmIcon(),
		OnTapped: func() {
			//             if appKind == "desktop" {
			//                 authSyncProvider(r.app, provider)
			//             } else {
			//                 dropboxClient.OAuth()
			//             }

			// check that new value for cloudURI is among supported providers
			val, err := cloudURI.Get()
			if err == nil {
				p := strings.Split(val, ":")[0]
				if !supportedProvider(p) {
					msg := fmt.Sprintf("provider %s is not supported in %v", p, supportedProviders)
					appLog("ERROR", msg, errors.New("user input error"))
					return
				}
				provider = p
			} else {
				appLog("ERROR", "fail to get cloudURI", err)
				return
			}
			msg := fmt.Sprintf("Supported provider %s", provider)
			appLog("INFO", msg, nil)

			// perform oauth request with our cloud provider
			if provider == "dropbox" {
				dropboxClient.OAuth()
			} else {
				msg := fmt.Sprintf("Cloud provider %s is not supported yet", provider)
				appLog("ERROR", msg, errors.New(msg))
				return
			}

			// update rclone path
			//             sconf := syncPath(r.app)
			//             err := ecmsync.EcmUpdateConfig(sconf, provider)
			//             if err != nil {
			//                 appLog("ERROR", "unable to update rclone config", err)
			//             }
			//             authSyncProvider(r.app, provider)
		},
	}
	return btn
}

// helper function to perform sync operation
func syncFunc(app fyne.App, vdir, src string, local bool) {
	dst := fmt.Sprintf("local:%s", vdir)
	// fetch local entry if it was set within widget
	if local {
		if val, err := localURI.Get(); err == nil {
			msg := fmt.Sprintf("syncFunc %s", src)
			appLog("INFO", msg, nil)
			src = val
		} else {
			msg := fmt.Sprintf("syncFunc %v", err)
			appLog("ERROR", msg, nil)
		}
		dst = vdir
	}

	// perform sync from dropbox to vault
	dir := app.Storage().RootURI().Path()
	fconf := fmt.Sprintf("%s/rclone.conf", dir)
	sconf := os.Getenv("ECM_SYNC_CONFIG")
	if sconf != "" {
		fconf = sconf
	}
	msg := fmt.Sprintf("config: %s, sync from %s to %s", fconf, src, dst)
	var err error
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		pref := app.Preferences()
		vdir := pref.String("VaultDirectory")
		arr := strings.Split(vdir, "/")
		vname := "Primary"
		if len(arr) > 1 {
			vname = arr[len(arr)-1]
		}
		rurl := fmt.Sprintf("%s/vault/%s/records?id=true", src, vname)
		appLog("INFO", msg, nil)
		err = ecmsync.SyncFromServer(rurl, dst)
	} else {
		appLog("INFO", msg, nil)
		err = ecmsync.EcmSync(fconf, src, dst)
	}
	if err != nil {
		msg := fmt.Sprintf("unable to sync from %s to %s", src, dst)
		appLog("ERROR", msg, err)
		syncStatus.Set(msg)
		return
	}
	log.Println("records are synced")
	// reset vault records
	_vault.Records = nil
	// read again vault records
	err = _vault.Read()
	if err != nil {
		msg := fmt.Sprintf("unable to read the vault records, %v", err)
		appLog("ERROR", msg, err)
		syncStatus.Set(msg)
		return
	}
	// refresh ui records
	//     r.vaultRecords.Refresh()
	if appRecords != nil {
		appRecords.Refresh()
	}
	msg = fmt.Sprintf("%s records are synced successfully", src)
	syncStatus.Set(msg)
	appLog("INFO", msg, nil)
}

// helper function to provide sync button to given destination
func (r *SyncUI) syncButton(src string, local bool) *widget.Button {
	// get vault dir from preferences
	pref := r.app.Preferences()
	vdir := pref.String("VaultDirectory")

	btn := &widget.Button{
		Text: "Sync",
		Icon: theme.HistoryIcon(),
		OnTapped: func() {
			syncFunc(r.app, vdir, src, local)
		},
	}
	return btn
}

// Token describes structure of OAuath token response
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expires      int64  `json:"expires_in"`
}

// helper function to check if token is valid for given cloud provider
func (r *SyncUI) isValidToken(provider string) bool {
	// TODO: we should generalize this function to any cloud provider
	if provider != "dropbox" {
		msg := fmt.Sprintf("Cloud provider %s is not supported", provider)
		appLog("ERROR", msg, errors.New(msg))
		return false
	}

	// get our config mod time in unix format
	sconf := syncPath(r.app)
	file, err := os.Stat(sconf)
	if err != nil {
		log.Println("sync path", err)
		appLog("ERROR", "unable to get file stats", err)
		return false
	}
	mtime := file.ModTime().Unix()
	//     log.Println("sconf", sconf, "mod time", mtime)

	// read out sync config map
	sdict, err := syncConfigMap(r.app)
	if err != nil {
		log.Println("unable to read sync config map", err)
		appLog("ERROR", "unable to read sync config map", err)
		return false
	}
	// check if our provider has token
	if vals, ok := sdict[provider]; ok {
		if strings.Contains(vals, "token") {
			for _, line := range strings.Split(vals, "\n") {
				if strings.HasPrefix(line, "token") {
					arr := strings.Split(line, " = ")
					if len(arr) == 2 {
						var token Token
						err := json.Unmarshal([]byte(arr[1]), &token)
						if err != nil {
							appLog("ERROR", "unable to unmarshal token", err)
							return false
						}
						// check if our token is valid
						if mtime+token.Expires > time.Now().Unix() {
							return true
						} else {
							// place refresh token request
							data, err := dropboxClient.RefreshToken(token.RefreshToken)
							if err == nil {
								// if successfull we get new token object
								var newToken Token
								err := json.Unmarshal(data, &newToken)
								if err == nil {
									// if successfull we update new token object
									newToken.RefreshToken = token.RefreshToken
									// and we get its binary representation
									newData, err := json.Marshal(data)
									if err == nil {
										// if successfull we update sync config
										updateSyncConfig(
											r.app,
											"dropbox",
											dropboxClient.ClientID,
											dropboxClient.ClientSecret,
											newData,
										)
									}
									return true

								}
							}
						}
					}
				}
			}
		}
	}
	return false
}

// global auth button pointer
var cloudAuthButton *widget.Button
var localSetButton *widget.Button

// helper function to build UI
func (r *SyncUI) buildUI() *fyne.Container {

	// create text box which will update text once sync is completed
	syncStatus = binding.NewString()
	syncStatus.Set("Sync status will appear here")
	statusText := widget.NewLabelWithData(syncStatus)
	statusText.Wrapping = fyne.TextWrapBreak

	// cloud storage URI
	cloudURI = binding.NewString()
	//     cloudURI.Set("dropbox:ECM")
	cloudURI.Set(getPrefValue(r.preferences, "cloud", "dropbox:ECM"))
	cloud := widget.NewEntryWithData(cloudURI)
	cloud.OnChanged = r.onCloudPathChanged
	//     dropbox := &widget.Entry{Text: "dropbox:ECM", OnSubmitted: r.onDropboxPathChanged}
	cloudStorage := widget.NewEntryWithData(cloudURI)
	cloudStorage.OnSubmitted = r.onCloudPathChanged
	cloudStorage.OnChanged = r.onCloudPathChanged

	// local storage URI
	lpath := os.Getenv("EXTERNAL_STORAGE")
	lpath = fmt.Sprintf("local:/%s/ECM", lpath)
	if appKind == "desktop" {
		home := os.Getenv("HOME")
		lpath = fmt.Sprintf("local:%s/.ecm", home)
	}
	localURI = binding.NewString()
	localValue, err := localURI.Get()
	// set localURI only we need to
	if err != nil || localValue == "" {
		// set localURI either from preference value or provided one
		localURI.Set(getPrefValue(r.preferences, "local", lpath))
		//         localURI.Set(lpath)
	}
	//local := &widget.Entry{Text: lpath, OnSubmitted: r.onLocalPathChanged}
	local := widget.NewEntryWithData(localURI)
	local.OnChanged = r.onLocalPathChanged

	cloudAuthButton = r.authButton("dropbox")
	cloudSync := colorButtonContainer(r.syncButton(cloud.Text, false), authColor)
	cloudAuth := colorButtonContainer(cloudAuthButton, authColor)
	localSync := colorButtonContainer(r.syncButton(local.Text, true), btnColor)
	//     btn := &widget.Button{}
	//     noAuth := colorButtonContainer(btn, grayColor)

	cloudLabel := widget.NewLabel("Cloud to vault")
	cloudLabel.TextStyle.Bold = true
	labelName := "local to vault"
	localLabel := widget.NewLabel(labelName)
	localLabel.TextStyle.Bold = true

	// by default we show auth button
	// check if token exist and it is valid, then we show sync button
	cloudContainer := container.NewGridWithColumns(2, cloud, cloudAuth)
	if r.isValidToken("dropbox") {
		cloudContainer = container.NewGridWithColumns(2, cloud, cloudSync)
	}
	localContainer := container.NewGridWithColumns(2, local, localSync)

	box := container.NewVBox(
		cloudLabel,
		cloudContainer,
		localLabel,
		localContainer,
		statusText,
	)
	return box
}
func (r *SyncUI) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Sync", Icon: theme.HistoryIcon(), Content: r.buildUI()}
}
