package app

import (
    "errors"
    "github.com/GavinGuan24/gofer/log"
    "github.com/GavinGuan24/gofer/views"
    "github.com/gdamore/tcell"
    "os"
)

var screen tcell.Screen
var root *rootView

type ApplicationDelegate interface {
    Launched()
}

type rootView struct {
    views.View
}

func Run(delegate ApplicationDelegate) {
    //初始化 screen. (screen的内容从root view 中获取, 在图层链上的所有视图, 只有root view 的super view 为空)
    if delegate == nil {
        log.Fatal(errors.New("ApplicationDelegate is Nil"))
    }
    if screen0, e0 := tcell.NewScreen(); e0 != nil {
        log.Fatal(e0)
    } else {
        screen = screen0
    }
    if e0 := screen.Init(); e0 != nil {
        log.Fatal(e0)
    }
    //
    log.Logger("Screen init finished.")
    root = &rootView{views.NewView()}
    root.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite).Normal())
    root.SetLocation(views.NewPoint(0, 0))
    w, h := screen.Size()
    root.SetWidth(w)
    root.SetHeight(h)
    root.SetZ(0)
    //
    delegate.Launched()
    log.Logger("delegate.Launched finished.")
    //
    drawUI(nil, nil, true)
}

func Stop() {
    if screen != nil {
        screen.Fini()
        log.Logger("screen finalized.")
    }
    log.Logger("App stopped.")
    os.Exit(0)
}

func drawUI(from *views.Point, to *views.Point, needClear bool) {
    if screen == nil {
       return
    }
    if from == nil {
       from = views.NewPoint(0, 0)
    }
    if to == nil {
       to = views.NewPoint(screen.Size())
    }
    content := root.GetContent(from, to)
    if needClear {
       screen.Clear()
    }
    //TODO: draw all content to screen
    screen.Show()
}

func syncUI() {
    drawUI(nil, nil, true)
    screen.Sync()
}
