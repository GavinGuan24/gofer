package views

import (
    "bytes"
    "fmt"
    "github.com/gdamore/tcell"
    "github.com/mattn/go-runewidth"
)

// 某视图中字符位的描述, 位置由使用 Rune 的程序维护
// combc 一般为 nil, 除非要处理 ZWJ(zero-width joiner) 问题
type Rune struct {
    mainc rune
    combc []rune
    style tcell.Style
}

func (r *Rune) String() string {
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
        return fmt.Sprintf("%s",buf.String())
    }
}

func (r *Rune) Width() int {
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

func NewRune(mainc rune, combc []rune, style tcell.Style) *Rune {
    return &Rune{mainc, combc, style}
}

//以下这个文档是为了以后优化 方法<code>func (r *Rune) Width() int</code>
//
//---------- rune 宽度判定(UTF-8), 关于 Unicode 12.0.0 文档的说明
// http://www.unicode.org/Public/UCD/latest/ucd/EastAsianWidth.txt
// 该文档表示, ta的内容是一个表, 用英文分号(;)分割, 每一行有两个字段
//     字段0: Unicode码点值或码点值范围
//     字段1: East_Asian_Width属性 (ea)
//-----------------------------------------------------------------
// 1. 对 ea 属性我进行以下解读
//     单字符宽度即, 在确定的字号下, ASCII码可以表示的英文字母所占的宽度
//
//     ea 可能的所有值
//
// # East_Asian_Width (ea)
// ea ; A         ; Ambiguous # 不明确的
// ea ; N         ; Neutral   # 中立的
// --------------------------
// ea ; H         ; Halfwidth # 明确的半宽度, 即为单字符宽度
// ea ; Na        ; Narrow    # 窄的, 在终端屏幕上可以处理为单字符宽度
// --------------------------
// ea ; F         ; Fullwidth # 明确的全宽度, 即为双字符宽度
// ea ; W         ; Wide      # 宽的, 在终端屏幕上可以处理为双字符宽度
//
// 所以在我看来
//     H/Na 处理为单字符宽度
//     F/W 处理为双字符宽度
//     A/N 人工判定(这类字符数量非常多, 没有处理前, 一律按照单字符宽度处理)
//         我想到的解决方法是: 使用golang写一个输出字符的程序, 截屏为图像. 然后让Python的图形库解析图像, 这样来判定(A/N)字符的宽度.
//         mac上截屏取图的golang代码如下文, 我不熟悉Python的图形分析, 所以没写.
//         同时我计算了一下时间(MacBook Pro Retina, 13-inch, Early 2015), 从 0 ~ 0x10FFFE 需要截屏 4个多小时. python 分析的话, 我就不清楚了
//
//package main
//
//import (
//    "fmt"
//    "os/exec"
//    "time"
//)
//
//func main() {
//    var i rune
//    for i = 0; i < 0x10FFFE; i++ {
//        TestCharWidth(i)
//        time.Sleep(time.Millisecond*time.Duration(11))
//        c := fmt.Sprintf("screencapture -xR 0,45,90,148 ./chars/%06X.jpg", i)
//        cmd := exec.Command("bash", "-c", c)
//        _, _ = cmd.Output()
//    }
//}
//
//func TestCharWidth(r rune) {
//    fmt.Printf("=>%X\n", r)
//    fmt.Printf("%c|\n", r)
//    fmt.Printf("|%c|\n", r)
//    fmt.Printf("e%c|\n", r)
//    fmt.Println("-----")
//}
//-----------------------------------------------------------------
// 2. 文档开头强调的强制
//     2.1 未明确列出的所有码点, 不论是否被分配, 一致视为 N. (单字符宽度)
//     2.2 U+3400..U+4DBF, U+4E00..U+9FFF, U+F900..U+FAFF 这几个块一致为 W. (双字符宽度)
//     2.3 U+20000..U+2FFFD, U+30000..U+3FFFD 这几层一致为 W. (双字符宽度)
//-----------------------------------------------------------------
// 3. 该文档中未提及零宽字符(例: u200d), 也没有提及提供字形变化的字符集COMBINING(例: u030a, u20e3)
//     这类字符有些宽度为0, 有些可以独立在行首, 有些需要作为另一个字符的COMBINING(有些会影响整体在屏幕上展示的结果的宽度)
//-----------------------------------------------------------------
// 参考文档, 我写下本文件时, Unicode 归档版本是12.0.0
// [unicode 东亚宽度 12.0.0](http://www.unicode.org/Public/UCD/latest/ucd/EastAsianWidth.txt, https://www.unicode.org/Public/12.0.0/ucd/EastAsianWidth.txt)
// [unicode 例子](https://www.jianshu.com/p/2c75d107c187)
// [编码对照表](https://blog.csdn.net/coolwu/article/details/79752396)
// [unicode 与安全](https://blog.csdn.net/P5dEyT322JACS/article/details/79454805)
//

