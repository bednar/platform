package nats

import (
	"errors"
	"io"
	"io/ioutil"

	stan "github.com/nats-io/go-nats-streaming"
	"go.uber.org/zap"
)

var ErrNoNatsConnection = errors.New("Nats connection has not been established. Call Open() first.")

type Publisher interface {
	// Open creates a connection to the nats server
	Open() error
	// Similar to Publish in nats
	Publish(subject string, r io.Reader) error
}

type publisher struct {
	ClientID   string
	Connection stan.Conn
	Logger     *zap.Logger
}

func NewPublisher(clientID string) *publisher {
	p := publisher{ClientID: clientID}
	return &p
}

func (p *publisher) Open() error {
	sc, err := stan.Connect(ServerName, p.ClientID)
	if err != nil {
		return err
	}
	p.Connection = sc
	return nil
}

func (p *publisher) Publish(subject string, r io.Reader) error {
	if p.Connection == nil {
		return ErrNoNatsConnection
	}

	ah := func(guid string, err error) {
		if err != nil {
			p.Logger.Info(err.Error())
		}
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	_, err = p.Connection.PublishAsync(subject, data, ah)
	return err
}
