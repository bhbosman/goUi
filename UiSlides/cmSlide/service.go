package cmSlide

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/pubSub"
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
)

type Service struct {
	onConnectionListChange     func(connectionList []IdAndName)
	onConnectionInstanceChange func(data ConnectionInstanceData)
	onData                     func() (IConnectionSlideData, error)
	state                      IFxService.State
	ctx                        context.Context
	cancelFunc                 context.CancelFunc
	cmdChannel                 chan interface{}
	pubSub                     *pubsub.PubSub
	ConnectionManagerHelper    goConnectionManager.IHelper
	ConnectionManager          goConnectionManager.IService
	subscribeChannel           *pubsub.NextFuncSubscription

	UniqueReferenceService interfaces.IUniqueReferenceService
	logger                 *zap.Logger
	goFunctionCounter      GoFunctionCounter.IService
}

func (self *Service) SetConnectionInstanceChange(cb func(data ConnectionInstanceData)) {
	self.onConnectionInstanceChange = cb
}

func (self *Service) SetConnectionListChange(cb func(connectionList []IdAndName)) {
	self.onConnectionListChange = cb
}

func NewService(
	parentContext context.Context,
	pubSub *pubsub.PubSub,
	onData func() (IConnectionSlideData, error),
	ConnectionManagerHelper goConnectionManager.IHelper,
	UniqueReferenceService interfaces.IUniqueReferenceService,
	logger *zap.Logger,
	goFunctionCounter GoFunctionCounter.IService,
	ConnectionManager goConnectionManager.IService,
) (*Service, error) {
	ctx, cancelFunc := context.WithCancel(parentContext)
	channel := make(chan interface{}, 32)

	return &Service{
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

func (self *Service) Send(message interface{}) error {
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

func (self *Service) OnStart(ctx context.Context) error {
	err := self.start(ctx)
	if err != nil {
		return err
	}
	self.state = IFxService.Started
	return nil
}

func (self *Service) OnStop(ctx context.Context) error {
	err := self.shutdown(ctx)
	self.cancelFunc()
	close(self.cmdChannel)
	self.state = IFxService.Stopped
	return err
}

func (self *Service) shutdown(_ context.Context) error {
	self.cancelFunc()
	return pubSub.Unsubscribe("FMD Manager Service", self.pubSub, self.goFunctionCounter, self.subscribeChannel)
}

func (self *Service) start(_ context.Context) error {
	data, err := self.onData()
	data.SetConnectionListChange(self.onConnectionListChange)
	data.SetConnectionInstanceChange(self.onConnectionInstanceChange)
	if err != nil {
		return err
	}
	go self.goStart(data)
	return nil
}

func (self *Service) State() IFxService.State {
	return self.state
}

func (self *Service) ServiceName() string {
	return "ConnectionSlideService"
}

func (self *Service) goStart(data IConnectionSlideData) {
	self.subscribeChannel = pubsub.NewNextFuncSubscription(goCommsDefinitions.CreateNextFunc(self.cmdChannel))
	self.pubSub.AddSub(self.subscribeChannel, self.ConnectionManagerHelper.PublishChannelName())

	ss := self.UniqueReferenceService.Next("ConnectionManagerReceiver")
	refreshSubChannel := pubsub.NewChannelSubscription(32)

	self.pubSub.AddSub(refreshSubChannel, ss)
	go func(refreshSubChannel *pubsub.ChannelSubscription) {
	loop:
		for {
			select {
			case unk, ok := <-refreshSubChannel.Data:
				if !ok {
					break loop
				}
				switch v := unk.(type) {
				default:
					_ = self.Send(v)
					break
				}
			}
		}
	}(refreshSubChannel)
	_ = self.ConnectionManager.Send(
		&goConnectionManager.RefreshDataTo{},
	)

	var messageReceived interface{}
	var ok bool

	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.ctx,
		data,
		[]ChannelHandler.ChannelHandler{
			{
				//BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					rr, e := ISendMessage.ChannelEventsForISendMessage(next.(ISendMessage.ISendMessage), message)
					return rr, e
				},
			},
			{
				//BreakOnSuccess: false,
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
			err := data.ShutDown()
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
