package connectionManagerSlide

import (
	"github.com/bhbosman/gocommon/Services/IFxService"
)

type IConnectionSlideService interface {
	IConnectionSlide
	IFxService.IFxServices
}
