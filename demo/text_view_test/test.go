package main

import (
    "fmt"
    "github.com/GavinGuan24/gofer/gofer"
    "github.com/GavinGuan24/gofer/widget"
    "github.com/gdamore/tcell"
)

var _root gofer.RootView

type agent struct {
    v1 *widget.TextView
}

func (a *agent) Launched(root gofer.RootView) {
    textview := widget.NewTextView(3)
    textview.SetWidth(36)
    a.v1 = textview
    textview.SetText("type your text here in English mode.")
    textview.SetLocation(gofer.NewPoint(6, 5))
    textview.SetStyle(tcell.StyleDefault.Foreground(tcell.NewRGBColor(255, 250, 227)).Background(tcell.NewRGBColor(0, 0, 0)))
    root.AddSubview(textview, root)
    _root = root
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
                    case event.Key() == tcell.KeyEnter:
                        a.v1.SetText("")
                        gofer.UpdateUI(_root, nil)
                    case event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 || event.Key() == tcell.KeyDelete:
                        runes := []rune(a.v1.Text())
                        length := len(runes)
                        if length > 0 {
                            a.v1.SetText(string(runes[:length-1]))
                            gofer.UpdateUI(_root, nil)
                        }
                    case validRune(event.Rune()):
                        runes := []rune(a.v1.Text())
                        if a.v1.TextWidth() >= 36 {
                            runes = []rune{}
                        }
                        runes = append(runes, event.Rune())
                        a.v1.SetText(string(runes))
                        gofer.UpdateUI(_root, nil)
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

func validRune(r rune) bool {
    return true
}

func main() {
    gofer.NewApp().Run(&agent{})
}
