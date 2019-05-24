package core

import (
    "fmt"
    "sync"
)

//更新UI的信息
type UpdateUiMsg struct {
    view View
    rect Rect
}

func (msg *UpdateUiMsg)GetView() View {
    return msg.view
}

func (msg *UpdateUiMsg)GetRect() Rect {
    return msg.rect
}

func NewUpdateUiMsg(view View, rect Rect) *UpdateUiMsg {
    return &UpdateUiMsg{view, rect}
}

///////////////////////////////////////////////////////////

// 向父视图发送更新UI的通知, Rect 坐标是基于当前视图的.
func UpdateUI(v View, rect Rect) {
    if v.SuperView() == nil {
        return
    }
    Receiver(v.SuperView()) <- NewUpdateUiMsg(v, rect)
}

// 父视图处理更新视图信息
func updateUiMsgHandler(v View, msg *UpdateUiMsg) {
    defer updateUiMsgErrorHandler()
    if kidRect := msg.GetRect(); kidRect == nil {
        //当接收到的子视图通知中的rect为空, 约定更新整个父视图(即当前视图)
        UpdateUI(v, NewRectByXY(0, 0, v.Width() - 1, v.Height() - 1))
    } else {
        if sv := msg.GetView(); sv != nil {
            //当前视图先将接收到的 Rect 转变到自身坐标轴, 然后向上层通知更新自身部分内容
            UpdateUI(v, kidRect.ChangeBaseByPoint(sv.Location().Reverse()))
        }
    }
}

// 视图更新失败处理
func updateUiMsgErrorHandler() {
    if o := recover(); o != nil {
        LogError(fmt.Sprintf("Unknown error when update UI: %v", o))
    }
}

// 获取当前视图内一个区域的内容, 并转化子视图内容至当前视图内容中
// 约定, 父视图承诺处理好子视图左右边界上可能出现的 2倍宽字符 越界问题
func GetMergeContent(v View, from Point, to Point) [][]Rune {
    basicV, ok := v.BaseView().(*basicView)
    if !ok {
        LogCrash("v cannot be a (*basicView)")
    }
    basicV.drawLock.Lock()
    defer basicV.drawLock.Unlock()

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
    //4. 计算需要的内容范围, 从子视图获取内容后, 迭代至现有内容中
    minInt, maxInt := func(a, b int) int {
        if a < b {
            return a
        }
        return b
    }, func(a, b int) int {
        if a > b {
            return a
        }
        return b
    }
    for _, sv := range svArr {
        //4.1 基于当前视图(父视图)坐标计算需要更新的范围, 并且变更坐标至子视图坐标系
        cRect := NewRectByXY(
            maxInt(from.X(), sv.Rect().From().X()),
            maxInt(from.Y(), sv.Rect().From().Y()),
            minInt(to.X(), sv.Rect().To().X()),
            minInt(to.Y(), sv.Rect().To().Y()))
        svCtRect := cRect.Copy().ChangeBaseByPoint(sv.Location())
        //4.2 迭代内容
        originContent = iteratingContent(
            GetMergeContent(sv, svCtRect.From(), svCtRect.To()), originContent,
            cRect, oRect.Copy(),
            svCtRect.From().X() == 0, svCtRect.To().X() == sv.Width()-1)
    }
    return originContent
}

// 迭代子视图的内容至现有视图内容中 (先处理左边界问题(2宽度留单的处理为符号 ›), 再使用子视图内容复写, 同时处理右边界问题)
//
// 以下参数均参考当前视图(父视图)的坐标系.
// svContent: 子视图内容
// oContent: 父视图内容
// cRect: svContent 的作用域
// oRect: oContent 的作用域
// watchLeft: 注意内容的左边界问题
// watchRight: 注意内容的右边界问题
func iteratingContent(svContent, oContent [][]Rune, cRect, oRect Rect, watchLeft, watchRight bool) [][]Rune {

    //将当前问题转化为相对做坐标的问题, 计算两个二维数组下标差值
    dCol, dRow := oRect.From().X()-cRect.From().X(), oRect.From().Y()-cRect.From().Y()
    cw, ch := cRect.Width(), cRect.Height()
    for row := 0; row < ch; row++ {
        for col := 0; col < cw; col++ {
            cRune := svContent[row][col]
            oRow, oCol := row-dRow, col-dCol
            if oRow < 0 || oCol < 0 {
                continue
            }
            if watchLeft && col == 0 && oCol-1 >= 0 {
                //处理左边界
                oRune_l1 := oContent[oRow][oCol-1]
                if oRune_l1.Width() == 2 {
                    oRune_l1.SetMainc('›')
                    oRune_l1.SetCombc(nil)
                }
            }
            if watchRight && col == cw-1 {
                //处理右边界
                if cRune.Width() == 2 {
                    cRune.SetMainc('›')
                    cRune.SetCombc(nil)
                }
            }
            // 内容复写
            oContent[oRow][oCol] = cRune
        }
    }
    return oContent
}


var receiverLock sync.Mutex
// 获取一个接收器. 接收器用于接收子视图的通知.
// 通知内容为一个 message.UpdateUiMsg, 该 msg的rect坐标轴为子视图的坐标轴.
// if message.rect == nil, 当前视图(父视图)向上层通知(当前视图的父视图/root view的上层为screen)更新自身所有内容
// 否则, 当前视图先将接收到的 Rect 转变到自身坐标轴, 然后向上层通知更新自身部分内容
func Receiver(v View) chan<- *UpdateUiMsg {
    receiverLock.Lock()
    defer receiverLock.Unlock()
    bv, ok := v.BaseView().(*basicView)
    if !ok {
        LogCrash("v cannot be a (*basicView)")
    }
    if bv.uiMsgQueue == nil {
        //赋予当前视图成为视图通知链一个节点的能力
        bv.uiMsgQueue = make(chan *UpdateUiMsg)
        go func() {
            for {
                updateUiMsgHandler(v, <-bv.uiMsgQueue)
            }
        }()
    }
    return bv.uiMsgQueue
}
