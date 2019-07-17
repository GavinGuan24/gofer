package gofer

import (
    "github.com/gdamore/tcell"
)

func init() {
    LogInfo("Package init: gofer.rootView")
}

// rootView 是所有视图链上的根节点, (与 app 协作)直接与 screen 交互.
// 行为上不一定与 <code>View interface</code> 一致.
// 以源码为准.
type rootView struct {
    *basicView
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

func (v *rootView) BaseView() View {
    return v.basicView
}

func (v *rootView) GetContent(from Point, to Point) [][]Rune {
    return v.basicView.GetContent(from, to)
}

func (v *rootView) UpdateUI(rect Rect) {
    Receiver(v) <- newUpdateUiMsg(nil, nil)
}

func (v *rootView) updateUiMsgHandler(msg *UpdateUiMsg) {
    defer updateUiMsgErrorHandler()
    if v.Screen == nil {
        return
    }
    if msg.GetRect() == nil {
        // root view 子身回调 或者 来自子视图的消息, 更新自身全部
        sw, sh := v.Screen.Size()
        content := GetMergeContent(v, NewPoint(0, 0), NewPoint(sw-1, sh-1))

        //处理二倍宽字符在更新过程中导致的异常空白区, 该bug是tcell本身的问题, 如果官方不修复, 就只能使用下面的代码段折中处理.
        for row := 0; row < sh; row++ {
            _, _, _, width := v.Screen.GetContent(sw-2, row)
            if width > 1 {
                v.Screen.Clear()
                v.Screen.Show()
            }
        }

        for row, line := range content {
            for col, cRune := range line {
                v.Screen.SetContent(col, row, cRune.Mainc(), cRune.Combc(), cRune.Style())
            }
        }
        v.Screen.Show()
        return
    }

    if sv := msg.GetView(); sv != nil {
        //按子视图给的范围, 从 root view 的角度获取内容, 并更新
        cRect := msg.GetRect().ChangeBaseByPoint(sv.Location().Reverse())
        cw, ch := cRect.Width(), cRect.Height()
        from, to := cRect.From(), cRect.To()
        content := GetMergeContent(v, from, to)
        for row := 0; row < ch; row++ {
            for col := 0; col < cw; col++ {
                cRune := content[row][col]
                v.Screen.SetContent(col+from.X(), row+from.Y(), cRune.Mainc(), cRune.Combc(), cRune.Style())
            }
        }
        v.Screen.Show()
        return
    }
}

func NewRootView() *rootView {
    root := &rootView{}
    root.basicView = &basicView{subviews: make([]View, 0, 4), rect: NewRectByXY(0, 0, 0, 0)}
    root.tagMap = make(map[string]View)
    //赋予root view 成为视图通知链根节点的能力
    root.uiMsgQueue = make(chan *UpdateUiMsg)
    go func() {
        for {
            root.updateUiMsgHandler(<-root.uiMsgQueue)
        }
    }()
    return root
}
