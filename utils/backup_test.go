package utils

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"testing"
)

// TestCopy function
func TestCopy(t *testing.T) {
	// Copy(src, dst string) (int64, error)

	// create and open a temporary file
	f, err := os.CreateTemp("", "tmpfile-")
	if err != nil {
		t.Fatal(err)
	}

	// close and remove the temporary file at the end of the program
	defer f.Close()
	defer os.Remove(f.Name())
	// write data to the temporary file
	data := []byte("test")
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}

	dst := fmt.Sprintf("%s.copy", f.Name())
	_, err = Copy(f.Name(), dst)
	if err != nil {
		t.Fatal(err)
	}

	// perform diff operation
	out, err := exec.Command("diff", f.Name(), dst).Output()
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "" {
		t.Fatalf("Files %s %s are different with %s", f.Name(), dst, string(out))
	}
}

// TestFiles function
func TestFiles(t *testing.T) {
	// Files(vdir string) ([]string, error)
	dir := "/etc"
	files, err := Files(dir)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(files)

	out, err := exec.Command("ls", dir).Output()
	if err != nil {
		t.Fatal(err)
	}
	if len(out) == 0 {
		t.Fatalf("Unable to list files in %s", dir)
	}
	var ufiles []string
	for _, f := range strings.Split(string(out), "\n") {
		fname := strings.Replace(f, "\n", "", -1)
		fname = strings.Trim(fname, " ")
		if fname == "" {
			continue
		}
		ufiles = append(ufiles, fname)
	}
	sort.Strings(ufiles)
	if len(files) != len(ufiles) {
		t.Fatalf("Files API fails:\nprovide\n%v\nexpect\n%v\n", files, ufiles)
	}

}

// TestBackupFile function
func TestBackupFile(t *testing.T) {
	// BackupFile(vdir, fid, bdir string) error

	fid := "123"
	tdir := strings.Replace(fmt.Sprintf("%s/ecm-test", os.TempDir()), "//", "/", -1)
	err := os.Mkdir(tdir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	fname := fmt.Sprintf("%s/%s", tdir, fid)
	bdir := strings.Replace(fmt.Sprintf("%s/backups", tdir), "//", "/", -1)
	err = os.Mkdir(bdir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	// create a file in temp dir
	f, err := os.Create(fname)
	if err != nil {
		t.Fatal(err)
	}

	// cleanup tdir and bdir
	defer os.RemoveAll(bdir)
	defer os.RemoveAll(tdir)

	// close and remove the temporary file at the end of the program
	defer f.Close()
	defer os.Remove(fname)
	// write data to the temporary file
	data := []byte("test")
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}
	// compare dirs
	tout, err := Files(tdir)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Content of %s, %v\n", tdir, tout)

	err = BackupFile(tdir, fid, bdir)
	if err != nil {
		t.Fatal(err)
	}

	bsubdir, err := BackupTDir(bdir)
	if err != nil {
		t.Fatal(err)
	}
	bout, err := Files(bsubdir)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Content of %s, %v\n", bsubdir, bout)
	if strings.Join(tout, " ") != strings.Join(bout, " ") {
		t.Fatalf("not equal content: %v != %v", tout, bout)
	}
}

// TestBackup function
func TestBackup(t *testing.T) {
	// Backup(vdir string, verbose int) error

	fid := "123"
	tdir := strings.Replace(fmt.Sprintf("%s/ecm-test", os.TempDir()), "//", "/", -1)
	err := os.Mkdir(tdir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	// cleanup tdir
	defer os.RemoveAll(tdir)

	// create a file in temp dir
	fname := fmt.Sprintf("%s/%s", tdir, fid)
	f, err := os.Create(fname)
	if err != nil {
		t.Fatal(err)
	}

	// close and remove the temporary file at the end of the program
	defer f.Close()
	defer os.Remove(fname)
	// write data to the temporary file
	data := []byte("test")
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}

	// content of tdir
	tout, err := Files(tdir)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Content of %s, %v\n", tdir, tout)
	verbose := 1

	err = Backup(tdir, verbose)
	if err != nil {
		t.Fatal(err)
	}

	bdir := strings.Replace(fmt.Sprintf("%s/backups", tdir), "//", "/", -1)
	bsubdir, err := BackupTDir(bdir)
	if err != nil {
		t.Fatal(err)
	}
	bout, err := Files(bsubdir)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Content of %s, %v\n", bsubdir, bout)

	// compare content of two dirs
	if strings.Join(tout, " ") != strings.Join(bout, " ") {
		t.Fatalf("not equal content: %v != %v", tout, bout)
	}
}
