package connectionManagerSlide

import "github.com/bhbosman/gocommon/Services/ISendMessage"

type IConnectionSlide interface {
	SetConnectionInstanceChange(cb func(data *ConnectionInstanceData))
	SetConnectionListChange(cb func(connectionList []IdAndName))
	ISendMessage.ISendMessage
}
