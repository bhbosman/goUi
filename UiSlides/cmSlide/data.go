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
	messageRouter              messageRouter.IMessageRouter
	onConnectionListChange     func(connectionList []IdAndName)
	onConnectionInstanceChange func(data ConnectionInstanceData)
	onSendMessageToService     func(interface{})
}

func (self *data) ResetConnectionParams(connectionId string) {
	if info, ok := self.ConnectionDataMap[connectionId]; ok {
		if info.CancelFunc != nil {
			info.ResetConnectionParams()
		}
	}

}

func (self *data) ResetAllConnectionParams() {
	for _, value := range self.ConnectionDataMap {
		if value.CancelFunc != nil {
			value.ResetConnectionParams()
		}
	}
}

func (self *data) DisconnectConnection(connectionId string) {
	if info, ok := self.ConnectionDataMap[connectionId]; ok {
		if info.CancelFunc != nil {
			info.CancelFunc()
		}
	}
}

func (self *data) DisconnectAllConnections() {
	for _, value := range self.ConnectionDataMap {
		if value.CancelFunc != nil {
			value.CancelFunc()
		}
	}

	// TODO: something to think about
	// At this point I can kill the current ConnectionDataMap. May be do this ??
	//self.ConnectionDataMap = make(map[string]*ConnectionInstanceData)
	//self.connectionListIsDirty = true
}

func (self *data) SendMessageToService(cb func(interface{})) {
	self.onSendMessageToService = cb
}

func (self *data) ShutDown() error {
	return nil
}

func NewData() (IConnectionSlideData, error) {
	result := &data{
		ConnectionDataMap: make(map[string]*ConnectionInstanceData),
		messageRouter:     messageRouter.NewMessageRouter(),
	}
	_ = result.messageRouter.Add(result.handleEmptyQueue)
	_ = result.messageRouter.Add(result.handleConnectionState)
	_ = result.messageRouter.Add(result.handleConnectionCreated)
	_ = result.messageRouter.Add(result.handleConnectionClosed)
	_ = result.messageRouter.Add(result.handlePublishInstanceDataFor)

	return result, nil
}

func (self *data) Send(data interface{}) error {
	self.messageRouter.Route(data)
	return nil
}

func (self *data) handlePublishInstanceDataFor(message *publishInstanceDataFor) error {
	if info, ok := self.ConnectionDataMap[message.Id]; ok {
		self.DoConnectionInstanceChange(info)
	}
	return nil
}

func (self *data) handleEmptyQueue(*messages.EmptyQueue) {
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
		dataInstance.update(message.Grid, message.KeyValue)
	}
}

func (self *data) handleConnectionClosed(message *model.ConnectionClosed) error {
	delete(self.ConnectionDataMap, message.ConnectionId)
	self.connectionListIsDirty = true
	return nil
}

func (self *data) handleConnectionCreated(message *model.ConnectionCreated) error {
	self.ConnectionDataMap[message.ConnectionId] = NewConnectionInstanceData(
		message.ConnectionId,
		true,
		message.CancelContext,
		message.CancelFunc,
		message.ConnectionName,
		message.ConnectionTime,
		message.NextFuncOutBoundChannel,
		message.NextFuncInBoundChannel,
	)

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
