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

// global variable we'll use to update sync status
var syncStatus binding.String

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

func (r *SyncUI) onLocalPathChanged(v string) {
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

			// perform oauth request with our dropbox client
			dropboxClient.OAuth()

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

// helper function to provide sync button to given destination
func (r *SyncUI) syncButton(src string) *widget.Button {
	// get vault dir from preferences
	pref := r.app.Preferences()
	vdir := pref.String("VaultDirectory")

	btn := &widget.Button{
		Text: "Sync",
		Icon: theme.HistoryIcon(),
		OnTapped: func() {
			// perform sync from dropbox to vault
			dir := r.app.Storage().RootURI().Path()
			fconf := fmt.Sprintf("%s/rclone.conf", dir)
			sconf := os.Getenv("ECM_SYNC_CONFIG")
			if sconf != "" {
				fconf = sconf
			}
			dst := fmt.Sprintf("local:%s", vdir)
			msg := fmt.Sprintf("config: %s, sync from %s to %s", fconf, src, dst)
			appLog("INFO", msg, nil)
			err := ecmsync.EcmSync(fconf, src, dst)
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
			r.vaultRecords.Refresh()
			msg = fmt.Sprintf("%s records are synced successfully", src)
			syncStatus.Set(msg)
			appLog("INFO", msg, nil)
		},
	}
	return btn
}

type Token struct {
	AccessToken string `json:"access_token"`
	Expire      int64  `json:"expires_in"`
}

// helper function to check if token is valid for given cloud provider
func (r *SyncUI) isValidToken(provider string) bool {
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
					//                     log.Println("found token", vals, arr, len(arr))
					if len(arr) == 2 {
						var token Token
						err := json.Unmarshal([]byte(arr[1]), &token)
						if err != nil {
							appLog("ERROR", "unable to unmarshal token", err)
							return false
						}
						// check if our token is valid
						//                         log.Println("token tstamp", mtime+token.Expire, " now ", time.Now().Unix())
						if mtime+token.Expire > time.Now().Unix() {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// helper function to build UI
func (r *SyncUI) buildUI() *fyne.Container {

	// create text box which will update text once sync is completed
	syncStatus = binding.NewString()
	syncStatus.Set("Sync status will appear here")
	statusText := widget.NewLabelWithData(syncStatus)
	statusText.Wrapping = fyne.TextWrapBreak

	// sync form container
	dropbox := &widget.Entry{Text: "dropbox:ECM", OnSubmitted: r.onDropboxPathChanged}
	lpath := os.Getenv("EXTERNAL_STORAGE")
	if lpath == "" {
		lpath = "sdcard"
	}
	lpath = fmt.Sprintf("local:/%s/ECM", lpath)
	if appKind == "desktop" {
		home := os.Getenv("HOME")
		lpath = fmt.Sprintf("local:%s/.ecm", home)
	}
	local := &widget.Entry{Text: lpath, OnSubmitted: r.onLocalPathChanged}

	dropboxSync := colorButtonContainer(r.syncButton(dropbox.Text), btnColor)
	dropboxAuth := colorButtonContainer(r.authButton("dropbox"), authColor)
	localSync := colorButtonContainer(r.syncButton(local.Text), btnColor)
	//     btn := &widget.Button{}
	//     noAuth := colorButtonContainer(btn, grayColor)

	dropboxLabel := widget.NewLabel("Dropbox to vault")
	dropboxLabel.TextStyle.Bold = true
	labelName := "local to vault"
	if appKind != "desktop" {
		labelName = "sdcard to vault"
	}
	localLabel := widget.NewLabel(labelName)
	localLabel.TextStyle.Bold = true

	// by default we show auth button
	dropboxContainer := container.NewGridWithColumns(2, dropbox, dropboxAuth)
	// check if token exist and it is valid, then we show sync button
	if r.isValidToken("dropbox") {
		dropboxContainer = container.NewGridWithColumns(2, dropbox, dropboxSync)
	}
	localContainer := container.NewGridWithColumns(2, local, localSync)

	box := container.NewVBox(
		dropboxLabel,
		dropboxContainer,
		localLabel,
		localContainer,
		statusText,
	)
	return box
}
func (r *SyncUI) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Sync", Icon: theme.HistoryIcon(), Content: r.buildUI()}
}
