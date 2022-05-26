package storage

// TODO: I need to put actual implementations

// SSHStorage provides file-system based storage
type SSHStorage struct {
	Path string
}

func NewSSHStorage(path string) *SSHStorage {
	return &SSHStorage{Path: path}
}

// Read implements Storage.Read method
func (f *SSHStorage) Read(rid string) ([]byte, error) {
	return []byte{}, nil

}

// Write implements Storage.Write method
func (f *SSHStorage) Write(fname string, rec []byte) error {
	return nil
}

// Records implement Storage Records method
func (f *SSHStorage) Records() ([]string, error) {
	return []string{}, nil
}
