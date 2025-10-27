package startup

type BaseConsumer struct {
	Ready   chan struct{}
	Brokers []string
	Group   string
}
