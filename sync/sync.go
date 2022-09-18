package sync

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	// ecm packages
	"github.com/vkuznet/ecm/utils"

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
// The code is based on https://rclone.org/ library and relies on sync module
func EcmSync(cpath, src, dst string) error {
	// create backup of our destination area
	err := utils.Backup(dst, 0) // non-verbose
	if err != nil {
		return err
	}

	// setup configuration for rclone
	if cpath != "" {
		config.SetConfigPath(cpath)
	}
	configfile.Install()

	// perform sync call of rclone, see https://rclone.org/commands/rclone_sync/
	args := []string{src, dst}
	createEmptySrcDirs := true
	fsrc, srcFileName, fdst := cmd.NewFsSrcFileDst(args)
	if srcFileName == "" {
		// if we need to sync: it will wiped out files at destination if they don't exist at src
		// return sync.Sync(context.Background(), fdst, fsrc, createEmptySrcDirs)
		// if we need to copy: it will preserve files at destination
		return sync.CopyDir(context.Background(), fdst, fsrc, createEmptySrcDirs)
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

// ServerRecord represents HTTP vault record
type ServerRecord struct {
	ID   string
	Data []byte
}

// helper function to perform sync operation from HTTP end-point
func SyncFromServer(rurl, dst string) error {
	// create backup of our destination area
	err := utils.Backup(dst, 0) // non-verbose
	if err != nil {
		return err
	}

	// perform HTTP call to our server and create new records at destination
	client := &http.Client{}
	resp, err := client.Get(rurl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// records represent list of HTTP vault records
	var records []ServerRecord
	err = json.Unmarshal(data, &records)
	if err != nil {
		return err
	}
	for _, rec := range records {
		fname := fmt.Sprintf("%s/%s", dst, rec.ID)
		file, err := os.Create(fname)
		if err != nil {
			return err
		}
		defer file.Close()
		file.Write(rec.Data)
	}
	return nil
}
