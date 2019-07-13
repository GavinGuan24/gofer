package main

import (
    "fmt"
    "github.com/GavinGuan24/gofer/gofer"
    "github.com/gdamore/tcell"
)

///////////////////////
type agent struct {
    v0 *View0
    v1 *View1
    v2 *View1
}

func (a *agent) Launched(root gofer.RootView) {
    v0 := NewView0()
    a.v0 = v0
    v0.SetLocation(gofer.NewPoint(1, 1))
    v0.SetWidth(18)
    v0.SetHeight(20)
    v0.SetStyle(tcell.StyleDefault.Background(tcell.NewRGBColor(225,225,225)))
    root.AddSubview(v0, root)

    v1 := NewView1()
    a.v1 = v1
    v1.SetLocation(gofer.NewPoint(0,0))
    v1.SetWidth(10)
    v1.SetHeight(2)
    v1.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkGray).Background(tcell.ColorAntiqueWhite))
    v1.SetZ(2)
    v0.AddSubview(v1, v0)

    v2 := NewView1()
    a.v2 = v2
    v2.SetLocation(gofer.NewPoint(0,0))
    v2.SetWidth(7)
    v2.SetHeight(2)
    v2.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkGray).Background(tcell.ColorGhostWhite))
    v2.SetZ(3)
    root.AddSubview(v2, root)

    //go func() {
    //   for {
    //       r := (int32)(rand.Intn(256))
    //       g := (int32)(rand.Intn(256))
    //       b := (int32)(rand.Intn(256))
    //       v0.SetStyle(tcell.StyleDefault.Background(tcell.NewRGBColor(225,225,225)).Foreground(tcell.NewRGBColor(r, g, b)))
    //       gofer.UpdateUI(v0, nil)
    //       time.Sleep(time.Second)
    //   }
    //}()
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

///////////////////////
type View0 struct {
    gofer.View
}

func (v *View0) GetContent(from gofer.Point, to gofer.Point) [][]gofer.Rune {
   content := v.View.GetContent(from, to)
    for _, line := range content {
        for col, ch := range line {
            if from.X() % 2 == 0 && col % 2 == 0 {
                ch.SetMainc('我')
            }
            if from.X() % 2 != 0 && col % 2 != 0 {
                ch.SetMainc('我')
            }
        }
    }
   return content
}

func NewView0() *View0 {
    view0 := &View0{}
    view0.View = gofer.NewView()
    return view0
}

////////////////////////
type View1 struct {
    gofer.View
}

func (v *View1) GetContent(from gofer.Point, to gofer.Point) [][]gofer.Rune {
    content := v.View.GetContent(from, to)

    
    for _, line := range content {
        for col, ch := range line {
            if from.X() % 2 == 0 && col % 2 == 0 {
                ch.SetMainc('你')
            }
            if from.X() % 2 != 0 && col % 2 != 0 {
                ch.SetMainc('你')
            }
        }
    }
    return content
}

func NewView1() *View1 {
    view1 := &View1{}
    view1.View = gofer.NewView()
    return view1
}

///////////////////////
func main() {
    gofer.NewApp().Run(&agent{})
}
