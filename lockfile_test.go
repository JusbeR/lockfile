package lockfile

import (
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLockFileWhenPathDoesNotExist(t *testing.T) {
	_, err := NewLockFile("/path/that/does/not/exist")
	assert.NotNil(t, err)
}

func TestNewLockFileWhenPathExists(t *testing.T) {
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	lockFile, err := NewLockFile(currentDir)
	assert.Nil(t, err)
	assert.Equal(t, currentDir, lockFile.path)
	assert.Equal(t, defaultLockFileName, lockFile.fileName)
}

func TestNewLockFileWhenPathAndFilenameGiven(t *testing.T) {
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	absolutePathPlusFileName := path.Join(currentDir, "test")
	lockFile, err := NewLockFile(absolutePathPlusFileName)
	assert.Nil(t, err)
	assert.Equal(t, path.Clean(currentDir), path.Clean(lockFile.path))
	assert.Equal(t, "test", lockFile.fileName)
}

func TestNewLockFileWhenOnlyFileNameGiven(t *testing.T) {
	fileName := "test"
	_, err := NewLockFile(fileName)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Invalid path/filename given(test)")
}

func TestLockFileHappyPath(t *testing.T) {
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	lockFile, _ := NewLockFile(currentDir)
	err := lockFile.Lock()
	assert.Nil(t, err)
	err = lockFile.Unlock()
	assert.Nil(t, err)
}

func TestLockFileWhenAutoNamedFileAlreadyMade(t *testing.T) {
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	lockFile, err := NewLockFile(currentDir)
	assert.Nil(t, err)
	assert.Equal(t, defaultLockFileName, lockFile.fileName)
	err = lockFile.Lock()
	assert.Nil(t, err)
	defer lockFile.Unlock()
	lockFile2, err := NewLockFile(currentDir)
	assert.Nil(t, err)
	assert.Equal(t, defaultLockFileName, lockFile2.fileName)
	err = lockFile2.Lock()
	assert.NotNil(t, err)
}

func TestLockFileWhenSelfNamedFileAlreadyMade(t *testing.T) {
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	absolutePathPlusFileName := path.Join(currentDir, "test")
	lockFile, err := NewLockFile(absolutePathPlusFileName)
	assert.Nil(t, err)
	assert.Equal(t, "test", lockFile.fileName)
	err = lockFile.Lock()
	assert.Nil(t, err)
	defer lockFile.Unlock()
	lockFile2, err := NewLockFile(absolutePathPlusFileName)
	assert.Nil(t, err)
	assert.Equal(t, "test", lockFile2.fileName)
	err = lockFile2.Lock()
	assert.NotNil(t, err)
}

func TestLockFileWhenLockingTwice(t *testing.T) {
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	lockFile, _ := NewLockFile(currentDir)
	err := lockFile.Lock()
	assert.Nil(t, err, "Should lock normally")
	err = lockFile.Lock()
	assert.NotNil(t, err, "Should not lock twice")
	err = lockFile.Unlock()
	assert.Nil(t, err, "Should unlock normally")
	err = lockFile.Lock()
	assert.Nil(t, err, "Should lock normally after unlock")
	err = lockFile.Unlock()
	assert.Nil(t, err, "Should unlock normally")
}

func TestLockWaitWhenResourceIsNotTaken(t *testing.T) {
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	lockFile, _ := NewLockFile(currentDir)
	err := lockFile.LockWait(time.Duration(time.Second * 5))
	assert.Nil(t, err)
	err = lockFile.Unlock()
	assert.Nil(t, err)
}
func unlockAfterDuration(lockfile *Lockfile, duration time.Duration) {
	time.Sleep(duration)
	lockfile.Unlock()
}
func TestLockWaitWhenResourceIsTakenWaitTimeout(t *testing.T) {
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	lockFile, _ := NewLockFile(currentDir)
	err := lockFile.Lock()
	assert.Nil(t, err)
	go unlockAfterDuration(&lockFile, time.Duration(time.Millisecond*200))
	lockFile2, _ := NewLockFile(currentDir)
	err = lockFile2.LockWait(time.Duration(time.Millisecond * 400))
	assert.Nil(t, err)
	err = lockFile2.Unlock()
	assert.Nil(t, err)
}

func TestLockWaitWhenResourceIsTakenTimoutExpires(t *testing.T) {
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	lockFile, _ := NewLockFile(currentDir)
	err := lockFile.Lock()
	assert.Nil(t, err)
	go unlockAfterDuration(&lockFile, time.Duration(time.Millisecond*400))
	lockFile2, _ := NewLockFile(currentDir)
	err = lockFile2.LockWait(time.Duration(time.Millisecond * 100))
	assert.NotNil(t, err)
	// TODO: a proper time pkg mocking needed.
	time.Sleep(time.Millisecond * 500)
}

func TestLockWaitWhenZeroDuration(t *testing.T) {
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	lockFile, _ := NewLockFile(currentDir)
	err := lockFile.LockWait(time.Duration(0))
	assert.NotNil(t, err)
}

func TestLockWaitWhenTooSmallDuration(t *testing.T) {
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	lockFile, _ := NewLockFile(currentDir)
	err := lockFile.LockWait(time.Duration(time.Millisecond * 90))
	assert.NotNil(t, err)
}
