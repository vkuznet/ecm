package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Copy performs backup copy of source to destination
// based on https://github.com/mactsouk/opensource.com/blob/master/cp1.go
func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		//         log.Printf("file '%s' does not exist, error %v", src, err)
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	err = os.Chmod(dst, 0600)
	if err != nil {
		log.Println("unable to change file permission of", dst)
	}

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// Files return list of vault files
func Files(vdir string) ([]string, error) {
	files, err := os.ReadDir(vdir)
	if err != nil {
		return []string{}, err
	}
	var out []string
	for _, f := range files {
		if f.Name() != "backups" {
			out = append(out, f.Name())
		}
	}
	return out, nil
}

// BackupFile perform backup copy of given file in vault dir
func BackupFile(vdir, fid, bdir string) error {
	fname := fmt.Sprintf("%s", filepath.Join(vdir, fid))
	if _, err := os.Stat(fname); err != nil {
		// backup file name with existing cipher
		tstamp := time.Now().Format(time.RFC3339)
		bname := filepath.Join(bdir, fmt.Sprintf("%s-%s", fid, tstamp))
		// make backup of our record
		_, err = Copy(fname, bname)
		return err
	}
	return nil
}

// Backup backups vault directory
func Backup(vdir string, verbose int) error {
	// create backups vault area
	bdir := filepath.Join(vdir, "backups")
	err := os.MkdirAll(bdir, 0755)
	if err != nil {
		log.Printf("unable to create %s, error %v", bdir, err)
		return err
	}

	// get list of files in vault area
	files, err := Files(vdir)
	if err != nil {
		return err
	}

	for _, fid := range files {
		// backup existing record if it exists
		err = BackupFile(vdir, fid, bdir)
		if err != nil {
			if verbose > 0 {
				log.Println("unable to make backup for record", fid, " error ", err)
			}
		}
	}
	return nil
}
