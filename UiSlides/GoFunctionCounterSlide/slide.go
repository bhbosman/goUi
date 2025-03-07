package GoFunctionCounterSlide

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type slide struct {
	slideOrderNumber int
	slideName        string
	app              *tview.Application
	table            *tview.Table
	plate            *tablePlate
	canDraw          bool
	names            []string
}

func (self *slide) OrderNumber() int {
	return self.slideOrderNumber
}

func (self *slide) Name() string {
	return self.slideName
}

func (self *slide) Toggle(b bool) {
	self.canDraw = b
	if b {
		go func() {
			self.SetConnectionListChange(self.names)
		}()
	}
}

func (self *slide) Draw(screen tcell.Screen) {
	self.table.Draw(screen)
}

func (self *slide) GetRect() (int, int, int, int) {
	return self.table.GetRect()
}

func (self *slide) SetRect(x, y, width, height int) {
	self.table.SetRect(x, y, width, height)
}

func (self *slide) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return self.table.InputHandler()
}

func (self *slide) Focus(delegate func(p tview.Primitive)) {
	self.table.Focus(delegate)
}

func (self *slide) HasFocus() bool {
	return self.table.HasFocus()
}

func (self *slide) Blur() {
	self.table.Blur()
}

func (self *slide) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return self.table.MouseHandler()
}

func (self *slide) Close() error {
	return nil
}

func (self *slide) UpdateContent() error {
	return nil
}

func (self *slide) init() {
	self.plate = newTablePlate([]string{})
	self.table = tview.NewTable()
	self.table.SetSelectable(true, false)
	self.table.SetContent(self.plate)
	self.table.SetBorder(true)
	self.table.SetFixed(1, 1)
	self.table.SetTitle("Registered go Functions")
}

func (self *slide) SetConnectionListChange(names []string) {
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

func newSlide(slideOrderNumber int, slideName string, app *tview.Application) *slide {
	return &slide{
		slideOrderNumber: slideOrderNumber,
		app:              app,
		slideName:        slideName,
	}
}
