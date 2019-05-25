package gofer

import "fmt"

type BasicRect struct {
    x1, y1, x2, y2 int
}

func (r *BasicRect) String() string {
    return fmt.Sprintf("r(%v, %v)", r.From(), r.To())
}

func (r *BasicRect) SetFrom(from Point) {
    if from == nil {
        r.x1, r.y1 = 0, 0
        return
    }
    r.x1, r.y1 = from.X(), from.Y()
}

func (r *BasicRect) From() Point {
    return NewPoint(r.x1, r.y1)
}

func (r *BasicRect) SetTo(to Point) {
    if to == nil {
        r.x2, r.y2 = 0, 0
        return
    }
    r.x2, r.y2 = to.X(), to.Y()
}

func (r *BasicRect) To() Point {
    return NewPoint(r.x2, r.y2)
}

func (r *BasicRect) SetWidth(w int) {
    if w < 0 {
        w = 0
    }
    r.x2 = r.x1 + w - 1
}

func (r *BasicRect) Width() int {
    return r.x2 - r.x1 + 1
}

func (r *BasicRect) SetHeight(h int) {
    if h < 0 {
        h = 0
    }
    r.y2 = r.y1 + h - 1
}

func (r *BasicRect) Height() int {
    return r.y2 - r.y1 + 1
}

func (r *BasicRect) WhereIsPoint(p Point) RectPointLoc {
    if p == nil {
        return RectOutside
    }
    if p.X() < r.x1-1 || p.X() > r.x2+1 || p.Y() < r.y1-1 || p.Y() > r.y2+1 {
        return RectOutside
    }
    if p.X() < r.x1 || p.X() > r.x2 || p.Y() < r.y1 || p.Y() > r.y2 {
        return RectAdjacent
    }
    if p.X() == r.x1 || p.X() == r.x2 || p.Y() == r.y1 || p.Y() == r.y2 {
        return RectBorder
    }
    return RectInside
}

func (r *BasicRect) ContainPoint(p Point) bool {
    if p == nil {
        return false
    }
    return r.WhereIsPoint(p) == RectInside || r.WhereIsPoint(p) == RectBorder
}

func (r *BasicRect) IntersectWith(r1 Rect) bool {
    if r1 == nil {
        return false
    }
    f1, t1 := r.From(), r.To()
    f2, t2 := r1.From(), r1.To()

    x1 := f1.X()
    if f2.X() > x1 {
        x1 = f2.X()
    }
    y1 := f1.Y()
    if f2.Y() > y1 {
        y1 = f2.Y()
    }
    x2 := t1.X()
    if t2.X() < x2 {
        x2 = t2.X()
    }
    y2 := t1.Y()
    if t2.Y() < y2 {
        y2 = t2.Y()
    }

    return x1 <= x2 && y1 <= y2
}

func (r *BasicRect) Copy() Rect {
    return &BasicRect{r.x1, r.y1, r.x2, r.y2}
}

func (r *BasicRect) ChangeBaseByXY(x, y int) (rn Rect) {
    r.x1, r.y1, r.x2, r.y2 = r.x1-x, r.y1-y, r.x2-x, r.y2-y
    return r
}

func (r *BasicRect) ChangeBaseByPoint(base Point) (rn Rect) {
    if base != nil {
        r.ChangeBaseByXY(base.X(), base.Y())
    }
    return r
}

func BasicRectByXY(x1, y1, x2, y2 int) *BasicRect {
    return &BasicRect{x1, y1, x2, y2}
}

func BasicRectByPoint(from, to Point) *BasicRect {
    rect := &BasicRect{}
    if from != nil {
        rect.x1, rect.y1 = from.X(), from.Y()
    }
    if to != nil {
        rect.x2, rect.y2 = to.X(), to.Y()
    }
    return rect
}
