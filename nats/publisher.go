package nats

import (
	"io"
	"io/ioutil"

	stan "github.com/nats-io/go-nats-streaming"
	"go.uber.org/zap"
)

type Publisher interface {
	// Similar to Publish in nats
	Publish(subject string, r io.Reader) error
}

type publisher struct {
	Connection stan.Conn
	Logger     *zap.Logger
}

func NewPublisher(clientID string) (Publisher, error) {
	sc, err := stan.Connect(ServerName, clientID)
	if err != nil {
		return nil, err
	}

	return &publisher{Connection: sc}, nil
}

func (p *publisher) Publish(subject string, r io.Reader) error {
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
