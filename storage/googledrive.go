package storage

// TODO: I need to put actual implementations

// GoogleDriveStorage provides file-system based storage
type GoogleDriveStorage struct {
	Path string
}

func NewGoogleDriveStorage(path string) *GoogleDriveStorage {
	return &GoogleDriveStorage{Path: path}
}

// Read implements Storage.Read method
func (f *GoogleDriveStorage) Read(rid string) ([]byte, error) {
	return []byte{}, nil

}

// Write implements Storage.Write method
func (f *GoogleDriveStorage) Write(fname string, rec []byte) error {
	return nil
}

// Records implement Storage Records method
func (f *GoogleDriveStorage) Records() ([]string, error) {
	return []string{}, nil
}
