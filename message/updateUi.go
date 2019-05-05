package message

import "github.com/GavinGuan24/gofer/views"

type UpdateUiMsg struct {
    view views.View
    rect views.Rect
}

func (msg *UpdateUiMsg)GetView() views.View {
    return msg.view
}

func (msg *UpdateUiMsg)GetRect() views.Rect {
    return msg.rect
}

func NewUpdateUiMsg(view views.View, rect views.Rect) *UpdateUiMsg {
    return &UpdateUiMsg{view, rect}
}
