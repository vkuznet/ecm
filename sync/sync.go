package sync

import (
	"bufio"
	"context"
	"fmt"
	"os"

	// load backend modules for rclone
	_ "github.com/rclone/rclone/backend/dropbox"
	_ "github.com/rclone/rclone/backend/local"
	_ "github.com/rclone/rclone/backend/pcloud"
	_ "github.com/rclone/rclone/backend/sftp"

	// rclone libraries
	"github.com/rclone/rclone/cmd"
	"github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/config/configfile"
	"github.com/rclone/rclone/fs/operations"
	"github.com/rclone/rclone/fs/rc"
	"github.com/rclone/rclone/fs/sync"
)

// EcmSync provides a sync interface between source and destination
// The code is based on https://rclone.org/ library
func EcmSync(cpath, src, dst string) error {
	if cpath != "" {
		config.SetConfigPath(cpath)
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

// EcmCreateConfig creates sync config to be used by ecm (and rclone)
func EcmCreateConfig(cname string) error {
	file, err := os.Create(cname)
	if err != nil {
		return err
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	w.Write([]byte("[local]\n"))
	w.Write([]byte("type = local\n\n"))
	w.Write([]byte("[dropbox]\n"))
	w.Write([]byte("type = dropbox\n"))
	w.Write([]byte("client_id = dropbox\n"))
	w.Write([]byte("client_secret = dropbox\n"))
	w.Write([]byte("env_auth = true\n\n"))
	w.Write([]byte("[pcloud]\n"))
	w.Write([]byte("type = pcloud\n"))
	w.Write([]byte("hostname = api.pcloud.com\n"))
	w.Flush()
	return nil
}

// EcmUpdateConfig update given config file and add appropriate tokens to cloud providers
func EcmUpdateConfig(cpath, provider string) error {
	fmt.Println("EcmUpdateCondig: ", cpath)
	if cpath != "" {
		config.SetConfigPath(cpath)
	}
	configfile.Install()

	// add tokens to cloud providers
	in := rc.Params{}
	in["env_auth"] = true
	in["config_refresh_token"] = true
	opts := config.UpdateRemoteOpt{}
	_, err := config.UpdateRemote(context.Background(), provider, in, opts)
	if err != nil {
		return err
	}
	return nil
}
