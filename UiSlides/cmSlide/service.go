package cmSlide

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/pubSub"
	"github.com/bhbosman/gocommon/services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"github.com/bhbosman/gocommon/services/interfaces"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
)

type service struct {
	onData                  func() (IConnectionSlideData, error)
	state                   IFxService.State
	ctx                     context.Context
	cancelFunc              context.CancelFunc
	cmdChannel              chan interface{}
	pubSub                  *pubsub.PubSub
	ConnectionManagerHelper goConnectionManager.IHelper
	ConnectionManager       goConnectionManager.IService
	subscribeChannel        *pubsub.NextFuncSubscription
	dataInstance            IConnectionSlideData
	UniqueReferenceService  interfaces.IUniqueReferenceService
	logger                  *zap.Logger
	goFunctionCounter       GoFunctionCounter.IService
}

func (self *service) ResetConnectionParams(connectionId string) {
	_, _ = CallIConnectionSlideResetConnectionParams(self.ctx, self.cmdChannel, false, connectionId)
}

func (self *service) ResetAllConnectionParams() {
	_, _ = CallIConnectionSlideResetAllConnectionParams(self.ctx, self.cmdChannel, false)
}

func (self *service) DisconnectAllConnections() {
	_, _ = CallIConnectionSlideDisconnectAllConnections(self.ctx, self.cmdChannel, false)
}

func (self *service) DisconnectConnection(connectionId string) {
	_, _ = CallIConnectionSlideDisconnectConnection(self.ctx, self.cmdChannel, false, connectionId)
}

func (self *service) SetConnectionInstanceChange(cb func(data ConnectionInstanceData)) {
	_, _ = CallIConnectionSlideSetConnectionInstanceChange(self.ctx, self.cmdChannel, false, cb)
}

func (self *service) SetConnectionListChange(cb func(connectionList []IdAndName)) {
	_, _ = CallIConnectionSlideSetConnectionListChange(self.ctx, self.cmdChannel, false, cb)
}

func newService(
	parentContext context.Context,
	pubSub *pubsub.PubSub,
	onData func() (IConnectionSlideData, error),
	ConnectionManagerHelper goConnectionManager.IHelper,
	UniqueReferenceService interfaces.IUniqueReferenceService,
	logger *zap.Logger,
	goFunctionCounter GoFunctionCounter.IService,
	ConnectionManager goConnectionManager.IService,
) (IConnectionSlideService, error) {
	ctx, cancelFunc := context.WithCancel(parentContext)
	channel := make(chan interface{}, 32)

	return &service{
		onData:                  onData,
		ctx:                     ctx,
		cancelFunc:              cancelFunc,
		cmdChannel:              channel,
		pubSub:                  pubSub,
		ConnectionManagerHelper: ConnectionManagerHelper,
		UniqueReferenceService:  UniqueReferenceService,
		logger:                  logger,
		goFunctionCounter:       goFunctionCounter,
		ConnectionManager:       ConnectionManager,
	}, nil
}

func (self *service) Send(message interface{}) error {
	send, err := ISendMessage.CallISendMessageSend(
		self.ctx,
		self.cmdChannel,
		false, // Todo: need to figure out why this is false. can not remember why
		message)
	if err != nil {
		return err
	}
	return send.Args0
}

func (self *service) OnStart(ctx context.Context) error {
	err := self.start(ctx)
	if err != nil {
		return err
	}
	self.state = IFxService.Started
	return nil
}

func (self *service) OnStop(ctx context.Context) error {
	err := self.shutdown(ctx)
	self.cancelFunc()
	close(self.cmdChannel)
	self.state = IFxService.Stopped
	return err
}

func (self *service) shutdown(_ context.Context) error {
	self.cancelFunc()
	return pubSub.Unsubscribe("FMD Manager Service", self.pubSub, self.goFunctionCounter, self.subscribeChannel)
}

func (self *service) start(_ context.Context) error {
	var err error
	self.dataInstance, err = self.onData()
	if err != nil {
		return err
	}
	go self.goStart(self.dataInstance)
	return nil
}

func (self *service) State() IFxService.State {
	return self.state
}

func (self *service) ServiceName() string {
	return "ConnectionSlideService"
}

func (self *service) goStart(dataInstance IConnectionSlideData) {
	self.subscribeChannel = pubsub.NewNextFuncSubscription(
		goCommsDefinitions.CreateNextFunc(self.cmdChannel),
	)
	self.pubSub.AddSub(self.subscribeChannel, self.ConnectionManagerHelper.PublishChannelName())
	_ = self.ConnectionManager.Send(
		&goConnectionManager.RefreshDataTo{
			PubSubBag: self.subscribeChannel,
		},
	)

	var messageReceived interface{}
	var ok bool

	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.ctx,
		dataInstance,
		[]ChannelHandler.ChannelHandler{
			{
				Cb: func(next interface{}, message interface{}) (bool, error) {
					rr, e := ChannelEventsForIConnectionSlide(next.(IConnectionSlide), message)
					return rr, e
				},
			},

			{
				Cb: func(next interface{}, message interface{}) (bool, error) {
					rr, e := ISendMessage.ChannelEventsForISendMessage(next.(ISendMessage.ISendMessage), message)
					return rr, e
				},
			},
			{
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if sm, ok := next.(ISendMessage.ISendMessage); ok {
						_ = sm.Send(message)
					}
					return true, nil
				},
			},
		},
		func() int {
			return len(self.cmdChannel)
		},
		goCommsDefinitions.CreateTryNextFunc(self.cmdChannel),
	)
loop:
	for {
		select {
		case <-self.ctx.Done():
			err := dataInstance.ShutDown()
			if err != nil {
				self.logger.Error(
					"error on done",
					zap.Error(err))
			}
			break loop
		case messageReceived, ok = <-self.cmdChannel:
			if !ok {
				return
			}
			b, err := channelHandlerCallback(messageReceived)
			if err != nil || b {
				return
			}
		}
	}
}
