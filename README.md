# Simple file lock

Simple file lock to share a resource between processes.

## Get it
    go get github.com/JusbeR/lockfile
    
## Usage
    import "github.com/JusbeR/lockfile"
    ...
    resourceLock, err := lockfile.NewLockFile("./resourcelock")
    if err != nil {
        log.Fatalln("Could not create lock, maybe file path is unaccessible or something:", err)
    }
    if err := resourceLock.Lock(); err != nil {
        log.Println("Someone else is using the resource")
    } else {
        defer resourceLock.Unlock()
        // Use resource
    }

Or alternatively

    ...
    if err := resourceLock.LockWait( 60 * time.Second ); err != nil {
        log.Println("Someone else is using the resource and not giving it up even after one minute of waiting")
    } else {
    ...