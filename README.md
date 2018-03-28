# Simple file lock

Simple file lock to share a resource between processes.

## Usage
    resourceLock, err := NewLockFile("/etc/myswlocks/")
    if err != nil {
        log.Fatalln("Could not create lock, maybe file path is unaccessible or something:", err)
    }
    if err := resourceLock.Lock(); err != nil {
        log.Println("Someone else is using the resource")
    }
    defer resourceLock.Unlock()
    // Use resource
