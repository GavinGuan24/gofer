package views

import (
    . "github.com/GavinGuan24/gofer/views/code"
    "github.com/gdamore/tcell"
)

// View 是对视图的抽象, 这里的默认实现是 <code>type basicView struct</code>
type View interface {

    // 增减子视图
    AddSubview(subview View) (ok bool)
    RemoveSubview(subview View) (ok bool)
    AddSubviewWithTag(subview View, tag string) (ok bool)
    RemoveSubviewWithTag(tag string) (ok bool)
    // 获取所有子视图
    GetSubviews() []View
    // 查找不到时, 返回 nil
    GetSubviewWithTag(tag string) View

    // 获取父视图
    SuperView() View

    // 视图宽高与其在父视图中的坐标
    SetLocation(loc Point)
    Location() Point
    SetWidth(w int)
    SetHeight(h int)
    Width() int
    Height() int

    // Z 是视图优先级, 默认为0
    // 约定父视图处理子视图优先级时, 优先级越低的子视图越容易被其他子视图遮挡.
    // 当多个子视图优先级相等时, 越后面被添加进父视图的, 优先级越高
    SetZ(z int)
    Z() int

    // 当前视图的默认风格
    SetStyle(style tcell.Style)
    Style() tcell.Style

    // 获取当前视图的内容. from, to 是当前视图内的两点, 不会为 nil
    //
    // 如果想自行实现一个特定功能的视图, 实现该方法即可控制自定义视图的内容.
    // 约定, 父视图承诺处理好子视图左右边界上可能出现的 2倍宽字符 越界问题.
    // 该问题会受到影响(以何种方式处理子视图优先级)
    //
    // 这里使用二维数组作为出参, 就算外面要求给出部分视图内容, 数组元素(0,0)也一定是from这个点的数据, 约定上层会自行转化坐标, 所以(0,0)不一定是当前视图的左上角
    GetContent(from Point, to Point) [][]*Rune

    // 向父视图发送更新UI的通知, Rect 坐标是基于当前视图的. 详情参考方法<code>basicView.receiver()</code>
    UpdateUI(rect Rect)
}

func NewView() View {
    return BasicView()
}