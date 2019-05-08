package core

import (
    "errors"
    "github.com/GavinGuan24/gofer/log"
    "github.com/gdamore/tcell"
    "os"
)

type ApplicationDelegate interface {
    Launched(root View)
    //返回值为是否退出
    WillStop(code int) bool
}

type app struct {
    root *rootView
}

func NewApp() *app {
    return &app{}
}

func (ap *app) Run(delegate ApplicationDelegate) {
    //安检 app 代理
    if delegate == nil {
        log.Fatal(errors.New("ApplicationDelegate is Nil"))
    }
    //初始化根视图
    ap.root = RootView()
    if screen0, e0 := tcell.NewScreen(); e0 != nil {
        log.Fatal(e0)
    } else {
        ap.root.Screen = screen0
    }
    if e0 := ap.root.Screen.Init(); e0 != nil {
        log.Fatal(e0)
    }
    ap.root.Screen.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite).Normal())
    log.Info("Screen init finished.")

    //通知自己的监听器, 更新视图
    ap.root.UpdateUI(nil)

    exit := make(chan int)

    go func() {
        for {
            event := ap.root.Screen.PollEvent()
            switch event := event.(type) {
            case *tcell.EventKey:
                if event.Key() == tcell.KeyCtrlQ {
                    exit <- 0
                }
            default:
            }
        }
    }()

    //暴露部分权利交给代理
    log.Info("async: will call delegate.Launched.")
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
        log.Info("screen finalized.")
    }
    log.Info("App stopped.")
    os.Exit(0)
}
