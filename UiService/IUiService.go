package UiService

import (
	"github.com/bhbosman/goUi/ui"
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"github.com/rivo/tview"
)

type IUi interface {
	Build(*tview.Application, ...ui.ISlideFactory) (ui.IPrimitiveCloser, error)
}

type IUiService interface {
	IUi
	IFxService.IFxServices
}

type IUiData interface {
	IUi
	IDataShutDown.IDataShutDown
	ISendMessage.ISendMessage
}
