package ConnectionSlide

import "github.com/bhbosman/gocommon/Services/ISendMessage"

type IConnectionSlide interface {
	ISendMessage.ISendMessage
}

type IConnectionSlideData interface {
	IConnectionSlide
}
