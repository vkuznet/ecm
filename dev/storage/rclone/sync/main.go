package main

import (
	"context"
	"flag"
	"fmt"

	_ "github.com/rclone/rclone/backend/all"
	//     _ "github.com/rclone/rclone/backend/dropbox"
	//     _ "github.com/rclone/rclone/backend/local"
	"github.com/rclone/rclone/cmd"
	rconfig "github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/config/configfile"
	"github.com/rclone/rclone/fs/operations"
	"github.com/rclone/rclone/fs/sync"
)

func main() {
	err := EcmSync()
	fmt.Println("sync error", err)
}
func EcmSync() error {
	var config string
	flag.StringVar(&config, "config", "", "config file")
	var dst string
	flag.StringVar(&dst, "dst", "", "destination folder, e.g. dropbox:TMP")
	var src string
	flag.StringVar(&src, "src", "", "source folder, e.g. local:/tmp")
	flag.Parse()

	if config != "" {
		rconfig.SetConfigPath(config)
		fmt.Println("config", rconfig.GetConfigPath())
	}
	configfile.Install()

	args := []string{src, dst}
	createEmptySrcDirs := true
	fsrc, srcFileName, fdst := cmd.NewFsSrcFileDst(args)
	if srcFileName == "" {
		return sync.Sync(context.Background(), fdst, fsrc, createEmptySrcDirs)
	}
	return operations.CopyFile(context.Background(), fdst, fsrc, srcFileName, srcFileName)

}
