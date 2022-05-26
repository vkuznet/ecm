package storage

// TODO: I need to put actual implementations

// DropboxStorage provides file-system based storage
type DropboxStorage struct {
	Path string
}

func NewDropboxStorage(path string) *DropboxStorage {
	return &DropboxStorage{Path: path}
}

// Read implements Storage.Read method
func (f *DropboxStorage) Read(rid string) ([]byte, error) {
	return []byte{}, nil

}

// Write implements Storage.Write method
func (f *DropboxStorage) Write(fname string, rec []byte) error {
	return nil
}

// Records implement Storage Records method
func (f *DropboxStorage) Records() ([]string, error) {
	return []string{}, nil
}
