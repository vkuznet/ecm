package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	theme "fyne.io/fyne/v2/theme"
	widget "fyne.io/fyne/v2/widget"
	ecmsync "github.com/vkuznet/ecm/sync"
)

// SyncUI represents UI SyncUI
type SyncUI struct {
	preferences  fyne.Preferences
	window       fyne.Window
	app          fyne.App
	dropbox      *widget.Entry
	syncButton   *widget.Button
	vaultRecords *vaultRecords
}

func newSyncUI(a fyne.App, w fyne.Window, v *vaultRecords) *SyncUI {
	return &SyncUI{
		app:         a,
		window:      w,
		preferences: a.Preferences(),
	}
}

func (r *SyncUI) onDropboxPathChanged(v string) {
	r.preferences.SetString("dropbox", v)
}

// helper function to build UI
func (r *SyncUI) buildUI() *container.Scroll {

	// sync form container
	dpath := "dropbox:ECM"
	r.dropbox = &widget.Entry{Text: dpath, OnSubmitted: r.onDropboxPathChanged}

	// get vault dir from preferences
	pref := r.app.Preferences()
	vdir := pref.String("VaultDirectory")

	r.syncButton = &widget.Button{
		Text: "Sync",
		Icon: theme.DownloadIcon(),
		OnTapped: func() {
			// perform sync from dropbox to vault
			dir := r.app.Storage().RootURI().Path()
			fconf := fmt.Sprintf("%s/ecmsync.conf", dir)
			log.Println("config", fconf)
			log.Printf("sync from %s to %s", dpath, vdir)
			err := ecmsync.EcmSync(fconf, dpath, vdir)
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
	box := container.NewVBox(
		container.NewGridWithColumns(2, r.dropbox, r.syncButton),
		&widget.Label{},
	)

	return container.NewScroll(box)
}
func (r *SyncUI) tabItem() *container.TabItem {
	return &container.TabItem{Text: "", Icon: theme.DownloadIcon(), Content: r.buildUI()}
}
