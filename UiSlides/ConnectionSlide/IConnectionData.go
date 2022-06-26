package ConnectionSlide

import (
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
)

type IConnectionSlide interface {
	SetConnectionInstanceChange(cb func(data *ConnectionInstanceData))
	SetConnectionListChange(cb func(connectionList []IdAndName))
	ISendMessage.ISendMessage
}

type IConnectionSlideData interface {
	IConnectionSlide
}

type IConnectionSlideService interface {
	IConnectionSlide
	IFxService.IFxServices
}
