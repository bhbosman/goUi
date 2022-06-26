package connectionManagerSlide

import (
	"context"
	"github.com/bhbosman/gocommon/model"
	"time"
)

type ConnectionInstanceData struct {
	ConnectionId   string
	isDirty        bool
	CancelContext  context.Context
	CancelFunc     context.CancelFunc
	Name           string
	ConnectionTime time.Time
	Grid           []model.LineData
}
