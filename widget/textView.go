package widget

import (
    "bytes"
    "github.com/GavinGuan24/gofer/gofer"
    "strings"
)

// 单行的文本视图, 会忽略 "\n", "\r". 建议展示较短的单行文本
type TextView struct {
    gofer.View
    text []gofer.Rune
}

func (v *TextView) GetContent(from gofer.Point, to gofer.Point) [][]gofer.Rune {
    content := v.View.GetContent(from, to)
    //
    flag, step := 0, 0
    for _, tRune := range v.text {
        tRune.SetStyle(v.Style())
        runeWidth := tRune.Width()
        flag += runeWidth
        //未处理到需要展示的字符位置
        if from.X() + 1 > flag {
            continue
        }
        if flag-from.X()-step == runeWidth {
            //刚好处理到一个完整字符(单/二倍宽字符)
            content[0][step] = tRune
            step += runeWidth
        } else {
            //刚好处理到一个二倍宽字符(被截断)
            if from.X() == to.X() {
                content[0][step] = gofer.NewRune(gofer.TextPadRight, nil, v.Style())
            } else {
                content[0][step] = gofer.NewRune(gofer.TextPadLeft, nil, v.Style())
            }
            step++
        }
        //已经处理到 to(point)
        if (from.X() < to.X() && flag >= to.X()+1) || from.X() == to.X() {
           break
        }
    }
    return content
}

func (v *TextView) SetText(text string) {
    strings.ReplaceAll(text, "\r", "")
    strings.ReplaceAll(text, "\n", "")
    v.text = gofer.StringToRunes(text, v.Style())
}

func (v *TextView) Text() string {
    length := len(v.text)
    if length == 0 {
        return ""
    }
    if length == 1 {
        return v.text[0].String()
    }
    var buf bytes.Buffer
    for _, tRune := range v.text {
        buf.WriteRune(tRune.Mainc())
        for _, trune := range tRune.Combc() {
            buf.WriteRune(trune)
        }
    }
    return buf.String()
}

func (v *TextView) TextWidth() int {
    textWidth := 0
    for _, tRune := range v.text {
        textWidth += tRune.Width()
    }
    return textWidth
}

func (v *TextView) SetHeight(h int) {
    //do nothing. the height is always equals to 1.
}

func NewTextView(width int) *TextView {
    view1 := &TextView{}
    view1.View = gofer.NewView()
    view1.View.SetWidth(width)
    view1.View.SetHeight(1)
    return view1
}
