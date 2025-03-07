package cmSlide

import "github.com/bhbosman/gocommon/services/IDataShutDown"

type IConnectionSlideData interface {
	IConnectionSlide
	IDataShutDown.IDataShutDown
}
