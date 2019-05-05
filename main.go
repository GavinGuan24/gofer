package main

import (
    "github.com/GavinGuan24/gofer/app"
    "time"
)

func main() {
    app.Run(&appDele{})
    time.Sleep(time.Duration(time.Second*5))
    app.Stop()
}

type appDele struct {

}

func (ad *appDele) Launched() {
}

