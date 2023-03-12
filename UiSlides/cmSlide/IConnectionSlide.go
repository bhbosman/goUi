package cmSlide

import "github.com/bhbosman/gocommon/services/ISendMessage"

type IConnectionSlide interface {
	SetConnectionInstanceChange(cb func(data ConnectionInstanceData))
	SetConnectionListChange(cb func(connectionList []IdAndName))
	ISendMessage.ISendMessage
	DisconnectAllConnections()
	DisconnectConnection(connectionId string)
	ResetConnectionParams(connectionId string)
	ResetAllConnectionParams()
}
