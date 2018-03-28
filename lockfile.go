package lockfile

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

const defaultLockFileName = ".lockfile-gthsf4563"

type Lockfile struct {
	path     string
	fileName string
}

// NewLockFile is used to initiate new lockfile. path should be known directory for each process that uses same lock.
// path should be writable. You can also define absolute file path+name for lock. e.g. '/etc/mylocks/.mylock'
func NewLockFile(path string) (Lockfile, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		dir, fileName := filepath.Split(path)
		_, err := os.Stat(dir)
		if err != nil {
			return Lockfile{}, fmt.Errorf("Invalid path/filename given(%v)", path)
		}
		return Lockfile{path: dir, fileName: fileName}, nil
	}
	if fileInfo.IsDir() {
		return Lockfile{path: path, fileName: defaultLockFileName}, nil
	}
	dir, fileName := filepath.Split(path)
	_, err = os.Stat(dir)
	if err != nil {
		return Lockfile{}, fmt.Errorf("Invalid path/filename given(%v)", path)
	}
	return Lockfile{path: dir, fileName: fileName}, nil
}

// Lock is used to aquire new lock. if err != nil lock could not be aquired and you should try it later again.
func (lockFile *Lockfile) Lock() error {
	file, err := os.OpenFile(path.Join(lockFile.path, lockFile.fileName), os.O_RDONLY|os.O_CREATE|os.O_EXCL, 0666)
	defer file.Close()
	return err
}

// Unlock is used to free the resource
func (lockFile *Lockfile) Unlock() error {
	return os.Remove(path.Join(lockFile.path, lockFile.fileName))
}
