# gofer

[gofer](https://github.com/GavinGuan24/gofer) 是一个基于 [tcell](https://github.com/gdamore/tcell) 的 TUI 框架(相对于GUI, 远程访问一台主机时, TUI更加高效).
ta 实现了基本的 view, app 概念. 且充分考虑Unicode排版问题, 充分利用 tcell 的按需更新的特性. 如果你使用 ta, 你可以快速构建一个基于终端中的 TUI 项目.

当然, superView(father) 可以添加 subview(kid). 如果你有过 iOS 开发经历, 你会发现一些视图的实现逻辑是相似的.
我实现了 view 的内容更新时的整个响应链逻辑(以golang的风格), view 的内容更新仅需要调用 `core.UpdateUI(view0, rect0)`.
因为TUI项目的视图层级关系不应过于复杂, 又同时有较多的快捷键支持的需求, 所以我并未像iOS的coco框架一样直接实现整个事件响应链.
因为我希望使用者自行转发监听到的事件给对应的view.(关于事件响应链的处理方式, 我在考虑是否应该进行一些封装)

后续, 我将封装一些常用的 widget.

## version

|version|date|commit id|
|:---:|:---:|:---:|
|0.0.3|20190717|master|
|0.0.2|20190603|ebb54da62|
|0.0.1|20190525|9fb25da9d|

## change
##### 20190717
在 widget 中新增 `TextView` (一个单行显示的文本框).
对 `interface Rune` 的默认实现 `basicRune` 的 `String()` 方法做性能优化.

##### 20190603
`EventResize` 交给 `ApplicationDelegate` 处理; `UpdateUiMsg` 私有化, 避免下游开发者的困扰.

## bugfix

##### 20190717
对 [tcell](https://github.com/gdamore/tcell) 的一个发生在 `screen` 右边缘的bug进行补救.
对 `app` 在调用 `delegate.Launched()` 后, 可能出现的视图未初始化的问题进行修复.

##### 20190603
`rootView` 处理 subview 消息时, 仅更新该 subview 导致内容优先级出错.

