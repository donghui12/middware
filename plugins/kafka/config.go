package kafka

type Config struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
	Offset  int64    `yaml:"offset"`
}
