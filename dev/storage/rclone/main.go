package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	_ "github.com/rclone/rclone/backend/all"
	//     _ "github.com/rclone/rclone/backend/dropbox"
	//     _ "github.com/rclone/rclone/backend/local"
	"github.com/rclone/rclone/cmd"
	rconfig "github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/config/configfile"
	"github.com/rclone/rclone/fs/operations"
)

func main() {
	var config string
	flag.StringVar(&config, "config", "", "config file")
	var dst string
	flag.StringVar(&dst, "dst", "", "ls destination folder, e.g. dropbox:TMP")
	flag.Parse()

	if config != "" {
		rconfig.SetConfigPath(config)
	}
	fmt.Println("### config", rconfig.GetConfigPath())
	configfile.Install()

	args := []string{dst}
	fmt.Println("### args", args)
	fsrc := cmd.NewFsSrc(args)
	res := operations.List(context.Background(), fsrc, os.Stdout)
	fmt.Println("results", res)
}
