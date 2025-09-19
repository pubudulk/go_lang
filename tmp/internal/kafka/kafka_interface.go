package kafka

type KafkaEvent interface {
	GetTopic() string
	GetVersion() int
	GetPayload() ([]byte, error)
}
