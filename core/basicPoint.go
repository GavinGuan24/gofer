package core


type basicPoint struct {
    x, y int
}

func (p *basicPoint) SetX(x int) {
    if p == nil {
        return
    }
    p.x = x
}

func (p *basicPoint) X() int {
    if p == nil {
        return 0
    }
    return p.x
}

func (p *basicPoint) SetY(y int) {
    if p == nil {
        return
    }
    p.y = y
}

func (p *basicPoint) Y() int {
    if p == nil {
        return 0
    }
    return p.y
}

func (p *basicPoint) Reverse() (pn Point) {
    if p == nil {
        return nil
    }
    return BasicPoint(-p.x, -p.y)
}

func (p *basicPoint) Copy() (pn Point) {
    return BasicPoint(p.x, p.y)
}

func (p *basicPoint) ChangeBaseByXY(x, y int) (pn Point) {
    p.x, p.y = p.x-x, p.y-y
    return p
}

func (p *basicPoint) ChangeBaseByPoint(base Point) (pn Point) {
    if base != nil {
        p.ChangeBaseByXY(base.X(), base.Y())
    }
    return p
}

func BasicPoint(x, y int) *basicPoint {
    return &basicPoint{x, y}
}
