package cmSlide

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type slide struct {
	service          IConnectionSlideService
	connectionList   *tview.Table
	table            *tview.Table
	tableStrings     *tview.Table
	textView         *tview.TextView
	actionList       *tview.List
	next             tview.Primitive
	app              *tview.Application
	canDraw          bool
	connectionPlate  *connectionPlate
	slideName        string
	slideOrderNumber int
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
		self.app.ForceDraw()
	}
}

func (self *slide) UpdateContent() error {
	return nil
}

func (self *slide) Close() error {
	return nil
}

func (self *slide) Draw(screen tcell.Screen) {
	if self.canDraw {
		self.next.Draw(screen)
	}
}

func (self *slide) GetRect() (int, int, int, int) {
	return self.next.GetRect()
}

func (self *slide) SetRect(x, y, width, height int) {
	self.next.SetRect(x, y, width, height)
}

func (self *slide) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return self.next.InputHandler()
}

func (self *slide) Focus(delegate func(p tview.Primitive)) {
	self.next.Focus(delegate)
}

func (self *slide) HasFocus() bool {
	return self.next.HasFocus()
}

func (self *slide) Blur() {
	self.next.Blur()
}

type MouseHandlerCallback = func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive)

func (self *slide) MouseHandler() MouseHandlerCallback {
	return self.next.MouseHandler()
}

func (self *slide) SetConnectionListChange(list []IdAndName) {
	self.app.QueueUpdate(
		func() {
			self.connectionPlate = newConnectionPlace(list)
			rows, columns := self.connectionList.GetSelectable()
			self.connectionList.SetContent(self.connectionPlate)
			self.connectionList.SetSelectable(rows, columns)
			if len(list) == 0 {
				self.table.SetContent(nil)
				self.textView.Clear()
			}
			if self.canDraw {
				self.app.ForceDraw()
			}
		},
	)
}

func (self *slide) SetConnectionInstanceChange(data ConnectionInstanceData) {
	self.app.QueueUpdate(
		func() {
			row, _ := self.connectionList.GetSelection()
			if text, ok := self.connectionPlate.GetItem(row); ok {
				if text.Id == data.ConnectionId {
					if data.Grid != nil {
						tableData := newConnectionPlateContent(data.Grid)
						if tableData != nil {
							self.table.SetContent(tableData)
						}
					}

					if data.KeyValue != nil {
						stringValueData := newStringsPlate(data.KeyValue)
						if stringValueData != nil {
							self.tableStrings.SetContent(stringValueData)
							self.tableStrings.ScrollToBeginning()
						}
					}

					self.textView.Clear()
					_, _ = fmt.Fprintf(self.textView, "Name: %v\n", data.Name)
					_, _ = fmt.Fprintf(self.textView, "Id: %v\n", data.ConnectionId)
					_, _ = fmt.Fprintf(self.textView, "Connect Time: %v, (%v)\n", data.ConnectionTime.Format(time.RFC3339), time.Now().Sub(data.ConnectionTime))
				}
				if self.canDraw {
					self.app.ForceDraw()
				}
			}
		},
	)
}

func (self *slide) init() {
	self.connectionList = tview.NewTable() //.ShowSecondaryText(true)
	self.connectionList.SetSelectionChangedFunc(
		func(row, column int) {
			row, _ = self.connectionList.GetSelection()
			if item, ok := self.connectionPlate.GetItem(row); ok {
				_ = self.service.Send(
					&publishInstanceDataFor{
						Id:   item.Id,
						Name: item.Name,
					},
				)
			}
		},
	)
	self.connectionList.SetSelectedFunc(
		func(row, column int) {
			self.actionList.SetCurrentItem(0)
			self.app.SetFocus(self.actionList)
		},
	)
	self.connectionList.SetSelectable(true, false)
	self.connectionList.SetBorder(true).SetTitle("Active Connections")
	self.connectionList.SetFixed(1, 1)
	self.actionList = tview.NewList().ShowSecondaryText(false)
	self.actionList.SetBorder(true).SetTitle("Actions")
	self.actionList.AddItem("..", "", 0,
		func() {
			self.app.SetFocus(self.connectionList)
		},
	)
	self.actionList.AddItem("Disconnect", "", 0,
		func() {
			row, _ := self.connectionList.GetSelection()
			if item, ok := self.connectionPlate.GetItem(row); ok {
				self.service.DisconnectConnection(item.Id)
			}
			self.actionList.SetCurrentItem(0)
			self.app.SetFocus(self.connectionList)
		},
	)
	self.actionList.AddItem("Disconnect All", "", 0,
		func() {
			self.service.DisconnectAllConnections()
			self.actionList.SetCurrentItem(0)
			self.app.SetFocus(self.connectionList)
		},
	)
	self.actionList.AddItem("Reset Counters", "", 0,
		func() {
			row, _ := self.connectionList.GetSelection()
			if item, ok := self.connectionPlate.GetItem(row); ok {
				self.service.ResetConnectionParams(item.Id)
			}
			self.actionList.SetCurrentItem(0)
			self.app.SetFocus(self.connectionList)
		},
	)
	self.actionList.AddItem("Reset All Counters", "", 0,
		func() {
			self.service.ResetAllConnectionParams()
			self.actionList.SetCurrentItem(0)
			self.app.SetFocus(self.connectionList)
		},
	)

	self.tableStrings = tview.NewTable()
	self.tableStrings.SetTitle("Strings")
	self.tableStrings.SetBorder(true)
	self.tableStrings.SetFixed(1, 1)

	self.table = tview.NewTable()
	self.table.SetTitle("Connection Stack").SetBorder(true)
	self.textView = tview.NewTextView()
	self.textView.SetTitle("Connection Information").SetBorder(true)
	self.table.SetBorder(true)
	self.next = tview.NewFlex().
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(tview.NewFlex().
					SetDirection(tview.FlexRow).
					AddItem(self.connectionList, 0, 3, true).
					AddItem(self.actionList, 7, 2, false),
					0,
					3,
					true).
				AddItem(tview.NewFlex().
					SetDirection(tview.FlexRow).
					AddItem(self.textView, 5, 0, false).
					AddItem(self.table, 10, 6, false).
					AddItem(self.tableStrings, 0, 6, false),
					0, 6, false),
			0,
			1,
			true)

}

func newConnectionSlide(
	slideName string,
	app *tview.Application,
	service IConnectionSlideService,
) (*slide, error) {
	result := &slide{
		slideName: slideName,
		service:   service,
		app:       app,
	}
	result.init()
	return result, nil
}
