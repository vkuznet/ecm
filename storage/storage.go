package storage

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Storage defines generic storage interface
type Storage interface {
	// Read reades from storage
	Read(rid string) ([]byte, error)
	// Write writes given record to storage using givne file name
	Write(fname string, rec []byte) error
	// Records return list of record ids from storage
	Records() ([]string, error)
}

// FileStorage provides file-system based storage
type FileStorage struct {
	Path string
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{Path: path}
}

// Read implements Storage.Read method
func (f *FileStorage) Read(rid string) ([]byte, error) {
	files, err := ioutil.ReadDir(f.Path)
	if err != nil {
		return []byte{}, err
	}
	var base, fname, ext string
	for _, f := range files {
		base = filepath.Base(f.Name())
		ext = filepath.Ext(f.Name())
		fname = strings.Replace(base, ext, "", -1)
		if fname == rid {
			break
		}
	}
	fileName := filepath.Join(f.Path, base)
	data, err := os.ReadFile(fileName)
	return data, err

}

// Write implements Storage.Write method
func (f *FileStorage) Write(fname string, rec []byte) error {
	fileName := filepath.Join(f.Path, fname)
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("unable to create file name", fileName, " error ", err)
		return err
	}
	defer file.Close()
	_, err = file.Write(rec)
	return err
}

// Records implement Storage Records method
func (f *FileStorage) Records() ([]string, error) {
	var records []string
	files, err := ioutil.ReadDir(f.Path)
	if err != nil {
		return records, err
	}
	for _, f := range files {
		base := filepath.Base(f.Name())
		ext := filepath.Ext(f.Name())
		fname := strings.Replace(base, ext, "", -1)
		records = append(records, fname)
	}
	return records, nil
}
