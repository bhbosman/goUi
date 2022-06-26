package ConnectionSlide

import (
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
)

type IConnectionSlide interface {
	ISendMessage.ISendMessage
}

type IConnectionSlideData interface {
	IConnectionSlide
}

type IConnectionSlideService interface {
	IConnectionSlide
	IFxService.IFxServices
}
