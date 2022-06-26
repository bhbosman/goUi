package ConnectionSlide

import (
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"sort"
)

type Data struct {
	connectionListIsDirty      bool
	ConnectionDataMap          map[string]*ConnectionInstanceData
	messageRouter              *messageRouter.MessageRouter
	onConnectionListChange     func(connectionList []IdAndName)
	onConnectionInstanceChange func(data *ConnectionInstanceData)
}

func NewData() *Data {
	result := &Data{
		ConnectionDataMap: make(map[string]*ConnectionInstanceData),
		messageRouter:     messageRouter.NewMessageRouter(),
	}
	_ = result.messageRouter.Add(result.handleEmptyQueue)
	_ = result.messageRouter.Add(result.handleConnectionState)
	_ = result.messageRouter.Add(result.handleConnectionCreated)
	_ = result.messageRouter.Add(result.handleConnectionClosed)
	_ = result.messageRouter.Add(result.handlePublishInstanceDataFor)
	_ = result.messageRouter.Add(result.handleDisconnectConnection)
	return result
}

func (self *Data) Send(data interface{}) error {
	_, err := self.messageRouter.Route(data)
	return err
}

func (self *Data) handleDisconnectConnection(message *DisconnectConnection) error {
	if info, ok := self.ConnectionDataMap[message.ConnectionId]; ok {
		if info.CancelFunc != nil {
			info.CancelFunc()
		}
	}
	return nil
}
func (self *Data) handlePublishInstanceDataFor(message *PublishInstanceDataFor) error {
	if info, ok := self.ConnectionDataMap[message.Id]; ok {
		self.DoConnectionInstanceChange(info)
	}
	return nil
}
func (self *Data) handleEmptyQueue(_ *messages.EmptyQueue) error {
	if self.connectionListIsDirty {
		self.DoConnectionListChange()
		self.connectionListIsDirty = false
	}
	for _, connectionData := range self.ConnectionDataMap {
		if connectionData.isDirty {
			self.DoConnectionInstanceChange(connectionData)
			connectionData.isDirty = false
		}
	}
	return nil
}

func (self *Data) handleConnectionState(message *model.ConnectionState) error {
	if data, ok := self.ConnectionDataMap[message.ConnectionId]; ok {
		data.isDirty = true
		data.CancelContext = message.CancelContext
		data.CancelFunc = message.CancelFunc
		data.Name = message.Name
		//data.Status = message.Status
		data.ConnectionTime = message.ConnectionTime
		data.Grid = message.Grid
	}
	return nil
}

func (self *Data) handleConnectionClosed(message *model.ConnectionClosed) error {
	delete(self.ConnectionDataMap, message.ConnectionId)
	self.connectionListIsDirty = true
	return nil
}

func (self *Data) handleConnectionCreated(message *model.ConnectionCreated) error {
	self.ConnectionDataMap[message.ConnectionId] = &ConnectionInstanceData{
		isDirty:      true,
		ConnectionId: message.ConnectionId,
		Name:         message.ConnectionName,
	}
	self.connectionListIsDirty = true
	return nil
}

func (self *Data) DoConnectionListChange() {
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
func (self *Data) DoConnectionInstanceChange(data *ConnectionInstanceData) {
	if self.onConnectionInstanceChange != nil {
		self.onConnectionInstanceChange(data)
	}
}

func (self *Data) SetConnectionInstanceChange(cb func(data *ConnectionInstanceData)) {
	self.onConnectionInstanceChange = cb
}
func (self *Data) SetConnectionListChange(cb func(connectionList []IdAndName)) {
	self.onConnectionListChange = cb
}
