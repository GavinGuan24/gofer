package core

type Point interface {
    SetX(x int)
    X() int
    SetY(y int)
    Y() int

    // 把当前点视为向量终点, (0, 0)视为起点, 返回反向向量的终点
    // 例: 当前为(3, 4), 该方法返回 (-3, -4)
    Reverse() (pn Point)
    // 复制一个自己
    Copy() (pn Point)

    // 使用 base 作为新的坐标轴原点, 更新自身坐标值.
    // 假设运算符号为 $, 示例如下:
    // p_old(3,4) $ base(-1, -2) = p_new(4, 6)
    ChangeBaseByXY(x, y int) (pn Point)
    ChangeBaseByPoint(base Point) (pn Point)
}

func NewPoint(x, y int) Point {
    return BasicPoint(x, y)
}
