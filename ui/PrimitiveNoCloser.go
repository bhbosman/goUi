package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PrimitiveNoCloser struct {
	primitive tview.Primitive
}

func (self *PrimitiveNoCloser) Draw(screen tcell.Screen) {
	self.primitive.Draw(screen)
}

func (self *PrimitiveNoCloser) GetRect() (int, int, int, int) {
	return self.primitive.GetRect()
}

func (self *PrimitiveNoCloser) SetRect(x, y, width, height int) {
	self.primitive.SetRect(x, y, width, height)
}

func (self *PrimitiveNoCloser) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return self.primitive.InputHandler()
}

func (self *PrimitiveNoCloser) Focus(delegate func(p tview.Primitive)) {
	self.primitive.Focus(delegate)
}

func (self *PrimitiveNoCloser) HasFocus() bool {
	return self.primitive.HasFocus()
}

func (self *PrimitiveNoCloser) Blur() {
	self.primitive.Blur()
}

func (self *PrimitiveNoCloser) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return self.primitive.MouseHandler()
}

func NewPrimitiveNoCloser(primitive tview.Primitive) *PrimitiveNoCloser {
	return &PrimitiveNoCloser{primitive: primitive}
}

func (self *PrimitiveNoCloser) Close() error {
	return nil
}
