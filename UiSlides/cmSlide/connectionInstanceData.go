package cmSlide

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
	KeyValue       []model.KeyValue
}

func (self *ConnectionInstanceData) update(
	//CancelContext context.Context,
	//CancelFunc context.CancelFunc,
	//Name string,
	//ConnectionTime time.Time,
	Grid []model.LineData,
	KeyValue []model.KeyValue,
) {
	self.isDirty = true
	//self.CancelContext = CancelContext
	//self.CancelFunc = CancelFunc
	//self.Name = Name
	//self.ConnectionTime = ConnectionTime
	self.Grid = Grid
	self.KeyValue = KeyValue

}

func NewConnectionInstanceData(connectionId string, isDirty bool, cancelContext context.Context, cancelFunc context.CancelFunc, name string, connectionTime time.Time, grid []model.LineData, keyValue []model.KeyValue) *ConnectionInstanceData {
	return &ConnectionInstanceData{
		ConnectionId:   connectionId,
		isDirty:        isDirty,
		CancelContext:  cancelContext,
		CancelFunc:     cancelFunc,
		Name:           name,
		ConnectionTime: connectionTime,
		Grid:           grid,
		KeyValue:       keyValue,
	}
}
