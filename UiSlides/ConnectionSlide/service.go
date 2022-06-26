package ConnectionSlide

import (
	"context"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
	"github.com/cskr/pubsub"
)

type Service struct {
	onConnectionListChange     func(connectionList []IdAndName)
	onConnectionInstanceChange func(data *ConnectionInstanceData)
	onData                     func() (IConnectionSlideData, error)
	state                      IFxService.State
	ctx                        context.Context
	cancelFunc                 context.CancelFunc
	cmdChannel                 chan interface{}
	pubSub                     *pubsub.PubSub
}

func (self *Service) SetConnectionInstanceChange(cb func(data *ConnectionInstanceData)) {
	self.onConnectionInstanceChange = cb
}

func (self *Service) SetConnectionListChange(cb func(connectionList []IdAndName)) {
	self.onConnectionListChange = cb
}

func NewService(
	parentContext context.Context,
	pubSub *pubsub.PubSub,
	onData func() (IConnectionSlideData, error),
) (*Service, error) {
	ctx, cancelFunc := context.WithCancel(parentContext)
	channel := make(chan interface{}, 32)

	return &Service{
		onData:     onData,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		cmdChannel: channel,
		pubSub:     pubSub,
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
	data, err := IDataShutDown.CallIDataShutDownShutDown(self.ctx, self.cmdChannel, true)
	if err != nil {
		return err
	}
	self.cancelFunc()
	return data.Args0
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

	pubSubChannel := self.pubSub.Sub("ActiveConnectionStatus")
	defer func(pubSubChannel chan interface{}) {
		// unsubscribe on different go routine to avoid deadlock
		go func(pubSubChannel chan interface{}) {
			self.pubSub.Unsub(pubSubChannel)
			for range pubSubChannel {
			}
		}(pubSubChannel)
	}(pubSubChannel)

	var messageReceived interface{}
	var ok bool

	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.ctx,
		data,
		[]ChannelHandler.ChannelHandler{
			{
				PubSubHandler:  false,
				BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					rr, e := ISendMessage.ChannelEventsForISendMessage(next.(ISendMessage.ISendMessage), message)
					return rr, e
				},
			},
			{
				PubSubHandler:  true,
				BreakOnSuccess: false,
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
		})
loop:
	for self.ctx.Err() == nil {
		select {
		case <-self.ctx.Done():
			break loop
		case messageReceived, ok = <-self.cmdChannel:
			if !ok {
				return
			}
			b, err := channelHandlerCallback(messageReceived, false)
			if err != nil || b {
				return
			}
		case messageReceived, ok = <-pubSubChannel:
			if !ok {
				return
			}
			b, err := channelHandlerCallback(messageReceived, true)
			if err != nil || b {
				return
			}
		}
	}
}
