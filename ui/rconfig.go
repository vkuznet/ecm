package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	_ "embed"

	fyne "fyne.io/fyne/v2"
)

//go:embed "rclone.conf"
var ecmconfig []byte

func WriteSyncConfig(app fyne.App) {
	dir := app.Storage().RootURI().Path()
	fname := fmt.Sprintf("%s/ecmsync.conf", dir)
	file, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	w.Write(ecmconfig)
	w.Flush()
}
