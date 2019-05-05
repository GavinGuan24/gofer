package code

import (
    "fmt"
    "github.com/GavinGuan24/gofer/log"
    "github.com/GavinGuan24/gofer/message"
    . "github.com/GavinGuan24/gofer/views"
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
    updateUiMsgQueue chan *message.UpdateUiMsg

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
    var realSubview = subview.(*basicView)
    realSubview.setSuperView(v)
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
            var realSubview = subview.(*basicView)
            realSubview.setSuperView(nil)
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

func (v *basicView) setSuperView(view View) {
    v.superView = view
}

func (v *basicView) SuperView() View {
    return v.superView
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

func (v *basicView) GetContent(from Point, to Point) [][]*Rune {
    //TODO: 这里使用数组, 就算外面要求给出部分视图内容, 也只能用二维数据返回, 也就是说, 外界需要自行转化坐标, 数组元素0, 0 不一定是当前视图的左上角, 但一定是from这个点
    x1, y1, x2, y2 := from.X(), from.Y(), to.X(), to.Y()
    lines := make([][]*Rune, 0, y2-y1+1)
    var words []*Rune
    for row := y1; row <= y2; row++ {
        for col := x1; col <= x2; col++ {
            if col == x1 {
                words = make([]*Rune, 0, x2-x1+1)
            }
            words = append(words, NewRune(' ', nil, v.Style()))
            if col == x2 {
                lines = append(lines, words)
            }
        }
    }
}

// 获取当前视图内一个区域的内容.
//
// 如果是展示文本的视图, 请实现者自行解决以下问题
// 1. 中文等 2倍宽字符
// 2. ZWJ(zero-width joiner) 问题
//
// 约定, 父视图承诺处理好子视图左右边界上可能出现的 2倍宽字符 越界问题
// 该问题会受到 以何种方式处理子视图优先级的 影响
func (v *basicView) getContent(from Point, to Point) [][]*Rune {
    v.drawLock.Lock()
    defer v.drawLock.Unlock()
    if from == nil || !v.rect.ContainPoint(from) {
        from = v.rect.From()
    }
    if to == nil || !v.rect.ContainPoint(to) {
        to = v.rect.To()
    }



    //v.GetSubviews()

    //TODO: 遍历子视图


    //      0. 调用 GetContent 获取当前视图空白内容的 originContent []Rune
    //originContent := v.GetContent(from, to)

    //      1. 找出范围内的
    //      2. 且按优先级排序 ==> drawViewArr

    /*

func main() {
    arr := []string{"a1", "a2", "b0", "c1", "a3", "c2"}

    for i := 0; i < len(arr) - 1; i++ {
        for j := i + 1; j < len(arr); j++ {
            if cmp(arr[i], arr[j]) == 1 || cmp(arr[i], arr[j]) == 0 {
                arr[i], arr[j] = arr[j], arr[i]
            }
        }
    }

    fmt.Println(arr)
}

func cmp(s1, s2 string) int {
    if s1[0:1] < s2[0:1] {
        return -1
    }
    if s1[0:1] > s2[0:1] {
        return 1
    }
    if s1[1:2] < s2[1:2] {
        return -1
    }
    if s1[1:2] > s2[1:2] {
        return 1
    }
    return 0
}

*/

    //TODO: 遍历 drawViewArr
    //      1. 换算子视图坐标轴的 (from, to) call sv.GetContent()
    //      2. 迭代 drawRuneArr, 先处理左边界问题(2宽度留单的处理为 …), 再使用子视图内容复写, 同时处理右边界问题

    return lines
}

// 提供一个接收器. 接收器用于接收子视图的通知.
// 通知内容为一个 message.UpdateUiMsg, 该 msg的rect坐标轴为子视图的坐标轴.
// if message.rect == nil, 当前视图(父视图)向上层通知(当前视图的父视图/root view的上层为screen)更新自身所有内容
// 否则, 当前视图先将接收到的 Rect 转变到自身坐标轴, 然后向上层通知更新自身部分内容
func (v *basicView) receiver() chan<- *message.UpdateUiMsg {
    return v.updateUiMsgQueue
}

func (v *basicView) UpdateUI(rect Rect) {
    if v.superView == nil {
        return
    }
    var realSuperView = v.superView.(*basicView)
    realSuperView.receiver() <- message.NewUpdateUiMsg(v, rect)
}

//----------

func (v *basicView) updateUiMsgHandler(msg *message.UpdateUiMsg) {
    defer func() {
        log.Logger(fmt.Sprintf("Unknown error when update UI: %v\n", recover()))
    }()
    if kidRect := msg.GetRect(); kidRect == nil {
        v.UpdateUI(v.rect)
    } else {
        //TODO: 当前视图先将接收到的 Rect 转变到自身坐标轴, 然后向上层通知更新自身部分内容
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
    bv.updateUiMsgQueue = make(chan *message.UpdateUiMsg)
    go func() {
        for {
            bv.updateUiMsgHandler(<-bv.updateUiMsgQueue)
        }
    }()
    return bv
}




