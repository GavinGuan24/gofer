package core

// 矩形Rect 与 点Point 的位置关系
type RectPointLoc int

const (
    // 在矩形内, 并不是矩形边界
    RectInside RectPointLoc = iota
    // 在矩形内, 是矩形边界
    RectBorder
    // 在矩形外, 与其边界相邻
    RectAdjacent
    // 在矩形外, 不与其边界相邻
    RectOutside

    // 这里是图示, 矩形是0与1(含1). 当矩形长宽其中一个等于 1, 位置关系就没有所谓的 RectInside
    // 3333333
    // 3222223
    // 3211123
    // 3210123
    // 3211123
    // 3222223
    // 3333333
)

type Rect interface {
    SetFrom(from Point)
    From() Point
    SetTo(to Point)
    To() Point
    SetWidth(w int)
    Width() int
    SetHeight(h int)
    Height() int
    // 获取Point相对于矩形的位置描述. 如果Point是nil, 永远处理为 RectOutside
    WhereIsPoint(p Point) RectPointLoc
    // 判定点是否被矩形包含
    ContainPoint(p Point) bool
    // 是否与另一个矩形相交
    IntersectWith(r1 Rect) bool

    // 复制一个自己
    Copy() Rect

    // 使用 base 作为新的坐标轴原点, 更新自身坐标值.
    // 可参考 views.Point.ChangeBaseByXXX()
    ChangeBaseByXY(x, y int) (rn Rect)
    ChangeBaseByPoint(base Point) (rn Rect)
}

func NewRectByXY(x1, y1, x2, y2 int) Rect {
    return BasicRectByXY(x1, y1, x2, y2)
}

func NewRectByPoint(from, to Point) Rect {
    return BasicRectByPoint(from, to)
}
