# Autorestart

`autorestart` is used for autorestarting Go apps in development or staging environments, but you can try to use in production at your own risk. Works well with `go build` and `rsync`.
It designed to be as lightweight as possible, it does not uses `fsnotify` it just periodically poll `os.Stat(filename)`

Where is `filename` is a self binary by default, but you can setup to watch `tmp/restart.txt` or something else.
On file change it will call `syscall.Exec(selfbinary)` or you can use function `SendSIGUSR2` (useful for grace restart) or write your own.

## Quick start

    go get github.com/slayer/autorestart

#### Basic usage

```go
package main

import "github.com/slayer/autorestart"

func main() {
    autorestart.StartWatcher()
    http.ListenAndServe(":8080", nil) // for example
}
```


#### Extended usage

```go
package main

import (
    "log"
    "http"
    "github.com/slayer/autorestart"
)

func main() {
    // set period
    autorestart.WatchPeriod = 3 * time.Second
    // custom file to watch
    autorestart.WatchFilename = "tmp/restart.txt"
    // custom restart function
    autorestart.RestartFunc = autorestart.SendSIGUSR2 // usefull for `github.com/facebookgo/grace`

    // or
    autorestart.RestartFunc = func () {
        if proc, err := os.FindProcess(os.Getpid()); err == nil {
            proc.Signal(syscall.SIGHUP)
        }
    }

    // Notifier
    restart := autorestart.GetNotifier()
    go func() {
        <- restart
        log.Printf("I will restart shortly")
    }()

    autorestart.StartWatcher()
    http.ListenAndServe(":8080", nil) // for example
}
```


## Licence

Apache 2.0
