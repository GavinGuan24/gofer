package gofer

import (
    "fmt"
    "github.com/gdamore/tcell"
    "sync"
)

func init() {
    LogInfo("Package init: gofer.basicView")
}

type basicView struct {
    // 保护 GetContent()
    drawLock sync.Mutex

    // 作用于 subviews 读写
    subviewsLock sync.Mutex

    // 作用于 tag - subviews 读写
    tagLock sync.Mutex
    // tagged subviews
    tagMap map[string]View

    //接收子视图通知
    uiMsgQueue chan *UpdateUiMsg

    // 视图基本描述
    rect      Rect
    z         int
    style     tcell.Style
    superView View
    subviews  []View
}

func (v *basicView) String() string {
    return fmt.Sprintf("v%p(%v)", v, v.rect)
}

func (v *basicView) AddSubview(subview View, realSuper View) (ok bool) {
    v.drawLock.Lock()
    defer v.drawLock.Unlock()
    v.subviewsLock.Lock()
    defer v.subviewsLock.Unlock()
    if subview == nil {
        return false
    }
    for _, sv := range v.subviews {
        if sv == subview {
            // 原则: 禁止重复添加子视图, 判定操作失败
            return false
        }
    }
    v.subviews = append(v.subviews, subview)
    subview.SetSuperView(realSuper)
    return true
}

func (v *basicView) RemoveSubview(subview View) (ok bool) {
    v.drawLock.Lock()
    defer v.drawLock.Unlock()
    v.subviewsLock.Lock()
    defer v.subviewsLock.Unlock()
    if subview == nil || len(v.subviews) == 0 {
        return false
    }
    found := false
    newSubviews := make([]View, 0, len(v.subviews))
    for _, sv := range v.subviews {
        if !found && sv == subview {
            found = true
            subview.SetSuperView(nil)
            continue
        }
        newSubviews = append(newSubviews, sv)
    }
    if !found {
        return false
    } else {
        v.subviews = newSubviews
        return true
    }
}

func (v *basicView) AddSubviewWithTag(subview View, tag string, realSuper View) (ok bool) {
    v.tagLock.Lock()
    defer v.tagLock.Unlock()
    if subview == nil || v.tagMap[tag] != nil {
        return false
    }
    if v.AddSubview(subview, realSuper) {
        v.tagMap[tag] = subview
        return true
    }
    return false
}

func (v *basicView) RemoveSubviewWithTag(tag string) (ok bool) {
    v.tagLock.Lock()
    defer v.tagLock.Unlock()
    if subview := v.tagMap[tag]; subview != nil {
        v.tagMap[tag] = nil
        return v.RemoveSubview(subview)
    }
    return false
}

func (v *basicView) GetSubviews() []View {
    v.subviewsLock.Lock()
    defer v.subviewsLock.Unlock()
    subviewSli := make([]View, 0, len(v.subviews))
    for _, sv := range v.subviews {
        subviewSli = append(subviewSli, sv)
    }
    return subviewSli
}

func (v *basicView) GetSubviewWithTag(tag string) View {
    v.tagLock.Lock()
    defer v.tagLock.Unlock()
    return v.tagMap[tag]
}

func (v *basicView) SetSuperView(view View) {
    v.superView = view
}

func (v *basicView) SuperView() View {
    return v.superView
}

func (v *basicView) Rect() Rect {
    return v.rect.Copy()
}

func (v *basicView) SetLocation(loc Point) {
    if loc == nil {
        loc = NewPoint(0, 0)
    }
    w, h := v.rect.Width(), v.rect.Height()
    v.rect.SetFrom(loc)
    v.rect.SetWidth(w)
    v.rect.SetHeight(h)
}

func (v *basicView) Location() Point {
    return v.rect.From()
}

func (v *basicView) SetWidth(w int) {
    v.rect.SetWidth(w)
}

func (v *basicView) SetHeight(h int) {
    v.rect.SetHeight(h)
}

func (v *basicView) Width() int {
    return v.rect.Width()
}

func (v *basicView) Height() int {
    return v.rect.Height()
}

func (v *basicView) SetZ(z int) {
    v.z = z
}

func (v *basicView) Z() int {
    return v.z
}

func (v *basicView) SetStyle(style tcell.Style) {
    v.style = style
}

func (v *basicView) Style() tcell.Style {
    return v.style
}

func (v *basicView) BaseView() View {
    return v
}

func (v *basicView) GetContent(from Point, to Point) [][]Rune {
    //就算调用方要求给出部分视图内容, 数组元素(0,0)也一定是from这个点的数据, 约定调用方会自行转化坐标
    x1, y1, x2, y2 := from.X(), from.Y(), to.X(), to.Y()
    h := y2 - y1 + 1
    w := x2 - x1 + 1
    lines := make([][]Rune, 0, h)
    var words []Rune
    for row := 0; row < h; row++ {
        for col := 0; col < w; col++ {
            if col == 0 {
                words = make([]Rune, 0, w)
            }
            words = append(words, BasicRune(' ', nil, v.Style()))
            if col == w-1 {
                lines = append(lines, words)
            }
        }
    }
    return lines
}

//----------

func NewBasicView() *basicView {
    bv := &basicView{subviews: make([]View, 0, 4), rect: NewRectByXY(0, 0, 0, 0)}
    bv.tagMap = make(map[string]View)
    return bv
}
