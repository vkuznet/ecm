package sync

import (
	"context"

	// load backend modules for rclone
	_ "github.com/rclone/rclone/backend/dropbox"
	_ "github.com/rclone/rclone/backend/local"
	_ "github.com/rclone/rclone/backend/pcloud"
	_ "github.com/rclone/rclone/backend/sftp"

	// rclone libraries
	"github.com/rclone/rclone/cmd"
	rconfig "github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/config/configfile"
	"github.com/rclone/rclone/fs/operations"
	"github.com/rclone/rclone/fs/sync"
)

// EcmSync provides a sync interface between source and destination
// The code is based on https://rclone.org/ library
func EcmSync(config, src, dst string) error {
	if config != "" {
		rconfig.SetConfigPath(config)
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
