package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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
func readSyncConfig(app fyne.App) error {
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
func updateSyncConfig(app fyne.App, provider string, data []byte) error {
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
	for _, line := range lines {
		if line == match {
			out = append(out, line)
			// write our token data afterwards
			out = append(out, string(data))
		} else {
			out = append(out, line)
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
func (r *SyncUI) onPCloudPathChanged(v string) {
	r.preferences.SetString("pcloud", v)
}
func (r *SyncUI) onSftpPathChanged(v string) {
	r.preferences.SetString("sftp", v)
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
			if appKind == "desktop" {
				authSyncProvider(r.app, provider)
			} else {
				authDropbox()
			}
		},
	}
	return btn
}

// helper function to provide sync button to given destination
func (r *SyncUI) syncButton(dst string) *widget.Button {
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
			msg := fmt.Sprintf("config: %s, sync from %s to %s", fconf, dst, vdir)
			appLog("INFO", msg, nil)
			err := ecmsync.EcmSync(fconf, dst, vdir)
			if err != nil {
				msg := fmt.Sprintf("unable to sync from %s to %s", dst, vdir)
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
			msg = fmt.Sprintf("%s records are synced successfully", dst)
			syncStatus.Set(msg)
			appLog("INFO", msg, nil)
		},
	}
	return btn
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
	pcloud := &widget.Entry{Text: "pcloud:ECM", OnSubmitted: r.onPCloudPathChanged}
	sftp := &widget.Entry{Text: "sftp:ECM", OnSubmitted: r.onSftpPathChanged}
	lpath := os.Getenv("EXTERNAL_STORAGE")
	if lpath == "" {
		lpath = "sdcard"
	}
	lpath = fmt.Sprintf("local:/%s/ECM", lpath)
	if appKind == "desktop" {
		home := os.Getenv("HOME")
		lpath = fmt.Sprintf("local:%s/.ecm", home)
	} else {
	}
	local := &widget.Entry{Text: lpath, OnSubmitted: r.onLocalPathChanged}

	dropboxSync := colorButtonContainer(r.syncButton(dropbox.Text), btnColor)
	dropboxAuth := colorButtonContainer(r.authButton("dropbox"), authColor)
	pcloudSync := colorButtonContainer(r.syncButton(pcloud.Text), btnColor)
	pcloudAuth := colorButtonContainer(r.authButton("pcloud"), authColor)
	sftpSync := colorButtonContainer(r.syncButton(sftp.Text), btnColor)
	localSync := colorButtonContainer(r.syncButton(local.Text), btnColor)
	noAuth := colorButtonContainer(r.authButton("noauth"), grayColor)

	dropboxLabel := widget.NewLabel("Dropbox to vault")
	dropboxLabel.TextStyle.Bold = true
	pcloudLabel := widget.NewLabel("PCloud to vault")
	pcloudLabel.TextStyle.Bold = true
	sftpLabel := widget.NewLabel("Sftp to vault")
	sftpLabel.TextStyle.Bold = true
	labelName := "local to vault"
	if appKind != "desktop" {
		labelName = "sdcard to vault"
	}
	localLabel := widget.NewLabel(labelName)
	localLabel.TextStyle.Bold = true

	box := container.NewVBox(
		dropboxLabel,
		container.NewGridWithColumns(3, dropbox, dropboxSync, dropboxAuth),
		pcloudLabel,
		container.NewGridWithColumns(3, pcloud, pcloudSync, pcloudAuth),
		sftpLabel,
		container.NewGridWithColumns(3, sftp, sftpSync, noAuth),
		localLabel,
		container.NewGridWithColumns(3, local, localSync, noAuth),
		statusText,
	)
	return box
}
func (r *SyncUI) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Sync", Icon: theme.HistoryIcon(), Content: r.buildUI()}
}
