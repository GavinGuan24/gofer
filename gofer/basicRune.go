package gofer

import (
    "bytes"
    "fmt"
    "github.com/gdamore/tcell"
    "github.com/mattn/go-runewidth"
)

// 某视图中字符位的描述, 位置由使用 basicRune 的程序维护
// combc 一般为 nil, 除非要处理 ZWJ(zero-width joiner) 问题
type basicRune struct {
    mainc rune
    combc []rune
    style tcell.Style
}

func (r *basicRune) SetMainc(mainc rune) {
    if r == nil {
        return
    }
    r.mainc = mainc
}
func (r *basicRune) Mainc() rune {
    if r == nil {
        return Empty
    }
    return r.mainc
}
func (r *basicRune) SetCombc(combc []rune) {
    if r == nil {
        return
    }
    r.combc = combc
}
func (r *basicRune) Combc() []rune {
    if r == nil {
        return nil
    }
    return r.combc
}
func (r *basicRune) SetStyle(style tcell.Style) {
    if r == nil {
        return
    }
    r.style = style
}
func (r *basicRune) Style() tcell.Style {
    if r == nil {
        return tcell.StyleDefault
    }
    return r.style
}

func (r *basicRune) String() string {
    if r == nil {
        return ""
    }
    if r.combc == nil || len(r.combc) == 0 {
        return fmt.Sprintf("%c", r.mainc)
    } else {
        var buf bytes.Buffer
        buf.WriteRune(r.mainc)
        for _, t := range r.combc {
            buf.WriteRune(t)
        }
        return fmt.Sprintf("%s", buf.String())
    }
}

func (r *basicRune) Width() int {
    if r.combc == nil || len(r.combc) == 0 {
        return runewidth.RuneWidth(r.mainc)
    }
    for _, t := range r.combc {
        if t == 0x200D {
            return 2
        }
        if runewidth.RuneWidth(t) > 1 {
            return 2
        }
    }
    return 1
}

func BasicRune(mainc rune, combc []rune, style tcell.Style) Rune {
    return &basicRune{mainc, combc, style}
}
