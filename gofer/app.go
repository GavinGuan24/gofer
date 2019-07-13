package gofer

import (
    "errors"
    "github.com/gdamore/tcell"
    "os"
)

type ApplicationDelegate interface {
    Launched(root RootView)
    //返回值为是否继续执行退出
    WillStop(code int) bool
    //是否启用鼠标
    EnableMouse() bool
    //返回一个事件监听者
    EventListener() chan<- tcell.Event
}

type app struct {
    root *rootView
    listener chan<- tcell.Event
}

func NewApp() *app {
    return &app{}
}

func (ap *app) Run(delegate ApplicationDelegate) {
    //安检 app 代理
    if delegate == nil {
        logFatal(errors.New("ApplicationDelegate is Nil"))
    }
    //初始化根视图
    ap.root = NewRootView()
    if screen0, e0 := tcell.NewScreen(); e0 != nil {
        logFatal(e0)
    } else {
        ap.root.Screen = screen0
    }
    if e0 := ap.root.Screen.Init(); e0 != nil {
        logFatal(e0)
    }
    if delegate.EnableMouse() {
        ap.root.Screen.EnableMouse()
    }
    ap.root.Screen.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite).Normal())
    //将screen 暴露给 log.go
    screenInLog = ap.root.Screen
    LogInfo("Screen init finished.")

    //更新视图
    UpdateUI(ap.root, nil)

    exit := make(chan int)

    listener := delegate.EventListener()
    go func() {
        for {
            event := ap.root.Screen.PollEvent()
            if event == nil {
                continue
            }
            switch event := event.(type) {
            case *tcell.EventKey:
                if event.Key() == tcell.KeyCtrlQ {
                    exit <- 0
                } else {
                    listener <- event
                }
            case *tcell.EventResize:
                ap.root.Sync()
                UpdateUI(ap.root, nil)
                listener <- event
            default:
                listener <- event
            }
        }
    }()

    //暴露部分权利交给代理(这里不用再调用一次UpdateUI(); 因为事件循环中, EventResize 会自动触发UpdateUI(); 窗口初始化后, 框架会自动发送一次EventResize)
    LogInfo("Async: will call delegate.Launched.")
    go delegate.Launched(ap.root)

    for {
        select {
        case code := <-exit:
            if delegate.WillStop(code) {
                ap.stop()
            }
        default:
        }
    }
}

func (ap *app) stop() {
    if ap.root != nil {
        ap.root.Screen.Fini()
        LogInfo("Screen finalized.")
    }
    if ap.listener != nil {
        close(ap.listener)
    }
    LogInfo("App stopped.\n\n")
    os.Exit(0)
}
