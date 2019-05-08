package core

import (
    "github.com/GavinGuan24/gofer/log"
    "github.com/gdamore/tcell"
)

func init() {
    log.Info("package init: core.rootView")
}

// rootView 是所有视图链上的根节点, (与 app 协作)直接与 screen 交互.
// 行为上不一定与 <code>View interface</code> 一致.
// 以源码为准.
type rootView struct {
    basicView
    tcell.Screen
}

func (v *rootView) SuperView() View {
    return nil
}

func (v *rootView) Location() Point {
    return NewPoint(0, 0)
}

func (v *rootView) Width() int {
    w, _ := v.Screen.Size()
    return w
}

func (v *rootView) Height() int {
    _, h := v.Screen.Size()
    return h
}

func (v *rootView) Z() int {
    return 0
}

func (v *rootView) SetStyle(style tcell.Style) {
    v.basicView.SetStyle(style)
}

func (v *rootView) Style() tcell.Style {
    return v.basicView.Style()
}

func (v *rootView) GetContent(from Point, to Point) [][]Rune {
    return v.basicView.GetContent(from, to)
}

func (v *rootView) UpdateUI(rect Rect) {
    v.Receiver() <- NewUpdateUiMsg(nil, nil)
}

func (v *rootView) updateUiMsgHandler(msg *UpdateUiMsg) {
    defer v.updateUiMsgErrorHandler()
    
    if v.Screen == nil {
        return
    }
    log.Debug("ggg")
    if msg.GetRect() == nil {
        // root view 子身回调 或者 来自子视图的消息, 更新自身全部
        sw, sh := v.Screen.Size()
        content := v.GetMergeContent(NewPoint(0, 0), NewPoint(sw-1, sh-1))
        v.Screen.Clear()
        for row, line := range content {
            for col, cRune := range line {
                v.Screen.SetContent(col, row, cRune.Mainc(), cRune.Combc(), cRune.Style())
            }
        }
        v.Screen.Show()
        return
    }

    if sv := msg.GetView(); sv != nil {
        //转化子视图的内容至root view 坐标, 更新至screen
        cRect := msg.GetRect()
        cw, ch := cRect.Width(), cRect.Height()
        //计算坐标差
        dCol, dRow:= -cRect.From().X() - sv.Location().X(), -cRect.From().Y() - sv.Location().Y()
        //更新内容
        content := sv.GetMergeContent(cRect.From(), cRect.To())
        for row := 0; row < ch; row++ {
            for col := 0; col < cw; col++ {
                cRune := content[row][col]
                v.Screen.SetContent(col-dCol, row-dRow, cRune.Mainc(), cRune.Combc(), cRune.Style())
            }
        }
        v.Screen.Show()
        return
    }
}

func RootView() *rootView {
    root := &rootView{}
    root.basicView = basicView{subviews: make([]View, 0, 4), rect: NewRectByXY(0, 0, 0, 0)}
    root.tagMap = make(map[string]View)
    //赋予root view 成为视图通知链根节点的能力
    root.updateUiMsgQueue = make(chan *UpdateUiMsg)
    go func() {
        for {
            root.updateUiMsgHandler(<-root.updateUiMsgQueue)
        }
    }()
    return root
}
