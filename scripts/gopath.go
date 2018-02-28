package main

import (
    "fmt"
    "go/build"
    "os"
)

func main() {
    gopath := os.Getenv("GOPATH")
    if gopath == "" {
        gopath = build.Default.GOPATH
    }
    fmt.Println(gopath)
}
