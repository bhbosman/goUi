package cmSlide

import (
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"sort"
)

type data struct {
	connectionListIsDirty      bool
	ConnectionDataMap          map[string]*ConnectionInstanceData
	messageRouter              *messageRouter.MessageRouter
	onConnectionListChange     func(connectionList []IdAndName)
	onConnectionInstanceChange func(data ConnectionInstanceData)
}

func (self *data) ShutDown() error {
	return nil
}

func NewData() (*data, error) {
	result := &data{
		ConnectionDataMap: make(map[string]*ConnectionInstanceData),
		messageRouter:     messageRouter.NewMessageRouter(),
	}
	_ = result.messageRouter.Add(result.handleEmptyQueue)
	_ = result.messageRouter.Add(result.handleConnectionState)
	_ = result.messageRouter.Add(result.handleConnectionCreated)
	_ = result.messageRouter.Add(result.handleConnectionClosed)
	_ = result.messageRouter.Add(result.handlePublishInstanceDataFor)
	_ = result.messageRouter.Add(result.handleDisconnectConnection)
	return result, nil
}

func (self *data) Send(data interface{}) error {
	self.messageRouter.Route(data)
	return nil
}

func (self *data) handleDisconnectConnection(message *DisconnectConnection) error {
	if info, ok := self.ConnectionDataMap[message.ConnectionId]; ok {
		if info.CancelFunc != nil {
			info.CancelFunc()
		}
	}
	return nil
}
func (self *data) handlePublishInstanceDataFor(message *publishInstanceDataFor) error {
	if info, ok := self.ConnectionDataMap[message.Id]; ok {
		self.DoConnectionInstanceChange(info)
	}
	return nil
}

func (self *data) handleEmptyQueue(_ *messages.EmptyQueue) {
	didSomething := self.connectionListIsDirty
	if self.connectionListIsDirty {
		self.DoConnectionListChange()
		self.connectionListIsDirty = false
	}
	for _, connectionData := range self.ConnectionDataMap {
		didSomething = didSomething || connectionData.isDirty
		if connectionData.isDirty {
			self.DoConnectionInstanceChange(connectionData)
			connectionData.isDirty = false
		}
	}
}

func (self *data) handleConnectionState(message *model.ConnectionState) {
	if dataInstance, ok := self.ConnectionDataMap[message.ConnectionId]; ok {
		dataInstance.isDirty = true
		dataInstance.CancelContext = message.CancelContext
		dataInstance.CancelFunc = message.CancelFunc
		dataInstance.Name = message.Name
		dataInstance.ConnectionTime = message.ConnectionTime
		dataInstance.Grid = message.Grid
	}
}

func (self *data) handleConnectionClosed(message *model.ConnectionClosed) error {
	delete(self.ConnectionDataMap, message.ConnectionId)
	self.connectionListIsDirty = true
	return nil
}

func (self *data) handleConnectionCreated(message *model.ConnectionCreated) error {
	self.ConnectionDataMap[message.ConnectionId] = &ConnectionInstanceData{
		isDirty:        true,
		ConnectionId:   message.ConnectionId,
		Name:           message.ConnectionName,
		CancelFunc:     message.CancelFunc,
		CancelContext:  message.CancelContext,
		ConnectionTime: message.ConnectionTime,
	}
	self.connectionListIsDirty = true
	return nil
}

func (self *data) DoConnectionListChange() {
	if self.onConnectionListChange != nil {
		ss := make([]string, 0, len(self.ConnectionDataMap))

		for key := range self.ConnectionDataMap {
			ss = append(ss, key)
		}
		sort.Strings(ss)
		cbData := make([]IdAndName, 0, len(self.ConnectionDataMap))

		for _, s := range ss {
			if info, ok := self.ConnectionDataMap[s]; ok {
				idAndName := IdAndName{
					Id:   info.ConnectionId,
					Name: info.Name,
				}
				cbData = append(cbData, idAndName)
			}
		}
		self.onConnectionListChange(cbData)
	}
}
func (self *data) DoConnectionInstanceChange(data *ConnectionInstanceData) {
	if self.onConnectionInstanceChange != nil {
		self.onConnectionInstanceChange(*data)
	}
}

func (self *data) SetConnectionInstanceChange(cb func(data ConnectionInstanceData)) {
	self.onConnectionInstanceChange = cb
}
func (self *data) SetConnectionListChange(cb func(connectionList []IdAndName)) {
	self.onConnectionListChange = cb
}
