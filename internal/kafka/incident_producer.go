package kafka

import (
	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
}

func NewProducer(brokerList []string) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokerList...),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, err
	}
	return &Producer{
		client: client,
	}, nil
}

func (p *Producer) ProduceMessage(event KafkaEvent) error {

	return nil
}
