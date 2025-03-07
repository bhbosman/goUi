package cmSlide

import (
	"context"
	"github.com/bhbosman/gocommon/model"
	"github.com/reactivex/rxgo/v2"
	"time"
)

type ConnectionInstanceData struct {
	ConnectionId            string
	isDirty                 bool
	CancelContext           context.Context
	CancelFunc              context.CancelFunc
	Name                    string
	ConnectionTime          time.Time
	Grid                    []model.LineData
	KeyValue                []model.KeyValue
	NextFuncOutBoundChannel rxgo.NextFunc
	NextFuncInBoundChannel  rxgo.NextFunc
}

func (self *ConnectionInstanceData) update(
	Grid []model.LineData,
	KeyValue []model.KeyValue,
) {
	self.isDirty = true
	self.Grid = Grid
	self.KeyValue = KeyValue

}

func (self *ConnectionInstanceData) ResetConnectionParams() {
	if self.NextFuncInBoundChannel != nil {
		self.NextFuncInBoundChannel(model.NewClearCounters())
	}
	if self.NextFuncOutBoundChannel != nil {
		self.NextFuncOutBoundChannel(model.NewClearCounters())
	}
}

func NewConnectionInstanceData(
	connectionId string,
	isDirty bool,
	cancelContext context.Context,
	cancelFunc context.CancelFunc,
	name string,
	connectionTime time.Time,
	nextFuncOutBoundChannel rxgo.NextFunc,
	nextFuncInBoundChannel rxgo.NextFunc,

) *ConnectionInstanceData {
	return &ConnectionInstanceData{
		ConnectionId:            connectionId,
		isDirty:                 isDirty,
		CancelContext:           cancelContext,
		CancelFunc:              cancelFunc,
		Name:                    name,
		ConnectionTime:          connectionTime,
		NextFuncOutBoundChannel: nextFuncOutBoundChannel,
		NextFuncInBoundChannel:  nextFuncInBoundChannel,
	}
}
