package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"go.uber.org/multierr"
	"io"
)

type PrimitiveWithCloser struct {
	primitive tview.Primitive
	closers   []io.Closer
}

func (self *PrimitiveWithCloser) Draw(screen tcell.Screen) {
	self.primitive.Draw(screen)
}

func (self *PrimitiveWithCloser) GetRect() (int, int, int, int) {
	return self.primitive.GetRect()
}

func (self PrimitiveWithCloser) SetRect(x, y, width, height int) {
	self.primitive.SetRect(x, y, width, height)
}

func (self *PrimitiveWithCloser) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return self.primitive.InputHandler()
}

func (self *PrimitiveWithCloser) Focus(delegate func(p tview.Primitive)) {
	self.primitive.Focus(delegate)
}

func (self *PrimitiveWithCloser) HasFocus() bool {
	return self.primitive.HasFocus()
}

func (self *PrimitiveWithCloser) Blur() {
	self.primitive.Blur()
}

func (self *PrimitiveWithCloser) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return self.primitive.MouseHandler()
}

func (self *PrimitiveWithCloser) Close() error {
	var err error
	err = nil
	for _, closer := range self.closers {
		err = multierr.Append(err, closer.Close())
	}
	return err
}

func NewPrimitiveWithCloser(primitive tview.Primitive, closers []io.Closer) *PrimitiveWithCloser {
	return &PrimitiveWithCloser{
		primitive: primitive,
		closers:   closers,
	}
}
