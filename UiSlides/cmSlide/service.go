package cmSlide

import (
	"context"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
	"github.com/bhbosman/gocommon/Services/interfaces"
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
	UniqueReferenceService     interfaces.IUniqueReferenceService
	logger                     *zap.Logger
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
	return nil
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
	defer func(cmdChannel <-chan interface{}) {
		//flush
		for range cmdChannel {
		}
	}(self.cmdChannel)

	pubSubChannel := self.pubSub.Sub(self.ConnectionManagerHelper.PublishChannelName())
	defer func(pubSubChannel chan interface{}) {
		// unsubscribe on different go routine to avoid deadlock
		go func(pubSubChannel chan interface{}) {
			self.pubSub.Unsub(pubSubChannel)
			for range pubSubChannel {
			}
		}(pubSubChannel)
	}(pubSubChannel)

	ss := self.UniqueReferenceService.Next("ConnectionManagerReceiver")
	refreshSubChannel := self.pubSub.Sub(ss)
	go func(refreshSubChannel chan interface{}) {
		for m := range refreshSubChannel {
			switch v := m.(type) {
			case *goConnectionManager.RefreshDataStart:
				_ = self.Send(v)
				break
			case *goConnectionManager.RefreshDataStop:
				_ = self.Send(v)
				self.pubSub.Unsub(refreshSubChannel, ss)
				break
			default:
				_ = self.Send(v)
				break
			}
		}
	}(refreshSubChannel)
	self.ConnectionManagerHelper.RefreshData(ss)

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
			return len(pubSubChannel) + len(self.cmdChannel)
		},
		//func(i interface{}) {
		//	select {
		//	case self.cmdChannel <- i:
		//		break
		//	default:
		//		break
		//	}
		//},
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
		case messageReceived, ok = <-pubSubChannel:
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
