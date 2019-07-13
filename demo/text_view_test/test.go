package main

import (
    "fmt"
    "github.com/GavinGuan24/gofer/gofer"
    "github.com/GavinGuan24/gofer/widget"
    "github.com/gdamore/tcell"
)

type agent struct {
    v1 gofer.View
    v2 gofer.View
}

func (a *agent) Launched(root gofer.RootView) {

    bg := gofer.NewView()
    bg.SetWidth(5)
    bg.SetHeight(3)
    bg.SetLocation(gofer.NewPoint(1, 1))
    bg.SetStyle(tcell.StyleDefault.Background(tcell.NewRGBColor(225,225,225)))
    root.AddSubview(bg, root)

    textview2 := widget.NewTextView(3)
    a.v2 = textview2
    textview2.SetText("a中")
    textview2.SetLocation(gofer.NewPoint(1, 1))
    textview2.SetStyle(tcell.StyleDefault.Foreground(tcell.NewRGBColor(255,250,227)).Background(tcell.NewRGBColor(0,0,0)))
    bg.AddSubview(textview2, bg)

    textview := widget.NewTextView(3)
    a.v1 = textview
    textview.SetText("a中")
    textview.SetLocation(gofer.NewPoint(6, 5))
    textview.SetStyle(tcell.StyleDefault.Foreground(tcell.NewRGBColor(255,250,227)).Background(tcell.NewRGBColor(0,0,0)))
    root.AddSubview(textview, root)
}

func (a *agent) WillStop(code int) bool {
    return true
}

func (a *agent) EnableMouse() bool {
    return false
}

func (a *agent) EventListener() chan<- tcell.Event {
    listener := make(chan tcell.Event)
    go func() {
        for {
            select {
            case event := <-listener:
                switch event := event.(type) {
                case *tcell.EventKey:
                    switch {
                    case event.Key() == tcell.KeyLeft:
                        a.v1.SetLocation(a.v1.Location().Left(1))
                        gofer.UpdateUI(a.v1, nil)
                    case event.Key() == tcell.KeyRight:
                        a.v1.SetLocation(a.v1.Location().Right(1))
                        gofer.UpdateUI(a.v1, nil)
                    case event.Key() == tcell.KeyUp:
                        a.v1.SetLocation(a.v1.Location().Up(1))
                        gofer.UpdateUI(a.v1, nil)
                    case event.Key() == tcell.KeyDown:
                        a.v1.SetLocation(a.v1.Location().Down(1))
                        gofer.UpdateUI(a.v1, nil)
                    case event.Rune() == 'w':
                        a.v2.SetLocation(a.v2.Location().Up(1))
                        gofer.UpdateUI(a.v2, nil)
                    case event.Rune() == 's':
                        a.v2.SetLocation(a.v2.Location().Down(1))
                        gofer.UpdateUI(a.v2, nil)
                    case event.Rune() == 'a':
                        a.v2.SetLocation(a.v2.Location().Left(1))
                        gofer.UpdateUI(a.v2, nil)
                    case event.Rune() == 'd':
                        a.v2.SetLocation(a.v2.Location().Right(1))
                        gofer.UpdateUI(a.v2, nil)
                    }
                case *tcell.EventMouse:
                    gofer.LogDebug(fmt.Sprintf("Event: Mouse."))
                case *tcell.EventResize:
                    gofer.LogDebug(fmt.Sprintf("Event: Resize Window."))
                default:
                    gofer.LogWarn(fmt.Sprintf("Unknown event %v", event))
                }
            default:
            }
        }
    }()
    return listener
}

func main() {
    gofer.NewApp().Run(&agent{})
}
