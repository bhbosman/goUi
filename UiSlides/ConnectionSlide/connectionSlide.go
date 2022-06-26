package ConnectionSlide

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
	"github.com/cskr/pubsub"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type ConnectionSlide struct {
	data           IConnectionData
	connectionList *tview.List
	table          *tview.Table
	textView       *tview.TextView
	actionList     *tview.List
	next           tview.Primitive
	ctx            context.Context
	cancelFunc     context.CancelFunc
	channel        chan interface{}
	pubSub         *pubsub.PubSub
	app            *tview.Application
}

func (self *ConnectionSlide) UpdateContent() error {
	return nil
}

func (self *ConnectionSlide) Close() error {
	self.cancelFunc()
	close(self.channel)
	return nil
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

func (self *ConnectionSlide) goRun() {
	defer func(cmdChannel <-chan interface{}) {
		//flush
		for range cmdChannel {
		}
	}(self.channel)

	pubSubChannel := self.pubSub.Sub("ActiveConnectionStatus")
	defer func(pubSubChannel chan interface{}) {
		// unsubscribe on different go routine to avoid deadlock
		go func(pubSubChannel chan interface{}) {
			self.pubSub.Unsub(pubSubChannel)
			for range pubSubChannel {
			}
		}(pubSubChannel)
	}(pubSubChannel)

	var messageReceived interface{}
	var ok bool

	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.ctx,
		self.data,
		[]ChannelHandler.ChannelHandler{
			{
				PubSubHandler:  false,
				BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					return ISendMessage.ChannelEventsForISendMessage(next.(ISendMessage.ISendMessage), message)
				},
			},
			{
				PubSubHandler:  true,
				BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if sm, ok := next.(ISendMessage.ISendMessage); ok {
						_ = sm.Send(message)
					}
					return true, nil
				},
			},
		},
		func() int {
			return len(pubSubChannel) + len(self.channel)
		})
loop:
	for {
		select {
		case <-self.ctx.Done():
			break loop
		case messageReceived, ok = <-self.channel:
			if !ok {
				return
			}
			b, err := channelHandlerCallback(messageReceived, false)
			if err != nil || b {
				return
			}
		case messageReceived, ok = <-pubSubChannel:
			if !ok {
				return
			}
			b, err := channelHandlerCallback(messageReceived, true)
			if err != nil || b {
				return
			}
		}
	}
}

func (self *ConnectionSlide) SetConnectionListChange(list []IdAndName) {
	self.app.QueueUpdateDraw(func() {
		idx := self.connectionList.GetCurrentItem()
		self.connectionList.Clear()
		for _, s := range list {
			self.connectionList.AddItem(s.Id, s.Name, 0, func() {
				self.actionList.SetCurrentItem(0)
				self.app.SetFocus(self.actionList)
			})
		}
		self.connectionList.SetCurrentItem(idx)
		if len(list) == 0 {
			self.table.SetContent(nil)
			self.textView.Clear()
		}
	})
}

func (self *ConnectionSlide) SetConnectionInstanceChange(data *ConnectionData) {
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
	self.connectionList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		_, _ = ISendMessage.CallISendMessageSend(self.ctx, self.channel, false, &PublishInstanceDataFor{
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
		_, _ = ISendMessage.CallISendMessageSend(self.ctx, self.channel, false,
			NewDisconnectConnection(text, secondary))

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
	applicationContext context.Context,
	pubSub *pubsub.PubSub,
	app *tview.Application,
) *ConnectionSlide {
	ctx, cancelFunc := context.WithCancel(applicationContext)
	channel := make(chan interface{}, 32)

	data := NewData()
	result := &ConnectionSlide{
		data:       data,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		channel:    channel,
		pubSub:     pubSub,
		app:        app,
	}
	result.init()
	data.SetConnectionListChange(result.SetConnectionListChange)
	data.SetConnectionInstanceChange(result.SetConnectionInstanceChange)
	go result.goRun()
	return result
}
