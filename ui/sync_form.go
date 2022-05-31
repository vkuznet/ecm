package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	widget "fyne.io/fyne/v2/widget"
	ecmsync "github.com/vkuznet/ecm/sync"
)

// SyncUI represents UI SyncUI
type SyncUI struct {
	preferences  fyne.Preferences
	window       fyne.Window
	app          fyne.App
	vaultRecords *vaultRecords
}

func newUISync(a fyne.App, w fyne.Window, v *vaultRecords) *SyncUI {
	return &SyncUI{
		app:         a,
		window:      w,
		preferences: a.Preferences(),
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
				log.Println("unable to sync", err)
			}
			log.Println("records are synced")
			// reset vault records
			_vault.Records = nil
			// read again vault records
			err = _vault.Read()
			if err != nil {
				log.Println("unable to read the vault records", err)
			}
			// refresh ui records
			r.vaultRecords.Refresh()
		},
	}
	return btn
}

// helper function to build UI
func (r *SyncUI) buildUI() *container.Scroll {

	// sync form container
	dropbox := &widget.Entry{Text: "dropbox:ECM", OnSubmitted: r.onDropboxPathChanged}
	pcloud := &widget.Entry{Text: "pcloud:ECM", OnSubmitted: r.onPCloudPathChanged}
	sftp := &widget.Entry{Text: "sftp:ECM", OnSubmitted: r.onSftpPathChanged}
	dstDir := _vault.Directory
	if appKind != "desktop" {
		dstDir = "mobile storage"
	}
	dst := &widget.Entry{Text: dstDir}
	dst.Disable()

	// button to sync
	btnTo := &widget.Button{
		Text: "",
		Icon: rightArrowImage.Resource,
	}
	//     btnTo := canvas.NewImageFromResource(rightArrowImage.Resource)
	btnToContainer := colorButtonContainer(btnTo, btnColor)

	dropboxButtonContainer := colorButtonContainer(r.syncButton(dropbox.Text), btnColor)
	rowDropbox := container.NewHBox(dropbox, btnToContainer, dst, dropboxButtonContainer)
	pcloudButtonContainer := colorButtonContainer(r.syncButton(pcloud.Text), btnColor)
	rowPCloud := container.NewHBox(pcloud, btnToContainer, dst, pcloudButtonContainer)
	sftpButtonContainer := colorButtonContainer(r.syncButton(sftp.Text), btnColor)
	rowSftp := container.NewHBox(sftp, btnToContainer, dst, sftpButtonContainer)

	box := container.NewVBox(
		//         container.NewGridWithColumns(2, dropbox, r.syncButton),
		widget.NewLabel("Dropbox"),
		rowDropbox,
		widget.NewLabel("Pcloud"),
		rowPCloud,
		widget.NewLabel("Sftp"),
		rowSftp,
		&widget.Label{},
	)

	return container.NewScroll(box)
}
func (r *SyncUI) tabItem() *container.TabItem {
	//     return &container.TabItem{Text: "Sync", Icon: theme.DownloadIcon(), Content: r.buildUI()}
	return &container.TabItem{Text: "Sync", Icon: syncImage.Resource, Content: r.buildUI()}
}
