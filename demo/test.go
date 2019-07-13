package main

import (
    "fmt"
    "github.com/GavinGuan24/gofer/gofer"
    "github.com/gdamore/tcell"
)

func main() {
    myRunes := gofer.StringToRunes("abc, 我是", tcell.StyleDefault)
    for _, vv := range myRunes {
        fmt.Printf("%c, %v\n", vv.Mainc(), vv.Combc())
    }
}


