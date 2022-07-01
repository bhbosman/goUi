package GoFunctionCounterSlide

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Slide struct {
	next    tview.Primitive
	app     *tview.Application
	table   *tview.Table
	plate   *tablePlate
	canDraw bool
	names   []string
}

func (self *Slide) Toggle(b bool) {
	self.canDraw = b
	if b {
		go func() {
			self.SetConnectionListChange(self.names)
		}()
	}
}

func NewSlide(app *tview.Application) *Slide {
	return &Slide{
		app: app,
	}
}

func (self *Slide) Draw(screen tcell.Screen) {
	self.next.Draw(screen)
}

func (self *Slide) GetRect() (int, int, int, int) {
	return self.next.GetRect()
}

func (self *Slide) SetRect(x, y, width, height int) {
	self.next.SetRect(x, y, width, height)
}

func (self *Slide) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return self.next.InputHandler()
}

func (self *Slide) Focus(delegate func(p tview.Primitive)) {
	self.next.Focus(delegate)
}

func (self *Slide) HasFocus() bool {
	return self.next.HasFocus()
}

func (self *Slide) Blur() {
	self.next.Blur()
}

func (self *Slide) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return self.next.MouseHandler()
}

func (self *Slide) Close() error {
	return nil
}

func (self *Slide) UpdateContent() error {
	return nil
}

func (self *Slide) init() {
	self.plate = newTablePlate([]string{})
	self.table = tview.NewTable()
	self.table.SetSelectable(true, false)
	self.table.SetContent(self.plate)
	self.table.SetBorder(true)
	self.table.SetFixed(1, 1)
	self.table.SetTitle("Registered go Functions")
	self.next = self.table
}

func (self *Slide) SetConnectionListChange(names []string) {
	self.names = names
	if self.canDraw {
		self.app.QueueUpdate(
			func() {
				self.plate = newTablePlate(self.names)
				self.table.SetContent(self.plate)
				//
				if self.canDraw {
					self.app.ForceDraw()
				}
			},
		)
	}
}
