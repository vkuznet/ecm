package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	binding "fyne.io/fyne/v2/data/binding"
	widget "fyne.io/fyne/v2/widget"
	ecmsync "github.com/vkuznet/ecm/sync"
)

// global variable we'll use to update sync status
var syncStatus binding.String

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

// helper function to provide sync button to given destination
func (r *SyncUI) syncButton(dst string) *widget.Button {
	// get vault dir from preferences
	pref := r.app.Preferences()
	vdir := pref.String("VaultDirectory")

	btn := &widget.Button{
		Text: "Sync",
		Icon: syncImage.Resource,
		OnTapped: func() {
			// perform sync from dropbox to vault
			dir := r.app.Storage().RootURI().Path()
			fconf := fmt.Sprintf("%s/ecmsync.conf", dir)
			log.Println("config", fconf)
			log.Printf("sync from %s to %s", dst, vdir)
			err := ecmsync.EcmSync(fconf, dst, vdir)
			if err != nil {
				msg := fmt.Sprintf("unable to sync, %v", err)
				syncStatus.Set(msg)
				log.Println(msg)
			}
			log.Println("records are synced")
			// reset vault records
			_vault.Records = nil
			// read again vault records
			err = _vault.Read()
			if err != nil {
				msg := fmt.Sprintf("unable to read the vault records, %v", err)
				syncStatus.Set(msg)
				log.Println(msg)
			}
			// refresh ui records
			r.vaultRecords.Refresh()
			msg := fmt.Sprintf("%si records are synced successfully", dst)
			syncStatus.Set(msg)
		},
	}
	return btn
}

// helper function to build UI
func (r *SyncUI) buildUI() *container.Scroll {

	// create text box which will update text once sync is completed
	syncStatus = binding.NewString()
	syncStatus.Set("Sync status will appear here")
	statusText := widget.NewLabelWithData(syncStatus)

	// sync form container
	dropbox := &widget.Entry{Text: "dropbox:ECM", OnSubmitted: r.onDropboxPathChanged}
	pcloud := &widget.Entry{Text: "pcloud:ECM", OnSubmitted: r.onPCloudPathChanged}
	sftp := &widget.Entry{Text: "sftp:ECM", OnSubmitted: r.onSftpPathChanged}

	dropboxSync := colorButtonContainer(r.syncButton(dropbox.Text), btnColor)
	pcloudSync := colorButtonContainer(r.syncButton(pcloud.Text), btnColor)
	sftpSync := colorButtonContainer(r.syncButton(sftp.Text), btnColor)

	dropboxLabel := widget.NewLabel("Dropbox to vault")
	dropboxLabel.TextStyle.Bold = true
	pcloudLabel := widget.NewLabel("PCloud to vault")
	pcloudLabel.TextStyle.Bold = true
	sftpLabel := widget.NewLabel("Sftp to vault")
	sftpLabel.TextStyle.Bold = true

	box := container.NewVBox(
		dropboxLabel,
		container.NewGridWithColumns(2, dropbox, dropboxSync),
		pcloudLabel,
		container.NewGridWithColumns(2, pcloud, pcloudSync),
		sftpLabel,
		container.NewGridWithColumns(2, sftp, sftpSync),
		statusText,
		//         &widget.Label{},
	)

	return container.NewScroll(box)
}
func (r *SyncUI) tabItem() *container.TabItem {
	//     return &container.TabItem{Text: "Sync", Icon: theme.DownloadIcon(), Content: r.buildUI()}
	return &container.TabItem{Text: "Sync", Icon: syncImage.Resource, Content: r.buildUI()}
}
