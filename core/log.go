package core

import (
    "bytes"
    "fmt"
    "github.com/gdamore/tcell"
    "os"
    "path/filepath"
    "strings"
    "time"
)

var screenInLog tcell.Screen
var logger *os.File

func init() {
    file, e := os.OpenFile("./gofer.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
    if e != nil {
        logFatal(e)
    }
    logger = file
}

func LogDebug(str string) {
    var buf bytes.Buffer
    buf.WriteString("\u001b[0;0;36mDebug > ")
    buf.WriteString(str)
    buf.WriteString("\u001b[0m")
    logout(buf.String())
}

func LogInfo(str string) {
    var buf bytes.Buffer
    buf.WriteString("Info > ")
    buf.WriteString(str)
    logout(buf.String())
}

func LogWarn(str string) {
    var buf bytes.Buffer
    buf.WriteString("\u001b[38;2;179;135;29mWarn > ")
    buf.WriteString(str)
    buf.WriteString("\u001b[0m")
    logout(buf.String())
}

func LogError(str string) {
    var buf bytes.Buffer
    buf.WriteString("\u001b[0;0;31mError > ")
    buf.WriteString(str)
    buf.WriteString("\u001b[0m")
    logout(buf.String())
}

func LogCrash(str string) {
    var buf bytes.Buffer
    buf.WriteString("\u001b[1;48;2;254;218;49;31mCrash > ")
    buf.WriteString(str)
    buf.WriteString("\u001b[0m")
    logout(buf.String())
    if screenInLog != nil {
        screenInLog.Fini()
    }
    logFilename := logger.Name()
    dirPath, _ := filepath.Abs(filepath.Dir(logFilename))
    logFatal(fmt.Errorf("app is terminated.\nthe crash log has been saved in the (%v) file", dirPath + logFilename[strings.LastIndex(logFilename, "/"):]))
}

//https://gochannel.org/links/link/snapshot/7352
// 前景 背景 颜色
// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色
//
// 模式
//  0  终端默认设置
//  1  高亮显示
//  4  使用下划线
//  5  闪烁
//  7  反白显示
//  8  不可见
// 极简的格式为
// %c[%d;%d;%dm ==> 着色, 模式, 背景色, 前景色
// %s
// %c[0m
//func color(str string) string {
//    var buf bytes.Buffer
//    buf.WriteString("\u001b[0;0;31m")
//    buf.WriteString(str)
//    buf.WriteString("\u001b[0m")
//    return buf.String()
//}

func logout(str string) {
    var buf bytes.Buffer
    buf.WriteString(time.Now().Format("2006-01-02 15:04:05.000 "))
    buf.WriteString(str)
    fmt.Fprintln(logger, buf.String())
}

// 启动时的失败, TUI的渲染会占用 stdout, 所以输出到 stderr
func logFatal(e error) {
    fmt.Fprint(os.Stderr, "%v\n", e.Error())
    os.Exit(1)
}

