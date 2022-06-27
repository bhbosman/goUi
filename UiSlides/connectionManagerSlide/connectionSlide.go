package connectionManagerSlide

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type ConnectionSlide struct {
	service        IConnectionSlideService
	connectionList *tview.List
	table          *tview.Table
	textView       *tview.TextView
	actionList     *tview.List
	next           tview.Primitive
	app            *tview.Application
}

func (self *ConnectionSlide) UpdateContent() error {
	return nil
}

func (self *ConnectionSlide) Close() error {
	return nil
	//return self.service.OnStop(context.Background())
}

func (self *ConnectionSlide) Draw(screen tcell.Screen) {
	self.next.Draw(screen)
}

func (self *ConnectionSlide) GetRect() (int, int, int, int) {
	return self.next.GetRect()
}

func (self *ConnectionSlide) SetRect(x, y, width, height int) {
	self.next.SetRect(x, y, width, height)
}

func (self *ConnectionSlide) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return self.next.InputHandler()
}

func (self *ConnectionSlide) Focus(delegate func(p tview.Primitive)) {
	self.next.Focus(delegate)
}

func (self *ConnectionSlide) HasFocus() bool {
	return self.next.HasFocus()
}

func (self *ConnectionSlide) Blur() {
	self.next.Blur()
}

func (self *ConnectionSlide) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return self.next.MouseHandler()
}

func (self *ConnectionSlide) SetConnectionListChange(list []IdAndName) {
	self.app.QueueUpdateDraw(
		func() {
			idx := self.connectionList.GetCurrentItem()
			self.connectionList.Clear()
			for _, s := range list {
				self.connectionList.AddItem(s.Id, s.Name, 0,
					func() {
						self.actionList.SetCurrentItem(0)
						self.app.SetFocus(self.actionList)
					},
				)
			}
			self.connectionList.SetCurrentItem(idx)
			if len(list) == 0 {
				self.table.SetContent(nil)
				self.textView.Clear()
			}
		})
}

func (self *ConnectionSlide) SetConnectionInstanceChange(data *ConnectionInstanceData) {
	self.app.QueueUpdateDraw(func() {
		index := self.connectionList.GetCurrentItem()
		text, _ := self.connectionList.GetItemText(index)
		if text == data.ConnectionId {
			if data != nil && data.Grid != nil {
				tableData := newConnectionPlateContent(data.Grid)
				if tableData != nil {
					self.table.SetContent(tableData)
				}
				self.textView.Clear()
				_, _ = fmt.Fprintf(self.textView, "Name: %v\n", data.Name)
				_, _ = fmt.Fprintf(self.textView, "Id: %v\n", data.ConnectionId)
				_, _ = fmt.Fprintf(self.textView, "Connect Time: %v, (%v)\n", data.ConnectionTime.Format(time.RFC3339), time.Now().Sub(data.ConnectionTime))
			}
		}
	})
}

func (self *ConnectionSlide) init() {
	self.connectionList = tview.NewList().ShowSecondaryText(true)
	self.connectionList.SetBorder(true).SetTitle("Active Connections")
	self.connectionList.SetChangedFunc(
		func(index int, mainText string, secondaryText string, shortcut rune) {
			_ = self.service.Send(
				&PublishInstanceDataFor{
					Id:   mainText,
					Name: secondaryText,
				})
		})

	self.actionList = tview.NewList().ShowSecondaryText(false)
	self.actionList.SetBorder(true).SetTitle("Actions")
	self.actionList.AddItem("..", "", 0, func() {
		self.app.SetFocus(self.connectionList)
	})
	self.actionList.AddItem("Disconnect", "", 0, func() {
		index := self.connectionList.GetCurrentItem()
		text, secondary := self.connectionList.GetItemText(index)
		_ = self.service.Send(
			NewDisconnectConnection(
				text,
				secondary,
			),
		)

		// Todo: this throws a panic on the first line. need to fix tview
		//self.connectionList.RemoveItem(index)
		self.actionList.SetCurrentItem(0)
		self.app.SetFocus(self.connectionList)
	})
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
					AddItem(self.actionList, 4, 2, false),
					0,
					3,
					true).
				AddItem(tview.NewFlex().
					SetDirection(tview.FlexRow).
					AddItem(self.textView, 5, 0, false).
					AddItem(self.table, 0, 6, false), 0, 6, false),
			0,
			1,
			true)

}

func NewConnectionSlide(
	app *tview.Application,
	service *Service,
) (*ConnectionSlide, error) {
	result := &ConnectionSlide{
		service: service,
		app:     app,
	}
	result.service.SetConnectionListChange(result.SetConnectionListChange)
	result.service.SetConnectionInstanceChange(result.SetConnectionInstanceChange)
	result.init()
	return result, nil
}
