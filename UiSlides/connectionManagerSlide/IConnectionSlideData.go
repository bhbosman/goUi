package connectionManagerSlide

import "github.com/bhbosman/gocommon/Services/IDataShutDown"

type IConnectionSlideData interface {
	IConnectionSlide
	IDataShutDown.IDataShutDown
}
