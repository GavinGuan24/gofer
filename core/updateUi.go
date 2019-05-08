package core

type UpdateUiMsg struct {
    view View
    rect Rect
}

func (msg *UpdateUiMsg)GetView() View {
    return msg.view
}

func (msg *UpdateUiMsg)GetRect() Rect {
    return msg.rect
}

func NewUpdateUiMsg(view View, rect Rect) *UpdateUiMsg {
    return &UpdateUiMsg{view, rect}
}
