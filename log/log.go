package log

import (
    "fmt"
    "os"
    "time"
)

var logger *os.File

func init() {
    file, e := os.OpenFile("./gofer.log", os.O_CREATE|os.O_RDWR, 0644)
    if e != nil {
        Fatal(e)
    }
    logger = file
}

func Fatal(e error) {
    fmt.Fprint(os.Stderr, "%v\n", e.Error())
    os.Exit(1)
}

func Logger(str string) {
    fmt.Fprintf(logger, "%s -- %s\n", time.Now().Format("2006-01-02 15:04:05.000"), str)
}

