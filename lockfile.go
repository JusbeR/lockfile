package lockfile

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"
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
	if err != nil {
		return fmt.Errorf("Failed to aquire lock: %v", err)
	}
	defer file.Close()
	return nil
}

// Unlock is used to free the resource
func (lockFile *Lockfile) Unlock() error {
	return os.Remove(path.Join(lockFile.path, lockFile.fileName))
}

// LockWait is like Lock(), but waits for resource to become available. If resource can't be aquired after timeout, error is returned.
// Note that this is not for realtime use. timeout is waited at least as long as stated, but it can be longer too.
// When many processes try to aquire the same lock, it is randomly selected who gets it.
// Timouts smaller than 100ms are not supported
func (lockFile *Lockfile) LockWait(timeout time.Duration) error {
	if timeout < time.Duration(time.Millisecond*100) {
		return fmt.Errorf("Invalid timeout(%v)", timeout)
	}
	now := time.Now()
	endWaiting := now.Add(timeout)
	var err error
	for now.Before(endWaiting) {
		err = lockFile.Lock()
		if err != nil {
			time.Sleep(time.Duration(time.Millisecond * 100))
			now = time.Now()
			continue
		}
		return nil
	}
	return err
}
