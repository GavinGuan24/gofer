# gofer

[gofer](https://github.com/GavinGuan24/gofer) 是一个基于 [tcell](https://github.com/gdamore/tcell) 的 TUI 框架(相对于GUI, 远程访问一台主机时, TUI更加高效).
ta 实现了基本的 view, app 概念. 如果你使用 ta, 你可以快速构建一个基于终端中的 TUI 项目.

当然, superView(father) 可以添加 subview(kid). 如果你有过 iOS 开发经历, 你会发现一些实现逻辑是相似的.
我实现了 view 的内容更新时的整个响应链逻辑(以golang的风格), view 的内容更新仅需要调用 `core.UpdateUI(view0, rect0)`.
因为TUI项目的视图层级关系不应过于复杂, 又同时有较多的快捷键支持的需求, 所以我并未像iOS的coco框架一样直接实现整个事件响应链.
因为我希望使用者自行转发监听到的事件给对应的view.

后续, 我将封装一些常用的 widget.

