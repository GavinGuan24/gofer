package core

import (
    "fmt"
    "github.com/GavinGuan24/gofer/log"
    "github.com/gdamore/tcell"
    "sync"
)

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
    updateUiMsgQueue chan *UpdateUiMsg

    // 视图基本描述
    rect      Rect
    z         int
    style     tcell.Style
    superView View
    subviews  []View
}

func (v *basicView) AddSubview(subview View) (ok bool) {
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
    subview.SetSuperView(v)
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

func (v *basicView) AddSubviewWithTag(subview View, tag string) (ok bool) {
    v.tagLock.Lock()
    defer v.tagLock.Unlock()
    if subview == nil || v.tagMap[tag] != nil {
        return false
    }
    if v.AddSubview(subview) {
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
    v.rect.SetFrom(loc)
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

// 获取当前视图内一个区域的内容, 并转化子视图内容至当前视图内容中
// 约定, 父视图承诺处理好子视图左右边界上可能出现的 2倍宽字符 越界问题
func (v *basicView) GetMergeContent(from Point, to Point) [][]Rune {
    v.drawLock.Lock()
    defer v.drawLock.Unlock()
    //1. 调用 GetContent 获取当前视图原始内容的 originContent []basicRune
    originContent := v.GetContent(from, to)
    //2. 遍历子视图, 找出范围内的子视图
    svArr := make([]View, 0, 4)
    oRect := NewRectByPoint(from, to)
    for _, sv := range v.GetSubviews() {
        if sv.Rect().IntersectWith(oRect) {
            svArr = append(svArr, sv)
        }
    }
    //3. 按优先级排序范围内的子视图
    if len(svArr) > 1 {
        for i := 0; i < len(svArr)-1; i++ {
            for j := i + 1; j < len(svArr); j++ {
                if svArr[i].Z() > svArr[j].Z() {
                    svArr[i], svArr[j] = svArr[j], svArr[i]
                }
            }
        }
    }
    //4. 根据子视图位置计算需要的内容范围(坐标基于子视图), 获取内容后, 迭代至现有内容中
    for _, sv := range svArr {
        var svCtRect Rect
        if sv.Rect().ContainPoint(from) {
            if sv.Rect().ContainPoint(to) {
                svCtRect = NewRectByPoint(from, to).ChangeBaseByPoint(sv.Location())
            } else {
                svCtRect = NewRectByPoint(from.ChangeBaseByPoint(sv.Location()), sv.Rect().To())
            }
        } else {
            if sv.Rect().ContainPoint(to) {
                svCtRect = NewRectByPoint(sv.Rect().From(), to.ChangeBaseByPoint(sv.Location()))
            } else {
                svCtRect = sv.Rect().Copy()
            }
        }


        originContent = v.iteratingContent(
            sv.GetMergeContent(svCtRect.From(), svCtRect.To()), originContent,
            svCtRect.ChangeBaseByPoint(sv.Location().Copy().Reverse()), oRect.Copy(),
            svCtRect.From().X() == 0, svCtRect.To().X() == sv.Width()-1)
    }
    return originContent
}

// 迭代子视图的内容至现有视图内容中 (先处理左边界问题(2宽度留单的处理为符号 …), 再使用子视图内容复写, 同时处理右边界问题)
//
// 以下参数均参考当前视图(父视图)的坐标系.
// svContent: 子视图内容
// oContent: 父视图内容
// cRect: svContent 的作用域
// oRect: oContent 的作用域
// watchLeft: 注意内容的左边界问题
// watchRight: 注意内容的右边界问题
func (v *basicView) iteratingContent(svContent, oContent [][]Rune, cRect, oRect Rect, watchLeft, watchRight bool) [][]Rune {

    //将当前问题转化为相对做坐标的问题, 计算两个二维数组下标差值
    dCol, dRow := oRect.From().X()-cRect.From().X(), oRect.From().Y()-cRect.From().Y()
    cw, ch := cRect.Width(), cRect.Height()
    for row := 0; row < ch; row++ {
        for col := 0; col < cw; col++ {
            cRune := svContent[row][col]
            oRow, oCol := row-dRow, col-dCol
            if watchLeft && col == 0 {
                //处理左边界
                oRune_l1 := oContent[oRow][oCol-1]
                if oRune_l1.Width() == 2 {
                    oRune_l1.SetMainc('…')
                    oRune_l1.SetCombc(nil)
                }
            }
            if watchRight && col == cw-1 {
                //处理右边界
                if cRune.Width() == 2 {
                    cRune.SetMainc('…')
                    cRune.SetCombc(nil)
                }
            }
            // 内容复写
            oContent[oRow][oCol] = cRune
        }
    }
    return oContent
}

// 提供一个接收器. 接收器用于接收子视图的通知.
// 通知内容为一个 message.UpdateUiMsg, 该 msg的rect坐标轴为子视图的坐标轴.
// if message.rect == nil, 当前视图(父视图)向上层通知(当前视图的父视图/root view的上层为screen)更新自身所有内容
// 否则, 当前视图先将接收到的 Rect 转变到自身坐标轴, 然后向上层通知更新自身部分内容
func (v *basicView) Receiver() chan<- *UpdateUiMsg {
    return v.updateUiMsgQueue
}

func (v *basicView) UpdateUI(rect Rect) {
    if v.superView == nil {
        return
    }
    v.superView.Receiver() <- NewUpdateUiMsg(v, rect)
}

//----------

func (v *basicView) updateUiMsgHandler(msg *UpdateUiMsg) {
    defer func() {
        if o := recover(); o != nil {
            log.Logger(fmt.Sprintf("Unknown error when update UI: %v\n", o))
        }
    }()
    if kidRect := msg.GetRect(); kidRect == nil {
        //当接收到的子视图通知中的rect为空, 约定更新整个父视图(即当前视图)
        v.UpdateUI(v.rect)
    } else {
        if sv := msg.GetView(); sv != nil {
            //当前视图先将接收到的 Rect 转变到自身坐标轴, 然后向上层通知更新自身部分内容
            v.UpdateUI(kidRect.ChangeBaseByPoint(sv.Location().Reverse()))
        }
    }
}

// 基于当前视图坐标系的点, 是否在矩形外
func (v *basicView) innerPointOutOfRect(from, to, p Point) bool {
    return p.X() < from.X() || p.X() > to.X() || p.Y() < from.Y() || p.Y() > to.Y()
}

//----------

func BasicView() *basicView {
    bv := &basicView{subviews: make([]View, 0, 4), rect: NewRectByXY(0, 0, 0, 0)}
    bv.tagMap = make(map[string]View)
    //赋予当前视图成为视图通知链一个节点的能力
    bv.updateUiMsgQueue = make(chan *UpdateUiMsg)
    go func() {
        for {
            bv.updateUiMsgHandler(<-bv.updateUiMsgQueue)
        }
    }()
    return bv
}
