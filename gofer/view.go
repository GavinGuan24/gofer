package gofer

import (
    "fmt"
    "github.com/gdamore/tcell"
)

// View 是对 [视图] 的抽象, 这里的默认实现是 <code>type basicView struct</code>
type View interface {

    fmt.Stringer

    // 增减子视图, 因为golang 是 内嵌型"继承", 所以需要真实的父类传入
    AddSubview(subview View, realSuper View) (ok bool)
    RemoveSubview(subview View) (ok bool)
    AddSubviewWithTag(subview View, tag string, realSuper View) (ok bool)
    RemoveSubviewWithTag(tag string) (ok bool)
    // 获取所有子视图
    GetSubviews() []View
    // 查找不到时, 返回 nil
    GetSubviewWithTag(tag string) View

    // 父视图, 不要随意调用 SetSuperView, 除非你知道自己在做什么
    SetSuperView(view View)
    SuperView() View

    // 获取视图的矩阵信息
    Rect() Rect

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

    // 返回实现该接口的基类
    // 因为golang的继承是内嵌, 不是真的继承, 所以约定这里返回实现<code>View interface</code>的基类
    // 即, 所有"子类"均返回其内嵌("继承")的实现该接口的基类.
    // 比如该接口的默认实现是 *basicView, 那么所有 basicView 的子类均返回内嵌的 *basicView
    BaseView() View

    // 获取当前视图的内容. from, to 是当前视图内的两点, 不会为 nil
    //
    // 如果想自行实现一个特定功能的视图, 实现该方法即可控制自定义视图的内容.
    // 如果是展示文本的视图, 请实现者自行解决以下问题, 父视图会强制处理 子视图左右边界上可能出现的 2倍宽字符 越界问题.
    //      1. 中文等 2倍宽字符
    //      2. ZWJ(zero-width joiner) 问题
    //
    // 约定, 父视图承诺处理好子视图左右边界上可能出现的 2倍宽字符 越界问题.
    // 该问题会受到影响(以何种方式处理子视图优先级)
    //
    // 这里使用二维数组作为出参, 就算 [调用方] 要求给出部分视图内容, 数组元素(0,0)也一定是 [from] 这个点的数据, 约定 [调用方] 会自行转化坐标, 所以(0,0)不一定是当前视图的左上角
    GetContent(from Point, to Point) [][]Rune
}

func NewView() View {
    return NewBasicView()
}