package nats

import (
	stan "github.com/nats-io/go-nats-streaming"
)

type Subscriber interface {
	// Open creates a connection to the nats server
	Open() error
	Subscribe(subject, group string, handler Handler) error
}

type subscriber struct {
	ClientID   string
	Connection stan.Conn
}

func NewSubscriber(clientID string) *subscriber {
	return &subscriber{ClientID: clientID}
}

func (s *subscriber) Open() error {
	sc, err := stan.Connect(ServerName, s.ClientID)
	if err != nil {
		return err
	}
	s.Connection = sc
	return nil
}

type messageHandler struct {
	handler Handler
	sub     subscription
}

func (mh *messageHandler) handle(m *stan.Msg) {
	mh.handler.Process(mh.sub, &message{m: m})
}

func (s *subscriber) Subscribe(subject, group string, handler Handler) error {
	if s.Connection == nil {
		return ErrNoNatsConnection
	}

	mh := messageHandler{handler: handler}
	sub, err := s.Connection.QueueSubscribe(subject, group, mh.handle, stan.DurableName(group), stan.SetManualAckMode(), stan.MaxInflight(25))
	if err != nil {
		return err
	}
	mh.sub = subscription{sub: sub}
	return nil
}
