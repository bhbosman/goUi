package cmSlide

import (
	"github.com/bhbosman/gocommon/services/IFxService"
)

type IConnectionSlideService interface {
	IConnectionSlide
	IFxService.IFxServices
}
