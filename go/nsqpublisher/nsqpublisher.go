package nsqpublisher

import (
	"log"

	"github.com/nsqio/go-nsq"
)

type Publisher struct {
	producer *nsq.Producer
}

func NewPublisher(address string) (*Publisher, error) {
	cfg := nsq.NewConfig()

	prod, err := nsq.NewProducer(address, cfg)
	log.Println("Address", address)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		producer: prod,
	}, nil
}

func (p *Publisher) Publish(topic string, body []byte) error {
	return p.producer.Publish(topic, body)
}

func (p *Publisher) Close() error {
	p.producer.Stop()
	return nil
}
